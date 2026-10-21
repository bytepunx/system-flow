// The browser side of flai's component harness (S-0438, ADR-0163). plugin.mjs
// serves it as virtual:flai-components/harness.mjs, so that svelte resolves
// from the sub-project. It runs the setup module, then renders one component
// at a time and talks to the page that frames it by postMessage, under the
// protocol of ADR-0163: every message is {source: "flai-components",
// version: 1, type, ...}. It accepts messages only from the origin ?origin=
// names, from its parent window, and posts only to it; with no origin it
// renders ?c= and posts nothing. Framed, with an origin, it renders only the
// component the page selects, which the page does on every ready, so a
// reloaded frame never shows the component its URL first named.
//
// It frames the component with a margin, the padding harness.html gives
// #flai-component, and posts the size of all the component draws, a fixed
// dialog among it (S-0449). A select removes everything the component before
// rendered, in #flai-component and in the body.
//
// A select message may carry functions and snippets as a props message does,
// so that a component that calls a function prop, or renders a snippet, as it
// mounts has it. A props message carries every prop the page sets, not a
// change: Frame spreads them over the mounted component, which takes them in
// place, keeping its own state, and a prop left out falls back to its
// default.
//
// A tokens message (S-0441) is whole as well: its light and dark custom
// properties replace the last message's, in one <style data-flai-tokens>
// kept last in the head, and a message with none removes it.
//
// Framed, a component's request to the app's own server, such as a fetch of
// /api/..., is not sent (S-0449): the dev server's server code has no
// dashboard token and no host channel, so it would answer 500. The harness
// answers a fetch 501 with a JSON error, an EventSource with an error event
// and no connection, and a WebSocket (ADR-0182) with an error event, then a
// close event with code 1006, and no connection, and posts each as an event,
// named fetch, EventSource, or WebSocket, so that the page's event log shows
// it. Vite's and flai's own requests, a module or a static asset, Vite's HMR
// and ping sockets, and another host's go through.
import { createRawSnippet, flushSync, mount, unmount } from 'svelte';
import Frame from 'virtual:flai-components/Frame.svelte';
import { components, root, setup, setupPath } from 'virtual:flai-components';

const source = 'flai-components';
const version = 1;
const params = new URL(location.href).searchParams;
const origin = pageOrigin(params.get('origin'));
const target = document.getElementById('flai-component');

// current is the mounted component: its path, its Frame, whether it failed
// since it last rendered, the boundary's reset after it did, and the loggers
// and snippets it was given, kept so that a props message that does not
// change one passes the same function again.
let current = null;
// wanted is the path the last select named, which an error is about until
// that component mounts.
let wanted = '';
let selection = 0;
// loading is whether the last select is importing its component's module,
// and pending the last props message that came meanwhile, which the
// component mounts with in place of the select's props.
let loading = false;
let pending = null;
let size = '';
// margin is the space around the component, harness.html's padding of target.
const margin = 10;
// own is what the body holds once the setup module has run; anything else in
// it a component added, such as a portal.
let own = new Set();
// tokenStyle is the style element the last tokens message wrote, or null.
let tokenStyle = null;

// pageOrigin is value when it is an origin, such as http://127.0.0.1:4242,
// or empty.
function pageOrigin(value) {
	if (!value) return '';
	try {
		return new URL(value).origin === value ? value : '';
	} catch {
		return '';
	}
}

function post(message) {
	if (origin) parent.postMessage({ source, version, ...message }, origin);
}

// forTheApp is whether a request for url, with the Accept header accept, is
// for the app's own server code: to this dev server, and none of Vite's or
// flai's own, under /@, /__, or /node_modules/, nor Vite's ping, nor a
// module or a static asset, whose last segment has a file extension.
function forTheApp(url, accept) {
	if (url.origin !== location.origin || accept.includes('text/x-vite-ping')) return false;
	const p = url.pathname;
	if (p.startsWith('/@') || p.startsWith('/__') || p.startsWith('/node_modules/')) return false;
	return !/\.[^/.]+$/.test(p.slice(p.lastIndexOf('/') + 1));
}

// requested is the URL, method, and Accept header of fetch's arguments,
// taking the method and headers from init over the Request's, or null when
// they do not name a URL, which fetch itself then refuses.
function requested(input, init) {
	try {
		const req = input instanceof Request ? input : null;
		const url = new URL(req ? req.url : String(input), document.baseURI);
		const method = String(init?.method ?? req?.method ?? 'GET').toUpperCase();
		const accept = new Headers(init?.headers ?? req?.headers).get('accept') ?? '';
		return { url, method, accept };
	} catch {
		return null;
	}
}

