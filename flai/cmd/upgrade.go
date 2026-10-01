package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/prompt"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/upgrade"
)

func newUpgradeCmd(a *app) *cobra.Command {
	var tplRepo, ref string
	var vars []string
	var dryRun, force, keepAll, replaceAll, relock bool
	c := &cobra.Command{
		Use:   "upgrade",
		Short: "Bring this project to the latest template version",
		Long: `Fetch the template (config, or --template and --ref), compare its version with
system-flow.yaml, and apply the difference per ADR-0015: add new files,
merge marker files (CLAUDE.md, conventions) above the marker, replace files
the project has not changed since they were applied, and report the rest as
conflicts. In a terminal each conflict offers keep, replace, or a diff;
otherwise --keep-all or --replace-all is required and nothing changes without
one. The manifest and lock are updated only when no conflict is left
unresolved. A dirty git tree is refused unless --force.

Each template variable is rendered with, in order: --var; the manifest's own
field for project_name, project_key, description, owner, and repo_url; the
value system-flow.lock.yaml recorded; the template's default, named in the
output when taken. A required variable with none of these is named and
nothing is changed. --var re-applies the version the project is at.`,
		Example: `  flai upgrade --dry-run
  flai upgrade
  flai upgrade --keep-all           # scripts and CI: never overwrite edits
  flai upgrade --template ./template --force
  flai upgrade --relock             # hand-assembled project: record the current files at this version
  flai upgrade --var team=billing   # a variable the template added, or a new value for one`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			mf := repo.Manifest
			sourceChanged := tplRepo != ""
			if tplRepo == "" {
				tplRepo, ref = mf.Template.Repo, mf.Template.Ref
			}
			src, m, err := a.resolveTemplate(tplRepo, ref, false)
			if err != nil {
				return err
			}
			given, err := givenVars(vars)
			if err != nil {
				return err
			}
			lk, err := lock.Load(repo.Root)
			if err != nil {
				return err
			}
			// --var re-applies the version the project is at, to change a value.
			if !relock && m.Version == mf.Template.Version && !force && len(given) == 0 {
				fmt.Fprintf(a.out, "already at template %s (use --force to re-apply)\n", m.Version)
				return nil
			}
			vals, defaulted, err := upgradeVars(m, mf, lk, given)
			if err != nil {
				return err
			}
			if !a.jsonOut {
				for _, name := range defaulted {
					fmt.Fprintf(a.out, "%s has no recorded value: rendering with its default %q (flai upgrade --var %s=value sets it)\n", name, vals[name], name)
				}
			}
			opt := template.Options{Vars: vals, Layout: mf.Layout, Source: src}
			if relock {
				l, err := upgrade.Relock(repo.Root, m, src, opt, a.now())
				if err != nil {
					return err
				}
				l.Vars = lockVars(opt.Vars)
				if err := lock.Save(repo.Root, l); err != nil {
					return err
				}
				if err := writeTemplateFields(filepath.Join(repo.Root, manifest.File), src, sourceChanged, m.Version, a.now()); err != nil {
					return err
				}
				fmt.Fprintf(a.out, "relocked %d files at template %s\n", len(l.Files), m.Version)
				return nil
			}
			if !force && !dryRun {
				if st, err := a.runner.Run(repo.Root, "git", "status", "--porcelain"); err == nil && strings.TrimSpace(st) != "" {
					return fmt.Errorf("working tree has uncommitted changes; commit or stash them so the upgrade is one reviewable diff, or pass --force")
				}
			}
			plan, err := upgrade.Compute(repo.Root, m, src, opt, lk, mf.Template.Version)
			if err != nil {
				return err
			}
			if a.jsonOut && dryRun {
				return a.printJSON(plan)
			}
			fmt.Fprintln(a.out, upgrade.Describe(plan))
			if plan.NoLock {
				fmt.Fprintln(a.out, "no system-flow.lock.yaml: every differing file counts as a conflict (flai upgrade --relock records the current state)")
			}
			for _, ch := range plan.Changes {
				if ch.Class != upgrade.Same {
					fmt.Fprintf(a.out, "  %-9s %s\n", ch.Class, ch.Path)
				}
			}
			if dryRun {
				fmt.Fprintln(a.out, "dry run: nothing changed")
				return nil
			}
			policy := upgrade.Policy{}
			conflicts := plan.Conflicts()
			interactive := !a.yes && a.isTerminal()
			switch {
			case len(conflicts) == 0:
			case keepAll:
				for _, p := range conflicts {
					policy[p] = "keep"
				}
			case replaceAll:
				for _, p := range conflicts {
					policy[p] = "replace"
				}
			case interactive:
				for _, p := range conflicts {
					choice, err := a.askConflict(repo.Root, src.Dir, m, opt, p)
					if err != nil {
						return err
					}
					policy[p] = choice
				}
			default:
				return &exitError{code: 1, msg: fmt.Sprintf("%d conflict(s) need a decision: %s\nre-run with --keep-all, --replace-all, or in a terminal", len(conflicts), strings.Join(conflicts, ", "))}
			}
			written, err := upgrade.Apply(repo.Root, plan, policy)
			if err != nil {
				return err
			}
			nl := upgrade.NewLock(plan, src, m.Version, a.now())
			nl.Vars = lockVars(opt.Vars)
			if err := lock.Save(repo.Root, nl); err != nil {
				return err
			}
			if err := writeTemplateFields(filepath.Join(repo.Root, manifest.File), src, sourceChanged, m.Version, a.now()); err != nil {
				return err
			}
			kept := 0
			for _, v := range policy {
				if v == "keep" {
					kept++
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"plan": plan, "written": written, "kept": kept})
			}
			fmt.Fprintf(a.out, "upgraded to template %s: %d files written, %d conflicts kept\n", m.Version, len(written), kept)
			if kept > 0 {
				fmt.Fprintln(a.out, "kept files stay divergent from the template and will be conflicts again next upgrade")
			}
			return nil
		},
	}
	f := c.Flags()
	f.StringVar(&tplRepo, "template", "", "template git URL or local directory (default: the manifest's template.repo)")
	f.StringVar(&ref, "ref", "", "template branch, tag, or commit (default: the manifest's template.ref)")
	f.BoolVar(&dryRun, "dry-run", false, "print the plan and change nothing")
	f.BoolVar(&force, "force", false, "re-apply the same version and allow a dirty tree")
	f.BoolVar(&keepAll, "keep-all", false, "keep every conflicting project file")
	f.BoolVar(&replaceAll, "replace-all", false, "replace every conflicting project file with the template's")
	f.BoolVar(&relock, "relock", false, "record the current files at the template's version without changing them")
	f.StringArrayVar(&vars, "var", nil, "set a template variable the manifest does not hold, name=value (repeatable); recorded in the lock")
	return c
}

