package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai serve project manages the projects flai serve serves (S-0121): the
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

A project below a folder named for import, or below the folder flai serve
was started in, is served without being registered. Removing one puts its
folder on flai serve's list of removed projects, which it does not serve
from below a folder; add takes it off the list, and it is served from there
again.

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
	if repo.Manifest.Key == "" {
		return fmt.Errorf("%s not served: %w", folder, errNoKey)
	}
	// one the operator removed from below a folder flai serve serves from is
	// served from there again, as it was, rather than registered (S-0123)
	root := mainRootOf(repo)
	restored, err := a.serveDir().RestoreRoot(root)
	if err != nil {
		return err
	}
	below := ""
	if restored {
		below = a.servedBelow(root)
	}
	s, err := a.dashboardSettings(repo, "", "", 0, "")
	if err != nil {
		return err
	}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: root}
	if below == "" {
		if e, err = a.registerProject(repo, s); err != nil {
			return fmt.Errorf("%s not served: %w", folder, err)
		}
	}
	st, started, herr := a.ensureHost()
	if a.jsonOut {
		out := map[string]any{"project": e, "host_started": started}
		if below != "" {
			out["served_below"] = below
		} else {
			out["dashboard"] = s.url()
		}
		if restored {
			out["restored"] = true
		}
		if herr != nil {
			out["host_error"] = herr.Error()
		} else {
			out["host_pid"] = st.PID
		}
		return a.printJSON(out)
	}
	if below != "" {
		fmt.Fprintf(a.out, "%s (%s) is off the list of removed projects, so flai serve serves it again from below %s, as it did before it was removed\n", e.Key, e.Root, below)
	} else {
		if restored {
			fmt.Fprintf(a.out, "%s (%s) is off the list of removed projects\n", e.Key, e.Root)
		}
		fmt.Fprintf(a.out, "%s (%s) is registered with flai serve, for the dashboard at %s\n", e.Key, e.Root, s.url())
	}
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
			folder = mainRootOf(repo)
		}
	}
	// a key first: a folder named like a key must not stand for the project it is in
	byKey := func(e serve.Entry) bool { return arg != "" && e.Key == arg }
	byRoot := func(e serve.Entry) bool { return e.Root == folder }
	below := a.projectsBelow()
	var found *serve.Entry
	registered := false
	for _, same := range []func(serve.Entry) bool{byKey, byRoot} {
		for i := range all {
			if found == nil && same(all[i]) {
				found, registered = &all[i], true
			}
		}
		// one served from below a folder, not registered (S-0120, S-0102), is
		// put on the list of removed projects instead (S-0123)
		for _, f := range below {
			if found == nil && same(f.Entry) {
				found = &f.Entry
			}
		}
	}
	if found == nil {
		return fmt.Errorf("no project %s is served; flai serve project list shows those that are", arg)
	}
	if registered {
		if err := dir.Unregister(found.Root); err != nil {
			return err
		}
	}
	// below a folder flai serve serves from, it would go on serving it from
	// there: it is kept off with the list
	from := a.servedBelow(found.Root)
	if from != "" {
		if err := dir.RemoveRoot(found.Root); err != nil {
			return err
		}
	}
	if a.jsonOut {
		out := map[string]any{"removed": found, "registered": registered}
		if from != "" {
			out["listed_below"] = from
		}
		return a.printJSON(out)
	}
	name := found.Key
	if name == "" {
		name = found.Name
	}
	if registered {
		fmt.Fprintf(a.out, "%s (%s) is no longer registered, and none of its files was touched\n", name, found.Root)
	} else {
		fmt.Fprintf(a.out, "%s (%s) is no longer served, and none of its files was touched\n", name, found.Root)
	}
	if from != "" {
		fmt.Fprintf(a.out, "  it is below %s, which flai serve serves from, so it is on the list of removed projects, and not served from there\n", from)
	}
	if _, alive := dir.ReadStatus(time.Now()); alive {
		fmt.Fprintln(a.out, "  flai serve drops it within a second, and the dashboard's switcher with it")
	}
	fmt.Fprintf(a.out, "  the dashboard keeps running for the other projects; flai serve project add %s serves it again\n", found.Root)
	return nil
}

