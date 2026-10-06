package guard

import (
	"fmt"
	"slices"
	"strings"
)

// The shared paths are the operator's (S-0295, ADR-0096): claims.shared in
// the manifest decides which overlaps hold a ready story, so an agent that
// could add to it could free its own story. Every session flai serve starts,
// a story's agent, the planner, the orchestrator, and the analyzer, and every
// sub-agent, is refused shared_paths_edit and flai shared add and remove,
// whatever its role or permissions. Listing and checking them are reads, open
// to all. The operator's own session, which flai serve did not start, changes
// them.

// StartedByEnv is the variable flai serve sets, to StartedByServe, in each
// session it starts.
const (
	StartedByEnv   = "FLAI_STARTED_BY"
	StartedByServe = "flai-serve"
)

// SharedPathsEdit is the MCP tool that changes claims.shared, the
// operator's alone.
const SharedPathsEdit = "shared_paths_edit"

// sharedEdits are the subcommands of flai shared that change the list.
var sharedEdits = []string{"add", "remove"}

// operatorsList says why no agent changes the shared paths (ADR-0096).
const operatorsList = "the shared paths, claims.shared in system-flow.yaml, decide which overlaps hold a story, so an agent that changed them could free its own; the list is the operator's to change (ADR-0096)"

// askForShared ends the refusal of a session's own call: what it may do
// instead.
const askForShared = "Ask the operator on a thread with thread_open, naming the pattern and what it would free; listing and checking the shared paths, with shared_paths or flai shared list and check, are open to you."

// sharing are the rules that find a change of the shared paths on a command
// line: flai shared add and remove, and no git command.
var sharing = rules{
	flai: func(cmd, sub string, _ []string) (string, string) {
		if cmd == "shared" && slices.Contains(sharedEdits, sub) {
			return operatorsList, ""
		}
		return "", ""
	},
}

// served says whether flai serve started the session: Served says so, and so
// do a role or a story, which only flai serve sets.
func (g Guard) served() bool { return g.Served || g.Role != "" || g.Story != "" }

// shared decides on a call that changes the shared paths, made by a sub-agent
// or in a session flai serve started: the zero Refusal when the call changes
// nothing there or the operator's own session makes it. The orchestrator's
// own refusal names the call, so that its activity document logs it.
func (g Guard) shared(e Event) Refusal {
	if e.AgentID == "" && !g.served() {
		return Refusal{}
	}
	var call, what string
	switch e.ToolName {
	case MCPPrefix + SharedPathsEdit:
		call, what = SharedPathsEdit, "call "+SharedPathsEdit
	case "Bash":
		for _, words := range commands(e.ToolInput.Command) {
			if why, _ := g.refuse(words, sharing); why != "" {
				call = strings.Join(words, " ")
				what = fmt.Sprintf("run %q", call)
				break
			}
		}
	}
	if call == "" {
		return Refusal{}
	}
	if e.AgentID != "" {
		who := e.AgentType
		if who == "" {
			who = "unnamed"
		}
		return Refusal{Why: fmt.Sprintf("a sub-agent (%s) cannot %s: %s. Put the change you need in your final message.", who, what, operatorsList)}
	}
	r := Refusal{Why: fmt.Sprintf("%s cannot %s: %s. %s", g.who(), what, operatorsList, askForShared)}
	if g.Role == RoleOrchestrate {
		r.Call = call
	}
	return r
}

// who names the session flai serve started, in a refusal.
func (g Guard) who() string {
	switch {
	case g.Role == RolePlan:
		return "the planner"
	case g.Role == RoleOrchestrate:
		return "the orchestrator"
	case g.Role == RoleAnalyze:
		return "the analyzer"
	case g.Role != "":
		return "a session flai serve started with the role " + g.Role
	case g.Story != "":
		return "the story's agent"
	}
	return "a session flai serve started"
}
