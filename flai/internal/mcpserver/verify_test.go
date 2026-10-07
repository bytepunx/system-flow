package mcpserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// verifyFixture is testFixture with the story's worktree on its own branch
// from main and its narrative's Current state and Next steps written; change
// commits a file on the story's branch.
func verifyFixture(t *testing.T) (f *fixture, change func(rel string)) {
	t.Helper()
	f, git := testFixture(t)
	wt := f.repo.WorktreePath(f.story.ID)
	git(f.repo.Root, "worktree", "add", "-q", "-b", "story/"+f.story.ID, wt, "main")
	path := f.repo.NarrativePath(f.story.ID)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	narrative := string(data)
	for _, h := range []string{"## Current state\n", "## Next steps\n"} {
		if !strings.Contains(narrative, h) {
			t.Fatalf("the narrative has no %q: %s", h, narrative)
		}
		narrative = strings.Replace(narrative, h, h+"\nWritten.\n", 1)
	}
	if err := os.WriteFile(path, []byte(narrative), 0o644); err != nil {
		t.Fatal(err)
	}
	return f, func(rel string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(wt, filepath.FromSlash(rel)), []byte("changed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(wt, "add", "-A")
		git(wt, "commit", "-q", "-m", "change "+rel)
	}
}

// steps are the answer's steps as name=state, and its findings' messages.
func steps(t *testing.T, out map[string]any) (states, found string) {
	t.Helper()
	list, ok := out["steps"].([]any)
	if !ok {
		t.Fatalf("no steps in %v", out)
	}
	var s, f []string
	for _, x := range list {
		step := x.(map[string]any)
		s = append(s, fmt.Sprintf("%s=%s", step["name"], step["state"]))
		fs, _ := step["findings"].([]any)
		for _, y := range fs {
			f = append(f, fmt.Sprint(y.(map[string]any)["message"]))
		}
	}
	return strings.Join(s, " "), strings.Join(f, "\n")
}

// S-0270: the tool is listed, saying what it runs, that a finding outside
// the story is a note, and that the verifier is kept for the review; it
// answers a pass, with each step and the tier the branch's change selects,
// and stores the answer as the story's last.
func TestVerifyAnswersAPassAndStoresIt(t *testing.T) {
	f, change := verifyFixture(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	listed := false
	for _, tool := range res.Tools {
		if tool.Name == "verify" {
			listed = true
			for _, want := range []string{"what the close-out runs", "rebase", "narrative", "flai check --strict scoped to the story", "tiers the branch's changes select",
				"outside the story is a note", "records no issue", "flai verify --json", "the verifier is kept for the review", "acceptance criteria and the conventions"} {
				if !strings.Contains(tool.Description, want) {
					t.Errorf("the description does not say %q: %s", want, tool.Description)
				}
			}
		}
	}
	if !listed {
		t.Fatal("the tool verify is not listed")
	}

	change("ok/a.txt")
	out, failed := f.call(t, "verify", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	states, found := steps(t, out)
	if out["passed"] != true || out["stopped_at"] != nil || out["story"] != f.story.ID || out["base"] != "main" ||
		states != "rebase=passed sync=passed narrative=passed check=passed ok=passed" || found != "" {
		t.Errorf("a pass: %v (steps %s, findings %q)", out, states, found)
	}
	last, ok, err := verify.LastReport(f.repo, f.story.ID)
	if err != nil || !ok || !last.Passed || last.Commit != out["commit"] || last.Commit == "" {
		t.Errorf("the stored report: %+v, %v, %v; answered commit %v", last, ok, err, out["commit"])
	}
}

// S-0270: a step that fails is an answer naming it, with its findings, and
// the steps after it not reached; it is stored as the story's last too.
func TestVerifyAnswersTheStepThatFailed(t *testing.T) {
	f, change := verifyFixture(t)
	change("ok/a.txt")
	change("bad/b.txt")
	out, failed := f.call(t, "verify", map[string]any{"story": f.story.ID, "max": 1})
	if failed != "" {
		t.Fatal(failed)
	}
	states, found := steps(t, out)
	if out["passed"] != false || out["stopped_at"] != "bad" ||
		states != "rebase=passed sync=passed narrative=passed check=passed ok=passed bad=failed" || !strings.Contains(found, "bad/b.txt:3: it broke") {
		t.Errorf("a failing tier: %v (steps %s, findings %q)", out, states, found)
	}
	if last, ok, err := verify.LastReport(f.repo, f.story.ID); err != nil || !ok || last.Passed || last.StoppedAt != "bad" {
		t.Errorf("the stored report: %+v, %v, %v", last, ok, err)
	}

	path := f.repo.NarrativePath(f.story.ID)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), "## Next steps\n\nWritten.\n", "## Next steps\n\n-\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	out, failed = f.call(t, "verify", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	states, found = steps(t, out)
	if out["passed"] != false || out["stopped_at"] != "narrative" ||
		states != "rebase=passed sync=passed narrative=failed check=not-reached ok=not-reached bad=not-reached" || !strings.Contains(found, "## Next steps is empty") {
		t.Errorf("an empty narrative section: %v (steps %s, findings %q)", out, states, found)
	}
}

// S-0270: a story that is not there, an item that is not a story, a story
// with no worktree, and a max under 1 are refused, saying what to do.
func TestVerifyRefusesWhatItCannotVerify(t *testing.T) {
	f, _ := verifyFixture(t)
	other, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"story": " "}, "name the story to verify"},
		{map[string]any{"story": "S-0099"}, "S-0099 is not a work item here"},
		{map[string]any{"story": f.task.ID}, f.task.ID + " is a task; name a story"},
		{map[string]any{"story": other.ID}, other.ID + " has no worktree at .flai-cache/worktrees/" + other.ID + "; open it with flai stream open " + other.ID},
		{map[string]any{"story": f.story.ID, "max": -1}, "max -1 is not a number of findings"},
	} {
		if _, failed := f.call(t, "verify", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("%v: %q, want %q", c.args, failed, c.want)
		}
	}
	if _, ok, err := verify.LastReport(f.repo, f.story.ID); ok || err != nil {
		t.Errorf("a refusal stored a report: %v, %v", ok, err)
	}
}
