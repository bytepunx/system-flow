// S-0420: the role page, from a fixture agent.status's roles with a role the project defines,
// reviewer, which reviews a story: its host action on, with what its triggers say, or off with the
// command that turns it on; the run under way with its item and its stream read by the role's name;
// its activity document; its runs with their item, cost, and outcome; Run with an item through
// /api/roles/reviewer, saying what flai answered; Stop and Start for the orchestrator alone; and the
// page's own controls and settings in their places. Then /workflow/<name>, which shows a role flai
// reports and says so for one it does not, and loads again when a run of the role starts or ends.
// ADR-0249: the page draws the role's own cost and use per day from flai stats, in its strategic
// colour, before its log.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { createRawSnippet, flushSync, mount, unmount } from 'svelte';
import type { RoleState, RoleView } from '$lib/roles';
import type { ActivityDocument, RoleRun } from '$lib/strategic';
import type { Change } from '$lib/changes';
import { CATEGORICAL } from '$lib/viz/palette';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace('[id]', params.id ?? '')
}));
const at = vi.hoisted(() => ({ params: { role: 'reviewer' } as Record<string, string> }));
vi.mock('$app/state', () => ({
	page: {
		get params() {
			return at.params;
		}
	}
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

// the role's charts draw on ECharts, whose option is what the test reads
const setOption = vi.fn();
vi.mock('echarts/core', () => ({
	use: vi.fn(),
	init: vi.fn(() => ({ setOption, resize: vi.fn(), dispose: vi.fn() }))
}));
vi.mock('echarts/charts', () => ({ ScatterChart: {}, LineChart: {}, BarChart: {} }));
vi.mock('echarts/components', () => ({
	GridComponent: {},
	TooltipComponent: {},
	LegendComponent: {},
	MarkLineComponent: {}
}));
vi.mock('echarts/renderers', () => ({ CanvasRenderer: {} }));
globalThis.ResizeObserver ??= class {
	observe() {}
	disconnect() {}
	unobserve() {}
} as unknown as typeof ResizeObserver;

import RolePage from './RolePage.svelte';
import RoleRoute from '../../routes/workflow/[role]/+page.svelte';

const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const one = (id: string) => document.querySelector<HTMLElement>(`[data-testid="${id}"]`);
const all = (id: string) => [...document.querySelectorAll<HTMLElement>(`[data-testid="${id}"]`)];
const text = (el: Element | null) => el!.textContent!.replace(/\s+/g, ' ').trim();
const button = (id: string) => one(id) as HTMLButtonElement | null;
function type(el: HTMLInputElement, value: string) {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
}

const reviewing: RoleRun = {
	story: '',
	item: 'S-0001',
	role: 'review',
	role_name: 'reviewer',
	agent: 'reviewer',
	harness: 'claude-code',
	model: 'opus',
	command: 'claude',
	pid: 42,
	started: '2026-10-28T10:00:00Z'
};
const reviewed: RoleRun = {
	...reviewing,
	item: 'S-0002',
	pid: 41,
	started: '2026-10-27T10:00:00Z',
	ended: '2026-10-27T10:20:00Z',
	outcome: 'worked'
};

/** The reviewer as a fixture agent.status's roles carry it, reviewing S-0001 after S-0002. */
const reviewer = (over: Partial<RoleState> = {}): RoleState => ({
	name: 'reviewer',
	role: 'review',
	path: 'design/workflows/roles/reviewer.md',
	item: ['story'],
	trigger: ['asked', 'schedule'],
	action: { name: 'review', enabled: true },
	run: reviewing,
	running: [reviewing],
	runs: [reviewed],
	held: false,
	...over
});
const orchestrator = (held: boolean): RoleState => ({
	name: 'orchestrator',
	role: 'orchestrate',
	path: 'design/workflows/roles/orchestrator.md',
	item: [],
	trigger: ['event'],
	action: { name: 'orchestrate', enabled: true },
	run: null,
	running: [],
	runs: [],
	held
});
const activity = (over: Partial<ActivityDocument> = {}): ActivityDocument => ({
	kind: 'reviewer',
	accrued_cost: 0.9,
	accrued_seconds: 1200,
	tasks_completed: 1,
	last_run: '2026-10-27T10:20:00Z',
	path: 'wip/agents/reviewer.md',
	entries: [
		{
			at: '2026-10-27T10:20:00Z',
			summary: 'reviewed S-0002: two findings',
			items: ['S-0002'],
			seconds: 1200,
			cost: 0.9,
			estimated: false
		}
	],
	...over
});
const stream = {
	story: '',
	agent: 'reviewer',
	running: true,
	from: 0,
	next: 0,
	size: 0,
	entries: []
};

describe('the role page (S-0420)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let changed = 0;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		changed = 0;
		document.body.innerHTML = '';
	});
	const show = async (role: RoleState, props: Record<string, unknown> = {}) => {
		c = mount(RolePage, {
			target: document.body,
			props: { role, activity: activity(), onchanged: () => changed++, ...props }
		});
		await settle();
	};
	/** Answers the stream's reads, and a POST to /api/roles/<name> with posted. */
	const backend = (posted: () => ReturnType<typeof answer>) =>
		api.mockImplementation(async (url: string, init?: RequestInit) =>
			init?.method === 'POST' ? posted() : answer(stream)
		);

	it("shows the reviewer's action on, its run under way with its item and stream, its log, and its runs", async () => {
		backend(() => answer({}));
		await show(reviewer());
		expect(text(one('reviewer-action'))).toContain('The review host action');
		expect(text(one('reviewer-action-on'))).toBe(
			'On: flai serve starts the reviewer when you Run it and when its schedule comes round.'
		);
		expect(one('reviewer-current-item')!.querySelector('a')!.getAttribute('href')).toBe(
			'/items/S-0001'
		);
		expect(text(one('reviewer-current'))).toContain('reviewer (claude-code, opus)');
		expect(api).toHaveBeenCalledWith('/api/agent-stream/reviewer?role=reviewer');
		expect(text(one('reviewer-activity'))).toContain('wip/agents/reviewer.md');
		expect(text(one('reviewer-entry-summary'))).toBe('reviewed S-0002: two findings');
		expect(all('reviewer-run-item').map(text)).toEqual(['S-0001', 'S-0002']);
		expect(all('reviewer-run-outcome').map(text)).toEqual(['running', 'worked']);
		expect(all('reviewer-run-cost').map(text)).toEqual(['—', '$0.900']);
		expect(one('reviewer-stop')).toBeNull();
		expect(one('reviewer-start')).toBeNull();
	});

	it('runs the reviewer on the item typed, saying what flai answered, once one is typed', async () => {
		backend(() =>
			answer({ role: 'review', item: 'S-0003', agent: 'reviewer', command: 'claude', pid: 43 })
		);
		await show(reviewer());
		expect(one('reviewer-item')!.getAttribute('placeholder')).toBe('S-0001');
		expect(button('reviewer-run-now')!.disabled).toBe(true);
		type(one('reviewer-item') as HTMLInputElement, ' S-0003 ');
		expect(button('reviewer-run-now')!.disabled).toBe(false);
		button('reviewer-run-now')!.click();
		await settle();
		const post = api.mock.calls.find(([, init]) => init?.method === 'POST')!;
		expect(post[0]).toBe('/api/roles/reviewer');
		expect(JSON.parse(post[1].body)).toEqual({ action: 'run', item: 'S-0003' });
		expect(text(one('reviewer-said'))).toBe(
			'started claude to run the reviewer on S-0003 as reviewer (pid 43)'
		);
		expect(changed).toBe(1);
	});

	it("says flai's refusal of a Run", async () => {
		backend(() => answer({ error: 'S-0003 is done: the reviewer runs on a story not done' }, 400));
		await show(reviewer());
		type(one('reviewer-item') as HTMLInputElement, 'S-0003');
		button('reviewer-run-now')!.click();
		await settle();
		expect(text(one('reviewer-said'))).toBe(
			'refused: S-0003 is done: the reviewer runs on a story not done'
		);
	});

	it('is off with the command that turns it on, and Run disabled saying why', async () => {
		await show(reviewer({ action: { name: 'review', enabled: false }, run: null, running: [] }));
		expect(text(one('reviewer-action-off'))).toContain(
			'Off: flai serve starts no reviewer for this project. On the host, in the project, run flai serve enable review'
		);
		expect(button('reviewer-run-now')!.disabled).toBe(true);
		expect((one('reviewer-item') as HTMLInputElement).disabled).toBe(true);
		expect(text(one('reviewer-controls-off'))).toContain('flai serve enable review');
	});

	it('runs a role on no item once at a time, and offers no Run to a role the operator cannot ask for', async () => {
		const auditing = { ...reviewing, item: undefined, role_name: 'auditor', agent: 'auditor' };
		await show(
			reviewer({
				name: 'auditor',
				role: 'audit',
				item: [],
				action: { name: 'audit', enabled: true },
				run: auditing,
				running: [auditing],
				runs: []
			})
		);
		expect(one('auditor-item')).toBeNull();
		expect(one('auditor-current-item')).toBeNull();
		expect(button('auditor-run-now')!.disabled).toBe(true);
		expect(one('auditor-run-busy')).not.toBeNull();
		if (c) unmount(c);
		document.body.innerHTML = '';
		await show(reviewer({ trigger: ['event'] }));
		expect(one('reviewer-run-now')).toBeNull();
		expect(one('reviewer-item')).toBeNull();
	});

	it('offers the orchestrator Stop while it is not held, and Start while it is', async () => {
		backend(() => answer({ held: true, pid: 42 }));
		await show(orchestrator(false), { activity: activity({ kind: 'orchestrator', entries: [] }) });
		expect(text(one('orchestrator-action-on'))).toContain('for as long as the action stays on');
		expect(one('orchestrator-run-now')).toBeNull();
		expect(one('orchestrator-start')).toBeNull();
		button('orchestrator-stop')!.click();
		await settle();
		const post = api.mock.calls.find(([, init]) => init?.method === 'POST')!;
		expect(post[0]).toBe('/api/roles/orchestrator');
		expect(JSON.parse(post[1].body)).toEqual({ action: 'stop' });
		expect(text(one('orchestrator-said'))).toBe(
			'stopped the orchestrator orchestrator (pid 42) and held it stopped'
		);
		unmount(c!);
		document.body.innerHTML = '';
		await show(orchestrator(true), { activity: activity({ kind: 'orchestrator', entries: [] }) });
		expect(text(one('orchestrator-action-held'))).toContain('until you Start it');
		expect(one('orchestrator-stop')).toBeNull();
		expect(button('orchestrator-start')!.disabled).toBe(false);
	});

	it("puts a page's own controls in place of Run, and its settings under the host action", async () => {
		const actions = createRawSnippet(() => ({
			render: () => '<p data-testid="own-controls">focus</p>'
		}));
		const settings = createRawSnippet(() => ({
			render: () => '<p data-testid="the-settings">roles.reviewer.agent</p>'
		}));
		await show(reviewer(), { actions, settings });
		expect(one('reviewer-action')!.contains(one('own-controls'))).toBe(true);
		expect(one('reviewer-run-now')).toBeNull();
		expect(one('reviewer-settings')!.contains(one('the-settings'))).toBe(true);
		expect(
			one('reviewer-settings')!.compareDocumentPosition(one('reviewer-current')!) &
				Node.DOCUMENT_POSITION_FOLLOWING
		).toBeTruthy();
	});
});

