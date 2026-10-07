package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// An agent's stream (S-0142): what its session wrote to the log flai serve
// gave it, read from a byte offset and cut into short entries, so that the
// dashboard can show an agent at work and flai serve agent stream can follow
// it. Claude Code writes stream-json, one event per line; any other line is
// shown as it is. Nothing is kept: every read goes to the log.

// Kinds of stream entries.
const (
	StreamSession  = "session"  // the session started
	StreamText     = "text"     // the agent said something
	StreamThinking = "thinking" // the agent's thinking, when the harness shows it
	StreamTool     = "tool"     // the agent called a tool
	StreamResult   = "result"   // a tool answered
	StreamTask     = "task"     // a background task started or ended
	StreamEnd      = "end"      // the session ended
	StreamOutput   = "output"   // a line that is not a stream-json event
)

// Bounds on one read of a stream.
const (
	// StreamTailBytes is how far from the end a read with no offset starts.
	StreamTailBytes = 256 << 10
	// streamReadBytes is how much of the log one read consumes at most.
	streamReadBytes = 1 << 20
	// streamLineBytes is the longest line that is parsed; a longer one is
	// consumed and shown only by its size.
	streamLineBytes = 1 << 20
	// StreamEntries is how many entries one read returns at most: the newest.
	StreamEntries = 200
	// Characters kept of one entry's text, by kind.
	streamTextRunes   = 2000
	streamResultRunes = 400
	streamToolRunes   = 300
)

// StreamEntry is one thing the agent did or said.
type StreamEntry struct {
	Kind string `json:"kind"`
	// Tool names the tool of a call.
	Tool string `json:"tool,omitempty"`
	Text string `json:"text"`
	// Error marks a tool result or a session end that failed.
	Error bool `json:"error,omitempty"`
}

// StreamRead is one read of a story's newest agent's stream, of an item's
// newest planner's, or of the project's newest orchestrator's or analyzer's.
type StreamRead struct {
	// Story is the story whose agent this is; empty on any other stream.
	Story string `json:"story"`
	// Item is the epic or story a planner plans; empty on any other stream.
	Item string `json:"item,omitempty"`
	// Role is the strategic role read by RoleStream, orchestrate or analyze;
	// empty on a story's agent's stream and a planner's.
	Role  string `json:"role,omitempty"`
	Agent string `json:"agent"`
	// Started names the run: a new run writes a new log, so an offset from
	// another run means nothing.
	Started string `json:"started"`
	Ended   string `json:"ended,omitempty"`
	Running bool   `json:"running"`
	Outcome string `json:"outcome,omitempty"`
	// From is where this read began and Next where the next one should; the
	// log is Size bytes long, and More says the read stopped before its end.
	From int64 `json:"from"`
	Next int64 `json:"next"`
	Size int64 `json:"size"`
	More bool  `json:"more,omitempty"`
	// Skipped counts the entries left out to keep the newest StreamEntries.
	Skipped int           `json:"skipped,omitempty"`
	Entries []StreamEntry `json:"entries"`
}

// Stream reads the stream of the newest agent flai serve started for story,
// from byte offset after, or its last StreamTailBytes when after is negative
// or past the end of the log; hostapi.ErrNoAgent when it started none. It
// returns whole lines only: a line still being written is read next time.
func Stream(st AgentState, story string, after int64) (*StreamRead, error) {
	run := st.Stories[story]
	if run == nil {
		return nil, fmt.Errorf("%s: %w", story, hostapi.ErrNoAgent)
	}
	out, err := streamRun(run, story, after)
	if err != nil {
		return nil, err
	}
	out.Story = story
	return out, nil
}

// PlanStream reads the stream of the newest planner flai serve started for
// item, an epic or a story, as Stream reads a story's agent's; an error that
// is hostapi.ErrNoAgent when it started none.
func PlanStream(st AgentState, item string, after int64) (*StreamRead, error) {
	run := st.Plans[item]
	if run == nil {
		return nil, noPlanner(item)
	}
	out, err := streamRun(run, item, after)
	if err != nil {
		return nil, err
	}
	out.Item = item
	return out, nil
}

// RoleStream reads the stream of the newest run of a strategic role flai
// serve started for the project, as Stream reads a story's agent's: the
// orchestrator's for conventions.RoleOrchestrate (S-0218), the analyzer's for
// conventions.RoleAnalyze (S-0228). It is an error that is hostapi.ErrNoAgent
// when it started none, and another for any other role.
func RoleStream(st AgentState, role string, after int64) (*StreamRead, error) {
	var run *AgentRun
	var what string
	switch role {
	case conventions.RoleOrchestrate:
		run, what = st.Orchestrator, "orchestrator"
	case conventions.RoleAnalyze:
		run, what = st.Analyzer, "analyzer"
	default:
		return nil, fmt.Errorf("there is no stream for the role %q: name %s or %s", role, conventions.RoleOrchestrate, conventions.RoleAnalyze)
	}
	if run == nil {
		return nil, noRun(what)
	}
	out, err := streamRun(run, "the "+what, after)
	if err != nil {
		return nil, err
	}
	out.Role = role
	return out, nil
}

