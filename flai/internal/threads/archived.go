package threads

import (
	"fmt"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// OnItems returns the threads still open or answered whose anchor is one of
// the items, matched in any padding, in thread ID order.
func OnItems(r *workitem.Repo, ids []string) ([]*Thread, error) {
	want := map[string]bool{}
	for _, id := range ids {
		want[workitem.CanonicalID(id)] = true
	}
	all, err := List(r)
	if err != nil {
		return nil, err
	}
	var out []*Thread
	for _, th := range all {
		if th.Open() && th.Anchor.Item != "" && want[workitem.CanonicalID(th.Anchor.Item)] {
			out = append(out, th)
		}
	}
	return out, nil
}

// ResolveOnItems resolves the threads OnItems returns for the items, each
// with a dated entry giving the reason, mirrors the narratives of their
// stories, and returns the resolved threads' IDs.
func ResolveOnItems(r *workitem.Repo, ids []string, author, reason string, now time.Time) ([]string, error) {
	open, err := OnItems(r, ids)
	if err != nil {
		return nil, err
	}
	var resolved []string
	stories := map[string]bool{}
	for _, th := range open {
		if _, err := Resolve(r, th.ID, author, reason, now); err != nil {
			return resolved, fmt.Errorf("resolve %s on %s: %w; resolve it with flai thread resolve %s", th.ID, th.Anchor.Item, err, th.ID)
		}
		resolved = append(resolved, th.ID)
		if s := StoryOf(r, th); s != "" {
			stories[s] = true
		}
	}
	// An archived story's narrative has left wip/agents, and MirrorNarrative
	// skips it; a task archived alone may leave its story's narrative live.
	var ss []string
	for s := range stories {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	for _, s := range ss {
		if err := MirrorNarrative(r, s); err != nil {
			return resolved, fmt.Errorf("mirror the threads into the narrative of %s: %w", s, err)
		}
	}
	return resolved, nil
}
