// Package inbox builds an agent's inbox: the threads awaiting it, the stories
// ready to pull in pull order, what is unpublished, and what others changed
// since it last looked, read and advanced through the agent's cursor. The MCP
// tool inbox and flai task done answer the same inbox from it (ADR-0107). It
// lists the open conversations of the agent's story too, and reports a
// message to that story once, as an event (ADR-0120, S-0331). The events
// that happen in flai serve and flai accept, planned, ready-empty,
// release-due, and an agent's run started, ended, failed, or stopped, are
// appended to a log of their own, which it reads back with the rest; that a
// thread came to await the agent it works out from the thread files
// (ADR-0211).
package inbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Options say whose inbox to read, in which project, and how.
type Options struct {
	Repo  *workitem.Repo
	Agent string // whose inbox; "agent" when empty
	Story string // only threads on this story and its tasks, in any zero padding
	All   bool   // every open thread, those awaiting the operator or another too
	// Operator lists the threads awaiting the operator too, as the
	// orchestrator's role does (ADR-0181).
	Operator bool
	Now      func() time.Time
	// Own is the agent's own story, whose open conversations the inbox lists
	// under messages (S-0331), and whose threads, and its tasks', await the
	// agent (ADR-0181); none lists none, and the agent takes part only in the
	// threads it wrote to.
	Own string
	// Runner asks git what is accepted and not yet published (ADR-0067),
	// and the newest flai release in the project's history (S-0181);
	// without one the inbox says nothing of either.
	Runner  execx.Runner
	Version string // the running flai's, compared with the newest flai release
	// Budget reads the host's budget state, each provider's that is set
	// (S-0534, ADR-0209 §5); nil, or none read, leaves the inbox without one.
	Budget func() []Budget
}

// Budget is one provider's budget state as an agent's inbox carries it
// (S-0534, ADR-0209 §5): flai serve set it when a run ended on the provider's
// usage limit, and starts no agent on the provider while it holds. Times are
// in TimeFormat, UTC.
type Budget struct {
	Provider string `json:"provider" jsonschema:"the provider the state is kept under, or the harness, such as claude-code, when the runs named no provider"`
	Limit    string `json:"limit,omitempty" jsonschema:"the error the harness named, such as usage_limit_reached; absent when it named none"`
	Status   int    `json:"status,omitempty" jsonschema:"the HTTP status of the refusal, such as 429; absent when the harness named none"`
	Since    string `json:"since" jsonschema:"when the state was set"`
	Until    string `json:"until" jsonschema:"when the limit resets, as the API reported it or a backoff guessed it, and flai serve may start a run on the provider again"`
	Probing  bool   `json:"probing,omitempty" jsonschema:"a run flai serve started once until passed runs now, to learn whether the limit has reset; the state clears when it ends on anything but the limit"`
}

