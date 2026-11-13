package guard

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// RolePlan is the planner's role, as flai serve sets it in FLAI_ROLE.
const RolePlan = "plan"

// Backlog is the one state the planner may move an item to.
const Backlog = "backlog"

// MCPPlans are flai's MCP tools the planner's baseline role file lists
// besides MCPReads, which word its refusals; item_move it may call only to
// move an item to backlog.
var MCPPlans = []string{"activity_log", "inbox", "item_edit", "item_new", "thread_open", "thread_reply", "wait_for_events"}

// cliPlans are the flai commands the planner's baseline role file lists
// besides cliReads, in the form cliReads has, which word its refusals; flai
// move it may run only to move an item to backlog.
var cliPlans = map[string][]string{
	"edit":    nil,
	"epic":    {"new"},
	"issue":   {"new", "bump"},
	"story":   {"new"},
	"task":    {"new"},
	"thread":  {"new", "reply"},
	"touches": nil,
}

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
var MCPOrchestrates = []string{"activity_log", "inbox", "thread_open", "verify", "wait_for_events"}

// cliOrchestrates are the flai commands the orchestrator may run besides
// cliReads, whatever its permissions, each with the subcommands it may run:
// flai queue with none lists the merge queue (S-0403).
var cliOrchestrates = map[string][]string{
	"issue":  {"new", "bump"},
	"queue":  {""},
	"thread": {"new"},
}

// QueueAdd is the MCP tool through which the orchestrator enqueues a story it
// has verified for the merge queue to land, while
// orchestration.permissions.accept_reviews is on (S-0403).
const QueueAdd = "queue_add"

// Accept is the MCP tool through which the orchestrator accepts a story it
// has verified, running flai accept --by orchestrator on the host with the
// host's configuration, while orchestration.permissions.accept_reviews is
// on (ADR-0278).
const Accept = "accept"

// orchestratorAccepts says why no session but the orchestrator calls Accept
// (ADR-0278).
const orchestratorAccepts = "it is the orchestrator's alone, while orchestration.permissions.accept_reviews is on; the operator accepts a story with flai accept <id> --by <name> on the host, or from the board (ADR-0093)"

// queueEnqueues are the subcommands of flai queue the orchestrator runs while
// accept_reviews is on, as itself; it runs no other that writes (S-0403).
var queueEnqueues = []string{"add"}

// ReleasePublish is the MCP tool through which the orchestrator alone
// publishes, while orchestration.permissions.publish is on (S-0222).
const ReleasePublish = "release_publish"

// gitPublishes are the git commands that reach the remote, which the
// orchestrator leaves to release_publish; git tag in a form that only lists
// is a read (gitListForms).
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

// MCPAnalyzes are flai's MCP tools the analyzer's baseline role file lists
// besides MCPReads, which word its refusals: issue_new and issue_bump file
// and bump an issue for a finding of its report (S-0224).
var MCPAnalyzes = []string{"activity_log", "inbox", "issue_bump", "issue_new", "thread_open", "thread_reply", "wait_for_events"}

// itemWrites are flai's MCP tools that write a work item, issue_story's
// draft story among them, which the analyzer never calls.
var itemWrites = []string{"issue_story", "item_edit", "item_move", "item_new"}

// cliAnalyzes are the flai commands the analyzer's baseline role file lists
// besides cliReads, each with its subcommands, which word its refusals: it
// files and bumps the issues its findings call for (S-0224).
var cliAnalyzes = map[string][]string{
	"issue": {"new", "bump"},
}

