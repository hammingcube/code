package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/google/go-github/github"
	"github.com/labstack/gommon/log"
	"github.com/maddyonline/code"
	"os"
)

type Eval struct {
	Generator *code.Input `json:"generator"`
	Solution  *code.Input `json:"solution"`
	Test      *code.Input `json:"test"`
}

func updateInput(client *github.Client, gist *github.Gist, inputs ...*code.Input) {
	for _, input := range inputs {
		log.Info("Input: %+v\n", input)
		for i, file := range input.Files {
			log.Info("%+v", file)
			if file.Content == "" {
				if gistId, content := file.Id, gist.Files[github.GistFilename(file.Name)].Content; gistId == "" && content != nil {
					input.Files[i].Content = *content
					continue
				}
				if _, fetchedFile, err := fetchFile(client, file.Id, file.Sha, file.Name); err == nil {
					log.Info("Fetching %s from sha %s", file.Name, file.Id)
					input.Files[i].Content = fetchedFile.Content
				}
			}
		}
		//log.Info("Input: %+v\n", input)
	}
}

func main() {
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", os.Getenv("RUNNER_BINARY"), "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)
	id := "4f1bae999b5fbea43624"
	mainRun(id, runner)
}

var ErrNotFound = errors.New("File not found")

func fetchFile(client *github.Client, id, sha, filename string) (*github.Gist, *code.File, error) {
	var gist *github.Gist
	var err error
	if sha != "" {
		gist, _, err = client.Gists.GetRevision(id, sha)
	} else {
		gist, _, err = client.Gists.Get(id)
	}
	if err != nil {
		return gist, nil, err
	}
	if content := gist.Files[github.GistFilename(filename)].Content; content != nil {
		return gist, &code.File{
			Name:    filename,
			Content: *content,
		}, nil
	}
	return gist, nil, ErrNotFound
}

func mainRun(id string, runner *code.Runner) {
	client := github.NewClient(nil)
	gist, file, err := fetchFile(client, id, "", "eval.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("%+v", file)

	v := &Eval{}
	err = json.Unmarshal([]byte(file.Content), v)
	log.Info("%+v\n", v)
	log.Info("%+v\n %+v\n %+v\n", v.Generator, v.Solution, v.Test)
	updateInput(client, gist, v.Generator, v.Solution, v.Test)
	code.Evaluate(v.Generator, v.Solution, v.Test, runner)
}
