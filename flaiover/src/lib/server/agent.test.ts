import { afterEach, describe, expect, it } from 'vitest';
import http from 'node:http';
import type { AddressInfo } from 'node:net';
import WebSocket from 'ws';
import { AgentError, AgentRegistry, proof, REQUIRED_METHODS } from './agent';

const KEY = 'agent-credential-for-tests';

type Flai = { ws: WebSocket; requests: Record<string, unknown>[]; closed: Promise<number> };

/** A registry on a real HTTP server, the way server.js wires it. */
async function serve(registry: AgentRegistry) {
	const server = http.createServer((_req, res) => res.end('ok'));
	server.on('upgrade', (req, socket, head) => registry.handleUpgrade(req, socket, head));
	await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
	const port = (server.address() as AddressInfo).port;
	return { server, url: `ws://127.0.0.1:${port}/agent`, port };
}

/** The flai end of the handshake, with the key it believes in, for the named project. */
function connect(
	url: string,
	key: string,
	answer: (m: { id: number; method: string; params: Record<string, unknown> }) => unknown = () => ({
		name: 'Harbour'
	}),
	headers: Record<string, string> = {},
	methods: string[] = REQUIRED_METHODS,
	project: { key: string; name: string } = { key: 'harbour', name: 'Harbour' }
): Promise<Flai> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(url, { headers });
		const flai: Flai = {
			ws,
			requests: [],
			closed: new Promise((r) => ws.on('close', (code) => r(code)))
		};
		const mine = 'c'.repeat(32);
		let proven = false;
		ws.on('open', () =>
			ws.send(
				JSON.stringify({
					jsonrpc: '2.0',
					id: 'hello',
					method: 'hello',
					params: {
						protocol: 1,
						nonce: mine,
						flai: '9.9.9',
						project,
						methods
					}
				})
			)
		);
		ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (!proven) {
				proven = true;
				if (m.result.proof !== proof(KEY, 'dashboard', mine, m.result.nonce))
					return reject(new Error('the dashboard did not prove itself'));
				ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						method: 'hello.prove',
						params: { proof: proof(key, 'flai', m.result.nonce, mine) }
					})
				);
				return setTimeout(() => resolve(flai), 30);
			}
			flai.requests.push(m);
			if (m.method === '$/cancel') return;
			const result = answer(m);
			if (result !== undefined) ws.send(JSON.stringify({ jsonrpc: '2.0', id: m.id, result }));
		});
		ws.on('error', reject);
		ws.on('unexpected-response', (_req, res) => reject(new Error(`status ${res.statusCode}`)));
	});
}

