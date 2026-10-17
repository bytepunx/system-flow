//go:build !windows

package components

// The integration test of the component harness (S-0438, ADR-0163): Serve
// starts a real Vite dev server, the sub-project's own, on the fixture
// sub-project testdata/harness, whose node_modules links to flaiover's, and
// the test reads the hot updates Vite sends over its HMR WebSocket as the
// fixture's component and setup stylesheet change. Then (S-0440) it loads
// the harness in a headless Chromium and checks that a props message updates
// a component in place, keeping its own state, and (S-0441) that a tokens
// message lays its values over the stylesheet's. It needs node and
// flaiover's dependencies installed, so it skips without them, and under
// -short; the browser checks skip without Playwright's Chromium. The build
// tag keeps it off Windows: the fixture links node_modules with symlinks,
// and the test checks vite's process group.

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
	harness := get(t, client, base+"/@id/virtual:flai-components/harness.mjs", vite)
	for _, want := range []string{"case 'tokens':", "'data-flai-tokens'", "new MutationObserver(keepLast)"} {
		if !strings.Contains(harness, want) {
			t.Errorf("the harness module lacks %s:\n%s", want, harness)
		}
	}
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

	t.Run("a props message updates the component in place", func(t *testing.T) {
		checkInPlace(t, dir, s.URL(), vite)
	})
	t.Run("a tokens message lays its values over the stylesheet's", func(t *testing.T) {
		checkTokens(t, dir, s.URL(), vite)
	})
}

// inPlaceComponent is the fixture component whose own state, open, a props
// message must leave as it is.
const inPlaceComponent = "src/lib/Stateful.svelte"

// pageScriptStart and pageScriptEnd wrap the body of a node script that loads
// the harness page in a headless Chromium, with Playwright from flaiover's
// dependencies, and reports, as one line of JSON, what the body found, what
// the page logged as an error or a warning, and what the harness posted. The
// page loads the harness at the top, with ?origin= its own origin, so that it
// is its own parent: it posts the page's messages to itself, as the
// dashboard posts them into the frame, and receives what the harness posts.
// The body runs once the harness is ready, with the script's arguments after
// the harness URL as args, and send and seen to post and await messages. The
// script reports skip, and why, when Playwright or its browser is not
// installed.
const pageScriptStart = `const [url, ...args] = process.argv.slice(2);
const report = { errors: [], messages: [] };
let chromium;
try {
	({ chromium } = await import('@playwright/test'));
} catch (e) {
	console.log(JSON.stringify({ skip: 'Playwright is not installed: ' + String(e?.message ?? e).split('\n')[0] }));
	process.exit(0);
}
let browser;
try {
	browser = await chromium.launch();
} catch (e) {
	const why = String(e?.message ?? e);
	if (/Executable doesn't exist|missing dependencies/.test(why)) {
		console.log(JSON.stringify({ skip: 'Playwright has no Chromium to run: ' + why.split('\n')[0] }));
		process.exit(0);
	}
	throw e;
}
const page = await browser.newPage();
try {
	page.on('pageerror', (e) => report.errors.push(String(e)));
	page.on('console', (m) => {
		if (m.type() === 'error' || m.type() === 'warning') report.errors.push(m.type() + ': ' + m.text());
	});
	await page.addInitScript(() => {
		window.__flai = [];
		addEventListener('message', (e) => {
			const m = e.data;
			// The page's own messages come back to it too; only the harness's are kept.
			if (m?.source === 'flai-components' && !['select', 'props', 'theme', 'tokens'].includes(m.type)) window.__flai.push(m);
		});
	});
	const seen = (types, n) =>
		page.waitForFunction(([types, n]) => window.__flai.filter((m) => types.includes(m.type)).length >= n, [types, n], { timeout: 30000 });
	const send = (m) => page.evaluate((m) => postMessage({ source: 'flai-components', version: 1, ...m }, location.origin), m);
	await page.goto(url + '?origin=' + encodeURIComponent(new URL(url).origin));
	await seen(['ready'], 1);
`

