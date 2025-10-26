package process

import (
	"os/exec"
	"errors"
	"strings"
	"os"
	"fmt"
)

type ProcessManager struct {
	cmd *exec.Cmd
	execCommandParts []string
}

func GetProcessManager(execCmd *string) (*ProcessManager, error) {
	if *execCmd == "" {
		return nil, errors.New("process error: No executable (-exec) file path provided")
	}

	execCommandParts := strings.Split(*execCmd, " ")

	return &ProcessManager{ execCommandParts: execCommandParts }, nil
}

func (pm *ProcessManager) StartProcess() error {
	if len(pm.execCommandParts) == 0 {
		return errors.New("process error: Process manager not created")
	}
	if err := pm.StopProcess(); err != nil {
		return errors.New(fmt.Sprintf("process error: %s", err.Error()))
	}
	
	cmd := exec.Command(pm.execCommandParts[0], pm.execCommandParts[1:]...)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("process error: %s", err.Error()))
	}

	pm.cmd = cmd

	return nil
}

func (pm *ProcessManager) StopProcess() error {
	if pm.cmd == nil || pm.cmd.Process == nil {
		return nil
	}

	// kill process
	KillProcess(pm.cmd)

	pm.cmd = nil

	return nil
}