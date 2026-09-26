package workitem

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Holds on ready stories (S-0128, ADR-0046). A story's claim is its touches
// and those of its tasks that are not done or cancelled. A ready story whose
// claim overlaps the claim of a story in progress or in review is held, and
// so is one that names in after: a story that is not done (S-0130): flai
// serve does not start its agent and wait_for_work does not offer it. The
// operator's own moves only warn.

// Hold reason codes.
const (
	HoldOverlap   = "overlap"    // its claim overlaps an open story's
	HoldNoTouches = "no-touches" // its claim or an open story's is empty
	HoldAfter     = "after"      // a story it names in after: is not done (S-0130)
)

// Hold is why a ready story waits: a reason code, and one line with the
// evidence and what clears it.
type Hold struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

// Holds judges ready stories against the stories open at one reading of the
// repository.
type Holds struct {
	projects []manifest.Project
	tasks    map[string][]*Item // open tasks by story
	open     []openClaim
	stories  map[string]*Item // every story given, archived ones included
	// lookup finds a story named in after: that was not given, such as an
	// archived one when only the active items were read.
	lookup func(id string) *Item
}

type openClaim struct {
	id    string
	label string // how the reason names its state: in progress, in review
	paths []string
}

// NewHolds reads the open stories' claims from items; projects are the
// manifest's sub-projects, whose names and tags a claim reads as their paths.
func NewHolds(items []*Item, projects []manifest.Project) *Holds {
	h := &Holds{projects: projects, tasks: map[string][]*Item{}, stories: map[string]*Item{}}
	for _, it := range items {
		if !it.Archived && it.Type == Task && !it.Closed() {
			h.tasks[it.Parent] = append(h.tasks[it.Parent], it)
		}
		if it.Type == Story {
			h.stories[it.ID] = it
		}
	}
	for _, it := range items {
		if it.Archived || it.Type != Story {
			continue
		}
		switch it.Status {
		case InProgress:
			h.Open(it, "in progress")
		case Review:
			h.Open(it, "in review")
		}
	}
	return h
}

// Holds judges ready stories against items, as NewHolds does, and finds a
// story named in after: in the archive when items do not hold it.
func (r *Repo) Holds(items []*Item) *Holds {
	h := NewHolds(items, r.Manifest.Projects)
	h.lookup = func(id string) *Item {
		if it, err := r.Get(id); err == nil && it.Type == Story {
			return it
		}
		return nil
	}
	return h
}

// Open counts story as open from now on, named with label: flai serve's
// launcher opens a story whose agent it has started before the agent moves
// it to in progress, so that one look does not start two that overlap.
func (h *Holds) Open(story *Item, label string) {
	h.open = append(h.open, openClaim{id: story.ID, label: label, paths: h.Claim(story)})
	sort.SliceStable(h.open, func(i, j int) bool { return h.open[i].id < h.open[j].id })
}

// Claim is what story claims: its touches and those of its open tasks, each
// a sub-project's path when it names the sub-project or one of its tags,
// without duplicates.
func (h *Holds) Claim(story *Item) []string {
	var out []string
	seen := map[string]bool{}
	add := func(entries []string) {
		for _, e := range entries {
			p := h.path(e)
			if p != "" && !seen[p] {
				seen[p] = true
				out = append(out, p)
			}
		}
	}
	add(story.Touches)
	for _, t := range h.tasks[story.ID] {
		add(t.Touches)
	}
	return out
}

// path is a touches entry as it is compared: a component's name or tag
// becomes its path, and anything else is taken as it is.
func (h *Holds) path(entry string) string {
	e := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(entry), "./"), "/")
	for _, p := range h.projects {
		if e == p.Name || contains(p.Tags, e) {
			return strings.TrimSuffix(p.Path, "/")
		}
	}
	return e
}

// Of is why story is held, or nil when it is not. Every story that holds it
// is named, in ID order. A story held both by after: and by an overlap is
// held (after), and its reason says both.
func (h *Holds) Of(story *Item) *Hold {
	after, overlap := h.after(story), h.overlap(story)
	switch {
	case after == nil:
		return overlap
	case overlap == nil:
		return after
	}
	return &Hold{Code: HoldAfter, Reason: after.Reason + "; also " + overlap.Reason}
}

