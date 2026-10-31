package cmd

import (
	"strings"
	"testing"
)

// S-0552 criterion 2: while the merge queue batches and is not held, flai move
// <story> done from review enqueues the story as flai accept does, and lands
// nothing: the story stays in review and the batch under way is not dropped.
func TestMoveToDoneEnqueuesWhileTheQueueBatches(t *testing.T) {
	root, _ := orchestratedQueueProject(t, true)
	mergeQueueOn(t, root)
	queueUnderWay(t, root, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	out, errOut, code := runIn(t, root, "move", "S-0001", "done", "--by", "alex")
	if want := "enqueued S-0001 in the merge queue at position 2 as alex, and it stays in review"; code != 0 || !strings.Contains(out, want) {
		t.Fatalf("flai move done must answer %q: %d %s\n%s", want, code, errOut, out)
	}
	assertNothingAccepted(t, root, mainBefore)
	st := readQueueState(t, root)
	if len(st.Enqueued) != 2 || st.Enqueued[1].Story != "S-0001" || st.Enqueued[1].By != "alex" {
		t.Errorf("enqueued %+v, want S-0001 after S-0005, by alex", st.Enqueued)
	}
	assertBatchUnderWay(t, st)
}
