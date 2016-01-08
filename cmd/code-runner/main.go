package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/maddyonline/code-runner"
)

const stdin = `{
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
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", ".", "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)
	input := &code.Input{}
	err := json.Unmarshal([]byte(stdin), input)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	output, err := runner.Run(input)
	if err != nil {
		fmt.Printf("%v", err)
	} else {
		fmt.Printf("%#v", output)
	}
}
