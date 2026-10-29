package inbox

import (
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Agent is the kind of the event that an agent flai serve runs started,
// ended, failed, or was stopped (S-0532): ID, Type, and Title are the item the
// run concerns, a story for a story's agent and the epic or story planned for
// a planner, or, for a run on no item, the role's name, of type RoleType; To
// is one of AgentStarted, AgentEnded, AgentFailed, and AgentStopped; By is
// the agent's name; Fields hold FieldRole and FieldWhy. It is kept in the
// event log with the other kinds flai serve logs: the MCP server runs in the
// agent's process and cannot read flai serve's own records.
const Agent = "agent"

// The values of an agent event's To: what became of the run.
const (
	AgentStarted = "started" // flai serve started it
	AgentEnded   = "ended"   // it ended on its own, its work done or waiting for an answer
	AgentFailed  = "failed"  // it could not be started, or ended without its work done
	AgentStopped = "stopped" // the operator stopped it
)

// RoleType is the Type of an agent event of a run on no item, the
// orchestrator's, the analyzer's, the curator's, or a project role's: its ID
// and Title are the role's name.
const RoleType = "role"

// The keys of an agent event's Fields, beside those of the other kinds.
const (
	FieldRole = "role" // agent: the run's kind, conventions.RoleStory, RolePlan, RoleOrchestrate, RoleAnalyze, RoleCurate, or a project role's name
	FieldWhy  = "why"  // agent: why it came to To, in the words flai serve's journal has
)

// AgentEvent is the agent event of the run of the agent by, of kind role, on
// the item id of type typ and title, or on no item, id and title the role's
// name and typ RoleType, that came to to for why. An empty role or why is
// left out of Fields.
func AgentEvent(id, typ, title, to, by, role, why string) LoggedEvent {
	fields := map[string]string{}
	if role != "" {
		fields[FieldRole] = role
	}
	if why != "" {
		fields[FieldWhy] = why
	}
	if len(fields) == 0 {
		fields = nil
	}
	return LoggedEvent{Kind: Agent, ID: id, Type: typ, Title: title, To: to, By: by, Fields: fields}
}

// AgentEventsOf are the agent events in the log on the item or role id, in
// any zero padding, stamped at since or after it, oldest first: what the
// agents of its runs did, for a thread that tells the operator of them.
func AgentEventsOf(repo *workitem.Repo, id string, since time.Time) []LoggedEvent {
	want := workitem.CanonicalID(id)
	since = since.UTC().Truncate(time.Second)
	var out []LoggedEvent
	for _, e := range LoggedEvents(repo) {
		if e.Kind != Agent || workitem.CanonicalID(e.ID) != want {
			continue
		}
		at, err := time.Parse(workitem.TimeFormat, e.At)
		if err != nil || at.Before(since) {
			continue
		}
		out = append(out, e)
	}
	// the log is in the order written, and a writer that stamped its own time
	// may have written after a later one
	sort.SliceStable(out, func(i, j int) bool { return out[i].At < out[j].At })
	return out
}

// describeAgent says an agent event in words: whose run, what became of it,
// and why, when fields say. Without fields it says an item's run is its
// agent's.
func describeAgent(c workitem.Change, fields map[string]string) string {
	var who string
	switch {
	case c.Type == RoleType:
		who = "the " + c.ID + "'s run"
	case fields[FieldRole] == conventions.RolePlan:
		who = c.ID + " " + c.Title + "'s planner"
	default:
		who = c.ID + " " + c.Title + "'s agent"
	}
	if c.By != "" && c.By != c.ID {
		who += " " + c.By
	}
	became := c.To
	switch c.To {
	case AgentStopped:
		became = "was stopped"
	case "":
		became = "changed"
	}
	why := ""
	if w := fields[FieldWhy]; w != "" {
		why = ": " + w
	}
	return who + " " + became + why + letPass
}
