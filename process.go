package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type IProcess interface {
	Execute()
}

type Process struct {
	execCommandParts []string
	restart          bool
}

func (process *Process) Execute() error {
	cmd := exec.Command(process.execCommandParts[0], process.execCommandParts[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func createProcess(execCmd *string) (*Process, error) {
	if *execCmd == "" {
		return nil, errors.New("process error: No executable (-exec) file path provided")
	}
	execCommandParts := strings.Split(*execCmd, " ")
	process := &Process{execCommandParts: execCommandParts}

	return process, nil
}
