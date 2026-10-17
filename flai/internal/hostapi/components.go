package hostapi

// The component library's methods (S-0437, S-0438, ADR-0163):
// components.list, the catalogue flai components list prints, and
// components.status, how each sub-project's dev server stands, are reads;
// components.start and components.stop are writes under the components host
// action, answered in flai serve's own process, which supervises the dev
// servers, through Host.Components. Nothing here starts or ends a process.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ComponentServer is how one sub-project's component dev server stands, as
// components.status and components.start answer it.
type ComponentServer struct {
	// Project is the sub-project's name in the manifest.
	Project string `json:"project"`
	// Running is whether its dev server runs and answers.
	Running bool `json:"running"`
	// Starting is whether its dev server was started and does not answer yet.
	Starting bool `json:"starting,omitempty"`
	// AlreadyRunning is set by components.start when it answered the dev
	// server that ran rather than start one.
	AlreadyRunning bool `json:"already_running,omitempty"`
	// URL is the harness page's address, http://127.0.0.1:<port>/__flai/components/.
	URL string `json:"url,omitempty"`
	// Port is the loopback port the dev server listens on.
	Port int `json:"port,omitempty"`
	// PID is the process ID of the sub-project's vite.
	PID int `json:"pid,omitempty"`
	// Since is when vite was started, in RFC 3339, UTC.
	Since string `json:"since,omitempty"`
	// LastError is why its last dev server failed to start or exited on its
	// own, until one starts again.
	LastError string `json:"last_error,omitempty"`
}

// ComponentStopped is what components.stop answers.
type ComponentStopped struct {
	// Project is the sub-project's name in the manifest.
	Project string `json:"project"`
	// Stopped is whether a dev server ran and was stopped; false when none ran.
	Stopped bool `json:"stopped"`
	// PID is the process ID of the vite that was stopped.
	PID int `json:"pid,omitempty"`
}

// ComponentsStatus is what components.status answers.
type ComponentsStatus struct {
	// Enabled is whether the components host action is on for the project.
	Enabled bool `json:"enabled"`
	// Enable is what the operator runs on the host to turn it on.
	Enable string `json:"enable"`
	// IdleSeconds is components.idle: a dev server without a
	// components.status for this long is stopped. Absent when flai serve
	// runs no dev servers here.
	IdleSeconds int `json:"idle_seconds,omitempty"`
	// Servers are the supported sub-projects', or the named one's, in the
	// manifest's order.
	Servers []ComponentServer `json:"servers"`
}

// ComponentsProject is one sub-project in components.list's answer, as
// flai components list --json prints it.
type ComponentsProject struct {
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Kind       string                 `json:"kind"`
	Supported  bool                   `json:"supported"`
	Setup      string                 `json:"setup,omitempty"`
	Components []components.Component `json:"components"`
}

// ComponentServers are the sub-projects' component dev servers flai serve
// supervises, one per project and sub-project (ADR-0163).
type ComponentServers interface {
	// Start starts the dev server of sub, a sub-project of the project at
	// root whose manifest names it project, and answers it once it answers;
	// when one runs, or is starting, it answers that one instead.
	Start(ctx context.Context, root, project string, sub manifest.Project) (ComponentServer, error)
	// Stop stops the dev server of the sub-project named sub, when one runs.
	Stop(root, sub string) (ComponentStopped, error)
	// Status is how the dev server of the sub-project named sub stands, and
	// marks it used, so that it is not stopped as idle.
	Status(root, sub string) ComponentServer
	// Idle is how long a dev server may go without a Status before it is
	// stopped: components.idle.
	Idle() time.Duration
}

// nodeRunner runs node for components.list; tests replace it.
var nodeRunner components.Runner = components.System{}

// noComponentServers refuses a start or a stop asked of a flai that runs no
// dev servers.
func noComponentServers() *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no component dev servers: flai serve runs them, so ask through the dashboard flai serve answers, or run flai components serve in a shell on the host"}
}

// componentsManifest is the manifest of the project p.
func componentsManifest(p channel.Project) (manifest.Manifest, *channel.Error) {
	m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
	if err != nil {
		return manifest.Manifest{}, failed(fmt.Errorf("cannot read the project's manifest to find its sub-projects: %w", err))
	}
	return m, nil
}

// subProjectOf is the sub-project raw's project names, or the one supported
// when it names none, as flai components serve --project picks it.
func subProjectOf(p channel.Project, raw json.RawMessage) (manifest.Manifest, manifest.Project, *channel.Error) {
	in, e := decode[struct {
		Project string `json:"project"`
	}](raw)
	if e != nil {
		return manifest.Manifest{}, manifest.Project{}, e
	}
	m, e := componentsManifest(p)
	if e != nil {
		return manifest.Manifest{}, manifest.Project{}, e
	}
	sub, err := components.PickSubProject(m.Projects, strings.TrimSpace(in.Project))
	if err != nil {
		return manifest.Manifest{}, manifest.Project{}, bad("%s", err)
	}
	return m, sub, nil
}

// hostRelative is msg with the project's root left out of the paths it
// names, which are then relative to the repository, as the dashboard sees
// them.
func hostRelative(root, msg string) string {
	return strings.ReplaceAll(msg, root+string(filepath.Separator), "")
}

// componentRefusal is the error a dev server's start answers with: a
// sub-project without vite is the operator's to fix, a rule; anything else
// failed.
func componentRefusal(root string, err error) *channel.Error {
	code := channel.CodeInternal
	if errors.Is(err, components.ErrViteMissing) {
		code = Rule
	}
	return &channel.Error{Code: code, Message: hostRelative(root, err.Error())}
}

