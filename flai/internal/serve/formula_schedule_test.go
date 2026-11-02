//go:build !windows

package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// nightlyFormula is a formula with no inputs whose file names a schedule,
// the one the tests' formula scheduler follows.
const nightlyFormula = `---
title: Nightly check
updated: 2026-11-01
status: active
name: nightly
description: Checks the tree every night.
schedule: %s
steps:
  - id: check
    title: Check the tree
    run:
      flai: check --strict
---

# Nightly
`

// formulaRun is one run of the fake flai the formula scheduler ran.
type formulaRun struct {
	dir  string
	env  []string
	name string
	args []string
}

// formulaScheduleLab is an agentLab with a formula scheduler whose clock the
// test sets, the formula action on or off as on says, and a fake flai that
// makes a ready story for each formula run, as flai formula run --ready
// does, unless refuse says why it refuses.
type formulaScheduleLab struct {
	*agentLab
	s      *formulaScheduler
	clock  time.Time
	on     bool
	refuse string
	runs   []formulaRun
}

// newFormulaScheduleLab is a formulaScheduleLab whose nightly formula names
// the schedule spec, with the formula action on, and whose scheduler starts
// at at.
func newFormulaScheduleLab(t *testing.T, spec string, at time.Time) *formulaScheduleLab {
	t.Helper()
	lab := &formulaScheduleLab{agentLab: newAgentLab(t), clock: at, on: true}
	lab.cfg.Enabled = false // no story's agent is started: the instance waits in ready
	lab.schedule(spec)
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	o.Host.Enabled = func(action, root string) bool {
		return lab.on && action == hostapi.ActionFormula && root == lab.root
	}
	lab.s = newFormulaScheduler(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	lab.s.exe = func() (string, error) { return "/opt/flai/bin/flai", nil }
	lab.s.run = func(_ context.Context, dir string, env []string, name string, args ...string) ([]byte, []byte, error) {
		lab.runs = append(lab.runs, formulaRun{dir: dir, env: env, name: name, args: args})
		if lab.refuse != "" {
			return nil, []byte("Error: " + lab.refuse + "\n"), errors.New("exit status 4")
		}
		id := lab.ready("Nightly check")
		return fmt.Appendf(nil, `{"story":%q,"status":"ready","commit":"abc1234"}`, id), nil, nil
	}
	return lab
}

// schedule writes the nightly formula with the schedule spec.
func (lab *formulaScheduleLab) schedule(spec string) {
	lab.t.Helper()
	dir := filepath.Join(lab.root, "design", "workflows", "formulas")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "nightly.md"), fmt.Appendf(nil, nightlyFormula, spec), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// lookAt has the scheduler look at t.
func (lab *formulaScheduleLab) lookAt(t time.Time) {
	lab.clock = t
	lab.s.look(context.Background())
}

// made are the instances the journal says the schedule made.
func (lab *formulaScheduleLab) made() []hostapi.Entry {
	var out []hostapi.Entry
	for _, e := range lab.entries() {
		if e.Action == hostapi.ActionFormula {
			out = append(out, e)
		}
	}
	return out
}

// S-0422: a formula's schedule makes its instance when it comes round with
// the formula action on, by running this flai's formula run --ready
// --autocommit in the main checkout, as flai serve and with no inputs, once
// however many times were missed between two looks, and journals it under
// the formula action with the schedule as its trigger.
func TestAFormulasScheduleMakesItsInstanceReady(t *testing.T) {
	lab := newFormulaScheduleLab(t, "daily", time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC))
	lab.lookAt(lab.clock)
	if len(lab.runs) != 0 {
		t.Fatalf("ran before the schedule came round: %+v", lab.runs)
	}
	lab.lookAt(time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)) // three midnights missed
	if len(lab.runs) != 1 {
		t.Fatalf("ran %d times for the times missed, want once", len(lab.runs))
	}
	run := lab.runs[0]
	if run.dir != lab.root || run.name != "/opt/flai/bin/flai" || strings.Join(run.args, " ") != "formula run nightly --ready --autocommit --json" {
		t.Errorf("ran %s %v in %s, want this flai's formula run nightly --ready --autocommit --json in the main checkout", run.name, run.args, run.dir)
	}
	if !slices.Contains(run.env, "FLAI_AGENT="+queueAgent) || slices.ContainsFunc(run.env, func(kv string) bool { return strings.HasPrefix(kv, "FLAI_STORY=") }) {
		t.Errorf("the run is not flai serve's: %v", run.env)
	}
	repo, err := workitem.Open(lab.root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := repo.List(false)
	if err != nil {
		t.Fatal(err)
	}
	var stories []*workitem.Item
	for _, it := range all {
		if it.Type == workitem.Story && it.Status == workitem.Ready {
			stories = append(stories, it)
		}
	}
	if len(stories) != 1 {
		t.Fatalf("want the instance in ready: %v", stories)
	}
	made := lab.made()
	if len(made) != 1 || made[0].Method != "serve.formula" || made[0].Outcome != "done" || made[0].By != "flai serve" ||
		made[0].Detail != "made "+stories[0].ID+", the instance of formula nightly, in ready, started by schedule daily; committed abc1234" {
		t.Errorf("journal: %+v", made)
	}
	if !strings.Contains(lab.logText(), `msg="scheduled formula instance made" component=serve project=t formula=nightly trigger="schedule daily" story=`+stories[0].ID) {
		t.Errorf("the instance is not logged:\n%s", lab.logText())
	}
	lab.lookAt(lab.clock.Add(time.Minute))
	if len(lab.runs) != 1 {
		t.Errorf("the missed times came round again: %d runs", len(lab.runs))
	}
	lab.lookAt(time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC))
	if len(lab.runs) != 2 || len(lab.made()) != 2 {
		t.Errorf("the next midnight ran %d times in all, want two", len(lab.runs))
	}
}