// noPlanner is an item flai serve has started no planner for: it says so,
// and is hostapi.ErrNoAgent to whoever asks, as a story with no agent is.
type noPlanner string

// Error names the item.
func (n noPlanner) Error() string { return "flai serve has started no planner for " + string(n) }

// Is makes errors.Is find hostapi.ErrNoAgent in it.
func (noPlanner) Is(target error) bool { return target == hostapi.ErrNoAgent }

// noRun is a strategic role, the orchestrator or the analyzer, flai serve
// has started no run of for the project: it says so, and is
// hostapi.ErrNoAgent to whoever asks, as noPlanner is.
type noRun string

// Error names the role's agent.
func (n noRun) Error() string { return "flai serve has started no " + string(n) + " for this project" }

// Is makes errors.Is find hostapi.ErrNoAgent in it.
func (noRun) Is(target error) bool { return target == hostapi.ErrNoAgent }

// streamRun reads run's log as Stream says; id names what the run is for in
// its errors.
func streamRun(run *AgentRun, id string, after int64) (*StreamRead, error) {
	out := &StreamRead{Agent: run.Agent, Started: run.Started, Ended: run.Ended, Running: run.live(), Outcome: run.Outcome, Entries: []StreamEntry{}}
	if run.Log == "" {
		if run.Error != "" {
			out.Entries = append(out.Entries, StreamEntry{Kind: StreamEnd, Text: "could not be started: " + run.Error, Error: true})
		}
		return out, nil
	}
	f, err := os.Open(run.Log)
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading the stream of %s: %w", id, err)
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("reading the stream of %s: %w", id, err)
	}
	out.Size = info.Size()
	tail := after < 0 || after > out.Size
	if tail {
		after = max(0, out.Size-StreamTailBytes)
	}
	out.From, out.Next = after, after
	if _, err := f.Seek(after, io.SeekStart); err != nil {
		return nil, fmt.Errorf("reading the stream of %s: %w", id, err)
	}
	r := bufio.NewReaderSize(f, 64<<10)
	if tail && after > 0 && !lineStart(f, after) {
		// the tail starts inside a line: begin at the next one
		_, n, _, whole, err := readLine(r, 0)
		if err != nil || !whole {
			return out, nil
		}
		out.Next += n
	}
	for out.Next-out.From < streamReadBytes {
		line, n, long, whole, err := readLine(r, streamLineBytes)
		if !whole {
			if err != nil && !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("reading the stream of %s: %w", id, err)
			}
			break
		}
		out.Next += n
		if long {
			out.Entries = append(out.Entries, StreamEntry{Kind: StreamOutput, Text: fmt.Sprintf("(a line of %d bytes, too long to show)", n)})
			continue
		}
		out.Entries = append(out.Entries, streamEntries(line)...)
	}
	out.More = out.Next < out.Size
	if len(out.Entries) > StreamEntries {
		out.Skipped = len(out.Entries) - StreamEntries
		out.Entries = out.Entries[out.Skipped:]
	}
	return out, nil
}

// lineStart says whether offset begins a line: the byte before it ends one.
func lineStart(f *os.File, offset int64) bool {
	b := make([]byte, 1)
	_, err := f.ReadAt(b, offset-1)
	return err == nil && b[0] == '\n'
}

// readLine reads up to and including the next newline, keeping at most keep
// bytes of it without the newline; long says it kept less than the line. n
// counts every byte consumed, and whole says the newline was found: a line
// cut off by the end of the log is not.
func readLine(r *bufio.Reader, keep int) (line []byte, n int64, long, whole bool, err error) {
	for {
		chunk, err := r.ReadSlice('\n')
		n += int64(len(chunk))
		body := bytes.TrimSuffix(chunk, []byte("\n"))
		room := keep - len(line)
		if len(body) > room {
			long = true
		}
		line = append(line, body[:max(0, min(room, len(body)))]...)
		switch {
		case err == nil:
			return bytes.TrimSuffix(line, []byte("\r")), n, long, true, nil
		case errors.Is(err, bufio.ErrBufferFull):
			continue
		default:
			return line, n, long, false, err
		}
	}
}

