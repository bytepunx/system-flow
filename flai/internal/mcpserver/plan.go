package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting the planner on the host (S-0208, ADR-0075): an operator's own
// agent can have flai do what flai plan does, under the same plan host
// action. flai guard refuses the tool to a sub-agent, as it refuses every
// tool that is not a read (ADR-0060), and to the planner, whose MCPPlans
// leave it out.

// PlanStart starts the planner for item, an epic or a story, in the project
// at root, as flai plan does, on by's word. A refusal is an error that says
// why.
type PlanStart func(ctx context.Context, root, item, by string) (PlanStarted, error)

// PlanStarted is the planner run started.
type PlanStarted struct {
	Item    string `json:"item"`
	Agent   string `json:"agent"`
	Harness string `json:"harness,omitempty"`
	Command string `json:"command,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Log     string `json:"log,omitempty" jsonschema:"the planner's output, on the host"`
	Session string `json:"session,omitempty"`
	Started string `json:"started,omitempty"`
}

// PlanIn names the epic or the story to plan.
type PlanIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"epic or story ID such as E-0016 or S-0208 (any zero padding)"`
}

func (in PlanIn) project() string { return in.Project }

const planDescription = "Start the planner for an epic or a story on this host now, as flai plan does (S-0208, ADR-0075): in the project's main checkout, with the project's planning agent (planning.agent over agent in system-flow.yaml), it drafts an epic's stories, enriches a story with touches, a forecast, and a cost of delay, and for an epic with stories revisits each one not done or cancelled, writing through flai and moving nothing past backlog. It does not wait for the planner: the run is recorded where flai serve tracks agents, and flai serve settles it once it ends. For the operator's own agent: refused to an agent flai serve started, but for the orchestrator, which may ask for an epic that flai plan --candidates lists (in the backlog with no stories, or open with every story done or cancelled and one done, and no planner running for it or awaiting the operator) while orchestration.permissions.plan_backlog_epics is on, and whose run records orchestrator, in place of asked, as what started it. Needs the plan host action on for the project (flai serve enable plan), as the dashboard's Plan does. Refused, saying why, for a task or an ID that is neither an epic's nor a story's, for an item that is archived, done, or cancelled, while a planner runs for the item, and when nothing can start it."

func (s *server) plan(ctx context.Context, _ *mcp.CallToolRequest, in PlanIn) (*mcp.CallToolResult, PlanStarted, error) {
	item := workitem.CanonicalID(in.ID)
	switch workitem.TypeOfID(item) {
	case workitem.Epic, workitem.Story:
	case workitem.Task:
		return nil, PlanStarted{}, fmt.Errorf("%s is a task; the planner plans an epic or a story", item)
	default:
		return nil, PlanStarted{}, fmt.Errorf("%q is not an epic's or a story's ID; the planner plans an E-nnnn or an S-nnnn", in.ID)
	}
	if s.plans == nil {
		return nil, PlanStarted{}, fmt.Errorf("this flai mcp cannot start the planner; on the host run flai plan %s", item)
	}
	out, err := s.plans(ctx, projectRoot(s.repo), item, s.agent)
	if err != nil {
		return nil, PlanStarted{}, err
	}
	return nil, out, nil
}
