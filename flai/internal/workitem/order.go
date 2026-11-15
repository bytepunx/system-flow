package workitem

import (
	"fmt"
	"slices"
	"sort"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// The pull order (S-0057). The board's order list names stories in the order
// they should be pulled: the stories of the pull's source first, then those
// of each state before it that reads in the pull order, in the order they
// should be refined, nearest the pull first (ADR-0301 §2). Only the states
// the workflow gives an order have one. A story the list does not name comes
// after the ones it does, by ID, so a queue nobody has prioritised still has
// one reading.

// Ordered reports whether a state's column has a pull order in the baseline
// workflow, as OrderedIn says.
func Ordered(state string) bool { return OrderedIn(baselineFlow, state) }

// OrderedIn reports whether the state of w named state, or that an alias
// names, reads in the pull order (ADR-0301 §2).
func OrderedIn(w workflow.Workflow, state string) bool {
	s, ok := w.State(state)
	return ok && s.Order
}

// PullSequence is the active stories of a state in pull order: those the list
// names first, in its order, then the rest by ID.
func PullSequence(order []string, items []*Item, state string) []string {
	in := map[string]bool{}
	for _, it := range items {
		if !it.Archived && it.Type == Story && it.Status == state {
			in[it.ID] = true
		}
	}
	seen := map[string]bool{}
	var out, rest []string
	for _, id := range order {
		if in[id] && !seen[id] {
			out = append(out, id)
			seen[id] = true
		}
	}
	for id := range in {
		if !seen[id] {
			rest = append(rest, id)
		}
	}
	sort.Slice(rest, func(i, j int) bool { return lessID(rest[i], rest[j]) })
	return append(out, rest...)
}

// Placement says where a story goes in its column's pull order: before or
// after another story of the same column, or at its top or bottom.
type Placement struct {
	Before, After string
	Top, Bottom   bool
}

func (p Placement) validate() error {
	n := 0
	for _, set := range []bool{p.Before != "", p.After != "", p.Top, p.Bottom} {
		if set {
			n++
		}
	}
	if n != 1 {
		return fmt.Errorf("rule: say where it goes with exactly one of --before, --after, --top, --bottom")
	}
	return nil
}

// Place puts a story of a state that reads in the pull order at a position in
// its column's pull order and rewrites the list: the stories of each such
// state, nearest the pull first, and nothing else (ADR-0301 §2). Stories of a
// state other than the pull's source that were never placed and still come
// last by ID stay unnamed, so one placement does not make the list name the
// whole queue. The pull's source's stories are always named, as Move names
// them.
// Refusals are rule errors, like Move's, so the dashboard shows them as the
// reason and not as a failure.
func (b *Board) Place(items []*Item, id string, p Placement) error {
	if err := p.validate(); err != nil {
		return err
	}
	byID := map[string]*Item{}
	for _, it := range items {
		if !it.Archived {
			byID[it.ID] = it
		}
	}
	it, ok := byID[id]
	if !ok {
		return fmt.Errorf("rule: %s is not an active item", id)
	}
	if it.Type != Story {
		return fmt.Errorf("rule: %s is a %s; only stories are in the pull order", id, it.Type)
	}
	w := b.Flow()
	if !OrderedIn(w, it.Status) {
		return fmt.Errorf("rule: %s is %s; only %s stories have a pull order", id, it.Status, and(orderedStates(w)))
	}
	ref := p.Before + p.After
	if ref != "" {
		other, ok := byID[ref]
		switch {
		case ref == id:
			return fmt.Errorf("rule: %s cannot be placed relative to itself", id)
		case !ok:
			return fmt.Errorf("rule: %s is not an active item", ref)
		case other.Type != Story:
			return fmt.Errorf("rule: %s is a %s; only stories are in the pull order", ref, other.Type)
		case other.Status != it.Status:
			return fmt.Errorf("rule: %s is %s and %s is %s; the pull order is within one column, and flai move changes the column", id, it.Status, ref, other.Status)
		}
	}

	named := map[string]bool{}
	for _, x := range b.Order {
		named[x] = true
	}
	var seq []string
	for _, x := range PullSequence(b.Order, items, it.Status) {
		if x != id {
			seq = append(seq, x)
		}
	}
	at := len(seq)
	switch {
	case p.Top:
		at = 0
	case ref != "":
		for i, x := range seq {
			if x == ref {
				at = i
				if p.After != "" {
					at = i + 1
				}
			}
		}
	}
	seq = append(seq[:at], append([]string{id}, seq[at:]...)...)

	source := pullSource(w)
	states := orderedStates(w)
	slices.Reverse(states) // nearest the pull first
	order := []string{}
	for _, st := range states {
		col := seq
		if st != it.Status {
			col = PullSequence(b.Order, items, st)
		}
		if st != source {
			col = col[:orderNamedEnd(col, named)]
		}
		order = append(order, col...)
	}
	b.Order = order
	return nil
}

// orderNamedEnd is where the unplaced tail of col, a column's pull sequence,
// starts: the stories after the last one named, that still come by ID. They
// are left unnamed, since reading the list back puts unnamed stories last by
// ID, which is where they already are.
func orderNamedEnd(col []string, named map[string]bool) int {
	end := len(col)
	for end > 0 {
		x := col[end-1]
		if named[x] || (end < len(col) && !lessID(x, col[end])) {
			break
		}
		end--
	}
	return end
}

// PlaceReadyLast names a story that has just entered the pull's source after
// the stories of that state the list already names and before any story of
// another state, wherever the list had it before (ADR-0301 §2).
func (b *Board) PlaceReadyLast(id string, items []*Item) {
	if items == nil {
		b.AppendToOrder(id)
		return
	}
	b.RemoveFromOrder(id)
	source := pullSource(b.Flow())
	queued := map[string]bool{}
	for _, it := range items {
		if !it.Archived && it.Type == Story && it.Status == source && it.ID != id {
			queued[it.ID] = true
		}
	}
	at := 0
	for i, x := range b.Order {
		if queued[x] {
			at = i + 1
		}
	}
	b.Order = append(b.Order[:at], append([]string{id}, b.Order[at:]...)...)
}
