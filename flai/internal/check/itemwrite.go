package check

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Write names what a work item write touches (ADR-0197, ADR-0213, S-0598).
type Write struct {
	// Items are the IDs of the items written. Before a creation the new item
	// does not exist yet, so give its parent.
	Items []string
	// Paths are the files the write makes, slash-separated and relative to
	// the repository's root or its main checkout's, or absolute. A path that
	// does not exist yet is skipped.
	Paths []string
}

// RunItemWrite runs only the rules a work item write can affect (S-0598):
// the item rules on the written items, their parents, and their children,
// with the duplicate-ID pass over every item; the after, overlap, touches,
// component-tag, unaccepted, and epic rules over every item; the board; the
// narratives of the written items' stories, and the index when the write
// wrote it; and the markdown lint and conflict markers of the written files.
// It skips every rule that reads no work item, among them the layout, the
// documentation, the ADR index, the conventions, roles, formulas, issues,
// experiments, analysis, the digest, the design tokens, the component
// examples, which start node, threads, and messages. A checked creation or
// edit runs it before and after the write, in place of Run.
func RunItemWrite(repo *workitem.Repo, now time.Time, w Write) (*Result, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	c := &checker{repo: repo, now: now.UTC(), items: items, byID: map[string]*workitem.Item{}, res: &Result{Items: len(items)}, write: newWriteScope(repo, w)}
	// Run's order, less the rules a write cannot affect
	c.workItems()
	c.narratives()
	c.overlap()
	c.componentTag()
	c.after()
	c.touchesFolder()
	c.afterHub()
	c.afterReason()
	c.unaccepted()
	c.epicLags()
	c.board()
	c.markdown()
	c.conflictMarkers()
	sortFindings(c.res.Findings)
	return c.res, nil
}

// writeScope is what a RunItemWrite reads beyond the item list.
type writeScope struct {
	ids []string
	// paths are the cleaned absolute forms of the paths written, each
	// relative path under the root and under the main checkout; files are
	// those that exist, sorted.
	paths map[string]bool
	files []string
	// items are the IDs oneItem runs on, and stories the stories whose
	// narratives are read, by canonical ID; pick fills them.
	items   map[string]bool
	stories map[string]bool
}

func newWriteScope(repo *workitem.Repo, w Write) *writeScope {
	s := &writeScope{ids: w.Items, paths: map[string]bool{}}
	for _, p := range w.Paths {
		var cands []string
		if filepath.IsAbs(p) {
			cands = []string{filepath.Clean(p)}
		} else {
			for _, root := range []string{repo.Root, repo.MainRoot} {
				if root != "" {
					cands = append(cands, filepath.Join(root, filepath.FromSlash(p)))
				}
			}
		}
		for _, c := range cands {
			if s.paths[c] {
				continue
			}
			s.paths[c] = true
			if info, err := os.Stat(c); err == nil && !info.IsDir() {
				s.files = append(s.files, c)
			}
		}
	}
	sort.Strings(s.files)
	return s
}

// pick works out, once c.byID is built, the items oneItem runs on: those
// written, named or at a path written, their parents, and their children, as
// a child's rules read its parent and report on it; and the stories whose
// narratives are read: each written story, and each written task's story.
func (s *writeScope) pick(c *checker) {
	s.items, s.stories = map[string]bool{}, map[string]bool{}
	written := map[string]*workitem.Item{}
	for _, id := range s.ids {
		if it := c.named(id); it != nil {
			written[it.ID] = it
		}
	}
	for _, it := range c.items {
		if s.paths[filepath.Clean(it.Path)] {
			written[it.ID] = it
		}
	}
	for id, it := range written {
		s.items[id] = true
		parent := c.named(it.Parent)
		if parent != nil {
			s.items[parent.ID] = true
		}
		for _, ch := range workitem.Children(c.items, id) {
			s.items[ch.ID] = true
		}
		switch {
		case it.Type == workitem.Story:
			s.stories[workitem.CanonicalID(id)] = true
		case it.Type == workitem.Task && parent != nil && parent.Type == workitem.Story:
			s.stories[workitem.CanonicalID(parent.ID)] = true
		}
	}
}

// story reports whether the narrative rules read story's narrative: always in
// a whole run, and in a write's only for the stories it wrote.
func (s *writeScope) story(id string) bool {
	return s == nil || s.stories[workitem.CanonicalID(id)]
}

// within reports whether path is inside dir.
func within(path, dir string) bool {
	return strings.HasPrefix(filepath.Clean(path), filepath.Clean(dir)+string(filepath.Separator))
}
