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
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
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
	// its full name, when it names one; for a landing, the story's rebased
	// head in the merge queue batch that verified it (ADR-0176).
	Verified string `json:"verified,omitempty"`
	// Batch is the merge queue batch the story lands from, whose
	// verification the acceptance takes, when it is the queue's landing
	// (ADR-0176).
	Batch string `json:"batch,omitempty"`
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
	// UndoesMain are the paths the story's record says its rebased branch
	// still undoes of the main branch's changes (ADR-0155), each not kept
	// also in Blockers.
	UndoesMain []UndonePath `json:"undoes_main,omitempty"`
	// ForgeBlockers are the forge blockers under pull-request acceptance
	// (ADR-0169), each also in Blockers.
	ForgeBlockers []Blocker `json:"forge_blockers,omitempty"`
	// PullRequest is the URL of the story's pull request as flai serve kept
	// it, when it kept one.
	PullRequest string `json:"pull_request,omitempty"`
	// Remote is how the story's pull request is merged through the remote
	// under pull-request acceptance, in place of the local merge (ADR-0170);
	// nil under local acceptance.
	Remote *RemoteMerge `json:"remote,omitempty"`
}

// RemoteMerge is the merge of a story's pull request on the forge under
// pull-request acceptance (ADR-0170): what the preview read of it from the
// state flai serve kept, and, once the real run has made it, what the merge
// came to.
type RemoteMerge struct {
	// PullRequest is the pull request's number as flai serve kept it; 0 while
	// it kept none.
	PullRequest int `json:"pull_request,omitempty"`
	// Head is the head the forge gate passed at, the kept status's head
	// commit: a pull request whose head has moved since is not merged.
	Head string `json:"head,omitempty"`
	// Method is acceptance.merge_method, or its default.
	Method forge.MergeMethod `json:"method"`
	// AlreadyMerged is a pull request flai serve read as merged on the forge
	// already, which the acceptance accepts with no second merge.
	AlreadyMerged bool `json:"already_merged,omitempty"`
	// ThroughServe is a merge flai accept hands to flai serve on the host,
	// since the forge's token cannot be read in its own environment.
	ThroughServe bool `json:"through_serve,omitempty"`
	// Outcome is what the real run's merge came to: merged, or
	// already-merged; SHA is the commit the merge left at the head of
	// origin's main, empty for a pull request merged before.
	Outcome string `json:"outcome,omitempty"`
	SHA     string `json:"sha,omitempty"`
	// BranchDeleted is whether the story's branch was deleted on origin, and
	// BranchReason why not when it was not; the merge stands either way.
	BranchDeleted bool   `json:"branch_deleted,omitempty"`
	BranchReason  string `json:"branch_reason,omitempty"`
}

// Codes of the forge blockers an acceptance under pull-request acceptance
// adds (ADR-0169).
const (
	BlockPullRequestMissing = "pull_request_missing" // flai serve has kept no pull request for the story
	BlockForgeStale         = "forge_stale"          // the kept state was read at a commit other than the branch head
	BlockCheckFailing       = "check_failing"        // a required check failed
	BlockCheckPending       = "check_pending"        // a required check has not finished
	BlockChangesRequested   = "changes_requested"    // a reviewer's latest decisive review requests changes
	BlockApprovalsMissing   = "approvals_missing"    // fewer reviewers approve than the base branch requires
)

// UndonePath is a path the story's branch undoes of the main branch's change
// (ADR-0155): main's lines the branch removes, and the operator's reason when
// the acceptance keeps the removal.
type UndonePath struct {
	storygit.Undone
	KeptReason string `json:"kept_reason,omitempty"`
}

