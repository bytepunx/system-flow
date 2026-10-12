package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// dashboardRefLabel is the label every dashboard container carries naming the
// image reference it was started for, such as ghcr.io/bytepunx/flaiover:latest:
// a container started from an image ID, as a restart starts one (I-0116),
// has that ID as its {{.Config.Image}}, and this keeps the tag to show.
const dashboardRefLabel = "io.bytepunx.flai.ref"

// allowUnsignedEnv is the variable that, set to 1, lets the dashboard serve a
// flai with no release stamp (release-signing.md § The development case).
const allowUnsignedEnv = "FLAIOVER_ALLOW_UNSIGNED"

// containerArgs is the docker run argument list every dashboard container
// shares: a port publish, the two secrets (ADR-0031), allowUnsignedEnv=1 when
// allowUnsigned, and, off Windows, running as this user so the secrets
// (readable by this user only) can be read. publish is a full --publish
// value: a fixed host port for the real container, or "127.0.0.1::<port>" for
// an upgrade's temporary one, which docker assigns an ephemeral port for. It
// runs image, a reference or an image ID, and labels the container with ref,
// the reference to show for it.
func containerArgs(dir, name, publish, image, ref string, allowUnsigned bool) []string {
	args := []string{"run", "--detach", "--rm", "--name", name, "--publish", publish}
	args = append(args, tokenArgs(dir)...)
	args = append(args, agentKeyArgs(dir)...)
	if allowUnsigned {
		args = append(args, "--env", allowUnsignedEnv+"=1")
	}
	if runtime.GOOS != "windows" {
		args = append(args, "--user", strconv.Itoa(os.Getuid())+":"+strconv.Itoa(os.Getgid()))
	}
	args = append(args, "--label", dashboardRefLabel+"="+ref, image)
	return args
}

// startContainer runs a dashboard container from ref, labelled with ref.
func (a *app) startContainer(dir, name, publish, ref string, allowUnsigned bool) (string, error) {
	return a.startContainerImage(dir, name, publish, ref, ref, allowUnsigned)
}

// startContainerImage runs a dashboard container from image, labelled with
// ref, the reference to show for it, allowing unsigned builds when
// allowUnsigned.
func (a *app) startContainerImage(dir, name, publish, image, ref string, allowUnsigned bool) (string, error) {
	return a.runner.Run("", "docker", containerArgs(dir, name, publish, image, ref, allowUnsigned)...)
}

// keepAllowance is s allowing unsigned builds when the dashboard flai
// dashboard recorded was started allowing them: a restart or an upgrade keeps
// what the start was given, a one-start --allow-unsigned or --build included,
// and the settings can only add to it.
func keepAllowance(dir string, s dashboardSettings) dashboardSettings {
	if r, ok := readDashboardRecord(dir); ok && r.AllowUnsigned {
		s.AllowUnsigned = true
	}
	return s
}

// imageID is the local image ID (docker image inspect's {{.Id}}) of ref,
// which must already be pulled.
func (a *app) imageID(ref string) (string, error) {
	out, err := a.runner.Run("", "docker", "image", "inspect", ref, "--format", "{{.Id}}")
	return strings.TrimSpace(out), err
}

// containerImageID is the local image ID a running container was started
// from ({{.Image}}, comparable directly with imageID's result).
func (a *app) containerImageID(name string) (string, error) {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{.Image}}", name)
	return strings.TrimSpace(out), err
}

