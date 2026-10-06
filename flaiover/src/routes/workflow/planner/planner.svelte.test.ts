// S-0259: the planner's page loads /api/planner, and loads it again when the planner's activity
// document or a work item changes, when flai serve says a run started or ended, and every 15 seconds
// while a run is under way; a reload that fails says so above what was last shown.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { PlannerView } from '$lib/planner';
import type { SettingsView } from '$lib/settings';
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

import PlannerPage from './+page.svelte';

const settle = async () => {
	for (let i = 0; i < 10; i++) await Promise.resolve();
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, status: 200, json: async () => body });
const refused = (error: string) => ({
	ok: false,
	status: 500,
	statusText: 'Internal Server Error',
	json: async () => ({ error })
});
const changed = (path: string, kind: ChangeKind) => {
	for (const x of events.followed) if (x.kinds.includes(kind)) x.f([{ path, kind }]);
};
const view = (over: Partial<PlannerView> = {}): PlannerView => ({
	plan_enabled: true,
	activity: {
		kind: 'planner',
		accrued_cost: 0,
		accrued_seconds: 0,
		tasks_completed: 0,
		last_run: '',
		path: 'wip/agents/planner.md',
		entries: []
	},
	runs: [],
	...over
});
const loads = () => api.mock.calls.filter(([url]) => url === '/api/planner').length;
const settingsReads = () =>
	api.mock.calls.filter(([url, init]) => url === '/api/settings' && !init?.method).length;
/** A settings.get answer with the planning block as flai lists it, and a key of another block. */
const settingsGet = (): SettingsView => ({
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
					key: 'planning.agent',
					kind: 'agent',
					meaning: "The planner's agent, over the project's.",
					set: false
				},
				{
					key: 'planning.replan',
					kind: 'choice',
					values: ['never', 'deterministic', 'agent'],
					default: 'deterministic',
					meaning: 'What flai serve does when the pull order changes.',
					set: false
				},
				{
					key: 'planning.schedule',
					kind: 'cron',
					meaning: 'When flai serve runs the planner over the ready column.',
					value: 'daily',
					set: true
				},
				{
					key: 'planning.hour_rate',
					kind: 'number',
					meaning: 'What an hour of work costs.',
					value: 80,
					set: true
				},
				{
					key: 'planning.cycle',
					kind: 'duration',
					default: '168h',
					meaning: "The period a cost of delay's time lost is counted over.",
					set: false
				},
				{
					key: 'planning.default_duration',
					kind: 'duration',
					default: '24h',
					meaning: 'The work a story is forecast to take with too little history.',
					set: false
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
});

