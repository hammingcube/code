package main

import (
	"encoding/json"
	"flag"
	"github.com/google/go-github/github"
	"github.com/maddyonline/code"
	"log"
	"os"
)

type Eval struct {
	Generator *code.Input `json:"generator"`
	Solution  *code.Input `json:"solution"`
	Test      *code.Input `json:"test"`
}

func updateInput(gist *github.Gist, inputs ...*code.Input) {
	for _, input := range inputs {
		log.Printf("Input: %+v\n", input)
		for i, file := range input.Files {
			if file.Content == "" {
				log.Printf("File %s is empty\n", file.Name)
				content := *gist.Files[github.GistFilename(file.Name)].Content
				//log.Printf("But, file content from gist is as follows.\n%s\n", content)
				input.Files[i].Content = content
			}
		}
		//log.Printf("Input: %+v\n", input)
	}
}

func main() {
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", os.Getenv("RUNNER_BINARY"), "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)
	mainRun(runner)
}

func mainRun(runner *code.Runner) {
	client := github.NewClient(nil)
	gist, _, err := client.Gists.Get("4f1bae999b5fbea43624")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", gist.Files["eval.json"])
	v := &Eval{}
	err = json.Unmarshal([]byte(*gist.Files["eval.json"].Content), v)
	log.Printf("%+v\n", v)
	log.Printf("%+v\n %+v\n %+v\n", v.Generator, v.Solution, v.Test)
	updateInput(gist, v.Generator, v.Solution, v.Test)
	code.Evaluate(v.Generator, v.Solution, v.Test, runner)
}
