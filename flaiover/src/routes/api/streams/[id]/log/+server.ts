import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST { entry }: append to the story's narrative log (flai's stream.log on the host). */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { entry?: string };
		if (!body.entry?.trim()) throw new RepoError(400, 'entry is required');
		return (await repo().write('stream.log', { id: params.id, entry: body.entry })).data;
	});
