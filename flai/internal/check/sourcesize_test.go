package check

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// maxSourceLines is the most lines a non-test Go file under flai may hold
// (I-0298): a file past it holds more than one responsibility, which
// code-quality.md's SRP rule asks to split.
const maxSourceLines = 2000

// TestNoSourceFileIsOver2000Lines walks the flai module and fails, naming
// each file and its count, on any non-test Go file over maxSourceLines. It
// reads files only, so it runs under -short too.
func TestNoSourceFileIsOver2000Lines(t *testing.T) {
	over, err := sourceFilesOver("../..", maxSourceLines)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range over {
		t.Errorf("%s is %d lines, over %d: split it by responsibility within its package (I-0298)", f.path, f.lines, maxSourceLines)
	}
}

// TestSourceFilesOverFindsALongFile checks that the walk counts lines and
// leaves out test files and testdata folders.
func TestSourceFilesOverFindsALongFile(t *testing.T) {
	root := t.TempDir()
	write := func(rel string, lines int) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, bytes.Repeat([]byte("//\n"), lines), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a/long.go", 11)
	write("a/short.go", 10)
	write("a/long_test.go", 50)
	write("a/testdata/long.go", 50)
	write("b/notes.txt", 50)
	over, err := sourceFilesOver(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(over) != 1 || over[0].path != filepath.Join("a", "long.go") || over[0].lines != 11 {
		t.Fatalf("over = %+v, want only a/long.go at 11 lines", over)
	}
}

type sourceFile struct {
	path  string
	lines int
}

// sourceFilesOver lists the non-test Go files under root, outside testdata,
// with more than limit lines, by path relative to root.
func sourceFilesOver(root string, limit int) ([]sourceFile, error) {
	var over []sourceFile
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		n := bytes.Count(data, []byte("\n"))
		if len(data) > 0 && data[len(data)-1] != '\n' {
			n++
		}
		if n > limit {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			over = append(over, sourceFile{path: rel, lines: n})
		}
		return nil
	})
	sort.Slice(over, func(i, j int) bool { return over[i].path < over[j].path })
	return over, err
}
