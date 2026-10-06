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
// column by a policy). A refusal names that permission; a call no permission
// allows, such as an edit of a file, a commit, or a story placed by hand in
// the pull order, it never makes. Its sub-agents are held as every sub-agent
// is.
package guard

import (
	"encoding/json"
	"fmt"
	"maps"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Event is the part of a Claude Code PreToolUse hook's input the guard
// reads.
type Event struct {
	ToolName  string `json:"tool_name"`
	ToolInput Input  `json:"tool_input"`
	// AgentID is set only when a sub-agent makes the call; AgentType names
	// the sub-agent's definition. A session started with --agent may carry
	// an agent type of its own, so the ID is what marks a sub-agent.
	AgentID   string `json:"agent_id"`
	AgentType string `json:"agent_type"`
}

// Input is the part of a tool call's input the guard reads: Bash's command;
// item_move's item and the state it moves the item to, and plan's item;
// item_new's type and whether it makes a story a draft, and item_edit's
// draft; and the file Edit or Write changes. Fields are the names of every
// field the input gives, whatever the guard reads of it, so that an
// item_edit that finalizes a draft is told from one that changes more, and
// a draft false given from one left out.
type Input struct {
	Command  string   `json:"command"`
	ID       string   `json:"id"`
	To       string   `json:"to"`
	Type     string   `json:"type"`
	Draft    bool     `json:"draft"`
	FilePath string   `json:"file_path"`
	Fields   []string `json:"-"`
}

// UnmarshalJSON decodes a tool call's input and the names of its fields.
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
	return nil
}

// MCPPrefix starts the name of each of flai's MCP tools as Claude Code
// knows them, from a server named flai.
const MCPPrefix = "mcp__flai__"

// MCPReads are flai's MCP tools a sub-agent may call: they read and use no
// agent's identity.
var MCPReads = []string{"board", "doc_get", "doc_search", "item_get", "order_by_policy", "prime", "promote_candidates", "release_evaluate", "thread_get", "who_touches"}

// cliReads are the flai commands a sub-agent may run, each with the
// subcommands it may run; nil allows the command whatever follows it, and ""
// allows it with no subcommand.
var cliReads = map[string][]string{
	"board":    {""},
	"check":    nil,
	"cod":      nil,
	"doc":      {"search", "show"},
	"forecast": nil,
	"help":     nil,
	"issue":    {"list"},
	"prime":    nil,
	"show":     nil,
	"stats":    nil,
	"stream":   {"diff"},
	"thread":   {"list", "show"},
	"touches":  {"suggest"},
	"version":  nil,
}

// flagReads are the flai commands a sub-agent may run only in the form that
// reads (S-0217, S-0219): with one of the flags that make them read, without
// any flag that makes them write, and in the form form allows when it is
// set. flai order --by computes an order and --apply writes it; flai order
// that names a story places it; flai release without --evaluate releases;
// and flai plan without --candidates starts the planner, which with it takes
// no item.
var flagReads = map[string]struct {
	with    []string
	without []string
	form    func(words []string) bool
}{
	"order":   {with: []string{"--by"}, without: []string{"--apply"}, form: ordersByPolicy},
	"plan":    {with: []string{"--candidates"}},
	"promote": {with: []string{"--candidates", "--drafts"}},
	"release": {with: []string{"--evaluate"}},
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

// fileEdits are the tools that change files, which the planner never uses.
var fileEdits = []string{"Edit", "NotebookEdit", "Write"}

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

// askOperator ends each of the orchestrator's refusals.
const askOperator = "Ask the operator with thread_open on the item if it needs doing (strategic-agents.md, ADR-0060)."

// plansEpics says why the orchestrator asks for the planner on an epic
// alone.
const plansEpics = "it asks for the planner on an epic alone, one flai plan --candidates lists"

// byHand says why the orchestrator never places a story by hand in the pull
// order (S-0219).
const byHand = "it orders the ready column by its policy, with flai order --by <policy> --apply, and never places a story by hand, which is the operator's"

// Guard decides on the calls of one flai: Commands are the names of its
// commands, so that a word flai on a command line counts as running flai
// only when a command of its follows. Role is the session's role, from
// FLAI_ROLE: RolePlan holds the session's own calls to planning,
// RoleOrchestrate holds them to Permissions, the project's
// orchestration.permissions, and any other leaves them alone.
type Guard struct {
	Commands    []string
	Role        string
	Permissions manifest.Permissions
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
// after flai; git says why a git command other than a read is refused.
type rules struct {
	flai func(cmd, sub string, rest []string) (why, needs string)
	git  string
}

// subAgent are a sub-agent's rules: flai's reads and git's.
var subAgent = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		if reads(cmd, sub, rest) {
			return "", ""
		}
		return "flai commands that change work items, threads, narratives, or releases are the story's agent's", ""
	},
	git: "git commands that change the worktree, the index, branches, or history are the story's agent's",
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
	git: "it runs only git's reads",
}

