// Package harness turns a story's agent into a command to start (S-0104).
//
// A story says which harness works it, with which model, and a few options
// (its agent, S-0103). The operator says, on the host, which program each
// harness is and what an agent it starts may do. An adapter puts the two
// together into an argument list, which flai serve runs as it stands, never
// through a shell.
//
// What the story says is checked here before it goes anywhere: anyone who
// can edit a story in the dashboard can write it, so a story names tunables
// the adapter knows, never flags, programs, or permissions.
package harness

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Command is the harness that runs the operator's own command (S-0079).
const Command = "command"

// Request is what an agent is started for: a story, or, for the planner, an
// epic or a story to plan (S-0208), or, for the orchestrator, the whole
// project (S-0218), or, for the analyzer, the whole project with a focus
// or none (S-0223), or, for the curator, the whole project since its last
// run or from the start (S-0431), or, for a role a project defines, its
// item or none, with the run's inputs (S-0417).
type Request struct {
	Story   string          // the story's ID, checked by the caller; empty for a strategic agent
	Root    string          // the project's directory, where it runs
	Project string          // the project's key, for a session's name
	Agent   *manifest.Agent // the story's agent, or the strategic agent's; nil when it has none
	Name    string          // FLAI_AGENT the session works under
	Flai    string          // this flai's executable, the agent's MCP server
	// Role is conventions.RolePlan when the agent is the planner,
	// conventions.RoleOrchestrate when it is the orchestrator,
	// conventions.RoleAnalyze when it is the analyzer,
	// conventions.RoleCurate when it is the curator, the verb its file's
	// role names when it is a role the project defines, and empty when it
	// works a story.
	Role string
	// RoleFile is the role file of a role the project defines (ADR-0207),
	// when the agent runs as one: its name, its verb, the items it takes,
	// its harness definition, and its body (S-0417). Name is then RoleAgent
	// of its name and Item. Nil for the built-in roles, which run as their
	// own definitions, and for an agent working a story. Never read from a
	// start request.
	RoleFile *roles.Role `json:"-"`
	// Builtin is the role file a built-in role runs from, the project's or
	// the baseline flai carries (roles.Builtin): its body, the project's
	// additions below the marker included, is the prompt of the definition
	// claude-code passes it (S-0419, ADR-0207 §4). Nil for any other agent,
	// and for a built-in role whose caller gives none, whose definition keeps
	// the prompt of its .claude/agents file. Never read from a start request.
	Builtin *roles.Role `json:"-"`
	// Inputs are a defined role's run's inputs, in the order given, which
	// its prompt names (S-0417); none for any other agent. Never read from a
	// start request.
	Inputs []Input `json:"-"`
	// Item is the epic or story the planner plans, or a defined role runs
	// on, checked by the caller; empty when the agent works a story,
	// orchestrates, analyzes, or curates, or a defined role runs on none.
	Item string
	// AskedOn is the thread the operator wrote on since Item was planned
	// before, when the orchestrator asked for the planner on Item again for
	// it (S-0524, ADR-0222); empty for any other run.
	AskedOn string
	// Focus is what the analyzer looks for, one of Focuses; empty for all
	// of them, and for any other agent.
	Focus string
	// Since is when the curator's run takes in from, a YYYY-MM-DD day in UTC
	// or an RFC 3339 time: the since the operator gave, or else the time of
	// the curator's last run (S-0431, ADR-0196). It is empty on the
	// curator's first run, which takes in everything, and for any other
	// agent.
	Since string
	// Session names the harness's session, so that it can be resumed; a
	// harness that has no sessions ignores it.
	Session string
	// Answered is the thread the agent asked on that has been answered, or
	// the conversation (MS-nnnn) another story's agent wrote on to its story:
	// the agent ended while it waited, and is started again to go on (S-0104,
	// S-0335).
	Answered string
	// SentBack is what flai sent the story back from review to in-progress
	// for, on its pull request, when flai serve started its agent for that
	// (S-0394, ADR-0167); nil otherwise.
	SentBack *SentBack
	// Restart says how the story's last agent ended, when the operator has
	// had a new one started for it (S-0116); empty when it entered ready.
	Restart string
	// AutoRestart says which automatic restart this is, such as "1 of 2",
	// when flai serve, not the operator, started the agent again after the
	// last one ended with its story in progress (S-0294, ADR-0108); empty
	// otherwise.
	AutoRestart string
	// AfterLimit is set when flai serve started the agent again on its own
	// because its last one ended on the harness's usage limit and the host's
	// budget state no longer holds its provider, spending no automatic
	// restart (S-0534, ADR-0209 §5); Restart says how the last one ended.
	AfterLimit bool
	// Commit is the story's worktree, when the operator has had an agent
	// started to commit what it holds and nothing else (S-0140); empty
	// otherwise.
	Commit string
	// Started is set when the operator has had the agent started now, from
	// the story's page or with flai serve agent start (S-0115), rather than
	// flai serve when the story entered ready. Past says what that start
	// went past: a hold's reason, a full in-progress limit, or nothing
	// (S-0182).
	Started bool
	Past    []string
	// Begun says where a story in progress was begun, when this host has
	// had no agent for it and the operator has one started here (S-0177,
	// ADR-0064); nil otherwise.
	Begun *Begun
	// Shares are the shares in force that cleared a hold on the story, each
	// made by the agent of the story that held it (ADR-0134); none otherwise.
	Shares []Share
	// ProviderName and Provider are where the agent's model calls go, the
	// provider its agent names, resolved and checked by the caller
	// (ADR-0129); empty and nil when it names none, and the harness calls
	// its own default.
	ProviderName string
	Provider     *manifest.Provider
	// Queue is set when the merge queue runs for the project, as
	// mergequeue.Batches holds (ADR-0173): the orchestrator enqueues each
	// story it has verified rather than accepting it (S-0403). It is unset
	// for any other agent.
	Queue bool
	// Reasons are why flai serve started the orchestrator, one line each,
	// such as "moved S-0123 to review: <title>" (ADR-0209 §2): its session
	// reads them in FLAI_REASONS, one a line. None for any other agent.
	Reasons []string
	// PrimeReasons are the form flai prime --reason takes of each of Reasons,
	// in turn, such as "review:S-0123" (ADR-0217 §6): the start prompt writes
	// each beside its line and tells the orchestrator to prime with them,
	// each once. flai serve sets them only when every reason has a form; none
	// when one has not, such as asked or schedule, or for a run by hand, and
	// the orchestrator primes with the general pack. None for any other agent.
	PrimeReasons []string
	// Inner is the harness pilot runs inside a container, one of Adapters,
	// as pilot's start request names it (S-0507); empty for a start on the
	// host, which is every start flai serve makes today.
	Inner string `json:"inner,omitempty"`
	// MCPURLEnv names the variable that holds the URL of the host's HTTP MCP
	// server, for an agent in a container, which cannot start flai's server
	// on stdio (S-0507). Set, the adapter configures the HTTP server by the
	// variable's name, and never by its value; empty, the stdio server.
	MCPURLEnv string `json:"mcp_url_env,omitempty"`
	// MCPTokenEnv names the variable that holds the HTTP server's bearer
	// token; required with MCPURLEnv.
	MCPTokenEnv string `json:"mcp_token_env,omitempty"`
	// MCPHostEnv names the variable that holds the Host header the HTTP
	// server listens under, when it is not the URL's: the server's MCP SDK
	// refuses a Host that is not loopback as DNS rebinding, and a container
	// dials host.docker.internal. Optional.
	MCPHostEnv string `json:"mcp_host_env,omitempty"`
	// Cache is the folder mounted in the container where the harness keeps
	// its configuration, its sessions, and its MCP file, so that a later
	// container resumes the session (S-0507); empty on the host.
	Cache string `json:"cache,omitempty"`
	// Resume is set by pilot when its cache records Session already, so that
	// the harness goes on in it rather than start it again (S-0507); never
	// read from a start request.
	Resume bool `json:"-"`
}

// Input is one of a defined role's run's inputs: its name and its value.
type Input struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RoleAgent is the FLAI_AGENT a defined role's run works under (ADR-0207
// §2): the role's name, or <name>-<item> for a run on an item.
func RoleAgent(name, item string) string {
	if item == "" {
		return name
	}
	return name + "-" + item
}

// SentBack is a send back of a story from review to in-progress, which flai
// makes on a failing required check or a requested change on its pull
// request (ADR-0167), as the story's agent is told it.
type SentBack struct {
	Thread string // the thread flai serve opened on the story for it
	What   string // what failed, as the move's reason names it
	// Resumed is set when the agent goes on in its last run's session,
	// Request.Session; a new session is told how to take the story up.
	Resumed bool
	// Queue is set when the merge queue sent the story back, from a batch
	// that failed with it, rather than its pull request (ADR-0176).
	Queue bool
	// Acceptance is set when the orchestrator's acceptance sent the story
	// back over a rebase it could not settle (ADR-0189).
	Acceptance bool
}

// then is what follows the agent's fix of a story sent back on thread: the
// push of its pull request, from the merge queue its batching again, or from
// the orchestrator's acceptance its acceptance again.
func (b *SentBack) then(story string) string {
	if b.Acceptance {
		return fmt.Sprintf("flai serve resolves %[2]s once %[1]s is back in review, and the orchestrator accepts it again", story, b.Thread)
	}
	if b.Queue {
		return fmt.Sprintf("flai serve resolves %[2]s once %[1]s is back in review at a new commit, and the merge queue batches it again once it is enqueued", story, b.Thread)
	}
	return fmt.Sprintf("flai serve pushes the fix, and resolves %[2]s once %[1]s is back in review at a new commit", story, b.Thread)
}

// resumes says whether r goes on in a session an earlier run had: after an
// answer, after a send back with its last run's session, or when pilot's
// cache records the session already (S-0507).
func (r Request) resumes() bool {
	return r.Session != "" && (r.Resume || r.Answered != "" || r.SentBack != nil && r.SentBack.Resumed)
}

// Share is the sharing of paths that cleared a hold on the story an agent is
// started for (ADR-0134): the conversation it was made on, the story that
// held it and shared, the paths shared, and the split of who changes what.
type Share struct {
	Conversation string
	Holder       string
	Paths        []string
	Split        string
}

