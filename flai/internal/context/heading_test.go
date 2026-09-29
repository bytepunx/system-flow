package context

import (
	"strings"
	"testing"
)

const headed = `---
title: The CLI
---

# The CLI

Intro.

## Commands

### flai prime --story

Prints the pack.

#### Budget

Fits 80 KB.

### flai doc

Documents.

## Dashboard

### flai prime --story

Not this one.

## Notes

End.
`

func TestSectionByTextPathAndSlugWithWhatIsBelowIt(t *testing.T) {
	for _, h := range []string{"Commands › flai prime --story", "commands > flai prime --story", "The CLI › Commands › flai prime --story", "#commands > flai-prime---story"} {
		p, err := Section("design/system/cli.md", headed, h)
		if err != nil {
			t.Fatalf("%s: %v", h, err)
		}
		if p.Heading != "Commands › flai prime --story" || p.Line != 11 || p.Text != "### flai prime --story\n\nPrints the pack.\n\n#### Budget\n\nFits 80 KB.\n\n" || p.Size != len(p.Text) || p.Path != "design/system/cli.md" {
			t.Errorf("%s: %+v", h, p)
		}
	}
	p, err := Section("design/system/cli.md", headed, "Commands")
	if err != nil || !strings.HasPrefix(p.Text, "## Commands\n") || !strings.Contains(p.Text, "Documents.") || strings.Contains(p.Text, "Dashboard") {
		t.Errorf("a level-2 section with those below it: %+v %v", p, err)
	}
	if p, err := Section("design/system/cli.md", headed, "## notes"); err != nil || p.Text != "## Notes\n\nEnd.\n" {
		t.Errorf("by lower-case text with its hashes: %+v %v", p, err)
	}
	if p, err := Section("design/system/cli.md", headed, "The CLI"); err != nil || !strings.HasPrefix(p.Text, "# The CLI\n") || !strings.HasSuffix(p.Text, "End.\n") {
		t.Errorf("the title is the whole body: %+v %v", p, err)
	}
}

func TestSectionRefusesAnAmbiguousOrUnknownHeading(t *testing.T) {
	_, err := Section("design/system/cli.md", headed, "flai prime --story")
	if err == nil || !strings.Contains(err.Error(), `"Commands › flai prime --story" (line 11)`) || !strings.Contains(err.Error(), `"Dashboard › flai prime --story" (line 25)`) {
		t.Errorf("ambiguous: %v", err)
	}
	_, err = Section("design/system/cli.md", headed, "Install")
	want := `no heading "Install" in design/system/cli.md; its headings are: "The CLI", "Commands", "Commands › flai prime --story", "Commands › flai prime --story › Budget", "Commands › flai doc", "Dashboard", "Dashboard › flai prime --story", "Notes"`
	if err == nil || err.Error() != want {
		t.Errorf("unknown:\n got %v\nwant %s", err, want)
	}
	for _, h := range []string{"", " # ", "›"} {
		if _, err := Section("design/system/cli.md", headed, h); err == nil {
			t.Errorf("%q must be refused", h)
		}
	}
}
