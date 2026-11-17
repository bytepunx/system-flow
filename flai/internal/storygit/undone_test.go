package storygit

import (
	"encoding/json"
	"errors"
	"fmt"
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

// ADR-0309: each hunk's lines are numbered from its header, with or without
// a count, removed lines in the old file and added ones in the new; file
// headers and git's note of a missing newline are not lines.
func TestDiffHunksNumbersLinesFromEachHunksHeader(t *testing.T) {
	diff := "diff --git a/a.md b/a.md\nindex 1..2 100644\n--- a/a.md\n+++ b/a.md\n" +
		"@@ -1 +1,2 @@\n-old  \n+--- new\n+++ also\n" +
		"@@ -7,2 +8,0 @@ func f() {\n-\t}\n-}\n" +
		"@@ -20,0 +20 @@\n+last\n\\ No newline at end of file"
	added, removed, err := diffHunks(diff)
	if err != nil {
		t.Fatal(err)
	}
	if want := []numbered{{1, "--- new"}, {2, "++ also"}, {20, "last"}}; !reflect.DeepEqual(added, want) {
		t.Errorf("added %v, want %v", added, want)
	}
	if want := []numbered{{1, "old"}, {7, "\t}"}, {8, "}"}}; !reflect.DeepEqual(removed, want) {
		t.Errorf("removed %v, want %v", removed, want)
	}
	if _, _, err := diffHunks("@@ -x +1 @@\n+a"); err == nil || !strings.Contains(err.Error(), "@@ -x +1 @@") {
		t.Errorf("a header with no line number read: %v", err)
	}
}

// ADR-0309: main's line is undone only where the branch removes the line of
// the number main added it at, whatever the text the branch removes elsewhere.
func TestAtPositionsCountsMainsLinesTheBranchRemovesWhereMainAddedThem(t *testing.T) {
	for _, c := range []struct {
		name           string
		added, removed []numbered
		want           []string
	}{
		{"removes main's line", []numbered{{4, "two"}}, []numbered{{4, "two"}}, []string{"two"}},
		{"same text elsewhere", []numbered{{12, "\t}"}}, []numbered{{30, "\t}"}}, []string{}},
		{"removes some of main's lines", []numbered{{1, "a"}, {2, "b"}, {3, "c"}}, []numbered{{3, "c"}, {1, "a"}}, []string{"a", "c"}},
		{"keeps main's change", []numbered{{1, "a"}}, nil, []string{}},
	} {
		if got := atPositions(c.added, c.removed); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

// I-0360: in a markdown file, main's line is kept when a line the branch adds
// holds every word of it, punctuation at the words' ends left out, in any
// order, each added line keeping one of main's at most; a dropped or changed
// word, a line restoring one main removed, a line with no words, and a file
// that is not markdown keep the multiset of whole lines.
func TestNotKeptLeavesMainsLinesABranchLineExtends(t *testing.T) {
	row := "| `flai check` | Validates the repo. |"
	cases := []struct {
		name, path                       string
		undone, branchAdded, mainRemoved []string
		want                             []string
	}{
		{"extended row", "a.md", []string{row}, []string{"| `flai check` | Validates the repo. It also checks images. |"}, nil, []string{}},
		{"word inserted inside", "a.md", []string{"one two three"}, []string{"one and two three"}, nil, []string{}},
		{"word moved", "a.md", []string{"a, and b, c"}, []string{"a, b, and c, d"}, nil, []string{}},
		{"bracket dropped", "a.md", []string{"(ADR-0256), ADR-0296)"}, []string{"(ADR-0256, ADR-0296)"}, nil, []string{}},
		{"upper-case extension", "A.MD", []string{"one two"}, []string{"one two three"}, nil, []string{}},
		{"word dropped", "a.md", []string{"one two three"}, []string{"one three four"}, nil, []string{"one two three"}},
		{"word changed", "a.md", []string{"one two"}, []string{"one too"}, nil, []string{"one two"}},
		{"word once of twice", "a.md", []string{"a a b"}, []string{"a b c"}, nil, []string{"a a b"}},
		{"revert adds nothing", "a.md", []string{"one two"}, nil, nil, []string{"one two"}},
		{"restores main's removed line", "a.md", []string{"one two"}, []string{"one two three"}, []string{"one two three"}, []string{"one two"}},
		{"one added line keeps one", "a.md", []string{"one two", "one two"}, []string{"one two three"}, nil, []string{"one two"}},
		{"each added line keeps one", "a.md", []string{"a b", "c d"}, []string{"c d e", "a b f"}, nil, []string{}},
		{"blank line", "a.md", []string{""}, []string{"anything"}, nil, []string{""}},
		{"punctuation alone", "a.md", []string{"---"}, []string{"anything"}, nil, []string{"---"}},
		{"not markdown", "a.go", []string{"x := 1"}, []string{"// x := 1"}, nil, []string{"x := 1"}},
	}
	for _, c := range cases {
		if got := notKept(c.path, c.undone, c.branchAdded, c.mainRemoved); !reflect.DeepEqual(got, c.want) {
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
	md := []string{"flai-cli.md"}
	rg.script(fileDiff("flai-cli.md", "@@ -54 +54 @@\n-| `flai check` | main's |\n+| `flai check` | main's and the story's |"), diffAllArgs("main", "story/S-0001", md)...)
	rg.script(fileDiff("flai-cli.md", "@@ -54 +54 @@\n-| `flai check` | old |\n+| `flai check` | main's |"), diffAllArgs("b0", "main", md)...)
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
	// b.md's branch removes nothing, so main's diff of it is not read
	g.script(fileDiff("a.md", "@@ -1 +0,0 @@\n-y"), diffAllArgs("main", "story/S-0001", []string{"a.md", "b.md"})...)
	g.script(fileDiff("a.md", "@@ -0,0 +1,2 @@\n+y\n+x"), diffAllArgs("b0", "main", []string{"a.md"})...)

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

	g.script("", diffAllArgs("main", "story/S-0001", []string{"a.md"})...)
	if got, err := RecordedUndone(g, repo, "S-0001"); err != nil || len(got) != 0 {
		t.Fatalf("restored: %+v %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(repo.CacheDir(), "undone", "S-0001.json")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the empty record is left: %v", err)
	}
}

// The Go file of I-0385 and I-0124: board.go's old base holds a loop, and
// main adds a function after it, at lines 8 to 13 of its file, ending in a
// tab and a brace and in a brace.
const (
	goMainAdds    = "@@ -7,0 +8,6 @@\n+\n+func a() {\n+\tif ok {\n+\t\treturn\n+\t}\n+}"
	goLoopReplace = "@@ -4,3 +4 @@\n-\tfor _, x := range xs {\n-\t\tuse(x)\n-\t}\n+\tuseAll(xs)"
	goRevert      = "@@ -8,6 +7,0 @@\n-\n-func a() {\n-\tif ok {\n-\t\treturn\n-\t}\n-}"
)

// goMainLines are the lines goMainAdds adds, as text.
var goMainLines = []string{"", "func a() {", "\tif ok {", "\t\treturn", "\t}", "}"}

// I-0385, ADR-0309: a branch that replaces board.go's loop removes the loop's
// closing brace, a tab and a brace, at line 6, which the text compare matched
// to the one main added at line 12: it undoes nothing. I-0124: a branch whose
// revert, replayed over main's change, removes main's function from
// revert.go, at the lines main added it, undoes every line of it.
func TestFindUndoneCountsMainsLinesByTheirPosition(t *testing.T) {
	g := undoneGit("c1", "b0", []string{"board.go", "revert.go"}, []string{"board.go", "revert.go"})
	g.script(goMainAdds, diffArgs("b0", "c1", "board.go")...)
	g.script(goLoopReplace, diffArgs("c1", "HEAD", "board.go")...)
	g.script(goMainAdds, diffArgs("b0", "c1", "revert.go")...)
	g.script(goRevert, diffArgs("c1", "HEAD", "revert.go")...)

	got, err := FindUndone(g, "/wt", "c1", "b0")
	if err != nil {
		t.Fatal(err)
	}
	if want := []Undone{{Path: "revert.go", Lines: goMainLines, OldBase: "b0"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("found %+v, want %+v", got, want)
	}
}

// I-0385, ADR-0309: board.go, recorded by the text compare for the brace the
// branch removes at another line than main added it, drops out of the record
// at its next read; revert.go, whose branch removes main's lines where main
// added them, stays.
func TestRecordedUndoneDropsAPathTheTextCompareRecorded(t *testing.T) {
	repo, g := undoneRepo(t)
	found := []Undone{{Path: "board.go", Lines: []string{"\t}"}, OldBase: "b0"}, {Path: "revert.go", Lines: goMainLines, OldBase: "b0"}}
	if err := recordUndone(repo.CacheDir(), "S-0002", found, time.Date(2026, 11, 15, 18, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	both := []string{"board.go", "revert.go"}
	g.script(fileDiff("board.go", goLoopReplace)+"\n"+fileDiff("revert.go", goRevert), diffAllArgs("main", "story/S-0002", both)...)
	g.script(fileDiff("board.go", goMainAdds)+"\n"+fileDiff("revert.go", goMainAdds), diffAllArgs("b0", "main", both)...)

	got, err := RecordedUndone(g, repo, "S-0002")
	if err != nil {
		t.Fatal(err)
	}
	if want := found[1:]; !reflect.DeepEqual(got, want) {
		t.Errorf("still undone %+v, want %+v", got, want)
	}
	if rec := readRecord(t, repo, "S-0002"); len(rec.Paths) != 1 || rec.Paths[0].Path != "revert.go" {
		t.Errorf("record not rewritten without board.go: %+v", rec)
	}
}

// fileDiff is what git prints of path in a diff of several files: its
// header, with the tab git puts after a name holding a space, then hunks.
func fileDiff(path, hunks string) string {
	end := ""
	if strings.Contains(path, " ") {
		end = "\t"
	}
	return "diff --git a/" + path + " b/" + path + "\nindex 1111111..2222222 100644\n--- a/" + path + end + "\n+++ b/" + path + end + "\n" + hunks
}

// countingGit is a scriptedGit that counts the git processes it runs, by
// their subcommand.
type countingGit struct {
	scriptedGit
	runs map[string]int
}

func (g countingGit) Run(dir, name string, args ...string) (string, error) {
	sub := args[0]
	if sub == "-c" {
		sub = args[2]
	}
	g.runs[sub]++
	return g.scriptedGit.Run(dir, name, args...)
}

// I-0349: a record of 20 paths with one old base is read with one git diff
// of the branch and one of main since that base, beside the rev-parse of the
// main branch, not one or two diffs per path; the answer is each path's.
func TestRecordedUndoneRunsOneDiffForTheBranchAndOnePerOldBase(t *testing.T) {
	repo, s := undoneRepo(t)
	g := countingGit{s, map[string]int{}}
	var found []Undone
	var paths []string
	branch, mains := "", ""
	for i := range 20 {
		p := fmt.Sprintf("p%02d.md", i)
		paths = append(paths, p)
		found = append(found, Undone{Path: p, Lines: []string{"main's " + p}, OldBase: "b0"})
		branch += fileDiff(p, "@@ -3 +2,0 @@\n-main's "+p) + "\n"
		mains += fileDiff(p, "@@ -2,0 +3 @@\n+main's "+p) + "\n"
	}
	if err := recordUndone(repo.CacheDir(), "S-0003", found, time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	g.script(branch, diffAllArgs("main", "story/S-0003", paths)...)
	g.script(mains, diffAllArgs("b0", "main", paths)...)

	got, err := RecordedUndone(g, repo, "S-0003")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, found) {
		t.Errorf("still undone %+v, want %+v", got, found)
	}
	if want := map[string]int{"rev-parse": 1, "diff": 2}; !reflect.DeepEqual(g.runs, want) {
		t.Errorf("git runs %v, want %v", g.runs, want)
	}
}

// I-0349: a record whose paths were found at two old bases reads main once
// per base, and only for the paths the branch still removes lines of; a path
// with a space is split out of the diff by its header, and a path the branch
// no longer changes drops out.
func TestRecordedUndoneReadsMainOncePerOldBase(t *testing.T) {
	repo, s := undoneRepo(t)
	g := countingGit{s, map[string]int{}}
	found := []Undone{
		{Path: "a.md", Lines: []string{"x"}, OldBase: "b0"},
		{Path: "d.go", Lines: goMainLines, OldBase: "b1"},
		{Path: "gone.md", Lines: []string{"g"}, OldBase: "b1"},
		{Path: "with space.md", Lines: []string{"y"}, OldBase: "b0"},
	}
	if err := recordUndone(repo.CacheDir(), "S-0004", found, time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	g.script(fileDiff("a.md", "@@ -2 +1,0 @@\n-x")+"\n"+fileDiff("d.go", goRevert)+"\n"+fileDiff("with space.md", "@@ -1 +0,0 @@\n-y"),
		diffAllArgs("main", "story/S-0004", []string{"a.md", "d.go", "gone.md", "with space.md"})...)
	g.script(fileDiff("a.md", "@@ -1,0 +2 @@\n+x")+"\n"+fileDiff("with space.md", "@@ -0,0 +1 @@\n+y"), diffAllArgs("b0", "main", []string{"a.md", "with space.md"})...)
	g.script(fileDiff("d.go", goMainAdds), diffAllArgs("b1", "main", []string{"d.go"})...)

	got, err := RecordedUndone(g, repo, "S-0004")
	if err != nil {
		t.Fatal(err)
	}
	if want := []Undone{found[0], found[1], found[3]}; !reflect.DeepEqual(got, want) {
		t.Errorf("still undone %+v, want %+v", got, want)
	}
	if want := map[string]int{"rev-parse": 1, "diff": 3}; !reflect.DeepEqual(g.runs, want) {
		t.Errorf("git runs %v, want %v", g.runs, want)
	}
	if rec := readRecord(t, repo, "S-0004"); len(rec.Paths) != 3 {
		t.Errorf("record not rewritten without gone.md: %+v", rec)
	}
}

// I-0349: a path git quotes in its header, here one holding a quote, matches
// no header of the batched diff, so it is read with a diff of its own, and a
// deleted file is split out by its header as any other.
func TestReadHunksOfReadsAQuotedPathOnItsOwn(t *testing.T) {
	g := scriptedGit{}
	quoted := `q"x.md`
	paths := []string{"del.md", quoted}
	g.script("diff --git a/del.md b/del.md\ndeleted file mode 100644\nindex 1111111..0000000\n--- a/del.md\n+++ /dev/null\n@@ -1,2 +0,0 @@\n-a\n-b\n"+
		`diff --git "a/q\"x.md" "b/q\"x.md"`+"\nindex 1111111..2222222 100644\n@@ -1 +1 @@\n-old\n+new",
		diffAllArgs("b0", "c1", paths)...)
	g.script("@@ -1 +1 @@\n-old\n+new", diffArgs("b0", "c1", quoted)...)

	got, err := readHunksOf(g, "/wt", "b0", "c1", paths)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]hunks{
		"del.md": {removed: []numbered{{1, "a"}, {2, "b"}}},
		quoted:   {added: []numbered{{1, "new"}}, removed: []numbered{{1, "old"}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("hunks %+v, want %+v", got, want)
	}
}
