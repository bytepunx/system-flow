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

// readQueueState reads the merge queue's state file of the project at root.
func readQueueState(t *testing.T, root string) mergequeue.State {
	t.Helper()
	st, err := mergequeue.ReadState(filepath.Join(root, ".flai-cache"))
	if err != nil {
		t.Fatal(err)
	}
	return st
}

// S-0399 criterion 1: flai queue add enqueues an eligible story in review
// with who enqueued it, FLAI_AGENT, and when; a second add keeps its place
// and changes nothing.
func TestQueueAddRecordsWhoAndWhen(t *testing.T) {
	root := queueProject(t)
	t.Setenv("FLAI_AGENT", "orchestrator")
	out, errOut, code := runIn(t, root, "queue", "add", "S-0001")
	if code != 0 || out != "S-0001 enqueued by orchestrator at 2026-09-15T21:00:00Z, 1 in the queue\n" {
		t.Fatalf("flai queue add: %d %q %s", code, out, errOut)
	}
	at := time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)
	if st := readQueueState(t, root); len(st.Enqueued) != 1 || st.Enqueued[0] != (mergequeue.Enqueued{Story: "S-0001", By: "orchestrator", At: at}) {
		t.Errorf("enqueued %+v", st.Enqueued)
	}

	t.Setenv("FLAI_AGENT", "alex")
	js, errOut, code := runInAt(t, root, at.Add(time.Hour), "queue", "add", "s-1", "--json")
	var ch mergequeue.Change
	if code != 0 || json.Unmarshal([]byte(js), &ch) != nil {
		t.Fatalf("flai queue add --json: %d %s %s", code, js, errOut)
	}
	if ch.Action != "add" || ch.Story != "S-0001" || ch.Changed || ch.By != "alex" || len(ch.State.Enqueued) != 1 || ch.State.Enqueued[0].By != "orchestrator" {
		t.Errorf("a second add: %+v", ch)
	}
	out, _, _ = runIn(t, root, "queue", "add", "S-0001")
	if out != "S-0001 is enqueued already, by orchestrator at 2026-09-15T21:00:00Z; nothing changed\n" {
		t.Errorf("a second add says %q", out)
	}
	if out, _, _ := runIn(t, root, "queue"); !strings.Contains(out, "1. S-0001") || !strings.Contains(out, " enqueued\n") {
		t.Errorf("flai queue after add:\n%s", out)
	}
}

// S-0399 criterion 1: flai queue add refuses a story that is not eligible,
// naming each reason by kind, and writes nothing.
func TestQueueAddRefusesAnIneligibleStoryNamingEachReason(t *testing.T) {
	root := queueProject(t)
	_, errOut, code := runIn(t, root, "queue", "add", "S-0002")
	if code == 0 {
		t.Fatal("flai queue add enqueued an ineligible story")
	}
	for _, want := range []string{"S-0002 is not enqueued: it is not eligible for the merge queue", "verify_missing: S-0002 has no flai verify result", "criteria_unticked: acceptance criteria 2, 3 are unticked", "thread_open: thread TH-0001 is not resolved"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal does not say %q:\n%s", want, errOut)
		}
	}
	if _, err := os.Stat(mergequeue.StatePath(filepath.Join(root, ".flai-cache"))); err == nil {
		t.Error("a refused add wrote the state")
	}
}

// S-0399 criterion 1: flai queue add refuses a story not in review, and an
// item that is not a story.
func TestQueueAddRefusesAStoryNotInReview(t *testing.T) {
	root := queueProject(t)
	if _, errOut, code := runIn(t, root, "story", "new", "Later", "--epic", "E-0001"); code != 0 {
		t.Fatal(errOut)
	}
	for id, want := range map[string]string{
		"S-0003": "S-0003 is not enqueued: the merge queue takes a story in review, and S-0003 is backlog",
		"T-0001": "T-0001 is not enqueued: the merge queue takes stories only",
	} {
		if _, errOut, code := runIn(t, root, "queue", "add", id); code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("flai queue add %s: %d, want %q:\n%s", id, code, want, errOut)
		}
	}
}

// S-0399 criterion 1, ADR-0173: under pull-request acceptance flai queue add
// refuses, saying the queue runs under local acceptance only.
func TestQueueAddRefusesUnderPullRequestAcceptance(t *testing.T) {
	root := queueProject(t)
	m := filepath.Join(root, "system-flow.yaml")
	_ = os.WriteFile(m, []byte(read(t, m)+"acceptance:\n  mode: pull-request\n"), 0o644)
	_, errOut, code := runIn(t, root, "queue", "add", "S-0001")
	if code == 0 || !strings.Contains(errOut, "S-0001 is not enqueued: the merge queue runs under local acceptance only") {
		t.Errorf("flai queue add: %d %s", code, errOut)
	}
}

