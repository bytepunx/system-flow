import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { AgentStreamRead } from '$lib/activity';

/**
 * GET ?after=N: what the newest agent flai serve started for a story said and did, read by flai on
 * the host from the log it gave the agent, from byte offset N, or the log's tail without one
 * (S-0142). Asked fresh each time: the log grows on the host and no file of the project changes
 * with it. 404 when flai started no agent for the story.
 *
 * GET ?plan&after=N: the same of the newest planner flai serve started for the epic or story the path
 * names (S-0259), whatever value plan has; 404 when flai started no planner for the item.
 */
export const GET: RequestHandler = ({ params, url }) =>
	respond(async () => {
		const raw = url.searchParams.get('after');
		const after = raw === null ? undefined : Math.max(0, Math.floor(Number(raw)) || 0);
		const subject = url.searchParams.has('plan') ? { plan: params.story } : { story: params.story };
		return repo().ask<AgentStreamRead>('agent.stream', {
			...subject,
			...(after === undefined ? {} : { after })
		});
	});
