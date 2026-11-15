package workitem

import (
	"fmt"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
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

// EpicFollows is the state epicID's stories put it in under the baseline
// workflow, as EpicFollowsIn reads it (S-0200).
func EpicFollows(items []*Item, epicID string) (string, bool) {
	return EpicFollowsIn(baselineFlow, items, epicID)
}

// EpicFollowsIn is the state epicID's stories put it in under w, archived
// stories included and those in a state of the cancelled kind left out; ok
// is false when no story counts. It reads kinds and roles, never a state's
// name (ADR-0301 §2, ADR-0076): the accepted state when every story is
// accepted; the acceptance's source when every story is there or accepted;
// the pull's target when any story has started, in an active, waiting, or
// accepted state; the pull's source when any story is there; else the
// initial state (S-0200).
func EpicFollowsIn(w workflow.Workflow, items []*Item, epicID string) (string, bool) {
	waiting, queued := acceptSource(w), pullSource(w)
	n, accepted, late, started, pullable := 0, 0, 0, 0, 0
	for _, s := range Children(items, epicID) {
		if s.Type != Story {
			continue
		}
		st, _ := w.State(s.Status)
		if st.Kind == workflow.KindCancelled {
			continue
		}
		n++
		switch {
		case st.Kind == workflow.KindAccepted:
			accepted++
			late++
		case st.Name == waiting:
			late++
		case st.Name == queued:
			pullable++
		}
		if startedKind(st.Kind) {
			started++
		}
	}
	switch {
	case n == 0:
		return "", false
	case accepted == n:
		return acceptedState(w), true
	case late == n:
		return waiting, true
	case started > 0:
		return pullTarget(w), true
	case pullable > 0:
		return queued, true
	}
	return initialState(w), true
}

// followPathIn is the walk an epic takes under w: the workflow's states in
// board order, those of the cancelled kind left out, which have no place in
// it (ADR-0301 §2).
func followPathIn(w workflow.Workflow) []string {
	var out []string
	for _, s := range w.States {
		if s.Kind != workflow.KindCancelled {
			out = append(out, s.Name)
		}
	}
	return out
}

// followRankIn is state's place along path, by its current name when it is
// an alias; a state off the path, of the cancelled kind or unknown, ranks
// first, as the walk has always ranked it.
func followRankIn(w workflow.Workflow, path []string, state string) int {
	if s, ok := w.State(state); ok {
		state = s.Name
	}
	if i := slices.Index(path, state); i >= 0 {
		return i
	}
	return 0
}

// followStepIn is the state an epic of rank from along path moves to next on
// its way to rank to: the nearest in that direction a transition for an epic
// takes it to, so that it walks one transition at a time and passes each
// state between that it may enter (ADR-0301 §2). It is the next state along
// path when no transition leads to any, which the move then refuses.
func followStepIn(w workflow.Workflow, path []string, from, to int) string {
	step := 1
	if to < from {
		step = -1
	}
	for i := from + step; i != to+step; i += step {
		if _, ok := w.Find(Epic, path[from], path[i]); ok {
			return path[i]
		}
	}
	return path[from+step]
}

// Follow moves story's epic in memory, one allowed transition at a time, the
// way story's move from from takes it under the repository's workflow:
// forward with a forward move or a cancellation, back with a back move or a
// return from a cancelled state, never into the accepted state unless accept
// (ADR-0301 §2). The caller saves the epic it returns (S-0200).
func (r *Repo) Follow(items []*Item, story *Item, from, by string, now time.Time, accept bool) (*Item, *Followed, error) {
	to := story.Status
	return r.follow(items, story, story.Parent, from, to, fmt.Sprintf("follows %s, which moved to %s", story.ID, to), by, now, accept)
}

// FollowParent moves in memory the epics story left and joined, when an edit
// gave it the parent it has instead of former: the one it joined as though
// the story had entered it from the initial state, the one it left as though
// the story had been cancelled out of it, so forward only, short of the
// accepted state (ADR-0301 §2). The caller saves the epics it returns
// (S-0200).
func (r *Repo) FollowParent(items []*Item, story *Item, former, by string, now time.Time) ([]*Item, []*Followed, error) {
	w := r.Flow()
	gone := *story
	gone.Status, gone.Parent = cancelledState(w), former
	var epics []*Item
	var walks []*Followed
	for _, f := range []struct {
		story        *Item
		from, reason string
	}{
		{&gone, story.Status, fmt.Sprintf("follows %s, which left it for %s", story.ID, orDefault(story.Parent, "no epic"))},
		{story, initialState(w), fmt.Sprintf("follows %s, which joined it from %s", story.ID, orDefault(former, "no epic"))},
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
// reason on every step, along the repository's workflow: by rank in its walk,
// clamped to the acceptance's source unless accept (ADR-0301 §2).
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
	w := r.Flow()
	target, ok := EpicFollowsIn(w, all, epic.ID)
	if !ok {
		return nil, nil, nil
	}
	path := followPathIn(w)
	rank := func(state string) int { return followRankIn(w, path, state) }
	fromGone := w.Kind(from) == workflow.KindCancelled
	toGone := w.Kind(to) == workflow.KindCancelled
	forward := toGone || (!fromGone && rank(to) > rank(from))
	back := fromGone || (!toGone && rank(to) < rank(from))
	if waiting := acceptSource(w); forward && !accept && waiting != "" && rank(target) > rank(waiting) {
		target = waiting
	}
	if !slices.Contains(path, target) {
		// A workflow that fills no such role gives the epic nowhere to go.
		return nil, nil, nil
	}
	switch {
	case forward && rank(epic.Status) < rank(target):
	case back && rank(epic.Status) > rank(target):
	default:
		return nil, nil, nil
	}
	was := &Followed{ID: epic.ID, Type: epic.Type, Title: epic.Title, From: epic.Status, Story: story.ID}
	for epic.Status != target {
		next := followStepIn(w, path, rank(epic.Status), rank(target))
		// An epic's own warnings (an open blocked interval) are the
		// operator's, who moves epics; the story's move does not repeat them.
		if _, err := r.Move(epic, next, MoveOptions{By: by, Reason: reason, Now: now, Items: all}); err != nil {
			return nil, nil, fmt.Errorf("epic %s cannot follow %s to %s, so nothing was changed: %w", epic.ID, story.ID, next, err)
		}
	}
	was.To = epic.Status
	return epic, was, nil
}
