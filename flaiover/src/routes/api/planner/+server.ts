// The planner page's read (S-0259): whether the plan host action is on, the planner's activity
// document, and its runs, each asked of flai on the host.

import { repo } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { HostAgent, PlanRun } from '$lib/activity';
import {
	plannerActivity,
	runsNewestFirst,
	type PlannerDocument,
	type PlannerView
} from '$lib/planner';
import type { RequestHandler } from './$types';

/**
 * GET: whether the plan host action is on (project.info's host_actions.plan), the planner's
 * activity document, and each item's newest planner run from agent.status's plans, newest first.
 * A flai that cannot be asked reads as the action off and no runs; the activity document is the
 * page's subject, so failing to read it is the answer's error.
 */
export const GET: RequestHandler = () =>
	respond(async (): Promise<PlannerView> => {
		const [plan_enabled, activity, runs] = await Promise.all([
			planEnabled(),
			repo().ask<PlannerDocument>('activity.document', { kind: 'planner' }).then(plannerActivity),
			planRuns()
		]);
		return { plan_enabled, activity, runs };
	});

async function planEnabled(): Promise<boolean> {
	try {
		const info = await repo().ask<{ host_actions?: Record<string, boolean> }>('project.info');
		return info.host_actions?.plan === true;
	} catch {
		return false;
	}
}

async function planRuns(): Promise<PlanRun[]> {
	try {
		const status = await repo().ask<HostAgent>('agent.status');
		return runsNewestFirst(Object.values(status.state?.plans ?? {}));
	} catch {
		return [];
	}
}
