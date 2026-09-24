import { describe, expect, it, vi } from 'vitest';
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
		expect(calls[1].params).toEqual({ archived: true, bodies: true });

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
							age_in_column: '1h',
							age_in_column_seconds: 3600
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
					review: null
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
