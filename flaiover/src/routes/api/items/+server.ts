import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET /api/items?type=story&status=review&archived=false. flai is asked for the type and state, and
 * for the archive unless archived=false (S-0162): a page that does not show the archive says so.
 */
export const GET: RequestHandler = ({ url }) =>
	respond(async () => {
		const archived = url.searchParams.get('archived');
		let items = await repo().items({
			type: url.searchParams.get('type') ?? undefined,
			status: url.searchParams.get('status') ?? undefined,
			archive: archived !== 'false'
		});
		if (archived === 'true') items = items.filter((it) => it.archived);
		return items.map((it) => {
			const { body, ...rest } = it;
			void body;
			return rest;
		});
	});

/**
 * POST { type, title, nature?, parent?, tags?, topics?, touches?, agent?, cost_of_delay?, body }: create an epic or a story with the
 * designer's markdown as its body (S-0059; flai's item.new on the host, S-0075). flai does it as one
 * step (ADR-0016): the item from the project's template with the next ID, linked into its parent,
 * checked with it in place, committed on its own, owned by the manifest's owner. 400 when an
 * argument is not what it should be; a story's agent ({ harness?, model?, config?, roles? }, S-0103, S-0189) goes over
 * the project's default, which fills in the rest, role by role; its cost of delay inputs
 * ({ revenue_per_week?, penalty_per_week?, time_lost_per_cycle? } as text, S-0204) are set by the operator; 422 { error, findings } when the check refuses it, and then
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
			topics: body.topics,
			touches: body.touches,
			agent: body.agent,
			cost_of_delay: body.cost_of_delay,
			body: body.body
		});
		return { ...data, log: warnings };
	});
