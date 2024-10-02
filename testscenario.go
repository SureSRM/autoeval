package main

import (
	"fmt"
	"sync"
	"time"
)


type PeekBuffer struct {
	lock   sync.Mutex
	buffer []byte

	pointer int
}

func (pb *PeekBuffer) add(data []byte) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	pb.buffer = append(pb.buffer, data...)
}

func (pb *PeekBuffer) Peek() []byte {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	bytes := pb.buffer[pb.pointer:]
	return bytes
}

func (pb *PeekBuffer) Seek(n int) error {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	if len(pb.buffer) < n {
		return fmt.Errorf("Error seeking: buffer length is %d, but tried to seek %d", pb.buffer, n)
	}
	pb.pointer += n
	return nil
}

func (pb *PeekBuffer) Dump() []byte {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	return pb.buffer
}

type TestScenario struct {
	PeekBuffer
	CLIProcess
}

func (process *TestScenario) captureOutput() {
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-tick:
			bytes, err := process.Read()
			if err == nil && len(bytes) > 0 {
				process.add(bytes)
			} else {
				if err.Error() == "EOF" {
                    // TODO: Implement EOF detection for Win and Linux
					// process.add([]byte("<EOF>"))
					break
				}
				// TODO: Check if it is correct toignore Read() errors
				// I want to ignore the empty response
				// My asumption is that it does not matter if the process crashed,
				// because I handle that independently
			}
		}
	}
}

func (process *TestScenario) Launch() error {
	err := process.Start()
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
		CLIProcess: process,

	}

	return testScenario, nil
}
