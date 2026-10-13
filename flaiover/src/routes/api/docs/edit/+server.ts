import { repo, RepoError, storyParam } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET ?path=&story=: a document as the editor loads it (flai's doc.show on the host: content, hash,
 * mode), from the story's worktree when story names one (S-0378), else the main checkout.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const path = url.searchParams.get('path');
		if (!path) throw new RepoError(400, 'path is required');
		const story = storyParam(url.searchParams.get('story'));
		return (await repo().run<Record<string, unknown>>('doc.show', { path, story })).data;
	});
