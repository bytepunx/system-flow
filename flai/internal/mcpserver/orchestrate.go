package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator's deterministic operations (S-0217): the ready column's
// order by a policy, the backlog stories that could go to ready, and whether
// the release policy is met. Each is a read, answered as its command's
// --json prints it, so flai guard passes it to a sub-agent and the planner.
// The one write beside them, release_publish, publishes by that policy
// (S-0222), and flai guard gives it to the orchestrator alone.

const orderByPolicyDescription = "The ready column's order by a policy, as flai order --by <policy> --json prints it (S-0217), with the figure each story was ordered by. It computes the order and writes nothing; applying it is flai order --by <policy> --apply on the host. cod orders by cost of delay value per week, highest first; wsjf by that value over the forecast duration in hours, highest first; throughput by forecast duration, shortest first; fifo by created, oldest first. A story missing the figure its policy needs goes after those with it, in its current order, and names what it lacks in missing; ties keep the current order. policy is one of cod, wsjf, throughput, or fifo; without it, the project's orchestration.policy in system-flow.yaml, fifo when that is not set. order is the board's whole pull order as it stands."

const promoteCandidatesDescription = "The backlog stories that could go to ready, as flai promote --candidates --json prints them (S-0217), and every other backlog story with each reason it cannot. A backlog story is a candidate when it is not a draft; meets the definition of ready (a goal, acceptance criteria with a checkbox, and an epic that is not cancelled); would not be held if it were ready (it declares touches that overlap no story in progress or in review, and every story it names in after is done); and has a forecast duration and a cost of delay value. The candidates are ordered by the project's orchestration.policy, fifo when it is not set, as order_by_policy orders the ready column, each with its figure. limit caps the candidates listed, 0 or none for all; those beyond it are not listed. It writes nothing: moving a candidate to ready is item_move."

const releaseEvaluateDescription = "Whether the project's release policy, orchestration.release in system-flow.yaml, is met, with its figures, as flai release --evaluate --json prints it (S-0217). It weighs the stories accepted and not yet released (pending): their count and their cost of delay values summed per week, naming those with no value in unvalued. threshold is met when that summed value is at or over the manifest's value, or the count at or over its count. theme is met when every story of its epic or tag, cancelled ones left out, is accepted and one at least is not yet released; not_accepted names the rest. judgement, the default when no policy is set, is never met by itself: the call is the orchestrator's or the operator's. reason says why in one sentence. It tags, bumps, and pushes nothing: releasing is the operator's."

const releasePublishDescription = "Publish what is accepted and not yet released, by the project's release policy (S-0222), as the dashboard's Publish does (publish.run): flai release --pending on the host, under the push host action, never forced, which applies, tags, and pushes everything accepted and unreleased. The host's journal records the run with the orchestrator as who asked and reason as why, apart from the operator's. reason is one sentence, which you also log: the figure that was met, or under judgement why the unreleased work is coherent and complete; under threshold or theme it is the evaluation's reason when left out. It weighs the batch first, as release_evaluate does, and refuses, changing nothing, saying why and what would allow it: while nothing is accepted and not yet released; while orchestration.release.whole_epics holds the batch back, naming each story held and its epic; under threshold or theme, while the policy is not met, with its figures; under judgement, without a reason; and while the push host action is off, which publish.run needs (flai serve enable push). When flai release --pending refuses with exit 3, because the remote has release tags newer than this clone's (S-0174) or its branch moved, the refusal is flai's message, unchanged, beginning conflict:. It returns the policy and the evaluation published by, with its figures, the versions and tags released, and the items bundled. The orchestrator's alone, while orchestration.permissions.publish is on."

// OrderByPolicyIn names the policy to order the ready column by.
type OrderByPolicyIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Policy  string `json:"policy,omitempty" jsonschema:"cod, wsjf, throughput, or fifo; the project's orchestration.policy when left out"`
}

func (in OrderByPolicyIn) project() string { return in.Project }

