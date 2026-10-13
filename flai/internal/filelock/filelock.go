// Package filelock takes an exclusive lock on a file, held by the operating
// system for the open file, so that the lock goes when its holder closes the
// file or dies (I-0140). A lock is never broken for its age: a lock made by
// creating a file and broken when the file is old lets two writers that both
// find it old both take it, and one loses its write.
//
// On unix it is flock(2) and on Windows LockFileEx, each tried without
// blocking and tried again every poll until the caller's wait has passed. A
// blocking lock in a goroutine was the alternative, but it cannot be called
// off at a deadline: the goroutine and its open file would stay behind until
// the holder let go.
//
// The lock belongs to the open file, not to the process: two goroutines of
// one process that each call Lock exclude each other as two processes do. Go
// opens files close-on-exec, so a program the holder starts does not inherit
// the lock.
package filelock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// poll is how long Lock sleeps between tries while another holds the lock.
const poll = 10 * time.Millisecond

// Lock takes an exclusive lock on the file at path, creating the file and
// its folder when missing, waiting up to wait while another holds it, or
// with wait <= 0 until it is free; unlock releases it and leaves the file.
func Lock(path string, wait time.Duration) (unlock func(), err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("lock %s: create its folder: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("lock %s: open it: %w", path, err)
	}
	deadline := time.Now().Add(wait)
	for {
		held, err := tryLock(f)
		if err != nil {
			_ = f.Close()
			return nil, fmt.Errorf("lock %s: %w", path, err)
		}
		if held {
			// The file stays: removing it would let a waiter that has it
			// open lock a file no longer at path while a newcomer locks a
			// new one there, and both would go ahead.
			return func() {
				_ = unlockFile(f)
				_ = f.Close()
			}, nil
		}
		if wait > 0 && !time.Now().Before(deadline) {
			_ = f.Close()
			return nil, fmt.Errorf("lock %s: another holder kept it for all of %s: try again once it is done, or stop the process that has the file open", path, wait)
		}
		time.Sleep(poll)
	}
}
