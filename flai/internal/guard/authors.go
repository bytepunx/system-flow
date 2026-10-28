package guard

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A served agent writes only as itself (S-0369, ADR-0143): a thread entry's
// author, a move's, a tick's, or a message's, is whatever the command that
// writes it says, --by, else FLAI_AGENT, else the config's author, so a
// story's agent, the planner, the analyzer, the curator, or a role a role file
// defines that named another could answer its own permission thread as the
// operator. The guard runs in the hook's environment, which the agent's
// command line cannot change, and so tells the agent's own name from the one
// a command names. The orchestrator keeps its own rule (S-0220); the
// operator's own session is not held.

// The variables that say whom flai writes as, and in which session.
const (
	agentEnv   = "FLAI_AGENT"
	sessionEnv = "FLAI_SESSION"
)

// identities are the variables a served agent leaves as flai serve set them.
var identities = []string{agentEnv, sessionEnv}

// authorFlags are the flags with which a flai command names who acted.
var authorFlags = []string{"--by", "--placed-by"}

// byNamesNoAuthor are the flai commands whose --by names no author: order's
// names a policy, and stats' a grouping.
var byNamesNoAuthor = []string{"order", "stats"}

// bySubNamesNoAuthor are the flai subcommands, command and subcommand, whose
// --by names no author: digest absorb's names the document that absorbed an
// entry (ADR-0199).
var bySubNamesNoAuthor = [][2]string{{"digest", "absorb"}}

// authorValues are flaiValues and the author flags: flai's flags that take a
// value before its command, as in flai --by alex thread reply.
var authorValues = map[string]bool{"--by": true, "--config": true, "--placed-by": true}

// declarers are the shell's builtins whose words assign variables, as in
// export FLAI_AGENT=alex.
var declarers = []string{"declare", "export", "local", "readonly", "typeset"}

// leaveIdentities ends a refusal of a change to the environment.
const leaveIdentities = "leave FLAI_AGENT and FLAI_SESSION as flai serve set them"

// authors says why a served agent's own shell call is refused, or "" when it
// is not (S-0369, ADR-0143): a story's agent's, the planner's, the
// analyzer's, the curator's, or a defined role's call that names an author
// other than its own with --by or --placed-by, or sets, empties, or unsets
// FLAI_AGENT or FLAI_SESSION to anything but its own, in any simple command
// of the line.
func (g Guard) authors(c Call) string {
	if c.Subagent || c.Kind != Shell || !g.heldToSelf(c) {
		return ""
	}
	for _, words := range commands(c.Command) {
		if why := g.impersonates(c, words); why != "" {
			return fmt.Sprintf("%s cannot run %q: %s.", g.who(c), strings.Join(words, " "), why)
		}
	}
	return ""
}

// heldToSelf says whether c's session is one ADR-0143 holds to its own name:
// a story's agent, or another session flai serve started with no role, the
// planner, the analyzer, the curator, or a role a role file defines (S-0416,
// ADR-0207 §5).
func (g Guard) heldToSelf(c Call) bool {
	return c.Role == RolePlan || c.Role == RoleAnalyze || c.Role == RoleCurate || DefinedRole(c.Role) || c.Role == "" && g.served(c)
}

// selves are the names a served agent writes as: Agent, the hook's
// FLAI_AGENT, and the name flai serve gives its run, agent-<Story>,
// planner-<Item>, analyzer, or curator, or, for a defined role, its role
// file's name, <name>-<Item> on an item (ADR-0207 §2); none when no role file
// names it.
func (g Guard) selves(c Call) []string {
	var run string
	switch {
	case c.Role == "" && c.Story != "":
		run = "agent-" + workitem.CanonicalID(c.Story)
	case c.Role == RolePlan && g.Item != "":
		run = workitem.ActivityPlanner + "-" + workitem.CanonicalID(g.Item)
	case c.Role == RoleAnalyze:
		run = workitem.ActivityAnalyzer
	case c.Role == RoleCurate:
		run = workitem.ActivityCurator
	case DefinedRole(c.Role) && g.RoleFile != nil && g.RoleFile.Name != "":
		run = g.RoleFile.Name
		if g.Item != "" {
			run += "-" + workitem.CanonicalID(g.Item)
		}
	}
	var names []string
	for _, n := range []string{g.Agent, run} {
		if n != "" && !slices.Contains(names, n) {
			names = append(names, n)
		}
	}
	return names
}

