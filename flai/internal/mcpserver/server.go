// Package mcpserver serves a system-flow repository to agents over the
// Model Context Protocol (ADR-0020). Files stay the record: every tool reads
// and writes the same files, through the same code, as the flai CLI. The
// server is for agents, so it never accepts a story; that is the operator's.
package mcpserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Options configure the server.
type Options struct {
	Repo    *workitem.Repo
	Agent   string           // the caller's identity on thread entries and transitions
	Version string           // flai version, reported to clients
	Now     func() time.Time // default time.Now
	Poll    time.Duration    // wait_for_events polling interval, default 250ms
	MaxWait time.Duration    // upper bound for one wait_for_events call, default 5m
	// Runner runs git to ask whether an acceptance is unpushed (S-0063).
	// Without one the server says nothing about it.
	Runner execx.Runner
	// Closing, when closed, ends every held wait_for_events as if its time had
	// passed, so a server over HTTP can stop without cutting agents off
	// (S-0076). The SDK does not tie a session's call to its HTTP request.
	Closing <-chan struct{}
}

type server struct {
	repo    *workitem.Repo
	agent   string
	now     func() time.Time
	poll    time.Duration
	maxWait time.Duration
	runner  execx.Runner
	closing <-chan struct{}
}

// New builds the MCP server with its tools and resources.
func New(opt Options) *mcp.Server {
	s := &server{repo: opt.Repo, agent: opt.Agent, now: opt.Now, poll: opt.Poll, maxWait: opt.MaxWait, runner: opt.Runner, closing: opt.Closing}
	if s.agent == "" {
		s.agent = "agent"
	}
	if s.now == nil {
		s.now = time.Now
	}
	if s.poll <= 0 {
		s.poll = 250 * time.Millisecond
	}
	if s.maxWait <= 0 {
		s.maxWait = 5 * time.Minute
	}
	srv := mcp.NewServer(&mcp.Implementation{Name: "flai", Title: "system-flow repository", Version: opt.Version}, &mcp.ServerOptions{
		Instructions: "This server is the agent's view of a system-flow repository. Call inbox at the start of every turn or session, at every task transition, and before moving a story to review: it lists threads awaiting you, the stories ready to pull in pull order, and what others changed since you last looked (at most 50 changes, the newest; changes_omitted counts older ones that are not reported again; your first look covers the last 24 hours of stories and epics only, so use board and item_get for how things stand). When nothing is in progress and can_pull is true, pull the first ready story without waiting to be told. An agent that stays running holds wait_for_events when idle; one that ends its turn calls inbox when it starts again, and nothing in between is lost. Reply to threads with thread_reply and ask the designer questions with thread_open. Stories are accepted by the operator only: item_move refuses to move a story or epic to done. A change that says an item was cancelled with a parent means the parent was cancelled and took it along: if it is your story or one of its tasks, stop work on it, log that in the narrative, and leave its branch and worktree alone. When inbox reports unpushed, an acceptance was made where nothing could push it: on the host run git fetch, then flai push --pending, before anything else; it never forces, and if it refuses because the remote moved, merge and run it again.",
	})
	mcp.AddTool(srv, &mcp.Tool{Name: "inbox", Description: "What needs this agent: unresolved threads (awaiting is 'you' when the last entry is not yours), the stories ready to pull in pull order with can_pull from the in-progress limit, and the changes others made to work items since this agent last looked, reported once: at most 50, newest kept, with changes_omitted counting the older ones left out. A first look covers 24 hours of stories and epics only. Filter by story to see only threads on a story and its tasks."}, s.inbox)
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_get", Description: "One thread with all of its dated entries."}, s.threadGet)
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_open", Description: "Open a thread on a repository path (optionally a heading in it) or a work item ID, to ask the designer a question or record a discussion."}, s.threadOpen)
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_reply", Description: "Add an entry to a thread as this agent. A reply from anyone but the opener marks the thread answered."}, s.threadReply)
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_resolve", Description: "Close a thread, optionally saying what settled it."}, s.threadResolve)
	mcp.AddTool(srv, &mcp.Tool{Name: "item_get", Description: "A work item by ID (any zero padding): front matter, body, and children."}, s.itemGet)
	mcp.AddTool(srv, &mcp.Tool{Name: "item_move", Description: "Transition a work item with the workflow rules enforced. Moving an item to cancelled also cancels everything open under it, and the result lists what went with it. Refuses to move a story or epic to done: acceptance is the operator's."}, s.itemMove)
	mcp.AddTool(srv, &mcp.Tool{Name: "doc_get", Description: "A markdown document under the design, docs, or wip folders, by repository path."}, s.docGet)
	mcp.AddTool(srv, &mcp.Tool{Name: "who_touches", Description: "In-progress and in-review items whose touches cover a path; ask before editing a path someone else is working on."}, s.whoTouches)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_events", Description: "Return what others changed since this agent last looked, at once when there is something already, otherwise block until a thread, work item, or narrative changes or the timeout passes. Hold this when idle to react to the designer within a second. At most 50 events, newest kept; events_omitted counts the rest."}, s.waitForEvents)
	mcp.AddTool(srv, &mcp.Tool{Name: "board", Description: "The kanban board as flai board --json prints it: cards per column, WIP limits, the pull order, and limit breaches. Stories only unless all is set."}, s.board)
	for _, key := range []string{"design", "docs"} {
		srv.AddResourceTemplate(&mcp.ResourceTemplate{
			Name:        key,
			Description: "Markdown documents under the repository's " + key + " folder",
			MIMEType:    "text/markdown",
			URITemplate: "flai://" + key + "/{+path}",
		}, s.readResource)
	}
	return srv
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
	Awaiting  string         `json:"awaiting" jsonschema:"'you' when the last entry is not yours, else 'other'"`
}

