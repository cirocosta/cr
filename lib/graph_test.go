package lib

import (
	"testing"

	"github.com/hashicorp/terraform/dag"
	"github.com/stretchr/testify/require"
)

func TestBuildDependencyGraph(t *testing.T) {
	var testCases = []struct{
		desc string
		jobs []*Job
		expected string
		shouldFail bool
	}{
		{
			desc: "nil should fail",
			shouldFail: true,
		},
	}

	var (
		err error
		graph dag.AcyclicGraph
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			graph, err = BuildDependencyGraph(tc.jobs)
			if tc.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
