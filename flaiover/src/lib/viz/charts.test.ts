import { describe, expect, it } from 'vitest';
import {
	aging,
	build,
	burnUp,
	cfd,
	completionCost,
	completionTime,
	completedIn,
	bucketLabel,
	bucketsFor,
	controls,
	cost,
	costPerItem,
	costSpent,
	cycleTime,
	estimates,
	human,
	KINDS,
	normalise,
	stateShare,
	throughput,
	timeInState,
	hasSpend,
	titleOf,
	tokenRate,
	tokensPerDollar,
	tokensPerItem,
	tokensSpent,
	models,
	withUsage,
	type Bucket,
	type Report
} from './charts';
import { CATEGORICAL, modelSlot, modelSymbol, theme, TYPE_SLOT } from './palette';

// Spend over time as flai lays it out (S-0163): stories on 3 August (two models) and on
// 5 August, the day between them empty, and tasks on 3 August.
const storyDays: Bucket[] = [
	{
		at: '2026-08-03T00:00:00Z',
		items: 2,
		tokens: 3500000,
		cost: 1.75,
		seconds: 5400,
		estimated: true,
		tokens_per_item: 1750000,
		cost_per_item: 0.875,
		tokens_per_minute: 38888.9,
		tokens_per_dollar: 2000000,
		mean_tokens: 3500000,
		mean_cost: 1.75,
		models: [
			{
				model: 'claude-haiku-4-5',
				items: 1,
				tokens: 1000000,
				cost: 0.25,
				seconds: 3600,
				tokens_per_item: 1000000,
				cost_per_item: 0.25,
				tokens_per_minute: 16666.7,
				tokens_per_dollar: 4000000
			},
			{
				model: 'claude-opus-5-5',
				items: 2,
				tokens: 2500000,
				cost: 1.5,
				seconds: 5400,
				estimated: true,
				tokens_per_item: 1250000,
				cost_per_item: 0.75,
				tokens_per_minute: 27777.8,
				tokens_per_dollar: 1666666.7
			}
		]
	},
	{
		at: '2026-08-04T00:00:00Z',
		items: 0,
		tokens: 0,
		cost: 0,
		seconds: 0,
		mean_tokens: 1750000,
		mean_cost: 0.875
	},
	{
		at: '2026-08-05T00:00:00Z',
		items: 1,
		tokens: 1000000,
		cost: 0.5,
		seconds: 0,
		tokens_per_item: 1000000,
		cost_per_item: 0.5,
		tokens_per_dollar: 2000000,
		mean_tokens: 1500000,
		mean_cost: 0.75,
		models: [
			{
				model: 'claude-opus-5-5',
				items: 1,
				tokens: 1000000,
				cost: 0.5,
				seconds: 0,
				tokens_per_item: 1000000,
				cost_per_item: 0.5,
				tokens_per_dollar: 2000000
			}
		]
	}
];
const taskDays: Bucket[] = [
	{
		at: '2026-08-03T00:00:00Z',
		items: 4,
		tokens: 2000000,
		cost: 1,
		seconds: 1200,
		tokens_per_item: 500000,
		cost_per_item: 0.25,
		tokens_per_minute: 100000,
		tokens_per_dollar: 2000000,
		mean_tokens: 2000000,
		mean_cost: 1,
		models: [
			{
				model: 'claude-opus-5-5',
				items: 4,
				tokens: 2000000,
				cost: 1,
				seconds: 1200,
				tokens_per_item: 500000,
				cost_per_item: 0.25,
				tokens_per_minute: 100000,
				tokens_per_dollar: 2000000
			}
		]
	}
];
const none = { items: 0, tokens: 0, cost: 0, seconds: 0, models: [], buckets: [] };

