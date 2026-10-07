package storygit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// lockedGit answers each git run with the next of its outputs, failing as
// execx.System does when the output is git's index-lock message or "fatal".
type lockedGit struct {
	outs []string
	runs int
}

func (g *lockedGit) Run(_, name string, args ...string) (string, error) {
	out := g.outs[min(g.runs, len(g.outs)-1)]
	g.runs++
	if strings.Contains(out, "fatal") {
		return out, fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), errors.New("exit status 128"), out)
	}
	return out, nil
}

func (g *lockedGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *lockedGit) LookPath(name string) (string, error) { return name, nil }

const lockedOut = "fatal: Unable to create '/repo/.git/index.lock': File exists.\n\n" +
	"Another git process seems to be running in this repository, e.g.\n" +
	"an editor opened by 'git commit'. Please make sure all processes\n" +
	"are terminated then try again."

// shortWaits makes RunPastIndexLock wait n times, each for d.
func shortWaits(t *testing.T, n int, d time.Duration) {
	t.Helper()
	saved := indexLockWaits
	indexLockWaits = make([]time.Duration, n)
	for i := range indexLockWaits {
		indexLockWaits[i] = d
	}
	t.Cleanup(func() { indexLockWaits = saved })
}

func TestRunPastIndexLockRunsASucceedingCommandOnce(t *testing.T) {
	shortWaits(t, 3, time.Millisecond)
	g := &lockedGit{outs: []string{"ok"}}
	out, err := RunPastIndexLock(g, "/repo", "add", "-A")
	if err != nil || out != "ok" || g.runs != 1 {
		t.Fatalf("out %q, err %v, runs %d; want ok, nil, 1", out, err, g.runs)
	}
}

// I-0100: a commit that meets another process's index lock is run again
// once the lock is gone.
func TestRunPastIndexLockRunsAgainWhileTheLockIsHeld(t *testing.T) {
	shortWaits(t, 3, time.Millisecond)
	g := &lockedGit{outs: []string{lockedOut, lockedOut, "committed"}}
	out, err := RunPastIndexLock(g, "/repo", "commit", "-q", "-m", "x")
	if err != nil || out != "committed" || g.runs != 3 {
		t.Fatalf("out %q, err %v, runs %d; want committed, nil, 3", out, err, g.runs)
	}
}

func TestRunPastIndexLockGivesUpWithGitsOutputWhenTheLockStaysHeld(t *testing.T) {
	shortWaits(t, 3, time.Millisecond)
	g := &lockedGit{outs: []string{lockedOut}}
	_, err := RunPastIndexLock(g, "/repo", "commit", "-q", "-m", "x")
	if err == nil {
		t.Fatal("no error while the lock stayed held")
	}
	if g.runs != 4 {
		t.Errorf("runs = %d, want 4: the first and one after each wait", g.runs)
	}
	for _, want := range []string{lockedOut, "held the index lock in /repo throughout the 3ms flai waited"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error lacks %q:\n%v", want, err)
		}
	}
}

func TestRunPastIndexLockReturnsAnyOtherFailureAtOnce(t *testing.T) {
	shortWaits(t, 3, time.Millisecond)
	other := "fatal: not a git repository (or any of the parent directories): .git"
	g := &lockedGit{outs: []string{other, "ok"}}
	_, err := RunPastIndexLock(g, "/repo", "add", "-A")
	if err == nil || !strings.Contains(err.Error(), other) || strings.Contains(err.Error(), "index lock") {
		t.Fatalf("err = %v, want git's own failure alone", err)
	}
	if g.runs != 1 {
		t.Errorf("runs = %d, want 1", g.runs)
	}
}

// Integration, I-0100 reproduced: another process holds the index lock when
// the commit starts and lets it go 300ms later, well after git's first try
// has failed; the commit waits for it and lands, and the lock file is left
// for its holder to remove.
func TestRunPastIndexLockCommitsOnceARealLockIsReleased(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	shortWaits(t, 100, 20*time.Millisecond)
	write(t, dir, "a.md", "two\n")
	git(t, dir, "add", "-A")
	lock := filepath.Join(dir, ".git", "index.lock")
	if err := os.WriteFile(lock, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	released := make(chan error, 1)
	go func() {
		time.Sleep(300 * time.Millisecond)
		released <- os.Remove(lock)
	}()
	if out, err := RunPastIndexLock(r, dir, "commit", "-q", "-m", "second"); err != nil {
		t.Fatalf("commit: %v\n%s", err, out)
	}
	if err := <-released; err != nil {
		t.Fatalf("the holder could not remove its lock: %v", err)
	}
	if out, err := r.Run(dir, "git", "log", "-1", "--format=%s"); err != nil || out != "second" {
		t.Errorf("last commit = %q, %v; want second", out, err)
	}
}
