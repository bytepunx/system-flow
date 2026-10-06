package serve

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Planning again on its own (S-0211, ADR-0084).
//
// With the plan host action on, flai serve plans again without the operator
// asking. A replanner per served project looks whenever the project's
// launcher does: when its work items or threads change, when an agent ends,
// and every minute. It acts on three triggers.
//
//   - An edit. Someone other than a planner edits a backlog or ready story's
//     goal, criteria, or touches, and the story has a forecast or a cost of
//     delay value older than the edit; or edits its cost of delay's inputs
//     after its value. The planner is queued for the story. A story never
//     planned, with neither figure, is not.
//   - Work ahead completes or the order changes: a story is accepted or
//     cancelled, or the pull order changes. planning.replan decides: never
//     does nothing; deterministic, the default, plays the board out again
//     for every ready and backlog story with a forecast, keeping its own
//     duration, and writes the delivery and basis where the delivery moved,
//     stamped by flai, in one commit; agent does the same and queues the
//     planner for each story it wrote.
//   - planning.schedule. Each time it comes round, the planner is queued for
//     every ready story. Times missed in between come round once.
//
// The queue holds a story once: a trigger for a story already queued is
// added to its entry. It runs one planner at a time per project, with the
// launcher of flai serve, so that the run's end is seen at once, through the
// checks Plan makes; an entry they refuse is dropped. The run records its
// triggers, joined by semicolons, and its activity entry says them.
//
// What it has seen is kept in memory only: edits, moves, and scheduled times
// from its start are acted on, and none that passed while flai serve was
// down. While the plan host action is off it acts on nothing and keeps up
// with what happens, so that turning it on does not act on the past.

// replanFields are the fields whose edit makes a planned story due for the
// planner again (ADR-0084), in the order a trigger names them; an edit of
// the cost of delay is judged apart.
var replanFields = []string{"goal", "criteria", "touches"}

// costOfDelay is the field an edit of a cost of delay names.
const costOfDelay = "cost_of_delay"

// triggerReordered is the trigger of a change of the pull order.
const triggerReordered = "reordered"

// queued is a story the planner is queued for, and what queued it, in the
// order said.
type queued struct {
	item     string
	triggers []string
}

// moved is a story whose forecast delivery a replay moved, and the delivery
// and basis it moved to.
type moved struct {
	id, delivery, basis string
}

// replanner plans again for one project on its own (S-0211, ADR-0084).
type replanner struct {
	o       Options
	e       Entry
	starter *launcher

	mu sync.Mutex
	// since is the second from which changes and edits are new, and seen
	// those stamped then or later that have been seen, by key, with when.
	since time.Time
	seen  map[string]time.Time
	order []string // the pull order as last seen
	// cadence follows planning.schedule.
	cadence cadence
	queue   []queued
	// started is the item of the planner run this replanner last started.
	started string
	said    map[string]bool // the manifest's refusals warned of
}

// newReplanner is the replanner for project e, starting planners with
// starter, the project's launcher. What happened before now is not news to
// it.
func newReplanner(o Options, e Entry, starter *launcher) *replanner {
	now := o.Now().UTC()
	r := &replanner{o: o, e: e, starter: starter, since: now.Truncate(time.Second), seen: map[string]time.Time{}, cadence: cadence{mark: now}, said: map[string]bool{}}
	if repo, err := workitem.Open(e.Root); err == nil {
		if board, err := repo.LoadBoard(); err == nil {
			r.order = slices.Clone(board.Order)
		}
		r.cadence.spec = strings.TrimSpace(repo.Manifest.Planning.Schedule)
	}
	return r
}

// look acts on what happened since the last look: it queues the planner on
// an edit and on the schedule, replans when work ahead completes or the
// order changes, and starts the next queued planner when none it started
// runs.
func (r *replanner) look(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.o.Now().UTC()
	repo, err := workitem.Open(r.e.Root)
	if err != nil {
		r.warnOnce("replanner could not read the project", err)
		return
	}
	items, err := repo.List(true)
	if err != nil {
		r.warnOnce("replanner could not read the project", err)
		return
	}
	board, err := repo.LoadBoard()
	if err != nil {
		r.warnOnce("replanner could not read the project", err)
		return
	}
	var changes []workitem.Change
	var edits []itemedit.Notice
	changes, edits, r.since, r.seen = news(items, itemedit.Notices(repo), r.since, r.seen, now)
	reordered := workitem.Reordered(r.order, board.Order)
	r.order = slices.Clone(board.Order)
	sched := r.schedule(repo.Manifest.Planning, now)
	if r.o.Agent == nil || !r.o.Agent(r.e.Root).Plan {
		r.queue, r.cadence.mark = nil, now
		return
	}
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	for _, n := range edits {
		if t := editTrigger(byID[n.ID], n); t != "" {
			r.queue = enqueue(r.queue, n.ID, t)
		}
	}
	ahead := moveTriggers(changes)
	if reordered {
		ahead = append(ahead, triggerReordered)
	}
	if len(ahead) > 0 {
		r.replan(repo, items, board, now, ahead)
	}
	if r.cadence.due(sched, now) {
		for _, id := range workitem.PullSequence(board.Order, items, workitem.Ready) {
			r.queue = enqueue(r.queue, id, scheduleTrigger(sched))
		}
	}
	r.drain(ctx)
}

