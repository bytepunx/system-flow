package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/importer"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

type importOptions struct {
	newOptions
	dryRun   bool
	commit   bool
	trailers []string
}

func newImportCmd(a *app) *cobra.Command {
	var o importOptions
	c := &cobra.Command{
		Use:   "import [dir]",
		Short: "Bring an existing repository under the system-flow standard",
		Long: `Analyse a repository, propose the layout, and apply it: create the
documentation folders, render the template files that do not already exist,
offer to move existing documentation folders and markdown files into the
new structure (git mv when tracked), detect code sub-projects by their build
files, and write system-flow.yaml. Nothing existing is overwritten.

Interactive in a terminal; --yes accepts every default; --dry-run prints the
proposal and stops.

--commit (S-0098) then runs the repository's tests and commits the import:
the host's checks when flai serve checks set names any, else what the
repository has (its own Makefile's test target, go test, the package
manager's test script, cargo test, pytest), and commits exactly the paths
the import wrote or moved, only when every test passed or none were found.
It needs a git repository with no uncommitted changes. When a test fails,
the imported files are left uncommitted, the answer says which failed, and
flai exits with code 5.

Once the manifest is written, a flai host running for this config is told
to serve the project (S-0117): it is registered with flai serve as flai
dashboard registers it, and the dashboard shows it in its switcher. With no
host running, nothing is registered, and flai dashboard in the project
serves it.`,
		Example: `  flai import --dry-run
  flai import
  flai import ../legacy --yes --layout design=architecture
  flai import --var description="Billing platform" --var owner=core
  flai import ../legacy --yes --commit`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			return a.runImport(dir, o)
		},
	}
	f := c.Flags()
	f.StringVar(&o.templateRepo, "template", "", "template git URL or local directory (default: config template.repo)")
	f.StringVar(&o.ref, "ref", "", "template branch, tag, or commit (default: config template.ref)")
	f.StringArrayVar(&o.vars, "var", nil, "set a template variable, name=value (repeatable)")
	f.StringArrayVar(&o.layout, "layout", nil, "name a layout folder, key=name (repeatable); defaults are proposed from what exists")
	f.BoolVar(&o.dryRun, "dry-run", false, "print the proposal and change nothing")
	f.BoolVar(&o.force, "force", false, "proceed even if system-flow.yaml already exists")
	f.BoolVar(&o.commit, "commit", false, "then run the repository's tests and commit the import when they pass")
	f.StringArrayVar(&o.trailers, "trailer", nil, "trailer line for the import's commit (repeatable, with --commit)")
	return c
}

