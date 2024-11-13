package main

type CLIProcess interface {
	Start() error
	Stop() error
	Read() ([]byte, error)
	Write(string) error
}
