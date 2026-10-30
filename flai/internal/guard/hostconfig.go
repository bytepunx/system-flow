package guard

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// The host configuration is the operator's (S-0511, ADR-0220): it holds what
// flai serve may do on this host, and agentEnv drops FLAI_CONFIG from a
// served agent, so a write of an agent's would land in the operator's own
// file. Every session flai serve starts and every sub-agent is refused each
// flai serve, flai config, and flai template command that writes it or turns
// a host action on or off, whatever its role or permissions; the reads stay
// open to this rule, and the operator's own session is not held.

// hostWrites are, for each flai serve sub-command that has its own, the
// sub-commands of it that write the host configuration in every form.
var hostWrites = map[string][]string{
	"agent":      {"set", "clear"},
	"checks":     {"set", "timeout", "clear"},
	"components": {"idle", "max-servers"},
	"notify":     {"set", "clear"},
	"import":     {"add", "remove"},
}

// hostSettings are, for flai serve agent harness and provider, the flags
// that make them write; without one they read.
var hostSettings = map[string][]string{
	"harness":  {"--program", "--guard", "--deny-protected", "--reset"},
	"provider": {"--api", "--base-url", "--key-env", "--spend-key-env", "--model", "--reset"},
}

// forgeSettings are the flags that make flai serve forge write.
var forgeSettings = []string{"--token-env", "--reset"}

// hostValues are the flags, flai's own and those of the flai serve
// sub-commands, that take a value as a word of their own, so that the
// sub-command's own sub-command is the third word that is not a flag.
var hostValues = map[string]bool{"--api": true, "--auto-restarts": true, "--base-url": true, "--command": true, "--config": true, "--guard": true, "--key-env": true, "--model": true, "--name": true, "--program": true, "--spend-key-env": true, "--token-env": true, "--url": true}

// operatorsHost says why no agent writes the host configuration (ADR-0220).
const operatorsHost = "the host configuration, ~/.flai/config.json or the file FLAI_CONFIG names, holds what the operator lets flai serve do on this host, its host actions, agent command, harnesses, providers, checks, forge token variable, reminder channel, and import folders, and flai serve hands an agent no copy of its own, so an agent that wrote it could widen its own permissions or name a program the host runs; only the operator changes it, in a shell or through the dashboard's settings host action (ADR-0220)"

// askForHost ends the refusal of a session's own write of the host
// configuration: what it does instead.
const askForHost = "Ask the operator on a thread with thread_open, naming the setting and why you need it; flai serve status, actions, journal, and the show and list forms read and are open to you."

// hosting are the rules that find a write of the host configuration on a
// command line, and no git command.
var hosting = rules{
	flai: func(cmd, sub string, rest []string) (string, string) {
		if writesHost(cmd, sub, rest) {
			return operatorsHost, ""
		}
		return "", ""
	},
}

// writesHost says whether a flai command, with its subcommand and the words
// after flai, writes the host configuration: flai config set, flai template
// use, flai serve enable and disable, and the flai serve forms ADR-0220
// lists. A setting flag counts whatever its value, as flai counts it.
func writesHost(cmd, sub string, rest []string) bool {
	switch {
	case cmd == "config":
		return sub == "set"
	case cmd == "template":
		return sub == "use"
	case cmd != "serve":
		return false
	}
	setsAny := func(flags []string) bool {
		return slices.ContainsFunc(flags, func(f string) bool { return named(rest, f) })
	}
	switch sub {
	case "enable", "disable":
		return true
	case "forge":
		return setsAny(forgeSettings)
	}
	var third string
	if words := positionals(rest, hostValues); len(words) > 2 {
		third = words[2]
	}
	if slices.Contains(hostWrites[sub], third) {
		return true
	}
	flags, ok := hostSettings[third]
	if sub != "agent" || !ok {
		return false
	}
	return setsAny(flags) || third == "harness" && slices.Contains(rest, "--")
}

// hostConfig decides on a command line that writes the host configuration,
// by a sub-agent or in a session flai serve started: the zero Verdict when it
// writes none or the operator's own session runs it. The orchestrator's
// refusal names the call, so that its activity document logs it.
func (g Guard) hostConfig(c Call) Verdict {
	if c.Kind != Shell || !c.Subagent && !g.served(c) {
		return Verdict{}
	}
	for _, words := range commands(c.Command) {
		if why, _ := g.refuse(words, hosting); why == "" {
			continue
		}
		call := strings.Join(words, " ")
		if c.Subagent {
			return deny(fmt.Sprintf("a sub-agent (%s) cannot run %q: %s. Put the setting you need in your final message.", cmp.Or(c.Agent, "unnamed"), call, operatorsHost))
		}
		v := deny(fmt.Sprintf("%s cannot run %q: %s. %s", g.who(c), call, operatorsHost, askForHost))
		if c.Role == RoleOrchestrate {
			v.Call = call
		}
		return v
	}
	return Verdict{}
}
