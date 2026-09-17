import { repo, RepoError } from '$lib/server/repo';
import { stats } from '$lib/server/stats';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/stats?since=30d&type=story&by=nature — flai stats --json */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const r = repo();
		await r.watch();
		const q = {
			since: url.searchParams.get('since') ?? undefined,
			type: url.searchParams.get('type') ?? undefined,
			by: url.searchParams.get('by') ?? undefined
		};
		try {
			return await stats(r, q);
		} catch (e) {
			if (e instanceof RepoError) throw e;
			throw new RepoError(400, e instanceof Error ? e.message : String(e));
		}
	});
