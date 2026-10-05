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
	"regexp"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Command is the harness that runs the operator's own command (S-0079).
const Command = "command"

// Request is what an agent is started for: a story, or, for the planner, an
// epic or a story to plan (S-0208), or, for the orchestrator, the whole
// project (S-0218).
type Request struct {
	Story   string          // the story's ID, checked by the caller; empty for a strategic agent
	Root    string          // the project's directory, where it runs
	Project string          // the project's key, for a session's name
	Agent   *manifest.Agent // the story's agent, the planner's, or the orchestrator's; nil when it has none
	Name    string          // FLAI_AGENT the session works under
	Flai    string          // this flai's executable, the agent's MCP server
	// Role is conventions.RolePlan when the agent is the planner,
	// conventions.RoleOrchestrate when it is the orchestrator, and empty
	// when it works a story.
	Role string
	// Item is the epic or story the planner plans, checked by the caller;
	// empty when the agent works a story or orchestrates.
	Item string
	// Session names the harness's session, so that it can be resumed; a
	// harness that has no sessions ignores it.
	Session string
	// Answered is the thread the agent asked on that has been answered: the
	// agent ended while it waited, and is started again to go on (S-0104).
	Answered string
	// Restart says how the story's last agent ended, when the operator has
	// had a new one started for it (S-0116); empty when it entered ready.
	Restart string
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
}

