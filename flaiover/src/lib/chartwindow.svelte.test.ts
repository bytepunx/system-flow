// S-0168: the window of the charts page, remembered per browser, 30 days by default.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ChartWindow, parseWindow } from './chartwindow.svelte';

const KEY = 'flaiover-chart-window';

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('parseWindow', () => {
	it('reads a window the page offers, and 30 days for anything else', () => {
		expect(parseWindow('7d')).toBe('7d');
		expect(parseWindow('365d')).toBe('365d');
		for (const raw of [null, '', '8d', '30', 'week']) expect(parseWindow(raw)).toBe('30d');
	});
});

describe('ChartWindow', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('starts at 30 days in a fresh browser', () => {
		expect(new ChartWindow().since).toBe('30d');
	});

	it('writes the window chosen at once, and a new page reads it back', () => {
		const first = new ChartWindow();
		first.set('7d');
		expect(first.since).toBe('7d');
		expect(localStorage.getItem(KEY)).toBe('7d');
		expect(new ChartWindow().since).toBe('7d');
	});

	it('keeps the window for the page when storage refuses it', () => {
		vi.stubGlobal('localStorage', {
			getItem: () => {
				throw new Error('denied');
			},
			setItem: () => {
				throw new Error('denied');
			}
		});
		const w = new ChartWindow();
		expect(w.since).toBe('30d');
		w.set('90d');
		expect(w.since).toBe('90d');
	});
});
