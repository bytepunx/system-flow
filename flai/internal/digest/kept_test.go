package digest

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// countEntryReads replaces CheckFolder's read with one that records the
// files read, until the test ends.
func countEntryReads(t *testing.T) func() []string {
	t.Helper()
	var mu sync.Mutex
	var read []string
	readEntryFile = func(name string) ([]byte, error) {
		mu.Lock()
		read = append(read, name)
		mu.Unlock()
		return os.ReadFile(name)
	}
	t.Cleanup(func() { readEntryFile = os.ReadFile })
	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := slices.Clone(read)
		slices.Sort(out)
		read = nil
		return out
	}
}

// keptDigest writes a digest folder of entries dated an hour ago, so that
// their problems may be kept: a clean one, one with a bad kind, one whose
// source S-0999 names nothing, one absorbed by a file that is not there, and
// the README, with a cache folder beside the folder. It returns the folder,
// the cache folder, and the entries' paths.
func keptDigest(t *testing.T) (dir, cache string, paths []string) {
	t.Helper()
	root := t.TempDir()
	dir = filepath.Join(root, "design", "digest")
	cache = filepath.Join(root, ".flai-cache")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	with := func(pairs ...string) string { return strings.NewReplacer(pairs...).Replace(good) }
	entries := []struct{ name, content string }{
		{"D-0001-clean.md", with("D-0003", "D-0001")},
		{"D-0002-kind.md", with("D-0003", "D-0002", "kind: decision", "kind: rumour")},
		{"D-0003-the-digest-is-one-file-per-entry.md", with("S-0428,", "S-0999,")},
		{"D-0004-gone.md", with("D-0003", "D-0004", "status: active", "status: absorbed", `absorbed_by: ""`, "absorbed_by: design/gone.md#rules")},
	}
	writeOldEntry(t, filepath.Join(dir, Index), "# The digest\n")
	for _, e := range entries {
		p := filepath.Join(dir, e.name)
		writeOldEntry(t, p, e.content)
		paths = append(paths, p)
	}
	return dir, cache, paths
}

// writeOldEntry writes body to path dated an hour ago.
func writeOldEntry(t *testing.T, path, body string) {
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

// validateSerially is what CheckFolder must answer: Validate on each entry.
func validateSerially(t *testing.T, paths []string, env Env) []FileProblems {
	t.Helper()
	var out []FileProblems
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if _, ps := Validate(filepath.Base(p), string(data), env); len(ps) > 0 {
			out = append(out, FileProblems{Path: p, Problems: ps})
		}
	}
	return out
}

func mustCheckFolder(t *testing.T, dir, cache string, env Env) []FileProblems {
	t.Helper()
	got, err := CheckFolder(dir, cache, env)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// S-0677: a second check of an unchanged digest reads no entry, and answers
// what Validate answers on each.
func TestASecondCheckReadsNoUnchangedEntry(t *testing.T) {
	dir, cache, paths := keptDigest(t)
	reads := countEntryReads(t)
	first := mustCheckFolder(t, dir, cache, env)
	if got := reads(); !reflect.DeepEqual(got, paths) {
		t.Fatalf("first check read %v, want %v", got, paths)
	}
	want := validateSerially(t, paths, env)
	if len(want) != 3 {
		t.Fatalf("the fixture should have three entries with problems, got %+v", want)
	}
	if !reflect.DeepEqual(first, want) {
		t.Fatalf("first check = %+v, want %+v", first, want)
	}
	if _, err := os.Stat(KeptPath(cache, dir)); err != nil {
		t.Fatalf("no kept problems after the first check: %v", err)
	}
	second := mustCheckFolder(t, dir, cache, env)
	if got := reads(); len(got) != 0 {
		t.Fatalf("second check read %v, want none", got)
	}
	if !reflect.DeepEqual(second, first) {
		t.Fatalf("kept problems = %+v, want %+v", second, first)
	}
}

// S-0677: an entry whose size or modification time changed is read again,
// alone, and its new problems are kept.
func TestAnEditedEntryIsCheckedAgainAlone(t *testing.T) {
	dir, cache, paths := keptDigest(t)
	reads := countEntryReads(t)
	mustCheckFolder(t, dir, cache, env)
	reads()
	data, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	writeOldEntry(t, paths[0], strings.Replace(string(data), "status: active", "status: stale", 1))
	got := mustCheckFolder(t, dir, cache, env)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[0]}) {
		t.Fatalf("check after one edit read %v, want only %s", r, paths[0])
	}
	if want := validateSerially(t, paths, env); !reflect.DeepEqual(got, want) {
		t.Fatalf("problems = %+v, want %+v", got, want)
	}
	mustCheckFolder(t, dir, cache, env)
	if r := reads(); len(r) != 0 {
		t.Fatalf("check after the edit was kept read %v, want none", r)
	}
}

