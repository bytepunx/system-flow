package mcpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// queueFixture is verifyFixture served to the orchestrator, with its story
// in review, accept_reviews on when permit says so, a change to a file under
// the story's touches committed on its branch, and flai verify passed at
// its head, which head is: a story the merge queue finds eligible.
func queueFixture(t *testing.T, permit bool) (f *fixture, head string) {
	t.Helper()
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	f, _, git := verifyFixtureWith(t, func(o *Options) { o.Agent = workitem.ActivityOrchestrator })
	f.repo.Manifest.Orchestration.Permissions.AcceptReviews = permit
	wt := f.repo.WorktreePath(f.story.ID)
	file := filepath.Join(wt, "flai", "internal", "mcpserver", "x.go")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("package mcpserver\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(wt, "add", "-A")
	git(wt, "commit", "-q", "-m", "feat: x")
	for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
		if _, err := f.repo.Transition(f.task, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	story, err := f.repo.Get(f.story.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Transition(story, workitem.Review, "alex", "", t0); err != nil {
		t.Fatal(err)
	}
	head = gitIn(t, f.repo.Root, "rev-parse", "story/"+f.story.ID)
	if err := verify.SaveReport(f.repo, verify.Report{Story: f.story.ID, Commit: head, Passed: true}); err != nil {
		t.Fatal(err)
	}
	return f, head
}

// S-0403 criterion 1: queue_add enqueues a story in review as the
// orchestrator, recording the commit its verifier checked and its evidence
// as given, and answers the change.
func TestQueueAddEnqueuesWithTheOrchestratorsEvidence(t *testing.T) {
	f, head := queueFixture(t, true)
	evidence := "Verdict: pass\n- 1: `flai/internal/mcpserver/x.go`"
	out, failed := f.call(t, "queue_add", map[string]any{"story": f.story.ID, "verified": head, "evidence": evidence})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["action"] != "add" || out["story"] != f.story.ID || out["changed"] != true || out["by"] != workitem.ActivityOrchestrator {
		t.Errorf("the change: %v", out)
	}
	st, err := mergequeue.ReadState(f.repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Enqueued) != 1 {
		t.Fatalf("enqueued %+v", st.Enqueued)
	}
	if e := st.Enqueued[0]; e.Story != f.story.ID || e.By != workitem.ActivityOrchestrator || e.Verified != head || e.Evidence != evidence || !e.At.Equal(*f.clock) {
		t.Errorf("enqueued %+v", e)
	}
}

// S-0403 criterion 1: queue_add is refused, and nothing written, while
// accept_reviews is off, for evidence that leaves a criterion out or names
// no changed file for it, and for a verified commit that is not the head.
func TestQueueAddRefusesAsTheOrchestratorsAcceptanceWould(t *testing.T) {
	off, head := queueFixture(t, false)
	good := "Verdict: pass\n- 1: `flai/internal/mcpserver/x.go`"
	if _, failed := off.call(t, "queue_add", map[string]any{"story": off.story.ID, "verified": head, "evidence": good}); !strings.Contains(failed, "only with orchestration.permissions.accept_reviews, which is off") {
		t.Errorf("accept_reviews off: %q", failed)
	}
	f, head := queueFixture(t, true)
	main := gitIn(t, f.repo.Root, "rev-parse", "main")
	for name, c := range map[string]struct{ verified, evidence, want string }{
		"criterion missing": {head, "Verdict: pass\n- 2: `flai/internal/mcpserver/x.go`", "criterion_unevidenced: the evidence has no item for acceptance criterion 1 (works)"},
		"file unchanged":    {head, "Verdict: pass\n- 1: `flai/internal/mcpserver/y.go`", "criterion_unevidenced: the evidence for acceptance criterion 1 (works) names no file the branch changes"},
		"no verdict":        {head, "- 1: `flai/internal/mcpserver/x.go`", "the evidence has no Verdict line"},
		"not the head":      {"main", good, "unverified: --verified " + main[:12] + " is not the head of story/" + f.story.ID},
	} {
		_, failed := f.call(t, "queue_add", map[string]any{"story": f.story.ID, "verified": c.verified, "evidence": c.evidence})
		if !strings.Contains(failed, c.want) {
			t.Errorf("%s: %q, want %q", name, failed, c.want)
		}
	}
	for _, x := range []*fixture{off, f} {
		if _, err := os.Stat(mergequeue.StatePath(x.repo.CacheDir())); err == nil {
			t.Error("a refused queue_add wrote the state")
		}
	}
}
