// A state's kind, colour, and moves, read from the lanes and transitions board.get answers
// (ADR-0301 §2, §3; S-0758), so that no view compares a status with a state's name.
import { laneOf, laneTitle, type Flow, type LaneKind, type Move } from './lanes';

/** The kind of the board's lane of that name; none for a state the board has no lane for. */
export function kindOf(flow: Flow, state: string | undefined): LaneKind | undefined {
	return state === undefined ? undefined : laneOf(flow, state)?.kind;
}

/** Whether the state is a lane of kind queue, waiting to be pulled. */
export function isQueue(flow: Flow, state: string | undefined): boolean {
	return kindOf(flow, state) === 'queue';
}

/** Whether the state is a lane of kind active, being worked. */
export function isActive(flow: Flow, state: string | undefined): boolean {
	return kindOf(flow, state) === 'active';
}

/** Whether the state is a lane of kind waiting, waiting on review. */
export function isWaiting(flow: Flow, state: string | undefined): boolean {
	return kindOf(flow, state) === 'waiting';
}

/** Whether the state is a lane of kind accepted. */
export function isAccepted(flow: Flow, state: string | undefined): boolean {
	return kindOf(flow, state) === 'accepted';
}

/** Whether the state is a lane of kind cancelled. */
export function isCancelled(flow: Flow, state: string | undefined): boolean {
	return kindOf(flow, state) === 'cancelled';
}

/** Whether the state is closed: accepted or cancelled. */
export function isClosed(flow: Flow, state: string | undefined): boolean {
	return isAccepted(flow, state) || isCancelled(flow, state);
}

/** Whether the state is open: any state not closed, an unknown one included. */
export function isOpen(flow: Flow, state: string | undefined): boolean {
	return !isClosed(flow, state);
}

/** Whether the state is underway: being worked or waiting on review. */
export function isUnderway(flow: Flow, state: string | undefined): boolean {
	return isActive(flow, state) || isWaiting(flow, state);
}

/** Whether an agent may pull an item from the state: some move from it has the effect pull. */
export function isPullable(flow: Flow, state: string | undefined): boolean {
	return (
		state !== undefined &&
		flow.transitions.some((m) => m.effects.includes('pull') && m.from.includes(state))
	);
}

/** The name of the board's first lane of kind accepted. */
export function acceptedState(flow: Flow): string | undefined {
	return flow.lanes.find((l) => l.kind === 'accepted')?.name;
}

/** The name of the board's first lane of kind cancelled. */
export function cancelledState(flow: Flow): string | undefined {
	return flow.lanes.find((l) => l.kind === 'cancelled')?.name;
}

/** The heading of the state's lane: its title, else the state's name. */
export function stateTitle(flow: Flow, state: string | undefined): string {
	const lane = state === undefined ? undefined : laneOf(flow, state);
	return lane ? laneTitle(lane) : (state ?? '');
}

/** The first move, by anyone, that accepts an item of `type` in `from`. */
export function acceptMove(flow: Flow, from: string, type: string): Move | undefined {
	return flow.transitions.find(
		(m) => m.effects.includes('accept') && m.from.includes(from) && m.types.includes(type)
	);
}

/** The first move a person may make that sends an item of `type` in `from` back. */
export function sendBackMove(flow: Flow, from: string, type: string): Move | undefined {
	return flow.transitions.find(
		(m) =>
			m.operator &&
			m.effects.includes('send-back') &&
			m.from.includes(from) &&
			m.types.includes(type)
	);
}

/** One move the item page offers. */
export type ItemMove = { to: string; title: string; reason: boolean; effects: string[] };

/**
 * The moves the item page offers an item of `type` in `state`: those a person may make to a later
 * lane, then those that send it back, then those that cancel it, one per target; a plain move to an
 * earlier lane, and reopening, are left to the board.
 */
export function itemMoves(flow: Flow, state: string | undefined, type: string): ItemMove[] {
	if (state === undefined) return [];
	const at = (name: string) => flow.lanes.findIndex((l) => l.name === name);
	const from = at(state);
	const group = (m: Move) =>
		m.effects.includes('cancel') ? 2 : m.effects.includes('send-back') ? 1 : 0;
	const moves = new Map<string, Move>();
	for (const m of flow.transitions) {
		if (!m.operator || !m.from.includes(state) || !m.types.includes(type)) continue;
		if (moves.has(m.to) || (group(m) === 0 && at(m.to) <= from)) continue;
		moves.set(m.to, m);
	}
	return [...moves.values()]
		.sort((a, b) => group(a) - group(b) || at(a.to) - at(b.to))
		.map((m) => ({
			to: m.to,
			title: stateTitle(flow, m.to),
			reason: m.reason,
			effects: m.effects
		}));
}

/**
 * The colour slots in palette.ts's CATEGORICAL each kind's lanes take, in board order; with the
 * baseline they give backlog 6, ready 3, in-progress 0, review 7, done 5, and cancelled 1.
 */
export const KIND_SLOTS: Record<LaneKind, number[]> = {
	queue: [6, 3],
	active: [0],
	waiting: [7],
	accepted: [5],
	cancelled: [1]
};

const SPARE_SLOTS = [2, 4];

/**
 * The state's colour slot in palette.ts's CATEGORICAL: the nth lane of a kind takes that kind's nth
 * slot, a lane past them the next spare slot; none once the spares run out, or for an unknown state.
 */
export function stateSlot(flow: Flow, state: string | undefined): number | undefined {
	const seen: Partial<Record<LaneKind, number>> = {};
	let spare = 0;
	for (const lane of flow.lanes) {
		const n = seen[lane.kind] ?? 0;
		seen[lane.kind] = n + 1;
		const slot = n < KIND_SLOTS[lane.kind].length ? KIND_SLOTS[lane.kind][n] : SPARE_SLOTS[spare++];
		if (lane.name === state) return slot;
	}
	return undefined;
}
