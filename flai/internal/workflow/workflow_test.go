package workflow

import (
	"slices"
	"strings"
	"testing"
)

// adrBaseline is ADR-0301 §7's baseline file, as written there.
const adrBaseline = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - from: in-progress
    to: review
    rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: in-progress, to: done, types: [task], rules: [children-closed]}
  - {from: review, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: in-progress, to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: in-progress}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: in-progress}}
  - {name: review, title: Review time, sum: [review]}
---

The board's body is prose for the operator; flai reads nothing of it.
`

func mustParse(t *testing.T, content string) Workflow {
	t.Helper()
	w, p := Parse("design/workflows/board.md", content)
	if p != nil {
		t.Fatalf("Parse: %s", p)
	}
	return w
}

// S-0737: the ADR's baseline parses into its six states, twelve transitions,
// and four timers, each with its line and its defaults filled in.
func TestParseReadsTheADRsBaseline(t *testing.T) {
	w := mustParse(t, adrBaseline)
	if w.Path != "design/workflows/board.md" || w.Title != "Board" || w.Updated != "2026-11-14" || w.Status != "active" {
		t.Errorf("header: %q %q %q %q", w.Path, w.Title, w.Updated, w.Status)
	}
	if got, want := w.Names(), []string{"backlog", "ready", "in-progress", "review", "done", "cancelled"}; !slices.Equal(got, want) {
		t.Errorf("Names: %v, want %v", got, want)
	}
	backlog := w.States[0]
	if backlog.Kind != KindQueue || !backlog.Initial || !backlog.Order || backlog.Limit || backlog.Line != 6 {
		t.Errorf("backlog: %+v", backlog)
	}
	if backlog.Title != "backlog" || !slices.Equal(backlog.Types, Types) || backlog.Aliases == nil || len(backlog.Aliases) != 0 {
		t.Errorf("backlog's defaults: title %q, types %v, aliases %#v", backlog.Title, backlog.Types, backlog.Aliases)
	}
	if r := w.States[3]; r.Kind != KindWaiting || !r.Limit || r.Order || r.Initial {
		t.Errorf("review: %+v", r)
	}
	if len(w.Transitions) != 12 {
		t.Fatalf("%d transitions, want 12", len(w.Transitions))
	}
	first := w.Transitions[0]
	if !slices.Equal(first.From, []string{"backlog"}) || first.To != "ready" || !slices.Equal(first.By, []string{ByAnyone}) || !slices.Equal(first.Types, Types) || !slices.Equal(first.Effects, []string{EffectCommit}) || first.Line != 13 {
		t.Errorf("transitions[0]: %+v", first)
	}
	toReview := w.Transitions[2]
	if toReview.Line != 15 || len(toReview.Rules) != 5 || toReview.Rules[3] != (RuleRef{Rule: "criteria-ticked", Severity: SeverityWarn}) || toReview.Rules[0] != (RuleRef{Rule: "has-task"}) {
		t.Errorf("transitions[2]: %+v", toReview)
	}
	if len(toReview.Effects) != 0 || toReview.Effects == nil {
		t.Errorf("transitions[2]'s effects: %#v, want an empty list", toReview.Effects)
	}
	if c := w.Transitions[9]; !slices.Equal(c.From, []string{"backlog", "ready", "in-progress"}) || c.To != "cancelled" {
		t.Errorf("transitions[9]: %+v", c)
	}
	if a := w.Transitions[3]; !slices.Equal(a.Types, []string{TypeEpic, TypeStory}) || !slices.Equal(a.By, []string{ByAcceptance}) {
		t.Errorf("transitions[3]: %+v", a)
	}
	if len(w.Timers) != 4 {
		t.Fatalf("%d timers, want 4", len(w.Timers))
	}
	lead, cycle, queue, review := w.Timers[0], w.Timers[1], w.Timers[2], w.Timers[3]
	if *lead.From != (Point{At: PointCreated}) || *lead.To != (Point{At: PointCompleted}) || lead.Until != nil || len(lead.Sum) != 0 || lead.Counts != CountsAccepted || !slices.Equal(lead.Types, Types) || lead.Line != 28 {
		t.Errorf("lead: %+v", lead)
	}
	if *cycle.From != (Point{Edge: EdgeFirst, State: "in-progress"}) || cycle.From.String() != "{first: in-progress}" {
		t.Errorf("cycle: %+v", cycle.From)
	}
	if queue.From != nil || !slices.Equal(queue.Sum, []string{"ready"}) || *queue.Until != (Point{Edge: EdgeFirst, State: "in-progress"}) {
		t.Errorf("queue: %+v", queue)
	}
	if review.Title != "Review time" || !slices.Equal(review.Sum, []string{"review"}) || review.Until != nil {
		t.Errorf("review: %+v", review)
	}
}

// minimal is the smallest file Parse takes: two states and a move.
const minimal = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true}
  - {name: done, kind: accepted}
transitions:
  - {from: backlog, to: done}
---
`

