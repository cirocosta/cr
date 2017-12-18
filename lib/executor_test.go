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

func TestResolveEnvironment(t *testing.T) {
	var testCases = []struct {
		desc        string
		job         *Job
		globalEnv   map[string]string
		state       *RenderState
		expected    map[string]string
		shouldError bool
	}{
		{
			desc:        "nil",
			shouldError: true,
		},
		{
			desc:     "empty environment uses none",
			job:      &Job{},
			state:    &RenderState{},
			expected: map[string]string{},
		},
		{
			desc: "custom job environment",
			job: &Job{
				Env: map[string]string{"FOO": "BAR"},
			},
			state:    &RenderState{},
			expected: map[string]string{"FOO": "BAR"},
		},
		{
			desc: "custom job environment with global",
			job: &Job{
				Env: map[string]string{"FOO": "BAR"},
			},
			globalEnv: map[string]string{"CAZ": "BAZ"},
			state:     &RenderState{},
			expected:  map[string]string{"FOO": "BAR", "CAZ": "BAZ"},
		},
		{
			desc:      "just global",
			job:       &Job{},
			globalEnv: map[string]string{"CAZ": "BAZ"},
			state:     &RenderState{},
			expected:  map[string]string{"CAZ": "BAZ"},
		},
		{
			desc: "job overriding global",
			job: &Job{
				Env: map[string]string{"CAZ": "LOL"},
			},
			globalEnv: map[string]string{"CAZ": "BAZ"},
			state:     &RenderState{},
			expected:  map[string]string{"CAZ": "LOL"},
		},
		{
			desc: "custom job environment with templating",
			job: &Job{
				Env: map[string]string{"FOO": "{{ .Jobs.Job1.Output }}"},
			},
			state: &RenderState{
				Jobs: map[string]*Job{
					"Job1": {Output: "lol"},
				},
			},
			expected: map[string]string{"FOO": "lol"},
		},
		{
			desc:      "global environment with templating",
			globalEnv: map[string]string{"CAZ": "{{ .Jobs.Job1.Output }}"},
			job:       &Job{},
			state: &RenderState{
				Jobs: map[string]*Job{
					"Job1": {Output: "lol"},
				},
			},
			expected: map[string]string{"CAZ": "lol"},
		},
	}

	var (
		err    error
		actual map[string]string

		e = Executor{
			config: &Config{},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			e.config.Env = tc.globalEnv

			actual, err = e.ResolveJobEnv(tc.job, tc.state)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(actual))

			for k, _ := range tc.expected {
				assert.Equal(t, tc.expected[k], actual[k])
			}
		})
	}
}
