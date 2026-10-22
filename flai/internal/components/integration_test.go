//go:build !windows

package components

// The integration test of the component harness (S-0438, ADR-0163): Serve
// starts a real Vite dev server, the sub-project's own, on the fixture
// sub-project testdata/harness, whose node_modules links to flaiover's, and
// the test reads the hot updates Vite sends over its HMR WebSocket as the
// fixture's component and setup stylesheet change. Then (S-0440) it loads
// the harness in a headless Chromium and checks that a props message updates
// a component in place, keeping its own state, (S-0441) that a tokens
// message lays its values over the stylesheet's, and (S-0449) how the
// harness frames a component: in a margin, a fixed dialog whole, with
// nothing of the component before, only once the page selects it, and with
// the props the page sent while it loaded, and that it answers the
// component's requests for the app's own server itself, and (S-0476) from
// the answers the page sends where one matches. It needs node and
// flaiover's dependencies installed from this tree's lockfile, so it skips
// without them, fails on an install from another lockfile, and skips under
// -short; the browser checks skip without Playwright's Chromium. The build
// tag keeps it off Windows: the fixture links node_modules with symlinks,
// and the test checks vite's process group.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
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
	t.Run("the component sits in a 10px margin, which its size includes", func(t *testing.T) {
		checkMargin(t, dir, s.URL(), vite)
	})
	t.Run("a fixed dialog is framed whole", func(t *testing.T) {
		checkDialog(t, dir, s.URL(), vite)
	})
	t.Run("a select leaves nothing of the component before", func(t *testing.T) {
		checkNothingLeft(t, dir, s.URL(), vite)
	})
	t.Run("framed, the harness mounts only what the page selects", func(t *testing.T) {
		checkOnlySelected(t, dir, s.URL()+"?c="+url.QueryEscape(helloComponent), vite)
	})
	t.Run("a props message sent while a select loads reaches the component", func(t *testing.T) {
		checkPropsWhileLoading(t, dir, s.URL(), vite)
	})
	t.Run("a request for the app's own server is answered by the harness", func(t *testing.T) {
		checkFetches(t, dir, s.URL(), vite)
	})
	t.Run("the harness answers the app's requests from the answers the page sends", func(t *testing.T) {
		checkAnswers(t, dir, s.URL(), vite)
	})
}

// inPlaceComponent is the fixture component whose own state, open, a props
// message must leave as it is.
const inPlaceComponent = "src/lib/Stateful.svelte"

// pageScriptStart and pageScriptEnd wrap the body of a node script that loads
// the harness page in a headless Chromium, with Playwright from flaiover's
// dependencies, and reports, as one line of JSON, what the body found, what
// the page logged as an error or a warning, and what the harness posted. The
// page loads the harness URL at the top, with ?origin= its own origin added, so that it
// is its own parent: it posts the page's messages to itself, as the
// dashboard posts them into the frame, and receives what the harness posts.
// The body runs once the harness is ready, with the script's arguments after
// the harness URL as args, send and seen to post and await messages, and
// sockets, the path and query of each WebSocket the page opened, and whether
// it received Vite's connected message. The script reports skip, and why,
// when Playwright or its browser is not installed.
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
	const sockets = [];
	page.on('websocket', (ws) => {
		const u = new URL(ws.url());
		const s = { path: u.pathname + u.search, connected: false };
		sockets.push(s);
		ws.on('framereceived', (f) => {
			if (String(f.payload).includes('"type":"connected"')) s.connected = true;
		});
	});
	const send = (m) => page.evaluate((m) => postMessage({ source: 'flai-components', version: 1, ...m }, location.origin), m);
	// lastSize is the last size the harness posted, once one follows its
	// last rendered, or within 2s.
	const lastSize = async () => {
		await page
			.waitForFunction(() => window.__flai.slice(window.__flai.findLastIndex((m) => m.type === 'rendered')).some((m) => m.type === 'size'), null, { timeout: 2000 })
			.catch(() => {});
		return page.evaluate(() => window.__flai.filter((m) => m.type === 'size').at(-1) ?? null);
	};
	const framed = new URL(url);
	framed.searchParams.set('origin', framed.origin);
	await page.goto(framed.href);
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

// The fixture components the framing checks (S-0449) select: Hello renders
// a line, Dialog a fixed backdrop centring a 300px dialog, and Leaky leaves
// an element ahead of its own and one in the body when it is unmounted.
const (
	helloComponent  = "src/lib/Hello.svelte"
	dialogComponent = "src/lib/Dialog.svelte"
	leakyComponent  = "src/lib/Leaky.svelte"
)

