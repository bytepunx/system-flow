package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/taskdone"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// closeClock is the good fixture's check time, a day after S-004 started.
var closeClock = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

// closeFixture is the good fixture copied to a new folder: S-004 in progress
// with its narrative, and its task T-003 in progress.
func closeFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("../metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	return dir
}

// closeProject is the good fixture in a git repository whose main branch
// holds it, with S-004's branch checked out in its worktree, which it
// returns. Integration: it runs real git.
func closeProject(t *testing.T) (*workitem.Repo, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	t.Setenv("FLAI_SESSION", "s1")
	dir := closeFixture(t)
	writeIn(t, dir, ".gitignore", ".flai-cache/\n")
	gitIn(t, dir, "init", "-q", "-b", "main")
	gitIn(t, dir, "add", "-A")
	gitIn(t, dir, "commit", "-q", "-m", "project")
	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	wt := repo.WorktreePath("S-004")
	gitIn(t, dir, "worktree", "add", "-q", "-b", storygit.Branch("S-004"), wt, "main")
	return repo, wt
}

// closer is an MCP server on repo as S-004's agent, at closeClock.
func closer(t *testing.T, repo *workitem.Repo) *fixture {
	t.Helper()
	cs := connectServer(t, Options{Repo: repo, Runner: execx.System{}, Agent: "agent-S-004", Version: "test", Now: func() time.Time { return closeClock }})
	return &fixture{repo: repo, cs: cs}
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

// t003 is T-003's status.
func t003(t *testing.T, repo *workitem.Repo) string {
	t.Helper()
	it, err := repo.Get("T-003")
	if err != nil {
		t.Fatal(err)
	}
	return it.Status
}

// field is the value at a dotted path in a decoded answer, nil when absent.
func field(out map[string]any, path string) any {
	var v any = out
	for _, k := range strings.Split(path, ".") {
		m, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		v = m[k]
	}
	return v
}

// ADR-0107: task_done commits, tells the other open stories whose claim
// covers the commit (S-0333), syncs, moves the task to done, logs, widens the
// touches, checks, and answers the inbox in one call.
func TestTaskDoneClosesATaskInOneCall(t *testing.T) {
	repo, wt := closeProject(t)
	writeIn(t, repo.Root, "wip/kanban/stories/S-006-six.md", "---\nid: S-006\ntype: story\nnature: feature\ntitle: Six\nstatus: in-progress\nparent: E-001\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\ntouches: [README.md]\n---\n\n# S-006 Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n")
	writeIn(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")

	out, failed := closer(t, repo).call(t, "task_done", map[string]any{"task": "T-3", "message": "docs: [S-004] T-003 the readme\n\nWhy."})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["stopped"] != "" || out["task"] != "T-003" || out["story"] != "S-004" {
		t.Fatalf("stopped at %v: %v", out["stopped"], out["error"])
	}
	if field(out, "commit.subject") != "docs: [S-004] T-003 the readme" || field(out, "sync.synced") != true || field(out, "move.state") != workitem.Done {
		t.Errorf("commit %v, sync %v, move %v", out["commit"], out["sync"], out["move"])
	}
	if field(out, "log.entry") != "docs: [S-004] T-003 the readme" || field(out, "check.passed") != true || field(out, "inbox.agent") != "agent-S-004" {
		t.Errorf("log %v, check %v, inbox %v", out["log"], out["check"], field(out, "inbox.agent"))
	}
	if told, _ := out["told"].([]any); len(told) != 1 || field(told[0].(map[string]any), "story") != "S-006" || field(told[0].(map[string]any), "conversation") != "MS-0001" {
		t.Errorf("told: %v", out["told"])
	}
	if st := t003(t, repo); st != workitem.Done {
		t.Errorf("T-003 is %s", st)
	}
	if dirty := gitIn(t, wt, "status", "--porcelain"); dirty != "" {
		t.Errorf("the worktree is not clean: %s", dirty)
	}
}

// sibling creates a second in-progress task of S-004, Second, touching
// touches, sets T-003's touches to mine, and returns the new task's ID.
func sibling(t *testing.T, repo *workitem.Repo, mine, touches []string) string {
	t.Helper()
	t003, err := repo.Get("T-003")
	if err != nil {
		t.Fatal(err)
	}
	t003.Touches = mine
	if err := repo.Save(t003); err != nil {
		t.Fatal(err)
	}
	it, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Second", Parent: "S-004", Touches: touches, Now: closeClock})
	if err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.TransitionAll(it, to, "agent", "", closeClock, false); err != nil {
			t.Fatal(err)
		}
	}
	return it.ID
}

// I-0108, ADR-0128: with two open tasks' files in the worktree, task_done
// on one commits its own file alone, answers the other's in left, leaves it
// uncommitted, and goes on past the sync it refuses.
func TestTaskDoneLeavesAnotherOpenTasksPaths(t *testing.T) {
	repo, wt := closeProject(t)
	sibling(t, repo, []string{"README.md"}, []string{"src/two.go"})
	writeIn(t, wt, "README.md", "# good\n\nT-003's.\n")
	writeIn(t, wt, "src/two.go", "package two\n")

	out, failed := closer(t, repo).call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["stopped"] != "" {
		t.Fatalf("stopped at %v: %v", out["stopped"], out["error"])
	}
	if paths, _ := field(out, "commit.paths").([]any); len(paths) != 1 || paths[0] != "README.md" {
		t.Errorf("commit: %v", out["commit"])
	}
	if left, _ := out["left"].([]any); len(left) != 1 || left[0] != "src/two.go" {
		t.Errorf("left: %v", out["left"])
	}
	if field(out, "sync.synced") != false || field(out, "sync.stopped") != storygit.StopUncommitted || field(out, "move.state") != workitem.Done {
		t.Errorf("sync %v, move %v", out["sync"], out["move"])
	}
	if untracked := gitIn(t, wt, "ls-files", "--others", "--exclude-standard"); untracked != "src/two.go" {
		t.Errorf("untracked: %q", untracked)
	}
}

// ADR-0128: task_done with nothing to commit needs no message: it runs to
// the end, answers no commit and left empty, and logs that the task was
// closed.
func TestTaskDoneClosesWithNothingToCommitAndNoMessage(t *testing.T) {
	repo, _ := closeProject(t)

	out, failed := closer(t, repo).call(t, "task_done", map[string]any{"task": "T-003"})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["stopped"] != "" || out["commit"] != nil {
		t.Fatalf("stopped at %v: %v; commit %v", out["stopped"], out["error"], out["commit"])
	}
	if left, ok := out["left"].([]any); !ok || len(left) != 0 {
		t.Errorf("left: %v", out["left"])
	}
	if field(out, "log.entry") != "Closed T-003: T3" || field(out, "move.state") != workitem.Done {
		t.Errorf("log %v, move %v", out["log"], out["move"])
	}
}

// A sync that stops on a conflict is the answer, not a tool error: stopped
// names the sync, which lists the conflicting path, and the task is not
// moved.
func TestTaskDoneAnswersARefusedSync(t *testing.T) {
	repo, wt := closeProject(t)
	writeIn(t, repo.Root, "README.md", "# good\n\nChanged on main.\n")
	gitIn(t, repo.Root, "commit", "-q", "-am", "docs: main's readme")
	writeIn(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")

	out, failed := closer(t, repo).call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
	if failed != "" {
		t.Fatalf("a stop is a tool error: %s", failed)
	}
	if out["stopped"] != taskdone.StepSync || !strings.Contains(out["error"].(string), "git rebase --continue") {
		t.Fatalf("stopped at %v: %v", out["stopped"], out["error"])
	}
	if c, _ := field(out, "sync.conflicts").([]any); len(c) != 1 || c[0] != "README.md" {
		t.Errorf("sync: %v", out["sync"])
	}
	if out["move"] != nil || out["inbox"] != nil {
		t.Errorf("steps after the sync ran: move %v, inbox %v", out["move"], out["inbox"])
	}
	if st := t003(t, repo); st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
	}
}

// A check that fails in the story is the answer too, after the move: the
// findings are in it, and the task is done.
func TestTaskDoneAnswersAFailedCheck(t *testing.T) {
	repo, wt := closeProject(t)
	writeIn(t, wt, "docs/notes.md", "# Notes\n\nNo front matter.\n")

	out, failed := closer(t, repo).call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 notes"})
	if failed != "" {
		t.Fatalf("a stop is a tool error: %s", failed)
	}
	if out["stopped"] != taskdone.StepCheck || field(out, "check.passed") != false || out["inbox"] != nil {
		t.Fatalf("stopped at %v: %v; check %v", out["stopped"], out["error"], out["check"])
	}
	findings, _ := field(out, "check.findings").([]any)
	found := false
	for _, f := range findings {
		m, _ := f.(map[string]any)
		path, _ := m["path"].(string)
		found = found || m["rule"] == "doc.front-matter" && strings.HasSuffix(filepath.ToSlash(path), "docs/notes.md")
	}
	if !found {
		t.Errorf("findings: %v", findings)
	}
	if st := t003(t, repo); st != workitem.Done {
		t.Errorf("T-003 is %s", st)
	}
}

// task_done answers what flai task done --json prints for the same case,
// field for field: its output is taskdone.Result, which the command prints.
// Two copies of the project close the same task, one through the tool and
// one through taskdone.Run, and answer the same JSON but for the commit's
// hash and the project's folder.
func TestTaskDoneAnswersWhatTheCommandPrints(t *testing.T) {
	if reflect.TypeFor[TaskDoneOut]() != reflect.TypeFor[taskdone.Result]() {
		t.Fatal("TaskDoneOut is a type of its own, not taskdone.Result")
	}
	const message = "docs: [S-004] T-003 notes"

	viaTool, wt := closeProject(t)
	writeIn(t, wt, "docs/notes.md", "# Notes\n\nNo front matter.\n")
	out, failed := closer(t, viaTool).call(t, "task_done", map[string]any{"task": "T-003", "message": message})
	if failed != "" {
		t.Fatal(failed)
	}

	direct, wt := closeProject(t)
	writeIn(t, wt, "docs/notes.md", "# Notes\n\nNo front matter.\n")
	res, err := taskdone.Run(context.Background(), taskdone.Options{Repo: direct, Runner: execx.System{}, Task: "T-003", Message: message, Agent: "agent-S-004", Session: "s1", Now: func() time.Time { return closeClock }, Version: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Stopped != taskdone.StepCheck || res.Commit == nil {
		t.Fatalf("the case is not a failed check: stopped %q: %s", res.Stopped, res.Error)
	}

	got, want := canonical(t, out, viaTool.Root), canonical(t, res, direct.Root)
	if got != want {
		t.Errorf("task_done answered\n%s\nflai task done --json prints\n%s", got, want)
	}
}

// canonical is v as JSON with sorted keys, its commit's hash and root
// written out of it.
func canonical(t *testing.T, v any, root string) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if c, ok := m["commit"].(map[string]any); ok {
		c["hash"] = "<hash>"
	}
	if data, err = json.MarshalIndent(m, "", "  "); err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(data), root, "<root>")
}

// task_done refuses, as a tool error and before any step, what it cannot
// close: an item that is not a task, and a task whose story has no worktree,
// with a message or without one, which the schema does not require
// (ADR-0128). It runs no git.
func TestTaskDoneRefusesWhatItCannotStart(t *testing.T) {
	repo, err := workitem.Open(closeFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	f := closer(t, repo)
	for _, tc := range []struct {
		name string
		args map[string]any
		want string
	}{
		{"a story", map[string]any{"task": "S-004", "message": "m"}, "S-004 is a story"},
		{"no worktree", map[string]any{"task": "T-003", "message": "m"}, "S-004 has no worktree"},
		{"no worktree, no message", map[string]any{"task": "T-003"}, "S-004 has no worktree"},
		{"no worktree, a blank message", map[string]any{"task": "T-003", "message": " "}, "S-004 has no worktree"},
	} {
		if out, failed := f.call(t, "task_done", tc.args); !strings.Contains(failed, tc.want) {
			t.Errorf("%s: answered %v, failed %q, want %q", tc.name, out, failed, tc.want)
		}
	}
	if st := t003(t, repo); st != workitem.InProgress {
		t.Errorf("T-003 is %s", st)
	}
}

// The tool is advertised with task alone required (ADR-0128), and the server's
// instructions send the agent to it at every task transition.
func TestTaskDoneIsAdvertisedAndInstructed(t *testing.T) {
	f := setup(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var tool *mcp.Tool
	for _, tl := range res.Tools {
		if tl.Name == "task_done" {
			tool = tl
		}
	}
	if tool == nil {
		t.Fatal("task_done is not advertised")
	}
	schema, err := json.Marshal(tool.InputSchema)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(schema), `"required":["task"]`) {
		t.Errorf("schema: %s", schema)
	}
	if in := f.cs.InitializeResult().Instructions; !strings.Contains(in, "At every task transition, close the task with task_done") {
		t.Errorf("the instructions do not send the agent to task_done: %s", in)
	}
}
