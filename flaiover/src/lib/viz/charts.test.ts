import { describe, expect, it } from 'vitest';
import {
	agentWaiting,
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
	deliveryAccuracy,
	doneIn,
	errorFacets,
	FLOW_KINDS,
	forecastAccuracy,
	forecastByModel,
	forecastRows,
	driftedIn,
	FORECAST_KINDS,
	hasClaims,
	holdTime,
	human,
	humanSigned,
	isClaimsKind,
	isForecastKind,
	KINDS,
	parallelism,
	touchesDrift,
	CLAIMS_KINDS,
	percentile,
	PLANNING_KINDS,
	spreadOf,
	normalise,
	onTimeShare,
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
	strategicRatio,
	strategicRows,
	STRATEGIC_KINDS,
	withUsage,
	COD_COLUMNS,
	COD_KINDS,
	codDayRows,
	codIncurred,
	codOrder,
	codOrderRows,
	codOutstanding,
	codWeekRows,
	hasCostOfDelay,
	hasOrder,
	orderSummary,
	withoutValueNow,
	type Bucket,
	type Claims,
	type CodDay,
	type CostOfDelayReport,
	type ErrorFilter,
	type ItemMetrics,
	type Report,
	type StrategicDay,
	type WaitWeek
} from './charts';
import {
	CATEGORICAL,
	modelSlot,
	modelSymbol,
	ORDER_SLOT,
	STATE_SLOT,
	theme,
	TYPE_SLOT
} from './palette';

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

