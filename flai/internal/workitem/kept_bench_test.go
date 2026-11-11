package workitem

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// keptCorpus is a copy of this repository's item folders, current and
// archived, in a project of its own under a temporary folder, with its kept
// items files written by one list (S-0676, I-0310). The copies are dated an
// hour back, so that every parse is trusted and kept. It skips without the
// monorepo and under -short.
func keptCorpus(b *testing.B) *Repo {
	b.Helper()
	if testing.Short() {
		b.Skip("copies every item of the monorepo")
	}
	src, err := Open(".")
	if err != nil {
		b.Skip("not in a project:", err)
	}
	root := b.TempDir()
	manifest, err := os.ReadFile(filepath.Join(src.Root, "system-flow.yaml"))
	if err != nil {
		b.Skip("no manifest:", err)
	}
	mustWriteB(b, filepath.Join(root, "system-flow.yaml"), manifest)
	if err := os.Mkdir(filepath.Join(root, ".flai-cache"), 0o755); err != nil {
		b.Fatal(err)
	}
	dst, err := Open(root)
	if err != nil {
		b.Fatal(err)
	}
	past := time.Now().Add(-time.Hour)
	n := 0
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			to := dst.ItemDir(typ, archived)
			if err := os.MkdirAll(to, 0o755); err != nil {
				b.Fatal(err)
			}
			files, _ := filepath.Glob(filepath.Join(src.ItemDir(typ, archived), "*.md"))
			for _, f := range files {
				data, err := os.ReadFile(f)
				if err != nil {
					b.Fatal(err)
				}
				path := filepath.Join(to, filepath.Base(f))
				mustWriteB(b, path, data)
				if err := os.Chtimes(path, past, past); err != nil {
					b.Fatal(err)
				}
				n++
			}
		}
	}
	if n == 0 {
		b.Skip("no items found")
	}
	freshStore(b)
	if _, err := dst.List(true); err != nil {
		b.Fatal(err)
	}
	return dst
}

func mustWriteB(b *testing.B, path string, data []byte) {
	b.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		b.Fatal(err)
	}
}

// freshStore gives the package a store with nothing in memory, as a new
// process has, and puts the old one back when the benchmark ends.
func freshStore(b *testing.B) {
	old := parsed
	parsed = newStore()
	b.Cleanup(func() { parsed = old })
}

// keptFolders is each item folder of r with its kept items file's path.
type keptFolderPath struct{ dir, path string }

func keptFolders(r *Repo) []keptFolderPath {
	var out []keptFolderPath
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			dir := key(r.ItemDir(typ, archived))
			out = append(out, keptFolderPath{dir, keptItemsPath(r.CacheDir(), dir)})
		}
	}
	return out
}

