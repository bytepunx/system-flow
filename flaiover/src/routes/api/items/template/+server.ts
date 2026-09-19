import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET ?type=epic|story: the body the project's item template gives that type, below the heading
 * (`flai <type> new --print-body`). The "new" form starts from it, so a project that changed its
 * template gets its own sections (S-0059).
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const type = url.searchParams.get('type');
		if (type !== 'epic' && type !== 'story') throw new RepoError(400, 'type must be epic or story');
		const { data } = await flai<{ type: string; body: string }>(repo().root, [
			type,
			'new',
			'--print-body'
		]);
		return data;
	});
