package taskdone

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// clock is the good fixture's check time, a day after S-004 started.
var clock = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

const fixture = "../metrics/testdata/good"

// copyFixture copies the good fixture, whose S-004 is in progress with its
// narrative and its task T-003 in progress, to a new folder.
func copyFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(fixture)); err != nil {
		t.Fatal(err)
	}
	return dir
}

// project is the good fixture in a git repository whose main branch holds
// it, with S-004's branch checked out in its worktree, which it returns.
// Integration: it runs real git.
func project(t *testing.T) (*workitem.Repo, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	dir := copyFixture(t)
	write(t, dir, ".gitignore", ".flai-cache/\n")
	git(t, dir, "init", "-q", "-b", "main")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "project")
	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	wt := repo.WorktreePath("S-004")
	git(t, dir, "worktree", "add", "-q", "-b", storygit.Branch("S-004"), wt, "main")
	return repo, wt
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := (execx.System{}).Run(dir, "git", args...)
	if err != nil {
		t.Fatal(out, err)
	}
	return out
}

// options close T-003 with message as agent-S-004 at the clock's time.
func options(repo *workitem.Repo, message string) Options {
	return Options{Repo: repo, Runner: execx.System{}, Task: "T-003", Message: message, Agent: "agent-S-004", Session: "s1", Now: func() time.Time { return clock }}
}

func get(t *testing.T, repo *workitem.Repo, id string) *workitem.Item {
	t.Helper()
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return it
}