// schedule is planning.schedule parsed, nil when it is unset or cannot be
// parsed, which is warned of once. A schedule set or changed comes round
// first at its next time from now.
func (r *replanner) schedule(plan manifest.Planning, now time.Time) *cron.Schedule {
	r.cadence.see(plan.Schedule, now)
	sched, err := plan.PlanSchedule()
	if err != nil {
		r.warnOnce("planning schedule not valid", err)
		return nil
	}
	return sched
}

// replan plays the board out again after triggers as planning.replan says:
// it writes the forecast deliveries that moved, commits them in one commit
// unless autocommit is off, and under agent queues the planner for each
// story it wrote. A policy that is not valid is warned of once, and nothing
// is done.
func (r *replanner) replan(repo *workitem.Repo, items []*workitem.Item, board *workitem.Board, now time.Time, triggers []string) {
	policy, err := repo.Manifest.Planning.ReplanPolicy()
	if err != nil {
		r.warnOnce("planning replan not valid", err)
		return
	}
	if policy == manifest.ReplanNever {
		return
	}
	moves, errs := replays(items, board, repo.Manifest.Planning, now)
	for _, err := range errs {
		r.starter.warn("forecast not replayed", "err", err)
	}
	var files, written []string
	for _, m := range moves {
		delivery, basis := m.delivery, m.basis
		res, err := itemedit.Apply(repo, gitOf(r.o), m.id, itemedit.Change{Forecast: &itemedit.ForecastEdit{Delivery: &delivery, Basis: &basis}},
			itemedit.Options{By: "flai", NoCommit: true, Now: now})
		if err != nil {
			r.starter.warn("forecast not replanned", "item", m.id, "err", err)
			continue
		}
		written = append(written, m.id)
		for _, f := range res.Files {
			if !slices.Contains(files, f) {
				files = append(files, f)
			}
		}
	}
	why := strings.Join(triggers, "; ")
	if len(files) > 0 && repo.Manifest.Autocommit() {
		dir := repo.Root
		if repo.MainRoot != "" {
			dir = repo.MainRoot
		}
		if err := commitFiles(gitOf(r.o), dir, files, "chore: replan forecasts after "+why); err != nil {
			r.starter.warn("replanned forecasts not committed", "files", len(files), "err", err)
		}
	}
	r.starter.log("forecasts replanned", "policy", policy, "moved", len(written), "triggers", why)
	if policy == manifest.ReplanAgent {
		for _, id := range written {
			for _, t := range triggers {
				r.queue = enqueue(r.queue, id, t)
			}
		}
	}
}

// drain starts the planner for the first queued story with no planner
// running, unless a planner this replanner started still runs. An entry the
// checks refuse is dropped and logged, and the next one tried; one that
// cannot be read is dropped and warned of.
func (r *replanner) drain(ctx context.Context) {
	st := r.o.Dir.AgentStates()[r.e.Root]
	if r.started != "" && st.Plans[r.started].running() {
		return
	}
	r.started = ""
	for {
		i := pick(r.queue, func(item string) bool { return st.Plans[item].running() })
		if i < 0 {
			return
		}
		q := r.queue[i]
		r.queue = slices.Delete(r.queue, i, i+1)
		trigger := strings.Join(q.triggers, "; ")
		ps, err := planCheck(r.o, r.e, q.item)
		var no *Refused
		switch {
		case errors.As(err, &no):
			r.starter.log("queued planner dropped", "item", q.item, "trigger", trigger, "reason", no.Why)
			continue
		case err != nil:
			r.starter.warn("queued planner dropped", "item", q.item, "trigger", trigger, "err", err)
			continue
		}
		r.starter.mu.Lock()
		r.starter.plan(ctx, ps.cfg, ps.item, ps.agent, trigger)
		r.starter.mu.Unlock()
		r.started = ps.item
		return
	}
}

