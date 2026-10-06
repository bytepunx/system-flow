package hostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
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
// actions, its default agent, its shared paths, its MCP token) needs it on for
// that project.

var (
	agentName   = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
	checkName   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,31}$`)
	harnessName = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	// itemLike is an entry flai shared check reads as a work item's ID, not
	// a path.
	itemLike     = regexp.MustCompile(`^[EeSsTt]-\d+$`)
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

// servedProject is what settings.serve and settings.unserve are asked about:
// the project's folder, and its key when it has one, for the journal.
type servedProject struct {
	Root string `json:"root"`
	Key  string `json:"key"`
}

// served is settings.serve (flai serve project add) or settings.unserve
// (flai serve project remove), about the project in the folder given.
func served(what, verb string) spec {
	build := func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		in, e := decode[servedProject](raw)
		if e != nil {
			return nil, "", e
		}
		if !filepath.IsAbs(in.Root) || filepath.Clean(in.Root) != in.Root || strings.ContainsAny(in.Root, "\x00\n\r") {
			return nil, "", bad("%q is not an absolute, clean path", in.Root)
		}
		if in.Key != "" && !checkName.MatchString(in.Key) {
			return nil, "", bad("%q is not a project key", in.Key)
		}
		return []string{"serve", "project", verb, "--", in.Root}, "", nil
	}
	return spec{action: ActionSettings, describe: settingsDone(what), say: sayArgs(build), build: build,
		about: func(raw json.RawMessage) (string, string) {
			in, _ := decode[servedProject](raw)
			return in.Key, in.Root
		}}
}

// sharedEdit is settings.shared's command line: flai shared add or remove
// with the one pattern given.
func sharedEdit(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[struct {
		Action  string `json:"action"`
		Pattern string `json:"pattern"`
	}](raw)
	if e != nil {
		return nil, "", e
	}
	if in.Action != "add" && in.Action != "remove" {
		return nil, "", bad("action is add or remove, not %q", in.Action)
	}
	pattern := strings.TrimSpace(in.Pattern)
	if pattern == "" {
		return nil, "", bad("give the pattern to %s, a path relative to the repository root, such as design/adrs", in.Action)
	}
	if e := argList("the pattern", []string{pattern}, false); e != nil {
		return nil, "", e
	}
	return []string{"shared", in.Action, "--autocommit", "--trailer=" + Trailer, "--", pattern}, "", nil
}

// sharedCheck is settings.shared_check's command line: flai shared check
// with each path given, then the story.
func sharedCheck(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[struct {
		Paths []string `json:"paths"`
		Story string   `json:"story"`
	}](raw)
	if e != nil {
		return nil, "", e
	}
	if len(in.Paths) == 0 && in.Story == "" {
		return nil, "", bad("give paths, the paths or touches entries to check, story, a story's ID whose claim to check, or both")
	}
	args := []string{"shared", "check", "--"}
	for _, p := range in.Paths {
		p = strings.TrimSpace(p)
		switch {
		case p == "":
			return nil, "", bad("an empty entry in paths names no path; give a path or a touches entry, such as docs/users/flai.md")
		case strings.ContainsAny(p, "\x00\n\r"):
			return nil, "", bad("%q in paths holds a line break or NUL", p)
		case itemLike.MatchString(p):
			return nil, "", bad("%q in paths is a work item's ID, not a path; give a story's ID as story", p)
		}
		args = append(args, p)
	}
	if in.Story != "" {
		if e := needStory(in.Story); e != nil {
			return nil, "", e
		}
		args = append(args, in.Story)
	}
	return args, "", nil
}

// manifestChange is what settings.manifest is asked: the strategic agents'
// settings to write, each key's value as JSON, and the keys to remove so
// that their defaults apply (S-0229).
type manifestChange struct {
	Set   map[string]json.RawMessage `json:"set"`
	Unset []string                   `json:"unset"`
}

// manifestEdit is settings.manifest's command line: flai manifest set with
// --unset for each key to remove and key=value for each to write, committed
// with the dashboard's trailer. Each key must be one flai manifest set writes
// and each value a boolean, a number, or a line of text, or for an agent an
// agent object; whether the value suits its key is flai's to say. A refusal
// here names each field and why in Data, as flai's own refusal does.
func manifestEdit(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[manifestChange](raw)
	if e != nil {
		return nil, "", e
	}
	if len(in.Set) == 0 && len(in.Unset) == 0 {
		return nil, "", bad("nothing to change: give set, the keys to write with their values, unset, the keys to remove, or both")
	}
	var problems []manifest.Problem
	refuse := func(key, reason string) { problems = append(problems, manifest.Problem{Field: key, Reason: reason}) }
	inCatalog := func(key string) (manifest.Setting, bool) {
		s, ok := manifest.SettingFor(key)
		if !ok {
			refuse(key, "is not a setting the dashboard writes; give one of the keys settings.get lists under strategic")
		}
		return s, ok
	}
	args := []string{"manifest", "set", "--autocommit", "--trailer=" + Trailer}
	for _, key := range in.Unset {
		if _, ok := inCatalog(key); ok {
			args = append(args, "--unset="+key)
		}
	}
	keys := make([]string, 0, len(in.Set))
	for k := range in.Set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var assignments []string
	for _, key := range keys {
		s, ok := inCatalog(key)
		if !ok {
			continue
		}
		v, reason := settingArg(s, in.Set[key])
		if reason != "" {
			refuse(key, reason)
			continue
		}
		assignments = append(assignments, key+"="+v)
	}
	if len(problems) > 0 {
		return nil, "", &channel.Error{Code: channel.CodeInvalidParams,
			Message: (&manifest.RefusedError{Problems: problems}).Error(),
			Data:    map[string]any{"refused": problems}}
	}
	if len(assignments) > 0 {
		args = append(append(args, "--"), assignments...)
	}
	return args, "", nil
}

// settingArg is a setting's value as flai manifest set reads it from its
// command line: true or false, a number written out in digits, a line of
// text as it stands, or an agent as compact JSON; or why the value is not
// one of these.
func settingArg(s manifest.Setting, raw json.RawMessage) (string, string) {
	if s.Kind == manifest.KindAgent {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		var a manifest.Agent
		if !bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) || dec.Decode(&a) != nil {
			return "", `is an agent, an object such as {"harness":"claude-code","model":"claude-sonnet-5"} with harness, model, config, and roles only; {} unsets it`
		}
		if err := a.Validate(); err != nil {
			return "", err.Error()
		}
		out, err := json.Marshal(a)
		if err != nil {
			return "", err.Error()
		}
		return string(out), ""
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if dec.Decode(&v) != nil {
		return "", "is not a value"
	}
	switch v := v.(type) {
	case bool:
		return strconv.FormatBool(v), ""
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return "", fmt.Sprintf("%s is not a number flai can write", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), ""
	case string:
		if strings.ContainsAny(v, "\x00\n\r") {
			return "", "holds a line break or NUL; give one line"
		}
		return v, ""
	case nil:
		return "", "is null; give a value, or name the key in unset so that its default applies"
	}
	return "", fmt.Sprintf("is a %s setting, given as a boolean, a number, or text, not a list or an object", s.Kind)
}

func settingsSpecs() map[string]spec {
	project := func(what string, b func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)) spec {
		return spec{action: ActionSettings, describe: settingsDone(what), say: sayArgs(b), build: b}
	}
	hostwide := func(what string, b func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)) spec {
		return spec{action: ActionSettings, hostwide: true, describe: settingsDone(what), say: sayArgs(b), build: b}
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
			if shellOnly[in.Action] {
				return nil, "", bad("the %s action is the operator's shell tool, outside the workflow (ADR-0067), and no dashboard turns it on or off: run %s or flai serve disable %s in a shell on the host", in.Action, EnableCommand(in.Action), in.Action)
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
			roles, e := roleArgs(in.Agent)
			return append(args, roles...), "", e
		}),

		// settings.shared: one pattern added to or removed from the project's
		// shared paths, claims.shared in its manifest (S-0295, ADR-0096), by
		// flai shared add or remove with --autocommit, which commits
		// system-flow.yaml on its own. flai exits 1 refusing a pattern that is not valid, one to add that is
		// listed already, or one to remove that is not, and with a manifest it
		// cannot read or write: each is answered as a rule, with its reason.
		"settings.shared": {action: ActionSettings, exits: map[int]int{1: Rule}, describe: settingsDone("the project's shared paths changed"),
			say: sayArgs(sharedEdit), build: sharedEdit},

		// settings.shared_check: whether paths or touches entries, and the
		// entries of a story's claim, lie inside the shared paths, as flai
		// shared check --json answers. It changes nothing, so like
		// settings.get it needs no host action.
		"settings.shared_check": {reads: true, build: sharedCheck},

		// settings.manifest: the strategic agents' settings in the project's
		// system-flow.yaml (S-0229), written by flai manifest set with
		// --autocommit, which checks the manifest as it would be as flai check
		// does. flai exits 1 refusing a value, with {"refused": [{field,
		// reason}]} on standard output: it is answered as Refused, with that
		// in Data, so that the panel can say what is wrong beside each field.
		"settings.manifest": {action: ActionSettings, exits: map[int]int{1: Refused}, describe: settingsDone("the project's strategic agents' settings changed"),
			say: sayArgs(manifestEdit), build: manifestEdit},

		// settings.agent: the agent's name and command, kept for every
		// project. command null removes the command (and with it the name,
		// as flai serve agent clear does). attended_minutes is retired
		// (S-0116, ADR-0043): a dashboard that still sends it is not refused,
		// and it changes nothing.
		"settings.agent": hostwide("the agent's name or command changed", func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
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
				return nil, "", bad("nothing to set: give name or command")
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

		// settings.serve and settings.unserve: a project flai serve serves on
		// the dashboard, registered or no longer (S-0122), as flai serve
		// project add and remove do in a shell. Gated by the settings action
		// of the project served or removed, not of the one the request
		// arrives through: a project to serve has no connection yet, and
		// one to remove may have lost its own.
		"settings.serve":   served("a project served", "add"),
		"settings.unserve": served("a project no longer served", "remove"),

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

// sayArgs journals a settings change as the flai command it ran, which says
// exactly what became what: "flai serve enable push", "flai serve import add
// -- /home/me/git". The trailer and --json are left out; a token never is in
// a command, only in its answer.
func sayArgs(b func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)) func(json.RawMessage) string {
	return func(raw json.RawMessage) string {
		args, _, e := b(channel.Project{}, raw)
		if e != nil {
			return ""
		}
		out := []string{"flai"}
		for _, a := range args {
			if strings.HasPrefix(a, "--trailer=") {
				continue
			}
			out = append(out, a)
		}
		return strings.Join(out, " ")
	}
}
