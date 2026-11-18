package verify

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
)

// verifyAgain runs Verify a minute after the fixture's last run, with what
// the tiers ran reset first, failing the test on an error.
func (f *verifyFixture) verifyAgain(t *testing.T) Report {
	t.Helper()
	f.proc.ran, f.proc.now = nil, f.proc.now.Add(time.Minute)
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

// ADR-0229: a close-out asked again at a clean commit the record holds a
// passing verdict for answers it, saying so and when that run was, runs no
// tier and no step, and leaves the story's last report as it was.
func TestVerifyAnswersThePassingVerdictTheRecordHoldsAndRunsNothing(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	first, err := Verify(context.Background(), f.opts)
	if err != nil || !first.Passed || !first.Clean || first.FromRecord {
		t.Fatalf("first run %+v, %v", first, err)
	}
	if _, err := os.Stat(RecordPath(f.opts.Project, "S-004", "0123abc")); err != nil {
		t.Fatalf("the clean run is not kept in the record: %v", err)
	}
	// a step that ran would find the narrative empty now
	f.editNarrative(t, "## Current state\n", "## Current state\n\n## Unrelated\n")
	f.editNarrative(t, "## Next steps\n", "## Next steps\n\n## Unrelated too\n")
	second := f.verifyAgain(t)
	if !second.FromRecord || !second.Passed || !second.RanAt.Equal(first.RanAt) || len(f.proc.ran) != 0 {
		t.Fatalf("second ask %+v ran %q; want the first run's verdict from the record, running nothing", second, f.proc.ran)
	}
	second.FromRecord = false
	if !reflect.DeepEqual(second, first) {
		t.Errorf("the record's verdict %+v\nwant the first run's %+v", second, first)
	}
	if got, want := (Report{FromRecord: true, RanAt: first.RanAt, Commit: "0123abc"}).RecordText(),
		"answered from the record of 2026-09-01T12:00:00Z at 0123abc; ran nothing"; got != want {
		t.Errorf("record text %q, want %q", got, want)
	}
	if first.RecordText() != "" {
		t.Errorf("a run says it came from the record: %q", first.RecordText())
	}
	if last, ok, err := LastReport(f.opts.Project, "S-004"); err != nil || !ok || last.FromRecord || !reflect.DeepEqual(last, first) {
		t.Errorf("the last report %+v, %v, %v; want the first run's, untouched", last, ok, err)
	}
	rep, ok, err := Recorded(f.opts)
	if err != nil || !ok || !rep.FromRecord || !rep.RanAt.Equal(first.RanAt) || len(f.proc.ran) != 0 {
		t.Errorf("Recorded %+v, %v, %v", rep, ok, err)
	}
}

// ADR-0229: the record answers only for the commit, the clean worktree, the
// tiers, and the main branch its verdict covered, and never when asked for a
// fresh run; main ahead only by commits the sync step passes over is covered.
func TestVerifyAnswersFromTheRecordOnlyWhatItsVerdictCovers(t *testing.T) {
	t.Parallel()
	moveMain := func(paths string) func(*testing.T, *verifyFixture) {
		return func(_ *testing.T, f *verifyFixture) {
			f.git.out["rev-parse --verify --quiet main^{commit}"] = "9abcdef\n"
			f.git.out["diff --no-renames --name-only HEAD...main"] = paths
			f.git.out["rev-list --abbrev-commit HEAD..main"] = "a1a1a1a"
			f.git.out["merge-base --is-ancestor fedcba9 9abcdef"] = ""
		}
	}
	for _, tc := range []struct {
		name   string
		first  func(*testing.T, *verifyFixture)
		change func(*testing.T, *verifyFixture)
		record bool
	}{
		{name: "the same clean commit", change: func(*testing.T, *verifyFixture) {}, record: true},
		{name: "main ahead by wip commits the branch does not change", change: moveMain("wip/kanban/stories/S-009-nine.md"), record: true},
		{name: "a new commit", change: func(_ *testing.T, f *verifyFixture) { f.git.out["rev-parse HEAD"] = "4567def\n" }},
		{name: "fresh", change: func(_ *testing.T, f *verifyFixture) { f.opts.Fresh = true }},
		{name: "an uncommitted change", change: func(_ *testing.T, f *verifyFixture) {
			f.git.out["status --porcelain --untracked-files=all"] = " M flai/internal/x/x.go"
		}},
		{name: "an untracked file", change: func(_ *testing.T, f *verifyFixture) {
			f.git.out["status --porcelain --untracked-files=all"] = "?? scratch.txt"
		}},
		{name: "main ahead by a commit the sync step does not pass over", change: moveMain("flai/internal/y/y.go")},
		{name: "main not ahead of the record's base", change: func(t *testing.T, f *verifyFixture) {
			moveMain("wip/kanban/stories/S-009-nine.md")(t, f)
			delete(f.git.out, "merge-base --is-ancestor fedcba9 9abcdef")
		}},
		{name: "a tier's new command", change: func(_ *testing.T, f *verifyFixture) {
			f.opts.Tiers = append([]Tier{}, storyTiers...)
			f.opts.Tiers[4].Command = []string{"smoke.sh", "--all"}
		}},
		{name: "every tier asked for after some were left to the queue",
			first: func(_ *testing.T, f *verifyFixture) {
				f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), true
			},
			change: func(_ *testing.T, f *verifyFixture) { f.opts.LeaveToQueue = false }},
		{name: "a verdict that did not pass",
			first:  func(_ *testing.T, f *verifyFixture) { f.proc.steps["smoke.sh"] = step{exit: 1} },
			change: func(_ *testing.T, f *verifyFixture) { f.proc.steps["smoke.sh"] = step{} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newVerifyFixture(t)
			if tc.first != nil {
				tc.first(t, f)
			}
			first, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			tc.change(t, f)
			rep := f.verifyAgain(t)
			if rep.FromRecord != tc.record {
				t.Fatalf("from the record %v, want %v: %+v", rep.FromRecord, tc.record, rep)
			}
			if tc.record {
				if !rep.RanAt.Equal(first.RanAt) || len(f.proc.ran) != 0 {
					t.Errorf("answered %+v, ran %q", rep, f.proc.ran)
				}
				return
			}
			if rep.RanAt.Equal(first.RanAt) || len(rep.Steps) < 4 || rep.Steps[0].Name != StepRebase {
				t.Errorf("did not run again: %+v", rep)
			}
			if last, _, err := LastReport(f.opts.Project, "S-004"); err != nil || !last.RanAt.Equal(rep.RanAt) {
				t.Errorf("the last report is not the new run's: %+v, %v", last, err)
			}
		})
	}
}

