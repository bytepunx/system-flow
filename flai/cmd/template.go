package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/publish"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

func newTemplateCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "template",
		Short: "Inspect, refresh, and switch the template source",
	}
	c.AddCommand(newTemplateShowCmd(a), newTemplateUpdateCmd(a), newTemplateUseCmd(a), newTemplatePushCmd(a))
	return c
}

// publishTemplate pushes a local template directory to its remote.
func (a *app) publishTemplate(dir, remote, ref string, tag, force, dryRun bool) (*publish.Result, error) {
	m, err := template.LoadManifest(dir)
	if err != nil {
		return nil, err
	}
	if remote == "" {
		remote = m.Publish.Repo
	}
	if ref == "" {
		ref = m.Publish.Ref
	}
	cacheDir := config.Default().CacheDir
	if cfg, _, err := a.loadConfig(); err == nil {
		cacheDir = cfg.CacheDir
	}
	abs, _ := filepath.Abs(dir)
	return publish.Run(a.runner, publish.Options{Dir: dir, Remote: remote, Ref: ref, CacheDir: config.ExpandHome(cacheDir), Tag: tag, Force: force, DryRun: dryRun, Source: abs})
}

func newTemplatePushCmd(a *app) *cobra.Command {
	var remote, ref string
	var tag, force, dryRun bool
	c := &cobra.Command{
		Use:   "push [dir]",
		Short: "Publish a locally developed template to its git remote",
		Long: `Clone the remote branch (creating it from the default branch if missing),
replace its contents with the local template, commit with the template
version, and push. --tag also creates and pushes v<version>, refusing if it
exists. Git errors are reported verbatim; nothing is retried, and nothing is
force-pushed without --force. The directory defaults to config template.repo
when that is a local path; the remote and branch default to publish.repo and
publish.ref in template.yaml.`,
		Example: `  flai template push --dry-run
  flai template push ./template --tag
  flai template push ./template --remote git@github.com:me/my-template.git --ref main`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := ""
			if len(args) == 1 {
				dir = args[0]
			} else {
				cfg, _, err := a.loadConfig()
				if err != nil {
					return err
				}
				if !template.IsLocal(cfg.Template.Repo) {
					return fmt.Errorf("config template.repo (%s) is not a local directory; pass the template directory as an argument", cfg.Template.Repo)
				}
				dir = config.ExpandHome(cfg.Template.Repo)
			}
			res, err := a.publishTemplate(dir, remote, ref, tag, force, dryRun)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			switch {
			case res.Nothing:
				fmt.Fprintf(a.out, "nothing to push: %s %s matches template %s\n", res.Remote, res.Ref, res.Version)
			case dryRun:
				fmt.Fprintf(a.out, "would push %d change(s) to %s %s as \"%s\"", len(res.Files), res.Remote, res.Ref, strings.SplitN(res.Message, "\n", 2)[0])
				if res.Tag != "" {
					fmt.Fprintf(a.out, " and tag %s", res.Tag)
				}
				fmt.Fprintln(a.out)
				for _, f := range res.Files {
					fmt.Fprintf(a.out, "  %s\n", f)
				}
				fmt.Fprintln(a.out, "dry run: remote untouched")
			default:
				fmt.Fprintf(a.out, "pushed %s to %s %s (%d change(s), commit %s", res.Version, res.Remote, res.Ref, len(res.Files), orDefault(res.Commit, "none"))
				if res.Tag != "" {
					fmt.Fprintf(a.out, ", tag %s", res.Tag)
				}
				if res.Created {
					fmt.Fprint(a.out, ", branch created")
				}
				fmt.Fprintln(a.out, ")")
			}
			return nil
		},
	}
	c.Flags().StringVar(&remote, "remote", "", "git URL (default: publish.repo in template.yaml)")
	c.Flags().StringVar(&ref, "ref", "", "branch (default: publish.ref in template.yaml)")
	c.Flags().BoolVar(&tag, "tag", false, "also create and push v<version>")
	c.Flags().BoolVar(&force, "force", false, "force push and overwrite an existing tag")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "report the changes and commit message; push nothing")
	return c
}

func newTemplateShowCmd(a *app) *cobra.Command {
	var repo, ref string
	c := &cobra.Command{
		Use:   "show",
		Short: "Print the template source, cache location, and manifest summary",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, m, err := a.resolveTemplate(repo, ref, false)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{
					"repo": src.Repo, "ref": src.Ref, "local": src.Local, "dir": src.Dir,
					"version": m.Version, "min_flai": m.MinFlai, "description": m.Description,
					"variables": m.Variables, "layout": m.Layout,
				})
			}
			fmt.Fprintf(a.out, "repo:     %s\n", src.Repo)
			if !src.Local {
				fmt.Fprintf(a.out, "ref:      %s\n", orDefault(src.Ref, "(default branch)"))
			}
			fmt.Fprintf(a.out, "dir:      %s\n", src.Dir)
			fmt.Fprintf(a.out, "version:  %s (min flai %s)\n", m.Version, orDefault(m.MinFlai, "any"))
			if m.Description != "" {
				fmt.Fprintf(a.out, "about:    %s\n", m.Description)
			}
			fmt.Fprintln(a.out, "variables:")
			for _, v := range m.Variables {
				req := ""
				if v.Required {
					req = " (required)"
				}
				fmt.Fprintf(a.out, "  %-14s %s%s", v.Name, v.Prompt, req)
				if v.Default != "" {
					fmt.Fprintf(a.out, " [default: %s]", v.Default)
				}
				fmt.Fprintln(a.out)
			}
			fmt.Fprintln(a.out, "layout:")
			for _, k := range m.LayoutKeys() {
				fmt.Fprintf(a.out, "  %-14s %s\n", k, m.Layout[k])
			}
			return nil
		},
	}
	c.Flags().StringVar(&repo, "template", "", "template git URL or local directory (default: config)")
	c.Flags().StringVar(&ref, "ref", "", "branch, tag, or commit (default: config)")
	return c
}

func newTemplateUpdateCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Re-fetch the configured template into the cache",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, m, err := a.resolveTemplate("", "", true)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"repo": src.Repo, "ref": src.Ref, "dir": src.Dir, "version": m.Version})
			}
			if src.Local {
				fmt.Fprintf(a.out, "%s is a local directory, nothing to fetch (version %s)\n", src.Repo, m.Version)
				return nil
			}
			fmt.Fprintf(a.out, "updated %s to version %s in %s\n", src.Repo, m.Version, src.Dir)
			return nil
		},
	}
}

func newTemplateUseCmd(a *app) *cobra.Command {
	var ref string
	c := &cobra.Command{
		Use:   "use <repo>",
		Short: "Set the template source in the config file",
		Example: `  flai template use git@github.com:me/system-flow-template.git --ref my-branch
  flai template use ./template`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Template = config.Template{Repo: args[0], Ref: ref}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(cfg.Template)
			}
			fmt.Fprintf(a.out, "template.repo = %s\ntemplate.ref = %s\n", cfg.Template.Repo, orDefault(ref, "(default branch)"))
			return nil
		},
	}
	c.Flags().StringVar(&ref, "ref", "", "branch, tag, or commit; empty means the default branch")
	return c
}

func orDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
