package verify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scriptGit answers the git commands it is given by their arguments and
// refuses every other, as git refuses a repository it is not run in.
type scriptGit struct {
	out map[string]string
	err map[string]error
	ran []string
}

func (g *scriptGit) Run(dir, name string, args ...string) (string, error) {
	return g.RunInput(dir, name, "", args...)
}

func (g *scriptGit) RunInput(_, name, _ string, args ...string) (string, error) {
	key := strings.Join(args, " ")
	g.ran = append(g.ran, name+" "+key)
	if err, ok := g.err[key]; ok {
		return g.out[key], err
	}
	if out, ok := g.out[key]; ok {
		return out, nil
	}
	return "", errors.New("fatal: not a git repository")
}

func (g *scriptGit) LookPath(name string) (string, error) { return name, nil }

// storyTiers are tiers like this repository's close-out runs: Go's tests
// and the markdown lint by path, the full Go tests and the template's by
// path but only all, and a smoke test for every story.
var storyTiers = []Tier{
	{Name: "go-test", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}},
	{Name: "markdown", Paths: []string{"**/*.md"}, Command: []string{"lint-md", "{files}"}},
	{Name: "integration", AllOnly: true, Dir: "flai", Paths: []string{"flai/**"}, Command: []string{"go", "test", "{packages}"}},
	{Name: "template", AllOnly: true, Paths: []string{"template/**"}, Command: []string{"template-test.sh"}},
	{Name: "smoke", AllOnly: true, Command: []string{"smoke.sh"}},
}

// verifyFixture is a project whose S-004 is in progress with its narrative
// written, checked out at its root, which stands for the story's worktree:
// its branch changed one Go file, contains main, and has no rebase left.
type verifyFixture struct {
	root string
	git  *scriptGit
	proc *fakeProc
	opts StoryOptions
}

func newVerifyFixture(t *testing.T) *verifyFixture {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	git := &scriptGit{out: map[string]string{
		"rev-parse HEAD": "0123abc\n",
		"rev-parse --verify --quiet main^{commit}":     "fedcba9\n",
		"diff --name-only --diff-filter=d main...HEAD": "flai/internal/x/x.go",
		"status --porcelain --untracked-files=all":     "",
		"rev-parse --git-path rebase-merge":            ".git/rebase-merge",
		"rev-parse --git-path rebase-apply":            ".git/rebase-apply",
		"diff --no-renames --name-only HEAD...main":    "",
		"rev-list --abbrev-commit HEAD..main":          "",
	}}
	proc := &fakeProc{steps: map[string]step{
		"go":       {took: 2 * time.Second},
		"smoke.sh": {took: 500 * time.Millisecond},
	}, now: time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/internal/x/x.go": {}}
	return &verifyFixture{root: root, git: git, proc: proc, opts: StoryOptions{
		Story: "s-4", Project: repo, Worktree: root, Base: "main", Tiers: storyTiers, Git: git,
		RunOptions: RunOptions{FS: fsys, Proc: proc, Now: proc.clock},
	}}
}

// editNarrative rewrites the story's narrative in the fixture, failing the
// test when old is not in it.
func (f *verifyFixture) editNarrative(t *testing.T, old, new string) {
	t.Helper()
	rel := "wip/agents/S-004.md"
	p := filepath.Join(f.root, rel)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), old) {
		t.Fatalf("%s has no %q", rel, old)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(data), old, new, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

// states are the report's steps as name:state.
func states(rep Report) []string {
	var out []string
	for _, s := range rep.Steps {
		out = append(out, s.Name+":"+string(s.State))
	}
	return out
}

func TestVerifyRunsEveryStepAndTheTiersTheBranchSelects(t *testing.T) {
	f := newVerifyFixture(t)
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:passed", "integration:passed", "smoke:passed"}
	if !rep.Passed || rep.StoppedAt != "" || !reflect.DeepEqual(states(rep), want) {
		t.Fatalf("report %+v, steps %v", rep, states(rep))
	}
	if rep.Story != "S-004" || rep.Commit != "0123abc" || rep.Base != "main" || !reflect.DeepEqual(rep.Paths, []string{"flai/internal/x/x.go"}) ||
		!rep.RanAt.Equal(time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)) || rep.DurationMS != 4500 || rep.Duration != "4.5s" {
		t.Errorf("report %+v", rep)
	}
	if !reflect.DeepEqual(f.proc.ran, []string{"go test ./internal/x", "go test ./...", "smoke.sh"}) || f.proc.dirs[0] != filepath.Join(f.root, "flai") {
		t.Errorf("ran %q in %q", f.proc.ran, f.proc.dirs)
	}
	for i, env := range f.proc.envs {
		if !slices.Contains(env, "CLOSE_OUT_STORY=S-004") {
			t.Errorf("%s ran with %q, not the story", f.proc.ran[i], env)
		}
	}
	zero := 0
	if got := rep.Steps[4]; !reflect.DeepEqual(got, Step{Name: "go-test", Tier: true, State: Passed, DurationMS: 2000, Duration: "2s",
		Command: []string{"go", "test", "./internal/x"}, Dir: "flai", ExitCode: &zero}) {
		t.Errorf("go-test step %+v", got)
	}
	if rep.Steps[0].Tier || rep.Steps[0].Command != nil {
		t.Errorf("rebase step %+v", rep.Steps[0])
	}
}

// S-0249: a finding outside the story is a note, not a failure: the
// fixture's epic is behind its stories, and another story's branch was
// never merged.
func TestVerifyAnswersTheCheckFindingsOutsideTheStoryAsNotes(t *testing.T) {
	f := newVerifyFixture(t)
	refs := filepath.Join(f.root, ".git", "refs", "heads", "story")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refs, "S-002"), []byte("0123456789abcdef0123456789abcdef01234567\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Passed || rep.Steps[3].State != Passed || len(rep.Steps[3].Findings) != 0 {
		t.Fatalf("the check failed on findings outside the story: %+v", rep)
	}
	rules := map[string]Note{}
	for _, n := range rep.Notes {
		rules[n.Rule] = n
	}
	unaccepted, ok := rules["story.unaccepted"]
	if !ok || unaccepted.Level != "warning" || unaccepted.Path != "wip/archive/kanban/stories/S-002-two.md" ||
		!strings.Contains(unaccepted.Message, "story/S-002 was never merged") {
		t.Errorf("notes %+v", rep.Notes)
	}
	if _, ok := rules["epic.lags-stories"]; !ok {
		t.Errorf("notes %+v lack the epic behind its stories", rep.Notes)
	}
}

