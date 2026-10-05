package mcpserver

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator's deterministic operations (S-0217): the ready column's
// order by a policy, the backlog stories that could go to ready, and whether
// the release policy is met. Each is a read, answered as its command's
// --json prints it, so flai guard passes it to a sub-agent and the planner.

const orderByPolicyDescription = "The ready column's order by a policy, as flai order --by <policy> --json prints it (S-0217), with the figure each story was ordered by. It computes the order and writes nothing; applying it is flai order --by <policy> --apply on the host. cod orders by cost of delay value per week, highest first; wsjf by that value over the forecast duration in hours, highest first; throughput by forecast duration, shortest first; fifo by created, oldest first. A story missing the figure its policy needs goes after those with it, in its current order, and names what it lacks in missing; ties keep the current order. policy is one of cod, wsjf, throughput, or fifo; without it, the project's orchestration.policy in system-flow.yaml, fifo when that is not set. order is the board's whole pull order as it stands."

const promoteCandidatesDescription = "The backlog stories that could go to ready, as flai promote --candidates --json prints them (S-0217), and every other backlog story with each reason it cannot. A backlog story is a candidate when it is not a draft; meets the definition of ready (a goal, acceptance criteria with a checkbox, and an epic that is not cancelled); would not be held if it were ready (it declares touches that overlap no story in progress or in review, and every story it names in after is done); and has a forecast duration and a cost of delay value. The candidates are ordered by the project's orchestration.policy, fifo when it is not set, as order_by_policy orders the ready column, each with its figure. limit caps the candidates listed, 0 or none for all; those beyond it are not listed. It writes nothing: moving a candidate to ready is item_move."

const releaseEvaluateDescription = "Whether the project's release policy, orchestration.release in system-flow.yaml, is met, with its figures, as flai release --evaluate --json prints it (S-0217). It weighs the stories accepted and not yet released (pending): their count and their cost of delay values summed per week, naming those with no value in unvalued. threshold is met when that summed value is at or over the manifest's value, or the count at or over its count. theme is met when every story of its epic or tag, cancelled ones left out, is accepted and one at least is not yet released; not_accepted names the rest. judgement, the default when no policy is set, is never met by itself: the call is the orchestrator's or the operator's. reason says why in one sentence. It tags, bumps, and pushes nothing: releasing is the operator's."

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

// addOrchestrateTools registers the orchestrator's reads.
func addOrchestrateTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "order_by_policy", Description: orderByPolicyDescription}, route(p, (*server).orderByPolicy))
	mcp.AddTool(srv, &mcp.Tool{Name: "promote_candidates", Description: promoteCandidatesDescription}, route(p, (*server).promoteCandidates))
	mcp.AddTool(srv, &mcp.Tool{Name: "release_evaluate", Description: releaseEvaluateDescription}, route(p, (*server).releaseEvaluate))
}
