// Package mergequeue works out which stories in review the merge queue would
// batch, in the order they last entered review, and why each other story in
// review would not, and reads and writes the queue's state (ADR-0173). It
// verifies, merges, and moves nothing. flai queue, flai board, the MCP board,
// and the host read board.get each carry its Queue as queue.
package mergequeue

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Kinds of the reasons a story in review is not eligible.
const (
	KindVerifyMissing    = "verify_missing"    // flai verify has no result for the story
	KindVerifyFailed     = "verify_failed"     // its last flai verify failed
	KindVerifyStale      = "verify_stale"      // its last flai verify passed at a commit other than its branch's head
	KindBranchUnread     = "branch_unread"     // the head of its branch cannot be read
	KindCriteriaUnticked = "criteria_unticked" // an acceptance criterion is unticked
	KindThreadOpen       = "thread_open"       // a thread on the story or one of its tasks is not resolved
	KindUndoneMain       = "undone_main"       // a path is recorded as undoing main's change (ADR-0155)
	KindUnreadable       = "unreadable"        // its last flai verify result, its record of undone paths, or its kluster namespace's record cannot be read
	KindDirtyMain        = "dirty_main"        // its change would overwrite files the main checkout has not committed outside the shared set (ADR-0178, ADR-0248)
	KindKluster          = "kluster"           // the kluster host action is on and its namespace was not deployed, or its checks did not pass, at its head (ADR-0254 §5a)
)

// Queue is the merge queue as it stands: the stories in review that are
// eligible, in order, each other story in review with its reasons, the
// state's enqueued stories, batch, and outcomes, and the window the next
// batch fills in.
type Queue struct {
	// Mode is the project's acceptance mode, local or pull-request.
	Mode string `json:"mode"`
	// Off says the queue does not run in Mode, and OffReason why; then no
	// story is listed and the state is not read.
	Off       bool   `json:"off,omitempty"`
	OffReason string `json:"off_reason,omitempty"`
	// Eligible are the stories the queue would batch, by when they last
	// entered review, oldest first, then by ID.
	Eligible []Entry `json:"eligible"`
	// Ineligible are the other stories in review, in the same order, each
	// with every reason it is not eligible.
	Ineligible []Entry    `json:"ineligible"`
	Enqueued   []Enqueued `json:"enqueued"`
	Batch      *Batch     `json:"batch,omitempty"`
	Outcomes   []Outcome  `json:"outcomes"`
	// Held says the queue builds no batch until it is resumed, and HeldBy
	// and HeldAt, in UTC, who held it and when, as the state has them.
	Held   bool      `json:"held,omitempty"`
	HeldBy string    `json:"held_by,omitempty"`
	HeldAt time.Time `json:"held_at,omitzero"`
	// Preflight is the hold of the landing's preflight, as the state has
	// it (ADR-0234): when it began, each finding on the main checkout with
	// its file, line, rule, and message, and the thread flai serve opened
	// for it; nil while the preflight does not hold the queue.
	Preflight *PreflightHold `json:"preflight,omitempty"`
	// Window is the window the next batch fills in before it is built
	// (S-0552): set only while the queue is not held, by the operator or
	// the preflight, no batch is under way or the last one has a result
	// and was dropped, so that the next is built in its place, stories are
	// enqueued and eligible, fewer than merge_queue.batch_max,
	// merge_queue.coalesce is more than zero, and the window has not
	// ended; nil otherwise.
	Window *Window `json:"window,omitempty"`
}

// Window is the window a batch of the merge queue fills in (S-0552): it
// opens when the first of its stories was enqueued and lasts
// merge_queue.coalesce, and the batch is built when it ends or holds
// BatchMax stories, whichever comes first. A story enqueued while a batch
// ran, or left enqueued once one was dropped or landed, keeps when it was
// enqueued, so it opens no new window.
type Window struct {
	// Ends is when the window ends, in UTC: when the first of Stories was
	// enqueued, plus merge_queue.coalesce.
	Ends time.Time `json:"ends"`
	// Stories are the stories the batch would hold were it built now, in
	// the order they were enqueued.
	Stories []string `json:"stories"`
	// BatchMax is merge_queue.batch_max: with as many stories, the batch
	// is built at once.
	BatchMax int `json:"batch_max"`
}

