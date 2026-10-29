package usage

import (
	"bytes"
	"encoding/json"
	"slices"
	"strings"
	"time"
)

// Claude Code's stream-json, one JSON record per line, read into events
// (S-0352, ADR-0130). Five kinds of record hold them; every other line is
// passed over unread:
//
//   - system with subtype init is a SessionStart, with the session's model.
//   - assistant is one content block of a call to the model: a Call with the
//     message's ID, model, and usage, then a ToolCall for each tool_use, with
//     its name, ID, and input's description, prompt, command, and file_path.
//     A message streamed as several blocks is several lines with its ID, and
//     the usage on each is what that line says, not the message's count.
//   - user holds a ToolResult for each tool_result, with the tool_use it
//     answers and is_error; the text of the result only for the tools whose
//     result usage reads, as only those are decoded.
//   - result is an End, with the run's final text, total_cost_usd, usage,
//     and modelUsage: each model's tokens and costUSD. One with is_error
//     ends on the usage limit when its api_error_status is 429, or when a
//     rate_limit_event rejected the run before it: its Limit is the
//     api_error, the status, and the time that event said the limit resets.
//   - rate_limit_event holds no event: its rate_limit_info, with status
//     rejected, says the run is on the usage limit until resetsAt, in Unix
//     seconds, for the result after it; with any other status, it is not.
//
// Each line's timestamp, session_id, parent_tool_use_id, and task_description
// are on every event it holds.

// waitTool is the name Claude Code gives the flai MCP tool wait_for_events.
const waitTool = "mcp__flai__wait_for_events"

// resultTools are the tools whose result the measurement reads.
var resultTools = []string{waitTool}

// claudeCodeEnd marks the one record that holds an End.
var claudeCodeEnd = []byte(`"type":"result"`)

// claudeCodeRecords mark the records that hold events, so that no other line
// is decoded.
var claudeCodeRecords = [][]byte{
	[]byte(`"type":"assistant"`), []byte(`"type":"user"`), claudeCodeEnd, []byte(`"subtype":"init"`), []byte(`"type":"rate_limit_event"`),
}

// limitStatus is a 429, Too Many Requests: the status of an error on a limit.
const limitStatus = 429

// claudeCodeReader reads one measurement's Claude Code logs, in order: it
// keeps the IDs of the calls whose result it reads, as a result's line does
// not name its tool, and the limit a rate_limit_event put the run on, as
// the result after it does not say when the limit resets.
type claudeCodeReader struct {
	reads map[string]bool
	// limit is the limit a rejecting rate_limit_event put the run on, with
	// its reset time; nil when none has since the run's start or its last
	// End, or the newest said the run is not on a limit.
	limit *Limit
}

// ClaudeCode is a Reader of Claude Code's stream-json, for the logs of one
// measurement read in order.
func ClaudeCode() Reader {
	return &claudeCodeReader{reads: map[string]bool{}}
}

// ccUsage is the usage of a message or of a run.
type ccUsage struct {
	Input      int64 `json:"input_tokens"`
	Output     int64 `json:"output_tokens"`
	CacheRead  int64 `json:"cache_read_input_tokens"`
	CacheWrite int64 `json:"cache_creation_input_tokens"`
}

// tokens are u's tokens; nil when u is.
func (u *ccUsage) tokens() *Tokens {
	if u == nil {
		return nil
	}
	return &Tokens{Input: u.Input, Output: u.Output, CacheRead: u.CacheRead, CacheWrite: u.CacheWrite}
}

// ccLine is the part of a stream-json record events are read from.
type ccLine struct {
	Type            string `json:"type"`
	Subtype         string `json:"subtype"`
	Timestamp       string `json:"timestamp"`
	Result          string `json:"result"`
	SessionID       string `json:"session_id"`
	Model           string `json:"model"`
	ParentToolUseID string `json:"parent_tool_use_id"`
	TaskDescription string `json:"task_description"`
	Message         *struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Content []struct {
			Type  string `json:"type"`
			ID    string `json:"id"`
			Name  string `json:"name"`
			Input struct {
				Description string `json:"description"`
				Prompt      string `json:"prompt"`
				Command     string `json:"command"`
				FilePath    string `json:"file_path"`
			} `json:"input"`
			ToolUseID string `json:"tool_use_id"`
			IsError   bool   `json:"is_error"`
		} `json:"content"`
		Usage *ccUsage `json:"usage"`
	} `json:"message"`
	ModelUsage map[string]struct {
		Input      int64   `json:"inputTokens"`
		Output     int64   `json:"outputTokens"`
		CacheRead  int64   `json:"cacheReadInputTokens"`
		CacheWrite int64   `json:"cacheCreationInputTokens"`
		Cost       float64 `json:"costUSD"`
	} `json:"modelUsage"`
	CostUSD *float64 `json:"total_cost_usd"`
	Usage   *ccUsage `json:"usage"`
	// IsError, APIError, and APIErrorStatus are a result's error: apart
	// from a tool_result's is_error, which is in the message's content.
	IsError        bool   `json:"is_error"`
	APIError       string `json:"api_error"`
	APIErrorStatus int    `json:"api_error_status"`
	RateLimitInfo  *struct {
		Status   string `json:"status"`
		ResetsAt int64  `json:"resetsAt"`
	} `json:"rate_limit_info"`
}