// S-0399 criterion 1: flai queue remove takes a story out and drops the
// batch that holds it; a story not enqueued is refused, named.
func TestQueueRemoveDequeuesAndDropsTheBatch(t *testing.T) {
	root := queueProject(t)
	t.Setenv("FLAI_AGENT", "alex")
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	writeQueueState(t, root, mergequeue.State{
		Enqueued: []mergequeue.Enqueued{{Story: "S-0001", By: "orchestrator", At: at}, {Story: "S-0004", By: "alex", At: at}},
		Batch:    &mergequeue.Batch{Stories: []string{"S-0001", "S-0004"}, Step: "verify"},
	})
	out, errOut, code := runIn(t, root, "queue", "remove", "S-0001")
	if code != 0 || out != "S-0001 removed from the merge queue by alex at 2026-09-15T21:00:00Z; the batch that holds it is dropped\n" {
		t.Fatalf("flai queue remove: %d %q %s", code, out, errOut)
	}
	st := readQueueState(t, root)
	if len(st.Enqueued) != 1 || st.Enqueued[0].Story != "S-0004" {
		t.Errorf("enqueued %+v", st.Enqueued)
	}
	if st.Batch == nil || st.Batch.Dropped == nil || st.Batch.Dropped.Reason != "S-0001 was removed from the queue by alex" || st.Batch.Dropped.By != "alex" {
		t.Errorf("batch %+v", st.Batch)
	}

	_, errOut, code = runIn(t, root, "queue", "remove", "S-0001")
	if code == 0 || !strings.Contains(errOut, "S-0001 is not enqueued in the merge queue, so there is nothing to remove") {
		t.Errorf("a second remove: %d %s", code, errOut)
	}

	js, errOut, code := runIn(t, root, "queue", "remove", "S-0004", "--json")
	var ch mergequeue.Change
	if code != 0 || json.Unmarshal([]byte(js), &ch) != nil {
		t.Fatalf("flai queue remove --json: %d %s %s", code, js, errOut)
	}
	if ch.Action != "remove" || ch.Story != "S-0004" || !ch.Changed || len(ch.State.Enqueued) != 0 {
		t.Errorf("remove --json: %+v", ch)
	}
}

// S-0399 criterion 4: flai queue hold holds the queue with who and when,
// flai queue shows it held, and flai queue resume clears it; each says when
// nothing changed.
func TestQueueHoldAndResume(t *testing.T) {
	root := queueProject(t)
	t.Setenv("FLAI_AGENT", "alex")
	out, errOut, code := runIn(t, root, "queue", "hold")
	if code != 0 || out != "merge queue held by alex at 2026-09-15T21:00:00Z: it builds no batch until flai queue resume\n" {
		t.Fatalf("flai queue hold: %d %q %s", code, out, errOut)
	}
	if st := readQueueState(t, root); !st.Held || st.HeldBy != "alex" || !st.HeldAt.Equal(time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)) {
		t.Errorf("held %+v", st)
	}
	t.Setenv("FLAI_AGENT", "orchestrator")
	if out, _, _ := runIn(t, root, "queue", "hold"); out != "merge queue held already, by alex at 2026-09-15T21:00:00Z; nothing changed\n" {
		t.Errorf("a second hold says %q", out)
	}
	if out, _, _ := runIn(t, root, "queue"); !strings.HasPrefix(out, "held by alex at 2026-09-15T21:00:00Z: no batch is built until flai queue resume\neligible (1)\n") {
		t.Errorf("flai queue while held:\n%s", out)
	}
	js, _, _ := runIn(t, root, "queue", "--json")
	var q mergequeue.Queue
	if err := json.Unmarshal([]byte(js), &q); err != nil || !q.Held || q.HeldBy != "alex" || q.HeldAt.IsZero() {
		t.Errorf("flai queue --json while held: %v %s", err, js)
	}

	js, errOut, code = runIn(t, root, "queue", "resume", "--json")
	var ch mergequeue.Change
	if code != 0 || json.Unmarshal([]byte(js), &ch) != nil {
		t.Fatalf("flai queue resume --json: %d %s %s", code, js, errOut)
	}
	if ch.Action != "resume" || !ch.Changed || ch.By != "orchestrator" || ch.State.Held || ch.State.HeldBy != "" {
		t.Errorf("resume --json: %+v", ch)
	}
	if out, _, _ := runIn(t, root, "queue", "resume"); out != "merge queue is not held; nothing changed\n" {
		t.Errorf("a second resume says %q", out)
	}
	if out, _, _ := runIn(t, root, "queue"); strings.Contains(out, "held") {
		t.Errorf("flai queue after resume:\n%s", out)
	}
}

