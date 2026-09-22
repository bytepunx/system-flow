//go:build !windows

package serve

import (
	"os"
	"os/exec"
	"syscall"
)

// Alive reports whether pid names a process this user can signal.
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	return err == nil && p.Signal(syscall.Signal(0)) == nil
}

// Detach makes cmd outlive the terminal and the process that started it.
func Detach(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} }

// Terminate asks the process to stop.
func Terminate(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGTERM)
}

// TerminateGroup asks the process group led by pid to stop (S-0082): Detach
// makes a detached process its own session and group leader, so a signal to
// the negative of its PID reaches it and everything it spawned, not just
// the one process.
func TerminateGroup(pid int) error {
	if pid <= 0 {
		return os.ErrInvalid
	}
	return syscall.Kill(-pid, syscall.SIGTERM)
}

// KillGroup forces the process group led by pid to stop.
func KillGroup(pid int) error {
	if pid <= 0 {
		return os.ErrInvalid
	}
	return syscall.Kill(-pid, syscall.SIGKILL)
}
