package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// moveRefusal runs flai move id done --json in root and answers its exit
// code, what it printed, and the refusal in it.
func moveRefusal(t *testing.T, root, id string) (int, string, unfixedRefusal) {
	t.Helper()
	out, errOut, code := runIn(t, root, "move", id, "done", "--json")
	var got struct {
		Refused unfixedRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("move --json printed no refusal: %v\n%s\n%s", err, out, errOut)
	}
	return code, out, got.Refused
}

// I-0145: with --json, a move to done of a story whose branch holds lint the
// fix leaves is refused as flai accept --json refuses it: exit 4 and the same
// refused object, each finding naming the story's branch as its checkout.
func TestMoveToDoneJSONRefusesABranchsLintAsAcceptDoes(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"docs/**\"]\n    fix_command: [\"true\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/bad.md", "fine\nBAD line\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")

	code, moved, r := moveRefusal(t, root, "S-0001")
	if code != exitDocRefused {
		t.Errorf("exit %d, want %d", code, exitDocRefused)
	}
	if r.ID != "S-0001" || len(r.Findings) != 1 || r.Findings[0].Checkout != "story" || r.Findings[0].Story != "S-0001" ||
		r.Findings[0].Path != "docs/bad.md" || r.Findings[0].Line != 2 || r.Findings[0].Rule != "MD009" {
		t.Errorf("the refusal's data: %+v", r)
	}
	// the same refusal after the move also says the move left S-0001 in review
	accepted, _, _ := runIn(t, root, "--json", "accept", "S-0001")
	if moved != accepted {
		t.Errorf("move printed\n%s\nflai accept printed\n%s", moved, accepted)
	}
}

// I-0145: with --json, a move to done refused for an uncommitted main-checkout
// file whose lint the fix leaves exits 4 with the same refused object as flai
// accept --json, naming the main checkout; without --json it is the text it
// was, with exit 1.
func TestMoveToDoneJSONRefusesTheMainCheckoutsLintAsAcceptDoes(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"wip/notes/**\"]\n    fix_command: [\"true\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "new\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	if err := os.MkdirAll(filepath.Join(root, "wip/notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/notes/n.md"), []byte("BAD\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, moved, r := moveRefusal(t, root, "S-0001")
	if code != exitDocRefused {
		t.Errorf("exit %d, want %d", code, exitDocRefused)
	}
	if len(r.Findings) != 1 || r.Findings[0].Checkout != "main" || r.Findings[0].Story != "" ||
		r.Findings[0].Path != "wip/notes/n.md" || r.Findings[0].Line != 1 || r.Findings[0].Rule != "MD009" {
		t.Errorf("the refusal's findings: %+v", r.Findings)
	}
	accepted, _, _ := runIn(t, root, "--json", "accept", "S-0001")
	if moved != accepted {
		t.Errorf("move printed\n%s\nflai accept printed\n%s", moved, accepted)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "done"); code != 1 || !strings.Contains(errOut, "wip/notes/n.md:1:3 error MD009") {
		t.Errorf("without --json: exit %d\n%s", code, errOut)
	}
}

// I-0145, ADR-0155: with --json, a move to done of a branch that undoes the
// main branch's change exits 4 with the same refused object as flai accept
// --json, the path under undoes_main, and changes nothing.
func TestMoveToDoneJSONRefusesABranchThatUndoesMainAsAcceptDoes(t *testing.T) {
	root, wt := undoingStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	moved, _, code := runIn(t, root, "move", "S-0001", "done", "--json")
	var refused struct {
		Refused undoneRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(moved), &refused); err != nil || code != exitDocRefused {
		t.Fatalf("--json must exit %d with the refusal as data: %d %v %s", exitDocRefused, code, err, moved)
	}
	if r := refused.Refused; r.ID != "S-0001" || len(r.UndoesMain) != 1 || r.UndoesMain[0].Path != "docs/guide.md" {
		t.Errorf("refusal: %+v", r)
	}
	accepted, _, _ := runIn(t, root, "accept", "S-0001", "--json")
	if moved != accepted {
		t.Errorf("move printed\n%s\nflai accept printed\n%s", moved, accepted)
	}
	unchanged(t, root, wt, mainBefore)
}
