// S-0163: the charts page offers the charts of spend over time, each with the controls it uses
// and a table of what it plots.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const at = vi.hoisted(() => ({ params: { kind: 'tokens-spent' }, moved: () => {} }));
// SvelteKit keeps the page when it goes to another chart: its params are read as they change
vi.mock('$app/state', async () => {
	const { createSubscriber } = await import('svelte/reactivity');
	const follows = createSubscriber((update) => {
		at.moved = update;
		return () => {};
	});
	return {
		page: {
			get params() {
				follows();
				return at.params;
			}
		}
	};
});
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
vi.mock('$lib/events', () => ({ follow: () => () => {} }));
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

import ChartsPage from './+page.svelte';
import { chartWindow } from '$lib/chartwindow.svelte';

const opus = (items: number, tokens: number, cost: number, seconds: number) => ({
	model: 'claude-opus-5-5',
	items,
	tokens,
	cost,
	seconds,
	tokens_per_item: tokens / items,
	cost_per_item: cost / items,
	tokens_per_minute: tokens / (seconds / 60),
	tokens_per_dollar: tokens / cost
});
const bucketOf = (when: string, items: number, tokens: number, cost: number, seconds: number) => {
	const { model, ...spend } = opus(items, tokens, cost, seconds);
	void model;
	return {
		at: when,
		...spend,
		mean_tokens: tokens,
		mean_cost: cost,
		models: [opus(items, tokens, cost, seconds)]
	};
};
const typeOf = (buckets: ReturnType<typeof bucketOf>[]) => {
	const sum = (k: 'items' | 'tokens' | 'cost' | 'seconds') => buckets.reduce((n, b) => n + b[k], 0);
	const { model, ...spend } = opus(
		sum('items') || 1,
		sum('tokens'),
		sum('cost') || 1,
		sum('seconds') || 1
	);
	void model;
	return { ...spend, items: sum('items'), models: [], buckets };
};
const summary = {
	completed: 3,
	cancelled: 0,
	cancellation_rate: 0,
	throughput_per_week: 1,
	wip: 1,
	cycle_time: { count: 0, p50_seconds: 0, p85_seconds: 0, max_seconds: 0, mean_seconds: 0 },
	lead_time: { count: 0, p50_seconds: 0, p85_seconds: 0, max_seconds: 0, mean_seconds: 0 },
	queue_time: { count: 0, p50_seconds: 0, p85_seconds: 0, max_seconds: 0, mean_seconds: 0 },
	flow_efficiency: 1,
	time_in_state_share: {}
};
/** A report as flai answers it, in the bucket asked for. */
const reportIn = (bucket: string) => ({
	generated_at: '2026-09-29T21:00:00Z',
	type: 'story',
	window_days: 30,
	window_start: '2026-08-30T21:00:00Z',
	summary,
	items: [{ id: 'S-0001', usage: { models: [] } }],
	throughput: [],
	burnup: {},
	cfd: [],
	usage: {
		items: 3,
		tokens: 36000000,
		cost: 15,
		seconds: 1800,
		models: [{ model: 'claude-opus-5-5', tokens: 36000000, cost: 15, items: 3 }],
		bucket,
		spend: {
			epic: typeOf([]),
			story: typeOf([
				bucketOf(
					bucket === 'hour' ? '2026-09-29T19:00:00Z' : '2026-09-29T00:00:00Z',
					3,
					36e6,
					15,
					1800
				)
			]),
			task: typeOf([bucketOf('2026-09-29T00:00:00Z', 8, 24e6, 10, 1200)])
		}
	}
});

const answer = (body: unknown) => ({ ok: true, json: async () => body });
/** The local midnight of a date, where a chart draws flai's UTC day or week of it (S-0329). */
const midnight = (date: string) => {
	const [y, m, d] = date.split('-').map(Number);
	return new Date(y, m - 1, d).getTime();
};
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const text = (sel: string) =>
	document.querySelector(sel)?.textContent?.replace(/\s+/g, ' ').trim() ?? null;
const choose = async (testid: string, value: string) => {
	const el = document.querySelector<HTMLSelectElement>(`[data-testid="${testid}"]`)!;
	el.value = value;
	el.dispatchEvent(new Event('change', { bubbles: true }));
	await settle();
};
/** The chart a page draws, 380 pixels high; the share chart of Time in State is 120 (S-0485). */
const chartEl = (height = 380) => document.querySelector(`[style*="height: ${height}px"]`);
/** The page's table view, the details element summarised "table view" (S-0485). */
const tableView = () =>
	[...document.querySelectorAll('details')].find(
		(d) => d.querySelector('summary')?.textContent?.trim() === 'table view'
	) ?? null;
/** The summary strip at the end of the controls row, whatever it holds (S-0485). */
const summaryStrip = () =>
	document
		.querySelector('[data-testid="window"]')
		?.closest('div')
		?.querySelector(':scope > span') ?? null;
/** Asserts each element is in the page and follows the one before it in the DOM (S-0485). */
const expectOrder = (...parts: (string | Element | null)[]) => {
	const els = parts.map((p) => (typeof p === 'string' ? document.querySelector(p) : p));
	const name = (i: number) => (typeof parts[i] === 'string' ? (parts[i] as string) : `part ${i}`);
	els.forEach((el, i) => expect(el, name(i)).not.toBeNull());
	for (let i = 1; i < els.length; i++)
		expect(
			els[i - 1]!.compareDocumentPosition(els[i]!) & Node.DOCUMENT_POSITION_FOLLOWING,
			`${name(i - 1)} before ${name(i)}`
		).toBeTruthy();
};
const EXPLAINER = '[data-testid="explainer"]';
const WINDOW = '[data-testid="window"]';