const report: Report = {
	generated_at: '2026-09-01T12:00:00Z',
	type: 'story',
	window_days: 30,
	window_start: '2026-08-02T12:00:00Z',
	summary: {
		completed: 2,
		cancelled: 1,
		cancellation_rate: 1 / 3,
		throughput_per_week: 0.47,
		wip: 1,
		cycle_time: {
			count: 2,
			p50_seconds: 86400,
			p85_seconds: 93600,
			max_seconds: 93600,
			mean_seconds: 90000
		},
		lead_time: {
			count: 2,
			p50_seconds: 174600,
			p85_seconds: 183600,
			max_seconds: 183600,
			mean_seconds: 179100
		},
		queue_time: {
			count: 2,
			p50_seconds: 86400,
			p85_seconds: 86400,
			max_seconds: 86400,
			mean_seconds: 86400
		},
		flow_efficiency: 0.88,
		time_in_state_share: { backlog: 0.015, ready: 0.48, 'in-progress': 0.34, review: 0.16 }
	},
	items: [
		{
			id: 'S-001',
			type: 'story',
			nature: 'feature',
			title: 'One',
			status: 'done',
			parent: 'E-001',
			created: '2026-08-01T09:00:00Z',
			started: '2026-08-02T10:00:00Z',
			completed: '2026-08-03T12:00:00Z',
			lead_time_seconds: 183600,
			cycle_time_seconds: 93600,
			queue_time_seconds: 86400,
			blocked_seconds: 21600,
			time_in_state_seconds: { backlog: 3600, ready: 86400, 'in-progress': 86400, review: 7200 },
			estimate_seconds: 72000,
			estimate_error: 0.3,
			usage: {
				source: 'log',
				tokens: 3000000,
				cost: 1.5,
				seconds: 3600,
				tokens_per_hour: 3000000,
				models: [
					{ model: 'claude-opus-5-5', tokens: 2000000, cost: 1.2, tokens_per_hour: 2000000 },
					{ model: 'claude-haiku-4-5', tokens: 1000000, cost: 0.3, tokens_per_hour: 1000000 }
				]
			}
		},
		{
			id: 'S-002',
			type: 'story',
			nature: 'improvement',
			title: 'Two',
			status: 'done',
			parent: 'E-001',
			created: '2026-08-10T09:00:00Z',
			started: '2026-08-11T09:30:00Z',
			completed: '2026-08-12T09:30:00Z',
			lead_time_seconds: 174600,
			cycle_time_seconds: 86400,
			blocked_seconds: 0,
			time_in_state_seconds: { backlog: 1800, ready: 86400, 'in-progress': 36000, review: 50400 },
			usage: {
				source: 'sum',
				tokens: 500000,
				cost: 0.25,
				seconds: 1800,
				tokens_per_hour: 1000000,
				estimated: true,
				models: [{ model: 'claude-opus-5-5', tokens: 500000, cost: 0.25, tokens_per_hour: 1000000 }]
			}
		},
		{
			id: 'S-004',
			type: 'story',
			nature: 'feature',
			title: 'Four',
			status: 'in-progress',
			parent: 'E-001',
			created: '2026-08-25T09:00:00Z',
			started: '2026-08-31T10:00:00Z',
			blocked_seconds: 0,
			time_in_state_seconds: {},
			age_seconds: 93600
		}
	],
	throughput: [
		{ week: '2026-W32', start: '2026-08-03', done: 1, by_nature: { feature: 1 } },
		{ week: '2026-W33', start: '2026-08-10', done: 1, by_nature: { improvement: 1 } }
	],
	burnup: {
		all: [
			{ date: '2026-08-01', scope: 1, done: 0 },
			{ date: '2026-09-01', scope: 4, done: 3 }
		],
		'E-001': [{ date: '2026-09-01', scope: 4, done: 3 }]
	},
	cfd: [
		{
			date: '2026-09-01',
			counts: { backlog: 0, ready: 0, 'in-progress': 1, review: 0, done: 3, cancelled: 1 }
		}
	],
	aging: [
		{
			id: 'S-004',
			title: 'Four',
			status: 'in-progress',
			age_seconds: 93600,
			over_p85: false,
			blocked: false,
			nature: 'feature',
			parent: 'E-001',
			age: '1d2h'
		}
	],
	usage: {
		items: 2,
		tokens: 3500000,
		cost: 1.75,
		seconds: 5400,
		estimated: true,
		models: [
			{ model: 'claude-haiku-4-5', tokens: 1000000, cost: 0.3, items: 1 },
			{ model: 'claude-opus-5-5', tokens: 2500000, cost: 1.45, items: 2 }
		],
		done: [
			{ at: '2026-08-03T12:00:00Z', id: 'S-001', done: 1, tokens: 3000000, cost: 1.5 },
			{ at: '2026-08-12T09:30:00Z', id: 'S-002', done: 2, tokens: 3500000, cost: 1.75 }
		],
		bucket: 'day',
		spend: {
			epic: none,
			story: {
				items: 3,
				tokens: 4500000,
				cost: 2.25,
				seconds: 5400,
				models: [],
				buckets: storyDays
			},
			task: { items: 4, tokens: 2000000, cost: 1, seconds: 1200, models: [], buckets: taskDays }
		},
		by_model: {
			'claude-haiku-4-5': [
				{ at: '2026-08-03T12:00:00Z', id: 'S-001', done: 1, tokens: 1000000, cost: 0.3 }
			],
			'claude-opus-5-5': [
				{ at: '2026-08-03T12:00:00Z', id: 'S-001', done: 1, tokens: 2000000, cost: 1.2 },
				{ at: '2026-08-12T09:30:00Z', id: 'S-002', done: 2, tokens: 2500000, cost: 1.45 }
			]
		}
	}
};

