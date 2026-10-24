package digest

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// asked is TH-0003's question, as claude asked it of alex.
var asked = []string{"2026-10-20T09:00:00Z claude", "Which colour?"}

// saved writes TH-0003, titled "Which colour?", with asked's question and
// the entries given after it, as a reply saves it, and reads it back.
func saved(t *testing.T, r *workitem.Repo, entries ...string) *threads.Thread {
	t.Helper()
	write(t, r.MainRoot, "wip/threads/TH-0003-colour.md", threadFile("TH-0003", "Which colour?", "answered", append(copied(asked), entries...)...))
	th, err := threads.Get(r, "TH-0003")
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// copied is a copy of list, so appending to it leaves list alone.
func copied(list []string) []string { return append([]string(nil), list...) }

// digestFiles is the names of the files in the digest folder under root.
func digestFiles(t *testing.T, root string) []string {
	t.Helper()
	files, err := os.ReadDir(filepath.Join(root, "design/digest"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	var out []string
	for _, f := range files {
		out = append(out, f.Name())
	}
	return out
}

// seedsNothingFor fails the test when Seed would write an entry for the
// thread's answer of ref, so the entry OnAnswer wrote is the only one.
func seedsNothingFor(t *testing.T, r *workitem.Repo, ref string) {
	t.Helper()
	res, err := Seed(r, fakeGit{}, SeedOptions{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range res.Written {
		if w.Source == "TH-0003" {
			t.Errorf("seed would write %+v", w)
		}
	}
	for _, s := range res.Skipped {
		if s.Source == "TH-0003" && s.Ref == ref {
			return
		}
	}
	t.Errorf("seed does not find TH-0003's answer of %s recorded: %+v", ref, res.Skipped)
}

// A first answer on a thread someone else opened is written as Seed would
// write it: its title, body, kind, source, date, topics, and paths; and a
// later seed finds it recorded and writes no second.
func TestOnAnswerWritesAFirstAnswer(t *testing.T) {
	r := harvestRepo(t)
	th := saved(t, r, "2026-10-20T09:01:00Z alex", "Blue. It matches the logo.")
	got := OnAnswer(r, fakeGit{}, th)
	want := &AnswerEntry{ID: "D-0001", Path: "design/digest/D-0001-which-colour-blue.md", Action: ActionWritten}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("OnAnswer = %+v, want %+v", got, want)
	}
	c := Answers(th, "alex", Lookup(r))[0]
	e, err := Get(r, "D-0001")
	if err != nil {
		t.Fatal(err)
	}
	if e.Title != seedTitle(c.Title) || e.Kind != "answer" || !reflect.DeepEqual(e.Sources, []string{"TH-0003"}) || e.Date != "2026-10-20" || e.Status != "active" ||
		!reflect.DeepEqual(e.Topics, c.Topics) || !reflect.DeepEqual(e.Paths, c.Paths) || e.Body != strings.TrimSpace(seedBody(c, "alex")) {
		t.Errorf("entry %+v\nfrom %+v", e, c)
	}
	lintClean(t, e.Path, readFile(t, r, e.Path))
	seedsNothingFor(t, r, "2026-10-20T09:01:00Z")
}

// A second answer on the same thread is appended to the entry the first
// wrote, under a sentence holding its time, before its source line, keeping
// all it said; no second file is written, the same answer again writes
// nothing, and a later seed finds both recorded.
func TestOnAnswerUpdatesTheEntryOnASecondAnswer(t *testing.T) {
	r := harvestRepo(t)
	first := []string{"2026-10-20T09:01:00Z alex", "Blue."}
	OnAnswer(r, fakeGit{}, saved(t, r, first...))
	before, err := Get(r, "D-0001")
	if err != nil {
		t.Fatal(err)
	}
	th := saved(t, r, append(copied(first), "2026-10-20T09:05:00Z claude", "Which blue?", "2026-10-20T09:06:00Z alex", "Navy.")...)
	got := OnAnswer(r, fakeGit{}, th)
	if want := (&AnswerEntry{ID: "D-0001", Path: before.Path, Action: ActionUpdated}); !reflect.DeepEqual(got, want) {
		t.Fatalf("OnAnswer = %+v, want %+v", got, want)
	}
	after, err := Get(r, "D-0001")
	if err != nil {
		t.Fatal(err)
	}
	kept := strings.TrimSuffix(before.Body, "Source: TH-0003.")
	if want := kept + "The operator, alex, answered TH-0003 again at 2026-10-20T09:06:00Z.\n\nThe answer:\n\n```text\nNavy.\n```\n\nSource: TH-0003."; after.Body != want {
		t.Errorf("body:\n%s\nwant:\n%s", after.Body, want)
	}
	if after.Title != before.Title || after.Kind != before.Kind {
		t.Errorf("the update changed the front matter: %+v", after)
	}
	lintClean(t, after.Path, readFile(t, r, after.Path))
	if files := digestFiles(t, r.Root); len(files) != 1 {
		t.Errorf("files %v, want one", files)
	}
	if again := OnAnswer(r, fakeGit{}, th); again != nil {
		t.Errorf("the same answer again: %+v", again)
	}
	if body, _ := Get(r, "D-0001"); body.Body != after.Body {
		t.Errorf("the same answer again changed the body:\n%s", body.Body)
	}
	seedsNothingFor(t, r, "2026-10-20T09:01:00Z")
	seedsNothingFor(t, r, "2026-10-20T09:06:00Z")
}

// digestEntry is an entry file with the sources, status, and body given.
func digestEntry(id, slug, sources, status, absorbedBy, body string) (string, string) {
	return "design/digest/" + id + "-" + slug + ".md", "---\nid: " + id + "\ntitle: " + slug + "\nkind: answer\nsources: [" + sources + "]\ntopics: []\npaths: []\ndate: 2026-10-20\nstatus: " + status + "\nabsorbed_by: \"" + absorbedBy + "\"\nasked_again: []\ncited: []\n---\n\n# " + id + " " + slug + "\n\n" + body + "\n"
}

// When the curator merged the thread's entry into another, the thread's
// source moved there, the answer updates that entry, the highest numbered
// active one citing the thread, and writes none; when only absorbed or
// compacted entries cite it, the answer is a new fact and is written.
func TestOnAnswerUpdatesTheEntryTheCuratorMergedInto(t *testing.T) {
	r := harvestRepo(t)
	write(t, r.Root, "design/system/colour.md", "# Colour\n")
	for _, e := range [][]string{
		{"D-0001", "which-colour-blue", "S-0001", "compacted", "", "The first answer, its source moved to D-0003."},
		{"D-0002", "an-older-merge", "TH-0003", "active", "", "An older merge."},
		{"D-0003", "colours-and-logos", "S-0001, TH-0003", "active", "", "Two facts the curator merged: blue, and the logo.\n\nSources: S-0001, TH-0003."},
		{"D-0004", "colour-in-the-design", "TH-0003", "absorbed", "design/system/colour.md", "The design records it."},
	} {
		p, content := digestEntry(e[0], e[1], e[2], e[3], e[4], e[5])
		write(t, r.Root, p, content)
	}
	got := OnAnswer(r, fakeGit{}, saved(t, r, "2026-10-20T09:01:00Z alex", "Navy, now."))
	if got == nil || got.ID != "D-0003" || got.Action != ActionUpdated || got.Warning != "" {
		t.Fatalf("OnAnswer = %+v, want D-0003 updated", got)
	}
	e, _ := Get(r, "D-0003")
	if !strings.HasPrefix(e.Body, "Two facts the curator merged: blue, and the logo.\n\nThe operator, alex, answered TH-0003 again at 2026-10-20T09:01:00Z.") || !strings.HasSuffix(e.Body, "```\n\nSources: S-0001, TH-0003.") {
		t.Errorf("D-0003's body:\n%s", e.Body)
	}
	if files := digestFiles(t, r.Root); len(files) != 4 {
		t.Errorf("files %v, want the four there were", files)
	}

	only := harvestRepo(t)
	p, content := digestEntry("D-0004", "colour-in-the-design", "TH-0003", "absorbed", "design/system/colour.md", "The design records it.")
	write(t, only.Root, p, content)
	write(t, only.Root, "design/system/colour.md", "# Colour\n")
	if got := OnAnswer(only, fakeGit{}, saved(t, only, "2026-10-20T09:01:00Z alex", "Navy, now.")); got == nil || got.ID != "D-0005" || got.Action != ActionWritten {
		t.Errorf("with only an absorbed entry citing the thread: %+v, want D-0005 written", got)
	}
}

// A confirmation is written with the confirmed recommendation's text.
func TestOnAnswerWritesAConfirmedRecommendation(t *testing.T) {
	r := harvestRepo(t)
	got := OnAnswer(r, fakeGit{}, saved(t, r,
		"2026-10-20T09:01:00Z orchestrator (recommendation)", "Green, as the plan says.",
		"2026-10-20T09:02:00Z alex", "Confirmed the recommendation of 2026-10-20T09:01:00Z orchestrator."))
	if got == nil || got.Action != ActionWritten {
		t.Fatalf("OnAnswer = %+v", got)
	}
	e, _ := Get(r, got.ID)
	if e.Kind != "answer" || e.Title != "Which colour: Green, as the plan says" || !strings.Contains(e.Body, "The answer:\n\n```text\nGreen, as the plan says.\n```") {
		t.Errorf("entry %+v", e)
	}
}

// An owner's reply over a pending recommendation is a recommendation entry
// naming the overruled text.
func TestOnAnswerWritesAnOverruledRecommendation(t *testing.T) {
	r := harvestRepo(t)
	got := OnAnswer(r, fakeGit{}, saved(t, r,
		"2026-10-20T09:01:00Z orchestrator (recommendation)", "Green, as the plan says.",
		"2026-10-20T09:02:00Z alex", "Red: green clashes."))
	if got == nil || got.Action != ActionWritten {
		t.Fatalf("OnAnswer = %+v", got)
	}
	e, _ := Get(r, got.ID)
	if e.Kind != "recommendation" || !strings.Contains(e.Body, "The recommendation it overruled:\n\n```text\nGreen, as the plan says.\n```\n\nThe operator's answer:\n\n```text\nRed: green clashes.\n```") {
		t.Errorf("entry %+v", e)
	}
	lintClean(t, e.Path, readFile(t, r, e.Path))
}

// A thread whose newest entry is no operator's answer writes nothing.
func TestOnAnswerWritesNothingForWhatIsNoAnswer(t *testing.T) {
	answer := []string{"2026-10-20T09:01:00Z alex", "Blue."}
	for _, c := range []struct {
		name    string
		file    string
		entries []string
	}{
		{name: "an agent's reply", entries: append(copied(answer), "2026-10-20T09:02:00Z claude", "Blue it is.")},
		{name: "a recommendation", entries: []string{"2026-10-20T09:01:00Z orchestrator (recommendation)", "Green."}},
		{name: "a resolution after the answer", entries: append(copied(answer), "2026-10-20T09:02:00Z alex", "Resolved.")},
		{name: "the operator's reply on their own thread", file: threadFile("TH-0003", "Which colour?", "answered",
			"2026-10-20T09:00:00Z alex", "Which colour?", "2026-10-20T09:01:00Z claude", "Blue?", "2026-10-20T09:02:00Z alex", "Blue.")},
		{name: "a permission thread", file: threadFile("TH-0003", "Allow Write .claude/settings.json?", "answered",
			"2026-10-20T09:00:00Z claude", "May I?", "2026-10-20T09:01:00Z alex", "yes")},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := harvestRepo(t)
			th := saved(t, r, c.entries...)
			if c.file != "" {
				write(t, r.Root, "wip/threads/TH-0003-colour.md", c.file)
				th, _ = threads.Get(r, "TH-0003")
			}
			if got := OnAnswer(r, fakeGit{}, th); got != nil {
				t.Errorf("OnAnswer = %+v, want nil", got)
			}
			if files := digestFiles(t, r.Root); len(files) != 0 {
				t.Errorf("wrote %v", files)
			}
		})
	}
	if got := OnAnswer(harvestRepo(t), fakeGit{}, nil); got != nil {
		t.Errorf("a nil thread: %+v", got)
	}
}

// An answer with a secret's shape writes nothing, first or later, and the
// warning names the rule, never the text; the thread is left as it is.
func TestOnAnswerWritesNoSecret(t *testing.T) {
	r := harvestRepo(t)
	th := saved(t, r, "2026-10-20T09:01:00Z alex", "The key is AKIAIOSFODNN7EXAMPLE.")
	thread := readFile(t, r, "wip/threads/TH-0003-colour.md")
	got := OnAnswer(r, fakeGit{}, th)
	if got == nil || got.ID != "" || got.Action != "" || !strings.Contains(got.Warning, "rule digest.secret") || !strings.Contains(got.Warning, "TH-0003 at 2026-10-20T09:01:00Z") || strings.Contains(got.Warning, "AKIA") {
		t.Fatalf("OnAnswer = %+v", got)
	}
	if files := digestFiles(t, r.Root); len(files) != 0 {
		t.Errorf("wrote %v", files)
	}
	if after := readFile(t, r, "wip/threads/TH-0003-colour.md"); after != thread {
		t.Errorf("the thread changed:\n%s", after)
	}

	later := harvestRepo(t)
	first := []string{"2026-10-20T09:01:00Z alex", "Blue."}
	OnAnswer(later, fakeGit{}, saved(t, later, first...))
	before := readFile(t, later, "design/digest/D-0001-which-colour-blue.md")
	got = OnAnswer(later, fakeGit{}, saved(t, later, append(copied(first), "2026-10-20T09:02:00Z alex", "token=AKIAIOSFODNN7EXAMPLE")...))
	if got == nil || !strings.Contains(got.Warning, "rule digest.secret") || strings.Contains(got.Warning, "AKIA") {
		t.Errorf("a later secret: %+v", got)
	}
	if after := readFile(t, later, "design/digest/D-0001-which-colour-blue.md"); after != before {
		t.Errorf("a later secret changed the entry:\n%s", after)
	}
}

// Called from a story's worktree, OnAnswer writes the entry in the main
// checkout, where the thread is, and nothing in the worktree.
func TestOnAnswerWritesInTheMainCheckout(t *testing.T) {
	main := harvestRepo(t)
	wt := t.TempDir()
	write(t, wt, "system-flow.yaml", readFile(t, main, "system-flow.yaml"))
	r, err := workitem.Open(wt)
	if err != nil {
		t.Fatal(err)
	}
	r.MainRoot = main.Root
	got := OnAnswer(r, fakeGit{}, saved(t, r, "2026-10-20T09:01:00Z alex", "Blue."))
	if got == nil || got.Action != ActionWritten {
		t.Fatalf("OnAnswer = %+v", got)
	}
	if files := digestFiles(t, main.Root); !reflect.DeepEqual(files, []string{"D-0001-which-colour-blue.md"}) {
		t.Errorf("main checkout's digest %v", files)
	}
	if files := digestFiles(t, wt); len(files) != 0 {
		t.Errorf("the worktree's digest %v", files)
	}
}

// A write that fails, a file where the digest folder would be, only warns:
// what failed, on what answer, and what takes it in instead.
func TestOnAnswerOnlyWarnsWhenTheWriteFails(t *testing.T) {
	r := harvestRepo(t)
	write(t, r.Root, "design/digest", "not a folder\n")
	got := OnAnswer(r, fakeGit{}, saved(t, r, "2026-10-20T09:01:00Z alex", "Blue."))
	if got == nil || got.ID != "" || got.Action != "" ||
		!strings.Contains(got.Warning, "TH-0003 at 2026-10-20T09:01:00Z") || !strings.Contains(got.Warning, "design/digest") || !strings.Contains(got.Warning, "the next flai digest seed or curator run takes it in") {
		t.Errorf("OnAnswer = %+v", got)
	}
}
