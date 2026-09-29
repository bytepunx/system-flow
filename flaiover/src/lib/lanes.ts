// What each lane of the board offers in its menu (S-0167), by flai's workflow rules: a story
// moves forward from backlog to ready, and back one column from ready, in-progress, review, or
// cancelled; done is final (ADR-0055). Only ready, in-progress, and review carry a WIP limit.

export const LANES = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'] as const;
export type Lane = (typeof LANES)[number];

/** The lanes an item can be created in; any other starts in backlog. */
export const START_LANES = ['backlog', 'ready', 'in-progress'] as const;
export type StartLane = (typeof START_LANES)[number];

/** The lane a new item starts in, from the lane its menu was opened on. */
export function startLane(lane: string | null | undefined): StartLane {
	return (START_LANES as readonly string[]).includes(lane ?? '') ? (lane as StartLane) : 'backlog';
}

/** The moves that take a new item, made in backlog, to its starting lane, in order. */
export function movesTo(lane: StartLane): Lane[] {
	return START_LANES.slice(1, START_LANES.indexOf(lane) + 1);
}

/** Where the lane menu moves stories forward to: from backlog only. */
export function forwardOf(lane: string): Lane | null {
	return lane === 'backlog' ? 'ready' : null;
}

const back: Record<string, Lane> = {
	ready: 'backlog',
	'in-progress': 'ready',
	review: 'in-progress',
	cancelled: 'backlog'
};

/** Where the lane menu moves stories back to; none from backlog or done. */
export function backOf(lane: string): Lane | null {
	return back[lane] ?? null;
}

/** flai asks why a story is sent back from review (workflow.md). */
export function needsReason(from: string, to: string): boolean {
	return from === 'review' && to === 'in-progress';
}

export const LIMITED_LANES = ['ready', 'in-progress', 'review'] as const;

/** Whether the lane carries a WIP limit that the menu can change. */
export function hasLimit(lane: string): boolean {
	return (LIMITED_LANES as readonly string[]).includes(lane);
}