func (s *server) orderByPolicy(ctx context.Context, _ *mcp.CallToolRequest, in OrderByPolicyIn) (*mcp.CallToolResult, workitem.ReadyOrder, error) {
	policy := in.Policy
	if policy == "" {
		policy = s.repo.Manifest.Orchestration.PolicyOrDefault()
	}
	if !slices.Contains(workitem.OrderPolicies, policy) {
		return nil, workitem.ReadyOrder{}, fmt.Errorf("unknown order policy %q: use one of %s", policy, strings.Join(workitem.OrderPolicies, ", "))
	}
	defer perf.Track(ctx, "order.by")()
	got, _, _, err := s.repo.ReadyOrderByPolicy(policy)
	if err != nil {
		return nil, workitem.ReadyOrder{}, err
	}
	return nil, got, nil
}

// PromoteCandidatesIn caps the candidates listed.
type PromoteCandidatesIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Limit   int    `json:"limit,omitempty" jsonschema:"list at most this many candidates; 0 or left out lists them all"`
}

func (in PromoteCandidatesIn) project() string { return in.Project }

func (s *server) promoteCandidates(ctx context.Context, _ *mcp.CallToolRequest, in PromoteCandidatesIn) (*mcp.CallToolResult, workitem.Candidates, error) {
	if in.Limit < 0 {
		return nil, workitem.Candidates{}, fmt.Errorf("limit is a number of candidates, 0 for all; got %d", in.Limit)
	}
	defer perf.Track(ctx, "promote.candidates")()
	got, err := s.repo.PromotionCandidates(in.Limit)
	if err != nil {
		return nil, workitem.Candidates{}, err
	}
	return nil, got, nil
}

// ReleaseEvaluateIn takes nothing but the project.
type ReleaseEvaluateIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
}

func (in ReleaseEvaluateIn) project() string { return in.Project }

func (s *server) releaseEvaluate(ctx context.Context, _ *mcp.CallToolRequest, _ ReleaseEvaluateIn) (*mcp.CallToolResult, release.Evaluation, error) {
	// what is released is read from git tags
	if s.runner == nil {
		return nil, release.Evaluation{}, fmt.Errorf("this flai mcp reads no git history, so it cannot say what is released; on the host run flai release --evaluate")
	}
	r := execx.Timed(ctx, s.runner)
	if err := execx.Require(r, "git", "What is released is read from git tags; install git."); err != nil {
		return nil, release.Evaluation{}, err
	}
	defer perf.Track(ctx, "release.evaluate")()
	ev, err := release.EvaluateRepo(r, s.repo.Root, s.repo.Manifest, s.repo)
	if err != nil {
		return nil, release.Evaluation{}, err
	}
	return nil, ev, nil
}

// ReleasePublishIn says why the orchestrator publishes.
type ReleasePublishIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Reason  string `json:"reason,omitempty" jsonschema:"one sentence, which you also log: the policy figure that was met, or under judgement why the unreleased work is coherent and complete; needed under judgement"`
}

func (in ReleasePublishIn) project() string { return in.Project }

// ReleasePublished is the release release_publish cut, with the policy and
// the figures it was cut by.
type ReleasePublished struct {
	Policy     string             `json:"policy"`
	Reason     string             `json:"reason" jsonschema:"why it was published: the reason given, else the evaluation's"`
	Evaluation release.Evaluation `json:"evaluation" jsonschema:"the release policy's figures it was published by, as release_evaluate gives them"`
	Versions   []ReleasedVersion  `json:"versions" jsonschema:"each component released, with its version before and after and its items"`
	Tags       []string           `json:"tags" jsonschema:"the release tags made"`
	Items      []string           `json:"items" jsonschema:"the IDs of the items bundled, each once"`
	Pushed     bool               `json:"pushed" jsonschema:"false when the branch tracks no remote: the release is tagged here only"`
	Published  []string           `json:"published,omitempty" jsonschema:"a template component published to its own remote"`
	Warnings   []string           `json:"warnings"`
}

