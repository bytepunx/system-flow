package digest

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
)

// repoLint is the repository's markdownlint configuration, .markdownlint.yaml,
// which the design folder's entries are linted with.
const repoLint = "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"

// lintClean fails the test with each finding markdownlint would report on
// content, a whole markdown file, under the repository's configuration.
func lintClean(t *testing.T, name, content string) {
	t.Helper()
	c, err := mdlint.Parse([]byte(repoLint), false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range c.Lint(content) {
		t.Errorf("%s:%d: %s", name, f.Line, f)
	}
	if t.Failed() {
		t.Logf("%s:\n%s", name, content)
	}
}

// seeded is the entries a first run over harvestRepo writes, in the order
// it writes them.
var seeded = []SeedEntry{
	{ID: "D-0001", Path: "design/digest/D-0001-which-shape-left.md", Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:05:00Z", Title: "Which shape: Left"},
	{ID: "D-0002", Path: "design/digest/D-0002-which-shape-narrow-too.md", Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:10:00Z", Title: "Which shape: Narrow, too"},
	{ID: "D-0003", Path: "design/digest/D-0003-i-0001-went-wrong.md", Kind: "cause", Source: "I-0001", Title: "I-0001 went wrong"},
	{ID: "D-0004", Path: "design/digest/D-0004-which-store-files-git-keeps-them.md", Kind: "recommendation", Source: "TH-0002", Ref: "2026-10-22T08:02:00Z", Title: "Which store: Files: git keeps them"},
}

// Seed writes one entry for each answer, overruled recommendation, and
// cause, two answers on one thread each their own, dated the day each was
// learned, citing its source, its body naming the answer's time, and
// linting clean; no decision is written.
func TestSeedWritesAnEntryForEachCandidate(t *testing.T) {
	r := harvestRepo(t)
	res, err := Seed(r, fakeGit{}, SeedOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res.Written, seeded) || len(res.Skipped) != 0 || res.DryRun {
		t.Fatalf("written %+v\nskipped %+v", res.Written, res.Skipped)
	}
	if want := map[string]int{"answer": 2, "recommendation": 1, "cause": 1}; !reflect.DeepEqual(res.Listed, want) {
		t.Errorf("listed %v, want %v", res.Listed, want)
	}
	if got := res.ByKind(); got["answer"] != (SeedTally{Listed: 2, Written: 2}) || got["cause"] != (SeedTally{Listed: 1, Written: 1}) || got["decision"] != (SeedTally{}) {
		t.Errorf("by kind %+v", got)
	}
	dates := []string{"2026-10-19", "2026-10-19", "2026-10-21", "2026-10-22"}
	says := []string{"answered TH-0001 at 2026-10-19T10:05:00Z.", "answered TH-0001 at 2026-10-19T10:10:00Z.", "I-0001 was closed on 2026-10-21.", "answered TH-0002 at 2026-10-22T08:02:00Z."}
	for i, w := range seeded {
		e, err := Get(r, w.ID)
		if err != nil {
			t.Fatal(err)
		}
		if e.Kind != w.Kind || !reflect.DeepEqual(e.Sources, []string{w.Source}) || e.Date != dates[i] || e.Status != "active" || !strings.Contains(e.Body, says[i]) || !strings.HasSuffix(e.Body, "Source: "+w.Source+".") {
			t.Errorf("%s: %+v", w.ID, e)
		}
		lintClean(t, w.Path, readFile(t, r, w.Path))
	}
	if e, _ := Get(r, "D-0004"); !strings.Contains(e.Body, "The recommendation it overruled:\n\n```text\nA database.\n```") {
		t.Errorf("the recommendation's body:\n%s", e.Body)
	}
}

// A second run finds every candidate recorded and writes nothing.
func TestSeedTwiceWritesNothingTheSecondTime(t *testing.T) {
	r := harvestRepo(t)
	if _, err := Seed(r, fakeGit{}, SeedOptions{}); err != nil {
		t.Fatal(err)
	}
	res, err := Seed(r, fakeGit{}, SeedOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var reasons []string
	for _, s := range res.Skipped {
		reasons = append(reasons, s.Source+" "+s.Reason)
	}
	want := []string{"TH-0001 cited by D-0001", "TH-0001 cited by D-0002", "I-0001 cited by D-0003", "TH-0002 cited by D-0004"}
	if len(res.Written) != 0 || !reflect.DeepEqual(reasons, want) {
		t.Errorf("second run wrote %+v, skipped %q; want %q", res.Written, reasons, want)
	}
	if files, _ := os.ReadDir(filepath.Join(r.Root, "design/digest")); len(files) != 4 {
		t.Errorf("%d files after two runs, want 4", len(files))
	}
}

// A dry run reports what it would write, with no IDs, and writes nothing.
func TestSeedDryRunWritesNothing(t *testing.T) {
	r := harvestRepo(t)
	res, err := Seed(r, fakeGit{}, SeedOptions{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	var want []SeedEntry
	for _, w := range seeded {
		w.ID, w.Path = "", ""
		want = append(want, w)
	}
	if !res.DryRun || !reflect.DeepEqual(res.Written, want) || len(res.Skipped) != 0 {
		t.Errorf("dry run: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(r.Root, "design/digest")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a dry run made the digest folder: %v", err)
	}
}

// An entry the check refuses, an answer with a key's shape, is skipped and
// named by the rule it broke, not its text, and the rest are written.
func TestSeedSkipsARefusedEntryAndGoesOn(t *testing.T) {
	r := harvestRepo(t)
	write(t, r.Root, "wip/threads/TH-0003-key.md", threadFile("TH-0003", "Which key?", "answered",
		"2026-10-20T09:00:00Z claude", "Which key does the deploy use?",
		"2026-10-20T09:01:00Z alex", "The key is AKIAIOSFODNN7EXAMPLE."))
	res, err := Seed(r, fakeGit{}, SeedOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want := []SeedSkip{{Kind: "answer", Source: "TH-0003", Ref: "2026-10-20T09:01:00Z", Reason: "refused: digest.secret"}}
	if !reflect.DeepEqual(res.Skipped, want) || len(res.Written) != 4 || res.Listed["answer"] != 3 {
		t.Errorf("written %+v\nskipped %+v", res.Written, res.Skipped)
	}
	files, _ := os.ReadDir(filepath.Join(r.Root, "design/digest"))
	for _, f := range files {
		if strings.Contains(readFile(t, r, "design/digest/"+f.Name()), "AKIA") {
			t.Errorf("%s holds the key", f.Name())
		}
	}
	if len(files) != 4 {
		t.Errorf("%d files, want 4", len(files))
	}
}

// A cause with no text or no fix is skipped; an issue is recorded by any
// entry citing it; a thread's candidate only by an entry citing the thread
// whose body holds the candidate's time.
func TestSkipReason(t *testing.T) {
	entries := []Entry{
		{ID: "D-0001", Sources: []string{"th-1"}, Body: "An answer of 2026-10-19T10:05:00Z."},
		{ID: "D-0002", Sources: []string{"I-1", "S-0001"}, Body: "Anything."},
		{ID: "D-0003", Sources: []string{"TH-0001"}, Body: "Also 2026-10-19T10:05:00Z."},
		{Sources: []string{"TH-0002"}, Body: "Written this run at 2026-10-22T08:02:00Z."},
	}
	for _, c := range []struct {
		c    Candidate
		want string
	}{
		{Candidate{Kind: "cause", Source: "I-0001", Fix: "f"}, "no cause"},
		{Candidate{Kind: "cause", Source: "I-0001", Text: "t", Fix: " "}, "no fix"},
		{Candidate{Kind: "cause", Source: "I-0001", Text: "t", Fix: "f"}, "cited by D-0002"},
		{Candidate{Kind: "cause", Source: "I-0002", Text: "t", Fix: "f"}, ""},
		{Candidate{Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:05:00Z"}, "cited by D-0001, D-0003"},
		{Candidate{Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:10:00Z"}, ""},
		{Candidate{Kind: "recommendation", Source: "TH-0002", Ref: "2026-10-22T08:02:00Z"}, "cited by an entry this run writes"},
		{Candidate{Kind: "answer", Source: "TH-0009", Ref: "2026-10-19T10:05:00Z"}, ""},
	} {
		if got := skipReason(c.c, entries); got != c.want {
			t.Errorf("skipReason(%+v) = %q, want %q", c.c, got, c.want)
		}
	}
}

// nasty is free text with what markdown would read as a heading, a table, a
// tab, trailing spaces, a run of blank lines, code fences of three and four
// backticks, a bold-only line, and a list.
const nasty = "# A heading in the text\n\n| a | b |\n|---|---|\n| 1 |\n\n\tA tabbed line   \n\n\n\nA fence:\n```go\nx := 1\n```\nAnd four: ````\n**Bold only**\n* a list\n"

// A body quotes each text in a fenced block that holds it whatever it is,
// and an entry with it lints clean.
func TestSeedBodyLintsClean(t *testing.T) {
	for _, c := range []Candidate{
		{Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:05:00Z", Question: nasty, Text: nasty},
		{Kind: "recommendation", Source: "TH-0001", Ref: "2026-10-19T10:05:00Z", Question: nasty, Text: nasty, Overruled: nasty},
		{Kind: "cause", Source: "I-0001", Date: "2026-10-21", Text: nasty, Fix: nasty},
	} {
		body := seedBody(c, "alex")
		if strings.Contains(body, "\t") || !strings.Contains(body, "`````text\n# A heading in the text\n\n| a | b |") || !strings.Contains(body, "\n    A tabbed line\n\nA fence:\n") {
			t.Errorf("%s body:\n%s", c.Kind, body)
		}
		lintClean(t, c.Kind, "---\nid: D-0001\ntitle: A fact\n---\n\n# D-0001 A fact\n\n"+body)
	}
	if body := seedBody(Candidate{Kind: "answer", Source: "TH-0001", Ref: "2026-10-19T10:05:00Z", Question: "Which?\n\nMore context.", Text: "This."}, "alex"); body != "The operator, alex, answered TH-0001 at 2026-10-19T10:05:00Z.\n\nWhat was asked:\n\n```text\nWhich?\n```\n\nThe answer:\n\n```text\nThis.\n```\n\nSource: TH-0001.\n" {
		t.Errorf("a plain answer's body:\n%s", body)
	}
}

// A question is quoted to its first paragraph, cut on a word near
// questionLimit.
func TestFirstParagraphIsCutOnAWord(t *testing.T) {
	long := firstParagraph(strings.Repeat("word ", 200) + "\n \nSecond paragraph.")
	if n := utf8.RuneCountInString(long); n > questionLimit+1 || !strings.HasSuffix(long, "word…") || strings.Contains(long, "Second") {
		t.Errorf("a long question, %d characters: %q", n, long)
	}
	if got := firstParagraph("  Short.\r\nStill the first.\r\n\r\nSecond."); got != "Short.\nStill the first." {
		t.Errorf("a short question: %q", got)
	}
}
