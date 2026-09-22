package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// containerArgs is the docker run argument list every dashboard container
// shares: a port publish, the two secrets (ADR-0031), and, off Windows,
// running as this user so the secrets (readable by this user only) can be
// read. publish is a full --publish value: a fixed host port for the real
// container, or "127.0.0.1::<port>" for an upgrade's temporary one, which
// docker assigns an ephemeral port for.
func containerArgs(dir, name, publish, ref string) []string {
	args := []string{"run", "--detach", "--rm", "--name", name, "--publish", publish}
	args = append(args, tokenArgs(dir)...)
	args = append(args, agentKeyArgs(dir)...)
	if runtime.GOOS != "windows" {
		args = append(args, "--user", strconv.Itoa(os.Getuid())+":"+strconv.Itoa(os.Getgid()))
	}
	args = append(args, ref)
	return args
}

func (a *app) startContainer(dir, name, publish, ref string) (string, error) {
	return a.runner.Run("", "docker", containerArgs(dir, name, publish, ref)...)
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
is running, this stops it and starts it again from the exact image reference
it was running, never a fresh pull, so a floating tag such as latest cannot
silently upgrade it; flai dashboard upgrade does that deliberately. If it is
not running, this starts it, the same as flai dashboard. Every project
registered with flai serve keeps its registration: only the container's
process restarts.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardRestart(image, tag, port, bind)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "image tag (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on (default: manifest, then config, then 0.0.0.0)")
	return c
}

func (a *app) runDashboardRestart(image, tag string, port int, bind string) error {
	repo, err := a.project()
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
	wasRunning, err := a.containerRunning(s.Name)
	if err != nil {
		return err
	}
	ref := s.ref()
	if wasRunning {
		if running, _ := a.containerInfo(s.Name); running != "" {
			ref = running
		}
		a.logger().Info("stopping dashboard for restart", "component", "dashboard", "container", s.Name)
		if _, err := a.runner.Run("", "docker", "stop", s.Name); err != nil {
			return fmt.Errorf("stop %s: %w", s.Name, err)
		}
	}
	id, err := a.startContainer(dir, s.Name, fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort), ref)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"container": s.Name, "was_running": wasRunning, "ref": ref, "id": short(id), "url": s.url()})
	}
	verb := "started"
	if wasRunning {
		verb = "restarted"
	}
	fmt.Fprintf(a.out, "%s %s (%s) at %s\n", verb, s.Name, ref, s.url())
	return nil
}

func newDashboardCheckCmd(a *app) *cobra.Command {
	var image, tag string
	c := &cobra.Command{
		Use:   "check",
		Short: "Report whether a newer dashboard image is available, without changing anything",
		Long: `Pulls the configured (or --tag) image and compares it against what the running
container was started from. Changes nothing: flai dashboard upgrade applies
it. Reports plainly, not as an error, when the container is not running.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardCheck(image, tag)
		},
	}
	c.Flags().StringVar(&image, "image", "", "image name (default: manifest, then config)")
	c.Flags().StringVar(&tag, "tag", "", "image tag to check against (default: manifest, then config)")
	return c
}

func (a *app) runDashboardCheck(image, tag string) error {
	repo, err := a.project()
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
	if err := a.pullDashboardImage(s); err != nil {
		return err
	}
	latest, err := a.imageID(s.ref())
	if err != nil {
		return err
	}
	current, err := a.containerImageID(s.Name)
	if err != nil {
		return err
	}
	available := current != "" && latest != "" && current != latest
	if a.jsonOut {
		return a.printJSON(map[string]any{"container": s.Name, "running": true, "ref": s.ref(), "upgrade_available": available})
	}
	if available {
		fmt.Fprintf(a.out, "an update to %s is available; flai dashboard upgrade applies it\n", s.ref())
	} else {
		fmt.Fprintf(a.out, "%s is already running %s\n", s.Name, s.ref())
	}
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

func newDashboardUpgradeCmd(a *app) *cobra.Command {
	var image, tag, bind string
	var port int
	c := &cobra.Command{
		Use:   "upgrade",
		Short: "Pull a newer dashboard image and swap to it, only once it answers healthy",
		Long: `Pulls the configured (or --tag) image and, if it differs from what is running,
starts it as a second, temporary container on a loopback port of its own,
waits for it to answer /_health, and only then stops the running container
and starts the new image at the real name and port. The running container is
never stopped until the replacement has proven healthy: if it does not
become healthy in time, the temporary container is removed and the running
one is left exactly as it was, and this reports why.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardUpgrade(image, tag, port, bind)
		},
	}
	f := c.Flags()
	f.StringVar(&image, "image", "", "image name (default: manifest, then config)")
	f.StringVar(&tag, "tag", "", "image tag to upgrade to (default: manifest, then config)")
	f.IntVar(&port, "port", 0, "host port to publish (default: manifest, then config)")
	f.StringVar(&bind, "bind", "", "host address to publish on (default: manifest, then config, then 0.0.0.0)")
	return c
}