// window is the window the batch of picks, the stories pick chose from
// enqueued at most limit, fills in at now, coalesce long; nil when there is
// no window to wait out: coalesce is not more than zero, picks is empty or
// holds limit stories, or the window ended. A story's enqueue time is the
// one enqueued records for it.
func window(picks []Entry, enqueued []Enqueued, limit int, coalesce time.Duration, now time.Time) *Window {
	if coalesce <= 0 || len(picks) == 0 || len(picks) >= limit {
		return nil
	}
	at := map[string]time.Time{}
	for _, en := range enqueued {
		at[workitem.CanonicalID(en.Story)] = en.At
	}
	var first time.Time
	stories := make([]string, 0, len(picks))
	for i, e := range picks {
		if t := at[workitem.CanonicalID(e.Story)]; i == 0 || t.Before(first) {
			first = t
		}
		stories = append(stories, e.Story)
	}
	ends := first.Add(coalesce).UTC()
	if !now.Before(ends) {
		return nil
	}
	return &Window{Ends: ends, Stories: stories, BatchMax: limit}
}

// Entry is a story in review as the queue sees it.
type Entry struct {
	Story string `json:"story"`
	Title string `json:"title"`
	// EnteredReview is when the story last entered review, in UTC.
	EnteredReview time.Time `json:"entered_review"`
	// Head is its branch's head, read only when its last flai verify
	// passed.
	Head string `json:"head,omitempty"`
	// Enqueued says the state has the story enqueued.
	Enqueued bool     `json:"enqueued,omitempty"`
	Reasons  []Reason `json:"reasons,omitempty"`
}

// Reason is one reason a story in review is not eligible: a stable kind, a
// sentence saying what fails and what to do, and what it names.
type Reason struct {
	Kind     string   `json:"kind"`
	Text     string   `json:"text"`
	Criteria []int    `json:"criteria,omitempty"` // the unticked criteria's numbers
	Threads  []string `json:"threads,omitempty"`  // the threads not resolved
	Paths    []string `json:"paths,omitempty"`    // the paths undoing main's change, or the uncommitted files the change would overwrite
	// Command is what clears the reason where one command does: for
	// KindDirtyMain, the git stash push that sets the files aside.
	Command string `json:"command,omitempty"`
}

// OffUnderPullRequest is why the queue does not run under pull-request
// acceptance.
const OffUnderPullRequest = "the merge queue runs under local acceptance only: under acceptance.mode pull-request the forge's required checks gate each story, and a batch rebases each story, so every landing would stop on forge_stale (ADR-0169, ADR-0173)"

// Batches reports whether the merge queue can batch for a project, as flai
// serve builds a batch: the merge-queue host action is on for it (actionOn),
// its manifest lets the orchestrator accept reviews, and acceptance is local.
// A hold does not count, since a held queue batches once resumed. While it
// can, flai verify leaves the tiers marked queue to the batch (ADR-0174).
func Batches(actionOn bool, m manifest.Manifest) bool {
	return actionOn && m.Orchestration.Permissions.AcceptReviews && !m.Acceptance.PullRequest()
}

// Compute is the queue of repo's stories in review, reading each story
// branch's head and its record of undone paths with git through r, and the
// queue's state from the main checkout's .flai-cache. Under pull-request
// acceptance it is off, with no story listed. A story whose own records
// cannot be read is ineligible for it; the items, the threads, or the state
// that cannot be read are an error. Its Window is as it stands now, read
// with the main branch's head; a head git cannot read leaves it nil. While
// no story is in review it reads no thread and runs no git (S-0604). It
// reads every branch head with one git for-each-ref and the main checkout's
// git status beside it, and both only while a story in review has a passing
// flai verify result, so the git it runs does not grow with the stories in
// review (S-0698).
func Compute(r execx.Runner, repo *workitem.Repo) (Queue, error) {
	return ComputeAt(r, repo, time.Now())
}

