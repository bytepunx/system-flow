import { describe, expect, it } from 'vitest';
import {
	aging,
	build,
	burnUp,
	cfd,
	completionCost,
	completionTime,
	cost,
	cycleTime,
	estimates,
	human,
	KINDS,
	normalise,
	stateShare,
	throughput,
	timeInState,
	tokenRate,
	models,
	type Report
} from './charts';
import { CATEGORICAL, modelSlot, theme } from './palette';

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
		const t = timeInState(report, light) as {
			xAxis: { data: string[] };
			series: { data: number[] }[];
		};
		expect(t.xAxis.data).toEqual(['S-001', 'S-002']);
		expect(t.series[2].data).toEqual([24, 10]);
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
		const tr = tokenRate(report, light) as {
			series: { name: string; itemStyle: { color: string }; data: { value: [string, number] }[] }[];
		};
		expect(tr.series.map((s) => s.name)).toEqual(['claude-haiku-4-5', 'claude-opus-5-5']);
		expect(tr.series[1].itemStyle.color).toBe(CATEGORICAL.light[0]);
		expect(tr.series[0].itemStyle.color).toBe(CATEGORICAL.light[5]);
		expect(tr.series[1].data.map((d) => d.value[1])).toEqual([2, 1]);
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
