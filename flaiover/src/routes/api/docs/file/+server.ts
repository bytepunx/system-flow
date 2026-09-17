import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/docs/file?path=design/system/overview.md */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		return repo().docFile(path);
	});
