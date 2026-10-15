import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

// /api/design/tokens reads the token file through flai's file.get and saves it through file.save
// (ADR-0147, S-0383), against a stubbed host: the flai side of the two methods is a later story's.
const TOKENS = 'design/system/design-system.tokens.json';
const CONTENT = '{\n  "color": {}\n}\n';
const HASH = 'a'.repeat(64);

type Handler = (event: never) => Promise<Response>;

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// design.get and project.info answer raw, as flai's reads do; file.get and file.save as commands do
const RAW = new Set(['design.get', 'project.info']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return RAW.has(method) ? s.data : { data: s.data, warnings: ['saved'] };
}) as Ask;

const design = (exists: boolean) => ({
	data: { path: TOKENS, exists, ...(exists ? {} : { no_token_file: 'run flai design init' }) }
});

describe('/api/design/tokens (S-0383)', () => {
	let GET: Handler;
	let PUT: Handler;
	const get = () => GET({ url: new URL('http://x/api/design/tokens') } as never);
	const put = (body: unknown) =>
		PUT({
			request: new Request('http://x/api/design/tokens', {
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
		script = { 'project.info': { data: { name: 'p', key: 'p' } }, 'design.get': design(true) };
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('reads the token file design.get names through file.get, with its hash', async () => {
		script['file.get'] = { data: { path: TOKENS, content: CONTENT, hash: HASH } };
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ path: TOKENS, content: CONTENT, hash: HASH });
		expect(asks('file.get')[0].params).toEqual({ path: TOKENS });
	});

	it('answers 404 when the project has no token file, without asking file.get', async () => {
		script['design.get'] = design(false);
		const r = await get();
		expect(r.status).toBe(404);
		expect((await r.json()).error).toContain('flai design init');
		expect(asks('file.get')).toEqual([]);
	});

	it('answers 501 to a read and a save on a host flai without file.get and file.save', async () => {
		const read = await get();
		expect(read.status).toBe(501);
		const error = (await read.json()).error;
		expect(error).toContain('too old to read the design token file');
		expect(error).toContain('upgrade');
		const save = await put({ path: TOKENS, content: CONTENT, hash: HASH });
		expect(save.status).toBe(501);
		expect((await save.json()).error).toContain('too old to save the design token file');
	});

	it('saves the token file through file.save on main and forgets design.get', async () => {
		script['file.save'] = { data: { path: TOKENS, hash: 'b'.repeat(64), committed: true } };
		await get();
		const r = await put({ path: TOKENS, content: CONTENT, hash: HASH, message: 'warmer primary' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			path: TOKENS,
			hash: 'b'.repeat(64),
			committed: true,
			log: ['saved']
		});
		const params = asks('file.save')[0].params;
		expect(params).toMatchObject({
			path: TOKENS,
			content: CONTENT,
			hash: HASH,
			message: 'warmer primary'
		});
		expect(params.story).toBeUndefined();
		const before = asks('design.get').length;
		await get();
		expect(asks('design.get').length).toBe(before + 1);
	});

	it('refuses to save another path with 400, before asking file.save', async () => {
		const r = await put({ path: 'design/system/other.json', content: CONTENT, hash: HASH });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain(TOKENS);
		expect(asks('file.save')).toEqual([]);
	});

	it('requires path, content, and hash', async () => {
		const bodies = [
			{ content: CONTENT, hash: HASH },
			{ path: TOKENS, hash: HASH },
			{ path: TOKENS, content: CONTENT }
		];
		for (const body of bodies) expect((await put(body)).status, JSON.stringify(body)).toBe(400);
		expect(asks('file.save')).toEqual([]);
	});

	it('answers a stale hash with 409, the current content, its hash, and a diff', async () => {
		const current = '{\n  "color": { "primary": {} }\n}\n';
		script['file.save'] = {
			error: new AgentError(502, `${TOKENS} changed after it was loaded`, -32009, {
				current,
				hash: 'c'.repeat(64),
				diff: '+mine'
			})
		};
		const r = await put({ path: TOKENS, content: CONTENT, hash: HASH });
		expect(r.status).toBe(409);
		expect(await r.json()).toMatchObject({
			error: `${TOKENS} changed after it was loaded`,
			current,
			hash: 'c'.repeat(64),
			diff: '+mine'
		});
	});

	it("answers flai's refusal with 422 and its findings", async () => {
		const findings = [{ rule: 'design.tokens', path: TOKENS, line: 2, message: 'bad colour' }];
		script['file.save'] = {
			error: new AgentError(502, 'flai check refused the save', -32010, { findings })
		};
		const r = await put({ path: TOKENS, content: CONTENT, hash: HASH });
		expect(r.status).toBe(422);
		expect(await r.json()).toMatchObject({ error: 'flai check refused the save', findings });
	});
});
