package serve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// call is one call to opus as Claude Code streams it, at at.
func streamCall(session, id string, at time.Time, read int) string {
	return fmt.Sprintf(`{"type":"assistant","timestamp":%q,"session_id":%q,"message":{"id":%q,"model":"claude-opus-5-5","content":[],"usage":{"input_tokens":1,"output_tokens":1,"cache_read_input_tokens":%d,"cache_creation_input_tokens":0}}}`,
		at.UTC().Format(time.RFC3339Nano), session, id, read)
}

// streamResult is the end of a run of session: its totals, as reported.
func streamResult(session string, read int, cost float64) string {
	return fmt.Sprintf(`{"type":"result","subtype":"success","session_id":%q,"modelUsage":{"claude-opus-5-5":{"inputTokens":2,"outputTokens":300,"cacheReadInputTokens":%d,"cacheCreationInputTokens":0,"costUSD":%v}}}`, session, read, cost)
}

// taskWorked makes a task under story, in progress from `from` and done at
// `to`, or still in progress when to is zero.
func (lab *agentLab) taskWorked(story, title string, from, to time.Time) string {
	lab.t.Helper()
	tk, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: title, Parent: story, Owner: "agent", Now: lab.now})
	if err != nil {
		lab.t.Fatal(err)
	}
	for _, step := range []struct {
		to string
		at time.Time
	}{{workitem.Ready, from}, {workitem.InProgress, from}, {workitem.Done, to}} {
		if step.at.IsZero() {
			break
		}
		if _, err := lab.repo.Transition(tk, step.to, "agent", "", step.at); err != nil {
			lab.t.Fatal(err)
		}
	}
	return tk.ID
}

func (lab *agentLab) usageOf(id string) *usage.Usage {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return it.Usage
}

