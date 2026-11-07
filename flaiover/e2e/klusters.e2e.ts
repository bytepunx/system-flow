// S-0586: the Klusters page leads an operator from a project with no kluster: block to a running
// cluster, in Chromium, with no shell and no file edit. The page's own requests are answered in the
// browser by a fake host that keeps what the page wrote, so each step reads back what the step
// before it saved: the tools (kubectl missing, with its install line), turning kluster on (which
// writes the block and kluster.yaml and commits them), the config's profile changed from spire to
// signet, the deploy command saved as an argument list, Up, and then the managing view with the
// cluster's nodes and pods. A second case has the settings host action off: one notice says so,
// the block's fields cannot be saved, and the config's profile still can. Each step is saved as a
// full-page screenshot in the test's output folder. No visible text names a flai command or tells
// the operator to edit system-flow.yaml, outside the settings notice and the install line.
import { expect, test, type Page, type Route, type TestInfo } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const token = () => readFileSync(join(dir, 'token'), 'utf8').trim();

const HASH = 'a'.repeat(64);
const SETTINGS_ENABLE = 'flai serve enable settings';
const CONFIG = 'name: p-local\nprovider: k3d\nprofile: spire\n';

/** A flai command, or an instruction to edit the manifest by hand. */
const COMMAND = /\bflai (serve|kluster|manifest|config|check|edit)\b|edit[^.]*system-flow\.yaml/i;

type Setting = {
	key: string;
	kind: 'text' | 'command' | 'duration';
	default?: string;
	meaning: string;
	value?: string | string[];
	set: boolean;
};

const setting = (key: string, kind: Setting['kind'], def?: string): Setting => ({
	key,
	kind,
	...(def ? { default: def } : {}),
	meaning: `what ${key.replace('kluster.', '')} does`,
	set: false
});

