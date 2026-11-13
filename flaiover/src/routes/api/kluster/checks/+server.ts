import { klusterStoryParam, repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/checks { story }: the project's checks run against the story's namespace
 * started, flai's kluster.checks under the kluster host action (S-0663), then the deployments, as
 * GET /api/kluster/deployments answers them. A story that is not an ID is refused with 400 before
 * flai is asked. 403 { error, action, enable } while the action is off; 400 with flai's sentence
 * when it refuses by rule, such as a namespace not deployed at the story's head; 501 when the host
 * flai is older than kluster.checks.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as { story?: unknown };
		await repo().klusterChecks(klusterStoryParam(body.story));
		return repo().klusterDeployments();
	});
