package code

import (
	"encoding/json"
	"errors"
	"github.com/google/go-github/github"
	"github.com/labstack/gommon/log"
)

type EvalContext struct {
	Generator *Input `json:"generator"`
	Solution  *Input `json:"solution"`
	Test      *Input `json:"test"`
}

func updateInput(client *github.Client, gist *github.Gist, inputs ...*Input) {
	for _, input := range inputs {
		log.Info("Input: %+v\n", input)
		if input == nil {
			continue
		}
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

var ErrNotFound = errors.New("File not found")

func fetchFile(client *github.Client, id, sha, filename string) (*github.Gist, *File, error) {
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
		return gist, &File{
			Name:    filename,
			Content: *content,
		}, nil
	}
	return gist, nil, ErrNotFound
}

func GistFetch(id string) *EvalContext {
	client := github.NewClient(nil)
	gist, file, err := fetchFile(client, id, "", "eval.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("%+v", file)

	v := &EvalContext{}
	err = json.Unmarshal([]byte(file.Content), v)
	log.Info("%+v\n", v)
	log.Info("%+v\n %+v\n %+v\n", v.Generator, v.Solution, v.Test)
	updateInput(client, gist, v.Generator, v.Solution, v.Test)
	return &EvalContext{v.Generator, v.Solution, v.Test}
}

func GistEvaluate(id string, runner *Runner) *Result {
	evalContext := GistFetch(id)
	result := Evaluate(evalContext.Generator, evalContext.Solution, evalContext.Test, runner)
	log.Info("Result: %v", result.Correct)
	return result
}
