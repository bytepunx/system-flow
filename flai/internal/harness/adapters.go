package harness

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ClaudeCode is Anthropic's Claude Code, run headless (claude -p).
const ClaudeCode = "claude-code"

type claudeCode struct{}

// DefaultHost is the operator's default for claude-code: edits in the
// project, any shell command (git, flai, the project's tests), and flai's
// MCP tools; nothing else without asking, and headless there is no one to ask.
func (claudeCode) DefaultHost() Host {
	return Host{Program: "claude", Args: []string{"--permission-mode", "acceptEdits", "--allowedTools", "Bash,mcp__flai"}}
}

var claudeCodeTakes = map[string]option{
	"effort":         {flag: "--effort", value: regexp.MustCompile(`^(low|medium|high|xhigh|max)$`), says: "low, medium, high, xhigh, or max"},
	"max_budget_usd": {flag: "--max-budget-usd", value: regexp.MustCompile(`^\d{1,5}(\.\d{1,2})?$`), says: "an amount of dollars, such as 5 or 12.50"},
	"fallback_model": {flag: "--fallback-model", value: regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@+-]{0,127}$`), says: "a model's name"},
}

// Start runs claude -p with the prompt, the story's model and options, flai's
// MCP server, and the operator's arguments last.
func (c claudeCode) Start(r Request, host Host) (Start, error) {
	var model string
	var config map[string]string
	if r.Agent != nil {
		model, config = r.Agent.Model, r.Agent.Config
	}
	opts, err := options(ClaudeCode, config, claudeCodeTakes)
	if err != nil {
		return Start{}, err
	}
	def := c.DefaultHost()
	if host.Program == "" {
		host.Program = def.Program
	}
	if host.Args == nil {
		host.Args = def.Args
	}
	// The agent reaches flai through this flai's own MCP server on stdio,
	// whatever the project's .mcp.json says, which a headless session would
	// not have approved, and no other: not the operator's own connectors and
	// plugins, which S-0104's trial found loaded otherwise. The operator adds
	// others with --mcp-config in the harness's arguments. The server is told
	// the agent's name as an argument: the harness's own settings can put a
	// FLAI_AGENT of their own into the servers it starts, which made every
	// agent flai serve started one name (S-0114, I-0037).
	args := []string{"mcp"}
	if r.Name != "" {
		args = append(args, "--agent", r.Name)
	}
	mcp, err := json.Marshal(map[string]any{"mcpServers": map[string]any{"flai": map[string]any{"type": "stdio", "command": r.Flai, "args": args}}})
	if err != nil {
		return Start{}, err
	}
	argv := []string{host.Program, "-p", Prompt(r), "--output-format", "stream-json", "--verbose",
		"--mcp-config", string(mcp), "--strict-mcp-config", "--name", strings.TrimSpace(r.Project + " " + r.Story)}
	// A session of its own, so that an agent that ended waiting for an
	// answer goes on with what it knew.
	switch {
	case r.Session != "" && r.Answered != "":
		argv = append(argv, "--resume", r.Session)
	case r.Session != "":
		argv = append(argv, "--session-id", r.Session)
	}
	if model != "" {
		argv = append(argv, "--model", model)
	}
	argv = append(argv, opts...)
	argv = append(argv, host.Args...)
	return Start{Harness: ClaudeCode, Argv: argv}, nil
}

// command is the operator's own command (S-0079): run as it stands, with
// {story}, {root}, {model}, and {harness} replaced in its arguments. The
// story's config is not interpreted: it is handed over as FLAI_AGENT_CONFIG,
// a JSON object, for the command to read if it wants. Started again after
// its question was answered, it has FLAI_ANSWERED, the thread's ID.
type command struct{}

// DefaultHost is nothing: the command has no default, the operator writes it.
func (command) DefaultHost() Host { return Host{} }

// Start replaces the placeholders in the operator's command.
func (command) Start(r Request, host Host) (Start, error) {
	if host.Program == "" {
		return Start{}, fmt.Errorf("no command is set on the host (flai serve agent set -- <program> [args...])")
	}
	var model, harness string
	config := map[string]string{}
	if r.Agent != nil {
		model, harness = r.Agent.Model, r.Agent.Harness
		for k, v := range r.Agent.Config {
			config[k] = v
		}
	}
	if harness == "" {
		harness = Command
	}
	rep := strings.NewReplacer("{story}", r.Story, "{root}", r.Root, "{model}", model, "{harness}", harness)
	argv := []string{rep.Replace(host.Program)}
	for _, a := range host.Args {
		argv = append(argv, rep.Replace(a))
	}
	cfg, err := json.Marshal(config)
	if err != nil {
		return Start{}, err
	}
	env := []string{"FLAI_MODEL=" + model, "FLAI_HARNESS=" + harness, "FLAI_AGENT_CONFIG=" + string(cfg)}
	if r.Answered != "" {
		env = append(env, "FLAI_ANSWERED="+r.Answered)
	}
	return Start{Harness: Command, Argv: argv, Env: env}, nil
}
