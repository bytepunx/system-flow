package cmd

import (
	"strings"
	"testing"
)

// S-0744 criteria 4 and 6: on a board with a state renamed and one added
// between in-progress and review, flai formula stop reads in progress, done,
// and cancelled by their kinds: the instance and the stopping step's task
// are in testing, the added state of the active kind; the step's task goes
// to shipped along testing's transition into the accepted kind; and the
// steps after it go to abandoned, the cancelled kind's state, with the
// stop's reason.
func TestFormulaStopFollowsTheWorkflowsTransitions(t *testing.T) {
	useWorkflow(t, renamedWorkflow)
	asOperator(t)
	root := formulaProject(t, map[string]string{"publish-like.md": publishLike})
	for _, args := range [][]string{
		{"formula", "run", "publish-like", "--ready"},
		{"move", "S-0001", "in-progress"}, {"move", "S-0001", "testing"},
		{"stream", "open", "S-0001", "--no-branch"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "shipped"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"}, {"move", "T-0002", "testing"},
	} {
		mustRun(t, root, args...)
	}
	out, errOut, code := runStdin(t, root, "not met\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2")
	if code != 0 {
		t.Fatalf("formula stop: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"Step evaluate stopped the run with exit 2; skipped T-0003, T-0004.",
		"T-0002 → shipped (step evaluate, exit 2)",
		"T-0003 → abandoned: step evaluate stopped the run with exit 2",
		"T-0004 → abandoned: step evaluate stopped the run with exit 2",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("formula stop's output lacks %q:\n%s", want, out)
		}
	}
	for id, want := range map[string]string{"T-0001": "shipped", "T-0002": "shipped", "T-0003": "abandoned", "T-0004": "abandoned", "S-0001": "testing"} {
		if got := stateOf(t, root, id); got != want {
			t.Errorf("%s is %q, want %s", id, got, want)
		}
	}
	for _, id := range []string{"T-0003", "T-0004"} {
		if task := itemFile(t, root, id); !strings.Contains(task, ": moved to abandoned: step evaluate stopped the run with exit 2\n") {
			t.Errorf("%s is not cancelled with the stop's reason:\n%s", id, task)
		}
	}
}

// S-0744 criterion 4: a stop is refused, writing nothing, when the stopping
// step's task is in a state of a kind other than active, though the state is
// named by an alias the board keeps: acceptance-review, once review, is of
// the waiting kind.
func TestFormulaStopRefusesAStepNotActiveOnTheWorkflow(t *testing.T) {
	useWorkflow(t, renamedWorkflow)
	asOperator(t)
	root := formulaProject(t, map[string]string{"publish-like.md": publishLike})
	for _, args := range [][]string{
		{"formula", "run", "publish-like", "--ready"},
		{"move", "S-0001", "in-progress"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "shipped"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"}, {"move", "T-0002", "testing"}, {"move", "T-0002", "review"},
	} {
		mustRun(t, root, args...)
	}
	before := snapshotItems(t, root)
	_, errOut, code := runStdin(t, root, "not met\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2")
	if want := "rule: T-0002, the task of step evaluate, is acceptance-review, not in progress"; code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("formula stop must refuse with %q: exit %d\n%s", want, code, errOut)
	}
	after := snapshotItems(t, root)
	for path, was := range before {
		if after[path] != was {
			t.Errorf("a refused stop changed %s", path)
		}
	}
}