// Adapter builds the command for one harness.
type Adapter interface {
	// DefaultHost is the program and arguments used when the operator has
	// set none.
	DefaultHost() Host
	// Start builds the command. It refuses what the story says that the
	// harness does not take.
	Start(req Request, host Host) (Start, error)
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
// answered and goes on in the session that already holds the rest; one
// started to commit a worktree is told to do only that (S-0140); one the
// operator started past a hold or a full limit is told what it went past
// (S-0182). Every agent working a story is told how its issues are recorded
// and that the operator chooses at acceptance which become stories (S-0198).
// The planner is asked to plan its item instead (planPrompt), and the
// orchestrator to keep the project's work moving (orchestratePrompt). Only
// the claude-code adapter sends a prompt.
func Prompt(r Request) string {
	switch r.Role {
	case conventions.RolePlan:
		return planPrompt(r)
	case conventions.RoleOrchestrate:
		return orchestratePrompt(r)
	}
	if r.Commit != "" && r.Answered == "" {
		return fmt.Sprintf(`You are %[1]s, started by flai serve on this host because the operator asked for the work left uncommitted in story %[2]s's worktree, %[3]s, to be committed: %[2]s is in review, and it cannot be accepted until that worktree is clean.

Do only this. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md, for how commits are made here. In the worktree, read git status and git diff. Commit the changes on story/%[2]s, in commits whose messages name %[2]s and say what changed and why, after the lint and tests the project runs for what they touch; fix what they find only when it is part of the same work. Discard a file only when it is plainly output that does not belong in the repository, and say which in the story's narrative with flai stream log %[2]s. Do not move %[2]s or its tasks, do not change anything the uncommitted work does not already change, and do not start other work. When git status in the worktree is clean, log what you committed with flai stream log %[2]s and end.

If a change cannot be committed without the designer deciding something, ask with the flai MCP tool thread_open on %[2]s, then call the flai MCP tool wait_for_events, again each time it returns, until the thread has an answer, and go on.`, r.Name, r.Story, r.Commit)
	}
	if r.Answered != "" {
		return fmt.Sprintf(`The designer has answered your question %[3]s on %[2]s. Read the answer with the flai MCP tool thread_get (or flai thread show %[3]s), then go on working %[2]s to review as before.

%[4]s`, r.Name, r.Story, r.Answered, rules(r))
	}
	why := "because it entered ready."
	switch {
	case r.Begun != nil:
		why = fmt.Sprintf("because the operator started it here, and this host has had no agent for it: it was begun %[1]s. Its branch and worktree may not be on this host, and what was not committed and pushed there is not here. Reconcile before you do anything else: run flai stream open %[2]s, which keeps the narrative and checks out story/%[2]s from this clone, else from the remote, else new from the main branch, and says which; read the narrative's Current state, Next steps, and log, and the story's tasks; compare them with what is committed on the branch; and go on from what is committed rather than starting over, doing again what the narrative says was done and is not there. Log what you found with flai stream log %[2]s.%[3]s", r.Begun.Said(), r.Story, answeredSince(r.Begun))
	case r.Restart != "":
		why = fmt.Sprintf("because the operator restarted it: its last agent %s. The story may already be in progress, with a narrative, a branch, and a worktree: if so, read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and go on from there rather than starting over.", r.Restart)
	case r.Started && len(r.Past) > 0:
		why = fmt.Sprintf("because the operator started it now, from the story's page or with flai serve agent start, past what kept flai serve from starting it: %s. That was the operator's decision: pull %[2]s as they asked, though flai warns that it is held or that the limit is full. Keep its touches to what it changes; do not narrow them only to clear the hold. Where its claim overlaps another open story's, say in the narrative which paths the two share.", strings.Join(r.Past, "; "), r.Story)
	case r.Started:
		why = "because the operator started it now, from the story's page or with flai serve agent start."
	}
	return fmt.Sprintf(`You are %[1]s, started by flai serve on this host to work story %[2]s in the project at %[3]s, %[5]s

Work %[2]s to review, and no other story. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md: prime your session with flai prime --story %[2]s (or the flai MCP tool prime), which prints the conventions that apply and what the story names whole, and briefs the design and ADRs its topics and links select, within a size budget. A brief is not the document: when one bears on the story, read it, or its section that does, with the flai MCP tool doc_get and its heading (flai doc show --heading on the host) before relying on it or changing what it describes, and find sections by their words with doc_search. Open the story with flai stream open %[2]s, write its tasks if it has none, and work them in the worktree that prints. When a task is done, commit its changes, with its docs and work-item updates, on story/%[2]s; then run flai stream sync %[2]s, and resolve each conflict it lists in the worktree, git add it, and git rebase --continue; then run the tests for what the task changed and commit any fix they need. flai stream sync does the branch's git work and refuses while anything is uncommitted: never start a rebase or merge by hand. Keep the narrative's Current state and Next steps true, and call the flai MCP tool inbox at every task transition.

%[6]s

%[4]s`, r.Name, r.Story, r.Root, rules(r), why, delegation(r))
}

// planPrompt is what the planner is asked to do (S-0208): plan its item, an
// epic or a story, as strategic-agents.md says, through flai alone, drafting
// a story's tasks or revisiting the open ones it has (S-0255); ask the
// operator on the item for an input it owns that is missing; and end with a
// summary, which flai serve logs as the run's activity (ADR-0079). An epic's
// planner writes its stories as drafts, summarises the plan in one thread on
// the epic, proposes there what it would change in the stories the epic
// already has, and names in its summary the stories it created and revisited
// (S-0209). A story's planner drafts its tasks, or revisits the open ones,
// and names in its summary the tasks it created and revisited (S-0255).
// It enriches a story from flai touches suggest, flai forecast, and flai cod,
// and records why under a ### Planning heading in the story's Notes (S-0210).
func planPrompt(r Request) string {
	kind := workitem.TypeOfID(r.Item)
	story := r.Item
	if kind == workitem.Epic {
		story = "S-nnnn"
	}
	enrich := fmt.Sprintf("its predicted touches, a forecast, and a cost of delay value worked out from the operator's inputs. "+
		"Run flai touches suggest %[1]s, adding the paths its goal, criteria, and linked design name when it declares no touches, and predict its touches from what that lists, its goal and criteria, the design documents it links, and the code layout, keeping every touch it already declares. "+
		"Run flai forecast %[1]s and flai cod %[1]s. "+
		"Review each figure, adjust it where you have a reason and state the reason, and write the touches, the forecast (duration, delivery, and basis), and the cost of delay value through flai: item_edit, or flai edit and flai touches. "+
		"In the story's Notes, under a ### Planning heading that is yours to rewrite, record where each touch came from (declared, co-change, design, or layout) and why each figure stands or was adjusted, and leave the rest of the Notes as it was", story)
	work := fmt.Sprintf("It is a story: enrich it with %[1]s. "+
		"If it has no tasks, draft the tasks that deliver its outcome, each with ## Work and ## Done when in its body, a nature, tags, touches (the paths it changes), and after (the tasks of the story it waits for), so that they form layers as work-management.md says, and create each in the backlog with the flai MCP tool item_new, type task and parent the story, or flai task new. "+
		"If it has tasks, revisit each one not done or cancelled against the story's outcome: re-enrich its touches and after with item_edit, create the tasks the outcome still lacks, and propose in the plan's thread any task you would split, merge, or drop. "+
		"Never cancel a task, or rewrite the title, Work, or Done when of a task you did not write, without asking on that thread. "+
		"Tasks carry no topics: when a task reaches a topic the story lacks, add the topic to the story. "+
		"Make each task you write pass flai check --strict and the markdown lint: flai refuses one that does not, so fix what the refusal names and write it again.\n\n"+
		"Open one thread on %[2]s that summarises the plan. In the plan's thread, name the tasks, their order and layers, and the assumptions you made.", enrich, r.Item)
	summary := "End with a one-line summary that names by ID the tasks you created and the tasks you revisited"
	if kind == workitem.Epic {
		work = fmt.Sprintf("It is an epic. If it has no stories, draft the stories that deliver its outcome, each with a goal, acceptance criteria as checkboxes, a nature, tags, topics, touches, and after, and create each with draft true in the backlog (item_new's draft, or flai story new --draft). If it has stories, revisit each one not done or cancelled against the epic's outcome, and enrich it again as you would a story: %[1]s. Size stories as work-management.md says, and make each one you write pass flai check --strict.\n\nOpen one thread on %[2]s that summarises the plan: the stories, their order (their after), and the assumptions you made. In that same thread, propose each story you would split, merge, add, or drop, and create drafts for the additions only: never cancel a finalized story or rewrite its words without asking.", enrich, r.Item)
		summary = "End with a one-line summary that names the stories you created and the stories you revisited"
	}
	return fmt.Sprintf(`You are %[1]s, the planner, started by flai serve on this host because the operator asked for %[2]s to be planned, in the project at %[3]s.

Plan %[2]s, and nothing else, as design/conventions/strategic-agents.md says under As the planner. Prime your session with flai prime --role plan --%[4]s %[2]s (or the flai MCP tool prime with role plan and %[4]s %[2]s), which prints the conventions you work by and what %[2]s names whole, and briefs the design its topics select. A brief is not the document: read the section that bears on the plan with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox. Read %[2]s with item_get, and what it links with item_get and doc_get. Hand wide search of the code, such as for a story's touches, to the explorer with the Agent tool.

%[5]s

Work in the main checkout and write only through flai: the flai MCP tools item_new and item_edit, or the flai CLI. Never edit a file yourself, code or anything else, never move an item past backlog, and never finalize a draft. Never overwrite the operator's inputs: a cost of delay's inputs, a story's estimate, and a finalized story's words; never cancel or rewrite a finalized story without asking.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on %[2]s, your recommended answer first, plan what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

%[6]s: flai serve logs the run's activity in wip/agents/planner.md with it.`, r.Name, r.Item, r.Root, kind, work, summary)
}

// orchestratePrompt is what the orchestrator is asked to do (S-0218): keep
// the project's work moving as strategic-agents.md says, only as far as the
// permissions the operator sets in orchestration.permissions allow, taking
// every order, candidate, and release figure from flai's commands rather
// than working it out; log each decision with activity_log, its reason and
// the policy figure behind it; and wait on wait_for_events between
// decisions, without ending, since flai serve runs it while the orchestrate
// host action is on. It never edits a file or works a story, and asks the
// operator where flai guard refuses it rather than working around the
// refusal.
func orchestratePrompt(r Request) string {
	return fmt.Sprintf(`You are %[1]s, the orchestrator, started by flai serve on this host because the operator turned on the orchestrate host action for the project at %[2]s.

Keep the project's work moving, and do nothing else, as design/conventions/strategic-agents.md says under As the orchestrator. Prime your session with flai prime --role orchestrate (or the flai MCP tool prime with role orchestrate), which prints the conventions you work by and briefs the design your role's topic selects. A brief is not the document: read the section that bears on a decision with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox, and read the board with the flai MCP tool board. Hand wide search to the explorer with the Agent tool.

Act only within the permissions the operator sets in system-flow.yaml under orchestration.permissions, each off by default, and by its policy, orchestration.policy. Ask the planner to plan a backlog epic, with the flai MCP tool plan, only while plan_backlog_epics is on; finalize a draft only while finalize_drafts is on; promote a story to ready only while promote_to_ready is on; order the ready column only while order_ready is on; answer a thread, or recommend an answer, only as answer_threads says; accept a story only while accept_reviews is on; publish a release only while publish is on. Do none of it while its permission is off, and if a permission is unclear, ask; do not act.

Take every figure from flai's commands and never do the arithmetic yourself: the ready column's order from flai order --by (the flai MCP tool order_by_policy), the stories that could go to ready from flai promote --candidates (promote_candidates), and whether a release is due from flai release --evaluate (release_evaluate).

Log each decision when you have made it with the flai MCP tool activity_log, kind orchestrator: what you did, on which items, why, and the policy figure behind it. Then hold the flai MCP tool wait_for_events, again each time it returns, and when something has changed, call inbox, read the board, and decide again. Repeat without ending: flai serve runs you for as long as the orchestrate host action is on.

Work in the main checkout and write only through flai: the flai MCP tools, or the flai CLI. Never edit code or documents, and never work a story yourself. flai guard refuses a call outside your permissions and names the permission it needs: never work around a refusal, by another tool, another command, or the shell. Ask the operator instead with the flai MCP tool thread_open on the item the decision concerns, your recommended answer first, and do what needs no answer meanwhile.`, r.Name, r.Root)
}

// roleEnv tells a strategic agent's session its role, which flai guard reads
// to hold it to its work: the planner's its item as well (S-0208), the
// orchestrator's nothing more, since it works the whole project (S-0218).
// Nothing for an agent working a story. A role flai does not start, a
// strategic agent with a story of its own, a planner with no epic or story
// to plan, and an orchestrator given an item are refused.
func roleEnv(r Request) ([]string, error) {
	switch {
	case r.Role == "":
		return nil, nil
	case r.Role != conventions.RolePlan && r.Role != conventions.RoleOrchestrate:
		return nil, fmt.Errorf("flai cannot start an agent in role %q; it starts the planner, role %s, the orchestrator, role %s, and an agent for a story", r.Role, conventions.RolePlan, conventions.RoleOrchestrate)
	case r.Role == conventions.RoleOrchestrate && r.Story != "":
		return nil, fmt.Errorf("the orchestrator was given story %s to work; it works no story of its own, so start it with none", r.Story)
	case r.Role == conventions.RoleOrchestrate && r.Item != "":
		return nil, fmt.Errorf("the orchestrator was given %s; it works the whole project, so start it with no item", r.Item)
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
// tasks are long, whose work it reviews, commits, and syncs itself (S-0176,
// S-0197, design/conventions/work-management.md), naming the task's ID in
// each task sub-agent's description so that flai measures the task by its
// calls (S-0230), and to sync once more before the close-out run (ADR-0069);
// and to have the verifier run the close-out once, without a pipe or a file,
// and read its last line, naming in its prompt any step already known to stop
// (S-0266); and to make independent edits and commands in one turn, and move a
// task it has just written to ready and in-progress in one command (S-0268).
func delegation(r Request) string {
	return fmt.Sprintf(`Keep your own context for decisions and edits, and hand noisy work to sub-agents with the Agent tool: code and document search across many files to the explorer, and test, lint, and flai check runs and long logs to the verifier. A sub-agent starts with nothing but your prompt: give it the worktree's path, %[1]s and the task's ID, the question, what you already know, and the shape of the answer you want, and ask for a summary with paths and lines, not raw output. Do it yourself when that is quicker: one file you know, one short command. A sub-agent cannot change work items or threads; a question it returns for the designer is yours to ask with thread_open. Make independent edits and commands in one turn, as several tool calls in one message: consecutive edits to one file, reads of files you already know, and commands that do not wait on each other.

Plan %[1]s's tasks as you write them: give each the paths it touches (--touches) and the tasks of %[1]s it must wait for (--after), so that tasks with no after between them and no path in common form layers that can run together, and record the layers and why each task waits in the narrative's Decisions. Move a task you have just written to ready and in-progress in one command, such as flai move T-nnnn ready && flai move T-nnnn in-progress: flai move refuses a task straight from backlog to in-progress. Work the plan layer by layer, handing each task to a task sub-agent with the Agent tool (a fork where you can), so that its reading stays out of your context; run a layer's tasks at once only when they are long, and then wait for all of them. Wait for a sub-agent run in the background through the harness's notice that it has finished, not by polling the flai MCP tool wait_for_events, which reports work items and threads, not sub-agents. Name the task's ID in each task sub-agent's description, so that flai measures the task by its calls. A task sub-agent edits only what its task touches, runs only its own tests, and never commits, syncs, or writes through flai. Review each one's work yourself, fix what falls short, commit it, sync and test as above, and move the task; only you commit, sync the stream, move items, and talk to the designer.

While you work, run only the tests for what you changed, and leave the whole suite, the lint, and flai check to the verifier rather than running them yourself as well. Before you move %[1]s to review, commit everything, run flai stream sync %[1]s again and resolve what it lists, then have one fresh verifier run the whole suite, the lint, and flai check in the worktree, through the project's close-out script where it has one (it refuses a branch that does not contain the main branch), and check the diff against the acceptance criteria and the conventions. Tell it to run the close-out once, in one command and without a pipe or a file, and to read its last line, which names the outcome and the step it stopped at; and name in its prompt any step you already know will stop, and why, so that it reports that stop and checks the steps after it rather than running the close-out again. Fix what it finds yourself, never through a sub-agent, commit, and have one more fresh verifier run the same and confirm the fixes. A verifier's passing run is the story's run before review: do not repeat it.`, r.Story)
}

func rules(r Request) string {
	return fmt.Sprintf(`Record friction, defects, and blockers you hit with flai issue new, or flai issue bump when the issue exists: each instance names %[1]s. Make no story for them yourself: %[1]s's review page lists them, checked, so the operator chooses at acceptance which become backlog stories. Make one with flai issue story or the flai MCP tool issue_story only when the operator asks, or when flai check warns issues.no-story.

When you need the designer to decide something, ask with the flai MCP tool thread_open on %[1]s, then call the flai MCP tool wait_for_events, again each time it returns, until the thread has an answer, and go on. If you end while the question is open, flai starts you again when it is answered. When every acceptance criterion is met, commit everything outstanding in the worktree, so that git status there is clean, run flai stream sync %[1]s again and resolve what it lists, then move %[1]s to review with flai move %[1]s review and end: the move is refused while anything is uncommitted. If you cannot go on, block the story with flai block %[1]s --reason and say why in its narrative, then end.`, r.Story)
}
