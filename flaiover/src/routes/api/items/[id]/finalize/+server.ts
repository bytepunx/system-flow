import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST: finalize a draft story as the designer (flai's item.finalize on the host, S-0201): the
 * draft flag is cleared and flai records who finalized it; the dashboard does not name anyone.
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => (await repo().write('item.finalize', { id: params.id })).data);
