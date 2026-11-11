package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// klusterStoryInReview is orchestratedStoryInReview with accept_reviews on
// and a kluster: block, naming kluster.yaml, committed on main after the
// story's branch was made; it returns the main checkout and the branch's
// head.
func klusterStoryInReview(t *testing.T) (root, head string) {
	t.Helper()
	root, head = orchestratedStoryInReview(t, true)
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	writeFiles(t, root, map[string]string{
		"system-flow.yaml": string(m) + "kluster:\n  config: kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n",
		"kluster.yaml":     "name: t-kluster\nprofile: authstar\nprovider: k3d\n",
	})
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: run the project's kluster")
	return root, head
}

// klusterOn turns the kluster host action on for the project at root in t's
// configuration, and the checks host action with checks.
func klusterOn(t *testing.T, root string, checks bool) {
	t.Helper()
	actions := map[string][]string{"kluster": {root}}
	if checks {
		actions["checks"] = []string{root}
	}
	cfg, err := json.Marshal(map[string]any{"host_actions": actions})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configOf(t), cfg, 0o644); err != nil {
		t.Fatal(err)
	}
}

// writeS0001Deploy writes r as the deployment record of S-0001's namespace,
// s-0001, in the main checkout at root, as a kluster deploy keeps it; with
// r's Commit empty it removes the record, as for a namespace never deployed.
func writeS0001Deploy(t *testing.T, root string, r kluster.Record) {
	t.Helper()
	path := kluster.RecordPath(root, "s-0001")
	if r.Commit == "" {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
		return
	}
	r.Story, r.Namespace = "S-0001", "s-0001"
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// klusterAcceptances are the command lines of S-0001's three acceptances,
// the operator's, the orchestrator's with its verifier's commit head and
// evidence, and the merge queue's landing from landingBatch, by name.
func klusterAcceptances(t *testing.T, head string) map[string][]string {
	t.Helper()
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	return map[string][]string{
		"the operator's":     {"accept", "S-0001", "--by", "alex"},
		"the orchestrator's": {"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev},
		"the queue's":        {"accept", "S-0001", "--by", "alex", "--batch", landingBatch},
	}
}

// klusterRefusal is the start of the blocker kluster for S-0001 (ADR-0254
// §5a).
const klusterRefusal = "the kluster host action is on, so S-0001 is accepted only once its namespace was deployed at its head"

// I-0268, ADR-0254 §5a: while the kluster host action is on, the operator's
// acceptance, the orchestrator's, and the merge queue's landing are each
// refused with the blocker kluster, before anything is merged, for a
// namespace never deployed, one deployed at an older head, a deploy that
// failed, and checks that failed against it at the head, each naming why.
func TestEveryAcceptanceIsRefusedWithKlusterUntilItsNamespacePassedAtItsHead(t *testing.T) {
	t.Parallel()
	root, head := klusterStoryInReview(t)
	klusterOn(t, root, true)
	landingState(t, root, head, "alex", nil)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	for _, c := range []struct {
		name   string
		record kluster.Record
		why    string
	}{
		{"never deployed", kluster.Record{}, "S-0001's namespace s-0001 was never deployed"},
		{"deployed at an older head", kluster.Record{Commit: mainBefore, Outcome: kluster.OutcomePassed, ChecksCommit: mainBefore, ChecksOutcome: kluster.OutcomePassed},
			"the last deploy of S-0001's namespace s-0001 was at " + mainBefore[:12] + ", not at the story's head " + head[:12] + ", and it passed"},
		{"deploy failed", kluster.Record{Commit: head, Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy},
			"the deploy of S-0001's namespace s-0001 at " + head[:12] + " failed at the deploy step"},
		{"checks failed", kluster.Record{Commit: head, Outcome: kluster.OutcomePassed, ChecksCommit: head, ChecksOutcome: kluster.OutcomeFailed},
			"the deploy of S-0001's namespace s-0001 at " + head[:12] + " passed, and the checks run against it there ended failed"},
	} {
		writeS0001Deploy(t, root, c.record)
		for name, args := range klusterAcceptances(t, head) {
			_, errOut, code := runIn(t, root, args...)
			if code == 0 || !strings.Contains(errOut, klusterRefusal) || !strings.Contains(errOut, c.why) || !strings.Contains(errOut, "flai accept S-0001 --without-kluster (ADR-0254)") {
				t.Errorf("%s: %s acceptance must be refused with kluster, naming %q: %d\n%s", c.name, name, c.why, code, errOut)
			}
		}
		assertNothingAccepted(t, root, mainBefore)
	}
}

// ADR-0254 §5a: once the deploy and the checks run against it passed at the
// story's head, or the deploy alone with the checks host action off, each
// acceptance lands the story, with nothing written for the gate.
func TestEveryAcceptanceLandsOnceItsNamespacePassedAtItsHead(t *testing.T) {
	t.Parallel()
	for _, checks := range []bool{true, false} {
		for _, name := range []string{"the operator's", "the orchestrator's", "the queue's"} {
			root, head := klusterStoryInReview(t)
			klusterOn(t, root, checks)
			landingState(t, root, head, "alex", nil)
			r := kluster.Record{Commit: head, Outcome: kluster.OutcomePassed}
			if checks {
				r.ChecksCommit, r.ChecksOutcome = head, kluster.OutcomePassed
			}
			writeS0001Deploy(t, root, r)
			args := klusterAcceptances(t, head)[name]
			if out, errOut, code := runInAt(t, root, landedAt, args...); code != 0 {
				t.Errorf("checks %v: %s acceptance must land the story: %d %s\n%s", checks, name, code, errOut, out)
				continue
			}
			if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done || strings.Contains(story.Body, "### Accepted without the kluster gate") {
				t.Errorf("checks %v: %s acceptance must accept the story with nothing written for the gate: %s %v\n%s", checks, name, story.Status, story.Archived, story.Body)
			}
		}
	}
}

// ADR-0254 §5a: --without-kluster lands the operator's acceptance, by flai
// accept and flai move alike, and writes why the gate held the story, who
// accepted, and when in its Notes.
func TestWithoutKlusterLandsTheStoryAndRecordsItInItsNotes(t *testing.T) {
	t.Parallel()
	root, _ := klusterStoryInReview(t)
	klusterOn(t, root, true)
	if out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--without-kluster"); code != 0 {
		t.Fatalf("--without-kluster lands the story: %d %s\n%s", code, errOut, out)
	}
	section := "### Accepted without the kluster gate\n\n- Why: S-0001's namespace s-0001 was never deployed: no deploy of it is recorded\n- By: alex\n- At: 2026-09-17T09:00:00Z\n"
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done || !strings.Contains(story.Body, "## Notes\n") || !strings.HasSuffix(story.Body, section) {
		t.Errorf("the archived story's Notes must end with\n%s\ngot %s %v\n%s", section, story.Status, story.Archived, story.Body)
	}

	moved, head := klusterStoryInReview(t)
	klusterOn(t, moved, true)
	writeS0001Deploy(t, moved, kluster.Record{Commit: head, Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy})
	if out, errOut, code := runInAt(t, moved, landedAt, "move", "S-0001", "done", "--by", "alex", "--without-kluster"); code != 0 {
		t.Fatalf("flai move done --without-kluster lands the story: %d %s\n%s", code, errOut, out)
	}
	section = "### Accepted without the kluster gate\n\n- Why: the deploy of S-0001's namespace s-0001 at " + head[:12] + " failed at the deploy step\n- By: alex\n- At: 2026-09-17T09:00:00Z\n"
	if story := acceptedS0001(t, moved); !story.Archived || !strings.HasSuffix(story.Body, section) {
		t.Errorf("a failed deploy's Notes must end with\n%s\ngot\n%s", section, story.Body)
	}
}

// ADR-0254 §5a: --without-kluster is the operator's alone, refused with --by
// orchestrator and with --batch before anything changes.
func TestWithoutKlusterIsRefusedForTheOrchestratorAndTheQueuesLanding(t *testing.T) {
	t.Parallel()
	root, head := klusterStoryInReview(t)
	klusterOn(t, root, true)
	landingState(t, root, head, "alex", nil)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev, "--without-kluster"}, "rule: --without-kluster is the operator's: the orchestrator never lands S-0001 before its namespace's deploy and checks pass at its head"},
		{[]string{"accept", "S-0001", "--by", "alex", "--batch", landingBatch, "--without-kluster"}, "rule: --without-kluster is the operator's acceptance by hand, and --batch " + landingBatch + " is the merge queue's landing"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("flai %v must be refused with %q: %d\n%s", c.args, c.want, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)
}