func TestVerifyStopsAtTheFirstStepThatFails(t *testing.T) {
	notReachedAfter := func(i int) []string {
		all := []string{"rebase", "sync", "narrative", "check", "go-test", "integration", "smoke"}
		var out []string
		for j, n := range all {
			switch {
			case j < i:
				out = append(out, n+":passed")
			case j == i:
				out = append(out, n+":failed")
			default:
				out = append(out, n+":not-reached")
			}
		}
		return out
	}
	for _, tc := range []struct {
		name    string
		setup   func(t *testing.T, f *verifyFixture)
		stopped int
		want    Finding
	}{
		{name: "a rebase left unfinished", stopped: 0, setup: func(t *testing.T, f *verifyFixture) {
			dir := t.TempDir()
			f.git.out["rev-parse --git-path rebase-merge"] = dir
		}, want: Finding{Name: "rebase", Message: "a rebase is in progress; finish it with git rebase --continue, or undo it with git rebase --abort, then verify again"}},
		{name: "a branch behind main", stopped: 1, setup: func(_ *testing.T, f *verifyFixture) {
			f.git.out["diff --no-renames --name-only HEAD...main"] = "flai/go.mod\nwip/kanban/board.md"
			f.git.out["rev-list --abbrev-commit HEAD..main"] = "c3c3c3c\nb2b2b2b\na1a1a1a"
		}, want: Finding{Name: "sync", Message: "the branch does not contain main, 3 commits behind it, which change paths outside wip/: flai/go.mod; run flai stream sync S-004, resolve what it reports, and verify again"}},
		{name: "a main branch git cannot read", stopped: 1, setup: func(_ *testing.T, f *verifyFixture) {
			delete(f.git.out, "diff --no-renames --name-only HEAD...main")
		}, want: Finding{Name: "sync", Message: "cannot tell whether the branch contains main: list what main changed since the checkout left it: fatal: not a git repository; check that main is a branch of the repository"}},
		{name: "an empty Next steps", stopped: 2, setup: func(t *testing.T, f *verifyFixture) {
			f.editNarrative(t, "## Next steps\n1. n\n", "## Next steps\n1.\n\n")
		}, want: Finding{Name: "narrative", Path: "wip/agents/S-004.md", Line: 17, Message: "## Next steps is empty; write it, then verify again"}},
		{name: "no narrative", stopped: 2, setup: func(t *testing.T, f *verifyFixture) {
			if err := os.Remove(filepath.Join(f.root, "wip/agents/S-004.md")); err != nil {
				t.Fatal(err)
			}
		}, want: Finding{Name: "narrative", Path: "wip/agents/S-004.md"}},
		{name: "a check warning on the story's own narrative", stopped: 3, setup: func(t *testing.T, f *verifyFixture) {
			f.editNarrative(t, "\n## Open questions\n", "\n")
		}, want: Finding{Name: "narrative.section", Path: "wip/agents/S-004.md", Line: 1, Message: `warning: missing "## Open questions" section`}},
		{name: "a failing tier", stopped: 4, setup: func(_ *testing.T, f *verifyFixture) {
			f.proc.steps["go"] = step{stdout: "--- FAIL: TestX\n", exit: 1}
		}, want: Finding{Name: "go-test", Message: "--- FAIL: TestX\nexit status 1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			tc.setup(t, f)
			rep, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			want := notReachedAfter(tc.stopped)
			if rep.Passed || rep.StoppedAt != strings.Split(want[tc.stopped], ":")[0] || !reflect.DeepEqual(states(rep), want) {
				t.Fatalf("stopped at %q, steps %v, want %v", rep.StoppedAt, states(rep), want)
			}
			got := rep.Steps[tc.stopped].Findings
			if len(got) != 1 {
				t.Fatalf("findings %+v", got)
			}
			if tc.want.Message == "" {
				// the message carries the system's own words for a missing file
				if got[0].Name != tc.want.Name || got[0].Path != tc.want.Path || !strings.Contains(got[0].Message, "no narrative for S-004") {
					t.Errorf("finding %+v", got[0])
				}
			} else if !reflect.DeepEqual(got[0], tc.want) {
				t.Errorf("finding\n%+v\nwant\n%+v", got[0], tc.want)
			}
			if tc.stopped < 4 && len(f.proc.ran) != 0 {
				t.Errorf("tiers ran after %s failed: %q", rep.StoppedAt, f.proc.ran)
			}
			for _, s := range rep.Steps[tc.stopped+1:] {
				if s.DurationMS != 0 || s.Findings != nil {
					t.Errorf("a step not reached has a result: %+v", s)
				}
			}
			if s := rep.Steps[len(rep.Steps)-1]; !s.Tier || !reflect.DeepEqual(s.Command, []string{"smoke.sh"}) {
				t.Errorf("a tier not reached does not say what it would run: %+v", s)
			}
		})
	}
}

// S-0342: a tier that another selected tier covers is skipped, naming it,
// whether the steps before the tiers pass or one of them fails, and the
// stored report keeps what covers it.
func TestVerifySkipsATierThatASelectedTierCovers(t *testing.T) {
	covering := slices.Clone(storyTiers)
	covering[2].Covers = []string{"go-test"}
	skipped := Step{Name: "go-test", Tier: true, State: Skipped, Command: []string{"go", "test", "./internal/x"}, Dir: "flai", CoveredBy: "integration"}
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, f *verifyFixture)
		want  []string
		ran   []string
	}{
		{name: "every step passes", setup: func(*testing.T, *verifyFixture) {},
			want: []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:skipped", "integration:passed", "smoke:passed"},
			ran:  []string{"go test ./...", "smoke.sh"}},
		{name: "the check fails", setup: func(t *testing.T, f *verifyFixture) {
			f.editNarrative(t, "\n## Open questions\n", "\n")
		}, want: []string{"rebase:passed", "sync:passed", "narrative:passed", "check:failed", "go-test:skipped", "integration:not-reached", "smoke:not-reached"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.opts.Tiers = covering
			tc.setup(t, f)
			rep, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(states(rep), tc.want) || rep.Passed != (tc.ran != nil) {
				t.Fatalf("passed %v, steps %v, want %v", rep.Passed, states(rep), tc.want)
			}
			if !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Errorf("ran %q, want %q", f.proc.ran, tc.ran)
			}
			if got := rep.Steps[4]; !reflect.DeepEqual(got, skipped) {
				t.Errorf("go-test step\n%+v\nwant\n%+v", got, skipped)
			}
			if rep.Steps[5].CoveredBy != "" {
				t.Errorf("integration step %+v names a tier that covers it", rep.Steps[5])
			}
			if last, ok, err := LastReport(f.opts.Project, "S-004"); err != nil || !ok || !reflect.DeepEqual(last.Steps[4], skipped) {
				t.Errorf("stored go-test step %+v, %v, %v", last.Steps, ok, err)
			}
		})
	}
}

// ADR-0135: the sync step passes over commits of main that change only wip/
// paths the branch does not change, and names them in its note.
func TestVerifySyncPassesOverCommitsThatChangeOnlyWipPathsTheBranchDoesNotChange(t *testing.T) {
	var twelve []string
	for i := range 12 {
		twelve = append(twelve, fmt.Sprintf("c%06d", i))
	}
	for _, tc := range []struct {
		name, diff, commits, status string
		passed                      bool
		message                     string
	}{
		{name: "one replan", diff: "wip/kanban/stories/S-009-nine.md", commits: "a1a1a1a", passed: true,
			message: "passed over 1 commit of main that changes only wip/ paths the branch does not change: a1a1a1a"},
		{name: "twelve commits, a rename, a deletion", diff: "wip/archive/kanban/stories/S-009-nine.md\nwip/kanban/board.md\nwip/kanban/stories/S-009-nine.md",
			commits: strings.Join(twelve, "\n"), passed: true,
			message: "passed over 12 commits of main that change only wip/ paths the branch does not change: " + strings.Join(twelve[:10], ", ") + " and 2 more"},
		{name: "commits whose changes undo each other", commits: "b2b2b2b\na1a1a1a", passed: true,
			message: "passed over 2 commits of main that change only wip/ paths the branch does not change: b2b2b2b, a1a1a1a"},
		{name: "an acceptance", diff: "flai/a.go\nflai/b.go\nflai/c.go\nflai/d.go\nwip/kanban/board.md", commits: "a1a1a1a",
			message: "the branch does not contain main, 1 commit behind it, which changes paths outside wip/: flai/a.go, flai/b.go, flai/c.go and 1 more; run flai stream sync S-004, resolve what it reports, and verify again"},
		{name: "a path with the folder's name but outside it", diff: "wipe.md", commits: "a1a1a1a",
			message: "the branch does not contain main, 1 commit behind it, which changes paths outside wip/: wipe.md; run flai stream sync S-004, resolve what it reports, and verify again"},
		{name: "a wip path the branch has not committed", diff: "wip/kanban/board.md\nwip/agents/S-004.md", commits: "a1a1a1a", status: "M wip/agents/S-004.md",
			message: "the branch does not contain main, 1 commit behind it, which changes wip/ paths the branch changes too: wip/agents/S-004.md; run flai stream sync S-004, resolve what it reports, and verify again"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.git.out["diff --no-renames --name-only HEAD...main"] = tc.diff
			f.git.out["rev-list --abbrev-commit HEAD..main"] = tc.commits
			f.git.out["status --porcelain --untracked-files=all"] = tc.status
			rep, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			sync := rep.Steps[1]
			if rep.Passed != tc.passed || sync.Name != StepSync {
				t.Fatalf("passed %v, steps %v", rep.Passed, states(rep))
			}
			if tc.passed {
				if sync.State != Passed || sync.Note != tc.message || len(sync.Findings) != 0 {
					t.Errorf("sync step %+v; want the note %q", sync, tc.message)
				}
				for _, n := range rep.Notes {
					if n.Message == tc.message {
						t.Errorf("the step's note is a check note too: %+v", n)
					}
				}
				return
			}
			if sync.State != Failed || sync.Note != "" || len(sync.Findings) != 1 || sync.Findings[0].Message != tc.message {
				t.Errorf("sync step %+v\nwant the finding %q", sync, tc.message)
			}
		})
	}
}

func TestSyncOnlyRunsOnlyTheRebaseAndSyncStepsAndStoresNoReport(t *testing.T) {
	f := newVerifyFixture(t)
	f.git.out["diff --no-renames --name-only HEAD...main"] = "wip/kanban/board.md"
	f.git.out["rev-list --abbrev-commit HEAD..main"] = "a1a1a1a"
	rep, err := SyncOnly(f.opts)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Passed || rep.Story != "S-004" || rep.Commit != "0123abc" || !reflect.DeepEqual(states(rep), []string{"rebase:passed", "sync:passed"}) ||
		!strings.Contains(rep.Steps[1].Note, "a1a1a1a") {
		t.Errorf("report %+v", rep)
	}
	if len(f.proc.ran) != 0 {
		t.Errorf("tiers ran: %q", f.proc.ran)
	}
	if _, ok, err := LastReport(f.opts.Project, "S-004"); ok || err != nil {
		t.Errorf("SyncOnly stored a report: %v, %v", ok, err)
	}

	f.git.out["rev-parse --git-path rebase-merge"] = t.TempDir()
	rep, err = SyncOnly(f.opts)
	if err != nil || rep.Passed || rep.StoppedAt != StepRebase || !reflect.DeepEqual(states(rep), []string{"rebase:failed", "sync:not-reached"}) {
		t.Errorf("report %+v, %v", rep, err)
	}
	if _, err := SyncOnly(StoryOptions{Story: "S-004"}); err == nil {
		t.Error("SyncOnly without a project answered no error")
	}
}

