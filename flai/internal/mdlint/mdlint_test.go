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

// Every markdown file in the monorepo passes markdownlint-cli2 (make
// lint-md, in CI), so mdlint must report nothing on any of them: a finding
// here is one markdownlint would not make.
func TestRepositoryLintsClean(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	c, err := Load(root)
	if err != nil || c == nil {
		t.Fatalf("config: %v %v", c, err)
	}
	skip := map[string]bool{"node_modules": true, "testdata": true, "bin": true, ".flai-cache": true, ".svelte-kit": true, ".git": true}
	var got []string
	n := 0
	err = filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if e.IsDir() {
			if skip[e.Name()] || rel == filepath.Join("flaiover", "build") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		n++
		for _, f := range c.Lint(string(data)) {
			got = append(got, fmt.Sprintf("%s:%d %s", rel, f.Line, f))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n < 100 {
		t.Fatalf("linted %d files; expected the monorepo's", n)
	}
	diff(t, nil, got)
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