// ThreadDetail is a thread with its entries.
type ThreadDetail struct {
	ThreadSummary
	Participants []string        `json:"participants"`
	EntryList    []threads.Entry `json:"entry_list"`
}

func (s *server) summary(th *threads.Thread) ThreadSummary {
	entries := th.Entries()
	out := ThreadSummary{ID: th.ID, Title: th.Title, Status: th.Status, Anchor: th.Anchor, Story: threads.StoryOf(s.repo, th), Updated: th.Updated, Entries: len(entries), Awaiting: "other"}
	if n := len(entries); n > 0 {
		out.LastBy, out.LastEntry = entries[n-1].Author, entries[n-1].Text
		if out.LastBy != s.agent {
			out.Awaiting = "you"
		}
	}
	return out
}

func (s *server) detail(th *threads.Thread) ThreadDetail {
	return ThreadDetail{ThreadSummary: s.summary(th), Participants: th.Participants, EntryList: th.Entries()}
}

// InboxIn filters the inbox.
type InboxIn struct {
	Story string `json:"story,omitempty" jsonschema:"only threads on this story and its tasks"`
	All   bool   `json:"all,omitempty" jsonschema:"include threads awaiting someone else"`
}

// InboxOut is the agent's inbox.
type InboxOut struct {
	Agent       string               `json:"agent"`
	Unpushed    *pending.Unpushed    `json:"unpushed,omitempty" jsonschema:"an acceptance made in this clone and not pushed, listed on every call while it is true: push it from the host with git fetch and then flai push --pending, which never forces"`
	AwaitingYou int                  `json:"awaiting_you"`
	Threads     []ThreadSummary      `json:"threads"`
	Ready       []workitem.BoardCard `json:"ready" jsonschema:"stories ready to pull, in pull order; listed on every call"`
	CanPull     bool                 `json:"can_pull" jsonschema:"whether the in-progress limit leaves room to pull one"`
	Changes     []Event              `json:"changes" jsonschema:"what others changed since this agent last looked, newest kept, at most 50; reported once. A first look under a new name covers the last 24 hours of stories and epics only"`
	Omitted     int                  `json:"changes_omitted" jsonschema:"how many older changes were left out of changes because of the cap; they are not reported later"`
}

