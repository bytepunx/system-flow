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

// adrDoc is an ADR file numbered n with title and status.
func adrDoc(n, title, status string) string {
	return "---\nid: ADR-" + n + "\ntitle: " + title + "\nstatus: " + status + "\ndate: 2026-10-31\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-" + n + " " + title + "\n\n## Decision\n\nWe decide.\n"
}

// writeADRs writes files, by name, in the project's design/adrs.
func writeADRs(t *testing.T, root string, files map[string]string) {
	t.Helper()
	dir := filepath.Join(root, "design", "adrs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// S-0550, ADR-0228 §1 and §3: an ADR index behind its folder is one
// adr.index warning on its first line that differs, naming the rows; Fix
// rewrites it to what flai adr index writes, and the check then reports
// nothing for it.
func TestFixRewritesAnADRIndexBehindItsFolder(t *testing.T) {
	repo := issuesProject(t)
	if _, err := issues.WriteSummary(repo, now); err != nil {
		t.Fatal(err)
	}
	intro := "# Architecture Decision Records\n\nThe table is generated.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n"
	writeADRs(t, repo.Root, map[string]string{
		"0000-template.md": adrDoc("0000", "Template", "template"),
		"0001-first.md":    adrDoc("0001", "First", "accepted"),
		"0002-second.md":   adrDoc("0002", "Second", "proposed"),
		"README.md": intro + "| [0001](0001-first.md) | First | proposed |\n" +
			"| [0007](0007-gone.md) | Gone | accepted |\n\nText below stays.\n",
	})
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var found []Finding
	for _, f := range res.Findings {
		if f.Rule == "adr.index" {
			found = append(found, f)
		}
	}
	want := "the ADR index differs from what its ADR files render: row 0001 differs, row 0002 is missing and row 0007 links no ADR file; run flai check --fix or flai adr index, which rewrite it"
	if len(found) != 1 || found[0].Level != Warning || found[0].Path != "design/adrs/README.md" || found[0].Line != 7 || found[0].Message != want {
		t.Fatalf("the index behind its folder should be one warning on line 7:\n%+v\nwant message %q", found, want)
	}
	fixed, err := Fix(repo, res, now)
	if err != nil {
		t.Fatal(err)
	}
	if want := []Fixed{{Rule: "adr.index", Path: "design/adrs/README.md"}}; !slices.Equal(fixed, want) {
		t.Errorf("fixed %+v, want %+v", fixed, want)
	}
	data, err := os.ReadFile(filepath.Join(repo.Root, "design", "adrs", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := intro + "| [0001](0001-first.md) | First | accepted |\n| [0002](0002-second.md) | Second | proposed |\n\nText below stays.\n"; string(data) != want {
		t.Errorf("the index is not what flai adr index writes:\n%s\nwant:\n%s", data, want)
	}
	if got := rules(t, repo, "adr.index"); len(got) != 0 {
		t.Errorf("after the fix the check still reports: %q", got)
	}
}

// The template's index, an empty table beside 0000-template.md, is current,
// and so is an index with no table and no ADRs to list; a project with no ADR
// index has nothing to report.
func TestADRIndexOfTheTemplateOrOfNoIndexIsNotReported(t *testing.T) {
	repo := issuesProject(t)
	if got := rules(t, repo, "adr.index"); len(got) != 0 {
		t.Errorf("no ADR folder: %q", got)
	}
	writeADRs(t, repo.Root, map[string]string{"0000-template.md": adrDoc("0000", "Template", "template"), "0001-first.md": adrDoc("0001", "First", "accepted")})
	if got := rules(t, repo, "adr.index"); len(got) != 0 {
		t.Errorf("no index: %q", got)
	}
	_ = os.Remove(filepath.Join(repo.Root, "design", "adrs", "0001-first.md"))
	writeADRs(t, repo.Root, map[string]string{"README.md": "# Architecture Decision Records\n\nText.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n"})
	if got := rules(t, repo, "adr.index"); len(got) != 0 {
		t.Errorf("the template's index: %q", got)
	}
	writeADRs(t, repo.Root, map[string]string{"README.md": "# adrs\n"})
	if got := rules(t, repo, "adr.index"); len(got) != 0 {
		t.Errorf("an index with no table and no ADRs to list: %q", got)
	}
}
