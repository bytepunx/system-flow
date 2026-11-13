package hostapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
)

// ActionPush is the host action that lets a dashboard publish (S-0078,
// ADR-0067): publish.run tags what was accepted since each component's last
// release and pushes the branch and the tags with the operator's own
// credentials. Publishing is the one way accepted work reaches the remote.
const ActionPush = "push"

// ActionForge is the host action under which flai serve pushes each story
// in review under pull-request acceptance to origin and opens or updates its
// pull request, with the host's git credentials and the token the host
// names for the forge (S-0393, ADR-0165). flai serve does it from what it
// sees in the project's files. Under it, too, the story's pull request is
// merged on the forge when the story is accepted (ADR-0170), and accept.run
// asks for it of a project under pull-request acceptance only. Off, which
// is the default, nothing is pushed or merged. The agents flai serve starts
// never hold the token, and flai guard refuses them flai forge push.
const ActionForge = "forge"

// ActionAutoPublish is the host action that has every flai push --pending,
// wherever it runs, first tag a release of everything accumulated (S-0094).
// Off, which is the default, a push tags nothing. Since ADR-0067 it and flai
// push --pending are the operator's shell tools, outside the workflow: no
// method asks for it, the dashboard neither sees nor changes it
// (DashboardActions leaves it out, settings.action refuses it), and flai
// push reads it.
const ActionAutoPublish = "auto-publish"

// ActionAutoApprove is the host action that has flai's MCP tool
// permission_prompt allow, without asking, a flai serve agent's Edit or Write
// of a path Claude Code protects, such as a file in a .claude/ folder or
// .mcp.json but never one in .git, in its own in-progress story's worktree,
// which Claude Code run headless refuses without a person's approval
// (S-0257, ADR-0106). On, it also lets that story's sub-agents make such
// writes, which flai guard otherwise refuses them at once (S-0299). A story
// whose branch changes such a path is accepted by the operator only, so
// nothing written this way is merged unread (S-0286). Off, which is the
// default, permission_prompt asks the story's owner on a thread and waits
// for allow. Like auto-publish it is the operator's shell tool (ADR-0067):
// no method asks for it, the dashboard neither sees nor changes it, and flai
// mcp reads it at every request.
const ActionAutoApprove = "auto-approve"

// ActionAutoUpgrade is the host action that has flai serve ask the host it
// runs under to install a newer flai release, and restart on it, when the
// project's history holds a newer flai/vX.Y.Z tag than the flai serving it
// (S-0387, ADR-0158). Off, which is the default, flai serve asks the
// operator on a thread instead. It replaces the binary on the operator's
// machine, so like auto-approve it is the operator's shell tool (ADR-0067):
// no method asks for it, and the dashboard neither sees nor changes it.
const ActionAutoUpgrade = "auto-upgrade"

// ActionAgent is the host action that starts a story's agent when the story
// becomes ready, the in-progress limit has room, and review is under its
// limit (S-0079, S-0104, ADR-0043, S-0243). flai serve performs that itself, from what it sees in the
// project's files. Four methods ask for it: agent.restart, a new agent for
// a story whose agent dropped or failed (S-0116), agent.start, a ready
// story's agent now (S-0115), agent.commit, an agent to commit what a
// story in review left uncommitted in its worktree (S-0140), and agent.stop,
// the story's agent ended (S-0170).
const ActionAgent = "agent"

// ActionDashboard is the host action that restarts, upgrades, or stops the
// dashboard container with Docker on this host (S-0081).
const ActionDashboard = "dashboard"

// ActionChecks is the host action that runs the commands named in the
// manifest or the host's configuration, in a story's worktree, and reports
// the outcome (S-0082); since S-0273 also the manifest's test tiers, through
// flai test.
const ActionChecks = "checks"

// ActionHost is the host action that has flai host start, stop, or restart
// flai serve and the MCP servers, and install the newest flai and restart on
// it (S-0106).
const ActionHost = "host"

// ActionSettings is the host action that lets a dashboard change the host's
// settings for a project (S-0105): the other host actions, the agent and its
// harnesses, the checks, the import folders, and the tokens. Those kept for
// every project need it enabled for every project. It is turned on and off
// in a shell on the host alone: no method changes it.
const ActionSettings = "settings"

// ActionPlan is the host action that starts the planner for an epic or a
// story on the operator's word (S-0208, ADR-0075): plan.run, from the item's
// page, and flai plan in a shell. flai serve never starts a planner by
// itself.
const ActionPlan = "plan"