// ADR-0229: a run that began with something uncommitted is the story's last
// report but is kept in no record, since its commit is not what it verified.
func TestVerifyKeepsNoRecordOfARunWithSomethingUncommitted(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	f.git.out["status --porcelain --untracked-files=all"] = " M flai/internal/x/x.go"
	rep, err := Verify(context.Background(), f.opts)
	if err != nil || !rep.Passed || rep.Clean {
		t.Fatalf("report %+v, %v", rep, err)
	}
	if _, err := os.Stat(RecordPath(f.opts.Project, "S-004", "0123abc")); !os.IsNotExist(err) {
		t.Errorf("a run with something uncommitted is in the record: %v", err)
	}
}

// ADR-0229: a close-out that committed what a passing run verified keeps the
// verdict at its commit, and a later ask there answers from the record.
func TestKeepAtLetsALaterAskAtThatCommitAnswerFromTheRecord(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	f.git.out["status --porcelain --untracked-files=all"] = " M flai/internal/x/x.go"
	ran, err := Verify(context.Background(), f.opts)
	if err != nil || !ran.Passed || ran.Clean {
		t.Fatalf("the run before the commit %+v, %v", ran, err)
	}
	f.git.out["status --porcelain --untracked-files=all"] = ""
	f.git.out["rev-parse HEAD"] = "4567def\n"
	if _, ok, err := Recorded(f.opts); ok || err != nil {
		t.Fatalf("the record holds the commit before KeepAt: %v, %v", ok, err)
	}
	if err := KeepAt(f.opts.Project, "4567def", ran); err != nil {
		t.Fatal(err)
	}
	rep := f.verifyAgain(t)
	if !rep.FromRecord || !rep.Passed || rep.Commit != "4567def" || !rep.RanAt.Equal(ran.RanAt) || len(f.proc.ran) != 0 {
		t.Errorf("the ask at the close-out's commit %+v ran %q; want the kept verdict, running nothing", rep, f.proc.ran)
	}
	for _, s := range rep.SuiteRuns {
		if s.Commit != "4567def" || !s.Clean {
			t.Errorf("suite %s kept at %s, clean %v; want it verified at exactly the close-out's commit", s.Name, s.Commit, s.Clean)
		}
	}
	// a last report keeps an earlier run's suite at another commit, which
	// the record at the close-out's commit leaves out
	earlier := verifyrecord.Suite{Name: "docs", Passed: true, Commit: "0000aaa", RanAt: ran.RanAt.Add(-time.Hour), Tiers: []verifyrecord.Tier{{Name: "md", State: verifyrecord.Passed}}}
	ran.SuiteRuns = append(ran.SuiteRuns, earlier)
	if err := KeepAt(f.opts.Project, "89abcde", ran); err != nil {
		t.Fatal(err)
	}
	if kept, _, err := readRecord(f.opts.Project, "S-004", "89abcde"); err != nil || len(kept.SuiteRuns) != len(ran.SuiteRuns)-1 ||
		slices.ContainsFunc(kept.SuiteRuns, func(s verifyrecord.Suite) bool { return s.Name == "docs" }) {
		t.Errorf("the record at the close-out's commit holds %+v, %v; want the run's own suites alone", kept.SuiteRuns, err)
	}
	if err := KeepAt(f.opts.Project, "", ran); err == nil {
		t.Error("KeepAt kept a verdict at no commit")
	}
}

