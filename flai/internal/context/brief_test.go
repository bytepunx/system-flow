package context

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/topics"
)

func TestDecisionSentence(t *testing.T) {
	adr := func(decision string) string {
		return "---\nid: ADR-0001\ntitle: T\n---\n\n# ADR-0001 T\n\n## Context\n\nWhy.\n\n## Decision\n\n" + decision + "\n\n## Consequences\n\nSome.\n"
	}
	for _, c := range []struct {
		decision, want string
		ok             bool
	}{
		{"**A pack fits a budget.** Conventions load whole.", "A pack fits a budget.", true},
		{"**MCP over HTTP.** flaiover serves `/mcp` at the root. More.", "MCP over HTTP. flaiover serves `/mcp` at the root.", true},
		{"flai reads `config.json` and\nwrites it back. Then more.", "flai reads `config.json` and writes it back.", true},
		{"Version 1.2.3 of `a.b` is pinned!", "Version 1.2.3 of `a.b` is pinned!", true},
		{"- a list first.", "", false},
		{"| a | table |", "", false},
		{"1. numbered.", "", false},
		{"### A heading", "", false},
		{"No full stop at all", "No full stop at all", false},
	} {
		got, ok := DecisionSentence(adr(c.decision))
		if got != c.want || ok != c.ok {
			t.Errorf("%q: got %q %v, want %q %v", c.decision, got, ok, c.want, c.ok)
		}
	}
	if _, ok := DecisionSentence("---\nid: ADR-0001\n---\n\n# ADR-0001\n\n## Context\n\nOnly.\n"); ok {
		t.Error("an ADR without a Decision section has a sentence")
	}
}

func TestFirstParagraph(t *testing.T) {
	for raw, want := range map[string]string{
		"---\ntitle: X\n---\n\n# X\n\nFirst line\nand second.\n\nSecond paragraph.\n": "First line and second.",
		"---\ntitle: X\n---\n\n# X\n\n## Part\n\nUnder a heading.\n":                  "Under a heading.",
		"No front matter at all.\n":   "No front matter at all.",
		"---\ntitle: X\n---\n\n# X\n": "",
	} {
		d, err := ParseDoc("design/system/x.md", KindSystem, raw)
		if err != nil {
			t.Fatal(err)
		}
		if got := d.FirstParagraph(); got != want {
			t.Errorf("%q: got %q, want %q", raw, got, want)
		}
	}
}

