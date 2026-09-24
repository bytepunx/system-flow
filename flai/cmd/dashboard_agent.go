package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The channel between the dashboard and flai on the host (ADR-0029). flai
// dashboard makes the credential the two prove to each other, hands it to
// the container as a secret, registers the project with flai serve, and
// starts flai host, which runs flai serve, when it is not running. The credential is not the login
// token: it opens nothing but the agent endpoint, and it is never sent.

const (
	agentKeyFileName  = "dashboard.agent-key"
	agentKeyMountPath = "/run/secrets/flaiover_agent_key"
	agentKeyFileEnv   = "FLAIOVER_AGENT_KEY_FILE"
)

// agentKeyPath is the shared credential's file, beside flai serve's state,
// next to flai's config file (S-0080): every project the dashboard serves
// proves itself with the one file, as the login token is now the one file.
func agentKeyPath(dir string) string { return filepath.Join(dir, agentKeyFileName) }

// ensureAgentKey creates the credential when missing, mode 0600.
func ensureAgentKey(dir string) error {
	p := agentKeyPath(dir)
	if data, err := os.ReadFile(p); err == nil && strings.TrimSpace(string(data)) != "" {
		return nil
	}
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(base64.RawURLEncoding.EncodeToString(buf)+"\n"), tokenFileMode)
}

// agentKeyArgs hand the credential to the container the way the token is.
func agentKeyArgs(dir string) []string {
	return []string{
		"--mount", "type=bind,source=" + agentKeyPath(dir) + ",target=" + agentKeyMountPath + ",readonly",
		"--env", agentKeyFileEnv + "=" + agentKeyMountPath,
	}
}

// dialURL is where flai on this host reaches the container's published port.
func (s dashboardSettings) dialURL() string {
	host := s.Bind
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return fmt.Sprintf("http://%s:%d", host, s.Port)
}

// connectServe registers the project with flai serve and makes sure flai
// host runs, which runs flai serve (S-0106). It reports what it did; a
// failure here leaves a dashboard that works as before, so it is told and not
// fatal.
func (a *app) connectServe(repo *workitem.Repo, s dashboardSettings) (note string) {
	if repo == nil {
		return a.connectServeFolder(s)
	}
	entry := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: s.Root, URL: s.dialURL(), KeyFile: agentKeyPath(string(a.serveDir()))}
	if entry.Key == "" {
		return "  host flai: not connected, the manifest has no key (flai check says how to add one)\n"
	}
	if err := a.serveDir().Register(entry); err != nil {
		return "  host flai: not registered: " + err.Error() + "\n"
	}
	st, started, err := a.ensureHost()
	switch {
	case err != nil:
		return "  host flai: registered, but flai host did not start: " + err.Error() + "\n    start it with: flai host start\n"
	case started:
		return fmt.Sprintf("  host flai: flai host started (pid %d); it runs flai serve, which will connect to this dashboard; flai host status shows it\n", st.PID)
	}
	return fmt.Sprintf("  host flai: flai host is running (pid %d); its flai serve will connect to this dashboard\n", st.PID)
}

// hostFlaiStatus is the part of flai serve's state about this project, and
// whether flai host, which runs flai serve, runs (S-0106).
type hostFlaiStatus struct {
	HostRunning bool   `json:"host_running"`
	HostPID     int    `json:"host_pid,omitempty"`
	Running     bool   `json:"running"`
	PID         int    `json:"pid,omitempty"`
	Since       string `json:"since,omitempty"`
	Registered  bool   `json:"registered"`
	Connected   bool   `json:"connected"`
	ConnSince   string `json:"connected_since,omitempty"`
	LastError   string `json:"last_error,omitempty"`
	Projects    int    `json:"projects"`
}

func (a *app) hostFlai(root string) hostFlaiStatus {
	out := hostFlaiStatus{}
	if h, alive := a.hostDir().ReadStatus(time.Now()); alive {
		out.HostRunning, out.HostPID = true, h.PID
	}
	st, err := a.readServeStatus()
	if err != nil {
		return out
	}
	out.Projects = len(st.Projects)
	for _, p := range st.Projects {
		if p.Root == root {
			out.Registered = true
		}
	}
	if st.Running {
		out.Running, out.PID, out.Since = true, st.Status.PID, st.Status.Started
		if c, ok := st.Status.Connections[root]; ok {
			out.Connected, out.ConnSince, out.LastError = c.Connected, c.Since, c.LastError
		}
	}
	return out
}

