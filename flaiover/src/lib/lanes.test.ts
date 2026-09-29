import { describe, expect, it } from 'vitest';
import { backOf, forwardOf, hasLimit, LANES, movesTo, needsReason, startLane } from './lanes';

// S-0167, ADR-0055: what each lane's menu offers
describe('lanes', () => {
	it('starts a new item in backlog, ready, or in-progress, and anywhere else in backlog', () => {
		expect(LANES.map((l) => startLane(l))).toEqual([
			'backlog',
			'ready',
			'in-progress',
			'backlog',
			'backlog',
			'backlog'
		]);
		expect(startLane(null)).toBe('backlog');
		expect(startLane('nonsense')).toBe('backlog');
		expect(movesTo('backlog')).toEqual([]);
		expect(movesTo('ready')).toEqual(['ready']);
		expect(movesTo('in-progress')).toEqual(['ready', 'in-progress']);
	});

	it('moves forward from backlog only, and back from all but backlog and done', () => {
		expect(LANES.map((l) => forwardOf(l))).toEqual(['ready', null, null, null, null, null]);
		expect(LANES.map((l) => backOf(l))).toEqual([
			null,
			'backlog',
			'ready',
			'in-progress',
			null,
			'backlog'
		]);
		expect(needsReason('review', 'in-progress')).toBe(true);
		expect(needsReason('in-progress', 'ready')).toBe(false);
	});

	it('changes the WIP limit of ready, in-progress, and review only', () => {
		expect(LANES.filter(hasLimit)).toEqual(['ready', 'in-progress', 'review']);
	});
});
