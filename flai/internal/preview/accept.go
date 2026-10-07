package preview

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conflictmark"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/protected"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Acceptance is what flai accept did, or under DryRun what it would do.
// Acceptance merges, moves to done, and archives; it computes no release, no
// tag, and no push (S-0087) — see flai release --pending for that, run when
// the operator chooses to publish what has accumulated on main.
type Acceptance struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	By       string   `json:"by,omitempty"` // who accepted it, as its done transition records; set by the real run
	DryRun   bool     `json:"dry_run,omitempty"`
	Resumed  bool     `json:"resumed,omitempty"` // the item was already done: never archived, or archived but never committed (CommitOnly)
	Branch   string   `json:"branch,omitempty"`
	Merged   bool     `json:"merged"`
	Archived int      `json:"archived"`
	Blockers []string `json:"blockers,omitempty"` // what would stop acceptance before it changes anything
	// CommitOnly says the item is done and archived but its acceptance
	// commit is missing: only the commit and the notices are left (I-0100).
	CommitOnly bool `json:"commit_only,omitempty"`
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
	// Protected are the files the story's branch changes on a path Claude
	// Code protects, and OperatorOnly says, when there are any, that only the
	// operator accepts the story and who that is (ADR-0106).
	Protected    []string `json:"protected,omitempty"`
	OperatorOnly string   `json:"operator_only,omitempty"`
	// ResolvedThreads are the threads still open or answered on what the
	// acceptance archives, which it resolves (I-0073), or under DryRun would.
	ResolvedThreads []string `json:"resolved_threads,omitempty"`
}

// Accept is what accepting a story or an epic would do, worked out before
// the first change so that an item is never left half accepted: an item
// that cannot be accepted at all is an error, and everything that would
// stop the acceptance midway is a blocker. by is who accepts: the probe's
// transitions are made by them, and a story whose branch changes a path
// Claude Code protects is refused unless they are its operator (ADR-0106);
// log takes a warning about a worktree git cannot open.
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
	res := &Acceptance{ID: it.ID, Status: it.Status, DryRun: true}
	switch {
	case it.Status == workitem.Done && !it.Archived:
		// Done without acceptance: finish the job rather than refuse it.
		res.Resumed = true
	case it.Status == workitem.Done && useGit && uncommitted(r, repo.MainRoot, it.Path):
		// Done and archived, its archived file not committed: its acceptance
		// commit failed (I-0100), so the commit is finished. The item's own
		// file decides, not the wip folder, where flai leaves other changes
		// uncommitted that an old item's subject must not take in.
		res.Resumed, res.CommitOnly = true, true
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
	if res.CommitOnly {
		// The merge, the transition, and the archive are behind it, so
		// nothing they check applies; the commit names the epic it archived.
		res.Epic = archivedEpic(r, repo, it)
		threadsToResolve(repo, it, res)
		return res, nil
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
		// a change to a path Claude Code protects takes effect for every
		// agent once merged, so only the operator accepts it (ADR-0106)
		if diff, err := storygit.StoryDiff(r, repo, it.ID); err != nil {
			res.Blockers = append(res.Blockers, fmt.Sprintf("cannot read what %s changes to check it for paths Claude Code protects: %v", res.Branch, err))
		} else if res.Protected = protected.Changed(changedPaths(diff)); len(res.Protected) > 0 {
			who := operatorsOf(repo, it)
			res.OperatorOnly = fmt.Sprintf("only the operator (%s) accepts %s: its branch changes paths Claude Code protects", who, it.ID)
			if !opts.Orchestrator && !isOperator(repo, it, by) {
				res.Blockers = append(res.Blockers, fmt.Sprintf("only the operator (%s) accepts %s, not %s: its branch changes paths Claude Code protects, %s; the operator reads them in the diff and accepts it (ADR-0106)", who, it.ID, by, strings.Join(res.Protected, ", ")))
			}
		}
	}
	if opts.Orchestrator {
		res.Evidence = opts.Evidence
		res.OrchestratorBlockers, res.Verified = orchestratorBlockers(r, repo, it, opts, res.Protected)
		for _, b := range res.OrchestratorBlockers {
			res.Blockers = append(res.Blockers, b.Message)
		}
	}
	threadsToResolve(repo, it, res)
	return res, nil
}

