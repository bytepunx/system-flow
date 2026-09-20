import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: what accepting the item would do (flai's accept.preview on the host: the branch, the release
 * plan, blockers, uncommitted files), so that the confirmation can show it first (S-0046).
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data, warnings } = await repo().run<Record<string, unknown>>(
			'accept.preview',
			{ id: params.id },
			{ timeoutMs: 60000 }
		);
		return { ...data, warnings };
	});
