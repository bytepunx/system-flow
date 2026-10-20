package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// S-0402 criterion 4, ADR-0178: while the main checkout holds an
// uncommitted file a story in review changes, flai queue and flai board name
// the story, the file, and the command that frees it, and flai queue --json
// carries them as a dirty_main reason; once the file is clean, neither does.
func TestQueueAndBoardNameTheStoryAnUncommittedFileHoldsBack(t *testing.T) {
	root, head := orchestratedStoryInReview(t, false)
	if err := verify.SaveReport(openProject(t, root), verify.Report{Story: "S-0001", Commit: head, Passed: true}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cli/feature.go"), []byte("the operator's own\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const free = "git stash push -- cli/feature.go"

	out, errOut, code := runIn(t, root, "queue")
	if code != 0 || !strings.Contains(out, "dirty_main: S-0001's change would overwrite cli/feature.go, which the main checkout has not committed") || !strings.Contains(out, free) {
		t.Errorf("flai queue must name S-0001, the file, and %q: %d %s\n%s", free, code, errOut, out)
	}
	out, errOut, code = runIn(t, root, "board")
	if want := "held back: S-0001 by the main checkout's uncommitted cli/feature.go; free it: " + free; code != 0 || !strings.Contains(out, want) {
		t.Errorf("flai board must say %q under review: %d %s\n%s", want, code, errOut, out)
	}
	out, _, _ = runIn(t, root, "queue", "--json")
	var q mergequeue.Queue
	if err := json.Unmarshal([]byte(out), &q); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(q.Ineligible) != 1 || !slices.ContainsFunc(q.Ineligible[0].Reasons, func(r mergequeue.Reason) bool {
		return r.Kind == mergequeue.KindDirtyMain && slices.Equal(r.Paths, []string{"cli/feature.go"}) && r.Command == free
	}) {
		t.Errorf("flai queue --json must carry the dirty_main reason with the file and the command: %+v", q.Ineligible)
	}

	if err := os.Remove(filepath.Join(root, "cli/feature.go")); err != nil {
		t.Fatal(err)
	}
	if out, _, _ := runIn(t, root, "queue"); strings.Contains(out, "dirty_main") {
		t.Errorf("once the file is clean flai queue names no hold:\n%s", out)
	}
	if out, _, _ := runIn(t, root, "board"); strings.Contains(out, "held back") {
		t.Errorf("once the file is clean flai board names no hold:\n%s", out)
	}
}

// S-0402 criterion 3, ADR-0178: flai queue names each story the batch's
// landing held back, with the files and the command that frees it.
func TestQueueNamesTheStoriesALandingHeldBack(t *testing.T) {
	root := tempProject(t)
	writeQueueState(t, root, mergequeue.State{Batch: &mergequeue.Batch{
		Stories: []string{"S-0001", "S-0002"}, Step: mergequeue.StepLanding + " S-0002",
		Landing: &mergequeue.Landing{Stories: []string{"S-0001", "S-0002"}, Landed: []string{}, HeldBack: []mergequeue.LeftOut{{Story: "S-0001", Paths: []string{"a.txt", "b c.md"}}}},
	}})
	out, errOut, code := runIn(t, root, "queue")
	if want := "  held back: S-0001, by the main checkout's uncommitted a.txt, b c.md; free it: git stash push -- a.txt 'b c.md'\n"; code != 0 || !strings.Contains(out, want) {
		t.Errorf("flai queue must say %q: %d %s\n%s", want, code, errOut, out)
	}
}
