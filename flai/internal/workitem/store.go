package workitem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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
//
// A root whose files a watcher looks at for this process, as flai serve's
// does, may keep a snapshot of each of its item folders (S-0606): a list of a
// fresh folder answers from what the folder last listed, reading and statting
// nothing. The folder goes stale when the store is told of a change in it, or
// when this process writes in it, and its next list looks at it as above.
type store struct {
	mu      sync.Mutex
	folders map[string]*folder // by cleaned absolute path
	roots   map[string]bool    // cleaned absolute roots whose folders keep a snapshot
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

	looks int // times the folder was read and its files statted, for tests

	marks atomic.Uint64 // changes the store was told of, so a look that one raced stays stale
	// changes counts the looks and reads that found a file's item added,
	// changed, or removed; with marks it is the folder's generation (S-0608).
	changes atomic.Uint64

	smu      sync.Mutex // guards the snapshot below, never held across a look
	on       bool       // its root keeps a snapshot
	snap     []*Item    // what the folder last listed, in file name order; never handed out
	have     bool       // snap holds a look
	snapOK   bool       // every file of that look was trusted
	snapMark uint64     // marks when that look began

	// bad is the files the folder's latest look could not read, in file name
	// order, whether or not it keeps a snapshot; smu guards it too.
	bad []Unreadable
}

// Unreadable is an item file a list skipped because it could not be stat'd,
// read, or parsed, and the error that took it (I-0344).
type Unreadable struct {
	Path string // the file, under the folder as the list was given it
	Err  error  // the error the stat, read, or parse gave
}

type kept struct {
	info    os.FileInfo
	item    *Item           // never handed out; callers get a copy
	trusted bool            // read later than racyWindow after it was modified
	raw     json.RawMessage // item's JSON in the kept items file, once made
}

var parsed = &store{folders: map[string]*folder{}, roots: map[string]bool{}}

// SnapshotItems makes the item folders under root keep a snapshot between
// lists, for a process that tells the store of every change under them with
// ItemChanged, as flai serve's watcher does.
func SnapshotItems(root string) { parsed.snapshot(root, true) }

// StopSnapshotItems drops root's snapshots, so its lists look at the folders
// again as every flai process does.
func StopSnapshotItems(root string) { parsed.snapshot(root, false) }

// ItemChanged tells the store a file or folder changed: rel is relative to
// root, as the watcher reports it, or absolute. The item folder that holds it
// is looked at again on its next list. A path outside the item folders is
// ignored.
func ItemChanged(root, rel string) {
	path := rel
	if !filepath.IsAbs(path) {
		path = filepath.Join(key(root), filepath.FromSlash(rel))
	}
	parsed.stale(filepath.Dir(path))
	parsed.stale(path)
}

// ItemsChanged makes every item folder under root be looked at again on its
// next list, as after another process wrote items there.
func ItemsChanged(root string) {
	r := key(root)
	parsed.mu.Lock()
	defer parsed.mu.Unlock()
	for dir, f := range parsed.folders {
		if under(dir, r) {
			f.marks.Add(1)
		}
	}
}

// key is dir as the store keys it: cleaned and absolute.
func key(dir string) string {
	if abs, err := filepath.Abs(dir); err == nil {
		return abs
	}
	return filepath.Clean(dir)
}

// under reports whether dir is root or inside it; both are keys.
func under(dir, root string) bool {
	return dir == root || strings.HasPrefix(dir, root+string(filepath.Separator))
}

func (s *store) snapshot(root string, on bool) {
	r := key(root)
	s.mu.Lock()
	defer s.mu.Unlock()
	if on {
		s.roots[r] = true
	} else {
		delete(s.roots, r)
	}
	for dir, f := range s.folders {
		if under(dir, r) {
			f.setOn(s.keeps(dir))
		}
	}
}

// keeps reports whether dir is under a root that keeps a snapshot. The caller
// holds s.mu.
func (s *store) keeps(dir string) bool {
	for r := range s.roots {
		if under(dir, r) {
			return true
		}
	}
	return false
}

// stale marks the folder at dir, when the store has one, to be looked at
// again.
func (s *store) stale(dir string) {
	s.mu.Lock()
	f := s.folders[key(dir)]
	s.mu.Unlock()
	if f != nil {
		f.marks.Add(1)
	}
}