// KeptRemoval is a path whose removal of the main branch's lines the operator
// keeps at acceptance, and why (ADR-0155).
type KeptRemoval struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
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
// beside every blocker Accept reports, and no blocker for a path the branch
// undoes of main that opts.Keep keeps (ADR-0155). With opts.Batch, the
// queue's landing, a blocker for each reason the merge queue does not let the
// story land from that batch now, and the story's rebased head in the batch
// as the verified commit (ADR-0176). Under pull-request
// acceptance, a story's acceptance, whoever accepts, also has a blocker for
// each forge condition of ADR-0169 that fails, read from the state flai serve
// keeps of its pull request.
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
		// a path the branch undoes of main's change waits for the operator
		// to restore it or keep the removal (ADR-0155)
		if undone, err := storygit.RecordedUndone(r, repo, it.ID); err != nil {
			res.Blockers = append(res.Blockers, fmt.Sprintf("cannot read which paths %s undoes of the main branch's changes: %v", res.Branch, err))
		} else {
			res.UndoesMain = UndonePaths(undone, opts.Keep)
			for _, u := range res.UndoesMain {
				if u.KeptReason == "" {
					res.Blockers = append(res.Blockers, UndoneBlocker(it.ID, u.Undone))
				}
			}
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
	// the forge's checks and reviews gate every acceptance, the
	// orchestrator's included (ADR-0169)
	if it.Type == workitem.Story && repo.Manifest.Acceptance.PullRequest() && !res.Resumed {
		forgeGate(r, repo, it.ID, res)
	}
	// the queue's landing takes the batch's verification of the story's
	// rebased head (ADR-0176); for the orchestrator, in place of its own
	var landing []string
	if opts.Batch != "" {
		res.Batch = opts.Batch
		landing = batchLanding(r, repo, it.ID, res)
		if !opts.Orchestrator {
			res.Blockers = append(res.Blockers, landing...)
		}
	}
	if opts.Orchestrator {
		res.Evidence = opts.Evidence
		res.OrchestratorBlockers, res.Verified = orchestratorBlockers(r, repo, it, opts, res.Protected, res.Verified, landing)
		for _, b := range res.OrchestratorBlockers {
			res.Blockers = append(res.Blockers, b.Message)
		}
	}
	threadsToResolve(repo, it, res)
	return res, nil
}

