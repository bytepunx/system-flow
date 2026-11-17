// Which items being worked or waiting on review say they are changing a document
// (ADR-0019, S-0037). Shared by the explorer and the editor.
import type { Flow } from './lanes';
import { isUnderway } from './workflow';

export type Worker = { id: string; title: string; status?: string; touches?: string[] };

/**
 * Items underway in the flow, in a lane of kind active or waiting, or of no status given, whose
 * `touches` cover the path: the path itself or a folder above it.
 */
export function touching(
	path: string,
	items: Worker[],
	flow: Flow = { lanes: [], transitions: [] }
): Worker[] {
	if (!path) return [];
	return items.filter(
		(w) =>
			(w.status === undefined || isUnderway(flow, w.status)) &&
			(w.touches ?? []).some((t) => {
				const p = t.replace(/\/$/, '');
				return path === p || path.startsWith(p + '/');
			})
	);
}
