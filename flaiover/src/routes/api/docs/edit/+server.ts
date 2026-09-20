import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET ?path=: a document as the editor loads it (flai's doc.show on the host: content, hash, mode). */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		return (await repo().run<Record<string, unknown>>('doc.show', { path })).data;
	});