// warnOnce warns of err with msg the first time it is seen.
func (r *replanner) warnOnce(msg string, err error) {
	if key := msg + "\x00" + err.Error(); !r.said[key] {
		r.said[key] = true
		r.starter.warn(msg, "err", err)
	}
}

// news are the changes in items and the edit notices that are new to a
// watermark at since that has seen seen, oldest first, and the watermark
// past them at now. A change or notice is new when it is stamped in since's
// second or later and has not been seen; the keys seen are kept while they
// are stamped in now's second or later.
func news(items []*workitem.Item, notices []itemedit.Notice, since time.Time, seen map[string]time.Time, now time.Time) ([]workitem.Change, []itemedit.Notice, time.Time, map[string]time.Time) {
	next := now.UTC().Truncate(time.Second)
	nextSeen := map[string]time.Time{}
	for k, at := range seen {
		if !at.Before(next) {
			nextSeen[k] = at
		}
	}
	fresh := func(key, stamp string) bool {
		at, err := time.Parse(workitem.TimeFormat, stamp)
		if err != nil || at.Before(since) {
			return false
		}
		if _, ok := seen[key]; ok {
			return false
		}
		if !at.Before(next) {
			nextSeen[key] = at
		}
		return true
	}
	var changes []workitem.Change
	for _, c := range workitem.Changes(items, since, "", nil) {
		if fresh(c.Key(), c.At) {
			changes = append(changes, c)
		}
	}
	var edits []itemedit.Notice
	for _, n := range notices {
		if fresh(noticeKey(n), n.At) {
			edits = append(edits, n)
		}
	}
	return changes, edits, next, nextSeen
}

// noticeKey identifies an edit notice, for telling apart edits within one
// second.
func noticeKey(n itemedit.Notice) string {
	return n.ID + "|edited|" + strings.Join(n.Changed, ",") + "|" + n.By + "|" + n.At
}

// editTrigger is the trigger an edit gives story it: "edited", the fields
// that make it due, and who edited; "" when it gives none (ADR-0084). It
// gives none for an item that is not a backlog or ready story on the board,
// and for a planner's edit. A goal, criteria, or touches edit makes it due
// when its forecast or its cost of delay value is older than the edit; a
// cost of delay edit, when its inputs are now newer than its value. A stamp
// that cannot be read is not older.
func editTrigger(it *workitem.Item, n itemedit.Notice) string {
	if it == nil || it.Archived || it.Type != workitem.Story || (it.Status != workitem.Backlog && it.Status != workitem.Ready) ||
		strings.HasPrefix(n.By, workitem.ActivityPlanner+"-") {
		return ""
	}
	at, err := time.Parse(workitem.TimeFormat, n.At)
	if err != nil {
		return ""
	}
	var fields []string
	if older(forecastAt(it), at) || older(valueAt(it), at) {
		for _, f := range replanFields {
			if slices.Contains(n.Changed, f) {
				fields = append(fields, f)
			}
		}
	}
	if slices.Contains(n.Changed, costOfDelay) && it.CostOfDelay.Stale() {
		fields = append(fields, costOfDelay)
	}
	if len(fields) == 0 {
		return ""
	}
	return "edited " + strings.Join(fields, ", ") + " by " + orSomeone(n.By)
}

// forecastAt is when a story's forecast was written, "" when it has none.
func forecastAt(it *workitem.Item) string {
	if it.Forecast.IsZero() {
		return ""
	}
	return it.Forecast.At
}

// valueAt is when a story's cost of delay value was set, "" when it has
// none.
func valueAt(it *workitem.Item) string {
	if it.CostOfDelay == nil || it.CostOfDelay.Value == nil {
		return ""
	}
	return it.CostOfDelay.At
}

// older reports whether the stamp is a time before t.
func older(stamp string, t time.Time) bool {
	at, err := time.Parse(workitem.TimeFormat, stamp)
	return err == nil && at.Before(t)
}

// moveTriggers are the triggers among changes of work ahead completing, in
// order: accepted for a story moved to done, cancelled for one moved to
// cancelled. A task's or an epic's move is none.
func moveTriggers(changes []workitem.Change) []string {
	var out []string
	for _, c := range changes {
		if c.Type != workitem.Story || c.Kind != workitem.Moved {
			continue
		}
		switch c.To {
		case workitem.Done:
			out = append(out, "accepted "+c.ID)
		case workitem.Cancelled:
			out = append(out, "cancelled "+c.ID)
		}
	}
	return out
}

