package itemedit

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// S-0468: Apply refuses a task a folder touch it did not have, writing
// nothing, with the folder, its file count, and what to do instead; it takes
// a file and a path to be made, keeps a folder touch the task already had,
// and leaves a story's folder touch alone (ADR-0188).
func TestApplyRefusesATaskAFolderTouch(t *testing.T) {
	repo := planningRepo(t)
	for _, p := range []string{"flai/cmd/a.go", "flai/cmd/b.go", "flai/cmd/c.go"} {
		full := filepath.Join(repo.Root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	task, err := repo.Get("T-0001")
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(task.Path)
	if err != nil {
		t.Fatal(err)
	}

	_, it, err := edit(t, repo, "T-0001", Change{Title: str("Renamed"), Touches: &[]string{"flai/cmd/a.go", "flai/cmd"}}, "alex", at)
	var inv *InvalidError
	if !errors.As(err, &inv) || !strings.Contains(err.Error(), "T-0001's touch flai/cmd is a folder holding 3 files: a task takes no folder touch (ADR-0188); name the files the task changes, or leave the folder to the story's touches with a ### Planning note") {
		t.Fatalf("a task's folder touch must be refused as the caller's mistake: %v", err)
	}
	if after, _ := os.ReadFile(it.Path); string(after) != string(before) || it.Title != "Task" {
		t.Errorf("nothing is written for a refused edit:\n%s", after)
	}

	if _, it, err = edit(t, repo, "T-0001", Change{Touches: &[]string{"flai/cmd/a.go", "flai/cmd/new.go"}}, "alex", at); err != nil || !slices.Equal(it.Touches, []string{"flai/cmd/a.go", "flai/cmd/new.go"}) {
		t.Errorf("a file and a path to be made are taken: %v %v", it.Touches, err)
	}

	// a task written before the rule keeps its folder touch through an edit
	it.Touches = []string{"flai/cmd"}
	if err := repo.Save(it); err != nil {
		t.Fatal(err)
	}
	if _, it, err = edit(t, repo, "T-0001", Change{Title: str("Kept"), Touches: &[]string{"flai/cmd/", "flai/cmd/a.go"}}, "alex", at); err != nil || it.Title != "Kept" || !slices.Equal(it.Touches, []string{"flai/cmd", "flai/cmd/a.go"}) {
		t.Errorf("a folder touch the task had stays: %v %v", it.Touches, err)
	}

	if _, s, err := edit(t, repo, "S-0001", Change{Touches: &[]string{"flai/cmd"}}, "alex", at); err != nil || !slices.Equal(s.Touches, []string{"flai/cmd"}) {
		t.Errorf("a story keeps a folder touch: %v %v", s.Touches, err)
	}
}
