import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import type { ImageJob, ImageLog, ImagesState } from '$lib/images';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

// /api/images and /api/images/retry ask flai's images.get, images.log, and images.ensure (S-0757)
// of a stubbed host: making a real image needs Docker, which the flai side's tests fake too.
type Handler = (event: never) => Promise<Response>;

const ENABLE = 'flai serve enable agent';
const ID = '20261113T093000Z-001';
const REF = 'flai-layer/p:claude-code';

const job = (over: Partial<ImageJob> = {}): ImageJob => ({
	id: ID,
	images: [REF],
	state: 'running',
	ways: [],
	step: 'agent image: building the layer',
	requested: '2026-11-13T09:30:00Z',
	started: '2026-11-13T09:30:01Z',
	log_size: 120,
	...over
});

const state = (over: Partial<ImagesState> = {}): ImagesState => ({
	flai: '1.94.0',
	images: [
		{
			kind: 'layer',
			inner: 'claude-code',
			ref: REF,
			source: 'layer',
			dockerfile: 'agent.Dockerfile',
			present: true,
			id: 'sha256:abc',
			labels: { flai: '1.93.0', recipe: 'r1', base: 'sha256:def' },
			state: 'stale',
			reasons: [{ kind: 'flai', message: 'built for flai 1.93.0, not 1.94.0' }],
			want: { flai: '1.94.0', recipe: 'r1', base: 'sha256:def' }
		}
	],
	running: null,
	queued: [],
	ended: [
		job({
			id: '20261113T080000Z-001',
			state: 'ended',
			step: 'docker build',
			ended: '2026-11-13T08:02:00Z',
			outcome: 'failed',
			cause: 'docker build exited 1',
			duration_ms: 120000
		})
	],
	...over
});

const read = (over: Partial<ImageLog> = {}): ImageLog => ({
	id: ID,
	state: 'running',
	exists: true,
	offset: 12,
	next: 30,
	size: 30,
	text: 'agent image: layer\n',
	...over
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// the reads answer raw, as flai's do; the write as commands do
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return method === 'images.ensure' ? { data: s.data, warnings: [] } : s.data;
}) as Ask;

describe('/api/images (S-0757)', () => {
	let GET: Handler;
	let RETRY: Handler;
	const get = (query = '') => GET({ url: new URL(`http://x/api/images${query}`) } as never);
	const retry = () =>
		RETRY({ request: new Request('http://x/api/images/retry', { method: 'POST' }) } as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
		RETRY = (await import('./retry/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'images.get': { data: state() },
			'images.log': { data: read() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('answers the images and the jobs as flai gives them, asking afresh each time', async () => {
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(state());
		script['images.get'] = { data: state({ running: job(), unjudged: 'docker is not on PATH' }) };
		const again = (await (await get()).json()) as ImagesState;
		expect(again.running).toMatchObject(job());
		expect(again.unjudged).toBe('docker is not on PATH');
		expect(asks('images.get').map((a) => a.params)).toEqual([{}, {}]);
		expect(asks('images.log')).toEqual([]);
	});

	it("reads a job's log from an offset, and from 0 when none is given", async () => {
		const r = await get(`?job=${ID}&offset=12`);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(read());
		expect((await get(`?job=${ID}`)).status).toBe(200);
		expect(asks('images.log').map((a) => a.params)).toEqual([
			{ id: ID, offset: 12 },
			{ id: ID, offset: 0 }
		]);
		expect(asks('images.get')).toEqual([]);
	});

	it('refuses a query that is not a job and an offset with 400 before asking flai', async () => {
		for (const query of [
			'?offset=12',
			'?job=',
			'?job=../x',
			'?job=20261113T093000Z',
			`?job=${ID}&offset=-1`,
			`?job=${ID}&offset=1.5`,
			`?job=${ID}&offset=x`
		]) {
			const r = await get(query);
			expect(r.status, query).toBe(400);
			expect((await r.json()).error, query).toMatch(/^(job|offset) must be/);
		}
		expect(asked.filter((a) => a.method.startsWith('images.'))).toEqual([]);
	});

	it("answers 404 with flai's sentence for a job the project does not have", async () => {
		const none = `this project has no image job ${ID}: give one images.get lists`;
		script['images.log'] = { error: new AgentError(502, none, -32004) };
		const r = await get(`?job=${ID}`);
		expect(r.status).toBe(404);
		expect((await r.json()).error).toBe(none);
	});

	it('starts a job with a request ID, then answers the images and jobs with it started', async () => {
		const queued: ImageJob = {
			id: ID,
			images: [REF],
			state: 'queued',
			ways: [],
			requested: '2026-11-13T09:30:00Z',
			log_size: 0
		};
		script['images.ensure'] = { data: { ...queued, ahead: 2 } };
		script['images.get'] = { data: state({ queued: [queued] }) };
		const r = await retry();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			...state({ queued: [queued] }),
			started: { ...queued, ahead: 2 }
		});
		expect(asks('images.ensure')).toHaveLength(1);
		const { request_id, ...rest } = asks('images.ensure')[0].params;
		expect(typeof request_id).toBe('string');
		expect(rest).toEqual({});
		expect(asks('images.get')).toHaveLength(1);
	});

	it('answers 403 with the command that enables the agent action while it is off', async () => {
		script['images.ensure'] = {
			error: new AgentError(
				502,
				`the host action "agent" is not enabled for this project. On the host, in the project, run: ${ENABLE}`,
				-32012,
				{ action: 'agent', enable: ENABLE }
			)
		};
		const r = await retry();
		expect(r.status).toBe(403);
		const body = await r.json();
		expect(body).toMatchObject({ action: 'agent', enable: ENABLE });
		expect(body.error).toContain(ENABLE);
		expect(asks('images.get')).toEqual([]);
	});

	it("answers a busy job's refusal word for word as 400, with the job", async () => {
		const busy = `the image job ${ID}, which makes ${REF}, runs since 2026-11-13T09:30:01Z, so no other job is started for the same images: follow its output with images.log, and ask again once it has ended`;
		script['images.ensure'] = { error: new AgentError(502, busy, -32011, { job: job() }) };
		const r = await retry();
		expect(r.status).toBe(400);
		expect(await r.json()).toMatchObject({ error: busy, job: job() });
		expect(asks('images.get')).toEqual([]);
	});

	it("answers flai's refusal when the agents need no image word for word as 400", async () => {
		const none =
			"the project's agents need no agent image flai makes: none runs on the container harness and its manifest names no agent.image.layer, or each image they run in is the operator's own, named with flai serve agent harness container --image, which flai never makes; so no image job is started";
		script['images.ensure'] = { error: new AgentError(502, none, -32011) };
		const r = await retry();
		expect(r.status).toBe(400);
		const body = await r.json();
		expect(body.error).toBe(none);
		expect(body.job).toBeUndefined();
	});

	it('answers 501 telling to upgrade a host flai without the images methods', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		for (const r of [await get(), await get(`?job=${ID}&offset=0`), await retry()]) {
			expect(r.status).toBe(501);
			const error = (await r.json()).error as string;
			expect(error).toContain('the host flai is too old');
			expect(error).toContain('flai self-upgrade');
		}
	});
});
