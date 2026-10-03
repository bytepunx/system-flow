import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RemoteTags, Unplanned } from '$lib/publish';
import type { RequestHandler } from './$types';

/**
 * One component's batch: everything accepted and unreleased for it since its last tag, at the
 * highest delivery level among it (S-0087, release.PendingPlan).
 */
export type PendingPlan = {
	component: { name: string; path: string; kind: string };
	level: string;
	from: string;
	to: string;
	tag?: string;
	version?: string;
	items: { id: string; title: string; level: string }[];
	files: string[];
};

/**
 * GET: what publishing now would release (flai's release.pending, `flai release --pending
 * --dry-run`), and whether the operator has enabled the push host action, which gates publishing
 * (ADR-0031, S-0078, ADR-0067). `remote` is how the clone's release tags and branch stand against
 * its remote's when that has something to say (S-0174, S-0195), and `unplanned` the accepted items
 * no plan covers (I-0024). No flai connected reads as nothing pending, not an error: the host flai banner
 * already says why.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		try {
			const { data } = await repo().run<{
				plans?: PendingPlan[];
				remote?: RemoteTags;
				unplanned?: Unplanned[];
			}>('publish.preview');
			return {
				plans: data?.plans ?? [],
				remote: data?.remote ?? null,
				unplanned: data?.unplanned ?? [],
				push_enabled: await pushEnabled()
			};
		} catch (e) {
			if (e instanceof RepoError && e.status === 503)
				return { plans: [], remote: null, unplanned: [], push_enabled: false };
			throw e;
		}
	});

async function pushEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.push === true;
	} catch {
		return false;
	}
}

/**
 * POST: publish everything accumulated (flai's `publish.run`, the push host action): applies,
 * tags, and pushes every pending component's release together (flai's own `flai release
 * --pending`). 403 with what enables it while the operator has not; 409 with the reason when the
 * remote has moved. Nothing here can enable it: that is `flai serve enable push`, in a shell on
 * the host.
 */
export const POST: RequestHandler = () =>
	respond(async () => {
		const { data, warnings } = await repo().write<{
			plans?: PendingPlan[];
			tags?: string[];
			pushed?: boolean;
			published?: string[];
		}>('publish.run', {}, { timeoutMs: 300000 });
		return { ...data, log: warnings };
	});
