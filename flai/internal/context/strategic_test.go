package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
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

// answersRepo is strategicRepo with a digest and threads. Its digest holds
// an entry for each way the orchestrator's pack selects one or leaves it
// out:
//
//   - D-0001, an answer to the question TH-0001 asks again;
//   - D-0002, an answer to the question TH-0002 asks;
//   - D-0003, a decision, and D-0004, a cause, in TH-0001's words;
//   - D-0005, an absorbed answer in TH-0001's words;
//   - D-0006, a compacted recommendation in TH-0002's words;
//   - D-0007, an answer in the words of TH-0003 and TH-0004 alone;
//   - D-0008, the newest answer, sharing fewer words with TH-0001 than D-0001;
//   - D-0009, an answer in TH-0002's words, too large for the share of a
//     16 KB budget.
//
// TH-0001 and TH-0002 await the operator; the operator answered TH-0003
// last, and TH-0004 is resolved.
func answersRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	repo, _, _ := strategicRepo(t)
	entry := func(id, kind, title, status, date, body string) {
		absorbed := `""`
		if status == "absorbed" {
			absorbed = "design/system/orchestration.md"
		}
		text := "---\nid: " + id + "\ntitle: " + title + "\nkind: " + kind + "\nsources: [TH-0099]\ntopics: [orchestration]\npaths: []\ndate: " + date + "\nstatus: " + status + "\nabsorbed_by: " + absorbed + "\nasked_again: []\ncited: []\n---\n\n# " + id + " " + title + "\n"
		if body != "" {
			text += "\n" + body + "\n"
		}
		p := filepath.Join(repo.Root, "design", "digest", id+"-"+kind+".md")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	entry("D-0001", "answer", "Split a story whose forecast runs over five days", "active", "2026-10-01", "The operator answered: split a story whose forecast runs over five days into two stories.")
	entry("D-0002", "answer", "The dashboard publishes on the stable release channel", "active", "2026-10-02", "Release builds of the dashboard go out on the stable channel.")
	entry("D-0003", "decision", "Split a story whose forecast runs over five days", "active", "2026-10-03", "Split a story whose forecast runs over five days.")
	entry("D-0004", "cause", "A story whose forecast runs over five days was split late", "active", "2026-10-04", "Split a story whose forecast runs over five days.")
	entry("D-0005", "answer", "Split a story whose forecast runs over five days", "absorbed", "2026-10-05", "Split a story whose forecast runs over five days.")
	entry("D-0006", "recommendation", "Release channel for dashboard builds", "compacted", "2026-10-06", "")
	entry("D-0007", "answer", "Wombats are quokkas", "active", "2026-10-07", "Zebra wombats.")
	entry("D-0008", "answer", "Forecasts count days", "active", "2026-10-08", "Each forecast counts days.")
	entry("D-0009", "answer", "Stable release channel", "active", "2026-10-09", strings.TrimSpace(strings.Repeat("The stable release channel publishes dashboard builds. ", 60)))
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	open := func(title, by, text string) *threads.Thread {
		th, err := threads.New(repo, threads.NewOptions{Title: title, On: "design/system/orchestration.md", Author: by, Text: text, Now: now})
		if err != nil {
			t.Fatal(err)
		}
		return th
	}
	open("Split a story over five days?", "agent-S-0001", "Should the planner split a story whose forecast runs over five days?")
	open("Which channel publishes the dashboard?", "orchestrator", "Which release channel publishes the dashboard builds?")
	answered := open("Wombats or quokkas?", "agent-S-0002", "Are wombats quokkas?")
	if _, err := threads.Reply(repo, answered.ID, "designer", "Zebra.", now); err != nil {
		t.Fatal(err)
	}
	resolved := open("Quokkas and wombats again?", "agent-S-0003", "Are quokkas wombats?")
	if _, err := threads.Resolve(repo, resolved.ID, "agent-S-0003", "asked elsewhere", now); err != nil {
		t.Fatal(err)
	}
	return repo
}