describe('the planner page (S-0259)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		events.followed.length = 0;
		events.agents.length = 0;
		vi.useRealTimers();
		document.body.innerHTML = '';
	});
	const show = async () => {
		c = mount(PlannerPage, { target: document.body });
		await settle();
	};

	it('loads the planner and shows it', async () => {
		api.mockResolvedValue(answer(view()));
		await show();
		expect(api).toHaveBeenCalledWith('/api/planner');
		expect(document.querySelector('h1')!.textContent).toBe('Planner');
		expect(document.body.textContent).toContain('wip/agents/planner.md');
		expect(document.querySelector('[data-testid="planner-current-none"]')).not.toBeNull();
	});

	it("loads again when the planner's document changes, and not for another narrative", async () => {
		api.mockResolvedValue(answer(view()));
		await show();
		expect(loads()).toBe(1);
		changed('wip/agents/S-0001.md', 'narrative');
		await settle();
		expect(loads()).toBe(1);

		api.mockResolvedValue(answer(view({ activity: { ...view().activity, tasks_completed: 3 } })));
		changed('wip/agents/planner.md', 'narrative');
		await settle();
		expect(loads()).toBe(2);
		expect(document.querySelector('[data-testid="planner-activities"]')!.textContent!.trim()).toBe(
			'3'
		);
	});

	it('loads again when a work item changes and when a run starts or ends', async () => {
		api.mockResolvedValue(answer(view()));
		await show();
		changed('wip/kanban/stories/S-0259-x.md', 'item');
		await settle();
		expect(loads()).toBe(2);
		events.agents.forEach((f) => f('E-0016'));
		await settle();
		expect(loads()).toBe(3);
	});

	it('loads again every 15 seconds while a run is under way, and stops when none is', async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		const run = {
			item: 'E-0016',
			agent: 'planner',
			command: 'claude',
			started: '2026-10-03T10:00:00Z'
		};
		api.mockImplementation(async (url: string) =>
			url === '/api/planner'
				? answer(view({ runs: [run] }))
				: answer({ entries: [], running: true, from: 0, next: 0, size: 0 })
		);
		await show();
		expect(loads()).toBe(1);
		await vi.advanceTimersByTimeAsync(15000);
		await settle();
		expect(loads()).toBe(2);

		api.mockImplementation(async (url: string) =>
			url === '/api/planner'
				? answer(view({ runs: [{ ...run, ended: '2026-10-03T10:05:00Z', outcome: 'worked' }] }))
				: answer({ entries: [], running: false, from: 0, next: 0, size: 0 })
		);
		await vi.advanceTimersByTimeAsync(15000);
		await settle();
		expect(loads()).toBe(3);
		await vi.advanceTimersByTimeAsync(30000);
		await settle();
		expect(loads()).toBe(3);
	});

	it('keeps what it showed when a reload fails, with the error above it', async () => {
		api.mockResolvedValue(answer(view()));
		await show();
		api.mockResolvedValue(refused('flai is away'));
		changed('wip/agents/planner.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(document.querySelector('[data-testid="planner-action"]')).not.toBeNull();

		api.mockResolvedValue(answer(view()));
		changed('wip/agents/planner.md', 'narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')).toBeNull();
	});

	it('says only the error when nothing was ever loaded', async () => {
		api.mockResolvedValue(refused('flai is away'));
		await show();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(document.body.textContent).not.toContain('Loading…');
		expect(document.querySelector('[data-testid="planner-action"]')).toBeNull();
	});

	// S-0229: the planner's settings, read apart from the planner on arrival, after a save, and when
	// system-flow.yaml changes, never on the 15-second timer.
	describe('its settings (S-0229)', () => {
		const backend = (planner: PlannerView, sent: unknown[] = []) =>
			api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
				if (url === '/api/planner') return answer(planner);
				if (url === '/api/settings' && init?.method === 'POST') {
					sent.push(JSON.parse(init.body!));
					return answer({ set: [], unset: [], warnings: [] });
				}
				if (url === '/api/settings') return answer(settingsGet());
				return answer({ entries: [], running: true, from: 0, next: 0, size: 0 });
			});
		const q = <T extends HTMLElement = HTMLElement>(id: string) =>
			document.querySelector<T>(`[data-testid="${id}"]`);

		it("shows the planning block's keys in its settings panel", async () => {
			backend(view());
			await show();
			expect(settingsReads()).toBe(1);
			const rows = [...document.querySelectorAll('[data-testid^="row-"]')].map((r) =>
				r.getAttribute('data-testid')!.slice('row-'.length)
			);
			expect(rows).toEqual([
				'planning.agent',
				'planning.replan',
				'planning.schedule',
				'planning.hour_rate',
				'planning.cycle',
				'planning.default_duration'
			]);
			expect(q<HTMLInputElement>('input-planning.schedule')!.value).toBe('daily');
		});

		it('saves a changed key, then reads the settings again', async () => {
			const sent: unknown[] = [];
			backend(view(), sent);
			await show();
			const rate = q<HTMLInputElement>('input-planning.hour_rate')!;
			rate.value = '95';
			rate.dispatchEvent(new Event('input', { bubbles: true }));
			flushSync();
			q<HTMLButtonElement>('strategic-save')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'manifest', set: { 'planning.hour_rate': 95 } }]);
			expect(settingsReads()).toBe(2);
		});

		it('reads them again when system-flow.yaml changes, and not on the timer of a run', async () => {
			vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
			backend(
				view({
					runs: [
						{
							item: 'E-0016',
							agent: 'planner',
							command: 'claude',
							started: '2026-10-03T10:00:00Z'
						}
					]
				})
			);
			await show();
			await vi.advanceTimersByTimeAsync(15000);
			await settle();
			expect(loads()).toBe(2);
			expect(settingsReads()).toBe(1);
			changed('system-flow.yaml', 'project');
			await settle();
			expect(settingsReads()).toBe(2);
		});

		it('says so when flai gives no settings', async () => {
			const old = settingsGet();
			delete old.host!.strategic;
			api.mockImplementation(async (url: string) => answer(url === '/api/settings' ? old : view()));
			await show();
			expect(q('strategic-absent')).not.toBeNull();
			expect(q('strategic-planning')).toBeNull();
		});
	});
});