// askConflict prompts for one conflict, with a diff on request.
func (a *app) askConflict(root, tplDir string, m template.Manifest, opt template.Options, path string) (string, error) {
	for {
		choice, err := a.prompts().Select("Conflict: "+path+" was changed in this project and in the template",
			[]prompt.Option{{Label: "keep the project's version", Value: "keep"}, {Label: "replace with the template's version", Value: "replace"}, {Label: "show diff", Value: "diff"}}, "keep")
		if err != nil {
			return "", err
		}
		if choice != "diff" {
			return choice, nil
		}
		tmp, err := os.CreateTemp("", "flai-upgrade-*")
		if err != nil {
			return "", err
		}
		var content []byte
		o := opt
		o.Collect = func(rel string, c []byte, _ os.FileMode) {
			if rel == path {
				content = c
			}
		}
		if _, err := template.Render(m, tplDir, root, o); err != nil {
			return "", err
		}
		_ = os.WriteFile(tmp.Name(), content, 0o644)
		out, _ := a.runner.Run(root, "git", "diff", "--no-index", "--", filepath.Join(root, filepath.FromSlash(path)), tmp.Name())
		_ = os.Remove(tmp.Name())
		fmt.Fprintln(a.out, out)
	}
}

// heldVar is a template variable system-flow.yaml holds: the manifest key
// it is read from, and its value there.
type heldVar struct{ key, value string }

