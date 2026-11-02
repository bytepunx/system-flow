package harness

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// Container is the harness that runs a story's agent in a Docker container,
// whose entrypoint, pilot, starts the inner harness its agent's config.inner
// names (S-0508, ADR-0212).
const Container = "container"

// ContainerLogin is where the container finds the operator's login to the
// inner harness, mounted read-only.
const ContainerLogin = "/run/claude-login"

// DefaultImages are the image each inner harness runs in when the host names
// none, by the inner harness's name.
var DefaultImages = map[string]string{ClaudeCode: "flai-agent-claude-code:local"}

// innerLogin are the inner harnesses that cannot start without the
// operator's login mounted.
var innerLogin = map[string]bool{ClaudeCode: true}

// The container's variables of the host's HTTP MCP server, which flai serve
// hands the docker client (serve.MCPURLEnv and serve.MCPTokenEnv), and of the
// container's own flai, which keeps its config and cache in the story's cache
// folder rather than the container's home.
const (
	containerURLEnv   = "FLAI_MCP_URL"
	containerTokenEnv = "FLAI_MCP_TOKEN"
	containerConfig   = "FLAI_CONFIG"
	containerCacheDir = "FLAI_CACHE_DIR"
)

// containerEnv are the variables every container is passed by name, in order,
// before the provider's key.
var containerEnv = []string{"FLAI_AGENT", "FLAI_SESSION", "FLAI_STORY", "FLAI_STARTED_BY", containerURLEnv, containerTokenEnv, containerConfig, containerCacheDir}

// containerSession is a session that names a container and pilot records in
// its cache as one file's name.
var containerSession = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// container runs the inner harness in a container. inner is the inner
// harness For resolved from the agent's config, for Capabilities and Reader,
// which take no request; empty is claude-code, as Adapters holds it. Start
// resolves the inner harness from its request's agent.
type container struct{ inner string }

// DefaultHost is docker, with no arguments: the container's mounts and
// variables are fixed, and the host's image, network, and login say the rest.
func (container) DefaultHost() Host { return Host{Program: "docker"} }

// Capabilities are the inner harness's, which pilot runs unchanged.
func (c container) Capabilities() Capabilities { return Adapters[c.innerName()].Capabilities() }

// Reader is the inner harness's: pilot relays its log unchanged.
func (c container) Reader() usage.Reader { return Adapters[c.innerName()].Reader() }

// MCPOverHTTP is true: the agent in the container reaches flai's MCP server
// on the host over HTTP, so flai serve hands the docker client its URL and
// token (serve.MCPAsker).
func (container) MCPOverHTTP() bool { return true }

func (c container) innerName() string {
	if c.inner == "" {
		return ClaudeCode
	}
	return c.inner
}

// innerOf is the inner harness a's config.inner names, claude-code when it
// names none. It refuses one Adapters lacks, a container, and the operator's
// command, which runs on the host.
func innerOf(a *manifest.Agent) (string, error) {
	inner := ClaudeCode
	if a != nil && a.Config["inner"] != "" {
		inner = a.Config["inner"]
	}
	if _, ok := Adapters[inner]; !ok || inner == Container || inner == Command {
		return "", fmt.Errorf("container cannot run the inner harness %q; set config inner to one of %s", inner, strings.Join(inners(), ", "))
	}
	return inner, nil
}

// inners are the inner harnesses a container runs, sorted.
func inners() []string {
	return slices.DeleteFunc(Names(), func(n string) bool { return n == Container || n == Command })
}

// ContainerName is the name of story's container in session, by which flai
// serve stops it.
func ContainerName(story, session string) string {
	if session == "" {
		return "flai-agent-" + story
	}
	return "flai-agent-" + story + "-" + session
}

// ContainerCache is the folder under the project at root that keeps story's
// inner harness's configuration and sessions and the container's flai's,
// mounted at the same path.
func ContainerCache(root, story string) string {
	return filepath.Join(root, ".flai-cache", "agents", story)
}

