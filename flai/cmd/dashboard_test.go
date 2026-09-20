package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// fakeRunner records docker calls and answers from a script.
type fakeRunner struct {
	calls    []string
	missing  bool            // docker not on PATH
	images   map[string]bool // images present
	running  map[string]bool // containers running
	private  map[string]bool // images that need a login to pull
	loggedIn bool
	noToken  bool   // gh has no token
	identity bool   // git config has a user
	mounts   string // mount destinations of a running container, one per line
}

func (f *fakeRunner) RunInput(dir, name, input string, args ...string) (string, error) {
	if name == "docker" && args[0] == "login" {
		f.calls = append(f.calls, "docker "+strings.Join(args, " ")+" <"+strings.TrimSpace(input)+">")
		if strings.TrimSpace(input) == "" {
			return "", fmt.Errorf("docker login: empty password")
		}
		f.loggedIn = true
		return "Login Succeeded", nil
	}
	return f.Run(dir, name, args...)
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if name == "docker" && f.missing {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeRunner) Run(dir, name string, args ...string) (string, error) {
	f.calls = append(f.calls, name+" "+strings.Join(args, " "))
	if name == "gh" {
		switch strings.Join(args, " ") {
		case "auth token":
			if f.noToken {
				return "", fmt.Errorf("gh auth token: not logged in")
			}
			return "ghp_fake", nil
		case "api user --jq .login":
			return "tester", nil
		}
		return "", nil
	}
	if name == "git" {
		switch {
		case args[0] == "rev-parse":
			return "abc1234", nil
		case args[0] == "config" && f.identity && args[1] == "user.name":
			return "Test User", nil
		case args[0] == "config" && f.identity && args[1] == "user.email":
			return "test@example.com", nil
		}
		return "", fmt.Errorf("git: no tags")
	}
	if name != "docker" {
		return "", nil
	}
	switch args[0] {
	case "image":
		if f.images[args[2]] {
			return "[]", nil
		}
		return "", fmt.Errorf("docker image inspect %s: exit status 1\nError: No such image", args[2])
	case "pull":
		if f.private[args[2]] && !f.loggedIn {
			return "", fmt.Errorf("docker pull --quiet %s: exit status 1\nError response from daemon: error from registry: unauthorized", args[2])
		}
		f.images[args[2]] = true
		return args[2], nil
	case "build":
		f.images[args[4]] = true
		return "sha256:built", nil
	case "run":
		f.running[args[4]] = true
		return "0123456789abcdef", nil
	case "ps":
		name := strings.TrimSuffix(strings.TrimPrefix(args[len(args)-1], "name=^/"), "$")
		if f.running[name] {
			return "0123456789ab", nil
		}
		return "", nil
	case "inspect":
		if strings.Contains(strings.Join(args, " "), ".Mounts") {
			return f.mounts, nil
		}
		return "ghcr.io/bytepunx/flaiover:0.2.0 5555", nil
	case "stop":
		delete(f.running, args[1])
		return args[1], nil
	case "logs":
		return "hello from container", nil
	}
	return "", nil
}

func runWith(t *testing.T, dir string, r *fakeRunner, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: dir, runner: r, clock: func() time.Time { return time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC) },
		// No test starts a real flai serve: the test binary is not flai.
		serveStarter: func() (serve.Status, bool, error) { return serve.Status{PID: 4242}, true, nil }}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		var ee *exitError
		if errors.As(err, &ee) {
			code = ee.code
		} else {
			a.fail(err)
			code = 1
		}
	}
	return out.String(), errOut.String(), code
}

