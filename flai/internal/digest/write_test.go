package digest

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeGit answers storygit.FolderNames as git would for a repository on
// main with one story branch, story/S-0009, whose digest folder holds
// onBranch, and no worktrees.
type fakeGit struct{ onBranch []string }

var errNoGit = errors.New("not a git command this fake answers")

func (f fakeGit) Run(_, _ string, args ...string) (string, error) {
	switch {
	case len(args) >= 2 && args[0] == "rev-parse" && args[1] == "--abbrev-ref":
		return "main", nil
	case len(args) >= 1 && args[0] == "for-each-ref":
		return "refs/heads/story/S-0009", nil
	case len(args) >= 3 && args[0] == "ls-tree" && args[2] == "refs/heads/story/S-0009":
		var out []string
		for _, n := range f.onBranch {
			out = append(out, "100644 blob 0123abcd\tdesign/digest/"+n)
		}
		return strings.Join(out, "\x00"), nil
	}
	return "", errNoGit
}

func (fakeGit) RunInput(_, _, _ string, _ ...string) (string, error) { return "", errNoGit }
func (fakeGit) LookPath(string) (string, error)                      { return "", errNoGit }

var everything = Env{Known: func(string) bool { return true }, Exists: func(string) bool { return true }}

func readFile(t *testing.T, r *workitem.Repo, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// ADR-0191: New allocates the first ID, names the file for the title's slug,
// writes the front matter in its order with the date in UTC, the sources in
// their padded form, the heading, and the body, creating the folder.
func TestNewWritesAnEntry(t *testing.T) {
	r := newRepo(t)
	now := time.Date(2026, 10, 23, 23, 30, 0, 0, time.FixedZone("west", -3*3600))
	res, err := New(r, fakeGit{}, NewOptions{
		Title: "  The digest:\n one file per entry ", Kind: "decision", Sources: []string{"s-428", "TH-476", "S-0428"},
		Topics: []string{"digest", " "}, Paths: []string{"flai/internal/digest"}, Body: "\nTH-0476 confirmed it.\n\n", Now: now, Env: everything,
	})
	if err != nil {
		t.Fatal(err)
	}
	wantPath := "design/digest/D-0001-the-digest-one-file-per-entry.md"
	if res.ID != "D-0001" || res.Path != wantPath || res.Title != "The digest: one file per entry" || res.Entry.Path != wantPath {
		t.Errorf("result %+v", res)
	}
	want := `---
id: D-0001
title: "The digest: one file per entry"
kind: decision
sources: [S-0428, TH-0476]
topics: [digest]
paths: [flai/internal/digest]
date: 2026-10-24
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0001 The digest: one file per entry

TH-0476 confirmed it.
`
	if got := readFile(t, r, wantPath); got != want {
		t.Errorf("file:\n%s\nwant:\n%s", got, want)
	}
	if e, err := Get(r, "D-1"); err != nil || e.Body != "TH-0476 confirmed it." || e.Status != "active" {
		t.Errorf("Get = %+v, %v", e, err)
	}
}

// ADR-0191: the next ID is one more than the highest anywhere, here or on a
// story branch; a gap is not filled, and what is not an entry is not counted.
func TestNextNumberIsOneMoreThanTheHighestAnywhere(t *testing.T) {
	r := newRepo(t)
	write(t, r.Root, "design/digest/README.md", "# Digest\n")
	write(t, r.Root, "design/digest/D-0001-one.md", "x\n")
	write(t, r.Root, "design/digest/D-0004-four.md", "x\n")
	write(t, r.Root, "design/digest/notes-D-0099.md", "x\n")
	if n, err := NextNumber(r, fakeGit{}); err != nil || n != 5 {
		t.Errorf("here only: %d, %v; want 5", n, err)
	}
	if n, err := NextNumber(r, fakeGit{onBranch: []string{"D-0002-two.md", "D-0007-on-a-branch.md"}}); err != nil || n != 8 {
		t.Errorf("with a story branch: %d, %v; want 8", n, err)
	}
	if n, err := NextNumber(newRepo(t), nil); err != nil || n != 1 {
		t.Errorf("no folder, no git: %d, %v; want 1", n, err)
	}
	res, err := New(r, fakeGit{onBranch: []string{"D-0007-on-a-branch.md"}}, NewOptions{Title: "Eight", Kind: "cause", Sources: []string{"I-0001"}, Env: everything})
	if err != nil || res.ID != "D-0008" {
		t.Errorf("New after a branch's D-0007: %+v, %v", res, err)
	}
}

// ADR-0191: a slug is cut at 80 characters, on a word.
func TestSlugIsCutOnAWord(t *testing.T) {
	long := strings.Repeat("abcdefghi ", 12)
	s := slug(long)
	if len(s) > 80 || strings.HasSuffix(s, "-") || !strings.HasSuffix(s, "abcdefghi") {
		t.Errorf("slug %q", s)
	}
	if slug("!!!") != "entry" {
		t.Errorf("a title with no letters: %q", slug("!!!"))
	}
}

// ADR-0191: a write the check refuses leaves nothing written, not the file
// and not a folder it created, and names each problem.
func TestNewRefusesAndLeavesNothing(t *testing.T) {
	for _, c := range []struct {
		name string
		opts NewOptions
		says []string
	}{
		{"unknown kind and no sources", NewOptions{Title: "A fact", Kind: "rumour", Env: everything}, []string{"kind is \"rumour\"", "sources names nothing"}},
		{"a source that names nothing", NewOptions{Title: "A fact", Kind: "answer", Sources: []string{"S-0999"}, Env: Env{Known: func(string) bool { return false }}}, []string{"S-0999 names nothing"}},
		{"a secret in the body", NewOptions{Title: "A fact", Kind: "answer", Sources: []string{"TH-0001"}, Body: "The key is AKIAIOSFODNN7EXAMPLE.", Env: everything}, []string{"digest.secret", "never its value"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := newRepo(t)
			_, err := New(r, fakeGit{}, c.opts)
			var refused *RefusedError
			if !errors.As(err, &refused) {
				t.Fatalf("want a refusal, got %v", err)
			}
			for _, s := range c.says {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("error does not say %q:\n%v", s, err)
				}
			}
			if _, err := os.Stat(filepath.Join(r.Root, "design")); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("the folders New created are still there: %v", err)
			}
		})
	}
	r := newRepo(t)
	write(t, r.Root, "design/digest/D-0001-kept.md", good)
	if _, err := New(r, fakeGit{}, NewOptions{Title: "A fact", Kind: "rumour", Sources: []string{"S-0001"}, Env: everything}); err == nil {
		t.Fatal("want a refusal")
	}
	entries, _ := os.ReadDir(filepath.Join(r.Root, "design/digest"))
	if len(entries) != 1 || entries[0].Name() != "D-0001-kept.md" {
		t.Errorf("a folder that was there keeps only what it had: %v", entries)
	}
}

