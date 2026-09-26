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
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Options configure the server.
type Options struct {
	// Repo is the project served. Leave it nil and set Folder to serve every
	// system-flow project in a folder and below it instead (S-0101).
	Repo    *workitem.Repo
	Folder  string
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
	// Rescan is how often a folder is looked through again for projects;
	// 5 seconds when zero.
	Rescan time.Duration
}

type server struct {
	repo    *workitem.Repo
	key     string // the manifest's key, else the folder's name
	agent   string
	now     func() time.Time
	poll    time.Duration
	maxWait time.Duration
	runner  execx.Runner
	closing <-chan struct{}

	// when wait_for_work last answered: a thread written to since then wakes it
	workMu    sync.Mutex
	workSince time.Time
}

// newServer is the server for one project, with the defaults filled in.
func newServer(opt Options, repo *workitem.Repo) *server {
	s := &server{repo: repo, agent: opt.Agent, now: opt.Now, poll: opt.Poll, maxWait: opt.MaxWait, runner: opt.Runner, closing: opt.Closing}
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
	s.key = repo.Manifest.Key
	if s.key == "" {
		s.key = filepath.Base(projectRoot(repo))
	}
	return s
}

// projectRoot is the main checkout's root, where wip lives.
func projectRoot(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

// New builds the MCP server with its tools and resources: for one project,
// or, with Options.Folder and no Repo, for every project in a folder (S-0101).
func New(opt Options) *mcp.Server {
	if opt.Repo == nil {
		return newFolderServer(opt)
	}
	s := newServer(opt, opt.Repo)
	one := single{s}
	srv := mcp.NewServer(&mcp.Implementation{Name: "flai", Title: "system-flow repository", Version: opt.Version}, &mcp.ServerOptions{
		Instructions: "This server is the agent's view of a system-flow repository. Call inbox at the start of every turn or session, at every task transition, and before moving a story to review: it lists threads awaiting you, the stories ready to pull in pull order, and what others changed since you last looked (at most 50 changes, the newest; changes_omitted counts older ones that are not reported again; your first look covers the last 24 hours of stories and epics only, so use board and item_get for how things stand). Stories are yours to pull without being told. Whenever you have no story of your own in progress, call wait_for_work and do what it answers: pull the story it names (item_move it to in-progress, then flai stream open on the host), answer the threads it names, or go back to your own story. It answers as soon as a story is ready and the in-progress limit leaves room, and waits otherwise; when it times out, call it again, so that an idle agent is always waiting for the next story rather than stopping. An agent that ends its turn instead calls inbox when it starts again, and nothing in between is lost; wait_for_events reports every change, for an agent that wants the changes themselves. Reply to threads with thread_reply and ask the designer questions with thread_open. Stories are accepted by the operator only: item_move refuses to move a story or epic to done. A change of kind edited means someone changed an item's own words with flai edit or from the dashboard, and to names what (title, nature, tags, touches, agent, parent, goal, criteria, notes, body): if it is your story, read it again with item_get before you go on, because its criteria or its title may no longer be what you are working to. A change that says an item was cancelled with a parent means the parent was cancelled and took it along: if it is your story or one of its tasks, stop work on it, log that in the narrative, and leave its branch and worktree alone. A change of kind overlapped means a story was accepted (cause) and changed paths (to) that an open story claims: if it is your story, run flai stream sync on it and the tests before you go on. When inbox reports unpushed, an acceptance was made where nothing could push it: on the host run git fetch, then flai push --pending, before anything else; it never forces, and if it refuses because the remote moved, merge and run it again.",
	})
	mcp.AddTool(srv, &mcp.Tool{Name: "inbox", Description: inboxDescription}, route(one, (*server).inbox))
	addProjectTools(srv, one)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_events", Description: "Return what others changed since this agent last looked, at once when there is something already, otherwise block until a thread, work item, or narrative changes or the timeout passes. Hold this when idle to react to the designer within a second. At most 50 events, newest kept; events_omitted counts the rest."}, s.waitForEvents)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_work", Description: "What to do when you have nothing to work on (S-0097). Answers at once when there is something: reason resume with your own story still in progress; thread with threads awaiting you written to since it last answered; pull with the first ready story that is not held when the in-progress limit leaves room for it (a story whose touches overlap a story in progress or in review, or that names in after a story not yet done, is held: it keeps its place and is offered once clear) (pull it: item_move it to in-progress, then flai stream open on the host; if item_move says it is already in-progress, another agent pulled it first: call wait_for_work again). Otherwise it waits until one of those is true, however long it takes, up to timeout_seconds; timed_out then says whether it is waiting for room (a story is ready, the limit is full), for a held story to be clear (held: every ready story is held, and ready says why each is), or for a story to be ready: call it again. Move your story to review first: while one of yours is in progress, it answers resume."}, s.waitForWork)
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
	Project   string         `json:"project,omitempty" jsonschema:"the project the thread is in, when the server serves more than one"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story,omitempty" jsonschema:"only threads on this story and its tasks"`
	All     bool   `json:"all,omitempty" jsonschema:"include threads awaiting someone else"`
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
	// true: a done story stays visible until it is published (S-0087), which
	// NewBoardView decides from PendingIDs, not from this flag.
	items, err := s.repo.List(true)
	if err != nil {
		return workitem.BoardView{}, err
	}
	board, err := s.repo.LoadBoard()
	if err != nil {
		return workitem.BoardView{}, err
	}
	view := workitem.NewBoardView(items, board, s.now(), all, release.PendingIDs(s.runner, s.repo.Root, s.repo.Manifest, s.repo), s.repo.Manifest.Projects)
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	All     bool   `json:"all,omitempty" jsonschema:"include epics and tasks"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"thread ID such as TH-0001"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id"`
	Text    string `json:"text"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id"`
	Reason  string `json:"reason,omitempty" jsonschema:"what settled it"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"work item ID in any zero padding, such as S-39 or S-0039"`
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
	After       []string              `json:"after,omitempty" jsonschema:"a story's: the stories it waits for until they are done"`
	Created     string                `json:"created"`
	Updated     string                `json:"updated"`
	Archived    bool                  `json:"archived"`
	Blocked     bool                  `json:"blocked"`
	Transitions []workitem.Transition `json:"transitions"`
	Path        string                `json:"path"`
	Body        string                `json:"body"`
	Children    []ItemBrief           `json:"children"`
	// Agent is a story's agent, and DefaultAgent the project's (S-0103).
	Agent        *manifest.Agent `json:"agent,omitempty"`
	DefaultAgent *manifest.Agent `json:"default_agent,omitempty"`
	// Hash is what item_edit takes to refuse a change made meanwhile.
	Hash string `json:"hash"`
}