const pageScriptEnd = `
} catch (e) {
	report.errors.push('the script: ' + String(e?.message ?? e));
} finally {
	report.messages = await page.evaluate(() => window.__flai).catch(() => []);
	await browser.close();
}
console.log(JSON.stringify(report));
`

// inPlaceScript is the body of a page script that reports what the page
// showed as the harness selected inPlaceComponent, the test clicked its
// button, and a props message changed its label and gave it a snippet.
const inPlaceScript = `	const [component] = args;
	const look = () =>
		page.evaluate(() => {
			const el = document.querySelector('.stateful');
			if (!el) return null;
			const slot = el.querySelector('.slot');
			// Svelte's anchors, comments, are left out of the snippet's markup.
			const markup = slot ? slot.innerHTML.replace(/<!--.*?-->/g, '') : null;
			return { label: el.querySelector('button').textContent, open: !!el.querySelector('.open'), slot: markup, marked: el.dataset.mark === 'first' };
		});
	await send({ type: 'select', component, props: { label: 'One' }, functions: ['onToggle'], snippets: ['children'] });
	await seen(['rendered', 'error'], 1);
	// A remount would replace the element, and the mark with it.
	await page.evaluate(() => {
		document.querySelector('.stateful').dataset.mark = 'first';
	});
	report.mounted = await look();
	await page.click('.stateful button');
	await seen(['event'], 1);
	report.opened = await look();
	await send({ type: 'props', props: { label: 'Two', children: 'Some <b>text</b>' }, functions: ['onToggle'], snippets: ['children'] });
	await seen(['rendered', 'error'], 2);
	report.changed = await look();
	await page.click('.stateful button');
	await seen(['event'], 2);
	report.closed = await look();`

// inPlaceLook is what the page showed of inPlaceComponent: its button's
// label, whether it was open, its snippet's markup, and whether its element
// is the one it first mounted.
type inPlaceLook struct {
	Label  string  `json:"label"`
	Open   bool    `json:"open"`
	Slot   *string `json:"slot"`
	Marked bool    `json:"marked"`
}

// harnessMessage is a message the harness posted, as far as the test reads
// it.
type harnessMessage struct {
	Type      string          `json:"type"`
	Component string          `json:"component"`
	Name      string          `json:"name"`
	Args      json.RawMessage `json:"args"`
	Message   string          `json:"message"`
}

// runPage writes the page script around body to name in the fixture
// sub-project's folder dir, so that it imports Playwright from its
// node_modules, runs it with args, the harness URL first, and decodes its
// report into report. It fails the test for each error the page logged, and
// returns what the harness posted. It skips the test when Playwright or its
// Chromium is not installed, saying that byHand, a check in a browser, then
// stands in for it.
func runPage(t *testing.T, dir, name, body, byHand string, vite func() string, report any, args ...string) []harnessMessage {
	t.Helper()
	script := filepath.Join(dir, name)
	if err := os.WriteFile(script, []byte(pageScriptStart+body+pageScriptEnd), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "node", append([]string{script}, args...)...)
	cmd.Dir = dir
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("the page script %s failed: %v\n%s%s%s", name, err, out, stderr.String(), vite())
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	last := []byte(lines[len(lines)-1])
	var r struct {
		Skip     string           `json:"skip"`
		Errors   []string         `json:"errors"`
		Messages []harnessMessage `json:"messages"`
	}
	if err := json.Unmarshal(last, &r); err != nil {
		t.Fatalf("the page script %s's report is not JSON: %v\n%s%s", name, err, out, stderr.String())
	}
	if r.Skip != "" {
		t.Skip(r.Skip + "; check in a browser " + byHand)
	}
	if err := json.Unmarshal(last, report); err != nil {
		t.Fatalf("the page script %s's report does not decode: %v\n%s", name, err, last)
	}
	for _, e := range r.Errors {
		t.Errorf("the page logged %s", e)
	}
	return r.Messages
}

