package mdlint

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// lintSerially is what flai check did before the pool: walk, read, and lint
// one file at a time.
func lintSerially(t *testing.T, c *Config, paths []string) []FileFindings {
	t.Helper()
	var out []FileFindings
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if f := c.Lint(string(data)); len(f) > 0 {
			out = append(out, FileFindings{Path: p, Findings: f})
		}
	}
	return out
}

func TestLintFolderMatchesSerialOnSeveralWorkers(t *testing.T) {
	if runtime.GOMAXPROCS(0) < 2 {
		defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	}
	c, err := Parse([]byte(`{"default": true}`), false, true)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	var paths []string
	for i, name := range []string{"a.md", "b/c.md", "b/d/e.md", "f.md", "g.md", "h/i.md", "j.md", "k.md"} {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		// Every other file breaks MD009 and MD047 on its own lines; the rest are clean.
		body := "# Title\n\nText.\n"
		if i%2 == 0 {
			body = "# Title\n\nText " + strings.Repeat("x", i) + "  \n\nEnd"
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}
	// Not markdown, and a link that dangles: both are left out.
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing.md"), filepath.Join(dir, "b", "dangling.md")); err != nil {
		t.Fatal(err)
	}
	// filepath.WalkDir's order: lexical within each folder.
	walked := []string{paths[0], paths[1], paths[2], paths[3], paths[4], paths[5], paths[6], paths[7]}
	want := lintSerially(t, c, walked)
	if len(want) != 4 {
		t.Fatalf("fixture has %d files with findings, want 4: %+v", len(want), want)
	}
	for range 20 {
		got := c.LintFolder(dir)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("LintFolder =\n%+v\nwant\n%+v", got, want)
		}
	}
}

func TestLintFolderSkipsUnreadableFiles(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads a file whatever its mode")
	}
	c, err := Parse([]byte(`{"default": true}`), false, true)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	locked := filepath.Join(dir, "locked.md")
	if err := os.WriteFile(locked, []byte("x"), 0o000); err != nil {
		t.Fatal(err)
	}
	open := filepath.Join(dir, "open.md")
	if err := os.WriteFile(open, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := c.LintFolder(dir)
	if len(got) != 1 || got[0].Path != open {
		t.Fatalf("LintFolder = %+v, want only %s", got, open)
	}
}

func TestLintFolderMissingOrNil(t *testing.T) {
	var none *Config
	if got := none.LintFolder(t.TempDir()); got != nil {
		t.Fatalf("nil configuration found %+v", got)
	}
	c, err := Parse([]byte(`{"default": true}`), false, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.LintFolder(filepath.Join(t.TempDir(), "absent")); got != nil {
		t.Fatalf("missing folder found %+v", got)
	}
}
