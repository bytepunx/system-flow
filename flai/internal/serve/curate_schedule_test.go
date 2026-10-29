//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// curationLab is a curateLab with a curate scheduler whose clock the test
// sets.
type curationLab struct {
	*agentLab
	s     *curateScheduler
	clock time.Time
}

// newCurationLab is a curationLab whose manifest sets digest.schedule to
// spec, none when it is empty, and whose scheduler starts at at.
func newCurationLab(t *testing.T, spec string, at time.Time) *curationLab {
	t.Helper()
	lab := &curationLab{agentLab: curateLab(t), clock: at}
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if spec != "" {
		m += "digest:\n  schedule: \"" + spec + "\"\n"
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	lab.s = newCurateScheduler(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

// lookAt has the scheduler look at t.
func (lab *curationLab) lookAt(t time.Time) {
	lab.clock = t
	lab.s.look(context.Background())
}

// accept makes a ready story and moves it through to done at the lab's
// clock, as its agent and the operator would, and returns its ID.
func (lab *curationLab) accept(title string) string {
	lab.t.Helper()
	id := lab.ready(title)
	task, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: id, Owner: "builder-" + id, Now: lab.clock})
	if err != nil {
		lab.t.Fatal(err)
	}
	for _, m := range [][2]string{{id, workitem.InProgress}, {task.ID, workitem.Ready}, {task.ID, workitem.InProgress}, {task.ID, workitem.Done}, {id, workitem.Review}, {id, workitem.Done}} {
		it, err := lab.repo.Get(m[0])
		if err != nil {
			lab.t.Fatal(err)
		}
		if m == [2]string{id, workitem.Done} { // its criterion met
			data, _ := os.ReadFile(it.Path)
			if err := os.WriteFile(it.Path, []byte(strings.Replace(string(data), "- [ ] works", "- [x] works", 1)), 0o644); err != nil {
				lab.t.Fatal(err)
			}
			if it, err = lab.repo.Get(id); err != nil {
				lab.t.Fatal(err)
			}
		}
		if _, err := lab.repo.Transition(it, m[1], "alex", "", lab.clock); err != nil {
			lab.t.Fatalf("move %s to %s: %v", m[0], m[1], err)
		}
	}
	return id
}

// starts are the journal entries of the curator runs started, oldest first.
func (lab *curationLab) starts() []hostapi.Entry {
	var out []hostapi.Entry
	for _, e := range lab.entries() {
		if e.Action == hostapi.ActionCurate && e.Outcome == "done" {
			out = append(out, e)
		}
	}
	return out
}

// ended waits for the curator's newest run, the nth started, to end.
func (lab *curationLab) ended(n int) {
	lab.t.Helper()
	waitFor(lab.t, "the curator run ended", func() bool {
		r := lab.state().Curator
		return len(lab.starts()) == n && r != nil && r.Ended != ""
	})
}

// S-0432, ADR-0199: a story accepted with the curate action on and no
// curator running starts the curator once, with the acceptance as its
// trigger, journalled and logged as such, from the curator's last activity.
func TestAnAcceptanceStartsTheCurator(t *testing.T) {
	lab := newCurationLab(t, "", time.Now().UTC())
	lab.lookAt(lab.clock)
	if n := len(lab.starts()); n != 0 {
		t.Fatalf("started with nothing accepted: %d", n)
	}
	lab.curatorRuns("Took in one entry.", runStart)
	story := lab.accept("Accepted")
	lab.lookAt(lab.clock.Add(time.Second))
	starts := lab.starts()
	if len(starts) != 1 {
		t.Fatalf("started %d runs after one acceptance, want one", len(starts))
	}
	run := lab.state().Curator
	if run.Trigger != "accepted "+story || run.Since != "" {
		t.Errorf("run = %+v, want started by accepted %s, from its last activity", run, story)
	}
	if !strings.Contains(starts[0].Detail, "to curate, everything, its first run, as curator (pid ") || !strings.Contains(starts[0].Detail, "), started by accepted "+story+"; log ") {
		t.Errorf("the journal entry does not name its trigger: %q", starts[0].Detail)
	}
	if !strings.Contains(lab.logText(), `msg="curator started" component=serve project=t since="" trigger="accepted `+story+`"`) {
		t.Errorf("the start is not logged with its trigger:\n%s", lab.logText())
	}
	lab.ended(1)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 {
		t.Errorf("the acceptance started the curator again: %d runs", n)
	}
	waitFor(t, "its activity logged with the acceptance as its trigger", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityCurator)
		return err == nil && len(doc.Entries) == 1 && doc.Entries[0].Trigger == "accepted "+story
	})
}

