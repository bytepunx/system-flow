// S-0228: the analyzer's page loads /api/analyzer, and loads it again when the analyzer's activity
// document changes, when flai serve says an analyzer run started or ended, and after a Run; a
// reload that fails says so above what was last shown.
// S-0229: it shows the analysis block of settings.get in its settings panel, and saves a key
// through /api/settings and reads the settings again.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { SettingsView } from '$lib/settings';
import type { AnalyzerView } from '$lib/strategic';
import type { Change, ChangeKind } from '$lib/changes';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[id]', params.id ?? '').replace('[...path]', params.path ?? '')
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

import AnalyzerPage from './+page.svelte';

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
function type(el: HTMLInputElement, value: string) {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
}
const changed = (path: string, kind: ChangeKind) => {
	for (const x of events.followed) if (x.kinds.includes(kind)) x.f([{ path, kind }]);
};

/** A settings.get answer with the analysis block as flai lists it, and a key of another block. */
function settingsGet(): SettingsView {
	return {
		here: true,
		everywhere: false,
		enable: 'flai serve enable settings',
		enable_everywhere: 'flai serve enable settings --everywhere',
		host: {
			actions: [],
			agent: { name: 'agent', command: [], harnesses: {} },
			checks: { commands: [], timeout_minutes: 10 },
			import_roots: [],
			strategic: {
				editable: true,
				enable: 'flai serve enable settings',
				settings: [
					{
						key: 'orchestration.policy',
						kind: 'choice',
						values: ['cod', 'wsjf', 'throughput', 'fifo'],
						default: 'fifo',
						meaning: 'How the ready column is ordered.',
						set: false
					},
					{
						key: 'analysis.agent',
						kind: 'agent',
						meaning: "The analyzer's agent, over the project's.",
						value: { model: 'claude-sonnet-5' },
						set: true
					},
					{
						key: 'analysis.schedule',
						kind: 'cron',
						meaning: 'When flai serve runs the analyzer.',
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
	api.mock.calls.filter(([url, init]) => url === '/api/analyzer' && !init?.method).length;
const view = (over: Partial<AnalyzerView> = {}): AnalyzerView => ({
	enabled: true,
	activity: {
		kind: 'analyzer',
		accrued_cost: 0,
		accrued_seconds: 0,
		tasks_completed: 0,
		last_run: '',
		path: 'wip/agents/analyzer.md',
		entries: []
	},
	run: null,
	runs: [],
	...over
});
/** Answers /api/analyzer with a, and every other request with the settings. */
const backend = (a: AnalyzerView) =>
	api.mockImplementation(async (url: string) =>
		answer(200, url === '/api/analyzer' ? a : settingsGet())
	);

describe('the analyzer page (S-0228)', () => {
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
		c = mount(AnalyzerPage, { target: document.body });
		await settle();
	};

	it('loads the analyzer and shows it above its settings', async () => {
		backend(view({ enabled: false }));
		await show();
		expect(api).toHaveBeenCalledWith('/api/analyzer');
		expect(document.querySelector('h1')!.textContent).toBe('Analyzer');
		expect(document.body.textContent).toContain('wip/agents/analyzer.md');
		expect(document.body.textContent).not.toContain('will be shown here');
		expect(q('analyzer-action-off')).not.toBeNull();
		expect(q('analyzer-current-none')).not.toBeNull();
		const settings = q('strategic-analysis')!;
		expect(
			q('analyzer-runs')!.compareDocumentPosition(settings) & Node.DOCUMENT_POSITION_FOLLOWING
		).toBeTruthy();
	});

	it("loads again when the analyzer's document changes, and not for another narrative", async () => {
		backend(view());
		await show();
		expect(loads()).toBe(1);
		changed('wip/agents/S-0001.md', 'narrative');
		await settle();
		expect(loads()).toBe(1);

		const found = 'Two risks (report design/analysis/2026-10-06-risk.md)';
		backend(
			view({
				activity: {
					...view().activity,
					tasks_completed: 1,
					entries: [
						{
							at: '2026-10-06T10:10:00Z',
							summary: found,
							trigger: 'asked',
							items: [],
							seconds: 600,
							cost: 1.1,
							estimated: false
						}
					]
				}
			})
		);
		changed('wip/agents/analyzer.md', 'narrative');
		await settle();
		expect(loads()).toBe(2);
		expect(q('analyzer-entry-summary')!.textContent).toBe(found);
		expect(q('analyzer-entry-report')!.querySelector('a')!.getAttribute('href')).toBe(
			'/docs/design/analysis/2026-10-06-risk.md'
		);
	});

	it('loads again when flai serve says an analyzer run started or ended, and not for another agent', async () => {
		backend(view());
		await show();
		for (const id of ['S-0228', 'orchestrator']) events.agents.forEach((f) => f(id));
		await settle();
		expect(loads()).toBe(1);

		const running = {
			story: '',
			agent: 'analyzer',
			command: 'claude',
			pid: 51,
			started: '2026-10-06T10:00:00Z',
			focus: 'risk'
		};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/analyzer') return answer(200, view({ run: running, runs: [running] }));
			if (url === '/api/settings') return answer(200, settingsGet());
			return answer(200, { entries: [], running: true, from: 0, next: 0, size: 0 });
		});
		events.agents.forEach((f) => f('analyzer'));
		await settle();
		expect(loads()).toBe(2);
		expect(q('analyzer-current-focus')!.textContent).toBe('risk');
		expect(q<HTMLButtonElement>('analyzer-run')!.disabled).toBe(true);
	});

	it('loads again after a Run', async () => {
		api.mockImplementation(async (url: string, init?: { method?: string }) => {
			if (url === '/api/analyzer' && init?.method === 'POST')
				return answer(200, {
					focus: 'all',
					agent: 'analyzer',
					command: 'claude',
					pid: 88,
					started: '2026-10-06T10:00:00Z'
				});
			if (url === '/api/analyzer') return answer(200, view());
			return answer(200, settingsGet());
		});
		await show();
		q<HTMLButtonElement>('analyzer-run')!.click();
		await settle();
		expect(q('analyzer-said')!.textContent).toBe(
			'started claude to analyze, focus all, as analyzer (pid 88)'
		);
		expect(loads()).toBe(2);
	});

	it('keeps what it showed when a reload fails, with the error above it', async () => {
		backend(view());
		await show();
		api.mockImplementation(async (url: string) =>
			url === '/api/analyzer' ? answer(502, { error: 'flai is away' }) : answer(200, settingsGet())
		);
		changed('wip/agents/analyzer.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(q('analyzer-action')).not.toBeNull();

		backend(view());
		changed('wip/agents/analyzer.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')).toBeNull();
	});

	// S-0229: the analyzer's settings, read apart from the analyzer on arrival and after a save.
	it("shows the analysis block's keys", async () => {
		backend(view());
		await show();
		expect(api).toHaveBeenCalledWith('/api/settings');
		const rows = [...document.querySelectorAll('[data-testid^="row-"]')].map((r) =>
			r.getAttribute('data-testid')!.slice('row-'.length)
		);
		expect(rows).toEqual(['analysis.agent', 'analysis.schedule']);
		expect(q<HTMLInputElement>('input-analysis.agent-model')!.value).toBe('claude-sonnet-5');
	});

	it('saves a changed key, then reads the settings again', async () => {
		const sent: unknown[] = [];
		api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
			if (url === '/api/analyzer') return answer(200, view());
			if (init?.method === 'POST') {
				sent.push(JSON.parse(init.body!));
				return answer(200, { set: [], unset: [], warnings: [] });
			}
			return answer(200, settingsGet());
		});
		await show();
		type(q<HTMLInputElement>('input-analysis.schedule')!, 'daily');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent).toEqual([{ kind: 'manifest', set: { 'analysis.schedule': 'daily' } }]);
		expect(reads()).toBe(2);
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
});
