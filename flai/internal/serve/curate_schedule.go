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
// What is tried again is a run that ended on the harness's usage limit
// (S-0514, ADR-0214). While the curator's last run ended on it and the
// limit has not reset, it starts nothing: acceptances are held, and a
// schedule time that comes round is taken as seen and folded into the
// retry. At the first look once the limit has reset, it starts that run
// again, once: from its since, with its trigger and the held acceptances
// joined, marked a retry. A retry that ends on the limit again is not
// retried, nor is a run on a limit that named no reset time. A retry the
// checks refuse is logged and not tried again while this scheduler lives.
//
// What it has seen is kept in memory only: acceptances and times that
// passed while flai serve was down do not count. The hold and the retry are
// read from the curator's last run in serve/agents.json, so a flai serve
// started meanwhile still holds, and retries at its first look after the
// reset. While the curate host action is off it starts nothing, retries
// nothing, and keeps up with the acceptances and the schedule, so that
// turning it on does not act on the past; a retry that came due meanwhile
// starts at the first look once it is on.

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
	// pending are the stories accepted while a curator ran, or while its
	// last run's usage limit had not reset, in ID order.
	pending []string
	// retried is the start of the run on the usage limit a retry was tried
	// for, so that one the checks refuse is not tried at every look.
	retried string
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
// it. While the curator's last run ended on the harness's usage limit and the
// limit has not reset, it starts nothing and holds both; once it has reset,
// it starts that run again, with the held acceptances, once.
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
	last := s.o.Dir.AgentStates()[s.e.Root].Curator
	resets := limitResetAt(last)
	if now.Before(resets) {
		return // held: the acceptances and a time come round wait for the retry
	}
	retry := !resets.IsZero() && !last.Retry && s.retried != last.Started
	var trigger, since, resetAt string
	switch {
	case len(s.pending) > 0 && last.running():
		return // the look its end brings starts the run
	case retry:
		trigger, since, resetAt, s.retried = retryTrigger(last.Trigger, s.pending), last.Since, last.ResetsAt, last.Started
	case len(s.pending) > 0:
		trigger = acceptedTrigger(s.pending)
	case due:
		trigger = scheduleTrigger(sched)
	default:
		return
	}
	s.pending = nil
	cs, err := curateCheck(s.o, s.e, since)
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
	s.starter.curate(ctx, cs.cfg, cs.agent, cs.since, trigger, resetAt)
}

// limitResetAt is when the usage limit run ended on resets, when run ended
// failed on the harness's usage limit and the harness named a reset time
// (ADR-0214), and the zero time otherwise.
func limitResetAt(run *AgentRun) time.Time {
	if run == nil || run.Outcome != OutcomeFailed || run.ResetsAt == "" {
		return time.Time{}
	}
	resets, err := time.Parse(time.RFC3339, run.ResetsAt)
	if err != nil {
		return time.Time{}
	}
	return resets
}

// acceptedPrefix begins the trigger of a run acceptances started.
const acceptedPrefix = "accepted "

// acceptedTrigger is the trigger of a run the stories accepted started,
// "accepted S-0001, S-0002", the stories in the order given.
func acceptedTrigger(stories []string) string {
	return acceptedPrefix + strings.Join(stories, ", ")
}

// retryTrigger is the trigger of the run that starts again one started by
// trigger, with the stories accepted while it was held, pending, joined
// (ADR-0214): an acceptances' trigger names them beside its own, in ID
// order; a schedule's gives way to theirs, as when both come round at one
// look (ADR-0199); any other, the operator's asked, is kept, with theirs
// after it, separated by a semicolon. With none held, it is trigger.
func retryTrigger(trigger string, pending []string) string {
	switch {
	case len(pending) == 0:
		return trigger
	case strings.HasPrefix(trigger, acceptedPrefix):
		stories := slices.Clone(pending)
		for _, id := range strings.Split(strings.TrimPrefix(trigger, acceptedPrefix), ", ") {
			if !slices.Contains(stories, id) {
				stories = append(stories, id)
			}
		}
		return acceptedTrigger(inIDOrder(stories))
	case strings.HasPrefix(trigger, "schedule "):
		return acceptedTrigger(pending)
	default:
		return trigger + "; " + acceptedTrigger(pending)
	}
}

// accepted is pending with the stories changes move to done added, each
// once, in ID order.
func accepted(pending []string, changes []workitem.Change) []string {
	for _, c := range changes {
		if c.Type == workitem.Story && c.Kind == workitem.Moved && c.To == workitem.Done && !slices.Contains(pending, c.ID) {
			pending = append(pending, c.ID)
		}
	}
	return inIDOrder(pending)
}

// inIDOrder sorts ids in the order of their numbers and returns them.
func inIDOrder(ids []string) []string {
	// zero-padded IDs: a longer one is a larger number
	slices.SortFunc(ids, func(a, b string) int { return cmp.Or(cmp.Compare(len(a), len(b)), strings.Compare(a, b)) })
	return ids
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