// Inbox is the agent's inbox, as the MCP tool inbox answers it.
type Inbox struct {
	Agent       string   `json:"agent"`
	Unpublished []string `json:"unpublished,omitempty" jsonschema:"the accepted items no release has covered yet, by ID: what publishing would send to the remote. Information for the operator, or for an agent the operator asks to publish; not a step to take on your own"`
	// FlaiOutdated is set while this flai is older than the newest flai
	// release in the project's history (S-0181).
	FlaiOutdated *release.Outdated `json:"flai_outdated,omitempty" jsonschema:"this MCP server's flai is older than the newest flai release tagged in the project's history, so it may lack rules and fields the project uses; listed on every call while it is true: tell the designer, who upgrades the host with the command named"`
	// Budget is the host's budget state, one entry per provider whose state
	// is set (S-0534); absent when none is.
	Budget      []Budget        `json:"budget,omitempty" jsonschema:"the providers this host is out of budget on, each with its limit and until, the reset time: while yours is listed, do not retry or start work that calls the model in a loop; write your narrative's state and end, and flai serve starts you again once the limit resets (ADR-0209 §5). Absent when no provider's state is set"`
	AwaitingYou int             `json:"awaiting_you" jsonschema:"threads awaiting you, those with awaiting 'you' (ADR-0181); threads awaiting the operator or another never count here, nor do messages"`
	Threads     []ThreadSummary `json:"threads"`
	// ThreadsOmitted counts the open threads on other items that a
	// story-scoped read leaves out of Threads where it would have listed them
	// (ADR-0179): without All, only those awaiting the agent (ADR-0181), and
	// those awaiting the operator with Operator; zero, and absent, without a
	// story.
	ThreadsOmitted int `json:"threads_omitted,omitempty" jsonschema:"open threads on other items awaiting you that this story's inbox leaves out, counted rather than listed (ADR-0179, ADR-0181): call inbox to list them"`
	// EntriesShortened counts the listed threads whose last entry or pending
	// recommendation Fit shortened (ADR-0179).
	EntriesShortened int                  `json:"entries_shortened,omitempty" jsonschema:"how many listed threads' last entries and pending recommendations were shortened to fit the answer (ADR-0179); read them whole with thread_get"`
	Messages         []MessageSummary     `json:"messages,omitempty" jsonschema:"the open conversations of this agent's story with the agents of other open stories (ADR-0120), in ID order; absent when it has none, or this session has no story of its own"`
	Ready            []workitem.BoardCard `json:"ready" jsonschema:"stories ready to pull, in pull order; listed on every call"`
	CanPull          bool                 `json:"can_pull" jsonschema:"whether the in-progress limit leaves room and review is under its limit, so that one may be pulled"`
	PullHold         string               `json:"pull_hold,omitempty" jsonschema:"why no story may be pulled now: the in-progress limit is full, or review is full and waits on acceptance"`
	Changes          []Event              `json:"changes" jsonschema:"what others changed since this agent last looked, newest kept, at most 50; reported once. A first look under a new name covers the last 24 hours of stories and epics only"`
	Omitted          int                  `json:"changes_omitted" jsonschema:"how many older changes were left out of changes because of the cap; they are not reported later"`
}

// Read builds the agent's inbox and advances its cursor past the changes it
// reports, so that no other reader of the same agent's inbox reports them
// again.
func Read(ctx context.Context, opt Options) (Inbox, error) {
	if opt.Agent == "" {
		opt.Agent = "agent"
	}
	if opt.Now == nil {
		opt.Now = time.Now
	}
	done := perf.Track(ctx, "threads.read")
	all, err := threads.List(opt.Repo)
	done()
	if err != nil {
		return Inbox{}, err
	}
	out := Inbox{Agent: opt.Agent, Threads: []ThreadSummary{}}
	reader := Reader{Agent: opt.Agent, Own: opt.Own}
	story := workitem.CanonicalID(opt.Story)
	for _, th := range all {
		if !th.Open() {
			continue
		}
		sum := Summarize(opt.Repo, reader, th)
		listed := opt.All || sum.Awaiting == "you" || (opt.Operator && sum.Awaiting == "operator")
		if opt.Story != "" && workitem.CanonicalID(sum.Story) != story {
			// counted where it would have been listed (ADR-0179)
			if listed {
				out.ThreadsOmitted++
			}
			continue
		}
		if sum.Awaiting == "you" {
			out.AwaitingYou++
		}
		if listed {
			out.Threads = append(out.Threads, sum)
		}
	}
	// one list for the board and the changes (S-0156)
	done = perf.Track(ctx, "repo.list")
	items, err := opt.Repo.List(true)
	done()
	if err != nil {
		return Inbox{}, err
	}
	view, err := Board(ctx, opt.Repo, opt.Runner, opt.Now, items, false)
	if err != nil {
		return Inbox{}, err
	}
	out.Ready, out.CanPull, out.PullHold, out.Unpublished = view.ReadyInPullOrder(), view.CanPull(), view.PullHold(), view.Unpublished
	out.FlaiOutdated = release.FlaiOutdated(opt.Runner, opt.Repo.Root, opt.Version)
	if out.Ready == nil {
		out.Ready = []workitem.BoardCard{}
	}
	if opt.Own != "" {
		done = perf.Track(ctx, "messages.read")
		out.Messages, err = Conversations(opt.Repo, opt.Own)
		done()
		if err != nil {
			return Inbox{}, err
		}
	}
	// the inbox lists the story's conversations as they stand, so it gives
	// CatchUp no story for messages; the reader's own story still says whom
	// a thread awaits (ADR-0181)
	done = perf.Track(ctx, "changes.read")
	out.Changes, out.Omitted, err = CatchUp(opt.Repo, reader, "", opt.Operator, opt.Now(), items)
	done()
	if err != nil {
		return Inbox{}, err
	}
	if opt.Story != "" {
		out.Changes = onStory(out.Changes, story)
	}
	if opt.Budget != nil {
		out.Budget = opt.Budget()
	}
	return out, nil
}

