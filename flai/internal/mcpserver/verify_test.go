package mcpserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// verifyFixture is testFixture with the story's worktree on its own branch
// from main and its narrative's Current state and Next steps written; change
// commits a file on the story's branch.
func verifyFixture(t *testing.T) (f *fixture, change func(rel string)) {
	t.Helper()
	f, change, _ = verifyFixtureWith(t, nil)
	return f, change
}

// verifyFixtureWith is verifyFixture with the server's options changed by
// with as well, when it is not nil, and the git it runs.
func verifyFixtureWith(t *testing.T, with func(*Options)) (f *fixture, change func(rel string), git func(dir string, args ...string)) {
	t.Helper()
	f, git = testFixtureWith(t, with)
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
	}, git
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

// S-0311: the orchestrator's server answers verify, and runs the tiers under
// the verify role rather than its own, so a tier whose tests make flai's
// writes is not refused as the orchestrator (TH-0260).
func TestVerifyRunsTheTiersUnderTheVerifyRoleForTheOrchestrator(t *testing.T) {
	// the server reads its role when the fixture makes it
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	f, change := verifyFixture(t)
	// on the story's branch the ok tier passes only under the verify role
	manifest := strings.Replace(testManifest, `[sh, -c, "exit 0"]`, `[sh, -c, 'test "$FLAI_ROLE" = verify || { echo "FLAI_ROLE is $FLAI_ROLE, not verify"; exit 1; }']`, 1)
	if err := os.WriteFile(filepath.Join(f.repo.WorktreePath(f.story.ID), "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	change("ok/a.txt")

	out, failed := f.call(t, "verify", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatalf("the orchestrator's verify was refused: %s", failed)
	}
	states, found := steps(t, out)
	if out["passed"] != true || states != "rebase=passed sync=passed narrative=passed check=passed ok=passed" {
		t.Errorf("the orchestrator's run: %v (steps %s, findings %q)", out, states, found)
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

// I-0141, ADR-0158: verify refuses tiers its flai does not take, in the
// story's worktree, as the tool test does: naming both versions and the
// upgrade when the project's history holds a newer flai release, and the
// manifest's fix when it does not. No report is stored either way.
func TestVerifySaysTheFlaiIsOlderThanTheProjectWhenTheTiersRefuse(t *testing.T) {
	for _, c := range []struct {
		version string
		want    []string
		not     string
	}{
		{"1.40.1", []string{i0141Problem, "this flai is 1.40.1, older than flai 1.49.0 in the project's history", "upgrade it with flai host upgrade"}, "flai manifest set tests="},
		{"1.49.0", []string{i0141Problem, "flai manifest set tests="}, "upgrade it"},
	} {
		f, _, git := verifyFixtureWith(t, func(o *Options) { o.Version = c.version })
		git(f.repo.Root, "tag", "-a", "flai/v1.49.0", "-m", "flai 1.49.0", "main")
		if err := os.WriteFile(filepath.Join(f.repo.WorktreePath(f.story.ID), "system-flow.yaml"), []byte(i0141Manifest), 0o644); err != nil {
			t.Fatal(err)
		}
		_, failed := f.call(t, "verify", map[string]any{"story": f.story.ID})
		for _, want := range c.want {
			if !strings.Contains(failed, want) {
				t.Errorf("flai %s: %q, want it to say %q", c.version, failed, want)
			}
		}
		if strings.Contains(failed, c.not) {
			t.Errorf("flai %s: %q, want it not to say %q", c.version, failed, c.not)
		}
		if _, ok, err := verify.LastReport(f.repo, f.story.ID); ok || err != nil {
			t.Errorf("flai %s: a refusal stored a report: %v, %v", c.version, ok, err)
		}
	}
}

// ADR-0174: while the host says the merge queue can batch for the project,
// asked with its main checkout's root at each call, verify leaves the tiers
// marked queue to it as queued steps and passes; with all, or while the queue
// cannot batch, it runs them. A run that runs them reuses the tiers a run
// that left them passed at the same patch, and runs only those it left
// (ADR-0221), and two that run them alike reuse them all.
func TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch(t *testing.T) {
	batches, asked := true, ""
	f, change, _ := verifyFixtureWith(t, func(o *Options) {
		o.QueueBatches = func(root string) bool {
			asked = root
			return batches
		}
	})
	manifest := testManifest + "  - name: late\n    command: [sh, -c, \"exit 0\"]\n    paths: [\"ok/**\"]\n    queue: true\n"
	if err := os.WriteFile(filepath.Join(f.repo.WorktreePath(f.story.ID), "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	change("ok/a.txt")
	for _, c := range []struct {
		name    string
		batches bool
		args    map[string]any
		want    string
	}{
		{"the queue batches", true, map[string]any{"story": f.story.ID}, "ok=passed late=queued"},
		{"all", true, map[string]any{"story": f.story.ID, "all": true}, "ok=reused late=passed"},
		{"the queue cannot batch", false, map[string]any{"story": f.story.ID}, "ok=reused late=reused"},
	} {
		batches, asked = c.batches, ""
		out, failed := f.call(t, "verify", c.args)
		if failed != "" {
			t.Fatalf("%s: %s", c.name, failed)
		}
		states, _ := steps(t, out)
		if want := "rebase=passed sync=passed narrative=passed check=passed " + c.want; out["passed"] != true || states != want {
			t.Errorf("%s: passed %v, steps %s, want %s", c.name, out["passed"], states, want)
		}
		if c.args["all"] == nil && asked != projectRoot(f.repo) {
			t.Errorf("%s: asked for %q, want the project's root %q", c.name, asked, projectRoot(f.repo))
		}
	}

	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range res.Tools {
		if tool.Name != "verify" {
			continue
		}
		for _, want := range []string{"marks queue", "queued, left to the merge queue", "all leaves no tier to the queue"} {
			if !strings.Contains(tool.Description, want) {
				t.Errorf("the description does not say %q: %s", want, tool.Description)
			}
		}
	}
}
