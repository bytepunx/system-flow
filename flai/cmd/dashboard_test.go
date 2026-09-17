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
	calls   []string
	missing bool            // docker not on PATH
	images  map[string]bool // images present
	running map[string]bool // containers running
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if name == "docker" && f.missing {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeRunner) Run(dir, name string, args ...string) (string, error) {
	f.calls = append(f.calls, name+" "+strings.Join(args, " "))
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
		f.images[args[2]] = true
		return args[2], nil
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
		"--publish 127.0.0.1:5555:3000",
		"--volume " + root + ":/project",
		"--env PROJECT_DIR=/project",
		"--user " + fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q in calls:\n%s", want, joined)
		}
	}
	if !strings.Contains(out, "http://localhost:5555") || !strings.Contains(out, "flai dashboard stop") {
		t.Errorf("output: %s", out)
	}
	// precedence: flag port beats manifest port; manifest port beat config's 4242
	out, _, _ = runWith(t, root, &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}, "dashboard", "--port", "6000", "--json")
	if !strings.Contains(out, `"url": "http://localhost:6000"`) {
		t.Errorf("flag precedence: %s", out)
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
