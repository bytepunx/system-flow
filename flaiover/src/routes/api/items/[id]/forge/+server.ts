import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: a story's pull request, its required checks and its reviews (S-0394, ADR-0167), flai's
 * forge.status, a read of what flai serve last read from the forge and kept: `{ story,
 * pull_request }`, the pull request null when flai serve has read none. Nothing here calls the
 * forge. A refusal, such as an ID that is not a story's, is answered with flai's status and words.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => (await repo().run('forge.status', { id: params.id })).data);
