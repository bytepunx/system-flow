package serve

import (
	"errors"
	"fmt"
	"maps"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting a story's agent on the container harness (S-0508, ADR-0212), or
// a strategic agent's, the planner's, the orchestrator's, the analyzer's, or
// the curator's (S-0509, ADR-0247).
//
// The adapter builds the docker run; flai serve readies the host for it, in
// startOf, which every start goes through. Before the start it refuses one
// that cannot succeed, each refusal saying what to do: docker not on PATH,
// the inner harness's image not on this host, and, through the adapter, a
// login the inner harness needs and the host has not set, and, once the HTTP
// MCP server's URL is known, a URL on loopback that a native engine's bridge
// cannot reach (containerReachable, S-0619, I-0254). It reads the inner
// harness's version from the image and records it on the run; a Claude Code
// version with no record of ADR-0106's check through the agent's provider is
// checked inside the image, beside the start (S-0618), and one whose check
// failed is warned of: the host's claude is not probed for it. Since docker
// would make an empty folder, the root user's, where a missing mount should
// be, it opens a story's worktree when it has none, and makes the main
// checkout's .git/hooks for a strategic agent's when it is missing; and it
// makes the cache folder of the story or the role with the link to the
// operator's login. The token and the keys reach the docker
// client by environment, and pilot's start request its standard input
// (spawnWith). A stop stops the container by its name.

// containerHost is the name a container on a network other than the host's
// reaches the host by, as the adapter's --add-host maps it.
const containerHost = "host.docker.internal"

// containerProgram is the docker client the container harness runs on the
// host cfg says.
func containerProgram(cfg AgentConfig) string {
	if p := cfg.host(harness.Container).Program; p != "" {
		return p
	}
	return harness.Adapters[harness.Container].DefaultHost().Program
}

// innerOf is the inner harness agent's config.inner names, claude-code when
// it names none; the adapter refuses one it cannot run.
func innerOf(agent *manifest.Agent) string {
	if agent != nil && agent.Config["inner"] != "" {
		return agent.Config["inner"]
	}
	return harness.ClaudeCode
}

// containerImageHost is the container harness's host cfg says, for an agent of
// the project at root whose inner harness is inner, with Images[inner] the
// image the container runs: the project's layer image when its manifest
// names agent.image.layer (S-0510), else the host's, else the default; layer
// says it is the layer's. It is the one place that image is decided, so that
// the check that it is on this host, the version read from it, and the
// adapter's docker run agree. Images is a copy; cfg's map is never changed.
// The layer applies only to an inner harness with a harness image, which the
// layer is built over. A manifest that cannot be read refuses the start: one
// without its layer would run the agent without the project's tools.
func containerImageHost(root string, cfg AgentConfig, inner string) (host harness.Host, layer bool, err error) {
	host = cfg.host(harness.Container)
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return harness.Host{}, false, fmt.Errorf("read the project's %s for the agent image layer it may name: %w; mend it (flai check says what is wrong) and start the agent again", manifest.File, err)
	}
	image := ""
	if harness.DefaultImages[inner] != "" {
		image = m.LayerImage(inner)
	}
	layer = image != ""
	if image == "" {
		image = host.Images[inner]
	}
	if image == "" {
		image = harness.DefaultImages[inner]
	}
	if image == "" {
		return host, false, nil
	}
	host.Images = maps.Clone(host.Images)
	if host.Images == nil {
		host.Images = map[string]string{}
	}
	host.Images[inner] = image
	return host, layer, nil
}

