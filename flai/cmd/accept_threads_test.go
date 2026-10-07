package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0073: an acceptance resolves the threads still open or answered on what
// it archives (the story, its tasks, and the epic that followed it), in the
// acceptance commit, so flai check does not warn threads.archived after it.
// A thread on another item is left as it is, and the dry run names the
// threads and changes nothing.
func TestAcceptResolvesTheThreadsOnWhatItArchives(t *testing.T) {
	root := storyInReviewWithThreads(t)
	status, head := gitIn(t, root, "status", "--porcelain"), gitIn(t, root, "rev-parse", "HEAD")
	paths := threadPaths(t, root)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if got := strings.Join(res.ResolvedThreads, " "); got != "TH-0001 TH-0002 TH-0003" || len(res.Blockers) != 0 {
		t.Errorf("the dry run names the threads it would resolve, not as blockers: %q %v", got, res.Blockers)
	}
	if out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(out, "would resolve TH-0001, TH-0002, TH-0003, open or answered on what it archives\n") {
		t.Errorf("the dry run's text names them:\n%s", out)
	}
	if gitIn(t, root, "status", "--porcelain") != status || gitIn(t, root, "rev-parse", "HEAD") != head {
		t.Error("the dry run changes nothing")
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("the acceptance: %d %s\n%s", code, errOut, out)
	}
	if !strings.Contains(out, "resolved TH-0001, TH-0002, TH-0003, open or answered on what it archived\n") {
		t.Errorf("the output names the threads resolved:\n%s", out)
	}
	assertResolvedByAcceptance(t, root, "TH-0001", "TH-0002", "TH-0003")
	if th := thread(t, root, "TH-0004"); th.Status != "open" {
		t.Errorf("a thread on another item is left open: %s", th.Status)
	}
	committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	for _, id := range []string{"TH-0001", "TH-0002", "TH-0003"} {
		if !strings.Contains(committed, paths[id]) {
			t.Errorf("%s is in the acceptance commit:\n%s", id, committed)
		}
	}
	if strings.Contains(committed, paths["TH-0004"]) {
		t.Errorf("a thread the acceptance did not resolve is not in its commit:\n%s", committed)
	}
	if s := gitIn(t, root, "status", "--porcelain"); s != "" {
		t.Errorf("everything is committed: %q", s)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "threads.archived") {
		t.Errorf("flai check finds no thread open on an archived item:\n%s", out)
	}
}

// An acceptance whose commit failed is finished by running it again, and the
// finish resolves a thread left open on what it archived, as an older flai's
// acceptance left them, into the commit it makes.
func TestFinishingAnAcceptanceResolvesTheThreadsLeftOpen(t *testing.T) {
	root := storyInReviewWithThreads(t)
	paths := threadPaths(t, root)
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	_ = os.MkdirAll(filepath.Dir(hook), 0o755)
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 {
		t.Fatal("the acceptance whose commit failed exits non-zero")
	}
	// TH-0002 as the acceptance found it, open on the archived task
	gitIn(t, root, "checkout", "HEAD", "--", paths["TH-0002"])
	if th := thread(t, root, "TH-0002"); th.Status != "open" {
		t.Fatalf("TH-0002 is open again: %s", th.Status)
	}
	if out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(out, "would resolve TH-0002, open or answered on what it archives\n") {
		t.Errorf("the dry run names the thread the finish would resolve:\n%s", out)
	}

	_ = os.Remove(hook)
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001") || !strings.Contains(out, "resolved TH-0002, open or answered on what it archived\n") {
		t.Fatalf("run again, the acceptance resolves the thread and commits: %d %s %s", code, out, errOut)
	}
	assertResolvedByAcceptance(t, root, "TH-0001", "TH-0002", "TH-0003")
	committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, paths["TH-0002"]) {
		t.Errorf("TH-0002 is in the acceptance commit:\n%s", committed)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "threads.archived") {
		t.Errorf("flai check finds no thread open on an archived item:\n%s", out)
	}
}

// storyInReviewWithThreads is orchestratedStoryInReview's project, whose
// S-0001 is its epic E-0001's only story, with S-0002 under no epic and these
// threads, committed: TH-0001 alex's, answered on S-0001, and claude's open
// ones, TH-0002 on its task T-0001, TH-0003 on E-0001, and TH-0004 on S-0002.
// The clock is fixed, so an entry by alex after alex's would join it.
func storyInReviewWithThreads(t *testing.T) string {
	t.Helper()
	root, _ := orchestratedStoryInReview(t, false)
	for _, args := range [][]string{
		{"story", "new", "Other"},
		{"thread", "new", "--on", "S-0001", "--by", "alex", "Is the flag right?", "Say which."},
		{"thread", "new", "--on", "T-0001", "--by", "claude", "Is the task done?", "Say so."},
		{"thread", "new", "--on", "E-0001", "--by", "claude", "Is the epic whole?", "Say so."},
		{"thread", "new", "--on", "S-0002", "--by", "claude", "When is it next?", "Say when."},
		{"thread", "reply", "TH-0001", "--by", "claude", "Use the flag."},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	if th := thread(t, root, "TH-0001"); th.Status != "answered" {
		t.Fatalf("TH-0001 is answered: %s", th.Status)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: threads")
	return root
}

// assertResolvedByAcceptance fails unless each thread is resolved by alex
// with S-0001's acceptance as the reason.
func assertResolvedByAcceptance(t *testing.T, root string, ids ...string) {
	t.Helper()
	for _, id := range ids {
		th := thread(t, root, id)
		es := th.Entries()
		last := es[len(es)-1]
		if th.Status != "resolved" || last.Author != "alex" || last.Text != "Resolved: S-0001 was accepted" {
			t.Errorf("%s is resolved by the acceptance: status %s, last entry %+v", id, th.Status, last)
		}
	}
}

// threadPaths are the threads' files relative to root, by ID.
func threadPaths(t *testing.T, root string) map[string]string {
	t.Helper()
	paths := map[string]string{}
	for _, id := range []string{"TH-0001", "TH-0002", "TH-0003", "TH-0004"} {
		paths[id] = filepath.ToSlash(mustRel(t, root, thread(t, root, id).Path))
	}
	return paths
}

// thread reads thread id of the project at root.
func thread(t *testing.T, root, id string) *threads.Thread {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	th, err := threads.Get(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return th
}