// scheduleDue reports whether sched has come round since mark by now: its
// next time after mark is not after now. Several times missed come round
// once.
func scheduleDue(sched cron.Schedule, mark, now time.Time) bool {
	next := sched.Next(mark)
	return !next.IsZero() && !next.After(now)
}

// cadence follows a manifest schedule from look to look, for the replanner
// and the analyzer's schedule (S-0211, S-0223): spec is the schedule as last
// seen, and mark when it last came round or was set.
type cadence struct {
	spec string
	mark time.Time
}

// see takes spec as the schedule at now: one set or changed comes round
// first at its next time from now.
func (c *cadence) see(spec string, now time.Time) {
	if spec = strings.TrimSpace(spec); spec != c.spec {
		c.spec, c.mark = spec, now
	}
}

// due reports whether sched, nil when there is none, has come round since
// the mark by now, and marks now when it has, so that times missed in
// between come round once.
func (c *cadence) due(sched *cron.Schedule, now time.Time) bool {
	if sched == nil || !scheduleDue(*sched, c.mark, now) {
		return false
	}
	c.mark = now
	return true
}

// scheduleTrigger is the trigger of a run a schedule started: "schedule"
// and the schedule as written, as "schedule daily".
func scheduleTrigger(sched *cron.Schedule) string {
	return "schedule " + sched.String()
}

// enqueue adds trigger for item to q: to item's entry when it has one and
// the trigger is not in it yet, and as a new entry at the end when it has
// none.
func enqueue(q []queued, item, trigger string) []queued {
	for i := range q {
		if q[i].item == item {
			if !slices.Contains(q[i].triggers, trigger) {
				q[i].triggers = append(q[i].triggers, trigger)
			}
			return q
		}
	}
	return append(q, queued{item: item, triggers: []string{trigger}})
}

// pick is the index of the first entry in q whose item has no planner
// running, -1 when there is none.
func pick(q []queued, running func(item string) bool) int {
	for i, e := range q {
		if !running(e.item) {
			return i
		}
	}
	return -1
}

// replays plays the board out again for every ready and backlog story with a
// forecast duration, in pull order, each keeping its own duration
// (planning.Replay), and returns those whose delivery moved with the new
// delivery and basis, and the replays that failed. Every story is played out
// against the same items: the writes keep the durations the others are
// played out with.
func replays(items []*workitem.Item, board *workitem.Board, plan manifest.Planning, now time.Time) ([]moved, []error) {
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	var order []string
	if board != nil {
		order = board.Order
	}
	var out []moved
	var errs []error
	for _, id := range append(workitem.PullSequence(order, items, workitem.Ready), workitem.PullSequence(order, items, workitem.Backlog)...) {
		it := byID[id]
		if it.Forecast.IsZero() || it.Forecast.Duration == "" {
			continue
		}
		res, err := planning.Replay(items, board, plan, id, now)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !sameStamp(res.Delivery, it.Forecast.Delivery) {
			out = append(out, moved{id: id, delivery: res.Delivery, basis: res.Basis})
		}
	}
	return out, errs
}

// sameStamp reports whether two timestamps name the same time, or are the
// same text when either cannot be read.
func sameStamp(a, b string) bool {
	ta, errA := time.Parse(workitem.TimeFormat, a)
	tb, errB := time.Parse(workitem.TimeFormat, b)
	if errA != nil || errB != nil {
		return a == b
	}
	return ta.Equal(tb)
}

// commitFiles commits files, and nothing else, removals included, on the
// current branch of the checkout at dir with msg. A dir that is not in a git
// repository has nothing to commit to.
func commitFiles(r execx.Runner, dir string, files []string, msg string) error {
	if _, err := r.Run(dir, "git", "rev-parse", "--is-inside-work-tree"); err != nil {
		return nil //nolint:nilerr // no repository means nothing to commit to, by design
	}
	if _, err := r.Run(dir, "git", append([]string{"add", "-A", "--"}, files...)...); err != nil {
		return fmt.Errorf("git add in %s: %w", dir, err)
	}
	if _, err := r.Run(dir, "git", append([]string{"commit", "-q", "-m", msg, "--"}, files...)...); err != nil {
		return fmt.Errorf("git commit in %s: %w", dir, err)
	}
	return nil
}