// Start runs docker with the story's worktree, its cache folder, and the
// operator's login read-only, mounted and nothing else; the variables of
// containerEnv and the provider's key passed by name alone, the values in the
// docker client's environment; and pilot's start request on standard input,
// which names the inner harness, the cache, and the MCP server's variables,
// never a value. It refuses a request with no story, the operator's
// arguments, an inner harness with no image, a login the inner harness needs
// and the host has not set, and what the inner harness itself refuses of the
// request, as pilot would.
func (container) Start(r Request, host Host) (Start, error) {
	if r.Story == "" {
		return Start{}, fmt.Errorf("the container harness runs a story's agent alone, and this start is for none: a strategic agent's mounts are S-0509's; run it on another harness")
	}
	if filepath.Base(r.Story) != r.Story || r.Story == ".." {
		return Start{}, fmt.Errorf("the story %q is not a story's ID", r.Story)
	}
	if len(host.Args) > 0 {
		return Start{}, fmt.Errorf("the container harness takes no arguments of the operator's, since its mounts and variables are fixed: clear them, and set its image, network, and login with flai serve agent harness container")
	}
	if !filepath.IsAbs(r.Root) {
		return Start{}, fmt.Errorf("the project's directory %q is not an absolute path, which a mount needs", r.Root)
	}
	if r.Session != "" && (!containerSession.MatchString(r.Session) || r.Session == "." || r.Session == "..") {
		return Start{}, fmt.Errorf("the session %q is not a session's name; name it with letters, digits, '.', '_', and '-' alone", r.Session)
	}
	inner, err := innerOf(r.Agent)
	if err != nil {
		return Start{}, err
	}
	image := host.Images[inner]
	if image == "" {
		image = DefaultImages[inner]
	}
	if image == "" {
		return Start{}, fmt.Errorf("no image is set for the inner harness %s: set one with flai serve agent harness container --image %s=<image>", inner, inner)
	}
	if innerLogin[inner] && host.Login == "" {
		return Start{}, fmt.Errorf("the inner harness %s signs in from the operator's login, and the host names none: set it with flai serve agent harness container --login <path>", inner)
	}
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", r.Story)
	cache := ContainerCache(r.Root, r.Story)
	if err := mountable(r.Root, worktree, cache, host.Login); err != nil {
		return Start{}, err
	}
	pr := pilotRequest(r, inner, worktree, cache)
	if _, err := Adapters[inner].Start(pr, Host{}); err != nil {
		return Start{}, fmt.Errorf("the inner harness %s refuses the start: %w", inner, err)
	}
	stdin, err := json.Marshal(pr)
	if err != nil {
		return Start{}, fmt.Errorf("the start request for pilot: %w", err)
	}
	network := host.Network
	if network == "" {
		network = "bridge"
	}
	program := host.Program
	if program == "" {
		program = "docker"
	}
	argv := []string{program, "run", "--rm", "-i", "--name", ContainerName(r.Story, r.Session),
		"--user", strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid()),
		"--network", network, "--add-host", "host.docker.internal:host-gateway",
		"-v", worktree + ":" + worktree, "-w", worktree,
		"-v", cache + ":" + cache}
	if host.Login != "" {
		argv = append(argv, "-v", host.Login+":"+ContainerLogin+":ro")
	}
	for _, name := range containerEnv {
		argv = append(argv, "-e", name)
	}
	var from map[string]string
	if p := r.Provider; p != nil && p.KeyEnv != "" {
		if !envName.MatchString(p.KeyEnv) {
			return Start{}, fmt.Errorf("provider %s's key_env %q is not a variable's name; name the variable with letters, digits, and underscores", r.ProviderName, p.KeyEnv)
		}
		argv = append(argv, "-e", p.KeyEnv)
		from = map[string]string{p.KeyEnv: p.KeyEnv}
	}
	argv = append(argv, image)
	env := []string{containerConfig + "=" + filepath.Join(cache, "flai", "config.json"), containerCacheDir + "=" + filepath.Join(cache, "flai", "cache")}
	return Start{Harness: Container, Argv: argv, Env: env, EnvFrom: from, Stdin: stdin}, nil
}

// pilotRequest is the start request pilot reads for r: r in the container's
// terms, rooted in the worktree, which is all of the project it mounts, its
// agent the inner harness's without inner, and the inner harness, the cache,
// and the MCP server's variables named.
func pilotRequest(r Request, inner, worktree, cache string) Request {
	pr := r
	pr.Root, pr.Flai = worktree, ""
	pr.RoleFile, pr.Builtin, pr.Inputs, pr.Resume = nil, nil, nil, false
	pr.Inner, pr.Cache, pr.MCPURLEnv, pr.MCPTokenEnv, pr.MCPHostEnv = inner, cache, containerURLEnv, containerTokenEnv, ""
	if a := r.Agent; a != nil {
		ia := *a
		ia.Harness = inner
		ia.Config = maps.Clone(a.Config)
		delete(ia.Config, "inner")
		if len(ia.Config) == 0 {
			ia.Config = nil
		}
		pr.Agent = &ia
	}
	return pr
}

// mountable refuses a mount that reaches what a container never mounts: the
// home, ~/.flai, the main checkout, or the Docker socket; and a path that a
// mount's colon-separated form cannot carry.
func mountable(root, worktree, cache, login string) error {
	for _, p := range []string{worktree, cache, login} {
		if strings.Contains(p, ":") {
			return fmt.Errorf("the path %q holds a colon, which a mount cannot carry", p)
		}
	}
	if login == "" {
		return nil
	}
	if !filepath.IsAbs(login) {
		return fmt.Errorf("the login %q is not an absolute path: set it with flai serve agent harness container --login <path>", login)
	}
	never := []string{root, "/var/run/docker.sock", "/run/docker.sock"}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		never = append(never, home, filepath.Join(home, ".flai"))
	}
	login = filepath.Clean(login)
	for _, p := range never {
		if p = filepath.Clean(p); p == login || strings.HasPrefix(p, login+string(filepath.Separator)) || login == string(filepath.Separator) {
			return fmt.Errorf("the login %s mounts %s, which a container never mounts: set the folder of the operator's login alone with flai serve agent harness container --login <path>", login, p)
		}
	}
	return nil
}
