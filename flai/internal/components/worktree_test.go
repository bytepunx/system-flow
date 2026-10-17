package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// fakePackageManager answers as a package manager would without running
// one: each Run is recorded, then calls install with the folder it ran in,
// when set, and answers stdout, stderr, and err.
type fakePackageManager struct {
	missing bool
	install func(dir string)
	stdout  string
	stderr  string
	err     error
	block   bool

	mu    sync.Mutex
	calls []fakeCall
}

func (f *fakePackageManager) LookPath(name string) (string, error) {
	if f.missing {
		return "", exec.ErrNotFound
	}
	return "/usr/bin/" + name, nil
}

func (f *fakePackageManager) Run(ctx context.Context, dir, name, input string, args ...string) ([]byte, []byte, error) {
	f.mu.Lock()
	f.calls = append(f.calls, fakeCall{dir: dir, input: input, args: append([]string{name}, args...)})
	f.mu.Unlock()
	if f.block {
		<-ctx.Done()
		return []byte("resolving"), nil, ctx.Err()
	}
	if f.install != nil {
		f.install(dir)
	}
	return []byte(f.stdout), []byte(f.stderr), f.err
}

// ran is each command Run was given, with the folder it ran in first.
func (f *fakePackageManager) ran() [][]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([][]string, len(f.calls))
	for i, c := range f.calls {
		out[i] = append([]string{c.dir}, c.args...)
	}
	return out
}

// installVite is a fake install that puts a vite launcher in the
// node_modules/.bin of the sub-project in dir.
func installVite(t *testing.T, dir string) func(string) {
	t.Helper()
	return func(string) {
		bin := filepath.Join(dir, "node_modules", ".bin")
		if err := os.MkdirAll(bin, 0o755); err != nil {
			t.Error(err)
			return
		}
		if err := os.WriteFile(filepath.Join(bin, viteBin), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Error(err)
		}
	}
}

// checkout is a checkout's root holding a sub-project at packages/ui, and
// the sub-project's folder.
func checkout(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "packages", "ui")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return root, dir
}

// write writes an empty file at each path.
func write(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestServeWorktreeInstallsTheDependenciesThenServes(t *testing.T) {
	serve, state := serveFixture(t, "serve")
	dir := filepath.Join(serve.Root, serve.Sub.Path)
	stash := filepath.Join(state, "node_modules")
	if err := os.Rename(filepath.Join(dir, "node_modules"), stash); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, "package-lock.json"))
	pm := &fakePackageManager{install: func(at string) {
		if err := os.Rename(stash, filepath.Join(at, "node_modules")); err != nil {
			t.Error(err)
		}
	}}
	var logged bytes.Buffer
	s, err := ServeWorktree(context.Background(), WorktreeOptions{
		ServeOptions: serve,
		Runner:       pm,
		Logger:       slog.New(slog.NewTextHandler(&logged, nil)),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Stop() })
	if got, want := pm.ran(), [][]string{{dir, "npm", "ci"}}; !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("ran %q, want %q", got, want)
	}
	resp, err := http.Get(s.URL())
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s: %d", s.URL(), resp.StatusCode)
	}
	for _, want := range []string{"level=INFO", "component=components", "dir=" + dir, "lock_file=package-lock.json", `command="npm ci"`, "duration="} {
		if !strings.Contains(logged.String(), want) {
			t.Errorf("the log lacks %s:\n%s", want, logged.String())
		}
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	waitGone(t, s.PID())
}

