// Chart option builders: pure functions from a flai stats report to an
// ECharts option, so they are unit-testable without a DOM. One y-axis per
// chart, thin marks, legends for two or more series, tooltips everywhere.
import { colorFor, NATURE_SLOT, STATE_SLOT, type Theme } from './palette';

export type Distribution = {
	count: number;
	p50_seconds: number;
	p85_seconds: number;
	max_seconds: number;
	mean_seconds: number;
};
export type Summary = {
	group?: string;
	completed: number;
	cancelled: number;
	cancellation_rate: number;
	throughput_per_week: number;
	wip: number;
	cycle_time: Distribution;
	lead_time: Distribution;
	queue_time: Distribution;
	flow_efficiency: number;
	time_in_state_share: Record<string, number>;
};
export type ItemMetrics = {
	id: string;
	type: string;
	nature: string;
	title: string;
	status: string;
	parent?: string;
	created: string;
	started?: string;
	completed?: string;
	lead_time_seconds?: number;
	cycle_time_seconds?: number;
	queue_time_seconds?: number;
	blocked_seconds: number;
	time_in_state_seconds: Record<string, number>;
	estimate_seconds?: number;
	estimate_error?: number;
	age_seconds?: number;
};
/**
 * Older flai builds emit null for empty lists; give every list the charts
 * iterate a value so a selection with no items draws an empty chart instead
 * of throwing (S-0045).
 */
export function normalise(r: Report): Report {
	return {
		...r,
		items: r.items ?? [],
		throughput: r.throughput ?? [],
		cfd: r.cfd ?? [],
		aging: r.aging ?? [],
		burnup: r.burnup ?? {}
	};
}

export type Report = {
	generated_at: string;
	type: string;
	window_days: number;
	window_start: string;
	summary: Summary;
	groups?: Summary[];
	items: ItemMetrics[];
	throughput: { week: string; start: string; done: number; by_nature: Record<string, number> }[];
	burnup: Record<string, { date: string; scope?: number; done?: number }[]>;
	cfd: { date: string; counts: Record<string, number> }[];
	aging: {
		id: string;
		title: string;
		status: string;
		age_seconds: number;
		over_p85: boolean;
		blocked: boolean;
		nature: string;
		parent?: string;
		age: string;
	}[];
};

export const KINDS = [
	'cycle-time',
	'burn-up',
	'cfd',
	'time-in-state',
	'throughput',
	'aging',
	'estimates'
] as const;
export type Kind = (typeof KINDS)[number];
export const TITLES: Record<Kind, string> = {
	'cycle-time': 'Cycle time',
	'burn-up': 'Burn-up',
	cfd: 'Cumulative flow',
	'time-in-state': 'Time in state',
	throughput: 'Throughput',
	aging: 'Aging work in progress',
	estimates: 'Estimate versus actual'
};

const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done'];
export const hours = (s: number) => Math.round((s / 3600) * 10) / 10;
export const days = (s: number) => Math.round((s / 86400) * 100) / 100;

/** Human duration for axes and tooltips. */
export function human(seconds: number): string {
	if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
	if (seconds < 86400) return `${Math.round((seconds / 3600) * 10) / 10}h`;
	return `${Math.round((seconds / 86400) * 10) / 10}d`;
}

