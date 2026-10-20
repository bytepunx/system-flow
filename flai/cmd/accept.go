package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// acceptOptions are shared by `flai accept` and `flai move <story> done`:
// there is one way for a story to become done, and it is acceptance (S-0046).
type acceptOptions struct {
	by       string
	trailers []string
	dryRun   bool
	// verified and evidence are the orchestrator's (ADR-0093): the commit its
	// verifier passed, and the file holding its evidence, "-" for stdin.
	verified string
	evidence string
	stdin    func() io.Reader // the command's standard input, for --evidence -
	// batch is the merge queue batch whose verification the acceptance
	// takes, as flai serve lands the story from it (ADR-0176).
	batch string
	// keepRemovals are the paths whose removal of the main branch's lines the
	// operator keeps, and reasons why: one per path in order, or one for all
	// (ADR-0155).
	keepRemovals []string
	reasons      []string
	// sharedReason is the command's own --reason, for flai move, which
	// defines --reason before the acceptance's flags: one reason for all.
	sharedReason func() []string
}

func newAcceptCmd(a *app) *cobra.Command {
	var o acceptOptions
	c := &cobra.Command{
		Use:   "accept <id>",
		Short: "Accept an item: merge its branch, move to done, archive, commit",
		Long: `The operator's acceptance step as one command, per
design/conventions/work-management.md and git.md:

  0. rebase the story branch and fast-forward it into the main branch
  1. move the item to done (its rules apply: children closed, criteria checked);
     a story's epic follows it, to done when it was the epic's last open story
  2. flai archive for the item and its children and narrative, and for an
     epic that followed its story to done, the epic and its cancelled stories;
     every thread still open or answered on what is archived is resolved,
     as "<id> was accepted", by whoever accepts (I-0073), and every
     conversation still open of a story archived is closed, as "<id> was
     accepted", or "<id> was archived" for an epic's cancelled story (ADR-0120)
  3. git commit the work item, the archive, the threads, and the
     conversations; while another git process holds the index lock, git add
     and git commit are run again, for about nine seconds in all, and the
     lock is never removed
  4. tell every story in progress or in review whose touches cover a path
     the merge changed which paths those are, for its agent's MCP inbox

Under acceptance.mode: pull-request (ADR-0170), step 0 merges no branch on
this machine. Once the pull request's checks and reviews pass (ADR-0169), it
refuses, naming them, while the main branch has commits not on origin's main
that change paths outside the wip folder; merges the story's pull request on
the forge by acceptance.merge_method, at the head the checks passed at, and
deletes the story's branch on origin; fetches main and brings the main branch
up to origin's, rebasing its unpublished wip commits onto it in a temporary
worktree and moving it with git reset --keep, which keeps uncommitted edits;
and removes the story's worktree and branch. The merge is made here when the
forge's token can be read here, while the forge host action is on;
otherwise flai serve on the host makes it, and flai accept waits for its
answer (the dashboard's acceptance never hands it over). A pull request
merged on the forge already is accepted with no second merge, and a merge the
forge refuses stops the acceptance, with its reason, before anything changes
here. Steps 1 to 4 are as under local acceptance.

Acceptance computes no release, creates no tag, and pushes nothing (S-0087):
that is a deliberate step of its own, run when the operator chooses to
publish what has accumulated on main, not tied to any one item. See flai
release --pending.

flai move <story> done from review runs exactly this. An item that is already
done but was never archived (an older flai, a hand edit) is completed from
step 0 without a second transition. An item that is done and archived while
its archived file is not committed is one whose commit failed: the error
kept git's output and named this command. It is completed by resolving the
threads and closing the conversations left open on what it archived, as
step 2 does, and from step 3, committing what the main checkout holds under
the usual subject, and step 4 tells the open stories the paths the story's
commits changed. A done, archived item whose file is committed is refused as
already done.
The orchestrator completes neither: that is the operator's. --dry-run
changes nothing.

A story whose branch changes a path Claude Code protects (a .claude folder,
.mcp.json, and the others of ADR-0106) is accepted by its operator only: the
story's owner or the project's owner, or anyone but the orchestrator when
neither is named. Anyone else is refused, before anything is merged, with the
files named; --dry-run lists them.

Markdown flai accept commits passes the project's lint (S-0367): the fix
command of each test tier that has one runs over what the main checkout has
not committed before anything changes and again before step 3's commit, and
step 0's rebase runs it over what the story branch changes. When a tier's
check still fails on those files, the acceptance is refused with each file,
line, and error, before anything changes or before the merge; correct them
and accept again. --dry-run runs neither. With --json that refusal exits 4
and prints {"refused": {"id", "reason", "findings"}}, each finding with its
tier, checkout (main, or story with the story's ID), path, line, column,
level, rule, alias, and message, as markdownlint-cli2 gives them (S-0378).

--by orchestrator is the orchestrator's acceptance (ADR-0093), refused unless
orchestration.permissions.accept_reviews is on; under FLAI_ROLE=orchestrate no
other acceptance is allowed. It is refused, before anything is merged, unless
--verified names the story branch's head, every acceptance criterion is
ticked, every file the branch changes is under the story's touches and on no
path Claude Code protects, no thread on the story or its tasks is open, and
--evidence, a Verdict: line and one item "- <n>: <files>" per criterion,
names a changed file for each of them.
The evidence is written, with the commit, under ### Accepted by the
orchestrator in the story's Notes. With --dry-run the evidence is optional.

A story whose rebased branch undoes the main branch's change to a path, as a
revert replayed over main's identical change does (ADR-0155, I-0124), is
refused before anything is merged while the path is not kept: each path flai
stream sync recorded and its branch still undoes, and each path the
acceptance's own rebase finds. --dry-run lists each under undoes_main, with
main's lines the branch removes, and each not kept as a blocker. Restore
main's version with flai stream restore <id> --path <path>, send the story
back to in-progress, or, as the operator, keep the removal with
--keep-removal <path> (repeatable) and --reason "<why>": one --reason per
--keep-removal in order, or one for them all. A --keep-removal naming a path
the record does not hold, one without a reason, a reason on more than one
line, and any with --by orchestrator are refused before anything changes:
an undone path waits for the operator. Each kept path and its reason are
written, with the commit, under ### Removals kept in the story's Notes. With
--json the refusal exits 4 and prints {"refused": {"id", "reason",
"undoes_main"}}, each path with its path, lines, and old_base.

--batch <id> is the merge queue's landing (ADR-0176): flai serve runs it for
each story of a batch that passed, in queue order, as whoever enqueued the
story. It takes the batch's verification in place of the story head's own:
the verified commit is the story's rebased head in the batch. It is refused,
before anything is merged, unless the queue's batch is the one named, is not
dropped, and lets the story land, the story is the next of it to land, the
story's branch is still at the head the batch rebased a copy of, and the main
branch is where the batch's landing left it. With --by orchestrator it needs
no --verified, which it is refused beside, and no --evidence; the batch and
its verified commit are written under ### Accepted by the orchestrator, with
the evidence when given. The story is recorded as landed in the queue's
state, and the batch kept, where an acceptance without --batch takes the
story out of the queue and drops the batch that holds it. Under
acceptance.mode: pull-request, which has no merge queue, it is refused.`,
		Example: `  flai accept S-031 --by alex
  flai accept S-031 --by alex --dry-run
  flai accept S-031 --by orchestrator --verified 4f1c2a9 --dry-run
  flai accept S-031 --by orchestrator --verified 4f1c2a9 --evidence evidence.md
  flai accept S-031 --by alex --keep-removal docs/guide.md --reason "the story rewrites the section"
  flai accept S-031 --by orchestrator --batch 20261019-041200`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			res, err := a.acceptItem(repo, it, o)
			if err != nil {
				return a.acceptRefusal(err)
			}
			return a.printAccept(res)
		},
	}
	addAcceptFlags(c, &o)
	c.Flags().BoolVar(&o.dryRun, "dry-run", false, "show what would block acceptance and change nothing")
	c.Flags().StringVar(&o.batch, "batch", "", "the merge queue batch whose verification the acceptance takes, as flai serve lands it (ADR-0176)")
	return c
}

