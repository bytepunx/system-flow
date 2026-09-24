package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting a story's agent on the operator's word (S-0116, ADR-0043).
//
// The launcher starts an agent when a story enters ready. The operator can
// also have one started from a shell or the story's page: flai serve agent
// restart, for a story whose agent dropped or failed. It runs in a process
// of its own, not in flai serve. It starts the agent the way the launcher
// does and records the run in serve/agents.json. The serving flai then
// tracks the run as one of its own: the dot on the card, the outcome once
// settleOrphans sees the process gone, and the resume on an answer.

// Refused is a start the operator asked for that is not allowed now; its
// message says why, and what would allow it.
type Refused struct{ Why string }

func (r *Refused) Error() string { return r.Why }

func refused(format string, args ...any) error { return &Refused{Why: fmt.Sprintf(format, args...)} }

// StartNow starts story's agent in project e as its files name it, records
// the run in serve/agents.json, and journals it, as the launcher does.
// Whether it should be started is the caller's to judge; restart says in the
// agent's prompt why it was started again, empty when it entered ready.
func StartNow(ctx context.Context, o Options, e Entry, story, restart string) (*AgentRun, error) {
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
	l := newLauncher(o, e)
	l.start(ctx, o.Agent(e.Root), readyStory{ID: it.ID, Agent: it.Agent, Restart: restart})
	run := o.Dir.AgentStates()[e.Root].Stories[it.ID]
	if run == nil {
		return nil, fmt.Errorf("the run for %s was not recorded in %s", it.ID, o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
}

// Restart starts a new agent, in a new session, for a story in ready or in
// progress whose last agent flai serve started has ended or dropped. It is
// refused while the agent action is off for the project, when the story is
// in another state, when no agent was started for it, while its agent runs
// or waits for an answer, when nothing can start it, and for a story in
// ready while the in-progress limit is full.
func Restart(ctx context.Context, o Options, e Entry, story string) (*AgentRun, error) {
	cfg := o.Agent(e.Root)
	if !cfg.Enabled {
		return nil, refused("the agent host action is off for this project: flai serve enable agent")
	}
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
	if it.Status != workitem.Ready && it.Status != workitem.InProgress {
		return nil, refused("%s is in %s; only a story in ready or in-progress is restarted", it.ID, it.Status)
	}
	st := o.Dir.AgentStates()[e.Root]
	run := st.Stories[it.ID]
	switch {
	case run == nil:
		return nil, refused("flai serve has started no agent for %s, so there is none to restart", it.ID)
	case run.live() && Alive(run.PID):
		return nil, refused("%s's agent is running (pid %d, started %s)", it.ID, run.PID, run.Started)
	case run.Outcome == OutcomeAsked:
		return nil, refused("%s's agent is waiting for an answer to %s; answering it starts the agent again", it.ID, run.Thread)
	}
	if (it.Agent == nil || it.Agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return nil, refused("%s names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])", it.ID)
	}
	if it.Status == workitem.Ready {
		stories, free, err := readyStories(e.Root)
		if err != nil {
			return nil, err
		}
		reserved := 0
		for _, s := range stories {
			if r := st.Stories[s.ID]; r.live() && s.ID != it.ID {
				reserved++
			}
		}
		if free >= 0 && reserved >= free {
			return nil, refused("the in-progress limit leaves no room for %s", it.ID)
		}
	}
	why := "ended"
	if run.Why != "" {
		why = run.Why
	} else if run.live() {
		why = "dropped: its process is gone"
	}
	return StartNow(ctx, o, e, it.ID, why)
}

// lockFile holds path, made for the purpose, until the returned func is
// called: processes that write the same state take turns. A lock older than
// staleLock is the leftover of a process that died holding it, and is taken.
func lockFile(path string) (unlock func()) {
	deadline := time.Now().Add(5 * time.Second)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			_ = f.Close()
			return func() { _ = os.Remove(path) }
		}
		if info, serr := os.Stat(path); serr == nil && time.Since(info.ModTime()) > staleLock {
			_ = os.Remove(path)
			continue
		}
		if !errors.Is(err, os.ErrExist) || time.Now().After(deadline) {
			// Write without it rather than not at all: the state is advice to
			// the dashboard, and a lost race costs one run's record.
			return func() {}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

const staleLock = 10 * time.Second