func (s *server) inbox(_ context.Context, _ *mcp.CallToolRequest, in InboxIn) (*mcp.CallToolResult, InboxOut, error) {
	all, err := threads.List(s.repo)
	if err != nil {
		return nil, InboxOut{}, err
	}
	out := InboxOut{Agent: s.agent, Threads: []ThreadSummary{}}
	story := workitem.CanonicalID(in.Story)
	for _, th := range all {
		if !th.Open() {
			continue
		}
		sum := s.summary(th)
		if in.Story != "" && sum.Story != story {
			continue
		}
		if sum.Awaiting == "you" {
			out.AwaitingYou++
		} else if !in.All {
			continue
		}
		out.Threads = append(out.Threads, sum)
	}
	view, err := s.boardView(false)
	if err != nil {
		return nil, InboxOut{}, err
	}
	out.Ready, out.CanPull, out.Unpushed = view.ReadyInPullOrder(), view.CanPull(), view.Unpushed
	if out.Ready == nil {
		out.Ready = []workitem.BoardCard{}
	}
	if out.Changes, out.Omitted, err = s.catchUp(); err != nil {
		return nil, InboxOut{}, err
	}
	return nil, out, nil
}

func (s *server) boardView(all bool) (workitem.BoardView, error) {
	items, err := s.repo.List(false)
	if err != nil {
		return workitem.BoardView{}, err
	}
	board, err := s.repo.LoadBoard()
	if err != nil {
		return workitem.BoardView{}, err
	}
	view := workitem.NewBoardView(items, board, s.now(), all)
	root := s.repo.MainRoot
	if root == "" {
		root = s.repo.Root
	}
	if u := pending.Detect(s.runner, root); u.Pending() {
		view.Unpushed = u
	}
	return view, nil
}

// BoardIn selects what the board shows.
type BoardIn struct {
	All bool `json:"all,omitempty" jsonschema:"include epics and tasks"`
}

func (s *server) board(_ context.Context, _ *mcp.CallToolRequest, in BoardIn) (*mcp.CallToolResult, workitem.BoardView, error) {
	view, err := s.boardView(in.All)
	if view.Breaches == nil {
		view.Breaches = []string{}
	}
	if view.Order == nil {
		view.Order = []string{}
	}
	return nil, view, err
}

// ThreadIDIn names a thread.
type ThreadIDIn struct {
	ID string `json:"id" jsonschema:"thread ID such as TH-0001"`
}

func (s *server) threadGet(_ context.Context, _ *mcp.CallToolRequest, in ThreadIDIn) (*mcp.CallToolResult, ThreadDetail, error) {
	th, err := threads.Get(s.repo, in.ID)
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	return nil, s.detail(th), nil
}

// ThreadOpenIn opens a thread.
type ThreadOpenIn struct {
	On      string `json:"on" jsonschema:"repository path or work item ID the thread is about"`
	Heading string `json:"heading,omitempty" jsonschema:"a heading in the document; it must exist"`
	Title   string `json:"title"`
	Text    string `json:"text" jsonschema:"the first entry"`
}

func (s *server) threadOpen(_ context.Context, _ *mcp.CallToolRequest, in ThreadOpenIn) (*mcp.CallToolResult, ThreadDetail, error) {
	th, err := threads.New(s.repo, threads.NewOptions{Title: in.Title, On: in.On, Heading: in.Heading, Author: s.agent, Text: in.Text, Now: s.now()})
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	return nil, s.detail(th), s.mirror(th)
}

// ThreadReplyIn replies to a thread.
type ThreadReplyIn struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (s *server) threadReply(_ context.Context, _ *mcp.CallToolRequest, in ThreadReplyIn) (*mcp.CallToolResult, ThreadDetail, error) {
	th, err := threads.Reply(s.repo, in.ID, s.agent, in.Text, s.now())
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	return nil, s.detail(th), s.mirror(th)
}