// addAcceptFlags registers the flags acceptance shares with move.
func addAcceptFlags(c *cobra.Command, o *acceptOptions) {
	if c.Flags().Lookup("by") == nil {
		c.Flags().StringVar(&o.by, "by", "", "who accepted (default: config author)")
	}
	c.Flags().StringArrayVar(&o.trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().StringVar(&o.verified, "verified", "", "with --by orchestrator: the commit its verifier passed, the story branch's head")
	c.Flags().StringVar(&o.evidence, "evidence", "", "with --by orchestrator: file holding its evidence, a Verdict: line and - <n>: <files> per criterion (- reads standard input)")
	c.Flags().StringArrayVar(&o.keepRemovals, "keep-removal", nil, "as the operator: accept the branch's removal of the main branch's lines at this recorded path, with --reason (repeatable)")
	if f := c.Flags().Lookup("reason"); f == nil {
		c.Flags().StringArrayVar(&o.reasons, "reason", nil, "with --keep-removal: why the removal is kept, one per --keep-removal in order or one for all (repeatable)")
	} else {
		o.sharedReason = func() []string {
			if f.Changed {
				return []string{f.Value.String()}
			}
			return nil
		}
	}
	o.stdin = c.InOrStdin
}

// accepted is what flai accept did, or under DryRun would do: the preview's
// acceptance and the conversations it closes (ADR-0120).
type accepted struct {
	*preview.Acceptance
	// ClosedMessages are the conversations still open of the stories the
	// acceptance archives, which it closes, or under DryRun would.
	ClosedMessages []string `json:"closed_messages,omitempty"`
	// closing names each of ClosedMessages with the story it was open with,
	// for the text output.
	closing []string
}

// acceptItem runs the acceptance flow for a story or an epic.
func (a *app) acceptItem(repo *workitem.Repo, it *workitem.Item, o acceptOptions) (*accepted, error) {
	opts, err := a.acceptance(repo, it, o)
	if err != nil {
		return nil, err
	}
	res, err := preview.AcceptWith(a.runner, repo, it, orDefault(o.by, a.author()), a.now(), a.logger(), opts)
	if err != nil {
		return nil, err
	}
	if opts.Orchestrator && res.Resumed {
		left := "is done but was never archived"
		if res.CommitOnly {
			left = "is done and archived but its acceptance was never committed"
		}
		return nil, fmt.Errorf("rule: %s %s; the orchestrator accepts only a story in review, so the operator completes this acceptance with flai accept %s", it.ID, left, it.ID)
	}
	acc := &accepted{Acceptance: res}
	messagesToClose(repo, it, acc)
	res.DryRun = o.dryRun
	if !o.dryRun {
		res.By = orDefault(o.by, a.author())
	}
	// the merge queue's landing leaves them as they are (S-0402)
	if dirty := res.Uncommitted; len(dirty) > 0 && !a.yes && !o.dryRun && o.batch == "" {
		return nil, fmt.Errorf("working tree has uncommitted changes outside wip (%s); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them", strings.Join(dirty, ", "))
	}
	if err := keptHeld(it.ID, res.UndoesMain, opts.Keep); err != nil {
		return nil, err
	}
	// under pull-request acceptance the merge is made through the remote,
	// here or by flai serve, as the token can be read (ADR-0170)
	if res.Remote != nil {
		res.Remote.ThroughServe = !a.remoteRoute(repo).direct
	}
	if len(res.Blockers) > 0 && !o.dryRun {
		if left := notKept(res.UndoesMain); len(left) > 0 {
			return nil, undoneRefused(it.ID, res.Blockers, left)
		}
		return nil, fmt.Errorf("%s cannot be accepted yet: %s", it.ID, strings.Join(res.Blockers, "; "))
	}
	if o.dryRun {
		return acc, nil
	}
	acc.ClosedMessages, acc.closing = nil, nil
	useGit, hasBranch := a.inGitWorkTree(repo.MainRoot), res.Branch != ""
	// what the acceptance commits from the main checkout passes the lint
	// its fix leaves, or nothing changes (S-0367, TH-0413)
	if useGit {
		if err := a.refuseUnfixedUncommitted(context.Background(), it.ID, repo.Root, leftUncommitted(res)); err != nil {
			return nil, err
		}
	}
	if res.CommitOnly {
		return a.finishCommit(repo, it, acc, o)
	}

	// 0. bring the story branch into the main branch (ADR-0019): under
	// pull-request acceptance by merging its pull request on the forge and
	// bringing main up to origin's (ADR-0170); a resumed acceptance's merge
	// is behind it
	var changed []string // what the merge brought, for the stories still open
	switch {
	case it.Type != workitem.Story:
	case repo.Manifest.Acceptance.PullRequest():
		if res.Remote != nil && useGit {
			if changed, err = a.acceptThroughRemote(context.Background(), repo, it, res); err != nil {
				return nil, err
			}
		}
	default:
		before := ""
		if hasBranch {
			before = a.headOf(repo.MainRoot)
		}
		merged, err := a.mergeStoryBranch(repo, it.ID, opts.Keep)
		if err != nil {
			a.logAcceptConflict(repo, it.ID, err)
			return nil, err
		}
		res.Merged = merged
		if merged {
			a.acceptStep(it, "merged", res.Branch+" rebased and fast-forwarded into the main branch")
			if before != "" {
				if changed, err = a.changedSince(repo.MainRoot, before); err != nil {
					a.logger().Warn("overlap notices not sent", "component", "accept", "item", it.ID, "err", err)
				}
			}
		}
	}
	// 0a. a story accepted by hand leaves the merge queue (ADR-0173); one
	// the queue lands is recorded landed, its batch kept (ADR-0176)
	if it.Type == workitem.Story {
		a.leaveQueue(repo, it.ID, o.batch, res.By)
	}

	// 1. done (epics: walk through review if needed)
	items, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, err
	}
	from := it.Status
	if !res.Resumed {
		for _, st := range preview.StepsToDone(it.Status) {
			if _, err := repo.Move(it, st, workitem.MoveOptions{By: orDefault(o.by, a.author()), Now: a.now(), Items: items, Board: board}); err != nil {
				return nil, err
			}
		}
		// the orchestrator's evidence goes into the acceptance commit with
		// the story (ADR-0093)
		if opts.Orchestrator {
			it.Body = withNotesSection(it.Body, orchestratorNotes(res, a.now()))
		}
		// so does each removal of main's lines the operator keeps (ADR-0155)
		if len(opts.Keep) > 0 {
			it.Body = withNotesSection(it.Body, keptNotes(opts.Keep, res.By, a.now()))
		}
		if err := repo.Save(it); err != nil {
			return nil, err
		}
		if _, err := repo.RollUp(it); err != nil {
			return nil, err
		}
		if it.Type == workitem.Story {
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return nil, err
			}
		}
	}
	res.Status = it.Status
	a.acceptStep(it, "done", it.ID+" moved to done")
	// 1a. the story's epic follows it, to done with its last open story
	// (S-0200); the story's roll-up has already summed the epic's usage
	if it.Type == workitem.Story && !res.Resumed {
		if res.Epic, err = a.followAccepted(repo, it, from, orDefault(o.by, a.author())); err != nil {
			return nil, err
		}
	}
	// 2. archive
	items, _ = repo.List(false)
	ap, err := repo.PlanArchive(items, preview.ArchivedWith(items, it.ID, res.Epic))
	if err != nil {
		return nil, err
	}
	if err := repo.Archive(ap); err != nil {
		return nil, err
	}
	res.Archived = len(ap.Items)
	a.acceptStep(it, "archived", fmt.Sprintf("%d item(s) and the narrative archived", res.Archived))
	archived := make([]string, len(ap.Items))
	for i, x := range ap.Items {
		archived[i] = x.ID
	}
	if err := a.resolveAccepted(repo, it, res, archived, o.by); err != nil {
		return nil, err
	}
	if err := a.closeAccepted(repo, it, acc, archived, o.by); err != nil {
		return nil, err
	}
	if err := a.refreshIndex(repo); err != nil {
		return nil, err
	}
	if !useGit {
		return acc, nil
	}
	// 3. commit: the item, the archive, and nothing else. No release is
	// computed, no tag created, no push made (S-0087) — that is a publish, a
	// deliberate step of its own over everything accumulated on main, not
	// tied to any one item's acceptance. See flai release --pending.
	if err := a.commitAcceptance(repo, it, res, o.trailers); err != nil {
		return nil, err
	}
	// 4. tell the stories still open what changed under their claim (S-0132)
	a.tellAccepted(repo, it, res, orDefault(o.by, a.author()), changed)
	return acc, nil
}