func narrative(t *testing.T, repo *workitem.Repo) string {
	t.Helper()
	data, err := os.ReadFile(repo.NarrativePath("S-004"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// ADR-0107: one call commits, syncs, moves the task to done, logs the
// subject, widens the task's and the story's touches with what each does
// not cover, checks the story clean, and answers the inbox. Called again on
// the done task, it commits nothing, skips the move, and runs the rest.
func TestRunClosesATaskInOneCall(t *testing.T) {
	repo, wt := project(t)
	task := get(t, repo, "T-003")
	task.Touches = []string{"README.md"}
	if err := repo.Save(task); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	write(t, wt, "docs/README.md", "# Docs\n\nChanged too.\n")

	res, err := Run(context.Background(), options(repo, "feat: [S-004] T-003 the readmes\n\nWhy, at length."))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Error != "" || res.ExitCode() != 0 || res.Task != "T-003" || res.Story != "S-004" {
		t.Fatalf("stopped at %q: %s; check %+v", res.Stopped, res.Error, res.Check)
	}
	if res.Commit == nil || res.Commit.Subject != "feat: [S-004] T-003 the readmes" || strings.Join(res.Commit.Paths, ",") != "README.md,docs/README.md" {
		t.Fatalf("commit: %+v", res.Commit)
	}
	if head := git(t, wt, "rev-parse", "HEAD"); res.Commit.Hash != head {
		t.Errorf("commit hash %s, HEAD %s", res.Commit.Hash, head)
	}
	if dirty := git(t, wt, "status", "--porcelain"); dirty != "" {
		t.Errorf("the worktree is not clean: %s", dirty)
	}
	if res.Sync == nil || !res.Sync.Synced {
		t.Errorf("sync: %+v", res.Sync)
	}
	if res.Move == nil || res.Move.State != workitem.Done || res.Move.Skipped || get(t, repo, "T-003").Status != workitem.Done {
		t.Errorf("move: %+v", res.Move)
	}
	if res.Log == nil || res.Log.Entry != "feat: [S-004] T-003 the readmes" || res.Log.Stream != "S-004" {
		t.Errorf("log: %+v", res.Log)
	}
	if n := narrative(t, repo); !strings.Contains(n, "### 2026-09-01T12:00:00Z\nfeat: [S-004] T-003 the readmes\n") || strings.Contains(n, "Why, at length") {
		t.Errorf("narrative:\n%s", n)
	}
	// the task covered README.md already; the story covered nothing
	if res.Touches == nil || strings.Join(res.Touches.Task, ",") != "docs/README.md" || strings.Join(res.Touches.Story, ",") != "README.md,docs/README.md" {
		t.Errorf("touches: %+v", res.Touches)
	}
	if got := get(t, repo, "T-003").Touches; strings.Join(got, ",") != "README.md,docs/README.md" {
		t.Errorf("T-003 touches %v", got)
	}
	if got := get(t, repo, "S-004").Touches; strings.Join(got, ",") != "README.md,docs/README.md" {
		t.Errorf("S-004 touches %v", got)
	}
	if res.Check == nil || !res.Check.Passed || res.Inbox == nil || res.Inbox.Agent != "agent-S-004" {
		t.Errorf("check %+v, inbox %+v", res.Check, res.Inbox)
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"stopped":""`, `"followed":[]`, `"warnings":[]`, `"overlaps":[]`, `"notes":[`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("json lacks %s: %s", want, data)
		}
	}

	again, err := Run(context.Background(), options(repo, "fix: [S-004] T-003 nothing left"))
	if err != nil {
		t.Fatal(err)
	}
	if again.Stopped != "" || again.Commit != nil || again.Move == nil || !again.Move.Skipped || again.Move.State != workitem.Done {
		t.Fatalf("again: stopped %q %s, commit %+v, move %+v", again.Stopped, again.Error, again.Commit, again.Move)
	}
	if len(again.Touches.Task) != 0 || len(again.Touches.Story) != 0 || again.Inbox == nil {
		t.Errorf("again: touches %+v, inbox %+v", again.Touches, again.Inbox)
	}
	if !strings.Contains(narrative(t, repo), "\nfix: [S-004] T-003 nothing left\n") {
		t.Error("the second call did not log")
	}
}

// A sync that stops on a conflict with the main branch stops the run before
// the move, with the paths and how to continue; the task stays in progress
// and nothing is logged. Called again with the rebase unfinished, the commit
// step refuses.
func TestRunStopsAtARefusedSyncBeforeTheMove(t *testing.T) {
	repo, wt := project(t)
	write(t, repo.Root, "README.md", "# good\n\nChanged on main.\n")
	git(t, repo.Root, "commit", "-q", "-am", "docs: main's readme")
	write(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	before := narrative(t, repo)

	res, err := Run(context.Background(), options(repo, "feat: [S-004] T-003 the readme"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != StepSync || res.ExitCode() != 3 || res.Commit == nil {
		t.Fatalf("stopped at %q: %s, commit %+v", res.Stopped, res.Error, res.Commit)
	}
	if res.Sync == nil || res.Sync.Stopped != storygit.StopConflicts || strings.Join(res.Sync.Conflicts, ",") != "README.md" || !strings.Contains(res.Sync.Continue, "flai task done T-003 again") {
		t.Errorf("sync: %+v", res.Sync)
	}
	if !strings.Contains(res.Error, "README.md") || !strings.Contains(res.Error, "git rebase --continue") {
		t.Errorf("error: %s", res.Error)
	}
	if res.Move != nil || res.Log != nil || res.Touches != nil || res.Check != nil || res.Inbox != nil {
		t.Errorf("steps after the sync ran: %+v", res)
	}
	if st := get(t, repo, "T-003").Status; st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
	}
	if narrative(t, repo) != before {
		t.Error("the narrative was logged to")
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"stopped":"sync"`, `"move":null`, `"inbox":null`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("json lacks %s: %s", want, data)
		}
	}

	res, err = Run(context.Background(), options(repo, "feat: [S-004] T-003 the readme"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != StepCommit || res.ExitCode() != 1 || res.Commit != nil || res.Sync != nil || !strings.Contains(res.Error, "git rebase --continue") {
		t.Errorf("with the rebase unfinished: stopped %q: %s", res.Stopped, res.Error)
	}
	if st := get(t, repo, "T-003").Status; st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
	}
}

// A check that fails inside the story stops the run with its findings after
// the move and the log: the task is done, and the inbox is not read.
func TestRunStopsAtAFailedCheckAfterTheMoveAndTheLog(t *testing.T) {
	repo, wt := project(t)
	write(t, wt, "docs/notes.md", "# Notes\n\nNo front matter.\n")

	res, err := Run(context.Background(), options(repo, "docs: [S-004] T-003 notes"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != StepCheck || res.ExitCode() != 4 {
		t.Fatalf("stopped at %q: %s", res.Stopped, res.Error)
	}
	if res.Move == nil || res.Move.State != workitem.Done || get(t, repo, "T-003").Status != workitem.Done {
		t.Errorf("move: %+v", res.Move)
	}
	if res.Log == nil || !strings.Contains(narrative(t, repo), "\ndocs: [S-004] T-003 notes\n") {
		t.Errorf("log: %+v", res.Log)
	}
	if res.Touches == nil || !slices.Contains(res.Touches.Task, "docs/notes.md") {
		t.Errorf("touches: %+v", res.Touches)
	}
	if res.Check == nil || res.Check.Passed {
		t.Fatalf("check: %+v", res.Check)
	}
	found := slices.ContainsFunc(res.Check.Findings, func(f check.Finding) bool {
		return f.Rule == "doc.front-matter" && strings.HasSuffix(filepath.ToSlash(f.Path), "docs/notes.md")
	})
	if !found {
		t.Errorf("findings: %+v", res.Check.Findings)
	}
	if res.Inbox != nil || !strings.Contains(res.Error, "flai task done T-003 again") {
		t.Errorf("inbox %+v, error %s", res.Inbox, res.Error)
	}
}

// Run refuses, before any step, what it cannot close: no message, an item
// that is not a task, and a task whose story has no worktree.
func TestRunRefusesWhatItCannotStart(t *testing.T) {
	repo, err := workitem.Open(copyFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, task, message, want string
	}{
		{"no message", "T-003", " ", "a commit message is needed"},
		{"a story", "S-004", "m", "S-004 is a story"},
		{"no worktree", "T-003", "m", "S-004 has no worktree at .flai-cache/worktrees/S-004; open one with flai stream open S-004"},
	} {
		o := options(repo, tc.message)
		o.Task = tc.task
		res, err := Run(context.Background(), o)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: %v, want %q", tc.name, err, tc.want)
		}
		if res.Stopped != "" || res.Commit != nil {
			t.Errorf("%s: a step ran: %+v", tc.name, res)
		}
	}
	if get(t, repo, "T-003").Status != workitem.InProgress {
		t.Error("T-003 moved")
	}
}

// ADR-0107's exit statuses: 3 for the sync, 4 for the check, 1 for any other
// step, and 0 when every step ran.
func TestExitCode(t *testing.T) {
	for stopped, want := range map[string]int{"": 0, StepCommit: 1, StepSync: 3, StepMove: 1, StepLog: 1, StepTouches: 1, StepCheck: 4, StepInbox: 1} {
		if got := (Result{Stopped: stopped}).ExitCode(); got != want {
			t.Errorf("stopped at %q: %d, want %d", stopped, got, want)
		}
	}
}

// S-0331: the inbox a closed task answers lists the open conversations of the
// task's story under messages, as the MCP tool inbox does for its agent.
func TestRunAnswersTheStorysMessages(t *testing.T) {
	repo, _ := project(t)
	write(t, repo.Root, "wip/kanban/stories/S-006-six.md", "---\nid: S-006\ntype: story\nnature: feature\ntitle: Six\nstatus: in-progress\nparent: E-001\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\n---\n\n# S-006 Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n")
	if _, err := messages.Send(repo, messages.SendOptions{From: "S-006", To: "S-004", Author: "agent-S-006", Text: "Who changes the readme?", Now: clock}); err != nil {
		t.Fatal(err)
	}

	res, err := Run(context.Background(), options(repo, "fix: [S-004] T-003 nothing to commit"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Inbox == nil || len(res.Inbox.Messages) != 1 {
		t.Fatalf("stopped %q %s, check %+v, inbox %+v", res.Stopped, res.Error, res.Check, res.Inbox)
	}
	if m := res.Inbox.Messages[0]; m.With != "S-006" || m.Awaiting != "you" || m.Last.Text != "Who changes the readme?" {
		t.Errorf("message %+v", m)
	}
}
