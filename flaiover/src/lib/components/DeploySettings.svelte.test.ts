// S-0595, ADR-0263: DeploySettings edits the kluster: block's deploy setup through settings.manifest:
// the namespace, one build command or one build per component, the deploy and teardown commands,
// and the timeout, each command one argument a line. Components are rows, added, removed, and moved,
// saved as the JSON list flai manifest set takes; switching between one build and components unsets
// the other, since a block takes one. flai's refusals are shown beside their field. With settings
// off nothing can be saved. After a save it reads kluster.status again and gives it to saved.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { KlusterSetting, KlusterStatus } from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params?: Record<string, string>) =>
		params ? route.replace(/\[(?:\.\.\.)?(\w+)\]/g, (_, k: string) => params[k]) : route
}));

import DeploySettings from './DeploySettings.svelte';

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

const KIND: Record<string, KlusterSetting['kind']> = {
	'kluster.build': 'command',
	'kluster.deploy': 'command',
	'kluster.teardown': 'command',
	'kluster.timeout': 'duration',
	'kluster.components': 'components'
};
const setting = (key: string, over: Partial<KlusterSetting> = {}): KlusterSetting => ({
	key,
	kind: KIND[key] ?? 'text',
	meaning: `what ${key} does`,
	set: false,
	...over
});
const SETTINGS = [
	setting('kluster.config', { value: 'kluster.yaml', set: true }),
	setting('kluster.namespace', { default: 'p' }),
	setting('kluster.build', { value: ['scripts/build.sh'], set: true }),
	setting('kluster.components'),
	setting('kluster.deploy', { value: ['scripts/kluster-deploy.sh'], set: true }),
	setting('kluster.teardown'),
	setting('kluster.timeout', { default: '30m' })
];
const COMPONENTS = SETTINGS.map((s) =>
	s.key === 'kluster.build'
		? { ...s, value: undefined, set: false }
		: s.key === 'kluster.components'
			? {
					...s,
					value: [
						{ name: 'api', build: ['scripts/build.sh', 'api'] },
						{ name: 'web', build: ['scripts/build.sh', 'web'] }
					],
					set: true
				}
			: s
);
const STATUS = { enabled: true, enable: 'x', block: true, settings: SETTINGS } as KlusterStatus;

const q = (id: string) => document.querySelector(`[data-testid="${id}"]`);
const all = (id: string) => [...document.querySelectorAll(`[data-testid="${id}"]`)];
const buttonNamed = (name: string) =>
	[...document.querySelectorAll<HTMLButtonElement>('button')].find(
		(b) => b.textContent!.trim() === name
	);
const writesTo = (path: string) =>
	api.mock.calls.filter(([p, init]) => p === path && (init as RequestInit)?.method === 'POST');
const sent = (path: string, n = 0) =>
	JSON.parse(String((writesTo(path)[n][1] as RequestInit).body)) as Record<string, unknown>;
const type = (el: HTMLInputElement | HTMLTextAreaElement, value: string) => {
	el.value = value;
	el.dispatchEvent(new Event('input'));
	flushSync();
};
const fieldOf = (key: string) =>
	q(`field-${key}`)!.querySelector<HTMLInputElement | HTMLTextAreaElement>('input, textarea')!;
const pick = (value: 'one' | 'components') => {
	const r = q('kluster-build-mode')!.querySelector<HTMLInputElement>(`input[value="${value}"]`)!;
	r.checked = true;
	r.dispatchEvent(new Event('change'));
	flushSync();
};

