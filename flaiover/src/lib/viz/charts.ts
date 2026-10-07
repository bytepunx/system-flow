// Chart option builders: pure functions from a flai stats report to an
// ECharts option, so they are unit-testable without a DOM. One y-axis per
// chart, save the share on time beside Delivery Accuracy's errors, thin
// marks, legends for two or more series, tooltips everywhere.
import {
	colorFor,
	modelSlot,
	modelSymbol,
	NATURE_SLOT,
	STATE_SLOT,
	STRATEGIC_SLOT,
	TYPE_SLOT,
	TYPE_SYMBOL,
	type Theme
} from './palette';
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
	/** The forecast, and the errors against it and the estimate (S-0205); absent from older flai. */
	forecast_seconds?: number;
	forecast_error_seconds?: number;
	delivery_error_seconds?: number;
	estimate_error_seconds?: number;
	/** The item's agent model, `(none)` without one (ADR-0111); absent from an older flai. */
	model?: string;
	age_seconds?: number;
	usage?: ItemUsage;
};
/** The count of a set of errors and the p50 and p85 of their absolute values, absent when empty. */
export type ErrorSpread = { count: number; p50_seconds?: number; p85_seconds?: number };
/** One kind of error spread over the items that have it, in all, per nature, and per model. */
export type ErrorStats = ErrorSpread & {
	by_nature: Record<string, ErrorSpread>;
	by_model: Record<string, ErrorSpread>;
};
/** How far forecasts and estimates were from what happened, over the items done in the window. */
export type Forecasts = { forecast: ErrorStats; delivery: ErrorStats; estimate: ErrorStats };
/** The per-item error fields the planning charts read. */
export type ErrorField =
	| 'forecast_error_seconds'
	| 'delivery_error_seconds'
	| 'estimate_error_seconds';
/** What a model spent on an item, or on the items in a window (S-0143). */
export type ModelSpend = {
	model: string;
	tokens: number;
	cost: number;
	/** Absent from a flai older than S-0163, which sends tokens_per_hour. */
	tokens_per_minute?: number;
	tokens_per_hour?: number;
	items?: number;
};
export type ItemUsage = {
	source: string;
	tokens: number;
	cost: number;
	seconds: number;
	tokens_per_minute?: number;
	tokens_per_hour?: number;
	estimated?: boolean;
	models: ModelSpend[];
	/** Absent from a flai older than S-0225. */
	strategic?: StrategicShare[];
};
/**
 * What one kind of strategic agent spent on an item, apart from its agents' figures and every
 * per-model value (S-0225, ADR-0083).
 */
export type StrategicShare = {
	kind: string;
	tokens: number;
	cost: number;
	seconds: number;
	estimated: boolean;
};
/** What strategic agents spent on the items of a set that carry it, apart from its other values. */
export type StrategicSpend = { items: number; tokens: number; cost: number; seconds: number };
/**
 * What was spent on a set of items, and what that comes to per item, per minute of agent work,
 * and per dollar; a value whose divisor is zero is absent (S-0163, ADR-0053).
 */
export type Spend = {
	items: number;
	tokens: number;
	cost: number;
	seconds: number;
	estimated?: boolean;
	tokens_per_item?: number;
	cost_per_item?: number;
	/** Absent from a flai older than S-0169; see minutesPerItem. */
	minutes_per_item?: number;
	tokens_per_minute?: number;
	tokens_per_dollar?: number;
};
/**
 * The mean agent minutes an item took (S-0169), worked out from the seconds when a flai older
 * than S-0169 does not send it; none when no item took agent time.
 */
export function minutesPerItem(s: Spend): number | undefined {
	if (s.minutes_per_item !== undefined) return s.minutes_per_item;
	return s.seconds > 0 && s.items > 0 ? s.seconds / 60 / s.items : undefined;
}
export type ModelShare = Spend & { model: string };
/** The items done in one bucket of time, which starts at `at`, with the running means so far. */
export type Bucket = Spend & {
	at: string;
	mean_tokens: number;
	mean_cost: number;
	models?: ModelShare[];
	/** Absent from a flai older than S-0225. */
	strategic?: StrategicSpend;
};
export type TypeSpend = Spend & {
	models: ModelShare[];
	buckets: Bucket[];
	strategic?: StrategicSpend;
};
export const BUCKETS = ['hour', 'day', 'week'] as const;
export type BucketSize = (typeof BUCKETS)[number];
/** The longest window flai lays out by the hour, in days. */
export const MAX_HOUR_WINDOW_DAYS = 31;
export const TYPES = ['epic', 'story', 'task'] as const;
export type UsageReport = {
	items: number;
	tokens: number;
	cost: number;
	seconds: number;
	estimated?: boolean;
	models: ModelSpend[];
	/** Absent from a flai older than S-0163, with spend. */
	bucket?: BucketSize;
	spend?: Record<string, TypeSpend>;
	/** Absent from a flai older than S-0225. */
	strategic?: StrategicSpend & {
		estimated: boolean;
		kinds: (StrategicSpend & { kind: string })[];
	};
};
/**
 * Older flai builds emit null for empty lists; give every list the charts
 * iterate a value so a selection with no items draws an empty chart instead
 * of throwing (S-0045).
 */
