import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as finalize.test.ts does (S-0388): what restoring changes is flai's,
// tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const PATH = 'flai/cmd/accept.go';
const REQUIRED = 'path is required: one the branch undoes of main, to restore';
const NO_PATH = [{}, { path: 3 }, { path: '  ' }, { paths: [PATH] }];

describe('restore route over the channel (S-0388)', () => {
	let post: Handler;
	const restore = (id: string, body: unknown) =>
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

	it('asks flai to restore the one path for the story and passes on its answer', async () => {
		script = {
			'stream.restore': { answer: { data: { id: 'S-0388', restored: [PATH] }, warnings: [] } }
		};
		const r = await restore('S-0388', { path: PATH });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ id: 'S-0388', restored: [PATH] });
		const write = asked.find((a) => a.method === 'stream.restore');
		expect(write?.params).toMatchObject({ id: 'S-0388', paths: [PATH] });
		expect(String(write?.params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it.each(NO_PATH)('answers 400 for %j before asking flai', async (body) => {
		const r = await restore('S-0388', body);
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(REQUIRED);
		expect(asked).toEqual([]);
	});

	it('passes on why flai refused, in its words', async () => {
		const why = 'S-0388 does not undo main at flai/cmd/other.go';
		const refusal = new AgentError(422, why, -32010, { undoes_main: [PATH] });
		script = { 'stream.restore': { error: refusal } };
		const r = await restore('S-0388', { path: 'flai/cmd/other.go' });
		expect(r.status).toBe(422);
		expect(await r.json()).toMatchObject({ error: why, undoes_main: [PATH] });
	});
});
