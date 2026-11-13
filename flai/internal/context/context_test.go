package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const quality = `---
title: Code quality
order: 60
topics: [all]
---

# Code quality

What every change carries.

## Rules

- Tests accompany the change.

### Go <!-- topics: go -->

- Standard library testing only.

#### Lint

- golangci-lint v2.

### Svelte <!-- topics: dashboard -->

- vitest.

## When in doubt <!-- topics: dashboard -->

- Shorter.

` + conventions.Marker + `

## Project additions <!-- topics: dashboard -->

- flaiover runs in Docker.
`

func TestFilterAllKeepsTheFileByteForByte(t *testing.T) {
	c, err := Filter("design/conventions/code-quality.md", quality, []string{"dashboard", "go", topics.All})
	if err != nil {
		t.Fatal(err)
	}
	if c.Text != quality {
		t.Errorf("text changed:\n%s", c.Text)
	}
	if len(c.LeftOut) != 0 || len(c.Omitted()) != 0 {
		t.Errorf("left out %v", c.LeftOut)
	}
	plain := "---\ntitle: X\n---\n\n# X\n\n- a rule\n"
	if c, _ := Filter("x.md", plain, []string{topics.All}); c.Text != plain || !reflect.DeepEqual(c.Topics, []string{topics.All}) {
		t.Errorf("no topics reads as all: %q %v", c.Text, c.Topics)
	}
}

func TestFilterLeavesOutSectionsAndKeepsMarkerAndAdditions(t *testing.T) {
	c, err := Filter("design/conventions/code-quality.md", strings.Replace(quality, "topics: [all]", "topics: [code]", 1), []string{"cli", "go", "code", topics.All})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"title: Code quality", "# Code quality", "### Go <!-- topics: go -->", "#### Lint", conventions.Marker + "\n\n## Project additions", "- Tests accompany"} {
		if !strings.Contains(c.Text, want) {
			t.Errorf("missing %q:\n%s", want, c.Text)
		}
	}
	for _, gone := range []string{"vitest", "Svelte", "When in doubt", "Shorter", "Docker"} {
		if strings.Contains(c.Text, gone) {
			t.Errorf("kept %q:\n%s", gone, c.Text)
		}
	}
	want := []string{
		"code-quality.md § Rules › Svelte (dashboard)",
		"code-quality.md § When in doubt (dashboard)",
		"code-quality.md § Project additions (dashboard)",
	}
	if got := c.Omitted(); !reflect.DeepEqual(got, want) {
		t.Errorf("omitted %q", got)
	}
	if len(c.Kept) != 4 || c.Kept[2].Heading != "Go" || !reflect.DeepEqual(c.Kept[3].Topics, []string{"go"}) {
		t.Errorf("kept %+v", c.Kept)
	}
}

func TestFilterKeepsTheHeadingAboveAKeptSection(t *testing.T) {
	doc := "# T <!-- topics: dashboard -->\n\nIntro.\n\n## Rules\n\n- r\n\n### Go <!-- topics: go -->\n\n- g\n\n#### Lint\n\n- l\n"
	c, err := Filter("t.md", doc, []string{"go", topics.All})
	if err != nil {
		t.Fatal(err)
	}
	want := "# T <!-- topics: dashboard -->\n## Rules\n### Go <!-- topics: go -->\n\n- g\n\n#### Lint\n\n- l\n"
	if c.Text != want {
		t.Errorf("text:\n%q", c.Text)
	}
	if got := c.Omitted(); !reflect.DeepEqual(got, []string{"t.md § T (dashboard)"}) {
		t.Errorf("omitted %q", got)
	}
	if len(c.LeftOut) != 2 {
		t.Errorf("json keeps every section left out: %+v", c.LeftOut)
	}
}

// I-0261: a story whose 40 touches reach one sub-project of three topics
// gets a pack whose topics encode in under 2 KB, where one claim source per
// touches entry and topic took about 10 KB.
func TestAPacksTopicsFoldAStorysTouches(t *testing.T) {
	touches := make([]string, 40)
	for i := range touches {
		touches[i] = fmt.Sprintf("flai/internal/pkg%02d/file.go", i)
	}
	repo, story := storyProject(t, "", touches)
	p, err := ForStory(repo, story.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if names := topics.Names(p.Topics); !reflect.DeepEqual(names, []string{"flai", "cli", "go", topics.Code, topics.All}) {
		t.Fatalf("topics %v", names)
	}
	data, err := json.Marshal(p.Topics)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) >= 2048 {
		t.Errorf("topics encode in %d bytes, want under 2048:\n%s", len(data), data)
	}
}

// storyProject is a project of one Go sub-project and one story with
// touches, whose manifest ends with extra, and that story.
func storyProject(t *testing.T, extra string, touches []string) (*workitem.Repo, *workitem.Item) {
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
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: flai\n    path: flai\n    kind: go\n    tags: [cli]\n"+extra)
	write("design/conventions/README.md", "# Conventions\n\n- [code-quality.md](code-quality.md)\n")
	write("design/conventions/code-quality.md", quality)
	for _, d := range []string{"flai", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "A wide story", Owner: "alex", Touches: touches,
		Now: time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC), Body: "## Goal\n\nBuild it.\n\n## Acceptance criteria\n\n- [ ] built\n"})
	if err != nil {
		t.Fatal(err)
	}
	return repo, story
}

