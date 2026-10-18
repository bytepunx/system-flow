package hostapi

// The component library's methods (S-0437, S-0438, ADR-0163):
// components.list, the catalogue flai components list prints, and
// components.status, how each sub-project's dev server stands, are reads;
// components.start and components.stop are writes under the components host
// action, answered in flai serve's own process, which supervises the dev
// servers, through Host.Components. Nothing here starts or ends a process.
// Each takes an optional story (S-0442): the catalogue and the dev servers
// of that story's worktree, while it is in progress or in review, in place
// of the main checkout's.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ComponentServer is how one sub-project's component dev server stands, as
// components.status and components.start answer it.
type ComponentServer struct {
	// Project is the sub-project's name in the manifest.
	Project string `json:"project"`
	// Story is the story whose worktree the dev server serves; absent for
	// the main checkout's.
	Story string `json:"story,omitempty"`
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
	// Story is the story whose worktree's dev server it was; absent for the
	// main checkout's.
	Story string `json:"story,omitempty"`
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
	Name       string            `json:"name"`
	Path       string            `json:"path"`
	Kind       string            `json:"kind"`
	Supported  bool              `json:"supported"`
	Setup      string            `json:"setup,omitempty"`
	Components []ListedComponent `json:"components"`
}

// ListedComponent is one component in components.list's answer: the
// component and its props, its examples, and where its examples file is
// (S-0443, ADR-0164).
type ListedComponent struct {
	components.Component
	// Examples are its examples file's, in the file's order; empty when it
	// has no file. A file that does not fit the format still gives the
	// examples that parsed, an example's props {} when they are not an
	// object, and flai check reports what does not fit.
	Examples []components.Example `json:"examples"`
	// ExamplesPath is its examples file's path relative to the checkout,
	// given whether or not the file exists, so that the page knows where a
	// first example goes.
	ExamplesPath string `json:"examples_path"`
}

// ListComponents is the catalogue of projects, sub-projects of the checkout
// at root whose manifest is m, in their order: a Svelte one with its
// components, their props read with node through r, and their examples,
// any other kind unsupported. Its error is the one components.List gives,
// which wraps components.ErrNodeMissing or components.ErrSvelteMissing when
// a tool is missing, or that of an examples file that cannot be read.
func ListComponents(ctx context.Context, r components.Runner, root string, m manifest.Manifest, projects []manifest.Project) ([]ComponentsProject, error) {
	fsys := os.DirFS(root)
	listed := make([]ComponentsProject, 0, len(projects))
	for _, sp := range projects {
		lp := ComponentsProject{Name: sp.Name, Path: sp.Path, Kind: sp.Kind, Components: []ListedComponent{}}
		if s, ok := sp.ComponentSettings(); ok {
			comps, err := components.List(ctx, r, root, sp.Path, s.Include, s.Exclude)
			if err != nil {
				return nil, err
			}
			lp.Supported, lp.Setup = true, s.Setup
			for _, c := range comps {
				file := components.ExamplesPath(m, sp.Name, c.Path)
				examples, _, err := components.Examples(fsys, file)
				if err != nil {
					return nil, err
				}
				for i := range examples {
					if examples[i].Props == nil {
						examples[i].Props = map[string]json.RawMessage{}
					}
				}
				lp.Components = append(lp.Components, ListedComponent{Component: c, Examples: examples, ExamplesPath: file})
			}
		}
		listed = append(listed, lp)
	}
	return listed, nil
}

// ComponentServers are the sub-projects' component dev servers flai serve
// supervises, one per project, sub-project, and checkout: the main
// checkout's, story empty, or a story's worktree's (ADR-0163, S-0442).
type ComponentServers interface {
	// Start starts the dev server of sub, a sub-project of the project at
	// root whose manifest names it project, and answers it once it answers:
	// in the main checkout when story is empty, and otherwise in dir, the
	// story's worktree. When one runs, or is starting, it answers that one
	// instead. A start beyond components.max_servers is refused with an
	// error that wraps ErrComponentServersFull.
	Start(ctx context.Context, root, project string, sub manifest.Project, story, dir string) (ComponentServer, error)
	// Stop stops the dev server of the sub-project named sub in story's
	// worktree, or the main checkout when story is empty, when one runs.
	Stop(root, sub, story string) (ComponentStopped, error)
	// Status is how the dev server of the sub-project named sub in story's
	// worktree, or the main checkout when story is empty, stands, and marks
	// it used, so that it is not stopped as idle.
	Status(root, sub, story string) ComponentServer
	// Idle is how long a dev server may go without a Status before it is
	// stopped: components.idle.
	Idle() time.Duration
}

// ErrComponentServersFull is wrapped by the error a start beyond
// components.max_servers is refused with; it is the operator's to settle by
// stopping one, a rule.
var ErrComponentServersFull = errors.New("components.max_servers reached")

// nodeRunner runs node for components.list; tests replace it.
var nodeRunner components.Runner = components.System{}

