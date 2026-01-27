//go:build !windows

package plan

import (
	"os/exec"
	"syscall"
)

func configureProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessTree(pid int) error {
	// negative pid targets the process group
	_ = syscall.Kill(-pid, syscall.SIGKILL)
	return nil
}
