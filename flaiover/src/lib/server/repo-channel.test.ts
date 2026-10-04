import { describe, expect, it, vi } from 'vitest';
import { doneLane } from '$lib/publish';
import { AgentError } from './agent';
import { board } from './board';
import { Repo, RepoError, type Ask } from './repo';

// S-0073: the project, items, threads, and the board are asked of flai on the host.
const project = { name: 'Harbour', key: 'harbour', layout: { design: 'd', docs: 'o', wip: 'w' } };

function fake(answers: Record<string, unknown>) {
	const calls: { method: string; params: Record<string, unknown> }[] = [];
	const ask = vi.fn(async (method: string, params: Record<string, unknown> = {}) => {
		calls.push({ method, params });
		const a = answers[method];
		if (a instanceof Error) throw a;
		return a;
	}) as unknown as Ask;
	return { ask, calls };
}

describe('Repo over the channel', () => {
	it('asks once and keeps the answer until flai says a file changed', async () => {
		const { ask, calls } = fake({ 'project.info': project, 'items.list': [] });
		const r = new Repo('/nowhere', ask);
		const seen: string[] = [];
		r.on('change', (p) => seen.push(p));
		await r.manifest();
		await r.layout();
		await r.items();
		await r.items();
		expect(calls.map((c) => c.method)).toEqual(['project.info', 'items.list']);
		expect(calls[1].params).toEqual({ archived: false });

		// a flai that went away: nothing it said is shown as current
		r.forget();
		await r.manifest();
		expect(calls.at(-1)?.method).toBe('project.info');
		calls.pop();

		r.changed('wip/kanban/stories/S-0001-a.md');
		expect(seen).toEqual(['wip/kanban/stories/S-0001-a.md']);
		await r.items();
		expect(calls.map((c) => c.method)).toEqual(['project.info', 'items.list', 'items.list']);
	});

	it('asks for the archive and a type only when asked, never for bodies, each query kept on its own (S-0162)', async () => {
		const { ask, calls } = fake({ 'items.list': [], 'items.count': { active: 3, archived: 9 } });
		const r = new Repo('/nowhere', ask);
		await r.items();
		await r.items({ archive: true });
		await r.items({ type: 'epic', archive: true });
		await r.items({ type: 'epic', status: 'ready' });
		await r.items({ archive: true });
		expect(await r.itemCount()).toEqual({ active: 3, archived: 9 });
		expect(calls.map((c) => c.params)).toEqual([
			{ archived: false },
			{ archived: true },
			{ archived: true, type: 'epic' },
			{ archived: false, type: 'epic', status: 'ready' },
			{}
		]);
		expect(calls.at(-1)?.method).toBe('items.count');
	});

	it('gives an item the shape the pages expect from what flai marshals', async () => {
		const flaiItem = {
			id: 'S-0001',
			type: 'story',
			nature: 'feature',
			title: 'Berths',
			status: 'ready',
			parent: '',
			owner: '',
			created: '2026-09-20T08:00:00Z',
			updated: '2026-09-20T08:01:00Z',
			transitions: null,
			blocked: [
				{ from: '2026-09-20T08:02:00Z', until: '', reason: 'tide' },
				{ from: '2026-09-19T08:00:00Z', until: '2026-09-19T09:00:00Z', reason: 'fog' }
			],
			estimate: '',
			stream: '',
			tags: null,
			path: 'wip/kanban/stories/S-0001-berths.md',
			archived: false,
			body: '# S-0001\n'
		};
		const { ask, calls } = fake({ 'item.get': { item: flaiItem, children: null } });
		const { item, children } = await new Repo('/nowhere', ask).itemById(' S-1 ');
		expect(calls[0]).toEqual({ method: 'item.get', params: { id: 'S-1' } });
		expect(children).toEqual([]);
		expect(item.parent).toBeUndefined();
		expect(item.transitions).toEqual([]);
		expect(item.tags).toBeUndefined();
		expect(item.blocked?.map((b) => b.until)).toEqual([undefined, '2026-09-19T09:00:00Z']);
	});

	it('carries a story’s task plan from item.get, and none when flai sent none (S-0176)', async () => {
		const flaiItem = {
			id: 'S-0001',
			type: 'story',
			nature: 'feature',
			title: 'Berths',
			status: 'in-progress',
			created: '2026-09-20T08:00:00Z',
			updated: '2026-09-20T08:01:00Z',
			transitions: null,
			blocked: null,
			tags: null,
			path: 'wip/kanban/stories/S-0001-berths.md',
			archived: false,
			body: '# S-0001\n'
		};
		const plan = {
			tasks: [
				{ id: 'T-0001', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] },
				{ id: 'T-0002', state: 'in-progress' },
				{ id: 'T-0003', state: 'ready' }
			],
			layers: [['T-0002', 'T-0003'], ['T-0001']]
		};
		const got = await new Repo(
			'/nowhere',
			fake({ 'item.get': { item: flaiItem, children: null, plan } }).ask
		).itemById('S-0001');
		expect(got.plan).toEqual(plan);
		const none = await new Repo(
			'/nowhere',
			fake({ 'item.get': { item: flaiItem, children: null } }).ask
		).itemById('S-0001');
		expect(none).not.toHaveProperty('plan');
	});

	it('answers with the status the route should give: no flai, not found, a bad argument', async () => {
		const cases: [Error, number][] = [
			[new AgentError(503, 'no host flai is connected; run flai dashboard'), 503],
			[new AgentError(502, 'S-0099 not found', -32004), 404],
			[new AgentError(502, '"x" is not a work item ID', -32602), 400],
			[new AgentError(504, 'the host flai did not answer item.get in time'), 504]
		];
		for (const [err, status] of cases) {
			const r = new Repo('/nowhere', fake({ 'item.get': err }).ask);
			const got = await r.itemById('S-0099').catch((e) => e);
			expect(got).toBeInstanceOf(RepoError);
			expect(got).toMatchObject({ status, message: err.message });
		}
	});

	it('does not keep a failure: the next question is asked again', async () => {
		let n = 0;
		const ask = (async () => {
			if (n++ === 0) throw new AgentError(503, 'no host flai is connected');
			return project;
		}) as unknown as Ask;
		const r = new Repo('/nowhere', ask);
		await expect(r.manifest()).rejects.toMatchObject({ status: 503 });
		await expect(r.manifest()).resolves.toMatchObject({ name: 'Harbour' });
	});

	it('asks threads with the resolved ones, and by anchor without a trailing slash', async () => {
		const { ask, calls } = fake({ 'threads.list': [] });
		const r = new Repo('/nowhere', ask);
		await r.threads();
		await r.threadsFor('design/system/');
		expect(calls.map((c) => c.params)).toEqual([{ all: true }, { on: 'design/system', all: true }]);
	});

	it('lays the board out as flai gave it, counting age from entered_at', async () => {
		const { ask, calls } = fake({
			'board.get': {
				columns: {
					ready: [
						{
							id: 'S-0002',
							type: 'story',
							title: 'Cranes',
							nature: 'feature',
							parent: 'E-0001',
							parent_title: 'Quay',
							status: 'ready',
							entered_at: '2026-09-20T08:00:00Z',
							blocked: false,
							draft: true,
							age_in_column: '1h',
							age_in_column_seconds: 3600,
							tasks: { ready: 1, waiting: 2, in_progress: 1, done: 3, layers: 3 }
						},
						{
							id: 'S-0001',
							type: 'story',
							title: 'Berths',
							nature: 'feature',
							status: 'ready',
							entered_at: '2026-09-20T07:00:00Z',
							blocked: true,
							age_in_column: '2h',
							age_in_column_seconds: 7200
						}
					],
					review: null,
					done: [
						{
							id: 'S-0003',
							type: 'story',
							title: 'Moorings',
							nature: 'feature',
							status: 'done',
							entered_at: '2026-09-19T08:00:00Z',
							blocked: false,
							archived: true,
							age_in_column: '1d',
							age_in_column_seconds: 90000
						},
						{
							id: 'T-0001',
							type: 'task',
							title: 'Bollards',
							nature: 'feature',
							parent: 'S-0002',
							parent_title: 'Cranes',
							status: 'done',
							entered_at: '2026-09-20T06:00:00Z',
							blocked: false,
							age_in_column: '3h',
							age_in_column_seconds: 10800
						}
					]
				},
				wip_limits: { ready: 5, 'in-progress': 2, review: 3 },
				order: ['S-0002', 'S-0001'],
				breaches: null
			}
		});
		const b = await board(new Repo('/nowhere', ask), new Date('2026-09-20T09:00:00Z'));
		expect(calls[0]).toEqual({ method: 'board.get', params: { all: true } });
		expect(calls).toHaveLength(1);
		expect(b.columns.ready.map((c) => [c.id, c.age_seconds, c.parent_title, c.blocked])).toEqual([
			['S-0002', 3600, 'Quay', false],
			['S-0001', 7200, undefined, true]
		]);
		expect(Object.keys(b.columns)).toEqual([
			'backlog',
			'ready',
			'in-progress',
			'review',
			'done',
			'cancelled'
		]);
		expect(b.columns.review).toEqual([]);
		expect(b.order).toEqual(['S-0002', 'S-0001']);
		// a story's tasks by state and its plan's layers (S-0176); none for a story without tasks
		expect(b.columns.ready[0].tasks).toEqual({
			ready: 1,
			waiting: 2,
			in_progress: 1,
			done: 3,
			layers: 3
		});
		expect(b.columns.ready[1]).not.toHaveProperty('tasks');
		// a draft story says so (S-0201); flai omits the flag otherwise, and so does the card
		expect(b.columns.ready[0].draft).toBe(true);
		expect(b.columns.ready[1]).not.toHaveProperty('draft');
		// an archived done card says so, or the done lane cannot leave it out while the clone lags its
		// remote's tags (I-0060); flai omits the flag otherwise, and so does the card
		expect(b.columns.done[0].archived).toBe(true);
		expect(b.columns.done[1]).not.toHaveProperty('archived');
		const tagsBehind = {
			remote: 'origin',
			behind: [{ component: 'flai', remote: 'flai/v1.0.0' }],
			fix: 'git fetch --tags',
			message: 'behind'
		};
		expect(doneLane(b.columns.done, tagsBehind).map((c) => c.id)).toEqual(['T-0001']);
	});

	it('repeats a write once, with the same request ID, when flai was lost and came back', async () => {
		const sent: Record<string, unknown>[] = [];
		let lose = true;
		const ask = (async (_m: string, params: Record<string, unknown> = {}) => {
			sent.push(params);
			if (lose) {
				lose = false;
				throw new AgentError(502, 'the host flai went away before it answered');
			}
			return { data: { id: 'S-0001', status: 'ready' }, warnings: [] };
		}) as unknown as Ask;
		const r = new Repo('/nowhere', ask, async () => true);
		await expect(r.write('item.move', { id: 'S-0001', to: 'ready' })).resolves.toMatchObject({
			data: { status: 'ready' }
		});
		expect(sent).toHaveLength(2);
		expect(sent[0].request_id).toMatch(/^[0-9a-f-]{36}$/);
		expect(sent[1].request_id).toBe(sent[0].request_id);

		// flai did not come back: the loss is the answer, and nothing is sent again
		lose = true;
		sent.length = 0;
		const gone = new Repo('/nowhere', ask, async () => false);
		await expect(gone.write('item.move', { id: 'S-0001', to: 'ready' })).rejects.toMatchObject({
			status: 502
		});
		expect(sent).toHaveLength(1);

		// what flai refused is not retried: a rule is an answer
		sent.length = 0;
		const refusing = (async (_m: string, params: Record<string, unknown> = {}) => {
			sent.push(params);
			throw new AgentError(502, 'S-0001 cannot go from review to cancelled', -32011);
		}) as unknown as Ask;
		await expect(
			new Repo('/nowhere', refusing, async () => true).write('item.move', { id: 'S-0001', to: 'x' })
		).rejects.toMatchObject({ status: 400 });
		expect(sent).toHaveLength(1);

		// a write whose success ends the connection is not repeated: the flai that came back has
		// no record of it and would do it again (S-0107, a serve restart that ran twice)
		lose = true;
		sent.length = 0;
		await expect(
			new Repo('/nowhere', ask, async () => true).write(
				'host.process',
				{ process: 'serve', action: 'restart' },
				{ retry: false }
			)
		).rejects.toMatchObject({ status: 502 });
		expect(sent).toHaveLength(1);
	});

	it('gives a conflict and a refusal their status and what flai sent with them', async () => {
		const conflict = (async () => {
			throw new AgentError(502, 'the file changed', -32009, { hash: 'h2', current: '# theirs' });
		}) as unknown as Ask;
		await expect(new Repo('/nowhere', conflict).write('doc.save', {})).rejects.toMatchObject({
			status: 409,
			data: { hash: 'h2' }
		});
		const refused = (async () => {
			throw new AgentError(502, 'flai check has 1 finding(s)', -32010, { findings: [{}] });
		}) as unknown as Ask;
		await expect(new Repo('/nowhere', refused).write('item.new', {})).rejects.toMatchObject({
			status: 422,
			data: { findings: [{}] }
		});
	});
});
