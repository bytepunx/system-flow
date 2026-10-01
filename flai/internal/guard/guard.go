// Package guard refuses what a sub-agent may not do (ADR-0060): Claude Code
// runs flai guard before each tool call in a project made from the template,
// with the call on standard input, and says which sub-agent makes it. A
// sub-agent reads; it does not change a work item, a thread, or the
// repository's history, because the story's agent alone acts for the story
// (ADR-0059). The story's agent's own calls carry no agent type and are
// never refused.
package guard

import (
	"fmt"
	"path"
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
	// AgentType names the sub-agent making the call; empty for the session's
	// own agent.
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

// Check says why a call is refused, or "" when it is not.
func Check(e Event) string {
	if e.AgentType == "" {
		return ""
	}
	if tool, ok := strings.CutPrefix(e.ToolName, MCPPrefix); ok {
		if slices.Contains(MCPReads, tool) {
			return ""
		}
		return fmt.Sprintf("a sub-agent (%s) cannot call %s: it reads, and only the story's agent changes work items and threads or reads the inbox (ADR-0059). Put what you need done, or the question for the designer, in your final message.", e.AgentType, tool)
	}
	if e.ToolName != "Bash" {
		return ""
	}
	for _, words := range commands(e.ToolInput.Command) {
		if why := refuse(words); why != "" {
			return fmt.Sprintf("a sub-agent (%s) cannot run %q: %s (ADR-0060). Put what you need done in your final message; the story's agent does it.", e.AgentType, strings.Join(words, " "), why)
		}
	}
	return ""
}

// refuse says why one simple command is refused, or "" when it is not.
func refuse(words []string) string {
	words = program(words)
	if len(words) == 0 {
		return ""
	}
	switch name := path.Base(words[0]); name {
	case "sh", "bash", "zsh", "eval":
		for i, w := range words[1:] {
			if name == "eval" || w == "-c" {
				rest := strings.Join(words[i+1:], " ")
				if w == "-c" && i+2 < len(words) {
					rest = words[i+2]
				}
				for _, c := range commands(rest) {
					if why := refuse(c); why != "" {
						return why
					}
				}
				return ""
			}
		}
		return ""
	case "flai", "flai.sh":
		cmd, sub := subcommands(words[1:], map[string]bool{"--config": true})
		subs, ok := cliReads[cmd]
		if ok && (subs == nil || slices.Contains(subs, sub)) {
			return ""
		}
		return "flai commands that change work items, threads, narratives, or releases are the story's agent's"
	case "git":
		cmd, _ := subcommands(words[1:], map[string]bool{"-C": true, "-c": true, "--git-dir": true, "--work-tree": true})
		if cmd == "" || slices.Contains(gitReads, cmd) {
			return ""
		}
		return "git commands that change the worktree, the index, branches, or history are the story's agent's"
	}
	return ""
}

// program drops what runs a command rather than being it: leading
// environment assignments and env, command, exec, nohup, time, and sudo.
func program(words []string) []string {
	for len(words) > 0 {
		w := words[0]
		switch {
		case strings.Contains(w, "=") && !strings.HasPrefix(w, "-") && !strings.Contains(w, "/"):
		case slices.Contains([]string{"env", "command", "exec", "nohup", "time", "sudo"}, w):
		default:
			return words
		}
		words = words[1:]
	}
	return words
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
