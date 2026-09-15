package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// Archive moves the planned files into wip/archive.
func (r *Repo) Archive(plan *ArchivePlan) error {
	for _, it := range plan.Items {
		dest := filepath.Join(r.ItemDir(it.Type, true), filepath.Base(it.Path))
		if err := moveFile(it.Path, dest); err != nil {
			return err
		}
		it.Path, it.Archived = dest, true
	}
	for _, src := range plan.Narratives {
		dest := filepath.Join(r.ArchiveDir(), "agents", filepath.Base(src))
		if err := moveFile(src, dest); err != nil {
			return err
		}
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

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