// noComponentServers refuses a start or a stop asked of a flai that runs no
// dev servers.
func noComponentServers() *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no component dev servers: flai serve runs them, so ask through the dashboard flai serve answers, or run flai components serve in a shell on the host"}
}

// componentsParams are what every components method takes: the
// sub-project, as sub_project, and the story whose worktree it is in. The
// sub-project is not project: project in a call's params is the channel's
// project key alone, which the dashboard stamps on every call, so a method
// that read it would get the key, not the sub-project (I-0163).
type componentsParams struct {
	SubProject string `json:"sub_project"`
	Story      string `json:"story"`
}

// checkout is where a components method reads and serves: the project's
// main checkout, or a story's worktree.
type checkout struct {
	// story is the story's ID, empty for the main checkout.
	story string
	// dir is the checkout's root, absolute.
	dir string
	// m is the checkout's manifest.
	m manifest.Manifest
}

// rel is msg with the checkout's root and then the project's left out of
// the paths it names, which are then relative to the checkout, as the
// dashboard sees them.
func (c checkout) rel(root, msg string) string {
	return hostRelative(root, hostRelative(c.dir, msg))
}

// checkoutOf is the checkout of story, which must be in progress or in
// review and have a worktree, or the project's main checkout when story is
// empty, with its manifest. A story no item has is NotFound.
func checkoutOf(p channel.Project, story string) (checkout, *channel.Error) {
	c := checkout{dir: p.Root}
	if story = strings.TrimSpace(story); story != "" {
		if e := needID(story); e != nil {
			return checkout{}, e
		}
		repo, err := workitem.Open(p.Root)
		if err != nil {
			return checkout{}, failed(err)
		}
		id, dir, err := ComponentsWorktree(repo, story)
		var refused *storyRefused
		switch {
		case errors.Is(err, workitem.ErrNotFound):
			return checkout{}, &channel.Error{Code: NotFound, Message: fmt.Sprintf("%s is not a work item here; name a story in progress or in review with a worktree, or leave story out for the main checkout", story)}
		case errors.As(err, &refused):
			return checkout{}, bad("%s", hostRelative(p.Root, err.Error()))
		case err != nil:
			return checkout{}, failed(errors.New(hostRelative(p.Root, err.Error())))
		}
		c.story, c.dir = id, dir
	}
	m, err := manifest.Load(filepath.Join(c.dir, manifest.File))
	if err != nil {
		return checkout{}, failed(fmt.Errorf("cannot read the project's manifest to find its sub-projects: %s", c.rel(p.Root, err.Error())))
	}
	c.m = m
	return c, nil
}

// storyRefused is ComponentsWorktree's refusal of a story whose components
// flai does not serve.
type storyRefused struct{ msg string }

func (e *storyRefused) Error() string { return e.msg }

// ComponentsWorktree is the canonical ID and the worktree of the story id
// names in repo, whose components flai lists and serves while it is in
// progress or in review (S-0442). An ID no item has wraps
// workitem.ErrNotFound; an item that is not a story, a story in any other
// state, and a story without a worktree are refused, as ADR-0145 refuses
// one, saying how to open it.
func ComponentsWorktree(repo *workitem.Repo, id string) (string, string, error) {
	it, err := repo.Get(id)
	if err != nil {
		return "", "", fmt.Errorf("%w; name a story in progress or in review with a worktree, or leave the story out for the main checkout", err)
	}
	switch {
	case it.Type != workitem.Story:
		return "", "", &storyRefused{fmt.Sprintf("%s is a %s; name a story (S-nnnn) in progress or in review, whose worktree holds the components it changes", it.ID, it.Type)}
	case it.Archived || it.Status != workitem.InProgress && it.Status != workitem.Review:
		return "", "", &storyRefused{fmt.Sprintf("%s is %s: flai serves a story's components only while it is in progress or in review; name one that is, or leave the story out for the main checkout", it.ID, it.Status)}
	}
	wt := repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", "", &storyRefused{fmt.Sprintf("%s has no worktree; open it with flai stream open %s on the host, or leave the story out for the main checkout", it.ID, it.ID)}
	}
	return it.ID, wt, nil
}

// subProjectOf is the checkout raw's story names and the sub-project its
// sub_project names there, or the one supported when it names none, as flai
// components serve --project picks it.
func subProjectOf(p channel.Project, raw json.RawMessage) (checkout, manifest.Project, *channel.Error) {
	in, e := decode[componentsParams](raw)
	if e != nil {
		return checkout{}, manifest.Project{}, e
	}
	c, e := checkoutOf(p, in.Story)
	if e != nil {
		return checkout{}, manifest.Project{}, e
	}
	sub, err := components.PickSubProject(c.m.Projects, strings.TrimSpace(in.SubProject))
	if err != nil {
		return checkout{}, manifest.Project{}, bad("%s", err)
	}
	return c, sub, nil
}

