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

	it('lists the flow charts and the usage charts apart, and draws tokens per day with its table', async () => {
		await open('tokens-spent');
		expect(
			[...document.querySelectorAll('[data-testid="charts-flow"] a')].map((l) => l.textContent)
		).toEqual([
			'Cycle Time',
			'Burn-up',
			'Cumulative Flow',
			'Time in State',
			'Throughput',
			'Agent Waiting'
		]);
		expect(
			[...document.querySelectorAll('[data-testid="charts-usage"] a')].map((l) => l.textContent)
		).toEqual([
			'Tokens / Min',
			'Tokens / Day',
			'Tokens per item',
			'Tokens / $',
			'$ / Day',
			'$ / Work Type',
			'$ / Item',
			'Avg. Time / Model',
			'Avg. Cost / Model'
		]);
		expect(text('h1')).toBe('Tokens / Day');
		expect(asked()).toEqual(['/api/stats?since=30d&type=story&bucket=day']);
		// no epic narrows what flai summed
		expect(controls()).toEqual(['window', 'per', 'type']);
		expect(text('[data-testid="usage-summary"]')).toContain(
			'per story 12.0M tokens, $5.00 · 1.2M tokens per agent minute · 2.4M tokens per dollar'
		);
		expect(text('[data-testid="spend-note"]')).toContain('The dashed line is the mean per day');
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
		expect(text('[data-testid="spend-note"]')).toContain('Each day holds the items');
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

	it('lists the planning charts as a third group, each opening its chart', async () => {
		await open('forecast-accuracy');
		const links = [...document.querySelectorAll('[data-testid="charts-planning"] a')];
		expect(links.map((l) => [l.textContent, l.getAttribute('href')])).toEqual([
			['Forecast Accuracy', '/charts/forecast-accuracy'],
			['Delivery Accuracy', '/charts/delivery-accuracy'],
			['Forecast Error / Model', '/charts/forecast-by-model'],
			['Parallelism', '/charts/parallelism'],
			['Hold Time', '/charts/hold-time'],
			['Touches Drift', '/charts/touches-drift'],
			['CoD Outstanding', '/charts/cod-outstanding'],
			['CoD Incurred', '/charts/cod-incurred'],
			['CoD by Order', '/charts/cod-order']
		]);
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

	it('lists the strategic charts as a fourth group and states the ratio on Strategic Cost', async () => {
		await open('strategic-cost');
		const links = [...document.querySelectorAll('[data-testid="charts-strategic"] a')];
		expect(links.map((l) => [l.textContent, l.getAttribute('href')])).toEqual([
			['Strategic Cost', '/charts/strategic-cost'],
			['Strategic Use', '/charts/strategic-use']
		]);
		const navs = [...document.querySelectorAll('nav')].map((n) => n.getAttribute('data-testid'));
		expect(navs).toEqual(['charts-flow', 'charts-usage', 'charts-planning', 'charts-strategic']);
		expect(text('h1')).toBe('Strategic Cost');
		expect(controls()).toEqual(['window', 'type']);
		expect(text('[data-testid="strategic-summary"]')).toBe(
			'strategic agents $1.50 · 2.3h · 3 stories completed'
		);
		expect(text('[data-testid="strategic-note"]')).toContain(
			'The bars are what planning, orchestration, and analysis cost each day.'
		);
		expect(text('[data-testid="strategic-ratio"]')).toBe(
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
		expect(text('[data-testid="strategic-note"]')).toBe(
			"The bars are the hours the planner, the orchestrator, and the analyzer worked each day. The lines are the mean cycle time and the mean waiting of the stories completed that day. Waiting or cycle time falling while the strategic agents' hours rise is the return on their time."
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
		expect(text('[data-testid="strategic-note"]')).toContain('of the tasks completed that day');
	});

	it('says why there is no ratio', async () => {
		answered = { items: [], strategic_days: days.map((d) => ({ ...d, completed: 0 })) };
		await open('strategic-cost');
		expect(document.querySelector('[data-testid="strategic-ratio"]')).toBeNull();
		expect(text('[data-testid="strategic-note"]')).toContain(
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
		expect(text('[data-testid="strategic-note"]')).toContain(
			'The agents spent nothing in the window, so there is no ratio.'
		);
	});

	it('says the flai on the host is older when it sends no strategic days', async () => {
		answered = { items };
		await open('strategic-use');
		expect(text('[data-testid="strategic-older"]')).toContain('flai self-upgrade');
		expect(document.querySelector('[data-testid="strategic-note"]')).toBeNull();
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

	it('lists the three cost of delay charts under planning, after the forecast and claims charts, with the window and the type to choose', async () => {
		await open('cod-outstanding');
		expect(
			[...document.querySelectorAll('[data-testid="charts-planning"] a')].map((l) => l.textContent)
		).toEqual([
			'Forecast Accuracy',
			'Delivery Accuracy',
			'Forecast Error / Model',
			'Parallelism',
			'Hold Time',
			'Touches Drift',
			'CoD Outstanding',
			'CoD Incurred',
			'CoD by Order'
		]);
		expect(
			[...document.querySelectorAll('nav')].map((n) => [
				n.dataset.testid,
				n.classList.contains('mb-3')
			])
		).toEqual([
			['charts-flow', false],
			['charts-usage', false],
			['charts-planning', false],
			['charts-strategic', true]
		]);
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
