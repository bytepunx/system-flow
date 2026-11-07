// S-0581: the Deployments page lists main's namespace first, then each story's, then the stories in
// review with no deployment, each with its story and title linking its item page, its namespace,
// commit, when it was deployed, its outcome, each step's log, and how to reach it for a demo; it
// deploys a story only where that fits, tears one down once confirmed, and promotes main, through
// /api/kluster/*; says while a job runs and asks again until it ends; shows flai's refusal word for
// word; and says when the kluster host action is off or nothing is deployed.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { KlusterDeployment, KlusterDeployments, KlusterJob } from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params?: Record<string, string>) =>
		params ? route.replace(/\[(?:\.\.\.)?(\w+)\]/g, (_, k: string) => params[k]) : route
}));

import DeploymentsPage from './+page.svelte';

const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

const ENABLE = 'flai serve enable kluster';
const HEAD = '0123456789abcdef0123456789abcdef01234567';
const OLD = 'fedcba9876543210fedcba9876543210fedcba98';

const access = (ns: string) => ({
	kubeconfig: '.flai-cache/kluster/kubeconfig',
	context: 'k3d-authstar',
	namespace: ns,
	list: `kubectl --kubeconfig .flai-cache/kluster/kubeconfig --context k3d-authstar --namespace ${ns} get services`,
	services: [
		{
			name: 'web',
			host: `web.${ns}.svc.cluster.local`,
			ports: [{ name: 'http', port: 8080, protocol: 'TCP' }],
			port_forward: `kubectl --kubeconfig .flai-cache/kluster/kubeconfig --context k3d-authstar --namespace ${ns} port-forward service/web 8080:8080`
		}
	]
});
const story = (id: string, over: Partial<KlusterDeployment> = {}): KlusterDeployment => {
	const ns = id.toLowerCase();
	return {
		story: id,
		namespace: ns,
		default: false,
		commit: HEAD,
		tag: `${ns}-0123456789ab`,
		checkout: `.flai-cache/worktrees/${id}`,
		images: [],
		started_at: '2026-11-07T01:00:00Z',
		ended_at: '2026-11-07T01:04:00Z',
		outcome: 'passed',
		steps: [
			{
				name: 'build',
				started_at: '2026-11-07T01:00:00Z',
				ended_at: '2026-11-07T01:02:00Z',
				outcome: 'passed',
				log: `.flai-cache/kluster/logs/${ns}/build.log`,
				log_tail: `built ${ns}\n`
			},
			{
				name: 'deploy',
				started_at: '2026-11-07T01:02:00Z',
				ended_at: '2026-11-07T01:04:00Z',
				outcome: 'passed',
				log: `.flai-cache/kluster/logs/${ns}/deploy.log`,
				log_tail: `applied ${ns}\n`
			}
		],
		title: `Story ${id}`,
		status: 'review',
		head: HEAD,
		access: access(ns),
		...over
	};
};
const main = (): KlusterDeployment => ({
	...story('S-0000', {}),
	story: '',
	namespace: 'authstar',
	default: true,
	tag: 'main-0123456789ab',
	title: undefined,
	status: undefined,
	access: access('authstar')
});
const list = (over: Partial<KlusterDeployments> = {}): KlusterDeployments => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	cluster: 'authstar',
	default_namespace: 'authstar',
	deployments: [main(), story('S-0581'), story('S-0583', { head: OLD })],
	review: [{ story: 'S-0584', title: 'Not yet deployed', head: HEAD }],
	...over
});
const job = (
	action: KlusterJob['action'],
	id: string | undefined,
	over: Partial<KlusterJob> = {}
) =>
	({
		action,
		cluster: 'authstar',
		...(id ? { story: id, namespace: id.toLowerCase() } : { namespace: 'authstar' }),
		started_at: '2026-11-07T02:00:00Z',
		log: `.flai-cache/kluster/${action}.log`,
		...over
	}) as KlusterJob;

