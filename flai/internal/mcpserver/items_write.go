package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ItemNewIn creates a work item (S-0103).
type ItemNewIn struct {
	Project string          `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Type    string          `json:"type" jsonschema:"story, task, or epic"`
	Title   string          `json:"title"`
	Nature  string          `json:"nature,omitempty" jsonschema:"feature (the default), improvement, remediation, research, or experiment"`
	Parent  string          `json:"parent,omitempty" jsonschema:"a story's epic (optional), a task's story (required)"`
	Tags    []string        `json:"tags,omitempty"`
	Touches []string        `json:"touches,omitempty" jsonschema:"paths or components the work changes"`
	Topics  []string        `json:"topics,omitempty" jsonschema:"a story's or epic's: what it is about beyond the components its tags and touches reach, such as logging or release"`
	After   []string        `json:"after,omitempty" jsonschema:"what it waits for until they are done: a story's stories (it is held in ready meanwhile), a task's tasks of the same story; flai check runs with it, and an entry that does not exist, a task of another story, or a cycle refuses the creation"`
	Agent   *manifest.Agent `json:"agent,omitempty" jsonschema:"a story's agent: harness, model, config, and roles (explore, verify: each a harness, model, and config for that sub-agent), over the project's default, which fills in what is not given, role by role"`
	Body    string          `json:"body,omitempty" jsonschema:"the goal, criteria, and notes below the heading; the template's empty sections when not given"`
	Draft   bool            `json:"draft,omitempty" jsonschema:"a story's only: true makes it a draft, which cannot go to ready until the operator finalizes it; give it for a story you wrote for the operator to review"`
}

func (in ItemNewIn) project() string { return in.Project }

// ItemNewOut is the item created, and the stories in progress whose claims
// a new task grew its story's claim into.
type ItemNewOut struct {
	ItemOut
	Overlaps []itemedit.Overlapping `json:"overlaps,omitempty" jsonschema:"a new task's: the other stories in progress whose claim covers a path its touches added to its story's claim (I-0059); the task stands, and both stories are told as an overlapped change: coordinate with their agents before you change those paths"`
}

func (s *server) itemNew(ctx context.Context, _ *mcp.CallToolRequest, in ItemNewIn) (*mcp.CallToolResult, ItemNewOut, error) {
	nature := in.Nature
	if nature == "" {
		nature = "feature"
	}
	owner := s.repo.Manifest.Owner
	if owner == "" {
		owner = s.agent
	}
	opt := workitem.NewOptions{Type: in.Type, Title: strings.TrimSpace(in.Title), Nature: nature, Parent: in.Parent, Owner: owner,
		Tags: in.Tags, Touches: in.Touches, Topics: in.Topics, After: in.After, Agent: in.Agent, Body: in.Body, Draft: in.Draft, Now: s.now()}
	var watch *itemedit.ClaimWatch
	if in.Type == workitem.Task && in.Parent != "" {
		watch = itemedit.WatchClaim(s.repo, in.Parent)
	}
	var it *workitem.Item
	var err error
	if len(in.After) > 0 || strings.TrimSpace(in.Body) != "" {
		// an after: entry that names nothing, or forms a cycle, is the check's
		// to find, so a creation that sets one is checked, as the CLI's --after
		// is (S-0176). A body its author wrote is checked as --body-stdin is: a
		// missing section or a lint finding refuses it and leaves nothing
		// (S-0209). The template's own body is not, as flai story new's is not.
		var res *itemnew.Result
		if res, err = itemnew.Create(s.repo, s.runner, itemnew.Options{New: opt}); err == nil {
			it = res.Item
		}
	} else {
		it, err = s.repo.Create(opt)
	}
	if err != nil {
		return nil, ItemNewOut{}, refusal(err)
	}
	out, err := s.itemOut(ctx, it)
	return nil, ItemNewOut{ItemOut: out, Overlaps: s.grown(watch, it.ID)}, err
}

// grown is what the write to id grew its story's claim into, told to both
// stories (I-0059). A failure is logged, and the write stands: the report is
// advisory.
func (s *server) grown(w *itemedit.ClaimWatch, id string) []itemedit.Overlapping {
	told, err := w.Grown(s.agent, s.now())
	if err != nil && s.logger != nil {
		s.logger.Warn("overlap notices not sent", "component", "mcp", "agent", s.agent, "item", id, "err", err.Error())
	}
	return told
}

// ItemEditIn changes an item's own words (S-0103): only what is given changes.
type ItemEditIn struct {
	Project string    `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string    `json:"id"`
	Hash    string    `json:"hash,omitempty" jsonschema:"the hash item_get gave; a change made meanwhile is then refused instead of overwritten"`
	Title   *string   `json:"title,omitempty"`
	Nature  *string   `json:"nature,omitempty"`
	Tags    *[]string `json:"tags,omitempty" jsonschema:"replaces the tags; an empty list removes them"`
	Touches *[]string `json:"touches,omitempty" jsonschema:"replaces the touches; an empty list removes them"`
	Topics  *[]string `json:"topics,omitempty" jsonschema:"a story's or epic's: replaces what it is about, such as logging or release; an empty list removes them"`
	After   *[]string `json:"after,omitempty" jsonschema:"a story's or a task's: replaces what it waits for until they are done, a story's stories (it is held in ready meanwhile), a task's tasks of the same story; an empty list removes them"`
	Parent  *string   `json:"parent,omitempty"`
	// Agent replaces a story's agent; ClearAgent removes it.
	Agent      *manifest.Agent `json:"agent,omitempty" jsonschema:"replaces the story's agent with exactly this harness, model, config, and roles; roles left out are removed"`
	ClearAgent bool            `json:"clear_agent,omitempty" jsonschema:"removes the story's agent"`
	Body       *string         `json:"body,omitempty" jsonschema:"replaces everything below the heading"`
	// Draft, CostOfDelay, and Forecast are the item's planning data (S-0199),
	// each block stamped with this agent as who set it.
	Draft            *bool          `json:"draft,omitempty" jsonschema:"a story's: true makes a story in the backlog a draft; false is refused, since finalizing a draft is the operator's, but to the orchestrator, which finalizes a complete draft with false alone"`
	CostOfDelay      *CostOfDelayIn `json:"cost_of_delay,omitempty" jsonschema:"a story's or epic's cost of delay: only the keys given change; with clear_cost_of_delay it replaces the cost of delay"`
	ClearCostOfDelay bool           `json:"clear_cost_of_delay,omitempty" jsonschema:"removes the cost of delay; to remove one amount, give this and cost_of_delay with the keys to keep"`
	Forecast         *ForecastIn    `json:"forecast,omitempty" jsonschema:"a story's forecast: only the keys given change; with clear_forecast it replaces the forecast"`
	ClearForecast    bool           `json:"clear_forecast,omitempty" jsonschema:"removes the story's forecast"`
}

