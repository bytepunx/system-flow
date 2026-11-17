import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as restore.test.ts does (S-0788): what the job does is flai's, tested
// there; here it is what the route asks for and how it passes flai's answer on (S-0804).
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const LOG = '/repo/.flai-cache/verify/jobs/S-0788.log';
const STARTED = {
	id: 'S-0788',
	job: 'S-0788-20261116T223000Z',
	started: true,
	running: false,
	log: LOG,
	reason: `rebase and resolve job S-0788-20261116T223000Z started for S-0788: flai serve rebases story/S-0788 onto the main branch and closes it out at the rebased head; its output is in ${LOG}`
};
const JOB = {
	id: 'S-0788-20261116T223000Z',
	story: 'S-0788',
	kind: 'resolve',
	by: 'alex',
	started: '2026-11-16T22:30:00Z',
	step: 'flai verify',
	tier: 'go-test',
	steps: [{ name: 'close-out', kind: 'step', source: 'resolve', started: '2026-11-16T22:30:02Z' }],
	log: LOG,
	running: true
};

describe('resolve route over the channel (S-0788, S-0804)', () => {
	let post: Handler;
	let get: Handler;
	const resolveStory = (id: string) =>
		post({ params: { id }, request: new Request('http://x', { method: 'POST' }) } as never);
	const jobOf = (id: string) => get({ params: { id } } as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const route = await import('./+server');
		post = route.POST as unknown as Handler;
		get = route.GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('asks flai for the rebase and resolve job and passes on its answer', async () => {
		script = { 'accept.resolve': { answer: { data: STARTED, warnings: [] } } };
		const r = await resolveStory('S-0788');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(STARTED);
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

	it('reads the story’s last job with resolve.status, and asks nothing of a write', async () => {
		script = { 'resolve.status': { answer: { data: { id: 'S-0788', job: JOB }, warnings: [] } } };
		const r = await jobOf('S-0788');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ id: 'S-0788', job: JOB });
		expect(asked).toEqual([{ method: 'resolve.status', params: { id: 'S-0788' } }]);
	});

	it('answers the story alone when it has had no job', async () => {
		script = { 'resolve.status': { answer: { data: { id: 'S-0788' }, warnings: [] } } };
		const r = await jobOf('S-0788');
		const body = await r.json();
		expect(body).toMatchObject({ id: 'S-0788' });
		expect(body.job).toBeUndefined();
	});

	it('answers 501 in flai’s words when the host flai is older than resolve.status', async () => {
		const why = 'method not found: resolve.status';
		script = { 'resolve.status': { error: new AgentError(501, why, -32601) } };
		const r = await jobOf('S-0788');
		expect(r.status).toBe(501);
		expect((await r.json()).error).toBe(why);
	});
});
