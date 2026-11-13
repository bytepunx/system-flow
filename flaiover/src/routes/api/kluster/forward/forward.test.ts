import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type KlusterDeployments,
	type KlusterForward
} from '$lib/server/repo';

// /api/kluster/forward asks flai's kluster.forward.start and kluster.forward.stop (S-0586) of a
// stubbed host: a real forward needs a cluster and kubectl, which the flai side's tests fake too.
type Handler = (event: never) => Promise<Response>;

const ENABLE = 'flai serve enable kluster';

const forward = (over: Partial<KlusterForward> = {}): KlusterForward => ({
	namespace: 's-0586',
	service: 'web',
	port: 8080,
	local_port: 41234,
	address: '0.0.0.0',
	started_at: '2026-11-07T10:00:00Z',
	running: true,
	...over
});
const deployments = (forwards: KlusterForward[]): KlusterDeployments => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	cluster: 'authstar',
	default_namespace: 'authstar',
	deployments: [
		{
			story: 'S-0586',
			namespace: 's-0586',
			default: false,
			commit: '0123456789abcdef0123456789abcdef01234567',
			tag: 's-0586-0123456789ab',
			checkout: '.flai-cache/worktrees/S-0586',
			images: [],
			started_at: '2026-11-07T01:00:00Z',
			outcome: 'passed',
			steps: [],
			access: {
				kubeconfig: '.flai-cache/kluster/kubeconfig',
				context: 'k3d-authstar',
				namespace: 's-0586',
				list: '',
				services: [{ name: 'web', ports: [{ port: 8080 }] }],
				forwards
			}
		}
	],
	review: []
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
const WRITES = new Set(['kluster.forward.start', 'kluster.forward.stop']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return WRITES.has(method) ? { data: s.data, warnings: [] } : s.data;
}) as Ask;

describe('/api/kluster/forward (S-0586)', () => {
	let POST: Handler;
	const post = (body?: unknown) =>
		POST({
			request: new Request('http://x/api/kluster/forward', {
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
			'kluster.forward.start': { data: forward() },
			'kluster.forward.stop': { data: forward({ running: false }) },
			'kluster.deployments': { data: deployments([forward()]) }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('starts a forward, then answers the deployments with it', async () => {
		const r = await post({ namespace: 's-0586', service: 'web', port: 8080, on: true });
		expect(r.status).toBe(200);
		expect(asks('kluster.forward.start').map((a) => a.params)).toMatchObject([
			{ namespace: 's-0586', service: 'web', port: 8080 }
		]);
		const body = (await r.json()) as KlusterDeployments;
		expect(body.deployments[0].access.forwards).toEqual([forward()]);
	});

	it("starts and stops a forward into a story in progress's namespace as into one in review (S-0663)", async () => {
		const inProgress = (forwards: KlusterForward[]): KlusterDeployments => {
			const d = deployments(forwards);
			const [row] = d.deployments;
			return {
				...d,
				deployments: [
					{
						...row,
						story: 'S-0005',
						namespace: 's-0005',
						column: 'in-progress',
						status: 'in-progress',
						access: { ...row.access, namespace: 's-0005' }
					}
				]
			};
		};
		const on = forward({ namespace: 's-0005' });
		script['kluster.forward.start'] = { data: on };
		script['kluster.deployments'] = { data: inProgress([on]) };
		const started = await post({ namespace: 's-0005', service: 'web', port: 8080, on: true });
		expect(started.status).toBe(200);
		expect(asks('kluster.forward.start').map((a) => a.params)).toMatchObject([
			{ namespace: 's-0005', service: 'web', port: 8080 }
		]);
		const body = (await started.json()) as KlusterDeployments;
		expect(body.deployments[0]).toMatchObject({ column: 'in-progress', namespace: 's-0005' });
		expect(body.deployments[0].access.forwards).toEqual([on]);

		script['kluster.forward.stop'] = { data: { ...on, running: false } };
		script['kluster.deployments'] = { data: inProgress([]) };
		const stopped = await post({ namespace: 's-0005', service: 'web', port: 8080, on: false });
		expect(stopped.status).toBe(200);
		expect(asks('kluster.forward.stop').map((a) => a.params)).toMatchObject([
			{ namespace: 's-0005', service: 'web', port: 8080 }
		]);
		expect(((await stopped.json()) as KlusterDeployments).deployments[0].access.forwards).toEqual(
			[]
		);
	});

	it('stops a forward with on false', async () => {
		script['kluster.deployments'] = { data: deployments([]) };
		const r = await post({ namespace: 's-0586', service: 'web', port: 8080, on: false });
		expect(r.status).toBe(200);
		expect(asks('kluster.forward.stop')).toHaveLength(1);
		expect(asks('kluster.forward.start')).toHaveLength(0);
	});

	it('refuses a target that is not a namespace, a service, and a port before asking flai', async () => {
		for (const body of [
			{ namespace: 'S_0586', service: 'web', port: 8080 },
			{ namespace: 's-0586', service: '', port: 8080 },
			{ namespace: 's-0586', service: 'web', port: 0 },
			{ namespace: 's-0586', service: 'web', port: '8080' },
			undefined
		]) {
			const r = await post(body);
			expect(r.status).toBe(400);
		}
		expect(asks('kluster.forward.start')).toEqual([]);
	});

	it("answers 403 while the action is off, and flai's refusal word for word as 400", async () => {
		script['kluster.forward.start'] = {
			error: new AgentError(
				502,
				`the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`,
				-32012,
				{ action: 'kluster', enable: ENABLE }
			)
		};
		expect((await post({ namespace: 's-0586', service: 'web', port: 8080 })).status).toBe(403);
		const none =
			'no port-forward of service/web port 8080 in the namespace s-0586 is running; nothing to stop';
		script['kluster.forward.stop'] = { error: new AgentError(502, none, -32011) };
		const r = await post({ namespace: 's-0586', service: 'web', port: 8080, on: false });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(none);
	});

	it('answers 501 telling to upgrade a host flai without the methods', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		const r = await post({ namespace: 's-0586', service: 'web', port: 8080 });
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('the host flai is too old');
	});
});