describe('the charts page (S-0163)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let older = false;
	beforeEach(() => {
		older = false;
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([{ id: 'E-0011', title: 'Agent status' }]);
			const bucket = new URL(url, 'http://localhost').searchParams.get('bucket') ?? 'day';
			const r = reportIn(bucket);
			if (older)
				return answer({ ...r, usage: { ...r.usage, bucket: undefined, spend: undefined } });
			return answer(r);
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	const asked = () => api.mock.calls.map(([u]) => u as string).filter((u) => u.includes('stats'));
	const controls = () =>
		[...document.querySelectorAll('label')].map((l) => l.textContent!.trim().split(/\s+/)[0]);

	it('draws tokens per day with its table', async () => {
		await open('tokens-spent');
		expect(text('h1')).toBe('Tokens / Day');
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
		// no epic narrows what flai summed
		expect(controls()).toEqual(['window', 'per', 'type']);
		expect(text('[data-testid="usage-summary"]')).toContain(
			'per story 12.0M tokens, $5.00 · 1.2M tokens per agent minute · 2.4M tokens per dollar'
		);
		// the note under the chart is in its explainer (S-0485)
		expect(document.querySelector('[data-testid="spend-note"]')).toBeNull();
		expect(text('[data-testid="explainer-how"]')).toContain(
			'Each day holds the items that entered done in it, in UTC'
		);
		expect(text('[data-testid="explainer-how"]')).toContain('The dashed line is the mean per day');
		const row = [...document.querySelectorAll('[data-testid="spend-table"] tbody tr')][0];
		expect([...row.querySelectorAll('td')].map((td) => td.textContent?.trim())).toEqual([
			'2026-09-29',
			'story',
			'3',
			'36.0M',
			'$15.00',
			'30.0',
			'12.0M',
			'$5.00',
			'10.0',
			'1.2M',
			'2.4M',
			'36.0M',
			'$15.00'
		]);
		const drawn = setOption.mock.calls.at(-1)![0] as { series: { name: string; type: string }[] };
		expect(drawn.series.map((s) => [s.name, s.type])).toEqual([
			['claude-opus-5-5', 'bar'],
			['mean per day', 'line']
		]);
	});

	it('asks flai again by the hour, names the chart for it, and goes back to days over a long window', async () => {
		await open('cost-spent');
		await choose('bucket', 'hour');
		expect(asked().at(-1)).toBe('/api/stats?since=30d&type=story&bucket=hour');
		expect(text('h1')).toBe('$ / Hour');
		// an hour in the local zone, New York's where the tests run (S-0329)
		expect(text('[data-testid="spend-table"] tbody td')).toBe('2026-09-29 15:00 EDT');
		const since = document.querySelector<HTMLSelectElement>('label select')!;
		since.value = '90d';
		since.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(asked().at(-1)).toBe('/api/stats?since=90d&type=story&bucket=day');
		const offered = [...document.querySelectorAll('[data-testid="bucket"] option')];
		expect(offered.map((o) => o.textContent)).toEqual(['day', 'week']);
	});

	it('draws one line per item type per item, with no type to choose (S-0169)', async () => {
		await open('cost-per-item');
		expect(text('h1')).toBe('$ / Work Type');
		expect(controls()).toEqual(['window', 'per']);
		const drawn = setOption.mock.calls.at(-1)![0] as { series: { name: string }[] };
		expect(drawn.series.map((s) => s.name)).toEqual(['story', 'task']);
		const rows = [...document.querySelectorAll('[data-testid="spend-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].slice(1, 3).map((td) => td.textContent?.trim())
		);
		expect(rows).toEqual([
			['story', '3'],
			['task', '8']
		]);
	});

	it('draws the mean agent time per item of the type chosen, one line per model (S-0169)', async () => {
		await open('time-per-model');
		expect(text('h1')).toBe('Avg. Time / Model');
		expect(controls()).toEqual(['window', 'per', 'type']);
		const drawn = setOption.mock.calls.at(-1)![0] as {
			series: { name: string; data: { value: [string, number] }[] }[];
			yAxis: { name: string };
		};
		// 3 stories over 30 agent minutes
		expect(drawn.series.map((s) => [s.name, s.data.map((d) => d.value[1])])).toEqual([
			['claude-opus-5-5', [10]]
		]);
		expect(drawn.yAxis.name).toBe('agent minutes per story');
	});

	it('draws the mean dollars per item of the type chosen, one line per model (S-0169)', async () => {
		await open('cost-per-model');
		expect(text('h1')).toBe('Avg. Cost / Model');
		expect(controls()).toEqual(['window', 'per', 'type']);
		const drawn = setOption.mock.calls.at(-1)![0] as { series: { name: string }[] };
		expect(drawn.series.map((s) => s.name)).toEqual(['claude-opus-5-5']);
		// the note and the star under the chart are in its explainer (S-0485)
		expect(document.querySelector('[data-testid="spend-note"]')).toBeNull();
		expect(document.body.textContent).not.toContain('* estimated in part');
		expect(text('[data-testid="explainer-how"]')).toContain('Each day holds the items');
		expect(text('[data-testid="explainer-how"]')).toContain(
			"A star marks a cost estimated in part: a task's share of its story's session, or a run that ended without its totals."
		);
		expect(
			[...document.querySelectorAll('[data-testid="spend-table"] tbody tr')].map(
				(tr) => tr.querySelectorAll('td')[1].textContent
			)
		).toEqual(['story', 'claude-opus-5-5']);
	});

	it('says the flai on the host is older when it sends no spend over time', async () => {
		older = true;
		await open('token-rate');
		expect(text('[data-testid="spend-none"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="spend-note"]')).toBeNull();
		expect(document.querySelectorAll('[data-testid="spend-table"] tbody tr').length).toBe(0);
	});

	it('lists what strategic agents spent apart in the $ / Item table (S-0225)', async () => {
		// just after midnight in UTC: the evening before in New York, where the tests run
		const done = '2026-09-29T02:00:00Z';
		const planner = { kind: 'planner', tokens: 2000, cost: 0.5, seconds: 60, estimated: true };
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			const r = reportIn('day');
			return answer({
				...r,
				items: [
					{
						id: 'S-0001',
						status: 'done',
						completed: done,
						usage: {
							tokens: 1000,
							cost: 2,
							seconds: 600,
							models: [{ model: 'claude-opus-5-5', tokens: 1000, cost: 2 }],
							strategic: [planner]
						}
					},
					{
						id: 'S-0002',
						status: 'done',
						completed: done,
						usage: { tokens: 0, cost: 0, seconds: 0, models: [], strategic: [planner] }
					}
				]
			});
		});
		await open('cost');
		const rows = [...document.querySelectorAll('[data-testid="usage-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent!.trim())
		);
		expect(rows.map((r) => [r[0], r[1], r[2]])).toEqual([
			['S-0001', '2026-09-28', 'claude-opus-5-5'],
			['S-0001', '2026-09-28', 'strategic'],
			['S-0002', '2026-09-28', 'strategic']
		]);
		expect(rows[2][5]).toMatch(/\*$/);
	});

	it('keeps the epic and the table per item on the charts of S-0143', async () => {
		await open('cost');
		expect(text('h1')).toBe('$ / Item');
		expect(controls()).toEqual(['window', 'type', 'epic']);
		expect(text('[data-testid="usage-table"] thead')).toContain('tokens/agent minute');
	});

	it('lists the longest waits under agent waiting only, and says when flai sends none (S-0215)', async () => {
		let waiting: unknown = {
			weeks: [],
			empty_wakes: { count: 0 },
			longest: [
				{
					item: 'S-0001',
					kind: 'thread',
					thread: 'TH-0003',
					started: '2026-09-29T09:00:00Z',
					seconds: 3600
				}
			]
		};
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			return answer({ ...reportIn('day'), waiting });
		});
		await open('agent-waiting');
		expect(text('h2')).toBe('Longest waits');
		expect(text('[data-testid="wait-table"] tbody')).toContain('TH-0003');
		unmount(c!);
		await open('cycle-time');
		expect(document.querySelector('[data-testid="wait-table"]')).toBeNull();
		unmount(c!);
		waiting = undefined;
		await open('agent-waiting');
		expect(text('[data-testid="waiting-older"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="wait-table"]')).toBeNull();
		expect(document.querySelector('[data-testid="wait-none"]')).toBeNull();
	});
});

