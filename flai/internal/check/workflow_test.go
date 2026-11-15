package check

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// boardRel is the workflow file's path in the test projects, whose design
// folder is design.
const boardRel = "design/workflows/board.md"

// workflowFindings are the findings of the workflow rules alone: the
// workflow's (workflow) and each item's workflow.item (itemStates).
func workflowFindings(t *testing.T, repo *workitem.Repo) []Finding {
	t.Helper()
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	c := &checker{repo: repo, res: &Result{}}
	c.workflow()
	for _, it := range items {
		c.itemStates(it)
	}
	sortFindings(c.res.Findings)
	return c.res.Findings
}

// baselineText is the baseline workflow file flai carries.
func baselineText(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// withFlow is cancelledProject running the baseline as edit changes it.
func withFlow(t *testing.T, edit func(w *workflow.Workflow)) *workitem.Repo {
	t.Helper()
	repo := cancelledProject(t)
	w := workflow.Baseline()
	edit(&w)
	repo.Workflow = &w
	return repo
}

// hasFinding is the finding of code whose message holds every one of says.
func hasFinding(fs []Finding, code string, says ...string) (Finding, bool) {
	for _, f := range fs {
		if f.Rule != code {
			continue
		}
		if !slices.ContainsFunc(says, func(s string) bool { return !strings.Contains(f.Message, s) }) {
			return f, true
		}
	}
	return Finding{}, false
}

// S-0741: the baseline flai carries gives no workflow finding, run as the
// repository's workflow and as the project's own file, and neither do the
// items of a project on it.
func TestBaselineWorkflowHasNoFindings(t *testing.T) {
	repo := cancelledProject(t)
	if got := workflowFindings(t, repo); len(got) != 0 {
		t.Errorf("the baseline: %+v", got)
	}
	writeFile(t, filepath.Join(repo.Root, filepath.FromSlash(boardRel)), baselineText(t))
	if got := workflowFindings(t, repo); len(got) != 0 {
		t.Errorf("the baseline as the project's file: %+v", got)
	}
	good, err := workitem.Open("../metrics/testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	if got := workflowFindings(t, good); len(got) != 0 {
		t.Errorf("the good fixture: %+v", got)
	}
}

// S-0741: a workflow written to fail each rule gives the rule's finding, at
// the path the project's file would have, naming the state, the transition,
// or the timer, and saying what to do.
func TestAWorkflowThatFailsARuleGivesItsFinding(t *testing.T) {
	all := slices.Clone(workflow.Types)
	cases := []struct {
		name  string
		edit  func(w *workflow.Workflow)
		code  string
		level string
		says  []string
	}{
		{"a state no transition reaches", func(w *workflow.Workflow) {
			w.States = append(w.States, workflow.State{Name: "parked", Kind: workflow.KindQueue, Types: all})
			w.Transitions = append(w.Transitions, workflow.Transition{From: []string{"parked"}, To: "backlog", Types: all, By: []string{workflow.ByAnyone}})
		}, codeUnreachable, Warning, []string{"state parked is entered by no chain of transitions from the initial state backlog", "add a transition into it"}},
		{"a state reached only from an unreachable one", func(w *workflow.Workflow) {
			w.States = append(w.States, workflow.State{Name: "parked", Kind: workflow.KindQueue, Types: all}, workflow.State{Name: "shelved", Kind: workflow.KindQueue, Types: all})
			w.Transitions = append(w.Transitions,
				workflow.Transition{From: []string{"parked"}, To: "shelved", Types: all, By: []string{workflow.ByAnyone}},
				workflow.Transition{From: []string{"shelved"}, To: "cancelled", Types: all, By: []string{workflow.ByAnyone}})
		}, codeUnreachable, Warning, []string{"state shelved is entered by no chain"}},
		{"a state with no way to close", func(w *workflow.Workflow) {
			w.States = append(w.States, workflow.State{Name: "stuck", Kind: workflow.KindWaiting, Types: all})
			w.Transitions = append(w.Transitions, workflow.Transition{From: []string{"review"}, To: "stuck", Types: all, By: []string{workflow.ByAnyone}})
		}, codeDeadEnd, Error, []string{"state stuck leads by no chain of transitions to done or cancelled", "add a transition out of it"}},
		{"no pull", func(w *workflow.Workflow) {
			w.Transitions[1].Effects = nil
		}, codeNoPull, Error, []string{"no transition gives stories the pull effect", "from a queue state to an active state effects: [pull]"}},
		{"a pull for tasks alone", func(w *workflow.Workflow) {
			w.Transitions[1].Types = []string{workflow.TypeTask}
			w.Transitions = append(w.Transitions, workflow.Transition{From: []string{"ready"}, To: "in-progress", Types: []string{workflow.TypeEpic, workflow.TypeStory}, By: []string{workflow.ByAnyone}})
		}, codeNoPull, Error, []string{"no transition gives stories the pull effect"}},
		{"two pulls", func(w *workflow.Workflow) {
			w.Transitions = append(w.Transitions, workflow.Transition{From: []string{"backlog"}, To: "in-progress", Types: all, By: []string{workflow.ByAnyone}, Effects: []string{workflow.EffectPull}})
		}, codeNoPull, Error, []string{"transitions[12] (backlog to in-progress) gives stories the pull effect, as transitions[1] (ready to in-progress) does", "take pull out of the effects of all but one"}},
		{"a pull from an active state", func(w *workflow.Workflow) {
			w.Transitions[1].Effects = nil
			w.Transitions[7].Effects = []string{workflow.EffectPull}
		}, codeNoPull, Error, []string{"transitions[7] (in-progress to ready) gives stories the pull effect from in-progress, an active state", "the pull leaves a queue state"}},
		{"a pull into a queue state", func(w *workflow.Workflow) {
			w.Transitions[1].Effects = nil
			w.Transitions[7].Effects = []string{workflow.EffectPull}
		}, codeNoPull, Error, []string{"gives stories the pull effect into ready, a queue state", "the pull enters an active state"}},
		{"no acceptance", func(w *workflow.Workflow) {
			w.Transitions[3].Effects = nil
		}, codeNoAccept, Error, []string{"no transition gives stories the accept effect", "from a waiting state to the accepted state effects: [accept]"}},
		{"an acceptance anyone makes", func(w *workflow.Workflow) {
			w.Transitions[3].By = []string{workflow.ByAnyone}
		}, codeNoAccept, Error, []string{"transitions[3] (review to done) accepts stories by [anyone]", "give it by: [acceptance]"}},
		{"an acceptance from an active state", func(w *workflow.Workflow) {
			w.Transitions[3].Effects = nil
			w.Transitions[5].Types = all
			w.Transitions[5].By = []string{workflow.ByAcceptance}
			w.Transitions[5].Effects = []string{workflow.EffectAccept}
		}, codeNoAccept, Error, []string{"transitions[5] (in-progress to done) gives stories the accept effect from in-progress, an active state", "the acceptance leaves a waiting state"}},
		{"two acceptances", func(w *workflow.Workflow) {
			w.Transitions[4].Types = all
			w.Transitions[4].Effects = []string{workflow.EffectAccept}
		}, codeNoAccept, Error, []string{"transitions[4] (review to done) gives stories the accept effect, as transitions[3] (review to done) does"}},
		{"no initial state", func(w *workflow.Workflow) {
			w.States[0].Initial = false
		}, codeState, Error, []string{"the workflow has no initial state", "mark exactly one queue state initial: true"}},
		{"two initial states", func(w *workflow.Workflow) {
			w.States[1].Initial = true
		}, codeState, Error, []string{"the workflow has 2 initial states, backlog, ready"}},
		{"two accepted states", func(w *workflow.Workflow) {
			w.States[5].Kind = workflow.KindAccepted
		}, codeState, Error, []string{"the workflow has 2 accepted states, done, cancelled"}},
		{"no cancelled state", func(w *workflow.Workflow) {
			w.States[5].Kind = workflow.KindAccepted
		}, codeState, Error, []string{"the workflow has no cancelled state", "give exactly one state kind: cancelled"}},
		{"a transition out of the accepted state", func(w *workflow.Workflow) {
			w.Transitions = append(w.Transitions, workflow.Transition{From: []string{"done"}, To: "review", Types: all, By: []string{workflow.ByAnyone}})
		}, codeState, Error, []string{"transitions[12] (done to review) leaves done, the accepted state", "remove done from the transition's from"}},
		{"a rule flai lacks", func(w *workflow.Workflow) {
			w.Transitions[0].Rules = append(w.Transitions[0].Rules, workflow.RuleRef{Rule: "signed-off"})
		}, workflow.CodeRule, Error, []string{"transitions[0] (backlog to ready) names rule signed-off, which flai does not implement"}},
		{"a timer summing no state", func(w *workflow.Workflow) {
			w.Timers[3].Sum = []string{"checking"}
		}, workflow.CodeTimer, Error, []string{`timer review's sum names "checking", which is no state of the workflow`, "or add the state"}},
		{"a timer from no state", func(w *workflow.Workflow) {
			w.Timers[1].From = &workflow.Point{Edge: workflow.EdgeFirst, State: "started"}
		}, workflow.CodeTimer, Error, []string{`timer cycle's from names "started"`}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := workflowFindings(t, withFlow(t, tc.edit))
			f, ok := hasFinding(got, tc.code, tc.says...)
			if !ok {
				t.Fatalf("no %s saying %q in %+v", tc.code, tc.says, got)
			}
			if f.Level != tc.level || f.Path != boardRel || f.Line != 1 {
				t.Errorf("%s: level %s, at %s:%d; want %s, at %s:1", tc.code, f.Level, f.Path, f.Line, tc.level, boardRel)
			}
		})
	}
}

