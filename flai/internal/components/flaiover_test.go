//go:build !windows

package components

// The integration test of this repository's own component library (S-0444):
// Serve starts flaiover's own Vite dev server, with the components key
// system-flow.yaml gives it, on a copy of flaiover whose node_modules links
// to the installed one, and the test asks it for the setup stylesheet and
// for each component the catalogue lists, as the harness imports them. A
// component that does not compile, or imports what does not resolve, is
// answered with Vite's error, and fails the test. It needs node and
// flaiover's dependencies installed, so it skips without them, and under
// -short. The build tag keeps it off Windows, as the harness's integration
// test is kept off.

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

func TestEveryListedFlaioverComponentCompilesInTheHarness(t *testing.T) {
	if testing.Short() {
		t.Skip("starts flaiover's own Vite dev server; the integration tier runs it")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not on PATH, and flaiover's vite runs on it")
	}
	modules := flaioverModules(t)
	repo, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Load(filepath.Join(repo, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	sub, err := PickSubProject(m.Projects, "flaiover")
	if err != nil {
		t.Fatal(err)
	}
	settings, _ := sub.ComponentSettings()
	if settings.Setup == "" {
		t.Fatalf("system-flow.yaml gives flaiover no components setup, and its components need the app's stylesheet")
	}

	// A copy, so that the dev server's optimizer cache and SvelteKit's
	// generated files land in it, not in the repository.
	root := t.TempDir()
	dir := filepath.Join(root, filepath.FromSlash(sub.Path))
	copySource(t, filepath.Join(repo, filepath.FromSlash(sub.Path)), dir)
	linkModules(t, filepath.Join(dir, "node_modules"), modules)

	comps, err := Find(root, sub.Path, settings.Include, settings.Exclude)
	if err != nil {
		t.Fatal(err)
	}
	if len(comps) == 0 {
		t.Fatalf("the catalogue of flaiover lists no component under %v", settings.Include)
	}

	s, err := Serve(context.Background(), ServeOptions{Root: root, Project: m.Name, Sub: sub, Timeout: serveTimeout})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Stop(); err != nil {
			t.Error(err)
		}
		waitGone(t, s.PID())
		waitGroupGone(t, s.PID())
	})
	vite := func() string { return "\nvite's last lines:\n" + string(s.out.last()) }
	base := strings.TrimSuffix(s.URL(), HarnessPath)
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	t.Cleanup(client.CloseIdleConnections)

	if status, body := fetch(t, client, base+"/"+settings.Setup); status != http.StatusOK {
		t.Fatalf("the setup %s does not compile: %d\n%s%s", settings.Setup, status, body, vite())
	}
	for _, c := range comps {
		if status, body := fetch(t, client, base+"/"+c.Path); status != http.StatusOK {
			t.Errorf("%s does not compile for the harness: %d\n%s", c.Path, status, body)
		}
	}

	// The control: a component that does not compile is answered with an
	// error, so that the checks above can fail.
	broken := "src/lib/components/FlaiBroken.svelte"
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(broken)), []byte("<p>{#if open}</p>\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if status, body := fetch(t, client, base+"/"+broken); status == http.StatusOK {
		t.Errorf("a component that does not compile was answered 200:\n%s%s", body, vite())
	}
}

// fetch is the status and body of a GET of u, failing the test only when
// the request cannot be made.
func fetch(t *testing.T, client *http.Client, u string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", u, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET %s: %v", u, err)
	}
	return resp.StatusCode, string(body)
}

// copySource copies the sub-project at src to dst, leaving out what an
// install or a build makes: node_modules, build, and .svelte-kit.
func copySource(t *testing.T, src, dst string) {
	t.Helper()
	skip := []string{"node_modules", "build", ".svelte-kit"}
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		if d.IsDir() && slices.Contains(skip, d.Name()) {
			return filepath.SkipDir
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
