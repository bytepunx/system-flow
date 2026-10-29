package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A role a project defines, started by flai serve for the events its
// trigger names (S-0418, ADR-0207 §3, ADR-0211).
//
// At each look, for every role a file under the design folder's
// workflows/roles defines, other than a built-in one, whose trigger names an
// event and whose host action is on, flai serve reads the inbox's events
// since the role's last look as the role, under a cursor of its own
// (roleCursor), and queues those its triggers match as reasons, by the
// orchestrator's matching (reasonOf). A thread event is the role's when the
// thread awaits one of its agents: the role's name for a role that takes no
// item, and the role's name and the thread's item (harness.RoleAgent) for a
// role that takes items, which is how its runs write. What the role's own
// agents changed is no reason for it.
//
// A role that takes no item has one queue; a role that takes items has one
// for each item a reason is on, its ID, or, for a thread, the item the
// thread is on, when that is of a kind the role takes; a reason on no such
// item is not queued. A queue is due once orchestration.coalesce has passed
// since its first reason was queued, and with no run of the role on its item
// going, its reasons start one run whose trigger names them all
// (roleTrigger); those that come while it runs wait for the next. The queues
// are kept per project and per role in serve/role-reasons.json, so that a
// flai serve started again loses none. A start the role's checks refuse
// drops its reasons and is logged; one the budget state holds keeps them for
// the next look.
//
// While a role's host action is off its queue is emptied and nothing is
// read; the look that finds it turned on takes the events since as history,
// so that turning it on does not act on the past. A role file that does not
// load is skipped, and flai check reports it.

// roleCursorPrefix starts the name of the inbox cursor flai serve reads a
// defined role's events under, so that no agent's own cursor moves.
const roleCursorPrefix = "flai-serve-role-"

// roleCursor is the inbox cursor flai serve reads the role name's events
// under.
func roleCursor(name string) string { return roleCursorPrefix + name }

// roleQueue is what flai serve keeps of one role a project defines between
// looks, and between one flai serve and the next.
type roleQueue struct {
	// On is whether the role's host action was on at the last look; a look
	// that finds it on when this is false finds it turned on.
	On bool `json:"on,omitempty"`
	// Told are the threads that awaited one of the role's agents at the last
	// look, each by its ID and its newest entry's time (threadKey), so that
	// each such entry is a reason once.
	Told []string `json:"told,omitempty"`
	// Items are the reasons that wait for a run, by the item the run is on,
	// "" for a role that takes no item.
	Items map[string]reasonQueue `json:"items,omitempty"`
}

// empty says whether q keeps nothing worth writing.
func (q roleQueue) empty() bool { return !q.On && len(q.Told) == 0 && len(q.Items) == 0 }

// roleEvents starts the roles one project defines for the events their
// triggers name (S-0418).
type roleEvents struct {
	o       Options
	e       Entry
	starter *launcher

	mu sync.Mutex
	// said are the failures warned of, so that a look every minute does not
	// warn of them again.
	said map[string]bool
	// held are the holds logged (heldOnce), by role name and item, so that
	// each run held is logged once per budget state.
	held map[string]map[string]bool
	// wake is sent on when the reasons queued come due between looks, and
	// timer sends it.
	wake  chan struct{}
	timer *time.Timer
}

// newRoleEvents is the role starter for project e, starting the runs with
// starter, the project's launcher.
func newRoleEvents(o Options, e Entry, starter *launcher) *roleEvents {
	return &roleEvents{o: o, e: e, starter: starter, said: map[string]bool{}, held: map[string]map[string]bool{}, wake: make(chan struct{}, 1)}
}

// look queues the reasons the events since the last look give each role
// whose action is on, empties the queue of each whose action is off, and
// starts a run for each queue that is due with no run of its role on its
// item going.
func (r *roleEvents) look(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg AgentConfig
	if r.o.Agent != nil {
		cfg = r.o.Agent(r.e.Root)
	}
	now := r.o.Now().UTC()
	repo, err := workitem.Open(r.e.Root)
	if err != nil {
		r.warnOnce("role events could not read the project", err)
		return
	}
	all, problems, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		r.warnOnce("role files not read", err)
		return
	}
	evented := eventRoles(all)
	qs := r.gather(repo, cfg, evented, problems, now)
	coalesce, err := repo.Manifest.Orchestration.CoalesceDuration()
	if err != nil {
		r.warnOnce("orchestration coalesce not valid", err)
		coalesce = manifest.DefaultCoalesce
	}
	r.startDue(ctx, cfg, evented, qs, coalesce, now)
}

