import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { HostAgent } from '$lib/activity';

/**
 * GET: what flai on the host knows of agents it starts when a story becomes ready (S-0079): whether
 * the operator enabled that for this project, the name of the command, what runs, the last start or
 * failure, why a ready story waits, and since S-0104 what each story's agent is doing (`stories`). It is asked of flai each time: the state lies on the host,
 * and no file of the project changes with it. Since S-0263 it says too whether the plan host action
 * is on (project.info's host_actions.plan, `plan_enabled`), for the board's card menu; a flai that
 * cannot say reads as off. Read-only by design: the dashboard cannot stop or configure an agent,
 * and starts one only through the story page's Start agent (S-0115) and Retry (S-0116, S-0118,
 * ../items/[id]/agent), and the planner through Plan (S-0208, ../items/[id]/plan).
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		try {
			const [status, plan_enabled] = await Promise.all([
				repo().ask<HostAgent>('agent.status'),
				planEnabled()
			]);
			return { ...status, plan_enabled };
		} catch (e) {
			// no flai connected is said by the banner; here it only means nothing is known
			if (e instanceof RepoError && e.status === 503) return { enabled: false };
			throw e;
		}
	});

async function planEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.plan === true;
	} catch {
		return false;
	}
}
