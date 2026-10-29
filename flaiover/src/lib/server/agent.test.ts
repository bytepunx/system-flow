import { afterEach, describe, expect, it, vi } from 'vitest';
import http from 'node:http';
import { generateKeyPairSync, sign, type KeyObject } from 'node:crypto';
import { mkdtempSync, writeFileSync } from 'node:fs';
import type { AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import WebSocket from 'ws';
import {
	AgentError,
	AgentRegistry,
	type AgentHub,
	ANALYZER,
	CLOSE_HELD,
	CLOSE_UNSIGNED,
	heldBudget,
	ORCHESTRATOR,
	PROTOCOL,
	proof,
	REQUIRED_METHODS,
	registry as processRegistry,
	resetAgent,
	type RegistryOptions
} from './agent';
import { log } from './log';

const KEY = 'agent-credential-for-tests';

// S-0237: every flai here sends a stamp signed by a release key the test makes, and the registry is
// given that key in place of the embedded one, so a flai is signed unless a test says otherwise.
function keyPair() {
	const { publicKey, privateKey } = generateKeyPairSync('ec', { namedCurve: 'P-256' });
	return { pem: publicKey.export({ type: 'spki', format: 'pem' }).toString(), privateKey };
}
const release = keyPair();
const other = keyPair();
function stamp(statement: string, key: KeyObject = release.privateKey) {
	return { statement, signature: sign('sha256', Buffer.from(statement), key).toString('base64') };
}
const SIGNED = stamp('flai 9.9.9 abc1234');
const DASHBOARD = stamp('flaiover 1.2.3 def5678');

type Flai = {
	ws: WebSocket;
	requests: Record<string, unknown>[];
	closed: Promise<number>;
	/** The dashboard's answer to hello. */
	hello?: Record<string, unknown>;
};

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
	project: { key: string; name: string } = { key: 'harbour', name: 'Harbour' },
	extra: Record<string, unknown> = {}
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
						protocol: PROTOCOL,
						nonce: mine,
						flai: '9.9.9',
						project,
						methods,
						release: SIGNED,
						...extra
					}
				})
			)
		);
		ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (!proven) {
				proven = true;
				flai.hello = m.result;
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

/** A flai that only says hello, as a refused one gets no further: how the dashboard closed, and
 * everything it sent first. No release stamp unless params gives one. */
function sayHello(
	url: string,
	params: Record<string, unknown>
): Promise<{ code: number; reason: string; messages: unknown[] }> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(url);
		const messages: unknown[] = [];
		ws.on('open', () =>
			ws.send(
				JSON.stringify({
					jsonrpc: '2.0',
					id: 'hello',
					method: 'hello',
					params: {
						protocol: PROTOCOL,
						nonce: 'c'.repeat(32),
						flai: '9.9.9',
						project: { key: 'harbour', name: 'Harbour' },
						methods: REQUIRED_METHODS,
						...params
					}
				})
			)
		);
		ws.on('message', (data) => messages.push(JSON.parse(data.toString())));
		ws.on('close', (code, reason) => resolve({ code, reason: reason.toString(), messages }));
		ws.on('error', reject);
	});
}

