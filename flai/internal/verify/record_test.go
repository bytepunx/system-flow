package verify

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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
