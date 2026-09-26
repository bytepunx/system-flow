package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai serve project manages the projects flai serve serves (S-0118): the
// registry flai dashboard writes, changed and read on its own, with the
// reason a project is not showing on the dashboard.

// errNoKey is a manifest with no key: flai serve names a project by it.
var errNoKey = errors.New("its " + manifest.File + " has no key (flai check says how to add one)")

// registerProject registers repo with flai serve for the dashboard s, making
// the agent credential when there is none. flai dashboard, flai import, and
// flai serve project add all register through it.
func (a *app) registerProject(repo *workitem.Repo, s dashboardSettings) (serve.Entry, error) {
	dir := a.serveDir()
	e := a.serveEntry(repo, s)
	if e.Key == "" {
		return e, errNoKey
	}
	if err := ensureAgentKey(string(dir)); err != nil {
		return e, fmt.Errorf("agent credential: %w", err)
	}
	return e, dir.Register(e)
}

func newServeProjectCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "project",
		Short: "The projects flai serve serves on the dashboard: add, remove, and list them, and why one is not showing",
		Long: `flai serve serves each project registered with it on the dashboard its entry
names. flai dashboard in a project registers it, and so does an import from the
board; these commands do it on their own, without starting or stopping the
dashboard, and say why a project is not showing.

Not to be confused with flai serve import, which names the folders whose git
repositories the board offers to import.`,
		Example: `  flai serve project add            # the project in this folder
  flai serve project add ~/git/blog
  flai serve project list
  flai serve project remove blog`,
	}
	c.AddCommand(
		&cobra.Command{
			Use:   "add [folder]",
			Short: "Serve the project in a folder (default: this one) on the dashboard",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.serveProjectAdd(args)
			},
		},
		&cobra.Command{
			Use:   "remove <key|folder>",
			Short: "Stop serving a project; none of its files is touched, and the dashboard keeps running",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.serveProjectRemove(args[0])
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "Every project served and how it is, the repositories offered for import, and the projects not served",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				return a.serveProjectList()
			},
		},
	)
	return c
}

func (a *app) serveProjectAdd(args []string) error {
	folder, err := a.workingDir()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		folder = resolveIn(folder, args[0])
	}
	if info, err := os.Stat(folder); err != nil || !info.IsDir() {
		return fmt.Errorf("%s is not a folder", folder)
	}
	if _, err := os.Stat(filepath.Join(folder, manifest.File)); err != nil {
		return fmt.Errorf("%s has no %s: flai import there makes it a project", folder, manifest.File)
	}
	repo, err := workitem.Open(folder)
	if err != nil {
		return err
	}
	s, err := a.dashboardSettings(repo, "", "", 0, "")
	if err != nil {
		return err
	}
	e, err := a.registerProject(repo, s)
	if err != nil {
		return fmt.Errorf("%s not served: %w", folder, err)
	}
	st, started, herr := a.ensureHost()
	if a.jsonOut {
		out := map[string]any{"project": e, "dashboard": s.url(), "host_started": started}
		if herr != nil {
			out["host_error"] = herr.Error()
		} else {
			out["host_pid"] = st.PID
		}
		return a.printJSON(out)
	}
	fmt.Fprintf(a.out, "%s (%s) is served on the dashboard at %s\n", e.Key, e.Root, s.url())
	switch {
	case herr != nil:
		fmt.Fprintf(a.out, "  host flai: registered, but flai host did not start: %s\n    start it with: flai host start\n", herr)
	case started:
		fmt.Fprintf(a.out, "  host flai: flai host started (pid %d); it runs flai serve, which connects the project to the dashboard\n", st.PID)
	default:
		fmt.Fprintf(a.out, "  host flai: flai host is running (pid %d); its flai serve connects the project within a second\n", st.PID)
	}
	fmt.Fprintln(a.out, "  no dashboard running? flai dashboard starts it; flai serve project list shows how each project is")
	return nil
}

// resolveIn is path made absolute against cwd, the folder flai runs in.
func resolveIn(cwd, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(cwd, path)
}

func (a *app) serveProjectRemove(arg string) error {
	dir := a.serveDir()
	all, err := dir.Projects()
	if err != nil {
		return err
	}
	cwd, err := a.workingDir()
	if err != nil {
		return err
	}
	folder := resolveIn(cwd, arg)
	// a folder inside the project, or a worktree of it, names the project too
	if _, err := os.Stat(folder); err == nil {
		if repo, err := workitem.Open(folder); err == nil {
			if folder = repo.Root; repo.MainRoot != "" {
				folder = repo.MainRoot
			}
		}
	}
	// a key first: a folder named like a key must not stand for the project it is in
	var found *serve.Entry
	for _, same := range []func(serve.Entry) bool{
		func(e serve.Entry) bool { return e.Key == arg },
		func(e serve.Entry) bool { return e.Root == folder },
	} {
		for i := range all {
			if found == nil && same(all[i]) {
				found = &all[i]
			}
		}
	}
	if found == nil {
		if st, alive := dir.ReadStatus(time.Now()); alive {
			for _, e := range st.FolderProjects {
				if e.Key == arg || e.Root == folder {
					return fmt.Errorf("%s is not registered: flai serve serves it because it was started in %s, which it is below", e.Key, st.Folder)
				}
			}
		}
		return fmt.Errorf("no project %s is served; flai serve project list shows those that are", arg)
	}
	if err := dir.Unregister(found.Root); err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"removed": found})
	}
	fmt.Fprintf(a.out, "%s (%s) is no longer served, and none of its files was touched\n", found.Key, found.Root)
	if _, alive := dir.ReadStatus(time.Now()); alive {
		fmt.Fprintln(a.out, "  flai serve drops it within a second, and the dashboard's switcher with it")
	}
	fmt.Fprintln(a.out, "  the dashboard keeps running for the other projects; flai serve project add serves it again")
	return nil
}