func TestDashboardLifecycle(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	serveDir := string(serve.DirFor(cfg))
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: My Proj\nkey: m\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n"), 0o644)

	missing := &fakeRunner{missing: true, images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, missing, "dashboard")
	if code == 0 || !strings.Contains(errOut, "docker is required") || !strings.Contains(errOut, "Install Docker") {
		t.Fatalf("missing docker: %d %s", code, errOut)
	}

	_ = out
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code = runWith(t, root, f, "dashboard", "--tag", "0.2.0")
	if code != 0 {
		t.Fatalf("run: %s", errOut)
	}
	joined := strings.Join(f.calls, "\n")
	for _, want := range []string{
		"docker image inspect ghcr.io/bytepunx/flaiover:0.2.0",
		"docker pull --quiet ghcr.io/bytepunx/flaiover:0.2.0",
		"--name flaiover",
		"--publish 0.0.0.0:5555:3000",
		"--mount type=bind,source=" + filepath.Join(serveDir, "dashboard.token") + ",target=/run/secrets/flaiover_token,readonly",
		"--env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token",
		"--user " + fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q in calls:\n%s", want, joined)
		}
	}
	if !strings.Contains(out, "http://localhost:5555") || !strings.Contains(out, "flai dashboard stop") || !strings.Contains(out, "http://localhost:5555/login#token=") || !strings.Contains(out, "no file of any project") {
		t.Errorf("output: %s", out)
	}
	// ADR-0031: a port and two secrets, and nothing of the project or of the
	// person who started it.
	var run string
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker run ") {
			run = c
		}
	}
	for _, gone := range []string{"--volume", root + ":", "PROJECT_DIR", "GIT_", ".git/", "push_key", "known_hosts", "/etc/passwd", "gitignore"} {
		if strings.Contains(run, gone) {
			t.Errorf("%q has no place in the container's arguments:\n%s", gone, run)
		}
	}
	if n := strings.Count(run, "--mount "); n != 2 || strings.Count(run, "target=/run/secrets/") != 2 || strings.Count(run, ",readonly") != 2 {
		t.Errorf("the only mounts are the two secrets, read-only:\n%s", run)
	}
	tokenFile := filepath.Join(serveDir, "dashboard.token")
	tokenData, err := os.ReadFile(tokenFile)
	if err != nil || len(strings.TrimSpace(string(tokenData))) < 40 {
		t.Fatalf("token file: %v %q", err, tokenData)
	}
	if st, _ := os.Stat(tokenFile); st.Mode().Perm() != 0o600 {
		t.Errorf("token file mode %v", st.Mode().Perm())
	}
	if strings.Contains(errOut, strings.TrimSpace(string(tokenData))) {
		t.Error("token must not be logged")
	}
	// token command: stable, then rotated with a restart of the running container
	out, _, code = runWith(t, root, f, "dashboard", "token")
	if code != 0 || !strings.Contains(out, "token: "+strings.TrimSpace(string(tokenData))) || !strings.Contains(out, "login: http://localhost:5555/login#token=") {
		t.Errorf("token: %d %s", code, out)
	}
	out, _, code = runWith(t, root, f, "dashboard", "token", "--rotate")
	rotated, _ := os.ReadFile(tokenFile)
	if code != 0 || string(rotated) == string(tokenData) || !strings.Contains(out, "token: "+strings.TrimSpace(string(rotated))) || !strings.Contains(strings.Join(f.calls, "\n"), "docker restart flaiover") || !strings.Contains(out, "restarted flaiover") {
		t.Errorf("rotate: %d %s", code, out)
	}
	// precedence: flag port beats manifest port; manifest port beat config's 4242
	out, _, _ = runWith(t, root, &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}, "dashboard", "--port", "6000", "--json")
	if !strings.Contains(out, `"url": "http://localhost:6000"`) {
		t.Errorf("flag precedence: %s", out)
	}
	// the host's git identity stays on the host: flai serve commits there (S-0075)
	idr := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}, identity: true}
	_, _, _ = runWith(t, root, idr, "dashboard")
	if j := strings.Join(idr.calls, "\n"); strings.Contains(j, "Test User") || strings.Contains(j, "git config") {
		t.Errorf("the identity is not read, let alone passed:\n%s", j)
	}
	// --bind restricts the published address; the manifest can set it too
	bnd := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}
	out, _, _ = runWith(t, root, bnd, "dashboard", "--bind", "127.0.0.1")
	if !strings.Contains(strings.Join(bnd.calls, "\n"), "--publish 127.0.0.1:5555:3000") || !strings.Contains(out, "reachable from this host only") {
		t.Errorf("--bind: %s\n%s", out, strings.Join(bnd.calls, "\n"))
	}
	// present image is not pulled unless --pull
	g := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}
	_, _, _ = runWith(t, root, g, "dashboard")
	if strings.Contains(strings.Join(g.calls, "\n"), "docker pull") {
		t.Error("pulled although present")
	}
	h := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}
	_, _, _ = runWith(t, root, h, "dashboard", "--pull")
	if !strings.Contains(strings.Join(h.calls, "\n"), "docker pull --quiet ghcr.io/bytepunx/flaiover:latest") {
		t.Error("--pull should pull")
	}
	// already running for this same project: still registers, says so, starts no second container
	out, _, _ = runWith(t, root, f, "dashboard")
	if !strings.Contains(out, "already runs") || !strings.Contains(out, "now also serves") || strings.Count(strings.Join(f.calls, "\n"), "docker run ") != 1 {
		t.Errorf("second start: %s\n%s", out, strings.Join(f.calls, "\n"))
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "running at http://localhost:5555 (ghcr.io/bytepunx/flaiover:0.2.0)") {
		t.Errorf("status must report the running container: %s", out)
	}
	if strings.Contains(out, "older flai") || strings.Contains(out, "push") || strings.Contains(out, "read-only") {
		t.Errorf("status says nothing of keys or git paths any more: %s", out)
	}
	// a container an older flai started still has the project mounted: say so
	f.mounts = "/run/secrets/flaiover_token\n" + root + "\n" + root + "/.git/hooks\n"
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "started by an older flai") || !strings.Contains(out, root) || !strings.Contains(out, "flai dashboard stop") {
		t.Errorf("a stale container: %s", out)
	}
	js, _, _ := runWith(t, root, f, "dashboard", "status", "--json")
	if !strings.Contains(js, `"stale_mounts"`) || strings.Contains(js, "git_read_only") || strings.Contains(js, "push_key") {
		t.Errorf("status --json: %s", js)
	}
	f.mounts = ""
	out, _, _ = runWith(t, root, f, "dashboard", "logs")
	if !strings.Contains(out, "hello from container") {
		t.Errorf("logs: %s", out)
	}
	out, _, code = runWith(t, root, f, "dashboard", "stop")
	if code != 0 || !strings.Contains(out, "stopped flaiover") || f.running["flaiover"] {
		t.Errorf("stop: %d %s", code, out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "stop")
	if !strings.Contains(out, "not running") {
		t.Errorf("stop twice: %s", out)
	}
}

