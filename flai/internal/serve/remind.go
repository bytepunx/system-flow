package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Reminding the operator of a thread that awaits them (S-0498, ADR-0205,
// I-0151). A reminder per served project looks whenever the project's
// launcher does: every minute, and when its work items or threads change.
// With the manifest's threads.remind_after set, it reminds the operator,
// through the channel the host configuration keeps (Send), of each open
// thread the inbox marks awaiting them (ADR-0181), on a document or on a
// work item that is not done or cancelled, whose last entry is older than
// remind_after, and again every threads.remind_every while it still awaits
// them. A new entry starts the count over from that entry. With no channel
// set, the reminder is logged instead, and that no channel is set is warned
// of once per project per run.
//
// Nothing is written in the repository: a thread entry would end the
// agent's wait (ADR-0091) and change what the thread awaits (ADR-0090). What
// was reminded, and when, is kept per project in reminders.json in the serve
// folder, so that a restart repeats nothing; a thread that no longer awaits
// the operator is dropped from it at the next look.

// reminderReader is who the reminder reads the threads as: a name no author
// takes, so that no thread awaits it, "you", and inbox.Summarize says of
// each whether it awaits the operator, as the operator's own view does.
const reminderReader = "flai serve (reminders)"

func (d Dir) remindersFile() string { return filepath.Join(string(d), "reminders.json") }

// reminded is what was reminded of one thread: the time of the entry the
// reminder was for, and when it was last sent.
type reminded struct {
	For  time.Time `json:"for"`
	Sent time.Time `json:"sent"`
}

// remindersMu orders the read-modify-write of reminders.json by the
// projects' reminders, which look at once.
var remindersMu sync.Mutex

// reminders reads what was reminded, by project root, then by thread; none
// is an empty map.
func (d Dir) reminders() (map[string]map[string]reminded, error) {
	out := map[string]map[string]reminded{}
	data, err := os.ReadFile(d.remindersFile())
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w; remove it and flai serve reminds each thread that awaits the operator once more", d.remindersFile(), err)
	}
	return out, nil
}

// projectReminders is what was reminded in the project at root.
func (d Dir) projectReminders(root string) (map[string]reminded, error) {
	remindersMu.Lock()
	defer remindersMu.Unlock()
	all, err := d.reminders()
	if err != nil {
		return nil, err
	}
	if all[root] == nil {
		return map[string]reminded{}, nil
	}
	return all[root], nil
}

// recordReminders keeps state as what was reminded in the project at root;
// an empty state removes the project from the file.
func (d Dir) recordReminders(root string, state map[string]reminded) error {
	remindersMu.Lock()
	defer remindersMu.Unlock()
	all, err := d.reminders()
	if err != nil {
		return err
	}
	if len(state) == 0 {
		delete(all, root)
	} else {
		all[root] = state
	}
	return d.write(d.remindersFile(), all)
}

// reminder reminds the operator of the threads of one project that await
// them; its clock, channel, and sender are fields, so that it is tested
// without a host configuration or a network.
type reminder struct {
	e       Entry
	dir     Dir
	now     func() time.Time
	log     *slog.Logger
	channel func() config.Notify // read at each look; nil is no channel
	send    func(ctx context.Context, ch config.Notify, r Reminder) error
	looks   *sync.WaitGroup // a look in the background is counted here

	busy atomic.Bool
	mu   sync.Mutex
	// warned is whether this run has warned that no channel is set.
	warned bool
	said   map[string]bool // the problems warned of, so each is warned once
}

// newReminder is the reminder for project e; a look in the background is
// counted in looks, which flai serve waits for when it stops.
func newReminder(o Options, e Entry, looks *sync.WaitGroup) *reminder {
	return &reminder{e: e, dir: o.Dir, now: o.Now, log: o.Logger, channel: o.Notify, send: Send, looks: looks}
}

// soon looks in the background unless a look still runs: a send waits on
// the webhook or the command for seconds, and the project's loop that calls
// it must not.
func (r *reminder) soon(ctx context.Context) {
	if ctx.Err() != nil || !r.busy.CompareAndSwap(false, true) {
		return
	}
	r.looks.Add(1)
	go func() {
		defer r.looks.Done()
		defer r.busy.Store(false)
		r.look(ctx)
	}()
}

