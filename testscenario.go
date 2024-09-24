package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
    "path/filepath"

	"github.com/iyzyi/aiopty/pty"
)


type CLIProcess struct {
	cmd string
    args []string
	pty *pty.Pty
}

func newCLIProcess(command string, args []string) (*CLIProcess, error) {
    // filepath.Abs() calls filepath.Clean() which translates the separators to the OS's default separator
    cmd, err := filepath.Abs(command)
    if err != nil {
        return nil, err
    }


	process := &CLIProcess{
		cmd: cmd,
		args: args,
		pty: nil,
	}

	return process, nil
}

func (p *CLIProcess) start() error {
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

func (p *CLIProcess) read() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := p.pty.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (p *CLIProcess) Write(input string) error {
	_, err := p.pty.Write([]byte(input + "\n"))
	return err
}

func (p *CLIProcess) Stop() error {
	return p.pty.Close()
}

type PeekBuffer struct {
	lock   sync.Mutex
	buffer bytes.Buffer
}

func (pb *PeekBuffer) add(data []byte) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	pb.buffer.Write(data)
}

func (pb *PeekBuffer) Peek() []byte {
	pb.lock.Lock()
	defer pb.lock.Unlock()

    bytes := pb.buffer.Bytes()
	return bytes
}

func (pb *PeekBuffer) Seek(n int) error {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.buffer.Len() < n {
		return fmt.Errorf("Error seeking: buffer length is %d, but tried to seek %d", pb.buffer.Len(), n)
	}

	remaining := pb.buffer.Bytes()[n:]
	pb.buffer.Reset()
	pb.buffer.Write(remaining)

	return nil
}

type TestScenario struct {
	PeekBuffer
	CLIProcess
}

func (process *TestScenario) captureOutput() {
	for {
		time.Sleep(10 * time.Millisecond)
		bytes, err := process.read()
		if err == nil && len(bytes) > 0 {
			process.add(bytes)
		} else {
			if err.Error() == "EOF" {
				break
			}
			// TODO: Check if it is correct toignore Read() errors
			// I want to ignore the empty response
			// My asumption is that it does not matter if the process crashed,
			// because I handle that independently
		}
	}
}

func (process *TestScenario) Launch() error {
	err := process.start()
	if err != nil {
		return err
	}
	go process.captureOutput()
	return nil
}

func NewTestScenario(command string, args []string) (*TestScenario, error) {
	process, err := newCLIProcess(command, args)
	if err != nil {
		return nil, err
	}

	testScenario := &TestScenario{
		PeekBuffer: PeekBuffer{},
		CLIProcess: *process,
	}

	return testScenario, nil
}
