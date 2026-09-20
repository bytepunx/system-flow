import { describe, expect, it } from 'vitest';
import { resolve } from 'node:path';
import { Repo } from './repo';
import { flaiAsk } from './testing';
import { adrs, search } from './search';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const monorepo = resolve('..');

// The index and the ranking are flai's (internal/search); here the dashboard's view of the answers.
describe('search on the fixture, through flai', () => {
	const r = new Repo(fixture, flaiAsk(fixture));
	const find = async (q: string, docs = false) => (await search(r, q, docs)).hits;
	it('indexes items and design files, docs only on request', async () => {
		expect((await search(r, '')).indexed).toBeGreaterThanOrEqual(12);
		const byId = await find('S-001');
		expect(byId[0].itemId).toBe('S-001');
		expect(byId[0].route).toBe('/items/S-001');
		expect(byId[0].kind).toBe('item');
		// the hit carries what the results page colours it by (S-0055)
		expect(byId[0].type).toBe('story');
		expect(byId[0].nature).toBe('feature');
		const byTitle = await find('Session start');
		expect(byTitle[0].path).toBe('design/conventions/session-start.md');
		expect(byTitle[0].route).toBe('/docs/design/conventions/session-start.md');
		expect(byTitle[0].nature).toBeUndefined();
		const noDocs = await find('docs');
		expect(noDocs.every((h) => h.scope !== 'docs')).toBe(true);
	});
	it('snippets show the match', async () => {
		const hits = await find('session');
		expect(hits[0].snippet.toLowerCase()).toContain('session');
	});
});

describe('adrs on the monorepo', () => {
	it('lists ADRs with the supersession chain', async () => {
		const list = await adrs(new Repo(monorepo, flaiAsk(monorepo)));
		expect(list.length).toBeGreaterThanOrEqual(15);
		expect(list[0].id).toBe('ADR-0001');
		expect(list.find((a) => a.id === 'ADR-0002')?.supersededBy).toEqual([]);
		expect(list.find((a) => a.id === 'ADR-0018')?.refines).toEqual(['ADR-0016']);
		expect(list.find((a) => a.id === 'ADR-0002')?.refines).toEqual([]);
		expect(list.every((a) => a.status === 'accepted')).toBe(true);
		expect(list.every((a) => /^\d{4}-\d{2}-\d{2}$/.test(a.date))).toBe(true);
	});
	it('finds a body phrase across wip and design', async () => {
		const { hits } = await search(new Repo(monorepo, flaiAsk(monorepo)), 'front matter');
		expect(hits.length).toBeGreaterThan(0);
		expect(hits.some((h) => h.path.startsWith('design/'))).toBe(true);
	});
});
