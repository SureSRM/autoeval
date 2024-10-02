//go:build windows

package main

import (
	"io"
    "os"
	"os/exec"
	"path/filepath"
	"strings"
)

type WinCLIProcess struct {
	command string
	args    []string

	cmd *exec.Cmd

	stdoutPipe io.ReadCloser
	stdinPipe  io.WriteCloser
}

func newCLIProcess(command string, args []string) (CLIProcess, error) {
	// filepath.Abs() calls filepath.Clean() which translates the separators to the OS's default separator
	if strings.Contains(command, "/") {
		var err error
		command, err = filepath.Abs(command)

        if _, err := os.Stat(command); os.IsNotExist(err) {
            // If file does not exist, we add the .exe extension
            // This is just for the Windows impl
            command += ".exe"
        }
		if err != nil {
			return nil, err
		}
	}

	cmd := exec.Command(command, args...)

	process := &WinCLIProcess{
		command: command,
		args:    args,
		cmd:     cmd,
	}

	return process, nil
}

func (p *WinCLIProcess) Start() error {
    var err error

	p.stdoutPipe, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	p.stdinPipe, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = p.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (p *WinCLIProcess) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := p.stdoutPipe.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (p *WinCLIProcess) Write(input string) error {
	_, err := io.Copy(p.stdinPipe, strings.NewReader(input+"\r\n"))
	return err
}

func (p *WinCLIProcess) Stop() error {
	return p.cmd.Wait()
}
