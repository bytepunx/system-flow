import { describe, expect, it } from 'vitest';
import type { Flow } from './lanes';
import { addedStateFlow, baselineFlow, testingFlow } from './lanes.fixture';
import {
	acceptMove,
	acceptedState,
	cancelledState,
	isAccepted,
	isActive,
	isCancelled,
	isClosed,
	isOpen,
	isPullable,
	isQueue,
	isUnderway,
	isWaiting,
	itemMoves,
	kindOf,
	sendBackMove,
	stateSlot,
	stateTitle
} from './workflow';

const empty: Flow = { lanes: [], transitions: [] };
const targets = (flow: Flow, state: string, type: string) =>
	itemMoves(flow, state, type).map((m) => m.to);

// S-0758: a state's kind, colour, and moves, read from the lanes and transitions board.get answers
describe('workflow', () => {
	it('reads each state kind from its lane', () => {
		expect(baselineFlow.lanes.map((l) => kindOf(baselineFlow, l.name))).toEqual([
			'queue',
			'queue',
			'active',
			'waiting',
			'accepted',
			'cancelled'
		]);
		expect(kindOf(testingFlow, 'doing')).toBe('active');
		expect(kindOf(addedStateFlow, 'testing')).toBe('active');
		expect(kindOf(baselineFlow, 'blocked')).toBeUndefined();
		expect(kindOf(baselineFlow, undefined)).toBeUndefined();
	});

	it('answers each kind predicate from the lane, not the name', () => {
		expect(isQueue(baselineFlow, 'ready')).toBe(true);
		expect(isActive(baselineFlow, 'in-progress')).toBe(true);
		expect(isActive(testingFlow, 'doing')).toBe(true);
		expect(isActive(testingFlow, 'in-progress')).toBe(false);
		expect(isWaiting(baselineFlow, 'review')).toBe(true);
		expect(isAccepted(baselineFlow, 'done')).toBe(true);
		expect(isCancelled(baselineFlow, 'cancelled')).toBe(true);
		expect(isQueue(baselineFlow, 'review')).toBe(false);
	});

	it('counts accepted and cancelled closed, and anything else open', () => {
		expect(baselineFlow.lanes.map((l) => isClosed(baselineFlow, l.name))).toEqual([
			false,
			false,
			false,
			false,
			true,
			true
		]);
		expect(isOpen(baselineFlow, 'review')).toBe(true);
		expect(isOpen(baselineFlow, 'done')).toBe(false);
		expect(isOpen(baselineFlow, 'blocked')).toBe(true);
		for (const flow of [testingFlow, addedStateFlow]) {
			expect(isOpen(flow, 'testing')).toBe(true);
			expect(isClosed(flow, 'testing')).toBe(false);
			expect(isUnderway(flow, 'testing')).toBe(true);
		}
	});

	it('counts active and waiting underway', () => {
		expect(baselineFlow.lanes.map((l) => isUnderway(baselineFlow, l.name))).toEqual([
			false,
			false,
			true,
			true,
			false,
			false
		]);
	});

	it('finds the state an item is pulled from by the pull effect', () => {
		for (const flow of [baselineFlow, testingFlow, addedStateFlow]) {
			expect(flow.lanes.filter((l) => isPullable(flow, l.name)).map((l) => l.name)).toEqual([
				'ready'
			]);
		}
		expect(isPullable(baselineFlow, undefined)).toBe(false);
	});

	it('names the accepted and cancelled states', () => {
		for (const flow of [baselineFlow, testingFlow, addedStateFlow]) {
			expect(acceptedState(flow)).toBe('done');
			expect(cancelledState(flow)).toBe('cancelled');
		}
	});

	it('titles a state by its lane, else by its name', () => {
		const titled: Flow = {
			lanes: [{ name: 'review', title: 'In review', kind: 'waiting' }],
			transitions: []
		};
		expect(stateTitle(titled, 'review')).toBe('In review');
		expect(stateTitle(baselineFlow, 'review')).toBe('review');
		expect(stateTitle(baselineFlow, 'blocked')).toBe('blocked');
		expect(stateTitle(baselineFlow, undefined)).toBe('');
	});

	it('finds the accept move by its effect', () => {
		expect(acceptMove(baselineFlow, 'review', 'story')?.to).toBe('done');
		expect(acceptMove(baselineFlow, 'review', 'task')).toBeUndefined();
		expect(acceptMove(testingFlow, 'review', 'epic')?.by).toEqual(['acceptance']);
		expect(acceptMove(baselineFlow, 'in-progress', 'story')).toBeUndefined();
	});

	it('finds the send-back move by its effect', () => {
		expect(sendBackMove(baselineFlow, 'review', 'story')?.to).toBe('in-progress');
		expect(sendBackMove(testingFlow, 'review', 'story')?.to).toBe('doing');
		expect(sendBackMove(addedStateFlow, 'review', 'task')?.to).toBe('in-progress');
		expect(sendBackMove(addedStateFlow, 'testing', 'story')).toBeUndefined();
	});

	it('offers the baseline moves forward, then back with a reason, then cancel', () => {
		const f = baselineFlow;
		expect(targets(f, 'backlog', 'story')).toEqual(['ready', 'cancelled']);
		expect(targets(f, 'ready', 'story')).toEqual(['in-progress', 'cancelled']);
		expect(targets(f, 'in-progress', 'story')).toEqual(['review', 'cancelled']);
		expect(targets(f, 'in-progress', 'task')).toEqual(['review', 'done', 'cancelled']);
		expect(targets(f, 'review', 'story')).toEqual(['done', 'in-progress']);
		expect(targets(f, 'done', 'story')).toEqual([]);
		expect(targets(f, 'cancelled', 'story')).toEqual([]);
		expect(itemMoves(f, 'review', 'story')).toEqual([
			{ to: 'done', title: 'done', reason: false, effects: ['accept'] },
			{ to: 'in-progress', title: 'in-progress', reason: true, effects: ['send-back'] }
		]);
		expect(itemMoves(f, 'ready', 'story')[1]).toEqual({
			to: 'cancelled',
			title: 'cancelled',
			reason: true,
			effects: ['cancel']
		});
	});

	it('offers an added lane its moves, and leaves its plain move back to the board', () => {
		const f = addedStateFlow;
		expect(targets(f, 'in-progress', 'story')).toEqual(['testing', 'review', 'cancelled']);
		expect(targets(f, 'in-progress', 'task')).toEqual(['testing', 'review', 'done', 'cancelled']);
		expect(targets(f, 'testing', 'story')).toEqual(['review', 'cancelled']);
		expect(targets(f, 'testing', 'task')).toEqual(['review', 'done', 'cancelled']);
		expect(targets(f, 'review', 'story')).toEqual(['done', 'in-progress']);
		expect(targets(testingFlow, 'review', 'story')).toEqual(['done', 'doing']);
	});

	it('gives the baseline the colour slots palette.ts gives its states by name', () => {
		const slots = baselineFlow.lanes.map((l) => [l.name, stateSlot(baselineFlow, l.name)]);
		expect(Object.fromEntries(slots)).toEqual({
			backlog: 6,
			ready: 3,
			'in-progress': 0,
			review: 7,
			done: 5,
			cancelled: 1
		});
		expect(stateSlot(baselineFlow, 'blocked')).toBeUndefined();
		expect(stateSlot(baselineFlow, undefined)).toBeUndefined();
	});

	it('gives a lane past its kind slots the next spare slot, then none', () => {
		expect(stateSlot(addedStateFlow, 'in-progress')).toBe(0);
		expect(stateSlot(addedStateFlow, 'testing')).toBe(2);
		expect(stateSlot(addedStateFlow, 'review')).toBe(7);
		expect(stateSlot(testingFlow, 'doing')).toBe(0);
		expect(stateSlot(testingFlow, 'testing')).toBe(2);
		const active = (name: string) => ({ name, kind: 'active' as const });
		const crowded: Flow = { lanes: ['a', 'b', 'c', 'd'].map(active), transitions: [] };
		expect(crowded.lanes.map((l) => stateSlot(crowded, l.name))).toEqual([0, 2, 4, undefined]);
	});

	it('reads nothing from an empty flow', () => {
		expect(kindOf(empty, 'ready')).toBeUndefined();
		expect(isOpen(empty, 'done')).toBe(true);
		expect(isClosed(empty, 'done')).toBe(false);
		expect(isUnderway(empty, 'in-progress')).toBe(false);
		expect(isPullable(empty, 'ready')).toBe(false);
		expect(acceptedState(empty)).toBeUndefined();
		expect(cancelledState(empty)).toBeUndefined();
		expect(acceptMove(empty, 'review', 'story')).toBeUndefined();
		expect(sendBackMove(empty, 'review', 'story')).toBeUndefined();
		expect(itemMoves(empty, 'ready', 'story')).toEqual([]);
		expect(stateSlot(empty, 'ready')).toBeUndefined();
		expect(stateTitle(empty, 'ready')).toBe('ready');
	});
});
