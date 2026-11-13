package serve

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting agents when stories become ready (S-0079, S-0104, ADR-0029).
//
// flai serve sees the project's files change. It starts an agent for each
// story in ready that has none, in pull order, if all of this holds: the
// agent action is enabled for the project; no agent has been started for the
// story since it entered ready, whether that was before or after flai serve
// began to serve it (S-0112: serve/agents.json remembers each story's runs
// across a restart), or its agent has been changed since the last one was
// started (S-0116); the in-progress limit leaves room, counting an agent
// started for a story still in ready as a story in progress; and review is
// under its limit (S-0243), so that review does not grow without end while
// the operator batches acceptance. Each ready story it does not start, and
// has no agent running, is named in the state's waiting with the reason, and
// logged when its reason changes.
//
// The story says which harness works it, with which model and options (its
// agent, S-0103); package harness turns that into a command, with the
// operator's program and arguments for the harness. A story that names no
// harness is started with the operator's command, when one is set.
//
// Nobody attending holds a ready story back any more (S-0116, ADR-0043): the
// in-progress limit does, a full review does (S-0243), and so does a claim
// (S-0128, ADR-0046). A ready story whose claim overlaps the claim of a
// story in progress, in review, or in ready with an agent started for it is
// held: it is skipped, keeps its place, and is started first once it is
// clear. An agent holding wait_for_work and the launcher may both go for a
// story; the second move to in-progress is refused, and the loser pulls the
// next one.
//
// Each command is run as it stands in the project's directory, never through
// a shell. The story's ID is flai's own reading of the files, checked against
// the form of an ID before it is put anywhere.

// AgentConfig is the operator's say about starting agents, read afresh at
// every look so that enabling, disabling, and a new command need no restart.
type AgentConfig struct {
	Enabled bool
	// Plan is whether the plan host action is enabled for the project: the
	// operator may have the planner started for an epic or a story (S-0208).
	Plan bool
	// Orchestrate is whether the orchestrate host action is enabled for the
	// project: flai serve runs its orchestrator while it is (S-0218).
	Orchestrate bool
	// Analyze is whether the analyze host action is enabled for the project:
	// the operator may have the analyzer started for it (S-0223).
	Analyze bool
	// Curate is whether the curate host action is enabled for the project:
	// the operator may have the curator started for it (S-0431, ADR-0196).
	Curate bool
	// RoleActions are whether each host action that a role the project
	// defines names in its file is enabled for the project, by action: the
	// operator may have the role started while its action is (S-0417,
	// ADR-0207). An action not here is off.
	RoleActions map[string]bool
	Command     []string
	Name        string // FLAI_AGENT is this and the story's ID; "agent" when empty
	// Harnesses are the program and arguments each harness runs with on this
	// host, the command among them when one is set (S-0104).
	Harnesses map[string]harness.Host
	// Providers are the host's providers, the configuration's
	// agent.providers: where the model calls of an agent that names one go,
	// with the project's manifest's overrides laid over (ADR-0129).
	Providers map[string]manifest.Provider
	// Flai is this flai's executable: the agents reach flai's MCP server
	// through it.
	Flai string
	// AutoRestarts is how many times in a row flai serve restarts a story's
	// agent that ended with its story in progress, on its own; 0 never does
	// (S-0294, ADR-0108).
	AutoRestarts int
	// Lifts are what the operator lifts on each harness on this host, the
	// command among them, by harness name (ADR-0131).
	Lifts map[string]HarnessLifts
	// AutoApprove is whether the auto-approve host action is enabled for the
	// project: a protected write is allowed without asking (ADR-0086).
	AutoApprove bool
	// ForgeTokenEnvs are the variables the host reads each forge's token
	// from, as flai serve forge names them (ADR-0165): no agent flai serve
	// starts gets one (S-0393).
	ForgeTokenEnvs []string
}

// HarnessLifts are what the operator lifts on one harness, taking the risk
// that comes with it (ADR-0131).
type HarnessLifts struct {
	// GuardNone lets a story's roles run without flai guard on a harness
	// that cannot run it.
	GuardNone bool
	// DenyProtected is the operator's say that their arguments for the
	// harness deny the protected paths, so a harness that cannot ask before
	// a protected write may start a story's agent.
	DenyProtected bool
}

