package verify

import (
	"encoding/json"
	"slices"
	"testing"
)

// S-0378: each markdownlint-cli2 line of a plain tier's output is a finding
// with its path, line, column, level, rule, and alias when it gives them,
// and the checkout to correct it in; other lines are none.
func TestFileFindingsGivesEachMarkdownlintLineAsData(t *testing.T) {
	t.Parallel()
	out := "lint-md: markdownlint-cli2 0.20.0\n" +
		"Summary: 3 error(s)\n" +
		"docs/a.md:3 error MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]\n" +
		"docs/a.md:7:1 error MD004/ul-style Unordered list style [Expected: asterisk; Actual: dash]\n" +
		"/wt/sub/b.md:1 MD041/first-line-heading/first-line-h1 First line in a file should be a top-level heading\n" +
		"c.md:2:5 warning MD013 Line length\n" +
		"exit status 1"
	res := Result{Tiers: []TierResult{
		{Name: "markdown", Dir: "sub", State: Failed, Findings: []Finding{{Name: "markdown", Message: out}}},
		{Name: "passed", State: Passed, Findings: []Finding{{Path: "p.go", Message: "not reported"}}},
		{Name: "go", State: Failed, Findings: []Finding{{Name: "errcheck", Path: "x.go", Line: 4, Message: "unchecked "}}},
	}}
	got := FileFindings("/wt", CheckoutStory, "S-0009", res)
	want := []FileFinding{
		{Path: "sub/docs/a.md", Line: 3, Level: "error", Rule: "MD012", Alias: "no-multiple-blanks", Message: "Multiple consecutive blank lines [Expected: 1; Actual: 2]"},
		{Path: "sub/docs/a.md", Line: 7, Column: 1, Level: "error", Rule: "MD004", Alias: "ul-style", Message: "Unordered list style [Expected: asterisk; Actual: dash]"},
		{Path: "sub/b.md", Line: 1, Rule: "MD041", Alias: "first-line-heading", Message: "First line in a file should be a top-level heading"},
		{Path: "sub/c.md", Line: 2, Column: 5, Level: "warning", Rule: "MD013", Message: "Line length"},
		{Tier: "go", Path: "x.go", Line: 4, Rule: "errcheck", Message: "unchecked"},
	}
	for i := range want {
		if want[i].Tier == "" {
			want[i].Tier = "markdown"
		}
		want[i].Checkout, want[i].Story = "story", "S-0009"
	}
	if !slices.Equal(got, want) {
		t.Errorf("got:\n%+v\nwant:\n%+v", got, want)
	}
}

// flai accept --json prints the findings and the dashboard reads them: their
// field names are its contract, and none is an empty list, not null.
func TestFileFindingsPrintAsTheRefusalPrintsThem(t *testing.T) {
	t.Parallel()
	none, err := json.Marshal(FileFindings("/repo", CheckoutMain, "", Result{}))
	if err != nil || string(none) != "[]" {
		t.Errorf("no findings print as %s (%v), want []", none, err)
	}
	f := FileFinding{Tier: "markdown", Checkout: CheckoutStory, Story: "S-0001", Path: "docs/bad.md", Line: 2, Column: 3,
		Level: "error", Rule: "MD009", Alias: "no-trailing-spaces", Message: "Trailing spaces"}
	b, err := json.Marshal(f)
	want := `{"tier":"markdown","checkout":"story","story":"S-0001","path":"docs/bad.md","line":2,"column":3,"level":"error","rule":"MD009","alias":"no-trailing-spaces","message":"Trailing spaces"}`
	if err != nil || string(b) != want {
		t.Errorf("got %s (%v), want %s", b, err, want)
	}
	b, _ = json.Marshal(FileFinding{Tier: "go", Checkout: CheckoutMain, Path: "x.go", Message: "m"})
	if want := `{"tier":"go","checkout":"main","path":"x.go","message":"m"}`; string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}

func TestCheckoutRelIsRelativeToTheCheckout(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ dir, p, want string }{
		{"", "docs/a.md", "docs/a.md"},
		{"sub", "a.md", "sub/a.md"},
		{"sub", "/wt/sub/b.md", "sub/b.md"},
		{"sub", "/elsewhere/c.md", "/elsewhere/c.md"},
		{"", "/wtx/d.md", "/wtx/d.md"},
	} {
		if got := checkoutRel("/wt", c.dir, c.p); got != c.want {
			t.Errorf("checkoutRel(/wt, %q, %q) = %q, want %q", c.dir, c.p, got, c.want)
		}
	}
}

func TestWithoutLeftLeavesOutEachPathLeftAndWhatLiesBelowIt(t *testing.T) {
	t.Parallel()
	left := []string{"design/a.md", "design/analysis/"}
	got := WithoutLeft([]string{"design/a.md", "design/analysis/r.md", "design/analysisx.md", "cli/b.go"}, left)
	if !slices.Equal(got, []string{"design/analysisx.md", "cli/b.go"}) {
		t.Errorf("WithoutLeft = %v", got)
	}
	if paths := []string{"a", "b"}; !slices.Equal(WithoutLeft(paths, nil), paths) {
		t.Error("with nothing left, WithoutLeft is paths")
	}
	if !IsLeft("design/analysis", left) || IsLeft("design", left) {
		t.Error("IsLeft names a folder left and what lies below it, not what lies above")
	}
}
