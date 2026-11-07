import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

// /api/kluster/config reads the kluster config file the kluster: block names through flai's
// file.get and saves it through file.save (ADR-0147, S-0586), against a stubbed host.
const CONFIG = 'kluster.yaml';
const CONTENT = 'name: p-local\nprovider: k3d\nprofile: spire\n';
const HASH = 'a'.repeat(64);

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// kluster.status and project.info answer raw, as flai's reads do; file.get and file.save as commands
const RAW = new Set(['kluster.status', 'project.info']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return RAW.has(method) ? s.data : { data: s.data, warnings: ['saved'] };
}) as Ask;

const status = (block: boolean) => ({
	data: { enabled: true, enable: '', block, ...(block ? { config: CONFIG } : {}) }
});

describe('/api/kluster/config (S-0586)', () => {
	let GET: Handler;
	let PUT: Handler;
	const get = () => GET({ url: new URL('http://x/api/kluster/config') } as never);
	const put = (body: unknown) =>
		PUT({
			request: new Request('http://x/api/kluster/config', {
				method: 'PUT',
				body: JSON.stringify(body)
			})
		} as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		const mod = await import('./+server');
		GET = mod.GET as unknown as Handler;
		PUT = mod.PUT as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = { 'project.info': { data: { name: 'p', key: 'p' } }, 'kluster.status': status(true) };
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('reads the config file kluster.status names through file.get, with its hash', async () => {
		script['file.get'] = { data: { path: CONFIG, content: CONTENT, hash: HASH } };
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ path: CONFIG, content: CONTENT, hash: HASH });
		expect(asks('file.get')[0].params).toEqual({ path: CONFIG });
	});

	it('answers 404 with no block, without asking file.get', async () => {
		script['kluster.status'] = status(false);
		const r = await get();
		expect(r.status).toBe(404);
		expect(asks('file.get')).toEqual([]);
	});

	it('saves the config file through file.save, and refuses any other path', async () => {
		script['file.save'] = { data: { path: CONFIG, hash: 'b'.repeat(64), committed: 'abc1234' } };
		const r = await put({ path: CONFIG, content: CONTENT, hash: HASH, message: 'profile' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ path: CONFIG, committed: 'abc1234', log: ['saved'] });
		expect(asks('file.save')[0].params).toMatchObject({
			path: CONFIG,
			content: CONTENT,
			hash: HASH,
			message: 'profile'
		});
		const other = await put({ path: 'system-flow.yaml', content: 'x', hash: HASH });
		expect(other.status).toBe(400);
		expect((await other.json()).error).toContain('only the kluster config file kluster.yaml');
		expect(asks('file.save')).toHaveLength(1);
	});

	it('answers 400 without a path, content, or hash', async () => {
		const r = await put({ path: CONFIG });
		expect(r.status).toBe(400);
		expect(asks('file.save')).toEqual([]);
	});
});
