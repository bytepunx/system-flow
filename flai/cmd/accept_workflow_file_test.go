package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0301 §7, I-0380: a story whose branch changes the workflow file is
// accepted by its operator alone. The orchestrator's acceptance is refused
// before anything is merged, its dry run naming the blocker workflow_file, and
// so is an agent's on its own name; the operator's passes.
func TestOnlyTheOperatorAcceptsAStoryThatChangesTheWorkflowFile(t *testing.T) {
	t.Parallel()
	root, head := workflowFileStoryInReview(t)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	const wf = "design/workflows/board.md"

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n"))
	if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet") || !strings.Contains(errOut, "story/S-0001 changes the workflow file "+wf+", which is code the operator runs, so only the operator (alex) accepts S-0001, never the orchestrator or the merge queue, whatever the orchestrator judges (ADR-0301 §7)") {
		t.Errorf("the orchestrator's acceptance must be refused naming the workflow file: %d\n%s", code, errOut)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(res.OrchestratorBlockers) != 1 || res.OrchestratorBlockers[0].Code != preview.BlockWorkflowFile || !strings.Contains(res.OrchestratorBlockers[0].Message, "ADR-0301 §7") {
		t.Errorf("the workflow file must be the orchestrator's one blocker, naming ADR-0301 §7: %+v", res.OrchestratorBlockers)
	}
	if res.WorkflowFile != wf || len(res.Protected) != 0 {
		t.Errorf("the preview must name the workflow file and no protected path: %s", out)
	}

	want := "only the operator (alex) accepts S-0001, not agent-S-0001: its branch changes the workflow file " + wf + ", which is code the operator runs; the operator reads it in the diff and accepts it, never the orchestrator or the merge queue (ADR-0301 §7)"
	for _, args := range [][]string{
		{"accept", "S-0001", "--by", "agent-S-0001"},
		{"move", "S-0001", "done", "--by", "agent-S-0001"},
	} {
		_, errOut, code := runIn(t, root, args...)
		if code == 0 || !strings.Contains(errOut, want) || strings.Contains(errOut, "Claude Code protects") {
			t.Errorf("flai %v by an agent on its own name must be refused naming the workflow file: %d\n%s", args, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the operator's dry run: %d %s", code, errOut)
	}
	res = preview.Acceptance{}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.OperatorOnly != "only the operator (alex) accepts S-0001: its branch changes the workflow file "+wf+" (ADR-0301 §7)" || res.WorkflowFile != wf || len(res.Blockers) != 0 {
		t.Errorf("the operator's dry run must say only the operator accepts it, unblocked: %s", out)
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("the operator's acceptance: %d %s\n%s", code, errOut, out)
	}
	res = preview.Acceptance{}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Status != workitem.Done || !res.Merged || res.WorkflowFile != wf {
		t.Errorf("the operator accepts it, the workflow file still named: %+v", res)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:"+wf), "states:") {
		t.Error("the story branch is merged")
	}
}

// ADR-0301 §7: the merge queue never lands a story whose branch changes the
// workflow file, so while the queue batches and is not held the operator's
// acceptance lands it by hand rather than enqueue it where the queue would
// refuse it, and leaves the queue as it was.
func TestTheOperatorLandsAWorkflowFileStoryByHandWhileTheQueueBatches(t *testing.T) {
	t.Parallel()
	root, head := workflowFileStoryInReview(t)
	mergeQueueOn(t, root)
	queueUnderWay(t, root, false)
	// a hand landing while the queue batches waits for a verify at the head
	// (ADR-0221)
	storeS0001PassingVerify(t, root, head, gitIn(t, root, "rev-parse", "main"))

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run")
	if code != 0 || !strings.Contains(out, "only the operator (alex) accepts S-0001: its branch changes the workflow file design/workflows/board.md (ADR-0301 §7)\n") {
		t.Errorf("the dry run must name the workflow file with no trailing list: %d %s\n%s", code, errOut, out)
	}
	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("the operator's acceptance while the queue batches: %d %s\n%s", code, errOut, out)
	}
	var res accepted
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Queue != nil || res.Status != workitem.Done || !res.Merged {
		t.Errorf("the operator lands it by hand, not through the queue: %+v %+v", res.Acceptance, res.Queue)
	}
	if st := readQueueState(t, root); len(st.Enqueued) != 1 || st.Enqueued[0].Story != "S-0005" {
		t.Errorf("enqueued %+v, want S-0005 only", st.Enqueued)
	}
}

// A story whose branch changes neither the workflow file nor a protected
// path has no operator-only sentence and no workflow_file blocker.
func TestAStoryThatLeavesTheWorkflowFileIsNotOperatorOnly(t *testing.T) {
	t.Parallel()
	root, head := orchestratedStoryInReview(t, true)
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("dry run: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.WorkflowFile != "" || res.OperatorOnly != "" || len(res.OrchestratorBlockers) != 0 {
		t.Errorf("a story that changes cli alone must have no operator-only sentence and no blocker: %s", out)
	}
}

// workflowFileStoryInReview is orchestratedStoryInReview whose story, owned
// by alex, also commits the baseline as design/workflows/board.md on its
// branch, under its touches; it returns the main checkout and the branch's
// head.
func workflowFileStoryInReview(t *testing.T) (root, head string) {
	t.Helper()
	root, _ = orchestratedStoryInReview(t, true)
	setOwner(t, root, "alex")
	if _, errOut, code := runIn(t, root, "touches", "S-0001", "--add", "design/workflows"); code != 0 {
		t.Fatal(errOut)
	}
	baseline, err := os.ReadFile("../internal/workflow/baseline/board.md")
	if err != nil {
		t.Fatal(err)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	if err := os.MkdirAll(filepath.Join(wt, "design", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, "design", "workflows", "board.md"), baseline, 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] change the board")
	return root, gitIn(t, root, "rev-parse", "story/S-0001")
}
