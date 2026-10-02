package mcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// primeInstructions tells an agent to prime with its story's pack and to fetch
// what the pack briefs (ADR-0049, S-0148).
const primeInstructions = "When you start work on a story, call prime with it: the pack holds the conventions that apply and what the story names, whole, and briefs the design and ADRs its topics and links select, within a size budget. A brief is not the document: when one bears on the story, read it, or its section that does, with doc_get and its heading before relying on it or changing what it describes, and find sections by their words with doc_search."

const inboxDescription = "What needs this agent: unresolved threads (awaiting is 'you' when the last entry is not yours), the stories ready to pull in pull order with can_pull from the in-progress limit (one whose touches overlap a story in progress or in review, or that names in after a story not yet done, carries held with the reason and what clears it: it is not offered), and the changes others made to work items since this agent last looked, reported once: at most 50, newest kept, with changes_omitted counting the older ones left out. A first look covers 24 hours of stories and epics only. Filter by story to see only threads on a story and its tasks."

// projects is what the tools are served for: one project, or every project
// in a folder (S-0101).
type projects interface {
	all() []*server
	// pick is the project a tool call names: by key or folder, or none when
	// there is only one.
	pick(name string) (*server, error)
}

type single struct{ s *server }

func (o single) all() []*server { return []*server{o.s} }

func (o single) pick(name string) (*server, error) {
	if name == "" || name == o.s.key || name == projectRoot(o.s.repo) {
		return o.s, nil
	}
	return nil, fmt.Errorf("this server serves only %s (%s), not %q", o.s.key, projectRoot(o.s.repo), name)
}

// scoped is a tool input that can name its project.
type scoped interface{ project() string }

// route serves a per-project tool for the project the call names.
func route[In scoped, Out any](p projects, h func(*server, context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)) mcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, Out, error) {
		s, err := p.pick(in.project())
		if err != nil {
			var zero Out
			return nil, zero, err
		}
		return h(s, ctx, req, in)
	}
}

