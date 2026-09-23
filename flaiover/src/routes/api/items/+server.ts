import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** GET /api/items?type=story&status=review&archived=false */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const type = url.searchParams.get('type');
		const status = url.searchParams.get('status');
		const archived = url.searchParams.get('archived');
		let items = await repo().items();
		if (type) items = items.filter((it) => it.type === type);
		if (status) items = items.filter((it) => it.status === status);
		if (archived === 'true' || archived === 'false')
			items = items.filter((it) => it.archived === (archived === 'true'));
		return items.map((it) => {
			const { body, ...rest } = it;
			void body;
			return rest;
		});
	});

/**
 * POST { type, title, nature?, parent?, tags?, touches?, agent?, body }: create an epic or a story with the
 * designer's markdown as its body (S-0059; flai's item.new on the host, S-0075). flai does it as one
 * step (ADR-0016): the item from the project's template with the next ID, linked into its parent,
 * checked with it in place, committed on its own, owned by the manifest's owner. 400 when an
 * argument is not what it should be; a story's agent ({ harness?, model?, config? }, S-0103) goes over the
 * project's default, which fills in the rest; 422 { error, findings } when the check refuses it, and then
 * nothing was created.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		const { data, warnings } = await repo().write<Record<string, unknown>>('item.new', {
			type: body.type,
			title: body.title,
			nature: body.nature,
			parent: body.parent,
			tags: body.tags,
			touches: body.touches,
			agent: body.agent,
			body: body.body
		});
		return { ...data, log: warnings };
	});