func newDashboardRestartCmd(a *app) *cobra.Command {
	var image, tag, bind string
	var port int
	c := &cobra.Command{
		Use:   "restart",
		Short: "Stop and start the dashboard container again, with whatever image it is already running",
		Long: `Cycles the shared dashboard container's process without changing its image: if it
is running, this stops it and starts again the exact image it was running, by
its image ID, never a fresh pull and never what its tag names now, so a
floating tag such as latest cannot silently upgrade it, even after flai
dashboard check pulled a newer one; flai dashboard upgrade does that
deliberately. It still names the reference the image was started for: for a
flaiover release, ghcr.io/bytepunx/flaiover@sha256:<digest>. If it is not
running, this starts the dashboard flai last recorded, by its image ID while
that is still here, else by its reference, without resolving dashboard.tag
again; with no record, or with --image or --tag, it starts what flai
dashboard would, a release refused starting nothing. Every project
registered with flai serve keeps its registration: only the container's
process restarts. It allows unsigned builds (FLAIOVER_ALLOW_UNSIGNED=1) when
the start it repeats did, by --allow-unsigned, --build, or the setting, or
dashboard.allow_unsigned is on now.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardRestart(cmd.Context(), image, tag, port, bind)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "latest or a flaiover release X.Y.Z, when no container runs; for another --image, its tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on (default: manifest, then config, then 0.0.0.0)")
	return c
}

func (a *app) runDashboardRestart(ctx context.Context, image, tag string, port int, bind string) error {
	repo, err := a.projectOrNone()
	if err != nil {
		return err
	}
	if err := a.requireDocker(); err != nil {
		return err
	}
	s, err := a.dashboardSettings(repo, image, tag, port, bind)
	if err != nil {
		return err
	}
	dir := string(a.serveDir())
	s = keepAllowance(dir, s)
	wasRunning, err := a.containerRunning(s.Name)
	if err != nil {
		return err
	}
	var ref, runImage string
	if wasRunning {
		// The image ID, not the tag: flai dashboard check may have pulled a
		// newer image under the tag since the container started (I-0116).
		ref = s.ref()
		if running, _ := a.containerInfo(s.Name); running != "" {
			ref = running
		}
		runImage = ref
		if current, _ := a.containerImageID(s.Name); current != "" {
			runImage = current
		}
		a.logger().Info("stopping dashboard for restart", "component", "dashboard", "container", s.Name, "image", runImage)
		if err := a.stopAndRelease(s.Name); err != nil {
			return err
		}
	} else if ref, runImage, err = a.restartImage(ctx, s, image != "" || tag != ""); err != nil {
		return err
	}
	id, err := a.startDashboardImage(dir, s, runImage, ref)
	if err != nil {
		if wasRunning {
			return fmt.Errorf("restart %s: it stopped but did not start again from %s: %w; the dashboard is down: run flai dashboard to start it", s.Name, refAndImage(ref, runImage), err)
		}
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"container": s.Name, "was_running": wasRunning, "ref": ref, "image": runImage, "id": short(id), "url": s.url()})
	}
	verb := "started"
	if wasRunning {
		verb = "restarted"
	}
	fmt.Fprintf(a.out, "%s %s (%s) at %s\n", verb, s.Name, ref, s.url())
	return nil
}

// restartImage is what a restart starts when no container runs: the recorded
// dashboard, by its image ID while that is still here, as flai host's watch
// starts it (ADR-0118), unless named, --image or --tag, asks for another;
// with no record, or named, what flai dashboard starts, which for the signed
// repository is the release s's tag resolves to, by digest (ADR-0070). ref is
// the reference to show for it; image, what docker runs, is pulled first when
// it is a reference not present for the daemon's platform.
func (a *app) restartImage(ctx context.Context, s dashboardSettings, named bool) (ref, image string, err error) {
	if r, ok := readDashboardRecord(string(a.serveDir())); ok && !named {
		if ref, err = a.pinnedRecordRef(ctx, r.Ref); err != nil {
			return "", "", err
		}
		image = ref
		if ref == r.Ref && r.Image != "" {
			if id, err := a.imageID(r.Image); err == nil && id != "" {
				image = r.Image
			}
		}
	} else {
		img, err := a.chooseDashboardImage(ctx, s, false)
		if err != nil {
			return "", "", err
		}
		ref, image = img.Ref, img.Ref
	}
	if image == ref && a.needsPull(ref) {
		if err := a.pullDashboardRef(s.Image, ref); err != nil {
			return "", "", err
		}
	}
	return ref, image, nil
}

// refAndImage names what a restart ran: the reference, and the image ID when
// it ran one other than the reference.
func refAndImage(ref, image string) string {
	if image == ref {
		return ref
	}
	return fmt.Sprintf("%s (image %s)", ref, image)
}

func newDashboardCheckCmd(a *app) *cobra.Command {
	var image, tag string
	c := &cobra.Command{
		Use:   "check",
		Short: "Report whether a newer dashboard image is available, without changing anything",
		Long: `Resolves the configured (or --tag) tag as flai dashboard does and compares the
result with what the running container runs. For ghcr.io/bytepunx/flaiover it
verifies the release's signed digest list and asks whether a digest of the
running image is on it, pulling nothing; a release refused is an error saying
why. For --image naming another repository it pulls the tag and compares
image IDs, and reports the image unsigned. Changes nothing: flai dashboard
upgrade applies it. Reports plainly, not as an error, when the container is
not running.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardCheck(cmd.Context(), image, tag)
		},
	}
	c.Flags().StringVar(&image, "image", "", "image name (default: manifest, then config)")
	c.Flags().StringVar(&tag, "tag", "", "latest or a flaiover release X.Y.Z to check against; for another --image, its tag (default: manifest, then config)")
	return c
}

