import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * POST { to, reason?, include_uncommitted?, dry_run? }: flai's item.move on the host (S-0075). flai
 * validates the arguments, builds the command, and records the move as the manifest's owner: a move
 * made on the board is the designer's (S-0058), and the dashboard does not get to name anyone.
 * include_uncommitted is the designer's choice in the acceptance confirmation to put uncommitted
 * files outside wip into the acceptance commit, and only for a move to done (S-0051). dry_run, only
 * for a move to cancelled, changes nothing and answers with what the cancellation would take with it
 * (S-0070, ADR-0028).
 */
type Moved = {
	id: string;
	status: string;
	warnings?: string[];
	/** A move to cancelled: what went, or under dry_run would go, with it, and what stays on disk. */
	cancelled?: { id: string; type: string; title: string; from: string }[];
	left_behind?: { id: string; narrative?: string; branch?: string; worktree?: string }[];
};

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			to?: string;
			reason?: string;
			include_uncommitted?: boolean;
			dry_run?: boolean;
		};
		if (!body.to) throw new RepoError(400, 'to is required');
		if (body.dry_run === true && body.to === 'cancelled') {
			const { data } = await repo().run<Moved>('item.move.preview', { id: params.id });
			return data;
		}
		const { data, warnings } = await repo().write<Moved>(
			'item.move',
			{
				id: params.id,
				to: body.to,
				reason: body.reason,
				include_uncommitted: body.include_uncommitted === true
			},
			// A story moved to done is accepted: a merge, a release, and a push take their time.
			{ timeoutMs: body.to === 'done' ? 600000 : 60000 }
		);
		return { ...data, warnings: [...(data?.warnings ?? []), ...warnings] };
	});
