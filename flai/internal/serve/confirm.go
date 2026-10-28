package serve

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Confirming an orchestrator's recommendation the operator left unanswered
// (S-0502, ADR-0206, I-0195). A confirmer per served project looks whenever
// the project's thread reminder does: every minute, and when its work items
// or threads change. While the project's orchestrate host action is on and
// the manifest's orchestration.permissions.confirm_after is set, it confirms,
// as flai serve (threads.ConfirmByTimeout), the pending recommendation of the
// orchestrator's on each open thread that a story's agent or a planner
// opened, that is no permission prompt's, on a document or on a work item
// that is not done or cancelled, when the recommendation is not marked for
// the operator alone, has waited confirm_after, and has no entry by the
// operator after it. The thread is then answered, and the launcher, which
// sees the change under threads, starts the asking agent again as on the
// operator's confirmation.
//
// Each confirmation is logged, and recorded in the orchestrator's activity
// document as an activity of no seconds and no cost naming the thread and
// its item: flai serve spent nothing on it, and the orchestrator's run is
// not charged for it.

// ConfirmTrigger is the trigger of the activity entry a confirmation by
// timeout is recorded as: readActivity measures the orchestrator's next
// activity from its own last entry, not from this one.
const ConfirmTrigger = "orchestration.permissions.confirm_after"

// confirmer confirms the recommendations of one project that waited
// confirm_after on the operator; its clock and the test of the orchestrate
// host action are fields, so that it is tested without a host
// configuration.
type confirmer struct {
	e   Entry
	now func() time.Time
	log *slog.Logger
	// orchestrating says whether the orchestrate host action is on for the
	// project, read at each look.
	orchestrating func() bool

	mu   sync.Mutex
	said map[string]bool // the problems warned of, so each is warned once
}

// newConfirmer is the confirmer for project e. The orchestrate host action
// is read as the orchestrator's look reads it, from o.Agent.
func newConfirmer(o Options, e Entry) *confirmer {
	return &confirmer{e: e, now: o.Now, log: o.Logger,
		orchestrating: func() bool { return o.Agent != nil && o.Agent(e.Root).Orchestrate }}
}

// look confirms each recommendation that is due, while the orchestrate host
// action is on and confirm_after is set. It runs in the caller's loop, as
// the orchestrator's look does: it writes a thread and the activity
// document, and waits on no network.
func (c *confirmer) look(ctx context.Context) {
	if ctx.Err() != nil || !c.orchestrating() {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	repo, err := workitem.Open(c.e.Root)
	if err != nil {
		c.warnOnce("recommendations not confirmed: the project could not be read", err)
		return
	}
	after, err := repo.Manifest.Orchestration.Permissions.ConfirmAfterDuration()
	if err != nil {
		c.warnOnce("recommendations not confirmed for a setting not valid", err)
		return
	}
	if after == 0 {
		return
	}
	all, err := threads.List(repo)
	if err != nil {
		c.warnOnce("recommendations not confirmed: the threads could not be read", err)
		return
	}
	now := c.now().UTC()
	for _, th := range all {
		if ctx.Err() != nil {
			return
		}
		rec, at, ok := dueConfirmation(repo, th, after, now)
		if !ok {
			continue
		}
		c.confirm(repo, th, rec, at, after, now)
	}
}

// dueConfirmation answers, of th, the orchestrator's recommendation flai
// serve confirms at now and when it was written; ok is false when th does
// not qualify (ADR-0206 item 3). threads.ConfirmByTimeout refuses again what
// it can tell, the operator-only mark, the age, and the operator's entry
// since; they are tested here first so that a thread that is not due is
// passed over without a warning.
func dueConfirmation(repo *workitem.Repo, th *threads.Thread, after time.Duration, now time.Time) (rec *threads.Entry, at time.Time, ok bool) {
	if !th.Open() || !th.OpenedByAgent() || permissionPrompt(th) {
		return nil, at, false
	}
	rec = th.PendingRecommendation()
	if rec == nil || !workitem.IsOrchestrator(rec.Author) || rec.OperatorOnly {
		return nil, at, false
	}
	at, err := time.Parse(workitem.TimeFormat, rec.At)
	if err != nil || now.Sub(at) < after {
		return nil, at, false
	}
	if item := th.Anchor.Item; item != "" {
		// an item flai cannot read is not known to be finished
		if it, err := repo.Get(item); err == nil && (it.Status == workitem.Done || it.Status == workitem.Cancelled) {
			return nil, at, false
		}
	}
	story := threads.StoryOf(repo, th)
	if th.WroteSince(*rec, func(name string) bool { return threads.OperatorName(repo, story, name) }) {
		return nil, at, false
	}
	return rec, at, true
}

// permissionPrompt says whether th is the thread a permission prompt opens
// for a write it holds, which only the operator answers (ADR-0206 item 3).
// ask.Asker.Hold titles it "Allow <tool> <path>?" and exports no test of it,
// so the title is matched here.
func permissionPrompt(th *threads.Thread) bool {
	return strings.HasPrefix(th.Title, "Allow ") && strings.HasSuffix(th.Title, "?")
}

// confirm confirms rec, written at at, on th, and logs the confirmation and
// records it in the orchestrator's activity document, with the trigger
// ConfirmTrigger. A confirmation that
// fails is warned of once per thread and error.
func (c *confirmer) confirm(repo *workitem.Repo, th *threads.Thread, rec *threads.Entry, at time.Time, after time.Duration, now time.Time) {
	if _, err := threads.ConfirmByTimeout(repo, th.ID, after, now); err != nil {
		c.warnOnce("recommendation not confirmed by timeout", err, "thread", th.ID, "item", th.Anchor.Item)
		return
	}
	waited := math.Round(now.Sub(at).Hours()*10) / 10
	c.log.Info("recommendation confirmed by timeout", "component", "serve", "project", c.e.Key, "thread", th.ID, "item", th.Anchor.Item,
		"waited_hours", waited, "confirm_after", after.String())
	items := []string{th.ID}
	if th.Anchor.Item != "" {
		items = append(items, th.Anchor.Item)
	}
	setting := strings.TrimSpace(repo.Manifest.Orchestration.Permissions.ConfirmAfter)
	entry := workitem.ActivityEntry{At: now, Items: items, Trigger: ConfirmTrigger, Summary: fmt.Sprintf(
		"flai serve confirmed the orchestrator's recommendation of %s on %s after %s with no answer from the operator (orchestration.permissions.confirm_after)",
		rec.At, th.ID, setting)}
	// appended as is, with no seconds and no cost: LogActivity would measure
	// the orchestrator's run and charge its time to this entry
	if _, err := repo.AppendActivity(workitem.ActivityOrchestrator, entry); err != nil {
		c.warnOnce("confirmation by timeout not logged in the orchestrator's activity", err, "thread", th.ID, "item", th.Anchor.Item)
	}
}

// warnOnce warns of err with msg and args the first time the three are seen
// together. The error is told by its innermost cause, so that a write that
// fails every minute on a temporary file of a new name each time is warned
// of once.
func (c *confirmer) warnOnce(msg string, err error, args ...any) {
	cause := err
	for inner := errors.Unwrap(cause); inner != nil; inner = errors.Unwrap(cause) {
		cause = inner
	}
	key := fmt.Sprint(msg, "\x00", cause.Error(), "\x00", args)
	if c.said == nil {
		c.said = map[string]bool{}
	}
	if c.said[key] {
		return
	}
	c.said[key] = true
	c.log.Warn(msg, append([]any{"component", "serve", "project", c.e.Key, "err", err.Error()}, args...)...)
}
