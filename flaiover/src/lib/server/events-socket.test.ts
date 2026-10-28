// ADR-0182: the live events reach the browser over a WebSocket, an upgrade of /api/events routed ahead
// of the agent registry, authenticated as a GET of /api/events is, before the socket opens. Tested on
// a real HTTP server that hands every upgrade to the router, the way server.js and the dev server do.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import http from 'node:http';
import { EventEmitter } from 'node:events';
import { mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import type { AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import WebSocket from 'ws';
import { exposeAgentUpgrade, registry, resetAgent } from './agent';
import { setToken } from './auth';
import { CLOSE_GOING_AWAY, EVENTS_PATH, EventsSocket } from './events-socket';
import { knownRepos, useRepo, type Repo } from './repo';

const TOKEN = 'dashboard-token-for-tests';
const BEARER = { authorization: `Bearer ${TOKEN}` };

/** A Repo that announces what a test emits and answers nothing. */
class FakeRepo extends EventEmitter {
	watched = 0;
	async watch(): Promise<void> {
		this.watched++;
	}
}

/** A client's socket, with every message it was sent, parsed, in order. */
type Client = { ws: WebSocket; next: () => Promise<unknown>; closed: Promise<number> };

/** Open a socket, or learn the status the upgrade was refused with. */
function dial(
	url: string,
	headers: Record<string, string> = {},
	opt: { autoPong?: boolean } = {}
): Promise<Client | { status: number }> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(url, { headers, ...opt });
		const got: unknown[] = [];
		const waiting: ((m: unknown) => void)[] = [];
		ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			const w = waiting.shift();
			if (w) w(m);
			else got.push(m);
		});
		const closed = new Promise<number>((r) => ws.on('close', (code) => r(code)));
		const next = () =>
			got.length ? Promise.resolve(got.shift()) : new Promise((r) => waiting.push(r));
		ws.on('open', () => resolve({ ws, next, closed }));
		ws.on('unexpected-response', (_req, res) => resolve({ status: res.statusCode ?? 0 }));
		ws.on('error', reject);
	});
}

async function opened(url: string, headers: Record<string, string> = BEARER): Promise<Client> {
	const c = await dial(url, headers);
	if (!('ws' in c)) throw new Error(`refused with ${c.status}`);
	return c;
}