// checkInPlace runs inPlaceScript against the harness at harnessURL, and
// checks that the component kept its own state, and its element, through a
// props message, rendered the snippet's text as text, and that its function
// prop's logger posted each call.
func checkInPlace(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Mounted *inPlaceLook `json:"mounted"`
		Opened  *inPlaceLook `json:"opened"`
		Changed *inPlaceLook `json:"changed"`
		Closed  *inPlaceLook `json:"closed"`
	}
	messages := runPage(t, dir, "inplace.mjs", inPlaceScript, "that a props message keeps the component's own state", vite, &r,
		harnessURL, inPlaceComponent)
	slot := "<span>Some &lt;b&gt;text&lt;/b&gt;</span>"
	for _, c := range []struct {
		when string
		got  *inPlaceLook
		want inPlaceLook
	}{
		{"once mounted", r.Mounted, inPlaceLook{Label: "One", Marked: true}},
		{"once clicked", r.Opened, inPlaceLook{Label: "One", Open: true, Marked: true}},
		{"after the props message", r.Changed, inPlaceLook{Label: "Two", Open: true, Slot: &slot, Marked: true}},
		{"once clicked again", r.Closed, inPlaceLook{Label: "Two", Slot: &slot, Marked: true}},
	} {
		switch {
		case c.got == nil:
			t.Errorf("%s, the page showed no %s", c.when, inPlaceComponent)
		case c.got.Label != c.want.Label || c.got.Open != c.want.Open || c.got.Marked != c.want.Marked ||
			(c.got.Slot == nil) != (c.want.Slot == nil) || (c.got.Slot != nil && *c.got.Slot != *c.want.Slot):
			t.Errorf("%s, the page showed %s, want %s", c.when, lookString(c.got), lookString(&c.want))
		}
	}
	var events []string
	rendered := 0
	for _, m := range messages {
		switch m.Type {
		case "error":
			t.Errorf("the harness posted an error about %s: %s", m.Component, m.Message)
		case "rendered":
			if m.Component == inPlaceComponent {
				rendered++
			}
		case "event":
			events = append(events, m.Component+" "+m.Name+" "+string(m.Args))
		}
	}
	if rendered != 2 {
		t.Errorf("the harness posted rendered %d times, want 2: once mounted, once for the props message", rendered)
	}
	// The event is the click's by its type, and the function by its name.
	want := []string{
		inPlaceComponent + ` onToggle [true,{"type":"click"},"[function anonymous]"]`,
		inPlaceComponent + ` onToggle [false,{"type":"click"},"[function anonymous]"]`,
	}
	if !slices.Equal(events, want) {
		t.Errorf("the harness posted the events\n%s\nwant\n%s", strings.Join(events, "\n"), strings.Join(want, "\n"))
	}
}

