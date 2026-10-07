package serve

import (
	"fmt"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Stop ends a story's agent on the operator's word (S-0170): flai serve agent
// stop, or Stop on the dashboard's activity page. An agent that runs is asked
// to stop, with everything it started (its process group), and made to once
// o.StopGrace has passed; one that ended waiting for an answer is not started
// again when the answer comes. Either way the run ends stopped, and its story
// is left where it is, with its worktree as the agent left it: it gets no
// agent until the operator retries it or moves it back to ready.
//
// The process is signalled only while it is still the one flai started
// (Owns): a run whose PID the system has given to another process since is
// only recorded as stopped. Stop is refused, saying why, for a story flai
// serve started no agent for, and for one whose agent is not running and not
// waiting for an answer.
func Stop(o Options, e Entry, story string) (*AgentRun, error) {
	if !StoryID.MatchString(story) {
		return nil, refused("%q is not a story's ID", story)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return nil, err
	}
	it, err := repo.Get(story)
	if err != nil {
		return nil, err
	}
	if it.Type != workitem.Story {
		return nil, refused("%s is %s; only a story has an agent", it.ID, it.Type)
	}
	run := o.Dir.AgentStates()[e.Root].Stories[it.ID]
	switch {
	case run == nil:
		return nil, refused("flai serve has started no agent for %s, so there is none to stop", it.ID)
	case run.live(), run.Outcome == OutcomeAsked:
	default:
		return nil, refused("%s's agent is not running: it ended %s, %s", it.ID, run.Ended, outcomeOf(run))
	}
	now := time.Now
	if o.Now != nil {
		now = o.Now
	}
	l := newLauncher(o, e)
	l.handOver = true
	at := now().UTC().Format(time.RFC3339)
	// Marked before the signal, so that the serving flai, which settles the
	// run once the process is gone, records the stop and not a failure.
	o.Dir.updateAgent(e.Root, func(s *AgentState) {
		if cur := s.Stories[it.ID]; cur.same(run) {
			marked := *cur
			marked.Stopped = at
			s.put(&marked)
		}
	})
	signalled, killed := false, false
	if run.running() {
		signalled = true
		_ = TerminateGroup(run.PID)
		if !gone(run.PID, o.stopGrace()) {
			killed = true
			_ = KillGroup(run.PID)
			gone(run.PID, 2*time.Second)
		}
	}
	var out *AgentRun
	o.Dir.updateAgent(e.Root, func(s *AgentState) {
		cur := s.Stories[it.ID]
		if !cur.same(run) {
			return // another agent was started for the story meanwhile
		}
		ended := *cur
		ended.Stopped = at
		if ended.Ended == "" {
			ended.Ended = now().UTC().Format(time.RFC3339)
		}
		ended.stopped()
		s.put(&ended)
		out = &ended
	})
	if out == nil {
		return nil, fmt.Errorf("%s's agent was replaced while it was stopped; it is not recorded as stopped", it.ID)
	}
	what := "recorded it as stopped: it was " + run.awaits() + ", and is not started again when that comes"
	switch {
	case killed:
		what = fmt.Sprintf("killed its process group (pid %d), which had not ended %s after it was asked to", run.PID, o.stopGrace())
	case signalled:
		what = fmt.Sprintf("ended its process group (pid %d)", run.PID)
	case run.live():
		what = fmt.Sprintf("recorded it as stopped: pid %d is no longer the agent flai started", run.PID)
	}
	if o.Host.Record != nil {
		o.Host.Record(hostapi.Entry{At: at, Action: hostapi.ActionAgent, Method: "serve.agent", Project: e.Key, Root: e.Root, By: "flai serve",
			Outcome: "done", Detail: fmt.Sprintf("stopped %s's agent %s: %s", it.ID, run.Agent, what)})
	}
	l.log("agent stopped", "story", it.ID, "pid", run.PID, "signalled", signalled, "killed", killed)
	l.measure(it.ID)
	return out, nil
}

// outcomeOf says how an ended run ended, for a refusal.
func outcomeOf(run *AgentRun) string {
	switch {
	case run.Why != "":
		return run.Outcome + ": " + run.Why
	case run.Outcome != "":
		return run.Outcome
	case run.Error != "":
		return "could not start: " + run.Error
	}
	return "how is not recorded"
}

// stopGrace is how long Stop waits for an agent it asked to stop before it
// kills it.
func (o Options) stopGrace() time.Duration {
	if o.StopGrace > 0 {
		return o.StopGrace
	}
	return 10 * time.Second
}

// gone waits up to d for pid to be gone, and says whether it is.
func gone(pid int, d time.Duration) bool {
	deadline := time.Now().Add(d)
	for Alive(pid) {
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
	return true
}
