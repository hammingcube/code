package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

const stdinBytes = `{
  "language": "cpp",
  "files": [
    {
      "name": "main.cpp",
      "content": "# include <iostream>\nusing namespace std;\nint main() {string s;while(cin >> s) {cout << s.size() << endl;}}"
    }, {
    "name": "_stdin_",
    "content": "abc\nhello"
    }
  ]
}`

func main() {
	workDir, _ := filepath.Abs(".")
	args := []string{"docker", "run", "--rm", "-i", "-v", fmt.Sprintf("%s:/app", workDir), "-w", "/app", "rsmmr/clang", "./runner"}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader([]byte(stdinBytes))
	err := cmd.Run()
	fmt.Printf("stdout:\n%s\nstdin:\n%s\nerr:\n%v\n", stdout.String(), stderr.String(), err)
}
