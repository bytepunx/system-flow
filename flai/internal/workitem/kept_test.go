package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// richDoc is an item with every kind of front matter the kept items file must
// restore: lists empty and full, a block, an agent with roles, usage with
// turns and strategic entries, a cost of delay whose one stamp ReadItem
// moves onto its inputs, a forecast, a finalized draft, a formula, and a key
// this flai does not know.
const richDoc = `---
id: S-0007
type: story
nature: improvement
title: "Rich: every field"
status: done
parent: E-0001
owner: alex
created: 2026-10-16T05:40:19Z
updated: 2026-10-31T06:23:24Z
transitions:
  - to: ready
    at: 2026-10-28T13:05:18Z
    by: alex
  - to: done
    at: 2026-10-31T06:23:24Z
    by: alex
blocked:
  - from: 2026-10-30T00:56:44Z
    until: 2026-10-30T01:57:41Z
    reason: awaiting the operator on TH-0760
  - from: 2026-10-30T13:34:42Z
    reason: still waiting
tags: []
topics: [cli, planning]
touches: [flai/internal/workitem/store.go, design/system/flai-cli.md]
after: [S-0006]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
  roles:
    explore:
      model: claude-haiku-5-5
usage:
  source: log
  seconds: 5900
  priced_by: harness
  turns:
    - day: 2026-10-30
      ceremony: 9
      work: 127
  models:
    - model: claude-opus-5-5
      input: 438
      output: 119148
      cache_read: 40340969
      cache_write: 1069270
      cost: 17.9763
  strategic:
    - kind: planner
      seconds: 120
      estimated: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
  value: 12.5
  by: planner-S-0251
  at: 2026-10-05T05:52:47Z
forecast:
  duration: 1h
  delivery: 2026-10-30T11:45:00Z
  basis: "Its own forecast: 1h"
  by: flai
  at: 2026-10-30T00:12:18Z
finalized:
  by: alex
  at: 2026-10-28T13:05:00Z
formula: publish
colour: |
  blue
  and green
---

# S-0007 Rich: every field

## Goal

Everything.
`

// itemDoc is a plain item with ID id.
func itemDoc(id, title string) string {
	return fmt.Sprintf("---\nid: %s\ntype: epic\nnature: feature\ntitle: %s\nstatus: backlog\nowner: alex\ncreated: 2026-10-01T00:00:00Z\nupdated: 2026-10-01T00:00:00Z\ntransitions: []\ntags: []\n---\n\n# %s %s\n", id, title, id, title)
}

// keptFolder makes a project folder with an item folder of n epics and, with
// cache, its .flai-cache folder, each file modified an hour ago so that its
// parse is trusted. It writes the files itself, so that no save makes the
// cache folder.
func keptFolder(t *testing.T, n int, cache bool) (dir, cacheDir string) {
	t.Helper()
	root := t.TempDir()
	dir = filepath.Join(root, "wip", "kanban", "epics")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	cacheDir = filepath.Join(root, ".flai-cache")
	if cache {
		if err := os.Mkdir(cacheDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	past := time.Now().Add(-time.Hour)
	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("E-%04d", i)
		path := filepath.Join(dir, id+"-epic.md")
		if err := os.WriteFile(path, []byte(itemDoc(id, fmt.Sprintf("Epic %d", i))), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, past, past); err != nil {
			t.Fatal(err)
		}
	}
	return dir, cacheDir
}

// newStore is a store with nothing in memory, as a new process has.
func newStore() *store { return &store{folders: map[string]*folder{}} }

// coldList lists dir with a new store and gives its items and how many files
// it parsed.
func coldList(t *testing.T, dir, cache string) ([]*Item, int) {
	t.Helper()
	s := newStore()
	items, err := s.list(dir, cache)
	if err != nil {
		t.Fatal(err)
	}
	return items, s.folder(dir).parses
}

// readAll is every item file of dir as ReadItem reads it, in name order.
func readAll(t *testing.T, dir string) []*Item {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	var out []*Item
	for _, f := range files {
		it, err := ReadItem(f)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, it)
	}
	return out
}

func sameItems(t *testing.T, got, want []*Item) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("listed %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("item %d:\n got %#v\nwant %#v", i, got[i], want[i])
		}
	}
}