// onStory is changes without the thread changes on threads of other stories
// than story: a story's inbox counts those threads in ThreadsOmitted rather
// than lists them (ADR-0179), and does not list them as changes either.
func onStory(changes []Event, story string) []Event {
	kept := changes[:0:0]
	for _, e := range changes {
		if e.Kind == Thread && workitem.CanonicalID(e.Cause) != story {
			continue
		}
		kept = append(kept, e)
	}
	return kept
}

// ShortEntry is the most bytes of a listed thread's last entry, and of its
// pending recommendation's text, that Fit keeps when it shortens them, before
// the ellipsis it appends (ADR-0179).
const ShortEntry = 500

// Fit cuts in until size, the encoded size of the answer that carries in, is
// at most limit (ADR-0179), and reports whether it is. It cuts in place, so
// that an answer holding in by pointer measures the cut: first the changes
// from the oldest, the fewest that make the answer fit, each counted in
// Omitted; then, when the answer is over without any changes, the last entry
// and the pending recommendation's text of each listed thread to ShortEntry
// bytes with an ellipsis, each thread so shortened counted once in
// EntriesShortened. It never writes to the slices or entries in held before:
// a caller that keeps a copy of *in keeps the inbox whole.
func Fit(in *Inbox, limit int, size func() int) bool {
	if in == nil || size() <= limit {
		return true
	}
	changes, omitted := in.Changes, in.Omitted
	drop := func(k int) {
		in.Changes, in.Omitted = changes[k:], omitted+k
	}
	// the size falls as more changes are dropped: the fewest that fit
	k := sort.Search(len(changes)+1, func(k int) bool {
		drop(k)
		return size() <= limit
	})
	if k <= len(changes) {
		drop(k)
		return true
	}
	drop(len(changes))
	listed := make([]ThreadSummary, len(in.Threads))
	copy(listed, in.Threads)
	for i := range listed {
		th := &listed[i]
		cut := false
		if s, ok := shorten(th.LastEntry); ok {
			th.LastEntry, cut = s, true
		}
		if rec := th.PendingRecommendation; rec != nil {
			if s, ok := shorten(rec.Text); ok {
				short := *rec
				short.Text = s
				th.PendingRecommendation, cut = &short, true
			}
		}
		if cut {
			in.EntriesShortened++
		}
	}
	in.Threads = listed
	return size() <= limit
}

// shorten is s cut to ShortEntry bytes at a rune boundary, with an ellipsis,
// and whether it was longer.
func shorten(s string) (string, bool) {
	if len(s) <= ShortEntry {
		return s, false
	}
	n := ShortEntry
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n] + "…", true
}