// S-0422: with the formula action off, the schedule makes nothing and says
// nothing, and turning it on does not act on the times that passed.
func TestAFormulasScheduleMakesNothingWhileTheActionIsOff(t *testing.T) {
	lab := newFormulaScheduleLab(t, "daily", time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC))
	lab.on = false
	lab.lookAt(time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC))
	if len(lab.runs) != 0 || len(lab.made()) != 0 || strings.Contains(lab.logText(), "scheduled formula") {
		t.Fatalf("the schedule acted while the action is off: %+v\n%s", lab.runs, lab.logText())
	}
	lab.on = true
	lab.lookAt(lab.clock.Add(time.Minute))
	if len(lab.runs) != 0 {
		t.Errorf("turning the action on acted on the past: %d runs", len(lab.runs))
	}
	lab.lookAt(time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC))
	if len(lab.runs) != 1 {
		t.Errorf("the next midnight ran %d times, want once", len(lab.runs))
	}
}

// S-0422: a run flai formula run refuses, as for a formula flai check
// refuses, is logged at info with what it said, journals nothing, and is not
// tried again until the schedule next comes round.
func TestARefusedScheduledFormulaWaitsForItsNextTime(t *testing.T) {
	lab := newFormulaScheduleLab(t, "0 6 * * *", time.Date(2026, 11, 2, 5, 0, 0, 0, time.UTC))
	lab.refuse = "formula nightly has 1 problem, so nothing was made"
	lab.lookAt(time.Date(2026, 11, 2, 6, 0, 0, 0, time.UTC))
	if len(lab.runs) != 1 || len(lab.made()) != 0 {
		t.Fatalf("runs %d, journal %+v", len(lab.runs), lab.made())
	}
	if !strings.Contains(lab.logText(), `level=INFO msg="scheduled formula not run" component=serve project=t formula=nightly trigger="schedule 0 6 * * *" reason="Error: formula nightly has 1 problem, so nothing was made"`) {
		t.Errorf("the refusal is not logged at info:\n%s", lab.logText())
	}
	lab.refuse = ""
	lab.lookAt(time.Date(2026, 11, 2, 7, 0, 0, 0, time.UTC))
	if len(lab.runs) != 1 {
		t.Errorf("the refused run was tried again before its next time: %d runs", len(lab.runs))
	}
	lab.lookAt(time.Date(2026, 11, 3, 6, 0, 0, 0, time.UTC))
	if len(lab.runs) != 2 || len(lab.made()) != 1 {
		t.Errorf("the next time made %d runs and %d instances, want two and one", len(lab.runs), len(lab.made()))
	}
}

// S-0422: a schedule that does not parse is warned of once and makes
// nothing; a formula with no schedule is not followed.
func TestAFormulaScheduleThatDoesNotParseMakesNothing(t *testing.T) {
	lab := newFormulaScheduleLab(t, "nightly", time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC))
	lab.lookAt(time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC))
	lab.lookAt(time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC))
	if len(lab.runs) != 0 {
		t.Fatalf("a schedule that does not parse ran: %+v", lab.runs)
	}
	if n := strings.Count(lab.logText(), `msg="formula schedule not valid"`); n != 1 {
		t.Errorf("warned %d times, want once:\n%s", n, lab.logText())
	}
	lab.schedule(`""`)
	lab.lookAt(time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC))
	if len(lab.runs) != 0 || len(lab.s.follows) != 0 {
		t.Errorf("a formula with no schedule is followed: %+v %v", lab.runs, lab.s.follows)
	}
}