describe("the role's cost and use per day (ADR-0249)", () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
	});
	// three days of the window: the planner on the first two, the reviewer, a fixture role the
	// project defines, on the last two, estimated in part on the last
	const report = {
		window_start: '2026-10-25T00:00:00Z',
		generated_at: '2026-10-27T12:00:00Z',
		type: 'story',
		items: [],
		strategic_days: [
			{
				date: '2026-10-25',
				agents: { planner: { cost: 1, seconds: 600 } },
				cost: 1,
				seconds: 600,
				completed: 0
			},
			{
				date: '2026-10-26',
				agents: {
					planner: { cost: 0.5, seconds: 300 },
					reviewer: { cost: 0.4, seconds: 1800 }
				},
				cost: 0.9,
				seconds: 2100,
				completed: 1
			},
			{
				date: '2026-10-27',
				agents: { reviewer: { cost: 0.5, seconds: 1800, estimated: true } },
				cost: 0.5,
				seconds: 1800,
				completed: 0
			}
		]
	};
	const stats = (body: unknown, status = 200) =>
		api.mockImplementation(async (url: string) =>
			url.startsWith('/api/stats') ? answer(body, status) : answer(stream)
		);
	const show = async (role: RoleState = reviewer()) => {
		c = mount(RolePage, { target: document.body, props: { role, activity: activity() } });
		await settle();
	};
	type Drawn = {
		series: { name: string; itemStyle: { color: string }; data: { value: [number, number] }[] }[];
		yAxis: { name: string };
	};
	const drawn = () => setOption.mock.calls.map(([o]) => o as Drawn);

	it("draws the reviewer's own cost and hours per day from flai stats, before its log", async () => {
		stats(report);
		await show();
		expect(api).toHaveBeenCalledWith('/api/stats?since=30d&type=story&bucket=day');
		expect(text(one('reviewer-usage-summary'))).toBe(
			'Over the last 30d the reviewer spent $0.900 and worked 1h, on 2 days.'
		);
		// what it spent each day, its own series alone, though the planner spent beside it
		await vi.waitFor(() => expect(setOption).toHaveBeenCalled());
		const cost = drawn().at(-1)!;
		expect(cost.yAxis.name).toBe('US dollars');
		expect(cost.series.map((s) => s.name)).toEqual(['reviewer']);
		expect(cost.series[0].data.map((d) => d.value[1])).toEqual([0, 0.4, 0.5]);
		// in the colour Strategic Cost gives it beside the planner: the first slot after the built-ins'
		expect(cost.series[0].itemStyle.color).toBe(CATEGORICAL.light[2]);
		expect(button('reviewer-usage-cost')!.getAttribute('aria-pressed')).toBe('true');
		// and the hours it worked
		button('reviewer-usage-use')!.click();
		await settle();
		const use = drawn().at(-1)!;
		expect(use.yAxis.name).toBe('hours');
		expect(use.series.map((s) => s.name)).toEqual(['reviewer']);
		expect(use.series[0].data.map((d) => d.value[1])).toEqual([0, 0.5, 0.5]);
		expect(use.series[0].itemStyle.color).toBe(CATEGORICAL.light[2]);
		expect(button('reviewer-usage-use')!.getAttribute('aria-pressed')).toBe('true');
		expect(
			one('reviewer-usage')!.compareDocumentPosition(one('reviewer-activity')!) &
				Node.DOCUMENT_POSITION_FOLLOWING
		).toBeTruthy();
		expect(one('reviewer-usage-error')).toBeNull();
	});

	it('says the role has no activity in the window, an older flai, or what flai stats refused', async () => {
		stats(report);
		await show(reviewer({ name: 'auditor' }));
		expect(text(one('auditor-usage-none'))).toBe(
			'No activity of the auditor ended in the last 30d.'
		);
		expect(setOption).not.toHaveBeenCalled();
		unmount(c!);
		document.body.innerHTML = '';
		stats({ ...report, strategic_days: undefined });
		await show();
		expect(text(one('reviewer-usage-older'))).toContain('flai self-upgrade');
		unmount(c!);
		document.body.innerHTML = '';
		stats({ error: 'flai is away' }, 502);
		await show();
		expect(text(one('reviewer-usage-error'))).toBe('flai stats could not be read: flai is away');
	});
});