// manifestVars are the template variables the manifest holds, by variable.
func manifestVars(mf manifest.Manifest) map[string]heldVar {
	return map[string]heldVar{
		"project_name": {"name", mf.Name}, "project_key": {"key", mf.Key}, "description": {"description", mf.Description},
		"owner": {"owner", mf.Owner}, repoURLVar: {"repo", mf.Repo},
	}
}

// upgradeVars resolves every template variable for a re-render, in order:
// --var, the manifest's own field, the value the lock recorded, the
// template's default; a required variable recorded empty takes its default.
// defaulted names the variables that took their default. A required variable
// left empty, an unknown --var, and a --var for a value the manifest holds
// are refused before anything is rendered (I-0040).
func upgradeVars(m template.Manifest, mf manifest.Manifest, lk *lock.Lock, given map[string]string) (vars map[string]any, defaulted []string, err error) {
	held := manifestVars(mf)
	known := map[string]bool{}
	for _, v := range m.Variables {
		known[v.Name] = true
	}
	for k := range given {
		if !known[k] {
			return nil, nil, fmt.Errorf("unknown variable %q; template defines %s", k, strings.Join(varNames(m), ", "))
		}
		if h, ok := held[k]; ok {
			return nil, nil, fmt.Errorf("%s is read from %s in %s; change it there, not with --var", k, h.key, manifest.File)
		}
	}
	vars = map[string]any{}
	var missing []string
	for _, v := range m.Variables {
		val, ok := given[v.Name]
		if h, isHeld := held[v.Name]; !ok && isHeld {
			val, ok = h.value, true
		}
		if !ok && lk != nil {
			val, ok = lk.Vars[v.Name]
			ok = ok && (val != "" || !v.Required)
		}
		if !ok {
			def, err := template.EvalDefault(v, template.Data(m, template.Options{Vars: vars}))
			if err != nil {
				return nil, nil, err
			}
			val = def
			defaulted = append(defaulted, v.Name)
		}
		if v.Required && val == "" {
			if h, isHeld := held[v.Name]; isHeld {
				missing = append(missing, fmt.Sprintf("%s (%s in %s)", v.Name, h.key, manifest.File))
			} else {
				missing = append(missing, v.Name)
			}
		}
		vars[v.Name] = val
	}
	if len(missing) > 0 {
		return nil, nil, fmt.Errorf("required variables have no value: %s; give each with --var name=value, or in %s where it says, and run flai upgrade again: nothing was changed", strings.Join(missing, ", "), manifest.File)
	}
	return vars, defaulted, nil
}

// writeTemplateFields updates template.version and applied in the manifest
// in place, and repo and ref only when the operator pointed at a different
// source, touching nothing else.
func writeTemplateFields(path string, src template.Source, sourceChanged bool, version string, now time.Time) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	repo := src.Repo
	if src.Local {
		repo = src.Dir
	}
	lines := strings.Split(string(data), "\n")
	inTemplate := false
	for i, l := range lines {
		switch {
		case strings.HasPrefix(l, "template:"):
			inTemplate = true
		case inTemplate && strings.HasPrefix(l, "  repo:") && sourceChanged:
			lines[i] = fmt.Sprintf("  repo: %q", repo)
		case inTemplate && strings.HasPrefix(l, "  ref:") && sourceChanged:
			lines[i] = fmt.Sprintf("  ref: %q", src.Ref)
		case inTemplate && strings.HasPrefix(l, "  version:"):
			lines[i] = fmt.Sprintf("  version: %q", version)
		case inTemplate && strings.HasPrefix(l, "  applied:"):
			lines[i] = "  applied: " + now.UTC().Format("2006-01-02T15:04:05Z")
		case inTemplate && !strings.HasPrefix(l, "  "):
			inTemplate = false
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}