describe('the window (S-0166)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const now = '2026-09-29T21:00:00Z';
	/** A report for the window asked for, with an item done a day ago and one done 20 days ago. */
	const reportFor = (since: string) => {
		const days = parseInt(since);
		const done = (id: string, completed: string) => ({
			id,
			type: 'story',
			nature: 'feature',
			title: id,
			status: 'done',
			created: '2026-09-01T00:00:00Z',
			started: '2026-09-01T00:00:00Z',
			completed,
			cycle_time_seconds: 3600,
			blocked_seconds: 0,
			time_in_state_seconds: {}
		});
		const items = [done('S-0001', '2026-09-28T21:00:00Z'), done('S-0002', '2026-09-09T21:00:00Z')];
		const start = new Date(Date.parse(now) - days * 86400e3).toISOString().replace('.000', '');
		return {
			...reportIn('day'),
			generated_at: now,
			window_days: days,
			window_start: start,
			items: items.filter((i) => i.completed >= start)
		};
	};
	const pending: { since: string; release: () => void }[] = [];
	let hold = false;
	let epicsFail = false;
	beforeEach(() => {
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items'))
				return epicsFail
					? { ok: false, statusText: 'Service Unavailable', json: async () => ({ error: 'down' }) }
					: answer([]);
			const since = new URL(url, 'http://localhost').searchParams.get('since') ?? '30d';
			if (hold) await new Promise<void>((release) => pending.push({ since, release }));
			return answer(reportFor(since));
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		hold = false;
		epicsFail = false;
		pending.length = 0;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	type Drawn = { xAxis: { min: number; max: number }; series: { data: { id: string }[] }[] };
	const drawn = () => setOption.mock.calls.at(-1)![0] as Drawn;
	/** The chart on a time axis last drawn: time in state draws its share of lead time under it. */
	const onTime = () =>
		setOption.mock.calls
			.map(([o]) => o as Drawn & { xAxis: { type: string } })
			.filter((o) => o.xAxis.type === 'time')
			.at(-1)!;
	const ids = () => drawn().series.flatMap((s) => s.data.map((d) => d.id));
	const asked = () =>
		api.mock.calls.map(([u]) => u as string).filter((u) => u.startsWith('/api/stats'));

	it('asks flai for the window chosen and redraws the axis and the points to it', async () => {
		at.params.kind = 'cycle-time';
		c = mount(ChartsPage, { target: document.body });
		await settle();
		expect(drawn().xAxis.min).toBe(Date.parse('2026-08-30T21:00:00Z'));
		expect(drawn().xAxis.max).toBe(Date.parse(now));
		expect(ids().sort()).toEqual(['S-0001', 'S-0002']);
		await choose('window', '7d');
		expect(api.mock.calls.at(-1)![0]).toBe('/api/stats?since=7d&type=story&bucket=day');
		expect(drawn().xAxis.min).toBe(Date.parse('2026-09-22T21:00:00Z'));
		expect(ids()).toEqual(['S-0001']);
	});

	it('draws the window chosen last, whatever order the answers come in', async () => {
		at.params.kind = 'cycle-time';
		c = mount(ChartsPage, { target: document.body });
		await settle();
		hold = true;
		await choose('window', '90d');
		await choose('window', '7d');
		expect(pending.map((p) => p.since)).toEqual(['90d', '7d']);
		pending[1].release();
		await settle();
		pending[0].release();
		await settle();
		expect(drawn().xAxis.min).toBe(Date.parse('2026-09-22T21:00:00Z'));
		expect(ids()).toEqual(['S-0001']);
	});

	it('opens every chart at the window chosen before, after leaving the page (S-0168)', async () => {
		at.params.kind = 'burn-up';
		c = mount(ChartsPage, { target: document.body });
		await settle();
		await choose('window', '7d');
		unmount(c);
		document.body.innerHTML = '';
		for (const kind of ['cycle-time', 'time-in-state']) {
			api.mockClear();
			at.params.kind = kind;
			c = mount(ChartsPage, { target: document.body });
			await settle();
			expect(asked()).toEqual(['/api/stats?since=7d&type=story&bucket=day']);
			expect(document.querySelector<HTMLSelectElement>('[data-testid="window"]')!.value).toBe('7d');
			// cycle time from the window's start; time in state from half a day before the day that holds
			// it, drawn on its date's tick
			expect(onTime().xAxis.min).toBe(
				kind === 'cycle-time'
					? Date.parse('2026-09-22T21:00:00Z')
					: midnight('2026-09-22') - 12 * 3600e3
			);
			unmount(c);
			document.body.innerHTML = '';
		}
		c = undefined;
	});

	it('draws time in state by the day on an axis that follows the window (S-0168)', async () => {
		at.params.kind = 'time-in-state';
		c = mount(ChartsPage, { target: document.body });
		await settle();
		const half = 12 * 3600e3;
		expect(onTime().xAxis.min).toBe(midnight('2026-08-30') - half);
		expect(onTime().xAxis.max).toBe(midnight('2026-09-29') + half);
		const days = () =>
			(onTime().series[0].data as unknown as { ids: string[] }[]).map((d) => d.ids);
		expect(days()).toEqual([['S-0002'], ['S-0001']]);
		await choose('window', '7d');
		expect(onTime().xAxis.min).toBe(midnight('2026-09-22') - half);
		expect(days()).toEqual([['S-0001']]);
	});

	it('draws the charts without waiting for the epics, and when they cannot be read (S-0168)', async () => {
		epicsFail = true;
		at.params.kind = 'cycle-time';
		c = mount(ChartsPage, { target: document.body });
		await settle();
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
		expect(ids().sort()).toEqual(['S-0001', 'S-0002']);
		expect(document.querySelectorAll('label select')[2].querySelectorAll('option').length).toBe(1);
	});
});

describe('the planning charts (S-0212)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const story = (id: string, completed: string, rest: Record<string, unknown>) => ({
		id,
		type: 'story',
		nature: 'feature',
		title: `Story ${id}`,
		status: 'done',
		created: '2026-09-01T00:00:00Z',
		started: '2026-09-01T00:00:00Z',
		completed,
		blocked_seconds: 0,
		time_in_state_seconds: {},
		...rest
	});
	const spread = (count: number, seconds: number) => ({
		count,
		p50_seconds: seconds,
		p85_seconds: seconds
	});
	// flai's own figures: the p50 and p85 of the absolute errors, in all, per nature, and per model
	const forecasts = {
		forecast: {
			count: 2,
			p50_seconds: 7200,
			p85_seconds: 10800,
			by_nature: { feature: spread(1, 7200), remediation: spread(1, 10800) },
			by_model: { 'claude-opus-5-5': spread(1, 7200), '(none)': spread(1, 10800) }
		},
		delivery: {
			...spread(1, 86400),
			by_nature: { feature: spread(1, 86400) },
			by_model: { 'claude-opus-5-5': spread(1, 86400) }
		},
		estimate: {
			count: 2,
			p50_seconds: 1800,
			p85_seconds: 3600,
			by_nature: { feature: { count: 2, p50_seconds: 1800, p85_seconds: 3600 } },
			by_model: { 'claude-opus-5-5': spread(1, 1800), '(none)': spread(1, 3600) }
		}
	};
	const items = [
		story('S-0001', '2026-09-10T12:00:00Z', { estimate_error_seconds: -3600 }),
		story('S-0002', '2026-09-20T12:00:00Z', {
			nature: 'remediation',
			forecast_seconds: 14400,
			cycle_time_seconds: 3600,
			forecast_error_seconds: -10800
		}),
		story('S-0003', '2026-09-28T12:00:00Z', {
			model: 'claude-opus-5-5',
			forecast_seconds: 21600,
			cycle_time_seconds: 28800,
			forecast_error_seconds: 7200,
			delivery_error_seconds: -86400,
			estimate_error_seconds: 1800
		})
	];
	let answered: Record<string, unknown> = {};
	beforeEach(() => {
		answered = { items, forecasts };
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			const q = new URL(url, 'http://localhost').searchParams;
			return answer({
				...reportIn(q.get('bucket') ?? 'day'),
				type: q.get('type'),
				...answered
			});
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	/** Goes to another chart on the same page, as a link in the header does. */
	const go = async (kind: string) => {
		at.params.kind = kind;
		at.moved();
		await settle();
	};
	const asked = () =>
		api.mock.calls.map(([u]) => u as string).filter((u) => u.startsWith('/api/stats'));
	const controls = () =>
		[...document.querySelectorAll('label')].map((l) => l.textContent!.trim().split(/\s+/)[0]);
	const options = (testid: string) =>
		[...document.querySelectorAll(`[data-testid="${testid}"] option`)].map((o) => o.textContent);
	type Drawn = { series: { name: string; data: { id: string }[] }[] };
	const drawn = () => setOption.mock.calls.at(-1)![0] as Drawn;
	const points = () => drawn().series.map((s) => [s.name, s.data.map((d) => d.id)]);
	const rows = () =>
		[...document.querySelectorAll('[data-testid="forecast-table"] tbody tr')].map((tr) =>
			tr.querySelector('a')!.textContent!.trim()
		);

	it('draws each planning chart it goes to from the one report', async () => {
		await open('forecast-accuracy');
		expect(text('h1')).toBe('Forecast Accuracy');
		expect(points()).toEqual([
			['forecast error', ['S-0002', 'S-0003']],
			['estimate error', ['S-0001', 'S-0003']]
		]);
		await go('delivery-accuracy');
		expect(text('h1')).toBe('Delivery Accuracy');
		expect(drawn().series[0]).toMatchObject({ name: 'delivery error', data: [{ id: 'S-0003' }] });
		await go('forecast-by-model');
		expect(text('h1')).toBe('Forecast Error / Model');
		expect(drawn().series.map((s) => s.name)).toEqual(['(none)', 'claude-opus-5-5']);
		// the report was read once: every planning chart reads the same stories
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
	});

	it('narrows a chart by nature and by model, with flai figures for one and the rows for both', async () => {
		await open('forecast-accuracy');
		expect(controls()).toEqual(['window', 'nature', 'model']);
		expect(options('nature')).toEqual(['all', 'feature', 'remediation']);
		expect(options('model')).toEqual(['all', '(none)', 'claude-opus-5-5']);
		expect(text('[data-testid="planning-summary"]')).toBe(
			'2 stories with a forecast in the window · forecast error p50 2h p85 3h · delivery error p50 1d p85 1d'
		);
		expect(rows()).toEqual(['S-0003', 'S-0002', 'S-0001']);
		await choose('nature', 'remediation');
		expect(points()).toEqual([['forecast error', ['S-0002']]]);
		expect(rows()).toEqual(['S-0002']);
		expect(text('[data-testid="planning-summary"]')).toBe(
			'1 story with a forecast (remediation) · forecast error p50 3h p85 3h · delivery error p50 - p85 -'
		);
		await choose('nature', 'feature');
		await choose('model', '(none)');
		expect(rows()).toEqual(['S-0001']);
		expect(text('[data-testid="planning-summary"]')).toBe(
			'0 stories with a forecast (feature, (none)) · forecast error p50 - p85 - · delivery error p50 - p85 -'
		);
		await choose('model', 'claude-opus-5-5');
		expect(rows()).toEqual(['S-0003']);
		expect(text('[data-testid="planning-summary"]')).toBe(
			'1 story with a forecast (feature, claude-opus-5-5) · forecast error p50 2h p85 2h · delivery error p50 1d p85 1d'
		);
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
	});

	it('shows the error per model by the bucket, narrowed by nature only', async () => {
		await open('forecast-by-model');
		expect(controls()).toEqual(['window', 'per', 'nature']);
		expect(document.querySelector('[data-testid="model"]')).toBeNull();
		expect(options('nature')).toEqual(['all', 'feature', 'remediation']);
		await choose('bucket', 'week');
		expect(asked().at(-1)).toBe('/api/stats?since=30d&type=story&bucket=week');
		await choose('nature', 'feature');
		expect(drawn().series.map((s) => s.name)).toEqual(['claude-opus-5-5']);
	});

	it('asks flai for stories, whichever type was chosen on another chart', async () => {
		await open('cycle-time');
		const type = [...document.querySelectorAll('label')]
			.find((l) => l.textContent!.trim().startsWith('type'))!
			.querySelector('select')!;
		type.value = 'task';
		type.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(asked()).toEqual([
			'/api/stats?since=30d&type=story&bucket=day',
			'/api/stats?since=30d&type=task&bucket=day'
		]);
		await go('forecast-accuracy');
		expect(asked().at(-1)).toBe('/api/stats?since=30d&type=story&bucket=day');
		expect(text('[data-testid="planning-summary"]')).toContain('2 stories with a forecast');
		await go('delivery-accuracy');
		await go('cycle-time');
		expect(asked().slice(2)).toEqual([
			'/api/stats?since=30d&type=story&bucket=day',
			'/api/stats?since=30d&type=task&bucket=day'
		]);
	});

	it('says how a forecast is set when no story in the window has one', async () => {
		answered = {
			items: [items[0]],
			forecasts: {
				forecast: { count: 0, by_nature: {}, by_model: {} },
				delivery: { count: 0, by_nature: {}, by_model: {} },
				estimate: { ...spread(1, 3600), by_nature: {}, by_model: {} }
			}
		};
		await open('forecast-accuracy');
		expect(text('[data-testid="forecast-none"]')).toContain(
			'flai edit <story> --forecast-duration 6h --forecast-delivery <UTC time>'
		);
		expect(text('[data-testid="forecast-none"]')).toContain('The planner sets one');
		expect(text('[data-testid="planning-summary"]')).toBe(
			'0 stories with a forecast in the window · forecast error p50 - p85 - · delivery error p50 - p85 -'
		);
		expect(rows()).toEqual(['S-0001']);
		// the notice above the chart, the filters and the summary after its explainer (S-0485)
		expectOrder(
			'h1',
			'[data-testid="forecast-none"]',
			chartEl(),
			EXPLAINER,
			WINDOW,
			'[data-testid="nature"]',
			'[data-testid="model"]',
			'[data-testid="planning-summary"]',
			tableView()
		);
	});

	it('says the flai on the host is older when it sends no forecast errors', async () => {
		answered = { items: [], forecasts: undefined };
		await open('delivery-accuracy');
		expect(text('[data-testid="forecasts-older"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="forecast-none"]')).toBeNull();
		expect(document.querySelector('[data-testid="planning-summary"]')).toBeNull();
	});
});

