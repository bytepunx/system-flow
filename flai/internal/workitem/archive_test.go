package workitem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// I-0139: a narrative archived with two blank lines under ## Open questions,
// as a hand edit can leave them, has one there in the archive, while two
// blank lines inside fenced code and every other line stay as they were.
func TestArchiveCollapsesANarrativesBlankRuns(t *testing.T) {
	r := inProgressStream(t)
	src := r.NarrativePath("S-0001")
	n, err := ReadNarrative(src)
	if err != nil {
		t.Fatal(err)
	}
	fenced := "```text\none\n\n\ntwo\n```"
	n.Body = strings.Replace(n.Body, "## Open questions\n\n- Which?\n\n", "## Open questions\n\n\n", 1)
	n.Body = strings.Replace(n.Body, "Why.\n", "Why.\n\n"+fenced+"\n", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, src)

	s, _ := r.Get("S-0001")
	s.Status = Done
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	items, _ := r.List(false)
	plan, err := r.PlanArchive(items, []string{"S-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Archive(plan); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("the narrative is still in wip/agents: %v", err)
	}
	got := readFile(t, filepath.Join(r.ArchiveDir(), "agents", "S-0001.md"))
	if !strings.Contains(got, "## Open questions\n\n## Log") {
		t.Errorf("want one blank line under ## Open questions:\n%s", got)
	}
	if !strings.Contains(got, fenced) {
		t.Errorf("blank lines inside fenced code must stay:\n%s", got)
	}
	if want := strings.Replace(before, "## Open questions\n\n\n", "## Open questions\n\n", 1); got != want {
		t.Errorf("only the blank run may change:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// Archiving refuses a narrative whose place in the archive is taken, and
// leaves it where it was.
func TestArchiveRefusesANarrativeAlreadyArchived(t *testing.T) {
	r := inProgressStream(t)
	src := r.NarrativePath("S-0001")
	dest := filepath.Join(r.ArchiveDir(), "agents", "S-0001.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("taken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.Archive(&ArchivePlan{Narratives: []string{src}}); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("want a refusal, got %v", err)
	}
	if _, err := os.Stat(src); err != nil {
		t.Errorf("the narrative must stay where it was: %v", err)
	}
	if readFile(t, dest) != "taken\n" {
		t.Error("the archived file must be left as it was")
	}
}
