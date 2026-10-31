package check

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
)

// S-0550: Fix rewrites summary.md that differs from its issue files to what
// flai issue summary writes, and the check then reports nothing for it.
func TestFixRewritesAStaleIssuesSummary(t *testing.T) {
	repo := issuesProject(t, [2]string{"I-0001", "open"}, [2]string{"I-0002", "open"})
	summary := filepath.Join(repo.Root, "design", "issues", "summary.md")
	if err := os.WriteFile(summary, []byte("# Open issues\n\n- [I-0001](I-0001-issue.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if got := rules(t, repo, "issues.summary"); len(got) != 1 {
		t.Fatalf("the stale summary should be reported once: %q", got)
	}
	fixed, err := Fix(repo, res, now)
	if err != nil {
		t.Fatal(err)
	}
	if want := []Fixed{{Rule: "issues.summary", Path: "design/issues/summary.md"}}; !slices.Equal(fixed, want) {
		t.Errorf("fixed %+v, want %+v", fixed, want)
	}
	data, err := os.ReadFile(summary)
	if err != nil {
		t.Fatal(err)
	}
	list, err := issues.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != issues.Summary(list, now) {
		t.Errorf("summary.md is not what flai issue summary writes:\n%s", data)
	}
	if got := rules(t, repo, "issues.summary"); len(got) != 0 {
		t.Errorf("after the fix the check still reports: %q", got)
	}
}

// S-0550: Fix runs no fixer for a rule with none, for a finding outside the
// story the run is scoped to, or for a rule the run did not report, and
// leaves the files as they were.
func TestFixRewritesOnlyWhatTheRunReportsInsideIt(t *testing.T) {
	repo := issuesProject(t, [2]string{"I-0001", "open"})
	summary := filepath.Join(repo.Root, "design", "issues", "summary.md")
	stale := "# Open issues\n"
	if err := os.WriteFile(summary, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, res := range map[string]*Result{
		"no findings":     {},
		"no fixer":        {Findings: []Finding{{Level: Error, Rule: "item.front-matter", Path: "wip/kanban/stories/S-001-x.md"}}},
		"outside a story": {Findings: []Finding{{Level: Warning, Rule: "issues.summary", Path: "design/issues/summary.md", Outside: true}}},
	} {
		fixed, err := Fix(repo, res, now.Add(time.Hour))
		if err != nil || len(fixed) != 0 {
			t.Errorf("%s: fixed %+v, %v; want nothing", name, fixed, err)
		}
		if data, err := os.ReadFile(summary); err != nil || string(data) != stale {
			t.Errorf("%s: summary.md was rewritten: %v\n%s", name, err, data)
		}
	}
}