// written is v as a write's answer, with no warnings.
func written(v any) (any, *channel.Error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, failed(err)
	}
	return Written{Data: data, Warnings: []string{}}, nil
}

// componentSpecs are the writes under the components host action. They run
// no flai command: build checks the params and do asks flai serve's
// supervisor.
func componentSpecs() map[string]spec {
	check := func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		_, _, e := subProjectOf(p, raw)
		return nil, "", e
	}
	return map[string]spec{
		// components.start: the sub-project's own Vite dev server with flai's
		// harness, started or, when it runs, answered as it is.
		"components.start": {action: ActionComponents, describe: describeComponentStart, build: check,
			do: func(ctx context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				m, sub, e := subProjectOf(p, raw)
				if e != nil {
					return nil, e
				}
				if host.Components == nil {
					return nil, noComponentServers()
				}
				s, err := host.Components.Start(ctx, p.Root, m.Name, sub)
				if err != nil {
					return nil, componentRefusal(p.Root, err)
				}
				s.LastError = hostRelative(p.Root, s.LastError)
				return written(s)
			}},
		// components.stop: the sub-project's dev server stopped, when one runs.
		"components.stop": {action: ActionComponents, describe: describeComponentStop, build: check,
			do: func(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				_, sub, e := subProjectOf(p, raw)
				if e != nil {
					return nil, e
				}
				if host.Components == nil {
					return nil, noComponentServers()
				}
				s, err := host.Components.Stop(p.Root, sub.Name)
				if err != nil {
					return nil, failed(errors.New(hostRelative(p.Root, err.Error())))
				}
				return written(s)
			}},
	}
}

// describeComponentStart says for the journal which sub-project's dev
// server was started, or answered as it ran, and where.
func describeComponentStart(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var s ComponentServer
	_ = json.Unmarshal(w.Data, &s)
	if s.AlreadyRunning {
		return "done", fmt.Sprintf("%s's component dev server was running at %s (pid %d)", s.Project, s.URL, s.PID)
	}
	return "done", fmt.Sprintf("started %s's component dev server at %s (pid %d)", s.Project, s.URL, s.PID)
}

// describeComponentStop says for the journal which sub-project's dev server
// was stopped, or that none ran.
func describeComponentStop(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var s ComponentStopped
	_ = json.Unmarshal(w.Data, &s)
	if !s.Stopped {
		return "done", fmt.Sprintf("%s's component dev server was not running", s.Project)
	}
	return "done", fmt.Sprintf("stopped %s's component dev server (pid %d)", s.Project, s.PID)
}

// componentMethods are the component library's reads.
func componentMethods(host Host) map[string]channel.Method {
	return map[string]channel.Method{
		// components.status: each supported sub-project's dev server, or the
		// named one's, whether it runs, where, since when, and its last
		// error; asking marks each used, so that the page that shows one
		// keeps it from being stopped as idle.
		"components.status": func(_ context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[struct {
				Project string `json:"project"`
			}](raw)
			if e != nil {
				return nil, e
			}
			m, e := componentsManifest(p)
			if e != nil {
				return nil, e
			}
			var subs []manifest.Project
			if name := strings.TrimSpace(in.Project); name != "" {
				sub, err := components.PickSubProject(m.Projects, name)
				if err != nil {
					return nil, bad("%s", err)
				}
				subs = []manifest.Project{sub}
			} else {
				for _, sp := range m.Projects {
					if manifest.SupportsComponents(sp.Kind) {
						subs = append(subs, sp)
					}
				}
			}
			out := ComponentsStatus{Enabled: host.enabled(ActionComponents, p.Root), Enable: EnableCommand(ActionComponents), Servers: make([]ComponentServer, 0, len(subs))}
			if host.Components != nil {
				out.IdleSeconds = int(host.Components.Idle() / time.Second)
			}
			for _, sub := range subs {
				s := ComponentServer{Project: sub.Name}
				if host.Components != nil {
					s = host.Components.Status(p.Root, sub.Name)
					s.LastError = hostRelative(p.Root, s.LastError)
				}
				out.Servers = append(out.Servers, s)
			}
			return out, nil
		},

		// components.list: the catalogue, as flai components list --json
		// prints it: every sub-project in the manifest's order, or the one
		// named, a Svelte one with its components and their props.
		"components.list": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[struct {
				Project string `json:"project"`
			}](raw)
			if e != nil {
				return nil, e
			}
			m, e := componentsManifest(p)
			if e != nil {
				return nil, e
			}
			projects := m.Projects
			if name := strings.TrimSpace(in.Project); name != "" {
				i := slices.IndexFunc(projects, func(sp manifest.Project) bool { return sp.Name == name })
				if i < 0 {
					return nil, bad("%s", components.UnknownSubProject(name, projects))
				}
				projects = projects[i : i+1]
			}
			listed := make([]ComponentsProject, 0, len(projects))
			for _, sp := range projects {
				lp := ComponentsProject{Name: sp.Name, Path: sp.Path, Kind: sp.Kind, Components: []components.Component{}}
				if s, ok := sp.ComponentSettings(); ok {
					comps, err := components.List(ctx, nodeRunner, p.Root, sp.Path, s.Include, s.Exclude)
					if errors.Is(err, components.ErrNodeMissing) || errors.Is(err, components.ErrSvelteMissing) {
						return nil, &channel.Error{Code: Rule, Message: hostRelative(p.Root, err.Error())}
					}
					if err != nil {
						return nil, failed(errors.New(hostRelative(p.Root, err.Error())))
					}
					lp.Supported, lp.Setup, lp.Components = true, s.Setup, comps
				}
				listed = append(listed, lp)
			}
			return map[string]any{"projects": listed}, nil
		},
	}
}
