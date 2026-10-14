import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { path }: restore main's change to one path the story's branch undoes (flai's
 * stream.restore on the host, S-0388), as flai accept --dry-run lists it under undoes_main.
 * 400 for a missing or blank path; flai's refusal (the story not in review, a worktree not
 * clean, a path its record does not hold) is passed on in its words (ADR-0016).
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		if (typeof body.path !== 'string' || !body.path.trim())
			throw new RepoError(400, 'path is required: one the branch undoes of main, to restore');
		return (await repo().write('stream.restore', { id: params.id, paths: [body.path] })).data;
	});
