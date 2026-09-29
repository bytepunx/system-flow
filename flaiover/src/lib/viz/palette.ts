// Categorical palette led by the brand hues (S-0044): slot 0 is the brand
// blue (#054a91 family), slot 2 the brand cyan (#23b5d3 family), slot 5 the
// brand green (#119822 family); the rest are supplementary hues stepped to
// the same bands. Both modes pass the dataviz validator (2026-09-18) against
// the chart surfaces below. Slots are assigned in fixed order and never
// cycled; scatter colours by at most three groups.
export const CATEGORICAL = {
	light: ['#2f6ec4', '#d0483f', '#0e8fa9', '#b5820a', '#c0559a', '#1d9a2e', '#6f5bc8', '#b0632a'],
	dark: ['#3f86d9', '#e05a52', '#1e9fbb', '#b88b1f', '#d55f9d', '#33a646', '#8b7ce6', '#c27a3c']
} as const;

// The workflow states get fixed slots so a state keeps its colour across charts.
export const STATE_SLOT: Record<string, number> = {
	backlog: 6,
	ready: 3,
	'in-progress': 0,
	review: 7,
	done: 5,
	cancelled: 1
};
// Natures likewise.
export const NATURE_SLOT: Record<string, number> = {
	feature: 0,
	improvement: 5,
	remediation: 7,
	research: 6,
	experiment: 4
};

// Models get fixed slots by family, so a model keeps its colour on every chart and every
// selection (S-0143). Another model takes slot 3 or 4, by its name, never 1 or 7, which the
// charts keep for cancelled work and for work over the p85.
export const MODEL_SLOT: Record<string, number> = { opus: 0, sonnet: 2, haiku: 5, fable: 6 };
export function modelSlot(model: string): number {
	for (const [family, slot] of Object.entries(MODEL_SLOT)) if (model.includes(family)) return slot;
	let h = 0;
	for (const ch of model) h = (h * 31 + ch.charCodeAt(0)) >>> 0;
	return h % 2 === 0 ? 3 : 4;
}

// A model's mark on a line, by family, so that two models are told apart by more than colour:
// opus and fable are close for a reader who sees little red or green (S-0163).
export const MODEL_SYMBOL: Record<string, string> = {
	opus: 'circle',
	sonnet: 'rect',
	haiku: 'triangle',
	fable: 'diamond'
};
export function modelSymbol(model: string): string {
	for (const [family, symbol] of Object.entries(MODEL_SYMBOL))
		if (model.includes(family)) return symbol;
	return 'roundRect';
}

// Item types on the charts that compare them (S-0163): the three slots that stay apart from
// each other in both modes for every pair (dataviz validator, all pairs, 2026-09-29), with a
// mark each.
export const TYPE_SLOT: Record<string, number> = { epic: 6, story: 5, task: 4 };
export const TYPE_SYMBOL: Record<string, string> = {
	epic: 'diamond',
	story: 'circle',
	task: 'triangle'
};

export type Theme = {
	dark: boolean;
	surface: string;
	text: string;
	textSecondary: string;
	grid: string;
	series: readonly string[];
};

export function theme(dark: boolean): Theme {
	return dark
		? {
				dark,
				surface: '#272220',
				text: '#f7f3e3',
				textSecondary: '#a79a95',
				grid: '#433b38',
				series: CATEGORICAL.dark
			}
		: {
				dark,
				surface: '#fdfbf3',
				text: '#1c1917',
				textSecondary: '#645853',
				grid: '#e3ddcc',
				series: CATEGORICAL.light
			};
}

export function colorFor(
	t: Theme,
	slots: Record<string, number>,
	key: string,
	fallbackIndex: number
): string {
	const slot = slots[key] ?? fallbackIndex;
	return t.series[slot % t.series.length];
}
