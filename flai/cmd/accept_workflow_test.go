package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0744 criteria 3 and 6: on a board with a state renamed and one added
// between in-progress and review, flai accept refuses a story in the added
// state in today's words, naming it and the state the acceptance leaves
// from, and its dry run of a story in that state finds nothing to stop it
// and changes nothing.
func TestAcceptTakesAStoryOnlyFromTheStateTheAcceptanceLeaves(t *testing.T) {
	root := renamedProject(t)
	mustRun(t, root, "story", "new", "Shipped")
	file := filepath.Join(root, "wip/kanban/stories/S-0002-shipped.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, root, "task", "new", "--story", "S-0002", "Step")
	for _, m := range [][]string{{"S-0002", "ready"}, {"S-0002", "in-progress"}, {"S-0002", "testing"}, {"T-0001", "ready"}, {"T-0001", "in-progress"}, {"T-0001", "shipped"}} {
		mustRun(t, root, "move", m[0], m[1])
	}

	_, errOut, code := runIn(t, root, "accept", "S-0002", "--by", "alex")
	if want := "S-0002 is testing; a story is accepted from acceptance-review"; code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("a story in testing must be refused with %q: exit %d\n%s", want, code, errOut)
	}
	if got := stateOf(t, root, "S-0002"); got != "testing" {
		t.Errorf("a refused acceptance changes nothing: status %q, want testing", got)
	}

	mustRun(t, root, "move", "S-0002", "acceptance-review")
	if out := mustRun(t, root, "accept", "S-0002", "--by", "alex", "--dry-run"); !strings.Contains(out, "dry run: nothing changed") || strings.Contains(out, "block") {
		t.Errorf("the dry run of a story in acceptance-review finds nothing to stop it:\n%s", out)
	}
	if got := stateOf(t, root, "S-0002"); got != "acceptance-review" {
		t.Errorf("a dry run changes nothing: status %q, want acceptance-review", got)
	}
}

// I-0377, S-0842: on a board whose accepted state is shipped, flai accept
// walked a story to shipped and then refused to archive it, "S-0002 is
// shipped; only done or cancelled items can be archived", because an item
// read as closed only in done or cancelled. The acceptance archives the story,
// its task, and its narrative in shipped and commits them.
func TestAcceptArchivesAStoryInTheAcceptedStateOfTheWorkflow(t *testing.T) {
	root := renamedProject(t)
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}} {
		gitIn(t, root, args...)
	}
	mustRun(t, root, "story", "new", "Shipped")
	file := filepath.Join(root, "wip/kanban/stories/S-0002-shipped.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, root, "task", "new", "--story", "S-0002", "Step")
	for _, m := range [][]string{{"S-0002", "ready"}, {"S-0002", "in-progress"}} {
		mustRun(t, root, "move", m[0], m[1])
	}
	mustRun(t, root, "stream", "open", "S-0002")
	for _, m := range [][]string{{"T-0001", "ready"}, {"T-0001", "in-progress"}, {"T-0001", "shipped"}, {"S-0002", "testing"}, {"S-0002", "acceptance-review"}} {
		mustRun(t, root, "move", m[0], m[1])
	}

	out, errOut, code := runIn(t, root, "accept", "S-0002", "--by", "alex")
	if code != 0 {
		t.Fatalf("a story in acceptance-review is accepted to shipped and archived: exit %d\n%s%s", code, out, errOut)
	}
	for _, f := range []string{"kanban/stories/S-0002-shipped.md", "kanban/tasks/T-0001-step.md"} {
		b, err := os.ReadFile(filepath.Join(root, "wip/archive", f))
		if err != nil || !strings.Contains(string(b), "status: shipped\n") {
			t.Errorf("%s is archived in shipped: %v\n%s", f, err, b)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "wip/archive/agents/S-0002.md")); err != nil {
		t.Errorf("the story's narrative is archived: %v", err)
	}
	if got := gitIn(t, root, "log", "-1", "--format=%s"); !strings.HasPrefix(got, "chore: [S-0002] accept and archive") {
		t.Errorf("the acceptance is committed: last commit %q", got)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("the acceptance leaves the working tree clean:\n%s", st)
	}
}
