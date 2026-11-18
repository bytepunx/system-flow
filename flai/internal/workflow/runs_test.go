package workflow

import (
	"slices"
	"strings"
	"testing"
)

// integrationFlow is a board with integration between review and done, the
// state the acceptance leaves from, which runs the integration suite and the
// cluster's deploy alone and sends a story back to in-progress; review gives
// no runs and so runs nothing.
const integrationFlow = `---
title: Board
updated: 2026-11-18
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: integration, kind: waiting, aliases: [int], runs: {suites: [integration], kluster: [deploy]}}
  - {name: done, kind: accepted}
transitions:
  - {from: backlog, to: in-progress, effects: [pull]}
  - {from: in-progress, to: review}
  - {from: review, to: integration}
  - {from: integration, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: [in-progress, review, integration], to: done, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: integration, to: in-progress, effects: [send-back]}
---
`

// S-0810: on the baseline, review is the one state that runs anything, the
// cluster's deploy and checks and no suite.
func TestStateRunsOnBaseline(t *testing.T) {
	w := Baseline()
	if got := w.KlusterStates(); !slices.Equal(got, []string{"review"}) {
		t.Errorf("KlusterStates = %v, want [review]", got)
	}
	r := w.StateRuns("review")
	if r == nil {
		t.Fatal("review runs nothing")
	}
	if !r.RunsStep(KlusterDeploy) || !r.RunsStep(KlusterChecks) {
		t.Errorf("review's runs %+v lack deploy or checks", r)
	}
	if names, all := r.SuiteNames(); names != nil || all {
		t.Errorf("review runs suites %v, %v; want none", names, all)
	}
	if r.Line != 9 {
		t.Errorf("review's runs are on line %d, want 9", r.Line)
	}
	for _, s := range []string{"backlog", "ready", "in-progress", "done", "cancelled", "nowhere"} {
		if r := w.StateRuns(s); r != nil {
			t.Errorf("%s runs %+v, want nothing", s, r)
		}
	}
}

// S-0810: a state runs what its runs name, suites and kluster steps each on
// their own; a state with no runs in a file that gives the keys runs
// nothing, and an alias names its state.
func TestStateRunsOfAFileGivingThem(t *testing.T) {
	w := mustParse(t, integrationFlow)
	if got := w.KlusterStates(); !slices.Equal(got, []string{"integration"}) {
		t.Errorf("KlusterStates = %v, want [integration]", got)
	}
	for _, name := range []string{"integration", "int"} {
		r := w.StateRuns(name)
		if r == nil {
			t.Fatalf("%s runs nothing", name)
		}
		if names, all := r.SuiteNames(); !slices.Equal(names, []string{"integration"}) || all {
			t.Errorf("%s runs suites %v, %v; want [integration], false", name, names, all)
		}
		if !r.RunsStep(KlusterDeploy) || r.RunsStep(KlusterChecks) {
			t.Errorf("%s's runs %+v: want deploy and not checks", name, r)
		}
	}
	for _, s := range []string{"backlog", "in-progress", "review", "done"} {
		if r := w.StateRuns(s); r != nil {
			t.Errorf("%s runs %+v, want nothing", s, r)
		}
	}
	every := mustParse(t, strings.Replace(integrationFlow, "{suites: [integration], kluster: [deploy]}", "{suites: all}", 1))
	r := every.StateRuns("integration")
	if names, all := r.SuiteNames(); names != nil || !all {
		t.Errorf("suites: all gives %v, %v; want none, true", names, all)
	}
	if r.RunsStep(KlusterDeploy) || every.KlusterStates() != nil {
		t.Errorf("a state running suites alone names a cluster step: %+v", r)
	}
}

// S-0810 (TH-1289): a file that gives none of ADR-0312's keys runs the
// baseline's deploy and checks on each state a story's acceptance leaves
// from, whatever it is called, and nothing elsewhere.
func TestStateRunsOfAFileGivingNoKeys(t *testing.T) {
	for _, c := range []struct {
		name, text, waiting string
	}{
		{"ADR-0301's baseline", adrBaseline, "review"},
		{"renamed", renamedFlow, "awaiting"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := mustParse(t, c.text)
			if got := w.KlusterStates(); !slices.Equal(got, []string{c.waiting}) {
				t.Errorf("KlusterStates = %v, want [%s]", got, c.waiting)
			}
			r := w.StateRuns(c.waiting)
			if r == nil || !slices.Equal(r.Kluster, KlusterSteps) || r.Suites != nil || r.Line != 0 {
				t.Errorf("%s runs %+v, want the baseline's kluster [deploy, checks] with no line", c.waiting, r)
			}
			for _, s := range w.Names() {
				if s != c.waiting && w.StateRuns(s) != nil {
					t.Errorf("%s runs something", s)
				}
			}
		})
	}
}

// S-0810: a nil Runs runs no suite and names no step, and StateRuns gives a
// copy, so a change to it leaves the workflow as it was.
func TestRunsNilAndCopy(t *testing.T) {
	var none *Runs
	if names, all := none.SuiteNames(); names != nil || all || none.RunsStep(KlusterDeploy) {
		t.Error("a nil Runs runs something")
	}
	w := mustParse(t, integrationFlow)
	r := w.StateRuns("integration")
	r.Suites.Names[0], r.Kluster[0] = "changed", "changed"
	names, _ := r.SuiteNames()
	names[0] = "again"
	again := w.StateRuns("integration")
	if again.Suites.Names[0] != "integration" || again.Kluster[0] != KlusterDeploy {
		t.Errorf("a change to StateRuns' answer changed the workflow: %+v", again)
	}
	b := Baseline().StateRuns("review")
	b.Kluster[0] = "changed"
	if got := mustParse(t, adrBaseline).StateRuns("review"); got.Kluster[0] != KlusterDeploy {
		t.Errorf("a change to the baseline's runs reached a file giving no keys: %+v", got)
	}
}
