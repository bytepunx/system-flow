// Package inbox builds an agent's inbox: the threads awaiting it, the stories
// ready to pull in pull order, what is unpublished, and what others changed
// since it last looked, read and advanced through the agent's cursor. The MCP
// tool inbox and flai task done answer the same inbox from it (ADR-0107). It
// lists the open conversations of the agent's story too, and reports a
// message to that story once, as an event (ADR-0120, S-0331).
package inbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

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
	All   bool   // threads awaiting someone else too
	Now   func() time.Time
	// Own is the agent's own story, whose open conversations the inbox lists
	// under messages (S-0331); none lists none.
	Own string
	// Runner asks git what is accepted and not yet published (ADR-0067),
	// and the newest flai release in the project's history (S-0181);
	// without one the inbox says nothing of either.
	Runner  execx.Runner
	Version string // the running flai's, compared with the newest flai release
}

// Inbox is the agent's inbox, as the MCP tool inbox answers it.
type Inbox struct {
	Agent       string   `json:"agent"`
	Unpublished []string `json:"unpublished,omitempty" jsonschema:"the accepted items no release has covered yet, by ID: what publishing would send to the remote. Information for the operator, or for an agent the operator asks to publish; not a step to take on your own"`
	// FlaiOutdated is set while this flai is older than the newest flai
	// release in the project's history (S-0181).
	FlaiOutdated *release.Outdated    `json:"flai_outdated,omitempty" jsonschema:"this MCP server's flai is older than the newest flai release tagged in the project's history, so it may lack rules and fields the project uses; listed on every call while it is true: tell the designer, who upgrades the host with the command named"`
	AwaitingYou  int                  `json:"awaiting_you" jsonschema:"threads awaiting you; messages never count here"`
	Threads      []ThreadSummary      `json:"threads"`
	Messages     []MessageSummary     `json:"messages,omitempty" jsonschema:"the open conversations of this agent's story with the agents of other open stories (ADR-0120), in ID order; absent when it has none, or this session has no story of its own"`
	Ready        []workitem.BoardCard `json:"ready" jsonschema:"stories ready to pull, in pull order; listed on every call"`
	CanPull      bool                 `json:"can_pull" jsonschema:"whether the in-progress limit leaves room and review is under its limit, so that one may be pulled"`
	PullHold     string               `json:"pull_hold,omitempty" jsonschema:"why no story may be pulled now: the in-progress limit is full, or review is full and waits on acceptance"`
	Changes      []Event              `json:"changes" jsonschema:"what others changed since this agent last looked, newest kept, at most 50; reported once. A first look under a new name covers the last 24 hours of stories and epics only"`
	Omitted      int                  `json:"changes_omitted" jsonschema:"how many older changes were left out of changes because of the cap; they are not reported later"`
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
	story := workitem.CanonicalID(opt.Story)
	for _, th := range all {
		if !th.Open() {
			continue
		}
		sum := Summarize(opt.Repo, opt.Agent, th)
		if opt.Story != "" && sum.Story != story {
			continue
		}
		if sum.Awaiting == "you" {
			out.AwaitingYou++
		} else if !opt.All {
			continue
		}
		out.Threads = append(out.Threads, sum)
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
	// the inbox lists the story's conversations as they stand, as it lists
	// threads, so its changes are the work items' only: no story for CatchUp
	done = perf.Track(ctx, "changes.read")
	out.Changes, out.Omitted, err = CatchUp(opt.Repo, opt.Agent, "", opt.Now(), items)
	done()
	if err != nil {
		return Inbox{}, err
	}
	return out, nil
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
	Awaiting  string         `json:"awaiting" jsonschema:"'you' when the last entry is not yours, else 'other'; 'other' too on a thread you opened while a recommendation on it awaits the operator's confirmation, which is no answer until they confirm it"`
	// PendingRecommendation is the recommendation awaiting the operator's
	// confirmation (ADR-0090), nil when there is none.
	PendingRecommendation *threads.Entry `json:"pending_recommendation" jsonschema:"the recommendation awaiting the operator's confirmation, null when there is none: the thread still awaits the operator until they confirm it or answer otherwise"`
	Project               string         `json:"project,omitempty" jsonschema:"the project the thread is in, when the server serves more than one"`
}

// Summarize is th as agent sees it. It awaits the agent when the last entry
// is someone else's, but for the agent that opened it while a recommendation
// on it awaits the operator's confirmation: that is no answer to its question
// until the operator confirms it (ADR-0090). Anyone else, the operator and
// the recommendation's author included, sees it by its last entry.
func Summarize(repo *workitem.Repo, agent string, th *threads.Thread) ThreadSummary {
	entries := th.Entries()
	out := ThreadSummary{ID: th.ID, Title: th.Title, Status: th.Status, Anchor: th.Anchor, Story: threads.StoryOf(repo, th), Updated: th.Updated, Entries: len(entries), Awaiting: "other", PendingRecommendation: th.PendingRecommendation()}
	if n := len(entries); n > 0 {
		out.LastBy, out.LastEntry = entries[n-1].Author, entries[n-1].Text
		asked := out.PendingRecommendation != nil && th.Opener() == agent
		if out.LastBy != agent && !asked {
			out.Awaiting = "you"
		}
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
	Last     messages.Entry `json:"last" jsonschema:"the newest entry: its time, author, story, and text"`
	Awaiting string         `json:"awaiting" jsonschema:"'you' when the conversation awaits this story, whose agent answers it with message_reply before going on, else 'other'"`
	Project  string         `json:"project,omitempty" jsonschema:"the project the conversation is in, when the server serves more than one"`
}

// Conversations are story's open conversations as its agent sees them, in ID
// order: one that reads as closed is left out (ADR-0120).
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
		sum := MessageSummary{ID: c.ID, Title: c.Title, With: otherStory(c, own), About: c.About, Updated: c.Updated, Entries: len(entries), Awaiting: "other"}
		if sum.About == nil {
			sum.About = []string{}
		}
		if n := len(entries); n > 0 {
			sum.Last = entries[n-1]
		}
		if workitem.CanonicalID(c.Awaiting()) == own {
			sum.Awaiting = "you"
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
// separated, and At the time of the newest such entry.
const Message = "message"

// MaxPaths is how many paths an overlap's summary names; To has them all.
const MaxPaths = 10

// Event is one thing that changed since the agent last looked.
type Event struct {
	workitem.Change
	Summary string `json:"summary" jsonschema:"the change in words"`
	Project string `json:"project,omitempty" jsonschema:"the project it happened in, when the server serves more than one"`
}

// CatchUp returns what others than agent changed since its cursor, among
// items listed archive included, and the messages to story, the agent's own,
// when it is given, the newest MaxEvents of them, and how many older ones it
// left out, and advances the cursor past all of them: what a look leaves out
// never comes back in a later one. A first look, with no cursor, is told of
// stories and epics only; a day of task transitions is history to an agent
// that has just arrived, not news.
func CatchUp(repo *workitem.Repo, agent, story string, now time.Time, items []*workitem.Item) ([]Event, int, error) {
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, 0, err
	}
	now = now.UTC().Truncate(time.Second)
	path := cursorPath(repo, agent)
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
	noticed := func(c workitem.Change, summary string) {
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
		events = append(events, Event{Change: c, Summary: summary})
	}
	for _, n := range itemedit.Notices(repo) {
		if n.By != agent {
			c := workitem.Change{ID: n.ID, Type: n.Type, Title: n.Title, Kind: Edited, To: strings.Join(n.Changed, ","), By: n.By, At: n.At}
			noticed(c, Describe(c))
		}
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
		noticed(c, summary)
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
					noticed(m, Describe(m))
				}
			}
		}
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
		about := ""
		if c.To != "" {
			about = " about " + namePaths(c.To)
		}
		from := c.Cause + "'s agent"
		if c.By != "" {
			from += ", " + c.By + ","
		}
		return from + " wrote in " + c.ID + " " + c.Title + about + ": read it with message_get and answer it with message_reply before you go on"
	}
	return c.ID + " " + c.Kind
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
