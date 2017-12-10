package lib

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

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
	Jobs []Job `yaml:"Jobs"`
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

func (j Job) Name () string {
	return j.Id
}

func ConfigFromFile(file string) (config Config, err error) {
	finfo, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.Wrapf(err,
				"configuration file %s not found",
				file)
			return
		}

		err = errors.Wrapf(err,
			"unexpected error looking for config file %s",
			file)
		return
	}

	configContent, err := ioutil.ReadAll(finfo)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly read config file %s",
			file)
		return
	}

	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly parse yaml config file %s",
			file)
		return
	}

	return
}
