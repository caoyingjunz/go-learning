package main

import (
	"fmt"
	osexec "os/exec"
)

// ErrExecutableNotFound is returned if the executable is not found.
var ErrExecutableNotFound = osexec.ErrNotFound

type Interface interface {
	Command(cmd string, args ...string) Cmd

	// LookPath wraps os/exec.LookPath
	LookPath(file string) (string, error)
}

type Cmd interface {
	Run() error
	CombinedOutput() ([]byte, error)
}

func NewExec() (Interface, error) {
	execer := New()
	if _, err := execer.LookPath("python3"); err != nil {
		return nil, fmt.Errorf("%s is required for sail runtime", "python3")
	}

	return execer, nil
}

// Implements Interface in terms of really exec()ing.
type executor struct{}

// New returns a new Interface which will os/exec to run commands.
func New() Interface {
	return &executor{}
}

// Command is part of the Interface interface.
func (executor *executor) Command(cmd string, args ...string) Cmd {
	return osexec.Command(cmd, args...)
}

// LookPath is part of the Interface interface
func (executor *executor) LookPath(file string) (string, error) {
	return osexec.LookPath(file)
}

func main() {
	exec, err := NewExec()
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("ls", "-al").CombinedOutput()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
