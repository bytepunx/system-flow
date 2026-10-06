package preview

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conflictmark"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
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
	By       string   `json:"by,omitempty"` // who accepts it, as its done transition records
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
	// Epic is the walk the story's epic takes with it: done, and archived
	// with it, when the story is the epic's last open one (S-0200).
	Epic *workitem.Followed `json:"epic,omitempty"`
	// OrchestratorBlockers are the conditions of an acceptance by the
	// orchestrator that fail (ADR-0093), each also in Blockers.
	OrchestratorBlockers []Blocker `json:"orchestrator_blockers,omitempty"`
	// Verified is the commit the orchestrator's verifier passed, resolved to
	// its full name, when it names one.
	Verified string `json:"verified,omitempty"`
	// Evidence is the orchestrator's evidence, when given.
	Evidence *Evidence `json:"evidence,omitempty"`
}

// Accept is what accepting a story or an epic would do, worked out before
// the first change so that an item is never left half accepted: an item
// that cannot be accepted at all is an error, and everything that would
// stop the acceptance midway is a blocker. by is who the probe's
// transitions are made by; log takes a warning about a worktree git cannot
// open.
func Accept(r execx.Runner, repo *workitem.Repo, it *workitem.Item, by string, now time.Time, log *slog.Logger) (*Acceptance, error) {
	return AcceptWith(r, repo, it, by, now, log, AcceptOptions{})
}

// AcceptWith is Accept with what the acceptance adds: under
// opts.Orchestrator, a blocker for each condition of ADR-0093 that fails,
// beside every blocker Accept reports.
func AcceptWith(r execx.Runner, repo *workitem.Repo, it *workitem.Item, by string, now time.Time, log *slog.Logger, opts AcceptOptions) (*Acceptance, error) {
	// Without git (a project that is not a repository) acceptance is the
	// transition and the archive; with git it also merges and commits.
	useGit := storygit.InWorkTree(r, repo.MainRoot)
	if it.Type == workitem.Task {
		return nil, fmt.Errorf("%s is a task; accept its story instead", it.ID)
	}
	res := &Acceptance{ID: it.ID, Status: it.Status, By: by, DryRun: true}
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
			walked := true
			for _, st := range StepsToDone(it.Status) {
				if _, err := repo.Move(&probe, st, workitem.MoveOptions{By: by, Now: now, Items: items}); err != nil {
					res.Blockers = append(res.Blockers, strings.TrimPrefix(err.Error(), "rule: "))
					walked = false
					break
				}
			}
			if walked && it.Type == workitem.Story {
				res.Blockers = append(res.Blockers, epicAccepted(repo, res, &probe, it.Status, by, now)...)
			}
		}
	}
	// A nature with no release rule is refused here, before the merge.
	if _, err := release.LevelFor(it); err != nil {
		res.Blockers = append(res.Blockers, err.Error())
	}
	// An experiment is accepted only with its results document (ADR-0066).
	if experiment.Needs(it) {
		if b := resultsBlocker(r, repo, it, useGit); b != "" {
			res.Blockers = append(res.Blockers, b)
		}
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
		// a conflict git-added unresolved is refused by the merge, so the
		// preview names it before Accept is pressed (S-0276)
		if b, err := ConflictMarkers(r, repo, it.ID); err != nil {
			res.Blockers = append(res.Blockers, fmt.Sprintf("cannot read %s for merge conflict markers: %v", res.Branch, err))
		} else if b != "" {
			res.Blockers = append(res.Blockers, b)
		}
	}
	if opts.Orchestrator {
		res.Evidence = opts.Evidence
		res.OrchestratorBlockers, res.Verified = orchestratorBlockers(r, repo, it, opts)
		for _, b := range res.OrchestratorBlockers {
			res.Blockers = append(res.Blockers, b.Message)
		}
	}
	return res, nil
}

