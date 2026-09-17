// Validated categorical palette (design-system reference instance, both modes
// checked with the dataviz validator on 2026-09-17). Slots are assigned in
// fixed order and never cycled; the first three are safe for all-pairs forms
// such as scatter, so scatter colours by at most three groups.
export const CATEGORICAL = {
	light: ['#2a78d6', '#eb6834', '#1baf7a', '#eda100', '#e87ba4', '#008300', '#4a3aa7', '#e34948'],
	dark: ['#3987e5', '#d95926', '#199e70', '#c98500', '#d55181', '#008300', '#9085e9', '#e66767']
} as const;

// The workflow states get fixed slots so a state keeps its colour across charts.
export const STATE_SLOT: Record<string, number> = {
	backlog: 6,
	ready: 3,
	'in-progress': 0,
	review: 1,
	done: 2,
	cancelled: 7
};
// Natures likewise.
export const NATURE_SLOT: Record<string, number> = {
	feature: 0,
	improvement: 2,
	remediation: 1,
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
				surface: '#1a1a19',
				text: '#ffffff',
				textSecondary: '#c3c2b7',
				grid: '#383835',
				series: CATEGORICAL.dark
			}
		: {
				dark,
				surface: '#fcfcfb',
				text: '#0b0b0b',
				textSecondary: '#52514e',
				grid: '#e6e5e1',
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