// ReleasedVersion is one component's release in a publish.
type ReleasedVersion struct {
	Component string                `json:"component"`
	Level     string                `json:"level"`
	From      string                `json:"from"`
	To        string                `json:"to"`
	Items     []release.PendingItem `json:"items"`
}

// releasePublish publishes the batch release_evaluate weighs when the
// release policy allows it (S-0222). Every refusal comes before anything
// changes, but flai release --pending's own, which is a conflict.
func (s *server) releasePublish(ctx context.Context, _ *mcp.CallToolRequest, in ReleasePublishIn) (*mcp.CallToolResult, ReleasePublished, error) {
	// flai guard gives the tool to the orchestrator alone, while publish is
	// on; a session it does not stand in front of is held the same, as
	// item_move holds the orchestrator's move to ready (S-0219)
	if !s.orchestrator() {
		return nil, ReleasePublished{}, fmt.Errorf("release_publish is the orchestrator's, while orchestration.permissions.publish is on; any other agent publishes only when the operator asks, with flai release --pending on the host (ADR-0067)")
	}
	if err := s.repo.OrchestratorPermits(manifest.PermitPublish, "publishes", "the most recently accepted story"); err != nil {
		return nil, ReleasePublished{}, err
	}
	if s.runner == nil {
		return nil, ReleasePublished{}, fmt.Errorf("this flai mcp reads no git history, so it cannot say what is released; on the host run flai release --evaluate")
	}
	if s.publish == nil {
		return nil, ReleasePublished{}, fmt.Errorf("this flai mcp cannot publish: publishing is the operator's, with flai release --pending on the host")
	}
	r := execx.Timed(ctx, s.runner)
	if err := execx.Require(r, "git", "What is released is read from git tags; install git."); err != nil {
		return nil, ReleasePublished{}, err
	}
	done := perf.Track(ctx, "release.evaluate")
	ev, err := release.EvaluateRepo(r, s.repo.Root, s.repo.Manifest, s.repo)
	done()
	if err != nil {
		return nil, ReleasePublished{}, err
	}
	reason := strings.TrimSpace(in.Reason)
	if err := mayPublish(ev, reason); err != nil {
		return nil, ReleasePublished{}, err
	}
	if reason == "" {
		reason = ev.Reason
	}
	p := channel.Project{Key: s.key, Name: s.repo.Manifest.Name, Root: projectRoot(s.repo)}
	w, e := s.publish.Publish(ctx, p, "mcp.release_publish", s.orchestratorBy(), ev.Policy+": "+reason)
	switch {
	case e == nil:
	case e.Code == hostapi.Conflict:
		// flai's own message, which reading its answer took conflict: off
		return nil, ReleasePublished{}, errors.New("conflict: " + e.Message)
	case e.Code == hostapi.Disabled:
		return nil, ReleasePublished{}, fmt.Errorf("release_publish publishes as publish.run does, which needs the push host action: %s", e.Message)
	default:
		return nil, ReleasePublished{}, errors.New(e.Message)
	}
	return released(ev, reason, w)
}

// mayPublish says why the batch ev weighs is not published for reason, or
// nil when it may be.
func mayPublish(ev release.Evaluation, reason string) error {
	if ev.Count == 0 {
		return fmt.Errorf("nothing is accepted and not yet released, so there is nothing to publish; it publishes after the next acceptance")
	}
	if len(ev.HeldByEpic) > 0 {
		held := make([]string, len(ev.HeldByEpic))
		for i, h := range ev.HeldByEpic {
			state := "not found"
			if h.EpicStatus != "" {
				state = h.EpicStatus
			}
			held[i] = fmt.Sprintf("%s, whose epic %s is %s", h.ID, h.Epic, state)
		}
		return fmt.Errorf("orchestration.release.whole_epics holds the batch back while a story's epic is in neither review nor done: %s. It publishes once each of those epics is in review or done", strings.Join(held, "; "))
	}
	switch {
	case ev.Policy == manifest.ReleaseJudgement && reason == "":
		return fmt.Errorf("the release policy is judgement, so publishing is your call and needs your reason: give reason, one sentence saying why the unreleased work is coherent and complete, which you also log (%s)", ev.Reason)
	case ev.Policy != manifest.ReleaseJudgement && !ev.Met:
		return fmt.Errorf("the release policy %s is not met: %s. It publishes once it is; release_evaluate gives the figures", ev.Policy, ev.Reason)
	}
	return nil
}

