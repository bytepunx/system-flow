package mcpserver

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// Starting the analyzer on the host (S-0223): an operator's own agent can
// have flai do what flai analyze does, under the same analyze host action.
// flai guard refuses the tool to a sub-agent, as it refuses every tool that
// is not a read (ADR-0060), and to the planner, the orchestrator, and the
// analyzer, whose tools leave it out.

// AnalyzeStart starts the analyzer in the project at root, looking for
// focus, one of harness.Focuses, or for all of them when focus is empty, as
// flai analyze does, on by's word. A refusal is an error that says why.
type AnalyzeStart func(ctx context.Context, root, focus, by string) (AnalyzeStarted, error)

// AnalyzeStarted is the analyzer run started.
type AnalyzeStarted struct {
	Focus   string `json:"focus" jsonschema:"what the run looks for: bottlenecks, intent, risk, or all"`
	Agent   string `json:"agent"`
	Harness string `json:"harness,omitempty"`
	Command string `json:"command,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Log     string `json:"log,omitempty" jsonschema:"the analyzer's output, on the host"`
	Session string `json:"session,omitempty"`
	Started string `json:"started,omitempty"`
	Trigger string `json:"trigger,omitempty" jsonschema:"what started the run: asked, for the operator's asking"`
}

// AnalyzeIn names what the analyzer looks for.
type AnalyzeIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Focus   string `json:"focus,omitempty" jsonschema:"bottlenecks, intent, or risk; all three when not given"`
}

func (in AnalyzeIn) project() string { return in.Project }

const analyzeDescription = "Start the analyzer for the project on this host now, as flai analyze does (S-0223): in the project's main checkout, with the project's analysis agent (analysis.agent over agent in system-flow.yaml), it reads the metrics, the design, the code, and the issues, looks for what focus names (bottlenecks in the flow, gaps between design/system and the code as intent, or technical and security risks as risk), or for all three without one, and writes one report under design/analysis, editing nothing else. It does not wait for the analyzer: the run is recorded where flai serve tracks agents, and flai serve settles it once it ends and logs the report it wrote and its cost. For the operator's own agent: refused to an agent flai serve started. Needs the analyze host action on for the project (flai serve enable analyze), as the dashboard's Analyze does. Refused, saying why, for a focus other than bottlenecks, intent, or risk, while an analyzer runs for the project (one runs at a time), and when nothing can start it."

func (s *server) analyze(ctx context.Context, _ *mcp.CallToolRequest, in AnalyzeIn) (*mcp.CallToolResult, AnalyzeStarted, error) {
	focus := strings.TrimSpace(in.Focus)
	if focus != "" && !slices.Contains(harness.Focuses, focus) {
		return nil, AnalyzeStarted{}, fmt.Errorf("the analyzer takes the focus %s, or none for all of them, and %q is none of them", strings.Join(harness.Focuses, ", "), in.Focus)
	}
	if s.analyses == nil {
		how := "flai analyze"
		if focus != "" {
			how += " --focus " + focus
		}
		return nil, AnalyzeStarted{}, fmt.Errorf("this flai mcp cannot start the analyzer; on the host run %s", how)
	}
	out, err := s.analyses(ctx, projectRoot(s.repo), focus, s.agent)
	if err != nil {
		return nil, AnalyzeStarted{}, err
	}
	return nil, out, nil
}
