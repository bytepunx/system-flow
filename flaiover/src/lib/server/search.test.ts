import { describe, expect, it } from 'vitest';
import { resolve } from 'node:path';
import { Repo } from './repo';
import { SearchIndex, adrs, snippet } from './search';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const monorepo = resolve('..');

describe('SearchIndex on the fixture', () => {
	const idx = new SearchIndex(new Repo(fixture));
	it('indexes items and design files, docs only on request', async () => {
		await idx.build();
		expect(idx.size()).toBeGreaterThanOrEqual(12);
		const byId = await idx.search('S-001');
		expect(byId[0].itemId).toBe('S-001');
		expect(byId[0].route).toBe('/items/S-001');
		expect(byId[0].kind).toBe('item');
		const byTitle = await idx.search('Session start');
		expect(byTitle[0].path).toBe('design/conventions/session-start.md');
		expect(byTitle[0].route).toBe('/docs/design/conventions/session-start.md');
		const noDocs = await idx.search('docs');
		expect(noDocs.every((h) => h.scope !== 'docs')).toBe(true);
	});
	it('snippets show the match', () => {
		const s = snippet('one two three commit at landing four five six', ['landing']);
		expect(s).toContain('landing');
		expect(snippet('short', ['nope'])).toBe('short');
	});
});

describe('adrs on the monorepo', () => {
	it('lists ADRs with the supersession chain', async () => {
		const list = await adrs(new Repo(monorepo));
		expect(list.length).toBeGreaterThanOrEqual(15);
		expect(list[0].id).toBe('ADR-0001');
		expect(list.find((a) => a.id === 'ADR-0002')?.supersededBy).toEqual([]);
		expect(list.every((a) => a.status === 'accepted')).toBe(true);
		expect(list.every((a) => /^\d{4}-\d{2}-\d{2}$/.test(a.date))).toBe(true);
	});
	it('finds a body phrase across wip and design', async () => {
		const idx = new SearchIndex(new Repo(monorepo));
		const hits = await idx.search('front matter');
		expect(hits.length).toBeGreaterThan(0);
		expect(hits.some((h) => h.path.startsWith('design/'))).toBe(true);
	});
});
