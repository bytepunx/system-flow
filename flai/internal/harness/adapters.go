package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
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
	agents, err := subAgents(r.Root, r.Agent)
	if err != nil {
		return Start{}, err
	}
	if agents != "" {
		argv = append(argv, "--agents", agents)
	}
	argv = append(argv, host.Args...)
	return Start{Harness: ClaudeCode, Argv: argv}, nil
}

// claudeCodeRoles are the sub-agent definitions, in the project's
// .claude/agents/, that each role runs as (ADR-0059).
var claudeCodeRoles = map[string]string{manifest.RoleExplore: "explorer", manifest.RoleVerify: "verifier"}

// subAgents is the --agents JSON that runs each of the story's roles on its
// own model (S-0189): the project's definition of the role's sub-agent, its
// front matter and its prompt, with the role's model over the definition's.
// A session's --agents outranks the project's definitions, so the file stays
// the source of everything but the model. Empty when no role names a model.
// A role sub-agents cannot run, because it names another harness or config,
// or has no definition, refuses the start.
func subAgents(root string, a *manifest.Agent) (string, error) {
	if a == nil || len(a.Roles) == 0 {
		return "", nil
	}
	out := map[string]map[string]any{}
	for _, n := range a.RoleNames() {
		role := a.Roles[n]
		def, ok := claudeCodeRoles[n]
		switch {
		case !ok:
			return "", fmt.Errorf("claude-code has no sub-agent for role %q; it runs %s", n, strings.Join(roleNames(), " and "))
		case role.Harness != "" && role.Harness != ClaudeCode:
			return "", fmt.Errorf("claude-code cannot run role %s on harness %q: a sub-agent runs in the story's own session", n, role.Harness)
		case len(role.Config) > 0:
			return "", fmt.Errorf("claude-code takes no config for role %s; it takes a model", n)
		case role.Model == "":
			continue
		}
		path := filepath.Join(root, ".claude", "agents", def+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("role %s runs the sub-agent %s, and its definition could not be read: %w", n, def, err)
		}
		agent, err := definition(data)
		if err != nil {
			return "", fmt.Errorf("role %s: %s: %w", n, path, err)
		}
		agent["model"] = role.Model
		out[def] = agent
	}
	if len(out) == 0 {
		return "", nil
	}
	js, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func roleNames() []string {
	out := make([]string, 0, len(claudeCodeRoles))
	for n := range claudeCodeRoles {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// definition reads a Claude Code sub-agent definition, YAML front matter
// and a prompt, as --agents takes one: every front-matter key but name,
// tools lists as arrays, and the body as prompt.
func definition(data []byte) (map[string]any, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	rest, ok := strings.CutPrefix(text, "---\n")
	if !ok {
		return nil, fmt.Errorf("no front matter")
	}
	front, body, ok := strings.Cut(rest, "\n---\n")
	if !ok {
		return nil, fmt.Errorf("its front matter does not end")
	}
	var fm map[string]any
	if err := yaml.Unmarshal([]byte(front), &fm); err != nil {
		return nil, fmt.Errorf("front matter: %w", err)
	}
	if d, _ := fm["description"].(string); d == "" {
		return nil, fmt.Errorf("it has no description")
	}
	delete(fm, "name")
	for _, k := range []string{"tools", "disallowedTools"} {
		if s, ok := fm[k].(string); ok {
			var list []string
			for _, t := range strings.Split(s, ",") {
				if t = strings.TrimSpace(t); t != "" {
					list = append(list, t)
				}
			}
			fm[k] = list
		}
	}
	fm["prompt"] = strings.TrimSpace(body)
	return fm, nil
}

// command is the operator's own command (S-0079): run as it stands, with
// {story}, {root}, {model}, and {harness} replaced in its arguments. The
// story's config is not interpreted: it is handed over as FLAI_AGENT_CONFIG,
// a JSON object, for the command to read if it wants, and its roles as
// FLAI_AGENT_ROLES, when it has any (S-0189). Started again after
// its question was answered, it has FLAI_ANSWERED, the thread's ID; started
// to commit what a story's worktree holds, FLAI_COMMIT, the worktree (S-0140).
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
	if r.Agent != nil && len(r.Agent.Roles) > 0 {
		roles, err := json.Marshal(r.Agent.Roles)
		if err != nil {
			return Start{}, err
		}
		env = append(env, "FLAI_AGENT_ROLES="+string(roles))
	}
	if r.Answered != "" {
		env = append(env, "FLAI_ANSWERED="+r.Answered)
	}
	if r.Commit != "" {
		env = append(env, "FLAI_COMMIT="+r.Commit)
	}
	return Start{Harness: Command, Argv: argv, Env: env}, nil
}
