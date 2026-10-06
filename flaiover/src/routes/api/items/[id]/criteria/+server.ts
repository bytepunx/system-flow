import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { hash, tick?: number[], untick?: number[] }: tick or untick an item's acceptance criteria
 * by their numbers, from 1 in the order they appear, as the designer (flai's item.criteria on the
 * host, which is `flai criteria tick` or `untick`, S-0282): those boxes change and no other byte of
 * the body, committed, and the dashboard does not name anyone. One call ticks or unticks, not
 * both. The hash is the one the item was loaded with: 409 { error, current, hash } when the item
 * changed after it was read; 422 { error, ... } when flai refuses (an archived or closed item, or
 * what flai check finds); 400 for a number there is not, or a body that is not this shape.
 * Answers flai's edit result with `criteria` ([{ n, text, ticked }]) as they now are.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
		if (typeof body.hash !== 'string' || !body.hash)
			throw new RepoError(400, 'hash is required: the one the item was loaded with');
		const change: Record<string, unknown> = { id: params.id, hash: body.hash };
		for (const f of ['tick', 'untick'] as const) {
			const v = body[f];
			if (v === undefined || v === null) continue;
			if (!Array.isArray(v) || !v.every((n) => Number.isInteger(n) && n >= 1))
				throw new RepoError(400, `${f} is a list of criterion numbers, from 1`);
			if (v.length) change[f] = v;
		}
		if (!change.tick && !change.untick)
			throw new RepoError(400, 'name the criteria to tick or untick by their numbers');
		const { data, warnings } = await repo().write<Record<string, unknown>>('item.criteria', change);
		return { ...data, log: warnings };
	});
