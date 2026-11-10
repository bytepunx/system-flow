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

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// Container is the harness that runs a story's agent, or a strategic agent's
// (S-0509, ADR-0247), in a Docker container, whose entrypoint, pilot, starts
// the inner harness its agent's config.inner names (S-0508, ADR-0212).
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

// containerEnv are the variables a story's container is passed by name, in
// order, before the provider's key.
var containerEnv = []string{"FLAI_AGENT", "FLAI_SESSION", "FLAI_STORY", "FLAI_STARTED_BY", containerURLEnv, containerTokenEnv, containerConfig, containerCacheDir}

// roleContainerEnv are the variables a strategic agent's container is passed
// by name, in order, before the provider's key: a story's, with FLAI_ROLE in
// FLAI_STORY's place (ADR-0247 §2).
var roleContainerEnv = []string{"FLAI_AGENT", "FLAI_SESSION", "FLAI_ROLE", "FLAI_STARTED_BY", containerURLEnv, containerTokenEnv, containerConfig, containerCacheDir}

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
	if !slices.Contains(ContainerInners(), inner) {
		return "", fmt.Errorf("container cannot run the inner harness %q; set config inner to one of %s", inner, strings.Join(ContainerInners(), ", "))
	}
	return inner, nil
}

// ContainerInners are the inner harnesses a container runs, sorted: every
// adapter but the container, the host's command, and opencode, which has no
// agent image yet and reaches flai's MCP server on the host alone (S-0406).
func ContainerInners() []string {
	return slices.DeleteFunc(Names(), func(n string) bool { return n == Container || n == Command || n == OpenCode })
}

// ContainerName is the name of the container in session of subject, a
// story's ID or a strategic agent's role, such as plan, by which flai serve
// stops it.
func ContainerName(subject, session string) string {
	if session == "" {
		return "flai-agent-" + subject
	}
	return "flai-agent-" + subject + "-" + session
}

// ContainerCache is the folder under the project at root that keeps the inner
// harness's configuration and sessions and the container's flai's for
// subject, a story's ID or a strategic agent's role, such as plan, mounted at
// the same path.
func ContainerCache(root, subject string) string {
	return filepath.Join(root, ".flai-cache", "agents", subject)
}

