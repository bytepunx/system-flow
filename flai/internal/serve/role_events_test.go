//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// watcherRole is a role a project defines that runs on no item, behind the
// watch host action, started for a story moved to review, a draft's edit, a
// planner's run that wrote items, and a thread that awaits it.
const watcherRole = `---
title: Watcher
updated: 2026-10-29
status: active
name: watcher
role: watch
definition:
  description: Watches the project, for flai serve.
  tools: [Read, Grep]
  model: inherit
action: watch
trigger:
  - asked
  - {event: moved, when: to-review}
  - {event: edited, when: draft-changed}
  - {event: planned}
  - {event: thread}
item: none
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the watcher: read what changed and say what it means.
`

// checkerRole is a role a project defines that runs on a story, behind the
// check host action, started for a story moved to review and a thread on a
// story that awaits its run on it.
const checkerRole = `---
title: Checker
updated: 2026-10-29
status: active
name: checker
role: check
definition:
  description: Checks a story, for flai serve.
  tools: [Read, Grep]
  model: inherit
action: check
trigger:
  - {event: moved, when: to-review}
  - {event: thread}
item: [story]
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the checker: check the story and say what it lacks.
`

// roleEventsLab is a roleLab whose project defines watcher and checker too,
// with every role's action on, and whose role starter reads the clock at.
type roleEventsLab struct {
	*agentLab
	ev *roleEvents
	at time.Time
}

// newRoleEventsLab is a roleEventsLab whose items made from here on are an
// hour old, history to the first look.
func newRoleEventsLab(t *testing.T) *roleEventsLab {
	t.Helper()
	lab := &roleEventsLab{agentLab: roleLab(t)}
	dir := filepath.Join(lab.root, "design", "workflows", "roles")
	for name, body := range map[string]string{"watcher.md": watcherRole, "checker.md": checkerRole} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lab.cfg.RoleActions = map[string]bool{"review": true, "audit": true, "watch": true, "check": true}
	o := lab.o
	o.Now = func() time.Time { return lab.at }
	lab.ev = newRoleEvents(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	lab.now = time.Now().Add(-time.Hour)
	return lab
}

// lookAt has the role starter look at at.
func (lab *roleEventsLab) lookAt(at time.Time) {
	lab.at = at
	lab.ev.look(context.Background())
}

// queued is the role name's queue as kept in the serve folder.
func (lab *roleEventsLab) queued(name string) roleQueue {
	return lab.o.Dir.roleQueues()[lab.root][name]
}

// moveAt moves item id to state to at, as its agent would.
func (lab *roleEventsLab) moveAt(id, to string, at time.Time) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, to, "builder", "", at); err != nil {
		lab.t.Fatal(err)
	}
}

// working makes a story with a task, in progress, as history.
func (lab *roleEventsLab) working(title string) string {
	lab.t.Helper()
	story := lab.backlog(title, nil)
	lab.toReady(story)
	lab.moveAt(story, workitem.InProgress, lab.now)
	if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: story, Owner: "builder-" + story, Now: lab.now}); err != nil {
		lab.t.Fatal(err)
	}
	return story
}

// asks opens a thread on item by author, as history, which awaits the
// operator, and returns its ID.
func (lab *roleEventsLab) asks(item, author string) string {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which way", On: item, Author: author, Text: "Which way?", Now: lab.now})
	if err != nil {
		lab.t.Fatal(err)
	}
	return th.ID
}

// answer answers thread id as the operator at, so that it awaits its
// opener.
func (lab *roleEventsLab) answer(id string, at time.Time) {
	lab.t.Helper()
	if _, err := threads.Reply(lab.repo, id, "alex", "This way.", at); err != nil {
		lab.t.Fatal(err)
	}
}

// ended waits until the role name's run on item has ended.
func (lab *roleEventsLab) ended(name, item string) {
	lab.t.Helper()
	waitFor(lab.t, "the run of "+name+" "+onItem(item)+" ended", func() bool {
		r := lab.roleRun(name, item)
		return r != nil && r.Ended != ""
	})
}

// logged counts the lines logged with msg for the role name.
func (lab *roleEventsLab) logged(msg, name string) int {
	lab.logs.mu.Lock()
	defer lab.logs.mu.Unlock()
	n := 0
	for line := range strings.Lines(lab.logs.buf.String()) {
		if strings.Contains(line, `msg="`+msg+`"`) && strings.Contains(line, "role="+name+" ") {
			n++
		}
	}
	return n
}

