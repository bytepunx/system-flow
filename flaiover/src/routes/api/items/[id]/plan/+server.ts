import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { HostAgent, PlanRun } from '$lib/activity';
import type { RequestHandler } from './$types';

/**
 * GET: whether the plan host action is on for the project (project.info's
 * host_actions.plan, S-0208), asked of flai each time like checks_enabled,
 * and the item's newest planner run from agent.status's plans, or null. A
 * flai that cannot be asked reads as the action off and no run.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async () => {
		const [plan_enabled, run] = await Promise.all([planEnabled(), newestRun(params.id)]);
		return { plan_enabled, run };
	});

async function planEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.plan === true;
	} catch {
		return false;
	}
}

async function newestRun(id: string): Promise<PlanRun | null> {
	try {
		const status = await repo().ask<HostAgent>('agent.status');
		return status.state?.plans?.[id] ?? null;
	} catch {
		return null;
	}
}

/**
 * POST: the plan host action (flai's plan.run, S-0208): the planner started
 * now for the epic or story, which answers the run. flai judges whether it
 * may and says why not, as a 400 (a planner already running for the item,
 * an item done, cancelled, or archived); it is 403 with what enables it
 * while the operator has not (flai serve enable plan).
 */
export const POST: RequestHandler = ({ params }) =>
	respond(async () => {
		const { data } = await repo().write('plan.run', { id: params.id });
		return data as Record<string, unknown>;
	});