func TestDashboardPrivateRegistryAndBuild(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	root := tempProject(t)
	img := "ghcr.io/bytepunx/flaiover:latest"

	// Unauthorized pull, token from gh: login and retry.
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{img: true}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	joined := strings.Join(f.calls, "\n")
	if code != 0 || !strings.Contains(joined, "docker login ghcr.io --username tester --password-stdin <ghp_fake>") || strings.Count(joined, "docker pull --quiet "+img) != 2 || !strings.Contains(out, "http://localhost:4242") {
		t.Fatalf("login retry: %d %s %s\n%s", code, out, errOut, joined)
	}
	if strings.Contains(errOut, "ghp_fake") {
		t.Error("token must not be logged")
	}

	// No token anywhere: the error names the fix.
	g := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{img: true}, noToken: true}
	_, errOut, code = runWith(t, root, g, "dashboard")
	if code == 0 || !strings.Contains(errOut, "read:packages") || strings.Contains(strings.Join(g.calls, "\n"), "docker login") {
		t.Errorf("no token: %d %s", code, errOut)
	}

	// Docker Hub images get no login attempt.
	h := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{"nginx:latest": true}}
	_, errOut, code = runWith(t, root, h, "dashboard", "--image", "nginx")
	if code == 0 || strings.Contains(strings.Join(h.calls, "\n"), "docker login") {
		t.Errorf("hub image: %d %s", code, errOut)
	}

	// --build outside the monorepo refuses; inside it builds and runs flaiover:local.
	b := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code = runWith(t, root, b, "dashboard", "--build")
	if code == 0 || !strings.Contains(errOut, "flaiover/Dockerfile") {
		t.Errorf("build without dockerfile: %d %s", code, errOut)
	}
	_ = os.MkdirAll(filepath.Join(root, "flaiover"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "flaiover", "Dockerfile"), []byte("FROM scratch\n"), 0o644)
	out, errOut, code = runWith(t, root, b, "dashboard", "--build")
	joined = strings.Join(b.calls, "\n")
	if code != 0 || !strings.Contains(joined, "docker build -f "+filepath.Join(root, "flaiover", "Dockerfile")+" -t flaiover:local --build-arg FLAI_VERSION=dev --build-arg FLAI_COMMIT=abc1234") || !strings.Contains(joined, "flaiover:local") || strings.Contains(joined, "docker pull") || !strings.Contains(out, "image flaiover:local") {
		t.Errorf("build: %d %s %s\n%s", code, out, errOut, joined)
	}
}

// S-0077: the push key is retired. Asked for, it is named, explained, and
// ignored, and the dashboard starts.
func TestDashboardPushKeyIsRetired(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	img := map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}
	f := &fakeRunner{images: img, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "--push-key", "/home/me/.ssh/deploy")
	if code != 0 || !strings.Contains(out, "--push-key: retired and ignored") || !strings.Contains(out, "flai push --pending") || !strings.Contains(out, "ADR-0031") {
		t.Fatalf("flag: %d %s %s", code, out, errOut)
	}
	if j := strings.Join(f.calls, "\n"); strings.Contains(j, "deploy") || strings.Contains(j, "ssh-keygen") {
		t.Errorf("the key is not looked at or passed:\n%s", j)
	}
	if _, _, code := runWith(t, root, f, "config", "set", "dashboard.push_key", "/home/me/.ssh/deploy"); code != 0 {
		t.Fatal("an old config still loads, and the key can still be set, so it can be cleared")
	}
	out, _, _ = runWith(t, root, &fakeRunner{images: img, running: map[string]bool{}}, "dashboard")
	if !strings.Contains(out, "dashboard.push_key: retired and ignored") {
		t.Errorf("config: %s", out)
	}
	js, _, _ := runWith(t, root, &fakeRunner{images: img, running: map[string]bool{}}, "dashboard", "--json")
	if !strings.Contains(js, `"retired"`) || strings.Contains(js, `"push_key"`) || strings.Contains(js, `"mount"`) {
		t.Errorf("--json: %s", js)
	}
	if help, _, _ := runWith(t, root, f, "dashboard", "--help"); strings.Contains(help, "push-key") {
		t.Error("the retired flags are hidden")
	}
}
