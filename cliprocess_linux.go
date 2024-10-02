//go:build linux

package main

import (
	"github.com/iyzyi/aiopty/pty"
	// "io"
	"path/filepath"
	"strings"
)

type LinuxCLIProcess struct {
	cmd  string
	args []string
	pty  *pty.Pty
}

func newCLIProcess(command string, args []string) (CLIProcess, error) {
	// filepath.Abs() calls filepath.Clean() which translates the separators to the OS's default separator
	if strings.Contains(command, "/") {
		var err error
		command, err = filepath.Abs(command)

		if err != nil {
			return nil, err
		}
	}

	process := &LinuxCLIProcess{
		cmd:  command,
		args: args,
		pty:  nil,
	}

	return process, nil
}

func (p *LinuxCLIProcess) Start() error {
	pty, err := pty.OpenWithOptions(&pty.Options{
		Path: p.cmd,
		Args: append([]string{p.cmd}, p.args...),
	})

	if err != nil {
		return err
	}
	p.pty = pty

	return nil
}

func (p *LinuxCLIProcess) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := p.pty.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (p *LinuxCLIProcess) Write(input string) error {
    _, err := p.pty.Write([]byte(input + "\n"))
	return err
}

func (p *LinuxCLIProcess) Stop() error {
	return p.pty.Close()
}