// leaveQueue takes story, which by accepted, out of the merge queue in the
// main checkout's .flai-cache and records its outcome: by hand, with batch
// "", marking dropped the batch that holds it; as the landing of merge queue
// batch batch, recording it landed there and keeping the batch (ADR-0176).
// The story is merged by then, so a state that cannot be written is warned
// of, and the acceptance goes on.
func (a *app) leaveQueue(repo *workitem.Repo, story, batch, by string) {
	var err error
	if batch != "" {
		_, err = mergequeue.Landed(repo.CacheDir(), story, batch, by, a.now())
	} else {
		_, err = mergequeue.AcceptedByHand(repo.CacheDir(), story, by, a.now())
	}
	if err != nil {
		a.logger().Warn("accepted story not taken out of the merge queue", "component", "accept", "item", story, "batch", batch, "err", err)
	}
}

// logAcceptConflict writes a line of the conflict log (S-0337) when err is
// the acceptance's own rebase of story's branch stopped on conflicts, with
// the files flai resolves itself left out of its paths; a rebase found in
// progress, or any other refusal, is not one.
func (a *app) logAcceptConflict(repo *workitem.Repo, story string, err error) {
	var stop *syncStopped
	if !errors.As(err, &stop) || !stop.RebaseInProgress || stop.Found {
		return
	}
	c := storygit.Conflict{At: a.now(), Story: story, Where: storygit.WhereAccept, Paths: stop.Conflicts}
	storygit.RecordConflict(repo.CacheDir(), c, issues.SyncFiles(repo, story, a.runner, a.now), a.logger())
}

