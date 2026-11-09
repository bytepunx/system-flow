import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as plan.test.ts does: the route passes agent.status on and says, from
// project.info, whether the plan host action is on, for the board's card menu (S-0263).
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
let asked: { method: string; params: Record<string, unknown> }[] = [];
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

describe('host-agent route over the channel (S-0263)', () => {
	let get: (query?: string) => Promise<Response>;
	const status = { enabled: true, state: { command: 'claude', stories: {} } };

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = (query = '') =>
			(mod.GET as unknown as Handler)({
				url: new URL(`http://dashboard/api/host-agent${query}`)
			} as never);
	});
	beforeEach(() => {
		script = {};
		asked = [];
	});
	afterAll(() => useRepo(null));

	it('passes agent.status on with whether the plan host action is on', async () => {
		script = {
			'agent.status': { answer: status },
			'project.info': { answer: { host_actions: { plan: true } } }
		};
		expect(await (await get()).json()).toMatchObject({ ...status, plan_enabled: true });

		script['project.info'] = { answer: { host_actions: { agent: true } } };
		expect((await (await get()).json()).plan_enabled).toBe(false);
	});

	it('asks agent.status for the open and live runs alone, with no finished (S-0611)', async () => {
		// the board and the activity page ask every 15 s and show the open stories' runs alone, so
		// the route leaves finished out
		script = { 'agent.status': { answer: status }, 'project.info': { answer: {} } };
		await get();
		expect(asked.filter((a) => a.method === 'agent.status').map((a) => a.params)).toEqual([{}]);
	});

	it('asks agent.status for the runs of finished work too with finished=1 (S-0611)', async () => {
		// a done or cancelled story's agent panel shows its last run, which the default leaves out
		script = { 'agent.status': { answer: status }, 'project.info': { answer: {} } };
		await get('?finished=1');
		expect(asked.filter((a) => a.method === 'agent.status').map((a) => a.params)).toEqual([
			{ finished: true }
		]);
	});

	it('reads the plan host action as off when project.info cannot be asked', async () => {
		script = {
			'agent.status': { answer: status },
			'project.info': { error: new AgentError(500, 'project.info failed') }
		};
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ ...status, plan_enabled: false });
	});

	it('knows nothing rather than fail with no flai connected', async () => {
		script = {
			'agent.status': { error: new AgentError(503, 'no flai connected') },
			'project.info': { error: new AgentError(503, 'no flai connected') }
		};
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ enabled: false });
	});
});
