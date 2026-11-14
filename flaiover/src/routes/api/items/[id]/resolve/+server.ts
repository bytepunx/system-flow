import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST: send a story in review that cannot be accepted back to its agent to rebase and resolve
 * what holds it (flai's accept.resolve on the host, ADR-0302). It accepts nothing: flai moves the
 * story to in-progress and opens the acceptance's send-back thread (ADR-0248), and the answer is
 * that send back, { id, reason, sent_back, branch, base, conflicts, thread, blockers }. flai's
 * refusal (nothing holds the acceptance, the story not in review) is passed on in its words with
 * what it sent, { refused: { id, reason } } (ADR-0016).
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => (await repo().write('accept.resolve', { id: params.id })).data);
