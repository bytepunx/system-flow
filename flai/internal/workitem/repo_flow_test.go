package workitem_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flowBoardRel is the workflow file's path in the test project, whose design
// folder is design.
const flowBoardRel = "design/workflows/board.md"

// flowBaseline is the baseline workflow file's text with each pair of old and
// new strings replaced once.
func flowBaseline(t *testing.T, edits ...string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for i := 0; i+1 < len(edits); i += 2 {
		if !strings.Contains(text, edits[i]) {
			t.Fatalf("the baseline lacks %q", edits[i])
		}
		text = strings.Replace(text, edits[i], edits[i+1], 1)
	}
	return text
}

// parkedBoard is the baseline with a state of its own, parked, that items
// reach from backlog and leave back to it.
func parkedBoard(t *testing.T) string {
	return flowBaseline(t,
		"  - {name: cancelled, kind: cancelled}\n", "  - {name: cancelled, kind: cancelled}\n  - {name: parked, kind: queue}\n",
		"  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog}\n  - {from: backlog, to: parked}\n  - {from: parked, to: backlog}\n")
}

// noPullBoard is the baseline with no transition giving the pull effect.
func noPullBoard(t *testing.T) string {
	return flowBaseline(t, "effects: [pull], ", "")
}

// flowProject is a project root with a minimal manifest and, unless board is
// empty, board as its workflow file.
func flowProject(t *testing.T, board string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if board != "" {
		writeBoard(t, root, board)
	}
	return root
}

// writeBoard writes board as the project's workflow file.
func writeBoard(t *testing.T, root, board string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(flowBoardRel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0755: Open runs the project's own workflow file, and the baseline when
// the project has none.
func TestOpenRunsTheProjectsWorkflowFile(t *testing.T) {
	repo, err := workitem.Open(flowProject(t, parkedBoard(t)))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(repo.Flow().Names(), "parked") || repo.Flow().Path != flowBoardRel {
		t.Errorf("Open runs %s at %s; want the project's file, with parked", repo.Flow().Names(), repo.Flow().Path)
	}
	repo, err = workitem.Open(flowProject(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	if got := repo.Flow().Path; got != workflow.BaselinePath {
		t.Errorf("a project with no file runs %s; want the baseline", got)
	}
}

// S-0755: a workflow file flai cannot run refuses Open, naming the finding;
// OpenUnchecked keeps the refusal and runs the baseline.
func TestOpenRefusesAWorkflowFileItCannotRun(t *testing.T) {
	root := flowProject(t, noPullBoard(t))
	repo, err := workitem.Open(root)
	var le *workflow.LoadError
	if repo != nil || !errors.As(err, &le) || le.Code != workflow.CodeNoPull {
		t.Fatalf("Open gave %v, %v; want a refusal naming %s", repo, err, workflow.CodeNoPull)
	}
	if !strings.Contains(err.Error(), root) || !strings.Contains(err.Error(), flowBoardRel+":12: "+workflow.CodeNoPull) {
		t.Errorf("the refusal reads %q; want the project and the finding's path, line, and code", err)
	}
	repo, err = workitem.OpenUnchecked(root)
	if err != nil {
		t.Fatal(err)
	}
	if !errors.As(repo.FlowErr(), &le) || le.Code != workflow.CodeNoPull {
		t.Errorf("OpenUnchecked keeps %v; want the %s refusal", repo.FlowErr(), workflow.CodeNoPull)
	}
	if got := repo.Flow().Path; got != workflow.BaselinePath {
		t.Errorf("OpenUnchecked runs %s while the file is refused; want the baseline", got)
	}
}

// S-0755: ReloadFlow reads the file again once it changes: an edit that
// adds a state, one that breaks it, which keeps the workflow run before, a
// fix, and a removal, which runs the baseline.
func TestReloadFlowFollowsTheFile(t *testing.T) {
	root := flowProject(t, "")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	copied := *repo
	writeBoard(t, root, parkedBoard(t))
	if err := repo.ReloadFlow(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(copied.Flow().Names(), "parked") {
		t.Errorf("after an edit adding parked the repository runs %s", copied.Flow().Names())
	}
	writeBoard(t, root, noPullBoard(t))
	var le *workflow.LoadError
	if err := repo.ReloadFlow(); !errors.As(err, &le) || le.Code != workflow.CodeNoPull {
		t.Errorf("after an edit taking the pull out ReloadFlow gives %v; want the %s refusal", err, workflow.CodeNoPull)
	}
	if !errors.As(repo.FlowErr(), &le) || !slices.Contains(repo.Flow().Names(), "parked") {
		t.Errorf("while refused: FlowErr %v, running %s; want the refusal, and the workflow run before", repo.FlowErr(), repo.Flow().Names())
	}
	if err := repo.ReloadFlow(); !errors.As(err, &le) {
		t.Errorf("an unchanged refused file reloads to %v; want the refusal still", err)
	}
	writeBoard(t, root, parkedBoard(t))
	if err := repo.ReloadFlow(); err != nil || repo.FlowErr() != nil {
		t.Errorf("after the fix: %v, %v; want no refusal", err, repo.FlowErr())
	}
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(flowBoardRel))); err != nil {
		t.Fatal(err)
	}
	if err := repo.ReloadFlow(); err != nil {
		t.Fatal(err)
	}
	if got := repo.Flow().Path; got != workflow.BaselinePath {
		t.Errorf("after the file's removal the repository runs %s; want the baseline", got)
	}
}
