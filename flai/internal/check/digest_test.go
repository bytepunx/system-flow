package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// digestEntry is an entry's text with its front matter keys in ADR-0191's
// order, each on its own line from line 2, over replacing a key's value; an
// empty value leaves the key out, an x_ key holding its line so the others
// keep theirs. Its body is on line 17.
func digestEntry(id string, over map[string]string, body string) string {
	keys := []string{"id", "title", "kind", "sources", "topics", "paths", "date", "status", "absorbed_by", "asked_again", "cited"}
	values := map[string]string{
		"id": id, "title": "The fact", "kind": "decision", "sources": "[E-0001]", "topics": "[cli]", "paths": "[]",
		"date": "2026-10-23", "status": "active", "absorbed_by": `""`, "asked_again": "[]", "cited": "[]",
	}
	for k, v := range over {
		values[k] = v
	}
	var b strings.Builder
	b.WriteString("---\n")
	for _, k := range keys {
		if values[k] == "" {
			b.WriteString("x_" + k + ": \"\"\n")
			continue
		}
		b.WriteString(k + ": " + values[k] + "\n")
	}
	b.WriteString("---\n\n# " + id + " The fact\n\n" + body + "\n")
	return b.String()
}

// S-0428: flai check reports digest.entry for an entry whose front matter
// does not parse or misses a field, whose kind or status is unknown, whose
// source names nothing, or whose absorbed_by names no file, and digest.secret
// for text with the shape of a token, each an error on its line; a clean
// entry, the folder's README.md, and files that are not markdown are not
// reported, and the folder is checked by its own rules only.
func TestDigestEntriesAreChecked(t *testing.T) {
	repo := planningProject(t, "")
	mustItem(t, repo, workitem.Epic, "An epic", "")
	dir := filepath.Join(repo.Root, "design", "digest")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	token := "ghp_" + strings.Repeat("a1B2", 9)
	write("README.md", "# The digest\n\nNo front matter, and not an entry.\n")
	write("notes.txt", "not an entry\n")
	write("sub/D-0099-nested.md", "not checked\n")
	write("D-0001-clean.md", digestEntry("D-0001", nil, "What was decided, and why."))
	write("D-0002-absorbed.md", digestEntry("D-0002", map[string]string{"status": "absorbed", "absorbed_by": "system-flow.yaml"}, "Recorded."))
	write("D-0003-parse.md", "---\nid: D-0003\ntitle: [unclosed\n---\n\n# D-0003 The fact\n")
	write("D-0004-missing.md", digestEntry("D-0004", map[string]string{"title": ""}, "No title."))
	write("D-0005-kind.md", digestEntry("D-0005", map[string]string{"kind": "guess"}, "A guess."))
	write("D-0006-status.md", digestEntry("D-0006", map[string]string{"status": "open"}, "Open."))
	write("D-0007-source.md", digestEntry("D-0007", map[string]string{"sources": "[S-0999]"}, "From nowhere."))
	write("D-0008-gone.md", digestEntry("D-0008", map[string]string{"status": "absorbed", "absorbed_by": "design/gone.md"}, "Gone."))
	write("D-0009-secret.md", digestEntry("D-0009", nil, "The token is "+token+" today."))

	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][]Finding{}
	for _, f := range res.Findings {
		if !strings.HasPrefix(filepath.ToSlash(f.Path), "design/digest/") {
			continue
		}
		if !strings.HasPrefix(f.Rule, "digest.") {
			t.Errorf("the folder is checked by its own rules only: %+v", f)
			continue
		}
		got[filepath.Base(f.Path)] = append(got[filepath.Base(f.Path)], f)
	}
	for _, name := range []string{"README.md", "notes.txt", "D-0099-nested.md", "D-0001-clean.md", "D-0002-absorbed.md"} {
		if fs := got[name]; len(fs) != 0 {
			t.Errorf("%s: want no findings, got %+v", name, fs)
		}
	}
	for _, want := range []struct {
		name, rule string
		line       int
		says       string
	}{
		{"D-0003-parse.md", "digest.entry", 0, "front matter does not parse"},
		{"D-0004-missing.md", "digest.entry", 1, "front matter needs title"},
		{"D-0005-kind.md", "digest.entry", 4, `kind is "guess"`},
		{"D-0006-status.md", "digest.entry", 9, `status is "open"`},
		{"D-0007-source.md", "digest.entry", 5, "source S-0999 names nothing in the project or its archive"},
		{"D-0008-gone.md", "digest.entry", 10, "absorbed_by is design/gone.md, which names no file"},
		{"D-0009-secret.md", "digest.secret", 17, "a GitHub token (ghp_… of 40 characters)"},
	} {
		fs := got[want.name]
		if len(fs) != 1 {
			t.Errorf("%s: want one finding, got %+v", want.name, fs)
			continue
		}
		f := fs[0]
		if f.Level != Error || f.Rule != want.rule || (want.line != 0 && f.Line != want.line) || !strings.Contains(f.Message, want.says) {
			t.Errorf("%s: %+v, want a %s error on line %d saying %q", want.name, f, want.rule, want.line, want.says)
		}
		if strings.Contains(f.Message, token) {
			t.Errorf("%s: the finding shows the secret: %s", want.name, f.Message)
		}
	}
}

// S-0428: a project with no digest folder has nothing checked.
func TestNoDigestFolderIsNotChecked(t *testing.T) {
	repo := planningProject(t, "")
	for _, rule := range []string{"digest.entry", "digest.secret"} {
		if got := findings(t, repo, rule); len(got) != 0 {
			t.Errorf("%s with no folder: %+v", rule, got)
		}
	}
}
