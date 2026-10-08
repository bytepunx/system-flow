import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
const conversation = (id: string, closed: boolean, updated: string) => ({
	id,
	title: id,
	from: 'S-0001',
	to: 'S-0002',
	about: [],
	status: closed ? 'closed' : 'open',
	closed,
	closed_reason: closed ? 'S-0002 was accepted' : '',
	awaiting: closed ? '' : 'S-0002',
	participants: null,
	created: '2026-10-01T09:00:00Z',
	updated,
	path: `wip/messages/${id}-x.md`,
	entries: []
});
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	if (method === 'project.info') return { layout: { design: 'design', docs: 'docs', wip: 'wip' } };
	if (method === 'messages.list')
		return [
			conversation('MS-0001', true, '2026-10-05T09:00:00Z'),
			conversation('MS-0002', false, '2026-10-02T09:00:00Z'),
			conversation('MS-0003', false, '2026-10-04T09:00:00Z'),
			conversation('MS-0004', true, '2026-10-06T09:00:00Z'),
			conversation('MS-0005', false, '2026-10-04T09:00:00Z')
		];
	throw new Error(`unscripted method ${method}`);
}) as Ask;

// S-0336: the conversations between stories, the project's or one story's, open ones first.
describe('/api/messages', () => {
	let GET: Handler;
	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		useRepo(new Repo('/nowhere', channel));
		asked = [];
	});
	afterAll(() => useRepo(null));

	const get = async (query: string) =>
		(await (await GET({ url: new URL(`http://x/api/messages${query}`) } as never)).json()) as {
			id: string;
			participants: string[];
		}[];
	const listed = () => asked.filter((a) => a.method === 'messages.list').map((a) => a.params);

	it("serves the project's open conversations, the newest activity first, ties by ID", async () => {
		const got = await get('');
		expect(got.map((c) => c.id)).toEqual(['MS-0003', 'MS-0005', 'MS-0002']);
		expect(got[0].participants).toEqual([]);
		expect(listed()).toEqual([{ all: true }]);
	});

	it('serves the closed too with all=1, after the open ones', async () => {
		expect((await get('?all=1')).map((c) => c.id)).toEqual([
			'MS-0003',
			'MS-0005',
			'MS-0002',
			'MS-0004',
			'MS-0001'
		]);
	});

	it("asks for one story's conversations with story", async () => {
		expect((await get('?story=S-1')).map((c) => c.id)).toEqual(['MS-0003', 'MS-0005', 'MS-0002']);
		expect(listed()).toEqual([{ story: 'S-1', all: true }]);
	});

	it("serves one story's closed conversations too with story and all=1", async () => {
		expect((await get('?story=S-0001&all=1')).map((c) => c.id)).toHaveLength(5);
		expect(listed()).toEqual([{ story: 'S-0001', all: true }]);
	});
});
