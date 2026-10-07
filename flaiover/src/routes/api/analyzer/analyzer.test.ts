// The analyzer page's read and its Run over the channel (S-0228), scripted as the planner route's
// test is: what flai answers is flai's, tested there; here it is what the route asks for and how it
// answers.

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

const ended = {
	story: '',
	agent: 'analyzer',
	harness: 'claude-code',
	command: 'claude',
	pid: 42,
	log: '/home/op/.flai/serve/analyzer.log',
	started: '2026-10-03T10:00:00Z',
	ended: '2026-10-03T10:20:00Z',
	outcome: 'worked',
	focus: 'risk',
	trigger: 'asked',
	report: 'design/analysis/2026-10-03-risk.md'
};

const document = {
	kind: 'analyzer',
	accrued_cost: 0.9,
	accrued_seconds: 1200,
	tasks_completed: 1,
	last_run: '2026-10-03T10:20:00Z',
	path: 'wip/agents/analyzer.md',
	entries: [
		{
			at: '2026-10-03T10:20:00Z',
			summary: 'analyzed risk: design/analysis/2026-10-03-risk.md',
			items: null,
			seconds: 1200,
			cost: 0.9,
			estimated: false
		}
	]
};

const started = {
	focus: 'intent',
	agent: 'analyzer',
	harness: 'claude-code',
	command: 'claude',
	pid: 43,
	log: '/home/op/.flai/serve/analyzer.log',
	session: 'sess-1',
	started: '2026-10-03T11:00:00Z',
	trigger: 'asked'
};

describe('/api/analyzer', () => {
	let get: Handler;
	let post: Handler;
	const doGet = () => get({} as never);
	const doPost = (body: unknown) =>
		post({
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = mod.GET as unknown as Handler;
		post = mod.POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { answer: { host_actions: { analyze: true } } },
			'activity.document': { answer: document },
			'agent.status': { answer: { enabled: false, state: { command: '', analyzer: ended } } }
		};
	});
	afterAll(() => useRepo(null));

	it("answers whether the analyze action is on, the analyzer's activity, and its runs", async () => {
		const earlier = {
			...ended,
			pid: 41,
			started: '2026-10-02T10:00:00Z',
			ended: '2026-10-02T10:30:00Z'
		};
		script['agent.status'] = {
			answer: { enabled: false, state: { analyzer: ended, past_analyzers: [earlier] } }
		};
		const r = await doGet();
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ enabled: true, run: null, runs: [ended, earlier] });
		expect(body.activity).toMatchObject({ kind: 'analyzer', accrued_cost: 0.9 });
		expect(body.activity.entries[0].items).toEqual([]);
		expect(asked.find((a) => a.method === 'activity.document')?.params).toEqual({
			kind: 'analyzer'
		});
	});

	it('answers the run under way as the current run', async () => {
		const running = { ...ended, ended: undefined, outcome: undefined, report: undefined };
		script['agent.status'] = { answer: { enabled: false, state: { analyzer: running } } };
		const body = await (await doGet()).json();
		expect(body.run).toMatchObject({ pid: 42, focus: 'risk' });
		expect(body.run.ended).toBeUndefined();
	});

	it('reads as off with no runs when flai cannot say, or has no run', async () => {
		script['project.info'] = { answer: { host_actions: { orchestrate: true } } };
		script['agent.status'] = { answer: { enabled: false } };
		let body = await (await doGet()).json();
		expect(body).toMatchObject({ enabled: false, run: null, runs: [] });

		script['project.info'] = { error: new AgentError(503, 'no flai connected') };
		script['agent.status'] = { error: new AgentError(503, 'no flai connected') };
		const r = await doGet();
		expect(r.status).toBe(200);
		body = await r.json();
		expect(body).toMatchObject({ enabled: false, run: null, runs: [] });
	});

	it("is flai's error when the activity document cannot be read", async () => {
		script['activity.document'] = { error: new AgentError(503, 'no flai connected') };
		const r = await doGet();
		expect(r.status).toBe(503);
		expect((await r.json()).error).toContain('no flai connected');
	});

	it('asks flai to run the analyzer with the focus, or with none, and answers the run', async () => {
		script['analyze.run'] = { answer: { data: started, warnings: [] } };
		let r = await doPost({ focus: 'intent' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ focus: 'intent', pid: 43 });
		for (const body of [{}, { focus: '' }, { focus: null }, null]) {
			r = await doPost(body);
			expect(r.status).toBe(200);
		}
		const writes = asked.filter((a) => a.method === 'analyze.run');
		expect(writes.map(({ params: { request_id, ...rest } }) => [typeof request_id, rest])).toEqual([
			['string', { focus: 'intent' }],
			['string', {}],
			['string', {}],
			['string', {}],
			['string', {}]
		]);
	});

	it('says why flai refused, and what enables the action when it is off', async () => {
		script['analyze.run'] = {
			error: new AgentError(
				400,
				'the analyzer is already running for this project (pid 42, started 2026-10-03T10:00:00Z)',
				-32011
			)
		};
		let r = await doPost({ focus: 'risk' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('already running');
		script['analyze.run'] = {
			error: new AgentError(
				403,
				'the host action "analyze" is not enabled for this project. On the host, in the project, run: flai serve enable analyze',
				-32012
			)
		};
		r = await doPost({});
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable analyze');
	});

	it('refuses a focus that is not one, asking flai nothing', async () => {
		for (const focus of ['all', 'speed', 3]) {
			const r = await doPost({ focus });
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('focus must be one of bottlenecks, intent, risk');
		}
		expect(asked.filter((a) => a.method === 'analyze.run')).toEqual([]);
	});
});
