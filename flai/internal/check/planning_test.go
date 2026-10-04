package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planningProject is an empty project whose system-flow.yaml ends with extra.
func planningProject(t *testing.T, extra string) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"+extra), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// findings are a run's findings under rule.
func findings(t *testing.T, repo *workitem.Repo, rule string) []Finding {
	t.Helper()
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var out []Finding
	for _, f := range res.Findings {
		if f.Rule == rule {
			out = append(out, f)
		}
	}
	return out
}

// S-0199: a draft story past backlog is warned about on its draft line; a
// draft in backlog, and a story that is not a draft, are not.
func TestADraftPastBacklogIsWarned(t *testing.T) {
	repo := planningProject(t, "")
	set := func(title, status string, draft bool) *workitem.Item {
		t.Helper()
		s := mustItem(t, repo, workitem.Story, title, "")
		s.Status, s.Draft = status, draft
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	ready := set("Draft ready", workitem.Ready, true)
	set("Draft backlog", workitem.Backlog, true)
	set("Final ready", workitem.Ready, false)
	got := findings(t, repo, "story.draft")
	if len(got) != 1 {
		t.Fatalf("want one story.draft, on %s: %+v", ready.ID, got)
	}
	f := got[0]
	want := ready.ID + " is a draft in ready; finalize it (flai edit " + ready.ID + " --no-draft) or move it back to backlog"
	if f.Level != Warning || !strings.Contains(f.Path, ready.ID) || f.Line != keyLine(ready.Path, "draft") || f.Line == 1 || f.Message != want {
		t.Errorf("got %+v, want a warning on %s's draft line saying %q", f, ready.ID, want)
	}
}

// S-0199: each thing wrong with the manifest's planning block is an error on
// its key; a valid block is not reported.
func TestPlanningSettingsAreChecked(t *testing.T) {
	repo := planningProject(t, "planning:\n  currency: EUR\n  hour_rate: 80\n  cycle: 336h\n")
	if got := findings(t, repo, "manifest.planning"); len(got) != 0 {
		t.Errorf("a valid planning block: %+v", got)
	}
	repo = planningProject(t, "planning:\n  currency: euro\n  hour_rate: -5\n  cycle: 2w\n")
	got := findings(t, repo, "manifest.planning")
	if len(got) != 3 {
		t.Fatalf("want three errors, currency, hour_rate, and cycle: %+v", got)
	}
	for i, want := range []string{`planning.currency "euro"`, "planning.hour_rate -5", `planning.cycle "2w"`} {
		f := got[i]
		if f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
			t.Errorf("finding %d: %+v, want an error on line 8 saying %q", i, f, want)
		}
	}
}

// S-0248: a front-matter problem whose message holds "; " is one finding on
// its key's line, not two with the second on line 1; two problems stay two.
func TestAFrontMatterProblemHoldingASemicolonIsOneFinding(t *testing.T) {
	repo := planningProject(t, "")
	s := mustItem(t, repo, workitem.Story, "Roles", "")
	s.Agent = &manifest.Agent{Harness: "claude-code", Roles: map[string]manifest.Role{"verify": {}}}
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	got := findings(t, repo, "item.front-matter")
	if len(got) != 1 {
		t.Fatalf("want one finding, the message whole: %+v", got)
	}
	f := got[0]
	if f.Line != keyLine(s.Path, "agent") || f.Line == 1 || !strings.HasPrefix(f.Message, "agent role verify sets nothing; give it") {
		t.Errorf("want an error on the agent line saying the role sets nothing: %+v", f)
	}
	s.Agent = &manifest.Agent{Harness: "Claude Code", Model: "opus 5"}
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	if got := findings(t, repo, "item.front-matter"); len(got) != 2 {
		t.Errorf("want two findings, the harness and the model: %+v", got)
	}
}

// S-0199: a malformed cost of delay is reported as the item's front matter,
// on the cost_of_delay line.
func TestAMalformedCostOfDelayIsReported(t *testing.T) {
	repo := planningProject(t, "")
	s := mustItem(t, repo, workitem.Story, "Costly", "")
	value := -3.0
	s.CostOfDelay = &workitem.CostOfDelay{Value: &value, By: "alex", At: "2026-09-01T12:00:00Z"}
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	got := findings(t, repo, "item.front-matter")
	if len(got) != 1 {
		t.Fatalf("want one finding, the message whole: %+v", got)
	}
	f := got[0]
	if f.Level != Error || f.Line != keyLine(s.Path, "cost_of_delay") || f.Line == 1 || !strings.HasPrefix(f.Message, "cost_of_delay.value -3 is negative: write") {
		t.Errorf("want an error on the cost_of_delay line saying the value is negative: %+v", f)
	}
}