// Begun is where, when, and by whom a story was begun on another host, or
// outside flai serve on this one, as its files say.
type Begun struct {
	By    string // who last moved it to in-progress
	At    string // when
	Agent string // the agent its narrative names, when it has one
	Host  string // the host its narrative was opened on, when it says
	// Threads are the story's threads with an entry by someone other than
	// By or Agent since At.
	Threads []string
}

// Said is where the story was begun, as a clause: by whom, when, and where.
func (b *Begun) Said() string {
	said := "by " + orSomeone(b.By)
	if b.At != "" {
		said += ", who moved it to in-progress at " + b.At
	}
	if b.Agent != "" && b.Agent != b.By {
		said += "; its narrative names " + b.Agent
	}
	if b.Host != "" {
		return said + ", on the host " + b.Host
	}
	return said + ", on another host"
}

func orSomeone(who string) string {
	if who == "" {
		return "someone"
	}
	return who
}

// Host is the operator's say about one harness, from the host's configuration.
type Host struct {
	// Program is what is run; the adapter's default when empty.
	Program string
	// Args are the operator's own arguments: what the agent may do. Nil means
	// the adapter's default; an empty list means none.
	Args []string
}

// Start is what to run.
type Start struct {
	Harness string
	Argv    []string // program first
	Env     []string // added to flai serve's own
	// EnvFrom maps a variable of the child's to the name of a variable of
	// flai serve's own environment, whose value flai serve copies into it as
	// it spawns the child: a provider's key (ADR-0129). An adapter never
	// reads a key, and neither Argv, Env, nor Files ever carries one, so a
	// Start can be logged and recorded whole.
	EnvFrom map[string]string
	// Files are what the starter writes before it spawns the child: the MCP
	// file of an agent in a container, which pilot writes (S-0507). None for
	// a start on the host.
	Files []File
}

// File is a file a Start needs written before its child is spawned: where,
// what it holds, and its permissions.
type File struct {
	Path string
	Data []byte
	Mode os.FileMode
}

// Capabilities are what flai serve may count on a harness to do (ADR-0130).
type Capabilities struct {
	// Resume is that flai serve may go on with an agent in the session it
	// ran before, rather than start a new one.
	Resume bool
	// Guard is that flai serve may run flai guard as a hook on the agent's
	// calls, so that a call the guard refuses is never made.
	Guard bool
	// Ask is that flai serve may hold a call to ask the owner on a thread,
	// as permission_prompt does, and answer it allow or deny.
	Ask bool
	// RoleModels is that flai serve may give each sub-agent role its own
	// model.
	RoleModels bool
	// StrategicAgents is that flai serve may run a strategic agent, the
	// planner, the orchestrator, the analyzer, or the curator, from its
	// definition.
	StrategicAgents bool
	// Cost is that the harness reports the cost of its calls, for flai to
	// measure.
	Cost bool
}

// Adapter builds the command for one harness.
type Adapter interface {
	// DefaultHost is the program and arguments used when the operator has
	// set none.
	DefaultHost() Host
	// Start builds the command. It refuses what the story says that the
	// harness does not take.
	Start(req Request, host Host) (Start, error)
	// Capabilities are what flai serve may count on the harness to do.
	Capabilities() Capabilities
	// Reader is a fresh reader of the harness's log into usage events, for
	// the logs of one measurement read in order (ADR-0130).
	Reader() usage.Reader
}

// Adapters are the harnesses flai can start, by name.
var Adapters = map[string]Adapter{
	ClaudeCode: claudeCode{},
	Command:    command{},
}

