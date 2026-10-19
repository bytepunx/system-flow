import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST: enqueue the story in the merge queue (flai's queue.add on the host, S-0399, ADR-0173), as
 * the project's owner, which flai records with when. DELETE: take it out (queue.remove), dropping
 * the batch that holds it. Each answers flai's change: { action, story, changed, by, at, state }.
 * A refusal, a story not in review or not eligible, or one not enqueued, is 422 with flai's
 * sentence as text and, for a story not eligible, each reason as reasons, { kind, text }, beside
 * error, to show verbatim (ADR-0016). It is 403 with what enables it while the operator has not
 * (flai serve enable merge-queue).
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => (await repo().write('queue.add', { id: params.id })).data);

export const DELETE: RequestHandler = ({ params }) =>
	respond(async () => (await repo().write('queue.remove', { id: params.id })).data);