// I-0119: a close-out verified the branch, and while it ran flai committed
// wip/ on main; the last check, that the branch contains main, failed though
// nothing the run verified had changed. This runs the close-out's last
// check, SyncOnly, on a real repository with the story's worktree.
func TestSyncOnlyPassesOverWipCommitsOnMainDuringACloseOut(t *testing.T) {
	root := realRepo(t)
	if err := os.CopyFS(root, os.DirFS("../metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, ".gitignore", ".flai-cache/\n")
	writeFile(t, root, "flai/x.go", "package x\n")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "init")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	worktree := repo.WorktreePath("S-004")
	runGit(t, root, "worktree", "add", "-q", "-b", "story/S-004", worktree)
	writeFile(t, worktree, "flai/x.go", "package x // S-004\n")
	runGit(t, worktree, "commit", "-q", "-am", "S-004's change")
	opts := StoryOptions{Story: "S-004", Project: repo, Base: "main", Git: execx.System{}}
	sync := func(t *testing.T) Step {
		t.Helper()
		rep, err := SyncOnly(opts)
		if err != nil {
			t.Fatal(err)
		}
		return rep.Steps[1]
	}
	if s := sync(t); s.State != Passed || s.Note != "" {
		t.Fatalf("the branch contains main, yet sync %+v", s)
	}

	// flai replans on main while the close-out runs
	writeFile(t, root, "wip/kanban/board.md", "replanned\n")
	runGit(t, root, "commit", "-q", "-am", "chore: replan forecasts")
	replan := runGit(t, root, "rev-parse", "--short", "HEAD")
	if s := sync(t); s.State != Passed || !strings.Contains(s.Note, "passed over 1 commit of main") || !strings.HasSuffix(s.Note, ": "+replan) {
		t.Errorf("after a replan on main, sync %+v; want it passed with a note naming %s", s, replan)
	}

	// the story's own narrative changed in its worktree and on main
	writeFile(t, worktree, "wip/agents/S-004.md", "the story's\n")
	writeFile(t, root, "wip/agents/S-004.md", "main's\n")
	runGit(t, root, "commit", "-q", "-am", "chore: narrative")
	if s := sync(t); s.State != Failed || !strings.Contains(s.Findings[0].Message, "paths the branch changes too: wip/agents/S-004.md") {
		t.Errorf("after main changed a wip path the branch changes, sync %+v", s)
	}
	runGit(t, worktree, "checkout", "-q", "--", "wip/agents/S-004.md")
	if s := sync(t); s.State != Passed || !strings.Contains(s.Note, "2 commits") {
		t.Errorf("once the branch no longer changes it, sync %+v", s)
	}

	// another story is accepted on main
	writeFile(t, root, "flai/y.go", "package x\n")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "merge: S-0005")
	if s := sync(t); s.State != Failed || !strings.Contains(s.Findings[0].Message, "outside wip/: flai/y.go; run flai stream sync S-004") {
		t.Errorf("after an acceptance on main, sync %+v", s)
	}
	if _, ok, err := LastReport(repo, "S-004"); ok || err != nil {
		t.Errorf("SyncOnly stored a report: %v, %v", ok, err)
	}
}

func TestVerifyCapsTheFailingStepsFindings(t *testing.T) {
	f := newVerifyFixture(t)
	f.editNarrative(t, "## Current state\ns\n\n## Next steps\n1. n\n", "## Current state\n\n## Next steps\n-\n")
	f.opts.Max = 1
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	if s := rep.Steps[2]; len(s.Findings) != 1 || s.Omitted != 1 || !strings.Contains(s.Findings[0].Message, "Current state") {
		t.Errorf("narrative step %+v", s)
	}
}

func TestVerifyRefusesWhatIsNotAStoryWithAWorktree(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		change     func(f *verifyFixture)
	}{
		{"no such story", "S-0099", func(f *verifyFixture) { f.opts.Story = "S-0099" }},
		{"a task", "it is a task", func(f *verifyFixture) { f.opts.Story = "T-003" }},
		{"no worktree", "flai stream open S-004", func(f *verifyFixture) { f.opts.Worktree = filepath.Join(f.root, "nowhere") }},
		{"a worktree git cannot read", "read the commit", func(f *verifyFixture) { delete(f.git.out, "rev-parse HEAD") }},
		{"no project", "no project", func(f *verifyFixture) { f.opts.Project = nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			tc.change(f)
			if _, err := Verify(context.Background(), f.opts); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %v, want one naming %q", err, tc.want)
			}
			if f.opts.Project == nil {
				return
			}
			if _, ok, err := LastReport(f.opts.Project, "S-004"); ok || err != nil {
				t.Errorf("a refused run stored a report: %v, %v", ok, err)
			}
		})
	}
}