describe('/workflow/<name> (S-0420)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		events.followed.length = 0;
		events.agents.length = 0;
		at.params = { role: 'reviewer' };
		document.body.innerHTML = '';
	});
	const loads = () =>
		api.mock.calls.filter(([url, init]) => url === '/api/roles/reviewer' && !init?.method).length;
	const routes = (view: RoleView | null, status = 200) =>
		api.mockImplementation(async (url: string) =>
			url.startsWith('/api/roles/')
				? answer(
						view ?? {
							error:
								'the project runs no role named auditor; it runs analyzer, orchestrator, planner, reviewer'
						},
						status
					)
				: answer(stream)
		);
	const show = async () => {
		c = mount(RoleRoute, { target: document.body });
		await settle();
	};

	it('shows the role agent.status reports by the name', async () => {
		routes({ role: reviewer(), activity: activity() });
		await show();
		expect(loads()).toBe(1);
		expect(document.querySelector('h1')!.textContent).toBe('Reviewer');
		expect(document.body.textContent).toContain('design/workflows/roles/reviewer.md');
		expect(one('reviewer-action-on')).not.toBeNull();
		expect(one('role-not-found')).toBeNull();
	});

	it('says so for a role flai does not report', async () => {
		at.params = { role: 'auditor' };
		routes(null, 404);
		await show();
		expect(text(one('role-not-found'))).toBe(
			'the project runs no role named auditor; it runs analyzer, orchestrator, planner, reviewer'
		);
		expect(one('auditor-action')).toBeNull();
	});

	it('says a failed read above what it last showed', async () => {
		routes({ role: reviewer(), activity: activity() });
		await show();
		routes({ error: 'flai is away' } as never, 502);
		for (const agent of events.agents) agent('reviewer');
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(one('reviewer-action')).not.toBeNull();
	});

	it("loads again when a run of the role starts or ends, or its document changes, and not for another's", async () => {
		routes({ role: reviewer(), activity: activity() });
		await show();
		for (const agent of events.agents) agent('analyzer');
		await settle();
		expect(loads()).toBe(1);
		for (const agent of events.agents) agent('S-0003');
		await settle();
		expect(loads()).toBe(2);
		for (const x of events.followed)
			if (x.kinds.includes('narrative')) x.f([{ path: 'wip/agents/S-0001.md', kind: 'narrative' }]);
		await settle();
		expect(loads()).toBe(2);
		for (const x of events.followed)
			if (x.kinds.includes('narrative'))
				x.f([{ path: 'wip/agents/reviewer.md', kind: 'narrative' }]);
		await settle();
		expect(loads()).toBe(3);
	});
});