// Board is the board of items, listed archive included: a done story stays
// on it until it is published (S-0087), which runner tells when there is one.
func Board(ctx context.Context, repo *workitem.Repo, runner execx.Runner, now func() time.Time, items []*workitem.Item, all bool) (workitem.BoardView, error) {
	done := perf.Track(ctx, "board.load")
	board, err := repo.LoadBoard()
	done()
	if err != nil {
		return workitem.BoardView{}, err
	}
	var ids map[string]bool // without a runner, nothing is known to be unpublished
	if runner != nil {
		done = perf.Track(ctx, "release.pending")
		ids = release.PendingIDs(execx.Timed(ctx, runner), repo.Root, repo.Manifest, repo)
		done()
	}
	return workitem.NewBoardView(items, board, now(), all, ids, repo.Holds(items)), nil
}

// ---- threads ----

// ThreadSummary is a thread without its entries.
type ThreadSummary struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Status    string         `json:"status"`
	Anchor    threads.Anchor `json:"anchor"`
	Story     string         `json:"story,omitempty" jsonschema:"the story this thread belongs to, when anchored to a story or task"`
	Updated   string         `json:"updated"`
	Entries   int            `json:"entries"`
	LastBy    string         `json:"last_by"`
	LastEntry string         `json:"last_entry" jsonschema:"text of the most recent entry"`
	Awaiting  string         `json:"awaiting" jsonschema:"'you' when you take part in the thread (you are among its participants) or it is on your own story or one of its tasks, and its last entry is not yours, but not on a thread you opened while a recommendation on it awaits the operator's confirmation, which is no answer until they confirm it; else 'operator' when it awaits the operator, a recommendation on it pending or its last entry by neither the project's owner nor its story's; else 'other' (ADR-0181)"`
	// PendingRecommendation is the recommendation awaiting the operator's
	// confirmation (ADR-0090), nil when there is none.
	PendingRecommendation *threads.Entry `json:"pending_recommendation" jsonschema:"the recommendation awaiting the operator's confirmation, null when there is none: the thread still awaits the operator until they confirm it or answer otherwise"`
	Project               string         `json:"project,omitempty" jsonschema:"the project the thread is in, when the server serves more than one"`
}

// Reader is who reads a thread (ADR-0181).
type Reader struct {
	Agent string // the reader's name; "agent" when empty
	Own   string // the reader's own story, in any zero padding; "" for none
	// Cursor names the cursor CatchUp reads and advances, in place of
	// Agent's, so that flai serve can read as an agent without moving that
	// agent's own cursor (ADR-0209); "" for Agent's.
	Cursor string
}

// Summarize is th as r sees it (ADR-0181). It awaits r, "you", when r takes
// part in it, among its participants, or it is on r's own story or one of its
// tasks, and its last entry is someone else's; but not for the agent that
// opened it while a recommendation on it awaits the operator's confirmation:
// that is no answer to its question until the operator confirms it
// (ADR-0090). Any other open thread awaits the operator, "operator", while a
// recommendation on it is pending or its last entry is by none of the
// operator's names (threads.OperatorName), and is "other" otherwise, as is a
// thread with no entries.
func Summarize(repo *workitem.Repo, r Reader, th *threads.Thread) ThreadSummary {
	if r.Agent == "" {
		r.Agent = "agent"
	}
	entries := th.Entries()
	out := ThreadSummary{ID: th.ID, Title: th.Title, Status: th.Status, Anchor: th.Anchor, Story: threads.StoryOf(repo, th), Updated: th.Updated, Entries: len(entries), Awaiting: "other", PendingRecommendation: th.PendingRecommendation()}
	n := len(entries)
	if n == 0 {
		return out
	}
	out.LastBy, out.LastEntry = entries[n-1].Author, entries[n-1].Text
	pending := out.PendingRecommendation != nil
	ownStory := r.Own != "" && out.Story != "" && workitem.CanonicalID(out.Story) == workitem.CanonicalID(r.Own)
	takesPart := ownStory || slices.Contains(th.Participants, r.Agent)
	asked := pending && th.Opener() == r.Agent
	switch {
	case takesPart && out.LastBy != r.Agent && !asked:
		out.Awaiting = "you"
	case pending || !threads.OperatorName(repo, out.Story, out.LastBy):
		out.Awaiting = "operator"
	}
	return out
}