// containerAllowed refuses the start of an agent, a story's or a strategic
// one, on the container harness when docker is not on PATH or the image of
// its inner harness, as agent names it, is not on this host, saying what to
// do. An inner harness with no image is the adapter's to refuse.
func (l *launcher) containerAllowed(cfg AgentConfig, agent *manifest.Agent) error {
	program := containerProgram(cfg)
	if _, err := l.docker.LookPath(program); err != nil {
		return fmt.Errorf("the container harness runs %s, which is not on flai serve's PATH: install Docker, or name its client with flai serve agent harness container --program <path>", program)
	}
	inner := innerOf(agent)
	host, layer, err := containerImageHost(l.entry.Root, cfg, inner)
	if err != nil {
		return err
	}
	image := host.Images[inner]
	if image == "" {
		return nil
	}
	out, err := l.docker.Run(l.entry.Root, program, "image", "inspect", "--format", "{{.Id}}", image)
	switch {
	case err == nil:
		return nil
	case layer:
		return fmt.Errorf("the image %[1]s, the project's layer for the inner harness %[2]s that agent.image.layer in %[3]s names, is not on this host (%[4]s): build it with make agent-images, or with scripts/agent-image.sh --layer over the harness image already built", image, inner, manifest.File, dockerSaid(out, err))
	}
	return fmt.Errorf("the image %[1]s the container harness runs the inner harness %[2]s in is not on this host (%[3]s): build it with make agent-images, or name another with flai serve agent harness container --image %[2]s=<image>", image, inner, dockerSaid(out, err))
}

// containerReady readies the host for the start of an agent, run, on the
// container harness for req, a story's or a strategic agent's, once the
// adapter has taken it: it reads the inner harness's version from its image
// onto run, makes the main checkout's .git/hooks for a strategic agent's,
// and makes the cache folder of the story or the role, as the adapter named
// it from req. A story's worktree is opened before the adapter is asked
// (containerWorktree, startOf), since the adapter reads its git to mount it.
func (l *launcher) containerReady(cfg AgentConfig, run *AgentRun, req harness.Request) error {
	inner := innerOf(req.Agent)
	if inner == harness.ClaudeCode {
		host, _, err := containerImageHost(l.entry.Root, cfg, inner)
		if err != nil {
			return err
		}
		l.containerVersion(cfg, run, host.Images[inner], req.ProviderName)
	}
	subject := req.Story
	if subject == "" {
		subject = req.Role
		if err := containerHooks(l.entry.Root); err != nil {
			return err
		}
	}
	return containerCache(l.entry.Root, subject, inner == harness.ClaudeCode)
}

// containerVersion records on run the version of Claude Code image holds, as
// claude --version prints it in a container of its own with no network, and
// has it checked against permission_prompt inside image (ADR-0106, S-0618)
// when this host has no record of that version through provider, the one the
// agent's start names, "" for none: the check runs beside the start, which
// never waits for it. A version whose record through provider is a failure
// is warned of. With no image check, outside flai serve's run, a version this
// host has no passed check of through any provider is warned of instead. A
// version that cannot be read is warned of, and the start goes on.
func (l *launcher) containerVersion(cfg AgentConfig, run *AgentRun, image, provider string) {
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "run", "--rm", "--network", "none", "--entrypoint", "claude", image, "--version")
	fields := strings.Fields(out)
	if err != nil || len(fields) == 0 {
		l.warn("container claude code version unreadable", append(run.containerSaid(), "image", image, "err", dockerSaid(out, err))...)
		return
	}
	version := fields[0]
	run.HarnessVersion = version
	if l.imageCheck == nil {
		passed, err := l.dir.claudePassed(version)
		if err != nil {
			l.warn("claude code checks unreadable", "version", version, "image", image, "err", err.Error())
			return
		}
		if !passed {
			l.warn("container claude code version unchecked", "version", version, "image", image)
		}
		return
	}
	checks, err := l.dir.claudeChecks()
	if err != nil {
		l.warn("claude code checks unreadable", "version", version, "image", image, "err", err.Error())
		return
	}
	through := checkedThrough(provider)
	rec, recorded := checks[version][through]
	switch {
	case !recorded:
		l.imageCheck(l.entry.Root, cfg, image, version, provider)
	case rec.Outcome != ClaudeCheckPassed:
		l.warn("container claude code version failed its check", "version", version, "image", image, "provider", through)
	}
}

