import { repo, storyParam, subProjectParam } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST /api/components/start { project?, story? }: a sub-project's component dev server started,
 * flai's components.start under the components host action (S-0439), then every dev server's status,
 * as GET /api/components answers it under status. project may be left out when one sub-project is
 * Svelte. With story, the dev server runs in that story's worktree and the status is its (S-0442);
 * a malformed ID is refused with 400 before flai is asked. 403 { error, action, enable } while the
 * action is off; 400 when flai refuses by rule, such as a sub-project without vite or a story with
 * no worktree; 500 with flai's message, Vite's last lines of output with it, when the dev server
 * fails to start; 501 when the host flai is older than components.start.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = ((await request.json().catch(() => null)) ?? {}) as {
			project?: unknown;
			story?: unknown;
		};
		const project = subProjectParam(body.project);
		const story = storyParam(body.story);
		await repo().componentsStart(project, story);
		return repo().componentsStatus(story);
	});
