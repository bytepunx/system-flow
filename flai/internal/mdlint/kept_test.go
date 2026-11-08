package mdlint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"
)

// countReads replaces the folder lint's read with one that records the files
// read, until the test ends.
func countReads(t *testing.T) func() []string {
	t.Helper()
	var mu sync.Mutex
	var read []string
	readFile = func(name string) ([]byte, error) {
		mu.Lock()
		read = append(read, name)
		mu.Unlock()
		return os.ReadFile(name)
	}
	t.Cleanup(func() { readFile = os.ReadFile })
	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := slices.Clone(read)
		slices.Sort(out)
		read = nil
		return out
	}
}

// keptFolder writes a folder of markdown files, some with findings, dated an
// hour ago so that their findings may be kept, and a cache folder beside it.
func keptFolder(t *testing.T) (dir, cache string, paths []string) {
	t.Helper()
	root := t.TempDir()
	dir = filepath.Join(root, "wip")
	cache = filepath.Join(root, ".flai-cache")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	for i, name := range []string{"a.md", "b/c.md", "b/d.md", "e.md"} {
		p := filepath.Join(dir, name)
		body := "# Title\n\nText.\n"
		if i%2 == 0 {
			body = "# Title\n\nText  \n\nEnd"
		}
		writeOld(t, p, body)
		paths = append(paths, p)
	}
	return dir, cache, paths
}

// writeOld writes body to path dated an hour ago.
func writeOld(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
}

func mustParse(t *testing.T, config string) *Config {
	t.Helper()
	c, err := Parse([]byte(config), false, true)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestASecondFolderLintReadsNoUnchangedFile(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	c := mustParse(t, `{"default": true}`)
	reads := countReads(t)
	first := c.LintFolderKept(dir, cache)
	if got := reads(); len(got) != len(paths) {
		t.Fatalf("first lint read %v, want all %d files", got, len(paths))
	}
	if want := lintSerially(t, c, paths); !reflect.DeepEqual(first, want) {
		t.Fatalf("first lint = %v, want %v", first, want)
	}
	if _, err := os.Stat(KeptPath(cache, dir)); err != nil {
		t.Fatalf("no kept findings after the first lint: %v", err)
	}
	second := c.LintFolderKept(dir, cache)
	if got := reads(); len(got) != 0 {
		t.Fatalf("second lint read %v, want none", got)
	}
	if !reflect.DeepEqual(second, first) {
		t.Fatalf("kept findings = %v, want %v", second, first)
	}
	// A new process parses the configuration again and gets the same hash.
	third := mustParse(t, `{"default": true}`).LintFolderKept(dir, cache)
	if got := reads(); len(got) != 0 {
		t.Fatalf("a lint with the configuration parsed again read %v, want none", got)
	}
	if !reflect.DeepEqual(third, first) {
		t.Fatalf("kept findings = %v, want %v", third, first)
	}
}

func TestAnEditedFileIsLintedAgainAlone(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	c := mustParse(t, `{"default": true}`)
	reads := countReads(t)
	c.LintFolderKept(dir, cache)
	reads()
	writeOld(t, paths[1], "# Title\n\nNow it has trailing spaces  \n\nEnd")
	got := c.LintFolderKept(dir, cache)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[1]}) {
		t.Fatalf("lint after one edit read %v, want only %s", r, paths[1])
	}
	if want := lintSerially(t, c, paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	c.LintFolderKept(dir, cache)
	if r := reads(); len(r) != 0 {
		t.Fatalf("lint after the edit was kept read %v, want none", r)
	}
}

func TestAFileJustWrittenIsLintedAgain(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	c := mustParse(t, `{"default": true}`)
	// Written now, within the racy window: its findings are not kept.
	if err := os.WriteFile(paths[0], []byte("# Fresh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reads := countReads(t)
	c.LintFolderKept(dir, cache)
	reads()
	c.LintFolderKept(dir, cache)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[0]}) {
		t.Fatalf("second lint read %v, want only the file written now", r)
	}
}