// margin is the space the harness keeps around the component, in CSS
// pixels, and dialogHeight the height of the fixture Dialog's dialog.
const (
	margin       = 10
	dialogHeight = 300
)

// rect is an element's box in the viewport, as getBoundingClientRect gives it.
type rect struct {
	Top    float64 `json:"top"`
	Left   float64 `json:"left"`
	Bottom float64 `json:"bottom"`
	Right  float64 `json:"right"`
}

// sizeMessage is a size message the harness posted.
type sizeMessage struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// marginScript is the body of a page script that reports where the page
// drew inPlaceComponent, which has no margin of its own, and the size the
// harness last posted.
const marginScript = `	const [component] = args;
	await send({ type: 'select', component, props: { label: 'One' } });
	await seen(['rendered', 'error'], 1);
	report.size = await lastSize();
	report.box = await page.evaluate(() => {
		const r = document.querySelector('.stateful').getBoundingClientRect();
		return { top: r.top, left: r.left, bottom: r.bottom, right: r.right };
	});`

// checkMargin runs marginScript against the harness at harnessURL, and
// checks that the component sits margin from the frame's top and left, and
// that the size posted is its box with margin on each side.
func checkMargin(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Box  *rect        `json:"box"`
		Size *sizeMessage `json:"size"`
	}
	messages := runPage(t, dir, "margin.mjs", marginScript, "that the component sits 10px from the frame's edges", vite, &r,
		harnessURL, inPlaceComponent)
	noErrors(t, messages)
	if r.Box == nil || r.Size == nil {
		t.Fatalf("the page reported the box %+v and the size %+v, want both", r.Box, r.Size)
	}
	if r.Box.Top != margin || r.Box.Left != margin {
		t.Errorf("the component's top left is %v, %v, want %d, %d", r.Box.Top, r.Box.Left, margin, margin)
	}
	if w, h := math.Ceil(r.Box.Right+margin), math.Ceil(r.Box.Bottom+margin); r.Size.Width != w || r.Size.Height != h {
		t.Errorf("the harness posted the size %vx%v for the box %+v, want %vx%v: the box and %dpx on each side", r.Size.Width, r.Size.Height, *r.Box, w, h, margin)
	}
}

// dialogScript is the body of a page script that sizes the viewport as the
// page sizes its frame, to each height the harness posts, 0 as unsized, until
// it posts no other: first from the page's unsized 160px, then from 800px,
// taller than the dialog. It reports, for each, the heights it set, and
// where the dialog then is in a viewport how tall.
const dialogScript = `	const [component] = args;
	const unsized = 160;
	const heights = () => page.evaluate(() => window.__flai.filter((m) => m.type === 'size').map((m) => m.height));
	const posted = (n) =>
		page.waitForFunction((n) => window.__flai.filter((m) => m.type === 'size').length > n, n, { timeout: 2000 }).catch(() => {});
	const settle = async () => {
		const steps = [];
		for (let i = 0; i < 20; i++) {
			const all = await heights();
			const h = all.at(-1) || unsized;
			const { width, height } = page.viewportSize();
			if (h === height) break;
			steps.push(h);
			await page.setViewportSize({ width, height: h });
			await posted(all.length);
		}
		const at = await page.evaluate(() => {
			const r = document.querySelector('.dialog').getBoundingClientRect();
			return { top: r.top, bottom: r.bottom, viewport: innerHeight };
		});
		return { steps, ...at };
	};
	await page.setViewportSize({ width: 800, height: unsized });
	await send({ type: 'select', component });
	await seen(['rendered', 'error'], 1);
	await lastSize();
	report.grown = await settle();
	const n = (await heights()).length;
	await page.setViewportSize({ width: 800, height: 800 });
	await posted(n);
	report.shrunk = await settle();`

// dialogSettled is where the fixture Dialog's dialog was once the viewport
// settled, how tall the viewport then was, and the heights set on the way.
type dialogSettled struct {
	Steps    []float64 `json:"steps"`
	Top      float64   `json:"top"`
	Bottom   float64   `json:"bottom"`
	Viewport float64   `json:"viewport"`
}