func (a *app) runImport(dir string, o importOptions) error {
	if _, err := givenVars(o.vars); err != nil {
		return err
	}
	an, err := importer.Scan(dir)
	if err != nil {
		return err
	}
	if an.Manifest && !o.force {
		return fmt.Errorf("%s already has %s; use flai check, or --force to re-import", an.Root, manifest.File)
	}
	if o.commit && !o.dryRun {
		if err := a.importPreflight(an); err != nil {
			return err
		}
	}
	ownMakefile := ownMakefileAt(an.Root)
	src, m, err := a.resolveTemplate(o.templateRepo, o.ref, false)
	if err != nil {
		return err
	}
	interactive := !a.yes && !o.dryRun && a.isTerminal()

	// layout names: flags, then existing folders, then template defaults
	layout := importer.Layout{}
	given, err := parsePairs(o.layout, "layout")
	if err != nil {
		return err
	}
	for _, k := range m.LayoutKeys() {
		switch {
		case given[k] != "":
			layout[k] = given[k]
		case an.Existing[k] != "":
			layout[k] = an.Existing[k]
		default:
			layout[k] = m.Layout[k]
		}
	}
	if interactive {
		for _, k := range m.LayoutKeys() {
			v := layout[k]
			if err := huh.NewForm(huh.NewGroup(huh.NewInput().Title("Folder for " + k).Description("Default from the template or an existing folder").Value(&v))).Run(); err != nil {
				return err
			}
			if strings.TrimSpace(v) != "" {
				layout[k] = strings.TrimSpace(v)
			}
		}
	}
	plan := importer.BuildPlan(an, layout)

	if a.jsonOut && o.dryRun {
		tests, from := a.importTests(an.Root, ownMakefile, plan.Projects)
		if tests == nil {
			tests = []importer.TestCommand{}
		}
		return a.printJSON(map[string]any{"analysis": an, "plan": plan, "tests": tests, "tests_from": from})
	}
	if !a.jsonOut {
		a.printProposal(an, plan, src, m)
	}
	if o.dryRun {
		fmt.Fprintln(a.out, "\ndry run: nothing changed")
		return nil
	}
	if interactive {
		ok := true
		if err := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Apply this proposal?").Value(&ok))).Run(); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("import cancelled")
		}
	}

	// 1. folder moves (whole candidate folders such as adr/ -> design/adrs)
	var moved, kept []importer.Move
	for _, mv := range plan.FolderMoves {
		do := true
		if interactive {
			if err := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(fmt.Sprintf("Move %s/ into %s/?", mv.From, mv.To)).Value(&do))).Run(); err != nil {
				return err
			}
		}
		if !do {
			kept = append(kept, mv)
			continue
		}
		ms, err := importer.MoveFolderContents(a.runner, an.Root, mv.From, mv.To, an.Git)
		if err != nil {
			return fmt.Errorf("move %s: %w", mv.From, err)
		}
		moved = append(moved, ms...)
	}

	// 2. render the template without overwriting anything
	origin := ""
	if an.Git {
		origin = a.originURL(an.Root)
	}
	vars, err := a.collectVars(m, newOptions{vars: o.vars, defaults: !interactive, origin: origin}, filepath.Base(an.Root))
	if err != nil {
		return err
	}
	res, err := template.Render(m, src.Dir, an.Root, template.Options{Vars: vars, Layout: layout, Source: src, Force: false})
	if err != nil {
		return err
	}

	// 3. loose markdown: ask per file, or leave in place
	dests := []string{"leave", layout["design"] + "/system", layout["design"] + "/adrs", layout["docs"] + "/users", layout["docs"] + "/operators", layout["docs"] + "/contributors"}
	skipAll := false
	for _, md := range plan.Markdown {
		if !interactive || skipAll {
			break
		}
		choice := "leave"
		opts := make([]huh.Option[string], 0, len(dests)+1)
		for _, d := range dests {
			opts = append(opts, huh.NewOption(d, d))
		}
		opts = append(opts, huh.NewOption("skip all remaining", "skip"))
		if err := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title("Where does " + md + " belong?").Options(opts...).Value(&choice))).Run(); err != nil {
			return err
		}
		switch choice {
		case "leave":
		case "skip":
			skipAll = true
		default:
			to := filepath.ToSlash(filepath.Join(choice, filepath.Base(md)))
			if _, err := os.Stat(filepath.Join(an.Root, to)); err == nil {
				fmt.Fprintf(a.out, "  kept %s: %s already exists\n", md, to)
				continue
			}
			if err := importer.MoveFile(a.runner, an.Root, md, to, an.Git); err != nil {
				return err
			}
			moved = append(moved, importer.Move{From: md, To: to})
		}
	}

	// 4. manifest with detected projects (the template rendered the base file)
	if err := writeProjects(filepath.Join(an.Root, manifest.File), plan.Projects); err != nil {
		return err
	}

	// 5. check
	repo, err := workitem.Open(an.Root)
	if err != nil {
		return err
	}
	resCheck, err := check.Run(repo, a.now())
	if err != nil {
		return err
	}
	var committed *importCommit
	if o.commit {
		c, err := a.commitImport(an.Root, ownMakefile, plan.Projects, res.Written, moved, o.trailers)
		if err != nil {
			return err
		}
		committed = &c
	}
	// 6. served by the host flai, so that the dashboard shows it (S-0117)
	served := a.serveImported(repo)
	if a.jsonOut {
		out := map[string]any{"root": an.Root, "layout": layout, "written": res.Written, "skipped": res.Skipped, "moved": moved, "kept": kept, "projects": plan.Projects, "check": resCheck, "serve": served}
		if committed != nil {
			out["key"] = repo.Manifest.Key
			out["name"] = repo.Manifest.Name
			out["commit"] = committed
		}
		if err := a.printJSON(out); err != nil {
			return err
		}
		if committed != nil && !committed.Committed {
			return &exitError{code: exitImportNotCommitted, msg: committed.Reason}
		}
		return nil
	}
	fmt.Fprintf(a.out, "\nImported %s\n  %d template files written, %d existing files kept\n", an.Root, len(res.Written), len(res.Skipped))
	for _, mv := range moved {
		fmt.Fprintf(a.out, "  moved %s -> %s\n", mv.From, mv.To)
	}
	for _, mv := range kept {
		fmt.Fprintf(a.out, "  kept %s/ in place (proposed %s/)\n", mv.From, mv.To)
	}
	if len(plan.Projects) > 0 {
		fmt.Fprintf(a.out, "  %d sub-projects recorded in %s\n", len(plan.Projects), manifest.File)
	}
	fmt.Fprintf(a.out, "  flai check: %d errors, %d warnings\n", resCheck.Errors, resCheck.Warnings)
	for _, f := range resCheck.Findings {
		fmt.Fprintf(a.out, "    %s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
	}
	a.printImportServed(served)
	if committed != nil {
		a.printImportCommit(*committed)
		if !committed.Committed {
			return &exitError{code: exitImportNotCommitted, msg: committed.Reason}
		}
		fmt.Fprintf(a.out, "\nNext:\n  read %s/conventions/README.md\n  flai epic new \"First deliverable\"\n", layout["design"])
		return nil
	}
	fmt.Fprintf(a.out, "\nNext:\n  review the moves with git status\n  read %s/conventions/README.md\n  flai epic new \"First deliverable\"\n", layout["design"])
	return nil
}