// hostRelative is msg with the project's root left out of the paths it
// names, which are then relative to the repository, as the dashboard sees
// them.
func hostRelative(root, msg string) string {
	return strings.ReplaceAll(msg, root+string(filepath.Separator), "")
}

// componentRefusal is the error a dev server's start in c answers with: a
// sub-project without vite or without a lock file to install it by, and a
// start beyond components.max_servers, are the operator's to settle, a
// rule; anything else failed.
func componentRefusal(root string, c checkout, err error) *channel.Error {
	code := channel.CodeInternal
	if errors.Is(err, components.ErrViteMissing) || errors.Is(err, components.ErrNoLockFile) || errors.Is(err, ErrComponentServersFull) {
		code = Rule
	}
	return &channel.Error{Code: code, Message: c.rel(root, err.Error())}
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
		// harness, in the main checkout or the story's worktree, started or,
		// when it runs, answered as it is.
		"components.start": {action: ActionComponents, describe: describeComponentStart, build: check,
			do: func(ctx context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				c, sub, e := subProjectOf(p, raw)
				if e != nil {
					return nil, e
				}
				if host.Components == nil {
					return nil, noComponentServers()
				}
				s, err := host.Components.Start(ctx, p.Root, c.m.Name, sub, c.story, c.dir)
				if err != nil {
					return nil, componentRefusal(p.Root, c, err)
				}
				s.LastError = c.rel(p.Root, s.LastError)
				return written(s)
			}},
		// components.stop: the sub-project's dev server, the main checkout's or
		// the story's worktree's, stopped, when one runs.
		"components.stop": {action: ActionComponents, describe: describeComponentStop, build: check,
			do: func(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				c, sub, e := subProjectOf(p, raw)
				if e != nil {
					return nil, e
				}
				if host.Components == nil {
					return nil, noComponentServers()
				}
				s, err := host.Components.Stop(p.Root, sub.Name, c.story)
				if err != nil {
					return nil, failed(errors.New(c.rel(p.Root, err.Error())))
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
		return "done", fmt.Sprintf("%s component dev server was running at %s (pid %d)", whose(s.Story, s.Project), s.URL, s.PID)
	}
	return "done", fmt.Sprintf("started %s component dev server at %s (pid %d)", whose(s.Story, s.Project), s.URL, s.PID)
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
		return "done", fmt.Sprintf("%s component dev server was not running", whose(s.Story, s.Project))
	}
	return "done", fmt.Sprintf("stopped %s component dev server (pid %d)", whose(s.Story, s.Project), s.PID)
}

// whose names the dev server of sub in the journal: ui's for the main
// checkout's, S-0441's ui for a story's worktree's.
func whose(story, sub string) string {
	if story == "" {
		return sub + "'s"
	}
	return story + "'s " + sub
}

// componentMethods are the component library's reads.
func componentMethods(host Host) map[string]channel.Method {
	return map[string]channel.Method{
		// components.status: each supported sub-project's dev server, or the
		// one sub_project names, in the main checkout or the story's
		// worktree, whether it runs, where, since when, and its last error;
		// asking marks each used, so that the page that shows one keeps it
		// from being stopped as idle.
		"components.status": func(_ context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[componentsParams](raw)
			if e != nil {
				return nil, e
			}
			c, e := checkoutOf(p, in.Story)
			if e != nil {
				return nil, e
			}
			m := c.m
			var subs []manifest.Project
			if name := strings.TrimSpace(in.SubProject); name != "" {
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
				s := ComponentServer{Project: sub.Name, Story: c.story}
				if host.Components != nil {
					s = host.Components.Status(p.Root, sub.Name, c.story)
					s.LastError = c.rel(p.Root, s.LastError)
				}
				out.Servers = append(out.Servers, s)
			}
			return out, nil
		},

		// components.list: the catalogue, as flai components list --json
		// prints it: every sub-project in the manifest's order, or the one
		// sub_project names, a Svelte one with its components, their props,
		// and their examples, read from the main checkout or the story's
		// worktree.
		"components.list": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[componentsParams](raw)
			if e != nil {
				return nil, e
			}
			c, e := checkoutOf(p, in.Story)
			if e != nil {
				return nil, e
			}
			projects := c.m.Projects
			if name := strings.TrimSpace(in.SubProject); name != "" {
				i := slices.IndexFunc(projects, func(sp manifest.Project) bool { return sp.Name == name })
				if i < 0 {
					return nil, bad("%s", components.UnknownSubProject(name, projects))
				}
				projects = projects[i : i+1]
			}
			listed, err := ListComponents(ctx, nodeRunner, c.dir, c.m, projects)
			switch {
			case errors.Is(err, components.ErrNodeMissing) || errors.Is(err, components.ErrSvelteMissing):
				return nil, &channel.Error{Code: Rule, Message: c.rel(p.Root, err.Error())}
			case err != nil:
				return nil, failed(errors.New(c.rel(p.Root, err.Error())))
			}
			return map[string]any{"projects": listed}, nil
		},
	}
}
