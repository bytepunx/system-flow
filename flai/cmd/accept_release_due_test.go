package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// releaseDueEvents are the release-due events in root's event log.
func releaseDueEvents(t *testing.T, root string) []inbox.LoggedEvent {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var due []inbox.LoggedEvent
	for _, e := range inbox.LoggedEvents(repo) {
		if e.Kind == inbox.ReleaseDue {
			due = append(due, e)
		}
	}
	return due
}

// commitPolicy sets root's orchestration.release to yaml and commits it, so
// that the acceptance finds nothing uncommitted outside wip.
func commitPolicy(t *testing.T, root, yaml string) {
	t.Helper()
	setReleasePolicy(t, root, yaml)
	gitIn(t, root, "add", "system-flow.yaml")
	gitIn(t, root, "commit", "-q", "-m", "chore: release policy")
}

// ADR-0211: the acceptance of a story that changes the component makes a
// threshold of one story met, and appends one release-due event with both
// answers and the acceptor.
func TestAcceptAppendsReleaseDueWhenTheReleaseBecomesDue(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	commitPolicy(t, root, "    policy: threshold\n    count: 1\n")

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	assertAccepted(t, root, remote)
	due := releaseDueEvents(t, root)
	if len(due) != 1 {
		t.Fatalf("one release-due event, got %+v", due)
	}
	e := due[0]
	if e.ID != "S-0001" || e.Title != "What we found" || e.By != "alex" || e.Type != workitem.Story || e.At == "" {
		t.Errorf("the event names the story, its title, and the acceptor: %+v", e)
	}
	if got, want := e.Fields[inbox.FieldBefore], "not met: no accepted story waiting for a release"; got != want {
		t.Errorf("before:\n got %s\nwant %s", got, want)
	}
	if got, want := e.Fields[inbox.FieldAfter], "met: 1 accepted story not yet released, worth 0 USD/week"; got != want {
		t.Errorf("after:\n got %s\nwant %s", got, want)
	}
}

// ADR-0211: under judgement, never met by itself, flai move S-nnnn done
// adding a pending story changes the figures the answer holds, so it is
// release-due too.
func TestMoveDoneAppendsReleaseDueWhenTheFiguresChange(t *testing.T) {
	root, _ := researchProject(t, "feature", true)
	if _, errOut, code := runIn(t, root, "move", "S-0001", "done", "--by", "alex"); code != 0 {
		t.Fatalf("move to done: %d %s", code, errOut)
	}
	due := releaseDueEvents(t, root)
	if len(due) != 1 || due[0].Fields[inbox.FieldBefore] != "not met: no accepted story waiting for a release" ||
		due[0].Fields[inbox.FieldAfter] != "not met: 1 accepted story not yet released, worth 0 USD/week" {
		t.Errorf("one release-due event, from no pending story to one, still not met: %+v", due)
	}
}

// ADR-0211: a research story that changes no component adds nothing pending,
// so the answer is the same before and after, and no event is appended.
func TestAcceptAppendsNoReleaseDueWhenTheAnswerStays(t *testing.T) {
	root, remote := researchProject(t, "research", false)
	commitPolicy(t, root, "    policy: threshold\n    count: 1\n")

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	assertAccepted(t, root, remote)
	if due := releaseDueEvents(t, root); len(due) != 0 {
		t.Errorf("an acceptance that leaves the answer as it was appends no release-due event: %+v", due)
	}
}

// ADR-0211: an evaluation that fails, here over a template component whose
// template.yaml is missing, is warned of and appends nothing, and the
// acceptance is made as it would be without it.
func TestAcceptStandsWhenTheReleaseCannotBeEvaluated(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	path := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(data), "    kind: go\n", "    kind: go\n  - name: tpl\n    path: tpl\n    kind: template\n", 1)
	if broken == string(data) {
		t.Fatalf("no component to add a template after: %s", data)
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "system-flow.yaml")
	gitIn(t, root, "commit", "-q", "-m", "chore: a template with no template.yaml")
	if _, _, code := runIn(t, root, "release", "--evaluate"); code == 0 {
		t.Fatal("the release cannot be evaluated with tpl/template.yaml missing")
	}

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("a failed evaluation fails no acceptance: %d %s", code, errOut)
	}
	assertAccepted(t, root, remote)
	if !strings.Contains(errOut, "release not evaluated") {
		t.Errorf("the failed evaluation is warned of: %s", errOut)
	}
	if due := releaseDueEvents(t, root); len(due) != 0 {
		t.Errorf("a failed evaluation appends nothing: %+v", due)
	}
}
