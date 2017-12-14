package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/cr/lib"
	"github.com/rs/zerolog"
)

var version string = "dev"

type cliArgs struct {
	lib.Runtime
}

func (r cliArgs) Version() string {
	return "cr - the concurrent runner - version=" + version
}

var (
	args = &cliArgs{
		lib.Runtime{
			File:          "./.cr.yml",
			LogsDirectory: "/tmp",
			Stdout:        false,
			Graph:         false,
		},
	}
	logger = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
	ui = lib.NewUi()
)

func must(err error) {
	if err == nil {
		return
	}

	logger.Fatal().
		Err(err).
		Msg("main execution failed")
}

func main() {
	arg.MustParse(args)

	rand.Seed(time.Now().UnixNano())
	log.SetOutput(ioutil.Discard)

	cfg, err := lib.ConfigFromFile(args.File)
	must(err)

	cfg.OnJobStatusChange = func(a *lib.Activity) {
		ui.WriteActivity(a)
	}

	cfg.Runtime = args.Runtime

	executor, err := lib.New(&cfg)
	must(err)

	if args.Graph {
		fmt.Println(executor.GetDotGraph())
		os.Exit(0)
	}

	fmt.Printf(`
	Starting execution.

	Logs directory:	%s
	`+"\n", cfg.Runtime.LogsDirectory)

	err = executor.Execute(context.Background())
	must(err)
}
