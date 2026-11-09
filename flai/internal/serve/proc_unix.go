//go:build !windows

package serve

import (
	"os"
	"os/exec"
	"syscall"
)

// Alive reports whether pid names a process this user can signal and that
// has not ended: a zombie, which signal 0 still reaches until its parent
// reaps it, is not alive where the system says so (I-0299).
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil || p.Signal(syscall.Signal(0)) != nil {
		return false
	}
	_, _, ended, err := process(pid)
	return err != nil || !ended
}

// Owns reports whether pid is still the process Detach started (S-0170):
// alive, not a zombie, the leader of its own session, and, when start is
// known (Started), started then. A PID the system has given to another
// process since, after a reboot, is not, and is never signalled for the run
// that had it.
func Owns(pid int, start int64) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil || p.Signal(syscall.Signal(0)) != nil {
		return false
	}
	sid, at, ended, err := process(pid)
	return err == nil && !ended && sid == pid && (start == 0 || at == 0 || at == start)
}

// Started is when pid started as the system counts it, which no two
// processes of one boot share with one PID: 0 where it is not known.
func Started(pid int) int64 {
	_, at, _, err := process(pid)
	if err != nil {
		return 0
	}
	return at
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
