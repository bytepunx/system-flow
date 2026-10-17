import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted (S-0394, as checks.test.ts scripts it): what flai on the host would answer
// for forge.status. Reading the kept state is flai's business, tested there; here it is what the
// route asks for and what it does with the answer.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const PULL_REQUEST = {
	story: 'S-0394',
	pull_request: 7,
	url: 'https://github.com/bytepunx/system-flow/pull/7',
	head_sha: '0123456789abcdef0123456789abcdef01234567',
	read_at: '2026-10-17T10:00:00Z',
	checks: [
		{
			name: 'flai',
			state: 'failed',
			url: 'https://github.com/bytepunx/system-flow/actions/runs/1',
			summary: 'go test failed',
			since: '2026-10-17T09:58:00Z'
		}
	],
	reviews: [
		{
			author: 'alex',
			state: 'changes-requested',
			url: 'https://github.com/bytepunx/system-flow/pull/7#pullrequestreview-1',
			body: 'Name the check.',
			submitted_at: '2026-10-17T09:59:00Z'
		}
	]
};

describe('the forge route over the channel', () => {
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

	it("answers the story's pull request, its checks and its reviews, as flai kept them", async () => {
		script = {
			'forge.status': {
				answer: { data: { story: 'S-0394', pull_request: PULL_REQUEST }, warnings: [] }
			}
		};
		const r = await doGet('S-0394');
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ story: 'S-0394', pull_request: PULL_REQUEST });
		expect(asked[0]).toEqual({ method: 'forge.status', params: { id: 'S-0394' } });
	});

	it('answers a null pull request for a story flai serve has read none for', async () => {
		script = {
			'forge.status': { answer: { data: { story: 'S-0001', pull_request: null }, warnings: [] } }
		};
		const r = await doGet('S-0001');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ story: 'S-0001', pull_request: null });
	});

	it("answers flai's refusal with its status and words", async () => {
		script = {
			'forge.status': {
				error: new AgentError(
					400,
					"T-0001 is a task; forge.status reads a story's pull request, so give a story's ID, such as S-0001",
					-32602
				)
			}
		};
		const r = await doGet('T-0001');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain("give a story's ID");
	});
});
