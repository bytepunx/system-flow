package serve

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A session as claude -p --output-format stream-json writes it, with the
// events a reader does not follow among those it does.
var sessionLines = []string{
	`{"type":"system","subtype":"init","model":"claude-opus-5-5","session_id":"s"}`,
	`{"type":"system","subtype":"thinking_tokens","estimated_tokens":50}`,
	`{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"","signature":"x"}]}}`,
	`{"type":"assistant","message":{"content":[{"type":"text","text":"Priming the session."}]}}`,
	`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"ls -la","description":"List files"}}]}}`,
	`{"type":"rate_limit_event","rate_limit_info":{"status":"allowed"}}`,
	`{"type":"user","message":{"content":[{"type":"tool_result","content":"a\nb","is_error":false}]}}`,
	`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"/tmp/x.go","limit":5}}]}}`,
	`{"type":"user","message":{"content":[{"type":"tool_result","content":[{"type":"text","text":"no such file"}],"is_error":true}]}}`,
	`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"mcp__flai__inbox","input":{"story":"S-1","all":true}}]}}`,
	`{"type":"system","subtype":"task_started","description":"Run the tests"}`,
	`{"type":"tool_progress","tool_name":"Bash","heartbeat":true}`,
	`{"type":"system","subtype":"task_notification","status":"completed","summary":"Run the tests"}`,
	`{"type":"result","subtype":"success","is_error":false,"num_turns":12,"duration_ms":65000,"total_cost_usd":1.5,"result":"Moved to review."}`,
}

var sessionEntries = []StreamEntry{
	{Kind: StreamSession, Text: "session started (claude-opus-5-5)"},
	{Kind: StreamText, Text: "Priming the session."},
	{Kind: StreamTool, Tool: "Bash", Text: "List files"},
	{Kind: StreamResult, Text: "a\nb"},
	{Kind: StreamTool, Tool: "Read", Text: "/tmp/x.go"},
	{Kind: StreamResult, Text: "no such file", Error: true},
	{Kind: StreamTool, Tool: "mcp__flai__inbox", Text: `all=true story="S-1"`},
	{Kind: StreamTask, Text: "started: Run the tests"},
	{Kind: StreamTask, Text: "completed: Run the tests"},
	{Kind: StreamEnd, Text: "session ended (success, 12 turns, 65s, $1.50): Moved to review."},
}

// streamState is a project's agents with one run for S-0001 whose log holds
// text, running or not.
func streamState(t *testing.T, text string, running bool) (AgentState, string) {
	t.Helper()
	log := filepath.Join(t.TempDir(), "sf-S-0001.log")
	if err := os.WriteFile(log, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	run := &AgentRun{Story: "S-0001", Agent: "agent-S-0001", Started: "2026-09-29T05:00:00Z", Log: log, PID: 42}
	if !running {
		run.Ended, run.Outcome = "2026-09-29T05:10:00Z", OutcomeWorked
	}
	return AgentState{Stories: map[string]*AgentRun{"S-0001": run}}, log
}

func sameEntries(t *testing.T, got, want []StreamEntry) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d:\n%+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestStreamReadsASessionFromTheStart(t *testing.T) {
	text := strings.Join(sessionLines, "\n") + "\n"
	st, _ := streamState(t, text, false)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, sessionEntries)
	if got.From != 0 || got.Next != int64(len(text)) || got.Size != int64(len(text)) || got.More {
		t.Errorf("from %d next %d size %d more %v, want 0, %d, %d, false", got.From, got.Next, got.Size, got.More, len(text), len(text))
	}
	if got.Running || got.Agent != "agent-S-0001" || got.Started != "2026-09-29T05:00:00Z" || got.Outcome != OutcomeWorked {
		t.Errorf("run = %+v", got)
	}
}

func TestStreamGoesOnFromAnOffsetAndLeavesAPartialLine(t *testing.T) {
	first := strings.Join(sessionLines[:4], "\n") + "\n"
	st, log := streamState(t, first, true)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, sessionEntries[:2])
	if !got.Running || got.Next != int64(len(first)) {
		t.Fatalf("running %v next %d, want true, %d", got.Running, got.Next, len(first))
	}
	// the agent writes one line and half of the next
	half := sessionLines[6][:20]
	f, err := os.OpenFile(log, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = fmt.Fprint(f, sessionLines[4]+"\n"+half)
	got, err = Stream(st, "S-0001", got.Next)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, sessionEntries[2:3])
	if got.From != int64(len(first)) || got.Next != int64(len(first)+len(sessionLines[4])+1) || !got.More {
		t.Errorf("from %d next %d more %v", got.From, got.Next, got.More)
	}
	// the rest of the line arrives
	_, _ = fmt.Fprint(f, sessionLines[6][20:]+"\n")
	_ = f.Close()
	got, err = Stream(st, "S-0001", got.Next)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, sessionEntries[3:4])
	if got.More {
		t.Error("more after reading to the end")
	}
}