describe('the claims charts (S-0214)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const story = (id: string, status: string, completed: string) => ({
		id,
		type: 'story',
		nature: 'feature',
		title: `Story ${id}`,
		status,
		created: '2026-08-01T00:00:00Z',
		started: '2026-08-01T00:00:00Z',
		completed,
		blocked_seconds: 0,
		time_in_state_seconds: {}
	});
	// in the window of reportIn, 30 August 21:00 to 29 September 21:00: S-0001 and S-0002 done in
	// it, S-0003 cancelled in it, S-0004 done before it
	const items = [
		story('S-0001', 'done', '2026-09-22T12:00:00Z'),
		// done just after midnight in UTC: the evening of the 23rd in New York, where the tests run
		story('S-0002', 'done', '2026-09-24T02:00:00Z'),
		story('S-0003', 'cancelled', '2026-09-25T12:00:00Z'),
		story('S-0004', 'done', '2026-08-01T12:00:00Z')
	];
	const drift = (id: string, outside: string[], unchanged: string[]) => ({
		id,
		committed: [...outside, 'flai/cmd/stats.go'].sort(),
		outside,
		unchanged,
		outside_count: outside.length,
		unchanged_count: unchanged.length
	});
	const claims = {
		limit: 2,
		days: [
			{ date: '2026-09-27', in_progress: 1, held: 0 },
			{ date: '2026-09-28', in_progress: 2, held: 1 },
			{ date: '2026-09-29', in_progress: 3, held: 2 }
		],
		weeks: [
			{
				week: '2026-W39',
				start: '2026-09-21',
				held_seconds: { overlap: 7200, after: 3600, 'no-touches': 0 },
				stories: 2,
				exact: 1,
				exact_share: 0.5
			},
			{
				week: '2026-W40',
				start: '2026-09-28',
				held_seconds: { overlap: 0, after: 0, 'no-touches': 5400 },
				stories: 0,
				exact: 0
			}
		],
		drift: [
			drift('S-0001', [], []),
			drift('S-0002', ['docs/a.md', 'docs/b.md'], ['flaiover']),
			drift('S-0003', ['docs/c.md'], []),
			drift('S-0004', ['docs/d.md'], [])
		]
	};
	let answered: Record<string, unknown> = {};
	beforeEach(() => {
		answered = { items, claims };
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			const q = new URL(url, 'http://localhost').searchParams;
			return answer({
				...reportIn(q.get('bucket') ?? 'day'),
				type: q.get('type'),
				...answered
			});
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	const go = async (kind: string) => {
		at.params.kind = kind;
		at.moved();
		await settle();
	};
	const asked = () =>
		api.mock.calls.map(([u]) => u as string).filter((u) => u.startsWith('/api/stats'));
	const controls = () =>
		[...document.querySelectorAll('label')].map((l) => l.textContent!.trim().split(/\s+/)[0]);
	type Drawn = { series: { name: string; data: unknown[] }[] };
	const drawn = () => setOption.mock.calls.at(-1)![0] as Drawn;
	const cells = (testid: string) =>
		[...document.querySelectorAll(`[data-testid="${testid}"] tbody tr`)].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent!.replace(/\s+/g, ' ').trim())
		);
	const notes = () =>
		['forecasts-older', 'forecast-none', 'planning-summary', 'claims-older', 'drift-none'].filter(
			(id) => document.querySelector(`[data-testid="${id}"]`) !== null
		);

	it('draws parallelism per day with held and the limit, from stories over the window chosen', async () => {
		await open('cycle-time');
		const type = [...document.querySelectorAll('label')]
			.find((l) => l.textContent!.trim().startsWith('type'))!
			.querySelector('select')!;
		type.value = 'task';
		type.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		await go('parallelism');
		expect(text('h1')).toBe('Parallelism');
		// stories, whichever type was chosen on another chart
		expect(asked().at(-1)).toBe('/api/stats?since=30d&type=story&bucket=day');
		expect(controls()).toEqual(['window']);
		expect(drawn().series.map((s) => [s.name, s.data])).toEqual([
			[
				'in progress',
				[
					['2026-09-27', 1],
					['2026-09-28', 2],
					['2026-09-29', 3]
				]
			],
			[
				'held',
				[
					['2026-09-27', 0],
					['2026-09-28', 1],
					['2026-09-29', 2]
				]
			],
			[
				'limit',
				[
					['2026-09-27', 2],
					['2026-09-28', 2],
					['2026-09-29', 2]
				]
			]
		]);
		expect(cells('parallelism-table')).toEqual([
			['2026-09-27', '1', '0', '2'],
			['2026-09-28', '2', '1', '2'],
			['2026-09-29', '3', '2', '2']
		]);
		// no forecast note on a chart that is not drawn from forecasts
		expect(notes()).toEqual([]);
		await choose('window', '7d');
		expect(asked().at(-1)).toBe('/api/stats?since=7d&type=story&bucket=day');
	});

	it('draws hold time per week in hours, stacked by reason', async () => {
		await open('hold-time');
		expect(text('h1')).toBe('Hold Time');
		expect(controls()).toEqual(['window']);
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
		const bars = drawn().series.map((s) => [
			s.name,
			(s.data as { value: [number, number] }[]).map((d) => d.value)
		]);
		// each week on its Monday's tick
		const w39 = midnight('2026-09-21');
		const w40 = midnight('2026-09-28');
		expect(bars).toEqual([
			[
				'overlap',
				[
					[w39, 2],
					[w40, 0]
				]
			],
			[
				'after',
				[
					[w39, 1],
					[w40, 0]
				]
			],
			[
				'empty claim',
				[
					[w39, 0],
					[w40, 1.5]
				]
			]
		]);
		expect(cells('hold-time-table')).toEqual([
			['2026-W39', '2026-09-21', '2h', '1h', '0m', '3h'],
			['2026-W40', '2026-09-28', '0m', '0m', '1.5h', '1.5h']
		]);
		expect(notes()).toEqual([]);
	});

	it('draws touches drift per story done in the window, with the weekly share of exact touches', async () => {
		await open('touches-drift');
		expect(text('h1')).toBe('Touches Drift');
		expect(controls()).toEqual(['window']);
		const series = drawn().series;
		const bars = series
			.filter((s) => s.name !== 'exact touches per week')
			.map((s) => [
				s.name,
				(s.data as { id: string; value: [string, number] }[]).map((d) => [d.id, d.value[1]])
			]);
		// S-0003 cancelled and S-0004 done before the window are left out
		expect(bars).toEqual([
			[
				'outside its touches',
				[
					['S-0001', 0],
					['S-0002', 2]
				]
			],
			[
				'touches unchanged',
				[
					['S-0001', 0],
					['S-0002', 1]
				]
			]
		]);
		const line = series.find((s) => s.name === 'exact touches per week')!;
		// a week without a story done has no share, a gap
		expect((line.data as { value: [number, number | null] }[]).map((d) => d.value[1])).toEqual([
			0.5,
			null
		]);
		expect(cells('drift-table')).toEqual([
			['S-0001 Story S-0001', '2026-09-22', '0', '0'],
			['S-0002 Story S-0002', '2026-09-23', '2', '1']
		]);
		expect(cells('exact-table')).toEqual([
			['2026-W39', '2026-09-21', '2', '1', '50%'],
			['2026-W40', '2026-09-28', '0', '0', '-']
		]);
		expect(notes()).toEqual([]);
	});

	it('says git could not be read when flai sends no drift', async () => {
		answered = {
			items,
			claims: {
				...claims,
				weeks: claims.weeks.map(({ week, start, held_seconds }) => ({ week, start, held_seconds })),
				drift: undefined
			}
		};
		await open('touches-drift');
		expect(notes()).toEqual(['drift-none']);
		expect(text('[data-testid="drift-none"]')).toContain('could not read git');
		expect(cells('drift-table')).toEqual([]);
		expect(cells('exact-table')).toEqual([
			['2026-W39', '2026-09-21', '-', '-', '-'],
			['2026-W40', '2026-09-28', '-', '-', '-']
		]);
		// hold time is still drawn: the reasons do not need git
		await go('hold-time');
		expect(notes()).toEqual([]);
		expect(cells('hold-time-table').length).toBe(2);
	});

	// coordination as flai sends it (S-0337) over the weeks of the window: in W39 a conversation
	// opened and a conflict found early over two paths, and a story started on a share that would
	// have waited 20 hours; in W40 the conversation agreed and closed with its story, a conflict met
	// at sync and one at acceptance, and two stories on a share that would have waited 2 and 4 hours
	const coordination = {
		logged: true,
		weeks: [
			{
				week: '2026-W39',
				start: '2026-09-21',
				conversations: { opened: 1, agreed: 0, escalated: 0, closed_with_story: 0 },
				conflicts: {
					found: 1,
					found_paths: 2,
					met_sync: 0,
					met_sync_paths: 0,
					met_accept: 0,
					met_accept_paths: 0
				},
				shares: { started: 1, held_seconds: 72000 }
			},
			{
				week: '2026-W40',
				start: '2026-09-28',
				conversations: { opened: 0, agreed: 1, escalated: 0, closed_with_story: 1 },
				conflicts: {
					found: 0,
					found_paths: 0,
					met_sync: 1,
					met_sync_paths: 1,
					met_accept: 1,
					met_accept_paths: 1
				},
				shares: { started: 2, held_seconds: 21600 }
			}
		]
	};
	const coordinationNotes = () =>
		['coordination-older', 'conflicts-unread', 'claims-older'].filter(
			(id) => document.querySelector(`[data-testid="${id}"]`) !== null
		);

	it('draws coordination per week with the stories started on a share and the hold time saved (S-0337)', async () => {
		answered = { items, claims, coordination };
		await open('coordination');
		expect(text('h1')).toBe('Coordination');
		expect(controls()).toEqual(['window']);
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
		const w39 = midnight('2026-09-21');
		const w40 = midnight('2026-09-28');
		expect(
			drawn().series.map((s) => [
				s.name,
				(s.data as { value: [number, number] }[]).map((d) => d.value)
			])
		).toEqual([
			[
				'found early',
				[
					[w39, 1],
					[w40, 0]
				]
			],
			[
				'met late',
				[
					[w39, 0],
					[w40, 2]
				]
			],
			[
				'opened',
				[
					[w39, 1],
					[w40, 0]
				]
			],
			[
				'agreed',
				[
					[w39, 0],
					[w40, 1]
				]
			],
			[
				'escalated',
				[
					[w39, 0],
					[w40, 0]
				]
			],
			[
				'closed with story',
				[
					[w39, 0],
					[w40, 1]
				]
			]
		]);
		expect(text('[data-testid="coordination-shares"]')).toBe(
			'Stories started on a share in the window: 3 · hold time saved: 26h.'
		);
		expect(cells('coordination-table')).toEqual([
			['2026-W39', '2026-09-21', '1', '0', '0', '0', '1', '0', '0', '1', '20h'],
			['2026-W40', '2026-09-28', '0', '1', '0', '1', '0', '1', '1', '2', '6h']
		]);
		expect(coordinationNotes()).toEqual([]);
		// the note under the chart is in its explainer (S-0485)
		expect(document.querySelector('[data-testid="coordination-note"]')).toBeNull();
		expect(text('[data-testid="explainer-how"]')).toContain(
			'The bars are the conflicts a trial merge found while both stories were open, early'
		);
		expect(text('[data-testid="explainer-how"]')).toContain(
			'is what the stories started on a share would have waited for their holder'
		);
		await choose('window', '7d');
		expect(asked().at(-1)).toBe('/api/stats?since=7d&type=story&bucket=day');
	});

	it('says the conflict log could not be read, and draws the conversations alone (S-0337)', async () => {
		const weeks = coordination.weeks.map(({ week, start, conversations, shares }) => ({
			week,
			start,
			conversations,
			shares
		}));
		answered = { items, claims, coordination: { logged: false, weeks } };
		await open('coordination');
		expect(coordinationNotes()).toEqual(['conflicts-unread']);
		expect(text('[data-testid="conflicts-unread"]')).toContain('.flai-cache/conflicts.jsonl');
		expect(drawn().series.map((s) => s.name)).toEqual([
			'opened',
			'agreed',
			'escalated',
			'closed with story'
		]);
		expect(cells('coordination-table').map((r) => r.slice(6, 9))).toEqual([
			['-', '-', '-'],
			['-', '-', '-']
		]);
		expect(text('[data-testid="coordination-shares"]')).toContain('3');
	});

	it('puts the notice, the chart, the shares, the explainer, the window, and the table in order (S-0485)', async () => {
		answered = { items, claims, coordination };
		await open('coordination');
		expectOrder(
			'h1',
			chartEl(),
			'[data-testid="coordination-shares"]',
			EXPLAINER,
			WINDOW,
			summaryStrip(),
			tableView()
		);
		unmount(c!);
		document.body.innerHTML = '';
		answered = { items, claims, coordination: { ...coordination, logged: false } };
		await open('coordination');
		expectOrder(
			'h1',
			'[data-testid="conflicts-unread"]',
			chartEl(),
			'[data-testid="coordination-shares"]',
			EXPLAINER,
			WINDOW,
			tableView()
		);
	});

	it('says the flai on the host is older when it sends no coordination (S-0337)', async () => {
		answered = { items, claims };
		await open('coordination');
		expect(coordinationNotes()).toEqual(['coordination-older']);
		expect(text('[data-testid="coordination-older"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="coordination-shares"]')).toBeNull();
		expect(cells('coordination-table')).toEqual([]);
		// the other claims charts are drawn as before
		await go('hold-time');
		expect(coordinationNotes()).toEqual([]);
	});

	it('says the flai on the host is older when it sends no held stories or time held', async () => {
		const days = claims.days.map(({ date, in_progress }) => ({ date, in_progress }));
		answered = { items, claims: { limit: 2, days, drift: claims.drift } };
		await open('parallelism');
		expect(notes()).toEqual(['claims-older']);
		expect(text('[data-testid="claims-older"]')).toContain('flai self-upgrade');
		expect(cells('parallelism-table').map((r) => r[2])).toEqual(['-', '-', '-']);
		for (const kind of ['hold-time', 'touches-drift']) {
			await go(kind);
			expect(notes()).toEqual(['claims-older']);
		}
		answered = { items, claims: undefined };
		await choose('window', '7d');
		expect(notes()).toEqual(['claims-older']);
	});

	// the merge queue as flai sends it (S-0404) over the weeks of the window: in W39 a batch of two
	// that landed after 1 and 3 hours; in W40 a red batch of one bisected and sent back after 2 hours
	const queue = {
		weeks: [
			{
				week: '2026-W39',
				start: '2026-09-21',
				landed: 2,
				sent_back: 0,
				wait_mean_seconds: 7200,
				wait_max_seconds: 10800,
				batches: 1,
				batch_size_mean: 2,
				batch_size_max: 2,
				bisections: 0
			},
			{
				week: '2026-W40',
				start: '2026-09-28',
				landed: 0,
				sent_back: 1,
				wait_mean_seconds: 7200,
				wait_max_seconds: 7200,
				batches: 1,
				batch_size_mean: 1,
				batch_size_max: 1,
				bisections: 1
			}
		]
	};
	const queueNotes = () =>
		['queue-none', 'queue-note', 'waiting-older'].filter(
			(id) => document.querySelector(`[data-testid="${id}"]`) !== null
		);

	it('draws the merge queue per week beside agent waiting, with the window alone (S-0404)', async () => {
		answered = { items, queue };
		await open('queue');
		expect(text('h1')).toBe('Merge Queue');
		expect(controls()).toEqual(['window']);
		const w39 = midnight('2026-09-21');
		const w40 = midnight('2026-09-28');
		expect(
			drawn().series.map((s) => [
				s.name,
				(s.data as { value: [number, number | null] }[]).map((d) => d.value)
			])
		).toEqual([
			[
				'mean wait',
				[
					[w39, 2],
					[w40, 2]
				]
			],
			[
				'longest wait',
				[
					[w39, 3],
					[w40, 2]
				]
			],
			[
				'mean batch size',
				[
					[w39, 2],
					[w40, 1]
				]
			],
			[
				'bisections',
				[
					[w39, 0],
					[w40, 1]
				]
			]
		]);
		// the note under the chart is in its explainer (S-0485)
		expect(queueNotes()).toEqual([]);
		expect(text('[data-testid="explainer-what"]')).toContain(
			'The lines are the mean and the longest time the stories the merge queue landed or sent back each week waited in it'
		);
		expect(text('[data-testid="explainer-what"]')).toContain(
			'the mean number of stories in the batches it built, and as bars the red batches it bisected'
		);
		expect(text('[data-testid="explainer-how"]')).toContain('from its enqueueing');
		expect(cells('queue-table')).toEqual([
			['2026-W39', '2026-09-21', '2', '0', '2h', '3h', '1', '2', '2', '0'],
			['2026-W40', '2026-09-28', '0', '1', '2h', '2h', '1', '1', '1', '1']
		]);
	});

	it('says the flai on the host sends no merge queue when the report has none (S-0404)', async () => {
		answered = { items };
		await open('queue');
		expect(queueNotes()).toEqual(['queue-none']);
		expect(text('[data-testid="queue-none"]')).toContain('sends no merge queue');
		expect(text('[data-testid="queue-none"]')).toContain('logged no event yet');
		expect(text('[data-testid="queue-none"]')).toContain('flai self-upgrade');
		expect(drawn().series).toEqual([]);
		expect(cells('queue-table')).toEqual([]);
		// the other charts say nothing of it
		await go('cycle-time');
		expect(queueNotes()).toEqual([]);
	});

	// the digest as flai sends it (S-0434) over two weeks of the window: in W39 five threads to the
	// operator for two stories completed, one a repeat; in W40 two threads and no story completed
	const digest = {
		weeks: [
			{
				week: '2026-W39',
				start: '2026-09-21',
				threads_to_operator: 5,
				stories_completed: 2,
				per_story: 2.5,
				asked_again: 1
			},
			{
				week: '2026-W40',
				start: '2026-09-28',
				threads_to_operator: 2,
				stories_completed: 0,
				asked_again: 0
			}
		],
		entries: { active: 120, compacted: 3, absorbed: 5 }
	};
	const digestNotes = () =>
		['digest-none', 'digest-entries', 'digest-never-curated'].filter(
			(id) => document.querySelector(`[data-testid="${id}"]`) !== null
		);

	it('draws the threads per completed story with the repeats stacked, and the entries beside it (S-0434)', async () => {
		answered = { items, digest };
		await open('digest');
		expect(text('h1')).toBe('Digest');
		expect(controls()).toEqual(['window']);
		const w39 = midnight('2026-09-21');
		const w40 = midnight('2026-09-28');
		expect(
			drawn().series.map((s) => [
				s.name,
				(s as { stack?: string }).stack,
				(s.data as { value: [number, number | null] }[]).map((d) => d.value)
			])
		).toEqual([
			[
				'other threads',
				'digest',
				[
					[w39, 2],
					[w40, null]
				]
			],
			[
				'asked again',
				'digest',
				[
					[w39, 0.5],
					[w40, null]
				]
			]
		]);
		expect(digestNotes()).toEqual(['digest-entries']);
		expect(text('[data-testid="digest-entries"]')).toBe(
			'Digest entries now: active 120 · compacted 3 · absorbed 5.'
		);
		expect(text('[data-testid="explainer-what"]')).toContain(
			'the threads opened to the operator per story completed that week'
		);
		// a flai older than S-0504 sends no curator runs, and asked_again 0 where it counted none
		expect(cells('digest-table')).toEqual([
			['2026-W39', '2026-09-21', '5', '2', '2.50', '1', '-'],
			['2026-W40', '2026-09-28', '2', '0', '-', '0', '-']
		]);
		expectOrder('h1', chartEl(), '[data-testid="digest-entries"]', EXPLAINER, WINDOW, tableView());
	});

	it('says the flai on the host sends no digest when the report has none (S-0434)', async () => {
		answered = { items };
		await open('digest');
		expect(digestNotes()).toEqual(['digest-none']);
		expect(text('[data-testid="digest-none"]')).toContain('sends no digest');
		expect(text('[data-testid="digest-none"]')).toContain('flai self-upgrade');
		expect(drawn().series).toEqual([]);
		expect(cells('digest-table')).toEqual([]);
		// the other charts say nothing of it
		await go('queue');
		expect(digestNotes()).toEqual([]);
	});

	it('draws the curator runs, and marks a week no run has reached as unmeasured (S-0504)', async () => {
		// one curator run in W39, the newest, so W40 is a week no run has reached
		const curated = {
			...digest,
			weeks: [
				{ ...digest.weeks[0], curator_runs: 1 },
				{ ...digest.weeks[1], asked_again: undefined, curator_runs: 0 }
			],
			curator: { runs: 1, last_run: '2026-09-23T05:00:00Z' }
		};
		answered = { items, digest: curated };
		await open('digest');
		expect(drawn().series.map((s) => s.name)).toEqual([
			'other threads',
			'asked again',
			'curator runs'
		]);
		expect(digestNotes()).toEqual(['digest-entries']);
		expect(text('[data-testid="explainer-what"]')).toContain("the curator's runs each week");
		expect(cells('digest-table')).toEqual([
			['2026-W39', '2026-09-21', '5', '2', '2.50', '1', '1'],
			['2026-W40', '2026-09-28', '2', '0', '-', '-', '0']
		]);
	});

	it('says the curator has never run, so no repeat was measured (S-0504, I-0207)', async () => {
		const never = {
			...digest,
			weeks: digest.weeks.map((wk) => ({ ...wk, asked_again: undefined, curator_runs: 0 })),
			curator: { runs: 0 }
		};
		answered = { items, digest: never };
		await open('digest');
		expect(digestNotes()).toEqual(['digest-entries', 'digest-never-curated']);
		expect(text('[data-testid="digest-never-curated"]')).toBe(
			'The curator has never run, so no repeat was measured: curate is off on the host, or its ' +
				'first run is still to come.'
		);
		expect(cells('digest-table')).toEqual([
			['2026-W39', '2026-09-21', '5', '2', '2.50', '-', '0'],
			['2026-W40', '2026-09-28', '2', '0', '-', '-', '0']
		]);
		expectOrder(
			'h1',
			chartEl(),
			'[data-testid="digest-entries"]',
			'[data-testid="digest-never-curated"]',
			EXPLAINER,
			WINDOW,
			tableView()
		);
	});
});

