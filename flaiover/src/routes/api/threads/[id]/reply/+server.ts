import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** POST { text }: reply to a thread as the designer (flai's thread.reply on the host). */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { text?: string };
		if (!body.text?.trim()) throw new RepoError(400, 'text is required');
		const { data, warnings } = await repo().write<Record<string, unknown>>('thread.reply', {
			id: params.id,
			text: body.text
		});
		return { ...data, warnings };
	});
