import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import type { LogRead } from '$lib/logs';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask } from '$lib/server/repo';

// /api/log asks flai's log.read (S-0785) of a stubbed host: which paths flai reads is flai's to
// say, and its own tests say it.
type Handler = (event: never) => Promise<Response>;

const PATH = '.flai-cache/kluster/logs/s-0595/build-api.log';

const read = (over: Partial<LogRead> = {}): LogRead => ({
	path: PATH,
	exists: true,
	offset: 12,
	next: 30,
	size: 30,
	text: '#2 [build 1/4] FROM\n',
	...over
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return s.data;
}) as Ask;

describe('/api/log (S-0785)', () => {
	let GET: Handler;
	const get = (query = '') => GET({ url: new URL(`http://x/api/log${query}`) } as never);
	const reads = () => asked.filter((a) => a.method === 'log.read');

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'log.read': { data: read() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('reads a log from an offset, and from 0 when none is given, asking afresh each time', async () => {
		const r = await get(`?path=${encodeURIComponent(PATH)}&offset=12`);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(read());
		const anew = read({ offset: 0, next: 4, size: 4, text: 'new\n', reset: true });
		script['log.read'] = { data: anew };
		const again = (await (await get(`?path=${encodeURIComponent(PATH)}`)).json()) as LogRead;
		expect(again).toMatchObject({ offset: 0, text: 'new\n', reset: true });
		expect(reads().map((a) => a.params)).toEqual([
			{ path: PATH, offset: 12 },
			{ path: PATH, offset: 0 }
		]);
	});

	it('passes an absolute or ~/ path to flai as given', async () => {
		const given = ['/home/op/.flai/serve/agents/s-0785.log', '~/.flai/serve/flai-serve.log'];
		for (const path of given) {
			script['log.read'] = { data: read({ path }) };
			expect((await get(`?path=${encodeURIComponent(path)}`)).status).toBe(200);
		}
		expect(reads().map((a) => a.params.path)).toEqual(given);
	});

	it('answers a log that is not there as flai does', async () => {
		script['log.read'] = { data: read({ exists: false, offset: 0, next: 0, size: 0, text: '' }) };
		const r = await get(`?path=${encodeURIComponent(PATH)}`);
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ exists: false, text: '' });
	});

	it('refuses a missing path and a bad offset with 400 before asking flai', async () => {
		for (const query of [
			'',
			'?offset=0',
			'?path=',
			`?path=${PATH}&offset=-1`,
			`?path=${PATH}&offset=1.5`,
			`?path=${PATH}&offset=x`,
			`?path=${PATH}&offset=`
		]) {
			const r = await get(query);
			expect(r.status, query).toBe(400);
			expect((await r.json()).error, query).toMatch(/^(path|offset) must/);
		}
		expect(reads()).toEqual([]);
	});

	it("keeps flai's refusal of a path it does not read word for word", async () => {
		const refusal =
			'log.read reads a log under the project root, ~/.flai, or the story worktrees, not /etc/passwd';
		for (const code of [-32011, -32602]) {
			script['log.read'] = { error: new AgentError(502, refusal, code) };
			const r = await get('?path=/etc/passwd');
			expect(r.status).toBe(400);
			expect((await r.json()).error).toBe(refusal);
		}
	});

	it('answers 501 telling to upgrade a host flai without log.read', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		const r = await get(`?path=${encodeURIComponent(PATH)}`);
		expect(r.status).toBe(501);
		const error = (await r.json()).error as string;
		expect(error).toContain('the host flai is too old to read a log');
		expect(error).toContain('flai self-upgrade');
	});
});
