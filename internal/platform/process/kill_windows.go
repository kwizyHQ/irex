//go:build windows

package process

import (
	"os/exec"
	"strconv"
	"syscall"
)

func configureProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func killProcessTree(pid int) error {
	// taskkill /PID <pid> /T /F kills the process tree on Windows
	return exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
}