describe('AgentRegistry and AgentHub', () => {
	const cleanup: (() => void)[] = [];
	afterEach(() => {
		for (const f of cleanup.splice(0)) f();
	});
	async function setup(opt = {}) {
		const registry = new AgentRegistry(KEY, opt);
		const s = await serve(registry);
		cleanup.push(() => {
			registry.close();
			s.server.close();
		});
		return { registry, hub: registry.hub('harbour'), ...s };
	}

	it('says 503 until a flai has proven itself, then asks it named methods for its project', async () => {
		const { hub, url } = await setup();
		expect(hub.status()).toEqual({ configured: true, connected: false });
		await expect(hub.ask('project.info')).rejects.toMatchObject({ status: 503 });

		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(hub.status()).toMatchObject({
			connected: true,
			flai: '9.9.9',
			serves: { key: 'harbour', name: 'Harbour' }
		});
		await expect(hub.ask('project.info')).resolves.toEqual({ name: 'Harbour' });
		expect(flai.requests[0]).toMatchObject({
			method: 'project.info',
			params: { project: 'harbour' }
		});
	});

	it('serves nothing to a connection that cannot prove the credential', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, 'some-other-key');
		expect(await flai.closed).toBe(4401);
		expect(hub.status().connected).toBe(false);
	});

	it('refuses a browser: any Origin header, and any other path', async () => {
		const { url, port } = await setup();
		await expect(connect(url, KEY, undefined, { Origin: 'http://127.0.0.1' })).rejects.toThrow(
			'status 403'
		);
		await expect(connect(`ws://127.0.0.1:${port}/elsewhere`, KEY)).rejects.toThrow('status 404');
	});

	it('closes a proven connection that names no project', async () => {
		const { url } = await setup();
		const flai = await connect(url, KEY, undefined, {}, REQUIRED_METHODS, { key: '', name: '' });
		expect(await flai.closed).toBe(4400);
	});

	it('is 503 for everyone when the dashboard was given no credential', async () => {
		const registry = new AgentRegistry(null);
		const s = await serve(registry);
		cleanup.push(() => {
			registry.close();
			s.server.close();
		});
		expect(registry.hub('harbour').status()).toEqual({ configured: false, connected: false });
		expect(registry.configured()).toBe(false);
		await expect(connect(s.url, KEY)).rejects.toThrow('status 503');
	});

	it('fails the request in flight when flai goes away, and the next one at once', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => undefined); // never answers
		const asked = hub.ask('project.info');
		await new Promise((r) => setTimeout(r, 20));
		flai.ws.terminate();
		await expect(asked).rejects.toMatchObject({ status: 502 });
		await expect(hub.ask('project.info')).rejects.toBeInstanceOf(AgentError);
		expect(hub.status().connected).toBe(false);
	});

	it('gives up on a method that is not answered, and tells flai to cancel it', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => undefined);
		cleanup.push(() => flai.ws.terminate());
		await expect(hub.ask('project.info', {}, 50)).rejects.toMatchObject({ status: 504 });
		await new Promise((r) => setTimeout(r, 30));
		expect(flai.requests.at(-1)).toMatchObject({ method: '$/cancel' });
	});

	it('passes on the error flai answers with', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => undefined);
		cleanup.push(() => flai.ws.terminate());
		flai.ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (m.id)
				flai.ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						id: m.id,
						error: { code: -32601, message: 'flai offers no method shell.run' }
					})
				);
		});
		await expect(hub.ask('shell.run')).rejects.toMatchObject({
			status: 502,
			code: -32601,
			message: 'flai offers no method shell.run'
		});
	});

	it('lets a newer proven connection replace the older one, for the same project', async () => {
		const { hub, url } = await setup();
		const first = await connect(url, KEY, () => ({ from: 'first' }));
		const second = await connect(url, KEY, () => ({ from: 'second' }));
		cleanup.push(() => second.ws.terminate());
		expect(await first.closed).toBe(4000);
		await expect(hub.ask('project.info')).resolves.toEqual({ from: 'second' });
	});

	it('marks a flai that stops answering pings as gone', async () => {
		const { hub, url } = await setup({ pingMs: 40 });
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		// A frozen process: the socket stays open and nothing answers, pongs included.
		flai.ws.pong = () => {};
		(
			flai.ws as unknown as { _receiver: { removeAllListeners: (e: string) => void } }
		)._receiver.removeAllListeners('ping');
		await new Promise((r) => setTimeout(r, 200));
		expect(hub.status().connected).toBe(false);
	});

	it("passes on flai's word that a file changed, and says when a flai has connected", async () => {
		const { hub, url } = await setup();
		const events: string[] = [];
		hub.on('connected', () => events.push('connected'));
		hub.on('gone', () => events.push('gone'));
		hub.on('change', (path: string) => events.push(path));
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		flai.ws.send(
			JSON.stringify({
				jsonrpc: '2.0',
				method: 'change',
				params: { project: 'harbour', path: 'wip/kanban/stories/S-0001-a.md' }
			})
		);
		// not a path, and not a notification: neither is announced
		flai.ws.send(JSON.stringify({ jsonrpc: '2.0', method: 'change', params: { path: 7 } }));
		flai.ws.send(
			JSON.stringify({ jsonrpc: '2.0', id: 99, method: 'change', params: { path: 'x' } })
		);
		await new Promise((r) => setTimeout(r, 50));
		expect(events).toEqual(['connected', 'wip/kanban/stories/S-0001-a.md']);
		flai.ws.terminate();
		await new Promise((r) => setTimeout(r, 50));
		expect(events.at(-1)).toBe('gone');
	});

	it('says what a flai older than the dashboard does not offer', async () => {
		const { hub, url } = await setup();
		const old = await connect(url, KEY, undefined, {}, ['project.info', 'board.get']);
		cleanup.push(() => old.ws.terminate());
		const missing = hub.status().missing ?? [];
		expect(missing).toContain('accept.run');
		expect(missing).not.toContain('board.get');
		old.ws.terminate();
		await new Promise((r) => setTimeout(r, 30));
		const current = await connect(url, KEY);
		cleanup.push(() => current.ws.terminate());
		expect(hub.status().missing).toEqual([]);
	});

	it('hands a request its progress, by the request ID', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => undefined);
		cleanup.push(() => flai.ws.terminate());
		flai.ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (m.method !== 'accept.run') return;
			for (const step of ['merged', 'committed'])
				flai.ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						method: '$/progress',
						params: { id: m.id, value: step }
					})
				);
			// progress for a request that is not waiting is dropped
			flai.ws.send(
				JSON.stringify({ jsonrpc: '2.0', method: '$/progress', params: { id: 9999, value: 'x' } })
			);
			flai.ws.send(JSON.stringify({ jsonrpc: '2.0', id: m.id, result: { data: 'done' } }));
		});
		const steps: unknown[] = [];
		await expect(hub.ask('accept.run', {}, 2000, (v) => steps.push(v))).resolves.toEqual({
			data: 'done'
		});
		expect(steps).toEqual(['merged', 'committed']);
	});

	it('keeps what flai sent with an error', async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => undefined);
		cleanup.push(() => flai.ws.terminate());
		flai.ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (m.id)
				flai.ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						id: m.id,
						error: { code: -32009, message: 'the file changed', data: { hash: 'h2' } }
					})
				);
		});
		await expect(hub.ask('doc.save')).rejects.toMatchObject({ code: -32009, data: { hash: 'h2' } });
	});

	// S-0080: one dashboard, several projects, one shared credential.
	it('serves two projects over two connections at once, without either replacing the other', async () => {
		const { registry, url } = await setup();
		const first = await connect(url, KEY, () => ({ from: 'harbour' }), {}, REQUIRED_METHODS, {
			key: 'harbour',
			name: 'Harbour'
		});
		const second = await connect(url, KEY, () => ({ from: 'quay' }), {}, REQUIRED_METHODS, {
			key: 'quay',
			name: 'Quay'
		});
		cleanup.push(() => {
			first.ws.terminate();
			second.ws.terminate();
		});
		await new Promise((r) => setTimeout(r, 30));
		expect(
			await Promise.race([first.closed, new Promise((r) => setTimeout(() => r('open'), 50))])
		).toBe('open');
		expect(registry.hub('harbour').status()).toMatchObject({
			connected: true,
			serves: { key: 'harbour' }
		});
		expect(registry.hub('quay').status()).toMatchObject({
			connected: true,
			serves: { key: 'quay' }
		});
		await expect(registry.hub('harbour').ask('project.info')).resolves.toEqual({ from: 'harbour' });
		await expect(registry.hub('quay').ask('project.info')).resolves.toEqual({ from: 'quay' });

		// losing one project's connection does not touch the other's
		first.ws.terminate();
		await new Promise((r) => setTimeout(r, 30));
		expect(registry.hub('harbour').status().connected).toBe(false);
		expect(registry.hub('quay').status().connected).toBe(true);

		expect(registry.list().map((p) => p.key)).toEqual(['harbour', 'quay']);
	});

	it('tells a project nothing has ever connected for apart from one that merely is not connected now', async () => {
		const { registry } = await setup();
		expect(registry.peek('harbour')).toBeDefined(); // set() in the test's own setup()
		expect(registry.peek('never-seen')).toBeUndefined();
	});

	it('solo() is the one connected project when there is exactly one, so a caller naming none keeps working', async () => {
		const { registry, hub, url } = await setup();
		expect(registry.solo()).toBe(hub); // vivified by the test's own setup(), before any connection
		expect(registry.solo().status()).toEqual({ configured: true, connected: false });

		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(registry.solo()).toBe(hub);
		expect(registry.solo().status().connected).toBe(true);

		// a second project appears: naming none is now ambiguous, not silently the first project
		const second = await connect(url, KEY, undefined, {}, REQUIRED_METHODS, {
			key: 'quay',
			name: 'Quay'
		});
		cleanup.push(() => second.ws.terminate());
		expect(registry.solo()).not.toBe(hub);
		expect(registry.solo().status()).toEqual({ configured: true, connected: false });
		// each project on its own is still reachable by name
		expect(registry.hub('harbour').status().connected).toBe(true);
		expect(registry.hub('quay').status().connected).toBe(true);
	});
});
