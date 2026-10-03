import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as agent.test.ts does (S-0201): what finalizing changes and who it records is
// flai's, tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

describe('finalize route over the channel (S-0201)', () => {
	let post: Handler;
	const finalize = (id: string) =>
		post({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: '{}' })
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		post = (await import('./+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('asks flai to finalize the story, naming nobody, and passes on its answer', async () => {
		script = {
			'item.finalize': { answer: { data: { id: 'S-0201', draft: false }, warnings: [] } }
		};
		const r = await finalize('S-0201');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ id: 'S-0201', draft: false });
		const write = asked.find((a) => a.method === 'item.finalize');
		expect(write?.params).toMatchObject({ id: 'S-0201' });
		expect(String(write?.params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
		expect(write?.params).not.toHaveProperty('by');
	});

	it('passes on why flai refused', async () => {
		script = {
			'item.finalize': { error: new AgentError(400, 'S-0201 is not a draft', -32011) }
		};
		const r = await finalize('S-0201');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe('S-0201 is not a draft');
	});
});
