package usage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Measuring an agent's usage from its log (S-0143). The reader of the
// harness that wrote the log turns its lines into events (ADR-0130), and two
// kinds of event say what was spent:
//
//   - End, at the end of each run, reports per model the tokens and the cost
//     of the whole session so far, subagents included. A session resumed in
//     a later run (after an answer) reports its totals again, cumulative, so
//     a session's newest End is its total.
//   - Call, for each line of each call to the model, with the call's time,
//     message ID, model, and tokens. A message logged over several lines is
//     one call, with the most any of its lines says of each kind of token.
//     With Claude Code its input and cache tokens are exact; its output
//     tokens are the first chunk's only, so they are a lower bound.
//
// A session's reported totals are taken as they are. What came after a
// session's newest End, a run still going or one that died, is counted from
// its calls and its cost estimated.
//
// A sub-agent's calls are in its story's agent's session and carry the ID of
// the tool call that started it, an Agent (or Task) call, as their
// ParentCall (S-0230, I-0054). A task is given each session's reported
// totals in the share of their calls' input and cache tokens that are its
// own: the calls of the sub-agents started for it, wherever they fall in
// time, and an even share, among the tasks in progress at the time, of every
// other call made while it was in progress. It is marked estimated.
//
// A story's empty wakes (S-0272, ADR-0105) are counted as the logs are read:
// each call to the flai MCP tool wait_for_events by the session's own agent,
// a ToolCall with no ParentCall, is paired by its ID with its ToolResult, and
// counted when that result is not an error and reports timed_out with no
// events, no changed paths, and no end. A call whose result is not in the
// logs is not counted.
//
// A story's turns (S-0293) are gathered as the logs are read, each message of
// the session's own agent with the tools it called, and classified and
// counted by day when the record is totalled: turns.go says how.

// Rates are dollars per token for each model, from what logs reported.
type Rates map[string]float64

// call is one call to a model.
type call struct {
	at      time.Time
	model   string
	session string
	usage   Model
	// agent is the ID of the tool call that started the sub-agent that made
	// it; empty for the session's own agent.
	agent string
	// covered says an End reported after it counts it.
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
	// Result is the final text of the run's newest End; empty when it has
	// none or that End carries no text.
	Result string `json:"result,omitempty"`
}

// Record is what a story's agent logs say, read in the order of their runs.
type Record struct {
	Runs  []Run
	calls []call
	// reported is each session's newest End, by model.
	reported map[string]map[string]Model
	// started are the sub-agents' starts, by the ID of their tool call.
	started map[string]start
	// waits are the session's own agent's wait_for_events calls whose result
	// has not been read yet, by the ID of their tool call.
	waits map[string]bool
	// emptied are the IDs of those calls whose result was an empty wake.
	emptied map[string]bool
	// turns are the session's own agent's messages, in the order they were
	// first seen, and turnByID their index by message ID.
	turns    []turn
	turnByID map[string]int
}

// start is what started a sub-agent: its Agent call's description and
// prompt, or, when the call was not seen, the description its calls carry.
type start struct {
	description, prompt string
}

// Read reads the logs with r, the reader of the harness that wrote them,
// oldest run first. A log that is missing is skipped: flai serve may not
// have written it yet.
func Read(r Reader, paths ...string) (*Record, error) {
	rec := &Record{
		reported: map[string]map[string]Model{}, started: map[string]start{}, waits: map[string]bool{},
		emptied: map[string]bool{}, turnByID: map[string]int{},
	}
	byID := map[string]int{} // message id to its index in calls
	for _, p := range paths {
		if err := rec.read(r, p, byID); err != nil {
			return nil, err
		}
	}
	return rec, nil
}

// eachLine hands each line of the log at path to read; a log that is
// missing has none.
func eachLine(path string, read func(line []byte)) error {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading the agent log %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	r := bufio.NewReaderSize(f, 256<<10)
	for {
		raw, err := r.ReadBytes('\n')
		if len(raw) > 0 {
			read(raw)
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading the agent log %s: %w", path, err)
		}
	}
}

func (rec *Record) read(r Reader, path string, byID map[string]int) error {
	run := Run{Path: path}
	model := "" // the session's model, for an End that names none
	err := eachLine(path, func(line []byte) {
		events, _ := r.Read(line)
		for _, e := range events {
			if !e.At.IsZero() {
				if run.Start.IsZero() || e.At.Before(run.Start) {
					run.Start = e.At
				}
				if e.At.After(run.End) {
					run.End = e.At
				}
			}
			switch e.Kind {
			case SessionStart:
				model = e.Model
			case Call:
				rec.addDescribed(e)
				rec.addTurn(e)
				rec.addCall(e, byID)
			case ToolCall:
				rec.addStart(e)
				rec.addWait(e)
				rec.addTool(e)
			case ToolResult:
				rec.addWake(e)
			case End:
				run.Result = e.Result
				rec.addResult(e, model)
			}
		}
	})
	if err != nil {
		return err
	}
	if !run.Start.IsZero() {
		rec.Runs = append(rec.Runs, run)
	}
	return nil
}

// addStart records the sub-agent a tool call starts: an Agent or Task call.
func (rec *Record) addStart(e Event) {
	if (e.Tool == "Agent" || e.Tool == "Task") && e.ToolID != "" {
		rec.started[e.ToolID] = start{description: e.Description, prompt: e.Prompt}
	}
}

// addDescribed records the description of the sub-agent a call is from when
// its start was not seen.
func (rec *Record) addDescribed(e Event) {
	if p := e.ParentCall; p != "" && e.ParentDescription != "" {
		if _, ok := rec.started[p]; !ok {
			rec.started[p] = start{description: e.ParentDescription}
		}
	}
}

