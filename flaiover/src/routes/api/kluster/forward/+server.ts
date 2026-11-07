import { klusterForwardParam, repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/kluster/forward { namespace, service, port, on }: a port-forward of the service's port
 * started (on true) or stopped (on false), flai's kluster.forward.start or kluster.forward.stop under
 * the kluster host action (S-0586), then the deployments, as GET /api/kluster/deployments answers
 * them, with the namespace's forwards in its access. flai serve runs kubectl port-forward bound to
 * 0.0.0.0 on a free port of the host. A target that is not a namespace, a service, and a port is
 * refused with 400 before flai is asked. 403 { error, action, enable } while the action is off; 400
 * with flai's sentence when it refuses by rule; 501 when the host flai is older than these methods.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as Record<string, unknown>;
		const target = klusterForwardParam(body);
		if (body.on === false) await repo().klusterForwardStop(target);
		else await repo().klusterForwardStart(target);
		return repo().klusterDeployments();
	});
