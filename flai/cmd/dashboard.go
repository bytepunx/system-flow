package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// containerPort is what flaiover listens on inside the image.
const containerPort = 3000

type dashboardSettings struct {
	Image, Tag string
	Port       int
	Bind       string // host address the port is published on
	Name       string
	Root       string
}

// defaultBind publishes on every interface so the dashboard is reachable
// from other hosts (operator's call, S-0035); --bind 127.0.0.1 restricts it.
const defaultBind = "0.0.0.0"

// dashboardSettings resolves image, tag, and port: flags, then the
// manifest's dashboard section, then config.
func (a *app) dashboardSettings(repo *workitem.Repo, image, tag string, port int, bind string) (dashboardSettings, error) {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return dashboardSettings{}, err
	}
	s := dashboardSettings{Image: cfg.Dashboard.Image, Tag: cfg.Dashboard.Tag, Port: cfg.Dashboard.Port, Bind: cfg.Dashboard.Bind}
	// Outside any project (S-0101) the configuration alone decides, and Root
	// is the folder it was run in.
	var m manifest.Manifest
	if repo == nil {
		if s.Root, err = a.workingDir(); err != nil {
			return dashboardSettings{}, err
		}
	} else {
		m, s.Root = repo.Manifest, repo.Root
		// The dashboard serves the main checkout: wip lives there, and a linked
		// worktree on its own has no repository to point at (ADR-0019).
		if repo.MainRoot != "" {
			s.Root = repo.MainRoot
		}
	}
	if d := m.Dashboard; d.Bind != "" {
		s.Bind = d.Bind
	}
	if bind != "" {
		s.Bind = bind
	}
	if s.Bind == "" {
		s.Bind = defaultBind
	}
	if d := m.Dashboard; d.Image != "" {
		s.Image = d.Image
	}
	if d := m.Dashboard; d.Tag != "" {
		s.Tag = d.Tag
	}
	if d := m.Dashboard; d.Port != 0 {
		s.Port = d.Port
	}
	if image != "" {
		s.Image = image
	}
	if tag != "" {
		s.Tag = tag
	}
	if port != 0 {
		s.Port = port
	}
	if s.Image == "" || s.Port == 0 {
		return dashboardSettings{}, fmt.Errorf("dashboard image and port must be set (config dashboard.*, system-flow.yaml dashboard, or flags)")
	}
	if s.Tag == "" {
		s.Tag = "latest"
	}
	// One container serves every project the host flai serves (S-0080): its
	// name is fixed, not derived from a project's own name any more.
	s.Name = sharedContainerName
	return s, nil
}

// sharedContainerName is the one dashboard container's name, for every
// project on this host.
const sharedContainerName = "flaiover"

func (s dashboardSettings) ref() string { return s.Image + ":" + s.Tag }
func (s dashboardSettings) url() string { return fmt.Sprintf("http://localhost:%d", s.Port) }

func newDashboardCmd(a *app) *cobra.Command {
	var image, tag, bind, pushKeyFlag, pushHostsFlag string
	var port int
	var pull, attach, open, build, noServe bool
	c := &cobra.Command{
		Use:   "dashboard",
		Short: "Make sure the one flaiover dashboard runs and serves this project",
		Long: `One dashboard container serves every project on this host (S-0080): if it is
not running, this pulls the image (unless missing) and starts it, published on
every interface (--bind 127.0.0.1 to restrict) on the configured port; if it
already runs, for this project or another, this registers the project with it
and starts no second container. The container is given that port, the one
login token, and the one credential flai serve proves itself with, kept
beside flai serve's state, and nothing else: no file of any project is
mounted into it (ADR-0031). Everything it shows and changes it asks of flai
serve on this host, which this command registers the project with and
starts. Image, tag, port, and bind come from flags, then the dashboard
section of system-flow.yaml, then config; they matter only the first time,
when they decide what the shared container is started with.`,
		Example: `  flai dashboard
  flai dashboard --port 8080 --pull
  flai dashboard --attach          # stream logs until Ctrl-C (the container keeps running)
  flai dashboard --build           # build flaiover:local from this monorepo instead of pulling
  flai dashboard stop`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboard(image, tag, port, bind, pushKeyFlag, pushHostsFlag, pull, attach, open, build, noServe)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "image tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on, 0.0.0.0 for every interface (default: manifest, then config, then 0.0.0.0)")
	f.StringVar(&pushKeyFlag, "push-key", "", "retired (ADR-0031): the container pushes nothing")
	f.StringVar(&pushHostsFlag, "push-known-hosts", "", "retired (ADR-0031)")
	_ = f.MarkHidden("push-key")
	_ = f.MarkHidden("push-known-hosts")
	f.BoolVar(&pull, "pull", false, "pull the image even if present")
	f.BoolVar(&attach, "attach", false, "follow the container logs after starting")
	f.BoolVar(&open, "open", false, "open the dashboard in a browser")
	f.BoolVar(&noServe, "no-serve", false, "do not register the project with flai serve or start it; the dashboard then has no flai on the host to ask (ADR-0029)")
	f.BoolVar(&build, "build", false, "build the image from flaiover/ in this repository as flaiover:local and run that")
	c.AddCommand(newDashboardStopCmd(a), newDashboardStatusCmd(a), newDashboardLogsCmd(a), newDashboardTokenCmd(a),
		newDashboardRestartCmd(a), newDashboardCheckCmd(a), newDashboardUpgradeCmd(a))
	return c
}

