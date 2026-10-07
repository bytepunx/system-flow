// Package guard refuses what a sub-agent may not do (ADR-0060): Claude Code
// runs flai guard before each tool call in a project made from the template,
// with the call on standard input, and says which sub-agent makes it. A
// sub-agent reads; it does not change a work item, a thread, or the
// repository's history, because the story's agent alone acts for the story
// (ADR-0059). The story's agent's own calls carry no agent ID and are never
// refused. The guard is not a shell: it finds the programs a command line
// runs well enough to stop a sub-agent that follows its instructions
// carelessly, not one that sets out to hide a command (a backslash inside a
// name, a command held in a variable).
//
// A sub-agent does not write a file in a .claude/ folder, with Edit,
// MultiEdit, Write, or NotebookEdit, unless the operator enabled the
// project's auto-approve host action (S-0299): Claude Code asks for each such
// write, and flai's permission_prompt asks the operator on a thread, which
// would hold the call, and the sub-agent's layer with it, until they answer.
// The sub-agent puts the file's content in its final message instead, for
// the story's agent to write.
//
// A planner session, one flai serve starts with the role plan, is held to
// planning (strategic-agents.md): its own calls, which carry no agent ID,
// may read, write work items and threads through flai, and move an item to
// backlog, but not edit files with Edit, Write, or NotebookEdit, move an
// item further, accept, publish, work a story, or change the repository's
// history. A story it creates is a draft for the operator to finalize
// (S-0209). Its sub-agents are held as every sub-agent is.
//
// An orchestrator session, one flai serve starts with the role orchestrate,
// is held to the permissions the operator gives it in the manifest's
// orchestration.permissions (S-0218): its own calls may read, open threads,
// record issues, and log its activities, and do each of the rest only while
// the permission that allows it is on (S-0219 maps planning an epic,
// finalizing a draft, promoting a story to ready, and ordering the ready
// column by a policy). With plan_backlog_stories it asks for the planner on
// a story, and gives a story cost of delay inputs with an edit that changes
// nothing else, never a value or a removal (S-0328, ADR-0119); flai judges
// which story. With accept_reviews it accepts a story in review, with
// flai accept or flai move to done, only as itself, --by orchestrator
// (S-0221, ADR-0093). A refusal names that permission; a call no permission
// allows, such as an edit of a file, a commit, or a story placed by hand in
// the pull order, it never makes. With publish it publishes through
// release_publish alone, and never with flai release, flai push, git push, or
// git tag, whatever its permissions (S-0222); no other session calls
// release_publish. Its sub-agents are held as every sub-agent is.
//
// An analyzer session, one flai serve starts with the role analyze, is held
// to its report (S-0223): its own calls may read, flai stats among them, log
// its activities, and open and reply to threads, and may edit, with Edit,
// Write, or NotebookEdit, only a file under the manifest's design folder's
// analysis/, its report and the folder's README.md that lists it. It files
// and bumps an issue for each actionable finding (S-0224), with flai issue
// new and bump or the MCP tools issue_new and issue_bump, which write the
// issue through flai, never through a file edit. It authors no stories:
// item_new, item_edit, item_move, issue_story, and flai issue story, story,
// epic, and task are refused it, as is every other write through flai, flai
// issue close among them, and git's. Its sub-agents, the explorer it hands
// search to among them, are held as every sub-agent is.
//
// Its calls on threads are held by who opened the thread and by
// answer_threads (S-0220): on a thread it opened it follows up and resolves,
// but never recommends or answers its own question; on one a story's
// planner opened, its opener planner-S-nnnn, plan_backlog_stories lets it
// answer, with or without a source, recommend, and resolve (S-0328,
// ADR-0119); on another's it replies only while answer_threads is on, as a
// recommendation for the operator to confirm, or, while it is autonomous, as
// an answer that cites its source, and it never resolves one. Confirming a
// recommendation is the operator's alone (ADR-0090).
//
// A story's agent, in a session flai serve starts with FLAI_STORY and no
// role, is refused its own wait_for_events while a sub-agent of its session
// runs and no thread on its story awaits the designer (S-0285): the call
// reports work items and threads, not sub-agents, so it would run to its
// timeout. Record keeps each session's running sub-agents, from Claude
// Code's SubagentStart and SubagentStop hooks.
//
// No session flai serve starts, and no sub-agent, changes the manifest's
// shared paths, with shared_paths_edit or flai shared add or remove
// (S-0295, ADR-0096): they decide which overlaps hold a story. Only the
// operator's own session does; listing and checking them are reads.
package guard

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/protected"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Event is the part of a Claude Code hook's input the guard reads: a
// PreToolUse's tool call, or a SubagentStart's or SubagentStop's sub-agent,
// as HookEventName says, in the session SessionID.
type Event struct {
	HookEventName string `json:"hook_event_name"`
	SessionID     string `json:"session_id"`
	ToolName      string `json:"tool_name"`
	ToolInput     Input  `json:"tool_input"`
	// AgentID is set only when a sub-agent makes the call, or on the
	// sub-agent a SubagentStart or SubagentStop reports; AgentType names the
	// sub-agent's definition. A session started with --agent may carry an
	// agent type of its own, so the ID is what marks a sub-agent.
	AgentID   string `json:"agent_id"`
	AgentType string `json:"agent_type"`
}

