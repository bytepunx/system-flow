package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/execx"
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
	s := dashboardSettings{Image: cfg.Dashboard.Image, Tag: cfg.Dashboard.Tag, Port: cfg.Dashboard.Port, Bind: cfg.Dashboard.Bind, Root: repo.Root}
	// The dashboard serves the main checkout: wip lives there, and a linked
	// worktree on its own has no repository to point at (ADR-0019).
	if repo.MainRoot != "" {
		s.Root = repo.MainRoot
	}
	if d := repo.Manifest.Dashboard; d.Bind != "" {
		s.Bind = d.Bind
	}
	if bind != "" {
		s.Bind = bind
	}
	if s.Bind == "" {
		s.Bind = defaultBind
	}
	if d := repo.Manifest.Dashboard; d.Image != "" {
		s.Image = d.Image
	}
	if d := repo.Manifest.Dashboard; d.Tag != "" {
		s.Tag = d.Tag
	}
	if d := repo.Manifest.Dashboard; d.Port != 0 {
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
	name := strings.ToLower(repo.Manifest.Name)
	name = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, name)
	s.Name = "flaiover-" + strings.Trim(name, "-")
	return s, nil
}

func (s dashboardSettings) ref() string { return s.Image + ":" + s.Tag }
func (s dashboardSettings) url() string { return fmt.Sprintf("http://localhost:%d", s.Port) }

func newDashboardCmd(a *app) *cobra.Command {
	var image, tag, bind string
	var port int
	var pull, attach, open, build bool
	c := &cobra.Command{
		Use:   "dashboard",
		Short: "Run the flaiover dashboard against this project in Docker",
		Long: `Pull the flaiover image if it is missing and run it detached with this
repository mounted read-write at its own host path, so the links git keeps
for story worktrees resolve in the container, published on every interface
(--bind 127.0.0.1 to restrict) on the configured port, as the current user
so files it writes keep your ownership. Image, tag, port, and bind come from
flags, then the dashboard section of system-flow.yaml, then config.`,
		Example: `  flai dashboard
  flai dashboard --port 8080 --pull
  flai dashboard --attach          # stream logs until Ctrl-C (the container keeps running)
  flai dashboard --build           # build flaiover:local from this monorepo instead of pulling
  flai dashboard stop`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboard(image, tag, port, bind, pull, attach, open, build)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "image tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on, 0.0.0.0 for every interface (default: manifest, then config, then 0.0.0.0)")
	f.BoolVar(&pull, "pull", false, "pull the image even if present")
	f.BoolVar(&attach, "attach", false, "follow the container logs after starting")
	f.BoolVar(&open, "open", false, "open the dashboard in a browser")
	f.BoolVar(&build, "build", false, "build the image from flaiover/ in this repository as flaiover:local and run that")
	c.AddCommand(newDashboardStopCmd(a), newDashboardStatusCmd(a), newDashboardLogsCmd(a), newDashboardTokenCmd(a))
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

// gitIdentityArgs passes the host's git user into the container as the
// standard GIT_AUTHOR_* and GIT_COMMITTER_* variables, when it has one.
func (a *app) gitIdentityArgs(root string) []string {
	if _, err := a.runner.LookPath("git"); err != nil {
		return nil
	}
	name, err1 := a.runner.Run(root, "git", "config", "user.name")
	email, err2 := a.runner.Run(root, "git", "config", "user.email")
	if err1 != nil || err2 != nil || name == "" || email == "" {
		return nil
	}
	return []string{
		"--env", "GIT_AUTHOR_NAME=" + name, "--env", "GIT_AUTHOR_EMAIL=" + email,
		"--env", "GIT_COMMITTER_NAME=" + name, "--env", "GIT_COMMITTER_EMAIL=" + email,
	}
}

// excludesMountPath is where the host's global git excludes file appears in
// the container.
const excludesMountPath = "/run/flaiover/gitignore"

// gitExcludesArgs gives the container the host's global git excludes, so a
// file ignored only there is not reported as uncommitted by git in the
// container, which refused an acceptance from the board (I-0019). The file
// is the one git uses on the host: core.excludesFile when set, else git's
// default location. It is mounted read-only and named to git through
// GIT_CONFIG_* variables, like the identity, because the container has no
// home to hold a git config. Nothing is passed when there is no such file.
func (a *app) gitExcludesArgs(root string) []string {
	if _, err := a.runner.LookPath("git"); err != nil {
		return nil
	}
	file, _ := a.runner.Run(root, "git", "config", "--type=path", "--get", "core.excludesFile")
	file = strings.TrimSpace(file)
	if file == "" {
		base := os.Getenv("XDG_CONFIG_HOME")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil
			}
			base = filepath.Join(home, ".config")
		}
		file = filepath.Join(base, "git", "ignore")
	}
	if st, err := os.Stat(file); err != nil || !st.Mode().IsRegular() || strings.Contains(file, ",") {
		return nil
	}
	return []string{
		"--mount", "type=bind,source=" + file + ",target=" + excludesMountPath + ",readonly",
		"--env", "GIT_CONFIG_COUNT=1",
		"--env", "GIT_CONFIG_KEY_0=core.excludesFile",
		"--env", "GIT_CONFIG_VALUE_0=" + excludesMountPath,
	}
}

