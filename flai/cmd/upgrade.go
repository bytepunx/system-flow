package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/upgrade"
)

func newUpgradeCmd(a *app) *cobra.Command {
	var tplRepo, ref string
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
unresolved. A dirty git tree is refused unless --force.`,
		Example: `  flai upgrade --dry-run
  flai upgrade
  flai upgrade --keep-all           # scripts and CI: never overwrite edits
  flai upgrade --template ./template --force
  flai upgrade --relock             # hand-assembled project: record the current files at this version`,
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
			opt := template.Options{Vars: manifestVars(mf), Layout: mf.Layout, Source: src}
			if relock {
				l, err := upgrade.Relock(repo.Root, m, src, opt, a.now())
				if err != nil {
					return err
				}
				if err := lock.Save(repo.Root, l); err != nil {
					return err
				}
				if err := writeTemplateFields(filepath.Join(repo.Root, manifest.File), src, sourceChanged, m.Version, a.now()); err != nil {
					return err
				}
				fmt.Fprintf(a.out, "relocked %d files at template %s\n", len(l.Files), m.Version)
				return nil
			}
			if m.Version == mf.Template.Version && !force {
				fmt.Fprintf(a.out, "already at template %s (use --force to re-apply)\n", m.Version)
				return nil
			}
			if !force && !dryRun {
				if st, err := a.runner.Run(repo.Root, "git", "status", "--porcelain"); err == nil && strings.TrimSpace(st) != "" {
					return fmt.Errorf("working tree has uncommitted changes; commit or stash them so the upgrade is one reviewable diff, or pass --force")
				}
			}
			lk, err := lock.Load(repo.Root)
			if err != nil {
				return err
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
			if err := lock.Save(repo.Root, upgrade.NewLock(plan, src, m.Version, a.now())); err != nil {
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
	return c
}

// askConflict prompts for one conflict, with a diff on request.
func (a *app) askConflict(root, tplDir string, m template.Manifest, opt template.Options, path string) (string, error) {
	for {
		choice := "keep"
		if err := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title("Conflict: "+path+" was changed in this project and in the template").
			Options(huh.NewOption("keep the project's version", "keep"), huh.NewOption("replace with the template's version", "replace"), huh.NewOption("show diff", "diff")).Value(&choice))).Run(); err != nil {
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

// manifestVars rebuilds template variables from the manifest so a re-render
// produces the same text the project was created with.
func manifestVars(mf manifest.Manifest) map[string]any {
	return map[string]any{
		"project_name": mf.Name, "project_key": mf.Key, "description": mf.Description,
		"owner": mf.Owner, "repo_url": mf.Repo,
	}
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
