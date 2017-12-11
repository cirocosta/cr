package lib

import (
	"bytes"
	"os/exec"
	"time"
)

// Execution represents the instantiation of a command
// whose execution can be limited and tracked.
type Execution struct {
	Argv      []string
	ExitCode  int
	Output    bytes.Buffer
	StartTime time.Time
	EndTime   time.Time
	Id        string

	cmd *exec.Cmd
}

// Config aggregates all the types of cofiguration
// that can be retrieved from a `.cr.yml` configuration
// file.
type Config struct {

	// Runtime contains the CLI and runtime configuration
	// to be applied when running `cr`.
	Runtime Runtime `yaml:"Runtime"`

	// Env defines environment variables that should
	// be applied to every execution
	Env map[string]string `yaml:"Env"`

	// Jobs lists the jobs to be executed.
	Jobs []*Job `yaml:"Jobs"`
}

// Runtime aggragates CLI and runtime configuration
// to be applied when running `cr`
type Runtime struct {
	File   string `arg:"help:path the configuration file" yaml:"File"`
	Stdout bool   `arg:"help:log executions to stdout" yaml:"Stdout"`
	Graph  bool   `arg:"help:output the execution graph" yaml:"Graph"`
}

// Job defines a unit of execution that at some point
// in time gets its command defined in `run` executed.
// It might happen to never be executed if a dependency
// is never met.
type Job struct {

	// Name is the name of the job being executed
	Id string `yaml:"Id"`

	// Run is a command to execute in the context
	// of a default shell.
	Run string `yaml:"Run"`

	// Whether the output of the execution should
	// be stored or not.
	CaptureOutput bool `yaml:"CaptureOutput"`

	// Env stores the extra environment to add to the
	// command execution.
	Env map[string]string `yaml:"Env"`

	// StartTime is the timestamp at the moment of
	// the initiation of the execution of the
	// command.
	StartTime time.Time `yaml:"-"`

	// FinishTime is the timestamp of the end execution
	// of the command.
	FinishTime time.Time `yaml:"-"`

	// ExitCode stores the result exit-code of the command.
	ExitCode int `yaml:"-"`

	// Output is the output captured once the command
	// has been executed.
	Output string `yaml:"-"`

	// DependsOn lists a series of jobs that the job depends
	// on to start its execution.
	DependsOn []string `yaml:"DependsOn,flow"`
}

func (j Job) Name() string {
	return j.Id
}
