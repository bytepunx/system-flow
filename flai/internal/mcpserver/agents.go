package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting a story's agent on the host (S-0177, ADR-0064): an operator's own
// agent can have flai do what flai serve agent start and restart do, under
// the same agent host action, so that a story begun on another host gets an
// agent here without the dashboard. flai guard refuses both to a sub-agent,
// as it refuses every tool that is not a read (ADR-0060).

// AgentStart starts or restarts story's agent in the project at root, as
// flai serve agent start or restart does, on by's word; verb is "start" or
// "restart". A refusal is an error that says why.
type AgentStart func(ctx context.Context, verb, root, story, by string) (AgentStarted, error)

// AgentStarted is the agent started, or the one queued.
type AgentStarted struct {
	Story   string `json:"story"`
	Agent   string `json:"agent"`
	Harness string `json:"harness,omitempty"`
	Command string `json:"command,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Log     string `json:"log,omitempty" jsonschema:"the agent's output, on the host"`
	Started string `json:"started,omitempty"`
	Queued  string `json:"queued,omitempty" jsonschema:"set when no agent started yet: flai serve starts it when the in-progress limit has room and nothing holds the story"`
}

// AgentIn names the story whose agent to start.
type AgentIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story" jsonschema:"story ID such as S-0177 (any zero padding)"`
}

func (in AgentIn) project() string { return in.Project }

const (
	agentStartDescription   = "Start a ready story's agent on this host now, as flai serve agent start does, whatever flai serve's own rules say about when: past a hold or a full in-progress limit, with a warning (S-0115). Needs the agent host action on for the project (flai serve enable agent), as the dashboard's Start agent does. Refused, saying why, when the story is not in ready, while its agent runs or waits for an answer, or when nothing can start it. For a story in progress, use agent_restart."
	agentRestartDescription = "Start a new agent on this host for a story in ready or in progress whose agent dropped or failed, or that this host has had no agent for, as flai serve agent restart does (S-0116, S-0177, ADR-0064). A story in progress begun on another host gets an agent told who began it, where, and when, and to go on from what is committed. A story in ready with the in-progress limit full, or held, is queued instead. Needs the agent host action on for the project (flai serve enable agent), as the dashboard's Retry and Start agent do. Refused, saying why, while this host's agent for it runs or waits for an answer, when one is already queued, or when nothing can start it."
)

func (s *server) agentStart(ctx context.Context, _ *mcp.CallToolRequest, in AgentIn) (*mcp.CallToolResult, AgentStarted, error) {
	return s.startAgent(ctx, "start", in)
}

func (s *server) agentRestart(ctx context.Context, _ *mcp.CallToolRequest, in AgentIn) (*mcp.CallToolResult, AgentStarted, error) {
	return s.startAgent(ctx, "restart", in)
}

func (s *server) startAgent(ctx context.Context, verb string, in AgentIn) (*mcp.CallToolResult, AgentStarted, error) {
	story := workitem.CanonicalID(in.Story)
	if workitem.TypeOfID(story) != workitem.Story {
		return nil, AgentStarted{}, fmt.Errorf("%q is not a story's ID", in.Story)
	}
	if s.agents == nil {
		return nil, AgentStarted{}, fmt.Errorf("this flai mcp cannot start agents; on the host run flai serve agent %s %s", verb, story)
	}
	out, err := s.agents(ctx, verb, projectRoot(s.repo), story, s.agent)
	if err != nil {
		return nil, AgentStarted{}, err
	}
	return nil, out, nil
}
