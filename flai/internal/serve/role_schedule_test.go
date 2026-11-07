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

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// roleScheduleLab is a roleLab whose auditor and reviewer are started on
// roles.<name>.schedule, with a role scheduler whose clock the test sets.
type roleScheduleLab struct {
	*agentLab
	s     *roleScheduler
	clock time.Time
}

// newRoleScheduleLab is a roleScheduleLab whose manifest sets the auditor's
// schedule to spec, none when it is empty, and whose scheduler starts at at.
func newRoleScheduleLab(t *testing.T, spec string, at time.Time) *roleScheduleLab {
	t.Helper()
	lab := &roleScheduleLab{agentLab: roleLab(t), clock: at}
	dir := filepath.Join(lab.root, "design", "workflows", "roles")
	for name, body := range map[string]string{"auditor.md": auditorRole, "reviewer.md": reviewerRole} {
		role := strings.TrimSuffix(name, ".md")
		body = strings.Replace(body, "  - asked\n", "  - {schedule: roles."+role+".schedule}\n", 1)
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lab.schedule(spec, "")
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	lab.s = newRoleScheduler(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

// schedule sets roles.auditor.schedule to auditor and roles.reviewer.schedule
// to reviewer, each unset when it is empty.
func (lab *roleScheduleLab) schedule(auditor, reviewer string) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nroles:\n"
	for name, spec := range map[string]string{"auditor": auditor, "reviewer": reviewer} {
		if spec != "" {
			m += "  " + name + ":\n    schedule: \"" + spec + "\"\n"
		}
	}
	if auditor == "" && reviewer == "" {
		m = strings.TrimSuffix(m, "roles:\n")
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// lookAt has the scheduler look at t.
func (lab *roleScheduleLab) lookAt(t time.Time) {
	lab.clock = t
	lab.s.look(context.Background())
}

// starts counts the auditor's runs started, from the journal.
func (lab *roleScheduleLab) starts() int {
	n := 0
	for _, e := range lab.entries() {
		if e.Method == MethodRole && strings.HasPrefix(e.Detail, "started ") {
			n++
		}
	}
	return n
}

// ended waits for the auditor's newest run to end.
func (lab *roleScheduleLab) ended() {
	lab.t.Helper()
	waitFor(lab.t, "the auditor's run ended", func() bool { r := lab.roleRun("auditor", ""); return r != nil && r.Ended != "" })
}

// S-0418: roles.<name>.schedule starts the role on no item when it comes
// round, with the schedule as its trigger on its run and its activity, once
// however many times were missed between two looks, and not before.
func TestARoleScheduleStartsTheRole(t *testing.T) {
	lab := newRoleScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamFinal("Audited the project.", 2000, 0.5)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-auditor"), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 6, 23, 59, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 {
		t.Fatalf("started before the schedule came round: %d", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)) // three midnights missed
	if n := lab.starts(); n != 1 {
		t.Fatalf("started %d runs for the times missed, want one", n)
	}
	if run := lab.roleRun("auditor", ""); run == nil || run.Trigger != "schedule daily" || run.Item != "" {
		t.Errorf("run = %+v, want one on no item started by schedule daily", run)
	}
	lab.ended()
	waitFor(t, "its activity logged with the schedule as its trigger", func() bool {
		doc, err := lab.repo.Activity("auditor")
		return err == nil && len(doc.Entries) == 1 && doc.Entries[0].Trigger == "schedule daily"
	})
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 1 {
		t.Errorf("the missed times came round again: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 2 {
		t.Errorf("the next midnight started %d runs in all, want two", n)
	}
	lab.ended()
}

// S-0418: with the role's host action off, its schedule starts nothing and
// says nothing, and turning it on does not act on the times that passed.
func TestARoleScheduleStartsNothingWhileItsActionIsOff(t *testing.T) {
	lab := newRoleScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.cfg.RoleActions = map[string]bool{"review": true}
	lab.lookAt(time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 || lab.roleRun("auditor", "") != nil || strings.Contains(lab.logText(), "scheduled role") {
		t.Fatalf("the schedule acted while audit is off: %d runs\n%s", n, lab.logText())
	}
	lab.cfg.RoleActions = map[string]bool{"review": true, "audit": true}
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 0 {
		t.Errorf("turning audit on acted on the past: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 {
		t.Errorf("the next midnight started %d runs, want one", n)
	}
	lab.ended()
}

// S-0418: a time that comes round while the role runs is refused, logged
// with the role, the trigger, and why, and not tried again once the run
// ends; the next time starts it.
func TestARefusedScheduledRoleWaitsForTheNextTime(t *testing.T) {
	lab := newRoleScheduleLab(t, "0 6 * * *", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.hold()
	lab.lookAt(time.Date(2026, 10, 7, 6, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 {
		t.Fatalf("started %d runs, want one", n)
	}
	lab.lookAt(time.Date(2026, 10, 8, 6, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 {
		t.Errorf("started beside a running auditor: %d runs", n)
	}
	if !strings.Contains(lab.logText(), `msg="scheduled role not started" component=serve project=t role=auditor trigger="schedule 0 6 * * *" reason="the role auditor is already running on no item`) {
		t.Errorf("the refusal is not logged:\n%s", lab.logText())
	}
	lab.release("auditor")
	lab.ended()
	lab.lookAt(lab.clock.Add(time.Minute))
	lab.lookAt(lab.clock.Add(time.Hour))
	if n := lab.starts(); n != 1 {
		t.Errorf("the refused time was tried again before the next: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 6, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 2 {
		t.Errorf("the next time started %d runs in all, want two", n)
	}
	lab.ended()
}

// S-0418: a schedule set later comes round first at its next time from
// then; a role whose file takes items is not started under an entry with no
// condition, and that is said once each time its schedule comes round.
func TestASetRoleScheduleComesRoundAtItsNextTime(t *testing.T) {
	lab := newRoleScheduleLab(t, "", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock)
	lab.clock = time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC)
	lab.schedule("0 9 * * *", "0 9 * * *")
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 8, 8, 59, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 || strings.Contains(lab.logText(), "scheduled role") {
		t.Fatalf("a schedule just set came round at a time before it was set: %d runs\n%s", n, lab.logText())
	}
	lab.lookAt(time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 1 || lab.roleRun("auditor", "").Trigger != "schedule 0 9 * * *" {
		t.Fatalf("the schedule set did not start the auditor at its next time: %d runs", n)
	}
	lab.ended()
	if n := strings.Count(lab.logText(), `msg="scheduled role not started" component=serve project=t role=reviewer trigger="schedule 0 9 * * *" reason="the role runs on an item`); n != 1 {
		t.Errorf("the reviewer, which takes items, was said not started %d times, want once:\n%s", n, lab.logText())
	}
	if r := lab.state().Roles["reviewer"]; len(r) != 0 {
		t.Errorf("the reviewer was started on %v", r)
	}
}

// S-0418: a time the schedule comes round at while the budget state holds
// starts nothing, logged once, and the first look once it clears starts it,
// once, with the schedule as its trigger.
func TestAScheduledRoleWaitsWhileTheBudgetStateHoldsIt(t *testing.T) {
	lab := newRoleScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.budgetHeldOn(harness.Command, anHour())
	lab.lookAt(time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC))
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 0 || lab.roleRun("auditor", "") != nil {
		t.Fatalf("started while held: %d runs", n)
	}
	if n := strings.Count(lab.logText(), `msg="scheduled role held"`); n != 1 {
		t.Errorf("the hold was logged %d times, want once:\n%s", n, lab.logText())
	}
	lab.budgetCleared()
	lab.lookAt(lab.clock.Add(time.Minute))
	if run := lab.roleRun("auditor", ""); lab.starts() != 1 || run == nil || run.Trigger != "schedule daily" {
		t.Fatalf("once the state cleared: %d runs, %+v", lab.starts(), run)
	}
	lab.ended()
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 1 {
		t.Errorf("the held time started %d runs, want one", n)
	}
}

// S-0418: a schedule that does not parse is warned of once and starts
// nothing.
func TestARoleScheduleThatDoesNotParseStartsNothing(t *testing.T) {
	lab := newRoleScheduleLab(t, "61 * * * *", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 {
		t.Errorf("a schedule that does not parse started %d runs", n)
	}
	if n := strings.Count(lab.logText(), `msg="role schedule not valid" component=serve project=t role=auditor err="roles.auditor.schedule \"61 * * * *\"`); n != 1 {
		t.Errorf("warned %d times, want once:\n%s", n, lab.logText())
	}
}

// eachReadyStory has the reviewer's schedule entry run it once for every
// ready story (ADR-0250 §5).
func (lab *roleScheduleLab) eachReadyStory() {
	lab.t.Helper()
	path := filepath.Join(lab.root, "design", "workflows", "roles", "reviewer.md")
	data, err := os.ReadFile(path)
	if err != nil {
		lab.t.Fatal(err)
	}
	body := strings.Replace(string(data), "{schedule: roles.reviewer.schedule}", "{schedule: roles.reviewer.schedule, when: each-ready-story}", 1)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// pullOrder sets the board's pull order.
func (lab *roleScheduleLab) pullOrder(order ...string) {
	lab.t.Helper()
	b, err := lab.repo.LoadBoard()
	if err != nil {
		lab.t.Fatal(err)
	}
	b.Order = order
	if err := b.Save(lab.clock.Format("2006-01-02")); err != nil {
		lab.t.Fatal(err)
	}
}

// reviewerStarts are the stories the reviewer's runs were started on, in
// the order they started, from the journal.
func (lab *roleScheduleLab) reviewerStarts(stories ...string) []string {
	var out []string
	for _, e := range lab.entries() {
		if e.Method != MethodRole || !strings.HasPrefix(e.Detail, "started ") {
			continue
		}
		for _, id := range stories {
			if strings.Contains(e.Detail, id) {
				out = append(out, id)
			}
		}
	}
	return out
}

// reviewed waits for the reviewer's run on each of stories to end.
func (lab *roleScheduleLab) reviewed(stories ...string) {
	lab.t.Helper()
	for _, id := range stories {
		waitFor(lab.t, "the reviewer's run on "+id+" ended", func() bool { r := lab.roleRun("reviewer", id); return r != nil && r.Ended != "" })
	}
}

// S-0571 (I-0262): a defined role that takes stories, under a schedule entry
// with when: each-ready-story, is started once for each ready story in the
// board's pull order when its schedule comes round, each on its story with
// the schedule as its trigger, and never on a story that is not ready.
func TestAnEachReadyStoryScheduleStartsTheRoleOnEveryReadyStory(t *testing.T) {
	lab := newRoleScheduleLab(t, "", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.eachReadyStory()
	first, second := lab.ready("First ready"), lab.ready("Second ready")
	waiting := lab.backlog("Still in backlog", nil)
	lab.pullOrder(second, first)
	lab.schedule("", "0 9 * * *")
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC))
	if got := lab.reviewerStarts(first, second, waiting); !slices.Equal(got, []string{second, first}) {
		t.Fatalf("the reviewer was started on %v, want %v in pull order\n%s", got, []string{second, first}, lab.logText())
	}
	for _, id := range []string{first, second} {
		if run := lab.roleRun("reviewer", id); run == nil || run.Trigger != "schedule 0 9 * * *" {
			t.Errorf("the run on %s = %+v, want one started by schedule 0 9 * * *", id, run)
		}
	}
	if lab.roleRun("reviewer", waiting) != nil {
		t.Errorf("the reviewer was started on %s, which is in backlog", waiting)
	}
	if strings.Contains(lab.logText(), `msg="scheduled role not started"`) {
		t.Errorf("a start was refused:\n%s", lab.logText())
	}
	lab.reviewed(first, second)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.reviewerStarts(first, second)); n != 2 {
		t.Errorf("the reviewer was started %d times by one time coming round, want two", n)
	}
}

// S-0571: a start of a role that takes stories that the budget state holds
// is kept for its story and tried again at each look until the state clears,
// once, and dropped for a story no longer ready.
func TestAnEachReadyStoryStartHeldByTheBudgetIsKeptPerStory(t *testing.T) {
	lab := newRoleScheduleLab(t, "", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.eachReadyStory()
	kept, left := lab.ready("Kept ready"), lab.ready("Left ready")
	lab.schedule("", "0 9 * * *")
	lab.lookAt(lab.clock)
	lab.budgetHeldOn(harness.Command, anHour())
	lab.lookAt(time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock.Add(time.Minute))
	if got := lab.reviewerStarts(kept, left); len(got) != 0 {
		t.Fatalf("started while held: %v", got)
	}
	if n := strings.Count(lab.logText(), `msg="scheduled role held"`); n != 1 {
		t.Errorf("the hold was logged %d times, want once:\n%s", n, lab.logText())
	}
	st, err := lab.repo.Get(left)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lab.repo.Transition(st, workitem.Backlog, "alex", "", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.budgetCleared()
	lab.lookAt(lab.clock.Add(time.Minute))
	if got := lab.reviewerStarts(kept, left); !slices.Equal(got, []string{kept}) {
		t.Fatalf("once the state cleared the reviewer was started on %v, want only %s", got, kept)
	}
	if run := lab.roleRun("reviewer", kept); run == nil || run.Trigger != "schedule 0 9 * * *" {
		t.Errorf("the held run = %+v, want one started by schedule 0 9 * * *", run)
	}
	lab.reviewed(kept)
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := len(lab.reviewerStarts(kept, left)); n != 1 {
		t.Errorf("the held start was tried again after it started: %d runs", n)
	}
}

// S-0571: a role that takes no item under each-ready-story alone is not
// started, and that is said when its schedule comes round.
func TestAnEachReadyStoryScheduleStartsNoRoleThatTakesNoStory(t *testing.T) {
	lab := newRoleScheduleLab(t, "0 9 * * *", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	path := filepath.Join(lab.root, "design", "workflows", "roles", "auditor.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), "{schedule: roles.auditor.schedule}", "{schedule: roles.auditor.schedule, when: each-ready-story}", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.ready("A ready story")
	lab.lookAt(time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 {
		t.Errorf("the auditor, which takes no story, was started %d times", n)
	}
	if !strings.Contains(lab.logText(), `msg="scheduled role not started" component=serve project=t role=auditor trigger="schedule 0 9 * * *" reason="the role takes no story`) {
		t.Errorf("why it was not started is not logged:\n%s", lab.logText())
	}
}
