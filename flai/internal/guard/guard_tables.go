package guard

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// MCPReads are flai's MCP tools a sub-agent may call: they read and use no
// agent's identity.
// verify is not among them: it asks flai serve for a close-out job, which
// commits in the story's worktree (ADR-0229), so only the story's agent, whose
// own calls pass, and the orchestrator, which holds it in MCPOrchestrates and
// its role file, call it (S-0556).
var MCPReads = []string{"board", "doc_get", "doc_search", "item_get", "order_by_policy", "prime", "promote_candidates", "release_evaluate", "shared_paths", Subprojects, "test", "thread_get", "who_touches"}

// Subprojects is flai's MCP tool that lists the manifest's sub-projects, a
// read open to every agent; SubprojectEdit is the one that adds, changes, and
// removes one (S-0491). A story's agent and the operator's own session call
// SubprojectEdit and run flai subproject add, set, and remove with no rule:
// the sub-projects are the project's to shape, on the story's branch. Every
// sub-agent, a task sub-agent among them, is refused both, as it is flai
// components config, since they rewrite system-flow.yaml, which is the story's
// agent's to change; the planner and the analyzer are refused them by their
// own rules, which allow no write they do not name, and the orchestrator is
// allowed both only while orchestration.permissions.edit_subprojects is on
// (ADR-0195).
const (
	Subprojects    = "subprojects"
	SubprojectEdit = "subproject_edit"
)

// MCPStoryReads are flai's MCP tools that read for the calling agent's
// story, which a story's sub-agent may call besides MCPReads: message_get
// reads a conversation of the story's (S-0331, ADR-0120). The planner, the
// orchestrator, and the analyzer have no story, so the tools refuse them, and
// their allowlists leave them out. message_send, message_reply,
// message_escalate, and message_share, which write a conversation, the third
// a thread too (ADR-0121), and the last a share that lifts a hold
// (ADR-0134), are the story's agent's alone.
var MCPStoryReads = []string{"message_get"}

// ADRNew is flai's MCP tool that records an ADR, which a sub-agent may call
// without commit (S-0287, ADR-0127).
const ADRNew = "adr_new"

// adrWrites are the flai adr subcommands a sub-agent may run, without
// adrCommits, the flags with which they commit what they wrote (S-0287,
// ADR-0127). The planner, the orchestrator, and the analyzer run none.
var (
	adrWrites  = []string{"accept", "new", "topics"}
	adrCommits = []string{"--autocommit", "--commit"}
)

// Why a sub-agent writes an ADR but does not commit it, and what it does
// instead (S-0287, ADR-0127).
const (
	commitsADR    = "the commit is the story's agent's, since it writes history and widens the story's touches"
	leavesADRToIt = "and leave the files uncommitted for the story's agent to review and commit (ADR-0060, ADR-0127)"
)

// cliReads are the flai commands a sub-agent may run, each with the
// subcommands it may run; nil allows the command whatever follows it, and ""
// allows it with no subcommand. flai components list reads the component
// library, and flai components alone prints its help (S-0477, ADR-0186).
// flai digest list and show read the digest (S-0428). flai subproject list
// reads the manifest's sub-projects, and flai subproject alone prints its
// help (S-0491). flai formula list and show read the formulas, and flai
// formula alone prints its help (S-0421).
var cliReads = map[string][]string{
	"board":      {""},
	"check":      nil,
	"cod":        nil,
	"components": {"", "list"},
	"criteria":   {"list"},
	"digest":     {"list", "show"},
	"doc":        {"search", "show"},
	"forecast":   nil,
	"formula":    {"", "list", "show"},
	"help":       nil,
	"issue":      {"list"},
	"message":    {"list", "show"},
	"prime":      nil,
	"shared":     {"list", "check"},
	"show":       nil,
	"stats":      nil,
	"stream":     {"diff"},
	"subproject": {"", "list"},
	"test":       nil,
	"thread":     {"list", "show"},
	"touches":    {"suggest"},
	"version":    nil,
}