// servedProject is one project flai serve serves, and how it is.
type servedProject struct {
	serve.Entry
	// Folder is true for a project served because flai serve was started in
	// a folder above it (S-0102), and not registered.
	Folder bool `json:"folder,omitempty"`
	// State is connected, connecting, not-connected, unavailable, or
	// not-running (flai serve is not).
	State     string `json:"state"`
	Since     string `json:"since,omitempty"`
	LastError string `json:"last_error,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// unservedProject is a system-flow project under an import folder that
// flai serve does not serve.
type unservedProject struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name"`
	Root string `json:"root"`
	// Next is what serves it.
	Next string `json:"next"`
}

type projectListing struct {
	Running    bool              `json:"running"`
	Served     []servedProject   `json:"served"`
	Candidates []serve.Candidate `json:"candidates"`
	Unserved   []unservedProject `json:"unserved"`
	Roots      []string          `json:"import_roots"`
}

func (a *app) listServedProjects() (projectListing, error) {
	dir := a.serveDir()
	registered, err := dir.Projects()
	if err != nil {
		return projectListing{}, err
	}
	out := projectListing{Served: []servedProject{}, Candidates: []serve.Candidate{}, Unserved: []unservedProject{}, Roots: []string{}}
	st, alive := dir.ReadStatus(time.Now())
	out.Running = alive
	describe := func(e serve.Entry, folder bool) servedProject {
		p := servedProject{Entry: e, Folder: folder}
		switch c, ok := st.Connections[e.Root]; {
		case !folder && e.Unavailable() != "":
			p.State, p.Reason = "unavailable", e.Unavailable()
		case !alive:
			p.State = "not-running"
		case ok && c.Connected:
			p.State, p.Since = "connected", c.Since
		case ok && c.LastError != "":
			p.State, p.LastError = "not-connected", c.LastError
		default:
			p.State = "connecting"
		}
		return p
	}
	served, taken := map[string]bool{}, map[string]bool{}
	for _, e := range registered {
		out.Served = append(out.Served, describe(e, false))
		served[e.Root], taken[e.Key] = true, true
	}
	if alive {
		for _, e := range st.FolderProjects {
			out.Served = append(out.Served, describe(e, true))
			served[e.Root], taken[e.Key] = true, true
		}
	}
	sort.Slice(out.Served, func(i, j int) bool { return out.Served[i].Key < out.Served[j].Key })

	cfg, _, err := a.loadConfig()
	if err != nil {
		return out, err
	}
	out.Roots = append(out.Roots, cfg.ImportRoots...)
	if found := serve.FindCandidates(cfg.ImportRoots, served, taken); found != nil {
		out.Candidates = found
	}
	seen := map[string]bool{}
	for _, r := range cfg.ImportRoots {
		for _, root := range workitem.FindProjects(r) {
			if served[root] || seen[root] {
				continue
			}
			seen[root] = true
			u := unservedProject{Name: filepath.Base(root), Root: root, Next: "flai serve project add " + root}
			if m, err := manifest.Load(filepath.Join(root, manifest.File)); err != nil {
				u.Next = "fix its " + manifest.File + ": " + err.Error()
			} else if u.Key, u.Name = m.Key, m.Name; m.Key == "" {
				u.Next = "give its " + manifest.File + " a key (flai check says how), then flai serve project add " + root
			}
			out.Unserved = append(out.Unserved, u)
		}
	}
	return out, nil
}

func (a *app) serveProjectList() error {
	l, err := a.listServedProjects()
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(l)
	}
	if !l.Running {
		fmt.Fprintln(a.out, "flai serve is not running, so no project is on the dashboard (flai serve start)")
	}
	if len(l.Served) == 0 {
		fmt.Fprintln(a.out, "no project served; flai serve project add in a project serves it")
	} else {
		fmt.Fprintln(a.out, "served:")
	}
	for _, p := range l.Served {
		state := p.State
		switch p.State {
		case "connected":
			state = "connected since " + p.Since
		case "not-connected":
			state = "not connected: " + p.LastError
		case "unavailable":
			state = "not served: " + p.Reason
		case "not-running":
			state = "not connected: flai serve is not running"
		}
		how := ""
		if p.Folder {
			how = " (below the folder flai serve was started in)"
		}
		fmt.Fprintf(a.out, "  %s  %s  %s%s\n    %s\n", p.Key, p.URL, state, how, p.Root)
	}
	if len(l.Candidates) > 0 {
		fmt.Fprintln(a.out, "offered for import on the board (flai serve import):")
		for _, c := range l.Candidates {
			fmt.Fprintf(a.out, "  %s  %s\n", c.Name, c.Root)
		}
	}
	if len(l.Unserved) > 0 {
		fmt.Fprintln(a.out, "projects under the import folders that are not served:")
		for _, u := range l.Unserved {
			fmt.Fprintf(a.out, "  %s  %s\n    %s\n", u.Name, u.Root, u.Next)
		}
	}
	return nil
}
