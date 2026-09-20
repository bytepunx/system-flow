package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// cancelProject is an epic with two stories, one in progress with a task and
// a narrative, one in backlog with a task.
func cancelProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	for _, args := range [][]string{
		{"config", "set", "author", "olive"},
		{"epic", "new", "Epic"},
		{"story", "new", "First", "--epic", "E-0001"},
		{"story", "new", "Second", "--epic", "E-0001"},
		{"task", "new", "One", "--story", "S-0001"},
		{"task", "new", "Two", "--story", "S-0002"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-first.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	for _, args := range [][]string{{"move", "S-0001", "ready"}, {"stream", "open", "S-0001"}, {"move", "S-0001", "in-progress"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	return root
}

func statusOf(t *testing.T, root, glob string) string {
	t.Helper()
	m, _ := filepath.Glob(filepath.Join(root, "wip/kanban", glob))
	if len(m) != 1 {
		t.Fatalf("%s: %v", glob, m)
	}
	data, _ := os.ReadFile(m[0])
	for _, l := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(l, "status: ") {
			return strings.TrimPrefix(l, "status: ")
		}
	}
	return ""
}

func TestMoveCancelledDryRunListsAndChangesNothing(t *testing.T) {
	root := cancelProject(t)
	out, errOut, code := runIn(t, root, "move", "E-0001", "cancelled", "--reason", "a different route", "--dry-run")
	if code != 0 {
		t.Fatal(errOut)
	}
	for _, want := range []string{"would also cancel 4 items", "story S-0001  in-progress", "task  T-0001  backlog", "story S-0002  backlog", "S-0001: narrative wip/agents/S-0001.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("dry run does not say %q:\n%s", want, out)
		}
	}
	if st := statusOf(t, root, "epics/E-0001-*.md"); st != "backlog" {
		t.Errorf("dry run moved the epic to %s", st)
	}
	if st := statusOf(t, root, "tasks/T-0001-*.md"); st != "backlog" {
		t.Errorf("dry run moved a task to %s", st)
	}
	if _, errOut, code := runIn(t, root, "move", "E-0001", "cancelled", "--dry-run"); code == 0 || !strings.Contains(errOut, "--reason") {
		t.Errorf("a dry run reports the refusal the real move would get: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0002", "ready", "--dry-run"); code == 0 || !strings.Contains(errOut, "nothing to preview") {
		t.Errorf("--dry-run on another move: %d %s", code, errOut)
	}
}

func TestMoveCancelledCascadesAndReportsJSON(t *testing.T) {
	root := cancelProject(t)
	out, errOut, code := runIn(t, root, "move", "E-0001", "cancelled", "--reason", "a different route", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var res cancelResult
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	ids := []string{}
	for _, c := range res.Cancelled {
		ids = append(ids, c.ID+":"+c.From)
	}
	if got := strings.Join(ids, " "); got != "S-0001:in-progress T-0001:backlog S-0002:backlog T-0002:backlog" {
		t.Errorf("cancelled: %s", got)
	}
	if res.Status != "cancelled" || len(res.LeftBehind) != 1 || res.LeftBehind[0].ID != "S-0001" || res.LeftBehind[0].Narrative == "" {
		t.Errorf("result: %+v", res)
	}
	for _, glob := range []string{"epics/E-0001-*.md", "stories/S-0001-*.md", "stories/S-0002-*.md", "tasks/T-0001-*.md", "tasks/T-0002-*.md"} {
		if st := statusOf(t, root, glob); st != "cancelled" {
			t.Errorf("%s is %s", glob, st)
		}
	}
	m, _ := filepath.Glob(filepath.Join(root, "wip/kanban/tasks/T-0002-*.md"))
	data, _ := os.ReadFile(m[0])
	if !strings.Contains(string(data), "E-0001 cancelled: a different route") || !strings.Contains(string(data), "by: olive") {
		t.Errorf("the task does not say why or who:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/agents/S-0001.md")); err != nil {
		t.Errorf("the narrative must be left alone: %v", err)
	}
}

func TestMoveCancelledAsksOnATerminal(t *testing.T) {
	root := cancelProject(t)
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	run := func(answer bool, args ...string) (string, string, int, int) {
		var out, errOut bytes.Buffer
		asked := 0
		tty := true
		a := &app{out: &out, errOut: &errOut, cwd: root, stdinIsTerminal: &tty,
			clock:   func() time.Time { return time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC) },
			confirm: func(string) (bool, error) { asked++; return answer, nil }}
		cmd := newRootCmdWith(a)
		cmd.SetArgs(args)
		code := 0
		if err := cmd.Execute(); err != nil {
			a.fail(err)
			code = 1
		}
		return out.String(), errOut.String(), code, asked
	}
	out, errOut, code, asked := run(false, "move", "E-0001", "cancelled", "--reason", "why")
	if code == 0 || asked != 1 || !strings.Contains(errOut, "nothing was cancelled") || !strings.Contains(out, "also cancels 4 items") {
		t.Errorf("declined: code %d asked %d\n%s\n%s", code, asked, out, errOut)
	}
	if st := statusOf(t, root, "stories/S-0002-*.md"); st != "backlog" {
		t.Errorf("declined, and S-0002 is %s", st)
	}
	// Nothing open under it: no question.
	if _, errOut, code, asked := run(false, "move", "T-0002", "cancelled", "--reason", "why"); code != 0 || asked != 0 {
		t.Errorf("a task: code %d asked %d %s", code, asked, errOut)
	}
	// --yes: no question.
	if _, errOut, code, asked := run(false, "move", "S-0001", "cancelled", "--reason", "why", "--yes"); code != 0 || asked != 0 {
		t.Errorf("--yes: code %d asked %d %s", code, asked, errOut)
	}
	out, errOut, code, asked = run(true, "move", "E-0001", "cancelled", "--reason", "why")
	if code != 0 || asked != 1 || !strings.Contains(out, "1 item cancelled with it") {
		t.Errorf("confirmed: code %d asked %d\n%s\n%s", code, asked, out, errOut)
	}
}
