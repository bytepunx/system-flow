package statsread

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// project writes a project with a sub-project, a board whose in-progress
// limit is 3, a story, and a thread on it, and opens it.
func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	board := "---\ntitle: Board\nstatus: active\nwip_limits:\n  in-progress: 3\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "One", Owner: "t", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := threads.New(repo, threads.NewOptions{Title: "Which way", On: "S-0001", Author: "agent", Text: "Left or right?", Now: now}); err != nil {
		t.Fatal(err)
	}
	return repo
}

// warnings are the events log took, decoded.
func warnings(t *testing.T, b *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(b.String()), "\n") {
		if l == "" {
			continue
		}
		var ev map[string]any
		if err := json.Unmarshal([]byte(l), &ev); err != nil {
			t.Fatal(err)
		}
		out = append(out, ev)
	}
	return out
}

func TestReadWithoutGitLeavesTheCommitsOutAndWarns(t *testing.T) {
	repo := project(t)
	var logged bytes.Buffer
	items, opt, err := Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || len(opt.Threads) != 1 || opt.WIPLimit != 3 || len(opt.Projects) != 1 || opt.Projects[0].Name != "cli" {
		t.Errorf("read %d items, options %+v", len(items), opt)
	}
	if opt.Commits != nil {
		t.Errorf("commits outside a repository: %v", opt.Commits)
	}
	evs := warnings(t, &logged)
	if len(evs) != 1 || evs[0]["level"] != "WARN" || evs[0]["component"] != "stats" || evs[0]["err"] == nil || evs[0]["detail"] == nil {
		t.Errorf("want one warning with err and detail, got %v", evs)
	}
}

func TestReadInARepositoryReadsTheCommits(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	repo := project(t)
	r := execx.System{}
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"},
		{"add", "-A"}, {"commit", "-q", "-m", "feat: [S-0001] the story"},
	} {
		if out, err := r.Run(repo.Root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	var logged bytes.Buffer
	_, opt, err := Read(r, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if got := opt.Commits["S-0001"]; len(got) != 1 || got[0] != "system-flow.yaml" {
		t.Errorf("commits = %v, want the manifest alone, wip left out", opt.Commits)
	}
	if logged.Len() != 0 {
		t.Errorf("warned in a repository: %s", logged.String())
	}
}

func TestReadStopsOnAnUnreadableThread(t *testing.T) {
	repo := project(t)
	if err := os.WriteFile(filepath.Join(threads.Dir(repo), "TH-0009-bad.md"), []byte("---\nid: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err == nil || !strings.Contains(err.Error(), "cannot read the threads") {
		t.Errorf("an unreadable thread: %v", err)
	}
}
