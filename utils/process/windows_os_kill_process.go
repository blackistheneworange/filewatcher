//go:build windows
// +build windows

package process

import (
	"os/exec"
	"strconv"
)

func KillProcess(cmd *exec.Cmd) error {
	processKiller := exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(cmd.Process.Pid))

	err := processKiller.Run()

	return err
}