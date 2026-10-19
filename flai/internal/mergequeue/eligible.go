// Package mergequeue works out which stories in review the merge queue would
// batch, in the order they last entered review, and why each other story in
// review would not, and reads the queue's state (ADR-0173). It verifies,
// merges, and moves nothing. flai queue, flai board, the MCP board, and the
// host read board.get each carry its Queue as queue.
package mergequeue

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
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
	KindUnreadable       = "unreadable"        // its last flai verify result or its record of undone paths cannot be read
)

// Queue is the merge queue as it stands: the stories in review that are
// eligible, in order, each other story in review with its reasons, and the
// state's enqueued stories, batch, and outcomes.
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
	Paths    []string `json:"paths,omitempty"`    // the paths undoing main's change
}

// OffUnderPullRequest is why the queue does not run under pull-request
// acceptance.
const OffUnderPullRequest = "the merge queue runs under local acceptance only: under acceptance.mode pull-request the forge's required checks gate each story, and a batch rebases each story, so every landing would stop on forge_stale (ADR-0169, ADR-0173)"

// Compute is the queue of repo's stories in review, reading each story
// branch's head and its record of undone paths with git through r, and the
// queue's state from the main checkout's .flai-cache. Under pull-request
// acceptance it is off, with no story listed. A story whose own records
// cannot be read is ineligible for it; the items, the threads, or the state
// that cannot be read are an error.
func Compute(r execx.Runner, repo *workitem.Repo) (Queue, error) {
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
	items, err := repo.List(false)
	if err != nil {
		return q, fmt.Errorf("list the items to find the stories in review for the merge queue: %w; fix the item file it names", err)
	}
	all, err := threads.List(repo)
	if err != nil {
		return q, fmt.Errorf("read the threads to check the stories in review for the merge queue: %w; fix the thread file it names", err)
	}
	enqueued := map[string]bool{}
	for _, e := range st.Enqueued {
		enqueued[workitem.CanonicalID(e.Story)] = true
	}
	var entries []Entry
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.Review {
			continue
		}
		e := Entry{Story: it.ID, Title: it.Title, EnteredReview: enteredReview(it), Enqueued: enqueued[workitem.CanonicalID(it.ID)]}
		var verifyReason *Reason
		e.Head, verifyReason = verified(r, repo, it.ID)
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
	return q, nil
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

// verified is the head of story's branch and nil when its last flai verify
// passed there; otherwise the reason it did not, with the head when it was
// read.
func verified(r execx.Runner, repo *workitem.Repo, story string) (string, *Reason) {
	rep, ok, err := verify.LastReport(repo, story)
	switch {
	case err != nil:
		return "", &Reason{Kind: KindUnreadable, Text: fmt.Sprintf("cannot read the last flai verify result of %s: %v", story, err)}
	case !ok:
		return "", &Reason{Kind: KindVerifyMissing, Text: fmt.Sprintf("%s has no flai verify result: run flai verify %s at the head of its branch", story, story)}
	case !rep.Passed:
		at := ""
		if rep.StoppedAt != "" {
			at = " at its step " + rep.StoppedAt
		}
		return "", &Reason{Kind: KindVerifyFailed, Text: fmt.Sprintf("the last flai verify of %s failed%s, at %s: fix it and run flai verify %s again", story, at, short(rep.Commit), story)}
	}
	branch := storygit.Branch(story)
	out, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", branch+"^{commit}")
	head := strings.TrimSpace(out)
	if err != nil || head == "" {
		why := "git names no commit"
		if err != nil {
			why = err.Error()
		}
		return "", &Reason{Kind: KindBranchUnread, Text: fmt.Sprintf("cannot read the head of %s to check that flai verify passed there: %s; a story is batched from its branch, so check that git can read it", branch, why)}
	}
	if rep.Commit != head {
		return head, &Reason{Kind: KindVerifyStale, Text: fmt.Sprintf("the last flai verify of %s passed at %s, but the head of %s is %s: run flai verify %s again at the head", story, short(rep.Commit), branch, short(head), story)}
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
