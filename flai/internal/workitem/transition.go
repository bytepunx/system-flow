package workitem

import (
	"fmt"
	"time"
)

// MoveResult is what a transition did beyond the item itself.
type MoveResult struct {
	Warnings []string
	// Cancelled lists the items cancelled with the one that was moved, in
	// the order of CancelPlan, each as it was before: a cancelled parent
	// takes everything open under it along (S-0070).
	Cancelled []Cascaded
}

// Cascaded is one item a cancellation takes with it, and the state it was in.
type Cascaded struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	From  string `json:"from"`
}

// CancelPlan lists what cancelling id would cancel with it: every open,
// unarchived story under an epic and every open task under those, or the
// open tasks under a story. Each story is followed by its tasks.
func CancelPlan(items []*Item, id string) []*Item {
	var out []*Item
	for _, c := range Children(items, id) {
		if c.Archived {
			continue
		}
		if !c.Closed() {
			out = append(out, c)
		}
		// A closed story can still hold open tasks (cancelled by an older
		// flai, or edited by hand); the cascade closes those too.
		out = append(out, CancelPlan(items, c.ID)...)
	}
	return out
}

// Transition is the shared implementation behind flai move and the MCP
// item_move tool: load the context, apply the rules, save the item and the
// board, and rewrite the stream index.
func (r *Repo) Transition(it *Item, to, by, reason string, now time.Time) (warnings []string, err error) {
	res, err := r.TransitionAll(it, to, by, reason, now)
	if err != nil {
		return nil, err
	}
	return res.Warnings, nil
}

// TransitionAll is Transition that also reports what went with the item. A
// move to cancelled cancels everything open under the item with the same
// actor and time, and a note naming the item that caused it. Every move is
// validated before any file is written, so a refusal changes nothing.
func (r *Repo) TransitionAll(it *Item, to, by, reason string, now time.Time) (*MoveResult, error) {
	items, err := r.List(false)
	if err != nil {
		return nil, err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return nil, err
	}
	res := &MoveResult{Cancelled: []Cascaded{}}
	res.Warnings, err = r.Move(it, to, MoveOptions{By: by, Reason: reason, Now: now, Items: items, Board: board})
	if err != nil {
		return nil, err
	}
	changed := []*Item{it}
	storyMoved := it.Type == Story
	if to == Cancelled {
		cause := fmt.Sprintf("%s cancelled: %s", it.ID, reason)
		for _, c := range CancelPlan(items, it.ID) {
			was := Cascaded{ID: c.ID, Type: c.Type, Title: c.Title, From: c.Status}
			if _, err := r.Move(c, Cancelled, MoveOptions{By: by, Reason: cause, Now: now, Items: items, Board: board, Cascade: true}); err != nil {
				return nil, fmt.Errorf("%s cannot be cancelled with %s, so nothing was changed: %w", c.ID, it.ID, err)
			}
			res.Cancelled = append(res.Cancelled, was)
			changed = append(changed, c)
			storyMoved = storyMoved || c.Type == Story
		}
	}
	for _, c := range changed {
		if err := r.Save(c); err != nil {
			return nil, err
		}
	}
	if storyMoved {
		if err := board.Save(now.Format("2006-01-02")); err != nil {
			return nil, err
		}
	}
	items, err = r.List(false)
	if err != nil {
		return nil, err
	}
	return res, r.WriteIndex(items, now)
}
