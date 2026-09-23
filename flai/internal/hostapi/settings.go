package hostapi

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The host's settings, changed from the dashboard (S-0105). Until then every
// one of them was the operator's to change in a shell on the host (ADR-0029);
// now each is a flai command a method builds, as the dashboard's other writes
// are, gated by the settings host action, which only a shell turns on and
// off. What the host's configuration keeps for every project (the agent's
// command, the harnesses, the checks, the import folders, the dashboard
// token) needs it on for every project; what belongs to one project (its host
// actions, its default agent, its MCP token) needs it on for that project.

var (
	agentName    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
	checkName    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,31}$`)
	harnessName  = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	settingsDone = func(what string) func(any, *channel.Error) (string, string) {
		return func(_ any, err *channel.Error) (string, string) {
			if err != nil {
				return "failed", err.Message
			}
			return "done", what
		}
	}
)

// EnableEverywhere is what the operator runs for a host-wide setting.
func EnableEverywhere(action string) string { return EnableCommand(action) + " --all-projects" }

// argList checks a command given as a list: a program, then its arguments,
// each one line with no NUL. It is run as it stands, never through a shell.
func argList(what string, list []string, needProgram bool) *channel.Error {
	if needProgram && (len(list) == 0 || strings.TrimSpace(list[0]) == "") {
		return bad("%s needs a program", what)
	}
	for _, a := range list {
		if strings.ContainsAny(a, "\x00\n\r") {
			return bad("%s: an argument holds a line break or NUL", what)
		}
	}
	return nil
}

func settingsSpecs() map[string]spec {
	project := func(what string, b func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)) spec {
		return spec{action: ActionSettings, describe: settingsDone(what), build: b}
	}
	hostwide := func(what string, b func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)) spec {
		return spec{action: ActionSettings, hostwide: true, describe: settingsDone(what), build: b}
	}
	return map[string]spec{
		// settings.action: another host action on or off for this project.
		"settings.action": project("a host action turned on or off", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Action string `json:"action"`
				On     *bool  `json:"on"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if _, ok := Actions[in.Action]; !ok {
				return nil, "", bad("%q is not a host action", in.Action)
			}
			if in.Action == ActionSettings {
				return nil, "", bad("the settings action is turned on and off in a shell on the host, and nowhere else")
			}
			if in.On == nil {
				return nil, "", bad("say on: true or on: false")
			}
			if *in.On {
				return []string{"serve", "enable", in.Action}, "", nil
			}
			return []string{"serve", "disable", in.Action}, "", nil
		}),

		// settings.default_agent: the project's default agent, replaced whole,
		// or removed with null; committed as the dashboard's writes are.
		"settings.default_agent": project("the project's default agent changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Agent *manifest.Agent `json:"agent"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Agent.IsZero() {
				return []string{"agent", "clear", "--autocommit", "--trailer=" + Trailer}, "", nil
			}
			if err := in.Agent.Validate(); err != nil {
				return nil, "", bad("%s", err)
			}
			args := []string{"agent", "set", "--replace", "--autocommit", "--trailer=" + Trailer}
			if in.Agent.Harness != "" {
				args = append(args, "--harness="+in.Agent.Harness)
			}
			if in.Agent.Model != "" {
				args = append(args, "--model="+in.Agent.Model)
			}
			for _, k := range in.Agent.ConfigKeys() {
				if strings.TrimSpace(in.Agent.Config[k]) == "" {
					return nil, "", bad("agent config %s has no value", k)
				}
				args = append(args, "--config="+k+"="+in.Agent.Config[k])
			}
			return args, "", nil
		}),

		// settings.agent: the agent's name, attended minutes, and command,
		// kept for every project. command null removes the command (and with
		// it the name and minutes, as flai serve agent clear does).
		"settings.agent": hostwide("the agent's name, minutes, or command changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			var in struct {
				Name     string          `json:"name"`
				Attended int             `json:"attended_minutes"`
				Command  json.RawMessage `json:"command"`
			}
			if e := params(raw, &in); e != nil {
				return nil, "", e
			}
			if string(in.Command) == "null" {
				return []string{"serve", "agent", "clear"}, "", nil
			}
			args := []string{"serve", "agent", "set"}
			if in.Name != "" {
				if !agentName.MatchString(in.Name) {
					return nil, "", bad("%q is not an agent's name: lower-case letters, digits, and -", in.Name)
				}
				args = append(args, "--name="+in.Name)
			}
			if in.Attended != 0 {
				if in.Attended < 1 || in.Attended > 1440 {
					return nil, "", bad("attended minutes are 1 to 1440")
				}
				args = append(args, "--attended-minutes="+strconv.Itoa(in.Attended))
			}
			if len(in.Command) > 0 {
				var cmd []string
				if err := json.Unmarshal(in.Command, &cmd); err != nil {
					return nil, "", bad("command is a list of strings, the program first")
				}
				if e := argList("the command", cmd, true); e != nil {
					return nil, "", e
				}
				args = append(append(args, "--"), cmd...)
			}
			if len(args) == 3 {
				return nil, "", bad("nothing to set: give name, attended_minutes, or command")
			}
			return args, "", nil
		}),

		// settings.harness: what a harness is on this host, and the arguments
		// that say what its agent may do; args null leaves them, [] means none.
		"settings.harness": hostwide("a harness's program or arguments changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			var in struct {
				Name    string          `json:"name"`
				Program string          `json:"program"`
				Args    json.RawMessage `json:"args"`
				Reset   bool            `json:"reset"`
			}
			if e := params(raw, &in); e != nil {
				return nil, "", e
			}
			if !harnessName.MatchString(in.Name) || in.Name == "command" {
				return nil, "", bad("%q is not a harness to set", in.Name)
			}
			args := []string{"serve", "agent", "harness", in.Name}
			if in.Reset {
				return append(args, "--reset"), "", nil
			}
			if in.Program != "" {
				if e := argList("the program", []string{in.Program}, true); e != nil {
					return nil, "", e
				}
				args = append(args, "--program="+in.Program)
			}
			if len(in.Args) > 0 && string(in.Args) != "null" {
				var list []string
				if err := json.Unmarshal(in.Args, &list); err != nil {
					return nil, "", bad("args is a list of strings")
				}
				if e := argList("the arguments", list, false); e != nil {
					return nil, "", e
				}
				args = append(append(args, "--"), list...)
			}
			if len(args) == 4 {
				return nil, "", bad("nothing to set: give program, args, or reset")
			}
			return args, "", nil
		}),

		// settings.check: one named check command, or null to remove it.
		"settings.check": hostwide("a check command changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			var in struct {
				Name    string          `json:"name"`
				Command json.RawMessage `json:"command"`
			}
			if e := params(raw, &in); e != nil {
				return nil, "", e
			}
			if !checkName.MatchString(in.Name) {
				return nil, "", bad("%q is not a check's name", in.Name)
			}
			if len(in.Command) == 0 || string(in.Command) == "null" {
				return []string{"serve", "checks", "clear", in.Name}, "", nil
			}
			var cmd []string
			if err := json.Unmarshal(in.Command, &cmd); err != nil {
				return nil, "", bad("command is a list of strings, the program first")
			}
			if e := argList("the check", cmd, true); e != nil {
				return nil, "", e
			}
			return append([]string{"serve", "checks", "set", "--name=" + in.Name, "--"}, cmd...), "", nil
		}),

		"settings.checks_timeout": hostwide("the checks' time limit changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Minutes int `json:"minutes"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Minutes < 1 || in.Minutes > 1440 {
				return nil, "", bad("the time limit is 1 to 1440 minutes")
			}
			return []string{"serve", "checks", "timeout", strconv.Itoa(in.Minutes)}, "", nil
		}),

		// settings.import: a folder whose repositories the board offers to
		// import, named or no longer.
		"settings.import": hostwide("an import folder added or removed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Folder string `json:"folder"`
				Add    *bool  `json:"add"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !filepath.IsAbs(in.Folder) || filepath.Clean(in.Folder) != in.Folder || strings.ContainsAny(in.Folder, "\x00\n\r") {
				return nil, "", bad("%q is not an absolute, clean path", in.Folder)
			}
			if in.Add == nil {
				return nil, "", bad("say add: true or add: false")
			}
			if *in.Add {
				return []string{"serve", "import", "add", "--", in.Folder}, "", nil
			}
			return []string{"serve", "import", "remove", "--", in.Folder}, "", nil
		}),

		// settings.mcp_token: a new token for this project's HTTP MCP server,
		// which is restarted with it; every agent connected over HTTP must be
		// given it again.
		"settings.mcp_token": project("the project's MCP token rotated", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"mcp", "token", "--rotate"}, "", nil
		}),

		// settings.dashboard_token: a new login token for the dashboard, for
		// every project it serves. The dashboard is not restarted: it takes
		// the new token from the answer and keeps the asking session.
		"settings.dashboard_token": hostwide("the dashboard token rotated", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "token", "--rotate", "--no-restart"}, "", nil
		}),
	}
}

// enabledEverywhere says whether an action is on for every project.
func (h Host) enabledEverywhere(action string) bool { return h.enabled(action, config.AllProjects) }