func (a *app) runDashboardCheck(ctx context.Context, image, tag string) error {
	repo, err := a.projectOrNone()
	if err != nil {
		return err
	}
	if err := a.requireDocker(); err != nil {
		return err
	}
	s, err := a.dashboardSettings(repo, image, tag, 0, "")
	if err != nil {
		return err
	}
	running, err := a.containerRunning(s.Name)
	if err != nil {
		return err
	}
	if !running {
		if a.jsonOut {
			return a.printJSON(map[string]any{"container": s.Name, "running": false})
		}
		fmt.Fprintf(a.out, "%s is not running; flai dashboard starts it\n", s.Name)
		return nil
	}
	img, err := a.chooseDashboardImage(ctx, s, false)
	if err != nil {
		return err
	}
	// A release is compared by digest, with nothing pulled (ADR-0070); an
	// unsigned image by the image ID its tag names once pulled, as before.
	var available bool
	if img.Release != "" {
		available = !img.runs(a.runningRepoDigests(s.Name))
	} else {
		if err := a.pullDashboardRef(s.Image, img.Ref); err != nil {
			return err
		}
		latest, err := a.imageID(img.Ref)
		if err != nil {
			return err
		}
		current, err := a.containerImageID(s.Name)
		if err != nil {
			return err
		}
		available = current != "" && latest != "" && current != latest
	}
	if a.jsonOut {
		out := map[string]any{"container": s.Name, "running": true, "ref": s.ref(), "image": img.Ref, "signed": img.Release != "", "upgrade_available": available}
		if img.Release != "" {
			out["release"] = img.Release
		}
		return a.printJSON(out)
	}
	if available {
		fmt.Fprintf(a.out, "an update to %s is available; flai dashboard upgrade applies it\n", img.Ref)
	} else {
		fmt.Fprintf(a.out, "%s is already running %s\n", s.Name, img.Ref)
	}
	fmt.Fprint(a.out, img.describe())
	return nil
}

// upgradeProbeSuffix names the temporary container an upgrade starts to
// prove a new image healthy before the running one is touched.
const upgradeProbeSuffix = "-upgrade-check"

// upgradeProbeAttempts and upgradeProbeInterval bound how long an upgrade
// waits for the temporary container to answer healthy: about ten seconds,
// as a count of tries rather than a wall-clock deadline, so a test can make
// every attempt instant without needing a moving clock.
const (
	upgradeProbeAttempts = 20
	upgradeProbeInterval = 500 * time.Millisecond
)

