package lib

import (
	"github.com/hashicorp/terraform/dag"
	"github.com/pkg/errors"
)

// BuildDependencyGraph creates a DAG with
// a dumb root aiming at providing the
// biggest possible parallelism to a series
// of job builds.
func BuildDependencyGraph(jobs []*Job) (g dag.AcyclicGraph, err error) {
	var (
		rootJob = &Job{
			Id: "_root",
		}
		jobsMap = map[string]*Job{
			"_root": rootJob,
		}
		job *Job
		dep string
	)

	if jobs == nil {
		err = errors.Errorf("jobs can't be nil")
		return
	}

	g.Add(rootJob)

	for _, job = range jobs {
		if job.Id == "" {
			err = errors.Errorf("job must have name")
			return
		}

		g.Add(job)
		jobsMap[job.Id] = job
	}

	for _, job = range jobs {
		if len(job.DependsOn) == 0 {
			g.Connect(dag.BasicEdge(rootJob, job))
			continue
		}

		for _, dep = range job.DependsOn {
			depJob, present := jobsMap[dep]
			if !present {
				err = errors.Errorf(
					"job %s has a dependency %s "+
						"that does not exist",
					job.Id, dep)
				return
			}

			g.Connect(dag.BasicEdge(depJob, job))
		}
	}

	_, err = g.Root()
	if err != nil {
		err = errors.Wrapf(err, "couldn't compute DAG root")
		return
	}

	if len(g.Cycles()) > 0 {
		err = errors.Errorf("graph contains cycles")
		return
	}

	return
}
