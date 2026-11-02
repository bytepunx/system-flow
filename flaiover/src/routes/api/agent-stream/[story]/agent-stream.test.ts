import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { AgentError } from '$lib/server/agent';

type Handler = (event: never) => Promise<Response>;

// The channel, scripted: what flai on the host would answer for agent.stream (S-0142). How flai
// reads the agent's log is flai's business, tested there; here it is what the route asks for.
let asked: { method: string; params: Record<string, unknown> }[] = [];
let answer: unknown;
let refuse: AgentError | undefined;
const channel = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	if (refuse) throw refuse;
	return answer;
}) as Ask;

describe('/api/agent-stream/[story]', () => {
	let get: Handler;
	const doGet = (story: string, query = '') =>
		get({
			params: { story },
			url: new URL(`http://x/api/agent-stream/${story}${query}`)
		} as never);

	beforeAll(async () => {
		useRepo(new Repo('/nowhere', channel));
		get = (await import('./+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		answer = { story: 'S-0142', running: true, from: 0, next: 10, entries: [] };
		refuse = undefined;
	});
	afterAll(() => useRepo(null));

	it("asks flai for the tail, or from the offset given, and answers with flai's read", async () => {
		let r = await doGet('S-0142');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(answer as object);
		r = await doGet('S-0142', '?after=4096');
		expect(r.status).toBe(200);
		expect((await doGet('S-0142', '?after=-5')).status).toBe(200);
		// project.info is the answer's project identity, asked once
		expect(asked.filter((a) => a.method !== 'project.info')).toEqual([
			{ method: 'agent.stream', params: { story: 'S-0142' } },
			{ method: 'agent.stream', params: { story: 'S-0142', after: 4096 } },
			{ method: 'agent.stream', params: { story: 'S-0142', after: 0 } }
		]);
	});

	// S-0259: the planner page reads a planner run's stream through the same route
	it("asks flai for the planner's stream of the item when the query has plan", async () => {
		answer = { item: 'E-0016', running: false, from: 0, next: 10, entries: [] };
		let r = await doGet('E-0016', '?plan');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(answer as object);
		r = await doGet('E-0016', '?plan=1&after=4096');
		expect(r.status).toBe(200);
		expect(asked.filter((a) => a.method !== 'project.info')).toEqual([
			{ method: 'agent.stream', params: { plan: 'E-0016' } },
			{ method: 'agent.stream', params: { plan: 'E-0016', after: 4096 } }
		]);
		refuse = new AgentError(502, 'flai serve has started no planner for E-0016', -32004);
		r = await doGet('E-0016', '?plan');
		expect(r.status).toBe(404);
		expect((await r.json()).error).toContain('no planner');
	});

	// S-0218: the activity page reads the orchestrator's stream through the same route
	it("asks flai for the orchestrator's stream when the query has orchestrator", async () => {
		answer = { running: true, from: 0, next: 10, entries: [] };
		let r = await doGet('orchestrator', '?orchestrator');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(answer as object);
		r = await doGet('orchestrator', '?orchestrator&after=4096');
		expect(r.status).toBe(200);
		expect(asked.filter((a) => a.method !== 'project.info')).toEqual([
			{ method: 'agent.stream', params: { orchestrator: true } },
			{ method: 'agent.stream', params: { orchestrator: true, after: 4096 } }
		]);
	});

	// S-0228: the orchestrator's and the analyzer's pages read their runs' streams by role
	it("asks flai for the newest run of the role's stream when the query has a role", async () => {
		answer = { story: '', agent: 'analyzer', running: true, from: 0, next: 10, entries: [] };
		let r = await doGet('analyzer', '?role=analyze');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(answer as object);
		r = await doGet('orchestrator', '?role=orchestrate&after=4096');
		expect(r.status).toBe(200);
		r = await doGet('S-0142', '?role=orchestrate&plan&orchestrator');
		expect(r.status).toBe(200);
		expect(asked.filter((a) => a.method !== 'project.info')).toEqual([
			{ method: 'agent.stream', params: { role: 'analyze' } },
			{ method: 'agent.stream', params: { role: 'orchestrate', after: 4096 } },
			{ method: 'agent.stream', params: { role: 'orchestrate' } }
		]);
		refuse = new AgentError(502, 'flai serve has started no analyzer for the project', -32004);
		r = await doGet('analyzer', '?role=analyze');
		expect(r.status).toBe(404);
		expect((await r.json()).error).toContain('no analyzer');
	});

	// S-0420: a role page reads any role's newest run by the role's name, and flai judges the name
	it("asks flai for a role's newest run by its name, and passes on flai's refusal of one it does not run", async () => {
		answer = {
			story: '',
			item: 'S-0001',
			role: 'review',
			running: true,
			from: 0,
			next: 10,
			entries: []
		};
		let r = await doGet('reviewer', '?role=reviewer&after=128');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(answer as object);
		expect(asked.filter((a) => a.method === 'agent.stream')).toEqual([
			{ method: 'agent.stream', params: { role: 'reviewer', after: 128 } }
		]);
		refuse = new AgentError(
			502,
			'role must be orchestrate, for the orchestrator, ...; not "ghost"',
			-32602
		);
		r = await doGet('ghost', '?role=ghost');
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('not "ghost"');
	});

	it('refuses a role that is not a name in kebab case, asking flai nothing', async () => {
		for (const query of ['?role=', '?role=Reviewer', '?role=-x', '?role=a%20b', '?role=a/b']) {
			const r = await doGet('analyzer', query);
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('role must be orchestrate, analyze, or a role');
		}
		expect(asked.filter((a) => a.method === 'agent.stream')).toEqual([]);
	});

	it("passes on flai's refusals: no agent is a 404, a bad ID a 400", async () => {
		refuse = new AgentError(502, 'flai serve has started no agent for S-0142', -32004);
		let r = await doGet('S-0142');
		expect(r.status).toBe(404);
		expect((await r.json()).error).toContain('no agent');
		refuse = new AgentError(502, '"x" is not a story\'s ID', -32602);
		r = await doGet('x');
		expect(r.status).toBe(400);
	});
});
