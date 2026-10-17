//go:build !windows

package components

// The integration test of the component harness (S-0438, ADR-0163): Serve
// starts a real Vite dev server, the sub-project's own, on the fixture
// sub-project testdata/harness, whose node_modules links to flaiover's, and
// the test reads the hot updates Vite sends over its HMR WebSocket as the
// fixture's component and setup stylesheet change. It needs node and
// flaiover's dependencies installed, so it skips without them, and under
// -short. The build tag keeps it off Windows: the fixture links
// node_modules with symlinks, and the test checks vite's process group.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// serveTimeout bounds Serve in the integration test: the first run
// optimizes the dependencies, which can be slow.
const serveTimeout = 2 * time.Minute

// updateTimeout bounds the wait for each hot update, and quietWindow is how
// long the test then waits to see that no full reload follows it.
const (
	updateTimeout = 30 * time.Second
	quietWindow   = 2 * time.Second
)

func TestServeCarriesHotUpdatesOfTheComponentAndTheStylesheet(t *testing.T) {
	if testing.Short() {
		t.Skip("starts the fixture sub-project's own Vite dev server; the integration tier runs it")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not on PATH, and the fixture sub-project's vite runs on it")
	}
	modules := flaioverModules(t)
	before := moduleNames(t, modules)
	t.Cleanup(func() {
		if after := moduleNames(t, modules); !slices.Equal(before, after) {
			t.Errorf("the test changed %s: its entries were\n%v\nand are\n%v", modules, before, after)
		}
	})

	root := t.TempDir()
	dir := filepath.Join(root, "ui")
	if err := os.CopyFS(dir, os.DirFS(filepath.Join("testdata", "harness"))); err != nil {
		t.Fatal(err)
	}
	linkModules(t, filepath.Join(dir, "node_modules"), modules)

	s, err := Serve(context.Background(), ServeOptions{
		Root:    root,
		Project: "fixture",
		Sub: manifest.Project{Name: "ui", Path: "ui", Kind: manifest.ProjectKindSvelte,
			Components: &manifest.Components{Setup: "src/app.css"}},
		Timeout: serveTimeout,
	})
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
	// No proxy, as the dev server is on loopback, and no Timeout: each
	// request, and the WebSocket dial, has a deadline of its own.
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	t.Cleanup(client.CloseIdleConnections)

	hmr := dialHMR(t, client, base, vite)
	if p := hmr.next(t, updateTimeout, vite); p.Type != "connected" {
		t.Fatalf("Vite's HMR socket first sent %s, want connected%s", p.raw, vite())
	}

	// The page, and what it imports, as a browser loading it would, so that
	// Vite's module graph holds the component and the stylesheet.
	page := get(t, client, base+HarnessPath+"?c=src/lib/Hello.svelte", vite)
	for _, want := range []string{"/@vite/client", "/@id/virtual:flai-components/harness.mjs", `id="flai-component"`} {
		if !strings.Contains(page, want) {
			t.Errorf("the harness page lacks %s:\n%s", want, page)
		}
	}
	get(t, client, base+"/@id/virtual:flai-components/harness.mjs", vite)
	get(t, client, base+"/@id/virtual:flai-components/Frame.svelte", vite)
	catalogue := get(t, client, base+"/@id/__x00__virtual:flai-components", vite)
	for _, want := range []string{`"src/lib/Hello.svelte": () => import(`, `export const setupPath = "src/app.css";`} {
		if !strings.Contains(catalogue, want) {
			t.Errorf("the catalogue lacks %s:\n%s", want, catalogue)
		}
	}
	if js := get(t, client, base+"/src/lib/Hello.svelte", vite); !strings.Contains(js, "Hello, ") {
		t.Errorf("the component's module lacks its text:\n%s", js)
	}
	if css := get(t, client, base+"/src/app.css", vite); !strings.Contains(css, "rebeccapurple") {
		t.Errorf("the stylesheet's module lacks its value:\n%s", css)
	}
	hmr.settle(t, vite)

	edit(t, filepath.Join(dir, "src", "lib", "Hello.svelte"), "Hello, ", "Hi there, ")
	hmr.awaitUpdate(t, "/src/lib/Hello.svelte", vite)
	hmr.quiet(t, vite)
	if js := get(t, client, base+"/src/lib/Hello.svelte", vite); !strings.Contains(js, "Hi there, ") {
		t.Errorf("the component's module does not carry the edit:\n%s", js)
	}

	edit(t, filepath.Join(dir, "src", "app.css"), "rebeccapurple", "teal")
	hmr.awaitUpdate(t, "/src/app.css", vite)
	hmr.quiet(t, vite)
	if css := get(t, client, base+"/src/app.css", vite); !strings.Contains(css, "teal") {
		t.Errorf("the stylesheet's module does not carry the rewrite:\n%s", css)
	}
}