// Read is the events of one line of Claude Code's stream-json.
func (r *claudeCodeReader) Read(line []byte) ([]Event, bool) {
	if len(line) == 0 || line[0] != '{' || !slices.ContainsFunc(claudeCodeRecords, func(m []byte) bool { return bytes.Contains(line, m) }) {
		return nil, false
	}
	var l ccLine
	if json.Unmarshal(line, &l) != nil {
		return nil, false
	}
	base := Event{Session: l.SessionID, ParentCall: l.ParentToolUseID, ParentDescription: l.TaskDescription}
	if at, err := time.Parse(time.RFC3339Nano, l.Timestamp); err == nil {
		base.At = at
	}
	switch {
	case l.Type == "system" && l.Subtype == "init":
		r.limit = nil
		base.Kind, base.Model = SessionStart, l.Model
		return []Event{base}, true
	case l.Type == "assistant":
		return r.calls(l, base), true
	case l.Type == "user":
		return r.results(l, base, line), true
	case l.Type == "result":
		e := end(l, base)
		e.Limit, r.limit = r.ended(l), nil
		return []Event{e}, true
	case l.Type == "rate_limit_event":
		r.limit = rateLimit(l)
		return nil, true
	}
	return nil, false
}

// rateLimit is the limit a rate_limit_event puts the run on: its reset
// time when it rejected the run; nil when it did not.
func rateLimit(l ccLine) *Limit {
	if l.RateLimitInfo == nil || l.RateLimitInfo.Status != "rejected" {
		return nil
	}
	lim := &Limit{}
	if l.RateLimitInfo.ResetsAt > 0 {
		lim.Resets = time.Unix(l.RateLimitInfo.ResetsAt, 0).UTC()
	}
	return lim
}

// ended is the limit a result line ended the run on: its error and status,
// with the reset time of the rate_limit_event that rejected the run before
// it; nil for a result that is no error, or an error on no limit.
func (r *claudeCodeReader) ended(l ccLine) *Limit {
	if !l.IsError || (l.APIErrorStatus != limitStatus && r.limit == nil) {
		return nil
	}
	lim := Limit{Error: l.APIError, Status: l.APIErrorStatus}
	if r.limit != nil {
		lim.Resets = r.limit.Resets
	}
	return &lim
}

// MayEnd says a line may be a result record, the one that holds an End.
func (r *claudeCodeReader) MayEnd(line []byte) bool {
	return len(line) > 0 && line[0] == '{' && bytes.Contains(line, claudeCodeEnd)
}

// calls are the Call of an assistant line and a ToolCall for each tool it
// uses; it notes the calls whose result is read.
func (r *claudeCodeReader) calls(l ccLine, base Event) []Event {
	c := base
	c.Kind = Call
	if l.Message == nil {
		return []Event{c}
	}
	c.MessageID, c.Model, c.Tokens = l.Message.ID, l.Message.Model, l.Message.Usage.tokens()
	out := []Event{c}
	for _, b := range l.Message.Content {
		if b.Type != "tool_use" {
			continue
		}
		t := base
		t.Kind, t.MessageID, t.Model = ToolCall, l.Message.ID, l.Message.Model
		t.Tool, t.ToolID = b.Name, b.ID
		t.Description, t.Prompt, t.Command, t.Path = b.Input.Description, b.Input.Prompt, b.Input.Command, b.Input.FilePath
		out = append(out, t)
		if b.ID != "" && slices.Contains(resultTools, b.Name) {
			r.reads[b.ID] = true
		}
	}
	return out
}

// results are the ToolResults of a user line, with the text of those whose
// result is read, decoded from raw.
func (r *claudeCodeReader) results(l ccLine, base Event, raw []byte) []Event {
	if l.Message == nil {
		return nil
	}
	var out []Event
	read := false
	for _, b := range l.Message.Content {
		if b.Type != "tool_result" {
			continue
		}
		e := base
		e.Kind, e.ToolID, e.IsError = ToolResult, b.ToolUseID, b.IsError
		out = append(out, e)
		read = read || r.reads[b.ToolUseID]
	}
	if !read {
		return out
	}
	texts := resultTexts(raw)
	for i := range out {
		if r.reads[out[i].ToolID] {
			out[i].Result = texts[out[i].ToolID]
		}
	}
	return out
}

// resultTexts are the texts of the tool_results in the user line raw, by
// the ID of the tool_use each answers; the first is kept for an ID answered
// twice.
func resultTexts(raw []byte) map[string]string {
	var l struct {
		Message struct {
			Content []struct {
				Type      string          `json:"type"`
				ToolUseID string          `json:"tool_use_id"`
				Content   json.RawMessage `json:"content"`
			} `json:"content"`
		} `json:"message"`
	}
	out := map[string]string{}
	if json.Unmarshal(raw, &l) != nil {
		return out
	}
	for _, b := range l.Message.Content {
		if _, seen := out[b.ToolUseID]; b.Type == "tool_result" && !seen {
			out[b.ToolUseID] = resultText(b.Content)
		}
	}
	return out
}

// resultText is a tool_result's content as text: the string it is, or its
// text blocks joined by newlines.
func resultText(content json.RawMessage) string {
	var text string
	if json.Unmarshal(content, &text) == nil {
		return text
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(content, &blocks) != nil {
		return ""
	}
	var texts []string
	for _, b := range blocks {
		if b.Type == "text" {
			texts = append(texts, b.Text)
		}
	}
	return strings.Join(texts, "\n")
}

// end is the End of a result line: its text, cost, usage, and totals per
// model in order of name.
func end(l ccLine, base Event) Event {
	e := base
	e.Kind, e.Result, e.CostUSD, e.Tokens = End, l.Result, l.CostUSD, l.Usage.tokens()
	for name, m := range l.ModelUsage {
		e.Models = append(e.Models, Model{Model: name, Input: m.Input, Output: m.Output, CacheRead: m.CacheRead, CacheWrite: m.CacheWrite, Cost: m.Cost})
	}
	slices.SortFunc(e.Models, func(a, b Model) int { return strings.Compare(a.Model, b.Model) })
	return e
}
