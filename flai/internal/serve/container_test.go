//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// containerLab is providerLab with a fake docker first on PATH, the
// container harness's login set, and the project's HTTP MCP server running,
// as flai host keeps it, at a loopback URL with a token. The fake answers
// image inspect, unless a file no-image is in the lab's output folder; a run
// with --entrypoint, as claude --version in the image; and stop, recording
// the name, failing when a file stop-fails is there. Any other run is the
// agent's: it records its arguments and environment in <story>.txt once it
// has read its standard input whole into <story>.stdin, and holds as the
// stub does.
func containerLab(t *testing.T) (*agentLab, string) {
	t.Helper()
	lab := providerLab(t)
	bin := t.TempDir()
	out := lab.outDir
	fake := "#!/bin/sh\n" +
		"out=\"" + out + "\"\n" +
		"case \"$1\" in\n" +
		"image) [ -f \"$out/no-image\" ] && { echo \"Error: No such image: $5\"; exit 1; }; echo sha256:fake; exit 0 ;;\n" +
		"stop) echo \"$2\" >> \"$out/docker-stopped\"; [ -f \"$out/stop-fails\" ] && { echo \"Error response from daemon: cannot stop\"; exit 1; }; exit 0 ;;\n" +
		"esac\n" +
		"for a in \"$@\"; do [ \"$a\" = --entrypoint ] && { echo \"2.1.312 (Claude Code)\"; exit 0; }; done\n" +
		"trap 'touch \"$out/ended-$$\"; exit 143' TERM\n" +
		"cat > \"$out/$FLAI_STORY.stdin\"\n" +
		"{ echo \"args: $*\"; echo \"story: $FLAI_STORY\"; echo \"agent: $FLAI_AGENT\"; echo \"session: $FLAI_SESSION\"; echo \"started_by: $FLAI_STARTED_BY\"; echo \"mcp: ${FLAI_MCP_URL-unset} ${FLAI_MCP_TOKEN-unset}\"; echo \"key: ${S0350_GATEWAY_KEY-unset}\"; echo \"config: ${FLAI_CONFIG-unset}\"; } > \"$out/$FLAI_STORY.txt\"\n" +
		"while [ -f \"$out/hold\" ] && [ ! -f \"$out/release-$FLAI_STORY\" ]; do sleep 0.05; done\n" +
		"touch \"$out/ended-$$\"\n"
	if err := os.WriteFile(filepath.Join(bin, "docker"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	login := t.TempDir()
	lab.cfg.Harnesses[harness.Container] = harness.Host{Login: login}
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
	return lab, token
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
// image's Claude Code version is read and recorded on the run, and a version
// this host has no passed check of is warned of. The cache folder is made the
// user's alone, with the link to the login.
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
// image not on this host, no login for an inner harness that needs one, and
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
			says: []string{"flai-agent-claude-code:local", "No such image", "make agent-images", "--image claude-code=<image>"}},
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
