package digest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// entryText is an entry file with the fields given and a body below its
// heading.
func entryText(id, title, date, status, absorbedBy, cited, body string) string {
	if absorbedBy == "" {
		absorbedBy = `""`
	}
	return "---\nid: " + id + "\ntitle: " + title + "\nkind: decision\nsources: [S-0001]\ntopics: []\npaths: []\ndate: " + date +
		"\nstatus: " + status + "\nabsorbed_by: " + absorbedBy + "\nasked_again: []\ncited: " + cited + "\n---\n\n# " + id + " " + title + "\n\n" + body + "\n"
}

// citingText is a work item or thread file updated at updated whose body is
// text.
func citingText(id, updated, text string) string {
	return "---\nid: " + id + "\nupdated: " + updated + "\n---\n\n# " + id + "\n\n" + text + "\n"
}

// decayNow is the time the decay tests judge at; with a 30-day window an
// entry last cited or dated before 2026-09-24T12:00:00Z is compacted.
var decayNow = time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)

const month = 720 * time.Hour

// decayRepo is a project whose entries are cited, or not, in each way
// ADR-0199 counts:
//
//   - D-0001, dated 2026-08-01, cited by the epic E-0002, the story S-0001
//     (updated 2026-10-20), and the thread TH-0001, its cited holding only
//     TH-0001 so far;
//   - D-0002, dated 2026-08-01, cited by nothing that counts: the agents
//     folder's index, a design document, and another entry name it;
//   - D-0003, dated 2026-10-10 and cited by nothing;
//   - D-0004, absorbed, dated 2026-08-01 and cited by nothing;
//   - D-0005, dated 2026-08-01, cited as D-5 by S-0002's archived narrative,
//     updated 2026-10-01;
//   - D-0006, dated 2026-07-01, cited by the archived task T-0003, updated
//     on the day 2026-08-15;
//   - D-0007, dated 2026-08-01, cited by the thread TH-0001, updated
//     2026-09-30.
func decayRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	r := newRepo(t)
	root := r.Root
	write(t, root, "design/digest/README.md", "# Digest\n\nNot an entry: D-0003.\n")
	write(t, root, "design/digest/D-0001-kept.md", entryText("D-0001", "Kept by a story", "2026-08-01", "active", "", "[TH-0001]", "Why it was kept."))
	write(t, root, "design/digest/D-0002-uncited.md", entryText("D-0002", "Uncited and old", "2026-08-01", "active", "", "[]", "It names D-0001, which does not count."))
	write(t, root, "design/digest/D-0003-recent.md", entryText("D-0003", "Recent", "2026-10-10", "active", "", "[]", "Learned lately."))
	write(t, root, "design/digest/D-0004-absorbed.md", entryText("D-0004", "Absorbed", "2026-08-01", "absorbed", "design/system/x.md", "[]", "Recorded in x."))
	write(t, root, "design/digest/D-0005-narrative.md", entryText("D-0005", "Cited by a narrative", "2026-08-01", "active", "", "[]", "From a narrative."))
	write(t, root, "design/digest/D-0006-task.md", entryText("D-0006", "Cited by a task long ago", "2026-07-01", "active", "", "[]", "From a task."))
	write(t, root, "design/digest/D-0007-thread.md", entryText("D-0007", "Cited by a thread", "2026-08-01", "active", "", "[]", "From a thread."))
	write(t, root, "design/system/x.md", "# X\n\nIt records D-0002 and D-0004.\n\n## Rules\n\nOne rule.\n")

	write(t, root, "wip/kanban/epics/E-0002-e.md", citingText("E-0002", "2026-09-01T00:00:00Z", "Builds on D-0001."))
	write(t, root, "wip/kanban/stories/S-0001-s.md", citingText("S-0001", "2026-10-20T10:00:00Z", "Relies on D-0001 (D-0001 again) and on D-0099, which is no entry."))
	write(t, root, "wip/archive/kanban/tasks/T-0003-t.md", citingText("T-0003", "2026-08-15", "Used D-0006."))
	write(t, root, "wip/threads/TH-0001-q.md", citingText("TH-0001", "2026-09-30T00:00:00Z", "See D-0007 and D-0001."))
	write(t, root, "wip/archive/agents/S-0002.md", "---\nstream: S-0002\ntitle: t\nupdated: 2026-10-01T00:00:00Z\nagent: a\nsession: s\n---\n\n# S-0002 t\n\n## Decisions\n\n- As D-5 says.\n")
	write(t, root, "wip/agents/index.md", "# Narratives\n\nD-0002 is mentioned here.\n")
	return r
}

