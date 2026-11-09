import { describe, expect, it } from 'vitest';
import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

// S-0626, I-0293: the dashboard's design named a /conventions page and a /streams page that
// flaiover no longer had. Every route its page table and its API table name must be a route
// under src/routes, so a row that outlives its route fails here, as TestReferenceIsCurrent
// fails on a stale docs/users/flai-reference.md.
const DESIGN = fileURLToPath(
	new URL('../../../../design/system/flaiover-dashboard.md', import.meta.url)
);
const ROUTES = fileURLToPath(new URL('../../routes', import.meta.url));

// Paths the design names that are served outside src/routes, each for the reason given.
const OUTSIDE: Record<string, string> = {
	'/api/events':
		'a WebSocket upgrade, routed in server.js and the dev server plugin in vite.config.ts (ADR-0182)'
};

const PAGE = ['+page.svelte', '+page.ts'];
const SERVER = ['+server.ts'];

/** The rows of the section's table whose first column is Route, after its header row, as cells. */
function table(markdown: string, heading: string): string[][] {
	const lines = markdown.split('\n');
	const start = lines.findIndex((l) => l.trim() === heading);
	if (start < 0) throw new Error(`no heading ${heading} in ${DESIGN}`);
	let i = start + 1;
	while (i < lines.length && !/^\|\s*Route\s*\|/.test(lines[i])) {
		if (lines[i].startsWith('## ')) throw new Error(`no Route table under ${heading}`);
		i++;
	}
	const rows: string[][] = [];
	for (i += 2; i < lines.length && lines[i].startsWith('|'); i++) rows.push(cells(lines[i]));
	return rows;
}

/** A table row's cells, split on the pipes a backslash does not escape. */
function cells(row: string): string[] {
	return row
		.replace(/^\|/, '')
		.replace(/\|\s*$/, '')
		.split(/(?<!\\)\|/)
		.map((c) => c.trim());
}

/** The URL paths a cell names in code spans, with the method, query, and anchor left out. */
function paths(cell: string): string[] {
	const out: string[] = [];
	for (const [, span] of cell.matchAll(/`([^`]+)`/g)) {
		const path = span.trim().replace(/^[A-Z]+\s+/, '');
		if (!path.startsWith('/')) continue;
		out.push(path.split(/[\s?#]/)[0]);
	}
	return out;
}

const isGroup = (name: string) => /^\(.+\)$/.test(name);
const isOptional = (name: string) => /^\[\[.+\]\]$/.test(name);
const isRest = (name: string) => /^\[\.\.\..+\]$/.test(name);
const isParam = (name: string) => /^\[[^[\]]+\]$/.test(name);
// The design writes a parameter as :id, <name>, or [story], and a rest parameter as [...path].
const isDocParam = (seg: string) => /^(:.+|<.+>|\[.+\])$/.test(seg);

/** Whether the folder, or a route below it, serves the segments with one of the files. */
function served(dir: string, segs: string[], files: string[]): boolean {
	const entries = readdirSync(dir, { withFileTypes: true });
	if (segs.length === 0 && entries.some((e) => e.isFile() && files.includes(e.name))) return true;
	for (const e of entries) {
		if (!e.isDirectory()) continue;
		const sub = join(dir, e.name);
		if (isGroup(e.name)) {
			if (served(sub, segs, files)) return true;
		} else if (isOptional(e.name)) {
			if (served(sub, segs, files) || (segs.length > 0 && served(sub, segs.slice(1), files)))
				return true;
		} else if (isRest(e.name)) {
			for (let k = 0; k <= segs.length; k++) if (served(sub, segs.slice(k), files)) return true;
		} else if (segs.length === 0) {
			continue;
		} else if (isParam(e.name)) {
			if (!/^\[\.\.\./.test(segs[0]) && served(sub, segs.slice(1), files)) return true;
		} else if (e.name === segs[0] && !isDocParam(segs[0])) {
			if (served(sub, segs.slice(1), files)) return true;
		}
	}
	return false;
}

function routeExists(path: string, files: string[]): boolean {
	if (path in OUTSIDE) return true;
	return served(ROUTES, path.split('/').filter(Boolean), files);
}

describe('the dashboard design names only routes flaiover has', () => {
	const design = readFileSync(DESIGN, 'utf8');

	it('serves every page and the data each page reads in the Views table', () => {
		const rows = table(design, '## Views');
		expect(rows.length).toBeGreaterThan(30);
		const missing: string[] = [];
		for (const row of rows) {
			for (const path of paths(row[0]))
				if (!routeExists(path, PAGE)) missing.push(`page ${path} (row ${row[0]})`);
			for (const path of paths(row[row.length - 1]).filter((p) => p.startsWith('/api/')))
				if (!routeExists(path, SERVER)) missing.push(`data ${path} (row ${row[0]})`);
		}
		expect(missing).toEqual([]);
	});

	it('serves every route in the API table', () => {
		const rows = table(design, '## API (S-0011)');
		expect(rows.length).toBeGreaterThan(30);
		const missing: string[] = [];
		for (const row of rows)
			for (const path of paths(row[0]))
				if (!routeExists(path, SERVER)) missing.push(`${path} (row ${row[0]})`);
		expect(missing).toEqual([]);
	});

	it('fails on a row whose route is gone', () => {
		expect(routeExists('/conventions', PAGE)).toBe(false);
		expect(routeExists('/api/streams', SERVER)).toBe(false);
		expect(routeExists('/items/:id', PAGE)).toBe(true);
		expect(routeExists('/docs/<path>', PAGE)).toBe(true);
		expect(routeExists('/api/agent-stream/[story]', SERVER)).toBe(true);
	});
});
