package workitem

import (
	"slices"
	"sort"
	"time"
)

// Kinds of change to a work item that the files record with a time.
const (
	Moved      = "moved"
	WasBlocked = "blocked"
	Unblocked  = "unblocked"
)

// Change is something that happened to a work item, read back from its
// front matter: a transition, or the start or end of a blocked interval.
// Nothing is recorded to produce it (S-0058).
type Change struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	To      string `json:"to,omitempty"`      // the state moved to
	By      string `json:"by,omitempty"`      // who, when the file says; blocked intervals do not
	Reason  string `json:"reason,omitempty"`  // why it was blocked
	Cause   string `json:"cause,omitempty"`   // the item whose cancellation took this one with it (ADR-0028)
	Follows string `json:"follows,omitempty"` // the story whose move this epic followed (S-0200)
	At      string `json:"at"`
}

// Key identifies a change, for telling apart changes within one second.
func (c Change) Key() string { return c.ID + "|" + c.Kind + "|" + c.To + "|" + c.At }

// Changes returns what happened after since, oldest first, leaving out
// transitions made by self. Timestamps have second resolution, so a change
// stamped with the very second of since is included unless seen names it:
// the caller keeps the keys it has already reported for that second.
func Changes(items []*Item, since time.Time, self string, seen map[string]bool) []Change {
	since = since.UTC().Truncate(time.Second)
	after := func(at string, c Change) bool {
		t, err := time.Parse(TimeFormat, at)
		if err != nil {
			return false
		}
		return t.After(since) || (t.Equal(since) && !seen[c.Key()])
	}
	byID := map[string]*Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	var out []Change
	for _, it := range items {
		base := Change{ID: it.ID, Type: it.Type, Title: it.Title}
		for _, tr := range it.Transitions {
			c := base
			c.Kind, c.To, c.By, c.At = Moved, tr.To, tr.By, tr.At
			if tr.To == Cancelled {
				c.Cause = CancelledWith(byID, it, tr.At)
			} else if it.Type == Epic {
				c.Follows = FollowedWith(byID, it, tr)
			}
			if tr.By != self && after(tr.At, c) {
				out = append(out, c)
			}
		}
		for _, b := range it.Blocked {
			c := base
			c.Kind, c.Reason, c.At = WasBlocked, b.Reason, b.From
			if after(b.From, c) {
				out = append(out, c)
			}
			if b.Until != "" {
				u := base
				u.Kind, u.At = Unblocked, b.Until
				if after(b.Until, u) {
					out = append(out, u)
				}
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].At != out[j].At {
			return out[i].At < out[j].At
		}
		return lessID(out[i].ID, out[j].ID)
	})
	return out
}

// Reordered reports whether the stories two pull orders share come in a
// different sequence. Stories entering and leaving the order are moves,
// reported as such (flai move appends a ready story and drops a started
// one); only a change of priority among them is news of its own. The MCP
// server's inbox and flai serve's replanner (S-0211, ADR-0084) both ask it.
func Reordered(before, after []string) bool {
	in := func(list []string) map[string]bool {
		m := map[string]bool{}
		for _, id := range list {
			m[id] = true
		}
		return m
	}
	inBefore, inAfter := in(before), in(after)
	var a, b []string
	for _, id := range before {
		if inAfter[id] {
			a = append(a, id)
		}
	}
	for _, id := range after {
		if inBefore[id] {
			b = append(b, id)
		}
	}
	return !slices.Equal(a, b)
}

// CancelledWith names the item whose cancellation at the same moment took it
// along: the highest ancestor cancelled at that time. Nothing records it; a
// cascade stamps every item with one time (ADR-0028).
func CancelledWith(byID map[string]*Item, it *Item, at string) string {
	cause := ""
	for p := byID[it.Parent]; p != nil; p = byID[p.Parent] {
		for _, tr := range p.Transitions {
			if tr.To == Cancelled && tr.At == at {
				cause = p.ID
			}
		}
		if p.Parent == p.ID {
			break
		}
	}
	return cause
}

// FollowedWith names the story whose move an epic's transition followed: a
// child story moved at the same time by the same actor, the lowest ID if
// several. Nothing records it; a story and its epic share one stamp (S-0200).
func FollowedWith(byID map[string]*Item, it *Item, tr Transition) string {
	found := ""
	for _, c := range byID {
		if c.Parent != it.ID || c.Type != Story || (found != "" && !lessID(c.ID, found)) {
			continue
		}
		for _, st := range c.Transitions {
			if st.At == tr.At && st.By == tr.By {
				found = c.ID
				break
			}
		}
	}
	return found
}
