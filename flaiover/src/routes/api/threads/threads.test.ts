import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

type Handler = (event: never) => Promise<Response>;

let owner: string | undefined;
const entry = (author: string) => ({
	at: `2026-09-26T07:00:0${author.length % 10}Z`,
	author,
	text: 'x'
});
const channel = (async (method: string) => {
	if (method === 'project.info')
		return { layout: { design: 'design', docs: 'docs', wip: 'wip' }, owner };
	if (method === 'threads.list')
		return [
			{ id: 'TH-0001', status: 'open', entries: [entry('agent-S-0001'), entry('dana')] },
			{ id: 'TH-0002', status: 'resolved', entries: [entry('designer')] }
		];
	throw new Error(`unscripted method ${method}`);
}) as Ask;

// S-0127: each entry says whether the operator wrote it, so the page can set the two apart.
describe('/api/threads marks the operator', () => {
	let GET: Handler;
	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => useRepo(new Repo('/nowhere', channel)));
	afterAll(() => useRepo(null));

	const get = async (query: string) =>
		(await (await GET({ url: new URL(`http://x/api/threads${query}`) } as never)).json()) as {
			id: string;
			entries: { author: string; operator: boolean }[];
		}[];

	it("marks the manifest owner's entries and only those", async () => {
		owner = 'dana';
		const got = await get('?on=S-0001');
		expect(got.map((t) => t.id)).toEqual(['TH-0001']);
		expect(got[0].entries.map((e) => [e.author, e.operator])).toEqual([
			['agent-S-0001', false],
			['dana', true]
		]);
	});

	it('takes "designer" as the operator when the manifest names no owner, as flai does', async () => {
		owner = undefined;
		const got = await get('?all=1');
		expect(got.flatMap((t) => t.entries.map((e) => [e.author, e.operator]))).toEqual([
			['agent-S-0001', false],
			['dana', false],
			['designer', true]
		]);
	});
});
