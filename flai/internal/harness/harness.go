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

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Command is the harness that runs the operator's own command (S-0079).
const Command = "command"

// Request is what an agent is started for.
type Request struct {
	Story   string          // the story's ID, checked by the caller
	Root    string          // the project's directory, where it runs
	Project string          // the project's key, for a session's name
	Agent   *manifest.Agent // the story's agent; nil when it has none
	Name    string          // FLAI_AGENT the session works under
	Flai    string          // this flai's executable, the agent's MCP server
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
// (S-0182). Only the claude-code adapter sends a prompt.
func Prompt(r Request) string {
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

Work %[2]s to review, and no other story. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md: prime your session with flai prime --story %[2]s (or the flai MCP tool prime), which prints the conventions that apply and what the story names whole, and briefs the design and ADRs its topics and links select, within a size budget. A brief is not the document: when one bears on the story, read it, or its section that does, with the flai MCP tool doc_get and its heading (flai doc show --heading on the host) before relying on it or changing what it describes, and find sections by their words with doc_search. Open the story with flai stream open %[2]s, write its tasks if it has none, and work them in the worktree that prints. Commit each task, keep the narrative's Current state and Next steps true, run flai stream sync %[2]s at every task transition, and call the flai MCP tool inbox there too.

%[6]s

%[4]s`, r.Name, r.Story, r.Root, rules(r), why, delegation(r))
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
// (ADR-0059, design/conventions/delegation.md).
func delegation(r Request) string {
	return fmt.Sprintf(`Keep your own context for decisions and edits, and hand noisy work to sub-agents with the Agent tool: code and document search across many files to the explorer, and test, lint, and flai check runs and long logs to the verifier. A sub-agent starts with nothing but your prompt: give it the worktree's path, %[1]s and the task's ID, the question, what you already know, and the shape of the answer you want, and ask for a summary with paths and lines, not raw output. Do it yourself when that is quicker: one file you know, one short command. A sub-agent cannot change work items or threads; a question it returns for the designer is yours to ask with thread_open.

While you work, run only the tests for what you changed, and leave the whole suite, the lint, and flai check to the verifier rather than running them yourself as well. Before you move %[1]s to review, commit everything, then have one fresh verifier run the whole suite, the lint, and flai check in the worktree, through the project's close-out script where it has one, and check the diff against the acceptance criteria and the conventions. Fix what it finds yourself, never through a sub-agent, commit, and have one more fresh verifier run the same and confirm the fixes. A verifier's passing run is the story's run before review: do not repeat it.`, r.Story)
}

func rules(r Request) string {
	return fmt.Sprintf(`When you need the designer to decide something, ask with the flai MCP tool thread_open on %[1]s, then call the flai MCP tool wait_for_events, again each time it returns, until the thread has an answer, and go on. If you end while the question is open, flai starts you again when it is answered. When every acceptance criterion is met, commit everything outstanding in the worktree, so that git status there is clean, then move %[1]s to review with flai move %[1]s review and end: the move is refused while anything is uncommitted. If you cannot go on, block the story with flai block %[1]s --reason and say why in its narrative, then end.`, r.Story)
}
