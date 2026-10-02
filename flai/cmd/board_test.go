package cmd

import (
	"encoding/json"
	"strings"
	"testing"
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
