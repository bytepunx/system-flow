// A work item's planning data (S-0199): its cost of delay and its forecast, as flai writes them in
// front matter and returns them in item.get, and how the item page says them. Amounts are in the
// project's currency, which item.get does not carry, so they are shown as plain numbers.

/** What a cost of delay is worked out from; an absent amount is unknown, not zero. */
export type CostInputs = {
	revenue_per_week?: number;
	penalty_per_week?: number;
	time_lost_per_cycle?: string; // a Go duration
};

/** What each week of waiting for an epic or story costs, and who set it and when. */
export type CostOfDelay = {
	inputs?: CostInputs;
	value?: number; // per week
	by?: string;
	at?: string;
};

/** When a story is expected to be delivered, what that rests on, and who set it and when. */
export type Forecast = {
	duration?: string; // a Go duration
	delivery?: string; // a UTC timestamp
	basis?: string;
	by?: string;
	at?: string;
};

/** An amount with thousands separators and at most two decimals, in no particular currency. */
export function amount(n: number): string {
	return n.toLocaleString('en-US', { maximumFractionDigits: 2 });
}

function setBy(by?: string, at?: string): string[] {
	if (by && at) return [`set by ${by} at ${at}`];
	if (by) return [`set by ${by}`];
	if (at) return [`set at ${at}`];
	return [];
}

/** The cost of delay in lines: its value per week, the inputs given, and who set it; none when absent. */
export function costOfDelayLines(c?: CostOfDelay): string[] {
	if (!c) return [];
	const inputs = [
		...(c.inputs?.revenue_per_week !== undefined
			? [`revenue ${amount(c.inputs.revenue_per_week)} per week`]
			: []),
		...(c.inputs?.penalty_per_week !== undefined
			? [`penalty ${amount(c.inputs.penalty_per_week)} per week`]
			: []),
		...(c.inputs?.time_lost_per_cycle ? [`${c.inputs.time_lost_per_cycle} lost per cycle`] : [])
	];
	const rest = [...(inputs.length ? [inputs.join(' · ')] : []), ...setBy(c.by, c.at)];
	if (c.value === undefined && !rest.length) return [];
	const head =
		c.value !== undefined ? `cost of delay: ${amount(c.value)} per week` : 'cost of delay';
	return [head, ...rest];
}

/** The forecast in lines: its duration and delivery, its basis, and who set it; none when absent. */
export function forecastLines(f?: Forecast): string[] {
	if (!f) return [];
	const when = [
		...(f.duration ? [`${f.duration} of work`] : []),
		...(f.delivery ? [`delivery ${f.delivery}`] : [])
	];
	const rest = [...(f.basis ? [f.basis] : []), ...setBy(f.by, f.at)];
	if (!when.length && !rest.length) return [];
	return [when.length ? `forecast: ${when.join(', ')}` : 'forecast', ...rest];
}
