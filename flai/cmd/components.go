package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// componentsRunner runs node for flai components list; tests replace it.
var componentsRunner components.Runner = components.System{}

// newComponentsCmd groups the commands on the component library: the
// catalogue of a project's Svelte components and the manifest key that says
// where they are (S-0437, ADR-0162).
func newComponentsCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "components",
		Short: "List the component library: each Svelte sub-project's components and their props",
		Long: `The component library is the components of the project's Svelte
sub-projects, those whose kind in system-flow.yaml is sveltekit or svelte
(ADR-0162). A sub-project's components are its .svelte files that match its
components.include globs, src/lib/**/*.svelte by default, and none of its
components.exclude globs; its components.setup names a module the harness
imports before any component, such as the app's stylesheet.

flai components list prints the catalogue, each component with the props it
takes, read by the sub-project's own svelte/compiler; flai components config
writes a sub-project's components key, which flai manifest set does not
reach; flai components serve runs a sub-project's own Vite dev server with
flai's harness, which renders one component by its path (ADR-0163).`,
		Example: `  flai components list
  flai components list --project flaiover --json
  flai components config flaiover --setup src/routes/layout.css
  flai components serve --project flaiover`,
	}
	c.AddCommand(newComponentsListCmd(a), newComponentsConfigCmd(a), newComponentsServeCmd(a))
	return c
}

// componentsProject is one sub-project in flai components list's answer.
type componentsProject struct {
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Kind       string                 `json:"kind"`
	Supported  bool                   `json:"supported"`
	Setup      string                 `json:"setup,omitempty"`
	Components []components.Component `json:"components"`
}