// flagReads are the flai commands a sub-agent may run only in the form that
// reads (S-0217, S-0219): with one of the flags that make them read, without
// any flag that makes them write, and in the form form allows when it is
// set. flai order --by computes an order and --apply writes it; flai order
// that names a story places it; flai release without --evaluate releases;
// flai plan without --candidates starts the planner, which with it takes
// no item; flai verify --record-issues writes issues, which without it it
// does not (S-0270); and flai verify --job asks flai serve for the close-out
// job, which commits in the story's worktree (S-0556, ADR-0229). flai role
// list, show, and definition read the role files, and flai role alone prints
// its help, but flai role run starts an agent and flai role definition
// --write writes the definitions into .flai-cache/agents (S-0568, I-0245).
// flai hostapi alone lists the host API's methods, and with a method on
// HostAPIReads it answers that read, --timing or not; every other method
// writes (S-0612, I-0307). flai serve reads in the forms serveReads lists,
// those ADR-0220 leaves open, and in none that writes the host
// configuration; flai serve alone runs the server, and its other forms start,
// stop, or commit (S-0680, I-0318).
// A command with no with flags reads in any form without its without flags.
var flagReads = map[string]struct {
	with    []string
	without []string
	form    func(words []string) bool
}{
	"hostapi": {form: readsHostAPI},
	"order":   {with: []string{"--by"}, without: []string{"--apply"}, form: ordersByPolicy},
	"plan":    {with: []string{"--candidates"}},
	"promote": {with: []string{"--candidates", "--drafts"}},
	"release": {with: []string{"--evaluate"}},
	"role":    {without: []string{"--write"}, form: readsRoles},
	"serve":   {form: readsServe},
	"verify":  {without: []string{"--job", "--record-issues"}},
}

// serveReads are the flai serve sub-commands that read in some form, each
// with the sub-commands of its own it reads with: nil in any form, "" with
// none. flai serve agent harness and provider, and flai serve forge, read
// bare or with a name alone, and write with a setting, which writesHost
// finds (ADR-0220).
var serveReads = map[string][]string{
	"actions":    nil,
	"agent":      {"", "show", "harness", "provider"},
	"checks":     {"show"},
	"components": {"show"},
	"forge":      nil,
	"import":     {"list"},
	"journal":    nil,
	"mcp":        {"", "show"},
	"notify":     {"show"},
	"status":     nil,
}

// readsServe says whether the words of flai serve read: a sub-command and
// form serveReads lists, which writes no host configuration.
func readsServe(words []string) bool {
	args := positionals(words, hostValues)
	if len(args) < 2 {
		return false
	}
	subs, ok := serveReads[args[1]]
	if !ok || writesHost("serve", args[1], words) {
		return false
	}
	var third string
	if len(args) > 2 {
		third = args[2]
	}
	return subs == nil || slices.Contains(subs, third)
}

// roleReads are the subcommands of flai role that read, "" for flai role
// alone.
var roleReads = []string{"", "definition", "list", "show"}

// readsRoles says whether the words of flai role name a subcommand that
// reads.
func readsRoles(words []string) bool {
	_, sub := subcommands(words, flaiValues)
	return slices.Contains(roleReads, sub)
}

// HostAPIReads are the methods of the host API's table (flai hostapi, ADR-0029)
// that write nothing: no work item, thread, document, file, or commit, no
// mark of what was seen, no setting, and no process started; at most the
// caches a flai read such as flai board keeps under .flai-cache (S-0612,
// I-0307). The guard cannot import flai/internal/hostapi, which imports it,
// so a test there keeps this list in step with the table.
var HostAPIReads = []string{
	"accept.preview", "activity.document", "activity.get", "adrs.list", "agent.status", "agent.stream",
	"board.get", "components.list", "design.get", "design.proposals", "doc.get", "doc.lint", "doc.show",
	"docs.tree", "forge.status", "inbox.designer", "issue.list", "item.get", "item.move.preview", "item.show",
	"items.count", "items.list", "kluster.status", "messages.get", "messages.list", "order.by", "project.info",
	"promote.candidates", "publish.preview", "release.evaluate", "search.query", "settings.get", "stats.get",
	"stream.diff", "threads.list", "verify.status",
}

// readsHostAPI says whether the words of flai hostapi read: they name no
// method, which lists the methods, or a method on HostAPIReads.
func readsHostAPI(words []string) bool {
	m := hostAPIMethod(words)
	return m == "" || slices.Contains(HostAPIReads, m)
}

// redirection begins a word that redirects a command's input or output, such
// as >f, 2>&1, or >, which the shell's parse leaves in a command's words.
var redirection = regexp.MustCompile(`^[0-9]*&?[<>]`)

// hostAPIMethod is the method the words of flai hostapi name, its first
// positional word after hostapi that is not a redirection or its target; ""
// when they name none.
func hostAPIMethod(words []string) string {
	args := positionals(words, flaiValues)
	for i := 1; i < len(args); i++ {
		w := args[i]
		if !redirection.MatchString(w) {
			return w
		}
		if strings.TrimLeft(w, "0123456789&<>|") == "" {
			i++ // the operator stands apart from its target, the next word
		}
	}
	return ""
}

// writesHostAPI says why a sub-agent does not run a host API method off
// HostAPIReads, naming it.
func writesHostAPI(words []string) string {
	return fmt.Sprintf("a sub-agent runs flai hostapi only with one of the host API's read methods, those on the guard's HostAPIReads, and %s is not one: the host API's write methods are the story's agent's (ADR-0060)", hostAPIMethod(words))
}