// flaioverModules is flaiover's installed node_modules, which has vite,
// svelte, and @sveltejs/vite-plugin-svelte: this checkout's, or, in a story
// worktree whose flaiover has none installed, the main checkout's. The test
// skips when neither has vite.
func flaioverModules(t *testing.T) string {
	t.Helper()
	repo, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	candidates := []string{filepath.Join(repo, "flaiover", "node_modules")}
	if git, err := exec.LookPath("git"); err == nil {
		common, err := exec.Command(git, "-C", repo, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
		if err == nil {
			checkout := filepath.Dir(strings.TrimSpace(string(common)))
			candidates = append(candidates, filepath.Join(checkout, "flaiover", "node_modules"))
		}
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, ".bin", viteBin)); err == nil {
			return c
		}
	}
	t.Skipf("flaiover's dependencies are not installed in %s, so the fixture has no vite to run; install them with scripts/flaiover-install.sh", strings.Join(slices.Compact(candidates), " or "))
	return ""
}

// linkModules makes dst a folder of links to each entry of the node_modules
// src, .bin and pnpm's .pnpm among them, except Vite's caches. dst is a
// folder rather than a link so that what Vite writes into node_modules,
// its optimized dependencies and its bundled config, lands in the fixture,
// not in src.
func linkModules(t *testing.T, dst, src string) {
	t.Helper()
	if err := os.Mkdir(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range moduleNames(t, src) {
		if strings.HasPrefix(name, ".vite") {
			continue
		}
		if err := os.Symlink(filepath.Join(src, name), filepath.Join(dst, name)); err != nil {
			t.Fatal(err)
		}
	}
}

// moduleNames are the names of the entries of the node_modules dir, sorted.
func moduleNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names
}

// get is the body of a GET of u, which must answer 200.
func get(t *testing.T, client *http.Client, u string, vite func() string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v%s", u, err, vite())
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET %s: %v%s", u, err, vite())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: %d\n%s%s", u, resp.StatusCode, body, vite())
	}
	return string(body)
}

