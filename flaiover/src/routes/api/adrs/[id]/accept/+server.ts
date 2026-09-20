import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST: accept a proposed ADR (flai's adr.accept on the host): status accepted, dated today, its
 * index row updated, committed on its own. From then on the ADR is immutable, in the editor too.
 * 422 { error, findings } when the check refuses it.
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data, warnings } = await repo().write<Record<string, unknown>>('adr.accept', {
			id: params.id
		});
		return { ...data, log: warnings };
	});
