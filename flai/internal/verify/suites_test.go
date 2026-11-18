package verify

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// baselineBoard is the baseline workflow file flai carries, as text a test
// edits into a project's own file.
func baselineBoard(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../workflow/baseline/board.md")
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// edited is the baseline file with each old of pairs replaced by the new
// after it, failing the test when the file lacks one.
func edited(t *testing.T, pairs ...string) string {
	t.Helper()
	board := baselineBoard(t)
	for i := 0; i < len(pairs); i += 2 {
		if !strings.Contains(board, pairs[i]) {
			t.Fatalf("the baseline board has no %q", pairs[i])
		}
		board = strings.Replace(board, pairs[i], pairs[i+1], 1)
	}
	return board
}

// parsedBoard parses a workflow file as a project's own.
func parsedBoard(t *testing.T, content string) workflow.Workflow {
	t.Helper()
	w, prob := workflow.Parse("design/workflows/board.md", content)
	if prob != nil {
		t.Fatalf("the workflow file does not parse: %s", prob)
	}
	return w
}

// The baseline's transitions, as edits find them.
const (
	toReview      = "    verify: {suites: all, severity: warn}"
	backToReady   = "{from: in-progress, to: ready, rules: [criteria-present, not-draft, after-done]}"
	taskDone      = "{from: in-progress, to: done, types: [task], rules: [children-closed]}"
	acceptVerify  = "effects: [accept], rules: [children-closed, criteria-ticked], verify: all}"
	reviewRuns    = ", runs: {kluster: [deploy, checks]}"
	verifyAtLine  = "verify_at: {batch: all, ci: all}\n"
	acceptNoSuite = "effects: [accept], rules: [children-closed, criteria-ticked]}"
)

// ADR-0312 §2: the suites a story needs are those the verify of each
// transition a story takes out of its state lists, each once in file order,
// a warn severity's as a refuse's; a task's transition does not count; all
// on one of them, or none giving verify, is every suite (nil); and the
// baseline's, as a file with none of the keys reads, is every suite out of
// in-progress and review, today's close-out (ADR-0312 §8).
func TestNeededSuitesAreTheVerifySuitesOfTheTransitionsOutOfTheState(t *testing.T) {
	listed := edited(t,
		toReview, "    verify: {suites: [static, behavior], severity: warn}",
		backToReady, strings.TrimSuffix(backToReady, "}")+", verify: [behavior, docs]}",
		taskDone, strings.TrimSuffix(taskDone, "}")+", verify: [slow]}",
		acceptVerify, strings.Replace(acceptVerify, "verify: all", "verify: [integration]", 1),
	)
	all := edited(t,
		toReview, "    verify: [static]",
		backToReady, strings.TrimSuffix(backToReady, "}")+", verify: all}",
	)
	none := edited(t, toReview+"\n", "", acceptVerify, acceptNoSuite, reviewRuns, "", verifyAtLine, "")
	onlyBatch := edited(t, toReview+"\n", "", acceptVerify, acceptNoSuite, reviewRuns, "", verifyAtLine, "verify_at: {batch: [static]}\n")
	for _, tc := range []struct {
		name, file, state string
		want              []string
	}{
		{"different suites out of in-progress, once each", listed, "in-progress", []string{"static", "behavior", "docs"}},
		{"the acceptance's out of review", listed, "review", []string{"integration"}},
		{"no transition out of ready gives verify", listed, "ready", nil},
		{"all on one transition out of in-progress", all, "in-progress", nil},
		{"the baseline, out of in-progress", baselineBoard(t), "in-progress", nil},
		{"the baseline, out of review", baselineBoard(t), "review", nil},
		{"a file with none of the keys reads as the baseline", none, "in-progress", nil},
		{"a file that gives no transition verify", onlyBatch, "in-progress", nil},
	} {
		w := parsedBoard(t, tc.file)
		if got := NeededSuites(w, tc.state); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: suites %#v, want %#v", tc.name, got, tc.want)
		}
	}
	if parsedBoard(t, none).UsesSuiteKeys() || !parsedBoard(t, onlyBatch).UsesSuiteKeys() {
		t.Error("the edits should leave none of the keys, and then verify_at alone")
	}
}

// StorySuites reads the project's own workflow file at the story's state,
// and refuses an ID that names no story.
func TestStorySuitesReadTheProjectsWorkflowAtTheStorysState(t *testing.T) {
	f := newVerifyFixture(t)
	if got, err := StorySuites(f.opts.Project, "S-4"); err != nil || got != nil {
		t.Errorf("with the baseline: %#v, %v; want every suite", got, err)
	}
	path := filepath.Join(f.root, "design", "workflows", "board.md")
	if err := os.WriteFile(path, []byte(edited(t, toReview, "    verify: [behavior]")), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := StorySuites(repo, "S-4"); err != nil || !reflect.DeepEqual(got, []string{"behavior"}) {
		t.Errorf("with the project's file: %#v, %v; want behavior", got, err)
	}
	if _, err := StorySuites(repo, "E-1"); err == nil || !strings.Contains(err.Error(), "name a story (S-nnnn)") {
		t.Errorf("an epic: %v", err)
	}
	if _, err := StorySuites(repo, "S-999"); err == nil || !strings.Contains(err.Error(), "name a story of this project") {
		t.Errorf("no such story: %v", err)
	}
}