func newComponentsListCmd(a *app) *cobra.Command {
	var project string
	c := &cobra.Command{
		Use:   "list",
		Short: "Print each Svelte sub-project's components and the props each takes",
		Long: `Print the component library: every sub-project of system-flow.yaml, in its
order, or only the one --project names. A sub-project of kind sveltekit or
svelte is listed with its components, each a .svelte file its components key
finds, by its path relative to the sub-project, and under it each prop: its
name, its kind (string, number, boolean, union with its literals, array,
object, function, snippet, or unknown), its type as the source writes it, its
default after =, and required when it has neither a default nor an optional
type. A component whose props cannot be read, such as one that does not
parse, is listed with the reason. A sub-project of any other kind is listed
as unsupported, with its kind.

Props are read by a script flai embeds, run with node in the sub-project's
folder on the sub-project's own svelte/compiler; no component is run. So
node must be on PATH and svelte installed in each Svelte sub-project listed:
when either is missing, list fails, exit status 1, saying what to install.

--json prints {"projects": [{"name", "path", "kind", "supported", "setup",
"components": [{"project", "path", "name", "props": [{"name", "kind",
"type", "default", "required", "literals"}], "error"}]}]}, setup only when
the sub-project names one, and type, default, literals, and error only when
they have a value.`,
		Example: `  flai components list
  flai components list --project flaiover
  flai components list --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			projects := repo.Manifest.Projects
			if project != "" {
				i := slices.IndexFunc(projects, func(p manifest.Project) bool { return p.Name == project })
				if i < 0 {
					return components.UnknownSubProject(project, projects)
				}
				projects = projects[i : i+1]
			}
			listed := make([]componentsProject, 0, len(projects))
			for _, p := range projects {
				lp := componentsProject{Name: p.Name, Path: p.Path, Kind: p.Kind, Components: []components.Component{}}
				if s, ok := p.ComponentSettings(); ok {
					comps, err := components.List(cmd.Context(), componentsRunner, repo.Root, p.Path, s.Include, s.Exclude)
					if err != nil {
						return err
					}
					lp.Supported, lp.Setup, lp.Components = true, s.Setup, comps
				}
				listed = append(listed, lp)
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"projects": listed})
			}
			return printComponents(a.out, listed)
		},
	}
	c.Flags().StringVar(&project, "project", "", "list only the sub-project of this name, one of the manifest's projects")
	return c
}

// printComponents writes flai components list's text: a block for each
// sub-project, its components under it, and each one's props aligned under
// that.
func printComponents(w io.Writer, listed []componentsProject) error {
	if len(listed) == 0 {
		_, err := fmt.Fprintln(w, "no sub-projects in system-flow.yaml: add one of kind sveltekit or svelte under projects to list its components")
		return err
	}
	var b bytes.Buffer
	for i, p := range listed {
		if i > 0 {
			b.WriteString("\n")
		}
		if !p.Supported {
			fmt.Fprintf(&b, "%s (%s): unsupported kind %s\n", p.Name, p.Path, p.Kind)
			continue
		}
		head := p.Kind
		if p.Setup != "" {
			head += ", setup " + p.Setup
		}
		fmt.Fprintf(&b, "%s (%s, %s): %s\n", p.Name, p.Path, head, plural(len(p.Components), "component"))
		for _, c := range p.Components {
			fmt.Fprintf(&b, "  %s\n", c.Path)
			switch {
			case c.Error != "":
				fmt.Fprintf(&b, "    error: %s\n", oneLine(c.Error))
			case len(c.Props) == 0:
				b.WriteString("    no props\n")
			default:
				tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
				for _, pr := range c.Props {
					fmt.Fprintf(tw, "    %s\t%s\t%s\t%s\n", pr.Name, propKind(pr), oneLine(pr.Type), propExtras(pr))
				}
				if err := tw.Flush(); err != nil {
					return err
				}
			}
		}
	}
	// tabwriter pads the type column of a prop with no default and not
	// required; drop the trailing spaces that leaves.
	for _, line := range strings.SplitAfter(b.String(), "\n") {
		if line == "" {
			continue
		}
		if _, err := io.WriteString(w, strings.TrimRight(line, " \n")+"\n"); err != nil {
			return err
		}
	}
	return nil
}

// propKind is a prop's kind, a union's followed by its literals.
func propKind(p components.Prop) string {
	if p.Kind == components.KindUnion && len(p.Literals) > 0 {
		return string(p.Kind) + " " + strings.Join(p.Literals, "|")
	}
	return string(p.Kind)
}

// propExtras is a prop's default after =, and required when it is.
func propExtras(p components.Prop) string {
	var parts []string
	if p.Default != "" {
		parts = append(parts, "= "+oneLine(p.Default))
	}
	if p.Required {
		parts = append(parts, "required")
	}
	return strings.Join(parts, " ")
}

// oneLine is s with each run of white space, line breaks among it, one
// space, so that a type or default written over lines keeps to one.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func newComponentsConfigCmd(a *app) *cobra.Command {
	var include, exclude []string
	var setup string
	var unset bool
	c := &cobra.Command{
		Use:   "config <sub-project>",
		Short: "Write a Svelte sub-project's components key: its include and exclude globs and its setup",
		Long: `Write the components key of the sub-project the argument names in
system-flow.yaml, keeping every other key, the order, and the comments.

Each flag given replaces its field of the key as it is; a field whose flag is
not given keeps its value, so --setup alone keeps the globs. --include and
--exclude are repeatable, and the globs given replace the field's whole list.
A glob is relative to the sub-project's folder, with segments separated by /,
each a path.Match pattern or ** alone, which stands for any number of
folders. --setup is a module of the sub-project, relative to its folder, that
the harness imports before any component; --setup "" removes it. --unset
removes the key, so that the defaults apply: include src/lib/**/*.svelte,
nothing excluded, no setup. It is given alone.

The manifest is checked as it would be after the change, as flai check checks
it. A sub-project whose kind is not sveltekit or svelte, a glob that does not
parse, or a setup that names no file under the sub-project's folder is
refused with the field and the reason, and nothing is written; the exit
status is 1, as for a workflow rule's refusal. With --json a refusal is
{"refused": [{"field": ..., "reason": ...}]} on standard output, and a change
{"project": ..., "components": {"include": [...], "exclude": [...], "setup":
...}}, components null after --unset.`,
		Example: `  flai components config flaiover --setup src/routes/layout.css
  flai components config ui --include 'src/lib/**/*.svelte' --include 'src/widgets/*.svelte'
  flai components config ui --exclude 'src/lib/internal/**' --json
  flai components config ui --unset`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			f := cmd.Flags()
			given := f.Changed("include") || f.Changed("exclude") || f.Changed("setup")
			switch {
			case unset && given:
				return fmt.Errorf("--unset removes sub-project %q's components key whole, so it takes no other flag: give --unset alone, or --include, --exclude, or --setup without it", name)
			case !unset && !given:
				return fmt.Errorf("nothing to write for sub-project %q's components: give --include, --exclude, or --setup to set a field, or --unset to remove the key, such as flai components config %s --setup src/app.css", name, name)
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			var next *manifest.Components
			if !unset {
				var c manifest.Components
				if i := slices.IndexFunc(repo.Manifest.Projects, func(p manifest.Project) bool { return p.Name == name }); i >= 0 && repo.Manifest.Projects[i].Components != nil {
					c = *repo.Manifest.Projects[i].Components
				}
				if f.Changed("include") {
					c.Include = include
				}
				if f.Changed("exclude") {
					c.Exclude = exclude
				}
				if f.Changed("setup") {
					c.Setup = setup
				}
				next = &c
			}
			if err := manifest.SetComponents(filepath.Join(repo.Root, manifest.File), name, next); err != nil {
				var refused *manifest.RefusedError
				if errors.As(err, &refused) {
					return a.refusedSettings(refused)
				}
				return err
			}
			return a.printComponentsConfig(name, next)
		},
	}
	f := c.Flags()
	f.StringArrayVar(&include, "include", nil, "a glob of the sub-project's component files, relative to its folder (repeatable; replaces the list)")
	f.StringArrayVar(&exclude, "exclude", nil, "a glob of files include finds that are not components (repeatable; replaces the list)")
	f.StringVar(&setup, "setup", "", `a module of the sub-project the harness imports before any component, such as src/app.css; "" removes it`)
	f.BoolVar(&unset, "unset", false, "remove the components key, so that the defaults apply")
	return c
}

// printComponentsConfig says what flai components config wrote: the
// sub-project's components key, or that it was removed when c is nil.
func (a *app) printComponentsConfig(name string, c *manifest.Components) error {
	if a.jsonOut {
		return a.printJSON(map[string]any{"project": name, "components": c})
	}
	if c == nil {
		_, err := fmt.Fprintf(a.out, "unset sub-project %s's components, so the defaults apply: include %s, nothing excluded, no setup\n", name, manifest.DefaultComponentInclude)
		return err
	}
	include := strings.Join(c.Include, ", ")
	if include == "" {
		include = manifest.DefaultComponentInclude + " (the default)"
	}
	exclude := strings.Join(c.Exclude, ", ")
	if exclude == "" {
		exclude = "none"
	}
	setup := c.Setup
	if setup == "" {
		setup = "none"
	}
	_, err := fmt.Fprintf(a.out, "set sub-project %s's components\n  include  %s\n  exclude  %s\n  setup    %s\n", name, include, exclude, setup)
	return err
}

func newComponentsServeCmd(a *app) *cobra.Command {
	var project string
	c := &cobra.Command{
		Use:   "serve",
		Short: "Run a Svelte sub-project's own Vite dev server with flai's component harness, in the foreground",
		Long: `Run the Vite dev server of the sub-project --project names with flai's
component harness, in the foreground, until Ctrl-C or SIGTERM stops it. The
sub-project is one of kind sveltekit or svelte in system-flow.yaml; --project
may be left out when the manifest has one such sub-project only.

flai writes a Vite config under .flai-cache/components/<project>/<sub-project>/
that loads the sub-project's own Vite config and merges flai's harness plugin
into it, then runs the sub-project's own vite, from its node_modules, in its
folder, on 127.0.0.1 and a free port (ADR-0163). Once the harness answers, it
prints the harness's URL, the port's /__flai/components/ on 127.0.0.1, and
under it the sub-project and how to stop it. Open the URL with ?c=<path>, a
component's path in the sub-project as flai components list prints it, to
render that component after the components key's setup module; a saved
change reaches the page as a hot update. Vite's own output goes to standard
error.

This runs the project's code and its dependencies on this machine, as the
sub-project's own dev server does, reachable from this machine only. The
dashboard has flai serve start the same dev server only once the operator
turns the components host action on (flai serve enable components).

A sub-project without vite installed is refused with what to install, and
a dev server that exits on its own ends the command, exit status 1, with its
last lines of output. --json prints {"project", "url", "port", "pid"} once
the harness answers.`,
		Example: `  flai components serve
  flai components serve --project flaiover
  flai components serve --project flaiover --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			sub, err := components.PickSubProject(repo.Manifest.Projects, project)
			if err != nil {
				return err
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			return a.serveComponents(ctx, repo.Root, repo.Manifest.Name, sub)
		},
	}
	c.Flags().StringVar(&project, "project", "", "serve the sub-project of this name, of kind sveltekit or svelte; needed when the manifest has more than one")
	return c
}

