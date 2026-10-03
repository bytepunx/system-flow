package issues

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0198: each instance names the story it was recorded for, first under its
// heading, and the issue lists those stories once each, in order.
func TestAnInstanceNamesItsStory(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Flaky thing", Class: "defect", Note: "first", Story: "s-198", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(is.Body, "### 2026-09-16T10:00:00Z\nStory: S-0198.\nfirst\n") {
		t.Errorf("new names its story first:\n%s", is.Body)
	}
	if err := Bump(is, "S-0200", "", "second", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := Bump(is, "", "", "third, for no story", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := Bump(is, "S-0198", "", "fourth", t0.Add(3*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(is.Body, "### 2026-09-16T11:00:00Z\nStory: S-0200.\nsecond\n\n### 2026-09-16T12:00:00Z\nthird, for no story\n\n") {
		t.Errorf("bump names its story, or none:\n%s", is.Body)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(Stories(back), " "); got != "S-0198 S-0200" {
		t.Errorf("stories in order, once each: %q", got)
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "defect", Story: "T-0001", Now: t0}); err == nil {
		t.Error("a task is not a story")
	}
	if err := Bump(is, "nonsense", "", "", t0); err == nil || is.Count != 4 {
		t.Errorf("a bad story is refused before anything changes: %v, count %d", err, is.Count)
	}
}

// S-0198 with I-0011: an instance joined to one in the same second names its
// story only when that instance does not already.
func TestTheSameSecondJoinNamesAStoryOnce(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Flaky thing", Class: "defect", Note: "first", Story: "S-0198", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range []struct{ story, note string }{{"S-0198", "second"}, {"S-0200", "third"}, {"", "fourth"}} {
		if err := Bump(is, b.story, "", b.note, t0); err != nil {
			t.Fatal(err)
		}
	}
	want := "### 2026-09-16T10:00:00Z\nStory: S-0198.\nfirst\n\nsecond\n\nStory: S-0200.\nthird\n\nfourth\n\n## Remediation\n"
	if !strings.Contains(is.Body, want) {
		t.Errorf("one heading, each story named once:\n%s", is.Body)
	}
	if strings.Count(is.Body, "### ") != 1 {
		t.Errorf("one heading for the second:\n%s", is.Body)
	}
}

func TestRecordedByListsTheIssuesAStoryRecorded(t *testing.T) {
	r := repo(t)
	a, _ := New(r, NewOptions{Title: "A", Class: "defect", Story: "S-0198", Now: t0})
	b, _ := New(r, NewOptions{Title: "B", Class: "defect", Now: t0})
	c, _ := New(r, NewOptions{Title: "C", Class: "defect", Story: "S-0200", Now: t0})
	if err := Bump(c, "S-0198", "", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, is := range RecordedBy([]*Issue{a, b, c}, "s-198") {
		got = append(got, is.ID)
	}
	if strings.Join(got, " ") != "I-0001 I-0003" {
		t.Errorf("recorded by S-0198: %v", got)
	}
	if len(RecordedBy([]*Issue{a, b, c}, "S-0999")) != 0 {
		t.Error("a story that recorded nothing")
	}
}

func TestForStory(t *testing.T) {
	lint, err := mdlint.Parse([]byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD032: false\nMD041: false\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ class, nature string }{
		{"defect", "remediation"}, {"blocker", "remediation"}, {"efficiency", "improvement"}, {"impression", "improvement"},
	} {
		for _, fix := range []string{"", "Pin the linter in the Makefile.\n\n- and say so in design/tech"} {
			is := &Issue{ID: "I-0056", Title: "Lint: version mismatch", Class: c.class, Path: "/x/design/issues/I-0056-lint-version-mismatch.md",
				Body: "\n# I-0056 Lint: version mismatch\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n\n## Remediation\n" + fix + "\n"}
			d := ForStory(is)
			if d.Title != is.Title || d.Nature != c.nature {
				t.Errorf("%s: %+v", c.class, d)
			}
			if !strings.HasPrefix(d.Body, "## Goal\n\nThis story remediates [I-0056](../../../design/issues/I-0056-lint-version-mismatch.md), \"Lint: version mismatch\".") {
				t.Errorf("goal links the issue:\n%s", d.Body)
			}
			if fix != "" && !strings.Contains(d.Body, "solution:\n\n"+fix+"\n\n## Acceptance criteria\n") {
				t.Errorf("goal recommends the remediation:\n%s", d.Body)
			}
			if fix == "" && !strings.Contains(d.Body, "propose one from its instances") {
				t.Errorf("goal asks for a solution:\n%s", d.Body)
			}
			criteria := d.Body[strings.Index(d.Body, "## Acceptance criteria\n"):strings.Index(d.Body, "## Tasks")]
			lines := strings.Split(strings.TrimSpace(criteria), "\n")
			if len(lines) != 3 || !strings.HasPrefix(lines[1], "- [ ] The cause I-0056 describes no longer occurs") || lines[2] != "- [ ] I-0056 is closed with `flai issue close I-0056 --reason` saying what fixed it" {
				t.Errorf("criteria, closing the issue last:\n%s", criteria)
			}
			if !strings.HasSuffix(d.Body, "\n\n## Tasks\n\n## Notes\n") {
				t.Errorf("tasks and notes:\n%s", d.Body)
			}
			doc := "---\nid: S-0301\n---\n\n# S-0301 " + d.Title + "\n\n" + d.Body
			if f := lint.Lint(doc); len(f) > 0 {
				t.Errorf("lint: %v\n%s", f, doc)
			}
		}
	}
}

func TestLinks(t *testing.T) {
	is := &Issue{ID: "I-0056", Path: "/x/design/issues/I-0056-lint-version-mismatch.md"}
	for body, want := range map[string]bool{
		"Remediates I-0056.": true,
		"I-0056":             true,
		"see [it](../../../design/issues/I-0056-lint-version-mismatch.md)": true,
		"Remediates I-00561.": false,
		"Remediates XI-0056.": false,
		"Remediates I-005.":   false,
		"nothing here":        false,
	} {
		if got := Links(body, is); got != want {
			t.Errorf("Links(%q) = %v", body, got)
		}
	}
	legacy := &Issue{ID: "I-012"}
	if !Links("fixes I-0012", legacy) || !Links("fixes I-012", legacy) {
		t.Error("a three-digit issue is named in either form")
	}
}

func TestNoStory(t *testing.T) {
	open := func(id string) *Issue { return &Issue{ID: id, Status: "open"} }
	list := []*Issue{open("I-0001"), open("I-0002"), open("I-0003"), open("I-0004"), {ID: "I-0005", Status: "closed"}}
	items := []*workitem.Item{
		{ID: "S-0010", Type: workitem.Story, Status: workitem.Backlog, Body: "Remediates I-0001."},
		{ID: "S-0011", Type: workitem.Story, Status: workitem.Review, Body: "and I-0002"},
		{ID: "S-0012", Type: workitem.Story, Status: workitem.Done, Body: "Remediated I-0003."},
		{ID: "S-0013", Type: workitem.Story, Status: workitem.Cancelled, Body: "I-0003"},
		{ID: "T-0014", Type: workitem.Task, Status: workitem.InProgress, Body: "I-0004"},
		{ID: "E-0015", Type: workitem.Epic, Status: workitem.InProgress, Body: "I-0004"},
	}
	var got []string
	for _, is := range NoStory(list, items) {
		got = append(got, is.ID)
	}
	if strings.Join(got, " ") != "I-0003 I-0004" {
		t.Errorf("open issues no open story links: %v", got)
	}
}

// The story an issue's file names round-trips through the file.
func TestTheStoryLineSurvivesReading(t *testing.T) {
	r := repo(t)
	if _, err := New(r, NewOptions{Title: "A", Class: "defect", Story: "S-0198", Now: t0}); err != nil {
		t.Fatal(err)
	}
	list, err := List(r)
	if err != nil || len(list) != 1 || filepath.Base(list[0].Path) != "I-0001-a.md" {
		t.Fatalf("list: %v", err)
	}
	if got := Stories(list[0]); len(got) != 1 || got[0] != "S-0198" {
		t.Errorf("stories: %v", got)
	}
}
