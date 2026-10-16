// The Vite plugin of flai's component harness (S-0438, ADR-0163). flai writes
// it, harness.mjs, and harness.html into .flai-cache/components/<project>/
// <sub-project>/, beside a vite.config.mjs that loads the sub-project's own
// config with importVite and merges this plugin into it, and runs the
// sub-project's own vite in its folder with that config.
//
// The plugin serves the harness page at /__flai/components/, and three
// modules that no file backs, so that each resolves its imports, svelte
// above all, from the sub-project rather than from .flai-cache:
//
//   virtual:flai-components/harness.mjs   the page's script, harness.mjs
//   virtual:flai-components/Frame.svelte  holds the rendered component in an
//                                         error boundary, and takes new props
//   virtual:flai-components               the catalogue: root, setupPath,
//                                         setup(), and components, each
//                                         component's path in the sub-project
//                                         mapped to a function importing it
//
// The catalogue finds components as flai's Find does (components.go): each
// .svelte file whose path in the sub-project matches an include glob and no
// exclude glob, outside node_modules and folders whose names begin with a
// dot. A component added or removed while the server runs reloads the page.
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';
import { pathToFileURL } from 'node:url';

const page = '/__flai/components';
const harnessId = 'virtual:flai-components/harness.mjs';
const frameId = 'virtual:flai-components/Frame.svelte';
const catalogueId = 'virtual:flai-components';
const resolvedCatalogueId = '\0virtual:flai-components';

// frame renders nothing until update gives it props, so that the component
// first renders with them, and renders inside the boundary, so that a throw
// reaches onerror, with a reset that renders it again.
const frame = `<script>
	let { Component, onerror } = $props();
	let props = $state.raw();
	export function update(next) {
		props = next;
	}
</script>

{#if props}
	<svelte:boundary {onerror}>
		<Component {...props} />
	</svelte:boundary>
{/if}
`;

// flaiComponents is the plugin. root is the sub-project's folder, absolute;
// files the folder flai wrote the harness into; include and exclude the
// components key's globs, as flai's CheckGlob accepts them; setup the
// module, relative to root, imported before any component, or empty.
export default function flaiComponents({ root, files, include, exclude = [], setup = '' }) {
	const globs = { include: include.map(compileGlob), exclude: exclude.map(compileGlob) };
	const isComponent = (rel) =>
		rel.endsWith('.svelte') &&
		rel.split('/').slice(0, -1).every((s) => s !== 'node_modules' && s !== '..' && !s.startsWith('.')) &&
		matchAny(globs.include, rel) &&
		!matchAny(globs.exclude, rel);
	return {
		name: 'flai-components',
		enforce: 'pre',
		config(user) {
			// The page shows the harness's error itself; Vite's overlay would
			// cover the frame with it.
			if (user.server?.hmr !== false) return { server: { hmr: { overlay: false } } };
		},
		resolveId(id) {
			if (id === catalogueId) return resolvedCatalogueId;
			if (id === harnessId || id === frameId) return id;
		},
		load(id) {
			if (id === resolvedCatalogueId) return catalogue(root, setup, find(root, isComponent));
			if (id === harnessId) return fs.readFileSync(path.join(files, 'harness.mjs'), 'utf8');
			if (id === frameId) return frame;
		},
		configureServer(server) {
			// Added here rather than in a returned hook, so that it runs
			// before SvelteKit's router, which answers every other path.
			server.middlewares.use(async (req, res, next) => {
				const url = new URL(req.url ?? '/', 'http://flai.invalid');
				if (url.pathname !== page && url.pathname !== page + '/') return next();
				try {
					const html = fs.readFileSync(path.join(files, 'harness.html'), 'utf8');
					const body = await server.transformIndexHtml(url.pathname + url.search, html, req.originalUrl);
					res.statusCode = 200;
					res.setHeader('Content-Type', 'text/html; charset=utf-8');
					res.setHeader('Cache-Control', 'no-store');
					res.end(body);
				} catch (e) {
					next(e);
				}
			});
			const changed = (file) => {
				const rel = slash(path.relative(root, file));
				if (!isComponent(rel)) return;
				const graph = server.environments?.client?.moduleGraph ?? server.moduleGraph;
				const mod = graph.getModuleById(resolvedCatalogueId);
				if (!mod) return;
				graph.invalidateModule(mod);
				(server.environments?.client?.hot ?? server.ws).send({ type: 'full-reload' });
			};
			server.watcher.on('add', changed);
			server.watcher.on('unlink', changed);
		}
	};
}

