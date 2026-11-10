package serve

import (
	"context"
	"fmt"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// StoryID matches a story ID: S-, then three or more digits.
var StoryID = regexp.MustCompile(`^S-\d{3,}$`)

// launcher starts agents for one project.
type launcher struct {
	dir    Dir
	entry  Entry
	config func(root string) AgentConfig
	record func(hostapi.Entry)
	now    func() time.Time
	log    func(msg string, args ...any)
	warn   func(msg string, args ...any)

	mu      sync.Mutex
	said    map[string]string // why each skipped story waits, as last logged
	waiting map[int]bool      // the PIDs of agents this launcher waits for
	// handOver leaves each agent it starts to the serving flai, which settles
	// the run when its process is gone, instead of waiting for it: set for a
	// start on the operator's word, from a process that exits once it has
	// started the agent (S-0115).
	handOver bool
	again    chan struct{} // an agent ended: look again
	// changed is told the run recorded as started, not started, or ended:
	// none of these changes a file of the project, so flai serve says so to
	// the dashboard (S-0154). Nil tells no one.
	changed func(run *AgentRun)
	// untold are the runs at the limit of automatic restarts whose thread
	// could not be opened, by story and start: warned of once, not at every
	// look (S-0294).
	untold map[string]bool
	// unasked are the asks about a hold that could not be written at the
	// last look, by held and holding story, with the error: warned of once
	// until the error changes, not at every look (ADR-0134).
	unasked map[string]string
	// holding are the stories whose agent's start the host's budget state
	// held since the look began, by ID, with the state that held it
	// (S-0534); heldSaid is what the last look logged of those it held
	// outside the ready column, so that a look every minute while a limit
	// holds does not repeat it.
	holding  map[string]Budget
	heldSaid string
	// git asks git for the trial merge that judges an overlap (ADR-0230),
	// the system's when nil.
	git execx.Runner
	// docker asks docker whether the container harness can start a story's
	// agent, and stops its container (S-0508).
	docker execx.Runner
	// engine is the operating system docker info named for the container
	// harness's engine, read once, when docker answered (containerEngine,
	// S-0619); engineMu guards it.
	engineMu sync.Mutex
	engine   string
	// imageCheck is Options.imageCheck: what checks a container image's
	// Claude Code version that has no record, beside the start (S-0618); nil
	// checks none.
	imageCheck func(root string, cfg AgentConfig, image, version, provider string)
}

// told records a run in the host's state and tells changed.
func (l *launcher) told(run *AgentRun) {
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(run) })
	if l.changed != nil {
		l.changed(run)
	}
}

// ended records a run as ended, as the operator's stop when they stopped it
// (S-0170): Stop marks the run before it signals the process, so the end
// the process makes is theirs, not a failure. A hold on the orchestrator
// stays on the run that ended (S-0228). It says whether the run was still
// recorded as going until now, so that its end is told once (S-0532): a run
// that Stop, or another look, recorded as ended first is recorded again, and
// told by whoever ended it.
func (l *launcher) ended(run *AgentRun) (first bool) {
	first = true
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) {
		cur := s.of(run)
		if !cur.same(run) {
			s.put(run)
			return
		}
		first = cur.Ended == ""
		if cur.Stopped != "" {
			run.Stopped = cur.Stopped
			run.stopped()
		}
		run.Held = cur.Held
		s.put(run)
	})
	if l.changed != nil {
		l.changed(run)
	}
	return first
}

func newLauncher(o Options, e Entry) *launcher {
	return &launcher{dir: o.Dir, entry: e, config: o.Agent, record: o.Host.Record, now: o.Now, waiting: map[int]bool{}, again: make(chan struct{}, 1), git: gitOf(o), docker: dockerOf(o), imageCheck: o.imageCheck,
		log: func(msg string, args ...any) {
			o.Logger.Info(msg, append([]any{"component", "serve", "project", e.Key}, args...)...)
		},
		warn: func(msg string, args ...any) {
			o.Logger.Warn(msg, append([]any{"component", "serve", "project", e.Key}, args...)...)
		}}
}