func brief(it *workitem.Item) ItemBrief {
	return ItemBrief{ID: it.ID, Type: it.Type, Title: it.Title, Status: it.Status, Touches: it.Touches}
}

func (s *server) itemOut(it *workitem.Item) (ItemOut, error) {
	items, err := s.repo.List(true)
	if err != nil {
		return ItemOut{}, err
	}
	data, err := os.ReadFile(it.Path)
	if err != nil {
		return ItemOut{}, err
	}
	out := ItemOut{Agent: it.Agent, DefaultAgent: s.repo.Manifest.Agent, Hash: docedit.Hash(string(data)), ItemBrief: brief(it), Nature: it.Nature, Parent: it.Parent, Owner: it.Owner, Tags: it.Tags, After: it.After, Created: it.Created, Updated: it.Updated, Archived: it.Archived, Transitions: it.Transitions, Path: s.rel(it.Path), Body: it.Body, Children: []ItemBrief{}}
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id"`
	To      string `json:"to" jsonschema:"one of backlog, ready, in-progress, review, done, cancelled"`
	Reason  string `json:"reason,omitempty" jsonschema:"required for cancelled and for review back to in-progress"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Path    string `json:"path" jsonschema:"repository path or component name"`
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
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Path    string `json:"path" jsonschema:"repository path of a markdown file under design, docs, or wip"`
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
	// an edit's notice is written after its files, and an overlap's has no
	// file of its own: a waiting agent wakes for them too
	for _, log := range noticeLogs(s.repo) {
		if info, err := os.Stat(log); err == nil {
			out[log] = fmt.Sprintf("%d/%d", info.ModTime().UnixNano(), info.Size())
		}
	}
	return out
}

// noticeLogs are the logs under .flai-cache whose lines arrive as events:
// edits (S-0085) and overlaps at acceptance (S-0132).
func noticeLogs(repo *workitem.Repo) []string {
	return []string{itemedit.NoticesPath(repo), itemedit.OverlapsPath(repo)}
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
				// the notices wake a waiting agent and are not a path of the
				// repository to show it: what they say arrives as events
				shown := changed[:0]
				logs := map[string]bool{}
				for _, log := range noticeLogs(s.repo) {
					logs[log] = true
				}
				for _, p := range changed {
					if !logs[p] {
						shown = append(shown, s.rel(p))
					}
				}
				changed = shown
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

// Each tool input names its project, when the server serves more than one.
func (in InboxIn) project() string         { return in.Project }
func (in BoardIn) project() string         { return in.Project }
func (in ThreadIDIn) project() string      { return in.Project }
func (in ThreadOpenIn) project() string    { return in.Project }
func (in ThreadReplyIn) project() string   { return in.Project }
func (in ThreadResolveIn) project() string { return in.Project }
func (in ItemIDIn) project() string        { return in.Project }
func (in ItemMoveIn) project() string      { return in.Project }
func (in WhoTouchesIn) project() string    { return in.Project }
func (in DocIn) project() string           { return in.Project }