// socketForTheApp is whether a WebSocket for url, with the subprotocols
// protocols, is for the app's own server: url, its ws: or wss: read as
// http: or https:, passes forTheApp, and no subprotocol is Vite's own, such
// as vite-hmr and vite-ping, which Vite's client opens on its base path.
function socketForTheApp(url, protocols) {
	const http = new URL(url.href);
	if (http.protocol === 'ws:') http.protocol = 'http:';
	else if (http.protocol === 'wss:') http.protocol = 'https:';
	if (!forTheApp(http, '')) return false;
	return ![].concat(protocols ?? []).some((p) => String(p).startsWith('vite-'));
}

// pathOf is url's path and query, as a component on this origin names it.
function pathOf(url) {
	return url.pathname + url.search;
}

// answerTheAppsRequests replaces fetch, EventSource, and WebSocket, before
// the setup module or any component runs, so that a request for the app's
// own server is answered here, and posted as an event of the component last
// selected.
function answerTheAppsRequests() {
	answerTheAppsSockets();
	const send = window.fetch.bind(window);
	window.fetch = function fetch(input, init) {
		const r = requested(input, init);
		if (!r || !forTheApp(r.url, r.accept)) return send(input, init);
		const path = pathOf(r.url);
		post({ type: 'event', component: wanted, name: 'fetch', args: [r.method, path] });
		const body = JSON.stringify({ error: `flai's component harness does not answer the app's own requests: ${r.method} ${path}` });
		return Promise.resolve(
			new Response(body, { status: 501, statusText: 'Not Implemented', headers: { 'content-type': 'application/json' } })
		);
	};
	const Native = window.EventSource;
	if (!Native) return;
	// An EventSource for the app's server opens no connection: it fails as
	// one refused would, with an error event, and is closed.
	window.EventSource = class EventSource extends EventTarget {
		static CONNECTING = 0;
		static OPEN = 1;
		static CLOSED = 2;
		constructor(url, init) {
			let u;
			try {
				u = new URL(String(url), document.baseURI);
			} catch {
				return new Native(url, init);
			}
			if (!forTheApp(u, 'text/event-stream')) return new Native(url, init);
			super();
			this.url = u.href;
			this.withCredentials = !!init?.withCredentials;
			this.readyState = 0;
			this.onopen = this.onmessage = this.onerror = null;
			post({ type: 'event', component: wanted, name: 'EventSource', args: [pathOf(u)] });
			setTimeout(() => {
				if (this.readyState === 2) return;
				this.readyState = 2;
				const e = new Event('error');
				this.onerror?.(e);
				this.dispatchEvent(e);
			});
		}
		close() {
			this.readyState = 2;
		}
	};
}

// answerTheAppsSockets replaces WebSocket. A WebSocket for the app's server
// (ADR-0182) opens no connection: it fails as one refused would, with an
// error event and then a close event, code 1006 and not clean, and is
// closed. Its on- handlers and its listeners both hear them.
function answerTheAppsSockets() {
	const Native = window.WebSocket;
	if (!Native) return;
	// Its states are on the class and on each socket, as the native one's are.
	const states = { CONNECTING: 0, OPEN: 1, CLOSING: 2, CLOSED: 3 };
	class WebSocket extends EventTarget {
		constructor(url, protocols) {
			let u;
			try {
				u = new URL(String(url), document.baseURI);
			} catch {
				return new Native(url, protocols);
			}
			if (!socketForTheApp(u, protocols)) return new Native(url, protocols);
			super();
			if (u.protocol === 'http:') u.protocol = 'ws:';
			else if (u.protocol === 'https:') u.protocol = 'wss:';
			this.url = u.href;
			this.readyState = 0;
			this.protocol = '';
			this.extensions = '';
			this.bufferedAmount = 0;
			this.binaryType = 'blob';
			for (const type of ['open', 'message', 'error', 'close']) {
				this['on' + type] = null;
				this.addEventListener(type, (e) => this['on' + type]?.call(this, e));
			}
			post({ type: 'event', component: wanted, name: 'WebSocket', args: [pathOf(u)] });
			setTimeout(() => {
				this.readyState = 3;
				this.dispatchEvent(new Event('error'));
				this.dispatchEvent(new CloseEvent('close', { code: 1006, reason: '', wasClean: false }));
			});
		}
		send() {
			if (this.readyState === 0) throw new DOMException("Failed to execute 'send' on 'WebSocket': Still in CONNECTING state.", 'InvalidStateError');
		}
		close() {
			if (this.readyState === 0) this.readyState = 2;
		}
	}
	for (const [name, value] of Object.entries(states)) {
		Object.defineProperty(WebSocket, name, { value, enumerable: true });
		Object.defineProperty(WebSocket.prototype, name, { value, enumerable: true });
	}
	window.WebSocket = WebSocket;
}