// ActionOrchestrate is the host action under which flai serve runs one
// orchestrator per project for as long as it is on (S-0218): started when it
// is turned on, started again when a run ends, and stopped when it is turned
// off. flai serve sees it at every look. Two methods ask for it (S-0228):
// orchestrate.stop, the run ended and the orchestrator held stopped while it
// stays on, and orchestrate.start, the hold lifted and the orchestrator
// started again.
const ActionOrchestrate = "orchestrate"

// ActionAnalyze is the host action that starts the analyzer for a project,
// with a focus or none, on the operator's word (S-0223): analyze.run, and
// flai analyze in a shell; and on analysis.schedule, which flai serve starts
// for all three focuses. One analyzer runs per project at a time.
const ActionAnalyze = "analyze"

// ActionCurate is the host action that starts the curator for a project,
// from a since or from its last run, on the operator's word (S-0431,
// ADR-0196): curate.run, flai curate in a shell, and the MCP tool curate.
// One curator runs per project at a time; nothing starts it on a schedule.
const ActionCurate = "curate"

// ActionFormula is the host action under which flai serve makes a formula's
// instance each time the schedule its file names comes round (S-0422,
// ADR-0238): it runs flai formula run <name> --ready --autocommit in the
// project's main checkout, so that the instance's agent starts as any ready
// story's does. flai serve does it from what it sees in the project's files;
// no method asks for it. Off, which is the default, no schedule makes one.
const ActionFormula = "formula"

// ActionComponents is the host action under which flai serve runs a
// sub-project's own Vite dev server, with flai's component harness, on
// 127.0.0.1 (S-0438, ADR-0163): components.start and components.stop ask
// for it. Starting one runs the project's code and dependencies on the host.
const ActionComponents = "components"

// ActionMergeQueue is the host action under which flai serve runs the merge
// queue, a loop of its own with no model (S-0399, ADR-0173): while
// orchestration.permissions.accept_reviews is on and acceptance.mode is
// local, it batches at most merge_queue.batch_max enqueued eligible stories
// in a worktree of its own and verifies the batch once. flai serve does it
// from what it sees in the project's files. Off, which is the default, no
// batch starts.
const ActionMergeQueue = "merge-queue"

// ActionKluster is the host action under which flai serve runs a project's
// kluster jobs, with no model, and the dashboard's Ops pages bring the
// project's cluster up or down and deploy, tear down, or promote through
// flai (S-0577, ADR-0254). flai kluster in a shell works whatever it says.
// Off, which is the default, no job starts and every namespace is left as it
// is; turning it off tears nothing down. flai serve enable refuses it for a
// project whose manifest has no kluster: block.
const ActionKluster = "kluster"

