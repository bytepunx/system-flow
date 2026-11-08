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
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Planning again on its own (S-0211, ADR-0084).
//
// With the plan host action on, flai serve plans again without the operator
// asking. A replanner per served project looks whenever the project's
// launcher does: when its work items or threads change, when an agent ends,
// and every minute. At each look it reads the planner's role file, the
// project's or the baseline flai carries (roles.Builtin), and queues the
// planner only for the trigger entries the file names, judged by matching as
// the defined roles' launchers judge theirs (ADR-0250 §4). The baseline's
// entries are these.
//
//   - {event: edited, when: stale-figures}. Someone other than a planner or
//     the story's own agent edits a backlog or ready story's goal, criteria,
//     or touches, and the story has a forecast or a cost of delay value older
//     than the edit; or edits its cost of delay's inputs after its value
//     (roles.StaleFigures). The planner is queued for the story. A story
//     never planned, with neither figure, is not.
//   - {event: accepted}, {event: cancelled}, and {event: reordered}, each
//     when replan-agent: work ahead completes or the order changes, and
//     planning.replan is agent. The planner is queued for each story the
//     replay moved (roles.ReplanAgent).
//   - {schedule: planning.schedule, when: each-ready-story}. Each time the
//     setting comes round, the planner is queued for every ready story, in
//     pull order (roles.EachReadyStory). Times missed in between come round
//     once. Each setting the file's schedule entries name keeps a cadence of
//     its own.
//
// An event entry with no condition queues the event's item, when it is of a
// kind the planner takes, and never for an edit by a planner or by the
// replay; {event: reordered} has no item, and queues nothing so. An entry the
// replanner does not serve, as flai check reports, queues nothing.
//
// Replaying the forecasts is flai's own work, which no trigger entry gates:
// when a story is accepted or cancelled, or the pull order changes,
// planning.replan decides. never does nothing; deterministic, the default,
// plays the board out again for every ready and backlog story with a
// forecast, keeping its own duration, and writes the delivery and basis
// where the delivery moved, stamped by flai, in one commit; agent does the
// same, and replan-agent's entries queue the planner for each story it
// wrote.
//
// The queue holds a story once: a trigger for a story already queued is
// added to its entry. It runs one planner at a time per project, with the
// launcher of flai serve, so that the run's end is seen at once, through the
// checks Plan makes; an entry they refuse is dropped, as is one for a story
// its agent works: in progress or in review, or with its agent running
// (I-0098). One whose start the budget state of the planner's provider holds
// stays queued, in its place, until the state lets it start (S-0534). The
// run records its triggers, joined by semicolons, and its activity entry
// says them.
//
// What it has seen is kept in memory only: edits, moves, and scheduled times
// from its start are acted on, and none that passed while flai serve was
// down. While the plan host action is off it acts on nothing and keeps up
// with what happens, so that turning it on does not act on the past.

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
	// cadences follow the settings the planner's schedule entries name, by
	// setting.
	cadences map[string]*cadence
	queue    []queued
	// started is the item of the planner run this replanner last started.
	started string
	said    map[string]bool // the manifest's refusals warned of, and the holds logged (heldOnce)
}

// newReplanner is the replanner for project e, starting planners with
// starter, the project's launcher. What happened before now is not news to
// it.
func newReplanner(o Options, e Entry, starter *launcher) *replanner {
	now := o.Now().UTC()
	r := &replanner{o: o, e: e, starter: starter, since: now.Truncate(time.Second), seen: map[string]time.Time{}, cadences: map[string]*cadence{}, said: map[string]bool{}}
	if repo, err := workitem.Open(e.Root); err == nil {
		if board, err := repo.LoadBoard(); err == nil {
			r.order = slices.Clone(board.Order)
		}
		if planner, _, err := roles.Builtin(repo.Root, repo.Manifest.Layout["design"], workitem.ActivityPlanner); err == nil {
			r.cadences = startCadences(planner, repo.Manifest, now)
		}
	}
	return r
}