// withTimer is minimal with one timer entry, on line 11.
func withTimer(entry string) string {
	return strings.Replace(minimal, "  - {from: backlog, to: done}\n", "  - {from: backlog, to: done}\ntimers:\n  - "+entry+"\n", 1)
}

// S-0737: a file that does not parse is refused with the field and the line
// that keep it from parsing, under workflow.file.
func TestParseRefusesWhatDoesNotParse(t *testing.T) {
	swap := func(old, repl string) string { return strings.Replace(minimal, old, repl, 1) }
	cases := []struct {
		name, content string
		line          int
		field, says   string
	}{
		{"no front matter", "title: Board\n", 1, "", "no front matter"},
		{"front matter not closed", "---\ntitle: Board\n", 1, "", "no front matter"},
		{"bad yaml", swap("status: active\n", "status: [active\n"), 5, "", "does not parse"},
		{"not a mapping", "---\n- a\n---\n", 2, "", "not a mapping"},
		{"unknown field", swap("status: active\n", "status: active\ncolour: blue\n"), 5, "colour", "unknown field"},
		{"unknown state field", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, colour: blue}"), 7, "states[1].colour", "unknown field"},
		{"no title", swap("title: Board\n", ""), 1, "title", "needs title"},
		{"updated not a date", swap("updated: 2026-11-14", "updated: soon"), 3, "updated", "not a date"},
		{"no status", swap("status: active\n", ""), 1, "status", "needs status"},
		{"unknown status", swap("status: active", "status: live"), 4, "status", "not one of draft"},
		{"no states", "---\ntitle: Board\nupdated: 2026-11-14\nstatus: active\ntransitions:\n  - {from: a, to: b}\n---\n", 1, "states", "needs states"},
		{"state not a mapping", swap("  - {name: done, kind: accepted}", "  - done"), 7, "states[1]", "a state is a mapping"},
		{"state without a name", swap("{name: done, kind: accepted}", "{kind: accepted}"), 7, "states[1].name", "needs name"},
		{"state name not kebab-case", swap("{name: done,", "{name: Done,"), 7, "states[1].name", "not kebab-case"},
		{"state name twice", swap("{name: done,", "{name: backlog,"), 7, "states[1].name", "also the name of states[0]"},
		{"state without a kind", swap("{name: done, kind: accepted}", "{name: done}"), 7, "states[1].kind", "needs kind"},
		{"unknown kind", swap("kind: accepted", "kind: finished"), 7, "states[1].kind", `kind "finished" is not one of`},
		{"unknown state type", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, types: [bug]}"), 7, "states[1].types[0]", `type "bug"`},
		{"type twice", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, types: [task, task]}"), 7, "states[1].types[1]", "named twice"},
		{"initial on an accepted state", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, initial: true}"), 7, "states[1].initial", "only a queue state is initial"},
		{"order on an accepted state", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, order: true}"), 7, "states[1].order", "pull order"},
		{"alias a current name", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, aliases: [backlog]}"), 7, "states[1].aliases[0]", "never a current state's name"},
		{"alias twice", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, aliases: [finished, finished]}"), 7, "states[1].aliases[1]", "also an alias of state done"},
		{"no transitions", swap("transitions:\n  - {from: backlog, to: done}\n", ""), 1, "transitions", "needs transitions"},
		{"transition without from", swap("{from: backlog, to: done}", "{to: done}"), 9, "transitions[0].from", "needs from"},
		{"transition from no state", swap("{from: backlog, to: done}", "{from: ready, to: done}"), 9, "transitions[0].from", `"ready" is no state of the file`},
		{"transition from a list with no state", swap("{from: backlog, to: done}", "{from: [backlog, ready], to: done}"), 9, "transitions[0].from[1]", `"ready" is no state`},
		{"transition without to", swap("{from: backlog, to: done}", "{from: backlog}"), 9, "transitions[0].to", "needs to"},
		{"transition to no state", swap("to: done}", "to: review}"), 9, "transitions[0].to", `"review" is no state`},
		{"unknown transition type", swap("to: done}", "to: done, types: [bug]}"), 9, "transitions[0].types[0]", `type "bug"`},
		{"unknown by", swap("to: done}", "to: done, by: [robot]}"), 9, "transitions[0].by[0]", `by "robot" is not one of`},
		{"unknown effect", swap("to: done}", "to: done, effects: [merge]}"), 9, "transitions[0].effects[0]", `effect "merge"`},
		{"rule not kebab-case", swap("to: done}", "to: done, rules: [Has_Task]}"), 9, "transitions[0].rules[0]", "not kebab-case"},
		{"rule mapping without rule", swap("to: done}", "to: done, rules: [{severity: warn}]}"), 9, "transitions[0].rules[0]", "needs rule"},
		{"unknown severity", swap("to: done}", "to: done, rules: [{rule: has-task, severity: loud}]}"), 9, "transitions[0].rules[0].severity", `severity "loud"`},
		{"transition twice", swap("  - {from: backlog, to: done}\n", "  - {from: backlog, to: done}\n  - {from: backlog, to: done, types: [task]}\n"), 10, "transitions[1]", "already moves a task from backlog to done"},
		{"timer without name", withTimer("{title: Lead, from: created, to: completed}"), 11, "timers[0].name", "needs name"},
		{"timer name not kebab-case", withTimer("{name: Lead, title: Lead, from: created, to: completed}"), 11, "timers[0].name", "not kebab-case"},
		{"timer without title", withTimer("{name: lead, from: created, to: completed}"), 11, "timers[0].title", "needs title"},
		{"timer with from alone", withTimer("{name: lead, title: Lead, from: created}"), 11, "timers[0].to", "give both"},
		{"timer with to alone", withTimer("{name: lead, title: Lead, to: completed}"), 11, "timers[0].from", "give both"},
		{"timer with neither form", withTimer("{name: lead, title: Lead}"), 11, "timers[0]", "needs from and to, or sum"},
		{"timer with both forms", withTimer("{name: lead, title: Lead, from: created, to: completed, sum: [backlog]}"), 11, "timers[0].sum", "not both"},
		{"until without sum", withTimer("{name: lead, title: Lead, from: created, to: completed, until: completed}"), 11, "timers[0].until", "until bounds a sum only"},
		{"point of no kind", withTimer("{name: lead, title: Lead, from: {last: {kind: busy}}, to: completed}"), 11, "timers[0].from", `kind "busy"`},
		{"point of another shape", withTimer("{name: lead, title: Lead, from: started, to: completed}"), 11, "timers[0].from", "a point is created, completed"},
		{"unknown counts", withTimer("{name: lead, title: Lead, from: created, to: completed, counts: all}"), 11, "timers[0].counts", `counts "all"`},
		{"timer name twice", strings.Replace(withTimer("{name: wait, title: Wait, sum: [backlog]}"), "sum: [backlog]}\n", "sum: [backlog]}\n  - {name: wait, title: Again, sum: [done]}\n", 1), 12, "timers[1].name", "also timers[0]'s"},
		{"verify at the top", swap("status: active\n", "status: active\nverify: all\n"), 5, "verify", "unknown field"},
		{"verify of another word", swap("to: done}", "to: done, verify: static}"), 9, "transitions[0].verify", "a single suite is a list of one"},
		{"verify of no value", swap("to: done}", "to: done, verify: null}"), 9, "transitions[0].verify", "verify gives no value"},
		{"verify of an empty list", swap("to: done}", "to: done, verify: []}"), 9, "transitions[0].verify", "give suites, or all"},
		{"verify naming all in a list", swap("to: done}", "to: done, verify: [all, static]}"), 9, "transitions[0].verify[0]", "write all alone"},
		{"verify naming a suite not kebab-case", swap("to: done}", "to: done, verify: [Static]}"), 9, "transitions[0].verify[0]", `suite "Static" is not kebab-case`},
		{"verify naming a suite twice", swap("to: done}", "to: done, verify: [static, static]}"), 9, "transitions[0].verify[1]", "suite static is named twice"},
		{"verify naming a mapping", swap("to: done}", "to: done, verify: [static, {a: b}]}"), 9, "transitions[0].verify[1]", "named by a word"},
		{"verify's mapping with another key", swap("to: done}", "to: done, verify: {suites: all, colour: blue}}"), 9, "transitions[0].verify.colour", "unknown field"},
		{"verify's mapping without suites", swap("to: done}", "to: done, verify: {severity: warn}}"), 9, "transitions[0].verify", "needs suites"},
		{"verify's suites of no value", swap("to: done}", "to: done, verify: {suites: null}}"), 9, "transitions[0].verify.suites", "verify's suites is all or a list"},
		{"verify's suites an empty list", swap("to: done}", "to: done, verify: {suites: []}}"), 9, "transitions[0].verify.suites", "give suites, or all"},
		{"verify's suites naming a suite twice", swap("to: done}", "to: done, verify: {suites: [static, static]}}"), 9, "transitions[0].verify.suites[1]", "named twice"},
		{"verify's unknown severity", swap("to: done}", "to: done, verify: {suites: all, severity: loud}}"), 9, "transitions[0].verify.severity", `verify's severity "loud" is not one of refuse, warn`},
		{"runs not a mapping", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: all}"), 7, "states[1].runs", "runs is a mapping"},
		{"runs of no value", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: null}"), 7, "states[1].runs", "runs gives no value"},
		{"runs naming nothing", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {}}"), 7, "states[1].runs", "state done's runs names nothing"},
		{"runs with another key", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {suites: all, colour: blue}}"), 7, "states[1].runs.colour", "unknown field"},
		{"runs' suites of no value", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {suites: null}}"), 7, "states[1].runs.suites", "runs gives suites no value"},
		{"runs' suites naming all in a list", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {suites: [static, all]}}"), 7, "states[1].runs.suites[1]", "write all alone"},
		{"runs' suites an empty list", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {suites: []}}"), 7, "states[1].runs.suites", "give suites, or all"},
		{"runs' unknown kluster step", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {kluster: [deploy, test]}}"), 7, "states[1].runs.kluster[1]", `kluster step "test" is not one of deploy, checks`},
		{"runs' kluster step twice", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {kluster: [checks, checks]}}"), 7, "states[1].runs.kluster[1]", "kluster step checks is named twice"},
		{"runs' kluster empty", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {kluster: []}}"), 7, "states[1].runs.kluster", "state done's kluster is empty"},
		{"runs' kluster of no value", swap("{name: done, kind: accepted}", "{name: done, kind: accepted, runs: {kluster: null}}"), 7, "states[1].runs.kluster", "runs gives kluster no value"},
		{"verify_at not a mapping", swap("transitions:\n", "verify_at: all\ntransitions:\n"), 8, "verify_at", "verify_at is a mapping"},
		{"verify_at of no value", swap("transitions:\n", "verify_at:\ntransitions:\n"), 8, "verify_at", "verify_at gives no value"},
		{"verify_at naming no point", swap("transitions:\n", "verify_at: {}\ntransitions:\n"), 8, "verify_at", "verify_at names no point"},
		{"verify_at with another point", swap("transitions:\n", "verify_at:\n  batch: all\n  nightly: [static]\ntransitions:\n"), 10, "verify_at.nightly", "unknown field"},
		{"verify_at's point of no value", swap("transitions:\n", "verify_at: {publish: null}\ntransitions:\n"), 8, "verify_at.publish", "verify_at gives publish no value"},
		{"verify_at's point of another word", swap("transitions:\n", "verify_at: {batch: every}\ntransitions:\n"), 8, "verify_at.batch", "suites are all, or a list"},
		{"verify_at's point an empty list", swap("transitions:\n", "verify_at: {batch: []}\ntransitions:\n"), 8, "verify_at.batch", "give suites, or all"},
		{"verify_at's point naming a suite twice", swap("transitions:\n", "verify_at:\n  ci: [static, static]\ntransitions:\n"), 9, "verify_at.ci[1]", "suite static is named twice"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, p := Parse("design/workflows/board.md", tc.content)
			if p == nil {
				t.Fatal("parsed, want a problem")
			}
			if p.Code != CodeFile || p.Path != "design/workflows/board.md" || p.Line != tc.line || p.Field != tc.field || !strings.Contains(p.Message, tc.says) {
				t.Errorf("got %s (line %d, field %q), want line %d, field %q, saying %q", p, p.Line, p.Field, tc.line, tc.field, tc.says)
			}
			if !strings.Contains(p.Error(), p.Message) {
				t.Errorf("Error %q lacks the message", p.Error())
			}
		})
	}
}

