package mcpserver

import (
	"context"
	"errors"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
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
	Agent   *manifest.Agent `json:"agent,omitempty" jsonschema:"a story's agent: harness, model, and config, over the project's default, which fills in what is not given"`
	Body    string          `json:"body,omitempty" jsonschema:"the goal, criteria, and notes below the heading; the template's empty sections when not given"`
}

func (in ItemNewIn) project() string { return in.Project }

func (s *server) itemNew(_ context.Context, _ *mcp.CallToolRequest, in ItemNewIn) (*mcp.CallToolResult, ItemOut, error) {
	nature := in.Nature
	if nature == "" {
		nature = "feature"
	}
	owner := s.repo.Manifest.Owner
	if owner == "" {
		owner = s.agent
	}
	it, err := s.repo.Create(workitem.NewOptions{Type: in.Type, Title: strings.TrimSpace(in.Title), Nature: nature, Parent: in.Parent, Owner: owner,
		Tags: in.Tags, Touches: in.Touches, Agent: in.Agent, Body: in.Body, Now: s.now()})
	if err != nil {
		return nil, ItemOut{}, err
	}
	out, err := s.itemOut(it)
	return nil, out, err
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
	Parent  *string   `json:"parent,omitempty"`
	// Agent replaces a story's agent; ClearAgent removes it.
	Agent      *manifest.Agent `json:"agent,omitempty" jsonschema:"replaces the story's agent with exactly this harness, model, and config"`
	ClearAgent bool            `json:"clear_agent,omitempty" jsonschema:"removes the story's agent"`
	Body       *string         `json:"body,omitempty" jsonschema:"replaces everything below the heading"`
}

func (in ItemEditIn) project() string { return in.Project }

// ItemEditOut is what an edit changed.
type ItemEditOut struct {
	ID        string   `json:"id"`
	Changed   []string `json:"changed" jsonschema:"title, nature, tags, touches, agent, parent, goal, criteria, notes, body"`
	Unchanged bool     `json:"unchanged,omitempty"`
	Hash      string   `json:"hash"`
}

func (s *server) itemEdit(_ context.Context, _ *mcp.CallToolRequest, in ItemEditIn) (*mcp.CallToolResult, ItemEditOut, error) {
	ch := itemedit.Change{Title: in.Title, Nature: in.Nature, Tags: in.Tags, Touches: in.Touches, Parent: in.Parent, Body: in.Body, Agent: in.Agent, ClearAgent: in.ClearAgent}
	if ch == (itemedit.Change{}) {
		return nil, ItemEditOut{}, errors.New("nothing to change: give title, nature, tags, touches, parent, agent, clear_agent, or body")
	}
	res, err := itemedit.Apply(s.repo, s.runner, in.ID, ch, itemedit.Options{Hash: in.Hash, By: s.agent, NoCommit: true, Now: s.now()})
	if err != nil {
		return nil, ItemEditOut{}, err
	}
	changed := res.Changed
	if changed == nil {
		changed = []string{}
	}
	return nil, ItemEditOut{ID: res.ID, Changed: changed, Unchanged: res.Unchanged, Hash: res.Hash}, nil
}
