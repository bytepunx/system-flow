//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// containerLab is providerLab with containerOn.
func containerLab(t *testing.T) (*agentLab, string) {
	t.Helper()
	lab := providerLab(t)
	return lab, containerOn(lab)
}

// containerOn puts a fake docker first on PATH for lab, sets the container
// harness's login, and runs the project's HTTP MCP server, as flai host keeps
// it, at a loopback URL with a token, which it returns. The fake answers
// image inspect as --format {{json .}} prints it, with an ID of its own for
// each image and as its labels the JSON object in a file image-labels in the
// lab's output folder, else none, recording the image in docker-inspected,
// unless a file no-image is there, or failing as a daemon that is down when
// inspect-fails is; info, recording its arguments in
// docker-info, as Docker Desktop's engine, as a native one when a file
// native-engine is there, failing when info-fails is; network inspect,
// recording its arguments in docker-network, naming the gateway 172.17.0.1,
// failing when network-fails is; a run with --entrypoint, as claude
// --version in the image, recording its arguments in docker-versioned; and
// stop, recording
// the name, failing when a file stop-fails is there. Any other run is the
// agent's, for its subject, FLAI_STORY's value, else FLAI_ROLE's: it reads
// its standard input whole into <subject>.stdin, prints stream-<subject> as
// its output when there is one, records its arguments and environment in
// <subject>.txt, and holds as the stub does.
func containerOn(lab *agentLab) string {
	t := lab.t
	t.Helper()
	bin := t.TempDir()
	out := lab.outDir
	fake := "#!/bin/sh\n" +
		"out=\"" + out + "\"\n" +
		"case \"$1\" in\n" +
		"image) echo \"$5\" >> \"$out/docker-inspected\"; [ -f \"$out/no-image\" ] && { echo \"Error: No such image: $5\"; exit 1; }; [ -f \"$out/inspect-fails\" ] && { echo \"Cannot connect to the Docker daemon\"; exit 1; }; labels=null; [ -f \"$out/image-labels\" ] && labels=$(cat \"$out/image-labels\"); echo \"{\\\"Id\\\":\\\"sha256:fake-$5\\\",\\\"Config\\\":{\\\"Labels\\\":$labels}}\"; exit 0 ;;\n" +
		"stop) echo \"$2\" >> \"$out/docker-stopped\"; [ -f \"$out/stop-fails\" ] && { echo \"Error response from daemon: cannot stop\"; exit 1; }; exit 0 ;;\n" +
		"info) echo \"$*\" >> \"$out/docker-info\"; [ -f \"$out/info-fails\" ] && { echo \"Cannot connect to the Docker daemon\"; exit 1; }; if [ -f \"$out/native-engine\" ]; then echo \"Ubuntu 24.04.1 LTS\"; else echo \"Docker Desktop\"; fi; exit 0 ;;\n" +
		"network) echo \"$*\" >> \"$out/docker-network\"; [ -f \"$out/network-fails\" ] && { echo \"Error: No such network: $3\"; exit 1; }; echo \"172.17.0.1 \"; exit 0 ;;\n" +
		"esac\n" +
		"for a in \"$@\"; do [ \"$a\" = --entrypoint ] && { echo \"$*\" >> \"$out/docker-versioned\"; echo \"2.1.312 (Claude Code)\"; exit 0; }; done\n" +
		"trap 'touch \"$out/ended-$$\"; exit 143' TERM\n" +
		"subject=\"${FLAI_STORY:-$FLAI_ROLE}\"\n" +
		"cat > \"$out/$subject.stdin\"\n" +
		"[ -f \"$out/stream-$subject\" ] && cat \"$out/stream-$subject\"\n" +
		"{ echo \"args: $*\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: ${FLAI_ROLE-unset}\"; echo \"agent: $FLAI_AGENT\"; echo \"session: $FLAI_SESSION\"; echo \"started_by: $FLAI_STARTED_BY\"; echo \"mcp: ${FLAI_MCP_URL-unset} ${FLAI_MCP_TOKEN-unset}\"; echo \"key: ${S0350_GATEWAY_KEY-unset}\"; echo \"config: ${FLAI_CONFIG-unset}\"; } > \"$out/$subject.txt\"\n" +
		"while [ -f \"$out/hold\" ] && [ ! -f \"$out/release-$subject\" ]; do sleep 0.05; done\n" +
		"touch \"$out/ended-$$\"\n"
	if err := os.WriteFile(filepath.Join(bin, "docker"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	if lab.cfg.Harnesses == nil {
		lab.cfg.Harnesses = map[string]harness.Host{}
	}
	lab.cfg.Harnesses[harness.Container] = harness.Host{Login: t.TempDir()}
	cache := filepath.Join(lab.root, ".flai-cache")
	if err := os.MkdirAll(cache, 0o700); err != nil {
		t.Fatal(err)
	}
	state := fmt.Sprintf(`{"pid": %d, "url": "http://127.0.0.1:4299/mcp", "updated": %q}`, os.Getpid(), time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(filepath.Join(cache, MCPStateFile), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	const token = "mcp-token-sentinel-S0508"
	if err := os.WriteFile(filepath.Join(cache, MCPTokenFile), []byte(token+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return token
}

// containerStory makes a story on the container harness with its worktree,
// as flai stream open leaves it, and moves it to ready.
func containerStory(lab *agentLab, title string, a *manifest.Agent) string {
	lab.t.Helper()
	id := lab.backlog(title, a)
	// A linked worktree, as git worktree add leaves it: its .git file names
	// its registration, whose commondir names the main checkout's .git, which
	// the adapter mounts beside it (ADR-0246).
	reg := filepath.Join(lab.root, ".git", "worktrees", id)
	worktree := filepath.Join(lab.root, ".flai-cache", "worktrees", id)
	for _, d := range []string{reg, worktree} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			lab.t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(reg, "commondir"), []byte("../..\n"), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: "+reg+"\n"), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.toReady(id)
	return id
}

// pilotRequest is the start request the fake docker read for story.
func (lab *agentLab) pilotRequest(story string) harness.Request {
	lab.t.Helper()
	data, err := os.ReadFile(filepath.Join(lab.outDir, story+".stdin"))
	if err != nil {
		lab.t.Fatal(err)
	}
	var r harness.Request
	if err := json.Unmarshal(data, &r); err != nil {
		lab.t.Fatalf("the start request on standard input is not JSON: %v\n%s", err, data)
	}
	return r
}

// S-0508: a story on the container harness is started as docker run, with
// pilot's start request on its standard input, closed once written, and the
// provider's key and the MCP server's token in the docker client's
// environment alone, the URL's host the container's name for the host. The
// image's Claude Code version is read and recorded on the run, and, with no
// image check, as outside flai serve's run, a version this host has no passed
// check of is warned of. The cache folder is made the user's alone, with the
// link to the login.
func TestAStoryOnTheContainerHarnessIsStartedWithItsRequestOnStandardInput(t *testing.T) {
	const key = "sk-sentinel-S0508-never-in-argv"
	t.Setenv("S0350_GATEWAY_KEY", key)
	lab, token := containerLab(t)
	ctx := context.Background()
	id := containerStory(lab, "In a container", &manifest.Agent{Harness: harness.Container, Provider: "gateway"})
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	r := lab.run(id)
	if r.Error != "" || r.PID == 0 || r.Command != "docker" || r.Harness != harness.Container || r.HarnessVersion != "2.1.312" {
		t.Fatalf("run: %+v", r)
	}
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	args, _, _ := strings.Cut(string(got), "\n")
	for _, secret := range []string{key, token, "4299"} {
		if strings.Contains(args, secret) {
			t.Errorf("docker's arguments hold %q: %s", secret, args)
		}
	}
	for _, want := range []string{"-e FLAI_MCP_TOKEN", "-e S0350_GATEWAY_KEY", "--name " + harness.ContainerName(id, r.Session), " flai-agent-claude-code:local"} {
		if !strings.Contains(args, want) {
			t.Errorf("docker's arguments lack %q: %s", want, args)
		}
	}
	for _, want := range []string{"story: " + id + "\n", "agent: " + r.Agent + "\n", "started_by: flai-serve\n", "mcp: http://host.docker.internal:4299/mcp " + token + "\n", "key: " + key + "\n",
		"config: " + filepath.Join(harness.ContainerCache(lab.root, id), "flai", "config.json") + "\n"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("docker's environment lacks %q:\n%s", want, got)
		}
	}
	if !strings.Contains(string(got), "session: 2") {
		t.Errorf("docker's environment has no FLAI_SESSION:\n%s", got)
	}
	req := lab.pilotRequest(id)
	if req.Story != id || req.Session != r.Session || req.Inner != harness.ClaudeCode || req.Cache != harness.ContainerCache(lab.root, id) || req.MCPTokenEnv != MCPTokenEnv || req.Agent == nil || req.Agent.Harness != harness.ClaudeCode {
		t.Errorf("the start request: %+v", req)
	}
	stdin, _ := os.ReadFile(filepath.Join(lab.outDir, id+".stdin"))
	cache := harness.ContainerCache(lab.root, id)
	if fi, err := os.Stat(cache); err != nil || fi.Mode().Perm() != 0o700 {
		t.Errorf("the cache folder is not the user's alone: %v %v", fi, err)
	}
	if fi, err := os.Stat(filepath.Join(cache, "flai")); err != nil || !fi.IsDir() {
		t.Errorf("no flai folder in the cache: %v", err)
	}
	if to, err := os.Readlink(filepath.Join(cache, "claude", ".credentials.json")); err != nil || to != harness.ContainerLogin+"/.credentials.json" {
		t.Errorf("the login's link: %q %v", to, err)
	}
	record, _ := os.ReadFile(lab.o.Dir.agents())
	if !strings.Contains(string(record), `"harness_version": "2.1.312"`) {
		t.Errorf("the run's record has no version:\n%s", record)
	}
	texts := map[string]string{"serve/agents.json": string(record), "the journal": fmt.Sprintf("%+v", lab.entries()), "flai serve's log": lab.logText(), "the start request": string(stdin)}
	for what, text := range texts {
		for _, secret := range []string{key, token} {
			if strings.Contains(text, secret) {
				t.Errorf("%s holds %q:\n%s", what, secret, text)
			}
		}
	}
	if n := lab.logged("container claude code version unchecked"); n != 1 || !strings.Contains(lab.logText(), "version=2.1.312") {
		t.Errorf("an unchecked version was warned of %d times:\n%s", n, lab.logText())
	}

	// a version this host has a passed check of, through any provider, is not
	// warned of; the cache's link is replaced
	if err := lab.o.Dir.recordClaudeCheck("2.1.312", "gateway", ClaudeCheck{Outcome: ClaudeCheckPassed}); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(cache, "claude", ".credentials.json")
	_ = os.Remove(link)
	_ = os.WriteFile(link, []byte("{}"), 0o600)
	again := containerStory(lab, "In a container again", &manifest.Agent{Harness: harness.Container})
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(again); return r != nil && !r.live() })
	if r := lab.run(again); r.Error != "" || r.HarnessVersion != "2.1.312" {
		t.Errorf("run: %+v", r)
	}
	if n := lab.logged("container claude code version unchecked"); n != 1 {
		t.Errorf("a checked version was warned of:\n%s", lab.logText())
	}
}

// I-0253, S-0618, ADR-0106: a container start of a claude-code agent whose
// image's Claude Code version has no record through the agent's provider has
// the check made inside the image, once, beside the start, which does not
// wait for it: before S-0618 the version was only warned of and no check ran.
// The record names the image. A version recorded through the provider is not
// checked again, and one whose check failed is warned of.
func TestAContainerStartChecksItsImagesUnrecordedClaudeCodeInsideIt(t *testing.T) {
	t.Setenv("S0350_GATEWAY_KEY", "sk-S0618-gateway")
	lab, _ := containerLab(t)
	ctx := context.Background()
	checker := newClaudeChecker(lab.o, func() []Entry { return []Entry{lab.l.entry} })
	release := make(chan struct{})
	var mu sync.Mutex
	var calls [][]string
	checker.runImage = func(_ context.Context, _ string, argv, _ []string) (string, error) {
		mu.Lock()
		calls = append(calls, argv)
		mu.Unlock()
		<-release
		// the image's claude writes what the check asks, through permission_prompt
		prompt := argv[slices.Index(argv, "-p")+1]
		_, rest, _ := strings.Cut(prompt, "create the file ")
		target, rest, _ := strings.Cut(rest, " with exactly this content")
		_, rest, _ = strings.Cut(rest, "on one line: ")
		content, _, _ := strings.Cut(rest, "\n")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		return "", os.WriteFile(target, []byte(content+"\n"), 0o644)
	}
	ran := func() int {
		mu.Lock()
		defer mu.Unlock()
		return len(calls)
	}
	lab.l.imageCheck = func(root string, cfg AgentConfig, image, version, provider string) {
		checker.considerImage(ctx, root, cfg, image, version, provider)
	}
	start := func(title string) {
		t.Helper()
		id := containerStory(lab, title, &manifest.Agent{Harness: harness.Container, Provider: "gateway"})
		lab.l.look(ctx, false)
		waitFor(t, title+" ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
		if r := lab.run(id); r.Error != "" || r.HarnessVersion != "2.1.312" {
			t.Fatalf("run: %+v", r)
		}
	}

	// the agent started and ended while the check was held
	start("In a container")
	waitFor(t, "the check ran", func() bool { return ran() == 1 })
	close(release)
	checker.wait()
	checks, err := lab.o.Dir.claudeChecks()
	if err != nil {
		t.Fatal(err)
	}
	if rec := checks["2.1.312"]["gateway"]; rec.Outcome != ClaudeCheckPassed || rec.Image != harness.DefaultImages[harness.ClaudeCode] {
		t.Errorf("the check's record: %+v in %+v", rec, checks)
	}
	argv := strings.Join(calls[0], " ")
	if !strings.Contains(argv, "--entrypoint claude "+harness.DefaultImages[harness.ClaudeCode]+" ") || !strings.Contains(argv, "--permission-prompt-tool") {
		t.Errorf("the check's docker run: %s", argv)
	}
	if n := lab.logged("claude code check started"); n != 1 || !strings.Contains(lab.logText(), "image="+harness.DefaultImages[harness.ClaudeCode]) {
		t.Errorf("the check was logged started %d times:\n%s", n, lab.logText())
	}
	if n := lab.logged("container claude code version unchecked"); n != 0 {
		t.Errorf("a version checked in its image was warned of as unchecked:\n%s", lab.logText())
	}

	// a version recorded passed through the provider is not checked again
	start("In a container again")
	checker.wait()
	if n := ran(); n != 1 {
		t.Errorf("a version recorded passed was checked again: %d checks", n)
	}

	// nor one recorded failed, which is warned of
	if err := lab.o.Dir.recordClaudeCheck("2.1.312", "gateway", ClaudeCheck{Outcome: ClaudeCheckFailed, Image: "flai-agent-claude-code:local"}); err != nil {
		t.Fatal(err)
	}
	start("In a container once more")
	checker.wait()
	if n := ran(); n != 1 {
		t.Errorf("a version recorded failed was checked again: %d checks", n)
	}
	if n := lab.logged("container claude code version failed its check"); n != 1 || !strings.Contains(lab.logText(), "provider=gateway") {
		t.Errorf("a version whose check failed was warned of %d times:\n%s", n, lab.logText())
	}
}

// S-0508: the container harness reaches the host's MCP server at
// host.docker.internal on any network but the host's, keeping the port.
func TestContainerMCPNamesTheHostOffTheHostNetwork(t *testing.T) {
	mcp := []string{MCPURLEnv + "=http://127.0.0.1:4299/mcp", MCPTokenEnv + "=t"}
	for network, want := range map[string]string{"": "http://host.docker.internal:4299/mcp", "bridge": "http://host.docker.internal:4299/mcp", "agents": "http://host.docker.internal:4299/mcp", "host": "http://127.0.0.1:4299/mcp"} {
		got := containerMCP(slices.Clone(mcp), network)
		if !slices.Equal(got, []string{MCPURLEnv + "=" + want, MCPTokenEnv + "=t"}) {
			t.Errorf("network %q: %v", network, got)
		}
	}
}

// S-0508: flai serve refuses a container start that cannot succeed before
// anything runs, saying what to do: docker not on PATH, the inner harness's
// image not on this host, naming flai image status (S-0735), a docker that
// cannot say whether it is, no login for an inner harness that needs one, and
// no worktree to mount in a project git is not in use in.
func TestAContainerStartIsRefusedWhatItLacks(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name   string
		set    func(lab *agentLab)
		bare   bool // no worktree made
		says   []string
		stdins bool
	}{
		{name: "no docker", set: func(lab *agentLab) {
			lab.cfg.Harnesses[harness.Container] = harness.Host{Program: filepath.Join(t.TempDir(), "docker"), Login: lab.cfg.Harnesses[harness.Container].Login}
		}, says: []string{"not on flai serve's PATH", "flai serve agent harness container --program"}},
		{name: "no image", set: func(lab *agentLab) { _ = os.WriteFile(filepath.Join(lab.outDir, "no-image"), nil, 0o644) },
			says: []string{"flai-agent-claude-code:local", "not on this host", "flai image status", "make agent-images", "--image claude-code=<image>"}},
		{name: "no layer image", set: func(lab *agentLab) {
			layerManifest(lab)
			_ = os.WriteFile(filepath.Join(lab.outDir, "no-image"), nil, 0o644)
		}, says: []string{"flai-agent-t-claude-code:local", "not on this host", "agent.image.layer in system-flow.yaml", "flai image status", "make agent-images", "scripts/agent-image.sh --layer"}},
		{name: "docker cannot say", set: func(lab *agentLab) { _ = os.WriteFile(filepath.Join(lab.outDir, "inspect-fails"), nil, 0o644) },
			says: []string{"flai-agent-claude-code:local", "Cannot connect to the Docker daemon", "start Docker"}},
		{name: "no login", set: func(lab *agentLab) { lab.cfg.Harnesses[harness.Container] = harness.Host{} },
			says: []string{"flai serve agent harness container --login"}},
		{name: "no worktree", bare: true, says: []string{"there is none at", "only in a git repository"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, _ := containerLab(t)
			if c.set != nil {
				c.set(lab)
			}
			var id string
			if c.bare {
				id = lab.readyWith("Refused", &manifest.Agent{Harness: harness.Container})
			} else {
				id = containerStory(lab, "Refused", &manifest.Agent{Harness: harness.Container})
			}
			lab.l.look(ctx, false)
			r := lab.run(id)
			if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 {
				t.Fatalf("started: %+v", r)
			}
			for _, s := range c.says {
				if !strings.Contains(r.Why, s) {
					t.Errorf("the refusal does not say %q: %s", s, r.Why)
				}
			}
			if _, err := os.Stat(filepath.Join(lab.outDir, id+".stdin")); err == nil {
				t.Error("the agent's container was run")
			}
		})
	}
}

// mcpListens rewrites the state of lab's HTTP MCP server, as containerOn
// wrote it, at its loopback URL, to say it listens on addrs too, as flai mcp
// http writes them since S-0732.
func mcpListens(lab *agentLab, addrs ...string) {
	lab.t.Helper()
	listens, err := json.Marshal(addrs)
	if err != nil {
		lab.t.Fatal(err)
	}
	state := fmt.Sprintf(`{"pid": %d, "url": "http://127.0.0.1:4299/mcp", "addrs": %s, "updated": %q}`, os.Getpid(), listens, time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(filepath.Join(lab.root, ".flai-cache", MCPStateFile), []byte(state), 0o600); err != nil {
		lab.t.Fatal(err)
	}
}

// I-0254, S-0619, S-0732: on a native Linux engine, which routes
// host.docker.internal to the bridge's gateway, a container start whose MCP
// server listens on no address the container reaches, neither the gateway
// nor every address, is refused before the agent's container runs, naming
// the automatic bind first, then the gateway's beside loopback, and the
// host's network as the other way; before S-0619 it started and its agent's
// first MCP call failed. An older server's state, with no addrs, listens at
// its URL's host alone. The gateway docker cannot name is named by the
// command that names it.
func TestAContainerStartANativeEnginesBridgeCannotReachIsRefused(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name    string
		network string
		addrs   []string // none: an older server's state, with no addrs
		files   []string
		says    []string
	}{
		{name: "bridge, an older server", files: []string{"native-engine"}, says: []string{
			"listens on 127.0.0.1:4299,", "(Ubuntu 24.04.1 LTS)", "host.docker.internal", `network "bridge", 172.17.0.1,`,
			"flai serve mcp bind auto", "flai serve mcp bind 127.0.0.1 172.17.0.1, then flai host restart mcp", "flai serve agent harness container --network host"}},
		{name: "bridge, loopback alone", addrs: []string{"127.0.0.1:4299"}, files: []string{"native-engine"}, says: []string{
			"listens on 127.0.0.1:4299,", `network "bridge", 172.17.0.1,`, "flai serve mcp bind auto", "flai serve mcp bind 127.0.0.1 172.17.0.1"}},
		{name: "an address that is not the gateway", addrs: []string{"127.0.0.1:4299", "10.0.0.5:4299"}, files: []string{"native-engine"}, says: []string{
			"listens on 127.0.0.1:4299, 10.0.0.5:4299,", "flai serve mcp bind auto"}},
		{name: "a network of its own", network: "agents", files: []string{"native-engine"}, says: []string{`network "agents"`, "flai serve mcp bind 127.0.0.1 172.17.0.1"}},
		{name: "gateway unreadable", files: []string{"native-engine", "network-fails"}, says: []string{
			"flai serve mcp bind 127.0.0.1 the gateway docker network inspect bridge names"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, _ := containerLab(t)
			host := lab.cfg.Harnesses[harness.Container]
			host.Network = c.network
			lab.cfg.Harnesses[harness.Container] = host
			if c.addrs != nil {
				mcpListens(lab, c.addrs...)
			}
			for _, f := range c.files {
				if err := os.WriteFile(filepath.Join(lab.outDir, f), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			id := containerStory(lab, "Refused", &manifest.Agent{Harness: harness.Container})
			lab.l.look(ctx, false)
			r := lab.run(id)
			if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 {
				t.Fatalf("started: %+v", r)
			}
			for _, s := range c.says {
				if !strings.Contains(r.Why, s) {
					t.Errorf("the refusal does not say %q: %s", s, r.Why)
				}
			}
			// the automatic bind is the first way named, then the gateway's
			// beside loopback, then the host's network
			auto, both, hostNet := strings.Index(r.Why, "flai serve mcp bind auto"), strings.Index(r.Why, "flai serve mcp bind 127.0.0.1 "), strings.Index(r.Why, "--network host")
			if auto < 0 || both < auto || hostNet < both || strings.Index(r.Why, "flai serve mcp bind") != auto {
				t.Errorf("the refusal does not name flai serve mcp bind auto first: %s", r.Why)
			}
			if _, err := os.Stat(filepath.Join(lab.outDir, id+".stdin")); err == nil {
				t.Error("the agent's container was run")
			}
		})
	}
}

// S-0732: a container start whose MCP server listens, beside loopback, on an
// address the container reaches starts with neither flai serve mcp bind
// <gateway> nor the host's network: on a native engine, the network's
// gateway, which the engine routes host.docker.internal to, or every address
// (0.0.0.0 or ::); on Docker Desktop, loopback alone. The agent is handed
// the URL at host.docker.internal on the server's port. Before S-0732 a
// native engine's start was refused whenever the URL was on loopback.
func TestAContainerStartPassesWhenTheServerListensWhereTheContainerReaches(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name   string
		addrs  []string
		native bool
	}{
		{name: "loopback and the gateway", addrs: []string{"127.0.0.1:4299", "172.17.0.1:4299"}, native: true},
		{name: "every IPv4 address", addrs: []string{"0.0.0.0:4299"}, native: true},
		{name: "every address", addrs: []string{"[::]:4299"}, native: true},
		{name: "docker desktop, loopback alone", addrs: []string{"127.0.0.1:4299"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, token := containerLab(t)
			mcpListens(lab, c.addrs...)
			if c.native {
				if err := os.WriteFile(filepath.Join(lab.outDir, "native-engine"), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			id := containerStory(lab, "Reaches the server", &manifest.Agent{Harness: harness.Container})
			lab.l.look(ctx, false)
			waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
			if r := lab.run(id); r.Error != "" || r.PID == 0 || strings.Contains(r.Why, "cannot reach flai's MCP server") {
				t.Fatalf("run: %+v", r)
			}
			got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
			if !strings.Contains(string(got), "mcp: http://host.docker.internal:4299/mcp "+token+"\n") {
				t.Errorf("the agent is not handed the server at host.docker.internal:\n%s", got)
			}
			if args, _, _ := strings.Cut(string(got), "\n"); strings.Contains(args, "--network host") {
				t.Errorf("the container runs on the host's network: %s", args)
			}
		})
	}
}

// S-0619: what the refusal leaves alone starts: a native engine with the
// container on the host's network, where loopback is the host's, and an
// engine docker info cannot name, which is warned of at each start rather
// than refused. Docker Desktop's engine, which every other container test
// runs on, is asked of once per launcher.
func TestAContainerStartTheBridgeCheckLeavesAloneStarts(t *testing.T) {
	ctx := context.Background()
	start := func(t *testing.T, lab *agentLab, title string) {
		t.Helper()
		id := containerStory(lab, title, &manifest.Agent{Harness: harness.Container})
		lab.l.look(ctx, false)
		waitFor(t, title+" ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
		// the fake agent ends at once, its story still ready, so the run's
		// outcome says so; it started, which is what is asked
		if r := lab.run(id); r.Error != "" || r.PID == 0 || strings.Contains(r.Why, "cannot reach flai's MCP server") {
			t.Fatalf("run: %+v", r)
		}
	}
	infos := func(lab *agentLab) int {
		data, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-info"))
		return strings.Count(string(data), "\n")
	}
	t.Run("host network", func(t *testing.T) {
		lab, _ := containerLab(t)
		_ = os.WriteFile(filepath.Join(lab.outDir, "native-engine"), nil, 0o644)
		host := lab.cfg.Harnesses[harness.Container]
		host.Network = "host"
		lab.cfg.Harnesses[harness.Container] = host
		start(t, lab, "On the host's network")
		if n := infos(lab); n != 0 {
			t.Errorf("docker info was asked %d times for a container on the host's network", n)
		}
	})
	t.Run("engine unreadable", func(t *testing.T) {
		lab, _ := containerLab(t)
		_ = os.WriteFile(filepath.Join(lab.outDir, "info-fails"), nil, 0o644)
		start(t, lab, "Engine unreadable")
		start(t, lab, "Engine unreadable again")
		if n := lab.logged("container engine unreadable"); n != 2 || !strings.Contains(lab.logText(), "Cannot connect to the Docker daemon") {
			t.Errorf("an unreadable engine was warned of %d times:\n%s", n, lab.logText())
		}
	})
	t.Run("docker desktop", func(t *testing.T) {
		lab, _ := containerLab(t)
		start(t, lab, "On Docker Desktop")
		start(t, lab, "On Docker Desktop again")
		if n := infos(lab); n != 1 {
			t.Errorf("docker info was asked %d times, not once per launcher", n)
		}
		if _, err := os.Stat(filepath.Join(lab.outDir, "docker-network")); err == nil {
			t.Error("the network's gateway was asked of on Docker Desktop")
		}
	})
}

// S-0619: a URL's host names loopback when it is localhost or an address in
// 127.0.0.0/8 or ::1, and nothing else does.
func TestLoopbackHost(t *testing.T) {
	for host, want := range map[string]bool{"127.0.0.1": true, "127.1.2.3": true, "::1": true, "localhost": true, "LocalHost": true,
		"172.17.0.1": false, "0.0.0.0": false, "host.docker.internal": false, "": false, "::": false} {
		if got := loopbackHost(host); got != want {
			t.Errorf("loopbackHost(%q) = %v", host, got)
		}
	}
}

// layerManifest makes the lab's project name a layer for its agents' image
// under agent.image.layer, and writes its Dockerfile, which the layer's image
// is judged by.
func layerManifest(lab *agentLab) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nagent:\n  image:\n    layer: images/agent/Dockerfile\n"
	if err := os.WriteFile(filepath.Join(lab.root, manifest.File), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(lab.root, "images", "agent"), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lab.root, "images", "agent", "Dockerfile"), []byte(layerDockerfile), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// layerDockerfile is the Dockerfile layerManifest writes.
const layerDockerfile = "ARG BASE\nFROM ${BASE}\n"

// S-0510: a story's agent on the container harness runs in the project's
// layer image, flai-agent-<key>-claude-code:local, when the manifest names
// agent.image.layer, in place of the host's image and the default: the image
// checked to be on this host, the one claude --version is read from, and the
// one docker runs are the same; under a layer the harness image it is built
// over is inspected too, before it, for the layer's base (S-0735). With no
// layer the host's image still wins over the default. The host's images are
// never changed.
func TestTheProjectsLayerImageRunsItsContainerAgents(t *testing.T) {
	ctx := context.Background()
	const hostImage = "host-agent:1"
	for _, c := range []struct {
		name      string
		layer     bool
		want      string
		inspected string
	}{
		{name: "layer", layer: true, want: "flai-agent-t-claude-code:local", inspected: hostImage + "\nflai-agent-t-claude-code:local\n"},
		{name: "no layer", want: hostImage, inspected: hostImage + "\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, _ := containerLab(t)
			if c.layer {
				layerManifest(lab)
			}
			host := lab.cfg.Harnesses[harness.Container]
			host.Images = map[string]string{harness.ClaudeCode: hostImage}
			lab.cfg.Harnesses[harness.Container] = host
			id := containerStory(lab, "In the image", &manifest.Agent{Harness: harness.Container})
			lab.l.look(ctx, false)
			waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
			if r := lab.run(id); r.Error != "" || r.PID == 0 || r.HarnessVersion != "2.1.312" {
				t.Fatalf("run: %+v", r)
			}
			got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
			args, _, _ := strings.Cut(string(got), "\n")
			if !strings.Contains(args+" ", " "+c.want+" ") {
				t.Errorf("docker run does not run %s: %s", c.want, args)
			}
			for _, other := range []string{hostImage, "flai-agent-claude-code:local", "flai-agent-t-claude-code:local"} {
				if other != c.want && strings.Contains(args, other) {
					t.Errorf("docker run names %s, not only %s: %s", other, c.want, args)
				}
			}
			if inspected, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-inspected")); string(inspected) != c.inspected {
				t.Errorf("docker image inspect was of %q, not %q", inspected, c.inspected)
			}
			if versioned, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-versioned")); !strings.Contains(string(versioned), " "+c.want+" --version") {
				t.Errorf("claude --version was not read from %s: %s", c.want, versioned)
			}
			if images := lab.cfg.Harnesses[harness.Container].Images; len(images) != 1 || images[harness.ClaudeCode] != hostImage {
				t.Errorf("the host's images were changed: %v", images)
			}
		})
	}
}

// S-0735, ADR-0295: an agent's container image whose labels are not those
// the project's flai must run it in is stale, and its start is warned of,
// with the image and why, naming flai image status, and goes on: flai serve
// refuses a stale image only once it can rebuild it (S-0751). A current
// image is not warned of, nor the image the host names with --image, which
// is the operator's own and never judged. Under a layer the layer is judged
// by its Dockerfile's recipe and the harness image it is built over.
func TestAStaleContainerImageIsWarnedOfAndTheStartGoesOn(t *testing.T) {
	ctx := context.Background()
	start := func(t *testing.T, lab *agentLab, title string) {
		t.Helper()
		id := containerStory(lab, title, &manifest.Agent{Harness: harness.Container})
		lab.l.look(ctx, false)
		waitFor(t, title+" ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
		if r := lab.run(id); r.Error != "" || r.PID == 0 {
			t.Fatalf("run: %+v", r)
		}
	}
	labels := func(t *testing.T, lab *agentLab, l map[string]string) {
		t.Helper()
		data, err := json.Marshal(l)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(lab.outDir, "image-labels"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const stale = "container agent image stale"
	t.Run("harness image", func(t *testing.T) {
		lab, _ := containerLab(t)
		start(t, lab, "In a stale image")
		log := lab.logText()
		if n := lab.logged(stale); n != 1 || !strings.Contains(log, "image="+harness.DefaultImages[harness.ClaudeCode]) || !strings.Contains(log, agentimage.LabelFlai+" label") || !strings.Contains(log, "flai image status") {
			t.Errorf("a stale image was warned of %d times:\n%s", n, log)
		}
		labels(t, lab, map[string]string{agentimage.LabelFlai: "dev"})
		start(t, lab, "In a current image")
		if n := lab.logged(stale); n != 1 {
			t.Errorf("a current image was warned of as stale:\n%s", lab.logText())
		}
	})
	t.Run("the host's own image", func(t *testing.T) {
		lab, _ := containerLab(t)
		host := lab.cfg.Harnesses[harness.Container]
		host.Images = map[string]string{harness.ClaudeCode: "host-agent:1"}
		lab.cfg.Harnesses[harness.Container] = host
		start(t, lab, "In the operator's image")
		if n := lab.logged(stale); n != 0 {
			t.Errorf("the host's --image image, which carries no labels, was judged stale:\n%s", lab.logText())
		}
	})
	t.Run("layer", func(t *testing.T) {
		lab, _ := containerLab(t)
		layerManifest(lab)
		labels(t, lab, map[string]string{agentimage.LabelFlai: "dev"})
		start(t, lab, "In a stale layer")
		log := lab.logText()
		if n := lab.logged(stale); n != 1 || !strings.Contains(log, "image=flai-agent-t-claude-code:local") || !strings.Contains(log, "kind=layer") || !strings.Contains(log, agentimage.LabelRecipe+" label") {
			t.Errorf("a stale layer was warned of %d times:\n%s", n, log)
		}
		labels(t, lab, map[string]string{agentimage.LabelFlai: "dev", agentimage.LabelRecipe: agentimage.Recipe([]byte(layerDockerfile), "", nil),
			agentimage.LabelBase: "sha256:fake-" + harness.DefaultImages[harness.ClaudeCode]})
		start(t, lab, "In a current layer")
		if n := lab.logged(stale); n != 1 {
			t.Errorf("a current layer was warned of as stale:\n%s", lab.logText())
		}
	})
}

// argMounts are the -v mounts of args, docker run's arguments as the fake
// docker recorded them, in their order.
func argMounts(args string) []harness.ContainerMount {
	fields := strings.Fields(args)
	var ms []harness.ContainerMount
	for i := 0; i+1 < len(fields); i++ {
		if fields[i] != "-v" {
			continue
		}
		parts := strings.SplitN(fields[i+1], ":", 3)
		ms = append(ms, harness.ContainerMount{Host: parts[0], Path: parts[1], ReadOnly: len(parts) == 3 && parts[2] == "ro"})
	}
	return ms
}

// recordedContainer checks that run, the run of subject on the container
// harness, records in agents.json the image it ran, the name of its
// container, and the mounts docker run was given, as the fake docker
// recorded them in <subject>.txt.
func (lab *agentLab) recordedContainer(subject string, run *AgentRun, image string) {
	t := lab.t
	t.Helper()
	got, _ := os.ReadFile(filepath.Join(lab.outDir, subject+".txt"))
	args, _, _ := strings.Cut(string(got), "\n")
	mounts := argMounts(args)
	if len(mounts) == 0 {
		t.Fatalf("docker run was given no mounts: %s", args)
	}
	name := harness.ContainerName(subject, run.Session)
	if run.Image != image || run.Container != name || !slices.Equal(run.Mounts, mounts) {
		t.Errorf("the run records image %q, container %q, and mounts %+v; want %q, %q, and %+v", run.Image, run.Container, run.Mounts, image, name, mounts)
	}
	record, _ := os.ReadFile(lab.o.Dir.agents())
	for _, want := range []string{`"image": "` + image + `"`, `"container": "` + name + `"`, `"mounts": [`} {
		if !strings.Contains(string(record), want) {
			t.Errorf("agents.json lacks %s:\n%s", want, record)
		}
	}
}

// journaled checks that the journal holds a start entry that begins with
// start: what the run started, its harness, and on the container harness its
// image and container, then what it was started for.
func (lab *agentLab) journaled(start string) {
	lab.t.Helper()
	for _, e := range lab.entries() {
		if strings.HasPrefix(e.Detail, start) {
			return
		}
	}
	lab.t.Errorf("the journal has no start entry beginning %q: %+v", start, lab.entries())
}

// I-0274: a run on the container harness records the image it ran, the
// container's name, and its mounts in agents.json, for a story's agent and a
// strategic agent's, so that a run that used the project's layer image can be
// shown to have done so once --rm has removed its container, and its journal
// start entry names the image and the container. Before S-0665 the record and
// the entry named only docker. A run on another harness records none of them
// and logs its start as before.
func TestAContainerRunRecordsItsImageContainerAndMounts(t *testing.T) {
	const layer = "flai-agent-t-claude-code:local"
	ctx := context.Background()
	lab, _ := roleContainerLab(t)
	layerManifest(lab)

	id := containerStory(lab, "Recorded in a container", &manifest.Agent{Harness: harness.Container})
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	run := lab.run(id)
	if run.Error != "" || run.Harness != harness.Container {
		t.Fatalf("run: %+v", run)
	}
	lab.recordedContainer(id, run, layer)
	lab.journaled("started docker (container, image " + layer + ", container " + harness.ContainerName(id, run.Session) + ") for " + id + " as ")

	epic := lab.epic.ID
	if !lab.planInContainer(epic) {
		t.Fatalf("the planner was not started: %+v", lab.planRun(epic))
	}
	plan := lab.planRun(epic)
	waitFor(t, "the planner's container ran", func() bool {
		data, _ := os.ReadFile(filepath.Join(lab.outDir, conventions.RolePlan+".txt"))
		return strings.Contains(string(data), "config: ")
	})
	lab.recordedContainer(conventions.RolePlan, plan, layer)
	lab.journaled("started docker (container, image " + layer + ", container " + harness.ContainerName(conventions.RolePlan, plan.Session) + ") to plan " + epic + " as ")

	host := lab.ready("On the host")
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(host); return r != nil && !r.live() })
	r := lab.run(host)
	if r.Harness == harness.Container || r.Image != "" || r.Container != "" || r.Mounts != nil {
		t.Errorf("a run on another harness records a container: %+v", r)
	}
	lab.journaled("started " + r.Command + " (" + r.Harness + ") for " + host + " as ")
}

// S-0510: containerImageHost decides the container's image: the project's layer
// for an inner harness with a harness image, else the host's, else the
// default, in a copy of the host's images; an inner harness with no harness
// image gets no layer, and is the adapter's to refuse. A manifest that cannot
// be read refuses.
func TestContainerImageHostPutsTheLayerImageInACopyOfTheHostsImages(t *testing.T) {
	lab := newAgentLab(t)
	images := map[string]string{harness.ClaudeCode: "host-agent:1"}
	cfg := AgentConfig{Harnesses: map[string]harness.Host{harness.Container: {Login: "/login", Images: images}}}
	bare := AgentConfig{Harnesses: map[string]harness.Host{}}
	for _, c := range []struct {
		name      string
		layer     bool
		cfg       AgentConfig
		inner     string
		want      string
		wantLayer bool
	}{
		{name: "host image", cfg: cfg, inner: harness.ClaudeCode, want: "host-agent:1"},
		{name: "default", cfg: bare, inner: harness.ClaudeCode, want: harness.DefaultImages[harness.ClaudeCode]},
		{name: "layer over the host's", layer: true, cfg: cfg, inner: harness.ClaudeCode, want: "flai-agent-t-claude-code:local", wantLayer: true},
		{name: "layer over the default", layer: true, cfg: bare, inner: harness.ClaudeCode, want: "flai-agent-t-claude-code:local", wantLayer: true},
		{name: "no harness image", layer: true, cfg: bare, inner: harness.Command, want: ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.layer {
				layerManifest(lab)
			} else {
				_ = os.WriteFile(filepath.Join(lab.root, manifest.File), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
			}
			host, layer, err := containerImageHost(lab.root, c.cfg, c.inner)
			if err != nil || host.Images[c.inner] != c.want || layer != c.wantLayer || host.Login != c.cfg.host(harness.Container).Login {
				t.Errorf("containerImageHost = %+v, %v, %v; want image %q, layer %v", host, layer, err, c.want, c.wantLayer)
			}
			if len(images) != 1 || images[harness.ClaudeCode] != "host-agent:1" {
				t.Errorf("the host's images were changed: %v", images)
			}
		})
	}
	_ = os.WriteFile(filepath.Join(lab.root, manifest.File), []byte("name: [\n"), 0o644)
	if _, _, err := containerImageHost(lab.root, cfg, harness.ClaudeCode); err == nil || !strings.Contains(err.Error(), "agent image layer") {
		t.Errorf("an unreadable manifest is not refused: %v", err)
	}
}

// S-0508, ADR-0058: the operator's stop of an agent on the container harness
// signals its docker client and stops its container by the name the adapter
// gave it, from the session the run recorded; a container that cannot be
// stopped is warned of, and the run ends stopped all the same.
func TestStoppingAContainerAgentStopsItsContainer(t *testing.T) {
	ctx := context.Background()
	for _, fails := range []bool{false, true} {
		t.Run(fmt.Sprintf("stop fails %v", fails), func(t *testing.T) {
			lab, _ := containerLab(t)
			if fails {
				_ = os.WriteFile(filepath.Join(lab.outDir, "stop-fails"), nil, 0o644)
			}
			lab.hold()
			id := containerStory(lab, "Stopped in a container", &manifest.Agent{Harness: harness.Container})
			lab.l.look(ctx, false)
			lab.trapped(t, id)
			started := lab.run(id)
			run, err := Stop(lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
			if err != nil {
				t.Fatal(err)
			}
			if run.Outcome != OutcomeStopped || Alive(started.PID) {
				t.Errorf("not stopped: %+v", run)
			}
			name := harness.ContainerName(id, started.Session)
			if got, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-stopped")); string(got) != name+"\n" {
				t.Errorf("docker stop was given %q, not %s", got, name)
			}
			j := lab.entries()
			detail := j[len(j)-1].Detail
			if fails {
				if !strings.Contains(detail, "docker stop "+name) || lab.logged("agent container not stopped") != 1 {
					t.Errorf("a failed docker stop was not said: %s\n%s", detail, lab.logText())
				}
				return
			}
			if !strings.Contains(detail, "stopped its container "+name) || lab.logged("agent container not stopped") != 0 {
				t.Errorf("journal: %s\n%s", detail, lab.logText())
			}
		})
	}
}

// S-0508, S-0104: an agent on the container harness that ended asking is
// started again on the answer in a container of the same session and cache
// folder, its request naming the session and what was answered, so that
// pilot resumes the inner harness's session.
func TestAContainerAgentThatEndedAskingResumesInItsSessionAndCache(t *testing.T) {
	lab, _ := containerLab(t)
	ctx := context.Background()
	lab.hold()
	id := containerStory(lab, "Asks in a container", &manifest.Agent{Harness: harness.Container})
	lab.l.look(ctx, false)
	lab.trapped(t, id)
	first := lab.run(id)
	before := lab.pilotRequest(id)
	firstArgs, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.release(id)
	waitFor(t, "it ends asking", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeAsked })
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	_ = os.Remove(filepath.Join(lab.outDir, id+".txt"))
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	lab.trapped(t, id)
	again := lab.run(id)
	after := lab.pilotRequest(id)
	if again.Session != first.Session || after.Session != before.Session || after.Session != first.Session || after.Answered != th.ID || after.Cache != before.Cache {
		t.Errorf("not resumed in its session and cache: first %+v, again %+v, request %+v", before, after, again)
	}
	againArgs, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	firstLine, _, _ := strings.Cut(string(firstArgs), "\n")
	againLine, _, _ := strings.Cut(string(againArgs), "\n")
	if firstLine != againLine || !strings.Contains(againLine, "--name "+harness.ContainerName(id, first.Session)) {
		t.Errorf("the container is not the same:\n%s\n%s", firstLine, againLine)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// roleContainerLab is containerLab for a strategic agent (S-0509): the
// project is a git checkout whose .git has its config and no hooks folder
// yet, and holds the four roles' Claude Code definitions.
func roleContainerLab(t *testing.T) (*agentLab, string) {
	t.Helper()
	lab, token := containerLab(t)
	lab.roleContainerReady()
	return lab, token
}

// roleContainerReady makes the lab's project a git checkout with its config
// and no hooks folder, and gives it the four roles' Claude Code definitions,
// which the inner harness's start reads.
func (lab *agentLab) roleContainerReady() {
	lab.t.Helper()
	if err := os.MkdirAll(filepath.Join(lab.root, ".git"), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lab.root, ".git", "config"), []byte("[core]\n\trepositoryformatversion = 0\n"), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	defs := filepath.Join(lab.root, ".claude", "agents")
	if err := os.MkdirAll(defs, 0o755); err != nil {
		lab.t.Fatal(err)
	}
	for _, def := range []string{"planner", "orchestrator", "analyzer", "curator"} {
		if err := os.WriteFile(filepath.Join(defs, def+".md"), []byte("---\nname: "+def+"\ndescription: The "+def+".\ntools: Read\n---\n\nYou are the "+def+".\n"), 0o644); err != nil {
			lab.t.Fatal(err)
		}
	}
}

// roleContainerStarted checks that run, the built-in strategic agent's in
// role, was started on the container harness as ADR-0247 says, its container
// and cache folder named for the role, and returns pilot's start request
// (strategicContainerStarted).
func (lab *agentLab) roleContainerStarted(role string, run *AgentRun) harness.Request {
	lab.t.Helper()
	return lab.strategicContainerStarted(role, role, run)
}

// strategicContainerStarted checks that run, a strategic agent's in role,
// whose container and cache folder are named for subject, the role itself
// for a built-in one and its name for one the project defines (S-0705), was
// started on the container harness as ADR-0247 says, and returns pilot's
// start request, which the fake docker recorded under role, FLAI_ROLE's
// value: docker run of the container named for subject and the session,
// mounting the main checkout read-write at its path as the working
// directory, its .git/hooks and .git/config read-only (ADR-0281), and
// subject's cache folder, and no story's worktree; FLAI_ROLE passed by name
// and FLAI_STORY not at all; the hooks folder and subject's cache folder,
// with the link to the login, made; and the image's Claude Code version
// recorded on the run.
func (lab *agentLab) strategicContainerStarted(role, subject string, run *AgentRun) harness.Request {
	t := lab.t
	t.Helper()
	if run == nil || run.Error != "" || run.PID == 0 || run.Command != "docker" || run.Harness != harness.Container || run.HarnessVersion != "2.1.312" {
		t.Fatalf("run: %+v", run)
	}
	out := filepath.Join(lab.outDir, role+".txt")
	waitFor(t, "the "+role+" container ran", func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "config: ")
	})
	got, _ := os.ReadFile(out)
	args, _, _ := strings.Cut(string(got), "\n")
	root := lab.root
	hooks := filepath.Join(root, ".git", "hooks")
	config := filepath.Join(root, ".git", "config")
	cache := harness.ContainerCache(root, subject)
	for _, want := range []string{"--name " + harness.ContainerName(subject, run.Session) + " ", " -v " + root + ":" + root + " ", " -w " + root + " ",
		" -v " + hooks + ":" + hooks + ":ro ", " -v " + config + ":" + config + ":ro ", " -v " + cache + ":" + cache + " ", " -e FLAI_ROLE ", " flai-agent-claude-code:local"} {
		if !strings.Contains(args, want) {
			t.Errorf("docker's arguments lack %q: %s", want, args)
		}
	}
	for _, never := range []string{"FLAI_STORY", filepath.Join(root, ".flai-cache", "worktrees")} {
		if strings.Contains(args, never) {
			t.Errorf("docker's arguments hold %q: %s", never, args)
		}
	}
	for _, want := range []string{"story: unset\n", "role: " + role + "\n", "agent: " + run.Agent + "\n", "started_by: flai-serve\n"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("docker's environment lacks %q:\n%s", want, got)
		}
	}
	if fi, err := os.Stat(hooks); err != nil || !fi.IsDir() {
		t.Errorf("no .git/hooks folder was made: %v", err)
	}
	if fi, err := os.Stat(cache); err != nil || fi.Mode().Perm() != 0o700 {
		t.Errorf("the cache folder %s is not the user's alone: %v %v", cache, fi, err)
	}
	if to, err := os.Readlink(filepath.Join(cache, "claude", ".credentials.json")); err != nil || to != harness.ContainerLogin+"/.credentials.json" {
		t.Errorf("the login's link: %q %v", to, err)
	}
	if record, _ := os.ReadFile(lab.o.Dir.agents()); !strings.Contains(string(record), `"harness_version": "2.1.312"`) {
		t.Errorf("the run's record has no version:\n%s", record)
	}
	req := lab.pilotRequest(role)
	if req.Role != role || req.Story != "" || req.Root != root || req.Session != run.Session || req.Inner != harness.ClaudeCode || req.Cache != cache || req.Agent == nil || req.Agent.Harness != harness.ClaudeCode {
		t.Errorf("the start request: %+v", req)
	}
	return req
}

// I-0258, S-0705: a role the project defines, on the container harness, is
// started in a container named for the role's name rather than its verb,
// which a built-in role or another defined role could share: flai serve makes
// the cache folder the adapter mounts under the name, with the link to the
// login, and none under the verb, and the run records the container by its
// name. A stop runs docker stop on that container, by the name the run
// recorded and, for a run recorded without it, by the role's name and the
// session. Before S-0705 flai serve made the cache folder under the verb,
// and the folder the container mounted was docker's to make, the root
// user's.
func TestADefinedRoleOnTheContainerHarnessIsReadiedAndStoppedByItsName(t *testing.T) {
	lab, _ := roleContainerLab(t)
	dir := filepath.Join(lab.root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "reviewer.md"), []byte(reviewerRole), 0o644); err != nil {
		t.Fatal(err)
	}
	story := lab.backlog("Reviewed in a container", nil)
	repo, err := workitem.Open(lab.root)
	if err != nil {
		t.Fatal(err)
	}
	f, err := definedRole(repo, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	lab.l.mu.Lock()
	started, _ := lab.l.role(context.Background(), lab.cfg, f, story, &manifest.Agent{Harness: harness.Container}, nil, TriggerAsked)
	lab.l.mu.Unlock()
	run := lab.roleRun("reviewer", story)
	if !started {
		t.Fatalf("not started: %+v", run)
	}
	req := lab.strategicContainerStarted("review", "reviewer", run)
	if req.RoleName != "reviewer" || req.Item != story || req.Name != "reviewer-"+story {
		t.Errorf("the start request names role %q, item %q, agent %q; want reviewer, %s, reviewer-%s", req.RoleName, req.Item, req.Name, story, story)
	}
	if _, err := os.Stat(harness.ContainerCache(lab.root, "review")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a cache folder was made under the role's verb: %v", err)
	}
	name := harness.ContainerName("reviewer", run.Session)
	if run.Container != name {
		t.Errorf("the run records container %q, want %s", run.Container, name)
	}
	waitFor(t, "it ended", func() bool { r := lab.roleRun("reviewer", story); return r != nil && !r.live() })
	unrecorded := *run
	unrecorded.Container = ""
	for _, r := range []*AgentRun{run, &unrecorded} {
		if said := lab.l.stopContainer(lab.cfg, r); said != "stopped its container "+name {
			t.Errorf("the stop of a run recording container %q says %q, want it stopped %s", r.Container, said, name)
		}
	}
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-stopped")); string(got) != name+"\n"+name+"\n" {
		t.Errorf("docker stop was given %q, not %s twice", got, name)
	}
}

// S-0705: a container run recorded without its container's name, as runs
// before S-0665 were, is stopped by the name the adapter gave it: its
// story's, a defined role's name, not its verb, or a built-in role's verb.
func TestAContainerRunWithoutItsRecordedNameIsNamedAsTheAdapterNamedIt(t *testing.T) {
	const session = "3f2a-b1"
	for _, c := range []struct {
		name string
		run  AgentRun
		want string
	}{
		{"a story's agent", AgentRun{Story: "S-0042"}, "S-0042"},
		{"a defined role on an item", AgentRun{Item: "S-0042", Role: "review", RoleName: "reviewer"}, "reviewer"},
		{"a defined role on no item", AgentRun{Role: "audit", RoleName: "auditor"}, "auditor"},
		{"a defined role in a built-in verb", AgentRun{Role: conventions.RoleCurate, RoleName: "librarian"}, "librarian"},
		{"the curator", AgentRun{Role: conventions.RoleCurate}, conventions.RoleCurate},
		{"the planner", AgentRun{Item: "E-0029"}, conventions.RolePlan},
		{"the analyzer", AgentRun{Focus: "risk"}, conventions.RoleAnalyze},
		{"the orchestrator", AgentRun{}, conventions.RoleOrchestrate},
	} {
		c.run.Session = session
		if got, want := harness.ContainerName(c.run.containerSubject(), c.run.Session), "flai-agent-"+c.want+"-"+session; got != want {
			t.Errorf("%s's container is named %s, want %s", c.name, got, want)
		}
	}
}

// S-0509, ADR-0247 §2: turning the orchestrate action off stops the
// orchestrator's run on the container harness, its docker client signalled
// and its container stopped by the name the adapter gave it, for the role and
// the session, and the stop's journal entry says so.
func TestTheOrchestratorsContainerIsStoppedWhenTheActionIsTurnedOff(t *testing.T) {
	lab := newOrchestrateLab(t)
	containerOn(lab.agentLab)
	lab.roleContainerReady()
	m := filepath.Join(lab.root, manifest.File)
	data, err := os.ReadFile(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(m, append(data, "orchestration:\n  agent:\n    harness: container\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	lab.roleContainerStarted(conventions.RoleOrchestrate, run)
	lab.cfg.Orchestrate = false
	lab.look()
	stopped := lab.ended(run.Session)
	if stopped.Outcome != OutcomeStopped || Alive(run.PID) {
		t.Errorf("not stopped: %+v", stopped)
	}
	name := harness.ContainerName(conventions.RoleOrchestrate, run.Session)
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-stopped")); string(got) != name+"\n" {
		t.Errorf("docker stop was given %q, not %s", got, name)
	}
	j := lab.orchestrateEntries()
	if last := j[len(j)-1]; !strings.Contains(last.Detail, "since orchestrate was turned off") || !strings.Contains(last.Detail, "; stopped its container "+name) {
		t.Errorf("journal: %+v", last)
	}
}

// S-0509, ADR-0247 §2: the operator's stop of an orchestrator on the
// container harness stops its container by its name too, and journals it.
func TestTheOperatorsStopOfTheOrchestratorStopsItsContainer(t *testing.T) {
	lab := newOrchestrateLab(t)
	containerOn(lab.agentLab)
	lab.roleContainerReady()
	m := filepath.Join(lab.root, manifest.File)
	data, err := os.ReadFile(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(m, append(data, "orchestration:\n  agent:\n    harness: container\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	lab.roleContainerStarted(conventions.RoleOrchestrate, run)
	if _, err := OrchestratorStop(lab.o, lab.entry()); err != nil {
		t.Fatal(err)
	}
	name := harness.ContainerName(conventions.RoleOrchestrate, run.Session)
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, "docker-stopped")); string(got) != name+"\n" {
		t.Errorf("docker stop was given %q, not %s", got, name)
	}
	j := lab.orchestrateEntries()
	if last := j[len(j)-1]; !strings.Contains(last.Detail, "on the operator's word") || !strings.Contains(last.Detail, "; stopped its container "+name) {
		t.Errorf("journal: %+v", last)
	}
}

// S-0509, ADR-0247 §5: a project whose strategic agents run on the container
// harness has the host's claude checked for none of them, since the image's
// Claude Code version is the one read and recorded.
func TestStrategicAgentsInContainersAreNoCauseToCheckTheHostsClaude(t *testing.T) {
	root := t.TempDir()
	manifestOn := func(agents string) {
		t.Helper()
		m := filepath.Join(root, manifest.File)
		if err := os.WriteFile(m, []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"+agents), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := manifest.Load(m); err != nil {
			t.Fatal(err)
		}
	}
	manifestOn("agent:\n  harness: container\nplanning:\n  agent:\n    harness: container\norchestration:\n  agent:\n    harness: container\nanalysis:\n  agent:\n    harness: container\ndigest:\n  agent:\n    harness: container\n")
	if _, uses := claudeCodeProject(root); uses {
		t.Error("a project whose agents all run in containers has the host's claude checked")
	}
	manifestOn("agent:\n  harness: container\nplanning:\n  agent:\n    harness: claude-code\n")
	if _, uses := claudeCodeProject(root); !uses {
		t.Error("a planner on the host's claude-code does not have it checked")
	}
}

// ADR-0281 §2: a strategic agent's start on the container harness is refused
// before anything runs when the main checkout's .git/config is missing or is
// a folder, since the container mounts it read-only and docker would make a
// missing one a folder the root user owns; the checkout is left as it was.
func TestAStrategicContainerIsRefusedACheckoutWithNoGitConfig(t *testing.T) {
	for _, c := range []struct {
		name string
		set  func(config string) error
	}{
		{"missing", os.Remove},
		{"a folder", func(config string) error {
			if err := os.Remove(config); err != nil {
				return err
			}
			return os.Mkdir(config, 0o755)
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, _ := roleContainerLab(t)
			config := filepath.Join(lab.root, ".git", "config")
			if err := c.set(config); err != nil {
				t.Fatal(err)
			}
			before, _ := os.Lstat(config)
			id := lab.epic.ID
			if lab.planInContainer(id) {
				t.Fatal("started")
			}
			r := lab.planRun(id)
			if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 {
				t.Fatalf("run: %+v", r)
			}
			for _, s := range []string{config, "is not a file", "run the agent on another harness"} {
				if !strings.Contains(r.Why, s) {
					t.Errorf("the refusal does not say %q: %s", s, r.Why)
				}
			}
			if _, err := os.Stat(filepath.Join(lab.outDir, conventions.RolePlan+".stdin")); err == nil {
				t.Error("the planner's container was run")
			}
			after, _ := os.Lstat(config)
			if (before == nil) != (after == nil) || (after != nil && after.IsDir() != before.IsDir()) {
				t.Errorf("the refusal changed .git/config: %v, then %v", before, after)
			}
		})
	}
}