// Input is the part of a tool call's input the guard reads: Bash's command;
// item_move's item and the state it moves the item to, plan's item, and the
// thread of thread_reply and thread_resolve; item_new's type and whether it
// makes a story a draft, and item_edit's draft; whether thread_reply is a
// recommendation and the source it cites; and the file Edit or Write
// changes, or the notebook NotebookEdit does. Fields are the names of every
// field the input gives, whatever the guard reads of it, so that an
// item_edit that finalizes a draft is told from one that changes more, and a
// draft false given from one left out. CostOfDelay is item_edit's
// cost_of_delay, each key it gives with its value as given, so that an edit
// that gives a story cost of delay inputs is told from one that sets a value
// or removes one (S-0328); it is nil when cost_of_delay is not an object.
type Input struct {
	Command        string                     `json:"command"`
	ID             string                     `json:"id"`
	To             string                     `json:"to"`
	Type           string                     `json:"type"`
	Draft          bool                       `json:"draft"`
	Recommendation bool                       `json:"recommendation"`
	Source         string                     `json:"source"`
	FilePath       string                     `json:"file_path"`
	NotebookPath   string                     `json:"notebook_path"`
	Fields         []string                   `json:"-"`
	CostOfDelay    map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON decodes a tool call's input, the names of its fields, and
// the keys of its cost_of_delay. A cost_of_delay that is not an object is
// left nil rather than failing the input, which flai guard would let through
// unread.
func (in *Input) UnmarshalJSON(data []byte) error {
	type input Input
	if err := json.Unmarshal(data, (*input)(in)); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	in.Fields = slices.Sorted(maps.Keys(fields))
	if raw, ok := fields["cost_of_delay"]; ok && json.Unmarshal(raw, &in.CostOfDelay) != nil {
		in.CostOfDelay = nil
	}
	return nil
}

// MCPPrefix starts the name of each of flai's MCP tools as Claude Code
// knows them, from a server named flai.
const MCPPrefix = "mcp__flai__"

// MCPReads are flai's MCP tools a sub-agent may call: they read and use no
// agent's identity.
// verify runs a story's checks and stores its last result in flai's cache,
// which is neither a work item, a thread, nor history (S-0270).
var MCPReads = []string{"board", "doc_get", "doc_search", "item_get", "order_by_policy", "prime", "promote_candidates", "release_evaluate", "shared_paths", "test", "thread_get", "verify", "who_touches"}

// cliReads are the flai commands a sub-agent may run, each with the
// subcommands it may run; nil allows the command whatever follows it, and ""
// allows it with no subcommand.
var cliReads = map[string][]string{
	"board":    {""},
	"check":    nil,
	"cod":      nil,
	"criteria": {"list"},
	"doc":      {"search", "show"},
	"forecast": nil,
	"help":     nil,
	"issue":    {"list"},
	"prime":    nil,
	"shared":   {"list", "check"},
	"show":     nil,
	"stats":    nil,
	"stream":   {"diff"},
	"test":     nil,
	"thread":   {"list", "show"},
	"touches":  {"suggest"},
	"version":  nil,
}

// flagReads are the flai commands a sub-agent may run only in the form that
// reads (S-0217, S-0219): with one of the flags that make them read, without
// any flag that makes them write, and in the form form allows when it is
// set. flai order --by computes an order and --apply writes it; flai order
// that names a story places it; flai release without --evaluate releases;
// flai plan without --candidates starts the planner, which with it takes
// no item; and flai verify --record-issues writes issues, which without it
// it does not (S-0270). A command with no with flags reads in any form
// without its without flags.
var flagReads = map[string]struct {
	with    []string
	without []string
	form    func(words []string) bool
}{
	"order":   {with: []string{"--by"}, without: []string{"--apply"}, form: ordersByPolicy},
	"plan":    {with: []string{"--candidates"}},
	"promote": {with: []string{"--candidates", "--drafts"}},
	"release": {with: []string{"--evaluate"}},
	"verify":  {without: []string{"--record-issues"}},
}

// orderValues are the flags of flai order, flai's own among them, that take
// a value; placements are those that place one story by hand.
var (
	orderValues = map[string]bool{"--after": true, "--before": true, "--by": true, "--config": true, "--keep-placed": true, "--placed-by": true}
	placements  = []string{"--after", "--before", "--bottom", "--placed-by", "--top"}
)

// gitReads are the git commands a sub-agent may run.
var gitReads = []string{"blame", "cat-file", "describe", "diff", "grep", "log", "ls-files", "ls-tree", "merge-base", "rev-list", "rev-parse", "shortlog", "show", "status"}

// RolePlan is the planner's role, as flai serve sets it in FLAI_ROLE.
const RolePlan = "plan"

// Backlog is the one state the planner may move an item to.
const Backlog = "backlog"

// MCPPlans are flai's MCP tools the planner may call besides MCPReads;
// item_move it may call only to move an item to backlog.
var MCPPlans = []string{"activity_log", "inbox", "item_edit", "item_new", "thread_open", "thread_reply", "wait_for_events"}

// cliPlans are the flai commands the planner may run besides cliReads, in
// the form cliReads has; flai move it may run only to move an item to
// backlog.
var cliPlans = map[string][]string{
	"edit":    nil,
	"epic":    {"new"},
	"issue":   {"new", "bump"},
	"story":   {"new"},
	"task":    {"new"},
	"thread":  {"new", "reply"},
	"touches": nil,
}

// fileEdits are the tools that change files, which the planner never uses,
// and the analyzer only on its report (S-0223).
var fileEdits = []string{"Edit", "NotebookEdit", "Write"}

// fileWrites are the tools a sub-agent writes a file with, none of which it
// may point at a file in a .claude/ folder while auto-approve is off.
var fileWrites = []string{"Edit", "MultiEdit", "NotebookEdit", "Write"}

// moveValues are the flags of flai move, flai's own among them, that take
// a value.
var moveValues = map[string]bool{"--by": true, "--config": true, "--reason": true}

// planner ends each of the planner's refusals: the rule the call breaks and
// what to do instead.
const planner = "it plans through flai and never edits code, moves an item past backlog, accepts, or publishes (strategic-agents.md, ADR-0060). Ask the operator with thread_open on the item, or say it in your final summary."

// drafts says why the planner may create a story only as a draft (S-0209).
const drafts = "the stories it writes are drafts for the operator to finalize, so give item_new draft true, or flai story new --draft"

// RoleOrchestrate is the orchestrator's role, as flai serve sets it in
// FLAI_ROLE.
const RoleOrchestrate = "orchestrate"

// Ready is the one state the orchestrator may move an item to, with
// promote_to_ready.
const Ready = "ready"

// MCPOrchestrates are flai's MCP tools the orchestrator may call besides
// MCPReads, whatever its permissions.
var MCPOrchestrates = []string{"activity_log", "inbox", "thread_open", "wait_for_events"}

// cliOrchestrates are the flai commands the orchestrator may run besides
// cliReads, whatever its permissions, each with the subcommands it may run.
var cliOrchestrates = map[string][]string{
	"issue":  {"new", "bump"},
	"thread": {"new"},
}

// ReleasePublish is the MCP tool through which the orchestrator alone
// publishes, while orchestration.permissions.publish is on (S-0222).
const ReleasePublish = "release_publish"

// gitPublishes are the git commands that reach the remote, which the
// orchestrator leaves to release_publish.
var gitPublishes = []string{"push", "tag"}

// publishesThrough says why the orchestrator never publishes but through
// release_publish (S-0222).
const publishesThrough = "it publishes through release_publish alone, while orchestration.permissions.publish is on, and reaches the remote no other way"

// orchestratorPublishes says why a session other than the orchestrator's
// never calls release_publish (S-0222, ADR-0067).
const orchestratorPublishes = "it is the orchestrator's alone, while orchestration.permissions.publish is on; other sessions publish only when the operator asks, with flai push --pending (ADR-0067)"

// RoleAnalyze is the analyzer's role, as flai serve sets it in FLAI_ROLE
// (S-0223).
const RoleAnalyze = "analyze"

// MCPAnalyzes are flai's MCP tools the analyzer may call besides MCPReads:
// issue_new and issue_bump file and bump an issue for a finding of its
// report (S-0224).
var MCPAnalyzes = []string{"activity_log", "inbox", "issue_bump", "issue_new", "thread_open", "thread_reply", "wait_for_events"}

// itemWrites are flai's MCP tools that write a work item, issue_story's
// draft story among them, which the analyzer never calls.
var itemWrites = []string{"issue_story", "item_edit", "item_move", "item_new"}

// cliAnalyzes are the flai commands the analyzer may run besides cliReads,
// each with the subcommands it may run: it files and bumps the issues its
// findings call for (S-0224).
var cliAnalyzes = map[string][]string{
	"issue": {"new", "bump"},
}

// cliAuthors are the flai commands that author a work item, with the
// subcommands that do, nil for every one; the analyzer runs none of them.
var cliAuthors = map[string][]string{
	"epic":  nil,
	"issue": {"story"},
	"story": nil,
	"task":  nil,
}

// analyzer ends each of the analyzer's refusals: the rule the call breaks
// and what to do instead.
const analyzer = "it reads the project and its metrics, files or bumps an issue for each actionable finding with flai issue new and bump, and edits nothing but its report under the design folder's analysis/ and the folder's README.md (strategic-agents.md, ADR-0060). Put what it found, and the stories it would suggest, in its report, or ask the operator with thread_open."

// authorsNoStories says why the analyzer never writes a work item.
const authorsNoStories = "the analyzer authors no stories; the stories its findings call for are the planner's and the operator's to write, flai issue story among them"

// analyzing are the analyzer's rules: flai's reads, flai stats among them,
// issue new and bump, and git's reads.
var analyzing = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		if reads(cmd, sub, rest) {
			return "", ""
		}
		if subs, ok := cliAnalyzes[cmd]; ok && slices.Contains(subs, sub) {
			return "", ""
		}
		if subs, ok := cliAuthors[cmd]; ok && (subs == nil || slices.Contains(subs, sub)) {
			return authorsNoStories, ""
		}
		return "of flai's commands it runs only those that read, flai stats among them, and issue new and bump", ""
	},
	git: func(string) string { return "it runs only git's reads" },
}

