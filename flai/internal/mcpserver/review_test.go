package mcpserver

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// statusRunner answers git status in a worktree with what it holds.
type statusRunner struct {
	execx.System
	status string
}

func (r *statusRunner) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && len(args) > 0 && args[0] == "status" {
		return r.status, nil
	}
	return "", nil
}

// S-0140: item_move refuses a story whose worktree holds uncommitted work,
// naming it, and moves it once the worktree is clean. The server hands its
// runner to the repo, so the rule asks git as flai move does.
func TestItemMoveRefusesReviewWithUncommittedWork(t *testing.T) {
	f := setup(t)
	git := &statusRunner{status: "?? docs/left.md\n M flai/cmd/move.go"}
	New(Options{Repo: f.repo, Runner: git})
	if f.repo.Git != execx.Runner(git) {
		t.Fatal("the server must give the repo its runner")
	}
	if err := os.MkdirAll(f.repo.WorktreePath(f.story.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{"ready", "in-progress", "done"} {
		if _, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": to}); failed != "" {
			t.Fatal(failed)
		}
	}
	_, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"})
	if !strings.Contains(failed, "has uncommitted changes (docs/left.md, flai/cmd/move.go)") || !strings.Contains(failed, "commit them on story/"+f.story.ID) {
		t.Fatalf("review with uncommitted work: %q", failed)
	}
	if it, _ := f.repo.Get(f.story.ID); it.Status != workitem.InProgress {
		t.Errorf("a refused move changes nothing: %s", it.Status)
	}
	git.status = ""
	if out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"}); failed != "" || out["status"] != "review" {
		t.Errorf("a clean worktree moves: %v %s", out, failed)
	}
}
