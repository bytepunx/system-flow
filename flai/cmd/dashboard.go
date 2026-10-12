package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
	"github.com/bytepunx/system-flow/flai/internal/serve"
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

// signedDashboardImage is the image repository the flaiover releases' signed
// digest lists name (ADR-0070), and the configured image by default.
const signedDashboardImage = "ghcr.io/bytepunx/flaiover"

// dashboardReleases are the release client's options, the token aside, that
// flai dashboard resolves dashboard.tag with: the flaiover releases of the
// repository self-upgrade reads, from FLAI_RELEASES_API when set, verified
// with the keys this flai trusts. Tests point it at a stand-in.
var dashboardReleases = func() selfupgrade.Options {
	return selfupgrade.Options{Repo: "bytepunx/system-flow", APIBase: strings.TrimSpace(os.Getenv("FLAI_RELEASES_API")), Keys: releaseKeys}
}

// dashboardImage is what a new dashboard container runs: a flaiover release's
// image by digest from its verified digest list, or, for --build and an image
// of another repository, the reference as given, unsigned.
type dashboardImage struct {
	Ref         string         // what is pulled and run: <repository>@sha256:… for a release, image:tag otherwise
	Release     string         // the flaiover release whose verified list names Ref; "" when unsigned
	List        dashboard.List // that verified list, which check and upgrade match a running image against; empty when unsigned
	Unreachable error          // why the release could not be reached, when its cached list was used
}

// runs reports whether one of repoDigests, a running image's RepoDigests, is
// on d's verified list: never, for an unsigned d.
func (d dashboardImage) runs(repoDigests []string) bool {
	return d.Release != "" && slices.ContainsFunc(repoDigests, d.List.Matches)
}

// describe is the line flai dashboard prints for where d comes from.
func (d dashboardImage) describe() string {
	switch {
	case d.Release == "":
		return fmt.Sprintf("  unsigned: no signed digest list names %s; it runs as given (--build, or an image of a repository other than %s)\n", d.Ref, signedDashboardImage)
	case d.Unreachable != nil:
		return fmt.Sprintf("  flaiover %s, run by digest from its verified digest list cached beside flai serve's state: the release could not be reached (%v)\n", d.Release, d.Unreachable)
	}
	return fmt.Sprintf("  flaiover %s, run by digest from its verified digest list (ADR-0070)\n", d.Release)
}

// dashboardCache is the cache of verified digest lists beside flai serve's
// state, trusting the keys a list must verify against.
func (a *app) dashboardCache(opt selfupgrade.Options) (dashboard.Cache, error) {
	keys, err := selfupgrade.TrustedKeys(opt)
	if err != nil {
		return dashboard.Cache{}, fmt.Errorf("reading the release keys to verify a flaiover digest list failed: %w; this flai build is broken, so install a release of flai", err)
	}
	return dashboard.NewCache(string(a.serveDir()), keys), nil
}

// chooseDashboardImage is what a new container runs (ADR-0070): for the
// signed repository, s's tag resolved to a flaiover release and run by the
// digest its verified list names, a release refused or not resolved being an
// error; for --build or another repository, s's reference as given.
func (a *app) chooseDashboardImage(ctx context.Context, s dashboardSettings, build bool) (dashboardImage, error) {
	if build || s.Image != signedDashboardImage {
		return dashboardImage{Ref: s.ref()}, nil
	}
	opt := dashboardReleases()
	opt.Token = a.githubToken()
	cache, err := a.dashboardCache(opt)
	if err != nil {
		return dashboardImage{}, err
	}
	res, err := dashboard.Resolver{Releases: opt, Cache: cache}.Resolve(ctx, s.Image, s.Tag)
	var refused *dashboard.Refusal
	switch {
	case errors.As(err, &refused):
		return dashboardImage{}, err
	case err != nil:
		return dashboardImage{}, fmt.Errorf("%w; nothing was pulled, started, or swapped", err)
	}
	if res.Cached {
		a.logger().Warn("flaiover release not reached, cached digest list used", "component", "dashboard", "release", res.Version, "err", res.Unreachable.Error())
	}
	a.logger().Info("dashboard release resolved", "component", "dashboard", "tag", s.Tag, "release", res.Version, "image", res.Reference)
	return dashboardImage{Ref: res.Reference, Release: res.Version, List: res.List, Unreachable: res.Unreachable}, nil
}