// flaiValues are flai's own flags that take a value.
var flaiValues = map[string]bool{"--config": true}

// finalizeValues are the flags flai edit may carry besides --no-draft when it
// only finalizes a draft and that take a value; finalizeFlags those that do
// not.
var (
	finalizeValues = map[string]bool{"--by": true, "--config": true, "--hash": true}
	finalizeFlags  = []string{"--json", "--verbose", "--yes", "-v", "-y"}
)

// finalizeFields are the fields item_edit's input may give when it only
// finalizes a draft.
var finalizeFields = []string{"draft", "hash", "id", "project"}

// costFields are the fields item_edit's input may give when it only gives a
// story cost of delay inputs, and costInputs the keys of its cost_of_delay
// that are inputs (S-0328, ADR-0119).
var (
	costFields = []string{"cost_of_delay", "hash", "id", "project"}
	costInputs = []string{"penalty_per_week", "revenue_per_week", "time_lost_per_cycle"}
)

// costInputFlags are the flags of flai edit that set a cost of delay input
// (S-0328).
var costInputFlags = []string{"--penalty-per-week", "--revenue-per-week", "--time-lost-per-cycle"}

// editValues are the flags of flai edit, flai's own among them, that take a
// value and that an edit the orchestrator makes may give.
var editValues = map[string]bool{"--by": true, "--config": true, "--hash": true, "--penalty-per-week": true, "--revenue-per-week": true, "--time-lost-per-cycle": true}

// askOperator ends each of the orchestrator's refusals.
const askOperator = "Ask the operator with thread_open on the item if it needs doing (strategic-agents.md, ADR-0060)."

// plansCandidates says why the orchestrator asks for the planner on an epic
// or a story alone (S-0328).
const plansCandidates = "it asks for the planner on an epic or a story alone, one flai plan --candidates lists"

// Why the orchestrator never makes an edit, by MCP and on the command line:
// it only finalizes a draft, or gives a story cost of delay inputs (S-0219,
// S-0328, ADR-0119).
const (
	editsMCP = "it edits an item only to finalize a draft, with item_edit draft false and nothing else, or to give a story cost of delay inputs, with cost_of_delay's revenue_per_week, penalty_per_week, and time_lost_per_cycle and nothing else: no value, no clear_cost_of_delay"
	editsCLI = "it edits an item only to finalize a draft, with flai edit --no-draft and nothing else, or to give a story cost of delay inputs, with --revenue-per-week, --penalty-per-week, and --time-lost-per-cycle and nothing else: no --cost-of-delay-value, no --clear-cost-of-delay"
)

// byHand says why the orchestrator never places a story by hand in the pull
// order (S-0219).
const byHand = "it orders the ready column by its policy, with flai order --by <policy> --apply, and never places a story by hand, which is the operator's"

// movesTo says why the orchestrator moves an item nowhere but to ready, and
// a story to done only as it accepts it (S-0221).
const movesTo = "it moves an item to ready and no further, but a story it accepts, to done with --by orchestrator while orchestration.permissions.accept_reviews is on (ADR-0093)"

// acceptsAsSelf says why the orchestrator accepts a story only as itself,
// whatever accept_reviews lets it do (S-0221).
const acceptsAsSelf = "it accepts a story only as itself, so give --by orchestrator: orchestration.permissions.accept_reviews lets it accept as nobody else (ADR-0093)"

// acceptsWithFlai says why the orchestrator never moves an item to done with
// item_move (S-0221).
const acceptsWithFlai = "item_move never moves an item to done; it accepts a story with flai accept <id> --by orchestrator --verified <commit> --evidence <file>, while orchestration.permissions.accept_reviews is on (ADR-0093)"

// threadValues are the flags of flai thread reply, resolve, and confirm,
// flai's own among them, that take a value.
var threadValues = map[string]bool{"--by": true, "--config": true, "--reason": true, "--source": true}

// Why the orchestrator never makes a thread call, whatever its permissions
// (S-0220).
const (
	confirms     = "confirming a recommendation is the operator's alone (ADR-0090)"
	writesAsSelf = "it writes to a thread as itself, so give no --by naming another"
	resolvesOwn  = "it resolves only a thread it opened; one another opened is for its opener or the operator to resolve"
	answersOwn   = "on a thread it opened it only follows up, with neither recommendation nor source: it never recommends or answers its own question"
)

// Guard decides on the calls of one flai: Commands are the names of its
// commands, so that a word flai on a command line counts as running flai
// only when a command of its follows. Role is the session's role, from
// FLAI_ROLE: RolePlan holds the session's own calls to planning,
// RoleOrchestrate holds them to Permissions, the project's
// orchestration.permissions, RoleAnalyze to its report, and any other leaves
// them alone. Opener says
// who opened the thread with an ID, for the orchestrator's calls on threads,
// itself or a story's planner among others; ThreadOpener reads it from a
// project's threads. A thread whose opener it cannot tell, or every thread
// when Opener is nil, is held as another's.
// RoleAnalyze holds them to reads and the analyzer's report (S-0223): Root
// is the project's root and Reports its reports folder, relative to Root and
// slash separated, as analysis.Dir gives it; with either "" no file may be
// edited.
//
// Served says flai serve started the session, as StartedByEnv set to
// StartedByServe says; a Role or a Story says so too, since only flai serve
// sets them. A session flai serve started never changes the shared paths
// (ADR-0096).
//
// Story is the story a story's agent's session works, from FLAI_STORY, ""
// outside one; with Role "" it holds the session's own wait_for_events to
// the sub-agents Running lists and the threads ThreadOpen tells of. Running
// lists the sub-agents running in the session, as the package's Running
// reads them from the session's record; ThreadOpen says whether an
// unresolved thread is on Story or one of its tasks, as StoryThreadOpen
// reads it. A nil Running lists no sub-agent and a nil ThreadOpen no open
// thread; an error from either lets the wait through.
//
// AutoApprove says the operator enabled the project's auto-approve host
// action; while it is off a sub-agent's write of a file in a .claude/ folder
// is refused (S-0299).
type Guard struct {
	Commands    []string
	Role        string
	Served      bool
	Permissions manifest.Permissions
	Opener      func(id string) (string, error)
	Root        string
	Reports     string
	Story       string
	Running     func() ([]SubAgent, error)
	ThreadOpen  func() (bool, error)
	AutoApprove bool
}