// ComputeAt is Compute with its Window as it stands at now, as flai queue
// reads it with its own clock.
func ComputeAt(r execx.Runner, repo *workitem.Repo, now time.Time) (Queue, error) {
	return ComputeWith(r, repo, now, Actions{})
}

// Actions are the host actions on for the project that gate a story's
// acceptance, and so its eligibility, as the caller reads them: Kluster, the
// kluster host action, and Checks, the checks host action (ADR-0254 §5a).
// The zero value has both off.
type Actions struct {
	Kluster bool
	Checks  bool
}

// ComputeWith is ComputeAt with the host actions act on for the project:
// while act.Kluster is on and the manifest has a kluster: block, a story in
// review is eligible only while kluster.AcceptGate passes it at the head of
// its branch, with act.Checks, as flai accept --batch refuses it otherwise.
func ComputeWith(r execx.Runner, repo *workitem.Repo, now time.Time, act Actions) (Queue, error) {
	q := Queue{Mode: repo.Manifest.Acceptance.ModeOrDefault(), Eligible: []Entry{}, Ineligible: []Entry{}, Enqueued: []Enqueued{}, Outcomes: []Outcome{}}
	if repo.Manifest.Acceptance.PullRequest() {
		q.Off, q.OffReason = true, OffUnderPullRequest
		return q, nil
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return q, err
	}
	q.Enqueued, q.Batch, q.Outcomes = st.Enqueued, st.Batch, st.Outcomes
	q.Held, q.HeldBy, q.HeldAt = st.Held, st.HeldBy, st.HeldAt
	q.Preflight = st.Preflight
	items, err := repo.List(false)
	if err != nil {
		return q, fmt.Errorf("list the items to find the stories in review for the merge queue: %w; fix the item file it names", err)
	}
	var review []*workitem.Item
	for _, it := range items {
		if it.Type == workitem.Story && it.Status == workitem.Review {
			review = append(review, it)
		}
	}
	if len(review) == 0 {
		// No story in review: no reason to work out, so no thread is read
		// and no git runs, and with nothing eligible there is no window.
		return q, nil
	}
	all, err := threads.List(repo)
	if err != nil {
		return q, fmt.Errorf("read the threads to check the stories in review for the merge queue: %w; fix the thread file it names", err)
	}
	reports := make([]lastVerify, len(review))
	anyPassed := false
	for i, it := range review {
		reports[i] = readLastVerify(repo, it.ID)
		anyPassed = anyPassed || reports[i].reason == nil
	}
	// Only a story whose last flai verify passed needs its branch's head,
	// and only its head can be held back by the main checkout's dirt, so
	// with none, no git runs; otherwise every head is one read, beside the
	// main checkout's status (S-0698).
	var g gitState
	if anyPassed {
		g = readGitState(r, repo)
	}
	enqueued := map[string]bool{}
	for _, e := range st.Enqueued {
		enqueued[workitem.CanonicalID(e.Story)] = true
	}
	var entries []Entry
	for i, it := range review {
		e := Entry{Story: it.ID, Title: it.Title, EnteredReview: enteredReview(it), Enqueued: enqueued[workitem.CanonicalID(it.ID)]}
		var verifyReason *Reason
		e.Head, verifyReason = g.verified(it.ID, reports[i])
		if verifyReason != nil {
			e.Reasons = append(e.Reasons, *verifyReason)
		}
		if reason := unticked(it); reason != nil {
			e.Reasons = append(e.Reasons, *reason)
		}
		if reason := openThreads(all, items, it.ID); reason != nil {
			e.Reasons = append(e.Reasons, *reason)
		}
		if reason := undoesMain(r, repo, it.ID); reason != nil {
			e.Reasons = append(e.Reasons, *reason)
		}
		if reason := g.dirt.holds(r, repo, it.ID, e.Head); reason != nil {
			e.Reasons = append(e.Reasons, *reason)
		}
		if reason := klusterHolds(r, repo, act, it.ID, e.Head); reason != nil {
			e.Reasons = append(e.Reasons, *reason)
		}
		entries = append(entries, e)
	}
	slices.SortFunc(entries, func(a, b Entry) int {
		return cmp.Or(a.EnteredReview.Compare(b.EnteredReview), strings.Compare(a.Story, b.Story))
	})
	for _, e := range entries {
		if len(e.Reasons) == 0 {
			q.Eligible = append(q.Eligible, e)
		} else {
			q.Ineligible = append(q.Ineligible, e)
		}
	}
	q.Window = filling(repo, q, g, now)
	return q, nil
}

