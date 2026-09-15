package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

func newTemplateCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "template",
		Short: "Inspect, refresh, and switch the template source",
	}
	c.AddCommand(newTemplateShowCmd(a), newTemplateUpdateCmd(a), newTemplateUseCmd(a))
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
