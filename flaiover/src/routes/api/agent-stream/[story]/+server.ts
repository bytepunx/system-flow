import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { AgentStreamRead } from '$lib/activity';

/**
 * GET ?after=N: what the newest agent flai serve started for a story said and did, read by flai on
 * the host from the log it gave the agent, from byte offset N, or the log's tail without one
 * (S-0142). Asked fresh each time: the log grows on the host and no file of the project changes
 * with it. 404 when flai started no agent for the story.
 */
export const GET: RequestHandler = ({ params, url }) =>
	respond(async () => {
		const raw = url.searchParams.get('after');
		const after = raw === null ? undefined : Math.max(0, Math.floor(Number(raw)) || 0);
		return repo().ask<AgentStreamRead>('agent.stream', {
			story: params.story,
			...(after === undefined ? {} : { after })
		});
	});
