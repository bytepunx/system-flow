import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/** POST { text } */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { text?: string };
		if (!body.text?.trim()) throw new RepoError(400, 'text is required');
		const { data, warnings } = await flai<Record<string, unknown>>(repo().root, [
			'thread',
			'reply',
			params.id,
			'--by',
			await designer(),
			body.text.trim()
		]);
		return { ...data, warnings };
	});