// A story's logs measure it whole, each task over the time it was in
// progress, and its epic is summed again.
func TestMeasureWritesAStoryItsTasksAndItsEpic(t *testing.T) {
	lab := newAgentLab(t)
	story := lab.backlog("Measured", nil)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	first := lab.taskWorked(story, "first", t0.Add(time.Minute), t0.Add(5*time.Minute))
	second := lab.taskWorked(story, "second", t0.Add(5*time.Minute), t0.Add(20*time.Minute))
	idle := lab.taskWorked(story, "idle", t0.Add(30*time.Minute), t0.Add(31*time.Minute))
	lab.taskWorked(story, "never started", time.Time{}, time.Time{})
	logs := filepath.Join(string(lab.o.Dir), "agents")
	_ = os.MkdirAll(logs, 0o700)
	write := func(name string, lines ...string) {
		if err := os.WriteFile(filepath.Join(logs, name), []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// priming before the first task, a call in each task, and the session
	// resumed in a second run; another project's log for the same story ID
	write("t-"+story+"-20260929T100000Z.log",
		streamCall("s", "m0", t0, 1000),
		streamCall("s", "m1", t0.Add(2*time.Minute), 1000),
		streamResult("s", 2000, 1.0))
	write("t-"+story+"-20260929T101000Z.log",
		streamCall("s", "m2", t0.Add(10*time.Minute), 2000),
		streamResult("s", 4000, 2.0))
	write("other-"+story+"-20260929T100000Z.log", streamResult("x", 5000, 50))

	m, err := Measure(lab.o.Dir, lab.root, "t", story, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Logs) != 2 {
		t.Errorf("logs = %v, want the project's two", m.Logs)
	}
	su := lab.usageOf(story)
	if su == nil || su.Source != usage.SourceLog || su.Estimated || su.Cost() != 2.0 || su.Models[0].CacheRead != 4000 {
		t.Fatalf("story usage = %+v, want the session's newest result", su)
	}
	// calls weigh 1001, 1001, and 2001 of 4003: the first task holds m1, the second m2
	f, s := lab.usageOf(first), lab.usageOf(second)
	if f == nil || !f.Estimated || diff(f.Cost(), 2.0*1001/4003) > 1e-3 || s == nil || diff(s.Cost(), 2.0*2001/4003) > 1e-3 {
		t.Fatalf("tasks = %+v and %+v", f, s)
	}
	if u := lab.usageOf(idle); u == nil || !u.Empty() {
		t.Errorf("a done task nothing was spent in = %+v, want measured as nothing", u)
	}
	if u := lab.usageOf(lab.epic.ID); u == nil || u.Source != usage.SourceSum || u.Cost() != 2.0 {
		t.Errorf("epic usage = %+v, want its story's", u)
	}
	if strings.Join(m.Changed, ",") != strings.Join([]string{first, second, idle, story, lab.epic.ID}, ",") {
		t.Errorf("changed = %v", m.Changed)
	}
	again, err := Measure(lab.o.Dir, lab.root, "t", story, true)
	if err != nil || len(again.Changed) != 0 {
		t.Errorf("measuring again changed %v (%v)", again.Changed, err)
	}
	if got, _ := lab.o.Dir.LoggedStories("t"); strings.Join(got, ",") != story {
		t.Errorf("logged stories = %v", got)
	}
}

// ADR-0083: measuring a story from its agents' logs keeps what the planner
// spent on it, and on its epic, and measuring again writes nothing.
func TestMeasureKeepsWhatThePlannerSpent(t *testing.T) {
	lab := newAgentLab(t)
	story := lab.backlog("Planned", nil)
	planned := &usage.Usage{Source: usage.SourceLog, Seconds: 61, Estimated: true, Models: []usage.Model{{Model: "claude-opus-5-5", Input: 2, Output: 300, CacheRead: 2000, Cost: 0.5}}}
	if _, err := lab.repo.ChargeStrategic(story, workitem.ActivityPlanner, planned); err != nil {
		t.Fatal(err)
	}
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	logs := filepath.Join(string(lab.o.Dir), "agents")
	_ = os.MkdirAll(logs, 0o700)
	if err := os.WriteFile(filepath.Join(logs, "t-"+story+"-20260929T100000Z.log"), []byte(streamCall("s", "m0", t0, 1000)+"\n"+streamResult("s", 1000, 1.0)+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := Measure(lab.o.Dir, lab.root, "t", story, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(m.Changed, ",") != story+","+lab.epic.ID {
		t.Errorf("changed = %v, want the story and its epic", m.Changed)
	}
	for _, id := range []string{story, lab.epic.ID} {
		u := lab.usageOf(id)
		if u == nil || u.Cost() != 1.0 || len(u.Strategic) != 1 || u.Strategic[0].Kind != workitem.ActivityPlanner || u.Strategic[0].Cost() != 0.5 || u.Strategic[0].Seconds != 61 {
			t.Errorf("%s's usage = %+v, want the agent's 1.0 USD and the planner's 0.5 kept", id, u)
		}
	}
	again, err := Measure(lab.o.Dir, lab.root, "t", story, true)
	if err != nil || len(again.Changed) != 0 {
		t.Errorf("measuring again changed %v (%v), want nothing", again.Changed, err)
	}
}

// subAgentCall is a call of opus in session s at at by the sub-agent the
// tool_use parent started, or by the story's agent when parent is empty,
// with content blocks.
func subAgentCall(id, parent string, at time.Time, read int, blocks string) string {
	p := "null"
	if parent != "" {
		p = fmt.Sprintf("%q", parent)
	}
	return fmt.Sprintf(`{"type":"assistant","timestamp":%q,"session_id":"s","parent_tool_use_id":%s,"message":{"id":%q,"model":"claude-opus-5-5","content":[%s],"usage":{"input_tokens":0,"output_tokens":0,"cache_read_input_tokens":%d,"cache_creation_input_tokens":0}}}`,
		at.UTC().Format(time.RFC3339Nano), p, id, blocks, read)
}

// Two tasks in progress at once, each worked by its own sub-agent, are each
// measured by their sub-agent's calls and half the story's agent's calls in
// their overlap, and together come to no more than the story.
func TestMeasureGivesEachTaskItsSubAgentAndAShareOfTheStorysAgent(t *testing.T) {
	lab := newAgentLab(t)
	story := lab.backlog("Planned", nil)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	first := lab.taskWorked(story, "first", t0.Add(time.Minute), t0.Add(10*time.Minute))
	second := lab.taskWorked(story, "second", t0.Add(time.Minute), t0.Add(20*time.Minute))
	start := func(id, task string) string {
		return fmt.Sprintf(`{"type":"tool_use","id":%q,"name":"Agent","input":{"description":"Task sub-agent: %s","prompt":"Work %s of %s."}}`, id, task, task, story)
	}
	logs := filepath.Join(string(lab.o.Dir), "agents")
	_ = os.MkdirAll(logs, 0o700)
	log := strings.Join([]string{
		subAgentCall("a1", "", t0.Add(time.Minute), 1000, start("toolu_A", first)+","+start("toolu_B", second)),
		subAgentCall("b1", "toolu_A", t0.Add(2*time.Minute), 2000, ""),
		subAgentCall("c1", "toolu_B", t0.Add(3*time.Minute), 3000, ""),
		subAgentCall("a2", "", t0.Add(4*time.Minute), 4000, ""),
		streamResult("s", 10000, 10.0),
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(logs, "t-"+story+"-20260929T100000Z.log"), []byte(log), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Measure(lab.o.Dir, lab.root, "t", story, true); err != nil {
		t.Fatal(err)
	}
	// of 10000: the first half of a1, b1, and half of a2; the second half
	// of a1, c1, and half of a2
	f, s := lab.usageOf(first), lab.usageOf(second)
	if f == nil || diff(f.Cost(), 4.5) > 1e-3 || s == nil || diff(s.Cost(), 5.5) > 1e-3 {
		t.Fatalf("tasks = %+v and %+v, want 4.5 and 5.5", f, s)
	}
	if story := lab.usageOf(story); f.Cost()+s.Cost() > story.Cost()+1e-9 {
		t.Errorf("the tasks cost %v, more than the story's %v", f.Cost()+s.Cost(), story.Cost())
	}
}

func diff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

// usageStub is an agent that writes a call to its log, stays while held,
// and ends with its session's totals.
func (lab *agentLab) usageStub() {
	lab.t.Helper()
	stub := filepath.Join(lab.outDir, "usage-agent")
	script := "#!/bin/sh\n" +
		"now=$(date -u +%Y-%m-%dT%H:%M:%S.000Z)\n" +
		`printf '{"type":"assistant","timestamp":"%s","session_id":"s","message":{"id":"m1","model":"claude-opus-5-5","content":[],"usage":{"input_tokens":1,"output_tokens":1,"cache_read_input_tokens":500,"cache_creation_input_tokens":0}}}\n' "$now"` + "\n" +
		"while [ -f \"" + lab.outDir + "/hold\" ] && [ ! -f \"" + lab.outDir + "/release-$FLAI_STORY\" ]; do sleep 0.05; done\n" +
		`echo '{"type":"result","subtype":"success","session_id":"s","modelUsage":{"claude-opus-5-5":{"inputTokens":1,"outputTokens":99,"cacheReadInputTokens":500,"cacheCreationInputTokens":0,"costUSD":0.75}}}'` + "\n" +
		"touch \"" + lab.outDir + "/ended-$$\"\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	lab.cfg.Command = []string{stub}
}

// A task done while its story's agent runs is measured at the next look, and
// once; the story is measured whole when the agent ends.
func TestAnAgentsStoryIsMeasuredWhenATaskIsDoneAndWhenItEnds(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.usageStub()
	lab.hold()
	story := lab.ready("Worked")
	lab.l.look(ctx, false)
	waitFor(t, "the agent runs", func() bool { return lab.run(story).live() })
	lab.move(story, workitem.InProgress)
	waitFor(t, "the agent's call is in its log", func() bool {
		data, _ := os.ReadFile(lab.run(story).Log)
		return strings.Contains(string(data), `"m1"`)
	})
	task := lab.taskWorked(story, "done mid-run", time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	lab.l.look(ctx, false)
	lab.l.look(ctx, false)
	tu := lab.usageOf(task)
	if tu == nil || tu.Tokens() != 502 || !tu.Estimated {
		t.Fatalf("task usage = %+v, want its call, estimated", tu)
	}
	if n := strings.Count(lab.logText(), `msg="usage measured"`); n != 1 {
		t.Errorf("measured %d times while the agent ran, want once", n)
	}
	lab.release(story)
	waitFor(t, "the story and its epic are measured", func() bool {
		u, e := lab.usageOf(story), lab.usageOf(lab.epic.ID)
		return u != nil && !u.Estimated && u.Cost() == 0.75 && e != nil && e.Cost() == 0.75
	})
	if u := lab.usageOf(task); u == nil || u.Cost() != 0.75 || u.Models[0].Output != 99 {
		t.Errorf("task usage after the run = %+v, want the session's totals, all in its window", u)
	}
}

func (lab *agentLab) logText() string {
	lab.logs.mu.Lock()
	defer lab.logs.mu.Unlock()
	return lab.logs.buf.String()
}
