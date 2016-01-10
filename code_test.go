package code

import (
	"encoding/json"
	"os"
	"testing"
)

var STDIN_EXAMPLE = File{
	Name:    "_stdin_",
	Content: "abc\nhello",
}

var EXPECTED_OUTPUT = &Output{Stdout: "3\n5\n", Stderr: "", Error: ""}

const PY_CODE = `
import sys

for line in sys.stdin:
  line = line.strip().rstrip()
  print(len(line))
`

const CPP_CODE = `
# include <iostream>
using namespace std;
int main() {
  string s;
  while(cin >> s) {
    cout << s.size() << endl;
  }
}
`

const GO_CODE = `
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
      fmt.Println(len(scanner.Text()))
  }
}
`

const JAVASCRIPT_CODE = `
var readline = require('readline');
var rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
});

rl.on('line', function(line){
    console.log(line.length);
})
`

const JAVA_CODE = `
import java.io.BufferedInputStream;
import java.util.Scanner;

public class MainClass {
    public static void main(String args[]) {
        Scanner stdin = new Scanner(new BufferedInputStream(System.in));
        while (stdin.hasNext()) {
          String word = stdin.next();
            System.out.println(word.length());
        }
    }
}
`

var TEST_CASES = []testpair{
	{MakeInput("python", "main.py", PY_CODE, STDIN_EXAMPLE), EXPECTED_OUTPUT},
	{MakeInput("cpp", "main.cpp", CPP_CODE, STDIN_EXAMPLE), EXPECTED_OUTPUT},
	{MakeInput("go", "main.go", GO_CODE, STDIN_EXAMPLE), EXPECTED_OUTPUT},
	{MakeInput("javascript", "main.js", JAVASCRIPT_CODE, STDIN_EXAMPLE), EXPECTED_OUTPUT},
	//{makeInput("java", "MainClass.java", JAVA_CODE), EXPECTED_OUTPUT},
}

const basic_example = `{
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

type testpair struct {
	input    *Input
	expected *Output
}

func addBasic(t *testing.T, testpairs []testpair) []testpair {
	input := &Input{}
	err := json.Unmarshal([]byte(basic_example), input)
	if err != nil {
		t.Error("got err: ", err)
	}

	return append(testpairs, testpair{input, EXPECTED_OUTPUT})
}

func makeTestPairs(t *testing.T) []testpair {
	testpairs := []testpair{}
	addBasic(t, testpairs)
	return testpairs
}

func TestRunner(t *testing.T) {
	runner := NewRunner(os.Getenv("RUNNER_BINARY"))
	var tests = TEST_CASES
	for _, pair := range tests {
		output, err := runner.Run(pair.input)
		if err != nil {
			t.Error(
				"For", pair.input,
				"expected err:", nil,
				"got", err,
			)
		}
		if *output != *pair.expected {
			t.Errorf("For %s, expected: %q, got: %q", pair.input.Language, pair.expected, output)
		}
	}
}
