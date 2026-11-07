import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type KlusterJob,
	type KlusterStatus,
	type KlusterWorkloads
} from '$lib/server/repo';

// /api/kluster, /api/kluster/up, and /api/kluster/down ask flai's kluster.status, kluster.up, and
// kluster.down (S-0580) of a stubbed host: running a real cluster needs kluster and Docker, which
// the flai side's tests fake too.
type Handler = (event: never) => Promise<Response>;

const ENABLE = 'flai serve enable kluster';
const NO_BLOCK =
	"this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; add one";

const job = (action: 'up' | 'down', ended = false): KlusterJob => ({
	action,
	cluster: 'authstar',
	started_at: '2026-11-07T01:00:00Z',
	log: `.flai-cache/kluster/${action}.log`,
	...(ended ? { ended_at: '2026-11-07T01:03:00Z', outcome: 'done' as const } : {})
});

const status = (over: Partial<KlusterStatus> = {}): KlusterStatus => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	settings: [],
	settings_on: true,
	settings_enable: 'flai serve enable settings',
	findings: [],
	deployable: true,
	tools: [],
	cluster: {
		name: 'authstar',
		profile: 'integration',
		provider: 'k3d',
		context: 'k3d-authstar',
		exists: true,
		running: true,
		age: '2h50m39s',
		age_seconds: 10239
	},
	...over
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// the read answers raw, as flai's do; the two writes as commands do
const WRITES = new Set(['kluster.up', 'kluster.down']);
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

describe('/api/kluster (S-0580)', () => {
	let GET: Handler;
	let UP: Handler;
	let DOWN: Handler;
	const get = () => GET({ url: new URL('http://x/api/kluster') } as never);
	const post = (fn: Handler) =>
		fn({ request: new Request('http://x/api/kluster/x', { method: 'POST' }) } as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
		UP = (await import('./up/+server')).POST as unknown as Handler;
		DOWN = (await import('./down/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'kluster.status': { data: status() }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('answers the status, asking flai afresh each time', async () => {
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(status());
		script['kluster.status'] = {
			data: status({ cluster: { ...status().cluster!, running: false, exists: false } })
		};
		const again = (await (await get()).json()) as KlusterStatus;
		expect(again.cluster?.running).toBe(false);
		expect(asks('kluster.status')).toHaveLength(2);
	});

	it('answers the status with the action off and with no block, as flai gives them', async () => {
		script['kluster.status'] = {
			data: { enabled: false, enable: ENABLE, block: false, no_block: NO_BLOCK }
		};
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({
			enabled: false,
			enable: ENABLE,
			block: false,
			no_block: NO_BLOCK
		});
	});

	for (const [name, action, fn] of [
		['up', 'kluster.up', () => UP],
		['down', 'kluster.down', () => DOWN]
	] as const) {
		it(`starts kluster ${name} with a request ID, then answers the status with the job running`, async () => {
			script[action] = { data: job(name) };
			script['kluster.status'] = { data: status({ running: job(name) }) };
			const r = await post(fn());
			expect(r.status).toBe(200);
			expect(await r.json()).toMatchObject(status({ running: job(name) }));
			expect(asks(action)).toHaveLength(1);
			expect(typeof asks(action)[0].params.request_id).toBe('string');
			expect(asks('kluster.status')).toHaveLength(1);
		});
	}

	it('answers 403 with the command that enables the action while it is off', async () => {
		script['kluster.up'] = { error: disabled() };
		script['kluster.down'] = { error: disabled() };
		for (const fn of [UP, DOWN]) {
			const r = await post(fn);
			expect(r.status).toBe(403);
			const body = await r.json();
			expect(body).toMatchObject({ action: 'kluster', enable: ENABLE });
			expect(body.error).toContain(ENABLE);
		}
		expect(asks('kluster.status')).toEqual([]);
	});

	it("answers flai's refusal by rule word for word as 400: no block, and a job that runs", async () => {
		const running =
			'kluster up of the cluster authstar runs since 2026-11-07T01:00:00Z; wait for it to end, as kluster.status says, then ask again';
		script['kluster.up'] = { error: new AgentError(502, NO_BLOCK, -32011) };
		script['kluster.down'] = { error: new AgentError(502, running, -32011) };
		const up = await post(UP);
		expect(up.status).toBe(400);
		expect((await up.json()).error).toBe(NO_BLOCK);
		const down = await post(DOWN);
		expect(down.status).toBe(400);
		expect((await down.json()).error).toBe(running);
	});

	it('answers 501 telling to upgrade a host flai without the three methods', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		for (const r of [await get(), await post(UP), await post(DOWN)]) {
			expect(r.status).toBe(501);
			const error = (await r.json()).error as string;
			expect(error).toContain('the host flai is too old');
			expect(error).toContain('flai self-upgrade');
		}
	});
});