const (
	localImage = "flaiover"
	localTag   = "local"
)

// pullDashboardImage pulls the image; when the registry refuses, it logs
// Docker in with the GitHub token sources and tries once more. The flaiover
// package on GHCR is private while the repository is.
func (a *app) pullDashboardImage(s dashboardSettings) error {
	a.logger().Info("pulling image", "component", "dashboard", "image", s.ref())
	_, err := a.runner.Run("", "docker", "pull", "--quiet", s.ref())
	if err == nil {
		return nil
	}
	registry := registryOf(s.Image)
	if registry == "" || !unauthorized(err) {
		return err
	}
	token := a.githubToken()
	if token == "" {
		return fmt.Errorf("%s refused the pull for %s and no token is available: set GITHUB_TOKEN or run `gh auth login`, then `gh auth refresh -h github.com -s read:packages` (the package is private while the repository is; making the package public also works)", registry, s.ref())
	}
	user := a.registryUser()
	a.logger().Info("logging into registry", "component", "dashboard", "registry", registry, "user", user)
	if _, lerr := a.runner.RunInput("", "docker", token+"\n", "login", registry, "--username", user, "--password-stdin"); lerr != nil {
		return fmt.Errorf("docker login %s failed: %w; the token needs the read:packages scope (`gh auth refresh -h github.com -s read:packages`) or the package must be public", registry, lerr)
	}
	if _, err := a.runner.Run("", "docker", "pull", "--quiet", s.ref()); err != nil {
		return fmt.Errorf("%w; logged into %s as %s but the pull was still refused: the token needs read:packages, or the package must be public", err, registry, user)
	}
	return nil
}

// buildDashboardImage builds flaiover:local from flaiover/ in a monorepo
// checkout, with the same version metadata as scripts/flaiover-image.sh.
func (a *app) buildDashboardImage(root, ref string) error {
	dockerfile := filepath.Join(root, "flaiover", "Dockerfile")
	if _, err := os.Stat(dockerfile); err != nil {
		return fmt.Errorf("--build needs flaiover/Dockerfile under %s (the system-flow monorepo); elsewhere pull the published image", root)
	}
	describe := func(match, prefix, def string) string {
		out, err := a.runner.Run(root, "git", "describe", "--tags", "--match", match, "--abbrev=0")
		if err != nil || out == "" {
			return def
		}
		return strings.TrimPrefix(out, prefix)
	}
	commit, err := a.runner.Run(root, "git", "rev-parse", "--short", "HEAD")
	if err != nil || commit == "" {
		commit = "unknown"
	}
	a.logger().Info("building image", "component", "dashboard", "image", ref, "dockerfile", dockerfile)
	_, err = a.runner.Run(root, "docker", "build", "-f", dockerfile, "-t", ref,
		"--build-arg", "FLAI_VERSION="+describe("flai/v*", "flai/v", "dev"),
		"--build-arg", "FLAI_COMMIT="+commit,
		"--build-arg", "FLAIOVER_VERSION="+describe("flaiover/v*", "flaiover/v", "0.0.0"),
		"--build-arg", "FLAI_DATE="+a.now().UTC().Format(time.RFC3339),
		root)
	return err
}