// checkDialog runs dialogScript against the harness at harnessURL, and
// checks that a frame sized as the harness says shows the whole dialog
// within the margin, growing from the unsized frame and shrinking back from
// a taller one, rather than growing without end (I-0160).
func checkDialog(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Grown  *dialogSettled `json:"grown"`
		Shrunk *dialogSettled `json:"shrunk"`
	}
	messages := runPage(t, dir, "dialog.mjs", dialogScript, "that a dialog is framed whole", vite, &r,
		harnessURL, dialogComponent)
	noErrors(t, messages)
	for _, c := range []struct {
		when string
		got  *dialogSettled
		most float64
	}{
		{"from the unsized frame", r.Grown, math.Inf(1)},
		// Each step halves the room left over, and the last rounds up a pixel.
		{"from a frame taller than the dialog", r.Shrunk, dialogHeight + 2*margin + 1},
	} {
		if c.got != nil {
			t.Logf("%s, the frame was sized %v and settled at %vpx", c.when, c.got.Steps, c.got.Viewport)
		}
		switch {
		case c.got == nil:
			t.Errorf("%s, the page reported nothing", c.when)
		case len(c.got.Steps) == 20:
			t.Errorf("%s, the frame did not settle: it was sized %v", c.when, c.got.Steps)
		case c.got.Viewport < dialogHeight+2*margin || c.got.Viewport > c.most:
			t.Errorf("%s, the frame settled at %vpx, want from %d to %v: the dialog and %dpx above and below; it was sized %v",
				c.when, c.got.Viewport, dialogHeight+2*margin, c.most, margin, c.got.Steps)
		case c.got.Top < margin || c.got.Bottom > c.got.Viewport-margin:
			t.Errorf("%s, the dialog is from %v to %v in a %vpx frame, want it %dpx within it", c.when, c.got.Top, c.got.Bottom, c.got.Viewport, margin)
		}
	}
}

// nothingLeftScript is the body of a page script that reports how many of
// Leaky's elements, its own and those it leaves, the page holds once it is
// selected, and once another component is.
const nothingLeftScript = `	const [leaky, other] = args;
	const count = () => page.evaluate(() => document.querySelectorAll('.leaky, .leaked').length);
	await send({ type: 'select', component: leaky });
	await seen(['rendered', 'error'], 1);
	report.mounted = await count();
	await send({ type: 'select', component: other });
	await seen(['rendered', 'error'], 2);
	report.after = await count();
	report.other = await page.evaluate(() => !!document.querySelector('#flai-component .hello'));`

// checkNothingLeft runs nothingLeftScript against the harness at harnessURL,
// and checks that selecting another component removes all Leaky rendered.
func checkNothingLeft(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Mounted int  `json:"mounted"`
		After   int  `json:"after"`
		Other   bool `json:"other"`
	}
	messages := runPage(t, dir, "nothingleft.mjs", nothingLeftScript, "that a select leaves nothing of the component before", vite, &r,
		harnessURL, leakyComponent, helloComponent)
	noErrors(t, messages)
	if r.Mounted != 3 {
		t.Fatalf("once %s mounted, the page held %d of its elements, want 3: its own, and the two it leaves", leakyComponent, r.Mounted)
	}
	if r.After != 0 {
		t.Errorf("once %s was selected, the page still held %d of %s's elements, want none", helloComponent, r.After, leakyComponent)
	}
	if !r.Other {
		t.Errorf("once %s was selected, the page did not show it", helloComponent)
	}
}

// onlySelectedScript is the body of a page script, for a harness URL that
// names a component as ?c=, that reports what the harness mounted and said
// it rendered 2s after its ready, and once the page selected another.
const onlySelectedScript = `	const [component] = args;
	const look = () =>
		page.evaluate(() => ({
			hello: !!document.querySelector('.hello'),
			stateful: !!document.querySelector('.stateful'),
			rendered: window.__flai.filter((m) => m.type === 'rendered').map((m) => m.component)
		}));
	await page.waitForTimeout(2000);
	report.before = await look();
	await send({ type: 'select', component });
	await seen(['rendered', 'error'], 1);
	report.after = await look();`

// shownLook is which fixture components the page showed, and those the
// harness said it rendered.
type shownLook struct {
	Hello    bool     `json:"hello"`
	Stateful bool     `json:"stateful"`
	Rendered []string `json:"rendered"`
}