func (h hostFlaiStatus) describe() string {
	switch {
	case !h.Running && h.HostRunning:
		return fmt.Sprintf("  host flai: flai host runs (pid %d) but flai serve does not; flai host status says why\n", h.HostPID)
	case !h.Running:
		return "  host flai: flai host is not running; start it with flai host start, which runs flai serve\n"
	case !h.Registered:
		return fmt.Sprintf("  host flai: flai serve runs (pid %d, since %s, %d project(s)) but this project is not registered with it; restart the dashboard\n", h.PID, h.Since, h.Projects)
	case h.Connected:
		return fmt.Sprintf("  host flai: connected since %s (flai serve pid %d, running since %s, serving %d project(s))\n", h.ConnSince, h.PID, h.Since, h.Projects)
	case h.LastError != "":
		return fmt.Sprintf("  host flai: flai serve runs (pid %d) but is not connected: %s\n", h.PID, h.LastError)
	}
	return fmt.Sprintf("  host flai: flai serve runs (pid %d) and is connecting\n", h.PID)
}

// connectServeFolder is connectServe outside any project (S-0101): the
// projects below the folder are registered, the dashboard is recorded so that
// flai serve offers it the repositories to import even when none is, the
// folder is named for import, as flai serve import add would, and flai serve
// is started when it is not running.
func (a *app) connectServeFolder(s dashboardSettings) (note string) {
	dir := a.serveDir()
	if err := dir.AddDashboard(serve.Dashboard{URL: s.dialURL(), KeyFile: agentKeyPath(string(dir))}); err != nil {
		return "  host flai: dashboard not recorded: " + err.Error() + "\n"
	}
	named := ""
	if cfg, path, err := a.loadConfig(); err != nil {
		named = "  import: this folder not named (" + err.Error() + ")\n"
	} else if !slices.Contains(cfg.ImportRoots, s.Root) {
		cfg.ImportRoots = append(cfg.ImportRoots, s.Root)
		if err := config.Save(path, cfg); err != nil {
			named = "  import: this folder not named (" + err.Error() + ")\n"
		} else {
			named = fmt.Sprintf("  import: %s is named for import (flai serve import remove . undoes it); its git repositories without system-flow.yaml are offered on the board\n", s.Root)
		}
	} else {
		named = fmt.Sprintf("  import: %s is named for import; its git repositories without system-flow.yaml are offered on the board\n", s.Root)
	}
	// the projects below the folder are served too, as flai mcp there serves them
	var keys []string
	for _, root := range workitem.FindProjects(s.Root) {
		repo, err := workitem.Open(root)
		if err != nil || repo.Manifest.Key == "" {
			named += fmt.Sprintf("  %s: not registered (no key in its system-flow.yaml; flai check says how to add one)\n", root)
			continue
		}
		if err := dir.Register(serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: root, URL: s.dialURL(), KeyFile: agentKeyPath(string(dir))}); err != nil {
			named += fmt.Sprintf("  %s: not registered: %s\n", root, err)
			continue
		}
		keys = append(keys, repo.Manifest.Key)
	}
	registered := "no project below this folder, so none is registered"
	if len(keys) > 0 {
		registered = fmt.Sprintf("registered the %d project(s) below this folder: %s", len(keys), strings.Join(keys, ", "))
	}
	st, started, err := a.ensureHost()
	switch {
	case err != nil:
		return named + "  host flai: flai host did not start: " + err.Error() + "\n    start it with: flai host start\n"
	case started:
		return named + fmt.Sprintf("  host flai: flai host started (pid %d) and runs flai serve; %s\n", st.PID, registered)
	}
	return named + fmt.Sprintf("  host flai: flai host is running (pid %d); %s\n", st.PID, registered)
}

// describeFolder is describe outside any project (S-0101): there is nothing
// to be registered or connected here, only flai serve itself to report.
func (h hostFlaiStatus) describeFolder() string {
	if !h.Running {
		return "  host flai: no project here; flai serve is not running, start it with flai host start (or flai dashboard)\n"
	}
	return fmt.Sprintf("  host flai: no project here; flai serve runs (pid %d, since %s, serving %d project(s))\n", h.PID, h.Since, h.Projects)
}
