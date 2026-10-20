package mcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const queueAddDescription = "The orchestrator's enqueue of a story in review in the merge queue, while orchestration.permissions.accept_reviews is on (S-0403, ADR-0093): in place of accepting a story it has verified, the orchestrator enqueues it, and the queue lands it as the orchestrator with this evidence. It records with the story who enqueued it (orchestrator) and when, verified, the commit the orchestrator's verifier checked, and evidence, the verifier's report as flai accept --by orchestrator takes it: one Verdict: line and one item \"- <n>: <files>\" per acceptance criterion, naming the changed files that meet it, with no heading. It is refused, and nothing changes, to anyone but the orchestrator; while accept_reviews is off; without verified or evidence; for evidence with a heading in it or that cannot be read so; for a story that flai queue add refuses (not in review, not eligible, naming every reason, or under acceptance.mode pull-request); and, as flai accept --by orchestrator refuses its acceptance, unless verified is the head of the story's branch, every criterion is ticked, every changed file is under the story's touches, no thread on the story or its tasks is open, the evidence names a changed file for every criterion, and the branch changes no path Claude Code protects, naming each condition that fails by its code. A story enqueued already keeps its place and what it was enqueued with, and changed is false. Answers the change as flai queue add --json prints it: action, story, changed, by, at, and the queue's state after it. flai queue add <story> --by orchestrator --verified <commit> --evidence <file> on the command line enqueues the same way."

// QueueAddIn names the story the orchestrator enqueues and what it vouches
// for it with.
type QueueAddIn struct {
	Project  string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story    string `json:"story" jsonschema:"the story in review to enqueue (S-nnnn, any zero padding)"`
	Verified string `json:"verified" jsonschema:"the commit the orchestrator's verifier checked: the head of the story's branch"`
	Evidence string `json:"evidence" jsonschema:"the evidence, as flai accept --by orchestrator takes it: a Verdict: line and one - <n>: <files> item per acceptance criterion, without headings"`
}

func (in QueueAddIn) project() string { return in.Project }

// queueAdd enqueues the story the call names as the orchestrator, with the
// commit its verifier checked and its evidence, as flai queue add --by
// orchestrator does.
func (s *server) queueAdd(ctx context.Context, _ *mcp.CallToolRequest, in QueueAddIn) (*mcp.CallToolResult, mergequeue.Change, error) {
	// flai guard gives the tool to the orchestrator; a session it does not
	// stand in front of is held the same, as release_publish is (S-0222)
	if !s.orchestrator() {
		return nil, mergequeue.Change{}, fmt.Errorf("queue_add is the orchestrator's, while orchestration.permissions.accept_reviews is on; anyone else enqueues a story with flai queue add %s on the host", in.Story)
	}
	// the queue's eligibility and the evidence's check read git
	if s.runner == nil {
		return nil, mergequeue.Change{}, fmt.Errorf("this flai mcp reads no git history, so it cannot check %s's branch against the evidence; on the host run flai queue add %s --by %s --verified <commit> --evidence <file>", in.Story, in.Story, workitem.ActivityOrchestrator)
	}
	ch, err := EnqueueByOrchestrator(execx.Timed(ctx, s.runner), s.repo, in.Story, in.Verified, in.Evidence, s.now(), s.logger)
	if err != nil {
		return nil, mergequeue.Change{}, err
	}
	return nil, ch, nil
}

// EnqueueByOrchestrator enqueues story in repo's merge queue as the
// orchestrator at at, reading git through r, with verified, the commit its
// verifier checked, and evidence, its text as given (ADR-0093, S-0403), so
// that the queue lands the story as the orchestrator with that evidence.
// flai queue add --by orchestrator and the MCP tool queue_add both enqueue
// through it. It is refused, changing nothing, while
// orchestration.permissions.accept_reviews is off; without verified or
// evidence; for evidence with a heading in it, or that preview.ParseEvidence
// cannot read; for a story mergequeue.Add refuses; and with a
// *mergequeue.Refusal naming each condition by its code, for a story the
// orchestrator's acceptance would be refused with that commit and evidence,
// as flai accept --by orchestrator checks them. log takes the acceptance
// preview's warnings; nil logs nothing.
func EnqueueByOrchestrator(r execx.Runner, repo *workitem.Repo, story, verified, evidence string, at time.Time, log *slog.Logger) (mergequeue.Change, error) {
	id := workitem.CanonicalID(story)
	if err := repo.OrchestratorPermits(manifest.PermitAcceptReviews, "enqueues a story for the merge queue to land", id); err != nil {
		return mergequeue.Change{}, fmt.Errorf("rule: %w", err)
	}
	verified, evidence = strings.TrimSpace(verified), strings.TrimSpace(evidence)
	if verified == "" || evidence == "" {
		return mergequeue.Change{}, fmt.Errorf("rule: the orchestrator enqueues %s only with the commit its verifier checked, the head of story/%s, and its evidence, a Verdict: line and one - <n>: <files> item per criterion: give both, --verified and --evidence to flai queue add, verified and evidence to queue_add (ADR-0093)", id, id)
	}
	for i, line := range strings.Split(evidence, "\n") {
		if strings.HasPrefix(line, "#") {
			return mergequeue.Change{}, fmt.Errorf("rule: line %d of the orchestrator's evidence is a heading (%s): write a Verdict: line and - <n>: <files> items without headings, since the queue's landing puts the evidence under ### Accepted by the orchestrator in the story's Notes", i+1, line)
		}
	}
	ev, err := preview.ParseEvidence(evidence)
	if err != nil {
		return mergequeue.Change{}, fmt.Errorf("rule: the orchestrator's evidence for %s cannot be read as ADR-0093 asks: %w", id, err)
	}
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	vouch := func(it *workitem.Item) error {
		res, err := preview.AcceptWith(r, repo, it, workitem.ActivityOrchestrator, at, log, preview.AcceptOptions{Orchestrator: true, Verified: verified, Evidence: ev})
		if err != nil {
			return fmt.Errorf("check the orchestrator's evidence for %s before enqueuing it: %w", it.ID, err)
		}
		if len(res.OrchestratorBlockers) == 0 {
			return nil
		}
		reasons := make([]mergequeue.Reason, 0, len(res.OrchestratorBlockers))
		for _, b := range res.OrchestratorBlockers {
			reasons = append(reasons, mergequeue.Reason{Kind: b.Code, Text: b.Message})
		}
		return &mergequeue.Refusal{Story: it.ID, Text: fmt.Sprintf("%s is not enqueued: the orchestrator's acceptance of it with this commit and evidence would be refused, so the queue could not land it as the orchestrator (ADR-0093)", it.ID), Reasons: reasons}
	}
	return mergequeue.AddWith(r, repo, story, workitem.ActivityOrchestrator, at, mergequeue.AddOptions{Verified: verified, Evidence: evidence, Vouch: vouch})
}

// addQueueTools registers the orchestrator's enqueue in the merge queue.
func addQueueTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "queue_add", Description: queueAddDescription}, route(p, (*server).queueAdd))
}
