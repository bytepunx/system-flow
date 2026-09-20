import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST { reason? }: resolve a thread as the designer (flai's thread.resolve on the host). */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { reason?: string };
		const { data, warnings } = await repo().write<Record<string, unknown>>('thread.resolve', {
			id: params.id,
			reason: body.reason
		});
		return { ...data, warnings };
	});
