import { repo, RepoError } from '$lib/server/repo';
import { flai } from '$lib/server/flai';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: whether the clone holds an acceptance its remote has not got (S-0063, ADR-0026). The answer
 * is flai's, asked offline of git (`flai push --pending --dry-run`): the dashboard needs no
 * credential to know, and never pushes from here. `unpushed` is null when there is nothing, when
 * what is ahead holds no acceptance, and when flai is not available to ask.
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
			const { data } = await flai<{ unpushed?: Unpushed }>(repo().root, [
				'push',
				'--pending',
				'--dry-run'
			]);
			return { unpushed: data?.unpushed ?? null };
		} catch (e) {
			// No flai (a read-only dashboard) is not an error worth a banner; a diverged clone is
			// something the operator should see.
			if (e instanceof RepoError && e.status === 503) return { unpushed: null };
			if (e instanceof RepoError) return { unpushed: null, problem: e.message };
			throw e;
		}
	});
