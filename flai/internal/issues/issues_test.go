package issues

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 16, 10, 0, 0, 0, time.UTC)

func repo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "design"), 0o755)
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestLifecycle(t *testing.T) {
	r := repo(t)
	if id := NextID(r, nil); id != "I-0001" {
		t.Fatalf("first id %s", id)
	}
	is, err := New(r, NewOptions{Title: "Lint: version mismatch", Class: "efficiency", Cost: "5m", Note: "found in S-004", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if is.ID != "I-0001" || is.Count != 1 || is.Cost != "5m" || !strings.HasSuffix(is.Path, "I-0001-lint-version-mismatch.md") {
		t.Fatalf("new: %+v", is)
	}
	raw, _ := os.ReadFile(is.Path)
	if !strings.Contains(string(raw), `title: "Lint: version mismatch"`) || !strings.Contains(string(raw), "### 2026-09-16T10:00:00Z\nfound in S-004") {
		t.Errorf("file:\n%s", raw)
	}
	back, err := Read(is.Path)
	if err != nil || back.Marshal() != string(raw) {
		t.Fatalf("round trip: %v", err)
	}
	if err := Bump(back, "", "15m", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if back.Count != 2 || back.Cost != "10m" || back.LastReported != "2026-09-16T11:00:00Z" || back.FirstReported != "2026-09-16T10:00:00Z" {
		t.Errorf("bump: %+v", back)
	}
	if !strings.Contains(back.Body, "### 2026-09-16T11:00:00Z\nagain\n\n## Remediation") {
		t.Errorf("instance not inserted before remediation:\n%s", back.Body)
	}
	if err := Bump(back, "", "", "no cost this time", t0.Add(2*time.Hour)); err != nil || back.Count != 3 || back.Cost != "10m" {
		t.Errorf("bump without cost: %v %+v", err, back)
	}
	second, _ := New(r, NewOptions{Title: "Two", Class: "defect", Now: t0})
	if second.ID != "I-0002" || second.Cost != "" {
		t.Errorf("second: %+v", second)
	}
	list, err := WriteSummary(r, t0)
	if err != nil || len(list) != 2 {
		t.Fatal(err)
	}
	sum, _ := os.ReadFile(filepath.Join(Dir(r), SummaryFile))
	s := string(sum)
	if !strings.Contains(s, "| [I-0001](I-0001-lint-version-mismatch.md) | efficiency | Lint: version mismatch | 3 | 10m | 30m |") || !strings.Contains(s, "| [I-0002](I-0002-two.md) | defect | Two | 1 | - | - |") {
		t.Errorf("summary:\n%s", s)
	}
	if strings.Index(s, "[I-0001]") > strings.Index(s, "[I-0002]") {
		t.Error("most expensive first")
	}
	if err := Close(back, "fixed by S-010", t0.Add(3*time.Hour)); err != nil || back.Status != "closed" || !strings.Contains(back.Body, "Closed 2026-09-16T13:00:00Z: fixed by S-010") {
		t.Errorf("close: %v %+v", err, back)
	}
	if err := Bump(back, "", "", "", t0); err == nil {
		t.Error("bump on closed should fail")
	}
	list, _ = WriteSummary(r, t0)
	if table := SummaryTable(list); strings.Contains(table, "I-0001") || !strings.Contains(table, "I-0002") {
		t.Errorf("closed issue in table:\n%s", table)
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "bug", Now: t0}); err == nil {
		t.Error("bad class should fail")
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "defect", Cost: "soon", Now: t0}); err == nil {
		t.Error("bad cost should fail")
	}
}

func TestValidate(t *testing.T) {
	good := &Issue{ID: "I-0001", Title: "t", Class: "defect", Status: "open", Count: 1, Cost: "5m", FirstReported: "2026-09-16T10:00:00Z", LastReported: "2026-09-16T10:00:00Z", Updated: "2026-09-16T10:00:00Z"}
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := *good
	bad.Class, bad.Status, bad.Count, bad.Cost = "meh", "maybe", 0, "lots"
	err := bad.Validate()
	for _, want := range []string{"class", "status", "count", "cost"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("expected %s in %v", want, err)
		}
	}
	for in, want := range map[string]string{"1h0m0s": "1h", "30m0s": "30m", "2h30m0s": "2h30m", "10s": "10s", "30s": "30s", "1m30s": "1m30s", "0s": "0s", "": ""} {
		if got := normalise(in); got != want {
			t.Errorf("normalise(%q) = %q, want %q", in, got, want)
		}
	}
}

