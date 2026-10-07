package issues

import (
	"strings"
	"testing"
	"time"
)

func TestFindOpenByTitle(t *testing.T) {
	list := []*Issue{
		{ID: "I-0003", Title: "a", Status: "open"},
		{ID: "I-0001", Title: "a", Status: "closed"},
		{ID: "I-0002", Title: "a", Status: "open"},
		{ID: "I-0004", Title: "b", Status: "open"},
	}
	if got := FindOpenByTitle(list, "a"); got == nil || got.ID != "I-0002" {
		t.Errorf("the lowest open issue titled a is I-0002, got %+v", got)
	}
	if got := FindOpenByTitle(list, "A"); got != nil {
		t.Errorf("a title matches exactly, got %+v", got)
	}
	if got := FindOpenByTitle([]*Issue{{ID: "I-0001", Title: "c", Status: "closed"}}, "c"); got != nil {
		t.Errorf("a closed issue is not found, got %+v", got)
	}
}

func TestHasInstance(t *testing.T) {
	is := &Issue{Body: "\n# I-0001 t\n\n## Description\nStory: S-0009.\nfound\n\n## Instances\n\n" +
		"### 2026-09-16T10:00:00Z\nStory: S-0001.\nfound:\n`a.md`: one\n\n" +
		"### 2026-09-16T11:00:00Z\nStory: S-0002.\nother\n\nStory: S-0003.\nthird\n\n" +
		"### 2026-09-16T12:00:00Z\nno story\n\n## Remediation\nStory: S-0004.\nfound:\n`a.md`: one\n"}
	for _, tc := range []struct {
		story, note string
		want        bool
	}{
		{"S-0001", "found:\n`a.md`: one", true},
		{"S-1", "  found:\n`a.md`: one\n", true}, // any padding, any surrounding space
		{"S-0001", "found:", false},              // a whole paragraph, not part of one
		{"S-0002", "found:\n`a.md`: one", false}, // another story's
		{"S-0003", "third", true},                // in a heading shared with another story
		{"S-0009", "found", false},               // Description is not an instance
		{"S-0004", "found:\n`a.md`: one", false}, // nor is Remediation
		{"", "no story", true},
		{"S-0001", "", false},
		{"not a story", "found:\n`a.md`: one", false},
	} {
		if got := HasInstance(is, tc.story, tc.note); got != tc.want {
			t.Errorf("HasInstance(%q, %q) = %v, want %v", tc.story, tc.note, got, tc.want)
		}
	}
}

func TestRecordOnce(t *testing.T) {
	r := repo(t)
	opt := NewOptions{Title: "flai check finds `x.y` outside the story at close-out", Class: "efficiency", Story: "S-0001", Note: "found:\n`a.md`: one", Now: t0}
	is, outcome, err := RecordOnce(r, opt)
	if err != nil || outcome != Opened || is.Count != 1 {
		t.Fatalf("first: %v %s %+v", err, outcome, is)
	}
	opt.Now = t0.Add(time.Hour)
	if is, outcome, err = RecordOnce(r, opt); err != nil || outcome != Already || is.Count != 1 {
		t.Fatalf("again: %v %s %+v", err, outcome, is)
	}
	opt.Note = "found:\n`a.md`: one\n`b.md`: two"
	if is, outcome, err = RecordOnce(r, opt); err != nil || outcome != Bumped || is.Count != 2 {
		t.Fatalf("other findings: %v %s %+v", err, outcome, is)
	}
	opt.Story = "S-0002"
	if is, outcome, err = RecordOnce(r, opt); err != nil || outcome != Bumped || is.Count != 3 {
		t.Fatalf("another story: %v %s %+v", err, outcome, is)
	}
	if err := Close(is, "fixed", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if is, outcome, err = RecordOnce(r, opt); err != nil || outcome != Opened || is.ID != "I-0002" || is.Count != 1 {
		t.Fatalf("after closing, a new issue: %v %s %+v", err, outcome, is)
	}
	reread, err := Read(is.Path)
	if err != nil || reread.Title != opt.Title {
		t.Fatalf("title round trip: %v %+v", err, reread)
	}
}

// S-0275: NewOrBump and RecordOnce report the issue file they wrote, whether
// they open the issue or bump it, and none when the occurrence was recorded
// already.
func TestNewOrBumpAndRecordOnceReportTheFileTheyWrote(t *testing.T) {
	r := repo(t)
	changed := func(is *Issue) string {
		t.Helper()
		got, err := is.Changed(r)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Join(got, " ")
	}
	const first, second = "design/issues/I-0001-review-waits-a-day.md", "design/issues/I-0002-plain.md"
	steps := []struct {
		name    string
		opt     NewOptions
		outcome Outcome
		want    string
	}{
		{"a report's finding opens", NewOptions{Title: "Review waits a day", Class: "efficiency", Report: report, Now: t0}, Opened, first},
		{"another report's bumps", NewOptions{Title: "Review waits a day", Class: "efficiency", Report: "design/analysis/2026-10-13-risk.md", Now: t0.Add(time.Hour)}, Bumped, first},
		{"the same report again", NewOptions{Title: "Review waits a day", Class: "efficiency", Report: report, Now: t0.Add(2 * time.Hour)}, Already, ""},
		{"no report opens", NewOptions{Title: "Plain", Class: "defect", Now: t0}, Opened, second},
	}
	for _, s := range steps {
		is, outcome, err := NewOrBump(r, s.opt)
		if err != nil || outcome != s.outcome {
			t.Fatalf("%s: %v %s", s.name, err, outcome)
		}
		if got := changed(is); got != s.want {
			t.Errorf("%s: changed %q, want %q", s.name, got, s.want)
		}
	}
	opt := NewOptions{Title: "Found at close-out", Class: "efficiency", Story: "S-0001", Note: "found", Now: t0}
	const third = "design/issues/I-0003-found-at-close-out.md"
	for _, want := range []struct {
		outcome Outcome
		changed string
	}{{Opened, third}, {Already, ""}} {
		is, outcome, err := RecordOnce(r, opt)
		if err != nil || outcome != want.outcome || changed(is) != want.changed {
			t.Fatalf("record once: %v %s %q, want %s %q", err, outcome, changed(is), want.outcome, want.changed)
		}
	}
	opt.Story = "S-0002"
	if is, outcome, err := RecordOnce(r, opt); err != nil || outcome != Bumped || changed(is) != third {
		t.Fatalf("record once for another story: %v %s %q", err, outcome, changed(is))
	}
}
