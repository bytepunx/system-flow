package storygit

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// NewTrialJudge is a judge for workitem holds: given a ready story held and a
// story in progress holder, it trial-merges holder's branch at held's base and
// answers the paths that conflict, none when clean, or an error when the trial
// cannot run. files is what a rebase resolves itself in repo, whose generated
// and merged files it leaves out of the conflicts as sync's trial merge does.
// It is nil without a runner or when repo's main checkout is no git
// repository: no story there has a branch to trial-merge, so an overlap holds
// as it did before (ADR-0230).
func NewTrialJudge(r execx.Runner, repo *workitem.Repo, files Files) func(held, holder string) ([]string, error) {
	if _, err := os.Stat(filepath.Join(repo.MainRoot, ".git")); r == nil || err != nil {
		return nil
	}
	j := &trialJudge{r: r, root: repo.MainRoot, files: files, answers: map[[2]string][]string{}}
	j.git = sync.OnceValue(j.gitTooOld)
	return j.judge
}

// trialJudge is NewTrialJudge's state: the git version check, run once, and
// the answers by the pair of commits trial-merged.
type trialJudge struct {
	r     execx.Runner
	root  string
	files Files
	git   func() error

	mu      sync.Mutex
	answers map[[2]string][]string
}

// judge answers the paths that conflict when holder's branch merges at held's
// base: held's own branch when it has one, else the main branch (S-0551).
func (j *trialJudge) judge(held, holder string) ([]string, error) {
	base, err := j.base(held)
	if err != nil {
		return nil, err
	}
	head, ok := j.commit(Branch(holder))
	if !ok {
		return nil, fmt.Errorf("%s has no branch %s yet to trial-merge at %s's base", holder, Branch(holder), held)
	}
	key := [2]string{base, head}
	j.mu.Lock()
	answer, seen := j.answers[key]
	j.mu.Unlock()
	if seen {
		return slices.Clone(answer), nil
	}
	if err := j.git(); err != nil {
		return nil, err
	}
	conflicts, err := TrialMerge(j.r, j.root, base, head)
	if err != nil {
		return nil, fmt.Errorf("trial merge of %s at %s's base in %s: %w", Branch(holder), held, j.root, err)
	}
	conflicts = withoutResolved(conflicts, j.files)
	j.mu.Lock()
	j.answers[key] = conflicts
	j.mu.Unlock()
	return slices.Clone(conflicts), nil
}

// base is the commit held's branch is at when it has one, as a story sent
// back does, else the main branch's.
func (j *trialJudge) base(held string) (string, error) {
	if sha, ok := j.commit(Branch(held)); ok {
		return sha, nil
	}
	main, err := MainBranch(j.r, j.root)
	if err != nil {
		return "", fmt.Errorf("read the main branch in %s for %s's base: %w", j.root, held, err)
	}
	sha, ok := j.commit(main)
	if !ok {
		return "", fmt.Errorf("main branch %s in %s resolves to no commit for %s's base", main, j.root, held)
	}
	return sha, nil
}

// commit is the commit the local branch is at, and false when it has none.
func (j *trialJudge) commit(branch string) (string, bool) {
	out, err := j.r.Run(j.root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch+"^{commit}")
	return out, err == nil && out != ""
}

// gitTooOld is an error naming git's version when it has no merge-tree
// --write-tree, or cannot be read.
func (j *trialJudge) gitTooOld() error {
	v, err := gitver.Installed(j.r)
	if err != nil {
		return fmt.Errorf("read git's version to trial-merge: %w; check that git is on PATH", err)
	}
	if !v.AtLeast(gitver.MergeTree) {
		return fmt.Errorf("git %s is older than %s, the first with merge-tree --write-tree; upgrade git to trial-merge", v, gitver.MergeTree)
	}
	return nil
}