// containerWorktree opens story's worktree when it has none, as flai story
// start does, from the remote's branch when its narrative says it was begun
// before (ADR-0064): the container mounts it, and docker would make an empty
// folder, the root user's, where it is missing. It refuses when there is
// still none, in a project git is not in use in.
func (l *launcher) containerWorktree(story string) error {
	repo, err := workitem.Open(l.entry.Root)
	if err != nil {
		return err
	}
	worktree := repo.WorktreePath(story)
	if _, err := os.Stat(worktree); err == nil {
		return nil
	}
	_, narrativeErr := os.Stat(repo.NarrativePath(story))
	if _, err := storygit.Open(storygit.OpenOptions{Runner: l.git, Repo: repo, Story: story, FromRemote: narrativeErr == nil}); err != nil {
		return fmt.Errorf("open %s's worktree, which the container harness mounts: %w", story, err)
	}
	if _, err := os.Stat(worktree); err != nil {
		return fmt.Errorf("the container harness mounts %s's worktree, and there is none at %s: flai opens a story's worktree only in a git repository, so run its agent on another harness", story, worktree)
	}
	return nil
}

// containerHooks makes the .git/hooks folder of the main checkout at root
// when it is missing, which a strategic agent's container mounts read-only
// over the checkout (ADR-0247 §1): docker would make a missing one the root
// user's. It refuses a checkout with no .git folder, which has no hooks to
// keep from the agent.
func containerHooks(root string) error {
	git := filepath.Join(root, ".git")
	if fi, err := os.Stat(git); err != nil || !fi.IsDir() {
		return fmt.Errorf("the container harness mounts a strategic agent's main checkout with its .git/hooks read-only, and %s has no .git folder: run the agent in a container only in a git repository's main checkout, or run it on another harness", root)
	}
	hooks := filepath.Join(git, "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		return fmt.Errorf("make %s, which the container harness mounts read-only and docker would make the root user's: %w; make it yourself, or run the agent on another harness", hooks, err)
	}
	return nil
}

// containerCache makes the cache folder of the container of subject, a
// story's ID or a strategic agent's role, under the project at root, the
// user's alone, with its flai and claude folders, and, when login, links
// claude/.credentials.json to the login the container mounts read-only,
// replacing what is there, so that Claude Code finds it.
func containerCache(root, subject string, login bool) error {
	cache := harness.ContainerCache(root, subject)
	for _, d := range []string{"flai", "claude"} {
		if err := os.MkdirAll(filepath.Join(cache, d), 0o700); err != nil {
			return fmt.Errorf("make the container's cache folder %s: %w", cache, err)
		}
	}
	if err := os.Chmod(cache, 0o700); err != nil {
		return fmt.Errorf("make the container's cache folder %s the user's alone: %w", cache, err)
	}
	if !login {
		return nil
	}
	link := filepath.Join(cache, "claude", ".credentials.json")
	if err := os.Remove(link); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("replace %s with the link to the login: %w; remove it and start the agent again", link, err)
	}
	if err := os.Symlink(harness.ContainerLogin+"/.credentials.json", link); err != nil {
		return fmt.Errorf("link %s to the login the container mounts: %w", link, err)
	}
	return nil
}

// containerMCP is mcp, the HTTP MCP server's variables as mcpEnv hands them,
// with the URL's host the container's name for the host, containerHost,
// unless network is the host's own, where the host's address reaches it.
func containerMCP(mcp []string, network string) []string {
	if network == "host" {
		return mcp
	}
	out := make([]string, 0, len(mcp))
	for _, kv := range mcp {
		if name, value, _ := strings.Cut(kv, "="); name == MCPURLEnv {
			if u, err := url.Parse(value); err == nil && u.Host != "" {
				port := u.Port()
				u.Host = containerHost
				if port != "" {
					u.Host = net.JoinHostPort(containerHost, port)
				}
				kv = name + "=" + u.String()
			}
		}
		out = append(out, kv)
	}
	return out
}

// dockerDesktop is what docker info says of the operating system of Docker
// Desktop's engine, which forwards the host gateway to the host's loopback.
const dockerDesktop = "Docker Desktop"

