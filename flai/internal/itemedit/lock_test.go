package itemedit

import (
	"bytes"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// editedStory notes S-0001 in an undo and saves the edit's change to it, tags
// of its own, as Apply does before its check; it returns the undo and the
// file as it was before the edit.
func editedStory(t *testing.T, repo *workitem.Repo, at time.Time) (*undo, []byte) {
	t.Helper()
	it, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	u := &undo{repo: repo, files: map[string]*file{}}
	u.note(it.Path)
	it.Tags, it.Updated = []string{"edited"}, at.Format(workitem.TimeFormat)
	if err := repo.Save(it); err != nil {
		t.Fatal(err)
	}
	u.saved(it)
	return u, u.files[it.Path].before
}

// I-0175: the check after an edit takes seconds; a move another process
// saved to the item meanwhile, from a copy it read before the edit, was undone
// when the check refused the edit and the files were put back byte for byte.
func TestRestoreKeepsAMoveSavedAfterTheEdit(t *testing.T) {
	repo := planningRepo(t)
	t1 := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	starter, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	u, _ := editedStory(t, repo, t1)
	starter.Status = workitem.InProgress
	starter.Transitions = append(starter.Transitions, workitem.Transition{To: workitem.InProgress, At: t1.Add(time.Minute).Format(workitem.TimeFormat), By: "agent"})
	starter.Updated = t1.Add(time.Minute).Format(workitem.TimeFormat)
	if err := repo.Save(starter); err != nil {
		t.Fatalf("the move after the edit was refused: %v", err)
	}
	if err := u.restore(); err != nil {
		t.Fatalf("restore: %v", err)
	}
	got, err := workitem.ReadItem(starter.Path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 0 {
		t.Errorf("the edit's tags were not taken out: %v", got.Tags)
	}
	if got.Status != workitem.InProgress || !slices.ContainsFunc(got.Transitions, func(tr workitem.Transition) bool { return tr.To == workitem.InProgress }) {
		t.Errorf("restore undid the move made after the edit: status %s, transitions %+v", got.Status, got.Transitions)
	}
}

// With nothing written since the edit, restore gives the file back as it was.
func TestRestoreGivesBackTheFileAsItWas(t *testing.T) {
	repo := planningRepo(t)
	u, before := editedStory(t, repo, time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC))
	if err := u.restore(); err != nil {
		t.Fatalf("restore: %v", err)
	}
	for path := range u.files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, before) {
			t.Errorf("restore wrote\n%s\nwant\n%s", data, before)
		}
	}
}
