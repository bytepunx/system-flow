// The cost of delay inputs as the new-item and edit forms hold them (S-0204): three optional text
// fields, what they say in one line, and what a save sends to flai. flai parses and validates the
// values and stamps the operator as who set them; the forms only collect the text.
import { amount, type CostOfDelay } from '$lib/planning';

/** The currency of amounts when the manifest's planning.currency is not set, as flai has it. */
export const DEFAULT_CURRENCY = 'USD';

/** The input keys, in the order the forms show them; flai's item.new and item.edit take the same names. */
export const COST_KEYS = ['revenue_per_week', 'penalty_per_week', 'time_lost_per_cycle'] as const;
export type CostKey = (typeof COST_KEYS)[number];

/** Each input as the text of its field; an empty field is an input not given. */
export type CostFields = Record<CostKey, string>;

/** The fields as an item's cost of delay fills them: amounts as plain numbers, empty when absent. */
export function costFields(c?: CostOfDelay): CostFields {
	const i = c?.inputs;
	return {
		revenue_per_week: i?.revenue_per_week !== undefined ? String(i.revenue_per_week) : '',
		penalty_per_week: i?.penalty_per_week !== undefined ? String(i.penalty_per_week) : '',
		time_lost_per_cycle: i?.time_lost_per_cycle ?? ''
	};
}

function money(v: string, currency: string): string {
	const n = Number(v);
	return `${Number.isFinite(n) ? amount(n) : v} ${currency}/week`;
}

/** The inputs given, in one line, such as "revenue 1,200 USD/week · 6h lost per cycle"; empty when none is. */
export function costSummary(f: CostFields, currency: string): string {
	const revenue = f.revenue_per_week.trim();
	const penalty = f.penalty_per_week.trim();
	const lost = f.time_lost_per_cycle.trim();
	return [
		...(revenue ? [`revenue ${money(revenue, currency)}`] : []),
		...(penalty ? [`penalty ${money(penalty, currency)}`] : []),
		...(lost ? [`${lost} lost per cycle`] : [])
	].join(' · ');
}

/**
 * What item.edit's cost_of_delay gets: each input whose field changed from what was loaded, an
 * emptied one as "" so that flai removes it; undefined when none changed. The planner's value is
 * never sent, so clearing every input leaves it.
 */
export function costPatch(before: CostFields, after: CostFields): Partial<CostFields> | undefined {
	const out: Partial<CostFields> = {};
	for (const k of COST_KEYS) if (after[k].trim() !== before[k].trim()) out[k] = after[k].trim();
	return Object.keys(out).length ? out : undefined;
}

/** What item.new's cost_of_delay gets: the inputs given; undefined when none is. */
export function costInputs(f: CostFields): Partial<CostFields> | undefined {
	const out: Partial<CostFields> = {};
	for (const k of COST_KEYS) if (f[k].trim()) out[k] = f[k].trim();
	return Object.keys(out).length ? out : undefined;
}