// startAllowed says whether a story's agent, with agent's roles, may start
// on harness name, which can do what caps says, under cfg's lifts
// (ADR-0131), and, when its roles then run without flai guard, the words
// that say so in the start's journal entry and log line. A harness
// that cannot run the guard is refused a story with roles unless the host
// sets guard none for it; one that cannot ask before a protected write is
// refused unless the project's auto-approve host action is on or the host
// sets deny_protected for it. Each refusal names what lifts it, and both
// are given when both hold.
func startAllowed(name string, caps harness.Capabilities, agent *manifest.Agent, cfg AgentConfig) (unguarded string, err error) {
	lifts := cfg.Lifts[name]
	var refused []string
	if !caps.Guard && agent != nil && len(agent.Roles) > 0 {
		roles := strings.Join(slices.Sorted(maps.Keys(agent.Roles)), ", ")
		if lifts.GuardNone {
			unguarded = fmt.Sprintf("its roles (%s) run without flai guard (guard: none on harness %s, ADR-0131)", roles, name)
		} else {
			refused = append(refused, fmt.Sprintf("harness %[1]s cannot run flai guard, so the story's roles (%[2]s) cannot run; flai serve agent harness %[1]s --guard none runs them unguarded (ADR-0131)", name, roles))
		}
	}
	if !caps.Ask && !cfg.AutoApprove && !lifts.DenyProtected {
		refused = append(refused, fmt.Sprintf("harness %[1]s cannot ask before a protected write, so flai serve does not start a story's agent on it; turn on the project's auto-approve host action (flai serve enable auto-approve), or say your arguments for it deny the protected paths with flai serve agent harness %[1]s --deny-protected (ADR-0131)", name))
	}
	if len(refused) > 0 {
		return "", errors.New(strings.Join(refused, "; and "))
	}
	return unguarded, nil
}

// host is what harness name runs with: the operator's, or for the command
// the command itself.
func (c AgentConfig) host(name string) harness.Host {
	if h, ok := c.Harnesses[name]; ok {
		return h
	}
	if name == harness.Command && len(c.Command) > 0 {
		return harness.Host{Program: c.Command[0], Args: c.Command[1:]}
	}
	return harness.Host{}
}

// commandSet is whether the host's command is set for the project.
func (c AgentConfig) commandSet() bool { return c.host(harness.Command).Program != "" }

// storyAgent is the agent flai serve starts for the story it (ADR-0261): its
// own, the project's default under it (ADR-0251), the host's command when
// neither names a harness and one is set, and otherwise flai's built-in
// default, claude-code.
func (c AgentConfig) storyAgent(repo *workitem.Repo, it *workitem.Item) *manifest.Agent {
	return repo.StartAgentOf(it, c.commandSet())
}

// startAgent is the agent that a, a story's resolved by AgentOf or a
// strategic role's, starts with (ADR-0261): a itself, the host's command when
// it names no harness and one is set, and otherwise flai's built-in default
// under it.
func (c AgentConfig) startAgent(a *manifest.Agent) *manifest.Agent {
	if c.commandSet() {
		return a
	}
	return a.OrBuiltin()
}

// refuseHarness refuses to start the story id's agent, resolved as
// storyAgent resolves it, when its agent runs elsewhere (ADR-0242) or
// refuseUnstartable refuses it; nil when it can be started.
func (c AgentConfig) refuseHarness(id string, agent *manifest.Agent) error {
	if harness.Elsewhere(agent) {
		return refused("%s", harness.ElsewhereRefusal(id))
	}
	return c.refuseUnstartable(id, agent)
}

// refuseUnstartable refuses to start whose agent, a story's ID or a role
// such as "the planner", when it names the host's command and none is set
// for the project, or names a harness flai cannot start, each refusal naming
// the dashboard page and the command that fix it (ADR-0261 § 4); nil
// otherwise.
func (c AgentConfig) refuseUnstartable(whose string, agent *manifest.Agent) error {
	name := ""
	if agent != nil {
		name = agent.Harness
	}
	if name == harness.Command && !c.commandSet() {
		return refused("%s", harness.CommandUnsetRefusal(whose))
	}
	if _, ok := harness.Adapters[name]; name != "" && name != harness.None && !ok {
		return refused("%s", harness.UnknownRefusal(whose, name))
	}
	return nil
}

