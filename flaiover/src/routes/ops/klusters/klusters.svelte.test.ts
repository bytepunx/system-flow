// S-0580, S-0586: the Klusters page leads an operator through setting up the project's kluster, in
// five steps that each say whether they are done (the host's tools, turning the kluster host action
// on, which writes the block, the config file and its profile, the commands that deploy, and Up),
// with no command to run and no file to edit but the one host step, turning settings on; then it
// manages the cluster: its name, profile, provider, state, and age, Up, and Down once confirmed,
// through /api/kluster, the job that runs and asked again until it ends, flai check's findings,
// the nodes and pods, and turning the action off. flai's refusals are shown word for word.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type {
	KlusterJob,
	KlusterSetting,
	KlusterStatus,
	KlusterTool,
	KlusterWorkloads
} from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params?: Record<string, string>) =>
		params ? route.replace(/\[(?:\.\.\.)?(\w+)\]/g, (_, k: string) => params[k]) : route
}));

import KlustersPage from './+page.svelte';

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
const SETTINGS_ENABLE = 'flai serve enable settings';
const NO_BLOCK =
	"this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; turn the kluster host action on";
const CONFIG = 'name: p-local\nprovider: k3d\nprofile: spire\n';
const HASH = 'a'.repeat(64);

const tool = (name: KlusterTool['name'], found = true): KlusterTool => ({
	name,
	found,
	...(found ? { path: `/usr/bin/${name}` } : {}),
	need: `flai needs ${name}`,
	...(name === 'docker' ? {} : { install: name === 'kubectl' ? 'kluster setup' : 'go install x' }),
	docs: `https://docs/${name}`
});
const TOOLS = [tool('kluster'), tool('docker'), tool('kubectl')];

const setting = (key: string, over: Partial<KlusterSetting> = {}): KlusterSetting => ({
	key,
	kind: ['kluster.build', 'kluster.deploy', 'kluster.teardown'].includes(key)
		? 'command'
		: key === 'kluster.timeout'
			? 'duration'
			: 'text',
	meaning: `what ${key} does`,
	set: false,
	...over
});
const SETTINGS = [
	setting('kluster.config', { value: 'kluster.yaml', set: true }),
	setting('kluster.namespace', { default: 'p' }),
	setting('kluster.build'),
	setting('kluster.deploy', { value: ['scripts/kluster-deploy.sh'], set: true }),
	setting('kluster.teardown'),
	setting('kluster.timeout', { default: '30m' })
];

const job = (action: 'up' | 'down', over: Partial<KlusterJob> = {}): KlusterJob => ({
	action,
	cluster: 'authstar',
	started_at: '2026-11-07T01:00:00Z',
	log: `.flai-cache/kluster/${action}.log`,
	...over
});
const status = (over: Partial<KlusterStatus> = {}): KlusterStatus => ({
	enabled: true,
	enable: ENABLE,
	block: true,
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
	settings: SETTINGS,
	settings_on: true,
	settings_enable: SETTINGS_ENABLE,
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
	findings: [],
	deployable: true,
	tools: TOOLS,
	...over
});
// A project with no block: the action off, the tools there, settings on.
const fresh = (over: Partial<KlusterStatus> = {}): KlusterStatus => ({
	enabled: false,
	enable: ENABLE,
	block: false,
	no_block: NO_BLOCK,
	settings: SETTINGS.map((s) => ({ ...s, value: undefined, set: false })),
	settings_on: true,
	settings_enable: SETTINGS_ENABLE,
	findings: [],
	deployable: false,
	tools: TOOLS,
	...over
});
const WORKLOADS: KlusterWorkloads = {
	block: true,
	cluster: 'authstar',
	nodes: [
		{ name: 'k3d-authstar-server-0', ready: true, roles: ['control-plane'], version: 'v1.31' }
	],
	namespaces: [
		{
			namespace: 'kube-system',
			pods: [{ name: 'coredns-1', phase: 'Running', ready: 1, containers: 1, restarts: 2 }]
		}
	]
};

const q = (id: string) => document.querySelector(`[data-testid="${id}"]`);
const buttonNamed = (name: string) =>
	[...document.querySelectorAll<HTMLButtonElement>('button')].find(
		(b) => b.textContent!.trim() === name
	);
const writesTo = (path: string) =>
	api.mock.calls.filter(
		([p, init]) => p === path && ['POST', 'PUT'].includes((init as RequestInit)?.method ?? '')
	);
const sent = (path: string, n = 0) =>
	JSON.parse(String((writesTo(path)[n][1] as RequestInit).body)) as Record<string, unknown>;