// readyStory is a story in ready, as the launcher needs it.
type readyStory struct {
	ID      string
	Agent   *manifest.Agent
	Entered time.Time
	// Restart says how its last agent ended, when the operator has it
	// started again (S-0116); empty when it entered ready.
	Restart string
	// Commit is the story's worktree, when the operator has its agent
	// started to commit what the worktree holds (S-0140).
	Commit string
	// Asked is the run that ended asking, once its question is answered: the
	// agent is started again in its session (S-0182). Answer is what answered
	// it: the thread, or the conversation a reply or a new message to the
	// story came on (S-0335).
	Asked  *AgentRun
	Answer string
	// SentBack is what flai sent a story in progress back from review for,
	// when flai serve starts its agent for that (S-0394, ADR-0167), and
	// Resume the last run, whose session the agent goes on in, when it had
	// one; nil for a new session.
	SentBack *harness.SentBack
	Resume   *AgentRun
	// Started is set when the operator has its agent started now (S-0115),
	// and Past says what that start went past: a hold's reason, a full
	// in-progress limit, a full review (S-0243), or nothing (S-0182).
	Started bool
	Past    []string
	// Begun says where a story in progress that this host has had no agent
	// for was begun, when the operator has one started here (S-0177).
	Begun *harness.Begun
	// AutoRestarts is how many automatic restarts in a row this start makes,
	// counting it, when flai serve restarts the story's agent on its own
	// (S-0294, ADR-0108), or, with AfterLimit, how many its last run had; 0
	// for any other start.
	AutoRestarts int
	// AfterLimit is set when flai serve starts the agent again because its
	// last run ended on the harness's usage limit, once its provider's
	// budget state no longer holds it (S-0534): the start spends no
	// automatic restart, and the run keeps its last run's count.
	AfterLimit bool
	item       *workitem.Item
}

// formulaOf is the formula story is an instance of (S-0422): its item's
// formula when the board read it, else the story's as root holds it, and
// empty when it is no instance or cannot be read.
func formulaOf(root string, story readyStory) string {
	if story.item != nil {
		return story.item.Formula
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return ""
	}
	it, err := repo.Get(story.ID)
	if err != nil {
		return ""
	}
	return it.Formula
}

// withTrial has repo hold a ready story on overlap only when a trial merge,
// asked of git, the system's when nil, conflicts (ADR-0230); outside a git
// checkout NewTrialJudge gives none, and an overlap holds as before.
func withTrial(repo *workitem.Repo, git execx.Runner) {
	if git == nil {
		git = execx.System{}
	}
	repo.Trial = storygit.NewTrialJudge(git, repo, issues.FilesIn(repo, repo.MainRoot, git, time.Now))
}

// readyStories are the ready stories in pull order, the claims of the open
// stories they are held by, how many more stories the in-progress limit
// leaves room for, -1 when there is none, why review being full holds every
// one back, "" when it does not (S-0243), and the items read, by ID. Each
// story's agent is the one it starts with: the project's default under it
// when it names no harness (ADR-0251).
func readyStories(repo *workitem.Repo) (ready []readyStory, holds *workitem.Holds, free int, review string, byID map[string]*workitem.Item, err error) {
	items, err := repo.List(false)
	if err != nil {
		return nil, nil, 0, "", nil, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, nil, 0, "", nil, err
	}
	byID = map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	// Pull order and the counts never look at done or archived items, so
	// which ones stay visible there (S-0087) makes no difference here.
	holds = repo.Holds(items)
	view := workitem.NewBoardView(items, board, time.Now(), false, nil, holds)
	for _, c := range view.ReadyInPullOrder() {
		if it := byID[c.ID]; it != nil && c.Type == workitem.Story && !c.Blocked {
			ready = append(ready, readyStory{ID: c.ID, Agent: repo.AgentOf(it), Entered: it.EnteredAt(), item: it})
		}
	}
	free = -1
	if limit, ok := view.WIPLimits[workitem.InProgress]; ok && limit > 0 {
		free = max(0, limit-view.Counts[workitem.InProgress])
	}
	return ready, holds, free, view.ReviewHold(), byID, nil
}

