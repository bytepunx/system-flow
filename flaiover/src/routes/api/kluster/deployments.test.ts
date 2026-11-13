import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type KlusterDeployment,
	type KlusterDeployments,
	type KlusterJob
} from '$lib/server/repo';

// /api/kluster/deployments, /deploy, /teardown, and /promote ask flai's kluster.deployments,
// kluster.deploy, kluster.teardown, and kluster.promote (S-0581) of a stubbed host: a real deploy
// needs kluster, Docker, and the project's own scripts, which the flai side's tests fake too.
type Handler = (event: never) => Promise<Response>;

const ENABLE = 'flai serve enable kluster';

const access = (ns: string) => ({
	kubeconfig: '.flai-cache/kluster/kubeconfig',
	context: 'k3d-authstar',
	namespace: ns,
	list: `kubectl --kubeconfig .flai-cache/kluster/kubeconfig --context k3d-authstar --namespace ${ns} get services`,
	services: []
});
const deployment = (over: Partial<KlusterDeployment> = {}): KlusterDeployment => ({
	story: 'S-0581',
	namespace: 's-0581',
	default: false,
	commit: '0123456789abcdef0123456789abcdef01234567',
	tag: 's-0581-0123456789ab',
	checkout: '.flai-cache/worktrees/S-0581',
	images: ['authstar/web:s-0581-0123456789ab'],
	started_at: '2026-11-07T01:00:00Z',
	ended_at: '2026-11-07T01:04:00Z',
	outcome: 'passed',
	steps: [],
	title: 'The Deployments page',
	status: 'review',
	head: '0123456789abcdef0123456789abcdef01234567',
	access: access('s-0581'),
	...over
});
const deployments = (over: Partial<KlusterDeployments> = {}): KlusterDeployments => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	cluster: 'authstar',
	default_namespace: 'authstar',
	deployments: [deployment()],
	review: [],
	...over
});
const job = (action: KlusterJob['action'], story?: string): KlusterJob => ({
	action,
	cluster: 'authstar',
	...(story ? { story, namespace: story.toLowerCase() } : { namespace: 'authstar' }),
	started_at: '2026-11-07T02:00:00Z',
	log: `.flai-cache/kluster/${action}.log`,
	logs: `.flai-cache/kluster/logs/${story ? story.toLowerCase() : 'authstar'}`
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// the read answers raw, as flai's do; the three writes as commands do
const WRITES = new Set(['kluster.deploy', 'kluster.teardown', 'kluster.promote']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return WRITES.has(method) ? { data: s.data, warnings: [] } : s.data;
}) as Ask;

const disabled = () =>
	new AgentError(
		502,
		`the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`,
		-32012,
		{ action: 'kluster', enable: ENABLE }
	);

describe('/api/kluster/deployments (S-0581)', () => {
	let GET: Handler;
	let DEPLOY: Handler;
	let TEARDOWN: Handler;
	let PROMOTE: Handler;
	const get = () => GET({ url: new URL('http://x/api/kluster/deployments') } as never);
	const post = (fn: Handler, body?: unknown) =>
		fn({
			request: new Request('http://x/api/kluster/x', {
				method: 'POST',
				...(body === undefined ? {} : { body: JSON.stringify(body) })
			})
		} as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		GET = (await import('./deployments/+server')).GET as unknown as Handler;
		DEPLOY = (await import('./deploy/+server')).POST as unknown as Handler;
		TEARDOWN = (await import('./teardown/+server')).POST as unknown as Handler;
		PROMOTE = (await import('./promote/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'kluster.deployments': { data: deployments() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('answers the deployments as flai lists them, asking afresh each time', async () => {
		const main = deployment({
			story: '',
			namespace: 'authstar',
			default: true,
			title: undefined,
			status: undefined,
			access: access('authstar')
		});
		script['kluster.deployments'] = {
			data: deployments({
				deployments: [main, deployment()],
				review: [{ story: 'S-0582', title: 'Another', head: 'abc' }]
			})
		};
		const r = await get();
		expect(r.status).toBe(200);
		const body = (await r.json()) as KlusterDeployments;
		expect(body.deployments.map((d) => d.namespace)).toEqual(['authstar', 's-0581']);
		expect(body.review).toEqual([{ story: 'S-0582', title: 'Another', head: 'abc' }]);
		await get();
		expect(asks('kluster.deployments')).toHaveLength(2);
	});

	it("passes a story in progress's row and the count against the limit through (S-0663)", async () => {
		const head = 'fedcba9876543210fedcba9876543210fedcba98';
		const progress = deployment({
			story: 'S-0005',
			namespace: 's-0005',
			title: 'A story in progress',
			status: 'in-progress',
			column: 'in-progress',
			head,
			at_head: false,
			asked_by: 'agent-S-0005',
			log: '.flai-cache/kluster/logs/s-0005',
			checks_commit: '0123456789abcdef0123456789abcdef01234567',
			checks_outcome: 'failed',
			checks_log: '.flai-cache/serve/checks/s-0005-20261113T020000Z.log',
			checks_asked_by: 'agent-S-0005',
			checks_log_tail: 'FAIL login\n',
			access: access('s-0005')
		});
		const counted = deployments({
			deployments: [progress, deployment({ column: 'review', at_head: true })],
			review: [{ story: 'S-0582', title: 'Another' }],
			undeployed: [
				{ story: 'S-0007', title: 'Not yet deployed', head, column: 'in-progress' },
				{ story: 'S-0582', title: 'Another', column: 'review' }
			],
			live: 1,
			holding: ['S-0005'],
			limit: 1,
			refusals: [
				{
					story: 'S-0007',
					by: 'agent-S-0007',
					at: '2026-11-13T03:00:00Z',
					limit: 1,
					reason: 'S-0005 holds the one story namespace kluster.agent_namespaces allows'
				}
			]
		});
		script['kluster.deployments'] = { data: counted };
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(counted);
	});

	it('answers an empty cluster, and a running deploy, as flai gives them', async () => {
		script['kluster.deployments'] = { data: deployments({ deployments: [], review: [] }) };
		expect(await (await get()).json()).toMatchObject({ deployments: [], review: [] });
		script['kluster.deployments'] = {
			data: deployments({ running: { ...job('deploy', 'S-0581'), log_tail: 'building\n' } })
		};
		expect(((await (await get()).json()) as KlusterDeployments).running).toMatchObject({
			action: 'deploy',
			story: 'S-0581',
			log_tail: 'building\n'
		});
	});

	for (const [name, method, fn, body, params] of [
		['deploy', 'kluster.deploy', () => DEPLOY, { story: 'S-0581' }, { story: 'S-0581' }],
		['teardown', 'kluster.teardown', () => TEARDOWN, { story: 'S-0581' }, { story: 'S-0581' }],
		['promote', 'kluster.promote', () => PROMOTE, undefined, {}]
	] as const) {
		it(`starts ${name} with its params and a request ID, then answers the deployments`, async () => {
			const started = job(name, name === 'promote' ? undefined : 'S-0581');
			script[method] = { data: started };
			script['kluster.deployments'] = { data: deployments({ running: started }) };
			const r = await post(fn(), body);
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject(deployments({ running: started }));
			expect(asks(method)).toHaveLength(1);
			const { request_id, ...rest } = asks(method)[0].params;
			expect(typeof request_id).toBe('string');
			expect(rest).toEqual(params);
			expect(asks('kluster.deployments')).toHaveLength(1);
		});
	}

	it('refuses a deploy or a teardown that names no story with 400, before flai is asked', async () => {
		for (const fn of [DEPLOY, TEARDOWN]) {
			for (const body of [undefined, {}, { story: 'main' }, { story: 581 }]) {
				const r = await post(fn, body);
				expect(r.status).toBe(400);
				expect((await r.json()).error).toContain('story must be the ID');
			}
		}
		expect(asks('kluster.deploy')).toEqual([]);
		expect(asks('kluster.teardown')).toEqual([]);
	});

	it('answers 403 with the command that enables the action while it is off', async () => {
		script['kluster.deploy'] = { error: disabled() };
		script['kluster.teardown'] = { error: disabled() };
		script['kluster.promote'] = { error: disabled() };
		for (const [fn, body] of [
			[DEPLOY, { story: 'S-0581' }],
			[TEARDOWN, { story: 'S-0581' }],
			[PROMOTE, undefined]
		] as const) {
			const r = await post(fn, body);
			expect(r.status).toBe(403);
			const out = await r.json();
			expect(out).toMatchObject({ action: 'kluster', enable: ENABLE });
			expect(out.error).toContain(ENABLE);
		}
		expect(asks('kluster.deployments')).toEqual([]);
	});

	it("answers flai's refusal by rule word for word as 400", async () => {
		const running =
			'kluster deploy of S-0582 in the namespace s-0582 runs since 2026-11-07T01:00:00Z; wait for it to end, as kluster.deployments says, then ask again';
		const none = 'S-0581 has no deployment recorded; nothing to tear down';
		script['kluster.deploy'] = { error: new AgentError(502, running, -32011) };
		script['kluster.teardown'] = { error: new AgentError(502, none, -32011) };
		const d = await post(DEPLOY, { story: 'S-0581' });
		expect(d.status).toBe(400);
		expect((await d.json()).error).toBe(running);
		const t = await post(TEARDOWN, { story: 'S-0581' });
		expect(t.status).toBe(400);
		expect((await t.json()).error).toBe(none);
	});

	it('answers 501 telling to upgrade a host flai without the four methods', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		for (const r of [
			await get(),
			await post(DEPLOY, { story: 'S-0581' }),
			await post(TEARDOWN, { story: 'S-0581' }),
			await post(PROMOTE)
		]) {
			expect(r.status).toBe(501);
			const error = (await r.json()).error as string;
			expect(error).toContain('the host flai is too old');
		}
	});
});
