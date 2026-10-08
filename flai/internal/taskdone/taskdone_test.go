package taskdone

import (
	"context"
	"encoding/json"
	"log/slog"
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

// story writes a story to the main checkout, in progress or in review, whose
// touches are its claim.
func story(t *testing.T, repo *workitem.Repo, id, status string, touches ...string) {
	t.Helper()
	moves := "  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\n"
	if status == workitem.Review {
		moves += "  - to: review\n    at: 2026-08-31T11:00:00Z\n    by: agent\n"
	}
	write(t, repo.Root, "wip/kanban/stories/"+id+"-other.md", "---\nid: "+id+"\ntype: story\nnature: feature\ntitle: Other "+id+"\nstatus: "+status+"\nparent: E-001\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T11:00:00Z\ntransitions:\n"+moves+"tags: []\ntouches: ["+strings.Join(touches, ", ")+"]\n---\n\n# "+id+" Other "+id+"\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n")
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
	// no other story is open: nobody is told, and no conversation is written
	if res.Told == nil || len(res.Told) != 0 {
		t.Errorf("told: %+v", res.Told)
	}
	if _, err := os.Stat(messages.Dir(repo)); !os.IsNotExist(err) {
		t.Errorf("a conversation was written: %v", err)
	}
	for _, want := range []string{`"stopped":""`, `"left":[]`, `"told":[]`, `"followed":[]`, `"warnings":[]`, `"overlaps":[]`, `"notes":[`} {
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
// and nothing is logged, but the other open stories have been told of the
// commit, which a second call would not make again. Called again with the
// rebase unfinished, the commit step refuses.
func TestRunStopsAtARefusedSyncBeforeTheMove(t *testing.T) {
	repo, wt := project(t)
	write(t, repo.Root, "README.md", "# good\n\nChanged on main.\n")
	git(t, repo.Root, "commit", "-q", "-am", "docs: main's readme")
	story(t, repo, "S-006", workitem.InProgress)
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
	if len(res.Told) != 1 || res.Told[0].Story != "S-006" || res.Told[0].Conversation != "MS-0001" {
		t.Errorf("told before the sync: %+v", res.Told)
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
	if res.Stopped != StepCommit || res.ExitCode() != 1 || res.Commit != nil || res.Told != nil || res.Sync != nil || !strings.Contains(res.Error, "git rebase --continue") {
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

// Run refuses, before any step, what it cannot close: an item that is not a
// task, and a task whose story has no worktree.
func TestRunRefusesWhatItCannotStart(t *testing.T) {
	repo, err := workitem.Open(copyFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, task, message, want string
	}{
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

// newTask creates an in-progress task of S-004 in the main checkout whose
// touches are touches, and returns its ID.
func newTask(t *testing.T, repo *workitem.Repo, title string, touches ...string) string {
	t.Helper()
	it, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: title, Parent: "S-004", Touches: touches, Now: clock})
	if err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.TransitionAll(it, to, "agent", "", clock, false); err != nil {
			t.Fatal(err)
		}
	}
	return it.ID
}

// setTouches sets id's touches in the main checkout.
func setTouches(t *testing.T, repo *workitem.Repo, id string, touches ...string) {
	t.Helper()
	it := get(t, repo, id)
	it.Touches = touches
	if err := repo.Save(it); err != nil {
		t.Fatal(err)
	}
}

// committed is the paths the commit hash changed, joined with commas.
func committed(t *testing.T, dir, hash string) string {
	t.Helper()
	return strings.ReplaceAll(git(t, dir, "show", "--name-only", "--no-renames", "--format=", hash), "\n", ",")
}

// I-0104, ADR-0128: two open tasks of one story, each with its own touches
// and its own changes in the worktree, close apart. The first close commits
// its own paths, a staged deletion among them, leaves the second's
// uncommitted and lists them, goes on past the sync they refuse, and does
// not widen any touches with them; the second close commits the rest and
// syncs.
func TestRunClosesTwoTasksOfOneStoryApart(t *testing.T) {
	repo, wt := project(t)
	setTouches(t, repo, "T-003", "src/one", "docs")
	second := newTask(t, repo, "Second", "src/two.go", "README.md")
	write(t, wt, "src/one/a.go", "package one\n")
	git(t, wt, "rm", "-q", "docs/README.md")
	write(t, wt, "src/two.go", "package two\n")
	write(t, wt, "README.md", "# good\n\nThe second task's.\n")

	res, err := Run(context.Background(), options(repo, "feat: [S-004] T-003 one"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.ExitCode() != 0 {
		t.Fatalf("stopped at %q: %s; sync %+v, check %+v", res.Stopped, res.Error, res.Sync, res.Check)
	}
	if res.Commit == nil || strings.Join(res.Commit.Paths, ",") != "docs/README.md,src/one/a.go" || committed(t, wt, res.Commit.Hash) != "docs/README.md,src/one/a.go" {
		t.Fatalf("commit: %+v", res.Commit)
	}
	if strings.Join(res.Left, ",") != "README.md,src/two.go" {
		t.Errorf("left: %v", res.Left)
	}
	if modified, untracked := git(t, wt, "diff", "--name-only"), git(t, wt, "ls-files", "--others", "--exclude-standard"); modified != "README.md" || untracked != "src/two.go" {
		t.Errorf("uncommitted: modified %q, untracked %q", modified, untracked)
	}
	if res.Sync == nil || res.Sync.Synced || res.Sync.Stopped != storygit.StopUncommitted || !strings.Contains(res.Sync.Continue, "the close that leaves none syncs story/S-004") {
		t.Errorf("sync: %+v", res.Sync)
	}
	if res.Move == nil || res.Move.State != workitem.Done || res.Log == nil || res.Log.Entry != "feat: [S-004] T-003 one" {
		t.Errorf("move %+v, log %+v", res.Move, res.Log)
	}
	if res.Touches == nil || len(res.Touches.Task) != 0 || strings.Join(res.Touches.Story, ",") != "docs/README.md,src/one/a.go" {
		t.Errorf("touches: %+v", res.Touches)
	}
	if got := get(t, repo, "T-003").Touches; strings.Join(got, ",") != "src/one,docs" {
		t.Errorf("T-003 touches %v", got)
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"left":["README.md","src/two.go"]`) {
		t.Errorf("json: %s", data)
	}

	o := options(repo, "feat: [S-004] second")
	o.Task = second
	res, err = Run(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Commit == nil || strings.Join(res.Commit.Paths, ",") != "README.md,src/two.go" || len(res.Left) != 0 {
		t.Fatalf("second: stopped at %q: %s, commit %+v, left %v", res.Stopped, res.Error, res.Commit, res.Left)
	}
	if res.Sync == nil || !res.Sync.Synced {
		t.Errorf("second: sync %+v", res.Sync)
	}
	if dirty := git(t, wt, "status", "--porcelain"); dirty != "" {
		t.Errorf("the worktree is not clean: %s", dirty)
	}
	if got := get(t, repo, second).Touches; strings.Join(got, ",") != "src/two.go,README.md" {
		t.Errorf("%s touches %v", second, got)
	}
}

// I-0108, ADR-0128: three open tasks of one layer, each with its file in the
// worktree, close in turn, each in a commit of its own file alone under its
// own message, and no task's touches widen with another's file.
func TestRunClosesALayerOfThreeTasksInACommitEach(t *testing.T) {
	repo, wt := project(t)
	setTouches(t, repo, "T-003", "src/a.go")
	ids := []string{"T-003", newTask(t, repo, "B", "src/b.go"), newTask(t, repo, "C", "src/c.go")}
	files := []string{"src/a.go", "src/b.go", "src/c.go"}
	for _, f := range files {
		write(t, wt, f, "package src\n")
	}
	base := git(t, wt, "rev-parse", "HEAD")

	var subjects []string
	for i, id := range ids {
		o := options(repo, "feat: [S-004] "+id+" its file")
		o.Task = id
		res, err := Run(context.Background(), o)
		if err != nil {
			t.Fatal(err)
		}
		if res.Stopped != "" {
			t.Fatalf("%s: stopped at %q: %s", id, res.Stopped, res.Error)
		}
		if res.Commit == nil || strings.Join(res.Commit.Paths, ",") != files[i] || committed(t, wt, res.Commit.Hash) != files[i] {
			t.Errorf("%s: commit %+v", id, res.Commit)
		}
		if strings.Join(res.Left, ",") != strings.Join(files[i+1:], ",") {
			t.Errorf("%s: left %v", id, res.Left)
		}
		if got := get(t, repo, id).Touches; strings.Join(got, ",") != files[i] {
			t.Errorf("%s touches %v", id, got)
		}
		subjects = append([]string{o.Message}, subjects...)
	}
	if log := git(t, wt, "log", "--format=%s", base+"..HEAD"); log != strings.Join(subjects, "\n") {
		t.Errorf("the branch's commits:\n%s", log)
	}
}

// ADR-0128: a changed file no task declares goes with the task closed first,
// whichever it is, and widens that task's touches.
func TestRunCommitsAnUndeclaredFileWithTheTaskClosedFirst(t *testing.T) {
	repo, wt := project(t)
	setTouches(t, repo, "T-003", "src/a.go")
	second := newTask(t, repo, "B", "src/b.go")
	for _, f := range []string{"src/a.go", "src/b.go", "src/extra.go"} {
		write(t, wt, f, "package src\n")
	}

	o := options(repo, "feat: [S-004] B")
	o.Task = second
	res, err := Run(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Commit == nil || strings.Join(res.Commit.Paths, ",") != "src/b.go,src/extra.go" || strings.Join(res.Left, ",") != "src/a.go" {
		t.Fatalf("stopped at %q: %s, commit %+v, left %v", res.Stopped, res.Error, res.Commit, res.Left)
	}
	if got := get(t, repo, second).Touches; strings.Join(got, ",") != "src/b.go,src/extra.go" {
		t.Errorf("%s touches %v", second, got)
	}
	if got := get(t, repo, "T-003").Touches; strings.Join(got, ",") != "src/a.go" {
		t.Errorf("T-003 touches %v", got)
	}
}

// ADR-0128: a close with nothing to commit needs no message: it runs to the
// end and logs that the task was closed.
func TestRunClosesWithNothingToCommitAndNoMessage(t *testing.T) {
	repo, _ := project(t)

	res, err := Run(context.Background(), options(repo, ""))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Commit != nil || res.Left == nil || len(res.Left) != 0 {
		t.Fatalf("stopped at %q: %s, commit %+v, left %v", res.Stopped, res.Error, res.Commit, res.Left)
	}
	if res.Move == nil || res.Move.State != workitem.Done || res.Log == nil || res.Log.Entry != "Closed T-003: T3" {
		t.Errorf("move %+v, log %+v", res.Move, res.Log)
	}
	if !strings.Contains(narrative(t, repo), "\nClosed T-003: T3\n") {
		t.Error("the close was not logged")
	}
}

// ADR-0128: a close with something to commit and no message stops at the
// commit step, asking for one, and leaves the worktree, its index, and the
// task as they were.
func TestRunStopsForAMessageWhenThereIsSomethingToCommit(t *testing.T) {
	repo, wt := project(t)
	write(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	head := git(t, wt, "rev-parse", "HEAD")

	res, err := Run(context.Background(), options(repo, " "))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != StepCommit || res.ExitCode() != 1 || !strings.Contains(res.Error, `a commit message is needed to commit README.md: flai task done T-003 -m "<message>"`) {
		t.Fatalf("stopped at %q: %s", res.Stopped, res.Error)
	}
	if res.Commit != nil || res.Sync != nil || res.Move != nil || res.Log != nil {
		t.Errorf("steps ran: %+v", res)
	}
	if git(t, wt, "rev-parse", "HEAD") != head || git(t, wt, "diff", "--name-only") != "README.md" || git(t, wt, "diff", "--cached", "--name-only") != "" {
		t.Error("the worktree changed")
	}
	if st := get(t, repo, "T-003").Status; st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
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

// S-0333: the commit is told to each other story in progress or in review
// whose claim covers a path it changed, a shared path included, on one
// conversation per pair, which a later commit's notice reuses, its about
// growing; a story whose claim is empty is told of every path, and one whose
// claim covers none of them is not told.
func TestRunTellsEachOpenStoryWhoseClaimCoversTheCommit(t *testing.T) {
	repo, wt := project(t)
	write(t, repo.Root, "system-flow.yaml", "version: 1\nname: good\nkey: g\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nclaims:\n  shared: [README.md]\n")
	// the repository's lint, so that a notice it rejects is not told
	write(t, repo.Root, ".markdownlint.yaml", "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n")
	story(t, repo, "S-006", workitem.InProgress, "README.md", "docs")
	story(t, repo, "S-007", workitem.Review)
	story(t, repo, "S-008", workitem.InProgress, "design")
	repo, err := workitem.Open(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	if shared := repo.Manifest.Claims.Shared; len(shared) != 1 || shared[0] != "README.md" {
		t.Fatalf("README.md is not a shared path: %v", shared)
	}
	write(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	write(t, wt, "src/main.go", "package main\n")

	res, err := Run(context.Background(), options(repo, "feat: [S-004] T-003 the readme\n\nWhy."))
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Commit == nil {
		t.Fatalf("stopped at %q: %s", res.Stopped, res.Error)
	}
	want := []Told{
		{Story: "S-006", Title: "Other S-006", Paths: []string{"README.md"}, Conversation: "MS-0001"},
		{Story: "S-007", Title: "Other S-007", Paths: []string{"README.md", "src/main.go"}, Conversation: "MS-0002"},
	}
	if !toldAs(res.Told, want) {
		t.Fatalf("told %+v, want %+v", res.Told, want)
	}
	c, err := messages.Get(repo, "MS-0001")
	if err != nil {
		t.Fatal(err)
	}
	hash := res.Commit.Hash[:7]
	text := "T-003 of S-004 changed paths S-006's claim covers.\n\nT-003, T3, committed " + hash + " on story/S-004, `feat: [S-004] T-003 the readme`, changing `README.md`. It reaches the main branch when S-004 is accepted; `git show " + hash + "` shows it until then. Reply here if it breaks your work, or adjust to it early."
	if es := c.Entries(); c.From != "S-004" || c.To != "S-006" || strings.Join(c.About, ",") != "README.md" || len(es) != 1 || es[0].Author != "agent-S-004" || es[0].Text != text {
		t.Errorf("MS-0001: %+v\n%s", c, c.Body)
	}
	if err := c.Validate(); err != nil {
		t.Error(err)
	}

	write(t, wt, "docs/README.md", "# Docs\n\nChanged later.\n")
	res, err = Run(context.Background(), options(repo, "docs: [S-004] T-003 the docs"))
	if err != nil {
		t.Fatal(err)
	}
	want = []Told{
		{Story: "S-006", Title: "Other S-006", Paths: []string{"docs/README.md"}, Conversation: "MS-0001"},
		{Story: "S-007", Title: "Other S-007", Paths: []string{"docs/README.md"}, Conversation: "MS-0002"},
	}
	if res.Stopped != "" || !toldAs(res.Told, want) {
		t.Fatalf("again: stopped %q %s, told %+v", res.Stopped, res.Error, res.Told)
	}
	if c, err = messages.Get(repo, "MS-0001"); err != nil {
		t.Fatal(err)
	}
	// the same agent at the same second joins the entry before it
	if strings.Join(c.About, ",") != "README.md,docs/README.md" || strings.Count(c.Body, "T-003 of S-004 changed paths S-006's claim covers.") != 2 || !strings.Contains(c.Body, "`docs: [S-004] T-003 the docs`, changing `docs/README.md`.") {
		t.Errorf("the pair's conversation is reused: %+v\n%s", c, c.Body)
	}
	if all, _ := messages.List(repo); len(all) != 2 {
		t.Errorf("one conversation per pair: %d", len(all))
	}
}

// S-0333: a notice that cannot be written is logged, is not in told, and does
// not stop the close.
func TestRunGoesOnWhenANoticeCannotBeWritten(t *testing.T) {
	repo, wt := project(t)
	story(t, repo, "S-006", workitem.InProgress)
	write(t, repo.Root, "wip/messages", "a file where the folder goes\n")
	write(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	var logs strings.Builder
	o := options(repo, "feat: [S-004] T-003 the readme")
	o.Logger = slog.New(slog.NewTextHandler(&logs, nil))

	res, err := Run(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != "" || res.Told == nil || len(res.Told) != 0 || get(t, repo, "T-003").Status != workitem.Done {
		t.Fatalf("stopped at %q: %s, told %+v", res.Stopped, res.Error, res.Told)
	}
	if !strings.Contains(logs.String(), "change notice not sent") || !strings.Contains(logs.String(), "item=S-006") {
		t.Errorf("logs:\n%s", logs.String())
	}
}

// toldAs reports whether told is want, in order.
func toldAs(told, want []Told) bool {
	return slices.EqualFunc(told, want, func(a, b Told) bool {
		return a.Story == b.Story && a.Title == b.Title && a.Conversation == b.Conversation && slices.Equal(a.Paths, b.Paths)
	})
}
