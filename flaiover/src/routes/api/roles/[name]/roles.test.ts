// A role page's read and its Run, Stop, and Start over the channel (S-0420), scripted as the
// orchestrator route's test is, from a fixture agent.status whose roles hold a defined role,
// reviewer, beside the built-in ones: what flai answers is flai's, tested there; here it is what the
// route asks for and how it answers.

import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const reviewing = {
	story: '',
	item: 'S-0001',
	role: 'review',
	role_name: 'reviewer',
	agent: 'reviewer',
	command: 'claude',
	pid: 42,
	started: '2026-10-28T10:00:00Z'
};
const reviewed = {
	...reviewing,
	item: 'S-0002',
	pid: 41,
	started: '2026-10-27T10:00:00Z',
	ended: '2026-10-27T10:20:00Z',
	outcome: 'worked'
};

const builtin = (name: string, role: string, action: string, held = false) => ({
	name,
	role,
	path: `design/workflows/roles/${name}.md`,
	item: name === 'planner' ? ['epic', 'story'] : null,
	trigger: ['asked'],
	action: { name: action, enabled: true },
	run: null,
	running: null,
	runs: null,
	held
});

/** agent.status as flai answers it since S-0420, with the reviewer reviewing S-0001. */
const status = () => ({
	enabled: true,
	state: {
		command: 'claude',
		roles: {
			planner: builtin('planner', 'plan', 'plan'),
			orchestrator: builtin('orchestrator', 'orchestrate', 'orchestrate'),
			analyzer: builtin('analyzer', 'analyze', 'analyze'),
			curator: builtin('curator', 'curate', 'curate'),
			reviewer: {
				name: 'reviewer',
				role: 'review',
				path: 'design/workflows/roles/reviewer.md',
				item: ['story'],
				trigger: ['asked', 'schedule'],
				action: { name: 'review', enabled: true },
				run: reviewing,
				running: [reviewing],
				runs: [reviewed],
				held: false
			}
		}
	}
});

const document = {
	kind: 'reviewer',
	accrued_cost: 0.9,
	accrued_seconds: 1200,
	tasks_completed: 1,
	last_run: '2026-10-27T10:20:00Z',
	path: 'wip/agents/reviewer.md',
	entries: [
		{
			at: '2026-10-27T10:20:00Z',
			summary: 'reviewed S-0002: two findings',
			items: null,
			seconds: 1200,
			cost: 0.9,
			estimated: false
		}
	]
};

const off = new AgentError(
	403,
	'the host action "review" is not enabled for this project. On the host, in the project, run: flai serve enable review',
	-32012
);

