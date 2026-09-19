import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/**
 * POST { to, reason?, by?, include_uncommitted? }. include_uncommitted is the designer's choice in the
 * acceptance confirmation to put uncommitted files outside wip into the acceptance commit; it becomes
 * flai's --yes, and only for a move to done (S-0051).
 */
type MoveBody = { to: string; reason?: string; by?: string; include_uncommitted?: boolean };

/** The flai arguments for a move. Exported for tests. */
export function _moveArgs(id: string, body: MoveBody): string[] {
	const args = ['move', id, body.to];
	if (body.reason) args.push('--reason', body.reason);
	if (body.by) args.push('--by', body.by);
	if (body.include_uncommitted === true && body.to === 'done') args.push('--yes');
	return args;
}

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			to?: string;
			reason?: string;
			by?: string;
			include_uncommitted?: boolean;
		};
		if (!body.to) throw new RepoError(400, 'to is required');
		// A move made on the board is the designer's, not the dashboard's: name them, as threads do.
		// flai would otherwise record FLAI_AGENT, which here is "flaiover" (S-0058).
		const args = _moveArgs(params.id, { ...(body as MoveBody), by: body.by || (await designer()) });
		const { data, warnings } = await flai<{ id: string; status: string; warnings: string[] }>(
			repo().root,
			args
		);
		return { ...data, warnings: [...(data?.warnings ?? []), ...warnings] };
	});