describe('the strategic charts (S-0216)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const story = (id: string, completed: string, rest: Record<string, unknown>) => ({
		id,
		type: 'story',
		nature: 'feature',
		title: `Story ${id}`,
		status: 'done',
		created: '2026-09-01T00:00:00Z',
		started: '2026-09-01T00:00:00Z',
		completed,
		blocked_seconds: 0,
		time_in_state_seconds: {},
		...rest
	});
	// one story done on the 28th that waited an hour, two on the 29th that waited an hour and two
	const items = [
		story('S-0001', '2026-09-28T12:00:00Z', {
			wait_threads_seconds: 1800,
			wait_review_seconds: 1800
		}),
		story('S-0002', '2026-09-29T09:00:00Z', { wait_threads_seconds: 3600 }),
		story('S-0003', '2026-09-29T15:00:00Z', { wait_review_seconds: 7200 })
	];
	// $1.50 over 3 stories: $0.50 a story, a tenth of the agents' $15 over 3
	const days = [
		{
			date: '2026-09-27',
			agents: { planner: { cost: 0.5, seconds: 1800 } },
			cost: 0.5,
			seconds: 1800,
			completed: 0
		},
		{
			date: '2026-09-28',
			agents: {
				planner: { cost: 0.25, seconds: 900, estimated: true },
				orchestrator: { cost: 0.25, seconds: 1800 }
			},
			cost: 0.5,
			seconds: 2700,
			completed: 1,
			cost_per_item: 4,
			cycle_time_seconds: 7200
		},
		{
			date: '2026-09-29',
			agents: { analyzer: { cost: 0.5, seconds: 3600 } },
			cost: 0.5,
			seconds: 3600,
			completed: 2,
			cost_per_item: 5.5,
			cycle_time_seconds: 10800
		}
	];
	let answered: Record<string, unknown> = {};
	beforeEach(() => {
		answered = { items, strategic_days: days };
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			const q = new URL(url, 'http://localhost').searchParams;
			return answer({ ...reportIn(q.get('bucket') ?? 'day'), type: q.get('type'), ...answered });
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	const controls = () =>
		[...document.querySelectorAll('label')].map((l) => l.textContent!.trim().split(/\s+/)[0]);
	const rows = () =>
		[...document.querySelectorAll('[data-testid="strategic-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent!.trim())
		);
	const heads = () =>
		[...document.querySelectorAll('[data-testid="strategic-table"] th')].map((th) =>
			th.textContent!.trim()
		);
	const series = () =>
		(setOption.mock.calls.at(-1)![0] as { series: { name: string }[] }).series.map((s) => s.name);

	it('states the ratio on Strategic Cost', async () => {
		await open('strategic-cost');
		expect(text('h1')).toBe('Strategic Cost');
		expect(controls()).toEqual(['window', 'type']);
		expect(text('[data-testid="strategic-summary"]')).toBe(
			'strategic agents $1.50 · 2.3h · 3 stories completed'
		);
		// the explanation under the chart is in its explainer, the ratio stays under it (S-0485)
		expect(document.querySelector('[data-testid="strategic-note"]')).toBeNull();
		expect(text('[data-testid="explainer-what"]')).toContain(
			'The bars are what planning, orchestration, and analysis cost each day'
		);
		expect(text('[data-testid="explainer-what"]')).toContain(
			'the dashed line is what the strategic agents add to each story over the window'
		);
		expect(text('[data-testid="explainer-how"]')).toContain(
			'A star marks a cost estimated in part: an activity of that agent that day carries an estimated cost.'
		);
		expect(document.body.textContent).not.toContain('* estimated in part');
		expect(text('[data-testid="strategic-figures"]')).not.toContain('The bars');
		expect(text('[data-testid="strategic-figures"] [data-testid="strategic-ratio"]')).toBe(
			'Over the window they spent $1.50 for 3 stories completed: $0.500 per story, 10% of the $5.00 the agents spent per story.'
		);
		expect(series()).toEqual(['planner', 'orchestrator', 'analyzer', 'mean per story']);
		expect(heads()).toEqual([
			'day',
			'planner',
			'orchestrator',
			'analyzer',
			'total',
			'completed',
			'mean per story'
		]);
		expect(rows()).toEqual([
			['2026-09-27', '$0.500', '-', '-', '$0.500', '0', '-'],
			['2026-09-28', '$0.250*', '$0.250', '-', '$0.500', '1', '$4.00'],
			['2026-09-29', '-', '-', '$0.500', '$0.500', '2', '$5.50']
		]);
	});

	it('draws Strategic Use in hours against the cycle time and waiting of the type chosen', async () => {
		await open('strategic-use');
		expect(text('h1')).toBe('Strategic Use');
		// the note under the chart is in its explainer (S-0485)
		expect(document.querySelector('[data-testid="strategic-note"]')).toBeNull();
		expect(document.querySelector('[data-testid="strategic-figures"]')).toBeNull();
		expect(text('[data-testid="explainer-what"]')).toContain(
			'The bars are the hours the planner, the orchestrator, and the analyzer worked each day'
		);
		expect(text('[data-testid="explainer-what"]')).toContain(
			'The lines are the mean cycle time and the mean waiting of the stories completed that day'
		);
		expect(text('[data-testid="explainer-read"]')).toContain(
			"Waiting or cycle time falling while the strategic agents' hours rise is the return on their time."
		);
		expect(document.querySelector('[data-testid="strategic-ratio"]')).toBeNull();
		expect(series()).toEqual([
			'planner',
			'orchestrator',
			'analyzer',
			'mean cycle time per story',
			'mean waiting per story'
		]);
		expect(heads().slice(5)).toEqual(['completed', 'mean cycle time', 'mean waiting']);
		expect(rows()).toEqual([
			['2026-09-27', '0.5h', '-', '-', '0.5h', '0', '-', '-'],
			['2026-09-28', '0.3h', '0.5h', '-', '0.8h', '1', '2h', '1h'],
			['2026-09-29', '-', '-', '1h', '1h', '2', '3h', '1.5h']
		]);
		const type = [...document.querySelectorAll('label')]
			.find((l) => l.textContent!.trim().startsWith('type'))!
			.querySelector('select')!;
		type.value = 'task';
		type.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(api.mock.calls.at(-1)![0]).toBe('/api/stats?since=30d&type=task&bucket=day');
		// the explainer names the type the report reads (S-0485)
		expect(text('[data-testid="explainer-what"]')).toContain('of the tasks completed that day');
	});

	it('says why there is no ratio', async () => {
		answered = { items: [], strategic_days: days.map((d) => ({ ...d, completed: 0 })) };
		await open('strategic-cost');
		expect(document.querySelector('[data-testid="strategic-ratio"]')).toBeNull();
		expect(text('[data-testid="strategic-figures"]')).toBe(
			'No stories were completed in the window, so there is no ratio.'
		);
		unmount(c!);
		document.body.innerHTML = '';
		answered = {
			items,
			strategic_days: days,
			usage: { items: 0, tokens: 0, cost: 0, seconds: 0 }
		};
		await open('strategic-cost');
		expect(text('[data-testid="strategic-figures"]')).toBe(
			'The agents spent nothing in the window, so there is no ratio.'
		);
	});

	it('puts the chart, the ratio, the explainer, and the controls with the summary in order (S-0485)', async () => {
		await open('strategic-cost');
		expectOrder(
			'h1',
			chartEl(),
			'[data-testid="strategic-figures"]',
			EXPLAINER,
			WINDOW,
			'[data-testid="strategic-summary"]',
			tableView()
		);
		unmount(c!);
		document.body.innerHTML = '';
		answered = { items };
		await open('strategic-cost');
		expectOrder('h1', '[data-testid="strategic-older"]', chartEl(), EXPLAINER, WINDOW, tableView());
	});

	it('says the flai on the host is older when it sends no strategic days', async () => {
		answered = { items };
		await open('strategic-use');
		expect(text('[data-testid="strategic-older"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="strategic-figures"]')).toBeNull();
		expect(document.querySelector('[data-testid="strategic-summary"]')).toBeNull();
		expect(rows()).toEqual([]);
	});
});

describe('the cost of delay charts (S-0213)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const now = '2026-09-29T21:00:00Z';
	const cols = (backlog: number, ready: number, inProgress: number, review: number) => ({
		backlog,
		ready,
		'in-progress': inProgress,
		review
	});
	const series = (by: string, points: [string, number, string?][]) => ({
		by,
		total: points.at(-1)![1],
		points: points.map(([at, incurred, id]) => (id ? { at, id, incurred } : { at, incurred }))
	});
	const order = {
		at: now,
		horizon: '2026-10-01T21:00:00Z',
		series: [
			series('current', [
				[now, 0],
				['2026-09-30T21:00:00Z', 100, 'S-0001'],
				['2026-10-01T21:00:00Z', 450, 'S-0002']
			]),
			series('cod', [
				[now, 0],
				['2026-09-30T21:00:00Z', 250, 'S-0002'],
				['2026-10-01T21:00:00Z', 330, 'S-0001']
			]),
			series('wsjf', [
				[now, 0],
				['2026-09-30T21:00:00Z', 100, 'S-0001'],
				['2026-10-01T21:00:00Z', 360, 'S-0002']
			])
		],
		saving: 120,
		cheaper: 'cod',
		left_out: ['S-0003', 'S-0004']
	};
	const costOfDelay = {
		days: [
			{
				date: '2026-09-28',
				outstanding: cols(1000, 500, 250, 0),
				without_value: cols(3, 1, 0, 0),
				incurred: 120
			},
			{
				date: '2026-09-29',
				outstanding: cols(1200, 500, 0, 250),
				without_value: cols(2, 1, 0, 0),
				incurred: 130.5
			}
		],
		weeks: [
			{ week: '2026-W39', start: '2026-09-21', incurred: 700 },
			{ week: '2026-W40', start: '2026-09-28', incurred: 250.5 }
		],
		order
	};
	/** What flai answers; each test changes the cost of delay it carries. */
	let cod: unknown;
	beforeEach(() => {
		cod = costOfDelay;
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			return answer({ ...reportIn('day'), generated_at: now, cost_of_delay: cod });
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	const rows = (testid: string) =>
		[...document.querySelectorAll(`[data-testid="${testid}"] tbody tr`)].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent!.trim())
		);

	it('offers the window and the type to choose on a cost of delay chart', async () => {
		await open('cod-outstanding');
		expect(text('h1')).toBe('CoD Outstanding');
		expect(
			[...document.querySelectorAll('label')].map((l) => l.textContent!.trim().split(/\s+/)[0])
		).toEqual(['window', 'type']);
	});

	it('states the items without a value now per column, and lists each day in the table', async () => {
		await open('cod-outstanding');
		expect(text('[data-testid="cod-without-value"]')).toBe(
			'Items without a value now, which add nothing to this chart: backlog 2 · ready 1 · in-progress 0 · review 0.'
		);
		expect(document.querySelector('[data-testid="cod-none"]')).toBeNull();
		expect(text('[data-testid="cod-day-table"] thead')).toBe(
			'datebacklogreadyin-progressreviewwithout value'
		);
		expect(rows('cod-day-table')).toEqual([
			['2026-09-29', '1,200', '500', '0', '250', '3'],
			['2026-09-28', '1,000', '500', '250', '0', '4']
		]);
	});

	it('says none is without a value when every item has one', async () => {
		cod = {
			...costOfDelay,
			days: costOfDelay.days.map((d) => ({ ...d, without_value: cols(0, 0, 0, 0) }))
		};
		await open('cod-incurred');
		expect(text('[data-testid="cod-without-value"]')).toBe('Items without a value now: none.');
	});

	it('lists each week with its incurred cost and running mean', async () => {
		await open('cod-incurred');
		expect(text('h1')).toBe('CoD Incurred');
		expect(text('[data-testid="cod-without-value"]')).toContain('backlog 2 · ready 1');
		expect(rows('cod-week-table')).toEqual([
			['2026-W40', '2026-09-28', '250.5', '475.25'],
			['2026-W39', '2026-09-21', '700', '700']
		]);
	});

	it('states what the cheaper order saves, the three totals, the stories left out, and each pull', async () => {
		await open('cod-order');
		expect(text('[data-testid="cod-saving"]')).toBe(
			'Ordering by cost of delay would save 120 on the pull order.'
		);
		expect(text('[data-testid="cod-totals"]')).toBe(
			"Projected from now until each ready story is pulled, in the project's currency: pull order 450 · by cost of delay 330 · by WSJF 360."
		);
		expect(text('[data-testid="cod-left-out"]')).toBe(
			'Ready stories left out, without a value or a forecast duration: S-0003, S-0004.'
		);
		// the order states the counts without a value too, beside the stories it left out
		expect(text('[data-testid="cod-without-value"]')).toMatch(/^Items without a value now/);
		expect(rows('cod-order-table')).toEqual([
			['pull order', '2026-09-30 17:00 EDT', 'S-0001', '100'],
			['pull order', '2026-10-01 17:00 EDT', 'S-0002', '450'],
			['by cost of delay', '2026-09-30 17:00 EDT', 'S-0002', '250'],
			['by cost of delay', '2026-10-01 17:00 EDT', 'S-0001', '330'],
			['by WSJF', '2026-09-30 17:00 EDT', 'S-0001', '100'],
			['by WSJF', '2026-10-01 17:00 EDT', 'S-0002', '360']
		]);
	});

	it('says the pull order is already the cheapest when it saves nothing, and none left out', async () => {
		cod = { ...costOfDelay, order: { ...order, saving: -30, cheaper: 'wsjf', left_out: [] } };
		await open('cod-order');
		expect(text('[data-testid="cod-saving"]')).toBe(
			'The pull order is already the cheapest of the three.'
		);
		expect(text('[data-testid="cod-left-out"]')).toBe(
			'Ready stories left out, without a value or a forecast duration: none.'
		);
	});

	it('puts the chart, the figures it states, the explainer, the controls, and the table in order (S-0485)', async () => {
		await open('cod-order');
		expectOrder(
			'h1',
			chartEl(),
			'[data-testid="cod-saving"]',
			'[data-testid="cod-totals"]',
			'[data-testid="cod-left-out"]',
			'[data-testid="cod-without-value"]',
			EXPLAINER,
			WINDOW,
			summaryStrip(),
			tableView()
		);
		expect(tableView()!.querySelector('[data-testid="cod-order-table"]')).not.toBeNull();
		unmount(c!);
		document.body.innerHTML = '';
		cod = undefined;
		await open('cod-outstanding');
		expectOrder('h1', '[data-testid="cod-none"]', chartEl(), EXPLAINER, WINDOW, tableView());
	});

	it('says the flai on the host is older when it sends no cost of delay, order, or counts', async () => {
		const older = (missing: string) =>
			text('[data-testid="cod-none"]')?.startsWith(`The flai on the host ${missing}`) &&
			text('[data-testid="cod-none"]')?.endsWith(
				'it is older than the dashboard. Upgrade it with flai self-upgrade.'
			);
		cod = undefined;
		await open('cod-outstanding');
		expect(older('sends no cost of delay')).toBe(true);
		expect(document.querySelector('[data-testid="cod-without-value"]')).toBeNull();
		expect(rows('cod-day-table')).toEqual([]);
		unmount(c!);
		document.body.innerHTML = '';

		cod = { days: costOfDelay.days, weeks: costOfDelay.weeks };
		await open('cod-order');
		expect(older('sends no projection of the pull order')).toBe(true);
		expect(document.querySelector('[data-testid="cod-saving"]')).toBeNull();
		expect(rows('cod-order-table')).toEqual([]);
		unmount(c!);
		document.body.innerHTML = '';

		cod = {
			...costOfDelay,
			days: costOfDelay.days.map(({ date, outstanding, incurred }) => ({
				date,
				outstanding,
				incurred
			}))
		};
		await open('cod-incurred');
		expect(older('does not count the items without a value')).toBe(true);
		expect(document.querySelector('[data-testid="cod-without-value"]')).toBeNull();
		expect(rows('cod-week-table').length).toBe(2);
	});
});