// pinnedRecordRef is what a restart runs for a recorded reference: ref
// itself, unless it names the signed repository by tag, as a record an older
// flai wrote does; then the release that tag resolves to, by digest, so that
// no restart runs the signed repository unverified (ADR-0070, ADR-0118).
func (a *app) pinnedRecordRef(ctx context.Context, ref string) (string, error) {
	if strings.Contains(ref, "@") {
		return ref, nil
	}
	repo, tag := ref, "latest"
	if i := strings.LastIndex(ref, ":"); i >= 0 && !strings.Contains(ref[i+1:], "/") {
		repo, tag = ref[:i], ref[i+1:]
	}
	if repo != signedDashboardImage {
		return ref, nil
	}
	img, err := a.chooseDashboardImage(ctx, dashboardSettings{Image: repo, Tag: tag}, false)
	if err != nil {
		return "", err
	}
	return img.Ref, nil
}

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
serve on this host, which this command registers the project with; it
starts flai host, which runs flai serve and the MCP servers, when none runs
(S-0106). Image, tag, port, and bind come from flags, then the dashboard
section of system-flow.yaml, then config; they matter only the first time,
when they decide what the shared container is started with.

For ghcr.io/bytepunx/flaiover, the default image, the tag names a flaiover
release: latest is the newest, X.Y.Z that one. flai downloads the release's
digest list, verifies its signature with the key built into it, and pulls and
runs ghcr.io/bytepunx/flaiover@sha256:<digest> (ADR-0070). Any other tag, or a
release with no list this flai can verify, is refused, saying why, and nothing
starts. When the releases cannot be reached, a list verified before and kept
beside flai serve's state is used, and that is said. --build, and --image
naming another repository, run image:tag as given and are reported unsigned.`,
		Example: `  flai dashboard
  flai dashboard --tag 0.42.3      # that flaiover release, run by digest
  flai dashboard --port 8080 --pull
  flai dashboard --attach          # stream logs until Ctrl-C (the container keeps running)
  flai dashboard --build           # build flaiover:local from this monorepo instead of pulling
  flai dashboard stop`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboard(cmd.Context(), image, tag, port, bind, pushKeyFlag, pushHostsFlag, pull, attach, open, build, noServe)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image repository; one other than ghcr.io/bytepunx/flaiover runs its tag as given, unsigned (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "latest or a flaiover release X.Y.Z, run by digest from its signed list; for another --image, its tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on, 0.0.0.0 for every interface (default: manifest, then config, then 0.0.0.0)")
	f.StringVar(&pushKeyFlag, "push-key", "", "retired (ADR-0031): the container pushes nothing")
	f.StringVar(&pushHostsFlag, "push-known-hosts", "", "retired (ADR-0031)")
	_ = f.MarkHidden("push-key")
	_ = f.MarkHidden("push-known-hosts")
	f.BoolVar(&pull, "pull", false, "pull the image even if present")
	f.BoolVar(&attach, "attach", false, "follow the container logs after starting")
	f.BoolVar(&open, "open", false, "open the dashboard in a browser")
	f.BoolVar(&noServe, "no-serve", false, "do not register the project with flai serve or start flai host; the dashboard then has no flai on the host to ask (ADR-0029)")
	f.BoolVar(&build, "build", false, "build the image from flaiover/ in this repository as flaiover:local and run that")
	c.AddCommand(newDashboardStopCmd(a), newDashboardStatusCmd(a), newDashboardLogsCmd(a), newDashboardTokenCmd(a),
		newDashboardRestartCmd(a), newDashboardCheckCmd(a), newDashboardUpgradeCmd(a), newDashboardVersionsCmd(a))
	return c
}

const (
	localImage = "flaiover"
	localTag   = "local"
)

// pullDashboardRef pulls ref, an image of repository image by tag or by
// digest, for the daemon's platform; when the registry refuses, it logs
// Docker in with the GitHub token sources and tries once more. The flaiover
// package on GHCR is private while the repository is.
func (a *app) pullDashboardRef(image, ref string) error {
	platform := a.dockerPlatform()
	a.logger().Info("pulling image", "component", "dashboard", "image", ref, "platform", platform)
	pull := func() error {
		args := []string{"pull", "--quiet"}
		if platform != "" {
			args = append(args, "--platform", platform)
		}
		_, err := a.runner.Run("", "docker", append(args, ref)...)
		if err != nil && noManifest(err) {
			return fmt.Errorf("%s has no image for %s: %w; pick a tag published for it, or build one from the system-flow monorepo with flai dashboard --build", ref, orDefault(platform, "this Docker daemon's platform"), err)
		}
		return err
	}
	err := pull()
	if err == nil {
		return nil
	}
	registry := registryOf(image)
	if registry == "" || !unauthorized(err) {
		return err
	}
	token := a.githubToken()
	if token == "" {
		return fmt.Errorf("%s refused the pull for %s and no token is available: set GITHUB_TOKEN or run `gh auth login`, then `gh auth refresh -h github.com -s read:packages` (the package is private while the repository is; making the package public also works)", registry, ref)
	}
	user := a.registryUser()
	a.logger().Info("logging into registry", "component", "dashboard", "registry", registry, "user", user)
	if _, lerr := a.runner.RunInput("", "docker", token+"\n", "login", registry, "--username", user, "--password-stdin"); lerr != nil {
		return fmt.Errorf("docker login %s failed: %w; the token needs the read:packages scope (`gh auth refresh -h github.com -s read:packages`) or the package must be public", registry, lerr)
	}
	if err := pull(); err != nil {
		if noManifest(err) {
			return err
		}
		return fmt.Errorf("%w; logged into %s as %s but the pull was still refused: the token needs read:packages, or the package must be public", err, registry, user)
	}
	return nil
}

// dockerPlatform is the daemon's platform, os/arch, such as linux/arm64, or
// "" when the daemon does not say. It is the daemon's and not flai's own
// GOOS/GOARCH: an amd64 flai under Rosetta talks to an arm64 daemon.
func (a *app) dockerPlatform() string {
	out, err := a.runner.Run("", "docker", "version", "--format", "{{.Server.Os}}/{{.Server.Arch}}")
	out = strings.TrimSpace(out)
	if err != nil || strings.HasPrefix(out, "/") || strings.HasSuffix(out, "/") {
		return ""
	}
	return out
}

// imagePlatform is the os/arch a local image was built for, or "" when it
// is not present.
func (a *app) imagePlatform(ref string) string {
	out, err := a.runner.Run("", "docker", "image", "inspect", ref, "--format", "{{.Os}}/{{.Architecture}}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// noManifest reports a pull the registry answered with no image for the
// platform asked for.
func noManifest(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "no matching manifest")
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
	return "  " + strings.Join(set, ", ") + ": retired and ignored. The container holds no git, no ssh, and no file of the project, so it has nothing to push with (ADR-0031). Publishing from the board (Publish) runs on this host with your own credentials, as flai release --pending does in a shell (ADR-0067). Clear the setting with: flai config set dashboard.push_key \"\"\n"
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

// needsPull says whether ref must be pulled before it is run: it is not
// present, or it was built for a platform other than the daemon's, as an
// amd64 image pulled on an arm64 host before an arm64 one was published.
func (a *app) needsPull(ref string) bool {
	have := a.imagePlatform(ref)
	if have == "" {
		return true
	}
	want := a.dockerPlatform()
	if want != "" && have != want {
		a.logger().Info("image is for another platform", "component", "dashboard", "image", ref, "have", have, "want", want)
		return true
	}
	return false
}

func (a *app) requireDocker() error {
	return execx.Require(a.runner, "docker", "Install Docker Engine 24 or newer (https://docs.docker.com/engine/install/) or Docker Desktop, and make sure the daemon is running.")
}

func (a *app) runDashboard(ctx context.Context, image, tag string, port int, bind, pushKeyFlag, pushHostsFlag string, pull, attach, open, build, noServe bool) error {
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
	var img dashboardImage
	if !already {
		if img, err = a.chooseDashboardImage(ctx, s, build); err != nil {
			return err
		}
		if build {
			if err := a.buildDashboardImage(s.Root, img.Ref); err != nil {
				return err
			}
		} else if pull || a.needsPull(img.Ref) {
			if err := a.pullDashboardRef(s.Image, img.Ref); err != nil {
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
	if already {
		// a container an older flai started has no record: keep it as it runs,
		// by its image ID too (I-0116), so that flai host watches it (S-0184)
		if _, ok := readDashboardRecord(dir); !ok {
			if ref, _ := a.containerInfo(s.Name); ref != "" {
				a.recordDashboard(dir, dashboardRecord{Name: s.Name, Ref: ref, Image: a.runningImageID(s.Name, ref), Publish: fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort)})
			}
		}
	} else {
		id, err = a.startDashboard(dir, s, img.Ref)
		if err != nil {
			return err
		}
		a.logger().Info("dashboard started", "component", "dashboard", "container", s.Name, "id", short(id), "url", s.url(), "bind", s.Bind)
	}
	serveNote := "  host flai: not started (--no-serve); flai host start runs flai serve, which connects it\n"
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
		if !already {
			out["image"], out["signed"] = img.Ref, img.Release != ""
			if img.Release != "" {
				out["release"] = img.Release
			}
			if img.Unreachable != nil {
				out["cached_list"] = true
			}
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
		fmt.Fprintf(a.out, "flaiover running at %s (%s)\n  log in with: %s\n  container %s, image %s, for %s\n    it is given this port, its login token, and the host flai's credential, and no file of any project (ADR-0031)\n  token: %s (flai dashboard token to print or rotate)\n  stop with: flai dashboard stop\n", url, reach, loginURL(url, token), s.Name, img.Ref, s.Root, tokenPath(dir))
		fmt.Fprint(a.out, img.describe())
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

// containerInfoFormat is the docker inspect template containerInfo reads: the
// image docker ran, the host port, and the dashboardRefLabel, separated by |
// (which no reference or port holds) so that an empty field, such as the
// label of a container an older flai started, shifts none of the others.
const containerInfoFormat = `{{.Config.Image}}|{{range $p, $b := .HostConfig.PortBindings}}{{(index $b 0).HostPort}}{{end}}|{{index .Config.Labels "` + dashboardRefLabel + `"}}`

// containerInfo reads the running container's image and host port, so
// status reflects what is actually running rather than the configuration.
// The image is the reference its dashboardRefLabel names, which a container
// started from an image ID needs (I-0116), or else its {{.Config.Image}}, for
// one an older flai started without the label.
func (a *app) containerInfo(name string) (image, url string) {
	out, err := a.runner.Run("", "docker", "inspect", "--format", containerInfoFormat, name)
	if err != nil {
		return "", ""
	}
	parts := strings.SplitN(strings.TrimSpace(out), "|", 3)
	image = strings.TrimSpace(parts[0])
	if len(parts) >= 2 && strings.TrimSpace(parts[1]) != "" {
		url = "http://localhost:" + strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 && strings.TrimSpace(parts[2]) != "" {
		image = strings.TrimSpace(parts[2])
	}
	return image, url
}

// imageProvenance is where a running container's image comes from: its
// repository digest, and the flaiover release whose verified digest list
// names it (ADR-0070).
type imageProvenance struct {
	Digest  string // sha256:… of the RepoDigests entry a list names, else of the first; "" when docker says none
	Release string // the release whose cached verified list names the image; "" when none does
}

// describe is the line flai dashboard status prints for p.
func (p imageProvenance) describe() string {
	switch {
	case p.Release != "":
		return fmt.Sprintf("  digest: %s (flaiover %s, on a signed list)\n", p.Digest, p.Release)
	case p.Digest != "":
		return fmt.Sprintf("  digest: %s (unsigned: no signed digest list names it)\n", p.Digest)
	}
	return "  unsigned: its image has no repository digest, so no signed digest list names it\n"
}

// runningRepoDigests are the RepoDigests of the named container's image, not
// its image ID, which is the listed digest in one image store and the
// configuration's digest in another (release-signing.md § The digest list):
// none when docker cannot say.
func (a *app) runningRepoDigests(name string) []string {
	id, err := a.containerImageID(name)
	if err != nil || id == "" {
		return nil
	}
	out, err := a.runner.Run("", "docker", "image", "inspect", id, "--format", "{{json .RepoDigests}}")
	var repoDigests []string
	if err == nil {
		err = json.Unmarshal([]byte(strings.TrimSpace(out)), &repoDigests)
	}
	if err != nil {
		a.logger().Warn("image digests not read", "component", "dashboard", "container", name, "image", id, "err", err.Error())
		return nil
	}
	return repoDigests
}

// runningProvenance matches the named container's runningRepoDigests against
// the verified lists cached beside flai serve's state, the newest release
// first.
func (a *app) runningProvenance(name string) imageProvenance {
	repoDigests := a.runningRepoDigests(name)
	if len(repoDigests) == 0 {
		return imageProvenance{}
	}
	cache, err := a.dashboardCache(dashboardReleases())
	if err != nil {
		a.logger().Warn("cached digest lists not read", "component", "dashboard", "err", err.Error())
		return imageProvenance{Digest: digestOf(repoDigests[0])}
	}
	lists, skipped := cache.All()
	if len(skipped) > 0 {
		a.logger().Warn("cached digest lists skipped", "component", "dashboard", "count", len(skipped), "err", errors.Join(skipped...).Error())
	}
	for _, l := range lists {
		for _, rd := range repoDigests {
			if l.Matches(rd) {
				return imageProvenance{Digest: digestOf(rd), Release: l.Version}
			}
		}
	}
	return imageProvenance{Digest: digestOf(repoDigests[0])}
}

// serveImage is flai serve's last verdict on the image the dashboard
// container runs (S-0238): the one it recorded for root, else, as one
// container serves every project, the one of the first root it recorded one
// for; false when flai serve does not run or has recorded none.
func (a *app) serveImage(root string) (serve.ImageVerdict, bool) {
	st, alive := a.serveDir().ReadStatus(time.Now())
	if !alive || len(st.Images) == 0 {
		return serve.ImageVerdict{}, false
	}
	if v, ok := st.Images[root]; ok {
		return v, true
	}
	return st.Images[slices.Sorted(maps.Keys(st.Images))[0]], true
}

// describeServeImage is the line flai dashboard status prints for what flai
// serve did with the image v judges. from is the same image's provenance:
// when its line already names the digest and release v lists, this one does
// not repeat them.
func describeServeImage(v serve.ImageVerdict, from imageProvenance) string {
	switch v.Kind {
	case serve.ImageListed:
		if v.Digest == from.Digest && v.Release == from.Release {
			return "  flai serve: dials it, signed\n"
		}
		return fmt.Sprintf("  flai serve: dials it, signed (flaiover %s, %s)\n", v.Release, v.Digest)
	case serve.ImageUnlisted:
		digests := "no RepoDigests"
		if len(v.RepoDigests) > 0 {
			digests = strings.Join(v.RepoDigests, ", ")
		}
		return fmt.Sprintf("  flai serve: refused it: the dashboard runs an unsigned image %s (%s); it does not dial it\n", v.Image, digests)
	case serve.ImageNotMeasured:
		return fmt.Sprintf("  flai serve: image not measured: %s; the release stamp alone is checked\n", v.Reason)
	}
	return ""
}

// digestOf is the digest of a RepoDigests entry, <repository>@sha256:….
func digestOf(repoDigest string) string {
	_, digest, _ := strings.Cut(repoDigest, "@")
	return digest
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
			if len(remaining) == 0 {
				// no project is left for it: flai host stops watching it (S-0184)
				forgetDashboard(string(a.serveDir()))
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
		Short: "Show whether the dashboard container runs and answers: running, not answering, or gone",
		Long: `Asks docker whether the shared dashboard container runs, and its published
