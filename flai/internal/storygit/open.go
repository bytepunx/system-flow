package storygit

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
// worktree that is there already is taken as it is.
func Open(o OpenOptions) (Opened, error) {
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	r, repo, id := o.Runner, o.Repo, o.Story
	if !InWorkTree(r, repo.MainRoot) {
		return Opened{}, nil
	}
	res := Opened{Branch: Branch(id), Worktree: repo.WorktreePath(id)}
	if _, err := os.Stat(res.Worktree); err == nil {
		return res, nil
	}
	base, err := MainBranch(r, repo.MainRoot)
	if err != nil {
		return Opened{}, err
	}
	if err := os.MkdirAll(filepath.Dir(res.Worktree), 0o755); err != nil {
		return Opened{}, err
	}
	res.From = FromLocal
	if o.FromRemote {
		// a worktree removed without git leaves its registration, which
		// refuses the branch a new one
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
	args := []string{"worktree", "add", "--quiet"}
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
