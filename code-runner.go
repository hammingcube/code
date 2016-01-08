package code

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	Stdout string
	Stderr string
	Error  string
}

type Runner struct {
	RunnerBinary string
	Verbose      bool
}

func NewRunner(pathToRunner string) *Runner {
	dir, _ := filepath.Abs(pathToRunner)
	return &Runner{RunnerBinary: dir}
}

func (r *Runner) Run(input *Input) (*Output, error) {
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
		"-w", "/app", "rsmmr/clang",
		runnerBinary}
	if r.Verbose {
		fmt.Println("Running...")
		fmt.Printf("%s", strings.Join(args, " "))
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader([]byte(inputBytes))
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
