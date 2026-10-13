//go:build !windows

package filelock

import (
	"errors"
	"os"
	"syscall"
)

// tryLock takes an exclusive flock on f without waiting, and reports false
// when another open file holds it.
func tryLock(f *os.File) (bool, error) {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, syscall.EWOULDBLOCK), errors.Is(err, syscall.EINTR):
		return false, nil
	default:
		return false, err
	}
}

// unlockFile releases the flock f holds.
func unlockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