// nameReleaseAttempts and nameReleaseInterval bound how long a restart or an
// upgrade waits, after stopping the dashboard, for Docker to finish removing
// the --rm container and release its name (I-0094): about five seconds, a
// count of tries like the health probe's.
const (
	nameReleaseAttempts = 10
	nameReleaseInterval = 500 * time.Millisecond
)

func newDashboardUpgradeCmd(a *app) *cobra.Command {
	var image, tag, bind string
	var port int
	var published bool
	var src releaseSource
	c := &cobra.Command{
		Use:   "upgrade",
		Short: "Pull a newer dashboard image, or a chosen release, and swap to it, only once it answers healthy",
		Long: `Resolves the configured (or --tag) tag as flai dashboard does. For
ghcr.io/bytepunx/flaiover that is latest, the newest flaiover release, or a
release X.Y.Z, whose signed digest list flai verifies; a release refused is
refused before anything is pulled, saying why. When no digest of the running
image is on the list, it pulls the release's image by digest,
ghcr.io/bytepunx/flaiover@sha256:<digest>. For --image naming another
repository it pulls the tag and compares image IDs, and reports the image
unsigned. When the image differs from what is running, it
starts it as a second, temporary container on a loopback port of its own,
waits for it to answer /_health, and only then stops the running container
and starts the new image at the real name and port. The running container is
never stopped until the replacement has proven healthy: if it does not
become healthy in time, the temporary container is removed and the running
one is left exactly as it was, and this reports why. Should the new image
then fail to start at the real name, the previous image is started there
again, and this still fails, saying so. The new container allows unsigned
builds (FLAIOVER_ALLOW_UNSIGNED=1) when the one it replaces was started
allowing them, or dashboard.allow_unsigned is on now.

--tag deploys that release, an earlier one included, the same way, for
this container once: dashboard.tag is not changed, so an upgrade without a
tag uses the configured one again (pin with flai config set dashboard.tag).
With --published the tag must be a published dashboard release, the bare
X.Y.Z of a flaiover/vX.Y.Z tag (flai dashboard versions lists them): any
other is refused, naming the published ones, and so is one with no digest
list this flai can verify, naming why and the deployable ones, before
anything is pulled. Only an --image of another repository, such as a
mirror's, runs a tag as given, unsigned.`,
		Example: `  flai dashboard upgrade
  flai dashboard upgrade --published --tag 0.4.0
  flai dashboard upgrade --image registry.example.com/flaiover --tag my-mirror-build`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if published {
				if tag == "" {
					return errors.New("--published checks the --tag given against the published dashboard releases: give --tag X.Y.Z with it (flai dashboard versions lists them)")
				}
				if err := requirePublishedDashboard(cmd.Context(), src.options(a, ""), tag); err != nil {
					return err
				}
			}
			return a.runDashboardUpgrade(cmd.Context(), image, tag, port, bind)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "latest or a flaiover release X.Y.Z to upgrade to, for this container once; for another --image, its tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on (default: manifest, then config, then 0.0.0.0)")
	f.BoolVar(&published, "published", false, "refuse a --tag that is not a published dashboard release, before anything is pulled")
	src.register(c, "GitHub repository whose flaiover/vX.Y.Z tags are the published dashboard releases (with --published)")
	return c
}

