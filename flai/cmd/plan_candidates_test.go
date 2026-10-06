//go:build !windows

package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0219: flai plan --candidates lists the epics the planner should plan,
// each with why, leaves out one a planner runs for and one whose planner
// waits on the operator, and writes nothing.
func TestPlanCandidatesCommand(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	// status sets an item's status in its file, as a fixture: the moves to
	// done are the operator's, and beside the point here.
	status := func(id, to string) {
		t.Helper()
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/*", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		if err := os.WriteFile(file[0], []byte(strings.Replace(string(s), "status: backlog\n", "status: "+to+"\n", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("epic", "new", "Empty")               // E-0001: a candidate
	run("epic", "new", "Finished stories")    // E-0002: a candidate
	run("epic", "new", "Open story")          // E-0003: not one
	run("epic", "new", "Done epic")           // E-0004: not one
	run("epic", "new", "Cancelled epic")      // E-0005: not one
	run("epic", "new", "Being planned")       // E-0006: left out, a planner runs
	run("epic", "new", "Waiting on operator") // E-0007: left out, its planner asked
	run("epic", "new", "Answered")            // E-0008: a candidate, its question answered
	run("story", "new", "Shipped", "--epic", "E-0002")
	run("story", "new", "Dropped", "--epic", "E-0002")
	run("story", "new", "Shipped too", "--epic", "E-0003")
	run("story", "new", "Still open", "--epic", "E-0003")
	run("story", "new", "Shipped", "--epic", "E-0004")
	status("S-0001", workitem.Done)
	status("S-0002", workitem.Cancelled)
	status("E-0002", workitem.InProgress)
	status("S-0003", workitem.Done)
	status("S-0005", workitem.Done)
	status("E-0004", workitem.Done)
	run("move", "E-0005", "cancelled", "--reason", "dropped")
	run("thread", "new", "--on", "E-0007", "--by", "planner-E-0007", "Which outcome?", "Say which.")
	run("thread", "new", "--on", "E-0008", "--by", "planner-E-0008", "Which users?", "Say which.")
	run("thread", "reply", "TH-0002", "--by", "alex", "These.")

	sleep := exec.Command("sleep", "60")
	serve.Detach(sleep)
	if err := sleep.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sleep.Process.Kill(); _ = sleep.Wait() })
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	ended := "2026-09-15T20:00:00Z"
	states := map[string]serve.AgentState{mainRootOf(repo): {Plans: map[string]*serve.AgentRun{
		"E-0006": {Item: "E-0006", Agent: "planner-E-0006", PID: sleep.Process.Pid, Start: serve.Started(sleep.Process.Pid), Started: "2026-09-15T20:30:00Z"},
		"E-0007": {Item: "E-0007", Agent: "planner-E-0007", Started: "2026-09-15T19:00:00Z", Ended: ended, Outcome: serve.OutcomeAsked, Thread: "TH-0001"},
		"E-0008": {Item: "E-0008", Agent: "planner-E-0008", Started: "2026-09-15T19:00:00Z", Ended: ended, Outcome: serve.OutcomeAsked, Thread: "TH-0002"},
		"E-0001": {Item: "E-0001", Agent: "planner-E-0001", PID: 999999, Started: "2026-09-15T19:00:00Z", Ended: ended, Outcome: serve.OutcomeWorked},
	}}}
	data, _ := json.Marshal(states)
	dir := string(serve.DirFor(cfg))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	before := treeOf(t, root)
	out := run("plan", "--candidates")
	for _, want := range []string{
		"epics to plan:\n  E-0001  Empty\n    - in the backlog with no stories",
		"  E-0002  Finished stories\n    - every story done or cancelled (1 done, 1 cancelled) while the epic is in-progress",
		"  E-0008  Answered\n    - in the backlog with no stories",
		"left out:\n  E-0006  Being planned\n    - a planner runs for it now (pid ",
		"  E-0007  Waiting on operator\n    - its planner asked on TH-0001, which awaits the operator: Which outcome?\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"E-0003", "E-0004", "E-0005"} {
		if strings.Contains(out, not) {
			t.Errorf("%s listed:\n%s", not, out)
		}
	}
	if !reflect.DeepEqual(treeOf(t, root), before) {
		t.Error("flai plan --candidates changed a file")
	}

	var j struct {
		Candidates []struct{ ID, Reason string } `json:"candidates"`
		LeftOut    []struct{ ID, Reason string } `json:"left_out"`
	}
	if err := json.Unmarshal([]byte(run("plan", "--candidates", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	var ids, left []string
	for _, c := range j.Candidates {
		ids = append(ids, c.ID)
		if c.Reason == "" {
			t.Errorf("%s has no reason", c.ID)
		}
	}
	for _, c := range j.LeftOut {
		left = append(left, c.ID)
	}
	if !reflect.DeepEqual(ids, []string{"E-0001", "E-0002", "E-0008"}) || !reflect.DeepEqual(left, []string{"E-0006", "E-0007"}) {
		t.Errorf("--json: candidates %v, left out %v", ids, left)
	}
	if !reflect.DeepEqual(treeOf(t, root), before) {
		t.Error("flai plan --candidates --json changed a file")
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"plan", "--candidates", "E-0001"}, "unknown command"},
		{[]string{"plan"}, "accepts 1 arg"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
}