export function normalise(r: Report): Report {
	const spend = r.usage?.spend
		? Object.fromEntries(
				Object.entries(r.usage.spend).map(([type, s]) => [
					type,
					{ ...s, models: s?.models ?? [], buckets: s?.buckets ?? [] }
				])
			)
		: undefined;
	return {
		...r,
		items: r.items ?? [],
		throughput: r.throughput ?? [],
		cfd: r.cfd ?? [],
		burnup: r.burnup ?? {},
		usage: {
			items: 0,
			tokens: 0,
			cost: 0,
			seconds: 0,
			...r.usage,
			models: r.usage?.models ?? [],
			spend
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
	/** Absent from a flai older than S-0143. */
	usage?: UsageReport;
	/** Absent from a flai older than S-0205. */
	forecasts?: Forecasts;
};

/** The charts of how work flows. */
export const FLOW_KINDS = ['cycle-time', 'burn-up', 'cfd', 'time-in-state', 'throughput'] as const;
/** The charts of what agents spent (S-0143, S-0163): they need items that carry usage. */
export const USAGE_KINDS = [
	'token-rate',
	'tokens-spent',
	'tokens-per-item',
	'tokens-per-dollar',
	'cost-spent',
	'cost-per-item',
	'cost',
	'time-per-model',
	'cost-per-model'
] as const;
/** The charts of forecasts and estimates against what happened (S-0212). */
export const PLANNING_KINDS = [
	'forecast-accuracy',
	'delivery-accuracy',
	'forecast-by-model'
] as const;
export type PlanningKind = (typeof PLANNING_KINDS)[number];
export const KINDS = [...FLOW_KINDS, ...USAGE_KINDS, ...PLANNING_KINDS] as const;
export type Kind = (typeof KINDS)[number];
export const TITLES: Record<Kind, string> = {
	'cycle-time': 'Cycle Time',
	'burn-up': 'Burn-up',
	cfd: 'Cumulative Flow',
	'time-in-state': 'Time in State',
	throughput: 'Throughput',
	'token-rate': 'Tokens / Min',
	'tokens-spent': 'Tokens / Day',
	'tokens-per-item': 'Tokens per item',
	'tokens-per-dollar': 'Tokens / $',
	'cost-spent': '$ / Day',
	'cost-per-item': '$ / Work Type',
	cost: '$ / Item',
	'time-per-model': 'Avg. Time / Model',
	'cost-per-model': 'Avg. Cost / Model',
	'forecast-accuracy': 'Forecast Accuracy',
	'delivery-accuracy': 'Delivery Accuracy',
	'forecast-by-model': 'Forecast Error / Model'
};
/** The charts drawn from spend over time, which flai lays out in buckets (S-0163). */
export const SPEND_KINDS: readonly Kind[] = [
	'token-rate',
	'tokens-spent',
	'tokens-per-item',
	'tokens-per-dollar',
	'cost-spent',
	'cost-per-item',
	'time-per-model',
	'cost-per-model'
];
/** The charts per item, one line per item type: every type is drawn, so none is chosen. */
export const PER_ITEM_KINDS: readonly Kind[] = ['tokens-per-item', 'cost-per-item'];
/** The windows the charts page offers. */
export const WINDOWS = ['1d', '7d', '30d', '90d', '365d'] as const;
/** The buckets a window can be laid out in: flai lays out by the hour over 31 days or less. */
export function bucketsFor(since: string): BucketSize[] {
	const days = since.endsWith('w') ? parseInt(since) * 7 : parseInt(since);
	return BUCKETS.filter((b) => b !== 'hour' || days <= MAX_HOUR_WINDOW_DAYS);
}
const isPlanning = (kind: Kind): kind is PlanningKind =>
	(PLANNING_KINDS as readonly Kind[]).includes(kind);
/**
 * The controls a chart uses. Spend over time is summed by flai, so no epic narrows it; a chart
 * per item shows every type, so none is chosen. The planning charts read the stories, narrowed by
 * nature and model; the chart per model shows every model, by the bucket.
 */
export function controls(kind: Kind) {
	if (isPlanning(kind))
		return {
			type: false,
			epic: false,
			bucket: kind === 'forecast-by-model',
			nature: true,
			model: kind !== 'forecast-by-model'
		};
	const spend = SPEND_KINDS.includes(kind);
	return {
		type: !PER_ITEM_KINDS.includes(kind),
		epic: !spend && !['cfd', 'throughput'].includes(kind),
		bucket: spend,
		nature: false,
		model: false
	};
}
/** A bucket's name in a title: Hour, Day, Week. */
const titled = (bucket: BucketSize) => bucket[0].toUpperCase() + bucket.slice(1);
/** A chart's title: the charts of what a bucket spent are named for the bucket (S-0169). */
export function titleOf(kind: Kind, bucket: BucketSize = 'day'): string {
	if (kind === 'tokens-spent') return `Tokens / ${titled(bucket)}`;
	if (kind === 'cost-spent') return `$ / ${titled(bucket)}`;
	return TITLES[kind];
}

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

/** The window a report covers, in milliseconds, from its start to its now; none when it says neither. */
function windowOf(r: Report): { start: number; end: number } | undefined {
	const start = Date.parse(r.window_start);
	const end = Date.parse(r.generated_at);
	return Number.isNaN(start) || Number.isNaN(end) ? undefined : { start, end };
}
/** Whether an item was completed in the report's window (S-0166): the items a chart per item plots. */
function doneInWindow(r: Report, i: ItemMetrics): boolean {
	if (!i.completed) return false;
	const w = windowOf(r);
	if (!w) return true;
	const at = Date.parse(i.completed);
	return at >= w.start && at <= w.end;
}
/** The items completed in the report's window. */
export const completedIn = (r: Report) => r.items.filter((i) => doneInWindow(r, i));
const DAY_MS = 86400e3;
/** The start of the hour, the UTC day, or the ISO week from its Monday that holds a moment. */
function floorTo(ms: number, bucket: BucketSize): number {
	if (bucket === 'hour') return Math.floor(ms / 3600e3) * 3600e3;
	const day = Math.floor(ms / DAY_MS) * DAY_MS;
	return bucket === 'day' ? day : day - ((new Date(day).getUTCDay() + 6) % 7) * DAY_MS;
}
/**
 * A time axis's ends: the report's window (S-0166), whatever the data it holds. A series by the day
 * starts at the day that holds the window's start, so that its first point is on the axis.
 */
function span(r: Report, daily = false): Opt {
	const w = windowOf(r);
	if (!w) return {};
	return { min: daily ? floorTo(w.start, 'day') : w.start, max: w.end };
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
/** A legend whose entries carry each series' own mark, where marks tell the series apart too. */
function marks(t: Theme, show: boolean): Opt {
	return { show, top: 0, textStyle: { color: t.textSecondary }, itemWidth: 14, itemHeight: 10 };
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
	const done = completedIn(r).filter(
		(i) => i.cycle_time_seconds !== undefined && (!epic || i.parent === epic)
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
		xAxis: axisX(t, { type: 'time', ...span(r) }),
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
		xAxis: axisX(t, { type: 'time', ...span(r, true) }),
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
		xAxis: axisX(t, { type: 'time', ...span(r, true) }),
		yAxis: axisY(t, { type: 'value', minInterval: 1 }),
		series
	});
}

/** The items of a day in time in state: the moment the UTC day starts, and those completed in it. */
type StateDay = { at: number; items: ItemMetrics[] };
type StateBar = { value: [number, number]; ids: string[] };
/** How many of a day's items the tooltip of time in state names before it counts the rest. */
const NAMED = 10;
/**
 * Time in state (S-0168): one stacked bar per UTC day of the window in which items were completed,
 * the mean hours per state of those items, on a time axis that spans the window.
 */
export function timeInState(r: Report, t: Theme, epic?: string): Opt {
	const byDay = new Map<number, ItemMetrics[]>();
	for (const i of completedIn(r).filter((i) => !epic || i.parent === epic)) {
		const at = floorTo(Date.parse(i.completed!), 'day');
		byDay.set(at, [...(byDay.get(at) ?? []), i]);
	}
	const daysDone: StateDay[] = [...byDay]
		.map(([at, items]) => ({ at, items }))
		.sort((a, b) => a.at - b.at);
	const mean = (items: ItemMetrics[], st: string) =>
		hours(items.reduce((n, i) => n + (i.time_in_state_seconds[st] ?? 0), 0) / items.length);
	const series = STATES.map((st) => ({
		name: st,
		type: 'bar',
		stack: 'state',
		barMaxWidth: 24,
		itemStyle: { color: colorFor(t, STATE_SLOT, st, 0), borderColor: t.surface, borderWidth: 1 },
		data: daysDone.map((d): StateBar => ({
			value: [d.at, mean(d.items, st)],
			ids: d.items.map((i) => i.id)
		}))
	}));
	return base(t, {
		useUTC: true,
		legend: legend(t, true),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'shadow' },
			formatter: (ps: { marker?: string; seriesName: string; data: StateBar }[]) => {
				const list = ps.filter((p) => p?.data);
				if (list.length === 0) return '';
				const ids = list[0].data.ids;
				const named =
					ids.slice(0, NAMED).join(', ') +
					(ids.length > NAMED ? ` and ${ids.length - NAMED} more` : '');
				const lines = list.map((p) => `${p.marker ?? ''}${p.seriesName}: ${p.data.value[1]} hours`);
				return `${new Date(list[0].data.value[0]).toISOString().slice(0, 10)}, mean of ${ids.length} ${ids.length === 1 ? r.type : plural(r.type)}<br/>${named}<br/>${lines.join('<br/>')}`;
			}
		}),
		xAxis: axisX(t, {
			type: 'time',
			minInterval: DAY_MS,
			...buckets(
				r,
				'day',
				daysDone.map((d) => d.at)
			)
		}),
		yAxis: axisY(t, {
			type: 'value',
			name: 'mean hours',
			nameTextStyle: { color: t.textSecondary }
		}),
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
/** The items completed in the window that carry usage, of an epic when one is chosen. */
export const withUsage = (r: Report, epic?: string) =>
	completedIn(r).filter(
		(i) => i.usage && i.usage.models.length > 0 && (!epic || i.parent === epic)
	);
/** The same, with the items only strategic agents spent on: the items $ / Item draws (ADR-0083). */
export const withCost = (r: Report, epic?: string) =>
	completedIn(r).filter(
		(i) =>
			i.usage &&
			(i.usage.models.length > 0 || (i.usage.strategic ?? []).length > 0) &&
			(!epic || i.parent === epic)
	);

/** The series flai laid out for a type; none from a flai older than S-0163. */
export function bucketsOf(r: Report, type: string = r.type): Bucket[] {
	return r.usage?.spend?.[type]?.buckets ?? [];
}
/** Whether the report carries spend over time at all. */
export const hasSpend = (r: Report) => r.usage?.spend !== undefined;
const bucketSize = (r: Report): BucketSize => r.usage?.bucket ?? 'day';
/** A bucket as a reader names it: the hour, the day, or the week from its Monday, in UTC. */
export function bucketLabel(at: string, bucket: BucketSize): string {
	const day = at.slice(0, 10);
	if (bucket === 'hour') return `${day} ${at.slice(11, 16)} UTC`;
	return bucket === 'week' ? `week of ${day}` : day;
}
const share = (b: Bucket, model: string) => (b.models ?? []).find((m) => m.model === model);
/** The models that spent anything in the series, in order of name. */
function modelsIn(buckets: Bucket[]): string[] {
	const names = new Set<string>();
	for (const b of buckets) for (const m of b.models ?? []) names.add(m.model);
	return [...names].sort();
}
const ALL = 'all models';
/**
 * The series of what strategic agents spent, on the cost charts only: apart from the models, and
 * estimated, since every charge is apportioned (S-0225, ADR-0083).
 */
const STRATEGIC = 'strategic';
const strategicColor = (t: Theme) => t.series[STRATEGIC_SLOT];
/** What strategic agents spent on an item, every kind summed, to four decimals as flai rounds it. */
export const strategicCost = (i: ItemMetrics) =>
	Math.round((i.usage?.strategic ?? []).reduce((n, s) => n + s.cost, 0) * 1e4) / 1e4;

/** A point's cost is estimated in part, or all of it is, as strategic spend is. */
type Point = { value: [string, number]; items: number; estimated?: boolean | 'all' };
type Hover = { marker?: string; seriesName: string; data: Point };
/** The tooltip of a chart over time: the bucket, then each series with its value and items. */
function hover(bucket: BucketSize, say: (v: number) => string, per: string) {
	return (ps: Hover | Hover[]) => {
		const list = (Array.isArray(ps) ? ps : [ps]).filter((p) => p?.data);
		if (list.length === 0) return '';
		const lines = list.map((p) => {
			const d = p.data;
			const of = per && d.items > 0 ? ` over ${d.items} ${d.items === 1 ? 'item' : 'items'}` : '';
			const est =
				d.estimated === 'all' ? ' (estimated)' : d.estimated ? ' (estimated in part)' : '';
			return `${p.marker ?? ''}${p.seriesName}: ${say(d.value[1])}${per}${of}${est}`;
		});
		return `${bucketLabel(list[0].data.value[0], bucket)}<br/>${lines.join('<br/>')}`;
	};
}
const BUCKET_MS: Record<BucketSize, number> = {
	hour: 3600e3,
	day: DAY_MS,
	week: 7 * DAY_MS
};
/**
 * The ends of a time axis of buckets: the report's window (S-0166), from the bucket that holds its
 * start to the one that holds its now, with half a bucket either side so that the mark of each is
 * whole. A report without a window runs from a bucket before the first drawn to one after the last.
 */
function buckets(r: Report, bucket: BucketSize, at: number[]): Opt {
	const size = BUCKET_MS[bucket];
	const w = windowOf(r);
	if (w)
		return { min: floorTo(w.start, bucket) - size / 2, max: floorTo(w.end, bucket) + size / 2 };
	return at.length > 0 ? { min: Math.min(...at) - size, max: Math.max(...at) + size } : {};
}
/** The axes of a chart of flai's spend buckets, in the report's bucket. */
function overTime(
	t: Theme,
	r: Report,
	series: { data: Point[] }[],
	name: string,
	say: (v: number) => string,
	per: string
) {
	const bucket = bucketSize(r);
	return bucketAxes(
		t,
		r,
		bucket,
		series.flatMap((s) => s.data.map((d) => Date.parse(d.value[0]))),
		series.length > 1,
		name,
		say,
		hover(bucket, say, per)
	);
}
/**
 * What the charts over time share: buckets are UTC, so the axis is; it spans the window as
 * `buckets` has it, or the moments `at` drawn without one, and its ticks are no finer than a
 * bucket.
 */
function bucketAxes(
	t: Theme,
	r: Report,
	bucket: BucketSize,
	at: number[],
	several: boolean,
	name: string,
	say: (v: number) => string,
	formatter: unknown
) {
	return {
		useUTC: true,
		// room for the legend above the axis name, which a legend of four would run into
		grid: { left: 64, right: 24, top: 64, bottom: 48, containLabel: false },
		legend: marks(t, several),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'line', lineStyle: { color: t.textSecondary, width: 1 } },
			formatter
		}),
		xAxis: axisX(t, { type: 'time', minInterval: BUCKET_MS[bucket], ...buckets(r, bucket, at) }),
		yAxis: axisY(t, {
			type: 'value',
			name,
			nameTextStyle: { color: t.textSecondary, align: 'left' },
			axisLabel: { color: t.textSecondary, formatter: say }
		})
	};
}
function line(
	t: Theme,
	name: string,
	color: string,
	symbol: string,
	data: Point[],
	dashed = false
) {
	return {
		name,
		type: 'line',
		showSymbol: data.length < 40,
		symbol,
		symbolSize: 9,
		lineStyle: { width: 2, color, type: dashed ? 'dashed' : 'solid' },
		itemStyle: { color, borderColor: t.surface, borderWidth: 2 },
		data
	};
}

