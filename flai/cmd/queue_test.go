package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// queueProject is orchestratedStoryInReview with its S-0001 verified at the
// head of its branch, and S-0002 in review an hour later, with criteria 2 and
// 3 unticked, no flai verify result, and a thread open on its task T-0002.
func queueProject(t *testing.T) string {
	t.Helper()
	root, head := orchestratedStoryInReview(t, false)
	if err := verify.SaveReport(openProject(t, root), verify.Report{Story: "S-0001", Commit: head, Passed: true}); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Half done", "--epic", "E-0001"); code != 0 {
		t.Fatal(errOut)
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0002-half-done.md")
	_ = os.WriteFile(file, []byte(strings.Replace(read(t, file), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] one\n- [ ] two\n- [ ] three\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0002", "ready"}, {"move", "S-0002", "in-progress"}, {"stream", "open", "S-0002"},
		{"task", "new", "Do half", "--story", "S-0002"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"}, {"move", "T-0002", "done"},
		{"thread", "new", "--on", "T-0002", "--by", "alex", "Is half enough?", "Say so."},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	if _, errOut, code := runInAt(t, root, time.Date(2026, 9, 15, 22, 0, 0, 0, time.UTC), "move", "S-0002", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root
}

// writeQueueState writes the merge queue's state file of the project at root.
func writeQueueState(t *testing.T, root string, st mergequeue.State) {
	t.Helper()
	path := mergequeue.StatePath(filepath.Join(root, ".flai-cache"))
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0398, ADR-0173: flai queue lists the eligible stories first, in queue
// order, then each other story in review with every reason it is not
// eligible, each by kind, then the batch and its step and the outcomes.
func TestQueueListsTheEligibleThenEveryReasonOfTheOthers(t *testing.T) {
	root := queueProject(t)
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	writeQueueState(t, root, mergequeue.State{
		Enqueued: []mergequeue.Enqueued{{Story: "S-0001", By: "alex", At: at}},
		Batch:    &mergequeue.Batch{Stories: []string{"S-0001"}, Step: "verify"},
		Outcomes: []mergequeue.Outcome{{Story: "S-0009", Outcome: "landed", At: at, Detail: "merged at abc123"}},
	})

	out, errOut, code := runIn(t, root, "queue")
	if code != 0 {
		t.Fatalf("flai queue: %d %s", code, errOut)
	}
	for _, want := range []string{
		"eligible (1)\n  1. S-0001 Ship it",
		"in review since 2026-09-15T21:00:00Z enqueued\nnot eligible (1)\n  S-0002 Half done",
		"in review since 2026-09-15T22:00:00Z\n",
		"       verify_missing: S-0002 has no flai verify result",
		"       criteria_unticked: acceptance criteria 2, 3 are unticked on S-0002",
		"       thread_open: thread TH-0001 is not resolved on S-0002 or its tasks",
		"batch: S-0001, step verify\n",
		"outcomes\n  S-0009 landed 2026-09-15T20:00:00Z: merged at abc123\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("flai queue says %q:\n%s", want, out)
		}
	}

	js, errOut, code := runIn(t, root, "queue", "--json")
	var q mergequeue.Queue
	if code != 0 || json.Unmarshal([]byte(js), &q) != nil {
		t.Fatalf("flai queue --json: %d %s %s", code, js, errOut)
	}
	if q.Mode != "local" || q.Off || len(q.Eligible) != 1 || q.Eligible[0].Story != "S-0001" || !q.Eligible[0].Enqueued {
		t.Errorf("eligible: %+v", q)
	}
	if len(q.Ineligible) != 1 || q.Ineligible[0].Story != "S-0002" {
		t.Fatalf("ineligible: %+v", q.Ineligible)
	}
	var kinds []string
	for _, r := range q.Ineligible[0].Reasons {
		kinds = append(kinds, r.Kind)
	}
	if got := strings.Join(kinds, " "); got != "verify_missing criteria_unticked thread_open" {
		t.Errorf("reason kinds: %s", got)
	}
	if q.Batch == nil || strings.Join(q.Batch.Stories, " ") != "S-0001" || q.Batch.Step != "verify" || len(q.Outcomes) != 1 {
		t.Errorf("batch and outcomes: %+v %+v", q.Batch, q.Outcomes)
	}
}

// S-0398: with no state file the queue shows no batch and no outcomes.
func TestQueueWithNoStateShowsNoBatch(t *testing.T) {
	root := queueProject(t)
	out, errOut, code := runIn(t, root, "queue")
	if code != 0 || !strings.Contains(out, "batch: none\n") || strings.Contains(out, "outcomes") {
		t.Errorf("flai queue: %d %s\n%s", code, errOut, out)
	}
	js, _, _ := runIn(t, root, "queue", "--json")
	if !strings.Contains(js, `"outcomes": []`) || strings.Contains(js, `"batch"`) {
		t.Errorf("flai queue --json: %s", js)
	}
}

// S-0398: with no story in review flai queue says so plainly.
func TestQueueWithNoStoryInReview(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	out, errOut, code := runIn(t, tempProject(t), "queue")
	if code != 0 || out != "no story is in review\nbatch: none\n" {
		t.Errorf("flai queue: %d %q %s", code, out, errOut)
	}
}

// S-0398, ADR-0169: under pull-request acceptance flai queue says the queue
// does not run, and why, and lists no story.
func TestQueueIsOffUnderPullRequestAcceptance(t *testing.T) {
	root := queueProject(t)
	m := filepath.Join(root, "system-flow.yaml")
	_ = os.WriteFile(m, []byte(read(t, m)+"acceptance:\n  mode: pull-request\n"), 0o644)

	out, errOut, code := runIn(t, root, "queue")
	if code != 0 || out != "merge queue off under acceptance.mode pull-request: "+mergequeue.OffUnderPullRequest+"\n" {
		t.Errorf("flai queue: %d %q %s", code, out, errOut)
	}
	js, _, _ := runIn(t, root, "queue", "--json")
	var q mergequeue.Queue
	if err := json.Unmarshal([]byte(js), &q); err != nil {
		t.Fatalf("%v: %s", err, js)
	}
	if !q.Off || q.Mode != "pull-request" || len(q.Eligible) != 0 || len(q.Ineligible) != 0 {
		t.Errorf("flai queue --json: %+v", q)
	}
}
