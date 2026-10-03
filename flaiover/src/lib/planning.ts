// A work item's planning data (S-0199): its cost of delay and its forecast, as flai writes them in
// front matter and returns them in item.get, and how the item page says them. Amounts are in the
// project's currency, which item.get does not carry, so they are shown as plain numbers.

/** What a cost of delay is worked out from, and who last set them and when; an absent amount is unknown, not zero. */
export type CostInputs = {
	revenue_per_week?: number;
	penalty_per_week?: number;
	time_lost_per_cycle?: string; // a Go duration
	by?: string;
	at?: string;
};

/** What each week of waiting for an epic or story costs; by and at are the value's (ADR-0080). */
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

/** A line with who set what it says and when after it, when either is known. */
function stamped(line: string, by?: string, at?: string): string {
	return [line, ...setBy(by, at)].join(', ');
}

/** The inputs given, in the order the forms take them; none when every amount is absent. */
function inputParts(i?: CostInputs): string[] {
	return [
		...(i?.revenue_per_week !== undefined
			? [`revenue ${amount(i.revenue_per_week)} per week`]
			: []),
		...(i?.penalty_per_week !== undefined
			? [`penalty ${amount(i.penalty_per_week)} per week`]
			: []),
		...(i?.time_lost_per_cycle ? [`${i.time_lost_per_cycle} lost per cycle`] : [])
	];
}

/** A time as flai writes it, 2006-01-02T15:04:05Z, in milliseconds; undefined when it is not one. */
function utc(s?: string): number | undefined {
	if (!s || !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/.test(s)) return undefined;
	const ms = Date.parse(s);
	if (Number.isNaN(ms)) return undefined;
	// a day out of range, such as February 30, parses to another day; flai refuses it, so do we
	return new Date(ms).toISOString() === s.replace('Z', '.000Z') ? ms : undefined;
}

/**
 * Whether the inputs changed after the value was set (ADR-0080): there are both, and the inputs'
 * at is later than the value's. Equal times, or a time that is not one, are not stale.
 */
export function costOfDelayStale(c?: CostOfDelay): boolean {
	if (c?.value === undefined || !inputParts(c.inputs).length) return false;
	const inputs = utc(c.inputs?.at);
	const value = utc(c.at);
	return inputs !== undefined && value !== undefined && inputs > value;
}

/**
 * The cost of delay in lines: its value per week with who set it, or that there is none yet, then
 * the inputs given with who set them; none when there is neither a value nor an input.
 */
export function costOfDelayLines(c?: CostOfDelay): string[] {
	const inputs = inputParts(c?.inputs);
	if (c?.value === undefined && !inputs.length) return [];
	const head =
		c?.value !== undefined
			? stamped(`cost of delay: ${amount(c.value)} per week`, c.by, c.at)
			: 'cost of delay: no value yet';
	return [
		head,
		...(inputs.length
			? [stamped(`inputs: ${inputs.join(' · ')}`, c?.inputs?.by, c?.inputs?.at)]
			: [])
	];
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
