// Reordering cards within a board column (S-0057). Only the lanes the workflow orders have a pull
// order, backlog and ready with the baseline (S-0753), and only stories are in it; the rules are
// flai's (design/system/workflow.md), and flai refuses what these helpers should never ask for.
// These decide which placement a drop, a button, or a key asks flai order for, and when there is
// nothing to ask.
import type { Lane } from './lanes';

export type Placement = { before: string } | { after: string } | { top: true } | { bottom: true };

type Orderable = { id: string; type: string; status: string };

/** Whether a card has a place in the pull order of an ordered lane that the designer may change. */
export function reorderable(
	card: Orderable | undefined,
	writable: boolean,
	lanes: Pick<Lane, 'name' | 'ordered'>[]
): boolean {
	return (
		!!card &&
		writable &&
		card.type === 'story' &&
		lanes.some((l) => l.ordered === true && l.name === card.status)
	);
}

/** Whether `dragged` may be dropped on or beside `target` to reorder: same ordered column, both stories. */
export function canReorderOnto(
	dragged: Orderable | undefined,
	target: Orderable | undefined,
	writable: boolean,
	lanes: Pick<Lane, 'name' | 'ordered'>[]
): boolean {
	return (
		reorderable(dragged, writable, lanes) &&
		reorderable(target, writable, lanes) &&
		dragged!.id !== target!.id &&
		dragged!.status === target!.status
	);
}

/**
 * The placement for dropping `id` above or below `target`, given the column's stories in
 * order. Null when the card would stay where it is, so that nothing is written and nothing
 * claims to have changed.
 */
export function dropPlacement(
	stories: string[],
	id: string,
	target: string,
	half: 'above' | 'below'
): Placement | null {
	const from = stories.indexOf(id);
	const onto = stories.indexOf(target);
	if (from < 0 || onto < 0 || id === target) return null;
	if (half === 'above' && onto === from + 1) return null;
	if (half === 'below' && onto === from - 1) return null;
	return half === 'above' ? { before: target } : { after: target };
}

/** The placement for dropping `id` in the column's empty space: last, unless it already is. */
export function endPlacement(stories: string[], id: string): Placement | null {
	if (!stories.includes(id) || stories.at(-1) === id) return null;
	return { bottom: true };
}

/** One step up or down for the buttons and the keyboard; null at the end of the column. */
export function stepPlacement(stories: string[], id: string, dir: 'up' | 'down'): Placement | null {
	const at = stories.indexOf(id);
	if (at < 0) return null;
	if (dir === 'up') return at > 0 ? { before: stories[at - 1] } : null;
	return at < stories.length - 1 ? { after: stories[at + 1] } : null;
}

/** Alt+ArrowUp and Alt+ArrowDown on a focused card; anything else is not ours. */
export function keyStep(e: {
	key: string;
	altKey: boolean;
	ctrlKey: boolean;
	metaKey: boolean;
	shiftKey: boolean;
}): 'up' | 'down' | null {
	if (!e.altKey || e.ctrlKey || e.metaKey || e.shiftKey) return null;
	return e.key === 'ArrowUp' ? 'up' : e.key === 'ArrowDown' ? 'down' : null;
}