// own says whether name is one a served agent writes as; "" never is.
func (g Guard) own(c Call, name string) bool {
	return name != "" && slices.Contains(g.selves(c), name)
}

// writesAs says, in a refusal, whom and where a served agent writes as:
// itself, by its names, for FLAI_AGENT and an author, and its own session for
// FLAI_SESSION.
func (g Guard) writesAs(c Call, variable string) string {
	if variable == sessionEnv {
		if g.Session == "" {
			return "a served agent writes only in its own session"
		}
		return "a served agent writes only in its own session, " + g.Session
	}
	if names := g.selves(c); len(names) > 0 {
		return "a served agent writes only as itself, " + strings.Join(names, " or ")
	}
	return "a served agent writes only as itself"
}

// impersonates says why one simple command writes as another, or "" when it
// does not: a flai command whose --by or --placed-by names another; an
// assignment of FLAI_AGENT or FLAI_SESSION not the agent's own before the
// program, after a declarer, or among env's own words; unset of either,
// env's -u of either, or env -i; and the same in the script a shell given -c,
// or eval, runs. A declarer and unset, builtins no wrapper runs, count only
// as the program; env, flai, a shell, and eval count wherever they stand, as
// refuse finds them.
func (g Guard) impersonates(c Call, words []string) string {
	program := true
	for i, w := range words {
		if assignment.MatchString(w) {
			if program {
				if why := g.assigns(c, w); why != "" {
					return why
				}
			}
			continue
		}
		rest := words[i+1:]
		switch base := path.Base(w); {
		case program && slices.Contains(declarers, w):
			for _, a := range rest {
				if assignment.MatchString(a) {
					if why := g.assigns(c, a); why != "" {
						return why
					}
				}
			}
		case program && w == "unset":
			for _, a := range rest {
				if slices.Contains(identities, a) {
					return g.unsets(c, a)
				}
			}
		case base == "env":
			if why := g.env(c, rest); why != "" {
				return why
			}
		case slices.Contains(shells, base):
			for j, f := range rest {
				if shellFlags.MatchString(f) && j+1 < len(rest) {
					if why := g.impersonatesIn(c, rest[j+1]); why != "" {
						return why
					}
					break
				}
			}
		case base == "eval":
			if why := g.impersonatesIn(c, strings.Join(rest, " ")); why != "" {
				return why
			}
		case base == "flai", base == "flai.sh":
			if why := g.names(c, rest); why != "" {
				return why
			}
		}
		program = false
	}
	return ""
}

// impersonatesIn checks each simple command of a script a shell is given.
func (g Guard) impersonatesIn(c Call, text string) string {
	for _, words := range commands(text) {
		if why := g.impersonates(c, words); why != "" {
			return why
		}
	}
	return ""
}

// names says why the words after flai name an author other than the agent's
// own, or "" when they do not: --by, but on order, stats, and digest
// absorb, and --placed-by, each as its last occurrence gives it. An empty value falls
// back to FLAI_AGENT and passes; so does a word flai that runs none of its
// commands, or asks for help.
func (g Guard) names(c Call, rest []string) string {
	cmd, sub := subcommands(rest, authorValues)
	if !slices.Contains(g.Commands, cmd) || slices.Contains(rest, "--help") || slices.Contains(rest, "-h") {
		return ""
	}
	for _, f := range authorFlags {
		if f == "--by" && (slices.Contains(byNamesNoAuthor, cmd) || slices.Contains(bySubNamesNoAuthor, [2]string{cmd, sub})) {
			continue
		}
		if v := value(rest, f); v != "" && !g.own(c, v) {
			return fmt.Sprintf("it names %s as its author with %s, and %s (ADR-0143); leave out %s, or give its own name", v, f, g.writesAs(c, agentEnv), f)
		}
	}
	return ""
}

