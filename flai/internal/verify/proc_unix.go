//go:build !windows

package verify

import (
	"os"
	"os/exec"
	"syscall"
)

// ownGroup starts cmd as the leader of a process group of its own.
func ownGroup(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} }

// terminateGroup asks the process group led by pid to stop.
func terminateGroup(pid int) error {
	if pid <= 0 {
		return os.ErrInvalid
	}
	return syscall.Kill(-pid, syscall.SIGTERM)
}

// killGroup forces the process group led by pid to stop.
func killGroup(pid int) error {
	if pid <= 0 {
		return os.ErrInvalid
	}
	return syscall.Kill(-pid, syscall.SIGKILL)
}