// registryOf returns the registry host of an image reference, or "" for
// Docker Hub.
func registryOf(image string) string {
	first, _, ok := strings.Cut(image, "/")
	if !ok || (!strings.Contains(first, ".") && !strings.Contains(first, ":") && first != "localhost") {
		return ""
	}
	return first
}

func unauthorized(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unauthorized") || strings.Contains(msg, "denied") || strings.Contains(msg, "authentication required")
}

// registryUser is the login name for docker login: the gh account, then
// GITHUB_ACTOR, then a placeholder (GHCR authenticates the token, not the name).
func (a *app) registryUser() string {
	if _, err := a.runner.LookPath("gh"); err == nil {
		if out, err := a.runner.Run("", "gh", "api", "user", "--jq", ".login"); err == nil && out != "" {
			return out
		}
	}
	if v := os.Getenv("GITHUB_ACTOR"); v != "" {
		return v
	}
	return "token"
}

// retiredPushKey says what replaced the push key when one is still asked
// for, by flag or in the host's config. It is a note, not an error: the
// dashboard starts, and what the setting was for is done another way.
func (a *app) retiredPushKey(keyFlag, hostsFlag string) string {
	set := []string{}
	if keyFlag != "" {
		set = append(set, "--push-key")
	}
	if hostsFlag != "" {
		set = append(set, "--push-known-hosts")
	}
	if cfg, _, err := a.loadConfig(); err == nil {
		if cfg.Dashboard.PushKey != "" {
			set = append(set, "dashboard.push_key")
		}
		if cfg.Dashboard.PushKnownHosts != "" {
			set = append(set, "dashboard.push_known_hosts")
		}
	}
	if len(set) == 0 {
		return ""
	}
	return "  " + strings.Join(set, ", ") + ": retired and ignored. The container holds no git, no ssh, and no file of the project, so it has nothing to push with (ADR-0031). An acceptance from the board is pushed from this host with your own credentials: flai push --pending. Clear the setting with: flai config set dashboard.push_key \"\"\n"
}

// secretMountPrefix is where the container's two secrets are mounted; a
// running container with a mount anywhere else was started by an older flai.
const secretMountPrefix = "/run/secrets/"

// staleMounts lists what a running container has mounted besides its two
// secrets: the project, from a flai older than ADR-0031.
func (a *app) staleMounts(container string) []string {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{range .Mounts}}{{.Destination}}{{\"\\n\"}}{{end}}", container)
	if err != nil {
		return nil
	}
	var stale []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" && !strings.HasPrefix(l, secretMountPrefix) {
			stale = append(stale, l)
		}
	}
	return stale
}

func (a *app) requireDocker() error {
	return execx.Require(a.runner, "docker", "Install Docker Engine 24 or newer (https://docs.docker.com/engine/install/) or Docker Desktop, and make sure the daemon is running.")
}

