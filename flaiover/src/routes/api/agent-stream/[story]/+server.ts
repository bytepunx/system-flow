import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';
import type { AgentStreamRead } from '$lib/activity';
import { streamRole } from '$lib/strategic';

/**
 * GET ?after=N: what the newest agent flai serve started for a story said and did, read by flai on
 * the host from the log it gave the agent, from byte offset N, or the log's tail without one
 * (S-0142). Asked fresh each time: the log grows on the host and no file of the project changes
 * with it. 404 when flai started no agent for the story.
 *
 * GET ?plan&after=N: the same of the newest planner flai serve started for the epic or story the path
 * names (S-0259), whatever value plan has; 404 when flai started no planner for the item.
 *
 * GET ?orchestrator&after=N: the same of the project's newest orchestrator run (S-0218), whatever the
 * path names; 404 when flai started no orchestrator for the project.
 *
 * GET ?role=<role>&after=N: the same of the project's newest run of that strategic role, orchestrate
 * or analyze by its verb (S-0228), or any role the project runs by its name (S-0420), whatever the
 * path names; 404 when flai started none. A role that is not one word in kebab case is 400 here, and
 * one the project does not run is flai's 400.
 */
export const GET: RequestHandler = ({ params, url }) =>
	respond(async () => {
		const raw = url.searchParams.get('after');
		const after = raw === null ? undefined : Math.max(0, Math.floor(Number(raw)) || 0);
		const role = url.searchParams.get('role');
		if (role !== null && !streamRole(role))
			throw new RepoError(
				400,
				`role must be orchestrate, analyze, or a role's name in kebab case, such as reviewer, not ${JSON.stringify(role)}`
			);
		const subject =
			role !== null
				? { role }
				: url.searchParams.has('orchestrator')
					? { orchestrator: true }
					: url.searchParams.has('plan')
						? { plan: params.story }
						: { story: params.story };
		return repo().ask<AgentStreamRead>('agent.stream', {
			...subject,
			...(after === undefined ? {} : { after })
		});
	});
