package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
)

const criteriaTickDescription = "Tick, or untick, a story's acceptance criteria by number once you have verified each (work-management.md): tick and untick name them, numbered from 1 in the order of the boxes under ## Acceptance criteria, as item_get's body and flai criteria list show them. Only those boxes change, every other byte of the item stays as it was, and nothing is committed. A number with no box, or one both to tick and to untick, is refused and nothing is written. Give the hash item_get returned so that a change made meanwhile is refused rather than ticking the wrong box. Returns the criteria after the write, with the new hash. The CLI's flai criteria and the dashboard do the same."

// CriteriaTickIn ticks and unticks a story's acceptance criteria by number
// (S-0282).
type CriteriaTickIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id"`
	Tick    []int  `json:"tick,omitempty" jsonschema:"the numbers of the criteria to tick, from 1, in the order of the boxes under ## Acceptance criteria"`
	Untick  []int  `json:"untick,omitempty" jsonschema:"the numbers of the criteria to untick, from 1, in the order of the boxes under ## Acceptance criteria"`
	Hash    string `json:"hash,omitempty" jsonschema:"the hash item_get gave; a change made meanwhile is then refused rather than ticking the wrong box"`
}

func (in CriteriaTickIn) project() string { return in.Project }

// CriteriaTickOut is the item after the write and its criteria.
type CriteriaTickOut struct {
	ID        string               `json:"id"`
	Path      string               `json:"path"`
	Hash      string               `json:"hash"`
	Changed   []string             `json:"changed" jsonschema:"criteria when a box changed, otherwise empty"`
	Unchanged bool                 `json:"unchanged,omitempty" jsonschema:"every box named was already as asked"`
	Criteria  []itemedit.Criterion `json:"criteria" jsonschema:"the acceptance criteria after the write, each with its number, its words, and whether it is ticked"`
}

func (s *server) criteriaTick(_ context.Context, _ *mcp.CallToolRequest, in CriteriaTickIn) (*mcp.CallToolResult, CriteriaTickOut, error) {
	v, err := itemedit.Show(s.repo, in.ID)
	if err != nil {
		return nil, CriteriaTickOut{}, err
	}
	// the item's state and the hash come before the numbers: an archived or
	// closed item, or one changed since it was read, is refused as item_edit
	// refuses it, whatever the numbers name
	if !v.Editable {
		return nil, CriteriaTickOut{}, &docedit.RefusedError{Path: v.Path, Reason: v.Reason}
	}
	hash := in.Hash
	if hash == "" {
		hash = v.Hash
	}
	if hash != v.Hash {
		return nil, CriteriaTickOut{}, &docedit.ConflictError{Path: v.Path, Hash: v.Hash}
	}
	body, err := itemedit.Tick(v.Body, in.Tick, in.Untick)
	if err != nil {
		return nil, CriteriaTickOut{}, err
	}
	res, err := itemedit.Apply(s.repo, s.runner, in.ID, itemedit.Change{Body: &body}, itemedit.Options{Hash: hash, By: s.agent, NoCommit: true, Now: s.now()})
	if err != nil {
		return nil, CriteriaTickOut{}, refusal(err)
	}
	changed := res.Changed
	if changed == nil {
		changed = []string{}
	}
	return nil, CriteriaTickOut{ID: res.ID, Path: res.Path, Hash: res.Hash, Changed: changed, Unchanged: res.Unchanged, Criteria: itemedit.Criteria(body)}, nil
}
