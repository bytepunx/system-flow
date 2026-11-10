package itemedit

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/filelock"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// checkedWriteFree tries the checked-write lock without waiting long, and
// says whether it was free.
func checkedWriteFree(t *testing.T, repo *workitem.Repo) bool {
	t.Helper()
	unlock, err := filelock.Lock(filepath.Join(repo.CacheDir(), "items-checked.lock"), 20*time.Millisecond)
	if err != nil {
		return false
	}
	unlock()
	return true
}

// I-0328: a checked edit holds the checked-write lock from its first save
// until the check keeps the edit or the restore puts it back, so that a
// reader under it, such as flai serve's look judging holds, never sees an
// edit the check refuses; and it lets the lock go once Apply returns.
func TestTheCheckedWriteLockIsHeldWhileTheEditIsChecked(t *testing.T) {
	at := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name string
		kept bool
	}{
		{"kept", true},
		{"refused", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := planningRepo(t)
			story, err := repo.Get("S-0001")
			if err != nil {
				t.Fatal(err)
			}
			was := snapshot(story.Path, filepath.Join(repo.AgentsDir(), "index.md"))
			var during []bool
			runCheck = func(r *workitem.Repo, now time.Time, w check.Write) (*check.Result, error) {
				during = append(during, checkedWriteFree(t, r))
				res, err := check.RunItemWrite(r, now, w)
				if err != nil || tc.kept || len(during) < 2 {
					return res, err
				}
				// the run with the edit in place refuses it
				res.Findings = append(res.Findings, check.Finding{Level: check.Error, Rule: "item.probe", Path: rel(r, story.Path), Message: "refused by the test"})
				res.Errors++
				return res, nil
			}
			t.Cleanup(func() { runCheck = check.RunItemWrite })
			_, err = Apply(repo, nil, "S-0001", Change{Tags: &[]string{"edited"}}, Options{By: "alex", NoCommit: true, Now: at})
			if tc.kept && err != nil {
				t.Fatalf("the edit is kept: %v", err)
			}
			if !tc.kept {
				if _, ok := docedit.IsRefused(err); !ok {
					t.Fatalf("want the edit refused, got %v", err)
				}
				sameFiles(t, was)
			}
			// the run before the write is free; the run with the edit in place holds it
			if !slices.Equal(during, []bool{true, false}) {
				t.Errorf("the lock was free before the write and held while the edit was checked, got free %v", during)
			}
			if !checkedWriteFree(t, repo) {
				t.Error("the lock is free once Apply returns")
			}
		})
	}
}

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
