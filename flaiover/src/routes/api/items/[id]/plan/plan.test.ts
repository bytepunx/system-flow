import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as agent.test.ts does (S-0208): whether the planner may start is flai's to
// judge, tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

const run = {
	item: 'E-0016',
	agent: 'planner',
	harness: 'claude-code',
	command: 'claude',
	pid: 42,
	log: '/home/op/.flai/serve/plan-E-0016.log',
	session: 'sess-1',
	started: '2026-10-03T10:00:00Z'
};

describe('plan route over the channel (S-0208)', () => {
	let get: Handler;
	let post: Handler;
	const doGet = (id: string) => get({ params: { id } } as never);
	const doPost = (id: string) =>
		post({ params: { id }, request: new Request('http://x', { method: 'POST' }) } as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = mod.GET as unknown as Handler;
		post = mod.POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('says whether the plan host action is on, with the item’s newest planner run', async () => {
		script = {
			'project.info': { answer: { host_actions: { plan: true, agent: false } } },
			'agent.status': {
				answer: { enabled: false, state: { command: '', plans: { 'E-0016': run } } }
			}
		};
		let body = await (await doGet('E-0016')).json();
		expect(body).toMatchObject({ plan_enabled: true, run: { item: 'E-0016', pid: 42 } });
		body = await (await doGet('S-0208')).json();
		expect(body).toMatchObject({ plan_enabled: true, run: null });

		script['project.info'] = { answer: { host_actions: { agent: true } } };
		body = await (await doGet('E-0016')).json();
		expect(body.plan_enabled).toBe(false);
	});

	it('reads as off with no run rather than fail when flai cannot be asked', async () => {
		script = {
			'project.info': { error: new AgentError(503, 'no flai connected') },
			'agent.status': { error: new AgentError(503, 'no flai connected') }
		};
		const r = await doGet('E-0016');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ plan_enabled: false, run: null });
	});

	it('asks flai to start the planner for the item and passes on the run', async () => {
		script = { 'plan.run': { answer: { data: run, warnings: [] } } };
		const r = await doPost('E-0016');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ item: 'E-0016', pid: 42, log: run.log });
		const write = asked.find((a) => a.method === 'plan.run');
		expect(write?.params).toMatchObject({ id: 'E-0016' });
		expect(typeof write?.params.request_id).toBe('string');
	});

	it('says why flai refused, and what enables the action when it is off', async () => {
		script = {
			'plan.run': {
				error: new AgentError(
					400,
					'the planner is already running for E-0016 (pid 42, started 2026-10-03T10:00:00Z); one item has one planner at a time',
					-32011
				)
			}
		};
		let r = await doPost('E-0016');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('already running for E-0016');
		script = {
			'plan.run': {
				error: new AgentError(
					403,
					'the host action "plan" is not enabled for this project. On the host, in the project, run: flai serve enable plan',
					-32012
				)
			}
		};
		r = await doPost('E-0016');
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable plan');
	});
});