// S-0181: an issue a newer flai wrote is listed, and a bump keeps the field
// this flai does not know.
func TestAnIssueWithAnUnknownFieldIsListedAndKept(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(is.Path)
	doc := strings.Replace(string(data), "\n---\n", "\nowner: sam\n---\n", 1)
	if err := os.WriteFile(is.Path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := List(r)
	if err != nil || len(list) != 1 || len(list[0].Unknown) != 1 || list[0].Unknown[0].Name != "owner" {
		t.Fatalf("listing: %v %+v", err, list)
	}
	if err := Bump(list[0], "", "", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(is.Path)
	if !strings.Contains(string(after), "count: 2\n") || !strings.Contains(string(after), "owner: sam\n---\n") {
		t.Errorf("a bump dropped the unknown field:\n%s", after)
	}
}

// I-0231: a long title made a file name that Claude Code's Edit could not
// write once its temporary suffix was added; the slug is cut at 80 bytes.
func TestNewCutsALongTitlesFileName(t *testing.T) {
	r := repo(t)
	title := "Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes"
	is, err := New(r, NewOptions{Title: title, Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(is.Path)
	if max := len(is.ID) + 1 + 80 + len(".md"); len(base) > max || !strings.HasPrefix(base, is.ID+"-claude-code-s-edit-") || !strings.HasSuffix(base, ".md") {
		t.Errorf("the file name is %d bytes, want at most %d on a word: %s", len(base), max, base)
	}
	if back, err := Read(is.Path); err != nil || back.Title != title {
		t.Errorf("the issue keeps its whole title: %v", err)
	}
}

// A project in git with no story branches numbers past what main holds and
// what its own checkout holds uncommitted.
// SummaryPath follows layout.design, and names the file WriteSummary writes.
func TestSummaryPath(t *testing.T) {
	r := repo(t)
	if got := SummaryPath(r); got != "design/issues/summary.md" {
		t.Errorf("summary path %q", got)
	}
	if got, want := filepath.Join(r.Root, filepath.FromSlash(SummaryPath(r))), filepath.Join(Dir(r), SummaryFile); got != want {
		t.Errorf("summary path %s is not where WriteSummary writes, %s", got, want)
	}
	r.Manifest.Layout["design"] = "./docs/design"
	if got := SummaryPath(r); got != "docs/design/issues/summary.md" {
		t.Errorf("summary path under another layout %q", got)
	}
}

// S-0275: New, Bump, BumpWith, and Close report the issue file they wrote,
// relative to the root under layout.design; an issue only read reports none,
// and one written outside the checkout is refused.
func TestChangedIsTheIssueFileEachWriteWrote(t *testing.T) {
	r := repo(t)
	r.Manifest.Layout["design"] = "notes/design"
	changed := func(is *Issue) string {
		t.Helper()
		got, err := is.Changed(r)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Join(got, " ")
	}
	is, err := New(r, NewOptions{Title: "Lint: version mismatch", Class: "efficiency", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	const file = "notes/design/issues/I-0001-lint-version-mismatch.md"
	if got := changed(is); got != file {
		t.Errorf("new: %q, want %q", got, file)
	}
	back, err := Get(r, "I-1")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := back.Changed(r); err != nil || got == nil || len(got) != 0 {
		t.Errorf("an issue only read wrote nothing: %q %v", got, err)
	}
	if err := Bump(back, "S-0275", "5m", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if got := changed(back); got != file {
		t.Errorf("bump: %q, want %q", got, file)
	}
	if err := BumpWith(r, back, BumpOptions{Report: "notes/design/analysis/2026-10-06-risk.md", Now: t0.Add(2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if got := changed(back); got != file {
		t.Errorf("bump with a report: %q, want %q", got, file)
	}
	if err := Close(back, "fixed by S-0275", t0.Add(3*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if got := changed(back); got != file {
		t.Errorf("close: %q, want %q", got, file)
	}
	if err := Bump(back, "", "", "", t0); err == nil || changed(back) != file {
		t.Errorf("a refused bump wrote, or forgot what the close wrote: %v %q", err, changed(back))
	}
	if _, err := WriteSummary(r, t0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(r.Root, filepath.FromSlash(SummaryPath(r)))); err != nil || SummaryPath(r) != "notes/design/issues/summary.md" {
		t.Errorf("WriteSummary wrote no file at SummaryPath %s: %v", SummaryPath(r), err)
	}
	other := repo(t)
	if _, err := back.Changed(other); err == nil || !strings.Contains(err.Error(), "outside the checkout") {
		t.Errorf("an issue written in another checkout: %v", err)
	}
}

func TestNextIDWithoutStoryBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	r, run := repo(t), execx.System{}
	if err := os.MkdirAll(Dir(r), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(Dir(r), "I-0001-committed.md"), []byte("one\n"), 0o644)
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "one"}} {
		if out, err := run.Run(r.Root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	_ = os.WriteFile(filepath.Join(Dir(r), "I-0002-uncommitted.md"), []byte("two\n"), 0o644)
	if id := NextID(r, run); id != "I-0003" {
		t.Errorf("next id %s, want I-0003", id)
	}
}

// past moves every file of the issues folder an hour back, so that what the
// kept store parses of them is trusted (S-0602).
func past(t *testing.T, r *workitem.Repo, ago time.Duration) {
	t.Helper()
	files, _ := filepath.Glob(filepath.Join(Dir(r), "*.md"))
	old := time.Now().Add(-ago)
	for _, f := range files {
		if err := os.Chtimes(f, old, old); err != nil {
			t.Fatal(err)
		}
	}
}

// ids is the IDs and counts of list, in its order.
func ids(list []*Issue) string {
	var b strings.Builder
	for _, is := range list {
		fmt.Fprintf(&b, "%s:%d:%s ", is.ID, is.Count, is.Status)
	}
	return b.String()
}

// S-0602: List keeps what it parsed and parses again only the issues that
// changed, were added, or are gone, a bump written in place at the same size
// included.
func TestListParsesAgainOnlyTheIssuesThatChanged(t *testing.T) {
	r := repo(t)
	one, _ := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	two, _ := New(r, NewOptions{Title: "Two", Class: "defect", Now: t0})
	if _, err := WriteSummary(r, t0); err != nil {
		t.Fatal(err)
	}
	past(t, r, time.Hour)
	if list, err := List(r); err != nil || ids(list) != "I-0001:1:open I-0002:1:open " {
		t.Fatalf("first list: %v %s", err, ids(list))
	}
	before := kept.Parses()
	if list, err := List(r); err != nil || ids(list) != "I-0001:1:open I-0002:1:open " {
		t.Fatalf("second list: %v %s", err, ids(list))
	}
	if n := kept.Parses() - before; n != 0 {
		t.Errorf("a second list with nothing changed parsed %d files", n)
	}

	// A bump through Save.
	got, err := Get(r, "I-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := Bump(got, "", "", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	past(t, r, 50*time.Minute)
	if list, _ := List(r); ids(list) != "I-0001:2:open I-0002:1:open " {
		t.Errorf("after a bump: %s", ids(list))
	}

	// A change in place at the same size, with a new modification time.
	data, _ := os.ReadFile(two.Path)
	same := strings.Replace(string(data), "count: 1\n", "count: 7\n", 1)
	if len(same) != len(data) {
		t.Fatal("the edit changed the size")
	}
	if err := os.WriteFile(two.Path, []byte(same), 0o644); err != nil {
		t.Fatal(err)
	}
	newer := time.Now().Add(-40 * time.Minute)
	if err := os.Chtimes(two.Path, newer, newer); err != nil {
		t.Fatal(err)
	}
	if list, _ := List(r); ids(list) != "I-0001:2:open I-0002:7:open " {
		t.Errorf("after an edit at the same size: %s", ids(list))
	}

	// Added and removed.
	three, _ := New(r, NewOptions{Title: "Three", Class: "defect", Now: t0})
	if list, _ := List(r); ids(list) != "I-0001:2:open I-0002:7:open I-0003:1:open " {
		t.Errorf("after an add: %s", ids(list))
	}
	if err := os.Remove(one.Path); err != nil {
		t.Fatal(err)
	}
	if list, _ := List(r); ids(list) != "I-0002:7:open I-0003:1:open " {
		t.Errorf("after a removal: %s", ids(list))
	}
	if _, err := Get(r, "I-0001"); err == nil {
		t.Error("Get found a removed issue")
	}

	// A caller's change does not reach the kept parse.
	list, _ := List(r)
	list[0].Count = 99
	list[0].Unknown = append(list[0].Unknown, workitem.Field{Name: "x"})
	if again, _ := List(r); again[0].Count != 7 || len(again[0].Unknown) != 0 {
		t.Errorf("a caller's change reached the kept issue: %+v", again[0])
	}
	if got, _ := Get(r, three.ID); got.Title != "Three" {
		t.Errorf("Get: %+v", got)
	}
}

// S-0602: two lists at once get the same issues.
func TestTwoListsAtOnceGetTheSameIssues(t *testing.T) {
	r := repo(t)
	for _, title := range []string{"One", "Two", "Three"} {
		if _, err := New(r, NewOptions{Title: title, Class: "defect", Now: t0}); err != nil {
			t.Fatal(err)
		}
	}
	past(t, r, time.Hour)
	var wg sync.WaitGroup
	got := make([]string, 8)
	for i := range got {
		wg.Go(func() {
			list, err := List(r)
			if err != nil {
				t.Error(err)
			}
			got[i] = ids(list)
		})
	}
	wg.Wait()
	for i, g := range got {
		if g != "I-0001:1:open I-0002:1:open I-0003:1:open " {
			t.Errorf("list %d: %s", i, g)
		}
	}
}
