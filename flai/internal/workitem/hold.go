package workitem

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Holds on ready stories (S-0128, ADR-0046). A story's claim is its touches,
// each folder among them narrowed to the touches its tasks name inside it,
// and the touches of its open tasks outside them (ADR-0096). A ready story whose
// claim overlaps the claim of a story in progress is held, and so is one that
// names in after: a story that is not done (S-0130): flai serve does not start
// its agent and wait_for_work does not offer it. A story in review holds
// nothing (S-0295, ADR-0096): its branch is finished and synced, and the
// notice at acceptance tells an overlapping story what changed. An overlap
// that lies wholly inside a pattern of the manifest's claims.shared holds
// nothing either (ADR-0096): many stories change those paths in separate
// sections or new files. An overlap between a ready story and the story in
// progress that holds it, on paths that story's agent shared with it in their
// conversation, holds nothing while the share is in force (ADR-0134, share.go).
// The operator's own moves only warn.

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
	shared   manifest.Claims    // an overlap wholly inside these holds nothing
	shares   []Share            // shares in force, IDs canonical (ADR-0134)
	tasks    map[string][]*Item // tasks by story, cancelled ones left out
	open     []openClaim
	stories  map[string]*Item // every story given, archived ones included
	// lookup finds a story named in after: that was not given, such as an
	// archived one when only the active items were read.
	lookup func(id string) *Item
}

type openClaim struct {
	id    string
	label string // how the reason names its state: in progress, its agent started
	paths []string
}

// NewHolds reads from items the claims of the stories in progress, the only
// stories that hold (ADR-0096); projects are the manifest's sub-projects,
// whose names and tags a claim reads as their paths.
func NewHolds(items []*Item, projects []manifest.Project) *Holds {
	h := &Holds{projects: projects, tasks: map[string][]*Item{}, stories: map[string]*Item{}}
	for _, it := range items {
		if !it.Archived && it.Type == Task && it.Status != Cancelled {
			h.tasks[it.Parent] = append(h.tasks[it.Parent], it)
		}
		if it.Type == Story {
			h.stories[it.ID] = it
		}
	}
	for _, it := range items {
		if !it.Archived && it.Type == Story && it.Status == InProgress {
			h.Open(it, "in progress")
		}
	}
	return h
}

// Holds judges ready stories against items, as NewHolds does, with the
// manifest's shared paths as SharedClaims reads them and the shares in force
// of the open conversations (ADR-0134), and finds a story named in after: or
// in a share in the archive when items do not hold it.
func (r *Repo) Holds(items []*Item) *Holds {
	h := NewHolds(items, r.Manifest.Projects).WithShared(r.SharedClaims())
	h.lookup = func(id string) *Item {
		if it, err := r.Get(id); err == nil && it.Type == Story {
			return it
		}
		return nil
	}
	return h.WithShares(r.Shares())
}

// SharedClaims is the manifest's claims as system-flow.yaml has them now,
// read again rather than taken from the manifest r was opened with: flai mcp
// keeps one Repo for as long as it runs, and flai shared or shared_paths_edit
// may have changed the list since. The manifest r was opened with gives them
// when the file cannot be read.
func (r *Repo) SharedClaims() manifest.Claims {
	m, err := manifest.Load(filepath.Join(r.Root, manifest.File))
	if err != nil {
		return r.Manifest.Claims
	}
	return m.Claims
}

// WithShared has h treat an overlap wholly inside a pattern of claims as no
// overlap (ADR-0096), and returns h.
func (h *Holds) WithShared(claims manifest.Claims) *Holds {
	h.shared = claims
	return h
}

// Overlaps says whether paths a and b, such as entries of two claims, overlap
// where it counts: they cover a common path, as PathsOverlap says, and the
// narrower of the two, the deeper one or either when they are equal, does not
// lie wholly inside a shared path (ADR-0096).
func (h *Holds) Overlaps(a, b string) bool {
	if !PathsOverlap(a, b) {
		return false
	}
	_, shared := h.shared.Covers(narrower(a, b))
	return !shared
}

// overlapsFor is Overlaps for mine, an entry of held's claim, and theirs, an
// entry of holder's: a pair whose narrower lies inside a path holder shared
// with held does not overlap either (ADR-0134).
func (h *Holds) overlapsFor(held, holder, mine, theirs string) bool {
	return h.Overlaps(mine, theirs) && !h.sharedBy(holder, held, narrower(mine, theirs))
}

