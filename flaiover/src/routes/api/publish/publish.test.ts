import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, RepoError, type Ask } from '$lib/server/repo';

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new Error(`unscripted method ${method}`);
	if (s.error) throw s.error;
	// project.info is asked directly (repo().ask, raw T); publish.preview and publish.run go
	// through repo().run/write, which wrap the channel's answer as { data, warnings }.
	return method === 'project.info' ? s.data : { data: s.data, warnings: [] };
}) as Ask;

describe('/api/publish (S-0087)', () => {
	let GET: Handler;
	let POST: Handler;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		GET = mod.GET as unknown as Handler;
		POST = mod.POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('answers everything pending and whether the push action is enabled', async () => {
		script['publish.preview'] = {
			data: {
				plans: [
					{
						component: { name: 'cli', path: 'cli', kind: 'go' },
						level: 'minor',
						from: '1.0.0',
						to: '1.1.0',
						tag: 'cli/v1.1.0',
						items: [{ id: 'S-0101', title: 'Fix one', level: 'patch' }],
						files: ['cli/a.go']
					}
				]
			}
		};
		script['project.info'] = { data: { host_actions: { push: true } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.plans).toHaveLength(1);
		expect(body.plans[0]).toMatchObject({ component: { name: 'cli' }, to: '1.1.0' });
		expect(body.push_enabled).toBe(true);
		expect(asked[0]).toEqual({ method: 'publish.preview', params: {} });
	});

	it('answers an empty plan as an empty list, not an error', async () => {
		script['publish.preview'] = { data: { plans: [] } };
		script['project.info'] = { data: { host_actions: { push: false } } };
		const body = await (await GET({} as never)).json();
		expect(body).toMatchObject({ plans: [], push_enabled: false });
	});

	it('treats no flai connected as nothing pending, not an error', async () => {
		script['publish.preview'] = { error: new RepoError(503, 'no host flai is connected') };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ plans: [], push_enabled: false });
	});

	it('runs the publish and returns what happened', async () => {
		script['publish.run'] = {
			data: { plans: [{ component: { name: 'cli' } }], tags: ['cli/v1.1.0'], pushed: true }
		};
		const r = await POST({} as never);
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.pushed).toBe(true);
		expect(body.tags).toEqual(['cli/v1.1.0']);
		expect(asked[0].method).toBe('publish.run');
		expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it('answers with the host action’s refusal when it is not enabled', async () => {
		script['publish.run'] = { error: new RepoError(403, 'the host action "push" is not enabled') };
		const r = await POST({} as never);
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('not enabled');
	});
});
