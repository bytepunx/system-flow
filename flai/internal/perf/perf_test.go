package perf

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"
)

// clock moves only when told, so that durations are exact.
type clock struct{ t time.Time }

func (c *clock) now() time.Time       { return c.t }
func (c *clock) step(d time.Duration) { c.t = c.t.Add(d) }
func newClock() *clock                { return &clock{t: time.Date(2026, 9, 29, 7, 0, 0, 0, time.UTC)} }

// events are the JSON events logged to b, one per line.
func events(t *testing.T, b *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(b.Bytes()), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var e map[string]any
		if err := json.Unmarshal(line, &e); err != nil {
			t.Fatalf("not a JSON event: %s", line)
		}
		out = append(out, e)
	}
	return out
}

func TestPhasesAreSummedCountedAndLongestFirst(t *testing.T) {
	c := newClock()
	ctx, rec := StartAt(context.Background(), c.now)
	for range 3 {
		done := Track(ctx, "exec.git")
		c.step(10 * time.Millisecond)
		done()
	}
	done := Track(ctx, "items.list")
	c.step(182400 * time.Microsecond)
	done()
	if got := rec.Elapsed(); got != 212400*time.Microsecond {
		t.Errorf("elapsed %s", got)
	}
	ph := rec.Phases()
	if len(ph) != 2 || ph[0].Name != "items.list" || ph[1].Name != "exec.git" || ph[1].Count != 3 || ph[1].Total != 30*time.Millisecond {
		t.Errorf("phases %+v", ph)
	}
	if got := rec.String(); got != "items.list=182.4 exec.git=30.0x3" {
		t.Errorf("string %q", got)
	}
}

func TestWithoutARecorderNothingIsRecorded(t *testing.T) {
	Track(context.Background(), "anything")()
	var r *Recorder
	r.Add("anything", time.Second) // a nil recorder is a no-op, not a panic
	if From(context.Background()) != nil {
		t.Error("a bare context has no recorder")
	}
}

func TestLogLevelFollowsTheThresholdAndWaitsAreNeverSlow(t *testing.T) {
	var b bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&b, &slog.HandlerOptions{Level: slog.LevelDebug}))
	for _, tc := range []struct {
		took  time.Duration
		wait  bool
		level string
	}{
		{100 * time.Millisecond, false, "DEBUG"},
		{500 * time.Millisecond, false, "INFO"},
		{2 * time.Minute, true, "DEBUG"},
	} {
		b.Reset()
		c := newClock()
		ctx, rec := StartAt(context.Background(), c.now)
		done := Track(ctx, "board.load")
		c.step(tc.took)
		done()
		rec.Log(logger, Answered{Transport: "channel", Method: "board.get", Project: "harbour", Bytes: 4060, Wait: tc.wait}, 500*time.Millisecond)
		ev := events(t, &b)
		if len(ev) != 1 {
			t.Fatalf("%s: %d events", tc.took, len(ev))
		}
		e := ev[0]
		if e["level"] != tc.level || e["msg"] != "request answered" || e["component"] != "perf" || e["transport"] != "channel" ||
			e["method"] != "board.get" || e["project"] != "harbour" || e["bytes"] != float64(4060) || e["duration_ms"] != float64(tc.took.Milliseconds()) {
			t.Errorf("%s: event %v", tc.took, e)
		}
		if e["phases"] == nil || e["err"] != nil {
			t.Errorf("%s: phases or err %v", tc.took, e)
		}
	}
}

func TestADebugEventIsNotBuiltWhenDebugIsOff(t *testing.T) {
	var b bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&b, nil))
	_, rec := Start(context.Background())
	rec.Log(logger, Answered{Transport: "mcp", Method: "inbox"}, time.Hour)
	if b.Len() != 0 {
		t.Errorf("logged at info: %s", b.String())
	}
}

func TestSlowReadsItsVariable(t *testing.T) {
	for v, want := range map[string]time.Duration{"": DefaultSlow, "2s": 2 * time.Second, "250ms": 250 * time.Millisecond, "soon": DefaultSlow, "-1s": DefaultSlow} {
		t.Setenv(SlowEnv, v)
		if got := Slow(); got != want {
			t.Errorf("%s=%q: %s, want %s", SlowEnv, v, got, want)
		}
	}
}

func TestExecNamesWhatRanNotOnWhat(t *testing.T) {
	for _, tc := range []struct {
		program string
		args    []string
		want    string
	}{
		{"git", []string{"-C", "/repo", "log", "--oneline", "main..story/S-0001"}, "exec.git.log"},
		{"git", []string{"-c", "core.quotepath=off", "status", "--porcelain"}, "exec.git.status"},
		{"git", []string{"rev-parse", "HEAD"}, "exec.git.rev-parse"},
		{"flai", []string{"stream", "diff", "S-0001", "--json"}, "exec.flai.stream.diff"},
		{"flai", []string{"show", "S-0001", "--json"}, "exec.flai.show"},
		{"flai", []string{"--json", "stats"}, "exec.flai.stats"},
		{"git", nil, "exec.git"},
	} {
		if got := Exec(tc.program, tc.args); got != tc.want {
			t.Errorf("%s %v: %q, want %q", tc.program, tc.args, got, tc.want)
		}
	}
}