describe('DeploySettings (S-0595, ADR-0263)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let writes: Record<string, ReturnType<typeof answer>> = {};
	const saved = vi.fn();

	const show = async (props: Record<string, unknown> = {}) => {
		c = mount(DeploySettings, {
			target: document.body,
			props: {
				settings: SETTINGS,
				settingsOn: true,
				settingsEnable: 'flai serve enable settings',
				deployable: true,
				saved,
				...props
			}
		});
		await settle();
	};

	beforeEach(() => {
		writes = {};
		saved.mockReset();
		api.mockReset();
		api.mockImplementation(async (path: string, init?: RequestInit) => {
			if (init?.method === 'POST') {
				const w = writes[path];
				if (!w) throw new Error(`unexpected POST ${path}`);
				return w;
			}
			if (path === '/api/kluster') return answer(STATUS);
			if (path === '/api/items?type=epic&archived=false')
				return answer([
					{ id: 'E-0001', title: 'Ship it', status: 'in-progress' },
					{ id: 'E-0002', title: 'Old', status: 'done' }
				]);
			throw new Error(`unexpected ${path}`);
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	it('saves a changed command as its argument list, and gives saved the status again', async () => {
		await show();
		type(fieldOf('kluster.deploy'), 'scripts/kluster-deploy.sh\n--wait\n');
		writes['/api/settings'] = answer({ commit: 'def5678abc' });
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({
			kind: 'manifest',
			set: { 'kluster.deploy': ['scripts/kluster-deploy.sh', '--wait'] }
		});
		expect(q('kluster-save-deploy')!.textContent).toContain('def5678');
		expect(saved).toHaveBeenCalledWith(STATUS);
	});

	it('says when nothing changed, and unsets an emptied command', async () => {
		await show();
		buttonNamed('Save')!.click();
		await settle();
		expect(writesTo('/api/settings')).toEqual([]);
		expect(q('kluster-save-deploy')!.textContent).toContain('Nothing changed');
		type(fieldOf('kluster.deploy'), '');
		writes['/api/settings'] = answer({});
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({ kind: 'manifest', unset: ['kluster.deploy'] });
	});

	it('shows the components as rows and saves them as the JSON list', async () => {
		await show({ settings: COMPONENTS });
		expect(all('component-row')).toHaveLength(2);
		expect(fieldOf('kluster.namespace')).toBeTruthy();
		expect(q('field-kluster.build')).toBeNull();
		buttonNamed('Add a component')!.click();
		flushSync();
		const names = all('component-name') as HTMLInputElement[];
		const builds = all('component-build') as HTMLTextAreaElement[];
		type(names[2], 'worker');
		type(builds[2], 'scripts/build.sh\nworker');
		writes['/api/settings'] = answer({});
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({
			kind: 'manifest',
			set: {
				'kluster.components': [
					{ name: 'api', build: ['scripts/build.sh', 'api'] },
					{ name: 'web', build: ['scripts/build.sh', 'web'] },
					{ name: 'worker', build: ['scripts/build.sh', 'worker'] }
				]
			}
		});
	});

	it('removes and moves a component', async () => {
		await show({ settings: COMPONENTS });
		[...document.querySelectorAll<HTMLButtonElement>('button')]
			.find((b) => b.getAttribute('aria-label') === 'Move web up')!
			.click();
		flushSync();
		expect((all('component-name') as HTMLInputElement[]).map((i) => i.value)).toEqual([
			'web',
			'api'
		]);
		(all('component-row')[1].querySelector('button:last-of-type') as HTMLButtonElement).click();
		flushSync();
		expect(all('component-row')).toHaveLength(1);
		writes['/api/settings'] = answer({});
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({
			kind: 'manifest',
			set: { 'kluster.components': [{ name: 'web', build: ['scripts/build.sh', 'web'] }] }
		});
	});

	it('switching from one build to components unsets kluster.build in the same save', async () => {
		await show();
		pick('components');
		expect(all('component-row')).toHaveLength(0);
		buttonNamed('Add a component')!.click();
		flushSync();
		type(all('component-name')[0] as HTMLInputElement, 'api');
		type(all('component-build')[0] as HTMLTextAreaElement, 'scripts/build.sh');
		writes['/api/settings'] = answer({});
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({
			kind: 'manifest',
			set: { 'kluster.components': [{ name: 'api', build: ['scripts/build.sh'] }] },
			unset: ['kluster.build']
		});
	});

	it('switching from components to one build unsets kluster.components', async () => {
		await show({ settings: COMPONENTS });
		pick('one');
		type(fieldOf('kluster.build'), 'scripts/build-all.sh');
		writes['/api/settings'] = answer({});
		buttonNamed('Save')!.click();
		await settle();
		expect(sent('/api/settings')).toEqual({
			kind: 'manifest',
			set: { 'kluster.build': ['scripts/build-all.sh'] },
			unset: ['kluster.components']
		});
	});

	it("shows a refused field with flai's reason beside it", async () => {
		await show({ settings: COMPONENTS });
		type(all('component-name')[0] as HTMLInputElement, 'API');
		const reason = 'kluster.components[1].name "API" is not a DNS label';
		writes['/api/settings'] = answer(
			{ error: 'refused', refused: [{ field: 'kluster.components', reason }] },
			422
		);
		buttonNamed('Save')!.click();
		await settle();
		expect(q('refused-kluster.components')!.textContent).toBe(reason);
		expect(q('kluster-save-deploy')!.textContent).toContain('flai refused what is marked');
		expect(saved).not.toHaveBeenCalled();
	});

	it('with settings off says the one host step and saves nothing', async () => {
		await show({ settingsOn: false });
		expect(q('deploy-settings-off')!.textContent).toContain('flai serve enable settings');
		expect(fieldOf('kluster.deploy').disabled).toBe(true);
		expect(buttonNamed('Save')!.disabled).toBe(true);
		expect((q('kluster-build-mode')!.querySelector('input') as HTMLInputElement).disabled).toBe(
			true
		);
	});

	it('leaves the settings notice to a page that says it elsewhere', async () => {
		await show({ settingsOn: false, notice: false });
		expect(q('deploy-settings-off')).toBeNull();
	});

	it('says no story is deployed while the deploy is unset', async () => {
		await show({ deployable: false });
		expect(q('kluster-no-deploy')!.textContent).toContain('no story is deployed');
	});

	it('offers no choice of components to a flai that lists none', async () => {
		await show({ settings: SETTINGS.filter((s) => s.key !== 'kluster.components') });
		expect(q('kluster-build-mode')).toBeNull();
		expect(fieldOf('kluster.build')).toBeTruthy();
	});

	it('makes a backlog story for an agent that names the components, and links it', async () => {
		await show({ settings: COMPONENTS });
		writes['/api/items'] = answer({ item: { id: 'S-0042' } });
		buttonNamed('Write these with an agent')!.click();
		await settle();
		const body = sent('/api/items');
		expect(body).toMatchObject({ type: 'story', parent: 'E-0001', nature: 'feature' });
		expect(String(body.title)).toContain('build, deploy, and teardown scripts');
		expect(String(body.body)).toContain('Its components are api, web');
		expect(q('kluster-agent-made')!.querySelector('a')!.getAttribute('href')).toBe('/items/S-0042');
	});
});