// importVite is the sub-project's own Vite, resolved from root as the
// sub-project's own config resolves it; a bare import of vite from
// .flai-cache would find none.
export async function importVite(root) {
	const require = createRequire(path.join(root, 'package.json'));
	let pkgFile;
	try {
		pkgFile = require.resolve('vite/package.json');
	} catch {
		throw new Error(`cannot resolve vite from ${root}: run npm install (or its package manager's install) there`);
	}
	const pkg = JSON.parse(fs.readFileSync(pkgFile, 'utf8'));
	return import(pathToFileURL(path.join(path.dirname(pkgFile), entry(pkg))).href);
}

// entry is the file a package's import of its root names.
function entry(pkg) {
	let e = pkg.exports?.['.'] ?? pkg.exports;
	if (e && typeof e === 'object') e = e.import ?? e.default;
	if (e && typeof e === 'object') e = e.default;
	return typeof e === 'string' ? e : (pkg.module ?? pkg.main ?? 'index.js');
}

// catalogue is the source of the catalogue module.
function catalogue(root, setup, rels) {
	const q = JSON.stringify;
	const lines = [`export const root = ${q(slash(root))};`, `export const setupPath = ${q(setup)};`];
	lines.push(
		setup
			? `export const setup = () => import(${q(slash(path.resolve(root, setup)))});`
			: 'export const setup = async () => {};'
	);
	lines.push('export const components = {');
	for (const rel of rels) lines.push(`\t${q(rel)}: () => import(${q(slash(path.join(root, rel)))}),`);
	lines.push('};');
	return lines.join('\n') + '\n';
}

// find is the components under root, by path relative to it in slash form.
function find(root, isComponent) {
	const out = [];
	const walk = (dir, rel) => {
		for (const d of fs.readdirSync(dir, { withFileTypes: true })) {
			const r = rel ? rel + '/' + d.name : d.name;
			if (d.isDirectory()) {
				if (d.name !== 'node_modules' && !d.name.startsWith('.')) walk(path.join(dir, d.name), r);
			} else if (d.isFile() && isComponent(r)) {
				out.push(r);
			}
		}
	};
	walk(root, '');
	return out.sort();
}

function slash(p) {
	return p.split(path.sep).join('/');
}

// compileGlob is a glob's segments, ** kept as itself and every other one a
// RegExp matching what Go's path.Match matches. A glob with no glob
// character matches the path itself and everything below it, as flai's
// matchAny does.
function compileGlob(g) {
	const segs = g.split('/').map((s) => (s === '**' ? s : compileSegment(s)));
	if (!/[*?[\\]/.test(g)) segs.push('**');
	return segs;
}

function compileSegment(seg) {
	const chars = [...seg];
	const lit = (c) => '\\u{' + c.codePointAt(0).toString(16) + '}';
	let re = '';
	for (let i = 0; i < chars.length; i++) {
		const c = chars[i];
		if (c === '*') {
			re += '[^/]*';
		} else if (c === '?') {
			re += '[^/]';
		} else if (c === '[') {
			let k = i + 1;
			let cls = '';
			if (chars[k] === '^') {
				cls = '^';
				k++;
			}
			while (k < chars.length && chars[k] !== ']') {
				let lo = chars[k++];
				if (lo === '\\') lo = chars[k++];
				cls += lit(lo);
				if (chars[k] === '-') {
					let hi = chars[++k];
					k++;
					if (hi === '\\') hi = chars[k++];
					cls += '-' + lit(hi);
				}
			}
			re += '[' + cls + ']';
			i = k;
		} else {
			re += lit(c === '\\' ? chars[++i] : c);
		}
	}
	return new RegExp('^' + re + '$', 'u');
}

function matchAny(globs, rel) {
	const segs = rel.split('/');
	return globs.some((g) => matchSegments(g, 0, segs, 0));
}

// matchSegments reports whether g's segments from i match s's from j: **
// matches zero or more whole segments, any other segment one.
function matchSegments(g, i, s, j) {
	if (i === g.length) return j === s.length;
	if (g[i] === '**') return matchSegments(g, i + 1, s, j) || (j < s.length && matchSegments(g, i, s, j + 1));
	return j < s.length && g[i].test(s[j]) && matchSegments(g, i + 1, s, j + 1);
}
