package usage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
// from its calls and its cost estimated. A window of time is given each
// session's reported totals in the share of their calls' input and cache
// tokens that fall in it, and is marked estimated.

// Rates are dollars per token for each model, from what logs reported.
type Rates map[string]float64

// call is one call to a model.
type call struct {
	at      time.Time
	model   string
	session string
	usage   Model
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
}

// Record is what a story's agent logs say, read in the order of their runs.
type Record struct {
	Runs  []Run
	calls []call
	// reported is each session's newest result, by model.
	reported map[string]map[string]Model
}

// Read reads the logs, oldest run first. A log that is missing is skipped:
// flai serve may not have written it yet.
func Read(paths ...string) (*Record, error) {
	rec := &Record{reported: map[string]map[string]Model{}}
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
	SessionID string `json:"session_id"`
	Model     string `json:"model"`
	Message   *struct {
		ID    string `json:"id"`
		Model string `json:"model"`
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
					rec.addCall(e, byID)
				case e.Type == "result":
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
	rec.calls = append(rec.calls, call{at: at, model: got.Model, session: e.SessionID, usage: got})
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

// Total is everything the record says was spent, with SourceLog; nil when it
// says nothing. What no result reported is priced at rates, or at the
// record's own when rates has none for its model.
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
	if u.Empty() {
		return nil
	}
	u.Tidy()
	return u
}

// Window is what was spent from from until to, with SourceLog and marked
// estimated; nil when nothing was. A zero to is open.
func (rec *Record) Window(from, to time.Time, rates Rates) *Usage {
	rates = rec.Rates().Merge(rates)
	in := func(t time.Time) bool { return !t.Before(from) && (to.IsZero() || t.Before(to)) }
	// each session and model's weight in the window and in all
	type key struct{ session, model string }
	part, whole := map[key]int64{}, map[key]int64{}
	u := &Usage{Source: SourceLog, Estimated: true}
	for _, c := range rec.calls {
		if !c.covered {
			if in(c.at) {
				u.addModel(estimate(c.usage, rates))
			}
			continue
		}
		k := key{c.session, c.model}
		whole[k] += c.weight()
		if in(c.at) {
			part[k] += c.weight()
		}
	}
	for session, totals := range rec.reported {
		for name, m := range totals {
			k := key{session, name}
			if part[k] > 0 && whole[k] > 0 {
				u.addModel(m.scaled(float64(part[k]) / float64(whole[k])))
			}
		}
	}
	for _, run := range rec.Runs {
		start, end := run.Start, run.End
		if start.Before(from) {
			start = from
		}
		if !to.IsZero() && end.After(to) {
			end = to
		}
		if end.After(start) {
			u.Seconds += int64(end.Sub(start).Seconds())
		}
	}
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

// Windows is what was spent over the spans, as Window is for one; nil when
// nothing was.
func (rec *Record) Windows(spans []Span, rates Rates) *Usage {
	var out *Usage
	for _, s := range spans {
		w := rec.Window(s.From, s.To, rates)
		if w == nil {
			continue
		}
		if out == nil {
			out = &Usage{Source: SourceLog}
		}
		out.Add(w)
	}
	if out != nil {
		out.Tidy()
	}
	return out
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