/** The fake host: what flai on the host holds for the project, changed by the page's writes. */
function host(settingsOn = true) {
	const s = {
		enabled: false,
		block: false,
		config: CONFIG,
		hash: HASH,
		running: false,
		upped: false,
		settings: [
			setting('kluster.config', 'text'),
			setting('kluster.namespace', 'text', 'p'),
			setting('kluster.build', 'command'),
			setting('kluster.deploy', 'command'),
			setting('kluster.teardown', 'command'),
			setting('kluster.timeout', 'duration', '30m')
		],
		sent: [] as { at: string; method: string; body: unknown }[]
	};
	const profile = () => /^profile:\s*(\S+)/m.exec(s.config)?.[1] ?? '';
	const deployable = () => {
		const d = s.settings.find((x) => x.key === 'kluster.deploy');
		return Array.isArray(d?.value) && d.value.length > 0;
	};
	const status = () => ({
		enabled: s.enabled,
		enable: 'flai serve enable kluster',
		block: s.block,
		...(s.block
			? {
					cluster: {
						name: 'p-local',
						profile: profile(),
						provider: 'k3d',
						context: 'k3d-p-local',
						exists: s.upped,
						running: s.upped,
						...(s.upped ? { age: '1m2s', age_seconds: 62 } : {})
					},
					config: 'kluster.yaml',
					config_file: {
						exists: true,
						name: 'p-local',
						provider: 'k3d',
						profile: profile(),
						default_profile: 'spire',
						profiles: ['spire', 'signet', 'authstar'],
						addons: []
					}
				}
			: {
					no_block:
						"this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it"
				}),
		...(s.upped
			? {
					last: {
						action: 'up',
						cluster: 'p-local',
						started_at: '2026-11-07T09:00:00Z',
						ended_at: '2026-11-07T09:01:02Z',
						outcome: 'passed',
						log: '.flai-cache/kluster/up.log',
						log_tail: 'cluster "p-local" created\nspire: ready'
					}
				}
			: {}),
		settings: s.settings,
		settings_on: settingsOn,
		settings_enable: SETTINGS_ENABLE,
		findings: [],
		deployable: deployable(),
		tools: [
			{
				name: 'kluster',
				found: true,
				path: '/usr/local/bin/kluster',
				need: 'makes the cluster',
				docs: 'https://github.com/bytepunx/kluster'
			},
			{
				name: 'docker',
				found: true,
				path: '/usr/bin/docker',
				need: 'runs the cluster',
				docs: 'https://docs.docker.com/engine/install/'
			},
			{
				name: 'kubectl',
				found: false,
				need: 'flai lists the cluster’s nodes and pods and deletes namespaces with it',
				install: 'kluster setup',
				docs: 'https://kubernetes.io/docs/tasks/tools/'
			}
		]
	});
	const workloads = () => ({
		block: true,
		cluster: 'p-local',
		context: 'k3d-p-local',
		nodes: s.upped
			? [
					{
						name: 'k3d-p-local-server-0',
						ready: true,
						roles: ['control-plane', 'master'],
						version: 'v1.31.5+k3s1'
					}
				]
			: [],
		namespaces: s.upped
			? [
					{
						namespace: 'kube-system',
						pods: [
							{
								name: 'coredns-6799fbcd5-x2k9q',
								phase: 'Running',
								ready: 1,
								containers: 1,
								restarts: 0
							}
						]
					},
					{
						namespace: 'spire',
						pods: [
							{ name: 'spire-server-0', phase: 'Running', ready: 2, containers: 2, restarts: 1 }
						]
					}
				]
			: [],
		...(s.upped
			? {}
			: { workloads_error: 'the cluster is not up yet, so it has no nodes or pods to list' })
	});

	async function answer(route: Route) {
		const req = route.request();
		const at = new URL(req.url()).pathname;
		const method = req.method();
		const body = method === 'GET' ? undefined : req.postDataJSON();
		const json = (b: unknown, code = 200) => route.fulfill({ status: code, json: b });
		if (method !== 'GET') s.sent.push({ at, method, body });
		if (at === '/api/kluster' && method === 'GET') {
			if (s.running) [s.running, s.upped] = [false, true];
			return json(status());
		}
		if (at === '/api/kluster/workloads') return json(workloads());
		if (at === '/api/kluster/config') {
			if (!s.block) return json({ error: 'no block' }, 404);
			if (method === 'GET') return json({ path: 'kluster.yaml', content: s.config, hash: s.hash });
			s.config = body.content;
			s.hash = 'b'.repeat(64);
			return json({ path: 'kluster.yaml', hash: s.hash, log: [] });
		}
		// as on the host, Up answers with the job running; the page's next read finds it ended
		if (at === '/api/kluster/up' && method === 'POST') {
			s.running = true;
			return json({
				...status(),
				running: {
					action: 'up',
					cluster: 'p-local',
					started_at: '2026-11-07T09:00:00Z',
					log: '.flai-cache/kluster/up.log'
				}
			});
		}
		if (at === '/api/settings' && method === 'POST') {
			if (body.kind === 'action') {
				s.enabled = !!body.on;
				if (body.on && !s.block) {
					s.block = true;
					const c = s.settings.find((x) => x.key === 'kluster.config')!;
					c.value = 'kluster.yaml';
					c.set = true;
					return json({
						action: 'kluster',
						enabled: true,
						kluster: {
							wrote: ['system-flow.yaml', 'kluster.yaml'],
							config: 'kluster.yaml',
							commit: 'abc1234'
						}
					});
				}
				return json({ action: 'kluster', enabled: s.enabled });
			}
			if (body.kind === 'manifest') {
				for (const [key, value] of Object.entries(body.set ?? {})) {
					const x = s.settings.find((y) => y.key === key)!;
					x.value = value as string | string[];
					x.set = true;
				}
				for (const key of body.unset ?? []) {
					const x = s.settings.find((y) => y.key === key)!;
					x.value = undefined;
					x.set = false;
				}
				return json({ commit: 'def5678' });
			}
		}
		if (at === '/api/items' && method === 'GET')
			return json([{ id: 'E-0001', title: 'Local cluster', status: 'in-progress' }]);
		return route.fallback();
	}
	return { s, answer };
}

/**
 * The page's visible text, in the layout's main area (the shell's banners about the host's flai
 * are the layout's, not this page's), but for the settings notice and each missing tool's install
 * line.
 */
const textOutsideTheHostSteps = (page: Page) =>
	page.evaluate(() => {
		const main = document.querySelector('main');
		if (!main) return '';
		let text = main.innerText;
		const left = [...main.querySelectorAll<HTMLElement>('[data-testid="kluster-settings-off"]')];
		for (const e of main.querySelectorAll<HTMLElement>('span'))
			if (e.textContent?.startsWith('Install it on the host with')) left.push(e);
		for (const e of left) text = text.replace(e.innerText, '');
		return text;
	});

async function shot(page: Page, info: TestInfo, name: string) {
	await page.screenshot({ path: info.outputPath(`${name}.png`), fullPage: true });
	expect(await textOutsideTheHostSteps(page)).not.toMatch(COMMAND);
}

const step = (page: Page, id: string) => page.getByTestId(id);

test.use({ viewport: { width: 1280, height: 900 } });

test.beforeEach(async ({ context, baseURL }) => {
	await context.addCookies([{ name: 'flaiover_session', value: token(), url: baseURL! }]);
});