// orderValues are the flags of flai order, flai's own among them, that take
// a value; placements are those that place one story by hand.
var (
	orderValues = map[string]bool{"--after": true, "--before": true, "--by": true, "--config": true, "--keep-placed": true, "--placed-by": true}
	placements  = []string{"--after", "--before", "--bottom", "--placed-by", "--top"}
)

// gitReads are the git commands a sub-agent may run.
var gitReads = []string{"blame", "cat-file", "describe", "diff", "grep", "log", "ls-files", "ls-tree", "merge-base", "rev-list", "rev-parse", "shortlog", "show", "status"}

// gitFlags are git's own flags, before its subcommand, that take a value.
var gitFlags = map[string]bool{"-C": true, "-c": true, "--git-dir": true, "--work-tree": true, "--namespace": true}

// gitListForms are the git commands that write in some forms and only list in
// others (S-0623, I-0289): git tag and git branch list when they name no tag
// or branch, or when they are told to list, with -l, --list, or a filter
// that puts them in list mode, and carry no flag that writes; git worktree
// lists with list alone. values are the flags that take the next word as
// their value when not given with = (one whose value is optional and that
// does not make the command list, such as --column, takes it only with =);
// writes the flags that write, and in a cluster of short flags the letters
// of writes' short flags; lists the flags that make them list whatever words
// follow.
var gitListForms = map[string]struct {
	values, writes, lists []string
	sub                   string
}{
	"tag": {
		values: []string{"--contains", "--no-contains", "--points-at", "--merged", "--no-merged", "--sort", "--format", "-u", "-m", "-F", "--local-user", "--message", "--file", "--cleanup"},
		writes: []string{"-a", "-s", "-u", "-m", "-F", "-d", "-f", "-e", "--annotate", "--sign", "--no-sign", "--local-user", "--message", "--file", "--delete", "--force", "--edit", "--create-reflog", "--cleanup", "--trailer"},
		lists:  []string{"-l", "--list", "--contains", "--no-contains", "--points-at", "--merged", "--no-merged"},
	},
	"branch": {
		values: []string{"--contains", "--no-contains", "--points-at", "--merged", "--no-merged", "--sort", "--format", "-u", "--set-upstream-to"},
		writes: []string{"-d", "-D", "-m", "-M", "-c", "-C", "-f", "-u", "-t", "--delete", "--move", "--copy", "--force", "--set-upstream-to", "--unset-upstream", "--edit-description", "--track", "--no-track", "--create-reflog", "--recurse-submodules"},
		lists:  []string{"-l", "--list", "--contains", "--no-contains", "--points-at", "--merged", "--no-merged"},
	},
	"worktree": {sub: "list"},
}

// gitSubcommand is git's subcommand among the words after git, past git's
// own flags, and the words after it.
func gitSubcommand(words []string) (cmd string, args []string) {
	for i := 0; i < len(words); i++ {
		w := words[i]
		if strings.HasPrefix(w, "-") {
			if gitFlags[w] {
				i++
			}
			continue
		}
		return w, words[i+1:]
	}
	return "", nil
}

// gitListsOnly says whether git cmd, with the words args after it, only
// lists, as gitListForms says.
func gitListsOnly(cmd string, args []string) bool {
	f, ok := gitListForms[cmd]
	switch {
	case !ok:
		return false
	case f.sub != "":
		_, sub := subcommands(append([]string{cmd}, args...), nil)
		return sub == f.sub
	}
	var positional, listing bool
	for i := 0; i < len(args); i++ {
		w := args[i]
		if w == "--" {
			positional = positional || i+1 < len(args)
			break
		}
		if redirection.MatchString(w) {
			if strings.TrimLeft(w, "0123456789&<>|") == "" {
				i++ // the operator stands apart from its target, the next word
			}
			continue
		}
		if !strings.HasPrefix(w, "-") || w == "-" {
			positional = true
			continue
		}
		name, _, valued := strings.Cut(w, "=")
		if !strings.HasPrefix(w, "--") && len(w) > 2 {
			// a cluster of short flags, such as -vv or -dr
			for _, c := range w[1:] {
				if slices.Contains(f.writes, "-"+string(c)) {
					return false
				}
				listing = listing || slices.Contains(f.lists, "-"+string(c))
			}
			continue
		}
		if slices.Contains(f.writes, name) {
			return false
		}
		if slices.Contains(f.lists, name) {
			listing = true
		}
		if !valued && slices.Contains(f.values, name) {
			i++
		}
	}
	return listing || !positional
}