// finishCommit completes the acceptance of it, done and archived, whose
// commit failed (I-0100): steps 3 and 4 only, with no merge, no transition,
// and no archive, since those happened.
func (a *app) finishCommit(repo *workitem.Repo, it *workitem.Item, acc *accepted, o acceptOptions) (*accepted, error) {
	res := acc.Acceptance
	// the threads and conversations the run that archived left open: an
	// older flai's, or one whose resolving or closing failed
	all, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	archived := preview.ThreadItems(all, it, res.Epic)
	if err := a.resolveAccepted(repo, it, res, archived, o.by); err != nil {
		return nil, err
	}
	if err := a.closeAccepted(repo, it, acc, archived, o.by); err != nil {
		return nil, err
	}
	if err := a.commitAcceptance(repo, it, res, o.trailers); err != nil {
		return nil, err
	}
	// The merge is behind this run, so what it changed cannot be read from
	// the head before it as step 0 reads it. The paths are taken from the
	// story's commits instead, outside the wip folder: what its branch
	// brought into the main branch.
	var changed []string
	if it.Type == workitem.Story {
		committed, err := storygit.Committed(a.runner, repo)
		if err != nil {
			a.logger().Warn("overlap notices not sent", "component", "accept", "item", it.ID, "err", err)
		}
		changed = committed[it.ID]
	}
	a.tellAccepted(repo, it, res, orDefault(o.by, a.author()), changed)
	return acc, nil
}

// resolveAccepted is step 2a: it resolves the threads still open or answered
// on the items the acceptance of it archived, as by, with the acceptance as
// the reason, so that none is left open on an archived item (I-0073). The
// thread files go into the acceptance commit.
func (a *app) resolveAccepted(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, archived []string, by string) error {
	resolved, err := threads.ResolveOnItems(repo, archived, orDefault(by, a.author()), it.ID+" was accepted", a.now())
	if err != nil {
		return fmt.Errorf("resolve the threads on what the acceptance of %s archived: %w\n%s, but not committed; once the thread above is fixed, run flai accept %s to resolve the threads left and commit", it.ID, err, acceptedState(it.ID, res), it.ID)
	}
	res.ResolvedThreads = resolved
	if len(resolved) > 0 {
		a.acceptStep(it, "threads", "resolved "+strings.Join(resolved, ", ")+", open or answered on what it archived")
	}
	return nil
}

// closeAccepted is step 2b: it closes the conversations still open of the
// stories the acceptance of it archived (ADR-0120), as by: the accepted
// story's as "<id> was accepted", and those of the cancelled stories of an
// epic archived with it as "<id> was archived", which is what happened to
// them. The files go into the acceptance commit.
func (a *app) closeAccepted(repo *workitem.Repo, it *workitem.Item, acc *accepted, archived []string, by string) error {
	by = orDefault(by, a.author())
	open, err := openConversations(repo, archived)
	var closed, rest []string
	if err == nil {
		closed, err = messages.CloseOn(repo, []string{it.ID}, by, "accepted", a.now())
	}
	if err == nil {
		rest, err = messages.CloseOn(repo, archived, by, "archived", a.now())
	}
	if err != nil {
		return fmt.Errorf("close the conversations of what the acceptance of %s archived: %w\n%s, but not committed; once the conversation above is fixed, run flai accept %s to close the conversations left and commit", it.ID, err, acceptedState(it.ID, acc.Acceptance), it.ID)
	}
	closed = append(closed, rest...)
	sort.Strings(closed)
	acc.ClosedMessages, acc.closing = closed, closingLines(among(open, closed), archived)
	if len(closed) > 0 {
		a.acceptStep(it, "messages", "closed "+strings.Join(closed, ", ")+", open on what it archived")
	}
	return nil
}

