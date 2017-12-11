package lib

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
)

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