func TestVerifyStoresTheLastReportForTheStory(t *testing.T) {
	f := newVerifyFixture(t)
	if _, ok, err := LastReport(f.opts.Project, "S-004"); ok || err != nil {
		t.Fatalf("a story never verified has a report: %v, %v", ok, err)
	}
	f.proc.steps["go"] = step{stdout: "boom\n", exit: 2}
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(f.root, ".flai-cache", "verify", "S-004.json"); ReportPath(f.opts.Project, "S-004") != want {
		t.Errorf("report path %s, want %s", ReportPath(f.opts.Project, "S-004"), want)
	}
	last, ok, err := LastReport(f.opts.Project, "s-4")
	if err != nil || !ok || !reflect.DeepEqual(last, rep) {
		t.Errorf("last report %+v, %v, %v\nwant %+v", last, ok, err, rep)
	}
	if err := os.WriteFile(ReportPath(f.opts.Project, "S-004"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LastReport(f.opts.Project, "S-004"); err == nil || !strings.Contains(err.Error(), "delete it and verify again") {
		t.Errorf("a report that does not read: %v", err)
	}
}

// S-0341: a run at the same head, base, and tiers as the last one runs the
// steps before the tiers and then only from the tier that failed, and a run
// after that names the run that ran each tier it reuses.
func TestVerifyResumesAtTheTierTheLastRunFailed(t *testing.T) {
	f := newVerifyFixture(t)
	f.proc.steps["smoke.sh"] = step{stdout: "network down\n", exit: 1, took: time.Second}
	first, err := Verify(context.Background(), f.opts)
	if err != nil || first.StoppedAt != "smoke" || first.BaseCommit != "fedcba9" || first.Fingerprint == "" {
		t.Fatalf("first run %+v, %v", first, err)
	}
	if n, _ := first.Reused(); n != 0 {
		t.Errorf("a first run reused %d tiers", n)
	}
	f.proc.steps["smoke.sh"] = step{took: time.Second}
	f.proc.ran, f.proc.now = nil, f.proc.now.Add(time.Minute)
	second, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:reused", "integration:reused", "smoke:passed"}
	if !second.Passed || !reflect.DeepEqual(states(second), want) || !reflect.DeepEqual(f.proc.ran, []string{"smoke.sh"}) {
		t.Fatalf("second run steps %v, ran %q", states(second), f.proc.ran)
	}
	zero := 0
	if got := second.Steps[4]; !reflect.DeepEqual(got, Step{Name: "go-test", Tier: true, State: Reused, DurationMS: 2000, Duration: "2s",
		Command: []string{"go", "test", "./internal/x"}, Dir: "flai", ExitCode: &zero, ReusedFrom: &first.RanAt}) {
		t.Errorf("reused step %+v", got)
	}
	if n, from := second.Reused(); n != 2 || !from.Equal(first.RanAt) {
		t.Errorf("reused %d from %v, want 2 from %v", n, from, first.RanAt)
	}
	if last, _, err := LastReport(f.opts.Project, "S-004"); err != nil || !reflect.DeepEqual(last, second) {
		t.Errorf("the stored report is not the whole run: %+v, %v", last, err)
	}
	f.proc.ran = nil
	third, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:reused", "integration:reused", "smoke:reused"}
	if !third.Passed || !reflect.DeepEqual(states(third), want) || len(f.proc.ran) != 0 {
		t.Fatalf("third run steps %v, ran %q", states(third), f.proc.ran)
	}
	for _, s := range third.Steps[4:] {
		wantFrom := first.RanAt
		if s.Name == "smoke" {
			wantFrom = second.RanAt
		}
		if s.ReusedFrom == nil || !s.ReusedFrom.Equal(wantFrom) {
			t.Errorf("%s reused from %v, want the run that ran it at %v", s.Name, s.ReusedFrom, wantFrom)
		}
	}
}

// S-0341, S-0342: a tier another selected tier covers does not stop a resumed
// run from reusing the tiers after it, and stays skipped rather than reused.
func TestVerifyResumesPastATierASelectedTierCovers(t *testing.T) {
	f := newVerifyFixture(t)
	f.opts.Tiers = slices.Clone(storyTiers)
	f.opts.Tiers[2].Covers = []string{"go-test"}
	f.proc.steps["smoke.sh"] = step{exit: 1}
	if first, err := Verify(context.Background(), f.opts); err != nil || first.StoppedAt != "smoke" {
		t.Fatalf("first run %+v, %v", first, err)
	}
	f.proc.steps["smoke.sh"] = step{}
	f.proc.ran = nil
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:skipped", "integration:reused", "smoke:passed"}
	if !rep.Passed || !reflect.DeepEqual(states(rep), want) || !reflect.DeepEqual(f.proc.ran, []string{"smoke.sh"}) {
		t.Fatalf("steps %v, ran %q", states(rep), f.proc.ran)
	}
	if s := rep.Steps[4]; s.CoveredBy != "integration" || s.ReusedFrom != nil {
		t.Errorf("go-test step %+v", s)
	}
	if n, _ := rep.Reused(); n != 1 {
		t.Errorf("reused %d tiers, want 1", n)
	}
}

// queueTiers are storyTiers with integration and smoke marked queue, as this
// repository marks them (ADR-0174).
func queueTiers() []Tier {
	tiers := slices.Clone(storyTiers)
	tiers[2].Queue, tiers[4].Queue = true, true
	return tiers
}

// ADR-0174: with LeaveToQueue, the selected tiers marked queue are not run
// but reported queued after the rest, whether the steps before the tiers
// pass or one fails, and the run passes without them; without it they run.
func TestVerifyLeavesTheQueueTiersToTheMergeQueue(t *testing.T) {
	queued := func(name string, argv ...string) Step {
		s := Step{Name: name, Tier: true, State: Queued, Command: argv}
		if name == "integration" {
			s.Dir = "flai"
		}
		return s
	}
	for _, tc := range []struct {
		name  string
		leave bool
		setup func(t *testing.T, f *verifyFixture)
		want  []string
		ran   []string
	}{
		{name: "left to the queue", leave: true, setup: func(*testing.T, *verifyFixture) {},
			want: []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:passed", "integration:queued", "smoke:queued"},
			ran:  []string{"go test ./internal/x"}},
		{name: "left to the queue after the check fails", leave: true, setup: func(t *testing.T, f *verifyFixture) {
			f.editNarrative(t, "\n## Open questions\n", "\n")
		}, want: []string{"rebase:passed", "sync:passed", "narrative:passed", "check:failed", "go-test:not-reached", "integration:queued", "smoke:queued"}},
		{name: "run without LeaveToQueue", setup: func(*testing.T, *verifyFixture) {},
			want: []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:passed", "integration:passed", "smoke:passed"},
			ran:  []string{"go test ./internal/x", "go test ./...", "smoke.sh"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), tc.leave
			tc.setup(t, f)
			rep, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(states(rep), tc.want) || rep.Passed != (tc.ran != nil) {
				t.Fatalf("passed %v, steps %v, want %v", rep.Passed, states(rep), tc.want)
			}
			if !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Errorf("ran %q, want %q", f.proc.ran, tc.ran)
			}
			if !tc.leave {
				return
			}
			want := []Step{queued("integration", "go", "test", "./..."), queued("smoke", "smoke.sh")}
			if !reflect.DeepEqual(rep.Steps[5:], want) {
				t.Errorf("queued steps\n%+v\nwant\n%+v", rep.Steps[5:], want)
			}
			if n, _ := rep.Reused(); n != 0 {
				t.Errorf("a first run reused %d tiers", n)
			}
		})
	}
}

// ADR-0174: a tier that a queue tier covers is not skipped for it while the
// queue tier is left out, since the covering tier does not run.
func TestVerifyRunsATierWhoseCoverIsLeftToTheMergeQueue(t *testing.T) {
	f := newVerifyFixture(t)
	f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), true
	f.opts.Tiers[2].Covers = []string{"go-test"}
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "go-test:passed", "integration:queued", "smoke:queued"}
	if !rep.Passed || !reflect.DeepEqual(states(rep), want) || !reflect.DeepEqual(f.proc.ran, []string{"go test ./internal/x"}) {
		t.Fatalf("steps %v, ran %q", states(rep), f.proc.ran)
	}
	if s := rep.Steps[4]; s.CoveredBy != "" {
		t.Errorf("go-test step %+v names a tier that covers it", s)
	}
}

// ADR-0174, I-0203: a run reuses the tiers the last run passed whether
// either left the queue tiers to the merge queue or ran them, since leaving
// them does not change the fingerprint, and runs what the last left queued:
// after a close-out that left them, flai verify --all runs only them.
func TestVerifyReusesAcrossLeavingToTheQueue(t *testing.T) {
	for _, tc := range []struct {
		name          string
		before, after bool
		covers        bool
		want          []string
		ran           []string
	}{
		{name: "left, then left again", before: true, after: true,
			want: []string{"go-test:reused", "integration:queued", "smoke:queued"}},
		{name: "left, then run", before: true, after: false,
			want: []string{"go-test:reused", "integration:passed", "smoke:passed"}, ran: []string{"go test ./...", "smoke.sh"}},
		{name: "run, then left", before: false, after: true,
			want: []string{"go-test:reused", "integration:queued", "smoke:queued"}},
		// go-test was skipped for the queue tier that covers it, which a run
		// that leaves the queue tiers does not run, so go-test runs
		{name: "run with a queue tier covering, then left", before: false, after: true, covers: true,
			want: []string{"go-test:passed", "integration:queued", "smoke:queued"}, ran: []string{"go test ./internal/x"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), tc.before
			if tc.covers {
				f.opts.Tiers[2].Covers = []string{"go-test"}
			}
			first, err := Verify(context.Background(), f.opts)
			if err != nil || !first.Passed {
				t.Fatalf("first run %+v, %v", first, err)
			}
			f.opts.LeaveToQueue = tc.after
			f.proc.ran, f.proc.now = nil, f.proc.now.Add(time.Minute)
			rep, err := Verify(context.Background(), f.opts)
			if err != nil || !rep.Passed {
				t.Fatalf("report %+v, %v", rep, err)
			}
			if rep.Fingerprint != first.Fingerprint {
				t.Errorf("the fingerprint changed with LeaveToQueue: %s, then %s", first.Fingerprint, rep.Fingerprint)
			}
			if got := tierStates(rep); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("tiers %v, want %v", got, tc.want)
			}
			if !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Errorf("ran %q, want %q", f.proc.ran, tc.ran)
			}
			if s := rep.Steps[4]; s.State == Reused && (s.ReusedFrom == nil || !s.ReusedFrom.Equal(first.RanAt)) {
				t.Errorf("go-test reused from %v, want the first run's %v", s.ReusedFrom, first.RanAt)
			}
		})
	}
}

// ADR-0174, I-0203: a change to a tier's command after a run that left the
// queue tiers to the merge queue changes the fingerprint, whether the tier
// ran or was left, so the next run reuses nothing.
func TestVerifyReusesNothingAcrossLeavingToTheQueueWhenATiersCommandChanges(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(tiers []Tier)
		ran    []string
	}{
		{"a tier that ran", func(tiers []Tier) { tiers[0].Command = []string{"go", "test", "-race", "{packages}"} },
			[]string{"go test -race ./internal/x", "go test ./...", "smoke.sh"}},
		{"a tier left to the queue", func(tiers []Tier) { tiers[4].Command = []string{"smoke.sh", "--full"} },
			[]string{"go test ./internal/x", "go test ./...", "smoke.sh --full"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), true
			first, err := Verify(context.Background(), f.opts)
			if err != nil || !first.Passed {
				t.Fatalf("first run %+v, %v", first, err)
			}
			tc.change(f.opts.Tiers)
			f.opts.LeaveToQueue = false
			f.proc.ran = nil
			rep, err := Verify(context.Background(), f.opts)
			if err != nil || !rep.Passed {
				t.Fatalf("report %+v, %v", rep, err)
			}
			if rep.Fingerprint == first.Fingerprint {
				t.Errorf("the fingerprint %s did not change with the command", rep.Fingerprint)
			}
			if n, _ := rep.Reused(); n != 0 || !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Errorf("reused %d, ran %q, want %q", n, f.proc.ran, tc.ran)
			}
		})
	}
}

