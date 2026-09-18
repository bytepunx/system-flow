import { repo } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/** POST { reason? } */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { reason?: string };
		const args = ['thread', 'resolve', params.id, '--by', await designer()];
		if (body.reason?.trim()) args.push('--reason', body.reason.trim());
		const { data, warnings } = await flai<Record<string, unknown>>(repo().root, args);
		return { ...data, warnings };
	});