// S-0399: flai queue shows the batch's result and why it was dropped.
func TestQueueShowsTheBatchsResultAndDrop(t *testing.T) {
	root := queueProject(t)
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	batch := &mergequeue.Batch{Stories: []string{"S-0001", "S-0004"}, Step: "verify"}
	batch.Result = &mergequeue.Result{StoppedAt: "integration", Summary: "TestShip failed", FinishedAt: at}
	batch.Dropped = &mergequeue.Dropped{Reason: "S-0001 was accepted by hand by alex", By: "alex", At: at}
	writeQueueState(t, root, mergequeue.State{Batch: batch})
	out, _, _ := runIn(t, root, "queue")
	want := "batch: S-0001, S-0004, step verify\n  result: failed at integration, finished 2026-09-15T20:00:00Z: TestShip failed\n  dropped by alex at 2026-09-15T20:00:00Z: S-0001 was accepted by hand by alex\n"
	if !strings.Contains(out, want) {
		t.Errorf("flai queue does not say %q:\n%s", want, out)
	}
	writeQueueState(t, root, mergequeue.State{Batch: &mergequeue.Batch{Stories: []string{"S-0001"}, Step: "land", Result: &mergequeue.Result{Passed: true, FinishedAt: at}}})
	if out, _, _ := runIn(t, root, "queue"); !strings.Contains(out, "batch: S-0001, step land\n  result: passed, finished 2026-09-15T20:00:00Z\n") {
		t.Errorf("flai queue with a passed batch:\n%s", out)
	}
}

// orchestratedQueueProject is orchestratedStoryInReview, with accept_reviews
// on when permit says so, and its S-0001 verified at the head of its branch,
// which head is, so that the merge queue finds it eligible.
func orchestratedQueueProject(t *testing.T, permit bool) (root, head string) {
	t.Helper()
	root, head = orchestratedStoryInReview(t, permit)
	if err := verify.SaveReport(openProject(t, root), verify.Report{Story: "S-0001", Commit: head, Passed: true}); err != nil {
		t.Fatal(err)
	}
	return root, head
}

// assertNotEnqueued fails unless the project at root has written no merge
// queue state, as a refused enqueue writes none.
func assertNotEnqueued(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(mergequeue.StatePath(filepath.Join(root, ".flai-cache"))); err == nil {
		t.Errorf("a refused add wrote the state: %+v", readQueueState(t, root).Enqueued)
	}
}

