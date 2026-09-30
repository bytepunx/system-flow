// The window the charts page covers (S-0168): 30 days until the user picks another. Like the
// theme, the choice is a per-browser convenience kept in localStorage, so every chart opens at it
// after a reload or a visit to another page, not only when reached from another chart.
import { WINDOWS } from '$lib/viz/charts';

export type Window = (typeof WINDOWS)[number];
const KEY = 'flaiover-chart-window';
const DEFAULT: Window = '30d';

/** The stored window, or the default for anything that is not one the page offers. */
export function parseWindow(raw: string | null): Window {
	return (WINDOWS as readonly string[]).includes(raw ?? '') ? (raw as Window) : DEFAULT;
}

function stored(): Window {
	try {
		return parseWindow(localStorage.getItem(KEY));
	} catch {
		return DEFAULT;
	}
}

export class ChartWindow {
	since = $state<Window>(stored());

	/** Choose a window and remember it. */
	set(since: string): void {
		this.since = parseWindow(since);
		try {
			localStorage.setItem(KEY, this.since);
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}
}

export const chartWindow = new ChartWindow();
