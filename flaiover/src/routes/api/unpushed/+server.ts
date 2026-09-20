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
			return { unpushed: data?.unpushed ?? null, push_enabled: await pushEnabled() };
		} catch (e) {
			// No flai connected is said by the banner; a diverged clone is something the operator should see.
			if (e instanceof RepoError && e.status === 503)
				return { unpushed: null, push_enabled: false };
			if (e instanceof RepoError)
				return { unpushed: null, problem: e.message, push_enabled: await pushEnabled() };
			throw e;
		}
	});

/**
 * Whether the operator enabled the push action for this project (S-0078). Asked of flai each time,
 * not remembered: it is a setting of the host, no file of the project changes with it, and so no
 * change notice would ever clear a remembered answer.
 */
async function pushEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.push === true;
	} catch {
		return false;
	}
}

/**
 * POST: push what was accepted and not pushed, and publish the template where its version moved
 * (flai's push.run, a host action). flai does it on the host as the operator, never forced. 403 with
 * what enables it while the operator has not; 409 with the reason when the remote has moved. Nothing
 * here can enable it: that is `flai serve enable push`, in a shell on the host.
 */
export const POST: RequestHandler = () =>
	respond(async () => {
		const { data, warnings } = await repo().write<Record<string, unknown>>(
			'push.run',
			{},
			{ timeoutMs: 300000 }
		);
		return { ...data, log: warnings };
	});
