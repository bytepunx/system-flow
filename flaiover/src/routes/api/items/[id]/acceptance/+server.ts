import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: what accepting this item would do (flai accept --dry-run): release plan, the branch to merge,
 * blockers, and uncommitted paths outside wip. Nothing changes. It runs without --yes so it checks what
 * the move that follows checks (S-0051).
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data, warnings } = await flai<Record<string, unknown>>(repo().root, [
			'accept',
			params.id,
			'--dry-run'
		]);
		return { ...data, warnings };
	});
