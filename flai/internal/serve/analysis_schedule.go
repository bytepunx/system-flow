package serve

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The analyzer on the schedules its role file names (S-0223, ADR-0250 §4).
//
// An analysis scheduler per served project looks whenever the project's
// launcher does: when its work items or threads change, when an agent ends,
// and every minute, which is when a schedule is seen to come round, as the
// replanner sees planning.schedule (S-0211). At each look it reads the
// analyzer's role file, the project's or the baseline flai carries
// (roles.Builtin), and follows each setting the file's schedule entries
// name, each with a cadence of its own: analysis.schedule in the baseline's
// {schedule: analysis.schedule}, or roles.analyzer.schedule. Each time a
// setting named by an entry with no condition comes round with the analyze
// host action on, it starts the analyzer for all three focuses, with the
// schedule as its trigger ("schedule daily"), through the checks Analyze
// makes; settings that come round at one look start one run, the first's in
// the file's order. A start they refuse, as while an analyzer runs, is
// logged and not tried again until a schedule next comes round. A start the
// budget state of the analyzer's provider holds is not lost: each look tries
// it again until it starts (S-0534). Times missed between two looks come
// round once, and a schedule set or changed, or newly named, comes round
// first at its next time from then. A file without a schedule entry starts
// no analysis on any schedule; the operator's asked still does. A file that
// cannot be read is warned of once and starts nothing.
//
// When analysis.formula names a formula, a schedule coming round makes that
// formula's instance in place of starting the analyzer (S-0424): it runs
// flai formula run <name> --ready --autocommit --json in the main checkout,
// as the formula scheduler does, journalled under the formula action with
// the setting and the schedule as its trigger ("analysis.schedule daily"),
// and the instance's role step starts the analyzer through the checks a
// role's run makes. A run flai formula run refuses is logged and not tried
// again until the schedule next comes round.
//
// What it has seen is kept in memory only: times that passed while flai
// serve was down do not come round. While the analyze host action is off it
// starts nothing and keeps up with the schedules, so that turning it on does
// not act on the past.

// analysisScheduler starts the analyzer for one project on the schedules
// its role file names (S-0223, ADR-0250 §4).
type analysisScheduler struct {
	o       Options
	e       Entry
	starter *launcher
	// run runs flai formula run for analysis.formula, and exe is this
	// flai's executable, as the formula scheduler's are.
	run outputRunner
	exe func() (string, error)

	mu sync.Mutex
	// cadences follow the settings the analyzer's schedule entries name, by
	// setting.
	cadences map[string]*cadence
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
	now := o.Now().UTC()
	s := &analysisScheduler{o: o, e: e, starter: starter, run: runOutput, exe: os.Executable, cadences: map[string]*cadence{}, said: map[string]bool{}}
	if repo, err := workitem.Open(e.Root); err == nil {
		if analyzer, _, err := roles.Builtin(repo.Root, repo.Manifest.Layout["design"], workitem.ActivityAnalyzer); err == nil {
			s.cadences = startCadences(analyzer, repo.Manifest, now)
		}
	}
	return s
}

// look starts the analyzer for all three focuses when a schedule its role
// file names with no condition has come round since the last look, or a run
// one came round for was held by the budget state, and the analyze host
// action is on, unless the checks refuse it. A start the budget state holds
// is kept for the next look, and logged once per state. With
// analysis.formula set, it makes that formula's instance instead.
func (s *analysisScheduler) look(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.o.Now().UTC()
	repo, err := workitem.Open(s.e.Root)
	if err != nil {
		s.warnOnce("analysis scheduler could not read the project", err)
		return
	}
	analyzer, _, err := roles.Builtin(repo.Root, repo.Manifest.Layout["design"], workitem.ActivityAnalyzer)
	if err != nil {
		s.warnOnce("analyzer role file not read", err) // no entry starts the analyzer
	}
	var scheds map[string]*cron.Schedule
	scheds, s.cadences = followSchedules(analyzer, s.cadences, repo.Manifest, now, func(err error) { s.warnOnce("analysis schedule not valid", err) })
	if s.o.Agent == nil || !s.o.Agent(s.e.Root).Analyze {
		markAll(s.cadences, now)
		s.heldDue = ""
		return
	}
	trigger, setting := s.heldDue, ""
	for _, set := range scheduleSettings(analyzer) {
		sched, ok := scheds[set]
		if !ok || !s.cadences[set].due(sched, now) {
			continue
		}
		if setting == "" && servedPlainly(analyzer.Trigger, roles.Occurrence{Schedule: set}) {
			trigger, setting = scheduleTrigger(sched), set
		}
	}
	if trigger == "" {
		return
	}
	s.heldDue = ""
	if name := repo.Manifest.Analysis.Formula; name != "" {
		if setting == "" { // a held start, which the formula's instance takes
			setting = "analysis.schedule"
		}
		s.makeFormula(ctx, name, strings.TrimSuffix(setting, "schedule")+trigger, now)
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

// warnOnce warns of err with msg the first time it is seen.
func (s *analysisScheduler) warnOnce(msg string, err error) {
	if key := msg + "\x00" + err.Error(); !s.said[key] {
		s.said[key] = true
		s.starter.warn(msg, "err", err)
	}
}
