// Chart option builders: pure functions from a flai stats report to an
// ECharts option, so they are unit-testable without a DOM. One y-axis per
// chart, thin marks, legends for two or more series, tooltips everywhere.
import { colorFor, modelSlot, NATURE_SLOT, STATE_SLOT, type Theme } from './palette';
import { count, dollars } from '$lib/usage';

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
	usage?: ItemUsage;
};
/** What a model spent on an item, or on the items in a window (S-0143). */
export type ModelSpend = {
	model: string;
	tokens: number;
	cost: number;
	tokens_per_hour?: number;
	items?: number;
};
export type ItemUsage = {
	source: string;
	tokens: number;
	cost: number;
	seconds: number;
	tokens_per_hour?: number;
	estimated?: boolean;
	models: ModelSpend[];
};
export type SpendPoint = { at: string; id: string; done: number; tokens: number; cost: number };
export type UsageReport = {
	items: number;
	tokens: number;
	cost: number;
	seconds: number;
	estimated?: boolean;
	models: ModelSpend[];
	done: SpendPoint[];
	by_model: Record<string, SpendPoint[]>;
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
		burnup: r.burnup ?? {},
		usage: {
			items: 0,
			tokens: 0,
			cost: 0,
			seconds: 0,
			...r.usage,
			models: r.usage?.models ?? [],
			done: r.usage?.done ?? [],
			by_model: r.usage?.by_model ?? {}
		}
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
	/** Absent from a flai older than S-0143. */
	usage?: UsageReport;
};

export const KINDS = [
	'cycle-time',
	'burn-up',
	'cfd',
	'time-in-state',
	'throughput',
	'aging',
	'estimates',
	'token-rate',
	'cost',
	'completion-time',
	'completion-cost'
] as const;
export type Kind = (typeof KINDS)[number];
export const TITLES: Record<Kind, string> = {
	'cycle-time': 'Cycle time',
	'burn-up': 'Burn-up',
	cfd: 'Cumulative flow',
	'time-in-state': 'Time in state',
	throughput: 'Throughput',
	aging: 'Aging work in progress',
	estimates: 'Estimate versus actual',
	'token-rate': 'Token rate',
	cost: 'Cost',
	'completion-time': 'Completion over time',
	'completion-cost': 'Completion against cost'
};
/** The usage charts (S-0143): they need items that carry usage. */
export const USAGE_KINDS: readonly Kind[] = [
	'token-rate',
	'cost',
	'completion-time',
	'completion-cost'
];

const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done'];
/** An item type's plural: epics, stories, tasks. */
export const plural = (type: string) => (type === 'story' ? 'stories' : `${type}s`);
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
			name: plural(r.type),
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

/** Every model the report names, in order of name. */
export function models(r: Report): string[] {
	const names = new Set<string>();
	for (const i of r.items) for (const m of i.usage?.models ?? []) names.add(m.model);
	for (const m of r.usage?.models ?? []) names.add(m.model);
	return [...names].sort();
}
/** A model's colour: its family's fixed slot, whichever models a report or a filter leaves. */
function modelColor(t: Theme, model: string): string {
	return t.series[modelSlot(model)];
}
const withUsage = (r: Report, epic?: string) =>
	r.items.filter((i) => i.usage && i.usage.models.length > 0 && (!epic || i.parent === epic));

/**
 * Token rate: one point per item with usage and agent time, per model, x when it was completed
 * (started, while open), y the model's tokens per agent hour. At most three models are coloured
 * (the all-pairs rule); the rest fold into other.
 */
export function tokenRate(r: Report, t: Theme, epic?: string): Opt {
	const all = models(r);
	const pts = withUsage(r, epic).flatMap((i) =>
		i
			.usage!.models.filter((m) => m.tokens_per_hour !== undefined && (i.completed || i.started))
			.map((m) => ({ i, m }))
	);
	const present = all.filter((name) => pts.some((p) => p.m.model === name));
	const top = present.slice(0, 3);
	const groups = [...top, ...(present.length > 3 ? ['other'] : [])];
	const series = groups.map((g) => ({
		name: g,
		type: 'scatter',
		symbolSize: 10,
		itemStyle: {
			color: g === 'other' ? t.textSecondary : modelColor(t, g),
			borderColor: t.surface,
			borderWidth: 2
		},
		data: pts
			.filter((p) => (g === 'other' ? !top.includes(p.m.model) : p.m.model === g))
			.map((p) => ({
				value: [p.i.completed ?? p.i.started, p.m.tokens_per_hour! / 1e6],
				id: p.i.id,
				title: p.i.title,
				model: p.m.model
			}))
	}));
	return base(t, {
		legend: legend(t, groups.length > 1),
		tooltip: tooltip(t, {
			formatter: (p: {
				data: { id: string; title: string; model: string; value: [string, number] };
			}) =>
				`${p.data.id} ${p.data.title}<br/>${p.data.model}: ${count(p.data.value[1] * 1e6)} tokens per agent hour`
		}),
		xAxis: axisX(t, { type: 'time' }),
		yAxis: axisY(t, {
			type: 'value',
			name: 'million tokens per agent hour',
			nameTextStyle: { color: t.textSecondary }
		}),
		series
	});
}

