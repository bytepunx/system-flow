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

/** What a lane's menu does: create an item there, move stories a column, or change the limit. */
export type LaneAction = 'create' | 'forward' | 'back' | 'limit';

/** The lane's menu, in order (S-0167); a card's menu offers it too, under the lane (S-0202). */
export function laneEntries(lane: string): { action: LaneAction; label: string }[] {
	const forward = forwardOf(lane);
	const backTo = backOf(lane);
	return [
		{ action: 'create' as const, label: 'Create item here' },
		...(forward
			? [{ action: 'forward' as const, label: `Move stories forward to ${forward}…` }]
			: []),
		...(backTo ? [{ action: 'back' as const, label: `Move stories back to ${backTo}…` }] : []),
		...(hasLimit(lane) ? [{ action: 'limit' as const, label: 'Change WIP limit…' }] : [])
	];
}

// What a lane holds, by type (S-0256): a short count below its title, the long form on hover.

/** How many epics, stories, and tasks a lane holds. */
export type LaneCounts = { epic: number; story: number; task: number };

/** How many epics, stories, and tasks the cards hold; other types are not counted. */
export function laneCounts(cards: { type: string }[]): LaneCounts {
	const n: LaneCounts = { epic: 0, story: 0, task: 0 };
	for (const c of cards) if (Object.hasOwn(n, c.type)) n[c.type as keyof LaneCounts]++;
	return n;
}

/** The short form, epics | stories | tasks: `1 | 2 | 4`. */
export function shortCounts(n: LaneCounts): string {
	return `${n.epic} | ${n.story} | ${n.task}`;
}

const many = (n: number, one: string, other: string) => `${n} ${n === 1 ? one : other}`;

/** The long form, each type pluralized: `1 epic, 2 stories, 0 tasks`. */
export function longCounts(n: LaneCounts): string {
	return [
		many(n.epic, 'epic', 'epics'),
		many(n.story, 'story', 'stories'),
		many(n.task, 'task', 'tasks')
	].join(', ');
}
