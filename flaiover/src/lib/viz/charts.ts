// Chart option builders: pure functions from a flai stats report to an
// ECharts option, so they are unit-testable without a DOM. One y-axis per
// chart, thin marks, legends for two or more series, tooltips everywhere.
import {
	colorFor,
	modelSlot,
	modelSymbol,
	NATURE_SLOT,
	STATE_SLOT,
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
	age_seconds?: number;
	usage?: ItemUsage;
};
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
};
export type SpendPoint = { at: string; id: string; done: number; tokens: number; cost: number };
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
	tokens_per_minute?: number;
	tokens_per_dollar?: number;
};
export type ModelShare = Spend & { model: string };
/** The items done in one bucket of time, which starts at `at`, with the running means so far. */
export type Bucket = Spend & {
	at: string;
	mean_tokens: number;
	mean_cost: number;
	models?: ModelShare[];
};
export type TypeSpend = Spend & { models: ModelShare[]; buckets: Bucket[] };
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
	done: SpendPoint[];
	by_model: Record<string, SpendPoint[]>;
	/** Absent from a flai older than S-0163, with spend. */
	bucket?: BucketSize;
	spend?: Record<string, TypeSpend>;
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
			by_model: r.usage?.by_model ?? {},
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

/** The charts of how work flows. */
export const FLOW_KINDS = [
	'cycle-time',
	'burn-up',
	'cfd',
	'time-in-state',
	'throughput',
	'aging',
	'estimates'
] as const;
/** The charts of what agents spent (S-0143, S-0163): they need items that carry usage. */
export const USAGE_KINDS = [
	'token-rate',
	'tokens-spent',
	'tokens-per-item',
	'tokens-per-dollar',
	'cost-spent',
	'cost-per-item',
	'cost',
	'completion-time',
	'completion-cost'
] as const;
export const KINDS = [...FLOW_KINDS, ...USAGE_KINDS] as const;
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
	'tokens-spent': 'Tokens per day',
	'tokens-per-item': 'Tokens per item',
	'tokens-per-dollar': 'Tokens per dollar',
	'cost-spent': 'Cost per day',
	'cost-per-item': 'Cost per item',
	cost: 'Cost by item',
	'completion-time': 'Completion over time',
	'completion-cost': 'Completion against cost'
};
/** The charts drawn from spend over time, which flai lays out in buckets (S-0163). */
export const SPEND_KINDS: readonly Kind[] = [
	'token-rate',
	'tokens-spent',
	'tokens-per-item',
	'tokens-per-dollar',
	'cost-spent',
	'cost-per-item'
];
/** The charts per item, which compare the item types, or the models on one type. */
export const PER_ITEM_KINDS: readonly Kind[] = ['tokens-per-item', 'cost-per-item'];
export type By = 'type' | 'model';
/** The windows the charts page offers. */
export const WINDOWS = ['1d', '7d', '30d', '90d', '365d'] as const;
/** The buckets a window can be laid out in: flai lays out by the hour over 31 days or less. */
export function bucketsFor(since: string): BucketSize[] {
	const days = since.endsWith('w') ? parseInt(since) * 7 : parseInt(since);
	return BUCKETS.filter((b) => b !== 'hour' || days <= MAX_HOUR_WINDOW_DAYS);
}
/**
 * The controls a chart uses. Spend over time is summed by flai, so no epic narrows it; a chart
 * per item by type shows every type, so none is chosen.
 */