// eventRoles are the roles among all that flai serve starts for events: not
// built-in, with a trigger that names an event.
func eventRoles(all []roles.Role) []roles.Role {
	var out []roles.Role
	for _, f := range all {
		if _, builtin := builtinRoles[f.Name]; builtin {
			continue
		}
		if slices.ContainsFunc(f.Trigger, func(t roles.Trigger) bool { return t.Kind == roles.TriggerEvent }) {
			out = append(out, f)
		}
	}
	return out
}

// gather reads the events since the last look of each role among evented
// whose action is on, queues those its triggers match, empties the queue of
// each whose action is off, lets go of the queues of roles no file defines
// any more, and returns the project's queues as written. A role found
// turned on takes the events as history and queues nothing. When the
// project's items or threads cannot be read, nothing is read or queued, and
// the queues are returned as kept.
func (r *roleEvents) gather(repo *workitem.Repo, cfg AgentConfig, evented []roles.Role, problems []roles.Problem, now time.Time) map[string]roleQueue {
	var on []roles.Role
	for _, f := range evented {
		if cfg.RoleActions[f.Action] {
			on = append(on, f)
		}
	}
	var items []*workitem.Item
	var ths []*threads.Thread
	if len(on) > 0 {
		var err error
		if items, err = repo.List(true); err != nil {
			r.warnOnce("role events could not read the project", err)
			return r.o.Dir.roleQueues()[r.e.Root]
		}
		if ths, err = threads.List(repo); err != nil {
			r.warnOnce("role events could not read the threads", err)
			return r.o.Dir.roleQueues()[r.e.Root]
		}
	}
	read := map[string][]inbox.Event{}
	awaiting := map[string][]inbox.Event{}
	for _, f := range on {
		read[f.Name] = r.catchUp(repo, f, items, now)
		awaiting[f.Name] = awaitingRole(repo, f, ths, now)
	}
	byID := itemsByID(items)
	queued := map[string]int{}
	qs, err := r.o.Dir.updateRoleQueues(r.e.Root, func(qs map[string]roleQueue) {
		for name := range qs {
			if !slices.ContainsFunc(evented, func(f roles.Role) bool { return f.Name == name }) && !unloaded(problems, name) {
				delete(qs, name)
			}
		}
		for _, f := range evented {
			if !cfg.RoleActions[f.Action] {
				delete(qs, f.Name)
				continue
			}
			q := qs[f.Name]
			events := slices.Clone(read[f.Name])
			told := make([]string, 0, len(awaiting[f.Name]))
			for _, e := range awaiting[f.Name] {
				if key := threadKey(e); !slices.Contains(q.Told, key) {
					events = append(events, e)
				}
				told = append(told, threadKey(e))
			}
			q.Told = told
			if !q.On {
				q.On = true // turned on: what came before is history
				qs[f.Name] = q
				continue
			}
			slices.SortStableFunc(events, func(a, b inbox.Event) int { return strings.Compare(a.At, b.At) })
			queued[f.Name] = q.queue(f, events, byID, now)
			qs[f.Name] = q
		}
	})
	if err != nil {
		r.starter.warn("role reasons not kept", "err", err)
	}
	for _, f := range on {
		if n := queued[f.Name]; n > 0 {
			waiting := 0
			for _, iq := range qs[f.Name].Items {
				waiting += len(iq.Reasons)
			}
			r.starter.log("role reasons queued", "role", f.Name, "queued", n, "waiting", waiting)
		}
	}
	return qs
}

// queue queues, at now, each of events that f's triggers match as a reason
// on the item f's run would take for it, and returns how many it queued.
// What f's own agents changed is no reason, nor is an event on no item f
// takes. byID are the project's items.
func (q *roleQueue) queue(f roles.Role, events []inbox.Event, byID map[string]*workitem.Item, now time.Time) int {
	n := 0
	for _, e := range events {
		if roleAgentOf(f, e.By) {
			continue
		}
		item, ok := reasonItem(f, e)
		if !ok {
			continue
		}
		reason, ok := reasonOf(e, f.Trigger, byID, func(awaits string) bool { return roleAgentOf(f, awaits) })
		if !ok {
			continue
		}
		if q.Items == nil {
			q.Items = map[string]reasonQueue{}
		}
		iq := q.Items[item]
		n += iq.add(now, reason)
		q.Items[item] = iq
	}
	return n
}

// catchUp is the events since f's cursor's last look as f reads them, its
// own changes left out, but for threads, which awaitingRole reads for each
// of f's agents. A failure is warned of once, and is no events.
func (r *roleEvents) catchUp(repo *workitem.Repo, f roles.Role, items []*workitem.Item, now time.Time) []inbox.Event {
	events, _, err := inbox.CatchUp(repo, inbox.Reader{Agent: harness.RoleAgent(f.Name, ""), Cursor: roleCursor(f.Name)}, "", false, now, items)
	if err != nil {
		r.warnOnce("role events could not read the project's events", err, "role", f.Name)
		return nil
	}
	return slices.DeleteFunc(events, func(e inbox.Event) bool { return e.Kind == inbox.Thread })
}