// edit rewrites the file at path with its first old replaced by with, as an
// editor saving it in place does.
func edit(t *testing.T, path, old, with string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), old) {
		t.Fatalf("%s lacks %q:\n%s", path, old, b)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(b), old, with, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

// hmrPayload is a message Vite sends its client over the HMR socket, as far
// as the test reads it, and its JSON.
type hmrPayload struct {
	Type    string      `json:"type"`
	Updates []hmrUpdate `json:"updates"`
	Err     hmrError    `json:"err"`
	raw     string
}

// hmrError is an error payload's error, as far as the test reports it.
type hmrError struct {
	Message string `json:"message"`
}

// hmrUpdate is one module's update in an update payload: js-update for a
// module that accepts it, css-update for a linked stylesheet.
type hmrUpdate struct {
	Type         string `json:"type"`
	Path         string `json:"path"`
	AcceptedPath string `json:"acceptedPath"`
}

// hmrClient is a connection to Vite's HMR socket, as Vite's own client
// makes it, and the payloads read from it, in order; the channel closes when
// the connection ends, with why in err.
type hmrClient struct {
	payloads chan hmrPayload
	err      error
}

// hmrToken finds the token Vite writes into the client it serves, which
// the client passes as the socket's ?token= (Vite 6 and later).
var hmrToken = regexp.MustCompile(`const wsToken = "([^"]*)";`)

// dialHMR connects to the HMR socket of the dev server at base, as Vite's
// client does: at the server's root, with the token from /@vite/client and
// the vite-hmr subprotocol. The connection closes when the test ends.
func dialHMR(t *testing.T, client *http.Client, base string, vite func() string) *hmrClient {
	t.Helper()
	m := hmrToken.FindStringSubmatch(get(t, client, base+"/@vite/client", vite))
	if m == nil {
		t.Fatalf("/@vite/client has no wsToken, so the test cannot connect as Vite's client does%s", vite())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ws := "ws" + strings.TrimPrefix(base, "http") + "/?token=" + url.QueryEscape(m[1])
	conn, resp, err := websocket.Dial(ctx, ws, &websocket.DialOptions{HTTPClient: client, Subprotocols: []string{"vite-hmr"}})
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		t.Fatalf("cannot connect to Vite's HMR socket at %s: %v%s", ws, err, vite())
	}
	conn.SetReadLimit(1 << 22)
	h := &hmrClient{payloads: make(chan hmrPayload, 64)}
	go func() {
		defer close(h.payloads)
		for {
			_, data, err := conn.Read(context.Background())
			if err != nil {
				h.err = err
				return
			}
			p := hmrPayload{raw: string(data)}
			if err := json.Unmarshal(data, &p); err != nil {
				h.err = fmt.Errorf("the payload Vite sent is not JSON: %s: %w", data, err)
				return
			}
			h.payloads <- p
		}
	}()
	t.Cleanup(func() { _ = conn.CloseNow() })
	return h
}

// next is the next payload, waiting at most within for it.
func (h *hmrClient) next(t *testing.T, within time.Duration, vite func() string) hmrPayload {
	t.Helper()
	p, ok, timedOut := h.read(within)
	switch {
	case timedOut:
		t.Fatalf("Vite sent nothing over its HMR socket within %s%s", within, vite())
	case !ok:
		t.Fatalf("Vite's HMR socket closed: %v%s", h.err, vite())
	}
	return p
}

// read is the next payload, ok false when the connection ended, and
// timedOut true when none came within.
func (h *hmrClient) read(within time.Duration) (p hmrPayload, ok, timedOut bool) {
	timer := time.NewTimer(within)
	defer timer.Stop()
	select {
	case p, ok = <-h.payloads:
		return p, ok, false
	case <-timer.C:
		return hmrPayload{}, true, true
	}
}

// settle reads what Vite sends while it finishes loading the page's
// modules, such as a reload once it has optimized a dependency it found
// late, until it has sent nothing for quietWindow, so that what the test
// reads next is what its edits cause.
func (h *hmrClient) settle(t *testing.T, vite func() string) {
	t.Helper()
	deadline := time.Now().Add(updateTimeout)
	for {
		if time.Now().After(deadline) {
			t.Fatalf("Vite kept sending over its HMR socket for %s after the page loaded%s", updateTimeout, vite())
		}
		p, ok, timedOut := h.read(quietWindow)
		switch {
		case timedOut:
			return
		case !ok:
			t.Fatalf("Vite's HMR socket closed: %v%s", h.err, vite())
		}
		t.Logf("before the edits, Vite sent %s", p.raw)
	}
}

// awaitUpdate waits for an update payload with an update of the module at
// path, failing on a full reload or an error first.
func (h *hmrClient) awaitUpdate(t *testing.T, path string, vite func() string) {
	t.Helper()
	deadline := time.Now().Add(updateTimeout)
	for {
		p := h.next(t, time.Until(deadline), vite)
		switch p.Type {
		case "full-reload":
			t.Fatalf("Vite sent a full reload, not a hot update of %s: %s%s", path, p.raw, vite())
		case "error":
			t.Fatalf("Vite sent an error, not a hot update of %s: %s%s", path, p.Err.Message, vite())
		case "update":
			if slices.ContainsFunc(p.Updates, func(u hmrUpdate) bool {
				return (u.Type == "js-update" || u.Type == "css-update") && (u.Path == path || u.AcceptedPath == path)
			}) {
				return
			}
		}
		t.Logf("waiting for the update of %s, Vite sent %s", path, p.raw)
	}
}

// quiet fails when Vite sends a full reload or an error within quietWindow.
func (h *hmrClient) quiet(t *testing.T, vite func() string) {
	t.Helper()
	deadline := time.Now().Add(quietWindow)
	for {
		p, ok, timedOut := h.read(time.Until(deadline))
		switch {
		case timedOut:
			return
		case !ok:
			t.Fatalf("Vite's HMR socket closed: %v%s", h.err, vite())
		case p.Type == "full-reload" || p.Type == "error":
			t.Fatalf("after the hot update, Vite sent %s%s", p.raw, vite())
		}
	}
}

// waitGroupGone waits for every process in the process group pgid, which
// vite led, to end, and fails, killing them, when one has not within 5s.
func waitGroupGone(t *testing.T, pgid int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !errors.Is(syscall.Kill(-pgid, 0), syscall.ESRCH) {
		if time.Now().After(deadline) {
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
			t.Fatalf("a process of vite's group %d is still running", pgid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
