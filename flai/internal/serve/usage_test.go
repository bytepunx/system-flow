package serve

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
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

// claudeCode is a story's agent that runs claude-code, whose logs are its
// stream-json.
var claudeCode = &manifest.Agent{Harness: harness.ClaudeCode}

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
	story := lab.backlog("Measured", claudeCode)
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

	m, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true)
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
	again, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true)
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
	story := lab.backlog("Planned", claudeCode)
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
	m, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true)
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
	again, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true)
	if err != nil || len(again.Changed) != 0 {
		t.Errorf("measuring again changed %v (%v), want nothing", again.Changed, err)
	}
}

// ADR-0130: the logs of a story whose agent ran the operator's command are
// read by the command's reader, which reads nothing, whatever they hold: the
// story whose agent names no harness, the one that names it, and the one
// that names claude-code whose newest run the command ran.
func TestMeasureReadsNothingOfACommandsLog(t *testing.T) {
	lab := newAgentLab(t)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	logs := filepath.Join(string(lab.o.Dir), "agents")
	_ = os.MkdirAll(logs, 0o700)
	ran := lab.backlog("Ran by the command", claudeCode)
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Story: ran, Harness: harness.Command, Agent: "builder-" + ran, Started: t0.Format(time.RFC3339), Ended: t0.Add(time.Hour).Format(time.RFC3339)})
	})
	for _, story := range []string{lab.backlog("No harness", nil), lab.backlog("The command", &manifest.Agent{Harness: harness.Command}), ran} {
		log := streamCall("s", "m0", t0, 1000) + "\n" + streamResult("s", 1000, 1.0) + "\n"
		if err := os.WriteFile(filepath.Join(logs, "t-"+story+"-20260929T100000Z.log"), []byte(log), 0o600); err != nil {
			t.Fatal(err)
		}
		m, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(m.Logs) != 1 || m.Usage != nil || len(m.Changed) != 0 || lab.usageOf(story) != nil {
			t.Errorf("%s measured %+v, usage %+v; want its log and nothing measured", story, m, lab.usageOf(story))
		}
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
	story := lab.backlog("Planned", claudeCode)
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
	if _, err := Measure(lab.o.Dir, lab.root, "t", story, nil, true); err != nil {
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

// usageStub is claude-code's program: an agent that writes a call to its
// log, stays while held, and ends with its session's totals.
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
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: stub}}
}

// A task done while its story's agent runs is measured at the next look, and
// once; the story is measured whole when the agent ends.
func TestAnAgentsStoryIsMeasuredWhenATaskIsDoneAndWhenItEnds(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.usageStub()
	lab.hold()
	story := lab.readyWith("Worked", claudeCode)
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

// fakeSpend is a gateway's spend reader that answers as answer says and
// keeps what it was asked.
type fakeSpend struct {
	mu     sync.Mutex
	answer func(run usage.RunRef) (*usage.Spent, error)
	asked  []usage.RunRef
}

func (f *fakeSpend) Run(_ context.Context, run usage.RunRef) (*usage.Spent, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.asked = append(f.asked, run)
	return f.answer(run)
}

func (f *fakeSpend) times() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.asked)
}

// fakeKeySpend is a fakeSpend that also reads the key's usage, as
// OpenRouter's reader does.
type fakeKeySpend struct {
	*fakeSpend
	usage float64
}

func (f fakeKeySpend) KeyUsage(context.Context) (float64, error) { return f.usage, nil }

// costs answers cost for every run.
func costs(cost float64) func(usage.RunRef) (*usage.Spent, error) {
	return func(usage.RunRef) (*usage.Spent, error) {
		return &usage.Spent{Cost: cost, By: usage.SpentBySession}, nil
	}
}

// gatewayReads has every gateway read by s for the test, and keeps the key
// each reader was made with.
func gatewayReads(t *testing.T, s usage.Spend) *[]string {
	t.Helper()
	was := spendFor
	var mu sync.Mutex
	keys := &[]string{}
	spendFor = func(_, key string) usage.Spend {
		mu.Lock()
		defer mu.Unlock()
		*keys = append(*keys, key)
		return s
	}
	t.Cleanup(func() { spendFor = was })
	return keys
}

// spendKey is the agent's key to the provider gateway in tests.
const spendKey = "sk-sentinel-S0357-never-logged"