// containerReachable refuses the start of an agent on the container harness
// whose HTTP MCP server, as mcp, mcpEnv's variables, names it, listens on
// loopback while the container joins network, a network other than the
// host's ("" is bridge, as the adapter runs it), on an engine other than
// Docker Desktop (S-0619, I-0254). Such an engine routes host.docker.internal
// to the network's gateway, a bridge address the server does not answer, so
// the agent's first MCP call would fail; Docker Desktop forwards the gateway
// to the host's loopback, where it answers (S-0505, S-0506). The refusal says
// to bind the server to the gateway, which it names, or to run the container
// on the host's network. An engine docker info cannot name is warned of and
// the start goes on: the check never refuses on a guess.
func (l *launcher) containerReachable(cfg AgentConfig, mcp []string, network string) error {
	if network == "host" {
		return nil
	}
	if network == "" {
		network = "bridge"
	}
	bind := ""
	for _, kv := range mcp {
		if name, value, _ := strings.Cut(kv, "="); name == MCPURLEnv {
			if u, err := url.Parse(value); err == nil && loopbackHost(u.Hostname()) {
				bind = u.Hostname()
			}
		}
	}
	if bind == "" {
		return nil
	}
	engine, ok := l.containerEngine(cfg)
	if !ok || engine == dockerDesktop {
		return nil
	}
	gateway := "the gateway docker network inspect " + network + " names"
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "network", "inspect", network, "--format", "{{range .IPAM.Config}}{{.Gateway}} {{end}}")
	if fields := strings.Fields(out); err == nil && len(fields) > 0 {
		gateway = fields[0]
	}
	return fmt.Errorf("the agent's container cannot reach flai's MCP server: the server listens on %[1]s, and this Docker engine (%[2]s) routes %[3]s to the gateway of the network %[4]q, %[5]s, which a server on %[1]s does not answer: bind the server there with flai serve mcp bind %[5]s, then flai host restart mcp, or run the container on the host's network with flai serve agent harness container --network host, and start the agent again", bind, engine, containerHost, network, gateway)
}

// containerEngine is the operating system docker info names for the engine
// cfg's docker client reaches, read once per launcher, and whether docker
// answered; an engine it cannot read is warned of, and asked again at the
// next start.
func (l *launcher) containerEngine(cfg AgentConfig) (string, bool) {
	l.engineMu.Lock()
	defer l.engineMu.Unlock()
	if l.engine != "" {
		return l.engine, true
	}
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "info", "--format", "{{.OperatingSystem}}")
	engine := strings.TrimSpace(out)
	if err != nil || engine == "" {
		l.warn("container engine unreadable", "err", dockerSaid(out, err))
		return "", false
	}
	l.engine = engine
	return engine, true
}

// loopbackHost says whether host, a URL's host name without its port, names
// this machine's loopback: localhost or an address in 127.0.0.0/8 or ::1.
func loopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// stopContainer stops run's container by the name the adapter gave it, from
// the story or the role (containerSubject) and the session the run recorded,
// once its docker client has been signalled: a client killed outright leaves
// its container running. It says what it did, for the stop's journal entry;
// a container gone with its client is no failure, and any other failure is
// warned of.
func (l *launcher) stopContainer(cfg AgentConfig, run *AgentRun) string {
	name := harness.ContainerName(run.containerSubject(), run.Session)
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "stop", name)
	switch {
	case err == nil:
		return "stopped its container " + name
	case strings.Contains(out, "No such container"):
		return "its container " + name + " had ended with it"
	}
	l.warn("agent container not stopped", append(run.containerSaid(), "container", name, "err", dockerSaid(out, err))...)
	return "could not stop its container " + name + ": stop it with docker stop " + name
}

// containerSubject is what the container of r, a run on the container
// harness, was named for, as the start's request named it to the adapter:
// its story, else the role of the strategic agent it is, told by the fields
// that tell the runs apart.
func (r *AgentRun) containerSubject() string {
	switch {
	case r.Story != "":
		return r.Story
	case r.Role != "":
		return r.Role
	case r.Item != "":
		return conventions.RolePlan
	case r.Focus != "":
		return conventions.RoleAnalyze
	}
	return conventions.RoleOrchestrate
}

// containerSaid is the field a line logged of r's container names it by: its
// story, or the role of a strategic agent's.
func (r *AgentRun) containerSaid() []any {
	if r.Story != "" {
		return []any{"story", r.Story}
	}
	return []any{"role", r.containerSubject()}
}

// dockerSaid is what a docker command said of its failure: its output,
// else the error, on one line.
func dockerSaid(out string, err error) string {
	if strings.TrimSpace(out) == "" && err != nil {
		out = err.Error()
	}
	return strings.Join(strings.Fields(out), " ")
}
