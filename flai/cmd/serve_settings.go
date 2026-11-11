package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

// flai serve settings: the host's settings as the dashboard's settings.get
// gives them for the project in the working directory (S-0506).
func newServeSettingsCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Print this host's settings as the dashboard's settings page has them for the project in this folder",
		Long: `Prints what the dashboard's settings page is given for the project in the
working directory: the host actions, here and everywhere; the agent, its
harnesses and providers; the checks; the folders named for import; the
projects flai serve serves; planning; and the project's HTTP MCP server,
whether it runs and the address flai host starts it on, mcp.bind or ` + config.DefaultMCPBind + `
by default, which flai serve mcp bind sets. --json prints every field as
the dashboard gets it. It prints no token and no variable's value.`,
		Example: `  flai serve settings
  flai serve settings --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showHostSettings() },
	}
}

// showHostSettings prints hostSettings for the project in the working
// directory, whole as JSON or as a summary.
func (a *app) showHostSettings() error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	settings := a.hostSettings(context.Background(), mainRootOf(repo))
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	var view hostSettingsView
	if err := json.Unmarshal(data, &view); err != nil {
		return fmt.Errorf("read the host settings for %s: %w", repo.Manifest.Key, err)
	}
	if view.Error != "" {
		return fmt.Errorf("read the host settings for %s: %s", repo.Manifest.Key, view.Error)
	}
	if a.jsonOut {
		return a.printJSON(settings)
	}
	writeHostSettings(a.out, repo.Manifest.Name, mainRootOf(repo), view)
	return nil
}

// hostSettingsView is the part of hostSettings the summary shows, read
// from the JSON the dashboard gets.
type hostSettingsView struct {
	Error   string `json:"error"`
	Actions []struct {
		Name       string `json:"name"`
		Means      string `json:"means"`
		Here       bool   `json:"here"`
		Everywhere bool   `json:"everywhere"`
	} `json:"actions"`
	Agent struct {
		Name      string   `json:"name"`
		Command   []string `json:"command"`
		Harnesses map[string]struct {
			Program string            `json:"program"`
			Args    []string          `json:"args"`
			Set     bool              `json:"set"`
			Images  map[string]string `json:"images"`
			Network string            `json:"network"`
			Login   string            `json:"login"`
			Outside []string          `json:"outside"`
		} `json:"harnesses"`
		Providers map[string]struct {
			API    string `json:"api"`
			KeyEnv string `json:"key_env"`
			KeySet bool   `json:"key_set"`
		} `json:"providers"`
	} `json:"agent"`
	DefaultAgent *struct {
		Harness  string `json:"harness"`
		Model    string `json:"model"`
		Provider string `json:"provider"`
	} `json:"default_agent"`
	AllowedHosts string `json:"allowed_hosts"`
	Checks       struct {
		Commands []struct {
			Name    string   `json:"name"`
			Command []string `json:"command"`
		} `json:"commands"`
		TimeoutMinutes int `json:"timeout_minutes"`
	} `json:"checks"`
	ImportRoots []string `json:"import_roots"`
	Shared      []string `json:"shared"`
	Projects    struct {
		Error    string            `json:"error"`
		Running  bool              `json:"running"`
		Served   []json.RawMessage `json:"served"`
		Unserved []json.RawMessage `json:"unserved"`
	} `json:"projects"`
	Planning *struct {
		Plan      bool   `json:"plan"`
		Replan    string `json:"replan"`
		ReplanSet bool   `json:"replan_set"`
		Schedule  string `json:"schedule"`
		Next      string `json:"next"`
	} `json:"planning"`
	MCP *struct {
		Running bool   `json:"running"`
		URL     string `json:"url"`
		PID     int    `json:"pid"`
		Bind    string `json:"bind"`
		BindSet bool   `json:"bind_set"`
	} `json:"mcp"`
}

// writeHostSettings writes the summary of v, the host settings for the
// project name at root.
func writeHostSettings(w io.Writer, name, root string, v hostSettingsView) {
	fmt.Fprintf(w, "host settings for %s (%s); flai serve settings --json gives every field\n", name, root)
	fmt.Fprintln(w, "actions:")
	for _, act := range v.Actions {
		fmt.Fprintf(w, "  %s: %s here, %s everywhere: %s\n", act.Name, onOff(act.Here), onOff(act.Everywhere), act.Means)
	}
	command := "none"
	if len(v.Agent.Command) > 0 {
		command = strings.Join(v.Agent.Command, " ")
	}
	fmt.Fprintf(w, "agent: name %s, command %s\n", orNone(v.Agent.Name), command)
	for _, n := range slices.Sorted(maps.Keys(v.Agent.Harnesses)) {
		h := v.Agent.Harnesses[n]
		if h.Program == "" {
			continue
		}
		from := "default"
		if h.Set {
			from = "set"
		}
		fmt.Fprintf(w, "  harness %s: %s (%s)\n", n, strings.TrimSpace(h.Program+" "+strings.Join(h.Args, " ")), from)
		if len(h.Outside) > 0 {
			fmt.Fprintf(w, "    outside: %s\n", strings.Join(h.Outside, ", "))
		}
		if n != config.ContainerHarness {
			continue
		}
		for _, inner := range slices.Sorted(maps.Keys(h.Images)) {
			fmt.Fprintf(w, "    image %s: %s\n", inner, h.Images[inner])
		}
		network := h.Network
		if network == "" {
			network = "Docker's default bridge"
		}
		login := h.Login
		if login == "" {
			login = "not set"
		}
		fmt.Fprintf(w, "    network %s, login %s\n", network, login)
	}
	for _, n := range slices.Sorted(maps.Keys(v.Agent.Providers)) {
		p := v.Agent.Providers[n]
		key := "not set"
		if p.KeySet {
			key = "set"
		}
		fmt.Fprintf(w, "  provider %s: %s, key %s %s\n", n, p.API, orNone(p.KeyEnv), key)
	}
	if d := v.DefaultAgent; d != nil && (d.Harness != "" || d.Model != "" || d.Provider != "") {
		fmt.Fprintf(w, "default agent: %s\n", strings.Join(nonEmpty(d.Harness, d.Model, d.Provider), " "))
	}
	fmt.Fprintf(w, "checks: timeout %d minutes\n", v.Checks.TimeoutMinutes)
	for _, c := range v.Checks.Commands {
		fmt.Fprintf(w, "  %s: %s\n", c.Name, strings.Join(c.Command, " "))
	}
	fmt.Fprintf(w, "import roots: %s\n", orNone(strings.Join(v.ImportRoots, ", ")))
	fmt.Fprintf(w, "allowed hosts: %s\n", orNone(strings.ReplaceAll(v.AllowedHosts, ",", ", ")))
	fmt.Fprintf(w, "shared paths: %s\n", orNone(strings.Join(v.Shared, ", ")))
	if v.Projects.Error != "" {
		fmt.Fprintf(w, "projects: %s\n", v.Projects.Error)
	} else {
		running := "not running"
		if v.Projects.Running {
			running = "running"
		}
		fmt.Fprintf(w, "projects: flai serve %s, %d served, %d not served\n", running, len(v.Projects.Served), len(v.Projects.Unserved))
	}
	if p := v.Planning; p != nil {
		replan := p.Replan
		if !p.ReplanSet {
			replan += " (default)"
		}
		schedule := orNone(p.Schedule)
		if p.Next != "" {
			schedule += ", next " + p.Next
		}
		fmt.Fprintf(w, "planning: plan %s, replan %s, schedule %s\n", onOff(p.Plan), replan, schedule)
	}
	if m := v.MCP; m != nil {
		state := "not running"
		if m.Running {
			state = fmt.Sprintf("running at %s (pid %d)", m.URL, m.PID)
		}
		bind := m.Bind + " (default; flai serve mcp bind <address> sets another)"
		if m.BindSet {
			bind = m.Bind + " (set; flai serve mcp clear goes back to the default)"
		}
		fmt.Fprintf(w, "mcp: %s; bind %s\n", state, bind)
	}
}

func onOff(on bool) string {
	if on {
		return "on"
	}
	return "off"
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

func nonEmpty(values ...string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