// addProjectTools registers the tools that act on one project.
func addProjectTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_get", Description: "One thread with all of its dated entries."}, route(p, (*server).threadGet))
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_open", Description: "Open a thread on a repository path (optionally a heading in it) or a work item ID, to ask the designer a question or record a discussion."}, route(p, (*server).threadOpen))
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_reply", Description: "Add an entry to a thread as this agent. A reply from anyone but the opener marks the thread answered."}, route(p, (*server).threadReply))
	mcp.AddTool(srv, &mcp.Tool{Name: "thread_resolve", Description: "Close a thread, optionally saying what settled it."}, route(p, (*server).threadResolve))
	mcp.AddTool(srv, &mcp.Tool{Name: "item_get", Description: "A work item by ID (any zero padding): front matter, body, and children."}, route(p, (*server).itemGet))
	mcp.AddTool(srv, &mcp.Tool{Name: "item_new", Description: "Create a work item (S-0103): a story (with an epic or none), a task under its story, or an epic. A story gets the project's default agent (harness, model, config, and roles for its sub-agents) with anything given in agent winning. after names what it waits for: a story's stories, a task's tasks of the same story; flai check runs with it and refuses an entry that does not exist, a task of another story, or a cycle. Written to the repository as flai story new writes it; nothing is committed."}, route(p, (*server).itemNew))
	mcp.AddTool(srv, &mcp.Tool{Name: "item_edit", Description: "Change an item's own words (S-0103): title, nature, tags, touches, parent, body, a story's or epic's topics (what it is about beyond its components; an empty list removes them), a story's or a task's after (what it waits for: a story's stories, a task's tasks of the same story; an empty list removes them), and a story's agent (agent replaces it; clear_agent removes it). Only what is given changes. Give the hash item_get returned so that a change made meanwhile is refused instead of overwritten. Nothing is committed."}, route(p, (*server).itemEdit))
	mcp.AddTool(srv, &mcp.Tool{Name: "item_move", Description: "Transition a work item with the workflow rules enforced. Moving an item to cancelled also cancels everything open under it, and the result lists what went with it. Refuses to move a story or epic to done: acceptance is the operator's. Refuses to move a story to review while its worktree has uncommitted changes: commit them first."}, route(p, (*server).itemMove))
	mcp.AddTool(srv, &mcp.Tool{Name: "doc_get", Description: "A markdown document under the design, docs, or wip folders, by repository path: its front matter and its body, or, given a heading, only that section and the sections below it (ADR-0049). A heading is its text, a heading path such as Commands › flai prime, or its anchor slug; one that names no section, or more than one, is refused with the headings to choose from."}, route(p, (*server).docGet))
	mcp.AddTool(srv, &mcp.Tool{Name: "doc_search", Description: "Search the sections of the design and docs folders, the conventions among them, by BM25 against a query (ADR-0049): at most 20, best first, each with its path, the document's title, its heading path below the title, its line, its first lines, and its size. Sections are cut as the prime pack cuts them, each down to the next heading of any level. Fetch one with doc_get and its heading."}, route(p, (*server).docSearch))
	mcp.AddTool(srv, &mcp.Tool{Name: "prime", Description: "The context pack for a story (ADR-0047), fitted to a size budget (ADR-0049), as flai prime --story --json prints it: the story's topics and where each came from; the budget, the pack's size, and exceeded naming the part that takes it over (conventions, named, or briefs; nothing is ranked then, and briefs are never cut); each convention the story's agent reads, those whose roles are empty or list story (ADR-0068), with the sections its topics leave out taken out, and one line per section left out; the number of open issues; then items, each with its step, reason, size, and text: what the story, its epic, and its tasks name, whole (named), except a document named only by its path written out and larger than an eighth of the budget, briefed with the reason named in <ID> (ADR-0050); a brief of each design and tech file its topics select, as its first paragraph and heading outline, and of each ADR its topics or one link step reach, as its decision sentence (briefed, with whole giving the document's size); and the sections ranked highest against the story, cut at their own headings, while the budget has room (ranked); and a catalog of the documents neither loaded nor briefed. Read a briefed document, or the section of it that bears on the story, with doc_get and its heading before changing what it describes, and find sections by their words with doc_search. budget is a size such as 80KB; the default is the project's prime.budget, else 80KB. A sub-agent the story's agent hands work to gives role, explore or verify (ADR-0059), and gets a smaller pack, half the budget: the conventions whose roles are empty or list it, the story's goal and acceptance criteria (goal), and briefs, never bodies, of what the story names and what its topics and one link step select, while the budget has room (briefs_left_out counts the rest); no issues, nothing ranked, no catalog."}, route(p, (*server).prime))
	mcp.AddTool(srv, &mcp.Tool{Name: "who_touches", Description: "In-progress and in-review items whose touches cover a path; ask before editing a path someone else is working on."}, route(p, (*server).whoTouches))
	mcp.AddTool(srv, &mcp.Tool{Name: "agent_start", Description: agentStartDescription}, route(p, (*server).agentStart))
	mcp.AddTool(srv, &mcp.Tool{Name: "agent_restart", Description: agentRestartDescription}, route(p, (*server).agentRestart))
	mcp.AddTool(srv, &mcp.Tool{Name: "board", Description: "The kanban board as flai board --json prints it: cards per column, WIP limits, the pull order, and limit breaches; a ready story whose touches overlap a story in progress or in review, or that names in after a story not yet done, carries held with the reason and what clears it. Stories only unless all is set."}, route(p, (*server).board))
}

// ---- a folder of projects (S-0101) ----

// folder serves every project in a folder, found again every Rescan, so that
// one created or imported while the agent works joins without a restart.
type folder struct {
	root   string
	opt    Options
	rescan time.Duration

	mu      sync.Mutex
	byRoot  map[string]*server
	list    []*server
	scanned time.Time
	refused map[string]string // why a project was left out, as last logged

	// when wait_for_work last answered, across every project
	workMu    sync.Mutex
	workSince time.Time
}