describe('the charts page beside the rail (S-0484)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			const bucket = new URL(url, 'http://localhost').searchParams.get('bucket') ?? 'day';
			return answer(reportIn(bucket));
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	/** Goes to another chart on the same page, as the router does on following a link. */
	const go = async (kind: string) => {
		at.params.kind = kind;
		at.moved();
		await settle();
	};

	it('has no list of charts of its own, and draws its controls and the chart in the page', async () => {
		await open('cost-spent');
		// the rail lists the charts: the page has no nav, sidebar, or links to them
		expect(document.querySelectorAll('nav').length).toBe(0);
		expect(document.querySelector('nav[aria-label="Charts"]')).toBeNull();
		expect(document.querySelector('[data-sidebar-content]')).toBeNull();
		expect(document.querySelector('[data-sidebar-entry]')).toBeNull();
		expect(document.querySelector('a[href^="/charts/"]')).toBeNull();
		const window = document.querySelector('[data-testid="window"]')!;
		expect(window.closest('div')!.classList).toContain('flex-wrap');
		expect(document.querySelector('[data-testid="bucket"]')).not.toBeNull();
		expect(text('h1')).toBe('$ / Day');
		expect(setOption).toHaveBeenCalled();
	});

	it('draws the chart gone to, leaving the hash alone', async () => {
		await open('cycle-time');
		expect(text('h1')).toBe('Cycle Time');
		await go('tokens-spent');
		expect(text('h1')).toBe('Tokens / Day');
		expect(location.hash).toBe('');
		await go('strategic-use');
		expect(text('h1')).toBe('Strategic Use');
	});
});

