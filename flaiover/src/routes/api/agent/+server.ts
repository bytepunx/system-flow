import { agent } from '$lib/server/agent';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/**
 * GET: whether a flai on the host has this dashboard connected (ADR-0029), and when it has, what it
 * says about the project: a round trip over the channel, so "connected" means it answers.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		const status = agent().status();
		if (!status.connected) return status;
		try {
			return { ...status, info: await agent().ask('project.info', {}, 2000) };
		} catch (err) {
			return { ...status, error: err instanceof Error ? err.message : String(err) };
		}
	});
