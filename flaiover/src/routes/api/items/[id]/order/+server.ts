import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { before | after | top | bottom }: place a story in the pull order (flai's item.order on the
 * host, S-0057, S-0075). Exactly one of them; flai refuses anything else, and an ID that is not one.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		const { data } = await repo().write<Record<string, unknown>>('item.order', {
			id: params.id,
			before: typeof body.before === 'string' ? body.before : undefined,
			after: typeof body.after === 'string' ? body.after : undefined,
			top: body.top === true,
			bottom: body.bottom === true
		});
		return data;
	});
