import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST: confirm a thread's pending recommendation as its answer, as the designer (flai's thread.confirm on the host, ADR-0090). */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data, warnings } = await repo().write<Record<string, unknown>>('thread.confirm', {
			id: params.id
		});
		return { ...data, warnings };
	});
