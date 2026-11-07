import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/down: kluster down of the project's cluster started, flai's kluster.down under
 * the kluster host action (S-0580), then the kluster's status, as GET /api/kluster answers it, the
 * job running in it. Refused as POST /api/kluster/up is: 403 while the action is off, 400 with
 * flai's sentence by rule, 501 when the host flai is older than kluster.down.
 */
export const POST: RequestHandler = () =>
	respond(async () => {
		await repo().klusterDown();
		return repo().klusterStatus();
	});
