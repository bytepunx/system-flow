package mcpserver

import (
	"context"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
)

// Edited is the kind of a change someone made to an item's own fields or
// body with flai edit; To names what of it changed.
const Edited = inbox.Edited

// Overlapped is the kind of the change to an open story when a story whose
// changes its claim covers is accepted, or when a write grows the claims of
// two stories in progress to overlap (S-0132, I-0059).
const Overlapped = inbox.Overlapped

// maxEvents is the most changes one look reports.
const maxEvents = inbox.MaxEvents

// maxPaths is how many paths an overlap's summary names.
const maxPaths = inbox.MaxPaths

// Event is one thing that changed since the agent last looked.
type Event = inbox.Event

// describe says a change in words.
var describe = inbox.Describe

// catchUp returns what others changed since the agent's cursor, and the
// messages to this session's own story (S-0331), the newest maxEvents of
// them, and how many older ones it left out, and advances the cursor past all
// of them, as inbox.CatchUp does.
func (s *server) catchUp(ctx context.Context) ([]Event, int, error) {
	items, err := s.listItems(ctx)
	if err != nil {
		return nil, 0, err
	}
	return inbox.CatchUp(s.repo, s.agent, s.recordingStory(""), s.now(), items)
}
