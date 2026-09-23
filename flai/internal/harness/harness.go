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
// else, the way the project's conventions say, asking the designer through
// flai when it needs them.
func Prompt(r Request) string {
	return fmt.Sprintf(`You are %[1]s, started by flai serve on this host to work story %[2]s in the project at %[3]s, because it entered ready.

Work %[2]s to review, and no other story. Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md: prime your session with flai prime --cat, open the story with flai stream open %[2]s, write its tasks if it has none, and work them in the worktree that prints. Commit each task, keep the narrative's Current state and Next steps true, run flai stream sync %[2]s at every task transition, and call the flai MCP tool inbox there too.

When you need the designer to decide something, ask with the flai MCP tool thread_open on %[2]s, then hold wait_for_events until the thread is answered, and go on. Do not end while a question you asked is open. When every acceptance criterion is met, move %[2]s to review with flai move %[2]s review and end. If you cannot go on, block the story with flai block %[2]s --reason and say why in its narrative, then end.`, r.Name, r.Story, r.Root)
}
