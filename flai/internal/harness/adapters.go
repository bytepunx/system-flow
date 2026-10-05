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

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ClaudeCode is Anthropic's Claude Code, run headless (claude -p).
const ClaudeCode = "claude-code"

type claudeCode struct{}

// PermissionPromptTool is the handler claude-code asks whenever a call would
// prompt: permission_prompt on the server named flai that Start passes in
// --mcp-config. Claude Code run headless refuses a write under a .claude/
// folder unless a person or a permission handler approves it; this one asks
// the operator on a thread, or allows at once where the operator turned on
// the auto-approve host action, and denies the rest (S-0257).
const PermissionPromptTool = "mcp__flai__permission_prompt"

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
// MCP server, flai's permission_prompt as its permission handler
// (PermissionPromptTool), and the operator's arguments last, so the handler
// applies with the default host arguments and with the operator's own. A
// planner's session runs as the project's planner definition (S-0208), and
// an orchestrator's as its orchestrator definition (S-0218), each with flai
// guard on its file edits and the same handler, which denies anything
// outside a story's worktree.
func (c claudeCode) Start(r Request, host Host) (Start, error) {
	var model string
	var config map[string]string
	if r.Agent != nil {
		model, config = r.Agent.Model, r.Agent.Config
	}
	env, err := roleEnv(r)
	if err != nil {
		return Start{}, err
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
	// <key> plan <item>, <key> orchestrate, or <key> <story>
	subject := r.Story
	if r.Role != "" {
		subject = strings.TrimSpace(r.Role + " " + r.Item)
	}
	argv := []string{host.Program, "-p", Prompt(r), "--output-format", "stream-json", "--verbose",
		"--mcp-config", string(mcp), "--strict-mcp-config", "--name", strings.TrimSpace(r.Project + " " + subject)}
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
	agents, err := subAgents(r)
	if err != nil {
		return Start{}, err
	}
	if agents != "" {
		argv = append(argv, "--agents", agents)
	}
	if r.Role != "" {
		argv = append(argv, "--agent", claudeCodeStrategic[r.Role])
	}
	argv = append(argv, "--permission-prompt-tool", PermissionPromptTool)
	argv = append(argv, host.Args...)
	return Start{Harness: ClaudeCode, Argv: argv, Env: env}, nil
}

// claudeCodeRoles are the sub-agent definitions, in the project's
// .claude/agents/, that each role runs as (ADR-0059).
var claudeCodeRoles = map[string]string{manifest.RoleExplore: "explorer", manifest.RoleVerify: "verifier"}

// claudeCodeStrategic are the definitions, in the project's .claude/agents/,
// that a strategic agent's session runs as, by its role: the planner's
// (S-0208) and the orchestrator's (S-0218). roleEnv refuses any other role
// before they are looked up.
var claudeCodeStrategic = map[string]string{conventions.RolePlan: "planner", conventions.RoleOrchestrate: "orchestrator"}

// subAgents is the --agents JSON that runs each of the agent's roles on its
// own model (S-0189): the project's definition of the role's sub-agent, its
// front matter and its prompt, with the role's model over the definition's.
// A session's --agents outranks the project's definitions, so the file stays
// the source of everything but the model. For a strategic agent, the planner
// or the orchestrator, it holds that agent's definition as well, with the
// agent's model when it names one, whatever its roles (S-0208, S-0218).
// Empty when nothing is to be passed. A role sub-agents cannot run, because
// it names another harness or config, or has no definition, refuses the
// start, as a strategic agent with no definition does.
func subAgents(r Request) (string, error) {
	out := map[string]map[string]any{}
	if r.Role != "" {
		def := claudeCodeStrategic[r.Role]
		agent, err := readDefinition(r.Root, def)
		if err != nil {
			return "", fmt.Errorf("the %[1]s's session runs as the agent %[1]s, and %[2]w; flai upgrade adds it from the template when it is missing", def, err)
		}
		if r.Agent != nil && r.Agent.Model != "" {
			agent["model"] = r.Agent.Model
		}
		out[def] = agent
	}
	a := r.Agent
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
		agent, err := readDefinition(r.Root, def)
		if err != nil {
			return "", fmt.Errorf("role %s runs the sub-agent %s, and %w", n, def, err)
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

// readDefinition reads the project's .claude/agents/<def>.md as --agents
// takes it. Its error is a clause, for what needed the definition to finish.
func readDefinition(root, def string) (map[string]any, error) {
	path := filepath.Join(root, ".claude", "agents", def+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("its definition could not be read: %w", err)
	}
	agent, err := definition(data)
	if err != nil {
		return nil, fmt.Errorf("its definition %s does not read: %w", path, err)
	}
	return agent, nil
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
// Started as the planner, it has FLAI_ROLE and FLAI_ITEM, and {story} is
// empty (S-0208); started as the orchestrator, FLAI_ROLE alone (S-0218).
type command struct{}

// DefaultHost is nothing: the command has no default, the operator writes it.
func (command) DefaultHost() Host { return Host{} }

// Start replaces the placeholders in the operator's command.
func (command) Start(r Request, host Host) (Start, error) {
	if host.Program == "" {
		return Start{}, fmt.Errorf("no command is set on the host (flai serve agent set -- <program> [args...])")
	}
	role, err := roleEnv(r)
	if err != nil {
		return Start{}, err
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
	env = append(env, role...)
	return Start{Harness: Command, Argv: argv, Env: env}, nil
}