// checkOnlySelected runs onlySelectedScript against harnessURL, which names
// helloComponent as ?c=, and checks that the framed harness mounts nothing
// until the page selects, and then only what it selects.
func checkOnlySelected(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Before *shownLook `json:"before"`
		After  *shownLook `json:"after"`
	}
	messages := runPage(t, dir, "onlyselected.mjs", onlySelectedScript, "that a framed harness renders only what the page selects", vite, &r,
		harnessURL, inPlaceComponent)
	noErrors(t, messages)
	for _, c := range []struct {
		when string
		got  *shownLook
		want shownLook
	}{
		{"before the page selected", r.Before, shownLook{Rendered: []string{}}},
		{"once the page selected " + inPlaceComponent, r.After, shownLook{Stateful: true, Rendered: []string{inPlaceComponent}}},
	} {
		switch {
		case c.got == nil:
			t.Errorf("%s, the page reported nothing", c.when)
		case c.got.Hello != c.want.Hello || c.got.Stateful != c.want.Stateful || !slices.Equal(c.got.Rendered, c.want.Rendered):
			t.Errorf("%s, the page showed %+v, want %+v", c.when, *c.got, c.want)
		}
	}
}

// propsWhileLoadingScript is the body of a page script that posts a select
// and, before the harness can have loaded the component's module, a props
// message, and reports the text the component then showed.
const propsWhileLoadingScript = `	const [component] = args;
	await page.evaluate((component) => {
		const m = { source: 'flai-components', version: 1 };
		postMessage({ ...m, type: 'select', component, props: { name: 'the select' } }, location.origin);
		postMessage({ ...m, type: 'props', props: { name: 'the props' } }, location.origin);
	}, component);
	await seen(['rendered', 'error'], 1);
	report.text = await page.evaluate(() => document.querySelector('.hello')?.textContent ?? null);`

// checkPropsWhileLoading runs propsWhileLoadingScript against the harness at
// harnessURL, and checks that the component mounted with the props message's
// props, which the page sends after the select on every ready.
func checkPropsWhileLoading(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Text *string `json:"text"`
	}
	messages := runPage(t, dir, "propswhileloading.mjs", propsWhileLoadingScript, "that props sent while a component loads reach it", vite, &r,
		harnessURL, helloComponent)
	noErrors(t, messages)
	if r.Text == nil || !strings.HasSuffix(*r.Text, "the props!") {
		t.Errorf("the component showed %v, want the props message's name, the props", stringOrNone(r.Text))
	}
}

// fetcherComponent is the fixture component that asks the app's own server,
// in its module and as it mounts, and fetches a module of the dev server's.
const fetcherComponent = "src/lib/Fetcher.svelte"

// fetchScript is the body of a page script that selects fetcherComponent
// and reports what it was answered, as it shows it, the path of every
// request for /api the page sent, the WebSockets it opened, and whether a
// socket on Vite's vite-ping, opened through the harness, opened.
const fetchScript = `	const [component] = args;
	const sent = [];
	page.on('request', (r) => {
		const u = new URL(r.url());
		if (u.pathname.startsWith('/api/')) sent.push(u.pathname);
	});
	await send({ type: 'select', component });
	await seen(['rendered', 'error'], 1);
	await page.waitForSelector('.fetched', { timeout: 30000 });
	report.answers = JSON.parse(await page.textContent('.fetched'));
	report.sent = sent;
	report.sockets = sockets;
	// Vite's client opens its HMR socket before the harness replaces
	// WebSocket, but reopens it, after a ping on vite-ping, through the
	// harness's: the ping must reach Vite.
	const hmr = sockets.find((s) => s.connected);
	report.ping = hmr
		? await page.evaluate(
				(path) =>
					new Promise((resolve) => {
						const ws = new WebSocket(path, 'vite-ping');
						ws.onopen = () => {
							ws.close();
							resolve('open');
						};
						ws.onerror = () => resolve('error');
					}),
				hmr.path.split('?')[0]
			)
		: '';`

// socketAnswer is how fetcherComponent's WebSocket ended, as it shows it:
// the events it heard, in order, its readyState as it closed, and its close
// event's code and wasClean, and whether its state was WebSocket.CLOSED.
type socketAnswer struct {
	Heard  string `json:"heard"`
	State  int    `json:"state"`
	Code   int    `json:"code"`
	Clean  bool   `json:"clean"`
	Closed bool   `json:"closed"`
}

// pageSocket is a WebSocket the page opened: its path and query, and whether
// it received Vite's connected message.
type pageSocket struct {
	Path      string `json:"path"`
	Connected bool   `json:"connected"`
}

// fetchAnswer is a response as fetcherComponent shows it.
type fetchAnswer struct {
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Type       string `json:"type"`
	Body       string `json:"body"`
}