async function select(path, props, functions, snippets) {
	const token = ++selection;
	wanted = path;
	loading = false;
	pending = null;
	close();
	const load = Object.hasOwn(components, path) ? components[path] : null;
	if (!load) {
		failed(path, new Error(`${path} is not in the catalogue: check its path, and the sub-project's components include and exclude`));
		return;
	}
	let mod;
	loading = true;
	try {
		mod = await load();
	} catch (e) {
		if (token === selection) {
			loading = false;
			failed(path, e);
		}
		return;
	}
	if (token !== selection) return;
	loading = false;
	const given = pending ?? { props, functions, snippets };
	pending = null;
	const c = { path, frame: null, failed: false, reset: null, loggers: new Map(), snippets: new Map() };
	current = c;
	const onerror = (error, reset) => {
		c.reset = reset;
		failed(path, error);
	};
	try {
		c.frame = mount(Frame, { target, props: { Component: mod.default, onerror } });
	} catch (e) {
		failed(path, e);
		return;
	}
	render(given.props, given.functions, given.snippets);
}

// render gives the mounted component props, renders it again if it had
// failed, and posts rendered unless it fails now.
function render(props, functions, snippets) {
	const c = current;
	if (!c?.frame) return;
	c.failed = false;
	try {
		if (props !== undefined) c.frame.update(live(c, props, functions, snippets));
		const reset = c.reset;
		c.reset = null;
		reset?.();
		flushSync();
	} catch (e) {
		failed(c.path, e);
	}
	if (!c.failed) post({ type: 'rendered', component: c.path });
}

// close unmounts the component and removes what it left: anything still in
// target, such as what a mount that threw had rendered, and anything it
// added to the body, such as a portal.
function close() {
	if (current?.frame) unmount(current.frame);
	current = null;
	target.replaceChildren();
	for (const n of [...document.body.childNodes]) {
		if (!own.has(n)) n.remove();
	}
}

function failed(path, error) {
	if (current?.path === path) current.failed = true;
	console.error(error);
	post({ type: 'error', component: path, message: String(error?.message ?? error), ...where(error) });
}

// where is the file, relative to the sub-project when it is in it, line,
// and column an error names: Vite's loc or id for a module that failed to
// compile, else the first frame of its stack in the sub-project's own code.
function where(error) {
	const loc = error?.loc;
	if (loc?.file || error?.id) {
		return clean({ file: inProject(String(loc?.file ?? error.id).split('?')[0]), line: loc?.line, column: loc?.column });
	}
	for (const m of String(error?.stack ?? '').matchAll(/(https?:\/\/[^\s()]+?):(\d+):(\d+)/g)) {
		const file = fileOf(m[1]);
		if (file && !file.startsWith('virtual:') && !file.includes('node_modules/')) {
			return { file, line: Number(m[2]), column: Number(m[3]) };
		}
	}
	return {};
}

function clean(o) {
	return Object.fromEntries(Object.entries(o).filter(([, v]) => v !== undefined && v !== ''));
}

function inProject(file) {
	return file.startsWith(root + '/') ? file.slice(root.length + 1) : file;
}

// fileOf is the file a module URL of this server names, or empty.
function fileOf(url) {
	let u;
	try {
		u = new URL(url);
	} catch {
		return '';
	}
	if (u.origin !== location.origin) return '';
	const p = decodeURIComponent(u.pathname);
	if (p.startsWith('/@id/')) return p.slice('/@id/'.length);
	if (p.startsWith('/@fs/')) return inProject(p.slice('/@fs'.length));
	return p.slice(1);
}