// orchestration are the orchestrator's rules under its permissions: flai's
// reads, the commands it always runs, and those its permissions allow; and
// git's reads.
func (g Guard) orchestration() rules {
	return rules{
		flai: func(cmd, sub string, rest []string) (string, string) {
			needs, never := orchestrated(cmd, sub, rest)
			switch {
			case never != "":
				return never, ""
			case needs != "" && !g.Permissions.Allows(needs):
				return "it needs " + needs, needs
			}
			return "", ""
		},
		git: "it runs only git's reads",
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
		if args := positionals(rest, flaiValues); len(args) > 1 && isEpic(args[1]) {
			return manifest.PermitPlanBacklogEpics, ""
		}
		return "", plansEpics
	case "edit":
		if finalizesOnly(rest) {
			return manifest.PermitFinalizeDrafts, ""
		}
		return "", "it edits an item only to finalize a draft, with flai edit --no-draft and nothing else"
	case "move":
		if args := positionals(rest, moveValues); len(args) > 2 && args[2] == Ready {
			return manifest.PermitPromoteToReady, ""
		}
		return "", "it moves an item to ready and no further"
	case "order":
		if ordersByPolicy(rest) {
			return manifest.PermitOrderReady, ""
		}
		return "", byHand
	case "accept":
		return manifest.PermitAcceptReviews, ""
	case "push":
		return manifest.PermitPublish, ""
	case "release":
		if given(rest, "--pending") {
			return manifest.PermitPublish, ""
		}
		return "", "it releases with flai release --pending alone"
	case "thread":
		if sub == "reply" {
			return manifest.PermitAnswerThreads, ""
		}
	}
	return "", "of flai's commands that write, it runs only thread new, issue new and bump, and those its permissions allow"
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

// isEpic says whether id is an epic's, in any zero padding.
func isEpic(id string) bool {
	n, ok := strings.CutPrefix(strings.ToUpper(id), "E-")
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

// reads says whether a flai command with its subcommand and the words after
// flai only reads.
func reads(cmd, sub string, rest []string) bool {
	if f, ok := flagReads[cmd]; ok {
		return slices.ContainsFunc(f.with, func(w string) bool { return given(rest, w) }) &&
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
	if e.AgentID == "" {
		switch g.Role {
		case RolePlan:
			return Refusal{Why: g.plan(e)}
		case RoleOrchestrate:
			return g.orchestrate(e)
		}
		return Refusal{}
	}
	who := e.AgentType
	if who == "" {
		who = "unnamed"
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

// orchestrate decides on the orchestrator's own call under its permissions.
func (g Guard) orchestrate(e Event) Refusal {
	if slices.Contains(fileEdits, e.ToolName) {
		call := strings.TrimSpace(e.ToolName + " " + e.ToolInput.FilePath)
		return refusedOrchestrator(call, "use "+e.ToolName, "it acts through flai and never edits a file", "")
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		in := e.ToolInput
		call, what := strings.Join(strings.Fields(tool+" "+in.ID+" "+in.To), " "), "call "+tool
		var needs, never string
		switch {
		case slices.Contains(MCPReads, tool), slices.Contains(MCPOrchestrates, tool):
			return Refusal{}
		case tool == "plan":
			what = "plan " + in.ID
			if isEpic(in.ID) {
				needs = manifest.PermitPlanBacklogEpics
			} else {
				never = plansEpics
			}
		case tool == "item_edit":
			what = "edit " + in.ID
			if finalizesDraft(in) {
				what, needs = "finalize "+in.ID, manifest.PermitFinalizeDrafts
			} else {
				never = "it edits an item only to finalize a draft, with item_edit draft false and nothing else"
			}
		case tool == "item_move":
			what = fmt.Sprintf("move %s to %s", in.ID, in.To)
			if in.To == Ready {
				needs = manifest.PermitPromoteToReady
			} else {
				never = "it moves an item to ready and no further"
			}
		case tool == "thread_reply":
			needs = manifest.PermitAnswerThreads
		default:
			never = "of flai's tools that write, it calls only thread_open and activity_log, and those its permissions allow"
		}
		if never == "" && g.Permissions.Allows(needs) {
			return Refusal{}
		}
		return refusedOrchestrator(call, what, never, needs)
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

// refusedOrchestrator is the orchestrator's refusal of call, said as what:
// for want of the permission needs, or, when needs is "", because it never
// makes the call, for why.
func refusedOrchestrator(call, what, why, needs string) Refusal {
	r := Refusal{Call: call}
	if needs != "" {
		r.Needs = "orchestration.permissions." + needs
		r.Why = fmt.Sprintf("the orchestrator cannot %s: it needs %s, which is off. %s", what, r.Needs, askOperator)
		return r
	}
	r.Why = fmt.Sprintf("the orchestrator cannot %s: the orchestrator never does it, whatever its permissions: %s. %s", what, why, askOperator)
	return r
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
			if !subcommand.MatchString(cmd) || slices.Contains(gitReads, cmd) {
				continue
			}
			return r.git, ""
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

// commands splits a shell command line into simple commands, each as its
// words with quotes removed: at ;, &, |, newlines, and the openings of
// subshells and command substitutions. It is not a shell; it is enough to
// find each program a line runs and the words after it.
func commands(line string) [][]string {
	var out [][]string
	var words []string
	var word strings.Builder
	inWord := false
	var quote rune
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
	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			word.WriteRune(r)
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ';' || r == '&' || r == '|' || r == '\n' || r == '(' || r == ')' || r == '`' || r == '{' || r == '}':
			split()
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
