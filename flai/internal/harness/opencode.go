package harness

import (
	"fmt"
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

// openCodeTakes are the config keys opencode takes: none yet.
var openCodeTakes = map[string]option{}

// openCodeSession is a session ID OpenCode makes, as --session takes one:
// ses_ and its letters and digits.
var openCodeSession = regexp.MustCompile(`^ses_[0-9A-Za-z]+$`)

// Start runs opencode run with the operator's arguments, its events as JSON,
// the story's model as provider/model, the session it resumes, and the
// prompt last, as its message. OpenCode makes its own session IDs, so a new
// run passes none, and a resumed one passes the ses_ ID it is given and
// refuses any other, which OpenCode would not find. Its environment keeps
// OpenCode from updating the host's binary. It works a story alone: a
// strategic role, a provider, a config key, a model OpenCode cannot split,
// and a start in a container are refused, each saying what to do.
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
	if _, err := options(OpenCode, config, openCodeTakes); err != nil {
		return Start{}, err
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
	env = append(env, "OPENCODE_DISABLE_AUTOUPDATE=1")
	return Start{Harness: OpenCode, Argv: argv, Env: env}, nil
}
