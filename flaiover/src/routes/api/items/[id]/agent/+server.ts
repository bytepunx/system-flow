import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

const ACTIONS = ['start', 'restart'] as const;
type Action = (typeof ACTIONS)[number];

/**
 * POST {action}: the agent host action. start (flai's agent.start, S-0115) has
 * a ready story's agent started now, whatever flai serve's own rules say about
 * when; restart (agent.restart, S-0116) has a new agent, in a new session,
 * started for a story in ready or in progress whose agent dropped or failed.
 * flai judges whether it may and says why not, as a 400; it is 403 with what
 * enables it while the operator has not (flai serve enable agent).
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { action?: string };
		const action = body.action;
		if (!action || !ACTIONS.includes(action as Action)) {
			throw new RepoError(400, `action must be one of ${ACTIONS.join(', ')}`);
		}
		const { data } = await repo().write(`agent.${action}`, { id: params.id });
		return data as Record<string, unknown>;
	});
