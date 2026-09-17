package workitem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCanonicalID(t *testing.T) {
	for in, want := range map[string]string{
		"S-32": "S-0032", "s-032": "S-0032", "S-0032": "S-0032", "T-12345": "T-12345",
		"E-1": "E-0001", "ADR-0003": "ADR-0003", "I-011": "I-011", "": "",
	} {
		if got := CanonicalID(in); got != want {
			t.Errorf("CanonicalID(%q) = %q, want %q", in, got, want)
		}
	}
	if got := widen("see S-032, T-100 and T-0100, ADR-0003, I-011, S-1."); got != "see S-0032, T-0100 and T-0100, ADR-0003, I-011, S-0001." {
		t.Errorf("widen: %q", got)
	}
}

func TestIDMigration(t *testing.T) {
	r := newProject(t)
	t0 := time.Date(2026, 9, 17, 8, 0, 0, 0, time.UTC)
	// A legacy three-digit epic, story, and task written by hand, plus a narrative.
	write := func(rel, body string) {
		p := filepath.Join(r.Root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("wip/kanban/epics/E-001-big.md", "---\nid: E-001\ntype: epic\nnature: feature\ntitle: Big\nstatus: backlog\nowner: a\ncreated: 2026-09-17T08:00:00Z\nupdated: 2026-09-17T08:00:00Z\ntransitions: []\ntags: []\n---\n\n# E-001 Big\n\n## Goal\nx\n")
	write("wip/kanban/stories/S-007-slice.md", "---\nid: S-007\ntype: story\nnature: feature\ntitle: Slice\nstatus: backlog\nparent: E-001\nowner: a\ncreated: 2026-09-17T08:00:00Z\nupdated: 2026-09-17T08:00:00Z\ntransitions: []\ntags: []\n---\n\n# S-007 Slice\n\n## Goal\nPart of E-001.\n")
	write("wip/archive/kanban/tasks/T-099-old.md", "---\nid: T-099\ntype: task\nnature: feature\ntitle: Old\nstatus: done\nparent: S-007\nowner: a\ncreated: 2026-09-17T08:00:00Z\nupdated: 2026-09-17T08:00:00Z\ntransitions: []\ntags: []\n---\n\n# T-099 Old\n")
	write("wip/agents/S-007.md", "---\nstream: S-007\nagent: a\nsession: s\nupdated: 2026-09-17T08:00:00Z\n---\n\n# S-007 Slice\n\n## Log\n")
	// Half-applied state: front matter already widened, file name not yet.
	write("wip/kanban/tasks/T-005-half.md", "---\nid: T-0005\ntype: task\nnature: feature\ntitle: Half\nstatus: backlog\nparent: S-0007\nowner: a\ncreated: 2026-09-17T08:00:00Z\nupdated: 2026-09-17T08:00:00Z\ntransitions: []\ntags: []\n---\n\n# T-0005 Half\n")
	write("design/system/notes.md", "---\ntitle: Notes\n---\n\nS-007 depends on E-001; ADR-0003 and I-011 stay.\n")
	write("README.md", "See S-007.\n")
	write("flai/internal/testdata/fixture.md", "S-007 must stay\n")
	// A four-digit item created by flai is untouched.
	if it := mustCreate(t, r, Story, "New", "E-001"); it.ID != "S-0008" {
		t.Fatalf("new ID should be four digits above the highest existing number, got %s", it.ID)
	}
	if got, err := r.Get("S-7"); err != nil || got.ID != "S-007" {
		t.Fatalf("short lookup before migration: %v %v", got, err)
	}

	plan, err := r.PlanIDMigration()
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Map) != 4 || plan.Map["S-007"] != "S-0007" || plan.Map["T-099"] != "T-0099" || plan.Map["T-005"] != "T-0005" {
		t.Errorf("map: %v", plan.Map)
	}
	froms := []string{}
	for _, mv := range plan.Renames {
		froms = append(froms, mv.From+">"+mv.To)
	}
	joined := strings.Join(froms, " ")
	for _, want := range []string{
		"wip/kanban/epics/E-001-big.md>wip/kanban/epics/E-0001-big.md",
		"wip/kanban/stories/S-007-slice.md>wip/kanban/stories/S-0007-slice.md",
		"wip/archive/kanban/tasks/T-099-old.md>wip/archive/kanban/tasks/T-0099-old.md",
		"wip/agents/S-007.md>wip/agents/S-0007.md",
		"wip/kanban/tasks/T-005-half.md>wip/kanban/tasks/T-0005-half.md",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing rename %s in %s", want, joined)
		}
	}
	rw := strings.Join(plan.Rewrites, " ")
	for _, want := range []string{"design/system/notes.md", "README.md", "wip/kanban/stories/S-007-slice.md", "wip/kanban/stories/S-0008-new.md", "wip/agents/S-007.md"} {
		if !strings.Contains(rw, want) {
			t.Errorf("missing rewrite %s in %s", want, rw)
		}
	}
	if strings.Contains(rw, "testdata") {
		t.Errorf("testdata must not be rewritten: %s", rw)
	}

	if err := r.ApplyIDMigration(plan, os.Rename); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(r.Root, "wip/kanban/stories/S-0007-slice.md")); err != nil {
		t.Fatal("story not renamed")
	}
	notes, _ := os.ReadFile(filepath.Join(r.Root, "design/system/notes.md"))
	if string(notes) != "---\ntitle: Notes\n---\n\nS-0007 depends on E-0001; ADR-0003 and I-011 stay.\n" {
		t.Errorf("notes: %s", notes)
	}
	it, err := r.Get("S-007")
	if err != nil || it.ID != "S-0007" || it.Parent != "E-0001" {
		t.Errorf("after migration: %+v %v", it, err)
	}
	items, _ := r.List(true)
	if len(items) != 5 {
		t.Errorf("expected 5 items, got %d", len(items))
	}
	n, err := r.LogStream("S-7", "after", StreamOptions{Now: t0})
	if err != nil || n.Stream != "S-0007" {
		t.Errorf("log by short id after migration: %+v %v", n, err)
	}
	again, _ := r.PlanIDMigration()
	if len(again.Map) != 0 || len(again.Rewrites) != 0 {
		t.Errorf("second run should be a no-op: %+v", again)
	}
}
