package usage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Measuring an agent's usage from its log (S-0143). Claude Code writes
// stream-json, one event per line, and two kinds of event say what was spent:
//
//   - result, at the end of each run, reports per model the tokens and the
//     cost of the whole session so far, subagents included. A session resumed
//     in a later run (after an answer) reports its totals again, cumulative,
//     so a session's newest result is its total.
//   - assistant, once per content block of each call to the model, with the
//     call's timestamp, message id, model, and usage. Its input and cache
//     tokens are exact once the repeats are dropped; its output tokens are the
//     first chunk's only, so they are a lower bound.
//
// A session's reported totals are taken as they are. What came after a
// session's newest result, a run still going or one that died, is counted
// from its calls and its cost estimated.
//
// A sub-agent's calls are in its story's agent's session and carry
// parent_tool_use_id, the ID of the Agent (or Task) tool_use that started it
// (S-0230, I-0054). A task is given each session's reported totals in the
// share of their calls' input and cache tokens that are its own: the calls
// of the sub-agents started for it, wherever they fall in time, and an even
// share, among the tasks in progress at the time, of every other call made
// while it was in progress. It is marked estimated.
//
// A story's empty wakes (S-0272, ADR-0105) are counted as the logs are read:
// each call to the flai MCP tool wait_for_events by the session's own agent,
// a tool_use without parent_tool_use_id, is paired by its ID with its
// tool_result, and counted when that result is not an error and reports
// timed_out with no events, no changed paths, and no end. A call whose result
// is not in the logs is not counted.

// Rates are dollars per token for each model, from what logs reported.
type Rates map[string]float64

// call is one call to a model.
type call struct {
	at      time.Time
	model   string
	session string
	usage   Model
	// agent is the ID of the tool_use that started the sub-agent that made
	// it; empty for the session's own agent.
	agent string
	// covered says a result reported after it counts it.
	covered bool
}

// weight is how much of a session's reported totals a call is given: its
// input and cache tokens, which the log reports exactly.
func (c call) weight() int64 { return c.usage.Input + c.usage.CacheRead + c.usage.CacheWrite }

// Run is one log: one run of an agent, from its first event to its last.
type Run struct {
	Path  string    `json:"path"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	// Result is the final text of the run's newest result event; empty when
	// it has none or that event carries no text.
	Result string `json:"result,omitempty"`
}

// Record is what a story's agent logs say, read in the order of their runs.
type Record struct {
	Runs  []Run
	calls []call
	// reported is each session's newest result, by model.
	reported map[string]map[string]Model
	// started are the sub-agents' starts, by the ID of their tool_use.
	started map[string]start
	// waits are the session's own agent's wait_for_events calls whose result
	// has not been read yet, by the ID of their tool_use.
	waits map[string]bool
	// emptyWakes counts those calls whose result was an empty wake.
	emptyWakes int
}

// waitTool is the name Claude Code gives the flai MCP tool wait_for_events.
const waitTool = "mcp__flai__wait_for_events"

// start is what started a sub-agent: its Agent call's description and
// prompt, or, when the call was not seen, the description its calls carry.
type start struct {
	description, prompt string
}

// Read reads the logs, oldest run first. A log that is missing is skipped:
// flai serve may not have written it yet.
func Read(paths ...string) (*Record, error) {
	rec := &Record{reported: map[string]map[string]Model{}, started: map[string]start{}, waits: map[string]bool{}}
	byID := map[string]int{} // message id to its index in calls
	for _, p := range paths {
		if err := rec.read(p, byID); err != nil {
			return nil, err
		}
	}
	return rec, nil
}

// line is the part of a stream-json event usage is measured from.
type line struct {
	Type      string `json:"type"`
	Subtype   string `json:"subtype"`
	Timestamp string `json:"timestamp"`
	// Result is a result event's final text.
	Result    string `json:"result"`
	SessionID string `json:"session_id"`
	Model     string `json:"model"`
	// ParentToolUseID is the tool_use that started the sub-agent the event
	// is from; empty for the session's own agent.
	ParentToolUseID string `json:"parent_tool_use_id"`
	// TaskDescription is the description of the Agent call that started
	// the sub-agent the event is from.
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
			} `json:"input"`
			// ToolUseID and IsError are a tool_result's: the tool_use it
			// answers, and whether the call failed.
			ToolUseID string `json:"tool_use_id"`
			IsError   bool   `json:"is_error"`
		} `json:"content"`
		Usage *struct {
			Input      int64 `json:"input_tokens"`
			Output     int64 `json:"output_tokens"`
			CacheRead  int64 `json:"cache_read_input_tokens"`
			CacheWrite int64 `json:"cache_creation_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
	ModelUsage map[string]struct {
		Input      int64   `json:"inputTokens"`
		Output     int64   `json:"outputTokens"`
		CacheRead  int64   `json:"cacheReadInputTokens"`
		CacheWrite int64   `json:"cacheCreationInputTokens"`
		Cost       float64 `json:"costUSD"`
	} `json:"modelUsage"`
	CostUSD float64 `json:"total_cost_usd"`
	Usage   *struct {
		Input      int64 `json:"input_tokens"`
		Output     int64 `json:"output_tokens"`
		CacheRead  int64 `json:"cache_read_input_tokens"`
		CacheWrite int64 `json:"cache_creation_input_tokens"`
	} `json:"usage"`
}

