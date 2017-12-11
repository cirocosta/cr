package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/cr/lib"
	"github.com/hashicorp/terraform/dag"
	"github.com/rs/zerolog"
)

var (
	args = &lib.Runtime{
		File:   "./.cr.yml",
		Stdout: true,
	}
	logger = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
)

func must(err error) {
	if err == nil {
		return
	}

	logger.Fatal().
		Err(err).
		Msg("main execution failed")
}

func WalkFunc(v dag.Vertex) (err error) {
	logger.Info().Interface("vertex", v).Msg("walk")
	return
}

func main() {
	arg.MustParse(args)
	log.SetOutput(ioutil.Discard)

	cfg, err := lib.ConfigFromFile(args.File)
	must(err)

	graph, err := lib.BuildDependencyGraph(cfg.Jobs)
	must(err)

	err = lib.TraverseAndExecute(context.Background(), &graph)
	must(err)
}
