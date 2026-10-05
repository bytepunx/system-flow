package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scopeProject copies the good fixture, whose S-004 is in progress with its
// narrative, its task T-003, and the thread TH-0001 anchored on it, and gives
// the archived S-002 a branch that was never merged: story.unaccepted.
func scopeProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	refs := filepath.Join(root, ".git", "refs", "heads", "story")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refs, "S-002"), []byte("0123456789abcdef0123456789abcdef01234567\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// edit rewrites a fixture file, failing the test when old is not in it.
func edit(t *testing.T, root, rel, old, new string) {
	t.Helper()
	p := filepath.Join(root, rel)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), old) {
		t.Fatalf("%s has no %q", rel, old)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(data), old, new, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runScoped(t *testing.T, root string, changed []string) (*Result, *Result) {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	whole, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	scoped, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := ScopeToStory(scoped, repo, "S-004", changed); err != nil {
		t.Fatal(err)
	}
	return whole, scoped
}

// I-0057: a story's close-out stopped at flai check --strict on
// story.unaccepted for another, archived story whose branch was never
// merged. Scoped to the story in progress, the run passes; unscoped, it
// fails as before.
func TestScopedCheckPassesOverAnotherStorysUnmergedBranch(t *testing.T) {
	whole, scoped := runScoped(t, scopeProject(t), nil)
	var unaccepted *Finding
	for i, f := range scoped.Findings {
		if f.Rule == "story.unaccepted" {
			unaccepted = &scoped.Findings[i]
		}
	}
	if unaccepted == nil || !strings.Contains(unaccepted.Path, "S-002") || !strings.Contains(unaccepted.Message, "story/S-002 was never merged") {
		t.Fatalf("the fixture should give story.unaccepted on S-002: %+v", scoped.Findings)
	}
	if whole.OK(true) || !whole.OK(false) {
		t.Errorf("unscoped, story.unaccepted must fail --strict only: %+v", whole)
	}
	if !unaccepted.Outside || unaccepted.Level != Warning {
		t.Errorf("story.unaccepted on S-002 should be a warning outside S-004: %+v", unaccepted)
	}
	if !scoped.OK(true) {
		t.Errorf("scoped to S-004, the run should pass --strict: %+v", scoped)
	}
	// The other finding is E-001's advisory epic.lags-stories: outside too,
	// and counted once, as outside, not again as advisory.
	if scoped.Outside != 2 || scoped.Advisory != 0 || scoped.Warnings != whole.Warnings || whole.Advisory != 1 {
		t.Errorf("outside %d advisory %d warnings %d, want 2, 0, %d", scoped.Outside, scoped.Advisory, scoped.Warnings, whole.Warnings)
	}
	for _, f := range whole.Findings {
		if f.Outside {
			t.Errorf("unscoped, nothing is outside: %+v", f)
		}
	}
}

// S-0249: a finding on what the story owns still fails --strict when the
// run is scoped to it; a wip.overlap, even on its own file, and an error on
// another story's file do not.
func TestScopedCheckFailsOnlyOnTheStorysOwnFindings(t *testing.T) {
	const overview = "design/system/overview.md"
	for _, tc := range []struct {
		name    string
		rule    string
		path    string // the finding's path ends with this
		changed []string
		fails   bool
		setup   func(t *testing.T, root string)
	}{
		{name: "the story's item file", rule: "item.heading", path: "S-004-four.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/stories/S-004-four.md", "\n## Notes", "\n")
		}},
		{name: "a task's file", rule: "item.heading", path: "T-003-t3.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/tasks/T-003-t3.md", "\n## Notes", "\n")
		}},
		{name: "the narrative", rule: "narrative.section", path: "S-004.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/agents/S-004.md", "\n## Open questions\n", "\n")
		}},
		{name: "a thread anchored on the story", rule: "threads.entries", path: "TH-0001-is-four-really-done.md", fails: true, setup: func(t *testing.T, root string) {
			p := filepath.Join(root, "wip/threads/TH-0001-is-four-really-done.md")
			data, _ := os.ReadFile(p)
			head, _, _ := strings.Cut(string(data), "\n### ")
			if err := os.WriteFile(p, []byte(head+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "a path in the story's diff", rule: "doc.title", path: overview, changed: []string{overview}, fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "a path in a folder the story adds", rule: "doc.title", path: overview, changed: []string{"design/system/"}, fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "a path outside the story's diff", rule: "doc.title", path: overview, changed: []string{"flai/main.go"}, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "wip.overlap with another story in progress", rule: "wip.overlap", path: "S-004-four.md", setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/stories/S-004-four.md", "tags: []\n", "tags: []\ntouches: [flai]\n")
			addStory(t, root, "S-006", "E-001", "touches: [flai/cmd]\n")
		}},
		{name: "an error on another story's file", rule: "item.parent-missing", path: "S-006-six.md", setup: func(t *testing.T, root string) {
			addStory(t, root, "S-006", "E-009", "")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeProject(t)
			tc.setup(t, root)
			whole, scoped := runScoped(t, root, tc.changed)
			var got *Finding
			for i, f := range scoped.Findings {
				if f.Rule == tc.rule && strings.HasSuffix(filepath.ToSlash(f.Path), tc.path) {
					got = &scoped.Findings[i]
				}
			}
			if got == nil {
				t.Fatalf("no %s on %s: %+v", tc.rule, tc.path, scoped.Findings)
			}
			if whole.OK(true) {
				t.Errorf("unscoped, the run should fail --strict: %+v", whole)
			}
			if got.Outside == tc.fails {
				t.Errorf("outside %v, want %v: %+v", got.Outside, !tc.fails, got)
			}
			if scoped.OK(true) == tc.fails {
				t.Errorf("scoped OK(strict) %v, want %v: %+v", scoped.OK(true), !tc.fails, scoped.Findings)
			}
		})
	}
}

// addStory writes an in-progress story under parent, with front matter
// extra, and its narrative.
func addStory(t *testing.T, root, id, parent, extra string) {
	t.Helper()
	story := "---\nid: " + id + "\ntype: story\nnature: feature\ntitle: Six\nstatus: in-progress\nparent: " + parent +
		"\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\n" + extra +
		"---\n\n# " + id + " Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/stories", id+"-six.md"), []byte(story), 0o644); err != nil {
		t.Fatal(err)
	}
	narrative, err := os.ReadFile(filepath.Join(root, "wip/agents/S-004.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/agents", id+".md"), []byte(strings.ReplaceAll(string(narrative), "S-004", id)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestScopeToStoryNamesAStory(t *testing.T) {
	repo, err := workitem.Open(scopeProject(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"S-099", "T-003", "nonsense"} {
		if err := ScopeToStory(&Result{}, repo, id, nil); err == nil || !strings.Contains(err.Error(), "scope the check to") {
			t.Errorf("%s: want an error naming the scope, got %v", id, err)
		}
	}
}