// ---- messages ----

// MessageSummary is a conversation of the agent's story without its entries,
// as the agent of that story sees it (S-0331).
type MessageSummary struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	With     string         `json:"with" jsonschema:"the other story of the two, whose agent this story's agent talks to"`
	About    []string       `json:"about" jsonschema:"the repository paths the conversation is about"`
	Updated  string         `json:"updated"`
	Entries  int            `json:"entries"`
	Last     messages.Entry `json:"last" jsonschema:"the newest entry: its time, author, story, text, and whether it is settled"`
	Awaiting string         `json:"awaiting" jsonschema:"'you' when the conversation awaits this story, whose agent answers it with message_reply before going on; empty while its newest entry from a story is settled, a reply that asks nothing, so that it awaits neither story and no reply is owed (ADR-0154); else 'other'"`
	Project  string         `json:"project,omitempty" jsonschema:"the project the conversation is in, when the server serves more than one"`
}

// Conversations are story's open conversations as its agent sees them, in ID
// order: one that reads as closed is left out (ADR-0120). One that is settled
// awaits neither story, and its Awaiting is empty (ADR-0154).
func Conversations(repo *workitem.Repo, story string) ([]MessageSummary, error) {
	all, err := messages.For(repo, story)
	if err != nil {
		return nil, err
	}
	own := workitem.CanonicalID(story)
	var out []MessageSummary
	for _, c := range all {
		if closed, _ := c.Closed(repo); closed {
			continue
		}
		entries := c.Entries()
		sum := MessageSummary{ID: c.ID, Title: c.Title, With: otherStory(c, own), About: c.About, Updated: c.Updated, Entries: len(entries)}
		if sum.About == nil {
			sum.About = []string{}
		}
		if n := len(entries); n > 0 {
			sum.Last = entries[n-1]
		}
		switch awaiting := c.Awaiting(); {
		case awaiting == "":
			// settled: it awaits neither story, and Awaiting stays empty
		case workitem.CanonicalID(awaiting) == own:
			sum.Awaiting = "you"
		default:
			sum.Awaiting = "other"
		}
		out = append(out, sum)
	}
	return out, nil
}

// otherStory is the story of c's two that is not own, a canonical ID.
func otherStory(c *messages.Conversation, own string) string {
	if workitem.CanonicalID(c.From) == own {
		return c.To
	}
	return c.From
}

// ---- changes ----

// A cursor is how far one agent has read the repository's changes. It is a
// read marker under .flai-cache, outside git: tooling never owns state
// (overview.md), and losing a cursor only repeats or skips the report of a
// change. Ready work and open threads are state and are always listed.
type cursor struct {
	Seen  string          `json:"seen"`            // TimeFormat; changes up to here were reported
	Keys  map[string]bool `json:"keys,omitempty"`  // changes stamped with that very second, already reported
	Order []string        `json:"order,omitempty"` // the pull order as last reported
	known bool            // a cursor file existed
}

// firstLook is how far back a new agent is told about.
const firstLook = 24 * time.Hour

// MaxEvents is the most changes one look reports; the newest are kept and the
// rest are counted, not listed. It is a constant, not a setting: the bound
// exists so that a look always fits a tool result. The first look of a new
// agent on 2026-09-19 was 209 changes and 68 KB, about 330 characters each,
// and did not fit (S-0061); fifty is about 16 KB, which leaves room for
// threads and ready work, and is more than an agent acts on in one look.
const MaxEvents = 50

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func cursorPath(repo *workitem.Repo, agent string) string {
	name := strings.Trim(unsafeName.ReplaceAllString(agent, "-"), "-.")
	if name == "" {
		name = "agent"
	}
	return filepath.Join(repo.CacheDir(), "mcp", name+".json")
}

