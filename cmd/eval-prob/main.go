package main

import (
	"flag"
	"fmt"
	"github.com/maddyonline/code"
	"github.com/zabawaba99/firego"
	"log"
	"os"
)

func main() {
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", os.Getenv("RUNNER_BINARY"), "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)

	fb := firego.New("https://thinkhike.firebaseio.com/problems", nil)
	problems := map[string]string{}

	if err := fb.Value(&problems); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", problems)
	if runner != nil && false {
	}
	ec := code.GistFetch(problems["prob-1"])
	fmt.Printf("%#v", ec)
}
