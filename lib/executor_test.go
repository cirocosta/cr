package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveJobDirectory(t *testing.T) {
	var testCases = []struct {
		desc        string
		job         *Job
		state       *RenderState
		expected    string
		shouldError bool
	}{
		{
			desc:        "nil",
			shouldError: true,
		},
		{
			desc: "empty directory uses default",
			job: &Job{
				Directory: "",
			},
			state:    &RenderState{},
			expected: ".",
		},
		{
			desc: "templated directory",
			job: &Job{
				Directory: "/{{ .Jobs.Job1.Output }}",
			},
			state: &RenderState{
				Jobs: map[string]*Job{
					"Job1": {Output: "lol"},
				},
			},
			expected: "/lol",
		},
	}

	var (
		err    error
		actual string

		e = Executor{}
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err = e.ResolveJobDirectory(tc.job, tc.state)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