// ThreadOpener says who opened the thread with an ID in r's wip/threads, in
// the main checkout, where every flai that writes to a thread reads it.
func ThreadOpener(r *workitem.Repo) func(id string) (string, error) {
	return func(id string) (string, error) {
		th, err := threads.Get(r, id)
		if err != nil {
			return "", err
		}
		return th.Opener(), nil
	}
}

// StoryThreadOpen says whether an unresolved thread in r's wip/threads is on
// story or on one of its tasks, as the inbox filtered by a story counts them.
func StoryThreadOpen(r *workitem.Repo, story string) func() (bool, error) {
	return func() (bool, error) {
		all, err := threads.List(r)
		if err != nil {
			return false, err
		}
		want := workitem.CanonicalID(story)
		return slices.ContainsFunc(all, func(th *threads.Thread) bool {
			return th.Open() && workitem.CanonicalID(threads.StoryOf(r, th)) == want
		}), nil
	}
}

// Refusal is the guard's decision on a call: Why it is refused, "" when it
// is not. For the orchestrator's own calls Call is the call refused, in one
// line, and Needs the permission that would allow it, as
// orchestration.permissions.<name>, or "" when none would.
type Refusal struct {
	Why   string
	Call  string
	Needs string
}

// rules are what one kind of caller may run on a command line: flai says
// why a flai command is refused, and the permission that would allow it, or
// "" when it is not refused, from its command, its subcommand, and the words
// after flai; git says why a git command other than a read is refused, from
// its subcommand, and nil refuses none.
type rules struct {
	flai func(cmd, sub string, rest []string) (why, needs string)
	git  func(cmd string) string
}

// subAgent are a sub-agent's rules: flai's reads and git's.
var subAgent = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		if reads(cmd, sub, rest) {
			return "", ""
		}
		return "flai commands that change work items, threads, narratives, or releases are the story's agent's", ""
	},
	git: func(string) string {
		return "git commands that change the worktree, the index, branches, or history are the story's agent's"
	},
}

// planning are the planner's rules: flai's reads and its planning commands,
// and git's reads.
var planning = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		switch {
		case reads(cmd, sub, rest):
			return "", ""
		case cmd == "move":
			if args := positionals(rest, moveValues); len(args) > 2 && args[2] == Backlog {
				return "", ""
			}
			return "it moves an item to backlog and no further", ""
		case cmd == "edit" && slices.ContainsFunc(rest, finalizes):
			return "finalizing a draft is the operator's", ""
		case cmd == "story" && sub == "new" && !drafted(rest):
			return drafts, ""
		}
		if subs, ok := cliPlans[cmd]; ok && (subs == nil || slices.Contains(subs, sub)) {
			return "", ""
		}
		return "of flai's commands that write, it runs only story new with --draft, epic new, task new, edit, touches, thread new and reply, issue new and bump, and move to backlog", ""
	},
	git: func(string) string { return "it runs only git's reads" },
}

// orchestration are the orchestrator's rules under its permissions: flai's
// reads, the commands it always runs, and those its permissions allow; and
// git's reads.
func (g Guard) orchestration() rules {
	return rules{
		flai: func(cmd, sub string, rest []string) (string, string) {
			if cmd == "thread" && slices.Contains([]string{"reply", "resolve", "confirm"}, sub) {
				return g.thread(cliThread(sub, rest))
			}
			needs, never := orchestrated(cmd, sub, rest)
			switch {
			case never != "":
				return nevers(never), ""
			case needs != "" && !g.Permissions.Allows(needs):
				return off(needs), needs
			case needs == manifest.PermitAcceptReviews && !workitem.IsOrchestrator(value(rest, "--by")):
				return acceptsAsSelf, needs
			}
			return "", ""
		},
		git: func(cmd string) string {
			if slices.Contains(gitPublishes, cmd) {
				return nevers(publishesThrough)
			}
			return nevers("it runs only git's reads")
		},
	}
}

// orchestrated says what the orchestrator needs to run a flai command: the
// permission that allows it, or why it never runs it; both "" when it runs
// it whatever its permissions.
func orchestrated(cmd, sub string, rest []string) (needs, never string) {
	if reads(cmd, sub, rest) {
		return "", ""
	}
	if subs, ok := cliOrchestrates[cmd]; ok && slices.Contains(subs, sub) {
		return "", ""
	}
	switch cmd {
	case "plan":
		args := positionals(rest, flaiValues)
		switch {
		case len(args) > 1 && isEpic(args[1]):
			return manifest.PermitPlanBacklogEpics, ""
		case len(args) > 1 && isStory(args[1]):
			return manifest.PermitPlanBacklogStories, ""
		}
		return "", plansCandidates
	case "edit":
		switch {
		case finalizesOnly(rest):
			return manifest.PermitFinalizeDrafts, ""
		case setsCostInputsOnly(rest):
			return manifest.PermitPlanBacklogStories, ""
		}
		return "", editsCLI
	case "move":
		args := positionals(rest, moveValues)
		switch {
		case len(args) > 2 && args[2] == Ready:
			return manifest.PermitPromoteToReady, ""
		case len(args) > 2 && args[2] == workitem.Done && isStory(args[1]):
			// the done transition of a story it accepts (ADR-0093); the
			// orchestration rules hold it to --by orchestrator
			return manifest.PermitAcceptReviews, ""
		}
		return "", movesTo
	case "order":
		if ordersByPolicy(rest) {
			return manifest.PermitOrderReady, ""
		}
		return "", byHand
	case "accept":
		return manifest.PermitAcceptReviews, ""
	case "push", "release":
		// flai release --evaluate is a read, let through above
		return "", publishesThrough
	}
	return "", "of flai's commands that write, it runs only thread new, issue new and bump, and those its permissions allow"
}

// threadCall is one of the orchestrator's calls on a thread: verb is reply,
// resolve, or confirm, and id the thread; for a reply, whether it is a
// recommendation and whether it cites a source; and by the author the call
// names in place of the orchestrator, "" when it names none.
type threadCall struct {
	verb, id, by            string
	recommendation, sourced bool
}

// cliThread is the thread call the words after flai make, of flai thread
// with its subcommand sub.
func cliThread(sub string, rest []string) threadCall {
	c := threadCall{verb: sub, by: value(rest, "--by"), recommendation: given(rest, "--recommend"), sourced: cites(value(rest, "--source"))}
	if args := positionals(rest, threadValues); len(args) > 2 {
		c.id = args[2]
	}
	return c
}

// cites says whether a reply's source names a file, as flai reads it.
func cites(source string) bool {
	s, err := threads.ParseSource(source)
	return err == nil && s.Path != ""
}

// value is the value words give flag last, as --flag value or --flag=value;
// "" when they give none.
func value(words []string, flag string) string {
	v := ""
	for i, w := range words {
		if w == flag && i+1 < len(words) {
			v = words[i+1]
		} else if s, ok := strings.CutPrefix(w, flag+"="); ok {
			v = s
		}
	}
	return v
}

