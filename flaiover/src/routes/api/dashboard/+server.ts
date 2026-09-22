import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** What flai dashboard status --json reports (S-0081), asked as dashboard.status, a read. */
export type DashboardStatus = {
	container: string;
	running: boolean;
	image?: string;
	url?: string;
	serves?: string[];
};

const ACTIONS = ['check', 'restart', 'upgrade', 'stop'] as const;
type Action = (typeof ACTIONS)[number];

/**
 * GET: what the dashboard container is running, and whether the operator has enabled the
 * dashboard host action (S-0081), which gates restart, upgrade, and stop. No flai connected reads
 * as not running, not an error: the host flai banner already says why.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		try {
			const { data } = await repo().run<DashboardStatus>(
				'dashboard.status',
				{},
				{ timeoutMs: 15000 }
			);
			return { ...data, dashboard_enabled: await dashboardEnabled() };
		} catch (e) {
			if (e instanceof RepoError && e.status === 503) {
				return { container: 'flaiover', running: false, dashboard_enabled: false };
			}
			throw e;
		}
	});

async function dashboardEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.dashboard === true;
	} catch {
		return false;
	}
}

/**
 * POST {action}: check reports whether a newer image is available, changing nothing (a read that
 * pulls, like push.pending's dry-run fetch — S-0081); restart, upgrade, and stop are the
 * dashboard host action, gated the same way push and publish are (flai serve enable dashboard).
 * Real Docker time: restart and upgrade get generous timeouts, upgrade longest of all since it
 * pulls an image and waits for a temporary container to answer healthy before it touches anything
 * running. 403 with what enables it while the operator has not; the running container is never
 * stopped until a new one has proven itself, so a failed upgrade is an ordinary error, not a
 * partial one.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as { action?: string };
		const action = body.action;
		if (!action || !ACTIONS.includes(action as Action)) {
			throw new RepoError(400, `action must be one of ${ACTIONS.join(', ')}`);
		}
		if (action === 'check') {
			const { data } = await repo().run('dashboard.check', {}, { timeoutMs: 120000 });
			return data as Record<string, unknown>;
		}
		const timeoutMs = action === 'upgrade' ? 300000 : 60000;
		const { data, warnings } = await repo().write(`dashboard.${action}`, {}, { timeoutMs });
		return { ...(data as Record<string, unknown>), log: warnings };
	});
