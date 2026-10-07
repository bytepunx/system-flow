// The orchestrator page's read and its stop and start (S-0228): whether the orchestrate host action
// is on, whether the operator holds the orchestrator stopped, its activity document, and its runs,
// each asked of flai on the host.

import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { StrategicStatus } from '$lib/server/agent';
import {
	activityDocument,
	currentRun,
	ORCHESTRATOR_ACTIONS,
	type OrchestratorAction,
	type OrchestratorView,
	type RawActivityDocument
} from '$lib/strategic';
import type { RequestHandler } from './$types';

/** An answer of flai's, or null when flai cannot be asked: the page then reads it as off and empty. */
const orNull = <T>(asked: Promise<T>): Promise<T | null> => asked.catch(() => null);

/**
 * GET: whether the orchestrate host action is on (project.info's host_actions.orchestrate), the
 * orchestrator's activity document, whose entries are its decisions, each summary with its reason,
 * and its runs from agent.status, newest first, with whether the operator holds it stopped. A flai that
 * cannot be asked reads as the action off, not held, and no runs; the activity document is the
 * page's subject, so failing to read it is the answer's error.
 */
export const GET: RequestHandler = () =>
	respond(async (): Promise<OrchestratorView> => {
		const [info, status, activity] = await Promise.all([
			orNull(repo().ask<{ host_actions?: Record<string, boolean> }>('project.info')),
			orNull(repo().ask<StrategicStatus>('agent.status')),
			repo()
				.ask<RawActivityDocument<'orchestrator'>>('activity.document', { kind: 'orchestrator' })
				.then(activityDocument)
		]);
		const newest = status?.state?.orchestrator ?? null;
		const runs = [...(newest ? [newest] : []), ...(status?.state?.past_orchestrators ?? [])];
		return {
			enabled: info?.host_actions?.orchestrate === true,
			held: newest?.held === true,
			activity,
			run: currentRun(runs),
			runs
		};
	});

/**
 * POST { action }: stop is flai's orchestrate.stop, which ends the orchestrator's run and holds it
 * stopped while the orchestrate action stays on; start is orchestrate.start, which lifts the hold.
 * Each answers what flai answered; it is 403 with what enables it while the operator has not (flai
 * serve enable orchestrate), and 400 for an action that is neither.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => null)) as { action?: unknown } | null;
		const action = body?.action;
		if (!ORCHESTRATOR_ACTIONS.includes(action as OrchestratorAction))
			throw new RepoError(
				400,
				`action must be one of ${ORCHESTRATOR_ACTIONS.join(', ')}, not ${JSON.stringify(action ?? null)}`
			);
		const { data } = await repo().write(`orchestrate.${action}`);
		return (data ?? {}) as Record<string, unknown>;
	});
