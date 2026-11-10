package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const verifyDescription = "Ask flai serve for a story's close-out (S-0546, ADR-0229): it runs the close-out in the story's worktree as a job, with no model, scripts/close-out.sh S-nnnn -m <message> where the worktree has the script, else flai verify S-nnnn --record-issues, so no tool call holds it and the harness's time limit does not cut it short. It answers at once: the job's id (started true, or running true when a job of the story was running already and is answered instead), or, when the worktree is clean at a commit the verify record holds a passing verdict for, that verdict with recorded true, its commit and ran_at, and nothing runs. message is the commit message the close-out commits with; all leaves no tier to the merge queue; fresh runs whatever the record holds. Then hold wait_for_events until it reports verified for the story, with the job, the commit, the verdict (passed, stopped with stopped_at, failed, or interrupted), and the duration, and read the result, each step with its findings, with flai verify S-nnnn --last; the job's output is in log. Call it before you move a story to review, instead of running the close-out in a turn, and keep the verifier sub-agent for the review of the diff against the story's acceptance criteria and the conventions. When flai serve is not running, it takes no request within ten seconds and the call says so: run the close-out script, or flai verify where the project has none, yourself. A story with no worktree, or an ID that is not a story, is refused."

// VerifyIn asks flai serve for one story's close-out (S-0546).
type VerifyIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story" jsonschema:"the story to close out (S-nnnn, any zero padding), in its worktree under .flai-cache/worktrees"`
	Message string `json:"message,omitempty" jsonschema:"the commit message the close-out commits with, as scripts/close-out.sh -m takes it; needed when the worktree has something uncommitted"`
	All     bool   `json:"all,omitempty" jsonschema:"leave no tier to the merge queue: run every tier the branch's changes select, those marked queue too"`
	Fresh   bool   `json:"fresh,omitempty" jsonschema:"run every step whatever the verify record holds, reusing nothing"`
}

// VerifyOut is flai serve's answer to the close-out request: the job it
// started or found running, or the verdict the verify record holds.
type VerifyOut struct {
	Job     string `json:"job,omitempty"`
	Story   string `json:"story"`
	Started bool   `json:"started"`
	Running bool   `json:"running,omitempty"`
	// Recorded is a verdict from the verify record at Commit, of the run at
	// RanAt; nothing ran.
	Recorded  bool      `json:"recorded,omitempty"`
	Verdict   string    `json:"verdict,omitempty"`
	StoppedAt string    `json:"stopped_at,omitempty"`
	Commit    string    `json:"commit,omitempty"`
	RanAt     time.Time `json:"ran_at,omitzero"`
	// Log is where the job's output is kept.
	Log string `json:"log,omitempty"`
}

func (in VerifyIn) project() string { return in.Project }

// verify asks flai serve for the close-out of the story the call names as a
// job, as flai verify --job does, and answers its answer; it runs nothing in
// this process. A refusal and a request flai serve did not take are errors.
func (s *server) verify(ctx context.Context, _ *mcp.CallToolRequest, in VerifyIn) (*mcp.CallToolResult, VerifyOut, error) {
	story, err := s.storyWithWorktree(in.Story)
	if err != nil {
		return nil, VerifyOut{}, err
	}
	ans, err := closeout.Ask(ctx, s.repo.CacheDir(), closeout.Request{Story: story, By: s.agent, Message: in.Message, All: in.All, Fresh: in.Fresh}, s.closeOutWait)
	if err != nil {
		return nil, VerifyOut{}, err
	}
	if ans.Refused() {
		return nil, VerifyOut{}, errors.New(ans.Reason)
	}
	return nil, VerifyOut{Job: ans.Job, Story: ans.Story, Started: ans.Started, Running: ans.Running(), Recorded: ans.Recorded,
		Verdict: ans.Verdict, StoppedAt: ans.StoppedAt, Commit: ans.Commit, RanAt: ans.RanAt, Log: ans.Log}, nil
}

// storyWithWorktree is the story id names, as its canonical ID, which must
// have a worktree.
func (s *server) storyWithWorktree(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", errors.New("name the story to verify, such as S-0042")
	}
	it, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, workitem.ErrNotFound) {
			return "", fmt.Errorf("%s is not a work item here; name a story with a worktree", id)
		}
		return "", err
	}
	if it.Type != workitem.Story {
		return "", fmt.Errorf("%s is a %s; name a story (S-nnnn): verify runs for a story's worktree", it.ID, it.Type)
	}
	wt := s.repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", fmt.Errorf("%s has no worktree at %s; open it with flai stream open %s on the host", it.ID, s.rel(wt), it.ID)
	}
	return it.ID, nil
}