func (a *app) printImportCommit(c importCommit) {
	switch c.Source {
	case "none":
		fmt.Fprintln(a.out, "  tests: none found (no Makefile test target, go.mod, package.json test script, Cargo.toml, or pytest setup)")
	case "host":
		fmt.Fprintln(a.out, "  tests: the checks this host names (flai serve checks)")
	}
	for _, t := range c.Tests {
		state := "passed"
		if !t.OK {
			state = fmt.Sprintf("FAILED (exit %d)", t.ExitCode)
		}
		fmt.Fprintf(a.out, "  %s: %s in %ds\n", t.Name, state, t.Seconds)
		if !t.OK {
			for _, line := range strings.Split(strings.TrimRight(t.Output, "\n"), "\n") {
				fmt.Fprintf(a.out, "      %s\n", line)
			}
		}
	}
	if c.Committed {
		fmt.Fprintf(a.out, "  committed as %s\n", c.Commit)
	} else {
		fmt.Fprintf(a.out, "  not committed: %s\n", c.Reason)
	}
}

func (a *app) printProposal(an *importer.Analysis, plan *importer.Plan, src template.Source, m template.Manifest) {
	fmt.Fprintf(a.out, "Import proposal for %s (template %s %s)\n", an.Root, src.Repo, m.Version)
	if an.Git {
		fmt.Fprintln(a.out, "  git repository: moves use git mv")
	}
	fmt.Fprintln(a.out, "  layout:")
	keys := make([]string, 0, len(plan.Layout))
	for k := range plan.Layout {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		note := ""
		if an.Existing[k] == plan.Layout[k] {
			note = " (exists)"
		}
		fmt.Fprintf(a.out, "    %-8s %s%s\n", k, plan.Layout[k], note)
	}
	if len(plan.Create) > 0 {
		fmt.Fprintf(a.out, "  create: %s\n", strings.Join(plan.Create, ", "))
	}
	for _, mv := range plan.FolderMoves {
		fmt.Fprintf(a.out, "  move folder: %s/ -> %s/\n", mv.From, mv.To)
	}
	if len(plan.Markdown) > 0 {
		fmt.Fprintf(a.out, "  markdown to place (%d): %s\n", len(plan.Markdown), strings.Join(plan.Markdown, ", "))
	}
	for _, p := range plan.Projects {
		fmt.Fprintf(a.out, "  sub-project: %s (%s) at %s\n", p.Name, p.Kind, p.Path)
	}
	fmt.Fprintln(a.out, "  template files that already exist are kept; nothing is overwritten")
}

// writeProjects replaces the projects list in a rendered manifest.
func writeProjects(path string, projects []importer.Project) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("projects:")
	if len(projects) == 0 {
		b.WriteString(" []\n")
	} else {
		b.WriteString("\n")
		for _, p := range projects {
			fmt.Fprintf(&b, "  - name: %s\n    path: %s\n    kind: %s\n", p.Name, p.Path, p.Kind)
		}
	}
	s := strings.Replace(string(data), "projects: []\n", b.String(), 1)
	return os.WriteFile(path, []byte(s), 0o644)
}
