// What agents spent on a work item (S-0143, ADR-0051): the usage block of its front matter,
// as flai writes it, and how the dashboard says it. What strategic agents such as the planner
// spent on it is carried apart, and an item with a forecast or estimate has an expected cost
// (S-0225, ADR-0083).

export type ModelUsage = {
	model: string;
	input: number;
	output: number;
	cache_read: number;
	cache_write: number;
	cost: number;
};

/** What one kind of strategic agent spent on an item: apart from its agents' figures, and estimated. */
export type StrategicUsage = {
	kind: string; // planner, orchestrator, or analyzer
	seconds: number;
	estimated: boolean;
	models: ModelUsage[];
};

export type Usage = {
	source: 'log' | 'sum';
	seconds: number;
	estimated?: boolean;
	models: ModelUsage[];
	strategic?: StrategicUsage[];
};

/** What an item is expected to cost before agents work it: its forecast or estimate priced at the project's cost per agent hour. */
export type ExpectedCost = {
	cost: number;
	from: 'forecast' | 'estimate';
	estimated: boolean;
};

export const modelTokens = (m: ModelUsage) => m.input + m.output + m.cache_read + m.cache_write;

export function tokensOf(u: Usage): number {
	return u.models.reduce((n, m) => n + modelTokens(m), 0);
}

export function costOf(u: Usage): number {
	return u.models.reduce((c, m) => c + m.cost, 0);
}

/** A count of tokens made short, as flai prints it: 950, 12.3K, 20.1M. */
export function count(n: number): string {
	if (n >= 1e9) return `${(n / 1e9).toFixed(1)}B`;
	if (n >= 1e6) return `${(n / 1e6).toFixed(1)}M`;
	if (n >= 1e4) return `${(n / 1e3).toFixed(1)}K`;
	return String(Math.round(n));
}

/** Dollars, to the cent, or to the tenth of a cent under a dollar. */
export function dollars(c: number): string {
	return c > 0 && c < 1 ? `$${c.toFixed(3)}` : `$${c.toFixed(2)}`;
}

/** A length of time in hours and minutes, or seconds under a minute. */
export function duration(seconds: number): string {
	const h = Math.floor(seconds / 3600);
	const m = Math.floor((seconds % 3600) / 60);
	if (h > 0) return m > 0 ? `${h}h${m}m` : `${h}h`;
	return m > 0 ? `${m}m` : `${Math.round(seconds)}s`;
}

/** The usage in a line: tokens, cost, agent time, and where the numbers came from. */
export function usageLine(u: Usage): string {
	const cost = dollars(costOf(u)) + (u.estimated ? ' (estimated)' : '');
	const from = u.source === 'sum' ? 'summed from its children' : "measured from its agents' logs";
	return `${count(tokensOf(u))} tokens · ${cost} · ${duration(u.seconds)} of agent work · ${from}`;
}

/** One model's line: its tokens by kind, and its cost. */
export function modelLine(m: ModelUsage): string {
	return `${m.model}: ${count(m.input)} in · ${count(m.output)} out · ${count(m.cache_read)} cache read · ${count(m.cache_write)} cache write · ${dollars(m.cost)}`;
}

/** Whether the agents spent anything: an empty measurement, or strategic usage alone, says not. */
export function spent(u: Usage | undefined | null): u is Usage {
	return !!u && (u.models.length > 0 || u.seconds > 0);
}

/** A line per kind of strategic agent that spent on the item, apart from its agents' usage. */
export function strategicLines(u: Usage | undefined | null): string[] {
	return (u?.strategic ?? []).map((s) => {
		const tokens = s.models.reduce((n, m) => n + modelTokens(m), 0);
		const cost = s.models.reduce((c, m) => c + m.cost, 0);
		return `${s.kind}, strategic: ${count(tokens)} tokens · ${dollars(cost)}${s.estimated ? ' (estimated)' : ''} · ${duration(s.seconds)}`;
	});
}

/** The expected cost in a line, marked as an estimate and saying what it was priced from. */
export function expectedCostLine(e: ExpectedCost | undefined | null): string | undefined {
	if (!e) return undefined;
	return `expected cost: ${dollars(e.cost)} (${e.estimated ? 'estimated, ' : ''}from the ${e.from})`;
}
