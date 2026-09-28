package context

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/topics"
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
	p, err := Build("/r", "S-0001", "A story", st, set, table)
	if err != nil {
		t.Fatal(err)
	}
	want := "design/conventions/README.md\n============================\n\n# Conventions\n\n" +
		"design/conventions/code-quality.md\n==================================\n\n" + quality +
		"\nopen issues\n===========\n\n" + table
	if p.Body() != want {
		t.Errorf("body is not prime --cat:\n%s", p.Body())
	}
	if p.Size.Bytes != len(want) || p.Size.Lines != strings.Count(want, "\n") || p.OpenIssues != 1 {
		t.Errorf("size %+v, open issues %d", p.Size, p.OpenIssues)
	}
	h := p.Header()
	for _, s := range []string{
		"S-0001 context pack\n===================\n\n",
		"story: S-0001 A story\n",
		"  go         own\n",
		"  dashboard  epic E-0001\n",
		fmt.Sprintf("size: %d bytes, %d lines below this header\n\n", len(want), strings.Count(want, "\n")),
	} {
		if !strings.Contains(h, s) {
			t.Errorf("header lacks %q:\n%s", s, h)
		}
	}
	if strings.Contains(h, "left out") {
		t.Errorf("nothing was left out:\n%s", h)
	}

	p, _ = Build("/r", "S-0001", "A story", st[:1], set, table)
	if !strings.HasSuffix(p.Body(), "\nleft out\n========\n\n- code-quality.md § Rules › Svelte (dashboard)\n- code-quality.md § When in doubt (dashboard)\n- code-quality.md § Project additions (dashboard)\n") {
		t.Errorf("left out list:\n%s", p.Body())
	}
	if !strings.Contains(p.Header(), "left out: 3 sections, listed at the end\n") {
		t.Errorf("header:\n%s", p.Header())
	}
}
