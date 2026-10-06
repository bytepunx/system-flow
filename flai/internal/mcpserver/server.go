// Model Context Protocol (ADR-0020). Files stay the record: every tool reads
// and writes the same files, through the same code, as the flai CLI. The
// server is for agents, so it never accepts a story; that is the operator's.
package mcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/search"
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
	MaxWait time.Duration    // upper bound for one wait_for_events or wait_for_work call, default LongestWait
	// After arms the deadline of a held wait_for_events or wait_for_work
	// call; default time.After.
	After func(time.Duration) <-chan time.Time
	// Heartbeat is how often a held wait sends a progress notification to a
	// call that carries a progress token; default a minute.
	Heartbeat time.Duration
	// Runner runs git to ask what is accepted and not yet published
	// (ADR-0067), and whether a story's worktree is committed before it goes
	// to review (S-0140). Without one the server says nothing about the
	// first and skips the second.
	Runner execx.Runner
	// Closing, when closed, ends every held wait_for_events as if its time had
	// passed, so a server over HTTP can stop without cutting agents off
	// (S-0076). The SDK does not tie a session's call to its HTTP request.
	Closing <-chan struct{}
	// Rescan is how often a folder is looked through again for projects;
	// 5 seconds when zero.
	Rescan time.Duration
	// Logger receives one "request answered" event per request, at info
	// when it took Slow or longer (perf.Slow() when zero), at debug
	// otherwise (S-0152); nil logs nothing.
	Logger *slog.Logger
	Slow   time.Duration
	// Agents starts or restarts a story's agent on the host for the tools
	// agent_start and agent_restart (S-0177); without it they say how to on
	// the host.
	Agents AgentStart
	// Plans starts the planner for an epic or a story on the host for the
	// tool plan (S-0208); without it the tool says how to on the host.
	Plans PlanStart
	// Activities logs a strategic agent's activity for the tool activity_log
	// (S-0206); without it the tool says it cannot.
	Activities ActivityLog
	// AutoApprove, called at every permission_prompt request with the
	// project's root, lets it allow a write under .claude/ in an in-progress
	// story's worktree without asking the operator on a thread when it
	// returns true (S-0257); nil asks every time.
	AutoApprove AutoApprove
}

type server struct {
	repo    *workitem.Repo
	key     string // the manifest's key, else the folder's name
	agent   string
	role    string // the caller's role, as FLAI_ROLE says it (S-0219)
	now     func() time.Time
	poll    time.Duration
	maxWait time.Duration
	after   func(time.Duration) <-chan time.Time
	beat    time.Duration
	runner  execx.Runner
	closing <-chan struct{}
	agents  AgentStart
	plans   PlanStart // starts the planner (S-0208)
	// activities logs a strategic agent's activity (S-0206)
	activities ActivityLog
	version    string // the running flai's, compared with the project's newest flai tag (S-0181)
	// autoApprove lets permission_prompt allow without asking (S-0257)
	autoApprove AutoApprove
	logger      *slog.Logger // nil logs nothing

	// when wait_for_work last answered: a thread written to since then wakes it
	workMu    sync.Mutex
	workSince time.Time
}

