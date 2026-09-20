// The board, as flai on the host lays it out (S-0073).
import type { Item, Repo } from './repo';
import { flaiBinary } from './flai';

export const STATES = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'] as const;

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

/** A card as flai's board.get gives it (internal/workitem BoardCard). */
type FlaiCard = Omit<Card, 'age_seconds'> & { age_in_column_seconds: number };
type FlaiBoard = {
	columns: Record<string, FlaiCard[] | null>;
	wip_limits: Record<string, number>;
	order: string[] | null;
};

/**
 * The board is flai's: its columns, limits, pull order, and the order of backlog and ready come from
 * board.get over the channel (S-0073), the same code as flai board. Age is counted from entered_at
 * here so that `now` means the same thing it always did to callers and tests.
 */
export async function board(repo: Repo, now = new Date()): Promise<Board> {
	const b = await repo.ask<FlaiBoard>('board.get', { all: true });
	const columns: Record<string, Card[]> = Object.fromEntries(STATES.map((s) => [s, []]));
	for (const [state, cards] of Object.entries(b.columns ?? {})) {
		columns[state] = (cards ?? []).map((c) => ({
			id: c.id,
			type: c.type,
			title: c.title,
			nature: c.nature,
			parent: c.parent || undefined,
			parent_title: c.parent_title || undefined,
			status: c.status,
			blocked: c.blocked,
			entered_at: c.entered_at,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(c.entered_at)) / 1000))
		}));
	}
	return {
		wip_limits: b.wip_limits,
		order: b.order ?? [],
		writable: (await flaiBinary()) !== null,
		columns
	};
}
