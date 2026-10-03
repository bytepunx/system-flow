package context

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// strategicRepo is a project with conventions read by different roles, a
// design file for each strategic role's topic, an open issue, and an epic
// whose story names a design file.
func strategicRepo(t *testing.T) (*workitem.Repo, *workitem.Item, *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	write := func(rel, s string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	conv := func(name, title, order, roles string) {
		write("design/conventions/"+name, "---\ntitle: "+title+"\nupdated: 2026-10-03\naudience: agent\norder: "+order+"\nstatus: active\n"+roles+"---\n\n# "+title+"\n\n## Rules\n\n- "+title+" rule.\n\n### Svelte <!-- topics: svelte -->\n\n- "+title+" svelte rule.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n")
	}
	write("design/conventions/README.md", "# Conventions\n\n- [safety.md](safety.md)\n")
	conv("safety.md", "Safety", "10", "")
	conv("stream.md", "Stream", "20", "roles: [story]\n")
	conv("strategy.md", "Strategy", "30", "roles: [plan, orchestrate, analyze]\n")
	conv("planning.md", "Planning", "40", "roles: [plan]\n")
	conv("orchestration.md", "Orchestration", "50", "roles: [orchestrate]\n")
	conv("analysis.md", "Analysis", "60", "roles: [analyze]\n")
	doc := func(name, title, topic string) {
		write("design/system/"+name+".md", "---\ntitle: "+title+"\nupdated: 2026-10-03\nstatus: active\n"+topic+"---\n\n# "+title+"\n\nThe "+name+" design.\n\n## Part\n\nThe body of the "+name+" part.\n")
	}
	doc("planning", "Planning", "topics: [planning]\n")
	doc("orchestration", "Orchestration", "topics: [orchestration]\n")
	doc("analysis", "Analysis", "topics: [analysis]\n")
	doc("named", "Named", "")
	doc("epic-named", "Epic named", "")
	write("design/issues/I-0001-slow.md", "---\nid: I-0001\ntitle: Slow\nclass: efficiency\nstatus: open\ncount: 1\ncost: 4m\nfirst_reported: 2026-10-01T08:00:00Z\nlast_reported: 2026-10-01T08:00:00Z\nupdated: 2026-10-01T08:00:00Z\n---\n\n# I-0001 Slow\n\n## Description\nSlow.\n")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Strategic agents", Owner: "alex", Topics: []string{"svelte"}, Now: now,
		Body: "## Goal\n\nKeep to [epic named](../../../design/system/epic-named.md).\n"})
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Roles prime", Parent: epic.ID, Owner: "alex", Topics: []string{"planning"}, Now: now,
		Body: "## Goal\n\nKeep to [named](../../../design/system/named.md).\n\n## Acceptance criteria\n\n- [ ] it primes\n"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "A task", Parent: story.ID, Owner: "alex", Now: now}); err != nil {
		t.Fatal(err)
	}
	return repo, epic, story
}

// paths is the paths of the pack's conventions, and of its items with
// their steps.
func paths(p *Pack) (conv []string, items map[string]string) {
	items = map[string]string{}
	for _, c := range p.Conventions {
		conv = append(conv, c.Path)
	}
	for _, it := range p.Items {
		items[it.Path] = it.Step
	}
	return conv, items
}

// topicOf is the pack's topic with this name, or nil.
func topicOf(p *Pack, name string) *topics.StoryTopic {
	for i := range p.Topics {
		if p.Topics[i].Topic == name {
			return &p.Topics[i]
		}
	}
	return nil
}

