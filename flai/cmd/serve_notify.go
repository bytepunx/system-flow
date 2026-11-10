package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

// flai serve notify: the channel flai serve reminds the operator through of
// a thread that has awaited them past the manifest's threads.remind_after
// (S-0498, ADR-0205), a webhook, a shell command, or both, kept on the host
// beside the providers and forges.
func newServeNotifyCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "notify",
		Short: "The channel flai serve reminds you through of a thread that awaits you, on this host",
		Long: `flai serve reminds you of each thread that has awaited you longer than a
project's threads.remind_after, and again every threads.remind_every, through
the channel you set here, for every project it serves on this host. It writes
nothing on the thread. With no channel set it only logs the reminder.

--url is a webhook, http or https, that each reminder is POSTed to as JSON:
one attempt, a five second timeout, an answer other than 2xx a failure.
--command is run through the host's shell, sh -c, with the reminder's JSON on
its standard input and a 30 second timeout. Set either or both; set changes
only what you give it. The JSON holds the project, the thread, its title and
item, how long it has awaited you and since when, a pending recommendation's
author, first line, and source, and the dashboard's link to it, plus the
entry the dashboard posts to dashboard.notify_url, so one receiver can take
both. No token and no file content is sent.`,
		Example: `  flai serve notify set --url https://ntfy.sh/<topic>
  flai serve notify set --command 'curl -s -H "Title: flai" -d @- https://ntfy.sh/<topic>'
  flai serve notify show --json
  flai serve notify clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showNotify() },
	}
	var target, command string
	set := &cobra.Command{
		Use:   "set [--url <url>] [--command <command>]",
		Short: "Set the webhook, the shell command, or both, that reminders are sent to",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			urlGiven, commandGiven := cmd.Flags().Changed("url"), cmd.Flags().Changed("command")
			if !urlGiven && !commandGiven {
				return fmt.Errorf("give --url <webhook>, --command <shell command>, or both; flai serve notify clear removes the channel")
			}
			if urlGiven && strings.TrimSpace(target) == "" {
				return fmt.Errorf("--url needs an http or https webhook, such as https://ntfy.sh/<topic>; flai serve notify clear removes the channel")
			}
			if commandGiven && strings.TrimSpace(command) == "" {
				return fmt.Errorf("--command needs a shell command that reads the reminder's JSON on its standard input; flai serve notify clear removes the channel")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			next := cfg.Notify
			if urlGiven {
				next.URL = strings.TrimSpace(target)
			}
			if commandGiven {
				next.Command = command
			}
			if cfg, err = cfg.WithNotify(next); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showNotify()
		},
	}
	set.Flags().StringVar(&target, "url", "", "an http or https webhook each reminder is POSTed to as JSON")
	set.Flags().StringVar(&command, "command", "", "a shell command run with each reminder's JSON on its standard input")
	clear := &cobra.Command{
		Use:   "clear",
		Short: "Remove the channel, so that flai serve only logs a reminder",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if err := config.Save(path, cfg.WithoutNotify()); err != nil {
				return err
			}
			return a.showNotify()
		},
	}
	c.AddCommand(set, clear,
		&cobra.Command{Use: "show", Short: "Print the webhook and the shell command reminders are sent to", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showNotify() }},
	)
	return c
}

// showNotify prints the host's reminder channel: the webhook and the shell
// command, each or none. The operator set both in their own shell, so the
// URL is printed whole, secret and all, as they wrote it.
func (a *app) showNotify() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	n := cfg.Notify
	if a.jsonOut {
		return a.printJSON(map[string]any{"url": n.URL, "command": n.Command})
	}
	if n == (config.Notify{}) {
		fmt.Fprintln(a.out, "no reminder channel is set; flai serve only logs a reminder (flai serve notify set --url <url> or --command <command>)")
		return nil
	}
	if n.URL != "" {
		fmt.Fprintf(a.out, "url: %s\n", n.URL)
	}
	if n.Command != "" {
		fmt.Fprintf(a.out, "command: %s\n", n.Command)
	}
	return nil
}