// S-0432, ADR-0199: acceptances while the curate action is off start
// nothing, then or once it is turned on; the next one after does.
func TestAcceptancesWhileCurateIsOffStartNothing(t *testing.T) {
	lab := newCurationLab(t, "", time.Now().UTC())
	lab.cfg.Curate = false
	lab.accept("First")
	lab.lookAt(lab.clock.Add(time.Second))
	lab.accept("Second")
	lab.lookAt(lab.clock.Add(time.Second))
	if n := len(lab.starts()); n != 0 || lab.state().Curator != nil || strings.Contains(lab.logText(), "curator not started") {
		t.Fatalf("acceptances acted while curate is off: %d runs, %+v\n%s", n, lab.state().Curator, lab.logText())
	}
	lab.cfg.Curate = true
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 0 {
		t.Fatalf("turning curate on acted on the acceptances made while it was off: %d runs", n)
	}
	third := lab.accept("Third")
	lab.lookAt(lab.clock.Add(time.Second))
	if n := len(lab.starts()); n != 1 || lab.state().Curator.Trigger != "accepted "+third {
		t.Errorf("an acceptance once curate is on: %d runs, %+v, want one started by accepted %s alone", n, lab.state().Curator, third)
	}
	lab.ended(1)
}

// S-0432, ADR-0199: acceptances seen while a curator runs start nothing
// then, and exactly one run at the look its end brings, its trigger naming
// every one of them in ID order.
func TestAcceptancesDuringARunStartOneRunAfterIt(t *testing.T) {
	lab := newCurationLab(t, "", time.Now().UTC())
	lab.hold()
	if !lab.curateHere("") {
		t.Fatalf("the operator's curator did not start: %+v", lab.state().Curator)
	}
	first := lab.accept("First")
	lab.lookAt(lab.clock.Add(time.Second))
	second := lab.accept("Second")
	lab.lookAt(lab.clock.Add(time.Second))
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 || lab.state().Curator.Trigger != TriggerAsked {
		t.Fatalf("started beside a running curator: %d runs, %+v", n, lab.state().Curator)
	}
	if strings.Contains(lab.logText(), "curator not started") {
		t.Errorf("the acceptances during the run were refused rather than kept:\n%s", lab.logText())
	}
	lab.release("curate")
	lab.ended(1)
	// the look the run's end brings, as the launcher's again does
	lab.lookAt(lab.clock.Add(time.Second))
	starts := lab.starts()
	want := "accepted " + first + ", " + second
	if len(starts) != 2 || lab.state().Curator.Trigger != want {
		t.Fatalf("after the run: %d runs, %+v, want one more started by %q", len(starts), lab.state().Curator, want)
	}
	if !strings.Contains(starts[1].Detail, "started by "+want+"; log ") || !strings.Contains(starts[0].Detail, "started by asked; log ") {
		t.Errorf("the journal entries do not name their triggers: %q, %q", starts[0].Detail, starts[1].Detail)
	}
	lab.ended(2)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 2 {
		t.Errorf("the acceptances started the curator again: %d runs", n)
	}
}