// cliAuthors are the flai commands that author a work item, with the
// subcommands that do, nil for every one; the analyzer runs none of them.
// flai formula run makes a story and its tasks (S-0421).
var cliAuthors = map[string][]string{
	"epic":    nil,
	"formula": {"run"},
	"issue":   {"story"},
	"story":   nil,
	"task":    nil,
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

// RoleCurate is the curator's role, as flai serve sets it in FLAI_ROLE
// (S-0431, ADR-0196).
const RoleCurate = "curate"

// MCPCurates are flai's MCP tools the curator's baseline role file lists
// besides MCPReads, which word its refusals.
var MCPCurates = []string{"activity_log", "inbox", "thread_open", "thread_reply", "wait_for_events"}

// cliCurates are the flai commands the curator's baseline role file lists
// besides cliReads, each with its subcommands, which word its refusals:
// flai digest harvest lists the candidates, new and edit write an entry or
// merge a candidate into the entry that holds it, absorb marks an entry a
// document records, cited writes each entry's citers, and compact compacts
// the entries uncited within the window (ADR-0199).
var cliCurates = map[string][]string{
	"digest": {"absorb", "cited", "compact", "edit", "harvest", "new"},
}

// curator ends each of the curator's refusals: the rule the call breaks and
// what to do instead.
const curator = "the curator writes digest entries only with flai digest new, edit, absorb, cited, and compact, and edits no file and authors no work item (strategic-agents.md, ADR-0196, ADR-0199). Ask the operator with thread_open when an entry needs their judgement."

// seedsUnjudged says why the curator never runs flai digest seed.
const seedsUnjudged = "flai digest seed writes an entry for every candidate without judging whether an entry already holds it"

// curating are the curator's rules: flai's reads, digest harvest, new,
// edit, absorb, cited, and compact, and git's reads.
var curating = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		switch {
		case reads(cmd, sub, rest):
			return "", ""
		case cmd == "digest" && sub == "seed":
			return seedsUnjudged, ""
		}
		if subs, ok := cliCurates[cmd]; ok && slices.Contains(subs, sub) {
			return "", ""
		}
		return "of flai's commands it runs only those that read, and digest harvest, new, edit, absorb, cited, and compact", ""
	},
	git: func(string) string { return "it runs only git's reads" },
}

// flaiValues are the flags that take a value and that flai reads before its
// command: its own, and --wait, which flai kluster deploy and teardown --job
// take and cobra reads wherever it stands, so that flai --wait 5m kluster
// deploy is read as the deploy it runs.
var flaiValues = map[string]bool{"--config": true, "--wait": true}

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

