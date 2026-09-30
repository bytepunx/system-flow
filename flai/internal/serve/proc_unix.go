//go:build !windows

package serve

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Alive reports whether pid names a process this user can signal.
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	return err == nil && p.Signal(syscall.Signal(0)) == nil
}

// Owns reports whether pid is still the process Detach started at started
// (S-0170): alive, the leader of its own session, and, where the system says
// when a process started, started within startSlack of then. A PID the
// system has given to another process since, after a reboot, is not, and is
// never signalled for the run that had it.
func Owns(pid int, started time.Time) bool {
	if !Alive(pid) {
		return false
	}
	sid, at, err := process(pid)
	if err != nil || sid != pid {
		return false
	}
	if at.IsZero() {
		return true
	}
	d := at.Sub(started)
	return d > -startSlack && d < startSlack
}

// startSlack is how far the start the system reports for a process may be
// from the start flai recorded: the clock the system counts from can drift
// from the wall clock, and a WSL host's does after a sleep.
const startSlack = 10 * time.Minute

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