// Start runs docker with the story's worktree; the two git mounts
// storygit.GitMounts gives it, the main checkout's .git read-only and the
// story's registration under it read-write (ADR-0246); the three of the main
// checkout storygit.MainMounts gives it, its wip and system-flow.yaml
// read-only and its .flai-cache/verify read-write, which Start makes on the
// host when it is absent, so flai verify runs in the container (ADR-0253);
// its cache folder; and the operator's login read-only, mounted in that
// order and nothing else; the variables of
// containerEnv and the provider's key passed by name alone, the values in the
// docker client's environment; and pilot's start request on standard input,
// which names the inner harness, the cache, and the MCP server's variables,
// never a value. A strategic agent's start mounts the main checkout instead
// (roleStart). It refuses a request that is neither a story's agent's nor a
// built-in strategic agent's (containerSubject), the operator's arguments, an
// inner harness with no image, a login the inner harness needs and the host
// has not set, a worktree whose git or main checkout it cannot mount, and
// what the inner harness itself refuses of the request, as pilot would.
func (container) Start(r Request, host Host) (Start, error) {
	if err := containerSubject(r); err != nil {
		return Start{}, err
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
	if r.Story == "" {
		return roleStart(r, host, inner, image)
	}
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", r.Story)
	cache := ContainerCache(r.Root, r.Story)
	if err := mountable(r.Root, worktree, cache, host.Login); err != nil {
		return Start{}, err
	}
	gitMounts, err := storygit.GitMounts(worktree)
	if err != nil {
		return Start{}, fmt.Errorf("the story's worktree %s has no git for the container: %w", worktree, err)
	}
	gitArgs, err := mountArgs(gitMounts)
	if err != nil {
		return Start{}, err
	}
	mainMounts, err := storygit.MainMounts(worktree)
	if err != nil {
		return Start{}, fmt.Errorf("the story's worktree %s has no main checkout for the container, so flai verify could not run in it: %w", worktree, err)
	}
	mainArgs, err := mountArgs(mainMounts)
	if err != nil {
		return Start{}, err
	}
	stdin, err := pilotStdin(pilotRequest(r, inner, worktree, cache))
	if err != nil {
		return Start{}, err
	}
	if err := createMounts(mainMounts); err != nil {
		return Start{}, err
	}
	argv := append(dockerRun(host, ContainerName(r.Story, r.Session)),
		"-v", worktree+":"+worktree, "-w", worktree)
	argv = append(argv, gitArgs...)
	argv = append(argv, mainArgs...)
	argv = append(argv, "-v", cache+":"+cache)
	return containerStart(argv, containerEnv, safeMain(mainMounts), r, host, image, cache, stdin)
}

// safeMain is git's safe.directory for the main checkout of a story's
// container, given as command-line configuration through GIT_CONFIG_COUNT,
// GIT_CONFIG_KEY_0, and GIT_CONFIG_VALUE_0, which git honours for it. Docker
// makes the main checkout's own folder as root, as the parent of its mounts,
// so git refuses to run there, as flai verify does to find the main branch,
// for dubious ownership (ADR-0253). None when ms names no main checkout.
func safeMain(ms []storygit.Mount) []string {
	for _, m := range ms {
		if filepath.Base(m.Path) == "system-flow.yaml" {
			return []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=safe.directory", "GIT_CONFIG_VALUE_0=" + filepath.Dir(m.Path)}
		}
	}
	return nil
}

// mountArgs is docker's -v of each of ms at its host path, in order, :ro on
// the read-only ones. It refuses a path with a colon, which a mount cannot
// carry.
func mountArgs(ms []storygit.Mount) ([]string, error) {
	var args []string
	for _, m := range ms {
		if strings.Contains(m.Path, ":") {
			return nil, fmt.Errorf("the path %q holds a colon, which a mount cannot carry", m.Path)
		}
		v := m.Path + ":" + m.Path
		if m.ReadOnly {
			v += ":ro"
		}
		args = append(args, "-v", v)
	}
	return args, nil
}

// createMounts makes on the host each folder of ms that carries Create and is
// absent, as the operator's, since Docker would make it as root and the
// container's user could not write to it (ADR-0253 §1).
func createMounts(ms []storygit.Mount) error {
	for _, m := range ms {
		if !m.Create {
			continue
		}
		if err := os.MkdirAll(m.Path, 0o755); err != nil {
			return fmt.Errorf("making %s on the host, which the container mounts read-write: %w; make it, owned by the user flai serve runs as", m.Path, err)
		}
	}
	return nil
}

// roleStart is Start for a strategic agent, the planner, the orchestrator,
// the analyzer, or the curator (S-0509, ADR-0247): it mounts the project's
// main checkout read-write at its host path as the working directory, the
// checkout's .git/hooks read-only over it, the role's cache folder, and the
// operator's login read-only, and nothing else; it passes the variables of
// roleContainerEnv by name, FLAI_ROLE's value in the docker client's
// environment, as the claude-code adapter's is on the host; and it names the
// container for the role.
func roleStart(r Request, host Host, inner, image string) (Start, error) {
	cache := ContainerCache(r.Root, r.Role)
	if err := mountable(r.Root, r.Root, cache, host.Login); err != nil {
		return Start{}, err
	}
	stdin, err := pilotStdin(pilotRequest(r, inner, r.Root, cache))
	if err != nil {
		return Start{}, err
	}
	hooks := filepath.Join(r.Root, ".git", "hooks")
	argv := append(dockerRun(host, ContainerName(r.Role, r.Session)),
		"-v", r.Root+":"+r.Root, "-w", r.Root,
		"-v", hooks+":"+hooks+":ro",
		"-v", cache+":"+cache)
	st, err := containerStart(argv, roleContainerEnv, nil, r, host, image, cache, stdin)
	if err != nil {
		return Start{}, err
	}
	st.Env = append(st.Env, "FLAI_ROLE="+r.Role)
	return st, nil
}

// containerSubject refuses a start the container harness cannot run: one
// that names neither a story nor a role, or both; a role a project defines in
// a role file of its own, whose file the start request does not carry to
// pilot; a role that is not a built-in strategic one; and a story that is not
// a story's ID.
func containerSubject(r Request) error {
	switch {
	case r.Story == "" && r.Role == "":
		return fmt.Errorf("the container harness runs a story's agent or a strategic agent, and this start names neither a story nor a role: start it for one")
	case r.Story != "" && r.Role != "":
		return fmt.Errorf("this start names story %s and role %q, and the container harness runs a story's agent or a strategic agent, not both: start it for one", r.Story, r.Role)
	case r.RoleFile != nil:
		return fmt.Errorf("role %q is one the project defines in a role file of its own, which the container harness cannot carry to pilot: run it on another harness until a story carries its file to pilot", r.Role)
	case r.Story == "" && !slices.Contains(conventions.StrategicRoles, r.Role):
		return fmt.Errorf("the container harness runs the strategic agents in role %s, and %q is none of them: start one of them, or run this agent on another harness", strings.Join(conventions.StrategicRoles, ", "), r.Role)
	case r.Story != "" && (filepath.Base(r.Story) != r.Story || r.Story == ".."):
		return fmt.Errorf("the story %q is not a story's ID", r.Story)
	}
	return nil
}

// pilotStdin is pr, pilot's start request, as pilot reads it on standard
// input, once the inner harness has been asked to start it, so that a start
// it refuses is refused here, as pilot would refuse it.
func pilotStdin(pr Request) ([]byte, error) {
	if _, err := Adapters[pr.Inner].Start(pr, Host{}); err != nil {
		return nil, fmt.Errorf("the inner harness %s refuses the start: %w", pr.Inner, err)
	}
	stdin, err := json.Marshal(pr)
	if err != nil {
		return nil, fmt.Errorf("the start request for pilot: %w", err)
	}
	return stdin, nil
}

// dockerRun is the host's container program's run of the container named
// name, before its mounts: removed when it ends, its standard input open, as
// the user, on the host's network, bridge when the host names none, and with
// the host reachable as host.docker.internal.
func dockerRun(host Host, name string) []string {
	program, rest := dockerUser(host)
	return append([]string{program, "run", "--rm", "-i", "--name", name}, rest...)
}

// dockerUser is the host's container program, docker when it names none, and
// the arguments of its run that make the container the user's, on the host's
// network, bridge when it names none, with the host reachable as
// host.docker.internal.
func dockerUser(host Host) (program string, args []string) {
	program, args = dockerAs(host)
	return program, append(args, "--add-host", "host.docker.internal:host-gateway")
}

// dockerAs is dockerUser without the host's name: the host's container
// program and the arguments that make its container the user's, on the
// host's network, bridge when it names none.
func dockerAs(host Host) (program string, args []string) {
	network := host.Network
	if network == "" {
		network = "bridge"
	}
	program = host.Program
	if program == "" {
		program = "docker"
	}
	return program, []string{"--user", strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid()), "--network", network}
}