// thread says why the orchestrator may not make a call on a thread, and the
// permission that would allow it; both "" when it may (S-0220). It never
// confirms a recommendation, writes as another, or resolves a thread it did
// not open but a story's planner's. On a thread it opened it follows up and
// resolves, whatever its permissions, but never recommends or answers. On a
// thread a story's planner opened, plan_backlog_stories lets it reply as an
// answer, with or without a source, or as a recommendation, and resolve the
// thread (S-0328, ADR-0119); while it is off, such a thread is held as
// another's, and a refusal names plan_backlog_stories, which would allow the
// call. On another's, while answer_threads is off it does not reply;
// recommend lets it reply with a recommendation; and autonomous lets it
// answer too, citing a source, so that an answer it cannot source goes to
// the operator as a recommendation.
func (g Guard) thread(c threadCall) (why, needs string) {
	switch {
	case c.verb == "confirm":
		return nevers(confirms), ""
	case c.by != "" && c.by != workitem.ActivityOrchestrator:
		return nevers(writesAsSelf), ""
	}
	who, stories := g.opener(c.id), manifest.PermitPlanBacklogStories
	own, planner := who == workitem.ActivityOrchestrator, storyPlanner(who)
	switch {
	case own && (c.recommendation || c.sourced):
		return nevers(answersOwn), ""
	case own, planner && g.Permissions.Allows(stories):
		return "", ""
	case planner && (c.verb == "resolve" || !g.Permissions.Allows(manifest.PermitAnswerThreads)):
		return off(stories), stories
	case c.verb == "resolve":
		return nevers(resolvesOwn), ""
	}
	why, needs = g.answers(c)
	if why != "" && planner {
		return why + "; on a thread a story's planner opened, orchestration.permissions." + stories + ", which is off, would let it reply so", stories
	}
	return why, needs
}

// answers says why answer_threads refuses the orchestrator a reply on
// another's thread, and the permission that would allow it; both "" when it
// may make it (S-0220).
func (g Guard) answers(c threadCall) (why, needs string) {
	mode, answers := g.Permissions.AnswerMode(), manifest.PermitAnswerThreads
	switch {
	case !g.Permissions.Allows(answers):
		return off(answers), answers
	case c.recommendation, mode == manifest.AnswerAutonomous && c.sourced:
		return "", ""
	case mode == manifest.AnswerAutonomous:
		return "it answers another's thread only citing what the answer rests on, so give source <path> or <path>#<heading> (--source on the command line), or post it as a recommendation for the operator to confirm with recommendation true (--recommend)", ""
	case c.sourced:
		return fmt.Sprintf("orchestration.permissions.%s is %s, so it replies to another's thread only as a recommendation for the operator to confirm: give recommendation true (--recommend); %s would let it answer citing its source", answers, mode, manifest.AnswerAutonomous), answers
	}
	return fmt.Sprintf("orchestration.permissions.%s is %s, so it replies to another's thread only as a recommendation for the operator to confirm: give recommendation true (--recommend); %s would let it answer only citing a source", answers, mode, manifest.AnswerAutonomous), ""
}

// opener is who opened the thread with id, as Opener tells; "" for a thread
// Opener cannot tell of, which is another's.
func (g Guard) opener(id string) string {
	if g.Opener == nil || id == "" {
		return ""
	}
	who, err := g.Opener(id)
	if err != nil {
		return ""
	}
	return who
}

// storyPlanner says whether who is a story's planner, as flai serve names
// the planner it runs for a story: planner-S-nnnn (S-0328).
func storyPlanner(who string) bool {
	id, ok := strings.CutPrefix(who, workitem.ActivityPlanner+"-")
	return ok && isStory(id)
}

// threadDoes is what a thread call does, said after "the orchestrator
// cannot".
func threadDoes(c threadCall) string {
	id := c.id
	if id == "" {
		id = "a thread"
	}
	switch {
	case c.verb != "reply":
		return c.verb + " " + id
	case c.recommendation:
		return "recommend an answer to " + id
	}
	return "reply to " + id
}

// ordersByPolicy says whether the words after flai, of flai order, order the
// ready column by a policy rather than place one story by hand: with --by,
// and with neither a story nor a flag that places one.
func ordersByPolicy(words []string) bool {
	return given(words, "--by") && len(positionals(words, orderValues)) == 1 &&
		!slices.ContainsFunc(placements, func(f string) bool { return named(words, f) })
}

// finalizesDraft says whether item_edit's input finalizes a draft and
// changes nothing else: it gives draft false, and no field but the item, its
// hash, and its project besides.
func finalizesDraft(in Input) bool {
	return !in.Draft && slices.Contains(in.Fields, "draft") &&
		!slices.ContainsFunc(in.Fields, func(f string) bool { return !slices.Contains(finalizeFields, f) })
}

// setsCostInputs says whether item_edit's input gives a story cost of delay
// inputs and changes nothing else (S-0328, ADR-0119): cost_of_delay with one
// input or more, none null or an empty string, which removes it, and no other
// key, value among them; and no field but the story, its hash, and its
// project besides. Whether the story may have them is flai's to judge.
func setsCostInputs(in Input) bool {
	if !isStory(in.ID) || len(in.CostOfDelay) == 0 ||
		slices.ContainsFunc(in.Fields, func(f string) bool { return !slices.Contains(costFields, f) }) {
		return false
	}
	for k, raw := range in.CostOfDelay {
		var v any
		if !slices.Contains(costInputs, k) || json.Unmarshal(raw, &v) != nil || v == nil || v == "" {
			return false
		}
	}
	return true
}

// isEpic says whether id is an epic's, in any zero padding.
func isEpic(id string) bool { return isOf(id, "E-") }

// isStory says whether id is a story's, in any zero padding.
func isStory(id string) bool { return isOf(id, "S-") }

// isOf says whether id is a number after prefix, in either case.
func isOf(id, prefix string) bool {
	n, ok := strings.CutPrefix(strings.ToUpper(id), prefix)
	_, err := strconv.Atoi(n)
	return ok && err == nil
}

// finalizesOnly says whether the words of flai edit finalize a draft and
// change nothing else.
func finalizesOnly(words []string) bool {
	if !given(words, "--no-draft") {
		return false
	}
	for i := 0; i < len(words); i++ {
		w := words[i]
		name, _, valued := strings.Cut(w, "=")
		switch {
		case !strings.HasPrefix(w, "-"), finalizes(w), slices.Contains(finalizeFlags, w):
		case finalizeValues[name]:
			if !valued {
				i++
			}
		default:
			return false
		}
	}
	return true
}

// setsCostInputsOnly says whether the words of flai edit give a story cost
// of delay inputs and change nothing else (S-0328, ADR-0119): one input flag
// or more, each with a value that is not empty, which would remove it,
// beside the flags that finalizing a draft allows but --no-draft.
func setsCostInputsOnly(words []string) bool {
	if args := positionals(words, editValues); len(args) < 2 || !isStory(args[1]) {
		return false
	}
	inputs := 0
	for i := 0; i < len(words); i++ {
		w := words[i]
		name, v, valued := strings.Cut(w, "=")
		switch {
		case !strings.HasPrefix(w, "-"), slices.Contains(finalizeFlags, w):
		case slices.Contains(costInputFlags, name):
			if !valued {
				if i+1 == len(words) {
					return false
				}
				i++
				v = words[i]
			}
			if v == "" {
				return false
			}
			inputs++
		case finalizeValues[name]:
			if !valued {
				i++
			}
		default:
			return false
		}
	}
	return inputs > 0
}