// look reminds the operator of each thread that is due a reminder, records
// what it reminded, and drops what no longer awaits them.
func (r *reminder) look(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	repo, err := workitem.Open(r.e.Root)
	if err != nil {
		r.warnOnce("thread reminders could not read the project", err)
		return
	}
	after, err := repo.Manifest.Threads.RemindAfterDuration()
	if err != nil {
		r.warnOnce("thread reminders off for a setting not valid", err)
		return
	}
	if after == 0 {
		return
	}
	every, err := repo.Manifest.Threads.RemindEveryDuration()
	if err != nil {
		r.warnOnce("thread reminders off for a setting not valid", err)
		return
	}
	all, err := threads.List(repo)
	if err != nil {
		r.warnOnce("thread reminders could not read the threads", err)
		return
	}
	was, err := r.dir.projectReminders(r.e.Root)
	if err != nil {
		r.warnOnce("thread reminders could not read what was reminded", err)
		return
	}
	state := map[string]reminded{}
	var ch config.Notify
	if r.channel != nil {
		ch = r.channel()
	}
	now := r.now().UTC()
	for _, th := range all {
		since, sum, ok := awaitsOperator(repo, th)
		if !ok {
			continue
		}
		prev, had := was[th.ID]
		if had {
			state[th.ID] = prev
		}
		if now.Sub(since) < after || (had && prev.For.Equal(since) && now.Sub(prev.Sent) < every) {
			continue
		}
		if ctx.Err() != nil {
			break
		}
		r.remind(ctx, ch, reminderOf(repo.Manifest.Name, th, sum, since, now))
		if ctx.Err() != nil {
			break // cut by serve stopping: the next run sends it
		}
		state[th.ID] = reminded{For: since, Sent: now}
	}
	if maps.EqualFunc(was, state, func(a, b reminded) bool { return a.For.Equal(b.For) && a.Sent.Equal(b.Sent) }) {
		return
	}
	if err := r.dir.recordReminders(r.e.Root, state); err != nil {
		r.warnOnce("thread reminders not recorded", err, "file", r.dir.remindersFile())
	}
}

// awaitsOperator answers, of an open thread that awaits the operator, as
// inbox.Summarize says (ADR-0181), on a document or on a work item that is
// not done or cancelled, when its last entry was written and its summary;
// ok is false for any other thread.
func awaitsOperator(repo *workitem.Repo, th *threads.Thread) (since time.Time, sum inbox.ThreadSummary, ok bool) {
	if !th.Open() {
		return time.Time{}, sum, false
	}
	sum = inbox.Summarize(repo, inbox.Reader{Agent: reminderReader}, th)
	if sum.Awaiting != "operator" {
		return time.Time{}, sum, false
	}
	if item := th.Anchor.Item; item != "" {
		// an item flai cannot read is not known to be finished
		if it, err := repo.Get(item); err == nil && (it.Status == workitem.Done || it.Status == workitem.Cancelled) {
			return time.Time{}, sum, false
		}
	}
	at := th.Created
	if es := th.Entries(); len(es) > 0 {
		at = es[len(es)-1].At
	}
	since, err := time.Parse(workitem.TimeFormat, at)
	if err != nil {
		return time.Time{}, sum, false
	}
	return since, sum, true
}

// reminderOf is the reminder of th, which has awaited the operator since
// since, made at now, in the project named project.
func reminderOf(project string, th *threads.Thread, sum inbox.ThreadSummary, since, now time.Time) Reminder {
	out := Reminder{Project: project, Thread: th.ID, Title: th.Title, Item: th.Anchor.Item,
		AwaitingSeconds: int64(now.Sub(since) / time.Second), Since: since, Href: ReminderHref(th.Anchor.Item, th.ID)}
	if rec := sum.PendingRecommendation; rec != nil {
		first, _, _ := strings.Cut(strings.TrimSpace(rec.Text), "\n")
		out.Recommendation = &ReminderRecommendation{Author: rec.Author, FirstLine: strings.TrimSpace(first), Source: rec.Source}
	}
	return out
}

// remind sends rem through ch and logs what became of it; with no channel
// set it logs the reminder, and warns once per run that no channel is set.
func (r *reminder) remind(ctx context.Context, ch config.Notify, rem Reminder) {
	hours := math.Round(float64(rem.AwaitingSeconds)/360) / 10
	attrs := []any{"component", "serve", "project", r.e.Key, "thread", rem.Thread, "item", rem.Item, "awaiting_hours", hours, "recommendation", rem.Recommendation != nil}
	if ch.URL == "" && strings.TrimSpace(ch.Command) == "" {
		if !r.warned {
			r.warned = true
			r.log.Warn("thread reminders have no channel", "component", "serve", "project", r.e.Key,
				"hint", "set one in a shell on the host with flai serve notify set --url <webhook> or --command <command>")
		}
		r.log.Info("thread awaits the operator", attrs...)
		return
	}
	if err := r.send(ctx, ch, rem); err != nil {
		r.log.Warn("thread reminder not sent", append(attrs, "err", err.Error())...)
		return
	}
	r.log.Info("thread reminded", attrs...)
}

// warnOnce warns of err with msg the first time it is seen.
func (r *reminder) warnOnce(msg string, err error, args ...any) {
	key := msg + "\x00" + err.Error()
	if r.said == nil {
		r.said = map[string]bool{}
	}
	if r.said[key] {
		return
	}
	r.said[key] = true
	r.log.Warn(msg, append([]any{"component", "serve", "project", r.e.Key, "err", err.Error()}, args...)...)
}
