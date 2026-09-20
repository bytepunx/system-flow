import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET ?type=epic|story: the body the project's template gives a new item (flai's item.template). */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const type = url.searchParams.get('type') ?? 'story';
		if (type !== 'epic' && type !== 'story') throw new RepoError(400, 'type must be epic or story');
		return (await repo().run<{ type: string; body: string }>('item.template', { type })).data;
	});
