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

// I-0360: in a markdown file, main's line is kept when a line the branch adds
// holds every word of it in order, each added line keeping one of main's at
// most; a dropped, changed, or reordered word, a line with no words, and a
// file that is not markdown keep the multiset of whole lines.
func TestNotKeptLeavesMainsLinesABranchLineExtends(t *testing.T) {
	row := "| `flai check` | Validates the repo. |"
	cases := []struct {
		name, path          string
		undone, branchAdded []string
		want                []string
	}{
		{"extended row", "a.md", []string{row}, []string{"| `flai check` | Validates the repo. It also checks images. |"}, []string{}},
		{"word inserted inside", "a.md", []string{"one two three"}, []string{"one and two three"}, []string{}},
		{"upper-case extension", "A.MD", []string{"one two"}, []string{"one two three"}, []string{}},
		{"word dropped", "a.md", []string{"one two three"}, []string{"one three four"}, []string{"one two three"}},
		{"word changed", "a.md", []string{"one two"}, []string{"one too"}, []string{"one two"}},
		{"words reordered", "a.md", []string{"one two"}, []string{"two one"}, []string{"one two"}},
		{"revert adds nothing", "a.md", []string{"one two"}, nil, []string{"one two"}},
		{"one added line keeps one", "a.md", []string{"one two", "one two"}, []string{"one two three"}, []string{"one two"}},
		{"each added line keeps one", "a.md", []string{"a b", "c d"}, []string{"c d e", "a b f"}, []string{}},
		{"blank line", "a.md", []string{""}, []string{"anything"}, []string{""}},
		{"not markdown", "a.go", []string{"x := 1"}, []string{"// x := 1"}, []string{"x := 1"}},
	}
	for _, c := range cases {
		if got := notKept(c.path, c.undone, c.branchAdded); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

// I-0360: FindUndone reports nothing where the branch extends main's markdown
// line, and RecordedUndone drops a path recorded before the rule once the
// branch's line keeps main's words.
func TestUndoneLeavesAMarkdownLineTheBranchExtends(t *testing.T) {
	g := undoneGit("c3", "b3", []string{"flai-cli.md", "code.go"}, []string{"flai-cli.md", "code.go"})
	g.script("@@ -54 +54 @@\n-| `flai check` | old |\n+| `flai check` | main's |", diffArgs("b3", "c3", "flai-cli.md")...)
	g.script("@@ -54 +54 @@\n-| `flai check` | main's |\n+| `flai check` | main's and the story's |", diffArgs("c3", "HEAD", "flai-cli.md")...)
	g.script("@@ -1 +1 @@\n-x := 0\n+x := 1", diffArgs("b3", "c3", "code.go")...)
	g.script("@@ -1 +1 @@\n-x := 1\n+// x := 1", diffArgs("c3", "HEAD", "code.go")...)
	got, err := FindUndone(g, "/wt", "c3", "b3")
	if err != nil {
		t.Fatal(err)
	}
	if want := []Undone{{Path: "code.go", Lines: []string{"x := 1"}, OldBase: "b3"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("found %+v, want %+v", got, want)
	}

	repo, rg := undoneRepo(t)
	if err := recordUndone(repo.CacheDir(), "S-0001", []Undone{{Path: "flai-cli.md", Lines: []string{"| `flai check` | main's |"}, OldBase: "b0"}}, time.Now()); err != nil {
		t.Fatal(err)
	}
	rg.script("@@ -54 +54 @@\n-| `flai check` | main's |\n+| `flai check` | main's and the story's |", diffArgs("main", "story/S-0001", "flai-cli.md")...)
	if got, err := RecordedUndone(rg, repo, "S-0001"); err != nil || len(got) != 0 {
		t.Errorf("recorded and extended: %+v %v; want nothing still undone", got, err)
	}
}

// undoneGit scripts the check of the worktree wt, rebased onto the commit
// onto from old, where main changed mainPaths between them and the branch
// changes branchPaths against onto. It scripts no read of main's name, so a
// check that reads main again fails.
func undoneGit(onto, old string, mainPaths, branchPaths []string) scriptedGit {
	g := scriptedGit{}
	g.script(strings.Join(mainPaths, "\n"), nameOnlyArgs(old, onto)...)
	g.script(strings.Join(branchPaths, "\n"), nameOnlyArgs(onto, "HEAD")...)
	return g
}

// ADR-0155: a path where the branch removes lines main added is reported
// with them; one where it removes only its own lines, or keeps main's, is
// not; a path the branch deletes removes every line main added.
func TestFindUndone(t *testing.T) {
	g := undoneGit("c1", "b0", []string{"gone.md", "keeps.md", "own.md", "undone.md", "mainonly.md"}, []string{"gone.md", "keeps.md", "own.md", "undone.md"})
	g.script("@@ -1 +1,2 @@\n line\n+two", diffArgs("b0", "c1", "undone.md")...)
	g.script("@@ -2 +1,0 @@\n-two", diffArgs("c1", "HEAD", "undone.md")...)
	g.script("@@ -0,0 +1 @@\n+main's", diffArgs("b0", "c1", "own.md")...)
	g.script("@@ -3 +3 @@\n-branch's old\n+branch's new", diffArgs("c1", "HEAD", "own.md")...)
	g.script("@@ -0,0 +1 @@\n+kept", diffArgs("b0", "c1", "keeps.md")...)
	g.script("@@ -2,0 +2 @@\n+more", diffArgs("c1", "HEAD", "keeps.md")...)
	g.script("@@ -1 +1,2 @@\n+a\n+b", diffArgs("b0", "c1", "gone.md")...)
	g.script("diff --git a/gone.md b/gone.md\ndeleted file mode 100644\n--- a/gone.md\n+++ /dev/null\n@@ -1,3 +0,0 @@\n-a\n-b\n-c", diffArgs("c1", "HEAD", "gone.md")...)

	got, err := FindUndone(g, "/wt", "c1", "b0")
	if err != nil {
		t.Fatal(err)
	}
	want := []Undone{{Path: "gone.md", Lines: []string{"a", "b"}, OldBase: "b0"}, {Path: "undone.md", Lines: []string{"two"}, OldBase: "b0"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("found %+v, want %+v", got, want)
	}
}

// A rebase whose old base is the commit it took has nothing to check.
func TestFindUndoneChecksNothingWhenMainHasNotMoved(t *testing.T) {
	got, err := FindUndone(scriptedGit{}, "/wt", "b0", "b0")
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("found %+v, %v", got, err)
	}
}

// I-0338: main gains a commit, c2, while the sync runs, after the rebase
// took c1. The check diffs against c1 alone, so the line only c2 adds, which
// the rebased branch lacks, is none the branch removes; reading main again
// would see c2 and report the path.
func TestFindUndoneDiffsAgainstTheCommitTheRebaseTookNotMainsNewerHead(t *testing.T) {
	g := undoneGit("c1", "b0", []string{"a.md"}, []string{"a.md"})
	g.script("@@ -1 +1,2 @@\n line\n+main's", diffArgs("b0", "c1", "a.md")...)
	g.script("@@ -3,0 +4 @@\n+branch's", diffArgs("c1", "HEAD", "a.md")...)
	// main as it is once c2 lands: what a check reading main's name would see
	g.script("c2", "rev-parse", "main")
	g.script("a.md", nameOnlyArgs("b0", "main")...)
	g.script("a.md", nameOnlyArgs("main", "HEAD")...)
	g.script("@@ -1 +1,3 @@\n line\n+main's\n+c2's", diffArgs("b0", "main", "a.md")...)
	g.script("@@ -3 +2,0 @@\n-c2's\n@@ -3,0 +4 @@\n+branch's", diffArgs("main", "HEAD", "a.md")...)

	got, err := FindUndone(g, "/wt", "c1", "b0")
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("found %+v, %v; want nothing undone against the commit the rebase took", got, err)
	}
}

// I-0338, I-0335: Rebase reads main's head once, as c1, and rebases onto,
// checks, settles, and folds against c1, so main gaining c2 meanwhile, which
// a read of main's name would see, reports nothing undone; the answer still
// names the branch.
func TestRebaseChecksAndFoldsAgainstTheCommitItRebasedOnto(t *testing.T) {
	repo, g := undoneRepo(t)
	if err := os.MkdirAll(repo.WorktreePath("S-0001"), 0o755); err != nil {
		t.Fatal(err)
	}
	g.script("", "status", "--porcelain")
	g.script("c1", "rev-parse", "--verify", "main^{commit}")
	g.script("b0", "merge-base", "c1", "HEAD")
	g.script("h0", "rev-parse", "HEAD")
	g.script("", "rebase", "c1")
	g.script("", nameOnlyArgs("b0", "h0")...)
	for k, v := range undoneGit("c1", "b0", []string{"a.md"}, []string{"a.md"}) {
		g[k] = v
	}
	g.script("@@ -1 +1,2 @@\n line\n+main's", diffArgs("b0", "c1", "a.md")...)
	g.script("@@ -3,0 +4 @@\n+branch's", diffArgs("c1", "HEAD", "a.md")...)
	// main once c2 lands, which the rebased branch lacks
	g.script("c2", "rev-parse", "main")
	g.script("a.md", nameOnlyArgs("b0", "main")...)
	g.script("a.md", nameOnlyArgs("main", "HEAD")...)
	g.script("@@ -1 +1,3 @@\n line\n+main's\n+c2's", diffArgs("b0", "main", "a.md")...)
	g.script("@@ -3 +2,0 @@\n-c2's\n@@ -3,0 +4 @@\n+branch's", diffArgs("main", "HEAD", "a.md")...)
	var foldedOnto string
	o := SyncOptions{Runner: g, Repo: repo, Story: &workitem.Item{ID: "S-0001"}, Now: time.Now()}
	o.Files.Fold = func(base string) ([]Folded, []string, error) {
		foldedOnto = base
		return nil, nil, nil
	}

	res, err := Rebase(o)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || res.Base != "main" || len(res.UndoesMain) != 0 || foldedOnto != "c1" {
		t.Errorf("synced %v, base %q, undoes %+v, folded onto %q; want main, nothing undone, folded onto c1", res.Synced, res.Base, res.UndoesMain, foldedOnto)
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