// BenchmarkFreshListWithKeptItems times a fresh process's Repo.List(true)
// with the kept items files in place, and each phase of it apart (T-2995).
func BenchmarkFreshListWithKeptItems(b *testing.B) {
	r := keptCorpus(b)
	folders := keptFolders(r)
	var size int64
	for _, f := range folders {
		if info, err := os.Stat(f.path); err == nil {
			size += info.Size()
		}
	}
	b.Logf("kept items files: %.1f MB", float64(size)/1e6)

	// The whole fresh list, as every one-shot flai pays it.
	b.Run("list", func(b *testing.B) {
		for b.Loop() {
			freshStore(b)
			if _, err := r.List(true); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Reading the kept items files, and nothing more.
	b.Run("phase=read", func(b *testing.B) {
		for b.Loop() {
			for _, f := range folders {
				if _, err := os.ReadFile(f.path); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	// readKeptItems: the read and the whole-file decode into keptRecords.
	b.Run("phase=decode-file", func(b *testing.B) {
		for b.Loop() {
			for _, f := range folders {
				if readKeptItems(f.path, f.dir) == nil {
					b.Fatal("no kept items for", f.dir)
				}
			}
		}
	})

	records := map[string]map[string]keptRecord{}
	for _, f := range folders {
		records[f.dir] = readKeptItems(f.path, f.dir)
	}

	// Each item's json.Unmarshal, on one goroutine.
	b.Run("phase=unmarshal-items", func(b *testing.B) {
		for b.Loop() {
			for _, recs := range records {
				for _, rec := range recs {
					var it Item
					if err := json.Unmarshal(rec.Item, &it); err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	})

	items := restoreAll(b, records)

	// Each item's Marshal for readAs, on one goroutine.
	b.Run("phase=marshal-readAs", func(b *testing.B) {
		for b.Loop() {
			for _, it := range items {
				_ = it.Marshal()
			}
		}
	})

	// lookup's stat of every item file, and the folder reads before them.
	b.Run("phase=stat", func(b *testing.B) {
		for b.Loop() {
			for _, f := range folders {
				entries, err := os.ReadDir(f.dir)
				if err != nil {
					b.Fatal(err)
				}
				for _, e := range entries {
					if _, err := os.Stat(filepath.Join(f.dir, e.Name())); err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	})

	// List's clone of every item and its sort by ID.
	b.Run("phase=clone-sort", func(b *testing.B) {
		for b.Loop() {
			out := make([]*Item, len(items))
			for i, it := range items {
				out[i] = it.clone()
			}
			sort.Slice(out, func(i, j int) bool { return lessID(out[i].ID, out[j].ID) })
		}
	})

	// Candidates: a model of lookAll's restore path (read, decode, stat,
	// restore on parsers() workers), with each change T-2995 names.
	gobs := map[string]string{}
	var gobSize int
	for _, f := range folders {
		data := encodeGob(b, f.dir, records[f.dir])
		gobs[f.dir] = strings.TrimSuffix(f.path, ".json") + ".gob"
		mustWriteB(b, gobs[f.dir], data)
		gobSize += len(data)
	}
	b.Logf("the same records as gob: %.1f MB", float64(gobSize)/1e6)
	lines := map[string]string{}
	var linesSize int
	for _, f := range folders {
		data := encodeLines(b, records[f.dir])
		lines[f.dir] = strings.TrimSuffix(f.path, ".json") + ".lines"
		mustWriteB(b, lines[f.dir], data)
		linesSize += len(data)
	}
	b.Logf("the same records a line each: %.1f MB", float64(linesSize)/1e6)
	candidates := []struct {
		name string
		opt  protoOpts
	}{
		{"candidate=json", protoOpts{readAs: true}},
		{"candidate=json-lazy-readAs", protoOpts{}},
		{"candidate=json-folders-at-once", protoOpts{readAs: true, together: true}},
		{"candidate=json-lazy-readAs-folders-at-once", protoOpts{together: true}},
		{"candidate=lines", protoOpts{readAs: true, lines: lines}},
		{"candidate=lines-folders-at-once", protoOpts{readAs: true, together: true, lines: lines}},
		{"candidate=lines-lazy-readAs", protoOpts{lines: lines}},
		{"candidate=lines-lazy-readAs-folders-at-once", protoOpts{together: true, lines: lines}},
		{"candidate=gob", protoOpts{readAs: true, gob: gobs}},
		{"candidate=gob-folders-at-once", protoOpts{readAs: true, together: true, gob: gobs}},
		{"candidate=gob-lazy-readAs", protoOpts{gob: gobs}},
		{"candidate=gob-lazy-readAs-folders-at-once", protoOpts{gob: gobs, together: true}},
	}
	for _, c := range candidates {
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				if n := protoList(b, folders, c.opt); n != len(items) {
					b.Fatalf("restored %d items, want %d", n, len(items))
				}
			}
		})
	}

	// The fresh list and every candidate in turn, round after round, so that
	// a host whose load moves weighs on each alike; each reports its median.
	b.Run("interleaved", func(b *testing.B) {
		times := make([][]time.Duration, len(candidates)+1)
		for b.Loop() {
			start := time.Now()
			freshStore(b)
			if _, err := r.List(true); err != nil {
				b.Fatal(err)
			}
			times[0] = append(times[0], time.Since(start))
			for i, c := range candidates {
				start := time.Now()
				protoList(b, folders, c.opt)
				times[i+1] = append(times[i+1], time.Since(start))
			}
		}
		b.ReportMetric(medianMS(times[0]), "list-ms")
		for i, c := range candidates {
			b.ReportMetric(medianMS(times[i+1]), strings.TrimPrefix(c.name, "candidate=")+"-ms")
		}
	})
}

// medianMS is the median of ds in milliseconds.
func medianMS(ds []time.Duration) float64 {
	s := slices.Clone(ds)
	slices.Sort(s)
	return float64(s[len(s)/2]) / 1e6
}

// restoreAll restores every record as the store would, for the phases that
// need the items.
func restoreAll(b *testing.B, records map[string]map[string]keptRecord) []*Item {
	b.Helper()
	var out []*Item
	for dir, recs := range records {
		for name, rec := range recs {
			it, err := restoreItem(filepath.Join(dir, name), rec)
			if err != nil {
				b.Fatal(err)
			}
			out = append(out, it)
		}
	}
	return out
}

// gobFile and gobRecord are a kept items file as gob would hold it: each item
// decoded into an Item with the file, not kept raw.
type gobFile struct {
	Version int
	Dir     string
	Files   []gobRecord
}

type gobRecord struct {
	Name    string
	Size    int64
	ModTime int64
	Item    *Item
}

func encodeGob(b *testing.B, dir string, recs map[string]keptRecord) []byte {
	b.Helper()
	gf := gobFile{Version: keptVersion, Dir: dir}
	for name, rec := range recs {
		var it Item
		if err := json.Unmarshal(rec.Item, &it); err != nil {
			b.Fatal(err)
		}
		it.Unknown = rec.Unknown
		gf.Files = append(gf.Files, gobRecord{Name: name, Size: rec.Size, ModTime: rec.ModTime, Item: &it})
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(gf); err != nil {
		b.Fatal(err)
	}
	return buf.Bytes()
}

// encodeLines is a folder's records one to a line, the key before the item so
// that a reader matches it without decoding the item:
// name TAB size TAB modtime TAB unknown-JSON TAB item-JSON.
func encodeLines(b *testing.B, recs map[string]keptRecord) []byte {
	b.Helper()
	var buf bytes.Buffer
	buf.WriteString("{\"version\":0}\n")
	for name, rec := range recs {
		unk := []byte("null")
		if len(rec.Unknown) > 0 {
			var err error
			if unk, err = json.Marshal(rec.Unknown); err != nil {
				b.Fatal(err)
			}
		}
		fmt.Fprintf(&buf, "%s\t%d\t%d\t%s\t%s\n", name, rec.Size, rec.ModTime, unk, rec.Item)
	}
	return buf.Bytes()
}

// parseLine splits one record line of encodeLines.
func parseLine(line []byte) (name string, key fileKey, item, unknown []byte, ok bool) {
	parts := bytes.SplitN(line, []byte{'\t'}, 5)
	if len(parts) != 5 {
		return "", fileKey{}, nil, nil, false
	}
	size, err1 := strconv.ParseInt(string(parts[1]), 10, 64)
	mod, err2 := strconv.ParseInt(string(parts[2]), 10, 64)
	if err1 != nil || err2 != nil {
		return "", fileKey{}, nil, nil, false
	}
	if string(parts[3]) == "null" {
		parts[3] = nil
	}
	return string(parts[0]), fileKey{size, mod}, parts[4], parts[3], true
}

type protoOpts struct {
	readAs   bool              // make readAs with Marshal, as restoreItem does
	together bool              // the folders at once, not one after another
	gob      map[string]string // read and decode this gob file, by folder, in place of the JSON file
	lines    map[string]string // read this line-per-record file, by folder, in place of the JSON file
}

// protoList restores every item of folders as a fresh list does when every
// kept entry matches, under opts, and gives how many it restored. It leaves
// out the store's bookkeeping and the write, which a list that matched
// everything skips.
func protoList(b *testing.B, folders []keptFolderPath, o protoOpts) int {
	var mu sync.Mutex
	total := 0
	one := func(f keptFolderPath) {
		n := protoFolder(b, f, o)
		mu.Lock()
		total += n
		mu.Unlock()
	}
	if !o.together {
		for _, f := range folders {
			one(f)
		}
		return total
	}
	var wg sync.WaitGroup
	for _, f := range folders {
		wg.Go(func() { one(f) })
	}
	wg.Wait()
	return total
}

func protoFolder(b *testing.B, f keptFolderPath, o protoOpts) int {
	type entry struct {
		key    fileKey
		raw    json.RawMessage
		unk    []Field
		unkRaw []byte // the line format's unknown fields, decoded with the item
		item   *Item
	}
	recs := map[string]entry{}
	switch {
	case o.gob != nil:
		data, err := os.ReadFile(o.gob[f.dir])
		if err != nil {
			b.Error(err)
			return 0
		}
		var gf gobFile
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&gf); err != nil {
			b.Error(err)
			return 0
		}
		for _, r := range gf.Files {
			recs[r.Name] = entry{key: fileKey{r.Size, r.ModTime}, item: r.Item}
		}
	case o.lines != nil:
		data, err := os.ReadFile(o.lines[f.dir])
		if err != nil {
			b.Error(err)
			return 0
		}
		for rest := data; len(rest) > 0; {
			line := rest
			if i := bytes.IndexByte(rest, '\n'); i >= 0 {
				line, rest = rest[:i], rest[i+1:]
			} else {
				rest = nil
			}
			name, key, raw, unk, ok := parseLine(line)
			if ok {
				recs[name] = entry{key: key, raw: raw, unkRaw: unk}
			}
		}
	default:
		data, err := os.ReadFile(f.path)
		if err != nil {
			b.Error(err)
			return 0
		}
		var kf keptItemsFile
		if err := json.Unmarshal(data, &kf); err != nil {
			b.Error(err)
			return 0
		}
		for _, r := range kf.Files {
			recs[r.Name] = entry{key: fileKey{r.Size, r.ModTime}, raw: r.Item, unk: r.Unknown}
		}
	}
	entries, err := os.ReadDir(f.dir)
	if err != nil {
		b.Error(err)
		return 0
	}
	var todo []entry
	var names []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		info, err := os.Stat(filepath.Join(f.dir, name))
		if err != nil {
			b.Error(err)
			return 0
		}
		if rec, ok := recs[name]; ok && rec.key == keyOf(info) {
			todo = append(todo, rec)
			names = append(names, name)
		}
	}
	restore := func(i int) {
		rec := todo[i]
		it := rec.item
		if it == nil {
			it = &Item{}
			if err := json.Unmarshal(rec.raw, it); err != nil {
				b.Error(err)
				return
			}
			it.Unknown = rec.unk
			if len(rec.unkRaw) > 0 {
				if err := json.Unmarshal(rec.unkRaw, &it.Unknown); err != nil {
					b.Error(err)
					return
				}
			}
		}
		it.Path = filepath.Join(f.dir, names[i])
		if o.readAs {
			it.readAs = it.Marshal()
		}
	}
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(parsers(), max(1, len(todo))) {
		wg.Go(func() {
			for i := range next {
				restore(i)
			}
		})
	}
	for i := range todo {
		next <- i
	}
	close(next)
	wg.Wait()
	return len(todo)
}
