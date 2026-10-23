package digest

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newRepo is a project with its design folder at design and nothing in it.
func newRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	write(t, root, "system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// write writes a file and the folders above it.
func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const goodName = "D-0003-the-digest-is-one-file-per-entry.md"

const good = `---
id: D-0003
title: The digest is one file per entry
kind: decision
sources: [S-0428, TH-0476]
topics: [digest]
paths: [flai/internal/digest]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0003 The digest is one file per entry

TH-0476 confirmed one file per entry under design/digest.
`

var env = Env{
	Known:  func(id string) bool { return id != "S-0999" },
	Exists: func(p string) bool { return p == "design/system/digest.md" },
}

func TestDirIsTheDesignFoldersDigest(t *testing.T) {
	m := manifest.Manifest{Layout: map[string]string{"design": "notes"}}
	if got := Dir(m); got != "notes/digest" {
		t.Errorf("Dir = %q, want notes/digest", got)
	}
	if !Skipped("README.md") || Skipped(goodName) {
		t.Error("only the README is skipped")
	}
}

// ADR-0191: an entry with every field there and well formed has no problems,
// active or absorbed by a file that exists.
func TestValidateAcceptsAWellFormedEntry(t *testing.T) {
	e, problems := Validate(goodName, good, env)
	if len(problems) != 0 {
		t.Fatalf("problems: %+v", problems)
	}
	want := Entry{ID: "D-0003", Title: "The digest is one file per entry", Kind: "decision", Sources: []string{"S-0428", "TH-0476"},
		Topics: []string{"digest"}, Paths: []string{"flai/internal/digest"}, Date: "2026-10-23", Status: "active",
		AskedAgain: []string{}, Cited: []string{}, Body: "TH-0476 confirmed one file per entry under design/digest."}
	if !reflect.DeepEqual(e, want) {
		t.Errorf("entry:\n got %+v\nwant %+v", e, want)
	}
	absorbed := strings.NewReplacer("status: active", "status: absorbed", `absorbed_by: ""`, "absorbed_by: design/system/digest.md",
		"asked_again: []", "asked_again: [TH-0500]", "cited: []", "cited: [S-0430, TH-0501]").Replace(good)
	if _, problems := Validate(goodName, absorbed, env); len(problems) != 0 {
		t.Errorf("absorbed: %+v", problems)
	}
	if _, problems := Validate(goodName, strings.Replace(good, "S-0428", "S-0999", 1), Env{}); len(problems) != 0 {
		t.Errorf("an Env without funcs checks no existence: %+v", problems)
	}
}

// ADR-0191: Validate names each way an entry falls short, the field and the
// line of its key, and what to do.
func TestValidateReportsEachProblem(t *testing.T) {
	without := func(key string) string {
		var kept []string
		for _, l := range strings.SplitAfter(good, "\n") {
			if !strings.HasPrefix(l, key+":") {
				kept = append(kept, l)
			}
		}
		return strings.Join(kept, "")
	}
	with := func(pairs ...string) string { return strings.NewReplacer(pairs...).Replace(good) }
	for _, c := range []struct {
		name, file, content string
		field               string
		line                int
		rule, says          string
	}{
		{"name without four digits", "D-3-x.md", good, "", 1, RuleEntry, "D-nnnn-<slug>.md"},
		{"name not a slug", "D-0003-The_Digest.md", good, "", 1, RuleEntry, "D-nnnn-<slug>.md"},
		{"no front matter", goodName, "# D-0003 The digest\n", "", 1, RuleEntry, "no front matter"},
		{"front matter does not parse", goodName, with("sources: [S-0428, TH-0476]", "sources: [S-0428, TH-0476"), "", 0, RuleEntry, "does not parse"},
		{"missing id", goodName, without("id"), "id", 1, RuleEntry, "needs id"},
		{"missing title", goodName, without("title"), "title", 1, RuleEntry, "needs title"},
		{"missing kind", goodName, without("kind"), "kind", 1, RuleEntry, "needs kind"},
		{"missing date", goodName, without("date"), "date", 1, RuleEntry, "needs date"},
		{"missing status", goodName, without("status"), "status", 1, RuleEntry, "needs status"},
		{"missing sources", goodName, without("sources"), "sources", 1, RuleEntry, "at least one"},
		{"empty sources", goodName, with("sources: [S-0428, TH-0476]", "sources: []"), "sources", 5, RuleEntry, "at least one"},
		{"id not an entry ID", goodName, with("id: D-0003", "id: X-3"), "id", 2, RuleEntry, "not an entry ID"},
		{"id not the file's", goodName, with("id: D-0003", "id: D-0004"), "id", 2, RuleEntry, "named for D-0003"},
		{"unknown kind", goodName, with("kind: decision", "kind: rumour"), "kind", 4, RuleEntry, "none of decision, answer, recommendation, cause"},
		{"unknown status", goodName, with("status: active", "status: stale"), "status", 9, RuleEntry, "none of active, absorbed, compacted"},
		{"date not a date", goodName, with("date: 2026-10-23", "date: 23/10/2026"), "date", 8, RuleEntry, "YYYY-MM-DD"},
		{"date out of range", goodName, with("date: 2026-10-23", "date: 2026-13-01"), "date", 8, RuleEntry, "YYYY-MM-DD"},
		{"source not an ID", goodName, with("TH-0476]", "PR-12]"), "sources", 5, RuleEntry, `"PR-12" is not the ID`},
		{"source names nothing", goodName, with("S-0428,", "S-0999,"), "sources", 5, RuleEntry, "S-0999 names nothing"},
		{"topic not a word", goodName, with("topics: [digest]", `topics: ["two words"]`), "topics", 6, RuleEntry, "not a word"},
		{"absorbed_by while active", goodName, with(`absorbed_by: ""`, "absorbed_by: design/system/digest.md"), "absorbed_by", 10, RuleEntry, "set status: absorbed"},
		{"absorbed without absorbed_by", goodName, with("status: active", "status: absorbed"), "absorbed_by", 10, RuleEntry, "absorbed_by is empty"},
		{"absorbed_by names no file", goodName, with("status: active", "status: absorbed", `absorbed_by: ""`, "absorbed_by: design/gone.md"), "absorbed_by", 10, RuleEntry, "names no file"},
		{"absorbed_by outside the project", goodName, with("status: active", "status: absorbed", `absorbed_by: ""`, "absorbed_by: ../elsewhere.md"), "absorbed_by", 10, RuleEntry, "not a path from the project root"},
		{"asked_again not a thread", goodName, with("asked_again: []", "asked_again: [S-0001]"), "asked_again", 11, RuleEntry, "not a thread ID"},
		{"cited an issue", goodName, with("cited: []", "cited: [I-0001]"), "cited", 12, RuleEntry, "not the ID of a work item"},
		{"a secret in the body", goodName, good + "\nThe key is ghp_" + strings.Repeat("a1B2", 9) + ".\n", "", 19, RuleSecret, "never its value"},
		{"a secret in the front matter", goodName, with("paths: [flai/internal/digest]", "paths: [AKIAIOSFODNN7EXAMPLE]"), "", 7, RuleSecret, "AWS access key"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, problems := Validate(c.file, c.content, env)
			if len(problems) != 1 {
				t.Fatalf("want one problem, got %+v", problems)
			}
			p := problems[0]
			if p.Field != c.field || p.Rule != c.rule || !strings.Contains(p.Message, c.says) || (c.line > 0 && p.Line != c.line) {
				t.Errorf("got %+v, want field %q line %d rule %s saying %q", p, c.field, c.line, c.rule, c.says)
			}
		})
	}
}