// gateways are the host's providers in tests: gateway, its key in
// S0357_KEY, and admin, whose spend log is read with the key in
// S0357_SPEND_KEY.
func gateways(t *testing.T) map[string]manifest.Provider {
	t.Helper()
	t.Setenv("S0357_KEY", spendKey)
	return map[string]manifest.Provider{
		"gateway": {API: manifest.APIAnthropicMessages, BaseURL: "https://gateway.example.com", KeyEnv: "S0357_KEY"},
		"admin":   {API: manifest.APIAnthropicMessages, BaseURL: "https://admin.example.com", KeyEnv: "S0357_KEY", SpendKeyEnv: "S0357_SPEND_KEY"},
	}
}

// agentLog writes the log of a run of story's agent started at stamp, and
// gives its path.
func (lab *agentLab) agentLog(story, stamp string, lines ...string) string {
	lab.t.Helper()
	dir := filepath.Join(string(lab.o.Dir), "agents")
	_ = os.MkdirAll(dir, 0o700)
	p := filepath.Join(dir, "t-"+story+"-"+stamp+".log")
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		lab.t.Fatal(err)
	}
	return p
}

// throughProvider writes the spend record of the run whose log is at log, as
// its start does.
func throughProvider(t *testing.T, log string, rec spendRecord) {
	t.Helper()
	if err := writeSpend(log, &rec); err != nil {
		t.Fatal(err)
	}
}

func spendOfLog(t *testing.T, log string) (*spendRecord, string) {
	t.Helper()
	r, err := readSpend(log)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(spendPath(log))
	return r, string(data)
}

// modelCall is one call to model as Claude Code streams it.
func modelCall(session, id, model string, at time.Time, read int) string {
	return strings.Replace(streamCall(session, id, at, read), "claude-opus-5-5", model, 1)
}

const haiku = "claude-haiku-4-5"

// S-0357, ADR-0132: a run through a provider is priced by its gateway: its
// cost replaces the harness's, apportioned over the run's models by their
// share of the harness's cost, with the harness's tokens; the story's task
// is scaled with it, both priced by the gateway, and its epic sums it. The
// price is kept beside the log, without the key, and the next measurement
// asks nothing.
func TestAGatewayPricesARunThroughAProvider(t *testing.T) {
	lab := newAgentLab(t)
	providers := gateways(t)
	gw := &fakeSpend{answer: costs(4.0)}
	keys := gatewayReads(t, gw)
	story := lab.backlog("Through the gateway", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	task := lab.taskWorked(story, "worked", t0.Add(-time.Minute), t0.Add(10*time.Minute))
	log := lab.agentLog(story, "20260929T100000Z",
		streamCall("s", "m0", t0, 1000),
		modelCall("s", "m1", haiku, t0.Add(time.Minute), 1000),
		`{"type":"result","subtype":"success","session_id":"s","modelUsage":{"claude-opus-5-5":{"inputTokens":1,"outputTokens":300,"cacheReadInputTokens":1000,"cacheCreationInputTokens":0,"costUSD":1.0},`+
			`"claude-haiku-4-5":{"inputTokens":1,"outputTokens":100,"cacheReadInputTokens":1000,"cacheCreationInputTokens":0,"costUSD":0.25}}}`)
	throughProvider(t, log, spendRecord{Provider: "gateway"})

	m, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Unread) != 0 || gw.times() != 1 {
		t.Fatalf("unread %+v, asked %d times; want the gateway asked once", m.Unread, gw.times())
	}
	sum := sha256.Sum256([]byte(spendKey))
	ref := gw.asked[0]
	if !slices.Equal(ref.Sessions, []string{"s"}) || !ref.From.Equal(t0) || !ref.To.Equal(t0.Add(time.Minute)) || ref.KeyBefore != nil || ref.KeyHash != hex.EncodeToString(sum[:]) {
		t.Errorf("asked for %+v, want the run's session, window, and the hash of its key", ref)
	}
	if !slices.Equal(*keys, []string{spendKey}) {
		t.Errorf("the reader was made with %d keys, want the agent's key alone", len(*keys))
	}
	su := lab.usageOf(story)
	if su == nil || su.PricedBy != usage.PricedGateway || !su.Estimated || diff(su.Cost(), 4.0) > 1e-9 || len(su.Models) != 2 {
		t.Fatalf("story usage = %+v, want the gateway's 4.0 apportioned over two models", su)
	}
	if h, o := su.Models[0], su.Models[1]; h.Model != haiku || diff(h.Cost, 0.8) > 1e-9 || o.Model != "claude-opus-5-5" || diff(o.Cost, 3.2) > 1e-9 || o.Output != 300 || o.CacheRead != 1000 {
		t.Errorf("models = %+v, want haiku 0.8 and opus 3.2 with the harness's tokens", su.Models)
	}
	if tu := lab.usageOf(task); tu == nil || tu.PricedBy != usage.PricedGateway || !tu.Estimated || diff(tu.Cost(), 4.0) > 1e-3 {
		t.Errorf("task usage = %+v, want the story's every call at the gateway's price", tu)
	}
	if eu := lab.usageOf(lab.epic.ID); eu == nil || eu.Priced() != usage.PricedGateway || diff(eu.Cost(), 4.0) > 1e-9 {
		t.Errorf("epic usage = %+v, want its story's, priced by the gateway", eu)
	}
	r, text := spendOfLog(t, log)
	if !r.priced() || *r.Cost != 4.0 || r.By != usage.SpentBySession || r.PricedBy != usage.PricedGateway || r.At == "" || strings.Contains(text, spendKey) {
		t.Errorf("spend record = %s, want the gateway's price and no key", text)
	}

	again, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true)
	if err != nil || len(again.Changed) != 0 || gw.times() != 1 {
		t.Errorf("measuring again changed %v (%v) and asked %d times, want nothing and once", again.Changed, err, gw.times())
	}
}

