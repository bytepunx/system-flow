package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const streamStateDescription = "Write a story's narrative's state as this agent (S-0271): current replaces what is under ## Current state and next what is under ## Next steps, each as markdown without a # or ## heading of its own. An empty text leaves its section as it is, but give at least one. Every other section of the narrative is left as it was and nothing is appended to its log: write the log with flai stream log on the host. The narrative's updated stamp, agent, and session are written as flai stream log writes them, wip/agents/index.md is written again, and nothing is committed. The same text again writes nothing, and changed is false. A story not in progress or in review is refused, as is one with no narrative (open it with flai stream open) and text the project's markdown lint rejects, with the findings; nothing is written then. Returns the narrative's path, the sections written, its updated stamp, and whether it changed. The CLI's flai stream state does the same."

// StreamStateIn is the text for a story's narrative's Current state and Next
// steps (S-0271).
type StreamStateIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story"`
	Current string `json:"current,omitempty" jsonschema:"the markdown for ## Current state, replacing what is under it; empty leaves the section as it is"`
	Next    string `json:"next,omitempty" jsonschema:"the markdown for ## Next steps, replacing what is under it; empty leaves the section as it is"`
}

func (in StreamStateIn) project() string { return in.Project }

// StreamStateOut is what the write did to the narrative.
type StreamStateOut struct {
	Story   string   `json:"story"`
	Path    string   `json:"path" jsonschema:"the narrative, relative to the project's root"`
	Written []string `json:"written" jsonschema:"the sections given text, Current state and Next steps, in that order"`
	Updated string   `json:"updated" jsonschema:"the narrative's updated stamp after the write, or as it was when nothing changed"`
	Changed bool     `json:"changed" jsonschema:"false when the narrative already held the text given, and nothing was written"`
}

func (s *server) streamState(_ context.Context, _ *mcp.CallToolRequest, in StreamStateIn) (*mcp.CallToolResult, StreamStateOut, error) {
	if strings.TrimSpace(in.Current) == "" && strings.TrimSpace(in.Next) == "" {
		return nil, StreamStateOut{}, fmt.Errorf("nothing to write to %s's narrative: give current, next, or both, the text for ## Current state and ## Next steps", in.Story)
	}
	res, err := s.repo.SetStreamState(in.Story, in.Current, in.Next, workitem.StreamOptions{Agent: s.agent, Session: s.session, Now: s.now()})
	if err != nil {
		return nil, StreamStateOut{}, streamRefusal(err)
	}
	return nil, StreamStateOut{Story: res.Stream, Path: s.rel(res.Path), Written: res.Written, Updated: res.Updated, Changed: res.Changed}, nil
}

// streamRefusal is err with what to do about it when the project's lint
// refused the text; a refusal of the story's state names it already.
func streamRefusal(err error) error {
	var lint *mdlint.Error
	if errors.As(err, &lint) {
		return fmt.Errorf("%w; change the text so that it passes, and call stream_state again", err)
	}
	return err
}
