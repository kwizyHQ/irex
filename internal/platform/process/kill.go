package process

import "os/exec"

// killProcessTree is implemented in OS-specific files.
func KillProcessTree(pid int) error { return killProcessTree(pid) }

func ConfigureProcessAttributes(cmd *exec.Cmd) {
	configureProcessAttributes(cmd)
}