/**
 * One line per model over time, of a value each bucket has per model; with two models or more,
 * a dashed line for all of them together.
 */
function perModel(
	r: Report,
	t: Theme,
	pick: (s: Spend) => number | undefined,
	name: string,
	say: (v: number) => string,
	per: string
): Opt {
	const buckets = bucketsOf(r);
	const names = modelsIn(buckets);
	const points = (of: (b: Bucket) => Spend | undefined): Point[] =>
		buckets.flatMap((b) => {
			const s = of(b);
			const v = s && pick(s);
			return s && v !== undefined
				? [{ value: [b.at, v] as [string, number], items: s.items, estimated: s.estimated }]
				: [];
		});
	const series = names.map((m) =>
		line(
			t,
			m,
			modelColor(t, m),
			modelSymbol(m),
			points((b) => share(b, m))
		)
	);
	if (names.length > 1)
		series.push(
			line(
				t,
				ALL,
				t.textSecondary,
				'emptyCircle',
				points((b) => b),
				true
			)
		);
	return base(t, { ...overTime(t, r, series, name, say, per), series });
}

/**
 * Token rate: per model, the tokens per minute of agent work over the items done in each bucket
 * (S-0163: an hour of agent work is longer than an item takes).
 */
export const tokenRate = (r: Report, t: Theme) =>
	perModel(
		r,
		t,
		(s) => s.tokens_per_minute,
		'tokens per agent minute',
		count,
		' tokens per minute'
	);

