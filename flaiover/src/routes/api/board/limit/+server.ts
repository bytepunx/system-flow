import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { column, limit }: set a column's WIP limit in wip/kanban/board.md, from the board's lane
 * menu (flai's board.limit on the host, S-0167). 0 is no limit; flai refuses a column without one.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		const { data } = await repo().write<Record<string, unknown>>('board.limit', {
			column: typeof body.column === 'string' ? body.column : undefined,
			limit: typeof body.limit === 'number' ? body.limit : undefined
		});
		return data;
	});