// generation is the folder at dir's count of the changes the store was told
// of, this process's writes among them, and of the items its looks and reads
// found added, changed, or removed. It never goes back.
func (s *store) generation(dir string) uint64 {
	f := s.folder(dir)
	return f.marks.Load() + f.changes.Load()
}

func (s *store) folder(dir string) *folder {
	k := key(dir)
	s.mu.Lock()
	defer s.mu.Unlock()
	f := s.folders[k]
	if f == nil {
		f = &folder{files: map[string]*kept{}, on: s.keeps(k)}
		s.folders[k] = f
	}
	return f
}

func (f *folder) setOn(on bool) {
	f.smu.Lock()
	defer f.smu.Unlock()
	f.on = on
	if !on {
		f.snap, f.have = nil, false
	}
}

// fresh returns the snapshot when the folder keeps one and nothing it was
// told of, or wrote, came since the look that made it.
func (f *folder) fresh() ([]*Item, bool) {
	f.smu.Lock()
	defer f.smu.Unlock()
	if f.on && f.have && f.snapOK && f.snapMark == f.marks.Load() {
		return f.snap, true
	}
	return nil, false
}

// keepSnapshot records the files a look could not read and, when the folder
// keeps a snapshot, the look that began at mark as the snapshot.
func (f *folder) keepSnapshot(items []*Item, bad []Unreadable, mark uint64, trusted bool) {
	f.smu.Lock()
	defer f.smu.Unlock()
	f.bad = bad
	if f.on {
		f.snap, f.have, f.snapOK, f.snapMark = items, true, trusted, mark
	}
}

// items returns the kept items of dir, in file name order: the snapshot when
// it is fresh, else what a look at the folder finds (with cache as list takes
// it). They are never changed and never handed out; callers copy what they
// hand on.
func (s *store) items(dir, cache string) ([]*Item, error) {
	return s.itemsOn(dir, cache, nil)
}