// S-0677: what Env answers is never kept: a source archived or removed, or
// an absorbed_by file deleted, since the kept run is reported on the next
// check without the entry being read, and one that comes back is no longer.
func TestARemovedSourceIsReportedAfterTheKeptRun(t *testing.T) {
	dir, cache, paths := keptDigest(t)
	all := Env{Known: func(string) bool { return true }, Exists: func(string) bool { return true }}
	reads := countEntryReads(t)
	first := mustCheckFolder(t, dir, cache, all)
	reads()
	if len(first) != 1 || first[0].Path != paths[1] {
		t.Fatalf("with every source known, want only the bad kind, got %+v", first)
	}
	gone := Env{Known: func(id string) bool { return id != "TH-0476" }, Exists: func(string) bool { return false }}
	got := mustCheckFolder(t, dir, cache, gone)
	if r := reads(); len(r) != 0 {
		t.Fatalf("check read %v, want none", r)
	}
	if want := validateSerially(t, paths, gone); !reflect.DeepEqual(got, want) {
		t.Fatalf("problems = %+v, want %+v", got, want)
	}
	var says []string
	for _, f := range got {
		for _, p := range f.Problems {
			says = append(says, p.Message)
		}
	}
	for _, want := range []string{"source TH-0476 names nothing", "absorbed_by is design/gone.md#rules, which names no file"} {
		if !slices.ContainsFunc(says, func(s string) bool { return strings.Contains(s, want) }) {
			t.Errorf("want a problem saying %q, got %q", want, says)
		}
	}
	if again := mustCheckFolder(t, dir, cache, all); !reflect.DeepEqual(again, first) {
		t.Errorf("with the sources back, problems = %+v, want %+v", again, first)
	}
}

// S-0677: a kept file that is damaged, of another folder, or kept under
// other rules is a miss: every entry is read and the problems are Validate's.
func TestAKeptFileThatDoesNotFitIsAMiss(t *testing.T) {
	for _, c := range []struct {
		name string
		kept func(dir string) string
	}{
		{"damaged", func(string) string { return "{not json" }},
		{"another folder", func(string) string {
			return `{"version":1,"root":"/elsewhere","hash":"` + rulesHash() + `","files":{}}`
		}},
		{"other rules", func(dir string) string {
			return `{"version":1,"root":"` + filepath.ToSlash(filepath.Clean(dir)) + `","hash":"old","files":{}}`
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir, cache, paths := keptDigest(t)
			if err := os.WriteFile(KeptPath(cache, dir), []byte(c.kept(dir)), 0o644); err != nil {
				t.Fatal(err)
			}
			reads := countEntryReads(t)
			got := mustCheckFolder(t, dir, cache, env)
			if r := reads(); !reflect.DeepEqual(r, paths) {
				t.Fatalf("read %v, want every entry", r)
			}
			if want := validateSerially(t, paths, env); !reflect.DeepEqual(got, want) {
				t.Fatalf("problems = %+v, want %+v", got, want)
			}
		})
	}
}

// S-0677: an entry changed within the racy window is not kept, so the next
// check reads it again; a project with no cache folder keeps nothing and
// makes none.
func TestARecentEntryAndNoCacheFolderKeepNothing(t *testing.T) {
	dir, cache, paths := keptDigest(t)
	data, err := os.ReadFile(paths[2])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(paths[2], data, 0o644); err != nil {
		t.Fatal(err)
	}
	reads := countEntryReads(t)
	mustCheckFolder(t, dir, cache, env)
	reads()
	mustCheckFolder(t, dir, cache, env)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[2]}) {
		t.Fatalf("second check read %v, want only the entry written just now", r)
	}

	none := filepath.Join(t.TempDir(), ".flai-cache")
	got := mustCheckFolder(t, dir, none, env)
	if want := validateSerially(t, paths, env); !reflect.DeepEqual(got, want) {
		t.Fatalf("problems = %+v, want %+v", got, want)
	}
	if _, err := os.Stat(none); !os.IsNotExist(err) {
		t.Fatalf("a check made the cache folder %s: %v", none, err)
	}
}

// S-0677: a digest folder that is not there is the error.
func TestCheckFolderOfNoFolderIsAnError(t *testing.T) {
	if _, err := CheckFolder(filepath.Join(t.TempDir(), "digest"), "", env); err == nil {
		t.Fatal("want an error for a folder that is not there")
	}
}