// checkFetches runs fetchScript against the harness at harnessURL, and
// checks that each request for the app's own server, by fetch, EventSource,
// or WebSocket, in the component's module or as it mounts, was answered by
// the harness, never sent, and posted as an event, and that a fetch of a
// module, and Vite's HMR socket, still reached the dev server.
func checkFetches(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Answers *struct {
			Module  int          `json:"module"`
			String  fetchAnswer  `json:"string"`
			URL     fetchAnswer  `json:"url"`
			Request fetchAnswer  `json:"request"`
			Asset   fetchAnswer  `json:"asset"`
			Stream  int          `json:"stream"`
			Socket  socketAnswer `json:"socket"`
		} `json:"answers"`
		Sent    []string     `json:"sent"`
		Sockets []pageSocket `json:"sockets"`
		Ping    string       `json:"ping"`
	}
	messages := runPage(t, dir, "fetches.mjs", fetchScript, "that a component's fetch of /api is answered 501 by the harness", vite, &r,
		harnessURL, fetcherComponent)
	noErrors(t, messages)
	if r.Answers == nil {
		t.Fatalf("%s showed no answers", fetcherComponent)
	}
	a := r.Answers
	if a.Module != http.StatusNotImplemented {
		t.Errorf("the module's fetch of /api/module was answered %d, want 501", a.Module)
	}
	for _, c := range []struct {
		what   string
		got    fetchAnswer
		method string
		path   string
	}{
		{"a string", a.String, "GET", "/api/x"},
		{"a URL, with a method in lower case", a.URL, "PUT", "/api/y?z=1"},
		{"a Request", a.Request, "POST", "/api/r"},
	} {
		want := fetchAnswer{
			Status: http.StatusNotImplemented, StatusText: "Not Implemented", Type: "application/json",
			Body: notAnswered(c.method, c.path),
		}
		if c.got != want {
			t.Errorf("the fetch by %s was answered %+v, want %+v", c.what, c.got, want)
		}
	}
	if a.Asset.Status != http.StatusOK || !strings.Contains(a.Asset.Body, "--fixture-ink") {
		t.Errorf("the fetch of /src/app.css was answered %+v, want the dev server's 200 with the stylesheet", a.Asset)
	}
	if a.Stream != 2 {
		t.Errorf("the EventSource for /api/events was in state %d at its error, want 2, closed", a.Stream)
	}
	if wantSocket := (socketAnswer{Heard: "error close", State: 3, Code: 1006, Clean: false, Closed: true}); a.Socket != wantSocket {
		t.Errorf("the WebSocket for /api/events was answered %+v, want %+v", a.Socket, wantSocket)
	}
	if len(r.Sent) != 0 {
		t.Errorf("the page sent %v, want no request for /api", r.Sent)
	}
	hmr := false
	for _, s := range r.Sockets {
		if strings.HasPrefix(s.Path, "/api/") {
			t.Errorf("the page opened a WebSocket to %s, want none for /api", s.Path)
		}
		hmr = hmr || s.Connected
	}
	if !hmr {
		t.Errorf("the page's WebSockets were %+v, want Vite's HMR socket connected", r.Sockets)
	}
	if r.Ping != "open" {
		t.Errorf("a WebSocket on vite-ping to Vite's HMR path ended %q, want open: the harness must let Vite's own sockets through", r.Ping)
	}
	var events []string
	for _, m := range messages {
		if m.Type == "event" {
			events = append(events, m.Component+" "+m.Name+" "+string(m.Args))
		}
	}
	want := []string{
		fetcherComponent + ` fetch ["GET","/api/module",501]`,
		fetcherComponent + ` fetch ["GET","/api/x",501]`,
		fetcherComponent + ` fetch ["PUT","/api/y?z=1",501]`,
		fetcherComponent + ` fetch ["POST","/api/r",501]`,
		fetcherComponent + ` EventSource ["/api/events?project=x",501]`,
		fetcherComponent + ` WebSocket ["/api/events?project=x"]`,
	}
	if !slices.Equal(events, want) {
		t.Errorf("the harness posted the events\n%s\nwant\n%s", strings.Join(events, "\n"), strings.Join(want, "\n"))
	}
}

// answeredComponent is the fixture component whose requests for the app's
// own server the page's answers answer.
const answeredComponent = "src/lib/Answered.svelte"

