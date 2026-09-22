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
	return method === 'project.info' ? s.data : { data: s.data, warnings: [] };
}) as Ask;

describe('/api/dashboard (S-0081)', () => {
	let GET: Handler;
	let POST: Handler;
	let post: (action: unknown) => Promise<Response>;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		GET = mod.GET as unknown as Handler;
		POST = mod.POST as unknown as Handler;
		post = (action) =>
			POST({
				request: new Request('http://x/api/dashboard', {
					method: 'POST',
					body: JSON.stringify({ action })
				})
			} as never);
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('answers what is running and whether the dashboard action is enabled', async () => {
		script['dashboard.status'] = {
			data: { container: 'flaiover', running: true, image: 'ghcr.io/bytepunx/flaiover:latest' }
		};
		script['project.info'] = { data: { host_actions: { dashboard: true } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ container: 'flaiover', running: true, dashboard_enabled: true });
		expect(asked[0]).toEqual({ method: 'dashboard.status', params: {} });
	});

	it('treats no flai connected as not running, not an error', async () => {
		script['dashboard.status'] = { error: new RepoError(503, 'no host flai is connected') };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ running: false, dashboard_enabled: false });
	});

	it('refuses a POST with no recognised action', async () => {
		const r = await post('pull');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('action must be one of');
		expect(asked).toHaveLength(0);
	});

	it('checks for an update as a read, changing nothing', async () => {
		script['dashboard.check'] = {
			data: { container: 'flaiover', running: true, upgrade_available: true }
		};
		const r = await post('check');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ upgrade_available: true });
		expect(asked[0]).toEqual({ method: 'dashboard.check', params: {} });
	});

	it('restarts, sending a request_id', async () => {
		script['dashboard.restart'] = {
			data: { container: 'flaiover', was_running: true, ref: 'ghcr.io/bytepunx/flaiover:latest' }
		};
		const r = await post('restart');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ was_running: true });
		expect(asked[0].method).toBe('dashboard.restart');
		expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
	});

	it('upgrades and returns what happened', async () => {
		script['dashboard.upgrade'] = {
			data: { container: 'flaiover', outcome: 'upgraded', from: '0.22.6', to: '0.22.7' }
		};
		const r = await post('upgrade');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ outcome: 'upgraded', to: '0.22.7' });
		expect(asked[0].method).toBe('dashboard.upgrade');
	});

	it('answers with the host action’s refusal when it is not enabled', async () => {
		script['dashboard.stop'] = {
			error: new RepoError(403, 'the host action "dashboard" is not enabled')
		};
		const r = await post('stop');
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('not enabled');
	});

	it('answers an upgrade that never became healthy as an ordinary failure, not a partial one', async () => {
		script['dashboard.upgrade'] = {
			error: new RepoError(
				500,
				'upgrade to ghcr.io/bytepunx/flaiover:latest failed: the new image did not answer healthy in time; flaiover keeps running unchanged'
			)
		};
		const r = await post('upgrade');
		expect(r.status).toBe(500);
		expect((await r.json()).error).toContain('keeps running unchanged');
	});
});
