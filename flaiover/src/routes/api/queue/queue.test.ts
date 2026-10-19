// The merge queue's read and its hold over the channel (S-0399), scripted as the orchestrator
// route's test is: what flai answers is flai's, tested there; here it is what the route asks for
// and how it answers.

import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const queue = {
	mode: 'local',
	eligible: [],
	ineligible: [],
	enqueued: [{ story: 'S-0001', by: 'alex', at: '2026-10-19T10:00:00Z' }],
	held: true,
	held_by: 'alex',
	held_at: '2026-10-19T10:05:00Z',
	batch: { id: '20261019-100600', stories: ['S-0001'], step: 'verifying go-test' },
	outcomes: []
};

describe('/api/queue', () => {
	let get: Handler;
	let post: Handler;
	const doPost = (body: unknown) =>
		post({
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = mod.GET as unknown as Handler;
		post = mod.POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { answer: { host_actions: { 'merge-queue': true } } },
			'board.get': { answer: { queue } }
		};
	});
	afterAll(() => useRepo(null));

	it('answers whether the merge-queue action is on and the queue as board.get carries it', async () => {
		const r = await get({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ enabled: true, queue });
	});

	it('reads as off with no queue when flai sends no action and no queue', async () => {
		script['project.info'] = { error: new AgentError(503, 'no flai connected') };
		script['board.get'] = { answer: {} };
		const body = await (await get({} as never)).json();
		expect(body).toMatchObject({ enabled: false, queue: null });
	});

	it("asks flai to hold and to resume, and answers flai's change", async () => {
		script['queue.hold'] = { answer: { data: { action: 'hold', changed: true }, warnings: [] } };
		let r = await doPost({ held: true });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ action: 'hold', changed: true });
		script['queue.hold'] = { answer: { data: { action: 'resume', changed: true }, warnings: [] } };
		r = await doPost({ held: false });
		expect(await r.json()).toMatchObject({ action: 'resume' });
		const writes = asked.filter((a) => a.method === 'queue.hold');
		expect(writes.map((w) => w.params.held)).toEqual([true, false]);
		for (const w of writes) expect(typeof w.params.request_id).toBe('string');
	});

	it('says why when the merge-queue action is off', async () => {
		script['queue.hold'] = {
			error: new AgentError(
				403,
				'the host action "merge-queue" is not enabled for this project. On the host, in the project, run: flai serve enable merge-queue',
				-32012
			)
		};
		const r = await doPost({ held: true });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable merge-queue');
	});

	it('refuses a held that is not true or false, asking flai nothing', async () => {
		for (const body of [{ held: 'yes' }, {}, { held: 1 }, null]) {
			const r = await doPost(body);
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('held is true to hold the merge queue');
		}
		expect(asked.filter((a) => a.method === 'queue.hold')).toEqual([]);
	});
});
