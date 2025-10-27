package process

import (
	"os/exec"
	"errors"
	"strings"
)

type ProcessManager struct {
	cmd *exec.Cmd
	execCommandParts []string
	process *Process
}

type Process interface {
	StartProcess() error
	StopProcess() error
}

func GetProcessManager(execCmd *string) (*ProcessManager, error) {
	if *execCmd == "" {
		return nil, errors.New("process error: No executable (-exec) file path provided")
	}

	execCommandParts := strings.Split(*execCmd, " ")

	return &ProcessManager{ execCommandParts: execCommandParts }, nil
}

func (pm *ProcessManager) CreateProcess(execCmd string) Process {
	return _CreateProcess(execCmd)
}