// edited is an entry written by hand, with a block list, which Edit keeps
// byte for byte except for what it is given.
const edited = `---
id: D-0003
title: The digest is one file per entry
kind: decision
sources: [S-0428]
topics:
  - old
  - older
paths: []
date: 2026-10-23
status: active
absorbed_by: ""
cited: [S-0001]
---

# D-0003 The digest is one file per entry

The first body.
`

// ADR-0191: Edit changes only the fields it is given, adds asked_again
// where the entry has none, replaces the body below the heading, and writes
// nothing for a change that changes nothing.
func TestEditChangesOnlyWhatItIsGiven(t *testing.T) {
	r := newRepo(t)
	rel := "design/digest/" + goodName
	write(t, r.Root, rel, edited)
	topics, asked, body := []string{"digest", "prime"}, []string{"th-5"}, "A new body.\n\nWith two paragraphs."
	res, err := Edit(r, "D-3", EditOptions{Topics: &topics, AskedAgain: &asked, Body: &body, Env: everything})
	if err != nil {
		t.Fatal(err)
	}
	want := strings.NewReplacer("topics:\n  - old\n  - older\n", "topics: [digest, prime]\n",
		"absorbed_by: \"\"\n", "absorbed_by: \"\"\nasked_again: [TH-0005]\n",
		"The first body.\n", "A new body.\n\nWith two paragraphs.\n").Replace(edited)
	if got := readFile(t, r, rel); got != want {
		t.Errorf("file:\n%s\nwant:\n%s", got, want)
	}
	if res.ID != "D-0003" || res.Path != rel || res.Entry.Body != body || len(res.Entry.AskedAgain) != 1 {
		t.Errorf("result %+v", res)
	}
	if _, err := Edit(r, "D-0003", EditOptions{Topics: &topics, Env: everything}); err != nil || readFile(t, r, rel) != want {
		t.Errorf("an edit that changes nothing: %v", err)
	}
	if _, err := Edit(r, "D-0004", EditOptions{}); err == nil {
		t.Error("an entry that is not there is an error")
	}
}

// ADR-0191: an edit the check refuses puts the file back as it was, and a
// problem the entry had before the edit does not refuse it.
func TestEditRestoresTheFileOnARefusal(t *testing.T) {
	r := newRepo(t)
	rel := "design/digest/" + goodName
	write(t, r.Root, rel, edited)
	none, secret, unknown := []string{}, "The token: ghp_"+strings.Repeat("Zz9", 12), []string{"S-0999"}
	for name, opts := range map[string]EditOptions{
		"no sources":             {Sources: &none, Env: everything},
		"a secret in the body":   {Body: &secret, Env: everything},
		"a source names nothing": {Sources: &unknown, Env: env},
	} {
		var refused *RefusedError
		if _, err := Edit(r, "D-0003", opts); !errors.As(err, &refused) || len(refused.Problems) != 1 {
			t.Errorf("%s: want one problem refused, got %v", name, err)
		}
		if got := readFile(t, r, rel); got != edited {
			t.Errorf("%s: the file was not put back:\n%s", name, got)
		}
	}
	write(t, r.Root, rel, strings.Replace(edited, "kind: decision", "kind: rumour", 1))
	paths := []string{"flai/internal/digest"}
	if _, err := Edit(r, "D-0003", EditOptions{Paths: &paths, Env: everything}); err != nil {
		t.Errorf("a problem the entry already had refused the edit: %v", err)
	}
	if got := readFile(t, r, rel); !strings.Contains(got, "paths: [flai/internal/digest]\n") {
		t.Errorf("paths not written:\n%s", got)
	}
}
