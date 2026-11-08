package workitem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

	// loaded is set once the folder's kept items file was read, by its first
	// list in the process; disk is what it held, by base name, for the files
	// not yet listed, and written the keys the file holds now, nil while that
	// is not known (S-0600).
	loaded  bool
	disk    map[string]keptRecord
	written map[string]fileKey
}

type kept struct {
	info    os.FileInfo
	item    *Item           // never handed out; callers get a copy
	trusted bool            // read later than racyWindow after it was modified
	raw     json.RawMessage // item's JSON in the kept items file, once made
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

// maxParsers bounds the workers a list parses a folder's files on (S-0600):
// past 16 the parse of this repository's 3,233 items gained little.
const maxParsers = 16

// parsers is how many workers a list parses on: the CPUs Go may use, at most
// maxParsers.
func parsers() int {
	return max(1, min(runtime.GOMAXPROCS(0), maxParsers))
}

// readItem parses one item file; a test replaces it to watch the pool.
var readItem = ReadItem

// look is one file of a folder as a list looks at it: its stat, and the
// item, the kept one or a fresh parse, or the error that took it.
type look struct {
	name   string
	info   os.FileInfo
	item   *Item
	err    error
	readAt time.Time
	parsed bool // item is a fresh parse, not yet kept

	// rec is the kept items file's entry whose key matches the file's stat,
	// for parseAll to restore in place of a parse; restored is set when it
	// did.
	rec      *keptRecord
	restored bool
}

// list returns a copy of every item in dir, in file name order. A folder
// that does not exist has none. It holds the folder's lock throughout, so
// that lists at once share one parse of each changed file: it stats every
// file, takes the kept item of each whose key still matches, in memory or,
// on the folder's first list in the process, in its kept items file in the
// cache folder cache (S-0600), parses the rest on a bounded pool of
// workers, keeps what they parsed, and writes the kept items file again when
// what it holds changed. With no cache, or a cache folder that does not
// exist, nothing is kept between processes.
func (s *store) list(dir, cache string) ([]*Item, error) {
	f := s.folder(dir)
	f.mu.Lock()
	defer f.mu.Unlock()
	path := keptItemsPath(cache, dir)
	if !f.loaded && path != "" {
		f.loaded = true
		f.disk = readKeptItems(path, dir)
		if f.disk != nil {
			f.written = make(map[string]fileKey, len(f.disk))
			for name, rec := range f.disk {
				f.written[name] = fileKey{size: rec.Size, mod: rec.ModTime}
			}
		}
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		clear(f.files)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	files := make([]*look, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		files = append(files, f.lookup(dir, name))
	}
	parseAll(dir, files)
	seen := make(map[string]bool, len(files))
	out := make([]*Item, 0, len(files))
	for _, l := range files {
		if l.err != nil {
			delete(f.files, l.name)
			if os.IsNotExist(l.err) {
				continue // removed since the folder was read
			}
			return nil, l.err
		}
		if l.parsed || l.restored {
			f.keep(dir, l)
		}
		seen[l.name] = true
		out = append(out, l.item.clone())
	}
	for name := range f.files {
		if !seen[name] {
			delete(f.files, name)
		}
	}
	f.disk = nil // every file now has its entry in memory, or is gone
	f.writeKept(path, dir)
	return out, nil
}

// lookup stats name and gives it the kept item when the key still matches;
// otherwise the file is left for a parse. The caller holds f.mu.
func (f *folder) lookup(dir, name string) *look {
	l := &look{name: name}
	l.info, l.err = os.Stat(filepath.Join(dir, name))
	if l.err != nil {
		return l
	}
	if k := f.files[name]; k != nil && k.trusted && os.SameFile(k.info, l.info) &&
		k.info.ModTime().Equal(l.info.ModTime()) && k.info.Size() == l.info.Size() {
		l.item = k.item
	} else if rec, ok := f.disk[name]; ok && (fileKey{size: rec.Size, mod: rec.ModTime}) == keyOf(l.info) {
		l.rec = &rec
	}
	return l
}

// parseAll parses every file lookup left without an item, on at most
// parsers() workers, restoring instead the item of one whose kept items
// file's entry matched, unless that entry does not decode. It writes only
// each file's own look, so it needs no lock of its own.
func parseAll(dir string, files []*look) {
	var todo []*look
	for _, l := range files {
		if l.err == nil && l.item == nil {
			todo = append(todo, l)
		}
	}
	parse := func(l *look) {
		if l.rec != nil {
			if it, err := restoreItem(filepath.Join(dir, l.name), *l.rec); err == nil {
				l.item, l.restored = it, true
				return
			}
			l.rec = nil // a damaged entry is a miss: parse the file
		}
		l.readAt = time.Now()
		l.item, l.err = readItem(filepath.Join(dir, l.name))
		l.parsed = l.err == nil
	}
	workers := min(parsers(), len(todo))
	if workers <= 1 {
		for _, l := range todo {
			parse(l)
		}
		return
	}
	next := make(chan *look)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for l := range next {
				parse(l)
			}
		})
	}
	for _, l := range todo {
		next <- l
	}
	close(next)
	wg.Wait()
}

// keep records a fresh parse, or an item restored from the kept items file,
// as the folder's item for its file. A parse is trusted once the read came
// racyWindow after the file's modification; a restored item was, since only
// trusted parses are written to the file. The item's unknown fields are
// warned of once in the process, as a parse of them is. The caller holds
// f.mu.
func (f *folder) keep(dir string, l *look) {
	WarnUnknown(filepath.Join(dir, l.name), l.item.Unknown)
	if l.restored {
		f.files[l.name] = &kept{info: l.info, item: l.item, trusted: true, raw: l.rec.Item}
		return
	}
	f.parses++
	// The kept items file's entry for the file, if any, no longer holds:
	// the file changed, or its entry would not decode.
	delete(f.written, l.name)
	f.files[l.name] = &kept{info: l.info, item: l.item, trusted: l.info.ModTime().Before(l.readAt.Add(-racyWindow))}
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
	l := &look{name: name, info: info, readAt: time.Now()}
	if l.item, err = readItem(path); err != nil {
		delete(f.files, name)
		return nil, err
	}
	f.keep(dir, l)
	return l.item.clone(), nil
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