const type = (el: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement, value: string) => {
	el.value = value;
	el.dispatchEvent(new Event(el instanceof HTMLSelectElement ? 'change' : 'input'));
	flushSync();
};
// A flai command the operator would have to run, or the manifest to edit by hand.
const COMMAND =
	/flai (serve|kluster|manifest|config|check)\b|edit system-flow\.yaml|to system-flow\.yaml/;
const fieldOf = (key: string) =>
	q(`field-${key}`)!.querySelector<HTMLInputElement | HTMLTextAreaElement>('input, textarea')!;

describe('Klusters page (S-0580, S-0586)', () => {
	let c: ReturnType<typeof mount> | undefined;
	// the answers GET /api/kluster gives, in turn; the last is kept
	let reads: ReturnType<typeof answer>[] = [];
	let writes: Record<string, ReturnType<typeof answer>> = {};
	let gets: Record<string, ReturnType<typeof answer>> = {};

	const show = async () => {
		c = mount(KlustersPage, { target: document.body });
		await settle();
	};

	beforeEach(() => {
		reads = [answer(status())];
		writes = {};
		gets = {
			'/api/kluster/workloads': answer(WORKLOADS),
			'/api/kluster/config': answer({ path: 'kluster.yaml', content: CONFIG, hash: HASH }),
			'/api/items?type=epic&archived=false': answer([
				{ id: 'E-0001', title: 'Ship it', status: 'in-progress' },
				{ id: 'E-0002', title: 'Old', status: 'done' }
			])
		};
		api.mockReset();
		api.mockImplementation(async (path: string, init?: RequestInit) => {
			if (init?.method === 'POST' || init?.method === 'PUT') {
				const w = writes[path];
				if (!w) throw new Error(`unexpected ${init.method} ${path}`);
				return w;
			}
			if (path === '/api/kluster') return reads.length > 1 ? reads.shift()! : reads[0];
			const g = gets[path];
			if (!g) throw new Error(`unexpected ${path}`);
			return g;
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.useRealTimers();
	});

	it("shows the cluster's name, profile, provider, state, and age", async () => {
		await show();
		const text = q('kluster-cluster')!.textContent!;
		for (const s of ['authstar', 'integration', 'k3d', 'Running', '2h50m39s']) {
			expect(text).toContain(s);
		}
		expect(buttonNamed('Up')).toBeDefined();
		expect(buttonNamed('Down')).toBeDefined();
	});

	it('offers no Up for a cluster that is running already', async () => {
		await show();
		const ups = [...document.querySelectorAll('button')].filter((b) => b.textContent === 'Up');
		expect(ups.length).toBeGreaterThan(0);
		for (const up of ups) expect(up.disabled).toBe(true);
	});

	it('folds the setup away once every step is done, and shows the nodes and pods', async () => {
		await show();
		const setup = q('kluster-setup') as HTMLDetailsElement;
		expect(setup.open).toBe(false);
		expect(
			[...setup.querySelectorAll('[data-testid="step-state"]')].map((s) => s.textContent)
		).toEqual(['Done', 'Done', 'Done', 'Done', 'Done']);
		expect(q('kluster-nodes')!.textContent).toContain('k3d-authstar-server-0');
		expect(q('kluster-nodes')!.textContent).toContain('control-plane');
		const pods = q('kluster-pods')!.textContent!;
		for (const s of ['kube-system', 'coredns-1', 'Running', '1/1 ready', '2 restarts']) {
			expect(pods).toContain(s);
		}
		const deployments = [...document.querySelectorAll('a')].find(
			(a) => a.textContent === 'Deployments'
		);
		expect(deployments!.getAttribute('href')).toBe('/ops/deployments');
	});

	it('says why the nodes and pods cannot be read', async () => {
		gets['/api/kluster/workloads'] = answer({
			...WORKLOADS,
			nodes: [],
			namespaces: [],
			workloads_error: 'flai has no kubeconfig for the cluster yet: bring it up first'
		});
		await show();
		expect(q('kluster-workloads-error')!.textContent!.trim()).toBe(
			'flai has no kubeconfig for the cluster yet: bring it up first'
		);
	});

	it('walks a project with no block from its tools to Up', async () => {
		const scaffold = status({
			cluster: {
				name: 'p-local',
				profile: 'spire',
				provider: 'k3d',
				context: 'k3d-p-local',
				exists: false,
				running: false
			},
			deployable: false,
			settings: SETTINGS.map((s) =>
				s.key === 'kluster.config' ? s : { ...s, value: undefined, set: false }
			)
		});
		reads = [answer(fresh())];
		await show();
		const setup = q('kluster-setup') as HTMLDetailsElement;
		expect(setup.open).toBe(true);
		// step 1: every tool found
		expect(q('step-tools')!.textContent).toContain('Done');
		// nothing names a command or a file to edit
		expect(document.body.textContent).not.toMatch(COMMAND);
		expect(q('kluster-cluster')).toBeNull();

		// step 2: turning it on writes the block and the config
		writes['/api/settings'] = answer({
			action: 'kluster',
			enabled: true,
			kluster: {
				wrote: ['system-flow.yaml', 'kluster.yaml'],
				config: 'kluster.yaml',
				commit: 'abc1234'
			}
		});
		reads = [answer(scaffold)];
		buttonNamed('Turn kluster on')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({ kind: 'action', action: 'kluster', on: true });
		expect(q('kluster-wrote')!.textContent).toContain('system-flow.yaml, kluster.yaml');
		expect(q('kluster-wrote')!.textContent).toContain('abc1234');
		expect(q('step-on')!.textContent).toContain('Done');

		// step 3: the config's profile, changed before the first Up
		expect(q('kluster-profile')!.textContent).toContain('spire');
		writes['/api/kluster/config'] = answer({ path: 'kluster.yaml', hash: 'b'.repeat(64) });
		type(q('kluster-profile-choice') as HTMLSelectElement, 'signet');
		reads = [answer({ ...scaffold, config_file: { ...scaffold.config_file!, profile: 'signet' } })];
		buttonNamed('Save the profile')!.click();
		await settle();
		// the picker reads the profile saved, not the one before
		expect((q('kluster-profile-choice') as HTMLSelectElement).value).toBe('signet');
		expect(sent('/api/kluster/config')).toMatchObject({
			path: 'kluster.yaml',
			hash: HASH,
			content: 'name: p-local\nprovider: k3d\nprofile: signet\n'
		});
		expect(q('kluster-config-said')!.textContent).toContain('Saved');

		// step 4: no deploy yet, so the page says nothing is deployed; set one
		expect(q('kluster-no-deploy')!.textContent).toContain('no story is deployed');
		type(fieldOf('kluster.deploy'), 'scripts/kluster-deploy.sh\n--wait\n');
		writes['/api/settings'] = answer({ commit: 'def5678' });
		reads = [answer(status({ cluster: scaffold.cluster }))];
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings', 1)).toEqual({
			kind: 'manifest',
			set: { 'kluster.deploy': ['scripts/kluster-deploy.sh', '--wait'] }
		});
		expect(q('kluster-save-deploy')!.textContent).toContain('def5678');
		expect(q('kluster-no-deploy')).toBeNull();

		// step 5: Up
		writes['/api/kluster/up'] = answer(status({ cluster: scaffold.cluster, running: job('up') }));
		q('step-up')!.querySelector('button')!.click();
		await settle();
		expect(writesTo('/api/kluster/up')).toHaveLength(1);
		expect(q('kluster-running')!.textContent).toContain('kluster up of the cluster authstar runs');
	});

	it('shows a missing tool with how to install it, and only a docs link for docker', async () => {
		reads = [
			answer(fresh({ tools: [tool('kluster'), tool('docker', false), tool('kubectl', false)] }))
		];
		await show();
		expect(q('step-tools')!.textContent).toContain('To do');
		const docker = q('tool-docker')!;
		expect(docker.textContent).toContain('missing');
		expect(docker.textContent).toContain('flai needs docker');
		expect(docker.querySelector('a')!.getAttribute('href')).toBe('https://docs/docker');
		expect(docker.querySelector('code')).toBeNull();
		expect(q('tool-kubectl')!.querySelector('code')!.textContent).toBe('kluster setup');
		expect(buttonNamed('Install')).toBeUndefined();
	});

	// S-0595: step 4 is DeploySettings, the form the Deployments page's Setup shares; its fields,
	// its refusals, and its agent story are tested in DeploySettings.svelte.test.ts (ADR-0263).
	it("holds the deploy setup form in step 4, and links the Deployments page's Setup", async () => {
		await show();
		const step = q('step-deploy')!;
		expect(step.querySelector('[data-testid="deploy-settings"]')).not.toBeNull();
		expect(step.querySelector('[data-testid="field-kluster.deploy"]')).not.toBeNull();
		expect(q('kluster-deploy-setup-link')!.getAttribute('href')).toBe('/ops/deployments#setup');
	});

	it('says once that settings is off, with the one host step, and saves nothing', async () => {
		reads = [answer(fresh({ settings_on: false }))];
		await show();
		const off = q('kluster-settings-off')!;
		expect(off.textContent).toContain('The settings host action is off');
		expect(off.textContent).toContain(SETTINGS_ENABLE);
		expect(buttonNamed('Turn kluster on')!.disabled).toBe(true);
		// the one command the page names is the settings step
		const text = document.body.textContent!.replace(off.textContent!, '');
		expect(text).not.toMatch(COMMAND);
		expect(document.querySelectorAll('[data-testid="kluster-settings-off"]')).toHaveLength(1);
	});

	it('saves the config file with settings off, but not the block or its commands', async () => {
		reads = [answer(status({ settings_on: false }))];
		await show();
		expect((q('kluster-profile-choice') as HTMLSelectElement).disabled).toBe(false);
		expect((q('kluster-config-text') as HTMLTextAreaElement).disabled).toBe(false);
		expect(buttonNamed('Save the config path')!.disabled).toBe(true);
		expect(fieldOf('kluster.deploy').disabled).toBe(true);
	});

	it('turns the kluster host action off', async () => {
		await show();
		writes['/api/settings'] = answer({ action: 'kluster', enabled: false });
		reads = [answer(status({ enabled: false }))];
		buttonNamed('Turn kluster off')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({ kind: 'action', action: 'kluster', on: false });
		expect(q('kluster-off')!.textContent).toContain('The kluster host action is off');
		expect(buttonNamed('Down')).toBeUndefined();
	});

	it("shows flai check's findings on the block", async () => {
		reads = [answer(status({ findings: ['kluster.timeout "0s" is not more than zero'] }))];
		await show();
		expect(q('kluster-findings')!.textContent).toContain('kluster.timeout "0s"');
	});

	it('says a cluster is stopped or not created', async () => {
		reads = [answer(status({ cluster: { ...status().cluster!, running: false, age: undefined } }))];
		await show();
		expect(q('kluster-state')!.textContent).toBe('Stopped');
		unmount(c!);
		reads = [answer(status({ cluster: { ...status().cluster!, running: false, exists: false } }))];
		c = mount(KlustersPage, { target: document.body });
		await settle();
		expect(q('kluster-state')!.textContent).toBe('Not created');
	});

	it('brings the cluster up and shows the job running, the end of its log, and asks again until it ends', async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		reads = [answer(status({ cluster: { ...status().cluster!, running: false } }))];
		await show();
		const running = job('up', { log_tail: 'creating cluster authstar\n' });
		writes['/api/kluster/up'] = answer(status({ running }));
		q('kluster-cluster')!.parentElement!.querySelector<HTMLButtonElement>('button')!.click();
		await settle();
		expect(writesTo('/api/kluster/up')).toHaveLength(1);
		expect(q('kluster-running')!.textContent).toContain('kluster up of the cluster authstar runs');
		expect(q('kluster-log')!.textContent).toBe('creating cluster authstar\n');
		expect(buttonNamed('Down')!.disabled).toBe(true);

		const done = job('up', {
			ended_at: '2026-11-07T01:03:00Z',
			outcome: 'done',
			log_tail: 'cluster authstar is up\n'
		});
		reads = [answer(status({ running })), answer(status({ last: done }))];
		const asked = () => api.mock.calls.filter(([p]) => p === '/api/kluster').length;
		const nodesRead = () => api.mock.calls.filter(([p]) => p === '/api/kluster/workloads').length;
		const before = asked();
		const nodesBefore = nodesRead();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 1);
		expect(q('kluster-running')).not.toBeNull();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 2);
		expect(q('kluster-running')).toBeNull();
		expect(q('kluster-last')!.textContent).toContain('done');
		expect(q('kluster-log')!.textContent).toBe('cluster authstar is up\n');
		// the nodes and pods are read again once the job ends
		expect(nodesRead()).toBe(nodesBefore + 1);
		// with nothing running, the page does not ask again
		vi.advanceTimersByTime(30000);
		await settle();
		expect(asked()).toBe(before + 2);
	});

	it('takes the cluster down only once the operator confirms', async () => {
		await show();
		writes['/api/kluster/down'] = answer(status({ running: job('down') }));
		buttonNamed('Down')!.click();
		flushSync();
		expect(q('kluster-confirm-down')!.textContent).toContain('Take the cluster authstar down?');
		expect(writesTo('/api/kluster/down')).toEqual([]);
		buttonNamed('Keep it')!.click();
		flushSync();
		expect(q('kluster-confirm-down')).toBeNull();

		buttonNamed('Down')!.click();
		flushSync();
		buttonNamed('Take it down')!.click();
		await settle();
		expect(writesTo('/api/kluster/down')).toHaveLength(1);
		expect(q('kluster-confirm-down')).toBeNull();
		expect(q('kluster-running')!.textContent).toContain(
			'kluster down of the cluster authstar runs'
		);
	});

	it("shows a failed job's error as flai gives it", async () => {
		reads = [
			answer(
				status({
					last: job('up', {
						ended_at: '2026-11-07T01:01:00Z',
						outcome: 'failed',
						error: "kluster is not on the host's PATH: install it, as design/tech/kluster.md says"
					})
				})
			)
		];
		await show();
		expect(q('kluster-last')!.textContent).toContain('failed');
		expect(q('kluster-last-error')!.textContent).toBe(
			"kluster is not on the host's PATH: install it, as design/tech/kluster.md says"
		);
	});

	it("links the job's output and a log its error names to the log page (S-0785)", async () => {
		const failed = '.flai-cache/kluster/logs/up-k3d.log';
		reads = [
			answer(
				status({
					last: job('up', {
						ended_at: '2026-11-07T01:01:00Z',
						outcome: 'failed',
						error: `k3d cluster create exited 1; read ${failed}`
					})
				})
			)
		];
		await show();
		const link = (id: string) =>
			q(id)!.querySelector<HTMLAnchorElement>('[data-testid="log-link"]')!;
		const output = link('kluster-job-log');
		expect(output.dataset.path).toBe('.flai-cache/kluster/up.log');
		expect(new URL(output.href).pathname).toBe('/ops/log');
		expect(q('kluster-job-log')!.textContent!.trim()).toBe(
			'Its output: .flai-cache/kluster/up.log'
		);
		expect(link('kluster-last-error').dataset.path).toBe(failed);
	});

	it("shows flai's refusal word for word", async () => {
		reads = [answer(status({ cluster: { ...status().cluster!, running: false } }))];
		await show();
		const refusal =
			'kluster up of the cluster authstar runs since 2026-11-07T01:00:00Z; wait for it to end, as kluster.status says, then ask again';
		writes['/api/kluster/up'] = answer({ error: refusal }, 400);
		q('kluster-cluster')!.parentElement!.querySelector<HTMLButtonElement>('button')!.click();
		await settle();
		expect(q('kluster-action-error')!.textContent!.trim()).toBe(refusal);
	});

	it('shows the status error kluster gave, word for word', async () => {
		const why = "kluster is not on the host's PATH";
		reads = [answer(status({ status_error: why }))];
		await show();
		expect(q('kluster-status-error')!.textContent!.trim()).toBe(why);
	});

	it('says the host action is off, with no command, in place of the buttons', async () => {
		reads = [answer(status({ enabled: false }))];
		await show();
		const off = q('kluster-off')!;
		expect(off.textContent).toContain('The kluster host action is off');
		expect(off.textContent).not.toMatch(COMMAND);
		expect(buttonNamed('Down')).toBeUndefined();
		expect(buttonNamed('Turn kluster on')).toBeDefined();
		// the cluster is still shown
		expect(q('kluster-cluster')!.textContent).toContain('authstar');
	});

	it('turns to the off state when a write is refused because the action went off', async () => {
		reads = [answer(status({ cluster: { ...status().cluster!, running: false } }))];
		await show();
		const message = `the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`;
		writes['/api/kluster/up'] = answer({ error: message, action: 'kluster', enable: ENABLE }, 403);
		q('kluster-cluster')!.parentElement!.querySelector<HTMLButtonElement>('button')!.click();
		await settle();
		expect(q('kluster-off')).not.toBeNull();
		expect(q('kluster-action-error')!.textContent!.trim()).toBe(message);
		expect(buttonNamed('Down')).toBeUndefined();
	});

	it('shows no cluster and no no-block notice with no block, only the setup', async () => {
		reads = [answer(fresh())];
		await show();
		expect(q('kluster-no-block')).toBeNull();
		expect(q('kluster-cluster')).toBeNull();
		expect(document.body.textContent).not.toContain(NO_BLOCK);
		expect(buttonNamed('Turn kluster on')).toBeDefined();
	});

	it('tells to upgrade flai on the host when it is too old', async () => {
		reads = [answer({ error: 'the host flai is too old to show the kluster' }, 501)];
		await show();
		expect(q('kluster-too-old')!.textContent).toContain('upgrade flai on the host');
		expect(buttonNamed('Up')).toBeUndefined();
	});
});
