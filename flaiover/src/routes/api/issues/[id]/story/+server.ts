import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { epic? }: make a backlog story from an open issue (flai's issue.story on the host, S-0198),
 * as the manifest's owner, committed on its own. The issue is read from the main checkout and is not
 * changed. 422 { error, findings } when the check refuses the story.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { epic?: string };
		const { data, warnings } = await repo().write<Record<string, unknown>>('issue.story', {
			id: params.id,
			...(body.epic ? { epic: body.epic } : {})
		});
		return { ...data, log: warnings };
	});
