package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
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
reach.`,
		Example: `  flai components list
  flai components list --project flaiover --json
  flai components config flaiover --setup src/routes/layout.css`,
	}
	c.AddCommand(newComponentsListCmd(a), newComponentsConfigCmd(a))
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
					return unknownSubProject(project, projects)
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

// unknownSubProject is the error for a sub-project name the manifest does
// not have, listing the names it does.
func unknownSubProject(name string, projects []manifest.Project) error {
	if len(projects) == 0 {
		return fmt.Errorf("system-flow.yaml has no sub-project %q, and names none under projects; add the sub-project there first", name)
	}
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return fmt.Errorf("system-flow.yaml has no sub-project %q; name one of %s", name, strings.Join(names, ", "))
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
