package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

func newConfigCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Read and write ~/.flai/config.json",
		Long: `Read and write the flai configuration file.

Keys:
  ` + strings.Join(config.Keys(), "\n  "),
	}
	c.AddCommand(newConfigGetCmd(a), newConfigSetCmd(a), newConfigPathCmd(a))
	return c
}

func newConfigGetCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Print the whole config, or one dotted key",
		Example: `  flai config get
  flai config get template.repo
  flai config get dashboard --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				return a.printJSON(cfg)
			}
			v, err := cfg.Get(args[0])
			if err != nil {
				return err
			}
			if _, isObj := v.(map[string]any); isObj || a.jsonOut {
				return a.printJSON(v)
			}
			_, err = fmt.Fprintln(a.out, v)
			return err
		},
	}
}

func newConfigSetCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set one dotted key and save",
		Example: `  flai config set template.repo git@github.com:me/system-flow-template.git
  flai config set template.ref my-branch
  flai config set dashboard.port 8080`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			updated, err := cfg.Set(args[0], args[1])
			if err != nil {
				return err
			}
			if err := config.Save(path, updated); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"path": path, "key": args[0], "value": args[1]})
			}
			_, err = fmt.Fprintf(a.out, "%s = %s\n", args[0], args[1])
			return err
		},
	}
}

func newConfigPathCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path in use",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.ResolvePath(a.configPath)
			if a.jsonOut {
				return a.printJSON(map[string]string{"path": path})
			}
			_, err := fmt.Fprintln(a.out, path)
			return err
		},
	}
}