// awaitingRole are the open threads among ths whose newest entry, someone
// else's and not after now, makes them await one of f's agents under the
// inbox's rule (inbox.Summarize, ADR-0181), as thread events: the role's
// name for a role that takes no item, and its run's agent on the thread's
// item for one that takes items, when the thread is on an item of a kind it
// takes. Fields name the agent awaited and the item or path the thread is
// on.
func awaitingRole(repo *workitem.Repo, f roles.Role, ths []*threads.Thread, now time.Time) []inbox.Event {
	var out []inbox.Event
	for _, th := range ths {
		if !th.Open() {
			continue
		}
		on := ""
		if th.Anchor.Item != "" {
			on = workitem.CanonicalID(th.Anchor.Item)
		}
		agent := harness.RoleAgent(f.Name, "")
		if len(f.Item) > 0 {
			if !slices.Contains(f.Item, workitem.TypeOfID(on)) {
				continue
			}
			agent = harness.RoleAgent(f.Name, on)
		}
		sum := inbox.Summarize(repo, inbox.Reader{Agent: agent}, th)
		if sum.Awaiting != "you" || sum.LastBy == agent {
			continue
		}
		entries := th.Entries() // at least one: a thread with none awaits another
		at := entries[len(entries)-1].At
		if t, err := time.Parse(workitem.TimeFormat, at); err == nil && t.After(now) {
			continue
		}
		if on == "" {
			on = th.Anchor.Path
		}
		cause := sum.Story
		if cause == "" {
			cause = th.Anchor.Path
		}
		c := workitem.Change{ID: th.ID, Type: inbox.Thread, Title: th.Title, Kind: inbox.Thread, To: "you", By: sum.LastBy, Cause: cause, At: at}
		out = append(out, inbox.Event{Change: c, Summary: inbox.Describe(c), Fields: map[string]string{inbox.FieldAwaits: agent, inbox.FieldOn: on}})
	}
	return out
}

// threadKey is what a thread event is told by: the thread's ID and its
// newest entry's time.
func threadKey(e inbox.Event) string { return e.ID + "|" + e.At }

// roleAgentOf says whether name is one of f's agents: its name, or, for a
// role that takes items, its run's agent on an item of a kind it takes.
func roleAgentOf(f roles.Role, name string) bool {
	if name == harness.RoleAgent(f.Name, "") {
		return true
	}
	item, ok := strings.CutPrefix(name, f.Name+"-")
	return ok && len(f.Item) > 0 && slices.Contains(f.Item, workitem.TypeOfID(item)) && harness.RoleAgent(f.Name, item) == name
}

// reasonItem is the item f's run for e is on, "" for a role that takes no
// item: e's own ID, or the item a thread is on, when it is of a kind f
// takes. ok is false when it is none, and e is no reason for f.
func reasonItem(f roles.Role, e inbox.Event) (item string, ok bool) {
	if len(f.Item) == 0 {
		return "", true
	}
	id := e.ID
	if e.Kind == inbox.Thread {
		id = e.Fields[inbox.FieldOn]
	}
	id = workitem.CanonicalID(id)
	return id, slices.Contains(f.Item, workitem.TypeOfID(id))
}

// unloaded says whether the role name's file is among those that did not
// load, whose queue is kept until it loads again.
func unloaded(problems []roles.Problem, name string) bool {
	return slices.ContainsFunc(problems, func(p roles.Problem) bool { return strings.TrimSuffix(path.Base(p.Path), ".md") == name })
}

// roleTrigger is what started a run for rs, as its run records it: each
// reason's kind and ID, "moved S-0001", once each, joined by semicolons.
func roleTrigger(rs []orchestrateReason) string {
	var said []string
	for _, r := range rs {
		s := r.Kind
		if r.ID != "" {
			s += " " + r.ID
		}
		if !slices.Contains(said, s) {
			said = append(said, s)
		}
	}
	return strings.Join(said, "; ")
}

// startDue starts a run for each queue in qs of a role among evented whose
// action is on that is due at now, by coalesce, with no run of its role on
// its item going, and has the serving loop look again when the first of the
// others comes due.
func (r *roleEvents) startDue(ctx context.Context, cfg AgentConfig, evented []roles.Role, qs map[string]roleQueue, coalesce time.Duration, now time.Time) {
	st := r.o.Dir.AgentStates()[r.e.Root]
	var next time.Time
	for _, f := range evented {
		if !cfg.RoleActions[f.Action] {
			continue
		}
		q := qs[f.Name]
		for _, item := range slices.Sorted(maps.Keys(q.Items)) {
			iq := q.Items[item]
			ready, at := iq.due(now, coalesce)
			if !ready {
				if !at.IsZero() && (next.IsZero() || at.Before(next)) {
					next = at
				}
				continue
			}
			if st.Roles[f.Name][item].running() {
				continue // its end has the loop look again
			}
			r.start(ctx, f.Name, item, iq.Reasons)
		}
	}
	if !next.IsZero() {
		r.wakeIn(next.Sub(now))
	}
}