// I-0203: QueueTiersNotRun says why a story's queue tiers cannot be taken as
// run at its head, and names those its last report at head did not run.
func TestQueueTiersNotRun(t *testing.T) {
	tier := func(name string, state State) Step { return Step{Name: name, Tier: true, State: state} }
	report := func(commit string, passed bool, steps ...Step) Report {
		rep := Report{Commit: commit, Passed: passed, Steps: append([]Step{{Name: StepCheck, State: Passed}}, steps...)}
		if !passed {
			rep.StoppedAt = "go-test"
		}
		return rep
	}
	const head = "0123abcdef0123abcdef"
	for _, tc := range []struct {
		name  string
		rep   Report
		ok    bool
		names []string
		why   string
	}{
		{"no report", Report{}, false, nil, "no flai verify has run on it"},
		{"at another commit", report("4567def0123", true, tier("go-test", Passed), tier("integration", Passed)), true,
			nil, "its last flai verify ran at 4567def, not at its head 0123abc"},
		{"failed at head", report(head, false, tier("go-test", Failed), tier("integration", Queued), tier("smoke", Queued)), true,
			[]string{"integration", "smoke"}, "its last flai verify stopped at go-test"},
		{"left to the queue at head", report(head, true, tier("go-test", Passed), tier("integration", Queued), tier("smoke", Queued)), true,
			[]string{"integration", "smoke"}, "its last flai verify did not run tiers integration and smoke, marked queue"},
		{"one left to the queue at head", report(head, true, tier("go-test", Reused), tier("smoke", Queued)), true,
			[]string{"smoke"}, "its last flai verify did not run tier smoke, marked queue"},
		{"all run at head", report(head, true, tier("go-test", Passed), tier("integration", Reused), tier("smoke", Passed)), true, nil, ""},
		{"a queue tier skipped for its cover at head", report(head, true, tier("go-test", Passed), tier("integration", Skipped), tier("smoke", Passed)), true, nil, ""},
		{"no queue tier selected at head", report(head, true, tier("markdown", Passed)), true, nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			names, why := QueueTiersNotRun(tc.rep, tc.ok, head, queueTiers())
			if !reflect.DeepEqual(names, tc.names) || why != tc.why {
				t.Errorf("answered %q, %q; want %q, %q", names, why, tc.names, tc.why)
			}
		})
	}
	t.Run("a close-out that left them", func(t *testing.T) {
		f := newVerifyFixture(t)
		f.opts.Tiers, f.opts.LeaveToQueue = queueTiers(), true
		rep, err := Verify(context.Background(), f.opts)
		if err != nil || !rep.Passed {
			t.Fatalf("report %+v, %v", rep, err)
		}
		names, why := QueueTiersNotRun(rep, true, "0123abc", f.opts.Tiers)
		if !reflect.DeepEqual(names, []string{"integration", "smoke"}) || why == "" {
			t.Errorf("answered %q, %q", names, why)
		}
	})
}

// S-0341: a run reuses nothing when what its tiers read or run may differ
// from the last run's, or when it is told not to.
func TestVerifyRunsEveryTierWhenTheLastRunIsNotTheSame(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(t *testing.T, f *verifyFixture)
	}{
		{"a new head", func(_ *testing.T, f *verifyFixture) { f.git.out["rev-parse HEAD"] = "4567def\n" }},
		{"a new base commit", func(_ *testing.T, f *verifyFixture) {
			f.git.out["rev-parse --verify --quiet main^{commit}"] = "9abcdef\n"
		}},
		{"a base git cannot read", func(_ *testing.T, f *verifyFixture) { delete(f.git.out, "rev-parse --verify --quiet main^{commit}") }},
		{"other tiers selected", func(_ *testing.T, f *verifyFixture) {
			f.git.out["diff --name-only --diff-filter=d main...HEAD"] = "flai/internal/x/x.go\ndocs/a.md"
		}},
		{"a new tier in the manifest", func(_ *testing.T, f *verifyFixture) {
			f.opts.Tiers = append(slices.Clone(storyTiers), Tier{Name: "vitest", Paths: []string{"flaiover/**"}, Command: []string{"vitest"}})
		}},
		{"a tier's new command", func(_ *testing.T, f *verifyFixture) {
			f.opts.Tiers = slices.Clone(storyTiers)
			f.opts.Tiers[0].Command = []string{"go", "test", "-race", "{packages}"}
		}},
		{"a report without a fingerprint", func(t *testing.T, f *verifyFixture) {
			last, _, err := LastReport(f.opts.Project, "S-004")
			if err != nil {
				t.Fatal(err)
			}
			last.Fingerprint = ""
			if err := SaveReport(f.opts.Project, last); err != nil {
				t.Fatal(err)
			}
		}},
		{"fresh", func(_ *testing.T, f *verifyFixture) { f.opts.Fresh = true }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.proc.steps["smoke.sh"] = step{exit: 1}
			if first, err := Verify(context.Background(), f.opts); err != nil || first.StoppedAt != "smoke" {
				t.Fatalf("first run %+v, %v", first, err)
			}
			tc.change(t, f)
			f.proc.steps["smoke.sh"] = step{}
			f.proc.ran = nil
			rep, err := Verify(context.Background(), f.opts)
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range rep.Steps {
				if s.State == Reused || s.ReusedFrom != nil {
					t.Errorf("reused %s: %+v", s.Name, s)
				}
			}
			if len(f.proc.ran) < 3 || !strings.HasPrefix(f.proc.ran[0], "go test") {
				t.Errorf("ran %q, not every tier", f.proc.ran)
			}
		})
	}
}

// S-0341: a run that stopped before the tiers leaves nothing to reuse.
func TestVerifyResumedAfterAStepBeforeTheTiersRunsEveryTier(t *testing.T) {
	f := newVerifyFixture(t)
	f.git.out["diff --no-renames --name-only HEAD...main"] = "flai/internal/y/y.go"
	f.git.out["rev-list --abbrev-commit HEAD..main"] = "a1a1a1a"
	if first, err := Verify(context.Background(), f.opts); err != nil || first.StoppedAt != "sync" {
		t.Fatalf("first run %+v, %v", first, err)
	}
	f.git.out["diff --no-renames --name-only HEAD...main"] = ""
	f.git.out["rev-list --abbrev-commit HEAD..main"] = ""
	rep, err := Verify(context.Background(), f.opts)
	if err != nil || !rep.Passed {
		t.Fatalf("report %+v, %v", rep, err)
	}
	if n, _ := rep.Reused(); n != 0 || !reflect.DeepEqual(f.proc.ran, []string{"go test ./internal/x", "go test ./...", "smoke.sh"}) {
		t.Errorf("reused %d, ran %q", n, f.proc.ran)
	}
}

func TestVerifyReportsTheContextEndingDuringATier(t *testing.T) {
	f := newVerifyFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	f.opts.Proc = &cancellingProc{cancel: cancel}
	rep, err := Verify(ctx, f.opts)
	if !errors.Is(err, context.Canceled) || rep.Passed || rep.StoppedAt != "go-test" {
		t.Errorf("report %+v, %v", rep, err)
	}
}

