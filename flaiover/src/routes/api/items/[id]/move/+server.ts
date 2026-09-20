import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import { designer } from '$lib/server/threads';
import type { RequestHandler } from './$types';

/**
 * POST { to, reason?, by?, include_uncommitted?, dry_run? }. include_uncommitted is the designer's choice
 * in the acceptance confirmation to put uncommitted files outside wip into the acceptance commit; it
 * becomes flai's --yes, and only for a move to done (S-0051). dry_run, only for a move to cancelled,
 * changes nothing and answers with what the cancellation would take with it (S-0070, ADR-0028).
 */
type MoveBody = {
	to: string;
	reason?: string;
	by?: string;
	include_uncommitted?: boolean;
	dry_run?: boolean;
};

/** The reason a preview runs with when the designer has not typed one yet: flai validates the move first. */
const PREVIEW_REASON = 'preview';

/** The flai arguments for a move. Exported for tests. */
export function _moveArgs(id: string, body: MoveBody): string[] {
	const args = ['move', id, body.to];
	const preview = body.dry_run === true && body.to === 'cancelled';
	if (body.reason || preview) args.push('--reason', body.reason || PREVIEW_REASON);
	if (preview) args.push('--dry-run');
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
			dry_run?: boolean;
		};
		if (!body.to) throw new RepoError(400, 'to is required');
		// A move made on the board is the designer's, not the dashboard's: name them, as threads do.
		// flai would otherwise record FLAI_AGENT, which here is "flaiover" (S-0058).
		const args = _moveArgs(params.id, { ...(body as MoveBody), by: body.by || (await designer()) });
		const { data, warnings } = await flai<{
			id: string;
			status: string;
			warnings: string[];
			/** A move to cancelled: what went, or under dry_run would go, with it, and what stays on disk. */
			cancelled?: { id: string; type: string; title: string; from: string }[];
			left_behind?: { id: string; narrative?: string; branch?: string; worktree?: string }[];
		}>(repo().root, args);
		return { ...data, warnings: [...(data?.warnings ?? []), ...warnings] };
	});
