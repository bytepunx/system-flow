package cmd

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// boardColumn is the lines flai board printed under the column named state,
// up to the next column's header.
func boardColumn(out, state string) []string {
	var lines []string
	in := false
	for line := range strings.SplitSeq(out, "\n") {
		header := line != "" && !strings.HasPrefix(line, " ")
		switch {
		case header && in:
			return lines
		case header:
			in = line == state || strings.HasPrefix(line, state+" (")
		case in:
			lines = append(lines, line)
		}
	}
	return lines
}

// S-0744 criteria 1 and 6: on a board with a state renamed and one added
// between in-progress and review, flai board prints the workflow's columns
// in its order, and an open story's task plan in the added state and in the
// renamed state that waits for the acceptance.
func TestBoardReadsTheWorkflowsColumnsAndKinds(t *testing.T) {
	root := renamedProject(t)
	mustRun(t, root, "task", "new", "Piece", "--story", "S-0001")
	for _, to := range []string{"ready", "in-progress", "testing"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	out := mustRun(t, root, "board")
	var headers []string
	for line := range strings.SplitSeq(out, "\n") {
		if line != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "pull order") {
			headers = append(headers, strings.Fields(line)[0])
		}
	}
	if want := []string{"backlog", "ready", "in-progress", "testing", "acceptance-review", "shipped", "abandoned"}; !slices.Equal(headers, want) {
		t.Errorf("the board's columns are %v, want the workflow's %v:\n%s", headers, want, out)
	}
	for _, state := range []string{"testing", "acceptance-review"} {
		if state == "acceptance-review" {
			mustRun(t, root, "move", "S-0001", "review")
			out = mustRun(t, root, "board")
		}
		col := strings.Join(boardColumn(out, state), "\n")
		if !strings.Contains(col, "S-0001") || !strings.Contains(col, "         tasks ") {
			t.Errorf("the open story in %s is listed with its task plan:\n%s", state, out)
		}
	}
}

// S-0744 criteria 1 and 6: a story's seeds leave out its tasks in the
// workflow's cancelled state, whatever it is named, as flai touches suggest
// reads them.
func TestTouchesSeedsLeaveOutATaskInTheCancelledKind(t *testing.T) {
	root := renamedProject(t)
	for _, args := range [][]string{
		{"task", "new", "Kept", "--story", "S-0001"}, {"task", "new", "Dropped", "--story", "S-0001"},
		{"touches", "T-0001", "docs/guide.md"}, {"touches", "T-0002", "web"},
		{"move", "T-0002", "abandoned", "--reason", "dropped", "--yes"},
	} {
		mustRun(t, root, args...)
	}
	if got := stateOf(t, root, "T-0002"); got != "abandoned" {
		t.Fatalf("T-0002 is %q, want abandoned", got)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	repo.Workflow = testWorkflow
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	seeds, err := storySeeds(repo, story, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"docs/guide.md"}; !slices.Equal(seeds, want) {
		t.Errorf("the seeds are %v, want %v: the abandoned task's touches are left out", seeds, want)
	}
}
