package workitem

import "time"

// Transition performs a state change end to end, the way every caller must:
// the workflow rules, the item file, the board's pull order, and the
// narrative index. flai move and the MCP server both go through it, so a
// transition means the same thing wherever it is made.
func (r *Repo) Transition(it *Item, to, by, reason string, now time.Time) (warnings []string, err error) {
	items, err := r.List(false)
	if err != nil {
		return nil, err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return nil, err
	}
	warnings, err = r.Move(it, to, MoveOptions{By: by, Reason: reason, Now: now, Items: items, Board: board})
	if err != nil {
		return nil, err
	}
	if err := r.Save(it); err != nil {
		return nil, err
	}
	if it.Type == Story {
		if err := board.Save(now.Format("2006-01-02")); err != nil {
			return nil, err
		}
	}
	items, err = r.List(false)
	if err != nil {
		return nil, err
	}
	return warnings, r.WriteIndex(items, now)
}
