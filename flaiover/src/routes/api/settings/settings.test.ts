import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, RepoError, type Ask } from '$lib/server/repo';
import { setToken, tokenMatches } from '$lib/server/auth';

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new Error(`unscripted method ${method}`);
	if (s.error) throw s.error;
	return method === 'settings.get' || method === 'project.info'
		? s.data
		: { data: s.data, warnings: [] };
}) as Ask;

describe('/api/settings (S-0105)', () => {
	let GET: Handler;
	let post: (body: unknown) => Promise<Response>;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		GET = mod.GET as unknown as Handler;
		const POST = mod.POST as unknown as Handler;
		post = (body) =>
			POST({
				request: new Request('http://x/api/settings', {
					method: 'POST',
					body: JSON.stringify(body)
				}),
				url: new URL('http://x/api/settings')
			} as never);
	});
	beforeEach(() => {
		asked = [];
		script = { 'project.info': { data: { name: 'Flow', key: 'flow' } } };
		setToken('old-token');
	});
	afterAll(() => {
		useRepo(null);
		setToken(null);
	});

	it('reads the host settings from flai', async () => {
		script['settings.get'] = {
			data: { here: true, everywhere: false, enable: 'flai serve enable settings' }
		};
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ here: true, everywhere: false });
	});

	it('passes a change to flai as settings.<kind>, without the kind', async () => {
		script['settings.action'] = { data: { ok: true } };
		const r = await post({ kind: 'action', action: 'push', on: true });
		expect(r.status).toBe(200);
		const call = asked.find((a) => a.method === 'settings.action')!;
		expect(call.params).toMatchObject({ action: 'push', on: true });
		expect(call.params.kind).toBeUndefined();
	});

	it('serves and removes a project through flai serve (S-0122)', async () => {
		script['settings.serve'] = { data: { project: { key: 'shop' } } };
		script['settings.unserve'] = { data: { removed: { key: 'blog' } } };
		expect((await post({ kind: 'serve', root: '/home/me/git/shop', key: 'shop' })).status).toBe(
			200
		);
		expect((await post({ kind: 'unserve', root: '/home/me/git/blog' })).status).toBe(200);
		expect(asked.filter((a) => a.method.startsWith('settings.'))).toMatchObject([
			{ method: 'settings.serve', params: { root: '/home/me/git/shop', key: 'shop' } },
			{ method: 'settings.unserve', params: { root: '/home/me/git/blog' } }
		]);
	});

	it('refuses a kind that is not a setting, and says what the host must enable', async () => {
		expect((await post({ kind: 'config', key: 'x' })).status).toBe(400);
		expect(asked.some((a) => a.method.startsWith('settings.'))).toBe(false);
		script['settings.agent'] = {
			error: new RepoError(403, 'enable it', {
				enable: 'flai serve enable settings --all-projects',
				hostwide: true
			})
		};
		const r = await post({ kind: 'agent', name: 'builder' });
		expect(r.status).toBe(403);
		expect(await r.json()).toMatchObject({
			enable: 'flai serve enable settings --all-projects',
			hostwide: true
		});
	});

	it('takes a rotated dashboard token at once and keeps the asking session logged in', async () => {
		script['settings.dashboard_token'] = {
			data: {
				token: 'new-token',
				login_url: 'http://x/login#token=new-token',
				rotated: true,
				restarted: false
			}
		};
		const r = await post({ kind: 'dashboard_token' });
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.login_url).toBe('http://x/login#token=new-token');
		expect(body.token).toBeUndefined();
		expect(r.headers.get('set-cookie')).toContain('flaiover_session=new-token');
		expect(tokenMatches('new-token')).toBe(true);
		expect(tokenMatches('old-token')).toBe(false);
	});
});
