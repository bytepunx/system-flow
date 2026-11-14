import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as restore.test.ts does (S-0788): what sending back changes is flai's,
// tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const SENT_BACK = {
	id: 'S-0788',
	reason:
		'S-0788 was sent back to in-progress: alex asked its agent to rebase story/S-0788 onto main and resolve 1 blocker, on TH-0001 (ADR-0302)',
	sent_back: true,
	branch: 'story/S-0788',
	base: 'main',
	conflicts: [],
	thread: 'TH-0001',
	blockers: ['story/S-0788 is 2 commits behind main: rebase it']
};

describe('resolve route over the channel (S-0788)', () => {
	let post: Handler;
	const resolveStory = (id: string) =>
		post({ params: { id }, request: new Request('http://x', { method: 'POST' }) } as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		post = (await import('./+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('asks flai to send the story back to resolve and passes on the send back', async () => {
		script = { 'accept.resolve': { answer: { data: SENT_BACK, warnings: [] } } };
		const r = await resolveStory('S-0788');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(SENT_BACK);
		const write = asked.find((a) => a.method === 'accept.resolve');
		expect(write?.params).toMatchObject({ id: 'S-0788' });
		expect(String(write?.params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it('passes on why flai refused, in its words', async () => {
		const why =
			"nothing holds S-0788's acceptance, so its agent has nothing to resolve: accept it with flai accept S-0788";
		const refusal = new AgentError(422, why, -32010, { refused: { id: 'S-0788', reason: why } });
		script = { 'accept.resolve': { error: refusal } };
		const r = await resolveStory('S-0788');
		expect(r.status).toBe(422);
		expect(await r.json()).toMatchObject({ error: why, refused: { id: 'S-0788', reason: why } });
	});
});