// addWait records a wait_for_events call the session's own agent makes; a
// sub-agent's is left out.
func (rec *Record) addWait(e Event) {
	if e.ParentCall == "" && e.Tool == waitTool && e.ToolID != "" {
		rec.waits[e.ToolID] = true
	}
}

// addWake pairs a tool's result with the wait_for_events call waiting for
// it, and counts it when it was an empty wake.
func (rec *Record) addWake(e Event) {
	if !rec.waits[e.ToolID] {
		return
	}
	delete(rec.waits, e.ToolID)
	if !e.IsError && emptyWake(e.Result) {
		rec.emptied[e.ToolID] = true
	}
}

// addTurn opens the turn of a call's message, at its first call, when the
// session's own agent made it; a sub-agent's is left out.
func (rec *Record) addTurn(e Event) {
	if e.MessageID == "" || e.ParentCall != "" {
		return
	}
	if _, ok := rec.turnByID[e.MessageID]; !ok {
		rec.turnByID[e.MessageID] = len(rec.turns)
		rec.turns = append(rec.turns, turn{at: e.At})
	}
}

// addTool gathers a tool call of the session's own agent into the turn of
// its message, which a call of the message opened.
func (rec *Record) addTool(e Event) {
	if e.MessageID == "" || e.ParentCall != "" {
		return
	}
	rec.addTurn(e)
	i := rec.turnByID[e.MessageID]
	rec.turns[i].calls = append(rec.turns[i].calls, toolCall{id: e.ToolID, name: e.Tool, command: e.Command, path: e.Path})
}

// emptyWake says the text of a wait_for_events result reports that the wait
// timed out with no events, no changed paths, and no end: the first JSON
// object in it says so.
func emptyWake(text string) bool {
	d := json.NewDecoder(strings.NewReader(text))
	for {
		var r struct {
			TimedOut bool              `json:"timed_out"`
			Events   []json.RawMessage `json:"events"`
			Changed  []json.RawMessage `json:"changed"`
			End      bool              `json:"end"`
		}
		err := d.Decode(&r)
		var notObject *json.UnmarshalTypeError
		switch {
		case err == nil:
			return r.TimedOut && len(r.Events) == 0 && len(r.Changed) == 0 && !r.End
		case !errors.As(err, &notObject):
			return false
		}
	}
}

// addCall counts a call to the model once per message: a repeat of its ID,
// another line of the message, keeps the most any line said.
func (rec *Record) addCall(e Event, byID map[string]int) {
	if e.Tokens == nil || e.MessageID == "" {
		return
	}
	u := e.Tokens
	got := Model{Model: e.Model, Input: u.Input, Output: u.Output, CacheRead: u.CacheRead, CacheWrite: u.CacheWrite}
	if i, ok := byID[e.MessageID]; ok {
		c := &rec.calls[i]
		c.usage.Input, c.usage.Output = max(c.usage.Input, got.Input), max(c.usage.Output, got.Output)
		c.usage.CacheRead, c.usage.CacheWrite = max(c.usage.CacheRead, got.CacheRead), max(c.usage.CacheWrite, got.CacheWrite)
		return
	}
	byID[e.MessageID] = len(rec.calls)
	rec.calls = append(rec.calls, call{at: e.At, model: got.Model, session: e.Session, usage: got, agent: e.ParentCall})
}

// addResult keeps what an End reports as its session's totals, by model, or
// for model, the session's, when it reports no breakdown by model.
func (rec *Record) addResult(e Event, model string) {
	totals := map[string]Model{}
	for _, m := range e.Models {
		totals[m.Model] = m
	}
	if len(totals) == 0 && e.Tokens != nil && model != "" {
		// a harness that reports no model breakdown: the session's model
		u := e.Tokens
		m := Model{Model: model, Input: u.Input, Output: u.Output, CacheRead: u.CacheRead, CacheWrite: u.CacheWrite}
		if e.CostUSD != nil {
			m.Cost = *e.CostUSD
		}
		totals[model] = m
	}
	if len(totals) == 0 {
		return
	}
	rec.reported[e.Session] = totals
	for i := range rec.calls {
		if rec.calls[i].session == e.Session {
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
// empty wakes and the turns by day its logs hold; nil when it says nothing. What no End
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
	u.EmptyWakes = len(rec.emptied)
	u.Turns = turnDays(rec.turns, rec.emptied)
	if u.Empty() {
		return nil
	}
	u.Tidy()
	return u
}

// estimate prices a call no End reported at its model's rate; with no
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
	owner := map[string]string{} // a sub-agent's tool call ID to its task
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

// endSkimmer is a Reader that tells, without decoding a line, whether it may
// hold an End, so that ReadRates decodes no other line.
type endSkimmer interface {
	// MayEnd says the line may hold an End.
	MayEnd(line []byte) bool
}

// ReadRates are the rates the Ends in the logs report, read with r, the
// reader of the harness that wrote them; when r can skim, without the rest
// of each log.
func ReadRates(r Reader, paths ...string) (Rates, error) {
	rec := &Record{reported: map[string]map[string]Model{}}
	skim, _ := r.(endSkimmer)
	for _, p := range paths {
		err := eachLine(p, func(line []byte) {
			if skim != nil && !skim.MayEnd(line) {
				return
			}
			events, _ := r.Read(line)
			for _, e := range events {
				if e.Kind == End {
					// a session per log is enough to weigh the rates by
					e.Session = p + "\x00" + e.Session
					rec.addResult(e, "")
				}
			}
		})
		if err != nil {
			return nil, err
		}
	}
	return rec.Rates(), nil
}
