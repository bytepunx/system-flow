package workitem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sort"
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
// with the kept items files in place, each phase of it apart (T-2995), and
// the same list as version 1 of the kept items file made it (T-2998).
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

	// readKeptItems: the read and the split into entries, items undecoded.
	b.Run("phase=split-file", func(b *testing.B) {
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

	// Before S-0676: version 1's one JSON document per folder, decoded whole
	// and the folders restored one after another, modelled on the same items.
	v1 := map[string]string{}
	var v1Size int
	for _, f := range folders {
		data := encodeV1(b, f.dir, records[f.dir])
		v1[f.dir] = strings.TrimSuffix(f.path, ".jsonl") + ".json"
		mustWriteB(b, v1[f.dir], data)
		v1Size += len(data)
	}
	b.Logf("the same items as version 1's JSON: %.1f MB", float64(v1Size)/1e6)
	b.Run("version1", func(b *testing.B) {
		for b.Loop() {
			if n := v1List(b, folders, v1); n != len(items) {
				b.Fatalf("restored %d items, want %d", n, len(items))
			}
		}
	})

	// The fresh list and version 1's model in turn, round after round, so
	// that a host whose load moves weighs on each alike; each reports its
	// median.
	b.Run("interleaved", func(b *testing.B) {
		var list, before []time.Duration
		for b.Loop() {
			start := time.Now()
			freshStore(b)
			if _, err := r.List(true); err != nil {
				b.Fatal(err)
			}
			list = append(list, time.Since(start))
			start = time.Now()
			v1List(b, folders, v1)
			before = append(before, time.Since(start))
		}
		b.ReportMetric(medianMS(list), "list-ms")
		b.ReportMetric(medianMS(before), "version1-ms")
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

// v1File and v1Record are a kept items file as version 1 held it (S-0600):
// one JSON document, each item raw within it.
type v1File struct {
	Version int        `json:"version"`
	Dir     string     `json:"dir"`
	Files   []v1Record `json:"files"`
}

type v1Record struct {
	Name    string          `json:"name"`
	Size    int64           `json:"size"`
	ModTime int64           `json:"mod_time"`
	Item    json.RawMessage `json:"item"`
	Unknown []Field         `json:"unknown,omitempty"`
}

func encodeV1(b *testing.B, dir string, recs map[string]keptRecord) []byte {
	b.Helper()
	vf := v1File{Version: 1, Dir: dir}
	for name, rec := range recs {
		var unknown []Field
		if len(rec.Unknown) > 0 {
			if err := json.Unmarshal(rec.Unknown, &unknown); err != nil {
				b.Fatal(err)
			}
		}
		vf.Files = append(vf.Files, v1Record{Name: name, Size: rec.Size, ModTime: rec.ModTime, Item: rec.Item, Unknown: unknown})
	}
	data, err := json.Marshal(vf)
	if err != nil {
		b.Fatal(err)
	}
	return data
}

// v1List restores every item of folders as a fresh list did with version 1's
// files when every kept entry matched, the folders one after another, and
// gives how many it restored. It leaves out the store's bookkeeping and the
// write, which a list that matched everything skips.
func v1List(b *testing.B, folders []keptFolderPath, files map[string]string) int {
	total := 0
	for _, f := range folders {
		total += v1Folder(b, f, files[f.dir])
	}
	return total
}

func v1Folder(b *testing.B, f keptFolderPath, path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		b.Error(err)
		return 0
	}
	var vf v1File
	if err := json.Unmarshal(data, &vf); err != nil {
		b.Error(err)
		return 0
	}
	recs := make(map[string]v1Record, len(vf.Files))
	for _, r := range vf.Files {
		recs[r.Name] = r
	}
	entries, err := os.ReadDir(f.dir)
	if err != nil {
		b.Error(err)
		return 0
	}
	var todo []v1Record
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
		if rec, ok := recs[name]; ok && (fileKey{rec.Size, rec.ModTime}) == keyOf(info) {
			todo = append(todo, rec)
		}
	}
	next := make(chan v1Record)
	var wg sync.WaitGroup
	for range min(parsers(), max(1, len(todo))) {
		wg.Go(func() {
			for rec := range next {
				var it Item
				if err := json.Unmarshal(rec.Item, &it); err != nil {
					b.Error(err)
					continue
				}
				it.Unknown = rec.Unknown
				it.Path = filepath.Join(f.dir, rec.Name)
				it.readAs = it.Marshal()
			}
		})
	}
	for _, rec := range todo {
		next <- rec
	}
	close(next)
	wg.Wait()
	return len(todo)
}