/** A week of waiting as flai lays it out (S-0205): what its items waited on threads and in review. */
function waitWeek(week: string, start: string, items = 0, threads = 0, review = 0): WaitWeek {
	const zero = { count: 0, total_seconds: 0 };
	const mean = (s: number) => (items > 0 ? { mean_seconds: s / items } : {});
	return {
		week,
		start,
		items,
		empty_wakes: 0,
		threads: { total_seconds: threads, ...mean(threads), orchestrator: zero, confirmed: zero },
		review: { total_seconds: review, ...mean(review) }
	};
}

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
	// one week per ISO week of the window, 27 July to 31 August: S-001 done in the week of
	// 3 August, S-002 in the week of 10 August, and the rest with no items
	waiting: {
		weeks: [
			waitWeek('2026-W31', '2026-07-27'),
			waitWeek('2026-W32', '2026-08-03', 1, 7200, 7200),
			waitWeek('2026-W33', '2026-08-10', 1, 3600, 50400),
			waitWeek('2026-W34', '2026-08-17'),
			waitWeek('2026-W35', '2026-08-24'),
			waitWeek('2026-W36', '2026-08-31')
		],
		empty_wakes: { count: 0 },
		longest: [
			{
				item: 'S-002',
				kind: 'review',
				started: '2026-08-11T19:30:00Z',
				ended: '2026-08-12T09:30:00Z',
				seconds: 50400,
				awaited: 'operator'
			}
		]
	},
	// the days S-001 and S-002 were completed; flai sends every day of the window
	strategic_days: [
		{
			date: '2026-08-03',
			agents: { planner: { cost: 0.4, seconds: 600, estimated: true } },
			cost: 0.4,
			seconds: 600,
			completed: 1,
			cost_per_item: 1.5,
			cycle_time_seconds: 93600
		},
		{
			date: '2026-08-12',
			agents: {},
			cost: 0,
			seconds: 0,
			completed: 1,
			cost_per_item: 0.25,
			cycle_time_seconds: 86400
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

// Cost of delay as flai stats --json prints it (S-0205, S-0213), over a window from Sunday
// 23 August 12:00 to 1 September 12:00: a day from the one that holds the start, and the ISO
// weeks from the one that holds it, Monday 17 August.
const codDay = (
	date: string,
	[backlog, ready, inProgress, review]: number[],
	[b, r, i, v]: number[],
	incurred: number
): CodDay => ({
	date,
	outstanding: { backlog, ready, 'in-progress': inProgress, review },
	without_value: { backlog: b, ready: r, 'in-progress': i, review: v },
	incurred
});
const cod: CostOfDelayReport = {
	days: [
		...['23', '24', '25', '26', '27', '28'].map((d) =>
			codDay(`2026-08-${d}`, [100, 50, 0, 0], [2, 1, 0, 0], 21.43)
		),
		codDay('2026-08-29', [100, 50, 25, 0], [2, 1, 0, 0], 21.43),
		codDay('2026-08-30', [80, 70, 25, 0], [2, 1, 0, 0], 21.43),
		codDay('2026-08-31', [80, 45, 25, 25], [3, 0, 1, 0], 17.86),
		codDay('2026-09-01', [80, 45, 0, 50], [3, 0, 0, 1], 8.93)
	],
	weeks: [
		{ week: '2026-W34', start: '2026-08-17', incurred: 150 },
		{ week: '2026-W35', start: '2026-08-24', incurred: 140.5 },
		{ week: '2026-W36', start: '2026-08-31', incurred: 30.25 }
	],
	// S-010 is worth 100 a week, S-011 50, S-012 10; S-013 has no forecast and is left out
	order: {
		at: '2026-09-01T12:00:00Z',
		horizon: '2026-09-08T12:00:00Z',
		series: [
			{
				by: 'current',
				total: 107.14,
				points: [
					{ at: '2026-09-01T12:00:00Z', incurred: 0 },
					{ at: '2026-09-01T12:00:00Z', id: 'S-012', incurred: 0 },
					{ at: '2026-09-05T12:00:00Z', id: 'S-010', incurred: 57.14 },
					{ at: '2026-09-08T12:00:00Z', id: 'S-011', incurred: 107.14 }
				]
			},
			{
				by: 'cod',
				total: 18.58,
				points: [
					{ at: '2026-09-01T12:00:00Z', incurred: 0 },
					{ at: '2026-09-01T12:00:00Z', id: 'S-010', incurred: 0 },
					{ at: '2026-09-03T12:00:00Z', id: 'S-011', incurred: 14.29 },
					{ at: '2026-09-04T12:00:00Z', id: 'S-012', incurred: 18.58 }
				]
			},
			{
				by: 'wsjf',
				total: 17.15,
				points: [
					{ at: '2026-09-01T12:00:00Z', incurred: 0 },
					{ at: '2026-09-01T12:00:00Z', id: 'S-011', incurred: 0 },
					{ at: '2026-09-02T12:00:00Z', id: 'S-010', incurred: 14.29 },
					{ at: '2026-09-03T12:00:00Z', id: 'S-012', incurred: 17.15 }
				]
			}
		],
		saving: 89.99,
		cheaper: 'wsjf',
		left_out: ['S-013']
	}
};
const codReport: Report = {
	...report,
	window_days: 9,
	window_start: '2026-08-23T12:00:00Z',
	cost_of_delay: cod
};

const light = theme(false);
const dark = theme(true);
/**
 * The local midnight of a date, where a chart draws flai's UTC day or week of that date (S-0329):
 * in New York, where the tests run, four or five hours after the UTC day starts.
 */
const midnight = (date: string) => {
	const [y, m, d] = date.slice(0, 10).split('-').map(Number);
	return new Date(y, m - 1, d).getTime();
};

describe('chart builders', () => {
	it('every kind builds with one y-axis and a tooltip', () => {
		for (const k of KINDS) {
			// the fixture carries no forecasts, claims, or cost of delay: those charts draw from their own
			const cod = (COD_KINDS as readonly string[]).includes(k);
			const fixture = isClaimsKind(k)
				? claiming
				: isForecastKind(k)
					? forecasting
					: cod
						? codReport
						: report;
			const o = build(k, fixture, light) as {
				useUTC?: boolean;
				yAxis: unknown;
				tooltip: unknown;
				series: unknown[];
			};
			// every axis and tooltip reads the local zone (S-0329)
			expect(o.useUTC, k).toBeUndefined();
			expect(o.yAxis, k).toBeDefined();
			// delivery accuracy and touches drift alone draw a share on a second axis
			if (k === 'delivery-accuracy' || k === 'touches-drift')
				expect((o.yAxis as unknown[]).length).toBe(2);
			else expect(Array.isArray(o.yAxis), `${k} must not use two y-axes`).toBe(false);
			expect(o.tooltip, k).toBeDefined();
			// whose stories carry no delivery errors: delivery accuracy has its own fixture below
			if (k !== 'delivery-accuracy') expect(o.series.length, k).toBeGreaterThan(0);
		}
	});
	it("every chart spans the report's window and plots only the items completed in it (S-0166)", () => {
		type Axis = { xAxis: { min?: number; max?: number; data?: string[] } };
		const end = Date.parse('2026-09-01T12:00:00Z');
		const hour = 3600e3;
		// the report's own window, 2 August 12:00 to 1 September 12:00
		const month = Date.parse('2026-08-02T12:00:00Z');
		expect((cycleTime(report, light) as Axis).xAxis).toMatchObject({ min: month, max: end });
		// a series by the day starts on the local midnight of the day that holds the window's start,
		// where its first date is drawn
		const day = midnight('2026-08-02');
		expect((burnUp(report, light) as Axis).xAxis).toMatchObject({ min: day, max: end });
		expect((cfd(report, light) as Axis).xAxis).toMatchObject({ min: day, max: end });
		// now just after midnight in UTC, the evening before here: the axis runs to the local midnight
		// of now's UTC day, where that day's point is drawn
		const late: Report = { ...report, generated_at: '2026-09-01T02:00:00Z' };
		expect((burnUp(late, light) as Axis).xAxis).toMatchObject({
			min: day,
			max: midnight('2026-09-01')
		});
		// time in state by the day (S-0168): from the day that holds the start to the one that holds
		// now, each on its date's tick, half a day either side
		expect((timeInState(report, light) as Axis).xAxis).toMatchObject({
			type: 'time',
			min: day - 12 * hour,
			max: midnight('2026-09-01') + 12 * hour
		});
		// spend over time: from the bucket that holds the start to the one that holds now, half a
		// bucket either side
		expect((tokensSpent(report, light) as Axis).xAxis).toMatchObject({
			min: day - 12 * hour,
			max: midnight('2026-09-01') + 12 * hour
		});
		const weekly = { ...report, usage: { ...report.usage!, bucket: 'week' as const } };
		expect((costSpent(weekly, light) as Axis).xAxis).toMatchObject({
			min: midnight('2026-07-27') - 84 * hour,
			max: midnight('2026-08-31') + 84 * hour
		});

		// a narrower window moves the axis and drops S-001, completed on 3 August
		const narrow: Report = { ...report, window_days: 7, window_start: '2026-08-10T00:00:00Z' };
		const from = Date.parse('2026-08-10T00:00:00Z');
		const ct = cycleTime(narrow, light) as Axis & { series: { data: { id: string }[] }[] };
		expect(ct.xAxis).toMatchObject({ min: from, max: end });
		expect(ct.series.flatMap((s) => s.data.map((d) => d.id))).toEqual(['S-002']);
		const tis = timeInState(narrow, light) as Axis & { series: { data: { ids: string[] }[] }[] };
		expect(tis.xAxis.min).toBe(midnight('2026-08-10') - 12 * hour);
		expect(tis.series[0].data.map((d) => d.ids)).toEqual([['S-002']]);
		expect((cost(narrow, light) as Axis).xAxis.data).toEqual(['S-002*']);
		expect((tokenRate(narrow, light) as Axis).xAxis.min).toBe(midnight('2026-08-10') - 12 * hour);
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
		type Bars = {
			series: {
				name: string;
				data: { value: [number, number]; ids: string[]; day: string }[];
			}[];
			tooltip: { formatter: (ps: unknown[]) => string };
		};
		const t = timeInState(report, light) as Bars;
		// a bar per UTC day with items completed, on its date's tick: the local midnight of its date
		expect(t.series[2].name).toBe('in-progress');
		expect(t.series[2].data).toEqual([
			{ value: [midnight('2026-08-03'), 24], ids: ['S-001'], day: '2026-08-03' },
			{ value: [midnight('2026-08-12'), 10], ids: ['S-002'], day: '2026-08-12' }
		]);
		// the tooltip names the bar's own date, not the local one of the moment the UTC day starts
		expect(
			t.tooltip.formatter(t.series.map((s) => ({ seriesName: s.name, data: s.data[0] })))
		).toMatch(/^2026-08-03, mean of 1 story<br\/>S-001<br\/>/);
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
	type Waits = {
		xAxis: { type: string; min: number; max: number };
		yAxis: { name: string };
		legend: { show: boolean };
		tooltip: { formatter: (ps: unknown[]) => string };
		series: {
			name: string;
			type: string;
			stack?: string;
			connectNulls?: boolean;
			itemStyle: { color: string };
			data: { value: [number, number | null]; week: WaitWeek }[];
		}[];
	};
	// a week from its Monday in UTC, drawn on its Monday's tick
	const monday = midnight;
	it('agent waiting stacks the hours waited on threads and in review per week, with the mean per story as a line', () => {
		const o = agentWaiting(report, light) as Waits;
		expect(o.series.map((s) => [s.name, s.type, s.stack])).toEqual([
			['threads', 'bar', 'waiting'],
			['review', 'bar', 'waiting'],
			['mean per story', 'line', undefined]
		]);
		expect(o.series[0].itemStyle.color).toBe(CATEGORICAL.light[2]);
		expect(o.series[1].itemStyle.color).toBe(CATEGORICAL.light[7]);
		// a bar per week, at its Monday, of the hours its items waited
		const starts = [
			'2026-07-27',
			'2026-08-03',
			'2026-08-10',
			'2026-08-17',
			'2026-08-24',
			'2026-08-31'
		];
		expect(o.series[0].data.map((d) => d.value)).toEqual(
			starts.map((s, i) => [monday(s), [0, 2, 1, 0, 0, 0][i]])
		);
		expect(o.series[1].data.map((d) => d.value)).toEqual(
			starts.map((s, i) => [monday(s), [0, 2, 14, 0, 0, 0][i]])
		);
		// the mean wait per story, threads and review together; a gap in a week with no items
		expect(o.series[2].data.map((d) => d.value)).toEqual(
			starts.map((s, i) => [monday(s), [null, 4, 15, null, null, null][i]])
		);
		expect(o.series[2].connectNulls).toBe(false);
		expect(o.yAxis.name).toBe('hours waited');
		expect(o.legend.show).toBe(true);
		expect(
			o.tooltip.formatter(o.series.map((s) => ({ seriesName: s.name, data: s.data[2] })))
		).toBe(
			'week of 2026-08-10 (2026-W33), 1 story done<br/>threads: 1h<br/>review: 14h<br/>mean per story: 15h'
		);
		// a week with no items has no mean to show
		expect(
			o.tooltip.formatter(o.series.slice(2).map((s) => ({ seriesName: s.name, data: s.data[0] })))
		).toBe('');
		expect((build('agent-waiting', report, light) as Waits).series).toEqual(o.series);
	});
	it('agent waiting divides the waits of a week by the items done in it', () => {
		const two: Report = {
			...report,
			waiting: { ...report.waiting!, weeks: [waitWeek('2026-W32', '2026-08-03', 2, 3600, 10800)] }
		};
		const o = agentWaiting(two, light) as Waits;
		expect(o.series.map((s) => s.data[0].value[1])).toEqual([1, 3, 2]);
		expect(
			o.tooltip.formatter(o.series.map((s) => ({ seriesName: s.name, data: s.data[0] })))
		).toBe(
			'week of 2026-08-03 (2026-W32), 2 stories done<br/>threads: 1h<br/>review: 3h<br/>mean per story: 2h'
		);
	});
	it("agent waiting spans the report's window in weeks (ADR-0054)", () => {
		const hour = 3600e3;
		// from the week that holds 2 August, the window's start, to the one that holds now, half a week
		// either side
		expect((agentWaiting(report, light) as Waits).xAxis).toMatchObject({
			type: 'time',
			min: monday('2026-07-27') - 84 * hour,
			max: monday('2026-08-31') + 84 * hour
		});
		// whatever flai sends: a report without waiting spans the window still
		const older = agentWaiting({ ...report, waiting: undefined }, light) as Waits;
		expect(older.xAxis).toMatchObject({ min: monday('2026-07-27') - 84 * hour });
	});
	it('draws an empty agent waiting from a report without waiting, or with null lists', () => {
		const older = build('agent-waiting', { ...report, waiting: undefined }, light) as Waits;
		expect(older.series).toEqual([]);
		const nulls = {
			...report,
			waiting: { weeks: null, empty_wakes: { count: 0 } }
		} as unknown as Report;
		expect((build('agent-waiting', nulls, light) as Waits).series).toEqual([]);
		expect(normalise(nulls).waiting).toMatchObject({ weeks: [], longest: [] });
		expect(normalise({ ...report, waiting: undefined }).waiting).toBeUndefined();
	});
	it('lists agent waiting under the flow charts, summed by flai over every epic', () => {
		expect(FLOW_KINDS).toContain('agent-waiting');
		expect(KINDS).toContain('agent-waiting');
		expect(titleOf('agent-waiting')).toBe('Agent Waiting');
		expect(controls('agent-waiting')).toEqual({
			type: true,
			epic: false,
			bucket: false,
			nature: false,
			model: false
		});
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
		data: { value: [number, number]; at: string; items: number; estimated?: boolean }[];
	};
	type Over = {
		series: Line[];
		legend: { show: boolean };
		useUTC?: boolean;
		xAxis: { type: string; min?: number; max?: number; minInterval: number };
		yAxis: { name: string; axisLabel: { formatter: (v: number) => string } };
		tooltip: { trigger: string; formatter: (p: unknown) => string };
	};
	// each point by its bucket's start as flai sends it, and its value
	const values = (l: Line) => l.data.map((d) => [d.at, d.value[1]]);
	it('token rate is tokens per agent minute over time, per model, with all of them dashed', () => {
		const o = tokenRate(report, light) as Over;
		expect(o.series.map((s) => s.name)).toEqual([
			'claude-haiku-4-5',
			'claude-opus-5-5',
			'all models'
		]);
		expect(o.xAxis.type).toBe('time');
		// the axis reads the local zone, each UTC day on its date's tick, and it spans the window's
		// days, not the one day drawn (S-0166)
		expect(o.useUTC).toBeUndefined();
		expect([o.xAxis.min, o.xAxis.max, o.xAxis.minInterval]).toEqual([
			midnight('2026-08-02') - 43200e3,
			midnight('2026-09-01') + 43200e3,
			86400e3
		]);
		expect(o.yAxis.name).toBe('tokens per agent minute');
		expect(o.yAxis.axisLabel.formatter(38888.9)).toBe('38.9K');
		// the day with nothing done has no point, and neither has the one with no agent time
		expect(values(o.series[1])).toEqual([['2026-08-03T00:00:00Z', 27777.8]]);
		expect(values(o.series[2])).toEqual([['2026-08-03T00:00:00Z', 38888.9]]);
		// a day's point sits on its date's tick: the local midnight of 3 August, not 20:00 on the 2nd
		expect(o.series[1].data[0].value[0]).toBe(midnight('2026-08-03'));
		expect(o.series[1].data[0].value[0]).toBe(Date.parse('2026-08-03T04:00:00Z'));
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
		// an hour in the local zone; a UTC day or week by its own date, which the local one is not
		expect(bucketLabel('2026-09-29T19:00:00Z', 'hour')).toBe('2026-09-29 15:00 EDT');
		expect(bucketLabel('2026-09-30T02:00:00Z', 'hour')).toBe('2026-09-29 22:00 EDT');
		expect(bucketLabel('2026-09-29T19:00:00Z', 'hour')).not.toContain('UTC');
		expect(bucketLabel('2026-09-29T00:00:00Z', 'day')).toBe('2026-09-29');
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
	type Cod = {
		useUTC?: boolean;
		legend: { show: boolean };
		xAxis: { type: string; min?: number; max?: number; minInterval?: number };
		yAxis: { name: string };
		tooltip: { formatter: (p: unknown) => string };
		series: {
			name: string;
			type: string;
			stack?: string;
			step?: string;
			symbol?: string;
			lineStyle: { type?: string; color: string };
			itemStyle: { color: string };
			areaStyle?: { color: string };
			data: unknown[];
		}[];
	};
	it('CoD Outstanding stacks the value per week of each column per day, in board order', () => {
		const o = codOutstanding(codReport, light) as Cod;
		expect(o.series.map((s) => s.name)).toEqual(['backlog', 'ready', 'in-progress', 'review']);
		expect(new Set(o.series.map((s) => s.stack))).toEqual(new Set(['outstanding']));
		// a band per column in the colour the cumulative flow gives its state, with the surface as seam
		for (const s of o.series) {
			expect(s.areaStyle?.color).toBe(CATEGORICAL.light[STATE_SLOT[s.name]]);
			expect(s.lineStyle.color).toBe(light.surface);
		}
		expect(o.series[0].data).toHaveLength(10);
		expect(o.series.map((s) => s.data.slice(-3))).toEqual([
			[
				['2026-08-30', 80],
				['2026-08-31', 80],
				['2026-09-01', 80]
			],
			[
				['2026-08-30', 70],
				['2026-08-31', 45],
				['2026-09-01', 45]
			],
			[
				['2026-08-30', 25],
				['2026-08-31', 25],
				['2026-09-01', 0]
			],
			[
				['2026-08-30', 0],
				['2026-08-31', 25],
				['2026-09-01', 50]
			]
		]);
		// by the day over the window: from the day that holds its start, on its date's tick, to its
		// now (ADR-0054)
		expect([o.xAxis.type, o.xAxis.min, o.xAxis.max]).toEqual([
			'time',
			midnight('2026-08-23'),
			Date.parse('2026-09-01T12:00:00Z')
		]);
		expect(o.yAxis.name).toBe('value per week');
		expect(o.legend.show).toBe(true);
		expect((build('cod-outstanding', codReport, light) as Cod).series).toEqual(o.series);
		// the items without a value now: the last day's, per column
		expect(withoutValueNow(codReport)).toEqual({
			backlog: 3,
			ready: 0,
			'in-progress': 0,
			review: 1
		});
		expect(Object.keys(withoutValueNow(codReport)!)).toEqual([...COD_COLUMNS]);
		const rows = codDayRows(codReport);
		expect(rows.map((d) => d.date).slice(0, 2)).toEqual(['2026-09-01', '2026-08-31']);
		expect(rows[0]).toEqual(cod.days[cod.days.length - 1]);
	});
	it('CoD Incurred is a bar per week with the running mean per week from the first, dashed', () => {
		const o = codIncurred(codReport, light) as Cod;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['incurred', 'bar'],
			['mean per week', 'line']
		]);
		type P = { value: [number, number]; at: string };
		const points = (i: number) => (o.series[i].data as P[]).map((d) => d.value);
		// each week on its Monday's tick
		expect(points(0)).toEqual([
			[midnight('2026-08-17'), 150],
			[midnight('2026-08-24'), 140.5],
			[midnight('2026-08-31'), 30.25]
		]);
		expect((o.series[0].data as P[]).map((d) => d.at)).toEqual([
			'2026-08-17T00:00:00Z',
			'2026-08-24T00:00:00Z',
			'2026-08-31T00:00:00Z'
		]);
		// the sum from the window's first week to this one over the number of those weeks:
		// 150, 290.5 / 2, 320.75 / 3, to two decimals
		expect(points(1).map((v) => v[1])).toEqual([150, 145.25, 106.92]);
		expect(o.series[1].lineStyle.type).toBe('dashed');
		expect(o.series[0].itemStyle.color).toBe(CATEGORICAL.light[0]);
		// UTC weeks over the window on an axis in the local zone: from the week that holds the start
		// to the one that holds now, half a week either side
		const halfWeek = 3.5 * 86400e3;
		expect(o.useUTC).toBeUndefined();
		expect([o.xAxis.min, o.xAxis.max, o.xAxis.minInterval]).toEqual([
			midnight('2026-08-17') - halfWeek,
			midnight('2026-08-31') + halfWeek,
			7 * 86400e3
		]);
		expect(o.legend.show).toBe(true);
		expect(
			o.tooltip.formatter([
				{ seriesName: 'incurred', data: o.series[0].data[1] },
				{ seriesName: 'mean per week', data: o.series[1].data[1] }
			])
		).toBe('week of 2026-08-24<br/>incurred: 140.5<br/>mean per week: 145.25');
		expect(codWeekRows(codReport)).toEqual([
			{ week: '2026-W36', start: '2026-08-31', incurred: 30.25, mean: 106.92 },
			{ week: '2026-W35', start: '2026-08-24', incurred: 140.5, mean: 145.25 },
			{ week: '2026-W34', start: '2026-08-17', incurred: 150, mean: 150 }
		]);
	});
	it('CoD by Order draws a step line per order from now to the last pull, and states the saving', () => {
		const o = codOrder(codReport, light) as Cod;
		expect(o.series.map((s) => s.name)).toEqual(['pull order', 'by cost of delay', 'by WSJF']);
		expect(o.series.map((s) => s.step)).toEqual(['end', 'end', 'end']);
		type P = { value: [number, number]; at: string; id?: string };
		const points = (i: number) => (o.series[i].data as P[]).map((d) => [d.id, d.at, d.value[1]]);
		expect(points(0)).toEqual([
			[undefined, '2026-09-01T12:00:00Z', 0],
			['S-012', '2026-09-01T12:00:00Z', 0],
			['S-010', '2026-09-05T12:00:00Z', 57.14],
			['S-011', '2026-09-08T12:00:00Z', 107.14]
		]);
		// each pull is drawn at its moment
		for (const d of o.series.flatMap((s) => s.data as P[]))
			expect(d.value[0]).toBe(Date.parse(d.at));
		expect(points(1).map((p) => p[2])).toEqual([0, 0, 14.29, 18.58]);
		expect(points(2).map((p) => p[0])).toEqual([undefined, 'S-011', 'S-010', 'S-012']);
		// a colour and a mark per order, apart for every pair
		expect(o.series.map((s) => s.itemStyle.color)).toEqual([
			CATEGORICAL.light[ORDER_SLOT.current],
			CATEGORICAL.light[ORDER_SLOT.cod],
			CATEGORICAL.light[ORDER_SLOT.wsjf]
		]);
		expect(new Set(o.series.map((s) => s.symbol)).size).toBe(3);
		// a projection: from its start to its horizon, not over the window (ADR-0112)
		expect([o.useUTC, o.xAxis.type, o.xAxis.min, o.xAxis.max]).toEqual([
			undefined,
			'time',
			Date.parse('2026-09-01T12:00:00Z'),
			Date.parse('2026-09-08T12:00:00Z')
		]);
		expect(o.yAxis.name).toBe('projected cost');
		expect(o.legend.show).toBe(true);
		expect(o.tooltip.formatter({ seriesName: 'pull order', data: o.series[0].data[2] })).toBe(
			'pull order<br/>S-010 pulled 2026-09-05 08:00 EDT<br/>cumulative 57.14'
		);
		expect(o.tooltip.formatter({ seriesName: 'by WSJF', data: o.series[2].data[0] })).toBe(
			'by WSJF<br/>now, 2026-09-01 08:00 EDT<br/>cumulative 0'
		);
		expect(orderSummary(codReport)).toEqual({
			totals: { current: 107.14, cod: 18.58, wsjf: 17.15 },
			saving: 89.99,
			cheaper: 'wsjf',
			left_out: ['S-013']
		});
		// the table: each projected pull of each order, in order of pull
		const rows = codOrderRows(codReport);
		expect(rows).toHaveLength(9);
		expect(rows.slice(0, 4)).toEqual([
			{ by: 'current', at: '2026-09-01T12:00:00Z', id: 'S-012', incurred: 0 },
			{ by: 'current', at: '2026-09-05T12:00:00Z', id: 'S-010', incurred: 57.14 },
			{ by: 'current', at: '2026-09-08T12:00:00Z', id: 'S-011', incurred: 107.14 },
			{ by: 'cod', at: '2026-09-01T12:00:00Z', id: 'S-010', incurred: 0 }
		]);
		expect(rows.map((r) => r.by)).toEqual([
			...Array(3).fill('current'),
			...Array(3).fill('cod'),
			...Array(3).fill('wsjf')
		]);
	});
	it('the cost of delay charts take the window, and the type but on the projection, and are named as the others are', () => {
		expect(controls('cod-outstanding')).toEqual({
			type: true,
			epic: false,
			bucket: false,
			nature: false,
			model: false
		});
		expect(controls('cod-incurred')).toEqual({
			type: true,
			epic: false,
			bucket: false,
			nature: false,
			model: false
		});
		expect(controls('cod-order')).toEqual({
			type: false,
			epic: false,
			bucket: false,
			nature: false,
			model: false
		});
		expect(COD_KINDS.map((k) => titleOf(k))).toEqual([
			'CoD Outstanding',
			'CoD Incurred',
			'CoD by Order'
		]);
		expect(KINDS.slice(-3)).toEqual([...COD_KINDS]);
	});
	it('the cost of delay charts are empty, not broken, from a flai that sends less', () => {
		// older than S-0205: no cost of delay at all
		expect(hasCostOfDelay(report)).toBe(false);
		expect(hasOrder(report)).toBe(false);
		for (const kind of COD_KINDS) {
			const o = build(kind, report, light) as Cod;
			expect(o.series, kind).toEqual([]);
			expect(o.legend.show, kind).toBe(false);
		}
		expect(withoutValueNow(report)).toBeUndefined();
		expect(orderSummary(report)).toBeUndefined();
		expect([codDayRows(report), codWeekRows(report), codOrderRows(report)]).toEqual([[], [], []]);
		// older than S-0213: no items without a value, and no order
		const older: Report = {
			...codReport,
			cost_of_delay: {
				days: cod.days.map((d) => ({
					date: d.date,
					outstanding: d.outstanding,
					incurred: d.incurred
				})),
				weeks: cod.weeks
			}
		};
		expect(hasCostOfDelay(older)).toBe(true);
		expect(hasOrder(older)).toBe(false);
		expect((codOutstanding(older, light) as Cod).series).toEqual(
			(codOutstanding(codReport, light) as Cod).series
		);
		expect((codIncurred(older, light) as Cod).series).toHaveLength(2);
		expect(withoutValueNow(older)).toBeUndefined();
		expect((codOrder(older, light) as Cod).series).toEqual([]);
		expect(orderSummary(older)).toBeUndefined();
		expect(codOrderRows(older)).toEqual([]);
		// no story placed: each order is its first point, and nothing is saved
		const none: Report = {
			...codReport,
			cost_of_delay: {
				...cod,
				order: {
					at: '2026-09-01T12:00:00Z',
					horizon: '2026-09-01T12:00:00Z',
					series: (['current', 'cod', 'wsjf'] as const).map((by) => ({
						by,
						total: 0,
						points: [{ at: '2026-09-01T12:00:00Z', incurred: 0 }]
					})),
					saving: 0,
					cheaper: 'cod',
					left_out: ['S-010', 'S-011']
				}
			}
		};
		expect((codOrder(none, light) as Cod).series.map((s) => s.data.length)).toEqual([1, 1, 1]);
		expect(codOrderRows(none)).toEqual([]);
		expect(orderSummary(none)).toEqual({
			totals: { current: 0, cod: 0, wsjf: 0 },
			saving: 0,
			cheaper: 'cod',
			left_out: ['S-010', 'S-011']
		});
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
		// the date of its completion in the local zone: midnight in UTC is the evening before here
		expect(o.tooltip.formatter({ seriesName: 'forecast error', data: o.series[0].data[2] })).toBe(
			'S-103 Story 103<br/>forecast error: -3h · 2026-08-19'
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
		// the chart by the bucket runs from the day that holds the start to the one that holds now,
		// half a day either side, each day on its date's tick
		const start = Date.parse('2026-08-02T12:00:00Z');
		const end = Date.parse('2026-09-01T12:00:00Z');
		const spans = {
			'forecast-accuracy': [start, end, 0],
			'delivery-accuracy': [start, end, 0],
			'forecast-by-model': [midnight('2026-08-02'), midnight('2026-09-01'), 43200e3]
		} as const;
		for (const kind of FORECAST_KINDS) {
			const o = build(kind, report, light) as Accuracy;
			const [from, to, half] = spans[kind];
			expect(o.series, kind).toEqual([]);
			expect(o.xAxis, kind).toMatchObject({ type: 'time', min: from - half, max: to + half });
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
		// the Planning group lists the forecast charts, then the claims charts (S-0214)
		expect(PLANNING_KINDS.map((k) => titleOf(k))).toEqual([
			'Forecast Accuracy',
			'Delivery Accuracy',
			'Forecast Error / Model',
			'Parallelism',
			'Hold Time',
			'Touches Drift'
		]);
		expect(PLANNING_KINDS).toEqual([...FORECAST_KINDS, ...CLAIMS_KINDS]);
		expect(KINDS).toEqual(expect.arrayContaining([...PLANNING_KINDS]));
	});

	// Delivery errors on the stories of `forecasting`: S-101 a day early, S-102 half a day late,
	// S-103 on the day, S-104 two days late, S-105 an hour early; S-106 is done before the window,
	// S-108 is cancelled, and S-107 carries no delivery forecast.
	const delivered: Record<string, number> = {
		'S-101': -86400,
		'S-102': 43200,
		'S-103': 0,
		'S-104': 172800,
		'S-105': -3600,
		'S-106': -60,
		'S-108': 600
	};
	const delivering: Report = {
		...forecasting,
		items: forecasting.items.map((i) =>
			i.id in delivered ? { ...i, delivery_error_seconds: delivered[i.id] } : i
		),
		forecasts: {
			...forecasting.forecasts!,
			delivery: {
				count: 5,
				p50_seconds: 43200,
				p85_seconds: 172800,
				by_nature: {
					feature: { count: 3, p50_seconds: 86400, p85_seconds: 172800 },
					improvement: { count: 2, p50_seconds: 3600, p85_seconds: 43200 }
				},
				by_model: {
					[opus]: { count: 3, p50_seconds: 3600, p85_seconds: 86400 },
					[haiku]: { count: 1, p50_seconds: 43200, p85_seconds: 43200 },
					'(none)': { count: 1, p50_seconds: 172800, p85_seconds: 172800 }
				}
			}
		}
	};
	type Week = {
		value: [number, number | null];
		week: string;
		start: string;
		count: number;
		on_time: number;
		share?: number;
	};
	type Delivery = Omit<Accuracy, 'yAxis' | 'tooltip'> & {
		yAxis: {
			name: string;
			min?: number;
			max?: number;
			axisLabel: { formatter: (v: number) => string };
		}[];
		tooltip: {
			formatter: (p: { seriesName: string; data: Scatter['data'][number] | Week }) => string;
		};
	};
	const delivery = (f: ErrorFilter, r: Report = delivering) =>
		build('delivery-accuracy', r, light, undefined, f) as Delivery;
	const refs = (o: Delivery) => o.series[0].markLine?.data.map((d) => [d.name, d.yAxis]);
	const weeksOf = (o: Delivery) =>
		(o.series.find((s) => s.name === 'on time per week')?.data ?? []) as unknown as Week[];
	const ids = (o: Delivery) => o.series[0].data.map((d) => d.id);
	const at = (s: string) => Date.parse(s);
	const day = 86400;

	it('delivery accuracy plots each story done in the window by its delivery error in days', () => {
		const o = deliveryAccuracy(delivering, light) as Delivery;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['delivery error', 'scatter'],
			['on time per week', 'line']
		]);
		// x completed, y completed minus the forecast date in days, unrounded so that flai's seconds
		// stand; done before the window, without a delivery forecast, cancelled, or not done: left out
		expect(points(o.series[0])).toEqual([
			['S-101', '2026-08-03T12:00:00Z', -1],
			['S-102', '2026-08-12T09:30:00Z', 0.5],
			['S-103', '2026-08-20T00:00:00Z', 0],
			['S-104', '2026-08-25T06:00:00Z', 2],
			['S-105', '2026-08-28T00:00:00Z', -3600 / day]
		]);
		// flai's p50 and p85 of the absolute delivery error, in days, either side of zero
		expect(refs(o)).toEqual(band(0.5, 2));
		expect(o.series[1]).toMatchObject({ yAxisIndex: 1, connectNulls: false });
		expect(o.yAxis.map((y) => y.name)).toEqual(['days after forecast date', 'on time']);
		expect(o.yAxis[0].axisLabel.formatter(-1)).toBe('-1d');
		expect(o.yAxis[1]).toMatchObject({ min: 0, max: 1 });
		expect(o.yAxis[1].axisLabel.formatter(0.5)).toBe('50%');
		expect(o.legend.show).toBe(true);
		expect(o.tooltip.formatter({ seriesName: 'delivery error', data: o.series[0].data[1] })).toBe(
			'S-102 Story 102<br/>delivery error: +12h · 2026-08-12'
		);
		expect(delivery({}).series).toEqual(o.series);
	});
	it('delivery accuracy gives the share on time per ISO week of the window, none in an empty week', () => {
		// from the week that holds the window's start, 2 August, to the one that holds now, 1
		// September; S-107, done in W33 without a delivery forecast, counts in neither points nor share
		expect(onTimeShare(delivering)).toEqual([
			{ week: '2026-W31', start: '2026-07-27', count: 0, on_time: 0 },
			{ week: '2026-W32', start: '2026-08-03', count: 1, on_time: 1, share: 1 },
			{ week: '2026-W33', start: '2026-08-10', count: 1, on_time: 0, share: 0 },
			{ week: '2026-W34', start: '2026-08-17', count: 1, on_time: 1, share: 1 },
			{ week: '2026-W35', start: '2026-08-24', count: 2, on_time: 1, share: 0.5 },
			{ week: '2026-W36', start: '2026-08-31', count: 0, on_time: 0 }
		]);
		// each at the middle of the part of its week in the window; an empty week is a gap, not 0
		const o = deliveryAccuracy(delivering, light) as Delivery;
		expect(weeksOf(o).map((w) => w.value)).toEqual([
			[at('2026-08-02T18:00:00Z'), null],
			[at('2026-08-06T12:00:00Z'), 1],
			[at('2026-08-13T12:00:00Z'), 0],
			[at('2026-08-20T12:00:00Z'), 1],
			[at('2026-08-27T12:00:00Z'), 0.5],
			[at('2026-08-31T18:00:00Z'), null]
		]);
		const weeks = weeksOf(o);
		expect(o.tooltip.formatter({ seriesName: 'on time per week', data: weeks[4] })).toBe(
			'week of 2026-08-24 (2026-W35)<br/>1 of 2 on or before the forecast date (50%)'
		);
		expect(o.tooltip.formatter({ seriesName: 'on time per week', data: weeks[0] })).toBe(
			'week of 2026-07-27 (2026-W31)<br/>no story with a delivery forecast'
		);
		// ISO weeks as flai names them across a new year: 30 December 2024 is in 2025-W01
		const turn = {
			...delivering,
			window_start: '2024-12-25T00:00:00Z',
			generated_at: '2025-01-08T00:00:00Z'
		};
		expect(onTimeShare(turn).map((w) => [w.week, w.start])).toEqual([
			['2024-W52', '2024-12-23'],
			['2025-W01', '2024-12-30'],
			['2025-W02', '2025-01-06']
		]);
	});
	it("delivery accuracy spans the report's window (ADR-0054)", () => {
		const o = deliveryAccuracy(delivering, light) as Delivery;
		expect(o.xAxis).toMatchObject({
			type: 'time',
			min: at('2026-08-02T12:00:00Z'),
			max: at('2026-09-01T12:00:00Z')
		});
		// a narrower window moves the axis, drops S-101, done on 3 August, and starts the weeks at W33
		const narrow = { ...delivering, window_start: '2026-08-10T00:00:00Z' };
		const n = deliveryAccuracy(narrow, light) as Delivery;
		expect(n.xAxis.min).toBe(at('2026-08-10T00:00:00Z'));
		expect(ids(n)).toEqual(['S-102', 'S-103', 'S-104', 'S-105']);
		expect(onTimeShare(narrow).map((w) => w.week)).toEqual([
			'2026-W33',
			'2026-W34',
			'2026-W35',
			'2026-W36'
		]);
		expect(weeksOf(n)[0].value).toEqual([at('2026-08-13T12:00:00Z'), 0]);
	});
	it("narrows delivery accuracy by nature and by model, with flai's percentiles for each", () => {
		const feature = delivery({ nature: 'feature' });
		expect(ids(feature)).toEqual(['S-101', 'S-103', 'S-104']);
		expect(refs(feature)).toEqual(band(1, 2));
		expect(weeksOf(feature).map((w) => w.value[1])).toEqual([null, 1, null, 1, 0, null]);
		const none = delivery({ model: '(none)' });
		expect(ids(none)).toEqual(['S-104']);
		expect(refs(none)).toEqual(band(2, 2));
		expect(onTimeShare(delivering, { model: '(none)' }).map((w) => w.share)).toEqual([
			undefined,
			undefined,
			undefined,
			undefined,
			0,
			undefined
		]);
		// under both, worked out from the stories shown: improvement and opus is S-105 alone
		const both = delivery({ nature: 'improvement', model: opus });
		expect(ids(both)).toEqual(['S-105']);
		expect(refs(both)).toEqual(band(3600 / day, 3600 / day));
		// the lines are flai's figures, not worked out again
		const nudged: Report = {
			...delivering,
			forecasts: {
				...delivering.forecasts!,
				delivery: {
					...delivering.forecasts!.delivery,
					by_nature: { feature: { count: 3, p50_seconds: 86401, p85_seconds: 172801 } }
				}
			}
		};
		expect(refs(delivery({ nature: 'feature' }, nudged))).toEqual(band(86401 / day, 172801 / day));
		// nothing shown: no points, no lines, and no share
		expect(delivery({ nature: 'research' }).series).toEqual([]);
		expect(errorFacets(delivering, 'delivery-accuracy')).toEqual({
			natures: ['feature', 'improvement'],
			models: ['(none)', haiku, opus]
		});
	});
	it('draws an empty delivery accuracy from a report without delivery errors', () => {
		// forecasts but no delivery errors: two axes over the window, nothing on them
		const o = deliveryAccuracy(forecasting, light) as Delivery;
		expect(o.series).toEqual([]);
		expect(o.yAxis.length).toBe(2);
		const empty = onTimeShare(forecasting);
		expect(empty.length).toBe(6);
		expect(empty.every((w) => w.count === 0 && w.share === undefined)).toBe(true);
		// errors without the spreads, from an older flai: points and shares, but no lines
		const older = delivery({}, { ...delivering, forecasts: undefined });
		expect(older.series[0].data.length).toBe(5);
		expect(older.series[0].markLine).toBeUndefined();
		expect(weeksOf(older).length).toBe(6);
		// null lists, and a report without a window, which has no weeks
		const nulls = { ...delivering, items: null } as unknown as Report;
		expect(delivery({}, nulls).series).toEqual([]);
		expect(onTimeShare(nulls).map((w) => w.count)).toEqual([0, 0, 0, 0, 0, 0]);
		expect(onTimeShare({ ...delivering, window_start: '' })).toEqual([]);
	});

	// The stories of `forecasting` and two more by opus in the week of 17 August, laid out by the
	// week: opus has three stories that week, S-103 (-3h) and S-110 (20m) on 20 August, and S-111
	// (-4000s) on 21 August.
	type ModelPoint = { value: [number, number | null]; at: string; count: number };
	type ByModel = {
		useUTC?: boolean;
		legend: { show: boolean };
		xAxis: { type: string; min?: number; max?: number; minInterval?: number };
		yAxis: { name: string; axisLabel: { formatter: (v: number) => string } };
		tooltip: {
			trigger: string;
			formatter: (ps: { marker?: string; seriesName: string; data: ModelPoint }[]) => string;
		};
		series: {
			name: string;
			type: string;
			symbol: string;
			showSymbol: boolean;
			itemStyle: { color: string };
			lineStyle: { color: string };
			data: ModelPoint[];
		}[];
	};
	const inBuckets = (bucket: 'hour' | 'day' | 'week', r: Report = forecasting): Report => ({
		...r,
		usage: { ...r.usage!, bucket },
		items: [
			...r.items,
			story('S-110', 'feature', '2026-08-20T12:00:00Z', { forecast_error_seconds: 1200 }),
			story('S-111', 'improvement', '2026-08-21T00:00:00Z', { forecast_error_seconds: -4000 })
		]
	});
	const weekly = inBuckets('week');
	const byModel = (r: Report, f: ErrorFilter = {}) =>
		build('forecast-by-model', r, light, undefined, f) as ByModel;
	const seriesOf = (o: ByModel, model: string) => o.series.find((s) => s.name === model)!;
	/** A model's points, the gaps left out: the bucket's day, the p50, and the stories under it. */
	const drawn = (o: ByModel, model: string) => {
		const points = seriesOf(o, model).data.filter((d) => d.value[1] !== null);
		return points.map((d) => [d.at.slice(0, 10), d.value[1], d.count]);
	};

	it('forecast by model plots the p50 absolute forecast error per bucket per model', () => {
		const o = forecastByModel(weekly, light) as ByModel;
		// a line per model in order of name; S-104, from a flai that sends no model, is (none)
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['(none)', 'line'],
			[haiku, 'line'],
			[opus, 'line']
		]);
		// nearest rank of the absolute errors, as flai: 1200, 4000, 10800 has 4000 for its p50
		expect(drawn(o, opus)).toEqual([
			['2026-08-03', 7200, 1],
			['2026-08-17', 4000, 3],
			['2026-08-24', 600, 1]
		]);
		expect(drawn(o, haiku)).toEqual([['2026-08-10', 1800, 1]]);
		expect(drawn(o, '(none)')).toEqual([['2026-08-24', 3600, 1]]);
		// every week of the window has a place on each line; a week without the model's stories is a
		// gap, not 0
		expect(seriesOf(o, opus).data.map((d) => d.value[1])).toEqual([
			null,
			7200,
			null,
			4000,
			600,
			null
		]);
		// a week on its Monday's tick
		expect(seriesOf(o, opus).data[0]).toEqual({
			value: [midnight('2026-07-27'), null],
			at: '2026-07-27T00:00:00.000Z',
			count: 0
		});
		// each model in its fixed colour and mark, every mark shown
		for (const m of [opus, haiku, '(none)']) {
			const s = seriesOf(o, m);
			expect(s.itemStyle.color, m).toBe(light.series[modelSlot(m)]);
			expect(s.lineStyle.color, m).toBe(light.series[modelSlot(m)]);
			expect(s.symbol, m).toBe(modelSymbol(m));
			expect(s.showSymbol, m).toBe(true);
		}
		expect(o.legend.show).toBe(true);
		expect(o.yAxis.name).toBe('p50 absolute forecast error');
		expect(o.yAxis.axisLabel.formatter(7200)).toBe('2h');
		// the tooltip names the bucket and each model with a point in it
		const week = 3;
		const tip = o.tooltip.formatter([
			{ marker: '', seriesName: opus, data: seriesOf(o, opus).data[week] },
			{ marker: '', seriesName: haiku, data: seriesOf(o, haiku).data[week] }
		]);
		expect(tip).toBe(`week of 2026-08-17<br/>${opus}: p50 1.1h over 3 stories`);
		expect(o.tooltip.formatter([{ seriesName: opus, data: seriesOf(o, opus).data[0] }])).toBe('');
		expect(o.tooltip.trigger).toBe('axis');
		expect(byModel(weekly).series).toEqual(o.series);
	});
	it('forecast by model takes its bucket from the report, as the spend charts do', () => {
		// by the day, S-103 and S-110 share 20 August: 1200 and 10800 has 1200 for its p50
		const o = byModel(inBuckets('day'));
		expect(drawn(o, opus)).toEqual([
			['2026-08-03', 7200, 1],
			['2026-08-20', 1200, 2],
			['2026-08-21', 4000, 1],
			['2026-08-28', 600, 1]
		]);
		// a day per place, 2 August to 1 September
		expect(seriesOf(o, opus).data.length).toBe(31);
		expect(o.tooltip.formatter([{ seriesName: haiku, data: seriesOf(o, haiku).data[10] }])).toBe(
			`2026-08-12<br/>${haiku}: p50 30m over 1 story`
		);
		// a report from a flai that names no bucket is by the day
		const daily = inBuckets('day');
		const unnamed = { ...daily, usage: { ...daily.usage!, bucket: undefined } };
		expect(byModel(unnamed).series).toEqual(o.series);
	});
	it("forecast by model spans the report's window in its buckets (ADR-0054)", () => {
		const hour = 3600e3;
		// from the bucket that holds the window's start to the one that holds now, half a bucket
		// either side, with ticks no finer than a bucket, in the local zone: a day or a week on its
		// date's tick, an hour at its moment
		const w = byModel(weekly);
		expect(w.useUTC).toBeUndefined();
		expect(w.xAxis).toMatchObject({
			type: 'time',
			minInterval: 168 * hour,
			min: midnight('2026-07-27') - 84 * hour,
			max: midnight('2026-08-31') + 84 * hour
		});
		expect(byModel(inBuckets('day')).xAxis).toMatchObject({
			minInterval: 24 * hour,
			min: midnight('2026-08-02') - 12 * hour,
			max: midnight('2026-09-01') + 12 * hour
		});
		const hourly = byModel(inBuckets('hour'));
		expect(hourly.xAxis).toMatchObject({
			minInterval: hour,
			min: Date.parse('2026-08-02T12:00:00Z') - hour / 2,
			max: Date.parse('2026-09-01T12:00:00Z') + hour / 2
		});
		expect(drawn(hourly, opus).length).toBe(5);
		// an hour at its moment, named in the local zone
		const first = seriesOf(hourly, opus).data.find((d) => d.value[1] !== null)!;
		expect(first.value[0]).toBe(Date.parse('2026-08-03T12:00:00Z'));
		expect(hourly.tooltip.formatter([{ seriesName: opus, data: first }])).toBe(
			`2026-08-03 08:00 EDT<br/>${opus}: p50 2h over 1 story`
		);
		// a narrower window moves the axis and drops S-101, done on 3 August
		const narrow = byModel({ ...weekly, window_start: '2026-08-10T00:00:00Z' });
		expect(narrow.xAxis.min).toBe(midnight('2026-08-10') - 84 * hour);
		expect(seriesOf(narrow, opus).data.map((d) => d.value[1])).toEqual([null, 4000, 600, null]);
		// a report without a window has a place per bucket drawn, a bucket either side; S-106, done
		// before the window, comes back
		const open = byModel({ ...weekly, window_start: '' });
		expect(drawn(open, opus)[0]).toEqual(['2026-07-20', 99999, 1]);
		expect(seriesOf(open, haiku).data.length).toBe(5);
		expect(open.xAxis).toMatchObject({
			min: midnight('2026-07-20') - 168 * hour,
			max: midnight('2026-08-24') + 168 * hour
		});
	});
	it('forecast by model narrows by nature and shows every model whatever the model filter', () => {
		const feature = byModel(weekly, { nature: 'feature' });
		// S-102, haiku's only story, is an improvement
		expect(feature.series.map((s) => s.name)).toEqual(['(none)', opus]);
		expect(drawn(feature, opus)).toEqual([
			['2026-08-03', 7200, 1],
			['2026-08-17', 1200, 2]
		]);
		expect(byModel(weekly, { model: haiku }).series).toEqual(byModel(weekly).series);
		expect(byModel(weekly, { nature: 'feature', model: haiku }).series).toEqual(feature.series);
		expect(byModel(weekly, { nature: 'research' }).series).toEqual([]);
		const nulls = { ...weekly, items: null } as unknown as Report;
		expect(byModel(nulls).series).toEqual([]);
	});
	it("forecast by model over a single bucket agrees with flai's p50 per model", () => {
		// the stories of `forecasting` moved into the week of 24 August, the window that week alone
		const moved: Record<string, string> = {
			'S-101': '2026-08-24T12:00:00Z',
			'S-102': '2026-08-25T09:30:00Z',
			'S-103': '2026-08-26T00:00:00Z',
			'S-104': '2026-08-27T06:00:00Z',
			'S-105': '2026-08-28T00:00:00Z',
			'S-107': '2026-08-26T00:00:00Z'
		};
		const oneWeek: Report = {
			...forecasting,
			window_start: '2026-08-24T00:00:00Z',
			generated_at: '2026-08-30T12:00:00Z',
			usage: { ...forecasting.usage!, bucket: 'week' },
			items: forecasting.items.map((i) => (i.id in moved ? { ...i, completed: moved[i.id] } : i))
		};
		const o = byModel(oneWeek);
		const flai = forecasting.forecasts!.forecast.by_model;
		expect(o.series.map((s) => s.name)).toEqual(Object.keys(flai).sort());
		for (const [m, s] of Object.entries(flai))
			expect(seriesOf(o, m).data, m).toEqual([
				{
					value: [midnight('2026-08-24'), s.p50_seconds],
					at: '2026-08-24T00:00:00.000Z',
					count: s.count
				}
			]);
	});
	it('gives the rows of the planning table: the stories done in the window with an error', () => {
		// newest first; S-107 has no error, S-106 is done before the window, S-108 is cancelled
		expect(forecastRows(forecasting).map((row) => row.id)).toEqual([
			'S-105',
			'S-104',
			'S-103',
			'S-102',
			'S-101'
		]);
		const rows = forecastRows(delivering);
		expect(rows.find((row) => row.id === 'S-101')).toEqual({
			id: 'S-101',
			title: 'Story 101',
			completed: '2026-08-03T12:00:00Z',
			nature: 'feature',
			model: opus,
			forecast_seconds: 7200,
			forecast_error_seconds: 7200,
			delivery_error_seconds: -86400,
			estimate_error_seconds: -3600
		});
		expect(rows.find((row) => row.id === 'S-104')).toMatchObject({
			model: '(none)',
			forecast_error_seconds: 3600,
			delivery_error_seconds: 172800,
			estimate_error_seconds: 900
		});
		// a story with a delivery error alone is a row; its cycle time comes with it
		const late = {
			...forecasting,
			items: forecasting.items.map((i) =>
				i.id === 'S-107' ? { ...i, delivery_error_seconds: 0, cycle_time_seconds: 86400 } : i
			)
		};
		const s107 = forecastRows(late).find((row) => row.id === 'S-107');
		expect(s107).toMatchObject({ delivery_error_seconds: 0, cycle_time_seconds: 86400 });
		expect(forecastRows(late).map((row) => row.id)).toEqual([
			'S-105',
			'S-104',
			'S-103',
			'S-107',
			'S-102',
			'S-101'
		]);
		// under the filter, nature and model
		expect(forecastRows(forecasting, { nature: 'feature' }).map((row) => row.id)).toEqual([
			'S-104',
			'S-103',
			'S-101'
		]);
		expect(forecastRows(forecasting, { nature: 'improvement', model: opus })).toEqual([
			expect.objectContaining({ id: 'S-105' })
		]);
		expect(forecastRows(forecasting, { model: '(none)' }).map((row) => row.id)).toEqual(['S-104']);
		// null lists from an older flai
		expect(forecastRows({ ...forecasting, items: null } as unknown as Report)).toEqual([]);
	});
});

// Claims as flai sends them (S-0205, ADR-0113) over the stories of `forecasting`, in the window
// of 2 August 12:00 to 1 September 12:00: a day per day of the window, a week per ISO week from
// 27 July to 31 August, and the drift of S-101 (two files outside), S-102 (exact), and S-103 (a file
// outside, a touch unchanged), done in the window; of S-106, done before it; of S-108, cancelled;
// and of S-109, in progress.
const windowDays = Array.from({ length: 31 }, (_, n) =>
	new Date(Date.parse('2026-08-02') + n * 86400e3).toISOString().slice(0, 10)
);
const held = (overlap: number, after: number, noTouches: number) => ({
	overlap,
	after,
	'no-touches': noTouches
});
const drift = (id: string, outside: string[], unchanged: string[]) => ({
	id,
	committed: [...outside, 'flai/internal/metrics/claims.go'].sort(),
	outside,
	unchanged,
	outside_count: outside.length,
	unchanged_count: unchanged.length
});
const claims: Claims = {
	limit: 2,
	days: windowDays.map((date, n) => ({ date, in_progress: n % 4, held: n === 10 ? 2 : 0 })),
	weeks: [
		{ week: '2026-W31', start: '2026-07-27', held_seconds: held(0, 0, 0), stories: 0, exact: 0 },
		{
			week: '2026-W32',
			start: '2026-08-03',
			held_seconds: held(7200, 3600, 0),
			stories: 1,
			exact: 0,
			exact_share: 0
		},
		{
			week: '2026-W33',
			start: '2026-08-10',
			held_seconds: held(0, 0, 1800),
			stories: 1,
			exact: 1,
			exact_share: 1
		},
		{
			week: '2026-W34',
			start: '2026-08-17',
			held_seconds: held(900, 0, 0),
			stories: 1,
			exact: 0,
			exact_share: 0
		},
		{ week: '2026-W35', start: '2026-08-24', held_seconds: held(0, 0, 0), stories: 0, exact: 0 },
		{ week: '2026-W36', start: '2026-08-31', held_seconds: held(0, 5400, 0), stories: 0, exact: 0 }
	],
	drift: [
		drift('S-101', ['flai/cmd/a.go', 'flai/cmd/b.go'], []),
		drift('S-102', [], []),
		drift('S-103', ['docs/x.md'], ['flaiover/src/lib/viz']),
		drift('S-106', ['flai/cmd/c.go'], []),
		drift('S-108', ['flai/cmd/d.go'], []),
		drift('S-109', ['flai/cmd/e.go'], ['docs'])
	]
};
const claiming: Report = { ...forecasting, claims };

describe('claims charts', () => {
	type Line = {
		name: string;
		type: string;
		stack?: string;
		yAxisIndex?: number;
		lineStyle: { color: string; type: string };
		itemStyle: { color: string };
		data: unknown[];
	};
	type Option = {
		useUTC?: boolean;
		legend: { show: boolean };
		xAxis: { type: string; min?: number; max?: number; minInterval?: number };
		yAxis: { name: string } | { name: string; min?: number; max?: number }[];
		tooltip: { formatter: (p: unknown) => string };
		series: Line[];
	};
	const day = 86400e3;
	const start = Date.parse('2026-08-02T12:00:00Z');
	const end = Date.parse('2026-09-01T12:00:00Z');
	const empty: Claims = { days: [], weeks: [], drift: [] };

	it('parallelism plots the items in progress and the stories held per day, and the limit', () => {
		const o = parallelism(claiming, light) as Option;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['in progress', 'line'],
			['held', 'line'],
			['limit', 'line']
		]);
		const [inProgress, heldStories, limit] = o.series;
		expect(inProgress.data.length).toBe(31);
		expect(inProgress.data.slice(0, 3)).toEqual([
			['2026-08-02', 0],
			['2026-08-03', 1],
			['2026-08-04', 2]
		]);
		expect(inProgress.data[30]).toEqual(['2026-09-01', 2]);
		expect(heldStories.data[10]).toEqual(['2026-08-12', 2]);
		expect(heldStories.data[11]).toEqual(['2026-08-13', 0]);
		// the limit is a flat dashed line over every day
		expect(new Set(limit.data.map((d) => (d as [string, number])[1]))).toEqual(new Set([2]));
		expect(limit.data.length).toBe(31);
		expect(limit.lineStyle).toMatchObject({ type: 'dashed', color: light.textSecondary });
		// in progress and held in their states' colours
		expect(inProgress.itemStyle.color).toBe(light.series[0]);
		expect(heldStories.itemStyle.color).toBe(light.series[3]);
		expect(o.legend.show).toBe(true);
		expect((o.yAxis as { name: string }).name).toBe('stories');
		// a series by the day: from the day that holds the window's start, on its date's tick, to now
		expect(o.xAxis).toMatchObject({
			type: 'time',
			min: midnight('2026-08-02'),
			max: end
		});
		expect((build('parallelism', claiming, light) as Option).series).toEqual(o.series);
	});
	it('draws no limit line when the board has none, and a gap where a day carries no held count', () => {
		const days = claims.days.map(({ date, in_progress }) => ({ date, in_progress }));
		const unlimited: Report = { ...claiming, claims: { ...claims, limit: undefined, days } };
		const o = parallelism(unlimited, light) as Option;
		expect(o.series.map((s) => s.name)).toEqual(['in progress', 'held']);
		// a flai older than S-0214 sends no held count
		expect(o.series[1].data[0]).toEqual(['2026-08-02', null]);
	});
	it('hold time stacks the hours held per week by reason', () => {
		const o = holdTime(claiming, light) as Option;
		expect(o.series.map((s) => [s.name, s.type, s.stack])).toEqual([
			['overlap', 'bar', 'held'],
			['after', 'bar', 'held'],
			['empty claim', 'bar', 'held']
		]);
		type Bar = { value: [number, number]; seconds: number };
		const hoursOf = (s: Line) => s.data.map((d) => (d as Bar).value[1]);
		expect(hoursOf(o.series[0])).toEqual([0, 2, 0, 0.25, 0, 0]);
		expect(hoursOf(o.series[1])).toEqual([0, 1, 0, 0, 0, 1.5]);
		expect(hoursOf(o.series[2])).toEqual([0, 0, 0.5, 0, 0, 0]);
		// each bar on its week's Monday's tick
		expect(o.series[0].data.map((d) => (d as Bar).value[0])).toEqual(
			['07-27', '08-03', '08-10', '08-17', '08-24', '08-31'].map((d) => midnight(`2026-${d}`))
		);
		// three reasons in three colours apart from each other
		expect(new Set(o.series.map((s) => s.itemStyle.color)).size).toBe(3);
		const week = 1;
		expect(
			o.tooltip.formatter(o.series.map((s) => ({ seriesName: s.name, data: s.data[week] })))
		).toBe(
			'week of 2026-08-03 (2026-W32), 3h held<br/>overlap: 2h<br/>after: 1h<br/>empty claim: 0m'
		);
		expect(o.tooltip.formatter([])).toBe('');
		// from the week that holds the window's start to the one that holds now, half a week either side
		expect(o.xAxis).toMatchObject({
			type: 'time',
			minInterval: 7 * day,
			min: midnight('2026-07-27') - 3.5 * day,
			max: midnight('2026-08-31') + 3.5 * day
		});
		expect(o.useUTC).toBeUndefined();
		expect((o.yAxis as { name: string }).name).toBe('hours held');
		expect((build('hold-time', claiming, light) as Option).series).toEqual(o.series);
	});
	it('picks the stories of the drift done in the window, cancelled ones left out', () => {
		// S-106 is done before the window, S-108 cancelled, S-109 in progress
		expect(driftedIn(claiming).map((s) => [s.id, s.title, s.completed])).toEqual([
			['S-101', 'Story 101', '2026-08-03T12:00:00Z'],
			['S-102', 'Story 102', '2026-08-12T09:30:00Z'],
			['S-103', 'Story 103', '2026-08-20T00:00:00Z']
		]);
		expect(driftedIn(forecasting)).toEqual([]);
	});
	it('touches drift stacks each story by paths outside and unchanged, with the exact share per week', () => {
		const o = touchesDrift(claiming, light) as Option;
		expect(o.series.map((s) => [s.name, s.type, s.yAxisIndex])).toEqual([
			['outside its touches', 'bar', undefined],
			['touches unchanged', 'bar', undefined],
			['exact touches per week', 'line', 1]
		]);
		expect(o.series[0].stack).toBe(o.series[1].stack);
		type Bar = { value: [string, number]; id: string; paths: string[] };
		const bars = (s: Line) => s.data.map((d) => [(d as Bar).id, ...(d as Bar).value]);
		expect(bars(o.series[0])).toEqual([
			['S-101', '2026-08-03T12:00:00Z', 2],
			['S-102', '2026-08-12T09:30:00Z', 0],
			['S-103', '2026-08-20T00:00:00Z', 1]
		]);
		expect(bars(o.series[1])).toEqual([
			['S-101', '2026-08-03T12:00:00Z', 0],
			['S-102', '2026-08-12T09:30:00Z', 0],
			['S-103', '2026-08-20T00:00:00Z', 1]
		]);
		// each week's share at the middle of its part in the window; a week without a story is a gap
		type Share = { value: [number, number | null] };
		expect(o.series[2].data.map((d) => (d as Share).value)).toEqual([
			[Date.parse('2026-08-02T18:00:00Z'), null],
			[Date.parse('2026-08-06T12:00:00Z'), 0],
			[Date.parse('2026-08-13T12:00:00Z'), 1],
			[Date.parse('2026-08-20T12:00:00Z'), 0],
			[Date.parse('2026-08-27T12:00:00Z'), null],
			[Date.parse('2026-08-31T18:00:00Z'), null]
		]);
		const tip = (s: number, d: number) =>
			o.tooltip.formatter({ seriesName: o.series[s].name, data: o.series[s].data[d] });
		expect(tip(0, 0)).toBe(
			'S-101 Story 101 · 2026-08-03<br/>outside its touches: 2<br/>flai/cmd/a.go, flai/cmd/b.go'
		);
		expect(tip(1, 1)).toBe('S-102 Story 102 · 2026-08-12<br/>touches unchanged: 0');
		// the date of its completion in the local zone: midnight in UTC is the evening before here
		expect(tip(0, 2)).toBe('S-103 Story 103 · 2026-08-19<br/>outside its touches: 1<br/>docs/x.md');
		expect(tip(2, 2)).toBe(
			'week of 2026-08-10 (2026-W33)<br/>1 of 1 story with exact touches (100%)'
		);
		const axes = o.yAxis as { name: string; min?: number; max?: number }[];
		expect(axes.map((a) => a.name)).toEqual(['paths', 'exact touches']);
		expect(axes[1]).toMatchObject({ min: 0, max: 1 });
		expect(o.legend.show).toBe(true);
		expect(o.xAxis).toMatchObject({ type: 'time', min: start, max: end });
		expect((build('touches-drift', claiming, light) as Option).series).toEqual(o.series);
	});
	it('draws touches drift empty when git could not be read, and hold time still', () => {
		// git unread: the weeks carry no stories, exact, or share
		const weeks = claims.weeks!.map((wk) => ({
			week: wk.week,
			start: wk.start,
			held_seconds: wk.held_seconds
		}));
		const unread: Report = { ...claiming, claims: { ...claims, weeks, drift: undefined } };
		expect(driftedIn(unread)).toEqual([]);
		const o = touchesDrift(unread, light) as Option;
		expect(o.series).toEqual([]);
		expect(o.xAxis).toMatchObject({ min: start, max: end });
		expect((holdTime(unread, light) as Option).series[0].data.length).toBe(6);
		// drift without a story done in the window: the shares alone
		const idle: Report = { ...claiming, claims: { ...claims, drift: [] } };
		expect((touchesDrift(idle, light) as Option).series.map((s) => s.name)).toEqual([
			'exact touches per week'
		]);
	});
	it('draws the claims charts empty over the window from an empty window or an older flai', () => {
		const spans = {
			parallelism: [midnight('2026-08-02'), end],
			'hold-time': [midnight('2026-07-27') - 3.5 * day, midnight('2026-08-31') + 3.5 * day],
			'touches-drift': [start, end]
		} as const;
		for (const r of [{ ...claiming, claims: empty }, forecasting]) {
			for (const kind of CLAIMS_KINDS) {
				const o = build(kind, r, light) as Option;
				expect(
					o.series.every((s) => s.data.length === 0),
					kind
				).toBe(true);
				expect(o.xAxis, kind).toMatchObject({ min: spans[kind][0], max: spans[kind][1] });
			}
		}
		expect((build('touches-drift', forecasting, light) as Option).series).toEqual([]);
		expect((build('parallelism', forecasting, light) as Option).series.length).toBe(2);
		// a flai older than S-0214 sends the days without weeks or held counts
		const older: Report = { ...forecasting, claims: { limit: 2, days: claims.days } };
		expect(hasClaims(claiming)).toBe(true);
		expect(hasClaims(older)).toBe(false);
		expect(hasClaims(forecasting)).toBe(false);
		expect((holdTime(older, light) as Option).series.every((s) => s.data.length === 0)).toBe(true);
	});
	it('reads the stories alone, with no filter, and plots no forecast error', () => {
		for (const kind of CLAIMS_KINDS) {
			expect(isClaimsKind(kind), kind).toBe(true);
			expect(isForecastKind(kind), kind).toBe(false);
			expect(controls(kind), kind).toEqual({
				type: false,
				epic: false,
				bucket: false,
				nature: false,
				model: false
			});
			expect(errorFacets(forecasting, kind), kind).toEqual({ natures: [], models: [] });
		}
		expect(FORECAST_KINDS.every(isForecastKind)).toBe(true);
	});
});

// flai stats --json --since 11d on this repository, 2026-10-07, trimmed: the strategic days, and of
// usage its totals. The planner spent from 4 October, the orchestrator on 7 October; on 26 to 28
// September the stories done carried no agents' usage, so those days have no cost per item.
const strategicDays: StrategicDay[] = [
	{
		date: '2026-09-26',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 17,
		cycle_time_seconds: 6918.64705882353
	},
	{
		date: '2026-09-27',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 1,
		cycle_time_seconds: 24374
	},
	{
		date: '2026-09-28',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 2,
		cycle_time_seconds: 77352.5
	},
	{
		date: '2026-09-29',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 30,
		cost_per_item: 5.7092,
		cycle_time_seconds: 3339.233333333333
	},
	{
		date: '2026-09-30',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 5,
		cost_per_item: 5.2406,
		cycle_time_seconds: 1104.2
	},
	{
		date: '2026-10-01',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 15,
		cost_per_item: 10.1552,
		cycle_time_seconds: 9572.533333333333
	},
	{
		date: '2026-10-02',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 11,
		cost_per_item: 7.7906,
		cycle_time_seconds: 13681.181818181818
	},
	{
		date: '2026-10-03',
		agents: {},
		cost: 0,
		seconds: 0,
		completed: 13,
		cost_per_item: 13.7473,
		cycle_time_seconds: 12154.384615384615
	},
	{
		date: '2026-10-04',
		agents: { planner: { cost: 8.7265, seconds: 3516, estimated: true } },
		cost: 8.7265,
		seconds: 3516,
		completed: 10,
		cost_per_item: 16.2799,
		cycle_time_seconds: 11044
	},
	{
		date: '2026-10-05',
		agents: { planner: { cost: 34.5781, seconds: 5576, estimated: true } },
		cost: 34.5781,
		seconds: 5576,
		completed: 15,
		cost_per_item: 7.302,
		cycle_time_seconds: 2108.133333333333
	},
	{
		date: '2026-10-06',
		agents: { planner: { cost: 56.688, seconds: 7156, estimated: true } },
		cost: 56.688,
		seconds: 7156,
		completed: 25,
		cost_per_item: 11.6428,
		cycle_time_seconds: 6123.8
	},
	{
		date: '2026-10-07',
		agents: {
			planner: { cost: 6.682, seconds: 578, estimated: true },
			orchestrator: { cost: 30.18, seconds: 28705, estimated: true }
		},
		cost: 36.862,
		seconds: 29283,
		completed: 18,
		cost_per_item: 11.9793,
		cycle_time_seconds: 5366.333333333333
	}
];
const strategic: Report = {
	...report,
	generated_at: '2026-10-07T09:21:02Z',
	window_days: 11,
	window_start: '2026-09-26T09:21:02Z',
	strategic_days: strategicDays,
	usage: {
		items: 128,
		tokens: 2789278953,
		cost: 1296.4735,
		seconds: 266955,
		estimated: true,
		models: []
	}
};

/** A story of this repository as flai stats --json reported it, with what it waited. */
const doneOn = (
	id: string,
	completed: string,
	cycle: number,
	threads: number | undefined,
	review: number | undefined,
	status = 'done'
): ItemMetrics => ({
	id,
	type: 'story',
	nature: 'feature',
	title: id,
	status,
	created: '2026-09-20T00:00:00Z',
	completed,
	cycle_time_seconds: status === 'done' ? cycle : undefined,
	blocked_seconds: 0,
	time_in_state_seconds: {},
	wait_threads_seconds: threads,
	wait_review_seconds: review
});
// The stories flai stats --json reported done on 26 to 28 September, 2026-10-07, with S-0139, the
// one cancelled; those done after 28 September left out, so their days have none done.
const waited: Report = {
	...strategic,
	items: [
		doneOn('S-0117', '2026-09-26T04:18:59Z', 2976, undefined, 2718),
		doneOn('S-0118', '2026-09-26T03:15:35Z', 945, 195, 33),
		doneOn('S-0119', '2026-09-26T03:23:38Z', 440, undefined, 25),
		doneOn('S-0120', '2026-09-26T05:39:27Z', 901, undefined, 11),
		doneOn('S-0121', '2026-09-26T06:02:02Z', 2242, 1093, 203),
		doneOn('S-0122', '2026-09-26T07:04:22Z', 3886, 1434, 274),
		doneOn('S-0123', '2026-09-26T07:27:34Z', 1351, undefined, 233),
		doneOn('S-0124', '2026-09-26T08:02:04Z', 2637, 1967, 118),
		doneOn('S-0125', '2026-09-26T17:48:17Z', 35087, 34515, 70),
		doneOn('S-0126', '2026-09-26T07:40:57Z', 613, undefined, 398),
		doneOn('S-0127', '2026-09-26T07:46:34Z', 356, undefined, 70),
		doneOn('S-0128', '2026-09-26T17:47:43Z', 35019, undefined, 34197),
		doneOn('S-0129', '2026-09-26T18:12:29Z', 1314, undefined, 517),
		doneOn('S-0130', '2026-09-26T20:14:28Z', 8458, 1457, 6968),
		doneOn('S-0131', '2026-09-26T21:02:26Z', 10683, 665, 10013),
		doneOn('S-0132', '2026-09-26T21:09:40Z', 10430, 602, 9791),
		doneOn('S-0133', '2026-09-26T17:53:35Z', 279, undefined, 37),
		doneOn('S-0139', '2026-09-26T20:27:09Z', 0, 0, undefined, 'cancelled'),
		doneOn('S-0140', '2026-09-27T03:56:25Z', 24374, 24271, 27),
		doneOn('S-0134', '2026-09-28T22:39:28Z', 153757, 129, 2014),
		doneOn('S-0135', '2026-09-28T22:56:03Z', 948, undefined, 79)
	],
	strategic_days: strategicDays.map((d) =>
		d.date <= '2026-09-28'
			? d
			: { ...d, completed: 0, cost_per_item: undefined, cycle_time_seconds: undefined }
	)
};

describe('strategic charts', () => {
	type Series = {
		name: string;
		type: string;
		stack?: string;
		itemStyle: { color: string };
		lineStyle?: { type: string };
		data: { value: [number, number]; at: string; items: number; estimated?: boolean }[];
		markLine?: { data: { name: string; yAxis: number }[] };
	};
	type Strategic = {
		series: Series[];
		legend: { show: boolean };
		useUTC?: boolean;
		xAxis: { type: string; min?: number; max?: number; minInterval: number };
		yAxis: { name: string; axisLabel: { formatter: (v: number) => string } };
		tooltip: { formatter: (p: unknown) => string };
	};
	const hour = 3600e3;

	it('strategic cost stacks what each strategic agent spent per day, with the cost per story', () => {
		const o = build('strategic-cost', strategic, light) as Strategic;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['planner', 'bar'],
			['orchestrator', 'bar'],
			['mean per story', 'line']
		]);
		expect(new Set(o.series.slice(0, 2).map((s) => s.stack))).toEqual(new Set(['strategic']));
		// one bar per day of the window, on its date's tick, 0 on a day a kind spent nothing
		expect(o.series[0].data.map((d) => d.at)).toEqual(
			strategicDays.map((d) => `${d.date}T00:00:00Z`)
		);
		expect(o.series[0].data.map((d) => d.value[0])).toEqual(
			strategicDays.map((d) => midnight(d.date))
		);
		expect(o.series[0].data.map((d) => d.value[1])).toEqual([
			0, 0, 0, 0, 0, 0, 0, 0, 8.7265, 34.5781, 56.688, 6.682
		]);
		expect(o.series[1].data.map((d) => d.value[1])).toEqual([...Array(11).fill(0), 30.18]);
		// every strategic entry is apportioned, so estimated, and the bars say so
		expect(o.series[0].data.map((d) => d.estimated)).toEqual([
			...Array(8).fill(undefined),
			true,
			true,
			true,
			true
		]);
		// the agents' cost per story done that day: no point on 26 to 28 September
		expect(o.series[2].data.map((d) => [d.at, d.value[1]])).toEqual([
			['2026-09-29T00:00:00Z', 5.7092],
			['2026-09-30T00:00:00Z', 5.2406],
			['2026-10-01T00:00:00Z', 10.1552],
			['2026-10-02T00:00:00Z', 7.7906],
			['2026-10-03T00:00:00Z', 13.7473],
			['2026-10-04T00:00:00Z', 16.2799],
			['2026-10-05T00:00:00Z', 7.302],
			['2026-10-06T00:00:00Z', 11.6428],
			['2026-10-07T00:00:00Z', 11.9793]
		]);
		// one axis, in dollars, over the window by the day, half a day either side
		expect(o.yAxis.name).toBe('US dollars');
		expect(o.yAxis.axisLabel.formatter(0.5)).toBe('$0.500');
		expect(o.useUTC).toBeUndefined();
		expect(o.xAxis).toMatchObject({
			type: 'time',
			minInterval: 24 * hour,
			min: midnight('2026-09-26') - 12 * hour,
			max: midnight('2026-10-07') + 12 * hour
		});
		expect(o.legend.show).toBe(true);
		expect(
			o.tooltip.formatter([
				{ seriesName: 'planner', data: o.series[0].data[10] },
				{ seriesName: 'orchestrator', data: o.series[1].data[10] },
				{ seriesName: 'mean per story', data: o.series[2].data[7] }
			])
		).toBe(
			'2026-10-06<br/>planner: $56.69 (estimated in part)<br/>orchestrator: $0.00<br/>mean per story: $11.64'
		);
		expect(STRATEGIC_KINDS).toEqual(['strategic-cost', 'strategic-use']);
		expect(titleOf('strategic-cost')).toBe('Strategic Cost');
		expect(titleOf('strategic-use')).toBe('Strategic Use');
		for (const k of STRATEGIC_KINDS)
			expect(controls(k)).toEqual({
				type: true,
				epic: false,
				bucket: false,
				nature: false,
				model: false
			});
	});
	it('strategic cost states what the strategic agents spent per story against the agents', () => {
		// $136.8546 over 162 stories done, against the agents' $1296.4735 over 128: what flai stats
		// prints as "$136.85 · 12h38m, beside 162 completed (usage $10.13 per item)"
		const ratio = strategicRatio(strategic)!;
		expect(ratio).toMatchObject({ cost: 136.8546, completed: 162 });
		expect(ratio.cost_per_item).toBeCloseTo(0.84478, 5);
		expect(ratio.agent_cost_per_item).toBeCloseTo(10.1287, 4);
		expect(ratio.share).toBeCloseTo(0.0834, 4);
		const o = build('strategic-cost', strategic, light) as Strategic;
		expect(o.series[2].markLine?.data).toEqual([
			{ yAxis: ratio.cost_per_item, name: 'strategic per story 8%' }
		]);
		// the type of the report names it
		const tasks = build('strategic-cost', { ...strategic, type: 'task' }, light) as Strategic;
		expect(tasks.series[2].name).toBe('mean per task');
		expect(tasks.series[2].markLine?.data[0].name).toBe('strategic per task 8%');
		// nothing completed, or nothing the agents spent: no ratio, and no line for it
		const idle: Report = {
			...strategic,
			strategic_days: strategicDays.map((d) => ({ ...d, completed: 0, cost_per_item: undefined }))
		};
		expect(strategicRatio(idle)).toBeUndefined();
		const none = build('strategic-cost', idle, light) as Strategic;
		expect(none.series[2].data).toEqual([]);
		expect(none.series[2].markLine).toBeUndefined();
		const unspent: Report = { ...strategic, usage: { ...strategic.usage!, items: 0, cost: 0 } };
		expect(strategicRatio(unspent)).toBeUndefined();
	});
	it('strategic cost gives each strategic agent its fixed colour, in a fixed order', () => {
		const analyzed: Report = {
			...strategic,
			strategic_days: strategicDays.map((d) =>
				d.date === '2026-10-02'
					? { ...d, agents: { analyzer: { cost: 1.25, seconds: 600 } }, cost: 1.25, seconds: 600 }
					: d
			)
		};
		const o = build('strategic-cost', analyzed, light) as Strategic;
		expect(o.series.map((s) => s.name)).toEqual([
			'planner',
			'orchestrator',
			'analyzer',
			'mean per story'
		]);
		expect(o.series.slice(0, 3).map((s) => s.itemStyle.color)).toEqual([
			CATEGORICAL.light[6],
			CATEGORICAL.light[5],
			CATEGORICAL.light[4]
		]);
		const night = build('strategic-cost', analyzed, dark) as Strategic;
		expect(night.series[2].itemStyle.color).toBe(CATEGORICAL.dark[4]);
		// an entry not estimated is not marked
		expect(o.tooltip.formatter([{ seriesName: 'analyzer', data: o.series[2].data[6] }])).toBe(
			'2026-10-02<br/>analyzer: $1.25'
		);
		// the strategic ratio counts it
		expect(strategicRatio(analyzed)!.cost).toBe(138.1046);
	});
	it('draws an empty strategic cost from a flai older than S-0205', () => {
		const older = { ...strategic, strategic_days: undefined };
		const o = build('strategic-cost', older, light) as Strategic;
		expect(o.series).toEqual([]);
		expect(o.legend.show).toBe(false);
		expect(strategicRatio(older)).toBeUndefined();
		expect(
			normalise({ ...strategic, strategic_days: null } as unknown as Report).strategic_days
		).toEqual([]);
		const use = build('strategic-use', older, light) as Strategic;
		expect(use.series).toEqual([]);
		expect(strategicRows(older)).toEqual([]);
	});
	it('strategic use stacks the strategic agents hours per day, with mean cycle time and waiting', () => {
		const o = build('strategic-use', waited, light) as Strategic;
		expect(o.series.map((s) => [s.name, s.type])).toEqual([
			['planner', 'bar'],
			['orchestrator', 'bar'],
			['mean cycle time per story', 'line'],
			['mean waiting per story', 'line']
		]);
		expect(new Set(o.series.slice(0, 2).map((s) => s.stack))).toEqual(new Set(['strategic']));
		// one bar per day of the window, in hours, 0 on a day a kind worked none, estimated as spent
		expect(o.series[0].data.map((d) => d.value[0])).toEqual(
			strategicDays.map((d) => midnight(d.date))
		);
		const h = (s: number) => s / 3600;
		expect(o.series[0].data.map((d) => d.value[1])).toEqual(
			[...Array(8).fill(0), 3516, 5576, 7156, 578].map(h)
		);
		expect(o.series[1].data.map((d) => d.value[1])).toEqual([...Array(11).fill(0), 28705].map(h));
		expect(o.series[1].data[11].estimated).toBe(true);
		// the same colours as strategic cost
		const c = build('strategic-cost', waited, light) as Strategic;
		expect(o.series.slice(0, 2).map((s) => s.itemStyle.color)).toEqual(
			c.series.slice(0, 2).map((s) => s.itemStyle.color)
		);
		// the lines have a point only on the days a story was done, in hours
		const at = ['2026-09-26T00:00:00Z', '2026-09-27T00:00:00Z', '2026-09-28T00:00:00Z'];
		expect(o.series[2].data.map((d) => d.at)).toEqual(at);
		expect(o.series[3].data.map((d) => d.value[0])).toEqual(at.map(midnight));
		expect(o.series[2].data.map((d) => d.value[1])).toEqual(
			[6918.64705882353, 24374, 77352.5].map(h)
		);
		// one axis, in hours, over the window by the day
		expect(o.yAxis.name).toBe('hours');
		expect(o.yAxis.axisLabel.formatter(6.75)).toBe('6.8h');
		expect(o.xAxis).toMatchObject({
			type: 'time',
			minInterval: 24 * hour,
			min: midnight('2026-09-26') - 12 * hour,
			max: midnight('2026-10-07') + 12 * hour
		});
		expect(o.legend.show).toBe(true);
		expect(
			o.tooltip.formatter([
				{ seriesName: 'planner', data: o.series[0].data[11] },
				{ seriesName: 'orchestrator', data: o.series[1].data[11] }
			])
		).toBe(
			'2026-10-07<br/>planner: 0.2h (estimated in part)<br/>orchestrator: 8h (estimated in part)'
		);
		expect(
			o.tooltip.formatter([
				{ seriesName: 'mean cycle time per story', data: o.series[2].data[0] },
				{ seriesName: 'mean waiting per story', data: o.series[3].data[0] }
			])
		).toBe('2026-09-26<br/>mean cycle time per story: 1.9h<br/>mean waiting per story: 1.8h');
	});
	it('strategic use waits the mean of the threads and review waits of the stories done that day', () => {
		const rows = strategicRows(waited);
		expect(rows.map((d) => d.date)).toEqual(strategicDays.map((d) => d.date));
		// 26 September: 17 stories done, the cancelled S-0139 left out, those done before the window's
		// start at 09:21 counted, as flai counts them; a missing wait counts 0
		expect(rows[0]).toMatchObject({ completed: 17, wait_seconds: 107604 / 17 });
		// 27 September: S-0140 alone, 24271 on threads and 27 in review
		expect(rows[1].wait_seconds).toBe(24298);
		// 28 September: S-0134 129 + 2014, S-0135 no thread and 79 in review
		expect(rows[2].wait_seconds).toBe(1111);
		// no story done, no wait
		expect(rows.slice(3).every((d) => d.wait_seconds === undefined)).toBe(true);
		// the items that wait are those the day's completed counts, and whose cycle time it means
		const done = waited.items.filter((i) => i.status === 'done');
		for (const d of rows.slice(0, 3)) {
			const of = done.filter((i) => i.completed!.startsWith(d.date));
			expect(of.length).toBe(d.completed);
			expect(of.reduce((n, i) => n + i.cycle_time_seconds!, 0) / of.length).toBeCloseTo(
				d.cycle_time_seconds!,
				9
			);
		}
		const o = build('strategic-use', waited, light) as Strategic;
		expect(o.series[3].data.map((d) => d.value[1])).toEqual(
			[107604 / 17, 24298, 1111].map((s) => s / 3600)
		);
	});
});