// Outcomes of an agent that has ended.
const (
	OutcomeWorked = "worked" // it left its story in review or done
	OutcomeAsked  = "asked"  // it ended waiting for an answer or a reply, and is started again when there is one
	OutcomeFailed = "failed" // it could not be started, or left its story anywhere else
	// OutcomeStopped is an agent the operator stopped (S-0170): its story is
	// left where it was, and gets no agent until a retry or a move to ready.
	OutcomeStopped = "stopped"
)

// AgentRun is one agent flai serve started, or failed to: a story's agent,
// the planner for an item (S-0208), the project's orchestrator (S-0218), or
// its analyzer (S-0223), or its curator (S-0431).
type AgentRun struct {
	Story string `json:"story"`
	// Item is the epic or story a planner run plans; empty on a story's
	// agent's run, whose story is Story, and Story is empty on a planner's.
	// Both are empty on the orchestrator's run, the analyzer's, and the
	// curator's, which work the whole project.
	Item    string `json:"item,omitempty"`
	Harness string `json:"harness,omitempty"`
	Model   string `json:"model,omitempty"`
	Command string `json:"command"` // the program's name, not its arguments
	Agent   string `json:"agent"`
	// Provider is the provider the run's agent named, which its model calls
	// went through, and KeyEnv the variable of flai serve's environment its
	// key was copied from, never the key (ADR-0129); both empty when it named
	// none, and on runs before S-0350.
	Provider string `json:"provider,omitempty"`
	KeyEnv   string `json:"key_env,omitempty"`
	// BudgetKey is what the run's budget state is kept under (budgetKeyOf),
	// such as "claude-code meclod", resolved before its start's gate; empty
	// on runs before S-0720, whose state is kept under budgetKey.
	BudgetKey string `json:"budget_key,omitempty"`
	PID       int    `json:"pid,omitempty"`
	// Start is when the process started as the system counts it (clock
	// ticks since boot on Linux), which tells it from another given the same
	// PID since (S-0170); 0 where it is not known, and on runs before it.
	Start   int64  `json:"start,omitempty"`
	Started string `json:"started"`
	Ended   string `json:"ended,omitempty"`
	Exit    *int   `json:"exit,omitempty"`
	Error   string `json:"error,omitempty"` // why it could not be started
	Log     string `json:"log,omitempty"`
	// Session is the harness's session, which a start after an answer resumes.
	Session string `json:"session,omitempty"`
	// Answered is what this run was started again for: the thread whose
	// question was answered, or the conversation a reply or a new message to
	// its story came on (S-0335).
	Answered string `json:"answered,omitempty"`
	// SentBack is the thread flai serve opened when flai sent the story back
	// from review to in-progress, when this run was started for it (S-0394,
	// ADR-0167).
	SentBack string `json:"sent_back,omitempty"`
	// Trigger is what started a planner run (ADR-0084): asked, for the
	// operator's asking, orchestrator, for the orchestrator's (S-0219), and
	// otherwise the replanner's triggers, separated by semicolons. It is empty on a story's agent's run and on planner runs
	// before S-0211. An analyzer run's and a curator run's are what its
	// starter said (S-0223, S-0431).
	Trigger string `json:"trigger,omitempty"`
	// Role is the mark of a run that the fields that name its work do not
	// tell from the orchestrator's: conventions.RoleCurate on the curator's
	// runs (S-0431, ADR-0196), so that one is never taken for the
	// orchestrator's run, and the verb its file gives on a run of a role the
	// project defines (S-0417); empty on every other run.
	Role string `json:"role,omitempty"`
	// RoleName is the name of the role the project defines whose run this
	// is, as its file under the design folder's workflows/roles gives it,
	// and Action the host action its file names, which its start and its
	// end are journalled under (S-0417, ADR-0207). Both are empty on every
	// other run.
	RoleName string `json:"role_name,omitempty"`
	Action   string `json:"action,omitempty"`
	// Since is when a curator run takes in from, as the harness was given it:
	// the since its starter gave, else the time of the curator's last
	// activity, and empty on its first run, which takes in everything
	// (S-0431).
	Since string `json:"since,omitempty"`
	// Entries are the digest entries a curator run wrote or changed, by ID,
	// once it has ended: those whose files under the design folder's digest/
	// changed since it started; empty when it changed none (S-0431).
	Entries []string `json:"entries,omitempty"`
	// Focus is what an analyzer run looks for (S-0223): one of
	// harness.Focuses, or harness.AllFocus when it was asked for none. It is
	// set on the analyzer's runs alone, and tells them from the
	// orchestrator's.
	Focus string `json:"focus,omitempty"`
	// Report is the report an analyzer run wrote, relative to the project's
	// root, once it has ended: the newest file under design/analysis changed
	// since it started; empty when it wrote none (S-0223).
	Report string `json:"report,omitempty"`
	// Instance is the formula's instance an analyzer run was started for, by
	// its agent, for a role step of its formula that names the analyzer: its
	// agent events are on the instance rather than on the analyzer's role
	// (S-0424). It is empty on every other run, and leaves Item empty, so that
	// the run is still the analyzer's (analyzes).
	Instance string `json:"instance,omitempty"`
	// Outcome is set once it has ended; Why says what went wrong, or what it
	// waits for, and Thread is the question to the operator it ended waiting
	// on.
	Outcome string `json:"outcome,omitempty"`
	Why     string `json:"why,omitempty"`
	Thread  string `json:"thread,omitempty"`
	// Conversations are the conversations with other stories' agents a run
	// that ended asking ended on (S-0335): those awaiting the reply of the
	// agent of a story WaitsOn names, and those holding a message to its
	// story that came while it ran and that it left unanswered. WaitsOn are
	// the stories whose agents' replies it waits on, each once, in the order
	// of their conversations.
	Conversations []string `json:"conversations,omitempty"`
	WaitsOn       []string `json:"waits_on,omitempty"`
	// StoryAgent is the story's agent as it was when this run was started,
	// the project's default under it when it named no harness (ADR-0251),
	// empty when neither named anything, and nil on runs recorded before
	// S-0116: a story whose agent has changed since is started again.
	StoryAgent *manifest.Agent `json:"story_agent,omitempty"`
	// Queued is when the operator asked for another agent after this one
	// ended, for its story in ready while the in-progress limit was full:
	// flai serve starts it as soon as there is room (S-0118).
	Queued string `json:"queued,omitempty"`
	// Stopped is when the operator stopped it (S-0170): however its process
	// ends from then, it ended stopped.
	Stopped string `json:"stopped,omitempty"`
	// Held is set on the orchestrator's run alone, while the operator holds
	// the orchestrator stopped (S-0228): flai serve does not start it again
	// until they start it, or turn the orchestrate action off and on.
	Held bool `json:"held,omitempty"`
	// AutoRestarts is how many automatic restarts in a row led to this run
	// (S-0294, ADR-0108): flai serve's restart on its own carries its last
	// run's count and one more, its start after the usage limit its last
	// run's count as it was (S-0534), and any other start begins again at 0.
	AutoRestarts int `json:"auto_restarts,omitempty"`
	// LimitThread is the thread flai serve opened to the operator when this
	// run ended at the limit of automatic restarts, so that it opens one
	// (S-0294, ADR-0108).
	LimitThread string `json:"limit_thread,omitempty"`
	// ResetsAt is when the harness's usage limit a run ended on resets, as an
	// RFC 3339 time in UTC; empty when it ended on no limit, or the harness
	// named no reset time (S-0514, ADR-0214): a curator run's since S-0514,
	// every run's since S-0534.
	ResetsAt string `json:"resets_at,omitempty"`
	// UsageLimited marks a run whose log ends on the harness's usage limit
	// (S-0534, ADR-0209 §5): its end set or extended its provider's budget
	// state, and it is no failed run that flai serve restarts on its own
	// (agent.auto_restarts) or starts again after orchestrateRetry.
	UsageLimited bool `json:"usage_limited,omitempty"`
	// Retry marks a run started again after the run before it ended on the
	// usage limit (S-0514, ADR-0214).
	Retry bool `json:"retry,omitempty"`
	// HarnessVersion is the version of the inner harness a container run's
	// image holds, as the harness prints it, such as 2.1.312 for Claude Code
	// (S-0508); empty on every other run, and when it could not be read.
	HarnessVersion string `json:"harness_version,omitempty"`
	// Image, Container, and Mounts are what a container run's docker run
	// was given: the image it ran, the project's layer image when the
	// manifest names one, its container's name, and its mounts in argv
	// order, so that the run can be shown afterwards to have used them once
	// --rm has removed its container (I-0274, S-0665); all empty on every
	// other run, and on runs before S-0665.
	Image     string                   `json:"image,omitempty"`
	Container string                   `json:"container,omitempty"`
	Mounts    []harness.ContainerMount `json:"mounts,omitempty"`
	// ImageJob is the id of the image job (ADR-0295 §5) that this run's
	// start asked for, or joined, to make the agent images it found missing
	// or stale, as image-jobs/<id>.job.json under flai serve's state folder
	// keeps it, and ImageHeld marks a story's run whose start waits for that
	// job: it has no process yet, and goes on once the job passes, or fails
	// with the job's cause (S-0751). Both are empty on a run whose start
	// needed no job; ImageJob stays on the run that went on, or failed.
	ImageJob  string `json:"image_job,omitempty"`
	ImageHeld bool   `json:"image_held,omitempty"`
}