// S-0432, ADR-0199: digest.schedule starts the curator when it comes round,
// with the schedule as its trigger, once however many times were missed; a
// time that comes round while a run goes starts none, is logged, and is not
// taken again once it ends.
func TestTheDigestScheduleStartsTheCurator(t *testing.T) {
	lab := newCurationLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.hold()
	lab.lookAt(lab.clock)
	if n := len(lab.starts()); n != 0 {
		t.Fatalf("started before the schedule came round: %d", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)) // three midnights missed
	starts := lab.starts()
	if len(starts) != 1 {
		t.Fatalf("started %d runs for the times missed, want one", len(starts))
	}
	run := lab.state().Curator
	if run.Trigger != "schedule daily" || !strings.Contains(starts[0].Detail, "started by schedule daily; log ") {
		t.Errorf("run = %+v, journal %q, want started by schedule daily", run, starts[0].Detail)
	}
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 {
		t.Errorf("the missed times came round again: %d runs", n)
	}
	// midnight comes round while the run goes: none is started, and it is said
	lab.lookAt(time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC))
	if n := len(lab.starts()); n != 1 || !lab.state().Curator.same(run) {
		t.Errorf("started beside a running curator: %d runs, %+v", n, lab.state().Curator)
	}
	if !strings.Contains(lab.logText(), `msg="curator not started on its own" component=serve project=t trigger="schedule daily" reason="the curator is already running for this project`) {
		t.Errorf("the refusal is not logged:\n%s", lab.logText())
	}
	lab.release("curate")
	lab.ended(1)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 {
		t.Errorf("the time that came round during the run was taken after it: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC))
	if n := len(lab.starts()); n != 2 {
		t.Errorf("the next midnight started %d runs in all, want two", n)
	}
	lab.ended(2)
}

// S-0432, ADR-0199: acceptances and the schedule coming round at one look
// start one run, with the acceptances' trigger.
func TestAcceptancesAndTheScheduleAtOneLookStartOneRun(t *testing.T) {
	lab := newCurationLab(t, "daily", time.Date(2026, 10, 6, 23, 59, 0, 0, time.UTC))
	lab.clock = time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC)
	stories := []string{lab.accept("Later"), lab.accept("Earlier")}
	slices.Sort(stories)
	lab.lookAt(lab.clock)
	want := "accepted " + strings.Join(stories, ", ")
	if n := len(lab.starts()); n != 1 || lab.state().Curator.Trigger != want {
		t.Fatalf("%d runs, %+v, want one started by %q", n, lab.state().Curator, want)
	}
	lab.ended(1)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 {
		t.Errorf("the schedule came round again after the acceptances' run: %d runs", n)
	}
}

// S-0432: acceptances are kept once each, in the order of their IDs, a
// longer zero-padded ID after a shorter one; a task's or an epic's move to
// done, and a story's to anything else, is none.
func TestAcceptedKeepsEachStoryOnceInIDOrder(t *testing.T) {
	changes := []workitem.Change{
		{ID: "S-10000", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Done},
		{ID: "S-0002", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Done},
		{ID: "S-0003", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Cancelled},
		{ID: "T-0001", Type: workitem.Task, Kind: workitem.Moved, To: workitem.Done},
		{ID: "E-0001", Type: workitem.Epic, Kind: workitem.Moved, To: workitem.Done},
		{ID: "S-0009", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Done},
	}
	if got := accepted([]string{"S-0009"}, changes); !slices.Equal(got, []string{"S-0002", "S-0009", "S-10000"}) {
		t.Errorf("accepted = %v", got)
	}
}

// limitedAt is when I-0212's curator run started, on the harness's usage
// limit, and limitReset when that limit reset.
var (
	limitedAt  = time.Date(2026, 10, 26, 16, 3, 52, 0, time.UTC)
	limitReset = time.Date(2026, 10, 28, 8, 0, 0, 0, time.UTC)
)

// limited records in serve/agents.json the curator's last run as I-0212's
// was: started by trigger at limitedAt, from since, and ended two seconds
// later failed on the usage limit, resetting at limitReset; retry marks it
// a retry itself. With reset false, the harness named no reset time.
func (lab *curationLab) limited(trigger, since string, retry, reset bool) {
	exit := 1
	run := &AgentRun{Agent: workitem.ActivityCurator, Role: "curate", Started: limitedAt.Format(time.RFC3339), Ended: limitedAt.Add(2 * time.Second).Format(time.RFC3339),
		Exit: &exit, Session: "s", Trigger: trigger, Since: since, Outcome: OutcomeFailed, Why: limitWhyText, ResetsAt: limitResetsAt, Retry: retry}
	if !reset {
		run.Why, run.ResetsAt = "ended on the harness's usage limit (usage_limit_reached, 429); it named no reset time", ""
	}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) { s.put(run) })
}

