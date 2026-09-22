package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// checksProject is a project with a real worktree directory for one story
// (RunChecks only needs the directory to exist, not a real git worktree)
// and, on the host's own config, one named check that always passes.
func checksProject(t *testing.T) (root, story string) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	story = "S-0001"
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.WorktreePath(story), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "ok", "--", "true"); code != 0 {
		t.Fatal("set")
	}
	return root, story
}

func TestChecksRunRefusesSomethingThatIsNotAStoryID(t *testing.T) {
	root := tempProject(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	if _, errOut, code := runIn(t, root, "checks", "run", "not-a-story"); code == 0 || !strings.Contains(errOut, "is not a story's ID") {
		t.Errorf("run: %d %s", code, errOut)
	}
}

func TestChecksStatusBeforeAnyRun(t *testing.T) {
	root, story := checksProject(t)
	out, _, code := runIn(t, root, "checks", "status", story, "--json")
	if code != 0 {
		t.Fatalf("status: %s", out)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["running"] != false || got["started"] != "" || got["outcome"] != nil {
		t.Errorf("nothing has run yet: %+v", got)
	}
}

func TestChecksRunEndToEndThroughTheCLI(t *testing.T) {
	root, story := checksProject(t)
	out, errOut, code := runIn(t, root, "checks", "run", story, "--json")
	if code != 0 {
		t.Fatalf("run: %d %s %s", code, out, errOut)
	}
	var run map[string]any
	if err := json.Unmarshal([]byte(out), &run); err != nil {
		t.Fatal(err)
	}
	if run["outcome"] != "passed" || run["running"] != false {
		t.Fatalf("run: %+v", run)
	}
	steps, ok := run["steps"].([]any)
	if !ok || len(steps) != 1 {
		t.Fatalf("steps: %+v", run["steps"])
	}

	// status now shows the same outcome, and tail can read the log from 0
	statusOut, _, code := runIn(t, root, "checks", "status", story, "--json")
	if code != 0 || !strings.Contains(statusOut, `"outcome": "passed"`) {
		t.Errorf("status after: %d %s", code, statusOut)
	}
	tailOut, tailErr, code := runIn(t, root, "checks", "tail", story, "--from", "0", "--wait", "0", "--json")
	if code != 0 {
		t.Fatalf("tail: %d %s %s", code, tailOut, tailErr)
	}
	var tailResult map[string]any
	if err := json.Unmarshal([]byte(tailOut), &tailResult); err != nil {
		t.Fatal(err)
	}
	if tailResult["running"] != false {
		t.Errorf("tail: %+v", tailResult)
	}
	if off, ok := tailResult["offset"].(float64); !ok || off <= 0 {
		t.Errorf("tail should have read something: %+v", tailResult)
	}

	// a second run is fine now that the first has ended
	if _, errOut, code := runIn(t, root, "checks", "run", story); code != 0 {
		t.Errorf("a second run once the first ended: %d %s", code, errOut)
	}
}

func TestChecksRunRefusesWithoutAWorktree(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "ok", "--", "true"); code != 0 {
		t.Fatal("set")
	}
	if _, errOut, code := runIn(t, root, "checks", "run", "S-0099"); code == 0 || !strings.Contains(errOut, "no worktree") || !strings.Contains(errOut, "flai stream open S-0099") {
		t.Errorf("run without a worktree: %d %s", code, errOut)
	}
}

func TestChecksCancelRefusesWhenNothingIsRunning(t *testing.T) {
	root, story := checksProject(t)
	if _, errOut, code := runIn(t, root, "checks", "cancel", story); code == 0 || !strings.Contains(errOut, "no checks run is active") {
		t.Errorf("cancel: %d %s", code, errOut)
	}
}
