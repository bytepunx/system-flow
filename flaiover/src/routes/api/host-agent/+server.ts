import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { HostAgent } from '$lib/activity';

/**
 * GET: what flai on the host knows of agents it starts when a story becomes ready (S-0079): whether
 * the operator enabled that for this project, the name of the command, what runs, the last start or
 * failure, why a ready story waits, and since S-0104 what each story's agent is doing (`stories`). It is asked of flai each time: the state lies on the host,
 * and no file of the project changes with it. Read-only by design: the dashboard cannot stop or
 * configure an agent, and starts one only through the story page's Restart agent (S-0116,
 * ../items/[id]/agent).
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		try {
			return await repo().ask<HostAgent>('agent.status');
		} catch (e) {
			// no flai connected is said by the banner; here it only means nothing is known
			if (e instanceof RepoError && e.status === 503) return { enabled: false };
			throw e;
		}
	});
