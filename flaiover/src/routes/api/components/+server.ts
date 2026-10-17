import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/components: the component catalogue and the dev servers, flai's components.list and
 * components.status on the host (S-0439), asked together and afresh each time, as
 * { catalogue: { projects }, status: { enabled, enable, idle_seconds?, servers } }. Asking keeps a
 * running dev server from being stopped as idle. 501 when the host flai is older than either; 400
 * when flai cannot read the props, such as without node or svelte.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		const [catalogue, status] = await Promise.all([repo().components(), repo().componentsStatus()]);
		return { catalogue, status };
	});
