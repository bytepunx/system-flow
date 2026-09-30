//go:build !linux && !windows

package serve

import "syscall"

// process is pid's session; when it started is not known here without a
// subprocess, so Owns goes by the session alone.
func process(pid int) (sid int, start int64, err error) {
	sid, err = syscall.Getsid(pid)
	return sid, 0, err
}
