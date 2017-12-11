package lib

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDependencyGraph(t *testing.T) {
	var testCases = []struct {
		desc       string
		jobs       []*Job
		expected   string
		shouldFail bool
	}{
		{
			desc:       "nil should fail",
			shouldFail: true,
		},
		{
			desc: "single job",
			jobs: []*Job{
				{
					Id: "job1",
				},
			},
			expected: `
_root
  job1
job1`,
		},
		{
			desc: "two jobs with no deps",
			jobs: []*Job{
				{
					Id: "job1",
				},
				{
					Id: "job2",
				},
			},
			expected: `
_root
  job1
  job2
job1
job2`,
		},
		{
			desc: "two jobs with single dependency",
			jobs: []*Job{
				{
					Id: "job1",
				},
				{
					Id: "job2",
					DependsOn: []string{
						"job1",
					},
				},
			},
			expected: `
_root
  job1
job1
  job2
job2`,
		},
		{
			desc: "thre jobs with two jobs depending in one",
			jobs: []*Job{
				{
					Id: "job1",
				},
				{
					Id: "job2",
					DependsOn: []string{
						"job1",
					},
				},
				{
					Id: "job3",
					DependsOn: []string{
						"job1",
					},
				},
			},
			expected: `
_root
  job1
job1
  job2
  job3
job2
job3`,
		},
		{
			desc: "three jobs with serial dependency",
			jobs: []*Job{
				{
					Id: "job1",
				},
				{
					Id: "job2",
					DependsOn: []string{
						"job1",
					},
				},
				{
					Id: "job3",
					DependsOn: []string{
						"job2",
					},
				},
			},
			expected: `
_root
  job1
job1
  job2
job2
  job3
job3`,
		},
		{
			desc: "cyclic dependency",
			jobs: []*Job{
				{
					Id: "job1",
					DependsOn: []string{
						"job2",
					},
				},
				{
					Id: "job2",
					DependsOn: []string{
						"job1",
					},
				},
			},
			shouldFail: true,
		},
	}

	var (
		err      error
		actual   string
		expected string
		graph    dag.AcyclicGraph
	)

	// TODO add a root dummy job

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			graph, err = BuildDependencyGraph(tc.jobs)
			if tc.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			actual = strings.Trim(graph.String(), "\n")
			expected = strings.Trim(tc.expected, "\n")

			assert.Equal(t, expected, actual)
		})
	}
}