// released reads what flai release --pending printed as the release, or
// says why it published nothing.
func released(ev release.Evaluation, reason string, w hostapi.Written) (*mcp.CallToolResult, ReleasePublished, error) {
	var said struct {
		Plans []struct {
			Component struct {
				Name string `json:"name"`
			} `json:"component"`
			Level string                `json:"level"`
			From  string                `json:"from"`
			To    string                `json:"to"`
			Items []release.PendingItem `json:"items"`
		} `json:"plans"`
		Tags      []string `json:"tags"`
		Pushed    bool     `json:"pushed"`
		PushError string   `json:"push_error"`
		Reason    string   `json:"reason"`
		Published []string `json:"published"`
	}
	if err := json.Unmarshal(w.Data, &said); err != nil {
		return nil, ReleasePublished{}, fmt.Errorf("flai release --pending answered what is not a release: %w", err)
	}
	if said.PushError != "" {
		return nil, ReleasePublished{}, fmt.Errorf("the release is applied and tagged here, but the push failed and nothing was forced: %s. The operator puts that right and publishes again; what is tagged is not tagged again", said.PushError)
	}
	if len(said.Plans) == 0 {
		why := said.Reason
		if why == "" {
			why = "nothing was pending"
		}
		return nil, ReleasePublished{}, fmt.Errorf("flai release --pending published nothing: %s", why)
	}
	out := ReleasePublished{Policy: ev.Policy, Reason: reason, Evaluation: ev, Versions: []ReleasedVersion{}, Tags: said.Tags, Items: []string{}, Pushed: said.Pushed, Published: said.Published, Warnings: w.Warnings}
	if out.Tags == nil {
		out.Tags = []string{}
	}
	if out.Warnings == nil {
		out.Warnings = []string{}
	}
	for _, p := range said.Plans {
		out.Versions = append(out.Versions, ReleasedVersion{Component: p.Component.Name, Level: p.Level, From: p.From, To: p.To, Items: p.Items})
		for _, it := range p.Items {
			if !slices.Contains(out.Items, it.ID) {
				out.Items = append(out.Items, it.ID)
			}
		}
	}
	return nil, out, nil
}

// orchestratorBy names the orchestrator as who publishes, with the name this
// flai mcp knows it by when that does not already, as flai plan's journal
// entry does.
func (s *server) orchestratorBy() string {
	switch {
	case s.agent == "":
		return workitem.ActivityOrchestrator
	case strings.HasPrefix(s.agent, workitem.ActivityOrchestrator):
		return s.agent
	}
	return fmt.Sprintf("%s (%s)", workitem.ActivityOrchestrator, s.agent)
}

// addOrchestrateTools registers the orchestrator's reads, and its publish.
func addOrchestrateTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "order_by_policy", Description: orderByPolicyDescription}, route(p, (*server).orderByPolicy))
	mcp.AddTool(srv, &mcp.Tool{Name: "promote_candidates", Description: promoteCandidatesDescription}, route(p, (*server).promoteCandidates))
	mcp.AddTool(srv, &mcp.Tool{Name: "release_evaluate", Description: releaseEvaluateDescription}, route(p, (*server).releaseEvaluate))
	mcp.AddTool(srv, &mcp.Tool{Name: "release_publish", Description: releasePublishDescription}, route(p, (*server).releasePublish))
}