// live is props as c's component takes them: each name in functions a
// logger that posts event with its arguments, and each name in snippets
// whose value is a string a snippet of that text.
function live(c, props, functions, snippets) {
	const out = { ...props };
	for (const name of names(functions)) {
		if (!c.loggers.has(name)) {
			c.loggers.set(name, (...args) => post({ type: 'event', component: c.path, name, args: args.map((a) => jsonSafe(a)) }));
		}
		out[name] = c.loggers.get(name);
	}
	for (const name of names(snippets)) {
		if (typeof out[name] !== 'string') continue;
		if (c.snippets.get(name)?.text !== out[name]) c.snippets.set(name, { text: out[name], snippet: textSnippet(out[name]) });
		out[name] = c.snippets.get(name).snippet;
	}
	return out;
}

function names(list) {
	return Array.isArray(list) ? list.filter((n) => typeof n === 'string') : [];
}

// textSnippet is a snippet rendering text, as text rather than markup, in a
// span, since a raw snippet renders one element; it takes no notice of the
// arguments the component renders it with.
function textSnippet(text) {
	return createRawSnippet(() => ({
		render: () => '<span></span>',
		setup: (span) => {
			span.textContent = text;
		}
	}));
}

// jsonSafe is value as far as JSON carries it: an event by its type and
// detail, a node by its tag, a function by its name, a cycle cut off.
function jsonSafe(value, seen = new WeakSet(), depth = 0) {
	switch (typeof value) {
		case 'string':
		case 'boolean':
			return value;
		case 'number':
			return Number.isFinite(value) ? value : String(value);
		case 'undefined':
			return null;
		case 'function':
			return `[function ${value.name || 'anonymous'}]`;
		case 'bigint':
		case 'symbol':
			return String(value);
	}
	if (value === null) return null;
	if (value instanceof Event) {
		return value instanceof CustomEvent ? { type: value.type, detail: jsonSafe(value.detail, seen, depth + 1) } : { type: value.type };
	}
	if (value instanceof Node) return `<${value.nodeName.toLowerCase()}>`;
	if (value instanceof Date) return value.toJSON();
	if (seen.has(value) || depth > 8) return '[…]';
	seen.add(value);
	if (Array.isArray(value)) return value.map((v) => jsonSafe(v, seen, depth + 1));
	return Object.fromEntries(Object.entries(value).map(([k, v]) => [k, jsonSafe(v, seen, depth + 1)]));
}

// tokens lays the custom properties in light and dark over the stylesheet's,
// under the selectors flai design css writes (ADR-0159), in the one style
// element tokenStyle, last in the head so that it wins; with none to lay, it
// removes the element. A message is whole: it replaces the last one.
function tokens(light, dark) {
	const css = block(':root', light) + block("[data-theme='dark']", dark);
	if (!css) {
		tokenStyle?.remove();
		tokenStyle = null;
		return;
	}
	if (!tokenStyle) {
		tokenStyle = document.createElement('style');
		tokenStyle.setAttribute('data-flai-tokens', '');
	}
	tokenStyle.textContent = css;
	document.head.append(tokenStyle);
}

// keepLast moves tokenStyle back to the end of the head when a stylesheet,
// such as a hot update's or a lazily loaded component's, lands after it.
// Its own move adds nothing after it, so it does not set itself off again.
function keepLast() {
	if (tokenStyle?.parentNode !== document.head) return;
	for (let n = tokenStyle.nextElementSibling; n; n = n.nextElementSibling) {
		if (n.matches('style, link[rel~="stylesheet"]')) {
			document.head.append(tokenStyle);
			return;
		}
	}
}

// isVars is whether v may be a tokens message's light or dark: an object of
// custom properties, or left out.
function isVars(v) {
	return v === undefined || (typeof v === 'object' && v !== null && !Array.isArray(v));
}