// threadsToResolve sets on res the threads still open or answered on what
// accepting it archives, which the acceptance resolves (I-0073). A thread
// that cannot be read is a blocker: the acceptance would stop on it after
// archiving.
func threadsToResolve(repo *workitem.Repo, it *workitem.Item, res *Acceptance) {
	all, err := repo.List(true)
	if err != nil {
		res.Blockers = append(res.Blockers, fmt.Sprintf("cannot list the items to find the threads accepting %s resolves: %v", it.ID, err))
		return
	}
	open, err := threads.OnItems(repo, ThreadItems(all, it, res.Epic))
	if err != nil {
		res.Blockers = append(res.Blockers, fmt.Sprintf("cannot read the threads on what accepting %s archives, which it resolves: %v; fix the thread file it names", it.ID, err))
		return
	}
	for _, th := range open {
		res.ResolvedThreads = append(res.ResolvedThreads, th.ID)
	}
}

// ThreadItems are the items whose threads accepting it resolves: what
// ArchivedWith names, and each story's tasks. items include archived ones, so
// that an acceptance finishing its commit names what it archived before.
func ThreadItems(items []*workitem.Item, it *workitem.Item, epic *workitem.Followed) []string {
	var ids []string
	for _, id := range ArchivedWith(items, it.ID, epic) {
		ids = append(ids, id)
		if workitem.TypeOfID(id) == workitem.Story {
			for _, c := range workitem.Children(items, id) {
				ids = append(ids, c.ID)
			}
		}
	}
	return ids
}

// uncommitted reports whether path has changes the checkout at root has not
// committed, staged or not; false when git cannot say.
func uncommitted(r execx.Runner, root, path string) bool {
	st, err := r.Run(root, "git", "status", "--porcelain", "--", path)
	return err == nil && strings.TrimSpace(st) != ""
}

// archivedEpic is the walk the epic of story took with it to done, when the
// epic's archived file is not committed either: the acceptance whose commit
// is missing archived it with the story. Nil otherwise.
func archivedEpic(r execx.Runner, repo *workitem.Repo, story *workitem.Item) *workitem.Followed {
	if story.Type != workitem.Story || story.Parent == "" {
		return nil
	}
	epic, err := repo.Get(story.Parent)
	if err != nil || epic.Type != workitem.Epic || epic.Status != workitem.Done || !epic.Archived || !uncommitted(r, repo.MainRoot, epic.Path) {
		return nil
	}
	return &workitem.Followed{ID: epic.ID, Type: epic.Type, Title: epic.Title, From: walkedFrom(epic.Transitions), To: workitem.Done, Story: story.ID}
}

// walkedFrom is the state an item was in before its last walk: the
// transitions made at the time of the last one are one walk, as a follow
// makes them; backlog when there is nothing before them.
func walkedFrom(ts []workitem.Transition) string {
	i := len(ts) - 1
	for i >= 0 && ts[i].At == ts[len(ts)-1].At {
		i--
	}
	if i < 0 {
		return workitem.Backlog
	}
	return ts[i].To
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

// operators are who accept a story whose branch changes a path Claude Code
// protects (ADR-0106): its owner and the project's owner, those named, each
// once, as ADR-0097 names who answers a permission thread. None means anyone
// but the orchestrator.
func operators(repo *workitem.Repo, it *workitem.Item) []string {
	var who []string
	for _, name := range []string{it.Owner, repo.Manifest.Owner} {
		if name != "" && !slices.Contains(who, name) {
			who = append(who, name)
		}
	}
	return who
}

// isOperator says whether by is the operator of it: never the orchestrator.
func isOperator(repo *workitem.Repo, it *workitem.Item, by string) bool {
	who := operators(repo, it)
	return !workitem.IsOrchestrator(by) && (len(who) == 0 || slices.Contains(who, by))
}

// operatorsOf names the operator of it for a sentence: the names, or who
// may accept when it and the project name no owner.
func operatorsOf(repo *workitem.Repo, it *workitem.Item) string {
	if who := operators(repo, it); len(who) > 0 {
		return strings.Join(who, " or ")
	}
	return fmt.Sprintf("anyone but the orchestrator, as neither %s nor the project names an owner", it.ID)
}

// changedPaths are the paths a branch's diff changes, a renamed file's old
// path beside its new one.
func changedPaths(diff *storygit.Diff) []string {
	var changed []string
	for _, f := range diff.Files {
		changed = append(changed, f.Path)
		if f.OldPath != "" {
			changed = append(changed, f.OldPath)
		}
	}
	return changed
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