// Only these events are decoded; every other line is passed over unread.
var wanted = [][]byte{[]byte(`"type":"assistant"`), []byte(`"type":"user"`), []byte(`"type":"result"`), []byte(`"subtype":"init"`)}

func (rec *Record) read(path string, byID map[string]int) error {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading the agent log %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	run := Run{Path: path}
	model := "" // the session's model, for a result that names none
	r := bufio.NewReaderSize(f, 256<<10)
	for {
		raw, err := r.ReadBytes('\n')
		if len(raw) > 0 && raw[0] == '{' && wantedLine(raw) {
			var e line
			if json.Unmarshal(raw, &e) == nil {
				if at, perr := time.Parse(time.RFC3339Nano, e.Timestamp); perr == nil {
					if run.Start.IsZero() || at.Before(run.Start) {
						run.Start = at
					}
					if at.After(run.End) {
						run.End = at
					}
				}
				switch {
				case e.Type == "system" && e.Subtype == "init":
					model = e.Model
				case e.Type == "assistant":
					rec.addStarts(e)
					rec.addWaits(e)
					rec.addCall(e, byID)
				case e.Type == "user":
					rec.addWakes(e, raw)
				case e.Type == "result":
					run.Result = e.Result
					rec.addResult(e, model)
				}
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading the agent log %s: %w", path, err)
		}
	}
	if !run.Start.IsZero() {
		rec.Runs = append(rec.Runs, run)
	}
	return nil
}

func wantedLine(raw []byte) bool {
	for _, w := range wanted {
		if bytes.Contains(raw, w) {
			return true
		}
	}
	return false
}

// addStarts records the sub-agents an event starts, and the description of
// the one it is from when its start was not seen.
func (rec *Record) addStarts(e line) {
	if e.Message != nil {
		for _, b := range e.Message.Content {
			if b.Type == "tool_use" && (b.Name == "Agent" || b.Name == "Task") && b.ID != "" {
				rec.started[b.ID] = start{description: b.Input.Description, prompt: b.Input.Prompt}
			}
		}
	}
	if p := e.ParentToolUseID; p != "" && e.TaskDescription != "" {
		if _, ok := rec.started[p]; !ok {
			rec.started[p] = start{description: e.TaskDescription}
		}
	}
}

// addWaits records the wait_for_events calls the session's own agent makes
// in an event; a sub-agent's are left out.
func (rec *Record) addWaits(e line) {
	if e.Message == nil || e.ParentToolUseID != "" {
		return
	}
	for _, b := range e.Message.Content {
		if b.Type == "tool_use" && b.Name == waitTool && b.ID != "" {
			rec.waits[b.ID] = true
		}
	}
}

// addWakes pairs the tool_results of an event with the wait_for_events
// calls waiting for them, and counts those that were empty wakes. Only a
// result that answers such a call has its content decoded, from raw.
func (rec *Record) addWakes(e line, raw []byte) {
	if e.Message == nil || len(rec.waits) == 0 {
		return
	}
	for _, b := range e.Message.Content {
		if b.Type != "tool_result" || !rec.waits[b.ToolUseID] {
			continue
		}
		delete(rec.waits, b.ToolUseID)
		if !b.IsError && emptyWake(raw, b.ToolUseID) {
			rec.emptyWakes++
		}
	}
}

// emptyWake says the tool_result for id in the event raw reports that the
// wait timed out with no events, no changed paths, and no end. Its content
// is the tool's JSON as a string, or in a list of text blocks.
func emptyWake(raw []byte, id string) bool {
	var e struct {
		Message struct {
			Content []struct {
				Type      string          `json:"type"`
				ToolUseID string          `json:"tool_use_id"`
				Content   json.RawMessage `json:"content"`
			} `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(raw, &e) != nil {
		return false
	}
	for _, b := range e.Message.Content {
		if b.Type != "tool_result" || b.ToolUseID != id {
			continue
		}
		var texts []string
		var text string
		var blocks []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}
		switch {
		case json.Unmarshal(b.Content, &text) == nil:
			texts = []string{text}
		case json.Unmarshal(b.Content, &blocks) == nil:
			for _, t := range blocks {
				if t.Type == "text" {
					texts = append(texts, t.Text)
				}
			}
		}
		for _, t := range texts {
			var r struct {
				TimedOut bool              `json:"timed_out"`
				Events   []json.RawMessage `json:"events"`
				Changed  []json.RawMessage `json:"changed"`
				End      bool              `json:"end"`
			}
			if json.Unmarshal([]byte(t), &r) == nil {
				return r.TimedOut && len(r.Events) == 0 && len(r.Changed) == 0 && !r.End
			}
		}
		return false
	}
	return false
}

func (rec *Record) addCall(e line, byID map[string]int) {
	if e.Message == nil || e.Message.Usage == nil || e.Message.ID == "" {
		return
	}
	u := e.Message.Usage
	got := Model{Model: e.Message.Model, Input: u.Input, Output: u.Output, CacheRead: u.CacheRead, CacheWrite: u.CacheWrite}
	if i, ok := byID[e.Message.ID]; ok {
		// a repeat for another content block: keep the most it said
		c := &rec.calls[i]
		c.usage.Input, c.usage.Output = max(c.usage.Input, got.Input), max(c.usage.Output, got.Output)
		c.usage.CacheRead, c.usage.CacheWrite = max(c.usage.CacheRead, got.CacheRead), max(c.usage.CacheWrite, got.CacheWrite)
		return
	}
	at, _ := time.Parse(time.RFC3339Nano, e.Timestamp)
	byID[e.Message.ID] = len(rec.calls)
	rec.calls = append(rec.calls, call{at: at, model: got.Model, session: e.SessionID, usage: got, agent: e.ParentToolUseID})
}

func (rec *Record) addResult(e line, model string) {
	totals := map[string]Model{}
	for name, m := range e.ModelUsage {
		totals[name] = Model{Model: name, Input: m.Input, Output: m.Output, CacheRead: m.CacheRead, CacheWrite: m.CacheWrite, Cost: m.Cost}
	}
	if len(totals) == 0 && e.Usage != nil && model != "" {
		// a harness that reports no model breakdown: the session's model
		u := e.Usage
		totals[model] = Model{Model: model, Input: u.Input, Output: u.Output, CacheRead: u.CacheRead, CacheWrite: u.CacheWrite, Cost: e.CostUSD}
	}
	if len(totals) == 0 {
		return
	}
	rec.reported[e.SessionID] = totals
	for i := range rec.calls {
		if rec.calls[i].session == e.SessionID {
			rec.calls[i].covered = true
		}
	}
}

// Rates are the dollars per token each model's reported totals come to in
// this record.
func (rec *Record) Rates() Rates {
	cost, tokens := map[string]float64{}, map[string]int64{}
	for _, totals := range rec.reported {
		for name, m := range totals {
			cost[name] += m.Cost
			tokens[name] += m.Tokens()
		}
	}
	out := Rates{}
	for name, t := range tokens {
		if t > 0 && cost[name] > 0 {
			out[name] = cost[name] / float64(t)
		}
	}
	return out
}

// Merge adds the rates of o for models r has none for.
func (r Rates) Merge(o Rates) Rates {
	out := Rates{}
	for k, v := range o {
		out[k] = v
	}
	for k, v := range r {
		out[k] = v
	}
	return out
}

// Total is everything the record says was spent, with SourceLog, and the
// empty wakes its logs hold; nil when it says nothing. What no result
// reported is priced at rates, or at the record's own when rates has none
// for its model.
func (rec *Record) Total(rates Rates) *Usage {
	rates = rec.Rates().Merge(rates)
	u := &Usage{Source: SourceLog}
	for _, totals := range rec.reported {
		for _, m := range totals {
			u.addModel(m)
		}
	}
	for _, c := range rec.calls {
		if !c.covered {
			u.addModel(estimate(c.usage, rates))
			u.Estimated = true
		}
	}
	for _, run := range rec.Runs {
		u.Seconds += int64(run.End.Sub(run.Start).Seconds())
	}
	u.EmptyWakes = rec.emptyWakes
	if u.Empty() {
		return nil
	}
	u.Tidy()
	return u
}

// estimate prices a call no result reported at its model's rate; with no
// rate its cost is left at nothing.
func estimate(m Model, rates Rates) Model {
	m.Cost = float64(m.Tokens()) * rates[m.Model]
	return m
}

// Span is an interval of time; a zero To is open.
type Span struct {
	From, To time.Time
}

// Tasks is what each task spent, by its ID, given the spans it was in
// progress: the calls of the sub-agents started for it, and an even share,
// among the tasks in progress at the time, of every other call made while it
// was in progress. Each is marked estimated, and its seconds are the time it
// was in progress while a run went on; a task that spent nothing is absent.
func (rec *Record) Tasks(spans map[string][]Span, rates Rates) map[string]*Usage {
	rates = rec.Rates().Merge(rates)
	ids := map[string]string{} // a task ID's number to the ID
	for id := range spans {
		ids[idNumber(id)] = id
	}
	// each task's weight of each session and model, and the whole's
	type key struct{ session, model string }
	part, whole := map[string]map[key]float64{}, map[key]float64{}
	out := map[string]*Usage{}
	at := func(id string) *Usage {
		if out[id] == nil {
			out[id] = &Usage{Source: SourceLog, Estimated: true}
		}
		return out[id]
	}
	owner := map[string]string{} // a sub-agent's tool_use ID to its task
	for _, c := range rec.calls {
		k := key{c.session, c.model}
		if c.covered {
			whole[k] += float64(c.weight())
		}
		// the tasks the call is given to: its sub-agent's, or those in
		// progress at the time
		var given []string
		if c.agent != "" {
			id, ok := owner[c.agent]
			if !ok {
				id = rec.started[c.agent].task(ids)
				owner[c.agent] = id
			}
			if id != "" {
				given = []string{id}
			}
		}
		if given == nil {
			for id, ss := range spans {
				if within(ss, c.at) {
					given = append(given, id)
				}
			}
		}
		share := 1 / float64(len(given))
		for _, id := range given {
			if !c.covered {
				at(id).addModel(estimate(c.usage, rates).scaled(share))
				continue
			}
			if part[id] == nil {
				part[id] = map[key]float64{}
			}
			part[id][k] += float64(c.weight()) * share
		}
	}
	for id, weights := range part {
		for session, totals := range rec.reported {
			for name, m := range totals {
				k := key{session, name}
				if weights[k] > 0 && whole[k] > 0 {
					at(id).addModel(m.scaled(weights[k] / whole[k]))
				}
			}
		}
	}
	for id, ss := range spans {
		if s := rec.seconds(ss); s > 0 {
			at(id).Seconds = s
		}
	}
	for id, u := range out {
		if u.Empty() {
			delete(out, id)
			continue
		}
		u.Tidy()
	}
	return out
}

// within says t falls in one of the spans.
func within(spans []Span, t time.Time) bool {
	for _, s := range spans {
		if !t.Before(s.From) && (s.To.IsZero() || t.Before(s.To)) {
			return true
		}
	}
	return false
}

// seconds is how long runs went on within the spans.
func (rec *Record) seconds(spans []Span) int64 {
	var n int64
	for _, s := range spans {
		for _, run := range rec.Runs {
			from, to := run.Start, run.End
			if from.Before(s.From) {
				from = s.From
			}
			if !s.To.IsZero() && to.After(s.To) {
				to = s.To
			}
			if to.After(from) {
				n += int64(to.Sub(from).Seconds())
			}
		}
	}
	return n
}

// taskID is a task's ID as an agent writes it in text, and as one alone.
var taskID, aTaskID = regexp.MustCompile(`\bT-(\d+)\b`), regexp.MustCompile(`^T-(\d+)$`)

// idNumber is a task ID's number without its zero padding, or the ID when
// it is not a task's.
func idNumber(id string) string {
	m := aTaskID.FindStringSubmatch(id)
	if m == nil {
		return id
	}
	n, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return id
	}
	return strconv.FormatUint(n, 10)
}

// task is the task a sub-agent was started for: the first of ids its
// description names, or when it names none, the first its prompt names;
// empty when neither names one.
func (s start) task(ids map[string]string) string {
	for _, text := range []string{s.description, s.prompt} {
		for _, m := range taskID.FindAllString(text, -1) {
			if id, ok := ids[idNumber(m)]; ok {
				return id
			}
		}
	}
	return ""
}

// ReadRates are the rates the results in the logs report, read without the
// rest of each log.
func ReadRates(paths ...string) (Rates, error) {
	rec := &Record{reported: map[string]map[string]Model{}}
	for _, p := range paths {
		f, err := os.Open(p)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("reading the agent log %s: %w", p, err)
		}
		r := bufio.NewReaderSize(f, 256<<10)
		for {
			raw, err := r.ReadBytes('\n')
			if len(raw) > 0 && raw[0] == '{' && bytes.Contains(raw, []byte(`"type":"result"`)) {
				var e line
				if json.Unmarshal(raw, &e) == nil && e.Type == "result" {
					// a session per log is enough to weigh the rates by
					e.SessionID = p + "\x00" + e.SessionID
					rec.addResult(e, "")
				}
			}
			if err != nil {
				break
			}
		}
		_ = f.Close()
	}
	return rec.Rates(), nil
}