// newServer is the server for one project, with the defaults filled in.
func newServer(opt Options, repo *workitem.Repo) *server {
	s := &server{repo: repo, agent: opt.Agent, now: opt.Now, poll: opt.Poll, maxWait: opt.MaxWait, after: opt.After, beat: opt.Heartbeat, runner: opt.Runner, closing: opt.Closing, agents: opt.Agents, plans: opt.Plans, activities: opt.Activities, version: opt.Version, autoApprove: opt.AutoApprove, logger: opt.Logger}
	if repo.Git == nil {
		repo.Git = opt.Runner // item_move asks git whether a story's worktree is committed (S-0140)
	}
	if s.agent == "" {
		s.agent = "agent"
	}
	// flai serve sets the orchestrator's role in its session, whose flai mcp
	// inherits it: the orchestrator's writes are held to its criteria (S-0219)
	s.role = os.Getenv("FLAI_ROLE")
	if s.now == nil {
		s.now = time.Now
	}
	if s.poll <= 0 {
		s.poll = 250 * time.Millisecond
	}
	if s.maxWait <= 0 {
		s.maxWait = LongestWait
	}
	if s.after == nil {
		s.after = time.After
	}
	if s.beat <= 0 {
		s.beat = heartbeat
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
		Instructions: "This server is the agent's view of a system-flow repository. Call inbox at the start of every turn or session, at every task transition, and before moving a story to review: it lists threads awaiting you, the stories ready to pull in pull order, and what others changed since you last looked (at most 50 changes, the newest; changes_omitted counts older ones that are not reported again; your first look covers the last 24 hours of stories and epics only, so use board and item_get for how things stand). Stories are yours to pull without being told. Whenever you have no story of your own in progress, call wait_for_work and do what it answers: pull the story it names (item_move it to in-progress, then flai stream open on the host), answer the threads it names, or go back to your own story. It answers as soon as a story is ready, the in-progress limit leaves room, and review is under its limit, and waits otherwise; when it times out, call it again, so that an idle agent is always waiting for the next story rather than stopping. An agent that ends its turn instead calls inbox when it starts again, and nothing in between is lost; wait_for_events reports every change, for an agent that wants the changes themselves. Reply to threads with thread_reply and ask the designer questions with thread_open. " + primeInstructions + " Commit everything in a story's worktree before you move it to review: item_move refuses a story whose worktree has uncommitted changes, because the operator cannot accept it. Stories are accepted by the operator only: item_move refuses to move a story or epic to done. A change of kind edited means someone changed an item's own words with flai edit or from the dashboard, and to names what (title, nature, tags, topics, touches, after, agent, parent, draft, cost_of_delay, forecast, goal, criteria, notes, body): if it is your story, read it again with item_get before you go on, because its criteria or its title may no longer be what you are working to. A change that says an epic moved following a story means that story's move took its epic along: nothing to do. A change that says an item was cancelled with a parent means the parent was cancelled and took it along: if it is your story or one of its tasks, stop work on it, log that in the narrative, and leave its branch and worktree alone. A change of kind overlapped means a story was accepted (cause) and changed paths (to) that an open story claims: if it is your story, run flai stream sync on it and the tests before you go on. When its cause is a story in progress, not an accepted one, a write grew the two stories' claims to overlap on those paths: coordinate with that story's agent on a thread before you change them, and narrow your touches if you can (I-0059). Publishing is the operator's: accepted work reaches the remote only when it is published (git fetch, then flai release --pending, or the board's Publish), and you publish only when the operator asks (ADR-0067); inbox's unpublished lists what is accepted and not yet published, for that. When inbox reports flai_outdated, the flai serving you is older than the newest flai release in the project's history and may lack rules the project relies on: tell the designer, who upgrades the host with the command it names, and go on.",
	})
	srv.AddReceivingMiddleware(timing(opt.Logger, nil, opt.Slow))
	mcp.AddTool(srv, &mcp.Tool{Name: "inbox", Description: inboxDescription}, route(one, (*server).inbox))
	addProjectTools(srv, one)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_events", Description: "Return what others changed since this agent last looked, at once when there is something already, otherwise block until a thread, work item, or narrative changes or the timeout passes: timeout_seconds, 60 by default and at most 1800 (30 minutes). Hold it when idle, as the orchestrator does between decisions, or to react within a second to the designer's answer on a thread awaiting them. It does not see a sub-agent finish: a story's agent waits for one by launching it with the Agent tool's run_in_background set to false, which returns the sub-agent's result as the tool's result, and flai guard refuses a story's agent this call while a sub-agent of its session runs and no thread on its story is open. At most 50 events, newest kept; events_omitted counts the rest."}, s.waitForEvents)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_work", Description: "What to do when you have nothing to work on (S-0097). Answers at once when there is something: reason resume with your own story still in progress; thread with threads awaiting you written to since it last answered; pull with the first ready story that is not held when the in-progress limit leaves room for it and review is under its limit (a story whose touches overlap a story in progress or in review, or that names in after a story not yet done, is held: it keeps its place and is offered once clear) (pull it: item_move it to in-progress, then flai stream open on the host; if item_move says it is already in-progress, another agent pulled it first: call wait_for_work again). Otherwise it waits until one of those is true, however long it takes, up to timeout_seconds; timed_out then says whether it is waiting for room (a story is ready, the in-progress limit is full), for review (a story is ready, review is full: no story is pulled until the operator accepts or sends one back), for a held story to be clear (held: every ready story is held, and ready says why each is), or for a story to be ready: call it again. Move your story to review first: while one of yours is in progress, it answers resume."}, s.waitForWork)
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
	Awaiting  string         `json:"awaiting" jsonschema:"'you' when the last entry is not yours, else 'other'; 'other' too on a thread you opened while a recommendation on it awaits the operator's confirmation, which is no answer until they confirm it"`
	// PendingRecommendation is the recommendation awaiting the operator's
	// confirmation (ADR-0090), nil when there is none.
	PendingRecommendation *threads.Entry `json:"pending_recommendation" jsonschema:"the recommendation awaiting the operator's confirmation, null when there is none: the thread still awaits the operator until they confirm it or answer otherwise"`
	Project               string         `json:"project,omitempty" jsonschema:"the project the thread is in, when the server serves more than one"`
}

