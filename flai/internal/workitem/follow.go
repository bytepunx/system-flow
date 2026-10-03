package workitem

import (
	"fmt"
	"time"
)

// Followed is an epic's whole walk for one story's move: the state it was in,
// the state it reached, and the story it followed (S-0200).
type Followed struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	From  string `json:"from"`
	To    string `json:"to"`
	Story string `json:"story"`
}

// followRank orders the open states and done along the board; cancelled has
// no place in it.
var followRank = map[string]int{Backlog: 0, Ready: 1, InProgress: 2, Review: 3, Done: 4}

// followPath is the board's states in order, the walk an epic takes one at a
// time.
var followPath = []string{Backlog, Ready, InProgress, Review, Done}

// EpicFollows is the state epicID's stories put it in, archived ones included
// and cancelled ones left out; ok is false when no story counts (S-0200).
func EpicFollows(items []*Item, epicID string) (string, bool) {
	n, done, late, started, ready := 0, 0, 0, 0, 0
	for _, s := range Children(items, epicID) {
		if s.Type != Story || s.Status == Cancelled {
			continue
		}
		n++
		switch s.Status {
		case Done:
			done++
			late++
			started++
		case Review:
			late++
			started++
		case InProgress:
			started++
		case Ready:
			ready++
		}
	}
	switch {
	case n == 0:
		return "", false
	case done == n:
		return Done, true
	case late == n:
		return Review, true
	case started > 0:
		return InProgress, true
	case ready > 0:
		return Ready, true
	}
	return Backlog, true
}

// Follow moves story's epic in memory, one allowed transition at a time, the
// way story's move from from takes it: forward with a forward move or a
// cancellation, back with a back move or a return from cancelled, never into
// done unless accept. The caller saves the epic it returns (S-0200).
func (r *Repo) Follow(items []*Item, story *Item, from, by string, now time.Time, accept bool) (*Item, *Followed, error) {
	to := story.Status
	return r.follow(items, story, story.Parent, from, to, fmt.Sprintf("follows %s, which moved to %s", story.ID, to), by, now, accept)
}

// FollowParent moves in memory the epics story left and joined, when an edit
// gave it the parent it has instead of former: the one it joined as though
// the story had entered it from backlog, the one it left as though the story
// had been cancelled out of it, so forward only, short of done. The caller
// saves the epics it returns (S-0200).
func (r *Repo) FollowParent(items []*Item, story *Item, former, by string, now time.Time) ([]*Item, []*Followed, error) {
	gone := *story
	gone.Status, gone.Parent = Cancelled, former
	var epics []*Item
	var walks []*Followed
	for _, f := range []struct {
		story        *Item
		from, reason string
	}{
		{&gone, story.Status, fmt.Sprintf("follows %s, which left it for %s", story.ID, orDefault(story.Parent, "no epic"))},
		{story, Backlog, fmt.Sprintf("follows %s, which joined it from %s", story.ID, orDefault(former, "no epic"))},
	} {
		epic, walk, err := r.follow(items, f.story, f.story.Parent, f.from, f.story.Status, f.reason, by, now, false)
		if err != nil {
			return nil, nil, err
		}
		if walk != nil {
			walk.Story = story.ID
			epics, walks = append(epics, epic), append(walks, walk)
		}
	}
	return epics, walks, nil
}

// follow walks epicID the way story's move from from to to takes it, with
// reason on every step.
func (r *Repo) follow(items []*Item, story *Item, epicID, from, to, reason, by string, now time.Time, accept bool) (*Item, *Followed, error) {
	// The caller's list may hold a copy of story read before it moved.
	all := make([]*Item, len(items))
	var epic *Item
	for i, it := range items {
		if it.ID == story.ID {
			it = story
		}
		all[i] = it
		if it.ID == epicID && it.Type == Epic {
			epic = it
		}
	}
	if epic == nil || epic.Archived || epic.Closed() {
		return nil, nil, nil
	}
	target, ok := EpicFollows(all, epic.ID)
	if !ok {
		return nil, nil, nil
	}
	forward := to == Cancelled || (from != Cancelled && followRank[to] > followRank[from])
	back := from == Cancelled || (to != Cancelled && followRank[to] < followRank[from])
	if forward && !accept && followRank[target] > followRank[Review] {
		target = Review
	}
	step := 0
	switch {
	case forward && followRank[epic.Status] < followRank[target]:
		step = 1
	case back && followRank[epic.Status] > followRank[target]:
		step = -1
	default:
		return nil, nil, nil
	}
	was := &Followed{ID: epic.ID, Type: epic.Type, Title: epic.Title, From: epic.Status, Story: story.ID}
	for epic.Status != target {
		next := followPath[followRank[epic.Status]+step]
		// An epic's own warnings (an open blocked interval) are the
		// operator's, who moves epics; the story's move does not repeat them.
		if _, err := r.Move(epic, next, MoveOptions{By: by, Reason: reason, Now: now, Items: all}); err != nil {
			return nil, nil, fmt.Errorf("epic %s cannot follow %s to %s, so nothing was changed: %w", epic.ID, story.ID, next, err)
		}
	}
	was.To = epic.Status
	return epic, was, nil
}