/** Tokens per dollar: per model, the tokens a dollar bought over the items done in each bucket. */
export const tokensPerDollar = (r: Report, t: Theme) =>
	perModel(r, t, (s) => s.tokens_per_dollar, 'tokens per US dollar', count, ' tokens per dollar');

/**
 * What each bucket spent, stacked by model, with the running mean per bucket as a line: tokens
 * per day, or cost per day. Cost per day stacks what strategic agents spent on top, as a series
 * of its own, when any bucket has it; the running mean stays the agents' (S-0225).
 */
function spentOverTime(r: Report, t: Theme, what: 'tokens' | 'cost'): Opt {
	const buckets = bucketsOf(r);
	const names = modelsIn(buckets);
	const bucket = bucketSize(r);
	const say = what === 'tokens' ? count : dollars;
	const bar = (name: string, color: string, data: Point[]) => ({
		name,
		type: 'bar',
		stack: 'spent',
		barMaxWidth: 24,
		itemStyle: { color, borderColor: t.surface, borderWidth: 1 },
		data
	});
	const bars = names.map((m) =>
		bar(
			m,
			modelColor(t, m),
			buckets.map((b): Point => {
				const s = share(b, m);
				return { value: [b.at, s?.[what] ?? 0], items: s?.items ?? 0, estimated: s?.estimated };
			})
		)
	);
	if (what === 'cost' && buckets.some((b) => (b.strategic?.cost ?? 0) > 0))
		bars.push(
			bar(
				STRATEGIC,
				strategicColor(t),
				buckets.map((b): Point => {
					const c = b.strategic?.cost ?? 0;
					return {
						value: [b.at, c],
						items: b.strategic?.items ?? 0,
						estimated: c > 0 ? 'all' : undefined
					};
				})
			)
		);
	const mean = line(
		t,
		`mean per ${bucket}`,
		t.textSecondary,
		'none',
		buckets.map((b) => ({
			value: [b.at, what === 'tokens' ? b.mean_tokens : b.mean_cost],
			items: 0
		})),
		true
	);
	const series = buckets.length > 0 ? [...bars, { ...mean, showSymbol: false }] : [];
	return base(t, {
		...overTime(t, r, series, what === 'tokens' ? 'tokens' : 'US dollars', say, ''),
		series
	});
}
export const tokensSpent = (r: Report, t: Theme) => spentOverTime(r, t, 'tokens');
export const costSpent = (r: Report, t: Theme) => spentOverTime(r, t, 'cost');

