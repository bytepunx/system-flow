// What kind of file of the project changed (S-0161), from its repo-relative path and the manifest's
// layout: the server forgets only the answers read from that kind, and pages ask again only for the
// kinds they show. Pure, so the server and the pages share one reading of a path.

/** The manifest's three folders, repo-relative. */
export type ChangeLayout = { design: string; docs: string; wip: string };

/**
 * project: system-flow.yaml. item: a work item or the board's policy, under the wip folder's kanban
 * or archive. narrative: under its agents. thread: under its threads. adr: under the design folder's
 * adrs. document: any other file under the three folders. other: anything else, or any path while
 * the layout is not known, which is read as possibly affecting everything.
 */
export type ChangeKind = 'project' | 'item' | 'narrative' | 'thread' | 'adr' | 'document' | 'other';

export type Change = { path: string; kind: ChangeKind };

export const MANIFEST = 'system-flow.yaml';

function trim(dir: string): string {
	return dir.replace(/^\.\//, '').replace(/\/+$/, '');
}

function under(path: string, dir: string): boolean {
	const d = trim(dir);
	return d !== '' && path.startsWith(d + '/');
}

/** The kind of the file at path, read against layout; null layout reads every path but the manifest as other. */
export function kindOf(path: string, layout: ChangeLayout | null): ChangeKind {
	if (path === MANIFEST) return 'project';
	if (!layout) return 'other';
	const wip = trim(layout.wip);
	if (under(path, `${wip}/kanban`) || under(path, `${wip}/archive`)) return 'item';
	if (under(path, `${wip}/agents`)) return 'narrative';
	if (under(path, `${wip}/threads`)) return 'thread';
	if (under(path, `${trim(layout.design)}/adrs`)) return 'adr';
	if (under(path, layout.design) || under(path, layout.docs) || under(path, wip)) return 'document';
	return 'other';
}

/** Whether any of changes is of one of kinds, or of a kind that may affect anything (project, other). */
export function affects(changes: Change[], ...kinds: ChangeKind[]): boolean {
	return changes.some((c) => c.kind === 'project' || c.kind === 'other' || kinds.includes(c.kind));
}
