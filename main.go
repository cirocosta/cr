package main

import (
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

	_, err := lib.ConfigFromFile(args.File)
	must(err)

	var g dag.AcyclicGraph

	g.Add(1)
	g.Add(2)
	g.Connect(dag.BasicEdge(1, 2))

	w := &dag.Walker{
		Callback: WalkFunc,
	}

	w.Update(&g)

	err = w.Wait()
	must(err)
}
