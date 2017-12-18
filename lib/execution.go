package lib

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultFailedExitCode int = 1
)

// init initializes the Execution parameters that rely on
// initialization. It also sets default parameters
// that might not be set.
func (e *Execution) init(ctx context.Context) (err error) {
	if len(e.Argv) == 0 {
		err = errors.Errorf("Argv must have at least one element")
		return
	}

	allEnv := os.Environ()
	for k, v := range e.Env {
		allEnv = append(allEnv, k+"="+v)
	}

	e.cmd = exec.CommandContext(ctx, e.Argv[0], e.Argv[1:]...)
	e.cmd.Stdout = e.Stdout
	e.cmd.Stderr = e.Stderr
	e.cmd.Dir = e.Directory
	e.cmd.Env = allEnv

	return
}

// Run is a blocking method that executes the desired command
// tying it to a context which, when cancelled, kills the process.
func (e *Execution) Run(ctx context.Context) (err error) {
	err = e.init(ctx)
	if err != nil {
		err = errors.Wrapf(err, "Couldn't initialize execution")
		return
	}

	e.StartTime = time.Now()
	err = e.cmd.Run()
	e.EndTime = time.Now()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			e.ExitCode = ws.ExitStatus()
		} else {
			e.ExitCode = defaultFailedExitCode
		}
	} else {
		ws := e.cmd.ProcessState.Sys().(syscall.WaitStatus)
		e.ExitCode = ws.ExitStatus()
	}

	return
}
