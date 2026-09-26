// S-0095: live changes reach each project's Repo, and so its /api/events stream, however late that
// project's flai connects. A Repo used to bind to whichever hub agent() resolved the moment it
// started watching: made at startup, before any flai, that was a placeholder that never connects,
// and every change afterwards went nowhere. Tested against the real process-wide registry with real
// WebSocket connections that prove the credential, the way flai serve does.
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import http from 'node:http';
import { mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import type { AddressInfo } from 'node:net';
import WebSocket from 'ws';
import { proof, registry, REQUIRED_METHODS, resetAgent, withProject } from './agent';
import { knownRepos, Repo, repo, useRepo, type Ask } from './repo';

const KEY = 'shared-credential-for-tests';

type Flai = { ws: WebSocket; change: (path: string) => void };

/** What a fake flai answers a request with; by default, its project's manifest for every method. */
type Answer = (method: string) => unknown;

function connect(
	url: string,
	project: { key: string; name: string },
	answer?: Answer
): Promise<Flai> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(url);
		const mine = 'd'.repeat(32);
		let proven = false;
		const flai: Flai = {
			ws,
			change: (path) =>
				ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						method: 'change',
						params: { project: project.key, path }
					})
				)
		};
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
				return setTimeout(() => resolve(flai), 30);
			}
			if (typeof m.id === 'number')
				ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						id: m.id,
						result: answer?.(m.method) ?? {
							name: project.name,
							version: 1,
							layout: { design: 'design', docs: 'docs', wip: 'wip' }
						}
					})
				);
		});
		ws.on('error', reject);
	});
}

/** The paths a Repo announces as changed, as /api/events would stream them. */
function heard(r: Repo): string[] {
	const out: string[] = [];
	r.on('change', (p: string) => out.push(p));
	return out;
}

const settle = () => new Promise((r) => setTimeout(r, 50));

describe('live changes per project (S-0095)', () => {
	let dir: string;
	let server: http.Server;
	let url: string;
	const flais: Flai[] = [];

	beforeEach(async () => {
		resetAgent();
		for (const key of [...knownRepos().keys()]) useRepo(null, key);
		dir = mkdtempSync(join(tmpdir(), 'flaiover-events-'));
		writeFileSync(join(dir, 'agent-key'), KEY + '\n');
		process.env.FLAIOVER_AGENT_KEY_FILE = join(dir, 'agent-key');
		server = http.createServer((_req, res) => res.end('ok'));
		server.on('upgrade', (req, socket, head) => registry().handleUpgrade(req, socket, head));
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		url = `ws://127.0.0.1:${(server.address() as AddressInfo).port}/agent`;
	});

	afterEach(() => {
		for (const f of flais.splice(0)) f.ws.terminate();
		for (const key of [...knownRepos().keys()]) useRepo(null, key);
		resetAgent();
		server.close();
		delete process.env.FLAIOVER_AGENT_KEY_FILE;
		rmSync(dir, { recursive: true, force: true });
	});

	it('a Repo watching from startup, before any flai, hears the first one to connect', async () => {
		const r = repo(); // the default project, as hooks.server.ts's init watches it
		await r.watch();
		const paths = heard(r);

		const flai = await connect(url, { key: 'harbour', name: 'Harbour' });
		flais.push(flai);
		flai.change('wip/kanban/stories/S-0001-a.md');
		await settle();

		expect(paths).toEqual(['system-flow.yaml', 'wip/kanban/stories/S-0001-a.md']);
	});

	it('two projects: each Repo hears its own flai only, even watching before either connected', async () => {
		const a = withProject('alpha', () => repo());
		const b = withProject('beta', () => repo());
		await a.watch();
		await b.watch();
		const heardA = heard(a);
		const heardB = heard(b);

		const fa = await connect(url, { key: 'alpha', name: 'Alpha' });
		const fb = await connect(url, { key: 'beta', name: 'Beta' });
		flais.push(fa, fb);
		fa.change('wip/kanban/board.md');
		fb.change('wip/kanban/stories/S-0002-b.md');
		await settle();

		expect(heardA).toEqual(['system-flow.yaml', 'wip/kanban/board.md']);
		expect(heardB).toEqual(['system-flow.yaml', 'wip/kanban/stories/S-0002-b.md']);
	});

	it('the default Repo stops following a project once a second one connects', async () => {
		const r = repo();
		await r.watch();
		const paths = heard(r);
		const fa = await connect(url, { key: 'alpha', name: 'Alpha' });
		flais.push(fa);
		const fb = await connect(url, { key: 'beta', name: 'Beta' });
		flais.push(fb);
		fa.change('after-a-second-project.md');
		await settle();
		// with two projects there is no "the one project": a page must name one, and does
		expect(paths).toEqual(['system-flow.yaml']);
	});

	it("a project's Repo asks that project's flai, even outside the request that made it", async () => {
		const fa = await connect(url, { key: 'alpha', name: 'Alpha' });
		const fb = await connect(url, { key: 'beta', name: 'Beta' });
		flais.push(fa, fb);
		const b = withProject('beta', () => repo());
		// no project context here, as for the notifier or a change handler
		expect((await b.manifest()).name).toBe('Beta');
	});

	// S-0117: a Repo made for a request that names its project was never watched, so its answers
	// were kept until the container restarted: stories accepted and archived elsewhere stayed in the
	// done lane of /api/board?project=sf while /api/board showed them gone.
	it("a named project's Repo from repo() forgets its board when that project's flai reports a change", async () => {
		let stories = ['S-0028', 'S-0115', 'S-0116'];
		const asked: string[] = [];
		const sf = await connect(url, { key: 'sf', name: 'SF' }, (method) => {
			asked.push(method);
			return method === 'board.get' ? { done: [...stories] } : undefined;
		});
		flais.push(sf);
		// as hooks.server.ts runs a request with ?project=sf; nothing calls watch() on it
		const r = withProject('sf', () => repo());

		expect(await r.boardView()).toEqual({ done: ['S-0028', 'S-0115', 'S-0116'] });
		stories = [];
		sf.change('wip/kanban/stories/S-0115-accepted-elsewhere.md');
		await settle();

		expect(await withProject('sf', () => repo().boardView())).toEqual({ done: [] });
		expect(asked.filter((m) => m === 'board.get')).toHaveLength(2);
	});

	it('a Repo given its own source does not listen, and watching twice listens once', async () => {
		const reg = registry();
		const before = reg.listenerCount('change');
		const own = new Repo('/own', (async () => ({})) as Ask);
		await own.watch();
		expect(reg.listenerCount('change')).toBe(before);

		const r = withProject('gamma', () => repo());
		await r.watch();
		await r.watch();
		expect(reg.listenerCount('change')).toBe(before + 1);
		expect(reg.listenerCount('connected')).toBeGreaterThan(0);
		expect(reg.listenerCount('gone')).toBeGreaterThan(0);
	});
});
