import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: what the story's branch changes against main, files and hunks (flai stream diff, S-0041). */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data } = await flai<Record<string, unknown>>(repo().root, [
			'stream',
			'diff',
			params.id
		]);
		return data;
	});
