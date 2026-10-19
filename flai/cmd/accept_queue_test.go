package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
)

// enqueueWithBatch writes a queue state with S-0001 and S-0005 enqueued and
// in the batch under verification.
func enqueueWithBatch(t *testing.T, root string) {
	t.Helper()
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	writeQueueState(t, root, mergequeue.State{
		Enqueued: []mergequeue.Enqueued{{Story: "S-0001", By: "orchestrator", At: at}, {Story: "S-0005", By: "orchestrator", At: at}},
		Batch:    &mergequeue.Batch{Stories: []string{"S-0001", "S-0005"}, Step: "verify"},
	})
}

// assertLeftTheQueue fails unless S-0001 is out of the queue, the batch that
// held it is dropped as accepted by hand by alex, and its outcome says so.
func assertLeftTheQueue(t *testing.T, root string) {
	t.Helper()
	st := readQueueState(t, root)
	if len(st.Enqueued) != 1 || st.Enqueued[0].Story != "S-0005" {
		t.Errorf("enqueued %+v, want S-0005 only", st.Enqueued)
	}
	if st.Batch == nil || st.Batch.Dropped == nil || st.Batch.Dropped.Reason != "S-0001 was accepted by hand by alex" || st.Batch.Dropped.By != "alex" {
		t.Errorf("batch %+v, want it dropped", st.Batch)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Story != "S-0001" || st.Outcomes[0].Outcome != mergequeue.OutcomeAcceptedByHand || st.Outcomes[0].Detail != "accepted by hand by alex" {
		t.Errorf("outcomes %+v", st.Outcomes)
	}
}

// S-0399 criterion 1, ADR-0173: flai accept by hand takes the story out of
// the queue, drops the batch that holds it, and records its outcome.
func TestAcceptByHandTakesTheStoryOutOfTheQueue(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	enqueueWithBatch(t, root)
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("flai accept: %d %s\n%s", code, errOut, out)
	}
	assertLeftTheQueue(t, root)
}

// S-0399 criterion 1: under pull-request acceptance too, the story merged
// through the remote leaves the queue.
func TestAcceptThroughTheRemoteTakesTheStoryOutOfTheQueue(t *testing.T) {
	root, _, _, _ := remoteProject(t, "")
	enqueueWithBatch(t, root)
	acceptJSON(t, root)
	assertLeftTheQueue(t, root)
}

// A queue state that cannot be written after the merge is warned of, and the
// acceptance is not failed for it.
func TestAcceptWarnsOfAQueueStateItCannotWrite(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	path := mergequeue.StatePath(filepath.Join(root, ".flai-cache"))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"enqueued": [`), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 || !strings.Contains(out, "accepted S-0001: done") {
		t.Fatalf("flai accept: %d %s\n%s", code, errOut, out)
	}
	var warned bool
	for _, line := range strings.Split(errOut, "\n") {
		if strings.Contains(line, path) && strings.Contains(line, `"level":"WARN"`) && strings.Contains(line, `"component":"accept"`) && strings.Contains(line, `"item":"S-0001"`) {
			warned = true
		}
	}
	if !warned {
		t.Errorf("no warning naming the queue's state:\n%s", errOut)
	}
}
