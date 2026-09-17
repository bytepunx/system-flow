import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { reason?: string };
		if (!body.reason?.trim()) throw new RepoError(400, 'reason is required');
		return (await flai(repo().root, ['block', params.id, '--reason', body.reason])).data;
	});
