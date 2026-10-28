package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestAcceptEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "design", "conventions"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "conventions", "README.md"), []byte("# c\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n    tags: [command]\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "cli"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "cli", "main.go"), []byte("package main\n"), 0o644)
	// flai keeps its cache, the event log among it, out of git (ADR-0211 §4)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, step := range [][]string{
		{"epic", "new", "E"}, {"story", "new", "Ship it", "--epic", "E-0001", "--tag", "command"}, {"task", "new", "T", "--story", "S-0001"},
	} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	sp := filepath.Join(root, "wip", "kanban", "stories", "S-0001-ship-it.md")
	b, _ := os.ReadFile(sp)
	_ = os.WriteFile(sp, []byte(strings.Replace(string(b), "- [ ]\n", "- [x] done\n", 1)), 0o644)
	for _, step := range [][]string{{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"}, {"move", "S-0001", "review"}} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	_ = os.WriteFile(filepath.Join(root, "cli", "feature.go"), []byte("package main\n"), 0o644)
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "feat: [S-0001] ship it"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	out, errOut, code := runIn(t, root, "release", "S-0001", "--dry-run")
	if code != 0 || !strings.Contains(out, "cli        0.0.0 -> 0.1.0  minor  delivered") {
		t.Fatalf("release dry run: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--trailer", "Co-Authored-By: t <t@t>")
	if code != 0 {
		t.Fatalf("accept: %s\n%s", errOut, out)
	}
	// the story, its task, and E-0001, whose only story it is (S-0200)
	if !strings.Contains(out, "accepted S-0001: done, 3 items archived, committed") || !strings.Contains(out, "E-0001 followed it from review to done, archived") {
		t.Errorf("summary: %s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "archive", "kanban", "stories", "S-0001-ship-it.md")); err != nil {
		t.Error("story not archived")
	}
	log, _ := r.Run(root, "git", "log", "-1", "--format=%B")
	if !strings.HasPrefix(log, "chore: [S-0001] accept and archive, with E-0001") || !strings.Contains(log, "Co-Authored-By: t <t@t>") {
		t.Errorf("commit message:\n%s", log)
	}
	st, _ := r.Run(root, "git", "status", "--porcelain")
	if strings.TrimSpace(st) != "" {
		t.Errorf("working tree dirty after accept:\n%s", st)
	}
	tag, _ := r.Run(root, "git", "tag", "-l", "cli/*")
	if strings.TrimSpace(tag) != "" {
		t.Errorf("acceptance alone creates no tag: %q", tag)
	}
	_, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "already done") {
		t.Errorf("second accept must refuse: %s", errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "archive", "kanban", "epics", "E-0001-e.md")); err != nil {
		t.Error("the epic is archived with its last story")
	}
	// an epic with no story is accepted straight from backlog: walked to
	// done, no release
	if _, errOut, code := runIn(t, root, "epic", "new", "Empty"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code = runIn(t, root, "accept", "E-0002", "--by", "alex")
	if code != 0 || !strings.Contains(out, "accepted E-0002: done") {
		t.Errorf("epic accept from backlog: %d %s %s", code, out, errOut)
	}
}

// I-0204: a story branch flai serve pushed tracks origin's copy. The
// acceptance rebase rewrites its commits, so git branch -d refused the branch
// as not merged to its upstream, after main had the work and the worktree was
// gone. The acceptance deletes it with -D and finishes.
func TestAcceptDeletesABranchWhoseUpstreamTheRebaseLeftBehind(t *testing.T) {
	root := pushedStoryInReview(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	if got := gitIn(t, root, "show", "main:design/x.md"); got != "the story's" {
		t.Errorf("main must carry the story's work: %q", got)
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("an accepted story's branch is removed")
	}
	if b, _ := os.ReadFile(filepath.Join(root, "wip/archive/kanban/stories/S-0001-pushed.md")); !strings.Contains(string(b), "status: done") {
		t.Error("the story is done and archived")
	}
}

// I-0204: an acceptance that stopped after its fast-forward, with the work
// on main, the worktree removed, and the branch left, finishes when run
// again: it fast-forwards nothing and deletes the branch.
func TestAcceptFinishesFromAFastForwardItAlreadyMade(t *testing.T) {
	root := pushedStoryInReview(t)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	gitIn(t, wt, "rebase", "-q", "main")
	gitIn(t, root, "merge", "-q", "--ff-only", "story/S-0001")
	gitIn(t, root, "worktree", "remove", wt)
	merged := gitIn(t, root, "rev-parse", "main")

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("accepting again must finish: %d %s", code, errOut)
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("the branch left behind is removed")
	}
	if got := gitIn(t, root, "rev-parse", "main~1"); got != merged {
		t.Errorf("main must be the merge it had, %s, and the acceptance commit: main~1 is %s", merged, got)
	}
}

// pushedStoryInReview makes a project with story S-0001 in review whose
// branch commits design/x.md and is pushed to a bare origin with its
// upstream set, as flai serve pushes it for a pull request, and main a
// commit past the branch's base, so the acceptance rebase rewrites it.
func pushedStoryInReview(t *testing.T) (root string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	origin := filepath.Join(t.TempDir(), "origin.git")
	gitIn(t, root, "init", "-q", "--bare", origin)
	gitIn(t, root, "remote", "add", "origin", origin)
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Pushed", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-pushed.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wt, "design", "x.md"), []byte("the story's\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] x")
	gitIn(t, wt, "push", "-q", "-u", "origin", "story/S-0001")
	_ = os.WriteFile(filepath.Join(root, "design", "y.md"), []byte("main's\n"), 0o644)
	gitIn(t, root, "add", "design/y.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: main's y")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root
}
