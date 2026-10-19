// The board, as flai on the host lays it out (S-0073).
import type { Item, Repo } from './repo';
import type { TaskSummary } from '$lib/taskplan';
import type { Share } from '$lib/activity';
import type { MergeQueueState } from '$lib/components/MergeQueue.svelte';

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
	/** A story an agent wrote that the operator has not yet finalized (S-0201); absent otherwise. */
	draft?: boolean;
	/** A done card flai keeps only because its release looks pending (S-0087); absent otherwise. */
	archived?: boolean;
	entered_at: string;
	age_seconds: number;
	/** A story's tasks by state and its plan's layers (S-0176); absent when it has no tasks. */
	tasks?: TaskSummary;
	/** The shares in force on a story started on or freed by them (S-0338); absent when none. */
	shared?: Share[];
	/** The paths an open story's branch undoes of main's change (ADR-0155, S-0388); absent when none. */
	undoes_main?: string[];
};

export type Board = {
	wip_limits: Record<string, number>;
	order: string[];
	writable: boolean;
	columns: Record<string, Card[]>;
	/** The merge queue beside the board (S-0398); absent from an older flai or when flai could not read it. */
	queue?: MergeQueueState;
};

/** A card as flai's board.get gives it (internal/workitem BoardCard). */
type FlaiCard = Omit<Card, 'age_seconds' | 'tasks' | 'shared' | 'undoes_main'> & {
	age_in_column_seconds: number;
	tasks?: TaskSummary | null;
	shared?: Share[] | null;
	undoes_main?: string[] | null;
};
type FlaiBoard = {
	columns: Record<string, FlaiCard[] | null>;
	wip_limits: Record<string, number>;
	order: string[] | null;
	queue?: MergeQueueState | null;
};

/**
 * The board is flai's: its columns, limits, pull order, and the order of backlog and ready come from
 * board.get over the channel (S-0073), the same code as flai board. Age is counted from entered_at
 * here so that `now` means the same thing it always did to callers and tests.
 */
export async function board(repo: Repo, now = new Date()): Promise<Board> {
	// Kept until flai says a file changed, like everything else asked: the board's one moving part,
	// a card's age, is counted here from entered_at.
	const b = await repo.boardView<FlaiBoard>();
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
			...(c.draft ? { draft: true } : {}),
			...(c.archived ? { archived: true } : {}),
			entered_at: c.entered_at,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(c.entered_at)) / 1000)),
			...(c.tasks ? { tasks: c.tasks } : {}),
			...(c.shared?.length ? { shared: c.shared } : {}),
			...(c.undoes_main?.length ? { undoes_main: c.undoes_main } : {})
		}));
	}
	return {
		wip_limits: b.wip_limits,
		order: b.order ?? [],
		// flai on the host does the writing (S-0075): a board that could be read can be written.
		writable: true,
		columns,
		...(b.queue ? { queue: b.queue } : {})
	};
}
