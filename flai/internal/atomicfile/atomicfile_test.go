package atomicfile

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// raceReads saves path over and over with write while another goroutine reads
// it as fast as it can, and counts the reads that saw neither version whole.
func raceReads(t *testing.T, write func(path string, data []byte) error) int64 {
	t.Helper()
	path := filepath.Join(t.TempDir(), "S-0001-story.md")
	a := []byte("---\nid: S-0001\ntitle: A\n---\n" + strings.Repeat("a", 64<<10))
	b := []byte("---\nid: S-0001\ntitle: B\n---\n" + strings.Repeat("b", 64<<10))
	if err := os.WriteFile(path, a, 0o644); err != nil {
		t.Fatal(err)
	}
	var torn atomic.Int64
	var stop atomic.Bool
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stop.Load() {
			data, err := os.ReadFile(path)
			whole := bytes.Equal(data, a) || bytes.Equal(data, b)
			if err != nil || !whole {
				torn.Add(1)
			}
		}
	}()
	for i := range 2000 {
		data := a
		if i%2 == 1 {
			data = b
		}
		if err := write(path, data); err != nil {
			t.Fatal(err)
		}
	}
	stop.Store(true)
	wg.Wait()
	return torn.Load()
}

// S-0100: a reader never sees a file WriteFile is replacing half-written; with
// os.WriteFile it does, which is how an agent's wait_for_events failed in CI
// with "no front matter".
func TestAReaderNeverSeesAHalfWrittenFile(t *testing.T) {
	if n := raceReads(t, func(p string, d []byte) error { return WriteFile(p, d, 0o644) }); n != 0 {
		t.Errorf("%d reads saw a file half-written", n)
	}
	if testing.Short() {
		return
	}
	// the old way, to show the test can tell
	if n := raceReads(t, func(p string, d []byte) error { return os.WriteFile(p, d, 0o644) }); n == 0 {
		t.Log("os.WriteFile happened not to be caught half-written this time")
	}
}

func TestWriteFileKeepsTheModeAndLeavesNothingBehind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "board.md")
	if err := WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("mode: %v %v", info, err)
	}
	if err := WriteFile(filepath.Join(dir, "missing", "x.md"), []byte("x"), 0o644); err == nil {
		t.Error("a write into a folder that does not exist must fail")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("left behind: %v", entries)
	}
}

// I-0137: WriteFile writes the content and mode given, replaces what was there,
// and leaves no temporary file; the syncs it adds for a crash cannot be
// simulated here, only that they do not fail on a real folder.
func TestWriteFileReplacesTheContentAndLeavesNoTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "S-0337.md")
	if err := WriteFile(path, []byte("first\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(path, []byte("second, longer\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "second, longer\n" {
		t.Errorf("content: got %q, want %q", data, "second, longer\n")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode: got %v, want %v", info.Mode().Perm(), os.FileMode(0o600))
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "S-0337.md" {
		t.Errorf("folder holds %v, want only S-0337.md", entries)
	}
}

func TestSyncDirSyncsAFolderAndReportsOneThatIsMissing(t *testing.T) {
	dir := t.TempDir()
	if err := syncDir(dir); err != nil {
		t.Errorf("syncing %s: %v", dir, err)
	}
	if runtime.GOOS == "windows" {
		return
	}
	if err := syncDir(filepath.Join(dir, "missing")); err == nil {
		t.Error("syncing a folder that does not exist must fail")
	}
}
