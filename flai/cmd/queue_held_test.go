package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// S-0402 criterion 4, ADR-0178: while the main checkout holds an
// uncommitted file a story in review changes, flai queue and flai board name
// the story, the file, and the command that frees it, and flai queue --json
// carries them as a dirty_main reason; once the file is clean, neither does.
func TestQueueAndBoardNameTheStoryAnUncommittedFileHoldsBack(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

// S-0549 criterion 1, ADR-0234: flai queue says the preflight holds the
// queue, since when, each file with its findings, and the thread to the
// operator once flai serve has opened it, beside the operator's hold; and
// that the batch's landing waits on main at a story.
func TestQueueShowsThePreflightHoldAndTheLandingWaitingOnMain(t *testing.T) {
	t.Parallel()
	since := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	hold := &mergequeue.PreflightHold{Since: since, Main: "abc", Findings: []verify.FileFinding{
		{Tier: "markdown", Checkout: verify.CheckoutMain, Path: "wip/kanban/tasks/T-0001-x.md", Line: 12, Rule: "MD051", Message: "Link fragments should be valid"},
		{Tier: mergequeue.PreflightCheckTier, Checkout: verify.CheckoutMain, Path: "system-flow.yaml", Rule: "manifest.schema", Message: "name is required"},
	}}
	q := mergequeue.Queue{Mode: "local", Held: true, HeldBy: "alex", HeldAt: since, Preflight: hold,
		Batch: &mergequeue.Batch{Stories: []string{"S-0001", "S-0002"}, Step: mergequeue.StepVerified, Result: &mergequeue.Result{Passed: true, FinishedAt: since},
			Stalled: &mergequeue.Stalled{Story: "S-0002", Reason: "flai accept S-0002: exit status 1", At: since}}}
	var out bytes.Buffer
	printQueue(&out, q)
	want := "held by alex at 2026-11-01T09:00:00Z: no batch is built until flai queue resume\n" +
		"held by the preflight since 2026-11-01T09:00:00Z: no batch is built until the main checkout is clean of these, and the queue resumes on its own\n" +
		"  wip/kanban/tasks/T-0001-x.md:12 MD051: Link fragments should be valid (markdown)\n" +
		"  system-flow.yaml manifest.schema: name is required (flai-check)\n"
	if !strings.HasPrefix(out.String(), want) {
		t.Errorf("flai queue must start with %q:\n%s", want, out.String())
	}
	if want := "  landing waits on main at S-0002 since 2026-11-01T09:00:00Z: flai accept S-0002: exit status 1\n"; !strings.Contains(out.String(), want) {
		t.Errorf("flai queue must say %q:\n%s", want, out.String())
	}
	hold.Thread = "TH-0007"
	out.Reset()
	printQueue(&out, mergequeue.Queue{Mode: "local", Preflight: hold})
	if want := "clean of these, and the queue resumes on its own; thread TH-0007\n"; !strings.Contains(out.String(), want) || strings.Contains(out.String(), "flai queue resume") {
		t.Errorf("flai queue must name the thread, and no operator's hold, with %q:\n%s", want, out.String())
	}
}
