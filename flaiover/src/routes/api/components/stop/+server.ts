import { repo, subProjectParam } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/components/stop { project? }: a sub-project's component dev server stopped when one
 * runs, flai's components.stop under the components host action (S-0439), then every dev server's
 * status, as GET /api/components answers it under status. 403 { error, action, enable } while the
 * action is off; 501 when the host flai is older than components.stop.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as { project?: unknown };
		await repo().componentsStop(subProjectParam(body.project));
		return repo().componentsStatus();
	});