// assigns says why an assignment, NAME=value, writes as another, or "" when
// it does not: FLAI_AGENT to a name not the agent's own, or FLAI_SESSION to a
// session not its own, the empty value included. The guard reads a word up
// to a $, so a value held in a variable or a command substitution reads as
// empty, and is not its own.
func (g Guard) assigns(c Call, w string) string {
	name, v, _ := strings.Cut(w, "=")
	if !slices.Contains(identities, name) || name == agentEnv && g.own(c, v) || name == sessionEnv && v != "" && v == g.Session {
		return ""
	}
	set := "to " + v
	if v == "" {
		set = "to nothing, or to a value the guard cannot read, such as a variable's"
	}
	return fmt.Sprintf("it sets %s %s, and %s (ADR-0143); %s", name, set, g.writesAs(c, name), leaveIdentities)
}

// unsets says why removing a variable of identities from the environment is
// refused.
func (g Guard) unsets(c Call, variable string) string {
	so := ""
	if variable == agentEnv {
		so = ", so that flai would write as the config's author"
	}
	return fmt.Sprintf("it unsets %s%s, and %s (ADR-0143); %s", variable, so, g.writesAs(c, variable), leaveIdentities)
}

// env says why env's own words, its flags and assignments before the command
// it runs, write as another, or "" when they do not: -i, -, or
// --ignore-environment, which clear FLAI_AGENT; -u or --unset of FLAI_AGENT or
// FLAI_SESSION; an assignment of either not the agent's own; and the same in
// -S's string, which env splits into its words.
func (g Guard) env(c Call, args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-" || a == "--ignore-environment":
			return g.clears(c)
		case a == "--unset" || a == "--chdir" || a == "--split-string":
			if i+1 == len(args) {
				return ""
			}
			i++
			if why := g.envValue(c, a, args[i]); why != "" {
				return why
			}
		case strings.HasPrefix(a, "--"):
			flag, v, _ := strings.Cut(a, "=")
			if why := g.envValue(c, flag, v); why != "" {
				return why
			}
		case strings.HasPrefix(a, "-"):
			flag, v, ok := envShort(a)
			if flag == "-i" {
				return g.clears(c)
			}
			if ok && v == "" && i+1 < len(args) {
				i++
				v = args[i]
			}
			if why := g.envValue(c, flag, v); why != "" {
				return why
			}
		case assignment.MatchString(a):
			if why := g.assigns(c, a); why != "" {
				return why
			}
		default:
			return ""
		}
	}
	return ""
}

// envShort reads a cluster of env's short flags, as -0iu NAME or -uNAME: -i
// when it clears the environment, else the last flag that takes a value, -u,
// -C, or -S, with the rest of the cluster as its value, and whether it takes
// one; "" when none does.
func envShort(cluster string) (flag, v string, takes bool) {
	for j, r := range cluster[1:] {
		switch r {
		case 'i':
			return "-i", "", false
		case 'u', 'C', 'S':
			return "-" + string(r), cluster[j+2:], true
		}
	}
	return "", "", false
}

// envValue says why one of env's flags with its value writes as another:
// -u or --unset of a variable of identities, or -S's or --split-string's
// string, read as env's words.
func (g Guard) envValue(c Call, flag, v string) string {
	switch flag {
	case "-u", "--unset":
		if slices.Contains(identities, v) {
			return g.unsets(c, v)
		}
	case "-S", "--split-string":
		return g.env(c, strings.Fields(v))
	}
	return ""
}

// clears says why env -i is refused: it runs its command with no
// FLAI_AGENT.
func (g Guard) clears(c Call) string {
	return fmt.Sprintf("it clears the environment, so that flai would write as the config's author, and %s (ADR-0143); %s", g.writesAs(c, agentEnv), leaveIdentities)
}