// Actions are the host actions there are, with what each lets a dashboard
// do, or, for one no dashboard sees (DashboardSees), what it has flai do.
var Actions = map[string]string{
	ActionPush:        "publish when you press Publish: tag what was accepted since each component's last release and push the branch and the tags, with your git credentials, never forced; a holder of the dashboard token can then publish everything accepted and unreleased",
	ActionForge:       "push each story that moves to review, while acceptance.mode in system-flow.yaml is pull-request, to origin with your git credentials, and open its one pull request against the main branch, or update the one it has, with the forge's token from the variable flai serve forge names (GITHUB_TOKEN for GitHub by default, ADR-0165), never forcing a branch but a story's own that a sync rebased, with --force-with-lease; a push or a request that fails is said on a thread on the story and tried again at the next look. The agents flai serve starts never get the token's variable, and flai guard refuses them flai forge push. Accepting a story merges its pull request on the forge by acceptance.merge_method and deletes its branch on origin (ADR-0170): the dashboard's acceptance runs flai accept here, with the token's variable, and an agent's flai accept asks flai serve to merge. Off, or with acceptance.mode local, nothing is pushed or merged, and under pull-request acceptance every acceptance is refused",
	ActionAutoPublish: "has flai push --pending, run in a shell, first tag a release of everything merged and unreleased since each component's last tag; off, pushing releases nothing and what is accepted waits for Publish or flai release --pending. It is the operator's shell tool, outside the workflow (ADR-0067): no dashboard sees or changes it",
	ActionAutoApprove: "has flai's MCP tool permission_prompt allow, without asking you, a flai serve agent's Edit or Write of a path Claude Code protects, such as a file in a .claude/ folder (its settings, hooks, and agent definitions) or .mcp.json but never .git, in its own in-progress story's worktree, which Claude Code run headless refuses without a person's approval, and lets that story's sub-agents make those writes too, which flai guard otherwise refuses them; off, permission_prompt asks the story's owner on a thread and waits for allow. A story whose branch changes such a path is accepted by you only, so you read those changes at acceptance (ADR-0106). It is the operator's shell tool (ADR-0067): no dashboard sees or changes it",
	ActionAutoUpgrade: "has flai serve ask flai host to install the newest flai release tagged in this project's history, verified as flai host upgrade --version installs it, and restart flai serve and the MCP servers on it, whenever the flai serving the project is older than that release, checked every minute; it waits while a dashboard write runs in flai serve, and asks again for a refused release no sooner than 15 minutes later; off, flai serve asks you on a thread on the project's system-flow.yaml, once for each newer release. It replaces the flai on this host, so it is the operator's shell tool (ADR-0067, ADR-0158): no dashboard sees or changes it",
	ActionAgent:       "start each story's agent, with the harnesses and the command you set with flai serve agent, on this machine and as you, whenever a story becomes ready, the in-progress limit has room, and review is under its limit, start a ready story's agent on demand, start or queue a new one for a story whose agent dropped or failed, start one to commit what a story in review left uncommitted in its worktree, and stop a story's agent, ending its process and everything it started; whoever can move a story to ready or press Start agent, Retry, Have an agent commit them, or Stop, a holder of the dashboard token included, then starts or stops it",
	ActionDashboard:   "restart the dashboard container, upgrade it to the image your configuration names, or stop it, with Docker on this host; an upgrade is never applied until the new image answers healthy, so a bad one leaves the running container untouched",
	ActionChecks:      "run the commands named in flai serve checks set or the manifest's checks:, and the test tiers the manifest's tests: declare through flai test, in a story's worktree or the main checkout, on this host, and cancel a checks run; whoever can open the review page then decides what runs there",
	ActionHost:        "have flai host start, stop, or restart flai serve and the MCP servers of every project on this host, and download the newest flai release with your GitHub credentials, install it over the flai on this host, and restart everything on it",
	ActionSettings:    "change this project's host settings: turn the other host actions on and off, set its default agent, and rotate its MCP token; enabled for every project, also the agent's command, the harnesses, the providers, the checks, the import folders, and the dashboard token. A holder of the dashboard token can then run any command on this host, as you; only a shell turns this off",
	ActionPlan:        "start the planner, with the project's planning agent (planning.agent over agent in system-flow.yaml) and the harnesses and the command you set with flai serve agent, on this machine and as you, in the project's main checkout, for an epic or a story when you press Plan or run flai plan; it writes work items and threads through flai and moves nothing past backlog; a holder of the dashboard token can then start it for any epic or story not done or cancelled",
	ActionOrchestrate: "run the orchestrator, with the project's orchestration agent (orchestration.agent over agent in system-flow.yaml) and the harnesses and the command you set with flai serve agent, on this machine and as you, in the project's main checkout, for as long as this is on: started again when it ends, a minute after a failure, and stopped when this is turned off; it keeps work moving through flai only as far as orchestration.permissions in system-flow.yaml allow, each off by default, and logs each decision in wip/agents/orchestrator.md; a holder of the dashboard token can then stop it and hold it stopped, and start it again, from its page",
	ActionComponents:  "start a Svelte sub-project's own Vite dev server, from its node_modules, with its own config, plugins, and dependencies and flai's component harness, on 127.0.0.1 and a port flai picks, on this machine and as you, when the Components page asks, in the main checkout or in the worktree of a story in progress or in review, whose dependencies flai first installs by the sub-project's lock file when they are missing, and stop it; at most components.max_servers (3 by default) run at once per project; flai serve stops one that nothing asked about for components.idle (30 minutes by default), a story's once the story leaves review, and every one when it stops. A holder of the dashboard token can then run the sub-project's code, and its branch's, on this host",
	ActionMergeQueue:  "run the merge queue in flai serve, with no model, on this machine and as you, while orchestration.permissions.accept_reviews in system-flow.yaml is on and acceptance.mode is local: it takes up to merge_queue.batch_max (4 by default) eligible stories enqueued in review, in queue order, rebases each one's branch onto the one before it from the main branch, in a worktree of its own under .flai-cache/worktrees, and verifies the batch once with the tiers flai verify runs, running the project's tests on this host; a story whose rebase stops is left out of the batch. No batch starts while flai queue hold holds the queue. Off, or under pull-request acceptance, no batch starts and enqueued stories wait",
	ActionAnalyze:     "start the analyzer, with the project's analysis agent (analysis.agent over agent in system-flow.yaml) and the harnesses and the command you set with flai serve agent, on this machine and as you, in the project's main checkout, when you run flai analyze or call analyze.run, looking for bottlenecks, intent, or risk, or all three, and for all three each time analysis.schedule in system-flow.yaml comes round; it writes one report under design/analysis and edits nothing else, and one runs per project at a time; a holder of the dashboard token can then start it whenever none runs",
	ActionCurate:      "start the curator, with the project's digest agent (digest.agent over agent in system-flow.yaml) and the harnesses and the command you set with flai serve agent, on this machine and as you, in the project's main checkout, when you run flai curate or call curate.run; it takes in what the archive and the threads gained since the day or time you give, or since its last run, as digest entries written through flai digest new and edit alone, edits no file, authors no work item, and leaves its entries uncommitted, and one runs per project at a time; a holder of the dashboard token can then start it whenever none runs",
	ActionKluster:     "run the project's local Kubernetes cluster, the one the kluster: block of its system-flow.yaml names, made by bytepunx/kluster with k3d and Docker, on this machine and as you: deploy each story that enters review into a namespace of its own with the project's own build, deploy, and teardown commands, tear it down when the story leaves review, deploy main into the default namespace on each acceptance, and bring the cluster up when a deploy finds it absent, but never take it down; and let the dashboard bring the cluster up or take it down, and deploy, tear down, or promote, through flai. Whoever can change those commands on a story's branch then runs them on this host. Off, no job starts and every namespace stays as it is; flai kluster in a shell works either way",
	ActionFormula:     "make a formula's instance each time the schedule its file under design/workflows/formulas names comes round, by running flai formula run <name> --ready --autocommit in the project's main checkout, on this machine and as you, with no inputs, so that the instance is committed in ready and its story's agent starts as any ready story's does; a formula flai check refuses is not made and is tried again when its schedule next comes round. Whoever can change a formula's file on the main branch then decides what runs on this host, and when",
}