// E-0016, S-0207: each strategic role's pack has the conventions whose
// roles are empty or list it, without the README, the open issues, its
// role's topic, and briefs of the design that topic selects.
func TestForStrategicProjectPacks(t *testing.T) {
	repo, _, _ := strategicRepo(t)
	for role, want := range map[string]struct {
		conv  []string
		brief string
		not   []string
	}{
		"orchestrate": {[]string{"safety.md", "strategy.md", "orchestration.md"}, "design/system/orchestration.md", []string{"design/system/planning.md", "design/system/analysis.md"}},
		"analyze":     {[]string{"safety.md", "strategy.md", "analysis.md"}, "design/system/analysis.md", []string{"design/system/planning.md", "design/system/orchestration.md"}},
	} {
		p, err := ForStrategic(repo, role, "", "", "")
		if err != nil {
			t.Fatalf("%s: %v", role, err)
		}
		conv, items := paths(p)
		for i, c := range conv {
			conv[i] = filepath.Base(c)
		}
		if !slices.Equal(conv, want.conv) || p.README != nil {
			t.Errorf("%s conventions: %v, readme %v", role, conv, p.README != nil)
		}
		if items[want.brief] != StepBriefed {
			t.Errorf("%s: %s not briefed: %v", role, want.brief, items)
		}
		for _, n := range want.not {
			if _, ok := items[n]; ok {
				t.Errorf("%s: chose %s: %v", role, n, items)
			}
		}
		for _, it := range p.Items {
			if it.Step == StepRanked || it.Step == StepNamed {
				t.Errorf("%s: %s is %s; a project pack names and ranks nothing", role, it.Path, it.Step)
			}
		}
		top := topicOf(p, roleTopic[role])
		if p.Role != role || p.Story != "" || p.Item != "" || p.OpenIssues != 1 || top == nil || top.Sources[0].Kind != FromRole || len(p.Catalog.NotLoaded) == 0 {
			t.Errorf("%s pack: role %q story %q item %q issues %d topics %+v catalog %v", role, p.Role, p.Story, p.Item, p.OpenIssues, p.Topics, p.Catalog)
		}
		if p.Budget != DefaultBudget {
			t.Errorf("%s budget %d, want the story's agent's %d", role, p.Budget, DefaultBudget)
		}
		head := p.Header() + p.Body()
		for _, w := range []string{"Project context pack for the " + role + " role\n", "role: " + role + ", the ", "for the whole project", "open issues", "- Strategy rule."} {
			if !strings.Contains(head, w) {
				t.Errorf("%s: missing %q:\n%s", role, w, head)
			}
		}
		for _, gone := range []string{"story:", "Stream rule", "svelte rule", "Planning rule"} {
			if strings.Contains(head, gone) {
				t.Errorf("%s: printed %q", role, gone)
			}
		}
	}
}

// E-0016, S-0207: the planner's pack for an epic has its topics with the
// role's, and loads what the epic names, not what its stories name.
func TestForStrategicPlansAnEpic(t *testing.T) {
	repo, epic, _ := strategicRepo(t)
	p, err := ForStrategic(repo, "plan", epic.ID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	conv, items := paths(p)
	if strings.Join(conv, ",") != "design/conventions/safety.md,design/conventions/strategy.md,design/conventions/planning.md" || p.README != nil {
		t.Errorf("conventions: %v", conv)
	}
	if items["design/system/epic-named.md"] != StepNamed || items["design/system/named.md"] == StepNamed || items["design/system/planning.md"] == "" {
		t.Errorf("items: %v", items)
	}
	if p.Item != epic.ID || p.Story != "" || p.Title != "Strategic agents" || p.OpenIssues != 1 || p.Topics[0].Topic != "planning" || topicOf(p, "svelte") == nil {
		t.Errorf("pack: item %q story %q title %q issues %d topics %+v", p.Item, p.Story, p.Title, p.OpenIssues, p.Topics)
	}
	out := p.Header() + p.Body()
	for _, w := range []string{epic.ID + " context pack for the plan role\n", "item: " + epic.ID + " Strategic agents\n", "role: plan, the planner, for " + epic.ID, "- Planning svelte rule."} {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q:\n%s", w, out)
		}
	}
}

// E-0016, S-0207: the planner's pack for a story has the story's topics with
// the role's, the role a source of a topic they share, and loads what the
// story names.
func TestForStrategicPlansAStory(t *testing.T) {
	repo, _, story := strategicRepo(t)
	p, err := ForStrategic(repo, "plan", "", story.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	_, items := paths(p)
	if items["design/system/named.md"] != StepNamed || items["design/system/planning.md"] == "" {
		t.Errorf("items: %v", items)
	}
	top := topicOf(p, "planning")
	if p.Item != story.ID || top == nil || len(top.Sources) != 2 || top.Sources[0].Kind != topics.FromOwn || top.Sources[1].Kind != FromRole || topicOf(p, "svelte") == nil {
		t.Errorf("pack: item %q topics %+v", p.Item, p.Topics)
	}
}

// E-0016, S-0207: a strategic role refuses what it is not primed for.
func TestForStrategicRefusals(t *testing.T) {
	repo, epic, story := strategicRepo(t)
	for _, c := range []struct{ role, epic, story, want string }{
		{"plan", "", "", "give --epic E-nnnn or --story S-nnnn"},
		{"plan", "", "T-0001", "T-0001 is task; give --epic E-nnnn or --story S-nnnn"},
		{"plan", "", "S-0999", "give --epic E-nnnn or --story S-nnnn"},
		{"plan", story.ID, "", "--epic " + story.ID + ": " + story.ID + " is story"},
		{"plan", "", epic.ID, "--story " + epic.ID + ": " + epic.ID + " is epic"},
		{"plan", epic.ID, story.ID, "give --epic or --story, not both"},
		{"orchestrate", "", story.ID, "the orchestrator primes for the whole project; give no --story or --epic"},
		{"analyze", epic.ID, "", "the analyzer primes for the whole project"},
		{"explore", "", "", "no such role; a strategic agent's role is plan, orchestrate, analyze"},
	} {
		if _, err := ForStrategic(repo, c.role, c.epic, c.story, ""); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %s%s: %v, want %q", c.role, c.epic, c.story, err, c.want)
		}
	}
}