// costFields are the fields item_edit's input may give when it only gives an
// epic or a story cost of delay inputs, and costInputs the keys of its
// cost_of_delay that are inputs (S-0328, ADR-0119).
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
// it only finalizes a draft, or gives an epic or a story cost of delay inputs
// (S-0219, S-0328, ADR-0119, I-0187).
const (
	editsMCP = "it edits an item only to finalize a draft, with item_edit draft false and nothing else, or to give an epic or a story cost of delay inputs, with cost_of_delay's revenue_per_week, penalty_per_week, and time_lost_per_cycle and nothing else: no value, no clear_cost_of_delay"
	editsCLI = "it edits an item only to finalize a draft, with flai edit --no-draft and nothing else, or to give an epic or a story cost of delay inputs, with --revenue-per-week, --penalty-per-week, and --time-lost-per-cycle and nothing else: no --cost-of-delay-value, no --clear-cost-of-delay"
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
const acceptsWithFlai = "item_move never moves an item to done; it accepts a story with flai accept <id> --by orchestrator --verified <commit> --evidence <file>, while orchestration.permissions.accept_reviews is on (ADR-0093), run on the host by the MCP tool accept with story, verified, and evidence (ADR-0278), or, while the merge queue runs, enqueues it with queue_add for the queue to land (S-0403)"

// enqueuesAsSelf says why the orchestrator enqueues a story only as itself,
// whatever accept_reviews lets it do (S-0403).
const enqueuesAsSelf = "it enqueues a story only as itself, so give --by orchestrator with --verified and --evidence, or call queue_add: orchestration.permissions.accept_reviews lets it enqueue as nobody else (ADR-0093)"

// enqueuesOnly says why the orchestrator never lands, takes out, holds, or
// resumes anything in the merge queue by hand (S-0403, ADR-0176).
const enqueuesOnly = "the merge queue lands, bisects, and sends back the stories enqueued itself, in flai serve, and taking one out, holding, and resuming the queue are the operator's; the orchestrator enqueues a story it has verified with queue_add, or flai queue add <id> --by orchestrator --verified <commit> --evidence <file>, while orchestration.permissions.accept_reviews is on"

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

// forgePushes are the subcommands of flai forge that push a story and open
// its pull request, the operator's alone (S-0393).
var forgePushes = []string{"push"}

// operatorsPush says why no agent runs flai forge push (S-0393).
const operatorsPush = "the forge host action has flai serve push a story's branch and open or update its pull request when the story moves to review, with the host's credentials and the forge's token, which no agent holds; flai forge push is the operator's (S-0393)"

// askForPush ends the refusal of a session's own flai forge push: what it
// does instead.
const askForPush = "Move the story to review and leave the push to flai serve, or ask the operator on a thread with thread_open; flai forge status and check read and are open to you."

// forging are the rules that find flai forge push on a command line, and no
// git command.
var forging = rules{
	flai: func(cmd, sub string, _ []string) (string, string) {
		if cmd == "forge" && slices.Contains(forgePushes, sub) {
			return operatorsPush, ""
		}
		return "", ""
	},
}

// klusterWrites are the subcommands of flai kluster that change the
// project's cluster or a namespace in it, the operator's alone (S-0577,
// ADR-0254 §7). flai kluster status reads and is open to every agent.
var klusterWrites = []string{"up", "down", "deploy", "teardown", "promote"}

// operatorsKluster says why no agent runs a flai kluster command that writes
// (S-0577, ADR-0254 §7).
const operatorsKluster = "flai kluster up, down, deploy, teardown, and promote change the project's kluster, shared infrastructure on the operator's machine, so they are the operator's, from a shell or the dashboard's Ops pages, and flai serve's behind the kluster host action (ADR-0254, safety.md)"

// askForKluster ends the refusal of a session's own flai kluster command
// that writes: what it does instead.
const askForKluster = "Ask the operator on a thread with thread_open; flai kluster status reads and is open to you."

// klusterRequests are the subcommands of flai kluster that, with --job, ask
// flai serve for a job on a story's namespace instead of running it: a
// request a story's agent makes for its own story (ADR-0292 §1, §2).
var klusterRequests = []string{"deploy", "teardown"}

// KlusterDeploy and KlusterTeardown are the MCP tools that make the requests
// of flai kluster deploy and teardown --job (ADR-0292 §1).
const (
	KlusterDeploy   = "kluster_deploy"
	KlusterTeardown = "kluster_teardown"
)

// klusterValues are the flags of flai and of a kluster request that take a
// value, so that the value is not read as the story it names.
var klusterValues = map[string]bool{"--config": true, "--wait": true}

// Why the guard refuses a request a story's agent alone makes, and what to do
// instead (ADR-0292 §1, §2).
const (
	storysRequest = "a request to deploy or tear down a story's namespace, flai kluster deploy or teardown with --job or the MCP tools kluster_deploy and kluster_teardown, is the story's agent's, for its own story alone (ADR-0292)"
	askOfStory    = "Put the deploy or teardown you need in your final message; the story's agent asks for it."
)

// storyAgent is a name flai serve gives a story's agent, agent-S-nnnn, with
// the story.
var storyAgent = regexp.MustCompile(`^agent-(S-\d+)$`)

// shellExpands are the characters with which the shell could make a word of
// a command line into other words, such as $X, `...`, a brace, a glob, a
// tilde, or an escape: on a line that has one, the guard cannot tell the
// words flai is given, which could hold --job=false.
const shellExpands = "$`{}*?[]~\\"

// writtenPlainly says why a request is let through only as written plainly.
const writtenPlainly = "a request is let through only as written plainly, with no $, `, brace, glob, ~, or backslash on the line, since the shell could expand one into --job=false and run the job itself (ADR-0292)"

// ownStory is the story of c's session, when it is a story's agent's, in
// canonical form: FLAI_STORY, else the story in its name, agent-S-nnnn; ""
// for a session with a role or with neither.
func (g Guard) ownStory(c Call) string {
	if c.Role != "" {
		return ""
	}
	if c.Story != "" {
		return workitem.CanonicalID(c.Story)
	}
	if m := storyAgent.FindStringSubmatch(g.Agent); m != nil {
		return workitem.CanonicalID(m[1])
	}
	return ""
}

// klusterRequest says whether the words after flai, rest, whose subcommand
// of kluster is sub, ask flai serve for a job rather than run one: deploy or
// teardown with --job set true, the last of its forms deciding as flags do,
// and before any --. story is the one story they name, canonical, or "" when
// they name none or more than one.
func klusterRequest(sub string, rest []string) (story string, job bool) {
	if !slices.Contains(klusterRequests, sub) {
		return "", false
	}
	var args []string
	for i := 0; i < len(rest); i++ {
		w := rest[i]
		switch {
		case w == "--":
			args = append(args, rest[i+1:]...)
			i = len(rest)
		case w == "--job":
			job = true
		case strings.HasPrefix(w, "--job="):
			v, err := strconv.ParseBool(strings.TrimPrefix(w, "--job="))
			job = err == nil && v
		case strings.HasPrefix(w, "-"):
			if klusterValues[w] {
				i++
			}
		default:
			args = append(args, w)
		}
	}
	// the words left are kluster, the subcommand, and what they name
	if named := args[min(2, len(args)):]; len(named) == 1 && isStory(named[0]) {
		story = workitem.CanonicalID(named[0])
	}
	return story, job
}

// klusterBools are flai's own flags that take no value, the ones flai knows
// between kluster and its subcommand; a subcommand's own, such as --job, it
// does not know there.
var klusterBools = []string{"--help", "--json", "--verbose", "--yes", "-h", "-v", "-y"}

// klusterWriteSub is the subcommand of flai kluster that writes which the
// words after flai, rest, may run, or "" when they run none: given, the
// first word after kluster's flags, or the one flai would take instead, as
// it reads a flag it does not know before the subcommand as taking the word
// after it. ambiguous says the two differ, so that the line makes no
// request.
func klusterWriteSub(given string, rest []string) (sub string, ambiguous bool) {
	at := slices.Index(rest, "kluster")
	var read string
	for i := at + 1; at >= 0 && i < len(rest); i++ {
		w := rest[i]
		if w == "--" {
			break
		}
		if strings.HasPrefix(w, "-") {
			if !strings.Contains(w, "=") && !slices.Contains(klusterBools, w) {
				i++
			}
			continue
		}
		read = w
		break
	}
	switch {
	case slices.Contains(klusterWrites, given):
		return given, read != given
	case slices.Contains(klusterWrites, read):
		return read, true
	}
	return "", false
}

// klustering are the rules that find on c's command line a flai kluster
// command that writes and that c may not run, and no git command: every one
// for a sub-agent, and for a session flai serve started but a story's agent's
// request for its own story, written plainly; for the operator's own
// session only a request for a story other than the one in its agent's name.
// Each says why and what to do.
func (g Guard) klustering(c Call) rules {
	own := g.ownStory(c)
	return rules{
		flai: func(cmd, given string, rest []string) (string, string) {
			if cmd != "kluster" {
				return "", ""
			}
			sub, ambiguous := klusterWriteSub(given, rest)
			if sub == "" {
				return "", ""
			}
			story, job := klusterRequest(sub, rest)
			job = job && !ambiguous
			switch {
			case c.Subagent && job:
				return storysRequest + ". " + askOfStory, ""
			case c.Subagent:
				return operatorsKluster + ". Put what you need done to the kluster in your final message.", ""
			case job && own != "" && strings.ContainsAny(c.Command, shellExpands):
				return fmt.Sprintf("%s. Write it as flai kluster %s %s --job, in a Bash call of its own, or call kluster_%s.", writtenPlainly, sub, own, sub), ""
			case job && own != "" && story == own:
				return "", ""
			case job && own != "":
				return fmt.Sprintf("a story's agent asks flai serve to deploy or tear down its own story's namespace alone, %s, and this names %s (ADR-0292). Ask for your own story's: flai kluster %s %s --job, or kluster_%s with no story.", own, cmp.Or(story, "no single story"), sub, own, sub), ""
			case g.served(c):
				return operatorsKluster + ". " + askForKluster, ""
			}
			return "", ""
		},
	}
}

// forgePush decides on a command line that runs flai forge push, by a
// sub-agent or in a session flai serve started: the zero Verdict when it runs
// none or the operator's own session runs it. The orchestrator's refusal
// names the call, so that its activity document logs it.
func (g Guard) forgePush(c Call) Verdict {
	return g.operatorsOnly(c, forging, operatorsPush, askForPush, "Put what you need pushed in your final message.")
}

// klusterWrite decides on a call that changes the project's kluster (S-0577)
// or asks flai serve to (ADR-0292): on a command line, as klustering finds
// it, and on the MCP tools kluster_deploy and kluster_teardown, which every
// sub-agent is refused, and a story's agent when they name a story other
// than its own. The zero Verdict for any other call, and for one the
// operator's own session makes; the orchestrator's refusal names the call, so
// that its activity document logs it.
func (g Guard) klusterWrite(c Call) Verdict {
	switch c.Kind {
	case FlaiTool:
		return g.klusterTool(c)
	case Shell:
	default:
		return Verdict{}
	}
	r := g.klustering(c)
	for _, words := range arguments(c.Command) {
		why, _ := g.refuse(words, r)
		if why == "" {
			continue
		}
		call := strings.Join(words, " ")
		v := deny(fmt.Sprintf("%s cannot run %q: %s", g.klusterCaller(c), call, why))
		if c.Role == RoleOrchestrate && !c.Subagent {
			v.Call = call
		}
		return v
	}
	return Verdict{}
}

// klusterTool decides on a call of kluster_deploy or kluster_teardown: every
// sub-agent is refused it, and a story's agent when it names a story other
// than its own; with no story the tool asks for the session's own.
func (g Guard) klusterTool(c Call) Verdict {
	if c.Tool != KlusterDeploy && c.Tool != KlusterTeardown {
		return Verdict{}
	}
	if c.Subagent {
		return deny(fmt.Sprintf("a sub-agent (%s) cannot call %s: %s. %s", cmp.Or(c.Agent, "unnamed"), c.Tool, storysRequest, askOfStory))
	}
	own := g.ownStory(c)
	if own == "" || c.Args.Story == "" || workitem.CanonicalID(c.Args.Story) == own {
		return Verdict{}
	}
	return deny(fmt.Sprintf("the story's agent cannot call %s for %s: a story's agent asks flai serve to deploy or tear down its own story's namespace alone, %s (ADR-0292). Call %s with no story, or with story %s.", c.Tool, c.Args.Story, own, c.Tool, own))
}

// klusterCaller names c's caller in a refusal of a flai kluster command.
func (g Guard) klusterCaller(c Call) string {
	switch {
	case c.Subagent:
		return fmt.Sprintf("a sub-agent (%s)", cmp.Or(c.Agent, "unnamed"))
	case g.ownStory(c) != "":
		return "the story's agent"
	}
	return g.who(c)
}

// operatorsOnly decides on a command line that runs a command found by
// found, the operator's alone, by a sub-agent or in a session flai serve
// started: refused with why and, for a sub-agent, subAsk, or for a served
// session, ask; the zero Verdict when it runs none or the operator's own
// session runs it. The orchestrator's refusal names the call.
func (g Guard) operatorsOnly(c Call, found rules, why, ask, subAsk string) Verdict {
	if c.Kind != Shell || !c.Subagent && !g.served(c) {
		return Verdict{}
	}
	for _, words := range arguments(c.Command) {
		if refused, _ := g.refuse(words, found); refused == "" {
			continue
		}
		call := strings.Join(words, " ")
		if c.Subagent {
			return deny(fmt.Sprintf("a sub-agent (%s) cannot run %q: %s. %s", cmp.Or(c.Agent, "unnamed"), call, why, subAsk))
		}
		v := deny(fmt.Sprintf("%s cannot run %q: %s. %s", g.who(c), call, why, ask))
		if c.Role == RoleOrchestrate {
			v.Call = call
		}
		return v
	}
	return Verdict{}
}
