package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A story begun on another host reaches this clone through git without its
// branch or worktree. flai stream open takes up its narrative and checks the
// branch out: fetched from the remote when it is there, from the main
// branch otherwise (ADR-0064).
func TestAStreamBegunOnAnotherHostIsReopened(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "agent-S-0001")
	t.Setenv("FLAI_SESSION", "there")

	// the other host: the project, pushed to a bare remote, with S-0001 and
	// S-0002 begun there; only S-0001's branch is pushed
	remote := filepath.Join(t.TempDir(), "remote.git")
	gitIn(t, filepath.Dir(remote), "init", "-q", "--bare", "-b", "main", remote)
	there := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "docs"} {
		_ = os.MkdirAll(filepath.Join(there, d), 0o755)
		_ = os.WriteFile(filepath.Join(there, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(there, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, there, "init", "-q", "-b", "main")
	gitIn(t, there, "config", "user.email", "t@t")
	gitIn(t, there, "config", "user.name", "t")
	gitIn(t, there, "remote", "add", "origin", remote)
	gitIn(t, there, "add", "-A")
	gitIn(t, there, "commit", "-q", "-m", "init")
	if _, errOut, code := runIn(t, there, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	for _, title := range []string{"Pushed", "Unpushed"} {
		if _, errOut, code := runIn(t, there, "story", "new", title, "--epic", "E-0001"); code != 0 {
			t.Fatal(errOut)
		}
	}
	for _, id := range []string{"S-0001", "S-0002"} {
		if _, errOut, code := runIn(t, there, "stream", "open", id); code != 0 {
			t.Fatal(errOut)
		}
	}
	wtThere := filepath.Join(there, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wtThere, "docs", "begun.md"), []byte("begun there\n"), 0o644)
	gitIn(t, wtThere, "add", "-A")
	gitIn(t, wtThere, "commit", "-q", "-m", "feat: [S-0001] begun there")
	gitIn(t, there, "add", "-A")
	gitIn(t, there, "commit", "-q", "-m", "chore: the work items")
	gitIn(t, there, "push", "-q", "origin", "main", "story/S-0001")

	// this host: a clone with the narratives and neither branch
	here := filepath.Join(t.TempDir(), "here")
	gitIn(t, filepath.Dir(here), "clone", "-q", remote, here)
	gitIn(t, here, "config", "user.email", "t@t")
	gitIn(t, here, "config", "user.name", "t")
	t.Setenv("FLAI_AGENT", "agent-S-0001-here")
	t.Setenv("FLAI_SESSION", "here")

	out, errOut, code := runIn(t, here, "stream", "open", "S-0001")
	if code != 0 || !strings.Contains(out, "reopened wip/agents/S-0001.md") || !strings.Contains(out, "branch story/S-0001 (fetched from origin) checked out at .flai-cache/worktrees/S-0001") {
		t.Fatalf("reopen with the branch on the remote: %d %s %s", code, out, errOut)
	}
	wt := filepath.Join(here, ".flai-cache", "worktrees", "S-0001")
	if got, err := os.ReadFile(filepath.Join(wt, "docs", "begun.md")); err != nil || string(got) != "begun there\n" {
		t.Errorf("the worktree should hold what was committed there: %q %v", got, err)
	}
	n, err := workitem.ReadNarrative(filepath.Join(here, "wip", "agents", "S-0001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if n.Agent != "agent-S-0001-here" || n.Session != "here" || n.Host != workitem.ThisHost() {
		t.Errorf("the narrative should record who reopened it, and where: %+v", n)
	}
	if !strings.Contains(n.Body, "## Log") || !strings.Contains(n.Body, "Stream opened.") {
		t.Errorf("the narrative's body should be kept: %q", n.Body)
	}

	out, errOut, code = runIn(t, here, "stream", "open", "S-0002", "--json")
	if code != 0 || !strings.Contains(out, `"from": "main"`) || !strings.Contains(out, `"reopened": true`) {
		t.Fatalf("reopen with no branch anywhere: %d %s %s", code, out, errOut)
	}
	if gitIn(t, filepath.Join(here, ".flai-cache", "worktrees", "S-0002"), "rev-parse", "--abbrev-ref", "HEAD") != "story/S-0002" {
		t.Error("S-0002's worktree should be on a new story branch")
	}

	_, errOut, code = runIn(t, here, "stream", "open", "S-0001")
	if code == 0 || !strings.Contains(errOut, "stream S-0001 already exists") {
		t.Errorf("a stream with its worktree is not opened again: %d %s", code, errOut)
	}
	_, errOut, code = runIn(t, there, "stream", "open", "S-0002", "--no-branch")
	if code == 0 || !strings.Contains(errOut, "already exists") {
		t.Errorf("--no-branch has nothing to reopen: %d %s", code, errOut)
	}
}
