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
//
// Slow are folders, relative to Root, on a slow tier (S-0606): on every look
// only their folders are statted, and a folder whose modification time or
// size moved is read again and its files statted, so a file added, renamed
// in, or removed is seen at once. Every file below a slow folder is statted
// only every SlowEvery, which is when an edit in place is seen, since it
// moves no folder's modification time. A slow folder below a fast path is
// left out of the fast walk. A second change to a folder within the file
// system's timestamp granularity of the first may wait for the next slow look.
type Watcher struct {
	Root      string
	Paths     []string
	Every     time.Duration // default 300ms
	Slow      []string
	SlowEvery time.Duration // default 10s

	now  func() time.Time                  // the clock; time.Now when nil
	stat func(string) (fs.FileInfo, error) // the slow tier's stat; os.Lstat when nil
	slow *slowTier
}

// slowTier is what the slow tier keeps between looks: each folder's own
// signature and what it held, and the last signature of each file below.
type slowTier struct {
	dirs     map[string]*dirState // by full path
	files    map[string]sig       // by repository-relative path
	lastFull time.Time
	looked   bool
}

type dirState struct {
	sig     sig
	subdirs []string // base names
	files   []string // base names
}

func sigOf(info fs.FileInfo) sig {
	return sig{mod: info.ModTime().UnixNano(), size: info.Size()}
}

// snapshot stats every regular file under the fast paths, and the slow paths
// as the slow tier allows, carrying the files it does not stat.
func (w *Watcher) snapshot() map[string]sig {
	out := map[string]sig{}
	slow := map[string]bool{}
	for _, p := range w.Slow {
		slow[filepath.Clean(filepath.Join(w.Root, p))] = true
	}
	add := func(path string, info fs.FileInfo) {
		if rel, ok := w.rel(path); ok {
			out[rel] = sigOf(info)
		}
	}
	for _, p := range w.Paths {
		full := filepath.Clean(filepath.Join(w.Root, p))
		if slow[full] {
			continue
		}
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
				return nil //nolint:nilerr // a folder that vanished or cannot be read is skipped, not fatal: the next look tries again
			}
			if d.IsDir() {
				if path != full && (strings.HasPrefix(d.Name(), ".") || slow[path]) {
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
	w.walkSlow(out)
	return out
}

func (w *Watcher) rel(path string) (string, bool) {
	rel, err := filepath.Rel(w.Root, path)
	if err != nil {
		return "", false
	}
	return filepath.ToSlash(rel), true
}

// walkSlow adds the files below the slow paths to out: statted when their
// folder changed or the slow interval came round, carried from the last look
// otherwise.
func (w *Watcher) walkSlow(out map[string]sig) {
	if len(w.Slow) == 0 {
		return
	}
	if w.slow == nil {
		w.slow = &slowTier{dirs: map[string]*dirState{}, files: map[string]sig{}}
	}
	t := w.slow
	now := time.Now()
	if w.now != nil {
		now = w.now()
	}
	every := w.SlowEvery
	if every <= 0 {
		every = 10 * time.Second
	}
	full := !t.looked || now.Sub(t.lastFull) >= every
	if full {
		t.lastFull, t.looked = now, true
	}
	seen := map[string]bool{}
	cur := map[string]sig{}
	for _, p := range w.Slow {
		w.lookDir(filepath.Clean(filepath.Join(w.Root, p)), full, seen, cur)
	}
	for dir := range t.dirs {
		if !seen[dir] {
			delete(t.dirs, dir)
		}
	}
	t.files = cur
	for rel, s := range cur {
		out[rel] = s
	}
}

// lookDir stats dir and, when it changed or full is set, reads it again and
// stats its files; otherwise it carries its files' last signatures. Then it
// looks at the folders below it.
func (w *Watcher) lookDir(dir string, full bool, seen map[string]bool, cur map[string]sig) {
	stat := w.stat
	if stat == nil {
		stat = os.Lstat
	}
	t := w.slow
	info, err := stat(dir)
	if err != nil || !info.IsDir() {
		return // gone: its files are reported removed
	}
	seen[dir] = true
	st := t.dirs[dir]
	read := false
	if s := sigOf(info); st == nil || st.sig != s || full {
		if fresh, ok := w.readDir(dir, s, stat, cur); ok {
			st, read = fresh, true
			t.dirs[dir] = st
		} else if st == nil {
			return // unreadable and never read: the next look tries again
		}
	}
	if !read { // unchanged, or unreadable now but read before: carried
		for _, name := range st.files {
			if rel, ok := w.rel(filepath.Join(dir, name)); ok {
				if s, kept := t.files[rel]; kept {
					cur[rel] = s
				}
			}
		}
	}
	for _, sub := range st.subdirs {
		w.lookDir(filepath.Join(dir, sub), full, seen, cur)
	}
}

// readDir reads dir again, statting each regular file into cur, and returns
// what it holds; false when it cannot be read.
func (w *Watcher) readDir(dir string, s sig, stat func(string) (fs.FileInfo, error), cur map[string]sig) (*dirState, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, false
	}
	st := &dirState{sig: s}
	for _, e := range entries {
		name, path := e.Name(), filepath.Join(dir, e.Name())
		if e.IsDir() {
			if !strings.HasPrefix(name, ".") {
				st.subdirs = append(st.subdirs, name)
			}
			continue
		}
		fi, err := stat(path)
		if err != nil || !fi.Mode().IsRegular() {
			continue
		}
		st.files = append(st.files, name)
		if rel, ok := w.rel(path); ok {
			cur[rel] = sigOf(fi)
		}
	}
	return st, true
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
