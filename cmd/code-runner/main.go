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

func process(output *code.Output, err error) {
	if err != nil {
		log.Info("%v", err)
	} else {
		log.Info("%#v", output)
	}
}

func mainRun(runner *code.Runner) {
	gen := []chan *code.Output{
		make(chan *code.Output),
		make(chan *code.Output),
	}
	results := make(chan struct {
		InputStr *string
		Output   *code.Output
	}, 2)
	go func() {
		input := code.MakeInput("cpp", "main.cpp", GEN_CODE, code.StdinFile(""))
		output, err := runner.Run(input)
		process(output, err)
		go func() { gen[0] <- output }()
		go func() { gen[1] <- output }()
	}()
	for i := 0; i < 2; i++ {
		go func(i int) {
			log.Info("Using index: %d\n", i)
			genOutput := <-gen[i]
			input := code.MakeInput("cpp", "main.cpp", CODES[i], code.StdinFile(genOutput.Stdout))
			output, err := runner.Run(input)
			process(output, err)
			results <- struct {
				InputStr *string
				Output   *code.Output
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
