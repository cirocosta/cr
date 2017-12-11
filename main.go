package main

import (
	"context"
	"fmt"
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
		Graph:  false,
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

func main() {
	arg.MustParse(args)
	log.SetOutput(ioutil.Discard)

	cfg, err := lib.ConfigFromFile(args.File)
	must(err)

	graph, err := lib.BuildDependencyGraph(cfg.Jobs)
	must(err)

	if args.Graph {
		dot := string(graph.Dot(&dag.DotOpts{}))
		fmt.Println(dot)
		os.Exit(0)
	}

	err = lib.TraverseAndExecute(context.Background(), &graph)
	must(err)
}
