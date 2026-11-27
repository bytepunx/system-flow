package workitem

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// ADR-0312 §3: the suites an acceptance needs are those the verify of the
// transition with the accept effect lists: the baseline's all; a list as
// given, as a copy; none when that transition gives no verify, though another
// does, or when the type has no acceptance; and all for a file that gives none
// of ADR-0312's keys, which reads as the baseline.
func TestAcceptSuitesAreTheVerifyOfTheAcceptance(t *testing.T) {
	const accept = "effects: [accept], rules: [children-closed, criteria-ticked]}"
	parsed := func(t *testing.T, old, new string) workflow.Workflow {
		t.Helper()
		if !strings.Contains(testFlowFile, old) {
			t.Fatalf("the fixture workflow has no %q", old)
		}
		w, p := workflow.Parse("design/workflows/board.md", strings.Replace(testFlowFile, old, new, 1))
		if p != nil {
			t.Fatalf("the edited workflow does not parse: %s", p)
		}
		return w
	}
	for _, tc := range []struct {
		name string
		w    workflow.Workflow
		typ  string
		want *workflow.Suites
	}{
		{"the baseline", workflow.Baseline(), Story, &workflow.Suites{All: true}},
		{"the baseline's epic", workflow.Baseline(), Epic, &workflow.Suites{All: true}},
		{"a task, which has no acceptance", workflow.Baseline(), Task, nil},
		{"a list", parsed(t, accept, strings.TrimSuffix(accept, "}")+", verify: [static, behavior]}"), Story,
			&workflow.Suites{Names: []string{"static", "behavior"}}},
		{"no verify on the acceptance", parsed(t, "{from: doing, to: testing}", "{from: doing, to: testing, verify: [static]}"), Story, nil},
		{"a file with no suite keys", *testFlow(t), Story, &workflow.Suites{All: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := AcceptSuites(tc.w, tc.typ)
			if got != nil {
				got.Line = 0
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AcceptSuites = %+v; want %+v", got, tc.want)
			}
		})
	}
	t.Run("a copy", func(t *testing.T) {
		w := parsed(t, accept, strings.TrimSuffix(accept, "}")+", verify: [static]}")
		AcceptSuites(w, Story).Names[0] = "changed"
		if got := AcceptSuites(w, Story).Names; !reflect.DeepEqual(got, []string{"static"}) {
			t.Errorf("after a caller changed its copy, AcceptSuites names %q; want static", got)
		}
	})
}