func (a *app) runDashboardUpgrade(image, tag string, port int, bind string) error {
	repo, err := a.project()
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
	a.logger().Info("pulling image", "component", "dashboard", "image", s.ref())
	if err := a.pullDashboardImage(s); err != nil {
		return err
	}
	running, err := a.containerRunning(s.Name)
	if err != nil {
		return err
	}
	if !running {
		id, err := a.startContainer(dir, s.Name, fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort), s.ref())
		if err != nil {
			return err
		}
		return a.printUpgradeResult(s, "started", "", s.ref(), id)
	}
	latest, err := a.imageID(s.ref())
	if err != nil {
		return err
	}
	current, err := a.containerImageID(s.Name)
	if err != nil {
		return err
	}
	if latest != "" && current == latest {
		return a.printUpgradeResult(s, "up-to-date", s.ref(), s.ref(), "")
	}

	probeName := s.Name + upgradeProbeSuffix
	_, _ = a.runner.Run("", "docker", "rm", "-f", probeName) // a previous failed attempt may have left one behind
	a.logger().Info("starting a temporary container to check the new image", "component", "dashboard", "container", probeName)
	if _, err := a.startContainer(dir, probeName, fmt.Sprintf("127.0.0.1::%d", containerPort), s.ref()); err != nil {
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
		return fmt.Errorf("upgrade to %s failed: %s; %s keeps running unchanged", s.ref(), detail, s.Name)
	}

	a.logger().Info("new image is healthy, swapping", "component", "dashboard", "container", s.Name)
	previousRef, _ := a.containerInfo(s.Name)
	if _, err := a.runner.Run("", "docker", "stop", s.Name); err != nil {
		return fmt.Errorf("stop %s: %w", s.Name, err)
	}
	id, err := a.startContainer(dir, s.Name, fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort), s.ref())
	if err != nil {
		return fmt.Errorf("the previous image stopped but the new one failed to start: %w (run flai dashboard to recover)", err)
	}
	return a.printUpgradeResult(s, "upgraded", previousRef, s.ref(), id)
}

func (a *app) printUpgradeResult(s dashboardSettings, outcome, from, to string, id string) error {
	if a.jsonOut {
		return a.printJSON(map[string]any{"container": s.Name, "outcome": outcome, "from": from, "to": to, "id": short(id), "url": s.url()})
	}
	switch outcome {
	case "up-to-date":
		fmt.Fprintf(a.out, "%s is already running %s\n", s.Name, to)
	case "started":
		fmt.Fprintf(a.out, "started %s (%s) at %s\n", s.Name, to, s.url())
	case "upgraded":
		fmt.Fprintf(a.out, "upgraded %s from %s to %s at %s\n", s.Name, from, to, s.url())
	}
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

// probeAddress is the host:port docker published the temporary container's
// containerPort on.
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
