//go:build !windows

package host

import (
	"os"
	"syscall"
)

// alive reports whether pid names a process this user can signal.
func alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	return err == nil && p.Signal(syscall.Signal(0)) == nil
}

// terminate asks a child to stop.
func terminate(p *os.Process) error { return p.Signal(syscall.SIGTERM) }

// Reexec replaces this process with the flai at exe, same arguments and
// environment, keeping its PID: what flai host does after an upgrade.
func Reexec(exe string) error {
	return syscall.Exec(exe, append([]string{exe}, os.Args[1:]...), os.Environ())
}