// readyChecked is readyStories read under the checked-write lock, so that no
// claim is judged by an item a checked creation wrote and may yet remove
// (I-0301). The item folders are looked at again under it, since a snapshot
// kept before a removal stays fresh until the watcher marks it, which may not
// have happened yet; the look reuses each unchanged file's parse. A lock it
// cannot take is warned of and ends the look, as a read that fails does.
func (l *launcher) readyChecked(repo *workitem.Repo) (ready []readyStory, holds *workitem.Holds, free int, review string, byID map[string]*workitem.Item, err error) {
	unlock, err := repo.LockCheckedWrites()
	if err != nil {
		l.warn("ready stories not read", "err", err.Error())
		return nil, nil, 0, "", nil, err
	}
	defer unlock()
	workitem.ItemsChanged(repo.Root)
	return readyStories(repo)
}

// agentStarted is how a hold names a story in ready whose agent has been
// started: it claims its paths from then, not once its agent pulls it.
const agentStarted = "its agent started"

// judgeRun records on run, a story's agent's run that has ended with exit,
// how it left its story: worked, asked, or failed, with why, and what an
// asked run waits on. It ended asking when a question of its own on the
// story is open and unanswered, when a conversation of its story awaits
// another story's agent's reply (S-0335), or when a message to its story came
// on a conversation after the run started and has had no reply: a reply that
// came before the end was judged, which would otherwise leave the agent
// waiting on what has come, so that it is started again at once to read it.
// So, too, when the question it wrote on a thread of the story after the run
// started was answered before the end was judged (I-0095).
func judgeRun(root string, run *AgentRun, exit *int) {
	run.Outcome, run.Why, run.Thread, run.Conversations, run.WaitsOn = OutcomeFailed, "", "", nil, nil
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		run.Why = fmt.Sprintf("ended (%s); the project could not be read: %v", code, err)
		return
	}
	it, err := repo.Get(run.Story)
	if err != nil {
		run.Why = fmt.Sprintf("ended (%s); %s could not be read: %v", code, run.Story, err)
		return
	}
	switch it.Status {
	case workitem.Review, workitem.Done:
		run.Outcome = OutcomeWorked
		return
	}
	th := asking(repo, run.Story, run.Agent)
	var since *threads.Thread
	if th == nil {
		since = askedAnswered(repo, run)
	}
	waits, unread := awaitingOther(repo, run.Story), unanswered(repo, run.Story, run.Started)
	if th != nil || since != nil || len(waits) > 0 || len(unread) > 0 {
		run.Outcome, run.Why = OutcomeAsked, waitWhy(th, since, waits, unread, run.Story)
		if th != nil {
			run.Thread = th.ID
		}
		if since != nil {
			run.Thread = since.ID
		}
		run.Conversations, run.WaitsOn = conversationIDs(slices.Concat(waits, unread)), others(waits, run.Story)
		return
	}
	run.Why = fmt.Sprintf("ended (%s) with %s in %s", code, run.Story, it.Status)
	for _, b := range it.Blocked {
		if b.Until == "" {
			run.Why += ", blocked: " + b.Reason
		}
	}
}

// endRun judges run, a story's agent's run that has ended with exit
// (judgeRun), records the end (ended), and then blocks its story when it
// ended asking on a thread (blockAsked): after the end is recorded, so that
// releaseAsked, which closes the interval of a story whose last run no
// longer waits on its thread, never takes the run still live for one that
// does. The other way round, a look that starts the agent again before the
// block is written finds the block at its next look, and releases it. The
// end is appended to the event log once, by the end that recorded the run
// as ended first (S-0532): a run Stop recorded as stopped before its
// process ended was told by Stop, and one judged again tells nothing more.
// A run whose log ends on the harness's usage limit is marked so
// (limitEnded), and one judged failed says the limit as its why, as the
// curator's does; that end feeds its provider's budget state, once, as every
// end does (budgetEnd), and restartEnded does not restart it (S-0534). On a
// harness that makes its own session IDs, the session its log names is
// recorded as the run's, which a start that resumes it goes on in; a log that
// names none leaves the one flai made, which the adapter refuses to resume.
func (l *launcher) endRun(run *AgentRun, exit *int) {
	judgeRun(l.entry.Root, run, exit)
	if r, ok := harness.Adapters[run.Harness].(harness.SessionReader); ok && run.Log != "" {
		log, err := os.ReadFile(run.Log)
		if err != nil {
			l.warn("run's log not read for its session", "log", run.Log, "err", err)
		}
		if session, ok := r.Session(log); ok {
			run.Session = session
		}
	}
	lim, _ := runLimit(run, l.warn)
	run.limitEnded(lim)
	if lim != nil && run.Outcome == OutcomeFailed {
		run.Why = limitWhy(lim)
	}
	first := l.ended(run)
	if first {
		appendAgentEnd(l.entry.Root, run, l.warn)
	}
	l.budgetEnd(run, lim, first)
	l.blockAsked(run)
}

