package verifyrecord

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var ranAt = time.Date(2026, 11, 18, 9, 0, 0, 0, time.UTC)

// tierStep is a tier's step in a stored report.
func tierStep(name string, state State, coveredBy string) Step {
	return Step{Name: name, Tier: true, State: state, CoveredBy: coveredBy}
}

// suite is a clean suite at commit, run at ranAt, with tiers.
func suite(name, commit string, passed bool, tiers ...Tier) Suite {
	return Suite{Name: name, Passed: passed, Commit: commit, Clean: true, RanAt: ranAt, Tiers: tiers}
}

// writeRecord writes data as the story's verdict at commit under a new cache
// folder, which it answers.
func writeRecord(t *testing.T, story, commit, data string) string {
	t.Helper()
	cache := t.TempDir()
	path := CommitPath(cache, story, commit)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return cache
}

// Group groups a report's tier steps by suite, at the report's head and
// base: a suite fails on a tier that failed or was not reached, or that was
// skipped for one that did; a tier left to the merge queue does not fail it.
func TestGroupGroupsTheTierStepsBySuite(t *testing.T) {
	t.Parallel()
	rec := Record{Commit: "c1", Base: "main", BaseCommit: "b1", MergeBase: "m1", Patch: "p1", RanAt: ranAt, Clean: true, Steps: []Step{
		{Name: "check", State: Passed},
		tierStep("lint", Skipped, "unit"),
		tierStep("unit", Passed, ""),
		tierStep("vet", Skipped, "integration"),
		tierStep("integration", Failed, ""),
		tierStep("smoke", NotReached, ""),
		tierStep("e2e", Queued, ""),
	}}
	suites := map[string]string{"lint": "static", "unit": "static", "vet": "static", "integration": "behavior", "smoke": "behavior"}
	got := Group(rec, func(tier string) string { return suites[tier] })
	at := func(s Suite) Suite {
		s.Commit, s.Base, s.BaseCommit, s.MergeBase, s.Patch = "c1", "main", "b1", "m1", "p1"
		return s
	}
	want := []Suite{
		at(suite("static", "c1", false, Tier{Name: "lint", State: Skipped, CoveredBy: "unit"}, Tier{Name: "unit", State: Passed},
			Tier{Name: "vet", State: Skipped, CoveredBy: "integration"})),
		at(suite("behavior", "c1", false, Tier{Name: "integration", State: Failed}, Tier{Name: "smoke", State: NotReached})),
		at(suite("e2e", "c1", true, Tier{Name: "e2e", State: Queued})),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("suites\n%+v\nwant\n%+v", got, want)
	}
	rec.Steps = rec.Steps[:3]
	if got := Group(rec, nil); len(got) != 2 || got[0].Name != "lint" || !got[0].Passed || got[1].Name != "unit" {
		t.Errorf("each tier its own suite: %+v", got)
	}
}

// A report written before suites were recorded reads, each tier its own
// suite or the one suiteOf gives it, with no fingerprint; one that holds
// suites answers them.
func TestSuitesOfReadsAReportFromBeforeSuites(t *testing.T) {
	t.Parallel()
	old := `{"story":"S-004","commit":"c1","base":"main","base_commit":"b1","fingerprint":"f","ran_at":"2026-11-18T09:00:00Z",` +
		`"duration_ms":10,"passed":true,"clean":true,"steps":[{"name":"check","state":"passed","duration_ms":1},` +
		`{"name":"unit","tier":true,"state":"passed","duration_ms":1},{"name":"smoke","tier":true,"state":"reused","reused_from":"2026-11-17T09:00:00Z"}]}`
	path := filepath.Join(t.TempDir(), "S-004.json")
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	rec, ok, err := Read(path)
	if err != nil || !ok || rec.Commit != "c1" || len(rec.Steps) != 3 || rec.SuiteRuns != nil {
		t.Fatalf("read %+v, %v, %v", rec, ok, err)
	}
	got := rec.SuitesOf(nil)
	if len(got) != 2 || got[0].Name != "unit" || got[1].Name != "smoke" || got[0].Fingerprint != "" || !got[1].Passed ||
		got[1].Commit != "c1" || got[1].BaseCommit != "b1" || !got[1].RanAt.Equal(ranAt) || got[1].Tiers[0].ReusedFrom == nil {
		t.Errorf("own suites %+v", got)
	}
	if got := rec.SuitesOf(func(string) string { return "behavior" }); len(got) != 1 || !reflect.DeepEqual(got[0].TierNames(), []string{"unit", "smoke"}) {
		t.Errorf("one suite %+v", got)
	}
	rec.SuiteRuns = []Suite{suite("static", "c1", true)}
	if got := rec.SuitesOf(nil); !reflect.DeepEqual(got, rec.SuiteRuns) {
		t.Errorf("held suites %+v", got)
	}
}