// ThreadResolveIn resolves a thread.
type ThreadResolveIn struct {
	ID     string `json:"id"`
	Reason string `json:"reason,omitempty" jsonschema:"what settled it"`
}

func (s *server) threadResolve(_ context.Context, _ *mcp.CallToolRequest, in ThreadResolveIn) (*mcp.CallToolResult, ThreadDetail, error) {
	th, err := threads.Resolve(s.repo, in.ID, s.agent, in.Reason, s.now())
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	return nil, s.detail(th), s.mirror(th)
}

func (s *server) mirror(th *threads.Thread) error {
	if story := threads.StoryOf(s.repo, th); story != "" {
		return threads.MirrorNarrative(s.repo, story)
	}
	return nil
}

// ---- items ----

// ItemIDIn names a work item.
type ItemIDIn struct {
	ID string `json:"id" jsonschema:"work item ID in any zero padding, such as S-39 or S-0039"`
}

// ItemBrief is a child or owner listing.
type ItemBrief struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Title   string   `json:"title"`
	Status  string   `json:"status"`
	Touches []string `json:"touches,omitempty"`
}

// ItemOut is a work item with its body and children.
type ItemOut struct {
	ItemBrief
	Nature      string                `json:"nature"`
	Parent      string                `json:"parent,omitempty"`
	Owner       string                `json:"owner,omitempty"`
	Tags        []string              `json:"tags"`
	Created     string                `json:"created"`
	Updated     string                `json:"updated"`
	Archived    bool                  `json:"archived"`
	Blocked     bool                  `json:"blocked"`
	Transitions []workitem.Transition `json:"transitions"`
	Path        string                `json:"path"`
	Body        string                `json:"body"`
	Children    []ItemBrief           `json:"children"`
}

func brief(it *workitem.Item) ItemBrief {
	return ItemBrief{ID: it.ID, Type: it.Type, Title: it.Title, Status: it.Status, Touches: it.Touches}
}

func (s *server) itemOut(it *workitem.Item) (ItemOut, error) {
	items, err := s.repo.List(true)
	if err != nil {
		return ItemOut{}, err
	}
	out := ItemOut{ItemBrief: brief(it), Nature: it.Nature, Parent: it.Parent, Owner: it.Owner, Tags: it.Tags, Created: it.Created, Updated: it.Updated, Archived: it.Archived, Transitions: it.Transitions, Path: s.rel(it.Path), Body: it.Body, Children: []ItemBrief{}}
	for _, b := range it.Blocked {
		if b.Until == "" {
			out.Blocked = true
		}
	}
	if out.Tags == nil {
		out.Tags = []string{}
	}
	if out.Transitions == nil {
		out.Transitions = []workitem.Transition{}
	}
	for _, c := range workitem.Children(items, it.ID) {
		out.Children = append(out.Children, brief(c))
	}
	return out, nil
}

func (s *server) itemGet(_ context.Context, _ *mcp.CallToolRequest, in ItemIDIn) (*mcp.CallToolResult, ItemOut, error) {
	it, err := s.repo.Get(in.ID)
	if err != nil {
		return nil, ItemOut{}, err
	}
	out, err := s.itemOut(it)
	return nil, out, err
}

// ItemMoveIn transitions an item.
type ItemMoveIn struct {
	ID     string `json:"id"`
	To     string `json:"to" jsonschema:"one of backlog, ready, in-progress, review, done, cancelled"`
	Reason string `json:"reason,omitempty" jsonschema:"required for cancelled and for review back to in-progress"`
}

// ItemMoveOut is the result of a transition.
type ItemMoveOut struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Warnings []string `json:"warnings"`
	// Cancelled is what a move to cancelled took with it: everything that was open under the item.
	Cancelled []workitem.Cascaded `json:"cancelled,omitempty" jsonschema:"items cancelled with this one, each with the state it was in"`
}