func (f *folder) all() []*server {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.scanned.IsZero() && time.Since(f.scanned) < f.rescan {
		return f.list
	}
	f.scanned = time.Now()
	byRoot := map[string]*server{}
	var list []*server
	for _, root := range workitem.FindProjects(f.root) {
		s := f.byRoot[root]
		if s == nil {
			repo, err := workitem.Open(root)
			if err != nil {
				f.leaveOut(root, err)
				continue
			}
			s = newServer(f.opt, repo)
		}
		byRoot[root] = s
		list = append(list, s)
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].key < list[j].key })
	f.byRoot, f.list = byRoot, list
	return list
}

// leaveOut logs once why a project under the folder is not served, such as
// a minimum flai above this one (S-0181). The caller holds f.mu.
func (f *folder) leaveOut(root string, err error) {
	if f.refused == nil {
		f.refused = map[string]string{}
	}
	if f.refused[root] == err.Error() {
		return
	}
	f.refused[root] = err.Error()
	log := f.opt.Logger
	if log == nil {
		log = slog.Default()
	}
	log.Warn("project not served", "component", "mcp", "root", root, "err", err.Error())
}

func (f *folder) pick(name string) (*server, error) {
	list := f.all()
	if len(list) == 0 {
		return nil, fmt.Errorf("there is no system-flow project in %s or below it", f.root)
	}
	if name == "" {
		if len(list) == 1 {
			return list[0], nil
		}
		return nil, fmt.Errorf("this server serves %d projects; name one with project: %s", len(list), f.names(list))
	}
	for _, s := range list {
		root := projectRoot(s.repo)
		if rel, err := filepath.Rel(f.root, root); name == s.key || name == root || (err == nil && name == filepath.ToSlash(rel)) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("no project %q here; the projects are: %s", name, f.names(list))
}

func (f *folder) names(list []*server) string {
	out := make([]string, len(list))
	for i, s := range list {
		out[i] = fmt.Sprintf("%s (%s)", s.key, f.rel(projectRoot(s.repo)))
	}
	return strings.Join(out, ", ")
}

// rel is a path relative to the folder, as the agent sees the projects.
func (f *folder) rel(p string) string {
	if rel, err := filepath.Rel(f.root, p); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(p)
}

func newFolderServer(opt Options) *mcp.Server {
	f := &folder{root: filepath.Clean(opt.Folder), opt: opt, rescan: opt.Rescan}
	if f.rescan <= 0 {
		f.rescan = 5 * time.Second
	}
	poll, maxWait, now := opt.Poll, opt.MaxWait, opt.Now
	if poll <= 0 {
		poll = 250 * time.Millisecond
	}
	if maxWait <= 0 {
		maxWait = 5 * time.Minute
	}
	if now == nil {
		now = time.Now
	}
	fw := &folderWaits{f: f, poll: poll, maxWait: maxWait, now: now, closing: opt.Closing}
	srv := mcp.NewServer(&mcp.Implementation{Name: "flai", Title: "system-flow projects in " + f.root, Version: opt.Version}, &mcp.ServerOptions{
		Instructions: "This server is the agent's view of every system-flow project in " + f.root + " and the folders below it (S-0101): it was started in a folder that is not itself a project. inbox, wait_for_work, and wait_for_events cover every project and say which one each thing is in; every other tool takes project, a key inbox lists (or the project's folder), and needs it whenever there is more than one project. Projects created or imported below the folder join within seconds. Call inbox at the start of every turn or session, at every task transition, and before moving a story to review. Stories are yours to pull without being told. Whenever you have no story of your own in progress, call wait_for_work and do what it answers: pull the story it names in the project it names (item_move it to in-progress with that project, then flai stream open in that project's folder on the host), answer the threads it names, or go back to your own story. It answers as soon as a story is ready in some project whose in-progress limit leaves room, and waits otherwise; when it times out, call it again. Reply to threads with thread_reply and ask the designer questions with thread_open. " + primeInstructions + " Commit everything in a story's worktree before you move it to review: item_move refuses a story whose worktree has uncommitted changes, because the operator cannot accept it. Stories are accepted by the operator only: item_move refuses to move a story or epic to done. A change of kind edited means someone changed an item's own words: if it is your story, read it again with item_get before you go on. A change that says an item was cancelled with a parent means the parent was cancelled and took it along: if it is your story or one of its tasks, stop work on it, log that in the narrative, and leave its branch and worktree alone. A change of kind overlapped means a story was accepted (cause) and changed paths (to) that an open story claims: if it is your story, run flai stream sync on it and the tests before you go on. When a project's inbox reports unpushed, an acceptance was made where nothing could push it: in that project's folder on the host run git fetch, then flai push --pending, before anything else.",
	})
	srv.AddReceivingMiddleware(timing(opt.Logger, nil, opt.Slow))
	mcp.AddTool(srv, &mcp.Tool{Name: "inbox", Description: "What needs this agent in every project here, one entry per project with its key, name, and folder: " + inboxDescription + " Name project to see only that one; story needs project when there is more than one."}, fw.inbox)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_work", Description: "What to do when you have nothing to work on, across every project here (S-0097, S-0101). Answers at once with reason resume and your own story still in progress, in whichever project; thread with threads awaiting you written to since it last answered, each with its project; pull with the first ready story that is not held, in key order of projects and pull order within one, of a project whose in-progress limit leaves room (pull it: item_move it to in-progress with that project, then flai stream open in its folder on the host; if item_move says it is already in-progress, another agent pulled it first: call wait_for_work again). Otherwise it waits until one of those is true, up to timeout_seconds; timed_out then says whether it is waiting for room, for a held story to be clear (held: every ready story with room is held, and ready says why each is), or for a story to be ready: call it again."}, fw.work)
	mcp.AddTool(srv, &mcp.Tool{Name: "wait_for_events", Description: "Return what others changed in any project here since this agent last looked, at once when there is something already, otherwise block until a thread, work item, or narrative changes in any of them or the timeout passes. Each event names its project; changed paths are relative to the folder. At most 50 events per project, newest kept; events_omitted counts the rest."}, fw.events)
	addProjectTools(srv, f)
	return srv
}

// folderWaits are the tools that look across every project in a folder.
type folderWaits struct {
	f       *folder
	poll    time.Duration
	maxWait time.Duration
	now     func() time.Time
	closing <-chan struct{}
}

// ProjectInbox is one project's inbox, in a folder of them.
type ProjectInbox struct {
	Project string `json:"project" jsonschema:"the key the other tools take as project"`
	Name    string `json:"name"`
	Folder  string `json:"folder" jsonschema:"the project's folder, relative to the one this server serves"`
	InboxOut
}

// FolderInboxOut is the agent's inbox across every project in a folder.
type FolderInboxOut struct {
	Agent       string         `json:"agent"`
	Folder      string         `json:"folder"`
	AwaitingYou int            `json:"awaiting_you" jsonschema:"threads awaiting you, in every project"`
	Projects    []ProjectInbox `json:"projects"`
}

func (fw *folderWaits) inbox(ctx context.Context, req *mcp.CallToolRequest, in InboxIn) (*mcp.CallToolResult, FolderInboxOut, error) {
	list := fw.f.all()
	out := FolderInboxOut{Agent: fw.agent(), Folder: fw.f.root, Projects: []ProjectInbox{}}
	if in.Project != "" {
		s, err := fw.f.pick(in.Project)
		if err != nil {
			return nil, FolderInboxOut{}, err
		}
		list = []*server{s}
	} else if in.Story != "" && len(list) > 1 {
		return nil, FolderInboxOut{}, fmt.Errorf("a story ID is one project's: name the project too (%s)", fw.f.names(list))
	}
	for _, s := range list {
		_, one, err := s.inbox(ctx, req, in)
		if err != nil {
			return nil, FolderInboxOut{}, fmt.Errorf("%s: %w", s.key, err)
		}
		for i := range one.Threads {
			one.Threads[i].Project = s.key
		}
		for i := range one.Changes {
			one.Changes[i].Project = s.key
		}
		out.AwaitingYou += one.AwaitingYou
		out.Projects = append(out.Projects, ProjectInbox{Project: s.key, Name: s.repo.Manifest.Name, Folder: fw.f.rel(projectRoot(s.repo)), InboxOut: one})
	}
	return nil, out, nil
}

func (fw *folderWaits) agent() string {
	if fw.f.opt.Agent != "" {
		return fw.f.opt.Agent
	}
	return "agent"
}

// ProjectWork is one project's ready stories and room, for wait_for_work.
type ProjectWork struct {
	Project string               `json:"project"`
	Folder  string               `json:"folder"`
	Ready   []workitem.BoardCard `json:"ready"`
	CanPull bool                 `json:"can_pull"`
}

// FolderWorkOut is WorkOut across every project in a folder.
type FolderWorkOut struct {
	Reason     string              `json:"reason" jsonschema:"resume, thread, pull, or empty when it timed out"`
	Project    string              `json:"project,omitempty" jsonschema:"the project of the story to resume or pull"`
	Folder     string              `json:"folder,omitempty" jsonschema:"that project's folder, where flai stream open runs"`
	Story      *workitem.BoardCard `json:"story,omitempty"`
	Threads    []ThreadSummary     `json:"threads" jsonschema:"threads awaiting you written to since wait_for_work last answered, each with its project"`
	Projects   []ProjectWork       `json:"projects" jsonschema:"every project's ready stories and whether its limit leaves room"`
	WaitingFor string              `json:"waiting_for,omitempty"`
	TimedOut   bool                `json:"timed_out"`
}

// decide is what an idle agent should do across every project, if anything.
func (fw *folderWaits) decide(ctx context.Context, since time.Time) (FolderWorkOut, error) {
	out := FolderWorkOut{Threads: []ThreadSummary{}, Projects: []ProjectWork{}}
	var pull *FolderWorkOut
	anyReady, anyRoom := false, false
	for _, s := range fw.f.all() {
		one, err := s.work(ctx, since)
		if err != nil {
			return FolderWorkOut{}, fmt.Errorf("%s: %w", s.key, err)
		}
		folder := fw.f.rel(projectRoot(s.repo))
		out.Projects = append(out.Projects, ProjectWork{Project: s.key, Folder: folder, Ready: one.Ready, CanPull: one.CanPull})
		anyReady = anyReady || len(one.Ready) > 0
		anyRoom = anyRoom || (one.CanPull && len(one.Ready) > 0)
		switch one.Reason {
		case WorkResume:
			res := out
			res.Reason, res.Project, res.Folder, res.Story = WorkResume, s.key, folder, one.Story
			return res, nil
		case WorkThread:
			for _, th := range one.Threads {
				th.Project = s.key
				out.Threads = append(out.Threads, th)
			}
		}
		// work() puts threads ahead of a pull; across projects a pull waits for
		// every project's threads to be looked at first
		if card := workitem.FirstClear(one.Ready); pull == nil && one.CanPull && card != nil {
			pull = &FolderWorkOut{Project: s.key, Folder: folder, Story: card}
		}
	}
	switch {
	case len(out.Threads) > 0:
		out.Reason = WorkThread
	case pull != nil:
		out.Reason, out.Project, out.Folder, out.Story = WorkPull, pull.Project, pull.Folder, pull.Story
	case anyRoom:
		out.WaitingFor = WaitingForHeld
	case anyReady:
		out.WaitingFor = WaitingForRoom
	default:
		out.WaitingFor = WaitingForReady
	}
	return out, nil
}

// snapshot fingerprints every project's watched files, and the projects there are.
func (fw *folderWaits) snapshot() map[string]string {
	out := map[string]string{}
	for _, s := range fw.f.all() {
		out["project:"+projectRoot(s.repo)] = s.key
		for p, v := range s.snapshot() {
			out[p] = v
		}
	}
	return out
}

func (fw *folderWaits) timeout(seconds int, def time.Duration) time.Duration {
	t := time.Duration(seconds) * time.Second
	if t <= 0 {
		t = def
	}
	if t > fw.maxWait {
		t = fw.maxWait
	}
	return t
}

// hold waits until check says it is done, looking again whenever a watched
// file in any project changes, or a project comes or goes.
func (fw *folderWaits) hold(ctx context.Context, timeout time.Duration, check func(changed []string) (bool, error)) (timedOut bool, err error) {
	before := fw.snapshot()
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	tick := time.NewTicker(fw.poll)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-deadline.C:
			return true, nil
		case <-fw.closing:
			return true, nil
		case <-tick.C:
			after := fw.snapshot()
			changed := diff(before, after)
			if len(changed) == 0 {
				continue
			}
			before = after
			if done, err := check(changed); done || err != nil {
				return false, err
			}
		}
	}
}