// ConflictMarkers is why story id's branch cannot be accepted for the merge
// conflict markers it carries, naming each path and line of the files it
// adds or changes against the main branch, or "" when it carries none. The
// acceptance refuses by it and its preview reports by it, so the two cannot
// disagree (S-0253, S-0276, I-0066).
func ConflictMarkers(r execx.Runner, repo *workitem.Repo, id string) (string, error) {
	branch := storygit.Branch(id)
	base, err := storygit.MainBranch(r, repo.MainRoot)
	if err != nil {
		return "", err
	}
	found, err := conflictmark.Branch(r, repo.MainRoot, base, branch)
	if err != nil || len(found) == 0 {
		return "", err
	}
	wt := repo.WorktreePath(id)
	where := "in " + rel(repo.MainRoot, wt)
	if _, err := os.Stat(wt); err != nil {
		where = fmt.Sprintf("on %s (open its worktree with flai stream open %s)", branch, id)
	}
	return fmt.Sprintf("%s carries merge conflict markers at %s; resolve each conflict %s keeping what both sides meant, remove the markers, commit, and accept again", branch, workitem.Shorten(found, 20), where), nil
}

// epicAccepted works out what the story's epic would do were the story,
// walked to done as probe from from, accepted, sets it on res, and returns
// what would stop it: the epic's walk, or the archive of the epic with
// everything of it still on the board.
func epicAccepted(repo *workitem.Repo, res *Acceptance, probe *workitem.Item, from, by string, now time.Time) []string {
	all, followed, err := EpicWalk(repo, probe, from, by, now, true)
	if err != nil {
		return []string{err.Error()}
	}
	res.Epic = followed
	if followed == nil || followed.To != workitem.Done {
		return nil
	}
	var open []*workitem.Item
	for _, x := range all {
		if !x.Archived {
			open = append(open, x)
		}
	}
	if _, err := repo.PlanArchive(open, ArchivedWith(open, probe.ID, followed)); err != nil {
		return []string{fmt.Sprintf("%s would follow %s to done but cannot be archived with it: %v", followed.ID, probe.ID, err)}
	}
	return nil
}

// EpicWalk is what the epic of probe, a story moved from from, would do with
// it, worked out on copies so that nothing real changes (S-0200). It returns
// the items, archived ones included, with probe and the epic's copy in them.
func EpicWalk(repo *workitem.Repo, probe *workitem.Item, from, by string, now time.Time, accept bool) ([]*workitem.Item, *workitem.Followed, error) {
	all, err := repo.List(true)
	if err != nil {
		return nil, nil, err
	}
	for i, x := range all {
		switch {
		case x.ID == probe.ID:
			all[i] = probe
		case x.ID == probe.Parent && x.Type == workitem.Epic:
			epic := *x
			epic.Transitions = append([]workitem.Transition(nil), x.Transitions...)
			all[i] = &epic
		}
	}
	_, followed, err := repo.Follow(all, probe, from, by, now, accept)
	return all, followed, err
}

// ArchivedWith is what accepting story archives: the story, and when its
// epic followed it to done, the epic's other stories still on the board
// (cancelled ones) and then the epic, so that nothing of it is left there.
func ArchivedWith(items []*workitem.Item, story string, epic *workitem.Followed) []string {
	ids := []string{story}
	if epic == nil || epic.To != workitem.Done {
		return ids
	}
	for _, c := range workitem.Children(items, epic.ID) {
		if c.ID != story && !c.Archived {
			ids = append(ids, c.ID)
		}
	}
	return append(ids, epic.ID)
}

// resultsBlocker is why an experiment story cannot be accepted for want of
// its results document, or "" when it has one: on its branch, which is what
// acceptance merges, or in the main checkout when there is no branch.
func resultsBlocker(r execx.Runner, repo *workitem.Repo, it *workitem.Item, useGit bool) string {
	dir := experiment.Dir(repo.Manifest)
	if branch := storygit.Branch(it.ID); useGit && storygit.BranchExists(r, repo.MainRoot, branch) {
		out, err := r.Run(repo.MainRoot, "git", "ls-tree", "--name-only", branch, dir+"/")
		if err != nil {
			return fmt.Sprintf("cannot list %s on %s: %v", dir, branch, err)
		}
		var names []string
		for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
			names = append(names, path.Base(l))
		}
		if _, ok := experiment.Find(names, it.ID); !ok {
			return experiment.Missing(repo.Manifest, it, "its branch "+branch)
		}
		return ""
	}
	entries, _ := os.ReadDir(filepath.Join(repo.MainRoot, filepath.FromSlash(dir)))
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if _, ok := experiment.Find(names, it.ID); !ok {
		return experiment.Missing(repo.Manifest, it, "main")
	}
	return ""
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