func TestASecondProcessListsWithoutParsingUnchangedFiles(t *testing.T) {
	dir, cache := keptFolder(t, 12, true)
	first, n := coldList(t, dir, cache)
	if n != 12 {
		t.Fatalf("the first process parsed %d files, want 12", n)
	}
	if _, err := os.Stat(keptItemsPath(cache, dir)); err != nil {
		t.Fatalf("the first process kept no items file: %v", err)
	}
	second, n := coldList(t, dir, cache)
	if n != 0 {
		t.Errorf("a second process parsed %d unchanged files", n)
	}
	sameItems(t, second, first)
	sameItems(t, second, readAll(t, dir))
}

func TestAColdListParsesOnlyTheFileThatChanged(t *testing.T) {
	dir, cache := keptFolder(t, 6, true)
	coldList(t, dir, cache)
	past := time.Now().Add(-30 * time.Minute)

	grown := filepath.Join(dir, "E-0002-epic.md")
	if err := os.WriteFile(grown, []byte(itemDoc("E-0002", "Epic two, longer")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(grown, past, past); err != nil {
		t.Fatal(err)
	}
	items, n := coldList(t, dir, cache)
	if n != 1 {
		t.Errorf("after one file grew, a new process parsed %d files, want 1", n)
	}
	if items[1].Title != "Epic two, longer" {
		t.Errorf("the grown file lists as %q", items[1].Title)
	}

	// The same bytes, only a later modification time.
	touched := filepath.Join(dir, "E-0005-epic.md")
	later := past.Add(time.Minute)
	if err := os.Chtimes(touched, later, later); err != nil {
		t.Fatal(err)
	}
	items, n = coldList(t, dir, cache)
	if n != 1 {
		t.Errorf("after one file's time moved, a new process parsed %d files, want 1", n)
	}
	sameItems(t, items, readAll(t, dir))
}

// keptLines is the kept items file at path, a line to an entry.
func keptLines(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// spoilLine replaces the n-th line, the header being 0, of the kept items
// file at path with what edit makes of it.
func spoilLine(t *testing.T, path string, n int, edit func(string) string) {
	t.Helper()
	lines := keptLines(t, path)
	lines[n] = edit(lines[n])
	mustWrite(t, path, strings.Join(lines, "\n")+"\n")
}

func TestAMissingDamagedForeignOrOlderKeptFileIsAMiss(t *testing.T) {
	for _, tc := range []struct {
		name  string
		spoil func(t *testing.T, path string)
		want  int // files parsed, of 5
	}{
		{"missing", func(t *testing.T, path string) { mustRemove(t, path) }, 5},
		{"not a header", func(t *testing.T, path string) {
			spoilLine(t, path, 0, func(string) string { return "not json at all" })
		}, 5},
		{"another folder", func(t *testing.T, path string) {
			spoilLine(t, path, 0, func(l string) string { return strings.Replace(l, `"dir":"`, `"dir":"/elsewhere`, 1) })
		}, 5},
		{"another version", func(t *testing.T, path string) {
			spoilLine(t, path, 0, func(l string) string {
				return strings.Replace(l, fmt.Sprintf(`"version":%d`, keptVersion), `"version":999`, 1)
			})
		}, 5},
		// Cut inside the third entry: the two before it still restore.
		{"truncated", func(t *testing.T, path string) {
			lines := keptLines(t, path)
			mustWrite(t, path, strings.Join(lines[:3], "\n")+"\n"+lines[3][:len(lines[3])/2])
		}, 3},
		{"an entry that does not decode as an item", func(t *testing.T, path string) {
			spoilLine(t, path, 1, func(l string) string { return strings.Replace(l, "\t{", "\t{\"id\":7,", 1) })
		}, 1},
		// I-0310: an entry that is not JSON at all spoils only its own file.
		// Version 1 was one JSON document, so one broken entry made the whole
		// file a miss, and reading it decoded every entry, matching or not.
		{"an entry that is not JSON", func(t *testing.T, path string) {
			spoilLine(t, path, 2, func(l string) string { return l[:strings.LastIndex(l, "\t")] + "\t{\"id\":" })
		}, 1},
		{"a line that is not an entry", func(t *testing.T, path string) {
			spoilLine(t, path, 4, func(l string) string { return strings.Replace(l, "\t", " ", 2) })
		}, 1},
		{"an unknown field that does not decode", func(t *testing.T, path string) {
			spoilLine(t, path, 5, func(l string) string { return strings.Replace(l, "\tnull\t", "\t[{\"Name\":7}]\t", 1) })
		}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, cache := keptFolder(t, 5, true)
			coldList(t, dir, cache)
			tc.spoil(t, keptItemsPath(cache, dir))
			items, n := coldList(t, dir, cache)
			sameItems(t, items, readAll(t, dir))
			if n != tc.want {
				t.Errorf("parsed %d files, want %d", n, tc.want)
			}
			if _, n := coldList(t, dir, cache); n != 0 {
				t.Errorf("the kept file was not written again whole: the next process parsed %d", n)
			}
		})
	}
}

// A kept items file version 1 wrote, one JSON document beside the new file's
// name, is never read: the first list parses every file, writes the new
// file, and leaves the old one for an older flai that may still use it.
func TestAVersion1KeptFileIsAMissAndTheNextListWritesTheNewOne(t *testing.T) {
	dir, cache := keptFolder(t, 4, true)
	path := keptItemsPath(cache, dir)
	if !strings.HasSuffix(path, ".jsonl") {
		t.Fatalf("the kept items file is %s, want a .jsonl", path)
	}
	old := strings.TrimSuffix(path, ".jsonl") + ".json"
	v1 := fmt.Sprintf(`{"version":1,"dir":%q,"files":[]}`, dir)
	mustWrite(t, old, v1)
	if _, n := coldList(t, dir, cache); n != 4 {
		t.Errorf("with only a version 1 file, a new process parsed %d files, want 4", n)
	}
	if got := readKeptItems(path, dir); len(got) != 4 {
		t.Errorf("the list wrote %d entries in the new file, want 4", len(got))
	}
	if data, err := os.ReadFile(old); err != nil || string(data) != v1 {
		t.Errorf("the version 1 file was changed: %q, %v", data, err)
	}
	if _, n := coldList(t, dir, cache); n != 0 {
		t.Errorf("a second process parsed %d files", n)
	}
}

// A name that needs escaping in JSON, a tab or a quote in it, keeps its line
// whole and restores.
func TestAKeptEntryNamedWithATabOrAQuoteRestores(t *testing.T) {
	dir, cache := keptFolder(t, 2, true)
	for _, name := range []string{"E-0003-a\ttab.md", `E-0004-a"quote.md`} {
		path := filepath.Join(dir, name)
		mustWrite(t, path, itemDoc(name[:6], "Odd"))
		past := time.Now().Add(-time.Hour)
		if err := os.Chtimes(path, past, past); err != nil {
			t.Fatal(err)
		}
	}
	coldList(t, dir, cache)
	items, n := coldList(t, dir, cache)
	if n != 0 {
		t.Errorf("parsed %d files, want every one restored", n)
	}
	sameItems(t, items, readAll(t, dir))
}

func TestAProjectWithoutACacheFolderKeepsNothing(t *testing.T) {
	dir, cache := keptFolder(t, 4, false)
	items, n := coldList(t, dir, cache)
	if n != 4 {
		t.Errorf("parsed %d files, want 4", n)
	}
	sameItems(t, items, readAll(t, dir))
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		t.Errorf("the list made %s: %v", cache, err)
	}
	if _, n := coldList(t, dir, cache); n != 4 {
		t.Errorf("a second process parsed %d files with nothing kept, want 4", n)
	}
}

func TestTheKeptFileIsWrittenWholeThroughARename(t *testing.T) {
	dir, cache := keptFolder(t, 30, true)
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			items, err := newStore().list(dir, cache)
			if err != nil || len(items) != 30 {
				t.Errorf("a list at once gave %d items, %v", len(items), err)
			}
		})
	}
	wg.Wait()
	left, err := filepath.Glob(filepath.Join(cache, ".items-*"))
	if err != nil || len(left) != 0 {
		t.Errorf("temporary files left in the cache folder: %v %v", left, err)
	}
	if got := readKeptItems(keptItemsPath(cache, dir), dir); len(got) != 30 {
		t.Errorf("the kept file holds %d entries, want 30", len(got))
	}
	if _, n := coldList(t, dir, cache); n != 0 {
		t.Errorf("after lists at once, a new process parsed %d files", n)
	}
}