func (a *app) runDashboardUpgrade(ctx context.Context, image, tag string, port int, bind string) error {
	repo, err := a.projectOrNone()
	if err != nil {
		return err
	}
	if err := a.requireDocker(); err != nil {
		return err
	}
	s, err := a.dashboardSettings(repo, image, tag, port, bind)
	if err != nil {
		return err
	}
	dir := string(a.serveDir())
	s = keepAllowance(dir, s)
	// A release refused is refused here, before anything is pulled (ADR-0070).
	img, err := a.chooseDashboardImage(ctx, s, false)
	if err != nil {
		return err
	}
	signed := img.Release != ""
	running, err := a.containerRunning(s.Name)
	if err != nil {
		return err
	}
	if !running {
		if err := a.pullDashboardRef(s.Image, img.Ref); err != nil {
			return err
		}
		id, err := a.startDashboard(dir, s, img.Ref)
		if err != nil {
			return err
		}
		return a.printUpgradeResult(s, img, "started", "", id)
	}
	// A release is compared by digest before it is pulled; an unsigned image
	// by the image ID its tag names once pulled, as before.
	if !signed {
		if err := a.pullDashboardRef(s.Image, img.Ref); err != nil {
			return err
		}
	}
	current, err := a.containerImageID(s.Name)
	if err != nil {
		return err
	}
	var upToDate bool
	if signed {
		upToDate = img.runs(a.runningRepoDigests(s.Name))
	} else {
		latest, err := a.imageID(img.Ref)
		if err != nil {
			return err
		}
		upToDate = latest != "" && current == latest
	}
	if upToDate {
		return a.printUpgradeResult(s, img, "up-to-date", img.Ref, "")
	}
	if signed {
		if err := a.pullDashboardRef(s.Image, img.Ref); err != nil {
			return err
		}
	}

	probeName := s.Name + upgradeProbeSuffix
	_, _ = a.runner.Run("", "docker", "rm", "-f", probeName) // a previous failed attempt may have left one behind
	a.logger().Info("starting a temporary container to check the new image", "component", "dashboard", "container", probeName)
	if _, err := a.startContainer(dir, probeName, fmt.Sprintf("127.0.0.1::%d", containerPort), img.Ref, s.AllowUnsigned); err != nil {
		return fmt.Errorf("start the temporary container to check the new image: %w", err)
	}
	healthy, healthErr := a.waitHealthy(probeName)
	_, _ = a.runner.Run("", "docker", "stop", probeName)

	if !healthy {
		detail := "the new image did not answer healthy in time"
		if healthErr != nil {
			detail = healthErr.Error()
		}
		// An exit, not a --json success payload with a "failed" field inside
		// it: this must fail the same way whether or not --json was given, so
		// a host action calling it with --json always sees a real error, not
		// a 200-shaped answer that happens to say "failed".
		return fmt.Errorf("upgrade to %s failed: %s; %s keeps running unchanged", img.Ref, detail, s.Name)
	}

	a.logger().Info("new image is healthy, swapping", "component", "dashboard", "container", s.Name)
	previousRef, _ := a.containerInfo(s.Name)
	if err := a.stopAndRelease(s.Name); err != nil {
		return err
	}
	id, err := a.startDashboard(dir, s, img.Ref)
	if err != nil {
		// The image ID, not the tag: the pull may have moved a floating tag
		// such as latest to the new image this just failed to start.
		return a.restorePreviousDashboard(dir, s, img.Ref, previousRef, orDefault(current, previousRef), err)
	}
	return a.printUpgradeResult(s, img, "upgraded", previousRef, id)
}

// restorePreviousDashboard starts the image the dashboard ran before an
// upgrade to newRef stopped it, under the real name and port, after newRef
// failed to start there. It always returns an error, with --json too, as the
// upgrade failed either way: one that says the previous image runs again, or,
// when it does not start either, that the dashboard is down.
func (a *app) restorePreviousDashboard(dir string, s dashboardSettings, newRef, previousRef, previousImage string, startErr error) error {
	failed := fmt.Errorf("upgrade to %s failed: the new image did not start after %s stopped: %w", newRef, s.Name, startErr)
	shown := orDefault(previousRef, previousImage)
	a.logger().Warn("new image did not start, starting the previous one again", "component", "dashboard", "container", s.Name, "image", previousImage)
	a.releaseName(s.Name) // a failed docker run can leave a created container holding the name
	if _, err := a.startDashboardImage(dir, s, previousImage, shown); err != nil {
		return fmt.Errorf("%w; starting the previous image %s again failed too: %w; the dashboard is down: run flai dashboard to start it", failed, shown, err)
	}
	return fmt.Errorf("%w; %s runs the previous image %s again", failed, s.Name, shown)
}