// look acts on what happened since the last look, as the planner's role
// file says: it queues the planner for the event and schedule entries the
// file names, replays the forecasts when work ahead completes or the order
// changes, and starts the next queued planner when none it started runs.
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
	planner, _, err := roles.Builtin(repo.Root, repo.Manifest.Layout["design"], workitem.ActivityPlanner)
	if err != nil {
		r.warnOnce("planner role file not read", err) // the replay goes on; no entry queues the planner
	}
	var changes []workitem.Change
	var edits []itemedit.Notice
	changes, edits, r.since, r.seen = news(items, itemedit.Notices(repo), r.since, r.seen, now)
	reordered := workitem.Reordered(r.order, board.Order)
	r.order = slices.Clone(board.Order)
	scheds := r.schedules(planner, repo.Manifest, now)
	var cfg AgentConfig
	if r.o.Agent != nil {
		cfg = r.o.Agent(r.e.Root)
	}
	if !cfg.Plan {
		r.queue = nil
		markAll(r.cadences, now)
		return
	}
	byID := itemsByID(items)
	for _, n := range edits {
		r.queueEdit(planner, byID[n.ID], n, cfg.Name)
	}
	if ahead := aheadOf(changes, reordered); len(ahead) > 0 {
		policy, written := r.replan(repo, items, board, now, triggersOf(ahead))
		for _, a := range ahead {
			r.queueAhead(planner, a, policy, written)
		}
	}
	for _, setting := range scheduleSettings(planner) {
		sched, ok := scheds[setting]
		if !ok || !r.cadences[setting].due(sched, now) {
			continue
		}
		for _, t := range matching(planner.Trigger, roles.Occurrence{Schedule: setting}) {
			if t.When != roles.CondEachReadyStory {
				continue // flai check reports it: a schedule gives the planner no item
			}
			for _, id := range roles.EachReadyStory(board.Order, items) {
				r.queue = enqueue(r.queue, id, scheduleTrigger(sched))
			}
		}
	}
	r.drain(ctx, byID)
}

// queueEdit queues the planner, whose role file is planner, for the edit n
// of it, as read after the edit, under each edited entry of the file that n
// matches: under stale-figures with editTrigger's trigger, and under no
// condition, or one such as draft-changed, with the fields the edit named,
// unless a planner or the replay made the edit, or it is of no kind the
// planner takes. agent is the configured name of the story agents.
func (r *replanner) queueEdit(planner roles.Role, it *workitem.Item, n itemedit.Notice, agent string) {
	o := roles.Occurrence{Event: inbox.Edited, Item: it, Edit: roles.Edit{Changed: n.Changed, By: n.By, At: n.At, Own: storyAgent(agent, n.ID)}}
	o.Draft = it != nil && it.Type == workitem.Story && it.Draft && (n.By != replanBy || !slices.Contains(n.Changed, "forecast"))
	for _, t := range matching(planner.Trigger, o) {
		switch {
		case t.When == "stale-figures":
			if trigger := editTrigger(it, n, agent); trigger != "" {
				r.queue = enqueue(r.queue, n.ID, trigger)
			}
		case it != nil && takes(planner, n.ID) && !roleAgentOf(planner, n.By) && n.By != replanBy:
			r.queue = enqueue(r.queue, n.ID, "edited "+strings.Join(n.Changed, ", ")+" by "+orSomeone(n.By))
		}
	}
}

// queueAhead queues the planner, whose role file is planner, for a, work
// ahead completing or the order changing, under each entry of the file for
// a's event that a matches: under replan-agent for each story the replay
// wrote, written, while policy, planning.replan, is agent; under no
// condition for a's own item, when it has one of a kind the planner takes.
// Each is queued with a's trigger.
func (r *replanner) queueAhead(planner roles.Role, a ahead, policy string, written []string) {
	for _, t := range matching(planner.Trigger, roles.Occurrence{Event: a.event, To: a.to, Replan: policy}) {
		var ids []string
		switch {
		case t.When == "replan-agent":
			ids = roles.ReplanAgent(policy, written)
		case t.When == "" && a.item != "" && takes(planner, a.item):
			ids = []string{a.item}
		}
		for _, id := range ids {
			r.queue = enqueue(r.queue, id, a.trigger)
		}
	}
}