func TestUnwrittenReadsSectionsAsTheCloseOutDoes(t *testing.T) {
	sections := []string{"Current state", "Next steps"}
	for _, tc := range []struct {
		name, narrative string
		want            []Finding
	}{
		{name: "both written", narrative: "## Current state\nworking\n## Next steps\n1. ship\n"},
		{name: "text after blank lines and a list marker", narrative: "## Current state\n\n  \n- x\n## Next steps\n  12.   y\n"},
		{name: "a subheading is content", narrative: "## Current state\n### Detail\n## Next steps\n* z\n"},
		{name: "bare markers and blanks are empty", narrative: "# T\n## Current state\n1.\n -\n*\n\n## Next steps\nn\n",
			want: []Finding{{Name: "narrative", Line: 2, Message: "## Current state is empty; write it, then verify again"}}},
		{name: "a section with nothing before the next", narrative: "## Current state\nc\n## Next steps\n## Log\nl\n",
			want: []Finding{{Name: "narrative", Line: 3, Message: "## Next steps is empty; write it, then verify again"}}},
		{name: "a missing section", narrative: "## Current state\r\nc\r\n",
			want: []Finding{{Name: "narrative", Message: "the narrative has no ## Next steps section; add it and write it, then verify again"}}},
		{name: "a heading with more after it is another section", narrative: "## Current state now\nc\n## Next steps\nn\n",
			want: []Finding{{Name: "narrative", Message: "the narrative has no ## Current state section; add it and write it, then verify again"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Unwritten(tc.narrative, sections); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

// rebaseTiers are tiers a real repository's story selects by its Go file:
// unit and, after docs, which it does not select, integration and smoke;
// integration is selected by lib/ too.
var rebaseTiers = []Tier{
	{Name: "unit", Paths: []string{"flai/**/*.go"}, Command: []string{"unit", "{files}"}},
	{Name: "docs", Paths: []string{"docs/**"}, Command: []string{"docs", "{files}"}},
	{Name: "integration", AllOnly: true, Paths: []string{"flai/**", "lib/**"}, Command: []string{"integration"}},
	{Name: "smoke", AllOnly: true, Paths: []string{"flai/**"}, Command: []string{"smoke"}},
}

// wipTiers are rebaseTiers with an integration tier that wip paths select
// too, as a markdown tier selects wip's markdown.
var wipTiers = []Tier{
	{Name: "unit", Paths: []string{"flai/**/*.go"}, Command: []string{"unit", "{files}"}},
	{Name: "integration", AllOnly: true, Paths: []string{"flai/**", "wip/**"}, Command: []string{"integration"}},
	{Name: "smoke", AllOnly: true, Paths: []string{"flai/**"}, Command: []string{"smoke"}},
}

// rebaseFixture is a real repository whose S-004 has a worktree, on a
// branch that changed flai/x.go, and the options that verify it.
type rebaseFixture struct {
	root, worktree string
	proc           *fakeProc
	opts           StoryOptions
}

func newRebaseFixture(t *testing.T) *rebaseFixture {
	t.Helper()
	root := realRepo(t)
	if err := os.CopyFS(root, os.DirFS("../metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, ".gitignore", ".flai-cache/\n")
	writeFile(t, root, "flai/x.go", "package x\n")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "init")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	worktree := repo.WorktreePath("S-004")
	runGit(t, root, "worktree", "add", "-q", "-b", "story/S-004", worktree)
	writeFile(t, worktree, "flai/x.go", "package x // S-004\n")
	runGit(t, worktree, "commit", "-q", "-am", "S-004's change")
	proc := &fakeProc{now: time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)}
	return &rebaseFixture{root: root, worktree: worktree, proc: proc, opts: StoryOptions{
		Story: "S-004", Project: repo, Base: "main", Tiers: rebaseTiers, Git: execx.System{},
		RunOptions: RunOptions{Proc: proc, Now: proc.clock},
	}}
}

// verify runs Verify and answers its report, failing the test on an error;
// what the tiers ran is reset first.
func (f *rebaseFixture) verify(t *testing.T) Report {
	t.Helper()
	f.proc.ran = nil
	f.proc.now = f.proc.now.Add(time.Minute)
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

// commitOnMain commits a file on main, as an acceptance does while the
// story's close-out runs.
func (f *rebaseFixture) commitOnMain(t *testing.T, rel string) {
	t.Helper()
	writeFile(t, f.root, rel, "main's\n")
	runGit(t, f.root, "add", "-A")
	runGit(t, f.root, "commit", "-q", "-m", "merge: another story")
}

// rebase rebases the story's branch onto main, as flai stream sync does.
func (f *rebaseFixture) rebase(t *testing.T) {
	t.Helper()
	runGit(t, f.worktree, "rebase", "-q", "main")
}

// tierStates are the report's tiers as name:state.
func tierStates(rep Report) []string {
	var out []string
	for _, s := range rep.Steps {
		if s.Tier {
			out = append(out, s.Name+":"+string(s.State))
		}
	}
	return out
}

// I-0136: every tier passed, another story was accepted on main while they
// ran, and the branch rebased onto it cleanly; the run after the rebase ran
// every tier again, though main's commit changed nothing they test.
func TestVerifyReusesTheTiersAcrossACleanRebaseWhenMainsPathsSelectNone(t *testing.T) {
	f := newRebaseFixture(t)
	first := f.verify(t)
	if !first.Passed || !reflect.DeepEqual(f.proc.ran, []string{"unit flai/x.go", "integration", "smoke"}) {
		t.Fatalf("first run %v, ran %q: %+v", states(first), f.proc.ran, first)
	}
	if first.MergeBase == "" || first.Patch == "" {
		t.Fatalf("the first run recorded merge base %q and patch %q", first.MergeBase, first.Patch)
	}
	f.commitOnMain(t, "other/notes.txt")
	f.rebase(t)
	second := f.verify(t)
	want := []string{"unit:reused", "integration:reused", "smoke:reused"}
	if !second.Passed || !reflect.DeepEqual(tierStates(second), want) || len(f.proc.ran) != 0 {
		t.Fatalf("after a clean rebase, tiers %v, ran %q; want %v and nothing run", tierStates(second), f.proc.ran, want)
	}
	if second.Commit == first.Commit || second.MergeBase == first.MergeBase || second.Patch != first.Patch {
		t.Errorf("commit %s, merge base %s, patch %s after the rebase; first %s, %s, %s",
			second.Commit, second.MergeBase, second.Patch, first.Commit, first.MergeBase, first.Patch)
	}
	note := "reused across the rebase onto main: the branch's own change is the same as at " + first.Commit[:7] +
		", and the 1 path main brought in selects none of the tiers reused"
	if s := second.Steps[4]; s.Note != note || !s.ReusedFrom.Equal(first.RanAt) {
		t.Errorf("first reused tier %+v; want the note %q", s, note)
	}
	for _, s := range second.Steps[5:] {
		if s.Note != "" {
			t.Errorf("%s has a note: %q", s.Name, s.Note)
		}
	}
	// a run after it, at the same head and base, reuses them with no note
	third := f.verify(t)
	if !reflect.DeepEqual(tierStates(third), want) || len(f.proc.ran) != 0 || third.Steps[4].Note != "" {
		t.Errorf("third run tiers %v, ran %q, note %q", tierStates(third), f.proc.ran, third.Steps[4].Note)
	}
}

// ADR-0156: a tier that main's new paths select, and every tier after it,
// runs again after the rebase; the tiers before it are reused.
func TestVerifyRunsAgainTheTiersMainsNewPathsSelectAndThoseAfterThem(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		want, ran  []string
		note       string
	}{
		{name: "a path integration selects", path: "lib/y.txt",
			want: []string{"unit:reused", "integration:passed", "smoke:passed"}, ran: []string{"integration", "smoke"},
			note: "run again from here: the 1 path main brought in since the last run selects integration: lib/y.txt"},
		{name: "a Go file of another package", path: "flai/y/y.go",
			want: []string{"unit:passed", "integration:passed", "smoke:passed"}, ran: []string{"unit flai/x.go", "integration", "smoke"},
			note: "run again from here: the 1 path main brought in since the last run selects unit, integration, smoke: flai/y/y.go"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newRebaseFixture(t)
			if first := f.verify(t); !first.Passed {
				t.Fatalf("first run %v", states(first))
			}
			f.commitOnMain(t, tc.path)
			f.rebase(t)
			rep := f.verify(t)
			if !rep.Passed || !reflect.DeepEqual(tierStates(rep), tc.want) || !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Fatalf("tiers %v, ran %q; want %v, %q", tierStates(rep), f.proc.ran, tc.want, tc.ran)
			}
			at := 4 + len(tc.want) - len(tc.ran)
			if s := rep.Steps[at]; s.Note != tc.note {
				t.Errorf("first tier run again %+v; want the note %q", s, tc.note)
			}
		})
	}
}

// ADR-0156: across a rebase, a run reuses nothing when the branch's own
// change is not the same, when it is told not to, or when the last report
// had uncommitted changes, which its patch id would not cover.
func TestVerifyReusesNothingAcrossARebaseWhenTheBranchsChangeMayDiffer(t *testing.T) {
	for _, tc := range []struct {
		name          string
		before, after func(t *testing.T, f *rebaseFixture)
	}{
		{name: "another commit on the branch", after: func(t *testing.T, f *rebaseFixture) {
			writeFile(t, f.worktree, "flai/notes.txt", "more\n")
			runGit(t, f.worktree, "add", "-A")
			runGit(t, f.worktree, "commit", "-q", "-m", "S-004's notes")
		}},
		{name: "fresh", after: func(_ *testing.T, f *rebaseFixture) { f.opts.Fresh = true }},
		{name: "a report with an uncommitted change", before: func(t *testing.T, f *rebaseFixture) {
			writeFile(t, f.worktree, "scratch.txt", "draft\n")
		}, after: func(t *testing.T, f *rebaseFixture) {
			if err := os.Remove(filepath.Join(f.worktree, "scratch.txt")); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newRebaseFixture(t)
			if tc.before != nil {
				tc.before(t, f)
			}
			first := f.verify(t)
			if !first.Passed {
				t.Fatalf("first run %v", states(first))
			}
			if tc.before != nil && first.Patch != "" {
				t.Errorf("a run with an uncommitted change recorded the patch %q", first.Patch)
			}
			f.commitOnMain(t, "other/notes.txt")
			tc.after(t, f)
			f.rebase(t)
			rep := f.verify(t)
			want := []string{"unit:passed", "integration:passed", "smoke:passed"}
			if !rep.Passed || !reflect.DeepEqual(tierStates(rep), want) || len(f.proc.ran) != 3 {
				t.Errorf("tiers %v, ran %q; want every tier run", tierStates(rep), f.proc.ran)
			}
		})
	}
}

// ADR-0135, ADR-0156: a commit on main that changes only wip/ paths the
// branch does not change is passed over by the sync step, and since the
// branch was not rebased, its tiers are reused.
func TestVerifyReusesTheTiersWhenMainGainsOnlyWipCommits(t *testing.T) {
	f := newRebaseFixture(t)
	first := f.verify(t)
	if !first.Passed {
		t.Fatalf("first run %v", states(first))
	}
	f.commitOnMain(t, "wip/notes.txt")
	rep := f.verify(t)
	want := []string{"unit:reused", "integration:reused", "smoke:reused"}
	if !rep.Passed || !reflect.DeepEqual(tierStates(rep), want) || len(f.proc.ran) != 0 {
		t.Fatalf("steps %v, ran %q; want %v", states(rep), f.proc.ran, want)
	}
	if s := rep.Steps[1]; s.State != Passed || !strings.Contains(s.Note, "passed over 1 commit of main") {
		t.Errorf("sync step %+v", s)
	}
	note := "reused though main moved on: the branch's own change is the same as at " + first.Commit[:7] + ", and it was not rebased"
	if rep.Commit != first.Commit || rep.BaseCommit == first.BaseCommit || rep.Steps[4].Note != note {
		t.Errorf("commit %s, base %s, note %q; first %s, %s; want the note %q", rep.Commit, rep.BaseCommit, rep.Steps[4].Note, first.Commit, first.BaseCommit, note)
	}
}

// I-0196: the branch and main each recorded an issue, so the rebase stopped
// on the issues summary, which the sync wrote again; the branch's summary
// differs after it, but its patch id leaves the summary out, and main's
// issue paths select no tier, so every tier is reused.
func TestVerifyReusesTheTiersAcrossARebaseThatWritesTheSummaryAgain(t *testing.T) {
	f := newRebaseFixture(t)
	record := func(root, title string) {
		t.Helper()
		r, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := issues.New(r, issues.NewOptions{Title: title, Class: "defect", Now: f.proc.now}); err != nil {
			t.Fatal(err)
		}
		if _, err := issues.WriteSummary(r, f.proc.now); err != nil {
			t.Fatal(err)
		}
		runGit(t, root, "add", "-A")
		runGit(t, root, "commit", "-q", "-m", "record "+title)
	}
	record(f.worktree, "Branch friction")
	first := f.verify(t)
	if !first.Passed || first.Patch == "" {
		t.Fatalf("first run %v, patch %q", states(first), first.Patch)
	}
	f.proc.now = f.proc.now.Add(time.Hour)
	record(f.root, "Main friction")
	if _, err := f.opts.Git.Run(f.worktree, "git", "rebase", "-q", "main"); err == nil {
		t.Fatal("the rebase did not stop on the summary")
	}
	// the sync writes the summary again in the story's worktree
	wt := *f.opts.Project
	wt.Root = f.worktree
	for _, g := range issues.Generated(&wt, "", func() time.Time { return f.proc.now }) {
		if err := g.Regenerate(); err != nil {
			t.Fatal(err)
		}
		runGit(t, f.worktree, "add", g.Path)
	}
	runGit(t, f.worktree, "-c", "core.editor=true", "rebase", "--continue")

	rep := f.verify(t)
	want := []string{"unit:reused", "integration:reused", "smoke:reused"}
	if !rep.Passed || !reflect.DeepEqual(tierStates(rep), want) || len(f.proc.ran) != 0 {
		t.Fatalf("after a rebase that wrote the summary again, steps %v, ran %q; want tiers %v and nothing run", states(rep), f.proc.ran, want)
	}
	if rep.MergeBase == first.MergeBase || rep.Patch != first.Patch {
		t.Errorf("merge base %s, patch %s after the rebase; first %s, %s", rep.MergeBase, rep.Patch, first.MergeBase, first.Patch)
	}
	if note := rep.Steps[4].Note; !strings.Contains(note, "reused across the rebase onto main") {
		t.Errorf("first reused tier's note %q", note)
	}
}

// patchDiff is the git diff BranchPatch runs for a fixture's branch from the
// merge base mb: the issues summary left out.
func patchDiff(mb string) string {
	return "diff --no-color --no-ext-diff --full-index " + mb + " HEAD -- :(top,exclude)design/issues/summary.md"
}

// ADR-0156, without real git: what a run after a rebase reuses, by what git
// answers for the branch's merge base, its patch id, and main's new paths.
func TestVerifyAcrossARebaseReusesUpToWhatMainsNewPathsSelect(t *testing.T) {
	covering := []Tier{
		{Name: "lint", Paths: []string{"flai/**/*.go"}, Command: []string{"lint", "{files}"}, Covers: []string{"check-y"}},
		{Name: "unit", Paths: []string{"flai/**/*.go"}, Command: []string{"unit", "{files}"}},
		{Name: "check-y", Paths: []string{"flai/**", "lib/**"}, Command: []string{"check-y"}},
	}
	every := []string{"unit flai/internal/x/x.go", "integration", "smoke"}
	for _, tc := range []struct {
		name   string
		tiers  []Tier
		change func(f *verifyFixture)
		want   []string
		ran    []string
		notes  []string
	}{
		{name: "main's paths select no tier", want: []string{"unit:reused", "integration:reused", "smoke:reused"},
			notes: []string{"reused across the rebase onto main: the branch's own change is the same as at 0123abc, and the 1 path main brought in selects none of the tiers reused"}},
		{name: "main's paths select integration", change: func(f *verifyFixture) { f.git.out["diff --no-renames --name-only mb1 mb2"] = "lib/y.txt" },
			want: []string{"unit:reused", "integration:passed", "smoke:passed"}, ran: []string{"integration", "smoke"},
			notes: []string{"reused across the rebase onto main: the branch's own change is the same as at 0123abc, and the 1 path main brought in selects none of the tiers reused", "run again from here: the 1 path main brought in since the last run selects integration: lib/y.txt"}},
		{name: "main's paths select a tier that one before it covers", tiers: covering,
			change: func(f *verifyFixture) { f.git.out["diff --no-renames --name-only mb1 mb2"] = "lib/y.txt" },
			want:   []string{"lint:passed", "unit:passed", "check-y:skipped"}, ran: []string{"lint flai/internal/x/x.go", "unit flai/internal/x/x.go"},
			notes: []string{"run again from here: the 1 path main brought in since the last run selects lint, check-y: lib/y.txt"}},
		{name: "main's paths are wip paths the branch does not change", tiers: wipTiers,
			change: func(f *verifyFixture) { f.git.out["diff --no-renames --name-only mb1 mb2"] = "wip/kanban/notes.md" },
			want:   []string{"unit:reused", "integration:reused", "smoke:reused"},
			notes:  []string{"reused across the rebase onto main: the branch's own change is the same as at 0123abc, and main brought in only wip paths the branch does not change"}},
		{name: "another patch", change: func(f *verifyFixture) { f.git.out["patch-id --stable"] = "p2 4567def" },
			want: []string{"unit:passed", "integration:passed", "smoke:passed"}, ran: every},
		{name: "fresh", change: func(f *verifyFixture) { f.opts.Fresh = true },
			want: []string{"unit:passed", "integration:passed", "smoke:passed"}, ran: every},
		{name: "main's new paths git cannot list", change: func(f *verifyFixture) { delete(f.git.out, "diff --no-renames --name-only mb1 mb2") },
			want: []string{"unit:passed", "integration:passed", "smoke:passed"}, ran: every},
		{name: "a merge base git cannot read", change: func(f *verifyFixture) { delete(f.git.out, "merge-base HEAD main") },
			want: []string{"unit:passed", "integration:passed", "smoke:passed"}, ran: every},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVerifyFixture(t)
			f.opts.Tiers = rebaseTiers
			if tc.tiers != nil {
				f.opts.Tiers = tc.tiers
			}
			diff := "diff --git a/flai/internal/x/x.go b/flai/internal/x/x.go"
			f.git.out["merge-base HEAD main"] = "mb1"
			f.git.out[patchDiff("mb1")] = diff
			f.git.out["patch-id --stable"] = "p1 0123abc"
			first, err := Verify(context.Background(), f.opts)
			if err != nil || !first.Passed || first.MergeBase != "mb1" || first.Patch != "p1" {
				t.Fatalf("first run %+v, %v", first, err)
			}
			// the branch is rebased onto main, which gained a commit
			f.git.out["rev-parse HEAD"] = "4567def"
			f.git.out["rev-parse --verify --quiet main^{commit}"] = "9abcdef"
			f.git.out["merge-base HEAD main"] = "mb2"
			f.git.out[patchDiff("mb2")] = diff
			f.git.out["patch-id --stable"] = "p1 4567def"
			f.git.out["diff --no-renames --name-only mb1 mb2"] = "other/notes.txt"
			if tc.change != nil {
				tc.change(f)
			}
			f.proc.ran = nil
			rep, err := Verify(context.Background(), f.opts)
			if err != nil || !rep.Passed {
				t.Fatalf("report %+v, %v", rep, err)
			}
			if !reflect.DeepEqual(tierStates(rep), tc.want) || !reflect.DeepEqual(f.proc.ran, tc.ran) {
				t.Fatalf("tiers %v, ran %q; want %v, %q", tierStates(rep), f.proc.ran, tc.want, tc.ran)
			}
			var notes []string
			for _, s := range rep.Steps[4:] {
				if s.Note != "" {
					notes = append(notes, s.Note)
				}
			}
			if !reflect.DeepEqual(notes, tc.notes) {
				t.Errorf("notes %q, want %q", notes, tc.notes)
			}
		})
	}
}

// ADR-0156: a report taken with something uncommitted records no patch id,
// so the run after a rebase reuses nothing.
func TestVerifyRecordsNoPatchWithAnUncommittedChange(t *testing.T) {
	f := newVerifyFixture(t)
	f.opts.Tiers = rebaseTiers
	f.git.out["status --porcelain --untracked-files=all"] = "?? scratch.txt"
	f.git.out["merge-base HEAD main"] = "mb1"
	f.git.out[patchDiff("mb1")] = "diff --git a/flai/internal/x/x.go b/flai/internal/x/x.go"
	f.git.out["patch-id --stable"] = "p1 0123abc"
	first, err := Verify(context.Background(), f.opts)
	if err != nil || !first.Passed || first.MergeBase != "mb1" || first.Patch != "" {
		t.Fatalf("first run %+v, %v", first, err)
	}
	f.git.out["status --porcelain --untracked-files=all"] = ""
	f.git.out["rev-parse HEAD"] = "4567def"
	f.git.out["merge-base HEAD main"] = "mb2"
	f.git.out[patchDiff("mb2")] = "diff --git a/flai/internal/x/x.go b/flai/internal/x/x.go"
	f.git.out["diff --no-renames --name-only mb1 mb2"] = "other/notes.txt"
	f.proc.ran = nil
	rep, err := Verify(context.Background(), f.opts)
	if err != nil || rep.Patch != "p1" {
		t.Fatalf("report %+v, %v", rep, err)
	}
	if n, _ := rep.Reused(); n != 0 || len(f.proc.ran) != 3 {
		t.Errorf("reused %d, ran %q; want every tier run", n, f.proc.ran)
	}
}

// tierOptions are the options that verify the fixture's checkout with
// Tiers, against main, with its tiers, git, and processes, keeping the
// tiers each callback is called with.
func (f *verifyFixture) tierOptions(started, ended *[]string) TierOptions {
	return TierOptions{
		Dir: f.root, Base: "main", Tiers: f.opts.Tiers, Git: f.git, RunOptions: f.opts.RunOptions,
		Started: func(tier string) { *started = append(*started, tier) },
		Ended:   func(s Step) { *ended = append(*ended, s.Name+":"+string(s.State)) },
	}
}

// S-0399: Tiers runs the tiers flai verify's tier step selects for
// what the checkout changed against the base, with none of the steps before
// them and no CLOSE_OUT_STORY, calls back as each starts and ends, and
// leaves the story's last report as it was.
func TestVerifyTiersRunsTheTiersTheChangesSelectAndStoresNoReport(t *testing.T) {
	f := newVerifyFixture(t)
	before := Report{Story: "S-004", Commit: "aaaa", Passed: false, StoppedAt: "check", Steps: []Step{}}
	if err := SaveReport(f.opts.Project, before); err != nil {
		t.Fatal(err)
	}
	var started, ended []string
	rep, err := Tiers(context.Background(), f.tierOptions(&started, &ended))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"go-test:passed", "integration:passed", "smoke:passed"}
	if !rep.Passed || rep.StoppedAt != "" || !reflect.DeepEqual(states(rep), want) {
		t.Fatalf("report %+v, steps %v, want %v", rep, states(rep), want)
	}
	if rep.Story != "" || rep.Commit != "0123abc" || rep.Base != "main" || rep.BaseCommit != "fedcba9" || rep.Fingerprint == "" ||
		!reflect.DeepEqual(rep.Paths, []string{"flai/internal/x/x.go"}) || rep.DurationMS != 4500 {
		t.Errorf("report %+v", rep)
	}
	if !reflect.DeepEqual(f.proc.ran, []string{"go test ./internal/x", "go test ./...", "smoke.sh"}) {
		t.Errorf("ran %q", f.proc.ran)
	}
	if !reflect.DeepEqual(started, []string{"go-test", "integration", "smoke"}) || !reflect.DeepEqual(ended, want) {
		t.Errorf("started %v, ended %v", started, ended)
	}
	for _, env := range f.proc.envs {
		if slices.ContainsFunc(env, func(e string) bool { return strings.HasPrefix(e, StoryEnv+"=") }) {
			t.Errorf("a tier ran with %q, want no %s", env, StoryEnv)
		}
	}
	for _, ran := range f.git.ran {
		if strings.Contains(ran, "rebase-merge") || strings.Contains(ran, "HEAD...main") {
			t.Errorf("git ran %q, a step before the tiers", ran)
		}
	}
	if last, ok, err := LastReport(f.opts.Project, "S-004"); err != nil || !ok || last.Commit != "aaaa" || last.StoppedAt != "check" {
		t.Errorf("the story's last report %+v, %v, %v; want it as it was", last, ok, err)
	}
}

