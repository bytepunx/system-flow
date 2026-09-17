// The board: columns from board.md limits and order, cards from the reader.
import { readFile } from 'node:fs/promises';
import { splitFrontMatter, type Item, type Repo } from './repo';
import { flaiBinary } from './flai';

export const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'] as const;

export type Card = {
	id: string;
	type: Item['type'];
	title: string;
	nature: string;
	parent?: string;
	status: string;
	blocked: boolean;
	entered_at: string;
	age_seconds: number;
};

export type Board = {
	wip_limits: Record<string, number>;
	order: string[];
	writable: boolean;
	columns: Record<string, Card[]>;
};

export function enteredAt(it: Item): string {
	return it.transitions.at(-1)?.at ?? it.created;
}

export async function board(repo: Repo, now = new Date()): Promise<Board> {
	const layout = await repo.layout();
	let limits: Record<string, number> = { ready: 5, 'in-progress': 2, review: 3 };
	let order: string[] = [];
	try {
		const raw = await readFile(repo.resolveInside(`${layout.wip}/kanban/board.md`), 'utf8');
		const { frontMatter } = splitFrontMatter(raw);
		const fm = (frontMatter ?? {}) as { wip_limits?: Record<string, number>; order?: string[] };
		if (fm.wip_limits) limits = { ...limits, ...fm.wip_limits };
		if (Array.isArray(fm.order)) order = fm.order.map(String);
	} catch {
		// no board.md: defaults
	}
	const columns: Record<string, Card[]> = Object.fromEntries(STATES.map((s) => [s, []]));
	for (const it of await repo.items()) {
		if (it.archived) continue;
		const entered = enteredAt(it);
		(columns[it.status] ??= []).push({
			id: it.id,
			type: it.type,
			title: it.title,
			nature: it.nature,
			parent: it.parent,
			status: it.status,
			blocked: (it.blocked ?? []).some((b) => !b.until),
			entered_at: entered,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(entered)) / 1000))
		});
	}
	return { wip_limits: limits, order, writable: (await flaiBinary()) !== null, columns };
}