func loadCursor(path string) cursor {
	var c cursor
	if data, err := os.ReadFile(path); err == nil && json.Unmarshal(data, &c) == nil && c.Seen != "" {
		c.known = true
	}
	return c
}

func (c cursor) since(now time.Time) time.Time {
	if t, err := time.Parse(workitem.TimeFormat, c.Seen); err == nil {
		return t
	}
	return now.Add(-firstLook)
}

// saveCursor is best effort: a cursor that cannot be written costs a
// repeated report, not a failed call.
func saveCursor(path string, c cursor) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	tmp := path + ".tmp"
	if os.WriteFile(tmp, data, 0o600) == nil {
		_ = os.Rename(tmp, path)
	}
}

// Edited is the kind of a change someone made to an item's own fields or
// body with flai edit; To names what of it changed.
const Edited = "edited"

// Overlapped is the kind of the change to an open story when a story whose
// changes its claim covers is accepted: Cause is the accepted story, To the
// paths, comma separated (S-0132). It is also the kind of the change to each
// of two stories in progress when a write grows the claim of one into the
// other's: Cause is the other story, To the paths the claim gained (I-0059).
const Overlapped = "overlapped"

// Message is the kind, and the type, of the change that a message came to the
// agent's story from the other story of a conversation, a new conversation or
// a reply (S-0331): ID and Title are the conversation's, Cause the sender's
// story, By its agent, To the paths the conversation is about, comma
// separated, and At the time of the newest such entry. The summary of a
// settled entry says that no reply is owed (ADR-0154).
const Message = "message"

// MaxPaths is how many paths an overlap's summary names; To has them all.
const MaxPaths = 10

// Event is one thing that changed since the agent last looked.
type Event struct {
	workitem.Change
	Summary string `json:"summary" jsonschema:"the change in words"`
	// Fields are a logged event's own, and a thread change's (ADR-0211),
	// under the Field keys.
	Fields  map[string]string `json:"fields,omitempty" jsonschema:"the event's own figures, by name, for the kinds flai serve and flai accept log, and for thread: ready-empty's room and limit (the ready count is in to), release-due's before and after, what flai release --evaluate answered before the acceptance and answers after it; agent's role, the run's kind (story, plan, orchestrate, analyze, curate, or a project role's name), and why, why it came to what to says (started, ended, failed, stopped); thread's awaits, whom the thread awaits by name (your agent name, or operator), and on, the item it is anchored to or its path; absent for other kinds"`
	Project string            `json:"project,omitempty" jsonschema:"the project it happened in, when the server serves more than one"`
}

