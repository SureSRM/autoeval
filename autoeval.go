package main

import (
	"fmt"
	"log"
    "flag"
	"os"
	"strings"
	"time"

	"github.com/Shopify/go-lua"
	"github.com/charmbracelet/lipgloss"
)

func color(color string) func(...string) string {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
}

var (
	red = color("1")
	green = color("2")
    yellow = color("3")
    blue = color("4")
    magenta = color("5")
    cyan = color("6")
    gray = color("7")
    black = color("8")
    light_red = color("9")
    light_green = color("10")
    light_yellow = color("11")
    light_blue = color("12")
    light_magenta = color("13")
    light_cyan = color("14")
    white = color("15")
)

func instruction(kind string, args ...string) string {

	// Combine the 'kind' and 'args' formatted
	result := blue(kind) + blue("(")

	for i, arg := range args {
		result += yellow(arg)
		if i < len(args)-1 {
			result += blue(",")
		}
	}

	result += blue(")")

    return result
}

var (
	process *TestScenario
)

func launch(l *lua.State) int {
	argn := l.Top()
	if argn < 1 {
		log.Fatalf(red("Command to execute must be provided."))
	}

	command := lua.CheckString(l, 1)

	args := []string{}
	for i := 2; i <= argn; i++ {
		arg := lua.CheckString(l, i)
		args = append(args, arg)
	}

	// commandln("launch", append([]string{command}, args...)...)
    items := append([]string{command}, args...)
    log.Println(instruction("launch", items...))

	if process != nil {
		process.Stop()
	}

	var err error
	process, err = NewTestScenario(command, args)

	if err != nil {
		log.Fatalf(red("Error preparing the command: %v"), err)
	}

	err = process.Launch()
	if err != nil {
		log.Fatalf(red("Error starting the command: %v"), err)
	}

	return 0
}

func stop(_ *lua.State) int {
	// commandln("stop")
    fmt.Println(instruction("stop"))

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
	tick := time.Tick(200 * time.Millisecond)

    // commandln("expect", expected)
    fmt.Print(instruction("expect", expected))

	for {
		output := string(process.Peek())

		index := strings.Index(output, expected)
		if index != -1 {
			endIndex := index + len(expected)

			process.Seek(endIndex)

			fmt.Println(green(" OK"))
			return 0
		}

		select {
		case <-timeoutChan:
			if points == 0 {
				fmt.Println(yellow("ERROR"))
			} else {
				fmt.Println(red("ERROR"))
				process.Stop()
				process = nil
			}
			return 0 // Timeout
		case <-tick:
			fmt.Print(gray("."))
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
    fmt.Println(instruction("write", input))

	err := process.Write(input)
	if err != nil {
		log.Fatalf(red("Error writing to the process: %v"), err)
	}

	return 0
}

func main() {
    log.SetFlags(0)

	if len(os.Args) < 2 {
        // Refactor from log.Fatal to fmt.Println and os.Exit
        // log.Fatal(red("Usage: autoeval <script.lua>"))
        fmt.Print(red("Usage: "))
        fmt.Println(yellow("autoeval <script.lua>"))
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
		log.Fatalf(red("Could not read the Lua file: %v"), err)
	}

	l := lua.NewState()
	lua.OpenLibraries(l)

	l.Register("launch", launch)
	l.Register("expect", expect)
	l.Register("write", write)
	l.Register("stop", stop)

	if err := lua.DoString(l, string(script)); err != nil {
		if strings.Contains(err.Error(), "attempt to call a nil value") {
            log.Fatalf(red("Command not supported. Try updating to the latest version."))
		}

        log.Fatalf(red("Error executing the Lua script: %v"), err)
	}
}
