package usage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// writeLog writes lines as one agent log and returns its path.
func writeLog(t *testing.T, name string, lines ...string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// assistant is one content block of a call, as Claude Code streams it.
func assistant(session, id, model, at string, in, out, read, write int) string {
	return `{"type":"assistant","timestamp":"` + at + `","session_id":"` + session + `","message":{"id":"` + id + `","model":"` + model +
		`","content":[{"type":"text","text":"x"}],"usage":{"input_tokens":` + itoa(in) + `,"output_tokens":` + itoa(out) +
		`,"cache_read_input_tokens":` + itoa(read) + `,"cache_creation_input_tokens":` + itoa(write) + `}}}`
}

func user(session, at string) string {
	return `{"type":"user","timestamp":"` + at + `","session_id":"` + session + `","message":{"content":[{"type":"tool_result","content":"ok"}]}}`
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

const opus, haiku = "claude-opus-5-5", "claude-haiku-4-5"

// A run of session s: two calls of opus, the first streamed as two blocks,
// and a subagent's call of haiku; the result reports more output than the
// stream shows, and the cost.
var run1 = []string{
	`{"type":"system","subtype":"init","model":"claude-opus-5-5","session_id":"s"}`,
	assistant("s", "m1", opus, "2026-09-29T10:00:00Z", 10, 1, 1000, 500),
	assistant("s", "m1", opus, "2026-09-29T10:00:01Z", 10, 7, 1000, 500),
	user("s", "2026-09-29T10:00:05Z"),
	`not json: a line the harness wrote to stderr`,
	assistant("s", "m2", haiku, "2026-09-29T10:01:00Z", 5, 1, 200, 100),
	`{"type":"system","subtype":"thinking_tokens","estimated_tokens":50}`,
	assistant("s", "m3", opus, "2026-09-29T10:10:00Z", 2, 1, 3000, 0),
	`{"type":"result","subtype":"success","session_id":"s","total_cost_usd":1.5,"modelUsage":{` +
		`"claude-opus-5-5":{"inputTokens":12,"outputTokens":400,"cacheReadInputTokens":4000,"cacheCreationInputTokens":500,"costUSD":1.2},` +
		`"claude-haiku-4-5":{"inputTokens":5,"outputTokens":50,"cacheReadInputTokens":200,"cacheCreationInputTokens":100,"costUSD":0.3}}}`,
}

func TestTotalOfOneRunIsWhatItsResultReports(t *testing.T) {
	rec, err := Read(writeLog(t, "a.log", run1...))
	if err != nil {
		t.Fatal(err)
	}
	u := rec.Total(nil)
	if u == nil || u.Source != SourceLog || u.Estimated {
		t.Fatalf("usage = %+v, want measured and reported", u)
	}
	want := []Model{
		{Model: haiku, Input: 5, Output: 50, CacheRead: 200, CacheWrite: 100, Cost: 0.3},
		{Model: opus, Input: 12, Output: 400, CacheRead: 4000, CacheWrite: 500, Cost: 1.2},
	}
	if len(u.Models) != 2 || u.Models[0] != want[0] || u.Models[1] != want[1] {
		t.Fatalf("models = %+v, want %+v", u.Models, want)
	}
	if u.Seconds != 600 {
		t.Errorf("seconds = %d, want 600 (first event to last)", u.Seconds)
	}
	if u.Tokens() != 5267 || u.Cost() != 1.5 {
		t.Errorf("tokens %d cost %v", u.Tokens(), u.Cost())
	}
}

// A session resumed in a second run reports its totals again, cumulative:
// the newest result is the session's total, not an addition to it.
func TestAResumedSessionCountsItsNewestResultOnly(t *testing.T) {
	first := writeLog(t, "a.log", run1...)
	second := writeLog(t, "b.log",
		assistant("s", "m4", opus, "2026-09-29T11:00:00Z", 1, 1, 5000, 0),
		`{"type":"result","subtype":"success","session_id":"s","modelUsage":{`+
			`"claude-opus-5-5":{"inputTokens":13,"outputTokens":500,"cacheReadInputTokens":9000,"cacheCreationInputTokens":500,"costUSD":1.7},`+
			`"claude-haiku-4-5":{"inputTokens":5,"outputTokens":50,"cacheReadInputTokens":200,"cacheCreationInputTokens":100,"costUSD":0.3}}}`)
	rec, err := Read(first, second)
	if err != nil {
		t.Fatal(err)
	}
	u := rec.Total(nil)
	if u.Cost() != 2.0 || u.Estimated {
		t.Fatalf("cost = %v estimated %v, want 2.0 as reported", u.Cost(), u.Estimated)
	}
	if len(rec.Runs) != 2 || u.Seconds != 600 {
		t.Errorf("runs %d seconds %d, want two runs and 600 s (the second has one event)", len(rec.Runs), u.Seconds)
	}
}

// A run with no result, still going or dead, is counted from its calls and
// its cost estimated at the rate other sessions reported.
func TestARunWithNoResultIsEstimated(t *testing.T) {
	first := writeLog(t, "a.log", run1...)
	second := writeLog(t, "b.log",
		`{"type":"system","subtype":"init","model":"claude-opus-5-5","session_id":"t"}`,
		assistant("t", "n1", opus, "2026-09-29T12:00:00Z", 2, 2, 4512, 0),
		assistant("t", "n2", "claude-fable-5-1", "2026-09-29T12:00:30Z", 1, 1, 100, 0))
	rec, err := Read(first, second)
	if err != nil {
		t.Fatal(err)
	}
	u := rec.Total(nil)
	if !u.Estimated {
		t.Fatal("a run no result reported is not marked estimated")
	}
	var got Model
	var fable Model
	for _, m := range u.Models {
		switch m.Model {
		case opus:
			got = m
		case "claude-fable-5-1":
			fable = m
		}
	}
	// opus reported 1.2 for 4912 tokens; t's call is 4516 tokens
	want := 1.2 + 4516*1.2/4912
	if got.Input != 14 || got.CacheRead != 8512 || diff(got.Cost, want) > 1e-4 {
		t.Errorf("opus = %+v, want cost %.4f", got, want)
	}
	if fable.Tokens() != 102 || fable.Cost != 0 {
		t.Errorf("fable = %+v: tokens counted, no rate to price them at", fable)
	}
	if priced := rec.Total(Rates{"claude-fable-5-1": 0.01}); diff(priced.Cost()-u.Cost(), 1.02) > 1e-4 {
		t.Errorf("a rate given for fable did not price it: %v", priced.Cost()-u.Cost())
	}
}

// window is what a task alone, in progress from from until to, is given.
func window(rec *Record, from, to time.Time) *Usage {
	return rec.Tasks(map[string][]Span{"T-0001": {{from, to}}}, nil)["T-0001"]
}

// A task worked alone is given each session's reported totals in the share
// of its calls' input and cache tokens that fall in its window.
func TestAWindowIsApportioned(t *testing.T) {
	rec, err := Read(writeLog(t, "a.log", run1...))
	if err != nil {
		t.Fatal(err)
	}
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	u := window(rec, at("2026-09-29T10:05:00Z"), time.Time{})
	if u == nil || !u.Estimated {
		t.Fatalf("window = %+v, want apportioned and estimated", u)
	}
	// opus: m3 weighs 3002 of 4512
	if len(u.Models) != 1 || u.Models[0].Model != opus || diff(u.Models[0].Cost, 1.2*3002/4512) > 1e-4 || u.Models[0].Output != 266 {
		t.Errorf("models = %+v", u.Models)
	}
	if u.Seconds != 300 {
		t.Errorf("seconds = %d, want the 300 s of the run in the window", u.Seconds)
	}
	first := window(rec, at("2026-09-29T10:00:00Z"), at("2026-09-29T10:05:00Z"))
	if diff(first.Cost()+u.Cost(), 1.5) > 1e-3 {
		t.Errorf("two windows that cover the run cost %v and %v, want 1.5 between them", first.Cost(), u.Cost())
	}
	if none := window(rec, at("2026-09-30T00:00:00Z"), time.Time{}); none != nil {
		t.Errorf("a window after the run = %+v, want nothing", none)
	}
}

func TestAMissingLogIsSkipped(t *testing.T) {
	rec, err := Read(filepath.Join(t.TempDir(), "gone.log"))
	if err != nil || rec.Total(nil) != nil {
		t.Fatalf("rec %+v err %v", rec, err)
	}
}

// Each run keeps the final text of its newest result; a run with none, or
// whose newest result carries no text, keeps nothing.
func TestARunKeepsItsNewestResultsText(t *testing.T) {
	answered := writeLog(t, "a.log",
		assistant("s", "m1", opus, "2026-09-29T10:00:00Z", 1, 1, 10, 0),
		`{"type":"result","subtype":"success","session_id":"s","result":"first answer","total_cost_usd":0.1}`,
		assistant("s", "m2", opus, "2026-09-29T10:01:00Z", 1, 1, 10, 0),
		`{"type":"result","subtype":"success","session_id":"s","result":"Planned S-0001.\nThen more.","total_cost_usd":0.2}`)
	silent := writeLog(t, "b.log",
		assistant("s", "m3", opus, "2026-09-29T11:00:00Z", 1, 1, 10, 0),
		`{"type":"result","subtype":"success","session_id":"s","result":"an answer"}`,
		`{"type":"result","subtype":"error_max_turns","session_id":"s"}`)
	going := writeLog(t, "c.log", assistant("s", "m4", opus, "2026-09-29T12:00:00Z", 1, 1, 10, 0))
	rec, err := Read(answered, silent, going)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Runs) != 3 {
		t.Fatalf("runs = %+v, want three", rec.Runs)
	}
	for i, want := range []string{"Planned S-0001.\nThen more.", "", ""} {
		if rec.Runs[i].Result != want {
			t.Errorf("run %d result = %q, want %q", i, rec.Runs[i].Result, want)
		}
	}
}

func TestSumAddsModelByModel(t *testing.T) {
	a := &Usage{Source: SourceLog, Seconds: 10, Models: []Model{{Model: opus, Input: 1, Cost: 0.5}}}
	b := &Usage{Source: SourceLog, Seconds: 5, Estimated: true, Models: []Model{{Model: haiku, Output: 3, Cost: 0.25}, {Model: opus, Input: 2, Cost: 0.25}}}
	s := Sum(a, nil, b)
	if s.Source != SourceSum || s.Seconds != 15 || !s.Estimated || len(s.Models) != 2 || s.Models[1].Input != 3 || s.Cost() != 1.0 {
		t.Fatalf("sum = %+v", s)
	}
	if Sum(nil, &Usage{}) != nil {
		t.Error("the sum of nothing is not nil")
	}
	if !Same(s, Sum(b, a)) || Same(s, a) {
		t.Error("Same does not compare what they say")
	}
}

func diff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestSummaryAndCount(t *testing.T) {
	u := &Usage{Source: SourceSum, Seconds: 3725, Estimated: true, Models: []Model{{Model: opus, Input: 256, Output: 89342, CacheRead: 19723140, CacheWrite: 327605, Cost: 8.1258}}}
	if got := u.Summary(); got != "20.1M tokens · $8.13 (estimated) · 1h2m5s of agent work · summed from its children" {
		t.Errorf("summary = %q", got)
	}
	if got := u.Models[0].String(); got != "claude-opus-5-5  input 256 · output 89.3K · cache read 19.7M · cache write 327.6K · $8.1258" {
		t.Errorf("model = %q", got)
	}
}

func TestWindowsAddUpAndRatesComeFromResultsAlone(t *testing.T) {
	p := writeLog(t, "a.log", run1...)
	rec, err := Read(p)
	if err != nil {
		t.Fatal(err)
	}
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	both := rec.Tasks(map[string][]Span{"T-0001": {{at("2026-09-29T10:00:00Z"), at("2026-09-29T10:05:00Z")}, {From: at("2026-09-29T10:05:00Z")}}}, nil)["T-0001"]
	if both.Source != SourceLog || diff(both.Cost(), 1.5) > 1e-3 || both.Seconds != 600 {
		t.Errorf("windows = %+v", both)
	}
	rates, err := ReadRates(p, filepath.Join(t.TempDir(), "gone.log"))
	if err != nil {
		t.Fatal(err)
	}
	if diff(rates[opus], 1.2/4912) > 1e-12 || diff(rates[haiku], 0.3/355) > 1e-12 {
		t.Errorf("rates = %v", rates)
	}
}

// opusCall is one call of opus in session s at at, that reads read tokens of
// cache: the story's agent's when parent is empty, else that of the
// sub-agent parent started, which carries described as its description.
// blocks are its content.
func opusCall(id, at string, read int, parent, described, blocks string) string {
	p := `null`
	if parent != "" {
		p = `"` + parent + `"`
	}
	return `{"type":"assistant","timestamp":"` + at + `","session_id":"s","parent_tool_use_id":` + p + `,"task_description":"` + described +
		`","message":{"id":"` + id + `","model":"` + opus + `","content":[` + blocks + `],"usage":{"input_tokens":0,"output_tokens":0,` +
		`"cache_read_input_tokens":` + itoa(read) + `,"cache_creation_input_tokens":0}}}`
}

// agentCall is an Agent tool_use block that starts sub-agent id.
func agentCall(id, description, prompt string) string {
	return `{"type":"tool_use","id":"` + id + `","name":"Agent","input":{"description":"` + description + `","subagent_type":"general-purpose","prompt":"` + prompt + `"}}`
}

// result reports session s's opus totals: read tokens of cache, a tenth of
// them of output, at a dollar per thousand read.
func result(read int) string {
	return `{"type":"result","subtype":"success","session_id":"s","modelUsage":{"claude-opus-5-5":{"inputTokens":0,"outputTokens":` + itoa(read/10) +
		`,"cacheReadInputTokens":` + itoa(read) + `,"cacheCreationInputTokens":0,"costUSD":` + strconv.FormatFloat(float64(read)/1000, 'f', -1, 64) + `}}}`
}

func span(from, to string) Span {
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	return Span{at(from), at(to)}
}

// Two tasks in progress at once, each worked by its own sub-agent, are each
// given their sub-agent's calls, wherever they fall, and half of the story's
// agent's calls while both were in progress; together no more than the story.
func TestTasksInProgressAtOnceShareTheStorysAgentAndKeepTheirSubAgents(t *testing.T) {
	rec, err := Read(writeLog(t, "a.log",
		opusCall("a0", "2026-09-29T10:00:00Z", 1000, "", "", ""),
		// one message starts both, a block a line; T-2 is the second task
		// unpadded, and the description wins over a prompt naming another
		opusCall("a1", "2026-09-29T10:01:00Z", 1000, "", "", agentCall("toolu_A", "Task sub-agent: T-0001 build", "Work T-0001.")),
		opusCall("a1", "2026-09-29T10:01:00Z", 1000, "", "", agentCall("toolu_B", "Task sub-agent: T-2 docs", "Sibling T-0001 is not yours.")),
		user("s", "2026-09-29T10:01:30Z"),
		opusCall("b1", "2026-09-29T10:02:00Z", 2000, "toolu_A", "Task sub-agent: T-0001 build", ""),
		opusCall("c1", "2026-09-29T10:03:00Z", 3000, "toolu_B", "Task sub-agent: T-2 docs", ""),
		opusCall("a2", "2026-09-29T10:04:00Z", 4000, "", "", ""),
		opusCall("a3", "2026-09-29T10:15:00Z", 600, "", "", ""),
		opusCall("b2", "2026-09-29T10:30:00Z", 500, "toolu_A", "Task sub-agent: T-0001 build", ""),
		result(12100)))
	if err != nil {
		t.Fatal(err)
	}
	got := rec.Tasks(map[string][]Span{
		"T-0001": {span("2026-09-29T10:01:00Z", "2026-09-29T10:10:00Z")},
		"T-0002": {span("2026-09-29T10:01:00Z", "2026-09-29T10:20:00Z")},
		"T-0003": {span("2026-09-29T11:00:00Z", "2026-09-29T11:10:00Z")},
	}, nil)
	// of 12100: T-0001 half of a1, b1, half of a2, and b2 after its window;
	// T-0002 half of a1, c1, half of a2, and a3; a0 is no task's
	for id, want := range map[string]float64{"T-0001": 500 + 2000 + 2000 + 500, "T-0002": 500 + 3000 + 2000 + 600} {
		u := got[id]
		if u == nil || !u.Estimated || u.Source != SourceLog || len(u.Models) != 1 {
			t.Fatalf("%s = %+v, want measured and estimated", id, u)
		}
		if m := u.Models[0]; m.CacheRead != int64(want) || m.Output != int64(want/10) || diff(m.Cost, want/1000) > 1e-4 {
			t.Errorf("%s = %+v, want %v read", id, m, want)
		}
	}
	if u, ok := got["T-0003"]; ok {
		t.Errorf("a task in progress after the run = %+v, want absent", u)
	}
	if sum, story := got["T-0001"].Cost()+got["T-0002"].Cost(), rec.Total(nil).Cost(); sum > story {
		t.Errorf("the tasks cost %v, more than the story's %v", sum, story)
	}
	if got["T-0001"].Seconds != 540 || got["T-0002"].Seconds != 1140 {
		t.Errorf("seconds = %d and %d, want each window's time while the run went on", got["T-0001"].Seconds, got["T-0002"].Seconds)
	}
}

// A sub-agent whose description names no task is the task's its prompt names
// first; one neither names a task of is shared as the story's agent's calls
// are; and one whose start was not seen is named by the description its
// calls carry. A call no result covers is estimated for its sub-agent's task.
func TestASubAgentIsNamedByItsPromptOrForNoTask(t *testing.T) {
	first := writeLog(t, "a.log",
		opusCall("a1", "2026-09-29T10:01:00Z", 100, "", "", agentCall("toolu_P", "Implement the parser", "You work T-0004 of S-0001; sibling T-0003 waits.")+","+
			agentCall("toolu_N", "Explore the code", "Read S-0001 and T-0099.")),
		opusCall("p1", "2026-09-29T10:02:00Z", 1000, "toolu_P", "Implement the parser", ""),
		opusCall("n1", "2026-09-29T10:05:00Z", 2000, "toolu_N", "Explore the code", ""),
		opusCall("d1", "2026-09-29T10:06:00Z", 400, "toolu_D", "Task sub-agent: T-0003 docs", ""),
		result(3500))
	second := writeLog(t, "b.log", opusCall("p2", "2026-09-29T10:08:00Z", 700, "toolu_P", "Implement the parser", ""))
	rec, err := Read(first, second)
	if err != nil {
		t.Fatal(err)
	}
	both := []Span{span("2026-09-29T10:00:00Z", "2026-09-29T10:10:00Z")}
	got := rec.Tasks(map[string][]Span{"T-0003": both, "T-0004": both}, nil)
	rate := 3.5 / 3850
	// of 3500: T-0003 d1, half of n1 and of a1; T-0004 p1, half of n1 and
	// of a1, and p2 estimated at the rate the result reports
	for id, want := range map[string]struct {
		read int64
		cost float64
	}{"T-0003": {400 + 1000 + 50, 1.45}, "T-0004": {1000 + 1000 + 50 + 700, 2.05 + 700*rate}} {
		u := got[id]
		if u == nil || len(u.Models) != 1 || !u.Estimated {
			t.Fatalf("%s = %+v", id, u)
		}
		if m := u.Models[0]; m.CacheRead != want.read || diff(m.Cost, want.cost) > 1e-4 {
			t.Errorf("%s = %+v, want %d read for %.4f", id, m, want.read, want.cost)
		}
	}
}
