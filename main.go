package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/cr/lib"
	"github.com/rs/zerolog"
)

var (
	args = &lib.Runtime{
		File:   "./.cr.yml",
		Stdout: true,
		Graph:  false,
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
	log.SetOutput(ioutil.Discard)

	cfg, err := lib.ConfigFromFile(args.File)
	must(err)

	cfg.OnJobStatusChange = func(a *lib.Activity) {
		ui.WriteActivity(a)
	}

	executor, err := lib.New(&cfg)
	must(err)

	if args.Graph {
		fmt.Println(executor.GetDotGraph())
		os.Exit(0)
	}

	err = executor.Execute(context.Background())
	must(err)
}
