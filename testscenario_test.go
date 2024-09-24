package main

import (
	"fmt"
	"testing"
	"time"
)

func expectOutput(t *testing.T, process *TestScenario, expected string) {
	time.Sleep(100 * time.Millisecond)
	bytes := process.Peek()
	process.Seek(len(bytes))
	if string(bytes) != expected {
		t.Errorf("Expected length %d, got %d", len(expected), len(bytes))
		t.Errorf("Expected '%s', got '%s'", expected, string(bytes))
	}
}

func TestNewTestScenario(t *testing.T) {
	process, err := NewTestScenario("./assets/tool.sh", []string{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	process.Launch()

	expectOutput(t, process, "Name: ")

	process.Write("AAA")

	expectOutput(t, process, "AAA\r\nHello, AAA\r\nAge: ")

	process.Write("32")

	expectOutput(t, process, "32\r\nYou are 32 years old\r\nYour arg was: \r\n")

	process.Stop()
}