// reads says whether a flai command with its subcommand and the words after
// flai only reads.
func reads(cmd, sub string, rest []string) bool {
	if f, ok := flagReads[cmd]; ok {
		return (len(f.with) == 0 || slices.ContainsFunc(f.with, func(w string) bool { return given(rest, w) })) &&
			!slices.ContainsFunc(f.without, func(w string) bool { return named(rest, w) }) &&
			(f.form == nil || f.form(rest))
	}
	subs, ok := cliReads[cmd]
	return ok && (subs == nil || slices.Contains(subs, sub))
}

// given says whether words give flag, as the last of its occurrences leaves
// it: bare, or with a value that is neither empty nor false.
func given(words []string, flag string) bool {
	on := false
	for _, w := range words {
		if w == flag {
			on = true
		} else if v, ok := strings.CutPrefix(w, flag+"="); ok {
			b, err := strconv.ParseBool(v)
			on = v != "" && (err != nil || b)
		}
	}
	return on
}

// named says whether words name flag at all, whatever its value.
func named(words []string, flag string) bool {
	return slices.ContainsFunc(words, func(w string) bool { return w == flag || strings.HasPrefix(w, flag+"=") })
}

// finalizes says whether a word of flai edit finalizes a draft.
func finalizes(w string) bool {
	return w == "--no-draft" || strings.HasPrefix(w, "--no-draft=")
}

// drafted says whether the words of flai story new make the story a draft:
// the last --draft decides, bare or with a value that is true. A value that
// is not a boolean, which flai refuses anyway, is not a draft.
func drafted(words []string) bool {
	draft := false
	for _, w := range words {
		if w == "--draft" {
			draft = true
		} else if v, ok := strings.CutPrefix(w, "--draft="); ok {
			draft, _ = strconv.ParseBool(v)
		}
	}
	return draft
}

// Check says why a call is refused, or "" when it is not.
func (g Guard) Check(e Event) string { return g.Decide(e).Why }

// Decide decides on a call: the zero Refusal lets it through.
func (g Guard) Decide(e Event) Refusal {
	if r := g.shared(e); r.Why != "" {
		return r
	}
	if e.AgentID == "" {
		if e.ToolName == MCPPrefix+ReleasePublish && g.Role != RoleOrchestrate && g.Role != RolePlan && g.Role != RoleAnalyze {
			return Refusal{Why: "only the orchestrator calls " + ReleasePublish + ": " + orchestratorPublishes}
		}
		switch g.Role {
		case RolePlan:
			return Refusal{Why: g.plan(e)}
		case RoleAnalyze:
			return Refusal{Why: g.analyze(e)}
		case RoleOrchestrate:
			return g.orchestrate(e)
		case "":
			return Refusal{Why: g.wait(e)}
		}
		return Refusal{}
	}
	who := e.AgentType
	if who == "" {
		who = "unnamed"
	}
	if file, ok := ProtectedWrite(e); ok && !g.AutoApprove {
		return Refusal{Why: fmt.Sprintf("a sub-agent (%s) cannot use %s on %s: a path Claude Code protects, such as a file in a .claude/ folder or .mcp.json, is written only with the operator's approval on a thread, which would hold this call, and the layer with it, until they answer (ADR-0086, ADR-0106). Put the file's whole new content in your final message; the story's agent writes it.", who, e.ToolName, file)}
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		if slices.Contains(MCPReads, tool) {
			return Refusal{}
		}
		return Refusal{Why: fmt.Sprintf("a sub-agent (%s) cannot call %s: it reads, and only the story's agent changes work items and threads or reads the inbox (ADR-0059). Put what you need done, or the question for the designer, in your final message.", who, tool)}
	}
	if e.ToolName != "Bash" {
		return Refusal{}
	}
	for _, words := range commands(e.ToolInput.Command) {
		if why, _ := g.refuse(words, subAgent); why != "" {
			return Refusal{Why: fmt.Sprintf("a sub-agent (%s) cannot run %q: %s (ADR-0060). Put what you need done in your final message; the story's agent does it.", who, strings.Join(words, " "), why)}
		}
	}
	return Refusal{}
}

// ProtectedWrite says whether e writes a path Claude Code protects other than
// one in .git, such as a file in a .claude/ folder or .mcp.json, and which
// file (S-0299, ADR-0106): Claude Code asks a person before any such write,
// so that a sub-agent's would wait on a thread for the operator.
// permission_prompt refuses a path in .git at once, so that one never waits.
func ProtectedWrite(e Event) (string, bool) {
	if !slices.Contains(fileWrites, e.ToolName) {
		return "", false
	}
	file := cmp.Or(e.ToolInput.FilePath, e.ToolInput.NotebookPath)
	if file == "" {
		return "", false
	}
	clean := filepath.Clean(file)
	return file, protected.Path(clean) && !protected.Git(clean)
}

// waitsOnSubAgents ends the refusal of a story's agent's wait_for_events: why
// it would not return, and how to wait instead.
const waitsOnSubAgents = "wait_for_events reports work items and threads, not sub-agents, so it would run to its timeout. Wait for a sub-agent by launching it with the Agent tool's run_in_background set to false, so that its result comes back as the tool's result; launch a layer's sub-agents in one message, each so. Use wait_for_events only for a thread awaiting the designer (delegation.md)."

// wait says why a story's agent's own call is refused, or "" when it is not
// (S-0285): wait_for_events while a sub-agent of its session runs, unless a
// thread on its story or one of its tasks is open, so that it waits on the
// designer, whose answer the call returns on.
func (g Guard) wait(e Event) string {
	if g.Story == "" || e.ToolName != MCPPrefix+"wait_for_events" || g.Running == nil {
		return ""
	}
	running, err := g.Running()
	if err != nil || len(running) == 0 {
		return ""
	}
	if g.ThreadOpen != nil {
		if open, err := g.ThreadOpen(); err != nil || open {
			return ""
		}
	}
	names := make([]string, len(running))
	for i, s := range running {
		who := s.Type
		if who == "" {
			who = "unnamed"
		}
		names[i] = fmt.Sprintf("%s %s", who, s.ID)
	}
	return fmt.Sprintf("the story's agent cannot call wait_for_events while its sub-agents run (%s) and no thread on %s awaits the designer: %s", strings.Join(names, ", "), g.Story, waitsOnSubAgents)
}

