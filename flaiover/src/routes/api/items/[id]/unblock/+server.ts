import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST: close the open blocked interval (flai's item.unblock on the host, S-0075). */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => (await repo().write('item.unblock', { id: params.id })).data);