// fixtureAnswers are the answers the page selects answeredComponent with
// (ADR-0187): first some the harness cannot use, then, of those it can, one
// for a path with no query and one for that path with a query after it, two
// for a path by method, a 204 with a body the harness drops, and two for an
// EventSource, one for a POST.
const fixtureAnswers = `[
	"not an answer",
	{"events": []},
	{"path": 5},
	{"path": "api/missing"},
	{"path": "/api/missing", "status": 700},
	{"path": "/api/missing", "headers": {"bad name": "x"}},
	{"path": "/api/missing", "events": "none"},
	{"path": "/api/module", "body": "from the module"},
	{"path": "/api/items", "body": "every item"},
	{"path": "/api/items?id=1", "body": {"id": 1}},
	{"path": "/api/accept", "method": "DELETE", "status": 500},
	{"path": "/api/accept", "method": "post", "status": 201, "headers": {"Content-Type": "text/csv", "X-Flai": "yes"}, "body": "a,b"},
	{"path": "/api/gone", "status": 204, "body": "dropped"},
	{"path": "/api/events", "method": "POST", "events": [{"data": "not for a GET"}]},
	{"path": "/api/events", "events": [{"data": "one", "id": "1"}, {"data": {"n": 2}}, {"event": "change", "data": {"kind": "edited"}, "id": "3"}, {"data": "after the close"}]}
]`

// answersScript is the body of a page script that selects answeredComponent
// with fixtureAnswers and reports what it was answered, as it shows it. Then
// it sends props messages, one without answers, one with an answer of its
// own, and one whose answers are not an array, and reports what the page's
// own fetch was answered after each, and the path of every request for /api
// the page sent.
const answersScript = `	const [component, answers] = args;
	const sent = [];
	page.on('request', (r) => {
		const u = new URL(r.url());
		if (u.pathname.startsWith('/api/')) sent.push(u.pathname);
	});
	const ask = (path) => page.evaluate((path) => fetch(path).then(async (r) => ({ status: r.status, body: await r.text() })), path);
	await send({ type: 'select', component, answers: JSON.parse(answers) });
	await seen(['rendered', 'error'], 1);
	await page.waitForSelector('.answered', { timeout: 30000 });
	report.answered = JSON.parse(await page.textContent('.answered'));
	await send({ type: 'props', props: {} });
	await seen(['rendered', 'error'], 2);
	report.without = await ask('/api/items?id=1');
	await send({ type: 'props', props: {}, answers: [{ path: '/api/late', body: 'late' }] });
	await seen(['rendered', 'error'], 3);
	report.late = await ask('/api/late');
	await send({ type: 'props', props: {}, answers: { path: '/api/late', body: 'late' } });
	await seen(['rendered', 'error'], 4);
	report.notArray = await ask('/api/late');
	report.sent = sent;`

// answeredFetch is a response as answeredComponent shows it: its status,
// its content-type and x-flai headers, empty when it has none, and its body.
type answeredFetch struct {
	Status int    `json:"status"`
	Type   string `json:"type"`
	Flai   string `json:"flai"`
	Body   string `json:"body"`
}

// answeredStream is what an EventSource of answeredComponent's heard, in
// order, and its readyState once it had closed it.
type answeredStream struct {
	Heard []string `json:"heard"`
	State int      `json:"state"`
}

// askedFetch is a response to the page's own fetch, as answersScript reports
// it.
type askedFetch struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// notAnswered is the body of the harness's 501 to method and path.
func notAnswered(method, path string) string {
	return fmt.Sprintf(`{"error":"flai's component harness does not answer the app's own requests: %s %s"}`, method, path)
}

