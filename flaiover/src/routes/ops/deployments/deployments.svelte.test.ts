// S-0581: the Deployments page lists main's namespace first, then each story's, then the stories in
// review with no deployment, each with its story and title linking its item page, its namespace,
// commit, when it was deployed, its outcome, each step's log, and how to reach it for a demo; it
// deploys a story only where that fits, tears one down once confirmed, and promotes main, through
// /api/kluster/*; says while a job runs and asks again until it ends; shows flai's refusal word for
// word; and says when the kluster host action is off or nothing is deployed. S-0586: the notices
// link the Klusters page's setup and name no command, and with no deploy set the page says deploys
// are not set up yet in place of Deploy and Promote. S-0595 (ADR-0263 §5): the page groups them in
// its own rail, Overview, Build, Logs, Access, and Setup, kept in the hash; Build shows one panel per
// component build with its status; with no deploy set, a link opens Setup in place of a call flai
// refuses; and Setup holds the deploy settings, so no one edits system-flow.yaml. S-0663: the
// stories in progress are listed beside those in review, each story namespace through DeploymentRow
// with its column, head, last deploy and checks run, and who asked; the page deploys, runs the
// checks of, and tears down a story in progress through /api/kluster/deploy, /checks, and
// /teardown; it counts the story namespaces against the limit and names the request refused for
// it; and a port-forward into a story in progress's namespace starts and stops as one in review's.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { KlusterDeployment, KlusterDeployments, KlusterJob } from '$lib/server/repo';
import { localTime } from '$lib/localtime';

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
// the rail's entry for a pane, and the pane shown
const entry = (id: string) =>
	document.querySelector<HTMLAnchorElement>(`[data-sidebar-entry="${id}"]`)!;
const paneShown = () =>
	document.querySelector('[data-sidebar-pane]')?.getAttribute('data-sidebar-pane');
const accessOf = (who: string) =>
	all('access-row').find((r) =>
		r.querySelector('[data-testid="deployment-who"]')!.textContent!.includes(who)
	)!;
const step = (name: string, outcome: string, over: Record<string, unknown> = {}) => ({
	name,
	outcome,
	...(outcome === 'pending'
		? {}
		: { started_at: '2026-11-07T01:00:00Z', log: `.flai-cache/kluster/logs/s-0590/${name}.log` }),
	...(outcome === 'passed' || outcome === 'failed' || outcome === 'skipped'
		? { ended_at: '2026-11-07T01:01:30Z' }
		: {}),
	...over
});