// serveComponents starts sub's dev server, prints where its harness answers,
// and runs until ctx ends, then stops it, or until it exits on its own, then
// returns why.
func (a *app) serveComponents(ctx context.Context, root, project string, sub manifest.Project) error {
	s, err := components.Serve(ctx, components.ServeOptions{Root: root, Project: project, Sub: sub, Output: a.errOut})
	if err != nil {
		return err
	}
	log := a.logger().With("component", "components", "project", sub.Name, "pid", s.PID(), "port", s.Port())
	log.Info("component dev server started")
	if err := a.printComponentsServe(sub, s); err != nil {
		return errors.Join(err, s.Stop())
	}
	select {
	case <-ctx.Done():
		if err := s.Stop(); err != nil {
			return err
		}
		log.Info("component dev server stopped")
		return nil
	case <-s.Done():
		return s.Err()
	}
}

// printComponentsServe says where sub's dev server s answers: its harness's
// URL on a line of its own, then the sub-project and how to stop it.
func (a *app) printComponentsServe(sub manifest.Project, s *components.DevServer) error {
	if a.jsonOut {
		return a.printJSON(map[string]any{"project": sub.Name, "url": s.URL(), "port": s.Port(), "pid": s.PID()})
	}
	_, err := fmt.Fprintf(a.out, "%s\n  sub-project %s (%s), Vite pid %d; Ctrl-C stops it\n", s.URL(), sub.Name, sub.Path, s.PID())
	return err
}
