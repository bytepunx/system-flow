package workitem

import (
	"fmt"
	"sort"
)

// The pull order (S-0057). The board's order list names stories in the order
// they should be pulled: ready stories first, then backlog stories in the
// order they should be refined. Only those two columns have an order. A
// story the list does not name comes after the ones it does, by ID, so a
// backlog nobody has prioritised still has one reading.

// Ordered reports whether a state's column has a pull order.
func Ordered(state string) bool { return state == Ready || state == Backlog }

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
		return fmt.Errorf("say where it goes with exactly one of --before, --after, --top, --bottom")
	}
	return nil
}

// Place puts a ready or backlog story at a position in its column's pull
// order and rewrites the list: ready stories, then backlog stories, and
// nothing else. Backlog stories that were never placed and still come last by
// ID stay unnamed, so one placement does not make the list name the whole
// backlog. Ready stories are always named, as Move names them.
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
		return fmt.Errorf("%s is not an active item", id)
	}
	if it.Type != Story {
		return fmt.Errorf("%s is a %s; only stories are in the pull order", id, it.Type)
	}
	if !Ordered(it.Status) {
		return fmt.Errorf("%s is %s; only backlog and ready stories have a pull order", id, it.Status)
	}
	ref := p.Before + p.After
	if ref != "" {
		other, ok := byID[ref]
		switch {
		case ref == id:
			return fmt.Errorf("%s cannot be placed relative to itself", id)
		case !ok:
			return fmt.Errorf("%s is not an active item", ref)
		case other.Type != Story:
			return fmt.Errorf("%s is a %s; only stories are in the pull order", ref, other.Type)
		case other.Status != it.Status:
			return fmt.Errorf("%s is %s and %s is %s; the pull order is within one column, and flai move changes the column", id, it.Status, ref, other.Status)
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

	ready, backlog := seq, PullSequence(b.Order, items, Backlog)
	if it.Status == Backlog {
		ready, backlog = PullSequence(b.Order, items, Ready), seq
	}
	// Leave the unplaced tail of the backlog unnamed: reading the list back
	// puts unnamed stories last by ID, which is where they already are.
	end := len(backlog)
	for end > 0 {
		x := backlog[end-1]
		if named[x] || (end < len(backlog) && !lessID(x, backlog[end])) {
			break
		}
		end--
	}
	b.Order = append(append([]string{}, ready...), backlog[:end]...)
	return nil
}

// PlaceReadyLast names a story that has just become ready after the ready
// stories the list already names and before any backlog story, wherever the
// list had it before.
func (b *Board) PlaceReadyLast(id string, items []*Item) {
	if items == nil {
		b.AppendToOrder(id)
		return
	}
	b.RemoveFromOrder(id)
	ready := map[string]bool{}
	for _, it := range items {
		if !it.Archived && it.Type == Story && it.Status == Ready && it.ID != id {
			ready[it.ID] = true
		}
	}
	at := 0
	for i, x := range b.Order {
		if ready[x] {
			at = i + 1
		}
	}
	b.Order = append(b.Order[:at], append([]string{id}, b.Order[at:]...)...)
}
