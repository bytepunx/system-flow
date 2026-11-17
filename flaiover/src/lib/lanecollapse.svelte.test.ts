// S-0462: which of the board's lanes are open on a narrow screen, remembered per browser, with
// Ready, In progress, and Review open by default; since S-0753 the lanes are the workflow's.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { defaultOpen, LaneCollapse, parseOpen } from './lanecollapse.svelte';
import { baselineFlow, testingFlow } from './lanes.fixture';

const { lanes: baseline } = baselineFlow;

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
			expect(parseOpen(raw, baseline)).toEqual(defaults);
	});

	it('reads a stored list back in the board order', () => {
		expect(parseOpen('["done","backlog"]', baseline)).toEqual(['backlog', 'done']);
		expect(parseOpen('[]', baseline)).toEqual([]);
	});

	// S-0753: the lanes and the default are the workflow's
	it("opens the lanes where the work moves, but the initial one, and reads only the board's lanes", () => {
		expect(defaultOpen(baseline)).toEqual(defaults);
		expect(defaultOpen(testingFlow.lanes)).toEqual(['ready', 'doing', 'testing', 'review']);
		expect(parseOpen(null, testingFlow.lanes)).toEqual(['ready', 'doing', 'testing', 'review']);
		expect(parseOpen('["testing","backlog"]', testingFlow.lanes)).toEqual(['backlog', 'testing']);
		expect(parseOpen('["in-progress"]', testingFlow.lanes)).toEqual(defaultOpen(testingFlow.lanes));
	});
});

describe('LaneCollapse', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('opens Ready, In progress, and Review in a fresh browser', () => {
		const lanes = new LaneCollapse(baseline);
		expect(lanes.open).toEqual(defaults);
		expect(lanes.isOpen('ready')).toBe(true);
		expect(lanes.isOpen('backlog')).toBe(false);
	});

	it('toggles a lane open and closed, writing the whole choice each time', () => {
		const lanes = new LaneCollapse(baseline);
		lanes.toggle('backlog');
		expect(lanes.open).toEqual(['backlog', 'ready', 'in-progress', 'review']);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(lanes.open);
		lanes.toggle('review');
		expect(lanes.isOpen('review')).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(['backlog', 'ready', 'in-progress']);
	});

	it('keeps the choice across a reload', () => {
		new LaneCollapse(baseline).toggle('ready');
		expect(new LaneCollapse(baseline).open).toEqual(['in-progress', 'review']);
	});

	it('ignores a corrupt stored value', () => {
		localStorage.setItem(KEY, '["ready",');
		expect(new LaneCollapse(baseline).open).toEqual(defaults);
	});

	it('reads the choice again when the board has other lanes, and keeps it for the same', () => {
		const lanes = new LaneCollapse();
		expect(lanes.open).toEqual([]);
		lanes.setLanes(baseline);
		expect(lanes.open).toEqual(defaults);
		lanes.toggle('ready');
		lanes.setLanes([...baseline]);
		expect(lanes.open).toEqual(['in-progress', 'review']);
		lanes.setLanes(testingFlow.lanes);
		expect(lanes.open).toEqual(['ready', 'doing', 'testing', 'review']);
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
		const lanes = new LaneCollapse(baseline);
		expect(lanes.open).toEqual(defaults);
		lanes.toggle('done');
		expect(lanes.isOpen('done')).toBe(true);
	});
});
