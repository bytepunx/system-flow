package conflictmark

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"
)

// countReads replaces a Kept's read with one that records the files read,
// until the test ends.
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

// writeOld writes body to path dated age ago.
func writeOld(t *testing.T, path, body string, age time.Duration) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-age)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
}

// keptTree writes a root of three markdown files dated an hour ago, one with
// a conflict, and a cache folder beside them, and returns the root, the
// cache folder, and the files' paths.
func keptTree(t *testing.T) (root, cache string, paths []string) {
	t.Helper()
	root = t.TempDir()
	cache = filepath.Join(root, ".flai-cache")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	files := []struct{ rel, body string }{
		{"README.md", "# Project\n"},
		{"design/system/a.md", "# A\n\n" + ours + " HEAD\nmine\n" + theirs + " topic\n"},
		{"wip/agents/S-0001.md", "# S-0001\n"},
	}
	for _, f := range files {
		p := filepath.Join(root, filepath.FromSlash(f.rel))
		writeOld(t, p, f.body, time.Hour)
		paths = append(paths, p)
	}
	return root, cache, paths
}

// scanAll scans paths with a Kept for root and cache, saves it, and returns
// each path's lines.
func scanAll(t *testing.T, root, cache string, paths []string) map[string][]int {
	t.Helper()
	k := OpenKept(cache, root)
	out := map[string][]int{}
	for _, p := range paths {
		lines, err := k.Lines(p)
		if err != nil {
			t.Fatalf("lines of %s: %v", p, err)
		}
		out[p] = lines
	}
	k.Save()
	return out
}

func TestASecondScanReadsNoUnchangedFile(t *testing.T) {
	root, cache, paths := keptTree(t)
	reads := countReads(t)
	cold := scanAll(t, root, cache, paths)
	if got := reads(); !reflect.DeepEqual(got, slices.Sorted(slices.Values(paths))) {
		t.Fatalf("a cold scan reads every file: %v", got)
	}
	if want := []int{3, 5}; !reflect.DeepEqual(cold[paths[1]], want) {
		t.Fatalf("cold lines of %s: %v, want %v", paths[1], cold[paths[1]], want)
	}
	if _, err := os.Stat(KeptPath(cache, root)); err != nil {
		t.Fatalf("the scan keeps its lines: %v", err)
	}
	warm := scanAll(t, root, cache, paths)
	if got := reads(); len(got) != 0 {
		t.Errorf("a warm scan reads no unchanged file, read %v", got)
	}
	if !reflect.DeepEqual(warm, cold) {
		t.Errorf("warm lines %v, want the cold scan's %v", warm, cold)
	}
}

func TestAnEditedFileIsReadAgainAlone(t *testing.T) {
	root, cache, paths := keptTree(t)
	reads := countReads(t)
	scanAll(t, root, cache, paths)
	reads()
	// The README gains a conflict, dated otherwise but still long ago.
	writeOld(t, paths[0], "# Project\n\n"+theirs+" main\n", 30*time.Minute)
	got := scanAll(t, root, cache, paths)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[0]}) {
		t.Errorf("only the edited file is read again, read %v", r)
	}
	if want := []int{3}; !reflect.DeepEqual(got[paths[0]], want) {
		t.Errorf("the marker added to a kept file is found: %v, want %v", got[paths[0]], want)
	}
	// The edit is kept in turn.
	scanAll(t, root, cache, paths)
	if r := reads(); len(r) != 0 {
		t.Errorf("a third scan reads nothing, read %v", r)
	}
}

func TestAFileNotScannedDropsOut(t *testing.T) {
	root, cache, paths := keptTree(t)
	scanAll(t, root, cache, paths)
	scanAll(t, root, cache, paths[:2])
	kept := readKept(KeptPath(cache, root), filepath.Clean(root))
	if _, ok := kept["wip/agents/S-0001.md"]; ok || len(kept) != 2 {
		t.Errorf("a file not scanned drops out of the kept lines: %v", kept)
	}
}

func TestARecentFileAndNoCacheFolderKeepNothing(t *testing.T) {
	root, cache, paths := keptTree(t)
	writeOld(t, paths[2], "# S-0001\n\nnew\n", 0)
	reads := countReads(t)
	scanAll(t, root, cache, paths)
	reads()
	scanAll(t, root, cache, paths)
	if r := reads(); !reflect.DeepEqual(r, []string{paths[2]}) {
		t.Errorf("a file changed within the racy window is read again, read %v", r)
	}

	none := filepath.Join(root, "no-cache")
	scanAll(t, root, none, paths)
	if _, err := os.Stat(none); !os.IsNotExist(err) {
		t.Errorf("a missing cache folder is not made: %v", err)
	}
	scanAll(t, root, "", paths)
	if r := reads(); len(r) != 2*len(paths) {
		t.Errorf("without a cache folder every file is read each scan, read %v", r)
	}
}

func TestAKeptFileThatDoesNotFitIsAMiss(t *testing.T) {
	root, cache, paths := keptTree(t)
	reads := countReads(t)
	for _, body := range []string{
		"not json",
		`{"version": 99, "root": "` + filepath.ToSlash(filepath.Clean(root)) + `", "files": {}}`,
		`{"version": 1, "root": "/elsewhere", "files": {"README.md": {"size": 10, "mtime": 0}}}`,
	} {
		if err := os.WriteFile(KeptPath(cache, root), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		got := scanAll(t, root, cache, paths)
		if r := reads(); len(r) != len(paths) {
			t.Errorf("kept file %q is a miss: read %v", body, r)
		}
		if want := []int{3, 5}; !reflect.DeepEqual(got[paths[1]], want) {
			t.Errorf("kept file %q: lines %v, want %v", body, got[paths[1]], want)
		}
	}
}

func TestAnUnreadableFileIsTheReadsError(t *testing.T) {
	root, cache, _ := keptTree(t)
	k := OpenKept(cache, root)
	if _, err := k.Lines(filepath.Join(root, "gone.md")); err == nil {
		t.Error("a missing file is an error")
	}
}
