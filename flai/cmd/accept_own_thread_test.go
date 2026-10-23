package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ownThreadWant is the start of the own_thread_open blocker for TH-0001 on
// S-0001, the thread the orchestrator opens in these tests.
const ownThreadWant = "thread TH-0001 on S-0001 is open, not resolved: Is the gate built?; the orchestrator opened it and nobody else wrote on it, so the orchestrator resolves it with a reason"

// orchestratorOpensAThread opens TH-0001 on S-0001 as the orchestrator, the
// concern it raises at step 4 of its acceptance.
func orchestratorOpensAThread(t *testing.T, root string) {
	t.Helper()
	if _, errOut, code := runIn(t, root, "thread", "new", "--on", "S-0001", "--by", "orchestrator", "Is the gate built?", "The story names a gate the diff does not build."); code != 0 {
		t.Fatal(errOut)
	}
}

// orchestratorsDryRun is the orchestrator's dry run of S-0001 at head, read
// from --json.
func orchestratorsDryRun(t *testing.T, root, head string) preview.Acceptance {
	t.Helper()
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return res
}

// blockersOf are the orchestrator's blockers of res with code.
func blockersOf(res preview.Acceptance, code string) []preview.Blocker {
	var out []preview.Blocker
	for _, b := range res.OrchestratorBlockers {
		if b.Code == code {
			out = append(out, b)
		}
	}
	return out
}

// ADR-0189: a thread the orchestrator opened on the story, that nobody else
// wrote on, is its own blocker, own_thread_open, saying the orchestrator
// resolves it; once the operator writes on it, it is thread_open as any
// other thread, with the text it has always had.
func TestTheOrchestratorsDryRunNamesAThreadOfItsOwnApart(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	orchestratorOpensAThread(t, root)

	res := orchestratorsDryRun(t, root, head)
	own := blockersOf(res, preview.BlockOwnThreadOpen)
	if len(own) != 1 || !strings.HasPrefix(own[0].Message, ownThreadWant) || !strings.Contains(own[0].Message, "flai thread resolve TH-0001 --reason") || !strings.Contains(own[0].Message, "(ADR-0189)") {
		t.Errorf("the dry run must name the orchestrator's own thread under %s: %+v", preview.BlockOwnThreadOpen, res.OrchestratorBlockers)
	}
	if others := blockersOf(res, preview.BlockThreadOpen); len(others) != 0 {
		t.Errorf("the orchestrator's own thread is not %s: %+v", preview.BlockThreadOpen, others)
	}

	if _, errOut, code := runIn(t, root, "thread", "reply", "TH-0001", "--by", "alex", "It is not needed; leave it."); code != 0 {
		t.Fatal(errOut)
	}
	res = orchestratorsDryRun(t, root, head)
	want := "thread TH-0001 on S-0001 is answered, not resolved: Is the gate built?"
	held := blockersOf(res, preview.BlockThreadOpen)
	if len(held) != 1 || held[0].Message != want {
		t.Errorf("a thread the operator wrote on is %s, %q: %+v", preview.BlockThreadOpen, want, res.OrchestratorBlockers)
	}
	if own := blockersOf(res, preview.BlockOwnThreadOpen); len(own) != 0 {
		t.Errorf("a thread the operator wrote on is not the orchestrator's own: %+v", own)
	}
}

// ADR-0189: the orchestrator's own thread holds its acceptance, naming the
// thread as its own; once the orchestrator resolves it with a reason, the
// acceptance goes through and the thread keeps the reason.
func TestTheOrchestratorAcceptsOnceItResolvesAThreadOfItsOwn(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	orchestratorOpensAThread(t, root)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev)
	if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet") || !strings.Contains(errOut, ownThreadWant) {
		t.Fatalf("the orchestrator's own thread holds the story, named as its own: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)

	if _, errOut, code := runIn(t, root, "thread", "resolve", "TH-0001", "--by", "orchestrator", "--reason", "the gate is out of scope; S-0002 builds it"); code != 0 {
		t.Fatal(errOut)
	}
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev); code != 0 {
		t.Fatalf("once its own thread is resolved the orchestrator accepts: %d %s\n%s", code, errOut, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	th, err := threads.Get(repo, "TH-0001")
	if err != nil {
		t.Fatal(err)
	}
	if th.Open() || !strings.Contains(th.Body, "Resolved: the gate is out of scope; S-0002 builds it") {
		t.Errorf("the thread is resolved and keeps its reason: %s\n%s", th.Status, th.Body)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
		t.Error("the story branch is merged")
	}
}
