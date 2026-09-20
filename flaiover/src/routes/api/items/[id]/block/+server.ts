import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST { reason }: open a blocked interval (flai's item.block on the host, S-0075). */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { reason?: string };
		if (!body.reason?.trim()) throw new RepoError(400, 'reason is required');
		return (await repo().write('item.block', { id: params.id, reason: body.reason })).data;
	});