// projectsBelow are the projects below the folder flai serve was started in
// and the folders named for import, served or not: from flai serve's state
// when it runs, and those it would find below the import folders either way.
func (a *app) projectsBelow() []serve.Found {
	var out []serve.Found
	if st, alive := a.serveDir().ReadStatus(time.Now()); alive {
		for _, e := range st.FolderProjects {
			out = append(out, serve.Found{Entry: e})
		}
		for _, e := range st.ImportProjects {
			out = append(out, serve.Found{Entry: e, Imported: true})
		}
		out = append(out, st.Unserved...)
	}
	return append(out, serve.FindBelow("", a.importRoots())...)
}

// servedBelow is the folder flai serve serves root from, if it is below one:
// a folder named for import, or the folder a running flai serve was started
// in (S-0102). Empty otherwise.
func (a *app) servedBelow(root string) string {
	if r := a.importRootOf(root); r != "" {
		return r
	}
	if st, alive := a.serveDir().ReadStatus(time.Now()); alive && st.Folder != "" && within(st.Folder, root) {
		return st.Folder
	}
	return ""
}

// within says whether root is folder or below it.
func within(folder, root string) bool {
	rel, err := filepath.Rel(folder, root)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// importRootOf is the folder named for import that root is served below, as
// flai serve finds them (S-0120), or empty.
func (a *app) importRootOf(root string) string {
	roots := a.importRoots()
	for _, f := range serve.FindBelow("", roots) {
		if f.Root != root {
			continue
		}
		for _, r := range roots {
			if within(r, root) {
				return r
			}
		}
	}
	return ""
}

// servedProject is one project flai serve serves, and how it is.
type servedProject struct {
	serve.Entry
	// From is registry for a registered project; folder for one served
	// because flai serve was started in a folder above it (S-0102); import
	// for one below a folder named for import (S-0120).
	From string `json:"from"`
	// State is connected, connecting, not-connected, unavailable, or
	// not-running (flai serve is not).
	State     string `json:"state"`
	Since     string `json:"since,omitempty"`
	LastError string `json:"last_error,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type projectListing struct {
	Running    bool              `json:"running"`
	Served     []servedProject   `json:"served"`
	Candidates []serve.Candidate `json:"candidates"`
	// Unserved are the projects below the folders named for import, or the
	// folder flai serve was started in, that it does not serve, and why (S-0120).
	Unserved []serve.Found `json:"unserved"`
	Roots    []string      `json:"import_roots"`
}

func (a *app) listServedProjects() (projectListing, error) {
	ss, err := a.readServeStatus()
	if err != nil {
		return projectListing{}, err
	}
	out := projectListing{Running: ss.Running, Served: []servedProject{}, Candidates: []serve.Candidate{}, Unserved: []serve.Found{}, Roots: []string{}}
	var st serve.Status
	if ss.Running {
		st = *ss.Status
	}
	describe := func(e serve.Entry, from string) servedProject {
		p := servedProject{Entry: e, From: from}
		switch c, ok := st.Connections[e.Root]; {
		case ss.Unavailable[e.Root] != "":
			p.State, p.Reason = "unavailable", ss.Unavailable[e.Root]
		case !ss.Running:
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
	add := func(entries []serve.Entry, from string) {
		for _, e := range entries {
			out.Served = append(out.Served, describe(e, from))
			served[e.Root], taken[e.Key] = true, true
		}
	}
	add(ss.Projects, "registry")
	add(st.FolderProjects, "folder")
	add(st.ImportProjects, "import")
	sort.Slice(out.Served, func(i, j int) bool { return out.Served[i].Key < out.Served[j].Key })

	cfg, _, err := a.loadConfig()
	if err != nil {
		return out, err
	}
	out.Roots = append(out.Roots, cfg.ImportRoots...)
	if found := serve.FindCandidates(cfg.ImportRoots, served, taken); found != nil {
		out.Candidates = found
	}
	if _, unserved := a.belowImportRoots(ss); unserved != nil {
		out.Unserved = unserved
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
		switch p.From {
		case "folder":
			how = " (below the folder flai serve was started in)"
		case "import":
			how = " (below a folder named for import)"
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
		fmt.Fprintln(a.out, "not served, below the folders named for import or the folder flai serve was started in:")
		for _, u := range l.Unserved {
			fmt.Fprintf(a.out, "  %s  %s\n    %s\n", u.Name, u.Root, u.Reason)
		}
	}
	return nil
}
