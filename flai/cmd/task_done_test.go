package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/taskdone"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// taskDoneClock is a day after the good fixture's S-004 started.
var taskDoneClock = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

// taskDoneProject is the good fixture in a git repository whose main branch
// holds it, with S-004's branch checked out in its worktree; it returns the
// main checkout and the worktree. Integration: it runs real git.
func taskDoneProject(t *testing.T) (root, wt string) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "agent-S-004")
	t.Setenv("FLAI_SESSION", "s1")
	root = t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	writeIn(t, root, ".gitignore", ".flai-cache/\n")
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "project")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	wt = repo.WorktreePath("S-004")
	gitIn(t, root, "worktree", "add", "-q", "-b", storygit.Branch("S-004"), wt, "main")
	return root, wt
}

func writeIn(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// otherStory is S-006, in progress with an empty claim, so that every path a
// commit of S-004 changes is told to it.
const otherStory = "---\nid: S-006\ntype: story\nnature: feature\ntitle: Six\nstatus: in-progress\nparent: E-001\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\n---\n\n# S-006 Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n"

// taskDone runs flai task done in the worktree at the fixture's clock.
func taskDone(t *testing.T, wt string, args ...string) (string, string, int) {
	t.Helper()
	return runInAt(t, wt, taskDoneClock, append([]string{"task", "done"}, args...)...)
}

func taskStatus(t *testing.T, root string) string {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get("T-003")
	if err != nil {
		t.Fatal(err)
	}
	return it.Status
}

// ADR-0107: flai task done commits, tells the other open stories whose claim
// covers the commit (S-0333), syncs, moves the task to done, logs, widens the
// touches, checks, and reads the inbox, a line for each; called again with
// --json it commits nothing, tells nobody, does not move the task again, and
// prints the result.
func TestTaskDoneClosesATaskAndAnswersTextAndJSON(t *testing.T) {
	root, wt := taskDoneProject(t)
	writeIn(t, root, "wip/kanban/stories/S-006-six.md", otherStory)
	writeIn(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")

	out, errOut, code := taskDone(t, wt, "T-003", "-m", "feat: [S-004] T-003 the readme\n\nWhy, at length.")
	if code != 0 {
		t.Fatalf("task done: %d\n%s\n%s", code, out, errOut)
	}
	head := gitIn(t, wt, "rev-parse", "HEAD")
	for _, want := range []string{
		"commit: " + short(head) + " feat: [S-004] T-003 the readme\ntold: S-006 (MS-0001) of README.md\nsync: story/S-004 is rebased onto main\n",
		"move: T-003 → done\n",
		": feat: [S-004] T-003 the readme\n",
		"touches: added to T-003: README.md\n",
		"touches: added to S-004: README.md\n",
		"check: passed",
		"inbox: ",
		"ready to pull, in pull order (can pull: ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "\nlog: S-004 at ") || strings.Contains(out, "stopped at") {
		t.Errorf("output:\n%s", out)
	}
	if st := taskStatus(t, root); st != workitem.Done {
		t.Errorf("T-003 is %s", st)
	}
	if dirty := gitIn(t, wt, "status", "--porcelain"); dirty != "" {
		t.Errorf("the worktree is not clean: %s", dirty)
	}

	out, errOut, code = taskDone(t, wt, "--json", "T-003", "-m", "fix: [S-004] nothing left", "--log", "checked it again")
	if code != 0 {
		t.Fatalf("task done --json: %d\n%s\n%s", code, out, errOut)
	}
	var res taskdone.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Task != "T-003" || res.Story != "S-004" || res.Stopped != "" || res.Commit != nil || res.Told == nil || len(res.Told) != 0 {
		t.Errorf("json: %+v", res)
	}
	if res.Sync == nil || !res.Sync.Synced || res.Move == nil || !res.Move.Skipped || res.Log == nil || res.Log.Entry != "checked it again" {
		t.Errorf("json: sync %+v, move %+v, log %+v", res.Sync, res.Move, res.Log)
	}
	if res.Check == nil || !res.Check.Passed || res.Inbox == nil || res.Inbox.Agent != "agent-S-004" {
		t.Errorf("json: check %+v, inbox %+v", res.Check, res.Inbox)
	}
}

// A sync that stops on a conflict with the main branch stops the run with
// exit 3, the paths, and how to continue, before the move. Called again with
// the rebase unfinished, the commit step refuses, with exit 1.
func TestTaskDoneStopsAtARefusedSync(t *testing.T) {
	root, wt := taskDoneProject(t)
	writeIn(t, root, "README.md", "# good\n\nChanged on main.\n")
	gitIn(t, root, "commit", "-q", "-am", "docs: main's readme")
	writeIn(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")

	out, errOut, code := taskDone(t, wt, "T-003", "-m", "feat: [S-004] T-003 the readme")
	if code != 3 {
		t.Fatalf("task done: %d, want 3\n%s\n%s", code, out, errOut)
	}
	for _, want := range []string{
		"commit: ",
		"sync: story/S-004 was not synced: the rebase onto main stopped on conflicts in 1 path:\n  README.md\n",
		"To continue: ",
		"git rebase --continue",
		"To abort: ",
		"stopped at sync: ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"move: ", "log: ", "check: ", "inbox: "} {
		if strings.Contains(out, not) {
			t.Errorf("a step after the sync ran (%q):\n%s", not, out)
		}
	}
	if !strings.Contains(errOut, "T-003 stopped at the sync step") {
		t.Errorf("error: %s", errOut)
	}
	if st := taskStatus(t, root); st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
	}

	out, errOut, code = taskDone(t, wt, "--json", "T-003", "-m", "feat: [S-004] T-003 the readme")
	if code != 1 {
		t.Fatalf("task done during the rebase: %d, want 1\n%s\n%s", code, out, errOut)
	}
	var res taskdone.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Stopped != taskdone.StepCommit || res.Sync != nil || !strings.Contains(res.Error, "git rebase --continue") {
		t.Errorf("json: stopped %q: %s", res.Stopped, res.Error)
	}
}

// A check that fails in the story stops the run with exit 4 and each
// finding, after the move and the log, and before the inbox.
func TestTaskDoneStopsAtAFailedCheck(t *testing.T) {
	root, wt := taskDoneProject(t)
	writeIn(t, wt, "docs/notes.md", "# Notes\n\nNo front matter.\n")

	out, errOut, code := taskDone(t, wt, "T-003", "-m", "docs: [S-004] T-003 notes")
	if code != 4 {
		t.Fatalf("task done: %d, want 4\n%s\n%s", code, out, errOut)
	}
	for _, want := range []string{
		"move: T-003 → done\n",
		"touches: added to T-003: docs/notes.md\n",
		"check: failed with ",
		"docs/notes.md:",
		": doc.front-matter: ",
		"stopped at check: ",
		"flai task done T-003 again",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "inbox: ") {
		t.Errorf("the inbox was read after a failed check:\n%s", out)
	}
	if !strings.Contains(errOut, "T-003 stopped at the check step") {
		t.Errorf("error: %s", errOut)
	}
	if st := taskStatus(t, root); st != workitem.Done {
		t.Errorf("T-003 is %s", st)
	}
}

// -m is required, and a story or an epic is refused with what moves it on
// instead; neither needs a project.
func TestTaskDoneRefusesWithoutAMessageOrATask(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := t.TempDir()
	for _, tc := range []struct {
		name string
		args []string
		want []string
	}{
		{"no message", []string{"T-3"}, []string{"required flag(s)", "message", "not set"}},
		{"a story", []string{"S-4", "-m", "m"}, []string{"S-0004 is a story", "flai move S-0004 review after its close-out"}},
		{"an epic", []string{"E-1", "-m", "m"}, []string{"E-0001 is an epic", "follows its stories"}},
	} {
		out, errOut, code := taskDone(t, dir, tc.args...)
		if code != 1 || out != "" {
			t.Errorf("%s: %d %q", tc.name, code, out)
		}
		for _, want := range tc.want {
			if !strings.Contains(errOut, want) {
				t.Errorf("%s: error lacks %q: %s", tc.name, want, errOut)
			}
		}
	}
}
