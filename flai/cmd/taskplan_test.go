package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// planProject is heldProject with three tasks under S-0001: T-0001 in
// progress, T-0002 after it, and T-0003 after none.
func planProject(t *testing.T) string {
	t.Helper()
	root := heldProject(t)
	for _, args := range [][]string{
		{"task", "new", "First", "--story", "S-0001"},
		{"task", "new", "Second", "--story", "S-0001", "--after", "T-0001"},
		{"task", "new", "Third", "--story", "S-0001"},
		{"move", "T-0001", "ready"},
		{"move", "T-0001", "in-progress"},
	} {
		if out, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s %s", args, code, out, errOut)
		}
	}
	return root
}

// S-0176: flai show prints a story's task plan, each task's state and what
// it waits for and the layers, and --json carries it as plan beside item
// and children; a story without tasks has none.
func TestShowPrintsTheTaskPlan(t *testing.T) {
	root := planProject(t)
	out, errOut, code := runIn(t, root, "show", "S-0001")
	if code != 0 {
		t.Fatal(errOut)
	}
	want := "  plan:\n" +
		"    T-0001  in-progress\n" +
		"    T-0002  waiting      for T-0001\n" +
		"    T-0003  ready\n" +
		"  layers:\n" +
		"    1  T-0001, T-0003\n" +
		"    2  T-0002\n"
	if !strings.HasSuffix(out, want) {
		t.Errorf("show:\n%s\nwant it to end:\n%s", out, want)
	}
	if out, _, _ := runIn(t, root, "show", "T-0002"); !strings.Contains(out, "  after: T-0001\n") || strings.Contains(out, "plan:") {
		t.Errorf("a task's show:\n%s", out)
	}

	out, errOut, code = runIn(t, root, "show", "S-0001", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	var plan any
	_ = json.Unmarshal(got["plan"], &plan)
	compact, _ := json.Marshal(plan)
	if string(compact) != `{"layers":[["T-0001","T-0003"],["T-0002"]],"tasks":[{"id":"T-0001","state":"in-progress"},{"after":["T-0001"],"id":"T-0002","state":"waiting","waiting_for":["T-0001"]},{"id":"T-0003","state":"ready"}]}` {
		t.Errorf("plan: %s", compact)
	}
	if got["item"] == nil || got["children"] == nil {
		t.Errorf("item and children: %s", out)
	}
	if out, _, _ := runIn(t, root, "show", "S-0002", "--json"); strings.Contains(out, `"plan"`) {
		t.Errorf("a story without tasks:\n%s", out)
	}
}

// S-0176: flai board gives a story with tasks its plan in counts, as text
// and as JSON, and none to a story without tasks.
func TestBoardCountsAStorysTasks(t *testing.T) {
	root := planProject(t)
	out, errOut, code := runIn(t, root, "board")
	if code != 0 {
		t.Fatal(errOut)
	}
	if !strings.Contains(out, "         tasks 1 ready, 1 waiting, 1 in progress, 0 done; 2 layers\n") || strings.Count(out, "tasks ") != 1 {
		t.Errorf("board:\n%s", out)
	}
	out, errOut, code = runIn(t, root, "board", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var view struct {
		Columns map[string][]struct {
			ID    string          `json:"id"`
			Tasks json.RawMessage `json:"tasks"`
		} `json:"columns"`
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatal(err)
	}
	ip := view.Columns["in-progress"]
	if len(ip) != 1 {
		t.Fatalf("in progress: %+v", ip)
	}
	var tasks bytes.Buffer
	_ = json.Compact(&tasks, ip[0].Tasks)
	if tasks.String() != `{"ready":1,"waiting":1,"in_progress":1,"done":0,"layers":2}` {
		t.Errorf("tasks: %s", tasks.String())
	}
	if r := view.Columns["ready"]; len(r) != 1 || r[0].Tasks != nil {
		t.Errorf("ready: %s %s", r[0].ID, r[0].Tasks)
	}
}

// S-0176: moving a task to in progress while a task of its after is open
// warns, and moves it.
func TestMovingAWaitingTaskWarnsAndMoves(t *testing.T) {
	root := planProject(t)
	if _, errOut, code := runIn(t, root, "move", "T-0002", "ready"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code := runIn(t, root, "move", "T-0002", "in-progress")
	if code != 0 || !strings.Contains(out, "T-0002 → in-progress") {
		t.Fatalf("move: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "T-0002 is waiting (after): waits for T-0001 (in progress); ready to start when T-0001 is done or cancelled") {
		t.Errorf("no warning: %s", errOut)
	}
	if _, errOut, _ := runIn(t, root, "move", "T-0003", "ready"); errOut != "" {
		t.Fatal(errOut)
	}
	if _, errOut, _ := runIn(t, root, "move", "T-0003", "in-progress"); strings.Contains(errOut, "waiting") {
		t.Errorf("a task that waits for nothing: %s", errOut)
	}
}
