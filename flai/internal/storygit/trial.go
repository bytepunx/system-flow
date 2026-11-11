package storygit

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	j := &trialJudge{r: r, root: repo.MainRoot, files: files}
	return j.judge
}

// trialJudge is NewTrialJudge's state: the runner, the main checkout, and
// what a rebase resolves there. It keeps no heads between calls, since a
// judge may outlive them, as the MCP server's lasts its session.
type trialJudge struct {
	r     execx.Runner
	root  string
	files Files
}

// judge answers the paths that conflict when holder's branch merges at held's
// base: held's own branch when it has one, else the main branch (S-0551).
func (j *trialJudge) judge(held, holder string) ([]string, error) {
	refs, err := readRefs(j.r, j.root)
	if err != nil {
		return nil, err
	}
	base, err := j.base(refs, held)
	if err != nil {
		return nil, err
	}
	head, ok := refs.heads[Branch(holder)]
	if !ok {
		return nil, fmt.Errorf("%s has no branch %s yet to trial-merge at %s's base", holder, Branch(holder), held)
	}
	key := [2]string{base, head}
	conflicts, kept := trialMerges.kept(j.root, key)
	if !kept {
		if err := gitReady(j.r); err != nil {
			return nil, err
		}
		conflicts, err = TrialMerge(j.r, j.root, base, head)
		if err != nil {
			return nil, fmt.Errorf("trial merge of %s at %s's base in %s: %w", Branch(holder), held, j.root, err)
		}
		trialMerges.keep(j.root, key, conflicts)
	}
	return withoutResolved(conflicts, j.files), nil
}

// base is the commit held's branch is at when it has one, as a story sent
// back does, else the main branch's.
func (j *trialJudge) base(refs refs, held string) (string, error) {
	if sha, ok := refs.heads[Branch(held)]; ok {
		return sha, nil
	}
	main := refs.main
	if main == "" {
		var err error
		if main, err = MainBranch(j.r, j.root); err != nil {
			return "", fmt.Errorf("read the main branch in %s for %s's base: %w", j.root, held, err)
		}
	}
	sha, ok := refs.heads[main]
	if !ok {
		return "", fmt.Errorf("main branch %s in %s resolves to no commit for %s's base", main, j.root, held)
	}
	return sha, nil
}

// refs is the commit each local branch is at, by branch name, and the
// branch the main checkout has checked out, "" when it is none of them.
type refs struct {
	heads map[string]string
	main  string
}

// refsFormat is one local branch a line, its fields split by tabs: its
// name, * when the checkout's HEAD is at it, the type and name of the object
// an annotated tag points at, empty for any other, and its own object's.
// Each line starts with the name and ends with an object's name, so a runner
// trimming the output, as execx.System does, cuts no field.
const refsFormat = "--format=%(refname)%09%(HEAD)%09%(*objecttype)%09%(*objectname)%09%(objecttype)%09%(objectname)"

// readRefs reads root's local branches with one for-each-ref. A branch at
// an annotated tag is at the tag's object, and one at anything but a
// commit, after that peel, is left out, as rev-parse's ^{commit} finds no
// commit for it. When HEAD is at no branch read, main is "", for MainBranch
// to say why.
func readRefs(r execx.Runner, root string) (refs, error) {
	out, err := r.Run(root, "git", "for-each-ref", refsFormat, "refs/heads/")
	if err != nil {
		return refs{}, fmt.Errorf("read the heads of the local branches in %s to trial-merge: %w", root, err)
	}
	read := refs{heads: map[string]string{}}
	for line := range strings.SplitSeq(out, "\n") {
		f := strings.Split(line, "\t")
		if len(f) != 6 {
			continue
		}
		branch, kind, sha := strings.TrimPrefix(f[0], "refs/heads/"), f[4], f[5]
		if f[2] != "" {
			kind, sha = f[2], f[3]
		}
		if kind == "commit" {
			read.heads[branch] = sha
		}
		if f[1] == "*" {
			read.main = branch
		}
	}
	return read, nil
}