// tokensScript is the body of a page script that reports what the page
// showed of the fixture stylesheet's --fixture-ink, and of the harness's
// tokens style element, before any tokens message, after one, after another
// stylesheet landed after its element, in the dark theme, after a malformed
// message, and after an empty one.
const tokensScript = `	const look = () =>
		page.evaluate(() => {
			const styles = document.querySelectorAll('style[data-flai-tokens]');
			return {
				ink: getComputedStyle(document.documentElement).getPropertyValue('--fixture-ink').trim(),
				styles: styles.length,
				last: styles.length === 1 && document.head.lastElementChild === styles[0],
				css: styles[0]?.textContent ?? ''
			};
		});
	const until = (f, arg) => page.waitForFunction(f, arg, { timeout: 30000 });
	const themed = (theme) => until((theme) => document.documentElement.dataset.theme === theme, theme);
	report.before = await look();
	await send({ type: 'tokens', light: { '--fixture-ink': 'red', '--fixture-bad': 'url(x)' }, dark: { 'fixture-ink': 'green' } });
	await until(() => getComputedStyle(document.documentElement).getPropertyValue('--fixture-ink').trim() === 'red');
	report.tried = await look();
	// A hot update, or a component's own styles, land after the element.
	await page.evaluate(async () => {
		const later = document.createElement('style');
		later.id = 'later';
		later.textContent = ':root { --fixture-ink: blue }';
		document.head.append(later);
		await new Promise((r) => setTimeout(r, 0));
	});
	report.appended = await look();
	await send({ type: 'theme', theme: 'dark' });
	await themed('dark');
	report.dark = await look();
	// The harness takes messages in order, so the malformed one is behind it
	// once the theme is light again.
	await send({ type: 'tokens', light: 'blue', dark: {} });
	await send({ type: 'theme', theme: 'light' });
	await themed('light');
	report.dropped = await look();
	await page.evaluate(() => document.getElementById('later').remove());
	await send({ type: 'tokens', light: {}, dark: {} });
	await until(() => !document.querySelector('style[data-flai-tokens]'));
	report.cleared = await look();`

// tokensLook is what the page showed: the computed --fixture-ink, how many
// tokens style elements the head held, whether the one is last in it, and
// its CSS.
type tokensLook struct {
	Ink    string `json:"ink"`
	Styles int    `json:"styles"`
	Last   bool   `json:"last"`
	CSS    string `json:"css"`
}

// checkTokens runs tokensScript against the harness at harnessURL, and
// checks that a tokens message's values win over the stylesheet's, in light
// and dark, from one style element that stays last in the head, that a
// malformed message is dropped, and that an empty one removes the element.
func checkTokens(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Before   *tokensLook `json:"before"`
		Tried    *tokensLook `json:"tried"`
		Appended *tokensLook `json:"appended"`
		Dark     *tokensLook `json:"dark"`
		Dropped  *tokensLook `json:"dropped"`
		Cleared  *tokensLook `json:"cleared"`
	}
	messages := runPage(t, dir, "tokens.mjs", tokensScript, "that a tokens message lays its values over the stylesheet's", vite, &r,
		harnessURL)
	for _, m := range messages {
		if m.Type == "error" {
			t.Errorf("the harness posted an error about %s: %s", m.Component, m.Message)
		}
	}
	if r.Before == nil || r.Before.Ink == "" || r.Before.Ink == "red" {
		t.Fatalf("before any tokens message, the page showed %+v, want the stylesheet's own --fixture-ink", r.Before)
	}
	css := ":root { --fixture-ink: red }\n[data-theme='dark'] { --fixture-ink: green }\n"
	for _, c := range []struct {
		when string
		got  *tokensLook
		want tokensLook
	}{
		{"after a tokens message", r.Tried, tokensLook{Ink: "red", Styles: 1, Last: true, CSS: css}},
		{"after a stylesheet landed after it", r.Appended, tokensLook{Ink: "red", Styles: 1, Last: true, CSS: css}},
		{"in the dark theme", r.Dark, tokensLook{Ink: "green", Styles: 1, Last: true, CSS: css}},
		{"after a malformed tokens message", r.Dropped, tokensLook{Ink: "red", Styles: 1, Last: true, CSS: css}},
		{"after an empty tokens message", r.Cleared, tokensLook{Ink: r.Before.Ink}},
	} {
		switch {
		case c.got == nil:
			t.Errorf("%s, the page reported nothing", c.when)
		case *c.got != c.want:
			t.Errorf("%s, the page showed %+v, want %+v", c.when, *c.got, c.want)
		}
	}
}

func lookString(l *inPlaceLook) string {
	slot := "no snippet"
	if l.Slot != nil {
		slot = "snippet " + *l.Slot
	}
	return fmt.Sprintf("{label %q, open %t, %s, the first element %t}", l.Label, l.Open, slot, l.Marked)
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
