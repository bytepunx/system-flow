package workflow

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// S-0806: NamedSuites lists every suite the file names, where, and on which
// line, transitions first, then states, then verify_at's points, and whether
// any place names all.
func TestNamedSuites(t *testing.T) {
	named, all := mustParse(t, suited).NamedSuites()
	want := []Named{
		{Suite: "static", Where: WhereTransition, Name: "review → integration", Line: 16},
		{Suite: "behavior", Where: WhereTransition, Name: "review → integration", Line: 16},
		{Suite: "behavior", Where: WhereTransition, Name: "review → done", Line: 22},
		{Suite: "integration", Where: WhereState, Name: "integration", Line: 11},
		{Suite: "smoke", Where: WhereState, Name: "integration", Line: 11},
		{Suite: "static", Where: WherePoint, Name: VerifyAtCI, Line: 27},
		{Suite: "behavior", Where: WherePoint, Name: VerifyAtCI, Line: 27},
	}
	if !slices.Equal(named, want) || !all {
		t.Errorf("NamedSuites: %+v, %v; want %+v, true", named, all, want)
	}
}

// S-0806: a file that names all alone lists no suite and says so; one that
// names some lists them without all; one that names none lists nothing.
func TestNamedSuitesOfAllSomeAndNone(t *testing.T) {
	named, all := mustParse(t, strings.Replace(minimal, "to: done}", "to: done, verify: all}", 1)).NamedSuites()
	if len(named) != 0 || !all {
		t.Errorf("verify: all: %+v, %v; want none, true", named, all)
	}
	some := strings.Replace(minimal, "  - {from: backlog, to: done}\n", "  - {from: backlog, to: done}\nverify_at: {publish: [smoke]}\n", 1)
	named, all = mustParse(t, some).NamedSuites()
	if want := []Named{{Suite: "smoke", Where: WherePoint, Name: VerifyAtPublish, Line: 10}}; !slices.Equal(named, want) || all {
		t.Errorf("verify_at: {publish: [smoke]}: %+v, %v; want %+v, false", named, all, want)
	}
	named, all = mustParse(t, minimal).NamedSuites()
	if named != nil || all {
		t.Errorf("a file naming none: %+v, %v", named, all)
	}
}

// S-0806 (TH-1289): a file uses ADR-0312's keys when any transition gives
// verify, any state runs, or the file verify_at, each alone; one that gives
// none of them does not, and reads as the baseline's.
func TestUsesSuiteKeys(t *testing.T) {
	if mustParse(t, minimal).UsesSuiteKeys() {
		t.Error("a file giving none of verify, runs, and verify_at uses them")
	}
	for name, text := range map[string]string{
		"verify":    strings.Replace(minimal, "to: done}", "to: done, verify: [static]}", 1),
		"verify_at": strings.Replace(minimal, "  - {from: backlog, to: done}\n", "  - {from: backlog, to: done}\nverify_at: {publish: [smoke]}\n", 1),
	} {
		if !mustParse(t, text).UsesSuiteKeys() {
			t.Errorf("a file giving %s alone does not use the keys", name)
		}
	}
	w := mustParse(t, minimal)
	w.States[0].Runs = &Runs{Kluster: []string{KlusterDeploy}}
	if !w.UsesSuiteKeys() {
		t.Error("a file giving a state's runs alone does not use the keys")
	}
	if !mustParse(t, suited).UsesSuiteKeys() {
		t.Error("the suited file does not use the keys")
	}
}

// S-0806: Clone copies a transition's verify, a state's runs, and verify_at,
// so a change to the copy leaves the original as it was.
func TestCloneCopiesVerifyRunsAndVerifyAt(t *testing.T) {
	w := mustParse(t, suited)
	c := w.Clone()
	c.Transitions[1].Verify.Suites.Names[0] = "changed"
	c.Transitions[1].Verify.Severity = SeverityWarn
	c.States[2].Runs.Suites.Names[0] = "changed"
	c.States[2].Runs.Kluster[0] = "changed"
	c.States[1].Runs.Line = 0
	c.VerifyAt.CI.Names[0] = "changed"
	c.VerifyAt.Batch.All = false
	if v := w.Transitions[1].Verify; v.Suites.Names[0] != "static" || v.Severity != "" {
		t.Errorf("the original's verify changed: %+v", v)
	}
	if r := w.States[2].Runs; r.Suites.Names[0] != "integration" || r.Kluster[0] != KlusterDeploy {
		t.Errorf("the original's runs changed: %+v", r)
	}
	if w.States[1].Runs.Line != 7 {
		t.Error("the original's runs is the copy's")
	}
	if w.VerifyAt.CI.Names[0] != "static" || !w.VerifyAt.Batch.All {
		t.Errorf("the original's verify_at changed: %+v", w.VerifyAt)
	}
	if m := mustParse(t, minimal).Clone(); m.VerifyAt != (VerifyAt{}) || m.Transitions[0].Verify != nil || m.States[0].Runs != nil {
		t.Errorf("a clone of a file giving none of them: %+v", m)
	}
}

// S-0806: a file that gives none of the keys reads as JSON as it did before
// them; one that does gives each where it is.
func TestJSONOfVerifyRunsAndVerifyAt(t *testing.T) {
	plain, err := json.Marshal(mustParse(t, minimal))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"verify_at"`, `"verify"`, `"runs"`} {
		if strings.Contains(string(plain), key) {
			t.Errorf("a file giving none of them has %s in its JSON: %s", key, plain)
		}
	}
	got, err := json.Marshal(mustParse(t, suited))
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{
		`"verify":{"suites":{"all":true,"line":15},"line":15}`,
		`"verify":{"suites":{"all":true,"line":20},"severity":"warn","line":19}`,
		`"runs":{"kluster":["deploy","checks"],"line":7}`,
		`"runs":{"suites":{"names":["integration","smoke"],"line":11},"kluster":["deploy"],"line":10}`,
		`"verify_at":{"batch":{"all":true,"line":25},"ci":{"names":["static","behavior"],"line":27}}`,
	} {
		if !strings.Contains(string(got), part) {
			t.Errorf("JSON lacks %s: %s", part, got)
		}
	}
}
