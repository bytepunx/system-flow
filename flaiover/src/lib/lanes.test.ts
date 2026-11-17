import { describe, expect, it } from 'vitest';
import {
	archiveOrder,
	backOf,
	forwardOf,
	hasLimit,
	laneCounts,
	laneEntries,
	longCounts,
	movesTo,
	needsReason,
	operatorMove,
	shortCounts,
	startLane,
	startLanes,
	type Flow
} from './lanes';
import { baselineFlow, testingFlow } from './lanes.fixture';

const LANES = baselineFlow.lanes.map((l) => l.name);

// S-0167, ADR-0055: what each lane's menu offers, read since S-0753 from the baseline's lanes and
// transitions as board.get answers them
describe('lanes', () => {
	const flow = baselineFlow;
	it('starts a new item in backlog, ready, or in-progress, and anywhere else in backlog', () => {
		expect(LANES.map((l) => startLane(flow, l))).toEqual([
			'backlog',
			'ready',
			'in-progress',
			'backlog',
			'backlog',
			'backlog'
		]);
		expect(startLane(flow, null)).toBe('backlog');
		expect(startLane(flow, 'nonsense')).toBe('backlog');
		expect(movesTo(flow, 'backlog')).toEqual([]);
		expect(movesTo(flow, 'ready')).toEqual(['ready']);
		expect(movesTo(flow, 'in-progress')).toEqual(['ready', 'in-progress']);
	});

	it('moves forward from backlog only, and back from all but backlog and done', () => {
		expect(LANES.map((l) => forwardOf(flow, l))).toEqual(['ready', null, null, null, null, null]);
		expect(LANES.map((l) => backOf(flow, l))).toEqual([
			null,
			'backlog',
			'ready',
			'in-progress',
			null,
			'backlog'
		]);
		expect(needsReason(flow, 'review', 'in-progress')).toBe(true);
		expect(needsReason(flow, 'in-progress', 'ready')).toBe(false);
	});

	it('changes the WIP limit of ready, in-progress, and review only', () => {
		expect(LANES.filter((l) => hasLimit(flow, l))).toEqual(['ready', 'in-progress', 'review']);
	});

	// S-0202: the lane's entries, shared by the lane's menu and a card's
	it("lists each lane's menu entries in order, with the column a move goes to", () => {
		expect(laneEntries(flow, 'backlog')).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'forward', label: 'Move stories forward to ready…' }
		]);
		expect(laneEntries(flow, 'review')).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'back', label: 'Move stories back to in-progress…' },
			{ action: 'limit', label: 'Change WIP limit…' }
		]);
		expect(laneEntries(flow, 'done').map((e) => e.action)).toEqual(['create']);
		expect(laneEntries(flow, 'cancelled').map((e) => e.action)).toEqual(['create', 'back']);
	});

	// S-0343: cancelled cards stay until archived
	it('offers Archive All on the cancelled lane only, and only while it holds cards', () => {
		expect(laneEntries(flow, 'cancelled', 3)).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'back', label: 'Move stories back to backlog…' },
			{ action: 'archive', label: 'Archive All' }
		]);
		expect(laneEntries(flow, 'cancelled', 0).map((e) => e.action)).not.toContain('archive');
		for (const lane of LANES.filter((l) => l !== 'cancelled'))
			expect(laneEntries(flow, lane, 3).map((e) => e.action)).not.toContain('archive');
	});

	// S-0753: a move only a cascade, an agent, or flai makes is not the board's to offer
	it('offers only a move a person may make, for the type it moves', () => {
		expect(operatorMove(flow, 'review', 'cancelled', 'story')).toBeUndefined();
		expect(operatorMove(flow, 'in-progress', 'cancelled', 'story')?.effects).toEqual(['cancel']);
		expect(operatorMove(flow, 'review', 'done', 'story')?.effects).toEqual(['accept']);
		expect(operatorMove(flow, 'review', 'done', 'task')?.effects).toEqual([]);
		expect(operatorMove(flow, 'in-progress', 'done', 'story')).toBeUndefined();
		expect(operatorMove(flow, 'backlog', 'in-progress', 'story')).toBeUndefined();
		const agents: Flow = {
			lanes: flow.lanes,
			transitions: flow.transitions.map((m) =>
				m.to === 'backlog' ? { ...m, by: ['agent'], operator: false } : m
			)
		};
		expect(backOf(agents, 'ready')).toBeNull();
		expect(backOf(agents, 'cancelled')).toBeNull();
		expect(laneEntries(agents, 'ready').map((e) => e.action)).toEqual(['create', 'limit']);
		const epicsOnly: Flow = {
			lanes: flow.lanes,
			transitions: flow.transitions.map((m) => (m.to === 'ready' ? { ...m, types: ['epic'] } : m))
		};
		expect(forwardOf(epicsOnly, 'backlog')).toBeNull();
		expect(backOf(epicsOnly, 'in-progress')).toBeNull();
	});

	it('archives stories first, then tasks, then epics, each in the order the lane holds them', () => {
		const cards = [
			{ id: 'E-0001', type: 'epic' },
			{ id: 'T-0002', type: 'task' },
			{ id: 'S-0003', type: 'story' },
			{ id: 'E-0004', type: 'epic' },
			{ id: 'T-0001', type: 'task' },
			{ id: 'S-0001', type: 'story' }
		];
		expect(archiveOrder(cards)).toEqual([
			'S-0003',
			'S-0001',
			'T-0002',
			'T-0001',
			'E-0001',
			'E-0004'
		]);
		expect(cards[0].id).toBe('E-0001');
		expect(archiveOrder([])).toEqual([]);
	});

	// S-0256: what a lane holds, below its title
	it('counts the epics, stories, and tasks a lane holds, and nothing else', () => {
		const cards = ['epic', 'story', 'task', 'story', 'task', 'task', 'task', 'bug', 'toString'].map(
			(type) => ({ type })
		);
		expect(laneCounts(cards)).toEqual({ epic: 1, story: 2, task: 4 });
		expect(laneCounts([])).toEqual({ epic: 0, story: 0, task: 0 });
	});

	it('gives the counts short, epics | stories | tasks, and long, each pluralized', () => {
		expect(shortCounts({ epic: 1, story: 2, task: 4 })).toBe('1 | 2 | 4');
		expect(longCounts({ epic: 1, story: 2, task: 4 })).toBe('1 epic, 2 stories, 4 tasks');
		expect(longCounts({ epic: 0, story: 1, task: 1 })).toBe('0 epics, 1 story, 1 task');
		expect(longCounts({ epic: 2, story: 0, task: 0 })).toBe('2 epics, 0 stories, 0 tasks');
	});
});

