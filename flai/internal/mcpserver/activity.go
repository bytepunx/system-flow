package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A strategic agent's activities (S-0206, ADR-0079): the planner, the
// orchestrator, the analyzer, the curator, or a role the project defines
// (ADR-0207) reports each activity it finishes, and flai measures it from the run's log and appends it to the
// kind's document under wip/agents. flai guard refuses the tool to a
// sub-agent, as it refuses every tool that is not a read (ADR-0060).

// ActivityLog logs an activity of the strategic agent kind in the project at
// root, as by reported it, and returns the entry and the document's totals.
// A refusal is an error that says why.
type ActivityLog func(ctx context.Context, root, kind, summary string, items []string, by string) (ActivityLogged, error)

// ActivityLogged is the activity logged, the totals of the document it was
// logged in, and what its cost was charged to.
type ActivityLogged struct {
	Entry    ActivityEntry  `json:"entry"`
	Activity ActivityTotals `json:"activity"`
	// ChargedTo are the items or issues the activity's cost was charged to,
	// by the kind's rule (ADR-0249): the item its run was started for
	// (ADR-0083), or the work items or the issues its activity named, a
	// share each (ADR-0095, S-0227).
	ChargedTo []string `json:"charged_to,omitempty" jsonschema:"the items or issues the cost was charged to under usage.strategic, each item summed up to its epic"`
	Charge    string   `json:"charge,omitempty" jsonschema:"what the cost was charged to, in a line: the items or issues, or the kind's project strategic total when none"`
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
	Kind    string   `json:"kind" jsonschema:"the strategic agent: planner, orchestrator, analyzer, curator, or the name of a role the project defines under its design folder's workflows/roles"`
	Summary string   `json:"summary" jsonschema:"one line saying what the activity did"`
	Items   []string `json:"items,omitempty" jsonschema:"the IDs of the items the activity touched, such as S-0230 or I-0042 (optional); under the named rule the cost is split evenly between the epics, stories, and tasks named, under the issues rule between the issues named"`
}

func (in ActivityIn) project() string { return in.Project }

const activityLogDescription = "Log an activity of a strategic agent, the planner, the orchestrator, the analyzer, the curator, or a role the project defines under its design folder's workflows/roles (ADR-0207), in its activity document, wip/agents/<kind>.md in the main checkout (S-0206, ADR-0079). Call it at the end of each activity, such as an item planned, the board ordered, an analysis written, or digest entries written, with your kind (a defined role's is its name), a one-line summary of what you did, and the IDs of the items you touched. flai measures the activity's wall-clock seconds and cost from your run's log, from the later of the run's start and the last entry to now, appends the entry, and adds it to the document's totals; the cost is the run's apportioned to the activity, marked estimated when it is. The cost is also charged under usage.strategic, under your kind, by the rule your role file names in charge (ADR-0249; with no role file, the planner's is item, the orchestrator's named, the analyzer's issues, and the curator's none): item, to the item your run was started for, a formula's instance included, summed up to its epic; named, split evenly between the epics, stories, and tasks you named (a thread, an issue, or anything else takes no share), each summed up to its epic; issues, split evenly between the issues you named (a work item or anything else takes no share); none, to nothing. What no item or issue takes is left in the kind's project strategic total. An activity outside a run flai serve logged is logged with no seconds and no cost. A run that ends logs the time since the last entry by itself, so an agent that does one thing and ends need not call it. Returns the entry, the document's totals, and what the cost was charged to. Not for a story's agent, whose usage is written on its story. Nothing is committed."

func (s *server) activityLog(ctx context.Context, _ *mcp.CallToolRequest, in ActivityIn) (*mcp.CallToolResult, ActivityLogged, error) {
	kind := strings.TrimSpace(in.Kind)
	kinds, err := s.repo.ActivityKinds()
	if err != nil {
		return nil, ActivityLogged{}, err
	}
	if !slices.Contains(kinds, kind) {
		return nil, ActivityLogged{}, fmt.Errorf("kind %q is not a strategic agent: give %s", in.Kind, orList(kinds))
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

// orList joins words as a sentence lists them: "a, b, or c".
func orList(words []string) string {
	if len(words) < 2 {
		return strings.Join(words, "")
	}
	return strings.Join(words[:len(words)-1], ", ") + ", or " + words[len(words)-1]
}
