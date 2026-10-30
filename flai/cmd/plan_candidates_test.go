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
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0219, S-0328: flai plan --candidates lists the epics and then the
// stories the planner should plan, each with why, leaves out one a planner
// runs for and one whose planner waits on the operator, and a story whose
// epic a planner runs for or whose planner ended with it unchanged since,
// and writes nothing.
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
	run("story", "new", "Unplanned")      // S-0006: a candidate
	run("story", "new", "Planned before") // S-0007: left out, unchanged since its planner ended
	run("story", "new", "Changed since")  // S-0008: a candidate, changed since its planner ended
	run("story", "new", "Asking")         // S-0009: left out, its planner asked
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
	run("thread", "new", "--on", "S-0009", "--by", "planner-S-0009", "What does it cost?", "Say how much.")

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
	s7, err := repo.Get("S-0007")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := time.Parse(workitem.TimeFormat, s7.Updated)
	if err != nil {
		t.Fatal(err)
	}
	endedAfter := updated.Add(time.Hour).Format(workitem.TimeFormat)
	states := map[string]serve.AgentState{mainRootOf(repo): {Plans: map[string]*serve.AgentRun{
		"E-0006": {Item: "E-0006", Agent: "planner-E-0006", PID: sleep.Process.Pid, Start: serve.Started(sleep.Process.Pid), Started: "2026-09-15T20:30:00Z"},
		"E-0003": {Item: "E-0003", Agent: "planner-E-0003", PID: sleep.Process.Pid, Start: serve.Started(sleep.Process.Pid), Started: "2026-09-15T20:40:00Z"},
		"S-0007": {Item: "S-0007", Agent: "planner-S-0007", Started: "2026-09-15T19:00:00Z", Ended: endedAfter, Outcome: serve.OutcomeWorked},
		"S-0008": {Item: "S-0008", Agent: "planner-S-0008", Started: "2026-09-15T19:00:00Z", Ended: ended, Outcome: serve.OutcomeWorked},
		"S-0009": {Item: "S-0009", Agent: "planner-S-0009", Started: "2026-09-15T19:00:00Z", Ended: endedAfter, Outcome: serve.OutcomeAsked, Thread: "TH-0003"},
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
		"to plan:\n  E-0001  Empty\n    - in the backlog with no stories",
		"  E-0002  Finished stories\n    - every story done or cancelled (1 done, 1 cancelled) while the epic is in-progress",
		"  E-0008  Answered\n    - in the backlog with no stories",
		"  S-0006  Unplanned\n    - in the backlog without a plan: no touches; no forecast duration; no cost of delay value; no tasks\n",
		"  S-0008  Changed since\n    - in the backlog without a plan: ",
		"left out:\n  E-0006  Being planned\n    - a planner runs for it now (pid ",
		"  E-0007  Waiting on operator\n    - its planner asked on TH-0001, which awaits the operator: Which outcome?\n",
		"  S-0004  Still open\n    - a planner runs for its epic E-0003 now (pid ",
		"  S-0007  Planned before\n    - its planner ran until " + endedAfter + " and it has not changed since",
		"  S-0009  Asking\n    - its planner asked on TH-0003, which awaits the operator: What does it cost?\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"  E-0003  ", "  E-0004  ", "  E-0005  "} {
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
	if !reflect.DeepEqual(ids, []string{"E-0001", "E-0002", "E-0008", "S-0006", "S-0008"}) || !reflect.DeepEqual(left, []string{"E-0006", "E-0007", "S-0004", "S-0007", "S-0009"}) {
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

// S-0524, ADR-0222, I-0188: flai plan --candidates lists an epic or a story
// planned before when the operator has written, since its newest planner run
// ended, on a thread on it or on an item under it, naming the thread; an
// entry by an agent, or one older than the run's end, lists nothing, and the
// holds of a running planner for its epic and of its planner's question
// awaiting the operator still leave it out.
func TestPlanCandidatesAskedAgainCommand(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(m, "owner: alex\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	at := func(hhmm string, args ...string) {
		t.Helper()
		when, err := time.Parse(time.RFC3339, "2026-09-15T"+hhmm+":00Z")
		if err != nil {
			t.Fatal(err)
		}
		if _, errOut, code := runInAt(t, root, when, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	plan := func(id string) {
		t.Helper()
		at("18:00", "edit", id, "--by", "alex", "--touches", "flai/cmd", "--forecast-duration", "6h", "--cost-of-delay-value", "900")
	}
	at("18:00", "epic", "new", "Confirmed on a story")         // E-0001: a candidate, I-0188
	at("18:00", "epic", "new", "Agents wrote")                 // E-0002: not one
	at("18:00", "epic", "new", "Wrote before its planner end") // E-0003: not one
	at("18:00", "epic", "new", "Its planner asks")             // E-0004: left out, its planner asked
	at("18:00", "epic", "new", "Being planned")                // E-0005: a planner runs
	at("18:00", "story", "new", "A new story confirmed", "--epic", "E-0001")
	at("18:00", "story", "new", "Agents only", "--epic", "E-0002")
	at("18:00", "story", "new", "Too early", "--epic", "E-0003")
	at("18:00", "story", "new", "Under the asking epic", "--epic", "E-0004")
	at("18:00", "story", "new", "Planned, asked on a task") // S-0005: a candidate
	at("18:00", "story", "new", "Under the running epic", "--epic", "E-0005")
	at("18:00", "story", "new", "Words changed, nobody wrote") // S-0007: not one
	at("18:00", "task", "new", "--story", "S-0004", "Part")
	at("18:00", "task", "new", "--story", "S-0005", "Part")
	at("18:00", "task", "new", "--story", "S-0006", "Part")
	at("18:00", "task", "new", "--story", "S-0007", "Part")
	plan("S-0005")
	plan("S-0006")
	plan("S-0007")
	at("20:30", "edit", "S-0007", "--by", "alex", "--title", "Words changed since, nobody wrote")
	at("19:00", "thread", "new", "--on", "S-0001", "--by", "agent-S-0001", "One more story?", "Should E-0024 gain one?")
	at("20:30", "thread", "reply", "TH-0001", "--by", "alex", "Yes, add it.")
	at("21:00", "thread", "reply", "TH-0001", "--by", "planner-E-0001", "Noted.")
	at("21:00", "thread", "new", "--on", "S-0002", "--by", "planner-E-0002", "A question", "Which?")
	at("21:00", "thread", "reply", "TH-0002", "--by", "orchestrator", "That one.")
	at("21:00", "thread", "new", "--on", "S-0003", "--by", "alex", "Too early", "Before its planner ended.")
	at("19:30", "thread", "new", "--on", "E-0004", "--by", "planner-E-0004", "Which outcome?", "Say which.")
	at("21:00", "thread", "new", "--on", "T-0001", "--by", "alex", "On a task", "Split it.")
	at("21:00", "thread", "new", "--on", "T-0002", "--by", "alex", "On a task", "Split it.")
	at("21:00", "thread", "new", "--on", "S-0006", "--by", "alex", "On the story", "Drop a task.")

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
	worked := func(id, end string) *serve.AgentRun {
		return &serve.AgentRun{Item: id, Agent: "planner-" + id, Started: "2026-09-15T19:00:00Z", Ended: end, Outcome: serve.OutcomeWorked}
	}
	states := map[string]serve.AgentState{mainRootOf(repo): {Plans: map[string]*serve.AgentRun{
		"E-0001": worked("E-0001", ended),
		"E-0002": worked("E-0002", ended),
		"E-0003": worked("E-0003", "2026-09-15T22:00:00Z"),
		"E-0004": {Item: "E-0004", Agent: "planner-E-0004", Started: "2026-09-15T19:00:00Z", Ended: ended, Outcome: serve.OutcomeAsked, Thread: "TH-0004"},
		"E-0005": {Item: "E-0005", Agent: "planner-E-0005", PID: sleep.Process.Pid, Start: serve.Started(sleep.Process.Pid), Started: "2026-09-15T20:30:00Z"},
		"S-0005": worked("S-0005", ended),
		"S-0006": worked("S-0006", ended),
		"S-0007": worked("S-0007", ended),
	}}}
	data, _ := json.Marshal(states)
	dir := string(serve.DirFor(cfg))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "plan", "--candidates")
	if code != 0 {
		t.Fatalf("flai plan --candidates: %s", errOut)
	}
	for _, want := range []string{
		"  E-0001  Confirmed on a story\n    - planned before; alex wrote on TH-0001 (S-0001) at 2026-09-15T20:30:00Z, after its planner ran until 2026-09-15T20:00:00Z\n",
		"  S-0005  Planned, asked on a task\n    - planned before; alex wrote on TH-0006 (T-0002) at 2026-09-15T21:00:00Z, after its planner ran until 2026-09-15T20:00:00Z\n",
		"left out:\n",
		"  E-0004  Its planner asks\n    - its planner asked on TH-0004, which awaits the operator: Which outcome?\n",
		"  S-0006  Under the running epic\n    - a planner runs for its epic E-0005 now (pid ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"  E-0002  ", "  E-0003  ", "  E-0005  ", "  S-0007  "} {
		if strings.Contains(out, not) {
			t.Errorf("%s listed:\n%s", not, out)
		}
	}

	out, errOut, code = runIn(t, root, "plan", "--candidates", "--json")
	if code != 0 {
		t.Fatalf("flai plan --candidates --json: %s", errOut)
	}
	var j struct {
		Candidates []struct{ ID, Thread string } `json:"candidates"`
		LeftOut    []struct{ ID, Thread string } `json:"left_out"`
	}
	if err := json.Unmarshal([]byte(out), &j); err != nil {
		t.Fatal(err)
	}
	threadOf := func(list []struct{ ID, Thread string }) map[string]string {
		got := map[string]string{}
		for _, c := range list {
			got[c.ID] = c.Thread
		}
		return got
	}
	if got, want := threadOf(j.Candidates), map[string]string{"E-0001": "TH-0001", "S-0001": "", "S-0002": "", "S-0003": "", "S-0004": "", "S-0005": "TH-0006"}; !reflect.DeepEqual(got, want) {
		t.Errorf("--json: candidates and their threads %v, want %v", got, want)
	}
	if got, want := threadOf(j.LeftOut), map[string]string{"E-0004": "TH-0005", "S-0006": "TH-0007"}; !reflect.DeepEqual(got, want) {
		t.Errorf("--json: left out and their threads %v, want %v", got, want)
	}
}