// ThreadDetail is a thread with its entries.
type ThreadDetail struct {
	ThreadSummary
	Participants []string        `json:"participants"`
	EntryList    []threads.Entry `json:"entry_list"`
}

// summary is th as this agent sees it. It awaits the agent when the last
// entry is someone else's, but for the agent that opened it while a
// recommendation on it awaits the operator's confirmation: that is no answer
// to its question until the operator confirms it (ADR-0090). Anyone else,
// the operator and the recommendation's author included, sees it by its last
// entry.
func (s *server) summary(th *threads.Thread) ThreadSummary {
	entries := th.Entries()
	out := ThreadSummary{ID: th.ID, Title: th.Title, Status: th.Status, Anchor: th.Anchor, Story: threads.StoryOf(s.repo, th), Updated: th.Updated, Entries: len(entries), Awaiting: "other", PendingRecommendation: th.PendingRecommendation()}
	if n := len(entries); n > 0 {
		out.LastBy, out.LastEntry = entries[n-1].Author, entries[n-1].Text
		asked := out.PendingRecommendation != nil && th.Opener() == s.agent
		if out.LastBy != s.agent && !asked {
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
	Agent       string   `json:"agent"`
	Unpublished []string `json:"unpublished,omitempty" jsonschema:"the accepted items no release has covered yet, by ID: what publishing would send to the remote. Information for the operator, or for an agent the operator asks to publish; not a step to take on your own"`
	// FlaiOutdated is set while this flai is older than the newest flai
	// release in the project's history (S-0181).
	FlaiOutdated *release.Outdated    `json:"flai_outdated,omitempty" jsonschema:"this MCP server's flai is older than the newest flai release tagged in the project's history, so it may lack rules and fields the project uses; listed on every call while it is true: tell the designer, who upgrades the host with the command named"`
	AwaitingYou  int                  `json:"awaiting_you"`
	Threads      []ThreadSummary      `json:"threads"`
	Ready        []workitem.BoardCard `json:"ready" jsonschema:"stories ready to pull, in pull order; listed on every call"`
	CanPull      bool                 `json:"can_pull" jsonschema:"whether the in-progress limit leaves room and review is under its limit, so that one may be pulled"`
	PullHold     string               `json:"pull_hold,omitempty" jsonschema:"why no story may be pulled now: the in-progress limit is full, or review is full and waits on acceptance"`
	Changes      []Event              `json:"changes" jsonschema:"what others changed since this agent last looked, newest kept, at most 50; reported once. A first look under a new name covers the last 24 hours of stories and epics only"`
	Omitted      int                  `json:"changes_omitted" jsonschema:"how many older changes were left out of changes because of the cap; they are not reported later"`
}

func (s *server) inbox(ctx context.Context, _ *mcp.CallToolRequest, in InboxIn) (*mcp.CallToolResult, InboxOut, error) {
	done := perf.Track(ctx, "threads.read")
	all, err := threads.List(s.repo)
	done()
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
	// one list for the board and the changes (S-0156)
	items, err := s.listItems(ctx)
	if err != nil {
		return nil, InboxOut{}, err
	}
	view, err := s.boardViewOf(ctx, items, false)
	if err != nil {
		return nil, InboxOut{}, err
	}
	out.Ready, out.CanPull, out.PullHold, out.Unpublished = view.ReadyInPullOrder(), view.CanPull(), view.PullHold(), view.Unpublished
	out.FlaiOutdated = release.FlaiOutdated(s.runner, s.repo.Root, s.version)
	if out.Ready == nil {
		out.Ready = []workitem.BoardCard{}
	}
	done = perf.Track(ctx, "changes.read")
	out.Changes, out.Omitted, err = s.catchUpWith(items)
	done()
	if err != nil {
		return nil, InboxOut{}, err
	}
	return nil, out, nil
}

func (s *server) boardView(ctx context.Context, all bool) (workitem.BoardView, error) {
	items, err := s.listItems(ctx)
	if err != nil {
		return workitem.BoardView{}, err
	}
	return s.boardViewOf(ctx, items, all)
}

// listItems lists every item, archive included: a done story stays visible
// until it is published (S-0087), which NewBoardView decides from
// PendingIDs, not from the list.
func (s *server) listItems(ctx context.Context) ([]*workitem.Item, error) {
	defer perf.Track(ctx, "repo.list")()
	return s.repo.List(true)
}

// boardViewOf is the board of items listed with listItems.
func (s *server) boardViewOf(ctx context.Context, items []*workitem.Item, all bool) (workitem.BoardView, error) {
	done := perf.Track(ctx, "board.load")
	board, err := s.repo.LoadBoard()
	done()
	if err != nil {
		return workitem.BoardView{}, err
	}
	var ids map[string]bool // without a runner, nothing is known to be unpublished
	if s.runner != nil {
		done = perf.Track(ctx, "release.pending")
		ids = release.PendingIDs(execx.Timed(ctx, s.runner), s.repo.Root, s.repo.Manifest, s.repo)
		done()
	}
	return workitem.NewBoardView(items, board, s.now(), all, ids, s.repo.Manifest.Projects), nil
}

// BoardIn selects what the board shows.
type BoardIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	All     bool   `json:"all,omitempty" jsonschema:"include epics and tasks"`
}

func (s *server) board(ctx context.Context, _ *mcp.CallToolRequest, in BoardIn) (*mcp.CallToolResult, workitem.BoardView, error) {
	view, err := s.boardView(ctx, in.All)
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
	Project        string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID             string `json:"id"`
	Text           string `json:"text"`
	Recommendation bool   `json:"recommendation,omitempty" jsonschema:"the reply is a recommendation: the thread keeps its status and awaits the operator, who confirms it to make it the answer; refused on a resolved thread and on a thread you opened"`
	Source         string `json:"source,omitempty" jsonschema:"what the reply rests on: a repository path, such as an ADR, a design document, or a convention, or <path>#<heading> for a heading in it; both must exist. Written as the entry's last line, Source: <path> § <heading>"`
}

func (s *server) threadReply(ctx context.Context, _ *mcp.CallToolRequest, in ThreadReplyIn) (*mcp.CallToolResult, ThreadDetail, error) {
	source, err := threads.ParseSource(in.Source)
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	marks := threads.Marks{Recommendation: in.Recommendation, Source: source}
	th, err := threads.ReplyWith(s.repo, in.ID, s.agent, in.Text, s.now(), marks)
	if err != nil {
		return nil, ThreadDetail{}, err
	}
	if err := s.mirror(th); err != nil {
		return nil, ThreadDetail{}, err
	}
	if err := s.logReply(ctx, th, marks.Recommendation); err != nil {
		return nil, ThreadDetail{}, fmt.Errorf("the reply on %s is posted, but the orchestrator's decision log did not take it: %w; log it with activity_log", th.ID, err)
	}
	return nil, s.detail(th), nil
}

// logReply records the orchestrator's reply in its decision log, its activity
// document (S-0218): the thread, whether the reply recommends or answers, and
// the source it cites, with the thread's story as its item, so that the log
// never depends on the orchestrator calling activity_log. Any other caller's
// reply, or a server that cannot log activities, logs nothing.
func (s *server) logReply(ctx context.Context, th *threads.Thread, recommendation bool) error {
	if !s.orchestrator() || s.activities == nil {
		return nil
	}
	summary := "Answered " + th.ID
	if recommendation {
		summary = "Recommended an answer on " + th.ID
	}
	entries := th.Entries()
	if src := entries[len(entries)-1].Source; src != nil {
		summary += ", citing " + src.Path
		if src.Heading != "" {
			summary += " § " + src.Heading
		}
	} else {
		summary += ", citing no source"
	}
	items := []string{}
	if story := threads.StoryOf(s.repo, th); story != "" {
		items = append(items, story)
	}
	_, err := s.activities(ctx, projectRoot(s.repo), workitem.ActivityOrchestrator, summary, items, s.agent)
	return err
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
	Topics      []string              `json:"topics,omitempty" jsonschema:"a story's or epic's own: what it is about beyond the components it reaches"`
	After       []string              `json:"after,omitempty" jsonschema:"what it waits for until they are done: a story's, stories; a task's, tasks of its story"`
	Created     string                `json:"created"`
	Updated     string                `json:"updated"`
	Archived    bool                  `json:"archived"`
	Blocked     bool                  `json:"blocked"`
	Transitions []workitem.Transition `json:"transitions"`
	Path        string                `json:"path"`
	Body        string                `json:"body"`
	Children    []ItemBrief           `json:"children"`
	// Plan is a story's task plan, when it has tasks (S-0176).
	Plan *workitem.Plan `json:"plan,omitempty" jsonschema:"a story's tasks, when it has any: each task's state (ready, waiting, in-progress, done, cancelled), its after, and the tasks it waits for; and layers, the task IDs that can run at once, in the order they can run"`
	// Agent is a story's agent, and DefaultAgent the project's (S-0103).
	Agent        *manifest.Agent `json:"agent,omitempty"`
	DefaultAgent *manifest.Agent `json:"default_agent,omitempty"`
	// Draft, CostOfDelay, and Forecast are the item's planning data, and
	// Currency what its amounts are in (S-0199).
	Draft       bool                  `json:"draft,omitempty" jsonschema:"a story an agent wrote that the operator has not finalized: it cannot go to ready"`
	CostOfDelay *workitem.CostOfDelay `json:"cost_of_delay,omitempty" jsonschema:"what each week of waiting for the item costs: the inputs it is worked out from with who set them and when, and the value per week with who set it and when (ADR-0080)"`
	Forecast    *workitem.Forecast    `json:"forecast,omitempty" jsonschema:"a story's expected duration and delivery, what they rest on, and who set it and when"`
	Currency    string                `json:"currency,omitempty" jsonschema:"the currency of the cost of delay's amounts"`
	// Hash is what item_edit takes to refuse a change made meanwhile.
	Hash string `json:"hash"`
}

func brief(it *workitem.Item) ItemBrief {
	return ItemBrief{ID: it.ID, Type: it.Type, Title: it.Title, Status: it.Status, Touches: it.Touches}
}

func (s *server) itemOut(ctx context.Context, it *workitem.Item) (ItemOut, error) {
	items, err := s.listItems(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	data, err := os.ReadFile(it.Path)
	if err != nil {
		return ItemOut{}, err
	}
	out := ItemOut{Agent: it.Agent, DefaultAgent: s.repo.Manifest.Agent, Hash: docedit.Hash(string(data)), ItemBrief: brief(it), Nature: it.Nature, Parent: it.Parent, Owner: it.Owner, Tags: it.Tags, Topics: it.Topics, After: it.After, Created: it.Created, Updated: it.Updated, Archived: it.Archived, Transitions: it.Transitions, Path: s.rel(it.Path), Body: it.Body, Children: []ItemBrief{},
		Draft: it.Draft, CostOfDelay: it.CostOfDelay, Forecast: it.Forecast}
	if it.CostOfDelay != nil {
		out.Currency = s.repo.Manifest.Planning.CurrencyCode()
	}
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
	children := workitem.Children(items, it.ID)
	for _, c := range children {
		out.Children = append(out.Children, brief(c))
	}
	if it.Type == workitem.Story {
		out.Plan = workitem.PlanOf(children, it.ID)
	}
	return out, nil
}

func (s *server) itemGet(ctx context.Context, _ *mcp.CallToolRequest, in ItemIDIn) (*mcp.CallToolResult, ItemOut, error) {
	done := perf.Track(ctx, "repo.get")
	it, err := s.repo.Get(in.ID)
	done()
	if err != nil {
		return nil, ItemOut{}, err
	}
	out, err := s.itemOut(ctx, it)
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
	// Followed is the epic that moved with this story (S-0200).
	Followed *workitem.Followed `json:"followed,omitempty" jsonschema:"the epic that moved with this story: the state it was in (from) and the state it reached (to)"`
}

func (s *server) itemMove(_ context.Context, _ *mcp.CallToolRequest, in ItemMoveIn) (*mcp.CallToolResult, ItemMoveOut, error) {
	it, err := s.repo.Get(in.ID)
	if err != nil {
		return nil, ItemMoveOut{}, err
	}
	if in.To == workitem.Done && it.Type != workitem.Task {
		return nil, ItemMoveOut{}, fmt.Errorf("%s is %s: moving it to done is acceptance, which only the operator does (flai accept, or the dashboard); move it to review and say what is ready", it.ID, articled(it.Type))
	}
	if in.To == workitem.Ready && s.orchestrator() {
		// the orchestrator promotes a candidate alone, while ready has room,
		// as flai move holds it (S-0219)
		if err := s.repo.Promotable(it); err != nil {
			return nil, ItemMoveOut{}, err
		}
	}
	if in.To == workitem.Ready && it.Type == workitem.Story && it.Draft {
		return nil, ItemMoveOut{}, fmt.Errorf("%s is a draft: finalizing it, which lets it go to ready, is the operator's; say in a thread or your narrative that it is ready to be finalized", it.ID)
	}
	// finalize is false: a move finalizes no draft; the orchestrator
	// finalizes a complete one with item_edit (S-0219)
	res, err := s.repo.TransitionAll(it, in.To, s.agent, in.Reason, s.now(), false)
	if err != nil {
		return nil, ItemMoveOut{}, err
	}
	warnings := res.Warnings
	if warnings == nil {
		warnings = []string{}
	}
	return nil, ItemMoveOut{ID: it.ID, Status: it.Status, Warnings: warnings, Cancelled: res.Cancelled, Followed: res.Followed}, nil
}

// orchestrator says whether the caller is the orchestrator (S-0219).
func (s *server) orchestrator() bool { return s.role == guard.RoleOrchestrate }

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
	Heading string `json:"heading,omitempty" jsonschema:"a heading in the document, by its text, a heading path such as Commands › flai prime, or its anchor slug: only that section and the sections below it are returned"`
}

// DocOut is a document split at its front matter, or one section of it.
type DocOut struct {
	Path        string `json:"path"`
	FrontMatter string `json:"front_matter" jsonschema:"the YAML front matter, without the --- fences"`
	Heading     string `json:"heading,omitempty" jsonschema:"with a heading asked for, the heading path of the section found; body is then that section and the sections below it"`
	Line        int    `json:"line,omitempty" jsonschema:"with a heading asked for, the line of the section's heading in the file"`
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
	if in.Heading != "" {
		part, err := ctxpack.Section(out.Path, string(data), in.Heading)
		if err != nil {
			return nil, DocOut{}, err
		}
		out.Heading, out.Line, out.Body = part.Heading, part.Line, part.Text
	}
	return nil, out, nil
}

// DocSearchIn is a search of the design and docs folders.
type DocSearchIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Query   string `json:"query" jsonschema:"words to rank sections by; stopwords and single letters are left out"`
	Limit   int    `json:"limit,omitempty" jsonschema:"at most this many sections; default and cap 20"`
}

