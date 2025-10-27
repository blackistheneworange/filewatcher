//go:build !windows
// +build !windows

package process

import (
	"os/exec"
	"errors"
	"strings"
	"os"
	"fmt"
)

type OsProcess struct {
	Process
	execCommandParts	[]string
}

func _CreateProcess(execCmd string) Process {
	process := &OsProcess{ execCommandParts: strings.Split(execCmd," ") }
	return process
}

func (process *OsProcess) StartProcess() error {
	if len(process.execCommandParts) == 0 {
		return errors.New("process error: Process not created")
	}
	if err := process.StopProcess(); err != nil {
		return errors.New(fmt.Sprintf("process error: %s", err.Error()))
	}
	
	cmd := exec.Command(process.execCommandParts[0], process.execCommandParts[1:]...)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("process error: %s", err.Error()))
	}

	return nil
}

func (process *OsProcess) StopProcess() error {
	return errors.New("process error: StopProcess not supported")
}