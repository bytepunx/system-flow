package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// activityRepo makes a project with the repository's markdownlint settings,
// a planner document of two activities, an orchestrator document of one, an
// orphan narrative, and a stray file under wip/agents. It returns the root
// and the planner document's text.
func activityRepo(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range []struct {
		kind    string
		at      time.Time
		seconds int64
		cost    float64
	}{
		{workitem.ActivityPlanner, now.Add(-2 * time.Hour), 600, 0.1234},
		{workitem.ActivityPlanner, now.Add(-time.Hour), 123, 0.2001},
		{workitem.ActivityOrchestrator, now.Add(-time.Hour), 60, 0.05},
	} {
		if _, err := repo.AppendActivity(a.kind, workitem.ActivityEntry{At: a.at, Summary: "Planned", Items: []string{"S-0001"}, Seconds: a.seconds, Cost: a.cost}); err != nil {
			t.Fatal(err)
		}
	}
	narrative := "---\nstream: S-9999\ntitle: Gone\nagent: t\nupdated: 2026-09-01T12:00:00Z\n---\n\n# S-9999 Gone\n\n## Context\n\n## Current state\n\n## Next steps\n\n## Decisions\n\n## Open questions\n\n## Log\n"
	_ = os.WriteFile(filepath.Join(root, "wip/agents/S-9999.md"), []byte(narrative), 0o644)
	_ = os.WriteFile(filepath.Join(root, "wip/agents/foo.md"), []byte("# Foo\n"), 0o644)
	data, err := os.ReadFile(repo.ActivityPath(workitem.ActivityPlanner))
	if err != nil {
		t.Fatal(err)
	}
	return root, string(data)
}

func runCheck(t *testing.T, root string) *Result {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func activityFindings(res *Result) []Finding {
	var out []Finding
	for _, f := range res.Findings {
		if strings.HasPrefix(f.Rule, "activity.") {
			out = append(out, f)
		}
	}
	return out
}

// S-0206: valid activity documents pass with no finding, are not taken for
// narratives, and are linted as wip; the orphan narrative is still reported
// and the stray file is not taken for an activity document.
func TestActivityDocumentsAreValidated(t *testing.T) {
	root, _ := activityRepo(t)
	res := runCheck(t, root)
	if got := activityFindings(res); len(got) != 0 {
		t.Errorf("valid documents have findings: %+v", got)
	}
	orphan := false
	for _, f := range res.Findings {
		base := filepath.Base(f.Path)
		if workitem.IsActivityDocument(base) {
			t.Errorf("finding on an activity document: %+v", f)
		}
		if base == "foo.md" {
			t.Errorf("finding on the stray file: %+v", f)
		}
		if f.Rule == "narrative.orphan" && base == "S-9999.md" {
			orphan = true
		}
	}
	if !orphan {
		t.Errorf("the orphan narrative is not reported: %+v", res.Findings)
	}
}

func TestActivityDocumentFaults(t *testing.T) {
	for _, tc := range []struct {
		name, old, new string
		level, rule    string
		line           int
		says           string
	}{
		{"unknown key", "kind: planner\n", "kind: planner\nowner: x\n", Error, "activity.document", 1, "front matter must hold only"},
		{"kind mismatch", "kind: planner\n", "kind: analyzer\n", Error, "activity.document", 1, `kind "analyzer" must be the file's name, "planner"`},
		{"bad log line", "- Seconds: 600\n", "- Seconds: ten\n", Error, "activity.document", 1, `seconds "ten" must be a whole number`},
		{"negative cost", "accrued_cost: 0.3235\n", "accrued_cost: -0.3235\n", Error, "activity.front-matter", 3, "accrued_cost is -0.3235 but must be zero or more"},
		{"negative seconds", "accrued_seconds: 723\n", "accrued_seconds: -1\n", Error, "activity.front-matter", 4, "accrued_seconds is -1 but must be zero or more"},
		{"negative count", "tasks_completed: 2\n", "tasks_completed: -2\n", Error, "activity.front-matter", 5, "tasks_completed is -2 but must be zero or more"},
		{"count", "tasks_completed: 2\n", "tasks_completed: 3\n", Warning, "activity.totals", 5, "tasks_completed is 3 but the log has 2 entries"},
		{"seconds", "accrued_seconds: 723\n", "accrued_seconds: 724\n", Warning, "activity.totals", 4, "accrued_seconds is 724 but the log's seconds sum to 723"},
		{"cost", "accrued_cost: 0.3235\n", "accrued_cost: 0.3236\n", Warning, "activity.totals", 3, "accrued_cost is 0.3236 but the log's costs sum to 0.3235"},
		{"last run", "last_run: 2026-09-01T11:00:00Z\n", "last_run: 2026-09-01T10:00:00Z\n", Warning, "activity.totals", 6, `last_run is "2026-09-01T10:00:00Z" but the newest log entry ended at "2026-09-01T11:00:00Z"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, doc := activityRepo(t)
			if !strings.Contains(doc, tc.old) {
				t.Fatalf("the planner document has no %q:\n%s", tc.old, doc)
			}
			path := filepath.Join(root, "wip/agents/planner.md")
			_ = os.WriteFile(path, []byte(strings.Replace(doc, tc.old, tc.new, 1)), 0o644)
			res := runCheck(t, root)
			got := activityFindings(res)
			if len(got) != 1 {
				t.Fatalf("want one finding: %+v", got)
			}
			f := got[0]
			if f.Level != tc.level || f.Rule != tc.rule || f.Path != filepath.Join("wip", "agents", "planner.md") || f.Line != tc.line {
				t.Errorf("finding: %+v", f)
			}
			if !strings.Contains(f.Message, tc.says) || !strings.Contains(f.Message, "restore it from git (git checkout -- wip/agents/planner.md)") {
				t.Errorf("message %q does not say %q and how to restore it", f.Message, tc.says)
			}
			if strings.Contains(f.Message, root) {
				t.Errorf("message repeats the absolute path: %q", f.Message)
			}
			// the fixture's other findings are not advisory either, so a
			// totals warning counted as advisory would show here
			if res.Advisory != 0 {
				t.Errorf("a totals warning is not advisory: it fails --strict like others: %+v", res)
			}
		})
	}
}

// The costs are rounded to four decimals, so a sum that drifts by less than
// half the last place is not a disagreement.
func TestActivityCostTolerance(t *testing.T) {
	root, doc := activityRepo(t)
	path := filepath.Join(root, "wip/agents/planner.md")
	_ = os.WriteFile(path, []byte(strings.Replace(doc, "accrued_cost: 0.3235\n", "accrued_cost: 0.32354\n", 1)), 0o644)
	if got := activityFindings(runCheck(t, root)); len(got) != 0 {
		t.Errorf("a drift within rounding is reported: %+v", got)
	}
}
