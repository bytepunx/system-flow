package mcpserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/storystart"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const storyStartDescription = "Start a ready story in one call (S-0274), with the answer flai story start --json prints: move it to in-progress as you, its epic following as item_move moves it (followed), open its narrative and its branch, story/<id>, in its worktree as flai stream open does (worktree, an absolute path, branch, and from, where the branch came from; all empty where the project is not a git repository), build its prime pack, and read your inbox as inbox does, save that it lists only the open threads on the story and its tasks and counts the others in threads_omitted: call inbox to list them. Work the story in that worktree. The whole answer, story, followed, worktree, branch, from, pack, and inbox, is kept to 40,000 bytes of JSON, as prime's parts are (ADR-0104, ADR-0179): pack is part 1 of the pack, as prime with the story and part 1 answers it, with part and parts, when the answer fits with it, leaving out the inbox's oldest changes, counted in changes_omitted, where that makes it fit; then call prime with the story and part 2, and so on up to parts. When part 1 does not fit even then, pack is only the pack's header, its story, title, topics, budget, exceeded, size, open_issues, and parts, with no part and no readme, left_out, conventions, items, or catalog: read every part with prime, from part 1; the inbox's oldest changes are left out again where that makes it fit, and where it does not fit even with none, each listed thread's last entry and pending recommendation is shortened, counted in entries_shortened: read them whole with thread_get. budget is the pack's size, such as 80KB, as prime takes it. A story that is not ready, is not a story, or is held, by an overlap with a story in progress or in review or by after, is refused with why and what clears it, and nothing is changed; one in progress already was pulled by another agent when wait_for_work offered it: call wait_for_work again. A step that fails after the move is an error naming the step and the command that finishes it; the story stays in progress."

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
	return nil, storystart.Fit(res), nil
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