/**
 * Cost: one bar per completed item with usage, in order of completion, stacked by model. An item
 * whose cost is estimated in part carries an asterisk on its label and says so in the tooltip.
 */
export function cost(r: Report, t: Theme, epic?: string): Opt {
	const all = models(r);
	const done = withUsage(r, epic)
		.filter((i) => i.completed)
		.sort((a, b) => (a.completed! < b.completed! ? -1 : a.completed! > b.completed! ? 1 : 0));
	const present = all.filter((name) =>
		done.some((i) => i.usage!.models.some((m) => m.model === name))
	);
	const series = present.map((name) => ({
		name,
		type: 'bar',
		stack: 'cost',
		barMaxWidth: 24,
		itemStyle: { color: modelColor(t, name), borderColor: t.surface, borderWidth: 1 },
		data: done.map((i) => i.usage!.models.find((m) => m.model === name)?.cost ?? 0)
	}));
	return base(t, {
		legend: legend(t, present.length > 1),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'shadow' },
			formatter: (ps: { dataIndex: number; seriesName: string; value: number }[]) => {
				const i = done[ps[0]?.dataIndex ?? 0];
				if (!i) return '';
				const lines = ps
					.filter((p) => p.value > 0)
					.map((p) => `${p.seriesName}: ${dollars(p.value)}`);
				return `${i.id} ${i.title}<br/>${lines.join('<br/>')}<br/>total ${dollars(i.usage!.cost)}${i.usage!.estimated ? ' (estimated in part)' : ''}`;
			}
		}),
		xAxis: axisX(t, {
			type: 'category',
			data: done.map((i) => (i.usage!.estimated ? `${i.id}*` : i.id)),
			axisLabel: { color: t.textSecondary, rotate: done.length > 12 ? 45 : 0 }
		}),
		yAxis: axisY(t, {
			type: 'value',
			name: 'US dollars',
			nameTextStyle: { color: t.textSecondary }
		}),
		series
	});
}

/** Items done cumulatively, per model, against time or against that model's cumulative cost. */
function completion(r: Report, t: Theme, against: 'time' | 'cost'): Opt {
	const all = models(r);
	const byModel = r.usage?.by_model ?? {};
	const present = all.filter((name) => (byModel[name] ?? []).length > 0);
	const series = present.map((name) => ({
		name,
		type: 'line',
		step: against === 'time' ? 'end' : undefined,
		showSymbol: (byModel[name] ?? []).length < 40,
		symbolSize: 8,
		lineStyle: { width: 2, color: modelColor(t, name) },
		itemStyle: { color: modelColor(t, name), borderColor: t.surface, borderWidth: 2 },
		data: (byModel[name] ?? []).map((p) => ({
			value: [against === 'time' ? p.at : p.cost, p.done],
			id: p.id,
			cost: p.cost,
			at: p.at
		}))
	}));
	return base(t, {
		legend: legend(t, present.length > 1),
		tooltip: tooltip(t, {
			formatter: (p: {
				seriesName: string;
				data: { id: string; cost: number; at: string; value: [unknown, number] };
			}) =>
				`${p.seriesName}: ${p.data.value[1]} done by ${p.data.id}<br/>${dollars(p.data.cost)} spent · ${p.data.at.slice(0, 10)}`
		}),
		xAxis: axisX(
			t,
			against === 'time'
				? { type: 'time' }
				: {
						type: 'value',
						name: 'US dollars spent',
						nameLocation: 'middle',
						nameGap: 28,
						nameTextStyle: { color: t.textSecondary }
					}
		),
		yAxis: axisY(t, {
			type: 'value',
			name: `${plural(r.type)} done`,
			nameTextStyle: { color: t.textSecondary },
			minInterval: 1
		}),
		series
	});
}
export const completionTime = (r: Report, t: Theme) => completion(r, t, 'time');
export const completionCost = (r: Report, t: Theme) => completion(r, t, 'cost');

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
		case 'token-rate':
			return tokenRate(r, t, epic);
		case 'cost':
			return cost(r, t, epic);
		case 'completion-time':
			return completionTime(r, t);
		case 'completion-cost':
			return completionCost(r, t);
	}
}
