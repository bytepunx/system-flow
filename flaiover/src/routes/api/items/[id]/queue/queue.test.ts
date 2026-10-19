// Enqueueing a story in the merge queue and taking it out over the channel (S-0399), scripted as
// the orchestrator route's test is: what flai answers is flai's, tested there; here it is what the
// route asks for and how it answers.

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

describe('/api/items/[id]/queue', () => {
	let post: Handler;
	let del: Handler;
	const event = { params: { id: 'S-0001' } } as never;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		post = mod.POST as unknown as Handler;
		del = mod.DELETE as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it("enqueues and dequeues the story, answering flai's change", async () => {
		script['queue.add'] = {
			answer: { data: { action: 'add', story: 'S-0001', changed: true }, warnings: [] }
		};
		script['queue.remove'] = {
			answer: { data: { action: 'remove', story: 'S-0001', changed: true }, warnings: [] }
		};
		let r = await post(event);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ action: 'add', story: 'S-0001' });
		r = await del(event);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ action: 'remove' });
		const writes = asked.filter((a) => a.method.startsWith('queue.'));
		expect(writes.map((w) => [w.method, w.params.id])).toEqual([
			['queue.add', 'S-0001'],
			['queue.remove', 'S-0001']
		]);
		for (const w of writes) expect(typeof w.params.request_id).toBe('string');
	});

	it("answers a refusal with flai's sentence and each reason, to show verbatim", async () => {
		const reasons = [{ kind: 'criteria_unticked', text: 'criterion 2 is unticked' }];
		script['queue.add'] = {
			error: new AgentError(422, 'S-0001 is not enqueued: it is not eligible', -32010, {
				story: 'S-0001',
				text: 'S-0001 is not enqueued: it is not eligible',
				reasons
			})
		};
		const r = await post(event);
		expect(r.status).toBe(422);
		const body = await r.json();
		expect(body.error).toContain('not eligible');
		expect(body.reasons).toEqual(reasons);
	});
});
