import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

let owner: string | undefined;
let asked: { method: string; params: Record<string, unknown> }[] = [];
let refuse: AgentError | undefined;
const entry = (author: string) => ({
	at: `2026-09-26T07:00:0${author.length % 10}Z`,
	author,
	text: 'x'
});
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	if (method === 'project.info')
		return { layout: { design: 'design', docs: 'docs', wip: 'wip' }, owner };
	if (method === 'threads.list')
		return [
			{ id: 'TH-0001', status: 'open', entries: [entry('agent-S-0001'), entry('dana')] },
			{ id: 'TH-0002', status: 'resolved', entries: [entry('designer')] }
		];
	if (method === 'thread.confirm') {
		if (refuse) throw refuse;
		return { data: { id: params.id, status: 'answered' }, warnings: [] };
	}
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

// ADR-0090: the operator confirms a pending recommendation, and it is the answer.
describe('/api/threads/:id/confirm', () => {
	let POST: Handler;
	beforeAll(async () => {
		POST = (await import('./[id]/confirm/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		useRepo(new Repo('/nowhere', channel));
		asked = [];
		refuse = undefined;
	});
	afterAll(() => useRepo(null));

	const confirm = (id: string) =>
		POST({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: '{}' })
		} as never);

	it("asks flai's thread.confirm with the thread's ID and passes on its answer", async () => {
		const r = await confirm('TH-0007');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ id: 'TH-0007', status: 'answered', warnings: [] });
		const write = asked.find((a) => a.method === 'thread.confirm');
		expect(write?.params).toMatchObject({ id: 'TH-0007' });
		expect(String(write?.params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
		expect(write?.params).not.toHaveProperty('by');
	});

	it('passes on why flai refused', async () => {
		refuse = new AgentError(400, 'TH-0007 has no recommendation awaiting confirmation', -32011);
		const r = await confirm('TH-0007');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe('TH-0007 has no recommendation awaiting confirmation');
	});
});