// awaitingOperator is the reason of the blocked interval flai serve opens on
// a story whose agent ended asking on thread (ADR-0205); awaitingReason
// reads the thread back out of it.
func awaitingOperator(thread string) string {
	return "awaiting the operator on " + thread
}

var awaitingReason = regexp.MustCompile(`^awaiting the operator on (TH-\d{4,})$`)

// askedOn is the thread that the open blocked interval of it names, when it
// is one flai serve opened as its agent ended asking (blockAsked), or "": an
// interval an agent or the operator wrote with flai block names no thread so.
func askedOn(it *workitem.Item) string {
	for _, b := range it.Blocked {
		if m := awaitingReason.FindStringSubmatch(b.Reason); b.Until == "" && m != nil {
			return m[1]
		}
	}
	return ""
}

// blockAsked blocks the story of run, which ended asking on a thread, with
// the reason that names the thread, so that the board shows it waiting on
// the operator and flai stats counts the wait as blocked time (ADR-0205,
// I-0151). It leaves alone a run that ended any other way or asked on
// conversations alone, a story not in progress, one blocked already, and a
// question answered before the block, which resume starts again at once.
// A failure is warned of: the run is started again on its answer as before.
func (l *launcher) blockAsked(run *AgentRun) {
	if run.Outcome != OutcomeAsked || run.Thread == "" {
		return
	}
	err := func() error {
		repo, err := workitem.Open(l.entry.Root)
		if err != nil {
			return err
		}
		if threadAnswered(repo, run) {
			return nil
		}
		it, err := repo.Get(run.Story)
		if err != nil {
			return err
		}
		if it.Status != workitem.InProgress || it.IsBlocked() {
			return nil
		}
		if err := workitem.BlockItem(it, awaitingOperator(run.Thread), l.now()); err != nil {
			return err
		}
		if err := repo.Save(it); err != nil {
			return err
		}
		l.log("story blocked on its agent's question", "story", run.Story, "thread", run.Thread)
		return nil
	}()
	if err != nil {
		l.warn("story not blocked on its agent's question", "story", run.Story, "thread", run.Thread, "err", err.Error())
	}
}

// unblockAsked closes the interval blockAsked opened on story for its
// agent's question on thread, and reports whether it did: not when the
// story's open interval is not that one. It reads the story again just
// before writing, as an agent may have changed it.
func (l *launcher) unblockAsked(repo *workitem.Repo, story, thread string) (bool, error) {
	it, err := repo.Get(story)
	if err != nil {
		return false, err
	}
	if thread == "" || askedOn(it) != thread {
		return false, nil
	}
	if err := workitem.UnblockItem(it, l.now()); err != nil {
		return false, err
	}
	if err := repo.Save(it); err != nil {
		return false, fmt.Errorf("close %s's blocked interval for %s: %w", story, thread, err)
	}
	return true, nil
}

// releaseAsked closes the interval blockAsked opened on each story, of
// items, the ones the look read, whose agent no longer waits on the thread
// it names: the story left in progress, sent back to ready or backlog,
// moved to review, or cancelled, or its last run is not the one that asked,
// its agent having been started again otherwise than on the answer, by the
// operator. A blocked story in ready is never started, so one sent back with
// its question open would otherwise wait there for good. It reports whether
// it closed any.
func (l *launcher) releaseAsked(repo *workitem.Repo, items map[string]*workitem.Item) bool {
	closed := false
	for id, run := range l.dir.AgentStates()[l.entry.Root].Stories {
		it := items[id]
		if it == nil {
			continue
		}
		thread := askedOn(it)
		if thread == "" || (it.Status == workitem.InProgress && asked(run) && run.Thread == thread) {
			continue
		}
		ok, err := l.unblockAsked(repo, id, thread)
		if err != nil {
			l.warn("story not unblocked from its agent's question", "story", id, "thread", thread, "err", err.Error())
			continue
		}
		if ok {
			closed = true
			l.log("story unblocked: its agent no longer waits on its question", "story", id, "thread", thread, "status", it.Status)
		}
	}
	return closed
}