// filling is the window q's next batch fills in at now, as RunBatch would
// pick it from the main branch's head, as g read it: nil while the queue is
// held, by the operator or the preflight, a batch is under way or waits to
// be landed or bisected, the queue has no window to wait out, or the main
// branch's head was not read. A story is eligible only once its flai verify
// passed, so g holds the heads whenever a story is.
func filling(repo *workitem.Repo, q Queue, g gitState, now time.Time) *Window {
	if q.Held || q.Preflight != nil || (q.Batch != nil && (q.Batch.Result == nil || q.Batch.Dropped == nil)) {
		return nil
	}
	coalesce := repo.Manifest.MergeQueueCoalesce()
	if coalesce <= 0 || len(enqueuedEligible(q)) == 0 || g.headsErr != nil {
		return nil
	}
	_, base, err := g.heads.Main()
	if err != nil {
		return nil
	}
	limit := repo.Manifest.BatchMax()
	return window(pick(q, base, limit), q.Enqueued, limit, coalesce, now)
}

// Board is the board with the merge queue beside it: what flai board --json,
// the MCP board, and the host read board.get each give.
type Board struct {
	workitem.BoardView
	// Queue is the merge queue as flai queue --json gives it; nil only when
	// it cannot be computed, which is logged, so the board is still served.
	Queue *Queue `json:"queue,omitempty"`
}

// OnBoard is v with repo's merge queue beside it, computed with Compute
// through r. A queue that cannot be computed is warned of on log, when there
// is one, and left off.
func OnBoard(ctx context.Context, r execx.Runner, repo *workitem.Repo, v workitem.BoardView, log *slog.Logger) Board {
	defer perf.Track(ctx, "merge.queue")()
	q, err := Compute(r, repo)
	if err != nil {
		if log != nil {
			log.Warn("the merge queue not read; left off the board", "component", "mergequeue", "err", err.Error())
		}
		return Board{BoardView: v}
	}
	return Board{BoardView: v, Queue: &q}
}

// enteredReview is when it last entered review, from its last transition to
// review; when it has none, when it entered its current state.
func enteredReview(it *workitem.Item) time.Time {
	for i := len(it.Transitions) - 1; i >= 0; i-- {
		if tr := it.Transitions[i]; tr.To == workitem.Review {
			if at, err := time.Parse(workitem.TimeFormat, tr.At); err == nil {
				return at.UTC()
			}
			break
		}
	}
	return it.EnteredAt().UTC()
}

// lastVerify is a story's last flai verify result as the queue reads it:
// the commit it passed at, or the reason it did not pass.
type lastVerify struct {
	commit string
	reason *Reason
}

// readLastVerify reads story's last flai verify result, which needs no git.
func readLastVerify(repo *workitem.Repo, story string) lastVerify {
	rep, ok, err := verify.LastReport(repo, story)
	switch {
	case err != nil:
		return lastVerify{reason: &Reason{Kind: KindUnreadable, Text: fmt.Sprintf("cannot read the last flai verify result of %s: %v", story, err)}}
	case !ok:
		return lastVerify{reason: &Reason{Kind: KindVerifyMissing, Text: fmt.Sprintf("%s has no flai verify result: run flai verify %s at the head of its branch", story, story)}}
	case !rep.Passed:
		at := ""
		if rep.StoppedAt != "" {
			at = " at its step " + rep.StoppedAt
		}
		return lastVerify{reason: &Reason{Kind: KindVerifyFailed, Text: fmt.Sprintf("the last flai verify of %s failed%s, at %s: fix it and run flai verify %s again", story, at, short(rep.Commit), story)}}
	}
	return lastVerify{commit: rep.Commit}
}

// gitState is the git the queue reads once for every story in review: each
// local branch's head and the branch checked out, or why they could not be
// read, and the main checkout's dirt.
type gitState struct {
	heads    storygit.Heads
	headsErr error
	dirt     dirt
}