func (in ItemEditIn) project() string { return in.Project }

// CostOfDelayIn is the keys of a cost of delay an edit sets. An amount cannot
// be removed alone, since a JSON number has no empty value: clear the cost of
// delay and give the keys to keep.
type CostOfDelayIn struct {
	RevenuePerWeek   *float64 `json:"revenue_per_week,omitempty" jsonschema:"input: the revenue the item brings each week once it is done, in the project's currency"`
	PenaltyPerWeek   *float64 `json:"penalty_per_week,omitempty" jsonschema:"input: what each week it is not done costs beyond revenue, in the project's currency"`
	TimeLostPerCycle *string  `json:"time_lost_per_cycle,omitempty" jsonschema:"input: the work lost each cycle it is not done, a Go duration such as 4h; an empty string removes it"`
	Value            *float64 `json:"value,omitempty" jsonschema:"the cost of delay per week, in the project's currency"`
}

// ForecastIn is the keys of a forecast an edit sets: nil leaves a key, an
// empty string removes it.
type ForecastIn struct {
	Duration *string `json:"duration,omitempty" jsonschema:"the agent time the story is expected to take, a Go duration such as 6h; an empty string removes it"`
	Delivery *string `json:"delivery,omitempty" jsonschema:"when it is expected to be done, a UTC timestamp such as 2026-10-09T17:00:00Z; an empty string removes it"`
	Basis    *string `json:"basis,omitempty" jsonschema:"what the forecast rests on, in one sentence; an empty string removes it"`
}

// edit is the itemedit change the keys given make.
func (in *CostOfDelayIn) edit() *itemedit.CostOfDelayEdit {
	if in == nil {
		return nil
	}
	return &itemedit.CostOfDelayEdit{RevenuePerWeek: amount(in.RevenuePerWeek), PenaltyPerWeek: amount(in.PenaltyPerWeek), TimeLostPerCycle: in.TimeLostPerCycle, Value: amount(in.Value)}
}

func (in *ForecastIn) edit() *itemedit.ForecastEdit {
	if in == nil {
		return nil
	}
	return &itemedit.ForecastEdit{Duration: in.Duration, Delivery: in.Delivery, Basis: in.Basis}
}

