// S-0485: whether a chart's explainer is open, remembered per browser, open on the first visit.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ChartExplainer, parseOpen } from './chartexplainer.svelte';

const KEY = 'flaiover-chart-explainer';

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('parseOpen', () => {
	it('reads closed only when the user closed it, and open for anything else', () => {
		expect(parseOpen('closed')).toBe(false);
		expect(parseOpen('open')).toBe(true);
		for (const raw of [null, '', 'false', '0', 'shut']) expect(parseOpen(raw)).toBe(true);
	});
});

describe('ChartExplainer', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('is open in a fresh browser', () => {
		expect(new ChartExplainer().open).toBe(true);
	});

	it('writes the choice at once, and a new page, on any chart, reads it back', () => {
		const first = new ChartExplainer();
		first.set(false);
		expect(first.open).toBe(false);
		expect(localStorage.getItem(KEY)).toBe('closed');
		expect(new ChartExplainer().open).toBe(false);
		first.set(true);
		expect(localStorage.getItem(KEY)).toBe('open');
		expect(new ChartExplainer().open).toBe(true);
	});

	it('keeps the choice for the page when storage refuses it', () => {
		vi.stubGlobal('localStorage', {
			getItem: () => {
				throw new Error('denied');
			},
			setItem: () => {
				throw new Error('denied');
			}
		});
		const e = new ChartExplainer();
		expect(e.open).toBe(true);
		e.set(false);
		expect(e.open).toBe(false);
	});
});