// Names are the adapters' names, sorted.
func Names() []string {
	out := make([]string, 0, len(Adapters))
	for n := range Adapters {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// For is the adapter a story's agent is started with, and its name. A story
// that names no harness is started with the operator's command, when one is
// set.
func For(a *manifest.Agent, commandSet bool) (string, Adapter, error) {
	name := ""
	if a != nil {
		name = a.Harness
	}
	if name == "" {
		if !commandSet {
			return "", nil, fmt.Errorf("the story names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])")
		}
		name = Command
	}
	ad, ok := Adapters[name]
	if !ok {
		return "", nil, fmt.Errorf("flai cannot start harness %q; it can start %s", name, strings.Join(Names(), ", "))
	}
	if name == Command && !commandSet {
		return "", nil, fmt.Errorf("the story's harness is the host's command, and none is set (flai serve agent set -- <program> [args...])")
	}
	return name, ad, nil
}

// option is a config key a harness takes, and the values it may have.
type option struct {
	flag  string
	value *regexp.Regexp
	says  string // what a value must be, for a refusal
}

// options checks a story's config against what a harness takes and returns
// the flags, in key order.
func options(harness string, config map[string]string, takes map[string]option) ([]string, error) {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []string
	for _, k := range keys {
		o, ok := takes[k]
		if !ok {
			known := make([]string, 0, len(takes))
			for n := range takes {
				known = append(known, n)
			}
			sort.Strings(known)
			return nil, fmt.Errorf("%s takes no config key %q; it takes %s", harness, k, strings.Join(known, ", "))
		}
		if !o.value.MatchString(config[k]) {
			return nil, fmt.Errorf("%s config %s=%q: it must be %s", harness, k, config[k], o.says)
		}
		out = append(out, o.flag, config[k])
	}
	return out, nil
}

// Prompt is what the agent is asked to do: work its story, and nothing
// else, the way the project's conventions say, handing noisy work and the
// check before review to sub-agents (ADR-0059), and asking the designer
// through flai when it needs them. A resumed agent is told its question was
// answered and goes on in the session that already holds the rest; one whose
// story flai sent back from review is told the thread that says why, in its
// last session or, with none, as a story taken up (S-0394); one
// started to commit a worktree is told to do only that (S-0140); one the
// operator started past a hold or a full limit is told what it went past
// (S-0182). Every agent working a story is told how its issues are recorded
// and that the operator chooses at acceptance which become stories (S-0198).
// An agent for a ready story begins it with story_start (S-0274); one for a
// story already in progress, which story_start refuses, primes and takes it
// up with flai stream open. One started on a share that cleared a hold on
// its story is told the share's conversation, paths, and split (ADR-0134,
// shared). The planner is asked to plan its item instead
// (planPrompt), the orchestrator to keep the project's work moving
// (orchestratePrompt), the analyzer to write a report (analyzePrompt), the
// curator to keep the digest (curatePrompt), and a role the project defines
// to do what its file's body says, on its item, with its run's inputs
// (definedPrompt).
// Only the claude-code adapter sends a prompt.
func Prompt(r Request) string {
	if r.RoleFile != nil {
		return definedPrompt(r)
	}
	switch r.Role {
	case conventions.RolePlan:
		return planPrompt(r)
	case conventions.RoleOrchestrate:
		return orchestratePrompt(r)
	case conventions.RoleAnalyze:
		return analyzePrompt(r)
	case conventions.RoleCurate:
		return curatePrompt(r)
	}
	if r.Commit != "" && r.Answered == "" {
		return fmt.Sprintf(`You are %[1]s, started by flai serve on this host because the operator asked for the work left uncommitted in story %[2]s's worktree, %[3]s, to be committed: %[2]s is in review, and it cannot be accepted until that worktree is clean.

Do only this. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md, for how commits are made here. In the worktree, read git status and git diff. Commit the changes on story/%[2]s, in commits whose messages name %[2]s and say what changed and why, after the lint and tests the project runs for what they touch; fix what they find only when it is part of the same work. Discard a file only when it is plainly output that does not belong in the repository, and say which in the story's narrative with flai stream log %[2]s. Do not move %[2]s or its tasks, do not change anything the uncommitted work does not already change, and do not start other work. When git status in the worktree is clean, log what you committed with flai stream log %[2]s and end.

If a change cannot be committed without the designer deciding something, ask with the flai MCP tool thread_open on %[2]s, and commit what does not wait on the answer. Then write the narrative's Current state and Next steps with flai stream state, saying what you asked and what is left to commit, and end: flai serve starts you again in this session when the thread is answered, and your first inbox holds the answer. Do not hold the flai MCP tool wait_for_events for an answer; when it answers end: true, write the narrative's Current state and Next steps with flai stream state, and end.`, r.Name, r.Story, r.Commit)
	}
	if strings.HasPrefix(r.Answered, "MS-") {
		return fmt.Sprintf(`Another story's agent has written to %[2]s on %[3]s. Read it with the flai MCP tool inbox, which lists %[2]s's conversations under messages, or message_get (flai message show %[3]s on the host), answer what awaits your reply with message_reply, then go on working %[2]s to review as before.

%[4]s`, r.Name, r.Story, r.Answered, rules(r))
	}
	if r.Answered != "" {
		return fmt.Sprintf(`The designer has answered your question %[3]s on %[2]s. Read the answer with the flai MCP tool thread_get (or flai thread show %[3]s), then go on working %[2]s to review as before.

%[4]s`, r.Name, r.Story, r.Answered, rules(r))
	}
	if b := r.SentBack; b != nil && b.Resumed {
		return fmt.Sprintf(`flai sent %[1]s back from review to in-progress, on thread %[2]s: %[3]s. Read the thread with the flai MCP tool thread_get (flai thread show %[2]s on the host); your inbox lists it, awaiting you. Fix the branch in %[1]s's worktree, close out, and move %[1]s to review again; %[5]s.

%[4]s`, r.Story, b.Thread, b.What, rules(r), b.then(r.Story))
	}
	why := "because it entered ready."
	open := startStory(r.Story, "Begin")
	switch {
	case r.Begun != nil, r.Restart != "" && r.AutoRestart != "", r.SentBack != nil:
		open = takeUp(r.Story, "Prime")
	case r.Restart != "":
		open = fmt.Sprintf("If %[1]s is still ready, ", r.Story) + startStory(r.Story, "begin") + " If it is in progress already, " + takeUp(r.Story, "prime")
	case r.Started && len(r.Past) > 0:
		open = startStory(r.Story, "Begin") + fmt.Sprintf(" If story_start refuses %[1]s because the board holds it, start it as the operator asked with flai move %[1]s in-progress, which only warns of the hold, then ", r.Story) + takeUp(r.Story, "prime")
	}
	switch {
	case r.SentBack != nil:
		why = fmt.Sprintf("because flai sent it back from review to in-progress, on thread %[1]s: %[2]s. This host has no session of its last agent to go on in. The story is in progress, with a narrative, a branch, and a worktree: read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and read the thread with the flai MCP tool thread_get (flai thread show %[1]s on the host); your first inbox lists it, awaiting you. Fix the branch there, close out, and move %[3]s to review again; %[4]s.", r.SentBack.Thread, r.SentBack.What, r.Story, r.SentBack.then(r.Story))
	case r.Begun != nil:
		why = fmt.Sprintf("because the operator started it here, and this host has had no agent for it: it was begun %[1]s. Its branch and worktree may not be on this host, and what was not committed and pushed there is not here. Reconcile before you do anything else: run flai stream open %[2]s, which keeps the narrative and checks out story/%[2]s from this clone, else from the remote, else new from the main branch, and says which; read the narrative's Current state, Next steps, and log, and the story's tasks; compare them with what is committed on the branch; and go on from what is committed rather than starting over, doing again what the narrative says was done and is not there. Log what you found with flai stream log %[2]s.%[3]s", r.Begun.Said(), r.Story, answeredSince(r.Begun))
	case r.Restart != "" && r.AutoRestart != "":
		why = fmt.Sprintf("because its last agent %s, and flai serve started it again on its own, automatic restart %s; at the limit it asks the operator instead. The story is in progress, with a narrative, a branch, and a worktree: read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and go on from there rather than starting over.", r.Restart, r.AutoRestart)
	case r.Restart != "" && r.AfterLimit:
		why = fmt.Sprintf("because its last agent %s, and flai serve started it again on its own once the host's usage limit no longer held its provider, spending no automatic restart. The story may already be in progress, with a narrative, a branch, and a worktree: if so, read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and go on from there rather than starting over.", r.Restart)
	case r.Restart != "":
		why = fmt.Sprintf("because the operator restarted it: its last agent %s. The story may already be in progress, with a narrative, a branch, and a worktree: if so, read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and go on from there rather than starting over.", r.Restart)
	case r.Started && len(r.Past) > 0:
		why = fmt.Sprintf("because the operator started it now, from the story's page or with flai serve agent start, past what kept flai serve from starting it: %s. That was the operator's decision: pull %[2]s as they asked, though flai warns that it is held or that the limit is full. Keep its touches to what it changes; do not narrow them only to clear the hold. Where its claim overlaps another open story's, say in the narrative which paths the two share.", strings.Join(r.Past, "; "), r.Story)
	case r.Started:
		why = "because the operator started it now, from the story's page or with flai serve agent start."
	}
	return fmt.Sprintf(`You are %[1]s, started by flai serve on this host to work story %[2]s in the project at %[3]s, %[5]s%[8]s

Work %[2]s to review, and no other story. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md. %[7]s A brief is not the document: when one bears on the story, read it, or its section that does, with the flai MCP tool doc_get and its heading (flai doc show --heading on the host) before relying on it or changing what it describes, and find sections by their words with doc_search. Write the story's tasks if it has none, or review the ones the planner drafted, and work them in that worktree. When a task is done, close it with flai task done T-nnnn -m "<message>" in the worktree (or the flai MCP tool task_done), with its docs and work-item updates in the change: it commits on story/%[2]s, runs flai stream sync %[2]s, moves the task to done, logs it in the narrative, widens the touches, runs flai check, and answers your inbox, stopping at the first step that fails. When the sync stops on conflicts, resolve each path it lists in the worktree, git add it, and git rebase --continue, then call it again; when the check stops, fix what it found and call it again. Then run the task's tests with flai test and the paths it changed (or the flai MCP tool test with them), and close any fix they need by calling flai task done again. flai stream sync does the branch's git work and refuses while anything is uncommitted: never start a rebase or merge by hand. Keep the narrative's Current state and Next steps true: rewrite them at every task transition with flai stream state %[2]s --current "<text>" --next "<text>" (or the flai MCP tool stream_state), never by editing the narrative.

%[6]s

%[4]s`, r.Name, r.Story, r.Root, rules(r), why, delegation(r), open, shared(r))
}

// shared tells an agent started on shares (ADR-0134), as a paragraph of its
// own, each share's conversation, the story that shared, the paths, and the
// split; to keep to the split; that its first inbox lists the conversation
// under messages; and to answer there with message_reply when the split no
// longer fits. Empty when it was started on none.
func shared(r Request) string {
	if len(r.Shares) == 0 {
		return ""
	}
	var each, convs []string
	for _, s := range r.Shares {
		paths := make([]string, len(s.Paths))
		for i, p := range s.Paths {
			paths[i] = "`" + p + "`"
		}
		each = append(each, fmt.Sprintf("On %s, %s's agent shared %s with %s, split so: \"%s\".", s.Conversation, s.Holder, strings.Join(paths, ", "), r.Story, strings.Join(strings.Fields(s.Split), " ")))
		if !slices.Contains(convs, s.Conversation) {
			convs = append(convs, s.Conversation)
		}
	}
	split, a := "the split", "the split"
	if len(r.Shares) > 1 {
		split, a = "each split", "a split"
	}
	id, read, there := convs[0], "read it", "there"
	if len(convs) > 1 {
		id, read, there = "<MS-nnnn>", "read each", "on its conversation"
	}
	return fmt.Sprintf(`

%[1]s works on a share (ADR-0134): its claim overlaps the claim of another story, whose agent shared the overlapping paths with it and split the work, so the overlap no longer holds it. %[2]s Keep to %[3]s: in the paths shared, change only what it gives %[1]s, and leave the rest to the story that shared. Your first inbox lists %[4]s under messages: %[5]s with the flai MCP tool message_get (flai message show %[6]s on the host). If %[7]s no longer fits the work, say so %[8]s with message_reply (flai message reply %[6]s on the host) before you change what it does not give %[1]s.`,
		r.Story, strings.Join(each, " "), split, strings.Join(convs, ", "), read, id, a, there)
}

// startStory is how an agent begins a ready story (S-0274): one call to the
// MCP tool story_start, or flai story start on the host, moves it to
// in-progress, opens its stream, primes, and answers the inbox, and the
// agent reads every part of the pack the answer does not hold. lead is its
// first word, Begin or begin.
func startStory(story, lead string) string {
	return fmt.Sprintf("%[2]s with the flai MCP tool story_start with %[1]s (flai story start %[1]s on the host): in one call it moves %[1]s to in-progress, opens its narrative and its branch in a worktree, primes your session, and answers your inbox. Work in the worktree it answers. Its pack holds the conventions that apply and what the story names whole, and briefs the design and ADRs its topics and links select, within a size budget, in parts: read every part. When the answer holds part 1, read part 2 on with the flai MCP tool prime and %[1]s; when it holds only the pack's header, read every part with prime.", story, lead)
}

// takeUp is how an agent takes up a story already in progress, which
// story_start refuses: it primes with flai prime --story and opens the
// stream again with flai stream open, which keeps the narrative. lead is
// its first word, Prime or prime, as the sentence it begins needs.
func takeUp(story, lead string) string {
	return fmt.Sprintf("%[2]s your session with flai prime --story %[1]s (or the flai MCP tool prime), which prints the conventions that apply and what the story names whole, and briefs the design and ADRs its topics and links select, within a size budget, and take %[1]s up with flai stream open %[1]s, which prints the worktree to work in.", story, lead)
}

// planPrompt is what the planner is asked to do (S-0208): plan its item, an
// epic or a story, as strategic-agents.md says, through flai alone, drafting
// a story's tasks or revisiting the open ones it has (S-0255); ask the
// operator on the item for an input it owns that is missing; and end with a
// summary, which flai serve logs as the run's activity (ADR-0079). An epic's
// planner writes its stories as drafts, summarises the plan in one thread on
// the epic, proposes there what it would change in the stories the epic
// already has, and names in its summary the stories it created and revisited
// (S-0209). It enriches each story it drafts as it would a story, drafts
// that story's tasks, and those of each draft it revisits that has none, in
// the story planner's words, names them and their layers in the epic's
// thread, and names by ID in its summary the tasks it created too (S-0300).
// A story's planner drafts its tasks, or revisits the open ones, and names
// in its summary the tasks it created and revisited (S-0255).
// It enriches a story from flai touches suggest, flai forecast, and flai cod,
// and records why under a ### Planning heading in the story's Notes (S-0210).
// It names files, not folders, in touches, keeping a folder only where files
// no task can name yet may be added, and says why there (S-0295, ADR-0096);
// it is told flai refuses a task's folder touch, and a story's its Planning
// note does not name, and that after names only a story whose accepted
// change the story needs on main (S-0468, ADR-0188). When the orchestrator
// asked for the item again, for what the operator wrote on r.AskedOn since
// it was planned before, its first paragraph says so in place of the
// operator's asking, and has the planner read that thread with thread_get
// and plan what the operator asked for there (S-0524, ADR-0222).
func planPrompt(r Request) string {
	kind := workitem.TypeOfID(r.Item)
	story := r.Item
	if kind == workitem.Epic {
		story = "S-nnnn"
	}
	enrich := fmt.Sprintf("its predicted touches, a forecast, and a cost of delay value worked out from the operator's inputs. "+
		"Run flai touches suggest %[1]s, adding the paths its goal, criteria, and linked design name when it declares no touches, and predict its touches from what that lists, its goal and criteria, the design documents it links, and the code layout, keeping every touch it already declares. "+
		"Name files, not folders, in the touches you write for the story and for each of its tasks: keep a folder touch only where the story may add files there that no task can name yet, since a folder touch claims every file below it and, while the story is in progress, holds every ready story that touches one; record each folder touch you kept, and why, under the story's ### Planning heading. "+
		"flai refuses the folder touches this asks you not to write, writing nothing: a task's always, and a story's unless the story's Notes name the folder under ### Planning, which one item_edit may add with the touch (ADR-0188). "+
		"Name in a story's after only a story whose accepted change it needs on main, never an order of preference, which the pull order sets. "+
		"Run flai forecast %[1]s and flai cod %[1]s. "+
		"Review each figure, adjust it where you have a reason and state the reason, and write the touches, the forecast (duration, delivery, and basis), and the cost of delay value through flai: item_edit, or flai edit and flai touches. "+
		"In the story's Notes, under a ### Planning heading that is yours to rewrite, record where each touch came from (declared, co-change, design, or layout) and why each figure stands or was adjusted, and leave the rest of the Notes as it was", story)
	// The story's planner and the epic's say the same about drafting a
	// story's tasks, in the same words (S-0300).
	draftTasks := "draft the tasks that deliver its outcome, each with ## Work and ## Done when in its body, a nature, tags, touches (the paths it changes), and after (the tasks of the story it waits for), so that they form layers as work-management.md says, and create each in the backlog with the flai MCP tool item_new, type task and parent the story, or flai task new"
	taskRules := "Tasks carry no topics: when a task reaches a topic the story lacks, add the topic to the story. " +
		"Make each task you write pass flai check --strict and the markdown lint: flai refuses one that does not, so fix what the refusal names and write it again."
	work := fmt.Sprintf("It is a story: enrich it with %[1]s. "+
		"If it has no tasks, %[3]s. "+
		"If it has tasks, revisit each one not done or cancelled against the story's outcome: re-enrich its touches and after with item_edit, create the tasks the outcome still lacks, and propose in the plan's thread any task you would split, merge, or drop. "+
		"Never cancel a task, or rewrite the title, Work, or Done when of a task you did not write, without asking on that thread. "+
		"%[4]s\n\n"+
		"Open one thread on %[2]s that summarises the plan. In the plan's thread, name the tasks, their order and layers, and the assumptions you made.", enrich, r.Item, draftTasks, taskRules)
	summary := "End with a one-line summary that names by ID the tasks you created and the tasks you revisited"
	if kind == workitem.Epic {
		work = fmt.Sprintf("It is an epic. If it has no stories, draft the stories that deliver its outcome, each with a goal, acceptance criteria as checkboxes, a nature, tags, topics, touches, and after, and create each with draft true in the backlog (item_new's draft, or flai story new --draft). "+
			"If it has stories, revisit each one not done or cancelled against the epic's outcome. "+
			"Enrich each story you draft, and again each one you revisit, as you would a story, S-nnnn being its ID: %[1]s. "+
			"Size stories as work-management.md says, and make each one you write pass flai check --strict.\n\n"+
			"For each story you draft, and each one you revisit that is a draft with no tasks, %[3]s. "+
			"A task carries no draft flag: it is a draft because its story is one. "+
			"%[4]s\n\n"+
			"Open one thread on %[2]s that summarises the plan: the stories, their order (their after), each story's tasks and their layers, and the assumptions you made. In that same thread, propose each story you would split, merge, add, or drop, and create drafts for the additions only: never cancel a finalized story or rewrite its words without asking.", enrich, r.Item, draftTasks, taskRules)
		summary = "End with a one-line summary that names by ID the stories and tasks you created and the stories you revisited"
	}
	why := fmt.Sprintf("because the operator asked for %s to be planned, in the project at %s.", r.Item, r.Root)
	if r.AskedOn != "" {
		why = fmt.Sprintf("in the project at %[3]s, because the orchestrator asked for %[1]s to be planned again: it was planned before, and the operator wrote on %[2]s since its last plan. "+
			"Read %[2]s with the flai MCP tool thread_get (flai thread show %[2]s on the host) and plan what the operator asked for there, drafting the additions and revisiting the rest as design/conventions/strategic-agents.md says.", r.Item, r.AskedOn, r.Root)
	}
	return fmt.Sprintf(`You are %[1]s, the planner, started by flai serve on this host %[3]s

Plan %[2]s, and nothing else, as design/conventions/strategic-agents.md says under As the planner. Prime your session with flai prime --role plan --%[4]s %[2]s (or the flai MCP tool prime with role plan and %[4]s %[2]s), which prints the conventions you work by and what %[2]s names whole, and briefs the design its topics select. A brief is not the document: read the section that bears on the plan with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox. Read %[2]s with item_get, and what it links with item_get and doc_get. Hand wide search of the code, such as for a story's touches, to the explorer with the Agent tool.

%[5]s

Work in the main checkout and write only through flai: the flai MCP tools item_new and item_edit, or the flai CLI. Never edit a file yourself, code or anything else, never move an item past backlog, and never finalize a draft. Never overwrite the operator's inputs: a cost of delay's inputs, a story's estimate, and a finalized story's words; never cancel or rewrite a finalized story without asking.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on %[2]s, your recommended answer first, plan what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

%[6]s: flai serve logs the run's activity in wip/agents/planner.md with it.`, r.Name, r.Item, why, kind, work, summary)
}

// orchestratePrompt is what the orchestrator is asked to do (S-0218): keep
// the project's work moving as strategic-agents.md says, only as far as the
// permissions the operator sets in orchestration.permissions allow, which
// with its policy and release policy it reads again before each decision,
// since the operator may change them while it runs and the guard holds each
// call to them as they are then (S-0229), taking
// every order, candidate, and release figure from flai's commands rather
// than working it out; log each action with activity_log, its reason and
// the policy figure that justified it; and end. Its first paragraph lists
// the reasons flai serve started it for, r.Reasons, one a markdown list
// item, and it acts on each, logs each action, then ends, since flai serve
// starts it again when there are more and wait_for_events answers it end
// once its reasons are done; a run with none, by hand, looks at the inbox
// and the board once, then ends (ADR-0209 §3, §4). When flai serve gives each
// reason its form for prime, r.PrimeReasons, one for each line, each line
// names its form, and it is told to prime with flai prime --role orchestrate
// and a --reason for each form, or MCP prime with role orchestrate and those
// reasons, each once, reading every part of a pack that holds the conventions
// those reasons need and what each reason needs, with no design briefs, so
// that it finds the rest with doc_search; otherwise, by hand or for a reason
// with no form, it primes with role orchestrate alone, for the general pack,
// which briefs the design its role's topic selects (ADR-0217 §6). It is told what to do with each of plan_backlog_epics,
// finalize_drafts, promote_to_ready, and order_ready, and by which command
// (S-0219), and what to do with the threads awaiting the operator, those
// inbox marks awaiting "operator" and lists in its role without all
// (ADR-0181), under each value of answer_threads, which it reads afresh
// each time since the operator may change it while it runs (S-0220), and
// how it accepts a story
// in review under accept_reviews: flai verify runs, or its last result is
// read, at the branch's head, its verifier checks the diff against each
// criterion (S-0270), flai accept --dry-run names the blockers, and it accepts with the verified
// commit and evidence for every criterion, or leaves the story in review and
// says on a thread what is missing (S-0221, ADR-0093). It is told that the
// main checkout's uncommitted files stop it only as dirty_main, that a thread
// of its own is own_thread_open, which it resolves with a reason once its
// concern is gone, that it opens a thread on a story in review only for a
// reason not to accept it, that a story flai accept sent back over a rebase
// it could not settle needs nothing more of it, and that it finishes an
// acceptance of its own that stopped after its move to done (ADR-0189).
// While publish is on,
// it evaluates the release policy after each acceptance and publishes
// through release_publish alone when the policy is met, or under judgement
// when it judges the batch coherent and complete, never a batch whole_epics
// holds back; it logs each release, each decision not to publish, and each
// refusal, which it also raises on a thread to the operator (S-0222). It
// never edits a file or works a story. A refusal from flai
// guard or from flai ends that attempt, which it logs and does not retry
// until something changes, and it never works around the refusal. While
// plan_backlog_stories is on, it starts the planner, after the epics, for
// each story flai plan --candidates lists, one at a time, and settles the
// threads a story's planner opens whatever answer_threads says: it approves
// a plan that fits the story, answers its questions, and chooses a cost of
// delay input, the figure the planner recommends unless the thread or the
// story gives a reason for an alternative, which it sets with item_edit as
// inputs only, never a value; it replies, resolves the thread, and logs
// each, and never confirms a recommendation (S-0328, ADR-0119). It marks a
// recommendation on money or a change to the host operator_only, and is told
// that while orchestration.permissions.confirm_after is set flai serve
// confirms its other recommendations on a thread a story's agent or a
// planner opened once they have waited that long unanswered (ADR-0206). The
// paragraph is unconditional: Request carries no permissions, and the
// orchestrator reads them afresh before each decision. While the
// merge queue runs for the project, r.Queue, it is told to enqueue each
// story it has verified with queue_add instead of running flai accept, to
// log each enqueue, and to leave landing, bisecting, holding, resuming,
// dequeuing, and the blockers the queue settles to the queue (S-0403).
func orchestratePrompt(r Request) string {
	queue := ""
	if r.Queue {
		queue = "\n\nThe merge queue runs for this project: enqueue each story you have verified with the flai MCP tool queue_add, with the commit and the evidence, instead of running flai accept, and log each enqueue with activity_log naming the story, its place in the queue, and the policy figure. Never land, bisect, hold, resume, or dequeue by hand, and open no thread for a blocker the queue settles: a verify not at the branch's head, a dirty main checkout, or a rebase that stops inside the shared paths."
	}
	why, act := fmt.Sprintf("You are %s, the orchestrator, started on this host for the project at %s with no reasons, as a run by hand is: look at the inbox and the board once, act on what they show, and then end.", r.Name, r.Root),
		"Act on what the inbox and the board show, once"
	prime := "Prime your session with flai prime --role orchestrate (or the flai MCP tool prime with role orchestrate), which prints the conventions you work by, briefs the design your role's topic selects, and loads the digest's answers and recommendations ranked against the threads awaiting the operator (ADR-0194). A brief is not the document: read the section that bears on a decision with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search."
	if len(r.Reasons) > 0 {
		lines, listed := r.Reasons, "listed here"
		if len(r.PrimeReasons) == len(r.Reasons) {
			lines, listed = make([]string, len(r.Reasons)), "listed here, each with the form flai prime takes it in"
			for i, l := range r.Reasons {
				lines[i] = l + " (prime: " + r.PrimeReasons[i] + ")"
			}
			prime = primeReasonsSentence(r.PrimeReasons)
		}
		why = fmt.Sprintf("You are %s, the orchestrator, started by flai serve on this host for the project at %s, while the orchestrate host action is on, for the reasons %s:\n\n- %s", r.Name, r.Root, listed, strings.Join(lines, "\n- "))
		act = "Act on each reason listed above"
	}
	return fmt.Sprintf(`%[3]s

Keep the project's work moving, and do nothing else, as design/conventions/strategic-agents.md says under As the orchestrator. %[5]s Call the flai MCP tool inbox, and read the board with the flai MCP tool board. Hand wide search to the explorer with the Agent tool.

Act only within the permissions the operator sets in system-flow.yaml under orchestration.permissions, each off by default, and by its policy, orchestration.policy. Ask the planner to plan a backlog epic, with the flai MCP tool plan, only while plan_backlog_epics is on; ask it to plan a backlog story, with plan, and settle the threads a story's planner opens, giving a backlog story its first cost of delay inputs with item_edit, only while plan_backlog_stories is on; set on an epic or a story the cost of delay inputs the operator decided, with item_edit, only while set_cost_of_delay is on; finalize a draft only while finalize_drafts is on; promote a story to ready only while promote_to_ready is on; order the ready column only while order_ready is on; answer a thread, or recommend an answer, only as answer_threads says; accept a story only while accept_reviews is on; publish a release only while publish is on; add, change, or remove a sub-project a story or a thread asks for only while edit_subprojects is on. Do none of it while its permission is off, and if a permission is unclear, ask; do not act. The operator may change your permissions, your policy, and the release policy, orchestration.release, while you run: read them again in system-flow.yaml in the main checkout before each decision rather than keep what you read at the start, and when flai guard's verdict on a call differs from what you read, the guard's verdict holds.

Take every figure from flai's commands and never do the arithmetic yourself: the epics to plan from flai plan --candidates, and after them the stories, the drafts complete enough to finalize from flai promote --drafts, the ready column's order from flai order --by (the flai MCP tool order_by_policy), the stories that could go to ready from flai promote --candidates (promote_candidates), and whether a release is due from flai release --evaluate (release_evaluate).

With each permission, do this. While plan_backlog_epics is on, run flai plan --candidates and start the planner with the flai MCP tool plan for each epic it lists, one at a time, and for no other. While plan_backlog_stories is on, after the epics, start the planner with the flai MCP tool plan for each story flai plan --candidates lists, those whose type is story, one at a time, and for no other. While finalize_drafts is on, run flai promote --drafts: finalize a draft it lists as complete, and whose criteria, touches, forecast, and value you judge consistent, with the flai MCP tool item_edit giving only its id and draft false; for any other draft, open one thread on the story saying what it lacks or what is inconsistent, once, and leave the draft as it is. While promote_to_ready is on, run flai promote --candidates and move its candidates to ready in its order with the flai MCP tool item_move while the ready column's limit has room: never a draft, and never a held story. While order_ready is on, run flai order --by <policy> --apply, with the policy orchestration.policy names, after each change to the ready column, yours or another's; it keeps in place a story placed by hand within the last day. Never place a story by hand yourself.

Reply on threads only as answer_threads says, and read it in system-flow.yaml in the main checkout each time before you act on threads: the operator may change it while you run. Take from inbox the threads it marks awaiting "operator", which it lists in your role without all true, leaving out the threads you opened, those whose last entry is yours, those whose pending_recommendation is not null, and those whose last entry is by neither a story's agent nor a planner, such as flai serve's notices. While answer_threads is off, leave them alone. While it is recommend, reply to each with the flai MCP tool thread_reply, recommendation true, and a source: the ADR, design section, or convention your answer rests on, as <path> or <path>#<heading>, read with doc_get first. When a digest entry in your pack answers a thread's question, it is the answer the operator gave last time: cite it, with the entry's path, design/digest/D-nnnn-<slug>.md, as the source, and name its D-nnnn in the reply. While it is autonomous, answer with thread_reply and a source when a source settles the question; post a recommendation instead, citing what it draws on, which escalates it to the operator, when no source settles it, or when it asks for the operator's judgement: a decision not yet recorded, a change of scope, or money, such as a cost of delay input, an estimate, or spend, save the cost of delay inputs a story's planner asks for while plan_backlog_stories is on. thread_reply logs the reply and its source in your decision log. Never resolve a thread you did not open, but a story's planner's while plan_backlog_stories is on, never answer a thread you opened, and never confirm a recommendation: the operator does.

While plan_backlog_stories is on, settle the threads a story's planner opened, those whose opener is planner-S-nnnn, whatever answer_threads says. Approve a plan whose tasks, touches, and figures fit the story: reply so with thread_reply, and resolve the thread with the flai MCP tool thread_resolve. Answer its questions with thread_reply, with a source when one settles them. For a cost of delay input it asks for, take the figure it recommends unless the thread or the story gives a reason for one of the alternatives it lists. Set it with the flai MCP tool item_edit, cost_of_delay on the story, giving the inputs only, revenue_per_week, penalty_per_week, or time_lost_per_cycle, never a value and nothing else: the planner works the value out on its next run. Then reply naming what you set, and resolve the thread. Log each approval, answer, and cost of delay you set with activity_log. When a question asks for a change of scope or a decision not yet recorded, or for a cost of delay input other than a backlog story's first, while neither it nor its epic has any, post a recommendation instead, which escalates it to the operator, and leave the thread open. Never confirm a recommendation: the operator does. While plan_backlog_stories is off, a story's planner's threads are the operator's, as answer_threads says.

While set_cost_of_delay is on, set a cost of delay input the operator decided on a thread, as their answer or a recommendation they confirmed, on the epic or the story it concerns, while that is not done, cancelled, or archived, whether or not it or its epic has inputs: item_edit, cost_of_delay, giving the inputs only, revenue_per_week, penalty_per_week, or time_lost_per_cycle, never a value, a removal, or a figure of your own (ADR-0224). Log it with activity_log, naming the thread.

Wherever you post a recommendation, under answer_threads or plan_backlog_stories, give thread_reply operator_only true with it when it is about money, such as spend, a paid service, or a price, or about a change to the host, such as installing or upgrading software, configuration outside the repository, or credentials: that answer is the operator's alone. While orchestration.permissions.confirm_after is set, flai serve confirms your other recommendations on a thread a story's agent or a planner opened once they have waited that long with no answer from the operator. Recommend there only an answer you would stand behind as the answer (ADR-0206).

While accept_reviews is on, take each story in review in turn. Verify it at the head of its branch, story/<S-nnnn>: read its last result with flai verify <S-nnnn> --last --json, and when it did not pass or its commit is not the branch's head (git rev-parse story/<S-nnnn>), run flai verify <S-nnnn> --json, or the flai MCP tool verify, which runs the tests, the lint, and flai check --strict in the story's worktree and answers the commit it verified. A run that did not pass is a blocker. Then hand its worktree, %[1]s/.flai-cache/worktrees/<S-nnnn>, and that result to the verifier with the Agent tool: it reads the result with flai verify <S-nnnn> --last rather than running the suite, checks the diff against each acceptance criterion, naming for each the changed files that meet it, and names the commit it verified. Run flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run with that commit, and read the blockers. Files the main checkout has not committed stop nothing: flai leaves them out of the acceptance, and the blocker dirty_main appears only while the story's change would overwrite one of them, naming the files and the command that sets them aside. A thread you opened on the story or one of its tasks that nobody else wrote on is the blocker own_thread_open, apart from thread_open: once its concern is gone, resolve it with the flai MCP tool thread_resolve and a reason, then accept. With no blocker and every criterion matched to changed files, write the evidence, a Verdict: line from the verifier's report and one list item per criterion, - <n>: <files>, and run flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -, passing the evidence on standard input with a heredoc in the shell, since you cannot write a file; then log the acceptance with activity_log, naming the story and the commit. Otherwise leave the story in review, open a thread on it with the flai MCP tool thread_open saying what is missing, each blocker and each criterion you could not check against the diff, unless a thread of yours on it says so already, and log that decision. You accept on flai's blockers and the story's criteria: open a thread on a story in review only for a reason not to accept it, and resolve it with thread_resolve and a reason once that reason is gone. A commit added after the verifier's run makes flai refuse: verify again. When flai accept says it sent the story back, its rebase stopped on conflicts flai cannot settle: flai aborted the rebase, moved the story back to in-progress, and opened a thread on which flai serve starts the story's agent. Log that with activity_log and do nothing more for the story: open no thread and move nothing. When an acceptance of yours stops after its move to done, before its archive or its commit, run flai accept <S-nnnn> --by orchestrator again, with no --verified and no evidence, once what stopped it is cleared, to finish it; one whose move to done anyone else made is the operator's to finish, and flai's refusal names who made it. Never move a story to done with item_move, and never accept a story while accept_reviews is off.%[2]s

While publish is on, after each acceptance, yours or one wait_for_events reports as a story moved to done, call the flai MCP tool release_evaluate. Under the threshold or theme policy, when the evaluation is met, call the flai MCP tool release_publish with the figure that was met as its reason. Under judgement, call release_publish only when you judge the unreleased work coherent and complete, with that reasoning as its reason. Under any policy, never publish a batch the evaluation holds back under whole_epics, the stories it names in held_by_epic. Log each release with activity_log: the policy, its figures, the versions and tags, and the items bundled, all as release_publish returned them. Log a decision not to publish too, with the evaluation's figures. When release_publish refuses, log the refusal with activity_log, and open a thread to the operator with the flai MCP tool thread_open on the most recently accepted story of the batch, with the refusal's words and what fixes it; do not try again until that thread is answered or the next acceptance. While publish is off, neither evaluate nor publish. Publish only through release_publish: flai guard refuses you flai release other than --evaluate, flai push, git push, and git tag, whatever your permissions.

While edit_subprojects is on, add, change, or remove a sub-project, an entry of projects in system-flow.yaml, only when a story or a thread asks for that change, never on your own judgement (ADR-0195). Read them with the flai MCP tool subprojects. Make the change in the main checkout with the shell: flai subproject add <name> --kind <kind>, with --path and --tags as asked, flai subproject set <name> with what changes, or flai subproject remove <name>, each with --autocommit and a --trailer giving the Co-Authored-By line that names your model, which commits system-flow.yaml and the README.md of a folder it created in one step. Do not use subproject_edit, which commits nothing. Log the change with activity_log, naming the story or thread that asked and the commit, and reply on the thread when one asked. While edit_subprojects is off, change no sub-project.

Log each action when you have taken it with the flai MCP tool activity_log, kind orchestrator: what you did, to which item, why, and the policy figure that justified it, as flai gave it: the story's value, its value over its duration, its forecast, or the candidate's rank. %[4]s, log each action, then end: flai serve starts you again when there is more. The flai MCP tool wait_for_events answers you end: true, why: reasons done, once your reasons are done: end then.

Work in the main checkout and write only through flai: the flai MCP tools, or the flai CLI. Never edit code or documents, and never work a story yourself. flai guard refuses a call outside your permissions and names the permission it needs: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt: log it with activity_log, with the refusal, and do not try it again until something it depends on changes. When a decision needs the operator, ask with the flai MCP tool thread_open on the item it concerns, your recommended answer first, and do what needs no answer meanwhile.`, r.Root, queue, why, act, prime)
}

// primeReasonsSentence tells the orchestrator to prime with its reasons'
// forms, each once in the order first given, as flai prime --reason flags
// and as MCP prime's reasons, reading every part of a pack that holds only
// what those reasons need (ADR-0217).
func primeReasonsSentence(forms []string) string {
	var once []string
	for _, f := range forms {
		if !slices.Contains(once, f) {
			once = append(once, f)
		}
	}
	flags, quoted := make([]string, len(once)), make([]string, len(once))
	for i, f := range once {
		flags[i], quoted[i] = "--reason "+f, `"`+f+`"`
	}
	return fmt.Sprintf("Prime your session with flai prime --role orchestrate %s (or the flai MCP tool prime with role orchestrate and reasons [%s]), reading every part. It prints the conventions you work by, with the sections those reasons do not need left out, and per reason what it needs: the thread and its item, the story in review and its last flai verify result, the plan, the draft, the release figures, or the board, with the digest's answers and recommendations ranked against the reasons' threads (ADR-0217). It briefs no design: find the section that bears on a decision by its words with the flai MCP tool doc_search, and read it with doc_get and its heading before relying on it.", strings.Join(flags, " "), strings.Join(quoted, ", "))
}

// Focuses are what the analyzer can be asked to look for (S-0223):
// bottlenecks in the flow, gaps between the design and the code (intent),
// and technical and security risks. A run asked for none looks for all
// three, and its report's focus is AllFocus.
var Focuses = []string{"bottlenecks", "intent", "risk"}

// AllFocus is the focus of a report from a run asked for no focus.
const AllFocus = "all"

// lookFor is the findings the analyzer looks for under each focus.
var lookFor = map[string]string{
	"bottlenecks": "bottlenecks in the flow of work, from the cumulative flow, the time items spend in each state, the time they wait, and the holds on them",
	"intent":      "gaps between what design/system says and what the code does: a design section the code does not meet, and code the design does not describe",
	"risk":        "technical and security risks in the code and in how it is built, tested, and released",
}

// analyzePrompt is what the analyzer is asked to do (S-0223): look for the
// findings its focus names, or all three kinds with none, as
// strategic-agents.md says; read the metrics from flai stats --json, the
// design with doc_search and doc_get, and the issues, handing wide search of
// the code to the explorer; write one report, design/analysis/<date>-<focus>.md,
// with its front matter and one section per finding with its evidence,
// severity, and estimated impact, and add it to design/analysis/README.md;
// file each actionable finding as an issue with its class, impact, evidence,
// and the report, bumping the open issue that already records it, and link
// each issue from its finding (S-0224); edit nothing else and author no
// stories, flai issue story among them, which flai guard enforces; and
// end with a one-line summary that names the report, which flai serve logs
// as the run's activity, as the planner's is (ADR-0079).
func analyzePrompt(r Request) string {
	focus := r.Focus
	look := "Look for " + lookFor[focus] + ", and for nothing else."
	if focus == "" {
		focus = AllFocus
		look = fmt.Sprintf("Look for three kinds of finding: %s; %s; and %s.", lookFor["bottlenecks"], lookFor["intent"], lookFor["risk"])
	}
	report := "design/analysis/<date>-" + focus + ".md"
	return fmt.Sprintf(`You are %[1]s, the analyzer, started by flai serve on this host because the operator asked for an analysis, focus %[2]s, of the project at %[3]s.

Analyze the project, and do nothing else, as design/conventions/strategic-agents.md says under As the analyzer. Prime your session with flai prime --role analyze (or the flai MCP tool prime with role analyze), which prints the conventions you work by and briefs the design your role's topic selects. Call the flai MCP tool inbox.

%[4]s

Read the metrics with flai stats --json, and take every figure from it rather than working it out yourself. Read the design with the flai MCP tools doc_search, which finds sections by their words, and doc_get and its heading: a brief is not the document, so read the section a finding rests on before relying on it. Read the issues under design/issues, design/issues/summary.md first. Hand wide search of the code to the explorer with the Agent tool.

Write one report, %[5]s, where <date> is today's date in UTC as YYYY-MM-DD. Give it front matter with title, updated, status draft while you write it and active once it is done, focus %[2]s, and the window its metrics cover, from and to, as dates. Give it one section per finding, with its evidence (the metric figures as flai gave them, the file paths, and the design sections quoted), its severity, and its estimated impact: the time it loses per cycle, or the revenue or penalty it puts at stake where the design states them. Add the report to design/analysis/README.md.

File each actionable finding as an issue, and make no story of it. List the open issues with flai issue list --json, or read design/issues/summary.md. When one already records the finding, under whatever title, bump it with flai issue bump <id> --report %[5]s and the finding's impact (or the flai MCP tool issue_bump). Otherwise file one with flai issue new "<title>" --class <class> --report %[5]s --json (or the flai MCP tool issue_new): the class from the finding, defect, efficiency, or impression for a risk with no measured instance; its impact as --time-lost-per-cycle (a duration such as 4h), or --revenue-per-week or --penalty-per-week (an amount in planning.currency), the figures you gave the finding; and --evidence, the evidence behind them. With --report, flai bumps an open issue of the same title rather than opening a second, and its outcome says which it did; either way the issue's Remediation section links your report. Link each issue you filed or bumped from its finding in the report, as [I-nnnn](../issues/<file>), the file name of the path flai returns. Never run flai issue story or flai issue close: the stories a finding calls for are the planner's and the operator's to write.

Edit nothing else: no code, no design, no issue, and no work item, and author no stories; the issues you file and bump, flai writes. flai guard refuses an edit outside design/analysis and any write to a work item: never work around a refusal, by another tool, another command, or the shell.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on your report, your recommended answer first, analyze what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

End with a one-line summary that names the report, %[5]s with its date: flai serve logs the run's activity in wip/agents/analyzer.md with it, with the run's cost, so you need not log it yourself with activity_log.`, r.Name, focus, r.Root, look, report)
}

// curatePrompt is what the curator is asked to do (S-0431, ADR-0196): keep
// the digest as strategic-agents.md says under As the curator, whether flai
// serve started it after an acceptance, on digest.schedule, or because the
// operator asked (S-0432, ADR-0199); prime with role curate first and call
// inbox; read what the archive and the threads gained since the run's since
// with flai digest harvest --json --since, or everything on its first run,
// when it has none, and the digest as it is with flai digest list --json and
// show; for each candidate write an entry with flai digest new when none
// holds the fact, or merge it into the one that does with flai digest edit,
// and add to that entry's asked_again a thread that asked what it already
// answered; in every run, look for absorption, searching the design with
// doc_search on each active entry's title and marking absorbed with flai
// digest absorb an entry a document states, then run flai digest compact
// once; never write a secret's value, only where the operator said it lives;
// edit no file and write no work item, which flai guard enforces; ask the
// operator on a thread only when an entry needs their judgement; and end
// with a one-line summary naming the entries written, merged, absorbed, and
// compacted, or that nothing was found to take in, which flai serve logs as
// the run's activity, as the analyzer's is (ADR-0079).
func curatePrompt(r Request) string {
	harvest := fmt.Sprintf("Take in what the archive and the threads gained since %[1]s, your last run or the time the operator gave: list the candidates with flai digest harvest --json --since %[1]s.", r.Since)
	if r.Since == "" {
		harvest = "This is your first run: no run of yours is logged in wip/agents/curator.md, so take in everything the archive and the threads hold: list the candidates with flai digest harvest --json, with no --since."
	}
	return fmt.Sprintf(`You are %[1]s, the curator, started by flai serve on this host to curate the digest of the project at %[2]s: after a story's acceptance, when the digest's schedule came round, or because the operator asked.

Keep the digest, one fact per file under design/digest, and do nothing else, as design/conventions/strategic-agents.md says under As the curator. Before anything else, prime your session with the flai MCP tool prime with role curate (flai prime --role curate on the host), which prints the conventions you work by and briefs the design your role's topics select. Call the flai MCP tool inbox.

%[3]s Read the digest as it is with flai digest list --json, and an entry whole with flai digest show <id>, before you write.

For each candidate, decide what fact it holds. When no entry holds that fact, write one with flai digest new "<fact>" --kind <kind> --source <id>, with the topics and paths it bears on, and a body of a few sentences: what was asked or decided, the answer and its reason, and links to its sources. When an entry already holds it, merge the candidate into that entry with flai digest edit <id>, its body and its sources. When a thread asked what an entry already answered, add the thread to that entry's asked_again with flai digest edit <id> --asked-again. Each list flag of flai digest edit replaces that list, so give the entry's own with what you add. Pass a body on standard input with --body-file - and a heredoc in the shell, since you cannot write a file.

Then, in every run, look for absorption: a fact the design now records need not be primed twice. List the active entries with flai digest list --status active --json. For each, search the design with the flai MCP tool doc_search on the entry's title, and read what it finds with the flai MCP tool doc_get and its heading. When a design document, a convention, or an ADR states the entry's fact, mark the entry absorbed by it with flai digest absorb <id> --by <path>, the document's path from the project root, and --heading <heading> when one section states it. Absorb an entry only when a document states its fact, not when one merely touches on it; leave it active otherwise.

After the absorptions, run flai digest compact once: it records each entry's citers, as flai digest cited does, and compacts each active entry that no work item, narrative, or thread cited within digest.compact_after, keeping its title and emptying its body. It deletes nothing.

Never write a secret's value in an entry, only where the operator said it lives: flai digest refuses text with the shape of a key or token, and writes nothing.

Edit nothing else: no file, and no work item; write entries only with flai digest new, edit, absorb, cited, and compact. flai guard refuses Edit and Write on every file, flai digest seed, and any write to a work item: never work around a refusal, by another tool, another command, or the shell.

When an entry would need the operator's judgement, ask with the flai MCP tool thread_open on the item the candidate comes from, your recommended answer first, take in what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

End with a one-line summary that names by ID the entries you wrote, those you merged, those you absorbed, and those flai digest compact compacted, or says that you found nothing to take in: flai serve logs the run's activity in wip/agents/curator.md with it, with the run's cost, so you need not log it yourself with activity_log.`, r.Name, r.Root, harvest)
}

// definedPrompt is what a role the project defines is asked to do (S-0417,
// ADR-0207 §4): what its file's body says, which is its definition's prompt
// too, on the run's item when it has one, read with item_get, with the run's
// inputs in the order given; prime with its verb, and its item when it has
// one, and call inbox; keep to what its file allows, as flai guard holds it;
// ask the operator on a thread for an input they own that is missing; and
// end with a one-line summary, which flai serve logs as the run's activity
// in wip/agents/<name>.md, as the planner's is (ADR-0079).
func definedPrompt(r Request) string {
	f := r.RoleFile
	on, item, prime, names := "", "This run is on no item: it works the project as your role file says.", fmt.Sprintf("flai prime --role %[1]s (or the flai MCP tool prime with role %[1]s)", f.Role), ""
	ask := "the item it concerns"
	if r.Item != "" {
		kind := workitem.TypeOfID(r.Item)
		on = " on " + r.Item
		item = fmt.Sprintf("This run is on %[1]s: read it with the flai MCP tool item_get, and what it links with item_get and doc_get.", r.Item)
		prime = fmt.Sprintf("flai prime --role %[1]s --%[2]s %[3]s (or the flai MCP tool prime with role %[1]s and %[2]s %[3]s)", f.Role, kind, r.Item)
		names = fmt.Sprintf(" and what %s names whole", r.Item)
		ask = r.Item
	}
	inputs := "The run has no inputs."
	if len(r.Inputs) > 0 {
		each := make([]string, len(r.Inputs))
		for i, in := range r.Inputs {
			each[i] = fmt.Sprintf("- %s: %s", in.Name, strings.Join(strings.Fields(in.Value), " "))
		}
		inputs = "The run's inputs, in the order given:\n\n" + strings.Join(each, "\n")
	}
	return fmt.Sprintf(`You are %[1]s, the %[2]s role this project defines in %[3]s, started by flai serve on this host because the operator asked for a run of it%[4]s, in the project at %[5]s. Your role file says what you are and what you do:

%[6]s

%[7]s Prime your session with %[8]s, which prints the conventions you work by%[9]s, and briefs the design your role's topics select. A brief is not the document: read the section that bears on your work with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox.

%[10]s

Work in the main checkout. flai guard holds every call you make to what your role file allows: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on %[11]s, your recommended answer first, do what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

End with a one-line summary of what you did: flai serve logs the run's activity in wip/agents/%[2]s.md with it, with the run's cost, so you need not log it yourself with activity_log.`,
		r.Name, f.Name, f.Path, on, r.Root, roleBody(f.Body), item, prime, names, inputs, ask)
}

// roleBody is a role file's body as its definition's prompt (ADR-0207 §4):
// the baseline text and the project's additions, without the marker's line
// and the blank lines that followed it.
func roleBody(body string) string {
	var out []string
	skip := false
	for _, l := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		blank := strings.TrimSpace(l) == ""
		switch {
		case strings.TrimSpace(l) == roles.Marker:
			skip = len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == ""
			continue
		case skip && blank:
			continue
		}
		skip = false
		out = append(out, l)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

// sinceForms are the forms a curator's since takes, as flai digest harvest
// --since reads it: a day in UTC, or an RFC 3339 time.
var sinceForms = []string{time.DateOnly, time.RFC3339}

// sinceReads says whether since is in one of sinceForms.
func sinceReads(since string) bool {
	return slices.ContainsFunc(sinceForms, func(form string) bool {
		_, err := time.Parse(form, since)
		return err == nil
	})
}

// ReasonsEnv is the variable that holds why flai serve started the
// orchestrator, one reason a line (ADR-0209).
const ReasonsEnv = "FLAI_REASONS"

// roleEnv tells a strategic agent's session its role, which flai guard reads
// to hold it to its work: the planner's its item as well (S-0208), the
// orchestrator's, the analyzer's, and the curator's nothing more, since each
// works the whole project (S-0218, S-0223, S-0431). Nothing for an agent
// working a story. A role flai does not start, a strategic agent with a
// story of its own, a planner with no epic or story to plan, an
// orchestrator, an analyzer, or a curator given an item, an analyzer with a
// focus it does not take, a curator with a since that does not read, and a
// focus or a since for any other agent are refused, as are inputs for any
// agent but a role the project defines. A role the project defines has its
// own (definedEnv). The orchestrator's reasons are FLAI_REASONS, one a
// line, and reasons for any other agent are refused (ADR-0209), as are forms
// for prime that are not one for each reason (ADR-0217 §6): the prompt
// carries the forms, and the environment does not.
func roleEnv(r Request) ([]string, error) {
	if len(r.Reasons) > 0 && (r.RoleFile != nil || r.Role != conventions.RoleOrchestrate) {
		return nil, fmt.Errorf("only the orchestrator is started for reasons; start this agent without them")
	}
	if len(r.PrimeReasons) > 0 && len(r.PrimeReasons) != len(r.Reasons) {
		return nil, fmt.Errorf("the orchestrator was given %d forms for prime for %d reasons; give one for each reason, or none", len(r.PrimeReasons), len(r.Reasons))
	}
	if r.RoleFile != nil {
		return definedEnv(r)
	}
	switch {
	case r.Role == "" && r.Focus == "" && r.Since == "" && len(r.Inputs) == 0:
		return nil, nil
	case len(r.Inputs) > 0:
		return nil, fmt.Errorf("only a role a project defines takes inputs; start this agent without them")
	case r.Role != conventions.RoleAnalyze && r.Focus != "":
		return nil, fmt.Errorf("only the analyzer takes a focus; start this agent without %q", r.Focus)
	case r.Role != conventions.RoleCurate && r.Since != "":
		return nil, fmt.Errorf("only the curator takes a since; start this agent without %q", r.Since)
	case !slices.Contains(conventions.StrategicRoles, r.Role):
		return nil, fmt.Errorf("flai cannot start an agent in role %q; it starts the planner, role %s, the orchestrator, role %s, the analyzer, role %s, the curator, role %s, a role a project defines, from its role file, and an agent for a story", r.Role, conventions.RolePlan, conventions.RoleOrchestrate, conventions.RoleAnalyze, conventions.RoleCurate)
	case r.Role == conventions.RoleCurate && r.Story != "":
		return nil, fmt.Errorf("the curator was given story %s to work; it works no story of its own, so start it with none", r.Story)
	case r.Role == conventions.RoleCurate && r.Item != "":
		return nil, fmt.Errorf("the curator was given %s; it curates the whole project's digest, so start it with no item", r.Item)
	case r.Role == conventions.RoleCurate && r.Since != "" && !sinceReads(r.Since):
		return nil, fmt.Errorf("the curator takes a since as a YYYY-MM-DD day in UTC or an RFC 3339 time, or none on its first run, and %q is neither", r.Since)
	case r.Role == conventions.RoleCurate:
		return []string{"FLAI_ROLE=" + r.Role}, nil
	case r.Role == conventions.RoleAnalyze && r.Story != "":
		return nil, fmt.Errorf("the analyzer was given story %s to work; it works no story of its own, so start it with none", r.Story)
	case r.Role == conventions.RoleAnalyze && r.Item != "":
		return nil, fmt.Errorf("the analyzer was given %s; it analyzes the whole project, so start it with no item", r.Item)
	case r.Role == conventions.RoleAnalyze && r.Focus != "" && !slices.Contains(Focuses, r.Focus):
		return nil, fmt.Errorf("the analyzer takes the focus %s, or none for all of them, and %q is none of them", strings.Join(Focuses, ", "), r.Focus)
	case r.Role == conventions.RoleAnalyze:
		return []string{"FLAI_ROLE=" + r.Role}, nil
	case r.Role == conventions.RoleOrchestrate && r.Story != "":
		return nil, fmt.Errorf("the orchestrator was given story %s to work; it works no story of its own, so start it with none", r.Story)
	case r.Role == conventions.RoleOrchestrate && r.Item != "":
		return nil, fmt.Errorf("the orchestrator was given %s; it works the whole project, so start it with no item", r.Item)
	case r.Role == conventions.RoleOrchestrate && len(r.Reasons) > 0:
		return []string{"FLAI_ROLE=" + r.Role, ReasonsEnv + "=" + strings.Join(r.Reasons, "\n")}, nil
	case r.Role == conventions.RoleOrchestrate:
		return []string{"FLAI_ROLE=" + r.Role}, nil
	case r.Story != "":
		return nil, fmt.Errorf("the planner for %s was given story %s to work; it works no story of its own, so start it with none", r.Item, r.Story)
	}
	if k := workitem.TypeOfID(r.Item); k != workitem.Epic && k != workitem.Story {
		return nil, fmt.Errorf("the planner plans an epic or a story, and %q is neither; give it an E-nnnn or an S-nnnn", r.Item)
	}
	return []string{"FLAI_ROLE=" + r.Role, "FLAI_ITEM=" + r.Item}, nil
}

// definedEnv tells the session of a role the project defines, which flai
// guard reads to hold it to its file, who it is, FLAI_AGENT, RoleAgent of its
// name and item; its verb, FLAI_ROLE; and its item, FLAI_ITEM, when the run
// has one (S-0417, ADR-0207 §2). A file with no name, or a sub-agent's name,
// or with no verb, or the story's agent's or a sub-agent's; a run in another
// role than its file's verb; one under another name than RoleAgent's; one
// given a story to work, which flai starts no defined role for; one on an
// item its file does not take; and a focus or a since are refused.
func definedEnv(r Request) ([]string, error) {
	f := r.RoleFile
	agent := RoleAgent(f.Name, r.Item)
	kind := workitem.TypeOfID(r.Item)
	switch {
	case f.Name == "":
		return nil, fmt.Errorf("the role file %s names no role; give it a name, as flai check says", f.Path)
	case slices.Contains(roles.ReservedNames, f.Name):
		return nil, fmt.Errorf("the role file %s is named %s, a sub-agent's definition; give the role another name, as flai check says", f.Path, f.Name)
	case f.Role == "" || slices.Contains(roles.ReservedRoles, f.Role):
		return nil, fmt.Errorf("the role file %s gives %s no role of its own, %q; give it a verb that is not %s, as flai check says", f.Path, f.Name, f.Role, strings.Join(roles.ReservedRoles, ", "))
	case r.Role != f.Role:
		return nil, fmt.Errorf("the role %s runs in role %s, as its file %s says, and was started in role %q; start it in its own", f.Name, f.Role, f.Path, r.Role)
	case r.Focus != "":
		return nil, fmt.Errorf("only the analyzer takes a focus; start this agent without %q", r.Focus)
	case r.Since != "":
		return nil, fmt.Errorf("only the curator takes a since; start this agent without %q", r.Since)
	case r.Story != "":
		return nil, fmt.Errorf("the role %s was given story %s to work; flai starts a role a project defines on its item, never to work a story, so start it with none, or with the story as its item when its file takes one", f.Name, r.Story)
	case r.Item != "" && len(f.Item) == 0:
		return nil, fmt.Errorf("the role %s was given %s, and its file %s says it takes no item; start it with none", f.Name, r.Item, f.Path)
	case r.Item != "" && !slices.Contains(f.Item, kind):
		return nil, fmt.Errorf("the role %s takes %s, as its file %s says, and %q is none; give it an item of that kind", f.Name, strings.Join(f.Item, " or "), f.Path, r.Item)
	case r.Name != agent:
		return nil, fmt.Errorf("the role %s's run works under the name %s, and was started as %q; start it as %s", f.Name, agent, r.Name, agent)
	}
	env := []string{"FLAI_AGENT=" + agent, "FLAI_ROLE=" + f.Role}
	if r.Item != "" {
		env = append(env, "FLAI_ITEM="+r.Item)
	}
	return env, nil
}

// answeredSince asks the agent to read the threads written to since its
// story was begun, when there are any.
func answeredSince(b *Begun) string {
	switch len(b.Threads) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(" %s, on the story, has an entry by someone else since it was begun: read it with the flai MCP tool thread_get before you go on.", b.Threads[0])
	}
	return fmt.Sprintf(" Its threads %s have entries by someone else since it was begun: read them with the flai MCP tool thread_get before you go on.", strings.Join(b.Threads, ", "))
}

// delegation tells the agent when to hand work to the explorer and verifier
// sub-agents the template defines, what to give them, and what comes back
// (ADR-0059, design/conventions/delegation.md), and to plan its tasks in
// layers and hand each task to a task sub-agent, a layer at once when its
// tasks are long, whose work it reviews and closes itself with flai task done
// (S-0176, S-0197, S-0269, design/conventions/work-management.md), naming the
// task's ID in each task sub-agent's description so that flai measures the
// task by its calls (S-0230), and to sync once more before the close-out run
// (ADR-0069);
// and to run the close-out, which runs flai verify, itself, once, without a
// pipe or a file, read its last line, fix what it names, and run it again,
// handing the verifier only the review of a large diff against the criteria
// and the conventions, which reads flai verify --last rather than running the
// suite (S-0266, S-0270); and to make independent edits and commands in one turn, and move a
// task it has just written to ready and in-progress in one command (S-0268);
// and to ask each task sub-agent which acceptance criteria its task meets, and
// tick them after its review (S-0282, ADR-0089); and to launch every sub-agent
// in the foreground and never end its turn or hold wait_for_events while one
// runs, since the session ends ten minutes after its turn does (S-0285); and
// that a write under a .claude/ folder is never a sub-agent's but its own,
// made once the layer is back, and that permission_prompt refuses it after at
// most four minutes unanswered, leaving its thread open, so the agent goes on,
// makes the same write again once the thread is answered, and ends when
// nothing else is left (S-0299, S-0309, ADR-0124); and to run the tests
// for what it changed itself with flai test, or the MCP tool test, on the
// paths it changed, leaving the whole suite and lint to the close-out
// (S-0273, S-0270).
func delegation(r Request) string {
	return fmt.Sprintf(`Keep your own context for decisions and edits, and hand noisy work to sub-agents with the Agent tool: code, document, and log search across many files to the explorer, and the review of a large diff against the acceptance criteria and the conventions to the verifier. A sub-agent starts with nothing but your prompt: give it the worktree's path, %[1]s and the task's ID, the question, what you already know, and the shape of the answer you want, and ask for a summary with paths and lines, not raw output. Do it yourself when that is quicker: one file you know, one short command. A sub-agent cannot change work items or threads; a question it returns for the designer is yours to ask with thread_open. Make independent edits and commands in one turn, as several tool calls in one message: consecutive edits to one file, reads of files you already know, and commands that do not wait on each other.

Plan %[1]s's tasks as you write them: give each the paths it touches (--touches) and the tasks of %[1]s it must wait for (--after), so that tasks with no after between them and no path in common form layers that can run together, and record the layers and why each task waits in the narrative's Decisions. Move a task you have just written to ready and in-progress in one command, such as flai move T-nnnn ready && flai move T-nnnn in-progress: flai move refuses a task straight from backlog to in-progress. Work the plan layer by layer, handing each task to a task sub-agent with the Agent tool (a fork where you can), so that its reading stays out of your context; run a layer's tasks at once only when they are long, and then wait for all of them. Launch every sub-agent with the Agent tool's run_in_background set to false, a layer's in one message, so that each result comes back as the tool's result however long the sub-agent runs, and go on from there. Never end your turn while a sub-agent runs in the background: Claude Code ends this session ten minutes after the turn ends, and the sub-agent with it. Never wait for a sub-agent with the flai MCP tool wait_for_events either, which reports work items and threads, not sub-agents, and which flai guard refuses while one runs: wait_for_events is for a thread awaiting the designer. Name the task's ID in each task sub-agent's description, so that flai measures the task by its calls. A task sub-agent edits only what its task touches, runs only its own tests with flai test, and never commits, syncs, or writes through flai; ask each to name in its final message the acceptance criteria its task meets, by their numbers in flai criteria list %[1]s. A write to a path Claude Code protects, such as a file in a .claude/ folder or .mcp.json, is never a sub-agent's: flai guard refuses it unless the operator has turned on auto-approve, so say so in the prompt of every sub-agent whose task changes such a file, and ask it to return the file's whole new content in its final message. Make those writes yourself once the layer's sub-agents are back, never while a layer runs: permission_prompt opens a thread on %[1]s and holds the write at most four minutes, then refuses it, naming the thread, which stays open. Go on with the work that does not need the write, and once the thread is answered make the same write again, with the same content: an allow lets it through, a refusal refuses it, and with no answer yet it holds the write at most four minutes again. When nothing is left but the answer, write the narrative's Current state and Next steps with flai stream state, naming the thread, and end: flai serve starts you again when the thread is answered, and your first inbox holds the answer. Review each one's work yourself, fix what falls short, close the task with flai task done and test it as above, then tick the criteria your review verified it meets; only you close tasks, commit, sync the stream, move items, and talk to the designer.

While you work, run only the tests for what you changed, yourself, with flai test and the paths you changed, or the flai MCP tool test with them: it runs the test and lint tiers the manifest's tests declare that those paths select, cheapest first, and answers pass or the first findings. Do not run go test, vitest, golangci-lint, or gofmt by hand and read their logs, and do not hand that run to a sub-agent. Leave the whole suite, the whole lint, and flai check to the close-out before review. Before you move %[1]s to review, commit everything, run flai stream sync %[1]s again and resolve what it lists, then run the close-out yourself, once, in the worktree: the project's close-out script where it has one, which runs flai verify %[1]s --record-issues and then commits (it refuses a branch that does not contain the main branch), or flai verify %[1]s (the flai MCP tool verify) where it has none. Run it as one Bash call with a timeout of 600000 ms, its exit status echoed after it, such as scripts/close-out.sh %[1]s -m "<message>"; echo "exit $?", never piping its output or redirecting it into a file, and read its last line, which names the outcome and the step it stopped at, and the findings of the step that failed. When it stops, fix what it names yourself, commit, and run it again; never hand the run to a sub-agent. Its passing run is the story's run before review: do not repeat it, and read it again with flai verify %[1]s --last. Then check the diff against the acceptance criteria and the conventions, naming the criteria it meets by number. Hand that review to one fresh verifier only when the diff is too large to read here: tell it the run passed, and at which commit, and to read the result with flai verify %[1]s --last rather than running the suite. Fix what it finds yourself, never through a sub-agent, commit, and run the close-out again.`, r.Story)
}

// rules is what every run working a story is told: how its issues are
// recorded (S-0198), in one call that commits them on the story's branch and
// widens its touches, as its ADRs are (S-0275), how to ask the designer, and that with nothing left but
// the answer, or when wait_for_events answers end, it writes the narrative's
// state and ends, for flai serve starts it again on the answer (S-0272); that
// it ticks each acceptance criterion through flai once it has verified it
// (S-0282, ADR-0089); and how it goes to review or blocks.
func rules(r Request) string {
	return fmt.Sprintf(`Record friction, defects, and blockers you hit with flai issue new --commit, or flai issue bump --commit when the issue exists, run in the worktree (or the flai MCP tools issue_new and issue_bump with commit): each instance names %[1]s, and the one call commits the issue on %[1]s's branch and adds it to %[1]s's touches, so follow it with no git add, git commit, or flai touches. Record an ADR the same way, with flai adr new --commit (or the flai MCP tool adr_new with commit). Make no story for them yourself: %[1]s's review page lists them, checked, so the operator chooses at acceptance which become backlog stories. Make one with flai issue story or the flai MCP tool issue_story only when the operator asks, or when flai check warns issues.no-story.

When you need the designer to decide something, ask with the flai MCP tool thread_open on %[1]s, and go on with the work of %[1]s that does not wait on the answer. When nothing is left but the answer, write the narrative's Current state and Next steps with flai stream state, saying what you asked and what you will do with each answer, and end: flai serve starts you again in this session when the thread is answered, and your first inbox holds the answer. Do not hold the flai MCP tool wait_for_events for an answer; when it answers end: true, do as its why says: write the narrative's Current state and Next steps with flai stream state, and end. Tick each acceptance criterion with flai criteria tick %[1]s <n> (or the flai MCP tool criteria_tick) as soon as you have verified it, never by editing %[1]s's file; leave one you cannot verify here unticked and say why in %[1]s's notes. When every acceptance criterion is met, commit everything outstanding in the worktree, so that git status there is clean, run flai stream sync %[1]s again and resolve what it lists, then move %[1]s to review with flai move %[1]s review and end: the move is refused while anything is uncommitted. If you cannot go on, block the story with flai block %[1]s --reason and say why in its narrative, then end.`, r.Story)
}
