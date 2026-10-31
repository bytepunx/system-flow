package storygit

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// trialCounter is the real runner, counting the merge-tree processes it
// starts and answering version for git version when that is set.
type trialCounter struct {
	execx.System
	merges  int
	version string
}

func (c *trialCounter) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && slices.Contains(args, "merge-tree") {
		c.merges++
	}
	if name == "git" && c.version != "" && slices.Equal(args, []string{"version"}) {
		return c.version, nil
	}
	return c.System.Run(dir, name, args...)
}

// trialProject is a project with S-0001 in progress, whose branch changes
// docs/guide.md, and the real runner, skipping when git has no merge-tree
// --write-tree. Integration: it runs real git.
func trialProject(t *testing.T) (*workitem.Repo, *trialCounter) {
	t.Helper()
	repo := syncProject(t)
	r := &trialCounter{}
	if v, err := gitver.Installed(r); err != nil || !v.AtLeast(gitver.MergeTree) {
		t.Skipf("git %v is older than %s: %v", v, gitver.MergeTree, err)
	}
	_, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] guide", "docs/guide.md")
	return repo, r
}

// judgeOnce judges held against S-0001, the holder, and fails on an error.
func judgeOnce(t *testing.T, judge func(held, holder string) ([]string, error), held string) []string {
	t.Helper()
	got, err := judge(held, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// S-0551: holder's branch merges at main cleanly when main has not changed
// what it did, so the trial answers no paths.
func TestTrialJudgeAnswersNoneForACleanMerge(t *testing.T) {
	repo, r := trialProject(t)
	judge := NewTrialJudge(r, repo, Files{})
	if got := judgeOnce(t, judge, "S-0002"); len(got) != 0 {
		t.Fatalf("conflicts: %v", got)
	}
}

// S-0551: when main has changed what holder's branch changed, the trial at
// a held story's base, main, answers the path.
func TestTrialJudgeAnswersTheConflictingPath(t *testing.T) {
	repo, r := trialProject(t)
	commitFiles(t, repo.MainRoot, "chore: guide", "docs/guide.md", "docs/other.md")
	judge := NewTrialJudge(r, repo, Files{})
	if got := judgeOnce(t, judge, "S-0002"); !slices.Equal(got, []string{"docs/guide.md"}) {
		t.Fatalf("conflicts: %v", got)
	}
}

// S-0551: a conflict only on a file sync writes again or merges is no
// conflict to the judge.
func TestTrialJudgeLeavesOutTheFilesSyncResolves(t *testing.T) {
	repo, r := trialProject(t)
	commitFiles(t, repo.MainRoot, "chore: guide", "docs/guide.md")
	generated := NewTrialJudge(r, repo, Files{Generated: []GeneratedFile{{Path: "docs/guide.md"}}})
	if got := judgeOnce(t, generated, "S-0002"); len(got) != 0 {
		t.Fatalf("generated conflicts: %v", got)
	}
	merged := NewTrialJudge(r, repo, Files{Merged: MergedFiles{Match: func(p string) bool { return p == "docs/guide.md" }}})
	if got := judgeOnce(t, merged, "S-0002"); len(got) != 0 {
		t.Fatalf("merged conflicts: %v", got)
	}
}

// S-0551: a holder with no branch yet cannot be trial-merged, and the error
// says so.
func TestTrialJudgeErrsWhenTheHolderHasNoBranch(t *testing.T) {
	repo, r := trialProject(t)
	judge := NewTrialJudge(r, repo, Files{})
	_, err := judge("S-0002", "S-0003")
	if err == nil || !strings.Contains(err.Error(), "S-0003 has no branch story/S-0003 yet") {
		t.Fatalf("err: %v", err)
	}
}

// S-0551: a held story with a branch of its own, as one sent back has, is
// judged at that branch, not at main.
func TestTrialJudgeUsesTheHeldStorysOwnBranchAsItsBase(t *testing.T) {
	repo, r := trialProject(t)
	_, wt := openStory(t, repo, "S-0002", "docs")
	commitFiles(t, wt, "feat: [S-0002] guide", "docs/guide.md")
	judge := NewTrialJudge(r, repo, Files{})
	if got := judgeOnce(t, judge, "S-0002"); !slices.Equal(got, []string{"docs/guide.md"}) {
		t.Fatalf("at S-0002's branch: %v", got)
	}
	if got := judgeOnce(t, judge, "S-0003"); len(got) != 0 {
		t.Fatalf("at main: %v", got)
	}
}

// S-0551: the judge runs merge-tree once while the two heads stay where
// they are, and again when one moves.
func TestTrialJudgeMergesOncePerPairOfHeads(t *testing.T) {
	repo, r := trialProject(t)
	judge := NewTrialJudge(r, repo, Files{})
	for range 3 {
		judgeOnce(t, judge, "S-0002")
	}
	if r.merges != 1 {
		t.Fatalf("merge-tree ran %d times on the same heads, want 1", r.merges)
	}
	commitFiles(t, repo.WorktreePath("S-0001"), "feat: [S-0001] more", "docs/more.md")
	judgeOnce(t, judge, "S-0002")
	if r.merges != 2 {
		t.Fatalf("merge-tree ran %d times after the holder moved, want 2", r.merges)
	}
}

// S-0551: a git older than merge-tree --write-tree cannot trial-merge, and
// the error names its version.
func TestTrialJudgeErrsOnAGitTooOld(t *testing.T) {
	repo, r := trialProject(t)
	r.version = "git version 2.30.1"
	judge := NewTrialJudge(r, repo, Files{})
	_, err := judge("S-0002", "S-0001")
	if err == nil || !strings.Contains(err.Error(), "git 2.30 is older than 2.38") {
		t.Fatalf("err: %v", err)
	}
	if r.merges != 0 {
		t.Fatalf("merge-tree ran %d times", r.merges)
	}
}

// S-0551: outside a git checkout, or without a runner, there is no judge, so
// an overlap holds as it did before rather than saying git failed.
func TestNoTrialJudgeOutsideAGitCheckoutOrWithoutARunner(t *testing.T) {
	t.Parallel()
	repo := &workitem.Repo{Root: t.TempDir()}
	repo.MainRoot = repo.Root
	if NewTrialJudge(execx.System{}, repo, Files{}) != nil {
		t.Error("a folder that is no git checkout should get no trial judge")
	}
	git := syncProject(t)
	if NewTrialJudge(nil, git, Files{}) != nil {
		t.Error("no runner should give no trial judge")
	}
}
