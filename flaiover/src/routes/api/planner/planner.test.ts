// The planner page's read over the channel (S-0259), scripted as the plan route's test is: what
// flai answers is flai's, tested there; here it is what the route asks for and how it answers.

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

const planRun = (item: string, started: string) => ({
	item,
	agent: 'planner',
	harness: 'claude-code',
	command: 'claude',
	pid: 42,
	log: `/home/op/.flai/serve/plan-${item}.log`,
	started
});

const document = {
	kind: 'planner',
	accrued_cost: 3.1,
	accrued_seconds: 2400,
	tasks_completed: 3,
	last_run: '2026-10-03T10:33:19Z',
	path: 'wip/agents/planner.md',
	entries: [
		{
			at: '2026-10-01T09:00:00Z',
			summary: 'planned E-0016',
			items: null,
			seconds: 300,
			cost: 0.4,
			estimated: false
		},
		{
			at: '2026-10-03T10:33:19Z',
			summary: 'planned S-0252',
			trigger: 'asked',
			items: ['S-0252', 'T-0831'],
			seconds: 1999,
			cost: 2.436,
			estimated: true
		}
	]
};

describe('/api/planner', () => {
	let get: Handler;
	const doGet = () => get({} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		get = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { answer: { host_actions: { plan: true } } },
			'activity.document': { answer: document },
			'agent.status': {
				answer: {
					enabled: false,
					state: {
						command: '',
						plans: {
							'E-0016': planRun('E-0016', '2026-10-01T08:55:00Z'),
							'S-0252': planRun('S-0252', '2026-10-03T10:00:00Z')
						}
					}
				}
			}
		};
	});
	afterAll(() => useRepo(null));

	it("answers whether the plan action is on, the planner's activity, and its runs newest first", async () => {
		const r = await doGet();
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.plan_enabled).toBe(true);
		expect(body.activity).toMatchObject({
			kind: 'planner',
			accrued_cost: 3.1,
			path: document.path
		});
		expect(body.activity.entries.map((e: { items: string[] }) => e.items)).toEqual([
			[],
			['S-0252', 'T-0831']
		]);
		expect(body.runs.map((x: { item: string }) => x.item)).toEqual(['S-0252', 'E-0016']);
		expect(asked.find((a) => a.method === 'activity.document')?.params).toEqual({
			kind: 'planner'
		});
	});

	it('reads as the action off with no runs when flai cannot say, or has no state', async () => {
		script['project.info'] = { answer: { host_actions: { agent: true } } };
		script['agent.status'] = { answer: { enabled: false } };
		let body = await (await doGet()).json();
		expect(body).toMatchObject({ plan_enabled: false, runs: [] });

		script['project.info'] = { error: new AgentError(503, 'no flai connected') };
		script['agent.status'] = { error: new AgentError(503, 'no flai connected') };
		script['activity.document'] = { answer: { ...document, entries: null } };
		const r = await doGet();
		expect(r.status).toBe(200);
		body = await r.json();
		expect(body).toMatchObject({ plan_enabled: false, runs: [] });
		expect(body.activity.entries).toEqual([]);
	});

	it("is flai's error when the activity document cannot be read", async () => {
		script['activity.document'] = { error: new AgentError(503, 'no flai connected') };
		const r = await doGet();
		expect(r.status).toBe(503);
		expect((await r.json()).error).toContain('no flai connected');
	});
});
