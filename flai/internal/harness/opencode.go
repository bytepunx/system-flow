package harness

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// OpenCode is OpenCode, run headless (opencode run), as agent-adapters.md
// § OpenCode reads it from OpenCode's source at v1.19.5 (S-0406).
const OpenCode = "opencode"

// openCode runs a story's agent on OpenCode, with no guard and no ask, so
// that ADR-0131 holds it to the story's own agent and to protected writes
// refused.
type openCode struct{}

// DefaultHost is opencode with no arguments of flai's: what the agent may do
// is OpenCode's permission configuration, and --auto, which approves every
// ask, is the operator's to add, since ADR-0131 refuses a protected write on
// a harness that cannot ask.
func (openCode) DefaultHost() Host {
	return Host{Program: "opencode", Args: []string{}}
}

// Capabilities are resume alone: opencode run --session goes on in a session,
// and flai has no guard, no ask, no role models, no strategic definitions,
// and no cost reader for it yet (ADR-0131).
func (openCode) Capabilities() Capabilities { return Capabilities{Resume: true} }

// Reader reads nothing yet: usage on OpenCode comes with S-0410.
func (openCode) Reader() usage.Reader { return unread{} }

// openCodeTakes are the config keys opencode takes: effort, as --variant,
// the model variant OpenCode names by its reasoning effort for the model's
// provider (provider/transform.ts at v1.19.5). A variant the model lacks
// leaves its default effort.
var openCodeTakes = map[string]option{
	"effort": {flag: "--variant", value: regexp.MustCompile(`^(none|minimal|low|medium|high|xhigh|max)$`), says: "none, minimal, low, medium, high, xhigh, or max"},
}

// openCodeSession is a session ID OpenCode makes, as --session takes one:
// ses_ and its letters and digits.
var openCodeSession = regexp.MustCompile(`^ses_[0-9A-Za-z]+$`)

// Start runs opencode run with the operator's arguments, its events as JSON,
// the story's model as provider/model, its effort as --variant, the session
// it resumes, and the prompt last, as its message. OpenCode makes its own
// session IDs, so a new run passes none, and a resumed one passes the ses_
// ID it is given and refuses any other, which OpenCode would not find. Its
// environment keeps OpenCode from updating the host's binary and gives it
// its configuration (ADR-0260): flai's MCP server in
// OPENCODE_CONFIG_CONTENT, and, when the host denies the protected paths,
// their deny rules in OPENCODE_PERMISSION, read over the project's own
// opencode.json and opencode.jsonc at r.Root. It works a story alone: a
// strategic role, a provider, a config key it does not take, a model
// OpenCode cannot split, and a start in a container are refused, each
// saying what to do.
func (openCode) Start(r Request, host Host) (Start, error) {
	var model string
	var config map[string]string
	provider := r.ProviderName
	if r.Agent != nil {
		model, config = r.Agent.Model, r.Agent.Config
		if r.Agent.Provider != "" {
			provider = r.Agent.Provider
		}
	}
	role := r.Role
	if r.RoleFile != nil {
		role = r.RoleFile.Name
	}
	switch {
	case role != "":
		return Start{}, fmt.Errorf("opencode cannot run the strategic role %s: flai runs a strategic agent from its definition, and has none for opencode yet; run the role on claude-code", role)
	case provider != "" || r.Provider != nil:
		return Start{}, fmt.Errorf("opencode cannot call the agent's provider %s yet: OpenCode calls the provider its model names, with its own login; leave provider out of the story's agent and name the model as provider/model, or run the story on claude-code", provider)
	case r.MCPURLEnv != "":
		return Start{}, fmt.Errorf("opencode cannot run in a container yet: it reaches flai's MCP server on the host alone; run the story on harness opencode on the host, or set config inner to claude-code")
	}
	env, err := roleEnv(r)
	if err != nil {
		return Start{}, err
	}
	flags, err := options(OpenCode, config, openCodeTakes)
	if err != nil {
		return Start{}, err
	}
	content, err := openCodeConfigContent(r)
	if err != nil {
		return Start{}, err
	}
	env = append(env, "OPENCODE_DISABLE_AUTOUPDATE=1", "OPENCODE_CONFIG_CONTENT="+content)
	if host.DenyProtected {
		if r.Root == "" {
			return Start{}, fmt.Errorf("opencode cannot deny the protected paths: the start names no project root, where flai reads the project's opencode.json and opencode.jsonc to keep its permission; name the project's root in the start")
		}
		permission, err := openCodePermission(r.Root, os.DirFS(r.Root))
		if err != nil {
			return Start{}, err
		}
		env = append(env, "OPENCODE_PERMISSION="+permission)
	}
	def := openCode{}.DefaultHost()
	if host.Program == "" {
		host.Program = def.Program
	}
	if host.Args == nil {
		host.Args = def.Args
	}
	argv := append([]string{host.Program, "run"}, host.Args...)
	argv = append(argv, "--format", "json")
	if model != "" {
		if p, id, ok := strings.Cut(model, "/"); !ok || p == "" || id == "" {
			return Start{}, fmt.Errorf("opencode takes the model as provider/model, which it splits at the first /, and the agent's model is %q; name its provider and model so, such as anthropic/claude-opus-5-5", model)
		}
		argv = append(argv, "--model", model)
	}
	argv = append(argv, flags...)
	if r.resumes() {
		if !openCodeSession.MatchString(r.Session) {
			return Start{}, fmt.Errorf("opencode resumes a session by the ses_ ID OpenCode made for it, and the start names %q, which it would not find; restart the story's agent in a new session with flai serve agent restart %s", r.Session, r.Story)
		}
		argv = append(argv, "--session", r.Session)
	}
	prompt, err := Prompt(r)
	if err != nil {
		return Start{}, err
	}
	argv = append(argv, prompt)
	return Start{Harness: OpenCode, Argv: argv, Env: env}, nil
}

// Session is the ses_ ID of the session an opencode run --format json ran
// in: the sessionID of its first event, the log's first line that is a JSON
// object with a type, since each event names the run's root session. It is
// false when no line is an event, or the first names no ID OpenCode makes.
func (openCode) Session(log []byte) (string, bool) {
	for line := range bytes.Lines(log) {
		var event struct {
			Type      string `json:"type"`
			SessionID string `json:"sessionID"`
		}
		if json.Unmarshal(line, &event) != nil || event.Type == "" {
			continue
		}
		if !openCodeSession.MatchString(event.SessionID) {
			return "", false
		}
		return event.SessionID, true
	}
	return "", false
}