// S-0418 criteria 1 and 3: a role that takes no item is queued for each
// event its trigger names, a move to review, a planner's run, a draft's
// edit, and a thread that awaits it, coalesced into one run once the window
// has passed since the first was queued and not before, whose trigger names
// each, joined by semicolons; the start takes them out of the queue.
func TestARoleIsQueuedForTheEventsItsTriggerNames(t *testing.T) {
	lab := newRoleEventsLab(t)
	story := lab.working("Reviewed")
	draft, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Drafted", Parent: lab.epic.ID, Owner: "alex", Draft: true, Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	thread := lab.asks(story, "watcher")
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	if q := lab.queued("watcher"); !q.On || len(q.Items) != 0 {
		t.Fatalf("the first look queued %+v, want the role on with nothing queued", q)
	}
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	planned := inbox.PlannedEvent(lab.epic.ID, workitem.Epic, "Epic", "planner", []string{draft.ID})
	planned.At = base.Add(2 * time.Second).Format(workitem.TimeFormat)
	if err := inbox.AppendEvent(lab.repo, planned); err != nil {
		t.Fatal(err)
	}
	touches := []string{"docs/drafted"}
	if _, err := itemedit.Apply(lab.repo, gitOf(lab.o), draft.ID, itemedit.Change{Touches: &touches}, itemedit.Options{By: "planner", NoCommit: true, Now: base.Add(3 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	lab.answer(thread, base.Add(4*time.Second))
	queued := base.Add(5 * time.Second)
	for _, at := range []time.Time{queued, queued.Add(manifest.DefaultCoalesce - time.Second)} {
		lab.lookAt(at)
		if r := lab.roleRun("watcher", ""); r != nil {
			t.Fatalf("started at %s, inside the window from %s: %+v", at, queued, r)
		}
	}
	if q := lab.queued("watcher").Items[""]; len(q.Reasons) != 4 || q.First != queued.Format(time.RFC3339) {
		t.Fatalf("queue = %+v, want four reasons first queued at %s", q, queued)
	}
	lab.lookAt(queued.Add(manifest.DefaultCoalesce))
	run := lab.roleRun("watcher", "")
	if run == nil {
		t.Fatal("not started once the window passed")
	}
	want := "moved " + story + "; planned " + lab.epic.ID + "; edited " + draft.ID + "; thread " + thread
	if run.Trigger != want || run.Agent != "watcher" {
		t.Errorf("run's trigger %q as %s, want %q as watcher", run.Trigger, run.Agent, want)
	}
	if q := lab.queued("watcher"); len(q.Items) != 0 {
		t.Errorf("the start left %+v queued", q.Items)
	}
	lab.ended("watcher", "")
	lab.lookAt(queued.Add(time.Hour))
	if r := lab.roleRun("watcher", ""); r.Session != run.Session {
		t.Errorf("started again with nothing queued: %+v", r)
	}
}

// S-0418: an event the trigger names under a condition that does not hold,
// a move elsewhere than review or an edit of a story that is no draft, and
// a thread that awaits the operator rather than the role, queue nothing; nor
// does what the role's own agent changed.
func TestAConditionThatDoesNotHoldQueuesNothing(t *testing.T) {
	lab := newRoleEventsLab(t)
	story := lab.working("Worked")
	other := lab.backlog("Other", nil)
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	lab.moveAt(other, workitem.Ready, base.Add(time.Second))
	touches := []string{"docs/other"}
	if _, err := itemedit.Apply(lab.repo, gitOf(lab.o), other, itemedit.Change{Touches: &touches}, itemedit.Options{By: "planner", NoCommit: true, Now: base.Add(time.Second)}); err != nil {
		t.Fatal(err)
	}
	if _, err := threads.New(lab.repo, threads.NewOptions{Title: "Ask", On: story, Author: "agent-" + story, Text: "Which way?", Now: base.Add(time.Second)}); err != nil {
		t.Fatal(err)
	}
	it, err := lab.repo.Get(story)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, workitem.Review, "watcher", "", base.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	lab.lookAt(base.Add(2 * time.Second))
	if q := lab.queued("watcher"); len(q.Items) != 0 {
		t.Fatalf("queued %+v, want nothing", q.Items)
	}
	lab.lookAt(base.Add(time.Hour))
	if r := lab.roleRun("watcher", ""); r != nil {
		t.Errorf("started for nothing: %+v", r)
	}
}

// S-0418 criterion 1: a role that takes items gets one run per item, on the
// story moved and on the story a thread that awaits its run there is on;
// a reason that comes while the item's run goes waits for it to end, and
// starts the next.
func TestARoleThatTakesItemsRunsOncePerItem(t *testing.T) {
	lab := newRoleEventsLab(t)
	reviewed, asked := lab.working("Reviewed"), lab.working("Asked")
	thread := lab.asks(asked, "checker-"+asked)
	lab.hold()
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	lab.moveAt(reviewed, workitem.Review, base.Add(time.Second))
	lab.answer(thread, base.Add(time.Second))
	lab.lookAt(base.Add(2 * time.Second))
	q := lab.queued("checker")
	if len(q.Items) != 2 || len(q.Items[reviewed].Reasons) != 1 || len(q.Items[asked].Reasons) != 1 {
		t.Fatalf("queue = %+v, want one reason on each story", q.Items)
	}
	started := base.Add(2*time.Second + manifest.DefaultCoalesce)
	lab.lookAt(started)
	first := lab.roleRun("checker", reviewed)
	if first == nil || first.Trigger != "moved "+reviewed || first.Agent != "checker-"+reviewed {
		t.Fatalf("the run on %s: %+v", reviewed, first)
	}
	if r := lab.roleRun("checker", asked); r == nil || r.Trigger != "thread "+thread || r.Agent != "checker-"+asked {
		t.Fatalf("the run on %s: %+v", asked, r)
	}
	it, err := lab.repo.Get(reviewed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, workitem.InProgress, "builder", "sent back", started.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	lab.moveAt(reviewed, workitem.Review, started.Add(2*time.Second))
	again := started.Add(3*time.Second + manifest.DefaultCoalesce)
	lab.lookAt(started.Add(3 * time.Second))
	lab.lookAt(again)
	if r := lab.roleRun("checker", reviewed); r.Session != first.Session {
		t.Fatalf("started while the run on %s goes: %+v", reviewed, r)
	}
	if q := lab.queued("checker").Items[reviewed]; len(q.Reasons) != 1 {
		t.Fatalf("queue on %s = %+v, want the move waiting", reviewed, q)
	}
	lab.release("checker-" + reviewed)
	lab.ended("checker", reviewed)
	lab.lookAt(again)
	if r := lab.roleRun("checker", reviewed); r.Session == first.Session || r.Trigger != "moved "+reviewed {
		t.Errorf("once the run ended: %+v, want the next for the move", r)
	}
	lab.release("checker-" + asked)
}

// S-0418 criterion 1: a start the role's checks refuse drops its reasons and
// logs why, once.
func TestARefusedStartDropsItsReasons(t *testing.T) {
	lab := newRoleEventsLab(t)
	story := lab.working("Reviewed")
	lab.cfg.Command = nil // no harness and no command: role.run refuses
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	lab.lookAt(base.Add(2 * time.Second))
	lab.lookAt(base.Add(2*time.Second + manifest.DefaultCoalesce))
	if r := lab.roleRun("watcher", ""); r != nil {
		t.Fatalf("started: %+v", r)
	}
	if q := lab.queued("watcher"); len(q.Items) != 0 {
		t.Errorf("the refusal kept %+v", q.Items)
	}
	lab.lookAt(base.Add(time.Hour))
	if n := lab.logged("role start refused", "watcher"); n != 1 {
		t.Errorf("the refusal was logged %d times, want once", n)
	}
}

// S-0418, S-0534: a start the budget state holds keeps its reasons queued,
// logged once, and the first look once it no longer holds starts the run.
func TestAHeldStartKeepsItsReasons(t *testing.T) {
	lab := newRoleEventsLab(t)
	story := lab.working("Reviewed")
	lab.budgetHeldOn(harness.Command, anHour())
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	lab.lookAt(base.Add(2 * time.Second))
	due := base.Add(2*time.Second + manifest.DefaultCoalesce)
	for range 2 {
		lab.lookAt(due)
	}
	if r := lab.roleRun("watcher", ""); r != nil {
		t.Fatalf("started while held: %+v", r)
	}
	if q := lab.queued("watcher").Items[""]; len(q.Reasons) != 1 {
		t.Fatalf("queue = %+v, want the move kept", q)
	}
	if n := lab.logged("role start held", "watcher"); n != 1 {
		t.Errorf("the hold was logged %d times, want once", n)
	}
	lab.budgetCleared()
	lab.lookAt(due.Add(time.Minute))
	if r := lab.roleRun("watcher", ""); r == nil || r.Trigger != "moved "+story {
		t.Fatalf("once cleared: %+v, want the run for the move", r)
	}
	if q := lab.queued("watcher"); len(q.Items) != 0 {
		t.Errorf("the start left %+v queued", q.Items)
	}
}

// S-0418: while a role's action is off nothing is queued for it and its
// queue is emptied; the look that finds it turned on takes what happened
// meanwhile as history, and what happens after is queued.
func TestARoleWhoseActionIsOffQueuesNothing(t *testing.T) {
	lab := newRoleEventsLab(t)
	first, second, third := lab.working("First"), lab.working("Second"), lab.working("Third")
	base := time.Now().UTC().Truncate(time.Second)
	lab.lookAt(base)
	lab.moveAt(first, workitem.Review, base.Add(time.Second))
	lab.lookAt(base.Add(2 * time.Second))
	if q := lab.queued("watcher").Items[""]; len(q.Reasons) != 1 {
		t.Fatalf("queue = %+v, want the move", q)
	}
	lab.cfg.RoleActions = map[string]bool{"check": true}
	lab.lookAt(base.Add(3 * time.Second))
	if q, ok := lab.o.Dir.roleQueues()[lab.root]["watcher"]; ok {
		t.Fatalf("turned off, the queue kept %+v", q)
	}
	lab.moveAt(second, workitem.Review, base.Add(4*time.Second))
	lab.lookAt(base.Add(5 * time.Second))
	lab.cfg.RoleActions = map[string]bool{"check": true, "watch": true}
	lab.lookAt(base.Add(6 * time.Second))
	if q := lab.queued("watcher"); !q.On || len(q.Items) != 0 {
		t.Fatalf("turned on, queued %+v, want what happened while off as history", q)
	}
	lab.moveAt(third, workitem.Review, base.Add(7*time.Second))
	lab.lookAt(base.Add(8 * time.Second))
	q := lab.queued("watcher").Items[""]
	if len(q.Reasons) != 1 || q.Reasons[0].ID != third {
		t.Fatalf("queue = %+v, want the move after it was turned on", q)
	}
	if r := lab.roleRun("watcher", ""); r != nil {
		t.Errorf("started: %+v", r)
	}
}

// A role's agents are its name and, for a role that takes items, its run's
// agent on one; a reason's item is the event's, or a thread's, of a kind the
// role takes; and a trigger names each reason once.
func TestARolesAgentsItemsAndTrigger(t *testing.T) {
	watcher := roles.Role{Name: "watcher"}
	checker := roles.Role{Name: "checker", Item: roles.Kinds{workitem.Story}}
	for _, c := range []struct {
		role  roles.Role
		agent string
		want  bool
	}{
		{watcher, "watcher", true},
		{watcher, "watcher-S-0001", false},
		{checker, "checker", true},
		{checker, "checker-S-0001", true},
		{checker, "checker-E-0001", false},
		{checker, "checkers", false},
	} {
		if got := roleAgentOf(c.role, c.agent); got != c.want {
			t.Errorf("roleAgentOf(%s, %s) = %t, want %t", c.role.Name, c.agent, got, c.want)
		}
	}
	thread := inbox.Event{Change: workitem.Change{Kind: inbox.Thread, ID: "TH-0001"}, Fields: map[string]string{inbox.FieldOn: "S-0002"}}
	onPath := inbox.Event{Change: workitem.Change{Kind: inbox.Thread, ID: "TH-0002"}, Fields: map[string]string{inbox.FieldOn: "docs/x.md"}}
	moved := inbox.Event{Change: workitem.Change{Kind: workitem.Moved, ID: "S-0001"}}
	planned := inbox.Event{Change: workitem.Change{Kind: inbox.Planned, ID: "E-0001"}}
	for _, c := range []struct {
		role roles.Role
		e    inbox.Event
		item string
		ok   bool
	}{
		{watcher, planned, "", true},
		{checker, moved, "S-0001", true},
		{checker, thread, "S-0002", true},
		{checker, onPath, "", false},
		{checker, planned, "", false},
	} {
		if item, ok := reasonItem(c.role, c.e); ok != c.ok || (ok && item != c.item) {
			t.Errorf("reasonItem(%s, %s %s) = %q, %t; want %q, %t", c.role.Name, c.e.Kind, c.e.ID, item, ok, c.item, c.ok)
		}
	}
	rs := []orchestrateReason{{Kind: workitem.Moved, ID: "S-0001", To: workitem.Review}, {Kind: workitem.Moved, ID: "S-0001", To: workitem.Backlog}, {Kind: inbox.Thread, ID: "TH-0002"}, {Kind: inbox.ReadyEmpty}}
	if got := roleTrigger(rs); got != "moved S-0001; thread TH-0002; ready-empty" {
		t.Errorf("trigger = %q", got)
	}
}