// checkAnswers runs answersScript against the harness at harnessURL, and
// checks that the harness answered each fetch and EventSource of the app's
// from the first answer that matches it, by path with its query before path
// alone, and by method, passing over those it cannot use; that it answered
// one no answer matches 501, or with an error; that each message's answers
// replaced the last's; that no request reached the dev server; and that the
// events it posted carry the status each request got.
func checkAnswers(t *testing.T, dir, harnessURL string, vite func() string) {
	t.Helper()
	var r struct {
		Answered *struct {
			Module  answeredFetch  `json:"module"`
			Query   answeredFetch  `json:"query"`
			Other   answeredFetch  `json:"other"`
			Post    answeredFetch  `json:"post"`
			Empty   answeredFetch  `json:"empty"`
			Events  answeredFetch  `json:"events"`
			Missing answeredFetch  `json:"missing"`
			Stream  answeredStream `json:"stream"`
			Silent  answeredStream `json:"silent"`
		} `json:"answered"`
		Without  askedFetch `json:"without"`
		Late     askedFetch `json:"late"`
		NotArray askedFetch `json:"notArray"`
		Sent     []string   `json:"sent"`
	}
	messages := runPage(t, dir, "answers.mjs", answersScript, "that the Components page's answers answer a component's fetch of /api", vite, &r,
		harnessURL, answeredComponent, fixtureAnswers)
	noErrors(t, messages)
	if r.Answered == nil {
		t.Fatalf("%s showed no answers", answeredComponent)
	}
	a := r.Answered
	text := "text/plain; charset=utf-8"
	fetched := func(t *testing.T, what string, got, want answeredFetch) {
		t.Helper()
		if got != want {
			t.Errorf("the fetch %s was answered %+v, want %+v", what, got, want)
		}
	}
	t.Run("an answered fetch gets the answer's status, headers, and body", func(t *testing.T) {
		fetched(t, "of /api/module, as the module was imported", a.Module, answeredFetch{Status: 200, Type: text, Body: "from the module"})
		fetched(t, "POST /api/accept, of the answers by method", a.Post, answeredFetch{Status: 201, Type: "text/csv", Flai: "yes", Body: "a,b"})
		fetched(t, "DELETE /api/gone, of a 204", a.Empty, answeredFetch{Status: 204})
	})
	t.Run("a fetch takes the answer for its query before the one for its path", func(t *testing.T) {
		fetched(t, "of /api/items?id=1", a.Query, answeredFetch{Status: 200, Type: "application/json", Body: `{"id":1}`})
		fetched(t, "of /api/items?id=2", a.Other, answeredFetch{Status: 200, Type: text, Body: "every item"})
	})
	t.Run("an answered EventSource opens and sends the answer's events until it is closed", func(t *testing.T) {
		want := answeredStream{State: 2, Heard: []string{
			"open|1",
			"onmessage|one|1", "message|one|1",
			`onmessage|{"n":2}|`, `message|{"n":2}|`,
			`change|{"kind":"edited"}|3`,
		}}
		if !slices.Equal(a.Stream.Heard, want.Heard) || a.Stream.State != want.State {
			t.Errorf("the EventSource for /api/events?project=x heard\n%s\nand was in state %d, want\n%s\nand state %d",
				strings.Join(a.Stream.Heard, "\n"), a.Stream.State, strings.Join(want.Heard, "\n"), want.State)
		}
	})
	t.Run("a request no answer matches is still answered 501", func(t *testing.T) {
		fetched(t, "of /api/missing, whose answers the harness cannot use", a.Missing,
			answeredFetch{Status: 501, Type: "application/json", Body: notAnswered("GET", "/api/missing")})
		fetched(t, "of /api/events?project=x, whose answers are an EventSource's", a.Events,
			answeredFetch{Status: 501, Type: "application/json", Body: notAnswered("GET", "/api/events?project=x")})
		if want := []string{"error|2"}; !slices.Equal(a.Silent.Heard, want) || a.Silent.State != 2 {
			t.Errorf("the EventSource for /api/silent heard %v and was in state %d, want %v and state 2", a.Silent.Heard, a.Silent.State, want)
		}
		if len(r.Sent) != 0 {
			t.Errorf("the page sent %v, want no request for /api", r.Sent)
		}
	})
	t.Run("each message's answers replace the last's", func(t *testing.T) {
		for _, c := range []struct {
			what string
			got  askedFetch
			want askedFetch
		}{
			{"of /api/items?id=1, after a props message without answers", r.Without, askedFetch{Status: 501, Body: notAnswered("GET", "/api/items?id=1")}},
			{"of /api/late, after a props message with an answer for it", r.Late, askedFetch{Status: 200, Body: "late"}},
			{"of /api/late, after a props message whose answers are not an array", r.NotArray, askedFetch{Status: 501, Body: notAnswered("GET", "/api/late")}},
		} {
			if c.got != c.want {
				t.Errorf("the page's fetch %s was answered %+v, want %+v", c.what, c.got, c.want)
			}
		}
	})
	t.Run("the posted events carry the status each request got", func(t *testing.T) {
		var events []string
		for _, m := range messages {
			if m.Type == "event" {
				events = append(events, m.Component+" "+m.Name+" "+string(m.Args))
			}
		}
		want := []string{
			answeredComponent + ` fetch ["GET","/api/module",200]`,
			answeredComponent + ` fetch ["GET","/api/items?id=1",200]`,
			answeredComponent + ` fetch ["GET","/api/items?id=2",200]`,
			answeredComponent + ` fetch ["POST","/api/accept",201]`,
			answeredComponent + ` fetch ["DELETE","/api/gone",204]`,
			answeredComponent + ` fetch ["GET","/api/events?project=x",501]`,
			answeredComponent + ` fetch ["GET","/api/missing",501]`,
			answeredComponent + ` EventSource ["/api/events?project=x",200]`,
			answeredComponent + ` EventSource ["/api/silent",501]`,
			answeredComponent + ` fetch ["GET","/api/items?id=1",501]`,
			answeredComponent + ` fetch ["GET","/api/late",200]`,
			answeredComponent + ` fetch ["GET","/api/late",501]`,
		}
		if !slices.Equal(events, want) {
			t.Errorf("the harness posted the events\n%s\nwant\n%s", strings.Join(events, "\n"), strings.Join(want, "\n"))
		}
	})
}