// imageHeld is a run whose start waits for its image job (S-0751): not
// ended, not refused, and with no process yet.
func (r *AgentRun) imageHeld() bool {
	return r != nil && r.ImageHeld && r.Ended == "" && r.Error == "" && r.PID == 0
}

// same says whether r and o are one run: the same start of the same
// process.
func (r *AgentRun) same(o *AgentRun) bool {
	return r != nil && o != nil && r.Story == o.Story && r.Item == o.Item && r.Started == o.Started && r.PID == o.PID
}

// limitEnded marks r, a run that ended on lim, the harness's usage limit, as
// such (S-0534): UsageLimited, and ResetsAt when the harness named a reset
// time. Nothing when lim is nil.
func (r *AgentRun) limitEnded(lim *usage.Limit) {
	if lim != nil {
		r.UsageLimited, r.ResetsAt = true, limitResets(lim)
	}
}

// budgetEnd feeds the end of run, on lim, the usage limit its log ends on,
// or nil for none, into its provider's budget state (budgetEnded), when
// first, the end that recorded the run as ended first (ended), so that a run
// settled twice is fed once (S-0534). Every end is fed, not only one on the
// limit, so that the probe's clean end clears the state. A run the operator
// stopped is not: its end says nothing of the API, and a probe stopped is
// claimed over once its process is gone.
func (l *launcher) budgetEnd(run *AgentRun, lim *usage.Limit, first bool) {
	if first && run.Outcome != OutcomeStopped {
		l.budgetEnded(run, lim)
	}
}

