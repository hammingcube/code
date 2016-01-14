package main

import (
	"flag"
	"github.com/maddyonline/code"
	"os"
)

func main() {
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", os.Getenv("RUNNER_BINARY"), "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)
	id := "4f1bae999b5fbea43624"
	if len(flag.Args()) > 0 {
		id = flag.Args()[0]
	}
	code.GistEvaluate(id, runner)
}
