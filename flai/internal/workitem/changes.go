package workitem

import (
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
	ID     string `json:"id"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Kind   string `json:"kind"`
	To     string `json:"to,omitempty"`     // the state moved to
	By     string `json:"by,omitempty"`     // who, when the file says; blocked intervals do not
	Reason string `json:"reason,omitempty"` // why it was blocked
	At     string `json:"at"`
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
	var out []Change
	for _, it := range items {
		base := Change{ID: it.ID, Type: it.Type, Title: it.Title}
		for _, tr := range it.Transitions {
			c := base
			c.Kind, c.To, c.By, c.At = Moved, tr.To, tr.By, tr.At
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
