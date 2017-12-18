package lib

import (
	"io"
	"os/exec"
	"time"
)

// Execution represents the instantiation of a command
// whose execution can be limited and tracked.
type Execution struct {
	Argv      []string
	ExitCode  int
	Env       map[string]string
	Directory string
	Stdout    io.Writer
	Stderr    io.Writer
	StartTime time.Time
	EndTime   time.Time

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

	// OnJobStatusChange is a callback function to be called
	// once per transition of job status.
	OnJobStatusChange func(a *Activity) `yaml:"-"`
}

// Runtime aggragates CLI and runtime configuration
// to be applied when running `cr`
type Runtime struct {
	// File denotes the path to the configuration file to
	// load
	File string `arg:"help:path the configuration file" yaml:"File"`

	// LogsDirectory indicates the path to the directory where logs
	// are sent to.
	LogsDirectory string `arg:"help:path to the directory where logs are sent to" yaml:"LogsDirectory"`

	// Stdout indicates whether the execution logs should be pipped
	// to stdout or not.
	Stdout bool `arg:"help:log executions to stdout" yaml:"Stdout"`

	// Graph indicates whether a dot graph should be output
	// or not.
	Graph bool `arg:"help:output the execution graph" yaml:"Graph"`

	// Directory denotes what to used as a current working directory
	// for the executions when a relative path is indicated in the
	// job description.
	Directory string `arg:"help:directory to be used as current working directory" yaml:"Directory"`
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

	// Directory names the absolute or relative path
	// to get into before executin the command.
	// By default it takes the value "." (current working
	// directory).
	Directory string `yaml:"Directory"`

	// Whether the output of the execution should
	// be stored or not.
	CaptureOutput bool `yaml:"CaptureOutput"`

	// Env stores the extra environment to add to the
	// command execution.
	Env map[string]string `yaml:"Env"`

	// StartTime is the timestamp at the moment of
	// the initiation of the execution of the
	// command.
	StartTime *time.Time `yaml:"-"`

	// EndTime is the timestamp of the end execution
	// of the command.
	EndTime *time.Time `yaml:"-"`

	// ExitCode stores the result exit-code of the command.
	ExitCode int `yaml:"-"`

	// Output is the output captured once the command
	// has been executed.
	Output string `yaml:"-"`

	// DependsOn lists a series of jobs that the job depends
	// on to start its execution.
	DependsOn []string `yaml:"DependsOn,flow"`

	// LogFilepath indicates the path to the file where the logs
	// of the job execution are sent to.
	LogFilepath string `yaml:"LogFilepath"`
}

func (j Job) Name() string {
	return j.Id
}
