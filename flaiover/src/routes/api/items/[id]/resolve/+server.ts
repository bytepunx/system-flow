import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { ResolveStarted, ResolveStatus } from '$lib/review';
import type { RequestHandler } from './$types';

/**
 * GET: the story's last close-out or Rebase and resolve job (S-0804), flai's resolve.status, a read
 * of what flai serve records of it under .flai-cache/verify/jobs: `{ id, job? }`, the job running or
 * ended with its step, tier, steps, and log, and whether it runs, or the story alone when it has had
 * none. The review page follows a running one with it. It runs nothing.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => (await repo().run<ResolveStatus>('resolve.status', { id: params.id })).data);

/**
 * POST: ask flai serve to rebase a story in review that cannot be accepted onto the main branch and
 * close it out again, with no agent (flai's accept.resolve on the host, S-0804). It accepts nothing
 * and answers at once with the job, `{ id, job, started, running, log, reason }`: the job's ID,
 * whether this request started it or one of the story ran already, where its output is, and what
 * it does. The job sends the story back to its agent only when the rebase stops on conflicts or a
 * step of the close-out fails. flai's refusal (nothing holds the acceptance, the story not in
 * review, flai serve not taking the request) is passed on in its words with what it sent,
 * `{ refused: { id, reason } }` (ADR-0016).
 */
export const POST: RequestHandler = ({ params }) =>
	respond(
		async () => (await repo().write<ResolveStarted>('accept.resolve', { id: params.id })).data
	);
