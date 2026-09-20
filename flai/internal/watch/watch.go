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

// Run calls emit with the repository-relative path of each file that was
// added, changed, or removed, until ctx ends. A file still being written is
// reported once it has looked the same for one tick, which is the debounce.
func (w *Watcher) Run(ctx context.Context, emit func(rel string)) {
	every := w.Every
	if every <= 0 {
		every = 300 * time.Millisecond
	}
	prev := w.snapshot()
	pending := map[string]sig{}
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		cur := w.snapshot()
		var ready []string
		for rel, was := range pending {
			now, exists := cur[rel]
			if !exists || now == was {
				ready = append(ready, rel)
				delete(pending, rel)
			}
		}
		for rel, s := range cur {
			if old, ok := prev[rel]; !ok || old != s {
				pending[rel] = s
			}
		}
		for rel := range prev {
			if _, ok := cur[rel]; !ok {
				if _, waiting := pending[rel]; !waiting {
					ready = append(ready, rel)
				}
			}
		}
		prev = cur
		sort.Strings(ready)
		for i, rel := range ready {
			if i == 0 || ready[i-1] != rel {
				emit(rel)
			}
		}
	}
}
