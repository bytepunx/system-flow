// Test support, never imported by the application: the lanes and transitions flai's board.get
// answers for the baseline workflow (flai/internal/workflow/baseline/board.md), as
// TestBoardViewLanesAreTheWorkflowsColumns and TestBoardViewTransitionsAreTheWorkflows in
// flai/internal/workitem give them, and for the workflow those tests also run: in-progress
// renamed doing, and testing, an active lane with a WIP limit, added between it and review; for
// the baseline with that testing lane added after in-progress, which keeps its name (S-0758); and
// for the baseline with its pullable queue, accepted, and cancelled lanes renamed queued, shipped,
// and dropped, so that a view that compares a state's name with ready, done, or cancelled shows it
// (S-0838, I-0399).
import type { Flow, Lane, Move } from './lanes';

const all = ['epic', 'story', 'task'];

/** A move by anyone, for every type, with no effect and no reason, unless given. */
const move = (from: string[], to: string, more: Partial<Move> = {}): Move => ({
	from,
	to,
	types: all,
	by: ['anyone'],
	effects: [],
	reason: false,
	operator: true,
	...more
});

const baselineLanes: Lane[] = [
	{ name: 'backlog', title: 'backlog', kind: 'queue', ordered: true, initial: true },
	{ name: 'ready', title: 'ready', kind: 'queue', wip_limit: 5, takes_limit: true, ordered: true },
	{ name: 'in-progress', title: 'in-progress', kind: 'active', wip_limit: 2, takes_limit: true },
	{ name: 'review', title: 'review', kind: 'waiting', wip_limit: 3, takes_limit: true },
	{ name: 'done', title: 'done', kind: 'accepted' },
	{ name: 'cancelled', title: 'cancelled', kind: 'cancelled' }
];

/** The baseline workflow's board, as flai answers it for a project with no workflow of its own. */
export const baselineFlow: Flow = {
	lanes: baselineLanes,
	transitions: [
		move(['backlog'], 'ready', { effects: ['commit'] }),
		move(['ready'], 'in-progress', { effects: ['pull'] }),
		move(['in-progress'], 'review'),
		move(['review'], 'done', { types: ['epic', 'story'], by: ['acceptance'], effects: ['accept'] }),
		move(['review'], 'done', { types: ['task'] }),
		move(['in-progress'], 'done', { types: ['task'] }),
		move(['review'], 'in-progress', { effects: ['send-back'], reason: true }),
		move(['in-progress'], 'ready'),
		move(['ready'], 'backlog'),
		move(['backlog', 'ready', 'in-progress'], 'cancelled', { effects: ['cancel'], reason: true }),
		move(['review'], 'cancelled', {
			by: ['cascade'],
			effects: ['cancel'],
			reason: true,
			operator: false
		}),
		move(['cancelled'], 'backlog', { effects: ['reopen'] })
	]
};

/** The workflow with in-progress renamed doing and testing added after it (S-0739). */
export const testingFlow: Flow = {
	lanes: [
		...baselineLanes.slice(0, 2),
		{ name: 'doing', title: 'doing', kind: 'active', wip_limit: 2, takes_limit: true },
		{ name: 'testing', title: 'testing', kind: 'active', wip_limit: 2, takes_limit: true },
		...baselineLanes.slice(3)
	],
	transitions: [
		move(['backlog'], 'ready', { effects: ['commit'] }),
		move(['ready'], 'doing', { effects: ['pull'] }),
		move(['doing'], 'testing'),
		move(['doing', 'testing'], 'review'),
		move(['review'], 'done', { types: ['epic', 'story'], by: ['acceptance'], effects: ['accept'] }),
		move(['review'], 'done', { types: ['task'] }),
		move(['doing', 'testing'], 'done', { types: ['task'] }),
		move(['review'], 'doing', { effects: ['send-back'], reason: true }),
		move(['testing'], 'doing'),
		move(['doing'], 'ready'),
		move(['ready'], 'backlog'),
		move(['backlog', 'ready', 'doing', 'testing'], 'cancelled', {
			effects: ['cancel'],
			reason: true
		}),
		move(['review'], 'cancelled', {
			by: ['cascade'],
			effects: ['cancel'],
			reason: true,
			operator: false
		}),
		move(['cancelled'], 'backlog', { effects: ['reopen'] })
	]
};

/** The baseline workflow with testing added after in-progress, which keeps its name (S-0758). */
export const addedStateFlow: Flow = {
	lanes: [
		...baselineLanes.slice(0, 3),
		{ name: 'testing', title: 'testing', kind: 'active', wip_limit: 2, takes_limit: true },
		...baselineLanes.slice(3)
	],
	transitions: [
		move(['backlog'], 'ready', { effects: ['commit'] }),
		move(['ready'], 'in-progress', { effects: ['pull'] }),
		move(['in-progress'], 'testing'),
		move(['in-progress', 'testing'], 'review'),
		move(['review'], 'done', { types: ['epic', 'story'], by: ['acceptance'], effects: ['accept'] }),
		move(['review'], 'done', { types: ['task'] }),
		move(['in-progress', 'testing'], 'done', { types: ['task'] }),
		move(['review'], 'in-progress', { effects: ['send-back'], reason: true }),
		move(['testing'], 'in-progress'),
		move(['in-progress'], 'ready'),
		move(['ready'], 'backlog'),
		move(['backlog', 'ready', 'in-progress', 'testing'], 'cancelled', {
			effects: ['cancel'],
			reason: true
		}),
		move(['review'], 'cancelled', {
			by: ['cascade'],
			effects: ['cancel'],
			reason: true,
			operator: false
		}),
		move(['cancelled'], 'backlog', { effects: ['reopen'] })
	]
};

/** The baseline workflow with ready renamed queued, done shipped, and cancelled dropped (S-0838). */
export const renamedFlow: Flow = {
	lanes: [
		baselineLanes[0],
		{ ...baselineLanes[1], name: 'queued', title: 'queued' },
		...baselineLanes.slice(2, 4),
		{ name: 'shipped', title: 'shipped', kind: 'accepted' },
		{ name: 'dropped', title: 'dropped', kind: 'cancelled' }
	],
	transitions: [
		move(['backlog'], 'queued', { effects: ['commit'] }),
		move(['queued'], 'in-progress', { effects: ['pull'] }),
		move(['in-progress'], 'review'),
		move(['review'], 'shipped', {
			types: ['epic', 'story'],
			by: ['acceptance'],
			effects: ['accept']
		}),
		move(['review'], 'shipped', { types: ['task'] }),
		move(['in-progress'], 'shipped', { types: ['task'] }),
		move(['review'], 'in-progress', { effects: ['send-back'], reason: true }),
		move(['in-progress'], 'queued'),
		move(['queued'], 'backlog'),
		move(['backlog', 'queued', 'in-progress'], 'dropped', { effects: ['cancel'], reason: true }),
		move(['review'], 'dropped', {
			by: ['cascade'],
			effects: ['cancel'],
			reason: true,
			operator: false
		}),
		move(['dropped'], 'backlog', { effects: ['reopen'] })
	]
};
