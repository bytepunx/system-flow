package serve

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The analyzer on analysis.schedule (S-0223).
//
// An analysis scheduler per served project looks whenever the project's
// launcher does: when its work items or threads change, when an agent ends,
// and every minute, which is when the schedule is seen to come round, as
// the replanner sees planning.schedule (S-0211). Each time analysis.schedule
// comes round with the analyze host action on, it starts the analyzer for
// all three focuses, with the schedule as its trigger ("schedule daily"),
// through the checks Analyze makes; a start they refuse, as while an
// analyzer runs, is logged and not tried again until the schedule next
// comes round. A start the budget state of the analyzer's provider holds is
// not lost: each look tries it again until it starts (S-0534). Times missed
// between two looks come round once, and a schedule set or changed comes
// round first at its next time from then.
//
// When analysis.formula names a formula, the schedule coming round makes
// that formula's instance in place of starting the analyzer (S-0424): it
// runs flai formula run <name> --ready --autocommit --json in the main
// checkout, as the formula scheduler does, journalled under the formula
// action with the trigger "analysis.schedule daily", and the instance's role
// step starts the analyzer through the checks a role's run makes. A run flai
// formula run refuses is logged and not tried again until the schedule next
// comes round.
//
// What it has seen is kept in memory only: times that passed while flai
// serve was down do not come round. While the analyze host action is off it
// starts nothing and keeps up with the schedule, so that turning it on does
// not act on the past.

// analysisScheduler starts the analyzer for one project on
// analysis.schedule (S-0223).
type analysisScheduler struct {
	o       Options
	e       Entry
	starter *launcher
	// run runs flai formula run for analysis.formula, and exe is this
	// flai's executable, as the formula scheduler's are.
	run outputRunner
	exe func() (string, error)

	mu sync.Mutex
	// cadence follows analysis.schedule.
	cadence cadence
	// heldDue is the trigger of the scheduled run the budget state of the
	// analyzer's provider held, which each look tries again until it starts
	// (S-0534); "" when none waits.
	heldDue string
	said    map[string]bool // the manifest's refusals warned of, and the holds logged (heldOnce)
}

// newAnalysisScheduler is the analysis scheduler for project e, starting
// the analyzer with starter, the project's launcher. Times before now do
// not come round.
func newAnalysisScheduler(o Options, e Entry, starter *launcher) *analysisScheduler {
	s := &analysisScheduler{o: o, e: e, starter: starter, run: runOutput, exe: os.Executable, cadence: cadence{mark: o.Now().UTC()}, said: map[string]bool{}}
	if repo, err := workitem.Open(e.Root); err == nil {
		s.cadence.see(repo.Manifest.Analysis.Schedule, s.cadence.mark)
	}
	return s
}

// look starts the analyzer for all three focuses when analysis.schedule has
// come round since the last look, or a run it came round for was held by the
// budget state, and the analyze host action is on, unless the checks refuse
// it. A start the budget state holds is kept for the next look, and logged
// once per state. With analysis.formula set, it makes that formula's
// instance instead.
func (s *analysisScheduler) look(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.o.Now().UTC()
	repo, err := workitem.Open(s.e.Root)
	if err != nil {
		s.warnOnce("analysis scheduler could not read the project", err)
		return
	}
	sched := s.schedule(repo.Manifest.Analysis, now)
	if s.o.Agent == nil || !s.o.Agent(s.e.Root).Analyze {
		s.cadence.mark, s.heldDue = now, ""
		return
	}
	trigger := s.heldDue
	if s.cadence.due(sched, now) {
		trigger = scheduleTrigger(sched)
	}
	if trigger == "" {
		return
	}
	s.heldDue = ""
	if name := repo.Manifest.Analysis.Formula; name != "" {
		s.makeFormula(ctx, name, "analysis."+trigger, now)
		return
	}
	as, err := analyzeCheck(s.o, s.e, "")
	var no *Refused
	switch {
	case errors.As(err, &no):
		s.starter.log("scheduled analyzer not started", "trigger", trigger, "reason", no.Why)
		return
	case err != nil:
		s.starter.warn("scheduled analyzer not started", "trigger", trigger, "err", err)
		return
	}
	s.starter.mu.Lock()
	defer s.starter.mu.Unlock()
	if _, held := s.starter.analyze(ctx, as.cfg, as.agent, "", trigger); held != nil {
		s.heldDue = trigger
		s.starter.heldOnce(s.said, "scheduled analyzer held", *held, "trigger", trigger)
	}
}

// makeFormula makes the instance of the formula analysis.formula names, with
// trigger, through the run the formula scheduler makes (S-0424): this flai's
// formula run <name> --ready --autocommit --json in the main checkout,
// journalled under the formula action, or logged when it is refused.
func (s *analysisScheduler) makeFormula(ctx context.Context, name, trigger string, now time.Time) {
	fs := &formulaScheduler{o: s.o, e: s.e, starter: s.starter, run: s.run, exe: s.exe}
	fs.make(ctx, formulas.Formula{Name: name}, trigger, now)
}

// schedule is analysis.schedule parsed, nil when it is unset or cannot be
// parsed, which is warned of once. A schedule set or changed comes round
// first at its next time from now.
func (s *analysisScheduler) schedule(a manifest.Analysis, now time.Time) *cron.Schedule {
	s.cadence.see(a.Schedule, now)
	sched, err := a.AnalysisSchedule()
	if err != nil {
		s.warnOnce("analysis schedule not valid", err)
		return nil
	}
	return sched
}

// warnOnce warns of err with msg the first time it is seen.
func (s *analysisScheduler) warnOnce(msg string, err error) {
	if key := msg + "\x00" + err.Error(); !s.said[key] {
		s.said[key] = true
		s.starter.warn(msg, "err", err)
	}
}