// awaitingOther are story's open conversations that await the other story's
// agent's reply; none when they cannot be read, as no question is asked when
// threads cannot be.
func awaitingOther(repo *workitem.Repo, story string) []*messages.Conversation {
	out, err := messages.AwaitingOther(repo, story)
	if err != nil {
		return nil
	}
	return out
}

// conversations is one read of a project's conversations, from which one
// Activity call judges every live run and every held story, rather than
// reading them all again for each (S-0603). Its methods on a nil one read
// what they need from the files, as before.
type conversations struct {
	all    []*messages.Conversation // each that read and parsed, in ID order
	err    error                    // the first that did not, on which messages.List fails
	shares []workitem.Share         // the shares of those stored open, as Repo.Conversations reads them
	asks   []workitem.Ask           // and their asks about a hold
}

// conversationsRead, when set, is called once Activity has read the
// conversations: the tests' proof that nothing reads them again.
var conversationsRead atomic.Pointer[func(repo *workitem.Repo)]

// readConversations reads repo's conversations once, for the lists
// messages.List gives and for the shares and asks Repo.Conversations gives,
// in file order as it reads them.
func readConversations(repo *workitem.Repo) *conversations {
	all, err := messages.ReadAll(repo)
	byFile := slices.Clone(all)
	slices.SortFunc(byFile, func(a, b *messages.Conversation) int { return strings.Compare(a.Path, b.Path) })
	docs := make([]workitem.ConversationDoc, 0, len(byFile))
	for _, c := range byFile {
		docs = append(docs, workitem.ConversationDoc{ID: c.ID, From: c.From, To: c.To, Status: c.Status, Shares: c.Shares, Body: c.Body})
	}
	out := &conversations{all: all, err: err}
	out.shares, out.asks = workitem.SharesAndAsks(docs)
	if read := conversationsRead.Load(); read != nil {
		(*read)(repo)
	}
	return out
}

// awaitingOther is awaitingOther from what c read.
func (c *conversations) awaitingOther(repo *workitem.Repo, story string) []*messages.Conversation {
	if c == nil {
		return awaitingOther(repo, story)
	}
	if c.err != nil {
		return nil
	}
	return messages.AwaitingOtherIn(repo, c.all, story)
}

// get is the conversation id, in any padding, as messages.Get finds it.
func (c *conversations) get(repo *workitem.Repo, id string) (*messages.Conversation, error) {
	if c == nil {
		return messages.Get(repo, id)
	}
	canon := messages.CanonicalID(id)
	for _, conv := range c.all {
		if messages.CanonicalID(conv.ID) == canon {
			return conv, nil
		}
	}
	return nil, fmt.Errorf("conversation %s not found in %s: flai message list --all lists them", canon, messages.Dir(repo))
}

// toSince is messages.ToSince from what c read.
func (c *conversations) toSince(repo *workitem.Repo, story string, t time.Time) ([]*messages.Conversation, error) {
	if c == nil {
		return messages.ToSince(repo, story, t)
	}
	if c.err != nil {
		return nil, c.err
	}
	return messages.ToSinceIn(c.all, story, t), nil
}

// unanswered are story's open conversations that await its reply to a
// message the other story's agent wrote after started, a run's start; none
// when that is not known.
func unanswered(repo *workitem.Repo, story, started string) []*messages.Conversation {
	at, err := time.Parse(time.RFC3339, started)
	if err != nil {
		return nil
	}
	since, err := messages.ToSince(repo, story, at)
	if err != nil {
		return nil
	}
	return slices.DeleteFunc(since, func(c *messages.Conversation) bool {
		closed, _ := c.Closed(repo)
		return closed || !owesReply(c, story)
	})
}

