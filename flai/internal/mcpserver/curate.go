package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Starting the curator on the host (S-0431, ADR-0196): an operator's own
// agent can have flai do what flai curate does, under the same curate host
// action. flai guard refuses the tool to a sub-agent, as it refuses every
// tool that is not a read (ADR-0060), and to the planner, the orchestrator,
// the analyzer, and the curator, whose tools leave it out.

// CurateStart starts the curator in the project at root, taking in what the
// archive and the threads gained since since, a YYYY-MM-DD day in UTC or an
// RFC 3339 time, or since its last run when since is empty, as flai curate
// does, on by's word. A refusal is an error that says why.
type CurateStart func(ctx context.Context, root, since, by string) (CurateStarted, error)

// CurateStarted is the curator run started.
type CurateStarted struct {
	Since   string `json:"since" jsonschema:"what the run takes in from: a day or a time, or empty for everything on the curator's first run"`
	Agent   string `json:"agent"`
	Harness string `json:"harness,omitempty"`
	Command string `json:"command,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Log     string `json:"log,omitempty" jsonschema:"the curator's output, on the host"`
	Session string `json:"session,omitempty"`
	Started string `json:"started,omitempty"`
	Trigger string `json:"trigger,omitempty" jsonschema:"what started the run: asked, for the operator's asking"`
}

// CurateIn names what the curator takes in from.
type CurateIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Since   string `json:"since,omitempty" jsonschema:"take in what was gained since this YYYY-MM-DD day in UTC or RFC 3339 time; since the curator's last run, or everything on its first, when not given"`
}

func (in CurateIn) project() string { return in.Project }

const curateDescription = "Start the curator for the project on this host now, as flai curate does (S-0431, ADR-0196): in the project's main checkout, with the project's digest agent (digest.agent over agent in system-flow.yaml), it reads what the archive and the threads gained since since (a YYYY-MM-DD day in UTC or an RFC 3339 time), or without one since its last run, or everything on its first, and takes it in as digest entries written through flai digest new and edit alone, editing no file and authoring no work item. It does not wait for the curator: the run is recorded where flai serve tracks agents, and flai serve settles it once it ends and logs the entries it wrote or changed and its cost in wip/agents/curator.md. For the operator's own agent: refused to an agent flai serve started. Needs the curate host action on for the project (flai serve enable curate). Refused, saying why, for a since in neither form, while a curator runs for the project (one runs at a time), and when nothing can start it."

func (s *server) curate(ctx context.Context, _ *mcp.CallToolRequest, in CurateIn) (*mcp.CallToolResult, CurateStarted, error) {
	since := strings.TrimSpace(in.Since)
	if s.curations == nil {
		how := "flai curate"
		if since != "" {
			how += " --since " + since
		}
		return nil, CurateStarted{}, fmt.Errorf("this flai mcp cannot start the curator; on the host run %s", how)
	}
	out, err := s.curations(ctx, projectRoot(s.repo), since, s.agent)
	if err != nil {
		return nil, CurateStarted{}, err
	}
	return nil, out, nil
}
