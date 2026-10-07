import { describe, expect, it } from 'vitest';
import {
	build,
	burnUp,
	cfd,
	completedIn,
	bucketLabel,
	bucketsFor,
	controls,
	cost,
	costPerItem,
	costPerModel,
	minutesPerItem,
	timePerModel,
	costSpent,
	cycleTime,
	doneIn,
	errorFacets,
	forecastAccuracy,
	human,
	humanSigned,
	KINDS,
	percentile,
	PLANNING_KINDS,
	spreadOf,
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
	spendRows,
	withUsage,
	type Bucket,
	type ErrorFilter,
	type ItemMetrics,
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
		minutes_per_item: 45,
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
				minutes_per_item: 60,
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
				minutes_per_item: 45,
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
			// the fixture carries no forecasts: the planning charts have their own
			if (!(PLANNING_KINDS as readonly string[]).includes(k))
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
	it('throughput shapes its data', () => {
		const th = throughput(report, light) as { series: { name: string; data: number[] }[] };
		expect(th.series.map((s) => s.name)).toEqual(['feature', 'improvement']);
		expect(th.series[0].data).toEqual([1, 0]);
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
		expect(titleOf('tokens-spent', 'week')).toBe('Tokens / Week');
		expect(titleOf('tokens-spent')).toBe('Tokens / Day');
		expect(titleOf('cost-spent', 'hour')).toBe('$ / Hour');
		expect(titleOf('cost-spent')).toBe('$ / Day');
		expect(titleOf('token-rate', 'hour')).toBe('Tokens / Min');
	});
	it('tokens and cost per item draw one line per item type', () => {
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
		const c = costPerItem(report, light) as Over;
		expect(c.yAxis.name).toBe('US dollars per item');
		expect(values(c.series[0])).toEqual([
			['2026-08-03T00:00:00Z', 0.875],
			['2026-08-05T00:00:00Z', 0.5]
		]);
		expect((build('cost-per-item', report, light) as Over).series.map((s) => s.name)).toEqual([
			'story',
			'task'
		]);
	});
	it('avg. cost per model is the mean dollars per item of the type of the report, one line per model', () => {
		const c = costPerModel(report, light) as Over;
		// with two models, a dashed line for all of them: the type's own cost per item
		expect(c.series.map((s) => s.name)).toEqual([
			'claude-haiku-4-5',
			'claude-opus-5-5',
			'all models'
		]);
		expect(values(c.series[1])).toEqual([
			['2026-08-03T00:00:00Z', 0.75],
			['2026-08-05T00:00:00Z', 0.5]
		]);
		expect(values(c.series[2])).toEqual([
			['2026-08-03T00:00:00Z', 0.875],
			['2026-08-05T00:00:00Z', 0.5]
		]);
		expect(c.yAxis.name).toBe('US dollars per story');
		expect(c.yAxis.axisLabel.formatter(0.75)).toBe('$0.750');
		expect(
			c.tooltip.formatter([{ seriesName: 'claude-haiku-4-5', data: c.series[0].data[0] }])
		).toBe('2026-08-03<br/>claude-haiku-4-5: $0.250 each over 1 item');
		expect((build('cost-per-model', report, light) as Over).series).toEqual(c.series);
		expect(titleOf('cost-per-model')).toBe('Avg. Cost / Model');
		expect(controls('cost-per-model')).toEqual({
			type: true,
			epic: false,
			bucket: true,
			nature: false,
			model: false
		});
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
		expect(controls('tokens-per-item')).toEqual({
			type: false,
			epic: false,
			bucket: true,
			nature: false,
			model: false
		});
		expect(controls('token-rate')).toEqual({
			type: true,
			epic: false,
			bucket: true,
			nature: false,
			model: false
		});
		expect(controls('cost')).toEqual({
			type: true,
			epic: true,
			bucket: false,
			nature: false,
			model: false
		});
		expect(controls('cfd').epic).toBe(false);
		expect(controls('cycle-time').epic).toBe(true);
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
	it('avg. time per model is the mean agent minutes per item of the type of the report, one line per model', () => {
		const o = timePerModel(report, light) as Over;
		expect(o.series.map((s) => s.name)).toEqual([
			'claude-haiku-4-5',
			'claude-opus-5-5',
			'all models'
		]);
		// 5 August took no agent time: no point
		expect(values(o.series[0])).toEqual([['2026-08-03T00:00:00Z', 60]]);
		expect(values(o.series[1])).toEqual([['2026-08-03T00:00:00Z', 45]]);
		expect(values(o.series[2])).toEqual([['2026-08-03T00:00:00Z', 45]]);
		expect(o.yAxis.name).toBe('agent minutes per story');
		expect(o.yAxis.axisLabel.formatter(90)).toBe('1.5h');
		expect(
			o.tooltip.formatter([{ seriesName: 'claude-haiku-4-5', data: o.series[0].data[0] }])
		).toBe('2026-08-03<br/>claude-haiku-4-5: 1h each over 1 item');
		expect((build('time-per-model', report, light) as Over).series).toEqual(o.series);
		expect(titleOf('time-per-model')).toBe('Avg. Time / Model');
		expect(controls('time-per-model')).toEqual({
			type: true,
			epic: false,
			bucket: true,
			nature: false,
			model: false
		});
		// a flai older than S-0169 sends the seconds only
		expect(minutesPerItem({ items: 2, tokens: 0, cost: 0, seconds: 5400 })).toBe(45);
		expect(minutesPerItem({ items: 1, tokens: 0, cost: 0, seconds: 0 })).toBeUndefined();
	});
	it('draws every chart from a report whose lists are null (older flai, empty selection)', () => {
		const empty = {
			...report,
			items: null,
			throughput: null,
			cfd: null,
			burnup: null,
			usage: undefined
		} as unknown as Report;
		for (const kind of KINDS) expect(() => build(kind, empty, light)).not.toThrow();
		expect(normalise(empty).cfd).toEqual([]);
	});
	// What the planner spent (S-0225, ADR-0083): on S-002 beside its agents, and on S-003 alone,
	// done on 4 August, a day whose bucket has no agents' items; S-005 carries neither and stays out.
	const planned: Report = {
		...report,
		items: [
			...report.items.map((i) =>
				i.id === 'S-002'
					? {
							...i,
							usage: {
								...i.usage!,
								strategic: [
									{ kind: 'planner', tokens: 100000, cost: 0.1, seconds: 300, estimated: true },
									{ kind: 'orchestrator', tokens: 1000, cost: 0.2, seconds: 10, estimated: true }
								]
							}
						}
					: i
			),
			...(['S-003', 'S-005'] as const).map((id) => ({
				...report.items[0],
				id,
				title: id === 'S-003' ? 'Three' : 'Five',
				completed: '2026-08-04T10:00:00Z',
				usage: {
					source: 'sum',
					tokens: 0,
					cost: 0,
					seconds: 0,
					models: [],
					strategic:
						id === 'S-003'
							? [{ kind: 'planner', tokens: 200000, cost: 0.4, seconds: 600, estimated: true }]
							: []
				}
			}))
		],
		usage: {
			...report.usage!,
			spend: {
				...report.usage!.spend!,
				story: {
					...report.usage!.spend!.story,
					strategic: { items: 1, tokens: 200000, cost: 0.4, seconds: 600 },
					buckets: storyDays.map((b) => ({
						...b,
						strategic:
							b.at === '2026-08-04T00:00:00Z'
								? { items: 1, tokens: 200000, cost: 0.4, seconds: 600 }
								: { items: 0, tokens: 0, cost: 0, seconds: 0 }
					}))
				}
			}
		}
	};
	it('$ / Item stacks strategic spend apart from the models, estimated, and only when there is some', () => {
		type Bars = {
			legend: { show: boolean };
			xAxis: { data: string[] };
			series: { name: string; stack: string; data: number[]; itemStyle: { color: string } }[];
			tooltip: {
				formatter: (p: { dataIndex: number; seriesName: string; value: number }[]) => string;
			};
		};
		const c = cost(planned, light) as Bars;
		// S-003, on which only the planner spent, has a bar; S-005, with neither, has none
		expect(c.xAxis.data).toEqual(['S-001', 'S-003*', 'S-002*']);
		expect(c.series.map((s) => [s.name, s.data])).toEqual([
			['claude-haiku-4-5', [0.3, 0, 0]],
			['claude-opus-5-5', [1.2, 0, 0.25]],
			['strategic', [0, 0.4, 0.3]]
		]);
		expect(new Set(c.series.map((s) => s.stack)).size).toBe(1);
		expect(c.series[2].itemStyle.color).toBe(CATEGORICAL.light[7]);
		expect((cost(planned, dark) as Bars).series[2].itemStyle.color).toBe(CATEGORICAL.dark[7]);
		expect(c.series.slice(0, 2).map((s) => s.itemStyle.color)).not.toContain(CATEGORICAL.light[7]);
		const at = (dataIndex: number) =>
			c.tooltip.formatter(
				c.series.map((s) => ({ dataIndex, seriesName: s.name, value: s.data[dataIndex] }))
			);
		// the total is the agents'; the strategic spend is apart from it, and estimated, its kinds
		// summed to four decimals: 0.1 and 0.2 make 0.3
		expect(at(2)).toBe(
			'S-002 Two<br/>claude-opus-5-5: $0.250<br/>total $0.250 (estimated in part)<br/>strategic: $0.300 (estimated)'
		);
		expect(at(1)).toBe('S-003 Three<br/>no agent spend<br/>strategic: $0.400 (estimated)');
		expect(at(0)).toBe(
			'S-001 One<br/>claude-haiku-4-5: $0.300<br/>claude-opus-5-5: $1.20<br/>total $1.50'
		);
		// one model and strategic spend: a legend for the two
		const one = cost({ ...planned, items: planned.items.slice(1, 2) }, light) as Bars;
		expect(one.series.map((s) => s.name)).toEqual(['claude-opus-5-5', 'strategic']);
		expect(one.legend.show).toBe(true);
		// without strategic spend, no series for it
		expect((cost(report, light) as Bars).series.map((s) => s.name)).not.toContain('strategic');
	});
	it("$ / bucket stacks strategic spend as a series of its own; the mean stays the agents'", () => {
		const c = costSpent(planned, light) as Over;
		expect(c.series.map((s) => [s.name, s.type])).toEqual([
			['claude-haiku-4-5', 'bar'],
			['claude-opus-5-5', 'bar'],
			['strategic', 'bar'],
			['mean per day', 'line']
		]);
		expect(c.series[2].stack).toBe('spent');
		expect(c.series[2].itemStyle.color).toBe(CATEGORICAL.light[7]);
		expect(values(c.series[2]).map((v) => v[1])).toEqual([0, 0.4, 0]);
		expect(values(c.series[3])).toEqual(values((costSpent(report, light) as Over).series[2]));
		expect(c.tooltip.formatter([{ seriesName: 'strategic', data: c.series[2].data[1] }])).toBe(
			'2026-08-04<br/>strategic: $0.400 (estimated)'
		);
		// without strategic spend, or of tokens, no series for it
		expect((costSpent(report, light) as Over).series.map((s) => s.name)).not.toContain('strategic');
		expect((tokensSpent(planned, light) as Over).series).toEqual(
			(tokensSpent(report, light) as Over).series
		);
		// the table of cost per bucket lists it, estimated, even for a bucket the agents spent nothing in
		const rows = spendRows(planned, 'cost-spent');
		expect(rows.filter((r) => r.of === 'strategic')).toEqual([
			{
				items: 1,
				tokens: 200000,
				cost: 0.4,
				seconds: 600,
				estimated: true,
				at: '2026-08-04T00:00:00Z',
				of: 'strategic',
				mean_tokens: undefined,
				mean_cost: undefined
			}
		]);
		expect(rows.some((r) => 'strategic' in r && r.of !== 'strategic')).toBe(false);
		expect(spendRows(planned, 'tokens-spent')).toEqual(spendRows(report, 'tokens-spent'));
	});
	it('strategic spend leaves the per-model and per-type charts as they were', () => {
		for (const kind of [
			'token-rate',
			'tokens-per-dollar',
			'time-per-model',
			'cost-per-model',
			'cost-per-item',
			'tokens-per-item'
		] as const) {
			expect((build(kind, planned, light) as Over).series, kind).toEqual(
				(build(kind, report, light) as Over).series
			);
		}
		for (const kind of [
			'token-rate',
			'tokens-per-dollar',
			'time-per-model',
			'cost-per-model'
		] as const)
			expect(spendRows(planned, kind), kind).toEqual(spendRows(report, kind));
	});
	it('dark theme swaps the palette and surface', () => {
		expect(dark.series[0]).toBe(CATEGORICAL.dark[0]);
		expect(dark.surface).not.toBe(light.surface);
		expect(human(93600)).toBe('1.1d');
		expect(human(2700)).toBe('45m');
	});
});

// Stories done in the report's window (2 August 12:00 to 1 September 12:00) with forecast errors,
// as flai sends them (S-0205, ADR-0111); S-104 comes from a flai that sends no model. Left out:
// S-106 done before the window, S-107 done with no forecast, S-108 cancelled, S-109 in progress.
const story = (
	id: string,
	nature: string,
	completed: string | undefined,
	extra: Partial<ItemMetrics> = {}
): ItemMetrics => ({
	id,
	type: 'story',
	nature,
	title: `Story ${id.slice(2)}`,
	status: 'done',
	created: '2026-07-01T09:00:00Z',
	completed,
	blocked_seconds: 0,
	time_in_state_seconds: {},
	model: 'claude-opus-5-5',
	...extra
});
const opus = 'claude-opus-5-5';
const haiku = 'claude-haiku-4-5';
const forecasting: Report = {
	...report,
	items: [
		story('S-101', 'feature', '2026-08-03T12:00:00Z', {
			forecast_seconds: 7200,
			forecast_error_seconds: 7200,
			estimate_error_seconds: -3600
		}),
		story('S-102', 'improvement', '2026-08-12T09:30:00Z', {
			model: haiku,
			forecast_error_seconds: -1800
		}),
		story('S-103', 'feature', '2026-08-20T00:00:00Z', { forecast_error_seconds: -10800 }),
		story('S-104', 'feature', '2026-08-25T06:00:00Z', {
			model: undefined,
			forecast_error_seconds: 3600,
			estimate_error_seconds: 900
		}),
		story('S-105', 'improvement', '2026-08-28T00:00:00Z', { forecast_error_seconds: 600 }),
		story('S-106', 'feature', '2026-07-20T00:00:00Z', { forecast_error_seconds: 99999 }),
		story('S-107', 'feature', '2026-08-15T00:00:00Z'),
		story('S-108', 'feature', '2026-08-16T00:00:00Z', {
			status: 'cancelled',
			forecast_error_seconds: 5000
		}),
		story('S-109', 'feature', undefined, { status: 'in-progress' })
	],
	forecasts: {
		forecast: {
			count: 5,
			p50_seconds: 3600,
			p85_seconds: 10800,
			by_nature: {
				feature: { count: 3, p50_seconds: 7200, p85_seconds: 10800 },
				improvement: { count: 2, p50_seconds: 600, p85_seconds: 1800 }
			},
			by_model: {
				[opus]: { count: 3, p50_seconds: 7200, p85_seconds: 10800 },
				[haiku]: { count: 1, p50_seconds: 1800, p85_seconds: 1800 },
				'(none)': { count: 1, p50_seconds: 3600, p85_seconds: 3600 }
			}
		},
		delivery: { count: 0, by_nature: {}, by_model: {} },
		estimate: {
			count: 2,
			p50_seconds: 900,
			p85_seconds: 3600,
			by_nature: { feature: { count: 2, p50_seconds: 900, p85_seconds: 3600 } },
			by_model: {
				[opus]: { count: 1, p50_seconds: 3600, p85_seconds: 3600 },
				'(none)': { count: 1, p50_seconds: 900, p85_seconds: 900 }
			}
		}
	}
};

describe('planning charts', () => {
	type Scatter = {
		name: string;
		type: string;
		data: { value: [string, number]; id: string; title: string }[];
		markLine?: { data: { name: string; yAxis: number }[] };
	};
	type Accuracy = {
		legend: { show: boolean };
		xAxis: { type: string; min?: number; max?: number };
		yAxis: { name: string; axisLabel: { formatter: (v: number) => string } };
		tooltip: { formatter: (p: { seriesName: string; data: Scatter['data'][number] }) => string };
		series: Scatter[];
	};
	const points = (s: Scatter) => s.data.map((d) => [d.id, d.value[0], d.value[1]]);
	const lines = (o: Accuracy) => o.series[0].markLine?.data.map((d) => [d.name, d.yAxis]);
	const accuracy = (f: ErrorFilter) =>
		build('forecast-accuracy', forecasting, light, undefined, f) as Accuracy;
	const band = (p50: number, p85: number) => [
		['p50', p50],
		['p50', -p50],
		['p85', p85],
		['p85', -p85]
	];

	it('forecast accuracy plots each story done in the window by its forecast and estimate errors', () => {
		const o = forecastAccuracy(forecasting, light) as Accuracy;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['forecast error', 'scatter'],
			['estimate error', 'scatter']
		]);
		// x completed, y actual minus forecast in seconds; the stories without an error, done before
		// the window, cancelled, or not done are left out
		expect(points(o.series[0])).toEqual([
			['S-101', '2026-08-03T12:00:00Z', 7200],
			['S-102', '2026-08-12T09:30:00Z', -1800],
			['S-103', '2026-08-20T00:00:00Z', -10800],
			['S-104', '2026-08-25T06:00:00Z', 3600],
			['S-105', '2026-08-28T00:00:00Z', 600]
		]);
		expect(points(o.series[1])).toEqual([
			['S-101', '2026-08-03T12:00:00Z', -3600],
			['S-104', '2026-08-25T06:00:00Z', 900]
		]);
		// the p50 and p85 of the absolute forecast error, flai's own, either side of zero
		expect(lines(o)).toEqual(band(3600, 10800));
		expect(o.series[1].markLine).toBeUndefined();
		expect(o.legend.show).toBe(true);
		expect(o.series[0].data[1].title).toBe('Story 102');
		expect(o.tooltip.formatter({ seriesName: 'forecast error', data: o.series[0].data[1] })).toBe(
			'S-102 Story 102<br/>forecast error: -30m · 2026-08-12'
		);
		expect(o.yAxis.name).toBe('actual minus forecast');
		expect(o.yAxis.axisLabel.formatter(7200)).toBe('+2h');
		expect(humanSigned(-10800)).toBe('-3h');
		expect(humanSigned(0)).toBe('+0m');
		expect((build('forecast-accuracy', forecasting, light) as Accuracy).series).toEqual(o.series);
	});
	it("forecast accuracy spans the report's window (ADR-0054)", () => {
		const o = forecastAccuracy(forecasting, light) as Accuracy;
		expect(o.xAxis).toMatchObject({
			type: 'time',
			min: Date.parse('2026-08-02T12:00:00Z'),
			max: Date.parse('2026-09-01T12:00:00Z')
		});
		// a narrower window moves the axis and drops S-101, done on 3 August
		const narrow = { ...forecasting, window_start: '2026-08-10T00:00:00Z' };
		const n = forecastAccuracy(narrow, light) as Accuracy;
		expect(n.xAxis.min).toBe(Date.parse('2026-08-10T00:00:00Z'));
		expect(n.series[0].data.map((d) => d.id)).toEqual(['S-102', 'S-103', 'S-104', 'S-105']);
		expect(n.series[1].data.map((d) => d.id)).toEqual(['S-104']);
		expect(doneIn(narrow).map((i) => i.id)).toEqual(['S-102', 'S-103', 'S-104', 'S-105', 'S-107']);
	});
	it("narrows forecast accuracy by nature and by model, with flai's percentiles for each", () => {
		const feature = accuracy({ nature: 'feature' });
		expect(feature.series[0].data.map((d) => d.id)).toEqual(['S-101', 'S-103', 'S-104']);
		expect(feature.series[1].data.map((d) => d.id)).toEqual(['S-101', 'S-104']);
		expect(lines(feature)).toEqual(band(7200, 10800));
		// a story from a flai that sends no model is (none)
		const none = accuracy({ model: '(none)' });
		expect(none.series[0].data.map((d) => d.id)).toEqual(['S-104']);
		expect(lines(none)).toEqual(band(3600, 3600));
		const fast = accuracy({ model: haiku });
		expect(fast.series.map((s) => s.name)).toEqual(['forecast error']);
		expect(lines(fast)).toEqual(band(1800, 1800));
		expect(fast.legend.show).toBe(false);
		// the lines are flai's figures, not worked out again
		const nudged: Report = {
			...forecasting,
			forecasts: {
				...forecasting.forecasts!,
				forecast: {
					...forecasting.forecasts!.forecast,
					by_nature: { feature: { count: 3, p50_seconds: 7201, p85_seconds: 10801 } }
				}
			}
		};
		const read = forecastAccuracy(nudged, light, { nature: 'feature' }) as Accuracy;
		expect(lines(read)).toEqual(band(7201, 10801));
	});
	it('works out the percentiles from the stories shown under both filters, as flai does', () => {
		// improvement and opus: S-105 alone, 600; neither the nature's (600, 1800) nor the model's
		const both = accuracy({ nature: 'improvement', model: opus });
		expect(both.series.length).toBe(1);
		expect(points(both.series[0])).toEqual([['S-105', '2026-08-28T00:00:00Z', 600]]);
		expect(lines(both)).toEqual(band(600, 600));
		// feature and opus: S-101 and S-103, 7200 and 10800 absolute
		const featureOpus = accuracy({ nature: 'feature', model: opus });
		expect(featureOpus.series[0].data.map((d) => d.id)).toEqual(['S-101', 'S-103']);
		expect(lines(featureOpus)).toEqual(band(7200, 10800));
		// nothing shown: no points and no lines
		expect(accuracy({ nature: 'research', model: opus }).series).toEqual([]);
		// nearest rank, at least 1, of the absolute errors: flai's figures for the window
		expect(spreadOf([7200, -1800, -10800, 3600, 600])).toEqual({
			count: 5,
			p50_seconds: 3600,
			p85_seconds: 10800
		});
		expect(spreadOf([-3600, 900])).toEqual({ count: 2, p50_seconds: 900, p85_seconds: 3600 });
		expect(spreadOf([])).toEqual({ count: 0 });
		expect(percentile([5], 1)).toBe(5);
		expect(percentile([1, 2, 3, 4], 0)).toBe(1);
		expect(percentile([1, 2, 3, 4], 100)).toBe(4);
	});
	it('lists the natures and models of the stories with an error each planning chart plots', () => {
		expect(errorFacets(forecasting, 'forecast-accuracy')).toEqual({
			natures: ['feature', 'improvement'],
			models: ['(none)', haiku, opus]
		});
		// no story carries a delivery error
		expect(errorFacets(forecasting, 'delivery-accuracy')).toEqual({ natures: [], models: [] });
		const nulls = { ...forecasting, items: null } as unknown as Report;
		expect(errorFacets(nulls, 'forecast-by-model')).toEqual({ natures: [], models: [] });
	});
	it('draws an empty planning chart from a report without forecasts or errors', () => {
		// the main fixture, as an older flai sends it
		expect(report.forecasts).toBeUndefined();
		for (const kind of PLANNING_KINDS) {
			const o = build(kind, report, light) as Accuracy;
			expect(o.series, kind).toEqual([]);
			expect(o.xAxis, kind).toMatchObject({
				type: 'time',
				min: Date.parse('2026-08-02T12:00:00Z'),
				max: Date.parse('2026-09-01T12:00:00Z')
			});
		}
		// errors without the spreads: points, but no lines
		const older: Report = { ...forecasting, forecasts: undefined };
		const o = build('forecast-accuracy', older, light) as Accuracy;
		expect(o.series[0].data.length).toBe(5);
		expect(o.series[0].markLine).toBeUndefined();
	});
	it('offers the planning charts nature and model filters on the stories', () => {
		expect(controls('forecast-accuracy')).toEqual({
			type: false,
			epic: false,
			bucket: false,
			nature: true,
			model: true
		});
		expect(controls('delivery-accuracy')).toMatchObject({ bucket: false, model: true });
		expect(controls('forecast-by-model')).toMatchObject({
			bucket: true,
			nature: true,
			model: false
		});
		expect(PLANNING_KINDS.map((k) => titleOf(k))).toEqual([
			'Forecast Accuracy',
			'Delivery Accuracy',
			'Forecast Error / Model'
		]);
		expect(KINDS).toEqual(expect.arrayContaining([...PLANNING_KINDS]));
	});
});
