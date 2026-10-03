package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const issueStoryDescription = "Make a backlog story that remediates an open issue (S-0198), as flai issue story makes one: it takes the issue's title; its nature is remediation for a defect or a blocker, improvement otherwise; its goal links the issue and carries the issue's recommended solution. That link is what ties the issue to the story; the issue's file is not changed. flai check runs with the story in place and refuses it, leaving nothing, if it reports anything the story introduces. A closed issue, or one an open story already links (named), is refused and nothing changes. epic is the story's epic (optional); story reads the issue from that story's worktree, where the issues it recorded are until it is accepted; the new story is made in wip/ as always. Nothing is committed."

// IssueStoryIn makes a story from an issue (S-0198).
type IssueStoryIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the issue, such as I-0007 (any zero padding)"`
	Epic    string `json:"epic,omitempty" jsonschema:"the story's epic (optional: a story need not belong to one)"`
	Story   string `json:"story,omitempty" jsonschema:"read the issue from this story's worktree when it has one"`
}

func (in IssueStoryIn) project() string { return in.Project }

// IssueStoryOut is the story made and the issue it remediates.
type IssueStoryOut struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Nature string `json:"nature"`
	Path   string `json:"path"`
	Issue  string `json:"issue"`
}

func (s *server) issueStory(ctx context.Context, _ *mcp.CallToolRequest, in IssueStoryIn) (*mcp.CallToolResult, IssueStoryOut, error) {
	is, err := issues.Get(issues.RepoFor(s.repo, in.Story), in.ID)
	if err != nil {
		return nil, IssueStoryOut{}, err
	}
	if is.Status != "open" {
		return nil, IssueStoryOut{}, fmt.Errorf("%s is %s, so no story was made for it; reopen it by editing its status, or record a new issue", is.ID, is.Status)
	}
	items, err := s.listItems(ctx)
	if err != nil {
		return nil, IssueStoryOut{}, err
	}
	if id := issues.LinkedBy(is, items); id != "" {
		return nil, IssueStoryOut{}, fmt.Errorf("%s is already linked by open story %s, so no story was made for it; work it in %s, or cancel %s first", is.ID, id, id, id)
	}
	owner := s.repo.Manifest.Owner
	if owner == "" {
		owner = s.agent
	}
	draft := issues.ForStory(is)
	res, err := itemnew.Create(s.repo, s.runner, itemnew.Options{New: workitem.NewOptions{
		Type: workitem.Story, Title: draft.Title, Nature: draft.Nature, Parent: in.Epic,
		Owner: owner, Body: draft.Body, Now: s.now(),
	}})
	if err != nil {
		return nil, IssueStoryOut{}, fmt.Errorf("making a story for %s: %w", is.ID, refusal(err))
	}
	return nil, IssueStoryOut{ID: res.Item.ID, Title: res.Item.Title, Nature: res.Item.Nature, Path: s.rel(res.Item.Path), Issue: is.ID}, nil
}
