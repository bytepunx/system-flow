import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { Issue } from '$lib/issues';
import type { RequestHandler } from './$types';

/**
 * GET ?story=S-nnnn: the project's issues, closed ones too, each with the stories its occurrences
 * name and the open story that links it (flai's issue.list, S-0198). With story, only the issues
 * that story recorded, read from its worktree while it has one.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const story = url.searchParams.get('story');
		return (await repo().run<Issue[]>('issue.list', story ? { story } : {})).data;
	});
