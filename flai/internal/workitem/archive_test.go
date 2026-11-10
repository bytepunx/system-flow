package workitem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
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

// I-0323: a Save of a copy read before an archive, waiting on the items lock
// while the archive holds it, reads after the move and is refused, writing
// nothing at the item's old path.
func TestArchiveMovesUnderTheItemsLockSoASaveWaitingIsRefused(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	c := copies(t, r, 2)
	stale, archived := c[0], c[1]
	old := stale.Path
	stale.Tags = []string{"usage"}
	orig := lockFile
	t.Cleanup(func() { lockFile = orig })
	var calls atomic.Int32
	waiting := make(chan struct{})
	saved := make(chan error, 1)
	lockFile = func(path string, wait time.Duration) (func(), error) {
		switch calls.Add(1) {
		case 1: // the archive's, held until the save is waiting on it
			unlock, err := orig(path, wait)
			if err != nil {
				return nil, err
			}
			go func() { saved <- r.Save(stale) }()
			<-waiting
			return unlock, nil
		case 2: // the save's
			close(waiting)
		}
		return orig(path, wait)
	}
	if err := r.Archive(&ArchivePlan{Items: []*Item{archived}}); err != nil {
		t.Fatal(err)
	}
	err := <-saved
	if !errors.Is(err, ErrItemMoved) {
		t.Errorf("the save waiting on the archive: %v, want ErrItemMoved", err)
	}
	if _, serr := os.Stat(old); !os.IsNotExist(serr) {
		t.Errorf("the save made the item again at its old path %s", old)
	}
	if got := onDisk(t, archived.Path); len(got.Tags) != 0 {
		t.Errorf("the refused save wrote the archived file: tags %v", got.Tags)
	}
}
