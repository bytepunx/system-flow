// S-0259: the planner's page loads /api/planner, and loads it again when the planner's activity
// document or a work item changes, when flai serve says a run started or ended, and every 15 seconds
// while a run is under way; a reload that fails says so above what was last shown.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { PlannerView } from '$lib/planner';
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
});