// itemsOn is items, its parses taking slots, a list's bound shared by the
// folders it lists at once (S-0676); with none, the look makes its own. The
// files the look skipped are recorded as the folder's unreadable ones; a
// fresh snapshot has none, since a look that skipped a file is not trusted.
func (s *store) itemsOn(dir, cache string, slots chan struct{}) ([]*Item, error) {
	f := s.folder(dir)
	if items, ok := f.fresh(); ok {
		return items, nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if items, ok := f.fresh(); ok {
		return items, nil // another caller looked while this one waited
	}
	mark := f.marks.Load()
	items, bad, trusted, err := f.lookAll(dir, cache, slots)
	if err != nil {
		return nil, err
	}
	f.keepSnapshot(items, bad, mark, trusted)
	return items, nil
}

// unreadable returns the files the latest look of the folder at dir could not
// read; none when the store never looked at it.
func (s *store) unreadable(dir string) []Unreadable {
	s.mu.Lock()
	f := s.folders[key(dir)]
	s.mu.Unlock()
	if f == nil {
		return nil
	}
	f.smu.Lock()
	defer f.smu.Unlock()
	return slices.Clone(f.bad)
}

// snapshotting reports whether dir keeps a snapshot, so that Get finds its
// item in what the folder listed rather than by a glob.
func (s *store) snapshotting(dir string) bool {
	f := s.folder(dir)
	f.smu.Lock()
	defer f.smu.Unlock()
	return f.on
}

// find returns a copy of the first item in dir, by file name, whose file name
// starts with prefix, as a glob of prefix*.md would; nil when none does. When
// only a file the folder could not read has such a name, its error is
// returned, as a get of that file would.
func (s *store) find(dir, prefix, cache string) (*Item, error) {
	items, err := s.items(dir, cache)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if strings.HasPrefix(filepath.Base(it.Path), prefix) {
			return it.clone(), nil
		}
	}
	for _, u := range s.unreadable(dir) {
		if strings.HasPrefix(filepath.Base(u.Path), prefix) {
			return nil, u.Err
		}
	}
	return nil, nil
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
// that does not exist has none.
func (s *store) list(dir, cache string) ([]*Item, error) {
	return s.listOn(dir, cache, nil)
}

// listOn is list, its parses taking slots as itemsOn's do.
func (s *store) listOn(dir, cache string, slots chan struct{}) ([]*Item, error) {
	items, err := s.itemsOn(dir, cache, slots)
	if err != nil {
		return nil, err
	}
	out := make([]*Item, len(items))
	for i, it := range items {
		out[i] = it.clone()
	}
	return out, nil
}

// lookAll looks at every file of dir and returns the kept items, never
// copied, the files it could not stat, read, or parse, and whether the look
// is trusted: every item read late enough to be trusted and no file skipped.
// A file that could not be read is left out and the rest are listed
// (I-0344); only a folder that cannot be read fails the look. The caller
// holds the folder's lock throughout, so that lists at once share one parse
// of each changed file: it stats every
// file, takes the kept item of each whose key still matches, in memory or,
// on the folder's first list in the process, in its kept items file in the
// cache folder cache (S-0600), parses the rest on a bounded pool of
// workers, each parse taking one of slots, keeps what they parsed, and writes
// the kept items file again when what it holds changed. With no cache, or a
// cache folder that does not exist, nothing is kept between processes.
func (f *folder) lookAll(dir, cache string, slots chan struct{}) ([]*Item, []Unreadable, bool, error) {
	f.looks++
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
		for name := range f.files {
			f.drop(name)
		}
		return nil, nil, true, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	files := make([]*look, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		files = append(files, f.lookup(dir, name))
	}
	parseAll(dir, files, slots)
	seen := make(map[string]bool, len(files))
	out := make([]*Item, 0, len(files))
	var bad []Unreadable
	trusted := true
	for _, l := range files {
		if l.err != nil {
			f.drop(l.name)
			if os.IsNotExist(l.err) {
				continue // removed since the folder was read
			}
			bad = append(bad, Unreadable{Path: filepath.Join(dir, l.name), Err: l.err})
			trusted = false // look again, so that a fix is seen
			continue
		}
		if l.parsed || l.restored {
			f.keep(dir, l)
		}
		trusted = trusted && f.files[l.name].trusted
		seen[l.name] = true
		out = append(out, l.item)
	}
	for name := range f.files {
		if !seen[name] {
			f.drop(name)
		}
	}
	f.disk = nil // every file now has its entry in memory, or is gone
	f.writeKept(path, dir)
	return out, bad, trusted, nil
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

// parseSlots is a bound of parsers() parses at once, for the folders of one
// list to share.
func parseSlots() chan struct{} { return make(chan struct{}, parsers()) }

// parseAll parses every file lookup left without an item, restoring instead
// the item of one whose kept items file's entry matched, unless that entry
// does not decode. Each parse or restore takes one of slots, so that at most
// parsers() run at once across the folders that share them; with no slots,
// the folder has its own. It writes only each file's own look, so it needs
// no lock of its own.
func parseAll(dir string, files []*look, slots chan struct{}) {
	var todo []*look
	for _, l := range files {
		if l.err == nil && l.item == nil {
			todo = append(todo, l)
		}
	}
	if slots == nil {
		slots = parseSlots()
	}
	parse := func(l *look) {
		slots <- struct{}{}
		defer func() { <-slots }()
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
	// A file read again only because it was too recent to trust, and read
	// the same, is no change.
	if old := f.files[l.name]; old == nil || old.item.readAs != l.item.readAs {
		f.changes.Add(1)
	}
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
		f.drop(name)
		return nil, err
	}
	if k := f.files[name]; k != nil && k.trusted && os.SameFile(k.info, info) &&
		k.info.ModTime().Equal(info.ModTime()) && k.info.Size() == info.Size() {
		return k.item.clone(), nil
	}
	l := &look{name: name, info: info, readAt: time.Now()}
	if l.item, err = readItem(path); err != nil {
		f.drop(name)
		return nil, err
	}
	f.keep(dir, l)
	return l.item.clone(), nil
}

// drop forgets name's kept item, counting a change when there was one. The
// caller holds f.mu.
func (f *folder) drop(name string) {
	if _, ok := f.files[name]; ok {
		delete(f.files, name)
		f.changes.Add(1)
	}
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