// after is why story waits for the stories it names in after:, or nil when
// each of them is done. One that was cancelled keeps the hold, and so does
// one that does not exist: flai check reports it, and the operator decides.
func (h *Holds) after(story *Item) *Hold {
	var waits, ids, notes []string
	seen := map[string]bool{}
	for _, e := range story.After {
		id := CanonicalID(e)
		if id == story.ID || seen[id] {
			continue
		}
		seen[id] = true
		it := h.story(id)
		switch {
		case it == nil:
			waits = append(waits, id+" (no such story)")
			notes = append(notes, id+" names no story, so fix after:")
		case it.Status == Done:
			continue
		case it.Status == Cancelled:
			waits = append(waits, id+" (cancelled)")
			notes = append(notes, id+" was cancelled, so drop it from after: if "+story.ID+" no longer needs it")
		default:
			waits = append(waits, id+" ("+stateLabel(it.Status)+")")
		}
		ids = append(ids, id)
	}
	if len(waits) == 0 {
		return nil
	}
	reason := fmt.Sprintf("held (after): waits for %s; starts when %s %s done", and(waits), and(ids), oneOrMany(len(ids), "is", "are"))
	if len(notes) > 0 {
		reason += "; " + strings.Join(notes, "; ")
	}
	return &Hold{Code: HoldAfter, Reason: reason}
}

// story finds a story by ID among those given, else through lookup.
func (h *Holds) story(id string) *Item {
	if it, ok := h.stories[id]; ok {
		return it
	}
	if h.lookup == nil {
		return nil
	}
	it := h.lookup(id)
	h.stories[id] = it
	return it
}

// stateLabel is how a reason names a story's state.
func stateLabel(status string) string {
	switch status {
	case Backlog:
		return "in backlog"
	case InProgress:
		return "in progress"
	case Review:
		return "in review"
	}
	return status
}

// overlap is why story's claim is held by the open stories', or nil.
func (h *Holds) overlap(story *Item) *Hold {
	claim := h.Claim(story)
	var others []openClaim
	for _, o := range h.open {
		if o.id != story.ID {
			others = append(others, o)
		}
	}
	if len(others) == 0 {
		return nil
	}
	if len(claim) == 0 {
		var named []string
		for _, o := range others {
			named = append(named, o.id+" ("+o.label+")")
		}
		return &Hold{Code: HoldNoTouches, Reason: fmt.Sprintf("held (no-touches): declares no touches, so it may change what %s %s; starts when it declares touches that overlap no open story's, or when %s", and(named), oneOrMany(len(others), "changes", "change"), clears(others))}
	}
	var code string
	var parts []string
	var by []openClaim
	for _, o := range others {
		part, c := holdBy(claim, o)
		if part == "" {
			continue
		}
		if code == "" {
			code = c
		}
		parts = append(parts, part)
		by = append(by, o)
	}
	if len(parts) == 0 {
		return nil
	}
	return &Hold{Code: code, Reason: fmt.Sprintf("held (%s): %s; starts when %s", code, strings.Join(parts, "; "), clears(by))}
}

// holdBy says how claim overlaps o's, if it does, naming the first pair.
func holdBy(claim []string, o openClaim) (part, code string) {
	who := o.id + " (" + o.label + ")"
	if len(o.paths) == 0 {
		return who + " declares no touches, so it may change anything", HoldNoTouches
	}
	for _, mine := range claim {
		for _, theirs := range o.paths {
			switch {
			case !PathsOverlap(mine, theirs):
			case mine == theirs:
				return fmt.Sprintf("touches %s, which %s touches too", mine, who), HoldOverlap
			case strings.HasPrefix(mine, theirs+"/"):
				return fmt.Sprintf("touches %s, inside %s which %s touches", mine, theirs, who), HoldOverlap
			default:
				return fmt.Sprintf("touches %s, which holds %s that %s touches", mine, theirs, who), HoldOverlap
			}
		}
	}
	return "", ""
}

// clears says what ends a hold by these stories.
func clears(by []openClaim) string {
	ids := make([]string, len(by))
	for i, o := range by {
		ids[i] = o.id
	}
	return fmt.Sprintf("%s %s accepted, cancelled, or sent back", and(ids), oneOrMany(len(ids), "is", "are"))
}

// PathsOverlap says whether two touches entries cover a common path: equal,
// or one a /-bounded prefix of the other (ADR-0019).
func PathsOverlap(a, b string) bool {
	a, b = strings.TrimSuffix(a, "/"), strings.TrimSuffix(b, "/")
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
}

func and(xs []string) string {
	switch len(xs) {
	case 0:
		return ""
	case 1:
		return xs[0]
	default:
		return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
	}
}

func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
