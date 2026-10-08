package mdlint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The fixtures in testdata/cases are linted with the configuration beside
// them; expected.txt is what markdownlint-cli2 reports for them
// (scripts/mdlint-fixtures.sh). For every rule mdlint implements, it must
// report the same rule on the same lines.
func TestFixturesMatchMarkdownlint(t *testing.T) {
	dir := filepath.Join("testdata", "cases")
	c, err := Load(dir)
	if err != nil || c == nil {
		t.Fatalf("config: %v %v", c, err)
	}
	implemented := map[string]bool{}
	for _, id := range Rules() {
		implemented[id] = true
	}
	data, err := os.ReadFile(filepath.Join(dir, "expected.txt"))
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, l := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if f := strings.Fields(l); len(f) == 2 && implemented[f[1]] {
			want = append(want, l)
		}
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	var got []string
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, x := range c.Lint(string(content)) {
			got = append(got, fmt.Sprintf("%s:%d %s", filepath.Base(f), x.Line, x.Rule))
		}
	}
	diff(t, want, got)
}

func diff(t *testing.T, want, got []string) {
	t.Helper()
	count := map[string]int{}
	for _, w := range want {
		count[w]++
	}
	for _, g := range got {
		count[g]--
	}
	var missing, extra []string
	for k, n := range count {
		for ; n > 0; n-- {
			missing = append(missing, k)
		}
		for ; n < 0; n++ {
			extra = append(extra, k)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing) > 0 {
		t.Errorf("markdownlint reports, mdlint does not:\n  %s", strings.Join(missing, "\n  "))
	}
	if len(extra) > 0 {
		t.Errorf("mdlint reports, markdownlint does not:\n  %s", strings.Join(extra, "\n  "))
	}
}

func TestConfig(t *testing.T) {
	dir := t.TempDir()
	if c, err := Load(dir); c != nil || err != nil {
		t.Fatalf("no configuration means no lint: %v %v", c, err)
	}
	var none *Config
	if f := none.Lint("# a\n\n# a\n"); f != nil {
		t.Errorf("nil config lints nothing: %v", f)
	}
	jsonc := `{
  // comments are allowed
  "default": false, /* and these */
  "no-trailing-spaces": true,
  "headings": true,
  "MD024": false,
  "MD012": { "maximum": 2 }
}`
	if err := os.WriteFile(filepath.Join(dir, ".markdownlint.jsonc"), []byte(jsonc), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(dir)
	if err != nil || c.Path != ".markdownlint.jsonc" {
		t.Fatalf("load: %v %v", c, err)
	}
	for id, want := range map[string]bool{"MD009": true, "MD026": true, "MD001": true, "MD024": false, "MD012": true, "MD004": false} {
		if c.Enabled(id) != want {
			t.Errorf("%s enabled %v, want %v", id, c.Enabled(id), want)
		}
	}
	got := c.Lint("# Title.\n\n# Title.\n\n\nend \n")
	var rules []string
	for _, f := range got {
		rules = append(rules, fmt.Sprintf("%d %s", f.Line, f.Rule))
	}
	if strings.Join(rules, ",") != "1 MD026,3 MD025,3 MD026,6 MD009" {
		t.Errorf("findings %v", rules)
	}

	cli2 := t.TempDir()
	_ = os.WriteFile(filepath.Join(cli2, ".markdownlint-cli2.yaml"), []byte("config:\n  default: true\n  MD047: false\nglobs: [\"**/*.md\"]\n"), 0o644)
	_ = os.WriteFile(filepath.Join(cli2, ".markdownlint.yaml"), []byte("default: false\n"), 0o644)
	c, err = Load(cli2)
	if err != nil || c.Path != ".markdownlint-cli2.yaml" || c.Enabled("MD047") || !c.Enabled("MD009") {
		t.Errorf("cli2 file first, rules under config: %+v %v", c, err)
	}
}

func TestIntroducedAndGuard(t *testing.T) {
	c, err := Parse([]byte("MD022: false\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	before := "# T\n\n## Log\n\n### 2026-09-29T07:00:00Z\nOne.\n\n### 2026-09-29T07:00:00Z\nTwo.\n"
	if f := c.Lint(before); len(f) != 1 || f[0].Rule != "MD024" || f[0].Line != 8 {
		t.Fatalf("fixture: %v", f)
	}
	// the known duplicate moves down a line; a trailing space is new
	after := "# T\n\nIntro.\n\n## Log\n\n### 2026-09-29T07:00:00Z\nOne.\n\n### 2026-09-29T07:00:00Z\nTwo.\n\nThree. \n"
	f := c.Introduced(before, after)
	if len(f) != 1 || f[0].Rule != "MD009" || f[0].Line != 13 {
		t.Fatalf("introduced: %v", f)
	}
	if !strings.Contains(f[0].String(), "MD009/no-trailing-spaces Trailing spaces [Expected: 0 or 2; Actual: 1]") {
		t.Errorf("string: %s", f[0])
	}

	dir := t.TempDir()
	if err := Guard(dir, "x.md", before, after); err != nil {
		t.Errorf("no configuration, no guard: %v", err)
	}
	_ = os.WriteFile(filepath.Join(dir, ".markdownlint.yaml"), []byte("MD022: false\n"), 0o644)
	err = Guard(dir, "wip/x.md", before, after)
	var le *Error
	if !errors.As(err, &le) || len(le.Findings) != 1 || !strings.Contains(err.Error(), "wip/x.md") || !strings.Contains(err.Error(), "line 13: MD009") {
		t.Errorf("guard: %v", err)
	}
	if err := Guard(dir, "wip/x.md", before, before+"\nFour.\n"); err != nil {
		t.Errorf("a clean change passes: %v", err)
	}
}

// S-0231's body as written indents a top-level list by one space, which
// markdownlint reports as MD007 on each item.
func TestUnorderedListIndentOfS0231(t *testing.T) {
	c, err := Load(filepath.Join("testdata", "cases"))
	if err != nil || c == nil {
		t.Fatalf("config: %v %v", c, err)
	}
	content, err := os.ReadFile(filepath.Join("testdata", "cases", "license-story.md"))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range c.Lint(string(content)) {
		got = append(got, fmt.Sprintf("%d %s", f.Line, f))
	}
	want := []string{
		"23 MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]",
		"24 MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]",
		"25 MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]",
		"26 MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("findings:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// markdownlint's MD007 options, with what markdownlint-cli2 0.20.0 reports
// for each.
func TestUnorderedListIndentOptions(t *testing.T) {
	doc := "# T\n\n- top\n  - two\n    - four\n\n  - two again\n"
	for _, tc := range []struct{ config, want string }{
		{"default: true\n", ""},
		{"MD007:\n  indent: 4\n", "4 Expected: 4; Actual: 2|5 Expected: 8; Actual: 4|7 Expected: 4; Actual: 2"},
		{"MD007:\n  start_indented: true\n", "3 Expected: 2; Actual: 0|4 Expected: 4; Actual: 2|5 Expected: 6; Actual: 4|7 Expected: 4; Actual: 2"},
		{"MD007:\n  start_indented: true\n  start_indent: 1\n  indent: 1\n", "3 Expected: 1; Actual: 0|5 Expected: 3; Actual: 4"},
		{"MD007:\n  start_indent: 3\n", ""}, // start_indent counts only with start_indented
		{"ul-indent: false\n", ""},
	} {
		c, err := Parse([]byte(tc.config), false, false)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, f := range c.Lint(doc) {
			if f.Rule == "MD007" {
				got = append(got, fmt.Sprintf("%d %s", f.Line, f.Detail))
			}
		}
		if strings.Join(got, "|") != tc.want {
			t.Errorf("%q:\n got  %s\n want %s", tc.config, strings.Join(got, "|"), tc.want)
		}
	}
}

// I-0072: a task body's code span with a space before its closing backtick
// reached main, and markdownlint-cli2 0.20.0 stopped a close-out on it.
// mdlint reports it as markdownlint does, and lets CommonMark's one-space
// padding through.
func TestSpaceInCodeSpanOfI0072(t *testing.T) {
	c, err := Parse([]byte("default: true\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ doc, want string }{
		{"# T\n\n- Run `- Trigger: ` here.\n", "3 MD038/no-space-in-code Spaces inside code span elements [Context: \"`- Trigger: `\"]"},
		{"# T\n\nPadded ` a ` and `` `x` `` and `   ` pass.\n", ""},
		{"# T\n\nTwo `  a  ` here.\n", "3 MD038/no-space-in-code Spaces inside code span elements [Context: \"`  a  `\"]|" +
			"3 MD038/no-space-in-code Spaces inside code span elements [Context: \"`  a  `\"]"},
		{"# T\n\nAcross `a\nb ` lines.\n", "4 MD038/no-space-in-code Spaces inside code span elements [Context: \"`a b `\"]"},
		{"# T\n\nA long `code span that runs past thirty characters ` here.\n", "3 MD038/no-space-in-code Spaces inside code span elements [Context: \"... runs past thirty characters `\"]"},
		{"# T\n\n```text\n`x ` in a fence\n```\n", ""},
	} {
		var got []string
		for _, f := range c.Lint(tc.doc) {
			got = append(got, fmt.Sprintf("%d %s", f.Line, f))
		}
		if strings.Join(got, "|") != tc.want {
			t.Errorf("%q:\n got  %s\n want %s", tc.doc, strings.Join(got, "|"), tc.want)
		}
	}
}

// I-0056: TH-0067's entries carried a bare email address, which
// markdownlint-cli2 0.20.0 reports as MD034, and they reached main because
// mdlint knew only http and https literals. A thread entry that brings one
// is now refused; one in a code span or an autolink is not.
func TestBareEmailOfI0056(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".markdownlint.yaml"), []byte("default: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := "# TH-0067\n\n## Entries\n\n### 2026-10-02T15:35:00Z alex\n\nThe invoice goes to:\n"
	err := Guard(dir, "wip/threads/TH-0067.md", before, before+"alex@example.com, copied to <alex@example.com>.\n")
	var le *Error
	if !errors.As(err, &le) || len(le.Findings) != 1 || !strings.Contains(err.Error(), "line 8: MD034/no-bare-urls") {
		t.Errorf("guard: %v", err)
	}
	if err := Guard(dir, "wip/threads/TH-0067.md", before, before+"`alex@example.com`, copied to <alex@example.com>.\n"); err != nil {
		t.Errorf("a code span and an autolink are not bare: %v", err)
	}
}

// I-0110: mdlint skipped a bare www. literal, which markdownlint-cli2 0.20.0
// reports as MD034, so a thread entry carrying one passed flai and failed a
// close-out. A thread entry that brings one is now refused; one in a code
// span or link text is not.
func TestBareWwwOfI0110(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".markdownlint.yaml"), []byte("default: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := "# TH-0300\n\n## Entries\n\n### 2026-10-07T08:55:00Z agent\n\nThe docs are at:\n"
	err := Guard(dir, "wip/threads/TH-0300.md", before, before+"www.example.com, mirrored at [www.example.org](https://www.example.org).\n")
	var le *Error
	if !errors.As(err, &le) || len(le.Findings) != 1 || !strings.Contains(err.Error(), "line 8: MD034/no-bare-urls") {
		t.Errorf("guard: %v", err)
	}
	if err := Guard(dir, "wip/threads/TH-0300.md", before, before+"`www.example.com`, mirrored at [www.example.org](https://www.example.org).\n"); err != nil {
		t.Errorf("a code span and link text are not bare: %v", err)
	}
}

// I-0070: TH-0101's entry quoted step 3 of a list as "> 3.", which
// markdownlint-cli2 0.20.0 reports as MD029, and it reached main because
// mdlint did not parse blockquotes. A thread entry that brings one is now
// refused.
func TestQuotedOrderedListOfI0070(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".markdownlint.yaml"), []byte("default: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := "# TH-0101\n\n## Entries\n\n### 2026-10-04T04:00:00Z agent\n\nStep 3 now begins:\n"
	err := Guard(dir, "wip/threads/TH-0101.md", before, before+"\n> 3. Plan what your item's state calls for.\n")
	var le *Error
	if !errors.As(err, &le) || len(le.Findings) != 1 || !strings.Contains(err.Error(), "line 9: MD029/ol-prefix Ordered list item prefix [Expected: 1; Actual: 3; Style: 1/1/1]") {
		t.Errorf("guard: %v", err)
	}
	if err := Guard(dir, "wip/threads/TH-0101.md", before, before+"\n> 3\\. Plan what your item's state calls for.\n"); err != nil {
		t.Errorf("an escaped number is text: %v", err)
	}
}