// S-0753, criterion 4: in-progress renamed doing, and testing added after it with a WIP limit
describe('lanes of a workflow with a state renamed and one added', () => {
	const flow = testingFlow;
	const lanes = flow.lanes.map((l) => l.name);
	it('moves back a lane to the renamed and from the added one, and forward from backlog only', () => {
		expect(lanes.map((l) => forwardOf(flow, l))).toEqual([
			'ready',
			null,
			null,
			null,
			null,
			null,
			null
		]);
		expect(lanes.map((l) => backOf(flow, l))).toEqual([
			null,
			'backlog',
			'ready',
			'doing',
			'doing',
			null,
			'backlog'
		]);
		expect(needsReason(flow, 'review', 'doing')).toBe(true);
		expect(needsReason(flow, 'testing', 'doing')).toBe(false);
		expect(lanes.filter((l) => hasLimit(flow, l))).toEqual(['ready', 'doing', 'testing', 'review']);
		expect(laneEntries(flow, 'testing')).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'back', label: 'Move stories back to doing…' },
			{ action: 'limit', label: 'Change WIP limit…' }
		]);
	});

	it('starts a new item in the lanes the operator may move it to before it waits', () => {
		expect(startLanes(flow)).toEqual(['backlog', 'ready', 'doing', 'testing']);
		expect(startLane(flow, 'testing')).toBe('testing');
		expect(startLane(flow, 'in-progress')).toBe('backlog');
		expect(startLane(flow, 'review')).toBe('backlog');
		expect(movesTo(flow, 'testing')).toEqual(['ready', 'doing', 'testing']);
		expect(movesTo(flow, 'review')).toEqual([]);
	});

	it('offers the added lane as a move a person may make, and not one only a cascade makes', () => {
		expect(operatorMove(flow, 'doing', 'testing', 'story')).toBeDefined();
		expect(operatorMove(flow, 'testing', 'review', 'story')).toBeDefined();
		expect(operatorMove(flow, 'ready', 'testing', 'story')).toBeUndefined();
		expect(operatorMove(flow, 'review', 'cancelled', 'story')).toBeUndefined();
	});
});