func TestAFileWrittenTooRecentlyToTellIsNotKept(t *testing.T) {
	dir, cache := keptFolder(t, 3, true)
	fresh := filepath.Join(dir, "E-0003-epic.md")
	now := time.Now()
	if err := os.Chtimes(fresh, now, now); err != nil {
		t.Fatal(err)
	}
	coldList(t, dir, cache)
	if _, n := coldList(t, dir, cache); n != 1 {
		t.Errorf("a new process parsed %d files, want only the one written too recently to trust", n)
	}
}

func TestKeptItemsRestoreAsReadItemReadsThem(t *testing.T) {
	dir, cache := keptFolder(t, 1, true)
	path := filepath.Join(dir, "S-0007-rich.md")
	mustWrite(t, path, richDoc)
	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatal(err)
	}
	want := readAll(t, dir)
	if len(want[1].Unknown) == 0 || want[1].CostOfDelay.Inputs.By == "" {
		t.Fatal("the fixture no longer carries an unknown key and a moved stamp")
	}
	coldList(t, dir, cache)
	got, n := coldList(t, dir, cache)
	if n != 0 {
		t.Fatalf("parsed %d files, want them restored", n)
	}
	sameItems(t, got, want)
	if got[1].Marshal() != want[1].Marshal() || got[1].readAs != want[1].readAs {
		t.Errorf("the restored item writes back as\n%s\nnot\n%s", got[1].Marshal(), want[1].Marshal())
	}
}