// messagesToClose sets on acc the conversations still open of the stories
// accepting it archives, which the acceptance closes (ADR-0120). One that
// cannot be read is a blocker, as a thread is: the acceptance would stop on
// it after archiving.
func messagesToClose(repo *workitem.Repo, it *workitem.Item, acc *accepted) {
	all, err := repo.List(true)
	if err != nil {
		acc.Blockers = append(acc.Blockers, fmt.Sprintf("cannot list the items to find the conversations accepting %s closes: %v", it.ID, err))
		return
	}
	ids := preview.ThreadItems(all, it, acc.Epic)
	open, err := openConversations(repo, ids)
	if err != nil {
		acc.Blockers = append(acc.Blockers, fmt.Sprintf("cannot read the conversations of what accepting %s archives, which it closes: %v; fix the conversation file it names", it.ID, err))
		return
	}
	for _, c := range open {
		acc.ClosedMessages = append(acc.ClosedMessages, c.ID)
	}
	acc.closing = closingLines(open, ids)
}

// openConversations are the conversations still stored open that name one of
// ids, in ID order: those messages.CloseOn closes on the same IDs.
func openConversations(repo *workitem.Repo, ids []string) ([]*messages.Conversation, error) {
	leaving := canonicalSet(ids)
	all, err := messages.List(repo)
	if err != nil {
		return nil, err
	}
	var out []*messages.Conversation
	for _, c := range all {
		if c.Status != messages.StatusClosed && (leaving[workitem.CanonicalID(c.From)] || leaving[workitem.CanonicalID(c.To)]) {
			out = append(out, c)
		}
	}
	return out, nil
}

// closingLines names each conversation of convs as the text output does: its
// ID and the story of its two that is not among ids, the one it was open
// with, or both stories when both are.
func closingLines(convs []*messages.Conversation, ids []string) []string {
	leaving := canonicalSet(ids)
	var lines []string
	for _, c := range convs {
		from, to := leaving[workitem.CanonicalID(c.From)], leaving[workitem.CanonicalID(c.To)]
		switch {
		case from && to:
			lines = append(lines, c.ID+", open between "+c.From+" and "+c.To)
		case to:
			lines = append(lines, c.ID+", open with "+c.From)
		default:
			lines = append(lines, c.ID+", open with "+c.To)
		}
	}
	return lines
}

// among are the conversations of convs whose IDs ids lists.
func among(convs []*messages.Conversation, ids []string) []*messages.Conversation {
	var out []*messages.Conversation
	for _, c := range convs {
		if slices.Contains(ids, c.ID) {
			out = append(out, c)
		}
	}
	return out
}

// canonicalSet is ids in canonical form, as a set.
func canonicalSet(ids []string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[workitem.CanonicalID(id)] = true
	}
	return set
}

// commitAcceptance is step 3: it runs the lint fix over what the main
// checkout has not committed (S-0367) and commits everything there under the
// acceptance's subject, but what the merge queue's landing leaves
// uncommitted (S-0402). git add and git commit are run again while another
// git process holds the index lock (I-0100). A commit that fails even so says
// what the acceptance has done and how to finish it.
func (a *app) commitAcceptance(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, trailers []string) error {
	left := leftUncommitted(res)
	a.fixUncommitted(context.Background(), it.ID, repo.Root, left)
	if _, err := storygit.RunPastIndexLock(a.runner, repo.Root, addArgs(left)...); err != nil {
		return notCommitted(it.ID, res, "not yet staged", err)
	}
	msg := acceptSubject(it.ID, res.Epic)
	for _, t := range trailers {
		msg += "\n\n" + t
	}
	if _, err := storygit.RunPastIndexLock(a.runner, repo.Root, "commit", "-q", "-m", msg); err != nil {
		return notCommitted(it.ID, res, "staged", err)
	}
	a.acceptStep(it, "committed", firstLine(msg))
	return nil
}

// acceptSubject is the subject of the acceptance commit of id, naming the
// epic that followed it to done.
func acceptSubject(id string, epic *workitem.Followed) string {
	subject := fmt.Sprintf("chore: [%s] accept and archive", id)
	if epic != nil && epic.To == workitem.Done {
		subject += ", with " + epic.ID
	}
	return subject
}

// notCommitted is the error of an acceptance of id whose commit failed with
// err, git's whole output in it: what is done, where its changes are, and
// that flai accept finishes it.
func notCommitted(id string, res *preview.Acceptance, changes string, err error) error {
	return fmt.Errorf("commit the acceptance of %s: %w\n%s, and its changes are %s in the main checkout but not committed; once git's error above is resolved, run flai accept %s to commit them", id, err, acceptedState(id, res), changes, id)
}

// acceptedState says what the acceptance of id has done: done and archived,
// and its branch merged when it was.
func acceptedState(id string, res *preview.Acceptance) string {
	state := id + " is done and archived"
	if res.Merged {
		state = res.Branch + " is merged and " + state
	}
	return state
}

// tellAccepted is step 4: it tells the stories still open which of changed
// their claims cover (S-0132).
func (a *app) tellAccepted(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, by string, changed []string) {
	res.Overlaps = a.tellOverlaps(repo, it, by, changed)
	if len(res.Overlaps) > 0 {
		ids := make([]string, len(res.Overlaps))
		for i, n := range res.Overlaps {
			ids[i] = n.ID
		}
		a.acceptStep(it, "told", "told "+strings.Join(ids, ", ")+" which changed paths their claims cover")
	}
}

