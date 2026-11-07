import { klusterStoryParam, repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/deploy { story }: the story's deploy into its namespace started, flai's
 * kluster.deploy under the kluster host action (S-0581), then the deployments, as GET
 * /api/kluster/deployments answers them, the deploy running in them. A story that is not an ID is
 * refused with 400 before flai is asked. 403 { error, action, enable } while the action is off; 400
 * with flai's sentence when it refuses by rule, such as a job that runs already; 501 when the host
 * flai is older than kluster.deploy.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as { story?: unknown };
		await repo().klusterDeploy(klusterStoryParam(body.story));
		return repo().klusterDeployments();
	});
