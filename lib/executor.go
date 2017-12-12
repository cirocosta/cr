package lib

import (
	"bytes"
	"context"
	"io"
	"os"
	"text/template"

	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Executor encapsulates the execution
// context of a graph of jobs.
type Executor struct {
	config *Config
	graph  *dag.AcyclicGraph
	logger zerolog.Logger
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
	e.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "executor").
		Logger()

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
	err = TraverseAndExecute(ctx, e.graph)
	if err != nil {
		err = errors.Wrapf(err, "jobs execution failed")
		return
	}

	return
}

// RenderState encapsulates the state that can
// be used when templating a given field.
type RenderState struct {
	Jobs map[string]*Job
}

// TemplateField takes a field string and a state.
// With that it applies the state in the template and
// generates a response.
func TemplateField(field string, state *RenderState) (res string, err error) {
	var (
		tmpl   *template.Template
		output bytes.Buffer
	)

	tmpl, err = template.New("tmpl").Parse(field)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to instantiate template for record '%s'",
			field)
		return
	}

	err = tmpl.Execute(&output, state)
	if err != nil {
		err = errors.Wrapf(err, "failed to execute template")
		return
	}

	res = output.String()

	return
}

func Execute(ctx context.Context, j *Job) (err error) {
	var (
		execution *Execution
		output    bytes.Buffer
		stdout    io.Writer = os.Stdout
		stderr    io.Writer = os.Stderr
	)

	if j.CaptureOutput {
		stdout = io.MultiWriter(&output, os.Stdout)
		stderr = io.MultiWriter(&output, os.Stderr)
	}

	execution = &Execution{
		Argv: []string{
			"/bin/bash",
			"-c",
			j.Run,
		},
		Stdout: stdout,
		Stderr: stderr,
	}

	err = execution.Run(ctx)
	if err != nil {
		err = errors.Wrapf(err, "command execution failed")
		return
	}

	j.Output = output.String()

	return
}

func CreateExecutor(ctx context.Context) dag.WalkFunc {
	return func(v dag.Vertex) error {
		job, ok := v.(*Job)
		if !ok {
			return errors.Errorf("vertex not a job")
		}

		if job.Id == "_root" {
			return nil
		}

		return Execute(ctx, job)
	}
}

// TraverseAndExecute goes through the graph
// provided and starts the execution of the jobs.
func TraverseAndExecute(ctx context.Context, g *dag.AcyclicGraph) (err error) {
	w := &dag.Walker{
		Callback: CreateExecutor(ctx),
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