// orchestrate decides on the orchestrator's own call under its permissions.
func (g Guard) orchestrate(e Event) Refusal {
	if slices.Contains(fileEdits, e.ToolName) {
		call := strings.TrimSpace(e.ToolName + " " + e.ToolInput.FilePath)
		return refusedOrchestrator(call, "use "+e.ToolName, nevers("it acts through flai and never edits a file"), "")
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		in := e.ToolInput
		call, what := strings.Join(strings.Fields(tool+" "+in.ID+" "+in.To), " "), "call "+tool
		var needs, never string
		switch {
		case slices.Contains(MCPReads, tool), slices.Contains(MCPOrchestrates, tool):
			return Refusal{}
		case tool == "thread_reply", tool == "thread_resolve":
			c := threadCall{verb: strings.TrimPrefix(tool, "thread_"), id: in.ID, recommendation: in.Recommendation, sourced: cites(in.Source)}
			why, permit := g.thread(c)
			if why == "" {
				return Refusal{}
			}
			if c.recommendation {
				call += " recommendation"
			}
			if c.sourced {
				call += " source " + strings.TrimSpace(in.Source)
			}
			return refusedOrchestrator(call, threadDoes(c), why, permit)
		case tool == "plan":
			what = "plan " + in.ID
			switch {
			case isEpic(in.ID):
				needs = manifest.PermitPlanBacklogEpics
			case isStory(in.ID):
				needs = manifest.PermitPlanBacklogStories
			default:
				never = plansCandidates
			}
		case tool == "item_edit":
			what = "edit " + in.ID
			switch {
			case finalizesDraft(in):
				what, needs = "finalize "+in.ID, manifest.PermitFinalizeDrafts
			case setsCostInputs(in):
				what, needs = "give "+in.ID+" cost of delay inputs", manifest.PermitPlanBacklogStories
			default:
				never = editsMCP
			}
		case tool == "item_move":
			what = fmt.Sprintf("move %s to %s", in.ID, in.To)
			switch in.To {
			case Ready:
				needs = manifest.PermitPromoteToReady
			case workitem.Done:
				never = acceptsWithFlai
			default:
				never = "it moves an item to ready and no further"
			}
		case tool == ReleasePublish:
			what, needs = "publish a release", manifest.PermitPublish
		default:
			never = "of flai's tools that write, it calls only thread_open and activity_log, thread_reply and thread_resolve as the thread allows, and those its permissions allow"
		}
		switch {
		case never != "":
			return refusedOrchestrator(call, what, nevers(never), "")
		case !g.Permissions.Allows(needs):
			return refusedOrchestrator(call, what, off(needs), needs)
		}
		return Refusal{}
	}
	if e.ToolName != "Bash" {
		return Refusal{}
	}
	r := g.orchestration()
	for _, words := range commands(e.ToolInput.Command) {
		if why, needs := g.refuse(words, r); why != "" {
			call := strings.Join(words, " ")
			return refusedOrchestrator(call, fmt.Sprintf("run %q", call), why, needs)
		}
	}
	return Refusal{}
}

// refusedOrchestrator is the orchestrator's refusal of call, said as what,
// for why; needs is the permission that would allow it, "" when none would.
func refusedOrchestrator(call, what, why, needs string) Refusal {
	r := Refusal{Call: call, Why: fmt.Sprintf("the orchestrator cannot %s: %s. %s", what, why, askOperator)}
	if needs != "" {
		r.Needs = "orchestration.permissions." + needs
	}
	return r
}

// off says the orchestrator is refused a call for want of the permission
// needs.
func off(needs string) string {
	return "it needs orchestration.permissions." + needs + ", which is off"
}

// nevers says the orchestrator never makes a call, whatever its
// permissions, for why.
func nevers(why string) string {
	return "the orchestrator never does it, whatever its permissions: " + why
}

// plan says why the planner's own call is refused, or "" when it is not.
func (g Guard) plan(e Event) string {
	if slices.Contains(fileEdits, e.ToolName) {
		return fmt.Sprintf("the planner cannot use %s: %s", e.ToolName, planner)
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		switch {
		case tool == "item_new" && e.ToolInput.Type == "story" && !e.ToolInput.Draft:
			return fmt.Sprintf("the planner cannot create a story that is not a draft: %s; %s", drafts, planner)
		case slices.Contains(MCPReads, tool), slices.Contains(MCPPlans, tool):
			return ""
		case tool == "item_move" && e.ToolInput.To == Backlog:
			return ""
		case tool == "item_move":
			return fmt.Sprintf("the planner cannot move %s to %s: %s", e.ToolInput.ID, e.ToolInput.To, planner)
		}
		return fmt.Sprintf("the planner cannot call %s: %s", tool, planner)
	}
	if e.ToolName != "Bash" {
		return ""
	}
	for _, words := range commands(e.ToolInput.Command) {
		if why, _ := g.refuse(words, planning); why != "" {
			return fmt.Sprintf("the planner cannot run %q: %s; %s", strings.Join(words, " "), why, planner)
		}
	}
	return ""
}

// analyze says why the analyzer's own call is refused, or "" when it is not
// (S-0223).
func (g Guard) analyze(e Event) string {
	if slices.Contains(fileEdits, e.ToolName) {
		file := e.ToolInput.FilePath
		if file == "" {
			file = e.ToolInput.NotebookPath
		}
		if g.inReports(file) {
			return ""
		}
		folder := "the design folder's analysis/"
		if g.Reports != "" {
			folder = strings.TrimSuffix(g.Reports, "/") + "/"
		}
		return fmt.Sprintf("the analyzer cannot use %s on %s: it edits only its report and the index, under %s; %s", e.ToolName, cmp.Or(file, "no file"), folder, analyzer)
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		switch {
		case slices.Contains(MCPReads, tool), slices.Contains(MCPAnalyzes, tool):
			return ""
		case slices.Contains(itemWrites, tool):
			return fmt.Sprintf("the analyzer cannot call %s: %s; %s", tool, authorsNoStories, analyzer)
		}
		return fmt.Sprintf("the analyzer cannot call %s: %s", tool, analyzer)
	}
	if e.ToolName != "Bash" {
		return ""
	}
	for _, words := range commands(e.ToolInput.Command) {
		if why, _ := g.refuse(words, analyzing); why != "" {
			return fmt.Sprintf("the analyzer cannot run %q: %s; %s", strings.Join(words, " "), why, analyzer)
		}
	}
	return ""
}

