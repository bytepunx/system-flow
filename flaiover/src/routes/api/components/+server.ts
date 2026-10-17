import { repo, storyParam } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/components?story=S-nnnn: the component catalogue and the dev servers, flai's
 * components.list and components.status on the host (S-0439), asked together and afresh each time,
 * as { catalogue: { projects }, status: { enabled, enable, idle_seconds?, servers } }. Asking keeps a
 * running dev server from being stopped as idle. With story, the catalogue in that open story's
 * worktree and its dev servers (S-0442): 400 when the ID is malformed or names no open story with a
 * worktree. 501 when the host flai is older than either; 400 when flai cannot read the props, such
 * as without node or svelte.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const story = storyParam(url.searchParams.get('story'));
		const [catalogue, status] = await Promise.all([
			repo().components(story),
			repo().componentsStatus(story)
		]);
		return { catalogue, status };
	});