// shellOnly are the host actions kept to the operator's shell (ADR-0067):
// no dashboard sees them or turns them on and off.
var shellOnly = map[string]bool{ActionAutoPublish: true, ActionAutoApprove: true, ActionAutoUpgrade: true}

// DashboardSees reports whether a dashboard is told of the host action name.
func DashboardSees(name string) bool {
	_, ok := Actions[name]
	return ok && !shellOnly[name]
}

// DefinedActions are the host actions the roles the project at root defines
// name in their files under the design folder's workflows/roles (S-0417,
// ADR-0207), each with what it lets a dashboard do: start the roles that
// name it, through role.run or flai role run. A built-in role's file, a file
// that does not load, and an action that is one of Actions name none; a
// project whose manifest or role files cannot be read defines none.
func DefinedActions(root string) map[string]string {
	byAction := map[string][]roles.Role{}
	var order []string
	for _, f := range definedRoles(root) {
		if len(byAction[f.Action]) == 0 {
			order = append(order, f.Action)
		}
		byAction[f.Action] = append(byAction[f.Action], f)
	}
	out := map[string]string{}
	for _, action := range order {
		each := make([]string, 0, len(byAction[action]))
		for _, f := range byAction[action] {
			each = append(each, fmt.Sprintf("start the role %[1]s, which %[2]s defines, with its agent (roles.%[1]s.agent over agent in system-flow.yaml) and the harnesses and the command you set with flai serve agent, on this machine and as you, in the project's main checkout, when you run flai role run %[1]s or call role.run, %[3]s; one run of it %[4]s at a time; a holder of the dashboard token can then start it", f.Name, f.Path, roleTakes(f), roleRunsOn(f)))
		}
		out[action] = strings.Join(each, "; and ")
	}
	return out
}

// definedRoles are the roles the project at root defines in role files that
// load, the built-in roles' files and those naming an action of Actions, or
// none, left out.
func definedRoles(root string) []roles.Role {
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return nil
	}
	all, _, err := roles.Load(root, m.Layout["design"])
	if err != nil {
		return nil
	}
	var out []roles.Role
	for _, f := range all {
		if _, fixed := Actions[f.Action]; f.Action == "" || fixed || builtinRole(f.Name) {
			continue
		}
		out = append(out, f)
	}
	return out
}

// builtinRole reports whether name is a built-in role's, which flai starts
// by commands of its own.
func builtinRole(name string) bool {
	for _, n := range roles.BuiltinActions {
		if n == name {
			return true
		}
	}
	return false
}

