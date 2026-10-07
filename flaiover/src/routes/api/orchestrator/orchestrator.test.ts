// The orchestrator page's read and its stop and start over the channel (S-0228), scripted as the
// planner route's test is: what flai answers is flai's, tested there; here it is what the route asks
// for and how it answers.

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

const run = {
	story: '',
	agent: 'orchestrator',
	harness: 'claude-code',
	command: 'claude',
	pid: 42,
	log: '/home/op/.flai/serve/orchestrator.log',
	started: '2026-10-03T10:00:00Z'
};

const document = {
	kind: 'orchestrator',
	accrued_cost: 1.2,
	accrued_seconds: 900,
	tasks_completed: 2,
	last_run: '2026-10-03T10:20:00Z',
	path: 'wip/agents/orchestrator.md',
	entries: [
		{
			at: '2026-10-03T10:10:00Z',
			summary: 'promoted S-0252: cost of delay 60.81 USD a week, first under policy cod',
			items: ['S-0252'],
			seconds: 300,
			cost: 0.4,
			estimated: false
		},
		{
			at: '2026-10-03T10:20:00Z',
			summary: 'answered TH-0107: the design names the split',
			items: null,
			seconds: 600,
			cost: 0.8,
			estimated: true
		}
	],
	refusals: [{ at: '2026-10-03T10:15:00Z', call: 'flai push', needs: '' }]
};

const off = new AgentError(
	403,
	'the host action "orchestrate" is not enabled for this project. On the host, in the project, run: flai serve enable orchestrate',
	-32012
);

describe('/api/orchestrator', () => {
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
			'project.info': { answer: { host_actions: { orchestrate: true } } },
			'activity.document': { answer: document },
			'agent.status': {
				answer: { enabled: false, state: { command: '', orchestrator: run } }
			}
		};
	});
	afterAll(() => useRepo(null));

	it('answers whether the action is on, the hold, its decisions, the run under way, and its runs', async () => {
		const earlier = {
			...run,
			pid: 41,
			started: '2026-10-02T10:00:00Z',
			ended: '2026-10-02T11:00:00Z'
		};
		script['agent.status'] = {
			answer: { enabled: false, state: { orchestrator: run, past_orchestrators: [earlier] } }
		};
		const r = await doGet();
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ enabled: true, held: false, run: { pid: 42 } });
		expect(body.runs).toEqual([run, earlier]);
		expect(body.activity).toMatchObject({
			kind: 'orchestrator',
			accrued_cost: 1.2,
			path: document.path,
			refusals: document.refusals
		});
		expect(body.activity.entries.map((e: { summary: string }) => e.summary)).toEqual(
			document.entries.map((e) => e.summary)
		);
		expect(body.activity.entries[1].items).toEqual([]);
		expect(asked.find((a) => a.method === 'activity.document')?.params).toEqual({
			kind: 'orchestrator'
		});
	});

	it('says the orchestrator is held, with no run under way, when the operator stopped it', async () => {
		const stopped = {
			...run,
			ended: '2026-10-03T10:30:00Z',
			outcome: 'stopped',
			why: 'stopped from the dashboard',
			held: true
		};
		script['agent.status'] = {
			answer: { enabled: false, state: { command: '', orchestrator: stopped } }
		};
		const body = await (await doGet()).json();
		expect(body).toMatchObject({ enabled: true, held: true, run: null });
		expect(body.runs).toEqual([stopped]);
	});

	it('reads as off, not held, with no runs when flai cannot say, or has no run', async () => {
		script['project.info'] = { answer: { host_actions: { agent: true } } };
		script['agent.status'] = { answer: { enabled: false, state: { command: '' } } };
		let body = await (await doGet()).json();
		expect(body).toMatchObject({ enabled: false, held: false, run: null, runs: [] });

		script['project.info'] = { error: new AgentError(503, 'no flai connected') };
		script['agent.status'] = { error: new AgentError(503, 'no flai connected') };
		script['activity.document'] = { answer: { ...document, entries: null } };
		const r = await doGet();
		expect(r.status).toBe(200);
		body = await r.json();
		expect(body).toMatchObject({ enabled: false, held: false, run: null, runs: [] });
		expect(body.activity.entries).toEqual([]);
	});

	it("is flai's error when the activity document cannot be read", async () => {
		script['activity.document'] = { error: new AgentError(503, 'no flai connected') };
		const r = await doGet();
		expect(r.status).toBe(503);
		expect((await r.json()).error).toContain('no flai connected');
	});

	it('asks flai to stop or start the orchestrator and answers what flai said', async () => {
		script['orchestrate.stop'] = { answer: { data: { held: true, pid: 42 }, warnings: [] } };
		script['orchestrate.start'] = { answer: { data: { held: false }, warnings: [] } };
		let r = await doPost({ action: 'stop' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ held: true, pid: 42 });
		r = await doPost({ action: 'start' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ held: false });
		const writes = asked.filter((a) => a.method.startsWith('orchestrate.'));
		expect(writes.map((w) => w.method)).toEqual(['orchestrate.stop', 'orchestrate.start']);
		for (const w of writes) {
			expect(Object.keys(w.params)).toEqual(['request_id']);
			expect(typeof w.params.request_id).toBe('string');
		}
	});

	it('says why when the orchestrate action is off', async () => {
		script['orchestrate.stop'] = { error: off };
		const r = await doPost({ action: 'stop' });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable orchestrate');
	});

	it('refuses an action that is neither stop nor start, asking flai nothing', async () => {
		for (const body of [{ action: 'restart' }, {}, { action: 7 }, 'stop', null]) {
			const r = await doPost(body);
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('action must be one of stop, start');
		}
		const r = await post({
			request: new Request('http://x', { method: 'POST', body: 'not json' })
		} as never);
		expect(r.status).toBe(400);
		expect(asked.filter((a) => a.method.startsWith('orchestrate.'))).toEqual([]);
	});
});
