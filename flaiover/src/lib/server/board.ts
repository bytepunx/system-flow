// The board: columns from board.md limits and order, cards from the reader.
import { readFile } from 'node:fs/promises';
import { splitFrontMatter, type Item, type Repo } from './repo';
import { flaiBinary } from './flai';

export const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'] as const;

/** The columns that have a pull order (S-0057); the rule is flai's, in design/system/workflow.md. */
export const ORDERED = ['backlog', 'ready'] as const;

const idNumber = (id: string) => Number(id.replace(/^\D+/, '')) || 0;

/**
 * A column's cards with its stories in pull order: those `order` names first, in its order,
 * then the rest by ID. Stories take the sequence in the places stories already hold, so epics
 * and tasks stay where they are. The same reading as flai board's.
 */
export function inPullOrder<T extends { id: string; type: string }>(
	cards: T[],
	order: string[]
): T[] {
	const rank = new Map<string, number>();
	for (const id of order) if (!rank.has(id)) rank.set(id, rank.size);
	const stories = cards
		.filter((c) => c.type === 'story')
		.sort((a, b) => {
			const ra = rank.get(a.id);
			const rb = rank.get(b.id);
			if (ra !== undefined && rb !== undefined) return ra - rb;
			if (ra !== undefined || rb !== undefined) return ra !== undefined ? -1 : 1;
			return idNumber(a.id) - idNumber(b.id) || a.id.localeCompare(b.id);
		});
	let next = 0;
	return cards.map((c) => (c.type === 'story' ? stories[next++] : c));
}

export type Card = {
	id: string;
	type: Item['type'];
	title: string;
	nature: string;
	parent?: string;
	/** Title of the parent item, so a card can name its epic or story without a second request. */
	parent_title?: string;
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
	const items = await repo.items();
	// Parents are looked up among every item, archived ones included.
	const titles = new Map(items.map((it) => [it.id, it.title]));
	for (const it of items) {
		if (it.archived) continue;
		const entered = enteredAt(it);
		(columns[it.status] ??= []).push({
			id: it.id,
			type: it.type,
			title: it.title,
			nature: it.nature,
			parent: it.parent,
			parent_title: it.parent ? titles.get(it.parent) : undefined,
			status: it.status,
			blocked: (it.blocked ?? []).some((b) => !b.until),
			entered_at: entered,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(entered)) / 1000))
		});
	}
	for (const state of ORDERED) columns[state] = inPullOrder(columns[state], order);
	return { wip_limits: limits, order, writable: (await flaiBinary()) !== null, columns };
}
