package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/go-lua"
)

const N = "\x1b[0m"
const R = "\x1b[31m"
const G = "\x1b[32m"
const Y = "\x1b[33m"
const B = "\x1b[34m"
const M = "\x1b[35m"
const C = "\x1b[36m"
const W = "\x1b[37m"

var (
	process *TestScenario
)

func launch(l *lua.State) int {
	argn := l.Top()
	if argn < 1 {
		log.Fatalf("Command to execute must be provided.")
	}

	command := lua.CheckString(l, 1)

	args := []string{}
	for i := 2; i <= argn; i++ {
		arg := lua.CheckString(l, i)
		args = append(args, arg)
	}

	fmt.Printf(B+"\nlaunch("+Y+"%s"+B+")\n", strings.Join(append([]string{command}, args...), B+", "+Y))

	if process != nil {
		process.Stop()
	}

	var err error
	process, err = NewTestScenario(command, args)

	if err != nil {
		log.Fatalf("Error preparing the command: %v", err)
	}

	err = process.Launch()
	if err != nil {
		log.Fatalf("Error starting the command: %v", err)
	}

	return 0
}

func stop(_ *lua.State) int {
	fmt.Printf(B + "stop()\n")

	if process != nil {
		process.Stop()
	}
	process = nil

	return 0
}

func expect(l *lua.State) int {
	// If process has being killed, do not perform any action
	if process == nil {
		return 0
	}

	expected := lua.CheckString(l, 1)
	points := lua.OptInteger(l, 2, 0)
	timeout := lua.OptInteger(l, 3, 1000)

	timeoutDuration := time.Duration(timeout) * time.Millisecond

	timeoutChan := time.After(timeoutDuration)
	tick := time.Tick(200 * time.Millisecond)

	fmt.Printf(B+"expect("+Y+"%s"+B+") ", expected)

	for {
		output := string(process.Peek())

		index := strings.Index(output, expected)
		if index != -1 {
			endIndex := index + len(expected)

			process.Seek(endIndex)

			fmt.Printf(G + "OK\n")
			return 0
		}

		select {
		case <-timeoutChan:
			if points == 0 {
				fmt.Printf(R + "ERROR\n")
			} else {
				fmt.Printf(R + "CRITICAL ERROR\n")
				process.Stop()
				process = nil
			}
			return 0 // Timeout
		case <-tick:
			fmt.Printf(". ")
			// Continue retrying
		}
	}
}

func write(l *lua.State) int {
	// If process has being killed, do not perform any action
	if process == nil {
		return 0
	}

	input := lua.CheckString(l, 1)

	fmt.Printf(B+"write("+Y+"%s"+B+")\n", input)

	err := process.Write(input)
	if err != nil {
		log.Fatalf("Error writing to the process: %v", err)
	}

	return 0
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("The path to the Lua script must be provided as an argument.")
	}

	scriptPath := os.Args[1]
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Could not read the Lua file: %v", err)
	}

	l := lua.NewState()
	lua.OpenLibraries(l)

	l.Register("launch", launch)
	l.Register("expect", expect)
	l.Register("write", write)
	l.Register("stop", stop)

	if err := lua.DoString(l, string(script)); err != nil {
		if strings.Contains(err.Error(), "attempt to call a nil value") {
			log.Fatalf(R + "\nCommand not supported. Try updating to the latest version." + N)
		}

		log.Fatalf(R+"\nError executing the Lua script: %v", err)
	}
}