// S-0357, ADR-0132: a gateway that holds nothing for a run leaves it priced
// by the harness and says so in flai serve's log, without the key, once;
// the next measurement asks again, and prices it when the gateway answers.
func TestAGatewayThatDoesNotAnswerLeavesTheHarnessFigureAndIsAskedAgain(t *testing.T) {
	lab := newAgentLab(t)
	lab.cfg.Providers = gateways(t)
	answer := fmt.Errorf("the gateway's log is empty: %w", usage.ErrNoSpend)
	gw := &fakeSpend{answer: func(usage.RunRef) (*usage.Spent, error) { return nil, answer }}
	gatewayReads(t, gw)
	story := lab.backlog("Unanswered", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	log := lab.agentLog(story, "20260929T100000Z", streamCall("s", "m0", t0, 1000), streamResult("s", 1000, 1.0))
	before := 10.0
	throughProvider(t, log, spendRecord{Provider: "gateway", KeyBefore: &before})

	lab.l.measure(story)
	su := lab.usageOf(story)
	if su == nil || su.PricedBy != usage.PricedHarness || su.Estimated || su.Cost() != 1.0 {
		t.Fatalf("story usage = %+v, want the harness's 1.0", su)
	}
	said := lab.logText()
	if n := strings.Count(said, "provider=gateway"); n != 1 || !strings.Contains(said, "level=WARN") || !strings.Contains(said, "log="+filepath.Base(log)) || !strings.Contains(said, "story="+story) || !strings.Contains(said, "is empty") {
		t.Errorf("flai serve's log says the spend was not read %d times, want once, with the story, log, and reason:\n%s", n, said)
	}
	if strings.Contains(said, spendKey) {
		t.Errorf("flai serve's log holds the key:\n%s", said)
	}
	if r, _ := spendOfLog(t, log); r.priced() {
		t.Errorf("spend record = %+v, want it left unpriced", r)
	}
	if ref := gw.asked[0]; ref.KeyBefore == nil || *ref.KeyBefore != 10.0 {
		t.Errorf("asked for %+v, want the key's usage at the start", ref)
	}

	gw.answer = costs(1.5)
	lab.l.measure(story)
	if gw.times() != 2 {
		t.Errorf("asked %d times, want again at the next measurement", gw.times())
	}
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedGateway || su.Estimated || su.Cost() != 1.5 {
		t.Errorf("story usage = %+v, want the gateway's 1.5 for its one model", su)
	}
	if r, _ := spendOfLog(t, log); !r.priced() || *r.Cost != 1.5 {
		t.Errorf("spend record = %+v, want it priced", r)
	}
}

// S-0357: a provider whose spend key's variable is unset leaves its run
// estimated, as its calls were, naming the variable and asking no gateway.
func TestAMissingSpendKeyLeavesTheRunEstimatedNamingTheVariable(t *testing.T) {
	lab := newAgentLab(t)
	providers := gateways(t)
	t.Setenv("S0357_SPEND_KEY", "")
	gw := &fakeSpend{answer: costs(9)}
	gatewayReads(t, gw)
	story := lab.backlog("Keyless", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	log := lab.agentLog(story, "20260929T100000Z", streamCall("s", "m0", t0, 1000))
	throughProvider(t, log, spendRecord{Provider: "admin"})
	m, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Unread) != 1 || m.Unread[0].Provider != "admin" || !strings.Contains(m.Unread[0].Reason, "S0357_SPEND_KEY") || strings.Contains(m.Unread[0].Reason, spendKey) || gw.times() != 0 {
		t.Errorf("unread = %+v, asked %d times; want the variable named and nothing asked", m.Unread, gw.times())
	}
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedEstimate || !su.Estimated {
		t.Errorf("story usage = %+v, want estimated", su)
	}
}

// S-0357: a provider whose gateway keeps no spend log flai reads, such as
// Anthropic's own API, leaves its run priced by the harness, and nothing is
// said of it, since there is nothing to read again.
func TestAProviderWithNoSpendLogIsPricedByTheHarnessQuietly(t *testing.T) {
	lab := newAgentLab(t)
	lab.cfg.Providers = gateways(t)
	gatewayReads(t, nil)
	story := lab.backlog("Anthropic's own", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	log := lab.agentLog(story, "20260929T100000Z", streamCall("s", "m0", t0, 1000), streamResult("s", 1000, 1.0))
	throughProvider(t, log, spendRecord{Provider: "gateway"})
	lab.l.measure(story)
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedHarness || su.Cost() != 1.0 {
		t.Errorf("story usage = %+v, want the harness's 1.0", su)
	}
	if said := lab.logText(); strings.Contains(said, "level=WARN") {
		t.Errorf("flai serve's log warns of a gateway with no spend log:\n%s", said)
	}
}

// S-0357: a run with no provider is priced by the harness and asks no
// gateway.
func TestARunWithNoProviderIsPricedByTheHarness(t *testing.T) {
	lab := newAgentLab(t)
	gw := &fakeSpend{answer: costs(9)}
	gatewayReads(t, gw)
	story := lab.backlog("No provider", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	lab.agentLog(story, "20260929T100000Z", streamCall("s", "m0", t0, 1000), streamResult("s", 1000, 1.0))
	if _, err := Measure(lab.o.Dir, lab.root, "t", story, gateways(t), true); err != nil {
		t.Fatal(err)
	}
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedHarness || su.Cost() != 1.0 || gw.times() != 0 {
		t.Errorf("story usage = %+v, asked %d times; want the harness's 1.0 and nothing asked", su, gw.times())
	}
}

// S-0357, ADR-0132: a story with a run priced by its gateway, one by the
// harness, and one estimated is the sum of the three, its harness figure
// for the first left out, and estimated, the least certain of them.
func TestAStorysRunsArePricedApartAndSumToTheLeastCertain(t *testing.T) {
	lab := newAgentLab(t)
	providers := gateways(t)
	gw := &fakeSpend{answer: func(run usage.RunRef) (*usage.Spent, error) {
		if !slices.Equal(run.Sessions, []string{"a"}) {
			t.Errorf("asked for %+v, want the gateway's run alone", run)
		}
		return &usage.Spent{Cost: 2.0, By: usage.SpentBySession}, nil
	}}
	gatewayReads(t, gw)
	story := lab.backlog("Mixed", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	through := lab.agentLog(story, "20260929T100000Z", streamCall("a", "a0", t0, 1000), streamResult("a", 1000, 1.0))
	lab.agentLog(story, "20260929T110000Z", streamCall("b", "b0", t0.Add(time.Hour), 1000), streamResult("b", 1000, 0.5))
	lab.agentLog(story, "20260929T120000Z", streamCall("c", "c0", t0.Add(2*time.Hour), 1000))
	throughProvider(t, through, spendRecord{Provider: "gateway"})

	harnessOnly, err := Measure(lab.o.Dir, lab.root, "t", story, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(harnessOnly.Unread) != 1 || gw.times() != 0 {
		t.Fatalf("without the host's providers: unread %+v, asked %d times; want the run unread", harnessOnly.Unread, gw.times())
	}
	m, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true)
	if err != nil {
		t.Fatal(err)
	}
	su := lab.usageOf(story)
	if su == nil || su.PricedBy != usage.PricedEstimate || !su.Estimated || diff(su.Cost(), harnessOnly.Usage.Cost()+1.0) > 1e-3 || su.Tokens() != harnessOnly.Usage.Tokens() {
		t.Errorf("story usage = %+v, want the harness's %v with the first run's 1.0 replaced by the gateway's 2.0, estimated", su, harnessOnly.Usage.Cost())
	}
	if len(m.Unread) != 0 || gw.times() != 1 {
		t.Errorf("unread %+v, asked %d times; want the gateway's run asked once", m.Unread, gw.times())
	}
}

// S-0357: a story's newest run, which flai serve's state says went through a
// provider and which kept no spend record, is not asked for while it goes
// on, and once it has ended is priced by the gateway, its record written.
func TestTheNewestRunIsPricedOnceItEnds(t *testing.T) {
	lab := newAgentLab(t)
	providers := gateways(t)
	gw := &fakeSpend{answer: costs(3.0)}
	gatewayReads(t, gw)
	story := lab.backlog("Going", claudeCode)
	t0 := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	log := lab.agentLog(story, "20260929T100000Z", streamCall("s", "m0", t0, 1000), streamResult("s", 1000, 1.0))
	run := &AgentRun{Story: story, Harness: harness.ClaudeCode, Agent: "builder-" + story, Provider: "gateway", KeyEnv: "S0357_KEY", PID: os.Getpid(), Started: t0.Format(time.RFC3339), Log: log}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) { s.put(run) })
	if _, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true); err != nil {
		t.Fatal(err)
	}
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedHarness || gw.times() != 0 {
		t.Fatalf("while it goes on: usage %+v, asked %d times; want the harness's and nothing asked", su, gw.times())
	}
	ended := *run
	ended.Ended = t0.Add(time.Hour).Format(time.RFC3339)
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) { s.put(&ended) })
	if _, err := Measure(lab.o.Dir, lab.root, "t", story, providers, true); err != nil {
		t.Fatal(err)
	}
	if su := lab.usageOf(story); su == nil || su.PricedBy != usage.PricedGateway || su.Cost() != 3.0 || gw.times() != 1 {
		t.Errorf("once it ended: usage %+v, asked %d times; want the gateway's 3.0", su, gw.times())
	}
	if r, _ := spendOfLog(t, log); !r.priced() || r.Provider != "gateway" {
		t.Errorf("spend record = %+v, want it written priced", r)
	}
}

