import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as finalize.test.ts does (S-0282): which boxes change, and that no other byte
// does, is flai's, tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const HASH = 'a'.repeat(64);

describe('criteria route over the channel (S-0282)', () => {
	let post: Handler;
	const send = (id: string, body: unknown) =>
		post({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
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

	it('asks flai to tick the criteria named, with the hash, naming nobody, and passes on its answer', async () => {
		script = {
			'item.criteria': {
				answer: {
					data: { id: 'S-0282', committed: true, criteria: [{ n: 1, text: 'x', ticked: true }] },
					warnings: ['S-0282 has 2 criteria unticked']
				}
			}
		};
		const r = await send('S-0282', { hash: HASH, tick: [1, 3] });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			id: 'S-0282',
			criteria: [{ n: 1, ticked: true }],
			log: ['S-0282 has 2 criteria unticked']
		});
		const write = asked.find((a) => a.method === 'item.criteria');
		expect(write?.params).toMatchObject({ id: 'S-0282', hash: HASH, tick: [1, 3] });
		expect(write?.params).not.toHaveProperty('untick');
		expect(write?.params).not.toHaveProperty('by');
		expect(String(write?.params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it('asks flai to untick', async () => {
		script = { 'item.criteria': { answer: { data: { id: 'S-0282' }, warnings: [] } } };
		const r = await send('S-0282', { hash: HASH, untick: [2], tick: [] });
		expect(r.status).toBe(200);
		const write = asked.find((a) => a.method === 'item.criteria');
		expect(write?.params).toMatchObject({ id: 'S-0282', hash: HASH, untick: [2] });
		expect(write?.params).not.toHaveProperty('tick');
	});

	it.each([
		[{ tick: [1] }, 'hash is required: the one the item was loaded with'],
		[{ hash: HASH }, 'name the criteria to tick or untick by their numbers'],
		[{ hash: HASH, tick: [] }, 'name the criteria to tick or untick by their numbers'],
		[{ hash: HASH, tick: [0] }, 'tick is a list of criterion numbers, from 1'],
		[{ hash: HASH, untick: [1.5] }, 'untick is a list of criterion numbers, from 1'],
		[{ hash: HASH, tick: '1' }, 'tick is a list of criterion numbers, from 1']
	])('refuses %j before asking flai', async (body, error) => {
		const r = await send('S-0282', body);
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(error);
		expect(asked.filter((a) => a.method === 'item.criteria')).toEqual([]);
	});

	it('passes on a conflict with what the item now is', async () => {
		script = {
			'item.criteria': {
				error: new AgentError(409, 'S-0282 changed', -32009, { hash: 'b'.repeat(64), current: 'x' })
			}
		};
		const r = await send('S-0282', { hash: HASH, tick: [1] });
		expect(r.status).toBe(409);
		expect(await r.json()).toMatchObject({ error: 'S-0282 changed', hash: 'b'.repeat(64) });
	});

	it('passes on why flai refused', async () => {
		script = {
			'item.criteria': {
				error: new AgentError(400, 'S-0282 has 3 criteria: there is no 9', -32011)
			}
		};
		const r = await send('S-0282', { hash: HASH, tick: [9] });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe('S-0282 has 3 criteria: there is no 9');
	});
});