// owesReply says whether c awaits story's reply: the other story wrote last.
func owesReply(c *messages.Conversation, story string) bool {
	return workitem.CanonicalID(c.Awaiting()) == workitem.CanonicalID(story)
}

// waitWhy says what a story's agent waits for: an answer to th, its question
// to the operator, when there is one; the reply of the other story's agent on
// each of waits; and, for a run that has ended, the answer to its question on
// since that came before the end was judged and a message on each of unread
// that it left unanswered.
func waitWhy(th, since *threads.Thread, waits, unread []*messages.Conversation, story string) string {
	var on, why []string
	if th != nil {
		on = append(on, "an answer to "+th.ID+": "+th.Title)
	}
	for _, c := range waits {
		on = append(on, c.Other(story)+"'s agent to reply on "+c.ID+": "+c.Title)
	}
	if len(on) > 0 {
		why = append(why, "waiting for "+strings.Join(on, "; and for "))
	}
	if since != nil {
		why = append(why, "ended asking on "+since.ID+", answered before the end was judged: "+since.Title)
	}
	for _, c := range unread {
		why = append(why, "ended with a message from "+c.Other(story)+"'s agent on "+c.ID+" unanswered: "+c.Title)
	}
	return strings.Join(why, "; ")
}

func conversationIDs(cs []*messages.Conversation) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.ID)
	}
	return out
}

// others are the stories story converses with in cs, each once, in order.
func others(cs []*messages.Conversation, story string) []string {
	var out []string
	for _, c := range cs {
		if o := c.Other(story); !slices.ContainsFunc(out, func(s string) bool { return workitem.CanonicalID(s) == workitem.CanonicalID(o) }) {
			out = append(out, o)
		}
	}
	return out
}

