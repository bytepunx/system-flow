package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// templateBoard is the workflow file the template ships, relative to this
// package.
var templateBoard = filepath.Join("..", "..", "..", "template", "root", "design", filepath.FromSlash(File))

// S-0755: the template ships the baseline flai carries as its
// design/workflows/board.md, byte for byte, front matter and body
// (ADR-0301 §7), so a project made from it runs what a project without it
// runs.
func TestBaselineIsTheTemplatesBoard(t *testing.T) {
	shipped, err := os.ReadFile(templateBoard)
	if err != nil {
		t.Fatal(err)
	}
	if string(shipped) != baselineFile {
		t.Errorf("the template's %s is not the embedded baseline, front matter and body; run: cp flai/internal/workflow/baseline/board.md template/root/design/%s", File, File)
	}
}

// todaysMoves is the transition map workitem enforced before the board was
// data (rules.go's allowed, S-0737), for every item type: a task also goes
// from in-progress to done. review to cancelled is left out: today only a
// parent's cancellation takes an item there (ADR-0028), which the baseline
// says with by: [cascade].
var todaysMoves = map[string][]string{
	"backlog":     {"ready", "cancelled"},
	"ready":       {"in-progress", "backlog", "cancelled"},
	"in-progress": {"review", "ready", "cancelled"},
	"review":      {"done", "in-progress"},
	"cancelled":   {"backlog"},
}

// today reports whether an item of type typ could move from one state to
// the other before the board was data, a cascade apart.
func today(typ, from, to string) bool {
	return slices.Contains(todaysMoves[from], to) || (typ == TypeTask && from == "in-progress" && to == "done")
}

// S-0737: the baseline flai carries parses, and its states are today's six,
// in board order, with ADR-0301 §2's kinds and flags.
func TestBaselineStatesAreTodays(t *testing.T) {
	w := Baseline()
	if w.Path != BaselinePath {
		t.Errorf("path %q, want %q", w.Path, BaselinePath)
	}
	want := []struct {
		name, kind            string
		initial, order, limit bool
	}{
		{"backlog", KindQueue, true, true, false},
		{"ready", KindQueue, false, true, true},
		{"in-progress", KindActive, false, false, true},
		{"review", KindWaiting, false, false, true},
		{"done", KindAccepted, false, false, false},
		{"cancelled", KindCancelled, false, false, false},
	}
	if len(w.States) != len(want) {
		t.Fatalf("states %v, want six", w.Names())
	}
	for i, s := range want {
		got := w.States[i]
		if got.Name != s.name || got.Kind != s.kind || got.Initial != s.initial || got.Order != s.order || got.Limit != s.limit {
			t.Errorf("state %d: %+v, want %+v", i, got, s)
		}
		if !slices.Equal(got.Types, Types) || len(got.Aliases) != 0 {
			t.Errorf("state %s: types %v aliases %v, want every type and no alias", got.Name, got.Types, got.Aliases)
		}
	}
}

// S-0737: pair by pair, for every item type, the baseline has a transition
// from one state to another exactly where today's map has the move, and
// review to cancelled only as a cascade. A change to the baseline that adds
// or drops a move fails here.
func TestBaselineMovesAreTodays(t *testing.T) {
	w := Baseline()
	for _, typ := range Types {
		for _, from := range w.Names() {
			for _, to := range w.Names() {
				tr, has := w.Find(typ, from, to)
				cascadeOnly := has && slices.Equal(tr.By, []string{ByCascade})
				if from == "review" && to == "cancelled" {
					if !cascadeOnly {
						t.Errorf("%s %s to %s: has %v by %v, want a transition by [cascade] alone", typ, from, to, has, tr.By)
					}
					continue
				}
				if got, want := has && !cascadeOnly, today(typ, from, to); got != want {
					t.Errorf("%s %s to %s: the baseline moves it %v, today %v", typ, from, to, got, want)
				}
			}
		}
	}
}

