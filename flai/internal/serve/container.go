package serve

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting a story's agent on the container harness (S-0508, ADR-0212).
//
// The adapter builds the docker run; flai serve readies the host for it.
// Before the start it refuses one that cannot succeed, each refusal saying
// what to do: docker not on PATH, the inner harness's image not on this
// host, and, through the adapter, a login the inner harness needs and the
// host has not set. It reads the inner harness's version from the image and
// records it on the run, warning when this host has no passed check of it
// (ADR-0106). It opens the story's worktree when it has none, since docker
// would make an empty folder where a missing mount should be, and makes the
// story's cache folder with the link to the operator's login. The token and
// the keys reach the docker client by environment, and pilot's start request
// its standard input (spawnWith). A stop stops the container by its name.

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

// containerImage is the image the container harness runs inner in on the
// host cfg says, "" when none is set and none is the default.
func containerImage(cfg AgentConfig, inner string) string {
	if image := cfg.host(harness.Container).Images[inner]; image != "" {
		return image
	}
	return harness.DefaultImages[inner]
}

// containerAllowed refuses the start of a story's agent, agent its story's,
// on the container harness when docker is not on PATH or the image of its
// inner harness is not on this host, saying what to do. An inner harness
// with no image is the adapter's to refuse.
func (l *launcher) containerAllowed(cfg AgentConfig, agent *manifest.Agent) error {
	program := containerProgram(cfg)
	if _, err := l.docker.LookPath(program); err != nil {
		return fmt.Errorf("the container harness runs %s, which is not on flai serve's PATH: install Docker, or name its client with flai serve agent harness container --program <path>", program)
	}
	inner := innerOf(agent)
	image := containerImage(cfg, inner)
	if image == "" {
		return nil
	}
	if out, err := l.docker.Run(l.entry.Root, program, "image", "inspect", "--format", "{{.Id}}", image); err != nil {
		return fmt.Errorf("the image %[1]s the container harness runs the inner harness %[2]s in is not on this host (%[3]s): build it with make agent-images, or name another with flai serve agent harness container --image %[2]s=<image>", image, inner, dockerSaid(out, err))
	}
	return nil
}

// containerReady readies the host for the start of a story's agent, run,
// agent its story's, on the container harness, once the adapter has taken
// it: it reads the inner harness's version from its image onto run, opens
// the story's worktree, and makes its cache folder.
func (l *launcher) containerReady(cfg AgentConfig, run *AgentRun, agent *manifest.Agent) error {
	inner := innerOf(agent)
	if inner == harness.ClaudeCode {
		l.containerVersion(cfg, run, containerImage(cfg, inner))
	}
	if err := l.containerWorktree(run.Story); err != nil {
		return err
	}
	return containerCache(l.entry.Root, run.Story, inner == harness.ClaudeCode)
}

// containerVersion records on run the version of Claude Code image holds, as
// claude --version prints it in a container of its own with no network, and
// warns when this host has no passed check of that version through any
// provider: the check in the image (ADR-0106) is not made here. A version
// that cannot be read is warned of, and the start goes on.
func (l *launcher) containerVersion(cfg AgentConfig, run *AgentRun, image string) {
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "run", "--rm", "--network", "none", "--entrypoint", "claude", image, "--version")
	fields := strings.Fields(out)
	if err != nil || len(fields) == 0 {
		l.warn("container claude code version unreadable", "story", run.Story, "image", image, "err", dockerSaid(out, err))
		return
	}
	run.HarnessVersion = fields[0]
	passed, err := l.dir.claudePassed(run.HarnessVersion)
	if err != nil {
		l.warn("claude code checks unreadable", "version", run.HarnessVersion, "image", image, "err", err.Error())
		return
	}
	if !passed {
		l.warn("container claude code version unchecked", "version", run.HarnessVersion, "image", image)
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

// containerCache makes the cache folder of story's container under the
// project at root, the user's alone, with its flai and claude folders, and,
// when login, links claude/.credentials.json to the login the container
// mounts read-only, replacing what is there, so that Claude Code finds it.
func containerCache(root, story string, login bool) error {
	cache := harness.ContainerCache(root, story)
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

// stopContainer stops run's container by the name the adapter gave it, from
// the story and the session the run recorded, once its docker client has
// been signalled: a client killed outright leaves its container running. It
// says what it did, for the stop's journal entry; a container gone with its
// client is no failure, and any other failure is warned of.
func (l *launcher) stopContainer(cfg AgentConfig, run *AgentRun) string {
	name := harness.ContainerName(run.Story, run.Session)
	out, err := l.docker.Run(l.entry.Root, containerProgram(cfg), "stop", name)
	switch {
	case err == nil:
		return "stopped its container " + name
	case strings.Contains(out, "No such container"):
		return "its container " + name + " had ended with it"
	}
	l.warn("agent container not stopped", "story", run.Story, "container", name, "err", dockerSaid(out, err))
	return "could not stop its container " + name + ": stop it with docker stop " + name
}

// dockerSaid is what a docker command said of its failure: its output,
// else the error, on one line.
func dockerSaid(out string, err error) string {
	if strings.TrimSpace(out) == "" && err != nil {
		out = err.Error()
	}
	return strings.Join(strings.Fields(out), " ")
}
