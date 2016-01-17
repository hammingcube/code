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
// Copyright (c) 2015 Elements of Programming Interviews. All rights reserved.

#include <cassert>
#include <iostream>
#include <random>
#include <string>

using std::boolalpha;
using std::cout;
using std::default_random_engine;
using std::endl;
using std::random_device;
using std::string;
using std::uniform_int_distribution;

string RandString(int len) {
  string ret;
  default_random_engine gen((random_device())());
  while (len--) {
    uniform_int_distribution<int> dis('a', 'z');
    ret += dis(gen);
  }
  return ret;
}

int main(int argc, char *argv[]) {
  default_random_engine gen((random_device())());
  for (int times = 0; times < 100; ++times) {
    string s;
    if (argc == 2) {
      s = argv[1];
    } else {
      uniform_int_distribution<int> dis(1, 10);
      s = RandString(dis(gen));
    }
    cout << s << endl;
  }
  return 0;
}
`

const CODE1 = `
#include <algorithm>
#include <string>
#include <iostream>

using std::string;

// Please update the following function
bool CanStringBeAPalindrome(string* s) {
  return false;
}

int main() {
    string s;
    while(std::cin >> s) {
        std::cout << CanStringBeAPalindrome(&s) << std::endl;
    }
    return 0;
}
`

const CODE2 = `
#include <algorithm>
#include <string>
#include <iostream>

using std::string;

bool CanStringBeAPalindrome(string* s) {
  sort(s->begin(), s->end());
  int odd_count = 0, num_curr_char = 1;

  for (int i = 1; i < s->size() && odd_count <= 1; ++i) {
    if ((*s)[i] != (*s)[i - 1]) {
      if (num_curr_char % 2) {
        ++odd_count;
      }
      num_curr_char = 1;
    } else {
      ++num_curr_char;
    }
  }
  if (num_curr_char % 2) {
    ++odd_count;
  }

  // A string can be permuted as a palindrome if the number of odd time
  // chars <= 1.
  return odd_count <= 1;
}

int main() {
    string s;
    while(std::cin >> s) {
        std::cout << CanStringBeAPalindrome(&s) << std::endl;
    }
    return 0;
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