// Every item of this repository restores from its kept items file as
// ReadItem reads it: each folder's items are written to a kept items file as
// a list writes them, read back, and restored.
func TestKeptItemsRestoreThisRepositorysItems(t *testing.T) {
	if testing.Short() {
		t.Skip("reads every item of the monorepo")
	}
	r, err := Open(".")
	if err != nil {
		t.Skip("not in a project:", err)
	}
	n := 0
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			dir := r.ItemDir(typ, archived)
			files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
			want := map[string]*Item{}
			var recs []keptRecord
			for _, f := range files {
				it, err := ReadItem(f)
				if err != nil {
					continue
				}
				raw, err := encodeKept(it)
				if err != nil {
					t.Fatal(err)
				}
				unknown, err := encodeUnknown(it.Unknown)
				if err != nil {
					t.Fatal(err)
				}
				name := filepath.Base(f)
				want[name] = it
				recs = append(recs, keptRecord{Name: name, Size: 1, ModTime: 2, Item: raw, Unknown: unknown})
			}
			path := filepath.Join(t.TempDir(), "items.jsonl")
			if err := writeKeptItems(path, dir, recs); err != nil {
				t.Fatal(err)
			}
			kept := readKeptItems(path, dir)
			if len(kept) != len(want) {
				t.Fatalf("%s: the kept file holds %d entries, want %d", dir, len(kept), len(want))
			}
			for name, w := range want {
				got, err := restoreItem(filepath.Join(dir, name), kept[name])
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}
				if !reflect.DeepEqual(got, w) {
					t.Errorf("%s restores as\n%#v\nnot\n%#v", name, got, w)
				}
				n++
			}
		}
	}
	if n == 0 {
		t.Skip("no items found")
	}
}

// Repo.List lists its folders at once (S-0676): it still gives every item of
// the current and archived folders once, marked archived by its folder and
// sorted by ID, and the error of a folder it cannot read.
func TestListGivesEveryFolderAtOnceInIDOrder(t *testing.T) {
	r := newProject(t)
	epic := mustCreate(t, r, Epic, "Epic", "")
	var stories []*Item
	for i := range 6 {
		stories = append(stories, mustCreate(t, r, Story, fmt.Sprintf("Story %d", i), epic.ID))
	}
	mustCreate(t, r, Task, "Task", stories[0].ID)
	for _, it := range stories[:3] {
		to := filepath.Join(r.ItemDir(Story, true), filepath.Base(it.Path))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(it.Path, to); err != nil {
			t.Fatal(err)
		}
	}
	items, err := r.List(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 8 {
		t.Fatalf("listed %d items, want 8", len(items))
	}
	for i, it := range items {
		if i > 0 && !lessID(items[i-1].ID, it.ID) {
			t.Errorf("%s listed after %s", it.ID, items[i-1].ID)
		}
		if want := strings.Contains(it.Path, string(filepath.Separator)+"archive"+string(filepath.Separator)); it.Archived != want {
			t.Errorf("%s: archived %t, from %s", it.ID, it.Archived, it.Path)
		}
	}
	tasks := r.ItemDir(Task, true)
	if err := os.MkdirAll(filepath.Dir(tasks), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, tasks, "a file where a folder should be")
	if _, err := r.List(true); err == nil {
		t.Error("a folder that cannot be read listed without an error")
	}
	if items, err := r.List(false); err != nil || len(items) != 5 {
		t.Errorf("without the archive: %d items, %v", len(items), err)
	}
}

func mustWrite(t *testing.T, path, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRemove(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}
