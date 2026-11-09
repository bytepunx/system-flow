//go:build !linux && !windows

package serve

import "syscall"

// process is pid's session; when it started, and whether it has ended and
// waits to be reaped, are not known here without a subprocess, so Owns goes
// by the session alone.
func process(pid int) (sid int, start int64, ended bool, err error) {
	sid, err = syscall.Getsid(pid)
	return sid, 0, false, err
}
