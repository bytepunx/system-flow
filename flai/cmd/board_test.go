package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
)

// S-0195, ADR-0067: flai board says what is accepted and not yet published,
// as information that publishing is the operator's, and tells no one to push.
func TestBoardSaysWhatIsUnpublished(t *testing.T) {
	root, _ := researchProject(t, "feature", true)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}

	board, _, _ := runIn(t, root, "board")
	if !strings.Contains(board, "accepted, not yet published: S-0001 (publishing is the operator's: git fetch, then flai release --pending)") {
		t.Errorf("flai board says what is unpublished and whose it is to publish:\n%s", board)
	}
	for _, never := range []string{"flai push", "not pushed"} {
		if strings.Contains(board, never) {
			t.Errorf("flai board says %q:\n%s", never, board)
		}
	}
	js, _, _ := runIn(t, root, "board", "--json")
	var view map[string]any
	if err := json.Unmarshal([]byte(js), &view); err != nil {
		t.Fatalf("board --json: %v %s", err, js)
	}
	if u, _ := view["unpublished"].([]any); len(u) != 1 || u[0] != "S-0001" {
		t.Errorf("board --json unpublished: %v", view["unpublished"])
	}
	if _, ok := view["unpushed"]; ok {
		t.Errorf("board --json has no unpushed: %v", view["unpushed"])
	}

	if _, errOut, code := runIn(t, root, "release", "--pending"); code != 0 {
		t.Fatalf("publish: %s", errOut)
	}
	if board, _, _ := runIn(t, root, "board"); strings.Contains(board, "not yet published") {
		t.Errorf("once published the board stops saying so:\n%s", board)
	}
}

// S-0398, ADR-0173: flai board carries the merge queue as flai queue gives
// it, in its JSON as queue, and as one line under the review column naming
// the queue's order and the batch; the reasons are flai queue's to give.
func TestBoardCarriesTheMergeQueue(t *testing.T) {
	root := queueProject(t)
	writeQueueState(t, root, mergequeue.State{Batch: &mergequeue.Batch{Stories: []string{"S-0001"}, Step: "verify"}})

	out, errOut, code := runIn(t, root, "board")
	if code != 0 || !strings.Contains(out, "\n         queue: S-0001; batch: S-0001, step verify\ndone") {
		t.Errorf("flai board names the queue's order and batch under review: %d %s\n%s", code, errOut, out)
	}
	if strings.Contains(out, "verify_missing") {
		t.Errorf("flai board gives the reasons, which are flai queue's:\n%s", out)
	}

	board, queue := boardAndQueueJSON(t, root)
	if !reflect.DeepEqual(board["queue"], queue) {
		t.Errorf("board --json queue:\n%v\nwant flai queue --json's:\n%v", board["queue"], queue)
	}
	q := queue.(map[string]any)
	if el := q["eligible"].([]any); len(el) != 1 || el[0].(map[string]any)["story"] != "S-0001" {
		t.Errorf("eligible: %v", q["eligible"])
	}
	if in := q["ineligible"].([]any); len(in) != 1 || in[0].(map[string]any)["story"] != "S-0002" || len(in[0].(map[string]any)["reasons"].([]any)) != 3 {
		t.Errorf("ineligible: %v", q["ineligible"])
	}
	if b := q["batch"].(map[string]any); b["step"] != "verify" {
		t.Errorf("batch: %v", b)
	}
}

// S-0398, ADR-0169: under pull-request acceptance flai board says the queue
// is off in its place, and its JSON lists no story as eligible.
func TestBoardSaysTheQueueIsOffUnderPullRequest(t *testing.T) {
	root := queueProject(t)
	m := filepath.Join(root, "system-flow.yaml")
	_ = os.WriteFile(m, []byte(read(t, m)+"acceptance:\n  mode: pull-request\n"), 0o644)

	out, errOut, code := runIn(t, root, "board")
	if code != 0 || !strings.Contains(out, "\n         queue off under acceptance.mode pull-request\ndone") || strings.Contains(out, "queue: ") {
		t.Errorf("flai board: %d %s\n%s", code, errOut, out)
	}
	board, queue := boardAndQueueJSON(t, root)
	q, _ := board["queue"].(map[string]any)
	if !reflect.DeepEqual(board["queue"], queue) || q["off"] != true || q["off_reason"] != mergequeue.OffUnderPullRequest || len(q["eligible"].([]any)) != 0 {
		t.Errorf("board --json queue: %v", board["queue"])
	}
}

// boardAndQueueJSON is flai board --json and flai queue --json at root.
func boardAndQueueJSON(t *testing.T, root string) (board map[string]any, queue any) {
	t.Helper()
	js, errOut, code := runIn(t, root, "board", "--json")
	if code != 0 || json.Unmarshal([]byte(js), &board) != nil {
		t.Fatalf("flai board --json: %d %s %s", code, js, errOut)
	}
	qs, errOut, code := runIn(t, root, "queue", "--json")
	if code != 0 || json.Unmarshal([]byte(qs), &queue) != nil {
		t.Fatalf("flai queue --json: %d %s %s", code, qs, errOut)
	}
	return board, queue
}
