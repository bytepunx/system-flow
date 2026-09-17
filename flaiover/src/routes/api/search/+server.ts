import { repo } from '$lib/server/repo';
import { searchIndex } from '$lib/server/search';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/search?q=term&docs=true */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const r = repo();
		await r.watch();
		const idx = searchIndex(r);
		const q = url.searchParams.get('q') ?? '';
		const hits = await idx.search(q, url.searchParams.get('docs') === 'true');
		return { query: q, indexed: idx.size(), hits };
	});
