import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as checks.test.ts scripts it: what flai on the host would answer for
// verify.status. How flai reads the stored report is flai's business, tested there; here it is what
// the route asks for and what it does with the answer.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const REPORT = {
	story: 'S-0270',
	commit: '0123456789abcdef',
	base: 'main',
	ran_at: '2026-10-06T22:00:00Z',
	duration_ms: 5100,
	duration: '5.1s',
	passed: true,
	steps: [{ name: 'rebase', state: 'passed', duration_ms: 3, duration: '3ms' }]
};

describe('verify route over the channel', () => {
	let get: Handler;
	const doGet = (id: string) => get({ params: { id } } as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		get = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('answers the story’s stored report from verify.status', async () => {
		script = { 'verify.status': { answer: { data: REPORT, warnings: [] } } };
		const r = await doGet('S-0270');
		expect(r.status).toBe(200);
		expect((await r.json()).report).toEqual(REPORT);
		expect(asked[0]).toEqual({ method: 'verify.status', params: { id: 'S-0270' } });
	});

	it('answers a null report when the story has none', async () => {
		script = { 'verify.status': { answer: { data: null, warnings: [] } } };
		const r = await doGet('S-0270');
		expect(r.status).toBe(200);
		expect((await r.json()).report).toBeNull();
	});

	it('answers a null report, not an error, when no flai is connected', async () => {
		script = {
			'verify.status': {
				error: new AgentError(
					503,
					'no host flai is connected; run flai dashboard in the project, or flai serve start'
				)
			}
		};
		const r = await doGet('S-0270');
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.report).toBeNull();
		expect(body.error).toBeUndefined();
	});

	it('says what flai said when it cannot read the stored report', async () => {
		script = {
			'verify.status': {
				error: new AgentError(
					502,
					'read the last verify report of S-0270: unexpected end of JSON input; delete it and verify again',
					-32603
				)
			}
		};
		const r = await doGet('S-0270');
		expect(r.status).toBe(500);
		expect((await r.json()).error).toContain('delete it and verify again');
	});
});