// suited is a file giving verify on a transition in each of its forms, runs
// on a state with suites, kluster, or both, and verify_at.
const suited = `---
title: Board
updated: 2026-11-17
status: active
states:
  - {name: backlog, kind: queue, initial: true}
  - {name: review, kind: waiting, runs: {kluster: [deploy, checks]}}
  - name: integration
    kind: waiting
    runs:
      suites: [integration, smoke]
      kluster: [deploy]
  - {name: done, kind: accepted}
transitions:
  - {from: backlog, to: review, verify: all}
  - {from: review, to: integration, verify: [static, behavior]}
  - from: integration
    to: done
    verify:
      suites: all
      severity: warn
  - {from: review, to: done, verify: {suites: [behavior]}}
  - {from: backlog, to: done}
verify_at:
  batch: all
  ci:
    - static
    - behavior
---
`

// S-0806: a transition's verify, a state's runs, and verify_at parse in each
// of their forms, each with its line (ADR-0312 §2, §4, §6).
func TestParseReadsVerifyRunsAndVerifyAt(t *testing.T) {
	w := mustParse(t, suited)
	verify := func(i int) Verify {
		t.Helper()
		if w.Transitions[i].Verify == nil {
			t.Fatalf("transitions[%d] has no verify", i)
		}
		return *w.Transitions[i].Verify
	}
	if v := verify(0); !v.Suites.All || v.Suites.Names != nil || v.Suites.Line != 15 || v.Severity != "" || v.Line != 15 {
		t.Errorf("verify: all: %+v", v)
	}
	if v := verify(1); v.Suites.All || !slices.Equal(v.Suites.Names, []string{"static", "behavior"}) || v.Suites.Line != 16 || v.Line != 16 {
		t.Errorf("verify: [static, behavior]: %+v", v)
	}
	if v := verify(2); !v.Suites.All || v.Suites.Line != 20 || v.Severity != SeverityWarn || v.Line != 19 {
		t.Errorf("verify: {suites: all, severity: warn}: %+v", v)
	}
	if v := verify(3); !slices.Equal(v.Suites.Names, []string{"behavior"}) || v.Severity != "" || v.Line != 22 {
		t.Errorf("verify: {suites: [behavior]}: %+v", v)
	}
	if w.Transitions[4].Verify != nil {
		t.Errorf("transitions[4] gives no verify, got %+v", *w.Transitions[4].Verify)
	}
	if w.States[0].Runs != nil {
		t.Errorf("backlog gives no runs, got %+v", *w.States[0].Runs)
	}
	if r := w.States[1].Runs; r == nil || r.Suites != nil || !slices.Equal(r.Kluster, KlusterSteps) || r.Line != 7 {
		t.Errorf("review's runs: %+v", r)
	}
	r := w.States[2].Runs
	if r == nil || r.Suites == nil || !slices.Equal(r.Suites.Names, []string{"integration", "smoke"}) || r.Suites.Line != 11 || !slices.Equal(r.Kluster, []string{KlusterDeploy}) || r.Line != 10 {
		t.Errorf("integration's runs: %+v", r)
	}
	at := w.VerifyAt
	if at.Batch == nil || !at.Batch.All || at.Batch.Line != 25 || at.Publish != nil || at.CI == nil || !slices.Equal(at.CI.Names, []string{"static", "behavior"}) || at.CI.Line != 27 {
		t.Errorf("verify_at: %+v", at)
	}
	if at.At(VerifyAtBatch) != at.Batch || at.At(VerifyAtPublish) != nil || at.At(VerifyAtCI) != at.CI || at.At("nightly") != nil {
		t.Error("VerifyAt.At does not give each point's suites")
	}
	if m := mustParse(t, minimal); m.VerifyAt != (VerifyAt{}) || m.Transitions[0].Verify != nil || m.States[0].Runs != nil {
		t.Errorf("a file giving none of them: %+v", m)
	}
}

