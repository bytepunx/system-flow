package mcpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// writeProjectBoard writes text as the workflow file of the project at root,
// as an operator's edit does while the server runs.
func writeProjectBoard(t *testing.T, root, text string) {
	t.Helper()
	path := filepath.Join(root, "design", "workflows", "board.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0755, ADR-0307 §3: the server reads the project's workflow file again at
// each tool call once it changed: a file written while it runs, adding
// testing, moves the story there; an edit adding staging moves it on; a
// broken file refuses the next call, naming the finding; and the fix serves
// calls again.
func TestEachToolCallFollowsTheWorkflowFile(t *testing.T) {
	t.Parallel()
	f := setup(t)
	if _, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "testing"}); failed == "" {
		t.Fatal("the baseline moved the story to testing, a state it lacks")
	}
	writeProjectBoard(t, f.repo.Root, mcpOtherWorkflow)
	if out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "testing"}); failed != "" || out["status"] != "testing" {
		t.Fatalf("after the file adds testing: %v %s", out, failed)
	}
	staged := strings.NewReplacer(
		"  - {name: testing, kind: active, limit: true}\n", "  - {name: testing, kind: active, limit: true}\n  - {name: staging, kind: active, limit: true}\n",
		"  - {from: testing, to: acceptance-review}\n", "  - {from: testing, to: staging}\n  - {from: staging, to: acceptance-review}\n",
		"[backlog, ready, in-progress, testing]", "[backlog, ready, in-progress, testing, staging]",
	).Replace(mcpOtherWorkflow)
	writeProjectBoard(t, f.repo.Root, staged)
	if out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "staging"}); failed != "" || out["status"] != "staging" {
		t.Fatalf("after an edit adding staging: %v %s", out, failed)
	}
	writeProjectBoard(t, f.repo.Root, strings.Replace(staged, "{from: ready, to: in-progress, effects: [pull]}", "{from: ready, to: in-progress}", 1))
	_, failed := f.call(t, "item_get", map[string]any{"id": f.story.ID})
	for _, want := range []string{"item_get is refused", f.repo.Root, "design/workflows/board.md:", workflow.CodeNoPull} {
		if !strings.Contains(failed, want) {
			t.Errorf("a call while the file is broken answers %q; want it to name %q", failed, want)
		}
	}
	writeProjectBoard(t, f.repo.Root, staged)
	if out, failed := f.call(t, "item_get", map[string]any{"id": f.story.ID}); failed != "" || out["status"] != "staging" {
		t.Errorf("after the fix: %v %s", out, failed)
	}
	if out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"}); failed != "" || out["status"] != "acceptance-review" {
		t.Errorf("after the fix the story moves on by the file: %v %s", out, failed)
	}
}