// trialKeep is how many merges a root keeps: past it, the merge asked for
// least recently, of heads that have moved or of a story that no longer
// overlaps, is dropped. It bounds a root by count rather than by the judges
// made since a merge was asked for, as a judge kept for a session, the MCP
// server's, makes no new judge however often it is asked.
const trialKeep = 256

// trialMerges is the raw TrialMerge answers kept for the process, by main
// checkout root and the pair of commits merged. A merge of two given commits
// answers the same paths every time, so a kept answer needs no other check;
// what a rebase resolves is taken out after, by each judge's own files.
var trialMerges = &mergeCache{roots: map[string]*rootMerges{}}

// mergeCache is the TrialMerge answers of every root, guarded by mu, since
// the judges of many projects, and many of one, ask at once.
type mergeCache struct {
	mu    sync.Mutex
	roots map[string]*rootMerges
}

// rootMerges is one root's kept merges and how many times they were asked
// for or kept, which orders them by when they last were.
type rootMerges struct {
	asks  int
	pairs map[[2]string]keptMerge
}

// keptMerge is a TrialMerge answer and its root's ask count when it was
// last asked for or kept.
type keptMerge struct {
	paths []string
	asked int
}

// kept is the answer kept for the merge of key in root, and false when none
// is kept.
func (c *mergeCache) kept(root string, key [2]string) ([]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rm := c.roots[root]
	if rm == nil {
		return nil, false
	}
	m, ok := rm.pairs[key]
	if ok {
		rm.asks++
		m.asked = rm.asks
		rm.pairs[key] = m
	}
	return m.paths, ok
}

// keep keeps paths as the answer for the merge of key in root, dropping the
// merge asked for least recently once root keeps more than trialKeep.
func (c *mergeCache) keep(root string, key [2]string, paths []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rm := c.roots[root]
	if rm == nil {
		rm = &rootMerges{pairs: map[[2]string]keptMerge{}}
		c.roots[root] = rm
	}
	rm.asks++
	rm.pairs[key] = keptMerge{paths: paths, asked: rm.asks}
	if len(rm.pairs) <= trialKeep {
		return
	}
	var oldest [2]string
	least := rm.asks
	for k, m := range rm.pairs {
		if m.asked < least {
			oldest, least = k, m.asked
		}
	}
	delete(rm.pairs, oldest)
}

// gitPassed is the runner whose git last had merge-tree --write-tree, so a
// judge asked of the same runner, the system's in flai, does not ask git's
// version again. It is kept by runner rather than for the process alone
// because another runner may answer for another git, as tests' runners do,
// and only a runner whose value can be compared is kept. An error is not
// kept, so a git installed later is seen.
var gitPassed struct {
	mu sync.Mutex
	r  execx.Runner
}

// gitReady is nil when r's git can trial-merge, asking its version only when
// r is not the runner whose git last could.
func gitReady(r execx.Runner) error {
	canKeep := reflect.ValueOf(r).Comparable()
	gitPassed.mu.Lock()
	passed := canKeep && gitPassed.r == r
	gitPassed.mu.Unlock()
	if passed {
		return nil
	}
	if err := gitTooOld(r); err != nil {
		return err
	}
	if canKeep {
		gitPassed.mu.Lock()
		gitPassed.r = r
		gitPassed.mu.Unlock()
	}
	return nil
}

// gitTooOld is an error naming git's version when it has no merge-tree
// --write-tree, or cannot be read.
func gitTooOld(r execx.Runner) error {
	v, err := gitver.Installed(r)
	if err != nil {
		return fmt.Errorf("read git's version to trial-merge: %w; check that git is on PATH", err)
	}
	if !v.AtLeast(gitver.MergeTree) {
		return fmt.Errorf("git %s is older than %s, the first with merge-tree --write-tree; upgrade git to trial-merge", v, gitver.MergeTree)
	}
	return nil
}