// inReports says whether file, as Edit, Write, or NotebookEdit names it, is
// under the reports folder: taken from Root when it is relative, cleaned, and
// with the symbolic links along the part of it that exists resolved, as the
// folder's are, so that neither .. nor a link leads out of it.
func (g Guard) inReports(file string) bool {
	if g.Root == "" || g.Reports == "" || file == "" {
		return false
	}
	if !filepath.IsAbs(file) {
		file = filepath.Join(g.Root, file)
	}
	dir := resolved(filepath.Join(g.Root, filepath.FromSlash(g.Reports)))
	rel, err := filepath.Rel(dir, resolved(filepath.Clean(file)))
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// resolved is p, clean and absolute, with the symbolic links of the longest
// part of it that exists resolved and the rest joined on as it stands.
func resolved(p string) string {
	rest := ""
	for {
		if r, err := filepath.EvalSymlinks(p); err == nil {
			return filepath.Join(r, rest)
		}
		parent := filepath.Dir(p)
		if parent == p {
			return filepath.Join(p, rest)
		}
		rest = filepath.Join(filepath.Base(p), rest)
		p = parent
	}
}

var (
	assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
	subcommand = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	shellFlags = regexp.MustCompile(`^-[a-zA-Z]*c[a-zA-Z]*$`)
)

// refuse says why r refuses one simple command, and the permission that
// would allow it, or "" when it does not. It does not trust the first word
// to be the program: wrappers such as env, sudo, timeout, xargs, and find
// -exec put it further along, so every word is looked at. A word is git run
// with a subcommand when a plain word follows it after git's own flags, and
// flai run with one when one of its commands follows; a shell given -c, or
// eval, has its script checked too.
func (g Guard) refuse(words []string, r rules) (why, needs string) {
	for i, w := range words {
		if assignment.MatchString(w) {
			continue
		}
		rest := words[i+1:]
		switch path.Base(w) {
		case "sh", "bash", "zsh", "dash", "ksh":
			for j, f := range rest {
				if shellFlags.MatchString(f) && j+1 < len(rest) {
					if why, needs := g.script(rest[j+1], r); why != "" {
						return why, needs
					}
					break
				}
			}
		case "eval":
			if why, needs := g.script(strings.Join(rest, " "), r); why != "" {
				return why, needs
			}
		case "flai", "flai.sh":
			cmd, sub := subcommands(rest, flaiValues)
			if !slices.Contains(g.Commands, cmd) || slices.Contains(rest, "--help") || slices.Contains(rest, "-h") {
				continue
			}
			if why, needs := r.flai(cmd, sub, rest); why != "" {
				return why, needs
			}
		case "git":
			cmd, _ := subcommands(rest, map[string]bool{"-C": true, "-c": true, "--git-dir": true, "--work-tree": true, "--namespace": true})
			if r.git == nil || !subcommand.MatchString(cmd) || slices.Contains(gitReads, cmd) {
				continue
			}
			return r.git(cmd), ""
		}
	}
	return "", ""
}

// script checks each simple command of a script a shell is given against r.
func (g Guard) script(text string, r rules) (why, needs string) {
	for _, c := range commands(text) {
		if why, needs := g.refuse(c, r); why != "" {
			return why, needs
		}
	}
	return "", ""
}

// subcommands is the first two words after a program's own flags; a flag in
// takesValue consumes the word after it.
func subcommands(words []string, takesValue map[string]bool) (cmd, sub string) {
	var found []string
	for i := 0; i < len(words) && len(found) < 2; i++ {
		w := words[i]
		if strings.HasPrefix(w, "-") {
			if takesValue[w] && len(found) == 0 {
				i++
			}
			continue
		}
		found = append(found, w)
	}
	if len(found) > 0 {
		cmd = found[0]
	}
	if len(found) > 1 {
		sub = found[1]
	}
	return cmd, sub
}

// positionals are the words that are not flags; a flag in takesValue
// consumes the word after it, wherever it stands.
func positionals(words []string, takesValue map[string]bool) []string {
	var out []string
	for i := 0; i < len(words); i++ {
		if strings.HasPrefix(words[i], "-") {
			if takesValue[words[i]] {
				i++
			}
			continue
		}
		out = append(out, words[i])
	}
	return out
}

// shells are the programs that run a script they read, so that a heredoc fed
// to one is read as commands.
var shells = []string{"sh", "bash", "zsh", "dash", "ksh"}

// heredoc is a here-document a line opens: the delimiter that ends its body,
// with quotes removed, and whether <<- strips the body's leading tabs.
type heredoc struct {
	delim string
	tabs  bool
}

// commands splits a shell command line into simple commands, each as its
// words with quotes removed: at ;, &, |, newlines, and the openings of
// subshells and command substitutions. It is not a shell; it is enough to
// find each program a line runs and the words after it. A heredoc's body is
// input, not commands, so it is left out, unless a shell is named on the
// line that opens it, as in bash <<EOF or cat <<EOF | sh, when its body is
// split as a script (S-0246).
func commands(line string) [][]string {
	var out [][]string
	var words []string
	var word strings.Builder
	inWord := false
	var quote rune
	var pending []heredoc
	start, arith, skip := 0, 0, 0
	end := func() {
		if inWord {
			words = append(words, word.String())
			word.Reset()
			inWord = false
		}
	}
	split := func() {
		end()
		if len(words) > 0 {
			out = append(out, words)
		}
		words = nil
	}
	for i, r := range line {
		if i < skip {
			continue
		}
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			word.WriteRune(r)
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == '<' && arith == 0 && strings.HasPrefix(line[i:], "<<") && !strings.HasPrefix(line[i:], "<<<") && !strings.HasSuffix(line[:i], "<"):
			end()
			h, next, ok := opens(line, i)
			if !ok {
				word.WriteString("<<")
				inWord, skip = true, i+2
				continue
			}
			pending, skip = append(pending, h), next
		case r == '\n' && len(pending) > 0:
			split()
			body, next := bodies(line, i+1, pending)
			if slices.ContainsFunc(out[start:], runsShell) {
				out = append(out, commands(body)...)
			}
			pending, start, skip = nil, len(out), next
		case r == '(' && strings.HasPrefix(line[i:], "(("):
			arith++
			split()
		case r == ')' && arith > 0 && strings.HasPrefix(line[i:], "))"):
			arith--
			split()
		case r == ';' || r == '&' || r == '|' || r == '\n' || r == '(' || r == ')' || r == '`' || r == '{' || r == '}':
			split()
			if r == '\n' {
				start = len(out)
			}
		case r == '$':
			end()
		case r == ' ' || r == '\t':
			end()
		default:
			word.WriteRune(r)
			inWord = true
		}
	}
	split()
	return out
}

// opens reads the heredoc operator at line[i:], << or <<-, and the delimiter
// word after it: the heredoc, and where the line goes on after the word. It
// is not one when no word follows.
func opens(line string, i int) (h heredoc, next int, ok bool) {
	j := i + 2
	if j < len(line) && line[j] == '-' {
		h.tabs = true
		j++
	}
	for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
		j++
	}
	var delim strings.Builder
	var quote byte
	for ; j < len(line); j++ {
		c := line[j]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
				continue
			}
			delim.WriteByte(c)
		case c == '\'' || c == '"':
			quote = c
		case c == '\\' && j+1 < len(line):
			j++
			delim.WriteByte(line[j])
		case strings.IndexByte(" \t\n;&|<>()`", c) >= 0:
			h.delim = delim.String()
			return h, j, h.delim != ""
		default:
			delim.WriteByte(c)
		}
	}
	h.delim = delim.String()
	return h, j, h.delim != ""
}

// bodies reads the bodies of the heredocs pending, in order, from line[from:]
// on: their text, and where the line goes on after the last one's delimiter.
// A body with no delimiter runs to the end, as in bash.
func bodies(line string, from int, pending []heredoc) (string, int) {
	var text strings.Builder
	at := from
	for _, h := range pending {
		for at < len(line) {
			l, _, _ := strings.Cut(line[at:], "\n")
			at = min(at+len(l)+1, len(line))
			mark := strings.TrimSuffix(l, "\r")
			if h.tabs {
				mark = strings.TrimLeft(mark, "\t")
			}
			if mark == h.delim {
				break
			}
			text.WriteString(l + "\n")
		}
	}
	return text.String(), at
}

// runsShell says whether a simple command names a shell, which reads a
// heredoc fed to it, or piped to it, as a script.
func runsShell(words []string) bool {
	return slices.ContainsFunc(words, func(w string) bool {
		return slices.Contains(shells, path.Base(w))
	})
}
