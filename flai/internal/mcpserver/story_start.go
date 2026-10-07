package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/storystart"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const storyStartDescription = "Start a ready story in one call (S-0274), with the answer flai story start --json prints: move it to in-progress as you, its epic following as item_move moves it (followed), open its narrative and its branch, story/<id>, in its worktree as flai stream open does (worktree, an absolute path, branch, and from, where the branch came from; all empty where the project is not a git repository), build its prime pack, and read your inbox as inbox does. Work the story in that worktree. The answer, story, followed, worktree, branch, from, pack, and inbox, is kept to 40,000 bytes of JSON, as prime's parts are (ADR-0104): pack is part 1 of the pack, as prime with the story and part 1 answers it, with part and parts, when the answer fits with it, leaving out the inbox's oldest changes, counted in changes_omitted, where that makes it fit; then call prime with the story and part 2, and so on up to parts. When part 1 does not fit even then, pack is only the pack's header, its story, title, topics, budget, exceeded, size, open_issues, and parts, with no part and no readme, left_out, conventions, items, or catalog: read every part with prime, from part 1. budget is the pack's size, such as 80KB, as prime takes it. A story that is not ready, is not a story, or is held, by an overlap with a story in progress or in review or by after, is refused with why and what clears it, and nothing is changed; one in progress already was pulled by another agent when wait_for_work offered it: call wait_for_work again. A step that fails after the move is an error naming the step and the command that finishes it; the story stays in progress."

// StoryStartIn names the story to start (S-0274).
type StoryStartIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story" jsonschema:"the ready story to start, S-nnnn, in any zero padding"`
	Budget  string `json:"budget,omitempty" jsonschema:"the size the prime pack fits, such as 80KB; default the project's prime.budget, else 80KB"`
}

func (in StoryStartIn) project() string { return in.Project }

// StoryStartOut is storystart's result itself, so that story_start answers
// key for key what flai story start --json prints, with the pack fitted to
// one tool result.
type StoryStartOut = storystart.Result

// storyStart starts the story as the server's agent. A refusal and a step
// that fails after the move are the tool's errors, saying what to do.
func (s *server) storyStart(ctx context.Context, _ *mcp.CallToolRequest, in StoryStartIn) (*mcp.CallToolResult, StoryStartOut, error) {
	res, err := storystart.Start(ctx, storystart.Options{
		Repo: s.repo, Runner: s.runner, Story: in.Story, Agent: s.agent, Session: s.session, Budget: in.Budget,
		Now: s.now, Version: s.version, RelativePaths: s.relativePaths, Log: s.logger,
	})
	var refused *storystart.Refused
	var step *storystart.StepError
	switch {
	case errors.As(err, &refused):
		return nil, StoryStartOut{}, startRefusal(refused)
	case errors.As(err, &step):
		return nil, StoryStartOut{}, stepFailure(step, res)
	case err != nil:
		return nil, StoryStartOut{}, err
	}
	return nil, fit(res), nil
}

// startRefusal is the refusal, with what an agent sent by wait_for_work does
// when another agent pulled the story first.
func startRefusal(r *storystart.Refused) error {
	if r.Status == workitem.InProgress && r.Hold == nil {
		return fmt.Errorf("%w; if wait_for_work offered it, another agent pulled it first: call wait_for_work again", r)
	}
	return r
}

// stepFailure is the failed step, with where the steps before it left the
// story: its worktree, and a pack built but not answered.
func stepFailure(step *storystart.StepError, res storystart.Result) error {
	err := error(step)
	if res.Worktree != "" {
		err = fmt.Errorf("%w; its worktree is %s, on branch %s", err, res.Worktree, res.Branch)
	}
	if res.Pack != nil {
		err = fmt.Errorf("%w; its pack is built: read it with prime and story %s, every part", err, step.Story)
	}
	return err
}

// fit is the answer story_start gives, within ctxpack.PartLimit bytes of
// JSON (ADR-0104). The pack is part 1 as prime cuts it, so that prime's
// parts 2 on follow, with the inbox whole when the answer fits with it,
// else with its oldest changes left out, counted in changes_omitted as the
// inbox's cap counts them. When part 1 does not fit even with no changes,
// the pack is its header alone, without a part number, the inbox whole,
// and every part is read with prime.
func fit(res storystart.Result) storystart.Result {
	res.Pack = res.Pack.Parts(ctxpack.PartLimit)[0]
	if encodedSize(res) <= ctxpack.PartLimit {
		return res
	}
	if res.Inbox != nil {
		whole := res.Inbox
		in := *whole
		res.Inbox = &in
		for len(in.Changes) > 0 && encodedSize(res) > ctxpack.PartLimit {
			in.Changes, in.Omitted = in.Changes[1:], in.Omitted+1
		}
		if encodedSize(res) <= ctxpack.PartLimit {
			return res
		}
		res.Inbox = whole
	}
	head := *res.Pack
	head.PartNumber, head.README = 0, nil
	head.Conventions, head.Omitted, head.Items = []ctxpack.Convention{}, []string{}, []ctxpack.Item{}
	head.Catalog = ctxpack.Catalog{NotLoaded: []ctxpack.Entry{}, InPart: []ctxpack.Entry{}}
	res.Pack = &head
	return res
}

// encodedSize is how many bytes v encodes to as JSON, as the MCP go-sdk
// encodes a tool's result.
func encodedSize(v any) int {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // a start's result holds strings, numbers, and lists, which always encode
	}
	return len(b)
}
