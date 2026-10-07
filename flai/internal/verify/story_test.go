package verify

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"

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
		"diff --name-only --diff-filter=d main...HEAD": "flai/internal/x/x.go",
		"status --porcelain --untracked-files=all":     "",
		"rev-parse --git-path rebase-merge":            ".git/rebase-merge",
		"rev-parse --git-path rebase-apply":            ".git/rebase-apply",
		"rev-list --count HEAD..main":                  "0",
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

// edit rewrites a file of the fixture, failing the test when old is not in
// it.
func (f *verifyFixture) edit(t *testing.T, rel, old, new string) {
	t.Helper()
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
			f.git.out["rev-list --count HEAD..main"] = "3"
		}, want: Finding{Name: "sync", Message: "the branch does not contain main, 3 commits behind it; run flai stream sync S-004, resolve what it reports, and verify again"}},
		{name: "a main branch git cannot read", stopped: 1, setup: func(_ *testing.T, f *verifyFixture) {
			delete(f.git.out, "rev-list --count HEAD..main")
		}, want: Finding{Name: "sync", Message: "cannot tell whether the branch contains main: fatal: not a git repository; check that main is a branch of the repository"}},
		{name: "an empty Next steps", stopped: 2, setup: func(t *testing.T, f *verifyFixture) {
			f.edit(t, "wip/agents/S-004.md", "## Next steps\n1. n\n", "## Next steps\n1.\n\n")
		}, want: Finding{Name: "narrative", Path: "wip/agents/S-004.md", Line: 17, Message: "## Next steps is empty; write it, then verify again"}},
		{name: "no narrative", stopped: 2, setup: func(t *testing.T, f *verifyFixture) {
			if err := os.Remove(filepath.Join(f.root, "wip/agents/S-004.md")); err != nil {
				t.Fatal(err)
			}
		}, want: Finding{Name: "narrative", Path: "wip/agents/S-004.md"}},
		{name: "a check warning on the story's own narrative", stopped: 3, setup: func(t *testing.T, f *verifyFixture) {
			f.edit(t, "wip/agents/S-004.md", "\n## Open questions\n", "\n")
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

func TestVerifyCapsTheFailingStepsFindings(t *testing.T) {
	f := newVerifyFixture(t)
	f.edit(t, "wip/agents/S-004.md", "## Current state\ns\n\n## Next steps\n1. n\n", "## Current state\n\n## Next steps\n-\n")
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
