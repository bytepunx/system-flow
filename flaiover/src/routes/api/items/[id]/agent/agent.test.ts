import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted as checks.test.ts does (S-0116): whether a restart is allowed is flai's to
// judge, tested there; here it is what the route asks for and how it passes flai's answer on.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { answer?: unknown; error?: AgentError }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method] ?? {};
	if (s.error) throw s.error;
	return s.answer;
}) as Ask;

describe('agent route over the channel (S-0116)', () => {
	let post: Handler;
	const doPost = (id: string, body: unknown) =>
		post({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		post = (await import('./+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('asks flai to restart the story agent and passes on the run', async () => {
		script = {
			'agent.restart': {
				answer: { data: { story: 'S-0116', agent: 'agent-S-0116', pid: 42 }, warnings: [] }
			}
		};
		const r = await doPost('S-0116', { action: 'restart' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ story: 'S-0116', agent: 'agent-S-0116', pid: 42 });
		const write = asked.find((a) => a.method === 'agent.restart');
		expect(write?.params).toMatchObject({ id: 'S-0116' });
		expect(typeof write?.params.request_id).toBe('string');
	});

	it('says why flai refused, and what enables the action when it is off', async () => {
		script = {
			'agent.restart': {
				error: new AgentError(400, "S-0116's agent is running (pid 42)", -32011)
			}
		};
		let r = await doPost('S-0116', { action: 'restart' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('is running');
		script = {
			'agent.restart': {
				error: new AgentError(
					403,
					'the host action "agent" is not enabled for this project. On the host, in the project, run: flai serve enable agent',
					-32012
				)
			}
		};
		r = await doPost('S-0116', { action: 'restart' });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('flai serve enable agent');
	});

	it('refuses an action it does not know without asking flai', async () => {
		const r = await doPost('S-0116', { action: 'stop' });
		expect(r.status).toBe(400);
		expect(asked.some((a) => a.method.startsWith('agent.'))).toBe(false);
	});
});
