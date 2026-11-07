import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask, type KlusterLog } from '$lib/server/repo';

// /api/kluster/log asks flai's kluster.log (ADR-0263) of a stubbed host: the flai side's tests read
// real log files; this one checks the query, what is asked, and what is answered.
type Handler = (event: never) => Promise<Response>;

const read = (over: Partial<KlusterLog> = {}): KlusterLog => ({
	namespace: 's-0595',
	step: 'build-api',
	log: '.flai-cache/kluster/logs/s-0595/build-api.log',
	exists: true,
	offset: 0,
	next: 12,
	size: 12,
	text: 'built api:1\n',
	outcome: 'running',
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

describe('/api/kluster/log (ADR-0263)', () => {
	let GET: Handler;
	const get = (query: string) => GET({ url: new URL(`http://x/api/kluster/log${query}`) } as never);
	const asks = () => asked.filter((a) => a.method === 'kluster.log');

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'kluster.log': { data: read() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it("reads a step's log from an offset and answers flai's object", async () => {
		const r = await get('?namespace=s-0595&step=build-api&offset=12');
		expect(r.status).toBe(200);
		expect(asks().map((a) => a.params)).toEqual([
			{ namespace: 's-0595', step: 'build-api', offset: 12 }
		]);
		expect(await r.json()).toMatchObject(read());
	});

	it('reads from 0 when no offset is given', async () => {
		for (const step of ['build', 'load', 'deploy', 'teardown']) {
			expect((await get(`?namespace=authstar&step=${step}`)).status).toBe(200);
		}
		expect(asks().map((a) => a.params.offset)).toEqual([0, 0, 0, 0]);
	});

	it('refuses a query that is not a namespace, a step, and an offset before asking flai', async () => {
		for (const query of [
			'',
			'?step=build',
			'?namespace=S_0595&step=build',
			'?namespace=s-0595',
			'?namespace=s-0595&step=../x',
			'?namespace=s-0595&step=build-',
			'?namespace=s-0595&step=build-API',
			'?namespace=s-0595&step=push',
			'?namespace=s-0595&step=build&offset=-1',
			'?namespace=s-0595&step=build&offset=1.5',
			'?namespace=s-0595&step=build&offset=x'
		]) {
			const r = await get(query);
			expect(r.status, query).toBe(400);
		}
		expect(asks()).toEqual([]);
	});

	it("answers flai's refusal word for word as 400", async () => {
		const bad = 'step must be build, build-<name>, load, deploy, or teardown, not "x"';
		script['kluster.log'] = { error: new AgentError(502, bad, -32602) };
		const r = await get('?namespace=s-0595&step=build');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(bad);
	});

	it('answers 501 telling to upgrade a host flai without kluster.log', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		const r = await get('?namespace=s-0595&step=build');
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('the host flai is too old');
	});
});
