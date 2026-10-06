package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0223: an analysis.schedule that does not parse and an analysis.agent that
// is not valid are errors on the analysis key, each saying what to write; a
// valid block, and no block, are not reported.
func TestAnalysisSettingsAreChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"analysis:\n  schedule: daily\n",
		"analysis:\n  schedule: \"0 6 * * 1\"\n  agent:\n    model: claude-sonnet-5\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.analysis"); len(got) != 0 {
			t.Errorf("a valid analysis block %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "analysis:\n  schedule: weekly\n  agent:\n    harness: Claude Code\n")
	got := findings(t, repo, "manifest.analysis")
	if len(got) != 2 {
		t.Fatalf("want two errors, the schedule and the agent: %+v", got)
	}
	for i, want := range []string{`analysis.schedule "weekly" has 1 fields, not five (minute, hour, day of month, month, day of week); write one such as`, `analysis.agent harness "Claude Code" is not a name such as claude-code`} {
		f := got[i]
		if f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
			t.Errorf("finding %d: %+v, want an error on line 8 saying %q", i, f, want)
		}
	}
	repo = planningProject(t, "analysis:\n  schedule: \"0 0 31 2 *\"\n")
	if got := findings(t, repo, "manifest.analysis"); len(got) != 1 || !strings.Contains(got[0].Message, "never comes round") {
		t.Errorf("want one error, a schedule that never comes round: %+v", got)
	}
}

// S-0223: flai check validates every report in design/analysis and warns of
// one that the folder's README.md does not list, and of a row that links to a
// report that is not there; a good, listed report is not reported, and the
// folder is not checked as other design documents are.
func TestAnalysisReportsAreChecked(t *testing.T) {
	repo := planningProject(t, "")
	dir := filepath.Join(repo.Root, "design", "analysis")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("README.md", "# analysis\n\n| Report | Focus | Window | Status |\n|--------|-------|--------|--------|\n| [Bottlenecks](2026-10-06-bottlenecks.md) | bottlenecks | 2026-09-01 to 2026-09-30 | active |\n| [Risk](2026-10-06-risk.md) | risk | 2026-09-01 to 2026-09-30 | active |\n| [Gone](2026-10-01-all.md) | all | 2026-09-01 to 2026-09-30 | active |\n")
	write("2026-10-06-bottlenecks.md", "---\ntitle: Bottlenecks\nupdated: 2026-10-06T08:00:00Z\nstatus: active\nfocus: bottlenecks\nfrom: 2026-09-01\nto: 2026-09-30\n---\n\n# Bottlenecks\n")
	write("2026-10-06-risk.md", "---\ntitle: Risk\nupdated: 2026-10-06T08:00:00Z\nstatus: active\nfocus: intent\nfrom: 2026-09-01\n---\n\n# Risk\n")
	write("2026-10-06-intent.md", "---\ntitle: Intent\nupdated: 2026-10-06\nstatus: draft\nfocus: intent\nfrom: 2026-09-01\nto: 2026-09-30\n---\n\n# Intent\n")

	got := findings(t, repo, "analysis.report")
	if len(got) != 2 {
		t.Fatalf("want two errors on the risk report, its window's missing end and its focus: %+v", got)
	}
	for i, want := range []struct {
		line int
		says string
	}{{1, "front matter needs to"}, {5, "focus is intent but the file is named for risk"}} {
		f := got[i]
		if f.Level != Error || f.Path != "design/analysis/2026-10-06-risk.md" || f.Line != want.line || f.Message != want.says {
			t.Errorf("finding %d: %+v, want an error on line %d saying %q", i, f, want.line, want.says)
		}
	}
	got = findings(t, repo, "analysis.index")
	if len(got) != 2 {
		t.Fatalf("want two warnings, the unlisted intent report and the row to a missing one: %+v", got)
	}
	unlisted, gone := got[0], got[1]
	if unlisted.Path != "design/analysis/2026-10-06-intent.md" {
		unlisted, gone = gone, unlisted
	}
	if unlisted.Level != Warning || unlisted.Path != "design/analysis/2026-10-06-intent.md" || !strings.Contains(unlisted.Message, "design/analysis/README.md does not list it") {
		t.Errorf("unlisted report: %+v", unlisted)
	}
	if gone.Level != Warning || gone.Path != "design/analysis/README.md" || gone.Line != 7 || !strings.Contains(gone.Message, "2026-10-01-all.md, which is not there") {
		t.Errorf("row to a missing report: %+v", gone)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if strings.HasPrefix(f.Path, "design/analysis/") && !strings.HasPrefix(f.Rule, "analysis.") {
			t.Errorf("the folder is checked by its own rules only: %+v", f)
		}
	}
}
