package itemedit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ForecastChange is one story's forecast edit in a batch.
type ForecastChange struct {
	ID   string
	Edit ForecastEdit
}

// Batch is what ApplyForecasts wrote and what it did not.
type Batch struct {
	// Written are the IDs of the items written, in the order of the changes.
	Written []string `json:"written"`
	// Files are every file written, relative to the checkout: the items'
	// and wip/agents/index.md.
	Files []string `json:"files"`
	// Refused holds, by ID, why each change that was not written was not.
	Refused map[string]error `json:"-"`
	// Warnings are the check's findings on the written items that refused
	// nothing.
	Warnings []check.Finding `json:"warnings"`
}

// runItems is the check ApplyForecasts runs over what it wrote, writeIndex
// its write of wip/agents/index.md, and beforeSave runs between its read and
// its save: variables, so that a test can count the check's runs and the
// index's writes and bring another writer's save in between.
var (
	runItems   = check.Items
	writeIndex = (*workitem.Repo).WriteIndex
	beforeSave = func() {}
)

// ApplyForecasts writes many stories' forecast edits in one pass (S-0599):
// one list of the items, one save of all of them under one hold of the items
// lock, one write of wip/agents/index.md, and one check of the work-item
// rules and the markdown lint over the items written. Each item gets the
// bytes Apply would write for its change. A change is refused alone, and the
// others written, when its item is missing, archived, or closed, when the
// edit leaves a valid item not valid, when another writer changed its
// forecast since the read (ADR-0184), or when the check finds an error on
// its own file, which puts that item back. It commits nothing.
func ApplyForecasts(repo *workitem.Repo, changes []ForecastChange, opt Options) (*Batch, error) {
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	res := &Batch{Written: []string{}, Files: []string{}, Refused: map[string]error{}, Warnings: []check.Finding{}}
	if len(changes) == 0 {
		return res, nil
	}
	all, err := repo.List(true)
	if err != nil {
		return nil, fmt.Errorf("list the items to replan their forecasts: %w", err)
	}
	byID := map[string]*workitem.Item{}
	for _, it := range all {
		byID[it.ID] = it
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	var edited []*workitem.Item
	before := map[string][]byte{}
	for _, ch := range changes {
		id := workitem.CanonicalID(ch.ID)
		it, ok := byID[id]
		if !ok {
			res.Refused[id] = fmt.Errorf("%s: %w", id, workitem.ErrNotFound)
			continue
		}
		if ok, reason := editable(it); !ok {
			res.Refused[id] = &docedit.RefusedError{Path: rel(repo, it.Path), Reason: reason}
			continue
		}
		data, err := os.ReadFile(it.Path)
		if err != nil {
			res.Refused[id] = fmt.Errorf("read %s: %w", rel(repo, it.Path), err)
			continue
		}
		wasValid := it.Validate() == nil
		e := ch.Edit
		changed, err := applyPlanning(it, Change{Forecast: &e}, opt, repo.Manifest.Planning.CurrencyCode())
		if err != nil {
			res.Refused[id] = err
			continue
		}
		if len(changed) == 0 {
			continue
		}
		if wasValid {
			if err := it.Validate(); err != nil {
				res.Refused[id] = invalid("%s: %v", it.ID, err)
				continue
			}
		}
		it.Body = fmt.Sprintf("# %s %s\n\n%s", it.ID, it.Title, strings.TrimLeft(below(it.Body), "\n"))
		it.Updated = now
		before[it.Path] = data
		edited = append(edited, it)
	}
	if len(edited) == 0 {
		return res, nil
	}

	beforeSave()
	saved, errs := repo.SaveAll(edited)
	for id, err := range errs {
		res.Refused[id] = err
	}
	if len(saved) == 0 {
		return res, nil
	}
	wrote := map[string][]byte{}
	for _, it := range saved {
		wrote[it.Path] = []byte(it.Marshal())
	}
	files := map[string]bool{}
	if items, err := repo.List(false); err == nil {
		index := filepath.Join(repo.AgentsDir(), "index.md")
		if err := writeIndex(repo, items, opt.Now); err != nil {
			// as Apply does, a failed index write puts back what was saved
			for _, it := range saved {
				if rerr := repo.Restore(it.Path, wrote[it.Path], before[it.Path]); rerr != nil {
					err = fmt.Errorf("%w; and putting %s back failed: %w", err, rel(repo, it.Path), rerr)
				}
			}
			return nil, fmt.Errorf("write %s after replanning %d forecasts: %w", rel(repo, index), len(saved), err)
		}
		files[rel(repo, index)] = true
	}

	// ---- the work-item rules and the lint over what was written: an error
	// on an item's own file puts that item back, as Apply's check does. The
	// saved items are the list's own copies, so the rules read them as
	// saved. ----
	after := runItems(repo, all, saved, opt.Now)
	refused := map[string][]check.Finding{}
	own := map[string]*workitem.Item{}
	for _, it := range saved {
		if p, err := filepath.Rel(repo.Root, it.Path); err == nil {
			own[filepath.ToSlash(p)] = it
		}
	}
	for _, f := range after.Findings {
		it := ownerOf(own, f.Path)
		if it != nil && f.Level == check.Error {
			refused[it.ID] = append(refused[it.ID], f)
			continue
		}
		res.Warnings = append(res.Warnings, f)
	}
	for _, it := range saved {
		if fs, bad := refused[it.ID]; bad {
			reason := fmt.Sprintf("flai check has %d finding(s) with this forecast; it was not written", len(fs))
			if err := repo.Restore(it.Path, wrote[it.Path], before[it.Path]); err != nil {
				reason += "; and putting the file back failed: " + err.Error()
				files[rel(repo, it.Path)] = true
			}
			res.Refused[it.ID] = &docedit.RefusedError{Path: rel(repo, it.Path), Reason: reason, Findings: fs}
			continue
		}
		res.Written = append(res.Written, it.ID)
		files[rel(repo, it.Path)] = true
		note(repo, Notice{At: now, By: opt.By, ID: it.ID, Title: it.Title, Type: it.Type, Changed: []string{"forecast"}})
	}
	for f := range files {
		res.Files = append(res.Files, f)
	}
	sort.Strings(res.Files)
	return res, nil
}

// ownerOf is the written item whose file a finding's path names, by its path
// from the project root or, as Apply matches, by its file name.
func ownerOf(own map[string]*workitem.Item, path string) *workitem.Item {
	p := filepath.ToSlash(path)
	if it, ok := own[p]; ok {
		return it
	}
	for o, it := range own {
		if strings.HasSuffix(p, "/"+filepath.Base(o)) {
			return it
		}
	}
	return nil
}