// CommandRun is one command flai runs in a story's container rather than on
// the host: the close-out job's flai verify, or one of the operator's checks
// (ADR-0253, I-0295).
type CommandRun struct {
	Root    string   // the project's main checkout, an absolute path
	Story   string   // the story whose worktree the command runs in
	Purpose string   // what the run is for, such as closeout or check, in the container's name
	ID      string   // tells this run's container from the story's others, such as the job's time
	Image   string   // the image the command runs in, the story's agent's
	Env     []string // variables passed by name alone, their values in the docker client's environment
	Argv    []string // the command: its program, the image's entrypoint, then its arguments
}

// CommandContainer is a CommandRun as docker runs it: its argv, the name of
// its container, by which a cancel stops it, and the variables, NAME=value,
// the docker client's environment must hold beside the run's own Env.
type CommandContainer struct {
	Name string
	Argv []string
	Env  []string
}

// ContainerCommand is the docker run of c's command in a story's container,
// removed when it ends, named flai-<purpose>-<story>-<id>, as the user, on
// the host's network: the story's worktree read-write as the working
// directory, the main checkout's .git read-only and the story's registration
// read-write over it (ADR-0246), and its wip and system-flow.yaml read-only
// and its .flai-cache/verify read-write, made on the host when absent
// (ADR-0253), in that order and nothing else; git's safe.directory for the
// main checkout and each of c.Env passed by name alone; and c.Argv's program
// as the entrypoint of c.Image. It mounts no login, no cache folder, and no
// Docker socket, and passes no MCP variable and no provider's key. It refuses
// what Start refuses of the story's mounts, a registration or a main checkout
// it cannot mount and a path holding a colon, and the operator's arguments,
// an image it is not given, an empty command, and a name that is not a
// variable's.
func ContainerCommand(c CommandRun, host Host) (CommandContainer, error) {
	switch {
	case len(host.Args) > 0:
		return CommandContainer{}, fmt.Errorf("the container harness takes no arguments of the operator's, since its mounts and variables are fixed: clear them, and set its image, network, and login with flai serve agent harness container")
	case !filepath.IsAbs(c.Root):
		return CommandContainer{}, fmt.Errorf("the project's directory %q is not an absolute path, which a mount needs", c.Root)
	case c.Story == "" || filepath.Base(c.Story) != c.Story || c.Story == "." || c.Story == "..":
		return CommandContainer{}, fmt.Errorf("the story %q is not a story's ID", c.Story)
	case !containerSession.MatchString(c.Purpose) || !containerSession.MatchString(c.ID):
		return CommandContainer{}, fmt.Errorf("the purpose %q and the id %q name no container; name them with letters, digits, '.', '_', and '-' alone", c.Purpose, c.ID)
	case c.Image == "":
		return CommandContainer{}, fmt.Errorf("no image is given for %s's %s container: set the story's agent's image, the project's agent.image.layer or flai serve agent harness container --image", c.Story, c.Purpose)
	case len(c.Argv) == 0 || c.Argv[0] == "":
		return CommandContainer{}, fmt.Errorf("no command is given to run in %s's container", c.Story)
	}
	for _, name := range c.Env {
		if !envName.MatchString(name) {
			return CommandContainer{}, fmt.Errorf("%q is not a variable's name, which %s's container is passed by name: name it with letters, digits, and underscores", name, c.Story)
		}
	}
	worktree := filepath.Join(c.Root, ".flai-cache", "worktrees", c.Story)
	if err := mountable(c.Root, worktree, worktree, ""); err != nil {
		return CommandContainer{}, err
	}
	gitMounts, err := storygit.GitMounts(worktree)
	if err != nil {
		return CommandContainer{}, fmt.Errorf("the story's worktree %s has no git for the container: %w", worktree, err)
	}
	gitArgs, err := mountArgs(gitMounts)
	if err != nil {
		return CommandContainer{}, err
	}
	mainMounts, err := storygit.MainMounts(worktree)
	if err != nil {
		return CommandContainer{}, fmt.Errorf("the story's worktree %s has no main checkout for the container, so flai verify could not run in it: %w", worktree, err)
	}
	mainArgs, err := mountArgs(mainMounts)
	if err != nil {
		return CommandContainer{}, err
	}
	if err := createMounts(mainMounts); err != nil {
		return CommandContainer{}, err
	}
	name := "flai-" + c.Purpose + "-" + c.Story + "-" + c.ID
	program, user := dockerAs(host)
	argv := append([]string{program, "run", "--rm", "--name", name}, user...)
	argv = append(argv, "-v", worktree+":"+worktree, "-w", worktree)
	argv = append(argv, gitArgs...)
	argv = append(argv, mainArgs...)
	env := safeMain(mainMounts)
	for _, kv := range env {
		n, _, _ := strings.Cut(kv, "=")
		argv = append(argv, "-e", n)
	}
	for _, n := range c.Env {
		argv = append(argv, "-e", n)
	}
	argv = append(argv, "--entrypoint", c.Argv[0], c.Image)
	argv = append(argv, c.Argv[1:]...)
	return CommandContainer{Name: name, Argv: argv, Env: env}, nil
}

