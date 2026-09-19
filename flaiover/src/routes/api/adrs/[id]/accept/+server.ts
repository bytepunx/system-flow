import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { _acceptArgs } from '../../+server';
import type { RequestHandler } from './$types';

/**
 * POST: accept a proposed ADR (`flai adr accept`): status accepted, dated today, its index row
 * updated, committed on its own. From then on the ADR is immutable, in the editor too.
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => {
		try {
			const { data, warnings } = await flai<Record<string, unknown>>(
				repo().root,
				_acceptArgs(params.id),
				{ exitStatus: { 4: 422 } }
			);
			return { ...data, log: warnings };
		} catch (e) {
			if (e instanceof RepoError && e.data?.refused)
				throw new RepoError(e.status, e.message, e.data.refused as Record<string, unknown>);
			throw e;
		}
	});