// look is called when the project's work items may have changed, when an
// agent this launcher started has ended, and now and then.
func (l *launcher) look(ctx context.Context, _ bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.config == nil {
		return
	}
	repo, err := workitem.Open(l.entry.Root)
	if err != nil {
		return
	}
	withTrial(repo, l.git)
	stories, holds, free, review, items, err := l.readyChecked(repo)
	if err != nil {
		return
	}
	if l.releaseAsked(repo, items) {
		// a story it unblocked may be ready to start, or no longer held
		if stories, holds, free, review, items, err = l.readyChecked(repo); err != nil {
			return
		}
	}
	// before the return on an empty ready column, which is the column at its
	// most drained (ADR-0211 §1), counting only the stories not held on after
	// (ADR-0279)
	l.noteReadyEmpty(repo, readyNotHeldOnAfter(stories, holds), free)
	// The first look after flai serve starts is like any other: what keeps a
	// restart from starting a story twice is its run in serve/agents.json,
	// not what happened to be ready when flai serve began (S-0112).
	l.settleOrphans()
	l.measureDone(l.dir.AgentStates()[l.entry.Root])
	l.resolveAcceptSentBack(repo)
	l.holding = nil
	defer l.sayHeld()
	cfg := l.config(l.entry.Root)
	if cfg.Enabled {
		l.resume(ctx, cfg, repo)
		// before restartEnded: a run from before the send back that ended
		// failed is started for the send back, spending no automatic restart
		l.startSentBack(ctx, cfg, repo)
		l.restartEnded(ctx, cfg, repo)
	}
	var sk skips
	defer func() { l.say(sk) }()
	if len(stories) == 0 {
		return
	}
	st := l.dir.AgentStates()[l.entry.Root]
	// Deferred, so that it runs on every way out of the look, the agent
	// action off and the limit full among them, with the claims of the agents
	// this look started opened in holds.
	defer l.askHolds(repo, stories, holds, items, st)
	if !cfg.Enabled {
		sk.add("the agent host action is off for this project", stories...)
		return
	}
	var todo, away, noRoom, reviewFull []readyStory
	byLimit := map[string][]readyStory{} // held by the budget state, by why
	reserved := 0
	for _, s := range stories {
		// what it starts with, flai's built-in default when nothing names a
		// harness and no command is set (ADR-0261)
		s.Agent = cfg.startAgent(s.Agent)
		run := st.Stories[s.ID]
		answer := ""
		if asked(run) {
			answer = answered(repo, run)
		}
		switch {
		case run.live():
			reserved++ // started, and not in progress yet: it has its agent
			holds.Open(s.item, agentStarted)
		case harness.Elsewhere(s.Agent):
			// its agent runs elsewhere, which no command on the host changes,
			// and is no failure (ADR-0242)
			away = append(away, s)
		case asked(run) && answer == "":
			sk.add(s.ID+"'s agent is "+run.awaits(), s)
		case asked(run):
			// answered while its story is in ready: started again in its
			// session, but, unlike resume, only when nothing holds it and
			// the limit has room (S-0182)
			s.Asked, s.Answer = run, answer
			todo = append(todo, s)
		case afterLimit(run):
			// its last agent ended on the usage limit before it pulled the
			// story: started again once the limit no longer holds its
			// provider, which the start's gate says, spending no restart
			s.AfterLimit, s.AutoRestarts, s.Restart = true, run.AutoRestarts, restartWhy(run)
			todo = append(todo, s)
		case run != nil && run.Queued == "" && !startedBefore(run, s.Entered) && !agentChanged(run, s.Agent):
			sk.add(tried(s.ID, run), s)
		default:
			if run != nil && run.Queued != "" {
				s.Restart = restartWhy(run) // the operator's retry, waiting for room
			}
			todo = append(todo, s)
		}
	}
	for _, s := range away {
		sk.add(s.ID+"'s agent runs elsewhere (harness none): flai serve starts nothing for it", s)
	}
	for _, s := range todo {
		switch h := holds.Of(s.item); {
		case h != nil:
			sk.add(s.ID+" "+h.Reason, s)
		case free >= 0 && reserved >= free:
			noRoom = append(noRoom, s)
		case review != "":
			// review full stops the pull, not the work: what is started
			// or in progress goes on (S-0243, I-0007)
			reviewFull = append(reviewFull, s)
		case l.start(ctx, cfg, s):
			reserved++
			holds.Open(s.item, agentStarted)
		default:
			if b, ok := l.holding[s.ID]; ok {
				// said with the ready stories that wait, not apart (sayHeld)
				delete(l.holding, s.ID)
				why := budgetHeldWhy(b, l.now())
				byLimit[why] = append(byLimit[why], s)
			}
		}
	}
	for _, why := range slices.Sorted(maps.Keys(byLimit)) {
		sk.add(why+"; it holds "+strings.Join(ids(byLimit[why]), ", "), byLimit[why]...)
	}
	if len(noRoom) > 0 {
		sk.add("the in-progress limit leaves no room for "+strings.Join(ids(noRoom), ", "), noRoom...)
	}
	if len(reviewFull) > 0 {
		sk.add(review+"; it holds "+strings.Join(ids(reviewFull), ", "), reviewFull...)
	}
}

// askHolds asks the agent of each story in progress that holds a ready story
// on overlap alone about the hold, in the two stories' conversation, with the
// paths it is held on (ADR-0134). It asks whether or not the agent action is
// on, the in-progress limit has room, or review is full: a held story is never
// started, so the ask matters most when the board is busy. messages.AskHold
// asks once per pair while the held story stays in ready, so a look that finds
// the same hold writes nothing. A story held by after: or by a story with no
// touches is not asked about (OverlapsBy), nor one whose agent runs, and a
// holder in ready with its agent started is not asked. An ask that cannot be
// written is warned of, once until its error changes, and the look goes on.
func (l *launcher) askHolds(repo *workitem.Repo, stories []readyStory, holds *workitem.Holds, items map[string]*workitem.Item, st AgentState) {
	now := l.now()
	failed := map[string]string{}
	for _, s := range stories {
		if st.Stories[s.ID].live() {
			continue
		}
		for _, o := range holds.OverlapsBy(s.item) {
			if h := items[o.By]; h == nil || h.Status != workitem.InProgress {
				continue
			}
			c, wrote, err := messages.AskHold(repo, messages.AskOptions{Held: s.ID, Holder: o.By, Paths: o.Paths, Now: now})
			key := s.ID + " " + o.By
			switch {
			case err != nil:
				failed[key] = err.Error()
				if l.unasked[key] != err.Error() {
					l.warn("hold not asked about", "held", s.ID, "holder", o.By, "err", err.Error())
				}
			case wrote:
				l.log("hold asked about", "held", s.ID, "holder", o.By, "conversation", c.ID)
			}
		}
	}
	l.unasked = failed
}