func (s *server) itemMove(_ context.Context, _ *mcp.CallToolRequest, in ItemMoveIn) (*mcp.CallToolResult, ItemMoveOut, error) {
	it, err := s.repo.Get(in.ID)
	if err != nil {
		return nil, ItemMoveOut{}, err
	}
	if in.To == workitem.Done && it.Type != workitem.Task {
		return nil, ItemMoveOut{}, fmt.Errorf("%s is %s: moving it to done is acceptance, which only the operator does (flai accept, or the dashboard); move it to review and say what is ready", it.ID, articled(it.Type))
	}
	res, err := s.repo.TransitionAll(it, in.To, s.agent, in.Reason, s.now())
	if err != nil {
		return nil, ItemMoveOut{}, err
	}
	warnings := res.Warnings
	if warnings == nil {
		warnings = []string{}
	}
	return nil, ItemMoveOut{ID: it.ID, Status: it.Status, Warnings: warnings, Cancelled: res.Cancelled}, nil
}

func articled(typ string) string {
	if typ == workitem.Epic {
		return "an epic"
	}
	return "a " + typ
}

// WhoTouchesIn asks who is working on a path.
type WhoTouchesIn struct {
	Path string `json:"path" jsonschema:"repository path or component name"`
}

// WhoTouchesOut lists the owners of a path.
type WhoTouchesOut struct {
	Path  string      `json:"path"`
	Items []ItemBrief `json:"items"`
}

func (s *server) whoTouches(_ context.Context, _ *mcp.CallToolRequest, in WhoTouchesIn) (*mcp.CallToolResult, WhoTouchesOut, error) {
	items, err := s.repo.List(false)
	if err != nil {
		return nil, WhoTouchesOut{}, err
	}
	want := strings.TrimSuffix(filepath.ToSlash(in.Path), "/")
	out := WhoTouchesOut{Path: want, Items: []ItemBrief{}}
	for _, it := range items {
		if it.Status != workitem.InProgress && it.Status != workitem.Review {
			continue
		}
		for _, t := range it.Touches {
			t = strings.TrimSuffix(t, "/")
			if t == want || strings.HasPrefix(want, t+"/") || strings.HasPrefix(t, want+"/") {
				out.Items = append(out.Items, brief(it))
				break
			}
		}
	}
	return nil, out, nil
}

// ---- documents ----

// DocIn names a document.
type DocIn struct {
	Path string `json:"path" jsonschema:"repository path of a markdown file under design, docs, or wip"`
}

// DocOut is a document split at its front matter.
type DocOut struct {
	Path        string `json:"path"`
	FrontMatter string `json:"front_matter" jsonschema:"the YAML front matter, without the --- fences"`
	Body        string `json:"body"`
}

// docPath resolves a repository path inside the layout folders, preferring
// the tree flai runs in (a story worktree) and falling back to the main
// checkout, where wip always lives.
func (s *server) docPath(rel string) (string, error) {
	abs, _, err := docedit.Resolve(s.repo, rel)
	return abs, err
}

func (s *server) docGet(_ context.Context, _ *mcp.CallToolRequest, in DocIn) (*mcp.CallToolResult, DocOut, error) {
	p, err := s.docPath(in.Path)
	if err != nil {
		return nil, DocOut{}, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, DocOut{}, err
	}
	out := DocOut{Path: filepath.ToSlash(filepath.Clean(in.Path)), Body: string(data)}
	if fm, body, err := workitem.SplitFrontMatter(string(data)); err == nil {
		out.FrontMatter, out.Body = fm, body
	}
	return nil, out, nil
}

func (s *server) readResource(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI
	rest, ok := strings.CutPrefix(uri, "flai://")
	if !ok {
		return nil, mcp.ResourceNotFoundError(uri)
	}
	key, path, ok := strings.Cut(rest, "/")
	dir := s.repo.Manifest.Layout[key]
	if !ok || dir == "" || (key != "design" && key != "docs") {
		return nil, mcp.ResourceNotFoundError(uri)
	}
	p, err := s.docPath(strings.TrimSuffix(dir, "/") + "/" + path)
	if err != nil {
		return nil, mcp.ResourceNotFoundError(uri)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, mcp.ResourceNotFoundError(uri)
	}
	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "text/markdown", Text: string(data)}}}, nil
}