// restart makes the lab's scheduler again, at the lab's clock, as a flai
// serve started again would: with nothing in memory.
func (lab *curationLab) restart() {
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	lab.s = newCurateScheduler(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
}

// I-0212, S-0514, ADR-0214: while the curator's last run ended on the usage
// limit and it has not reset, an acceptance and the schedule coming round
// start nothing; at the first look once it has reset, the run starts again,
// once, from its since, its trigger naming its stories and the one accepted
// meanwhile, marked a retry and journalled as one. The next time the
// schedule comes round starts a run as before.
func TestARunOnTheUsageLimitIsStartedAgainOnceItResets(t *testing.T) {
	lab := newCurationLab(t, "daily", limitedAt.Add(time.Minute))
	lab.limited("accepted S-0498, S-0501", "2026-10-20T09:00:00Z", false, true)
	lab.lookAt(lab.clock.Add(time.Minute))
	lab.clock = time.Date(2026, 10, 27, 10, 0, 0, 0, time.UTC)
	story := lab.accept("Accepted meanwhile")
	lab.lookAt(lab.clock.Add(time.Second))
	lab.lookAt(time.Date(2026, 10, 28, 0, 0, 30, 0, time.UTC)) // midnight came round
	lab.lookAt(limitReset.Add(-time.Second))
	if n := len(lab.starts()); n != 0 || lab.state().Curator.Ended == "" {
		t.Fatalf("started before the limit reset: %d runs, %+v", n, lab.state().Curator)
	}
	if strings.Contains(lab.logText(), "curator not started") {
		t.Errorf("the held starts were tried and refused rather than held:\n%s", lab.logText())
	}
	lab.lookAt(limitReset)
	starts := lab.starts()
	stories := []string{story, "S-0498", "S-0501"}
	slices.Sort(stories)
	want := "accepted " + strings.Join(stories, ", ")
	run := lab.state().Curator
	if len(starts) != 1 || run.Trigger != want || !run.Retry || run.Since != "2026-10-20T09:00:00Z" {
		t.Fatalf("at the reset: %d runs, %+v, want one retry started by %q from the retried run's since", len(starts), run, want)
	}
	if !strings.Contains(starts[0].Detail, "started by "+want+"; started again after the usage limit reset at "+limitResetsAt+"; log ") {
		t.Errorf("the journal entry does not say it is a retry: %q", starts[0].Detail)
	}
	lab.ended(1)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.starts()); n != 1 {
		t.Errorf("started again after the retry: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 29, 0, 0, 30, 0, time.UTC))
	if run := lab.state().Curator; len(lab.starts()) != 2 || run.Trigger != "schedule daily" || run.Retry {
		t.Errorf("the next midnight: %d runs, %+v, want a second started by schedule daily, no retry", len(lab.starts()), run)
	}
	lab.ended(2)
}

// I-0212, S-0514, ADR-0214: a flai serve started after the reset retries at
// its first look, from the run in serve/agents.json, with its trigger.
func TestARunOnTheUsageLimitIsStartedAgainByASchedulerMadeAfterTheReset(t *testing.T) {
	lab := newCurationLab(t, "", limitedAt.Add(time.Minute))
	lab.limited("accepted S-0498, S-0501", "", false, true)
	lab.lookAt(lab.clock.Add(time.Minute))
	lab.clock = limitReset.Add(26 * time.Hour)
	lab.restart()
	lab.lookAt(lab.clock)
	if run := lab.state().Curator; len(lab.starts()) != 1 || run.Trigger != "accepted S-0498, S-0501" || !run.Retry {
		t.Fatalf("first look after the reset: %d runs, %+v, want one retry started by accepted S-0498, S-0501", len(lab.starts()), run)
	}
	lab.ended(1)
}

// S-0514, ADR-0214: while curate is off, a run on the usage limit is not
// retried; at the first look once it is on again, it is.
func TestARunOnTheUsageLimitIsNotStartedAgainWhileCurateIsOff(t *testing.T) {
	lab := newCurationLab(t, "", limitedAt.Add(time.Minute))
	lab.limited(TriggerAsked, "", false, true)
	lab.cfg.Curate = false
	lab.lookAt(limitReset.Add(time.Minute))
	if n := len(lab.starts()); n != 0 || strings.Contains(lab.logText(), "curator not started") {
		t.Fatalf("retried while curate is off: %d runs\n%s", n, lab.logText())
	}
	lab.cfg.Curate = true
	lab.lookAt(lab.clock.Add(time.Minute))
	if run := lab.state().Curator; len(lab.starts()) != 1 || run.Trigger != TriggerAsked || !run.Retry {
		t.Fatalf("once curate is on: %d runs, %+v, want one retry started by asked", len(lab.starts()), run)
	}
	lab.ended(1)
}

// S-0514, ADR-0214: a run that failed with no usage limit, a run on a limit
// that named no reset time, and a retry that ended on the limit again are
// not started again; an acceptance after them starts a run of its own.
func TestRunsNotStartedAgain(t *testing.T) {
	for _, c := range []struct {
		name string
		put  func(lab *curationLab)
	}{
		{"failed with no limit", func(lab *curationLab) {
			exit := 1
			lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
				s.put(&AgentRun{Agent: workitem.ActivityCurator, Role: "curate", Started: limitedAt.Format(time.RFC3339), Ended: limitedAt.Add(time.Minute).Format(time.RFC3339),
					Exit: &exit, Trigger: "accepted S-0498", Outcome: OutcomeFailed, Why: "ended (exit 1)"})
			})
		}},
		{"a limit with no reset time", func(lab *curationLab) { lab.limited("accepted S-0498", "", false, false) }},
		{"a retry on the limit again", func(lab *curationLab) { lab.limited("accepted S-0498", "", true, true) }},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newCurationLab(t, "", limitedAt.Add(time.Minute))
			c.put(lab)
			lab.lookAt(limitReset.Add(time.Minute))
			lab.lookAt(lab.clock.Add(time.Hour))
			if n := len(lab.starts()); n != 0 {
				t.Fatalf("started again: %d runs, %+v", n, lab.state().Curator)
			}
			story := lab.accept("Accepted after")
			lab.lookAt(lab.clock.Add(time.Second))
			if run := lab.state().Curator; len(lab.starts()) != 1 || run.Trigger != "accepted "+story || run.Retry || strings.Contains(lab.starts()[0].Detail, "started again") {
				t.Fatalf("the acceptance after: %d runs, %+v, want one started by accepted %s alone, no retry", len(lab.starts()), run, story)
			}
			lab.ended(1)
		})
	}
}

// S-0514, ADR-0214: a retry's trigger is the retried run's with the
// acceptances held meanwhile joined: beside an acceptances' trigger's own,
// in ID order, in place of a schedule's, and after the operator's asked.
func TestRetryTriggerJoinsTheHeldAcceptances(t *testing.T) {
	for _, c := range []struct {
		trigger string
		pending []string
		want    string
	}{
		{"accepted S-0498, S-0501", nil, "accepted S-0498, S-0501"},
		{"accepted S-0498, S-0501", []string{"S-0499", "S-0501", "S-10000"}, "accepted S-0498, S-0499, S-0501, S-10000"},
		{"schedule daily", nil, "schedule daily"},
		{"schedule daily", []string{"S-0510"}, "accepted S-0510"},
		{TriggerAsked, nil, TriggerAsked},
		{TriggerAsked, []string{"S-0510", "S-0511"}, "asked; accepted S-0510, S-0511"},
	} {
		if got := retryTrigger(c.trigger, c.pending); got != c.want {
			t.Errorf("retryTrigger(%q, %v) = %q, want %q", c.trigger, c.pending, got, c.want)
		}
	}
}