// ADR-0194, S-0430 criterion 3: the orchestrator's pack loads the answers
// and recommendations ranked against each thread awaiting the operator,
// round by round, threads in ID order, each with its thread and rank; the
// answer to the question a thread asks again ranks first for it, ahead of
// a newer answer sharing fewer of its words. Decisions, causes, absorbed
// entries, and entries in the words of threads that do not await the
// operator never load.
func TestTheOrchestratorsPackLoadsTheAnswersRankedAgainstThreadsAwaitingTheOperator(t *testing.T) {
	repo := answersRepo(t)
	all, err := threads.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	awaiting := awaitingOperator(repo, all)
	if len(awaiting) != 2 || awaiting[0].ID != "TH-0001" || awaiting[1].ID != "TH-0002" {
		t.Fatalf("awaiting the operator: %d threads; want TH-0001 and TH-0002", len(awaiting))
	}
	entries, err := digest.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	tried := answersFor(entries, awaiting)
	var ids []string
	for _, it := range tried {
		ids = append(ids, it.ID)
	}
	if len(ids) < 2 || ids[0] != "D-0001" || tried[0].Reason != "thread TH-0001; rank 1" || !strings.HasPrefix(tried[1].Reason, "thread TH-0002; rank 1") {
		t.Fatalf("tried %v, first reasons %+v; want D-0001 for TH-0001 then TH-0002's best", ids, tried)
	}
	if i, j := slices.Index(ids, "D-0001"), slices.Index(ids, "D-0008"); j < 0 || i > j {
		t.Errorf("tried %v; want D-0008 after D-0001", ids)
	}
	for _, out := range []string{"D-0003", "D-0004", "D-0005", "D-0007"} {
		if slices.Contains(ids, out) {
			t.Errorf("tried %s: %v", out, ids)
		}
	}
	if len(ids) != 5 {
		t.Errorf("tried %v; want D-0001, D-0002, D-0006, D-0008, D-0009 once each", ids)
	}

	p, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	got := digestItems(p)
	if len(got) != 5 {
		t.Fatalf("digest items %v; want the five tried", got)
	}
	if one := got["D-0001"]; one.Reason != "thread TH-0001; rank 1" || one.Text != "# D-0001 Split a story whose forecast runs over five days\n\nThe operator answered: split a story whose forecast runs over five days into two stories.\n" || one.Size != len(one.Text) {
		t.Errorf("D-0001 %+v", one)
	}
	if c := got["D-0006"]; !strings.HasPrefix(c.Reason, "thread TH-0002; rank ") || !strings.HasSuffix(c.Reason, "; compacted") || c.Text != "" {
		t.Errorf("D-0006 %+v", c)
	}
	if r := got["D-0008"].Reason; !strings.HasPrefix(r, "thread TH-0001; rank ") || r == "thread TH-0001; rank 1" {
		t.Errorf("D-0008 reason %q", r)
	}
	body, header := p.Body(), p.Header()
	for _, want := range []string{strategicHead[groupDigest], "reason: thread TH-0001; rank 1\n\n# D-0001 Split a story", "\n- D-0006 Release channel for dashboard builds (design/digest/D-0006-recommendation.md; " + got["D-0006"].Reason + ")\n"} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %q:\n%s", want, body)
		}
	}
	for _, want := range []string{"the answers and recommendations of the digest that rank highest against the threads awaiting the operator", digestRow(5, 4)} {
		if !strings.Contains(header, want) {
			t.Errorf("header lacks %q:\n%s", want, header)
		}
	}
}

// ADR-0194: with no thread awaiting the operator, the orchestrator's pack
// loads no digest entry, though the digest holds answers in the words of
// threads the operator answered.
func TestTheOrchestratorsPackLoadsNoDigestWithNoThreadAwaitingTheOperator(t *testing.T) {
	repo := answersRepo(t)
	now := time.Date(2026, 10, 4, 9, 0, 0, 0, time.UTC)
	for _, id := range []string{"TH-0001", "TH-0002"} {
		if _, err := threads.Reply(repo, id, "designer", "Split it.", now); err != nil {
			t.Fatal(err)
		}
	}
	p, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := digestItems(p); len(got) != 0 || strings.Contains(p.Body(), "\ndigest\n======\n") {
		t.Errorf("digest items %v with no thread awaiting the operator", got)
	}
}

// ADR-0194: the orchestrator's digest step never takes the pack over its
// budget nor its items over an eighth of it; an entry too large for the
// share is skipped and the next tried.
func TestTheOrchestratorsDigestStaysWithinItsShare(t *testing.T) {
	repo := answersRepo(t)
	loaded := false
	for budget := 4096; budget <= 40*1024; budget += 512 {
		p, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", strconv.Itoa(budget))
		if err != nil {
			t.Fatal(err)
		}
		used := 0
		for _, it := range digestItems(p) {
			used += len(it.Printed())
		}
		switch {
		case p.Exceeded != "" && used > 0:
			t.Errorf("budget %d, exceeded %s: digest loads %d bytes", budget, p.Exceeded, used)
		case p.Exceeded == "" && (p.Size.Bytes > budget || used > budget/DigestShare):
			t.Errorf("budget %d: pack %d bytes, digest %d of a share of %d", budget, p.Size.Bytes, used, budget/DigestShare)
		}
		loaded = loaded || used > 0
	}
	if !loaded {
		t.Errorf("no budget loaded a digest entry")
	}
	p, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", "16KB")
	if err != nil {
		t.Fatal(err)
	}
	got := digestItems(p)
	if _, ok := got["D-0009"]; ok || got["D-0001"].ID == "" {
		t.Errorf("16 KB: digest items %v; want D-0001 and not D-0009", got)
	}
}

// ADR-0194: the planner's and the analyzer's packs take no digest step,
// though threads await the operator.
func TestThePlannerAndAnalyzerPacksLoadNoDigest(t *testing.T) {
	repo := answersRepo(t)
	stories, err := repo.List(false)
	if err != nil {
		t.Fatal(err)
	}
	story := ""
	for _, it := range stories {
		if it.Type == workitem.Story {
			story = it.ID
		}
	}
	plan, err := ForStrategic(repo, conventions.RolePlan, "", story, "")
	if err != nil {
		t.Fatal(err)
	}
	analyze, err := ForStrategic(repo, conventions.RoleAnalyze, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []*Pack{plan, analyze} {
		if got := digestItems(p); len(got) != 0 || strings.Contains(p.Header(), "the answers and recommendations") {
			t.Errorf("%s: digest items %v", p.Role, got)
		}
	}
}

// ADR-0194, S-0430 criterion 3: as --json prints it, each digest item of
// the orchestrator's pack has step digest, its ID, path, reason, and size.
func TestTheOrchestratorsDigestAsJSON(t *testing.T) {
	repo := answersRepo(t)
	p, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	var first map[string]any
	n := 0
	for _, it := range out.Items {
		if it["step"] == "digest" {
			if n == 0 {
				first = it
			}
			n++
		}
	}
	text := "# D-0001 Split a story whose forecast runs over five days\n\nThe operator answered: split a story whose forecast runs over five days into two stories.\n"
	if n != 5 || first["id"] != "D-0001" || first["path"] != "design/digest/D-0001-answer.md" || first["reason"] != "thread TH-0001; rank 1" || first["size"] != float64(len(text)) {
		t.Errorf("%d digest items, first %v", n, first)
	}
}
