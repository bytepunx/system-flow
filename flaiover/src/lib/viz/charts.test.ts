import { describe, expect, it } from 'vitest';
import {
	aging,
	build,
	burnUp,
	cfd,
	cycleTime,
	estimates,
	human,
	KINDS,
	stateShare,
	throughput,
	timeInState,
	type Report
} from './charts';
import { CATEGORICAL, theme } from './palette';

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
			estimate_error: 0.3
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
			time_in_state_seconds: { backlog: 1800, ready: 86400, 'in-progress': 36000, review: 50400 }
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
	]
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
		expect(o.series[1].itemStyle.color).toBe(CATEGORICAL.light[2]);
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
	it('dark theme swaps the palette and surface', () => {
		expect(dark.series[0]).toBe(CATEGORICAL.dark[0]);
		expect(dark.surface).not.toBe(light.surface);
		expect(human(93600)).toBe('1.1d');
		expect(human(2700)).toBe('45m');
	});
});