// S-0741: the project's workflow file is parsed in place of the repository's
// workflow: a problem keeping it from parsing is one finding under the
// problem's code, on its line, and nothing else is checked on the file; a file
// that parses gives the rules' findings on its entries' lines.
func TestTheProjectsWorkflowFileIsChecked(t *testing.T) {
	base := baselineText(t)
	swap := func(old, repl string) string {
		t.Helper()
		if !strings.Contains(base, old) {
			t.Fatalf("the baseline lacks %q", old)
		}
		return strings.Replace(base, old, repl, 1)
	}
	cases := []struct {
		name, content, code string
		line                int
		says                string
	}{
		{"a file that does not parse", "no front matter\n", workflow.CodeFile, 1, "the file has no front matter"},
		{"a field outside the allowed values", swap("{name: done, kind: accepted}", "{name: done, kind: finished}"), workflow.CodeFile, 10, `states[4].kind: state done's kind "finished" is not one of`},
		{"a rule flai lacks", swap("rules: [criteria-present, not-draft, after-done]}", "rules: [criteria-present, signed-off]}"), workflow.CodeRule, 13, "names rule signed-off, which flai does not implement"},
		{"a timer of no state", swap("sum: [review]}", "sum: [checking]}"), workflow.CodeTimer, 31, `timer review's sum names "checking", which is no state of the file`},
		{"no pull", swap("effects: [pull], ", ""), codeNoPull, 12, "no transition gives stories the pull effect"},
		{"a dead end", strings.Replace(swap("  - {name: cancelled, kind: cancelled}\n", "  - {name: cancelled, kind: cancelled}\n  - {name: stuck, kind: waiting}\n"), "  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog}\n  - {from: review, to: stuck}\n", 1),
			codeDeadEnd, 12, "state stuck leads by no chain of transitions"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := cancelledProject(t)
			writeFile(t, filepath.Join(repo.Root, filepath.FromSlash(boardRel)), tc.content)
			got := workflowFindings(t, repo)
			f, ok := hasFinding(got, tc.code, tc.says)
			if !ok {
				t.Fatalf("no %s saying %q in %+v", tc.code, tc.says, got)
			}
			if f.Level != Error || f.Path != boardRel || f.Line != tc.line {
				t.Errorf("%s: level %s at %s:%d; want an error at %s:%d", tc.code, f.Level, f.Path, f.Line, boardRel, tc.line)
			}
			if !strings.HasPrefix(tc.code, "workflow.no-") && tc.code != codeDeadEnd && len(got) != 1 {
				t.Errorf("a file that does not parse gives %d findings, want its problem alone: %+v", len(got), got)
			}
		})
	}
}

