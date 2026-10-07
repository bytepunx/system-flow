import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { VerifyReport } from '$lib/review';
import type { RequestHandler } from './$types';

/**
 * GET: a story's last verification (S-0270), flai's verify.status, a read of what flai verify
 * stored for it: `{ report }`, the report, or null when the story has none or no flai is connected.
 * The banner says flai is missing; here it only means nothing is known. Nothing runs it from here.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		try {
			const { data } = await repo().run<VerifyReport | null>('verify.status', { id: params.id });
			return { report: data ?? null };
		} catch (e) {
			if (e instanceof RepoError && e.status === 503) return { report: null };
			throw e;
		}
	});