func TestParseSize(t *testing.T) {
	for in, want := range map[string]int{"81920": 81920, "80KB": 81920, "80 kb": 81920, "80k": 81920, "80KiB": 81920, "1MB": 1 << 20, "100B": 100} {
		if got, err := ParseSize(in); err != nil || got != want {
			t.Errorf("ParseSize(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	for _, in := range []string{"", "0", "-5KB", "80GB", "eighty", "1.5KB"} {
		if _, err := ParseSize(in); err == nil {
			t.Errorf("ParseSize(%q) took it", in)
		}
	}
	if n, err := Budget("", ""); n != DefaultBudget || err != nil {
		t.Errorf("default %d %v", n, err)
	}
	if n, _ := Budget("", "100KB"); n != 100*1024 {
		t.Errorf("project %d", n)
	}
	if n, _ := Budget("10KB", "100KB"); n != 10*1024 {
		t.Errorf("flag %d", n)
	}
	if _, err := Budget("", "lots"); err == nil || !strings.Contains(err.Error(), "80KB") {
		t.Errorf("a bad project budget: %v", err)
	}
}

// budgetPack is a pack over the fixture documents with one small
// convention, for a story named by the sources, with the given budget.
func budgetPack(t *testing.T, budget int, storyTopics []string, sources []Source, query string) (*Pack, *Selection) {
	t.Helper()
	set := &conventions.Set{Dir: "/r/design/conventions", Files: []conventions.File{{Path: "design/conventions/a.md", Name: "a.md", Title: "A", Order: 10, Raw: "# A\n\n- Rule.\n"}}}
	var st []topics.StoryTopic
	for _, n := range storyTopics {
		st = append(st, topics.StoryTopic{Topic: n, Sources: []topics.Source{{Kind: topics.FromOwn, Item: "S-0001"}}})
	}
	p, err := Build("/r", "S-0001", "A story", st, set, "", budget)
	if err != nil {
		t.Fatal(err)
	}
	s := Design(fixture(t).Docs, storyTopics, sources, BriefOver(budget))
	p.AddDesign(s, query)
	return p, s
}

func ranked(p *Pack) []string {
	var out []string
	for _, it := range p.Items {
		if it.Step == StepRanked {
			l := it.Path
			if it.Label != "" {
				l += " § " + it.Label
			}
			out = append(out, l+" = "+it.Reason)
		}
	}
	return out
}

func TestRankedSectionsFillTheBudgetInRankOrder(t *testing.T) {
	query := "herons on the estuary at the turn of the tide"
	full, _ := budgetPack(t, DefaultBudget, []string{"cli"}, nil, query)
	all := ranked(full)
	if len(all) < 8 || all[0] != "design/system/plain.md § Heron migration = rank 1" {
		t.Fatalf("with room: %v", all)
	}
	if full.Size.Bytes > full.Budget || full.Exceeded != "" {
		t.Errorf("size %+v over %d, exceeded %q", full.Size, full.Budget, full.Exceeded)
	}

	base, _ := budgetPack(t, DefaultBudget, []string{"cli"}, nil, "")
	tight := base.Size.Bytes + 700
	p, s := budgetPack(t, tight, []string{"cli"}, nil, query)
	got := ranked(p)
	if len(got) == 0 || len(got) >= len(all) {
		t.Fatalf("a tight budget ranked %v of %v", got, all)
	}
	if got[0] != all[0] {
		t.Errorf("the best did not load first: %v", got)
	}
	if p.Size.Bytes > tight {
		t.Errorf("size %d over %d", p.Size.Bytes, tight)
	}
	if strings.Count(p.Header(), "  ranked  ") != len(got) {
		t.Errorf("header rows:\n%s", p.Header())
	}
	loaded := 0
	for _, d := range s.Docs {
		if s.Loaded(d) {
			loaded++
		}
	}
	if len(p.Catalog.NotLoaded) == 0 || loaded == len(s.Docs) {
		t.Errorf("nothing left for the catalog: %+v", p.Catalog)
	}
}

func TestTheBriefsAreKeptWhenTheyTakeThePackOverTheBudget(t *testing.T) {
	base, _ := budgetPack(t, DefaultBudget, []string{"cli", "go", "dashboard"}, nil, "")
	named := 0
	for _, it := range base.Items {
		if it.Step != StepBriefed {
			t.Fatalf("only briefs expected: %+v", it)
		}
		named++
	}
	p, _ := budgetPack(t, base.Size.Bytes-100, []string{"cli", "go", "dashboard"}, nil, "herons")
	if p.Exceeded != ExceededBriefs || len(p.Items) != named || len(ranked(p)) != 0 {
		t.Errorf("exceeded %q, %d items of %d, ranked %v", p.Exceeded, len(p.Items), named, ranked(p))
	}
	if !strings.Contains(p.Header(), "over budget: the conventions, what is named, and the briefs exceed it, so nothing is ranked; every brief is kept") {
		t.Errorf("header:\n%s", p.Header())
	}
}

func TestTheHeaderSaysWhatAloneExceedsTheBudget(t *testing.T) {
	p, _ := budgetPack(t, 1, []string{"cli"}, []Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0002"}}, "herons")
	if p.Exceeded != ExceededConventions || len(p.Items) != 0 || len(p.Catalog.NotLoaded) != 17 {
		t.Errorf("exceeded %q, %d items, catalog %d", p.Exceeded, len(p.Items), len(p.Catalog.NotLoaded))
	}
	if !strings.Contains(p.Header(), "over budget: the conventions alone exceed it") {
		t.Errorf("header:\n%s", p.Header())
	}

	conv, err := Build("/r", "S-0001", "A story", nil, &conventions.Set{Dir: "/r/design/conventions", Files: []conventions.File{{Path: "design/conventions/a.md", Name: "a.md", Title: "A", Order: 10, Raw: "# A\n\n- Rule.\n"}}}, "", 1)
	if err != nil {
		t.Fatal(err)
	}
	p, _ = budgetPack(t, conv.Size.Bytes+300, nil, []Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0002 and design/system/cli.md"}}, "herons")
	if p.Exceeded != ExceededNamed || len(ranked(p)) != 0 {
		t.Errorf("exceeded %q, ranked %v", p.Exceeded, ranked(p))
	}
	if !strings.Contains(p.Header(), "over budget: the conventions and what the story, its epic, and its tasks name exceed it") ||
		!strings.Contains(p.Header(), "  named       design/adrs/0002-second.md\n") {
		t.Errorf("header:\n%s", p.Header())
	}
}

func TestAPackSaysABriefOfWhatTheStoryNamesIsForItsSize(t *testing.T) {
	task := Source{ID: "T-0001", Path: "wip/kanban/tasks/T-0001-x.md", Body: "Update design/system/cli.md."}
	p, _ := budgetPack(t, 2000, nil, []Source{task}, "")
	var it *Item
	for i := range p.Items {
		if p.Items[i].Path == "design/system/cli.md" {
			it = &p.Items[i]
		}
	}
	if it == nil || it.Step != StepBriefed || it.Reason != "named in T-0001" {
		t.Fatalf("items %+v", p.Items)
	}
	body := p.Body()
	if !strings.Contains(body, "\ndesign/system/cli.md: The CLI (328 bytes; named in T-0001)\n"+namedBrief+"\n") {
		t.Errorf("brief not marked as named:\n%s", body)
	}
	if !strings.Contains(body, "name by a path written out and that are too large to load whole") {
		t.Errorf("briefs heading does not say what a named brief is:\n%s", body)
	}

	p, _ = budgetPack(t, DefaultBudget, nil, []Source{task}, "")
	if len(p.Items) == 0 || p.Items[0].Path != "design/system/cli.md" || p.Items[0].Step != StepNamed || p.Items[0].Heading != nil {
		t.Errorf("a document under an eighth of the budget is not loaded whole: %+v", p.Items)
	}
}
