import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type KlusterDeployments,
	type KlusterJob
} from '$lib/server/repo';

// /api/kluster/checks asks flai's kluster.checks (S-0663) of a stubbed host: a real checks run needs
// a deployed namespace and the project's checks, which the flai side's tests fake too.
type Handler = (event: never) => Promise<Response>;

const ENABLE = 'flai serve enable kluster';
const COMMIT = '0123456789abcdef0123456789abcdef01234567';

const deployments = (over: Partial<KlusterDeployments> = {}): KlusterDeployments => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	cluster: 'authstar',
	default_namespace: 'authstar',
	deployments: [
		{
			story: 'S-0005',
			namespace: 's-0005',
			default: false,
			commit: COMMIT,
			tag: 's-0005-0123456789ab',
			checkout: '.flai-cache/worktrees/S-0005',
			images: [],
			started_at: '2026-11-13T01:00:00Z',
			ended_at: '2026-11-13T01:04:00Z',
			outcome: 'passed',
			steps: [],
			title: 'A story in progress',
			status: 'in-progress',
			column: 'in-progress',
			head: COMMIT,
			at_head: true,
			asked_by: 'agent-S-0005',
			access: {
				kubeconfig: '.flai-cache/kluster/kubeconfig',
				context: 'k3d-authstar',
				namespace: 's-0005',
				list: '',
				services: []
			}
		}
	],
	review: [],
	...over
});
const run: KlusterJob = {
	action: 'checks',
	cluster: 'authstar',
	story: 'S-0005',
	namespace: 's-0005',
	started_at: '2026-11-13T02:00:00Z',
	log: '.flai-cache/serve/checks/s-0005-20261113T020000Z.log'
};

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// the read answers raw, as flai's do; the write as commands do
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return method === 'kluster.checks' ? { data: s.data, warnings: [] } : s.data;
}) as Ask;

describe('/api/kluster/checks (S-0663)', () => {
	let POST: Handler;
	const post = (body?: unknown) =>
		POST({
			request: new Request('http://x/api/kluster/checks', {
				method: 'POST',
				...(body === undefined ? {} : { body: JSON.stringify(body) })
			})
		} as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		POST = (await import('./+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'kluster.checks': { data: run },
			'kluster.deployments': { data: deployments() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it("starts a story in progress's checks with its story and a request ID, then answers the deployments", async () => {
		const passed = deployments({
			deployments: [
				{
					...deployments().deployments[0],
					checks_commit: COMMIT,
					checks_outcome: 'passed',
					checks_log: run.log,
					checks_log_tail: 'ok\n'
				}
			]
		});
		script['kluster.deployments'] = { data: passed };
		const r = await post({ story: 'S-0005' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(passed);
		expect(asks('kluster.checks')).toHaveLength(1);
		const { request_id, ...rest } = asks('kluster.checks')[0].params;
		expect(typeof request_id).toBe('string');
		expect(rest).toEqual({ story: 'S-0005' });
		expect(asks('kluster.deployments')).toHaveLength(1);
	});

	it('refuses a body that names no story with 400, before flai is asked', async () => {
		for (const body of [undefined, {}, { story: 'main' }, { story: 5 }]) {
			const r = await post(body);
			expect(r.status).toBe(400);
			expect((await r.json()).error).toContain('story must be the ID');
		}
		expect(asks('kluster.checks')).toEqual([]);
	});

	it('answers 403 with the command that enables the action while it is off', async () => {
		script['kluster.checks'] = {
			error: new AgentError(
				502,
				`the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`,
				-32012,
				{ action: 'kluster', enable: ENABLE }
			)
		};
		const r = await post({ story: 'S-0005' });
		expect(r.status).toBe(403);
		const out = await r.json();
		expect(out).toMatchObject({ action: 'kluster', enable: ENABLE });
		expect(out.error).toContain(ENABLE);
		expect(asks('kluster.deployments')).toEqual([]);
	});

	it("answers flai's refusal by rule word for word as 400", async () => {
		const stale =
			'the namespace s-0005 is not deployed at the head of story/S-0005; deploy S-0005 first, then run its checks';
		script['kluster.checks'] = { error: new AgentError(502, stale, -32011) };
		const r = await post({ story: 'S-0005' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(stale);
		expect(asks('kluster.deployments')).toEqual([]);
	});

	it('answers 501 telling to upgrade a host flai without kluster.checks', async () => {
		delete script['kluster.checks'];
		const r = await post({ story: 'S-0005' });
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('the host flai is too old');
	});
});
