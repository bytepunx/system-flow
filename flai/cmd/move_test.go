package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0552 criterion 2: while the merge queue batches and is not held, flai move
// <story> done from review enqueues the story as flai accept does, and lands
// nothing: the story stays in review and the batch under way is not dropped.
func TestMoveToDoneEnqueuesWhileTheQueueBatches(t *testing.T) {
	root, _ := orchestratedQueueProject(t, true)
	mergeQueueOn(t, root)
	queueUnderWay(t, root, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	out, errOut, code := runIn(t, root, "move", "S-0001", "done", "--by", "alex")
	if want := "enqueued S-0001 in the merge queue at position 2 as alex, and it stays in review"; code != 0 || !strings.Contains(out, want) {
		t.Fatalf("flai move done must answer %q: %d %s\n%s", want, code, errOut, out)
	}
	assertNothingAccepted(t, root, mainBefore)
	st := readQueueState(t, root)
	if len(st.Enqueued) != 2 || st.Enqueued[1].Story != "S-0001" || st.Enqueued[1].By != "alex" {
		t.Errorf("enqueued %+v, want S-0001 after S-0005, by alex", st.Enqueued)
	}
	assertBatchUnderWay(t, st)
}

// I-0347: S-0678 was accepted and archived with a blocked interval open.
// flai unblock reaches the archived item and closes the interval at its
// close, as the item was not blocked after it.
func TestUnblockClosesAnArchivedItemsIntervalAtItsClose(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	dir := filepath.Join(root, "wip", "archive", "kanban", "stories")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `---
id: S-0001
type: story
nature: improvement
title: Accepted while blocked
status: done
owner: alex
created: 2026-09-15T08:00:00Z
updated: 2026-09-15T10:15:07Z
transitions:
  - to: ready
    at: 2026-09-15T09:00:00Z
    by: alex
  - to: in-progress
    at: 2026-09-15T09:18:41Z
    by: agent-S-0001
  - to: review
    at: 2026-09-15T10:13:48Z
    by: agent-S-0001
  - to: done
    at: 2026-09-15T10:15:07Z
    by: orchestrator
blocked:
  - from: 2026-09-15T09:26:16Z
    reason: pause for urgent fix
---

# S-0001 Accepted while blocked

## Goal

g

## Acceptance criteria

- [x] c

## Notes
`
	path := filepath.Join(dir, "S-0001-accepted-while-blocked.md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "unblock", "S-0001")
	if want := "S-0001 unblocked at its close, 2026-09-15T10:15:07Z"; code != 0 || !strings.Contains(out, want) {
		t.Fatalf("flai unblock must answer %q: %d %s\n%s", want, code, errOut, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if !it.Archived || it.Path != path {
		t.Errorf("the item must stay in the archive at %s, got %s (archived %v)", path, it.Path, it.Archived)
	}
	if it.IsBlocked() || it.Blocked[0].Until != "2026-09-15T10:15:07Z" {
		t.Errorf("blocked = %+v, want the interval closed at the item's close", it.Blocked)
	}
	if _, errOut, code := runIn(t, root, "unblock", "S-0001"); code == 0 || !strings.Contains(errOut, "not blocked") {
		t.Errorf("a second unblock must say the item is not blocked: %d %s", code, errOut)
	}
}
