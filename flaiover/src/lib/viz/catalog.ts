// The chart catalogue: the kinds of chart, their groups, and their titles. It
// imports nothing, so the rail can list the charts without loading their
// builders (S-0484).

/** The charts of how work flows. */
export const FLOW_KINDS = [
	'cycle-time',
	'burn-up',
	'cfd',
	'time-in-state',
	'throughput',
	'agent-waiting',
	'queue'
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
	'time-per-model',
	'cost-per-model'
] as const;
/** The charts of forecasts and estimates against what happened (S-0212). */
export const FORECAST_KINDS = [
	'forecast-accuracy',
	'delivery-accuracy',
	'forecast-by-model'
] as const;
export type ForecastKind = (typeof FORECAST_KINDS)[number];
/**
 * The charts of what claims cost in parallelism and holds (S-0214), of how stories coordinated
 * around them (S-0337), and of touches drift (S-0214).
 */
export const CLAIMS_KINDS = ['parallelism', 'hold-time', 'coordination', 'touches-drift'] as const;
export type ClaimsKind = (typeof CLAIMS_KINDS)[number];
/** The charts the Planning group lists: forecasts, then claims. */
export const PLANNING_KINDS = [...FORECAST_KINDS, ...CLAIMS_KINDS] as const;
export type PlanningKind = (typeof PLANNING_KINDS)[number];
/** The charts of what the strategic agents add against delivery, by the day (S-0216). */
export const STRATEGIC_KINDS = ['strategic-cost', 'strategic-use'] as const;
/** The charts of what waiting costs and what the pull order will cost (S-0213). */
export const COD_KINDS = ['cod-outstanding', 'cod-incurred', 'cod-order'] as const;
export const KINDS = [
	...FLOW_KINDS,
	...USAGE_KINDS,
	...PLANNING_KINDS,
	...STRATEGIC_KINDS,
	...COD_KINDS
] as const;
export type Kind = (typeof KINDS)[number];
export const TITLES: Record<Kind, string> = {
	'cycle-time': 'Cycle Time',
	'burn-up': 'Burn-up',
	cfd: 'Cumulative Flow',
	'time-in-state': 'Time in State',
	throughput: 'Throughput',
	'agent-waiting': 'Agent Waiting',
	queue: 'Merge Queue',
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
	'forecast-by-model': 'Forecast Error / Model',
	parallelism: 'Parallelism',
	'hold-time': 'Hold Time',
	coordination: 'Coordination',
	'touches-drift': 'Touches Drift',
	'strategic-cost': 'Strategic Cost',
	'strategic-use': 'Strategic Use',
	'cod-outstanding': 'CoD Outstanding',
	'cod-incurred': 'CoD Incurred',
	'cod-order': 'CoD by Order'
};
/** The name of a group of charts. */
export type ChartGroupKey = 'flow' | 'usage' | 'planning' | 'strategic';
/** The groups the rail and the charts page list the charts under (S-0484). */
export const CHART_GROUPS: readonly {
	key: ChartGroupKey;
	label: string;
	kinds: readonly Kind[];
}[] = [
	{ key: 'flow', label: 'flow', kinds: FLOW_KINDS },
	{ key: 'usage', label: 'usage', kinds: USAGE_KINDS },
	{ key: 'planning', label: 'planning', kinds: [...PLANNING_KINDS, ...COD_KINDS] },
	{ key: 'strategic', label: 'strategic', kinds: STRATEGIC_KINDS }
];