// readGitState reads repo's branch heads with one git for-each-ref and,
// at the same time, the main checkout's uncommitted files with one git
// status, both through r (S-0698).
func readGitState(r execx.Runner, repo *workitem.Repo) gitState {
	var (
		d    storygit.Dirty
		dErr error
		done = make(chan struct{})
	)
	go func() {
		defer close(done)
		d, dErr = storygit.DirtyMain(r, repo)
	}()
	var g gitState
	g.heads, g.headsErr = storygit.ReadHeads(r, repo.MainRoot)
	<-done
	g.dirt = mainDirt(d, dErr, g.heads, g.headsErr)
	return g
}

// verified is the head of story's branch and nil when its last flai verify,
// last, passed there; otherwise the reason it did not, with the head when it
// was read.
func (g gitState) verified(story string, last lastVerify) (string, *Reason) {
	if last.reason != nil {
		return "", last.reason
	}
	branch := storygit.Branch(story)
	head, ok := g.heads.Commits[branch]
	if g.headsErr != nil || !ok {
		why := "git names no commit"
		if g.headsErr != nil {
			why = g.headsErr.Error()
		}
		return "", &Reason{Kind: KindBranchUnread, Text: fmt.Sprintf("cannot read the head of %s to check that flai verify passed there: %s; a story is batched from its branch, so check that git can read it", branch, why)}
	}
	if last.commit != head {
		return head, &Reason{Kind: KindVerifyStale, Text: fmt.Sprintf("the last flai verify of %s passed at %s, but the head of %s is %s: run flai verify %s again at the head", story, short(last.commit), branch, short(head), story)}
	}
	return head, nil
}

// unticked is the reason naming its unticked acceptance criteria, or nil
// when every one is ticked.
func unticked(it *workitem.Item) *Reason {
	var ns []int
	var names []string
	for _, c := range itemedit.Criteria(it.Body) {
		if !c.Ticked {
			ns = append(ns, c.N)
			names = append(names, strconv.Itoa(c.N))
		}
	}
	if len(ns) == 0 {
		return nil
	}
	what := "criterion " + names[0] + " is"
	if len(ns) > 1 {
		what = "criteria " + strings.Join(names, ", ") + " are"
	}
	return &Reason{Kind: KindCriteriaUnticked, Criteria: ns, Text: fmt.Sprintf("acceptance %s unticked on %s: tick each one verified with flai criteria tick %s <n>", what, it.ID, it.ID)}
}

