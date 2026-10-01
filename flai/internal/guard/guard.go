// Package guard refuses what a sub-agent may not do (ADR-0060): Claude Code
// runs flai guard before each tool call in a project made from the template,
// with the call on standard input, and says which sub-agent makes it. A
// sub-agent reads; it does not change a work item, a thread, or the
// repository's history, because the story's agent alone acts for the story
// (ADR-0059). The story's agent's own calls carry no agent ID and are never
// refused. The guard is not a shell: it finds the programs a command line
// runs well enough to stop a sub-agent that follows its instructions
// carelessly, not one that sets out to hide a command (a backslash inside a
// name, a command held in a variable).
package guard

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"
)

// Event is the part of a Claude Code PreToolUse hook's input the guard
// reads.
type Event struct {
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
	// AgentID is set only when a sub-agent makes the call; AgentType names
	// the sub-agent's definition. A session started with --agent may carry
	// an agent type of its own, so the ID is what marks a sub-agent.
	AgentID   string `json:"agent_id"`
	AgentType string `json:"agent_type"`
}

// MCPPrefix starts the name of each of flai's MCP tools as Claude Code
// knows them, from a server named flai.
const MCPPrefix = "mcp__flai__"

// MCPReads are flai's MCP tools a sub-agent may call: they read and use no
// agent's identity.
var MCPReads = []string{"board", "doc_get", "doc_search", "item_get", "prime", "thread_get", "who_touches"}

// cliReads are the flai commands a sub-agent may run, each with the
// subcommands it may run; nil allows the command whatever follows it, and ""
// allows it with no subcommand.
var cliReads = map[string][]string{
	"board":   {""},
	"check":   nil,
	"doc":     {"search", "show"},
	"help":    nil,
	"issue":   {"list"},
	"prime":   nil,
	"show":    nil,
	"stats":   nil,
	"stream":  {"diff"},
	"thread":  {"list", "show"},
	"version": nil,
}

// gitReads are the git commands a sub-agent may run.
var gitReads = []string{"blame", "cat-file", "describe", "diff", "grep", "log", "ls-files", "ls-tree", "merge-base", "rev-list", "rev-parse", "shortlog", "show", "status"}

// Guard decides on the calls of one flai: Commands are the names of its
// commands, so that a word flai on a command line counts as running flai
// only when a command of its follows.
type Guard struct {
	Commands []string
}

// Check says why a call is refused, or "" when it is not.
func (g Guard) Check(e Event) string {
	if e.AgentID == "" {
		return ""
	}
	who := e.AgentType
	if who == "" {
		who = "unnamed"
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		if slices.Contains(MCPReads, tool) {
			return ""
		}
		return fmt.Sprintf("a sub-agent (%s) cannot call %s: it reads, and only the story's agent changes work items and threads or reads the inbox (ADR-0059). Put what you need done, or the question for the designer, in your final message.", who, tool)
	}
	if e.ToolName != "Bash" {
		return ""
	}
	for _, words := range commands(e.ToolInput.Command) {
		if why := g.refuse(words); why != "" {
			return fmt.Sprintf("a sub-agent (%s) cannot run %q: %s (ADR-0060). Put what you need done in your final message; the story's agent does it.", who, strings.Join(words, " "), why)
		}
	}
	return ""
}

var (
	assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
	subcommand = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	shellFlags = regexp.MustCompile(`^-[a-zA-Z]*c[a-zA-Z]*$`)
)

// refuse says why one simple command is refused, or "" when it is not. It
// does not trust the first word to be the program: wrappers such as env,
// sudo, timeout, xargs, and find -exec put it further along, so every word
// is looked at. A word is git run with a subcommand when a plain word
// follows it after git's own flags, and flai run with one when one of its
// commands follows; a shell given -c, or eval, has its script checked too.
func (g Guard) refuse(words []string) string {
	for i, w := range words {
		if assignment.MatchString(w) {
			continue
		}
		rest := words[i+1:]
		switch path.Base(w) {
		case "sh", "bash", "zsh", "dash", "ksh":
			for j, f := range rest {
				if shellFlags.MatchString(f) && j+1 < len(rest) {
					if why := g.script(rest[j+1]); why != "" {
						return why
					}
					break
				}
			}
		case "eval":
			if why := g.script(strings.Join(rest, " ")); why != "" {
				return why
			}
		case "flai", "flai.sh":
			cmd, sub := subcommands(rest, map[string]bool{"--config": true})
			if !slices.Contains(g.Commands, cmd) || slices.Contains(rest, "--help") || slices.Contains(rest, "-h") {
				continue
			}
			if subs, ok := cliReads[cmd]; ok && (subs == nil || slices.Contains(subs, sub)) {
				continue
			}
			return "flai commands that change work items, threads, narratives, or releases are the story's agent's"
		case "git":
			cmd, _ := subcommands(rest, map[string]bool{"-C": true, "-c": true, "--git-dir": true, "--work-tree": true, "--namespace": true})
			if !subcommand.MatchString(cmd) || slices.Contains(gitReads, cmd) {
				continue
			}
			return "git commands that change the worktree, the index, branches, or history are the story's agent's"
		}
	}
	return ""
}

// script checks each simple command of a script a shell is given.
func (g Guard) script(text string) string {
	for _, c := range commands(text) {
		if why := g.refuse(c); why != "" {
			return why
		}
	}
	return ""
}

// subcommands is the first two words after a program's own flags; a flag in
// takesValue consumes the word after it.
func subcommands(words []string, takesValue map[string]bool) (cmd, sub string) {
	var found []string
	for i := 0; i < len(words) && len(found) < 2; i++ {
		w := words[i]
		if strings.HasPrefix(w, "-") {
			if takesValue[w] && len(found) == 0 {
				i++
			}
			continue
		}
		found = append(found, w)
	}
	if len(found) > 0 {
		cmd = found[0]
	}
	if len(found) > 1 {
		sub = found[1]
	}
	return cmd, sub
}

// commands splits a shell command line into simple commands, each as its
// words with quotes removed: at ;, &, |, newlines, and the openings of
// subshells and command substitutions. It is not a shell; it is enough to
// find each program a line runs and the words after it.
func commands(line string) [][]string {
	var out [][]string
	var words []string
	var word strings.Builder
	inWord := false
	var quote rune
	end := func() {
		if inWord {
			words = append(words, word.String())
			word.Reset()
			inWord = false
		}
	}
	split := func() {
		end()
		if len(words) > 0 {
			out = append(out, words)
		}
		words = nil
	}
	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			word.WriteRune(r)
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ';' || r == '&' || r == '|' || r == '\n' || r == '(' || r == ')' || r == '`' || r == '{' || r == '}':
			split()
		case r == '$':
			end()
		case r == ' ' || r == '\t':
			end()
		default:
			word.WriteRune(r)
			inWord = true
		}
	}
	split()
	return out
}