type Opt = Record<string, unknown>;
function base(t: Theme, extra: Opt = {}): Opt {
	return {
		backgroundColor: 'transparent',
		textStyle: { color: t.text },
		animationDuration: 300,
		grid: { left: 56, right: 24, top: 40, bottom: 48, containLabel: false },
		tooltip: {
			trigger: 'item',
			backgroundColor: t.surface,
			borderColor: t.grid,
			textStyle: { color: t.text }
		},
		...extra
	};
}
function tooltip(t: Theme, extra: Opt = {}): Opt {
	return {
		trigger: 'item',
		backgroundColor: t.surface,
		borderColor: t.grid,
		textStyle: { color: t.text },
		...extra
	};
}
function axisX(t: Theme, extra: Opt = {}): Opt {
	return {
		axisLine: { lineStyle: { color: t.grid } },
		axisLabel: { color: t.textSecondary },
		splitLine: { show: false },
		...extra
	};
}
function axisY(t: Theme, extra: Opt = {}): Opt {
	return {
		axisLine: { show: false },
		axisLabel: { color: t.textSecondary },
		splitLine: { lineStyle: { color: t.grid } },
		...extra
	};
}
function legend(t: Theme, show: boolean): Opt {
	return {
		show,
		top: 0,
		textStyle: { color: t.textSecondary },
		icon: 'circle',
		itemWidth: 10,
		itemHeight: 10
	};
}
function refLine(t: Theme, data: Opt[], formatter: string): Opt {
	return {
		silent: true,
		symbol: 'none',
		lineStyle: { type: 'dashed', color: t.textSecondary, width: 1 },
		label: { color: t.textSecondary, formatter },
		data
	};
}

/** Cycle time scatter: one point per completed item, p50 and p85 reference lines, coloured by nature (at most three groups per the all-pairs rule; the rest fold into Other). */
export function cycleTime(r: Report, t: Theme, epic?: string): Opt {
	const done = r.items.filter(
		(i) => i.completed && i.cycle_time_seconds !== undefined && (!epic || i.parent === epic)
	);
	const natures = [...new Set(done.map((i) => i.nature))];
	const top = natures.slice(0, 3);
	const groups = [...top, ...(natures.length > 3 ? ['other'] : [])];
	const series = groups.map((g, gi) => ({
		name: g,
		type: 'scatter',
		symbolSize: 10,
		itemStyle: { color: colorFor(t, NATURE_SLOT, g, gi), borderColor: t.surface, borderWidth: 2 },
		data: done
			.filter((i) => (g === 'other' ? !top.includes(i.nature) : i.nature === g))
			.map((i) => ({
				value: [i.completed, days(i.cycle_time_seconds!)],
				id: i.id,
				title: i.title
			})),
		markLine:
			gi === 0
				? refLine(
						t,
						[
							{ yAxis: days(r.summary.cycle_time.p50_seconds), name: 'p50' },
							{ yAxis: days(r.summary.cycle_time.p85_seconds), name: 'p85' }
						],
						'{b}'
					)
				: undefined
	}));
	return base(t, {
		legend: legend(t, groups.length > 1),
		tooltip: tooltip(t, {
			formatter: (p: { data: { id: string; title: string; value: [string, number] } }) =>
				`${p.data.id} ${p.data.title}<br/>${p.data.value[1]} days · ${p.data.value[0].slice(0, 10)}`
		}),
		xAxis: axisX(t, { type: 'time' }),
		yAxis: axisY(t, { type: 'value', name: 'days', nameTextStyle: { color: t.textSecondary } }),
		series
	});
}

/** Burn-up: scope and done lines for the whole set or one epic. */
export function burnUp(r: Report, t: Theme, epic = 'all'): Opt {
	const pts = r.burnup[epic] ?? r.burnup.all ?? [];
	const line = (name: string, key: 'scope' | 'done', slot: number) => ({
		name,
		type: 'line',
		showSymbol: false,
		lineStyle: { width: 2, color: t.series[slot] },
		itemStyle: { color: t.series[slot] },
		data: pts.map((p) => [p.date, p[key] ?? 0])
	});
	return base(t, {
		legend: legend(t, true),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'cross', label: { backgroundColor: t.grid, color: t.text } }
		}),
		xAxis: axisX(t, { type: 'time' }),
		yAxis: axisY(t, {
			type: 'value',
			name: r.type + 's',
			nameTextStyle: { color: t.textSecondary },
			minInterval: 1
		}),
		series: [line('scope', 'scope', 6), line('done', 'done', 2)]
	});
}