// S-0662: a story's pack in a project whose kluster: block names a deploy
// names its namespace and the commands its agent deploys and checks it
// with, in the header and in the JSON the MCP prime tool returns.
func TestAStorysPackInAKlusterProjectNamesTheDeployAndChecksCommands(t *testing.T) {
	repo, story := storyProject(t, "kluster:\n  config: kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n", []string{"flai/main.go"})
	p, err := ForStory(repo, story.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	ns := strings.ToLower(story.ID)
	h := p.Header()
	for _, want := range []string{
		"kluster: the story's namespace is " + ns + "\n",
		"  deploy  flai kluster deploy " + story.ID + " --job, or the MCP tool kluster_deploy\n",
		"  checks  flai checks run " + story.ID + " --job, or the MCP tool kluster_checks\n",
		"  exec    flai kluster exec " + story.ID + " -- <command>\n",
		"  env     flai kluster env " + story.ID + "\n",
		"  log     flai kluster log " + story.ID + " --step <step>\n",
		"  " + klusterReview + "\n",
	} {
		if !strings.Contains(h, want) {
			t.Errorf("header lacks %q:\n%s", want, h)
		}
	}
	part, err := p.Part(1, PartLimit)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(part)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Kluster map[string]string `json:"kluster"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"namespace":   ns,
		"deploy":      "flai kluster deploy " + story.ID + " --job",
		"deploy_tool": "kluster_deploy",
		"checks":      "flai checks run " + story.ID + " --job",
		"checks_tool": "kluster_checks",
		"exec":        "flai kluster exec " + story.ID + " -- <command>",
		"env":         "flai kluster env " + story.ID,
		"log":         "flai kluster log " + story.ID + " --step <step>",
		"review":      klusterReview,
	}
	if !reflect.DeepEqual(got.Kluster, want) {
		t.Errorf("kluster in the JSON %v, want %v", got.Kluster, want)
	}
	if !strings.Contains(klusterReview, "flai serve enable kluster") || !strings.Contains(klusterReview, "Notes") {
		t.Errorf("review sentence %q names neither the host action nor Notes", klusterReview)
	}
}

// S-0662: a project with no kluster: block, or one that names no deploy,
// gets neither the field nor the header's block.
func TestAStorysPackWithoutAKlusterDeployNamesNoCommands(t *testing.T) {
	for name, extra := range map[string]string{
		"no block":                  "",
		"no deploy":                 "kluster:\n  config: kluster.yaml\n",
		"deploys before review off": "kluster:\n  config: kluster.yaml\n  deploy: [scripts/deploy.sh]\n  agent_namespaces: 0\n",
	} {
		t.Run(name, func(t *testing.T) {
			repo, story := storyProject(t, extra, []string{"flai/main.go"})
			p, err := ForStory(repo, story.ID, "")
			if err != nil {
				t.Fatal(err)
			}
			if p.Kluster != nil {
				t.Errorf("kluster %+v", p.Kluster)
			}
			if h := p.Header(); strings.Contains(h, "kluster") {
				t.Errorf("header names the kluster:\n%s", h)
			}
			data, err := json.Marshal(p)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), `"kluster"`) {
				t.Errorf("JSON carries a kluster field")
			}
		})
	}
}

func TestBuildWithEverythingKeptIsPrimeCat(t *testing.T) {
	set := &conventions.Set{
		Dir:    "/r/design/conventions",
		README: "# Conventions\n",
		Files:  []conventions.File{{Path: "design/conventions/code-quality.md", Name: "code-quality.md", Title: "Code quality", Order: 60, Raw: quality}},
	}
	table := "| ID |\n|----|\n| I-1 |\n"
	st := []topics.StoryTopic{
		{Topic: "go", Sources: []topics.Source{{Kind: topics.FromOwn, Item: "S-0001"}}},
		{Topic: "dashboard", Sources: []topics.Source{{Kind: topics.FromEpic, Item: "E-0001"}}},
		{Topic: topics.All, Sources: []topics.Source{{Kind: topics.FromAll}}},
	}
	p, err := Build("/r", "S-0001", "A story", st, set, table, DefaultBudget)
	if err != nil {
		t.Fatal(err)
	}
	want := "design/conventions/README.md\n============================\n\n# Conventions\n\n" +
		"design/conventions/code-quality.md\n==================================\n\n" + quality +
		"\nopen issues\n===========\n\n" + table
	if p.Body() != want {
		t.Errorf("body is not prime --cat:\n%s", p.Body())
	}
	h := p.Header()
	if p.Size.Bytes != len(h)+len(want) || p.Size.Lines != strings.Count(h+want, "\n") || p.OpenIssues != 1 {
		t.Errorf("size %+v, open issues %d", p.Size, p.OpenIssues)
	}
	for _, s := range []string{
		"S-0001 context pack\n===================\n\n",
		"story: S-0001 A story\n",
		"  go         own\n",
		"  dashboard  epic E-0001\n",
		fmt.Sprintf("size: %d bytes, %d lines, this header included; budget %d bytes\n", len(h)+len(want), strings.Count(h+want, "\n"), DefaultBudget),
		fmt.Sprintf("  %6d  convention  design/conventions/code-quality.md\n", len(quality)),
		"      22  issues      1 open\n",
	} {
		if !strings.Contains(h, s) {
			t.Errorf("header lacks %q:\n%s", s, h)
		}
	}
	if strings.Contains(h, "left out") || strings.Contains(h, "over budget") {
		t.Errorf("nothing was left out:\n%s", h)
	}

	p, _ = Build("/r", "S-0001", "A story", st[:1], set, table, DefaultBudget)
	if !strings.HasSuffix(p.Body(), "\nleft out\n========\n\n- code-quality.md § Rules › Svelte (dashboard)\n- code-quality.md § When in doubt (dashboard)\n- code-quality.md § Project additions (dashboard)\n") {
		t.Errorf("left out list:\n%s", p.Body())
	}
	if !strings.Contains(p.Header(), "left out    3 convention sections, listed at the end\n") {
		t.Errorf("header:\n%s", p.Header())
	}
}
