import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { Repo, type Ask } from './repo';
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

// S-0201: flai marks a draft story's hit, and the dashboard passes the mark on with its route.
describe('search, as flai answers it', () => {
	it('passes a draft story’s flag through, and adds none where flai gave none', async () => {
		const hit = {
			kind: 'item',
			scope: 'wip',
			status: 'backlog',
			type: 'story',
			score: 1,
			snippet: ''
		};
		const ask = (async () => ({
			query: 'cranes',
			indexed: 2,
			hits: [
				{
					...hit,
					path: 'wip/kanban/stories/S-0002.md',
					itemId: 'S-0002',
					title: 'Cranes',
					draft: true
				},
				{ ...hit, path: 'wip/kanban/stories/S-0001.md', itemId: 'S-0001', title: 'Berths' }
			]
		})) as unknown as Ask;
		const { hits } = await search(new Repo('/nowhere', ask), 'cranes');
		expect(hits.map((h) => [h.itemId, h.route, h.draft])).toEqual([
			['S-0002', '/items/S-0002', true],
			['S-0001', '/items/S-0001', undefined]
		]);
		expect(hits[1]).not.toHaveProperty('draft');
	});
});

// Each call is a one-shot flai that builds the search index of the whole repository: about 2 s of
// CPU, and 4.8 to 6.1 s with 24 at once on a 24-core host, over the 5 s default (S-0501, I-0185).
describe('adrs on the monorepo', { timeout: 30000 }, () => {
	it('lists ADRs with the supersession chain', async () => {
		const list = await adrs(new Repo(monorepo, flaiAsk(monorepo)));
		expect(list.length).toBeGreaterThanOrEqual(15);
		expect(list[0].id).toBe('ADR-0001');
		expect(list.find((a) => a.id === 'ADR-0002')?.supersededBy).toEqual([]);
		expect(list.find((a) => a.id === 'ADR-0018')?.refines).toEqual(['ADR-0016']);
		expect(list.find((a) => a.id === 'ADR-0002')?.refines).toEqual([]);
		// a real repository holds proposed, superseded, and rejected ADRs as well as accepted ones
		// (ADR-0202 was rejected in favour of ADR-0203); the statuses allowed are those the ADR
		// template's status comment names (ADR-0223, I-0208), and a read that finds fewer fails here
		expect(list.find((a) => a.id === 'ADR-0001')?.status).toBe('accepted');
		const template = readFileSync(resolve(monorepo, 'design/adrs/0000-template.md'), 'utf8');
		const statusLine = template.split('\n').find((l) => l.startsWith('status:')) ?? '';
		const statuses = (statusLine.split('#')[1] ?? '').split('|').map((s) => s.trim());
		expect(statuses).toEqual(
			expect.arrayContaining(['proposed', 'accepted', 'superseded', 'deprecated', 'rejected'])
		);
		expect(list.filter((a) => !statuses.includes(a.status))).toEqual([]);
		expect(list.every((a) => /^\d{4}-\d{2}-\d{2}$/.test(a.date))).toBe(true);
	});
	it('finds a body phrase across wip and design', async () => {
		const { hits } = await search(new Repo(monorepo, flaiAsk(monorepo)), 'front matter');
		expect(hits.length).toBeGreaterThan(0);
		expect(hits.some((h) => h.path.startsWith('design/'))).toBe(true);
	});
});
