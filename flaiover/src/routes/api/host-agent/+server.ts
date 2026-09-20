import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: what flai on the host knows of agents it starts when a story becomes ready (S-0079): whether
 * the operator enabled that for this project, the name of the command, what runs, the last start or
 * failure, and why a ready story waits. It is asked of flai each time: the state lies on the host,
 * and no file of the project changes with it. Read-only by design: the dashboard cannot start, stop,
 * or configure an agent.
 */
export type AgentRun = {
	story: string;
	command: string;
	agent: string;
	pid?: number;
	started: string;
	ended?: string;
	exit?: number;
	error?: string;
	log?: string;
};
export type HostAgent = {
	enabled: boolean;
	state?: { command: string; running?: AgentRun | null; last?: AgentRun | null; waiting?: string };
};

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