// S-0737: who makes each move and what flai does with it are ADR-0301 §3's.
func TestBaselineByAndEffects(t *testing.T) {
	w := Baseline()
	cases := []struct {
		typ, from, to string
		by, effects   []string
	}{
		{TypeStory, "backlog", "ready", []string{ByAnyone}, []string{EffectCommit}},
		{TypeStory, "ready", "in-progress", []string{ByAnyone}, []string{EffectPull}},
		{TypeStory, "in-progress", "review", []string{ByAnyone}, []string{}},
		{TypeEpic, "review", "done", []string{ByAcceptance}, []string{EffectAccept}},
		{TypeStory, "review", "done", []string{ByAcceptance}, []string{EffectAccept}},
		{TypeTask, "review", "done", []string{ByAnyone}, []string{}},
		{TypeTask, "in-progress", "done", []string{ByAnyone}, []string{}},
		{TypeStory, "review", "in-progress", []string{ByAnyone}, []string{EffectSendBack}},
		{TypeStory, "in-progress", "ready", []string{ByAnyone}, []string{}},
		{TypeStory, "ready", "backlog", []string{ByAnyone}, []string{}},
		{TypeStory, "backlog", "cancelled", []string{ByAnyone}, []string{EffectCancel}},
		{TypeTask, "ready", "cancelled", []string{ByAnyone}, []string{EffectCancel}},
		{TypeEpic, "in-progress", "cancelled", []string{ByAnyone}, []string{EffectCancel}},
		{TypeStory, "review", "cancelled", []string{ByCascade}, []string{EffectCancel}},
		{TypeStory, "cancelled", "backlog", []string{ByAnyone}, []string{EffectReopen}},
	}
	for _, c := range cases {
		tr, ok := w.Find(c.typ, c.from, c.to)
		if !ok {
			t.Errorf("%s %s to %s: no transition", c.typ, c.from, c.to)
			continue
		}
		if !slices.Equal(tr.By, c.by) || !slices.Equal(tr.Effects, c.effects) {
			t.Errorf("%s %s to %s: by %v effects %v, want by %v effects %v", c.typ, c.from, c.to, tr.By, tr.Effects, c.by, c.effects)
		}
	}
	for _, typ := range []string{TypeEpic, TypeStory} {
		if tr, ok := w.WithEffect(typ, EffectPull); !ok || tr.To != "in-progress" {
			t.Errorf("%s: the pull is %+v, %v", typ, tr, ok)
		}
		if tr, ok := w.WithEffect(typ, EffectAccept); !ok || tr.To != "done" {
			t.Errorf("%s: the acceptance is %+v, %v", typ, tr, ok)
		}
	}
	if _, ok := w.WithEffect(TypeTask, EffectAccept); ok {
		t.Error("a task has an acceptance; its move to done is anyone's")
	}
}

// S-0738: each transition lists the rules Repo.Move applies on it today, in
// the order it applies them, as ADR-0301 §3's table names them, each one the
// vocabulary has, with the severity it gives. A rule written as name:warn
// sets warn; one written bare takes its rule's default.
func TestBaselineRulesAreTodays(t *testing.T) {
	w := Baseline()
	cases := []struct {
		typ, from, to string
		rules         []string
	}{
		{TypeStory, "backlog", "ready", []string{"criteria-present", "not-draft", "after-done"}},
		{TypeStory, "ready", "in-progress", []string{"not-blocked", "not-held", "task-after"}},
		{TypeStory, "in-progress", "review", []string{"has-task", "no-open-question", "clean-worktree", "criteria-ticked:warn", "kluster-verified"}},
		{TypeEpic, "review", "done", []string{"children-closed", "criteria-ticked"}},
		{TypeStory, "review", "done", []string{"children-closed", "criteria-ticked"}},
		{TypeTask, "review", "done", []string{"children-closed"}},
		{TypeTask, "in-progress", "done", []string{"children-closed"}},
		{TypeStory, "review", "in-progress", []string{"reason-given", "not-blocked", "task-after"}},
		{TypeStory, "in-progress", "ready", []string{"criteria-present", "not-draft", "after-done"}},
		{TypeStory, "ready", "backlog", nil},
		{TypeStory, "backlog", "cancelled", []string{"reason-given"}},
		{TypeTask, "ready", "cancelled", []string{"reason-given"}},
		{TypeEpic, "in-progress", "cancelled", []string{"reason-given"}},
		{TypeStory, "review", "cancelled", []string{"reason-given"}},
		{TypeStory, "cancelled", "backlog", []string{"parent-open"}},
	}
	for _, c := range cases {
		tr, ok := w.Find(c.typ, c.from, c.to)
		if !ok {
			t.Errorf("%s %s to %s: no transition", c.typ, c.from, c.to)
			continue
		}
		var got []string
		for _, r := range tr.Rules {
			if _, known := LookupRule(r.Rule); !known {
				t.Errorf("%s %s to %s: rule %s is not in the vocabulary", c.typ, c.from, c.to, r.Rule)
			}
			name := r.Rule
			if r.Severity != "" {
				name += ":" + r.Severity
			}
			got = append(got, name)
		}
		if !slices.Equal(got, c.rules) {
			t.Errorf("%s %s to %s: rules %v, want %v", c.typ, c.from, c.to, got, c.rules)
		}
	}
	tr, _ := w.Find(TypeStory, "in-progress", "review")
	if sev := tr.Rules[3].SeverityOf(); sev != SeverityWarn {
		t.Errorf("criteria-ticked to review is %s, want warn", sev)
	}
	if sev := tr.Rules[4].SeverityOf(); sev != SeverityWarn {
		t.Errorf("kluster-verified to review is %s, want its default, warn", sev)
	}
	tr, _ = w.Find(TypeStory, "review", "done")
	if sev := tr.Rules[1].SeverityOf(); sev != SeverityRefuse {
		t.Errorf("criteria-ticked to done is %s, want its default, refuse", sev)
	}
}

