import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/up: kluster up of the project's cluster started, flai's kluster.up under the
 * kluster host action (S-0580), then the kluster's status, as GET /api/kluster answers it, the job
 * running in it. 403 { error, action, enable } while the action is off; 400 with flai's sentence
 * when it refuses by rule: no kluster: block, a config it cannot read, or an up or down that runs
 * already; 501 when the host flai is older than kluster.up.
 */
export const POST: RequestHandler = () =>
	respond(async () => {
		await repo().klusterUp();
		return repo().klusterStatus();
	});
