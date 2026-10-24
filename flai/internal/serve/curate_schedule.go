package serve

import (
	"cmp"
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The curator after acceptances and on digest.schedule (S-0432, ADR-0199).
//
// A curate scheduler per served project looks whenever the project's
// launcher does: when its work items or threads change, when an agent ends,
// and every minute, as the analysis scheduler does (S-0223). With the
// curate host action on, it starts the curator when it sees a story
// accepted, moved to done, with the trigger "accepted S-nnnn". Acceptances
// seen while a curator runs are remembered, and start one run once it has
// ended, the look the end brings, its trigger naming every one of them in
// ID order ("accepted S-0001, S-0002"). Acceptances seen at one look start
// one run, and so do acceptances and the schedule coming round at one look,
// with the acceptances' trigger.
//
// Each time digest.schedule comes round, it starts the curator with the
// schedule as its trigger ("schedule daily"). Every start goes through the
// checks Curate makes, from the curator's last activity; a start they
// refuse, as while a curator runs, is logged and not tried again: for the
// schedule, until it next comes round. Times missed between two looks come
// round once, and a schedule set or changed comes round first at its next
// time from then.
//
// What it has seen is kept in memory only: acceptances and times that
// passed while flai serve was down do not count. While the curate host
// action is off it starts nothing and keeps up with the acceptances and the
// schedule, so that turning it on does not act on the past.

// curateScheduler starts the curator for one project after acceptances and
// on digest.schedule (S-0432, ADR-0199).
type curateScheduler struct {
	o       Options
	e       Entry
	starter *launcher

	mu sync.Mutex
	// since is the second from which moves are new, and seen those stamped
	// then or later that have been seen, by key, with when.
	since time.Time
	seen  map[string]time.Time
	// cadence follows digest.schedule.
	cadence cadence
	// pending are the stories accepted while a curator ran, in ID order.
	pending []string
	said    map[string]bool // the manifest's refusals warned of
}

// newCurateScheduler is the curate scheduler for project e, starting the
// curator with starter, the project's launcher. Acceptances and times before
// now do not count.
func newCurateScheduler(o Options, e Entry, starter *launcher) *curateScheduler {
	now := o.Now().UTC()
	s := &curateScheduler{o: o, e: e, starter: starter, since: now.Truncate(time.Second), seen: map[string]time.Time{}, cadence: cadence{mark: now}, said: map[string]bool{}}
	if repo, err := workitem.Open(e.Root); err == nil {
		s.cadence.see(repo.Manifest.Digest.Schedule, s.cadence.mark)
	}
	return s
}

// look starts the curator once with the curate host action on when stories
// were accepted since the last look, or while the curator last ran and none
// runs now, or when digest.schedule has come round, unless the checks refuse
// it.
func (s *curateScheduler) look(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.o.Now().UTC()
	repo, err := workitem.Open(s.e.Root)
	if err != nil {
		s.warnOnce("curate scheduler could not read the project", err)
		return
	}
	items, err := repo.List(true)
	if err != nil {
		s.warnOnce("curate scheduler could not read the project", err)
		return
	}
	var changes []workitem.Change
	changes, _, s.since, s.seen = news(items, nil, s.since, s.seen, now)
	sched := s.schedule(repo.Manifest.Digest, now)
	if s.o.Agent == nil || !s.o.Agent(s.e.Root).Curate {
		s.pending, s.cadence.mark = nil, now
		return
	}
	s.pending = accepted(s.pending, changes)
	due := s.cadence.due(sched, now)
	var trigger string
	switch {
	case len(s.pending) > 0 && s.o.Dir.AgentStates()[s.e.Root].Curator.running():
		return // the look its end brings starts the run
	case len(s.pending) > 0:
		trigger = "accepted " + strings.Join(s.pending, ", ")
	case due:
		trigger = scheduleTrigger(sched)
	default:
		return
	}
	s.pending = nil
	cs, err := curateCheck(s.o, s.e, "")
	var no *Refused
	switch {
	case errors.As(err, &no):
		s.starter.log("curator not started on its own", "trigger", trigger, "reason", no.Why)
		return
	case err != nil:
		s.starter.warn("curator not started on its own", "trigger", trigger, "err", err)
		return
	}
	s.starter.mu.Lock()
	defer s.starter.mu.Unlock()
	s.starter.curate(ctx, cs.cfg, cs.agent, cs.since, trigger)
}

// accepted is pending with the stories changes move to done added, each
// once, in ID order.
func accepted(pending []string, changes []workitem.Change) []string {
	for _, c := range changes {
		if c.Type == workitem.Story && c.Kind == workitem.Moved && c.To == workitem.Done && !slices.Contains(pending, c.ID) {
			pending = append(pending, c.ID)
		}
	}
	// zero-padded IDs: a longer one is a larger number
	slices.SortFunc(pending, func(a, b string) int { return cmp.Or(cmp.Compare(len(a), len(b)), strings.Compare(a, b)) })
	return pending
}

// schedule is digest.schedule parsed, nil when it is unset or cannot be
// parsed, which is warned of once. A schedule set or changed comes round
// first at its next time from now.
func (s *curateScheduler) schedule(d manifest.Digest, now time.Time) *cron.Schedule {
	s.cadence.see(d.Schedule, now)
	sched, err := d.DigestSchedule()
	if err != nil {
		s.warnOnce("digest schedule not valid", err)
		return nil
	}
	return sched
}

// warnOnce warns of err with msg the first time it is seen.
func (s *curateScheduler) warnOnce(msg string, err error) {
	if key := msg + "\x00" + err.Error(); !s.said[key] {
		s.said[key] = true
		s.starter.warn(msg, "err", err)
	}
}