// S-0741: an item whose status or history names a state the workflow neither
// has nor lists as an alias is a workflow.item error on its file, naming the
// item and the state; an alias of the state its items go to clears it.
func TestAnItemInAStateTheWorkflowLacks(t *testing.T) {
	rename := func(from, to string, alias bool) func(w *workflow.Workflow) {
		return func(w *workflow.Workflow) {
			for i := range w.States {
				if w.States[i].Name == from {
					w.States[i].Name = to
					if alias {
						w.States[i].Aliases = []string{from}
					}
				}
			}
			for i := range w.Transitions {
				for j := range w.Transitions[i].From {
					if w.Transitions[i].From[j] == from {
						w.Transitions[i].From[j] = to
					}
				}
				if w.Transitions[i].To == from {
					w.Transitions[i].To = to
				}
			}
			for i := range w.Timers {
				for j := range w.Timers[i].Sum {
					if w.Timers[i].Sum[j] == from {
						w.Timers[i].Sum[j] = to
					}
				}
			}
		}
	}
	repo := withFlow(t, rename("review", "checking", false))
	story, _ := repo.Get("S-0001")
	got := workflowFindings(t, repo)
	f, ok := hasFinding(got, codeItem, "S-0001's status names state review, which the workflow neither has nor lists as an alias", "add the state to "+boardRel+", or add review to the aliases of the state its items go to")
	if !ok {
		t.Fatalf("no workflow.item on S-0001's status in %+v", got)
	}
	if rel, _ := filepath.Rel(repo.Root, story.Path); f.Level != Error || f.Path != rel || f.Line != keyLine(story.Path, "status") {
		t.Errorf("workflow.item: level %s at %s:%d; want an error on %s's status line", f.Level, f.Path, f.Line, rel)
	}
	if n := len(slices.DeleteFunc(slices.Clone(got), func(f Finding) bool { return f.Rule != codeItem })); n != 1 {
		t.Errorf("%d workflow.item findings, want one for review, which S-0001's status and history both name: %+v", n, got)
	}
	repo = withFlow(t, rename("ready", "queued", false))
	if _, ok := hasFinding(workflowFindings(t, repo), codeItem, "S-0001's transitions[0] names state ready"); !ok {
		t.Errorf("no workflow.item on S-0001's history, which entered ready: %+v", workflowFindings(t, repo))
	}
	for _, from := range []string{"review", "ready"} {
		repo = withFlow(t, rename(from, "renamed", true))
		if got := workflowFindings(t, repo); len(got) != 0 {
			t.Errorf("%s renamed with its old name as an alias: %+v", from, got)
		}
	}
}