// S-0357, ADR-0132: a story's agent started through a provider whose reader
// joins by the key's usage keeps a spend record beside its log with the
// provider and the key's usage at its start, and no key; the agent is given
// no spend key's variable.
func TestAStartThroughAProviderKeepsTheKeysUsage(t *testing.T) {
	t.Setenv("S0350_GATEWAY_KEY", spendKey)
	lab := providerLab(t)
	gw := fakeKeySpend{fakeSpend: &fakeSpend{answer: costs(9)}, usage: 12.5}
	keys := gatewayReads(t, gw)
	id := lab.readyWith("Read at the start", &manifest.Agent{Harness: harness.ClaudeCode, Provider: "gateway"})
	lab.l.look(context.Background(), false)
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	r, text := spendOfLog(t, lab.run(id).Log)
	if r == nil || r.Provider != "gateway" || r.KeyBefore == nil || *r.KeyBefore != 12.5 || r.priced() {
		t.Fatalf("spend record = %s, want the provider and the key's usage at the start", text)
	}
	if strings.Contains(text, spendKey) || strings.Contains(text, "S0350_GATEWAY_KEY") {
		t.Errorf("spend record holds the key or its variable:\n%s", text)
	}
	if !slices.Equal(*keys, []string{spendKey}) {
		t.Errorf("the reader was made with %d keys, want the agent's", len(*keys))
	}
	if strings.Contains(lab.logText(), spendKey) {
		t.Errorf("flai serve's log holds the key:\n%s", lab.logText())
	}
}

// S-0357: the agents flai serve starts are not given a variable a provider
// names as its spend key's, unless another names it as its key's.
func TestAgentsAreNotGivenASpendKey(t *testing.T) {
	lab := newAgentLab(t)
	lab.cfg.Providers = map[string]manifest.Provider{
		"admin":  {KeyEnv: "AGENT_KEY", SpendKeyEnv: "MASTER_KEY"},
		"shared": {KeyEnv: "OTHER_KEY", SpendKeyEnv: "AGENT_KEY"},
	}
	got := lab.l.withoutSpendKeys([]string{"PATH=/bin", "MASTER_KEY=m", "AGENT_KEY=a", "OTHER_KEY=o"})
	if !slices.Equal(got, []string{"PATH=/bin", "AGENT_KEY=a", "OTHER_KEY=o"}) {
		t.Errorf("env = %v, want the master key alone left out", got)
	}
}
