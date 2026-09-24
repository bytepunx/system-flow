import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, RepoError, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

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
		script['host.status'] = { data: { running: true, status, dir: '/home/me/.flai/host' } };
		script['project.info'] = { data: { host_actions: { host: true, dashboard: false } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ running: true, version: '1.9.0', host_enabled: true });
		expect(body.children).toHaveLength(2);
		expect(asked[0]).toEqual({ method: 'host.status', params: {} });
	});

	it('reads the host action as off when project.info does not name it', async () => {
		script['host.status'] = { data: { running: true, status } };
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
		script['host.status'] = { data: { running: false, dir: '/home/me/.flai/host' } };
		script['project.info'] = { data: { host_actions: { host: true } } };
		const r = await GET({} as never);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			running: false,
			reason: 'flai host is not running',
			host_enabled: true
		});
	});

	it('names the host that runs for another config', async () => {
		script['host.status'] = {
			data: {
				running: false,
				elsewhere: { pid: 77, version: '1.9.0', config: '/home/me/.flai/config.json' }
			}
		};
		script['project.info'] = { data: {} };
		expect(await (await GET({} as never)).json()).toMatchObject({
			running: false,
			reason: "the machine's flai host (pid 77) runs for another config, /home/me/.flai/config.json"
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
			script['host.process'] = { data: status };
			const r = await post({ action, process });
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject({ version: '1.9.0' });
			expect(asked[0].method).toBe('host.process');
			expect(asked[0].params).toMatchObject({ process, action });
			expect(String(asked[0].params.request_id)).toMatch(/^[0-9a-f-]{36}$/);
		}
	});

	it('upgrades flai with no process named', async () => {
		script['host.upgrade'] = {
			data: { upgrade: { previous: '1.9.0', installed: '1.10.0' }, restarting: true }
		};
		const r = await post({ action: 'upgrade' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ restarting: true, upgrade: { installed: '1.10.0' } });
		expect(asked[0].method).toBe('host.upgrade');
		expect(asked[0].params.process).toBeUndefined();
	});

	it('does not repeat a write that ends its own connection when the connection is lost', async () => {
		// flai serve restarting drops the channel before it answers; the serve that comes back has
		// no record of the request, so a retry would restart it again (found live, S-0107)
		const sent: string[] = [];
		const losing = (async (method: string, params: Record<string, unknown> = {}) => {
			if (method.startsWith('host.')) sent.push(`${method} ${params.process ?? ''}`.trim());
			throw new AgentError(502, 'the host flai went away before it answered');
		}) as unknown as Ask;
		useRepo(new Repo('/nowhere', losing, async () => true));
		try {
			for (const body of [
				{ action: 'restart', process: 'serve' },
				{ action: 'stop', process: 'all' },
				{ action: 'upgrade' }
			]) {
				expect((await post(body)).status).toBe(502);
			}
			expect(sent).toEqual(['host.process serve', 'host.process all', 'host.upgrade']);
			// the MCP servers going down leaves serve's connection alone: the usual one retry stands
			sent.length = 0;
			await post({ action: 'restart', process: 'mcp' });
			expect(sent).toEqual(['host.process mcp', 'host.process mcp']);
		} finally {
			useRepo(new Repo('/nowhere', channel));
		}
	});

	it('answers with the host action’s refusal when it is not enabled', async () => {
		script['host.process'] = { error: new RepoError(403, 'the host action "host" is not enabled') };
		const r = await post({ action: 'restart', process: 'mcp' });
		expect(r.status).toBe(403);
		expect((await r.json()).error).toContain('not enabled');
	});
});
