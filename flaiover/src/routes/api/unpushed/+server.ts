import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: whether the clone holds an acceptance its remote has not got (S-0063, ADR-0026). The answer
 * is flai's on the host, asked offline of git (push.pending, which is `flai push --pending --dry-run`):
 * nothing is pushed from here. `unpushed` is null when there is nothing, when what is ahead holds no
 * acceptance, and when no flai is connected to ask.
 */
export type Unpushed = {
	branch: string;
	upstream: string;
	remote: string;
	commits: number;
	acceptances: string[];
	tags: string[];
	behind?: number;
	command: string;
};

export const GET: RequestHandler = () =>
	respond(async () => {
		try {
			const { data } = await repo().run<{ unpushed?: Unpushed }>('push.pending');
			return { unpushed: data?.unpushed ?? null };
		} catch (e) {
			// No flai connected is said by the banner; a diverged clone is something the operator should see.
			if (e instanceof RepoError && e.status === 503) return { unpushed: null };
			if (e instanceof RepoError) return { unpushed: null, problem: e.message };
			throw e;
		}
	});
