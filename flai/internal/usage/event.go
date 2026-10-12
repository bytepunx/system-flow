package usage

import "time"

// The neutral events usage is measured from (ADR-0130). Each harness's log
// has a shape of its own; a Reader turns its lines into these, one line at a
// time, and reports each line as it is: a message streamed over several
// lines is several events with one message ID, and a call's tokens are what
// that line says. Gathering the repeats, pairing a tool's call with its
// result, and totalling are the measurement's.

// Reader turns the lines of a harness's log into events.
type Reader interface {
	// Read is the events one line of the log holds, in the order the line
	// holds them; a line may hold none, or several, as a model call that
	// calls several tools does. ok says the reader read the line as a
	// record of the harness's log: false for a line that is none, such as
	// an error written to the log or a record it cannot decode, and for a
	// record of a kind it passes over unread; true for one it read, even
	// when the line holds no event.
	Read(line []byte) (events []Event, ok bool)
}

// Kind is what an event records.
type Kind string

// The kinds of event.
const (
	// SessionStart is a session's start, with its model.
	SessionStart Kind = "session_start"
	// Call is a call to a model, with its message ID, model, and tokens.
	Call Kind = "call"
	// ToolCall is a model call's use of a tool, with the tool's name, ID,
	// and the parts of its input usage reads.
	ToolCall Kind = "tool_call"
	// ToolResult is the result of a tool's call, by the call's ID.
	ToolResult Kind = "tool_result"
	// End is the end of a run, with what the harness reports it spent.
	End Kind = "end"
)

// Tokens are the tokens of a call, or of a run.
type Tokens struct {
	Input      int64
	Output     int64
	CacheRead  int64
	CacheWrite int64
}

// Event is one thing a harness's log records.
type Event struct {
	// Kind is what the event records.
	Kind Kind
	// At is when the harness logged it; zero when the line says not.
	At time.Time
	// Session is the harness's ID of the session it is in.
	Session string
	// Model is a SessionStart's model, or that of the message a Call or a
	// ToolCall is part of.
	Model string
	// MessageID is the ID of the model's message a Call, or a ToolCall,
	// is part of; a message the harness logs over several lines repeats it.
	MessageID string
	// Tokens are a Call's tokens as its line reports them, or what an End
	// reports for the run when it has no breakdown by model; nil when the
	// line reports none.
	Tokens *Tokens
	// CostUSD is the run's cost in US dollars as an End reports it; nil
	// when the harness reports none.
	CostUSD *float64
	// Models are an End's totals and cost per model, in order of name.
	Models []Model
	// ParentCall is the ID of the tool call that started the sub-agent the
	// event is from; empty for the session's own agent.
	ParentCall string
	// ParentDescription is the description of that tool call, as the
	// sub-agent's lines repeat it; empty when they do not.
	ParentDescription string
	// Tool is a ToolCall's tool, as the harness names it.
	Tool string
	// ToolID is the ID of a ToolCall, or of the call a ToolResult answers.
	ToolID string
	// Description and Prompt are a ToolCall's input of those names, as a
	// call that starts a sub-agent gives them.
	Description string
	Prompt      string
	// Command is a ToolCall's shell command; Path is the file it writes.
	Command string
	Path    string
	// IsError says a ToolResult's call failed.
	IsError bool
	// Result is the text of a ToolResult, for the tools whose result usage
	// reads (wait_for_events), or the final text of the run on an End.
	Result string
}
