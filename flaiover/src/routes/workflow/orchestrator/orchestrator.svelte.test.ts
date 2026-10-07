// S-0228: the orchestrator's page loads /api/orchestrator, and loads it again when the
// orchestrator's activity document changes and when flai serve says an orchestrator run started,
// ended, or was held; a reload that fails says so above what was last shown.
// S-0229: it shows the orchestration block of settings.get in its settings panel, saves a key
// through /api/settings and reads the settings again, is read-only while the settings host action is
// off, and reads the settings again when system-flow.yaml changes.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { SettingsView, StrategicSetting } from '$lib/settings';
import type { OrchestratorView } from '$lib/strategic';
import type { Change, ChangeKind } from '$lib/changes';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id ?? '')
}));
// what the page follows, to deliver change and agent events to as flai serve would
const events = vi.hoisted(() => ({
	followed: [] as { kinds: string[]; f: (changes: Change[]) => void }[],
	agents: [] as ((id: string) => void)[]
}));
vi.mock('$lib/events', () => ({
	follow: (kinds: string[], f: (changes: Change[]) => void) => {
		events.followed.push({ kinds, f });
		return () => {};
	},
	listen: (l: { agent?: (id: string) => void }) => {
		if (l.agent) events.agents.push(l.agent);
		return () => {};
	},
	debounced: (f: () => void) => Object.assign(() => f(), { stop: () => {} })
}));

import OrchestratorPage from './+page.svelte';