// acceptance is what o adds to the acceptance preview of it, after refusing
// what o may not ask: under FLAI_ROLE=orchestrate an acceptance as anyone but
// the orchestrator, --verified or --evidence for anyone else, the
// orchestrator's acceptance while orchestration.permissions.accept_reviews is
// off, and its real run without evidence or with evidence that cannot be read
// (ADR-0093); --batch under pull-request acceptance, which has no merge
// queue, or beside --verified, whose place it takes, while the orchestrator's
// landing from a batch needs no evidence (ADR-0176); and --keep-removal as
// keptRemovals refuses it (ADR-0155).
func (a *app) acceptance(repo *workitem.Repo, it *workitem.Item, o acceptOptions) (preview.AcceptOptions, error) {
	keep, err := keptRemovals(it.ID, o)
	if err != nil {
		return preview.AcceptOptions{}, err
	}
	orchestrator := workitem.IsOrchestrator(o.by)
	switch {
	case os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate && !orchestrator:
		given := "no --by"
		if o.by != "" {
			given = "--by " + o.by
		}
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator accepts a story only as itself, so give --by %s, not %s (ADR-0093)", workitem.ActivityOrchestrator, given)
	case o.batch != "" && repo.Manifest.Acceptance.PullRequest():
		return preview.AcceptOptions{}, fmt.Errorf("rule: --batch %s is the merge queue's landing, and %s; accept %s without --batch", o.batch, mergequeue.OffUnderPullRequest, it.ID)
	case o.batch != "" && o.verified != "":
		return preview.AcceptOptions{}, fmt.Errorf("rule: --batch %s takes the merge queue batch's verification of %s's rebased head in place of --verified's: give one or the other (ADR-0176)", o.batch, it.ID)
	case !orchestrator && (o.verified != "" || o.evidence != ""):
		return preview.AcceptOptions{}, fmt.Errorf("rule: --verified and --evidence are the orchestrator's acceptance: give them with --by %s, or leave them out", workitem.ActivityOrchestrator)
	case !orchestrator:
		return preview.AcceptOptions{Keep: keep, Batch: o.batch}, nil
	}
	if err := repo.OrchestratorPermits(manifest.PermitAcceptReviews, "accepts a story", it.ID); err != nil {
		return preview.AcceptOptions{}, fmt.Errorf("rule: %w", err)
	}
	opts := preview.AcceptOptions{Orchestrator: true, Verified: o.verified, Batch: o.batch}
	if o.evidence == "" {
		// the queue's eligibility checked what the evidence would vouch for
		// before the batch was built (ADR-0176)
		if o.dryRun || o.batch != "" {
			return opts, nil
		}
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator's acceptance of %s needs its evidence: --evidence <file> or -, a Verdict: line and one - <n>: <files> item per criterion (ADR-0093)", it.ID)
	}
	text, err := a.readEvidence(o)
	if err != nil {
		return preview.AcceptOptions{}, err
	}
	if opts.Evidence, err = preview.ParseEvidence(text); err != nil {
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator's evidence for %s cannot be read as ADR-0093 asks: %w", it.ID, err)
	}
	return opts, nil
}

// readEvidence is the text of the orchestrator's evidence: --evidence's file,
// or standard input for "-". A heading in it is refused, since the text goes
// under a heading of its own in the story's Notes.
func (a *app) readEvidence(o acceptOptions) (string, error) {
	var data []byte
	var err error
	if o.evidence == "-" {
		data, err = io.ReadAll(o.stdin())
	} else {
		path := o.evidence
		if !filepath.IsAbs(path) && a.cwd != "" {
			path = filepath.Join(a.cwd, path)
		}
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return "", fmt.Errorf("read the orchestrator's evidence from %s: %w; give a readable file, or - for standard input", o.evidence, err)
	}
	text := strings.TrimSpace(string(data))
	for i, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			return "", fmt.Errorf("rule: line %d of the orchestrator's evidence is a heading (%s): write a Verdict: line and - <n>: <files> items without headings, since flai puts the evidence under ### Accepted by the orchestrator in the story's Notes", i+1, line)
		}
	}
	return text, nil
}

// orchestratorNotes is what the orchestrator's acceptance writes under the
// story's Notes (ADR-0093): the commit its verifier passed, or, as the merge
// queue's landing, the story's rebased head in the batch that verified it and
// the batch (ADR-0176), when, and its evidence as given.
func orchestratorNotes(res *preview.Acceptance, now time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### Accepted by the orchestrator\n\n- Verified: %s\n", res.Verified)
	if res.Batch != "" {
		fmt.Fprintf(&b, "- Batch: %s\n", res.Batch)
	}
	fmt.Fprintf(&b, "- At: %s\n", now.UTC().Format(time.RFC3339))
	if res.Evidence != nil {
		fmt.Fprintf(&b, "\n%s\n", res.Evidence.Text)
	}
	return b.String()
}

