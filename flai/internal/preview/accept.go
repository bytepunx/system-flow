package preview

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Acceptance is what flai accept did, or under DryRun what it would do.
// Acceptance merges, moves to done, and archives; it computes no release, no
// tag, and no push (S-0087) — see flai release --pending for that, run when
// the operator chooses to publish what has accumulated on main.
type Acceptance struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	DryRun   bool     `json:"dry_run,omitempty"`
	Resumed  bool     `json:"resumed,omitempty"` // the item was already done but never archived
	Branch   string   `json:"branch,omitempty"`
	Merged   bool     `json:"merged"`
	Archived int      `json:"archived"`
	Blockers []string `json:"blockers,omitempty"` // what would stop acceptance before it changes anything
	// Uncommitted paths outside wip. The real run refuses them unless --yes
	// includes them in the acceptance commit; a dry run reports them so the
	// choice can be made before confirming (S-0051).
	Uncommitted []string `json:"uncommitted,omitempty"`
	// WorktreeUncommitted are the uncommitted paths in the story's worktree.
	// They block acceptance: the branch is merged as it is committed, and
	// the worktree is removed with it (S-0140).
	WorktreeUncommitted []string `json:"worktree_uncommitted,omitempty"`
	// The open stories told which paths the merge changed under their claim
	// (S-0132).
	Overlaps []itemedit.Overlap `json:"overlaps,omitempty"`
}

// Accept is what accepting a story or an epic would do, worked out before
// the first change so that an item is never left half accepted: an item
// that cannot be accepted at all is an error, and everything that would
// stop the acceptance midway is a blocker. by is who the probe's
// transitions are made by; log takes a warning about a worktree git cannot
// open.
func Accept(r execx.Runner, repo *workitem.Repo, it *workitem.Item, by string, now time.Time, log *slog.Logger) (*Acceptance, error) {
	// Without git (a project that is not a repository) acceptance is the
	// transition and the archive; with git it also merges and commits.
	useGit := storygit.InWorkTree(r, repo.MainRoot)
	if it.Type == workitem.Task {
		return nil, fmt.Errorf("%s is a task; accept its story instead", it.ID)
	}
	res := &Acceptance{ID: it.ID, Status: it.Status, DryRun: true}
	switch {
	case it.Status == workitem.Done && !it.Archived:
		// Done without acceptance: finish the job rather than refuse it.
		res.Resumed = true
	case it.Closed():
		return nil, fmt.Errorf("%s is already %s", it.ID, it.Status)
	case it.Type == workitem.Story && it.Status != workitem.Review:
		return nil, fmt.Errorf("%s is %s; a story is accepted from review", it.ID, it.Status)
	}
	if useGit {
		res.Uncommitted = storygit.DirtyOutsideWip(r, repo)
		if _, err := r.Run(repo.MainRoot, "git", "var", "GIT_COMMITTER_IDENT"); err != nil {
			res.Blockers = append(res.Blockers, "git has no committer identity here; set user.name and user.email (git config), or restart the dashboard with flai dashboard so it passes yours into the container")
		}
	} else if _, err := os.Stat(filepath.Join(repo.MainRoot, ".git")); err == nil {
		res.Blockers = append(res.Blockers, "this is a git repository but git cannot be run here; accept from a shell with flai accept "+it.ID)
	}
	// The workflow's own rules for done (open tasks, unticked criteria) are
	// checked on a copy before the branch is merged: a story that cannot be
	// done must not have its branch merged and its worktree removed first,
	// which is what happened until S-0041 found it from the review page.
	if !res.Resumed {
		if items, err := repo.List(false); err == nil {
			probe := *it
			probe.Transitions = append([]workitem.Transition(nil), it.Transitions...)
			for _, st := range StepsToDone(it.Status) {
				if _, err := repo.Move(&probe, st, workitem.MoveOptions{By: by, Now: now, Items: items}); err != nil {
					res.Blockers = append(res.Blockers, strings.TrimPrefix(err.Error(), "rule: "))
					break
				}
			}
		}
	}
	// A nature that is not accepted onto main (an experiment, ADR-0025) is
	// refused here, before the merge.
	if _, err := release.LevelFor(it); err != nil {
		res.Blockers = append(res.Blockers, err.Error())
	}
	// A story worktree git cannot open from here (I-0017): its links are
	// absolute host paths, and this process sees the repository somewhere
	// else. Say what to do instead of failing later with git's own error.
	if wt := repo.WorktreePath(it.ID); useGit && it.Type == workitem.Story {
		if _, err := os.Stat(wt); err == nil {
			if _, err := r.Run(wt, "git", "rev-parse", "--git-dir"); err != nil {
				log.Warn("story worktree cannot be opened by git", "component", "git", "worktree", rel(repo.MainRoot, wt), "err", err)
				res.Blockers = append(res.Blockers, fmt.Sprintf("git cannot open the story worktree %s from here: the paths git keeps for it do not exist in this environment, which happens when the dashboard sees the repository at a different path than the host does. Accept from a shell on the host with flai accept %s, or stop the dashboard and start it with flai dashboard, which mounts the repository at its host path", rel(repo.MainRoot, wt), it.ID))
			} else if dirty, err := repo.Uncommitted(it.ID); err == nil && len(dirty) > 0 {
				res.WorktreeUncommitted = dirty
				res.Blockers = append(res.Blockers, repo.UncommittedRule(it.ID, dirty, "accepting"))
			}
		}
	}
	if it.Type == workitem.Story && useGit && storygit.BranchExists(r, repo.MainRoot, storygit.Branch(it.ID)) {
		res.Branch = storygit.Branch(it.ID)
	}
	return res, nil
}

// StepsToDone is the walk from a state to done, one transition at a time.
func StepsToDone(status string) []string {
	path := []string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done}
	if status == workitem.Backlog {
		return path
	}
	for i, st := range path {
		if st == status {
			return path[i+1:]
		}
	}
	return []string{workitem.Done}
}
