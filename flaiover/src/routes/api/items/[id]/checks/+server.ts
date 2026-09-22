import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: the current or last checks run for a story (flai's checks.status, a
 * read, S-0082): which check is running, pass/fail so far, outcome and
 * duration once it has ended. A story with no run yet answers { running:
 * false }, not an error. Carries checks_enabled alongside it (project.info's
 * host_actions.checks), asked of flai each time like push's own pushEnabled:
 * a setting of the host, not of the project, so no change notice would ever
 * clear a remembered answer.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const [status, checks_enabled] = await Promise.all([
			repo().run('checks.status', { id: params.id }),
			checksEnabled()
		]);
		return { ...(status.data as Record<string, unknown>), checks_enabled };
	});

async function checksEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.checks === true;
	} catch {
		return false;
	}
}

const ACTIONS = ['run', 'cancel'] as const;
type Action = (typeof ACTIONS)[number];

/**
 * POST {action}: run is the checks host action (flai's checks.run): every
 * named check, in order, in the story's worktree, stopping at the first to
 * fail, bounded by the operator's own time limit. cancel stops an active
 * run (terminate, then kill). Both are 403 with what enables it while the
 * operator has not (flai serve enable checks). run can take real minutes:
 * ReviewChecks polls GET and the tail route below rather than waiting on
 * this response's own timeout.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { action?: string };
		const action = body.action;
		if (!action || !ACTIONS.includes(action as Action)) {
			throw new RepoError(400, `action must be one of ${ACTIONS.join(', ')}`);
		}
		const timeoutMs = action === 'run' ? 2 * 60 * 60 * 1000 : 30000;
		const { data, warnings } = await repo().write(
			`checks.${action}`,
			{ id: params.id },
			{ timeoutMs }
		);
		return { ...(data as Record<string, unknown>), log: warnings };
	});