// ADR-0229: a record that does not read is an error to Recorded, and Verify
// runs in its place and keeps the run there.
func TestVerifyRunsInPlaceOfARecordThatDoesNotRead(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	if _, ok, err := Recorded(f.opts); ok || err != nil {
		t.Fatalf("a story never verified has a record: %v, %v", ok, err)
	}
	if _, err := Verify(context.Background(), f.opts); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(RecordPath(f.opts.Project, "S-004", "0123abc"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := Recorded(f.opts); ok || err == nil || !strings.Contains(err.Error(), "delete it and verify again") {
		t.Errorf("a record that does not read: %v, %v", ok, err)
	}
	rep := f.verifyAgain(t)
	if rep.FromRecord || !rep.Passed {
		t.Fatalf("report %+v", rep)
	}
	if again, ok, err := Recorded(f.opts); err != nil || !ok || !again.RanAt.Equal(rep.RanAt) {
		t.Errorf("the run is not kept in the record: %+v, %v, %v", again, ok, err)
	}
}

// ADR-0312 §7: two runs of different suites at one commit both stay in the
// record and in the story's last report, and a later ask for one suite's
// tiers answers from the record with that suite alone, running nothing.
func TestTwoRunsOfDifferentSuitesAtOneCommitBothStayInTheRecord(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	static := []Tier{{Name: "go-test", Suite: "static", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}}}
	behavior := []Tier{{Name: "smoke", Suite: "behavior", AllOnly: true, Command: []string{"smoke.sh"}}}
	f.opts.Tiers = static
	first, err := Verify(context.Background(), f.opts)
	if err != nil || !first.Passed || len(first.SuiteRuns) != 1 {
		t.Fatalf("first run %+v, %v", first, err)
	}
	f.opts.Tiers = behavior
	second := f.verifyAgain(t)
	if second.FromRecord || !second.Passed || !reflect.DeepEqual(f.proc.ran, []string{"smoke.sh"}) {
		t.Fatalf("second run %+v ran %q; want behavior run, not answered from static's record", second, f.proc.ran)
	}
	verdicts, ok, err := verifyrecord.AtHead(f.opts.Project.CacheDir(), "S-004", "0123abc", nil)
	want := []verifyrecord.Verdict{
		{Suite: "static", Passed: true, AtHead: true, Tiers: []string{"go-test"}, Held: []string{"go-test"}},
		{Suite: "behavior", Passed: true, AtHead: true, Tiers: []string{"smoke"}, Held: []string{"smoke"}},
	}
	if err != nil || !ok || !reflect.DeepEqual(verdicts, want) {
		t.Errorf("the record at the head %+v, %v, %v\nwant %+v", verdicts, ok, err, want)
	}
	last, _, err := LastReport(f.opts.Project, "S-004")
	if err != nil || !last.PassedSuites("0123abc", []string{"static", "behavior"}) || last.PassedSuites("4567def", []string{"static"}) {
		t.Errorf("the last report's suites %+v, %v; want both passed at 0123abc alone", last.SuiteRuns, err)
	}
	f.opts.Tiers = static
	third := f.verifyAgain(t)
	if !third.FromRecord || !third.Passed || len(f.proc.ran) != 0 {
		t.Fatalf("third ask %+v ran %q; want static answered from the record", third, f.proc.ran)
	}
	if got := tierStates(third); !reflect.DeepEqual(got, []string{"go-test:passed"}) || len(third.SuiteRuns) != 1 || third.SuiteRuns[0].Name != "static" {
		t.Errorf("the answer's tiers %v and suites %+v; want static's alone", got, third.SuiteRuns)
	}
	if !third.PassedSuites("0123abc", []string{"static"}) || third.PassedSuites("0123abc", []string{"behavior"}) {
		t.Errorf("the answer covers %+v; want static and not behavior", third.SuiteRuns)
	}
}

// ADR-0312 §7: a verdict kept at a commit over one there before keeps the
// suites the earlier recorded that it did not run, with their tier steps,
// fails while one of them failed, and an entry that ran none of its tiers
// leaves the earlier one in place.
func TestKeepAtKeepsTheSuitesAnEarlierVerdictAtTheCommitRecorded(t *testing.T) {
	t.Parallel()
	project := newVerifyFixture(t).opts.Project
	tier := func(name string, state State) Step { return Step{Name: name, Tier: true, State: state} }
	suite := func(name string, passed bool, tier string, state verifyrecord.State) verifyrecord.Suite {
		return verifyrecord.Suite{Name: name, Passed: passed, Commit: "0123abc", Tiers: []verifyrecord.Tier{{Name: tier, State: state}}}
	}
	check := Step{Name: StepCheck, State: Passed}
	keep := func(rep Report) Report {
		t.Helper()
		rep.Story, rep.Commit = "S-004", "0123abc"
		if err := KeepAt(project, "0123abc", rep); err != nil {
			t.Fatal(err)
		}
		rec, ok, err := readRecord(project, "S-004", "0123abc")
		if err != nil || !ok {
			t.Fatalf("the record %v, %v", ok, err)
		}
		return rec
	}
	names := func(rep Report) []string {
		var out []string
		for _, s := range rep.SuiteRuns {
			out = append(out, s.Name)
		}
		return out
	}
	for _, step := range []struct {
		name      string
		rep       Report
		passed    bool
		stoppedAt string
		steps     []string
	}{
		{"behavior failed", Report{StoppedAt: "smoke", Steps: []Step{check, tier("smoke", Failed)},
			SuiteRuns: []verifyrecord.Suite{suite("behavior", false, "smoke", verifyrecord.Failed)}},
			false, "smoke", []string{"check:passed", "smoke:failed"}},
		{"static passed beside it", Report{Passed: true, Steps: []Step{check, tier("go-test", Passed)},
			SuiteRuns: []verifyrecord.Suite{suite("static", true, "go-test", verifyrecord.Passed)}},
			false, "smoke", []string{"check:passed", "go-test:passed", "smoke:failed"}},
		{"behavior passed again", Report{Passed: true, Steps: []Step{check, tier("smoke", Passed)},
			SuiteRuns: []verifyrecord.Suite{suite("behavior", true, "smoke", verifyrecord.Passed)}},
			true, "", []string{"check:passed", "smoke:passed", "go-test:passed"}},
		{"behavior left to the queue", Report{Passed: true, Steps: []Step{check, tier("smoke", Queued)},
			SuiteRuns: []verifyrecord.Suite{suite("behavior", true, "smoke", verifyrecord.Queued)}},
			true, "", []string{"check:passed", "smoke:passed", "go-test:passed"}},
	} {
		rec := keep(step.rep)
		if rec.Passed != step.passed || rec.StoppedAt != step.stoppedAt || !reflect.DeepEqual(states(rec), step.steps) {
			t.Errorf("%s: passed %v at %q, steps %v; want %v at %q, %v", step.name, rec.Passed, rec.StoppedAt, states(rec), step.passed, step.stoppedAt, step.steps)
		}
		if got := names(rec); len(rec.SuiteRuns) > 1 && !reflect.DeepEqual(got, []string{"behavior", "static"}) {
			t.Errorf("%s: suites %v, want behavior then static", step.name, got)
		}
	}
}

// A record kept before suites were recorded still reads, each tier its own
// suite, but covers no suite for a run, which runs in its place and keeps
// itself with its suites.
func TestVerifyRunsInPlaceOfARecordFromBeforeSuites(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	path := RecordPath(f.opts.Project, "S-004", "0123abc")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"story":"S-004","commit":"0123abc","base":"main","base_commit":"fedcba9","fingerprint":"f","ran_at":"2026-08-01T12:00:00Z",` +
		`"duration_ms":10,"passed":true,"paths":["flai/internal/x/x.go"],"clean":true,"steps":[{"name":"check","state":"passed","duration_ms":1},` +
		`{"name":"go-test","tier":true,"state":"passed","duration_ms":1},{"name":"smoke","tier":true,"state":"passed","duration_ms":1}]}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	verdicts, ok, err := verifyrecord.AtHead(f.opts.Project.CacheDir(), "S-004", "0123abc", nil)
	want := []verifyrecord.Verdict{
		{Suite: "go-test", Passed: true, AtHead: true, Tiers: []string{"go-test"}, Held: []string{"go-test"}},
		{Suite: "smoke", Passed: true, AtHead: true, Tiers: []string{"smoke"}, Held: []string{"smoke"}},
	}
	if err != nil || !ok || !reflect.DeepEqual(verdicts, want) {
		t.Errorf("the old record read %+v, %v, %v\nwant %+v", verdicts, ok, err, want)
	}
	if _, ok, err := Recorded(f.opts); ok || err != nil {
		t.Fatalf("the old record answered: %v, %v", ok, err)
	}
	rep, err := Verify(context.Background(), f.opts)
	if err != nil || rep.FromRecord || !rep.Passed || len(f.proc.ran) != 3 {
		t.Fatalf("report %+v ran %q, %v; want every tier run", rep, f.proc.ran, err)
	}
	if again, ok, err := Recorded(f.opts); err != nil || !ok || len(again.SuiteRuns) != 3 {
		t.Errorf("the run is not kept with its suites: %+v, %v, %v", again.SuiteRuns, ok, err)
	}
}

// The record has one shape: verifyrecord reads the report verify stores as
// verify gives it, and names the tier states in verify's words.
func TestTheRecordReadsTheReportVerifyStores(t *testing.T) {
	t.Parallel()
	f := newVerifyFixture(t)
	f.opts.Tiers = queueTiers()
	f.opts.Tiers[2].Covers, f.opts.LeaveToQueue = []string{"go-test"}, false
	rep, err := Verify(context.Background(), f.opts)
	if err != nil || !rep.Passed {
		t.Fatalf("report %+v, %v", rep, err)
	}
	rec, ok, err := verifyrecord.Read(ReportPath(f.opts.Project, "S-004"))
	if err != nil || !ok || !reflect.DeepEqual(rec, rep.record()) {
		t.Errorf("verifyrecord read %+v, %v, %v\nwant %+v", rec, ok, err, rep.record())
	}
	for state, word := range map[State]verifyrecord.State{Passed: verifyrecord.Passed, Failed: verifyrecord.Failed,
		NotReached: verifyrecord.NotReached, Reused: verifyrecord.Reused, Skipped: verifyrecord.Skipped, Queued: verifyrecord.Queued} {
		if string(state) != string(word) {
			t.Errorf("verify's %q is %q to the record", state, word)
		}
	}
}