func TestInstallRunsThePackageManagerOfTheLockFileInItsFolder(t *testing.T) {
	for _, tc := range []struct {
		name  string
		sub   []string // lock files in the sub-project's folder
		root  []string // lock files in the checkout's root
		want  []string
		inSub bool
	}{
		{"pnpm", []string{"pnpm-lock.yaml"}, nil, []string{"pnpm", "install", "--frozen-lockfile"}, true},
		{"npm", []string{"package-lock.json"}, nil, []string{"npm", "ci"}, true},
		{"yarn", []string{"yarn.lock"}, nil, []string{"yarn", "install", "--frozen-lockfile"}, true},
		{"a workspace's at the root", nil, []string{"pnpm-lock.yaml"}, []string{"pnpm", "install", "--frozen-lockfile"}, false},
		{"the sub-project's before the root's", []string{"yarn.lock"}, []string{"pnpm-lock.yaml"}, []string{"yarn", "install", "--frozen-lockfile"}, true},
		{"pnpm before npm in one folder", nil, []string{"package-lock.json", "pnpm-lock.yaml"}, []string{"pnpm", "install", "--frozen-lockfile"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, dir := checkout(t)
			for _, l := range tc.sub {
				write(t, filepath.Join(dir, l))
			}
			for _, l := range tc.root {
				write(t, filepath.Join(root, l))
			}
			pm := &fakePackageManager{install: installVite(t, dir)}
			ran, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: pm})
			if err != nil || !ran {
				t.Fatalf("Install: %v, %v; want it to run", ran, err)
			}
			at := root
			if tc.inSub {
				at = dir
			}
			if got, want := pm.ran(), [][]string{append([]string{at}, tc.want...)}; !slices.EqualFunc(got, want, slices.Equal) {
				t.Errorf("ran %q, want %q", got, want)
			}
		})
	}
}

func TestInstallSkipsASubProjectWithViteInstalled(t *testing.T) {
	root, dir := checkout(t)
	write(t, filepath.Join(dir, "pnpm-lock.yaml"))
	installVite(t, dir)("")
	pm := &fakePackageManager{}
	ran, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: pm})
	if err != nil || ran || len(pm.ran()) != 0 {
		t.Errorf("Install: %v, %v, ran %q; want nothing run", ran, err, pm.ran())
	}
}

func TestServeWorktreeReturnsAFailedInstallWithItsLastLines(t *testing.T) {
	root, dir := checkout(t)
	write(t, filepath.Join(dir, "pnpm-lock.yaml"))
	var stdout strings.Builder
	for i := 1; i <= 25; i++ {
		fmt.Fprintf(&stdout, "line %d\n", i)
	}
	pm := &fakePackageManager{stdout: stdout.String(), stderr: "ERR_PNPM_OUTDATED_LOCKFILE  Cannot install with \"frozen-lockfile\"", err: errors.New("exit status 1")}
	serve := ServeOptions{Root: root, Project: "proj", Sub: manifest.Project{Name: "ui", Path: "packages/ui", Kind: manifest.ProjectKindSvelteKit}}
	_, err := ServeWorktree(context.Background(), WorktreeOptions{ServeOptions: serve, Runner: pm})
	if err == nil {
		t.Fatal("ServeWorktree started a dev server after a failed install")
	}
	msg := err.Error()
	for _, want := range []string{"cannot serve sub-project ui's components: ", "pnpm install --frozen-lockfile in " + dir + " failed: exit status 1; fix what its last lines say", "\nline 25\n", "ERR_PNPM_OUTDATED_LOCKFILE"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error lacks %q:\n%s", want, msg)
		}
	}
	if strings.Contains(msg, "\nline 6\n") || !strings.Contains(msg, "\nline 7\n") {
		t.Errorf("error keeps other than the last %d lines:\n%s", tailLines, msg)
	}
}

func TestInstallWithoutALockFileSaysWhatToDo(t *testing.T) {
	root, dir := checkout(t)
	pm := &fakePackageManager{}
	_, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: pm})
	if !errors.Is(err, ErrNoLockFile) || len(pm.ran()) != 0 {
		t.Fatalf("error %v, ran %q; want ErrNoLockFile and nothing run", err, pm.ran())
	}
	for _, want := range []string{dir, "pnpm-lock.yaml, package-lock.json, yarn.lock", "up to " + root, "commit the sub-project's lock file"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error lacks %q: %v", want, err)
		}
	}
}