// DocSearchOut is the sections a search ranked highest.
type DocSearchOut struct {
	Query string              `json:"query"`
	Hits  []search.SectionHit `json:"hits"`
}

func (s *server) docSearch(_ context.Context, _ *mcp.CallToolRequest, in DocSearchIn) (*mcp.CallToolResult, DocSearchOut, error) {
	hits, err := ctxpack.SearchSections(s.repo, in.Query, in.Limit)
	if err != nil {
		return nil, DocSearchOut{}, err
	}
	return nil, DocSearchOut{Query: in.Query, Hits: hits}, nil
}

// ---- context pack ----

// PrimeIn names the story, or the role, to prime for.
type PrimeIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string `json:"story,omitempty" jsonschema:"story ID such as S-0138 (any zero padding); an archived story gets the pack it would get today. Needed but for role orchestrate or analyze, and for role plan given an epic"`
	Epic    string `json:"epic,omitempty" jsonschema:"with role plan, the epic the planner plans, such as E-0016, instead of a story"`
	Budget  string `json:"budget,omitempty" jsonschema:"the size the pack fits, such as 80KB; default the project's prime.budget, else 80KB, and half that with role explore or verify"`
	Role    string `json:"role,omitempty" jsonschema:"explore or verify: the smaller pack for a sub-agent of the story's agent in that role (ADR-0059); plan, orchestrate, or analyze: the pack for a strategic agent, the planner for an epic or a story, or the orchestrator or the analyzer for the whole project; empty for the story's agent's own pack"`
}