/** Cumulative flow: stacked areas per state with a surface-coloured seam between bands. */
export function cfd(r: Report, t: Theme): Opt {
	const series = STATES.map((st) => ({
		name: st,
		type: 'line',
		stack: 'flow',
		showSymbol: false,
		lineStyle: { width: 2, color: t.surface },
		areaStyle: { color: colorFor(t, STATE_SLOT, st, 0), opacity: 1 },
		itemStyle: { color: colorFor(t, STATE_SLOT, st, 0) },
		emphasis: { focus: 'series' },
		data: r.cfd.map((d) => [d.date, d.counts[st] ?? 0])
	}));
	return base(t, {
		legend: legend(t, true),
		tooltip: tooltip(t, { trigger: 'axis', axisPointer: { type: 'line' } }),
		xAxis: axisX(t, { type: 'time' }),
		yAxis: axisY(t, { type: 'value', minInterval: 1 }),
		series
	});
}

/** Time in state: one stacked bar per completed item, hours per state. */
export function timeInState(r: Report, t: Theme, epic?: string): Opt {
	const done = r.items.filter((i) => i.completed && (!epic || i.parent === epic));
	const series = STATES.map((st) => ({
		name: st,
		type: 'bar',
		stack: 'state',
		barMaxWidth: 24,
		itemStyle: { color: colorFor(t, STATE_SLOT, st, 0), borderColor: t.surface, borderWidth: 1 },
		data: done.map((i) => hours(i.time_in_state_seconds[st] ?? 0))
	}));
	return base(t, {
		legend: legend(t, true),
		tooltip: tooltip(t, { trigger: 'axis', axisPointer: { type: 'shadow' } }),
		xAxis: axisX(t, {
			type: 'category',
			data: done.map((i) => i.id),
			axisLabel: { color: t.textSecondary, rotate: done.length > 12 ? 45 : 0 }
		}),
		yAxis: axisY(t, { type: 'value', name: 'hours', nameTextStyle: { color: t.textSecondary } }),
		series
	});
}

/** Time-in-state share as one horizontal 100% bar (the process optimisation view). */
export function stateShare(r: Report, t: Theme): Opt {
	const share = r.summary.time_in_state_share ?? {};
	return base(t, {
		grid: { left: 24, right: 24, top: 40, bottom: 24 },
		legend: legend(t, true),
		tooltip: tooltip(t, {
			formatter: (p: { seriesName: string; value: number }) =>
				`${p.seriesName}: ${Math.round(p.value * 100)}%`
		}),
		xAxis: axisX(t, {
			type: 'value',
			max: 1,
			axisLabel: { color: t.textSecondary, formatter: (v: number) => `${Math.round(v * 100)}%` }
		}),
		yAxis: axisY(t, {
			type: 'category',
			data: ['share'],
			axisLabel: { show: false },
			splitLine: { show: false }
		}),
		series: STATES.map((st) => ({
			name: st,
			type: 'bar',
			stack: 'share',
			barMaxWidth: 24,
			itemStyle: { color: colorFor(t, STATE_SLOT, st, 0), borderColor: t.surface, borderWidth: 1 },
			label: {
				show: (share[st] ?? 0) > 0.08,
				color: t.surface,
				formatter: () => `${st} ${Math.round((share[st] ?? 0) * 100)}%`
			},
			data: [share[st] ?? 0]
		}))
	});
}

/** Throughput: completed items per ISO week, stacked by nature (fixed slots). */
export function throughput(r: Report, t: Theme): Opt {
	const natures = [...new Set(r.throughput.flatMap((w) => Object.keys(w.by_nature)))].sort();
	const series = natures.map((n, i) => ({
		name: n,
		type: 'bar',
		stack: 'week',
		barMaxWidth: 24,
		itemStyle: { color: colorFor(t, NATURE_SLOT, n, i), borderColor: t.surface, borderWidth: 1 },
		data: r.throughput.map((w) => w.by_nature[n] ?? 0)
	}));
	return base(t, {
		legend: legend(t, natures.length > 1),
		tooltip: tooltip(t, { trigger: 'axis', axisPointer: { type: 'shadow' } }),
		xAxis: axisX(t, { type: 'category', data: r.throughput.map((w) => w.week) }),
		yAxis: axisY(t, {
			type: 'value',
			minInterval: 1,
			name: 'done',
			nameTextStyle: { color: t.textSecondary }
		}),
		series
	});
}