// S-0737: the baseline's timers are today's lead, cycle, queue, and review
// time, in ADR-0301 §5's forms.
func TestBaselineTimersAreTodays(t *testing.T) {
	w := Baseline()
	want := []string{
		"lead: created to completed",
		"cycle: {first: in-progress} to completed",
		"queue: sum [ready] until {first: in-progress}",
		"review: sum [review]",
	}
	if len(w.Timers) != len(want) {
		t.Fatalf("%d timers, want %d", len(w.Timers), len(want))
	}
	for i, tm := range w.Timers {
		if got := timerForm(tm); got != want[i] {
			t.Errorf("timer %d: %s, want %s", i, got, want[i])
		}
		if tm.Counts != CountsAccepted || !slices.Equal(tm.Types, Types) {
			t.Errorf("timer %s: counts %s types %v", tm.Name, tm.Counts, tm.Types)
		}
	}
}

func timerForm(tm Timer) string {
	if tm.From != nil {
		return fmt.Sprintf("%s: %s to %s", tm.Name, tm.From, tm.To)
	}
	s := fmt.Sprintf("%s: sum %v", tm.Name, tm.Sum)
	if tm.Until != nil {
		s += fmt.Sprintf(" until %s", tm.Until)
	}
	return s
}

// S-0806: the baseline names the suites and cluster steps today runs, as
// ADR-0312 §8 gives them: every suite at warn on the move to review, every
// suite on the acceptance, the cluster's deploy and checks on entering
// review, and every suite in the batch and in CI, none at publish. No other
// transition gives verify and no other state gives runs.
func TestBaselineVerifyAndRunsAreTodays(t *testing.T) {
	w := Baseline()
	all := func(s *Suites) bool { return s != nil && s.All && len(s.Names) == 0 }
	toReview, ok := w.Find(TypeStory, "in-progress", "review")
	if !ok || toReview.Verify == nil || !all(&toReview.Verify.Suites) || toReview.Verify.Severity != SeverityWarn {
		t.Errorf("the move to review verifies %+v, want all at warn", toReview.Verify)
	}
	for _, typ := range []string{TypeEpic, TypeStory} {
		accept, ok := w.WithEffect(typ, EffectAccept)
		if !ok || accept.Verify == nil || !all(&accept.Verify.Suites) || accept.Verify.Severity != "" {
			t.Errorf("%s: the acceptance verifies %+v, want all at refuse", typ, accept.Verify)
		}
	}
	for _, tr := range w.Transitions {
		if tr.Verify == nil {
			continue
		}
		accepts := slices.Contains(tr.Effects, EffectAccept)
		toReview := slices.Equal(tr.From, []string{"in-progress"}) && tr.To == "review"
		if !accepts && !toReview {
			t.Errorf("transition %v to %s gives verify %+v, want none", tr.From, tr.To, tr.Verify)
		}
	}
	for _, s := range w.States {
		if s.Name != "review" {
			if s.Runs != nil {
				t.Errorf("state %s runs %+v, want nothing", s.Name, s.Runs)
			}
			continue
		}
		if s.Runs == nil || s.Runs.Suites != nil || !slices.Equal(s.Runs.Kluster, []string{KlusterDeploy, KlusterChecks}) {
			t.Errorf("review runs %+v, want kluster [deploy, checks] and no suites", s.Runs)
		}
	}
	if v := w.VerifyAt; !all(v.Batch) || !all(v.CI) || v.Publish != nil {
		t.Errorf("verify_at is batch %+v publish %+v ci %+v, want batch all, ci all, no publish", v.Batch, v.Publish, v.CI)
	}
}

// S-0737: each call to Baseline gives a copy of its own, so that one
// caller's change does not reach the next.
func TestBaselineIsACopy(t *testing.T) {
	a := Baseline()
	a.States[0].Name = "changed"
	a.States[1].Types[0] = "changed"
	a.Transitions[0].From[0] = "changed"
	a.Transitions[0].Rules[0].Rule = "changed"
	a.Timers[1].From.State = "changed"
	a.States[3].Runs.Kluster[0] = "changed"
	a.Transitions[2].Verify.Suites.All = false
	a.Transitions[2].Verify.Severity = "changed"
	a.VerifyAt.Batch.All = false
	a.VerifyAt.CI.All = false
	b := Baseline()
	if b.States[0].Name != "backlog" || b.States[1].Types[0] != TypeEpic || b.Transitions[0].From[0] != "backlog" ||
		b.Transitions[0].Rules[0].Rule != "criteria-present" || b.Timers[1].From.State != "in-progress" {
		t.Errorf("a change to one copy reached the next: %+v", b)
	}
	if b.States[3].Runs.Kluster[0] != KlusterDeploy || !b.Transitions[2].Verify.Suites.All || b.Transitions[2].Verify.Severity != SeverityWarn ||
		!b.VerifyAt.Batch.All || !b.VerifyAt.CI.All {
		t.Errorf("a change to one copy's verify, runs, or verify_at reached the next: %+v %+v %+v", b.States[3].Runs, b.Transitions[2].Verify, b.VerifyAt)
	}
}
