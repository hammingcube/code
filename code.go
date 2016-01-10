package code

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"os/exec"
	"path/filepath"
	"strings"
)

type File struct {
	Name    string
	Content string
}

type Input struct {
	Language string
	Files    []File
}

type Output struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"err"`
}

type Runner struct {
	RunnerBinary string
}

type langConfig struct {
	DockerImage string
	IsSupported bool
}

var languages = map[string]*langConfig{
	"assembly":     &langConfig{"", false},
	"bash":         &langConfig{"", false},
	"c":            &langConfig{"", false},
	"clojure":      &langConfig{"", false},
	"coffeescript": &langConfig{"", false},
	"csharp":       &langConfig{"", false},
	"d":            &langConfig{"", false},
	"elixir":       &langConfig{"", false},
	"cpp":          &langConfig{"glot/clang", true},
	"erlang":       &langConfig{"", false},
	"fsharp":       &langConfig{"", false},
	"haskell":      &langConfig{"", false},
	"idris":        &langConfig{"", false},
	"go":           &langConfig{"glot/golang", true},
	"java":         &langConfig{"glot/java", false},
	"javascript":   &langConfig{"glot/javascript", true},
	"julia":        &langConfig{"", false},
	"lua":          &langConfig{"", false},
	"nim":          &langConfig{"", false},
	"ocaml":        &langConfig{"", false},
	"perl":         &langConfig{"", false},
	"php":          &langConfig{"", false},
	"python":       &langConfig{"glot/python", true},
	"ruby":         &langConfig{"", false},
	"rust":         &langConfig{"", false},
	"scala":        &langConfig{"", false},
	"swift":        &langConfig{"", false},
}

func IsNotSupported(lang string) bool {
	return languages[lang] == nil || !languages[lang].IsSupported || languages[lang].DockerImage == ""
}

func NewRunner(pathToRunner string) *Runner {
	dir, _ := filepath.Abs(pathToRunner)
	return &Runner{RunnerBinary: dir}
}

func (r *Runner) Run(input *Input) (*Output, error) {
	if IsNotSupported(input.Language) {
		return &Output{Error: "Language not supported"}, nil
	}
	dockerImg := languages[input.Language].DockerImage
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	workDir, _ := filepath.Abs(".")
	runnerDir := filepath.Dir(r.RunnerBinary)
	runnerBinary := filepath.Join("/runner", filepath.Base(r.RunnerBinary))
	args := []string{
		"docker", "run", "--rm", "-i",
		"-v", fmt.Sprintf("%s:/app", workDir), // Mounted Work Directory
		"-v", fmt.Sprintf("%s:/runner", runnerDir), // Mounted Runner Directory
		"-w", "/app",
		dockerImg,
		runnerBinary}
	log.Info("Cmd: %s", strings.Join(args, " "))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader([]byte(inputBytes))
	log.Debug("Stdin: %s", inputBytes)
	err = cmd.Run()
	if stderr.String() != "" || err != nil {
		return nil, errors.New(fmt.Sprintf(`{"_cmd_stderr": %q, "_cmd_err": "%v"}`, stderr.String(), err))
	}
	output := &Output{}
	err = json.Unmarshal(stdout.Bytes(), &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
