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
	"errors"
	"fmt"
	"os"
	"path"
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

// None is the reserved harness of an agent that runs elsewhere, by hand or in
// a container no harness starts, for which flai serve starts nothing
// (ADR-0242). It is no adapter.
const None = "none"

// Elsewhere reports whether an agent runs elsewhere: whether its harness is
// None.
func Elsewhere(a *manifest.Agent) bool {
	return a != nil && a.Harness == None
}

// ElsewhereRefusal says why flai serve starts nothing for a story whose agent
// runs elsewhere, and what to do instead; story is its ID, or empty when the
// caller has none.
func ElsewhereRefusal(story string) string {
	who, start := "the story", "flai story start <story>"
	if story != "" {
		who, start = story, "flai story start "+story
	}
	return fmt.Sprintf("%s's agent runs elsewhere (harness none), so flai serve starts nothing for it: start it yourself with %s, or name a harness flai serve can start", who, start)
}

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
	// start request: pilot loads it by RoleName from the mounted checkout
	// (S-0705).
	RoleFile *roles.Role `json:"-"`
	// RoleName is the name of a defined role's file, which pilot's start
	// request carries in place of the file (S-0705); empty for any other
	// agent, and on the host, where RoleFile says it.
	RoleName string `json:"role_name,omitempty"`
	// Builtin is the role file a built-in role runs from, the project's or
	// the baseline flai carries (roles.Builtin): its body, the project's
	// additions below the marker included, is the prompt of the definition
	// claude-code passes it (S-0419, ADR-0207 §4). Nil for any other agent,
	// and for a built-in role whose caller gives none, whose definition keeps
	// the prompt of its .claude/agents file. Never read from a start request.
	Builtin *roles.Role `json:"-"`
	// Inputs are a defined role's run's inputs, in the order given, which
	// its prompt names (S-0417); none for any other agent. Pilot's start
	// request carries them, since the checkout it mounts does not (S-0705).
	Inputs []Input `json:"inputs,omitempty"`
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
	// Formula is the formula the story is an instance of, its formula front
	// matter, which flai formula run wrote (S-0421): its agent works the
	// story's tasks as the formula's steps (S-0422). Empty for a story made
	// any other way, and for any other agent.
	Formula string
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
	// Probe is set for a hand-run trial of an agent image (I-0224): its agent
	// only reads and reports, and works no item, not even Story, whose mounts
	// it runs in. Only a start request scripts/agent-container.sh --probe
	// writes sets it, never flai serve.
	Probe bool `json:"probe,omitempty"`
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
// request (ADR-0167), and on the other failures its fields name, as the
// story's agent is told it.
type SentBack struct {
	Thread string // the thread flai serve opened on the story for it
	What   string // what failed, as the move's reason names it
	// Resumed is set when the agent goes on in its last run's session,
	// Request.Session; a new session is told how to take the story up.
	Resumed bool
	// Queue is set when the merge queue sent the story back, from a batch
	// that failed with it, rather than its pull request (ADR-0176).
	Queue bool
	// Acceptance is set when an acceptance, whoever made it, sent the story
	// back over a rebase it could not settle (ADR-0248).
	Acceptance bool
	// Kluster is set when flai serve sent the story back because its deploy
	// into its kluster namespace, or its checks run against it, failed
	// (ADR-0254 §5a).
	Kluster bool
}

