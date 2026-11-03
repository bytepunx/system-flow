package mdlint_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Every markdown file in the monorepo passes markdownlint-cli2 (make
// lint-md, in CI), so mdlint must report nothing on any of them: a finding
// here is one markdownlint would not make.
//
// The integration tier's own: it lints the whole monorepo, so it skips under
// -short, as the close-out's go-test tier runs (S-0548). Under flai verify, in
// a close-out or a merge queue batch, the smoke tier's scripts/lint-md.sh
// lints the same files with markdownlint-cli2 itself, so this test runs there
// only when the stories change what its verdict depends on besides the
// markdown: mdlint's code, the markdownlint configuration, or lint-md.sh,
// which pins the markdownlint-cli2 version mdlint agrees with. A false
// positive the markdown alone brings is still caught in wip/, which flai check
// --strict lints with mdlint, and elsewhere unscoped: make integration, CI,
// and flai test --all.
func TestRepositoryLintsClean(t *testing.T) {
	if testing.Short() {
		t.Skip("the integration tier's own: lints the whole monorepo")
	}
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	// In a close-out (scripts/close-out.sh exports the story), the branch
	// holds wip/ as main last committed it: a finding in a wip file the story
	// does not change is main's, logged and not failed (I-0117). In a merge
	// queue batch, the same holds for a file none of its stories changes
	// (I-0174).
	var changed []string
	stories := scopeStories()
	story := strings.Join(stories, ",")
	if len(stories) > 0 {
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		changed = []string{}
		for _, id := range stories {
			c, err := storygit.StoryChanges(execx.System{}, repo, id)
			if err != nil {
				t.Fatal(err)
			}
			changed = append(changed, c...)
		}
		if !changesTheVerdict(changed) {
			t.Skipf("under flai verify for %s, which changes neither mdlint nor its configuration: the smoke tier's scripts/lint-md.sh lints the monorepo", story)
		}
	}
	failing, outside, n, err := lintTree(root, changed)
	if err != nil {
		t.Fatal(err)
	}
	if n < 100 {
		t.Fatalf("linted %d files; expected the monorepo's", n)
	}
	for _, f := range outside {
		t.Logf("outside %s: %s", story, f)
	}
	for _, f := range failing {
		t.Errorf("mdlint reports, markdownlint does not: %s", f)
	}
}

// scopeStories is the stories the run is scoped to: the close-out's story
// in CLOSE_OUT_STORY, else the merge queue batch's, comma-separated in
// MERGE_QUEUE_STORIES, else none.
func scopeStories() []string {
	if story := os.Getenv("CLOSE_OUT_STORY"); story != "" {
		return []string{story}
	}
	var stories []string
	for _, id := range strings.Split(os.Getenv(verify.MergeQueueStoriesEnv), ",") {
		if id = strings.TrimSpace(id); id != "" {
			stories = append(stories, id)
		}
	}
	return stories
}

// changesTheVerdict reports whether the changed paths hold what
// TestRepositoryLintsClean's verdict depends on besides the markdown: mdlint's
// code, a markdownlint configuration file, or scripts/lint-md.sh.
func changesTheVerdict(changed []string) bool {
	for _, p := range changed {
		if strings.HasPrefix(p, "flai/internal/mdlint/") || p == "scripts/lint-md.sh" || slices.Contains(mdlint.ConfigFiles, p) {
			return true
		}
	}
	return false
}

func TestChangesTheVerdict(t *testing.T) {
	for _, tc := range []struct {
		changed []string
		want    bool
	}{
		{nil, false},
		{[]string{"design/x.md", "flai/internal/check/check.go", "scripts/check.sh"}, false},
		{[]string{"design/x.md", "flai/internal/mdlint/rules.go"}, true},
		{[]string{"flai/internal/mdlint/"}, true},
		{[]string{".markdownlint.yaml"}, true},
		{[]string{".markdownlint-cli2.jsonc"}, true},
		{[]string{"scripts/lint-md.sh"}, true},
	} {
		if got := changesTheVerdict(tc.changed); got != tc.want {
			t.Errorf("changesTheVerdict(%v) = %v, want %v", tc.changed, got, tc.want)
		}
	}
}

// I-0117: a finding in main's wip fails only the close-out of a story that changes the file.
func TestMainsWipOutsideTheStoryOfI0117(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(".markdownlint.yaml", "default: false\nMD034: true\n")
	write("wip/agents/x.md", "# X\n\nwww.example.com\n")
	write("wip/agents/clean.md", "# Clean\n")
	for _, tc := range []struct {
		name             string
		changed          []string
		failing, outside []string
	}{
		{"not in a close-out", nil, []string{"wip/agents/x.md:3"}, nil},
		{"story leaves it alone", []string{"design/y.md"}, nil, []string{"wip/agents/x.md:3"}},
		{"story changes it", []string{"wip/agents/x.md"}, []string{"wip/agents/x.md:3"}, nil},
		{"story adds its folder", []string{"wip/agents/"}, []string{"wip/agents/x.md:3"}, nil},
	} {
		failing, outside, n, err := lintTree(root, tc.changed)
		if err != nil {
			t.Fatal(err)
		}
		if n != 2 {
			t.Errorf("%s: linted %d files, want 2", tc.name, n)
		}
		if !sameFiles(failing, tc.failing) || !sameFiles(outside, tc.outside) {
			t.Errorf("%s: failing %v outside %v, want %v and %v", tc.name, failing, outside, tc.failing, tc.outside)
		}
	}

	write("design/x.md", "# X\n\nwww.example.com\n")
	for _, changed := range [][]string{nil, {}, {"design/x.md"}} {
		failing, _, _, err := lintTree(root, changed)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.ContainsFunc(failing, func(f string) bool { return strings.HasPrefix(f, "design/x.md:3 MD034") }) {
			t.Errorf("changed %v: design/x.md must fail, got %v", changed, failing)
		}
	}
}

// sameFiles reports whether the findings are on exactly the file:line places.
func sameFiles(findings, want []string) bool {
	if len(findings) != len(want) {
		return false
	}
	for i, f := range findings {
		if !strings.HasPrefix(f, want[i]+" ") {
			return false
		}
	}
	return true
}

// lintTree splits the findings under root into failing and outside the story that changed them (nil: no close-out).
func lintTree(root string, changed []string) (failing, outside []string, n int, err error) {
	c, err := mdlint.Load(root)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("load the markdownlint configuration in %s: %w", root, err)
	}
	if c == nil {
		return nil, nil, 0, fmt.Errorf("%s has no markdownlint configuration: add .markdownlint.yaml", root)
	}
	inStory := func(rel string) bool {
		for _, p := range changed {
			if rel == p || (strings.HasSuffix(p, "/") && strings.HasPrefix(rel, p)) {
				return true
			}
		}
		return false
	}
	skip := map[string]bool{"node_modules": true, "testdata": true, "bin": true, ".flai-cache": true, ".svelte-kit": true, ".git": true}
	err = filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if e.IsDir() {
			if skip[e.Name()] || rel == "flaiover/build" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		// A link whose target is missing, as the strategic agents' committed
		// links into .flai-cache/agents are until flai role definition --write
		// runs (ADR-0236), is left out, as markdownlint-cli2 leaves it out.
		if e.Type()&os.ModeSymlink != 0 {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return nil
			}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		n++
		mains := changed != nil && strings.HasPrefix(rel, "wip/") && !inStory(rel)
		for _, f := range c.Lint(string(data)) {
			line := fmt.Sprintf("%s:%d %s", rel, f.Line, f)
			if mains {
				outside = append(outside, line)
			} else {
				failing = append(failing, line)
			}
		}
		return nil
	})
	return failing, outside, n, err
}