describe('the live events socket (ADR-0182)', () => {
	let dir: string;
	let server: http.Server;
	let base: string;
	let events: EventsSocket;
	let harbour: FakeRepo;
	const clients: WebSocket[] = [];

	/** The router on a real HTTP server, as server.js wires it, with a ping as often as pingMs. */
	async function serve(pingMs = 60000) {
		events = new EventsSocket(pingMs);
		exposeAgentUpgrade({ [EVENTS_PATH]: events });
		server = http.createServer((_req, res) => res.end('ok'));
		server.on('upgrade', (req, socket, head) =>
			globalThis.__flaioverAgentUpgrade?.(req, socket, head)
		);
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		base = `127.0.0.1:${(server.address() as AddressInfo).port}`;
	}

	const url = (path = `${EVENTS_PATH}?project=harbour`) => `ws://${base}${path}`;

	async function client(path?: string, headers?: Record<string, string>) {
		const c = await opened(url(path), headers);
		clients.push(c.ws);
		return c;
	}

	beforeEach(() => {
		resetAgent();
		dir = mkdtempSync(join(tmpdir(), 'flaiover-events-socket-'));
		writeFileSync(join(dir, 'agent-key'), 'agent-credential-for-tests');
		registry({ FLAIOVER_AGENT_KEY_FILE: join(dir, 'agent-key') }, { allowUnsigned: true });
		setToken(TOKEN);
		harbour = new FakeRepo();
		useRepo(harbour as unknown as Repo, 'harbour');
	});

	afterEach(async () => {
		for (const ws of clients.splice(0)) ws.terminate();
		events.close();
		globalThis.__flaioverAgentUpgrade = undefined;
		await new Promise((r) => server.close(r));
		for (const key of [...knownRepos().keys()]) useRepo(null, key);
		setToken(null);
		resetAgent();
		vi.unstubAllEnvs();
		rmSync(dir, { recursive: true, force: true });
	});

	it('refuses an upgrade with no token 401 before the socket opens', async () => {
		await serve();
		expect(await dial(url())).toEqual({ status: 401 });
		expect(await dial(url(), { cookie: 'flaiover_session=wrong' })).toEqual({ status: 401 });
		expect(harbour.watched).toBe(0);
		expect(harbour.listenerCount('change')).toBe(0);
	});

	it('refuses an upgrade from another site 403, cookie and all, before the socket opens', async () => {
		await serve();
		const cookie = { cookie: `flaiover_session=${TOKEN}` };
		expect(await dial(url(), { ...cookie, origin: 'http://elsewhere.example' })).toEqual({
			status: 403
		});
		expect(await dial(url(), { ...cookie, origin: 'not a url' })).toEqual({ status: 403 });
		expect(harbour.listenerCount('change')).toBe(0);
		// the dashboard's own page, which a browser names in Origin, is let in
		const own = await client(undefined, { ...cookie, origin: `http://${base}` });
		expect(await own.next()).toEqual({ type: 'ready' });
	});

	it("lets in an origin FLAIOVER_ALLOWED_HOSTS names, with a Host of the container's address (ADR-0193)", async () => {
		vi.stubEnv('FLAIOVER_ALLOWED_HOSTS', 'flow.example.com');
		await serve();
		const cookie = { cookie: `flaiover_session=${TOKEN}` };
		const tunnel = await client(undefined, { ...cookie, origin: 'https://flow.example.com' });
		expect(await tunnel.next()).toEqual({ type: 'ready' });
		expect(await dial(url(), { ...cookie, origin: 'http://elsewhere.example' })).toEqual({
			status: 403
		});
		expect(await dial(url(), { ...cookie, origin: 'https://flow.example.com:8443' })).toEqual({
			status: 403
		});
	});

	it('refuses the origin of a tunnel FLAIOVER_ALLOWED_HOSTS does not name 403', async () => {
		await serve();
		const cookie = { cookie: `flaiover_session=${TOKEN}` };
		expect(await dial(url(), { ...cookie, origin: 'https://flow.example.com' })).toEqual({
			status: 403
		});
	});

	it("lets in the dashboard's own page whatever the case of the Host it was sent with", async () => {
		await serve();
		const port = base.split(':')[1];
		const own = await client(undefined, {
			cookie: `flaiover_session=${TOKEN}`,
			host: `LocalHost:${port}`,
			origin: `http://localhost:${port}`
		});
		expect(await own.next()).toEqual({ type: 'ready' });
	});

	it("sends ready, then the project's changes and agents, each a JSON message with a type", async () => {
		await serve();
		const c = await client();
		expect(await c.next()).toEqual({ type: 'ready' });
		expect(harbour.watched).toBe(1);
		harbour.emit('change', 'wip/kanban/stories/S-0454.md', 'item');
		expect(await c.next()).toEqual({
			type: 'change',
			path: 'wip/kanban/stories/S-0454.md',
			kind: 'item'
		});
		harbour.emit('change', 'system-flow.yaml');
		expect(await c.next()).toEqual({ type: 'change', path: 'system-flow.yaml' });
		harbour.emit('agent', 'S-0454');
		expect(await c.next()).toEqual({ type: 'agent', story: 'S-0454' });
	});

	it("follows the project the URL names, and only that project's Repo", async () => {
		await serve();
		const other = new FakeRepo();
		useRepo(other as unknown as Repo, 'lighthouse');
		const c = await client(`${EVENTS_PATH}?project=lighthouse`);
		expect(await c.next()).toEqual({ type: 'ready' });
		expect(other.listenerCount('change')).toBe(1);
		expect(harbour.listenerCount('change')).toBe(0);
	});

	it('stops listening to the Repo when the socket closes', async () => {
		await serve();
		const c = await client();
		await c.next();
		expect(harbour.listenerCount('change')).toBe(1);
		expect(harbour.listenerCount('agent')).toBe(1);
		c.ws.close();
		await vi.waitFor(() => {
			expect(harbour.listenerCount('change')).toBe(0);
			expect(harbour.listenerCount('agent')).toBe(0);
		});
	});

	it('pings, and closes a socket that missed the pong of the ping before', async () => {
		await serve(40);
		const answering = await client();
		let pinged = 0;
		answering.ws.on('ping', () => pinged++);
		const silent = await dial(url(), BEARER, { autoPong: false });
		if (!('ws' in silent)) throw new Error('refused');
		clients.push(silent.ws);
		expect(await silent.closed).toBe(1006);
		await vi.waitFor(() => expect(pinged).toBeGreaterThan(2));
		expect(answering.ws.readyState).toBe(WebSocket.OPEN);
	});

	it('closes every events socket with 1001 when the server stops, so the server can close', async () => {
		await serve();
		const c = await client();
		await c.next();
		expect(harbour.listenerCount('change')).toBe(1);
		const stopped = new Promise((r) => server.close(r));
		globalThis.__flaioverAgentUpgrade?.close();
		// torn down at once, before the close handshake, which an upgraded socket's server does not wait for (I-0181, I-0205)
		expect(harbour.listenerCount('change')).toBe(0);
		expect(harbour.listenerCount('agent')).toBe(0);
		expect(await c.closed).toBe(CLOSE_GOING_AWAY);
		await stopped;
		expect(harbour.listenerCount('change')).toBe(0);
	});

	it('routes /agent to the registry as before, which still refuses a browser and other paths', async () => {
		await serve();
		const flai = await opened(url('/agent'), {});
		clients.push(flai.ws);
		expect(await dial(url('/agent'), { origin: `http://${base}` })).toEqual({ status: 403 });
		expect(await dial(url('/elsewhere'), BEARER)).toEqual({ status: 404 });
	});
});
