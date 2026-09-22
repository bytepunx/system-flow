import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted (S-0082, following accept.test.ts's own pattern): what flai on the host
// would answer for checks.status, checks.run, checks.cancel, and checks.tail. How the command line
// is built and what it runs is flai's business, tested there; here it is what the route asks for
// and what it does with the answer.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { progress?: unknown[]; answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}, opt) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	for (const step of s.progress ?? []) {
		opt?.onProgress?.(step);
		await new Promise((r) => setTimeout(r, 5));
	}
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

describe('checks routes over the channel', () => {
	let get: Handler;
	let post: Handler;
	let tailGet: Handler;
	const lines = async (r: Response) =>
		(await r.text())
			.trim()
			.split('\n')
			.filter((l) => l)
			.map((l) => JSON.parse(l));
	const doGet = (id: string) => get({ params: { id } } as never);
	const doPost = (id: string, body: unknown) =>
		post({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);
	const doTail = (id: string, from?: number) =>
		tailGet({
			params: { id },
			url: new URL(
				`http://x/api/items/${id}/checks/tail${from === undefined ? '' : `?from=${from}`}`
			)
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = mod.GET as unknown as Handler;
		post = mod.POST as unknown as Handler;
		tailGet = (await import('./tail/+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('carries checks_enabled from project.info alongside the run', async () => {
		script = {
			'checks.status': {
				answer: {
					data: {
						story: 'S-0082',
						running: true,
						current: 'flai',
						started: '2026-09-22T00:00:00Z'
					},
					warnings: []
				}
			},
			'project.info': { answer: { host_actions: { checks: true, push: false } } }
		};
		const r = await doGet('S-0082');
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ story: 'S-0082', running: true, checks_enabled: true });
		// respond() also asks project.info for the project header (ADR-0024); this route asks
		// it again itself for checks_enabled, so it appears at least twice, same as dashboard's.
		expect(asked[0]).toEqual({ method: 'checks.status', params: { id: 'S-0082' } });
		expect(asked.map((a) => a.method)).toContain('project.info');
	});

	it('says checks_enabled is false rather than fail when project.info cannot be asked', async () => {
		script = {
			'checks.status': { answer: { data: { story: 'S-0082', running: false }, warnings: [] } },
			'project.info': { error: new AgentError(503, 'no flai connected') }
		};
		const body = await (await doGet('S-0082')).json();
		expect(body).toMatchObject({ running: false, checks_enabled: false });
	});

	it('refuses a POST with no recognised action', async () => {
		const r = await doPost('S-0082', { action: 'bogus' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('run, cancel');
	});

	it('asks checks.run for run, with a request id, and answers what flai answered', async () => {
		script = {
			'checks.run': {
				answer: {
					data: { story: 'S-0082', running: false, outcome: 'passed' },
					warnings: ['a warning']
				}
			}
		};
		const r = await doPost('S-0082', { action: 'run' });
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ story: 'S-0082', outcome: 'passed', log: ['a warning'] });
		expect(asked[0].method).toBe('checks.run');
		expect(asked[0].params).toMatchObject({ id: 'S-0082' });
		expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it('asks checks.cancel for cancel', async () => {
		script = {
			'checks.cancel': {
				answer: { data: { story: 'S-0082', running: false, outcome: 'cancelled' }, warnings: [] }
			}
		};
		const body = await (await doPost('S-0082', { action: 'cancel' })).json();
		expect(body).toMatchObject({ outcome: 'cancelled' });
		expect(asked[0].method).toBe('checks.cancel');
	});

	it('is 403 with what enables it when the checks host action is off', async () => {
		script = {
			'checks.run': {
				error: new AgentError(
					403,
					'the host action "checks" is not enabled for this project. On the host, in the project, run: flai serve enable checks',
					-32012
				)
			}
		};
		const r = await doPost('S-0082', { action: 'run' });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable checks');
	});

	it('streams new lines as they arrive, then a done event with the offset and whether it is still running', async () => {
		script = {
			'checks.tail': {
				progress: [
					{ msg: 'line', text: '=== flai: scripts/flai-test.sh ===' },
					{ msg: 'line', text: 'ok  github.com/bytepunx/system-flow/flai/internal/serve' }
				],
				answer: { data: { offset: 512, running: true }, warnings: [] }
			}
		};
		const out = await lines(await doTail('S-0082', 0));
		expect(out.map((l) => l.event)).toEqual(['line', 'line', 'done']);
		expect(out[0].text).toContain('flai: scripts/flai-test.sh');
		expect(out[2]).toMatchObject({ offset: 512, running: true });
		expect(asked[0].method).toBe('checks.tail');
		expect(asked[0].params).toMatchObject({ id: 'S-0082', from: 0 });
	});

	it('reads the offset from the query string', async () => {
		script = {
			'checks.tail': {
				answer: { data: { offset: 900, running: false, outcome: 'passed' }, warnings: [] }
			}
		};
		const out = await lines(await doTail('S-0082', 512));
		expect(asked[0].params).toMatchObject({ from: 512 });
		expect(out[0]).toMatchObject({ event: 'done', offset: 900, running: false, outcome: 'passed' });
	});

	it('treats a negative or missing from as zero, not a negative offset', async () => {
		script = { 'checks.tail': { answer: { data: { offset: 10, running: false }, warnings: [] } } };
		await lines(await doTail('S-0082', -50));
		expect(asked[0].params.from).toBe(0);
	});

	it('ends the stream with an error event rather than throwing when flai refuses the tail', async () => {
		script = {
			'checks.tail': { error: new AgentError(404, 'no checks have run for S-0002', -32004) }
		};
		const out = await lines(await doTail('S-0002', 0));
		expect(out).toEqual([{ event: 'error', status: 404, error: 'no checks have run for S-0002' }]);
	});
});
