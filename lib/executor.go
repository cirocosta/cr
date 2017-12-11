package lib

import (
	"context"

	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
)

func Execute(ctx context.Context, j *Job) (err error) {
	var execution = &Execution{
		Argv: []string{
			"/bin/bash",
			"-c",
			j.Run,
		},
	}

	err = execution.Run(ctx)
	if err != nil {
		err = errors.Wrapf(err, "command execution failed")
		return
	}

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
func TraverseAndExecute(ctx context.Context, g dag.AcyclicGraph) (err error) {
	w := &dag.Walker{
		Callback: CreateExecutor(ctx),
	}

	err = w.Wait()
	if err != nil {
		err = errors.Wrapf(err,
			"execution of jobs failed")
		return
	}

	return
}
