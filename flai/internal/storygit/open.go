package storygit

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Where Open found a story's branch, as Opened.From says it, besides a
// remote's name.
const (
	// FromLocal is this clone's branch.
	FromLocal = "local"
	// FromMain is a new branch from the main branch.
	FromMain = "main"
)

// OpenOptions is what Open works with.
type OpenOptions struct {
	Runner execx.Runner
	Repo   *workitem.Repo
	Story  string
	// FromRemote fetches the story's branch from the remote when this clone
	// has none, as a story begun on another host is reopened (ADR-0064).
	FromRemote bool
	// RelativePaths reports whether to link a new worktree with relative
	// paths; it is asked only when Open adds one. Nil links an ordinary one.
	RelativePaths func() bool
	// Log takes the branch Open checked out; nil discards it.
	Log *slog.Logger
}

// Opened is the story branch Open checked out, all "" when git is not in
// use.
type Opened struct {
	Branch   string
	Worktree string // absolute
	// From is FromLocal, a remote's name, or FromMain; "" when the worktree
	// was there already.
	From string
}

// Open checks story/<id> out in its worktree under .flai-cache/worktrees
// (ADR-0019): this clone's branch when it has one, else, when FromRemote,
// the remote's, fetched (ADR-0064), else a new one from the main branch. A
// worktree that is there already is taken as it is. What a crash during a
// start leaves (I-0137) is made again when it holds no commits: a branch ref
// that resolves to nothing, and a folder git does not know as a worktree
// whose branch is missing or broken or whose files are all empty. Any other
// such folder is refused, as it may hold work.
func Open(o OpenOptions) (Opened, error) {
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	r, repo, id := o.Runner, o.Repo, o.Story
	if !InWorkTree(r, repo.MainRoot) {
		return Opened{}, nil
	}
	res := Opened{Branch: Branch(id), Worktree: repo.WorktreePath(id)}
	crashed := false
	if _, err := os.Stat(res.Worktree); err == nil {
		if IsWorktree(r, res.Worktree) {
			return res, nil
		}
		if err := removeCrashedWorktree(o, res); err != nil {
			return Opened{}, err
		}
		crashed = true
	}
	base, err := MainBranch(r, repo.MainRoot)
	if err != nil {
		return Opened{}, err
	}
	if err := os.MkdirAll(filepath.Dir(res.Worktree), 0o755); err != nil {
		return Opened{}, err
	}
	repaired, err := repairBrokenRef(o)
	if err != nil {
		return Opened{}, err
	}
	res.From = FromLocal
	if o.FromRemote || crashed || repaired {
		// a worktree removed without git leaves its registration, which
		// refuses the branch a new one; a crashed one's is all empty files,
		// which prune drops as an invalid gitdir file
		if _, err := r.Run(repo.MainRoot, "git", "worktree", "prune"); err != nil {
			return Opened{}, err
		}
	}
	if !BranchExists(r, repo.MainRoot, res.Branch) && o.FromRemote {
		if remote := Remote(r, repo.MainRoot); remote != "" {
			fetched, err := FetchBranch(r, repo.MainRoot, remote, res.Branch)
			if err != nil {
				return Opened{}, fmt.Errorf("%w; run flai stream open %s again when %s answers, or create the branch from %s with git branch %s %s first", err, id, remote, base, res.Branch, base)
			}
			if fetched {
				res.From = remote
			}
		}
	}
	// the ref is synced to disk before git returns, so a crash after it
	// leaves the branch whole (I-0137)
	args := []string{"-c", "core.fsync=reference", "worktree", "add", "--quiet"}
	if o.RelativePaths != nil && o.RelativePaths() {
		args = append(args, "--relative-paths")
	}
	if BranchExists(r, repo.MainRoot, res.Branch) {
		args = append(args, res.Worktree, res.Branch)
	} else {
		res.From = FromMain
		args = append(args, "-b", res.Branch, res.Worktree, base)
	}
	if _, err := r.Run(repo.MainRoot, "git", args...); err != nil {
		return Opened{}, err
	}
	o.Log.Info("story branch opened", "component", "git", "branch", res.Branch, "worktree", relPath(repo.MainRoot, res.Worktree), "base", base, "from", res.From)
	return res, nil
}

// removeCrashedWorktree removes a story's worktree folder that git does not
// know as one, as a crash during its start leaves it (I-0137), when it holds
// no work: the story's branch is missing or broken, so it holds no commits,
// or every file in the folder is empty, as a crash after git made the branch
// leaves the checkout. Any other such folder may hold work, and is refused.
func removeCrashedWorktree(o OpenOptions, res Opened) error {
	rel := relPath(o.Repo.MainRoot, res.Worktree)
	if BranchExists(o.Runner, o.Repo.MainRoot, res.Branch) && !allFilesEmpty(res.Worktree) {
		return fmt.Errorf("%s is there but git does not know it as a worktree, and %s holds commits, so it was left as it is: check that it holds no work you need, remove it with rm -rf %s, and run flai stream open %s again", rel, res.Branch, res.Worktree, o.Story)
	}
	if err := os.RemoveAll(res.Worktree); err != nil {
		return fmt.Errorf("could not remove %s, a folder git does not know as a worktree whose branch %s holds no commits: %w; remove it with rm -rf %s and run flai stream open %s again", rel, res.Branch, err, res.Worktree, o.Story)
	}
	o.Log.Warn("crashed story worktree removed", "component", "git", "worktree", rel, "branch", res.Branch)
	return nil
}

// allFilesEmpty reports whether every regular file under dir is empty, so the
// folder holds no work; one it cannot read through is not.
func allFilesEmpty(dir string) bool {
	held := false
	err := filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err != nil || !d.Type().IsRegular() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Size() > 0 {
			held = true
			return filepath.SkipAll
		}
		return nil
	})
	return err == nil && !held
}

// repairBrokenRef removes the story's branch ref when it resolves to nothing
// (brokenStoryRef), which git refuses to delete, replace, or make again, and
// its reflog when that is empty too; it reports whether it removed the ref.
// A reflog with entries is kept: they say where the branch was, and git
// appends to it when the branch is made again.
func repairBrokenRef(o OpenOptions) (bool, error) {
	common, ref, err := brokenStoryRef(o.Runner, o.Repo.MainRoot, o.Story)
	if err != nil || ref == "" {
		return false, err
	}
	refFile := filepath.Join(common, ref)
	held, err := os.ReadFile(refFile)
	if err != nil {
		return false, fmt.Errorf("could not read %s, the broken ref of %s: %w", refFile, Branch(o.Story), err)
	}
	if err := os.Remove(refFile); err != nil {
		return false, fmt.Errorf("could not remove %s, the broken ref of %s, which holds no commits: %w; remove it by hand and run flai stream open %s again", refFile, Branch(o.Story), err, o.Story)
	}
	reflog := "none"
	logFile := filepath.Join(common, "logs", ref)
	if info, err := os.Stat(logFile); err == nil {
		reflog = "kept"
		if info.Size() == 0 {
			if err := os.Remove(logFile); err != nil {
				return false, fmt.Errorf("could not remove %s, the empty reflog of %s: %w; remove it by hand and run flai stream open %s again", logFile, Branch(o.Story), err, o.Story)
			}
			reflog = "removed"
		}
	}
	o.Log.Warn("broken story branch ref removed", "component", "git", "branch", Branch(o.Story),
		"ref", relPath(o.Repo.MainRoot, refFile), "held", strings.TrimSpace(string(held)), "reflog", reflog)
	return true, nil
}
