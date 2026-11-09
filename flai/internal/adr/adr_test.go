package adr

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/adrindex"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0222: a quoted title keeps <, >, and & as themselves, and titleOf reads
// it back, escapes and all, as YAML reads it.
func TestYamlStringWritesNoHTMLEscapes(t *testing.T) {
	for _, c := range []struct{ title, want string }{
		{"Plain words", "Plain words"},
		{"Tokens: one per project", `"Tokens: one per project"`},
		{"Escaped as < and >", `"Escaped as < and >"`},
		{"Rock & roll", `"Rock & roll"`},
		{`Says "hi"`, `"Says \"hi\""`},
	} {
		got := yamlString(c.title)
		if got != c.want {
			t.Errorf("yamlString(%q) = %s, want %s", c.title, got, c.want)
		}
		if back := titleOf("id: ADR-0001\ntitle: " + got + "\nstatus: proposed"); back != c.title {
			t.Errorf("titleOf reads %s back as %q, want %q", got, back, c.title)
		}
	}
}

func TestTitleOfReadsAnEscapeAsTheCharacter(t *testing.T) {
	if got := titleOf("title: \"a \\u003c b \\u003e c \\u0026 d\"\n"); got != "a < b > c & d" {
		t.Errorf("titleOf = %q", got)
	}
	if got := titleOf("title: \"unclosed\n"); got != "unclosed" {
		t.Errorf("a quote YAML would refuse is trimmed as before: %q", got)
	}
}

// noGit is a runner with no git, so an ADR's number comes from the folder.
type noGit struct{}

var errNoGit = errors.New("no git in this test")

func (noGit) Run(string, string, ...string) (string, error)              { return "", errNoGit }
func (noGit) RunInput(string, string, string, ...string) (string, error) { return "", errNoGit }
func (noGit) LookPath(string) (string, error)                            { return "", errNoGit }

func adrText(n, title, status string) string {
	return "---\nid: ADR-" + n + "\ntitle: " + title + "\nstatus: " + status + "\ndate: 2026-10-01\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-" + n + " " + title + "\n\n## Context\n\nc\n\n## Decision\n\nWe decide.\n\n## Consequences\n\ne\n\n## Alternatives considered\n\nf\n"
}

// indexRepo is a project with ADRs 0001 and 0002 and an index whose row for
// 0001 carries a note its front matter does not hold.
func indexRepo(t *testing.T) (*workitem.Repo, string) {
	t.Helper()
	root := t.TempDir()
	write := func(rel, text string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write("design/adrs/0001-first.md", adrText("0001", "First", "accepted"))
	write("design/adrs/0002-second.md", adrText("0002", "Second", "accepted"))
	write("design/adrs/README.md", "# Architecture Decision Records\n\nGenerated.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n"+
		"| [0001](0001-first.md) | First | accepted, refined by 0002 |\n| [0002](0002-second.md) | Second | accepted |\n\nBelow.\n")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo, filepath.Join(root, "design", "adrs")
}

// indexIsRendered fails unless the index is what adrindex renders from the
// folder, and returns it.
func indexIsRendered(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := adrindex.Render(os.DirFS(dir), string(data))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Errorf("the index is not what flai adr index writes:\n%s\nwant:\n%s", data, want)
	}
	return string(data)
}

// ADR-0228: flai adr new, the supersede it does, and flai adr accept write
// the whole index again from the ADR files, dropping a note the front matter
// does not hold.
func TestNewAndAcceptRenderTheIndex(t *testing.T) {
	repo, dir := indexRepo(t)
	now := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	res, err := New(repo, noGit{}, Options{Title: "Third", Supersedes: []int{2}, Refines: []int{1}, Body: "## Context\n\nc\n\n## Decision\n\nWe decide again.\n", Now: now})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if want := []string{"design/adrs/0003-third.md", "design/adrs/README.md", "design/adrs/0002-second.md"}; !slices.Equal(res.Changed, want) {
		t.Errorf("changed %q, want %q", res.Changed, want)
	}
	index := indexIsRendered(t, dir)
	want := "# Architecture Decision Records\n\nGenerated.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n" +
		"| [0001](0001-first.md) | First | accepted |\n" +
		"| [0002](0002-second.md) | Second | accepted, superseded by 0003 |\n" +
		"| [0003](0003-third.md) | Third | proposed, supersedes 0002, refines 0001 |\n\nBelow.\n"
	if index != want {
		t.Errorf("after New:\n%s\nwant:\n%s", index, want)
	}

	res, err = Accept(repo, noGit{}, 3, Options{Now: now})
	if err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if want := []string{"design/adrs/0003-third.md", "design/adrs/README.md"}; !slices.Equal(res.Changed, want) {
		t.Errorf("changed %q, want %q", res.Changed, want)
	}
	if index := indexIsRendered(t, dir); !strings.Contains(index, "| [0003](0003-third.md) | Third | accepted, supersedes 0002, refines 0001 |\n") {
		t.Errorf("after Accept:\n%s", index)
	}
}

// An ADR folder with no index gets none from New or Accept.
func TestNewWritesNoIndexWhereThereIsNone(t *testing.T) {
	repo, dir := indexRepo(t)
	if err := os.Remove(filepath.Join(dir, "README.md")); err != nil {
		t.Fatal(err)
	}
	res, err := New(repo, noGit{}, Options{Title: "Third", Body: "## Decision\n\nWe decide again.\n", Now: time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if want := []string{"design/adrs/0003-third.md"}; !slices.Equal(res.Changed, want) {
		t.Errorf("changed %q, want %q", res.Changed, want)
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Errorf("no index is created: %v", err)
	}
}

// New, the supersede it does, Accept, and SetTopics replace each ADR rather
// than write it in place, so the folder's modification time moves and the
// watcher's slow tier sees the change at its next fast look (S-0606). The
// index is adrindex's to write.
func TestWritesReplaceTheFiles(t *testing.T) {
	repo, dir := indexRepo(t)
	now := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	stat := func(name string) os.FileInfo {
		t.Helper()
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		return info
	}
	replaced := func(name string, before os.FileInfo) {
		t.Helper()
		if os.SameFile(before, stat(name)) {
			t.Errorf("%s was written in place, not replaced", name)
		}
	}
	second := stat("0002-second.md")
	if _, err := New(repo, noGit{}, Options{Title: "Third", Supersedes: []int{2}, Body: "## Context\n\nc\n\n## Decision\n\nWe decide again.\n", Now: now}); err != nil {
		t.Fatalf("New: %v", err)
	}
	replaced("0002-second.md", second)

	third := stat("0003-third.md")
	if _, err := Accept(repo, noGit{}, 3, Options{Now: now}); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	replaced("0003-third.md", third)

	third = stat("0003-third.md")
	if _, err := SetTopics(repo, noGit{}, 3, []string{"all"}, Options{Now: now}); err != nil {
		t.Fatalf("SetTopics: %v", err)
	}
	replaced("0003-third.md", third)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("a temporary file was left behind: %s", e.Name())
		}
	}
}
