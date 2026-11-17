import { board } from '$lib/server/board';
import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import { isAccepted, isCancelled } from '$lib/workflow';
import type { RequestHandler } from './$types';

/**
 * POST { to, reason?, dry_run? }: flai's item.move on the host (S-0075). flai validates the
 * arguments, builds the command, and records the move as the manifest's owner: a move made on the
 * board is the designer's (S-0058), and the dashboard does not get to name anyone. What a move is
 * comes from the kind of its lane on the board's flow, never its name (S-0758, I-0399). A move to
 * the lane of kind accepted is an acceptance; the main checkout's uncommitted files are flai's to
 * commit or leave (ADR-0248), so the body carries no choice about them, and an include_uncommitted
 * an older page sends is ignored. dry_run, only for a move to the lane of kind cancelled, changes
 * nothing and answers with what the cancellation would take with it (S-0070, ADR-0028). An
 * acceptance that flai refuses answers 422 with the refused object's findings, each a file to open
 * on its own checkout and correct (ADR-0140), its undoes_main (ADR-0155), or its send back for a
 * rebase that stopped, sent_back with reason, conflicts, and thread (ADR-0248), beside error
 * (S-0465).
 */
type Moved = {
	id: string;
	status: string;
	warnings?: string[];
	/** A cancellation: what went, or under dry_run would go, with it, and what stays on disk. */
	cancelled?: { id: string; type: string; title: string; from: string }[];
	left_behind?: { id: string; narrative?: string; branch?: string; worktree?: string }[];
};

export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			to?: string;
			reason?: string;
			dry_run?: boolean;
		};
		if (!body.to) throw new RepoError(400, 'to is required');
		const flow = await board(repo());
		if (body.dry_run === true) {
			// A dry run never moves: one to a lane flai does not call cancelled is refused, as on a
			// flai too old to answer the board's lanes.
			if (!isCancelled(flow, body.to)) {
				throw new RepoError(400, `dry_run is only for a move to the cancelled lane: ${body.to}`);
			}
			const { data } = await repo().run<Moved>('item.move.preview', { id: params.id });
			return data;
		}
		const { data, warnings } = await repo().write<Moved>(
			'item.move',
			{
				id: params.id,
				to: body.to,
				reason: body.reason
			},
			// An acceptance takes a rebase, a merge, and an archive, which take their time.
			{ timeoutMs: isAccepted(flow, body.to) ? 600000 : 60000 }
		);
		return { ...data, warnings: [...(data?.warnings ?? []), ...warnings] };
	});