// openThreads is the reason naming the threads not resolved on story or one
// of its tasks, in thread ID order, or nil when there is none.
func openThreads(all []*threads.Thread, items []*workitem.Item, story string) *Reason {
	on := map[string]bool{workitem.CanonicalID(story): true}
	for _, x := range items {
		if x.Type == workitem.Task && workitem.CanonicalID(x.Parent) == workitem.CanonicalID(story) {
			on[workitem.CanonicalID(x.ID)] = true
		}
	}
	var ids []string
	for _, th := range all {
		if th.Open() && th.Anchor.Item != "" && on[workitem.CanonicalID(th.Anchor.Item)] {
			ids = append(ids, th.ID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	return &Reason{Kind: KindThreadOpen, Threads: ids, Text: fmt.Sprintf("%s on %s or its tasks: resolve each before the queue takes it", plural(ids, "thread", "is not resolved", "are not resolved"), story)}
}

// undoesMain is the reason naming the paths story's record says its branch
// still undoes of main's change (ADR-0155), or nil when there is none. A
// story with no record runs no git.
func undoesMain(r execx.Runner, repo *workitem.Repo, story string) *Reason {
	if _, err := os.Stat(storygit.UndoneRecordPath(repo.CacheDir(), story)); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	undone, err := storygit.RecordedUndone(r, repo, story)
	if err != nil {
		return &Reason{Kind: KindUnreadable, Text: fmt.Sprintf("cannot read which paths %s undoes of the main branch's change: %v", storygit.Branch(story), err)}
	}
	if len(undone) == 0 {
		return nil
	}
	paths := make([]string, 0, len(undone))
	for _, u := range undone {
		paths = append(paths, u.Path)
	}
	slices.Sort(paths)
	return &Reason{Kind: KindUndoneMain, Paths: paths, Text: fmt.Sprintf("%s undoes the main branch's change to %s, and the queue never batches it: restore main's version with flai stream restore %s --path <path>, keep the removal by accepting it with flai accept %s --keep-removal <path> --reason \"<why>\", or send %s back to in-progress (ADR-0155)", storygit.Branch(story), strings.Join(paths, ", "), story, story, story)}
}

// dirt is the main checkout's uncommitted files outside the shared set and
// the main branch's head, as the queue's eligibility reads them once for
// every story in review.
type dirt struct {
	files []string
	main  string
}

// mainDirt is the main checkout's dirt from what was read of it: its
// uncommitted files outside the shared set, wip, design, and the shared
// paths, which the landing commits (ADR-0248), d unless dErr, and the main
// branch's head from heads unless headsErr. A checkout git cannot read has
// none to hold a story back; the stories' other reasons say what git cannot
// read.
func mainDirt(d storygit.Dirty, dErr error, heads storygit.Heads, headsErr error) dirt {
	if dErr != nil || headsErr != nil || len(d.Rest) == 0 {
		return dirt{}
	}
	_, head, err := heads.Main()
	if err != nil {
		return dirt{}
	}
	return dirt{files: d.Rest, main: head}
}

// holds is the reason naming the uncommitted files of the main checkout that
// story's change, its branch at head from the commit it shares with main,
// would overwrite, so that the queue holds it back until they are clean, or
// nil when there are none or its head is not known (ADR-0178).
func (d dirt) holds(r execx.Runner, repo *workitem.Repo, story, head string) *Reason {
	if len(d.files) == 0 || head == "" {
		return nil
	}
	paths, err := storygit.ChangedAmong(r, repo.MainRoot, d.main, head, d.files)
	if err != nil {
		return &Reason{Kind: KindUnreadable, Text: fmt.Sprintf("cannot read which uncommitted files of the main checkout %s would overwrite: %v", storygit.Branch(story), err)}
	}
	if len(paths) == 0 {
		return nil
	}
	return &Reason{Kind: KindDirtyMain, Paths: paths, Command: FreeCommand(paths), Text: HeldBackText(story, paths)}
}

// klusterHolds is the reason the kluster gate refuses story's acceptance at
// head, its branch's head, as git in r judges a record's commit there
// (ADR-0292 §3), naming the last job of its namespace and its outcome, or
// nil while act.Kluster is off, the manifest has no kluster: block, its head
// is not known, as for a story whose flai verify did not pass, which is
// ineligible already, or the gate passes it (ADR-0254 §5a).
func klusterHolds(r execx.Runner, repo *workitem.Repo, act Actions, story, head string) *Reason {
	if !act.Kluster || repo.Manifest.Kluster == nil || head == "" {
		return nil
	}
	ok, why, err := kluster.AcceptGate(r, repo.MainRoot, repo.Manifest, story, head, act.Checks)
	if err != nil {
		return &Reason{Kind: KindUnreadable, Text: fmt.Sprintf("cannot tell from the kluster record of %s's namespace whether it was deployed at %s: %v", story, short(head), err)}
	}
	if ok {
		return nil
	}
	what := "its namespace's deploy passed"
	if act.Checks {
		what = "its namespace's deploy and the checks run against it passed"
	}
	return &Reason{Kind: KindKluster, Text: fmt.Sprintf("%s: while the kluster host action is on, the queue batches %s only once %s at the head of its branch (ADR-0254 §5a); fix what failed and move it to review again so that flai serve deploys it there, or deploy it with flai kluster deploy %s in a shell on the host", why, story, what, story)}
}

// plural names ids after one or many, with the verb that agrees.
func plural(ids []string, noun, one, many string) string {
	if len(ids) == 1 {
		return fmt.Sprintf("%s %s %s", noun, ids[0], one)
	}
	return fmt.Sprintf("%ss %s %s", noun, strings.Join(ids, ", "), many)
}

// short is a commit's name abbreviated to twelve characters.
func short(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}
