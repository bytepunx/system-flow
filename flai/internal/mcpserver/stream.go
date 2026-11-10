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

const streamStateDescription = "Write a story's narrative's state and record its context, decisions, and open questions as this agent (S-0271, S-0668): current replaces what is under ## Current state, next what is under ## Next steps, and context what is under ## Context, each as markdown without a # or ## heading of its own; each of decisions is appended as a bullet under ## Decisions, and each of questions as a bullet dated today under ## Open questions, above the block that mirrors threads, which is kept as it is. A question so written holds the move to review until it is answered, as a hand-written one does. Appending keeps what earlier sessions recorded. An empty text leaves its section as it is, but give at least one of the five. This is how an agent whose wip/ is read-only, in a container (ADR-0253), records them: the server runs on the host and writes the main checkout. Every other section of the narrative is left as it was and nothing is appended to its log: write the log with flai stream log on the host. The narrative's updated stamp, agent, and session are written as flai stream log writes them, wip/agents/index.md is written again, and nothing is committed. The same current, next, or context again writes nothing, and changed is false; a decision or question given again is appended again. A story not in progress or in review is refused, as is one with no narrative (open it with flai stream open) and text the project's markdown lint rejects, with the findings; nothing is written then. Returns the narrative's path, the sections written, its updated stamp, and whether it changed. The CLI's flai stream state does the same."

// StreamStateIn is the text for a story's narrative's Current state, Next
// steps, and Context, and the decisions and open questions to append to it
// (S-0271, S-0668).
type StreamStateIn struct {
	Project   string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story     string   `json:"story"`
	Current   string   `json:"current,omitempty" jsonschema:"the markdown for ## Current state, replacing what is under it; empty leaves the section as it is"`
	Next      string   `json:"next,omitempty" jsonschema:"the markdown for ## Next steps, replacing what is under it; empty leaves the section as it is"`
	Context   string   `json:"context,omitempty" jsonschema:"the markdown for ## Context, replacing what is under it; empty leaves the section as it is"`
	Decisions []string `json:"decisions,omitempty" jsonschema:"decisions to append, each as one bullet at the end of ## Decisions; what the section holds stays"`
	Questions []string `json:"questions,omitempty" jsonschema:"open questions to append, each as one bullet dated today under ## Open questions, above the block that mirrors threads; each holds the move to review until it is answered, as a hand-written one does"`
}

func (in StreamStateIn) project() string { return in.Project }

// StreamStateOut is what the write did to the narrative.
type StreamStateOut struct {
	Story   string   `json:"story"`
	Path    string   `json:"path" jsonschema:"the narrative, relative to the project's root"`
	Written []string `json:"written" jsonschema:"the sections given text, in the narrative's order: Context, Current state, Next steps, Decisions, and Open questions"`
	Updated string   `json:"updated" jsonschema:"the narrative's updated stamp after the write, or as it was when nothing changed"`
	Changed bool     `json:"changed" jsonschema:"false when the narrative already held the text given, and nothing was written"`
}

func (s *server) streamState(_ context.Context, _ *mcp.CallToolRequest, in StreamStateIn) (*mcp.CallToolResult, StreamStateOut, error) {
	if strings.TrimSpace(in.Current) == "" && strings.TrimSpace(in.Next) == "" && strings.TrimSpace(in.Context) == "" && blank(in.Decisions) && blank(in.Questions) {
		return nil, StreamStateOut{}, fmt.Errorf("nothing to write to %s's narrative: give current, next, or both, the text for ## Current state and ## Next steps, or context, decisions, or questions, for ## Context, ## Decisions, and ## Open questions", in.Story)
	}
	text := workitem.StreamText{Current: in.Current, Next: in.Next, Context: in.Context, Decisions: in.Decisions, Questions: in.Questions}
	res, err := s.repo.WriteStreamState(in.Story, text, workitem.StreamOptions{Agent: s.agent, Session: s.session, Now: s.now()})
	if err != nil {
		return nil, StreamStateOut{}, streamRefusal(err)
	}
	return nil, StreamStateOut{Story: res.Stream, Path: s.rel(res.Path), Written: res.Written, Updated: res.Updated, Changed: res.Changed}, nil
}

// blank is true when every text is empty or only spaces.
func blank(texts []string) bool {
	for _, t := range texts {
		if strings.TrimSpace(t) != "" {
			return false
		}
	}
	return true
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