// schedules are the settings the planner's schedule entries name, by
// setting, parsed from m, each seen by its cadence at now, so that one set or
// changed comes round first at its next time from now, and one newly named
// is followed from now. A setting that is unset, that the planner may not
// name, or that cannot be parsed is left out, and the latter two are warned
// of once. Cadences of settings the file no longer names are let go.
func (r *replanner) schedules(planner roles.Role, m manifest.Manifest, now time.Time) map[string]*cron.Schedule {
	scheds, cadences := followSchedules(planner, r.cadences, m, now, func(err error) { r.warnOnce("planning schedule not valid", err) })
	r.cadences = cadences
	return scheds
}

// beforeWrite runs between the replay and the write of the forecasts it
// moved, so that a test can bring another writer's save in between.
var beforeWrite = func() {}

// replan plays the board out again after triggers as planning.replan says:
// it writes the forecast deliveries that moved as at now, stamped at the
// write (I-0183), in one pass under the items lock with one check over the
// written items and one index write (itemedit.ApplyForecasts, S-0599), and
// commits them in one commit unless autocommit is off.
// It returns the policy and the stories it wrote, which replan-agent's
// entries queue the planner for under agent. A policy that is not valid is
// warned of once, and nothing is done.
func (r *replanner) replan(repo *workitem.Repo, items []*workitem.Item, board *workitem.Board, now time.Time, triggers []string) (policy string, written []string) {
	policy, err := repo.Manifest.Planning.ReplanPolicy()
	if err != nil {
		r.warnOnce("planning replan not valid", err)
		return "", nil
	}
	if policy == manifest.ReplanNever {
		return policy, nil
	}
	moves, errs := replays(items, board, repo.Manifest.Planning, now)
	for _, err := range errs {
		r.starter.warn("forecast not replayed", "err", err)
	}
	var files []string
	if len(moves) > 0 {
		changes := make([]itemedit.ForecastChange, 0, len(moves))
		for _, m := range moves {
			delivery, basis := m.delivery, m.basis
			changes = append(changes, itemedit.ForecastChange{ID: m.id, Edit: itemedit.ForecastEdit{Delivery: &delivery, Basis: &basis}})
		}
		beforeWrite()
		// the replay's own items, so that a save merges against what the
		// replay read: a forecast another writer changed since is refused
		batch, err := itemedit.ApplyForecasts(repo, items, changes, itemedit.Options{By: "flai", NoCommit: true, Now: r.o.Now().UTC()})
		if err != nil {
			r.starter.warn("forecasts not replanned", "moves", len(moves), "err", err)
		} else {
			for _, m := range moves {
				if err, refused := batch.Refused[workitem.CanonicalID(m.id)]; refused {
					r.starter.warn("forecast not replanned", "item", m.id, "err", err)
				} else if slices.Contains(batch.Written, workitem.CanonicalID(m.id)) {
					written = append(written, m.id)
				}
			}
			files = batch.Files
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
	return policy, written
}

// drain starts the planner for the first queued story with no planner
// running, unless a planner this replanner started still runs. An entry for
// a story its agent works, by byID and the story's newest run, or one the
// checks refuse, is dropped and logged, and the next one tried; one that
// cannot be read is dropped and warned of. One whose start the budget state
// of the planner's provider holds keeps its place, as does every entry after
// it, and the next look tries it again (S-0534); it is logged once per state.
func (r *replanner) drain(ctx context.Context, byID map[string]*workitem.Item) {
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
		if why := worked(q.item, byID[q.item], st.Stories[q.item]); why != "" {
			r.starter.log("queued planner dropped", "item", q.item, "trigger", trigger, "reason", why)
			continue
		}
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
		_, held := r.starter.plan(ctx, ps.cfg, ps.item, ps.agent, trigger)
		r.starter.mu.Unlock()
		if held != nil {
			r.queue = slices.Insert(r.queue, i, q)
			r.starter.heldOnce(r.said, "queued planner held", *held, "item", q.item, "trigger", trigger)
			return
		}
		r.started = ps.item
		return
	}
}

// worked is why the planner is not started on its own for story id, it as
// read, while its agent works it: it is in progress or in review, or run,
// its newest agent run, runs; "" when neither (I-0098).
func worked(id string, it *workitem.Item, run *AgentRun) string {
	if it != nil && (it.Status == workitem.InProgress || it.Status == workitem.Review) {
		return fmt.Sprintf("%s is %s; its agent works it, so the planner is not started for it on its own", id, it.Status)
	}
	if run.running() {
		return fmt.Sprintf("the agent is running for %s (pid %d, started %s); the planner is not started for it on its own", id, run.PID, run.Started)
	}
	return ""
}

// storyAgent is the name flai serve gives the run of story's agent, its
// FLAI_AGENT: the configured name, agent when there is none, and the
// story's ID.
func storyAgent(name, story string) string {
	if name == "" {
		name = "agent"
	}
	return name + "-" + story
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

// editTrigger is the trigger an edit gives story it under stale-figures:
// "edited", the fields that make it due, and who edited; "" when it gives
// none (ADR-0084, roles.StaleFigures). It gives none for an item that is not
// a backlog or ready story on the board, for a planner's edit, and for an
// edit by the story's own agent, whose configured name is agent (I-0098). A
// goal, criteria, or touches edit makes it due when its forecast or its cost
// of delay value is older than the edit; a cost of delay edit, when its
// inputs are now newer than its value. A stamp that cannot be read is not
// older.
func editTrigger(it *workitem.Item, n itemedit.Notice, agent string) string {
	if it == nil {
		return ""
	}
	fields := roles.StaleFigures(it, roles.Edit{Changed: n.Changed, By: n.By, At: n.At, Own: storyAgent(agent, it.ID)})
	if len(fields) == 0 {
		return ""
	}
	return "edited " + strings.Join(fields, ", ") + " by " + orSomeone(n.By)
}

// moveTriggers are the triggers among changes of work ahead completing, in
// order: accepted for a story moved to done, cancelled for one moved to
// cancelled. A task's or an epic's move is none.
func moveTriggers(changes []workitem.Change) []string {
	return triggersOf(aheadOf(changes, false))
}

// ahead is work ahead completing or the order changing: the trigger event it
// is, accepted, cancelled, or reordered; the state a move went to; the item
// it is on, "" for a reorder; and its trigger as a run records it,
// "accepted S-0001" or "reordered".
type ahead struct {
	event, to, item, trigger string
}

// aheadOf is the work ahead completing among changes, in order, a story
// moved to done or to cancelled, then the order changing when reordered
// says it did. A task's or an epic's move is none.
func aheadOf(changes []workitem.Change, reordered bool) []ahead {
	var out []ahead
	for _, c := range changes {
		if c.Type != workitem.Story || c.Kind != workitem.Moved {
			continue
		}
		switch c.To {
		case workitem.Done:
			out = append(out, ahead{event: "accepted", to: c.To, item: c.ID, trigger: "accepted " + c.ID})
		case workitem.Cancelled:
			out = append(out, ahead{event: "cancelled", to: c.To, item: c.ID, trigger: "cancelled " + c.ID})
		}
	}
	if reordered {
		out = append(out, ahead{event: triggerReordered, trigger: triggerReordered})
	}
	return out
}

// triggersOf are the triggers of as, in order.
func triggersOf(as []ahead) []string {
	out := make([]string, len(as))
	for i, a := range as {
		out[i] = a.trigger
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