func (fw *folderWaits) work(ctx context.Context, _ *mcp.CallToolRequest, in WorkIn) (*mcp.CallToolResult, FolderWorkOut, error) {
	f := fw.f
	f.workMu.Lock()
	since := f.workSince
	f.workMu.Unlock()
	if since.IsZero() {
		since = fw.now()
	}
	answer := func(out FolderWorkOut) FolderWorkOut {
		f.workMu.Lock()
		f.workSince = fw.now()
		f.workMu.Unlock()
		return out
	}
	out, err := fw.decide(ctx, since)
	if err != nil || out.Reason != "" {
		return nil, answer(out), err
	}
	timedOut, err := fw.hold(ctx, fw.timeout(in.TimeoutSeconds, fw.maxWait), func([]string) (bool, error) {
		next, err := fw.decide(ctx, since)
		out = next
		return err == nil && next.Reason != "", err
	})
	if err != nil {
		return nil, out, err
	}
	out.TimedOut = timedOut
	return nil, answer(out), nil
}

func (fw *folderWaits) catchUp(ctx context.Context) ([]Event, int, error) {
	events, omitted := []Event{}, 0
	for _, s := range fw.f.all() {
		some, n, err := s.catchUp(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", s.key, err)
		}
		for _, e := range some {
			e.Project = s.key
			events = append(events, e)
		}
		omitted += n
	}
	return events, omitted, nil
}

func (fw *folderWaits) events(ctx context.Context, _ *mcp.CallToolRequest, in WaitIn) (*mcp.CallToolResult, WaitOut, error) {
	events, omitted, err := fw.catchUp(ctx)
	if err != nil {
		return nil, WaitOut{}, err
	}
	if len(events) > 0 {
		return nil, WaitOut{Events: events, Omitted: omitted, Changed: []string{}}, nil
	}
	out := WaitOut{Events: []Event{}, Changed: []string{}}
	timedOut, err := fw.hold(ctx, fw.timeout(in.TimeoutSeconds, time.Minute), func(changed []string) (bool, error) {
		// the notices wake a waiting agent and are not a path to show it:
		// what they say arrives as events
		notices := map[string]bool{}
		for _, s := range fw.f.all() {
			for _, log := range noticeLogs(s.repo) {
				notices[log] = true
			}
		}
		for _, p := range changed {
			if strings.HasPrefix(p, "project:") || notices[p] {
				continue
			}
			out.Changed = append(out.Changed, fw.f.rel(p))
		}
		events, omitted, err := fw.catchUp(ctx)
		out.Events, out.Omitted = events, omitted
		return true, err
	})
	out.TimedOut = timedOut
	return nil, out, err
}