// CatchUp returns what others than r's agent changed since its cursor, among
// items listed archive included, the events flai serve and flai accept logged
// (ADR-0211), whoever caused them, save the agent events of r's agent's own
// runs (S-0532), the messages to story, the agent's own,
// when it is given, and the threads whose newest entry makes them await r, or,
// with operator, as the orchestrator's role reads them, the operator
// (ADR-0181, ADR-0211 §5); the newest MaxEvents of them, and how many older
// ones it left out. It advances the cursor past all of them: what a look
// leaves out never comes back in a later one. A first look, with no cursor, is
// told of stories and epics only; a day of task transitions is history to an
// agent that has just arrived, not news. r's Own says whom a thread awaits
// whether or not story is given, as the inbox lists the story's messages
// rather than reports them.
func CatchUp(repo *workitem.Repo, r Reader, story string, operator bool, now time.Time, items []*workitem.Item) ([]Event, int, error) {
	agent := r.Agent
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, 0, err
	}
	now = now.UTC().Truncate(time.Second)
	name := agent
	if r.Cursor != "" {
		name = r.Cursor
	}
	path := cursorPath(repo, name)
	cur := loadCursor(path)
	events := []Event{}
	next := cursor{Seen: now.Format(workitem.TimeFormat), Keys: map[string]bool{}, Order: board.Order}
	for _, c := range workitem.Changes(items, cur.since(now), agent, cur.Keys) {
		// the cursor passes every change, reported or not
		if c.At == next.Seen {
			next.Keys[c.Key()] = true
		}
		if !cur.known && c.Type == workitem.Task {
			continue
		}
		events = append(events, Event{Change: c, Summary: Describe(c)})
	}
	// Edits someone else made with flai edit or from the dashboard (S-0085).
	// They are not in the items' front matter, which is strict and shared
	// with older flai; flai edit notes them beside the cursors.
	since := cur.since(now).UTC().Truncate(time.Second)
	noticed := func(e Event) {
		c := e.Change
		at, err := time.Parse(workitem.TimeFormat, c.At)
		news := at.After(since) || (at.Equal(since) && !cur.Keys[c.Key()])
		if err != nil || !news {
			return
		}
		if c.At == next.Seen {
			next.Keys[c.Key()] = true
		}
		if !cur.known && c.Type == workitem.Task {
			return
		}
		events = append(events, e)
	}
	for _, n := range itemedit.Notices(repo) {
		if n.By != agent {
			c := workitem.Change{ID: n.ID, Type: n.Type, Title: n.Title, Kind: Edited, To: strings.Join(n.Changed, ","), By: n.By, At: n.At}
			noticed(Event{Change: c, Summary: Describe(c)})
		}
	}
	// What happened in flai serve and flai accept (ADR-0211): news to every
	// reader, whoever caused it, the orchestrator its acceptance's release-due
	// too; but an agent is not told of its own runs (S-0532), as it is not of
	// its own moves.
	for _, l := range LoggedEvents(repo) {
		if l.Kind == Agent && l.By == agent {
			continue
		}
		c := l.Change()
		noticed(Event{Change: c, Summary: describeLogged(c, l.Fields), Fields: l.Fields})
	}
	// Open stories that an acceptance changed paths under (S-0132), and
	// stories in progress whose claims a write grew to overlap (I-0059): the
	// change is on the story told, caused by the accepted or the other story,
	// and names the paths. Whoever accepted or wrote, it is news to the agent
	// of the story told.
	for _, o := range itemedit.Overlaps(repo) {
		c := workitem.Change{ID: o.ID, Type: workitem.Story, Title: o.Title, Kind: Overlapped, To: strings.Join(o.Paths, ","), By: o.By, Cause: o.Cause(), At: o.At}
		summary := Describe(c)
		if o.Accepted == "" {
			summary = describeGrown(c, o.Grew, o.Reached)
		}
		noticed(Event{Change: c, Summary: summary})
	}
	// Messages to the agent's story (S-0331): one change for each open
	// conversation whose newest entry by the other story is news, so that a
	// new conversation and a reply are each told once; the agent's own are not.
	if story != "" {
		convs, err := messages.For(repo, story)
		if err != nil {
			return nil, 0, err
		}
		for _, c := range convs {
			if e, ok := lastFromOther(c, workitem.CanonicalID(story), agent); ok {
				if closed, _ := c.Closed(repo); !closed {
					m := workitem.Change{ID: c.ID, Type: Message, Title: c.Title, Kind: Message, To: strings.Join(c.About, ","), By: e.Author, Cause: e.Story, At: e.At}
					noticed(Event{Change: m, Summary: describeMessage(m, e.Settled)})
				}
			}
		}
	}
	// Threads whose newest entry makes them await the reader, or the
	// operator for the orchestrator (ADR-0211 §5): no log, the thread files.
	awaiting, err := threadEvents(repo, r, operator, now)
	if err != nil {
		return nil, 0, err
	}
	for _, e := range awaiting {
		noticed(e)
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].At < events[j].At })
	omitted := 0
	if len(events) > MaxEvents {
		omitted = len(events) - MaxEvents
		events = events[omitted:] // oldest first, so the newest are at the end
	}
	// keys already reported for a second that is still the current one
	if cur.Seen == next.Seen {
		for k := range cur.Keys {
			next.Keys[k] = true
		}
	}
	if cur.known && workitem.Reordered(cur.Order, board.Order) {
		events = append(events, Event{
			Change:  workitem.Change{ID: "board", Type: "board", Title: "pull order", Kind: "reordered", At: next.Seen},
			Summary: "the pull order is now " + strings.Join(board.Order, ", "),
		})
	}
	saveCursor(path, next)
	return events, omitted, nil
}

