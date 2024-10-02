package main

import (
	"fmt"
	"testing"
	"time"
    "strings"
)

func expectOutput(t *testing.T, process *TestScenario, expected string) {
	time.Sleep(200 * time.Millisecond)

	output := string(process.Peek())

	index := strings.Index(output, expected)
	if index == -1 {
        t.Errorf("Expected:\n%s\n", expected)
        t.Errorf("Got:\n%s\n", output)
        return
	} else {
        endIndex := index + len(expected)
        process.Seek(endIndex)
	}
}

func TestNewTestScenario(t *testing.T) {
	process, err := NewTestScenario("./examples/tool.sh", []string{"arg1"})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	process.Launch()

	expectOutput(t, process, "Name: ")
	process.Write("AAA")
	expectOutput(t, process, "Hello, AAA")

	expectOutput(t, process, "Age: ")
	process.Write("32")
	expectOutput(t, process, "You are 32 years old")

	expectOutput(t, process, "Your arg was: arg1")

	process.Stop()
}