describe('AgentRegistry and AgentHub', () => {
	const cleanup: (() => void)[] = [];
	afterEach(() => {
		for (const f of cleanup.splice(0)) f();
	});
	async function setup(opt: RegistryOptions = {}) {
		const registry = new AgentRegistry(KEY, { peerKeys: [release.pem], stamp: DASHBOARD, ...opt });
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

	it("sends the channel's project key over any project a call's params carry (S-0448)", async () => {
		const { hub, url } = await setup();
		const flai = await connect(url, KEY, () => ({ projects: [] }));
		cleanup.push(() => flai.ws.terminate());
		await hub.ask('components.list', { project: 'flaiover', story: 'S-0448' });
		expect(flai.requests[0]).toMatchObject({ method: 'components.list' });
		expect(flai.requests[0].params).toEqual({ project: 'harbour', story: 'S-0448' });
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

	it('refuses a newer proven connection while the one the project has answers (S-0184)', async () => {
		const { hub, url } = await setup();
		const first = await connect(url, KEY, () => ({ from: 'first' }));
		cleanup.push(() => first.ws.terminate());
		const second = await connect(url, KEY, () => ({ from: 'second' }));
		cleanup.push(() => second.ws.terminate());
		expect(await second.closed).toBe(CLOSE_HELD);
		expect(first.ws.readyState).toBe(WebSocket.OPEN);
		await expect(hub.ask('project.info')).resolves.toEqual({ from: 'first' });
	});

	it('replaces the one the project has when it does not answer (S-0184)', async () => {
		const { hub, url } = await setup({ probeMs: 120 });
		const first = await connect(url, KEY, () => ({ from: 'first' }));
		cleanup.push(() => first.ws.terminate());
		// A frozen process: the socket stays open and nothing answers, pongs included.
		(
			first.ws as unknown as { _receiver: { removeAllListeners: (e: string) => void } }
		)._receiver.removeAllListeners('ping');
		const changed: string[] = [];
		hub.on('change', (path: string) => changed.push(path));
		const second = await connect(url, KEY, () => ({ from: 'second' }));
		cleanup.push(() => second.ws.terminate());
		// sent while it is still being judged: kept, and handled once it is adopted
		second.ws.send(
			JSON.stringify({ jsonrpc: '2.0', method: 'change', params: { path: 'wip/kanban/board.md' } })
		);
		expect(changed).toEqual([]);
		await new Promise((r) => setTimeout(r, 200));
		expect(second.ws.readyState).toBe(WebSocket.OPEN);
		expect(changed).toEqual(['wip/kanban/board.md']);
		await expect(hub.ask('project.info')).resolves.toEqual({ from: 'second' });
		first.ws.terminate();
		expect(await first.closed).not.toBe(CLOSE_HELD);
	});

	it('judges two newer connections in turn, and refuses both while the holder answers', async () => {
		const { hub, url } = await setup();
		const first = await connect(url, KEY, () => ({ from: 'first' }));
		cleanup.push(() => first.ws.terminate());
		const [second, third] = await Promise.all([
			connect(url, KEY, () => ({ from: 'second' })),
			connect(url, KEY, () => ({ from: 'third' }))
		]);
		cleanup.push(
			() => second.ws.terminate(),
			() => third.ws.terminate()
		);
		expect(await second.closed).toBe(CLOSE_HELD);
		expect(await third.closed).toBe(CLOSE_HELD);
		await expect(hub.ask('project.info')).resolves.toEqual({ from: 'first' });
	});

	it('adopts a newer connection at once when the one the project had has closed', async () => {
		const { hub, url } = await setup();
		const first = await connect(url, KEY, () => ({ from: 'first' }));
		first.ws.close();
		await first.closed;
		const second = await connect(url, KEY, () => ({ from: 'second' }));
		cleanup.push(() => second.ws.terminate());
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

	// S-0154: an agent starting or ending changes no file, and flai serve says so of itself.
	it("passes on flai serve's word that a story's agent started or ended", async () => {
		const { hub, url } = await setup();
		const stories: string[] = [];
		hub.on('agent', (story: string) => stories.push(story));
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		flai.ws.send(
			JSON.stringify({
				jsonrpc: '2.0',
				method: 'agent',
				params: { project: 'harbour', story: 'S-0154' }
			})
		);
		// S-0208: the planner for an epic, which names the item in place of the story
		flai.ws.send(
			JSON.stringify({
				jsonrpc: '2.0',
				method: 'agent',
				params: { project: 'harbour', item: 'E-0016' }
			})
		);
		// S-0218: the project's orchestrator, which names neither and has the role orchestrate
		flai.ws.send(
			JSON.stringify({
				jsonrpc: '2.0',
				method: 'agent',
				params: { project: 'harbour', role: 'orchestrate' }
			})
		);
		// S-0228: the project's analyzer, with the role analyze
		flai.ws.send(
			JSON.stringify({
				jsonrpc: '2.0',
				method: 'agent',
				params: { project: 'harbour', role: 'analyze' }
			})
		);
		// not a story, not a role the dashboard knows, and not a notification: none is announced
		flai.ws.send(JSON.stringify({ jsonrpc: '2.0', method: 'agent', params: { story: 7 } }));
		flai.ws.send(
			JSON.stringify({ jsonrpc: '2.0', method: 'agent', params: { project: 'harbour', role: 'x' } })
		);
		flai.ws.send(
			JSON.stringify({ jsonrpc: '2.0', id: 98, method: 'agent', params: { story: 'S-1' } })
		);
		await new Promise((r) => setTimeout(r, 50));
		expect(stories).toEqual(['S-0154', 'E-0016', ORCHESTRATOR, ANALYZER]);
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

	// S-0098: a repository offered for import connects as a candidate. It answers only the import
	// methods by design, so nothing counts as missing; it is listed while connected and not after
	// (it was imported, or its folder is no longer named); and it is never "the one project".
	it('adopts a candidate as a repository offered for import', async () => {
		const { registry, url } = await setup();
		const project = await connect(url, KEY, () => ({ name: 'Harbour' }));
		cleanup.push(() => project.ws.terminate());
		const offered = await connect(
			url,
			KEY,
			() => ({ tests: [] }),
			{},
			['import.preview', 'import.run'],
			{ key: 'import-widget', name: 'widget' },
			{ kind: 'candidate' }
		);
		expect(registry.peek('import-widget')?.status()).toMatchObject({
			connected: true,
			candidate: true,
			missing: []
		});
		expect(registry.list()).toEqual([
			expect.objectContaining({ key: 'harbour', connected: true }),
			expect.objectContaining({ key: 'import-widget', candidate: true, connected: true })
		]);
		expect(registry.list()[0].candidate).toBeUndefined();
		// one project and a candidate: a request with no ?project= still means the project
		expect(registry.solo()).toBe(registry.peek('harbour'));

		offered.ws.terminate();
		await offered.closed;
		await new Promise((r) => setTimeout(r, 20));
		expect(registry.list().map((p) => p.key)).toEqual(['harbour']);
	});

	// S-0121: flai serve says `removed` before it drops a project that left its registry. The
	// project leaves the list once its connection closes, while one that only lost its flai (a
	// restart) stays; a new connection for the key brings it back.
	it('drops a project its flai said was removed, and takes it back when it connects again', async () => {
		const { registry, url } = await setup();
		const harbour = await connect(url, KEY);
		cleanup.push(() => harbour.ws.terminate());
		const quay = await connect(url, KEY, () => ({ name: 'Quay' }), {}, REQUIRED_METHODS, {
			key: 'quay',
			name: 'Quay'
		});
		quay.ws.send(
			JSON.stringify({ jsonrpc: '2.0', method: 'removed', params: { project: 'quay' } })
		);
		quay.ws.close();
		await quay.closed;
		harbour.ws.terminate();
		await harbour.closed;
		await new Promise((r) => setTimeout(r, 20));
		expect(registry.list()).toEqual([
			expect.objectContaining({ key: 'harbour', connected: false })
		]);
		expect(registry.solo()).toBe(registry.peek('harbour'));

		const back = await connect(url, KEY, () => ({ name: 'Quay' }), {}, REQUIRED_METHODS, {
			key: 'quay',
			name: 'Quay'
		});
		cleanup.push(() => back.ws.terminate());
		expect(registry.list().map((p) => p.key)).toEqual(['harbour', 'quay']);
	});

	// S-0237: flai's release stamp is judged before the credential, and the dashboard answers with its own.
	it('answers a signed flai with its own stamp, adopts it, and marks it signed', async () => {
		const { registry, hub, url } = await setup();
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(PROTOCOL).toBe(2);
		expect(flai.hello).toMatchObject({ release: DASHBOARD });
		expect(hub.status()).toMatchObject({ connected: true, signed: true });
		expect(hub.status().refused).toBeUndefined();
		expect(registry.list()).toEqual([
			expect.objectContaining({ key: 'harbour', connected: true, signed: true })
		]);
	});

	it('answers without a stamp when it has none, as a development build', async () => {
		const { url } = await setup({ stamp: null });
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(flai.hello).toHaveProperty('nonce');
		expect(flai.hello).not.toHaveProperty('release');
	});

	it('closes an unsigned flai with 4403 before the proof, and keeps the refusal against its project', async () => {
		const { registry, url } = await setup();
		const cases: { params: Record<string, unknown>; reason: string; flai: string }[] = [
			{ params: {}, reason: 'missing stamp', flai: '9.9.9' },
			// a flai from before the stamp: protocol 1, no field, refused as unsigned, not as malformed
			{ params: { protocol: 1, flai: '1.40.0' }, reason: 'missing stamp', flai: '1.40.0' },
			{
				params: { release: { statement: '', signature: '' } },
				reason: 'missing stamp',
				flai: '9.9.9'
			},
			{ params: { release: DASHBOARD }, reason: 'wrong component flaiover', flai: '9.9.9' },
			{
				params: { release: stamp('flai 9.9.8 abc1234', other.privateKey) },
				reason: 'invalid signature',
				flai: '9.9.8'
			},
			{
				params: { release: { statement: 'flai 9.9.9', signature: SIGNED.signature } },
				reason: 'invalid signature',
				flai: '9.9.9'
			}
		];
		for (const c of cases) {
			const refused = await sayHello(url, c.params);
			// closed with nothing sent: no nonce to prove against, and no proof of the dashboard's
			expect(refused, c.reason).toEqual({
				code: CLOSE_UNSIGNED,
				reason: `unsigned peer: ${c.reason}`,
				messages: []
			});
			const kept = registry.peek('harbour')?.status().refused;
			expect(kept, c.reason).toMatchObject({ reason: c.reason, flai: c.flai });
			expect(new Date(kept!.at).toISOString()).toBe(kept!.at);
		}
		expect(CLOSE_UNSIGNED).toBe(4403);
		const [row] = registry.list();
		expect(row).toMatchObject({
			key: 'harbour',
			connected: false,
			refused: { reason: 'invalid signature', flai: '9.9.9' }
		});
		expect(row.signed).toBeUndefined();
	});

	it('forgets the refusal once a signed flai for the project is adopted', async () => {
		const { hub, url } = await setup();
		await sayHello(url, {});
		expect(hub.status()).toMatchObject({ connected: false, refused: { reason: 'missing stamp' } });
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(hub.status()).toMatchObject({ connected: true, signed: true });
		expect(hub.status().refused).toBeUndefined();
	});

	it('serves an unsigned flai, an older one too, marked unsigned, when told to allow it', async () => {
		const { hub, url } = await setup({ allowUnsigned: true });
		const flai = await connect(url, KEY, undefined, {}, REQUIRED_METHODS, undefined, {
			protocol: 1,
			release: undefined
		});
		cleanup.push(() => flai.ws.terminate());
		expect(hub.status()).toMatchObject({ connected: true, signed: false, unsigned: ['flai'] });
		expect(hub.status().refused).toBeUndefined();
	});

	it('refuses a stamp that does not verify or names another component even when told to allow unsigned', async () => {
		const { url } = await setup({ allowUnsigned: true });
		for (const release of [stamp('flai 9.9.9 abc1234', other.privateKey), DASHBOARD]) {
			const refused = await sayHello(url, { release });
			expect(refused.code, release.statement).toBe(CLOSE_UNSIGNED);
		}
	});

	it('still refuses a protocol it does not speak as a protocol, not as unsigned', async () => {
		const { url } = await setup();
		const refused = await sayHello(url, { protocol: 3, release: SIGNED });
		expect(refused).toMatchObject({ code: 4400, reason: 'protocol' });
		expect(refused.messages).toEqual([
			expect.objectContaining({
				error: expect.objectContaining({ message: 'this dashboard speaks protocol 2, flai sent 3' })
			})
		]);
	});

	it('makes a project known by its refusal alone, without making a proven one ambiguous', async () => {
		const { registry, url } = await setup(); // harbour, named here, is the one project
		await sayHello(url, { project: { key: 'quay', name: 'Quay' } });
		expect(registry.peek('quay')?.status()).toMatchObject({
			connected: false,
			refused: { reason: 'missing stamp' }
		});
		expect(registry.list().map((p) => p.key)).toEqual(['harbour', 'quay']);
		expect(registry.solo()).toBe(registry.peek('harbour'));
	});

	it('is the one project when a refusal is all it knows, and keeps at most eight such', async () => {
		const registry = new AgentRegistry(KEY, { peerKeys: [release.pem] });
		const s = await serve(registry);
		cleanup.push(() => {
			registry.close();
			s.server.close();
		});
		await sayHello(s.url, { project: { key: 'p0', name: '' } });
		expect(registry.solo()).toBe(registry.peek('p0'));
		expect(registry.solo().status().refused?.reason).toBe('missing stamp');
		for (let i = 1; i < 10; i++) await sayHello(s.url, { project: { key: `p${i}`, name: '' } });
		expect(registry.list()).toHaveLength(8);
		expect(registry.peek('p9')).toBeUndefined();
		expect(registry.solo().status()).toEqual({ configured: true, connected: false });
	});
});

// S-0239: the process registry reads FLAIOVER_ALLOW_UNSIGNED, which flai dashboard gives a container
// it starts as a development build, and each connection records which side it let through unsigned.
describe('the registry under FLAIOVER_ALLOW_UNSIGNED', () => {
	const cleanup: (() => void)[] = [];
	afterEach(() => {
		for (const f of cleanup.splice(0)) f();
		vi.restoreAllMocks();
	});

	const keyFile = join(mkdtempSync(join(tmpdir(), 'flaiover-agent-')), 'agent-key');
	writeFileSync(keyFile, KEY);

	/** The process registry, made from env as the server makes it, on a real HTTP server. */
	async function start(allow: string | undefined, opt: RegistryOptions = {}) {
		resetAgent();
		const env = { FLAIOVER_AGENT_KEY_FILE: keyFile, FLAIOVER_ALLOW_UNSIGNED: allow };
		const registry = processRegistry(env, { peerKeys: [release.pem], stamp: DASHBOARD, ...opt });
		const s = await serve(registry);
		cleanup.push(() => {
			resetAgent();
			s.server.close();
		});
		return { registry, hub: registry.hub('harbour'), ...s };
	}

	/** Whether a call to warn is the one that says a connection was let through unsigned. */
	const warnsUnsigned = ([fields]: unknown[]) =>
		typeof fields === 'object' && fields !== null && 'unsigned' in fields;

	/** A flai with no stamp, as scripts/flai.sh builds one, that completes the handshake. */
	async function unstamped(url: string, key = KEY) {
		const flai = await connect(url, key, undefined, {}, REQUIRED_METHODS, undefined, {
			release: undefined
		});
		cleanup.push(() => flai.ws.terminate());
		return flai;
	}

	it('closes a flai with no stamp with 4403 when the variable is unset or anything but 1', async () => {
		for (const allow of [undefined, '', '0', 'true']) {
			const { hub, url } = await start(allow);
			const refused = await sayHello(url, {});
			expect(refused, String(allow)).toMatchObject({
				code: CLOSE_UNSIGNED,
				reason: 'unsigned peer: missing stamp'
			});
			expect(hub.status().connected, String(allow)).toBe(false);
			for (const f of cleanup.splice(0)) f();
		}
	});

	it('serves a flai with no stamp when it is 1, recorded as the unsigned side flai, and warns once', async () => {
		const warn = vi.spyOn(log(), 'warn');
		const { registry, hub, url } = await start('1');
		await unstamped(url);
		expect(hub.status()).toMatchObject({ connected: true, signed: false, unsigned: ['flai'] });
		expect(registry.list()).toEqual([
			expect.objectContaining({
				key: 'harbour',
				connected: true,
				signed: false,
				unsigned: ['flai']
			})
		]);
		expect(warn.mock.calls.filter(warnsUnsigned)).toEqual([
			[
				{ component: 'agent', project: 'harbour', flai: '9.9.9', unsigned: ['flai'] },
				expect.any(String)
			]
		]);
	});

	it('forgets which side was unsigned once the flai disconnects', async () => {
		const { registry, hub, url } = await start('1');
		const flai = await unstamped(url);
		flai.ws.terminate();
		await flai.closed;
		await new Promise((r) => setTimeout(r, 20));
		expect(hub.status().connected).toBe(false);
		expect(hub.status().unsigned).toBeUndefined();
		expect(registry.list()[0].unsigned).toBeUndefined();
	});

	it('still refuses a bad signature or another component when it is 1', async () => {
		const { url } = await start('1');
		for (const stamped of [stamp('flai 9.9.9 abc1234', other.privateKey), DASHBOARD]) {
			const refused = await sayHello(url, { release: stamped });
			expect(refused.code, stamped.statement).toBe(CLOSE_UNSIGNED);
		}
	});

	it('still refuses the wrong credential when it is 1, signed or not', async () => {
		const { hub, url } = await start('1');
		const withoutStamp = await unstamped(url, 'some-other-key');
		expect(await withoutStamp.closed).toBe(4401);
		const signed = await connect(url, 'some-other-key');
		expect(await signed.closed).toBe(4401);
		expect(hub.status().connected).toBe(false);
	});

	it('records the dashboard as unsigned when this build has no stamp of its own, beside flai when both are', async () => {
		const { hub, url } = await start('1', { stamp: null });
		const signed = await connect(url, KEY);
		cleanup.push(() => signed.ws.terminate());
		expect(hub.status()).toMatchObject({ connected: true, signed: true, unsigned: ['dashboard'] });
		signed.ws.terminate();
		await signed.closed;
		await new Promise((r) => setTimeout(r, 20));
		await unstamped(url);
		expect(hub.status()).toMatchObject({
			connected: true,
			signed: false,
			unsigned: ['flai', 'dashboard']
		});
	});

	it('records the dashboard as unsigned with the variable unset when this build has no stamp, and warns once', async () => {
		const warn = vi.spyOn(log(), 'warn');
		const { registry, hub, url } = await start(undefined, { stamp: null });
		const signed = await connect(url, KEY);
		cleanup.push(() => signed.ws.terminate());
		expect(hub.status()).toMatchObject({ connected: true, signed: true, unsigned: ['dashboard'] });
		expect(registry.list()[0].unsigned).toEqual(['dashboard']);
		expect(warn.mock.calls.filter(warnsUnsigned)).toEqual([
			[
				expect.objectContaining({
					component: 'agent',
					project: 'harbour',
					unsigned: ['dashboard']
				}),
				expect.any(String)
			]
		]);
	});

	it('records nothing for a signed pair when it is 1, and warns nothing', async () => {
		const warn = vi.spyOn(log(), 'warn');
		const { registry, hub, url } = await start('1');
		const flai = await connect(url, KEY);
		cleanup.push(() => flai.ws.terminate());
		expect(hub.status()).toMatchObject({ connected: true, signed: true });
		expect(hub.status().unsigned).toBeUndefined();
		expect(registry.list()[0].unsigned).toBeUndefined();
		expect(warn.mock.calls.filter(warnsUnsigned)).toEqual([]);
	});

	it('lets an explicit option win over the variable', async () => {
		const { url } = await start('1', { allowUnsigned: false });
		expect((await sayHello(url, {})).code).toBe(CLOSE_UNSIGNED);
	});
});

// S-0534: the providers flai serve holds every agent start on, read from agent.status's state.budget.
describe('heldBudget', () => {
	const hold = {
		provider: 'claude-code',
		limit: 'usage_limit_reached',
		status: 429,
		since: '2026-10-28T05:30:00Z',
		until: '2026-10-28T09:00:00Z'
	};
	const hub = (answer: () => Promise<unknown>) => ({ ask: vi.fn(answer) as AgentHub['ask'] });
	const answering = (status: unknown) => hub(async () => status);

	it('asks agent.status and returns its budget', async () => {
		const h = answering({ enabled: true, state: { command: 'claude', budget: [hold] } });
		expect(await heldBudget(h, 1500)).toEqual([hold]);
		expect(h.ask).toHaveBeenCalledWith('agent.status', {}, 1500);
	});

	it('returns none while agent.status holds none', async () => {
		const answers = [
			{ enabled: false },
			{ enabled: true, state: { command: 'claude' } },
			{ enabled: true, state: { budget: null } },
			{ enabled: true, state: { budget: [] } }
		];
		for (const answer of answers) expect(await heldBudget(answering(answer), 1500)).toBeUndefined();
	});

	it('returns none when flai does not answer', async () => {
		const h = hub(async () => {
			throw new AgentError(504, 'the host flai did not answer agent.status in time');
		});
		expect(await heldBudget(h, 1500)).toBeUndefined();
	});
});
