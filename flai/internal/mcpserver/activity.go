package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A strategic agent's activities (S-0206, ADR-0079): the planner, the
// orchestrator, or the analyzer reports each activity it finishes, and flai
// measures it from the run's log and appends it to the kind's document under
// wip/agents. flai guard refuses the tool to a sub-agent, as it refuses every
// tool that is not a read (ADR-0060).

// ActivityLog logs an activity of the strategic agent kind in the project at
// root, as by reported it, and returns the entry and the document's totals.
// A refusal is an error that says why.
type ActivityLog func(ctx context.Context, root, kind, summary string, items []string, by string) (ActivityLogged, error)

// ActivityLogged is the activity logged and the totals of the document it
// was logged in.
type ActivityLogged struct {
	Entry    ActivityEntry  `json:"entry"`
	Activity ActivityTotals `json:"activity"`
}

// ActivityEntry is one activity as the log keeps it.
type ActivityEntry struct {
	At        string   `json:"at" jsonschema:"when the activity ended"`
	Summary   string   `json:"summary"`
	Items     []string `json:"items,omitempty"`
	Seconds   int64    `json:"seconds" jsonschema:"wall-clock seconds, measured from the run's log; 0 outside a run flai serve logged"`
	Cost      float64  `json:"cost" jsonschema:"US dollars, the run's cost apportioned to the activity"`
	Estimated bool     `json:"estimated" jsonschema:"the cost was apportioned or priced rather than reported"`
}

// ActivityTotals are an activity document's front matter.
type ActivityTotals struct {
	Kind           string  `json:"kind"`
	AccruedCost    float64 `json:"accrued_cost" jsonschema:"US dollars, every activity logged"`
	AccruedSeconds int64   `json:"accrued_seconds"`
	TasksCompleted int     `json:"tasks_completed" jsonschema:"the number of activities logged"`
	LastRun        string  `json:"last_run" jsonschema:"when the newest activity ended"`
}

// ActivityIn is the activity a strategic agent reports.
type ActivityIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Kind    string   `json:"kind" jsonschema:"the strategic agent: planner, orchestrator, or analyzer"`
	Summary string   `json:"summary" jsonschema:"one line saying what the activity did"`
	Items   []string `json:"items,omitempty" jsonschema:"the IDs of the items the activity touched, such as S-0230 (optional)"`
}

func (in ActivityIn) project() string { return in.Project }

const activityLogDescription = "Log an activity of a strategic agent, the planner, the orchestrator, or the analyzer, in its activity document, wip/agents/<kind>.md in the main checkout (S-0206, ADR-0079). Call it at the end of each activity, such as an item planned, the board ordered, or an analysis written, with your kind, a one-line summary of what you did, and the IDs of the items you touched. flai measures the activity's wall-clock seconds and cost from your run's log, from the later of the run's start and the last entry to now, appends the entry, and adds it to the document's totals; the cost is the run's apportioned to the activity, marked estimated when it is. An activity outside a run flai serve logged is logged with no seconds and no cost. A run that ends logs the time since the last entry by itself, so an agent that does one thing and ends need not call it. Returns the entry and the document's totals. Not for a story's agent, whose usage is written on its story. Nothing is committed."

func (s *server) activityLog(ctx context.Context, _ *mcp.CallToolRequest, in ActivityIn) (*mcp.CallToolResult, ActivityLogged, error) {
	kind := strings.TrimSpace(in.Kind)
	if !workitem.IsActivityKind(kind) {
		return nil, ActivityLogged{}, fmt.Errorf("kind %q is not a strategic agent: give planner, orchestrator, or analyzer", in.Kind)
	}
	summary := strings.Join(strings.Fields(in.Summary), " ")
	if summary == "" {
		return nil, ActivityLogged{}, errors.New("an activity needs a summary: say in one line what you did")
	}
	if s.activities == nil {
		return nil, ActivityLogged{}, errors.New("this flai mcp cannot log activities; flai serve logs a strategic run's activity by itself when the run ends")
	}
	items := make([]string, 0, len(in.Items))
	for _, it := range in.Items {
		items = append(items, workitem.CanonicalID(strings.TrimSpace(it)))
	}
	out, err := s.activities(ctx, projectRoot(s.repo), kind, summary, items, s.agent)
	if err != nil {
		return nil, ActivityLogged{}, err
	}
	return nil, out, nil
}