/** Aging WIP: horizontal bars of age since started against the p85 line; over-p85 items use the reserved red. */
export function aging(r: Report, t: Theme, epic?: string): Opt {
	const rows = r.aging.filter((a) => !epic || a.parent === epic);
	const p85 = days(r.summary.cycle_time.p85_seconds);
	return base(t, {
		grid: { left: 72, right: 24, top: 24, bottom: 40 },
		tooltip: tooltip(t, {
			formatter: (p: { name: string; value: number; dataIndex: number }) =>
				`${p.name} ${rows[p.dataIndex].title}<br/>${p.value} days${rows[p.dataIndex].blocked ? ' · blocked' : ''}`
		}),
		xAxis: axisX(t, {
			type: 'value',
			name: 'days since started',
			nameTextStyle: { color: t.textSecondary }
		}),
		yAxis: axisY(t, {
			type: 'category',
			data: rows.map((a) => a.id),
			inverse: true,
			splitLine: { show: false }
		}),
		series: [
			{
				name: 'age',
				type: 'bar',
				barMaxWidth: 24,
				itemStyle: { borderRadius: [0, 4, 4, 0] },
				data: rows.map((a) => ({
					value: days(a.age_seconds),
					itemStyle: { color: a.over_p85 ? t.series[7] : colorFor(t, STATE_SLOT, a.status, 0) }
				})),
				markLine: p85 > 0 ? refLine(t, [{ xAxis: p85 }], 'p85') : undefined
			}
		]
	});
}

/** Estimate versus actual: one point per estimated, completed item; the diagonal is the perfect estimate. */
export function estimates(r: Report, t: Theme): Opt {
	const pts = r.items.filter((i) => i.estimate_seconds && i.cycle_time_seconds !== undefined);
	const max = Math.max(
		1,
		...pts.flatMap((i) => [hours(i.estimate_seconds!), hours(i.cycle_time_seconds!)])
	);
	return base(t, {
		tooltip: tooltip(t, {
			formatter: (p: { data: { id: string; value: [number, number] } }) =>
				`${p.data.id}<br/>estimated ${p.data.value[0]}h · actual ${p.data.value[1]}h`
		}),
		xAxis: axisX(t, {
			type: 'value',
			name: 'estimated hours',
			nameTextStyle: { color: t.textSecondary },
			max
		}),
		yAxis: axisY(t, {
			type: 'value',
			name: 'actual hours',
			nameTextStyle: { color: t.textSecondary },
			max
		}),
		series: [
			{
				name: 'items',
				type: 'scatter',
				symbolSize: 10,
				itemStyle: { color: t.series[0], borderColor: t.surface, borderWidth: 2 },
				data: pts.map((i) => ({
					value: [hours(i.estimate_seconds!), hours(i.cycle_time_seconds!)],
					id: i.id
				})),
				markLine: {
					silent: true,
					symbol: 'none',
					lineStyle: { type: 'dashed', color: t.textSecondary, width: 1 },
					data: [[{ coord: [0, 0] }, { coord: [max, max] }]]
				}
			}
		]
	});
}

export function build(kind: Kind, report: Report, t: Theme, epic?: string): Opt {
	const r = normalise(report);
	switch (kind) {
		case 'cycle-time':
			return cycleTime(r, t, epic);
		case 'burn-up':
			return burnUp(r, t, epic || 'all');
		case 'cfd':
			return cfd(r, t);
		case 'time-in-state':
			return timeInState(r, t, epic);
		case 'throughput':
			return throughput(r, t);
		case 'aging':
			return aging(r, t, epic);
		case 'estimates':
			return estimates(r, t);
	}
}