// stopped makes an ended run read as the operator's stop, when they stopped
// it.
func (r *AgentRun) stopped() {
	r.Outcome, r.Why, r.Thread, r.Conversations, r.WaitsOn = OutcomeStopped, "stopped by the operator at "+r.Stopped, "", nil, nil
}

// agentChanged says whether the story's agent says something else than it
// did when run was started. A run that did not record it counts as unchanged.
func agentChanged(run *AgentRun, now *manifest.Agent) bool {
	return run.StoryAgent != nil && !run.StoryAgent.Same(now)
}

// orchestrates is the orchestrator's run: one for no story and no item, with
// no focus and no role of its own.
func (r *AgentRun) orchestrates() bool {
	return r.Story == "" && r.Item == "" && r.Focus == "" && r.Role == ""
}

// analyzes is the analyzer's run: one for no story and no item, with a focus
// (S-0223) and no role of its own.
func (r *AgentRun) analyzes() bool {
	return r.Story == "" && r.Item == "" && r.Focus != "" && r.Role == ""
}

// curates is the curator's run: one marked with its role (S-0431), whatever
// else it names, and no role the project defines.
func (r *AgentRun) curates() bool { return r.Role == conventions.RoleCurate && r.RoleName == "" }

// defined is a run of a role the project defines (S-0417), whatever else it
// names.
func (r *AgentRun) defined() bool { return r.RoleName != "" }

// live is a run that has not ended.
func (r *AgentRun) live() bool { return r != nil && r.Ended == "" && r.Error == "" && r.PID > 0 }

// running is a run that has not ended whose process is still the one
// started for it.
func (r *AgentRun) running() bool { return r.live() && Owns(r.PID, r.Start) }