// fallbackMount is where the repository goes when its host path cannot be
// a path in the container, and the image's own default.
const fallbackMount = "/project"

// containerMount returns the path the repository is mounted at in the
// container. Git links a story worktree to its repository with absolute
// paths, so they resolve in the container only when it sees the repository
// at its host path (ADR-0022). A Windows path cannot be one in a Linux
// container; mirrored is false then and the fixed mount is used.
func containerMount(goos, root string) (target string, mirrored bool) {
	if goos == "windows" || !strings.HasPrefix(root, "/") || strings.Contains(root, ":") {
		return fallbackMount, false
	}
	return root, true
}

func (a *app) requireDocker() error {
	return execx.Require(a.runner, "docker", "Install Docker Engine 24 or newer (https://docs.docker.com/engine/install/) or Docker Desktop, and make sure the daemon is running.")
}

func (a *app) runDashboard(image, tag string, port int, bind string, pull, attach, open, build bool) error {
	repo, err := a.project()
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
	if running, _ := a.containerRunning(s.Name); running {
		_, url := a.containerInfo(s.Name)
		fmt.Fprintf(a.out, "%s is already running at %s (flai dashboard stop to stop it)\n", s.Name, orDefault(url, s.url()))
		return nil
	}
	if build {
		if err := a.buildDashboardImage(repo.Root, s.ref()); err != nil {
			return err
		}
	} else if _, err := a.runner.Run("", "docker", "image", "inspect", s.ref()); err != nil || pull {
		if err := a.pullDashboardImage(s); err != nil {
			return err
		}
	}
	token, created, err := ensureToken(repo.MainRoot)
	if err != nil {
		return fmt.Errorf("dashboard token: %w", err)
	}
	if created {
		a.logger().Info("dashboard token created", "component", "dashboard", "file", relPath(repo.MainRoot, tokenPath(repo.MainRoot)))
	}
	mount, mirrored := containerMount(runtime.GOOS, s.Root)
	if !mirrored {
		a.logger().Warn("repository cannot be mounted at its host path", "component", "dashboard", "root", s.Root, "mount", mount,
			"effect", "stories with a branch cannot be accepted from the board",
			"fix", "accept from a shell with flai accept, or set worktrees.relative_paths (git 2.48 or newer) before opening stories")
	}
	args := []string{"run", "--detach", "--rm", "--name", s.Name,
		"--publish", fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort),
		"--volume", s.Root + ":" + mount,
		"--env", "PROJECT_DIR=" + mount,
	}
	args = append(args, tokenArgs(repo.MainRoot)...)
	// Acceptance from the dashboard commits as the person who started it
	// (S-0046): the container has no ~/.gitconfig of its own.
	args = append(args, a.gitIdentityArgs(repo.MainRoot)...)
	args = append(args, a.gitExcludesArgs(repo.MainRoot)...)
	if runtime.GOOS != "windows" {
		args = append(args, "--user", strconv.Itoa(os.Getuid())+":"+strconv.Itoa(os.Getgid()))
	}
	args = append(args, s.ref())
	id, err := a.runner.Run("", "docker", args...)
	if err != nil {
		return err
	}
	a.logger().Info("dashboard started", "component", "dashboard", "container", s.Name, "id", short(id), "url", s.url(), "bind", s.Bind)
	if a.jsonOut {
		return a.printJSON(map[string]any{"container": s.Name, "id": short(id), "image": s.ref(), "url": s.url(), "login_url": loginURL(s.url(), token), "bind": s.Bind, "port": s.Port, "mount": s.Root, "mount_target": mount})
	}
	reach := "reachable from this host only"
	if s.Bind == defaultBind {
		reach = "reachable on every interface of this host; the token is required, keep the host private"
	}
	fmt.Fprintf(a.out, "flaiover running at %s (%s)\n  log in with: %s\n  container %s, image %s, %s mounted read-write at %s\n  token: %s (flai dashboard token to print or rotate)\n  stop with: flai dashboard stop\n", s.url(), reach, loginURL(s.url(), token), s.Name, s.ref(), s.Root, mount, relPath(repo.MainRoot, tokenPath(repo.MainRoot)))
	if open {
		openBrowser(loginURL(s.url(), token))
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
		Short: "Stop the dashboard container for this project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
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
			if running, _ := a.containerRunning(s.Name); !running {
				fmt.Fprintf(a.out, "%s is not running\n", s.Name)
				return nil
			}
			if _, err := a.runner.Run("", "docker", "stop", s.Name); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"container": s.Name, "state": "stopped"})
			}
			fmt.Fprintf(a.out, "stopped %s\n", s.Name)
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
			repo, err := a.project()
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
			if a.jsonOut {
				return a.printJSON(map[string]any{"container": s.Name, "running": running, "url": url, "image": image})
			}
			if running {
				fmt.Fprintf(a.out, "%s running at %s (%s)\n", s.Name, url, image)
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
			repo, err := a.project()
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
