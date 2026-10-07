package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const verifyDescription = "Verify a story before it goes to review (S-0270): run, in the story's worktree, what the close-out runs, cheapest first, stopping at the first step that fails: no rebase is left unfinished (rebase), the branch contains the main branch (sync), the narrative's Current state and Next steps are written (narrative), flai check --strict scoped to the story passes (check), and then the test and lint tiers the branch's changes select, each run with CLOSE_OUT_STORY set to the story. Answers passed, stopped_at naming the step that failed, the commit and the main branch it was verified against, the paths the tiers were selected for, and each step with its state (passed, failed, or not-reached), its duration, and its findings, each a name, a path, a line, and a message, at most max across the run (5 by default; omitted counts the rest), a tier's with its command and exit code. A check finding outside the story is a note, in notes, not a failure. A failing step is an answer, not an error; a story with no worktree, or not a story, is refused. The answer is stored as the story's last, which flai verify --last prints. It commits nothing and records no issue. flai verify --json on the command line answers the same. Call it before you move a story to review, instead of having the verifier sub-agent run the tests, lint, and check; the verifier is kept for the review of the diff against the story's acceptance criteria and the conventions."

// VerifyIn verifies one story in its worktree (S-0270).
type VerifyIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story" jsonschema:"the story to verify (S-nnnn, any zero padding), in its worktree under .flai-cache/worktrees"`
	Max     int    `json:"max,omitempty" jsonschema:"the most findings to answer across the run; 5 when left out"`
}

func (in VerifyIn) project() string { return in.Project }

// verify answers what flai verify --json prints for the story the call
// names, run in its worktree with the tiers its branch's manifest declares.
func (s *server) verify(ctx context.Context, _ *mcp.CallToolRequest, in VerifyIn) (*mcp.CallToolResult, verify.Report, error) {
	if in.Max < 0 {
		return nil, verify.Report{}, fmt.Errorf("max %d is not a number of findings; give 1 or more, or leave it out for %d", in.Max, verify.DefaultMax)
	}
	story, wt, err := s.storyWorktree(in.Story)
	if err != nil {
		return nil, verify.Report{}, err
	}
	tiers, err := verify.CheckoutTiers(wt)
	if err != nil {
		return nil, verify.Report{}, err
	}
	git := s.runner
	if git == nil {
		git = execx.System{}
	}
	rep, err := verify.Verify(ctx, verify.StoryOptions{
		Story: story, Project: s.repo, Worktree: wt, Tiers: tiers, Git: execx.Timed(ctx, git),
		RunOptions: verify.RunOptions{Max: in.Max},
	})
	if err != nil {
		return nil, verify.Report{}, err
	}
	return nil, rep, nil
}

// storyWorktree is the story id names, as its canonical ID, and its
// worktree, which must exist.
func (s *server) storyWorktree(id string) (story, worktree string, err error) {
	if strings.TrimSpace(id) == "" {
		return "", "", errors.New("name the story to verify, such as S-0042")
	}
	it, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, workitem.ErrNotFound) {
			return "", "", fmt.Errorf("%s is not a work item here; name a story with a worktree", id)
		}
		return "", "", err
	}
	if it.Type != workitem.Story {
		return "", "", fmt.Errorf("%s is a %s; name a story (S-nnnn): verify runs for a story's worktree", it.ID, it.Type)
	}
	wt := s.repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", "", fmt.Errorf("%s has no worktree at %s; open it with flai stream open %s on the host", it.ID, s.rel(wt), it.ID)
	}
	return it.ID, wt, nil
}