const settle = async () => {
	for (let i = 0; i < 8; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const q = <T extends HTMLElement = HTMLElement>(id: string) =>
	document.querySelector<T>(`[data-testid="${id}"]`);
const answer = (status: number, body: unknown) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const changed = (path: string, kind: ChangeKind) => {
	for (const x of events.followed) if (x.kinds.includes(kind)) x.f([{ path, kind }]);
};

const permission = (name: string, value?: boolean): StrategicSetting => ({
	key: `orchestration.permissions.${name}`,
	kind: 'boolean',
	values: ['true', 'false'],
	default: 'false',
	meaning: `Lets the orchestrator ${name}.`,
	risk: `What can go wrong with ${name}.`,
	...(value === undefined ? { set: false } : { value, set: true })
});
const ORCHESTRATION = [
	'orchestration.permissions.plan_backlog_epics',
	'orchestration.permissions.plan_backlog_stories',
	'orchestration.permissions.finalize_drafts',
	'orchestration.permissions.promote_to_ready',
	'orchestration.permissions.order_ready',
	'orchestration.permissions.answer_threads',
	'orchestration.permissions.accept_reviews',
	'orchestration.permissions.publish',
	'orchestration.policy',
	'orchestration.release.policy',
	'orchestration.release.value',
	'orchestration.release.count',
	'orchestration.release.epic',
	'orchestration.release.tag',
	'orchestration.release.whole_epics'
];
/** A settings.get answer with the orchestration block as flai lists it, and a key of another block. */
function settingsGet(editable: boolean): SettingsView {
	return {
		here: editable,
		everywhere: false,
		enable: 'flai serve enable settings',
		enable_everywhere: 'flai serve enable settings --everywhere',
		host: {
			actions: [],
			agent: { name: 'agent', command: [], harnesses: {} },
			checks: { commands: [], timeout_minutes: 10 },
			import_roots: [],
			strategic: {
				editable,
				enable: 'flai serve enable settings',
				settings: [
					permission('plan_backlog_epics'),
					permission('plan_backlog_stories'),
					permission('finalize_drafts'),
					permission('promote_to_ready', true),
					permission('order_ready'),
					{
						key: 'orchestration.permissions.answer_threads',
						kind: 'choice',
						values: ['off', 'recommend', 'autonomous'],
						default: 'off',
						meaning: 'How the orchestrator answers threads.',
						risk: 'With autonomous, agents act on answers you did not give.',
						set: false
					},
					permission('accept_reviews'),
					permission('publish'),
					{
						key: 'orchestration.policy',
						kind: 'choice',
						values: ['cod', 'wsjf', 'throughput', 'fifo'],
						default: 'fifo',
						meaning: 'How the ready column is ordered.',
						value: 'cod',
						set: true
					},
					{
						key: 'orchestration.release.policy',
						kind: 'choice',
						values: ['judgement', 'threshold', 'theme'],
						default: 'judgement',
						meaning: 'When accepted stories not yet released are due a release.',
						set: false
					},
					{
						key: 'orchestration.release.value',
						kind: 'number',
						meaning: 'Under threshold, the unreleased cost of delay per week.',
						set: false
					},
					{
						key: 'orchestration.release.count',
						kind: 'number',
						whole: true,
						meaning: 'Under threshold, the number of accepted stories not yet released.',
						value: 3,
						set: true
					},
					{
						key: 'orchestration.release.epic',
						kind: 'text',
						meaning: 'Under theme, the epic.',
						set: false
					},
					{
						key: 'orchestration.release.tag',
						kind: 'text',
						meaning: 'Under theme, the tag.',
						set: false
					},
					{
						key: 'orchestration.release.whole_epics',
						kind: 'boolean',
						values: ['true', 'false'],
						default: 'false',
						meaning: 'Holds a release back while an epic is open.',
						set: false
					},
					{
						key: 'planning.schedule',
						kind: 'cron',
						meaning: 'When flai serve runs the planner.',
						set: false
					}
				]
			}
		}
	};
}
const reads = () =>
	api.mock.calls.filter(([url, init]) => url === '/api/settings' && !init?.method).length;
const loads = () =>
	api.mock.calls.filter(([url, init]) => url === '/api/orchestrator' && !init?.method).length;
const view = (over: Partial<OrchestratorView> = {}): OrchestratorView => ({
	enabled: true,
	held: false,
	activity: {
		kind: 'orchestrator',
		accrued_cost: 0,
		accrued_seconds: 0,
		tasks_completed: 0,
		last_run: '',
		path: 'wip/agents/orchestrator.md',
		entries: []
	},
	run: null,
	runs: [],
	...over
});
/** Answers /api/orchestrator with o, and every other request with the settings. */
const backend = (o: OrchestratorView, settings: SettingsView = settingsGet(true)) =>
	api.mockImplementation(async (url: string) =>
		answer(200, url === '/api/orchestrator' ? o : settings)
	);

describe('the orchestrator page (S-0228)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		events.followed.length = 0;
		events.agents.length = 0;
		document.body.innerHTML = '';
	});
	const show = async () => {
		c = mount(OrchestratorPage, { target: document.body });
		await settle();
	};

	it('loads the orchestrator and shows it above its settings', async () => {
		backend(view({ held: true }));
		await show();
		expect(api).toHaveBeenCalledWith('/api/orchestrator');
		expect(document.querySelector('h1')!.textContent).toBe('Orchestrator');
		expect(document.body.textContent).toContain('wip/agents/orchestrator.md');
		expect(document.body.textContent).not.toContain('will be shown here');
		expect(q('orchestrator-action-held')).not.toBeNull();
		expect(q('orchestrator-current-none')).not.toBeNull();
		const settings = q('strategic-orchestration')!;
		expect(
			q('orchestrator-runs')!.compareDocumentPosition(settings) & Node.DOCUMENT_POSITION_FOLLOWING
		).toBeTruthy();
	});

	it("loads again when the orchestrator's document changes, and not for another narrative", async () => {
		backend(view());
		await show();
		expect(loads()).toBe(1);
		changed('wip/agents/S-0001.md', 'narrative');
		await settle();
		expect(loads()).toBe(1);

		const decision = 'promoted S-0252: cost of delay 60.81 USD a week, first under policy cod';
		backend(
			view({
				activity: {
					...view().activity,
					tasks_completed: 1,
					entries: [
						{
							at: '2026-10-03T10:10:00Z',
							summary: decision,
							items: ['S-0252'],
							seconds: 300,
							cost: 0.4,
							estimated: false
						}
					]
				}
			})
		);
		changed('wip/agents/orchestrator.md', 'narrative');
		await settle();
		expect(loads()).toBe(2);
		expect(q('orchestrator-decision-summary')!.textContent).toBe(decision);
	});

	it('loads again when flai serve says an orchestrator run started, ended, or was held, and not for a story', async () => {
		backend(view());
		await show();
		events.agents.forEach((f) => f('S-0228'));
		await settle();
		expect(loads()).toBe(1);

		backend(view({ held: true }));
		events.agents.forEach((f) => f('orchestrator'));
		await settle();
		expect(loads()).toBe(2);
		expect(q('orchestrator-action-held')).not.toBeNull();
	});

	it('loads again after a Stop', async () => {
		const running = {
			story: '',
			agent: 'orchestrator',
			command: 'claude',
			pid: 42,
			started: '2026-10-03T10:00:00Z'
		};
		api.mockImplementation(async (url: string, init?: { method?: string }) => {
			if (url === '/api/orchestrator' && init?.method === 'POST')
				return answer(200, { held: true, orchestrator: { ...running, held: true } });
			if (url === '/api/orchestrator') return answer(200, view({ run: running, runs: [running] }));
			if (url === '/api/settings') return answer(200, settingsGet(true));
			return answer(200, { entries: [], running: true, from: 0, next: 0, size: 0 });
		});
		await show();
		q<HTMLButtonElement>('orchestrator-stop')!.click();
		await settle();
		expect(q('orchestrator-said')!.textContent).toContain('held it stopped');
		expect(loads()).toBe(2);
	});

	it('keeps what it showed when a reload fails, with the error above it', async () => {
		backend(view());
		await show();
		api.mockImplementation(async (url: string) =>
			url === '/api/orchestrator'
				? answer(502, { error: 'flai is away' })
				: answer(200, settingsGet(true))
		);
		changed('wip/agents/orchestrator.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(q('orchestrator-action')).not.toBeNull();

		backend(view());
		changed('wip/agents/orchestrator.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')).toBeNull();
	});

	// S-0229: the orchestrator's settings, read apart from the orchestrator on arrival, after a save,
	// and when system-flow.yaml changes.
	it("shows the orchestration block's keys, each permission with its risk", async () => {
		backend(view());
		await show();
		expect(api).toHaveBeenCalledWith('/api/settings');
		const rows = [...document.querySelectorAll('[data-testid^="row-"]')].map((r) =>
			r.getAttribute('data-testid')!.slice('row-'.length)
		);
		expect(rows).toEqual(ORCHESTRATION);
		for (const key of ORCHESTRATION.filter((k) => k.startsWith('orchestration.permissions.')))
			expect(q(`risk-${key}`)).not.toBeNull();
		expect(q<HTMLInputElement>('input-orchestration.permissions.promote_to_ready')!.checked).toBe(
			true
		);
		expect(q('strategic-readonly')).toBeNull();
	});

	it('saves a changed key, then reads the settings again', async () => {
		const sent: unknown[] = [];
		api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
			if (url === '/api/orchestrator') return answer(200, view());
			if (init?.method === 'POST') {
				sent.push(JSON.parse(init.body!));
				return answer(200, { set: [], unset: [], commit: 'abcdef123456', warnings: [] });
			}
			return answer(200, settingsGet(true));
		});
		await show();
		expect(reads()).toBe(1);
		q<HTMLInputElement>('input-orchestration.permissions.publish')!.click();
		flushSync();
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent).toEqual([
			{ kind: 'manifest', set: { 'orchestration.permissions.publish': true } }
		]);
		expect(reads()).toBe(2);
	});

	it('is read-only, with the reason, while the settings host action is off', async () => {
		backend(view(), settingsGet(false));
		await show();
		expect(q('strategic-readonly')!.textContent).toContain('settings host action is off');
		const inputs = document.querySelectorAll<HTMLInputElement>('[data-testid^="input-"]');
		expect(inputs.length).toBeGreaterThan(0);
		for (const el of inputs) expect(el.disabled).toBe(true);
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(true);
	});

	it('reads the settings again when system-flow.yaml changes, and not for another file', async () => {
		backend(view());
		await show();
		changed('wip/kanban/stories/S-0001-x.md', 'item');
		await settle();
		expect(reads()).toBe(1);
		changed('system-flow.yaml', 'project');
		await settle();
		expect(reads()).toBe(2);
	});

	it('says so when flai gives no settings, and shows an error it gives', async () => {
		const old = settingsGet(true);
		delete old.host!.strategic;
		backend(view(), old);
		await show();
		expect(q('strategic-absent')).not.toBeNull();
		expect(q('strategic-orchestration')).toBeNull();

		unmount(c!);
		document.body.innerHTML = '';
		api.mockResolvedValue(answer(502, { error: 'flai is away' }));
		await show();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(document.body.textContent).not.toContain('Loading…');
	});
});