// keptRemovals are o's --keep-removal paths, each with its --reason, after
// refusing what ADR-0155 does not allow: any keep by the orchestrator, since
// an undone path waits for the operator; a keep without a reason, an empty
// one, or one on more than one line; a number of reasons that is neither one
// nor one per path; a path named twice; and a reason with no keep, for flai
// accept, whose --reason says only why a removal is kept.
func keptRemovals(id string, o acceptOptions) ([]preview.KeptRemoval, error) {
	reasons := o.reasons
	if o.sharedReason != nil {
		reasons = o.sharedReason()
	}
	if len(o.keepRemovals) == 0 {
		if len(o.reasons) > 0 {
			return nil, errors.New("rule: --reason says why a --keep-removal is kept; give it with --keep-removal <path>, or leave it out")
		}
		return nil, nil
	}
	switch {
	case workitem.IsOrchestrator(o.by):
		return nil, fmt.Errorf("rule: the orchestrator keeps no removal of the main branch's lines: a path %s undoes of main waits for the operator, who restores it with flai stream restore %s --path <path> or keeps the removal with flai accept %s --keep-removal <path> --reason \"<why>\" (ADR-0155)", storyBranch(id), id, id)
	case len(reasons) == 0:
		return nil, fmt.Errorf("rule: --keep-removal %s needs a --reason saying why the removal of the main branch's lines is kept", strings.Join(o.keepRemovals, ", "))
	case len(reasons) != 1 && len(reasons) != len(o.keepRemovals):
		return nil, fmt.Errorf("rule: give one --reason for every --keep-removal, or one per --keep-removal in order: %s for %s", plural(len(reasons), "reason"), plural(len(o.keepRemovals), "path"))
	}
	var keep []preview.KeptRemoval
	for i, p := range o.keepRemovals {
		p = path.Clean(filepath.ToSlash(p))
		reason := strings.TrimSpace(reasons[min(i, len(reasons)-1)])
		switch {
		case reason == "":
			return nil, fmt.Errorf("rule: the --reason for --keep-removal %s is empty; say why the removal of the main branch's lines is kept", p)
		case strings.ContainsAny(reason, "\r\n"):
			return nil, fmt.Errorf("rule: the --reason for --keep-removal %s spans lines; give it on one line, since it is written as one item under ### Removals kept in the Notes of %s", p, id)
		case slices.ContainsFunc(keep, func(k preview.KeptRemoval) bool { return k.Path == p }):
			return nil, fmt.Errorf("rule: --keep-removal %s is given twice; name each path once", p)
		}
		keep = append(keep, preview.KeptRemoval{Path: p, Reason: reason})
	}
	return keep, nil
}

// keptHeld refuses a keep of a path that undone, the paths the story's record
// holds, does not name, naming those it does.
func keptHeld(id string, undone []preview.UndonePath, keep []preview.KeptRemoval) error {
	held := make([]string, len(undone))
	for i, u := range undone {
		held[i] = u.Path
	}
	for _, k := range keep {
		if slices.Contains(held, k.Path) {
			continue
		}
		holds := "holds no path"
		if len(held) > 0 {
			holds = "holds only " + strings.Join(held, ", ")
		}
		return fmt.Errorf("rule: --keep-removal %s names a path %s does not undo of the main branch's changes: its record %s; keep only a path flai accept %s --dry-run lists under undoes_main", k.Path, storyBranch(id), holds, id)
	}
	return nil
}

// notKept are the paths of undone the acceptance does not keep.
func notKept(undone []preview.UndonePath) []preview.UndonePath {
	var left []preview.UndonePath
	for _, u := range undone {
		if u.KeptReason == "" {
			left = append(left, u)
		}
	}
	return left
}

// undoneRefusal is the refusal of an acceptance whose story branch undoes
// the main branch's change to a path the operator does not keep (ADR-0155):
// its text, and the paths as data for flai accept --json, which prints it as
// {"refused": {...}} and exits 4.
type undoneRefusal struct {
	ID         string               `json:"id"`
	Reason     string               `json:"reason"`
	UndoesMain []preview.UndonePath `json:"undoes_main"`
}

// Error is the refusal's text.
func (e *undoneRefusal) Error() string { return e.Reason }

// undoneRefused is the refusal of id's acceptance for left, the paths not
// kept, giving blockers as why, or each path's blocker when blockers is nil.
func undoneRefused(id string, blockers []string, left []preview.UndonePath) *undoneRefusal {
	if blockers == nil {
		for _, u := range left {
			blockers = append(blockers, preview.UndoneBlocker(id, u.Undone))
		}
	}
	return &undoneRefusal{ID: id, Reason: fmt.Sprintf("%s cannot be accepted yet: %s", id, strings.Join(blockers, "; ")), UndoesMain: left}
}

// acceptRefusal is err as acceptRefused gives it, and under --json, when err
// is a refusal for paths the branch undoes of main, that refusal printed as
// {"refused": {...}} with exit 4, so the dashboard shows each path (ADR-0155).
func (a *app) acceptRefusal(err error) error {
	var r *undoneRefusal
	if !a.jsonOut || !errors.As(err, &r) {
		return a.acceptRefused(err)
	}
	_ = a.printJSON(map[string]any{"refused": r})
	return &exitError{code: exitDocRefused, msg: err.Error()}
}

// keptNotes is what an acceptance that keeps removals of the main branch's
// lines writes under the story's Notes (ADR-0155): each path, who kept it,
// when, and why.
func keptNotes(keep []preview.KeptRemoval, by string, now time.Time) string {
	var b strings.Builder
	b.WriteString("### Removals kept\n\n")
	for _, k := range keep {
		fmt.Fprintf(&b, "- `%s`, kept by %s at %s: %s\n", k.Path, by, now.UTC().Format(time.RFC3339), k.Reason)
	}
	return b.String()
}

// withNotesSection is body with section appended to the end of its "## Notes"
// section, which is added at the end of body when it has none.
func withNotesSection(body, section string) string {
	const notes = "## Notes"
	start := -1
	if strings.HasPrefix(body, notes+"\n") {
		start = 0
	} else if i := strings.Index(body, "\n"+notes+"\n"); i >= 0 {
		start = i + 1
	}
	if start < 0 {
		return strings.TrimRight(body, "\n") + "\n\n" + notes + "\n\n" + section
	}
	end := len(body)
	if i := strings.Index(body[start+len(notes):], "\n## "); i >= 0 {
		end = start + len(notes) + i + 1
	}
	head, tail := strings.TrimRight(body[:end], "\n"), body[end:]
	if tail != "" {
		return head + "\n\n" + section + "\n" + tail
	}
	return head + "\n\n" + section
}

