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

const status = {
	pid: 4100,
	version: '1.9.0',
	started: '2026-09-24T01:00:00Z',
	addr: '127.0.0.1:4241',
	children: [
		{ name: 'serve', state: 'running', pid: 4101, version: '1.9.0', restarts: 0 },
		{ name: 'mcp', root: '/src/a', state: 'running', pid: 4102, version: '1.9.0', restarts: 1 }
	]
};

describe('/api/host (S-0107)', () => {
	let GET: Handler;
	let post: (body: unknown) => Promise<Response>;

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		const mod = await import('./+server');
		GET = mod.GET as unknown as Handler;
		const POST = mod.POST as unknown as Handler;
		post = (body) =>
			POST({
				request: new Request('http://x/api/host', { method: 'POST', body: JSON.stringify(body) })
			} as never);
	});
	beforeEach(() => {
		asked = [];
		script = {};
	});
	afterAll(() => useRepo(null));

	it('answers the host, its processes, and whether the host action is enabled', async () => {
		script['host.status'] = { data: status };
		script['project.info'] = { data: { host_actions: { host: true, dashboard: false } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ running: true, version: '1.9.0', host_enabled: true });
		expect(body.children).toHaveLength(2);
		expect(asked[0]).toEqual({ method: 'host.status', params: {} });
	});

	it('reads the host action as off when project.info does not name it', async () => {
		script['host.status'] = { data: status };
		script['project.info'] = { data: { host_actions: { dashboard: true } } };
		expect(await (await GET({} as never)).json()).toMatchObject({ host_enabled: false });
	});

	it('treats no flai connected as no host, not an error', async () => {
		script['host.status'] = { error: new RepoError(503, 'no host flai is connected') };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			running: false,
			reason: 'no host flai is connected',
			host_enabled: false
		});
		expect(asked.map((a) => a.method)).toEqual(['host.status']);
	});

	it('treats a flai serve with no host above it as no host, with the reason', async () => {
		script['host.status'] = {
			error: new RepoError(502, 'flai host is not running; start it with flai host start')
		};
		script['project.info'] = { data: { host_actions: { host: true } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			running: false,
			reason: 'flai host is not running; start it with flai host start',
			host_enabled: true
		});
	});

	it('passes any other failure on', async () => {
		script['host.status'] = {
			error: new RepoError(500, 'flai host status: the state is unreadable')
		};
		expect((await GET({} as never)).status).toBe(500);
	});

	it('refuses an action it does not know, asking nothing', async () => {
		const r = await post({ action: 'pull' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('action must be one of');
		expect(asked).toHaveLength(0);
	});

	it('refuses start, stop, or restart without a process it knows', async () => {
		for (const body of [{ action: 'restart' }, { action: 'stop', process: 'flaiover' }]) {
			const r = await post(body);
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('process must be one of');
		}
		expect(asked).toHaveLength(0);
	});

	it('checks for a newer flai as a read, changing nothing', async () => {
		script['host.check'] = {
			data: { current: '1.9.0', latest: '1.10.0', upgrade_available: true }
		};
		const r = await post({ action: 'check' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ upgrade_available: true, latest: '1.10.0' });
		expect(asked[0]).toEqual({ method: 'host.check', params: {} });
	});

	it('starts, stops, and restarts the process named, sending a request_id', async () => {
		for (const [action, process] of [
			['start', 'mcp'],
			['stop', 'serve'],
			['restart', 'all']
		]) {
			asked = [];
			script[`host.${action}`] = { data: status };
			const r = await post({ action, process });
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject({ version: '1.9.0' });
			expect(asked[0].method).toBe(`host.${action}`);
			expect(asked[0].params.process).toBe(process);
			expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
		}
	});

	it('upgrades flai with no process named', async () => {
		script['host.upgrade'] = { data: { outcome: 'upgraded', from: '1.9.0', to: '1.10.0' } };
		const r = await post({ action: 'upgrade' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ outcome: 'upgraded', to: '1.10.0' });
		expect(asked[0].method).toBe('host.upgrade');
		expect(asked[0].params.process).toBeUndefined();
	});

	it('answers with the host action’s refusal when it is not enabled', async () => {
		script['host.restart'] = { error: new RepoError(403, 'the host action "host" is not enabled') };
		const r = await post({ action: 'restart', process: 'mcp' });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('not enabled');
	});
});
