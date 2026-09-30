// S-0163: the charts page offers the charts of spend over time, each with the controls it uses
// and a table of what it plots.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const at = vi.hoisted(() => ({ params: { kind: 'tokens-spent' } }));
vi.mock('$app/state', () => ({ page: at }));
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
		done: [],
		by_model: {},
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
		).toEqual(['Cycle Time', 'Burn-up', 'Cumulative Flow', 'Time in State', 'Throughput']);
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
			'Completion over time',
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
		expect(text('[data-testid="spend-table"] tbody td')).toBe('2026-09-29 19:00 UTC');
		const since = document.querySelector<HTMLSelectElement>('label select')!;
		since.value = '90d';
		since.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(asked().at(-1)).toBe('/api/stats?since=90d&type=story&bucket=day');
		const offered = [...document.querySelectorAll('[data-testid="bucket"] option')];
		expect(offered.map((o) => o.textContent)).toEqual(['day', 'week']);
	});

	it('compares the item types per item, or the models on one type', async () => {
		await open('cost-per-item');
		expect(controls()).toEqual(['window', 'per', 'compare']);
		let drawn = setOption.mock.calls.at(-1)![0] as { series: { name: string }[] };
		expect(drawn.series.map((s) => s.name)).toEqual(['story', 'task']);
		const rows = () =>
			[...document.querySelectorAll('[data-testid="spend-table"] tbody tr')].map((tr) =>
				[...tr.querySelectorAll('td')].slice(1, 3).map((td) => td.textContent?.trim())
			);
		expect(rows()).toEqual([
			['story', '3'],
			['task', '8']
		]);
		await choose('by', 'model');
		expect(controls()).toEqual(['window', 'per', 'compare', 'type']);
		drawn = setOption.mock.calls.at(-1)![0] as { series: { name: string }[] };
		expect(drawn.series.map((s) => s.name)).toEqual(['claude-opus-5-5']);
		expect(rows().map((r) => r[0])).toEqual(['story', 'claude-opus-5-5']);
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

	it('keeps the epic and the table per item on the charts of S-0143', async () => {
		await open('cost');
		expect(text('h1')).toBe('$ / Item');
		expect(controls()).toEqual(['window', 'type', 'epic']);
		expect(text('[data-testid="usage-table"] thead')).toContain('tokens/agent minute');
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
			// cycle time from the window's start; time in state from half a day before the day that holds it
			expect(onTime().xAxis.min).toBe(
				kind === 'cycle-time'
					? Date.parse('2026-09-22T21:00:00Z')
					: Date.parse('2026-09-22T00:00:00Z') - 12 * 3600e3
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
		expect(onTime().xAxis.min).toBe(Date.parse('2026-08-30T00:00:00Z') - half);
		expect(onTime().xAxis.max).toBe(Date.parse('2026-09-29T00:00:00Z') + half);
		const days = () =>
			(onTime().series[0].data as unknown as { ids: string[] }[]).map((d) => d.ids);
		expect(days()).toEqual([['S-0002'], ['S-0001']]);
		await choose('window', '7d');
		expect(onTime().xAxis.min).toBe(Date.parse('2026-09-22T00:00:00Z') - half);
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
