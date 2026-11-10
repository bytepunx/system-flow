// Package guard refuses what a sub-agent may not do (ADR-0060): Claude Code
// runs flai guard before each tool call in a project made from the template,
// with the call on standard input, and says which sub-agent makes it. A
// sub-agent reads; it does not change a work item, a thread, a conversation
// with another story's agent (S-0331), or the repository's history, because
// the story's agent alone acts for the story (ADR-0059). The story's agent's
// own calls carry no agent ID and are refused only as said below. The guard
// is not a shell: it finds the programs a command line runs well enough to
// stop a sub-agent that follows its instructions carelessly, not one that
// sets out to hide a command (a backslash inside a name, a command held in a
// variable).
//
// A harness's Reader reads its hook's input into a Call, which Decide decides
// on in no harness's terms, and answers the harness with the Verdict
// (S-0353); Claude Code's, under the harness name claude-code, is the only
// one.
//
// A sub-agent may write an ADR, a document like the design and the code it
// edits, with flai adr new, topics, and accept, or the MCP tool adr_new, but
// not commit it (S-0287, ADR-0127): with --commit or --autocommit, or
// adr_new's commit, they are refused, since a commit writes history and
// widens the story's touches, which are the story's agent's. Every other
// flai adr command is refused it as before.
//
// A sub-agent does not write a file in a .claude/ folder, with Edit,
// MultiEdit, Write, or NotebookEdit, unless the operator enabled the
// project's auto-approve host action (S-0299): Claude Code asks for each such
// write, and flai's permission_prompt asks the operator on a thread, which
// would hold the call, and the sub-agent's layer with it, until they answer.
// The sub-agent puts the file's content in its final message instead, for
// the story's agent to write.
//
// A read-only sub-agent, one whose definition grants neither Edit nor Write,
// as the explorer's and the verifier's do, or Claude Code's Explore or Plan,
// writes no file through the shell but under the temporary folder (S-0368,
// ADR-0141): a redirection, tee, sed -i, perl -i, cp, mv, rm, and the like,
// or a formatter's write flag, to any other target is refused it. A task
// sub-agent edits files by design, and writes through the shell as before.
//
// A planner session, one flai serve starts with the role plan, is held to
// planning (strategic-agents.md): its own calls, which carry no agent ID, may
// read, write work items and threads through flai as itself, and move an item
// to backlog, but not edit files with Edit, MultiEdit, Write, or
// NotebookEdit, move an item further, accept, publish, work a story, or
// change the repository's history. A story it creates is a draft for the
// operator to finalize (S-0209). Since S-0416 (ADR-0207) these rules are read
// from its role file, planner.md under the design folder's workflows/roles,
// or from flai's baseline copy of it when the project has none: it edits an
// item only with an edit that neither finalizes a draft nor sets a cost of
// delay input (ungated-edit), and a refusal the rules below would make is
// worded as they word it. Its sub-agents are held as every sub-agent is.
//
// An orchestrator session, one flai serve starts with the role orchestrate,
// is held to the permissions the operator gives it in the manifest's
// orchestration.permissions (S-0218): its own calls may read, open threads,
// record issues, and log its activities, and do each of the rest only while
// the permission that allows it is on (S-0219 maps planning an epic,
// finalizing a draft, promoting a story to ready, and ordering the ready
// column by a policy). With plan_backlog_stories it asks for the planner on
// a story. It gives a story cost of delay inputs with an edit that changes
// nothing else, never a value or a removal, with plan_backlog_stories or
// set_cost_of_delay, and an epic with set_cost_of_delay alone (S-0328,
// ADR-0119, S-0523, ADR-0224); flai judges which epic or story each allows.
// With accept_reviews it accepts a story in review, with flai accept or flai
// move to done, only as itself, --by orchestrator (S-0221, ADR-0093), and
// enqueues one for the merge queue to land, with queue_add or flai queue add,
// as itself too; it never lands one with flai accept --batch, nor takes one
// out of the queue, holds it, or resumes it, whatever its permissions
// (S-0403). A refusal names that permission; a call
// no permission allows, such as an edit of a file, a commit, or a story
// placed by hand in the pull order, it never makes. With publish it
// publishes through release_publish alone, and never with flai release,
// flai push, git push, or git tag that makes a tag, whatever its
// permissions (S-0222), though git tag that lists is a read (S-0623); no
// session with no role, nor a story's or a verifier's, calls
// release_publish. These rules are read from its role file, orchestrator.md,
// or flai's baseline copy of it, as the planner's are (S-0416), each
// permission entry's setting from orchestration.permissions at each call; a
// refusal names the call and the permission that would allow it, so that
// flai guard logs it in the orchestrator's activity document. Its sub-agents
// are held as every sub-agent is.
//
// An analyzer session, one flai serve starts with the role analyze, is held
// to its report (S-0223): its own calls may read, flai stats among them, log
// its activities, and open and reply to threads as itself, and may edit, with
// Edit, MultiEdit, Write, or NotebookEdit, only a file under the manifest's
// design folder's analysis/, its report and the folder's README.md that lists
// it. It files and bumps an issue for each actionable finding (S-0224), with
// flai issue new and bump or the MCP tools issue_new and issue_bump, which
// write the issue through flai, never through a file edit. It authors no
// stories:
// item_new, item_edit, item_move, issue_story, and flai issue story, story,
// epic, and task are refused it, as is every other write through flai, flai
// issue close among them, and git's. These rules are read from its role
// file, analyzer.md, or flai's baseline copy of it, as the planner's are
// (S-0416), with its files entry's {design} the design folder. Its
// sub-agents, the explorer it hands search to among them, are held as every
// sub-agent is.
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
// A curator session, one flai serve starts with the role curate, is held to
// the digest (S-0431, ADR-0196): its own calls may read, log its activities,
// open and reply to threads as itself, and run flai digest harvest, new, and
// edit, which list the candidates and write or merge an entry through flai,
// and, since S-0432 (ADR-0199), absorb, cited, and compact, which mark an
// entry absorbed, record each entry's citers, and compact the entries no one
// cited within the window. It edits no file, with Edit, MultiEdit, Write, or
// NotebookEdit, and authors no work item: flai digest seed, every write
// through flai but those, the tools and commands that start a strategic
// agent among them, and git's writes are refused it, each refusal saying
// that it writes entries only with flai digest new, edit, absorb, cited, and
// compact. These rules are read from its role file, curator.md, or flai's
// baseline copy of it, as the planner's are (S-0416). Its sub-agents are held
// as every sub-agent is.
//
// A session flai serve starts with a role a project's role file defines,
// under the design folder's workflows/roles (S-0416, ADR-0207), is held to
// that file: its own calls may make the reads every sub-agent has, and each
// MCP tool, flai command, and file edit the file lists, under the condition
// an entry names and, for an entry under a permission, while the operator's
// setting for it is on, roles.<name>.permissions.<p> for the role itself. An
// entry flai check reports as roles.grant is never granted. Every other call
// is refused, naming the role, the rule it breaks, and the operator to ask
// with thread_open; a role whose file is missing or does not parse is refused
// every call but a read, naming the file and flai check. Its sub-agents are
// held as every sub-agent is.
//
// Every role, the orchestrator included, is so held to its role file: the
// file decides, and the rules in Go of the planner, the orchestrator, the
// analyzer, and the curator only word a refusal of the file's that they
// make too. A call on a thread under own-thread, planner-thread, or
// answer-threads names no author with --by but the role itself.
//
// A story's agent, in a session flai serve starts with FLAI_STORY and no
// role, is refused its own wait_for_events while a sub-agent of its session
// runs and no thread on its story awaits the designer (S-0285): the call
// reports work items and threads, not sub-agents, so it would run to its
// timeout. RecordChange keeps each session's running sub-agents, from a
// harness's reports of each starting and stopping, Claude Code's
// SubagentStart and SubagentStop hooks.
//
// A story's agent, the planner, the analyzer, and the curator write only as
// themselves (S-0369, ADR-0143): an entry's author is whatever --by, else FLAI_AGENT,
// says, and permission_prompt takes an answer by its author, so their own
// Bash call is refused when a flai command's --by or --placed-by names
// another, or when it sets, empties, or unsets FLAI_AGENT or FLAI_SESSION to
// anything but their own. Their own name is FLAI_AGENT as the hook sees it,
// or the name flai serve gives the run; the orchestrator keeps its own rule.
//
// No session flai serve starts, and no sub-agent, changes the manifest's
// shared paths, with shared_paths_edit or flai shared add or remove
// (S-0295, ADR-0096): they decide which overlaps hold a story. Only the
// operator's own session does; listing and checking them are reads.
//
// The manifest's sub-projects are the project's to shape (S-0491): a story's
// agent changes them, with subproject_edit or flai subproject add, set, and
// remove, on its story's branch; no sub-agent does, a task sub-agent among
// them, nor the planner or the analyzer, and the orchestrator only while
// orchestration.permissions.edit_subprojects is on (ADR-0195). Listing them,
// with subprojects or flai subproject list, is a read open to all.
//
// No session flai serve starts, and no sub-agent, runs flai forge push
// (S-0393): the forge host action has flai serve push a story's branch and
// open its pull request on the move to review, with the host's credentials
// and the forge's token, which no agent holds; flai forge push is the
// operator's. flai forge check and status read, and are left alone.
//
// No session flai serve starts, and no sub-agent, runs flai kluster up,
// down, deploy, teardown, or promote (S-0577, ADR-0254 §7): they change the
// project's kluster, shared infrastructure on the operator's machine, which
// is the operator's, and flai serve's behind the kluster host action. flai
// kluster status reads, and is left alone.
//
// No session flai serve starts, and no sub-agent, writes the host
// configuration, ~/.flai/config.json or the file FLAI_CONFIG names (S-0511,
// ADR-0220): flai serve enable and disable; flai serve agent set and clear,
// and harness and provider with a setting or --reset; flai serve checks set,
// timeout, and clear; components idle and max-servers; forge with
// --token-env or --reset; notify set and clear; import add and remove; flai
// config set; and flai template use are refused it, whatever its role or
// permissions, since it holds what the operator lets flai serve do on the
// host. A sub-agent puts the setting it needs in its final message, and a
// served session asks the operator on a thread. The reads, such as flai serve
// status and the show and list forms, are left alone by this rule.
package guard

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Guard decides on the calls of one flai: Commands are the names of its
// commands, so that a word flai on a command line counts as running flai
// only when a command of its follows. Role is the session's role, from
// FLAI_ROLE: RolePlan holds the session's own calls to planning,
// RoleOrchestrate holds them to Permissions, the project's
// orchestration.permissions, RoleAnalyze to its report, RoleCurate to the
// digest, each through its role file, a defined role to its role file, and
// any other leaves them alone. Opener says
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
//
// ReadOnly says whether a sub-agent of a type is read-only, as ReadOnlyAgents
// reads it from a project's definitions; such a sub-agent's shell command
// that writes a file outside the temporary folder is refused (S-0368,
// ADR-0141). A nil ReadOnly holds no sub-agent read-only.
//
// Agent, Session, and Item are the hook's FLAI_AGENT, FLAI_SESSION, and
// FLAI_ITEM. A story's agent, the planner, the analyzer, and the curator
// write only as Agent or the name flai serve gives the run, agent-<Story>,
// planner-<Item>, analyzer, or curator, and only in Session: their own call
// that names another author with --by or --placed-by, or sets, empties, or
// unsets FLAI_AGENT or FLAI_SESSION to anything but these, is refused
// (S-0369, ADR-0143). A defined role writes only as Agent or its run name,
// its file's name, or <name>-<Item> on an item, and only in Session, so held
// too (ADR-0210). The orchestrator writes to a thread, accepts, and enqueues
// only as orchestrator, under its role file's conditions (S-0220, ADR-0093).
//
// A Role that a project's role file defines, any but the four above and the
// story's agent's and its sub-agents' verbs, holds the session's own calls to
// RoleFile, the role file of that role, nil when none was read (S-0416,
// ADR-0207). RoleProblem says why the file could not be read, missing or not
// parsing, naming the file or its folder and flai check; while it is set, or
// RoleFile is nil, the session reads and is refused every other call.
// RolePlan, RoleOrchestrate, RoleAnalyze, and RoleCurate hold them to
// RoleFile too, the project's file of the role, or, while it is nil, flai's
// baseline file of it; RoleProblem fails them closed as it does a defined
// role. Every role is so held to its role file; the rules in Go of the four
// only word its refusals. RolePermissions are the manifest's
// roles.<name>.permissions of the role; the orchestrator's stay in
// Permissions. Subcommands are each flai command's
// subcommands, from the command tree, which a role file's command entries
// are matched against; nil when not given. Design is the design folder,
// relative to Root and slash separated, for {design} in the file's files;
// while it is "", the parent of Reports stands for it when Reports is given.
type Guard struct {
	Commands        []string
	Role            string
	Served          bool
	Permissions     manifest.Permissions
	Opener          func(id string) (string, error)
	Root            string
	Reports         string
	Story           string
	Running         func() ([]SubAgent, error)
	ThreadOpen      func() (bool, error)
	AutoApprove     bool
	ReadOnly        func(agentType string) bool
	Agent           string
	Session         string
	Item            string
	RoleFile        *roles.Role
	RoleProblem     string
	RolePermissions manifest.Permissions
	Subcommands     map[string][]string
	Design          string
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

// Reader reads one harness's hook input for the guard and answers the
// harness with the guard's verdict, so that Decide sees no harness.
type Reader interface {
	// Read reads a hook's input into what it asks of the guard: a Call, in
	// g's role and story, or a sub-agent starting or stopping.
	Read(data []byte, g Guard) (Hook, error)
	// Answer tells the harness v, writing a refusal's reason to stderr, and
	// returns the exit status the hook ends with.
	Answer(v Verdict, stderr io.Writer) int
}

// Hook is what one hook input asks of the guard: a decision on Call, and on
// each call of Also, which the same tool call makes too, such as each further
// file an OpenCode patch writes; or, when SubAgent is set, a change to the
// record of running sub-agents.
type Hook struct {
	Call     Call
	Also     []Call
	SubAgent *SubAgentChange
}

// Calls are the calls h asks a decision on: Call, then each of Also.
func (h Hook) Calls() []Call { return append([]Call{h.Call}, h.Also...) }

// DecideAll decides on each of calls under g's policy, as Decide does, and
// answers the first refusal, else the first Ask, else lets them through: a
// tool call that makes several is refused when any one is.
func DecideAll(g Guard, calls []Call) Verdict {
	var asked Verdict
	for _, c := range calls {
		switch v := Decide(g, c); v.Outcome {
		case Deny:
			return v
		case Ask:
			if asked.Outcome == Allow {
				asked = v
			}
		}
	}
	return asked
}

// readers are the harnesses' readers, by the name flai guard --harness
// takes.
var readers = map[string]Reader{ClaudeCode: claudeCode{}, OpenCode: openCode{}}

// ReaderFor is the reader of a harness's hook input, by the harness's name,
// and whether there is one.
func ReaderFor(harness string) (Reader, bool) {
	r, ok := readers[harness]
	return r, ok
}

// Harnesses are the names of the harnesses with a reader, sorted.
func Harnesses() []string { return slices.Sorted(maps.Keys(readers)) }

// rules are what one kind of caller may run on a command line: flai says
// why a flai command is refused, and the permission that would allow it, or
// "" when it is not refused, from its command, its subcommand, and the words
// after flai; git says why a git command other than a read is refused, from
// its subcommand, and nil refuses none.
type rules struct {
	flai func(cmd, sub string, rest []string) (why, needs string)
	git  func(cmd string) string
}

// Why a read-only sub-agent does not serve the component library, and what
// it does instead (S-0477, ADR-0186); and why no sub-agent writes its
// manifest key.
const (
	servesComponents = "flai components serve writes a Vite config under .flai-cache and, with --story, may install the worktree's dependencies, and a read-only sub-agent writes no file but under the temporary folder (ADR-0141, ADR-0186); ask the story's agent, or a task sub-agent, to run it"
	configsManifest  = "flai components config rewrites system-flow.yaml, which is the story's agent's to change (ADR-0060, ADR-0186)"
	// editsSubprojects says why no sub-agent adds, changes, or removes a
	// sub-project (S-0491).
	editsSubprojects = "flai subproject add, set, and remove rewrite system-flow.yaml and create a sub-project's folder, which are the story's agent's to change, with flai subproject or subproject_edit (ADR-0060)"
)

// subAgent are a sub-agent's rules, read-only or not: flai's reads, the flai
// adr commands that write an ADR without committing it (S-0287, ADR-0127),
// flai components serve unless it is read-only (S-0477, ADR-0186), and git's
// reads; flai components config and flai subproject add, set, and remove
// (S-0491), which rewrite the manifest, are refused it. Each refusal cites
// its ADRs.
func subAgent(readOnly bool) rules {
	return rules{
		flai: func(cmd, sub string, rest []string) (string, string) {
			switch {
			case reads(cmd, sub, rest):
				return "", ""
			case cmd == "adr" && slices.Contains(adrWrites, sub):
				if slices.ContainsFunc(adrCommits, func(f string) bool { return named(rest, f) }) {
					return commitsADR + "; run it without --commit or --autocommit " + leavesADRToIt, ""
				}
				return "", ""
			case cmd == "components" && sub == "serve":
				if readOnly {
					return servesComponents, ""
				}
				return "", ""
			case cmd == "components" && sub == "config":
				return configsManifest, ""
			case cmd == "subproject":
				return editsSubprojects, ""
			case cmd == "hostapi":
				return writesHostAPI(rest), ""
			}
			return "flai commands that change work items, threads, conversations between stories, narratives, or releases are the story's agent's (ADR-0060)", ""
		},
		git: func(string) string {
			return "git commands that change the worktree, the index, branches, or history are the story's agent's (ADR-0060)"
		},
	}
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
			needs, never := orchestrated(cmd, sub, rest, g.Permissions)
			switch {
			case never != "":
				return nevers(never), ""
			case needs != "" && !g.Permissions.Allows(needs):
				return off(needs), needs
			case needs == manifest.PermitAcceptReviews && !workitem.IsOrchestrator(value(rest, "--by")):
				if cmd == "queue" {
					return enqueuesAsSelf, needs
				}
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
// it whatever its permissions. Its permissions p pick the one cost of delay
// inputs need (costNeeds).
func orchestrated(cmd, sub string, rest []string, p manifest.Permissions) (needs, never string) {
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
			return costNeeds(positionals(rest, editValues)[1], p), ""
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
		// --batch is the merge queue's landing, flai serve's alone (S-0403)
		if named(rest, "--batch") {
			return "", enqueuesOnly
		}
		return manifest.PermitAcceptReviews, ""
	case "queue":
		// flai queue with no subcommand lists, let through above
		if slices.Contains(queueEnqueues, sub) {
			return manifest.PermitAcceptReviews, ""
		}
		return "", enqueuesOnly
	case "push", "release":
		// flai release --evaluate is a read, let through above
		return "", publishesThrough
	case "subproject":
		// flai subproject list, and flai subproject alone, are reads, let
		// through above: add, set, and remove remain (ADR-0195)
		return manifest.PermitEditSubprojects, ""
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
func finalizesDraft(in Args) bool {
	return !in.Draft && slices.Contains(in.Fields, "draft") &&
		!slices.ContainsFunc(in.Fields, func(f string) bool { return !slices.Contains(finalizeFields, f) })
}

// setsCostInputs says whether item_edit's input gives an epic or a story cost
// of delay inputs and changes nothing else (S-0328, ADR-0119, I-0187):
// cost_of_delay with one input or more, none null or an empty string, which
// removes it, and no other key, value among them; and no field but the item,
// its hash, and its project besides. Whether the item may have them is
// flai's to judge.
func setsCostInputs(in Args) bool {
	if !isEpic(in.ID) && !isStory(in.ID) || len(in.CostOfDelay) == 0 ||
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

// costNeeds is the permission the orchestrator needs to give item id, an
// epic or a story, cost of delay inputs, under its permissions p: an epic
// needs set_cost_of_delay; a story needs either, so set_cost_of_delay while
// it is on and plan_backlog_stories, which a refusal names, otherwise
// (S-0523, ADR-0224).
func costNeeds(id string, p manifest.Permissions) string {
	if isEpic(id) || p.Allows(manifest.PermitSetCostOfDelay) {
		return manifest.PermitSetCostOfDelay
	}
	return manifest.PermitPlanBacklogStories
}

// setsCostInputsOnly says whether the words of flai edit give an epic or a
// story cost of delay inputs and change nothing else (S-0328, ADR-0119,
// I-0187): one input flag or more, each with a value that is not empty,
// which would remove it, beside the flags that finalizing a draft allows but
// --no-draft.
func setsCostInputsOnly(words []string) bool {
	if args := positionals(words, editValues); len(args) < 2 || !isEpic(args[1]) && !isStory(args[1]) {
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

// ProtectedWrite says whether c writes a path Claude Code protects other than
// one in .git, such as a file in a .claude/ folder or .mcp.json, and which
// file (S-0299, ADR-0106): Claude Code asks a person before any such write,
// so that a sub-agent's would wait on a thread for the operator.
// permission_prompt refuses a path in .git at once, so that one never waits.
func ProtectedWrite(c Call) (string, bool) { return protectedWrite(c) }