// S-0741: a timer naming a state the file lacks, in its sum or in a point, is
// refused under workflow.timer, naming the timer, its part, and the state.
func TestParseRefusesATimerOfNoState(t *testing.T) {
	cases := []struct {
		name, content string
		field, says   string
	}{
		{"sum of no state", withTimer("{name: wait, title: Wait, sum: [ready]}"), "timers[0].sum[0]", `timer wait's sum names "ready", which is no state of the file`},
		{"from of no state", withTimer("{name: lead, title: Lead, from: {first: ready}, to: completed}"), "timers[0].from", `timer lead's from names "ready"`},
		{"to of no state", withTimer("{name: lead, title: Lead, from: created, to: {last: review}}"), "timers[0].to", `timer lead's to names "review"`},
		{"until of no state", withTimer("{name: wait, title: Wait, sum: [backlog], until: {first: ready}}"), "timers[0].until", `timer wait's until names "ready"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, p := Parse("design/workflows/board.md", tc.content)
			if p == nil {
				t.Fatal("parsed, want a problem")
			}
			if p.Code != CodeTimer || p.Line != 11 || p.Field != tc.field || !strings.Contains(p.Message, tc.says) || !strings.Contains(p.Message, "name one of backlog, done") {
				t.Errorf("got %s (code %s, line %d, field %q), want workflow.timer on line 11, field %q, saying %q", p, p.Code, p.Line, p.Field, tc.field, tc.says)
			}
		})
	}
}

// S-0737: the lookups find states by name or alias and by kind, and the
// transitions an item of a type may take, the baseline's own included.
func TestLookups(t *testing.T) {
	w := mustParse(t, strings.Replace(adrBaseline, "{name: review, kind: waiting, limit: true}", "{name: review, kind: waiting, limit: true, aliases: [checking]}", 1))
	if s, ok := w.State("checking"); !ok || s.Name != "review" {
		t.Errorf("State(checking): %+v, %v", s, ok)
	}
	if _, ok := w.State("testing"); ok {
		t.Error("State(testing) found a state the file lacks")
	}
	if w.Kind("in-progress") != KindActive || w.Kind("checking") != KindWaiting || w.Kind("testing") != "" {
		t.Errorf("Kind: %q %q %q", w.Kind("in-progress"), w.Kind("checking"), w.Kind("testing"))
	}
	var queues []string
	for _, s := range w.OfKind(KindQueue) {
		queues = append(queues, s.Name)
	}
	if !slices.Equal(queues, []string{"backlog", "ready"}) {
		t.Errorf("OfKind(queue): %v", queues)
	}
	if s, ok := w.Initial(); !ok || s.Name != "backlog" {
		t.Errorf("Initial: %+v, %v", s, ok)
	}
	tos := func(ts []Transition) []string {
		var out []string
		for _, x := range ts {
			out = append(out, x.To)
		}
		return out
	}
	if got := tos(w.From(TypeTask, "in-progress")); !slices.Equal(got, []string{"review", "done", "ready", "cancelled"}) {
		t.Errorf("From(task, in-progress): %v", got)
	}
	if got := tos(w.From(TypeStory, "in-progress")); !slices.Equal(got, []string{"review", "ready", "cancelled"}) {
		t.Errorf("From(story, in-progress): %v", got)
	}
	if got := tos(w.From(TypeStory, "checking")); !slices.Equal(got, []string{"done", "in-progress", "cancelled"}) {
		t.Errorf("From(story, checking): %v", got)
	}
	if got := w.From(TypeStory, "done"); len(got) != 0 {
		t.Errorf("From(story, done): %v", got)
	}
	if tr, ok := w.Find(TypeStory, "review", "done"); !ok || !tr.HasEffect(EffectAccept) || !tr.Allows(ByAcceptance) || tr.Allows(ByAnyone) {
		t.Errorf("Find(story, review, done): %+v, %v", tr, ok)
	}
	if tr, ok := w.Find(TypeTask, "checking", "done"); !ok || tr.HasEffect(EffectAccept) || !tr.Allows(ByAnyone) || !tr.Takes(TypeTask) || tr.Takes(TypeStory) {
		t.Errorf("Find(task, checking, done): %+v, %v", tr, ok)
	}
	if _, ok := w.Find(TypeStory, "in-progress", "done"); ok {
		t.Error("Find(story, in-progress, done) found a move only a task makes")
	}
	if tr, ok := w.WithEffect(TypeStory, EffectPull); !ok || !slices.Equal(tr.From, []string{"ready"}) || tr.To != "in-progress" {
		t.Errorf("WithEffect(story, pull): %+v, %v", tr, ok)
	}
	if tr, ok := w.WithEffect(TypeStory, EffectSendBack); !ok || tr.To != "in-progress" {
		t.Errorf("WithEffect(story, send-back): %+v, %v", tr, ok)
	}
	if _, ok := w.WithEffect(TypeTask, EffectAccept); ok {
		t.Error("WithEffect(task, accept) found an acceptance for a task")
	}
	if !w.States[0].Takes(TypeEpic) {
		t.Error("backlog does not take epics")
	}
}

// S-0737: who may make a move is the transition's by, and a refusal names
// the move, who makes it, and what to do.
func TestPermits(t *testing.T) {
	w := mustParse(t, strings.Replace(adrBaseline, "  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog, by: [operator, flai]}\n", 1))
	find := func(typ, from, to string) Transition {
		t.Helper()
		tr, ok := w.Find(typ, from, to)
		if !ok {
			t.Fatalf("no transition %s %s to %s", typ, from, to)
		}
		return tr
	}
	accept, cascade, pull, back := find(TypeStory, "review", "done"), find(TypeStory, "review", "cancelled"), find(TypeStory, "ready", "in-progress"), find(TypeStory, "ready", "backlog")
	cases := []struct {
		name string
		tr   Transition
		m    Mover
		says string
	}{
		{"an acceptance accepts", accept, Mover{Who: MoverOperator, Acceptance: true}, ""},
		{"the operator's move is no acceptance", accept, Mover{Who: MoverOperator}, "only an acceptance moves an epic or a story from review to done: accept it with flai accept"},
		{"an agent's move is no acceptance", accept, Mover{Who: MoverAgent}, "only an acceptance"},
		{"a cascade cancels from review", cascade, Mover{Who: MoverOperator, Cascade: true}, ""},
		{"a move of its own does not", cascade, Mover{Who: MoverOperator}, "the move from review to cancelled is made only as a parent's cancellation takes the item along"},
		{"anyone pulls, a role too", pull, Mover{Who: MoverRole}, ""},
		{"anyone pulls, an agent too", pull, Mover{Who: MoverAgent}, ""},
		{"the operator moves back", back, Mover{Who: MoverOperator}, ""},
		{"flai moves back", back, Mover{Who: MoverFlai}, ""},
		{"an agent does not", back, Mover{Who: MoverAgent}, "only the operator or flai serve may move an item from ready to backlog"},
		{"nor a role", back, Mover{Who: MoverRole}, "only the operator or flai serve"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tr.Permits(tc.m)
			switch {
			case tc.says == "" && err != nil:
				t.Errorf("refused: %v", err)
			case tc.says != "" && (err == nil || !strings.Contains(err.Error(), tc.says)):
				t.Errorf("got %v, want %q", err, tc.says)
			}
		})
	}
	multi := find(TypeStory, "backlog", "cancelled")
	multi.By = []string{ByAgent}
	if err := multi.Permits(Mover{Who: MoverOperator}); err == nil || !strings.Contains(err.Error(), "only the story's agent may move an item from backlog, ready, or in-progress to cancelled") {
		t.Errorf("a move from three states: %v", err)
	}
}