// S-0403 criterion 1, ADR-0093: with accept_reviews on, flai queue add --by
// orchestrator enqueues a story it can vouch for, recording the commit its
// verifier checked and its evidence as given with the story.
func TestQueueAddByTheOrchestratorRecordsItsEvidence(t *testing.T) {
	root, head := orchestratedQueueProject(t, true)
	evidence := "Verdict: pass, tests and lint clean\n- 1: `cli/feature.go`"
	js, errOut, code := runStdin(t, root, evidence+"\n", "queue", "add", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-", "--json")
	var ch mergequeue.Change
	if code != 0 || json.Unmarshal([]byte(js), &ch) != nil {
		t.Fatalf("flai queue add --by orchestrator: %d %s %s", code, js, errOut)
	}
	if !ch.Changed || ch.By != "orchestrator" || len(ch.State.Enqueued) != 1 {
		t.Errorf("the change: %+v", ch)
	}
	want := mergequeue.Enqueued{Story: "S-0001", By: "orchestrator", At: time.Date(2026, 9, 19, 2, 0, 0, 0, time.UTC), Verified: head, Evidence: evidence}
	if st := readQueueState(t, root); len(st.Enqueued) != 1 || st.Enqueued[0] != want {
		t.Errorf("enqueued %+v, want %+v", st.Enqueued, want)
	}

	file := evidenceFile(t, "Verdict: other\n- 1: `cli/feature.go`\n")
	out, errOut, code := runIn(t, root, "queue", "add", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", file)
	if code != 0 || out != "S-0001 is enqueued already, by orchestrator at 2026-09-19T02:00:00Z; nothing changed\n" {
		t.Errorf("a second add: %d %q %s", code, out, errOut)
	}
	if st := readQueueState(t, root); len(st.Enqueued) != 1 || st.Enqueued[0].Evidence != evidence {
		t.Errorf("a second add keeps the first evidence: %+v", st.Enqueued)
	}
}

// S-0403 criterion 1: flai queue add --by orchestrator is refused while
// accept_reviews is off, naming the permission, and writes nothing.
func TestQueueAddByTheOrchestratorOnlyUnderAcceptReviews(t *testing.T) {
	root, head := orchestratedQueueProject(t, false)
	_, errOut, code := runIn(t, root, "queue", "add", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n"))
	want := "rule: the orchestrator enqueues a story for the merge queue to land only with orchestration.permissions.accept_reviews, which is off: ask the operator with thread_open on S-0001"
	if code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("flai queue add with accept_reviews off: %d, want %q:\n%s", code, want, errOut)
	}
	assertNotEnqueued(t, root)
}

// S-0403 criterion 1, ADR-0093: the orchestrator's enqueue is refused, and
// nothing written, as its acceptance would be: without the verified commit
// or the evidence, with evidence that has no Verdict line, that leaves a
// criterion out or names no changed file for it, or with a verified commit
// that is not the branch's head.
func TestQueueAddByTheOrchestratorHoldsTheEvidenceToTheAcceptancesRules(t *testing.T) {
	root, head := orchestratedQueueProject(t, true)
	good := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	initial := gitIn(t, root, "rev-parse", "main")
	for name, c := range map[string]struct {
		args []string
		want []string
	}{
		"no evidence":       {[]string{"--verified", head}, []string{"rule: the orchestrator enqueues S-0001 only with the commit its verifier checked", "give both, --verified and --evidence"}},
		"no verified":       {[]string{"--evidence", good}, []string{"give both, --verified and --evidence"}},
		"no verdict":        {[]string{"--verified", head, "--evidence", evidenceFile(t, "- 1: `cli/feature.go`\n")}, []string{"rule: the orchestrator's evidence for S-0001 cannot be read as ADR-0093 asks: the evidence has no Verdict line"}},
		"a heading":         {[]string{"--verified", head, "--evidence", evidenceFile(t, "# Report\nVerdict: pass\n- 1: `cli/feature.go`\n")}, []string{"rule: line 1 of the orchestrator's evidence is a heading (# Report)"}},
		"criterion missing": {[]string{"--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 2: `cli/feature.go`\n")}, []string{"S-0001 is not enqueued: the orchestrator's acceptance of it with this commit and evidence would be refused", "criterion_unevidenced: the evidence has no item for acceptance criterion 1 (the feature ships)"}},
		"file unchanged":    {[]string{"--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/other.go`\n")}, []string{"criterion_unevidenced: the evidence for acceptance criterion 1 (the feature ships) names no file the branch changes (it names cli/other.go)"}},
		"not the head":      {[]string{"--verified", initial, "--evidence", good}, []string{"unverified: --verified " + initial[:12] + " is not the head of story/S-0001, which is " + head[:12]}},
	} {
		t.Run(name, func(t *testing.T) {
			_, errOut, code := runIn(t, root, append([]string{"queue", "add", "S-0001", "--by", "orchestrator"}, c.args...)...)
			if code == 0 {
				t.Fatal("the orchestrator's enqueue was not refused")
			}
			for _, want := range c.want {
				if !strings.Contains(errOut, want) {
					t.Errorf("the refusal does not say %q:\n%s", want, errOut)
				}
			}
			assertNotEnqueued(t, root)
		})
	}
}

// S-0403 criterion 1: --verified and --evidence are refused without --by
// orchestrator, and under FLAI_ROLE=orchestrate any other --by is, as flai
// accept refuses them.
func TestQueueAddRefusesTheOrchestratorsFlagsToAnyoneElse(t *testing.T) {
	root, head := orchestratedQueueProject(t, true)
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	_, errOut, code := runIn(t, root, "queue", "add", "S-0001", "--evidence", ev)
	if code == 0 || !strings.Contains(errOut, "rule: --verified and --evidence are the orchestrator's enqueue: give them with --by orchestrator, or leave them out") {
		t.Errorf("--evidence without --by orchestrator: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "queue", "add", "S-0001", "--by", "alex", "--verified", head); code == 0 || !strings.Contains(errOut, "rule: --verified and --evidence are the orchestrator's enqueue") {
		t.Errorf("--verified with --by alex: %d %s", code, errOut)
	}
	t.Setenv("FLAI_ROLE", "orchestrate")
	for given, args := range map[string][]string{"--by alex": {"--by", "alex"}, "no --by": nil} {
		_, errOut, code := runIn(t, root, append([]string{"queue", "add", "S-0001"}, args...)...)
		if want := "rule: the orchestrator enqueues a story only as itself, so give --by orchestrator, not " + given; code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("%s under FLAI_ROLE=orchestrate: %d, want %q:\n%s", given, code, want, errOut)
		}
	}
	assertNotEnqueued(t, root)
}
