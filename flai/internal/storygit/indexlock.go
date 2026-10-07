package storygit

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// indexLockWaits is how long RunPastIndexLock waits before each further run
// while the index lock is held: growing from 100ms, about nine seconds in all.
// Another git process (flai serve's replan, an agent's commit, another
// acceptance) holds a checkout's index lock for well under a second, and a git
// command that writes the index meanwhile fails at once (I-0100). Tests
// shorten it.
var indexLockWaits = []time.Duration{
	100 * time.Millisecond, 200 * time.Millisecond, 400 * time.Millisecond,
	800 * time.Millisecond, 1600 * time.Millisecond,
	2 * time.Second, 2 * time.Second, 2 * time.Second,
}

// indexLockHeld is what git prints when another process holds the index
// lock: "fatal: Unable to create '<path>/index.lock': File exists."
const indexLockHeld = "index.lock': File exists"

// RunPastIndexLock runs git with args in dir, running it again while another process holds the index lock.
func RunPastIndexLock(r execx.Runner, dir string, args ...string) (string, error) {
	out, err := r.Run(dir, "git", args...)
	var waited time.Duration
	for _, wait := range indexLockWaits {
		if !lockHeld(out, err) {
			return out, err
		}
		// The lock file is never removed here: it may belong to a live process.
		time.Sleep(wait)
		waited += wait
		out, err = r.Run(dir, "git", args...)
	}
	if lockHeld(out, err) {
		return out, fmt.Errorf("%w\nanother git process held the index lock in %s throughout the %s flai waited; run the command again once that process has finished, or remove the index.lock file git names if no git process is running", err, dir, waited)
	}
	return out, err
}

// lockHeld reports whether a git run failed because the index lock was held.
func lockHeld(out string, err error) bool {
	return err != nil && strings.Contains(out, indexLockHeld)
}