// ImageClaude and ImageFlai are the claude and the flai an agent image
// holds: the claude-code image's claude on its PATH, and the flai the
// bedrock image every agent image is built over ships (ADR-0212).
const (
	ImageClaude = "claude"
	ImageFlai   = "/usr/local/bin/flai"
)

// ClaudeCodeImageCheck is the docker run of ADR-0106's check of the Claude
// Code image holds (S-0618): image's claude as the entrypoint, with args,
// in a container removed when it ends, as the user, on the host's network,
// as a story's container is (dockerUser); scratch, the check's temporary
// folder, mounted read-write at its host path with workdir below it the
// working directory, and the operator's login read-only at ContainerLogin,
// and nothing else; and each of names passed by name alone, its value in
// the docker client's environment. It refuses the operator's arguments, as
// Start does, an image or a login the host does not name, a login a
// container never mounts, a path a mount cannot carry, a working directory
// outside scratch, and a name that is not a variable's.
func ClaudeCodeImageCheck(host Host, image, scratch, workdir string, names, args []string) ([]string, error) {
	switch {
	case len(host.Args) > 0:
		return nil, fmt.Errorf("the container harness takes no arguments of the operator's, since its mounts and variables are fixed: clear them, and set its image, network, and login with flai serve agent harness container")
	case image == "":
		return nil, fmt.Errorf("no image is set for the inner harness %s: set one with flai serve agent harness container --image %s=<image>", ClaudeCode, ClaudeCode)
	case host.Login == "":
		return nil, fmt.Errorf("the inner harness %s signs in from the operator's login, and the host names none: set it with flai serve agent harness container --login <path>", ClaudeCode)
	case !filepath.IsAbs(scratch):
		return nil, fmt.Errorf("the check's folder %q is not an absolute path, which a mount needs", scratch)
	}
	if rel, err := filepath.Rel(scratch, workdir); err != nil || !filepath.IsLocal(rel) {
		return nil, fmt.Errorf("the check's working directory %s is not below its folder %s, the one the container mounts", workdir, scratch)
	}
	if err := mountable(scratch, scratch, workdir, host.Login); err != nil {
		return nil, err
	}
	program, user := dockerUser(host)
	argv := append([]string{program, "run", "--rm"}, user...)
	argv = append(argv, "-v", scratch+":"+scratch, "-w", workdir, "-v", host.Login+":"+ContainerLogin+":ro")
	for _, name := range names {
		if !envName.MatchString(name) {
			return nil, fmt.Errorf("%q is not a variable's name, which the check's container is passed by name: name it with letters, digits, and underscores", name)
		}
		argv = append(argv, "-e", name)
	}
	argv = append(argv, "--entrypoint", ImageClaude, image)
	return append(argv, args...), nil
}

