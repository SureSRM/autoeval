package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type colorFunction func(...string) string

func color(color string) colorFunction {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
}

// var (
// 	red           = color("1")
// 	green         = color("2")
// 	yellow        = color("3")
// 	blue          = color("4")
// 	magenta       = color("5")
// 	cyan          = color("6")
// 	gray          = color("7")
// 	black         = color("8")
// 	light_red     = color("9")
// 	light_green   = color("10")
// 	light_yellow  = color("11")
// 	light_blue    = color("12")
// 	light_magenta = color("13")
// 	light_cyan    = color("14")
// 	white         = color("15")
// )

var (
	base    = color("#5d98d3")
	base2   = color("#5eacf9")
	dim     = color("#d7ffae")
	accent  = color("#f7e165")
	accent2 = color("#a95ef9")
	success = color("#1bfcc4")
	err     = color("#e53939")
	err2    = color("#f49c3d")
)

type Logger struct {
	silent bool
}

func (l *Logger) setSilent(silent bool) {
	l.silent = silent
}

func (l *Logger) Instruction(kind string, args ...string) {
	var pickColor colorFunction
	switch kind {
	case "launch":
		pickColor = base2
	case "stop":
		pickColor = base2
	default:
		pickColor = base
	}

	// Combine the 'kind' and 'args' formatted
	result := pickColor(kind) + pickColor("(")

	for i, arg := range args {
		result += accent(arg)
		if i < len(args)-1 {
			result += pickColor(",")
		}
	}

	result += pickColor(")")

	switch kind {
	case "expect":
		fmt.Print(result)
	default:
		fmt.Println(result)
	}
}

func (l *Logger) Loading() {
	fmt.Print(dim("."))
}

func (l *Logger) OK() {
	fmt.Println(success(" OK"))
}

func (l *Logger) Error(severe bool) {
	if severe {
		fmt.Println(err(" ERROR"))
	} else {
		fmt.Println(err2(" ERROR"))
	}
}

var Log = Logger{silent: false}
