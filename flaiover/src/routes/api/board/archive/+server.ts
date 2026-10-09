import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { ids }: archive the done or cancelled items named, from the board's Archive and Archive All
 * menu entries (flai's board.archive on the host, S-0343). flai refuses an empty list and an item
 * still open.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		const { data } = await repo().write<Record<string, unknown>>('board.archive', {
			ids: Array.isArray(body.ids) ? body.ids.filter((id) => typeof id === 'string') : undefined
		});
		return data;
	});