func TestAChangedConfigurationLintsEveryFileAgain(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	reads := countReads(t)
	mustParse(t, `{"default": true}`).LintFolderKept(dir, cache)
	reads()
	c := mustParse(t, `{"default": true, "MD009": false}`)
	got := c.LintFolderKept(dir, cache)
	if r := reads(); len(r) != len(paths) {
		t.Fatalf("lint under a changed configuration read %v, want all %d files", r, len(paths))
	}
	if want := lintSerially(t, c, paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
}

func TestAConfigurationOfAnotherFormHasAnotherHash(t *testing.T) {
	data := []byte(`{"default": true}`)
	json, _ := Parse(data, false, true)
	yamlc, _ := Parse(data, false, false)
	cli2, _ := Parse([]byte(`{"config": {"default": true}}`), true, true)
	if json.Hash() == "" || json.Hash() == yamlc.Hash() || json.Hash() == cli2.Hash() {
		t.Fatalf("hashes %q %q %q: want three different, none empty", json.Hash(), yamlc.Hash(), cli2.Hash())
	}
	var none *Config
	if none.Hash() != "" {
		t.Fatal("a nil configuration has a hash")
	}
}

func TestADamagedOrForeignKeptFileIsAMiss(t *testing.T) {
	// lie keeps every file clean at its stat as it is now, as the kept file
	// of version, root, and hash.
	lie := func(t *testing.T, dir string, paths []string, version int, root, hash string) string {
		t.Helper()
		k := keptFindings{Version: version, Root: root, Hash: hash, Files: map[string]keptEntry{}}
		for _, p := range paths {
			info, err := os.Stat(p)
			if err != nil {
				t.Fatal(err)
			}
			k.Files[rel(dir, p)] = keptEntry{Size: info.Size(), ModTime: info.ModTime().UnixNano()}
		}
		data, err := json.Marshal(k)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	cases := map[string]func(t *testing.T, dir string, paths []string, hash string) string{
		"truncated": func(t *testing.T, dir string, paths []string, hash string) string {
			whole := lie(t, dir, paths, keptVersion, filepath.Clean(dir), hash)
			return whole[:len(whole)/2]
		},
		"not json": func(*testing.T, string, []string, string) string { return "not json" },
		"another root": func(t *testing.T, dir string, paths []string, hash string) string {
			return lie(t, dir, paths, keptVersion, "/elsewhere", hash)
		},
		"another hash": func(t *testing.T, dir string, paths []string, _ string) string {
			return lie(t, dir, paths, keptVersion, filepath.Clean(dir), "other")
		},
		"another format": func(t *testing.T, dir string, paths []string, hash string) string {
			return lie(t, dir, paths, keptVersion+1, filepath.Clean(dir), hash)
		},
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			dir, cache, paths := keptFolder(t)
			c := mustParse(t, `{"default": true}`)
			if err := os.WriteFile(KeptPath(cache, dir), []byte(content(t, dir, paths, c.Hash())), 0o644); err != nil {
				t.Fatal(err)
			}
			reads := countReads(t)
			got := c.LintFolderKept(dir, cache)
			if r := reads(); len(r) != len(paths) {
				t.Fatalf("read %v, want all %d files", r, len(paths))
			}
			if want := lintSerially(t, c, paths); !reflect.DeepEqual(got, want) {
				t.Fatalf("findings = %v, want %v", got, want)
			}
			c.LintFolderKept(dir, cache)
			if r := reads(); len(r) != 0 {
				t.Fatalf("the kept file was not written again: second lint read %v", r)
			}
		})
	}
}

func TestAMissingKeptFileIsAMiss(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	c := mustParse(t, `{"default": true}`)
	if got, want := c.LintFolderKept(dir, cache), lintSerially(t, c, paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
}

func TestNoCacheFolderKeepsNothing(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	if err := os.Remove(cache); err != nil {
		t.Fatal(err)
	}
	c := mustParse(t, `{"default": true}`)
	reads := countReads(t)
	for range 2 {
		if got, want := c.LintFolderKept(dir, cache), lintSerially(t, c, paths); !reflect.DeepEqual(got, want) {
			t.Fatalf("findings = %v, want %v", got, want)
		}
		if r := reads(); len(r) != len(paths) {
			t.Fatalf("read %v, want all %d files each time", r, len(paths))
		}
	}
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		t.Fatalf("the cache folder was made: %v", err)
	}
}

func TestARemovedFileIsDroppedFromTheKeptFindings(t *testing.T) {
	dir, cache, paths := keptFolder(t)
	c := mustParse(t, `{"default": true}`)
	c.LintFolderKept(dir, cache)
	if err := os.Remove(paths[0]); err != nil {
		t.Fatal(err)
	}
	got := c.LintFolderKept(dir, cache)
	if want := lintSerially(t, c, paths[1:]); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	if files := readKept(KeptPath(cache, dir), filepath.Clean(dir), c.Hash()); len(files) != len(paths)-1 {
		t.Fatalf("kept %d files, want %d", len(files), len(paths)-1)
	}
}

func TestKeptPathDiffersByFolder(t *testing.T) {
	if KeptPath("/c", "/a/wip") == KeptPath("/c", "/b/wip") {
		t.Fatal("two folders share one kept file")
	}
}
