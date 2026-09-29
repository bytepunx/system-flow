package search

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestTermsDropStopwordsAndRepeats(t *testing.T) {
	if got := Terms("The heron and the Tide, a heron: x"); !slices.Equal(got, []string{"heron", "tide"}) {
		t.Errorf("got %v", got)
	}
}

func TestSectionsRankEachSectionOnItsOwnAndSkipEmptyOnes(t *testing.T) {
	secs := []Section{
		{Path: "design/system/cli.md", Title: "The CLI", Heading: "", Level: 0, Line: 6, Text: "Intro about commands.\n"},
		{Path: "design/system/cli.md", Title: "The CLI", Heading: "Commands", Level: 2, Line: 8, Text: "## Commands\n"},
		{Path: "design/system/cli.md", Title: "The CLI", Heading: "Commands › flai prime", Level: 3, Line: 10, Text: "### flai prime\n\nPrints the context pack.\n\nThe pack fits a budget.\n"},
		{Path: "docs/users/flai.md", Title: "flai", Heading: "Heron", Level: 2, Line: 3, Text: "## Heron\n\nA bird that wades.\n"},
	}
	x := IndexSections(secs)
	if got := x.Rank("the pack"); !slices.Equal(got, []int{2}) {
		t.Errorf("rank %v", got)
	}
	if got := x.Rank("commands"); !slices.Contains(got, 0) || slices.Contains(got, 1) {
		t.Errorf("an empty section was ranked, or the text before the first heading was not: %v", got)
	}
	hits := x.Search("prime pack", 0)
	if len(hits) != 1 {
		t.Fatalf("hits %+v", hits)
	}
	h := hits[0]
	if h.Path != "design/system/cli.md" || h.Heading != "Commands › flai prime" || h.Line != 10 || h.Title != "The CLI" ||
		h.Lines != "Prints the context pack.\nThe pack fits a budget." || h.Size != len(secs[2].Text) || h.Score <= 0 {
		t.Errorf("hit %+v", h)
	}
	if got := x.Search("the a", 0); len(got) != 0 {
		t.Errorf("stopwords matched %+v", got)
	}
}

func TestSectionSearchReturnsAtMostTwenty(t *testing.T) {
	var secs []Section
	for i := range 30 {
		secs = append(secs, Section{Path: fmt.Sprintf("design/system/d%02d.md", i), Heading: "Part", Level: 2, Text: "## Part\n\nheron\n"})
	}
	x := IndexSections(secs)
	if got := len(x.Search("heron", 0)); got != MaxSectionHits {
		t.Errorf("default %d", got)
	}
	if got := len(x.Search("heron", 100)); got != MaxSectionHits {
		t.Errorf("over the cap %d", got)
	}
	if got := len(x.Search("heron", 3)); got != 3 {
		t.Errorf("limit 3 gave %d", got)
	}
}

func TestFirstLinesSkipBlanksAndCutLongLines(t *testing.T) {
	got := FirstLines("\n  one  \n\n"+strings.Repeat("x", 10)+"\nthree\nfour\n", 3, 5)
	if got != "one\nxxxxx…\nthree" {
		t.Errorf("got %q", got)
	}
}
