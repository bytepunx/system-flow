package cmd

import (
	"errors"
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
		Long: `Fetch the template and apply the difference between the version this
project is at and the version chosen, per ADR-0015 and ADR-0103.

The version: --ref names it (1.0.60 matches the tag v1.0.60) and wins over
template.ref in system-flow.yaml; it is recorded there and in the lock. With
no --ref, a git template is taken to its newest version tag when
template.ref follows releases: empty, the default branch such as main, or a
version tag. Another branch or a commit is used as given, fetched again; a
template with no version tags, a local directory, --template, and --relock
use the ref as given. When system-flow.yaml names a template.ref or
template.version the lock did not record, because it was edited, and that
release is not the newest, a terminal is asked which to apply: the version
system-flow.yaml names, the newest, or neither. Without a terminal, or with
--yes, nothing changes and the command exits 1 naming the flai upgrade --ref
to run for each. The version the project is at is the lock's, or the
manifest's when there is no lock. The ref and version applied are written to
system-flow.yaml and the lock. --dry-run says which version it would apply,
or that it would ask.

The difference: add new files, merge marker files (CLAUDE.md, conventions)
above the marker, replace files the project has not changed since they were
applied, and report the rest as conflicts. In a terminal each conflict offers
keep, replace, or a diff; otherwise --keep-all or --replace-all is required
and nothing changes without one. The manifest and lock are updated only when
no conflict is left unresolved. A dirty git tree is refused unless --force.

Each template variable is rendered with, in order: --var; the manifest's own
field for project_name, project_key, description, owner, and repo_url; the
value system-flow.lock.yaml recorded; the template's default, named in the
output when taken. A required variable with none of these is named and
nothing is changed. --var re-applies the version the project is at.`,
		Example: `  flai upgrade --dry-run
  flai upgrade
  flai upgrade --ref 1.0.60         # a chosen release, recorded in system-flow.yaml
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
			// An unreadable lock stops only an upgrade that has work to do.
			lk, lkErr := lock.Load(repo.Root)
			if lkErr != nil {
				lk = nil
			}
			at := projectVersion(mf.Template, lk)
			target, ok, err := a.pickTarget(mf.Template, lk, tplRepo, ref, relock, dryRun, at)
			if err != nil || !ok {
				return err
			}
			// The manifest records the repo when --template names another,
			// and the ref when the one applied is not the one it names.
			fields := templateFields{repo: tplRepo != "", ref: tplRepo != "" || target.Ref != mf.Template.Ref}
			if tplRepo == "" {
				tplRepo = mf.Template.Repo
			}
			src, m, err := a.resolveTemplate(tplRepo, target.Ref, false)
			if err != nil {
				return err
			}
			target.Version = m.Version
			if !a.jsonOut {
				fmt.Fprintf(a.out, "template %s at %s: %s\n", targetLabel(src), m.Version, target.Reason)
			}
			given, err := givenVars(vars)
			if err != nil {
				return err
			}
			// --var re-applies the version the project is at, to change a value.
			if !relock && m.Version == at && !force && len(given) == 0 {
				if !dryRun {
					settled, err := settleEdit(repo.Root, mf.Template, lk, src, m.Version, a.now())
					if err != nil {
						return err
					}
					if settled {
						fmt.Fprintf(a.out, "%s and the lock now record template %s at %s\n", manifest.File, targetLabel(src), m.Version)
					}
				}
				fmt.Fprintf(a.out, "already at template %s (use --force to re-apply)\n", m.Version)
				return nil
			}
			if lkErr != nil && !relock {
				return lkErr
			}
			if lkErr != nil {
				// --relock rewrites a lock it cannot read, without its values.
				a.logger().Warn("lock unreadable, relocking without its recorded variables", "component", "cmd", "err", lkErr)
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
				if err := writeTemplateFields(filepath.Join(repo.Root, manifest.File), src, fields, m.Version, a.now()); err != nil {
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
			plan, err := upgrade.Compute(repo.Root, m, src, opt, lk, at)
			if err != nil {
				return err
			}
			if a.jsonOut && dryRun {
				return a.printJSON(map[string]any{"target": target, "plan": plan})
			}
			if !a.jsonOut {
				fmt.Fprintln(a.out, upgrade.Describe(plan))
				if plan.NoLock {
					fmt.Fprintln(a.out, "no system-flow.lock.yaml: every differing file counts as a conflict (flai upgrade --relock records the current state)")
				}
				for _, ch := range plan.Changes {
					if ch.Class != upgrade.Same {
						fmt.Fprintf(a.out, "  %-9s %s\n", ch.Class, ch.Path)
					}
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
			if err := writeTemplateFields(filepath.Join(repo.Root, manifest.File), src, fields, m.Version, a.now()); err != nil {
				return err
			}
			kept := 0
			for _, v := range policy {
				if v == "keep" {
					kept++
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"target": target, "plan": plan, "written": written, "kept": kept})
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
	f.StringVar(&ref, "ref", "", "template branch, tag, or commit to apply, recorded as template.ref; 1.0.60 matches v1.0.60 (default: the newest version tag when the manifest's template.ref follows releases, else that ref)")
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

// Why an upgrade applies the ref it does (ADR-0103).
const (
	reasonNewest   = "the newest release"
	reasonManifest = manifest.File + " names it"
	reasonChosen   = "chosen at the prompt"
	reasonRef      = "--ref names it"
	reasonAsGiven  = "template.ref in " + manifest.File + " names a branch or commit, which does not follow releases"
	reasonNoTags   = "the template has no version tags, so template.ref in " + manifest.File + " is used as given"
	reasonLocal    = "a local template directory, used as it is"
	reasonTemplate = "--template names it"
	reasonRelock   = "--relock records the files at template.ref in " + manifest.File
)

// upgradeTarget is the template ref an upgrade applies, its version, and why.
type upgradeTarget struct {
	Ref     string `json:"ref"`
	Version string `json:"version"`
	Reason  string `json:"reason"`
}

// versionConflict is an edit to system-flow.yaml that names a release other
// than the newest: the operator chooses between them.
type versionConflict struct {
	Edited upgradeTarget `json:"system_flow_yaml"`
	Newest upgradeTarget `json:"newest"`
}

func (c versionConflict) String() string {
	return fmt.Sprintf("%s names template %s (%s), but the newest release is %s (%s)", manifest.File, c.Edited.Ref, c.Edited.Version, c.Newest.Ref, c.Newest.Version)
}

// refusal is the error of a run that would ask but cannot: it names the
// command that applies each choice.
func (c versionConflict) refusal() error {
	return &exitError{code: 1, msg: fmt.Sprintf("%s; nothing was changed. Run one of:\n  flai upgrade --ref %s   # the version %s names\n  flai upgrade --ref %s   # the newest release\nor run flai upgrade in a terminal, without --yes, to be asked", c, c.Edited.Ref, manifest.File, c.Newest.Ref)}
}

// projectVersion is the template version a project is at: the lock's, or
// the manifest's when there is no lock.
func projectVersion(mf manifest.Template, lk *lock.Lock) string {
	if lk != nil && lk.Template.Version != "" {
		return lk.Template.Version
	}
	return mf.Version
}

// manifestEdited reports whether system-flow.yaml names another template ref
// or version than the lock recorded, because the operator edited it.
func manifestEdited(mf manifest.Template, lk *lock.Template) bool {
	return lk != nil && (mf.Ref != lk.Ref || mf.Version != lk.Version)
}

// chooseTarget decides which ref an upgrade with no --ref applies to a git
// template (ADR-0103), from the template's remote, the manifest's template
// fields, and the lock's (nil without a lock). It answers a conflict instead
// when an edit to the manifest names a release other than the newest, and an
// error when an edited version has no release tag.
func chooseTarget(rm template.Remote, mf manifest.Template, lk *lock.Template) (upgradeTarget, *versionConflict, error) {
	newest, ok := rm.Latest()
	if !ok {
		return upgradeTarget{Ref: mf.Ref, Reason: reasonNoTags}, nil, nil
	}
	latest := upgradeTarget{Ref: newest.Name, Version: newest.Version, Reason: reasonNewest}
	asGiven := upgradeTarget{Ref: mf.Ref, Reason: reasonAsGiven}
	named := func(tag template.Tag) (upgradeTarget, *versionConflict, error) {
		edited := upgradeTarget{Ref: tag.Name, Version: tag.Version, Reason: reasonManifest}
		if tag.Name == newest.Name {
			edited.Reason = reasonManifest + ", the newest release"
			return edited, nil, nil
		}
		return upgradeTarget{}, &versionConflict{Edited: edited, Newest: latest}, nil
	}
	if !manifestEdited(mf, lk) {
		if rm.FollowsReleases(mf.Ref) {
			return latest, nil, nil
		}
		return asGiven, nil, nil
	}
	if mf.Ref != lk.Ref {
		if tag, ok := rm.Match(mf.Ref); ok {
			return named(tag)
		}
		if !rm.FollowsReleases(mf.Ref) {
			return asGiven, nil, nil
		}
	}
	if mf.Version == lk.Version {
		return latest, nil, nil
	}
	tag, ok := rm.TagFor(mf.Version)
	if !ok {
		names := make([]string, len(rm.Tags))
		for i, t := range rm.Tags {
			names[i] = t.Name
		}
		return upgradeTarget{}, nil, fmt.Errorf("%s names template.version %s, but the template has no release tag for it; its releases are %s. Set template.version to one of them, or run flai upgrade --ref <tag>: nothing was changed", manifest.File, mf.Version, strings.Join(names, ", "))
	}
	return named(tag)
}

// pickTarget resolves the ref flai upgrade applies: --ref, or with none the
// one chooseTarget decides for a git template, asking on a terminal when it
// answers a conflict. ok is false when nothing is to be applied: the operator
// chose to change nothing, or a dry run would ask; either has been said.
func (a *app) pickTarget(mf manifest.Template, lk *lock.Lock, tplRepo, ref string, relock, dryRun bool, at string) (target upgradeTarget, ok bool, err error) {
	repo := tplRepo
	if repo == "" {
		repo = mf.Repo
	}
	switch {
	case ref != "":
		// releaseRef (new.go) finds the release tag --ref names, so 1.0.60
		// finds v1.0.60; a ref that names none is used as given.
		_, given, _, err := a.releaseRef(repo, ref)
		if err != nil {
			return upgradeTarget{}, false, err
		}
		return upgradeTarget{Ref: given, Reason: reasonRef}, true, nil
	case tplRepo != "":
		return upgradeTarget{Reason: reasonTemplate}, true, nil
	case template.IsLocal(repo):
		return upgradeTarget{Ref: mf.Ref, Reason: reasonLocal}, true, nil
	case relock:
		return upgradeTarget{Ref: mf.Ref, Reason: reasonRelock}, true, nil
	}
	rm, err := template.ListRemote(a.runner, repo)
	if err != nil {
		return upgradeTarget{}, false, err
	}
	var lt *lock.Template
	if lk != nil {
		lt = &lk.Template
	}
	target, conflict, err := chooseTarget(rm, mf, lt)
	if err != nil || conflict == nil {
		return target, err == nil, err
	}
	return a.settleConflict(*conflict, at, dryRun)
}

// settleConflict asks a terminal which version to apply. A dry run says it
// would ask; without a terminal, or with --yes, the run is refused.
func (a *app) settleConflict(c versionConflict, at string, dryRun bool) (upgradeTarget, bool, error) {
	if dryRun {
		if a.jsonOut {
			return upgradeTarget{}, false, a.printJSON(map[string]any{"would_ask": c})
		}
		fmt.Fprintf(a.out, "%s\ndry run: flai upgrade would ask which to apply; nothing changed\n", c)
		return upgradeTarget{}, false, nil
	}
	if a.yes || !a.isTerminal() {
		return upgradeTarget{}, false, c.refusal()
	}
	choice, err := a.prompts().Select(c.String()+". Which version should flai upgrade apply?", []prompt.Option{
		{Label: fmt.Sprintf("%s (%s), the version %s names", c.Edited.Ref, c.Edited.Version, manifest.File), Value: "edited"},
		{Label: fmt.Sprintf("%s (%s), the newest release", c.Newest.Ref, c.Newest.Version), Value: "newest"},
		{Label: fmt.Sprintf("change nothing, stay at %s", at), Value: "none"},
	}, "edited")
	if errors.Is(err, prompt.ErrNoAnswer) {
		return upgradeTarget{}, false, c.refusal()
	}
	if err != nil {
		return upgradeTarget{}, false, err
	}
	switch choice {
	case "edited":
		c.Edited.Reason = reasonChosen
		return c.Edited, true, nil
	case "newest":
		c.Newest.Reason = reasonChosen
		return c.Newest, true, nil
	}
	fmt.Fprintf(a.out, "nothing changed: the project stays at template %s\n", at)
	return upgradeTarget{}, false, nil
}

// targetLabel names the template ref a source was resolved at.
func targetLabel(src template.Source) string {
	switch {
	case src.Local:
		return src.Dir
	case src.Ref == "":
		return src.Repo + " (its default branch)"
	}
	return src.Ref
}

// settleEdit records the ref and version of a git template in system-flow.yaml
// and the lock when the project is already at that version but the manifest
// was edited to name another, so the choice this run made is not asked again.
func settleEdit(root string, mf manifest.Template, lk *lock.Lock, src template.Source, version string, now time.Time) (bool, error) {
	if src.Local || lk == nil || !manifestEdited(mf, &lk.Template) {
		return false, nil
	}
	lk.Template.Ref, lk.Template.Version = src.Ref, version
	if err := lock.Save(root, lk); err != nil {
		return false, err
	}
	return true, writeTemplateFields(filepath.Join(root, manifest.File), src, templateFields{ref: true}, version, now)
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

// templateFields names the manifest's template fields an upgrade rewrites
// besides version and applied.
type templateFields struct{ repo, ref bool }

// writeTemplateFields updates template.version and applied in the manifest
// in place, and repo and ref when fields names them, touching nothing else.
func writeTemplateFields(path string, src template.Source, fields templateFields, version string, now time.Time) error {
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
		case inTemplate && strings.HasPrefix(l, "  repo:") && fields.repo:
			lines[i] = fmt.Sprintf("  repo: %q", repo)
		case inTemplate && strings.HasPrefix(l, "  ref:") && fields.ref:
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
