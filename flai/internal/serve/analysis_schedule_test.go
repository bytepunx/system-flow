//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scheduleLab is an analyzeLab with an analysis scheduler whose clock the
// test sets.
type scheduleLab struct {
	*agentLab
	s     *analysisScheduler
	clock time.Time
}

// newScheduleLab is a scheduleLab whose manifest sets analysis.schedule to
// spec, none when it is empty, and whose scheduler starts at at.
func newScheduleLab(t *testing.T, spec string, at time.Time) *scheduleLab {
	t.Helper()
	lab := &scheduleLab{agentLab: analyzeLab(t), clock: at}
	lab.schedule(spec)
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	lab.s = newAnalysisScheduler(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

// schedule sets analysis.schedule to spec, or unsets it when it is empty.
func (lab *scheduleLab) schedule(spec string) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if spec != "" {
		m += "analysis:\n  schedule: \"" + spec + "\"\n"
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// lookAt has the scheduler look at t.
func (lab *scheduleLab) lookAt(t time.Time) {
	lab.clock = t
	lab.s.look(context.Background())
}

// starts counts the analyzer runs started, from the journal.
func (lab *scheduleLab) starts() int {
	n := 0
	for _, e := range lab.entries() {
		if e.Action == hostapi.ActionAnalyze && e.Outcome == "done" {
			n++
		}
	}
	return n
}

// ended waits for the newest analyzer run to end.
func (lab *scheduleLab) ended() {
	lab.t.Helper()
	waitFor(lab.t, "the analyzer run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Ended != "" })
}

// S-0223: analysis.schedule starts the analyzer for all three focuses when
// it comes round, with the schedule as its trigger, once however many times
// were missed between two looks; a time that comes round while a run goes
// starts none, and is not taken again once it ends.
func TestTheAnalysisScheduleStartsTheAnalyzer(t *testing.T) {
	lab := newScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.hold()
	lab.analyzerRuns("2026-10-09-all.md", "Wrote design/analysis/2026-10-09-all.md.", runStart)
	lab.lookAt(lab.clock)
	if n := lab.starts(); n != 0 {
		t.Fatalf("started before the schedule came round: %d", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)) // three midnights missed
	if n := lab.starts(); n != 1 {
		t.Fatalf("started %d runs for the times missed, want one", n)
	}
	run := lab.state().Analyzer
	if run.Trigger != "schedule daily" || run.Focus != "all" {
		t.Errorf("run = %+v, want all three focuses started by schedule daily", run)
	}
	if got := lab.analyzerGiven(); !strings.Contains(got, "focus: unset\n") {
		t.Errorf("the scheduled analyzer was given a focus:\n%s", got)
	}
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 1 {
		t.Errorf("the missed times came round again: %d runs", n)
	}
	// midnight comes round while the run goes: none is started, and it is said
	lab.lookAt(time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 || !lab.state().Analyzer.same(run) {
		t.Errorf("started beside a running analyzer: %d runs, %+v", n, lab.state().Analyzer)
	}
	if !strings.Contains(lab.logText(), `msg="scheduled analyzer not started" component=serve project=t trigger="schedule daily" reason="the analyzer is already running for this project`) {
		t.Errorf("the refusal is not logged:\n%s", lab.logText())
	}
	lab.release("analyze")
	lab.ended()
	waitFor(t, "its activity logged with the schedule as its trigger", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityAnalyzer)
		return err == nil && len(doc.Entries) == 1 && doc.Entries[0].Trigger == "schedule daily"
	})
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 1 {
		t.Errorf("the time that came round during the run was taken after it: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 2 {
		t.Errorf("the next midnight started %d runs in all, want two", n)
	}
	lab.ended()
}

// S-0223: with the analyze action off, the schedule starts nothing and says
// nothing, and turning it on does not act on the times that passed.
func TestTheAnalysisScheduleStartsNothingWhileAnalyzeIsOff(t *testing.T) {
	lab := newScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.cfg.Analyze = false
	lab.lookAt(time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 || lab.state().Analyzer != nil || strings.Contains(lab.logText(), "scheduled analyzer") {
		t.Fatalf("the schedule acted while analyze is off: %d runs, %+v\n%s", n, lab.state().Analyzer, lab.logText())
	}
	lab.cfg.Analyze = true
	lab.lookAt(lab.clock.Add(time.Minute))
	if n := lab.starts(); n != 0 {
		t.Errorf("turning analyze on acted on the past: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 {
		t.Errorf("the next midnight started %d runs, want one", n)
	}
	lab.ended()
}

// S-0223: a schedule set or changed comes round first at its next time from
// then, not at a time that passed since the last came round.
func TestAChangedAnalysisScheduleComesRoundFirstAtItsNextTime(t *testing.T) {
	lab := newScheduleLab(t, "", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock)
	// set a day later: 09:00 on the 7th passed since the scheduler started
	lab.clock = time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC)
	lab.schedule("0 9 * * *")
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 8, 8, 59, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 {
		t.Fatalf("a schedule just set came round at a time before it was set: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 1 {
		t.Fatalf("the schedule set did not come round at its next time: %d runs", n)
	}
	lab.ended()
	// changed at 07:00 on the 9th: 06:00 came round since 09:00 on the 8th
	// under the new schedule, but before it was changed
	lab.clock = time.Date(2026, 10, 9, 7, 0, 0, 0, time.UTC)
	lab.schedule("0 6 * * *")
	lab.lookAt(lab.clock)
	if n := lab.starts(); n != 1 {
		t.Errorf("the changed schedule came round at a time before it changed: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC))
	if n := lab.starts(); n != 2 || lab.state().Analyzer.Trigger != "schedule 0 6 * * *" {
		t.Errorf("the changed schedule did not come round at its next time: %d runs, %+v", n, lab.state().Analyzer)
	}
	lab.ended()
}