// then is what follows the agent's fix of a story sent back on thread: the
// push of its pull request, from the merge queue its batching again, from an
// acceptance its acceptance again, or from its kluster namespace its deploy
// again in place.
func (b *SentBack) then(story string) string {
	if b.Kluster {
		return fmt.Sprintf("flai serve deploys %[1]s again in place once it is back in review, at a new commit or the same one, and resolves %[2]s once that deploy and its checks pass. When the cause is outside the story, such as the cluster's services, say so on %[2]s rather than change the branch for it", story, b.Thread)
	}
	if b.Acceptance {
		return fmt.Sprintf("flai serve resolves %[2]s once %[1]s is back in review, and it is accepted again", story, b.Thread)
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
	// Images are the container harness's image per inner harness, by the
	// inner harness's name; DefaultImages where it names none.
	Images map[string]string // container: the image per inner harness
	// Network is the container harness's Docker network; bridge when empty.
	Network string
	// Login is the folder of the operator's login to the inner harness, which
	// the container harness mounts read-only at ContainerLogin.
	Login string
	// DenyProtected is the operator's say that their arguments deny the
	// harness's protected paths (ADR-0131), from the host's deny_protected
	// for it: the opencode adapter then writes a deny rule for each into
	// OpenCode's permission map (ADR-0260).
	DenyProtected bool
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
	// Stdin is what the starter writes to the child's standard input, then
	// closes: pilot's start request, for the container harness. Empty for
	// every other adapter.
	Stdin []byte
	// Container is the container a container harness's start runs: its
	// image, its name, and its mounts, as its Argv gives them; nil for every
	// other adapter (I-0274).
	Container *ContainerRun
}

// ContainerRun is the container a start runs, read from its docker run: the
// image, the container's name, and each -v mount in the order the argv gives
// them. It carries no secret: the login mount names its path alone.
type ContainerRun struct {
	Image  string           `json:"image"`
	Name   string           `json:"name"`
	Mounts []ContainerMount `json:"mounts,omitempty"`
}

// ContainerMount is one -v mount of a container: the host's path, the path
// in the container, and whether it is read-only.
type ContainerMount struct {
	Host     string `json:"host"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"read_only,omitempty"`
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

// SessionReader is an adapter whose harness makes its own session IDs, so
// that flai serve learns from a run's log the session a later run resumes.
type SessionReader interface {
	// Session is the ID of the session the run that wrote log ran in, and
	// whether log names one.
	Session(log []byte) (string, bool)
}

// Adapters are the harnesses flai can start, by name.
var Adapters = map[string]Adapter{
	ClaudeCode: claudeCode{},
	Command:    command{},
	Container:  container{},
	OpenCode:   openCode{},
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

// For is the adapter a story's agent is started with, and its name. An
// agent that names no harness is started with the operator's command, when
// one is set, and otherwise with flai's built-in default, claude-code
// (ADR-0261); one whose agent runs elsewhere is refused. A container's
// adapter is the inner harness's its config.inner names, as Capabilities and
// Reader answer, and an inner harness it cannot run is refused.
func For(a *manifest.Agent, commandSet bool) (string, Adapter, error) {
	if Elsewhere(a) {
		return "", nil, errors.New(ElsewhereRefusal(""))
	}
	name := ""
	if a != nil {
		name = a.Harness
	}
	if name == "" {
		name = manifest.BuiltinHarness
		if commandSet {
			name = Command
		}
	}
	ad, ok := Adapters[name]
	if !ok {
		return "", nil, errors.New(UnknownRefusal("the agent", name))
	}
	if name == Command && !commandSet {
		return "", nil, errors.New(CommandUnsetRefusal("the agent"))
	}
	if name == Container {
		inner, err := innerOf(a)
		if err != nil {
			return "", nil, err
		}
		ad = container{inner: inner}
	}
	return name, ad, nil
}

// HarnessesPage and AgentsPage are the dashboard's settings pages that set
// the host's command and harness programs, and the project's default agent,
// by their paths, which the dashboard renders as links (ADR-0261 § 4).
const (
	HarnessesPage = "/settings/harnesses"
	AgentsPage    = "/settings/agents"
)

// CommandUnsetRefusal is the refusal for whose agent, such as "S-0042" or
// "the planner", whose harness is the host's command when none is set for
// the project: what is missing, the page that sets it, and the command
// (ADR-0261 § 4).
func CommandUnsetRefusal(whose string) string {
	return fmt.Sprintf("%s's harness is the host's command, and none is set for this project: set it on the dashboard's Harnesses page, %s, or with flai serve agent set -- <program> [args...]", whose, HarnessesPage)
}

// UnknownRefusal is the refusal for whose agent naming a harness flai cannot
// start: the harnesses it can, the page that sets the project's default
// agent, and the command (ADR-0261 § 4).
func UnknownRefusal(whose, name string) string {
	return fmt.Sprintf("flai cannot start harness %q; it can start %s: name one for %s on the dashboard's Agents page, %s, or with flai agent set --harness <name>", name, strings.Join(Names(), ", "), whose, AgentsPage)
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
		if !ok && len(takes) == 0 {
			return nil, fmt.Errorf("%s takes no config, and the agent sets %q; leave config out of the story's agent", harness, k)
		}
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
// shared). One whose story is a formula's instance is told to work its tasks
// as the formula's steps and what to do when one fails (S-0422, instance),
// and, answered, to run a failed step again (answeredStep). A strategic
// role, the planner, the orchestrator, the analyzer, the curator, or a role
// the project defines, is asked what its role file's start prompt says,
// rendered over the run's fields (startPrompt, ADR-0250). Its error is a
// start prompt that does not render, which stops the start. A probe, a trial
// of an agent image, is asked only to read and report, before any other case
// and whatever its story or role (probe, I-0224).
// Only the claude-code and opencode adapters send a prompt.
func Prompt(r Request) (string, error) {
	if r.Probe {
		return probe(r), nil
	}
	if r.RoleFile != nil {
		return startPrompt(r)
	}
	switch r.Role {
	case conventions.RolePlan, conventions.RoleOrchestrate, conventions.RoleAnalyze, conventions.RoleCurate:
		return startPrompt(r)
	}
	if r.Commit != "" && r.Answered == "" {
		return fmt.Sprintf(`You are %[1]s, started by flai serve on this host because the operator asked for the work left uncommitted in story %[2]s's worktree, %[3]s, to be committed: %[2]s is in review, and it cannot be accepted until that worktree is clean.

Do only this. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md, for how commits are made here. In the worktree, read git status and git diff. Commit the changes on story/%[2]s, in commits whose messages name %[2]s and say what changed and why, after the lint and tests the project runs for what they touch; fix what they find only when it is part of the same work. Discard a file only when it is plainly output that does not belong in the repository, and say which in the story's narrative with flai stream log %[2]s. Do not move %[2]s or its tasks, do not change anything the uncommitted work does not already change, and do not start other work. When git status in the worktree is clean, log what you committed with flai stream log %[2]s and end.

If a change cannot be committed without the designer deciding something, ask with the flai MCP tool thread_open on %[2]s, and commit what does not wait on the answer. Then write the narrative's Current state and Next steps with flai stream state, saying what you asked and what is left to commit, and end: flai serve starts you again in this session when the thread is answered, and your first inbox holds the answer. Do not hold the flai MCP tool wait_for_events for an answer; when it answers end: true, write the narrative's Current state and Next steps with flai stream state, and end.`, r.Name, r.Story, r.Commit), nil
	}
	if strings.HasPrefix(r.Answered, "MS-") {
		return fmt.Sprintf(`Another story's agent has written to %[2]s on %[3]s. Read it with the flai MCP tool inbox, which lists %[2]s's conversations under messages, or message_get (flai message show %[3]s on the host), answer what awaits your reply with message_reply, then go on working %[2]s to review as before.

%[4]s`, r.Name, r.Story, r.Answered, rules(r)), nil
	}
	if r.Answered != "" {
		return fmt.Sprintf(`The designer has answered your question %[3]s on %[2]s. Read the answer with the flai MCP tool thread_get (or flai thread show %[3]s), then go on working %[2]s to review as before.%[5]s

%[4]s`, r.Name, r.Story, r.Answered, rules(r), answeredStep(r)), nil
	}
	if b := r.SentBack; b != nil && b.Resumed {
		return fmt.Sprintf(`flai sent %[1]s back from review to in-progress, on thread %[2]s: %[3]s. Read the thread with the flai MCP tool thread_get (flai thread show %[2]s on the host); your inbox lists it, awaiting you. Fix the branch in %[1]s's worktree, close out, and move %[1]s to review again; %[5]s.

%[4]s`, r.Story, b.Thread, b.What, rules(r), b.then(r.Story)), nil
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

Work %[2]s to review, and no other story. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md. %[7]s A brief is not the document: when one bears on the story, read it, or its section that does, with the flai MCP tool doc_get and its heading (flai doc show --heading on the host) before relying on it or changing what it describes, and find sections by their words with doc_search. Write the story's tasks if it has none, or review the ones the planner drafted, and work them in that worktree. %[10]s Keep the narrative's Current state and Next steps true: rewrite them at every task transition with flai stream state %[2]s --current "<text>" --next "<text>" (or the flai MCP tool stream_state), never by editing the narrative.%[9]s

%[6]s

%[4]s`, r.Name, r.Story, r.Root, rules(r), why, delegation(r), open, shared(r), instance(r), closeTasks(r)), nil
}

// inContainer says the agent runs inside the container harness's container:
// pilot's start request names the inner harness it runs there (S-0507), and a
// start on the host names none. A story's agent there has the main checkout's
// .git read-only (ADR-0246), so its prompt closes tasks with the MCP tool
// task_done and the story with verify, whose commit and sync run on the host,
// in place of flai stream sync, the shell's flai task done, and git rebase,
// and moves items with item_move (I-0272).
func (r Request) inContainer() bool { return r.Inner != "" }

// closeTasks tells a story's agent how it closes each task: with flai task
// done in the worktree, or the MCP tool task_done, which syncs as flai stream
// sync does, resolving a sync's conflicts with git rebase --continue. In the
// container (Request.inContainer), whose main checkout's .git is read-only
// (ADR-0246), only task_done, whose commit and sync run on the host, closes a
// task, and a sync's conflicts go to the designer, since neither git add nor
// git rebase can write there (I-0272).
func closeTasks(r Request) string {
	if r.inContainer() {
		return fmt.Sprintf(`When a task is done, close it with the flai MCP tool task_done, with its docs and work-item updates in the change: it commits on story/%[1]s, syncs %[1]s's branch, moves the task to done, logs it in the narrative, widens the touches, runs flai check, and answers your inbox, stopping at the first step that fails. Its commit and its sync run on the host. When the sync stops on conflicts, ask the designer with the flai MCP tool thread_open on %[1]s, naming each path it lists, and go on with the work that does not wait on the answer; when the check stops, fix what it found and call it again. Then run the task's tests with flai test and the paths it changed (or the flai MCP tool test with them), and close any fix they need by calling task_done again. You run in a container, where the main checkout's .git is read-only: flai stream sync, the shell's flai task done, and git rebase cannot write there, so never run them, and never start a rebase or merge by hand.`, r.Story)
	}
	return fmt.Sprintf(`When a task is done, close it with flai task done T-nnnn -m "<message>" in the worktree (or the flai MCP tool task_done), with its docs and work-item updates in the change: it commits on story/%[1]s, runs flai stream sync %[1]s, moves the task to done, logs it in the narrative, widens the touches, runs flai check, and answers your inbox, stopping at the first step that fails. When the sync stops on conflicts, resolve each path it lists in the worktree, git add it, and git rebase --continue, then call it again; when the check stops, fix what it found and call it again. Then run the task's tests with flai test and the paths it changed (or the flai MCP tool test with them), and close any fix they need by calling flai task done again. flai stream sync does the branch's git work and refuses while anything is uncommitted: never start a rebase or merge by hand.`, r.Story)
}

// probe is what a trial of an agent image is asked: to read, report, and end.
func probe(r Request) string {
	story := ""
	if r.Story != "" {
		story = fmt.Sprintf(" You run in the mounts of story %[1]s, but %[1]s is not yours to work: do not start, take up, or change it.", r.Story)
	}
	return fmt.Sprintf(`You are %[1]s, started by hand as a trial of the agent image, not as any story's agent.%[2]s Do no work on any item: start, move, create, or edit none, and write to no thread, narrative, or file.

Check only this. Call the flai MCP tools that read, each once: inbox, board, item_get on an item the board lists, doc_get on design/conventions/README.md, and doc_search for a word in it. Then run flai version, git status, and git log -1 in your working directory. Never call the flai MCP tools that write: story_start, stream_state, task_done, item_new, item_edit, item_move, thread_open, thread_reply, or verify.

Report one line per check: what you called or ran, and what you saw, or the error it gave. Then end.`, r.Name, story)
}

// instance tells the agent of a formula's instance, a story flai formula run
// made (S-0421), as a paragraph of its own, how to work it (S-0422, TH-0471):
// its tasks are the formula's steps, already written, worked in their layers
// and in no other order, each as its task's Work names it, by the agent
// itself; a flai command or a script under scripts/ is run from the worktree
// and its output read, and a role step is started with role_run on the
// instance, the agent holding wait_for_events until the run ends; a task is
// closed with task_done when its step exits 0, or its role run ends worked,
// with the exit in its Notes; the end of a role run it started, a change of
// kind agent on the instance whose summary says a story's agent may let it
// pass, is the step's exit and not let pass (T-1775); a step that exits with
// an exit its stop lists stops the run cleanly: the agent runs flai formula
// stop with the step's output and takes the instance to review, asking
// nothing (S-0621, ADR-0271); a step that fails, a command or script that
// exits non-zero with no exit its stop lists or a role run that ends failed
// or asked, runs nothing after it: its task is blocked with its exit, the exit
// and the last 40 lines of its output go in its Notes, a thread on the story
// names the step and the exit, and the agent writes its Current state and
// Next steps and ends, as for any question (T-1776, answeredStep); the story
// goes to review, closed out with verify, when every task is done; and no
// step is invented or reworded. Empty for a story that is no instance, whose
// prompt is then as before. In the container (Request.inContainer) it names
// task_done alone, since the shell's flai task done cannot write there
// (I-0272).
func instance(r Request) string {
	if r.Formula == "" {
		return ""
	}
	shell := ` (flai task done T-nnnn -m "<message>" in the worktree)`
	if r.inContainer() {
		shell = ""
	}
	return fmt.Sprintf(`

%[1]s is an instance of the formula %[2]s, which flai formula run made from design/workflows/formulas/%[2]s.md: its tasks are the formula's steps, written already, so write, review, add, rewrite, or cancel none of them. Work them in their layers, a layer only once every task of the one before is done, and in no other order, each as its task's ## Work names it, and run each step yourself rather than through a task sub-agent. Run a flai command or a script under scripts/ from the root of %[1]s's worktree, and read its output. Start a role step with the flai MCP tool role_run (flai role run on the host), with the role the step names, the step's prompt and inputs, and %[1]s as the run's item, then hold the flai MCP tool wait_for_events until the run ends. When a step exits 0, or its role run ends worked, write its exit in the task's ## Notes and close the task with the flai MCP tool task_done%[3]s. The end of a role run you started for a step reaches wait_for_events as a change of kind agent on %[1]s by the role's agent: to failed when the run failed, and to ended when it asked, its why waiting for an answer to a thread, or when it worked, its why ended with its work done. That change is the step's exit: do not let it pass, though its summary says a story's agent may. When a command or script exits with an exit its formula lists in the step's stop, as its task's ## Done when names it, the step stops the run cleanly: run nothing after it, and run flai formula stop %[1]s --step <id> --exit <n> in the worktree with the step's output on standard input, which closes its task and cancels the steps after it. Then close %[1]s out with verify and move it to review, asking nothing and opening no thread. When a step fails, a command or script that exits other than 0 and with no exit its stop lists, or a role run that ends failed or asked, run nothing after it. Block its task with flai block T-nnnn --reason "<the step and its exit>" in the worktree, write the exit and the last 40 lines of the step's output in the task's ## Notes with the flai MCP tool item_edit, or by editing the task's body, and open a thread on %[1]s with the flai MCP tool thread_open naming the step and its exit. Then write the narrative's Current state and Next steps with the flai MCP tool stream_state, and end the session, as for any question: flai serve starts you again when the thread is answered. When every task is done, close %[1]s out with verify and move it to review as any story. Invent no step, and change no step's words.`, r.Story, r.Formula, shell)
}

// answeredStep tells the agent of a formula's instance, started again on the
// answer to its question (T-1776), that when it asked because a step failed
// it unblocks that step's task, runs the step again as the answer says, and
// goes on with the layers. Empty for a story that is no instance, whose
// answered prompt is then as before.
func answeredStep(r Request) string {
	if r.Formula == "" {
		return ""
	}
	return fmt.Sprintf(` %[1]s is an instance of the formula %[2]s: if you asked %[3]s because a step failed, unblock its task with flai unblock T-nnnn in the worktree, run that step again as the answer says, and go on with the layers from there, closing its task with task_done when it passes.`, r.Story, r.Formula, r.Answered)
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

// startPrompt is what a strategic role is asked to do (ADR-0250): the start
// prompt of the role file it runs from, rendered over the run's fields
// (startFields), with the project's Start prompt additions below the
// marker. The file is RoleFile for a role the project defines, Builtin for
// a built-in role whose caller gave it, else the baseline flai carries for
// the role. A file with no start renders the baseline's for a built-in role
// and flai's default for a defined one (roles.Role.StartTemplate). A
// template that does not render is an error naming the file and flai
// check, never a prompt.
func startPrompt(r Request) (string, error) {
	f := r.RoleFile
	if f == nil {
		f = r.Builtin
	}
	if f == nil {
		b, ok := roles.BaselineFor("design", r.Role)
		if !ok {
			return "", fmt.Errorf("flai carries no role file for role %q, and the start gave none; start a built-in role, plan, orchestrate, analyze, or curate", r.Role)
		}
		f = &b
	}
	return f.RenderStart(designOf(*f), startFields(r, *f))
}

// designOf is the design folder the role file f is under, as its path from
// the project root, <design>/workflows/roles/<name>.md, gives it; the
// request carries no design folder of its own. Only the baseline start a
// file with none falls back to reads it, which is the same under any.
func designOf(f roles.Role) string {
	dir := path.Dir(f.Path)
	if dir == roles.Folder {
		return "."
	}
	return strings.TrimSuffix(dir, "/"+roles.Folder)
}

// startFields are the run's fields r's start prompt renders over, from the
// role file f it runs from (ADR-0250 §2): the item's kind from its ID, each
// reason with its form for prime when every reason has one, and then the
// forms each once in the order first given, and each input's whitespace
// folded to single spaces.
func startFields(r Request, f roles.Role) roles.StartFields {
	s := roles.StartFields{
		Name: r.Name, Root: r.Root, Role: f.Role, RoleName: f.Name, File: f.Path, Body: roles.Prompt(f.Body),
		Item: r.Item, Kind: workitem.TypeOfID(r.Item), AskedOn: r.AskedOn,
		Primed: len(r.Reasons) > 0 && len(r.PrimeReasons) == len(r.Reasons),
		Queue:  r.Queue, Focus: r.Focus, Since: r.Since,
	}
	for i, line := range r.Reasons {
		reason := roles.StartReason{Line: line}
		if s.Primed {
			reason.Prime = r.PrimeReasons[i]
			if !slices.Contains(s.PrimeForms, reason.Prime) {
				s.PrimeForms = append(s.PrimeForms, reason.Prime)
			}
		}
		s.Reasons = append(s.Reasons, reason)
	}
	for _, in := range r.Inputs {
		s.Inputs = append(s.Inputs, roles.StartInput{Name: in.Name, Value: strings.Join(strings.Fields(in.Value), " ")})
	}
	return s
}

// Focuses are what the analyzer can be asked to look for (S-0223):
// bottlenecks in the flow, gaps between the design and the code (intent),
// and technical and security risks; the analyzer's start prompt names what
// each looks for. A run asked for none looks for all three, and its
// report's focus is AllFocus.
var Focuses = []string{"bottlenecks", "intent", "risk"}

// AllFocus is the focus of a report from a run asked for no focus.
const AllFocus = "all"

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
// and to ask flai serve for the close-out job with the MCP tool verify (flai
// verify --job on the host), hold wait_for_events until it reports verified
// for the story, never running the close-out script in a turn, read the
// verdict with flai verify --last, fix what it names, and ask again, running
// the close-out itself, once, without a pipe or a file, only when verify says
// flai serve is not running (S-0546, ADR-0229); handing the verifier only the
// review of a large diff against the criteria and the conventions, which reads
// flai verify --last rather than running the suite (S-0266, S-0270); and to make independent edits and commands in one turn, and move a
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
// (S-0273, S-0270). In the container (Request.inContainer) it closes tasks with
// task_done alone and syncs nothing itself before the close-out, whose commit
// and sync run on the host (ADR-0274, I-0272).
func delegation(r Request) string {
	closer := "flai task done"
	preReview := fmt.Sprintf("commit everything, run flai stream sync %[1]s again and resolve what it lists, then ask flai serve for the close-out with the flai MCP tool verify,", r.Story)
	mover := "in one command, such as flai move T-nnnn ready && flai move T-nnnn in-progress: flai move refuses"
	if r.inContainer() {
		closer = "task_done"
		preReview = "close every task with task_done, then ask flai serve for the close-out with the flai MCP tool verify, whose commit and sync run on the host,"
		mover = "with the flai MCP tool item_move, to ready and then to in-progress, since flai move cannot take the items lock in the container: item_move refuses"
	}
	return fmt.Sprintf(`Keep your own context for decisions and edits, and hand noisy work to sub-agents with the Agent tool: code, document, and log search across many files to the explorer, and the review of a large diff against the acceptance criteria and the conventions to the verifier. A sub-agent starts with nothing but your prompt: give it the worktree's path, %[1]s and the task's ID, the question, what you already know, and the shape of the answer you want, and ask for a summary with paths and lines, not raw output. Do it yourself when that is quicker: one file you know, one short command. A sub-agent cannot change work items or threads; a question it returns for the designer is yours to ask with thread_open. Make independent edits and commands in one turn, as several tool calls in one message: consecutive edits to one file, reads of files you already know, and commands that do not wait on each other.

Plan %[1]s's tasks as you write them: give each the paths it touches (--touches) and the tasks of %[1]s it must wait for (--after), so that tasks with no after between them and no path in common form layers that can run together, and record the layers and why each task waits in the narrative's Decisions. Move a task you have just written to ready and in-progress %[4]s a task straight from backlog to in-progress. Work the plan layer by layer, handing each task to a task sub-agent with the Agent tool (a fork where you can), so that its reading stays out of your context; run a layer's tasks at once only when they are long, and then wait for all of them. Launch every sub-agent with the Agent tool's run_in_background set to false, a layer's in one message, so that each result comes back as the tool's result however long the sub-agent runs, and go on from there. Never end your turn while a sub-agent runs in the background: Claude Code ends this session ten minutes after the turn ends, and the sub-agent with it. Never wait for a sub-agent with the flai MCP tool wait_for_events either, which reports work items and threads, not sub-agents, and which flai guard refuses while one runs: wait_for_events is for a thread awaiting the designer or the close-out job's verified event. Name the task's ID in each task sub-agent's description, so that flai measures the task by its calls. A task sub-agent edits only what its task touches, runs only its own tests with flai test, and never commits, syncs, or writes through flai; ask each to name in its final message the acceptance criteria its task meets, by their numbers in flai criteria list %[1]s. A write to a path Claude Code protects, such as a file in a .claude/ folder or .mcp.json, is never a sub-agent's: flai guard refuses it unless the operator has turned on auto-approve, so say so in the prompt of every sub-agent whose task changes such a file, and ask it to return the file's whole new content in its final message. Make those writes yourself once the layer's sub-agents are back, never while a layer runs: permission_prompt opens a thread on %[1]s and holds the write at most four minutes, then refuses it, naming the thread, which stays open. Go on with the work that does not need the write, and once the thread is answered make the same write again, with the same content: an allow lets it through, a refusal refuses it, and with no answer yet it holds the write at most four minutes again. When nothing is left but the answer, write the narrative's Current state and Next steps with flai stream state, naming the thread, and end: flai serve starts you again when the thread is answered, and your first inbox holds the answer. Review each one's work yourself, fix what falls short, close the task with %[2]s and test it as above, then tick the criteria your review verified it meets; only you close tasks, commit, sync the stream, move items, and talk to the designer.

While you work, run only the tests for what you changed, yourself, with flai test and the paths you changed, or the flai MCP tool test with them: it runs the test and lint tiers the manifest's tests declare that those paths select, cheapest first, and answers pass or the first findings. Do not run go test, vitest, golangci-lint, or gofmt by hand and read their logs, and do not hand that run to a sub-agent. Leave the whole suite, the whole lint, and flai check to the close-out before review. Before you move %[1]s to review, %[3]s with %[1]s and its commit message (flai verify %[1]s --job -m "<message>" on the host): it answers the close-out job's id at once, and the job runs the project's close-out script in the worktree, which runs flai verify %[1]s --record-issues and then commits (it refuses a branch that does not contain the main branch), or flai verify %[1]s where it has none. Never run the close-out script in a turn: hold the flai MCP tool wait_for_events, again each time it returns, until it reports verified for %[1]s, then read the verdict and the findings of the step that failed with flai verify %[1]s --last. When the verdict is not passed, fix what it names yourself, commit, and ask again; never hand the close-out to a sub-agent. When verify answers from the record, that is the passing run, at a commit already verified. Only when verify says flai serve is not running, run the close-out yourself, once, in the worktree: as one Bash call with a timeout of 600000 ms, its exit status echoed after it, such as scripts/close-out.sh %[1]s -m "<message>"; echo "exit $?", never piping its output or redirecting it into a file, and read its last line, which names the outcome and the step it stopped at, and the findings of the step that failed. When it stops, fix what it names yourself, commit, and run it again. The passing run, the job's or yours, is the story's run before review: do not repeat it, and read it again with flai verify %[1]s --last. Then check the diff against the acceptance criteria and the conventions, naming the criteria it meets by number. Hand that review to one fresh verifier only when the diff is too large to read here: tell it the run passed, and at which commit, and to read the result with flai verify %[1]s --last rather than running the suite. Fix what it finds yourself, never through a sub-agent, commit, and ask for the close-out again.`, r.Story, closer, preReview, mover)
}

// rules is what every run working a story is told: how its issues are
// recorded (S-0198), in one call that commits them on the story's branch and
// widens its touches, as its ADRs are (S-0275), how to ask the designer, and that with nothing left but
// the answer, or when wait_for_events answers end, it writes the narrative's
// state and ends, for flai serve starts it again on the answer (S-0272); that
// it ticks each acceptance criterion through flai once it has verified it
// (S-0282, ADR-0089); and how it goes to review or blocks: in the container
// (Request.inContainer), once task_done and verify have committed and synced
// on the host, never with flai stream sync (I-0272).
func rules(r Request) string {
	move := fmt.Sprintf("flai move %[1]s review", r.Story)
	ready := fmt.Sprintf("commit everything outstanding in the worktree, so that git status there is clean, run flai stream sync %[1]s again and resolve what it lists,", r.Story)
	if r.inContainer() {
		ready = fmt.Sprintf("close every task with the flai MCP tool task_done and %[1]s out with the flai MCP tool verify, which commit and sync on the host, so that git status in the worktree is clean,", r.Story)
		move = "the flai MCP tool item_move"
	}
	return fmt.Sprintf(`Record friction, defects, and blockers you hit with flai issue new --commit, or flai issue bump --commit when the issue exists, run in the worktree (or the flai MCP tools issue_new and issue_bump with commit): each instance names %[1]s, and the one call commits the issue on %[1]s's branch and adds it to %[1]s's touches, so follow it with no git add, git commit, or flai touches. Record an ADR the same way, with flai adr new --commit (or the flai MCP tool adr_new with commit). Make no story for them yourself: %[1]s's review page lists them, checked, so the operator chooses at acceptance which become backlog stories. Make one with flai issue story or the flai MCP tool issue_story only when the operator asks, or when flai check warns issues.no-story.

When you need the designer to decide something, ask with the flai MCP tool thread_open on %[1]s, and go on with the work of %[1]s that does not wait on the answer. When nothing is left but the answer, write the narrative's Current state and Next steps with flai stream state, saying what you asked and what you will do with each answer, and end: flai serve starts you again in this session when the thread is answered, and your first inbox holds the answer. Do not hold the flai MCP tool wait_for_events for an answer; when it answers end: true, do as its why says: write the narrative's Current state and Next steps with flai stream state, and end. Tick each acceptance criterion with flai criteria tick %[1]s <n> (or the flai MCP tool criteria_tick) as soon as you have verified it, never by editing %[1]s's file; leave one you cannot verify here unticked and say why in %[1]s's notes. When every acceptance criterion is met, %[2]s then move %[1]s to review with %[3]s and end: the move is refused while anything is uncommitted. If you cannot go on, block the story with flai block %[1]s --reason and say why in its narrative, then end.`, r.Story, ready, move)
}
