// Which in-progress or in-review items say they are changing a document
// (ADR-0019, S-0037). Shared by the explorer and the editor.

export type Worker = { id: string; title: string; status?: string; touches?: string[] };

/** Items in progress or review whose `touches` cover the path: the path itself or a folder above it. */
export function touching(path: string, items: Worker[]): Worker[] {
	if (!path) return [];
	return items.filter(
		(w) =>
			(w.status === undefined || w.status === 'in-progress' || w.status === 'review') &&
			(w.touches ?? []).some((t) => {
				const p = t.replace(/\/$/, '');
				return path === p || path.startsWith(p + '/');
			})
	);
}