describe('Deployments page (S-0581)', () => {
	let c: ReturnType<typeof mount> | undefined;
	// the answers GET /api/kluster/deployments gives, in turn; the last is kept
	let reads: ReturnType<typeof answer>[] = [];
	let writes: Record<string, ReturnType<typeof answer>> = {};
	// what GET /api/kluster, kluster.status, answers: only block and deployable are read here
	let status: ReturnType<typeof answer> = answer({ block: true, deployable: true });

	const show = async () => {
		c = mount(DeploymentsPage, { target: document.body });
		await settle();
	};
	// opens a pane from the rail, as the reader does
	const go = async (id: string) => {
		entry(id).click();
		await settle();
	};

	beforeEach(() => {
		history.replaceState(null, '', '/ops/deployments');
		reads = [answer(list())];
		writes = {};
		status = answer({ block: true, deployable: true });
		api.mockReset();
		api.mockImplementation(async (path: string) => {
			if (path === '/api/kluster/deployments') return reads.length > 1 ? reads.shift()! : reads[0];
			if (path === '/api/kluster') return status;
			if (path.startsWith('/api/kluster/log?'))
				return answer({ exists: true, offset: 0, next: 4, size: 4, text: 'log\n' });
			if (path.startsWith('/api/items?')) return answer([]);
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
			'S-0581 Story S-0581 in review',
			'S-0583 Story S-0583 in review',
			'S-0584 Not yet deployed in review'
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

	it('groups the page in its rail: Overview first, then Build, Logs, Access, and Setup', async () => {
		await show();
		const ids = [...document.querySelectorAll('[data-sidebar-entry]')].map((e) =>
			e.getAttribute('data-sidebar-entry')
		);
		expect(ids).toEqual(['overview', 'build', 'logs', 'access', 'setup']);
		expect(paneShown()).toBe('overview');
		expect(entry('overview').textContent).toContain('4');
		await go('setup');
		expect(paneShown()).toBe('setup');
		expect(window.location.hash).toBe('#setup');
		expect(all('deployment-row')).toEqual([]);
		await go('overview');
		expect(all('deployment-row').length).toBe(4);
	});

	it("shows each step's outcome on its row, and opens its log in Logs, with the log's path", async () => {
		await show();
		const steps = [
			...rowOf('S-0581').querySelectorAll<HTMLAnchorElement>('[data-testid="deployment-step"]')
		];
		expect(steps.map((s) => s.textContent!.trim())).toEqual(['build · passed', 'deploy · passed']);
		steps[1].click();
		await settle();
		expect(paneShown()).toBe('logs');
		expect(q('deployments-chosen')!.textContent).toContain('s-0581');
		const panels = all('step-panel');
		expect(panels.map((p) => p.getAttribute('data-step'))).toEqual(['build', 'deploy']);
		expect(panels[0].textContent).toContain('.flai-cache/kluster/logs/s-0581/build.log');
		// S-0785: each step's log is a link to the log page
		const link = panels[0].querySelector<HTMLAnchorElement>(
			'[data-testid="step-log"] [data-testid="log-link"]'
		)!;
		expect(link.dataset.path).toBe('.flai-cache/kluster/logs/s-0581/build.log');
		expect(new URL(link.href).pathname).toBe('/ops/log');
		expect(panels[0].querySelector('[data-testid="log-feed-status"]')!.textContent).toBe('passed');
	});

	it('builds one panel per component, each with its status: passed, running, and waiting', async () => {
		const building = story('S-0590', {
			outcome: 'running',
			ended_at: undefined,
			steps: [
				step('build-api', 'passed'),
				step('build-web', 'running'),
				step('build-worker', 'pending'),
				step('load', 'pending'),
				step('deploy', 'pending')
			]
		});
		reads = [
			answer(
				list({
					deployments: [main(), building],
					running: job('deploy', 'S-0590')
				})
			)
		];
		await show();
		await go('build');
		expect(q('deployments-chosen')!.textContent).toContain('s-0590');
		const panels = all('build-panel');
		expect(panels.map((p) => p.getAttribute('data-step'))).toEqual([
			'build-api',
			'build-web',
			'build-worker'
		]);
		expect(
			panels.map((p) => p.querySelector('[data-testid="log-feed-status"]')!.textContent)
		).toEqual(['passed', 'running', 'pending']);
		expect(panels.map((p) => p.querySelector('summary')!.textContent)).toEqual([
			expect.stringContaining('api'),
			expect.stringContaining('web'),
			expect.stringContaining('worker')
		]);
		expect(q('deployments-build-summary')!.textContent).toContain(
			'3 builds: 1 passed, 1 running, 1 waiting'
		);
		expect(q('build-load')!.querySelector('[data-testid="log-feed-status"]')!.textContent).toBe(
			'pending'
		);
		expect(entry('build').textContent).toContain('3');
	});

	it('shows a component whose build failed, and the builds after it skipped', async () => {
		const failed = story('S-0590', {
			outcome: 'failed',
			failed_step: 'build-web',
			steps: [
				step('build-api', 'passed'),
				step('build-web', 'failed', { exit_code: 2 }),
				step('build-worker', 'skipped'),
				step('load', 'skipped'),
				step('deploy', 'skipped')
			]
		});
		reads = [answer(list({ deployments: [main(), failed] }))];
		await show();
		await go('build');
		// main is the default without a job running: the reader picks the story's namespace
		const select = q('deployments-build-namespace') as HTMLSelectElement;
		select.value = 's-0590';
		select.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(
			all('build-panel').map((p) => p.querySelector('[data-testid="log-feed-status"]')!.textContent)
		).toEqual(['passed', 'failed', 'skipped']);
		expect(all('build-panel')[1].textContent).toContain('exit code 2');
		expect(q('deployments-build-summary')!.textContent).toContain('1 passed, 1 failed, 1 skipped');
		expect(q('deployments-chosen')!.textContent).toContain('at its build-web step');
	});

	it('says a deploy that built nothing names no build', async () => {
		const bare = story('S-0581', {
			steps: [step('build', 'skipped'), step('load', 'skipped'), step('deploy', 'passed')]
		});
		reads = [answer(list({ deployments: [bare] }))];
		await show();
		await go('build');
		expect(q('deployments-build-none')!.textContent).toContain('names no build command');
		expect(all('build-panel')).toEqual([]);
	});

	it('shows how to reach each deployment for a demo: its hosts and ports, with no command', async () => {
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
		await go('access');
		const demo = accessOf('main').querySelector('[data-testid="deployment-access"]')!;
		expect(demo.textContent).toContain('web.authstar.svc.cluster.local');
		expect(demo.textContent).toContain('port 8080');
		expect(buttonIn(demo, 'Start')).toBeDefined();
		expect(demo.textContent).not.toContain('kubectl');
		expect(demo.querySelector('[data-testid="deployment-port-forward"]')).toBeNull();
		expect(demo.querySelector('[data-testid="deployment-list"]')).toBeNull();
		expect(
			accessOf('S-0581').querySelector('[data-testid="deployment-services-error"]')!.textContent
		).toContain("kubectl is not on the host's PATH");
		expect(q('deployments-forwards-open')).toBeNull();
	});

	const running = (ns: string) => ({
		namespace: ns,
		service: 'web',
		port: 8080,
		local_port: 41234,
		address: '0.0.0.0',
		started_at: '2026-11-07T10:00:00Z',
		running: true
	});
	const withForwards = (forwards: ReturnType<typeof running>[]) =>
		list({
			deployments: [main(), story('S-0581', { access: { ...access('s-0581'), forwards } })]
		});

	it('starts a port-forward, links it on the host, says it is open to the network, and stops it', async () => {
		await show();
		await go('access');
		const row = () => accessOf('S-0581').querySelector('[data-testid="deployment-access"]')!;
		writes['/api/kluster/forward'] = answer(withForwards([running('s-0581')]));
		buttonIn(row(), 'Start')!.click();
		await settle();
		const [path, init] = posts().at(-1)!;
		expect(path).toBe('/api/kluster/forward');
		expect(JSON.parse(String((init as RequestInit).body))).toEqual({
			namespace: 's-0581',
			service: 'web',
			port: 8080,
			on: true
		});
		const link = row().querySelector<HTMLAnchorElement>('[data-testid="deployment-forward-link"]')!;
		expect(link.getAttribute('href')).toBe(`http://${window.location.hostname}:41234`);
		expect(q('deployments-forwards-open')!.textContent).toContain("reach the host's");
		expect(entry('access').textContent).toContain('1');
		expect(document.body.textContent).not.toMatch(/kubectl|flai (serve|kluster)/);

		writes['/api/kluster/forward'] = answer(withForwards([]));
		buttonIn(row(), 'Stop')!.click();
		await settle();
		expect(JSON.parse(String((posts().at(-1)![1] as RequestInit).body))).toMatchObject({
			on: false
		});
		expect(row().querySelector('[data-testid="deployment-forward-link"]')).toBeNull();
		expect(buttonIn(row(), 'Start')).toBeDefined();
		expect(q('deployments-forwards-open')).toBeNull();
	});

	it("shows flai's refusal of a forward beside its port, and an ended forward's output", async () => {
		reads = [
			answer(
				withForwards([
					{ ...running('s-0581'), running: false, output: 'lost connection to pod' } as never
				])
			)
		];
		await show();
		await go('access');
		const row = () => accessOf('S-0581').querySelector('[data-testid="deployment-access"]')!;
		expect(row().querySelector('[data-testid="deployment-forward-ended"]')!.textContent).toContain(
			'lost connection to pod'
		);
		const said =
			'the cluster authstar lists no service web with the port 8080 in the namespace s-0581; pick one kluster.deployments lists';
		writes['/api/kluster/forward'] = answer({ error: said }, 400);
		buttonIn(row(), 'Start')!.click();
		await settle();
		const port = row().querySelector('[data-testid="deployment-port"]')!;
		expect(port.querySelector('[data-testid="deployment-forward-said"]')!.textContent).toBe(said);
		expect(q('deployments-action-error')).toBeNull();
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
		expect(buttonIn(rowOf('S-0583'), 'Deploy')!.disabled).toBe(true);
		expect(buttonIn(rowOf('main'), 'Promote')!.disabled).toBe(true);
		await go('logs');
		expect(q('deployments-log')!.textContent).toBe('building s-0584\n');
		await go('overview');

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

	it("says the host action is off, linking the Klusters page's setup with no command, and disables every action", async () => {
		reads = [answer(list({ enabled: false }))];
		await show();
		const off = q('deployments-off')!;
		expect(off.textContent).not.toContain('flai ');
		expect(off.querySelector('a')!.getAttribute('href')).toBe('/ops/klusters');
		for (const b of document.querySelectorAll<HTMLButtonElement>(
			'[data-testid="deployment-row"] button'
		))
			expect(b.disabled).toBe(true);
	});

	it("takes a 403 on a write as the host action off, with flai's refusal word for word", async () => {
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
		expect(q('deployments-off')!.textContent).not.toContain('flai ');
		expect(q('deployments-off')!.querySelector('a')!.getAttribute('href')).toBe('/ops/klusters');
		expect(q('deployments-action-error')!.textContent!.trim()).toBe(
			`the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`
		);
	});

	it("says when the project has no kluster: block, linking the Klusters page's setup with no command", async () => {
		reads = [
			answer({
				...list({ deployments: [], review: [] }),
				block: false,
				no_block: `this project's system-flow.yaml has no kluster: block; turn the kluster host action on, with ${ENABLE}`
			})
		];
		status = answer({ block: false });
		await show();
		const nb = q('deployments-no-block')!;
		expect(nb.textContent).toContain('not set up yet');
		expect(nb.textContent).not.toContain('flai ');
		expect(nb.textContent).not.toContain('system-flow.yaml');
		expect(nb.querySelector('a')!.getAttribute('href')).toBe('/ops/klusters');
		expect(all('deployment-row')).toEqual([]);
	});

	it('with no deploy set, says deploys are not set up yet, linking the setup, in place of Deploy and Promote', async () => {
		status = answer({ block: true, deployable: false });
		await show();
		const nd = q('deployments-no-deploy')!;
		expect(nd.textContent).toContain('Deploys are not set up yet');
		expect(nd.textContent).not.toContain('flai ');
		expect(nd.querySelector('a')!.getAttribute('href')).toBe('#setup');
		expect(buttonIn(document, 'Deploy')).toBeUndefined();
		expect(buttonIn(document, 'Promote')).toBeUndefined();
		expect(rowOf('S-0584').querySelector('[data-testid="deployment-to-setup"]')).not.toBeNull();
		expect(rowOf('main').querySelector('[data-testid="deployment-to-setup"]')).not.toBeNull();
		// a deployed story's namespace can still be torn down
		expect(buttonIn(rowOf('S-0581'), 'Teardown')).toBeDefined();
		// the link opens Setup, and asks flai nothing
		rowOf('S-0584')
			.querySelector<HTMLAnchorElement>('[data-testid="deployment-to-setup"]')!
			.click();
		await settle();
		expect(paneShown()).toBe('setup');
		expect(posts()).toEqual([]);
		await go('overview');
		reads = [answer(list({ deployments: [], review: [] }))];
		if (c) unmount(c);
		document.body.innerHTML = '';
		await show();
		expect(q('deployments-no-main')).not.toBeNull();
		expect(buttonIn(q('deployments-no-main')!, 'Promote')).toBeUndefined();
	});

	it("holds the deploy setup in Setup, from kluster.status's settings, saving back into the page", async () => {
		const setting = (key: string, kind: string, value?: unknown) => ({
			key,
			kind,
			default: '',
			meaning: `the ${key}`,
			value,
			set: value !== undefined
		});
		status = answer({
			block: true,
			deployable: false,
			settings_on: true,
			settings_enable: 'flai serve enable settings',
			settings: [
				setting('kluster.config', 'text', 'kluster.yaml'),
				setting('kluster.namespace', 'text'),
				setting('kluster.build', 'command'),
				setting('kluster.components', 'components'),
				setting('kluster.deploy', 'command'),
				setting('kluster.teardown', 'command'),
				setting('kluster.timeout', 'duration')
			]
		});
		await show();
		expect(q('deployments-deployable')!.textContent).toContain('deploys not set up');
		await go('setup');
		const pane = q('deployments-setup-pane')!;
		expect(pane.querySelector('[data-testid="deploy-settings"]')).not.toBeNull();
		expect(pane.querySelector('[data-testid="field-kluster.deploy"]')).not.toBeNull();
		expect(
			pane
				.querySelector<HTMLAnchorElement>('[data-testid="deployments-setup"]')!
				.getAttribute('href')
		).toBe('/ops/klusters');
	});

	it('says where the setup is when kluster.status gives no settings', async () => {
		await show();
		await go('setup');
		expect(q('deployments-setup-unknown')!.textContent).toContain('Klusters page');
	});

	it('offers Deploy and Promote when kluster.status cannot be read, as from an older flai', async () => {
		status = answer({ error: 'no such method' }, 501);
		await show();
		expect(q('deployments-no-deploy')).toBeNull();
		expect(buttonIn(rowOf('S-0584'), 'Deploy')).toBeDefined();
		expect(buttonIn(rowOf('main'), 'Promote')).toBeDefined();
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

	// S-0663: two stories in progress, one at its head and checked, one behind its head, beside one
	// in review, and a story in progress and one in review not yet deployed, against a limit of 2,
	// with a deploy refused for it.
	const REASON =
		'kluster.agent_namespaces is 2, and S-0663 and S-0664 hold a namespace; tear one down, or wait for its story to be accepted';
	const checked = () =>
		story('S-0663', {
			status: 'in-progress',
			column: 'in-progress',
			at_head: true,
			asked_by: 'agent-S-0663',
			log: '.flai-cache/kluster/logs/s-0663',
			checks_commit: HEAD,
			checks_outcome: 'passed',
			checks_asked_by: 'agent-S-0663',
			checks_log: '.flai-cache/kluster/checks/s-0663.log',
			checks_log_tail: 'checks passed\n'
		});
	const live = (over: Partial<KlusterDeployments> = {}) =>
		list({
			deployments: [
				main(),
				story('S-0581', { column: 'review', at_head: true }),
				checked(),
				story('S-0664', {
					status: 'in-progress',
					column: 'in-progress',
					head: OLD,
					at_head: false,
					asked_by: 'agent-S-0664'
				})
			],
			review: [{ story: 'S-0584', title: 'In review, not deployed', head: HEAD }],
			undeployed: [
				{ story: 'S-0584', title: 'In review, not deployed', head: HEAD, column: 'review' },
				{ story: 'S-0665', title: 'In progress, not deployed', head: HEAD, column: 'in-progress' }
			],
			live: 2,
			holding: ['S-0663', 'S-0664'],
			limit: 2,
			refusals: [
				{
					story: 'S-0665',
					by: 'agent-S-0665',
					at: '2026-11-13T07:30:00Z',
					limit: 2,
					reason: REASON
				}
			],
			...over
		});
	const said = (el: Element, id: string) =>
		el.querySelector(`[data-testid="${id}"]`)!.textContent!.replace(/\s+/g, ' ').trim();

	it('lists each story in progress beside those in review, with its column, head, deploy, checks, and who asked', async () => {
		reads = [answer(live())];
		await show();
		expect(all('deployment-row').map((r) => r.getAttribute('data-story'))).toEqual([
			null,
			'S-0581',
			'S-0663',
			'S-0664',
			'S-0584',
			'S-0665'
		]);
		expect(entry('overview').textContent).toContain('6');
		const row = rowOf('S-0663');
		expect(said(row, 'deployment-column')).toBe('in progress');
		expect(said(row, 'deployment-commit')).toBe('0123456789ab');
		expect(row.querySelector('[data-testid="deployment-at-head"]')).not.toBeNull();
		expect(said(row, 'deployment-outcome')).toBe('passed');
		expect(said(row, 'deployment-asked-by')).toContain('asked by agent-S-0663');
		expect(said(row, 'deployment-checks-outcome')).toBe('passed');
		expect(said(row, 'deployment-checks-asked-by')).toBe('at 0123456789ab, asked by agent-S-0663');
		expect(said(row, 'deployment-checks-tail')).toBe('checks passed');
		expect(buttonIn(row, 'Run checks')).toBeDefined();
		expect(buttonIn(row, 'Deploy')).toBeUndefined();
		expect(buttonIn(row, 'Teardown')).toBeDefined();
		// behind its head: deployed again, not checked
		expect(rowOf('S-0664').querySelector('[data-testid="deployment-behind"]')).not.toBeNull();
		expect(buttonIn(rowOf('S-0664'), 'Deploy')).toBeDefined();
		expect(buttonIn(rowOf('S-0664'), 'Run checks')).toBeUndefined();
		expect(said(rowOf('S-0581'), 'deployment-column')).toBe('in review');
		// not deployed yet, from undeployed rather than review
		const waiting = rowOf('S-0665');
		expect(said(waiting, 'deployment-column')).toBe('in progress');
		expect(said(waiting, 'deployment-outcome')).toBe(
			'In progress, not deployed; its head is 0123456789ab.'
		);
		expect(buttonIn(waiting, 'Deploy')).toBeDefined();
		expect(buttonIn(waiting, 'Teardown')).toBeUndefined();
		expect(said(rowOf('S-0584'), 'deployment-outcome')).toContain('In review, not deployed');
	});

	it('deploys, runs the checks of, and tears down a story in progress, each through flai', async () => {
		reads = [answer(live())];
		await show();
		const body = (i: number) => {
			const [path, init] = posts()[i];
			return [path, JSON.parse(String((init as RequestInit).body))];
		};
		writes['/api/kluster/deploy'] = answer(live());
		buttonIn(rowOf('S-0665'), 'Deploy')!.click();
		await settle();
		buttonIn(rowOf('S-0664'), 'Deploy')!.click();
		await settle();
		writes['/api/kluster/checks'] = answer(live());
		buttonIn(rowOf('S-0663'), 'Run checks')!.click();
		await settle();
		writes['/api/kluster/teardown'] = answer(live({ running: job('teardown', 'S-0663') }));
		buttonIn(rowOf('S-0663'), 'Teardown')!.click();
		flushSync();
		expect(posts().length).toBe(3);
		buttonIn(rowOf('S-0663'), 'Tear it down')!.click();
		await settle();
		expect([0, 1, 2, 3].map(body)).toEqual([
			['/api/kluster/deploy', { story: 'S-0665' }],
			['/api/kluster/deploy', { story: 'S-0664' }],
			['/api/kluster/checks', { story: 'S-0663' }],
			['/api/kluster/teardown', { story: 'S-0663' }]
		]);
		expect(q('deployments-running')!.textContent).toContain('kluster teardown of S-0663');
		expect(buttonIn(rowOf('S-0663'), 'Run checks')!.disabled).toBe(true);
	});

	it("shows flai's refusal of a checks run word for word", async () => {
		const no =
			'the namespace s-0663 is deployed at fedcba987654, not at the head of story/S-0663; deploy first with flai kluster deploy S-0663 --job';
		reads = [answer(live())];
		writes['/api/kluster/checks'] = answer({ error: no }, 400);
		await show();
		buttonIn(rowOf('S-0663'), 'Run checks')!.click();
		await settle();
		expect(q('deployments-action-error')!.textContent!.trim()).toBe(no);
	});

	it('says how many story namespaces live against the limit, and which request was refused for it', async () => {
		reads = [answer(live())];
		await show();
		expect(said(document.body, 'deployments-limit')).toBe('2 of 2 story namespaces');
		const line = q('deployments-limit-line')!;
		expect(said(line, 'deployments-refused')).toBe(
			`Refused for the limit: S-0665's deploy, asked by agent-S-0665 at ${localTime('2026-11-13T07:30:00Z')}.`
		);
		expect(said(line, 'deployments-refused-reason')).toBe(REASON);
		expect(line.textContent!.replace(/\s+/g, ' ')).toContain(
			'2 of 2 story namespaces live for stories in progress: S-0663, S-0664. Stories in review do not count.'
		);
	});

	it('says deploys before review are off at a limit of 0, and nothing of the limit from an older flai', async () => {
		reads = [answer(live({ live: 0, holding: [], limit: 0, refusals: [] }))];
		await show();
		expect(said(document.body, 'deployments-limit')).toBe('deploys before review off');
		expect(said(document.body, 'deployments-limit-line')).toBe(
			'Deploys before review are off: no story in progress may hold a namespace. Stories in review do not count.'
		);
		expect(q('deployments-refused')).toBeNull();
		if (c) unmount(c);
		document.body.innerHTML = '';
		reads = [answer(list())];
		await show();
		expect(q('deployments-limit')).toBeNull();
		expect(q('deployments-limit-line')).toBeNull();
	});

	it("starts and stops a port-forward into a story in progress's namespace, as one in review's", async () => {
		reads = [answer(live())];
		await show();
		await go('access');
		const row = () => accessOf('S-0663').querySelector('[data-testid="deployment-access"]')!;
		const withForward = (forwards: ReturnType<typeof running>[]) =>
			live({
				deployments: [
					main(),
					{ ...checked(), access: { ...access('s-0663'), forwards } },
					story('S-0581', { column: 'review' })
				]
			});
		writes['/api/kluster/forward'] = answer(withForward([running('s-0663')]));
		buttonIn(row(), 'Start')!.click();
		await settle();
		expect(JSON.parse(String((posts().at(-1)![1] as RequestInit).body))).toEqual({
			namespace: 's-0663',
			service: 'web',
			port: 8080,
			on: true
		});
		expect(
			row().querySelector<HTMLAnchorElement>('[data-testid="deployment-forward-link"]')!.href
		).toContain(':41234');
		writes['/api/kluster/forward'] = answer(withForward([]));
		buttonIn(row(), 'Stop')!.click();
		await settle();
		expect(JSON.parse(String((posts().at(-1)![1] as RequestInit).body))).toMatchObject({
			namespace: 's-0663',
			on: false
		});
		expect(buttonIn(row(), 'Start')).toBeDefined();
	});

	it("opens Logs at a deploy's failed step from its row's log link", async () => {
		reads = [
			answer(
				live({
					deployments: [
						main(),
						story('S-0663', {
							column: 'in-progress',
							outcome: 'failed',
							failed_step: 'deploy',
							steps: [step('build', 'passed'), step('deploy', 'failed', { exit_code: 1 })]
						})
					]
				})
			)
		];
		await show();
		rowOf('S-0663').querySelector<HTMLAnchorElement>('[data-testid="deployment-logs"]')!.click();
		await settle();
		expect(paneShown()).toBe('logs');
		expect(q('deployments-chosen')!.textContent).toContain('s-0663');
		const open = all('step-panel').map(
			(p) => p.querySelector<HTMLDetailsElement>('[data-testid="log-feed"]')!.open
		);
		expect(open).toEqual([false, true]);
	});
});