const all = (id: string) => [...document.querySelectorAll(`[data-testid="${id}"]`)];
const q = (id: string) => document.querySelector(`[data-testid="${id}"]`);
const rowOf = (who: string) =>
	all('deployment-row').find((r) =>
		r.querySelector('[data-testid="deployment-who"]')!.textContent!.includes(who)
	)!;
const buttonIn = (el: Element | Document, name: string) =>
	[...el.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent!.trim() === name);
const posts = () =>
	api.mock.calls.filter(([, init]) => (init as RequestInit | undefined)?.method === 'POST');

describe('Deployments page (S-0581)', () => {
	let c: ReturnType<typeof mount> | undefined;
	// the answers GET /api/kluster/deployments gives, in turn; the last is kept
	let reads: ReturnType<typeof answer>[] = [];
	let writes: Record<string, ReturnType<typeof answer>> = {};

	const show = async () => {
		c = mount(DeploymentsPage, { target: document.body });
		await settle();
	};

	beforeEach(() => {
		reads = [answer(list())];
		writes = {};
		api.mockReset();
		api.mockImplementation(async (path: string) => {
			if (path === '/api/kluster/deployments') return reads.length > 1 ? reads.shift()! : reads[0];
			const w = writes[path];
			if (!w) throw new Error(`unexpected ${path}`);
			return w;
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.useRealTimers();
	});

	it("lists main first, then each story's namespace, then the stories in review not deployed", async () => {
		await show();
		const who = all('deployment-who').map((w) => w.textContent!.replace(/\s+/g, ' ').trim());
		expect(who).toEqual([
			'main',
			'S-0581 Story S-0581',
			'S-0583 Story S-0583',
			'S-0584 Not yet deployed'
		]);
		const link = rowOf('S-0581').querySelector('a')!;
		expect(link.getAttribute('href')).toBe('/items/S-0581');
		const text = rowOf('S-0581').textContent!;
		expect(text).toContain('s-0581');
		expect(rowOf('S-0581').querySelector('[data-testid="deployment-commit"]')!.textContent).toBe(
			'0123456789ab'
		);
		expect(
			rowOf('S-0581').querySelector('[data-testid="deployment-outcome"]')!.textContent
		).toContain('passed');
		expect(text).toMatch(/\d{4}-\d\d-\d\d \d\d:\d\d/);
		expect(
			rowOf('S-0584').querySelector('[data-testid="deployment-outcome"]')!.textContent
		).toContain('not deployed');
		expect(q('deployments-empty')).toBeNull();
		expect(q('deployments-no-main')).toBeNull();
	});

	it("links each step to its log, with the log's path and its end", async () => {
		await show();
		const steps = [...rowOf('S-0581').querySelectorAll('[data-testid="deployment-step"]')];
		expect(steps.map((s) => s.id)).toEqual(['log-s-0581-build', 'log-s-0581-deploy']);
		expect(steps[0].querySelector('summary')!.textContent).toContain('its log');
		expect(steps[0].textContent).toContain('.flai-cache/kluster/logs/s-0581/build.log');
		expect(steps[0].querySelector('[data-testid="deployment-step-log"]')!.textContent).toBe(
			'built s-0581\n'
		);
	});

	it('shows how to reach each deployment for a demo: its hosts, port-forwards, and services', async () => {
		reads = [
			answer(
				list({
					deployments: [
						main(),
						story('S-0581', {
							access: {
								...access('s-0581'),
								services: [],
								services_error: "kubectl is not on the host's PATH"
							}
						})
					]
				})
			)
		];
		await show();
		const demo = rowOf('main').querySelector('[data-testid="deployment-access"]')!;
		expect(demo.textContent).toContain('web.authstar.svc.cluster.local');
		expect(demo.textContent).toContain('8080');
		expect(demo.querySelector('[data-testid="deployment-port-forward"]')!.textContent).toBe(
			access('authstar').services[0].port_forward
		);
		expect(demo.querySelector('[data-testid="deployment-list"]')!.textContent).toBe(
			access('authstar').list
		);
		expect(demo.textContent).toContain('main checkout');
		expect(
			rowOf('S-0581').querySelector('[data-testid="deployment-services-error"]')!.textContent
		).toContain("kubectl is not on the host's PATH");
	});

	it('offers Deploy only where it fits: not deployed, or deployed behind its head', async () => {
		await show();
		expect(buttonIn(rowOf('S-0581'), 'Deploy')).toBeUndefined();
		expect(buttonIn(rowOf('S-0583'), 'Deploy')).toBeDefined();
		expect(rowOf('S-0583').querySelector('[data-testid="deployment-behind"]')).not.toBeNull();
		expect(buttonIn(rowOf('S-0584'), 'Deploy')).toBeDefined();
		expect(buttonIn(rowOf('main'), 'Deploy')).toBeUndefined();
		expect(buttonIn(rowOf('main'), 'Teardown')).toBeUndefined();
		expect(buttonIn(rowOf('main'), 'Promote')).toBeDefined();
		expect(buttonIn(rowOf('S-0581'), 'Teardown')).toBeDefined();
		expect(buttonIn(rowOf('S-0584'), 'Teardown')).toBeUndefined();
	});

	it('deploys a story, posting its ID, and shows the deploy running, the end of its log, and asks again until it ends', async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		await show();
		const running = job('deploy', 'S-0584', { log_tail: 'building s-0584\n' });
		writes['/api/kluster/deploy'] = answer(list({ running }));
		buttonIn(rowOf('S-0584'), 'Deploy')!.click();
		await settle();
		expect(posts().map(([p, init]) => [p, (init as RequestInit).body])).toEqual([
			['/api/kluster/deploy', JSON.stringify({ story: 'S-0584' })]
		]);
		expect(q('deployments-running')!.textContent).toContain(
			'kluster deploy of S-0584 in the namespace s-0584 runs'
		);
		expect(q('deployments-log')!.textContent).toBe('building s-0584\n');
		expect(buttonIn(rowOf('S-0583'), 'Deploy')!.disabled).toBe(true);
		expect(buttonIn(rowOf('main'), 'Promote')!.disabled).toBe(true);

		const done = job('deploy', 'S-0584', { ended_at: '2026-11-07T02:04:00Z', outcome: 'done' });
		const after = list({
			deployments: [main(), story('S-0581'), story('S-0583', { head: OLD }), story('S-0584')],
			review: [],
			last: done
		});
		reads = [answer(list({ running })), answer(after)];
		const asked = () => api.mock.calls.filter(([p]) => p === '/api/kluster/deployments').length;
		const before = asked();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 1);
		expect(q('deployments-running')).not.toBeNull();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 2);
		expect(q('deployments-running')).toBeNull();
		expect(q('deployments-last')!.textContent).toContain('done');
		expect(
			rowOf('S-0584').querySelector('[data-testid="deployment-outcome"]')!.textContent
		).toContain('passed');
		// with nothing running, the page does not ask again
		vi.advanceTimersByTime(30000);
		await settle();
		expect(asked()).toBe(before + 2);
	});

	it("shows a deploy another process runs, from the lock's job, and asks again while it runs", async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		const locked = story('S-0581', {
			outcome: 'running',
			ended_at: undefined,
			job: {
				kind: 'deploy',
				story: 'S-0581',
				namespace: 's-0581',
				started_at: '2026-11-07T02:00:00Z',
				pid: 42
			}
		});
		reads = [answer(list({ deployments: [main(), locked] })), answer(list())];
		await show();
		expect(q('deployments-running')!.textContent).toContain('kluster deploy of S-0581');
		expect(rowOf('S-0581').querySelector('[data-testid="deployment-job"]')).not.toBeNull();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(q('deployments-running')).toBeNull();
	});

	it('tears a namespace down only once the operator confirms', async () => {
		await show();
		buttonIn(rowOf('S-0581'), 'Teardown')!.click();
		flushSync();
		expect(posts()).toEqual([]);
		const confirm = q('deployment-confirm-teardown')!;
		expect(confirm.textContent).toContain('s-0581');
		buttonIn(confirm, 'Keep it')!.click();
		flushSync();
		expect(q('deployment-confirm-teardown')).toBeNull();
		expect(posts()).toEqual([]);

		writes['/api/kluster/teardown'] = answer(list({ running: job('teardown', 'S-0581') }));
		buttonIn(rowOf('S-0581'), 'Teardown')!.click();
		flushSync();
		buttonIn(q('deployment-confirm-teardown')!, 'Tear it down')!.click();
		await settle();
		expect(posts().map(([p, init]) => [p, (init as RequestInit).body])).toEqual([
			['/api/kluster/teardown', JSON.stringify({ story: 'S-0581' })]
		]);
		expect(q('deployments-running')!.textContent).toContain('kluster teardown of S-0581');
	});

	it('promotes main', async () => {
		await show();
		writes['/api/kluster/promote'] = answer(list({ running: job('promote', undefined) }));
		buttonIn(rowOf('main'), 'Promote')!.click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/kluster/promote']);
		expect(q('deployments-running')!.textContent).toContain(
			'kluster promote of main into the namespace authstar runs'
		);
	});

	it("shows flai's refusal word for word", async () => {
		const running =
			'kluster deploy of S-0582 in the namespace s-0582 runs since 2026-11-07T01:00:00Z; wait for it to end, as kluster.deployments says, then ask again';
		writes['/api/kluster/deploy'] = answer({ error: running }, 400);
		await show();
		buttonIn(rowOf('S-0584'), 'Deploy')!.click();
		await settle();
		expect(q('deployments-action-error')!.textContent!.trim()).toBe(running);
	});

	it('says the host action is off, linking the Features page, and disables every action', async () => {
		reads = [answer(list({ enabled: false }))];
		await show();
		const off = q('deployments-off')!;
		expect(off.textContent).toContain(ENABLE);
		expect(off.querySelector('a')!.getAttribute('href')).toBe('/settings/features');
		for (const b of document.querySelectorAll<HTMLButtonElement>(
			'[data-testid="deployment-row"] button'
		))
			expect(b.disabled).toBe(true);
	});

	it('takes a 403 on a write as the host action off, with its command', async () => {
		await show();
		writes['/api/kluster/promote'] = answer(
			{
				error: `the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`,
				action: 'kluster',
				enable: ENABLE
			},
			403
		);
		buttonIn(rowOf('main'), 'Promote')!.click();
		await settle();
		expect(q('deployments-off')!.textContent).toContain(ENABLE);
		expect(q('deployments-action-error')!.textContent).toContain('not enabled');
	});

	it('says when the project has no kluster: block, linking the design', async () => {
		reads = [
			answer({
				...list({ deployments: [], review: [] }),
				block: false,
				no_block: 'add a kluster: block'
			})
		];
		await show();
		const nb = q('deployments-no-block')!;
		expect(nb.textContent).toContain('add a kluster: block');
		expect(nb.querySelector('a')!.getAttribute('href')).toBe('/docs/design/system/kluster.md');
		expect(all('deployment-row')).toEqual([]);
	});

	it('says an empty cluster has nothing deployed, and offers Promote', async () => {
		reads = [answer(list({ deployments: [], review: [] }))];
		await show();
		expect(q('deployments-empty')!.textContent).toContain('Nothing is deployed');
		expect(all('deployment-row')).toEqual([]);
		writes['/api/kluster/promote'] = answer(list({ running: job('promote', undefined) }));
		buttonIn(q('deployments-no-main')!, 'Promote')!.click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/kluster/promote']);
	});

	it('says to upgrade a host flai older than kluster.deployments', async () => {
		reads = [answer({ error: 'the host flai is too old to show the deployments' }, 501)];
		await show();
		expect(q('deployments-too-old')!.textContent).toContain('too old');
	});
});
