package guard

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"time"
)

// ClaudeCode is the harness name of Claude Code, whose hooks run flai guard
// in a project made from the template, and the default of --harness.
const ClaudeCode = "claude-code"

// ClaudeCodeRefuses is the exit status with which a Claude Code hook refuses
// a tool call; the hook's standard error goes back to the caller.
const ClaudeCodeRefuses = 2

// The Claude Code hook events on a sub-agent, which the guard tells from a
// PreToolUse by hook_event_name. An input that names neither is a tool call,
// whether it names PreToolUse or, as the guard was first fed, no event.
const (
	HookSubagentStart = "SubagentStart"
	HookSubagentStop  = "SubagentStop"
)

// MCPPrefix starts the name of each of flai's MCP tools as Claude Code
// knows them, from a server named flai.
const MCPPrefix = "mcp__flai__"

// fileWrites are Claude Code's tools that write a file, each a FileWrite
// call: the planner and the orchestrator never use them, the analyzer only
// on its report (S-0223), and a sub-agent on a file in a .claude/ folder only
// while auto-approve is on.
var fileWrites = []string{"Edit", "MultiEdit", "NotebookEdit", "Write"}

// Event is the part of a Claude Code hook's input the guard reads: a
// PreToolUse's tool call, or a SubagentStart's or SubagentStop's sub-agent,
// as HookEventName says, in the session SessionID.
type Event struct {
	HookEventName string `json:"hook_event_name"`
	SessionID     string `json:"session_id"`
	ToolName      string `json:"tool_name"`
	ToolInput     Input  `json:"tool_input"`
	// AgentID is set only when a sub-agent makes the call, or on the
	// sub-agent a SubagentStart or SubagentStop reports; AgentType names the
	// sub-agent's definition. A session started with --agent may carry an
	// agent type of its own, so the ID is what marks a sub-agent.
	AgentID   string `json:"agent_id"`
	AgentType string `json:"agent_type"`
}

// Input is the part of a tool call's input the guard reads: Bash's command;
// the file Edit or Write changes, or the notebook NotebookEdit does; and a
// flai tool's Args. UnmarshalJSON fills Args' Fields and CostOfDelay.
type Input struct {
	Command      string `json:"command"`
	FilePath     string `json:"file_path"`
	NotebookPath string `json:"notebook_path"`
	Args
}

// UnmarshalJSON decodes a tool call's input, the names of its fields, and
// the keys of its cost_of_delay. A cost_of_delay that is not an object is
// left nil rather than failing the input, which flai guard would let through
// unread.
func (in *Input) UnmarshalJSON(data []byte) error {
	type input Input
	if err := json.Unmarshal(data, (*input)(in)); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	in.Fields = slices.Sorted(maps.Keys(fields))
	if raw, ok := fields["cost_of_delay"]; ok && json.Unmarshal(raw, &in.CostOfDelay) != nil {
		in.CostOfDelay = nil
	}
	return nil
}

// Refusal is the guard's decision on a Claude Code call: Why it is refused,
// "" when it is not. For the orchestrator's own calls Call is the call
// refused, in one line, and Needs the permission that would allow it, as
// orchestration.permissions.<name>, or "" when none would.
type Refusal struct {
	Why   string
	Call  string
	Needs string
}

// Check says why a Claude Code call is refused, or "" when it is not.
func (g Guard) Check(e Event) string { return g.Decide(e).Why }

// Decide decides on a Claude Code call, as Decide does on the Call it
// makes: the zero Refusal lets it through.
func (g Guard) Decide(e Event) Refusal {
	v := Decide(g, g.call(e))
	if v.Outcome == Allow {
		return Refusal{}
	}
	return Refusal{Why: v.Reason, Call: v.Call, Needs: v.Needs}
}

// call is the Call e makes, in g's role and story: Bash runs a command line,
// fileWrites write their file_path or notebook_path, and a tool of
// MCPPrefix is a flai tool; a call with an agent ID is a sub-agent's.
func (g Guard) call(e Event) Call {
	c := Call{Session: e.SessionID, Subagent: e.AgentID != "", Agent: e.AgentType, Role: g.Role, Story: g.Story, Name: e.ToolName}
	tool, flai := strings.CutPrefix(e.ToolName, MCPPrefix)
	switch {
	case e.ToolName == "Bash":
		c.Kind, c.Command = Shell, e.ToolInput.Command
	case slices.Contains(fileWrites, e.ToolName):
		c.Kind, c.Path = FileWrite, cmp.Or(e.ToolInput.FilePath, e.ToolInput.NotebookPath)
	case flai:
		c.Kind, c.Tool, c.Args = FlaiTool, tool, e.ToolInput.Args
	}
	return c
}

// subAgentChange is the change e reports to the record of running
// sub-agents, and whether it reports one: a SubagentStart's or
// SubagentStop's sub-agent in e's session.
func subAgentChange(e Event) (SubAgentChange, bool) {
	if e.HookEventName != HookSubagentStart && e.HookEventName != HookSubagentStop {
		return SubAgentChange{}, false
	}
	return SubAgentChange{Session: e.SessionID, ID: e.AgentID, Type: e.AgentType, Stopped: e.HookEventName == HookSubagentStop}, true
}

// Record notes in dir, the project's .flai-cache/guard, that e's sub-agent
// started (SubagentStart) or stopped (SubagentStop) at at, as RecordChange
// does. An event of another kind is ignored.
func Record(dir string, e Event, at time.Time) error {
	ch, ok := subAgentChange(e)
	if !ok {
		return nil
	}
	return RecordChange(dir, ch, at)
}

// claudeCode is Claude Code's Reader.
type claudeCode struct{}

// Read reads a Claude Code hook's input: a SubagentStart or SubagentStop is
// a change to the record of running sub-agents, and any other a tool call,
// in g's role and story.
func (claudeCode) Read(data []byte, g Guard) (Hook, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return Hook{}, err
	}
	if ch, ok := subAgentChange(e); ok {
		return Hook{SubAgent: &ch}, nil
	}
	return Hook{Call: g.call(e)}, nil
}

// Answer lets an allowed call through with exit 0, and refuses any other, an
// Ask among them, with ClaudeCodeRefuses and its reason on stderr, which
// Claude Code hands back to the caller.
func (claudeCode) Answer(v Verdict, stderr io.Writer) int {
	if v.Outcome == Allow {
		return 0
	}
	fmt.Fprintln(stderr, v.Reason)
	return ClaudeCodeRefuses
}
