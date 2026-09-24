import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** One process flai host keeps (S-0106): serve, or a project's MCP server with its root. */
export type HostChild = {
	name: string;
	root?: string;
	state: string;
	pid?: number;
	since?: string;
	version?: string;
	restarts: number;
	last_exit?: string;
	last_error?: string;
};

/** The running host: internal/host.Status (S-0106). */
export type HostStatus = {
	pid: number;
	version: string;
	started: string;
	updated?: string;
	addr?: string;
	children: HostChild[];
};

/** What flai host status --json reports (S-0106), asked as host.status, a read. */
export type HostReport = {
	running: boolean;
	status?: HostStatus;
	/** A host that holds the machine's address for another config file. */
	elsewhere?: { pid: number; version: string; config?: string };
};

/** What GET answers: the host as it stands, or why there is none to show. */
export type HostView =
	| (HostStatus & { running: true; host_enabled: boolean })
	| { running: false; reason: string; host_enabled: boolean };

const ACTIONS = ['check', 'start', 'stop', 'restart', 'upgrade'] as const;
type Action = (typeof ACTIONS)[number];
const PROCESSES = ['serve', 'mcp', 'all'] as const;

/**
 * GET: the host and the processes it keeps, and whether the operator has enabled the host action
 * (S-0107), which gates start, stop, restart, and upgrade. No flai connected, or no host running
 * for flai serve's config, reads as not running with the reason, not an error: the panel says
 * how to start one.
 */
export const GET: RequestHandler = () =>
	respond(async (): Promise<HostView> => {
		let report: HostReport;
		try {
			({ data: report } = await repo().run<HostReport>('host.status', {}, { timeoutMs: 15000 }));
		} catch (e) {
			if (e instanceof RepoError && e.status === 503) {
				return { running: false, reason: 'no host flai is connected', host_enabled: false };
			}
			throw e;
		}
		const host_enabled = await hostEnabled();
		if (report.running && report.status) return { ...report.status, running: true, host_enabled };
		const reason = report.elsewhere
			? `the machine's flai host (pid ${report.elsewhere.pid}) runs for another config, ${report.elsewhere.config ?? 'unnamed'}`
			: 'flai host is not running';
		return { running: false, reason, host_enabled };
	});

async function hostEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.host === true;
	} catch {
		return false;
	}
}

/**
 * POST {action, process}: check asks whether a newer flai is released, changing nothing, a read;
 * start, stop, and restart of serve, the MCP servers, or all (host.process), and upgrade, are the
 * host action, gated the way the dashboard action is (flai serve enable host). Stopping or restarting serve,
 * and an upgrade, end the very connection this request came on: the page expects that and polls
 * GET until it answers again. Upgrade gets the longest timeout, since it downloads a release.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as {
			action?: string;
			process?: string;
		};
		const action = body.action;
		if (!action || !ACTIONS.includes(action as Action)) {
			throw new RepoError(400, `action must be one of ${ACTIONS.join(', ')}`);
		}
		if (action === 'check') {
			const { data } = await repo().run('host.check', {}, { timeoutMs: 120000 });
			return data as Record<string, unknown>;
		}
		if (action === 'upgrade') {
			const { data, warnings } = await repo().write('host.upgrade', {}, { timeoutMs: 360000 });
			return { ...(data as Record<string, unknown>), log: warnings };
		}
		const process = body.process;
		if (!process || !PROCESSES.includes(process as (typeof PROCESSES)[number])) {
			throw new RepoError(400, `process must be one of ${PROCESSES.join(', ')}`);
		}
		const { data, warnings } = await repo().write(
			'host.process',
			{ process, action },
			{ timeoutMs: 60000 }
		);
		return { ...(data as Record<string, unknown>), log: warnings };
	});
