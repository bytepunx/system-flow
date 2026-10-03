import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted (S-0198): what flai on the host would answer for issue.list and
// issue.story. How the command line is built is flai's business, tested there; here it is what the
// routes ask for and what they do with the answer.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let answers: Record<string, unknown> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	return answers[method];
}) as Ask;

describe('issues routes over the channel', () => {
	let list: Handler;
	let story: Handler;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		list = (await import('./+server')).GET as unknown as Handler;
		story = (await import('./[id]/story/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		answers = {};
	});
	afterAll(() => useRepo(null));

	it('lists every issue, or a story’s own with ?story', async () => {
		answers['issue.list'] = { data: [{ id: 'I-0001' }], warnings: [] };
		const all = await list({ url: new URL('http://x/api/issues') } as never);
		expect(await all.json()).toEqual([{ id: 'I-0001' }]);
		await list({ url: new URL('http://x/api/issues?story=S-0198') } as never);
		const lists = asked.filter((a) => a.method === 'issue.list').map((a) => a.params);
		expect(lists).toEqual([{}, { story: 'S-0198' }]);
	});

	it('makes a story from an issue, with its epic when given, and answers what flai made', async () => {
		answers['issue.story'] = {
			data: { id: 'S-0201', title: 'Fixture was ignored', issue: 'I-0012', committed: true },
			warnings: ['a warning']
		};
		const r = await story({
			params: { id: 'I-0012' },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify({ epic: 'E-0002' }) })
		} as never);
		expect(await r.json()).toMatchObject({ id: 'S-0201', issue: 'I-0012', log: ['a warning'] });
		await story({
			params: { id: 'I-0013' },
			request: new Request('http://x', { method: 'POST', body: '{}' })
		} as never);
		const made = asked.filter((a) => a.method === 'issue.story').map((a) => a.params);
		expect(made[0]).toMatchObject({ id: 'I-0012', epic: 'E-0002' });
		expect(made[1]).toMatchObject({ id: 'I-0013' });
		expect(made[1]).not.toHaveProperty('epic');
		expect(made.every((p) => typeof p.request_id === 'string')).toBe(true);
	});
});