// start starts the role name on item for rs, the reasons queued for it,
// holding the launcher's lock. A start made, or one that failed and was
// recorded so, takes rs out of the queue; so does one the role's checks
// refuse, which is logged. One the budget state holds keeps them for the
// next look, logged once per state, and any other error keeps them and is
// warned of once.
func (r *roleEvents) start(ctx context.Context, name, item string, rs []orchestrateReason) {
	trigger := roleTrigger(rs)
	r.starter.mu.Lock()
	_, held, err := startRole(ctx, r.o, r.e, r.starter, name, item, trigger)
	r.starter.mu.Unlock()
	var no *Refused
	switch {
	case errors.As(err, &no):
		r.starter.log("role start refused", "role", name, "item", item, "trigger", trigger, "reason", no.Why)
	case err != nil:
		r.warnOnce("role not started", err, "role", name, "item", item, "trigger", trigger)
		return
	case held != nil:
		key := name + "\x00" + item
		if r.held[key] == nil {
			r.held[key] = map[string]bool{}
		}
		r.starter.heldOnce(r.held[key], "role start held", *held, "role", name, "item", item, "trigger", trigger)
		return
	}
	if _, err := r.o.Dir.updateRoleQueues(r.e.Root, func(qs map[string]roleQueue) {
		q := qs[name]
		iq, ok := q.Items[item]
		if !ok {
			return
		}
		iq.remove(rs)
		q.Items[item] = iq
		qs[name] = q
	}); err != nil {
		r.starter.warn("role reasons not kept", "role", name, "err", err)
	}
}

// wakeIn has the serving loop look again in d, when the reasons queued come
// due, in place of any earlier wake.
func (r *roleEvents) wakeIn(d time.Duration) {
	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(max(d, 0), func() {
		select {
		case r.wake <- struct{}{}:
		default:
		}
	})
}

// warnOnce warns of err with msg and args the first time it is seen.
func (r *roleEvents) warnOnce(msg string, err error, args ...any) {
	if key := msg + "\x00" + fmt.Sprint(args...) + "\x00" + err.Error(); !r.said[key] {
		r.said[key] = true
		r.starter.warn(msg, append(args, "err", err)...)
	}
}

// roleQueuesFile is where the queues of the roles every project defines are
// kept, by project root and role name, beside orchestrate-reasons.json, so
// that a flai serve started again loses none.
func (d Dir) roleQueuesFile() string { return filepath.Join(string(d), "role-reasons.json") }

var roleQueuesMu sync.Mutex

// roleQueues reads the queues of the roles every project defines, by project
// root and role name; none when the file cannot be read.
func (d Dir) roleQueues() map[string]map[string]roleQueue {
	out := map[string]map[string]roleQueue{}
	if data, err := os.ReadFile(d.roleQueuesFile()); err == nil {
		_ = json.Unmarshal(data, &out)
	}
	return out
}

// updateRoleQueues changes the queues of the roles the project at root
// defines with change, drops an item's queue with no reasons and a role's
// that keeps nothing, writes them when they changed, and returns them as
// written.
func (d Dir) updateRoleQueues(root string, change func(map[string]roleQueue)) (map[string]roleQueue, error) {
	roleQueuesMu.Lock()
	defer roleQueuesMu.Unlock()
	defer lockFile(d.roleQueuesFile() + ".lock")()
	all := d.roleQueues()
	before, _ := json.Marshal(all[root])
	qs := all[root]
	if qs == nil {
		qs = map[string]roleQueue{}
	}
	change(qs)
	for name, q := range qs {
		maps.DeleteFunc(q.Items, func(_ string, iq reasonQueue) bool { return len(iq.Reasons) == 0 })
		if len(q.Items) == 0 {
			q.Items = nil
		}
		if q.empty() {
			delete(qs, name)
		} else {
			qs[name] = q
		}
	}
	if len(qs) == 0 {
		delete(all, root)
	} else {
		all[root] = qs
	}
	if after, _ := json.Marshal(all[root]); string(after) == string(before) {
		return qs, nil
	}
	if err := d.write(d.roleQueuesFile(), all); err != nil {
		return qs, fmt.Errorf("write the roles' reasons for %s to %s: %w", root, d.roleQueuesFile(), err)
	}
	return qs, nil
}