// S-0399: Tiers stops at the first tier that fails, leaves the tiers
// after it not reached and calls back for none of them, and skips a tier a
// selected tier covers, as Verify does.
func TestVerifyTiersStopsAtTheFirstTierThatFails(t *testing.T) {
	f := newVerifyFixture(t)
	f.proc.steps["go"] = step{stdout: "", stderr: "x_test.go:3: boom\n", exit: 1, took: time.Second}
	var started, ended []string
	rep, err := Tiers(context.Background(), f.tierOptions(&started, &ended))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"go-test:failed", "integration:not-reached", "smoke:not-reached"}; rep.Passed || rep.StoppedAt != "go-test" || !reflect.DeepEqual(states(rep), want) {
		t.Errorf("passed %v at %q, steps %v, want %v", rep.Passed, rep.StoppedAt, states(rep), want)
	}
	if len(rep.Steps[0].Findings) == 0 || !reflect.DeepEqual(started, []string{"go-test"}) || !reflect.DeepEqual(ended, []string{"go-test:failed"}) {
		t.Errorf("findings %+v, started %v, ended %v", rep.Steps[0].Findings, started, ended)
	}

	covered := newVerifyFixture(t)
	covered.opts.Tiers = slices.Clone(storyTiers)
	covered.opts.Tiers[2].Covers = []string{"go-test"}
	started, ended = nil, nil
	rep, err = Tiers(context.Background(), covered.tierOptions(&started, &ended))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"go-test:skipped", "integration:passed", "smoke:passed"}; !rep.Passed || !reflect.DeepEqual(states(rep), want) || rep.Steps[0].CoveredBy != "integration" {
		t.Errorf("steps %v, want %v, covered by integration", states(rep), want)
	}
	if !reflect.DeepEqual(started, []string{"integration", "smoke"}) {
		t.Errorf("started %v, want the tiers that ran", started)
	}
}

// S-0399: a context that ends during a tier stops Tiers there, with
// the report so far and the context's error; the tiers after it are not
// reached.
func TestVerifyTiersReportsTheContextEnding(t *testing.T) {
	f := newVerifyFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	opts := f.tierOptions(new([]string), new([]string))
	opts.Proc = &cancellingProc{cancel: cancel}
	rep, err := Tiers(ctx, opts)
	if want := []string{"go-test:failed", "integration:not-reached", "smoke:not-reached"}; !errors.Is(err, context.Canceled) || rep.Passed || !reflect.DeepEqual(states(rep), want) {
		t.Errorf("report %v, %v", states(rep), err)
	}
}

// S-0399: a checkout whose HEAD git cannot read is an error with no report.
func TestVerifyTiersRefusesACheckoutGitCannotRead(t *testing.T) {
	f := newVerifyFixture(t)
	delete(f.git.out, "rev-parse HEAD")
	if rep, err := Tiers(context.Background(), f.tierOptions(new([]string), new([]string))); err == nil || !strings.Contains(err.Error(), f.root) || len(rep.Steps) != 0 {
		t.Errorf("report %+v, %v; want an error naming the checkout", rep, err)
	}
}
