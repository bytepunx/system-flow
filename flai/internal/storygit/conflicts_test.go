package storygit

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// conflictLog reads the conflict log in dir, failing the test on an error.
func conflictLog(t *testing.T, dir string) []Conflict {
	t.Helper()
	got, err := ReadConflicts(dir)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// S-0337: each conflict is one line, its time in UTC to the second and its
// paths sorted, read back as it was appended, oldest first.
func TestAppendAndReadConflicts(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".flai-cache")
	at := time.Date(2026, 10, 12, 17, 30, 15, 500, time.FixedZone("CEST", 2*3600))
	if err := AppendConflict(dir, Conflict{At: at, Story: "S-0001", Other: "S-0002", Where: WhereTrial, Paths: []string{"b.md", "a.md"}}); err != nil {
		t.Fatal(err)
	}
	if err := AppendConflict(dir, Conflict{At: at, Story: "S-0001", Where: WhereSync}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(ConflictsPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"at":"2026-10-12T15:30:15Z","story":"S-0001","other":"S-0002","where":"trial","paths":["a.md","b.md"]}` + "\n" +
		`{"at":"2026-10-12T15:30:15Z","story":"S-0001","where":"sync","paths":[]}` + "\n"
	if string(data) != want {
		t.Errorf("log:\n%s\nwant:\n%s", data, want)
	}
	utc := time.Date(2026, 10, 12, 15, 30, 15, 0, time.UTC)
	wantRead := []Conflict{
		{At: utc, Story: "S-0001", Other: "S-0002", Where: WhereTrial, Paths: []string{"a.md", "b.md"}},
		{At: utc, Story: "S-0001", Where: WhereSync, Paths: []string{}},
	}
	if got := conflictLog(t, dir); !reflect.DeepEqual(got, wantRead) {
		t.Errorf("read: %+v, want %+v", got, wantRead)
	}
	if err := AppendConflict(dir, Conflict{At: at, Story: "S-0001", Where: "rebase"}); err == nil {
		t.Error("a conflict met at no known place was appended")
	}
}

// A missing log is no conflicts; a log with a line that is not JSON is an
// error naming the line, and blank lines are skipped.
func TestReadConflictsMissingAndBadLine(t *testing.T) {
	dir := t.TempDir()
	if got, err := ReadConflicts(dir); err != nil || got != nil {
		t.Fatalf("missing log: %+v %v", got, err)
	}
	good := `{"at":"2026-10-12T15:30:15Z","story":"S-0001","where":"sync","paths":["a.md"]}`
	if err := os.WriteFile(ConflictsPath(dir), []byte(good+"\n\n"+good+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := conflictLog(t, dir); len(got) != 2 {
		t.Errorf("blank line not skipped: %+v", got)
	}
	if err := os.WriteFile(ConflictsPath(dir), []byte(good+"\n\nnot json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadConflicts(dir)
	if err == nil || got != nil || !strings.Contains(err.Error(), "line 3 of the conflict log") {
		t.Errorf("bad line: %+v %v", got, err)
	}
}

// ADR-0098, ADR-0126: the files flai resolves itself are left out of a
// line's paths, and a conflict on them alone is no line; a line that cannot
// be written is logged at warn, and nothing fails.
func TestRecordConflictLeavesOutResolvedFilesAndWarnsOnAFailedWrite(t *testing.T) {
	dir := t.TempDir()
	files := Files{
		Generated: []GeneratedFile{{Path: "design/issues/summary.md"}},
		Merged:    MergedFiles{Match: func(p string) bool { return p == "design/issues/I-0001-x.md" }},
	}
	RecordConflict(dir, Conflict{Story: "S-0001", Where: WhereSync, Paths: []string{"design/issues/summary.md", "design/issues/I-0001-x.md"}}, files, nil)
	if got := conflictLog(t, dir); len(got) != 0 {
		t.Fatalf("a conflict on resolved files alone was logged: %+v", got)
	}
	RecordConflict(dir, Conflict{Story: "S-0001", Where: WhereAccept, Paths: []string{"docs/guide.md", "design/issues/summary.md"}}, files, nil)
	if got := conflictLog(t, dir); len(got) != 1 || got[0].Where != WhereAccept || strings.Join(got[0].Paths, ",") != "docs/guide.md" {
		t.Fatalf("log: %+v", got)
	}

	var logged bytes.Buffer
	notADir := filepath.Join(dir, "file")
	if err := os.WriteFile(notADir, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	RecordConflict(notADir, Conflict{Story: "S-0001", Where: WhereSync, Paths: []string{"a.md"}}, Files{}, slog.New(slog.NewTextHandler(&logged, nil)))
	for _, want := range []string{"level=WARN", "component=storygit", "story=S-0001", "where=sync", "err="} {
		if !strings.Contains(logged.String(), want) {
			t.Errorf("the failed write's warning lacks %s: %s", want, logged.String())
		}
	}
}