// ADR-0191: List reads every entry in the folder, sorted by ID, and leaves
// the README out; a missing folder is no entries. Get finds an entry by its
// ID in any padding.
func TestListAndGetReadTheEntries(t *testing.T) {
	r := newRepo(t)
	if got, err := List(r); err != nil || len(got) != 0 {
		t.Fatalf("missing folder: %v, %v", got, err)
	}
	write(t, r.Root, "design/digest/README.md", "# Digest\n")
	write(t, r.Root, "design/digest/"+goodName, good)
	write(t, r.Root, "design/digest/D-0010-later.md", strings.NewReplacer("D-0003", "D-0010", "The digest is one file per entry", "Later").Replace(good))
	write(t, r.Root, "design/digest/D-0002-earlier.md", strings.NewReplacer("D-0003", "D-0002", "The digest is one file per entry", "Earlier").Replace(good))
	got, err := List(r)
	if err != nil {
		t.Fatal(err)
	}
	var ids, paths []string
	for _, e := range got {
		ids, paths = append(ids, e.ID), append(paths, e.Path)
	}
	if want := []string{"D-0002", "D-0003", "D-0010"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ids %v, want %v", ids, want)
	}
	if paths[1] != "design/digest/"+goodName {
		t.Errorf("path %q", paths[1])
	}
	for _, id := range []string{"D-3", "d-0003", "D-003", "D-0003"} {
		e, err := Get(r, id)
		if err != nil || e.ID != "D-0003" || e.Body == "" {
			t.Errorf("Get(%s) = %+v, %v", id, e, err)
		}
	}
	if _, err := Get(r, "D-0004"); err == nil || !strings.Contains(err.Error(), "flai digest list") {
		t.Errorf("an ID no entry has: %v", err)
	}
	if _, err := Get(r, "S-0003"); err == nil || !strings.Contains(err.Error(), "not a digest entry ID") {
		t.Errorf("not an entry ID: %v", err)
	}
	write(t, r.Root, "design/digest/D-0011-broken.md", "no front matter\n")
	if _, err := List(r); err == nil || !strings.Contains(err.Error(), "D-0011-broken.md") {
		t.Errorf("a broken entry is named: %v", err)
	}
}

// ADR-0191: a source is known when a work item, live or archived, a thread,
// or an issue has the ID, typed in any padding; absorbed_by exists when the
// file is under the checkout's root.
func TestEnvForKnowsWorkItemsThreadsAndIssues(t *testing.T) {
	r := newRepo(t)
	for _, rel := range []string{
		"wip/kanban/stories/S-0001-live.md", "wip/archive/kanban/tasks/T-0002-archived.md", "wip/kanban/epics/E-001-legacy.md",
		"wip/threads/TH-0003-thread.md", "design/issues/I-0004-issue.md", "design/system/x.md",
	} {
		write(t, r.Root, rel, "x\n")
	}
	e := EnvFor(r)
	for id, want := range map[string]bool{"S-0001": true, "S-1": true, "T-0002": true, "E-0001": true, "TH-3": true, "I-0004": true, "S-0005": false, "TH-0004": false, "I-0001": false} {
		if got := e.Known(id); got != want {
			t.Errorf("Known(%s) = %v, want %v", id, got, want)
		}
	}
	for p, want := range map[string]bool{"design/system/x.md": true, "design/system/y.md": false, "../x.md": false, "/etc/hostname": false} {
		if got := e.Exists(p); got != want {
			t.Errorf("Exists(%s) = %v, want %v", p, got, want)
		}
	}
}
