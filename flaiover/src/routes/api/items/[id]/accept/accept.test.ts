import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted (S-0041, S-0075): what flai on the host would send while it accepts, as
// progress for the request, then its answer or its error. How the command line is built is flai's
// business and tested there; here it is what the route asks for and what it streams to the browser.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: { progress?: unknown[]; answer?: unknown; error?: AgentError } = {};
const channel = (async (method: string, params: Record<string, unknown> = {}, opt) => {
	asked.push({ method, params });
	for (const step of script.progress ?? []) {
		opt?.onProgress?.(step);
		await new Promise((r) => setTimeout(r, 5));
	}
	if (script.error) throw script.error;
	return script.answer;
}) as Ask;

describe('acceptance and diff endpoints over the channel', () => {
	let accept: Handler;
	let diff: Handler;
	const lines = async (r: Response) =>
		(await r.text())
			.trim()
			.split('\n')
			.map((l) => JSON.parse(l));
	const post = (id: string, body: unknown) =>
		accept({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		accept = (await import('./+server')).POST as unknown as Handler;
		diff = (await import('../diff/+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('asks flai to accept and streams each step before the result', async () => {
		script = {
			progress: [
				{ level: 'INFO', msg: 'story branch merged', branch: 'story/S-0041' },
				{
					level: 'INFO',
					msg: 'acceptance step',
					step: 'merged',
					detail: 'story/S-0041 rebased and fast-forwarded into the main branch'
				},
				{
					level: 'INFO',
					msg: 'acceptance step',
					step: 'committed',
					detail: 'chore: [S-0041] accept and archive'
				},
				{
					level: 'WARN',
					msg: 'accepted locally but not pushed',
					detail: 'run: git push origin HEAD'
				},
				{
					level: 'INFO',
					msg: 'acceptance step',
					step: 'not-pushed',
					detail: 'accepted locally; run: git push origin HEAD'
				}
			],
			answer: {
				data: { id: 'S-0041', status: 'done', merged: true, tags: ['flaiover/v0.13.0'] },
				warnings: []
			}
		};
		const r = await post('S-0041', {});
		expect(r.status).toBe(200);
		expect(r.headers.get('content-type')).toBe('application/x-ndjson');
		const out = await lines(r);
		expect(out.map((l) => l.event)).toEqual([
			'progress',
			'progress',
			'warning',
			'progress',
			'done'
		]);
		expect(out[0]).toMatchObject({
			step: 'merged',
			msg: expect.stringContaining('fast-forwarded')
		});
		expect(out.at(-1).result).toMatchObject({ status: 'done', tags: ['flaiover/v0.13.0'] });
		expect(asked).toHaveLength(1);
		expect(asked[0].method).toBe('accept.run');
		expect(asked[0].params).toMatchObject({ id: 'S-0041', include_uncommitted: false });
		// a write: it carries a request ID, and names nobody: flai decides who the designer is
		expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
		expect(asked[0].params).not.toHaveProperty('by');
	});

	it('passes the choice to include uncommitted files only when it is true', async () => {
		script = { answer: { data: { id: 'S-0041', status: 'done' }, warnings: [] } };
		await (await post('S-0041', { include_uncommitted: true })).text();
		expect(asked[0].params.include_uncommitted).toBe(true);
		await (await post('S-0041', { include_uncommitted: 'yes' })).text();
		expect(asked[1].params.include_uncommitted).toBe(false);
	});

	it('ends with flai’s message verbatim when the acceptance fails', async () => {
		const message =
			'rebase of story/S-0041 onto main stopped with conflicts in docs/users/flai.md; resolve them in .flai-cache/worktrees/S-0041 and run flai stream sync';
		script = { error: new AgentError(502, message, -32603) };
		const out = await lines(await post('S-0041', {}));
		expect(out).toHaveLength(1);
		expect(out[0]).toEqual({ event: 'error', status: 500, error: message });
	});

	it('says the acceptance may have completed when the connection is lost on the way', async () => {
		script = {
			progress: [{ level: 'INFO', msg: 'acceptance step', step: 'merged', detail: 'merged' }],
			error: new AgentError(502, 'the host flai went away before it answered')
		};
		const out = await lines(await post('S-0041', {}));
		expect(out.map((l) => l.event)).toEqual(['progress', 'error']);
		expect(out[1].status).toBe(502);
		expect(out[1].error).toContain('may have completed on the host');
		expect(out[1].error).toContain('flai board');
	});

	it('passes the branch diff through', async () => {
		script = {
			answer: {
				data: {
					story: 'S-0041',
					branch: 'story/S-0041',
					files: [{ path: 'a.md', status: 'added', additions: 1, deletions: 0, patch: '+a' }]
				},
				warnings: []
			}
		};
		const r = await diff({ params: { id: 'S-0041' } } as never);
		expect(r.status).toBe(200);
		expect((await r.json()).files[0]).toMatchObject({ path: 'a.md', status: 'added' });
		expect(asked[0]).toEqual({ method: 'stream.diff', params: { id: 'S-0041' } });
	});

	it('answers a story without a branch with flai’s reason', async () => {
		script = {
			error: new AgentError(
				502,
				'S-0002 has no branch story/S-0002: it was never opened with flai stream open, or it has been merged and removed',
				-32603
			)
		};
		const r = await diff({ params: { id: 'S-0002' } } as never);
		expect(r.status).toBe(500);
		expect((await r.json()).error).toContain('has no branch story/S-0002');
	});
});
