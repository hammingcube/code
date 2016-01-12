package main

import (
	"flag"
	"github.com/labstack/gommon/log"
	"github.com/maddyonline/code"
	"os"
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

var STDIN_EXAMPLE = code.File{
	Name:    "_stdin_",
	Content: "abc\nhello",
}

const GEN_CODE = `
# include <iostream>
# include <vector>
using namespace std;
int main() {
	vector<string> vec({"ab", "hello", "really"});
	for(auto & v: vec) {
		cout << v << endl;
	}

}
`

const CODE1 = `
# include <iostream>
# include <vector>
using namespace std;
int main() {
	string s;
	while(cin >> s) {
		cout << s.size() << endl;
	}
}
`

const CODE2 = `
# include <iostream>
# include <vector>
using namespace std;
int main() {
	string s;
	while(cin >> s) {
		const char * str = s.c_str();
		int count = 0;
		for(const char *ptr = str; *ptr != '\0'; ptr++, count++);
		cout << count << endl;
	}
}
`

var CODES = []string{CODE1, CODE2}

func main() {
	var pathToRunner string
	flag.StringVar(&pathToRunner, "runner", os.Getenv("RUNNER_BINARY"), "path to runner binary")
	flag.Parse()
	runner := code.NewRunner(pathToRunner)
	mainRun(runner)
}

func runIt(runner *code.Runner) {
	input := code.MakeInput("cpp", "main.cpp", CODE2, STDIN_EXAMPLE)
	output, err := runner.Run(input)
	if err != nil {
		log.Info("%v", err)
	} else {
		log.Info("%#v", output)
	}
}

func mainRun(runner *code.Runner) {
	inputGen := code.MakeInput("cpp", "main.cpp", GEN_CODE, code.StdinFile(""))
	inputCode1 := code.MakeInput("cpp", "main.cpp", CODES[0], code.StdinFile(""))
	inputCode2 := code.MakeInput("cpp", "main.cpp", CODES[1], code.StdinFile(""))
	code.Evaluate(inputGen, inputCode1, inputCode2, runner)
}
