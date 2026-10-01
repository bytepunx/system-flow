// S-0176: a story's task plan as flai gives it, and the card's line for its tasks.
import { describe, expect, it } from 'vitest';
import { planFrom, stateLabel, tasksLine, tasksTitle, unplaced, waitsFor } from './taskplan';

describe('planFrom', () => {
	it('is undefined when flai sent no plan', () => {
		expect(planFrom(undefined)).toBeUndefined();
		expect(planFrom(null)).toBeUndefined();
	});

	it('keeps the contract’s shape and drops the lists flai sent empty or null', () => {
		expect(
			planFrom({
				tasks: [
					{ id: 'T-0001', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] },
					{ id: 'T-0002', state: 'in-progress', after: null, waiting_for: [] },
					{ id: 'T-0003', state: 'ready' }
				],
				layers: [['T-0002', 'T-0003'], null, ['T-0001']]
			})
		).toEqual({
			tasks: [
				{ id: 'T-0001', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] },
				{ id: 'T-0002', state: 'in-progress' },
				{ id: 'T-0003', state: 'ready' }
			],
			layers: [['T-0002', 'T-0003'], [], ['T-0001']]
		});
		expect(planFrom({ tasks: [], layers: null })).toEqual({ tasks: [], layers: [] });
	});
});

describe('a task of the plan', () => {
	it('says its state in words, and an unknown one as it is', () => {
		expect(['ready', 'waiting', 'in-progress', 'done', 'cancelled', 'odd'].map(stateLabel)).toEqual(
			['ready to start', 'waiting', 'in progress', 'done', 'cancelled', 'odd']
		);
	});

	it('waits for the tasks not yet done, or its whole after when flai named none', () => {
		const after = ['T-0002', 'T-0003'];
		expect(waitsFor({ id: 'T-0001', state: 'waiting', after, waiting_for: ['T-0003'] })).toEqual([
			'T-0003'
		]);
		expect(waitsFor({ id: 'T-0001', state: 'waiting', after })).toEqual(after);
		expect(waitsFor({ id: 'T-0001', state: 'ready' })).toEqual([]);
	});

	it('names the open tasks in no layer, which a cycle leaves out', () => {
		expect(
			unplaced({
				tasks: [
					{ id: 'T-0001', state: 'ready' },
					{ id: 'T-0002', state: 'waiting', after: ['T-0003'] },
					{ id: 'T-0003', state: 'waiting', after: ['T-0002'] },
					{ id: 'T-0004', state: 'done' },
					{ id: 'T-0005', state: 'cancelled' }
				],
				layers: [['T-0001']]
			})
		).toEqual(['T-0002', 'T-0003']);
	});
});

describe('the card’s line for a story’s tasks', () => {
	const summary = { ready: 1, waiting: 2, in_progress: 1, done: 3, layers: 3 };

	it('counts the open tasks by state and the layers', () => {
		expect(tasksLine(summary)).toBe('tasks 1 in progress · 1 ready · 2 waiting · 3 layers');
		expect(tasksTitle(summary)).toBe(
			'tasks: 1 in progress, 1 ready to start, 2 waiting, 3 done; the plan has 3 layers'
		);
	});

	it('leaves out a state with none, and says when every task is done', () => {
		expect(tasksLine({ ...summary, in_progress: 0, waiting: 0, layers: 1 })).toBe(
			'tasks 1 ready · 1 layer'
		);
		expect(tasksLine({ ready: 0, waiting: 0, in_progress: 0, done: 4, layers: 2 })).toBe(
			'tasks 4 done · 2 layers'
		);
	});
});
