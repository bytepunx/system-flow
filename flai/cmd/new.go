package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

type newOptions struct {
	templateRepo string
	ref          string
	vars         []string
	layout       []string
	defaults     bool
	force        bool
	noGit        bool
}

func newNewCmd(a *app) *cobra.Command {
	var o newOptions
	c := &cobra.Command{
		Use:   "new <dir>",
		Short: "Create a new monorepo from the template",
		Long: `Render the system-flow template into a new directory.

The template comes from --template (a git URL or a local directory), or from
template.repo in the config file. Variables are prompted for when running in a
terminal; supply them with --var or accept defaults with --defaults for
non-interactive use. Existing files are never overwritten unless --force.`,
		Example: `  flai new my-project
  flai new my-project --defaults --var description="Billing platform"
  flai new my-project --template ./template --no-git
  flai new my-project --template git@github.com:me/system-flow-template.git --ref my-branch
  flai new my-project --layout design=architecture`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runNew(args[0], o)
		},
	}
	f := c.Flags()
	f.StringVar(&o.templateRepo, "template", "", "template git URL or local directory (default: config template.repo)")
	f.StringVar(&o.ref, "ref", "", "template branch, tag, or commit (default: config template.ref)")
	f.StringArrayVar(&o.vars, "var", nil, "set a template variable, name=value (repeatable)")
	f.StringArrayVar(&o.layout, "layout", nil, "rename a layout folder, key=name, e.g. design=architecture (repeatable)")
	f.BoolVar(&o.defaults, "defaults", false, "use defaults for every variable not given with --var, never prompt")
	f.BoolVar(&o.force, "force", false, "overwrite existing files")
	f.BoolVar(&o.noGit, "no-git", false, "do not run git init in the new directory")
	return c
}

func (a *app) runNew(dir string, o newOptions) error {
	dest, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	src, m, err := a.resolveTemplate(o.templateRepo, o.ref, false)
	if err != nil {
		return err
	}

	vars, err := a.collectVars(m, o, filepath.Base(dest))
	if err != nil {
		return err
	}
	layout, err := parsePairs(o.layout, "layout")
	if err != nil {
		return err
	}
	for k := range layout {
		if _, ok := m.Layout[k]; !ok {
			return fmt.Errorf("unknown layout key %q; template defines %s", k, strings.Join(m.LayoutKeys(), ", "))
		}
	}

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	res, err := template.Render(m, src.Dir, dest, template.Options{
		Vars: vars, Layout: layout, Source: src, Force: o.force,
	})
	if err != nil {
		return err
	}
	if _, err := manifest.Load(filepath.Join(dest, manifest.File)); err != nil {
		return fmt.Errorf("rendered project has an invalid manifest: %w", err)
	}

	gitInit := false
	if !o.noGit {
		if _, err := os.Stat(filepath.Join(dest, ".git")); err != nil {
			if _, err := a.runner.LookPath("git"); err == nil {
				if _, err := a.runner.Run(dest, "git", "init", "--quiet", "--initial-branch=main"); err != nil {
					return err
				}
				gitInit = true
			}
		}
	}

	if a.jsonOut {
		return a.printJSON(map[string]any{
			"dir": dest, "template": src.Repo, "ref": src.Ref, "version": m.Version,
			"written": res.Written, "skipped": res.Skipped, "git_init": gitInit, "vars": vars,
		})
	}
	fmt.Fprintf(a.out, "Created %s from template %s (%s)\n", dest, src.Repo, m.Version)
	fmt.Fprintf(a.out, "  %d files written", len(res.Written))
	if len(res.Skipped) > 0 {
		fmt.Fprintf(a.out, ", %d existing files kept (use --force to overwrite)", len(res.Skipped))
	}
	fmt.Fprintln(a.out)
	if gitInit {
		fmt.Fprintln(a.out, "  git repository initialised on main, nothing committed")
	}
	fmt.Fprintf(a.out, "\nNext:\n  cd %s\n  flai epic new \"First deliverable\"\n  flai dashboard\n", dir)
	return nil
}

