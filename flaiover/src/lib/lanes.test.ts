import { describe, expect, it } from 'vitest';
import {
	archiveOrder,
	backOf,
	forwardOf,
	hasLimit,
	laneCounts,
	laneEntries,
	LANES,
	longCounts,
	movesTo,
	needsReason,
	shortCounts,
	startLane
} from './lanes';

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

	// S-0202: the lane's entries, shared by the lane's menu and a card's
	it("lists each lane's menu entries in order, with the column a move goes to", () => {
		expect(laneEntries('backlog')).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'forward', label: 'Move stories forward to ready…' }
		]);
		expect(laneEntries('review')).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'back', label: 'Move stories back to in-progress…' },
			{ action: 'limit', label: 'Change WIP limit…' }
		]);
		expect(laneEntries('done').map((e) => e.action)).toEqual(['create']);
		expect(laneEntries('cancelled').map((e) => e.action)).toEqual(['create', 'back']);
	});

	// S-0343: cancelled cards stay until archived
	it('offers Archive All on the cancelled lane only, and only while it holds cards', () => {
		expect(laneEntries('cancelled', 3)).toEqual([
			{ action: 'create', label: 'Create item here' },
			{ action: 'back', label: 'Move stories back to backlog…' },
			{ action: 'archive', label: 'Archive All' }
		]);
		expect(laneEntries('cancelled', 0).map((e) => e.action)).not.toContain('archive');
		for (const lane of LANES.filter((l) => l !== 'cancelled'))
			expect(laneEntries(lane, 3).map((e) => e.action)).not.toContain('archive');
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
