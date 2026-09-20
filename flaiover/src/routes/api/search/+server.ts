import { repo } from '$lib/server/repo';
import { search } from '$lib/server/search';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/search?q=term&docs=true */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const q = url.searchParams.get('q') ?? '';
		return search(repo(), q, url.searchParams.get('docs') === 'true');
	});
