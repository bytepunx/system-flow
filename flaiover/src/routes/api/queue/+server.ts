import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: the merge queue as board.get carries it (S-0398), or null when the flai on the host sends
 * none, and whether the merge-queue host action is on (project.info's host_actions, S-0399). Both
 * are asked of flai afresh each time, never remembered: the queue's state is in the host's
 * .flai-cache, whose changes no notice reports, so the board polls this while a batch is under
 * verification to show its step and result as they change.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		const [info, board] = await Promise.all([
			repo()
				.ask<{ host_actions?: Record<string, boolean> }>('project.info')
				.catch(() => null),
			repo().ask<{ queue?: unknown }>('board.get', {})
		]);
		return { enabled: info?.host_actions?.['merge-queue'] === true, queue: board.queue ?? null };
	});

/**
 * POST { held }: held true holds the merge queue, so that it builds no batch, and false resumes
 * it (flai's queue.hold on the host, S-0399), as the project's owner, which flai records with when.
 * It answers flai's change, { action, story, changed, by, at, state }; 400 for a held that is not
 * true or false, and 403 with what enables it while the operator has not (flai serve enable
 * merge-queue).
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => null)) as { held?: unknown } | null;
		if (typeof body?.held !== 'boolean')
			throw new RepoError(
				400,
				`held is true to hold the merge queue or false to resume it, not ${JSON.stringify(body?.held ?? null)}`
			);
		return (await repo().write('queue.hold', { held: body.held })).data;
	});
