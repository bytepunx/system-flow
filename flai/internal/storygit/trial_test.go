package storygit

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
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

// trialGit answers the git a trial judge runs in a project of no commits,
// counting each call by its subcommand: version as version says, the heads
// of the local branches from heads, with HEAD at main, and every merge
// clean.
type trialGit struct {
	version string
	heads   map[string]string
	calls   map[string]int
}

func (g *trialGit) Run(dir, name string, args ...string) (string, error) {
	if name != "git" || len(args) == 0 {
		return "", fmt.Errorf("%s %v is not git", name, args)
	}
	g.calls[args[0]]++
	switch args[0] {
	case "version":
		return g.version, nil
	case "for-each-ref":
		var out strings.Builder
		for branch, sha := range g.heads {
			head := " "
			if branch == "main" {
				head = "*"
			}
			fmt.Fprintf(&out, "refs/heads/%s\t%s\t\t\tcommit\t%s\n", branch, head, sha)
		}
		return strings.TrimSpace(out.String()), nil
	case "merge-tree":
		return "", nil
	}
	return "", fmt.Errorf("git %v is not a call the trial judge makes", args)
}

func (g *trialGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *trialGit) LookPath(name string) (string, error) { return name, nil }

// trialFake is a project in a folder of its own with a .git, so that kept
// merges of other tests are not its own, judged with a trialGit where main
// and S-0001's branch have a head and git can trial-merge.
func trialFake(t *testing.T) (*workitem.Repo, *trialGit) {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	git := &trialGit{
		version: "git version 2.47.0",
		heads:   map[string]string{"main": "m1", Branch("S-0001"): "a1"},
		calls:   map[string]int{},
	}
	return &workitem.Repo{Root: root, MainRoot: root}, git
}

// I-0313: the judge the next agent.status call makes over the same commits
// runs no merge-tree and no git version, and reads the heads of the local
// branches and the main branch with one for-each-ref.
func TestASecondTrialJudgeOverTheSameCommitsMergesNothing(t *testing.T) {
	repo, git := trialFake(t)
	judgeOnce(t, NewTrialJudge(git, repo, Files{}), "S-0002")
	if want := map[string]int{"for-each-ref": 1, "version": 1, "merge-tree": 1}; !maps.Equal(git.calls, want) {
		t.Fatalf("first judge's git calls: %v, want %v", git.calls, want)
	}
	clear(git.calls)
	judgeOnce(t, NewTrialJudge(git, repo, Files{}), "S-0002")
	if want := map[string]int{"for-each-ref": 1}; !maps.Equal(git.calls, want) {
		t.Fatalf("second judge's git calls: %v, want %v", git.calls, want)
	}
}

// I-0313: a kept merge is of two commits, so the next judge merges again
// once the holder's head has moved, and so does a judge kept for a session,
// as the MCP server's is, which reads the heads on every call.
func TestTrialJudgeMergesAMovedHeadAgain(t *testing.T) {
	repo, git := trialFake(t)
	kept := NewTrialJudge(git, repo, Files{})
	judgeOnce(t, kept, "S-0002")
	git.heads[Branch("S-0001")] = "a2"
	judgeOnce(t, NewTrialJudge(git, repo, Files{}), "S-0002")
	if git.calls["merge-tree"] != 2 {
		t.Fatalf("merge-tree ran %d times over a moved head, want 2", git.calls["merge-tree"])
	}
	git.heads[Branch("S-0001")] = "a3"
	judgeOnce(t, kept, "S-0002")
	if git.calls["merge-tree"] != 3 {
		t.Fatalf("merge-tree ran %d times once the kept judge's holder moved, want 3", git.calls["merge-tree"])
	}
}

// I-0313: a git too old to trial-merge is not kept as an answer, so the
// judge after it sees the git installed since, and asks its version once.
func TestTrialJudgeSeesAGitInstalledAfterOneTooOld(t *testing.T) {
	repo, git := trialFake(t)
	git.version = "git version 2.30.1"
	if _, err := NewTrialJudge(git, repo, Files{})("S-0002", "S-0001"); err == nil {
		t.Fatal("git 2.30 should be too old to trial-merge")
	}
	git.version = "git version 2.47.0"
	judgeOnce(t, NewTrialJudge(git, repo, Files{}), "S-0002")
	judgeOnce(t, NewTrialJudge(git, repo, Files{}), "S-0003")
	if git.calls["version"] != 2 {
		t.Fatalf("git version ran %d times, want 2: once too old, once new enough", git.calls["version"])
	}
}

// I-0313: a root keeps at most trialKeep merges, dropping the one asked
// for least recently.
func TestTheMergeAskedForLeastRecentlyIsDropped(t *testing.T) {
	c := &mergeCache{roots: map[string]*rootMerges{}}
	pair := func(i int) [2]string { return [2]string{"m1", fmt.Sprint(i)} }
	for i := range trialKeep {
		c.keep("root", pair(i), []string{})
	}
	c.kept("root", pair(0))
	c.keep("root", pair(trialKeep), []string{})
	if _, ok := c.kept("root", pair(1)); ok {
		t.Fatal("the merge asked for least recently is still kept")
	}
	if n := len(c.roots["root"].pairs); n != trialKeep {
		t.Fatalf("root keeps %d merges, want %d", n, trialKeep)
	}
	for _, i := range []int{0, 2, trialKeep} {
		if _, ok := c.kept("root", pair(i)); !ok {
			t.Fatalf("merge %d, asked for since, was dropped", i)
		}
	}
}

// I-0313: a kept merge is the raw answer, so each judge's own files take
// what a rebase resolves out of it, and a judge with none answers it whole.
// Integration: it runs real git.
func TestTrialJudgeTakesItsOwnFilesOutOfAKeptMerge(t *testing.T) {
	repo, r := trialProject(t)
	commitFiles(t, repo.MainRoot, "chore: guide", "docs/guide.md")
	if got := judgeOnce(t, NewTrialJudge(r, repo, Files{}), "S-0002"); !slices.Equal(got, []string{"docs/guide.md"}) {
		t.Fatalf("first judge's conflicts: %v", got)
	}
	generated := NewTrialJudge(r, repo, Files{Generated: []GeneratedFile{{Path: "docs/guide.md"}}})
	if got := judgeOnce(t, generated, "S-0002"); len(got) != 0 {
		t.Fatalf("a judge generating docs/guide.md got conflicts: %v", got)
	}
	if got := judgeOnce(t, NewTrialJudge(r, repo, Files{}), "S-0002"); !slices.Equal(got, []string{"docs/guide.md"}) {
		t.Fatalf("a judge resolving nothing after it got conflicts: %v", got)
	}
	if r.merges != 1 {
		t.Fatalf("merge-tree ran %d times over the same commits, want 1", r.merges)
	}
}