export function controls(kind: Kind, by: By) {
	const spend = SPEND_KINDS.includes(kind);
	const perItem = PER_ITEM_KINDS.includes(kind);
	return {
		type: !(perItem && by === 'type'),
		epic:
			!spend &&
			!['cfd', 'throughput', 'estimates', 'completion-time', 'completion-cost'].includes(kind),
		bucket: spend,
		by: perItem
	};
}
/** A chart's title: the charts of what a bucket spent are named for the bucket. */
export function titleOf(kind: Kind, bucket: BucketSize = 'day'): string {
	if (kind === 'tokens-spent') return `Tokens per ${bucket}`;
	if (kind === 'cost-spent') return `Cost per ${bucket}`;
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

type Point = { value: [string, number]; items: number; estimated?: boolean };
type Hover = { marker?: string; seriesName: string; data: Point };
/** The tooltip of a chart over time: the bucket, then each series with its value and items. */
function hover(bucket: BucketSize, say: (v: number) => string, per: string) {
	return (ps: Hover | Hover[]) => {
		const list = (Array.isArray(ps) ? ps : [ps]).filter((p) => p?.data);
		if (list.length === 0) return '';
		const lines = list.map((p) => {
			const d = p.data;
			const of = per && d.items > 0 ? ` over ${d.items} ${d.items === 1 ? 'item' : 'items'}` : '';
			return `${p.marker ?? ''}${p.seriesName}: ${say(d.value[1])}${per}${of}${d.estimated ? ' (estimated in part)' : ''}`;
		});
		return `${bucketLabel(list[0].data.value[0], bucket)}<br/>${lines.join('<br/>')}`;
	};
}
const BUCKET_MS: Record<BucketSize, number> = {
	hour: 3600e3,
	day: 86400e3,
	week: 7 * 86400e3
};
/**
 * What the charts over time share: buckets are UTC, so the axis is; it runs from a bucket before
 * the first drawn to one after the last, so that a single bucket is not a mark alone on an axis
 * of its own making, and its ticks are no finer than a bucket.
 */
function overTime(
	t: Theme,
	r: Report,
	series: { data: Point[] }[],
	name: string,
	say: (v: number) => string,
	per: string
) {
	const size = BUCKET_MS[bucketSize(r)];
	const at = series.flatMap((s) => s.data.map((d) => Date.parse(d.value[0])));
	const range = at.length > 0 ? { min: Math.min(...at) - size, max: Math.max(...at) + size } : {};
	return {
		useUTC: true,
		// room for the legend above the axis name, which a legend of four would run into
		grid: { left: 64, right: 24, top: 64, bottom: 48, containLabel: false },
		legend: marks(t, series.length > 1),
		tooltip: tooltip(t, {
			trigger: 'axis',
			axisPointer: { type: 'line', lineStyle: { color: t.textSecondary, width: 1 } },
			formatter: hover(bucketSize(r), say, per)
		}),
		xAxis: axisX(t, { type: 'time', minInterval: size, ...range }),
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
 * per day, or cost per day.
 */
function spentOverTime(r: Report, t: Theme, what: 'tokens' | 'cost'): Opt {
	const buckets = bucketsOf(r);
	const names = modelsIn(buckets);
	const bucket = bucketSize(r);
	const say = what === 'tokens' ? count : dollars;
	const bars = names.map((m) => ({
		name: m,
		type: 'bar',
		stack: 'spent',
		barMaxWidth: 24,
		itemStyle: { color: modelColor(t, m), borderColor: t.surface, borderWidth: 1 },
		data: buckets.map((b): Point => {
			const s = share(b, m);
			return { value: [b.at, s?.[what] ?? 0], items: s?.items ?? 0, estimated: s?.estimated };
		})
	}));
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

/**
 * What an item took on average in each bucket: one line per item type, or, by model, one line
 * per model over the items of the report's type that it worked on.
 */
function perItem(r: Report, t: Theme, what: 'tokens' | 'cost', by: By): Opt {
	const pick = (s: Spend) => (what === 'tokens' ? s.tokens_per_item : s.cost_per_item);
	const points = (buckets: Bucket[], of: (b: Bucket) => Spend | undefined): Point[] =>
		buckets.flatMap((b) => {
			const s = of(b);
			const v = s && pick(s);
			return s && v !== undefined
				? [{ value: [b.at, v] as [string, number], items: s.items, estimated: s.estimated }]
				: [];
		});
	const series =
		by === 'model'
			? modelsIn(bucketsOf(r)).map((m) =>
					line(
						t,
						m,
						modelColor(t, m),
						modelSymbol(m),
						points(bucketsOf(r), (b) => share(b, m))
					)
				)
			: TYPES.filter((type) => bucketsOf(r, type).some((b) => b.items > 0)).map((type) =>
					line(
						t,
						type,
						colorFor(t, TYPE_SLOT, type, 0),
						TYPE_SYMBOL[type],
						points(bucketsOf(r, type), (b) => b)
					)
				);
	const unit = by === 'model' ? r.type : 'item';
	return base(t, {
		...overTime(
			t,
			r,
			series,
			what === 'tokens' ? `tokens per ${unit}` : `US dollars per ${unit}`,
			what === 'tokens' ? count : dollars,
			what === 'tokens' ? ' tokens each' : ' each'
		),
		series
	});
}
export const tokensPerItem = (r: Report, t: Theme, by: By = 'type') => perItem(r, t, 'tokens', by);
export const costPerItem = (r: Report, t: Theme, by: By = 'type') => perItem(r, t, 'cost', by);

/** One row of a spend chart's table: a bucket, and whose spend in it the row is. */
export type SpendRow = Spend & { at: string; of: string; mean_tokens?: number; mean_cost?: number };
/**
 * What a chart over time plots, as rows: by type, each bucket of each item type in which items
 * were done; otherwise each such bucket of the report's type, whole, with its running means,
 * and then per model. Newest first.
 */
export function spendRows(r: Report, kind: Kind, by: By = 'type'): SpendRow[] {
	const rows: SpendRow[] = [];
	const put = (b: Bucket, of: string, s: Spend, means: boolean) =>
		rows.push({
			...s,
			at: b.at,
			of,
			mean_tokens: means ? b.mean_tokens : undefined,
			mean_cost: means ? b.mean_cost : undefined
		});
	if (PER_ITEM_KINDS.includes(kind) && by === 'type') {
		for (const type of TYPES)
			for (const b of bucketsOf(r, type)) if (b.items > 0) put(b, type, b, false);
	} else {
		for (const b of bucketsOf(r)) {
			if (b.items === 0) continue;
			const { models, ...whole } = b;
			put(b, models && models.length > 1 ? ALL : r.type, whole, true);
			for (const m of models ?? []) put(b, m.model, m, false);
		}
	}
	// newest bucket first; within a bucket, the order the rows were put in
	const seq = new Map(rows.map((row, i) => [row, i]));
	return rows.sort((a, b) => (a.at !== b.at ? (a.at < b.at ? 1 : -1) : seq.get(a)! - seq.get(b)!));
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

export function build(kind: Kind, report: Report, t: Theme, epic?: string, by: By = 'type'): Opt {
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
			return tokenRate(r, t);
		case 'tokens-spent':
			return tokensSpent(r, t);
		case 'tokens-per-item':
			return tokensPerItem(r, t, by);
		case 'tokens-per-dollar':
			return tokensPerDollar(r, t);
		case 'cost-spent':
			return costSpent(r, t);
		case 'cost-per-item':
			return costPerItem(r, t, by);
		case 'cost':
			return cost(r, t, epic);
		case 'completion-time':
			return completionTime(r, t);
		case 'completion-cost':
			return completionCost(r, t);
	}
}