describe('/api/roles/[name]', () => {
	let get: Handler;
	let post: Handler;
	const doGet = (name: string) => get({ params: { name } } as never);
	const doPost = (name: string, body: unknown) =>
		post({
			params: { name },
			request: new Request('http://x', {
				method: 'POST',
				body: typeof body === 'string' ? body : JSON.stringify(body)
			})
		} as never);
	const writes = () =>
		asked.filter((a) => a.method.endsWith('.run') || a.method.startsWith('orchestrate.'));

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = mod.GET as unknown as Handler;
		post = mod.POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'agent.status': { answer: status() },
			'activity.document': { answer: document }
		};
	});
	afterAll(() => useRepo(null));

	describe('GET', () => {
		it('answers the role agent.status reports by the name, with its activity document', async () => {
			const r = await doGet('reviewer');
			expect(r.status).toBe(200);
			const body = await r.json();
			expect(body.role).toEqual(status().state.roles.reviewer);
			expect(body.activity).toMatchObject({ kind: 'reviewer', path: 'wip/agents/reviewer.md' });
			expect(body.activity.entries[0].items).toEqual([]);
			expect(asked.find((a) => a.method === 'activity.document')?.params).toEqual({
				kind: 'reviewer'
			});
			// the page lists the role's ended runs, so it asks for them (S-0611)
			expect(asked.find((a) => a.method === 'agent.status')?.params).toEqual({ finished: true });
		});

		it("reads a role's lists flai answered null as empty", async () => {
			script['activity.document'] = { answer: { ...document, kind: 'curator' } };
			const body = await (await doGet('curator')).json();
			expect(body.role).toMatchObject({ item: [], running: [], runs: [], run: null });
		});

		it('is 404 naming the roles the project runs for a role it does not', async () => {
			const r = await doGet('auditor');
			expect(r.status).toBe(404);
			expect((await r.json()).error).toBe(
				'the project runs no role named auditor; it runs analyzer, curator, orchestrator, planner, reviewer'
			);
		});

		it('is 404 saying flai may be older when agent.status carries no roles', async () => {
			script['agent.status'] = { answer: { enabled: true, state: { command: 'claude' } } };
			const r = await doGet('reviewer');
			expect(r.status).toBe(404);
			expect((await r.json()).error).toContain('older than this dashboard');
		});

		it("is flai's error when agent.status or the activity document cannot be read", async () => {
			script['agent.status'] = { error: new AgentError(503, 'no flai connected') };
			let r = await doGet('reviewer');
			expect(r.status).toBe(503);
			script['agent.status'] = { answer: status() };
			script['activity.document'] = { error: new AgentError(503, 'no flai connected') };
			r = await doGet('reviewer');
			expect(r.status).toBe(503);
			expect((await r.json()).error).toContain('no flai connected');
		});

		it("refuses a name that is not a role's, asking flai nothing", async () => {
			for (const name of ['Reviewer', '../x', 'a b', '']) {
				const r = await doGet(name);
				expect(r.status).toBe(400);
				expect((await r.json()).error).toContain("is not a role's name");
			}
			expect(asked.filter((a) => a.method !== 'project.info')).toEqual([]);
		});
	});

	describe('POST', () => {
		it("runs the reviewer through role.run on the item given and answers flai's run", async () => {
			const started = {
				role: 'review',
				item: 'S-0003',
				agent: 'reviewer',
				command: 'claude',
				pid: 43
			};
			script['role.run'] = { answer: { data: started, warnings: [] } };
			const r = await doPost('reviewer', { action: 'run', item: ' S-0003 ' });
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject(started);
			const [w] = writes();
			expect(w.method).toBe('role.run');
			expect(w.params).toMatchObject({ name: 'reviewer', item: 'S-0003' });
			expect(typeof w.params.request_id).toBe('string');
		});

		it('runs a built-in by its own method', async () => {
			script['plan.run'] = { answer: { data: { item: 'E-0016' }, warnings: [] } };
			script['analyze.run'] = { answer: { data: { focus: 'all' }, warnings: [] } };
			script['curate.run'] = { answer: { data: {}, warnings: [] } };
			expect((await doPost('planner', { action: 'run', item: 'E-0016' })).status).toBe(200);
			expect((await doPost('analyzer', { action: 'run' })).status).toBe(200);
			expect((await doPost('curator', { action: 'run' })).status).toBe(200);
			expect(writes().map((w) => [w.method, Object.keys(w.params).sort()])).toEqual([
				['plan.run', ['id', 'request_id']],
				['analyze.run', ['request_id']],
				['curate.run', ['request_id']]
			]);
		});

		it("stops and starts the orchestrator's hold", async () => {
			script['orchestrate.stop'] = { answer: { data: { held: true, pid: 42 }, warnings: [] } };
			script['orchestrate.start'] = { answer: { data: { held: false }, warnings: [] } };
			let r = await doPost('orchestrator', { action: 'stop' });
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject({ held: true, pid: 42 });
			r = await doPost('orchestrator', { action: 'start' });
			expect(r.status).toBe(200);
			expect(writes().map((w) => w.method)).toEqual(['orchestrate.stop', 'orchestrate.start']);
		});

		it('is 400 for Stop or Start of a role with no hold, asking flai to write nothing', async () => {
			for (const action of ['stop', 'start']) {
				const r = await doPost('reviewer', { action });
				expect(r.status).toBe(400);
				expect((await r.json()).error).toBe(
					`only the orchestrator has a hold to ${action}; the reviewer has none`
				);
			}
			expect(writes()).toEqual([]);
		});

		it("says why when the role's host action is off", async () => {
			script['role.run'] = { error: off };
			const r = await doPost('reviewer', { action: 'run', item: 'S-0003' });
			expect(r.status).toBe(403);
			expect((await r.json()).error).toContain('flai serve enable review');
		});

		it('is 404 for a role flai does not report', async () => {
			const r = await doPost('auditor', { action: 'run' });
			expect(r.status).toBe(404);
			expect(writes()).toEqual([]);
		});

		it('refuses an action that is none of run, stop, and start, or an item that is not text', async () => {
			for (const body of [{ action: 'restart' }, {}, { action: 7 }, 'run', null, 'not json']) {
				const r = await doPost('reviewer', body);
				expect(r.status).toBe(400);
				expect((await r.json()).error).toContain('action must be one of run, stop, start');
			}
			const r = await doPost('reviewer', { action: 'run', item: 3 });
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('item must be');
			expect(writes()).toEqual([]);
		});
	});
});