// skips are the ready stories a look did not start, and why, in the order
// said.
type skips struct {
	whys  []string
	story map[string]string
}

func (k *skips) add(why string, stories ...readyStory) {
	if k.story == nil {
		k.story = map[string]string{}
	}
	k.whys = append(k.whys, why)
	for _, s := range stories {
		k.story[s.ID] = why
	}
}

// say puts why each skipped story waits in the state's waiting, and logs a
// story's reason when it is not the kind last logged, so that a look every
// minute does not repeat it.
func (l *launcher) say(k skips) {
	why := strings.Join(k.whys, "; ")
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = why })
	for _, id := range slices.Sorted(maps.Keys(k.story)) {
		if _, said := l.said[id]; !said || kind(l.said[id]) != kind(k.story[id]) {
			l.log("agent not started", "story", id, "why", k.story[id])
		}
	}
	l.said = k.story
}

// kind is a reason without the detail in its parentheses: when a run
// started, or how it ended, is not a new reason.
func kind(why string) string {
	if i := strings.Index(why, " ("); i >= 0 {
		return why[:i]
	}
	return why
}

// tried says of a story that has had its agent since it entered ready what
// became of it, and what starts another.
func tried(id string, run *AgentRun) string {
	what := "started " + run.Started
	switch {
	case run.Error != "":
		what += ", could not start: " + run.Error
	case run.Why != "":
		what += ", " + run.Outcome + ": " + run.Why
	case run.Outcome != "":
		what += ", " + run.Outcome
	}
	return id + " has had its agent since it entered ready (" + what + "); it gets another when it enters ready again, its agent is changed, or it is restarted (flai serve agent restart)"
}

func ids(stories []readyStory) []string {
	out := make([]string, len(stories))
	for i, s := range stories {
		out[i] = s.ID
	}
	return out
}

// startedBefore says whether run began before the story entered ready.
func startedBefore(run *AgentRun, entered time.Time) bool {
	started, err := time.Parse(time.RFC3339, run.Started)
	return err == nil && started.Before(entered.Truncate(time.Second))
}

// settleOrphans ends the runs whose process is gone while no launcher waits
// for it: flai serve was restarted while they ran, or a command started them
// on the operator's word and handed them over (S-0115, S-0116). A process
// that has the run's PID and is not the one started for it, after a reboot,
// is gone too (S-0170). A planner's run is settled the same way, and its
// activity logged (S-0208), and so are the orchestrator's (S-0218), the
// analyzer's (S-0223), the curator's (S-0431), and each defined role's
// (S-0417).
func (l *launcher) settleOrphans() {
	st := l.dir.AgentStates()[l.entry.Root]
	for _, run := range st.Stories {
		if !run.live() || l.waiting[run.PID] || run.running() {
			continue
		}
		ended := *run
		ended.Ended = l.now().UTC().Format(time.RFC3339)
		l.endRun(&ended, nil)
		l.log("agent ended, seen at a look", "story", run.Story, "pid", run.PID, "outcome", ended.Outcome)
		l.measure(run.Story)
	}
	for _, run := range st.Plans {
		if !run.live() || l.waiting[run.PID] || run.running() {
			continue
		}
		l.planEnded(run, nil)
	}
	if run := st.Orchestrator; run.live() && !l.waiting[run.PID] && !run.running() {
		l.orchestrateEnded(run, nil)
	}
	if run := st.Analyzer; run.live() && !l.waiting[run.PID] && !run.running() {
		l.analyzeEnded(run, nil)
	}
	if run := st.Curator; run.live() && !l.waiting[run.PID] && !run.running() {
		l.curateEnded(run, nil)
	}
	for _, runs := range st.Roles {
		for _, run := range runs {
			if run.live() && !l.waiting[run.PID] && !run.running() {
				l.roleEnded(run, nil)
			}
		}
	}
}