// roleAction is the host action the role name's file names in the project
// at root, which role.run is behind (S-0417), or "" when the project
// defines no role name: flai role run then refuses it in its own words.
func roleAction(p channel.Project, raw json.RawMessage) string {
	var in struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(raw, &in)
	for _, f := range definedRoles(p.Root) {
		if f.Name == in.Name {
			return f.Action
		}
	}
	return ""
}

// roleTakes says what a role's run is on, as its file says.
func roleTakes(f roles.Role) string {
	if len(f.Item) == 0 {
		return "on no item"
	}
	return "on a " + strings.Join(f.Item, " or ") + " not archived, done, or cancelled"
}

// roleRunsOn says how many runs of a role go at once: one per item, or one.
func roleRunsOn(f roles.Role) string {
	if len(f.Item) == 0 {
		return "for the project"
	}
	return "per item"
}

// Host is what the host decides and records about host actions (ADR-0029).
// The zero value enables nothing and records nothing.
type Host struct {
	// Enabled reports whether the operator enabled an action for a project.
	Enabled func(action, root string) bool
	// Record writes one entry of the host's journal.
	Record func(Entry)
	// Agent reports what flai serve knows of agents it starts for a project
	// (S-0079): what runs, the last start or failure, and why a ready story
	// waits. ctx is the request's, which carries its recorder, so the phases
	// of the state it builds are timed (S-0603). finished asks for the runs
	// of finished work and the past runs too; without it the state holds the
	// open and live runs alone (S-0611). Nil when nothing starts agents here.
	Agent func(ctx context.Context, root string, finished bool) any
	// AgentStream reads the stream of the newest agent flai serve started for
	// a story (S-0142), from a byte offset, or its tail when after is
	// negative; ErrNoAgent when it started none. Nil when nothing starts
	// agents here.
	AgentStream func(root, story string, after int64) (any, error)
	// PlanStream reads the stream of the newest planner flai serve started
	// for an epic or a story (S-0259), as AgentStream reads a story's agent's;
	// ErrNoAgent when it started none. Nil when nothing starts agents here.
	PlanStream func(root, item string, after int64) (any, error)
	// RoleStream reads the stream of the newest run of a strategic role flai
	// serve started for the project, as AgentStream reads a story's agent's:
	// the orchestrator's for orchestrate (S-0218), the analyzer's for analyze
	// (S-0228), and for any role's name the newest run of that role, on
	// whatever item (S-0420); ErrNoAgent, wrapped, when it started none. Nil
	// when nothing starts agents here.
	RoleStream func(root, role string, after int64) (any, error)
	// Settings reports the host's settings as they apply to a project
	// (S-0105), for the dashboard's settings page, recording its phases on
	// the context (S-0605). Nil when there are none.
	Settings func(ctx context.Context, root string) any
	// Requests keeps the records of detached writes where a restart of flai
	// serve does not lose them (S-0109). Nil keeps them in memory.
	Requests Requests
	// Components runs the sub-projects' component dev servers in flai
	// serve's own process (S-0438, ADR-0163). Nil when nothing runs them
	// here: components.status says none runs, and components.start refuses.
	Components ComponentServers
	// Kluster runs each project's kluster up and kluster down in flai serve's
	// own process (S-0580, ADR-0254). Nil when nothing runs them here:
	// kluster.status reports no job, and kluster.up and kluster.down refuse.
	Kluster *KlusterJobs
	// Images reads flai serve's image jobs and the project's agent images,
	// and starts a job (S-0757). Nil when nothing runs them here: the images
	// methods refuse.
	Images ImageJobs
}

// Entry is one host action asked for, whatever became of it.
type Entry struct {
	At        string `json:"at"`
	Action    string `json:"action"`
	Method    string `json:"method"`
	Project   string `json:"project"`
	Root      string `json:"root"`
	By        string `json:"by"` // who the dashboard acts for: the manifest's owner
	RequestID string `json:"request_id,omitempty"`
	Outcome   string `json:"outcome"` // done, failed, or disabled
	Detail    string `json:"detail,omitempty"`
	// DurationMS is how long a close-out job ran, and Verdict what it came
	// to, set on its end entry (ADR-0229 §6).
	DurationMS int64  `json:"duration_ms,omitempty"`
	Verdict    string `json:"verdict,omitempty"`
}

func (h Host) enabled(action, root string) bool {
	return h.Enabled != nil && h.Enabled(action, root)
}

// EnableCommand is what the operator runs, on the host, in the project.
func EnableCommand(action string) string { return "flai serve enable " + action }