// ---- events ----

// WaitIn bounds a wait.
type WaitIn struct {
	TimeoutSeconds int `json:"timeout_seconds,omitempty" jsonschema:"how long to wait, default 60, capped by the server"`
}

// WaitOut reports what changed.
type WaitOut struct {
	Events   []Event  `json:"events" jsonschema:"what others changed to work items since this agent last looked, newest kept, at most 50"`
	Omitted  int      `json:"events_omitted" jsonschema:"how many older events were left out because of the cap; they are not reported later"`
	Changed  []string `json:"changed" jsonschema:"repository paths that were added, modified, or removed while waiting"`
	TimedOut bool     `json:"timed_out"`
}

// watched returns the folders whose changes matter to an agent: threads,
// work items, and narratives, all in the main checkout.
func (s *server) watched() []string {
	return []string{threads.Dir(s.repo), s.repo.KanbanDir(), s.repo.AgentsDir()}
}

// snapshot fingerprints every markdown file under the watched folders.
// Polling is deliberate: inotify is unreliable across bind mounts and WSL,
// and the tree is small.
func (s *server) snapshot() map[string]string {
	out := map[string]string{}
	for _, dir := range s.watched() {
		_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				// A watched folder that does not exist yet (no threads so
				// far) or cannot be read simply contributes no files.
				return filepath.SkipDir
			}
			if d.IsDir() || !strings.HasSuffix(p, ".md") {
				return nil
			}
			if info, err := d.Info(); err == nil {
				out[p] = fmt.Sprintf("%d/%d", info.ModTime().UnixNano(), info.Size())
			}
			return nil
		})
	}
	return out
}

func diff(before, after map[string]string) []string {
	var changed []string
	for p, v := range after {
		if before[p] != v {
			changed = append(changed, p)
		}
	}
	for p := range before {
		if _, ok := after[p]; !ok {
			changed = append(changed, p)
		}
	}
	sort.Strings(changed)
	return changed
}

func (s *server) waitForEvents(ctx context.Context, _ *mcp.CallToolRequest, in WaitIn) (*mcp.CallToolResult, WaitOut, error) {
	timeout := time.Duration(in.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = time.Minute
	}
	if timeout > s.maxWait {
		timeout = s.maxWait
	}
	// Anything that happened between two calls is behind the cursor already:
	// report it now rather than wait for the next change.
	before := s.snapshot()
	if events, omitted, err := s.catchUp(); err != nil {
		return nil, WaitOut{}, err
	} else if len(events) > 0 {
		return nil, WaitOut{Events: events, Omitted: omitted, Changed: []string{}}, nil
	}
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	tick := time.NewTicker(s.poll)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, WaitOut{Events: []Event{}, Changed: []string{}}, ctx.Err()
		case <-deadline.C:
			return nil, WaitOut{Events: []Event{}, Changed: []string{}, TimedOut: true}, nil
		case <-s.closing: // nil, and so never ready, unless the server was given one
			return nil, WaitOut{Events: []Event{}, Changed: []string{}, TimedOut: true}, nil
		case <-tick.C:
			if changed := diff(before, s.snapshot()); len(changed) > 0 {
				for i, p := range changed {
					changed[i] = s.rel(p)
				}
				// Paths say something changed (a thread, a narrative, this
				// agent's own write); events say what others did to work items.
				events, omitted, err := s.catchUp()
				if err != nil {
					return nil, WaitOut{}, err
				}
				return nil, WaitOut{Events: events, Omitted: omitted, Changed: changed}, nil
			}
		}
	}
}

func (s *server) rel(p string) string {
	root := s.repo.MainRoot
	if root == "" {
		root = s.repo.Root
	}
	if rel, err := filepath.Rel(root, p); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(p)
}
