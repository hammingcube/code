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
	Id      string `json:"id"`
	Sha     string `json:"sha"`
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
	log.Info("Starting run...")
	if IsNotSupported(input.Language) {
		return &Output{}, errors.New(fmt.Sprintf("Language %s not supported", input.Language))
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

func StdinFile(content string) File {
	return File{
		Name:    "_stdin_",
		Content: content,
	}
}

func UpdateStdin(input *Input, stdinFile File) {
	for _, file := range input.Files {
		if file.Name == "_stdin_" {
			file.Content = stdinFile.Content
			return
		}
	}
	input.Files = append(input.Files, stdinFile)
	return
}

func MakeInput(language, name, content string, input File) *Input {
	return &Input{
		Language: language,
		Files: []File{
			File{
				Name:    name,
				Content: content,
			},
			input,
		},
	}
}

func Evaluate(inputGen, inputCode1, inputCode2 *Input, runner *Runner) {
	gen := []chan *Output{
		make(chan *Output),
		make(chan *Output),
	}
	results := make(chan struct {
		InputStr *string
		Output   *Output
	}, 2)
	go func() {
		output, err := runner.Run(inputGen)
		process(output, err)
		go func() { gen[0] <- output }()
		go func() { gen[1] <- output }()
	}()
	inputs := []*Input{inputCode1, inputCode2}
	for i := 0; i < 2; i++ {
		go func(i int) {
			log.Info("Using index: %d\n", i)
			genOutput := <-gen[i]
			input := inputs[i]
			UpdateStdin(input, StdinFile(genOutput.Stdout))
			output, err := runner.Run(input)
			process(output, err)
			results <- struct {
				InputStr *string
				Output   *Output
			}{&genOutput.Stdout, output}
		}(i)
	}
	out1, out2 := <-results, <-results

	if diff(out1.Output.Stdout, out2.Output.Stdout) {
		log.Info("Different on input %q: %q %q", *out1.InputStr, out1.Output.Stdout, out2.Output.Stdout)
	} else {
		log.Info("Identical on input %q", *out1.InputStr)
	}
}

func diff(a, b string) bool {
	return a != b
}

func process(output *Output, err error) {
	if err != nil {
		log.Info("%v", err)
	} else {
		log.Info("%#v", output)
	}
}
