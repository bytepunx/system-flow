package serve

import (
	"context"
	"fmt"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Start starts a ready story's agent now, on the operator's word (S-0115):
// flai serve agent start, or the story page's Start agent button. It is what
// the launcher does when a story enters ready, without the launcher's own
// rules about when: whether the story was ready before flai serve began,
// whether it had an agent since it entered ready, and whoever is attending.
// It is refused while the agent action is off for the project, when the
// story is not in ready, while its agent runs or waits for an answer, and
// when nothing can start it. A full in-progress limit does not refuse it,
// and nor does a hold (S-0128): the operator's word goes past both, as a
// move does, with a warning.
func Start(ctx context.Context, o Options, e Entry, story string) (*AgentRun, error) {
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
	if it.Status != workitem.Ready {
		why := "only a story in ready is started"
		if it.Status == workitem.InProgress {
			why += "; flai serve agent restart starts a new agent for one in progress whose agent dropped or failed"
		}
		return nil, refused("%s is in %s; %s", it.ID, it.Status, why)
	}
	st := o.Dir.AgentStates()[e.Root]
	switch run := st.Stories[it.ID]; {
	case run.live() && Alive(run.PID):
		return nil, refused("%s's agent is running (pid %d, started %s)", it.ID, run.PID, run.Started)
	case run != nil && run.Outcome == OutcomeAsked:
		return nil, refused("%s's agent is waiting for an answer to %s; answering it starts the agent again", it.ID, run.Thread)
	}
	if (it.Agent == nil || it.Agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return nil, refused("%s names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])", it.ID)
	}
	ok, hold, err := roomFor(e, st, it.ID)
	if err != nil {
		return nil, err
	}
	run, err := StartNow(ctx, o, e, it.ID, "")
	if err == nil && !ok && o.Logger != nil {
		o.Logger.Warn("agent started past the in-progress limit", "component", "serve", "story", it.ID,
			"detail", fmt.Sprintf("the in-progress limit was full; %s's agent was started on the operator's word", it.ID))
	}
	if err == nil && hold != nil && o.Logger != nil {
		o.Logger.Warn("agent started for a held story", "component", "serve", "story", it.ID,
			"detail", fmt.Sprintf("%s was %s; its agent was started on the operator's word", it.ID, hold.Reason))
	}
	return run, err
}

// roomFor says whether the in-progress limit leaves room for the ready story
// id, counting an agent running for another story still in ready as a story
// in progress, as the launcher does, and whether a claim holds it.
func roomFor(e Entry, st AgentState, id string) (bool, *workitem.Hold, error) {
	stories, holds, free, err := readyStories(e.Root)
	if err != nil {
		return false, nil, err
	}
	reserved := 0
	var self *workitem.Item
	for _, s := range stories {
		if r := st.Stories[s.ID]; r.live() && s.ID != id {
			reserved++
			holds.Open(s.item, agentStarted)
		}
		if s.ID == id {
			self = s.item
		}
	}
	var hold *workitem.Hold
	if self != nil {
		hold = holds.Of(self)
	}
	return free < 0 || reserved < free, hold, nil
}
