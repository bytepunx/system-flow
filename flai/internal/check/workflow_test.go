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
		{"a timer of no state", swap("sum: [review]}", "sum: [checking]}"), workflow.CodeTimer, lineOf(t, base, "sum: [review]}"), `timer review's sum names "checking", which is no state of the file`},
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

// edited is text with each of the lines in after put after its first line in
// turn, failing the test when text lacks one.
func edited(t *testing.T, text string, after ...[2]string) string {
	t.Helper()
	for _, a := range after {
		if !strings.Contains(text, a[0]) {
			t.Fatalf("the workflow file lacks %q", a[0])
		}
		text = strings.Replace(text, a[0], a[0]+a[1], 1)
	}
	return text
}

// lineStarting is the whole line of text, with its newline, that starts with
// prefix, failing the test when none does: an edit can then follow the line
// whatever the rest of it gives.
func lineStarting(t *testing.T, text, prefix string) string {
	t.Helper()
	for line := range strings.Lines(text) {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	t.Fatalf("the workflow file has no line starting %q", prefix)
	return ""
}

// qaBoard is the baseline file with a state of the project's own, qa, a
// waiting state between in-progress and review, which no rule finds fault
// with.
func qaBoard(t *testing.T) string {
	t.Helper()
	base := baselineText(t)
	return edited(t, base,
		[2]string{lineStarting(t, base, "  - {name: review, "), "  - {name: qa, kind: waiting}\n"},
		[2]string{"  - {from: ready, to: backlog}\n", "  - {from: in-progress, to: qa}\n  - {from: qa, to: review}\n"})
}

// writeBoard writes content as repo's workflow file and opens the project
// again as flai check opens it, without refusing the file.
func writeBoard(t *testing.T, repo *workitem.Repo, content string) *workitem.Repo {
	t.Helper()
	writeFile(t, filepath.Join(repo.Root, filepath.FromSlash(boardRel)), content)
	again, err := workitem.OpenUnchecked(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	return again
}

// qaProject is cancelledProject running qaBoard, with S-0002 moved through
// ready and in-progress into qa.
func qaProject(t *testing.T) *workitem.Repo {
	t.Helper()
	repo := writeBoard(t, cancelledProject(t), qaBoard(t))
	if err := repo.FlowErr(); err != nil {
		t.Fatalf("qaBoard is refused: %v", err)
	}
	s, err := repo.Get("S-0002")
	if err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress, "qa"} {
		if _, err := repo.Transition(s, to, "alex", "", now); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

// boardJudgments are the findings of Run on repo that judge the workflow, or
// an item or the board against it: those under a workflow. code, and
// item.sequence and board.order.
func boardJudgments(t *testing.T, repo *workitem.Repo) []Finding {
	t.Helper()
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	return slices.DeleteFunc(res.Findings, func(f Finding) bool {
		return !strings.HasPrefix(f.Rule, "workflow.") && f.Rule != "item.sequence" && f.Rule != "board.order"
	})
}

// S-0755: with the project's own workflow file, which flai can run, flai
// check judges the items against it: a story in a state of the file's own,
// entered by the file's transitions, gives no finding.
func TestItemsAreJudgedAgainstTheProjectsWorkflowFile(t *testing.T) {
	repo := qaProject(t)
	if got := boardJudgments(t, repo); len(got) != 0 {
		t.Errorf("S-0002 in qa, a state of the project's file: %+v", got)
	}
	// Without the file the same story is judged against the baseline, which
	// has no qa: the judgment the file's state clears.
	if err := os.Remove(filepath.Join(repo.Root, filepath.FromSlash(boardRel))); err != nil {
		t.Fatal(err)
	}
	baseline, err := workitem.OpenUnchecked(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := hasFinding(boardJudgments(t, baseline), codeItem, "S-0002's status names state qa"); !ok {
		t.Errorf("S-0002 in qa on the baseline gives no workflow.item: %+v", boardJudgments(t, baseline))
	}
}

// S-0755: with the project's workflow file refused, flai check reports the
// file's finding and leaves out the judgments of the items against the
// workflow, in a whole run and in a checked write's, since the repository
// runs the baseline in place of the board the project wrote (ADR-0307 §2).
func TestARefusedWorkflowFileGivesItsFindingAlone(t *testing.T) {
	repo := writeBoard(t, qaProject(t), edited(t, qaBoard(t),
		[2]string{"  - {name: cancelled, kind: cancelled}\n", "  - {name: stuck, kind: waiting}\n"},
		[2]string{"  - {from: qa, to: review}\n", "  - {from: review, to: stuck}\n"}))
	if repo.FlowErr() == nil {
		t.Fatal("a file with a dead end is not refused")
	}
	got := boardJudgments(t, repo)
	if len(got) != 1 || got[0].Rule != codeDeadEnd || got[0].Path != boardRel {
		t.Errorf("want the file's %s alone, at %s: %+v", codeDeadEnd, boardRel, got)
	}
	res, err := RunItemWrite(repo, now, Write{Items: []string{"S-0002"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if f.Rule == codeItem || f.Rule == "item.sequence" || f.Rule == "board.order" {
			t.Errorf("a checked write of S-0002 judges it against the baseline: %+v", f)
		}
	}
}

// S-0755: a project without the workflow file runs the baseline and gets no
// finding or advisory about the file it lacks (ADR-0307 §5).
func TestNoWorkflowFileGivesNoAdvisory(t *testing.T) {
	repo, err := workitem.OpenUnchecked(cancelledProject(t).Root)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if strings.HasPrefix(f.Rule, "workflow.") || f.Path == boardRel {
			t.Errorf("a project without %s: %+v", boardRel, f)
		}
	}
}

// suiteTiers are test tiers whose suites the workflow file's are held
// against: static and behavior with a manual tier beside behavior's, soak
// of a manual tier alone, and perf of a tier that is not manual beside one
// that is.
const suiteTiers = "tests:\n" +
	"  - {name: lint, command: [scripts/lint.sh], paths: [\"**\"], suite: static}\n" +
	"  - {name: unit, command: [scripts/unit.sh], paths: [\"**\"], suite: behavior}\n" +
	"  - {name: e2e, command: [scripts/e2e.sh], paths: [\"**\"], suite: behavior, manual: true}\n" +
	"  - {name: soak, command: [scripts/soak.sh], paths: [\"**\"], manual: true}\n" +
	"  - {name: bench, command: [scripts/bench.sh], paths: [\"**\"], suite: perf}\n" +
	"  - {name: bench-long, command: [scripts/bench.sh], paths: [\"**\"], suite: perf, manual: true}\n"

// suitesBoard is the baseline file with each all it gives swapped: review
// for the in-progress → review verify's suites, done for the review → done
// verify, and at for verify_at's value.
func suitesBoard(t *testing.T, review, done, at string) string {
	t.Helper()
	text := baselineText(t)
	for _, s := range []struct{ old, repl string }{
		{"    verify: {suites: all, severity: warn}\n", "    verify: {suites: " + review + ", severity: warn}\n"},
		{"verify: all}\n", "verify: " + done + "}\n"},
		{"verify_at: {batch: all, ci: all}\n", "verify_at: " + at + "\n"},
	} {
		if !strings.Contains(text, s.old) {
			t.Fatalf("the baseline lacks %q", s.old)
		}
		text = strings.Replace(text, s.old, s.repl, 1)
	}
	return text
}

// suiteFindings are the workflow.suite and workflow.suite-unnamed findings of
// a project whose manifest gives extra and whose workflow file is board, and
// the project, opened as flai check opens it.
func suiteFindings(t *testing.T, extra, board string) ([]Finding, *workitem.Repo) {
	t.Helper()
	repo := planningProject(t, extra)
	if board != "" {
		repo = writeBoard(t, repo, board)
	}
	return slices.DeleteFunc(workflowFindings(t, repo), func(f Finding) bool {
		return f.Rule != codeSuite && f.Rule != codeSuiteUnnamed
	}), repo
}

// S-0806 (ADR-0312 §9): a suite the workflow file names that no test tier
// has is a workflow.suite error on the line it is named on, naming the
// place, the suites the tiers have, and what to do; a suite of a tier that
// is not manual which nothing names is one workflow.suite-unnamed warning,
// naming only its tiers that are not manual; a suite of manual tiers alone
// gets none. Neither refuses the file.
func TestTheSuitesTheWorkflowFileNamesAreHeldAgainstTheTiers(t *testing.T) {
	board := suitesBoard(t, "[static, behaviour]", "[static]", "{batch: [static], ci: [static]}")
	got, repo := suiteFindings(t, suiteTiers, board)
	if err := repo.FlowErr(); err != nil {
		t.Fatalf("a file naming an unknown suite is refused: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want the error on behaviour and warnings on behavior and perf: %+v", got)
	}
	f, ok := hasFinding(got, codeSuite, "transition in-progress → review's verify names suite behaviour, which no test tier has", "the suites of system-flow.yaml's test tiers are static, behavior, soak, perf", "add suite: behaviour to the tier", "give all for every suite")
	if !ok {
		t.Fatalf("no workflow.suite on behaviour in %+v", got)
	}
	if want := lineOf(t, board, "verify: {suites: [static, behaviour]"); f.Level != Error || f.Path != boardRel || f.Line != want {
		t.Errorf("workflow.suite: level %s at %s:%d; want an error at %s:%d", f.Level, f.Path, f.Line, boardRel, want)
	}
	at := lineOf(t, board, "verify_at:")
	for _, w := range []struct{ suite, of string }{{"behavior", "of tier unit,"}, {"perf", "of tier bench,"}} {
		f, ok := hasFinding(got, codeSuiteUnnamed, "suite "+w.suite+", "+w.of, "so no step runs it", "give all for every suite", "manual: true")
		if !ok {
			t.Errorf("no workflow.suite-unnamed on %s naming only its tier that is not manual in %+v", w.suite, got)
			continue
		}
		if f.Level != Warning || f.Path != boardRel || f.Line != at {
			t.Errorf("workflow.suite-unnamed on %s: level %s at %s:%d; want a warning at %s:%d", w.suite, f.Level, f.Path, f.Line, boardRel, at)
		}
	}
	// Two tiers of one suite, neither manual, are named in one warning.
	twoUnit := strings.Replace(suiteTiers, "suite: behavior, manual: true}", "suite: behavior}", 1)
	got, _ = suiteFindings(t, twoUnit, board)
	if _, ok := hasFinding(got, codeSuiteUnnamed, "suite behavior, of tiers unit, e2e,"); !ok {
		t.Errorf("no one warning on behavior naming unit and e2e: %+v", got)
	}
	// A run reports both, and the error stops flai check.
	if res, err := Run(repo, now); err != nil {
		t.Fatal(err)
	} else if _, ok := hasFinding(res.Findings, codeSuite, "suite behaviour"); !ok {
		t.Errorf("Run gives no workflow.suite: %+v", res.Findings)
	}
}

// S-0806 (ADR-0312 §9): a file that names all anywhere runs every suite, so
// no tier's suite is unnamed; all is never an unknown suite, and a suite no
// tier has is still an error beside it.
func TestAWorkflowFileNamingAllGivesNoUnnamedSuite(t *testing.T) {
	board := suitesBoard(t, "[static, behaviour]", "all", "{batch: [static], ci: [static]}")
	got, _ := suiteFindings(t, suiteTiers, board)
	if len(got) != 1 || got[0].Rule != codeSuite || !strings.Contains(got[0].Message, "names suite behaviour") {
		t.Errorf("want the workflow.suite on behaviour alone: %+v", got)
	}
	if got, _ := suiteFindings(t, suiteTiers, suitesBoard(t, "[static, behavior, perf]", "[static]", "{batch: all}")); len(got) != 0 {
		t.Errorf("every suite named, and all at the batch: %+v", got)
	}
}

// S-0806 (ADR-0312 §9): the baseline names only all, so it gives neither
// finding, as the project's file or with no file; a manifest with no test
// tiers makes every suite the file names unknown.
func TestTheBaselineGivesNoSuiteFinding(t *testing.T) {
	for _, board := range []string{"", baselineText(t)} {
		if got, _ := suiteFindings(t, suiteTiers, board); len(got) != 0 {
			t.Errorf("the baseline (as the project's file: %t): %+v", board != "", got)
		}
		if got, _ := suiteFindings(t, "", board); len(got) != 0 {
			t.Errorf("the baseline with no test tiers (as the project's file: %t): %+v", board != "", got)
		}
	}
	board := strings.Replace(suitesBoard(t, "[static]", "all", "{ci: [static]}"), "runs: {kluster:", "runs: {suites: [static], kluster:", 1)
	got, _ := suiteFindings(t, "", board)
	if len(got) != 3 {
		t.Fatalf("want an error on static at each place that names it: %+v", got)
	}
	for _, at := range []string{"transition in-progress → review's verify", "state review's runs", "verify_at.ci"} {
		if _, ok := hasFinding(got, codeSuite, at+" names suite static", "system-flow.yaml declares no test tiers"); !ok {
			t.Errorf("no workflow.suite on %s's static saying the manifest has no tiers: %+v", at, got)
		}
	}
}
