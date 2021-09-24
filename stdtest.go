package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// TestConfig stores information relevant for a single test.
type TestConfig struct {
	// Name is the name of a test.
	Name string `json:"name"`
	// Input is a string representing standard input.
	Input string `json:"input"`
	// Output is a set of independent strings that are expected to be found in
	// the output.
	Output []string `json:"output"`
}

// openTestConfig opens a `.json` file, deserializes it and returns
// []TestConfig.
//
// The `.json` file must be compatible with []TestConfig type.
func openTestConfig(path string) ([]TestConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)
	var tests []TestConfig
	err = json.Unmarshal(bytes, &tests)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

// matchOutput checks if all strings in `testStrs` are found in output (in any
// order). If not, the error is returned.
func matchOutput(output string, testStrs []string) error {
	for _, str := range testStrs {
		if !strings.Contains(output, str) {
			return fmt.Errorf("%q not found in output %q", str, output)
		}
	}
	return nil
}

// testProgram runs a series of tests represented by `tests` for a Go program
// with a file found by `path`.
func testProgram(path string, tests []TestConfig) error {
	for _, test := range tests {
		cmd := exec.Command("go", "run", path)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, test.Input)
		}()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if err := matchOutput(string(out), test.Output); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	execFlag := flag.String("go_source_path", "", "Path to a Go program source code to run.")
	inputFlag := flag.String("tests_config", "", "Config path for tests.")
	flag.Parse()
	if *execFlag == "" || *inputFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tests, err := openTestConfig(*inputFlag)
	if err != nil {
		log.Fatal(err)
	}
	if err := testProgram(*execFlag, tests); err != nil {
		log.Fatal(err)
	}
}
