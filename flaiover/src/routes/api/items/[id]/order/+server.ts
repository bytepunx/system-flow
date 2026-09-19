import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { before | after | top | bottom }: place a ready or backlog story in its column's pull
 * order (S-0057). The write is flai order's (ADR-0016): the dashboard never edits board.md, and
 * flai's refusals (an epic, another state, a reference in another column) come back as 400 with
 * its reason.
 */
type OrderBody = { before?: string; after?: string; top?: boolean; bottom?: boolean };

const ID = /^[EST]-\d+$/;

/** The flai arguments for a placement. Exported for tests. */
export function _orderArgs(id: string, body: OrderBody): string[] {
	const given = [
		typeof body.before === 'string' && body.before !== '',
		typeof body.after === 'string' && body.after !== '',
		body.top === true,
		body.bottom === true
	].filter(Boolean).length;
	if (given !== 1)
		throw new RepoError(400, 'say where it goes with exactly one of before, after, top, bottom');
	const ref = body.before || body.after;
	// An ID, never a flag: what follows --before must not be read as an option.
	if (ref !== undefined && !ID.test(ref)) throw new RepoError(400, `${ref} is not an item ID`);
	if (!ID.test(id)) throw new RepoError(400, `${id} is not an item ID`);
	const args = ['order', id];
	if (body.before) args.push('--before', body.before);
	else if (body.after) args.push('--after', body.after);
	else if (body.top === true) args.push('--top');
	else args.push('--bottom');
	return args;
}

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as OrderBody;
		const { data } = await flai<{
			id: string;
			status: string;
			sequence: string[];
			order: string[];
		}>(repo().root, _orderArgs(params.id, body));
		return data;
	});
