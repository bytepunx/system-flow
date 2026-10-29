package serve

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// runConcerns is what run's agent event is on, and its role (S-0532): a
// story's agent's run is on its story, role story; a planner's on the item it
// plans, role plan; the orchestrator's, the analyzer's, and the curator's on
// no item, each by its role's name, roles orchestrate, analyze, and curate;
// and a project role's on no item, by the role it is marked with, the role
// its name. typ is inbox.RoleType for a run on no item, and empty for an
// item, whose type its ID tells.
func runConcerns(run *AgentRun) (id, typ, role string) {
	switch {
	case run.curates():
		return workitem.ActivityCurator, inbox.RoleType, conventions.RoleCurate
	case run.Story != "":
		return run.Story, "", conventions.RoleStory
	case run.Item != "":
		return run.Item, "", conventions.RolePlan
	case run.analyzes():
		return workitem.ActivityAnalyzer, inbox.RoleType, conventions.RoleAnalyze
	case run.orchestrates():
		return workitem.ActivityOrchestrator, inbox.RoleType, conventions.RoleOrchestrate
	}
	return run.Role, inbox.RoleType, run.Role
}

// agentEvent is the agent event of run that came to to for why, on the item
// or the role runConcerns says, the item's title read from repo. An item
// that cannot be read is named by its ID alone, of the type its ID tells:
// the event is news, and a title is not worth losing it for.
func agentEvent(repo *workitem.Repo, run *AgentRun, to, why string) inbox.LoggedEvent {
	id, typ, role := runConcerns(run)
	title := id
	if typ == "" {
		typ, title = workitem.TypeOfID(id), ""
		if it, err := repo.Get(id); err == nil {
			id, typ, title = it.ID, it.Type, it.Title
		}
	}
	return inbox.AgentEvent(id, typ, title, to, run.Agent, role, why)
}

// appendAgentEvent appends to the event log of the project at root the agent
// event of run that came to to, one of inbox.AgentStarted, AgentEnded,
// AgentFailed, and AgentStopped, for why, in the words of flai serve's
// journal (S-0532). It is the one way flai serve tells an agent's start, end,
// failure, or stop to the agents that read the inbox. A project it cannot
// open, or an event it cannot append, is told to warn and the run goes on:
// the event log is news, not state (ADR-0211 §4).
func appendAgentEvent(root string, run *AgentRun, to, why string, warn func(msg string, args ...any)) {
	repo, err := workitem.Open(root)
	if err != nil {
		warn("agent event not logged: the project could not be opened", "root", root, "agent", run.Agent, "to", to, "err", err)
		return
	}
	e := agentEvent(repo, run, to, why)
	if err := inbox.AppendEvent(repo, e); err != nil {
		warn("agent event not logged", "id", e.ID, "agent", run.Agent, "to", to, "err", err)
	}
}

// appendAgentEnd appends to the event log of the project at root the agent
// event of run, which has ended, as endEvent says it: the one call a run path
// makes at a run's end, as it calls appendAgentEvent with
// inbox.AgentStarted at its start (S-0532).
func appendAgentEnd(root string, run *AgentRun, warn func(msg string, args ...any)) {
	to, why := endEvent(run, warn)
	appendAgentEvent(root, run, to, why, warn)
}

// endEvent is what the agent event of run, which has ended, came to, and why
// (S-0532). A run the operator stopped is stopped, for its own why. A run
// whose log ends on the harness's usage limit failed, whatever else it was
// judged, for limitWhy's words, which name the limit and when it resets, or
// that the harness named no reset time, followed by the run's own why when
// it says more, and its exit. Otherwise a run that failed failed, for its own
// why with its exit beside it, and a run that worked or asked ended, for its
// own why, or a phrase that names its outcome when it has none. A log that
// cannot be read is told to warn, and read as ending on no limit.
func endEvent(run *AgentRun, warn func(msg string, args ...any)) (to, why string) {
	if run.Outcome == OutcomeStopped {
		if run.Why == "" {
			return inbox.AgentStopped, "stopped by the operator"
		}
		return inbox.AgentStopped, run.Why
	}
	if lim, _ := runLimit(run, warn); lim != nil {
		why = limitWhy(lim)
		// the curator's run says the limit as its why already
		if run.Why != "" && run.Why != why {
			why += "; " + run.Why
		}
		return inbox.AgentFailed, withExit(why, run.Exit)
	}
	switch run.Outcome {
	case OutcomeFailed:
		return inbox.AgentFailed, withExit(run.Why, run.Exit)
	case OutcomeAsked:
		if run.Why == "" {
			return inbox.AgentEnded, "ended waiting for an answer"
		}
	case OutcomeWorked:
		if run.Why == "" {
			return inbox.AgentEnded, "ended with its work done"
		}
	}
	if run.Why == "" {
		return inbox.AgentEnded, "ended"
	}
	return inbox.AgentEnded, run.Why
}

// withExit is why with exit beside it, in the words of flai serve's
// journal: why itself when it names the exit already, as the journal's
// "ended (exit 1) with S-0001 in in-progress" does, "ended (exit 1)" for no
// why, and why followed by "(exit 1)" otherwise. An exit nobody saw is
// left out of a why, and is the why, "ended (an exit code nobody saw)", when
// there is none.
func withExit(why string, exit *int) string {
	if exit == nil {
		if why == "" {
			return "ended (an exit code nobody saw)"
		}
		return why
	}
	code := fmt.Sprintf("(exit %d)", *exit)
	switch {
	case why == "":
		return "ended " + code
	case strings.Contains(why, code):
		return why
	}
	return why + " " + code
}
