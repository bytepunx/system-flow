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
	"slices"
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
// project (S-0218), or, for the analyzer, the whole project with a focus
// or none (S-0223).
type Request struct {
	Story   string          // the story's ID, checked by the caller; empty for a strategic agent
	Root    string          // the project's directory, where it runs
	Project string          // the project's key, for a session's name
	Agent   *manifest.Agent // the story's agent, or the strategic agent's; nil when it has none
	Name    string          // FLAI_AGENT the session works under
	Flai    string          // this flai's executable, the agent's MCP server
	// Role is conventions.RolePlan when the agent is the planner,
	// conventions.RoleOrchestrate when it is the orchestrator,
	// conventions.RoleAnalyze when it is the analyzer, and empty when it
	// works a story.
	Role string
	// Item is the epic or story the planner plans, checked by the caller;
	// empty when the agent works a story, orchestrates, or analyzes.
	Item string
	// Focus is what the analyzer looks for, one of Focuses; empty for all
	// of them, and for any other agent.
	Focus string
	// Session names the harness's session, so that it can be resumed; a
	// harness that has no sessions ignores it.
	Session string
	// Answered is the thread the agent asked on that has been answered, or
	// the conversation (MS-nnnn) another story's agent wrote on to its story:
	// the agent ended while it waited, and is started again to go on (S-0104,
	// S-0335).
	Answered string
	// Restart says how the story's last agent ended, when the operator has
	// had a new one started for it (S-0116); empty when it entered ready.
	Restart string
	// AutoRestart says which automatic restart this is, such as "1 of 2",
	// when flai serve, not the operator, started the agent again after the
	// last one ended with its story in progress (S-0294, ADR-0108); empty
	// otherwise.
	AutoRestart string
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
// An agent for a ready story begins it with story_start (S-0274); one for a
// story already in progress, which story_start refuses, primes and takes it
// up with flai stream open. One started on a share that cleared a hold on
// its story is told the share's conversation, paths, and split (ADR-0134,
// shared). The planner is asked to plan its item instead
// (planPrompt), the orchestrator to keep the project's work moving
// (orchestratePrompt), and the analyzer to write a report (analyzePrompt).
// Only the claude-code adapter sends a prompt.
func Prompt(r Request) string {
	switch r.Role {
	case conventions.RolePlan:
		return planPrompt(r)
	case conventions.RoleOrchestrate:
		return orchestratePrompt(r)
	case conventions.RoleAnalyze:
		return analyzePrompt(r)
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
	why := "because it entered ready."
	open := startStory(r.Story, "Begin")
	switch {
	case r.Begun != nil, r.Restart != "" && r.AutoRestart != "":
		open = takeUp(r.Story, "Prime")
	case r.Restart != "":
		open = fmt.Sprintf("If %[1]s is still ready, ", r.Story) + startStory(r.Story, "begin") + " If it is in progress already, " + takeUp(r.Story, "prime")
	case r.Started && len(r.Past) > 0:
		open = startStory(r.Story, "Begin") + fmt.Sprintf(" If story_start refuses %[1]s because the board holds it, start it as the operator asked with flai move %[1]s in-progress, which only warns of the hold, then ", r.Story) + takeUp(r.Story, "prime")
	}
	switch {
	case r.Begun != nil:
		why = fmt.Sprintf("because the operator started it here, and this host has had no agent for it: it was begun %[1]s. Its branch and worktree may not be on this host, and what was not committed and pushed there is not here. Reconcile before you do anything else: run flai stream open %[2]s, which keeps the narrative and checks out story/%[2]s from this clone, else from the remote, else new from the main branch, and says which; read the narrative's Current state, Next steps, and log, and the story's tasks; compare them with what is committed on the branch; and go on from what is committed rather than starting over, doing again what the narrative says was done and is not there. Log what you found with flai stream log %[2]s.%[3]s", r.Begun.Said(), r.Story, answeredSince(r.Begun))
	case r.Restart != "" && r.AutoRestart != "":
		why = fmt.Sprintf("because its last agent %s, and flai serve started it again on its own, automatic restart %s; at the limit it asks the operator instead. The story is in progress, with a narrative, a branch, and a worktree: read the narrative's Current state and Next steps, reconcile them with git status in the worktree, and go on from there rather than starting over.", r.Restart, r.AutoRestart)
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
// no task can name yet may be added, and says why there (S-0295, ADR-0096).
func planPrompt(r Request) string {
	kind := workitem.TypeOfID(r.Item)
	story := r.Item
	if kind == workitem.Epic {
		story = "S-nnnn"
	}
	enrich := fmt.Sprintf("its predicted touches, a forecast, and a cost of delay value worked out from the operator's inputs. "+
		"Run flai touches suggest %[1]s, adding the paths its goal, criteria, and linked design name when it declares no touches, and predict its touches from what that lists, its goal and criteria, the design documents it links, and the code layout, keeping every touch it already declares. "+
		"Name files, not folders, in the touches you write for the story and for each of its tasks: keep a folder touch only where the story may add files there that no task can name yet, since a folder touch claims every file below it and, while the story is in progress, holds every ready story that touches one; record each folder touch you kept, and why, under the story's ### Planning heading. "+
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
	return fmt.Sprintf(`You are %[1]s, the planner, started by flai serve on this host because the operator asked for %[2]s to be planned, in the project at %[3]s.

Plan %[2]s, and nothing else, as design/conventions/strategic-agents.md says under As the planner. Prime your session with flai prime --role plan --%[4]s %[2]s (or the flai MCP tool prime with role plan and %[4]s %[2]s), which prints the conventions you work by and what %[2]s names whole, and briefs the design its topics select. A brief is not the document: read the section that bears on the plan with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox. Read %[2]s with item_get, and what it links with item_get and doc_get. Hand wide search of the code, such as for a story's touches, to the explorer with the Agent tool.

%[5]s

Work in the main checkout and write only through flai: the flai MCP tools item_new and item_edit, or the flai CLI. Never edit a file yourself, code or anything else, never move an item past backlog, and never finalize a draft. Never overwrite the operator's inputs: a cost of delay's inputs, a story's estimate, and a finalized story's words; never cancel or rewrite a finalized story without asking.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on %[2]s, your recommended answer first, plan what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

%[6]s: flai serve logs the run's activity in wip/agents/planner.md with it.`, r.Name, r.Item, r.Root, kind, work, summary)
}

// orchestratePrompt is what the orchestrator is asked to do (S-0218): keep
// the project's work moving as strategic-agents.md says, only as far as the
// permissions the operator sets in orchestration.permissions allow, which
// with its policy and release policy it reads again before each decision,
// since the operator may change them while it runs and the guard holds each
// call to them as they are then (S-0229), taking
// every order, candidate, and release figure from flai's commands rather
// than working it out; log each action with activity_log, its reason and
// the policy figure that justified it; and wait on wait_for_events between
// decisions, without ending, since flai serve runs it while the orchestrate
// host action is on. It is told what to do with each of plan_backlog_epics,
// finalize_drafts, promote_to_ready, and order_ready, and by which command
// (S-0219), and what to do with the threads awaiting the operator under
// each value of answer_threads, which it reads afresh each time since the
// operator may change it while it runs (S-0220), and how it accepts a story
// in review under accept_reviews: flai verify runs, or its last result is
// read, at the branch's head, its verifier checks the diff against each
// criterion (S-0270), flai accept --dry-run names the blockers, and it accepts with the verified
// commit and evidence for every criterion, or leaves the story in review and
// says on a thread what is missing (S-0221, ADR-0093). While publish is on,
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
// each, and never confirms a recommendation (S-0328, ADR-0119).
func orchestratePrompt(r Request) string {
	return fmt.Sprintf(`You are %[1]s, the orchestrator, started by flai serve on this host because the operator turned on the orchestrate host action for the project at %[2]s.

Keep the project's work moving, and do nothing else, as design/conventions/strategic-agents.md says under As the orchestrator. Prime your session with flai prime --role orchestrate (or the flai MCP tool prime with role orchestrate), which prints the conventions you work by and briefs the design your role's topic selects. A brief is not the document: read the section that bears on a decision with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox, and read the board with the flai MCP tool board. Hand wide search to the explorer with the Agent tool.

Act only within the permissions the operator sets in system-flow.yaml under orchestration.permissions, each off by default, and by its policy, orchestration.policy. Ask the planner to plan a backlog epic, with the flai MCP tool plan, only while plan_backlog_epics is on; ask it to plan a backlog story, with plan, and settle the threads a story's planner opens, giving its story cost of delay inputs with item_edit, only while plan_backlog_stories is on; finalize a draft only while finalize_drafts is on; promote a story to ready only while promote_to_ready is on; order the ready column only while order_ready is on; answer a thread, or recommend an answer, only as answer_threads says; accept a story only while accept_reviews is on; publish a release only while publish is on. Do none of it while its permission is off, and if a permission is unclear, ask; do not act. The operator may change your permissions, your policy, and the release policy, orchestration.release, while you run: read them again in system-flow.yaml in the main checkout before each decision rather than keep what you read at the start, and when flai guard's verdict on a call differs from what you read, the guard's verdict holds.

Take every figure from flai's commands and never do the arithmetic yourself: the epics to plan from flai plan --candidates, and after them the stories, the drafts complete enough to finalize from flai promote --drafts, the ready column's order from flai order --by (the flai MCP tool order_by_policy), the stories that could go to ready from flai promote --candidates (promote_candidates), and whether a release is due from flai release --evaluate (release_evaluate).

With each permission, do this. While plan_backlog_epics is on, run flai plan --candidates and start the planner with the flai MCP tool plan for each epic it lists, one at a time, and for no other. While plan_backlog_stories is on, after the epics, start the planner with the flai MCP tool plan for each story flai plan --candidates lists, those whose type is story, one at a time, and for no other. While finalize_drafts is on, run flai promote --drafts: finalize a draft it lists as complete, and whose criteria, touches, forecast, and value you judge consistent, with the flai MCP tool item_edit giving only its id and draft false; for any other draft, open one thread on the story saying what it lacks or what is inconsistent, once, and leave the draft as it is. While promote_to_ready is on, run flai promote --candidates and move its candidates to ready in its order with the flai MCP tool item_move while the ready column's limit has room: never a draft, and never a held story. While order_ready is on, run flai order --by <policy> --apply, with the policy orchestration.policy names, after each change to the ready column, yours or another's; it keeps in place a story placed by hand within the last day. Never place a story by hand yourself.

Reply on threads only as answer_threads says, and read it in system-flow.yaml in the main checkout each time before you act on threads: the operator may change it while you run. Take from inbox the threads awaiting the operator, those whose status is open and whose last entry is by a story's agent, leaving out the threads you opened and those whose pending_recommendation is not null. While answer_threads is off, leave them alone. While it is recommend, reply to each with the flai MCP tool thread_reply, recommendation true, and a source: the ADR, design section, or convention your answer rests on, as <path> or <path>#<heading>, read with doc_get first. While it is autonomous, answer with thread_reply and a source when a source settles the question; post a recommendation instead, citing what it draws on, which escalates it to the operator, when no source settles it, or when it asks for the operator's judgement: a decision not yet recorded, a change of scope, or money, such as a cost of delay input, an estimate, or spend, save the cost of delay inputs a story's planner asks for while plan_backlog_stories is on. thread_reply logs the reply and its source in your decision log. Never resolve a thread you did not open, but a story's planner's while plan_backlog_stories is on, never answer a thread you opened, and never confirm a recommendation: the operator does.

While plan_backlog_stories is on, settle the threads a story's planner opened, those whose opener is planner-S-nnnn, whatever answer_threads says. Approve a plan whose tasks, touches, and figures fit the story: reply so with thread_reply, and resolve the thread with the flai MCP tool thread_resolve. Answer its questions with thread_reply, with a source when one settles them. For a cost of delay input it asks for, take the figure it recommends unless the thread or the story gives a reason for one of the alternatives it lists. Set it with the flai MCP tool item_edit, cost_of_delay on the story, giving the inputs only, revenue_per_week, penalty_per_week, or time_lost_per_cycle, never a value and nothing else: the planner works the value out on its next run, and flai refuses inputs once the story or its epic has any. Then reply naming what you set, and resolve the thread. Log each approval, answer, and cost of delay you set with activity_log. When a question asks for a change of scope or a decision not yet recorded, post a recommendation instead, which escalates it to the operator, and leave the thread open. Never confirm a recommendation: the operator does. While plan_backlog_stories is off, a story's planner's threads are the operator's, as answer_threads says.

While accept_reviews is on, take each story in review in turn. Verify it at the head of its branch, story/<S-nnnn>: read its last result with flai verify <S-nnnn> --last --json, and when it did not pass or its commit is not the branch's head (git rev-parse story/<S-nnnn>), run flai verify <S-nnnn> --json, or the flai MCP tool verify, which runs the tests, the lint, and flai check --strict in the story's worktree and answers the commit it verified. A run that did not pass is a blocker. Then hand its worktree, %[2]s/.flai-cache/worktrees/<S-nnnn>, and that result to the verifier with the Agent tool: it reads the result with flai verify <S-nnnn> --last rather than running the suite, checks the diff against each acceptance criterion, naming for each the changed files that meet it, and names the commit it verified. Run flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run with that commit, and read the blockers. With no blocker and every criterion matched to changed files, write the evidence, a Verdict: line from the verifier's report and one list item per criterion, - <n>: <files>, and run flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -, passing the evidence on standard input with a heredoc in the shell, since you cannot write a file; then log the acceptance with activity_log, naming the story and the commit. Otherwise leave the story in review, open a thread on it with the flai MCP tool thread_open saying what is missing, each blocker and each criterion you could not check against the diff, and log that decision. A commit added after the verifier's run makes flai refuse: verify again. Never move a story to done with item_move, and never accept a story while accept_reviews is off.

While publish is on, after each acceptance, yours or one wait_for_events reports as a story moved to done, call the flai MCP tool release_evaluate. Under the threshold or theme policy, when the evaluation is met, call the flai MCP tool release_publish with the figure that was met as its reason. Under judgement, call release_publish only when you judge the unreleased work coherent and complete, with that reasoning as its reason. Under any policy, never publish a batch the evaluation holds back under whole_epics, the stories it names in held_by_epic. Log each release with activity_log: the policy, its figures, the versions and tags, and the items bundled, all as release_publish returned them. Log a decision not to publish too, with the evaluation's figures. When release_publish refuses, log the refusal with activity_log, and open a thread to the operator with the flai MCP tool thread_open on the most recently accepted story of the batch, with the refusal's words and what fixes it; do not try again until that thread is answered or the next acceptance. While publish is off, neither evaluate nor publish. Publish only through release_publish: flai guard refuses you flai release other than --evaluate, flai push, git push, and git tag, whatever your permissions.

Log each action when you have taken it with the flai MCP tool activity_log, kind orchestrator: what you did, to which item, why, and the policy figure that justified it, as flai gave it: the story's value, its value over its duration, its forecast, or the candidate's rank. Then hold the flai MCP tool wait_for_events, again each time it returns, and when something has changed, call inbox, read the board, and decide again. Repeat without ending: flai serve runs you for as long as the orchestrate host action is on.

Work in the main checkout and write only through flai: the flai MCP tools, or the flai CLI. Never edit code or documents, and never work a story yourself. flai guard refuses a call outside your permissions and names the permission it needs: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt: log it with activity_log, with the refusal, and do not try it again until something it depends on changes. When a decision needs the operator, ask with the flai MCP tool thread_open on the item it concerns, your recommended answer first, and do what needs no answer meanwhile.`, r.Name, r.Root)
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

// roleEnv tells a strategic agent's session its role, which flai guard reads
// to hold it to its work: the planner's its item as well (S-0208), the
// orchestrator's and the analyzer's nothing more, since each works the whole
// project (S-0218, S-0223). Nothing for an agent working a story. A role
// flai does not start, a strategic agent with a story of its own, a planner
// with no epic or story to plan, an orchestrator or an analyzer given an
// item, an analyzer with a focus it does not take, and a focus for any
// other agent are refused.
func roleEnv(r Request) ([]string, error) {
	switch {
	case r.Role == "" && r.Focus == "":
		return nil, nil
	case r.Role != conventions.RoleAnalyze && r.Focus != "":
		return nil, fmt.Errorf("only the analyzer takes a focus; start this agent without %q", r.Focus)
	case r.Role != conventions.RolePlan && r.Role != conventions.RoleOrchestrate && r.Role != conventions.RoleAnalyze:
		return nil, fmt.Errorf("flai cannot start an agent in role %q; it starts the planner, role %s, the orchestrator, role %s, the analyzer, role %s, and an agent for a story", r.Role, conventions.RolePlan, conventions.RoleOrchestrate, conventions.RoleAnalyze)
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