// narrower is the deeper of two overlapping paths, the second when they are
// equal: one is the other or lies below it.
func narrower(a, b string) string {
	if len(strings.TrimSuffix(a, "/")) > len(strings.TrimSuffix(b, "/")) {
		return a
	}
	return b
}

// Open counts story as open from now on, named with label: flai serve's
// launcher opens a story whose agent it has started before the agent moves
// it to in progress, so that one look does not start two that overlap.
func (h *Holds) Open(story *Item, label string) {
	h.open = append(h.open, openClaim{id: story.ID, label: label, paths: h.Claim(story)})
	sort.SliceStable(h.open, func(i, j int) bool { return h.open[i].id < h.open[j].id })
}

// Claim is what story claims (ADR-0096), each entry read as a path first, a
// sub-project's when it names the sub-project or one of its tags, without
// duplicates. A story touch that holds touches of its tasks, done or open, is
// replaced by them; one that no task names inside, or that a task names whole,
// stays whole; and a touch of an open task outside every story touch is added.
// A done task's touch counts only where it narrows: the branch changed it.
func (h *Holds) Claim(story *Item) []string {
	var out []string
	seen := map[string]bool{}
	add := func(p string) {
		if p != "" && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	var mine, theirs, open []string // the story's touches, its tasks', its open tasks'
	for _, e := range story.Touches {
		mine = append(mine, h.path(e))
	}
	for _, t := range h.tasks[story.ID] {
		for _, e := range t.Touches {
			p := h.path(e)
			theirs = append(theirs, p)
			if !t.Closed() {
				open = append(open, p)
			}
		}
	}
	for _, s := range mine {
		var inside []string
		whole := false
		for _, p := range theirs {
			switch {
			case p == "" || !PathsOverlap(p, s):
			case strings.HasPrefix(p, s+"/"):
				inside = append(inside, p)
			default:
				whole = true // the task names the touch itself, or a folder around it
			}
		}
		if whole || len(inside) == 0 {
			add(s)
			continue
		}
		for _, p := range inside {
			add(p)
		}
	}
	for _, p := range open {
		if !within(p, mine) {
			add(p)
		}
	}
	return out
}

// within says whether p is one of entries or lies below one.
func within(p string, entries []string) bool {
	for _, e := range entries {
		if e != "" && (p == e || strings.HasPrefix(p, e+"/")) {
			return true
		}
	}
	return false
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
		return &Hold{Code: HoldNoTouches, Reason: fmt.Sprintf("held (no-touches): declares no touches, so it may change what %s %s; starts when it declares touches that overlap no story's in progress, or when %s", and(named), oneOrMany(len(others), "changes", "change"), clears(others))}
	}
	var code string
	var parts []string
	var by []openClaim
	for _, o := range others {
		part, c := h.holdBy(story.ID, claim, o)
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

// holdBy says how claim, held's, overlaps o's, if it does, naming the first
// pair that holds: a pair whose overlap lies wholly inside a shared path, or
// inside a path o shared with held, does not, so the pair named is one still
// held.
func (h *Holds) holdBy(held string, claim []string, o openClaim) (part, code string) {
	who := o.id + " (" + o.label + ")"
	if len(o.paths) == 0 {
		return who + " declares no touches, so it may change anything", HoldNoTouches
	}
	pairs := h.holding(held, claim, o)
	if len(pairs) == 0 {
		return "", ""
	}
	mine, theirs := pairs[0][0], pairs[0][1]
	switch {
	case mine == theirs:
		return fmt.Sprintf("touches %s, which %s touches too", mine, who), HoldOverlap
	case strings.HasPrefix(mine, theirs+"/"):
		return fmt.Sprintf("touches %s, inside %s which %s touches", mine, theirs, who), HoldOverlap
	default:
		return fmt.Sprintf("touches %s, which holds %s that %s touches", mine, theirs, who), HoldOverlap
	}
}

// holding is every pair of an entry of claim, held's, and one of o's that
// holds, in claim order, then o's.
func (h *Holds) holding(held string, claim []string, o openClaim) [][2]string {
	var out [][2]string
	for _, mine := range claim {
		for _, theirs := range o.paths {
			if h.overlapsFor(held, o.id, mine, theirs) {
				out = append(out, [2]string{mine, theirs})
			}
		}
	}
	return out
}

// clears says what ends a hold by these stories: each leaving in progress.
func clears(by []openClaim) string {
	ids := make([]string, len(by))
	for i, o := range by {
		ids[i] = o.id
	}
	if len(ids) == 1 {
		return ids[0] + " moves to review, is cancelled, or is sent back"
	}
	return and(ids) + " move to review, are cancelled, or are sent back"
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
