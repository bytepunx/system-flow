package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// componentsIdle is the operator's components.idle, read from the
// configuration at every look, so that a change takes effect without a
// restart; serve.DefaultComponentsIdle when it is not set or not readable.
func (a *app) componentsIdle() time.Duration {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.DefaultComponentsIdle
	}
	idle, _ := componentsIdleOf(cfg)
	return idle
}

// componentsIdleOf is cfg's components.idle, and whether it is set; the
// default when it is not, or when it is not a positive duration.
func componentsIdleOf(cfg config.Config) (time.Duration, bool) {
	d, err := time.ParseDuration(strings.TrimSpace(cfg.Components.Idle))
	if err != nil || d <= 0 {
		return serve.DefaultComponentsIdle, false
	}
	return d, true
}

// componentsMaxServers is the operator's components.max_servers, read from
// the configuration at every start, so that a change takes effect without a
// restart; serve.DefaultComponentsMaxServers when it is not set or not
// readable.
func (a *app) componentsMaxServers() int {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.DefaultComponentsMaxServers
	}
	n, _ := componentsMaxServersOf(cfg)
	return n
}

// componentsMaxServersOf is cfg's components.max_servers, and whether it is
// set; the default when it is not, or when it is not positive.
func componentsMaxServersOf(cfg config.Config) (int, bool) {
	if cfg.Components.MaxServers <= 0 {
		return serve.DefaultComponentsMaxServers, false
	}
	return cfg.Components.MaxServers, true
}

// flai serve components: the host's say about the component dev servers
// flai serve runs under the components host action (S-0438).
func newServeComponentsCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "components",
		Short: "How long flai serve keeps a component dev server nobody looks at, and how many it runs, on this host",
		Long: `flai serve runs a Svelte sub-project's own Vite dev server, with flai's
component harness, when the dashboard's Components page asks it to
(components.start), on 127.0.0.1 and a port flai picks. It runs the
project's code on this host, so it is a host action, off until you enable it
(flai serve enable components). Each start and stop is in flai serve journal.

flai serve stops a dev server that had no components.status, which the page
asks while it is open, for components.idle, 30 minutes by default, and every
dev server when it stops. flai serve components idle sets it, as a Go
duration such as 45m or 2h, for every project on this host; it takes effect
at once, with no restart.

The page also runs a dev server in the worktree of a story in progress or in
review, beside the main checkout's, on a port of its own (S-0442); flai serve
stops it once the story leaves review, to done, cancelled, or back to in
progress, or is no longer in progress or in review at all. At most
components.max_servers dev servers run or start at once per project, main's
and every story's together, 3 by default; a start beyond it is refused,
naming the ones running. flai serve components max-servers sets it, for
every project on this host, at once.`,
		Example: `  flai serve enable components
  flai serve components idle 45m
  flai serve components max-servers 4
  flai serve components show --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showComponentsSettings() },
	}
	idle := &cobra.Command{
		Use:   "idle <duration>",
		Short: "Stop a dev server after this long without a components.status (default 30m)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			given := strings.TrimSpace(args[0])
			d, err := time.ParseDuration(given)
			if err != nil || d <= 0 {
				return fmt.Errorf("components.idle must be a positive duration, such as 30m or 2h, not %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Components.Idle = given
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showComponentsSettings()
		},
	}
	maxServers := &cobra.Command{
		Use:   "max-servers <n>",
		Short: "Run at most this many component dev servers at once per project, main's and every story's (default 3)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			n, err := strconv.Atoi(strings.TrimSpace(args[0]))
			if err != nil || n <= 0 {
				return fmt.Errorf("components.max_servers must be a positive whole number, such as 3, not %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Components.MaxServers = n
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showComponentsSettings()
		},
	}
	c.AddCommand(idle, maxServers,
		&cobra.Command{Use: "show", Short: "Print components.idle, components.max_servers, and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showComponentsSettings() }},
	)
	return c
}

// showComponentsSettings prints components.idle, components.max_servers, and
// whether the components action is on for the project in the working
// directory.
func (a *app) showComponentsSettings() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := "", false
	if repo, err := a.project(); err == nil {
		here = mainRootOf(repo)
		enabled = cfg.ActionEnabled(hostapi.ActionComponents, here)
	}
	idle, set := componentsIdleOf(cfg)
	most, mostSet := componentsMaxServersOf(cfg)
	if a.jsonOut {
		return a.printJSON(map[string]any{"idle": shortDuration(idle), "idle_seconds": int(idle / time.Second), "idle_set": set, "max_servers": most, "max_servers_set": mostSet, "enabled_here": enabled})
	}
	byDefault := func(set bool) string {
		if set {
			return ""
		}
		return " (default)"
	}
	fmt.Fprintf(a.out, "idle: %s%s: a dev server without a components.status for this long is stopped\n", shortDuration(idle), byDefault(set))
	fmt.Fprintf(a.out, "max_servers: %d%s: at most this many dev servers run at once per project, main's and every story's together\n", most, byDefault(mostSet))
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the components action is on for this project; flai serve disable components turns it off")
		} else {
			fmt.Fprintln(a.out, "the components action is off for this project; flai serve enable components turns it on")
		}
	}
	return nil
}

// shortDuration is d as Go writes it, without the zero units it ends with:
// 30m rather than 30m0s.
func shortDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}