// event is the part of a stream-json line an entry is made from.
type event struct {
	Type        string  `json:"type"`
	Subtype     string  `json:"subtype"`
	Model       string  `json:"model"`
	Description string  `json:"description"`
	Summary     string  `json:"summary"`
	Status      string  `json:"status"`
	Result      string  `json:"result"`
	IsError     bool    `json:"is_error"`
	NumTurns    int     `json:"num_turns"`
	DurationMS  int64   `json:"duration_ms"`
	CostUSD     float64 `json:"total_cost_usd"`
	Message     *struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

type block struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Thinking string          `json:"thinking"`
	Name     string          `json:"name"`
	Input    json.RawMessage `json:"input"`
	Content  json.RawMessage `json:"content"`
	IsError  bool            `json:"is_error"`
}

// streamEntries turns one line of the log into entries: none for an event
// that says nothing a reader would follow (heartbeats, token counts, rate
// limits), the line itself for one that is not stream-json.
func streamEntries(line []byte) []StreamEntry {
	text := strings.TrimSpace(string(line))
	if text == "" {
		return nil
	}
	var e event
	if text[0] != '{' || json.Unmarshal(line, &e) != nil || e.Type == "" {
		return []StreamEntry{{Kind: StreamOutput, Text: cut(text, streamTextRunes)}}
	}
	switch e.Type {
	case "system":
		switch e.Subtype {
		case "init":
			say := "session started"
			if e.Model != "" {
				say += " (" + e.Model + ")"
			}
			return []StreamEntry{{Kind: StreamSession, Text: say}}
		case "task_started":
			return []StreamEntry{{Kind: StreamTask, Text: cut("started: "+e.Description, streamToolRunes)}}
		case "task_notification":
			return []StreamEntry{{Kind: StreamTask, Text: cut(e.Status+": "+e.Summary, streamToolRunes)}}
		}
	case "assistant", "user":
		if e.Message == nil {
			return nil
		}
		return blockEntries(e.Message.Content)
	case "result":
		say := fmt.Sprintf("session ended (%s", e.Subtype)
		if e.NumTurns > 0 {
			say += fmt.Sprintf(", %d turns", e.NumTurns)
		}
		if e.DurationMS > 0 {
			say += fmt.Sprintf(", %ds", e.DurationMS/1000)
		}
		if e.CostUSD > 0 {
			say += fmt.Sprintf(", $%.2f", e.CostUSD)
		}
		say += ")"
		if e.Result != "" {
			say += ": " + e.Result
		}
		return []StreamEntry{{Kind: StreamEnd, Text: cut(say, streamTextRunes), Error: e.IsError || e.Subtype != "success"}}
	}
	return nil
}

// blockEntries are the entries of a message's content: a string, or blocks.
func blockEntries(raw json.RawMessage) []StreamEntry {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		if s = strings.TrimSpace(s); s != "" {
			return []StreamEntry{{Kind: StreamText, Text: cut(s, streamTextRunes)}}
		}
		return nil
	}
	var blocks []block
	if json.Unmarshal(raw, &blocks) != nil {
		return nil
	}
	var out []StreamEntry
	for _, b := range blocks {
		switch b.Type {
		case "text":
			if t := strings.TrimSpace(b.Text); t != "" {
				out = append(out, StreamEntry{Kind: StreamText, Text: cut(t, streamTextRunes)})
			}
		case "thinking":
			if t := strings.TrimSpace(b.Thinking); t != "" {
				out = append(out, StreamEntry{Kind: StreamThinking, Text: cut(t, streamTextRunes)})
			}
		case "tool_use":
			out = append(out, StreamEntry{Kind: StreamTool, Tool: b.Name, Text: cut(toolInput(b.Input), streamToolRunes)})
		case "tool_result":
			out = append(out, StreamEntry{Kind: StreamResult, Text: cut(resultText(b.Content), streamResultRunes), Error: b.IsError})
		}
	}
	return out
}

// toolInput says what a tool was asked in a line: the description or the
// command of a shell call, the file of a file tool, the pattern of a search,
// and otherwise the input's fields, in order of name.
func toolInput(raw json.RawMessage) string {
	var in map[string]any
	if json.Unmarshal(raw, &in) != nil {
		return strings.TrimSpace(string(raw))
	}
	for _, k := range []string{"description", "command", "file_path", "path", "pattern", "query", "url", "prompt"} {
		if s, ok := in[k].(string); ok && strings.TrimSpace(s) != "" {
			return oneLine(s)
		}
	}
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v, _ := json.Marshal(in[k])
		parts = append(parts, k+"="+string(v))
	}
	return strings.Join(parts, " ")
}

// resultText is a tool result's text: a string, or its text blocks joined.
func resultText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return strings.TrimSpace(s)
	}
	var blocks []block
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && strings.TrimSpace(b.Text) != "" {
			parts = append(parts, strings.TrimSpace(b.Text))
		}
	}
	return strings.Join(parts, "\n")
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// cut keeps at most n characters of s, and says when it cut.
func cut(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n]) + "…"
}
