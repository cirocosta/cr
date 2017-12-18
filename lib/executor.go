package lib

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Executor encapsulates the execution
// context of a graph of jobs.
type Executor struct {
	config        *Config
	graph         *dag.AcyclicGraph
	logger        zerolog.Logger
	jobsMap       map[string]*Job
	logsDirectory string
}

// New instantiates a new Executor from
// the supplied configuration.
func New(cfg *Config) (e Executor, err error) {
	var finfo os.FileInfo

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

	if cfg.Runtime.LogsDirectory == "" {
		err = errors.Errorf("LogsDirectory must be specified")
		return
	}

	finfo, err = os.Stat(cfg.Runtime.LogsDirectory)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to look for logs directory %s",
			cfg.Runtime.LogsDirectory)
		return
	}

	if !finfo.IsDir() {
		err = errors.Errorf(
			"logs directory must be a directory %s",
			cfg.Runtime.LogsDirectory)
		return
	}

	e.logsDirectory = cfg.Runtime.LogsDirectory
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

func (e *Executor) ResolveJobDirectory(j *Job, renderState *RenderState) (res string, err error) {
	if j == nil || renderState == nil {
		err = errors.Errorf("job and renderState must be non-nil")
		return
	}

	switch j.Directory {
	case "":
		res = "."
	default:
		res, err = TemplateField(j.Directory, renderState)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't render Directory string")
			return
		}
	}

	return
}

func (e *Executor) ResolveJobLogFilepath(j *Job, renderState *RenderState) (res string, err error) {
	if j == nil || renderState == nil {
		err = errors.Errorf("job and renderState must be non-nil")
		return
	}

	switch j.LogFilepath {
	case "":
		res = path.Join(
			e.config.Runtime.LogsDirectory,
			j.Id)
	default:
		res, err = TemplateField(j.LogFilepath, renderState)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't render LogFilepath string")
			return
		}
	}

	return
}

func (e *Executor) ResolveJobRun(j *Job, renderState *RenderState) (res string, err error) {
	if j == nil || renderState == nil {
		err = errors.Errorf("job and renderState must be non-nil")
		return
	}

	switch j.Run {
	case "":
		res = ""
	default:
		res, err = TemplateField(j.Run, renderState)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't render run command")
			return
		}
	}

	return
}

func (e *Executor) ResolveJobEnv(j *Job, renderState *RenderState) (res map[string]string, err error) {
	res = map[string]string{}

	if j == nil || renderState == nil {
		err = errors.Errorf("job and renderState must be non-nil")
		return
	}

	var templateRes string
	for k, v := range e.config.Env {
		templateRes, err = TemplateField(v, renderState)
		if err != nil {
			err = errors.Errorf(
				"failed to template environment variable %s", k)
			return
		}

		res[k] = templateRes
	}

	for k, v := range j.Env {
		templateRes, err = TemplateField(v, renderState)
		if err != nil {
			err = errors.Errorf(
				"failed to template environment variable %s", k)
			return
		}

		res[k] = templateRes
	}

	return
}

// RunJob is a method invoked for each vertex
// in the execution graph except the root.
// TODO Split into a job preparation step and a
// job execution step.
func (e *Executor) RunJob(ctx context.Context, j *Job) (err error) {
	var (
		execution *Execution
		logFile   *os.File
		output    bytes.Buffer

		stdout      = []io.Writer{}
		stderr      = []io.Writer{}
		renderState = &RenderState{
			Jobs: e.jobsMap,
		}
	)

	if j.CaptureOutput {
		stdout = append(stdout, &output)
	}

	if e.config.Runtime.Stdout {
		stdout = append(stdout, os.Stdout)
		stderr = append(stdout, os.Stderr)
	}

	j.LogFilepath, err = e.ResolveJobLogFilepath(j, renderState)
	if err != nil {
		return
	}

	logFile, err = os.Create(j.LogFilepath)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create file for logging %s",
			j.LogFilepath)
		return
	}
	defer logFile.Close()

	stdout = append(stdout, logFile)
	stderr = append(stderr, logFile)

	j.Directory, err = e.ResolveJobDirectory(j, renderState)
	if err != nil {
		return
	}

	j.Env, err = e.ResolveJobEnv(j, renderState)
	if err != nil {
		return
	}

	j.Run, err = e.ResolveJobRun(j, renderState)
	if err != nil {
		return
	}

	if j.Run == "" {
		goto END
	}

	execution = &Execution{
		Argv: []string{
			"/bin/bash",
			"-c",
			j.Run,
		},
		Stdout:    io.MultiWriter(stdout...),
		Stderr:    io.MultiWriter(stderr...),
		Directory: j.Directory,
		Env:       j.Env,
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
