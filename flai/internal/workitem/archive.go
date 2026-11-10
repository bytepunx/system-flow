package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// ArchivePlan lists what Archive would move.
type ArchivePlan struct {
	Items      []*Item
	Narratives []string // source paths
}

// PlanArchive selects items to archive. With no IDs, every closed epic,
// closed story with its tasks, and closed task whose story is not active is
// selected. With IDs, each must be closed; a story brings its tasks.
func (r *Repo) PlanArchive(items []*Item, ids []string) (*ArchivePlan, error) {
	plan := &ArchivePlan{}
	byID := map[string]*Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	selected := map[string]bool{}
	add := func(it *Item) {
		if !selected[it.ID] && !it.Archived {
			selected[it.ID] = true
			plan.Items = append(plan.Items, it)
		}
	}
	if len(ids) == 0 {
		for _, it := range items {
			if it.Archived || !it.Closed() {
				continue
			}
			switch it.Type {
			case Epic:
				open := false
				for _, c := range Children(items, it.ID) {
					if !c.Archived {
						open = true
					}
				}
				if !open {
					add(it)
				}
			case Story:
				add(it)
				for _, c := range Children(items, it.ID) {
					add(c)
				}
			case Task:
				if p, ok := byID[it.Parent]; !ok || p.Archived || p.Closed() {
					add(it)
				}
			}
		}
	} else {
		for _, id := range ids {
			it, ok := byID[id]
			if !ok {
				return nil, fmt.Errorf("%s not found", id)
			}
			if it.Archived {
				return nil, fmt.Errorf("%s is already archived", id)
			}
			if !it.Closed() {
				return nil, fmt.Errorf("%s is %s; only done or cancelled items can be archived", id, it.Status)
			}
			add(it)
			if it.Type == Story {
				for _, c := range Children(items, it.ID) {
					if !c.Closed() {
						return nil, fmt.Errorf("%s has task %s in %s", it.ID, c.ID, c.Status)
					}
					add(c)
				}
			}
			if it.Type == Epic {
				for _, c := range Children(items, it.ID) {
					if !c.Archived && !selected[c.ID] {
						return nil, fmt.Errorf("%s has story %s still on the board; archive it first", it.ID, c.ID)
					}
				}
			}
		}
	}
	for _, it := range plan.Items {
		if it.Type == Story {
			if p := r.NarrativePath(it.ID); fileExists(p) {
				plan.Narratives = append(plan.Narratives, p)
			}
		}
	}
	sort.Slice(plan.Items, func(i, j int) bool { return lessID(plan.Items[i].ID, plan.Items[j].ID) })
	return plan, nil
}

// Archive moves the planned files into wip/archive. The items move under the
// items lock, so that a Save waiting on it reads after the move and refuses a
// copy read before it with ErrItemMoved rather than writing it back where it
// was (I-0323).
func (r *Repo) Archive(plan *ArchivePlan) error {
	if err := r.archiveItems(plan.Items); err != nil {
		return err
	}
	for _, src := range plan.Narratives {
		dest := filepath.Join(r.ArchiveDir(), "agents", filepath.Base(src))
		if err := moveNarrative(src, dest); err != nil {
			return err
		}
	}
	return nil
}

// archiveItems moves each item's file into wip/archive under the items lock.
func (r *Repo) archiveItems(items []*Item) error {
	if len(items) == 0 {
		return nil
	}
	unlock, err := r.lockItems()
	if err != nil {
		return fmt.Errorf("archive %s and the rest of the plan: nothing was moved: %w", items[0].ID, err)
	}
	defer unlock()
	for _, it := range items {
		dest := filepath.Join(r.ItemDir(it.Type, true), filepath.Base(it.Path))
		err := moveFile(it.Path, dest)
		// this process's own move is seen by its next list of either folder (S-0606)
		parsed.stale(filepath.Dir(it.Path))
		parsed.stale(filepath.Dir(dest))
		if err != nil {
			return err
		}
		it.Path, it.Archived = dest, true
	}
	return nil
}

func moveFile(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if fileExists(dest) {
		return fmt.Errorf("%s already exists", dest)
	}
	return os.Rename(src, dest)
}

// moveNarrative moves a narrative into the archive with each run of blank
// lines outside fenced code collapsed to one: no story owns an archived
// narrative, so a markdownlint MD012 finding a hand edit left in it would
// stay for good (I-0139).
func moveNarrative(src, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if fileExists(dest) {
		return fmt.Errorf("%s already exists", dest)
	}
	if err := atomicfile.WriteFile(dest, []byte(collapseBlankRuns(string(data))), 0o644); err != nil {
		return err
	}
	return os.Remove(src)
}

var fenceOpen = regexp.MustCompile("^ {0,3}(`{3,}|~{3,})")

// collapseBlankRuns leaves one blank line where text has two or more in a
// row, except inside a fenced code block, and every other line as it was.
func collapseBlankRuns(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	fence := ""
	for i, l := range lines {
		if fence != "" {
			if t := strings.TrimSpace(l); strings.HasPrefix(t, fence) && strings.Trim(t, fence[:1]) == "" {
				fence = ""
			}
		} else if m := fenceOpen.FindStringSubmatch(l); m != nil {
			fence = m[1]
		} else if i > 0 && strings.TrimSpace(l) == "" && strings.TrimSpace(lines[i-1]) == "" {
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
