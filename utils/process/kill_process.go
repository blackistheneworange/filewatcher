//go:build !windows
// +build !windows

package process

import (
	"syscall"
	"os/exec"
)

func KillProcess(cmd *exec.Cmd) error {
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		// force kill
		cmd.Process.Kill()
	}

	_, err := cmd.Process.Wait()

	return err
}