func (s *server) prime(_ context.Context, _ *mcp.CallToolRequest, in PrimeIn) (*mcp.CallToolResult, *ctxpack.Pack, error) {
	switch {
	case in.Epic != "" && in.Role != conventions.RolePlan:
		return nil, nil, fmt.Errorf("prime: epic primes the planner for an epic; give role plan too")
	case slices.Contains(conventions.StrategicRoles, in.Role):
		pack, err := ctxpack.ForStrategic(s.repo, in.Role, in.Epic, in.Story, in.Budget)
		return nil, pack, err
	case in.Story == "":
		return nil, nil, fmt.Errorf("prime: give story, or role %s", strings.Join(conventions.StrategicRoles, ", "))
	}
	if in.Role != "" {
		pack, err := ctxpack.ForRole(s.repo, in.Story, in.Role, in.Budget)
		return nil, pack, err
	}
	pack, err := ctxpack.ForStory(s.repo, in.Story, in.Budget)
	return nil, pack, err
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
	TimeoutSeconds int `json:"timeout_seconds,omitempty" jsonschema:"how long to wait in seconds, default 60, at most 1800 (30 minutes)"`
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

func (s *server) waitForEvents(ctx context.Context, req *mcp.CallToolRequest, in WaitIn) (*mcp.CallToolResult, WaitOut, error) {
	timeout := holdFor(in.TimeoutSeconds, eventsWait, s.maxWait)
	// Anything that happened between two calls is behind the cursor already:
	// report it now rather than wait for the next change.
	before := s.snapshot()
	if events, omitted, err := s.catchUp(ctx); err != nil {
		return nil, WaitOut{}, err
	} else if len(events) > 0 {
		return nil, WaitOut{Events: events, Omitted: omitted, Changed: []string{}}, nil
	}
	deadline := s.after(timeout)
	defer keepAlive(ctx, req, s.beat)()
	tick := time.NewTicker(s.poll)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, WaitOut{Events: []Event{}, Changed: []string{}}, ctx.Err()
		case <-deadline:
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
				events, omitted, err := s.catchUp(ctx)
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
func (in DocSearchIn) project() string     { return in.Project }
func (in PrimeIn) project() string         { return in.Project }
