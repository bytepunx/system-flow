package topics

import (
	"reflect"
	"strings"
	"testing"
)

const doc = `---
title: Code quality
topics: [code, go]
---

# Code quality

Intro.

## Rules

- A rule.

### Go <!-- topics: cli, go -->

- Go only.

#### Tests

- Inherited from Go.

## When in doubt

- Back to the file's.
`

func TestParseSectionsNestingAndInheritance(t *testing.T) {
	d, err := Parse(doc, []string{All})
	if err != nil {
		t.Fatal(err)
	}
	if !d.Declared || !reflect.DeepEqual(d.Topics, []string{"code", "go"}) {
		t.Fatalf("file topics = %v declared %v", d.Topics, d.Declared)
	}
	want := []struct {
		heading string
		level   int
		path    []string
		line    int
		own     []string
		topics  []string
	}{
		{"Code quality", 1, []string{"Code quality"}, 6, nil, []string{"code", "go"}},
		{"Rules", 2, []string{"Code quality", "Rules"}, 10, nil, []string{"code", "go"}},
		{"Go", 3, []string{"Code quality", "Rules", "Go"}, 14, []string{"cli", "go"}, []string{"cli", "go"}},
		{"Tests", 4, []string{"Code quality", "Rules", "Go", "Tests"}, 18, nil, []string{"cli", "go"}},
		{"When in doubt", 2, []string{"Code quality", "When in doubt"}, 22, nil, []string{"code", "go"}},
	}
	if len(d.Sections) != len(want) {
		t.Fatalf("got %d sections, want %d: %+v", len(d.Sections), len(want), d.Sections)
	}
	for i, w := range want {
		s := d.Sections[i]
		if s.Heading != w.heading || s.Level != w.level || !reflect.DeepEqual(s.Path, w.path) || s.Line != w.line || !reflect.DeepEqual(s.Own, w.own) || !reflect.DeepEqual(s.Topics, w.topics) {
			t.Errorf("section %d = %+v, want %+v", i, s, w)
		}
	}
	lines := strings.Split(doc, "\n")
	for _, s := range d.Sections {
		if !strings.HasPrefix(lines[s.Line-1], strings.Repeat("#", s.Level)+" ") {
			t.Errorf("line %d of %q is %q, not its heading", s.Line, s.Heading, lines[s.Line-1])
		}
	}
	var joined strings.Builder
	for _, s := range d.Sections {
		joined.WriteString(s.Text)
	}
	_, body, _ := strings.Cut(doc, "---\n\n")
	if joined.String() != body {
		t.Errorf("sections do not add up to the body:\n%s", joined.String())
	}
}

func TestParseDefaultWhenFileHasNoTopics(t *testing.T) {
	src := "---\ntitle: Git\n---\n\n# Git\n\n## Rules\n\n## Go <!-- topics: go -->\n"
	d, err := Parse(src, []string{All})
	if err != nil {
		t.Fatal(err)
	}
	if d.Declared || !reflect.DeepEqual(d.Topics, []string{All}) {
		t.Fatalf("a convention without topics is [all]: %v declared %v", d.Topics, d.Declared)
	}
	if got := d.Sections[1].Topics; !reflect.DeepEqual(got, []string{All}) {
		t.Errorf("Rules = %v, want [all]", got)
	}
	if got := d.Sections[2].Topics; !reflect.DeepEqual(got, []string{"go"}) {
		t.Errorf("Go = %v, want [go]", got)
	}
	d, err = Parse(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	if d.Topics != nil || d.Sections[1].Topics != nil {
		t.Errorf("a design file without topics has none: %v, %v", d.Topics, d.Sections[1].Topics)
	}
}

func TestHeadingCommentAtEndOfLine(t *testing.T) {
	cases := []struct {
		line  string
		level int
		text  string
		own   []string
		ok    bool
	}{
		{"## Go <!-- topics: cli, go -->", 2, "Go", []string{"cli", "go"}, true},
		{"## Go<!--topics:cli,go-->  ", 2, "Go", []string{"cli", "go"}, true},
		{"### Go <!-- topics: cli go cli -->", 3, "Go", []string{"cli", "go"}, true},
		{"## Go ##", 2, "Go", nil, true},
		{"## C#", 2, "C#", nil, true},
		{"## <!-- topics: a --> in the middle", 2, "<!-- topics: a --> in the middle", nil, true},
		{"## Go <!-- a note -->", 2, "Go <!-- a note -->", nil, true},
		{"#Go", 0, "", nil, false},
		{"Go <!-- topics: go -->", 0, "", nil, false},
	}
	for _, c := range cases {
		level, text, own, ok := Heading(c.line)
		if level != c.level || text != c.text || !reflect.DeepEqual(own, c.own) || ok != c.ok {
			t.Errorf("Heading(%q) = %d %q %v %v, want %d %q %v %v", c.line, level, text, own, ok, c.level, c.text, c.own, c.ok)
		}
	}
}

func TestFencedCodeIsNotAHeading(t *testing.T) {
	src := "# Doc\n\n```markdown\n# Not a heading <!-- topics: x -->\n```\n\n~~~\n## Nor this\n~~~\n\n## Real\n"
	d, err := Parse(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Sections) != 2 || d.Sections[1].Heading != "Real" || d.Sections[1].Line != 11 {
		t.Fatalf("sections = %+v", d.Sections)
	}
}

func TestTextBeforeTheFirstHeading(t *testing.T) {
	d, err := Parse("---\ntopics: [all]\n---\nPreamble.\n\n# Title\n", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Sections) != 2 || d.Sections[0].Level != 0 || d.Sections[0].Line != 4 || !reflect.DeepEqual(d.Sections[0].Topics, []string{All}) {
		t.Fatalf("sections = %+v", d.Sections)
	}
}

func TestTopicsThatAreNotAList(t *testing.T) {
	if _, err := Parse("---\ntopics: {a: b}\n---\n# T\n", nil); err == nil {
		t.Fatal("a map for topics should not parse")
	}
	d, err := Parse("# No front matter\n", []string{All})
	if err != nil || d.Declared || len(d.Sections) != 1 || d.Sections[0].Line != 1 {
		t.Fatalf("no front matter: %+v %v", d, err)
	}
}