describe('/api/kluster setup and /api/kluster/workloads (S-0586)', () => {
	let GET: Handler;
	let WORKLOADS: Handler;
	const get = (fn: Handler, path: string) => fn({ url: new URL(`http://x${path}`) } as never);

	const scaffold = status({
		deployable: false,
		settings_on: false,
		config: 'kluster.yaml',
		config_file: {
			exists: true,
			name: 'p-local',
			provider: 'k3d',
			profile: 'spire',
			default_profile: 'spire',
			profiles: ['spire', 'signet', 'authstar'],
			addons: []
		},
		settings: [
			{
				key: 'kluster.config',
				kind: 'text',
				meaning: 'The config.',
				value: 'kluster.yaml',
				set: true
			},
			{ key: 'kluster.deploy', kind: 'command', meaning: 'The deploy.', set: false },
			{ key: 'kluster.timeout', kind: 'duration', default: '30m', meaning: 'The most.', set: false }
		],
		tools: [
			{
				name: 'kluster',
				found: true,
				path: '/usr/local/bin/kluster',
				need: 'up',
				docs: 'https://x'
			},
			{
				name: 'docker',
				found: false,
				need: 'nodes',
				docs: 'https://docs.docker.com/engine/install/'
			},
			{ name: 'kubectl', found: false, need: 'pods', install: 'kluster setup', docs: 'https://y' }
		]
	});

	const workloads: KlusterWorkloads = {
		block: true,
		cluster: 'p-local',
		context: 'k3d-p-local',
		nodes: [
			{ name: 'k3d-p-local-server-0', ready: true, roles: ['control-plane'], version: 'v1.31' }
		],
		namespaces: [
			{
				namespace: 'kube-system',
				pods: [{ name: 'coredns-1', phase: 'Running', ready: 1, containers: 1, restarts: 0 }]
			}
		]
	};

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
		WORKLOADS = (await import('./workloads/+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'kluster.status': { data: scaffold },
			'kluster.workloads': { data: workloads }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it("answers what setup still lacks as flai gives it: the block's keys, its config, the tools", async () => {
		const r = await get(GET, '/api/kluster');
		expect(r.status).toBe(200);
		const body = (await r.json()) as KlusterStatus;
		expect(body).toMatchObject({ deployable: false, settings_on: false, config: 'kluster.yaml' });
		expect(body.config_file?.profile).toBe('spire');
		expect(body.settings?.map((s) => s.key)).toEqual([
			'kluster.config',
			'kluster.deploy',
			'kluster.timeout'
		]);
		expect(body.tools?.filter((t) => !t.found).map((t) => t.name)).toEqual(['docker', 'kubectl']);
	});

	it('answers the nodes and pods, asking flai afresh each time', async () => {
		const r = await get(WORKLOADS, '/api/kluster/workloads');
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(workloads);
		script['kluster.workloads'] = {
			data: { ...workloads, nodes: [], namespaces: [], workloads_error: 'no kubeconfig yet' }
		};
		const again = (await (
			await get(WORKLOADS, '/api/kluster/workloads')
		).json()) as KlusterWorkloads;
		expect(again.workloads_error).toBe('no kubeconfig yet');
		expect(asked.filter((a) => a.method === 'kluster.workloads')).toHaveLength(2);
	});

	it('answers 501 telling to upgrade a host flai without kluster.workloads', async () => {
		delete script['kluster.workloads'];
		const r = await get(WORKLOADS, '/api/kluster/workloads');
		expect(r.status).toBe(501);
		const error = (await r.json()).error as string;
		expect(error).toContain('the host flai is too old');
		expect(error).toContain('flai self-upgrade');
	});
});