func TestStreamReadsTheTailOfALongLog(t *testing.T) {
	line := `{"type":"assistant","message":{"content":[{"type":"text","text":"` + strings.Repeat("x", 50) + `"}]}}`
	var b strings.Builder
	for b.Len() < 3*StreamTailBytes {
		b.WriteString(line + "\n")
	}
	st, _ := streamState(t, b.String(), true)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	size := int64(b.Len())
	if got.From != size-StreamTailBytes || got.Next != size {
		t.Errorf("from %d next %d, want %d, %d", got.From, got.Next, size-StreamTailBytes, size)
	}
	if len(got.Entries) != StreamEntries || got.Skipped == 0 {
		t.Errorf("%d entries, %d skipped; want %d and some skipped", len(got.Entries), got.Skipped, StreamEntries)
	}
	for _, e := range got.Entries {
		if e.Kind != StreamText || e.Text != strings.Repeat("x", 50) {
			t.Fatalf("entry %+v: want every line whole, the first after the tail's start included", e)
		}
	}
	// an offset past the end, from a log since replaced, reads the tail too
	again, err := Stream(st, "S-0001", size+10)
	if err != nil {
		t.Fatal(err)
	}
	if again.From != got.From {
		t.Errorf("from %d past the end, want the tail at %d", again.From, got.From)
	}
}

func TestStreamReadsAMegabyteAtATime(t *testing.T) {
	line := `{"type":"user","message":{"content":[{"type":"tool_result","content":"` + strings.Repeat("y", 64<<10) + `"}]}}`
	var b strings.Builder
	for b.Len() < 2*streamReadBytes {
		b.WriteString(line + "\n")
	}
	st, _ := streamState(t, b.String(), true)
	got, err := Stream(st, "S-0001", 0)
	if err != nil {
		t.Fatal(err)
	}
	if !got.More || got.Next-got.From < streamReadBytes || got.Next-got.From > streamReadBytes+int64(len(line))+1 {
		t.Errorf("read %d bytes, more %v; want about %d and more", got.Next-got.From, got.More, streamReadBytes)
	}
}

func TestStreamShowsAnOverlongLineByItsSize(t *testing.T) {
	long := strings.Repeat("z", streamLineBytes+10)
	text := long + "\n" + sessionLines[3] + "\n"
	st, _ := streamState(t, text, true)
	got, err := Stream(st, "S-0001", 0)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, []StreamEntry{{Kind: StreamOutput, Text: fmt.Sprintf("(a line of %d bytes, too long to show)", len(long)+1)}})
	if !got.More || got.Next != int64(len(long)+1) {
		t.Fatalf("next %d more %v: the line fills one read, want %d and more", got.Next, got.More, len(long)+1)
	}
	got, err = Stream(st, "S-0001", got.Next)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, sessionEntries[1:2])
}

func TestStreamShowsOtherOutputAsItIs(t *testing.T) {
	st, _ := streamState(t, "starting agent for S-0001\r\n\n{not json\nError: no API key\n", false)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, []StreamEntry{
		{Kind: StreamOutput, Text: "starting agent for S-0001"},
		{Kind: StreamOutput, Text: "{not json"},
		{Kind: StreamOutput, Text: "Error: no API key"},
	})
}

func TestStreamCutsALongText(t *testing.T) {
	st, _ := streamState(t, `{"type":"assistant","message":{"content":[{"type":"text","text":"`+strings.Repeat("é", streamTextRunes+5)+`"}]}}`+"\n", true)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, []StreamEntry{{Kind: StreamText, Text: strings.Repeat("é", streamTextRunes) + "…"}})
}

func TestStreamOfAStoryWithNoAgent(t *testing.T) {
	st, _ := streamState(t, "", true)
	if _, err := Stream(st, "S-0002", -1); !errors.Is(err, ErrNoRun) {
		t.Errorf("err = %v, want ErrNoRun", err)
	}
}

func TestStreamOfARunThatHasNoLog(t *testing.T) {
	st := AgentState{Stories: map[string]*AgentRun{
		"S-0001": {Story: "S-0001", Started: "2026-09-29T05:00:00Z", Ended: "2026-09-29T05:00:00Z", Error: "claude: executable file not found in $PATH", Outcome: OutcomeFailed},
		"S-0002": {Story: "S-0002", Started: "2026-09-29T05:00:00Z", PID: 7, Log: filepath.Join(t.TempDir(), "gone.log")},
	}}
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, []StreamEntry{{Kind: StreamEnd, Text: "could not be started: claude: executable file not found in $PATH", Error: true}})
	if got.Running {
		t.Error("a run that could not be started reads as running")
	}
	got, err = Stream(st, "S-0002", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Entries) != 0 || !got.Running || got.Size != 0 {
		t.Errorf("a log not written yet: %+v", got)
	}
}

func TestStreamSaysAFailedSessionEnd(t *testing.T) {
	st, _ := streamState(t, `{"type":"result","subtype":"error_max_turns","is_error":true,"num_turns":3}`+"\n", false)
	got, err := Stream(st, "S-0001", -1)
	if err != nil {
		t.Fatal(err)
	}
	sameEntries(t, got.Entries, []StreamEntry{{Kind: StreamEnd, Text: "session ended (error_max_turns, 3 turns)", Error: true}})
}
