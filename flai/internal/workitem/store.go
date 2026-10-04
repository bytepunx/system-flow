package workitem

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

// racyWindow is how long after a file's modification time a read of it is
// still not trusted: a second write within the file system's timestamp
// granularity could leave the same time and size (the racy-git problem). A
// file read that soon after it changed is read again on the next look, until
// a read comes later than this. Two seconds covers the coarsest common
// granularity (FAT); on this host's ext4 it re-reads only files written in
// the last two seconds.
const racyWindow = 2 * time.Second

// store keeps the items List and Get parse, by file, so that a long-running
// flai (serve, mcp) reads again only the files that changed (S-0156). It is
// process-wide because the dashboard's methods open a Repo per request. A
// file's item is kept while the file is the same file with the modification
// time and size it had when it was read, and that time was older than the
// read by racyWindow. Items are handed out as copies.
type store struct {
	mu      sync.Mutex
	folders map[string]*folder
}

// folder is one item folder's files. Its lock is held while it is looked at,
// so that callers looking at once share one read of each changed file.
type folder struct {
	mu     sync.Mutex
	files  map[string]*kept // by base name
	parses int              // files parsed, for tests
}

type kept struct {
	info    os.FileInfo
	item    *Item // never handed out; callers get a copy
	trusted bool  // read later than racyWindow after it was modified
}

var parsed = &store{folders: map[string]*folder{}}

func (s *store) folder(dir string) *folder {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := s.folders[dir]
	if f == nil {
		f = &folder{files: map[string]*kept{}}
		s.folders[dir] = f
	}
	return f
}

// list returns a copy of every item in dir, in file name order. A folder that
// does not exist has none.
func (s *store) list(dir string) ([]*Item, error) {
	f := s.folder(dir)
	f.mu.Lock()
	defer f.mu.Unlock()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		clear(f.files)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(entries))
	out := make([]*Item, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		it, err := f.read(dir, name)
		if os.IsNotExist(err) {
			continue // removed since the folder was read
		}
		if err != nil {
			return nil, err
		}
		seen[name] = true
		out = append(out, it)
	}
	for name := range f.files {
		if !seen[name] {
			delete(f.files, name)
		}
	}
	return out, nil
}

// get returns a copy of the item in one file.
func (s *store) get(path string) (*Item, error) {
	f := s.folder(filepath.Dir(path))
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.read(filepath.Dir(path), filepath.Base(path))
}

// read returns a copy of name's item, parsing the file again only when it is
// not the file kept or may have changed since. The caller holds f.mu.
func (f *folder) read(dir, name string) (*Item, error) {
	path := filepath.Join(dir, name)
	info, err := os.Stat(path)
	if err != nil {
		delete(f.files, name)
		return nil, err
	}
	if k := f.files[name]; k != nil && k.trusted && os.SameFile(k.info, info) &&
		k.info.ModTime().Equal(info.ModTime()) && k.info.Size() == info.Size() {
		return k.item.clone(), nil
	}
	readAt := time.Now()
	it, err := ReadItem(path)
	if err != nil {
		delete(f.files, name)
		return nil, err
	}
	f.parses++
	WarnUnknown(path, it.Unknown)
	f.files[name] = &kept{info: info, item: it, trusted: info.ModTime().Before(readAt.Add(-racyWindow))}
	return it.clone(), nil
}

// clone copies the item and everything it refers to, so that a caller that
// changes its copy changes nothing another caller sees.
func (it *Item) clone() *Item {
	c := *it
	c.Transitions = slices.Clone(it.Transitions)
	c.Blocked = slices.Clone(it.Blocked)
	c.Tags = slices.Clone(it.Tags)
	c.Touches = slices.Clone(it.Touches)
	c.Topics = slices.Clone(it.Topics)
	c.After = slices.Clone(it.After)
	c.Unknown = slices.Clone(it.Unknown)
	c.Agent = it.Agent.Clone()
	c.Usage = it.Usage.Clone()
	c.CostOfDelay = it.CostOfDelay.clone()
	c.Forecast = it.Forecast.clone()
	c.Finalized = it.Finalized.clone()
	return &c
}
