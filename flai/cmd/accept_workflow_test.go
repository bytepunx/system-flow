package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0744 criteria 3 and 6: on a board with a state renamed and one added
// between in-progress and review, flai accept refuses a story in the added
// state in today's words, naming it and the state the acceptance leaves
// from, and its dry run of a story in that state finds nothing to stop it
// and changes nothing.
func TestAcceptTakesAStoryOnlyFromTheStateTheAcceptanceLeaves(t *testing.T) {
	root := renamedProject(t)
	mustRun(t, root, "story", "new", "Shipped")
	file := filepath.Join(root, "wip/kanban/stories/S-0002-shipped.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, root, "task", "new", "--story", "S-0002", "Step")
	for _, m := range [][]string{{"S-0002", "ready"}, {"S-0002", "in-progress"}, {"S-0002", "testing"}, {"T-0001", "ready"}, {"T-0001", "in-progress"}, {"T-0001", "shipped"}} {
		mustRun(t, root, "move", m[0], m[1])
	}

	_, errOut, code := runIn(t, root, "accept", "S-0002", "--by", "alex")
	if want := "S-0002 is testing; a story is accepted from acceptance-review"; code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("a story in testing must be refused with %q: exit %d\n%s", want, code, errOut)
	}
	if got := stateOf(t, root, "S-0002"); got != "testing" {
		t.Errorf("a refused acceptance changes nothing: status %q, want testing", got)
	}

	mustRun(t, root, "move", "S-0002", "acceptance-review")
	if out := mustRun(t, root, "accept", "S-0002", "--by", "alex", "--dry-run"); !strings.Contains(out, "dry run: nothing changed") || strings.Contains(out, "block") {
		t.Errorf("the dry run of a story in acceptance-review finds nothing to stop it:\n%s", out)
	}
	if got := stateOf(t, root, "S-0002"); got != "acceptance-review" {
		t.Errorf("a dry run changes nothing: status %q, want acceptance-review", got)
	}
}
