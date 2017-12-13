package lib

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Executor encapsulates the execution
// context of a graph of jobs.
type Executor struct {
	config  *Config
	graph   *dag.AcyclicGraph
	logger  zerolog.Logger
	jobsMap map[string]*Job
}

// New instantiates a new Executor from
// the supplied configuration.
func New(cfg *Config) (e Executor, err error) {
	if cfg == nil {
		err = errors.Errorf("cfg must be non-nill")
		return
	}

	graph, err := BuildDependencyGraph(cfg.Jobs)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create dependency graph")
		return
	}

	e.config = cfg
	e.graph = &graph
	e.jobsMap = map[string]*Job{}
	e.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "executor").
		Logger()

	for _, job := range cfg.Jobs {
		e.jobsMap[job.Id] = job
	}

	return
}

// GetDotGraph retrieves a `dot` visualization of
// the dependency graph.
func (e *Executor) GetDotGraph() (res string) {
	res = string(e.graph.Dot(&dag.DotOpts{}))
	return
}

// Execute initiates the parallel execution of the
// jobs.
func (e *Executor) Execute(ctx context.Context) (err error) {
	err = e.TraverseAndExecute(ctx, e.graph)
	if err != nil {
		err = errors.Wrapf(err, "jobs execution failed")
		return
	}

	return
}

// RunJob is a method invoked for each vertex
// in the execution graph except the root.
// TODO Split into a job preparation step and a
// job execution step.
func (e *Executor) RunJob(ctx context.Context, j *Job) (err error) {
	var (
		renderState = &RenderState{
			Jobs: e.jobsMap,
		}
		execution *Execution
		output    bytes.Buffer
		stdout    io.Writer = os.Stdout
		stderr    io.Writer = os.Stderr
		run       string
		directory string
	)

	if j.CaptureOutput {
		stdout = io.MultiWriter(&output, os.Stdout)
		stderr = io.MultiWriter(&output, os.Stderr)
	}

	switch j.Directory {
	case "":
		directory = "."
	default:
		directory, err = TemplateField(j.Directory, renderState)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't render directory string")
			return
		}
	}

	j.Directory = directory

	switch j.Run {
	case "":
		goto END
	default:
		run, err = TemplateField(j.Run, renderState)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't render run command")
			return
		}
	}

	j.Run = run

	execution = &Execution{
		Argv: []string{
			"/bin/bash",
			"-c",
			j.Run,
		},
		Stdout:    stdout,
		Stderr:    stderr,
		Directory: j.Directory,
	}

	if e.config.OnJobStatusChange != nil {
		e.config.OnJobStatusChange(&Activity{
			Type: ActivityStarted,
			Time: time.Now(),
			Job:  j,
		})
	}

	err = execution.Run(ctx)

	j.StartTime = &execution.StartTime
	j.EndTime = &execution.EndTime

	if err != nil {
		err = errors.Wrapf(err, "command execution failed")

		if e.config.OnJobStatusChange != nil {
			e.config.OnJobStatusChange(&Activity{
				Type: ActivityErrored,
				Time: time.Now(),
				Job:  j,
			})
		}

		return
	}

	j.Output = strings.TrimSpace(output.String())

END:
	if e.config.OnJobStatusChange != nil {
		e.config.OnJobStatusChange(&Activity{
			Type: ActivitySuccess,
			Time: time.Now(),
			Job:  j,
		})
	}

	return
}

func (e *Executor) CreateWalkFunc(ctx context.Context) dag.WalkFunc {
	return func(v dag.Vertex) error {
		job, ok := v.(*Job)
		if !ok {
			return errors.Errorf("vertex not a job")
		}

		if job.Id == "_root" {
			return nil
		}

		return e.RunJob(ctx, job)
	}
}

// TraverseAndExecute goes through the graph
// provided and starts the execution of the jobs.
func (e *Executor) TraverseAndExecute(ctx context.Context, g *dag.AcyclicGraph) (err error) {
	w := &dag.Walker{
		Callback: e.CreateWalkFunc(ctx),
	}

	w.Update(g)

	err = w.Wait()
	if err != nil {
		err = errors.Wrapf(err,
			"execution of jobs failed")
		return
	}

	return
}
