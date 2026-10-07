package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// lockingRunner is the real runner with another git process holding the main
// checkout's index lock for a moment: the first git add -A and the first git
// commit asked of it each meet .git/index.lock, made where git makes it, and
// the next run of each finds it gone.
type lockingRunner struct {
	execx.System
	lock string
	runs map[string]int
}

func (l *lockingRunner) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && len(args) > 1 && ((args[0] == "add" && args[1] == "-A") || (args[0] == "commit" && args[1] == "-q")) {
		l.runs[args[0]]++
		if l.runs[args[0]] == 1 {
			if err := os.WriteFile(l.lock, nil, 0o644); err != nil {
				return "", err
			}
		} else if err := os.Remove(l.lock); err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	return l.System.Run(dir, name, args...)
}

// I-0100: an acceptance whose commit meets the index lock another git process
// holds for a moment waits it out and completes, with one acceptance commit.
func TestAnAcceptanceWaitsOutAnIndexLockHeldForAMoment(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	before := gitIn(t, root, "rev-parse", "HEAD")
	lock := filepath.Join(root, ".git", "index.lock")
	r := &lockingRunner{lock: lock, runs: map[string]int{}}

	out, errOut, code := runWithApp(t, &app{cwd: root, runner: r}, "accept", "S-0001")
	if code != 0 || !strings.Contains(out, "accepted S-0001") {
		t.Fatalf("the acceptance waits out the lock: %d %s %s", code, out, errOut)
	}
	if r.runs["add"] != 2 || r.runs["commit"] != 2 {
		t.Errorf("git add and git commit each meet the lock once and run again: %v", r.runs)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the acceptance commit: %q", subject)
	}
	if n := gitIn(t, root, "rev-list", "--count", before+"..HEAD"); n != "2" {
		t.Errorf("the merged story commit and one acceptance commit: %s", n)
	}
	if status := gitIn(t, root, "status", "--porcelain"); status != "" {
		t.Errorf("everything is committed: %q", status)
	}
	if _, err := os.Stat(lock); !os.IsNotExist(err) {
		t.Errorf("the lock is the other process's to remove, and it did: %v", err)
	}
}

// I-0100: an acceptance whose commit fails even so keeps git's error and says
// how to finish; the story stays done and archived, its changes staged, and
// flai accept commits them, under the usual subject, as the operator's.
func TestAnAcceptanceWhoseCommitFailedIsFinishedByRunningItAgain(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, true)
	// an open story whose claim covers what S-0001's commits changed
	if _, errOut, code := runIn(t, root, "story", "new", "Next", "--touches", "cli"); code != 0 {
		t.Fatal(errOut)
	}
	next := filepath.Join(root, "wip/kanban/stories/S-0002-next.md")
	s, _ := os.ReadFile(next)
	_ = os.WriteFile(next, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] it works\n", 1)), 0o644)
	for _, st := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", "S-0002", st); code != 0 {
			t.Fatal(errOut)
		}
	}
	// git's commit fails, as it does while another process holds the index
	// lock beyond flai's wait
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	_ = os.MkdirAll(filepath.Dir(hook), 0o755)
	if err := os.WriteFile(hook, []byte("#!/bin/sh\necho 'pre-commit refuses the commit' >&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("the acceptance whose commit failed exits non-zero")
	}
	for _, want := range []string{"pre-commit refuses the commit", "story/S-0001 is merged and S-0001 is done and archived", "its changes are staged", "run flai accept S-0001 to commit them"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the error names %q:\n%s", want, errOut)
		}
	}
	for _, glob := range []string{"stories/S-0001-*.md", "epics/E-0001-*.md"} {
		m, _ := filepath.Glob(filepath.Join(root, "wip/archive/kanban", glob))
		if len(m) != 1 {
			t.Fatalf("%s is archived: %v", glob, m)
		}
		if data, _ := os.ReadFile(m[0]); !strings.Contains(string(data), "status: done") {
			t.Errorf("%s is done:\n%s", glob, data)
		}
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "feat: [S-0001] ship it" {
		t.Errorf("the branch is merged and nothing committed after it: %q", subject)
	}
	if staged := gitIn(t, root, "diff", "--cached", "--name-only"); !strings.Contains(staged, "wip/archive/kanban/stories/S-0001-ship-it.md") {
		t.Errorf("the archive is staged:\n%s", staged)
	}

	// the orchestrator does not finish it: that is the operator's
	_, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--dry-run")
	if code == 0 || !strings.Contains(errOut, "done and archived but its acceptance was never committed") || !strings.Contains(errOut, "the operator completes this acceptance with flai accept S-0001") {
		t.Errorf("the orchestrator is refused: %d %s", code, errOut)
	}

	status, head := gitIn(t, root, "status", "--porcelain"), gitIn(t, root, "rev-parse", "HEAD")
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run")
	if code != 0 || !strings.Contains(out, "would commit the acceptance of S-0001, done and archived but not committed: chore: [S-0001] accept and archive, with E-0001") || !strings.Contains(out, "dry run: nothing changed") {
		t.Errorf("the dry run says it would commit: %d %s %s", code, out, errOut)
	}
	if gitIn(t, root, "status", "--porcelain") != status || gitIn(t, root, "rev-parse", "HEAD") != head {
		t.Error("the dry run changes nothing")
	}

	_ = os.Remove(hook)
	out, errOut, code = runIn(t, root, "accept", "S-0001", "--trailer", "Accepted-by: t")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001: done and archived already, committed") {
		t.Fatalf("run again, the acceptance commits: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "told S-0002 it overlaps: cli/feature.go") {
		t.Errorf("the open story is told the paths the story's commits changed:\n%s", out)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the usual subject: %q", subject)
	}
	if body := gitIn(t, root, "log", "-1", "--format=%b"); body != "Accepted-by: t" {
		t.Errorf("the trailer: %q", body)
	}
	if n := gitIn(t, root, "rev-list", "--count", head+"..HEAD"); n != "1" {
		t.Errorf("one acceptance commit: %s", n)
	}
	if status := gitIn(t, root, "status", "--porcelain"); status != "" {
		t.Errorf("everything is committed: %q", status)
	}

	// another change under wip left uncommitted is not the acceptance's
	if f, err := os.OpenFile(next, os.O_APPEND|os.O_WRONLY, 0); err == nil {
		_, _ = f.WriteString("\nmore\n")
		_ = f.Close()
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code == 0 || !strings.Contains(errOut, "S-0001 is already done") {
		t.Errorf("a third run is refused as already done: %d %s", code, errOut)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the third run commits nothing: %q", subject)
	}
}
