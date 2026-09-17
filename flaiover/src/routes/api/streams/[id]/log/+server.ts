import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { entry?: string };
		if (!body.entry?.trim()) throw new RepoError(400, 'entry is required');
		return (await flai(repo().root, ['stream', 'log', params.id, body.entry])).data;
	});