/** What an item took on average in each bucket: one line per item type. */
function perItem(r: Report, t: Theme, what: 'tokens' | 'cost'): Opt {
	const pick = (s: Spend) => (what === 'tokens' ? s.tokens_per_item : s.cost_per_item);
	const series = TYPES.filter((type) => bucketsOf(r, type).some((b) => b.items > 0)).map((type) =>
		line(
			t,
			type,
			colorFor(t, TYPE_SLOT, type, 0),
			TYPE_SYMBOL[type],
			bucketsOf(r, type).flatMap((b): Point[] => {
				const v = pick(b);
				return v !== undefined
					? [{ value: [b.at, v], items: b.items, estimated: b.estimated }]
					: [];
			})
		)
	);
	return base(t, {
		...overTime(
			t,
			r,
			series,
			what === 'tokens' ? 'tokens per item' : 'US dollars per item',
			what === 'tokens' ? count : dollars,
			what === 'tokens' ? ' tokens each' : ' each'
		),
		series
	});
}
export const tokensPerItem = (r: Report, t: Theme) => perItem(r, t, 'tokens');
export const costPerItem = (r: Report, t: Theme) => perItem(r, t, 'cost');
/**
 * Avg. Time / Model (S-0169): per model, the mean agent time an item of the report's type took,
 * over the items done in each bucket that the model worked on.
 */
export const timePerModel = (r: Report, t: Theme) =>
	perModel(r, t, minutesPerItem, `agent minutes per ${r.type}`, (m) => human(m * 60), ' each');
/**
 * Avg. Cost / Model (S-0169): per model, the mean dollars an item of the report's type took over
 * the items done in each bucket that the model worked on.
 */
export const costPerModel = (r: Report, t: Theme) =>
	perModel(r, t, (s) => s.cost_per_item, `US dollars per ${r.type}`, dollars, ' each');

/** One row of a spend chart's table: a bucket, and whose spend in it the row is. */
export type SpendRow = Spend & { at: string; of: string; mean_tokens?: number; mean_cost?: number };
/**
 * What a chart over time plots, as rows: per item, each bucket of each item type in which items
 * were done; otherwise each such bucket of the report's type, whole, with its running means,
 * and then per model, and for cost per bucket what strategic agents spent in it (S-0225), which
 * may be all a bucket holds. Newest first.
 */
export function spendRows(r: Report, kind: Kind): SpendRow[] {
	const rows: SpendRow[] = [];
	const put = (b: Bucket, of: string, s: Spend, means: boolean) =>
		rows.push({
			...s,
			at: b.at,
			of,
			mean_tokens: means ? b.mean_tokens : undefined,
			mean_cost: means ? b.mean_cost : undefined
		});
	if (PER_ITEM_KINDS.includes(kind)) {
		for (const type of TYPES)
			for (const b of bucketsOf(r, type)) if (b.items > 0) put(b, type, b, false);
	} else {
		for (const b of bucketsOf(r)) {
			const { models, strategic, ...whole } = b;
			if (b.items > 0) {
				put(b, models && models.length > 1 ? ALL : r.type, whole, true);
				for (const m of models ?? []) put(b, m.model, m, false);
			}
			if (kind === 'cost-spent' && strategic && strategic.items > 0)
				put(b, STRATEGIC, { ...strategic, estimated: true }, false);
		}
	}
	// newest bucket first; within a bucket, the order the rows were put in
	const seq = new Map(rows.map((row, i) => [row, i]));
	return rows.sort((a, b) => (a.at !== b.at ? (a.at < b.at ? 1 : -1) : seq.get(a)! - seq.get(b)!));
}

/**
 * Cost: one bar per completed item with usage, in order of completion, stacked by model, with
 * what strategic agents spent on it on top as a series of its own (S-0225); an item on which only
 * they spent has that alone. An item whose cost is estimated in part, or that carries strategic
 * spend, which is estimated, carries an asterisk on its label and says so in the tooltip, where
 * the total is the agents' and the strategic spend stays apart from it.
 */
