package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storeS0001Report stores S-0001's last flai verify report in the project at
// root: run at commit, passed, with the tier test passed and the tier
// integration, marked queue, in the state given.
func storeS0001Report(t *testing.T, root, commit string, integration verify.State) {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	rep := verify.Report{Story: "S-0001", Commit: commit, Passed: true, Steps: []verify.Step{
		{Name: "check", State: verify.Passed},
		{Name: "test", Tier: true, State: verify.Passed},
		{Name: "integration", Tier: true, State: integration},
	}}
	if err := verify.SaveReport(repo, rep); err != nil {
		t.Fatal(err)
	}
}

// queueTiersRefusal is the start of the refusal of S-0001's acceptance by
// hand while the merge queue batches and its tiers marked queue did not run
// at its head (ADR-0221).
const queueTiersRefusal = "S-0001 cannot be accepted yet: S-0001's tiers marked queue have not run at its head: "

// I-0203, ADR-0221: while the merge queue batches, the operator's acceptance
// of a story whose close-out left integration to the queue at its branch's
// head is refused, naming integration and flai verify --all, before anything
// is merged; --dry-run lists it as the blocker queue_tiers. Once a report at
// its head ran every tier, it is accepted.
func TestAHandAcceptanceWaitsForTheQueueTiersWhileTheQueueBatches(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mergeQueueOn(t, root)
	storeS0001Report(t, root, head, verify.Queued)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	want := queueTiersRefusal + "its last flai verify did not run tier integration, marked queue; while the merge queue batches, a hand acceptance waits for them: run flai verify S-0001 --all"
	if code == 0 || !strings.Contains(errOut, want) || !strings.Contains(errOut, "flai queue add S-0001") || !strings.Contains(errOut, "--without-queue-tiers (ADR-0221)") {
		t.Fatalf("a hand acceptance without integration at the head must be refused with %q: %d\n%s", want, code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("flai accept --dry-run: %d %s\n%s", code, errOut, out)
	}
	var res accepted
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if b := res.QueueTiersBlocker; b == nil || b.Code != preview.BlockQueueTiers || len(res.QueueTiers) != 1 || res.QueueTiers[0] != "integration" || len(res.Blockers) != 1 || res.Blockers[0] != b.Message {
		t.Errorf("--dry-run must list the blocker queue_tiers naming integration: %+v %q %q", res.QueueTiersBlocker, res.QueueTiers, res.Blockers)
	}
	assertNothingAccepted(t, root, mainBefore)

	storeS0001Report(t, root, head, verify.Passed)
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("with integration passed at the head the story is accepted: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done || strings.Contains(story.Body, "### Accepted without the queue tiers") {
		t.Errorf("the story is accepted with no override in its Notes: %s %v\n%s", story.Status, story.Archived, story.Body)
	}
}

// ADR-0221: a report at a commit other than the branch's head, and no report
// at all, refuse the hand acceptance too, saying why.
func TestAHandAcceptanceWaitsForAReportAtTheStorysHead(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mergeQueueOn(t, root)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 || !strings.Contains(errOut, queueTiersRefusal+"no flai verify has run on it;") {
		t.Errorf("with no report the hand acceptance must be refused: %d\n%s", code, errOut)
	}
	storeS0001Report(t, root, mainBefore, verify.Passed)
	want := queueTiersRefusal + "its last flai verify ran at " + mainBefore[:7] + ", not at its head " + head[:7] + ";"
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("a report at an older commit must refuse the hand acceptance with %q: %d\n%s", want, code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// ADR-0221: --without-queue-tiers lands the story anyway and writes the tiers
// that did not run, who accepted, and when in its Notes; with no report it
// writes why.
func TestWithoutQueueTiersLandsTheStoryAndRecordsItInItsNotes(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mergeQueueOn(t, root)
	storeS0001Report(t, root, head, verify.Queued)
	if out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--without-queue-tiers"); code != 0 {
		t.Fatalf("--without-queue-tiers lands the story: %d %s\n%s", code, errOut, out)
	}
	section := "### Accepted without the queue tiers\n\n- Tiers: integration\n- By: alex\n- At: 2026-09-17T09:00:00Z\n"
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done || !strings.Contains(story.Body, "## Notes\n") || !strings.HasSuffix(story.Body, section) {
		t.Errorf("the archived story's Notes must end with\n%s\ngot %s %v\n%s", section, story.Status, story.Archived, story.Body)
	}

	none, _ := orchestratedStoryInReview(t, true)
	mergeQueueOn(t, none)
	if out, errOut, code := runInAt(t, none, landedAt, "move", "S-0001", "done", "--by", "alex", "--without-queue-tiers"); code != 0 {
		t.Fatalf("flai move done --without-queue-tiers lands the story: %d %s\n%s", code, errOut, out)
	}
	section = "### Accepted without the queue tiers\n\n- Why: no flai verify has run on it\n- By: alex\n- At: 2026-09-17T09:00:00Z\n"
	if story := acceptedS0001(t, none); !story.Archived || !strings.HasSuffix(story.Body, section) {
		t.Errorf("with no report the Notes must end with\n%s\ngot\n%s", section, story.Body)
	}
}

// ADR-0221: while the queue does not batch, a hand acceptance is as before,
// whatever the report left to the queue, and --without-queue-tiers writes
// nothing.
func TestAHandAcceptanceIsUnchangedWhileTheQueueDoesNotBatch(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mergeQueueOn(t, filepath.Join(t.TempDir(), "another"))
	storeS0001Report(t, root, head, verify.Queued)
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--without-queue-tiers"); code != 0 {
		t.Fatalf("with the queue off for the project the hand acceptance is unchanged: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || strings.Contains(story.Body, "### Accepted without the queue tiers") {
		t.Errorf("the story is accepted with nothing written for the override:\n%s", story.Body)
	}
}

// ADR-0221: --without-queue-tiers is the operator's alone, refused with --by
// orchestrator and with --batch before anything changes.
func TestWithoutQueueTiersIsRefusedForTheOrchestratorAndTheQueuesLanding(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev, "--without-queue-tiers"}, "rule: --without-queue-tiers is the operator's: the orchestrator never lands S-0001 without the tiers marked queue"},
		{[]string{"accept", "S-0001", "--by", "alex", "--batch", landingBatch, "--without-queue-tiers"}, "rule: --without-queue-tiers is the operator's acceptance by hand, and --batch " + landingBatch + " is the merge queue's landing"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("flai %v must be refused with %q: %d\n%s", c.args, c.want, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)
}
