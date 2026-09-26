import { afterEach, beforeAll, describe, expect, it } from 'vitest';
import http from 'node:http';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import type { AddressInfo } from 'node:net';
import WebSocket from 'ws';
import { GET } from './+server';
import { AgentRegistry, proof, registry, resetAgent, REQUIRED_METHODS } from '$lib/server/agent';

const KEY = 'agent-credential-for-tests';

beforeAll(() => {
	const dir = mkdtempSync(join(tmpdir(), 'flaiover-projects-'));
	const file = join(dir, 'agent-key');
	writeFileSync(file, KEY);
	process.env.FLAIOVER_AGENT_KEY_FILE = file;
});

// A minimal flai, over a real socket: enough to answer the three methods the glance route asks for
// (S-0080). Duplicated from agent.test.ts rather than shared, since it is the only thing this file
// needs from there.
function connect(
	url: string,
	answer: (method: string) => unknown,
	project: { key: string; name: string }
): Promise<WebSocket> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(url);
		const mine = 'c'.repeat(32);
		let proven = false;
		ws.on('open', () =>
			ws.send(
				JSON.stringify({
					jsonrpc: '2.0',
					id: 'hello',
					method: 'hello',
					params: { protocol: 1, nonce: mine, flai: '9.9.9', project, methods: REQUIRED_METHODS }
				})
			)
		);
		ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (!proven) {
				proven = true;
				ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						method: 'hello.prove',
						params: { proof: proof(KEY, 'flai', m.result.nonce, mine) }
					})
				);
				return setTimeout(() => resolve(ws), 30);
			}
			const result = answer(m.method);
			if (result !== undefined) ws.send(JSON.stringify({ jsonrpc: '2.0', id: m.id, result }));
			// undefined means: never answer, to test a stall.
		});
		ws.on('error', reject);
	});
}

describe('/api/projects (S-0080)', () => {
	afterEach(() => resetAgent());

	it('is empty when nothing has connected', async () => {
		const body = await (await GET({} as never)).json();
		expect(body.projects).toEqual([]);
	});

	it('lists every project the registry knows, connected or not, without a glance for one not connected', async () => {
		registry().hub('harbour');
		const body = await (await GET({} as never)).json();
		const row = body.projects.find((p: { key: string }) => p.key === 'harbour');
		expect(row).toMatchObject({ key: 'harbour', connected: false });
		expect(row.review).toBeUndefined();
	});

	it('gives each connected project a glance: stories in review, threads awaiting, an agent attending', async () => {
		const server = http.createServer();
		server.on('upgrade', (req, socket, head) =>
			(registry() as unknown as AgentRegistry).handleUpgrade(req, socket, head)
		);
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		const port = (server.address() as AddressInfo).port;
		const answer = (method: string) => {
			if (method === 'board.get')
				return { columns: { review: [{ type: 'story' }, { type: 'story' }, { type: 'task' }] } };
			if (method === 'inbox.designer') return { counts: { thread: 2 } };
			if (method === 'agent.status') return { state: { running: { story: 'S-0001' } } };
			return {};
		};
		const flai = await connect(`ws://127.0.0.1:${port}/agent`, answer, {
			key: 'harbour',
			name: 'Harbour'
		});
		try {
			const body = await (await GET({} as never)).json();
			const row = body.projects.find((p: { key: string }) => p.key === 'harbour');
			expect(row).toMatchObject({
				connected: true,
				review: 2, // the task in the review column does not count
				threadsAwaiting: 2,
				agentAttending: true
			});
		} finally {
			flai.terminate();
			server.close();
		}
	});

	it('a project whose flai does not answer shows without its glance, and does not hold up another’s', async () => {
		const server = http.createServer();
		server.on('upgrade', (req, socket, head) =>
			(registry() as unknown as AgentRegistry).handleUpgrade(req, socket, head)
		);
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		const port = (server.address() as AddressInfo).port;
		const stalled = await connect(`ws://127.0.0.1:${port}/agent`, () => undefined, {
			key: 'stalled',
			name: 'Stalled'
		});
		const fine = await connect(
			`ws://127.0.0.1:${port}/agent`,
			(method) =>
				method === 'board.get'
					? { columns: {} }
					: method === 'inbox.designer'
						? { counts: {} }
						: { state: {} },
			{ key: 'fine', name: 'Fine' }
		);
		try {
			const started = Date.now();
			const body = await (await GET({} as never)).json();
			const elapsed = Date.now() - started;
			const stalledRow = body.projects.find((p: { key: string }) => p.key === 'stalled');
			const fineRow = body.projects.find((p: { key: string }) => p.key === 'fine');
			expect(stalledRow).toMatchObject({ connected: true });
			expect(stalledRow.review).toBeUndefined();
			expect(fineRow).toMatchObject({ connected: true, review: 0, threadsAwaiting: 0 });
			// both projects were asked at once, on the glance route's own short timeout, not flai's
			// usual 15s default: the whole call returns in a few seconds, not one project's timeout
			// added to another's.
			expect(elapsed).toBeLessThan(10000);
		} finally {
			stalled.terminate();
			fine.terminate();
			server.close();
		}
	}, 15000);

	it('adds what flai serve serves: a served project not connected carries why, and one never connected here is listed (S-0122)', async () => {
		const server = http.createServer();
		server.on('upgrade', (req, socket, head) =>
			(registry() as unknown as AgentRegistry).handleUpgrade(req, socket, head)
		);
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		const port = (server.address() as AddressInfo).port;
		registry().hub('quay'); // named here once, not connected now
		const answer = (method: string) => {
			if (method === 'settings.get')
				return {
					host: {
						projects: {
							running: true,
							served: [
								{
									key: 'harbour',
									name: 'Harbour',
									root: '/r/harbour',
									state: 'connected',
									settings: true
								},
								{
									key: 'quay',
									name: 'Quay',
									root: '/r/quay',
									state: 'not-connected',
									last_error: 'dial: connection refused',
									settings: true
								},
								{
									key: 'dock',
									name: 'Dock',
									root: '/r/dock',
									state: 'unavailable',
									reason: 'no system-flow.yaml',
									settings: true
								}
							],
							unserved: [
								{
									key: 'loci',
									name: 'loci',
									root: '/i/loci',
									reason: 'not registered',
									settings: true
								}
							]
						}
					}
				};
			if (method === 'board.get') return { columns: {} };
			if (method === 'inbox.designer') return { counts: {} };
			return { state: {} };
		};
		const flai = await connect(`ws://127.0.0.1:${port}/agent`, answer, {
			key: 'harbour',
			name: 'Harbour'
		});
		try {
			const body = await (await GET({} as never)).json();
			const rows = Object.fromEntries(
				body.projects.map((p: { key: string }) => [p.key, p])
			) as Record<string, Record<string, unknown>>;
			expect(Object.keys(rows)).toEqual(['dock', 'harbour', 'quay']); // not loci: it is not served
			expect(rows.harbour).toMatchObject({ connected: true, served: true });
			expect(rows.harbour.lastError).toBeUndefined();
			expect(rows.quay).toMatchObject({
				connected: false,
				served: true,
				lastError: 'dial: connection refused'
			});
			expect(rows.dock).toMatchObject({
				name: 'Dock',
				connected: false,
				served: true,
				lastError: 'no system-flow.yaml'
			});
		} finally {
			flai.terminate();
			server.close();
		}
	});

	it('without a connected project to ask, lists what the registry knows, none of it marked served', async () => {
		registry().hub('quay');
		const body = await (await GET({} as never)).json();
		expect(body.projects).toEqual([{ key: 'quay', name: 'quay', connected: false }]);
	});
});
