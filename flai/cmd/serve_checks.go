package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// flai serve checks: the named commands a story in review is checked with
// (S-0082). Unlike agent, there can be several, so set names one at a time
// and clear takes an optional name; show lists all of them.
func newServeChecksCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "checks",
		Short: "The named commands a story in review is checked with, on this host",
		Long: `flai serve can run checks for a story in review, in the story's worktree, and
keep the outcome with it until it is accepted. It is a host action, off
until you enable it (flai serve enable checks), and by default nothing is
named here: the manifest's own checks: is used instead. Naming any command
here, on this host, uses this list instead of the manifest's, for this host
only.

Each command is an argument list, run as it stands in the worktree, as you,
never through a shell. In an argument {story} is replaced by the story's ID
and {root} by the worktree's directory; nothing else is interpreted. Every
named command runs in order; the first to fail stops the rest. One run per
story at a time; every run and every refusal is in flai serve journal.`,
		Example: `  flai serve checks set --name flai -- scripts/flai-test.sh
  flai serve checks set --name flaiover -- bash -c "cd flaiover && pnpm run check && pnpm run test:unit -- run && pnpm run lint"
  flai serve checks show
  flai serve enable checks
  flai serve checks clear flaiover
  flai serve checks timeout 20`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showChecksCommands() },
	}
	var name string
	set := &cobra.Command{
		Use:   "set --name <name> -- <program> [args...]",
		Short: "Add or replace one named command, as an argument list after --",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("--name is required")
			}
			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("the program is empty")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			replaced := false
			for i, nc := range cfg.Checks.Commands {
				if nc.Name == name {
					cfg.Checks.Commands[i].Command = args
					replaced = true
					break
				}
			}
			if !replaced {
				cfg.Checks.Commands = append(cfg.Checks.Commands, config.NamedCommand{Name: name, Command: args})
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	set.Flags().StringVar(&name, "name", "", "the check's name (required)")
	timeout := &cobra.Command{
		Use:   "timeout <minutes>",
		Short: "Bound how long one run of every command together may take (default 15)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			minutes, err := strconv.Atoi(args[0])
			if err != nil || minutes <= 0 {
				return fmt.Errorf("timeout must be a positive number of minutes, got %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Checks.TimeoutMinutes = minutes
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	clear := &cobra.Command{
		Use:   "clear [name]",
		Short: "Remove one named command, or every one when no name is given",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				cfg.Checks.Commands = nil
			} else {
				kept := cfg.Checks.Commands[:0]
				for _, nc := range cfg.Checks.Commands {
					if nc.Name != args[0] {
						kept = append(kept, nc)
					}
				}
				cfg.Checks.Commands = kept
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	c.AddCommand(set, timeout, clear,
		&cobra.Command{Use: "show", Short: "Print the named commands and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showChecksCommands() }},
	)
	return c
}

func (a *app) showChecksCommands() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := "", false
	if repo, err := a.project(); err == nil {
		here = mainRootOf(repo)
		enabled = cfg.ActionEnabled(hostapi.ActionChecks, here)
	}
	timeout := cfg.Checks.TimeoutMinutes
	if timeout <= 0 {
		timeout = 15
	}
	if a.jsonOut {
		cmds := cfg.Checks.Commands
		if cmds == nil {
			cmds = []config.NamedCommand{}
		}
		return a.printJSON(map[string]any{"commands": cmds, "timeout_minutes": timeout, "enabled_here": enabled})
	}
	if len(cfg.Checks.Commands) == 0 {
		fmt.Fprintln(a.out, "no command is named here; the manifest's own checks: is used instead")
	} else {
		for _, nc := range cfg.Checks.Commands {
			quoted := make([]string, len(nc.Command))
			for i, arg := range nc.Command {
				quoted[i] = fmt.Sprintf("%q", arg)
			}
			fmt.Fprintf(a.out, "%s: %s\n", nc.Name, strings.Join(quoted, " "))
		}
	}
	fmt.Fprintf(a.out, "timeout: %d minute%s\n", timeout, map[bool]string{true: "", false: "s"}[timeout == 1])
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the checks action is on for this project; flai serve disable checks turns it off")
		} else {
			fmt.Fprintln(a.out, "the checks action is off for this project; flai serve enable checks turns it on")
		}
	}
	return nil
}