export function cost(r: Report, t: Theme, epic?: string): Opt {
	const all = models(r);
	const done = withCost(r, epic).sort((a, b) =>
		a.completed! < b.completed! ? -1 : a.completed! > b.completed! ? 1 : 0
	);
	const present = all.filter((name) =>
		done.some((i) => i.usage!.models.some((m) => m.model === name))
	);
	const bar = (name: string, color: string, data: number[]) => ({
		name,
		type: 'bar',
		stack: 'cost',
		barMaxWidth: 24,
		itemStyle: { color, borderColor: t.surface, borderWidth: 1 },
		data
	});
	const series = present.map((name) =>
		bar(
			name,
			modelColor(t, name),
			done.map((i) => i.usage!.models.find((m) => m.model === name)?.cost ?? 0)
		)
	);
	const strategic = done.map(strategicCost);
	if (strategic.some((c) => c > 0)) series.push(bar(STRATEGIC, strategicColor(t), strategic));
	const planned = (i: ItemMetrics) => (i.usage!.strategic ?? []).some((s) => s.estimated);
	return base(t, {
		legend: legend(t, series.length > 1),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'shadow' },
			formatter: (ps: { dataIndex: number; seriesName: string; value: number }[]) => {
				const i = done[ps[0]?.dataIndex ?? 0];
				if (!i) return '';
				const lines = ps
					.filter((p) => p.value > 0 && p.seriesName !== STRATEGIC)
					.map((p) => `${p.seriesName}: ${dollars(p.value)}`);
				const agents =
					i.usage!.models.length > 0
						? `${lines.join('<br/>')}<br/>total ${dollars(i.usage!.cost)}${i.usage!.estimated ? ' (estimated in part)' : ''}`
						: 'no agent spend';
				const s = strategicCost(i);
				const apart =
					s > 0 ? `<br/>${STRATEGIC}: ${dollars(s)}${planned(i) ? ' (estimated)' : ''}` : '';
				return `${i.id} ${i.title}<br/>${agents}${apart}`;
			}
		}),
		xAxis: axisX(t, {
			type: 'category',
			data: done.map((i) => (i.usage!.estimated || planned(i) ? `${i.id}*` : i.id)),
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

/** The model of an item whose agent names none, and of every item from a flai before ADR-0111. */
export const NO_MODEL = '(none)';
/** The model an item's errors are grouped under (ADR-0111). */
export const modelOf = (i: ItemMetrics) => i.model ?? NO_MODEL;
/** The nature and the model a planning chart is narrowed to; neither when absent. */
export type ErrorFilter = { nature?: string; model?: string };
/** The errors each planning chart plots. */
export const ERRORS_OF: Record<PlanningKind, readonly ErrorField[]> = {
	'forecast-accuracy': ['forecast_error_seconds', 'estimate_error_seconds'],
	'delivery-accuracy': ['delivery_error_seconds'],
	'forecast-by-model': ['forecast_error_seconds']
};
/** The items done, not cancelled, in the report's window: those flai spreads the errors over. */
export const doneIn = (r: Report) => completedIn(r).filter((i) => i.status === 'done');
/** Whether an item is of the filter's nature and model. */
const passes = (i: ItemMetrics, f: ErrorFilter) =>
	(!f.nature || i.nature === f.nature) && (!f.model || modelOf(i) === f.model);
/**
 * The natures and the models of the items done in the window that carry an error a planning chart
 * plots, in order of name: what its nature and model selects offer.
 */
export function errorFacets(
	report: Report,
	kind: PlanningKind
): { natures: string[]; models: string[] } {
	const items = doneIn(normalise(report)).filter((i) =>
		ERRORS_OF[kind].some((f) => i[f] !== undefined)
	);
	return {
		natures: [...new Set(items.map((i) => i.nature))].sort(),
		models: [...new Set(items.map(modelOf))].sort()
	};
}
/** The nearest-rank percentile of values sorted ascending, as flai works it out (rank at least 1). */
export function percentile(sorted: readonly number[], p: number): number {
	return sorted[Math.max(1, Math.ceil((p / 100) * sorted.length)) - 1];
}
/** The count and the p50 and p85 of the absolute values of errors, as flai spreads them. */
export function spreadOf(errors: readonly number[]): ErrorSpread {
	if (errors.length === 0) return { count: 0 };
	const sorted = errors.map(Math.abs).sort((a, b) => a - b);
	return {
		count: sorted.length,
		p50_seconds: percentile(sorted, 50),
		p85_seconds: percentile(sorted, 85)
	};
}
/**
 * The spread of an error under a filter: flai's own, in all, per nature, or per model, so that the
 * figures match `flai stats --json`; under both, worked out from the errors shown.
 */
export function spreadFor(
	r: Report,
	which: keyof Forecasts,
	f: ErrorFilter,
	shown: readonly number[]
): ErrorSpread | undefined {
	const s = r.forecasts?.[which];
	if (f.nature && f.model) return spreadOf(shown);
	if (f.nature) return s?.by_nature?.[f.nature];
	if (f.model) return s?.by_model?.[f.model];
	return s;
}
/** A signed duration for axes and tooltips: later or longer than forecast is positive. */
export const humanSigned = (seconds: number) =>
	seconds < 0 ? `-${human(-seconds)}` : `+${human(seconds)}`;
type ErrorPoint = { value: [string, number]; id: string; title: string };
/** One point per item that carries the error: x completed, y the error in seconds. */
const errorPoints = (items: ItemMetrics[], field: ErrorField): ErrorPoint[] =>
	items.flatMap((i): ErrorPoint[] => {
		const e = i[field];
		return e === undefined ? [] : [{ value: [i.completed!, e], id: i.id, title: i.title }];
	});
/** The p50 and p85 of a spread either side of zero, in seconds over `unit`: the reference lines. */
function band(s: ErrorSpread | undefined, unit = 1): Opt[] {
	const at: [string, number | undefined][] = [
		['p50', s?.p50_seconds],
		['p85', s?.p85_seconds]
	];
	return at.flatMap(([name, v]) =>
		v === undefined ? [] : [v / unit, -v / unit].map((yAxis) => ({ yAxis, name }))
	);
}
/**
 * Forecast accuracy (S-0212): one point per story done in the window with a forecast, x completed,
 * y the forecast error; the estimate error as a second series. The p50 and p85 of the absolute
 * forecast error are drawn either side of zero, the band in which half and 85% of the errors fall.
 */
export function forecastAccuracy(r: Report, t: Theme, f: ErrorFilter = {}): Opt {
	const items = doneIn(r).filter((i) => passes(i, f));
	const forecast = errorPoints(items, 'forecast_error_seconds');
	const estimate = errorPoints(items, 'estimate_error_seconds');
	const s = spreadFor(r, 'forecast', f, forecast.map((p) => p.value[1]));
	const lines = band(s);
	const scatter = (name: string, color: string, symbol: string, data: ErrorPoint[]) => ({
		name,
		type: 'scatter',
		symbol,
		symbolSize: 10,
		itemStyle: { color, borderColor: t.surface, borderWidth: 2 },
		data
	});
	const series = [
		{
			...scatter('forecast error', t.series[0], 'circle', forecast),
			markLine: lines.length > 0 ? refLine(t, lines, '{b}') : undefined
		},
		scatter('estimate error', t.series[3], 'triangle', estimate)
	].filter((x) => x.data.length > 0);
	return base(t, {
		legend: marks(t, series.length > 1),
		tooltip: tooltip(t, {
			formatter: (p: { seriesName: string; data: ErrorPoint }) =>
				`${p.data.id} ${p.data.title}<br/>${p.seriesName}: ${humanSigned(p.data.value[1])} · ${p.data.value[0].slice(0, 10)}`
		}),
		xAxis: axisX(t, { type: 'time', ...span(r) }),
		yAxis: axisY(t, {
			type: 'value',
			name: 'actual minus forecast',
			nameTextStyle: { color: t.textSecondary, align: 'left' },
			axisLabel: { color: t.textSecondary, formatter: humanSigned }
		}),
		series
	});
}
const DAY_S = 86400;
/**
 * One ISO week of the delivery-accuracy chart: the stories with a delivery forecast done in it,
 * how many of them on or before the forecast date, and that as a share, absent when none was.
 */
export type OnTimeWeek = {
	week: string;
	start: string;
	count: number;
	on_time: number;
	share?: number;
};
/** The ISO week, as flai names it (2026-W32), that starts on a Monday at a moment in UTC. */
function isoWeek(monday: number): string {
	const thursday = monday + 3 * DAY_MS;
	const year = new Date(thursday).getUTCFullYear();
	const week = Math.floor((thursday - Date.UTC(year, 0, 1)) / BUCKET_MS.week) + 1;
	return `${year}-W${String(week).padStart(2, '0')}`;
}
/**
 * The share of the stories done in the window with a delivery forecast, of the filter's nature and
 * model, delivered on or before the forecast date, per ISO week from the one that holds the
 * window's start to the one that holds now (ADR-0054); none from a report without a window.
 */
export function onTimeShare(report: Report, f: ErrorFilter = {}): OnTimeWeek[] {
	const r = normalise(report);
	const w = windowOf(r);
	if (!w) return [];
	const weeks = new Map<number, OnTimeWeek>();
	for (let m = floorTo(w.start, 'week'); m <= w.end; m += BUCKET_MS.week)
		weeks.set(m, {
			week: isoWeek(m),
			start: new Date(m).toISOString().slice(0, 10),
			count: 0,
			on_time: 0
		});
	for (const i of doneIn(r)) {
		const e = i.delivery_error_seconds;
		const week = weeks.get(floorTo(Date.parse(i.completed!), 'week'));
		if (e === undefined || !week || !passes(i, f)) continue;
		week.count += 1;
		if (e <= 0) week.on_time += 1;
	}
	return [...weeks.values()].map((wk) =>
		wk.count > 0 ? { ...wk, share: wk.on_time / wk.count } : wk
	);
}
/** A share as a whole percentage. */
const percent = (v: number) => `${Math.round(v * 100)}%`;
/** A week's share on the time axis; null, a gap in the line, for a week without any. */
type SharePoint = OnTimeWeek & { value: [number, number | null] };
/**
 * Delivery accuracy (S-0212): one point per story done in the window with a delivery forecast,
 * x completed, y the days it was delivered after the forecast date, early below zero, with the p50
 * and p85 of the absolute delivery error either side of zero. On a second axis, the share of those
 * stories delivered on or before the forecast date per ISO week, drawn at the middle of the part of
 * the week in the window; a week without any is a gap.
 */
export function deliveryAccuracy(r: Report, t: Theme, f: ErrorFilter = {}): Opt {
	const items = doneIn(r).filter((i) => passes(i, f));
	const errors = errorPoints(items, 'delivery_error_seconds');
	const late = errors.map((p): ErrorPoint => ({
		...p,
		value: [p.value[0], p.value[1] / DAY_S]
	}));
	const s = spreadFor(r, 'delivery', f, errors.map((p) => p.value[1]));
	const lines = band(s, DAY_S);
	const w = windowOf(r);
	const weeks = w
		? onTimeShare(r, f).map((wk): SharePoint => {
				const from = Date.parse(wk.start);
				const mid = (Math.max(from, w.start) + Math.min(from + BUCKET_MS.week, w.end)) / 2;
				return { ...wk, value: [mid, wk.share ?? null] };
			})
		: [];
	const shares = weeks.some((wk) => wk.share !== undefined) ? weeks : [];
	const series = [
		{
			name: 'delivery error',
			type: 'scatter',
			symbolSize: 10,
			itemStyle: { color: t.series[0], borderColor: t.surface, borderWidth: 2 },
			data: late,
			markLine: lines.length > 0 ? refLine(t, lines, '{b}') : undefined
		},
		{
			name: 'on time per week',
			type: 'line',
			yAxisIndex: 1,
			connectNulls: false,
			symbol: 'emptyCircle',
			symbolSize: 8,
			lineStyle: { width: 2, color: t.series[2] },
			itemStyle: { color: t.series[2] },
			data: shares
		}
	].filter((x) => x.data.length > 0);
	const name = { color: t.textSecondary, align: 'left' };
	return base(t, {
		// room for the legend above the axis names, and for the share's labels on the right
		grid: { left: 56, right: 56, top: 64, bottom: 48, containLabel: false },
		legend: marks(t, series.length > 1),
		tooltip: tooltip(t, {
			formatter: (p: { seriesName: string; data: ErrorPoint | SharePoint }) => {
				const d = p.data;
				if (!('week' in d))
					return `${d.id} ${d.title}<br/>${p.seriesName}: ${humanSigned(d.value[1] * DAY_S)} · ${d.value[0].slice(0, 10)}`;
				const of =
					d.share === undefined
						? 'no story with a delivery forecast'
						: `${d.on_time} of ${d.count} on or before the forecast date (${percent(d.share)})`;
				return `week of ${d.start} (${d.week})<br/>${of}`;
			}
		}),
		xAxis: axisX(t, { type: 'time', ...span(r) }),
		yAxis: [
			axisY(t, {
				type: 'value',
				name: 'days after forecast date',
				nameTextStyle: name,
				axisLabel: { color: t.textSecondary, formatter: (v: number) => humanSigned(v * DAY_S) }
			}),
			axisY(t, {
				type: 'value',
				name: 'on time',
				min: 0,
				max: 1,
				nameTextStyle: { ...name, align: 'right' },
				splitLine: { show: false },
				axisLabel: { color: t.textSecondary, formatter: percent }
			})
		],
		series
	});
}
/**
 * A model's p50 absolute forecast error in a bucket, at its start, over the stories done in it;
 * null, a gap in the line, for a bucket without any.
 */
type ModelBucket = { value: [number, number | null]; count: number };
/** The starts of the buckets from the one holding the window's start to the one holding now. */
function bucketStarts(w: { start: number; end: number }, bucket: BucketSize): number[] {
	const at: number[] = [];
	for (let b = floorTo(w.start, bucket); b <= w.end; b += BUCKET_MS[bucket]) at.push(b);
	return at;
}
/**
 * Forecast error per model (S-0212): one line per model, with a point per bucket of the report's
 * size in which stories of that model with a forecast were done in the window, at the p50 of their
 * absolute forecast error as flai works it out; a bucket without any is a gap, not 0. Narrowed by
 * nature only: the chart shows every model, each in its fixed colour and mark.
 */
export function forecastByModel(r: Report, t: Theme, f: ErrorFilter = {}): Opt {
	const bucket = bucketSize(r);
	const byModel = new Map<string, Map<number, number[]>>();
	for (const i of doneIn(r)) {
		const e = i.forecast_error_seconds;
		if (e === undefined || !passes(i, { nature: f.nature })) continue;
		const at = floorTo(Date.parse(i.completed!), bucket);
		const of = byModel.get(modelOf(i)) ?? new Map<number, number[]>();
		of.set(at, [...(of.get(at) ?? []), e]);
		byModel.set(modelOf(i), of);
	}
	const w = windowOf(r);
	const drawn = [...byModel.values()].flatMap((of) => [...of.keys()]);
	const starts = w ? bucketStarts(w, bucket) : [...new Set(drawn)].sort((a, b) => a - b);
	const series = [...byModel.keys()].sort().map((m) => {
		const of = byModel.get(m)!;
		const data = starts.map((at): ModelBucket => {
			const errors = of.get(at) ?? [];
			return { value: [at, spreadOf(errors).p50_seconds ?? null], count: errors.length };
		});
		// a point between two gaps is drawn by its mark alone, so every mark shows
		return { ...line(t, m, modelColor(t, m), modelSymbol(m), []), showSymbol: true, data };
	});
	const noun = (n: number) => (n === 1 ? r.type : plural(r.type));
	const tip = (ps: { marker?: string; seriesName: string; data: ModelBucket }[]) => {
		const list = ps.filter((p) => p?.data && p.data.value[1] !== null);
		if (list.length === 0) return '';
		const lines = list.map((p) => {
			const n = p.data.count;
			return `${p.marker ?? ''}${p.seriesName}: p50 ${human(p.data.value[1]!)} over ${n} ${noun(n)}`;
		});
		const at = new Date(list[0].data.value[0]).toISOString();
		return `${bucketLabel(at, bucket)}<br/>${lines.join('<br/>')}`;
	};
	return base(t, {
		...bucketAxes(
			t,
			r,
			bucket,
			starts,
			series.length > 1,
			'p50 absolute forecast error',
			human,
			tip
		),
		series
	});
}
/** One story of a planning chart's table: its forecast, its cycle time, and its errors. */
export type ForecastRow = {
	id: string;
	title: string;
	completed: string;
	nature: string;
	model: string;
	forecast_seconds?: number;
	cycle_time_seconds?: number;
	forecast_error_seconds?: number;
	delivery_error_seconds?: number;
	estimate_error_seconds?: number;
};
const ERROR_FIELDS: readonly ErrorField[] = [
	'forecast_error_seconds',
	'delivery_error_seconds',
	'estimate_error_seconds'
];
/**
 * The rows of a planning chart's table: the stories done in the window with a forecast, delivery,
 * or estimate error, of the filter's nature and model, newest first.
 */
export function forecastRows(report: Report, f: ErrorFilter = {}): ForecastRow[] {
	return doneIn(normalise(report))
		.filter((i) => passes(i, f) && ERROR_FIELDS.some((e) => i[e] !== undefined))
		.sort((a, b) => Date.parse(b.completed!) - Date.parse(a.completed!))
		.map((i) => ({
			id: i.id,
			title: i.title,
			completed: i.completed!,
			nature: i.nature,
			model: modelOf(i),
			forecast_seconds: i.forecast_seconds,
			cycle_time_seconds: i.cycle_time_seconds,
			forecast_error_seconds: i.forecast_error_seconds,
			delivery_error_seconds: i.delivery_error_seconds,
			estimate_error_seconds: i.estimate_error_seconds
		}));
}

/** A chart's ECharts option: the planning charts narrowed by the filter, the others by the epic. */
export function build(
	kind: Kind,
	report: Report,
	t: Theme,
	epic?: string,
	filter: ErrorFilter = {}
): Opt {
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
		case 'token-rate':
			return tokenRate(r, t);
		case 'tokens-spent':
			return tokensSpent(r, t);
		case 'tokens-per-item':
			return tokensPerItem(r, t);
		case 'tokens-per-dollar':
			return tokensPerDollar(r, t);
		case 'cost-spent':
			return costSpent(r, t);
		case 'cost-per-item':
			return costPerItem(r, t);
		case 'cost':
			return cost(r, t, epic);
		case 'time-per-model':
			return timePerModel(r, t);
		case 'cost-per-model':
			return costPerModel(r, t);
		case 'forecast-accuracy':
			return forecastAccuracy(r, t, filter);
		case 'delivery-accuracy':
			return deliveryAccuracy(r, t, filter);
		case 'forecast-by-model':
			return forecastByModel(r, t, filter);
	}
}
