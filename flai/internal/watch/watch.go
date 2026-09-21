// Package watch tells flai serve which files of a project changed, so the
// dashboard can be told (ADR-0029). It polls: a walk of three folders of
// Markdown a few times a second costs little, behaves the same on every
// platform and filesystem (inotify, kqueue, and Windows each have their
// limits and surprises), needs no dependency, and can be tested with a clock.
package watch

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type sig struct {
	mod  int64
	size int64
}

// Watcher reports changes under Paths, which are relative to Root: folders
// are walked, a file is watched by itself.
type Watcher struct {
	Root  string
	Paths []string
	Every time.Duration // default 300ms
}

// snapshot stats every regular file under the watched paths.
func (w *Watcher) snapshot() map[string]sig {
	out := map[string]sig{}
	add := func(path string, info fs.FileInfo) {
		rel, err := filepath.Rel(w.Root, path)
		if err != nil {
			return
		}
		out[filepath.ToSlash(rel)] = sig{mod: info.ModTime().UnixNano(), size: info.Size()}
	}
	for _, p := range w.Paths {
		full := filepath.Join(w.Root, p)
		info, err := os.Stat(full)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			add(full, info)
			continue
		}
		_ = filepath.WalkDir(full, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// A folder that vanished or cannot be read is skipped, not fatal: the next look tries again.
				return nil //nolint:nilerr
			}
			if d.IsDir() {
				if path != full && strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if info, err := d.Info(); err == nil && info.Mode().IsRegular() {
				add(path, info)
			}
			return nil
		})
	}
	return out
}

// debounce holds what one tick needs from the ticks before it: the last
// snapshot, and the files that changed but have not yet looked the same for a
// whole tick. It has no clock and touches no files, so a test can feed it
// snapshots in any order and at any pace.
type debounce struct {
	prev    map[string]sig
	pending map[string]sig
}

func newDebounce(first map[string]sig) *debounce {
	return &debounce{prev: first, pending: map[string]sig{}}
}

// tick takes the next snapshot and returns, sorted and without repeats, the
// files that are ready to report: those that have looked the same since the
// last tick after a change, and those that were removed.
func (d *debounce) tick(cur map[string]sig) []string {
	var ready []string
	for rel, was := range d.pending {
		now, exists := cur[rel]
		if !exists || now == was {
			ready = append(ready, rel)
			delete(d.pending, rel)
		}
	}
	for rel, s := range cur {
		if old, ok := d.prev[rel]; !ok || old != s {
			d.pending[rel] = s
		}
	}
	for rel := range d.prev {
		if _, ok := cur[rel]; !ok {
			if _, waiting := d.pending[rel]; !waiting {
				ready = append(ready, rel)
			}
		}
	}
	d.prev = cur
	sort.Strings(ready)
	var out []string
	for _, rel := range ready {
		if len(out) == 0 || out[len(out)-1] != rel {
			out = append(out, rel)
		}
	}
	return out
}

// Run calls emit with the repository-relative path of each file that was
// added, changed, or removed, until ctx ends. A file still being written is
// reported once it has looked the same for one tick, which is the debounce.
func (w *Watcher) Run(ctx context.Context, emit func(rel string)) {
	every := w.Every
	if every <= 0 {
		every = 300 * time.Millisecond
	}
	d := newDebounce(w.snapshot())
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		for _, rel := range d.tick(w.snapshot()) {
			emit(rel)
		}
	}
}
