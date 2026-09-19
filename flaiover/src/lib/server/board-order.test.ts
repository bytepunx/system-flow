// The dashboard lays out backlog and ready in flai's pull order (S-0057). The cases mirror
// flai/internal/workitem/order_test.go: one rule, read in two places.
import { describe, expect, it } from 'vitest';
import { inPullOrder } from './board';

const story = (id: string) => ({ id, type: 'story' });

describe('inPullOrder', () => {
	it('puts named stories first in the order’s order, then the rest by ID', () => {
		const cards = ['S-0003', 'S-0004', 'S-0005', 'S-0006'].map(story);
		expect(inPullOrder(cards, ['S-0002', 'S-0001', 'S-0005']).map((c) => c.id)).toEqual([
			'S-0005',
			'S-0003',
			'S-0004',
			'S-0006'
		]);
	});
	it('sorts unnamed stories by number, not as text', () => {
		const cards = ['S-0100', 'S-099', 'S-0020'].map(story);
		expect(inPullOrder(cards, []).map((c) => c.id)).toEqual(['S-0020', 'S-099', 'S-0100']);
	});
	it('ignores names that are not in the column, and a name given twice', () => {
		const cards = ['S-0001', 'S-0002'].map(story);
		expect(inPullOrder(cards, ['S-0009', 'S-0002', 'S-0002', 'S-0001']).map((c) => c.id)).toEqual([
			'S-0002',
			'S-0001'
		]);
	});
	it('leaves epics and tasks where they are and orders the stories among them', () => {
		const cards = [
			{ id: 'E-0001', type: 'epic' },
			story('S-0001'),
			story('S-0002'),
			{ id: 'T-0001', type: 'task' }
		];
		expect(inPullOrder(cards, ['S-0002']).map((c) => c.id)).toEqual([
			'E-0001',
			'S-0002',
			'S-0001',
			'T-0001'
		]);
	});
	it('does not change the array it was given', () => {
		const cards = ['S-0001', 'S-0002'].map(story);
		inPullOrder(cards, ['S-0002']);
		expect(cards.map((c) => c.id)).toEqual(['S-0001', 'S-0002']);
	});
});
