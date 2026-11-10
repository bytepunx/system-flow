package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const acceptDescription = "The orchestrator's acceptance of a story in review, while orchestration.permissions.accept_reviews is on (ADR-0278, ADR-0093): flai accept <story> --by orchestrator --verified <commit> --evidence -, run by flai on the host with the host's configuration, wherever the orchestrator runs, so that it reads the host actions a container's flai cannot. verified is the commit the close-out passed at, as flai verify <story> --last names it after the MCP tool verify; evidence is the verifier's report, one Verdict: line and one item \"- <n>: <files>\" per acceptance criterion. With dry_run it runs flai accept <story> --by orchestrator --dry-run and answers the preview, its blockers and the orchestrator's among them, changing nothing; verified and evidence are optional then. It answers flai accept --json's result as data. It is refused, and nothing changes, to anyone but the orchestrator; while accept_reviews is off; under pull-request acceptance while the forge host action is off; and as flai accept --by orchestrator refuses: unless verified is the head of the story's branch, every criterion is ticked, every changed file is under the story's touches, no thread on the story or its tasks is open, and the evidence names a changed file for every criterion. While the merge queue batches, flai accept refuses the orchestrator's acceptance and names queue_add: enqueue the story with it instead. The orchestrator's alone, in place of flai accept in its shell, which refuses in a role's container."

// AcceptIn names the story the orchestrator accepts, the commit its verifier
// checked, and its evidence, or asks only for the preview.
type AcceptIn struct {
	Project  string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story    string `json:"story" jsonschema:"the story in review to accept (S-nnnn, any zero padding)"`
	Verified string `json:"verified,omitempty" jsonschema:"the commit the close-out passed at, the head of the story's branch, in hex; optional with dry_run"`
	Evidence string `json:"evidence,omitempty" jsonschema:"the evidence, as flai accept --by orchestrator takes it: a Verdict: line and one - <n>: <files> item per acceptance criterion; optional with dry_run"`
	DryRun   bool   `json:"dry_run,omitempty" jsonschema:"preview the acceptance and its blockers, changing nothing"`
}

func (in AcceptIn) project() string { return in.Project }

// Accepted is what flai accept --json answered, and the warnings it logged.
type Accepted struct {
	Story    string         `json:"story" jsonschema:"the story accepted, or previewed"`
	DryRun   bool           `json:"dry_run,omitempty" jsonschema:"true when this was only the preview"`
	Result   map[string]any `json:"result" jsonschema:"flai accept --json's answer: the preview with its blockers, or the acceptance"`
	Warnings []string       `json:"warnings,omitempty" jsonschema:"the warnings flai accept logged"`
}

// accept runs the orchestrator's acceptance on the host, through the
// Accepter flai mcp gives the server (ADR-0278).
func (s *server) accept(ctx context.Context, _ *mcp.CallToolRequest, in AcceptIn) (*mcp.CallToolResult, Accepted, error) {
	id := workitem.CanonicalID(strings.TrimSpace(in.Story))
	// flai guard gives the tool to the orchestrator; a session it does not
	// stand in front of is held the same, as queue_add is
	if !s.orchestrator() {
		return nil, Accepted{}, fmt.Errorf("accept is the orchestrator's, while orchestration.permissions.accept_reviews is on; the operator accepts %s with flai accept %s --by <name> on the host, or from the board", id, id)
	}
	if err := s.repo.OrchestratorPermits(manifest.PermitAcceptReviews, "accepts a story", id); err != nil {
		return nil, Accepted{}, err
	}
	if s.accepter == nil {
		return nil, Accepted{}, fmt.Errorf("this flai mcp cannot run flai accept on the host: on the host run flai accept %s --by %s --verified <commit> --evidence <file>", id, workitem.ActivityOrchestrator)
	}
	p := channel.Project{Key: s.key, Name: s.repo.Manifest.Name, Root: projectRoot(s.repo)}
	w, e := s.accepter.Accept(ctx, p, "mcp.accept", s.orchestratorBy(), id, strings.TrimSpace(in.Verified), in.Evidence, in.DryRun)
	if e != nil {
		return nil, Accepted{}, acceptRefused(e)
	}
	out := Accepted{Story: id, DryRun: in.DryRun, Warnings: w.Warnings}
	if len(w.Data) > 0 {
		if err := json.Unmarshal(w.Data, &out.Result); err != nil {
			return nil, Accepted{}, fmt.Errorf("flai accept %s answered what is not its result: %w", id, err)
		}
	}
	return nil, out, nil
}

// acceptRefused is a refusal of the acceptance as the tool's error: flai's
// message, with the findings or paths a Refused carries, so the orchestrator
// can name each on its thread.
func acceptRefused(e *channel.Error) error {
	msg := e.Message
	switch e.Code {
	case hostapi.Disabled:
		msg = "accept runs the acceptance as accept.run does, which needs the host action: " + msg
	case hostapi.Refused:
		if data, err := json.Marshal(e.Data); err == nil && e.Data != nil {
			msg += " " + string(data)
		}
	}
	return errors.New(msg)
}

// addAcceptTools registers the orchestrator's acceptance on the host.
func addAcceptTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "accept", Description: acceptDescription}, route(p, (*server).accept))
}
