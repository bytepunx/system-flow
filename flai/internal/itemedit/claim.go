package itemedit

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Claims that grow to overlap (I-0059). The pull-time hold keeps a story whose
// claim overlaps an open story's from starting (ADR-0046), but a story's claim
// grows after it starts, as its tasks are written or their touches widened,
// and nothing told the two agents until the close-out found wip.overlap. Every
// write of touches therefore watches the claim of the story it concerns, and
// tells both stories when what the claim gained reaches another story in
// progress. The write always stands: the report is advisory, as wip.overlap
// is (ADR-0019).

// Overlapping is another story in progress whose claim a write reached, and
// the paths the write added to the written story's claim that its claim
// covers.
type Overlapping struct {
	Story string   `json:"story" jsonschema:"the other story in progress"`
	Title string   `json:"title" jsonschema:"its title"`
	Paths []string `json:"paths" jsonschema:"the paths the write added that its claim covers"`
}

// A ClaimWatch is the claim of a story as it stood before a write, for Grown
// to compare with the claim after it.
type ClaimWatch struct {
	repo   *workitem.Repo
	story  string
	before map[string]bool
	err    error // why the claim could not be read, for Grown to return
}

// WatchClaim reads the claim of the story id is, or of the story a task id
// belongs to, before a write to its touches. An epic has no claim: its watch
// reports nothing. A claim that cannot be read is Grown's error, returned
// only once the write has succeeded: a write that names no item fails on its
// own.
func WatchClaim(repo *workitem.Repo, id string) *ClaimWatch {
	w := &ClaimWatch{repo: repo, before: map[string]bool{}}
	it, err := repo.Get(id)
	if err != nil {
		w.err = fmt.Errorf("read %s to watch its story's claim: %w", id, err)
		return w
	}
	switch it.Type {
	case workitem.Story:
		w.story = it.ID
	case workitem.Task:
		w.story = workitem.CanonicalID(it.Parent)
	default:
		return w
	}
	items, err := repo.List(false)
	if err != nil {
		w.err = fmt.Errorf("read the items to watch the claim of %s: %w", w.story, err)
		return w
	}
	holds := workitem.NewHolds(items, repo.Manifest.Projects)
	for _, it := range items {
		if it.ID == w.story {
			for _, p := range holds.Claim(it) {
				w.before[p] = true
			}
		}
	}
	return w
}

// Grown compares the watched story's claim, as the write left it, with each
// other story in progress, and returns those whose claim covers a path the
// claim gained, in ID order. Only the paths gained count, so an overlap the
// hold or an earlier write already allowed is not reported again: a folder
// narrowed to files inside it, as its tasks name them (ADR-0096), gains
// none, and a path whose overlap lies wholly inside a shared path of the
// manifest counts as no overlap. Only a story in progress is compared, and only with stories in
// progress: a story not yet started is held at pull time instead (ADR-0046).
// A story whose claim is empty is not reported: the hold names it, and every
// write would.
//
// Each overlap is recorded twice in the overlap notices, once for each story,
// stamped with by and now, for inbox and wait_for_events to report. Now is
// read after the write, not before it: an agent woken by the item's file may
// look, and move its cursor, in the second between the two.
func (w *ClaimWatch) Grown(by string, now time.Time) ([]Overlapping, error) {
	switch {
	case w == nil:
		return nil, nil
	case w.err != nil:
		return nil, w.err
	case w.story == "":
		return nil, nil
	}
	items, err := w.repo.List(false)
	if err != nil {
		return nil, fmt.Errorf("read the items to compare the claim of %s: %w", w.story, err)
	}
	holds := w.repo.Holds(items)
	var story *workitem.Item
	for _, it := range items {
		if it.Type == workitem.Story && it.ID == w.story {
			story = it
		}
	}
	if story == nil || story.Status != workitem.InProgress {
		return nil, nil
	}
	var gained []string
	for _, p := range holds.Claim(story) {
		if !w.held(p) {
			gained = append(gained, p)
		}
	}
	if len(gained) == 0 {
		return nil, nil
	}
	var out []Overlapping
	for _, it := range items {
		if it.Archived || it.Type != workitem.Story || it.ID == story.ID || it.Status != workitem.InProgress {
			continue
		}
		claim := holds.Claim(it)
		var paths []string
		for _, p := range gained {
			if coveredBy(holds, p, claim) {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			out = append(out, Overlapping{Story: it.ID, Title: it.Title, Paths: paths})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Story < out[j].Story })
	at := now.UTC().Format(workitem.TimeFormat)
	for _, o := range out {
		RecordOverlap(w.repo, Overlap{At: at, By: by, ID: story.ID, Title: story.Title, Grew: story.ID, Reached: o.Story, Paths: o.Paths})
		RecordOverlap(w.repo, Overlap{At: at, By: by, ID: o.Story, Title: o.Title, Grew: story.ID, Reached: o.Story, Paths: o.Paths})
	}
	return out, nil
}

// held says whether the claim before the write held p: p itself or a folder
// around it.
func (w *ClaimWatch) held(p string) bool {
	for {
		if w.before[p] {
			return true
		}
		i := strings.LastIndex(p, "/")
		if i < 0 {
			return false
		}
		p = p[:i]
	}
}

// coveredBy says whether path overlaps an entry of claim outside the shared
// paths holds reads (ADR-0096).
func coveredBy(holds *workitem.Holds, path string, claim []string) bool {
	for _, c := range claim {
		if holds.Overlaps(path, c) {
			return true
		}
	}
	return false
}
