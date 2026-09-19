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
	excludes string // git config core.excludesFile, already expanded
	// the push key (S-0062)
	remote     string // git remote get-url --push origin
	keyKind    string // what ssh-keygen -y -P "" says: "", "passphrase", "notkey"
	knownHosts string // what ssh-keygen -F prints for the remote's host
	pushMount  string // source of the push key mount of a running container
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
	if name == "ssh-keygen" {
		switch args[0] {
		case "-y":
			switch f.keyKind {
			case "passphrase":
				return "Load key \"k\": incorrect passphrase supplied to decrypt private key", fmt.Errorf("ssh-keygen: exit status 255")
			case "notkey":
				return "Load key \"k\": error in libcrypto", fmt.Errorf("ssh-keygen: exit status 255")
			}
			return "ssh-ed25519 AAAAC3Nza test", nil
		case "-l":
			return "256 SHA256:fakefingerprint alex at laptop (ED25519)", nil
		case "-F":
			if f.knownHosts == "" {
				return "", fmt.Errorf("ssh-keygen: exit status 1")
			}
			return "# Host " + args[1] + " found: line 3\n" + f.knownHosts + "\n", nil
		}
		return "", nil
	}
	if name == "git" {
		switch {
		case args[0] == "remote" && f.remote != "":
			return f.remote, nil
		case args[0] == "rev-parse":
			return "abc1234", nil
		case args[0] == "config" && f.identity && args[1] == "user.name":
			return "Test User", nil
		case args[0] == "config" && f.identity && args[1] == "user.email":
			return "test@example.com", nil
		case args[0] == "config" && f.excludes != "" && args[len(args)-1] == "core.excludesFile":
			return f.excludes, nil
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
		if strings.Contains(strings.Join(args, " "), pushKeyMountPath) {
			return f.pushMount, nil
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
	a := &app{out: &out, errOut: &errOut, cwd: dir, runner: r, clock: func() time.Time { return time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC) }}
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
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
		"--name flaiover-my-proj",
		"--publish 0.0.0.0:5555:3000",
		"--volume " + root + ":" + root, // at its host path, so worktree links resolve (ADR-0022)
		"--env PROJECT_DIR=" + root,
		"--mount type=bind,source=" + filepath.Join(root, ".flai-cache", "dashboard.token") + ",target=/run/secrets/flaiover_token,readonly",
		"--env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token",
		"--user " + fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q in calls:\n%s", want, joined)
		}
	}
	if !strings.Contains(out, "http://localhost:5555") || !strings.Contains(out, "flai dashboard stop") || !strings.Contains(out, "http://localhost:5555/login#token=") {
		t.Errorf("output: %s", out)
	}
	tokenFile := filepath.Join(root, ".flai-cache", "dashboard.token")
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
	if code != 0 || string(rotated) == string(tokenData) || !strings.Contains(out, "token: "+strings.TrimSpace(string(rotated))) || !strings.Contains(strings.Join(f.calls, "\n"), "docker restart flaiover-my-proj") || !strings.Contains(out, "restarted flaiover-my-proj") {
		t.Errorf("rotate: %d %s", code, out)
	}
	// precedence: flag port beats manifest port; manifest port beat config's 4242
	out, _, _ = runWith(t, root, &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}, "dashboard", "--port", "6000", "--json")
	if !strings.Contains(out, `"url": "http://localhost:6000"`) {
		t.Errorf("flag precedence: %s", out)
	}
	// the host's git identity travels into the container so acceptance can commit
	idr := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}, identity: true}
	_, _, _ = runWith(t, root, idr, "dashboard")
	if j := strings.Join(idr.calls, "\n"); !strings.Contains(j, "--env GIT_COMMITTER_NAME=Test User") || !strings.Contains(j, "--env GIT_AUTHOR_EMAIL=test@example.com") {
		t.Errorf("identity env missing:\n%s", j)
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
	// already running, status, logs, stop
	out, _, _ = runWith(t, root, f, "dashboard")
	if !strings.Contains(out, "already running") {
		t.Errorf("second start: %s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "running at http://localhost:5555 (ghcr.io/bytepunx/flaiover:0.2.0)") {
		t.Errorf("status must report the running container: %s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "logs")
	if !strings.Contains(out, "hello from container") {
		t.Errorf("logs: %s", out)
	}
	out, _, code = runWith(t, root, f, "dashboard", "stop")
	if code != 0 || !strings.Contains(out, "stopped flaiover-my-proj") || f.running["flaiover-my-proj"] {
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

func TestContainerMount(t *testing.T) {
	for _, c := range []struct {
		goos, root, target string
		mirrored           bool
	}{
		{"linux", "/home/a/repo", "/home/a/repo", true},
		{"darwin", "/Users/a/my repo", "/Users/a/my repo", true},
		{"windows", `C:\Users\a\repo`, "/project", false},
		{"linux", "C:/Users/a/repo", "/project", false},
		{"linux", "/home/a/re:po", "/project", false}, // a colon splits docker's --volume
	} {
		target, mirrored := containerMount(c.goos, c.root)
		if target != c.target || mirrored != c.mirrored {
			t.Errorf("containerMount(%s, %s) = %s, %v; want %s, %v", c.goos, c.root, target, mirrored, c.target, c.mirrored)
		}
	}
}

// A host path that cannot be a container path falls back to /project and
// says what that costs and what to do about it (ADR-0022).
func TestDashboardFallsBackWhenTheHostPathCannotBeMirrored(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := filepath.Join(t.TempDir(), "re:po") // a colon would split docker's --volume target
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: colon\nkey: c\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatalf("run: %s", errOut)
	}
	joined := strings.Join(f.calls, "\n")
	for _, want := range []string{"--volume " + root + ":/project", "--env PROJECT_DIR=/project"} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q in calls:\n%s", want, joined)
		}
	}
	for _, want := range []string{`"level":"WARN"`, `"mount":"/project"`, `"effect":`, "flai accept", "worktrees.relative_paths"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("warning should carry %s:\n%s", want, errOut)
		}
	}
	if !strings.Contains(out, "mounted read-write at /project") {
		t.Errorf("output should name the mount: %s", out)
	}
}

// The container gets the host's global git excludes, so what git ignores on
// the host is not uncommitted in the dashboard (I-0019, S-0051).
func TestDashboardPassesTheGlobalGitExcludes(t *testing.T) {
	run := func(t *testing.T, f *fakeRunner) string {
		t.Helper()
		t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
		if _, errOut, code := runWith(t, tempProject(t), f, "dashboard"); code != 0 {
			t.Fatalf("run: %s", errOut)
		}
		return strings.Join(f.calls, "\n")
	}
	env := "--env GIT_CONFIG_COUNT=1 --env GIT_CONFIG_KEY_0=core.excludesFile --env GIT_CONFIG_VALUE_0=/run/flaiover/gitignore"

	t.Run("the configured file", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "unused"))
		file := filepath.Join(t.TempDir(), "my-ignore")
		_ = os.WriteFile(file, []byte(".claude/settings.local.json\n"), 0o644)
		calls := run(t, &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, excludes: file})
		for _, want := range []string{"--mount type=bind,source=" + file + ",target=/run/flaiover/gitignore,readonly", env} {
			if !strings.Contains(calls, want) {
				t.Errorf("missing %q in calls:\n%s", want, calls)
			}
		}
	})
	t.Run("git's default location when none is configured", func(t *testing.T) {
		xdg := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", xdg)
		file := filepath.Join(xdg, "git", "ignore")
		_ = os.MkdirAll(filepath.Dir(file), 0o755)
		_ = os.WriteFile(file, []byte("x\n"), 0o644)
		calls := run(t, &fakeRunner{images: map[string]bool{}, running: map[string]bool{}})
		if !strings.Contains(calls, "source="+file+",target=/run/flaiover/gitignore,readonly") || !strings.Contains(calls, env) {
			t.Errorf("default excludes file not passed:\n%s", calls)
		}
	})
	t.Run("nothing when there is no such file", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		calls := run(t, &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, excludes: filepath.Join(t.TempDir(), "gone")})
		if strings.Contains(calls, "GIT_CONFIG_COUNT") || strings.Contains(calls, "/run/flaiover/gitignore") {
			t.Errorf("no excludes file, nothing to pass:\n%s", calls)
		}
	})
}