port whether /_health answers (S-0184): running when it does, not answering
when the container runs and it does not, gone when no container runs. Docker's
own verdict from the image's HEALTHCHECK is shown beside it. flai host
restarts a dashboard that is gone or not answering unless dashboard.no_restart
is set. For a running container it shows its image's digest and the flaiover
release whose signed digest list, kept beside flai serve's state, names it, or
that it is unsigned. While flai serve runs, a line beside it says what flai
serve did with that image at its last dial (S-0238): dials it, signed;
refused it, as the dashboard runs an unsigned image no signed list names; or
image not measured, as Docker where flai serve runs cannot say what the
dashboard runs, with why, so the release stamp alone is checked. --json gives
that verdict as serve_image.`,
		Args: cobra.NoArgs,
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
			state, health, err := a.dashboardState(s.Name)
			if err != nil {
				return err
			}
			running := state != host.DashboardGone
			image, url := s.ref(), s.url()
			if running {
				if i, u := a.containerInfo(s.Name); i != "" {
					image, url = i, orDefault(u, url)
				}
			}
			var stale []string
			var from imageProvenance
			var verdict serve.ImageVerdict
			judged := false
			if running {
				stale = a.staleMounts(s.Name)
				from = a.runningProvenance(s.Name)
				verdict, judged = a.serveImage(s.Root)
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
				out := map[string]any{"container": s.Name, "running": running, "state": state, "url": url, "image": image, "serves": names}
				if health != "" {
					out["health"] = health
				}
				if len(stale) > 0 {
					out["stale_mounts"] = stale
				}
				out["signed"] = from.Release != ""
				if from.Digest != "" {
					out["digest"] = from.Digest
				}
				if from.Release != "" {
					out["release"] = from.Release
				}
				if judged {
					out["serve_image"] = verdict
				}
				out["host_flai"] = a.hostFlai(s.Root)
				return a.printJSON(out)
			}
			if running {
				dockerSays := ""
				if health != "" {
					dockerSays = "; docker: " + health
				}
				if state == host.DashboardRunning {
					fmt.Fprintf(a.out, "%s running at %s (%s%s)\n", s.Name, url, image, dockerSays)
				} else {
					fmt.Fprintf(a.out, "%s not answering at %s (%s%s): the container runs and /_health does not answer; flai host restarts it unless dashboard.no_restart is set, or run flai dashboard restart\n", s.Name, url, image, dockerSays)
				}
				fmt.Fprint(a.out, from.describe())
				if judged {
					fmt.Fprint(a.out, describeServeImage(verdict, from))
				}
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
				fmt.Fprintf(a.out, "%s gone: no container runs; start it with flai dashboard\n", s.Name)
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