// amount is v as itemedit reads an amount, nil when it is not given.
func amount(v *float64) *string {
	if v == nil {
		return nil
	}
	s := strconv.FormatFloat(*v, 'f', -1, 64)
	return &s
}

// ItemEditOut is what an edit changed.
type ItemEditOut struct {
	ID        string   `json:"id"`
	Changed   []string `json:"changed" jsonschema:"title, nature, tags, topics, touches, after, agent, parent, draft, cost_of_delay, forecast, goal, criteria, notes, body"`
	Unchanged bool     `json:"unchanged,omitempty"`
	Hash      string   `json:"hash"`
	// Overlaps are the stories in progress whose claims the edit grew its
	// story's claim into (I-0059).
	Overlaps []itemedit.Overlapping `json:"overlaps,omitempty" jsonschema:"an edit of touches: the other stories in progress whose claim covers a path it added to its story's claim (I-0059); the edit stands, and both stories are told as an overlapped change: coordinate with their agents before you change those paths"`
}

func (s *server) itemEdit(_ context.Context, _ *mcp.CallToolRequest, in ItemEditIn) (*mcp.CallToolResult, ItemEditOut, error) {
	ch := itemedit.Change{Title: in.Title, Nature: in.Nature, Tags: in.Tags, Topics: in.Topics, Touches: in.Touches, After: in.After, Parent: in.Parent, Body: in.Body, Agent: in.Agent, ClearAgent: in.ClearAgent,
		Draft: in.Draft, CostOfDelay: in.CostOfDelay.edit(), ClearCostOfDelay: in.ClearCostOfDelay, Forecast: in.Forecast.edit(), ClearForecast: in.ClearForecast}
	if ch == (itemedit.Change{}) {
		return nil, ItemEditOut{}, errors.New("nothing to change: give title, nature, tags, topics, touches, after, parent, agent, clear_agent, body, draft, cost_of_delay, clear_cost_of_delay, forecast, or clear_forecast")
	}
	if in.Draft != nil && !*in.Draft {
		if !s.orchestrator() {
			// finalizing is the operator's, and the orchestrator's (S-0219)
			return nil, ItemEditOut{}, fmt.Errorf("%s: draft false would finalize it, which is the operator's: say in a thread or your narrative that it is ready to be finalized", in.ID)
		}
		if err := s.finalizes(in.ID, ch); err != nil {
			return nil, ItemEditOut{}, err
		}
	}
	var watch *itemedit.ClaimWatch
	if in.Touches != nil {
		watch = itemedit.WatchClaim(s.repo, in.ID)
	}
	res, err := itemedit.Apply(s.repo, s.runner, in.ID, ch, itemedit.Options{Hash: in.Hash, By: s.agent, NoCommit: true, Now: s.now()})
	if err != nil {
		return nil, ItemEditOut{}, refusal(err)
	}
	changed := res.Changed
	if changed == nil {
		changed = []string{}
	}
	return nil, ItemEditOut{ID: res.ID, Changed: changed, Unchanged: res.Unchanged, Hash: res.Hash, Overlaps: s.grown(watch, res.ID)}, nil
}

// finalizes says why the orchestrator may not make change ch to item id, a
// draft false, or nil when it may (S-0219): the change finalizes the draft
// and does nothing else, and the draft check (flai promote --drafts) finds
// the story complete. The finalized block then names the orchestrator, as
// any edit's names who made it.
func (s *server) finalizes(id string, ch itemedit.Change) error {
	if ch != (itemedit.Change{Draft: ch.Draft}) {
		return fmt.Errorf("%s: the orchestrator finalizes a draft with draft false and nothing else; it changes nothing else of an item", id)
	}
	it, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	lacks, err := s.repo.Finalizable(it)
	if err != nil {
		return err
	}
	if len(lacks) > 0 {
		return fmt.Errorf("%s is not complete, so the orchestrator does not finalize it (flai promote --drafts): %s", it.ID, strings.Join(lacks, "; "))
	}
	return nil
}

// refusal is err with, when flai check refused the change, the findings in
// it: they are what the agent fixes, not only how many there are.
func refusal(err error) error {
	r, ok := docedit.IsRefused(err)
	if !ok {
		return err
	}
	msgs := make([]string, len(r.Findings))
	for i, f := range r.Findings {
		msgs[i] = fmt.Sprintf("%s:%d: %s: %s", f.Path, f.Line, f.Rule, f.Message)
	}
	return fmt.Errorf("%w: %s", err, strings.Join(msgs, "; "))
}