// containerStart is the Start of argv, docker's run with its mounts, after
// which it mounts the operator's login read-only when the host names one,
// passes the variables of names, those of extra (NAME=value, none a secret),
// and the provider's key by name alone, and runs image; its environment points
// the container's flai into cache and holds extra's values, and stdin is
// pilot's start request. It refuses a provider's key_env that is not a
// variable's name.
func containerStart(argv, names, extra []string, r Request, host Host, image, cache string, stdin []byte) (Start, error) {
	if host.Login != "" {
		argv = append(argv, "-v", host.Login+":"+ContainerLogin+":ro")
	}
	for _, name := range names {
		argv = append(argv, "-e", name)
	}
	for _, kv := range extra {
		name, _, _ := strings.Cut(kv, "=")
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
	env := append([]string{containerConfig + "=" + filepath.Join(cache, "flai", "config.json"), containerCacheDir + "=" + filepath.Join(cache, "flai", "cache")}, extra...)
	return Start{Harness: Container, Argv: argv, Env: env, EnvFrom: from, Stdin: stdin, Container: containerRun(argv)}, nil
}

// containerRun is the container argv runs, read from argv itself so that the
// record and the run cannot disagree: the value of its --name, each -v value
// split into the host's path, the container's, and its :ro mark, and its last
// word, the image. No path of a mount holds a colon (mountable, mountArgs).
func containerRun(argv []string) *ContainerRun {
	run := &ContainerRun{Image: argv[len(argv)-1]}
	for i := 1; i < len(argv)-1; i++ {
		switch argv[i] {
		case "--name":
			run.Name = argv[i+1]
		case "-v":
			parts := strings.SplitN(argv[i+1], ":", 3)
			m := ContainerMount{Host: parts[0], Path: parts[0]}
			if len(parts) > 1 {
				m.Path = parts[1]
			}
			m.ReadOnly = len(parts) == 3 && parts[2] == "ro"
			run.Mounts = append(run.Mounts, m)
		}
	}
	return run
}

// pilotRequest is the start request pilot reads for r: r in the container's
// terms, rooted in root, all of the project it mounts, a story's worktree or
// the main checkout, its agent the inner harness's without inner, and the
// inner harness, the cache, and the MCP server's variables named.
func pilotRequest(r Request, inner, root, cache string) Request {
	pr := r
	pr.Root, pr.Flai = root, ""
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

// mountable refuses a login that reaches what a container never mounts: the
// home, ~/.flai, the main checkout, or the Docker socket; and a path that a
// mount's colon-separated form cannot carry. A story's agent's start mounts
// of the main checkout only its .git, read-only, with the story's
// registration under it read-write (ADR-0246), and its wip and
// system-flow.yaml read-only and its .flai-cache/verify read-write
// (ADR-0253), beside the story's worktree and cache folder. A strategic
// agent's mounts the main checkout itself (ADR-0247).
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