test('an operator sets up a project’s kluster from no block to a running cluster', async ({
	page,
	context
}, info) => {
	const { s, answer } = host();
	await context.route(/\/api\//, answer);
	await page.goto('/ops/klusters');

	// step 1: the tools, kubectl missing with its install line
	await expect(step(page, 'tool-kubectl')).toContainText('missing');
	await expect(step(page, 'tool-kubectl')).toContainText('kluster setup');
	await expect(step(page, 'tool-kluster')).toContainText('found');
	await shot(page, info, '1-no-block');

	// step 2: turning kluster on writes the block and the config, and commits them
	await page.getByRole('button', { name: 'Turn kluster on' }).click();
	await expect(step(page, 'kluster-wrote')).toContainText('system-flow.yaml, kluster.yaml');
	await expect(step(page, 'kluster-wrote')).toContainText('abc1234');
	await expect(step(page, 'step-on')).toContainText('Done');
	expect(s.sent.at(-1)).toMatchObject({
		at: '/api/settings',
		body: { kind: 'action', action: 'kluster', on: true }
	});
	await shot(page, info, '2-turned-on');

	// step 3: the profile reads spire, and is changed to signet
	await expect(step(page, 'kluster-profile')).toContainText('spire');
	await step(page, 'kluster-profile-choice').selectOption('signet');
	await page.getByRole('button', { name: 'Save the profile' }).click();
	await expect(step(page, 'kluster-config-said')).toContainText('Saved');
	await expect(step(page, 'kluster-profile')).toContainText('signet');
	expect(s.sent.at(-1)).toMatchObject({
		at: '/api/kluster/config',
		method: 'PUT',
		body: {
			path: 'kluster.yaml',
			hash: HASH,
			content: 'name: p-local\nprovider: k3d\nprofile: signet\n'
		}
	});
	await shot(page, info, '3-profile');

	// step 4: no deploy yet; the deploy command is saved as an argument list
	await expect(step(page, 'kluster-no-deploy')).toBeVisible();
	await step(page, 'field-kluster.deploy')
		.locator('textarea')
		.fill('scripts/kluster-deploy.sh\n--wait');
	await step(page, 'step-deploy').getByRole('button', { name: 'Save', exact: true }).click();
	await expect(step(page, 'kluster-save-deploy')).toContainText('def5678');
	expect(s.sent.at(-1)).toMatchObject({
		at: '/api/settings',
		body: { kind: 'manifest', set: { 'kluster.deploy': ['scripts/kluster-deploy.sh', '--wait'] } }
	});
	await expect(step(page, 'kluster-no-deploy')).toHaveCount(0);
	await expect(step(page, 'step-deploy')).toContainText('Done');
	await shot(page, info, '4-deploy');

	// step 5: Up
	await step(page, 'step-up').getByRole('button', { name: 'Up' }).click();
	await expect(step(page, 'kluster-running')).toContainText(
		'kluster up of the cluster p-local runs'
	);
	await shot(page, info, '5-up-runs');
	await expect(step(page, 'kluster-state')).toHaveText('Running', { timeout: 10_000 });
	expect(s.sent.at(-1)).toMatchObject({ at: '/api/kluster/up', method: 'POST' });

	// managing: the nodes and pods kubectl reports
	await expect(step(page, 'kluster-nodes')).toContainText('k3d-p-local-server-0');
	await expect(step(page, 'kluster-pods')).toContainText('spire-server-0');
	await expect(step(page, 'kluster-last')).toContainText('passed');
	await shot(page, info, '6-running');
});

test('with settings off, one notice says so, the block cannot be saved, and the profile can', async ({
	page,
	context
}, info) => {
	const { s, answer } = host(false);
	s.enabled = true;
	s.block = true;
	await context.route(/\/api\//, answer);
	await page.goto('/ops/klusters');

	await expect(step(page, 'kluster-settings-off')).toHaveCount(1);
	await expect(step(page, 'kluster-settings-off')).toContainText('One step on the host');
	await expect(step(page, 'kluster-settings-off')).toContainText(SETTINGS_ENABLE);
	await expect(step(page, 'field-kluster.deploy').locator('textarea')).toBeDisabled();
	await expect(page.getByRole('button', { name: 'Save the config path' })).toBeDisabled();
	await expect(step(page, 'kluster-profile-choice')).toBeEnabled();
	await step(page, 'kluster-profile-choice').selectOption('authstar');
	await expect(page.getByRole('button', { name: 'Save the profile' })).toBeEnabled();
	await shot(page, info, 'settings-off');
});