// Describe says a change in words.
func Describe(c workitem.Change) string {
	who := ""
	if c.By != "" {
		who = " by " + c.By
	}
	switch c.Kind {
	case workitem.Moved:
		if c.Cause != "" {
			return c.ID + " " + c.Title + " was cancelled with " + c.Cause + who
		}
		if c.Follows != "" {
			return c.ID + " " + c.Title + " moved to " + c.To + who + ", following " + c.Follows
		}
		return c.ID + " " + c.Title + " moved to " + c.To + who
	case workitem.WasBlocked:
		return c.ID + " " + c.Title + " was blocked: " + c.Reason
	case workitem.Unblocked:
		return c.ID + " " + c.Title + " was unblocked"
	case Edited:
		return c.ID + " " + c.Title + " was edited" + who + ": " + strings.ReplaceAll(c.To, ",", ", ") + ". Read it again before you go on"
	case Overlapped:
		return c.Cause + " was accepted" + who + " and changed " + namePaths(c.To) + ", which " + c.ID + " " + c.Title + " claims. Run flai stream sync " + c.ID + " and the tests before you go on"
	case Message:
		return describeMessage(c, false)
	case Planned, ReadyEmpty, ReleaseDue, Agent:
		return describeLogged(c, nil)
	case Thread:
		return describeThread(c)
	}
	return c.ID + " " + c.Kind
}

// describeMessage says a change of kind message: what to do with an entry
// that awaits a reply, or, for a settled one, that none is owed (ADR-0154).
func describeMessage(c workitem.Change, settled bool) string {
	about := ""
	if c.To != "" {
		about = " about " + namePaths(c.To)
	}
	from := c.Cause + "'s agent"
	if c.By != "" {
		from += ", " + c.By + ","
	}
	if settled {
		return from + " wrote in " + c.ID + " " + c.Title + about + ", settled: read it with message_get; no reply is owed"
	}
	return from + " wrote in " + c.ID + " " + c.Title + about + ": read it with message_get and answer it with message_reply before you go on"
}

// describeGrown says an overlapped change whose cause is not an acceptance:
// a write grew the claim of the story grew into that of reached, and both
// stories now claim the paths (I-0059).
func describeGrown(c workitem.Change, grew, reached string) string {
	who := ""
	if c.By != "" {
		who = ", written by " + c.By
	}
	return grew + "'s claim grew to overlap " + reached + "'s on " + namePaths(c.To) + who + ". Both stories claim them now: agree with " + c.Cause + "'s agent in the two stories' conversation (message_reply) who changes them first, before " + c.ID + " " + c.Title + " changes them"
}

// lastFromOther is c's newest entry written for the story of its two that is
// not own, by an agent other than agent, if there is one: an entry that closed
// it names no story, and one by own's agent is not news to it.
func lastFromOther(c *messages.Conversation, own, agent string) (messages.Entry, bool) {
	es := c.Entries()
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if e.Story != "" && workitem.CanonicalID(e.Story) != own && e.Author != agent {
			return e, true
		}
	}
	return messages.Entry{}, false
}

// namePaths names the comma-separated paths of an overlap, at most MaxPaths.
func namePaths(to string) string {
	paths := strings.Split(to, ",")
	if len(paths) > MaxPaths {
		return strings.Join(paths[:MaxPaths], ", ") + fmt.Sprintf(" and %d more", len(paths)-MaxPaths)
	}
	return strings.Join(paths, ", ")
}
