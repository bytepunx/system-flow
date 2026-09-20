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
		if (status.missing?.length)
			return {
				...status,
				error: `flai ${status.flai} on the host is older than this dashboard and lacks ${status.missing.length} of the things it asks for (${status.missing.slice(0, 3).join(', ')}${status.missing.length > 3 ? ', …' : ''}); upgrade flai on the host, then run flai serve stop and flai dashboard`
			};
		try {
			return { ...status, info: await agent().ask('project.info', {}, 2000) };
		} catch (err) {
			return { ...status, error: err instanceof Error ? err.message : String(err) };
		}
	});
