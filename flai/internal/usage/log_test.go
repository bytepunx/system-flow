package usage

import (
	"os"
	"path/filepath"
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

// A window is given each session's reported totals in the share of its
// calls' input and cache tokens that fall in it.
func TestAWindowIsApportioned(t *testing.T) {
	rec, err := Read(writeLog(t, "a.log", run1...))
	if err != nil {
		t.Fatal(err)
	}
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	u := rec.Window(at("2026-09-29T10:05:00Z"), time.Time{}, nil)
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
	first := rec.Window(at("2026-09-29T10:00:00Z"), at("2026-09-29T10:05:00Z"), nil)
	if diff(first.Cost()+u.Cost(), 1.5) > 1e-3 {
		t.Errorf("two windows that cover the run cost %v and %v, want 1.5 between them", first.Cost(), u.Cost())
	}
	if none := rec.Window(at("2026-09-30T00:00:00Z"), time.Time{}, nil); none != nil {
		t.Errorf("a window after the run = %+v, want nothing", none)
	}
}

func TestAMissingLogIsSkipped(t *testing.T) {
	rec, err := Read(filepath.Join(t.TempDir(), "gone.log"))
	if err != nil || rec.Total(nil) != nil {
		t.Fatalf("rec %+v err %v", rec, err)
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