const light = theme(false);
const dark = theme(true);

describe('chart builders', () => {
	it('every kind builds with one y-axis and a tooltip', () => {
		for (const k of KINDS) {
			const o = build(k, report, light) as { yAxis: unknown; tooltip: unknown; series: unknown[] };
			expect(o.yAxis, k).toBeDefined();
			expect(Array.isArray(o.yAxis), `${k} must not use two y-axes`).toBe(false);
			expect(o.tooltip, k).toBeDefined();
			expect(o.series.length, k).toBeGreaterThan(0);
		}
	});
	it("every chart spans the report's window and plots only the items completed in it (S-0166)", () => {
		type Axis = { xAxis: { min?: number; max?: number; data?: string[] } };
		const end = Date.parse('2026-09-01T12:00:00Z');
		const hour = 3600e3;
		// the report's own window, 2 August 12:00 to 1 September 12:00
		const month = Date.parse('2026-08-02T12:00:00Z');
		expect((cycleTime(report, light) as Axis).xAxis).toMatchObject({ min: month, max: end });
		expect((completionTime(report, light) as Axis).xAxis).toMatchObject({ min: month, max: end });
		// a series by the day starts on the day that holds the window's start
		const day = Date.parse('2026-08-02T00:00:00Z');
		expect((burnUp(report, light) as Axis).xAxis).toMatchObject({ min: day, max: end });
		expect((cfd(report, light) as Axis).xAxis).toMatchObject({ min: day, max: end });
		// time in state by the day (S-0168): from the day that holds the start to the one that holds
		// now, half a day either side
		expect((timeInState(report, light) as Axis).xAxis).toMatchObject({
			type: 'time',
			min: day - 12 * hour,
			max: Date.parse('2026-09-01T00:00:00Z') + 12 * hour
		});
		// spend over time: from the bucket that holds the start to the one that holds now, half a
		// bucket either side
		expect((tokensSpent(report, light) as Axis).xAxis).toMatchObject({
			min: day - 12 * hour,
			max: Date.parse('2026-09-01T00:00:00Z') + 12 * hour
		});
		const weekly = { ...report, usage: { ...report.usage!, bucket: 'week' as const } };
		expect((costSpent(weekly, light) as Axis).xAxis).toMatchObject({
			min: Date.parse('2026-07-27T00:00:00Z') - 84 * hour,
			max: Date.parse('2026-08-31T00:00:00Z') + 84 * hour
		});

		// a narrower window moves the axis and drops S-001, completed on 3 August
		const narrow: Report = { ...report, window_days: 7, window_start: '2026-08-10T00:00:00Z' };
		const from = Date.parse('2026-08-10T00:00:00Z');
		const ct = cycleTime(narrow, light) as Axis & { series: { data: { id: string }[] }[] };
		expect(ct.xAxis).toMatchObject({ min: from, max: end });
		expect(ct.series.flatMap((s) => s.data.map((d) => d.id))).toEqual(['S-002']);
		const tis = timeInState(narrow, light) as Axis & { series: { data: { ids: string[] }[] }[] };
		expect(tis.xAxis.min).toBe(from - 12 * hour);
		expect(tis.series[0].data.map((d) => d.ids)).toEqual([['S-002']]);
		expect((cost(narrow, light) as Axis).xAxis.data).toEqual(['S-002*']);
		const es = estimates(narrow, light) as { series: { data: unknown[] }[] };
		expect(es.series[0].data).toEqual([]);
		expect((tokenRate(narrow, light) as Axis).xAxis.min).toBe(from - 12 * hour);
		expect(completedIn(narrow).map((i) => i.id)).toEqual(['S-002']);
		expect(withUsage(narrow).map((i) => i.id)).toEqual(['S-002']);
	});
	it('cycle time carries p50 and p85 reference lines and colours by nature in fixed slots', () => {
		const o = cycleTime(report, light) as {
			series: {
				name: string;
				itemStyle: { color: string };
				markLine?: { data: { name: string; yAxis: number }[] };
			}[];
			legend: { show: boolean };
		};
		expect(o.series.map((s) => s.name)).toEqual(['feature', 'improvement']);
		expect(o.series[0].itemStyle.color).toBe(CATEGORICAL.light[0]);
		expect(o.series[1].itemStyle.color).toBe(CATEGORICAL.light[5]);
		expect(o.series[0].markLine?.data.map((d) => [d.name, d.yAxis])).toEqual([
			['p50', 1],
			['p85', 1.08]
		]);
		expect(o.legend.show).toBe(true);
		const one = cycleTime(report, light, 'E-999') as { series: { data: unknown[] }[] };
		expect(one.series.length).toBe(0);
	});
	it('burn-up uses the epic series when asked and falls back to all', () => {
		const e = burnUp(report, light, 'E-001') as { series: { data: unknown[] }[] };
		expect(e.series[0].data.length).toBe(1);
		const all = burnUp(report, light, 'E-404') as { series: { data: unknown[] }[] };
		expect(all.series[0].data.length).toBe(2);
	});
	it('cfd and time in state stack the states with the surface as the seam', () => {
		const c = cfd(report, light) as { series: { stack: string; lineStyle: { color: string } }[] };
		expect(new Set(c.series.map((s) => s.stack)).size).toBe(1);
		expect(c.series[0].lineStyle.color).toBe(light.surface);
		type Bars = { series: { name: string; data: { value: [number, number]; ids: string[] }[] }[] };
		const t = timeInState(report, light) as Bars;
		// a bar per day with items completed, at the day's start in UTC
		expect(t.series[2].name).toBe('in-progress');
		expect(t.series[2].data).toEqual([
			{ value: [Date.parse('2026-08-03T00:00:00Z'), 24], ids: ['S-001'] },
			{ value: [Date.parse('2026-08-12T00:00:00Z'), 10], ids: ['S-002'] }
		]);
		// items completed the same day share a bar: the mean hours per state
		const sameDay: Report = {
			...report,
			items: report.items.map((i) =>
				i.id === 'S-002' ? { ...i, completed: '2026-08-03T20:00:00Z' } : i
			)
		};
		const m = timeInState(sameDay, light) as Bars;
		expect(m.series.map((s) => s.data.map((d) => d.value[1]))).toEqual([
			[0.8],
			[24],
			[17],
			[8],
			[0]
		]);
		expect(m.series[0].data[0].ids).toEqual(['S-001', 'S-002']);
		expect((timeInState(report, light, 'E-999') as Bars).series[0].data).toEqual([]);
		const share = stateShare(report, light) as { series: { data: number[] }[] };
		expect(share.series[1].data[0]).toBeCloseTo(0.48);
	});
	it('throughput, aging, and estimates shape their data', () => {
		const th = throughput(report, light) as { series: { name: string; data: number[] }[] };
		expect(th.series.map((s) => s.name)).toEqual(['feature', 'improvement']);
		expect(th.series[0].data).toEqual([1, 0]);
		const ag = aging(report, dark) as {
			yAxis: { data: string[] };
			series: { data: { value: number }[] }[];
		};
		expect(ag.yAxis.data).toEqual(['S-004']);
		expect(ag.series[0].data[0].value).toBe(1.08);
		const es = estimates(report, light) as { series: { data: { value: number[] }[] }[] };
		expect(es.series[0].data[0].value).toEqual([20, 26]);
	});
	it('usage charts colour each model in a fixed slot by name, whatever a filter leaves', () => {
		expect(models(report)).toEqual(['claude-haiku-4-5', 'claude-opus-5-5']);
		const c = cost(report, light) as {
			xAxis: { data: string[] };
			series: { name: string; stack: string; data: number[]; itemStyle: { color: string } }[];
		};
		expect(c.xAxis.data).toEqual(['S-001', 'S-002*']);
		expect(c.series.map((s) => [s.name, s.data])).toEqual([
			['claude-haiku-4-5', [0.3, 0]],
			['claude-opus-5-5', [1.2, 0.25]]
		]);
		expect(new Set(c.series.map((s) => s.stack)).size).toBe(1);
		// an epic with only opus left: opus keeps its colour
		const one = cost({ ...report, items: [report.items[1]] }, light) as typeof c;
		expect(one.series.map((s) => s.name)).toEqual(['claude-opus-5-5']);
		expect(one.series[0].itemStyle.color).toBe(CATEGORICAL.light[0]);
		expect(modelSlot('claude-sonnet-5')).toBe(2);
		expect([3, 4]).toContain(modelSlot('gpt-9'));
	});
	type Line = {
		name: string;
		type: string;
		symbol?: string;
		stack?: string;
		itemStyle: { color: string };
		lineStyle?: { type: string; color: string };
		data: { value: [string, number]; items: number; estimated?: boolean }[];
	};
	type Over = {
		series: Line[];
		legend: { show: boolean };
		useUTC: boolean;
		xAxis: { type: string; min?: number; max?: number; minInterval: number };
		yAxis: { name: string; axisLabel: { formatter: (v: number) => string } };
		tooltip: { trigger: string; formatter: (p: unknown) => string };
	};
	const values = (l: Line) => l.data.map((d) => d.value);
	it('token rate is tokens per agent minute over time, per model, with all of them dashed', () => {
		const o = tokenRate(report, light) as Over;
		expect(o.series.map((s) => s.name)).toEqual([
			'claude-haiku-4-5',
			'claude-opus-5-5',
			'all models'
		]);
		expect(o.xAxis.type).toBe('time');
		// buckets are UTC, and the axis spans the window's days, not the one day drawn (S-0166)
		expect(o.useUTC).toBe(true);
		expect([o.xAxis.min, o.xAxis.max, o.xAxis.minInterval]).toEqual([
			Date.parse('2026-08-01T12:00:00Z'),
			Date.parse('2026-09-01T12:00:00Z'),
			86400e3
		]);
		expect(o.yAxis.name).toBe('tokens per agent minute');
		expect(o.yAxis.axisLabel.formatter(38888.9)).toBe('38.9K');
		// the day with nothing done has no point, and neither has the one with no agent time
		expect(values(o.series[1])).toEqual([['2026-08-03T00:00:00Z', 27777.8]]);
		expect(values(o.series[2])).toEqual([['2026-08-03T00:00:00Z', 38888.9]]);
		expect(o.series[1].itemStyle.color).toBe(CATEGORICAL.light[0]);
		expect(o.series[0].itemStyle.color).toBe(CATEGORICAL.light[5]);
		expect([o.series[0].symbol, o.series[1].symbol]).toEqual(['triangle', 'circle']);
		expect(o.series[2].lineStyle?.type).toBe('dashed');
		expect(o.legend.show).toBe(true);
		expect(
			o.tooltip.formatter([{ seriesName: 'claude-opus-5-5', data: o.series[1].data[0] }])
		).toBe(
			'2026-08-03<br/>claude-opus-5-5: 27.8K tokens per minute over 2 items (estimated in part)'
		);
		// one model: its line alone, and no legend
		const tasks = tokenRate({ ...report, type: 'task' }, light) as Over;
		expect(tasks.series.map((s) => s.name)).toEqual(['claude-opus-5-5']);
		expect(tasks.legend.show).toBe(false);
	});
	it('tokens and cost per bucket stack the models and carry the running mean', () => {
		const o = tokensSpent(report, light) as Over;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['claude-haiku-4-5', 'bar'],
			['claude-opus-5-5', 'bar'],
			['mean per day', 'line']
		]);
		expect(new Set(o.series.slice(0, 2).map((s) => s.stack))).toEqual(new Set(['spent']));
		// every model has a value in every bucket, so the stack lines up
		expect(values(o.series[0]).map((v) => v[1])).toEqual([1000000, 0, 0]);
		expect(values(o.series[1]).map((v) => v[1])).toEqual([2500000, 0, 1000000]);
		expect(values(o.series[2]).map((v) => v[1])).toEqual([3500000, 1750000, 1500000]);
		expect(o.legend.show).toBe(true);
		const c = costSpent(report, dark) as Over;
		expect(c.yAxis.name).toBe('US dollars');
		expect(values(c.series[1]).map((v) => v[1])).toEqual([1.5, 0, 0.5]);
		expect(values(c.series[2]).map((v) => v[1])).toEqual([1.75, 0.875, 0.75]);
		expect(c.series[1].itemStyle.color).toBe(CATEGORICAL.dark[0]);
		expect(
			c.tooltip.formatter([{ seriesName: 'claude-opus-5-5', data: c.series[1].data[0] }])
		).toBe('2026-08-03<br/>claude-opus-5-5: $1.50 (estimated in part)');
		const week = { ...report, usage: { ...report.usage!, bucket: 'week' as const } };
		expect((tokensSpent(week, light) as Over).series[2].name).toBe('mean per week');
		expect(titleOf('tokens-spent', 'week')).toBe('Tokens per week');
		expect(titleOf('cost-spent', 'hour')).toBe('Cost per hour');
		expect(titleOf('cost-spent')).toBe('Cost per day');
		expect(titleOf('token-rate', 'hour')).toBe('Token rate');
	});
	it('tokens and cost per item compare the types, or the models on the type of the report', () => {
		const o = tokensPerItem(report, light) as Over;
		// no epic carries usage: it has no line
		expect(o.series.map((s) => s.name)).toEqual(['story', 'task']);
		expect(values(o.series[0])).toEqual([
			['2026-08-03T00:00:00Z', 1750000],
			['2026-08-05T00:00:00Z', 1000000]
		]);
		expect(values(o.series[1])).toEqual([['2026-08-03T00:00:00Z', 500000]]);
		expect(o.series[0].itemStyle.color).toBe(CATEGORICAL.light[TYPE_SLOT.story]);
		expect(o.series[1].itemStyle.color).toBe(CATEGORICAL.light[TYPE_SLOT.task]);
		expect(o.series[0].symbol).not.toBe(o.series[1].symbol);
		expect(o.yAxis.name).toBe('tokens per item');
		const c = costPerItem(report, light, 'model') as Over;
		expect(c.series.map((s) => s.name)).toEqual(['claude-haiku-4-5', 'claude-opus-5-5']);
		expect(values(c.series[1])).toEqual([
			['2026-08-03T00:00:00Z', 0.75],
			['2026-08-05T00:00:00Z', 0.5]
		]);
		expect(c.yAxis.name).toBe('US dollars per story');
		expect(c.yAxis.axisLabel.formatter(0.75)).toBe('$0.750');
		expect(
			c.tooltip.formatter([{ seriesName: 'claude-haiku-4-5', data: c.series[0].data[0] }])
		).toBe('2026-08-03<br/>claude-haiku-4-5: $0.250 each over 1 item');
		expect((build('cost-per-item', report, light, undefined, 'model') as Over).series.length).toBe(
			2
		);
		expect((build('cost-per-item', report, light) as Over).series.map((s) => s.name)).toEqual([
			'story',
			'task'
		]);
	});
	it('tokens per dollar is what a dollar bought, per model', () => {
		const o = tokensPerDollar(report, light) as Over;
		expect(values(o.series[1])).toEqual([
			['2026-08-03T00:00:00Z', 1666666.7],
			['2026-08-05T00:00:00Z', 2000000]
		]);
		expect(values(o.series[2]).map((v) => v[1])).toEqual([2000000, 2000000]);
		expect(o.yAxis.name).toBe('tokens per US dollar');
	});
	it('offers each chart the controls it uses, and an hour over 31 days or less', () => {
		expect(bucketsFor('7d')).toEqual(['hour', 'day', 'week']);
		expect(bucketsFor('30d')).toEqual(['hour', 'day', 'week']);
		expect(bucketsFor('90d')).toEqual(['day', 'week']);
		expect(bucketsFor('12w')).toEqual(['day', 'week']);
		expect(controls('tokens-per-item', 'type')).toEqual({
			type: false,
			epic: false,
			bucket: true,
			by: true
		});
		expect(controls('tokens-per-item', 'model').type).toBe(true);
		expect(controls('token-rate', 'type')).toEqual({
			type: true,
			epic: false,
			bucket: true,
			by: false
		});
		expect(controls('cost', 'type')).toEqual({ type: true, epic: true, bucket: false, by: false });
		expect(controls('cfd', 'type').epic).toBe(false);
		expect(controls('cycle-time', 'model').epic).toBe(true);
	});
	it('names a bucket as a reader does, and says when a flai sends no spend', () => {
		expect(bucketLabel('2026-09-29T19:00:00Z', 'hour')).toBe('2026-09-29 19:00 UTC');
		expect(bucketLabel('2026-09-28T00:00:00Z', 'week')).toBe('week of 2026-09-28');
		expect(modelSymbol('claude-fable-5-1')).toBe('diamond');
		expect(modelSymbol('gpt-9')).toBe('roundRect');
		expect(hasSpend(normalise(report))).toBe(true);
		const older = normalise({
			...report,
			usage: { ...report.usage!, bucket: undefined, spend: undefined }
		});
		expect(hasSpend(older)).toBe(false);
		for (const kind of [
			'token-rate',
			'tokens-spent',
			'tokens-per-item',
			'tokens-per-dollar',
			'cost-spent',
			'cost-per-item'
		] as const) {
			const o = build(kind, older, light) as Over;
			expect(o.series, kind).toEqual([]);
			expect(o.legend.show, kind).toBe(false);
		}
		// a flai that sends a type with null lists
		const nulls = normalise({
			...report,
			usage: {
				...report.usage!,
				spend: { story: { items: 0, tokens: 0, cost: 0, seconds: 0 } }
			}
		} as unknown as Report);
		expect(nulls.usage?.spend?.story.buckets).toEqual([]);
	});
	it("completion charts lay out each model's items done against time and against cost", () => {
		const t = completionTime(report, light) as {
			xAxis: { type: string };
			series: { name: string; data: { value: [string | number, number] }[] }[];
		};
		expect(t.xAxis.type).toBe('time');
		expect(t.series[1].data.map((d) => d.value)).toEqual([
			['2026-08-03T12:00:00Z', 1],
			['2026-08-12T09:30:00Z', 2]
		]);
		const c = completionCost(report, light) as typeof t;
		expect(c.xAxis.type).toBe('value');
		expect(c.series[1].data.map((d) => d.value)).toEqual([
			[1.2, 1],
			[1.45, 2]
		]);
	});
	it('draws every chart from a report whose lists are null (older flai, empty selection)', () => {
		const empty = {
			...report,
			items: null,
			throughput: null,
			cfd: null,
			aging: null,
			burnup: null,
			usage: undefined
		} as unknown as Report;
		for (const kind of KINDS) expect(() => build(kind, empty, light)).not.toThrow();
		expect(normalise(empty).aging).toEqual([]);
	});
	it('dark theme swaps the palette and surface', () => {
		expect(dark.series[0]).toBe(CATEGORICAL.dark[0]);
		expect(dark.surface).not.toBe(light.surface);
		expect(human(93600)).toBe('1.1d');
		expect(human(2700)).toBe('45m');
	});
});
