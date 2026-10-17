// The browser side of flai's component harness (S-0438, ADR-0163). plugin.mjs
// serves it as virtual:flai-components/harness.mjs, so that svelte resolves
// from the sub-project. It runs the setup module, then renders one component
// at a time and talks to the page that frames it by postMessage, under the
// protocol of ADR-0163: every message is {source: "flai-components",
// version: 1, type, ...}. It accepts messages only from the origin ?origin=
// names, from its parent window, and posts only to it; with no origin it
// renders ?c= and posts nothing.
//
// A select message may carry functions and snippets as a props message does,
// so that a component that calls a function prop, or renders a snippet, as it
// mounts has it. A props message carries every prop the page sets, not a
// change: Frame spreads them over the mounted component, which takes them in
// place, keeping its own state, and a prop left out falls back to its
// default.
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
let size = '';

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

async function select(path, props, functions, snippets) {
	const token = ++selection;
	wanted = path;
	close();
	const load = Object.hasOwn(components, path) ? components[path] : null;
	if (!load) {
		failed(path, new Error(`${path} is not in the catalogue: check its path, and the sub-project's components include and exclude`));
		return;
	}
	let mod;
	try {
		mod = await load();
	} catch (e) {
		if (token === selection) failed(path, e);
		return;
	}
	if (token !== selection) return;
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
	render(props, functions, snippets);
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

function close() {
	if (!current) return;
	if (current.frame) unmount(current.frame);
	current = null;
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

function tokens(light, dark) {
	const css = block(':root', light) + block("[data-theme='dark']", dark);
	let style = document.getElementById('flai-tokens');
	if (!css) {
		style?.remove();
		return;
	}
	if (!style) {
		style = document.createElement('style');
		style.id = 'flai-tokens';
	}
	style.textContent = css;
	// Last in the head, so that it lies over the stylesheet's values.
	document.head.append(style);
}

// block is selector's rule setting each custom property in vars, by name
// with or without its leading --; a name or value that would break out of
// the rule is left out.
function block(selector, vars) {
	const decls = Object.entries(vars ?? {})
		.filter(([n, v]) => /^(--)?[\w-]+$/.test(n) && ['string', 'number'].includes(typeof v) && !/[;{}<]/.test(String(v)))
		.map(([n, v]) => `${n.startsWith('--') ? n : '--' + n}: ${v}`);
	return decls.length ? `${selector} { ${decls.join('; ')} }\n` : '';
}

function measure() {
	const rect = target.getBoundingClientRect();
	const width = Math.ceil(Math.max(rect.width, target.scrollWidth));
	const height = Math.ceil(Math.max(rect.height, target.scrollHeight));
	if (`${width}x${height}` === size) return;
	size = `${width}x${height}`;
	post({ type: 'size', width, height });
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
			render(m.props ?? {}, m.functions, m.snippets);
			break;
		case 'theme':
			if (m.theme === 'light' || m.theme === 'dark') document.documentElement.dataset.theme = m.theme;
			break;
		case 'tokens':
			tokens(m.light, m.dark);
			break;
	}
}

const asked = params.get('c') ?? '';
try {
	await setup();
} catch (e) {
	failed(asked, { message: `the setup module ${setupPath} failed: ${e?.message ?? e}`, loc: e?.loc, id: e?.id, stack: e?.stack });
}
if (import.meta.hot) {
	// A module that fails to compile, now or on a saved change, says where;
	// a hot update renders the component again, and posts rendered once it
	// renders.
	import.meta.hot.on('vite:error', ({ err }) => failed(current?.path ?? wanted, err));
	import.meta.hot.on('vite:afterUpdate', () => render());
}
new ResizeObserver(measure).observe(target);
window.addEventListener('message', receive);
post(asked ? { type: 'ready', component: asked } : { type: 'ready' });
if (asked) select(asked, {});