// Merge keeps every suite another run recorded, puts add's in their place,
// and lets no entry that ran none of its tiers displace one already there.
func TestMergeKeepsTheSuitesARunDidNotRun(t *testing.T) {
	t.Parallel()
	staticOld := suite("static", "c1", false, Tier{Name: "unit", State: Failed})
	behaviorOld := suite("behavior", "c1", true, Tier{Name: "smoke", State: Passed})
	e2eOld := suite("e2e", "c1", true, Tier{Name: "e2e", State: Passed})
	staticNew := suite("static", "c1", true, Tier{Name: "unit", State: Passed})
	behaviorQueued := suite("behavior", "c1", true, Tier{Name: "smoke", State: Queued})
	e2eUnreached := suite("e2e", "c1", false, Tier{Name: "e2e", State: NotReached})
	docsNew := suite("docs", "c1", true, Tier{Name: "md", State: Passed})

	merged, kept := Merge([]Suite{staticOld, behaviorOld, e2eOld}, []Suite{docsNew, staticNew, behaviorQueued, e2eUnreached})
	if want := []Suite{staticNew, behaviorOld, e2eOld, docsNew}; !reflect.DeepEqual(merged, want) {
		t.Errorf("merged\n%+v\nwant\n%+v", merged, want)
	}
	if !reflect.DeepEqual(kept, []string{"behavior", "e2e"}) {
		t.Errorf("kept %v", kept)
	}
	if merged, kept := Merge(nil, []Suite{behaviorQueued}); !reflect.DeepEqual(merged, []Suite{behaviorQueued}) || kept != nil {
		t.Errorf("into nothing %+v, kept %v", merged, kept)
	}
}

// AllPassed says every named suite passed at the commit.
func TestAllPassed(t *testing.T) {
	t.Parallel()
	suites := []Suite{suite("static", "c1", true), suite("behavior", "c1", false), suite("smoke", "c0", true)}
	for _, tc := range []struct {
		names []string
		want  bool
	}{
		{nil, true},
		{[]string{"static"}, true},
		{[]string{"static", "behavior"}, false},
		{[]string{"smoke"}, false},
		{[]string{"docs"}, false},
	} {
		if got := AllPassed(suites, "c1", tc.names); got != tc.want {
			t.Errorf("AllPassed(%v) = %v, want %v", tc.names, got, tc.want)
		}
	}
}

// AtHead answers per suite whether it passed at the head in a clean run,
// with which of its tiers held, following a skipped tier to the one that
// covered it in another suite.
func TestAtHeadAnswersEachSuitesVerdictAtTheHead(t *testing.T) {
	t.Parallel()
	data := `{"story":"S-004","commit":"c1","base":"main","ran_at":"2026-11-18T09:00:00Z","passed":false,"steps":[],"suite_runs":[
		{"name":"static","passed":true,"commit":"c1","clean":true,"ran_at":"2026-11-18T09:00:00Z","tiers":[{"name":"lint","state":"skipped","covered_by":"integration"},{"name":"unit","state":"reused"}]},
		{"name":"behavior","passed":true,"commit":"c1","clean":true,"ran_at":"2026-11-18T09:00:00Z","tiers":[{"name":"integration","state":"passed"},{"name":"smoke","state":"queued"}]},
		{"name":"docs","passed":false,"commit":"c1","clean":true,"ran_at":"2026-11-18T09:00:00Z","tiers":[{"name":"md","state":"failed"}]},
		{"name":"dirty","passed":true,"commit":"c1","ran_at":"2026-11-18T09:00:00Z","tiers":[{"name":"x","state":"passed"}]},
		{"name":"elsewhere","passed":true,"commit":"c0","clean":true,"ran_at":"2026-11-18T09:00:00Z","tiers":[{"name":"y","state":"passed"}]}]}`
	cache := writeRecord(t, "S-004", "c1", data)
	got, ok, err := AtHead(cache, "S-004", "c1", nil)
	want := []Verdict{
		{Suite: "static", Passed: true, AtHead: true, Tiers: []string{"lint", "unit"}, Held: []string{"lint", "unit"}},
		{Suite: "behavior", AtHead: true, Tiers: []string{"integration", "smoke"}, Held: []string{"integration"}, Queued: []string{"smoke"}},
		{Suite: "docs", AtHead: true, Tiers: []string{"md"}},
		{Suite: "dirty", Tiers: []string{"x"}, Held: []string{"x"}},
		{Suite: "elsewhere", Tiers: []string{"y"}, Held: []string{"y"}},
	}
	if err != nil || !ok || !reflect.DeepEqual(got, want) {
		t.Errorf("verdicts %+v, %v, %v\nwant %+v", got, ok, err, want)
	}
	if got, ok, err := AtHead(cache, "S-004", "c2", nil); ok || err != nil || got != nil {
		t.Errorf("no record at the head: %+v, %v, %v", got, ok, err)
	}
	broken := writeRecord(t, "S-004", "c1", "{")
	if _, ok, err := AtHead(broken, "S-004", "c1", nil); ok || err == nil || !strings.Contains(err.Error(), "delete it and verify again") {
		t.Errorf("a record that does not read: %v, %v", ok, err)
	}
}

// The record's files are where verify stores them in the project's cache.
func TestPaths(t *testing.T) {
	t.Parallel()
	if got, want := ReportPath("/p/.flai-cache", "S-004"), filepath.Join("/p/.flai-cache", "verify", "S-004.json"); got != want {
		t.Errorf("report path %s, want %s", got, want)
	}
	if got, want := CommitPath("/p/.flai-cache", "S-004", "c1"), filepath.Join("/p/.flai-cache", "verify", "S-004", "c1.json"); got != want {
		t.Errorf("commit path %s, want %s", got, want)
	}
}
