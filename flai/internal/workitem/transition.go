package workitem

import (
	"fmt"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// The helpers below read a move by what the workflow says of it (ADR-0301
// §2, §3): the kind of the state it reaches, the effects of its transition,
// and who may make it, so that flai move, flai accept, the formula commands,
// and the previews name no state.

// Accepts reports whether moving an item of type typ from one state to the
// other is an acceptance in w: its transition gives the accept effect. On the
// baseline, an epic's or a story's move from review to done.
func Accepts(w workflow.Workflow, typ, from, to string) bool {
	t, ok := w.Find(typ, from, to)
	return ok && t.HasEffect(workflow.EffectAccept)
}

// AcceptSource is the state the acceptance of an item of type typ leaves
// from in w, the first its transition names, empty when w has none.
func AcceptSource(w workflow.Workflow, typ string) string {
	t, ok := w.WithEffect(typ, workflow.EffectAccept)
	if !ok || len(t.From) == 0 {
		return ""
	}
	return t.From[0]
}

// AcceptTarget is the state the acceptance of an item of type typ reaches in
// w, else w's first state of the accepted kind.
func AcceptTarget(w workflow.Workflow, typ string) string {
	if t, ok := w.WithEffect(typ, workflow.EffectAccept); ok {
		return t.To
	}
	return FirstOfKind(w, workflow.KindAccepted)
}

// CancelledState is w's first state of the cancelled kind, where an item
// closes undelivered.
func CancelledState(w workflow.Workflow) string { return FirstOfKind(w, workflow.KindCancelled) }

// FirstOfKind is the name of w's first state of kind in board order, empty
// when there is none.
func FirstOfKind(w workflow.Workflow, kind string) string {
	if states := w.OfKind(kind); len(states) > 0 {
		return states[0].Name
	}
	return ""
}

// IsKind reports whether state, by its name or an alias, is of kind in w.
func IsKind(w workflow.Workflow, state, kind string) bool {
	return kind != "" && w.Kind(state) == kind
}

// PermitMove is nil when m may make the move of it to the state to in w, by
// its transition's by (ADR-0301 §3), else a rule refusal naming the
// transition, who may make it, and the command that does. A move w has no
// transition for, or one only a parent's cancellation makes, is Move's to
// refuse, in its own words.
func PermitMove(w workflow.Workflow, it *Item, to string, m workflow.Mover) error {
	t, ok := w.Find(it.Type, it.Status, to)
	if !ok || cascadeOnly(t) {
		return nil
	}
	err := t.Permits(m)
	if err == nil {
		return nil
	}
	command := fmt.Sprintf("flai move %s %s", it.ID, t.To)
	if len(t.By) == 1 && t.By[0] == workflow.ByAcceptance {
		command = "flai accept " + it.ID
	}
	return fmt.Errorf("rule: %s: %w (%s)", it.ID, err, command)
}

// MoveResult is what a transition did beyond the item itself.
type MoveResult struct {
	Warnings []string
	// Cancelled lists the items cancelled with the one that was moved, in
	// the order of CancelPlan, each as it was before: a cancelled parent
	// takes everything open under it along (S-0070).
	Cancelled []Cascaded
	// RolledUp lists the items above one that entered done whose usage was
	// summed again (S-0143).
	RolledUp []string
	// Followed is the walk the moved story's epic took with it, nil when it
	// took none (S-0200).
	Followed *Followed
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
// board, and rewrite the stream index. It does not finalize a draft story;
// TransitionAll does when asked.
func (r *Repo) Transition(it *Item, to, by, reason string, now time.Time) (warnings []string, err error) {
	res, err := r.TransitionAll(it, to, by, reason, now, false)
	if err != nil {
		return nil, err
	}
	return res.Warnings, nil
}

// TransitionAll is Transition that also reports what went with the item. A
// move to cancelled cancels everything open under the item with the same
// actor and time, and a note naming the item that caused it. Every move is
// validated before any file is written, so a refusal changes nothing. With
// finalize, a draft story moved to ready is finalized with the move (S-0199).
// A moved story's epic follows it, short of done, in the same write
// (S-0200); the stories a cancellation takes along do not move their epic.
func (r *Repo) TransitionAll(it *Item, to, by, reason string, now time.Time, finalize bool) (*MoveResult, error) {
	all, err := r.List(true)
	if err != nil {
		return nil, err
	}
	var items []*Item
	for _, x := range all {
		if !x.Archived {
			items = append(items, x)
		}
	}
	board, err := r.LoadBoard()
	if err != nil {
		return nil, err
	}
	res := &MoveResult{Cancelled: []Cascaded{}}
	from := it.Status
	res.Warnings, err = r.Move(it, to, MoveOptions{By: by, Reason: reason, Now: now, Items: items, Board: board, Finalize: finalize})
	if err != nil {
		return nil, err
	}
	// Move recorded the state's current name, never an alias (ADR-0301 §6)
	to = it.Status
	wf := r.Flow()
	changed := []*Item{it}
	storyMoved := it.Type == Story
	// a cancellation takes everything open under the item to the state it
	// reached, along the workflow's transitions (ADR-0028)
	if IsKind(wf, to, workflow.KindCancelled) {
		cause := fmt.Sprintf("%s cancelled: %s", it.ID, reason)
		for _, c := range CancelPlan(items, it.ID) {
			was := Cascaded{ID: c.ID, Type: c.Type, Title: c.Title, From: c.Status}
			if _, err := r.Move(c, to, MoveOptions{By: by, Reason: cause, Now: now, Items: items, Board: board, Cascade: true}); err != nil {
				return nil, fmt.Errorf("%s cannot be cancelled with %s, so nothing was changed: %w", c.ID, it.ID, err)
			}
			res.Cancelled = append(res.Cancelled, was)
			changed = append(changed, c)
			storyMoved = storyMoved || c.Type == Story
		}
	}
	if it.Type == Story {
		// The archive counts: a story accepted and archived still holds its epic.
		epic, followed, err := r.Follow(all, it, from, by, now, false)
		if err != nil {
			return nil, err
		}
		if followed != nil {
			changed = append(changed, epic)
			res.Followed = followed
		}
	}
	for _, c := range changed {
		if err := r.Save(c); err != nil {
			return nil, err
		}
	}
	if IsKind(wf, to, workflow.KindAccepted) {
		if res.RolledUp, err = r.RollUp(it); err != nil {
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