func TestInstallLooksNoHigherThanTheCheckoutsRoot(t *testing.T) {
	above := t.TempDir()
	write(t, filepath.Join(above, "package-lock.json"))
	root := filepath.Join(above, "worktree")
	dir := filepath.Join(root, "ui")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := Install(context.Background(), InstallOptions{Root: root + string(filepath.Separator), Dir: dir, Runner: &fakePackageManager{}})
	if !errors.Is(err, ErrNoLockFile) {
		t.Errorf("error %v, want ErrNoLockFile, with the lock file above the root not used", err)
	}
}

func TestInstallSaysWhenThePackageManagerIsMissing(t *testing.T) {
	root, dir := checkout(t)
	write(t, filepath.Join(dir, "yarn.lock"))
	_, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: &fakePackageManager{missing: true}})
	if err == nil || !strings.Contains(err.Error(), "yarn is not on PATH, and its lock file "+filepath.Join(dir, "yarn.lock")+" needs it; install yarn") {
		t.Errorf("error %v, want one saying to install yarn", err)
	}
}

func TestInstallStopsAPackageManagerThatTakesTooLong(t *testing.T) {
	root, dir := checkout(t)
	write(t, filepath.Join(dir, "package-lock.json"))
	_, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: &fakePackageManager{block: true}, Timeout: 20 * time.Millisecond})
	if err == nil || !strings.Contains(err.Error(), "npm ci in "+dir+" did not install the sub-project's dependencies within 20ms") || !strings.Contains(err.Error(), "\nresolving") {
		t.Errorf("error %v, want the timeout with the last lines", err)
	}
}

func TestInstallSaysWhenViteIsStillMissing(t *testing.T) {
	root, dir := checkout(t)
	write(t, filepath.Join(root, "package-lock.json"))
	_, err := Install(context.Background(), InstallOptions{Root: root, Dir: dir, Runner: &fakePackageManager{}})
	if !errors.Is(err, ErrViteMissing) || !strings.Contains(err.Error(), "after npm ci in "+root) || !strings.Contains(err.Error(), "add vite to the sub-project's own devDependencies") {
		t.Errorf("error %v, want ErrViteMissing naming the install, with the fix", err)
	}
}

func TestInstallRefuses(t *testing.T) {
	root, dir := checkout(t)
	for _, tc := range []struct {
		name string
		opts InstallOptions
		err  string
	}{
		{"relative folder", InstallOptions{Root: root, Dir: "packages/ui"}, "as absolute paths"},
		{"relative root", InstallOptions{Root: "wt", Dir: dir}, "as absolute paths"},
		{"outside the root", InstallOptions{Root: dir, Dir: root}, "it is not inside the checkout at " + dir},
		{"missing folder", InstallOptions{Root: root, Dir: filepath.Join(root, "web")}, "its folder is not there"},
	} {
		pm := &fakePackageManager{}
		tc.opts.Runner = pm
		if _, err := Install(context.Background(), tc.opts); err == nil || !strings.Contains(err.Error(), tc.err) || len(pm.ran()) != 0 {
			t.Errorf("%s: error %v, ran %q; want one containing %q, with nothing run", tc.name, err, pm.ran(), tc.err)
		}
	}
}

func TestServeWorktreeLeavesWhatServeRefusesToServe(t *testing.T) {
	root, _ := checkout(t)
	write(t, filepath.Join(root, "pnpm-lock.yaml"))
	pm := &fakePackageManager{}
	_, err := ServeWorktree(context.Background(), WorktreeOptions{ServeOptions: ServeOptions{Root: root, Project: "proj", Sub: pickProjects[1]}, Runner: pm})
	if err == nil || !strings.Contains(err.Error(), "its kind go is not sveltekit or svelte") || len(pm.ran()) != 0 {
		t.Errorf("error %v, ran %q; want Serve's refusal, with nothing installed", err, pm.ran())
	}
}