// followAccepted moves and saves the epic of story, accepted from from, the
// way its acceptance takes it, and says where it went; nil when it stays.
func (a *app) followAccepted(repo *workitem.Repo, story *workitem.Item, from, by string) (*workitem.Followed, error) {
	all, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	epic, followed, err := repo.Follow(all, story, from, by, a.now(), true)
	if err != nil || followed == nil {
		return nil, err
	}
	if err := repo.Save(epic); err != nil {
		return nil, err
	}
	a.acceptStep(story, "epic", fmt.Sprintf("%s followed %s to %s", epic.ID, story.ID, epic.Status))
	return followed, nil
}

// acceptStep logs one completed step of an acceptance as an info event with
// stable fields, so a client such as the dashboard can show progress while
// the command runs (S-0041). The message is fixed; what varies is in fields.
func (a *app) acceptStep(it *workitem.Item, step, detail string) {
	a.logger().Info("acceptance step", "component", "accept", "item", it.ID, "step", step, "detail", detail)
}

func (a *app) printAccept(acc *accepted) error {
	if a.jsonOut {
		return a.printJSON(acc)
	}
	res := acc.Acceptance
	if res.DryRun {
		for _, b := range res.Blockers {
			fmt.Fprintf(a.out, "blocked: %s\n", b)
		}
		if len(res.Uncommitted) > 0 {
			effect := "the real run refuses until they are committed or stashed, or --yes includes them"
			switch {
			case res.Batch != "":
				effect = "the merge queue's landing leaves them as they are, out of the acceptance commit (ADR-0178)"
			case a.yes:
				effect = "--yes includes them in the acceptance commit"
			}
			fmt.Fprintf(a.out, "uncommitted outside wip: %s; %s\n", strings.Join(res.Uncommitted, ", "), effect)
		}
		if res.OperatorOnly != "" {
			fmt.Fprintf(a.out, "%s: %s\n", res.OperatorOnly, strings.Join(res.Protected, ", "))
		}
		for _, u := range res.UndoesMain {
			if u.KeptReason != "" {
				fmt.Fprintf(a.out, "would keep the removal of the main branch's lines at %s: %s\n", u.Path, u.KeptReason)
			}
		}
		if res.Remote != nil {
			fmt.Fprintln(a.out, remotePlanLine(res))
		} else if res.Branch != "" {
			fmt.Fprintf(a.out, "would merge %s into the main branch and remove its worktree\n", res.Branch)
		}
		if res.CommitOnly {
			fmt.Fprintf(a.out, "would commit the acceptance of %s, done and archived but not committed: %s\n", res.ID, acceptSubject(res.ID, res.Epic))
		} else if e := res.Epic; e != nil {
			fmt.Fprintf(a.out, "would also move %s %s from %s to %s, following %s", e.ID, e.Title, e.From, e.To, e.Story)
			if e.To == workitem.Done {
				fmt.Fprint(a.out, ", and archive it")
			}
			fmt.Fprintln(a.out)
		}
		if len(res.ResolvedThreads) > 0 {
			fmt.Fprintf(a.out, "would resolve %s, open or answered on what it archives\n", strings.Join(res.ResolvedThreads, ", "))
		}
		for _, line := range acc.closing {
			fmt.Fprintf(a.out, "would close %s\n", line)
		}
		fmt.Fprintln(a.out, "dry run: nothing changed")
		return nil
	}
	verb := "accepted"
	if res.Resumed {
		verb = "completed the acceptance of"
	}
	if res.CommitOnly {
		fmt.Fprintf(a.out, "%s %s: done and archived already, committed", verb, res.ID)
	} else {
		fmt.Fprintf(a.out, "%s %s: done, %d items archived, committed", verb, res.ID, res.Archived)
	}
	remote := res.Remote != nil && res.Remote.Outcome != ""
	switch {
	case remote:
		fmt.Fprintf(a.out, ", %s, %s removed", remoteMergedLine(res.Remote), res.Branch)
	case res.Merged:
		fmt.Fprintf(a.out, ", %s merged and removed", res.Branch)
	}
	fmt.Fprintln(a.out, "; nothing released yet, run flai release --pending to publish")
	if remote && !res.Remote.BranchDeleted && res.Remote.BranchReason != "" {
		fmt.Fprintf(a.out, "origin's %s not deleted: %s\n", res.Branch, res.Remote.BranchReason)
	}
	for _, u := range res.UndoesMain {
		if u.KeptReason != "" {
			fmt.Fprintf(a.out, "kept the removal of the main branch's lines at %s; the reason is under ### Removals kept in the Notes of %s\n", u.Path, res.ID)
		}
	}
	if res.Batch != "" {
		fmt.Fprintf(a.out, "landed from merge queue batch %s, verified at %s\n", res.Batch, short(res.Verified))
	}
	if res.Evidence != nil {
		fmt.Fprintf(a.out, "accepted by the orchestrator at %s; its evidence is under ### Accepted by the orchestrator in the Notes of %s\n", short(res.Verified), res.ID)
	}
	if e := res.Epic; e != nil {
		fmt.Fprintf(a.out, "%s followed it from %s to %s", e.ID, e.From, e.To)
		if e.To == workitem.Done {
			fmt.Fprint(a.out, ", archived")
		}
		fmt.Fprintln(a.out)
	}
	if len(res.ResolvedThreads) > 0 {
		fmt.Fprintf(a.out, "resolved %s, open or answered on what it archived\n", strings.Join(res.ResolvedThreads, ", "))
	}
	for _, line := range acc.closing {
		fmt.Fprintf(a.out, "closed %s\n", line)
	}
	for _, n := range res.Overlaps {
		fmt.Fprintf(a.out, "told %s it overlaps: %s\n", n.ID, strings.Join(n.Paths, ", "))
	}
	return nil
}

func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return s
}
