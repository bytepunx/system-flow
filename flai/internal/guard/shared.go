package guard

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// The shared paths are the operator's (S-0295, ADR-0096): claims.shared in
// the manifest decides which overlaps hold a ready story, so an agent that
// could add to it could free its own story. Every session flai serve starts,
// a story's agent, the planner, the orchestrator, the analyzer, and the
// curator, and every sub-agent, is refused shared_paths_edit and flai shared
// add and remove, whatever its role or permissions. Listing and checking them
// are reads, open to all. The operator's own session, which flai serve did
// not start, changes them.

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

// served says whether flai serve started c's session: Served says so, and so
// do a role or a story, which only flai serve sets.
func (g Guard) served(c Call) bool { return g.Served || c.Role != "" || c.Story != "" }

// shared decides on a call that changes the shared paths, made by a sub-agent
// or in a session flai serve started: the zero Verdict when the call changes
// nothing there or the operator's own session makes it. The orchestrator's
// own refusal names the call, so that its activity document logs it.
func (g Guard) shared(c Call) Verdict {
	if !c.Subagent && !g.served(c) {
		return Verdict{}
	}
	var call, what string
	switch {
	case c.Kind == FlaiTool && c.Tool == SharedPathsEdit:
		call, what = SharedPathsEdit, "call "+SharedPathsEdit
	case c.Kind == Shell:
		for _, words := range commands(c.Command) {
			if why, _ := g.refuse(words, sharing); why != "" {
				call = strings.Join(words, " ")
				what = fmt.Sprintf("run %q", call)
				break
			}
		}
	}
	if call == "" {
		return Verdict{}
	}
	if c.Subagent {
		return deny(fmt.Sprintf("a sub-agent (%s) cannot %s: %s. Put the change you need in your final message.", cmp.Or(c.Agent, "unnamed"), what, operatorsList))
	}
	v := deny(fmt.Sprintf("%s cannot %s: %s. %s", g.who(c), what, operatorsList, askForShared))
	if c.Role == RoleOrchestrate {
		v.Call = call
	}
	return v
}

// who names the session flai serve started, in a refusal of c.
func (g Guard) who(c Call) string {
	switch {
	case c.Role == RolePlan:
		return "the planner"
	case c.Role == RoleOrchestrate:
		return "the orchestrator"
	case c.Role == RoleAnalyze:
		return "the analyzer"
	case c.Role == RoleCurate:
		return "the curator"
	case DefinedRole(c.Role) && g.RoleFile != nil && g.RoleFile.Name != "":
		return "the role " + g.RoleFile.Name
	case c.Role != "":
		return "a session flai serve started with the role " + c.Role
	case c.Story != "":
		return "the story's agent"
	}
	return "a session flai serve started"
}
