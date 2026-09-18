import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET: what accepting this item would do (flai accept --dry-run): release plan and the branch to merge. Nothing changes. */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data, warnings } = await flai<Record<string, unknown>>(repo().root, [
			'accept',
			params.id,
			'--dry-run',
			'--yes'
		]);
		return { ...data, warnings };
	});