// noErrors fails the test for each error the harness posted.
func noErrors(t *testing.T, messages []harnessMessage) {
	t.Helper()
	for _, m := range messages {
		if m.Type == "error" {
			t.Errorf("the harness posted an error about %s: %s", m.Component, m.Message)
		}
	}
}

func stringOrNone(s *string) string {
	if s == nil {
		return "nothing"
	}
	return strconv.Quote(*s)
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
// worktree whose flaiover has none installed, the main checkout's, but only
// one installed from this checkout's flaiover/pnpm-lock.yaml (I-0178). The
// test fails when the only installs with vite were made from another
// lockfile, and skips when neither has vite.
func flaioverModules(t *testing.T) string {
	t.Helper()
	repo, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	lockfile := filepath.Join(repo, "flaiover", "pnpm-lock.yaml")
	lock, err := os.ReadFile(lockfile)
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
	candidates = slices.Compact(candidates)
	modules, stale := installedFrom(lock, candidates)
	if modules != "" {
		return modules
	}
	if len(stale) > 0 {
		t.Fatalf("flaiover's dependencies in %s were installed from another pnpm-lock.yaml than %s, so the dev server would miss what this tree's components import; install them in this tree with %s --if-needed",
			strings.Join(stale, " and "), lockfile, filepath.Join(repo, "scripts", "flaiover-install.sh"))
	}
	t.Skipf("flaiover's dependencies are not installed in %s, so the fixture has no vite to run; install them with scripts/flaiover-install.sh", strings.Join(candidates, " or "))
	return ""
}

// installedFrom is the first of the node_modules candidates that has vite
// and was installed from lock, as pnpm's copy of the lockfile it installed
// from, .pnpm/lock.yaml, says; stale lists the candidates with vite that
// were installed from another lockfile, or whose copy is missing.
func installedFrom(lock []byte, candidates []string) (modules string, stale []string) {
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, ".bin", viteBin)); err != nil {
			continue
		}
		if installed, err := os.ReadFile(filepath.Join(c, ".pnpm", "lock.yaml")); err == nil && bytes.Equal(installed, lock) {
			return c, nil
		}
		stale = append(stale, c)
	}
	return "", stale
}

// I-0178: a story worktree with no install of its own borrowed the main
// checkout's, made before an acceptance added a dependency to the lockfile,
// and every components test then failed on an import vite could not find.
func TestInstalledFromRefusesAnInstallFromAnotherLockfile(t *testing.T) {
	root := t.TempDir()
	install := func(name, lock string) string {
		dir := filepath.Join(root, name, "node_modules")
		for _, p := range []string{filepath.Join(dir, ".bin", viteBin), filepath.Join(dir, ".pnpm", "lock.yaml")} {
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, ".bin", viteBin), nil, 0o755); err != nil {
			t.Fatal(err)
		}
		if lock != "" {
			if err := os.WriteFile(filepath.Join(dir, ".pnpm", "lock.yaml"), []byte(lock), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return dir
	}
	before := "packages:\n  vite: {}\n"
	after := before + "  '@lucide/svelte': {}\n"
	none := filepath.Join(root, "worktree", "node_modules")
	main := install("main", before)
	current := install("current", after)
	nocopy := install("nocopy", "")

	for _, c := range []struct {
		name       string
		candidates []string
		want       string
		stale      []string
	}{
		{"the worktree has none and main's is from the lockfile before", []string{none, main}, "", []string{main}},
		{"main's is from this lockfile", []string{none, current}, current, nil},
		{"the worktree's own comes first", []string{current, main}, current, nil},
		{"a stale own install is passed over for a current one", []string{main, current}, current, nil},
		{"an install without pnpm's copy of its lockfile", []string{nocopy}, "", []string{nocopy}},
		{"no install has vite", []string{none}, "", nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, stale := installedFrom([]byte(after), c.candidates)
			if got != c.want || !slices.Equal(stale, c.stale) {
				t.Errorf("installedFrom = %q, stale %v; want %q, stale %v", got, stale, c.want, c.stale)
			}
		})
	}
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
