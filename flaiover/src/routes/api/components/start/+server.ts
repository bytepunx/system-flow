import { repo, subProjectParam } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/components/start { project? }: a sub-project's component dev server started, flai's
 * components.start under the components host action (S-0439), then every dev server's status, as
 * GET /api/components answers it under status. project may be left out when one sub-project is
 * Svelte. 403 { error, action, enable } while the action is off; 400 when flai refuses by rule, such
 * as a sub-project without vite; 500 with flai's message, Vite's last lines of output with it, when
 * the dev server fails to start; 501 when the host flai is older than components.start.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as { project?: unknown };
		await repo().componentsStart(subProjectParam(body.project));
		return repo().componentsStatus();
	});