// resolveTemplate picks the source from flags or config, ensures it is
// available locally, and loads its manifest.
func (a *app) resolveTemplate(repo, ref string, refresh bool) (template.Source, template.Manifest, error) {
	cacheDir := config.Default().CacheDir
	if repo == "" {
		cfg, _, err := a.loadConfig()
		if err != nil {
			return template.Source{}, template.Manifest{}, err
		}
		repo, cacheDir = cfg.Template.Repo, cfg.CacheDir
		if ref == "" {
			ref = cfg.Template.Ref
		}
	} else if cfg, _, err := a.loadConfig(); err == nil {
		cacheDir = cfg.CacheDir
	}
	src, err := template.Resolve(repo, ref, cacheDir)
	if err != nil {
		return template.Source{}, template.Manifest{}, err
	}
	if !src.Cached() || refresh {
		if !a.jsonOut {
			fmt.Fprintf(a.errOut, "flai: fetching template %s", src.Repo)
			if src.Ref != "" {
				fmt.Fprintf(a.errOut, " @ %s", src.Ref)
			}
			fmt.Fprintln(a.errOut)
		}
	}
	if err := src.Ensure(a.runner, refresh); err != nil {
		return template.Source{}, template.Manifest{}, err
	}
	m, err := template.LoadManifest(src.Dir)
	if err != nil {
		return template.Source{}, template.Manifest{}, fmt.Errorf("%s: %w", src.Repo, err)
	}
	return src, m, nil
}

// collectVars merges --var values, evaluated defaults, and prompts.
func (a *app) collectVars(m template.Manifest, o newOptions, dirName string) (map[string]any, error) {
	given, err := parsePairs(o.vars, "var")
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, v := range m.Variables {
		known[v.Name] = true
	}
	for k := range given {
		if !known[k] {
			return nil, fmt.Errorf("unknown variable %q; template defines %s", k, strings.Join(varNames(m), ", "))
		}
	}
	interactive := !o.defaults && a.isTerminal()
	vars := map[string]any{}
	var missing []string
	for _, v := range m.Variables {
		if val, ok := given[v.Name]; ok {
			vars[v.Name] = val
			continue
		}
		def, err := template.EvalDefault(v, template.Data(m, template.Options{Vars: vars}))
		if err != nil {
			return nil, err
		}
		if v.Name == "project_name" && def == "" {
			def = dirName
		}
		val := def
		if interactive {
			val, err = a.prompt(v, def)
			if err != nil {
				return nil, err
			}
		}
		if v.Required && val == "" {
			missing = append(missing, v.Name)
		}
		vars[v.Name] = val
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("required variables not set: %s (use --var name=value)", strings.Join(missing, ", "))
	}
	return vars, nil
}

func (a *app) prompt(v template.Variable, def string) (string, error) {
	val := def
	title := v.Prompt
	if title == "" {
		title = v.Name
	}
	in := huh.NewInput().Title(title).Value(&val)
	if v.Required {
		in = in.Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("%s is required", v.Name)
			}
			return nil
		})
	}
	if err := huh.NewForm(huh.NewGroup(in)).Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}

func (a *app) isTerminal() bool {
	if a.stdinIsTerminal != nil {
		return *a.stdinIsTerminal
	}
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func parsePairs(items []string, flag string) (map[string]string, error) {
	out := map[string]string{}
	for _, it := range items {
		k, v, ok := strings.Cut(it, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("--%s expects name=value, got %q", flag, it)
		}
		out[k] = v
	}
	return out, nil
}

func varNames(m template.Manifest) []string {
	names := make([]string, 0, len(m.Variables))
	for _, v := range m.Variables {
		names = append(names, v.Name)
	}
	sort.Strings(names)
	return names
}

var _ execx.Runner = execx.System{}
