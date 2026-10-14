package storygit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scriptedGit answers each git run with the output scripted for its
// arguments, and fails a run it has no output for.
type scriptedGit map[string]string

func (g scriptedGit) Run(_, name string, args ...string) (string, error) {
	out, ok := g[strings.Join(args, " ")]
	if !ok {
		return "", errors.New(name + " " + strings.Join(args, " ") + ": not scripted")
	}
	return out, nil
}

func (g scriptedGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (scriptedGit) LookPath(name string) (string, error) { return name, nil }

// script sets the output of git run with args.
func (g scriptedGit) script(out string, args ...string) { g[strings.Join(args, " ")] = out }

func TestUndoneLinesCountsMainsAddedLinesTheBranchRemoves(t *testing.T) {
	for _, c := range []struct {
		name           string
		added, removed []string
		want           []string
	}{
		{"removes main's line", []string{"two"}, []string{"two"}, []string{"two"}},
		{"removes an unrelated line", []string{"two"}, []string{"one"}, []string{}},
		{"keeps main's change", []string{"two", "three"}, nil, []string{}},
		{"counted as a multiset", []string{"x", "y", "x", "x"}, []string{"x", "x", "z"}, []string{"x", "x"}},
	} {
		if got := undoneLines(c.added, c.removed); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

func TestDiffLinesReadsHunksAndSkipsHeaders(t *testing.T) {
	diff := "diff --git a/a.md b/a.md\nindex 1..2 100644\n--- a/a.md\n+++ b/a.md\n" +
		"@@ -1 +1,2 @@\n-old  \n+--- new\n+++ also\n\\ No newline at end of file"
	added, removed := diffLines(diff)
	if !reflect.DeepEqual(added, []string{"--- new", "++ also"}) || !reflect.DeepEqual(removed, []string{"old"}) {
		t.Errorf("added %q, removed %q", added, removed)
	}
}

// undoneGit scripts the check of the worktree wt, rebased onto main from
// old, where main changed mainPaths and the branch changes branchPaths
// against it.
func undoneGit(head, old string, mainPaths, branchPaths []string) scriptedGit {
	g := scriptedGit{}
	g.script(head, "rev-parse", "main")
	g.script(strings.Join(mainPaths, "\n"), nameOnlyArgs(old, "main")...)
	g.script(strings.Join(branchPaths, "\n"), nameOnlyArgs("main", "HEAD")...)
	return g
}

// ADR-0155: a path where the branch removes lines main added is reported
// with them; one where it removes only its own lines, or keeps main's, is
// not; a path the branch deletes removes every line main added.
func TestFindUndone(t *testing.T) {
	g := undoneGit("h1", "b0", []string{"gone.md", "keeps.md", "own.md", "undone.md", "mainonly.md"}, []string{"gone.md", "keeps.md", "own.md", "undone.md"})
	g.script("@@ -1 +1,2 @@\n line\n+two", diffArgs("b0", "main", "undone.md")...)
	g.script("@@ -2 +1,0 @@\n-two", diffArgs("main", "HEAD", "undone.md")...)
	g.script("@@ -0,0 +1 @@\n+main's", diffArgs("b0", "main", "own.md")...)
	g.script("@@ -3 +3 @@\n-branch's old\n+branch's new", diffArgs("main", "HEAD", "own.md")...)
	g.script("@@ -0,0 +1 @@\n+kept", diffArgs("b0", "main", "keeps.md")...)
	g.script("@@ -2,0 +2 @@\n+more", diffArgs("main", "HEAD", "keeps.md")...)
	g.script("@@ -1 +1,2 @@\n+a\n+b", diffArgs("b0", "main", "gone.md")...)
	g.script("diff --git a/gone.md b/gone.md\ndeleted file mode 100644\n--- a/gone.md\n+++ /dev/null\n@@ -1,3 +0,0 @@\n-a\n-b\n-c", diffArgs("main", "HEAD", "gone.md")...)

	got, err := FindUndone(g, "/wt", "main", "b0")
	if err != nil {
		t.Fatal(err)
	}
	want := []Undone{{Path: "gone.md", Lines: []string{"a", "b"}, OldBase: "b0"}, {Path: "undone.md", Lines: []string{"two"}, OldBase: "b0"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("found %+v, want %+v", got, want)
	}
}

// A rebase whose old base is main's head has nothing to check.
func TestFindUndoneChecksNothingWhenMainHasNotMoved(t *testing.T) {
	got, err := FindUndone(undoneGit("b0", "b0", nil, nil), "/wt", "main", "b0")
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("found %+v, %v", got, err)
	}
}

// undoneRepo is a project at a temporary folder, its main branch main.
func undoneRepo(t *testing.T) (*workitem.Repo, scriptedGit) {
	t.Helper()
	dir := t.TempDir()
	g := scriptedGit{}
	g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
	return &workitem.Repo{Root: dir, MainRoot: dir}, g
}

// readRecord reads story's record in repo, failing when it has none.
func readRecord(t *testing.T, repo *workitem.Repo, story string) UndoneRecord {
	t.Helper()
	data, err := os.ReadFile(UndoneRecordPath(repo.CacheDir(), story))
	if err != nil {
		t.Fatal(err)
	}
	var rec UndoneRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		t.Fatal(err)
	}
	return rec
}

// ADR-0155: the record holds each path found, replacing an earlier finding
// for the same path, with when it was found, in UTC.
func TestRecordUndoneReplacesAPathsEarlierFinding(t *testing.T) {
	repo, _ := undoneRepo(t)
	at := time.Date(2026, 10, 14, 20, 0, 0, 500, time.FixedZone("CEST", 2*3600))
	if err := recordUndone(repo.CacheDir(), "S-0001", []Undone{{Path: "b.md", Lines: []string{"x"}, OldBase: "b0"}, {Path: "a.md", Lines: []string{"y"}, OldBase: "b0"}}, at); err != nil {
		t.Fatal(err)
	}
	if err := recordUndone(repo.CacheDir(), "S-0001", []Undone{{Path: "b.md", Lines: []string{"z"}, OldBase: "b1"}}, at.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	utc := time.Date(2026, 10, 14, 18, 0, 0, 0, time.UTC)
	want := UndoneRecord{Story: "S-0001", Paths: []UndoneEntry{
		{Undone: Undone{Path: "a.md", Lines: []string{"y"}, OldBase: "b0"}, FoundAt: utc},
		{Undone: Undone{Path: "b.md", Lines: []string{"z"}, OldBase: "b1"}, FoundAt: utc.Add(time.Hour)},
	}}
	if got := readRecord(t, repo, "S-0001"); !reflect.DeepEqual(got, want) {
		t.Errorf("record %+v, want %+v", got, want)
	}
	data, _ := os.ReadFile(UndoneRecordPath(repo.CacheDir(), "S-0001"))
	for _, key := range []string{`"story": "S-0001"`, `"path": "a.md"`, `"lines": [`, `"old_base": "b0"`, `"found_at": "2026-10-14T18:00:00Z"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("record lacks %s:\n%s", key, data)
		}
	}
}

// ADR-0155: the reader gives the recorded paths the branch still undoes,
// checked against main now; a path whose branch no longer removes its lines
// drops out of the record, and the last one to go takes the record with it.
func TestRecordedUndoneDropsAPathTheBranchNoLongerRemoves(t *testing.T) {
	repo, g := undoneRepo(t)
	if got, err := RecordedUndone(g, repo, "S-0001"); err != nil || got == nil || len(got) != 0 {
		t.Fatalf("no record: %+v %v", got, err)
	}
	found := []Undone{{Path: "a.md", Lines: []string{"x", "y"}, OldBase: "b0"}, {Path: "b.md", Lines: []string{"z"}, OldBase: "b0"}}
	if err := recordUndone(repo.CacheDir(), "S-0001", found, time.Date(2026, 10, 14, 18, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	g.script("@@ -1 +0,0 @@\n-y", diffArgs("main", "story/S-0001", "a.md")...)
	g.script("", diffArgs("main", "story/S-0001", "b.md")...)

	got, err := RecordedUndone(g, repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if want := []Undone{{Path: "a.md", Lines: []string{"y"}, OldBase: "b0"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("still undone %+v, want %+v", got, want)
	}
	if rec := readRecord(t, repo, "S-0001"); len(rec.Paths) != 1 || rec.Paths[0].Path != "a.md" {
		t.Errorf("record not rewritten without b.md: %+v", rec)
	}

	g.script("", diffArgs("main", "story/S-0001", "a.md")...)
	if got, err := RecordedUndone(g, repo, "S-0001"); err != nil || len(got) != 0 {
		t.Fatalf("restored: %+v %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(repo.CacheDir(), "undone", "S-0001.json")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the empty record is left: %v", err)
	}
}
