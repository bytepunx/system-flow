package serve

import (
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