// stopAndRelease stops the named container and waits until Docker has
// released its name, so a container can be started under it again at once.
func (a *app) stopAndRelease(name string) error {
	if _, err := a.runner.Run("", "docker", "stop", name); err != nil {
		return fmt.Errorf("stop %s: %w", name, err)
	}
	a.releaseName(name)
	return nil
}

// releaseName waits until Docker no longer lists a container by this name in
// any state: docker stop returns before Docker has removed a --rm container,
// and a docker run under the name meanwhile fails with a conflict (I-0094).
// When nameReleaseAttempts is spent it removes the container by force; if
// that fails too, the start that follows fails and says why.
func (a *app) releaseName(name string) {
	for i := 0; i < nameReleaseAttempts; i++ {
		if listed, err := a.containerListed(name); err == nil && !listed {
			return
		}
		if i < nameReleaseAttempts-1 {
			a.sleep(nameReleaseInterval)
		}
	}
	a.logger().Warn("container name not released after stop, removing it", "component", "dashboard", "container", name)
	if _, err := a.runner.Run("", "docker", "rm", "-f", name); err != nil {
		a.logger().Warn("container not removed", "component", "dashboard", "container", name, "err", err.Error())
	}
}

// containerListed reports whether docker lists a container by this exact
// name in any state, stopped and being removed included.
func (a *app) containerListed(name string) (bool, error) {
	out, err := a.runner.Run("", "docker", "ps", "--all", "--quiet", "--filter", "name=^/"+name+"$")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// printUpgradeResult says what an upgrade to img did: outcome up-to-date,
// started, or upgraded, from the reference the container ran before.
func (a *app) printUpgradeResult(s dashboardSettings, img dashboardImage, outcome, from, id string) error {
	to := img.Ref
	if a.jsonOut {
		out := map[string]any{"container": s.Name, "outcome": outcome, "from": from, "to": to, "id": short(id), "url": s.url(), "signed": img.Release != ""}
		if img.Release != "" {
			out["release"] = img.Release
		}
		return a.printJSON(out)
	}
	switch outcome {
	case "up-to-date":
		fmt.Fprintf(a.out, "%s is already running %s\n", s.Name, to)
	case "started":
		fmt.Fprintf(a.out, "started %s (%s) at %s\n", s.Name, to, s.url())
	case "upgraded":
		fmt.Fprintf(a.out, "upgraded %s from %s to %s at %s\n", s.Name, from, to, s.url())
	}
	fmt.Fprint(a.out, img.describe())
	return nil
}

// waitHealthy polls the temporary container's published port for /_health
// until it answers or upgradeProbeAttempts is spent.
func (a *app) waitHealthy(name string) (bool, error) {
	addr, err := a.probeAddress(name)
	if err != nil {
		return false, fmt.Errorf("could not find the temporary container's published port: %w", err)
	}
	url := "http://" + addr + "/_health"
	for i := 0; i < upgradeProbeAttempts; i++ {
		if a.healthProbe(url) {
			return true, nil
		}
		if i < upgradeProbeAttempts-1 {
			a.sleep(upgradeProbeInterval)
		}
	}
	return false, nil
}

// probeAddress is the host:port docker published a container's containerPort
// on: the first, when it is published on more than one address.
func (a *app) probeAddress(name string) (string, error) {
	out, err := a.runner.Run("", "docker", "port", name, fmt.Sprintf("%d/tcp", containerPort))
	if err != nil {
		return "", err
	}
	addr := strings.TrimSpace(strings.SplitN(out, "\n", 2)[0])
	if addr == "" {
		return "", fmt.Errorf("docker port reported nothing for %s", name)
	}
	return addr, nil
}

// httpHealthProbe is the real health check: one GET, a short timeout, a 200
// required. Tests replace app.healthProbe instead of reaching the network.
func httpHealthProbe(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url) //nolint:gosec,noctx // localhost-only, one-shot health probe of a container this command just started
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode == http.StatusOK
}