describe("the chart page's order (S-0485)", () => {
	let c: ReturnType<typeof mount> | undefined;
	const done = (id: string, completed: string) => ({
		id,
		type: 'story',
		nature: 'feature',
		title: id,
		status: 'done',
		created: '2026-09-01T00:00:00Z',
		started: '2026-09-01T00:00:00Z',
		completed,
		cycle_time_seconds: 3600,
		lead_time_seconds: 7200,
		blocked_seconds: 0,
		time_in_state_seconds: { ready: 3600, 'in-progress': 3600 }
	});
	const items = [done('S-0001', '2026-09-28T21:00:00Z'), done('S-0002', '2026-09-09T21:00:00Z')];
	/** What flai answers beside the report of reportIn; each test changes it. */
	let answered: Record<string, unknown> = {};
	/** How flai answers the statistics: with the report, an error, or not at all. */
	let stats: 'report' | 'error' | 'never' = 'report';
	beforeEach(() => {
		answered = { items };
		stats = 'report';
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
		api.mockImplementation(async (url: string) => {
			if (url.startsWith('/api/items')) return answer([]);
			if (stats === 'error')
				return {
					ok: false,
					statusText: 'Bad Request',
					json: async () => ({ error: 'stats down' })
				};
			if (stats === 'never') return new Promise(() => {});
			const q = new URL(url, 'http://localhost').searchParams;
			return answer({ ...reportIn(q.get('bucket') ?? 'day'), type: q.get('type'), ...answered });
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		setOption.mockClear();
		document.body.innerHTML = '';
		chartWindow.set('30d');
	});
	const open = async (kind: string) => {
		at.params.kind = kind;
		c = mount(ChartsPage, { target: document.body });
		await settle();
	};
	const paragraph = (words: string) =>
		[...document.querySelectorAll('p')].find((p) => p.textContent?.includes(words)) ?? null;

	it('puts the title, the chart, the explainer, the controls with the summary, and the table view in order on a flow chart', async () => {
		await open('cycle-time');
		expect(text('h1')).toBe('Cycle Time');
		expectOrder('h1', chartEl(), EXPLAINER, WINDOW, summaryStrip(), tableView());
		expect(text('[data-testid="explainer"] summary')).toBe('About this chart');
		expect(document.querySelector<HTMLDetailsElement>(EXPLAINER)!.open).toBe(true);
		expect(text('[data-testid="explainer-what"]')).toContain('One point per story completed');
		expect(
			document.querySelector('[data-testid="explainer-metrics"]')!.getAttribute('href')
		).toContain('/docs/design/system/metrics.md#');
		expect(summaryStrip()!.textContent).toContain('completed 3');
		// nothing stands between the title and the chart when nothing qualifies it
		expect(document.querySelector('h1')!.nextElementSibling).toBe(chartEl());
	});

	it('keeps the share chart of Time in State with the chart, before the explainer', async () => {
		await open('time-in-state');
		expectOrder(
			'h1',
			chartEl(),
			'h2',
			chartEl(120),
			EXPLAINER,
			WINDOW,
			summaryStrip(),
			tableView()
		);
		expect(text('h2')).toBe('Share of lead time per state');
	});

	it('keeps the longest waits of Agent Waiting with the chart, before the explainer', async () => {
		answered = {
			items,
			waiting: {
				weeks: [],
				empty_wakes: { count: 0 },
				longest: [
					{
						item: 'S-0001',
						kind: 'thread',
						thread: 'TH-0003',
						started: '2026-09-29T09:00:00Z',
						seconds: 3600
					}
				]
			}
		};
		await open('agent-waiting');
		expectOrder(
			'h1',
			chartEl(),
			'h2',
			'[data-testid="wait-table"]',
			EXPLAINER,
			WINDOW,
			tableView()
		);
		expect(tableView()!.querySelector('[data-testid="wait-table"]')).toBeNull();
	});

	it('names the type and the bucket chosen in the explainer of a spend chart, its controls after it', async () => {
		await open('tokens-spent');
		await choose('bucket', 'week');
		expect(text('h1')).toBe('Tokens / Week');
		expectOrder(
			'h1',
			chartEl(),
			EXPLAINER,
			WINDOW,
			'[data-testid="bucket"]',
			'[data-testid="usage-summary"]',
			tableView()
		);
		expect(text('[data-testid="explainer-what"]')).toBe(
			'One bar per week: the tokens spent on the stories done in it, stacked by model, with the mean per week so far as a dashed line.'
		);
		expect(text('[data-testid="explainer-how"]')).toContain(
			'Each week holds the items that entered done in it, in UTC'
		);
		expect(text('[data-testid="explainer-how"]')).toContain('The dashed line is the mean per week');
		expect(document.querySelector('[data-testid="spend-note"]')).toBeNull();
		const type = [...document.querySelectorAll('label')]
			.find((l) => l.textContent!.trim().startsWith('type'))!
			.querySelector('select')!;
		type.value = 'task';
		type.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(text('[data-testid="explainer-what"]')).toContain(
			'the tokens spent on the tasks done in it'
		);
	});

	it('puts the notice between the title and the chart', async () => {
		const r = reportIn('day');
		answered = { items, usage: { ...r.usage, bucket: undefined, spend: undefined } };
		await open('token-rate');
		expectOrder('h1', '[data-testid="spend-none"]', chartEl(), EXPLAINER, WINDOW, tableView());
		expect(document.querySelector('[data-testid="spend-note"]')).toBeNull();
	});

	it('shows the controls while the report loads and when flai answers an error', async () => {
		stats = 'never';
		await open('cycle-time');
		expectOrder(paragraph('Loading…'), WINDOW);
		expect(document.querySelector('h1')).toBeNull();
		expect(document.querySelector(EXPLAINER)).toBeNull();
		expect(tableView()).toBeNull();
		unmount(c!);
		document.body.innerHTML = '';
		stats = 'error';
		await open('cycle-time');
		expectOrder(paragraph('stats down'), WINDOW);
		expect(document.querySelector('h1')).toBeNull();
		expect(document.querySelector(EXPLAINER)).toBeNull();
		expect(tableView()).toBeNull();
		// the window can still be changed, and asks flai again
		stats = 'report';
		await choose('window', '7d');
		expect(api.mock.calls.at(-1)![0]).toBe('/api/stats?since=7d&type=story&bucket=day');
		expectOrder('h1', chartEl(), EXPLAINER, WINDOW, tableView());
	});
});
