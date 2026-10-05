package issues

import (
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
