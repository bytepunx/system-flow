package hostapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

func withDocs(t *testing.T) channel.Project {
	t.Helper()
	p := harbour(t)
	files := map[string]string{
		"design/system/overview.md":    "---\ntitle: Overview\nupdated: 2026-09-20\nstatus: active\ntags: [a, b]\n---\n\n# Overview\n\ntext\n",
		"design/system/plain.md":       "# No front matter\n\nbody\n",
		"design/system/.hidden/x.md":   "---\ntitle: hidden\n---\n",
		"design/system/diagram.png":    "not markdown",
		"design/adrs/0000-template.md": "---\nid: ADR-0000\ntitle: Template\n---\n",
		"design/adrs/0001-first.md":    "---\nid: ADR-0001\ntitle: First\nstatus: accepted\ndate: 2026-09-01\nsupersedes: []\nsuperseded_by: [ADR-0002]\n---\n\n# ADR-0001\n",
		"design/adrs/0002-second.md":   "---\nid: ADR-0002\ntitle: Second\nstatus: accepted\ndate: 2026-09-02\nsupersedes: [ADR-0001]\nrefines: [ADR-0001]\n---\n",
		"design/adrs/README.md":        "# index\n",
		"docs/users/index.md":          "---\ntitle: Users\n---\n",
		"secret.md":                    "---\ntitle: not served\n---\n",
		"src/notes.md":                 "# not a document of the project\n",
	}
	for rel, content := range files {
		full := filepath.Join(p.Root, rel)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func find(nodes []*DocNode, path string) *DocNode {
	for _, n := range nodes {
		if n.Path == path {
			return n
		}
		if got := find(n.Children, path); got != nil {
			return got
		}
	}
	return nil
}

func TestDocsTree(t *testing.T) {
	var roots []*DocNode
	if err := call(t, withDocs(t), "docs.tree", `{}`, &roots); err != nil {
		t.Fatal(err)
	}
	if len(roots) != 3 || roots[0].Path != "design" || roots[1].Path != "docs" || roots[2].Path != "wip" {
		t.Fatalf("roots: %+v", roots)
	}
	ov := find(roots, "design/system/overview.md")
	if ov == nil || ov.Kind != "file" || ov.Title != "Overview" || ov.FrontMatter["updated"] != "2026-09-20" || ov.FrontMatter["status"] != "active" {
		t.Errorf("overview: %+v", ov)
	}
	if tags, ok := ov.FrontMatter["tags"].([]any); !ok || len(tags) != 2 {
		t.Errorf("tags: %#v", ov.FrontMatter["tags"])
	}
	if n := find(roots, "design/system/plain.md"); n == nil || n.Title != "" || n.FrontMatter != nil {
		t.Errorf("a file without front matter: %+v", n)
	}
	for _, gone := range []string{"design/system/.hidden/x.md", "design/system/diagram.png", "secret.md", "src/notes.md"} {
		if find(roots, gone) != nil {
			t.Errorf("%s is in the tree", gone)
		}
	}
	if n := find(roots, "wip/kanban/stories"); n == nil || n.Kind != "dir" || len(n.Children) == 0 {
		t.Errorf("wip stories: %+v", n)
	}
}

func TestDocGet(t *testing.T) {
	p := withDocs(t)
	var d Doc
	if err := call(t, p, "doc.get", `{"path":"design/system/overview.md"}`, &d); err != nil {
		t.Fatal(err)
	}
	if d.Path != "design/system/overview.md" || d.FrontMatter["title"] != "Overview" || !strings.HasPrefix(d.Body, "\n# Overview") || !strings.HasPrefix(d.Raw, "---\ntitle: Overview") {
		t.Errorf("%+v", d)
	}
	_ = call(t, p, "doc.get", `{"path":"design/system/plain.md"}`, &d)
	if d.FrontMatter != nil || d.Body != d.Raw {
		t.Errorf("no front matter: %+v", d)
	}
	outside := filepath.Join(t.TempDir(), "outside.md")
	_ = os.WriteFile(outside, []byte("secret\n"), 0o644)
	_ = os.Symlink(outside, filepath.Join(p.Root, "design/system/link.md"))
	refused := map[string]int{
		`{"path":"../outside.md"}`:                 channel.CodeInvalidParams,
		`{"path":"design/../../outside.md"}`:       channel.CodeInvalidParams,
		`{"path":"/etc/passwd"}`:                   channel.CodeInvalidParams,
		`{"path":"design/system/diagram.png"}`:     channel.CodeInvalidParams,
		`{"path":"secret.md"}`:                     channel.CodeInvalidParams,
		`{"path":"src/notes.md"}`:                  channel.CodeInvalidParams,
		`{"path":"system-flow.yaml"}`:              channel.CodeInvalidParams,
		`{"path":"design/system/link.md"}`:         channel.CodeInvalidParams,
		`{"path":""}`:                              channel.CodeInvalidParams,
		`{"path":"design/system/missing.md"}`:      NotFound,
		`{"path":".git/config"}`:                   channel.CodeInvalidParams,
		`{"path":"design/system/../../secret.md"}`: channel.CodeInvalidParams,
	}
	for params, code := range refused {
		if err := call(t, p, "doc.get", params, &d); err == nil || err.Code != code {
			t.Errorf("%s: %+v, want code %d", params, err, code)
		}
	}
}

func TestAdrsList(t *testing.T) {
	var list []Adr
	if err := call(t, withDocs(t), "adrs.list", `{}`, &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != "ADR-0001" || list[0].Date != "2026-09-01" || list[0].SupersededBy[0] != "ADR-0002" ||
		list[1].Supersedes[0] != "ADR-0001" || list[1].Refines[0] != "ADR-0001" || list[1].Path != "design/adrs/0002-second.md" || len(list[0].Refines) != 0 {
		t.Errorf("%+v", list)
	}
}
