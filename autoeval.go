package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/go-lua"
)

var (
	process *ProcessBuffer
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

	// commandln("launch", append([]string{command}, args...)...)
	items := append([]string{command}, args...)
	// log.Println(Instruction("launch", items...))
	Log.Instruction("launch", items...)

	if process != nil {
		process.Stop()
	}

	var err error
	process, err = NewProcessBuffer(command, args)

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
	// commandln("stop")
	// fmt.Println(Instruction("stop"))
	Log.Instruction("stop")

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

	// parse the string to integer and multiply by time.Millisecond
	timeoutDuration := time.Duration(timeout) * time.Millisecond

	timeoutChan := time.After(timeoutDuration)
	tick := time.Tick(10 * time.Millisecond)

	// commandln("expect", expected)
	// fmt.Print(Instruction("expect", expected))
	Log.Instruction("expect", expected)

	for {
		output := string(process.Peek())

		index := strings.Index(output, expected)
		if index != -1 {
			endIndex := index + len(expected)

			process.Seek(endIndex)

			// fmt.Println(green(" OK"))
			Log.OK()
			return 0
		}

		select {
		case <-timeoutChan:
			if points == 0 {
				// fmt.Println(yellow("ERROR"))
				Log.Error(false)
			} else {
				// fmt.Println(red("ERROR"))
				Log.Error(true)
				process.Stop()
				process = nil
			}
			return 0 // Timeout
		case <-tick:
			// fmt.Print(gray("."))
			// TODO: Chooso if print or not
			// Log.Loading()
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

	// commandln("write", input)
	// fmt.Println(Instruction("write", input))
	Log.Instruction("write", input)

	err := process.Write(input)
	if err != nil {
		log.Fatalf("Error writing to the process: %v", err)
	}

	return 0
}

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		// Refactor from log.Fatal to fmt.Println and os.Exit
		// log.Fatal(red("Usage: autoeval <script.lua>"))
		fmt.Print("Usage: ")
		fmt.Println("autoeval <script.lua>")
		os.Exit(1)
	}

	silentFlag := flag.Bool("silent", false, "Suppress logs")
	flag.Parse()
	if *silentFlag {
		// TODO: Set up logger
		// log.SetOutput(io.Discard) // Suprime los logs
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
			log.Fatalf("Command not supported. Try updating to the latest version.")
		}

		log.Fatalf("Error executing the Lua script: %v", err)
	}
}