func (a *app) runDashboard(image, tag string, port int, bind, pushKeyFlag, pushHostsFlag string, pull, attach, open, build, noServe bool) error {
	repo, err := a.projectOrNone()
	if err != nil {
		return err
	}
	if err := a.requireDocker(); err != nil {
		return err
	}
	if build {
		image, tag = localImage, localTag
	}
	s, err := a.dashboardSettings(repo, image, tag, port, bind)
	if err != nil {
		return err
	}
	dir := string(a.serveDir())
	retired := a.retiredPushKey(pushKeyFlag, pushHostsFlag)
	// A container started for another project may already serve this one
	// (S-0080): register this project with it too, and say so, instead of
	// starting a second container or refusing.
	already, _ := a.containerRunning(s.Name)
	if !already {
		if build {
			if err := a.buildDashboardImage(s.Root, s.ref()); err != nil {
				return err
			}
		} else if _, err := a.runner.Run("", "docker", "image", "inspect", s.ref()); err != nil || pull {
			if err := a.pullDashboardImage(s); err != nil {
				return err
			}
		}
	}
	token, created, err := ensureToken(dir)
	if err != nil {
		return fmt.Errorf("dashboard token: %w", err)
	}
	if created {
		a.logger().Info("dashboard token created", "component", "dashboard", "file", tokenPath(dir))
	}
	// The credential flai serve and the dashboard prove to each other (ADR-0029).
	if err := ensureAgentKey(dir); err != nil {
		return fmt.Errorf("agent credential: %w", err)
	}
	var id string
	if !already {
		id, err = a.startContainer(dir, s.Name, fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort), s.ref())
		if err != nil {
			return err
		}
		a.logger().Info("dashboard started", "component", "dashboard", "container", s.Name, "id", short(id), "url", s.url(), "bind", s.Bind)
	}
	serveNote := "  host flai: not started (--no-serve); flai serve start connects it\n"
	if !noServe {
		serveNote = a.connectServe(repo, s)
	}
	url := s.url()
	if already {
		if _, u := a.containerInfo(s.Name); u != "" {
			url = u
		}
	}
	if a.jsonOut {
		out := map[string]any{"container": s.Name, "already_running": already, "id": short(id), "url": url, "login_url": loginURL(url, token), "bind": s.Bind, "project": s.Root}
		out["host_flai"] = a.hostFlai(s.Root)
		if retired != "" {
			out["retired"] = strings.TrimSpace(retired)
		}
		return a.printJSON(out)
	}
	reach := "reachable from this host only"
	if s.Bind == defaultBind {
		reach = "reachable on every interface of this host; the token is required, keep the host private"
	}
	if already {
		fmt.Fprintf(a.out, "%s already runs at %s, and now also serves %s\n  log in with: %s\n  token: %s (flai dashboard token to print or rotate)\n  stop with: flai dashboard stop\n", s.Name, url, s.Root, loginURL(url, token), tokenPath(dir))
	} else {
		fmt.Fprintf(a.out, "flaiover running at %s (%s)\n  log in with: %s\n  container %s, image %s, for %s\n    it is given this port, its login token, and the host flai's credential, and no file of any project (ADR-0031)\n  token: %s (flai dashboard token to print or rotate)\n  stop with: flai dashboard stop\n", url, reach, loginURL(url, token), s.Name, s.ref(), s.Root, tokenPath(dir))
	}
	fmt.Fprint(a.out, retired)
	fmt.Fprint(a.out, serveNote)
	if open {
		openBrowser(loginURL(url, token))
	}
	if attach {
		cmd := exec.Command("docker", "logs", "--follow", s.Name)
		cmd.Stdout, cmd.Stderr = a.out, a.errOut
		return cmd.Run()
	}
	return nil
}

func (a *app) containerRunning(name string) (bool, error) {
	out, err := a.runner.Run("", "docker", "ps", "--quiet", "--filter", "name=^/"+name+"$")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// containerInfo reads the running container's image and host port, so
// status reflects what is actually running rather than the configuration.
func (a *app) containerInfo(name string) (image, url string) {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{.Config.Image}} {{range $p, $b := .HostConfig.PortBindings}}{{(index $b 0).HostPort}}{{end}}", name)
	if err != nil {
		return "", ""
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) >= 1 {
		image = parts[0]
	}
	if len(parts) >= 2 {
		url = "http://localhost:" + parts[1]
	}
	return image, url
}

func newDashboardStopCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop this project's dashboard; the shared container stops only when it was the last project served",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.projectOrNone()
			if err != nil {
				return err
			}
			if err := a.requireDocker(); err != nil {
				return err
			}
			s, err := a.dashboardSettings(repo, "", "", 0, "")
			if err != nil {
				return err
			}
			// flai serve stays for the other projects; it stops dialling this one.
			// One dashboard container serves every project on this host (S-0080),
			// so it stops only when this was the last one registered.
			if err := a.serveDir().Unregister(s.Root); err != nil {
				a.logger().Warn("project not unregistered from flai serve", "component", "dashboard", "err", err.Error())
			}
			remaining, err := a.serveDir().Projects()
			if err != nil {
				return err
			}
			if running, _ := a.containerRunning(s.Name); !running {
				if a.jsonOut {
					return a.printJSON(map[string]any{"container": s.Name, "state": "not-running"})
				}
				fmt.Fprintf(a.out, "%s is not running\n", s.Name)
				return nil
			}
			if len(remaining) > 0 {
				names := make([]string, len(remaining))
				for i, p := range remaining {
					names[i] = p.Key
				}
				if a.jsonOut {
					return a.printJSON(map[string]any{"container": s.Name, "state": "still-running", "serves": names})
				}
				fmt.Fprintf(a.out, "%s unregistered; %s keeps running for %s\n", s.Root, s.Name, strings.Join(names, ", "))
				return nil
			}
			if _, err := a.runner.Run("", "docker", "stop", s.Name); err != nil {
				return err
			}
			_ = a.serveDir().ForgetDashboard(s.dialURL())
			if a.jsonOut {
				return a.printJSON(map[string]string{"container": s.Name, "state": "stopped"})
			}
			fmt.Fprintf(a.out, "stopped %s (no project left registered)\n", s.Name)
			return nil
		},
	}
}

func newDashboardStatusCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether the dashboard container is running",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.projectOrNone()
			if err != nil {
				return err
			}
			if err := a.requireDocker(); err != nil {
				return err
			}
			s, err := a.dashboardSettings(repo, "", "", 0, "")
			if err != nil {
				return err
			}
			running, err := a.containerRunning(s.Name)
			if err != nil {
				return err
			}
			image, url := s.ref(), s.url()
			if running {
				if i, u := a.containerInfo(s.Name); i != "" {
					image, url = i, orDefault(u, url)
				}
			}
			var stale []string
			if running {
				stale = a.staleMounts(s.Name)
			}
			served, err := a.serveDir().Projects()
			if err != nil {
				return err
			}
			names := make([]string, len(served))
			for i, p := range served {
				names[i] = p.Key
			}
			if a.jsonOut {
				out := map[string]any{"container": s.Name, "running": running, "url": url, "image": image, "serves": names}
				if len(stale) > 0 {
					out["stale_mounts"] = stale
				}
				out["host_flai"] = a.hostFlai(s.Root)
				return a.printJSON(out)
			}
			if running {
				fmt.Fprintf(a.out, "%s running at %s (%s)\n", s.Name, url, image)
				if len(names) > 0 {
					fmt.Fprintf(a.out, "  serves: %s\n", strings.Join(names, ", "))
				}
				if len(stale) > 0 {
					// the shortest is the project itself; the rest lie inside it
					first := stale[0]
					for _, m := range stale {
						if len(m) < len(first) {
							first = m
						}
					}
					more := ""
					if len(stale) > 1 {
						more = fmt.Sprintf(" and %d more", len(stale)-1)
					}
					fmt.Fprintf(a.out, "  this container was started by an older flai and still has a project mounted (%s%s); restart it: flai dashboard stop, then flai dashboard\n", first, more)
				}
				if repo == nil {
					fmt.Fprint(a.out, a.hostFlai(s.Root).describeFolder())
				} else {
					fmt.Fprint(a.out, a.hostFlai(s.Root).describe())
				}
			} else {
				fmt.Fprintf(a.out, "%s not running; start with flai dashboard\n", s.Name)
			}
			return nil
		},
	}
}

func newDashboardLogsCmd(a *app) *cobra.Command {
	var follow bool
	c := &cobra.Command{
		Use:   "logs",
		Short: "Print the dashboard container's logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.projectOrNone()
			if err != nil {
				return err
			}
			if err := a.requireDocker(); err != nil {
				return err
			}
			s, err := a.dashboardSettings(repo, "", "", 0, "")
			if err != nil {
				return err
			}
			if follow {
				c := exec.Command("docker", "logs", "--follow", s.Name)
				c.Stdout, c.Stderr = a.out, a.errOut
				return c.Run()
			}
			out, err := a.runner.Run("", "docker", "logs", s.Name)
			if err != nil {
				return err
			}
			fmt.Fprintln(a.out, out)
			return nil
		},
	}
	c.Flags().BoolVarP(&follow, "follow", "f", false, "follow the log output")
	return c
}

func short(id string) string {
	id = strings.TrimSpace(id)
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func openBrowser(url string) {
	for _, candidate := range [][]string{{"xdg-open", url}, {"open", url}, {"cmd", "/c", "start", url}} {
		if _, err := exec.LookPath(candidate[0]); err == nil {
			_ = exec.Command(candidate[0], candidate[1:]...).Start()
			return
		}
	}
}