// unsafeValue is what a token's value must not carry into the stylesheet,
// as the Design page's preview refuses it (ADR-0149 § 4): url(, a brace, a
// semicolon, <, @, a backslash, or a comment, which would swallow the
// declarations after it.
const unsafeValue = /url\(|[{};<@\\]|\/\*/i;

// block is selector's rule setting each custom property in vars, by name
// with or without its leading --; a name or value that would break out of
// the rule is left out.
function block(selector, vars) {
	const decls = Object.entries(vars ?? {})
		.filter(([n, v]) => /^(--)?[\w-]+$/.test(n) && ['string', 'number'].includes(typeof v) && !unsafeValue.test(String(v)))
		.map(([n, v]) => `${n.startsWith('--') ? n : '--' + n}: ${v}`);
	return decls.length ? `${selector} { ${decls.join('; ')} }\n` : '';
}

// measure posts the size of all the component draws, with the margin on each
// side, when it changes: the union of target's box, its padding the margin,
// and the box of every element in it or in what it added to the body, a
// fixed or absolute one among them. A fixed element pinned to the viewport's
// top and bottom, such as a dialog's backdrop or a drawer, grows with the
// frame, so its own box is left out, but what it holds counts. Content above
// or left of the margin, such as a centred dialog taller than the frame,
// adds as much again, so that the frame grows until it fits. A horizontal
// scroll bar adds its height, so that it hides nothing. While a select loads
// its module, the last size stands.
function measure() {
	if (loading) return;
	const doc = document.documentElement;
	const t = target.getBoundingClientRect();
	const box = { top: Infinity, left: Infinity, bottom: t.bottom + scrollY - margin, right: t.right + scrollX - margin };
	for (const el of target.children) drawn(el, false, box);
	for (const n of document.body.children) {
		if (!own.has(n)) drawn(n, false, box);
	}
	const bar = doc.scrollWidth > doc.clientWidth ? innerHeight - doc.clientHeight : 0;
	const width = Math.ceil(box.right + margin + Math.max(0, margin - box.left));
	const height = Math.ceil(box.bottom + margin + Math.max(0, margin - box.top) + bar);
	if (`${width}x${height}` === size) return;
	size = `${width}x${height}`;
	post({ type: 'size', width, height });
}

// drawn widens box, in the document's coordinates, by el's box and those of
// the elements in it, leaving out one with no size, such as one not
// displayed, and a fixed one pinned to the viewport's top and bottom. A
// fixed element, and what it holds, is placed in the viewport, so the
// document's scroll does not move it.
function drawn(el, inFixed, box) {
	const isFixed = getComputedStyle(el).position === 'fixed';
	const fixed = inFixed || isFixed;
	const r = el.getBoundingClientRect();
	const pinned = isFixed && r.top <= 0.5 && r.bottom >= document.documentElement.clientHeight - 0.5;
	if ((r.width || r.height) && !pinned) {
		const x = fixed ? 0 : scrollX;
		const y = fixed ? 0 : scrollY;
		box.top = Math.min(box.top, r.top + y);
		box.left = Math.min(box.left, r.left + x);
		box.bottom = Math.max(box.bottom, r.bottom + y);
		box.right = Math.max(box.right, r.right + x);
	}
	for (const child of el.children) drawn(child, fixed, box);
}

function receive(event) {
	if (!origin || event.origin !== origin || event.source !== window.parent) return;
	const m = event.data;
	if (!m || m.source !== source || m.version !== version) return;
	switch (m.type) {
		case 'select':
			if (typeof m.component === 'string') select(m.component, m.props ?? {}, m.functions, m.snippets);
			break;
		case 'props':
			if (loading) pending = { props: m.props ?? {}, functions: m.functions, snippets: m.snippets };
			else render(m.props ?? {}, m.functions, m.snippets);
			break;
		case 'theme':
			if (m.theme === 'light' || m.theme === 'dark') document.documentElement.dataset.theme = m.theme;
			break;
		case 'tokens':
			if (isVars(m.light) && isVars(m.dark)) tokens(m.light, m.dark);
			break;
	}
}

const asked = params.get('c') ?? '';
if (origin) answerTheAppsRequests();
try {
	await setup();
} catch (e) {
	failed(asked, { message: `the setup module ${setupPath} failed: ${e?.message ?? e}`, loc: e?.loc, id: e?.id, stack: e?.stack });
}
own = new Set(document.body.childNodes);
if (import.meta.hot) {
	// A module that fails to compile, now or on a saved change, says where;
	// a hot update renders the component again, and posts rendered once it
	// renders.
	import.meta.hot.on('vite:error', ({ err }) => failed(current?.path ?? wanted, err));
	import.meta.hot.on('vite:afterUpdate', () => render());
}
// The size changes as target's box does, as the component's elements or a
// portal's change, and as the frame does, which moves a fixed dialog.
new ResizeObserver(measure).observe(target);
new MutationObserver(measure).observe(document.body, { childList: true, subtree: true, attributes: true, characterData: true });
window.addEventListener('resize', measure);
new MutationObserver(keepLast).observe(document.head, { childList: true });
window.addEventListener('message', receive);
post(asked ? { type: 'ready', component: asked } : { type: 'ready' });
if (asked && !origin) select(asked, {});
