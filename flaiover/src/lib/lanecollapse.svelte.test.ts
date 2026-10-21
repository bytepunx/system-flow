// S-0462: which of the board's lanes are open on a narrow screen, remembered per browser, with
// Ready, In progress, and Review open by default.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { LaneCollapse, parseOpen } from './lanecollapse.svelte';

const KEY = 'flaiover-board-lanes';
const defaults = ['ready', 'in-progress', 'review'];

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
	it('opens Ready, In progress, and Review when nothing is stored or what is stored is not a list of lanes', () => {
		for (const raw of [null, '', 'not json', '{}', '"ready"', 'null', '[1]', '["ready","nope"]'])
			expect(parseOpen(raw)).toEqual(defaults);
	});

	it('reads a stored list back in the board order', () => {
		expect(parseOpen('["done","backlog"]')).toEqual(['backlog', 'done']);
		expect(parseOpen('[]')).toEqual([]);
	});
});

describe('LaneCollapse', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('opens Ready, In progress, and Review in a fresh browser', () => {
		const lanes = new LaneCollapse();
		expect(lanes.open).toEqual(defaults);
		expect(lanes.isOpen('ready')).toBe(true);
		expect(lanes.isOpen('backlog')).toBe(false);
	});

	it('toggles a lane open and closed, writing the whole choice each time', () => {
		const lanes = new LaneCollapse();
		lanes.toggle('backlog');
		expect(lanes.open).toEqual(['backlog', 'ready', 'in-progress', 'review']);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(lanes.open);
		lanes.toggle('review');
		expect(lanes.isOpen('review')).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(['backlog', 'ready', 'in-progress']);
	});

	it('keeps the choice across a reload', () => {
		new LaneCollapse().toggle('ready');
		expect(new LaneCollapse().open).toEqual(['in-progress', 'review']);
	});

	it('ignores a corrupt stored value', () => {
		localStorage.setItem(KEY, '["ready",');
		expect(new LaneCollapse().open).toEqual(defaults);
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
		const lanes = new LaneCollapse();
		expect(lanes.open).toEqual(defaults);
		lanes.toggle('done');
		expect(lanes.isOpen('done')).toBe(true);
	});
});
