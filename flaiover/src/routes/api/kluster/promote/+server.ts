import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/promote: main's head deployed into the default namespace, started, flai's
 * kluster.promote under the kluster host action (S-0581), then the deployments, as GET
 * /api/kluster/deployments answers them. 403 { error, action, enable } while the action is off; 400
 * with flai's sentence when it refuses by rule; 501 when the host flai is older than
 * kluster.promote.
 */
export const POST: RequestHandler = () =>
	respond(async () => {
		await repo().klusterPromote();
		return repo().klusterDeployments();
	});
