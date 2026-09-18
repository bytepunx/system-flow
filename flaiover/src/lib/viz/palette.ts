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