// digestFolder is every file in the digest folder and its text.
func digestFolder(t *testing.T, r *workitem.Repo) map[string]string {
	t.Helper()
	dir := filepath.Join(r.Root, "design", "digest")
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[f.Name()] = string(data)
	}
	return out
}

// ADR-0199 §4: a citation is an entry's ID in any padding in a work item,
// kanban or archive, a narrative, counted as its story, or a thread; its
// date is the file's updated time; the digest folder, design documents, and
// files that are not narratives do not count, and IDs of no entry are left
// out.
func TestCitationsFindsEachCiter(t *testing.T) {
	r := decayRepo(t)
	entries, err := List(r)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Citations(r, entries)
	if err != nil {
		t.Fatal(err)
	}
	at := func(s string) time.Time {
		v, err := updatedTime(s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	want := map[string]Citation{
		"D-0001": {Citers: []string{"E-0002", "S-0001", "TH-0001"}, Latest: at("2026-10-20T10:00:00Z")},
		"D-0005": {Citers: []string{"S-0002"}, Latest: at("2026-10-01T00:00:00Z")},
		"D-0006": {Citers: []string{"T-0003"}, Latest: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		"D-0007": {Citers: []string{"TH-0001"}, Latest: at("2026-09-30T00:00:00Z")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("citations:\n got %v\nwant %v", got, want)
	}
}

// A citing file whose updated time cannot be read is an error naming it.
func TestCitationsNamesACiterWithNoDate(t *testing.T) {
	r := decayRepo(t)
	write(t, r.Root, "wip/kanban/stories/S-0009-s.md", citingText("S-0009", "last week", "Cites D-0003."))
	entries, err := List(r)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Citations(r, entries)
	if err == nil || !strings.Contains(err.Error(), "wip/kanban/stories/S-0009-s.md cites D-0003") || !strings.Contains(err.Error(), `"last week"`) {
		t.Errorf("error: %v", err)
	}
}

// ADR-0199 §4: Cited lists every entry with its citers and its latest
// citation, writes cited in ID order where it changed, and a second run
// changes nothing.
func TestCitedWritesEachEntrysCiters(t *testing.T) {
	r := decayRepo(t)
	list, err := Cited(r, EnvFor(r))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range list {
		got = append(got, fmt.Sprintf("%s %s %v %q %v", c.ID, c.Status, c.Cited, c.Latest, c.Changed))
	}
	want := []string{
		`D-0001 active [E-0002 S-0001 TH-0001] "2026-10-20T10:00:00Z" true`,
		`D-0002 active [] "" false`,
		`D-0003 active [] "" false`,
		`D-0004 absorbed [] "" false`,
		`D-0005 active [S-0002] "2026-10-01T00:00:00Z" true`,
		`D-0006 active [T-0003] "2026-08-15T00:00:00Z" true`,
		`D-0007 active [TH-0001] "2026-09-30T00:00:00Z" true`,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("cited:\n got %s\nwant %s", strings.Join(got, "\n     "), strings.Join(want, "\n     "))
	}
	text := readFile(t, r, "design/digest/D-0001-kept.md")
	if !strings.Contains(text, "\ncited: [E-0002, S-0001, TH-0001]\n---\n") || !strings.HasSuffix(text, "# D-0001 Kept by a story\n\nWhy it was kept.\n") {
		t.Errorf("D-0001 is written with its citers and keeps its body:\n%s", text)
	}
	if e, err := Get(r, "D-5"); err != nil || !reflect.DeepEqual(e.Cited, []string{"S-0002"}) {
		t.Errorf("a narrative's citation is its story's: %v %v", e.Cited, err)
	}

	before := digestFolder(t, r)
	again, err := Cited(r, EnvFor(r))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range again {
		if c.Changed {
			t.Errorf("the second run changes %s", c.ID)
		}
	}
	if after := digestFolder(t, r); !reflect.DeepEqual(after, before) {
		t.Error("the second run writes nothing")
	}
}

// ADR-0199 §6: Compact runs Cited, then compacts each active entry last
// cited, or dated when it has no citation, before the window: an uncited old
// entry and one whose only citation is old become compacted, their bodies
// emptied below the heading; a cited entry, one whose own date is inside the
// window, and an absorbed one are kept. A second run compacts nothing more.
func TestCompactKeepsCitedAndCompactsUncited(t *testing.T) {
	r := decayRepo(t)
	kept := map[string]string{}
	for _, f := range []string{"D-0003-recent.md", "D-0004-absorbed.md"} {
		kept[f] = readFile(t, r, "design/digest/"+f)
	}
	res, err := Compact(r, CompactOptions{Window: month, Now: decayNow, Env: EnvFor(r)})
	if err != nil {
		t.Fatal(err)
	}
	want := []Compacted{
		{ID: "D-0002", Title: "Uncited and old", Path: "design/digest/D-0002-uncited.md", Latest: "", Date: "2026-08-01"},
		{ID: "D-0006", Title: "Cited by a task long ago", Path: "design/digest/D-0006-task.md", Latest: "2026-08-15T00:00:00Z", Date: "2026-07-01"},
	}
	if !reflect.DeepEqual(res.Compacted, want) || res.Off || res.DryRun || res.Window != "720h0m0s" || res.Before != "2026-09-24T12:00:00Z" || len(res.Cited) != 7 {
		t.Errorf("result %+v", res)
	}
	if got := readFile(t, r, "design/digest/D-0002-uncited.md"); !strings.Contains(got, "\nstatus: compacted\n") || !strings.HasSuffix(got, "---\n\n# D-0002 Uncited and old\n") {
		t.Errorf("D-0002 is compacted to its title:\n%s", got)
	}
	if e, err := Get(r, "D-0006"); err != nil || e.Status != "compacted" || e.Body != "" || !reflect.DeepEqual(e.Cited, []string{"T-0003"}) {
		t.Errorf("D-0006 is compacted, its citer written: %+v %v", e, err)
	}
	for _, id := range []string{"D-0001", "D-0005", "D-0007"} {
		if e, err := Get(r, id); err != nil || e.Status != "active" || e.Body == "" {
			t.Errorf("%s is cited within the window and kept whole: %+v %v", id, e, err)
		}
	}
	for f, text := range kept {
		if got := readFile(t, r, "design/digest/"+f); got != text {
			t.Errorf("%s is left as it was:\n%s", f, got)
		}
	}

	before := digestFolder(t, r)
	again, err := Compact(r, CompactOptions{Window: month, Now: decayNow, Env: EnvFor(r)})
	if err != nil {
		t.Fatal(err)
	}
	if len(again.Compacted) != 0 {
		t.Errorf("the second run compacts %+v", again.Compacted)
	}
	if after := digestFolder(t, r); !reflect.DeepEqual(after, before) {
		t.Error("the second run writes nothing")
	}
}

// ADR-0199 §6: a dry run lists what it would compact and what cited would
// change, and writes nothing.
func TestCompactDryRunWritesNothing(t *testing.T) {
	r := decayRepo(t)
	before := digestFolder(t, r)
	res, err := Compact(r, CompactOptions{Window: month, Now: decayNow, DryRun: true, Env: EnvFor(r)})
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, c := range res.Compacted {
		ids = append(ids, c.ID)
	}
	if !res.DryRun || strings.Join(ids, ",") != "D-0002,D-0006" || !res.Cited[0].Changed {
		t.Errorf("dry run result %+v", res)
	}
	if after := digestFolder(t, r); !reflect.DeepEqual(after, before) {
		t.Error("a dry run writes nothing")
	}
}

// ADR-0199 §6: a zero window turns compaction off: nothing is compacted or
// written, and the result says so.
func TestCompactOffCompactsNothing(t *testing.T) {
	r := decayRepo(t)
	before := digestFolder(t, r)
	res, err := Compact(r, CompactOptions{Window: 0, Now: decayNow, Env: EnvFor(r)})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Off || len(res.Compacted) != 0 || res.Before != "" {
		t.Errorf("result %+v", res)
	}
	if after := digestFolder(t, r); !reflect.DeepEqual(after, before) {
		t.Error("compaction off writes nothing")
	}
	if _, err := Compact(r, CompactOptions{Window: -time.Hour, Now: decayNow}); err == nil || !strings.Contains(err.Error(), "negative") {
		t.Errorf("a negative window: %v", err)
	}
}

// anchors is an AnchorFunc that knows the headings of design/system/x.md.
func anchors(path, _, heading string) (string, error) {
	if path == "design/system/x.md" && heading == "Rules" {
		return "rules", nil
	}
	return "", fmt.Errorf("no heading %q in %s; its headings are: X, Rules", heading, path)
}

// ADR-0199 §5: Absorb sets status absorbed and absorbed_by, keeping the
// body; with a heading absorbed_by is <path>#<anchor>; absorbing again names
// the newer document; a compacted entry may be absorbed.
func TestAbsorbMarksAnEntryAbsorbed(t *testing.T) {
	r := decayRepo(t)
	write(t, r.Root, "design/system/y.md", "# Y\n")
	res, err := Absorb(r, "D-3", AbsorbOptions{By: "./design/system/x.md", Anchor: anchors, Env: EnvFor(r)})
	if err != nil {
		t.Fatal(err)
	}
	if res.ID != "D-0003" || res.Entry.Status != "absorbed" || res.Entry.AbsorbedBy != "design/system/x.md" || res.Entry.Body != "Learned lately." {
		t.Errorf("result %+v", res)
	}
	text := readFile(t, r, "design/digest/D-0003-recent.md")
	if !strings.Contains(text, "\nstatus: absorbed\nabsorbed_by: design/system/x.md\n") || !strings.HasSuffix(text, "# D-0003 Recent\n\nLearned lately.\n") {
		t.Errorf("file:\n%s", text)
	}

	res, err = Absorb(r, "D-0003", AbsorbOptions{By: "design/system/x.md", Heading: "Rules", Anchor: anchors, Env: EnvFor(r)})
	if err != nil || res.Entry.AbsorbedBy != "design/system/x.md#rules" {
		t.Errorf("with a heading: %+v %v", res, err)
	}
	if e, err := Get(r, "D-0003"); err != nil || e.AbsorbedBy != "design/system/x.md#rules" {
		t.Errorf("read back: %+v %v", e, err)
	}
	if _, problems := Validate("D-0003-recent.md", readFile(t, r, "design/digest/D-0003-recent.md"), EnvFor(r)); len(problems) != 0 {
		t.Errorf("the absorbed entry checks clean: %+v", problems)
	}

	if res, err := Absorb(r, "D-0003", AbsorbOptions{By: "design/system/y.md", Env: EnvFor(r)}); err != nil || res.Entry.AbsorbedBy != "design/system/y.md" {
		t.Errorf("absorbed again names the newer document: %+v %v", res, err)
	}
	if _, err := Compact(r, CompactOptions{Window: month, Now: decayNow, Env: EnvFor(r)}); err != nil {
		t.Fatal(err)
	}
	if res, err := Absorb(r, "D-0002", AbsorbOptions{By: "design/system/x.md", Env: EnvFor(r)}); err != nil || res.Entry.Status != "absorbed" {
		t.Errorf("a compacted entry is absorbed: %+v %v", res, err)
	}
}

// ADR-0199 §5: a path outside the project, one that names no file or a
// folder, and a heading the document does not have are refused, naming
// absorbed_by, and nothing is written.
func TestAbsorbRefusesAndWritesNothing(t *testing.T) {
	r := decayRepo(t)
	before := digestFolder(t, r)
	for _, c := range []struct {
		name string
		opts AbsorbOptions
		says string
	}{
		{"no path", AbsorbOptions{}, "no document is given"},
		{"outside the project", AbsorbOptions{By: "../x.md"}, "not a path from the project root"},
		{"no file", AbsorbOptions{By: "design/system/gone.md"}, "design/system/gone.md names no file"},
		{"a folder", AbsorbOptions{By: "design/system"}, "names no file"},
		{"no such heading", AbsorbOptions{By: "design/system/x.md", Heading: "Gone", Anchor: anchors}, `no heading "Gone" in design/system/x.md`},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.opts.Env = EnvFor(r)
			_, err := Absorb(r, "D-0001", c.opts)
			var refused *RefusedError
			if !errors.As(err, &refused) || len(refused.Problems) != 1 {
				t.Fatalf("want a refusal, got %v", err)
			}
			p := refused.Problems[0]
			if refused.Path != "design/digest/D-0001-kept.md" || p.Field != "absorbed_by" || p.Rule != RuleEntry || p.Line != 10 || !strings.Contains(p.Message, c.says) {
				t.Errorf("refusal %s %+v, want saying %q", refused.Path, p, c.says)
			}
		})
	}
	if after := digestFolder(t, r); !reflect.DeepEqual(after, before) {
		t.Error("a refused absorption writes nothing")
	}
	if _, err := Absorb(r, "D-0042", AbsorbOptions{By: "design/system/x.md"}); err == nil || !strings.Contains(err.Error(), "D-0042 is not in design/digest") {
		t.Errorf("an unknown entry: %v", err)
	}
}

// IDs sort by prefix, then by number, however wide.
func TestCompareIDs(t *testing.T) {
	ids := []string{"TH-0002", "S-10000", "T-0001", "S-0009", "E-0003"}
	slices.SortFunc(ids, compareIDs)
	if want := []string{"E-0003", "S-0009", "S-10000", "T-0001", "TH-0002"}; !slices.Equal(ids, want) {
		t.Errorf("sorted %v, want %v", ids, want)
	}
}
