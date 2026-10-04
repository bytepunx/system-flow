import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as plan.test.ts does: the route passes agent.status on and says, from
// project.info, whether the plan host action is on, for the board's card menu (S-0263).
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string) => {
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

describe('host-agent route over the channel (S-0263)', () => {
	let get: () => Promise<Response>;
	const status = { enabled: true, state: { command: 'claude', stories: {} } };

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		get = () => (mod.GET as unknown as Handler)({} as never);
	});
	beforeEach(() => {
		script = {};
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
