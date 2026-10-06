package hostapi

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	if ov := find(roots, "design/system/overview.md"); ov == nil || ov.Kind != "file" || ov.Title != "Overview" || ov.Name != "overview.md" {
		t.Errorf("overview: %+v", ov)
	}
	if n := find(roots, "design/system/plain.md"); n == nil || n.Title != "" {
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

// TestDocsTreeKeepsTitles: a warm walk reads no file; a file that changes is
// read again, and one added or removed is in the next answer or not (S-0162).
func TestDocsTreeKeepsTitles(t *testing.T) {
	p := withDocs(t)
	past := time.Now().Add(-time.Hour)
	_ = filepath.WalkDir(p.Root, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			_ = os.Chtimes(path, past, past)
		}
		return nil
	})
	reads := func() int {
		tt := titlesFor(p.Root)
		tt.mu.Lock()
		defer tt.mu.Unlock()
		return tt.reads
	}
	var roots []*DocNode
	if err := call(t, p, "docs.tree", `{}`, &roots); err != nil {
		t.Fatal(err)
	}
	cold := reads()
	if cold == 0 {
		t.Fatal("the first walk read no file")
	}
	if err := call(t, p, "docs.tree", `{}`, &roots); err != nil {
		t.Fatal(err)
	}
	if got := reads(); got != cold {
		t.Errorf("a warm walk read %d files", got-cold)
	}
	if ov := find(roots, "design/system/overview.md"); ov == nil || ov.Title != "Overview" {
		t.Errorf("a kept title: %+v", ov)
	}

	ov := filepath.Join(p.Root, "design/system/overview.md")
	if err := os.WriteFile(ov, []byte("---\ntitle: Overview, retitled\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(p.Root, "docs/users/index.md"))
	_ = os.WriteFile(filepath.Join(p.Root, "docs/users/new.md"), []byte("---\ntitle: New\n---\n"), 0o644)
	if err := call(t, p, "docs.tree", `{}`, &roots); err != nil {
		t.Fatal(err)
	}
	if got := reads() - cold; got != 2 {
		t.Errorf("read %d files after one changed and one was added, want 2", got)
	}
	if n := find(roots, "design/system/overview.md"); n == nil || n.Title != "Overview, retitled" {
		t.Errorf("a changed title: %+v", n)
	}
	if find(roots, "docs/users/index.md") != nil {
		t.Error("a removed file is in the tree")
	}
	if n := find(roots, "docs/users/new.md"); n == nil || n.Title != "New" {
		t.Errorf("an added file: %+v", n)
	}
	tt := titlesFor(p.Root)
	tt.mu.Lock()
	_, kept := tt.files[filepath.Join(p.Root, "docs/users/index.md")]
	tt.mu.Unlock()
	if kept {
		t.Error("a removed file's title is still kept")
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

// TestDocsAnalysisReports: the analyzer's reports under design/analysis are
// in the tree beside the folder's README, as any design folder is, and a
// report is served with its focus and window (S-0223).
func TestDocsAnalysisReports(t *testing.T) {
	p := withDocs(t)
	report := "---\ntitle: Bottlenecks in September\nupdated: 2026-10-06T08:00:00Z\nstatus: active\nfocus: bottlenecks\nfrom: 2026-09-01\nto: 2026-09-30\n---\n\n# Bottlenecks in September\n"
	for rel, content := range map[string]string{
		"design/analysis/README.md":                 "---\ntitle: analysis\n---\n\n# analysis\n",
		"design/analysis/2026-10-06-bottlenecks.md": report,
		"design/experiments/S-0001-first.md":        "---\ntitle: First experiment\n---\n",
		"design/analysis/.draft/2026-10-07-risk.md": "---\ntitle: hidden\n---\n",
	} {
		full := filepath.Join(p.Root, rel)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var roots []*DocNode
	if err := call(t, p, "docs.tree", `{}`, &roots); err != nil {
		t.Fatal(err)
	}
	analysis, experiments := find(roots, "design/analysis"), find(roots, "design/experiments")
	if analysis == nil || analysis.Kind != "dir" || experiments == nil || experiments.Kind != "dir" {
		t.Fatalf("analysis %+v, experiments %+v", analysis, experiments)
	}
	var names []string
	for _, n := range analysis.Children {
		names = append(names, n.Name)
	}
	if strings.Join(names, ",") != "2026-10-06-bottlenecks.md,README.md" {
		t.Errorf("design/analysis lists %v", names)
	}
	if n := find(roots, "design/analysis/2026-10-06-bottlenecks.md"); n == nil || n.Kind != "file" || n.Title != "Bottlenecks in September" {
		t.Errorf("the report: %+v", n)
	}
	if n := find(roots, "design/analysis/README.md"); n == nil || n.Title != "analysis" {
		t.Errorf("the README: %+v", n)
	}
	var d Doc
	if err := call(t, p, "doc.get", `{"path":"design/analysis/2026-10-06-bottlenecks.md"}`, &d); err != nil {
		t.Fatal(err)
	}
	fm := d.FrontMatter
	if fm["title"] != "Bottlenecks in September" || fm["status"] != "active" || fm["focus"] != "bottlenecks" ||
		fm["from"] != "2026-09-01" || fm["to"] != "2026-09-30" || fm["updated"] != "2026-10-06T08:00:00Z" ||
		!strings.HasPrefix(d.Body, "\n# Bottlenecks in September") {
		t.Errorf("the report served: %+v", d)
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
