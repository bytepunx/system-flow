// The board's lanes and the moves between them as the workflow defines them and flai's board.get
// answers them (ADR-0301 §2, §3; S-0753), and what each lane's menu offers (S-0167): every rule
// here reads the lanes and transitions it is given, and holds no list of states of its own.

/** What a lane is for (ADR-0301 §2). */
export type LaneKind = 'queue' | 'active' | 'waiting' | 'accepted' | 'cancelled';

/** One column of the board as board.get gives it; a flag or limit it does not have is absent. */
export type Lane = {
	name: string;
	title?: string;
	kind: LaneKind;
	wip_limit?: number;
	takes_limit?: boolean;
	ordered?: boolean;
	initial?: boolean;
};

/** One of the workflow's transitions as board.get gives it (internal/workitem BoardMove). */
export type Move = {
	from: string[];
	to: string;
	types: string[];
	by: string[];
	effects: string[];
	/** flai refuses the move without a reason. */
	reason: boolean;
	/** A person may make the move from the board. */
	operator: boolean;
};

/** The board's lanes in order and the workflow's moves in order. */
export type Flow = { lanes: Lane[]; transitions: Move[] };

/** The lane's heading: its title, else its name. */
export function laneTitle(lane: Lane): string {
	return lane.title || lane.name;
}

/** The board's lane of that name, if it has one. */
export function laneOf(flow: Flow, name: string): Lane | undefined {
	return flow.lanes.find((l) => l.name === name);
}

const at = (flow: Flow, name: string) => flow.lanes.findIndex((l) => l.name === name);

/**
 * The first move, in the workflow's order, that a person may make from the board taking an item of
 * `type` from `from` to `to`; none when the board may not offer it.
 */
export function operatorMove(flow: Flow, from: string, to: string, type: string): Move | undefined {
	return flow.transitions.find(
		(m) => m.operator && m.to === to && m.from.includes(from) && m.types.includes(type)
	);
}

/** The first move a person may make that cancels an item of `type` in `from`, if any. */
export function cancelMove(flow: Flow, from: string, type: string): Move | undefined {
	return flow.transitions.find(
		(m) =>
			m.operator && m.effects.includes('cancel') && m.from.includes(from) && m.types.includes(type)
	);
}

/** The lanes, in board order, that a person may move an item of `type` to from `lane`. */
function targets(flow: Flow, lane: string, type: string): Lane[] {
	return flow.lanes.filter((l) => l.name !== lane && operatorMove(flow, lane, l.name, type));
}

/**
 * Where the lane menu moves stories forward to: the nearest later lane of kind queue that a person
 * may move them to. Forward is within the queue only, which is the operator's to commit; past it a
 * story moves by its agent's pull and work, or by acceptance.
 */
export function forwardOf(flow: Flow, lane: string): string | null {
	const from = at(flow, lane);
	const later = targets(flow, lane, 'story').filter((l) => at(flow, l.name) > from);
	return later.find((l) => l.kind === 'queue')?.name ?? null;
}

/** Where the lane menu moves stories back to: the nearest earlier lane a person may move them to. */
export function backOf(flow: Flow, lane: string): string | null {
	const from = at(flow, lane);
	return targets(flow, lane, 'story').findLast((l) => at(flow, l.name) < from)?.name ?? null;
}

/** Whether flai refuses moving an item of `type` from `from` to `to` without a reason. */
export function needsReason(flow: Flow, from: string, to: string, type = 'story'): boolean {
	return operatorMove(flow, from, to, type)?.reason === true;
}

/** Whether the lane carries a WIP limit that the menu can change. */
export function hasLimit(flow: Flow, lane: string): boolean {
	return laneOf(flow, lane)?.takes_limit === true;
}

/** The lane flai makes a new item in: the workflow's initial lane. */
export function initialLane(flow: Flow): string | undefined {
	return (flow.lanes.find((l) => l.initial) ?? flow.lanes[0])?.name;
}

/**
 * The way a new item of `type` reaches each lane it may start in, from the initial lane: by moves a
 * person may make, each to a later lane of kind queue or active, so never into a lane that waits,
 * accepts, or cancels. The fewest moves, the nearer lane first where two are as few.
 */
function startPaths(flow: Flow, type: string): Map<string, string[]> {
	const first = initialLane(flow);
	const paths = new Map<string, string[]>();
	if (first === undefined) return paths;
	paths.set(first, []);
	const queue = [first];
	for (let lane = queue.shift(); lane !== undefined; lane = queue.shift()) {
		const from = at(flow, lane);
		for (const next of targets(flow, lane, type)) {
			if (paths.has(next.name) || at(flow, next.name) < from) continue;
			if (next.kind !== 'queue' && next.kind !== 'active') continue;
			paths.set(next.name, [...paths.get(lane)!, next.name]);
			queue.push(next.name);
		}
	}
	return paths;
}

/** The lanes an item of `type` can be created in, in board order; the initial lane first. */
export function startLanes(flow: Flow, type = 'story'): string[] {
	const paths = startPaths(flow, type);
	return flow.lanes.filter((l) => paths.has(l.name)).map((l) => l.name);
}

/** The lane a new item starts in, from the lane its menu was opened on: any other starts initial. */
export function startLane(flow: Flow, lane: string | null | undefined, type = 'story'): string {
	return startLanes(flow, type).includes(lane ?? '') ? lane! : (initialLane(flow) ?? '');
}

/** The moves that take a new item, made in the initial lane, to its starting lane, in order. */
export function movesTo(flow: Flow, lane: string, type = 'story'): string[] {
	return startPaths(flow, type).get(lane) ?? [];
}

/**
 * What a lane's menu does: create an item there, move stories a column, change the limit, or archive
 * every card it holds.
 */
export type LaneAction = 'create' | 'forward' | 'back' | 'limit' | 'archive';

/**
 * The lane's menu, in order (S-0167), for a lane holding `count` cards; a card's menu offers it too,
 * under the lane (S-0202).
 */
export function laneEntries(
	flow: Flow,
	lane: string,
	count = 0
): { action: LaneAction; label: string }[] {
	const forward = forwardOf(flow, lane);
	const backTo = backOf(flow, lane);
	return [
		{ action: 'create' as const, label: 'Create item here' },
		...(forward
			? [{ action: 'forward' as const, label: `Move stories forward to ${forward}…` }]
			: []),
		...(backTo ? [{ action: 'back' as const, label: `Move stories back to ${backTo}…` }] : []),
		...(hasLimit(flow, lane) ? [{ action: 'limit' as const, label: 'Change WIP limit…' }] : []),
		// cancelled cards stay until archived, and flai check fails on them (S-0343)
		...(laneOf(flow, lane)?.kind === 'cancelled' && count > 0
			? [{ action: 'archive' as const, label: 'Archive All' }]
			: [])
	];
}

/**
 * The IDs in the order Archive All archives them, one per call (S-0343): stories first, since each
 * takes its tasks along, then tasks, then epics, which flai refuses while a story of theirs is on
 * the board; each type in the order the lane holds it.
 */
export function archiveOrder(cards: { id: string; type: string }[]): string[] {
	const rank = (type: string) => (type === 'story' ? 0 : type === 'epic' ? 2 : 1);
	return [...cards].sort((a, b) => rank(a.type) - rank(b.type)).map((c) => c.id);
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