// batchLanding is every reason the merge queue does not let story id land
// now from res.Batch, as mergequeue.Landable gives them, a queue's state or
// main branch it cannot read among them, none when it may; when it may, it
// sets res.Verified to the story's rebased head in the batch, the commit the
// batch verified (ADR-0176).
func batchLanding(r execx.Runner, repo *workitem.Repo, id string, res *Acceptance) []string {
	head, reasons, err := mergequeue.Landable(r, repo, id, res.Batch)
	switch {
	case err != nil:
		return []string{fmt.Sprintf("cannot check that %s may land from merge queue batch %s: %v", id, res.Batch, err)}
	case len(reasons) > 0:
		return reasons
	case head.Rebased == "":
		return []string{fmt.Sprintf("merge queue batch %s kept no rebased head for %s, so the commit it verified is not known; the batch must be built again from the main branch", res.Batch, id)}
	}
	res.Verified = head.Rebased
	return nil
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

// forgeGate sets on res the forge blockers of story id under pull-request
// acceptance (ADR-0169), each also in Blockers, the URL of its pull request
// when flai serve kept one, and the merge through the remote that takes the
// local merge's place (ADR-0170). The kept state is judged against the head
// of res.Branch, the story's branch when it exists; a state or a head that
// cannot be read is itself a blocker, so the gate never opens unread. A pull
// request flai serve read as merged on the forge already passes: its checks
// and reviews were met where it merged, and it is accepted with no second
// merge.
func forgeGate(r execx.Runner, repo *workitem.Repo, id string, res *Acceptance) {
	st, kept, err := forge.ReadStatus(repo.CacheDir(), id)
	res.Remote = &RemoteMerge{Method: forge.MergeMethod(repo.Manifest.Acceptance.MergeMethodOrDefault())}
	if err == nil && kept {
		res.Remote.PullRequest, res.Remote.Head = st.PullRequest, st.HeadSHA
		res.Remote.AlreadyMerged = st.State == forge.PullRequestMerged
	}
	switch {
	case err != nil:
		res.ForgeBlockers = []Blocker{{Code: BlockPullRequestMissing, Message: fmt.Sprintf("cannot read the pull request flai serve kept for %s: %v (ADR-0169)", id, err)}}
	case kept && st.State == forge.PullRequestMerged:
		res.ForgeBlockers = nil
	case kept && res.Branch != "":
		out, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", res.Branch+"^{commit}")
		if err != nil {
			res.ForgeBlockers = []Blocker{{Code: BlockForgeStale, Message: fmt.Sprintf("cannot read the head of %s to check that flai serve read %s's pull request #%d there: %v; check that git can read the branch (ADR-0169)", res.Branch, id, st.PullRequest, err), URL: st.URL}}
		} else {
			res.ForgeBlockers = ForgeBlockers(id, st, kept, strings.TrimSpace(out))
		}
	default:
		res.ForgeBlockers = ForgeBlockers(id, st, kept, "")
	}
	if kept {
		res.PullRequest = st.URL
	}
	for _, b := range res.ForgeBlockers {
		res.Blockers = append(res.Blockers, b.Message)
	}
}

// ForgeBlockers are what in st, story id's pull request as flai serve kept
// it, stops its acceptance under pull-request acceptance (ADR-0169): kept
// says whether flai serve kept a status, and head is the story branch's head
// commit, in full, or "" when it is not known.
func ForgeBlockers(id string, st forge.Status, kept bool, head string) []Blocker {
	branch := storygit.Branch(id)
	if !kept {
		return []Blocker{{Code: BlockPullRequestMissing, Message: fmt.Sprintf("flai serve has kept no pull request for %s: under pull-request acceptance it pushes %s and opens the pull request while the story is in review, behind the forge host action, and reads its checks and reviews; accept once it has (ADR-0166, ADR-0169)", id, branch)}}
	}
	if head != "" && st.HeadSHA != head {
		return []Blocker{{Code: BlockForgeStale, Message: fmt.Sprintf("flai serve last read %s's pull request #%d at %s (at %s), but %s's head is %s: flai serve has not yet read the head, so its checks and reviews are not known; wait for its next look (ADR-0169)", id, st.PullRequest, shortSHA(st.HeadSHA), st.ReadAt, branch, shortSHA(head)), URL: st.URL}}
	}
	var out []Blocker
	at := shortSHA(st.HeadSHA)
	for _, c := range st.Checks {
		switch c.State {
		case forge.CheckFailed:
			found := ""
			if c.Summary != "" {
				found = ": " + c.Summary
			}
			out = append(out, Blocker{Code: BlockCheckFailing, Name: c.Name, URL: c.URL, Message: fmt.Sprintf("%s's pull request #%d has the required check %s failing at %s%s; fix it on %s, and flai serve reads it again once the branch is pushed (ADR-0169)", id, st.PullRequest, c.Name, at, found, branch)})
		case forge.CheckPending:
			out = append(out, Blocker{Code: BlockCheckPending, Name: c.Name, URL: c.URL, Message: fmt.Sprintf("%s's pull request #%d has the required check %s still pending at %s; accept once it passes (ADR-0169)", id, st.PullRequest, c.Name, at)})
		}
	}
	approvals := 0
	for _, d := range st.Decisions() {
		switch d.State {
		case forge.ReviewApproved:
			approvals++
		case forge.ReviewChangesRequested:
			out = append(out, Blocker{Code: BlockChangesRequested, Name: d.Author, URL: d.URL, Message: fmt.Sprintf("%s requested changes on %s's pull request #%d, and that review still stands; address it on the branch and have %s approve or dismiss it (ADR-0169)", d.Author, id, st.PullRequest, d.Author)})
		}
	}
	if approvals < st.RequiredApprovals {
		reviews := "approving reviews"
		if approvals == 1 {
			reviews = "approving review"
		}
		out = append(out, Blocker{Code: BlockApprovalsMissing, URL: st.URL, Message: fmt.Sprintf("%s's pull request #%d has %d %s of the %d its base branch requires; accept once enough reviewers approve (ADR-0169)", id, st.PullRequest, approvals, reviews, st.RequiredApprovals)})
	}
	return out
}

// shortSHA is a commit's name shortened to seven characters, as the forge
// shows a pull request's commits; short keeps twelve for a verified commit.
func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
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

// UndonePaths are the paths of undone, each with its reason when keep names
// it.
func UndonePaths(undone []storygit.Undone, keep []KeptRemoval) []UndonePath {
	var out []UndonePath
	for _, u := range undone {
		p := UndonePath{Undone: u}
		for _, k := range keep {
			if k.Path == u.Path {
				p.KeptReason = k.Reason
			}
		}
		out = append(out, p)
	}
	return out
}

// UndoneBlocker is why story id cannot be accepted while its branch undoes
// the main branch's change u, naming the ways out (ADR-0155).
func UndoneBlocker(id string, u storygit.Undone) string {
	lines := fmt.Sprintf("%d lines", len(u.Lines))
	if len(u.Lines) == 1 {
		lines = "1 line"
	}
	return fmt.Sprintf("%s undoes the main branch's change to %s, removing %s main added since the branch left it; restore main's version with flai stream restore %s --path %s, keep the removal with flai accept %s --keep-removal %s --reason \"<why>\", or send %s back to in-progress (ADR-0155)", storygit.Branch(id), u.Path, lines, id, u.Path, id, u.Path, id)
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
