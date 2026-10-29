package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mcphttp"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/usage"
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

// Capabilities are all of them: claude resumes a session, runs flai guard as
// its hook, asks through permission_prompt, takes a model per role in
// --agents, runs a strategic agent from its definition, and reports its cost
// in stream-json.
func (claudeCode) Capabilities() Capabilities {
	return Capabilities{Resume: true, Guard: true, Ask: true, RoleModels: true, StrategicAgents: true, Cost: true}
}

// Reader reads claude's stream-json.
func (claudeCode) Reader() usage.Reader { return usage.ClaudeCode() }

var claudeCodeTakes = map[string]option{
	"effort":         {flag: "--effort", value: regexp.MustCompile(`^(low|medium|high|xhigh|max)$`), says: "low, medium, high, xhigh, or max"},
	"max_budget_usd": {flag: "--max-budget-usd", value: regexp.MustCompile(`^\d{1,5}(\.\d{1,2})?$`), says: "an amount of dollars, such as 5 or 12.50"},
	"fallback_model": {flag: "--fallback-model", value: regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@+-]{0,127}$`), says: "a model's name"},
}

// Start runs claude -p with the prompt, the story's model and options, flai's
// MCP server, flai's permission_prompt as its permission handler
// (PermissionPromptTool), and the operator's arguments last, so the handler
// applies with the default host arguments and with the operator's own. A
// planner's session runs as the project's planner definition (S-0208), an
// orchestrator's as its orchestrator definition (S-0218), an analyzer's as
// its analyzer definition (S-0223), a curator's as its curator
// definition (S-0431), and a defined role's as the definition flai
// generates from its role file, passed in --agents under its name (S-0417),
// each with flai guard on its file
// edits and the same handler, which denies anything outside a story's
// worktree. Every session whose agent names a provider calls it rather
// than Anthropic (claudeCodeProvider). Started by pilot in a container, it
// reaches the host's HTTP MCP server through a file the Start's Files hold,
// and keeps its configuration and sessions under the request's cache
// (S-0507).
func (claudeCode) Start(r Request, host Host) (Start, error) {
	var model string
	var config map[string]string
	if r.Agent != nil {
		model, config = r.Agent.Model, r.Agent.Config
	}
	env, err := roleEnv(r)
	if err != nil {
		return Start{}, err
	}
	penv, from, err := claudeCodeProvider(r)
	if err != nil {
		return Start{}, err
	}
	env = append(env, penv...)
	opts, err := options(ClaudeCode, config, claudeCodeTakes)
	if err != nil {
		return Start{}, err
	}
	host = ClaudeCodeHost(host)
	mcp, files, err := claudeCodeMCP(r)
	if err != nil {
		return Start{}, err
	}
	// In a container, claude keeps its configuration and its sessions in the
	// mounted cache, so that a later container resumes the session; claude
	// makes the folder when it is missing, and what starts the container
	// links the login into it (S-0507).
	if r.Cache != "" {
		env = append(env, "CLAUDE_CONFIG_DIR="+filepath.Join(r.Cache, "claude"))
	}
	// <key> plan <item>, <key> orchestrate, <key> analyze [<focus>],
	// <key> curate, <key> <name> [<item>] for a defined role, or <key> <story>
	subject := r.Story
	switch {
	case r.RoleFile != nil:
		subject = strings.TrimSpace(r.RoleFile.Name + " " + r.Item)
	case r.Role != "":
		subject = strings.Join(strings.Fields(r.Role+" "+r.Item+" "+r.Focus), " ")
	}
	argv := append(headless(host.Program, Prompt(r), mcp), "--name", strings.TrimSpace(r.Project+" "+subject))
	// A session of its own, so that an agent that ended waiting for an
	// answer, or whose story was sent back, goes on with what it knew.
	switch {
	case r.resumes():
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
	switch {
	case r.RoleFile != nil:
		argv = append(argv, "--agent", r.RoleFile.Name)
	case r.Role != "":
		argv = append(argv, "--agent", claudeCodeStrategic[r.Role])
	}
	argv = append(argv, permissionArgs(host)...)
	return Start{Harness: ClaudeCode, Argv: argv, Env: env, EnvFrom: from, Files: files}, nil
}

// claudeCodeMCP is the --mcp-config of r's session, and the files it needs
// written. The agent reaches flai through flai's MCP server alone, whatever
// the project's .mcp.json says, which a headless session would not have
// approved: not the operator's own connectors and plugins, which S-0104's
// trial found loaded otherwise. The operator adds others with --mcp-config
// in the harness's arguments. On the host the server is this flai's own on
// stdio, told the agent's name as an argument: the harness's own settings
// can put a FLAI_AGENT of their own into the servers it starts, which made
// every agent flai serve started one name (S-0114, I-0037). In a container,
// where MCPURLEnv is set, it is the host's HTTP server, named in a file
// under the cache that the starter writes, pilot (S-0507).
func claudeCodeMCP(r Request) (string, []File, error) {
	if r.MCPURLEnv == "" {
		args := []string{"mcp"}
		if r.Name != "" {
			args = append(args, "--agent", r.Name)
		}
		mcp, err := mcpConfig(r.Flai, args)
		return mcp, nil, err
	}
	if err := checkContainer(r); err != nil {
		return "", nil, err
	}
	data, err := httpMCPConfig(r)
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(r.Cache, "mcp.json")
	return path, []File{{Path: path, Data: data, Mode: 0o600}}, nil
}

// envName is the name of an environment variable, as ${NAME} names it.
var envName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// checkContainer refuses a container start that cannot reach the host's
// HTTP MCP server: one that names no token's variable or no cache, or a
// variable that is not a variable's name.
func checkContainer(r Request) error {
	switch {
	case r.MCPTokenEnv == "":
		return fmt.Errorf("the start request names the MCP URL's variable %s but not its token's; set mcp_token_env to the variable that holds the HTTP server's bearer token", r.MCPURLEnv)
	case r.Cache == "":
		return fmt.Errorf("the start request names the MCP URL's variable %s but no cache; set cache to the mounted folder that keeps the harness's configuration and sessions", r.MCPURLEnv)
	case !filepath.IsAbs(r.Cache):
		return fmt.Errorf("the start request's cache %q is not an absolute path; set cache to the mounted folder's absolute path", r.Cache)
	}
	for _, v := range []struct{ field, name string }{{"mcp_url_env", r.MCPURLEnv}, {"mcp_token_env", r.MCPTokenEnv}, {"mcp_host_env", r.MCPHostEnv}} {
		if v.name != "" && !envName.MatchString(v.name) {
			return fmt.Errorf("the start request's %s %q is not a variable's name; name the variable with letters, digits, and underscores", v.field, v.name)
		}
	}
	return nil
}

// httpMCPConfig is the MCP file that names the host's HTTP server as the
// server named flai: its URL, its bearer token, and the Host it listens
// under each as ${NAME} of the variable r names, which claude expands as it
// reads the file, so that no token is ever written; and the agent's name in
// X-Flai-Agent, which the server names the session by.
func httpMCPConfig(r Request) ([]byte, error) {
	headers := map[string]string{"Authorization": "Bearer ${" + r.MCPTokenEnv + "}"}
	if r.Name != "" {
		headers[mcphttp.AgentHeader] = r.Name
	}
	if r.MCPHostEnv != "" {
		headers["Host"] = "${" + r.MCPHostEnv + "}"
	}
	mcp, err := json.Marshal(map[string]any{"mcpServers": map[string]any{"flai": map[string]any{"type": "http", "url": "${" + r.MCPURLEnv + "}", "headers": headers}}})
	if err != nil {
		return nil, fmt.Errorf("the MCP configuration for the host's HTTP server: %w", err)
	}
	return mcp, nil
}

// claudeCodeModels are the variables that set what claude's model aliases
// mean, by the alias a provider's models name.
var claudeCodeModels = []struct{ alias, env string }{
	{"opus", "ANTHROPIC_DEFAULT_OPUS_MODEL"},
	{"sonnet", "ANTHROPIC_DEFAULT_SONNET_MODEL"},
	{"haiku", "ANTHROPIC_DEFAULT_HAIKU_MODEL"},
}

// claudeCodeProvider is what claude needs to call the agent's provider
// rather than Anthropic (ADR-0129): its base URL, an empty
// ANTHROPIC_API_KEY so that a key flai serve inherited is not sent to the
// gateway, and what each alias the provider's models name means; and, to
// be copied from flai serve's environment as the child is spawned, the
// provider's key as ANTHROPIC_AUTH_TOKEN. Nothing when the agent names no
// provider. claude speaks anthropic-messages alone, so a provider that
// speaks another api is refused.
func claudeCodeProvider(r Request) ([]string, map[string]string, error) {
	p := r.Provider
	if p == nil {
		return nil, nil, nil
	}
	if p.API != manifest.APIAnthropicMessages {
		return nil, nil, fmt.Errorf("the agent's provider %[1]s speaks %[2]s, and claude-code speaks %[3]s alone; if the gateway speaks %[3]s too, say so with flai serve agent provider %[1]s --api %[3]s, or name another provider in the agent", r.ProviderName, p.API, manifest.APIAnthropicMessages)
	}
	env := []string{"ANTHROPIC_BASE_URL=" + p.BaseURL, "ANTHROPIC_API_KEY="}
	for _, m := range claudeCodeModels {
		if model, ok := p.Models[m.alias]; ok {
			env = append(env, m.env+"="+model)
		}
	}
	return env, map[string]string{"ANTHROPIC_AUTH_TOKEN": p.KeyEnv}, nil
}

// ClaudeCodeHost is host with claude-code's defaults where the operator set
// none: the program, and the arguments when they are nil.
func ClaudeCodeHost(host Host) Host {
	def := claudeCode{}.DefaultHost()
	if host.Program == "" {
		host.Program = def.Program
	}
	if host.Args == nil {
		host.Args = def.Args
	}
	return host
}

// ClaudeCodeCheck is the headless claude -p that checks a Claude Code
// against permission_prompt (ADR-0106): the prompt on model, with flai's MCP
// server, flai run with mcpArgs, as its only server, and the permission
// arguments a story's agent gets from host.
func ClaudeCodeCheck(host Host, flai string, mcpArgs []string, model, prompt string) ([]string, error) {
	host = ClaudeCodeHost(host)
	mcp, err := mcpConfig(flai, mcpArgs)
	if err != nil {
		return nil, err
	}
	argv := append(headless(host.Program, prompt, mcp), "--model", model)
	return append(argv, permissionArgs(host)...), nil
}

// headless is how every claude -p flai runs begins: the prompt, its events
// as stream-json, and the MCP servers of mcp and no others.
func headless(program, prompt, mcp string) []string {
	return []string{program, "-p", prompt, "--output-format", "stream-json", "--verbose", "--mcp-config", mcp, "--strict-mcp-config"}
}

// mcpConfig is the --mcp-config that starts flai's own MCP server on stdio,
// the program flai with args, as the server named flai.
func mcpConfig(flai string, args []string) (string, error) {
	mcp, err := json.Marshal(map[string]any{"mcpServers": map[string]any{"flai": map[string]any{"type": "stdio", "command": flai, "args": args}}})
	if err != nil {
		return "", fmt.Errorf("the MCP configuration for flai's server: %w", err)
	}
	return string(mcp), nil
}

// permissionArgs are what a claude-code session may do: flai's
// permission_prompt as its handler (PermissionPromptTool), then the
// operator's arguments, so that the handler applies with the default host
// arguments and with the operator's own.
func permissionArgs(host Host) []string {
	return append([]string{"--permission-prompt-tool", PermissionPromptTool}, host.Args...)
}

// claudeCodeRoles are the sub-agent definitions, in the project's
// .claude/agents/, that each role runs as (ADR-0059).
var claudeCodeRoles = map[string]string{manifest.RoleExplore: "explorer", manifest.RoleVerify: "verifier"}

// claudeCodeStrategic are the definitions, in the project's .claude/agents/,
// that a strategic agent's session runs as, by its role: the planner's
// (S-0208), the orchestrator's (S-0218), the analyzer's (S-0223), and the
// curator's (S-0431). roleEnv refuses any other role before they are looked
// up.
var claudeCodeStrategic = map[string]string{conventions.RolePlan: "planner", conventions.RoleOrchestrate: "orchestrator", conventions.RoleAnalyze: "analyzer", conventions.RoleCurate: "curator"}

// subAgents is the --agents JSON that runs each of the agent's roles on its
// own model (S-0189): the project's definition of the role's sub-agent, its
// front matter and its prompt, with the role's model over the definition's.
// A session's --agents outranks the project's definitions, so the file stays
// the source of everything but the model. For a strategic agent, the
// planner, the orchestrator, the analyzer, or the curator, it holds that
// agent's definition as well, with the agent's model when it names one,
// whatever its roles (S-0208, S-0218, S-0223, S-0431). For a role the
// project defines, it holds the definition flai generates from the role's
// file under its name (roleDefinition), and nothing is read from or written
// under .claude/ for it (S-0417, ADR-0207 §4).
// Empty when nothing is to be passed. A role sub-agents cannot run, because
// it names another harness or config, or has no definition, refuses the
// start, as a strategic agent with no definition does.
func subAgents(r Request) (string, error) {
	out := map[string]map[string]any{}
	switch {
	case r.RoleFile != nil:
		agent, err := roleDefinition(r.RoleFile, r.Agent)
		if err != nil {
			return "", err
		}
		out[r.RoleFile.Name] = agent
	case r.Role != "":
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

// inherit is a definition's model when it takes the session's.
const inherit = "inherit"

// roleDefinition is the definition flai generates from a role file, as
// --agents takes one (ADR-0207 §4): its definition's description and tools,
// its body without the marker's line as its prompt, and its model, the
// agent's when it names one, else the file's, left out when that is
// inherit, so that the session's holds. A file whose definition has no
// description or no tools, or whose body is empty, is refused, naming flai
// check, which reports it.
func roleDefinition(f *roles.Role, a *manifest.Agent) (map[string]any, error) {
	d := f.Definition
	prompt := roleBody(f.Body)
	var lacks []string
	if strings.TrimSpace(d.Description) == "" {
		lacks = append(lacks, "definition.description")
	}
	if len(d.Tools) == 0 {
		lacks = append(lacks, "definition.tools")
	}
	if prompt == "" {
		lacks = append(lacks, "body")
	}
	if len(lacks) > 0 {
		return nil, fmt.Errorf("the role %s runs as the agent %s, which flai generates from its file %s, and the file has no %s; flai check reports what it lacks", f.Name, f.Name, f.Path, strings.Join(lacks, " and no "))
	}
	agent := map[string]any{"description": d.Description, "tools": slices.Clone(d.Tools), "prompt": prompt}
	model := d.Model
	if a != nil && a.Model != "" {
		model = a.Model
	}
	if model != "" && model != inherit {
		agent["model"] = model
	}
	return agent, nil
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
// {story}, {root}, {model}, {harness}, and {provider} replaced in its
// program and arguments, {provider} by the name of the agent's provider, or
// nothing when it names none (ADR-0129). The
// story's config is not interpreted: it is handed over as FLAI_AGENT_CONFIG,
// a JSON object, for the command to read if it wants, and its roles as
// FLAI_AGENT_ROLES, when it has any (S-0189). Started again after
// its question was answered, it has FLAI_ANSWERED, the thread's ID; started
// to commit what a story's worktree holds, FLAI_COMMIT, the worktree (S-0140).
// Started as the planner, it has FLAI_ROLE and FLAI_ITEM, and {story} is
// empty (S-0208); started as the orchestrator, FLAI_ROLE alone (S-0218);
// started as the analyzer, FLAI_ROLE, and FLAI_FOCUS when it was asked for
// one (S-0223); started as the curator, FLAI_ROLE, and FLAI_SINCE when its
// run has a since (S-0431); started as a role a project defines, FLAI_AGENT,
// FLAI_ROLE, FLAI_ITEM when its run has an item, and FLAI_INPUTS, a JSON
// array of the run's inputs in order, each {name, value}, when it has any
// (S-0417), since the command has no prompt to say it. When
// the agent names a provider, it has FLAI_PROVIDER, its name,
// FLAI_PROVIDER_API, FLAI_PROVIDER_BASE_URL, and FLAI_PROVIDER_KEY_ENV, the
// name of the variable that holds its key, which the command reads itself:
// flai hands it no key, and refuses no api (ADR-0129). It takes no part in
// a container's start, and ignores the request's container fields (S-0507).
type command struct{}

// DefaultHost is nothing: the command has no default, the operator writes it.
func (command) DefaultHost() Host { return Host{} }

// Capabilities are none: flai can neither resume, guard, ask, nor read the
// operator's command.
func (command) Capabilities() Capabilities { return Capabilities{} }

// Reader reads nothing: flai does not know what the operator's command logs.
func (command) Reader() usage.Reader { return unread{} }

// unread is a reader of a log flai cannot read: no line of it is a record.
type unread struct{}

// Read is no event.
func (unread) Read([]byte) ([]usage.Event, bool) { return nil, false }

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
	provider := ""
	if r.Provider != nil {
		provider = r.ProviderName
	}
	rep := strings.NewReplacer("{story}", r.Story, "{root}", r.Root, "{model}", model, "{harness}", harness, "{provider}", provider)
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
	if r.Focus != "" {
		env = append(env, "FLAI_FOCUS="+r.Focus)
	}
	if r.Since != "" {
		env = append(env, "FLAI_SINCE="+r.Since)
	}
	if len(r.Inputs) > 0 {
		inputs, err := json.Marshal(r.Inputs)
		if err != nil {
			return Start{}, err
		}
		env = append(env, "FLAI_INPUTS="+string(inputs))
	}
	if p := r.Provider; p != nil {
		env = append(env, "FLAI_PROVIDER="+r.ProviderName, "FLAI_PROVIDER_API="+p.API, "FLAI_PROVIDER_BASE_URL="+p.BaseURL, "FLAI_PROVIDER_KEY_ENV="+p.KeyEnv)
	}
	return Start{Harness: Command, Argv: argv, Env: env}, nil
}
